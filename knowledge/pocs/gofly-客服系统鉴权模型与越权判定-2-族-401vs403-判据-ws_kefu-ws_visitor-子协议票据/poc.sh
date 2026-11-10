#!/usr/bin/env python3
# GoFly 客服系统鉴权判定三件套（HTTP + WS），零副作用、只读
# 1) 判定 /2/ 族认哪种票据（401=无票/票错，403=票有效但与 visitor_id 不匹配 ⇒ 提权判据）
for H in token authorization; do
  curl -sS -k --resolve example.com:443:203.0.113.88 \
    -H "$H: Bearer <VISITOR_TICKET>" \
    "https://example.com/2/messagesPages?page=1&pagesize=3&visitor_id=<OWN_VISITOR_ID>"
done
curl -sS -k --resolve example.com:443:203.0.113.88 \
  "https://example.com/2/messagesPages?page=1&pagesize=3&visitor_id=<OTHER_VISITOR_ID>"   # 403 => 无 IDOR

# 2) WS 子协议鉴权（票在 Sec-WebSocket-Protocol 里，浏览器没法自定义头）
#    客服侧:  ["hch-admin", "auth."+token]       访客侧: ?visitor_id=<id>  + ["hch-support", ticket]
python3 - <<'PY'
import socket, ssl, base64, os
def hs(path, protos, host="example.com", ip="203.0.113.88"):
    key = base64.b64encode(os.urandom(16)).decode()
    ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT); ctx.check_hostname=False; ctx.verify_mode=ssl.CERT_NONE
    s = ctx.wrap_socket(socket.create_connection((ip,443),timeout=10), server_hostname=host)
    req = [f"GET {path} HTTP/1.1", f"Host: {host}", "Upgrade: websocket", "Connection: Upgrade",
           f"Sec-WebSocket-Key: {key}", "Sec-WebSocket-Version: 13",
           "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/124"]
    if protos: req.append("Sec-WebSocket-Protocol: " + ", ".join(protos))
    s.sendall(("\r\n".join(req)+"\r\n\r\n").encode())
    buf=b""
    while b"\r\n\r\n" not in buf: buf += s.recv(4096)
    head, _, _ = buf.partition(b"\r\n\r\n")
    return head.decode("utf-8","replace").split("\r\n")[0]
print(hs("/ws_kefu", ["hch-admin","auth.<KEFU_JWT>"]))                       # 101 才是真客服票；其它=200+token失效
print(hs("/ws_visitor?visitor_id=<OWN>", ["hch-support","<VISITOR_TICKET>"])) # 101
print(hs("/ws_visitor?visitor_id=<OTHER>", ["hch-support","<VISITOR_TICKET>"]))# 401 => 票据与 id 精确绑定
PY

# 3) 未鉴权面清点（上游 GoFly 默认全是裸的，二开版可能已加门）
for p in /notice?kefu_id=agent /visitors_online /message_notice /visitor_login /uploadfile /uploadimg /about /config; do
  curl -sS -k -o /dev/null -w "$p %{http_code}\n" --resolve example.com:443:203.0.113.88 "https://example.com$p"
done