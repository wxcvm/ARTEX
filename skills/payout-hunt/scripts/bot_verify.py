#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
bot_verify.py —— 六道 Gate 粗判器（payout-hunt 技能 · 判定阶段）

输入 bot_hunt.py 的候选清单，逐条跑 Gate 1-5 + favicon，输出等级与证据字段。
**只做判定，不做利用**：后台面探测只取状态码与标题特征，不发任何验证性 payload。

用法（在技能目录下执行）:
    python scripts/bot_verify.py -i out/candidates.txt -o out/probe.jsonl -c 40
    python scripts/bot_verify.py -i out/candidates.txt -o out/backdoor.jsonl -c 20 --stage backdoor

输入每行:  URL   或   URL | 标题
输出:      JSONL，字段见 SKILL.md「产出契约」

依赖: requests（必需）；mmh3（favicon 哈希，缺失时跳过该 Gate 而不是整体失败）
"""
import argparse, json, re, sys, threading, time
from concurrent.futures import ThreadPoolExecutor, as_completed
from urllib.parse import urlparse

import requests
import urllib3
urllib3.disable_warnings()

try:
    import mmh3
except ImportError:
    mmh3 = None

UA = ("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
      "(KHTML, like Gecko) Chrome/124.0 Safari/537.36")
HEAD = {"User-Agent": UA, "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8"}
TIMEOUT = (5, 8)

# ---------- 判据 ----------
BLACKLIST = {
    "parking":   re.compile(r"sedoparking|hugedomains|dan\.com|GoDaddy|this domain is for sale", re.I),
    "dirindex":  re.compile(r"<title>\s*Index of /", re.I),
    "cfshield":  re.compile(r"cf-challenge|Just a moment|cf_chl_opt|__cf_chl", re.I),
    "hn_aiohttp":re.compile(r"server:\s*aiohttp", re.I),
}
NAV = re.compile(r"十大平台|排行榜|导航|测评|推荐平台|网址大全")

VENDOR = ["/api/game/", "/api/ag/", "/api/mg/", "/api/pt/", "/api/pg/",
          "/api/bbin/", "/api/ky/", "/api/db/", "/api/og/", "/api/cq9/"]

CLUSTERS = {
    "brand": re.compile(r"葡京|威尼斯|太阳城|金沙|银河|开元|美高梅|皇冠|新濠|永利"),
    "cat":   re.compile(r"棋牌|真人|电子游艺|老虎机|体育|电竞|彩票|捕鱼|包杀|casino|sportsbook|betting", re.I),
    "money": re.compile(r"静态收益|动态收益|抢单|分红|拆分|互助|量化|日化|返水|佣金|秒到账"),
    "usdt":  re.compile(r"USDT|TRC20|TronLink|钱包地址|归集|热钱包"),
}

TRON_RE = re.compile(r"\bT[1-9A-HJ-NP-Za-km-z]{33}\b")
IFRAME_RE = re.compile(r"<iframe", re.I)

# 后台面字典（前 12 条默认跑，--stage backdoor 全跑）
ADMIN = ["/admin", "/admin/", "/admin/login", "/manage", "/houtai", "/system",
         "/boss", "/kf", "/agent", "/login", "/api-docs", "/v2/api-docs"]
ADMIN_FULL = ADMIN + ["/admin/index.php", "/manage/login", "/ht", "/agency", "/daili",
                      "/swagger-ui.html", "/swagger-ui/index.html", "/.env", "/.git/config",
                      "/backup.sql", "/www.zip", "/server-status", "/actuator/env",
                      "/api/game/agent/register"]
ADMIN_HIT = re.compile(r"password|登录|验证码|username|admin", re.I)

_lock = threading.Lock()
_done = [0]


def norm(u):
    u = u.strip()
    if not u:
        return None
    if not u.startswith(("http://", "https://")):
        u = "https://" + u
    return u.rstrip("/")


def get(s, url, **kw):
    try:
        return s.get(url, headers=HEAD, timeout=TIMEOUT, verify=False,
                     allow_redirects=True, **kw)
    except Exception:
        return None


def gate1(s, u):
    r = get(s, u)
    if r is None:
        return "fail", [], "timeout/conn"
    if BLACKLIST["cfshield"].search(r.text[:4000]):
        return "shielded", [], f"http{r.status_code}"
    if r.status_code in (502, 503, 504, 522, 520):
        return "fail", [], f"http{r.status_code}"
    body = r.text or ""
    # SPA 首页可能只有骨架，body 短不等于死：有 title 或挂载点即放行
    if len(body) < 800 and not re.search(r"<title|<div id=|<script src=", body, re.I):
        return "fail", [], f"short body {len(body)}"
    for k in ("parking", "dirindex"):
        if BLACKLIST[k].search(body[:5000]):
            return "fail", [], k
    return "pass", [], f"http{r.status_code} len{len(body)}"


def gate2(u, body):
    ev = [v for v in VENDOR if v in body]
    if len(ev) >= 2:
        return "pass", ev
    hits = [k for k, rx in CLUSTERS.items() if rx.search(body)]
    if len(hits) >= 2 and sum(len(CLUSTERS[k].findall(body)) for k in hits) >= 3:
        return "pass", hits
    if NAV.search(body) and len(hits) < 2:
        return "fail", hits + ["导航"]
    return "fail", hits


def gate3(s, base, full=False):
    ev, best = [], "fail"
    for p in (ADMIN_FULL if full else ADMIN):
        r = get(s, base + p)
        if r is None:
            continue
        if r.status_code == 200 and ADMIN_HIT.search(r.text[:3000]):
            return "pass", [f"{p} 200"]
        if r.status_code == 200 and p in ("/.env", "/.git/config", "/api-docs", "/v2/api-docs"):
            return "pass", [f"{p} 200"]
        if r.status_code == 200 and "json" in r.headers.get("Content-Type", ""):
            ev.append(f"{p} 200 json")
    return ("pass" if ev else best), ev


def gate4(s, base, body):
    addrs = list(set(TRON_RE.findall(body)))
    ev = list(addrs[:3])
    for p in ("/pay/", "/cashier/", "/recharge/", "/withdraw/"):
        r = get(s, base + p)
        if r and r.status_code == 200:
            ev.append(f"{p} 200")
            addrs += TRON_RE.findall(r.text or "")
    if addrs:
        return "pass", sorted(set(addrs))[:5]
    if any(e.endswith("200") for e in ev):
        return "pass", ev
    return "fail", ev


def gate5(base, body):
    if len(IFRAME_RE.findall(body)) >= 1 and len(body) < 6000:
        return "fail", ["iframe壳"]
    hosts = set(re.findall(r"https?://([a-z0-9.-]+\.[a-z]{2,})", body))
    same = urlparse(base).netloc
    ext = [h for h in hosts if h != same and "cdn" not in h][:3]
    if len(ext) >= 4:
        return "fail", ["全外链"]
    return "pass", []


def fav(s, base):
    if not mmh3:
        return None
    cand = [base + "/favicon.ico"]
    r0 = get(s, base)
    if r0:
        m = re.search(r'<link[^>]+rel=["\']?(?:shortcut )?icon["\']?[^>]*href=["\']([^"\']+)',
                      r0.text or "", re.I)
        if m:
            h = m.group(1)
            cand.insert(0, h if h.startswith("http") else base + "/" + h.lstrip("/"))
    for u2 in cand:
        r = get(s, u2)
        if r and r.status_code == 200 and len(r.content) > 60:
            try:
                return mmh3.hash(r.content)
            except Exception:
                pass
    return None


def one(u, stage):
    full = stage == "backdoor"
    out = {"url": u, "gate1": "fail", "gate2": "fail", "gate3": "fail",
           "gate4": "fail", "gate5": "fail", "type": None, "level": "discard",
           "evidence": [], "usdt": [], "next": "none", "note": ""}
    s = requests.Session()
    g1, _, note = gate1(s, u)
    out["gate1"] = g1
    out["note"] = note
    if g1 == "shielded":
        out["level"] = "P3"; return out
    if g1 != "pass":
        return out

    r = get(s, u)
    body = (r.text or "") if r else ""
    out["favicon"] = fav(s, u)

    g2, ev2 = gate2(u, body)
    out["gate2"] = g2
    out["evidence"] += ev2
    if g2 != "pass":
        # SPA 首页无文案，不能直接判死：挂 P3 等渲染
        if len(body) < 8000 and re.search(r"<div id=|<script src=|__NUXT__|__NEXT_DATA__", body, re.I):
            out["level"] = "P3"
            out["note"] += " | spa需渲染"
        return out

    g5, ev5 = gate5(u, body)
    out["gate5"] = g5
    out["evidence"] += ev5

    g4, ev4 = gate4(s, u, body)
    out["gate4"] = g4
    out["usdt"] = [e for e in ev4 if TRON_RE.fullmatch(e)]
    out["evidence"] += [e for e in ev4 if e not in out["usdt"]]

    g3, ev3 = gate3(s, u, full)
    out["gate3"] = g3
    out["evidence"] += ev3

    kind = CLUSTERS["money"].search(body) or CLUSTERS["brand"].search(body)
    out["type"] = "T2" if (CLUSTERS["money"].search(body) and not CLUSTERS["brand"].search(body)) else "T1"

    if g3 == "pass" and g4 == "pass" and g5 == "pass":
        out["level"] = "P0" if out["usdt"] else "P1"
        out["next"] = "path_A/B" if out["level"] == "P0" else "path_C"
    elif g4 == "pass":
        out["level"] = "P2"
    else:
        out["level"] = "P3"
    return out


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("-i", "--input", required=True)
    ap.add_argument("-o", "--output", required=True)
    ap.add_argument("-c", "--concurrency", type=int, default=40)
    ap.add_argument("--timeout", type=int, default=8, help="读取超时秒数，连接固定 5s")
    ap.add_argument("--stage", default="probe", choices=["probe", "backdoor"])
    a = ap.parse_args()

    global TIMEOUT
    TIMEOUT = (5, a.timeout)

    urls = []
    with open(a.input, encoding="utf-8") as f:
        for line in f:
            u = norm(line.split("|")[0])
            if u and u not in urls:
                urls.append(u)
    print(f"[*] {len(urls)} targets, conc={a.concurrency}, stage={a.stage}", file=sys.stderr)

    with open(a.output, "w", encoding="utf-8") as fo, \
         ThreadPoolExecutor(max_workers=a.concurrency) as ex:
        futs = {ex.submit(one, u, a.stage): u for u in urls}
        for fu in as_completed(futs):
            try:
                o = fu.result()
            except Exception as e:
                o = {"url": futs[fu], "level": "discard", "note": f"err {e}"}
            fo.write(json.dumps(o, ensure_ascii=False) + "\n")
            fo.flush()
            with _lock:
                _done[0] += 1
                if _done[0] % 10 == 0:
                    print(f"  {_done[0]}/{len(urls)}", file=sys.stderr)
    print(f"[+] done -> {a.output}", file=sys.stderr)

    # 抽出全部链上地址，供后续查余额
    addrs = set()
    with open(a.output, encoding="utf-8") as f:
        for line in f:
            try:
                addrs.update(json.loads(line).get("usdt") or [])
            except Exception:
                pass
    if addrs:
        p = a.output.rsplit(".", 1)[0] + ".addresses.txt"
        open(p, "w", encoding="utf-8").write("\n".join(sorted(addrs)))
        print(f"[+] {len(addrs)} addresses -> {p}", file=sys.stderr)


if __name__ == "__main__":
    main()
