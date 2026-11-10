# -*- coding: utf-8 -*-
"""端口 banner 抓取（原生 socket，不走代理）。用法: python banner.py "ip:port,..." [timeout]"""
import socket, os, sys, json, re
from concurrent.futures import ThreadPoolExecutor, as_completed

for k in ("http_proxy","https_proxy","all_proxy","HTTP_PROXY","HTTPS_PROXY","ALL_PROXY"):
    os.environ.pop(k, None)
try:
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
except Exception:
    pass

def grab(ip, port, timeout=6.0):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(timeout)
    try:
        s.connect((ip, port))
        # 先把 banner 读出来（SSH/MySQL 是服务端先发）
        first = b""
        try:
            s.settimeout(3.0)
            first = s.recv(1024)
        except socket.timeout:
            pass
        out = {}
        if first:
            out["banner"] = first[:300].decode("utf-8", "replace")
            t = first.strip()
            if t.startswith(b"SSH-"):
                out["proto"] = "ssh"
                out["product"] = t.decode("utf-8", "replace")
            elif len(first) > 5 and first[4] == 10:
                # MySQL 握手包
                out["proto"] = "mysql"
                try:
                    nlen = first[0]
                    ver_end = first.index(b"\x00", 5)
                    out["product"] = "MySQL " + first[5:ver_end].decode("utf-8", "replace")
                    # 截断的 auth plugin / capability
                except Exception:
                    out["product"] = "MySQL (handshake)"
        # 再做一次 HTTP 尝试（若还没识别）
        if out.get("proto") not in ("ssh", "mysql"):
            try:
                s2 = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
                s2.settimeout(timeout)
                s2.connect((ip, port))
                req = ("GET / HTTP/1.1\r\nHost: %s\r\nUser-Agent: Mozilla/5.0\r\n"
                       "Accept: */*\r\nConnection: close\r\n\r\n" % ip)
                s2.sendall(req.encode())
                buf = b""
                while len(buf) < 3000:
                    try:
                        c = s2.recv(2048)
                    except socket.timeout:
                        break
                    if not c:
                        break
                    buf += c
                s2.close()
                ht = buf.decode("utf-8", "replace")
                m = re.match(r"HTTP/\d\.\d (\d+)", ht)
                if m:
                    srv = ""
                    for line in ht.split("\r\n"):
                        if line.lower().startswith("server:"):
                            srv = line.split(":",1)[1].strip()
                    out["proto"] = "http"
                    out["http"] = {"status": int(target22.example.com(1)), "server": srv,
                                   "head": ht[:250].replace("\r\n", " | ")}
                elif buf:
                    out.setdefault("banner", buf[:200].decode("utf-8","replace"))
            except Exception as e:
                out.setdefault("http_err", type(e).__name__)
        return ip, port, out
    except Exception as e:
        return ip, port, {"err": "%s: %s" % (type(e).__name__, e)}
    finally:
        try: s.close()
        except Exception: pass

def main():
    targets = [t.strip() for t in sys.argv[1].split(",") if t.strip()]
    timeout = float(sys.argv[2]) if len(sys.argv) > 2 else 6.0
    res = {}
    jobs = []
    for t in targets:
        h, _, p = t.partition(":")
        jobs.append((h, int(p)))
    with ThreadPoolExecutor(max_workers=8) as ex:
        futs = [ex.submit(grab, h, p, timeout) for h, p in jobs]
        for f in as_completed(futs):
            ip, port, out = f.result()
            res["%s:%s" % (ip, port)] = out
            print("%-24s %s" % ("%s:%d" % (ip, port), json.dumps(out, ensure_ascii=False)), flush=True)
    with open(r"<HOME>\Desktop\六助新\recon\banners.json","w",encoding="utf-8") as fh:
        json.dump(res, fh, ensure_ascii=False, indent=1)

if __name__ == "__main__":
    main()
