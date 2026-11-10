#!/usr/bin/env python3
"""从单点证书反查整个站群（同类集群扩张期最有效的扩面手法之一）。

原理：同一运营方在扩张期常把**同一张 TLS 证书**部署到所有同类节点上
     （省事：一张证书覆盖所有分站域名）。因此只要拿到任一节点的证书域名，
     用 FOFA `cert="<域名>"` 反查，就能一次性拿到全网所有带该证书的资产。

步骤：
1) 先确定一个"确定的站群资产"的域名（例如某台已确认节点的后台域名/业务域名）
2) crt.sh 看该域名的证书里有哪些名字（确认它是单域名证书还是多域名证书）
3) FOFA 用 cert= 反查，分页拉全（size=100, page=1..N）
4) 去重得到唯一 IP 列表与唯一域名列表 → 逐个 DNS 解析确认落点
5) 对新域名再做一次 crt.sh 扩展（子域命名规律通常可复用）
6) 全程 passive：只查公共数据源与 DNS，不主动连目标

FOFA 注意：key 需从 ~/workspace/.env 显式加载（shell 环境变量不继承）：
  Get-Content "$env:USERPROFILE\workspace\.env" | ? { $_ -match '^FOFA_KEY=' } | % { $env:FOFA_KEY = ($_ -split '=',2)[1] }
"""
import base64, json, subprocess, urllib.parse, os, time, collections

KEY = os.environ["FOFA_KEY"]
API = "https://target417.example.com/api/v1/search/all"
FIELDS = "host,ip,port,domain"
CERT_NAME = "example.com"      # ← 换成目标站群任一节点的证书域名

def fofa(query, page=1, size=100):
    qs = urllib.parse.urlencode({
        "key": KEY,
        "qbase64": base64.b64encode(query.encode()).decode(),
        "size": size, "page": page, "fields": FIELDS,
    })
    out = subprocess.run(["curl.exe", "-s", "--max-time", "60", API + "?" + qs],
                         capture_output=True, text=True).stdout
    return json.loads(out or "{}")

# 1) FOFA cert= 反查（限速：每页间隔 2s）
rows, page = [], 1
while True:
    j = fofa(f'cert="{CERT_NAME}"', page=page)
    if j.get("error"):
        print("ERR", j.get("errmsg")); break
    res = j.get("results") or []
    rows += res
    print(f"page {page}: {len(res)} rows")
    if len(res) < 100:
        break
    page += 1
    time.sleep(2)

# 2) 去重
ips = sorted({r[1] for r in rows if len(r) > 1})
doms = sorted({r[3] for r in rows if len(r) > 3 and r[3]})
print(f"\n唯一 IP {len(ips)}：\n" + "\n".join(ips))
print(f"\n唯一域名 {len(doms)}：\n" + "\n".join(doms))

# 3) 对新域名做 crt.sh 子域扩展（子域命名规律往往可复用，能定位更多机器）
#    curl.exe "https://crt.sh/?q=%25.<域名>&output=json"
