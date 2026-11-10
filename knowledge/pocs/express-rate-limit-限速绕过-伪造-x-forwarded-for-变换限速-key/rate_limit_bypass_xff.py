import random, json, urllib.request, urllib.error

def rnd_ip():
    return "%d.%d.%d.%d" % (random.randint(11, 222), random.randint(0, 255),
                            random.randint(0, 255), random.randint(1, 254))

OPENER = urllib.request.build_opener(urllib.request.ProxyHandler({}))  # 显式禁用代理

def call(url, method="GET", body=None, token=None, extra=None):
    xff = rnd_ip()
    hdrs = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/131.0 Safari/537.36",
        "Accept": "application/json, text/plain, */*",
        "Content-Type": "application/json",
        "X-Forwarded-For": xff,
        "X-Real-IP": xff,
        "Client-IP": xff,
        "X-Originating-IP": xff,
    }
    if token:
        hdrs["Authorization"] = "Bearer " + token
    if extra:
        hdrs.update(extra)
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(url, data=data, headers=hdrs, method=method)
    try:
        with OPENER.open(req, timeout=25) as r:
            return r.status, r.read(20000).decode("utf-8", "replace")
    except urllib.error.HTTPError as e:
        return e.code, e.read(20000).decode("utf-8", "replace")
    except Exception as e:
        return -1, str(e)

if __name__ == "__main__":
    print(call("http://203.0.113.105:15555/api/auth/google/url"))
