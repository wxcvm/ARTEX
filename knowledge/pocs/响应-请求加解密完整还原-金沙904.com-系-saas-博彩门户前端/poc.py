#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
请求加密（已字节级验证）:
  AES-CBC + ZeroPadding, iv="5421698523412578", key=EN(token)=md5(token+md5(token))[:16]
  body = {"encryptString": "<base64>"}

验证证据（浏览器 CryptoJS 与 Python 输出完全一致）:
  enc('AAAAAAAAAAAAAAAA', SK) = VHnu0ugizjJ5wMozdaWfKg==
  enc('{"a":1}',          SK) = GzScdsBhC4mmb+auz55Rxg==
"""
import base64, hashlib, json
from Crypto.Cipher import AES

QL = "5421698523412578"                      # IV，前端的 this.Ql
SK = "<16字符, 见 verified_note>"             # mm，用于本地 round-trip
RK = "fc361cdb770aebc2126cc0dac989c896"       # gm，32 字符


def md5h(s):
    if isinstance(s, str):
        s = s.encode()
    return hashlib.md5(s).hexdigest()


def EN(t):          # k1: token === n(systemToken) 时使用
    return md5h(t + md5h(t))[:16]


def UN(t, u):       # k2: token !== n 时使用
    return md5h(t + u)[:16]


def zp(b, bs=16):   # ZeroPadding（注意：不是 Pkcs7）
    r = len(b) % bs
    return b if r == 0 else b + b"\x00" * (bs - r)


def enc(d, k, iv=QL):
    raw = d if isinstance(d, bytes) else json.dumps(d, separators=(",", ":")).encode()
    return base64.b64encode(
        AES.new(k.encode(), AES.MODE_CBC, iv.encode()).encrypt(zp(raw))).decode()


def dec(ct, k, iv=QL):
    pt = AES.new(k.encode(), AES.MODE_CBC, iv.encode()).decrypt(
        base64.b64decode(ct + "=" * ((4 - len(ct) % 4) % 4)))
    return pt.rstrip(b"\x00")


# ---------------------------------------------------------------
# obfuscator.io 字符串表还原（bundle 里 Ht() 的还原方法）
#   Sa(e,t): index = e - 259
#   Ht(e,t) = Sa(e + 877, t)
#   a() = base64.b64decode 等价
#   c() = RC4，key 为**混淆器内建常量**（非 Ht 的第二参数）
#   $d 表需先按旋转 IIFE 旋转；旋转量用自校验表达式求解
# ---------------------------------------------------------------
def deobf_note():
    return """
在 bundle 中定位:
  function Sa(e,t){e=e-(7951+6*-1187+-570);...}   -> offset 259
  function Ht(e,t){return Sa(e- -877,t)}          -> +877
  function $d(){const e=[...];return $d=function(){return e},$d()}
  旋转 IIFE: (function(e,t){...})($d,-22*-45867+-62154+-187833)
本机 Vite bundle 可直接 import 取真实常量，无需静态还原:
  const mod = await import('/assets/theme-0/commonChunk.CMoJwNYm.js');
  mod.fP   // crypto 单例 (xN 实例): getSK getRK md5 sha256 hmacSha256 encrypt decrypt gk
  mod.a_   // CryptoJS 命名空间 dt: AES MD5 SHA256 enc mode{ECB,CBC} pad{Pkcs7,ZeroPadding}
"""


if __name__ == "__main__":
    # 自检：与浏览器输出比对
    assert enc("AAAAAAAAAAAAAAAA", "tgJt2ol!^E)E8(%P") == "VHnu0ugizjJ5wMozdaWfKg=="
    assert enc('{"a":1}', "tgJt2ol!^E)E8(%P") == "GzScdsBhC4mmb+auz55Rxg=="
    print("[+] 请求加密实现与浏览器字节级一致")
    print(deobf_note())
