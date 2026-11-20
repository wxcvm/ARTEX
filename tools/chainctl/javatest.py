#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Java 中间件族原语冒烟测试（本地夹具，127.0.0.1）。"""
import base64, json, os, sys, threading, urllib.parse
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import socket, tempfile, time

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import chainlib as L

FAILS = []
def check(cond, msg):
    print(("  [ok]  " if cond else "  [FAIL] ") + msg)
    if not cond:
        FAILS.append(msg)

# --- 1) T3 夹具 --------------------------------------------------------------
def t3_handler(c):
    data = c.recv(4096)
    if data.startswith(b"t3 "):
        c.sendall(b"HELO:10.3.6.0.false\r\n\r\n\x01\x00\x0a")
    else:
        c.sendall(b"Invalid header\r\n")
    c.close()

# --- 2) Java-Web 夹具（shiro/struts2/fastjson/druid/jenkins） ------------------
class JavaFixture(BaseHTTPRequestHandler):
    protocol_version = "HTTP/1.1"
    SHIRO_KEY = base64.b64decode("kPH+bIxk5D2deZiIxcaaaA==")

    def log_message(self, *a):
        pass

    def _send(self, code, body, ctype="text/html; charset=utf-8", extra=None):
        if isinstance(body, str):
            body = body.encode()
        self.send_response(code)
        self.send_header("Content-Type", ctype)
        self.send_header("Content-Length", str(len(body)))
        for k, v in (extra or {}).items():
            self.send_header(k, v)
        self.end_headers()
        self.wfile.write(body)

    def do_GET(self):
        p = urllib.parse.urlparse(self.path).path
        cookie = self.headers.get("Cookie", "")
        if p == "/login.jsp":
            # Shiro 模拟：rememberMe cookie 解密失败 → deleteMe；成功 → 无 deleteMe
            if "rememberMe=" in cookie:
                val = cookie.split("rememberMe=")[1].split(";")[0]
                ok = False
                try:
                    raw = base64.b64decode(val)
                    pt = _aes_decrypt(self.SHIRO_KEY, raw[:16], raw[16:])
                    ok = pt[:4] == b"\xac\xed\x00\x05"
                except Exception:
                    ok = False
                if ok:
                    return self._send(200, "welcome back")
                return self._send(200, "login page", extra={"Set-Cookie": "rememberMe=deleteMe; Path=/"})
            return self._send(200, "login page")
        if p == "/api":
            return self._send(200, "api root")
        if p == "/druid/index.html":
            return self._send(200, "<html><head><title>Druid Stat Index</title></head><body>druid</body></html>")
        if p == "/druid/weburi.json":
            return self._send(200, json.dumps({"ResultCode": 1, "Content": ["/api", "/admin"]})),
        if p in ("/druid/weburi.json", "/druid/datasource.json"):
            if p == "/druid/weburi.json":
                return self._send(200, json.dumps({"ResultCode": 1, "Content": ["/api", "/admin"]}), "application/json")
            return self._send(200, json.dumps({"ResultCode": 1, "Content": [{"name": "jdbc:mysql://db/app"}]}), "application/json")
        if p == "/druid/sql.json":
            return self._send(200, json.dumps({"ResultCode": 1, "Content": [{"SQL": "select * from users"}]}), "application/json")
        if p == "/script":
            return self._send(200, '<html><body><form name="script"><textarea name="script"></textarea></form>Groovy Script Console</body></html>')
        if p == "/crumbIssuer/api/xml":
            return self._send(200, '<crumbRequestField>Jenkins-Crumb</crumbRequestField><crumb>crumb123</crumb>')
        if p == "/json":
            return self._send(200, "parsed")
        self._send(404, "not found")

    def do_POST(self):
        p = urllib.parse.urlparse(self.path).path
        n = int(self.headers.get("Content-Length") or 0)
        body = self.rfile.read(n).decode("utf-8", "replace")
        ct = self.headers.get("Content-Type", "")
        # --- S2-045：错误页回显 Content-Type（%{expr} 已求值） ---
        if p == "/upload.action" or p == "/":
            if "%{" in ct:
                import re as _re
                m = _re.search(r"%\{(\d+)\*(\d+)\}", ct)
                if m:
                    res = str(int(m.group(1)) * int(m.group(2)))
                    return self._send(400, f"the request doesn't contain a multipart/form-data stream; content type: {ct} result {res}")
                return self._send(400, f"bad content type {ct}")
            return self._send(200, "ok")
        # --- fastjson ---
        if p == "/json":
            if body == '{"a":1':
                return self._send(400, "com.alibaba.fastjson.JSONException: expect ':' at 7")
            try:
                j = json.loads(body)
                if "@type" in j:
                    t = j["@type"]
                    if t == "com.alibaba.fastjson.JSONObject":
                        return self._send(400, "com.alibaba.fastjson.JSONException: autoType is not support")
                    if t == "com.sun.rowset.JdbcRowSetImpl":
                        # 模拟 <=1.2.24：直接实例化不报错
                        return self._send(200, "parsed")
                return self._send(200, "parsed")
            except Exception:
                return self._send(400, "parse error")
        # --- Jenkins /script ---
        if p == "/script":
            crumb = self.headers.get("Jenkins-Crumb", "")
            if crumb != "crumb123":
                return self._send(403, "CSRF crumb missing")
            form = urllib.parse.parse_qs(body)
            script = (form.get("script") or [""])[0]
            out = ""
            m = None
            import re as _re
            m = _re.search(r'println "([^"]*)"', script)
            if m:
                out = m.group(1)
            return self._send(200, out, "text/plain")
        self._send(404, "not found")


def _aes_decrypt(key, iv, ct):
    # 解密方向（测试夹具用，纯标准库）：AES-128-CBC + PKCS5 去填充
    sbox = L._aes_sbox()
    inv_sbox = [0] * 256
    for i, v in enumerate(sbox):
        inv_sbox[v] = i
    w = L._aes128_expand(key)
    gf = L._gf_mul

    def inv_subshift(s):
        t = [inv_sbox[b] for b in s]
        out = [0] * 16
        for c in range(4):
            for r in range(4):
                out[4 * c + r] = t[4 * ((c - r) % 4) + r]
        return out

    def inv_mixcol(s):
        out = [0] * 16
        for c in range(4):
            a = s[4 * c: 4 * c + 4]
            out[4 * c + 0] = gf(a[0], 14) ^ gf(a[1], 11) ^ gf(a[2], 13) ^ gf(a[3], 9)
            out[4 * c + 1] = gf(a[0], 9) ^ gf(a[1], 14) ^ gf(a[2], 11) ^ gf(a[3], 13)
            out[4 * c + 2] = gf(a[0], 13) ^ gf(a[1], 9) ^ gf(a[2], 14) ^ gf(a[3], 11)
            out[4 * c + 3] = gf(a[0], 11) ^ gf(a[1], 13) ^ gf(a[2], 9) ^ gf(a[3], 14)
        return out

    def addrk(s, rnd):
        for c in range(4):
            for r in range(4):
                s[4 * c + r] ^= w[4 * rnd + c][r]
        return s

    def dec_block(b):
        s = list(b)
        s = addrk(s, 10)
        for rnd in range(9, 0, -1):
            s = inv_mixcol(addrk(inv_subshift(s), rnd))
        s = addrk(inv_subshift(s), 0)
        return bytes(s)

    out = b""
    prev = iv
    for i in range(0, len(ct), 16):
        blk = ct[i:i + 16]
        out += bytes(x ^ y for x, y in zip(dec_block(blk), prev))
        prev = blk
    pad = out[-1]
    if 1 <= pad <= 16 and out[-pad:] == bytes([pad]) * pad:
        out = out[:-pad]
    return out


def main():
    tmp = tempfile.mkdtemp(prefix="chainctl-javatest-")
    # HTTP 夹具
    httpd = ThreadingHTTPServer(("127.0.0.1", 0), JavaFixture)
    web_port = httpd.server_address[1]
    threading.Thread(target=httpd.serve_forever, daemon=True).start()
    base = f"http://127.0.0.1:{web_port}"
    # T3 夹具
    tsrv = socket.socket()
    tsrv.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    tsrv.bind(("127.0.0.1", 0))
    tsrv.listen(4)
    t3_port = tsrv.getsockname()[1]
    def t3_loop():
        while True:
            try:
                c, _ = tsrv.accept()
            except OSError:
                return
            threading.Thread(target=t3_handler, args=(c,), daemon=True).start()
    threading.Thread(target=t3_loop, daemon=True).start()

    # ---- shiro（probe/exploit/verify 全链） ----
    plan = {"target": base, "profile": "web", "hints": {"shiro_login_url": "/login.jsp"},
            "primitives": ["shiro-rememberme"]}
    r = L.run_plan(plan, os.path.join(tmp, "shiro"), ["127.0.0.1"], execute=True, rps=1000)
    led = [json.loads(l) for l in open(r["ledger"])]
    for act, want in [("probe", True), ("exploit", True), ("verify", True)]:
        e = [x for x in led if x["action"] == act]
        check(bool(e) and e[0]["ok"] == want, f"shiro {act} → {e[0]['detail'][:70] if e else '无记录'}")

    # ---- struts2 ----
    plan = {"target": base, "profile": "web", "hints": {"struts_url": "/upload.action"},
            "primitives": ["struts2-ognl"]}
    r = L.run_plan(plan, os.path.join(tmp, "s2"), ["127.0.0.1"], execute=True, rps=1000)
    led = [json.loads(l) for l in open(r["ledger"])]
    for act, want in [("probe", True), ("exploit", True), ("verify", True)]:
        e = [x for x in led if x["action"] == act]
        check(bool(e) and e[0]["ok"] == want, f"struts2 {act} → {e[0]['detail'][:70] if e else '无记录'}")

    # ---- fastjson ----
    plan = {"target": base, "profile": "web", "hints": {"json_url": "/json"},
            "primitives": ["fastjson-rce"]}
    r = L.run_plan(plan, os.path.join(tmp, "fj"), ["127.0.0.1"], execute=True, rps=1000)
    led = [json.loads(l) for l in open(r["ledger"])]
    for act, want in [("probe", True), ("exploit", True), ("verify", True)]:
        e = [x for x in led if x["action"] == act]
        check(bool(e) and e[0]["ok"] == want, f"fastjson {act} → {e[0]['detail'][:70] if e else '无记录'}")

    # ---- druid ----
    plan = {"target": base, "profile": "web", "hints": {},
            "primitives": ["druid-unauth"]}
    r = L.run_plan(plan, os.path.join(tmp, "druid"), ["127.0.0.1"], execute=True, rps=1000)
    led = [json.loads(l) for l in open(r["ledger"])]
    for act, want in [("probe", True), ("exploit", True), ("verify", True)]:
        e = [x for x in led if x["action"] == act]
        check(bool(e) and e[0]["ok"] == want, f"druid {act} → {e[0]['detail'][:70] if e else '无记录'}")

    # ---- jenkins（--execute RCE） ----
    plan = {"target": base, "profile": "web", "hints": {},
            "primitives": ["jenkins-script"]}
    r = L.run_plan(plan, os.path.join(tmp, "jks"), ["127.0.0.1"], execute=True, rps=1000)
    led = [json.loads(l) for l in open(r["ledger"])]
    for act, want in [("probe", True), ("exploit", True), ("verify", True)]:
        e = [x for x in led if x["action"] == act]
        check(bool(e) and e[0]["ok"] == want, f"jenkins {act} → {e[0]['detail'][:70] if e else '无记录'}")

    # ---- weblogic T3 ----
    plan = {"target": f"http://127.0.0.1:{t3_port}", "profile": "db", "hints": {},
            "primitives": ["weblogic-t3iiop"]}
    r = L.run_plan(plan, os.path.join(tmp, "t3"), ["127.0.0.1"], execute=True, rps=1000)
    led = [json.loads(l) for l in open(r["ledger"])]
    for act, want in [("probe", True), ("exploit", True), ("verify", True)]:
        e = [x for x in led if x["action"] == act]
        check(bool(e) and e[0]["ok"] == want, f"t3 {act} → {e[0]['detail'][:70] if e else '无记录'}")

    # ---- 负例：非 T3 端口 ----
    plan = {"target": f"http://127.0.0.1:{web_port}", "profile": "db", "hints": {},
            "primitives": ["weblogic-t3iiop"]}
    r = L.run_plan(plan, os.path.join(tmp, "t3neg"), ["127.0.0.1"], execute=True, rps=1000)
    led = [json.loads(l) for l in open(r["ledger"])]
    e = [x for x in led if x["action"] == "probe"]
    check(bool(e) and not e[0]["ok"], f"t3 负例判否（{e[0]['detail'][:50] if e else '无'}）")

    print("\n" + "=" * 50)
    if FAILS:
        print(f"FAILED {len(FAILS)}")
        return 3
    print("PASS Java 中间件族冒烟测试全通过")
    return 0


if __name__ == "__main__":
    sys.exit(main())
