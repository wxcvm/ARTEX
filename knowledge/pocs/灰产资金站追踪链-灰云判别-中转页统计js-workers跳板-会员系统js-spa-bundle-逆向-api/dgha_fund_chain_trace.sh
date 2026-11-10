#!/bin/bash
# tgt9215集群 — 资金链自动追踪 + 资金站 API 逆向
# 实测生效：2026-09-28，从主 API 未授权接口泄露的资金链域名一路追到真实资金站源站
# 出口：必须在无代理环境（VPS 203.0.113.35）执行；本机 7892 代理会按 Host 分流产生假命中
set -u
UA='Mozilla/5.0 (Linux; Android 12; SM-G981B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/203.0.113.16 Mobile Safari/537.36'

echo "=== [1] 资金链入口：灰云域名跟随跳转（判别是否 CF 代理）"
for d in target13.example.com; do
  curl -skL --compressed --max-time 20 -A "$UA" -D - "https://$d/" -o /tmp/land.html 2>&1 \
    | grep -iE '^HTTP|^location|^server|^x-cache|cf-ray'
done
# 结果：target13.example.com -> target14.example.com（中转页，标题「正在载入，请稍等」）

echo "=== [2] 中转页/404 页里挖关联域名（hostArr 跳转池 + 统计脚本）"
curl -sk --compressed --max-time 15 -A "$UA" "https://target14.example.com/" -o /tmp/mj.html
grep -oE 'https://[a-zA-Z0-9._-]+' /tmp/mj.html | sort -u
grep -oE 'https://[a-zA-Z0-9._-]+/script\.js' /tmp/mj.html | sort -u
# 结果：hostArr=zad.teay-teay51[0-8]target15.example.com（6 个已轮换废弃）；404 页泄露统计域 target11.example.com

echo "=== [3] 统计域 script.js = 引流劫持脚本，挖 Workers 跳板"
curl -sk --max-time 20 -A "$UA" "https://target11.example.com/script.js" -o /tmp/s.js
grep -oE 'window\.location\.href = "[^"]+"' /tmp/s.js
grep -oE "fetch\('[^']+'" /tmp/s.js
# 结果：跳板 https://frite-morning-e253.misstom321.workers.dev/go；注释里藏后端 /get.php

echo "=== [4] Workers 跳板是纯 302 中继 —— 直接暴露真实落地站（含渠道码）"
curl -sk --max-time 20 -A "$UA" -D - -o /dev/null "https://frite-morning-e253.misstom321.workers.dev/go" \
  | grep -iE '^HTTP|^location'
# 结果：302 -> https://target16.example.com/k05/?shareName=ky007

echo "=== [5] 落地页引用的会员系统 JS 里硬编码资金站域名池与客服"
SUB='target17.example.com'
for f in index.js appinstall.js; do
  curl -sk --compressed --max-time 20 -A "$UA" "https://$SUB/xiaowei/ky58/$f" -o /tmp/$f
  echo "--- $f"
  grep -oE "'ky[0-9-]+\.com'" /tmp/$f | sort -u           # 资金站域名池
  grep -oE 'https://[a-zA-Z0-9._/-]+' /tmp/$f | sort -u    # 客服/指纹上报域
  grep -oE 'appKey[: ]+"[^"]+"' /tmp/$f
done
# 结果：autourldefault=tgt2905-172~target18.example.com；客服 target19.example.com；指纹上报 target20.example.com；appKey=bAlm3vhCBpgMhnX

echo "=== [6] 资金站域名是泛解析（阿里云 DNS）—— 必须用随机子域做 Host"
SUB2="$(head -c 8 /dev/urandom | od -An -tx1 | tr -d ' \n' | cut -c1-10).target12.example.com"
curl -skL --compressed --max-time 20 -A "$UA" -D - "https://$SUB2/?shareName=ky007" -o /tmp/spa.html \
  | grep -iE '^HTTP|^server|^x-cache|cf-ray'
grep -oE '<title>[^<]+' /tmp/spa.html
# 结果：200 / Server: nginx / 标题「开元棋牌」/ 无 cf-ray / IP 203.0.113.113（AWS HK）

echo "=== [7] 下载 SPA bundle 逆向全部后端 API 端点"
JS=$(grep -oE '/static/js/main\.[a-z0-9]+\.js' /tmp/spa.html | head -1)
curl -sk --compressed --max-time 25 -A "$UA" "https://$SUB2$JS" -o /tmp/main.js
grep -oiE 'baseUrl[^,;]{0,80}' /tmp/main.js
grep -oiE '"/(agent|member|capital|im_sys|discount)[a-zA-Z0-9_/]*"' /tmp/main.js | tr -d '"' | sort -u
# 结果：202 个端点；baseUrl=window.location.origin+"/member/commonRests"

echo "=== [8] 未授权接口验证（只读，绝不写/转账/下单）"
B="https://$SUB2"
for ep in /member/api/resources/banner /member/api/resources/logo /member/api/resources/icon; do
  printf '%s ' "$ep"
  curl -sk --compressed --max-time 12 -A "$UA" -D - -o /tmp/o "$B$ep" \
    -w 'CODE:%{http_code} SIZE:%{size_download}\n' | grep -iE '^HTTP|X-Application-Context|CODE'
  head -c 120 /tmp/o; echo
done
# 结果：三个资源接口均 200 + status:1 SUCCESS 未授权返回业务 JSON；
#       响应头泄露 X-Application-Context: application:redisson-cluster:8888
#       同站 /agent/* 资金接口返回 {"status":2,"msg":"您的账号未登录"} 需鉴权
