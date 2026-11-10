#!/usr/bin/env python3
# example.com（203.0.113.84）_glaxy_tus5r4_ 命名空间 API 探测器
# 用法: py -3 glaxy_probe.py 203.0.113.84 [domainName]
# 说明: 所有业务接口必须带 apiNameSpace 前缀 _glaxy_tus5r4_/；
#       去掉前缀会被 OpenResty 的 SPA try_files 回退页接住（3865B），看起来像"接口不存在/后端下线"。
import socket, ssl, json, sys, hashlib, random, time, re, base64

HOST = sys.argv[1] if len(sys.argv) > 1 else "203.0.113.84"
DOMAIN = sys.argv[2] if len(sys.argv) > 2 else "example.com"
NS = "_glaxy_tus5r4_"
XNS = "_glaxy_tus5r4_/_extra_"
APPID = "dd3b6a3102f54f588e0c8a37e1d96632"
PID = "a80db5262e14467b82d6dec666073364"
V = "1.0.0"
UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/122.0 Safari/537.36"

CTX = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
CTX.check_hostname = False
CTX.verify_mode = ssl.CERT_NONE


def md5(s):
    return hashlib.md5(s.encode()).hexdigest()


def http(path, body=None, headers=None, method="POST", port=443, timeout=15):
    hh = {"Host": HOST, "User-Agent": UA,
          "Accept": "application/json, text/plain, */*",
          "Content-Type": "application/json;charset=UTF-8",
          "Origin": "https://" + DOMAIN, "Referer": "https://" + DOMAIN + "/",
          "Connection": "close"}
    hh.update(headers or {})
    data = json.dumps(body or {}, separators=(",", ":")).encode()
    head = f"{method} {path} HTTP/1.1\r\n" + "".join(f"{k}: {v}\r\n" for k, v in hh.items())
    head += f"Content-Length: {len(data)}\r\n\r\n"
    with socket.create_connection((HOST, port), timeout=timeout) as s:
        with CTX.wrap_socket(s, server_hostname=HOST) as ss:
            ss.sendall(head.encode() + data)
            buf = b""
            while True:
                try:
                    d = ss.recv(65535)
                except socket.timeout:
                    break
                if not d:
                    break
                buf += d
    return buf


def dechunk(b):
    """解析 Transfer-Encoding: chunked（响应体开头形如 '54\\r\\n{...}' 不是乱码）"""
    out = b""
    while True:
        i = b.find(b"\r\n")
        if i < 0:
            break
        try:
            n = int(b[:i].split(b";")[0], 16)
        except ValueError:
            return b
        if n == 0:
            break
        out += b[i + 2:i + 2 + n]
        b = b[i + 2 + n + 2:]
    return out


def show(tag, buf):
    if buf.startswith(b"__ERR__"):
        print("  %-52s ERR %s" % (tag, buf.decode()[:90])); return
    body = buf[buf.find(b"\r\n\r\n") + 4:]
    body = dechunk(body) if buf[:200].find(b"chunked") >= 0 else body
    st = buf.split(b"\r\n", 1)[0].decode("latin1")
    print("  %-52s %-24s %s" % (tag, st, body[:160].decode("utf-8", "replace").strip()))


if __name__ == "__main__":
    print("=" * 100)
    print("A) 无签名 —— 期望 GW_899998（证明后端在线，非 404/502）")
    print("=" * 100)
    for p in [f"/{NS}/init", f"/{NS}/webToken", f"/{NS}/customer/login", f"/{NS}/game/queryGames"]:
        show("POST " + p, http(p, {}))

    print("\n" + "=" * 100)
    print("B) 签名算法复现（ca00.E=字符倒序，u()=6位随机数）")
    print("=" * 100)
    ts = str(int(time.time() * 1000))
    rand = "".join(str(int(9 * random.random())) for _ in range(6))
    qid = md5(ts + rand)
    js = json.dumps({"productId": PID}, separators=(",", ":"))
    sign = md5(js[::-1] + qid + APPID + V)
    hdr = {"Sign": sign, "AppId": APPID, "Qid": qid, "v": V,
           "domainName": DOMAIN, "needEncrypt": "1"}
    print("  Qid=%s\n  Sign=%s" % (qid, sign))
    show("POST /%s/webToken" % NS, http(f"/{NS}/webToken", {"productId": PID}, hdr))
    print("  (errCode=GW_880102 签名验证失败 = 签名头已被识别但值不符；GW_899998 = 未识别)")

    print("\n" + "=" * 100)
    print("C) 对照：同样路径去掉命名空间前缀 —— 这就是'后端下线'误判的来源")
    print("=" * 100)
    for p in ["/init", "/webToken", "/customer/login"]:
        show("POST " + p, http(p, {}))