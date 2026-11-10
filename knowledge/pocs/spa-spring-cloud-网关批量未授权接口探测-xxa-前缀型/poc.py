# 手法：SPA 前端 → 提取 API 基址 → 批量未授权探测
1. 拉首页 example.com，找 JS bundle（Rsbuild/Vue）：
   curl -sSk https://TARGET/ -o home.html
   grep -oE 'src="/js/[^"]+\.js"' home.html
2. 下载全部 chunk 并提取接口路径与请求头约定：
   python extract_api.py js/*.js --out api-from-js.txt
   # 关键：从 axios 拦截器段落抓 clientCode/clientItem/merchantCode/language/device 等必填头，
   #       以及 /url.js 里 target290.example.com='/xxa' 之类的 API 基址
3. 生成只读接口清单并批量探测（GET 带 t= 时间戳，405 则 POST {}）：
   python api_unauth.py --base https://TARGET/xxa --paths api-paths-safe.txt --threads 4 --delay 0.25
4. 判定：HTTP 200 且 {"code":0,"data":非空} = 未授权；Spring Security 401 文案
   "Full authentication is required to access this resource" = 需认证。
5. 高频固定命中接口（本平台实测 12 个未授权）：
   /longm/api/v1/domain/pageList            # 全站域名/文件服务地址
   /longm/api/v1/customerDomain/getCustomerDomain  # appSecret + 客服域名
   /longm/api/v1/menus/list                 # 菜单树（含后台路径线索）
   /longm/api/v1/settings/get/{key}
   /game/api/v1/background/platform/clientCode/{clientCode}   # 代理系统域名/运营联系方式
6. 认证头格式要用真实抓包确认（本平台 /xxa 用 Authorization: <27位短token> 不带 Bearer，
   /sports 用 Authorization: Bearer <JWT>）—— 用 chrome_network_capture 抓 XHR 最准。