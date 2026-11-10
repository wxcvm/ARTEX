#!/usr/bin/env python3
# MySQL 未授权/空口令 只读探测 + 取证（不写、不落文件、不上 UDF）
# 用法: python mysql_empty_root_probe.py HOST PORT [user]
import sys, socket, json
HOST = sys.argv[1]; PORT = int(sys.argv[2]); USER = sys.argv[3] if len(sys.argv) > 3 else "root"

def banner():
    s = socket.create_connection((HOST, PORT), timeout=8)
    hdr = s.recv(4)
    ln = hdr[0] | (hdr[1] << 8) | (hdr[2] << 16)
    data = b""
    while len(data) < ln:
        c = s.recv(ln - len(data))
        if not c: break
        data += c
    s.close()
    # protocol 10 -> 0x0A, 之后是 NUL 结尾的版本串
    return data[1:data.index(b"\x00", 1)].decode("latin1", "replace")

print("[banner] MySQL version =", banner())

try:
    import pymysql
    conn = pymysql.connect(host=HOST, port=PORT, user=USER, password="",
                           connect_timeout=8, read_timeout=20, charset="utf8mb4")
    cur = conn.cursor()
    cur.execute("select version(), current_user(), @@hostname, @@secure_file_priv")
    print("[login] OK empty password ->", cur.fetchone())
    for q in ("show grants for '%s'@'%%'" % USER, "show databases"):
        try:
            cur.execute(q)
            print("[%s] %s" % (q[:24], cur.fetchall()))
        except Exception as e:
            print("[%s] ERR %r" % (q[:24], e))
    conn.close()
except Exception as e:
    print("[login] FAILED:", repr(e))
