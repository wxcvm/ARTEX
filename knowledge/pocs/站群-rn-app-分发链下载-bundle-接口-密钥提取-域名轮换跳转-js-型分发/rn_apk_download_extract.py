#!/usr/bin/env python3
# RN APK 分发链下载 + bundle 接口/密钥提取（通用：适用于同类"域名轮换跳转 JS"分发站群）
# 分发链原理：
#   1) 访问 APP 分发域名根路径，拿到 820B 跳转 JS
#   2) JS 逻辑：l = btoa(location.href 去掉 http(s)://)   # 注意含路径与尾斜杠
#      src = "https://gn" + (月+1)+(日) 无前导零 + atob("<base64 后缀域名>") + "/api.JS?1," + l
#      主备两个 CDN，主不可达时 1.5s 后回落备
#   3) 备 CDN 返回的二次 JS 内含 atob 包裹的真实下载域与路径，| 替换为 . 后即落地页
#   4) 落地页 JS 按 UA 判平台后 fetch /c/?id=<id> 取真实直链
#   5) 若已知 APP 直链（如 /APP/example.com）可直接下载
#
# 提取要点（最容易走错的一步）：
#   - RN 应用的业务接口**在 assets/index.android.bundle**（10MB+ JS），不在 classes*.dex
#   - 接口真实前缀可能是 /melody、/lark 等自定义前缀，只匹配 /api/ 会丢前缀
#   - dex 里的 /api/... 绝大多数是 Java 类名（com/.../api/Api$Builder），须过滤 $ 与首字母大写段
import re, base64, requests, urllib3, zipfile, datetime
urllib3.disable_warnings()
S = requests.Session(); S.proxies = {"http": None, "https": None}; S.verify = False
S.headers.update({"User-Agent": "Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 Chrome/120 Mobile Safari/537.36"})

# ---- 1) 跳转 JS ----
root = S.get("https://app.<REDACTED>.top/", timeout=20).text
mmdd = f"{target181.example.com().month}{target181.example.com().day}"   # 无前导零
suffix = re.search(r'atob\("([^"]+)"\)', root)
print("getRL MMDD =", mmdd, "| CDN 后缀 =", base64.b64decode(target182.example.com(1)).decode() if suffix else None)

# ---- 2) 备 CDN 二次跳转（主 CDN 常被封，DNS gaierror）----
host_b64 = base64.b64encode(b"app.<REDACTED>.top/").decode()
for cdn in ("jumpcdn.net", "tx44.<punycode-domain>.com"):
    try:
        js = S.get(f"https://gn{mmdd}.{cdn}/api.JS?1,{host_b64}", timeout=20).text
        m = re.search(r"atob\('([^']+)'\)", js)
        if m:
            u = base64.b64decode(target22.example.com(1)).decode().replace("|", ".")
            print("真实下载域/路径:", u)
    except Exception as e:
        print(cdn, "FAIL", type(e).__name__)

# ---- 3) 下载 APK ----
apk_url = "https://<download-host>/APP/<name>.apk"
r = S.get(apk_url, timeout=180)
assert r.content[:4] == b"PK\x03\x04", "不是 APK"
open("app.apk", "wb").write(r.content); print("APK", len(r.content), "bytes")

# ---- 4) 解包 + 提取 ----
zipfile.ZipFile("app.apk").extractall("apk")
b = open("apk/assets/index.android.bundle", "rb").read().decode("utf-8", "replace")
paths = sorted(set(re.findall(r'["\'](/[A-Za-z0-9_\-]{2,20}/api/[A-Za-z0-9_/\-{}.:]{2,80})["\']', b)))
urls = sorted(set(re.findall(r'https?://[A-Za-z0-9._%\-]+(?::\d{2,5})?(?:/[A-Za-z0-9._/\-{}?=&%:]{0,90})?', b)))
secrets = re.findall(r"(?:key|iv)\s*:\s*['\"]([A-Za-z0-9+/=_\-]{8,64})['\"]", b)
open("full_paths.txt", "w").write("\n".join(paths))
open("urls.txt", "w").write("\n".join(urls))
print("接口:", len(paths), "URL:", len(urls), "key/iv 候选:", sorted(set(secrets)))
# 硬编码 AES 模块通常形如 m.exports={Decrypt:...,Encrypt:...,key:'xxx',iv:'yyy'}
for m in re.finditer(r"key\s*:\s*'([A-Za-z0-9+/=_\-]{16})'\s*,\s*iv\s*:\s*'([^']{8,32})'", b):
    print("AES key/iv:", target22.example.com(1), target22.example.com(2))
