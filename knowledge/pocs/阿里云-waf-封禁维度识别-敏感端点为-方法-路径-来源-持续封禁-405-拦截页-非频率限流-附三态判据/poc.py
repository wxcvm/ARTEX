#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""Diagnose whether the Aliyun WAF block is content-based or IP/reputation-based."""
import json
import sys
import time

import requests

try:
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
except Exception:
    pass

sys.path.insert(0, r"C:\rt\dgha")
import tgt2905_api as K

EP = "/member/memberManager/loginV3"
READ = "/member/kefuconfig/findProblemList"


def fire(tag, ep, body, hdr):
    s = requests.Session()
    s.proxies = K.PROXY
    ts = str(int(time.time() * 1000))
    h = {"User-Agent": K.UA, "Accept": "application/json, text/plain, */*", "Content-Type": "application/json"}
    h.update(hdr)
    try:
        r = target8.example.com(target29.example.com + ep, data=body.encode() if body else None, headers=h, timeout=20)
        t = r.text
        waf = "WAF" if "has been blocked" in t else "ok "
        print("%-38s http=%s %s len=%-6s %s" % (tag, r.status_code, waf, len(t), t[:120].replace("\n", " ")))
    except Exception as e:
        print("%-38s ERR %s" % (tag, e))


# 1. plain GET read (control for IP reputation)
s = requests.Session(); s.proxies = K.PROXY
s.headers.update({"User-Agent": K.UA})
try:
    r = s.get(target29.example.com + READ, timeout=20)
    print("%-38s http=%s %s len=%s" % ("GET read control", r.status_code,
                                       "WAF" if "has been blocked" in r.text else "ok", len(r.text)))
except Exception as e:
    print("GET control ERR", e)
time.sleep(6)

ts = str(int(time.time() * 1000))
fire("POST {} + Isauthent/isstep", EP, "{}", {"Isauthent": ts, "isstep": ts})
time.sleep(6)
fire("POST {a:b} + Isauthent/isstep", EP, '{"a":"b"}', {"Isauthent": ts, "isstep": ts})
time.sleep(6)
fire("POST {data:short} + Isauthent/isstep", EP, '{"data":"AAAA","AA":"BBBB"}',
     {"Isauthent": ts, "isstep": ts})
time.sleep(6)
fire("POST 64b64+64hex + Isauthent/isstep", EP,
     json.dumps({"data": "A" * 64, "AA": "B" * 64}), {"Isauthent": ts, "isstep": ts})
