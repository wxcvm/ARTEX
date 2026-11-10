#!/usr/bin/env python3
# 七牛：用 AK/SK 生成管理/上传签名，逐个端点验证权限（只读 listbucket / buckets / stat / domain list）
import base64, hashlib, hmac, urllib.parse, urllib.request

def b64(x): return base64.urlsafe_b64encode(x).decode().rstrip("=")
def sign(ak, sk, s):  # s = "METHOD URI\nHost: ...\n[Content-Type]\n\n[body]"
    return "QBox %s:%s" % (ak, b64(hmac.new(sk.encode(), s.encode(), hashlib.sha1).digest()))

AK, SK, BUCKET = "Lplztp2RiJ1-MujeBfWNIu2hTIDIdlxiSFDqH2N5", "yRfPFtH5c2K5UlDn7PV_GUAAPmzIkAdbAxNzyiCb", "20220909"

tests = [
    ("POST", "http://rsf.qiniu.com/list?bucket=%s&limit=3" % BUCKET, "/list", "POST /list?bucket=%s&limit=3\nHost: rsf.qiniu.com\n\n" % BUCKET),
    ("GET", "http://rs.qiniu.com/buckets", "/buckets", "GET /buckets\nHost: rs.qiniu.com\n\n"),
    ("GET", "http://rs.qiniu.com/stat/%s" % b64(BUCKET.encode()), None, "GET /stat/%s\nHost: rs.qiniu.com\n\n" % b64(BUCKET.encode())),
    ("POST", "http://api.qiniu.com/v6/domain/list?tbl=%s" % BUCKET, "/v6/domain/list", "POST /v6/domain/list?tbl=%s\nHost: api.qiniu.com\nContent-Type: application/x-www-form-urlencoded\n\n" % BUCKET),
]
for method, url, _, ss in tests:
    tok = sign(AK, SK, ss)
    req = urllib.request.Request(url, method=method, data=b"" if method == "POST" else None)
    req.add_header("Authorization", tok)
    if method == "POST": req.add_header("Content-Type", "application/x-www-form-urlencoded")
    try:
        r = urllib.request.build_opener(urllib.request.ProxyHandler({})).open(req, timeout=20)
        print(url, "->", r.status, r.read(300).decode())
    except urllib.error.HTTPError as e:
        print(url, "->", e.code, e.read(300).decode())   # 401 BadToken = 密钥无效
