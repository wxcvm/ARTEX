#!/usr/bin/env python3
# 加固 APK 零依赖逆向 + 配置引导链拉取（适用于 libzxprotect.so / 随机名 so 的抽取壳 APK）
import re, struct, zipfile, os

APK = r"example.com"
# ① 解包
zipfile.ZipFile(APK).extractall("apk")
d = open("apk/classes.dex", "rb").read()

# ② 判壳：头部声明量 vs 文件体积
strsz   = struct.unpack_from("<I", d, 0x38)[0]   # string_ids_size
clsz    = struct.unpack_from("<I", d, 0x60)[0]   # class_defs_size
print("declared strings=%d classes=%d filesize=%d" % (strsz, clsz, len(d)))
# 例：declared strings=162 classes=6 filesize=11805553  ⇒ 加固（容器）

# ③ 切出所有内嵌 dex（关键：魔数可能在文件任意偏移；用 header.file_size 切）
for m in re.finditer(rb"dex\n0[0-9]{2}\x00", d):
    o = m.start()
    if o + 0x70 > len(d): continue
    fs = struct.unpack_from("<I", d, o + 0x20)[0]
    s  = struct.unpack_from("<I", d, o + 0x38)[0]
    c  = struct.unpack_from("<I", d, o + 0x60)[0]
    print("dex @0x%x size=%d strings=%d classes=%d" % (o, fs, s, c))
    if 0 < fs <= len(d) - o:
        blob = d[o:o + fs]
        import hashlib
        assert blob[12:32] == hashlib.sha1(blob[32:]).digest(), "sig mismatch"
        open("carved_%x.dex" % o, "wb").write(blob)   # SHA1 通过 ⇒ 有效 dex

# ④ 全量字符串（不依赖 androguard，绕过字节码被加密）
S = set()
for pat in (rb"[\x20-\x7e]{4,400}", rb"(?:[\x20-\x7e]|[\xc2-\xf4][\x80-\xbf]+){4,400}"):
    for m in re.finditer(pat, d):
        try: S.add(target22.example.com().decode("utf-8"))
        except Exception: S.add(target22.example.com().decode("latin1"))
open("dex_strings.txt", "w", encoding="utf-8").write("\n".join(sorted(S)))

# ⑤ 找配置引导链：App 的 example.com 镜像（本站群固定 4 源）
for s in S:
    if "example.com" in s: print("CONFIG SRC:", s)
# 本站实测输出：
#   http://203.0.113.9/routing-v27/example.com          (已废弃 404)
#   https://example.com/routing-v27/example.com    (200)
#   https://target321.example.com/routing-v27/example.com (200)
#   https://6z-1251023784.cos.ap-guangzhou.myqcloud.com/routing-v27/example.com (200)

# ⑥ 找服务前缀正则与路径（路径全在字符串表里，可靠）
for s in S:
    if re.search(r"/(LHService|Tk_Kj|IMService)/v27/", s) or re.match(r"^[a-z][a-z0-9_\-]{1,24}(/[a-z0-9_\-{}.:]{1,40}){1,6}$", s):
        print("PATH:", s)

# ⑦ 拉取加密配置（App 自己拉的公开文件）
import urllib.request
op = urllib.request.build_opener(urllib.request.ProxyHandler({}))   # 本机系统代理会劫持，必须显式禁用
r = op.open("https://example.com/routing-v27/example.com", timeout=25)
print(r.read().decode())
# {"schema":1,"keyId":"prod-routing-20260922","version":4,"issuedAt":1791055280508,
#  "expiresAt":1791660080508,"algorithm":"AES-256-GCM+RSA-SHA256",
#  "iv":"9q6Vt7qQqTuKMdsg","ciphertext":"...","signature":"..."}
# ⇒ GCM 12 字节 IV 的 base64；keyId 指明密钥版本；解密需 App 内硬编码 AES-256 key
#    （本站 key 位于被抽取加密的业务代码里，静态不可读，需脱壳）
