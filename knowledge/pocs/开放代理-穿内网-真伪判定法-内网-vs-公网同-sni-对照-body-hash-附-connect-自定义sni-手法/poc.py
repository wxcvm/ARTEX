# 可复用：经 HTTP 代理建立 CONNECT 隧道 + 自定义 SNI 取目标响应
import socket, ssl, hashlib

PROXY = ("PROXY_HOST", 5566)

def conn_sni(target_ip, port, sni, path="/", timeout=15):
    """经 Squid CONNECT 到 target_ip:port，用指定 SNI 完成 TLS 握手并发 GET"""
    s = socket.create_connection(PROXY, timeout=timeout)
    s.sendall(("CONNECT %s:%d HTTP/1.1\r\nHost: %s:%d\r\n\r\n"
               % (target_ip, port, target_ip, port)).encode())
    resp = b""; s.settimeout(timeout)
    while b"\r\n\r\n" not in resp:
        c = s.recv(4096)
        if not c: break
        resp += c
    if b" 200 " not in resp.split(b"\r\n")[0]:
        raise RuntimeError("CONNECT rejected: " + resp.split(b"\r\n")[0].decode())
    ctx = ssl.create_default_context(); ctx.check_hostname = False; ctx.verify_mode = ssl.CERT_NONE
    ss = ctx.wrap_socket(s, server_hostname=sni)          # ← 关键：SNI 决定 nginx 选哪个 vhost
    ss.sendall(("GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: Mozilla/5.0\r\n"
                "Connection: close\r\n\r\n" % (path, sni)).encode())
    data = b""
    while True:
        c = ss.recv(16384)
        if not c: break
        data += c
    ss.close()
    head, _, body = data.partition(b"\r\n\r\n")
    return head.decode("utf-8", "replace").split("\r\n")[0], body

# 判定内网是否真有独占面（决定性）
def is_real_intranet_pivot(sni, internal_ip="172.17.0.1", public_ip="TARGET_PUBLIC_IP"):
    st1, b1 = conn_sni(internal_ip, 443, sni)
    st2, b2 = conn_sni(public_ip, 443, sni)       # 注意：public 侧也走同样的 SNI
    same = (st1 == st2) and (hashlib.sha256(b1).hexdigest() == hashlib.sha256(b2).hexdigest())
    return same, st1, len(b1), st2, len(b2)
# same=True  -> 内网入口只是宿主机自身 nginx，无内网独占服务，不算边界突破
# same=False -> 存在差异，继续挖差异点

# 本例实测（203.0.113.53）：
#   target82.example.com        : 内网 200/1166 sha 65d8428308fd41cf == 公网 200/1166 sha 65d8428308fd41cf
#   target38.example.com    : 同上（同一 New API 实例）
#   target83.example.com         : 内网 200/21114 sha 56c5e28a85ae88fa == 公网 200/21114 sha 56c5e28a85ae88fa
#   ⇒ 三例全等价 ⇒ 判定“无内网独占面”。
