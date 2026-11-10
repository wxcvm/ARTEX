#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# Webpack 异步 chunk 收割：从主包提取 chunkId->hash 映射，批量下载 /static/js/async/<id>.<hash>.js，
# 正则提取全部 /api/ 端点。比目录字典爆破准确率高一个量级（覆盖懒加载页面才调用的定制 API）。
import requests, re, time, sys
BASE = sys.argv[1].rstrip("/")          # e.g. https://example.com
MAIN = sys.argv[2]                      # 主包路径，如 /static/js/index.3137992311.js
S = requests.Session(); S.trust_env = True
S.headers.update({"User-Agent": "Mozilla/5.0", "Referer": BASE + "/"})
body = S.get(BASE + MAIN, timeout=60).text
chunks = {}
for k, v in re.findall(r'(\d{1,4}):"([a-f0-9]{8,20})"', body):
    chunks.setdefault(k, v)
print("chunk ids:", len(chunks))
# 主包本身也提取
urls = set(re.findall(r'["`/](/api/[A-Za-z0-9._/${}?:=&-]{2,70})', body))
for cid, h in sorted(chunks.items(), key=lambda x: int(x[0])):
    url = f"{BASE}/static/js/async/{cid}.{h}.js"
    for t in range(3):
        try:
            r = S.get(url, timeout=40)
            if r.status_code == 200:
                urls |= set(re.findall(r'["`/](/api/[A-Za-z0-9._/${}?:=&-]{2,70})', r.text))
            break
        except Exception:
            time.sleep(3 + t * 3)
    time.sleep(0.4)
for u in sorted(urls): print(u)
# 用法: python webpack_chunk_api_enum.py https://example.com /static/js/index.3137992311.js
# 主包路径从 GET / 首页 HTML 的 <script src> 取；模板串 "static/js/async/"+e+"."+{...}[e]+".js" 确认目录。
