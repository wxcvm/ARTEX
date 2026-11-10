#!/usr/bin/env python3
"""被动找跨段新资产：主域死亡时，用姊妹域的证书 + 机器命名规律反推落点。

适用场景：
  主域名已 NXDOMAIN（被动 DNS 源也查不到历史），但同团伙还有其它在用域名。
  站群/集群类目标的分站子域通常有固定命名规律（例：5 小写字母 + 2 位编号，
  如 cdgpy11-14 / yvrgh01-04 / kabjg51-54），每 4 台一组、每组落到同一 IP。
  把规律提取出来批量解析，就能发现指向新网段的新机器。

流程：
  1) crt.sh 拉证书名（%25 = URL 编码的 '%'，匹配所有子域）
  2) 提取唯一子域，去掉 *. 通配
  3) 批量 A 解析（注意：工具若是 Linux 二进制在 Windows 跑不了，
     直接用 PowerShell 的 Resolve-DnsName）
  4) 按"解析到的 IP"分组 → 每个 IP 就是一台（或一组）机器，
     不在已有资产库里的就是新资产
"""
import json, subprocess, sys, collections

DOMAINS = ["example.com", "example.com", "target56.example.com"]   # 姊妹域，替换成目标团伙的
PROXY = "socks5-hostname 127.0.0.1:1086"                  # 出口代理（按需）

def crt(domain):
    url = f"https://crt.sh/?q=%25.{domain}&output=json"
    raw = subprocess.run(["curl.exe", "-s", "--proxy", PROXY, "--max-time", "90", url],
                         capture_output=True, text=True).stdout
    try:
        return {n for row in json.loads(raw or "[]")
                for n in row.get("name_value", "").split("\n") if n and not n.startswith("*")}
    except Exception:
        return set()

# 1) 收集所有子域
subs = sorted({s for d in DOMAINS for s in crt(d)})
print(f"[crt.sh] unique subdomains = {len(subs)}")
open("subs.txt", "w").write("\n".join(subs))

# 2) 解析（PowerShell 侧执行，示例）：
#    $res = foreach ($d in Get-Content subs.txt) {
#      $r = Resolve-DnsName -Name $d -Type A -DnsOnly -EA SilentlyContinue | ? IPAddress
#      if ($r) { foreach ($x in $r) { "$d $($x.IPAddress)" } } }
#    $res | Set-Content dns.txt
# 3) 按 IP 分组：sort -k2 dns.txt | awk '{print $2}' | uniq -c
# 4) 逐 IP 只读验证是否同套系统（响应头指纹，例如 x-is-origin / X-Powered-By）
