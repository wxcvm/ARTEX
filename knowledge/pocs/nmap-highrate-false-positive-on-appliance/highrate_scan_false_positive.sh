#!/bin/bash
# 【方法论纠错】高速 SYN 扫描在带 SYN-flood 防护的设备上会产生海量假开放端口
#
# 实测事故：
#   nmap -Pn -n --open -T4 --min-rate 3000 --max-retries 1 -p- 203.0.113.40
#   → 报告 157 个开放端口（1297/1311/1320/1390/1538... 一直到 64738，密集随机分布）
#   → 但逐个 curl/TCP 连接全部 000 / 无响应
#
# 复核（正确做法）：
#   nmap -Pn -n -sT --open -T3 --max-retries 2 --scan-delay 8ms -p <同样端口列表> 
#   → 全部关闭。真实开放端口仅 8 个（80/443/5000/5001/5244/8090/8096/8080）
#
# 原因：
#   群晖 DSM、部分防火墙/IPS、F5 等设备带 SYN-flood 防护，在高速率 SYN 下会
#   对所有端口回 SYN-ACK（或随机回），造成 --min-rate 高速扫描的大规模假阳性。
#
# 判据（发现以下特征立刻降速复核）：
#   1. 开放端口数量异常多（几十~几百个）且端口号密集、随机、分布在高位段
#   2. 开放端口完全没有规律（不像服务端口的自然分布）
#   3. 目标同时开着 443/5000/8080 这类明确服务，但大量高位端口"也开着"
#   4. 对报告的"开放端口"做 TCP 连接或 curl 时全部无响应
#
# 正确姿势：
#   # 步骤1：快速扫（用于发现，不用于结论）
#   nmap -Pn -n -T4 --min-rate 1000 --open -p- TARGET -oN fast.txt
#
#   # 步骤2：对可疑结果做全连接低速复核（-sT 走完整 TCP 握手）
#   nmap -Pn -n -sT --open -T3 --max-retries 2 --scan-delay 8ms \
#         -p <fast.txt 中的端口列表> TARGET -oN verify.txt
#
#   # 步骤3：对复核通过的关键端口做真实协议探测（curl/banner grab）二次确认
#   for p in $(grep -oP '^\d+(?=/tcp)' verify.txt); do
#     curl -s --max-time 4 -o /dev/null -w "$p -> %{http_code}\n" "http://TARGET:$p/"
#   done
#
# 影响：不做这步复核，会把"157 个开放端口"当成真实攻击面去投入，
#       浪费大量时间并且把错误结论写进报告。
