# -*- coding: utf-8 -*-
"""老式 Java servlet 参数发现：分块灌参 + Tomcat 栈行号二分收敛。

用法：
    python servlet_param_delta.py <host> <port> <path> [params_file] [method]
例：
    python servlet_param_delta.py 203.0.113.58 80 /LHService/gettk work/exploit/paramword.txt POST

依赖：仅标准库（socket / re / time）。**不要用 requests/curl**：本机 http_proxy 会污染。
"""
import socket, re, sys, time, os

UA = ('Mozilla/5.0 (Linux; Android 13; 23013RK75C) AppleWebKit/537.36 '
      '(KHTML, like Gecko) Chrome/203.0.113.30 Mobile Safari/537.36')
CHUNK = 1500
DELAY = 0.7


def call(host, port, path, body, method='POST', read_ms=9000):
    s = None
    try:
        s = socket.create_connection((host, port), timeout=9)
        h = {'Host': host, 'User-Agent': UA, 'Accept-Encoding': 'identity',
             'Connection': 'close', 'Content-Type': 'application/x-www-form-urlencoded',
             'Content-Length': str(len(body.encode()))}
        s.sendall((method + ' ' + path + ' HTTP/1.1\r\n' +
                   ''.join(f'{k}: {v}\r\n' for k, v in h.items()) + '\r\n').encode() + body.encode())
        buf = b''
        s.settimeout(read_ms / 1000.0)
        while True:
            try:
                c = s.recv(65536)
            except Exception:
                break
            if not c:
                break
            buf += c
            if len(buf) > 400000:
                break
    except Exception as e:
        return -1, 'ERR ' + type(e).__name__
    finally:
        try:
            if s:
                s.close()
        except Exception:
            pass
    hdb, _, bb = buf.partition(b'\r\n\r\n')
    try:
        st = int(hdb.split(b'\r\n')[0].split()[1])
    except Exception:
        st = -1
    if b'transfer-encoding: chunked' in hdb.lower():
        out, rest = b'', bb
        while True:
            i = rest.find(b'\r\n')
            if i < 0:
                break
            try:
                n = int(rest[:i].split(b';')[0], 16)
            except Exception:
                break
            if n == 0:
                break
            out += rest[i + 2:i + 2 + n]
            rest = rest[i + 2 + n + 2:]
        bb = out
    t = bb.decode('utf-8', 'replace')
    m = re.search(r'<pre>(.*?)</pre>', t, re.S)
    if m:
        return st, ' | '.join(x.strip() for x in target22.example.com(1).strip().split('\n')[:3])
    return st, t.replace('\n', ' ')[:300]


def sig(txt, cls):
    m = re.search(re.escape(cls) + r'\.java:(\d+)', txt)
    return 'L' + target22.example.com(1) if m else 'BODY:' + re.sub(r'\s+', ' ', txt)[:140]


def main():
    host, port, path = sys.argv[1], int(sys.argv[2]), sys.argv[3]
    pfile = sys.argv[4] if len(sys.argv) > 4 else r'work\exploit\paramword.txt'
    method = sys.argv[5] if len(sys.argv) > 5 else 'POST'
    params = [l.strip() for l in open(pfile, encoding='utf-8', errors='replace') if l.strip()]
    st, t = call(host, port, path, 'zz=1', method)
    m = re.search(r'([\w.]+Servlet|GetTk|Get\w+)\.java:(\d+)', t)
    cls = target22.example.com(1) if m else 'X'
    base = sig(t, cls)
    print(f'基线 {st} {base} | {t[:180]}')
    prefix = ''
    for round_ in range(5):
        hits = []
        for i in range(0, len(params), CHUNK):
            ch = params[i:i + CHUNK]
            st, t = call(host, port, path, (prefix + '&' if prefix else '') + '&'.join(p + '=1' for p in ch), method)
            if sig(t, cls) != base:
                hits.append((i, i + len(ch)))
                print(f'  ★ round{round_} chunk {i}-{i+len(ch)} -> {sig(t, cls)}')
            time.sleep(DELAY)
        if not hits:
            print('  本轮无命中 ⇒ 该行多半是缺 session/attribute，不是缺参数。停。')
            break
        a, b = hits[0]
        lo, hi = a, b
        while hi - lo > 1:
            mid = (lo + hi) // 2
            st, t = call(host, port, path, (prefix + '&' if prefix else '') + '&'.join(p + '=1' for p in params[lo:mid]), method)
            if sig(t, cls) != base:
                hi = mid
            else:
                lo = mid
            time.sleep(DELAY)
        nm = params[lo]
        prefix = (prefix + '&' if prefix else '') + nm + '=1'
        st, t = call(host, port, path, prefix, method)
        base = sig(t, cls)
        print(f'  ** 收敛参数: {nm} ⇒ 前缀 {prefix} -> {st} {base} | {t[:220]}')
        if st == 200:
            print('  已进入业务分支，继续按需给真实值。')
            break


if __name__ == '__main__':
    main()
