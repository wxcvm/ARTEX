# -*- coding: utf-8 -*-
"""多目标业务路径探测驱动（原生 socket，不走代理）。
用法: python sweep.py "<ip:port,ip:port,...>" "<hosthdr1,hosthdr2|IP>" "<path1,path2>" [workers]
hosthdr 传 IP 表示用目标 IP 做 Host（默认）。
"""
import sys, os, json, re, socket, ssl
from concurrent.futures import ThreadPoolExecutor, as_completed

for k in ("http_proxy","https_proxy","all_proxy","HTTP_PROXY","HTTPS_PROXY","ALL_PROXY"):
    os.environ.pop(k, None)
try:
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
except Exception:
    pass

UA = "Mozilla/5.0 (Linux; Android 12; lz2025) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Mobile Safari/537.36"
SIGS = [("在线充值","在线充值"),("JSP_index","My JSP 'index.jsp' starting page"),
        ("fenleilist","LHService/fenleilist"),("LHService","LHService"),
        ("IMService","IMService"),("Tk_Kj","Tk_Kj"),("PayBmob","PayBmob"),
        ("UpLoadPic","UpLoadPic"),("六助","六助"),("lz2025","lz2025"),
        ("appLogin","appLogin"),("applogin","applogin"),("jielong","jielong"),
        ("v27","v27"),("springfox","springfox"),("Whitelabel","Whitelabel Error")]

def raw_http(ip, port, path, host, timeout=6.0, maxbytes=6000):
    https = port in (443, 8443)
    req = ("GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nAccept: */*\r\n"
           "Accept-Encoding: identity\r\nConnection: close\r\n\r\n" % (path, host, UA))
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(timeout)
    try:
        s.connect((ip, port))
        if https:
            ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
            ctx.check_hostname = False; ctx.verify_mode = ssl.CERT_NONE
            ctx.set_ciphers("DEFAULT@SECLEVEL=1")
            s = ctx.wrap_socket(s, server_hostname=host)
        s.sendall(req.encode())
        buf = b""
        while len(buf) < maxbytes:
            try: c = s.recv(4096)
            except socket.timeout: break
            if not c: break
            buf += c
        head, _, body = buf.partition(b"\r\n\r\n")
        ht = head.decode("utf-8","replace")
        st = 0
        m = re.match(r"HTTP/\d\.\d (\d+)", ht)
        if m: st = int(target22.example.com(1))
        srv = ""
        for line in ht.split("\r\n"):
            if line.lower().startswith("server:"): srv = line.split(":",1)[1].strip()
        return {"ok": True, "status": st, "server": srv, "len": len(body),
                "body": body.decode("utf-8","replace")}
    except Exception as e:
        return {"ok": False, "err": "%s: %s" % (type(e).__name__, e)}
    finally:
        try: s.close()
        except Exception: pass

def title_of(b):
    m = re.search(r"<title[^>]*>(.*?)</title>", b, re.I|re.S)
    return re.sub(r"\s+"," ",target22.example.com(1)).strip()[:90] if m else ""

def sigs(t):
    return [n for n,p in SIGS if p in t]

def main():
    targets = [t.strip() for t in sys.argv[1].split(",") if t.strip()]
    hostarg = sys.argv[2].strip() if len(sys.argv) > 2 else "IP"
    paths = [p.strip() for p in sys.argv[3].split(",") if p.strip()]
    workers = int(sys.argv[4]) if len(sys.argv) > 4 else 8

    tasks = []
    for t in targets:
        ip, _, port = t.partition(":")
        port = int(port)
        hl = [ip] if hostarg.upper() == "IP" else [h.strip() for h in hostarg.split(",") if h.strip()]
        for h in hl:
            for p in paths:
                tasks.append((ip, port, h, p))

    hits, rows = [], []
    with ThreadPoolExecutor(max_workers=workers) as ex:
        futs = {ex.submit(raw_http, ip, port, p, h): (ip, port, h, p) for ip, port, h, p in tasks}
        for f in as_completed(futs):
            ip, port, h, p = futs[f]
            r = f.result()
            if not r["ok"]:
                rows.append((ip, port, h, p, None, r["err"])); continue
            sg = sigs(r["body"])
            tl = title_of(r["body"])
            rows.append((ip, port, h, p, r, None))
            snip = re.sub(r"\s+"," ",r["body"][:130])
            if sg:
                hits.append({"ip": ip, "port": port, "host_hdr": h, "path": p,
                             "status": r["status"], "server": r["server"],
                             "len": r["len"], "title": tl, "sigs": sg, "snip": snip})
                print("[HIT ] %s:%d @%-20s %-34s st=%-4s len=%-6s %s | %s" % (
                    ip, port, h, p, r["status"], r["len"], sg, snip), flush=True)
            elif r["status"] not in (404, 0) and r["len"] > 100:
                print("[note] %s:%d @%-20s %-34s st=%-4s len=%-6s title=%r" % (
                    ip, port, h, p, r["status"], r["len"], tl), flush=True)
    with open(r"<HOME>\Desktop\六助新\recon\sweep-hits.json","w",encoding="utf-8") as fh:
        json.dump(hits, fh, ensure_ascii=False, indent=1)
    print("\n=== TASKS=%d HITS=%d ===" % (len(tasks), len(hits)))

if __name__ == "__main__":
    main()
