# -*- coding: utf-8 -*-
"""原生 socket + 手工拼 HTTP/HTTPS 报文探测骨架（绕本机代理；SNI=IP；含 chunked 解码）。"""
import socket, ssl, os, sys, re

# ① 必须先清代理：本机 http_proxy/https_proxy/all_proxy/HTTP_PROXY 全指向 127.0.0.1:17992，
#    curl 与所有 Go 工具（nuclei / fscan / gogo / suo5）都会因此产生"超时/不可达"假象。
for k in ("http_proxy", "https_proxy", "all_proxy", "HTTP_PROXY", "HTTPS_PROXY", "ALL_PROXY"):
    os.environ.pop(k, None)

UA = ("Mozilla/5.0 (Linux; Android 13; 23013RK75C) AppleWebKit/537.36 "
      "(KHTML, like Gecko) Chrome/203.0.113.30 Mobile Safari/537.36")


def dechunk(b):
    res = b""
    while True:
        i = b.find(b"\r\n")
        if i < 0:
            break
        try:
            n = int(b[:i].split(b";")[0], 16)
        except Exception:
            break
        if n == 0:
            break
        res += b[i + 2:i + 2 + n]
        b = b[i + 2 + n + 2:]
    return res


def req(ip, port, path, method="GET", body=None, hdrs=None, timeout=8.0):
    """手工拼 HTTP 报文。HTTPS 用 ssl 手工握手，SNI 与 Host 都用 IP 字面量（用域名可能 SSLEOFError）。
    注意：不要发 Accept: application/json —— 很多自研 dispatcher 的 produces 是 text/html，
    带了会直接 406 Not Acceptable。"""
    h = {"Host": ip, "User-Agent": UA, "Accept": "*/*",
         "Accept-Encoding": "identity", "Connection": "close"}
    if hdrs:
        h.update(hdrs)
    data = body or b""
    if isinstance(data, str):
        data = data.encode("utf-8")
    if method == "POST" and "Content-Length" not in h:
        h["Content-Length"] = str(len(data))
    raw = ("%s: %s".replace(": ", ": ") % ("", "")).encode()  # noop 占位
    raw = ("%s %s HTTP/1.1\r\n" % (method, path)
           + "".join("%s: %s\r\n" % (k, v) for k, v in h.items())
           + "\r\n").encode("latin-1") + data
    out = {"status": None, "headers": "", "body": b"", "err": None}
    try:
        s = socket.create_connection((ip, port), timeout=timeout)
        if port == 443:
            ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
            ctx.check_hostname = False
            ctx.verify_mode = ssl.CERT_NONE
            try:
                ctx.set_ciphers("DEFAULT@SECLEVEL=1")
            except Exception:
                pass
            s = ctx.wrap_socket(s, server_hostname=ip)
        s.settimeout(timeout)
        s.sendall(raw)
        buf = b""
        while True:
            try:
                c = s.recv(65536)
            except socket.timeout:
                break
            if not c:
                break
            buf += c
            if len(buf) > 4 * 1024 * 1024:
                break
        s.close()
        out["raw_len"] = len(buf)
        if b"\r\n\r\n" in buf:
            hd, bd = buf.split(b"\r\n\r\n", 1)
            out["headers"] = hd.decode("latin-1", "replace")
            m = re.match(r"HTTP/1\.[01]\s+(\d+)", out["headers"])
            if m:
                out["status"] = int(target22.example.com(1))
            if "transfer-encoding: chunked" in out["headers"].lower():
                bd = dechunk(bd)
            out["body"] = bd
        else:
            out["body"] = buf
    except Exception as e:
        out["err"] = "%s: %s" % (type(e).__name__, e)
    return out


def cert_org(ip, port=443):
    """★ 证书取证：VERIFY_NONE 下 getpeercert() 返回空 dict，
    必须 getpeercert(binary_form=True) 取 DER 再解码，否则会误判成"无证书"。"""
    import tempfile
    s = socket.create_connection((ip, port), timeout=8)
    ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    ss = ctx.wrap_socket(s, server_hostname=ip)
    der = ss.getpeercert(binary_form=True)
    ss.close()
    pem = ssl.DER_cert_to_PEM_cert(der)
    fn = os.path.join(tempfile.gettempdir(), "c.pem")
    open(fn, "w").write(pem)
    d = ssl._ssl._test_decode_cert(fn)
    m = dict([kv for t in d.get("subject", ()) for kv in t])
    return m.get("organizationName"), m.get("commonName"), d.get("subjectAltName"), d.get("issuer")


def whois_apnic(ip):
    """归属查询：target349.example.com:43 原始 TCP。"""
    s = socket.create_connection(("target349.example.com", 43), timeout=15)
    s.sendall(("%s\r\n" % ip).encode())
    buf = b""
    while True:
        c = s.recv(8192)
        if not c:
            break
        buf += c
    s.close()
    return buf.decode("utf-8", "replace")
