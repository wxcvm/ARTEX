# 发件域认证缺失（SPF / DKIM / DMARC 全无）判定与取证

## 一句话原理
发件人身份由收件方邮件网关用 SPF（源 IP 授权）、DKIM（签名）、DMARC（对齐策略）三者校验。
三个记录全缺 ⇒ 任何第三方都可以把 `From:` 写成 `*@品牌域`，而收件方没有任何依据拒收，
落在用户侧表现为正常邮件（可能进垃圾箱但不会被判伪造）。

**关键认知：有没有 MX 不影响能不能伪造。** 伪造发件人不需要该域具备收信能力，
封堵点在收件方的校验，不在发件域的 MX。

## 前置条件
- 只需要 DNS 出网能力，不接触目标主机，零凭据。
- 最好直连目标的**权威 NS** 查询；走本地递归解析器时负缓存会给出假 NXDOMAIN/NODATA。

## 取证四步（逐条留证）
1. **SPF**：`TXT <domain>` 是否存在 `v=spf1 ...`。
   只有 `google-site-verification=` / `protonmail-verification=` / `MS=msXXXX` 之类验证记录 ≠ 有 SPF。
2. **DMARC**：`TXT _dmarc.<domain>`。NXDOMAIN = 完全没有；NODATA = 无该记录；
   返回 `v=DMARC1; p=none` 属于「存在但监控模式」，拦截力弱，单独记一笔。
3. **DKIM**：`TXT <selector>._domainkey.<domain>`，逐个常见 selector 试，
   全部无记录（且没有从邮件头拿到真实 selector）才能判定 DKIM 缺失。
4. **MX / CAA（附加）**：MX 全 NODATA 说明该域不收信（纯 Web/分发域，常见于黑帽站群）；
   CAA 缺失说明未限制证书签发方，属旁证，**不能**当作邮件伪造的直接证据。

## 命令
```bash
# 先取权威 NS
nslookup -type=NS target506.example.com 8.8.8.8
# 再直连权威 NS（把 target507.example.com 换成实际 NS）
nslookup -type=TXT       target506.example.com            target507.example.com
nslookup -type=TXT       _dmarc.target506.example.com     target507.example.com
nslookup -type=TXT       default._domainkey.target506.example.com target507.example.com
nslookup -type=MX        target506.example.com            target507.example.com
nslookup -type=CAA       target506.example.com            target507.example.com
```

## 批量脚本 check_mail_auth.py
```python
#!/usr/bin/env python3
"""判定一组域是否缺失 SPF/DKIM/DMARC（只做 DNS 查询，不投递）。"""
import argparse, json, sys
import dns.resolver, dns.exception

SELECTORS = ["default", "selector1", "selector2", "google", "k1", "k2",
             "mail", "dkim", "s1", "s2", "smtp", "protonmail",
             "cm", "mandrill", "smtpapi"]

def mk(ns=None, timeout=5):
    r = dns.resolver.Resolver(configure=True)
    if ns:
        r.nameservers = [ns]
    r.timeout = timeout
    r.lifetime = timeout
    return r

def txt(name, r):
    try:
        return ["".join(x.decode() if isinstance(x, bytes) else x
                        for x in rr.strings) for rr in r.resolve(name, "TXT")]
    except (dns.resolver.NXDOMAIN, dns.resolver.NoAnswer):
        return []
    except dns.exception.Timeout:
        return None          # 超时 ≠ 不存在，必须重试或换 NS
    except Exception:
        return []

def q(name, rtype, r):
    try:
        return [str(x) for x in r.resolve(name, rtype)]
    except (dns.resolver.NXDOMAIN, dns.resolver.NoAnswer):
        return []
    except dns.exception.Timeout:
        return None
    except Exception:
        return []

def check(domain, ns=None):
    r = mk(ns)
    apex = txt(domain, r) or []
    spf = [t for t in apex if t.lower().startswith("v=spf1")]
    dmarc = [t for t in (txt("_dmarc." + domain, r) or [])
             if t.lower().startswith("v=dmarc1")]
    found_dkim = {}
    for s in SELECTORS:
        rec = txt(f"{s}._domainkey.{domain}", r)
        if rec:
            found_dkim[s] = rec
    mx = q(domain, "MX", r)
    caa = q(domain, "CAA", r)
    return {
        "domain": domain,
        "apex_txt": apex,
        "spf": spf,
        "dmarc": dmarc,
        "dmarc_record_present": bool(dmarc),
        "dmarc_policy": (dmarc[0].split("p=")[1].split(";")[0].strip()
                         if dmarc else None),
        "dkim_selectors_found": found_dkim,
        "mx_count": None if mx is None else len(mx),
        "caa_count": None if caa is None else len(caa),
        "verdict": "SPOOFABLE(no spf/dkim/dmarc)"
                   if not spf and not dmarc and not found_dkim else "PARTIAL/OK",
    }

if __name__ == "__main__":
    ap = argparse.ArgumentParser()
    ap.add_argument("--domains", required=True, help="一行一个域名")
    ap.add_argument("--ns", default=None, help="权威 NS IP，强烈建议指定")
    ap.add_argument("--json", default=None)
    a = ap.parse_args()
    doms = [l.strip() for l in open(a.domains) if l.strip()
            and not l.startswith("#")]
    out = [check(d, a.ns) for d in doms]
    for o in out:
        print(f"[{o['verdict']}] {o['domain']}  "
              f"spf={len(o['spf'])} dmarc={len(o['dmarc'])} "
              f"dkim={list(o['dkim_selectors_found'])} mx={o['mx_count']}")
    if a.json:
        json.dump(out, open(a.json, "w"), indent=2, ensure_ascii=False)
```

## 判据
- 命中：`spf=0` 且 `dmarc=0` 且 `dkim={}` → 判「可伪造发件域」。
- 部分命中：三条缺一到两条时按较弱结论写（如只有 SPF 缺 DMARC ⇒ 仍可绕过的余地大，
  因为 DMARC 才是强制对齐与拒收的依据）。
- `p=none`：记录存在但策略为监控，收件方仍可能放行，单独降级描述。

## 常见误判与失败原因
| 现象 | 真实原因 | 处理 |
|---|---|---|
| 全 NXDOMAIN | 查询打到了有负缓存的递归解析器 | 换权威 NS / 加 `+norecurse` 直问 |
| `TXT` 超时 | 权威 NS 限速或丢包 | 换 NS、加超时重试；超时≠不存在，不可直接下结论 |
| 找到 DKIM | 真实 selector 不在常见表里 | 只有拿到一封真实邮件头里的 `s=` 才能证否，否则写「未发现」而非「不存在」 |
| 只有 `protonmail-verification=`/`google-site-verification=` | 这是第三方所有权验证，不是 SPF | 不能当 SPF 证据 |
| 无 MX 就判「不能发信」 | 混淆了收信与伪造发件 | 伪造 From 不需要目标域 MX |
| `p=reject` 但仍判可伪造 | DMARC 的拒收由收件方执行，且子域/对齐方式有差别 | 有 `p=reject` 时应显著降级本条结论 |

## 影响面（报告怎么写）
- 可伪造 `From: *@<品牌域>`、`Return-Path` 与显示名，用于仿冒官方客服、充值通知、版本更新/APK 下载通知。
- 与「仿冒下载页/客户端下载链接劫持」类问题叠加时，链路完整度显著提升（邮件引流 → 仿冒页 → 恶意 APK）。
- 定级建议 medium（单看是配置缺失），若目标为金融/下载分发且无其它拦截手段可上调。
- **不要**在未经授权时做实际 SMTP 投递验证；DNS 层证据已足够支撑结论。
