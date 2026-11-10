import requests, hmac, hashlib, base64, json, urllib.parse

AK = "<AccessKey>"; SK = "<SecretKey>"; B = "<bucket>"
def ub(x): return base64.urlsafe_b64encode(x).decode()
def sign(s): return "QBox %s:%s" % (AK, ub(hmac.new(SK.encode(), s, hashlib.sha1).digest()))
def canon(p): return urllib.parse.urlencode(sorted(p.items())) if p else ""
def api(method, host, path, params=None, data=None, ctype="application/x-www-form-urlencoded"):
    q = canon(params); full = path + ("?" + q if q else "")
    ad = (full + "\n" + ctype + "\n" + data.decode()).encode() if (method=="POST" and data is not None) else (full+"\n").encode()
    h={"Authorization": sign(ad)}
    if data is not None: h["Content-Type"]=ctype
    return requests.request(method, host+full, headers=h, data=data, timeout=60)

print(api("GET","https://rs.qiniu.com","/buckets").text)
print(api("GET","https://rs.qiniu.com","/bucket/"+B).text)              # private:0/1
print(api("GET","https://api.qiniu.com","/v6/count",{"bucket":B,"begin":"2026-10-18","end":"2026-10-20","g":"day"}).text)
print(api("GET","https://api.qiniu.com","/v6/space",{"bucket":B,"begin":"2026-10-18","end":"2026-10-20","g":"day"}).text)
print(api("GET","https://api.qiniu.com","/v6/domain/list",{"tbl":B}).text)
print(api("GET","https://target465.example.com","/list",{"bucket":B,"delimiter":"/","limit":1000}).text[:800])
# 匿名读（无需任何凭据）：https://<绑定CDN域>/<key>
print(requests.get("https://<cdn-domain>/"+"/".join(["<tenant>","wx","<date>","<uuid>.png"]), timeout=60).status_code)
# 写/删（有权限时）：upload_token policy {"scope": B+":"+key, "deadline": now+3600}
# 删除：POST https://rs.qiniu.com/delete/<urlsafe_b64("bucket:key")>
