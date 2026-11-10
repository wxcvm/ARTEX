#!/bin/bash
# 上游/供应链域名「可控性」判定 —— 只做查询，不做任何注册/抢注/联系
# 用法：把待判域名填进 DOMAINS，在 VPS 上跑（本机出口常被 Cloudflare 挑战/被墙）
# 判定三选一：可控（已过期/待释放/parked/待售）| 不可控（在册+活业务+近期证书）| 不确定

DOMAINS="target360.example.com target364.example.com target365.example.com target366.example.com target367.example.com target368.example.com target369.example.com target370.example.com target371.example.com target372.example.com target373.example.com"

# ---------- 1) whois（默认 referral）----------
for d in $DOMAINS; do
  echo "### $d"
  whois "$d" 2>/dev/null | grep -iE 'domain name:|registrar:|creation date|created:|updated date|registry expiry|expiry|expiration|domain status:|name server:|dnssec:'
  sleep 3
done

# ---------- 2) RDAP 权威补位（whois 被 privacy/rate-limit 挡住时用）----------
# IANA bootstrap 找端点：curl -s https://target374.example.com/rdap/dns.json
rdap(){ curl -sS --max-time 25 -H 'Accept: application/rdap+json' "$2" \
  | jq -r '"status: \(.status//[]|join(","))", ([.events[]?|"\(.eventAction): \(.eventDate)"]|join("\n")), ([.nameservers[]?|"ns: \(.ldhName)"]|join("\n"))'; }
rdap x https://target375.example.com/com/v1/domain/EXAMPLE.com      # .com
rdap x https://target375.example.com/net/v1/domain/EXAMPLE.net      # .net
rdap x https://target376.example.com/rdap/domain/target377.example.com        # .space/.site (Radix)
rdap x https://target378.example.com/rdap/domain/target379.example.com  # .pro/.kim (Identity Digital)
# 注意：.vip/.xyz/.top 的 RDAP 端点因注册局而异，先查 IANA bootstrap

# ---------- 3) 指定注册局的 whois（部分注册局不支持默认 referral）----------
whois -h target380.example.com  target371.example.com     # NameCheap 客户
whois -h target381.example.com      target382.example.com         # Radix（可能 rate limit，隔 30s 重试）
whois -h target383.example.com x.com         # 注册局层（只给创建/到期/NS/状态，最稳）

# ---------- 4) DNS：解析链（判"是否还有解析 / 是不是 parked IP"）----------
for d in $DOMAINS; do
  echo "$d  NS=$(dig +short NS $d @203.0.113.1 | tr '\n' ' ')"
  dig "$d" A @203.0.113.1 +noall +answer +authority   # 看完整 CNAME 链
  sleep 2
done

# ---------- 5) HTTP/HTTPS：单次只读 GET，判 parked/活业务/NXDOMAIN ----------
# 每域名只打一次，间隔 ≥5s，不要骚扰别人的站
curl -sS -k -L --max-time 20 -A 'Mozilla/5.0' \
  -o /tmp/body.html -w '%{http_code} %{size_download} %{url_effective} %{remote_ip}\n' "$url"
grep -oiE '<title[^>]*>[^<]{0,120}' /tmp/body.html | head -1
sleep 5

# ---------- 6) CT（判"大约何时停用"）----------
curl -sS --max-time 40 "https://crt.sh/?q=%25.$d&output=json" -o ct.json
jq -r 'sort_by(.entry_timestamp)|reverse|.[0:5]|.[]|"\(.entry_timestamp) expires=\(.not_after) \(.common_name)"' ct.json

# ---------- 7) TLS 证书实际有效期（最准的"还活着吗"证据）----------
openssl s_client -connect $d:443 -servername $d </dev/null 2>/dev/null \
  | openssl x509 -noout -subject -issuer -dates

# ---------- 8) 家族挖掘（关键！单个域不可控不等于家族不可控）----------
# a) CT 历史里轮换的 common_name 往往就是同一供应商的兄弟域：
jq -r '.[].common_name' ct.json | sort -u
# b) NS 归并：同一组自建 NS 基本就是同一运营者
# c) 创建/到期时间戳比对：精确到同一秒 = 批量注册 = 同一批人
# d) 对挖出来的兄弟域重复 1)-7)
