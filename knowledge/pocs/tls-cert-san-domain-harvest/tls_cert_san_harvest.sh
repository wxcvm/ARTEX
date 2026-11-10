#!/bin/bash
# TLS 证书 SAN 反查域名 —— 纯 IP 目标段挖真实域名的最高效手法
#
# 场景：扫描只拿到 IP，不知道域名（无 PTR、无 DNS 记录可用）
# 原理：443/8443 等 HTTPS 端口握手时会下发服务器证书，证书的 subject CN 与 SAN 字段
#       包含该服务真实绑定的全部域名。用 SNI 指向 IP 本身也能触发证书下发。
# 收益：实测在纯 IP 段中一次挖出 20+ 个真实业务域名（含管理后台、对象存储、支付等子域），
#       并可用挖到的域名做 vhost 路由（同一 IP 上多个服务靠 Host 头区分）。

CIDR_IPS="$@"          # 传入 IP 列表或从 nmap 结果提取
PORTS="443 8443 9443 6443 10443"

for ip in $CIDR_IPS; do
  for port in $PORTS; do
    out=$(echo | timeout 6 openssl s_client -connect "$ip:$port" -servername "$ip" 2>/dev/null \
          | openssl x509 -noout -subject -issuer -ext subjectAltName 2>/dev/null | tr '\n' ' ')
    if [ -n "$out" ]; then
      cn=$(echo "$out" | grep -oP 'CN=\K[^,]+' | head -1)
      sans=$(echo "$out" | grep -oP 'DNS:\K[^, ]+' | tr '\n' ',' | sed 's/,$//')
      echo "$ip:$port | CN=$cn | SAN=$sans"
    fi
  done
done

# 拿到域名后关键一步：用域名重新访问同一 IP（vhost 路由），
# 常见会命中完全不同的服务：
#   curl -sk -H "Host: <domain>" https://IP/
#   curl -sk -H "Host: <domain>" https://IP:8443/
#
# 进阶：把域名喂给 crt.sh 继续扩面
#   curl -s "https://crt.sh/?q=%25.<domain>&output=json" | jq -r '.[].name_value' | sort -u
