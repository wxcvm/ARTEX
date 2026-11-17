#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
bot_hunt.py —— 三引擎检索封装（payout-hunt 技能 · 发现阶段）

同一个意图翻译成 FOFA / Shodan / Censys 三家语法，统一限速、退避、去重、落盘。
**只管找，不管判**；产物交给同目录的 bot_verify.py 做六道 Gate 判定。

模式:
  --keyword            跑内置关键词矩阵（可叠加 --geo / --species）
  --icon-hash N        按 favicon mmh3 反查整套站群（最强）
  --body STR           按页面独有字符串反查（静态资源名 / 接口前缀）
  --query RAW          原样查询，只发给指定引擎，不做翻译
  --from-sample URL    从活样本提指纹，输出 JSON 供上面三种模式复用手抄

示例（在技能目录下执行）:
  # 1) 先取指纹
  python scripts/bot_hunt.py --from-sample https://TARGET -o out/fp.json
  # 2) 用指纹铺站群
  python scripts/bot_hunt.py --engine fofa,shodan --icon-hash 1234567 -o out/squad.csv
  # 3) 关键词铺面
  python scripts/bot_hunt.py --engine fofa --keyword --geo sea --species usdt --pages 2 -o out/breadth.csv
  # 4) 判定
  python scripts/bot_verify.py -i out/breadth.csv -o out/probe.jsonl -c 40

环境变量:
  FOFA_KEY / SHODAN_KEY / CENSYS_ID + CENSYS_SECRET   三家凭据，缺哪家自动跳过哪家
  FOFA_ENDPOINT                                       覆盖 FOFA 端点（默认官方 fofa.info）

依赖: requests（必需）；mmh3（favicon 指纹，缺了相关模式降级为按标题/接口前缀反查）
"""
import argparse, csv, json, os, re, sys, time
from urllib.parse import quote, urlparse

import requests
import urllib3
urllib3.disable_warnings()

UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/124.0 Safari/537.36"

FOFA = os.environ.get("FOFA_ENDPOINT") or "https://fofa.info/api/v1/search/all"
SHODAN = "https://api.shodan.io/shodan/host/search"
CENSYS = "https://search.censys.io/api/v2/hosts/search"

# ---------------- 关键词矩阵（文档 §2.2 / §2.4） ----------------
SPECIES = {
    "brand": ['body="葡京"', 'body="威尼斯"', 'body="太阳城"', 'body="金沙"',
              'body="银河"', 'body="开元"', 'body="美高梅"', 'body="皇冠"'],
    "cat":   ['title="棋牌"', 'title="真人"', 'title="电子游艺"', 'title="体育"',
              'title="电竞"', 'title="彩票"', 'title="捕鱼"'],
    "money": ['body="静态收益"', 'body="动态收益"', 'body="抢单"', 'body="分红"',
              'body="拆分"', 'body="互助"', 'body="量化机器人"', 'body="日化"'],
    "usdt":  ['body="USDT充值"', 'body="TRC20"', 'body="TronLink"', 'body="钱包地址"',
              'body="秒到账"', 'body="归集"'],
    "agent": ['body="代理注册"', 'body="申请代理"', 'body="推广链接"', 'body="返水"',
              'body="佣金比例"'],
    "footer": ['body="技术支持"', 'body="系统源码"'],
}

# 地理：按 IP 落点，不按语言。CN 一律排除
GEO = {
    "sea":  ["SG", "MY", "PH", "KH", "LA", "TH", "VN", "ID"],
    "mena": ["AE", "SA", "QA", "KW", "BH"],
    "eu":   ["NL", "DE", "UA", "RU", "SC"],
    "all":  ["SG", "MY", "PH", "KH", "LA", "TH", "VN", "ID",
             "AE", "SA", "NL", "DE", "UA", "RU", "SC", "PA"],
}

PORTS = '(port="80"||port="443"||port="8080"||port="8443"||port="3000"||port="8888")'


# ---------------- 引擎适配 ----------------
class Engine:
    name = "base"
    qps = 1.0

    def __init__(self, qps):
        self.qps = qps
        self._last = 0.0
        self.calls = 0

    def throttle(self):
        gap = 1.0 / self.qps
        d = time.time() - self._last
        if d < gap:
            time.sleep(gap - d)
        self._last = time.time()

    def search(self, q, page):
        raise NotImplementedError

    def retry(self, fn, *a, **kw):
        """429 / 5xx / 超时退避重试，最多 4 次"""
        for i in range(4):
            try:
                self.throttle()
                self.calls += 1
                return fn(*a, **kw)
            except RateLimited:
                time.sleep(2 ** i * 3)
            except (requests.Timeout, requests.ConnectionError):
                time.sleep(2 ** i)
            except Exception as e:
                if i == 3:
                    print(f"[!] {self.name} 放弃: {e}", file=sys.stderr)
                    return []
                time.sleep(2 ** i)
        return []


class RateLimited(Exception):
    pass


class Fofa(Engine):
    name = "fofa"

    def search(self, q, page):
        key = os.environ.get("FOFA_KEY")
        if not key:
            print("[!] 缺 FOFA_KEY", file=sys.stderr); return []
        p = {"key": key, "qbase64": __import__("base64").b64encode(q.encode()).decode(),
             "fields": "host,ip,port,protocol,title,domain,country,city,server,icon_hash",
             "page": page, "size": 100}
        r = requests.get(FOFA, params=p, timeout=20,
                         headers={"User-Agent": UA})
        if r.status_code == 429:
            raise RateLimited
        j = r.json()
        if j.get("error"):
            msg = str(j.get("errmsg", ""))
            if "429" in msg or "limit" in msg.lower():
                raise RateLimited
            print(f"[!] fofa: {msg}", file=sys.stderr); return []
        return [{
            "host": x.get("host", ""), "ip": x.get("ip", ""), "port": x.get("port", ""),
            "protocol": x.get("protocol", ""), "title": x.get("title", ""),
            "domain": x.get("domain", ""), "country": x.get("country", ""),
            "city": x.get("city", ""), "server": x.get("server", ""),
            "icon_hash": x.get("icon_hash", ""), "source": "fofa", "query": q,
        } for x in j.get("results", [])]


class Shodan(Engine):
    name = "shodan"

    def search(self, q, page):
        key = os.environ.get("SHODAN_KEY")
        if not key:
            print("[!] 缺 SHODAN_KEY", file=sys.stderr); return []
        r = requests.get(SHODAN, params={"key": key, "query": q, "page": page},
                         timeout=20, headers={"User-Agent": UA})
        if r.status_code == 429:
            raise RateLimited
        j = r.json()
        if "error" in j:
            print(f"[!] shodan: {j['error']}", file=sys.stderr); return []
        out = []
        for m in j.get("matches", []):
            host = (m.get("hostnames") or [""])[0] or m.get("ip_str", "")
            out.append({
                "host": host, "ip": m.get("ip_str", ""), "port": m.get("port", ""),
                "protocol": m.get("transport", ""), "title": m.get("title", ""),
                "domain": "", "country": m.get("location", {}).get("country_code", ""),
                "city": m.get("location", {}).get("city", ""),
                "server": (m.get("http", {}) or {}).get("server", ""),
                "icon_hash": (m.get("http", {}) or {}).get("favicon", {}).get("hash", ""),
                "source": "shodan", "query": q,
            })
        return out


class Censys(Engine):
    name = "censys"

    def search(self, q, page):
        i, s = os.environ.get("CENSYS_ID"), os.environ.get("CENSYS_SECRET")
        if not (i and s):
            print("[!] 缺 CENSYS_ID/CENSYS_SECRET", file=sys.stderr); return []
        r = requests.post(CENSYS, json={"query": q, "per_page": 100, "cursor": page or ""},
                          auth=(i, s), timeout=25, headers={"User-Agent": UA})
        if r.status_code == 429:
            raise RateLimited
        j = r.json()
        out = []
        for m in j.get("result", {}).get("hits", []):
            out.append({
                "host": (m.get("name") or ""), "ip": m.get("ip", ""),
                "port": m.get("port", ""), "protocol": "",
                "title": "", "domain": "", "country": "",
                "city": "", "server": "",
                "icon_hash": "", "source": "censys", "query": q,
            })
        return out


ENGINES = {"fofa": Fofa, "shodan": Shodan, "censys": Censys}


# ---------------- 查询翻译 ----------------
def to_fofa(bare, geo):
    cc = "||".join(f'country="{c}"' for c in GEO[geo])
    return f'({bare})&&({cc})&&country!="CN"&&{PORTS}'


def to_shodan(bare):
    """body="x" / title="x" 直接可用；|| 改 OR 语义"""
    s = bare.replace('||', ' ')
    s = re.sub(r'body="([^"]+)"', r'"\1"', s)
    return s + ' -country:CN'


def build_pairs(a):
    """把各模式的条件拼成 bare 查询串"""
    pairs = []
    if getattr(a, "keyword", False):
        spe = [x.strip() for x in a.species.split(",") if x.strip() in SPECIES]
        for s in spe:
            pairs += SPECIES[s]
    if getattr(a, "icon_hash", None):
        pairs.append(f'icon_hash="{a.icon_hash}"')
    if getattr(a, "body", None):
        pairs.append(f'body="{a.body}"')
    return pairs


def build_jobs(engines, pairs, geo):
    jobs = []
    for e in engines:
        for bare in pairs:
            jobs.append((e, to_fofa(bare, geo) if e.name == "fofa" else to_shodan(bare)))
    return jobs


def run_jobs(jobs, pages, seen=None):
    """跑查询并去重，返回行列表"""
    seen = seen if seen is not None else set()
    rows = []
    for i, (eng, q) in enumerate(jobs, 1):
        for pg in range(1, pages + 1):
            res = eng.retry(eng.search, q, pg)
            new = 0
            for r in res:
                key = (r.get("ip") or "") + ":" + str(r.get("port") or "")
                if key in seen and key != ":":
                    continue
                seen.add(key)
                rows.append(r); new += 1
            print(f"  [{i}/{len(jobs)}] p{pg} {eng.name} +{new}  {q[:70]}", file=sys.stderr)
    return rows


COLS = ["host", "ip", "port", "protocol", "title", "domain", "country",
        "city", "server", "icon_hash", "source", "query"]


def write_csv(rows, path):
    with open(path, "w", newline="", encoding="utf-8-sig") as f:
        w = csv.DictWriter(f, fieldnames=COLS, extrasaction="ignore")
        w.writeheader(); w.writerows(rows)


def pick_engines(spec):
    """缺 key 的引擎整体跳过，只提示一次"""
    need = {"fofa": ["FOFA_KEY"], "shodan": ["SHODAN_KEY"],
            "censys": ["CENSYS_ID", "CENSYS_SECRET"]}
    usable = []
    for n in [x.strip() for x in spec.split(",") if x.strip()]:
        if n not in ENGINES:
            print(f"[!] 未知引擎 {n}", file=sys.stderr); continue
        miss = [k for k in need[n] if not os.environ.get(k)]
        if miss:
            print(f"[!] 跳过 {n}（缺 {','.join(miss)}）", file=sys.stderr)
        else:
            usable.append(n)
    return usable


def batch_mode(a):
    """一条命令：取样 → 扩线 → 落单 → 可选判定"""
    import subprocess
    seeds = [s for s in re.split(r"[,\s]+", a.batch) if s.strip()]
    if not seeds:
        print("[=] 没给种子", file=sys.stderr); return
    out = a.output.rstrip("/\\")
    os.makedirs(out, exist_ok=True)

    usable = pick_engines(a.engine)
    if not usable:
        print("[=] 没有可用引擎，先 export 凭据", file=sys.stderr); return
    engines = [ENGINES[n](a.qps) for n in usable]

    fps, pairs = [], []
    for s in seeds:
        fp = fingerprint(s)
        fps.append(fp)
        print(f"[fp] {s}  hash={fp['favicon_hash']} static={len(fp['static_paths'])} "
              f"api={len(fp['api_prefixes'])} {fp['note']}", file=sys.stderr)
        if fp["favicon_hash"] is not None:
            pairs.append(('icon_hash', f'icon_hash="{fp["favicon_hash"]}"'))
        for p in fp["static_paths"][:6]:
            pairs.append(('body', f'body="{p}"'))
        for p in fp["api_prefixes"][:4]:
            pairs.append(('body', f'body="{p}"'))

    json.dump(fps, open(f"{out}/fp.json", "w", encoding="utf-8"),
              ensure_ascii=False, indent=2)

    # 同组指纹去重
    seen_p, clean = set(), []
    for k, p in pairs:
        if p in seen_p:
            continue
        seen_p.add(p); clean.append(p)
    if not clean:
        print("[=] 种子里没有可用指纹（可能全被盾挡），换内页或手搓", file=sys.stderr)
        return

    print(f"[*] {len(seeds)} 个种子 → {len(clean)} 条指纹查询", file=sys.stderr)
    rows = run_jobs(build_jobs(engines, clean, a.geo), 1)

    # 已有结果先并入，同目录分批累积
    squad = f"{out}/squad.csv"
    if os.path.exists(squad):
        with open(squad, encoding="utf-8-sig") as f:
            rows = list(csv.DictReader(f)) + rows
        print(f"[*] 并入已有 squad.csv", file=sys.stderr)
    if not rows:
        print("[=] 扩线无结果", file=sys.stderr); return

    uniq, seen_h = [], set()
    for r in rows:
        h = (r.get("host") or "").strip()
        if not h or h in seen_h:
            continue
        seen_h.add(h); uniq.append(r)
    write_csv(uniq, squad)

    cand = f"{out}/candidates.txt"
    got = sorted({f"https://{r['host']}" for r in uniq if r.get("host") and "." in r["host"]})
    open(cand, "w", encoding="utf-8").write("\n".join(got))
    print(f"[+] {len(got)} 个候选 -> {cand}", file=sys.stderr)

    if a.verify:
        if not got:
            return
        v = os.path.join(os.path.dirname(os.path.abspath(__file__)), "bot_verify.py")
        print("[*] 进入判定...", file=sys.stderr)
        subprocess.run([sys.executable, v, "-i", cand, "-o", f"{out}/probe.jsonl",
                        "-c", str(a.conc), "--timeout", "8"])
        hist = {}
        for line in open(f"{out}/probe.jsonl", encoding="utf-8"):
            try:
                lv = json.loads(line).get("level", "discard")
            except Exception:
                continue
            hist[lv] = hist.get(lv, 0) + 1
        print(f"[+] 判定分布 {hist}  -> {out}/probe.jsonl", file=sys.stderr)
        print(f"[+] 复核 P0/P1: grep -E '\"P[01]\"' {out}/probe.jsonl", file=sys.stderr)


def main():
    ap = argparse.ArgumentParser()
    mode = ap.add_mutually_exclusive_group(required=True)
    mode.add_argument("--keyword", action="store_true")
    mode.add_argument("--icon-hash")
    mode.add_argument("--body")
    mode.add_argument("--query")
    mode.add_argument("--from-sample")
    mode.add_argument("--batch", help="种子 URL，逗号分隔；一条命令跑完 取样→扩线→落单→判定")
    ap.add_argument("--engine", default="fofa", help="fofa,shodan,censys 逗号分隔")
    ap.add_argument("--species", default="cat,money,usdt", help="关键词簇: " + ",".join(SPECIES))
    ap.add_argument("--geo", default="all", choices=list(GEO))
    ap.add_argument("--pages", type=int, default=1)
    ap.add_argument("--qps", type=float, default=1.0)
    ap.add_argument("--verify", action="store_true", help="batch 模式下接着跑判定")
    ap.add_argument("--conc", type=int, default=30, help="--verify 的并发数")
    ap.add_argument("-o", "--output", required=True, help="文件路径；batch 模式下当目录")
    a = ap.parse_args()

    if a.batch:
        batch_mode(a)
        return

    # 取指纹模式
    if a.from_sample:
        fp = fingerprint(a.from_sample)
        open(a.output, "w", encoding="utf-8").write(json.dumps(fp, ensure_ascii=False, indent=2))
        print(json.dumps(fp, ensure_ascii=False, indent=2))
        return

    names = [x.strip() for x in a.engine.split(",") if x.strip()]
    usable = pick_engines(a.engine)
    if not usable:
        print("[=] 没有可用引擎，先 export 凭据", file=sys.stderr); return
    engines = [ENGINES[n](a.qps) for n in usable]

    if a.query:
        jobs = [(e, a.query) for e in engines]
    else:
        jobs = build_jobs(engines, build_pairs(a), a.geo)

    print(f"[*] {len(jobs)} 条查询 × {a.pages} 页，引擎 {names}", file=sys.stderr)
    seen, rows = set(), []
    for i, (eng, q) in enumerate(jobs, 1):
        for pg in range(1, a.pages + 1):
            res = eng.retry(eng.search, q, pg)
            new = 0
            for r in res:
                key = (r.get("ip") or "") + ":" + str(r.get("port") or "")
                if key in seen and key != ":":
                    continue
                seen.add(key)
                rows.append(r); new += 1
            print(f"  [{i}/{len(jobs)}] p{pg} {eng.name} +{new}  {q[:70]}", file=sys.stderr)

    if not rows:
        print("[=] 无结果", file=sys.stderr); return

    cols = ["host", "ip", "port", "protocol", "title", "domain", "country",
            "city", "server", "icon_hash", "source", "query"]
    with open(a.output, "w", newline="", encoding="utf-8-sig") as f:
        w = csv.DictWriter(f, fieldnames=cols, extrasaction="ignore")
        w.writeheader(); w.writerows(rows)

    by_src = {}
    for r in rows:
        by_src[r["source"]] = by_src.get(r["source"], 0) + 1
    print(f"[+] {len(rows)} 条 -> {a.output}  ({by_src})", file=sys.stderr)
    print("[!] 只来自单引擎的命中，判定时要降一级", file=sys.stderr)


def fingerprint(url):
    """从活样本取指纹：favicon hash、静态资源名、接口前缀、服务器头"""
    u = url if url.startswith("http") else "https://" + url
    base = "{0.scheme}://{0.netloc}".format(urlparse(u))
    out = {"sample": u, "favicon_hash": None, "static_paths": [], "api_prefixes": [],
           "server": "", "title": "", "note": ""}
    r = None
    H = {
        "User-Agent": UA,
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
        "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
    }
    for i in range(3):          # SSL EOF / 瞬断，裸连接重试
        try:
            r = requests.get(u, headers=H, timeout=(5, 12), verify=False, allow_redirects=True)
            break
        except Exception as e:
            if i == 2:
                out["note"] = f"取不到首页: {type(e).__name__} {e}"
                return out
            time.sleep(1.5)
    body = r.text or ""
    out["server"] = r.headers.get("Server", "")
    m = re.search(r"<title[^>]*>([^<]{1,80})", body, re.I)
    out["title"] = m.group(1).strip() if m else ""
    out["static_paths"] = sorted(set(
        re.findall(r'/(?:static|assets|js|css)/[A-Za-z0-9_.-]{8,}', body)))[:8]
    out["api_prefixes"] = sorted(set(
        re.findall(r'"(/[a-z]{2,8}/api/|/api/[a-z]{2,10}/)"', body)))[:8]
    if mmh3_ok():
        import mmh3
        cand = [base + "/favicon.ico"]
        m2 = re.search(r'<link[^>]+rel=["\']?(?:shortcut )?icon["\']?[^>]*href=["\']([^"\']+)', body, re.I)
        if m2:
            h = m2.group(1)
            cand.insert(0, h if h.startswith("http") else base + "/" + h.lstrip("/"))
        for c in cand:
            try:
                rr = requests.get(c, headers=H, timeout=(5, 10), verify=False)
                if rr.status_code == 200 and len(rr.content) > 60:
                    out["favicon_hash"] = mmh3.hash(rr.content)
                    break
            except Exception:
                pass
    if not out["static_paths"] and not out["api_prefixes"] and not out["favicon_hash"]:
        out["note"] = "无可用指纹，换一个页面或人工看源码"
    if re.search(r"Just a moment|cf-challenge|Attention Required", out["title"], re.I):
        out["note"] = (out["note"] + " | CF盾页面，favicon 可能是云厂商的，换内页重取").strip(" |")
    return out


def mmh3_ok():
    try:
        import mmh3  # noqa
        return True
    except ImportError:
        return False


if __name__ == "__main__":
    main()
