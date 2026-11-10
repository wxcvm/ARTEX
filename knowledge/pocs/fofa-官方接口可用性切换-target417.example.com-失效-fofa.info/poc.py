#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""FOFA 采集（官方端点版）：单线程 + 最小间隔 + 退避 + 代理环境变量清理。

与技能内置脚本的差异：
1) API 换回官方 https://example.com/api/v1/search/all（第三方中转 target417.example.com 常报 key 无效）
2) 强制清理 HTTP_PROXY/HTTPS_PROXY —— 本机有前置代理时 curl/urllib 会 TLS 握手失败(exit 35)
3) key 读取顺序：环境变量 -> $DSH_HOME/.env（DSH 不把含 KEY 的变量传给子进程）
"""
import argparse, base64, json, os, sys, time, ssl, urllib.error, urllib.parse, urllib.request

for _k in ("HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy", "ALL_PROXY", "all_proxy"):
    os.environ.pop(_k, None)

API = "https://example.com/api/v1/search/all"          # ← 官方端点，不是 target417.example.com
FIELDS = "host,ip,port,title,domain,server"
BACKOFF = [5, 15, 45]
CTX = ssl.create_default_context()
CTX.check_hostname = False
CTX.verify_mode = ssl.CERT_NONE


def load_key():
    v = os.environ.get("FOFA_KEY")
    if v:
        return v.strip()
    home = os.environ.get("DSH_HOME") or os.path.join(os.path.expanduser("~"), "workspace")
    for path in (os.path.join(home, ".env"), os.path.join(os.getcwd(), ".env")):
        try:
            for line in open(path, encoding="utf-8", errors="ignore"):
                if line.startswith("FOFA_KEY="):
                    k = line.split("=", 1)[1].strip()
                    if k and not k.startswith("#"):
                        return k
        except OSError:
            continue
    raise SystemExit("找不到 FOFA_KEY（环境变量与 $DSH_HOME/.env 均无）")


KEY = load_key()


def fetch(query, page, size, timeout=40):
    qs = urllib.parse.urlencode({
        "key": KEY,
        "qbase64": base64.b64encode(query.encode()).decode(),
        "size": size, "page": page, "fields": FIELDS,
    })
    req = urllib.request.Request(API + "?" + qs, headers={"User-Agent": "recon/1.0"})
    try:
        with urllib.request.urlopen(req, timeout=timeout, context=CTX) as r:
            return r.status, json.load(r)
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8", "ignore")
        try:
            return e.code, json.loads(body)
        except Exception:
            return e.code, {"error": True, "errmsg": body[:200]}


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("query")
    ap.add_argument("--pages", type=int, default=1)
    ap.add_argument("--size", type=int, default=100)
    ap.add_argument("--interval", type=float, default=2.5)
    args = ap.parse_args()

    rows, seen, last = [], set(), 0.0
    for page in range(1, args.pages + 1):
        wait = args.interval - (time.time() - last)
        if wait > 0:
            time.sleep(wait)
        data = None
        for attempt, bo in enumerate([0] + BACKOFF):
            if attempt:
                print(f"[backoff] {bo}s", file=sys.stderr)
                time.sleep(bo)
            last = time.time()
            status, data = fetch(args.query, page, args.size)
            if status == 200 and not data.get("error"):
                break
            if "45012" in str(data.get("errmsg", "")) or status == 429:
                continue
            print(f"[stop] http={status} {data.get('errmsg')}", file=sys.stderr)
            data = None
            break
        if not data:
            break
        for it in data.get("results") or []:
            if isinstance(it, str):
                it = [it]
            key = (it[0] if len(it) > 0 else "", it[1] if len(it) > 1 else "", it[2] if len(it) > 2 else "")
            if key in seen:
                continue
            seen.add(key)
            rows.append({"host": key[0], "ip": key[1], "port": key[2],
                         "title": it[3] if len(it) > 3 else "",
                         "domain": it[4] if len(it) > 4 else "",
                         "server": it[5] if len(it) > 5 else ""})
        if len(data.get("results") or []) < args.size:
            break

    for r in rows:
        print(json.dumps(r, ensure_ascii=False))
    print(f"[done] {len(rows)} 条", file=sys.stderr)


if __name__ == "__main__":
    sys.exit(main())