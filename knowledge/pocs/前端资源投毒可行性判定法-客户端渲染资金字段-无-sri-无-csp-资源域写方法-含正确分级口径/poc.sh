# 目标
判定“改前端 JS 就能把全站访客的收款二维码/地址换掉”这条路是否具备前提条件（用于资金类站点投毒可行性评估）。

# 1) 先确认渲染层是否“客户端零校验”
在前端 bundle 里定位资金字段的渲染点（本例：二维码/地址）：
  grep -n "toAddress\|QRcode:!0\|value:.*toAddress" chunks/*.js
命中形如：
  {label:"C2C"===walletType?"支付宝收款码":"二维码", value: order.toAddress||order.imgUrl, QRcode:!0}
结论句式：**地址/二维码取值 100% 来自服务端订单对象，且渲染前无签名校验、无地址白名单** →
JS 一旦被替换，服务端“逐单随机地址”的防护即被绕过。

# 2) 完整性校验三项判定（各 1 个请求）
a) SRI：抓首页 HTML，统计 `integrity=` 出现次数
   curl.exe -s -k -A "<浏览器UA>" --proxy <隧道> -D hdr.txt -o example.com https://TARGET/
   PowerShell: ([regex]::Matches((Get-Content example.com -Raw),'integrity=')).Count
   0 → 无子资源完整性校验。
b) CSP：看 hdr.txt 是否有 `Content-Security-Policy`
   无 → 任何被注入/被替换的脚本都能执行；有则看 script-src 是否 'unsafe-inline'/通配。
c) 资源来源与指纹：列出 HTML 中全部 `<script src>`/`<link href>`
   全为同源 /assets/*（Vite 构建哈希文件名）→ 篡改点 = 站点静态根；若为第三方桶/CDN → 篡改点 = 该桶写权限。

# 3) 资源域 CORS 写方法判定（1 个 OPTIONS）
   curl.exe -s -k -I -A "<浏览器UA>" --proxy <隧道> -X OPTIONS https://<资源域>/<任意图片> \
     -H "Origin: https://TARGET" -H "Access-Control-Request-Method: PUT"
看是否回 `Access-Control-Allow-Origin: https://TARGET` + `Allow-Credentials: true` +
`Allow-Methods` 含 PUT/POST/DELETE/PATCH → 允许带凭据的跨域**写**方法（配合站内 XSS 或写接口即打通）。

# 4) 缓存层
看静态资源响应头 `Cache-Control`/`X-Cache`/`X-Debug-Static`
max-age 越大 → 篡改对老客户端生效越慢（本例 604800=7 天），评估时须写清。

# 结论模板
· 无 SRI + 无 CSP + 客户端零校验 + 资源同源 → **前置条件全部具备**；
· 此时 severity 取决于“资源写权限”是否取得：
  - 已能写 /assets 或替换 example.com 引用（上传落盘到站点根 / 静态目录写 / 备份还原）→ **critical（全站资金投毒）**；
  - 尚未取得 → 记 **medium**，把“缺什么写权限”写成明确卡点，不要写成已可利用。
