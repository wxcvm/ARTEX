# -*- coding: utf-8 -*-
"""非标准服务识别：空读 → HTTP → TLS → 乱包，四态对照（原生 socket，不走代理）。"""
import socket, ssl, os, sys, re
for k in ("http_proxy","https_proxy","all_proxy","HTTP_PROXY","HTTPS_PROXY","ALL_PROXY"):
    os.environ.pop(k, None)
try:
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
except Exception:
    pass

def rd(s, n=400, t=4.0):
    s.settimeout(t)
    try:
        return s.recv(n)
    except Exception:
        return b""

def show(ip, port):
    print("\n########## %s:%d ##########" % (ip, port))
    # 1) 空读（服务端先发协议：SSH/MySQL/SMTP/Redis...）
    try:
        s = socket.create_connection((ip, port), timeout=6)
        b = rd(s, 400, 4.0)
        print("  [空读] %r" % (b[:200],))
        s.close()
    except Exception as e:
        print("  [空读] 连接失败 %s: %s" % (type(e).__name__, e))
    # 2) HTTP
    try:
        s = socket.create_connection((ip, port), timeout=6)
        s.sendall(("GET / HTTP/1.1\r\nHost: %s\r\nUser-Agent: Mozilla/5.0\r\nAccept: */*\r\nConnection: close\r\n\r\n" % ip).encode())
        b = rd(s, 1200, 6.0)
        print("  [HTTP] %r" % (b[:400],))
        s.close()
    except Exception as e:
        print("  [HTTP] %s: %s" % (type(e).__name__, e))
    # 3) TLS
    try:
        s = socket.create_connection((ip, port), timeout=6)
        ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
        ctx.check_hostname = False; ctx.verify_mode = ssl.CERT_NONE
        try: ctx.set_ciphers("DEFAULT@SECLEVEL=1")
        except Exception: pass
        ss = ctx.wrap_socket(s, server_hostname=ip)
        der = ss.getpeercert(binary_form=True)
        names = re.findall(rb"[!-~]{5,80}", der or b"")
        cand = [n.decode("latin-1") for n in names if b"." in n and b" " not in n and not n.startswith(b"http")]
        print("  [TLS ] OK version=%s der=%dB names=%s" % (ss.version(), len(der or b""), cand[:12]))
        ss.close()
    except Exception as e:
        print("  [TLS ] %s: %s" % (type(e).__name__, e))
    # 4) 乱包（看是否回显错误）
    try:
        s = socket.create_connection((ip, port), timeout=6)
        s.sendall(b"\x00\x00\x00\x00\xff\xff\xff\xffGET / HTTP/1.0\r\n\r\n")
        b = rd(s, 500, 5.0)
        print("  [乱包] %r" % (b[:300],))
        s.close()
    except Exception as e:
        print("  [乱包] %s: %s" % (type(e).__name__, e))

TARGETS = [(t.split(":")[0], int(t.split(":")[1])) for t in sys.argv[1].split(",") if t.strip()]
for ip, port in TARGETS:
    show(ip, port)
