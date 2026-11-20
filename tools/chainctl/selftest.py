#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""chainctl 端到端自测：在 127.0.0.1 起一组"故意有漏洞"的夹具，跑完整链路并断言。

证明的是：探针 → 利用 → 双对照复测 → 通道接管 → 后渗透枚举 这条链在本机真跑得通，
而不是"看起来应该能跑"。任何断言失败即退出码 3。

覆盖：
  HTTP  Web: 指纹 / actuator 配置泄露 / Tomcat PUT 落地 / SQLi oracle / LFI+environ 投毒 /
             上传绕过矩阵 / JWT alg:none / JWT 弱密钥 / SSRF 差分 / WebShell 通道 / 后渗透枚举
  TCP       Redis(未授权 + 需认证负例) / MongoDB / memcached / PHP-FPM FastCGI 未授权 RCE
"""
from __future__ import annotations

import base64
import json
import os
import re
import socket
import struct
import sys
import tempfile
import threading
import time
import urllib.parse
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import chainlib as L

FAILS: list[str] = []


def check(cond: bool, msg: str) -> None:
    if cond:
        print("  [ok]  " + msg)
    else:
        print("  [FAIL] " + msg)
        FAILS.append(msg)


# ---------------------------------------------------------------------------
# HTTP 夹具
# ---------------------------------------------------------------------------
class WebFixture(BaseHTTPRequestHandler):
    protocol_version = "HTTP/1.1"
    files: dict = {}
    uploads: dict = {}

    def log_message(self, *a):  # 静音
        pass

    def _send(self, code: int, body: bytes | str, ctype="text/html; charset=utf-8", extra=None):
        if isinstance(body, str):
            body = body.encode()
        self.send_response(code)
        self.send_header("Content-Type", ctype)
        self.send_header("Content-Length", str(len(body)))
        self.send_header("Server", "nginx/1.24.0 (fixture)")
        for k, v in (extra or {}).items():
            self.send_header(k, v)
        self.end_headers()
        self.wfile.write(body)

    # --- 通道 contract：POST k=<key>&c=b64(cmd) → stdout 原文 ----------------
    def _contract_exec(self, key_expected="artex"):
        n = int(self.headers.get("Content-Length") or 0)
        raw = self.rfile.read(n).decode("utf-8", "replace")
        form = urllib.parse.parse_qs(raw)
        k = (form.get("k") or [""])[0]
        if k != key_expected:
            return self._send(403, "forbidden")
        cmd = base64.b64decode((form.get("c") or [""])[0]).decode("utf-8", "replace")
        out = self.fake_exec(cmd)
        self._send(200, out, "text/plain; charset=utf-8")

    @staticmethod
    def fake_exec(cmd: str) -> str:
        if cmd.startswith("echo "):
            return cmd[5:].strip()
        if cmd.strip() == "id":
            return "uid=33(www-data) gid=33(www-data) groups=33(www-data)"
        if "id;" in cmd or cmd.startswith("id"):
            return ("uid=33(www-data) gid=33(www-data) groups=33(www-data)\n"
                    "web01\nLinux web01 5.15.0-91-generic #101-Ubuntu SMP x86_64 GNU/Linux\n"
                    "PRETTY_NAME=\"Ubuntu 22.04.3 LTS\"")
        if "sudo -n -l" in cmd:
            return "Sorry, user www-data may not run sudo on web01."
        if "perm -4000" in cmd:
            return "/usr/bin/sudo\n/usr/bin/passwd\n/usr/bin/chsh"
        if "/etc/hosts" in cmd:
            return "127.0.0.1 localhost\n10.1.2.3 intranet-db\n10.1.2.4 intranet-ad"
        if "grep -iE" in cmd:
            return "DB_PASSWORD=hunter2\nAPI_TOKEN=tok_live_9f2"
        if "ip -br a" in cmd:
            return "eth0 UP 10.0.0.5/24 fe80::1/64"
        if "/etc/cron" in cmd:
            return "-rw-r--r-- 1 root root 1024 Jan 1 00:00 /etc/crontab"
        return "ok"

    def do_OPTIONS(self):
        self.send_response(200)
        self.send_header("Allow", "GET, POST, PUT, OPTIONS")
        self.send_header("Content-Length", "0")
        self.end_headers()

    def do_PUT(self):
        n = int(self.headers.get("Content-Length") or 0)
        body = self.rfile.read(n)
        WebFixture.files[self.path] = body
        self._send(201, "")

    def do_POST(self):
        p = urllib.parse.urlparse(self.path).path
        if p.endswith(".jsp") or p.endswith(".php"):
            return self._contract_exec()
        if p == "/upload":
            n = int(self.headers.get("Content-Length") or 0)
            body = self.rfile.read(n)
            m = re.search(rb'filename="([^"]+)"', body)
            fname = m.group(1).decode() if m else ""
            # 典型白名单缺陷：检查原始名后缀，存储时取分号前（IIS 语义）
            if not fname.lower().endswith((".jpg", ".png", ".txt")):
                return self._send(403, "extension not allowed")
            stored = fname.split(";")[0]
            WebFixture.uploads["/uploads/" + stored] = body
            return self._send(200, json.dumps({"url": "/uploads/" + stored}), "application/json")
        if p == "/__contract__":
            return self._contract_exec()
        self._send(404, "not found")

    def do_GET(self):
        u = urllib.parse.urlparse(self.path)
        p = u.path
        q = urllib.parse.parse_qs(u.query)
        if p == "/":
            if "path" in q:
                path = q["path"][0]
                ua = self.headers.get("User-Agent", "")
                if path.endswith("/proc/self/environ") or path == "/proc/self/environ":
                    if "<?php" in ua:
                        # 模拟 include 执行：<?php ?> 段被解析执行，UA 原文不再反射
                        m = re.search(r"echo '([^']+)'", ua)
                        env = ("PATH=/usr/bin\nGATEWAY_INTERFACE=CGI/1.1\nHTTP_ACCEPT=*/*\n"
                               + (m.group(1) if m else "")
                               + "\nuid=33(www-data) gid=33(www-data)\n")
                    else:
                        env = ("PATH=/usr/bin\nGATEWAY_INTERFACE=CGI/1.1\n"
                               "HTTP_USER_AGENT=%s\nHTTP_ACCEPT=*/*\n" % ua)
                    return self._send(200, env, "text/plain")
                if "passwd" in path:
                    return self._send(200, "root:x:0:0:root:/root:/bin/bash\nwww-data:x:33:33::/var/www:/usr/sbin/nologin\n",
                                      "text/plain")
            if "id" in q:
                v = q["id"][0]
                if "extractvalue" in v:
                    return self._send(200, "SQLSTATE: XPATH syntax error: '~8.0.32-artex'")
                if "and '1'='2" in v.replace(" ", "") or "1'='2" in v:
                    return self._send(200, "<html>no rows found</html>")
                if v.endswith("''"):
                    return self._send(200, "<html>user=guest uid=9 email=guest@corp.local</html>")
                if v.endswith("'"):
                    return self._send(200, "<html>SQL syntax error near '' at line 1</html>")
                if "1'='2" in v.replace(" ", "") or "1'='2" in v:
                    return self._send(200, "<html>no rows found</html>")
                if "1'='1" in v:
                    # 真实 SQL 语义：1 and '1'='1' 与原始值 1 同结果集
                    return self._send(200, "<html>user=guest</html>")
                return self._send(200, "<html>user=guest</html>")
            if "url" in q:
                url = q["url"][0]
                if re.search(r"127\.0\.0\.1:\d+/canary", url):
                    return self._send(200, "chain-canary-ok internal service")
                return self._send(200, "fetch error: connection refused")
            return self._send(200, "<html><head><title>ARTEX Fixture Portal</title></head>"
                                   "<body>welcome</body></html>")
        if p == "/actuator/env":
            return self._send(200, json.dumps({
                "activeProfiles": ["prod"],
                "propertySources": [{"name": "applicationConfig: [classpath:/application.properties]",
                                     "properties": {"spring.datasource.url": {"value": "jdbc:mysql://db:3306/app"},
                                                    "spring.datasource.password": {"value": "S3cr3t!pass"},
                                                    "jwt.signing.key": {"value": "k3y-0123456789"}}}]}),
                "application/json")
        if p == "/actuator/__nope__":
            return self._send(404, "not found")
        if p == "/me":
            auth = self.headers.get("Authorization", "")
            tok = self.headers.get("Cookie", "").replace("token=", "")
            if not tok:
                tok = auth.replace("Bearer ", "").strip()
            parts = tok.split(".")
            if len(parts) != 3:
                return self._send(401, "unauthorized")
            try:
                hdr = json.loads(L._b64u_dec(parts[0]))
            except Exception:
                return self._send(401, "bad jwt")
            if hdr.get("alg") == "none":
                return self._send(200, json.dumps({"role": "admin", "sub": "forged"}), "application/json")
            sig = L._b64u(L.hmac.new(b"secret", ("%s.%s" % (parts[0], parts[1])).encode(),
                                     L.hashlib.sha256).digest())
            if sig != parts[2]:
                return self._send(401, "bad signature")
            return self._send(200, json.dumps({"role": "user", "sub": "u1"}), "application/json")
        if p == "/shell.php" or p == "/shell.jsp":
            return self._contract_exec()
        if p in WebFixture.files:
            body = WebFixture.files[p]
            if p.endswith(".jsp") or p.endswith(".php"):
                if not body.strip().startswith(b"<"):
                    return self._send(200, body, "application/octet-stream")
                return self._contract_exec()
            return self._send(200, body, "application/octet-stream")
        if p in WebFixture.uploads:
            body = WebFixture.uploads[p]
            if p.endswith(".php"):
                return self._contract_exec()
            return self._send(200, body, "application/octet-stream")
        if p.startswith("/uploads/"):
            return self._send(404, "no such upload")
        self._send(404, "not found")


# ---------------------------------------------------------------------------
# TCP 夹具
# ---------------------------------------------------------------------------
class TcpFixture:
    def __init__(self, handler):
        self.sock = socket.socket()
        self.sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        self.sock.bind(("127.0.0.1", 0))
        self.sock.listen(16)
        self.port = self.sock.getsockname()[1]
        self.handler = handler
        self.stop = False
        threading.Thread(target=self._loop, daemon=True).start()

    def _loop(self):
        while not self.stop:
            try:
                c, _ = self.sock.accept()
            except OSError:
                return
            threading.Thread(target=self._serve, args=(c,), daemon=True).start()

    def _serve(self, c: socket.socket):
        try:
            c.settimeout(3)
            self.handler(c)
        except Exception:
            pass
        finally:
            try:
                c.close()
            except Exception:
                pass

    def close(self):
        self.stop = True
        try:
            self.sock.close()
        except Exception:
            pass


def redis_handler(auth: bool):
    # 状态必须跨连接共享（redis_cmd 每条命令一个连接）
    state = {"dir": "/var/lib/redis", "kv": {}}

    def h(c):
        f = c.makefile("rb")
        while True:
            line = f.readline()
            if not line:
                return
            if not line.startswith(b"*"):
                return
            n = int(line[1:].strip())
            args = []
            for _ in range(n):
                ln = int(f.readline()[1:].strip())
                args.append(f.read(ln + 2)[:ln].decode())
            cmd = args[0].lower() if args else ""
            if auth and cmd != "auth":
                c.sendall(b"-NOAUTH Authentication required.\r\n")
                continue
            if cmd == "info":
                body = ("# Server\r\nredis_version:7.2.4\r\nredis_mode:standalone\r\n").encode()
                c.sendall(b"$%d\r\n" % len(body) + body + b"\r\n")
            elif cmd == "config":
                c.sendall(b"*2\r\n$3\r\ndir\r\n$%d\r\n%s\r\n" % (len(state["dir"]), state["dir"].encode()))
            elif cmd == "set":
                state["kv"][args[1]] = args[2]
                c.sendall(b"+OK\r\n")
            elif cmd == "get":
                v = state["kv"].get(args[1])
                if v is None:
                    c.sendall(b"$-1\r\n")
                else:
                    c.sendall(b"$%d\r\n%s\r\n" % (len(v), v.encode()))
            elif cmd == "del":
                state["kv"].pop(args[1], None)
                c.sendall(b":1\r\n")
            else:
                c.sendall(b"-ERR unknown command\r\n")
    return h


def memcache_handler(c):
    data = c.recv(4096).decode("utf-8", "replace")
    if "stats" in data:
        c.sendall(b"STAT curr_items 12\r\nSTAT total_items 99\r\nSTAT version 1.6.21\r\nEND\r\n")
        return
    if "version" in data:
        c.sendall(b"VERSION 1.6.21\r\n")
        return
    # 未知命令必须回 ERROR（memcached 协议语义的阴性对照）
    c.sendall(b"ERROR\r\n")


def mongo_handler(c):
    head = b""
    while len(head) < 16:
        chunk = c.recv(16 - len(head))
        if not chunk:
            return
        head += chunk
    total, _r, _res, _op = struct.unpack("<iiii", head)
    rest = b""
    while len(rest) < total - 16:
        chunk = c.recv(total - 16 - len(rest))
        if not chunk:
            break
        rest += chunk
    doc = rest[5:]  # flagBits(4)+kind(1) 后即请求 BSON

    def bson(d: dict) -> bytes:
        body = b""
        for k, v in d.items():
            if isinstance(v, bool):
                body += b"\x08" + k.encode() + b"\x00" + (b"\x01" if v else b"\x00")
            elif isinstance(v, int):
                body += b"\x10" + k.encode() + b"\x00" + struct.pack("<i", v)
            elif isinstance(v, float):
                body += b"\x01" + k.encode() + b"\x00" + struct.pack("<d", v)
            elif isinstance(v, str):
                eb = v.encode() + b"\x00"
                body += b"\x02" + k.encode() + b"\x00" + struct.pack("<i", len(eb)) + eb
            elif isinstance(v, list) and v and isinstance(v[0], dict):
                sub = b""
                for i, item in enumerate(v):
                    inner = bson(item)
                    sub += b"\x03" + str(i).encode() + b"\x00" + inner
                sub = struct.pack("<i", len(sub) + 5) + sub + b"\x00"
                body += b"\x04" + k.encode() + b"\x00" + sub
        return struct.pack("<i", len(body) + 5) + body + b"\x00"

    if b"listDatabases" in doc:
        reply = {"ok": 1.0, "databases": [{"name": "admin", "sizeOnDisk": 4096},
                                          {"name": "app", "sizeOnDisk": 8192}]}
    else:
        reply = {"ok": 1.0, "maxWireVersion": 17, "isWritablePrimary": True}
    b = bson(reply)
    # 标准 OP_MSG kind-0 应答：flagBits(4) + kind(1) + BSON 文档（无集合名段）
    payload = struct.pack("<i", 0) + b"\x00" + b
    c.sendall(struct.pack("<iiii", 16 + len(payload), 1, 1, 2013) + payload)


def fcgi_handler(c):
    buf = b""
    params = b""
    stdin = b""
    while True:
        try:
            chunk = c.recv(65536)
        except socket.timeout:
            break
        if not chunk:
            break
        buf += chunk
        done = False
        while len(buf) >= 8:
            _v, rtype, rid, clen, plen, _res = struct.unpack(">BBHHBB", buf[:8])
            if len(buf) < 8 + clen + plen:
                break
            body = buf[8:8 + clen]
            buf = buf[8 + clen + plen:]
            if rtype == L.FCGI_PARAMS:
                params += body
            elif rtype == L.FCGI_STDIN:
                stdin += body
                if clen == 0:
                    done = True
        if done:
            break
    text = stdin.decode("utf-8", "replace")
    if b"auto_prepend_file=php://input" in params or b"allow_url_include=1" in params:
        m = re.search(r"echo '([^']+)'", text)
        out = (m.group(1) if m else "") + "\nuid=33(www-data) gid=33(www-data)\n"
    elif b"/etc/passwd" in params:
        out = "Primary script unknown"
    else:
        out = "Primary script unknown"
    outb = out.encode()
    pad = (8 - len(outb) % 8) % 8
    c.sendall(struct.pack(">BBHHBB", 1, L.FCGI_STDOUT, 1, len(outb), pad, 0) + outb + b"\x00" * pad)
    c.sendall(struct.pack(">BBHHBB", 1, L.FCGI_END, 1, 8, 0, 0) +
              struct.pack(">IBxxx", 0, 0))


# ---------------------------------------------------------------------------
# 断言
# ---------------------------------------------------------------------------
def load_ledger(path: str) -> list[dict]:
    with open(path, "r", encoding="utf-8") as f:
        return [json.loads(l) for l in f if l.strip()]


def got(ledger, pid, action, ok=None):
    for r in ledger:
        if r["primitive"] == pid and r["action"] == action and (ok is None or r["ok"] == ok):
            return r
    return None


def main() -> int:
    t0 = time.time()
    tmp = tempfile.mkdtemp(prefix="chainctl-selftest-")
    print("chainctl selftest · 夹具目录", tmp)

    # --- 起 HTTP 夹具 -----------------------------------------------------
    httpd = ThreadingHTTPServer(("127.0.0.1", 0), WebFixture)
    web_port = httpd.server_address[1]
    threading.Thread(target=httpd.serve_forever, daemon=True).start()
    base = f"http://127.0.0.1:{web_port}"

    class Canary(BaseHTTPRequestHandler):
        protocol_version = "HTTP/1.1"

        def log_message(self, *a):
            pass

        def do_GET(self):
            body = b"chain-canary-ok internal"
            self.send_response(200)
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)

    canary = ThreadingHTTPServer(("127.0.0.1", 0), Canary)
    canary_port = canary.server_address[1]
    threading.Thread(target=canary.serve_forever, daemon=True).start()

    # --- 起 TCP 夹具 ------------------------------------------------------
    f_redis = TcpFixture(redis_handler(auth=False))
    f_redis_auth = TcpFixture(redis_handler(auth=True))
    f_mem = TcpFixture(memcache_handler)
    f_mongo = TcpFixture(mongo_handler)
    f_fcgi = TcpFixture(fcgi_handler)
    print(f"  夹具端口 web={web_port} redis={f_redis.port} redis-auth={f_redis_auth.port} "
          f"mem={f_mem.port} mongo={f_mongo.port} fcgi={f_fcgi.port} canary={canary_port}")

    # --- 构造 JWT（HS256 / secret） ---------------------------------------
    hs_token = L.jwt_forge({"alg": "HS256", "typ": "JWT"},
                           {"sub": "u1", "role": "user", "exp": 9999999999}, b"secret")
    wrong_token = L.jwt_forge({"alg": "HS256", "typ": "JWT"},
                              {"sub": "u1", "role": "user"}, b"not-in-wordlist-zzz")

    hints = {
        "lfi_param": "path", "sqli_param": "id",
        "upload_url": "/upload", "upload_field": "file", "upload_dir_url": "/uploads/",
        "token": hs_token, "jwt_probe_url": "/me",
        "ssrf_param": "url", "canary_url": f"http://127.0.0.1:{canary_port}/canary",
        "shell_url": base + "/shell.php", "shell_key": "artex",
    }

    # ============================ 1) Web 链路 ============================
    print("\n[1] Web 链路（plan/profile=web → probe → exploit → verify）")
    plan = L.build_plan(base, "web", dict(hints))
    r = L.run_plan(plan, os.path.join(tmp, "web"), ["127.0.0.1"], execute=True, rps=1000)
    led = load_ledger(r["ledger"])
    print(f"  记录 {len(led)} 条")
    for pid, want in [
        ("spring-actuator", True), ("tomcat-put", True), ("sqli-oracle", True),
        ("lfi-environ-rce", True), ("upload-bypass", True), ("jwt-none", True),
        ("jwt-weak-secret", True), ("ssrf-diff", True),
    ]:
        e = got(led, pid, "exploit")
        check(bool(e) and e["ok"] == want, f"{pid} exploit 命中（{e['detail'][:60] if e else '无记录'}）")
        v = got(led, pid, "verify")
        check(bool(v) and v["ok"], f"{pid} verify 双对照成立（{v['detail'][:60] if v else '无记录'}）")

    # 证据落盘
    check(any(n.endswith(".json") for n in os.listdir(os.path.join(tmp, "web"))), "证据文件已落盘")

    # WebShell 通道 + 后渗透（同一台 web 夹具）
    print("\n[2] 通道接管与后渗透（shell-channel → post-unix-enum → lateral）")
    plan2 = {"target": base, "profile": "host", "hints": hints,
             "primitives": ["shell-channel", "post-unix-enum", "lateral-candidates"]}
    r2 = L.run_plan(plan2, os.path.join(tmp, "post"), ["127.0.0.1"], execute=True, rps=1000)
    led2 = load_ledger(r2["ledger"])
    ch = got(led2, "shell-channel", "exploit")
    check(bool(ch) and ch["ok"] and ch["gained"] == "shell-channel", "shell-channel 通道建立并实测")
    pe = got(led2, "post-unix-enum", "exploit")
    check(bool(pe) and pe["ok"], "post-unix-enum 完成（id/SUID/sudo/凭据/邻居）")
    lat = got(led2, "lateral-candidates", "exploit")
    blob = json.dumps([pe.get("evidence", {}) if pe else {}, lat.get("evidence", {}) if lat else {}])
    check(bool(lat) and pe and "10.1.2.3" in blob,
          "lateral 候选识别出内网邻居 10.1.2.3/10.1.2.4")

    # ============================ 3) DB/服务链路 =========================
    print("\n[3] 未授权服务族（Redis / MongoDB / memcached / Docker API / ES / PHP-FPM）")
    db_cases = [
        ("unauth-redis", f_redis.port, {"gained": "db-access+file-write"}),
        ("unauth-mongo", f_mongo.port, {"gained": "db-access"}),
        ("unauth-memcache", f_mem.port, {"gained": "db-access"}),
        ("php-fpm-unauth", f_fcgi.port, {"gained": "shell-channel", "hints": {"php_file": "/var/www/html/index.php"}}),
    ]
    db_dir = os.path.join(tmp, "db")
    for pid, port, exp in db_cases:
        h = dict(hints)
        h.update(exp.get("hints") or {})
        p = {"target": f"http://127.0.0.1:{port}", "profile": "db", "hints": h, "primitives": [pid]}
        rr = L.run_plan(p, db_dir, ["127.0.0.1"], execute=True, rps=1000)
        l = load_ledger(rr["ledger"])
        e = got(l, pid, "exploit")
        check(bool(e) and e["ok"] and e["gained"] == exp["gained"],
              f"{pid}:{port} → {e['detail'][:70] if e else '无记录'}")
        v = got(l, pid, "verify")
        check(bool(v) and v["ok"], f"{pid} verify 复测成立")

    # Docker / ES 用 HTTP 结构（本机用 web 夹具模拟对应端点）
    class DockerHandler(WebFixture):
        created: dict = {}

        def do_GET(self):
            p = urllib.parse.urlparse(self.path).path
            if p == "/version":
                return self._send(200, json.dumps({"ApiVersion": "1.43", "Version": "24.0.7"}),
                                  "application/json")
            if p == "/containers/json":
                containers = [{"Id": "abc123", "Names": ["/app"]}] + [
                    {"Id": cid, "Names": ["/" + name]} for name, cid in DockerHandler.created.items()]
                return self._send(200, json.dumps(containers), "application/json")
            if p == "/__nope__":
                return self._send(404, "not found")
            return self._send(404, "not found")

        def do_POST(self):
            p = urllib.parse.urlparse(self.path).path
            if p.startswith("/containers/create"):
                qs = urllib.parse.parse_qs(urllib.parse.urlparse(self.path).query)
                name = (qs.get("name") or ["unnamed"])[0]
                cid = "deadbeef" * 8
                DockerHandler.created[name] = cid
                return self._send(201, json.dumps({"Id": cid, "Warnings": []}), "application/json")
            if p.startswith("/containers/") and p.endswith("/start"):
                self.send_response(204)
                self.send_header("Content-Length", "0")
                self.end_headers()
                return
            self._send(404, "not found")

    class EsHandler(WebFixture):
        def do_GET(self):
            p = urllib.parse.urlparse(self.path).path
            if p == "/":
                return self._send(200, json.dumps({"cluster_name": "artex-es",
                                                   "version": {"number": "8.11.1"}}), "application/json")
            if p == "/_cat/indices":
                return self._send(200, json.dumps([{"index": "users"}, {"index": "orders"}]),
                                  "application/json")
            if p.endswith("/_count"):
                return self._send(200, json.dumps({"count": 1200000}), "application/json")
            return self._send(404, "not found")

    for cls, pid, port_expected in [(DockerHandler, "unauth-docker", 2375),
                                    (EsHandler, "unauth-es", 9200)]:
        srv = ThreadingHTTPServer(("127.0.0.1", 0), cls)
        port = srv.server_address[1]
        threading.Thread(target=srv.serve_forever, daemon=True).start()
        h = dict(hints)
        h["port"] = port
        p = {"target": f"http://127.0.0.1:{port}", "profile": "db", "hints": h, "primitives": [pid]}
        rr = L.run_plan(p, db_dir, ["127.0.0.1"], execute=True, rps=1000)
        l = load_ledger(rr["ledger"])
        e = got(l, pid, "exploit")
        check(bool(e) and e["ok"], f"{pid}:{port} → {e['detail'][:70] if e else '无记录'}")
        v = got(l, pid, "verify")
        check(bool(v) and v["ok"], f"{pid} verify 复测成立")

    # ============================ 4) 负例 ================================
    print("\n[4] 负例（必须判否，防假阳性）")
    p = {"target": f"http://127.0.0.1:{f_redis_auth.port}", "profile": "db",
         "hints": dict(hints), "primitives": ["unauth-redis"]}
    rr = L.run_plan(p, os.path.join(tmp, "neg"), ["127.0.0.1"], execute=True, rps=1000)
    l = load_ledger(rr["ledger"])
    e = got(l, "unauth-redis", "probe")
    check(bool(e) and not e["ok"], f"Redis 需认证负例判否（{e['detail'][:50] if e else '无记录'}）")

    h2 = dict(hints)
    h2["token"] = wrong_token
    p = {"target": base, "profile": "web", "hints": h2, "primitives": ["jwt-weak-secret"]}
    rr = L.run_plan(p, os.path.join(tmp, "neg"), ["127.0.0.1"], execute=True, rps=1000)
    l = load_ledger(rr["ledger"])
    e = got(l, "jwt-weak-secret", "exploit")
    check(bool(e) and not e["ok"], f"JWT 未知密钥负例判否（{e['detail'][:50] if e else '无记录'}）")

    # 范围护栏
    print("\n[5] 范围护栏")
    try:
        L.run_plan({"target": "http://169.254.169.254", "profile": "web", "hints": {}, "primitives": []},
                   os.path.join(tmp, "scope"), ["169.254.169.254"], execute=False)
        check(False, "云元数据应被拒绝")
    except L.ScopeError as e:
        check("红线" in str(e), f"云元数据被拒：{e}")
    try:
        L.run_plan({"target": "http://10.9.9.9", "profile": "web", "hints": {}, "primitives": []},
                   os.path.join(tmp, "scope"), ["10.0.0.0/24"], execute=False)
        check(False, "清单外目标应被拒绝")
    except L.ScopeError as e:
        check("不在授权清单" in str(e), f"清单外目标被拒：{e}")

    # ============================ 汇总 ==================================
    for x in (httpd, canary):
        x.shutdown()
    for f in (f_redis, f_redis_auth, f_mem, f_mongo, f_fcgi):
        f.close()

    print("\n" + "=" * 60)
    if FAILS:
        print(f"FAILED {len(FAILS)} 项：")
        for f in FAILS:
            print("  - " + f)
        print(f"耗时 {time.time()-t0:.1f}s，夹具目录 {tmp}")
        return 3
    print(f"PASS 全部断言通过，耗时 {time.time()-t0:.1f}s")
    print(f"证据与台账：{tmp}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
