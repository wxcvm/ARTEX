# Sub2API invalid_auth 限流打满 PoC
# 原理：middleware/api_key_auth.go:37 先于鉴权调用 rejectInvalidAuthAbuse()
#      -> ingress_reject.go:44-46 invalidAuthClientKey = normalizeIngressRejectIP(SecurityClientIP(c))
#      -> target486.example.com:2379-2380 默认 threshold=120 / block_seconds=60
# 结果：第 121 次 429 Retry-After: 60；桶被 /v1beta、/images 共享；伪造转发头无效

import sys, time, urllib.request, urllib.error

BASE = sys.argv[1].rstrip("/") if len(sys.argv) > 1 else "http://127.0.0.1:8080"
URL = BASE + "/v1/models"
N = 130

def get(headers=None):
    req = urllib.request.Request(URL, headers={"User-Agent": "Mozilla/5.0", **(headers or {})})
    try:
        r = urllib.request.urlopen(req, timeout=10)
        return r.status, r.headers.get("Retry-After"), r.read(200).decode("utf-8", "replace")
    except urllib.error.HTTPError as e:
        return e.code, e.headers.get("Retry-After"), e.read(200).decode("utf-8", "replace")
    except Exception as e:
        return -1, None, repr(e)

t0 = time.time()
first429 = None
for i in range(1, N + 1):
    st, ra, body = get({"x-api-key": f"sk-invalid-{i}"})
    print(f"{i}\t{st}\tRA={ra}\tt={time.time()-t0:.1f}s\t{body[:80]}")
    if st == 429 and first429 is None:
        first429 = i
        print(f"*** 429 first hit at request #{i}, Retry-After={ra}")
        break

# 封禁期内：伪造转发头一律应仍为 429（证明客户端头不能重置限流 key）
for h in ({"X-Forwarded-For": "203.0.113.77"},
          {"X-Real-IP": "203.0.113.78"},
          {"CF-Connecting-IP": "203.0.113.79"},
          {"Forwarded": "for=203.0.113.80"}):
    st, ra, body = get({"x-api-key": "sk-other", **h})
    print(f"[hdr] {list(h)[0]}\t{st}\tRA={ra}\t(期望 429)")

# 对照：非 apiKeyAuth 链应为 200
def probe_public():
    req = urllib.request.Request(BASE + "/api/v1/settings/public", headers={"User-Agent": "Mozilla/5.0"})
    try:
        r = urllib.request.urlopen(req, timeout=10); return r.status
    except urllib.error.HTTPError as e:
        return e.code
print(f"[ctrl] /api/v1/settings/public -> {probe_public()} (期望 200，证明仅 apiKeyAuth 链被封)")
