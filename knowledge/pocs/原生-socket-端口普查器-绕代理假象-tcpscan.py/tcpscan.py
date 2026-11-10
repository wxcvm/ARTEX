# -*- coding: utf-8 -*-
"""原生 socket 端口扫描器 —— 绝不走代理环境变量。
用法: python tcpscan.py <cidr_prefix|ip,ip> <ports> [workers] [timeout]
"""
import socket, sys, json, time, ipaddress
from concurrent.futures import ThreadPoolExecutor, as_completed

# 硬清代理环境变量（socket 本身不用，但避免任何库偷读）
import os
for k in ("http_proxy","https_proxy","all_proxy","HTTP_PROXY","HTTPS_PROXY","ALL_PROXY","no_proxy","NO_PROXY"):
    os.environ.pop(k, None)

def probe(ip, port, timeout=2.5):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.settimeout(timeout)
    try:
        s.connect((ip, port))
        try:
            s.settimeout(1.2)
            s.sendall(b"")
        except Exception:
            pass
        return (ip, port, True, None)
    except socket.timeout:
        return (ip, port, False, "timeout")
    except ConnectionRefusedError:
        return (ip, port, False, "refused")
    except OSError as e:
        return (ip, port, False, type(e).__name__)
    finally:
        try: s.close()
        except Exception: pass

def main():
    target = sys.argv[1]
    ports = [int(p) for p in sys.argv[2].split(",")]
    workers = int(sys.argv[3]) if len(sys.argv) > 3 else 12
    timeout = float(sys.argv[4]) if len(sys.argv) > 4 else 2.5

    if "/" in target:
        net = ipaddress.ip_network(target, strict=False)
        ips = [str(x) for x in net.hosts()]
    else:
        ips = [x.strip() for x in target.split(",") if x.strip()]

    tasks = [(ip, p) for ip in ips for p in ports]
    hits, misses = [], []
    t0 = time.time()
    with ThreadPoolExecutor(max_workers=workers) as ex:
        futs = {ex.submit(probe, ip, p, timeout): (ip, p) for ip, p in tasks}
        for f in as_completed(futs):
            ip, p, ok, err = f.result()
            if ok:
                hits.append({"ip": ip, "port": p})
                print("[OPEN] %s:%d" % (ip, p), flush=True)
            else:
                misses.append({"ip": ip, "port": p, "err": err})
    hits.sort(key=lambda x: (x["ip"], x["port"]))
    print("\n=== RESULT ===")
    print(json.dumps({"target": target, "ports": ports, "open": hits,
                      "open_count": len(hits), "probed": len(tasks),
                      "elapsed": round(time.time()-t0, 1)}, ensure_ascii=False, indent=1))

if __name__ == "__main__":
    main()
