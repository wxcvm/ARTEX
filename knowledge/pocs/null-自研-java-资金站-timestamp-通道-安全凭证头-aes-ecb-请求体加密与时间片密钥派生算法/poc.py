#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""tgt2905 fund-site TIMESTAMP/STATICAES channel client.

Reverse-engineered from /static/js/main.80e1a14b.js (webpack module 191):
  h = MD5(e).hex                       (password field)
  f = SHA256(e).hex  (lower)           u = SHA256(e).hex (UPPER)
  m(appId, ipayCustomerUrl, appKey, dhVideoList) =
        r = u(dhVideoList); i = appId+ipayCustomerUrl+appKey+r
        s = b64(SHA256(i)); while len(s)<50: s += b64(SHA256(s)); return s
  d(hexstr, n) = for f in range(n): pick m(...)[int(hexstr[2f:2f+2],16) % len]
  g(n)         = d( f( str( floor(nowSec / isIpayMergeChannel) ) ), n )      -> session key
  response decrypt: AES-ECB,Pkcs7, key = utf8(g(16)) for the TIMESTAMP channel
"""
import base64
import hashlib
import json
import sys
import time

import requests
from Crypto.Cipher import AES
from Crypto.Util.Padding import unpad

try:
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
except Exception:
    pass

HOST = "https://example.com"
PROXY = {"http": "http://127.0.0.1:7892", "https": "http://127.0.0.1:7892"}
UA = ("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
      "Chrome/203.0.113.20 Safari/537.36")
DEVICE_ID = "target26.example.com"

# getBaboConfig constants (unauthenticated /member/bb/api/getBaboConfig)
BABO = {
    "appId": "bf709056-6806-4252-98e6-3f9ca2022cf0",
    "ipayCustomerUrl": "https://target27.example.com",
    "appKey": "KgBPAHYAUQB5AGUAbAArAEEANQAvAG0ARgA/AHoAWQBCAFUAXAAlAD4AMQAuAEoA",
    "dhVideoList": "https://target28.example.com/CYnj/Hb7Ya",
    "isIpayMergeChannel": 300,
}


def b64(b):
    return base64.b64encode(b).decode()


def sha256_hex(s):
    return hashlib.sha256(s.encode()).hexdigest()


def sha256_hex_up(s):
    return hashlib.sha256(s.encode()).hexdigest().upper()


def md5_hex(s):
    return hashlib.md5(s.encode()).hexdigest()


def key_table():
    r = sha256_hex_up(BABO["dhVideoList"])
    i = BABO["appId"] + BABO["ipayCustomerUrl"] + BABO["appKey"] + r
    s = b64(hashlib.sha256(i.encode()).digest())
    while len(s) < 50:
        s += b64(hashlib.sha256(s.encode()).digest())
    return s


TABLE = key_table()


def pick(hexstr, n):
    p = ""
    for f in range(n):
        p += TABLE[int(hexstr[2 * f:2 * f + 2], 16) % len(TABLE)]
    return p


def session_key(n=16, offset=0):
    slot = int(time.time()) // BABO["isIpayMergeChannel"] + offset
    return pick(sha256_hex(str(slot)), n)


def aes_dec(b64data, key_utf8):
    try:
        ct = base64.b64decode(b64data)
        pt = AES.new(key_utf8.encode("utf-8"), AES.MODE_ECB).decrypt(ct)
        return unpad(pt, 16).decode("utf-8")
    except Exception as e:
        return None


def aes_enc(plain, key_utf8):
    data = plain.encode("utf-8")
    pad = 16 - len(data) % 16
    data += bytes([pad]) * pad
    return base64.b64encode(AES.new(key_utf8.encode("utf-8"), AES.MODE_ECB).encrypt(data)).decode()


def call(ep, payload=None, channel="TIMESTAMP", verbose=True, extra_headers=None):
    """POST to a TIMESTAMP-channel endpoint and return the decrypted JSON."""
    s = requests.Session()
    s.proxies = PROXY
    ts = str(int(time.time() * 1000))
    hdr = {"User-Agent": UA, "Accept": "application/json, text/plain, */*",
           "Content-Type": "application/json", "Origin": HOST, "Referer": HOST + "/",
           "Isauthent": ts, "isstep": ts}
    if extra_headers:
        hdr.update(extra_headers)
    body = json.dumps(payload if payload is not None else {})
    r = target8.example.com(HOST + ep, data=body.encode(), headers=hdr, timeout=20)
    txt = r.text
    if "has been blocked" in txt:
        return {"_waf": True, "_http": r.status_code, "_raw": txt[:200]}
    try:
        j = r.json()
    except Exception:
        return {"_raw": txt[:400], "_http": r.status_code}
    out = {"_http": r.status_code, "_raw_keys": list(j.keys())}
    if "status" in j:
        out.update(j)
        return out
    # needs decryption
    data = j.get("data")
    tried = {}
    cands = [("session", session_key(16, off)) for off in (-1, 0, 1)]
    if j.get("AA"):
        for ln in (16, 32, 24):
            cands.append(("AA-lower-%d" % ln, pick(str(j["AA"]).lower(), ln)))
            cands.append(("AA-upper-%d" % ln, pick(str(j["AA"]), ln)))
    for name, k in cands:
        pt = aes_dec(data, k)
        if pt:
            tried[name] = pt[:200]
            try:
                out.update(json.loads(pt))
                out["_key"] = name
                return out
            except Exception:
                out["_plaintext"] = pt[:600]
                out["_key"] = name
                return out
    out["_decrypt_failed"] = True
    out["_cands"] = {k: v for k, v in list(tried.items())}
    out["_AA"] = j.get("AA")
    out["_data"] = data
    return out


def enc_post(ep, payload, offset=0, verbose=True, extra_headers=None):
    """Proper TIMESTAMP-channel request: body = {data: AES-ECB(json, g(16)), AA: sha256(json).hex}"""
    s = requests.Session()
    s.proxies = PROXY
    ts = str(int(time.time() * 1000))
    key = session_key(16, offset)
    plain = json.dumps(payload, separators=(",", ":"))
    body = json.dumps({"data": aes_enc(plain, key), "AA": sha256_hex(plain)})
    hdr = {"User-Agent": UA, "Accept": "application/json, text/plain, */*",
           "Content-Type": "application/json", "Isauthent": ts, "isstep": ts}
    if extra_headers:
        hdr.update(extra_headers)
    r = target8.example.com(HOST + ep, data=body.encode(), headers=hdr, timeout=20)
    txt = r.text
    res = {"_http": r.status_code, "_key": key, "_plain": plain}
    if "has been blocked" in txt:
        res["_waf"] = True
        return res
    try:
        j = r.json()
    except Exception:
        res["_raw"] = txt[:400]
        return res
    if "status" in j:
        res.update(j)
        return res
    data = j.get("data")
    out = dict(res)
    out["_raw_keys"] = list(j.keys())
    out["_AA"] = j.get("AA")
    for nm, k in [("session", key), ("session+1", session_key(16, 1)),
                  ("session-1", session_key(16, -1)),
                  ("AA-lower16", pick(str(j.get("AA", "")).lower(), 16)),
                  ("AA-upper16", pick(str(j.get("AA", "")), 16))]:
        pt = aes_dec(data, k)
        if pt:
            try:
                out.update(json.loads(pt))
            except Exception:
                out["_plaintext"] = pt[:600]
            out["_deckey"] = nm
            return out
    out["_decrypt_failed"] = True
    out["_data"] = data
    return out


if __name__ == "__main__":
    print("key table len:", len(TABLE))
    print("session key now:", session_key(16))
    ep = sys.argv[1] if len(sys.argv) > 1 else "/member/memberManager/loginV3"
    payload = json.loads(sys.argv[2]) if len(sys.argv) > 2 else {}
    print("EP:", ep, "payload:", payload)
    res = call(ep, payload)
    print(json.dumps(res, ensure_ascii=False, indent=1)[:3000])
