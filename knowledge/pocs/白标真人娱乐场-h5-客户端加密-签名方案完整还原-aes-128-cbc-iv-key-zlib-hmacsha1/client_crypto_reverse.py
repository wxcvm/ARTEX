#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
自研真人娱乐场 H5 客户端加密/签名方案完整还原（从 example.com 内联 bundle 逆向）

== 原始证据（example.com 首页内联 JS） ==

密钥表（ENV: dev/test/test2/training/prerelease/release）：
  keyVersionMap = {"dev":"1.0.4","test":"1.0.4","test2":"1.0.4","training":"1.0.4",
                   "prerelease":"1.2.2","release":"1.2.2"}
  keyMap        = {"dev":"AA4194657AD89A56","test":"AA4194657AD89A56",
                   "test2":"AA4194657AD89A56","training":"AA4194657AD89A56",
                   "prerelease":"015CCB80A680E129","release":"015CCB80A680E129"}
  keyIdsMap     = {"dev":"tesbmb8fg38o7g5s","test":"tesbmb8fg38o7g5s",
                   "test2":"tesbmb8fg38o7g5s","training":"tesbmb8fg38o7g5s",
                   "prerelease":"preaacciqkhtg26b","release":"probinpjms7rfm26"}

环境域名表：
  envMap = {'0':'target173.example.com','1':'target173.example.com','2':'target173.example.com',
            '3':'target174.example.com','4':'target175.example.com','5':'target173.example.com'}

加密实现：
  function S(t,e){ return b.AES.encrypt(t,e,{iv:e,mode:b.mode.CBC,padding:b.pad.Pkcs7}).toString() }
  function E(t,e){ var n=b.AES.decrypt(t,e,{iv:e,mode:b.mode.CBC,padding:b.pad.Pkcs7}); return C(n) }
  x.encrypt = function(t,e){ "string"!=typeof t && (t=JSON.stringify(t));
                             var n=b.enc.Utf8.parse(e); return S(z(t), n) }   // z = pako.deflate
  x.decrypt = function(t,e){ return B(E(t, b.enc.Utf8.parse(e))) }            // B = inflate + Utf8
  x.createSign = function(t,e,n,i){ var r=t+e+n, a=y.HmacSHA1(r,i);
                                    return b.enc.Base64.stringify(a) }

请求构造：
  objData  = DataHandle.encrypt(JSON.stringify(vData), key)
  obj1     = {encrypted:true, gzipped:true, platform:"h5", version:keyVersion,
              application:'game_http', timestamp, nonce:0,
              sign:DataHandle.createSign(objData,0,timestamp,key), keyid:keyId}
  tokenStr = DataHandle.encrypt(JSON.stringify(obj1), key)
  头部：Content-Type: application/json;charset=UTF-8 / deviceType / model / deviceId /
        X-Request-Token / keyid / token
  网关：https://gateway.<envDomain>/game-http/player/getPlayerSetting?deviceType=6&playerId=<id>

== 结论 ==
· 算法已 100% 还原并复现成功（encrypt/decrypt 往返验证通过）
· 但 <X-Request-Token> 是**服务端签发**的会话 token，客户端 key 只用于加解密与请求体签名
  → 拿不到真实 token 时，伪造 token 一律返回 {"code":20000,"msg":"无效token"}
· 该网关**拒绝 IP 直连**（TLSV1_ALERT_INTERNAL_ERROR），必须域名直连
"""
import json, base64, zlib, hmac, hashlib, time
from Crypto.Cipher import AES
from Crypto.Util.Padding import pad, unpad

KEY_MAP = {
    "dev": "AA4194657AD89A56", "test": "AA4194657AD89A56", "test2": "AA4194657AD89A56",
    "training": "AA4194657AD89A56", "prerelease": "015CCB80A680E129", "release": "015CCB80A680E129",
}
KEY_ID_MAP = {
    "dev": "tesbmb8fg38o7g5s", "test": "tesbmb8fg38o7g5s", "test2": "tesbmb8fg38o7g5s",
    "training": "tesbmb8fg38o7g5s", "prerelease": "preaacciqkhtg26b", "release": "probinpjms7rfm26",
}
KEY_VERSION_MAP = {
    "dev": "1.0.4", "test": "1.0.4", "test2": "1.0.4", "training": "1.0.4",
    "prerelease": "1.2.2", "release": "1.2.2",
}
ENV_MAP = {"0": "target173.example.com", "1": "target173.example.com", "2": "target173.example.com",
           "3": "target174.example.com", "4": "target175.example.com", "5": "target173.example.com"}


def dh_encrypt(obj, key: str) -> str:
    """JSON → zlib deflate → AES-CBC(IV=key, Pkcs7) → Base64"""
    raw = obj if isinstance(obj, str) else json.dumps(obj, separators=(",", ":"), ensure_ascii=False)
    comp = zlib.compress(raw.encode("utf-8"))
    return base64.b64encode(AES.new(key.encode(), AES.MODE_CBC, key.encode()).encrypt(pad(comp, 16))).decode()


def dh_decrypt(b64: str, key: str) -> str:
    ct = base64.b64decode(b64)
    return zlib.decompress(unpad(AES.new(key.encode(), AES.MODE_CBC, key.encode()).decrypt(ct), 16)).decode("utf-8")


def dh_sign(obj_str: str, nonce, timestamp, key: str) -> str:
    """Base64(HmacSHA1(objData + nonce + timestamp, key))"""
    msg = "%s%s%s" % (obj_str, nonce, timestamp)
    return base64.b64encode(hmac.new(key.encode(), msg.encode(), hashlib.sha1).digest()).decode()


def build_headers(player_id, env="test2"):
    key, keyid, ver = KEY_MAP[env], KEY_ID_MAP[env], KEY_VERSION_MAP[env]
    obj_data = dh_encrypt(json.dumps({"playerId": player_id}), key)
    ts = int(time.time() * 1000)
    obj1 = {"encrypted": True, "gzipped": True, "platform": "h5", "version": ver,
            "application": "game_http", "timestamp": ts, "nonce": 0,
            "sign": dh_sign(obj_data, 0, ts, key), "keyid": keyid}
    return {"Content-Type": "application/json;charset=UTF-8", "deviceType": "6",
            "keyid": keyid, "token": dh_encrypt(json.dumps(obj1), key),
            "application": "game_http"}, ENV_MAP["1"]


if __name__ == "__main__":
    print("[self-test] encrypt/decrypt roundtrip")
    s = dh_encrypt({"playerId": 488283}, KEY_MAP["test2"])
    print("  enc:", s[:50])
    print("  dec:", dh_decrypt(s, KEY_MAP["test2"]))
    h, dom = build_headers(488283)
    print("[gateway]", "https://gateway.%s/game-http/player/getPlayerSetting?deviceType=6&playerId=488283" % dom)
    print("[headers]", {k: (v[:40] + "...") for k, v in h.items()})
