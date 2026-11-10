#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""FastAPI 应用未授权面快速测绘：拉 OpenAPI 规范 → 枚举全端点 → 分类鉴权状态 → 挑出无 SEC 的端点实测

用法: python fastapi_openapi_enum.py <host> <port> [base_path]
"""
import sys, json, socket

def http(host, port, path, method='GET', body=None, timeout=10):
    h = '%s %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\nUser-Agent: Mozilla/5.0\r\nAccept: application/json\r\n' % (method, path, host)
    if body is not None:
        h += 'Content-Type: application/json\r\nContent-Length: %d\r\n' % len(body)
    h += '\r\n'
    s = socket.create_connection((host, port), timeout=timeout); s.settimeout(timeout)
    s.sendall(h.encode() + (body.encode() if body else b''))
    d = b''
    while True:
        try: x = s.recv(65536)
        except: break
        if not x: break
        d += x
    s.close(); return d.decode('utf-8', 'ignore')

def main(host, port, base=''):
    raw = http(host, port, base + '/openapi.json')
    i = raw.find('{')
    if i < 0:
        print('[-] 未暴露 /openapi.json'); return
    spec = json.loads(raw[i:])
    print('[+] %s v%s — 端点 %d 个' % (spec.get('info', {}).get('title'), spec.get('info', {}).get('version'), len(spec.get('paths', {}))))
    schemas = spec.get('components', {}).get('schemas', {})
    noauth, auth = [], []
    for p in sorted(spec.get('paths', {})):
        for m, o in spec['paths'][p].items():
            if not isinstance(o, dict): continue
            (auth if o.get('security') else noauth).append((m.upper(), p, o.get('summary', '')))
    print('\n=== 无需鉴权(%d) —— 优先实测 ===' % len(noauth))
    for m, p, s in noauth:
        r = http(host, port, base + p, method=m if m in ('GET', 'POST') else 'GET',
                 body='{}' if m == 'POST' else None)
        code = r.split('\r\n')[0]
        body = r.split('\r\n\r\n', 1)[1][:160].replace('\n', ' ') if '\r\n\r\n' in r else ''
        print('  %-6s %-46s %-24s %s' % (m, p, code[:24], body))
    print('\n=== 需鉴权(%d) —— 已泄露路由地图 ===' % len(auth))
    for m, p, s in auth:
        print('  %-6s %-52s %s' % (m, p, s[:44]))
    # 提取 login / 内部端点 + request schema，用于后续针对性爆破
    print('\n=== 高价值目标 ===')
    for m, p, s in noauth + auth:
        if any(k in p.lower() for k in ('login', 'internal', 'admin', 'token', 'sign', 'callback', 'test_send')):
            print('  %-6s %s' % (m, p))
            op = spec['paths'][p].get(m.lower(), {})
            rb = op.get('requestBody', {})
            ref = str(rb)
            for name, sch in schemas.items():
                if name in ref:
                    print('         schema %s: %s' % (name, json.dumps(sch.get('properties', {}), ensure_ascii=False)[:200]))
                    print('         required: %s' % sch.get('required'))

if __name__ == '__main__':
    main(sys.argv[1], int(sys.argv[2]), sys.argv[3] if len(sys.argv) > 3 else '')