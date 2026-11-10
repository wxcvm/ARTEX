# -*- coding: utf-8 -*-
"""LHZS 节点签名探测器：CIDR 内 TCP 存活 → 打免签名接口，按响应体判定是否真节点。
用法: python lhprobe.py <cidr> <ports> [workers] [hosthdr]
判据（必须看正文，不看状态码）：
  HIT   : body 含 "routing-v27" 或 "LHService" 或 "在线充值" 或 "My JSP" 等 LHZS 特征
  GAT   : 400 INVALID_INPUT / 403 REQUEST_REJECTED  = LHZS 网关拒绝码（也是强归属证据）
  --    : 其它（无关站点/默认页）
"""
import socket, ssl, os, sys, json, re, ipaddress, time
from concurrent.futures import ThreadPoolExecutor, as_completed

for k in ("http_proxy","https_proxy","all_proxy","HTTP_PROXY","HTTPS_PROXY","ALL_PROXY"):
    os.environ.pop(k, None)
try:
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
except Exception:
    pass

UA = "Mozilla/5.0 (Linux; Android 12; lz2025) AppleWebKit/537.36"
PATHS = ["/LHService/v27/routing/health", "/LHService/"]
LH_SIGS = ["routing-v27", "在线充值", "My JSP", "六合", "fenleilist",
           "LHService", "Tk_Kj", "IMService", "PayBmob", "UpLoadPic"]
GATE_SIGS = ["INVALID_INPUT", "REQUEST_REJECTED", "request_rejected"]

def tcp_ok(ip, port, timeout=2.0):
    s = socket.socket(); s.settimeout(timeout)
    try:
        s.connect((ip, port)); return True
    except Exception:
        return False
    finally:
        try: s.close()
        except Exception: pass

def http(ip, port, path, host, timeout=6.0, maxb=4000):
    https = port in (443, 8443)
    req = ("GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nAccept: */*\r\n"
           "X-Security-Version: 2\r\nConnection: close\r\n\r\n" % (path, host, UA))
    s = socket.socket(); s.settimeout(timeout)
    try:
        s.connect((ip, port))
        if https:
            ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
            ctx.check_hostname = False; ctx.verify_mode = ssl.CERT_NONE
            try: ctx.set_ciphers("DEFAULT@SECLEVEL=1")
            except Exception: pass
            s = ctx.wrap_socket(s, server_hostname=host)
        s.sendall(req.encode())
        buf = b""
        while len(buf) < maxb:
            try: c = s.recv(3072)
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
        return {"ok": True, "status": st, "server": srv,
                "body": body.decode("utf-8","replace")[:3000], "len": len(body)}
    except Exception as e:
        return {"ok": False, "err": type(e).__name__}
    finally:
        try: s.close()
        except Exception: pass

def check(ip, port, hostarg):
    if not tcp_ok(ip, port):
        return None
    host = ip if hostarg.upper() == "IP" else hostarg
    out = {"ip": ip, "port": port, "host_hdr": host, "probes": {}}
    for p in PATHS:
        r = http(ip, port, p, host)
        if not r["ok"]:
            out["probes"][p] = r
            continue
        bl = r["body"]; hl = r["body"].lower()
        lh = [s for s in LH_SIGS if s.lower() in hl]
        gt = [s for s in GATE_SIGS if s.lower() in hl]
        out["probes"][p] = {"status": r["status"], "server": r["server"],
                            "len": r["len"], "lh_sigs": lh, "gate_sigs": gt,
                            "snip": re.sub(r"\s+"," ",bl[:160])}
    return out

def main():
    cidr = sys.argv[1]
    ports = [int(x) for x in sys.argv[2].split(",")]
    workers = int(sys.argv[3]) if len(sys.argv) > 3 else 10
    hostarg = sys.argv[4] if len(sys.argv) > 4 else "IP"
    ips = [str(x) for x in ipaddress.ip_network(cidr, strict=False).hosts()]
    tasks = [(ip, p) for ip in ips for p in ports]
    hits, gates, t0 = [], [], time.time()
    allrows = []
    with ThreadPoolExecutor(max_workers=workers) as ex:
        futs = {ex.submit(check, ip, p, hostarg): (ip, p) for ip, p in tasks}
        for f in as_completed(futs):
            ip, p = futs[f]
            try: r = f.result()
            except Exception: r = None
            if not r: continue
            allrows.append(r)
            anylh = any(pr.get("lh_sigs") for pr in r["probes"].values())
            anygate = any(pr.get("gate_sigs") for pr in r["probes"].values())
            if anylh:
                hits.append(r)
                print("[HIT ] %s:%d %s" % (ip, p, json.dumps(r["probes"], ensure_ascii=False)), flush=True)
            elif anygate:
                gates.append(r)
                print("[GATE] %s:%d %s" % (ip, p, json.dumps(r["probes"], ensure_ascii=False)), flush=True)
    with open(r"<HOME>\Desktop\六助新\recon\lhprobe-%s.json" % cidr.replace("/", "_").replace(".", "_"),
              "w", encoding="utf-8") as fh:
        json.dump({"hits": hits, "gates": gates, "all": allrows}, fh, ensure_ascii=False, indent=1)
    print("\n=== %s : HIT=%d GATE=%d ALL=%d  用时 %.0fs ===" % (cidr, len(hits), len(gates), len(allrows), time.time()-t0))

if __name__ == "__main__":
    main()
