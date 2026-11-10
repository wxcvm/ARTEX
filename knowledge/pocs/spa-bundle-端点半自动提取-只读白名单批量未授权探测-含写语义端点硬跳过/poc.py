#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""Batch unauthorised probe of tgt2905 fund-site endpoints. STRICTLY READ-ONLY.

Policy:
  * last path segment starting with find/query/get/select/list/my/agent -> READ  (GET probed)
  * last path segment matching a WRITE verb -> SKIPPED entirely (never requested)
  * anything else -> GET probed; POST {} retried only when classified "read-like"
Rate <=2 req/s. Output runs/tgt2905_probe.json
"""
import json
import os
import re
import time

import requests

HOST = "https://example.com"
PROXY = {"http": "http://127.0.0.1:7892", "https": "http://127.0.0.1:7892"}
OUTDIR = r"C:\rt\dgha\runs"
EPS = os.path.join(OUTDIR, "tgt2905_endpoints.json")
OUT = os.path.join(OUTDIR, "tgt2905_probe.json")
RATE = 0.55
UA = ("Mozilla/5.0 (Linux; Android 12; SM-G981B) AppleWebKit/537.36 (KHTML, like Gecko) "
      "Chrome/203.0.113.20 Mobile Safari/537.36")

READ_PREFIX = re.compile(r"^(find|query|get|select|list|my|agentquery|check|is|valid|logout|login)", re.I)
WRITE_WORDS = re.compile(
    r"(add|insert|update|delete|remove|modify|edit|apply|recharge|withdraw|receive|obtain|give|"
    r"send|readall|readmessage|disabled|callback|save|create|register|bind|change|storage|"
    r"inspection|cancel|reset|upload|exec|import|export|confirm|success|docash|"
    r"pay\b|paymethod/change|setdefault|logout|signup|getpocket|receivebonus|receivereward|"
    r"receivehandsel|receivereward|receivememberhandsel|signobtain|registration|registers|givepromotion)",
    re.I)
NOISE = re.compile(r"(__|^/content-|^/user-agent|^/system-(alert|ui)$|^/payVoucher$|^/bankFour|^/bankThree|"
                   r"^/banktwo|^/bank_card|^/bankDetail|^/payee$|^/payment$|^/userInfo$|^/userName$|"
                   r"^/userInterface$|^/systemConfig$|^/member/bankCard//)", re.I)
API = re.compile(r"^/(agent|member|capital|discount|pay|im_sys|sys|user)/")


def load():
    d = json.load(open(EPS, encoding="utf-8"))
    read, skip = [], []
    for p in d["endpoints"]:
        if NOISE.search(p) or not API.match(p):
            continue
        last = p.rstrip("/").split("/")[-1]
        if READ_PREFIX.match(last):
            read.append(p)
        elif WRITE_WORDS.search(last):
            skip.append(p)
        else:
            read.append(p)
    return sorted(set(read)), sorted(set(skip))


def classify(status, body):
    if status in (404, 405, 500, 502, 403, 400):
        return "none"
    try:
        j = json.loads(body)
    except Exception:
        return "other"
    if isinstance(j, dict):
        if j.get("status") == 1:
            return "SUCCESS"
        if j.get("status") == 2:
            return "auth"
        msg = str(j.get("msg", ""))
        if "未登录" in msg:
            return "auth"
        if "参数" in msg or "Missing" in msg or "Required" in msg or "请求方式" in msg:
            return "needparam"
    return "other"


def main():
    read, skip = load()
    print("READ endpoints:", len(read), " SKIPPED(write):", len(skip))
    json.dump({"skipped_write": skip}, open(os.path.join(OUTDIR, "tgt2905_skipped.json"), "w",
                                            encoding="utf-8"), ensure_ascii=False, indent=1)
    sess = requests.Session()
    sess.proxies = PROXY
    sess.headers.update({"User-Agent": UA, "Accept": "*/*"})
    results, n, t0 = [], 0, time.time()
    for ep in read:
        rec = {"ep": ep}
        try:
            r = sess.get(HOST + ep, timeout=20)
            body = r.text[:6000]
            st, cls = r.status_code, classify(r.status_code, body)
        except Exception as e:
            st, body, cls = 0, "ERR:%s" % e, "err"
        n += 1
        time.sleep(RATE)
        rec["GET"] = {"status": st, "cls": cls, "len": len(body), "body": body}
        print("[%-9s] %-3s %-6s %s" % (cls, st, rec["GET"]["len"], ep), flush=True)
        if cls in ("needparam", "none", "other"):
            try:
                r = target25.example.com(HOST + ep, data=b"{}", timeout=20,
                              headers={"Content-Type": "application/json"})
                body2 = r.text[:6000]
                st2, cls2 = r.status_code, classify(r.status_code, body2)
            except Exception as e:
                st2, body2, cls2 = 0, "ERR:%s" % e, "err"
            n += 1
            time.sleep(RATE)
            rec["POST"] = {"status": st2, "cls": cls2, "len": len(body2), "body": body2}
            print("[%-9s] %-3s %-6s %s  (POST)" % (cls2, st2, rec["POST"]["len"], ep), flush=True)
        results.append(rec)
    el = time.time() - t0
    meta = {"host": HOST, "requests": n, "seconds": round(el, 1), "rate": round(n / max(el, .001), 2),
            "read_eps": len(read), "skipped_write": len(skip)}
    json.dump({"meta": meta, "results": results}, open(OUT, "w", encoding="utf-8"),
              ensure_ascii=False, indent=1)
    hits = [r for r in results if any(r.get(m, {}).get("cls") == "SUCCESS" for m in ("GET", "POST"))]
    print("\n=== meta:", json.dumps(meta, ensure_ascii=False))
    print("=== SUCCESS hits:", len(hits))
    for r in hits:
        m = "GET" if r.get("GET", {}).get("cls") == "SUCCESS" else "POST"
        print("%-4s %-10s %s" % (m, r[m]["len"], r["ep"]))


if __name__ == "__main__":
    main()
