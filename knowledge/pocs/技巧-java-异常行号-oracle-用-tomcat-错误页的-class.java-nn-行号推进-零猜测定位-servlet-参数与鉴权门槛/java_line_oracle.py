# -*- coding: utf-8 -*-
"""Java 异常行号 oracle —— 零猜测定位 servlet 参数与鉴权门槛
用法：改 endpoint / CLASS / BASE / CANDS 后直接跑。
"""
import sys, re, socket, ssl, os, json, time
for _k in ("http_proxy","https_proxy","all_proxy","HTTP_PROXY","HTTPS_PROXY","ALL_PROXY"):
    os.environ.pop(_k, None)
try: sys.stdout.reconfigure(encoding="utf-8", errors="replace")
except Exception: pass

# ============ 目标配置 ============
IP, PORT, SNI_HOST = "203.0.113.57", 443, "203.0.113.57"
endpoint = "/LHService/tp"
CLASS    = "PutTp"            # 从无参响应里看到的类名（<pre> 里的 Xxx.java:NN）
BASE     = ""                # 参数钥匙串：如 "a=1&b=1&c=1&..."（可留空）
# ==================================

UA = ("Mozilla/5.0 (Linux; Android 13; 23013RK75C) AppleWebKit/537.36 "
      "(KHTML, like Gecko) Chrome/203.0.113.30 Mobile Safari/537.36")

def raw(method="GET", path="/", body=None, headers=None, timeout=20):
    """原生 socket + 手工 TLS；绝不读 http_proxy 环境变量"""
    hdr = {"Host": SNI_HOST, "User-Agent": UA, "Accept": "*/*",
           "Accept-Encoding": "identity", "Connection": "close"}
    if headers: hdr.update(headers)
    data = b""
    if body is not None:
        data = body.encode("utf-8") if isinstance(body, str) else body
        hdr.setdefault("Content-Type", "application/x-www-form-urlencoded")
        hdr["Content-Length"] = str(len(data))
    head = ("%s %s HTTP/1.1\r\n" % (method, path)
            + "".join("%s: %s\r\n" % (k, v) for k, v in hdr.items()) + "\r\n")
    s = socket.socket(); s.settimeout(timeout)
    s.connect((IP, PORT))
    if PORT in (443, 8443):
        ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
        ctx.check_hostname = False; ctx.verify_mode = ssl.CERT_NONE
        try: ctx.set_ciphers("DEFAULT@SECLEVEL=1")
        except Exception: pass
        s = ctx.wrap_socket(s, server_hostname=SNI_HOST)
    s.sendall(head.encode("utf-8", "replace") + data)
    buf = b""
    while True:
        try: c = s.recv(4096)
        except socket.timeout: break
        if not c: break
        buf += c
    s.close()
    return buf.decode("utf-8", "replace")

def line_of(text, cls=None):
    """返回 (行号 or None, 异常首行, 是否 NFE/其它)"""
    cls = cls or CLASS
    m = re.search(cls + r"\.java:(\d+)", text)
    ex = re.search(r"<pre>\s*([^\r\n<]+)", text)
    exs = target348.example.com(1).strip() if ex else ""
    return (int(target22.example.com(1)) if m else None, exs, text)

def probe(tag, path, method="GET", body=None):
    t = raw(method, path, body)
    ln, exs, _ = line_of(t)
    ln2 = ln if ln is not None else ("OTHER" if "Exception" in t else "NOERR")
    print("  %-34s %-5s line=%-7s exc=%s" % (tag, method, ln2, exs[:90]), flush=True)
    return ln2

def build(d, base=None):
    base = BASE if base is None else base
    ks = set(d.keys())
    parts = [x for x in base.split("&") if x and x.split("=")[0] not in ks]
    parts += ["%s=%s" % (k, v) for k, v in d.items()]
    return endpoint + ("?" + "&".join(parts) if parts else "")

def main():
    print("=== 1) 基准行号 ===")
    bare   = probe("裸请求", endpoint)
    keyed  = probe("参数钥匙串", build({})) if BASE else None

    print("\n=== 2) 单参数推进扫描（行号前进 = 参数被读取）===")
    CANDS = ["userid","id","type","flag","head","base64","imgName","imgStr","name",
             "filename","fileName","path","dir","url","tp","pic","img","image",
             "username","token","sign","key","fid","start","page","num","size",
             "body","content","data","file","up","bz","cz","json"]
    ref = keyed if BASE else bare
    advanced = []
    for pn in CANDS:
        ln = probe("附加 " + pn, build({pn: "1"}))
        if ln != ref:
            advanced.append((pn, ln)); print("      ^^^ 推进！", flush=True)
        time.sleep(0.8)
    print("\n=== 3) 结论 ===")
    if advanced:
        print("  被读取的参数: %s" % advanced)
        print("  → 针对这些参数继续给真实值，逼近业务分支")
    else:
        print("  行号停在 %s 不动 ⇒ 该行读的是 session/属性（非请求参数）" % ref)
        print("  → 判定：该 servlet 需要登录态，停止猜参数名")

if __name__ == "__main__":
    main()
