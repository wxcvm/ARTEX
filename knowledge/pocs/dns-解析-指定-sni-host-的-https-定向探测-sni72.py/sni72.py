# -*- coding: utf-8 -*-
"""DNS 解析 + 指定 SNI/Host 的 HTTPS 探测（原生 socket，不走代理）。"""
import socket, ssl, os, sys, re, json

for k in ("http_proxy","https_proxy","all_proxy","HTTP_PROXY","HTTPS_PROXY","ALL_PROXY"):
    os.environ.pop(k, None)
try:
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
except Exception:
    pass

DOMAINS = ["example.com","example.com","example.com","api3.example.com",
           "example.com","www.example.com","example.com","www.example.com","target332.example.com",
           "target333.example.com","target334.example.com","target322.example.com","target323.example.com"]

print("===== DNS 解析 =====")
res = {}
for d in DOMAINS:
    try:
        infos = socket.getaddrinfo(d, None, socket.AF_INET)
        ips = sorted({i[4][0] for i in infos})
        res[d] = ips
        print("%-22s -> %s" % (d, ", ".join(ips)))
    except Exception as e:
        res[d] = None
        print("%-22s -> DNS FAIL (%s)" % (d, type(e).__name__))

def probe(ip, port, host, sni, path, timeout=8.0):
    req = ("GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: Mozilla/5.0 (Linux; Android 12; lz2025)\r\n"
           "Accept: */*\r\nX-Security-Version: 2\r\nConnection: close\r\n\r\n" % (path, host))
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(timeout)
    try:
        s.connect((ip, port))
        ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
        ctx.check_hostname = False; ctx.verify_mode = ssl.CERT_NONE
        try: ctx.set_ciphers("DEFAULT@SECLEVEL=1")
        except Exception: pass
        s = ctx.wrap_socket(s, server_hostname=sni)
        s.sendall(req.encode())
        buf = b""
        while len(buf) < 5000:
            try: c = s.recv(4096)
            except socket.timeout: break
            if not c: break
            buf += c
        head, _, body = buf.partition(b"\r\n\r\n")
        ht = head.decode("utf-8","replace")
        st = 0
        m = re.match(r"HTTP/\d\.\d (\d+)", ht)
        if m: st = int(target22.example.com(1))
        return st, body.decode("utf-8","replace")[:300].replace("\n"," ")
    except Exception as e:
        return None, "%s: %s" % (type(e).__name__, e)
    finally:
        try: s.close()
        except Exception: pass

print("\n===== 203.0.113.11:443 定向 SNI/Host 探测 =====")
SNIS = ["example.com","example.com","example.com","example.com","www.example.com","203.0.113.11"]
PATHS = ["/", "/LHService/", "/PayBmob/", "/IMService/admin/", "/LHService/v27/routing/health"]
for sni in SNIS:
    for p in PATHS:
        st, body = probe("203.0.113.11", 443, sni, sni, p)
        print("  SNI/Host=%-20s %-34s st=%-5s %s" % (sni, p, st, body), flush=True)
