#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
aocai 站群前端第二套密钥完全复现（离线解密工具）
============================================================
【密钥来源链（逐字还原自前端 sD()/_d()）】
  GET /aocai-api/getWebConfigNew
      → 用第一套 key(a6bqLk97V45jXVmj, AES-128-ECB/PKCS7) 解密响应
      → data.key 是一段 JS 源码：
          vxsvdax=function(){return 'lzE235xVxhe$'.replace('xhe$','');};vxsvdax()
  前端 wap app.js：
      V = window['2D7MkKneA4']          // 运行时把上面那段源码注入全局
      W = eval(V)                        // => 'lzE235xV'
      key = W + V.substring(1, 5)        // => 'lzE235xV' + 'xsvd'
  _d(cipher, key):
      if (len(key) < 12) key += 'kvsL'   // 源码里的补齐常量（明文可见）
      CryptoJS.AES.decrypt(cipher, Utf8.parse(key), {mode: ECB, padding: Pkcs7})

【最终密钥】 lzE235xVxsvdkvsL    （16 字节 → AES-128-ECB/PKCS7）

【前端原文（解混淆后）】
  '_e':function(h){ i=window['2D7MkKneA4']; i=atob(i)+i.substring(A,B);
                    if(i.length<12) i+='UOvlNS1NVh';
                    return AES.encrypt(Utf8.parse(h), Utf8.parse(i),
                                       {mode:ECB, padding:Pkcs7}).toString(); }
  '_d':function(h,i){ if(i.length<12) i+='kvsL';
                      return Utf8.stringify(AES.decrypt(h, Utf8.parse(i),
                                        {mode:ECB, padding:Pkcs7})).trim(); }
  'sD':function r(s){ ... if (s.headers[..].indexOf(..) != 0)
                          s.data = _d(s.data, m + o + q);      // 分支A：第一套 key
                      else { v = window['2D7MkKneA4']; w = eval(v);
                             s.data = _d(s.data, w + v.substring(1,5)); } }  // 分支B：第二套

【用法】
  python aocai_key2_decrypt.py '<base64密文>'
  python aocai_key2_decrypt.py --file leak_verify.json
"""
import base64, json, re, sys
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes

KEY1 = b"a6bqLk97V45jXVmj"                      # 第一套（分支A: m+o+q）
WEB_CONFIG_KEY_JS = ("vxsvdax=function(){return "
    "'\\x6c\\x7a\\x45\\x32\\x33\\x35\\x78\\x56\\x78\\x68\\x65\\x24'"
    ".replace('\\x78\\x68\\x65\\x24','');};vxsvdax()")
_EVALED = "lzE235xV"                            # = eval(WEB_CONFIG_KEY_JS)
KEY2 = (_EVALED + WEB_CONFIG_KEY_JS[1:5] + "kvsL")[:16].encode()   # 第二套（分支B）


def _unpad(b):
    if not b:
        return b
    p = b[-1]
    return b[:-p] if 1 <= p <= 16 and b.endswith(bytes([p]) * p) else b


def dec(key: bytes, ct_b64: str, mode: str = "ECB") -> str:
    m = modes.ECB() if mode.upper() == "ECB" else modes.CBC(b"\x00" * 16)
    d = Cipher(algorithms.AES(key), m).decryptor()
    raw = base64.b64decode(ct_b64)
    return _unpad(d.update(raw) + d.finalize()).decode("utf-8", "ignore")


def dec_key2(ct: str, mode: str = "ECB") -> str:
    return dec(KEY2, ct, mode)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(__doc__); print("KEY2 =", KEY2)
    else:
        print("KEY2 ECB ->", dec_key2(sys.argv[1])[:800])

# ---- 实测验证 ----
# cipher: YMkbo/hL1KhY9BMcSRVUYIy4n+5DFScBCA8ztm1OMfLWnkH98ft2ntqBhHO3wYt+4tkSqtPOGxLEMd5O3PsMtcrMNG2q9YpiojQycew/ECY=
#   → {"code":"e00001","message":"登录已经过期,请重新登录"}
# cipher: +T6yNalHNALeZfloHtKkls4o+qQm9qjvpedV+KhiX3uPNxmqCipXg1ki33P8iZK1
#   → {"code":"success","message":"成功","total":0}
