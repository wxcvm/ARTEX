# 出口 IP 被宝塔 WAF 拉黑后的应急出口：免费 HTTP 代理池（含 Windows schannel 的坑）

## 适用场景
目标站对**我方所有固定出口 IP**（本机宽带 / 自建 VPS）返回 nginx 444（HTTPS 表现为
`schannel: server closed abruptly` 或 `Empty reply from server`，HTTP 表现为 444），
但 TCP 22/80/443/888 全 OPEN、第三方拨测（check-host.net）显示站点正常。
→ 这是**站点级 IP 黑名单**，唯一解是换出口 IP；换 vhost / 伪造 XFF / 改 UA / 畸形请求全部无效。

## 关键结论（实测）
1. **X-Forwarded-For / X-Real-IP / Client-IP / True-Client-IP / CF-Connecting-IP 全无效**——宝塔按 socket 源 IP 判定；空 SNI / 假域名 SNI / HTTP/1.0 / 双 Host 也无绕过空间。
2. **免费代理池可用但有寿命**：单个代理 IP 通常只撑 **2–4 个请求**就被目标拉黑 → 脚本必须"每 2 个请求换一个代理、失败重试 N 次"。
3. **Windows 自带 curl（schannel）必须加 `--ssl-no-revoke`**，否则吊销检查会让请求直接失败（表现为 000）。
4. 抓源建议走本机已有代理（例：`http://127.0.0.1:7892`）；`target289.example.com` 走 python urllib 常报 SSL UNEXPECTED_EOF，用 curl 抓更稳。
5. **判据要用目标站本身**（`GET /example.com` → `200|3933`），不要用百度等无关站点验活——只有"能拿到目标页"的代理才有用。
6. 免费代理会看到流量（cookie、马内容）——仅在没有其它出口时使用，并在报告里写明该取舍。

## 可复现流程
```
# 1) 抓源（走本机代理，curl）
curl -s --max-time 30 -A "Mozilla/5.0" --proxy http://127.0.0.1:7892 \
  https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt -o p1.txt
curl -s --max-time 30 -A "Mozilla/5.0" --proxy http://127.0.0.1:7892 \
  https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt -o p2.txt
# 2) 提取 ip:port 去重（正则 (\d{1,3}\.){3}\d{1,3}[:<\s/]*?\d{2,5}）
# 3) 并发验活：对每个代理请求 https://TARGET/example.com，保留 200 且 body>1000 的
#    （python 用 ssl.CERT_NONE + ProxyHandler；或 curl --proxy http://$p --ssl-no-revoke -k）
# 4) 得到可用池后，攻击链用"每 2 请求换代理 + 失败重试 40 次"
```
本机脚本：runs/proxy_hunt2.py（抓源+GET/POST 双判据）、runs/pw_attack2.py（轮换池 + 完整攻击链）、
runs/pw_hit_once.ps1（单代理原子序列一发命中）。

## 反面经验
- 只用"GET 能过"筛代理不够：很多免费代理能 GET 但 POST 返回 502 / 断连 → 必须**用 POST 做判据**（POST 一个无害表单，看是否 200）。
- 每发一次请求都在消耗代理 IP 的寿命，所以要把链路压到最短（登录→抓 verify→上传，3 发以内）。
