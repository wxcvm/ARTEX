# -*- coding: utf-8 -*-
"""TLS 证书抓取：CN + SAN 列表（原生 socket，不走代理）。用法: python certgrab.py "ip:port,..." """
import socket, ssl, os, sys, json, re
from concurrent.futures import ThreadPoolExecutor, as_completed

for k in ("http_proxy","https_proxy","all_proxy","HTTP_PROXY","HTTPS_PROXY","ALL_PROXY"):
    os.environ.pop(k, None)
try:
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
except Exception:
    pass

def get_cert(ip, port, timeout=8.0, sni=None):
    ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    try:
        ctx.set_ciphers("DEFAULT@SECLEVEL=1")
    except Exception:
        pass
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(timeout)
    try:
        s.connect((ip, port))
        ss = ctx.wrap_socket(s, server_hostname=sni or ip)
        der = ss.getpeercert(binary_form=True)
        cert = None
        try:
            cert = ss.getpeercert()
        except Exception:
            cert = None
        ver = ss.version()
        cipher = ss.cipher()
        ss.close()
        out = {"ok": True, "tls": ver, "cipher": (cipher or [None])[0]}
        if cert:
            subj = {k: v for t in cert.get("subject", ()) for k, v in t}
            iss = {k: v for t in cert.get("issuer", ()) for k, v in t}
            out["subject"] = subj
            out["issuer"] = iss
            out["notAfter"] = cert.get("notAfter")
            san = []
            for t in cert.get("subjectAltName", ()):
                san.append("%s:%s" % t)
            out["SAN"] = san
        else:
            out["note"] = "getpeercert() empty (verify_mode=CERT_NONE); parsing DER strings"
            out["der_len"] = len(der or b"")
            # 从 DER 里提可打印串（CN/SAN/issuer 都是明文 ASCII/UTF8）
            strs = re.findall(rb"[!-~]{4,120}", der or b"")
            names = []
            for s in strs:
                t = s.decode("latin-1")
                if "." in t and " " not in t and not t.startswith("http"):
                    names.append(t)
                elif t in ("CN", "O", "OU", "C", "ST", "L"):
                    continue
            # 去重保序
            seen, uniq = set(), []
            for n in names:
                if n not in seen:
                    seen.add(n); uniq.append(n)
            out["DER_names"] = uniq[:40]
        return out
    except Exception as e:
        return {"ok": False, "err": "%s: %s" % (type(e).__name__, e)}
    finally:
        try: s.close()
        except Exception: pass

def main():
    targets = [t.strip() for t in sys.argv[1].split(",") if t.strip()]
    res = {}
    jobs = []
    for t in targets:
        h, _, p = t.partition(":")
        jobs.append((h, int(p)))
    with ThreadPoolExecutor(max_workers=6) as ex:
        futs = {ex.submit(get_cert, h, p): (h, p) for h, p in jobs}
        for f in as_completed(futs):
            h, p = futs[f]
            out = f.result()
            res["%s:%s" % (h, p)] = out
            if out.get("ok"):
                cn = out.get("subject", {}).get("commonName", "") or "(n/a)"
                names = out.get("SAN") or out.get("DER_names") or "(none)"
                print("%-22s TLS=%-10s CN=%-40s NAMES=%s" % (
                    "%s:%d" % (h, p), out.get("tls"), cn, names), flush=True)
            else:
                print("%-22s ERR %s" % ("%s:%d" % (h, p), out.get("err")), flush=True)
    with open(r"<HOME>\Desktop\六助新\recon\certs.json","w",encoding="utf-8") as fh:
        json.dump(res, fh, ensure_ascii=False, indent=1)

if __name__ == "__main__":
    main()
