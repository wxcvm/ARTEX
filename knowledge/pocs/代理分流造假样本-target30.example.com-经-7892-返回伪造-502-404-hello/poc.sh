实测记录（2026-09-28）：

# 走代理（7892）——看起来"目标有反应"
$ curl -s -x http://127.0.0.1:7892 -k -A "Mozilla/5.0" -w "\n%{http_code}" https://example.com/api/v1/u/
hello
502
$ curl -s -x http://127.0.0.1:7892 -k -A "Mozilla/5.0" -w "\n%{http_code}" https://example.com/api/v1/v/
hello
404

# 本机直连复核 —— 完全不可达
$ curl -s --noproxy '*' -k -A "Mozilla/5.0" -w "%{http_code}|%{remote_ip}" https://example.com/api/v1/u/
000|            (schannel: 无法连接)
$ curl -s --noproxy '*' --path-as-is --resolve example.com:443:203.0.113.79 -k https://example.com/api/v1/u/
conn-fail
# 203.0.113.80 / 203.0.113.81 / 203.0.113.82 三个解析 IP 同样 conn-fail

结论：example.com 本机不可达；经代理返回的 {"hello", 502/404} 是代理自身派生的伪造响应，
"hello" 即代理占位体。**该域名的任何代理响应均不可作为证据**。

★ 由此推翻的错误推论：曾据 "example.com 上 /api/v1/u/ 返回 502（上游不可达）而 /api/v1/v/ 返回 404"
推断"该 nginx 存在 location /api/v1/u/，客服端网关在这台机器上"。该推论已作废。

正确入口（经直连验证）：
- https://203.0.113.85:3493/api/v1/v/init  POST {"cid":"<key>","vid":"0"} -> 200 + tk/vid/租户配置
- 该入口 --noproxy '*' 直连可稳定成功，无需域名、无需代理
