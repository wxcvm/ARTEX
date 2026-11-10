import json, time, hashlib, random, base64

NS = "_glaxy_tus5r4_"
V = "1.0.0"

def E(s):
    # ★ 关键：E 不是纯倒序 s[::-1]，而是「字符降序」sorted().reverse()
    # 混淆来源：lgclub app.js 里 m.E = q(e){return m(e).sort().reverse().join("")}
    return "".join(sorted(s, reverse=True))

def md5(s):
    return hashlib.md5(s.encode()).hexdigest()

def rand6():
    return "".join(str(int(9*random.random())) for _ in range(6))

def sign_headers(data, domain, appid, token=None, u2token=None):
    qid = md5(str(int(time.time()*1000)) + rand6())   # 时效性 nonce，必须现场重算
    js = json.dumps(data, separators=(",", ":"), ensure_ascii=False)
    tail = (token + (u2token or "")) if token else ""
    sig = md5(E(js) + qid + appid + V + tail)         # token 与 u2token 直接拼接、无分隔符
    h = {
        "Accept": "application/json, text/javascript, */*; q=0.01",
        "Content-Type": "application/json;charset=UTF-8",
        "X-Requested-With": "XMLHttpRequest",
        "X-WEBSITE-CODE": "CLUB_H86_PC",              # ★ 缺这个会被拒
        "appId": appid,
        "domainName": domain,
        "needEncrypt": "undefined",
        "orginType": "PC",                            # ★ 缺这个会被拒
        "productId": PRODUCT_ID,
        "qid": qid,
        "sign": sig,
        "v": V,
    }
    if token:
        h["token"] = token
    return h

# Step 1: 取匿名 webToken
r = post(f"/{NS}/webToken", {"productId": PRODUCT_ID}, sign_headers({"productId": PRODUCT_ID}, HOST, APPID))
info = json.loads(base64.b64decode(r["body"]["info"]))   # -> deviceId, u2token, token

# Step 2: 带 token 调 club/proxy（中央路由网关）
payload = {"uri": "/club/init", "requestMethod": "GET", "productId": PRODUCT_ID}
post(f"/{NS}/club/proxy", payload,
     sign_headers(payload, HOST, APPID, info["token"], info["u2token"]))
# -> errCode=0000，明文返回 cryptoKey / secretKey / uploadUrl / 内网 IP 等