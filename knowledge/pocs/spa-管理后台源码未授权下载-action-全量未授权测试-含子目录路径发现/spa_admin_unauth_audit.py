#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""管理后台 SPA 源码获取 + action 全量未授权测试

## 适用场景
现代 SPA 管理后台（React/Vue 打包成 ES module），前端源码未授权可下载，
源码里含完整 API 契约（action 名、路径、请求头规范、签名算法）。

## 关键手法

### 1. 找到真实的静态资源路径
SPA 常有路由 fallback：请求 /IMService/example.com 可能返回 404，而请求根路径
/IMService/ 返回首页 HTML。真实资源往往在 **子目录** 下（如 /IMService/admin/）。
读首页 HTML 的 <script src="..."> 拿到文件名，再枚举候选目录：
    /IMService/admin/ , /IMService/static/js/ , /IMService/assets/ , /IMService/js/
命中后（Content-Type: application/javascript）即拿到全部业务逻辑。

### 2. 顺着 import 语句把模块拉全
入口 JS 里会有 `import ... from './xxx.js'` 或 `import('./xxx.js')`，
正则提取后逐个下载 —— 通常能拿到 20~30 个子模块（api.js 定义接口契约，
各业务模块定义 action 名）。

### 3. 提取 action 名做全量未授权测试
从源码正则提取 action 名（形如 `action: 'a.b.c'`、`call('admin','query','x.y')`），
对聚合入口发**不带任何 session 的**请求，按响应码分类：
    401 => 接口存在且要求鉴权（鉴权健全）
    404 => 该 action 不属于此通道
    200 => 未授权可读（漏洞！）

## 重要：请求头规范（缺了会全部返回 400 INVALID_INPUT）
这类自研安全协议常要求：
    Content-Type: application/json
    X-Security-Version: 2
"""
import ssl, urllib.request, urllib.error, json, re, os, sys, time

PROXY = "http://127.0.0.1:7892"
ctx = ssl.create_default_context(); ctx.check_hostname = False; ctx.verify_mode = ssl.CERT_NONE
op = urllib.request.build_opener(
    urllib.request.ProxyHandler({"http": PROXY, "https": PROXY}),
    urllib.request.HTTPSHandler(context=ctx))
UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/131 Safari/537.36"


def fetch(url, timeout=30, n=3000000, extra=None):
    h = {"User-Agent": UA, "Connection": "close", "Accept": "*/*"}
    if extra:
        h.update(extra)
    r = urllib.request.Request(url, headers=h)
    try:
        with op.open(r, timeout=timeout) as resp:
            return resp.status, resp.read(n), dict(resp.headers)
    except urllib.error.HTTPError as e:
        return e.code, e.read(n), dict(e.headers)
    except Exception as e:
        return -1, f"{type(e).__name__}: {e}".encode(), {}


def download_spa(base, subdirs=("admin/", "static/js/", "assets/", "js/", "")):
    """下载 SPA 前端源码，返回 {filename: source}"""
    st, html, _ = fetch(base)
    if st != 200:
        print(f"[!] 首页 {base} -> {st}")
        return {}
    text = html.decode("utf-8", "replace")
    refs = re.findall(r'(?:src|href)\s*=\s*["\']([^"\']+\.(?:js|css))["\']', text)
    print(f"[+] 首页引用的资源: {refs}")
    got = {}
    for ref in refs:
        fname = ref.split("/")[-1].split("?")[0]
        for sub in subdirs:
            url = base.rstrip("/") + "/" + sub + fname
            st, raw, h = fetch(url, timeout=25)
            if st == 200 and len(raw) > 500:
                got[fname] = raw.decode("utf-8", "replace")
                print(f"    [200] {url}  ({len(raw)} bytes)")
                # 顺着 import 拉子模块
                for m in set(re.findall(r'''(?:from|import\s*\()\s*["\']\./([A-Za-z0-9_.-]+\.js)["\']''',
                                        got[fname])):
                    st2, raw2, _ = fetch(base.rstrip("/") + "/" + sub + m, timeout=25)
                    if st2 == 200 and len(raw2) > 100:
                        got[m] = raw2.decode("utf-8", "replace")
                        print(f"      [200] {m}  ({len(raw2)} bytes)")
                break
    return got


def extract_actions(sources):
    """从源码提取 action 名"""
    actions = set()
    for src in sources.values():
        for m in re.findall(r'action\s*[:=]\s*["\'`]([A-Za-z0-9_.\-]{2,60})["\'`]', src):
            actions.add(m)
        for m in re.findall(r'["\'`]([a-z][a-zA-Z0-9_]*(?:\.[a-z][a-zA-Z0-9_]*){1,3})["\'`]', src):
            actions.add(m)
    return sorted(actions)


def test_unauth(base_url, query_path, actions, timeout=15):
    """对聚合入口发未认证请求，按响应码分类"""
    res = {"200": [], "401": [], "404": [], "other": []}
    url = base_url.rstrip("/") + query_path
    for a in actions:
        body = json.dumps({"action": a, "payload": "{}", "operationId": "0" * 20}).encode()
        h = {"User-Agent": UA, "Content-Type": "application/json",
             "X-Security-Version": "2", "Connection": "close"}
        r = urllib.request.Request(url, data=body, headers=h, method="POST")
        try:
            with op.open(r, timeout=timeout) as resp:
                st, txt = resp.status, resp.read(1000).decode("utf-8", "replace")
        except urllib.error.HTTPError as e:
            st, txt = e.code, e.read(1000).decode("utf-8", "replace")
        except Exception as e:
            st, txt = -1, str(e)
        key = str(st) if str(st) in res else "other"
        res[key].append((a, txt[:80]))
        if st == 200:
            print(f"  [!!! 200] {a}  {txt[:150]}")
        time.sleep(0.04)
    return res


if __name__ == "__main__":
    BASE = "https://example.com/IMService/"
    srcs = download_spa(BASE)
    acts = extract_actions(srcs)
    print(f"\n[+] 提取到 {len(acts)} 个 action 名")
    res = test_unauth("https://example.com/IMService/v27", "/admin/query", acts)
    print(f"\n=== 未授权测试结果: 200={len(res['200'])}  401={len(res['401'])}  404={len(res['404'])}")
    if not res["200"]:
        print("    结论：管理 API 鉴权到位，无未授权接口")
