# -*- coding: utf-8 -*-
"""原生 socket HTTP 指纹：完全不走代理（硬清代理环境变量 + 手写 HTTP over socket）。
用法: python httpfp.py <ip:port,ip:port,...> [mode]
  mode=base  -> 只请求 /
  mode=path  -> 请求业务路径
输出 JSON
"""
import socket, ssl, os, sys, json, time, re
from concurrent.futures import ThreadPoolExecutor, as_completed

for k in ("http_proxy","https_proxy","all_proxy","HTTP_PROXY","HTTPS_PROXY","ALL_PROXY"):
    os.environ.pop(k, None)

try:
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")
except Exception:
    pass

UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36"

BIZ_PATHS = ["/PayBmob/", "/UpLoadPic/", "/LHService/", "/Tk_Kj/", "/IMService/",
             "/LHService/v27/routing/health"]

SIGS = [
    ("LHZS_在线充值",   "在线充值"),
    ("LHZS_JSP_index",  "My JSP 'index.jsp' starting page"),
    ("LHZS_fenleilist", "LHService/fenleilist"),
    ("LHZS_LHService",  "LHService"),
    ("LHZS_IMService",  "IMService"),
    ("LHZS_Tk_Kj",      "Tk_Kj"),
    ("LHZS_PayBmob",    "PayBmob"),
    ("LHZS_UpLoadPic",  "UpLoadPic"),
    ("LHZS_六助",       "六助"),
    ("LHZS_lz2025",     "lz2025"),
]

def raw_http(ip, port, path, host=None, https=None, timeout=6.0, maxbytes=8192):
    """手写 HTTP/1.1 over 原生 socket。返回 (status, headers, body, server, raw_head)"""
    if https is None:
        https = port in (443, 8443)
    if host is None:
        host = ip
    req = ("GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\n"
           "Accept: text/html,application/json,*/*\r\n"
           "Accept-Encoding: identity\r\nConnection: close\r\n\r\n" % (path, host, UA))
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(timeout)
    try:
        s.connect((ip, port))
        if https:
            ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
            ctx.check_hostname = False
            ctx.verify_mode = ssl.CERT_NONE
            ctx.set_ciphers("DEFAULT@SECLEVEL=1")
            s = ctx.wrap_socket(s, server_hostname=host)
        s.sendall(req.encode())
        buf = b""
        while len(buf) < maxbytes:
            try:
                chunk = s.recv(4096)
            except socket.timeout:
                break
            if not chunk:
                break
            buf += chunk
        head, _, body = buf.partition(b"\r\n\r\n")
        headtxt = head.decode("utf-8", "replace")
        stat = 0
        m = re.match(r"HTTP/\d\.\d (\d+)", headtxt)
        if m:
            stat = int(target22.example.com(1))
        server = ""
        for line in headtxt.split("\r\n"):
            if line.lower().startswith("server:"):
                server = line.split(":", 1)[1].strip()
        return {"ok": True, "status": stat, "server": server,
                "head": headtxt[:600],
                "body": body.decode("utf-8", "replace")[:maxbytes],
                "len": len(body)}
    except Exception as e:
        return {"ok": False, "err": "%s: %s" % (type(e).__name__, e)}
    finally:
        try: s.close()
        except Exception: pass

def title_of(body):
    m = re.search(r"<title[^>]*>(.*?)</title>", body, re.I | re.S)
    if m:
        return re.sub(r"\s+", " ", target22.example.com(1)).strip()[:120]
    return ""

def sigs_of(text):
    return [name for name, pat in SIGS if pat in text]

def one(target, mode):
    ip, _, port = target.partition(":")
    port = int(port)
    out = {"target": target}
    r = raw_http(ip, port, "/")
    if not r["ok"]:
        out["root"] = {"err": r["err"]}
        return out
    out["root"] = {"status": r["status"], "server": r["server"],
                   "title": title_of(r["body"]), "len": r["len"],
                   "sigs": sigs_of(r["body"]), "head": r["head"][:300]}
    if mode == "path":
        paths = {}
        for p in BIZ_PATHS:
            rr = raw_http(ip, port, p)
            if not rr["ok"]:
                paths[p] = {"err": rr["err"]}
            else:
                paths[p] = {"status": rr["status"], "len": rr["len"],
                            "title": title_of(rr["body"]),
                            "sigs": sigs_of(rr["body"]),
                            "snip": re.sub(r"\s+", " ", rr["body"][:180])}
        out["paths"] = paths
    return out

def main():
    targets = [t.strip() for t in sys.argv[1].split(",") if t.strip()]
    mode = sys.argv[2] if len(sys.argv) > 2 else "base"
    workers = int(sys.argv[3]) if len(sys.argv) > 3 else 8
    res = []
    with ThreadPoolExecutor(max_workers=workers) as ex:
        futs = {ex.submit(one, t, mode): t for t in targets}
        for f in as_completed(futs):
            t = futs[f]
            try:
                r = f.result()
            except Exception as e:
                r = {"target": t, "error": str(e)}
            res.append(r)
            root = r.get("root", {})
            print("[%s] st=%s srv=%s title=%r sigs=%s" % (
                t, root.get("status"), root.get("server"), root.get("title"), root.get("sigs")), flush=True)
    res.sort(key=lambda x: x["target"])
    with open(r"<HOME>\Desktop\六助新\recon\httpfp-%s.json" % mode, "w", encoding="utf-8") as fh:
        json.dump(res, fh, ensure_ascii=False, indent=1)
    print("\n=== JSON saved ===")

if __name__ == "__main__":
    main()
