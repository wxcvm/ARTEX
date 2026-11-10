#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""仓库全树密钥扫描：给定候选仓库，拉全树 → 逐文本文件 → 正则匹配密钥 + 品牌关键词。
关键点：(1) 用 git/trees?recursive=1 一次拿全树；(2) 对每个 <400KB 的文本文件拉 raw；
(3) 正则必须同时匹配 sk- 令牌与 key/secret/password/token 赋值；(4) 命中后人工区分
    占位符（your_key_here / sk-abc123 / sk-xxxxxxxx）与真实值（格式规整、多文件重复出现）。
"""
import re, time, urllib.parse, urllib.request
UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/124.0 Safari/537.36"
KEYWORDS = ["oneapi", "aigc", "tgt8272", "example", "tgt6403", "tgt1215"]
SECRET_RE = re.compile(r"sk-[A-Za-z0-9_\-]{16,}|"
                       r"(?:api[_-]?key|secret|password|passwd|token|SESSION_SECRET|CRYPTO_SECRET|JWT_SECRET)"
                       r"\s*[:=]\s*[\"']?[^\s\"',;]{6,}", re.I)

def get(url, accept="*/*", maxb=400000):
    try:
        with urllib.request.urlopen(urllib.request.Request(
                url, headers={"User-Agent": UA, "Accept": accept}), timeout=30) as r:
            return r.status, r.read(maxb).decode("utf-8", "replace")
    except Exception as e:
        return None, "ERR %s" % e

def sweep(repo):
    st, body = get("https://api.github.com/repos/%s/git/trees/HEAD?recursive=1" % repo,
                   accept="application/vnd.github+json", maxb=3000000)
    if st != 200:
        print("tree fail", st, body[:150]); return
    tree = __import__("json").loads(body).get("tree", [])
    cands = [x for x in tree if x.get("type") == "blob" and (x.get("size") or 0) < 400000
             and re.search(r"\.(env|ya?ml|json|properties|ini|conf|cfg|txt|md|py|js|ts|sh|xml|toml|example)$",
                           x["path"], re.I)]
    print("%s: %d entries, %d text candidates" % (repo, len(tree), len(cands)))
    for x in cands:
        st2, b2 = get("https://raw.githubusercontent.com/%s/HEAD/%s"
                      % (repo, urllib.parse.quote(x["path"])))
        if st2 != 200:
            continue
        kw = [k for k in KEYWORDS if k in b2.lower()]
        secs = sorted(set(target22.example.com(0)[:160] for m in SECRET_RE.finditer(b2)))
        if kw or secs:
            print("  [HIT] %-52s kw=%s" % (x["path"], ",".join(kw)))
            for s in secs[:10]:
                print("        %s" % s)
        time.sleep(0.25)

if __name__ == "__main__":
    import sys
    for repo in sys.argv[1:]:
        sweep(repo)
