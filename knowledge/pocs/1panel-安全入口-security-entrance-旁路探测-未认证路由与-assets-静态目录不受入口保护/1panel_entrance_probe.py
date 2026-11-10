import time, urllib.request, urllib.error

BASE = "http://<TARGET>:9527"
BLOCK_LEN = 1410   # 1Panel 安全入口拦截页固定长度（需先实测确认）

PATHS = ["/assets/", "/favicon.ico", "/api/v1/auth/login", "/api/v1/auth/captcha",
         "/api/v1/auth/mfa", "/api/v1/dashboard/base/os", "/api/v1/settings/search",
         "/api/v1/files/search", "/api/v1/hosts/ssh/search", "/debug/pprof/",
         "/metrics", "/health", "/swagger/", "/terminal", "/ws",
         "/api/v1/system/version", "/api/v1/dashboard/version", "/robots.txt"]

def get(path):
    req = urllib.request.Request(BASE + path, headers={"User-Agent": "Mozilla/5.0"})
    try:
        r = urllib.request.urlopen(req, timeout=10)
        return r.status, r.read(200000)
    except urllib.error.HTTPError as e:
        return e.code, e.read(200000)
    except Exception as e:
        return None, str(e).encode()

for p in PATHS:
    code, body = get(p)          # 注意：POST /api/v1/auth/login 需用 POST 才命中路由
    tag = "*** 未被入口拦截 ***" if len(body) != BLOCK_LEN else ""
    print("%-34s code=%-5s len=%-7s %s" % (p, code, len(body), tag))
    time.sleep(0.25)