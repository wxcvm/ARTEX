# 知识库索引

共 **448** 条，其中已验证 **378** 条。数据源 `knowledge.db`（SQLite + FTS5，可被 ARTEX 的 `knowledge_dir` 直接加载）。

## 分类分布

| category | 条数 |
|---|---|
| `info-leak` | 86 |
| `auth-bypass` | 74 |
| `unauthorized` | 57 |
| `rce` | 42 |
| `recon` | 32 |
| `methodology` | 23 |
| `file-upload` | 20 |
| `weak-password` | 18 |
| `frontend-reverse` | 16 |
| `app-audit` | 15 |
| `env` | 11 |
| `waf` | 10 |
| `path-traversal` | 9 |
| `ssrf` | 8 |
| `tunnel` | 8 |
| `sqli` | 6 |
| `deserialization` | 5 |
| `protocol` | 4 |
| `privesc` | 4 |

## 形态分布

| kind | 条数 |
|---|---|
| `script` | 199 |
| `poc` | 147 |
| `exp` | 46 |
| `template` | 44 |
| `payload` | 8 |
| `markdown` | 3 |
| `http` | 1 |

## 全量条目

| id | 标题 | category | 组件 | kind | 已验证 | tags |
|---|---|---|---|---|---|---|
| 1 | ThinkPHP 精确版本与绝对路径识别（异常页法） | info-leak | ThinkPHP 3.1.3 | script | Y | thinkphp,fingerprint,info-leak,路径泄露,版本识别 |
| 2 | 未鉴权接口回显完整 SQL —— 表结构/资金字段/拼接写法泄露 | info-leak | ThinkPHP 3.1.3 业务接口 | poc | Y | thinkphp,sql报错回显,信息泄露,资金接口,未鉴权,mysql |
| 3 | WG-lobby（新葡京系）前端加密模型还原 + OSS 预签名写入利用链 | unauthorized | WG-lobby / lobby-web v7.5.x | exp | Y | WG-lobby,AES,密钥派生,oss,presigned,未授权,反混淆 |
| 4 | WG-lobby 游客 token 数据面拖库（getSysInfo 取 token + chipher 重放 + 响应解密 + 接口普查翻页） | unauthorized | WG-lobby（WG娱乐城前端+Go 网关，/hall/api/*） | exp | Y | wg-lobby,未授权,数据拖库,游客token,AES-CBC,翻页 |
| 5 | 暴露数据库/SSH 服务弱口令与未授权探测脚本（MySQL/PostgreSQL/SSH 三合一） | weak-password |  | script |  | mysql,postgresql,ssh,弱口令,暴露服务,未授权 |
| 6 | NextAuth v4 JWE 会话令牌解密与伪造（弱 NEXTAUTH_SECRET 检测） | auth-bypass | NextAuth.js v4 | script |  | nextauth,nextjs,jwe,forge,session |
| 7 | New API / one-api 未授权面枚举与源码鉴权对照审计清单 | info-leak | New API (Calcium-Ion/new-api) / one-api | script | Y | new-api,one-api,未授权,信息泄露,源码对照,审计清单 |
| 8 | 1Panel 安全入口（Security Entrance）旁路探测：未认证路由与 /assets/ 静态目录不受入口保护 | auth-bypass | 1Panel | poc | Y | 1panel,security-entrance,未授权,目录遍历,信息泄露 |
| 9 | 1Panel 验证码绕过（CVE-2025-66507）+ RSA 公钥加密口令在线爆破脚本 | weak-password | 1Panel | script | Y | 1panel,captcha-bypass,CVE-2025-66507,爆破,RSA |
| 10 | 1Panel 登录口令加密格式（RSA 保护 AES 密钥的三段式）与验证码绕过 | weak-password | 1Panel | script | Y | 1panel,登录加密,RSA,AES,CVE-2025-66507,爆破 |
| 11 | express-rate-limit 限速绕过：伪造 X-Forwarded-For 变换限速 key | methodology | Node.js Express + express-rate-limit | script | Y | 限速绕过,rate-limit,trust-proxy,X-Forwarded-For,Express |
| 12 | 纯 Python AWS SigV4 实现：无 awscli/boto3 时实测 AWS 凭据有效性（STS GetCallerIdentity） | methodology | AWS STS / SigV4 | script | Y | aws,sts,sigv4,凭据验证,云凭据 |
| 13 | 直连本机 ChromeMcpServer 驱动用户真实 Chrome（MCP streamable-http 自建会话 + 裸 IP 导航缺陷绕法） | env | ChromeMcpServer / mcp-chrome 类浏览器扩展桥 + MCP streamable-http | script | Y | browser,chrome,mcp,真实登录态,ChromeMcpServer,<target>,工具坑 |
| 14 | Express 鉴权中间件路径判定绕过（大小写前缀 /API + URL 编码 .html），配套变体清单与只读收割顺序 | auth-bypass | Express（自研 AI API 中转平台管理面板） | script | Y | express,鉴权绕过,路径判定,大小写,url编码,信息泄露 |
| 15 | Dujiao-Next 发卡系统：订单/游客订单鉴权模型、API 契约、订单号格式与 IDOR 判定速查 | unauthorized | Dujiao-Next 发卡/续费系统（Go+Gin+GORM+Vue3） | template | Y | dujiao-next,发卡系统,IDOR,鉴权模型,订单号格式,源码审计,速查 |
| 16 | Dujiao-Next 用户中心与认证链路速查：/me·wallet·affiliate·2FA 鉴权事实 + 验证码节流与限流缺口 | auth-bypass | Dujiao-Next 发卡/续费系统（Go+Gin+GORM+Vue3）用户中心与认证链路 | template | Y | dujiao-next,用户中心,2fa,鉴权模型,验证码节流,限流缺口,速查 |
| 17 | cloudflare_temp_email <target> 鉴权模型与匿名端点 oracle（地址存在性枚举） | unauthorized | cloudflare_temp_email | template | Y | cloudflare_temp_email,临时邮箱,未授权,用户枚举,版本差异分析 |
| 18 | LLM 中转面板渠道 config.quota.url SSRF（触发三条件 + 绕过前置 Squid 直连 + 可达性判据表） | ssrf | 自研 AI API 中转平台管理面板（Express / LLM 管理中心） | exp | Y | ssrf,llm-panel,proxy-bypass,quota,writable-config |
| 19 | Microsoft/Google OAuth refresh_token 特征 grep 清单 + aigemini-org 13030 数据面阴性结论 | info-leak |  | script |  | oauth,refresh_token,microsoft,outlook,grep,negative-result |
| 20 | AIRouter 白标 AI 网关：管理面鉴权绕过尝试集（/API/ 大小写 + /admin%2Ehtml 编码变体）与被证伪的迁移性 | auth-bypass | AIRouter（自研 AI API 包月/白标网关，Express + OpenResty） | payload |  | express,鉴权绕过,case-sensitive,url-encoding,airouter,白标网关,证伪,判别清单 |
| 21 | LLM API key 零消费只读验活（DeepSeek 余额 / OpenRouter 分层 / 通用 /models） | info-leak | LLM API Key 只读验活（余额/权限面判定） | script | Y | llm,apikey,balance,readonly,zero-cost |
| 22 | gptcdk 卡密兑换站群自托管后端未授权卡密查询/枚举探测 | unauthorized | gptcdk 卡密兑换中心（自托管卡密后端） | poc | Y | 卡密兑换中心,gptcdk,未授权访问,枚举,fastapi,c ChatGPT卡台 |
| 23 | nps web 管理端 auth_key=md5(timestamp) 认证绕过（含隧道/客户端清单导出） | auth-bypass | ehang nps (web 管理端) | poc | Y | nps,内网穿透,auth_key,认证绕过,内网入口 |
| 24 | 自研卡密兑换站未授权批量查询接口 /api/v1/cdk/batch（20 张/请求，响应含套餐/订单） | unauthorized | 自研卡密兑换站（gptcdk 同源模板） | poc | Y | unauthorized,卡密,批量查询,gptcdk,接口枚举 |
| 25 | nps 管理端接口取数姿势：/client/list 与 /index/gettunnel 必须 POST（GET 返回模板页/空集），gettunnel 必须带 type | info-leak | ehang nps (web 管理端, beego AutoRouter) | script | Y | nps,取数姿势,POST,gettunnel,云安全组,内网穿透 |
| 26 | MSSQL TDS PRELOGIN 探针（47 字节，struct 构造 + assert）—— 长度写错会被误判为"伪造开放" | info-leak | Microsoft SQL Server (TDS) | script | Y | mssql,tds,prelogin,探针,踩坑,长度校验 |
| 27 | GET/POST 分支鉴权不一致导致 userId 水平越权遍历（含无盐 MD5 口令泄露） | unauthorized |  | script | Y | 越权,IDOR,未授权访问,无盐MD5,主键遍历 |
| 28 | 方法论陷阱：本地 HTTP 代理会按 Host 分流，制造"Host 头命中源站"的假象（必须无代理复核） | recon | 探测方法论 / HTTP 代理 | script | Y | 方法论,源站定位,Host头,代理陷阱,Cloudflare,对照实验 |
| 29 | 灰产资金站追踪链：灰云判别 → 中转页统计JS → Workers跳板 → 会员系统JS → SPA bundle 逆向 API | unauthorized |  | script | Y | 灰产,资金链,域名轮换,SPA逆向,API提取,Workers跳板,未授权 |
| 30 | 对象存储凭据有效性验证器（七牛 AK/SK 只读端点四连测） | info-leak | 七牛云 / 又拍云对象存储 | poc | Y | 七牛,又拍云,对象存储,未授权配置泄露,凭据验证 |
| 31 | 未授权越权接口 · 受控批量导出用户档案（IDOR 遍历 + 弱口令反查） | unauthorized |  | poc | Y | 未授权,水平越权,IDOR,用户档案导出,弱口令,MD5 |
| 32 | Druid StatViewServlet 暴露检测（含从 actuator/prometheus 的 statViewServlet 指标反推） | info-leak | Alibaba Druid StatViewServlet | poc | Y | druid,spring-boot,actuator,监控台,信息泄露 |
| 33 | 移动 App 侦察：从 pcap 配置分发文件定位真实后端 + 认证头校难强度对照实验 | info-leak | 移动 App API（okhttp 客户端 / Tomcat-JSP 与 Spring 混合后端） | script | Y | 移动App,抓包分析,配置泄露,认证绕过,签名校验,未授权 |
| 34 | SockJS + STOMP WebSocket 客户端实现要点（JSON 数组发送侧编码 / 版本协商 / CDP 抄帧法） | unauthorized | SockJS + STOMP over WebSocket | script | Y | websocket,sockjs,stomp,协议逆向,CDP |
| 35 | 用 Actuator prometheus 的 exception 标签做参数名/字段名 oracle（异常处理器统一返回 200 时的参数枚举法） | info-leak | Spring Boot Actuator (micrometer) + Spring MVC 参数绑定 | script | Y | spring-boot,actuator,prometheus,参数枚举,oracle,micrometer |
| 36 | Android App 自研 API 加密协议还原（DEX 反汇编定位 AES 密钥 + param 密文协议 + 登录三层校验差分定位） | frontend-reverse | Android App 自研 API 加密协议（AES param 封装） | script | Y | android,dex,reverse,aes,protocol,app-api,param-encryption |
| 37 | 阿里云验证码 2.0（AliyunCaptcha）自动化跑通并取得 captchaVerifyParam —— App/Web 登录风控绕过通用打法 | auth-bypass | 阿里云验证码 2.0 / AliyunCaptcha / FeiLin | script | Y | 阿里云验证码,captchaVerifyParam,风控绕过,App登录,自动化,Burp不可见 |
| 38 | SPA bundle 端点半自动提取 + 只读白名单批量未授权探测（含写语义端点硬跳过） | unauthorized | [范围外来源·方法通用] 前端 bundle 逆向 + 只读白名单批量未授权探测 | script | Y | 未授权,SPA逆向,批量探测,只读白名单,Java,Spring,范围外来源,方法本身通用 |
| 39 | 自研 Java 资金站 TIMESTAMP 通道：安全凭证头 + AES-ECB 请求体加密与时间片密钥派生算法 | unauthorized | [范围外来源·方法通用] 自研 Java 资金站（Spring + 阿里云 WAF + Redisson）PC/H5 传输层 | script |  | AES-ECB,SHA256,时间片密钥,前端逆向,Spring,资金站,传输层,范围外来源,方法本身通用 |
| 40 | 目标不可达的「三出口对照」诊断法 —— 区分目标侧故障 vs 我方出口被封 | env | 目标可达性诊断 / 出口对照 | poc | Y | 可达性诊断,出口对照,502,排障方法论,避免白跑 |
| 41 | 阿里云 WAF 封禁维度识别：敏感端点为「方法+路径+来源」持续封禁（405 拦截页），非频率限流；附三态判据 | waf | [范围外来源·方法通用] 阿里云 WAF / 抗 D（Anti-DDoS Pro，acw_tc cookie） | script | Y | 阿里云WAF,acw_tc,Anti-DDoS,封禁维度,方法论,别再踩坑,范围外来源,方法本身通用 |
| 42 | 代理分流造假样本：<target> 经 7892 返回伪造 502/404「hello」 | methodology | http-proxy / 7892 | poc | Y | 代理造假,7892,判定方法论,负结论 |
| 43 | OCS 在线客服系统：STOMP over SockJS 握手真值 + 客服侧 protobuf 信封 UserMsgPack | protocol | OCS 在线客服系统（Spring Boot + SockJS + STOMP + protobuf） | template | Y | OCS,SockJS,STOMP,protobuf,客服侧,UserMsgPack |
| 44 | SockJS 文本通道无法承载 protobuf 二进制：成因、对照实验与判据 | protocol | SockJS / Spring SockJS | poc | Y | SockJS,STOMP,protobuf,二进制,帧长度,负结论,方法论 |
| 45 | nginx location 绕过：%2e%2e 系列必然失败，及用"不存在路径对照"证伪 ..;/ 假穿越 | path-traversal | nginx location 匹配 / 后端 Tomcat-Spring | poc | Y | nginx,location绕过,路径归一化,..;,判据,负结论 |
| 46 | nps web 管理面板 auth_key=md5(timestamp) 认证绕过 + 隧道/客户端枚举 | unauthorized | ehang nps (0.26.10 / 0.26.9) | exp | Y | nps,内网穿透,认证绕过,tunnel,面板接管 |
| 47 | 密钥分发范式：4 段密钥/盐/公钥 XOR 0x5A 隐写在 4 个未授权图片接口的尾部 2048 字节（标记 [P1][P2][S][AES]） | info-leak | [范围外来源·方法通用] 自研 Java 资金站静态资源接口（/member/api/resources/*）+ 前端密钥加载器（webpack module 1483） | script | Y | 隐写,密钥分发,图片尾部,XOR-0x5A,AES-ECB,RSA-2048,SPKI,前端逆向,响应信封,同团伙复用,范围外来源,方法本身通用 |
| 48 | 阿里云 WAF 封禁维度修正实测：换出口无效（非 IP 维度）、裸 curl/Python 必被拦、真实浏览器可通过 | waf | [范围外来源·方法通用] 阿里云 WAF / 抗 D（acw_tc cookie）—— 封禁维度实测修正版 | script | Y | 阿里云WAF,acw_tc,封禁维度修正,浏览器可过,方法论,405误判,范围外来源,方法本身通用 |
| 49 | MQTT over WebSocket 未授权接入与 ACL 绕过（`#` 被拒 → `+/#` 全量放行，EMQX） | unauthorized | EMQX / MQTT over WebSocket | poc | Y | mqtt,emqx,websocket,acl-bypass,unauth,iot |
| 50 | Spring Boot Whitelabel Error Page 参数名 oracle（Accept: text/html 读出必需 @RequestParam 名） | info-leak | Spring Boot (Whitelabel Error Page) | script | Y | spring,java,参数名泄露,信息泄露,接口枚举,oracle |
| 51 | 1Panel 验证码绕过对照判定 + 三段式加密登录探针（可直接换字典做在线口令审计） | weak-password | 1Panel | script | Y | 1panel,验证码绕过,CVE-2025-66507,三段式加密,口令审计,对照实验 |
| 52 | Tomcat AJP 8009 探测 + Ghostcat(CVE-2020-1938) 文件读取（含「HTTP 代理 CONNECT 盲回 200」的识破判据） | rce | Apache Tomcat AJP (mod_jk / AJP13 connector) | script | Y | tomcat,ajp,ghostcat,8009,代理盲回,探测判据 |
| 53 | PowerShell 下 curl 内联 JSON body 被吃双引号 → 服务端全 500 被误判为目标故障 | methodology |  | script | Y | powershell,curl,json,false-positive,探测陷阱 |
| 54 | 1Panel 版本指纹（路由集合差集法，离线）：把面板压到具体 LTS 版本 | info-leak | 1Panel | script | Y | 1panel,版本指纹,路由差集,离线,swagger |
| 55 | 教训：1Panel 版本区间误判避坑指南（GHSA first_patched_version 漏掉 LTS 分支；正确判据 = swagger basePath + router 源码集合差集） | info-leak | 1Panel | poc | Y | 1panel,版本判断,教训,swagger basePath,路由集合差集,LTS分支,GHSA陷阱 |
| 56 | 本机红队环境事实与替代方案（Windows/ELF 平台限制） | env | 本机作业环境（Windows + redteam toolkit） | script | Y |  |
| 57 | 环境事实：本机出口走全局代理（127.0.0.1:7892）；仅代理不可用时才改走演练 VPS | env | 演练环境/出口链路 | script | Y | 环境事实,出口链路,TUN,Windows,工具白名单,VPS跳板,nuclei,FOFA_KEY |
| 58 | 1Panel 安全入口拦截页与 Swagger 拼接 ⇒ 未认证绕读完整 swagger（v1.10.28-LTS 实测）+ 无 swagger 时的版本夹逼探针 | info-leak | 1Panel | poc | Y | 1panel,swagger,security-entrance,绕读,版本指纹,路由存在性探针 |
| 59 | GM Collection 邮件工具站免登录开放接口未授权探测清单（Outlook/iCloud 邮件与凭据查询） | unauthorized |  | script | Y | unauthorized,email,outlook,icloud,temp-mail,cloudflare-worker,curl |
| 60 | New API 漏洞清单与未修复版本边界对照（CVE-2026-42339 SSRF / 33655 / 64865 / 30886） | ssrf | New API (QuantumNous) / one-api | template |  | new-api,one-api,ssrf,idor,cve,版本边界,源码审计,计费绕过,yaml |
| 61 | Vue3/Vite SPA 接口测绘法：入口 JS 只给路由，真实 API 在懒加载 chunk 的数据层 use*.js 里 | frontend-reverse | Vue3/Vite SPA 前端 | script | Y | 前端接口提取,Vite,Vue,chunk,SPA,接口测绘,只读 |
| 62 | 商城游客订单自认证头机制：Authorization: Guest base64url(email+"\n"+order_password)（含空口令/畸形头边界实测） | unauthorized | 糖果姐姐商城（自研多租户 SaaS 商城 v1.4.7） | poc | Y | 商城,游客订单,自认证头,Authorization,base64,越权验证,待验证 |
| 63 | 白标 SPA 商城溯源到上游开源项目（品牌名兜底值 + app_version）→ 转离线源码审计 | frontend-reverse | 白标商城 / Vue3+Vite SPA | poc | Y | 溯源,白标商城,开源项目识别,源码审计,版本锁定,默认密钥,零负载 |
| 64 | Webpack 异步 chunk 收割枚举 SPA 全量 API 面（new-api 定制构建实测） | info-leak | Webpack SPA（new-api 及一切 webpack 构建的后台） | script | Y | webpack,spa,api枚举,异步chunk,new-api,前端信息收集 |
| 65 | 通用方法：开源组件识别 → 拉源码离线审计 → 定向验证（找 advisory 之外的路由漏挂） | methodology | 通用方法（Web 靶标 / 开源组件） | script | Y | 源码审计,中间件,漏挂认证,gin,new-api,codeload,方法论 |
| 66 | 凭据来源侦察：公开代码托管泄露检索入口清单 + 仓库全树密钥扫描法 | info-leak |  | script | Y | 信息泄露,凭据来源,代码托管,供应链,被动收集 |
| 67 | New API / one-api 网关令牌有效性判别法（Invalid token vs User has been banned） | info-leak |  | poc | Y | new-api,one-api,令牌判别,代码泄露,凭据,零爆破 |
| 68 | 被动 DNS / 证书透明可用源清单与 CDN 误判防护（实测可用性矩阵） | recon |  | template | Y | 被动DNS,证书透明,源站,CDN误判防护,信息收集 |
| 69 | CVE-2026-81027 one-api/New API channelid 路径参数越权锁定信道（任意有效 token，无需余额） | auth-bypass | one-api / New API | exp |  | one-api,new-api,越权,CWE-862,channelid,提供商key,无需余额 |
| 70 | 环境事实（全体适用）：本机 Python urllib 走 IE 注册表代理 127.0.0.1:7892 导致请求超时，curl 正常 | env |  | script | Y | 环境事实,windows,urllib,代理,curl,超时,误判 |
| 71 | Discuz 7.x ajax.php checkusername 用户名枚举（零配额成本判定账号是否存在） | info-leak | Discuz | poc | Y | discuz,用户名枚举,信息收集,零配额 |
| 72 | phpwind <target> 前台 SQL 注入（actions/ajax/leaveword.php） | sqli | phpwind | poc |  | phpwind,sql注入,front-end |
| 73 | 智讯网站群管理系统 /api/ 未授权访问（form-urlencoded 传参 + action=index 泄露站点配置） | unauthorized | 智讯网网址导航系统（自研 PHP 站群程序） | poc | Y | php,站群,未授权,api,form-urlencoded,智讯网 |
| 74 | 中文站群"网站防火墙"WAF 识别判据（403/1200B 中文页 + 按扩展名拦截 + 多长度分类法） | waf | 网站防火墙（自研/宝塔类，中文拦截页） | poc | Y | WAF,网站防火墙,403,扩展名拦截,判据,phpwind,绕过矩阵未验证,站群通用 |
| 75 | phpwind 头像上传任意文件写入（job.php?action=uploadicon&step=2&filename=x.php） | file-upload | phpwind | exp |  | phpwind,getshell,文件上传,头像,写入 |
| 76 | phpwind 附件上传扩展名黑名单审计结论（前台变形不可行，口子在后台 $db_uploadfiletype） | file-upload | phpwind | poc | Y | phpwind,文件上传,黑名单,证伪,后台getshell |
| 77 | 六零导航页 LyLme Spage ≤<target> 后台 getshell（ajax_link.php submit=update 任意 zip 解压到站点根）+ 未授权 SSRF | rce | 六零导航页 LyLme Spage v1.6.0（含 v2.7.0 前各版本同族） | exp |  | lylme,六零导航页,导航站,后台getshell,任意文件下载,zip解压,SSRF,未授权,RCE |
| 78 | 智讯网站群：注册接口用户名枚举 + 口令大小写不敏感判定 | auth-bypass | 智讯网网址导航系统（自研 PHP 站群程序） | poc | Y | php,站群,用户枚举,口令策略,注册接口,智讯网 |
| 79 | 六零导航页 LyLme Spage admin_token 离线伪造器（SYS_KEY 硬编码 lylme_key → 任意配置泄露即通吃会话） | auth-bypass | 六零导航页 LyLme Spage（v1.6.0 及同族版本） | script | Y | lylme,六零导航页,SYS_KEY硬编码,admin_token,离线伪造,cookie伪造,会话绕过 |
| 80 | 网段存活探测纪律：nmap -Pn 会把全部 IP 判为 up，必须用 -sn 复核 | recon | nmap | script | Y | nmap,存活探测,假阳性,纪律,-Pn |
| 81 | FOFA_KEY 在本机需从 .dsh/.env 显式加载（shell 环境变量未继承） | env | FOFA API 客户端环境 | script | Y | fofa,环境变量,信息收集,踩坑 |
| 82 | 主域 NXDOMAIN 后从姊妹域证书按机器命名规律挖跨段新资产 | recon | crt.sh（证书透明日志） | script | Y | crt.sh,证书透明,子域枚举,站群,被动信息收集,跨段扩面 |
| 83 | SSRF 判据三陷阱：外层码≠目标状态 / icon 非空≠端口开放 / 快速返回≠存活 | ssrf | [通用方法] SSRF 响应判读 / get_head 类 URL 解析器 | poc | Y | SSRF,判据陷阱,get_head,误判,方法论,icon字段,端口探测 |
| 84 | 假阴性判据：路径级 403 ≠ IP 封禁 / curl 默认 UA 被拦 ≠ 不可达 / 按时间下线 ≠ 被封 | waf | [通用方法] 封禁判定 / 假阴性识别 | poc | Y | 假阴性,封禁判定,路径级403,curlUA,关站,对照实验,方法论,全队踩坑 |
| 85 | 批量扫描三条对照实验：戳破假存活 / 假 403 / 假全开 | recon | nmap / Test-NetConnection 批量扫描结果校验 | template | Y | 对照实验,假存活,假全开,假403,nmap,端口验证,信息收集,方法论 |
| 86 | JWT HS256 秘密密钥组合爆破器（无依赖，含 0-99999 数字后缀与词组组合） | tunnel |  | script | Y | jwt,hs256,bruteforce,authententication |
| 87 | 自研 Vue/Vite SPA 面板：从 bundle 递归还原完整 API 面（含 /admin 与开放 API 协议） | frontend-reverse |  | template | Y | spa,vite,vue,recon,api-surface |
| 88 | 智讯网站群 root_user_index 免认证读取全站用户表（含明文 opass）｜附 curpage 分页被忽略的踩坑 | unauthorized | 智讯网站群系统 (x-is-origin: true / PHP 8.1.7) | poc | Y | 站群,未授权,用户表泄露,uid令牌,分页陷阱 |
| 89 | 从单点证书反查整个站群（cert= 一把捞全网同证书资产） | recon | FOFA / crt.sh 证书反查 | script | Y | 证书反查,fofa,crt.sh,站群扩面,通用证书,信息收集 |
| 90 | 第三方 IP 归属库不可靠：必须用响应特征自行验证 | recon | IP 归属库 / 资产归属判定 | template | Y | IP归属,资产验证,误伤风险,合规,对照实验,信息收集 |
| 91 | 指令边界歧义处理：指定 IP 白名单 vs 禁扫网段（定向优先，显式确认） | methodology | 任务边界解析 / 授权范围判定 | template | Y | 授权边界,指令歧义,禁扫网段,白名单,协作规范,指挥意图 |
| 92 | LFI 注入点定位法：先找不经过守卫的可达文件目录，再用数组参数判定是否读参数 | path-traversal |  | script | Y | LFI,任意文件读取,注入点定位,数组参数,display_errors,方法论,Windows穿越,负结论 |
| 93 | 取证表述纪律：区分「实测」与「推断」—— 同一证据存在两个同样合理解释时，未排除的那个不得写进结论 | methodology |  | poc | Y | 取证纪律,实测与推断,对照实验,假阳性,自我纠错,表述规范,前缀拦截 |
| 94 | 从 TLS 证书 CN 反查不走 CDN 的子域 ⇒ 源站 IP 泄露 + 绕过 Cloudflare（含 HTTP.sys Host 头路由命中技巧） | info-leak |  | script | Y | recon,cloudflare,源站泄露,http.sys,host头路由,cert-cn |
| 95 | 同一 ApiInterface 跨控制器对照法：定位漏挂 [Authorize] 的控制器 | unauthorized |  | script | Y | unauthorized,aspnet,mvc,authorize,反射接口,对照实验 |
| 96 | nginx `/login` 通配重定向规则会造成 ffuf 假阳性：命中项必须用『任意词对照』核实 | methodology |  | script | Y | 假阳性,nginx重写,通配规则,目录爆破,ffuf,对照实验,Location相对路径,方法论 |
| 97 | <target>/.23 先验排除表：两条恒定假阳性 + /code 验证码端点（爆破前置） + 五处已穷尽面 | methodology |  | script | Y | 先验排除表,假阳性,nginx通配规则,整前缀认证,验证码端点,已穷尽面,449918,宝塔 |
| 98 | HTTP Range 精准抽取远端 zip 内任意单文件（不下载整包） | info-leak | 通用（IIS/nginx 目录浏览 / 对象存储 / 任何支持 Range 的 HTTP 文件） | script | Y | range,zip,源码泄露,目录浏览,取证 |
| 99 | vhost 真伪判定（SNI 双置位法）：curl --resolve 同时置位 SNI+Host + 逐 SNI 证书对比 | recon |  | script | Y | vhost,sni,iis,方法论,避坑 |
| 100 | MSSQL 1433 防火墙握手代理式伪装的判定法与 TDS PRELOGIN 探针 | waf | MSSQL / SQL Server 1433 可达性判定 | script | Y | mssql,1433,TDS,可达性判定,防火墙伪装,误判避免 |
| 101 | curl --resolve 端口陷阱：URL 不带端口会静默逃逸到真实 DNS（伪造出"目标有业务"的假象） | recon |  | script | Y | curl,取证陷阱,方法论,vhost,避坑,静默失败 |
| 102 | 证书域名提取陷阱：leaf SAN 与链内 OCSP/CRL/CA-Issuers 域名必须分开（否则误当业务域名） | recon |  | script | Y | 证书,SAN,AIA,CRL,取证陷阱,资产测绘,避坑 |
| 103 | 会话态不完整时后台功能页静默返回 200/0 字节 —— 必须把"0 字节 200"当"会话可疑"而非"页面为空" | methodology |  | script | Y | aspnet,后台,会话,误判,踩坑,geek |
| 104 | 浏览器实测必须用 hosts 映射：Chrome --host-resolver-rules 在本机 153 版不生效（附 Server 头复验判据与当即还原纪律） | env |  | script | Y | 浏览器实测,hosts,cloudflare,源站直连,踩坑,cdp |
| 105 | 带验证码且无回显差异的后台弱口令打法：验证码放大 6 倍识别 + 每组口令独立会话 + 以 302 为唯一成功信号 | weak-password |  | script | Y | 弱口令,验证码,ocr,后台,aspnet,独立会话,方法论 |
| 106 | 记分被 capped 时的正确处置：追加式更正（不改 target 规避去重，追加 chain + 报告显式呈现差异） | methodology |  | script | Y | 记分,去重,capped,数据完整性,自我纠错,方法论 |
| 107 | Discuz <target> failedlogins 登录失败计数绕过（伪造 Client-IP / X-Forwarded-For 实现无限爆破） | auth-bypass | Discuz 7.2 | poc | Y | discuz,brute-force,auth-bypass,ip-spoofing,failedlogins |
| 108 | Discuz 管理员 Cookie 伪造（authkey+UA 签名 / 由 SQL 注入取三要素） | auth-bypass |  | exp | Y | discuz,auth-bypass,cookie,authkey,authcode,sqli |
| 109 | SCLUB/168six：403 站点规则可被本站 Referer 绕过（探测假阴性） | auth-bypass | SCLUB 免費論壇平台 / Apache 2.2.3 (CentOS) + PHP 5.1.6 | poc | Y | SCLUB,Discuz,Apache,403绕过,Referer,防盗链,信息收集假阴性,防盗链规则 |
| 110 | Discuz 论坛平台（SCLUB 免费论坛）连接层封禁阈值探测方法 —— 区分「并发触发」与「累计触发」 | waf | SCLUB 免费论坛平台 / Discuz 7.2 | poc | Y | discuz,brute-force,rate-limit,waf-bypass,methodology,sclub |
| 117 | 未授权接口必须"重读+diff"：静态认知固化导致漏掉新增明文凭据（19 轮作业的真实教训） | info-leak |  | template | Y | 方法论,增量监控,未授权接口,认知固化,settings,明文凭据,作业纪律 |
| 118 | 配置类接口（settings/config/option）明文回显管理员口令与 API Key —— 字段特征与取用方法 | info-leak |  | template | Y | settings,明文口令,api_key,未授权访问,defaults泄露默认口令,自研面板 |
| 119 | 本地只读反代：把前端调用的 /api/ 改写为大写转发，让原版管理后台界面直接可用 | unauthorized |  | script | Y | 本地反代,大小写绕过,只读浏览,面板,Express,CaseSensitive,工具化 |
| 120 | 主列表漏报陷阱：明细统计必须逐对象走详情接口（本例少统计 30%） | info-leak |  | template | Y | 统计口径,漏报,主从结构,apikey_path,404误解,集合比对 |
| 121 | 凭据复用验证 vs 爆破的边界：拿到实证口令后，同集群各试 1~2 次的正确打法 | weak-password |  | template |  | 凭据复用,爆破边界,1Panel,RSA+AES三段式,Redis,横向,纪律 |
| 122 | SSRF 内网探测的两个陷阱：服务端自身卡顿伪装成"目标超时"，以及把公网可见服务误报成"内网突破" | info-leak |  | template | Y | SSRF,内网探测,超时误判,对照组,方法学,攻击面核对,报告口径 |
| 123 | 登录"成功"的真伪三重证伪：字段校验 / 鉴权中间件 / 浏览器注入会话+清 Cookie 对照 | auth-bypass |  | template | Y | 口令复用,假成功,会话验证,三重对照,计分纪律,凭据复用边界 |
| 124 | 路由与端点发现的四个高频误判：OPTIONS 不能证明路由存在、404 必须在三态对照下读、白名单形态特征、稳定性基线 | methodology |  | template | Y | 路由发现,OPTIONS,三态对照,404vs401,白名单识别,稳定性基线,收口判据 |
| 125 | cloud-mail 认证机制与探测要点（JWT + KV 双校验 / /init secret 预言机 / KV 键名规律） | auth-bypass | cloud-mail (maillab/cloud-mail) | script |  | cloud-mail,cloudflare-workers,jwt,kv,认证机制,预判否定 |
| 126 | cloudflare_temp_email 认证机制与攻击面（四套认证头 / address JWT 需查库 / admin_login 免认证爆破入口 / DISABLE_ADMIN_PASSWORD_CHECK） | auth-bypass | cloudflare_temp_email (dreamhunter2333) | script | Y | cloudflare_temp_email,cloudflare-workers,临时邮箱,认证机制,x-admin-auth,admin_login,爆破入口 |
| 127 | 前端 chunk 硬编码凭据 → 共享钥匙 → 密钥库路径枚举批量导出 | unauthorized | Cloudflare Workers 密钥库服务（keyvalue / local-apikeys 命名） | script | Y | 前端硬编码,凭据泄露,密钥库,Cloudflare Workers,路径枚举,未授权访问 |
| 128 | Vite dev server 暴露（Astro 生产误跑 dev）—— /@fs/ + ?raw 读取项目源码 | info-leak | Astro / Vite dev server | poc | Y | vite,dev-server,astro,源码泄露,任意文件读,/@fs,?raw,robots内部路径 |
| 129 | 1Panel v2 安全入口启用时的可打面与判定手册（匿名白名单路由 + 验证码是否为唯一前置） | info-leak | 1Panel | script | Y | 1panel,v2,安全入口,验证码,匿名路由,CVE-2025-66507,CVE-2025-66508,端口9527 |
| 130 | 1Panel 口令审计前置判据：必须先有 panel_public_key 与（安全入口启用时的）entrance，且无用户名 oracle | weak-password | 1Panel | script | Y | 1panel,口令审计,前置条件,panel_public_key,RSA三段式,安全入口,ddddocr,验证码,用户名oracle |
| 131 | 1Panel v2 登录验证码自动化攻破（ddddocr）+ ErrAuth/ErrEntrance 命中判据 + “有安全入口即不可行”结构性结论 | weak-password | 1Panel | script | Y | 1Panel,v2,验证码,ddddocr,OCR,ErrEntrance,安全入口,爆破可行性,端口9527 |
| 132 | 开放代理“穿内网”真伪判定法：内网 vs 公网同 SNI 对照 body hash（附 CONNECT+自定义SNI 手法） | tunnel | Squid / nginx SNI 路由 | script | Y | squid,开放代理,内网突破,boundary,SNI,vhost,对照实验,判据,CONNECT |
| 133 | 1Panel 面板登录验证码绕过（ignoreCaptcha=true → ErrAuth） | auth-bypass | 1Panel | poc | Y | 1panel,验证码绕过,登录,弱口令,明文口令 |
| 134 | 端口定性三态判据：中间盒黑洞 vs 真 filtered vs 真服务（同机对照 + 12 载荷 + SNI 矩阵） | info-leak | 端口定性 | script | Y | 端口定性,黑洞,filtered,对照实验,ECONNRESET,TLS,SNI,nmap,FOFA,8443 |
| 135 | Flask before_request 鉴权跳过 OPTIONS ⇒ 未授权枚举后台路由表与方法（含变量规则噪声判据） | info-leak | Flask | script | Y | flask,werkzeug,OPTIONS,before_request,未授权,路由枚举,Allow头,变量规则,噪声 |
| 136 | New API 支付回调/初始化免鉴权面判定清单（epay notify + stripe/creem/waffo webhook + /api/setup） | unauthorized |  | script | Y | new-api,epay,stripe,webhook,setup,未授权,支付回调,验签 |
| 137 | Discuz! <target> RELEASE 20101020 后台 Getshell 全路径（离线白盒审计手册） | rce | Discuz 7.2 后台 | exp |  | discuz,<target>,getshell,后台,into outfile,模板编译,uc_key,config.inc.php,sql执行 |
| 138 | 单请求型隧道下的内网渗透方法论（op=chk 主机发现 + 目标侧 PHP 绕过多轮协议限制） | tunnel | PHP WebShell 自建 SOCKS5 隧道（op=req 单请求模式） | script | Y | 内网渗透,隧道,WebShell,op=chk,NFS,RPC,open_basedir绕过,主机发现,静默失败,eval马,冰蝎 |
| 139 | phpwind <target> 站群通用探测判据集（登录接口响应语义 + 额度机制实测 + 未授权接口 + 目录 deny 清单） | weak-password | phpwind 8.7 | template | Y | phpwind,<target>,判据,登录接口,额度,锁定,XFF,regcheck,未授权,目录deny,上传补丁,站群通用 |
| 140 | phpwind <target> 补丁面速查表（<target> 白盒产出）：已修复 9 类 + 未验证缺陷 6 条 + 未穷尽面 4 个（全附源码行号） | app-audit | phpwind 8.7 (含 v8.7/20110913 与 8.7beta 同源分支) | template | Y | phpwind,<target>,补丁面,白盒审计,速查表,GBK,宽字节,addslashes,复用,同圈子,源码行号 |
| 141 | 博彩平台会话网关 /ssid1?url= 反射型 XSS + 开放重定向（可窃取 ssid1 会话 cookie，会员/代理/管理三端命中） | auth-bypass | 博彩系统出租平台 会话网关 /ssid1 | poc | Y | XSS,开放重定向,会话网关,ssid1,博彩平台,会话劫持,三端,跨站复用 |
| 142 | phpwind <target> 存储型 XSS 选点判据：checkVar 豁免名单 + prosign 零转义写入 + apartment 校验绕过 + groupid=7 拦截真相 + 签名展示位对照 | auth-bypass | phpwind 8.7 | poc | Y | phpwind,<target>,存储型XSS,prosign,checkVar豁免,apartment绕过,groupid7,展示位,站群通用 |
| 143 | phpwind <target> 注册三必填绕过（regreason + _hexie + apartment 地区ID）→ 注册成功且自动登录 | auth-bypass | phpwind 8.7（GBK） | poc | Y |  |
| 144 | phpwind <target>「decode 错位链」白盒审计结论：无需账号的裸引号→convert(/e) 路径不存在；CMS 投稿链裸引号入库→展示 RCE（代码级成立）；forum 发帖链被 addslashes+escapeChar 双保险挡住 | rce | phpwind 8.7 (require/bbscode.php [code]/e 链) | template |  | phpwind,<target>,decode错位,裸引号,code-e,convert,CMS投稿,白盒审计,负结论 |
| 145 | 入驻方资料接口未授权 IDOR（shop_user_id 任意遍历泄露商家微信/公众号） | unauthorized | 自研 uni-app 付费内容平台后端（/api/home_info 模块） | poc | Y | idor,unauthorized,business-logic,pii,wechat |
| 146 | 域名清单接口未授权泄露全平台 20 个生产域名 + 后端回源 IP | info-leak | 自研 uni-app 付费内容平台后端（/api/user 模块 + 后端健康检查） | poc | Y | info-leak,unauthorized,recon,subdomain,backend-ip |
| 147 | agent_id 鉴权分支绕过 —— 免登录读取付费内容平台正文（goodsHistory/goodsList/goodsDetail） | auth-bypass | ThinkPHP 5 自建「圣誉科技」付费内容平台 | poc | Y | auth-bypass,idor,未授权,付费内容,agent_id,thinkphp,thinkphp5 |
| 148 | SCLUB 平台 <target> 注入客户端修正（繁体歡迎+hex左截断+FROM包层+限流特征） | frontend-reverse |  | script | Y |  |
| 149 | SCLUB 平台 <target> 注入回显客户端（含三个易踩根因与语句约束表） | auth-bypass |  | script | Y | sql注入,union,回显截断,繁体字符,old-password,sclub,客户分库站群 |
| 150 | acg-faka（异次元发卡）免登录用户 API 访问控制模型与判定基准（secret/query/commodity/card） | unauthorized | acg-faka 异次元发卡 | template | Y | acg-faka,异次元发卡,发卡网,未授权,越权,列过滤盲注,源码判定 |
| 151 | Dujiao-Next Admin 后台指纹与 /api/v1/admin/** 管理 API 全清单提取（含未授权定性探测） | info-leak | Dujiao-Next Admin (dujiao-next, Go+Vue3) | template | Y | dujiao-next,独角数卡NEXT,发卡网,后台,API枚举,未授权,Vue3,源站直连 |
| 152 | 若依系前台登录「请求加密守卫」绕过（isEncrypt/AES-ECB + RSA encrypt-key，前端硬编码密钥对） | auth-bypass | RuoYi-Vue-Plus / 若依系（dromara）前后端分离后台 | script | Y |  |
| 153 | Dujiao-Next 版本指纹 + 上游补丁差分法（app_version 免认证泄露 → 缺失安全修复清单 → 逐条可达性实测） | info-leak | Dujiao-Next（独角数卡 Go 重写版） | script | Y | dujiao-next,独角数卡,发卡网,版本指纹,信息泄露,补丁差分,未授权 |
| 154 | agent 后台忘记密码流程逆向：未鉴权改密接口 getBackLoginPassword 完整字段契约与差分探测法 | auth-bypass | agent backoffice (/agent/api/v1) | script | Y | agent,backoffice,password-reset,unauth,getBackLoginPassword,forgetPassword,twoSt |
| 155 | X-API-XXX WAF 签名头逆向：路径绑定规则 + 验签缺失（任意十六进制即可穿透） | env | Next.js SPA + Moleculer 网关 | script | Y | waf,signature,x-api-xxx,nextjs,moleculer,replay,bypass |
| 156 | Next.js SSR pageProps.t 密文可解密 —— 前端硬编码 AES 密钥 V/GO`(h$-?93a8jg | frontend-reverse |  | poc | Y | 前端加密,SSR加密,硬编码密钥,信息泄露 |
| 157 | 主站登录 X-API-FINGER 指纹算法还原 + 6022 真实成因纠正（FINGER 非决定因素） | auth-bypass |  | poc |  | 登录绕过,风控指纹,前端逆向,FingerprintJS,murmur3 |
| 158 | 代理后台 loginFlow 两步验证流程完全逆向（isValidate/checkCode/XCode/XCodeId），卡点精确定位为极验票据 | auth-bypass |  | poc | Y | 登录流程,极验,两步验证,代理后台,错误码差异 |
| 159 | glaxy 网关 (_glaxy_tus5r4_ 命名空间) API 探测与签名复现 —— 含 10050 证伪与 8 证书多租户识别 | auth-bypass | glaxy 网关（Spring Boot + OpenResty/Tengine 1.15.1） | script | Y | python,springboot,签名网关,命名空间,地域封锁,信息收集,反证zabbix,多租户证书 |
| 160 | glaxy 业务网关签名复现（E=字符降序）+ 匿名未授权读取后端配置/密钥/内网IP | auth-bypass | glaxy 业务网关 (tus5r4) | exp | Y | 签名绕过,字符降序,未授权访问,glaxy网关,appid泄露,越权 |
| 161 | club-upload 未授权任意图片上传 + AES-ECB 响应解密（密钥 7EGWY67N31RFK7GY）+ FastDFS 落地回读 | unauthorized | <target> /club/upload | poc | Y | unauthorized-upload,file-storage,fastdfs,aes-ecb,frontend-key-recovery,<target>, |
| 162 | 纯 Python RTMP 握手与 chunk 层实现（修掉 message header 顺序与 S2 残留两个坑） | protocol | SRS (Simple Realtime Server) 5.0.210 | script | Y | rtmp,amf0,srs,protocol,fuzzing,harness,auth-probe |
| 163 | phpwind <target>（<target>:1852）补丁面黑盒实测矩阵：13 类已封堵 + 5 个门槛不可越 + write_config 落点存在但入口不可达 + 目录 deny 与限流判据 | recon | phpwind v8.7 (20111111) GBK + nginx | template | Y | phpwind,<target>,补丁面,黑盒实测,负结论,站群复用,write_config,速率限流,判据 |
| 164 | phpwind <target> download.php type=message 零鉴权越权读任意附件（游客可读） | unauthorized | phpwind 8.7 | poc | Y | phpwind,附件下载,越权,未授权访问,type参数,getDownloadFactory |
| 165 | <target> 字符串数组还原 + 前端 AES 密钥派生还原（Vite/Vue SPA 通用） | frontend-reverse | Vite/Vue SPA + <target> + CryptoJS | script | Y | <target>,JS混淆还原,RC4,字符串数组,AES,CryptoJS,前端加密,配置解密 |
| 166 | 响应/请求加解密完整还原（金沙904.com 系 SaaS 博彩门户前端） | frontend-reverse | 自研 lobby 前端 (Vite + CryptoJS) | script | Y |  |
| 167 | Vite SPA 多子应用 API 面递归提取（子应用发现 + manifest chunk 清单 + API 正则） | info-leak | Vite / Vue SPA 多子应用架构 | script | Y | Vite,SPA,子应用,manifest,API提取,前端逆向,攻击面发现 |
| 168 | 绕过前端响应加密 - 读取 Vue/Pinia 运行时 store 获取明文数据 | info-leak | Vue 3 + Pinia SPA（前端加密响应） | script | Y | 前端加密绕过,Vue,Pinia,运行时状态,响应解密,SPA,CDP |
| 169 | Vue Router 路由全量枚举 - 获取 SPA 完整页面地图 | info-leak | Vue Router SPA | script | Y | Vue Router,SPA,路由枚举,页面测绘,前端逆向,CDP |
| 170 | Java Servlet 未授权数据读取探测（POST + userid 参数 + 404/500/200 判别） | unauthorized | Java Servlet / Tomcat (LHService 系列) | poc | Y | java,servlet,tomcat,未授权访问,参数爆破,userid |
| 171 | SPA 管理后台源码未授权下载 + action 全量未授权测试（含子目录路径发现） | auth-bypass | IMService 管理后台 SPA (ECDSA P-256 设备签名协议) | poc | Y | spa,管理后台,源码泄露,未授权访问,action枚举,ecdsa |
| 172 | 菠菜站群 aocai-api 未授权读取加密站点配置（AES-ECB 密钥硬编码于前端 app.js） | unauthorized | 菠菜站群 prodSite5 / aocai-api（Tomcat） | poc | Y | 菠菜站群,prodSite5,aocai-api,Tomcat,未授权访问,AES,前端加密破解,信息泄露 |
| 173 | aocai 站群 API 网关未授权访问（userGuestRegister 签发 token + getWebConfigNew 配置泄露） | unauthorized | aocai 站群 API 网关（Spring Boot + MongoDB，前端 Vue/Flutter） | poc | Y | aocai,站群,未授权,spring-boot,mongodb,token签发,AES-ECB |
| 174 | 菠菜站群 prodSite5 未授权文件上传（attachment/upload + uploadChat，落地 Azure Blob 公网可读） | file-upload | 菠菜站群 prodSite5 / aocai-api（Spring + Tomcat 10.1.50） | poc | Y | 菠菜站群,prodSite5,aocai-api,未授权文件上传,Azure Blob,文件写入 |
| 175 | 前端混淆/Dart AOT 双路还原 WebSocket 上行协议（含心跳与可达性判定） | frontend-reverse | 菠菜站群 /chatWs WebSocket（Flutter aocai_love_v3 + Vue wap 前端） | script | Y | websocket,混淆还原,flutter,dart-aot,协议逆向,<target>,心跳 |
| 176 | aocai 站群前端双 AES 密钥完全复现（webConfig.key 派生，离线解密响应） | info-leak | aocai 站群前端（Vue2 + <target> + CryptoJS） | exp | Y | aocai,站群,前端加密,AES,ECB,密钥派生,离线解密,<target>,CryptoJS |
| 177 | WordPress XML-RPC system.multicall 批量口令爆破（含判定特征校准） | weak-password | WordPress XML-RPC (system.multicall) | script | Y | wordpress,xmlrpc,multicall,brute-force,旅 |
| 178 | WordPress REST API 未授权信息收割（用户枚举 + 全量页面媒体 + 版本指纹） | info-leak | WordPress REST API | script | Y | wordpress,rest-api,未授权,用户枚举,信息泄露 |
| 179 | 优客服 UKefu 路径 ACL 绕过（/chat/ 前缀放行 + sign 仅长度校验） | auth-bypass | 优客服 UKefu (Java/Spring Boot/Thymeleaf) + 自研 SWG/2.0 网关 | poc | Y | ukefu,优客服,java,spring,403绕过,访问控制,sign,路径黑名单 |
| 180 | Sunflower 2.0 博彩系统 API 未授权访问（含 TLS 拒绝绕过法） | unauthorized | Sunflower 2.0（自研白标体育博彩前端 SPA + .NET 后端 API + gocache 网关） | poc | Y | sunflower,白标博彩,未授权访问,api,gocache,tls绕过,ip直连,dotnet |
| 181 | 白标真人娱乐场 H5 客户端加密/签名方案完整还原（AES-128-CBC IV=key + zlib + HmacSHA1） | frontend-reverse | 自研白标真人娱乐场 H5 + 自研 Go 游戏网关（game-http） | script | Y | 白标博彩,加密逆向,AES,HmacSHA1,game-http,网关,token,内联bundle |
| 182 | 尾斜杠触发 Traefik 301 清理泄露隐藏路由 BASE（通用绕过手法） | info-leak | 自研加固平台（Go + React SPA + nginx + Traefik） | script | Y | traefik,nginx,路径清理,隐藏路由,访问控制绕过,尾斜杠 |
| 183 | 执行式去混淆：伪装运行环境跑混淆脚本，直接读回 window 配置（通用手法） | info-leak | 6686 博彩集团前端（rc4 变体混淆 + Spring Boot 网关 + Cloudflare） | script | Y | 前端混淆,去混淆,rc4,字符串表,node,配置泄露,6686 |
| 184 | Jackson 未知字段报错泄露内部类名与 DTO 字段清单（Spring Boot 接口测绘手法） | info-leak | Java Spring Boot + Jackson（com.iglottery.admin） | poc | Y | jackson,spring-boot,信息泄露,dto测绘,接口结构,未鉴权,6686 |
| 185 | 拼图滑块验证绕过：clientNonce 明文内含滑块正确坐标且可离线伪造 | auth-bypass | 自研博彩平台拼图滑块验证 + Spring Boot 网关 | poc | Y | 滑块验证码,拼图,clientNonce,验证码绕过,签名,博彩,6686,自研网关 |
| 186 | gameplat-cms /api/img/upload 未授权任意文件上传（Content-Type 可控 → 存储型 XSS） | file-upload | gameplat-cms / CN7 模板（博彩游戏平台 CMS） | poc | Y | 未授权,文件上传,对象存储,存储型XSS,博彩CMS |
| 187 | CN7 平台 /jumpto/img/&lt;url&gt; 未授权开放图片代理 SSRF | ssrf | CN7 CMS 3.1.5 博彩平台（多租户） | poc | Y | ssrf,图片代理,CN7,未授权 |
| 188 | OpenSSH regreSSHion CVE-2024-6387 崩溃验证（7etsuo PoC 落盘 + Python 复刻 v1/v2） | methodology | OpenSSH sshd (signal handler race) | exp |  | openssh,sshd,regresshion,race-condition,DoS |
| 189 | Vite + javascript-obfuscator SPA：还原 chunk 映射表 / 字符串表反混淆 / 服务层接口全量提取 | frontend-reverse | Vite + javascript-obfuscator SPA | script | Y | vite,webpack,chunk,javascript-obfuscator,反混淆,接口提取,SPA |
| 190 | OCS 在线客服中台 访客会话越权读取（vid IDOR） | auth-bypass | OCS 在线客服中台（自研，Vue+SpringBoot） | poc | Y | idor,越权,在线客服,OCS,cid,x-v-token,snowflake,未授权 |
| 191 | 自研 Web 后台接管四步法：app.js 提 API 面 + 真实域名绕 Origin 白名单 + 无锁定弱口令 + 数据面遍历 | weak-password | Live2 直播监听后台（自研 Go + 原生 JS SPA） | exp | Y | 弱口令,origin绕过,sni,自研后台,app.js,信息泄露 |
| 192 | 站群 RN APP 分发链下载 + bundle 接口/密钥提取（域名轮换跳转 JS 型分发） | info-leak | React Native APP（code-push-server 站群客户端） | script | Y | react-native,apk,逆向,接口提取,硬编码密钥,下载链 |
| 193 | PHPWind 零鉴权信息泄露双接口：sitemap.php 泄露全站帖子标题+作者，member.php 泄露全部会员名 | info-leak | PHPWind 7.3.2 | poc | Y | phpwind,信息泄露,sitemap,member,零鉴权,用户名枚举,站群 |
| 194 | PHPWind v4/5/6/7/8 hack/bank/index.php 银行排行榜 bank_sort.php 代码注入 RCE（官方称 <target>+ 才修复 → <target> 受影响） | rce | PHPWind | exp |  |  |
| 195 | PHPWind 7.5 apps/groups/share/index.php $route/$basePath 未初始化任意文件包含（前台零登录、可远程包含 getshell） | rce | PHPWind | exp |  |  |
| 196 | PHPWind 7.5 pw_api.php → api/class_base.php $mode 未过滤任意文件包含（来自 $_POST+$_GET，不依赖 register_globals） | rce | PHPWind | poc |  |  |
| 197 | PHPWind 7.x/8.x 后台 getshell 手法全谱（SQL INTO OUTFILE / 模板包导入 / 备份恢复 / cache.php+level.php eval 点 / 附件扩展名）+ 后台入口定位 | rce | PHPWind | template |  |  |
| 198 | PHPWind <target> 负面情报清单（register_globals 已移除 / mawhole eval 证否 / setstyles $ 转义 / hack.php 白名单锁死 / api 全 404）+ 唯一命中的 nuclei 模板 | app-audit | PHPWind | template |  |  |
| 199 | PHPWind <target> 负面情报：member.php 的 orderway/asc 排序参数未进入 SQL（纯回显） | app-audit | PHPWind 7.3.2 | template | Y | phpwind,<target>,负面情报,orderway,member.php,注入,判据,假阳性 |
| 200 | PHPWind 7.x hack/bank savearray → data/bbscache/bank_sort.php 用户名注入（写码+同请求include） | rce | PHPWind 7.3.2 / 7.5 SP3 / 8.7 | poc |  | phpwind,savearray,变量插值,写码执行,bbscache,<target>未修 |
| 201 | PHPWind global.php InitGP(array('u','a')) 装载 $a 但被无条件 unset —— ${@eval($a)} 取不到值的判据 | rce | PHPWind 7.5 SP3 global.php | poc |  | phpwind,InitGP,register_globals,变量注入,死payload判据 |
| 202 | PHPWind 7.x /eis 注入点：copyctrl('$read[colour]') —— 需 pw_threads.colour 列存在 | rce | PHPWind 7.x require/template.php + lib/staticpage.class.php | poc |  | phpwind,preg_replace,/e,复制控制,copyctrl,静态化 |
| 203 | PHPWind 7.x 入参清洗边界与 5 个 POST 豁免字段（payload 构造必须遵守） | app-audit | PHPWind 7.x global.php / require/common.php | script |  | phpwind,参数清洗,CheckVar,Char_cv,payload构造约束 |
| 204 | PHPWind uploadicon 头像上传在 7.5 SP3 的校验实现（<target> 需实测是否有白名单） | file-upload | PHPWind job.php / lib/upload.class.php | poc |  | phpwind,uploadicon,头像上传,任意文件写入,7.5已加校验 |
| 205 | PHPWind 7.x 后台入口机制：$admin_file = PHP_SELF，全站唯一入口 admin.php，install.php 存在则后台拒绝进入 | app-audit | PHPWind 7.x admin.php / admin/admincp.php | script |  | phpwind,后台入口,admin_file,PHP_SELF,adminjob,入口爆破 |
| 206 | PHPWind passport_client.php 通行证接口：verify=md5(action+userdb+forward+db_pptkey)，可改写任意用户密码并下发登录 Cookie | auth-bypass | PHPWind passport_client.php（通行证同步接口） | poc |  | phpwind,passport,通行证,认证绕过,账号接管,未授权面 |
| 207 | PHPWind 7.5 SP3 全站 /e 与代码求值点穷举（30 处 /e，仅 2 处可注入，均为 $read[colour]） | rce | PHPWind 7.5 SP3 全站 preg_replace /e 修饰符清单 | script |  | phpwind,preg_replace,/e,eis,代码求值,全量清单,<target>未修 |
| 208 | colour 判死与 7.5 等价字段 titlefont 证据链（/eis 注入必须的 ' 与 ) 都被前台清洗吃掉） | rce | PHPWind pw_threads.colour vs titlefont（7.5 SP3） | script |  | phpwind,colour,titlefont,/eis,判死,历史遗留字段,<target>未修 |
| 209 | PHPWind 7.5 SP3 上传扩展名变形的黑名单缺口（pht/phar/inc 未覆盖） | file-upload | PHPWind lib/upload/*.class.php 扩展名变形 | poc |  | phpwind,附件上传,扩展名变形,pht,后台白名单 |
| 210 | PHPWind mode.php / apps.php 路由审计：$m 与 $q 都不进路径，两侧均无文件包含（含 mode/o 真实文件清单） | app-audit | PHPWind mode.php / apps.php / require/app_route.php 路由审计 | script |  | phpwind,mode.php,apps.php,路由审计,selectMode,Pcv,文件包含,判死 |
| 211 | 照片上传（mode.php?m=o&q=photos&a=upload）判死：扩展名白名单硬编码 + 文件名全服务端生成 + GetImgSize 内容校验，.user.ini 不可行 | file-upload | PHPWind lib/upload/photoupload.class.php（相册照片上传） | poc |  | phpwind,照片上传,photoupload,扩展名白名单,user.ini,判死 |
| 212 | PHP /e 三类替换串（形态 A 死 / 形态 B 语法可行 / 形态 C 真危险）+ Char_cv 不处理括号与 CheckVar 白名单的实测印证 | rce | PHP 5.x /e 修饰符引擎语义（php-src 官方源码铁证） | script | Y | php,/e,preg_replace,php_addslashes,preg_do_eval,引擎级判据,否证 |
| 213 | PHPWind 7.5 SP3 全站 30 处 /e 逐条判定（含 atc_content 最强候选链、prosign 与 titlefont 链路判死） | rce | PHPWind 7.5 SP3 全部 /e 调用点逐条判定 + prosign/titlefont 链路 | script |  | phpwind,/e,逐条判定,atc_content,prosign,titlefont,CheckVar,判死 |
| 214 | PHPWind 7.5 SP3 全部 /e 替换串原文与形态判定：30 处无一为形态 B（参数位逃逸），全死 | rce | PHPWind 7.5 SP3 全部 /e 替换串原文与形态判定 | script |  | phpwind,/e,形态A,形态B,替换串原文,逐条判定,判死,参数位逃逸 |
| 215 | PHPWind <target> 实测判据集：prosign 写入姿势+引号实体化、admincheck 无自助入口、mode.php?m=o 各 q 可用性、载体长度硬上限 | recon | PHPWind 7.3.2 | template | Y | phpwind,<target>,prosign,admincheck,mode.php,判据,负面情报,站群复用 |
| 216 | PHPWind 7.5 SP3 反序列化判死：全站 0 个魔术方法（无链起点）+ 无 cookie/请求参数来源 | deserialization | PHPWind 7.5 SP3 unserialize 全量 + POP chain 可用性 | script |  | phpwind,反序列化,unserialize,POP chain,魔术方法,判死,版本错位 |
| 217 | PHPWind 7.5 SP3 cookie 登录态：StrCode 无 MAC 可逆，但伪造受 db_hash 封死（拿到 db_hash 即账号接管） | auth-bypass | PHPWind 7.5 SP3 cookie 登录态（StrCode / PwdCode / User_info） | script |  | phpwind,cookie,StrCode,PwdCode,User_info,db_hash,登录态,账号接管 |
| 218 | 站群同源判定：多 vhost 是否指向同一 DOCUMENT_ROOT（body diff 相似度 + 尺寸 = 4×主机名长度差） | recon | 多 vhost 站群（PHPWind/nginx/宝塔） | template | Y | 站群,同源判定,body,diff,DOCUMENT_ROOT,Cookie前缀,判据,负面情报 |
| 219 | 端口开放判定判据：代理 CONNECT 假成功陷阱 + 以应用层数据为准（宝塔 Pure-FTPd / nginx 403 实测） | recon | 宝塔主机端口复验（Pure-FTPd / nginx phpMyAdmin 端口） | template | Y | 端口复验,CONNECT假成功,应用层判据,宝塔,Pure-FTPd,phpMyAdmin,IP白名单,免费代理不可用,换出口 |
| 220 | 宝塔 Pure-FTPd 口令面判据：存在/不存在用户同一条 530（无法枚举），只有 230 算通；字典来源与 PUT 马战术 | weak-password | 宝塔面板 Pure-FTPd（MySQL 虚拟用户后端） | script | Y | 宝塔,Pure-FTPd,MySQL虚拟用户,530,用户名枚举不可行,230判据,FTP字典,上传webshell |
| 221 | EyouCMS 前台任意文件写入 GetShell（api/Uploadify::preview 后缀绕过） | file-upload | EyouCMS | exp |  | eyoucms,文件上传,getshell,前台未授权,Uploadify |
| 222 | EyouCMS 后台模板编辑 RCE（CVE-2026-7388，ThinkPHP {:} 标签绕过过滤） | rce | EyouCMS | exp |  | eyoucms,后台getshell,模板写马,代码注入,thinkphp,标签绕过,FilemanagerLogic |
| 224 | EyouCMS 前台 pre-auth SQL 注入（common.php GetSortData，参数 sort_asc） | sqli | EyouCMS | poc |  | eyoucms,sqli,前台未授权,pre-auth,GetSortData,sort_asc |
| 225 | EyouCMS 后台 Logo 上传不受限 → GetShell（CVE-2024-11211） | file-upload | EyouCMS | exp |  | eyoucms,后台getshell,文件上传,不受限上传,Logo |
| 226 | EyouCMS 未授权任意文件读 / 文件包含致代码执行（CVE-2026-1107，Diyajax check_userinfo） | unauthorized | EyouCMS | poc |  | eyoucms,文件包含,任意文件读,路径穿越,未授权,Diyajax,check_userinfo |
| 227 | EyouCMS 二次反序列化 RCE（CVE-2025-15375，api/Ajax.php arcpagelist 参数 attstr） | deserialization | EyouCMS | exp |  | eyoucms,反序列化,php对象注入,thinkphp,gadget链,arcpagelist |
| 228 | EyouCMS 后台接管跳板：官方 setpwd.php 密码/密保重置工具（配合任意文件写） | auth-bypass | EyouCMS | script |  | eyoucms,后台接管,密码重置,任意文件写,跳板,setpwd |
| 229 | trade.php 判死 + 全 NUL 成因（防护清零）+ 反推前人用的是 bank 链（真正该查 data/bbscache/bank_sort.php） | app-audit | PHPWind trade.php 与全 NUL 文件成因 | script |  | phpwind,trade.php,全NUL,查杀清零,写码点,反推前人payload,bank_sort |
| 230 | PHPWind 7.5/<target> 根目录标准文件清单对照：trade.php 是标准，moneylog.php 与 userinfo.php 是非标准 | app-audit | PHPWind 7.5/8.7 标准文件清单对照（识别非标准文件） | script |  | phpwind,标准文件清单,非标准文件,moneylog,userinfo,版本识别 |
| 231 | EasyPay 支付回调签名走私（return_url 内嵌 trade_status 伪造支付成功） | auth-bypass | EasyPay / 聚合支付回调（Go） | poc |  | payment,webhook,signature,parameter-smuggling,md5,go,business-logic |
| 232 | EyouCMS 后台任意文件复制 → 泄露数据库配置（CVE-2026-6561，Index::edit_adminlogo） | info-leak | EyouCMS | poc |  | eyoucms,后台,任意文件复制,配置泄露,数据库口令,edit_adminlogo |
| 233 | EyouCMS 版本判定速查（4 判据 / 最多 4 包锁定版本区间） | recon | EyouCMS | script |  | eyoucms,版本判定,指纹,踩点,省包,install.lock |
| 234 | EyouCMS 低价值 / 已修复条目汇总（反面清单，勿耗包数） | app-audit | EyouCMS | template |  | eyoucms,低价值,已修复,反面清单,xss,csrf,信息泄露 |
| 235 | <target> 站群跳转机全机地图与负面情报（唯一动态程序 PHPWind <target> + 全静态页 + 三套异机真入口） | recon | 站群地图与负面情报（宝塔域名池/跳转机 · 赌博类站群） | template | Y | 站群地图,负面情报,宝塔,跳转机,域名池,每日域名生成器,线路选择页,泛解析,PHPWind<target>,判死 |
| 236 | 易支付回调伪造可行性判定（签名域是否含攻击者可控值）——new-api 反面 / sub2api 正面 | auth-bypass | 易支付/EasyPay 兼容网关回调（go-epay 系） | poc |  | payment,webhook,signature,parameter-smuggling,easypay,new-api,sub2api,triage |
| 237 | 得分系统 G1 去重陷阱：同档位重复记分会被累加，禁止为补 step_id 关联而重记得分 | methodology | 得分系统/台账纪律 | template | Y | 得分系统,G1,去重陷阱,台账纪律,负面情报,禁止重记,step关联 |
| 238 | IDOR 判定铁律：必须做"同身份换 ID"对照实验，否则不算确认（含一次真实误报复盘） | auth-bypass | IDOR / 越权判定纪律 | template | Y | IDOR,越权,判定纪律,误报,对照实验,负面情报,phpwind |
| 239 | PHPWind 7.x 后台入口名定位：读 pw_config.db_adminfile（前台 $db_adminfile 渲染），不做文件名爆破 | info-leak | PHPWind 7.x | script |  | phpwind,后台入口,配置项,信息泄露,白盒审计 |
| 240 | Sub2API EasyPay 支付回调验签绕过（参数走私 / 签名串拼接不一致 → 伪造回调入账） | auth-bypass | Sub2API（EasyPay / 易支付） | exp | Y | 易支付,easypay,sub2api,验签绕过,参数走私,业务逻辑漏洞,伪造入账 |
| 241 | Sub2API 易支付(EasyPay) 回调验签绕过 —— 通用参数化利用脚本（零密钥伪造回调入账） | auth-bypass | Sub2API (Wei-Shaw/sub2api) + 易支付 EasyPay | exp | Y | sub2api,easypay,易支付,支付回调,验签绕过,参数走私,签名拼接不一致,零密钥,业务逻辑 |
| 242 | PHPWind hack 插件入口 hack.php?H_name=&lt;插件&gt; —— 绕过 hack/&lt;插件&gt;/index.php 的 9B Forbidden 守卫 | frontend-reverse | PHPWind 7.3.2 hack 插件框架 | poc | Y | phpwind,hack插件,入口发现,readover绕过,9B守卫 |
| 243 | PHPWind 7.x pw_ajax.php?action=pcdelimg fieldname 布尔盲注（wooyun-2010-0417，影响 phpwind 7/8，7.5 SP3 已修） | sqli | PHPWind | poc |  | phpwind,<target>,sql注入,盲注,pcdelimg,pw_ajax,后台入口,db_adminfile,apikey |
| 244 | PHPWind 7.x pw_api.php → api/class_Other.php threadscateGory 写码 RCE（wooyun-2010-0417 第二洞，7.5 SP3/8.x 已加 intval 修复） | rce | PHPWind | poc |  | phpwind,<target>,rce,写码,pw_api,class_Other,threadscateGory,info_class,sig,apikey |
| 245 | PHPWind 7.x pw_api.php → api/class_base.php $mode 任意文件包含（CN-VA10-02 点名文件，不依赖 register_globals） | rce | PHPWind | poc |  | phpwind,<target>,任意文件包含,LFI,pw_api,class_base,mode,CN-VA10-02,不依赖register_global |
| 246 | PHPWind <target> 情报审计负面清单 + 7.5已修/<target>未修补丁对照（tag写码/hack_rate LFI/faq/vote/apps 全部源码级证否） | app-audit | PHPWind | poc |  | phpwind,<target>,负面情报,证否,pw_var_export,hack/rate,faq,vote,apps,补丁差异,7.5已修,action |
| 247 | PHPWind <target> 生日显示写码链（6.3.2 前台 getshell 手法）成立性判定：机制成立但被注册名 Char_cv 引号实体化阻断 | rce | PHPWind | poc |  | phpwind,<target>,生日缓存,brith_cache,写码,前台getshell,Char_cv,引号实体化,负面情报,6.3.2手法 |
| 248 | PHPWind 7.x 探测方法论：CC 防护 pwDefendCc 导致的「0 字节假阴性」与四条规避规范 | recon | PHPWind 探测方法论 | script | Y | phpwind,方法论,CC防护,假阴性,探测规范,盲注,布尔判定 |
| 249 | PHPWind 7.x 入口 × action 全量枚举方法论（251 个 action + 10 类危险标记）+ 写点判死结论集（pw_var_export 转义 / Pcv 封穿越 / uploadicon 数学不可能 / mutiupload 三重锁） | path-traversal | PHPWind 7.x（7.3.2 / 7.5 SP3） | script | Y | phpwind,action枚举,攻击面测绘,写点判死,pw_var_export,Pcv,P_unlink,uploadicon,mutiupload,路径穿 |
| 250 | PHPWind 7.x 后台风格 CSS 编辑 → data/style/*_css.htm → updatecache_sy eval 代码注入 RCE | rce | PHPWind 7.x（admin/setstyles.php + admin/cache.php） | poc |  | phpwind,后台,风格模板,代码注入,eval,${},stripslashes抵消,data/style |
| 251 | PHPWind 写 PHP 文件 sink 的三面叠加扫描法（数据来源分类 + 变量路径补扫） | rce | PHPWind 7.3.2 / 7.5SP3 | script | Y | phpwind,白盒审计,sink扫描,写文件,数据来源分类,变量回溯 |
| 252 | PHPWind 写 .php 缓存文件判死边界速查表（含 bank_sort 例外与字段长度表） | app-audit | PHPWind 7.3.2 / 7.5SP3 | payload | Y | phpwind,白盒审计,判死边界,缓存写点,bank_sort,字段长度 |
| 253 | PHPWind <target> admin/bakup.php 任意 SQL 执行（bakindata 多语句 + bakuptable 表名注入） | sqli | PHPWind | poc |  | phpwind,bakup,任意SQL,多语句,备份恢复,权限绕过判据 |
| 254 | PHPWind <target> 系统性加固判死清单（bank_sort RCE 链 + pw_hack 写入面 + 权限模型全量审计） | rce | PHPWind | poc |  | phpwind,<target>,bank_sort,rce,判死清单,权限模型,pw_hack,复用价值 |
| 255 | GBK 站点后台表单「全量回放」必须走 latin-1 字节级往返（PHPWind <target> 实测） | app-audit | PHPWind 7.3.2 | script | Y | phpwind,gbk,utf-8,form-replay,encoding,latin-1,config-restore |
| 256 | PHPWind <target> 后台 Getshell 优先级手册（权限自检 148xx 判据 + P0 SQL OUTFILE → P5 附件扩展名 + 源码级判死清单） | rce | PHPWind 7.3.2 后台 getshell 优先级手册 | template |  | phpwind,<target>,后台getshell,操作手册,权限判据,判死清单,宝塔 |
| 257 | PHPWind <target> 后台附件白名单放开（GBK 全量表单回放）与上传失败分层判据 | info-leak | PHPWind 7.3.2 后台附件设置 / 宝塔 WAF | script | Y | phpwind,后台,附件白名单,全量表单回放,GBK,latin-1,attachdir,宝塔WAF,判据 |
| 258 | 出口 IP 被宝塔 WAF 拉黑后的应急出口：免费 HTTP 代理池（抓取/验活/轮换 + schannel 坑） | tunnel | 宝塔 WAF IP 黑名单绕过 / 免费 HTTP 代理池 | script | Y | 出口IP,宝塔WAF,IP黑名单,免费代理池,轮换,schannel,ssl-no-revoke,判据 |
| 259 | 宝塔网站防火墙 <target> 拦截边界与绕过（IP 层换出口 / URL 关键词混合大小写 / multipart 双拦 PHP 特征） | waf | 宝塔网站防火墙 9.8.9（BTWAF） | template | Y | 宝塔WAF,BTWAF,IP黑名单,关键词绕过,multipart,上传拦截,换出口IP,PHPWind,判据 |
| 260 | 出口 IP 与 Host 头是两个独立变量：环境代理静默改道导致宝塔 WAF 封禁被误读（curl 必须 --noproxy "*"） | waf | 宝塔 WAF 9.8.9 / 出口 IP 封禁 | template | Y | 宝塔WAF,出口IP,环境变量代理,noproxy,假阴性,判据,站群,PHPWind |
| 261 | PHPWind 7.x 前台/后台 POST 一律被 Showmsg('undefined_action') 拒绝 —— 根因是 PostCheck() 要求 Referer 的 host 严格等于 HTTP_HOST（用 IP 做 Referer 是致命陷阱） | app-audit | PHPWind 7.3.2 / 7.5（global.php PostCheck） | template | Y | phpwind,referer,undefined_action,PostCheck,POST失败根因,判据,GBK,IP直连陷阱 |
| 262 | 宝塔 WAF 来源封禁实测判据：multipart 上传后全出口 000/ProxyError，而同代理到 baidu 200（区分「代理失效」vs「目标拒我」） | waf |  | poc | Y |  |
| 263 | 宝塔网站防火墙 <target> 拦截页指纹（403/1200）＋ SQL 关键字载荷被整拦的实测判据＋代理出口两级验活法（--noproxy 自杀陷阱） | env | 宝塔网站防火墙 9.8.9 / 宝塔 nginx | template | Y | 宝塔WAF,<target>,SQL注入判据,403,1200,代理验活,--noproxy陷阱,myip,出口被封 |
| 264 | PHPWind <target>/7.5 后台写文件 getshell 三处写入端判死（Pcv 封穿越 / pw_var_export 转义数组键 / setcss 非写入端）+ data 目录"是否写入成功"的可靠判据 | path-traversal | PHPWind 7.3.2/7.5SP3 admin/setstyles.php + admincp.php PostLog + require/common.php Pcv/pw_var_export | template | Y | phpwind,后台getshell,判死,Pcv,路径穿越,setcss,PostLog,admin_record,pw_var_export,写文件判据,r |
| 265 | SPA+Spring Cloud 网关批量未授权接口探测（/xxa 前缀型） | unauthorized | Spring Cloud 多租户博彩 SaaS（/xxa 网关） | script | Y | 未授权访问,Spring Cloud,SPA,接口枚举,配置泄露 |
| 266 | 点云客服 dotclouds 访客端未授权 multipart 任意文件上传（net.dotclouds / Spring Boot） | file-upload | 点云客服系统（net.dotclouds，Udesk 私有化，Spring Boot） | poc | Y | 点云客服,dotclouds,udesk,未授权,文件上传,getshell |
| 267 | 点云客服 dotclouds 会话接口零鉴权读取（IDOR，queryChatHistory/queryConvMessage） | auth-bypass | 点云客服系统（net.dotclouds / Spring Boot） | poc | Y | 点云客服,dotclouds,未授权,IDOR,会话读取 |
| 268 | dysaas 平台家族统一上传接口任意内容写入（/xxa/longm/api/v1/file/upload） | file-upload | dysaas 博彩平台网关（/xxa · longm 文件服务） | poc | Y | 文件上传,扩展名绕过,对象存储,S3,任意文件写入,dysaas |
| 269 | dysaas 代理后台 API 前缀双写规律 + 匿名 clientCode 商户配置泄露 | unauthorized | dysaas 代理后台（代理后台 SPA · agentPlate） | poc | Y | 代理后台,未授权,路径前缀双写,信息泄露,dysaas,AES |
| 270 | 点云客服 dotclouds sendConvMessageFile：companyId 未 normalize 导致任意目录文件写入（路径穿越、可逃出 /data 到根文件系统） | file-upload | 点云客服系统（net.dotclouds，Udesk 私有化，Spring Boot） | exp | Y | 点云客服,dotclouds,udesk,未授权,文件上传,路径穿越,任意文件写入,spring-boot |
| 271 | dysaas 代理后台：CDN 边缘 IP 轮换绕过源站封禁 + 请求格式/AES + 匿名代理注册 appSave | unauthorized | dysaas 代理后台（agentPlate · 极星云 CDN 前端）+ 主站 /xxa 网关 | script | Y | dysaas,代理后台,CDN边缘IP轮换,绕过IP封禁,未授权,匿名代理注册,AES |
| 272 | 直播平台未授权注册绕过短信验证（/api/user/simple_regist 空 mobile_hash 直接签发 token） | auth-bypass | 龙门直播 / 多租户直播平台 API（Go gRPC-gateway） | exp | Y |  |
| 273 | 腾讯云 STS 临时凭证 + COS q-signature 手工签名（写入/边界判定）可用脚本 | unauthorized | 腾讯云 COS（STS 临时凭证 + q-signature） | script | Y |  |
| 274 | 多租户直播平台前端响应加密还原（AES-256-CBC 硬编码 key/IV + 时间戳尾部拼接） | info-leak | Vue SPA + Go gRPC-gateway（多租户直播平台「龙门直播」v10.8.21.9） | script | Y |  |
| 275 | 匿名注册/申请接口作账号存在性 oracle（用户名枚举通用手法 + dysaas 代理 SaaS 判定表） | info-leak | dysaas 代理 SaaS（Spring Boot + Spring Security） | script | Y | 用户枚举,存在性oracle,未授权,SpringBoot,CDN边缘直连 |
| 276 | /xxa/longm/api/v1/file/upload 未授权文件上传（multipart，无 Authorization） | file-upload | xxa longm file service (6T/dysaas 平台) | exp | Y | 未授权,文件上传,对象存储,S3,白名单绕过 |
| 277 | 密码喷洒：批量生成的运营账号共用默认口令（先证明登录无需验证码 + 按前端算法构造口令密文） | weak-password | dysaas 代理 SaaS 登录接口（Spring Boot，AES-CBC/Pkcs7 前端口令加密） | script | Y | 密码喷洒,弱口令,批量账号,验证码绕过,CDN边缘直连 |
| 278 | Swagger 驱动的水平越权批量探测（自己/他人 id 对照 + 未文档化参数注入 + 加密响应长度 oracle） | auth-bypass |  | poc | Y | idor,swagger,openapi,越权,加密响应,oracle |
| 279 | 多租户 SaaS 后台「已证否方向」坑位卡（上传无解析 / detail 非 IDOR / Bearer 前缀与单会话顶号 / 验证码漏校验 / 错误码语义表） | methodology | dysaas 代理 SaaS / 多租户后台（Spring Boot + Spring Security + 对象存储） | template | Y | 负向结论,坑位卡,多租户,SpringBoot,对象存储,单会话,Bearer |
| 280 | RN/Expo(Hermes) APK 逆向：提取后端签名密钥与加密算法 | info-leak | React Native / Expo (Hermes bytecode) | script | Y | android,apk,react-native,hermes,逆向,签名密钥,硬编码密钥 |
| 281 | CN7 /jumpto/img SSRF 内网/云元数据不可达（边缘代理阻断，实测负面） | ssrf | CN7 CMS 3.1.5 博彩平台（多租户） | poc | Y | ssrf,CN7,图片代理,内网不可达,负面 |
| 282 | 博彩站群(<target>:10666) Java 反序列化/JNDI/Log4j RCE 全否 — 严格 Jackson 非 Fastjson + /ecb 加密网关可达真实后端 | deserialization | Spring Boot + Jackson (若依风格 code+message+data) / OpenResty+AES-OFB ecb 网关 | poc | Y | fastjson,Jackson,反序列化,JNDI,Log4j2,actuator,ecb,若依,SpringBoot,负面结论 |
| 283 | CN7/gameplat 10666 管理组件面（actuator/druid/swagger/若依）全否 — SPA 兜底 + Spring 全局鉴权过滤器 + 版本泄露 5.3.0-RELEASE | info-leak | CN7 CMS / gameplat 5.3.0-RELEASE (Spring Boot + nginx + OpenResty /ecb 网关) | poc | Y |  |
| 284 | CN7/gameplat 5.3.0-RELEASE 无公开版本级 Nday + 同产品实例全为静态 CDN 克隆 + nginx ^/api 前缀路由测绘 — 版本指纹 Nday 匹配判死清单 | info-leak | CN7 CMS / gameplat 5.3.0-RELEASE (Spring Boot + nginx + OpenResty /ecb 网关) | poc | Y | CN7,gameplat,博彩CMS,版本指纹,Nday,nignx路由,多租户 |
| 285 | 博彩 CMS「构建哈希」指纹：body="&lt;buildhash&gt;" 在 FOFA 一键枚举同产品同构建部署舰队（含每租户一端口划分） | info-leak | CN7 CMS / gameplat（前端 webpack/Next.js 静态 SPA + Spring Boot） | poc | Y | 指纹,FOFA,构建哈希,博彩CMS,资产测绘,CN7,gameplat |
| 286 | 博彩站群 <target>:10666 Java 反序列化/JNDI/Log4j2 全否（含 OOB 有效性复核与 /ecb 网关格式） | deserialization | CN7/gameplat CMS 5.3.0-RELEASE (Spring Boot + 严格 Jackson + Logback) | markdown | Y | negative,deserialization,jndi,log4j,fastjson,jackson,spring-boot,ecb,waf-bypass, |
| 287 | CN7/gameplat 多租户SaaS：未授权 /api/tenant/setting/getPersonalCenter 跨租户泄露（含跨租户运营账号） | info-leak | CN7 / gameplat CMS (Spring Boot, 5.3.0-RELEASE) | poc | Y | 多租户,cross-tenant,unauth,info-leak,CN7,gameplat,SaaS |
| 288 | CN7/gameplat 多租户博彩 SaaS 未授权跨租户配置泄露 (getPersonalCenter/getConfig) | info-leak | CN7 gameplat 5.3.0-RELEASE (Spring Boot + nginx/OpenResty) | http | Y | multi-tenant,broken-tenant-isolation,unauthorized,CN7,gameplat,info-leak |
| 289 | CN7/gameplat /api/config/promote/* 是推广跳转(catch-all 302)而非受门控的配置管理组——无 host/租户绕过面 | auth-bypass | CN7 gameplat 5.3.0-RELEASE (Spring Boot + nginx) | poc | Y |  |
| 290 | Web 缓存投毒/缓存欺骗面一站式测绘脚本（含缓存键归一化与 unkeyed 头判定） | recon | CDNGreat / 通用 HTTP 缓存层 | script | Y | cache-poisoning,web-cache-deception,host-header,smuggling,http-proxy,cache-key |
| 291 | CDNGreat 缓存/防护层指纹识别与缓存范围、缓存键判定方法 | recon | CDNGreat 缓存层 | script | Y | cdngreat,cache,fingerprint,waf,host-header,cl-cache-status |
| 292 | fakerToken 服务无鉴权探测：/fbsport/v1/getToken 与 /ob/v2/getTryToken 泄漏后端域名与 token | info-leak | fakerToken (白标体育供应商 guest token 服务) | script | Y | fakerToken,guest token,未授权,信息泄露,requestData,form编码,线路域名路由表,白标供应商,体育平台,ob体育 |
| 293 | 对象存储未授权列举/写入判定（含 SPA 回落假阳性排除 + 写探针回滚） | unauthorized | AmazonS3 / 对象存储静态托管 | script | Y | 对象存储,S3,OSS,未授权,桶列举,ListBucketResult,写探针,假阳性排除,SPA回落,静态托管 |
| 294 | Vite/Vue3 minify SPA 接口调用点与请求体还原器（三层别名映射） | frontend-reverse | Vite 打包 SPA + Vue3（minify 混淆） | script | Y | Vite,Vue3,minify,混淆,别名映射,三层映射,接口还原,请求体还原,前端逆向,SPA,静态分析 |
| 295 | 用户可控字段 → 无鉴权公开接口渲染 → 面向全站访客的投毒链（四步验证 + 如实边界） | unauthorized | 通用链：用户可控字段写入接口 + 未鉴权公开渲染接口（本次为 RB88 博彩站 SPA） | script | Y | 投毒链,头像外链,未授权渲染,公开榜单,内容注入,XSS判定,清洗器,safe-url,博彩SPA |
| 296 | 前端硬编码 requestData 常量重放 → 未授权签发第三方平台 token 与线路域名路由表 | info-leak | 白标体育供应商 guest token 服务 + 前端硬编码 requestData 常量 | script | Y | 硬编码常量,requestData,重放,guest token,未授权,fakerToken,线路域名路由表,白标供应商,体育平台 |
| 297 | fakerToken 硬编码 requestData 常量重放：未认证签发 OB 平台登录 token 与线路后端域名表（RB88/OB 站群通用） | unauthorized | fakerToken（RB88/OB 体育站群内部 token 签发服务，Spring Boot） | poc | Y | fakerToken,rb88,博彩站群,未授权,token签发,硬编码常量,requestData,信息泄露 |
| 298 | Vite SPA 资金接口静态还原 + 二维码渲染定性（含带会话的页面内 fetch 验证） | info-leak | Vite 打包 SPA / Spring Boot 资金后端 | script | Y |  |
| 299 | 统一 20000 未登录拦截的绕过穷尽判定法（含 body×CT 归因控制实验，防 116 误报） | auth-bypass |  | script | Y | auth-bypass,拦截器绕过,20000,Content-Type,白名单前缀,对照实验,防误报 |
| 300 | cmwall「过移动墙」服务管理面板 <target> 未授权读写（38888）+ 零风险写验证流程 | auth-bypass | cmwall「过移动墙」服务管理面板 v1.35 | poc | Y | cmwall,过移动墙,38888,未授权,面板,流量劫持,红队 |
| 301 | 前端资源投毒可行性判定法：客户端渲染资金字段 + 无 SRI/无 CSP/资源域写方法（含正确分级口径） | frontend-reverse |  | script | Y | 前端投毒,SRI,CSP,资源域,CORS,二维码,资金链路,可行性判定 |
| 302 | 支付回调/通知端点探测法：静态零引用判定 + 未知路由三态判别 + 按站点命名风格派生字典（含负结论口径） | auth-bypass |  | script | Y | 支付回调,notify,伪造到账,路由三态判别,前端零引用,字典探测,负结论口径 |
| 303 | fakerToken/OB 的 api= 参数解密（硬编码 AES-ECB 密钥 OBTY20220712OBTY）与 requestData 不可解结论 | info-leak |  | script | Y | 硬编码密钥,AES-ECB,fakerToken,OB,api参数解密,离线,bundle审计 |
| 304 | cmwall「过移动墙」wizard 面板 静态分析速查 + 命令注入实现基础（含厂商公开包/部署路径/注入验证要点） | rce | cmwall「过移动墙」服务设置向导面板（wizard，Go 实现，UPX 打包） | script | Y | cmwall,<target>,wizard,过移动墙,命令注入,UPX,面板,红队 |
| 305 | WAF 绕过/能力修正：直连 HTTPS 清空 SNI（+任意 Host）→ 200；带真实 SNI（<target>）→ RST（按 SNI 丢包） | waf | nginx TLS vhost + 源 IP/SNI 型 WAF（<target> 站群） | script | Y | WAF绕过,SNI,直连,TLS,站群,API,能力修正 |
| 306 | 加固 APK（zxprotect 抽取壳）零依赖逆向 + config2.json 配置引导链提取 | info-leak | 加固 APK（zxprotect 抽取壳）+ 站群 App 配置引导链 | script | Y | android,apk,加固,抽取壳,zxprotect,dex切分,config引导链,C段指纹比对,信息收集 |
| 307 | LHZS/LHService v27 网关免签名面定位 + 两层签名校验拒绝码判据 | unauthorized | LHZS/LHService v27 API 网关（nginx 前置 + 双层签名校验，ECDSA P-256） | poc | Y | LHZS,网关,免签名,未授权,签名校验,拒绝码判据,routing-health,v27 |
| 308 | 原生 socket 端口普查器（绕代理假象）tcpscan.py | recon | 红队侦察基础设施 | poc | Y | 扫端口,原生socket,代理绕过,段普查,recon |
| 309 | LHZS/靶标节点签名探测器（严格判据）lhprobe.py | recon | 红队侦察基础设施 | poc | Y | 资产归属,严格判据,免签名接口,签名指纹,recon,LHZS |
| 310 | TLS 证书抓取 + 从 DER 提 CN/SAN（证书归属判定）certgrab.py | info-leak | TLS / 证书透明 | poc | Y | 证书,归属判定,SAN,DER,无代理,tls |
| 311 | 多目标业务路径 + Host(vhost) 覆盖探测器 sweep.py | unauthorized | 红队侦察基础设施 | poc | Y | vhost,Host头,业务路径,归属判定,recon |
| 312 | 非标准服务四态识别器 rawprobe.py | recon | 红队侦察基础设施 | poc | Y | 服务识别,黑洞端口,协议判定,recon |
| 313 | 端口 banner 抓取器（自动判 SSH/MySQL/HTTP）banner.py | recon | 红队侦察基础设施 | poc | Y | banner,服务识别,SSH,MySQL,版本指纹,recon |
| 314 | Web 指纹探测（标题 + LHZS 特征签名）httpfp.py | recon | 红队侦察基础设施 | poc | Y | Web指纹,标题,UTF8,中文判定,recon |
| 315 | 单目标业务路径 + vhost 探测（先看正文再判归属）bizprobe.py | unauthorized | 红队侦察基础设施 | poc | Y | 业务路径,vhost,正文判读,归属判定,recon |
| 316 | DNS 解析 + 指定 SNI/Host 的 HTTPS 定向探测 sni72.py | recon | TLS / vhost 路由 | poc | Y | SNI,DNS,归属判定,vhost,历史节点 |
| 317 | nginx→Tomcat `..;/` 归一化差异：穿 location 白名单 + 跨 context 横向（含绕开 nginx 死上游 504）；Manager IP 阀 33 变体实测无效 | path-traversal | nginx (location 白名单) + Apache Tomcat 8.5.78/8.5.81/9.0.122 多实例 | poc | Y | nginx,location绕过,路径归一化,..;,跨context,Tomcat,Manager,RemoteAddrValve,负结论 |
| 318 | Tomcat 拿 shell 路线的快速判死：PUT 405 + OPTIONS 无 PUT ⇒ CVE-2017-12615/12617/2025-24813 全灭；WEB-INF 25 种写法全 404；Tribes RCE 与 Ghostcat 版本不匹配 | rce | Apache Tomcat 8.5.x / 9.0.x | poc | Y | tomcat,负结论,readonly,PUT,405,CVE-2017-12615,CVE-2017-12617,CVE-2025-24813,CVE-202 |
| 319 | 全量参数钥匙：批量候选参数名越过老 servlet 的 Integer.parseInt 门槛（分组探测 + 二分收敛，通用可复用） | unauthorized | 老 JSP/Tomcat servlet（LHZS 系 2015 年代码） | script | Y | servlet,未授权访问,参数爆破,老代码,Tomcat,六合 |
| 320 | Tomcat Manager 快速判死：403 页字节级指纹（同版本 stock 页跨机器 sha1 完全一致）+ 端口面普查（只暴露 nginx ⇒ 头伪造必然无效，Tribes/8009 一并判死） | recon | Apache Tomcat Manager（RemoteAddrValve / host-manager） | poc | Y | tomcat,manager,RemoteAddrValve,负结论,stock403指纹,端口面,CVE-2026-34486,Ghostcat,头伪造无效 |
| 321 | 盲 SSRF 判定法：差分三值（死地址基线/内网/外带）+ A/B/A 交替复测识破抖动 —— 附"零差异=无 SSRF"实测样例与真 SSRF 对照特征 | ssrf | 老式 Java Servlet + 境外链路（高抖动）环境下的 SSRF 判定 | script | Y | ssrf,盲ssrf,判定法,差分,A/B/A复测,网络抖动,负结论,老servlet |
| 322 | 六助(LHZS) App <target> 客户端接口字段名与请求体结构全表 | protocol | 六助 LHZS App (com.ms6.lhzs / com.marksix) | payload |  |  |
| 323 | LHZS v27 管理台：认证链精确字段 + admin/query 信封 + release.publish SSRF + storage.configure 落盘可控 + action 全表 | auth-bypass | LHZS v27 管理台（SecurityClient / ECDSA P-256 + LHZSV2 canonical） | payload |  |  |
| 324 | 【技巧】Java 异常行号 oracle：用 Tomcat 错误页的 Class.java:NN 行号推进，零猜测定位 servlet 参数与鉴权门槛 | methodology |  | script | Y | java,servlet,参数发现,异常行号,oracle,judge,鉴权判定 |
| 325 | 【技巧】「参数钥匙串」批量命中老 Java servlet 数据分支：390 个候选参数各 =1，一次请求绕过猜参数名（单请求拿到 143KB/408 条用户记录） | methodology |  | script | Y | java,servlet,参数钥匙串,NumberFormatException,未认证,批量参数,老系统 |
| 326 | 【漏洞】nginx→Tomcat 归一化穿越（..;/）：绕开 nginx location 白名单，未认证直达 Tomcat ROOT/docs/examples 并确认 Manager/host-manager 存在（IP 阀纯头伪造绕不过） | path-traversal |  | poc | Y | nginx,tomcat,path-traversal,归一化,manager,RemoteAddrValve,绕过 |
| 327 | 原生 socket / 手工 TLS 探测骨架（绕本机 HTTP 代理污染 + binary_form 证书取证 + 定制404/Spring404/406 三态判存） | env | 自研「Json接口」Spring Boot 框架 / nginx 前置 | poc | Y | 原生socket,手写HTTP,绕本机代理,SNI=IP,证书取证,指纹判存,chunked |
| 328 | 老式 Java servlet「分块灌参 + 栈行号二分收敛」参数发现法（Tomcat 报错栈当 oracle） | info-leak | 老式 Java Servlet（Tomcat 默认错误页 + Integer.parseInt 型参数） | script | Y | java,servlet,参数爆破,tomcat,栈泄露,delta,二分 |
| 329 | UpLoadPic 三个入口对照表（upload2=唯一真写点/写 OSS；无扩展名分支=只回 URL 模板不落盘；upload 的 head 四字段=真字段但空响应、token 模式=无条件假成功） | file-upload | LHZS UpLoadPic 应用（com.lhzs.Service.UploadPicServlet，Tomcat 8.5.78/9.0.122 混部） | poc | Y | uploadpic,upload2,未认证文件上传,OSS,假成功,head字段,负结论,判据,LHZS |
| 330 | 【方法·负面知识】四步判死「无条件假成功」写端点 + 参数级 delta 证明：UploadPicServlet 只读 userid 一个参数、无文件名/路径可控、无文件落地 | methodology |  | script | Y | false-positive,假成功,写端点判定,UploadPicServlet,jar,参数级证明,delta,negative-knowledge |
| 331 | 行号 oracle 参数隔离法（附 PutTp 判死实测）：用 Tomcat 错误页的 `类名.java:行号` 逐级隔离参数、二分求最小集、判 POST 是否独立分支 | methodology | 老式 Java Servlet（Tomcat 默认错误页泄露行号） | poc | Y | 老servlet,Tomcat,错误页,行号oracle,参数隔离,最小参数集,doPost包装,负结论,PutTp |
| 332 | 【技巧】把老 servlet 入参面压到 1 个：最小必需参数集二分 + 行号 oracle 定位 parseInt 参数 + 用「-操作过于频繁」当分支存在判据（实测 Login4 收敛到 flag 一个参数） | auth-bypass |  | script | Y | java,servlet,flag分支,最小参数集,二分,行号oracle,值域枚举,限流判据,登录入口 |
| 333 | 【方法论·五步判据】老 Java servlet 读写判别：参数要补齐（防假否定）+ 读回要跨时间窗（防缓存假否定）+ 值域枚举看文案推进而非限流触发（防假肯定）+ 探针换假号（防锁桶） | methodology |  | script | Y | java,servlet,判据,假否定,假肯定,Redis缓存,写后读回,限流,值域枚举,flag分支,老系统 |
| 334 | 【oracle】判定「验证码到底有没有真的发出去」：用错码探「用码」分支，「请先获取验证码」vs「验证码错误」两句文案即判据（零成本、不消耗真码、无需目标手机） | unauthorized |  | script | Y | 短信验证码,发码判定,oracle,状态机,存根,假成功,参数名对照,限流保护 |
| 335 | 复现冲突四步法：纯读稳定性 → 裸/戴钥匙对照 → 逐个剔除定最小集 → 跨时间窗读回（写侧+读侧都做对才不误判） | methodology | 老式 Java Servlet（写类接口）+ 只读读回接口（带 Redis 缓存） | poc | Y | 未认证写,复现冲突,最小可复现请求,参数钥匙,Redis缓存,写后读回,判据口径,假成功 |
| 336 | jadx 卡死时的替代：最小 DEX 解析器 + 字符串表邻接/簇分析（无需反编译） | frontend-reverse |  | script | Y | android,apk,dex,静态分析,jadx替代,字符串邻接,不反编译 |
| 337 | 如何证明"我的会话确实生效了"：同一端点有无会话的响应差异自证法（400 动作级 vs 403 路径级） | methodology | lz2.apk / v27 后端（/IMService/v27） | script | Y | 会话自证,假阴性排除,方法论,黑盒验证,错误码分级 |
| 338 | Go 自研代理服务（管理口 + 代理口）免凭据探测与绕过穷尽核查脚本 | unauthorized | Go 自研代理/会话代理服务 (session-proxy / session-proxy-admin) | script | Y | go,代理面板,basic认证,open-proxy,未授权,指纹,穷尽测试 |
| 339 | FastAdmin CVE-2024-7928 任意文件读取（/ajax/lang lang 参数） | path-traversal | FastAdmin | poc | Y | FastAdmin,ThinkPHP,任意文件读取,path-traversal,CVE-2024-7928 |
| 340 | PHP .user.ini 两用技巧：追加 open_basedir=/ 破除目录限制 + auto_prepend_file 绕过 .php 写入黑名单落马（冰蝎 v4 双通道） | rce | PHP .user.ini (PHP_INI_PERDIR) / php-fpm / nginx 多站点 | script | Y | php,user.ini,auto_prepend_file,open_basedir,绕过黑名单,冰蝎,落马,宝塔 |
| 341 | PHP .user.ini 覆盖导致 open_basedir 逃逸（可写子目录 → 整机文件系统访问） | privesc | PHP 5.4 + php-fpm + .user.ini（user_ini.filename） | payload | Y | php,open_basedir,user.ini,逃逸,宝塔,phpwind,提权 |
| 342 | 宝塔防篡改只锁 .php + 替代后缀不被解析的实测判据（.php5/.phtml/.inc 换后缀拿不到 webshell） | app-audit | 宝塔网站防篡改 + nginx location 匹配 + php-fpm security.limit_extensions | template | Y | 宝塔,防篡改,php后缀,nginx,405,phtml,php5,webshell判死,phpwind |
| 343 | boyue(webman) 博彩站群会员端 —— upload-qrcode 任意扩展名写入 + 同源存储型 XSS（WithdrawController::uploadQrCode） | file-upload | boyue / webman 博彩站群会员端 | exp | Y | webman,boyue,博彩站群,任意文件上传,同源XSS,file-upload |
| 344 | boyue(webman+Pear Admin) 博彩站群后台指纹 / 路由表 / 密钥 / 免登录路由 速查手册 | info-leak | boyue（webman + Pear Admin 彩票/博彩管理系统） | script | Y | webman,boyue,pear-admin,博彩站群,gofly,指纹,路由表,免登录 |
| 345 | 本机 fake-ip DNS + 规则型本地代理下，用 CONNECT+自定义 SNI 直达 nginx vhost（绕过假性「目标封 IP」） | tunnel | 本地 HTTP 代理（Clash/v2ray 类）+ 目标 nginx vhost | script | Y | fake-ip,DNS劫持,SNI,CONNECT,代理,vhost,出口 |
| 346 | webman 插件路由指向非插件命名空间控制器 → AccessControl 中间件失效导致后台接口未授权 | unauthorized | webman + Pear Admin（PHP 常驻框架） | script | Y | webman,PearAdmin,未授权访问,中间件失效,路由命名空间,后台接口 |
| 347 | webman 路由大小写不敏感 → 绕过按路径控制器段精确匹配的「功能关闭」开关（可救回代码生成器写文件 → Monitor 热重载 RCE） | auth-bypass |  | script | Y | webman,workerman,路由,大小写绕过,WAF绕过,功能开关,代码生成器,file_put_contents,热重载,RCE |
| 348 | webman 默认(自动)路由绕过路由级 AuthMiddleware —— 未授权直达控制器方法（通用） | auth-bypass | webman (workerman/webman-framework) 默认自动路由 | poc | Y | webman,默认路由,自动路由,鉴权绕过,未授权访问,auth-bypass,php,workerman |
| 349 | webman-admin TableController::crud 代码生成器 docblock 注入 → 写 app/controller/*.php → Monitor 热重载 RCE | rce | webman-admin / Pear Admin（boyue 系博彩后台）TableController 代码生成器 | exp |  | webman,webman-admin,pear-admin,代码生成器,docblock注入,rce,php |
| 350 | webman/Pear Admin 目标「上传→执行」全判定手册（include/LFI/phar/模板 四路排除 + 扩展名名单测绘） | path-traversal | webman 2.1 / workerman 5.1.4 / Pear Admin 0.6.33 | script | Y | webman,workerman,pear-admin,php,code-audit,lfj,lfi,phar,upload,rce,oracle |
| 351 | webman 响应加密（X-Encrypted: 1）硬编码 AES 密钥解密 / 加密请求体构造 | frontend-reverse | webman 2.x + app\Utils\Crypto | script | Y | webman,响应加密,X-Encrypted,AES-256-CBC,硬编码密钥,固定IV |
| 352 | webman 自动路由绕过鉴权中间件（/{controller}/{action} 零凭证直连） | auth-bypass | webman（自动路由）+ 自研鉴权中间件 | script | Y | webman,自动路由,鉴权绕过,AuthMiddleware,未授权访问,RateLimit绕过 |
| 353 | webman 自动路由作为「方法存在性 oracle」：零凭证枚举部署版控制器方法（区分零鉴权直出 / 需登录 / 只判角色） | unauthorized | webman（PHP 常驻框架，含 Pear Admin 后台） | script | Y | webman,自动路由,未授权,方法枚举,oracle,鉴权测绘 |
| 354 | webman 自动路由未授权接口普查 + Crypto 硬编码密钥解密响应 | info-leak | webman | script | Y | webman,未授权,自动路由,信息泄露,AES硬编码 |
| 355 | nginx vhost oracle 枚举（Host 头区分真实虚拟主机，识别同名后端多品牌） | recon | nginx vhost | script | Y | nginx,vhost,信息收集,Host头,多品牌 |
| 356 | 跨 vhost 同源 XSS 可执行性判定矩阵（多域名站群：/uploads 可达性 × /app/admin 交叉对照） | unauthorized | webman 站群 / 多 vhost 共享 public + 会员端任意扩展名上传 | script | Y | 同源判定,vhost,多域名站群,存储型XSS,web root,跨域会话 |
| 357 | webman-admin(Vue Art Design Pro) 后台会话与请求加密规范 + 403 开关大小写绕过 + crud docblock 注入开火参数 | auth-bypass | webman-admin 定制后台（Vue3 Art Design Pro 前端 + plugin/admin 后端） | script | Y | webman-admin,ArtDesignPro,pinia,localStorage,AES-256-CBC,X-Encrypted,403开关绕过,doc |
| 358 | boyue/PearAdmin 充值渠道链上收款地址未授权替换 POC（payset-update / configs_data[trc20]） | unauthorized | boyue/博悦 彩票系统（webman + Pear Admin）充值渠道配置 | exp |  | webman,PearAdmin,未授权写入,充值渠道,链上地址,资金流,boyue |
| 359 | webman「显式路由 vs 自动路由」双姿势鉴权/闸门 oracle（可复用：一次判定整族方法被哪一层拦） | auth-bypass | webman 自动路由 / 交付版闸门 DeliveryMode | script | Y | webman,自动路由,鉴权绕过,交付版闸门,oracle,hch7777 |
| 360 | webman 自动路由「形参绑定」未授权身份/ID 注入（$request->all() → 控制器形参） | auth-bypass | webman (workerman/webman-framework v2.x) + workerman v5.x | poc | Y | webman,php,参数绑定,未授权,身份注入,IDOR |
| 361 | webman「形参名 oracle」：用 MissingInput/InputType 错误页在线反推控制器方法签名 | info-leak | webman (workerman/webman-framework v2.x, debug=true) | script | Y | webman,php,方法签名枚举,oracle,信息泄露,参数绑定 |
| 362 | SSRF 原语筛查法：从形参到外发 URL 的四步数据流溯源（含二阶配置写入者追踪） | ssrf | PHP / webman / 通用白盒筛查 | script | Y | ssrf,url参数,数据流溯源,白盒筛查,guzzle,二阶ssrf |
| 363 | 零副作用验证判据套件：身份注入 / 原密码强校验 / 验证码校验（用不存在 uid 与现值回写） | auth-bypass | 任意 Web 登录态系统（PHP/webman 类） | script | Y | 判据,零副作用,身份注入,验证码,改密,风控 |
| 364 | 路径穿越有效性判定：归一化等价对照法（nginx alias 前缀映射场景） | path-traversal | nginx-alias / %2f 编码穿越 | script | Y | nginx,alias,路径穿越,%2f,归一化,判定方法 |
| 365 | 任意后缀上传 → 同源存储型 XSS → localStorage JWT 会话窃取 → 无限速资金密码 oracle（账号接管/盗资金组合链） | file-upload | 好彩汇 H5 会员端（PHP + nginx 1.18.0，前端 Vue3 存 hch_token 于 localStorage） | exp | Y | 任意文件上传,存储型XSS,localStorage,JWT窃取,资金密码爆破,oracle无限试 |
| 366 | GoFly 客服系统鉴权模型与越权判定（/2/ 族 401vs403 判据 + ws_kefu/ws_visitor 子协议票据） | unauthorized | GoFly / goflylivechat（taoshihan1991）及二开版 | script | Y | gofly,goflylivechat,websocket,子协议鉴权,IDOR,客服系统,未授权 |
| 367 | 上游/供应链域名「可控性」判定流水线（whois+RDAP+DNS+HTTP+CT+TLS+家族挖掘） | recon | 通用 · 供应链域名接管面判定 | script | Y | 供应链,域名接管,whois,rdap,ct,recon |
| 368 | suo5 classic 隧道断链自愈 + PHP<target> 冰蝎马 openssl_decrypt("AES128") 返回 false 导致 200 空响应 | tunnel | suo5 v2.2.0 / PHP 5.4 / 冰蝎 v4 载荷 | script | Y | suo5,socks5,隧道,PHP<target>,冰蝎,openssl,运维 |
| 369 | ThinkPHP3 后台 Cookie 会话伪造（admin_user=username + admin_id=password哈希） | auth-bypass | ThinkPHP 3.1.3 后台 | poc | Y | thinkphp,会话伪造,cookie,auth-bypass,后台 |
| 370 | ThinkPHP3 后台 Api::article_save 任意文件写入（data:// 伪协议）→ WebShell | file-upload | ThinkPHP3 自研后台 Api::article_save | exp | Y | thinkphp,任意文件写入,data协议,getshell,rce |
| 371 | 冰蝎马 PHP 5.6 适配：openssl_decrypt("AES128") 失效须改 "aes-128-ecb" | rce | 冰蝎 Behinder PHP 马 × PHP 5.6 | payload | Y | 冰蝎,behinder,webshell,php5.6,openssl,AES128,aes-128-ecb,踩坑 |
| 372 | ThinkPHP <target> 站群 CMS：extends Action 的 Admin 控制器绕鉴权 + 微信签名可自算 → 未授权文件写入（Wechat::index_1） | auth-bypass | 六合彩/聚宝阁站群 CMS（ThinkPHP 3.1.3，署名 SomCNS「新开发kkjj」） | exp | Y | thinkphp,thinkphp<target>,鉴权绕过,extends Action,微信签名,任意文件写入,站点群,kindeditor,模板源码泄露 |
| 373 | PHP-FPM FastCGI 未授权访问 → PHP_ADMIN_VALUE 解除 open_basedir / 清空 disable_functions（纯 PHP 客户端，无需 shell） | info-leak | php-fpm (FastCGI) / 宝塔面板 / PHP 5.6-7.2 | exp | Y | php-fpm,FastCGI,未授权,open_basedir,disable_functions,宝塔,unix socket,横向通用,无shell |
| 374 | ThinkPHP <target> 站群 CMS：Home/Order::get_html1 零凭证覆盖模板文件（App/Tpl/Home/Index_index.html）→ 模板注入前置原语 | unauthorized | 六合彩/聚宝阁站群 CMS（ThinkPHP 3.1.3，署名 SomCNS「新开发kkjj」）— Home/OrderAction | poc | Y | thinkphp,thinkphp<target>,未授权,模板注入,任意文件写入,模板覆盖,get_html1,extends Action,站群 |
| 375 | ThinkPHP <target> 「未授权模板写入 → 编译执行」RCE 链（Home/Order::get_html1 → App/Tpl → App/Runtime/Cache） | rce | ThinkPHP 3.1.3 / 站群 CMS（Home/Order::get_html1） | exp | Y | thinkphp,<target>,template-injection,rce,unauth,站群,abrdns,php |
| 376 | Admin/Wechat 免鉴权签名绕过（token 空串伪造 sha1 签名） | auth-bypass | 站群 CMS Admin/WechatAction | poc | Y | wechat,签名绕过,未授权,站群,signature-bypass,thinkphp |
| 377 | ThinkPHP <target> 站群 CMS「模板烘焙」RCE 链完整打法 + 两个否证 oracle（get_html1 零凭证覆盖模板 / abrdns 抓取源陷阱 / Public/log 无回显失败判定） | rce | 六合彩/聚宝阁站群 CMS（ThinkPHP 3.1.3，署名 SomCNS「新开发kkjj」，jbg/thb/xmg 集群） | script | Y | thinkphp,thinkphp<target>,模板注入,模板烘焙,rce,get_html1,TMPL_DENY_PHP,parseKey,字段名注入,a |
| 378 | 冰蝎马 Python 自动化客户端（PHP <target> 需把 cipher 改成 aes-128-ecb） | tunnel | 冰蝎马 / Behinder AES-128-ECB | script | Y | behinder,冰蝎,AES-128-ECB,PHP5.6,webshell,自动化,客户端 |
| 379 | ThinkPHP 3：VAR_FILTERS=filter_exp 让 `exp` 数组注入全线失效（静默丢字段、无痕）——判定向量 + 无效绕过清单 | auth-bypass | ThinkPHP 3.1.3（VAR_FILTERS=filter_exp） | poc | Y | thinkphp,thinkphp3,VAR_FILTERS,filter_exp,exp数组注入,parseSet,静默丢弃,否证,判定向量,站群 |
| 380 | webman+nginx 站群：/app/admin/** 被 nginx 静态 location 遮蔽——路径/头/方法/差分编码/vhost 全不可绕（阴性） | unauthorized | nginx webman PearAdmin | markdown | Y | nginx,webman,vhost-admin,path-bypass,negative,static-location |
| 381 | 站群通用：article.ename = 创建者 admin 用户名 + Home/Index::bbs?post=&lt;id&gt; 零凭证渲染点（枚举后台账号名 / 内容可控渲染） | info-leak | 站群 ThinkPHP 3.1.3（天花板/聚宝阁/kkjjer 同套代码） | payload | Y | 站群,thinkphp<target>,零凭证,渲染点,账号枚举,article_save,get_html1 |
| 382 | UEditor PHP 零鉴权面盲探（config 泄露 / listimage 路径枚举 / upload 目录存在性 / catchimage SSRF） | unauthorized | UEditor (PHP) | poc | Y | ueditor,unauthorized,file-upload,mkdir,ssrf,thinkphp |
| 383 | ThinkPHP<target> 站群 CMS 通杀 RCE：admin 表哈希即凭据 → Home/Article::edit_cont 写 article.title → get_html1 烘焙模板 → 编译 <?php 执行 | rce | ThinkPHP 3.1.3 站群 CMS（署名 SomCNS） | poc | Y | ThinkPHP<target>,站群,RCE,模板烘焙,article.title,get_html1,open_basedir,Host头,admin哈希轮 |
| 384 | ThinkPHP<target> 站群通用链：后台 mass assignment 写 article.title → 零凭证 get_html1 模板烘焙 → 首页 RCE | rce | ThinkPHP 3.1.3 站群（天花板/大白菜/送钱宝宝/聚宝阁 同 CORE_DB 族） | exp | Y | thinkphp,站群,模板注入,模板烘焙,mass-assignment,get_html1,凭据复用,rce,全站群通用 |
| 385 | PHP-FPM unix socket 未授权 + PHP_ADMIN_VALUE 清空 open_basedir/disable_functions → 全盘任意文件读写（宝塔环境通用） | privesc | PHP-FPM (php-cgi unix socket) | exp | Y | php-fpm,fastcgi,open_basedir,宝塔,bypass,任意文件读写 |
| 386 | 对手检测能力画像：站群运营方 ThinkPHP3 后门自查脚本 jnbgg.sh（= 我们的规避指南） | recon | ThinkPHP 3.x 站点 / 黑帽站群运营方 | script | Y | 对手画像,检测规避,IOC,ThinkPHP3,站群,应急处置 |
| 387 | php-fpm FastCGI 沙箱逃逸：PHP_ADMIN_VALUE 清零 open_basedir 与 disable_functions（宝塔环境通用） | privesc | 宝塔(aapanel) + php-fpm (unix socket) + nginx | exp | Y | php-fpm,FastCGI,open_basedir绕过,disable_functions绕过,宝塔,aapanel,提权,沙箱逃逸,通用武器 |
| 388 | 站群兄弟节点发现：从目标自己的 vhost 访问日志 Referer 反查（日志分块读 + 一条请求判定同站群） | recon | 站群 / 多品牌黑帽集群（nginx + 宝塔） | script | Y | 站群,信息收集,Referer,vhost日志,兄弟节点,DNS,指纹 |
| 389 | php-fpm FastCGI sendmail_path 劫持 → 绕过 disable_functions 拿命令执行（mail() 的 C 层 popen） | privesc | php-fpm / PHP 5.6 / mail() / disable_functions | exp | Y | php-fpm,fastcgi,disable_functions,mail,sendmail_path,rce,linux |
| 390 | 熊猫体育 TY SDK /yewu40 签名通道无防重放 + checkId 不校验（试玩 token 实测） | auth-bypass | 熊猫体育 TY SDK (<target> / <target>) | template | Y | TY SDK,熊猫体育,白标体育,签名通道,防重放缺失,checkId,试玩token,<target>,ue3238 |
| 391 | VerneMQ-over-WebSocket 匿名消息总线快速评估（无鉴权 CONNECT/SUB/PUB + 直连要求 + 空闲判定） | unauthorized | VerneMQ (Erlang/Cowboy, MQTT-over-WS) | script | Y | mqtt,vernemq,websocket,anonymous,unauthorized,message-bus,recon |
| 392 | PHPWind 7.x/8.x 版本指纹识别（robots.txt + meta generator + 路径特征） | info-leak | PHPWind | poc | Y |  |
| 393 | FOFA 官方接口可用性切换（<target> 失效 → <target>） | recon | FOFA API | script | Y |  |
| 394 | 宝塔 WAF 对备份/源码文件扩展名的拦截：本站无「保留路径」的可用绕过（尾斜杠/%3f 只改路径，读不到文件） | info-leak | 宝塔网站防火墙 (BTWAF) | template | Y | 宝塔WAF,BTWAF,403-1200,备份文件,源码泄露,尾斜杠绕过,%3f,WAF绕过,db_hash,PHPWind |
| 395 | PHPWind 站点精确版本与数据库表前缀的无认证识别 | info-leak | PHPWind | poc | Y | phpwind,指纹,版本识别,表前缀,信息泄露 |
| 396 | 宝塔 per-source 封禁下的新出口：免费 HTTP 代理池可 CONNECT 任意端口(含 FTP:21)绕开源级拉黑 | tunnel | 宝塔/Pure-FTPd per-source IP 封禁 | poc | Y | 宝塔,per-source封禁,出口IP,免费代理池,FTP,CONNECT任意端口,绕WAF,新出口 |
| 397 | PHPWind <target> 后台 admin.php 锁定被 X-Forwarded-For 绕过（db_xforwardip=1 时）+ per-source CC 封禁才是爆破真正瓶颈 | weak-password | PHPWind 7.3.2 | poc | Y | phpwind,<target>,后台爆破,XFF,锁定绕过,per-source CC,db_xforwardip |
| 398 | 宝塔 WAF 按真实源IP CC 封禁下，用免费HTTP代理池逐请求轮换出口放大后台口令喷洒吞吐（含判据陷阱） | weak-password | 宝塔网站防火墙(BTWAF) 9.8.9 + PHPWind 后台 admin.php | script | Y | 宝塔WAF,CC封禁,出口IP,免费代理池,口令喷洒,XForwardedFor,PHPWind,admin.php,判据陷阱 |
| 399 | WordPress/WooCommerce 匿名攻击面清单（版本泄露+REST未授权+Store API非授权写入+XML-RPC爆破） | info-leak | WordPress / WooCommerce | script | Y | wordpress,woocommerce,unauthorized,info-leak,brute-force,rest-api,xmlrpc |
| 400 | 博彩 PHPWind 站群实战结论集（<target>/24）—— 判死路线 + 8 条可复用判据 | env | 博彩站群（PHPWind 7.3.2 + SafeLine WAF + 宝塔） | script | Y | 博彩站群,PHPWind<target>,SafeLineWAF,宝塔,判死结论,判据集,站群分群,出口被封,lgt,GBK |
| 401 | SSH 批量作业「banner 抖动 vs 真防护」判据 —— 单发复验可救回被误判的整批目标 | rce | OpenSSH | script | Y | ssh,bruteforce,误判防护,fail2ban,抖动复验,paramiko |
| 402 | FastAPI/OpenAPI 未授权测绘 —— openapi.json 泄露全端点 + 按 security 标记自动分类实测 | info-leak | FastAPI / OpenAPI | script | Y | fastapi,openapi,unauthorized,info-leak,swagger |
| 403 | ThinkPHP5 /kj 开奖数据接口未授权读取 | unauthorized | ThinkPHP5 自研开奖API | poc | Y | thinkphp5,unauthorized,api,json,博彩 |
| 404 | 自研BBS 站群 JSON API 硬编码密钥 + 接口协议逆向 | unauthorized | 自研BBS (bbsapiext.php) | poc | Y | 自研cms,api,硬编码密钥,注册接口,站群 |
| 405 | PHPWind 后台登录判据纠正 —— 失败态=14540B+页内计数文案；失败也下发 cookie；必须带 cookie 复访验证 | auth-bypass | PHPWind | script | Y | phpwind,后台,登录判据,假阳性,gbk,尝试计数器 |
| 406 | 页面源码挖自动生成 vhost —— 从 HTML/JS 提取域名并批量解析（补字典爆破盲区） | info-leak | 通用（Web 页面源码） | script | Y | recon,vhost,页面源码,自动生成域名,subdomain,方法论 |
| 407 | MinIO Console 两段式登录逆向 + Access Key 用户名枚举 | info-leak | MinIO Console | script | Y | minio,console,用户枚举,用户名枚举,认证,两段式登录 |
| 408 | TLS 证书 SAN 反查域名（纯 IP 段挖真实业务域名 + vhost 路由） | info-leak | 通用手法 / TLS 证书 | script | Y | 信息收集,证书,域名挖掘,vhost,TLS,SAN,纯IP目标 |
| 409 | AJ-Captcha 验证码协议复现 + secretKey 明文下发利用（芋道/RuoYi 系列） | auth-bypass | AJ-Captcha (anji-plus) / 芋道 yudao | poc | Y | AJ-Captcha,anji-plus,验证码,芋道,yudao,RuoYi,secretKey,AES,协议逆向 |
| 410 | TLS 证书 SAN 反查域名 + 同 IP vhost 路由发现（纯 IP 段挖真实域名，避免错过 target） | info-leak | TLS 证书 / nginx vhost / 泛域名 | script | Y | 信息泄露,TLS证书,SAN,vhost发现,域名挖掘,纯IP段,内网服务发现,可复用 |
| 411 | 母舰域名编号枚举 —— 一次画出攻击者全部节点（56 节点/10 网段实测） | info-leak | DNS / 基础设施测绘 | script | Y | recon,dns,母舰,基础设施测绘,编号枚举,资产发现 |
| 412 | SPA 后台突破 —— 从前端 api.js 提取全接口清单 + 定向弱口令（实测 admin/admin123 拿到超管） | weak-password | Node.js Express 管理后台 | script | Y | spa,react,api提取,弱口令,未授权,后台突破 |
| 413 | Synology DSM 未授权侦察 —— 两端点泄露 API 清单 / 主机名 / 公网状态 / 全部已装套件 | info-leak | Synology DSM 7.x | poc | Y | synology,DSM,未授权,信息泄露,套件枚举,攻击面,群晖,侦察 |
| 414 | 方法论纠错：高速 SYN 扫描在 SYN-flood 防护设备上产生海量假开放端口（必须低速复核） | recon | nmap / 通用扫描方法论 | script | Y | nmap,扫描,假阳性,方法论,SYN-flood,群晖,端口扫描,复核 |
| 415 | 群晖 DSM 套件前端 bundle 未鉴权扒取 + API/协议逆向（含容器终端 WebSocket 协议） | info-leak | Synology DSM 7.x 套件前端资源 | script | Y | synology,DSM,群晖,套件,静态资源,未鉴权,JS逆向,API清单,WebSocket,攻击面 |
| 416 | 群晖 DSM code 120 判读纠错 —— 参数校验先于权限，不可误判为未授权 | methodology | Synology DSM 错误码 | script | Y | 群晖,DSM,错误码,判读,未授权,误报,参数校验,测试方法论 |
| 417 | webman 站群(无PHP-FPM) upload-qrcode 上传→RCE 判死：/uploads 的 .php 硬403且不可绕、替代后缀只回源码、全 vhost 无解释器 | file-upload | webman / workerman / nginx | poc | Y | webman,file-upload,rce,negative,nginx,no-fpm,php后缀 |
| 418 | webman「写入 Monitor 监听目录→热重载 RCE」判死：落盘路径不可控 + 插件 admin 命名空间自动路由不可达 + /app/admin 被 nginx 遮蔽（hch7777 站群） | rce | webman / workerman / Pear Admin | poc | Y | webman,workerman,monitor,热重载,file-write,自动路由,命名空间,nginx,negative |
| 419 | shopnn(<target>) mock 登录认证绕过：alipay-mock 按 phone 无口令接管任意已知手机号账号 | auth-bypass | shopnn (Node/Express 二手电商) | poc | Y | auth-bypass,account-takeover,mock-login,node,express,jwt,shopnn |
| 420 | shopnn (<target>) 可利用性判死集：mock 登录回显 MySQL 错误（=信息泄露）；img-proxy SSRF 无内网/写面绕过 | info-leak | shopnn (Node/Express + Vue3 电商) | poc | Y | info-leak,mysql-error,express,mock-login,ssrf,blocked_url,redirect,self-port-gua |
| 422 | GoFly 客服附件卡片 safeChatAssetPath(%2e%2e) 绕过 → 同源存储型 XSS 主动投递链 | file-upload | GoFly / goflylivechat 二开版 (客服系统) | script | Y | gofly,kefu,xss,stored-xss,file-upload,sanitizer-bypass,phishing,same-origin |
| 423 | SafeLine/宝塔站群 WAF 精确触发签名（union select 需带操作数、相邻 ../、sleep() 括号）与全站路径无豁免实测判据 | sqli | SafeLine-WAF（PHPWind 7.3.2 GBK 站群） | poc | Y | waf,safeline,宝塔,phpwind,sqli绕过,全站拦截,union select,路径穿越,代理池 |
| 424 | PHPDisk <target> adminlogin 账号枚举 + bind_user 无频控口令校验双通道 | auth-bypass | PHPDisk 6.8 | script | Y | phpdisk,auth-bypass,user-enum,brute-force,php |
| 425 | PHPDisk 6.x adminlogin 注入面三点连测 + 后台 gid=1 校验 | auth-bypass | PHPDisk 6.8 | script | Y | phpdisk,sqli-triage,dosafe,padding-bypass,auth,admincp |
| 426 | New API 定制构建「全关闭」实例：注册关/setup已初始化/SMTP缺省下的负面情报与速判清单 | info-leak | New API (QuantumNous/new-api) v0.12.15-ledger-balance-probe 定制构建 | poc | Y | new-api,one-api,信息泄露,负面情报,未授权面,注册关闭,setup假象,限速预算 |
| 427 | go-captcha 拼图行为验证码破解：Sobel 梯度模板匹配定位 + /api/captcha/verify 直接换 token（4/4 实测通过） | auth-bypass | go-captcha (drag/拼图行为验证码) | exp | Y | go-captcha,验证码绕过,行为验证码,模板匹配,sobel,new-api,modelsell |
| 428 | New API Token Quota 面板 /token-quota.config.json 配置泄露 + 令牌状态判别法（unavailable=存在但禁用） | info-leak | New API (modelsell 定制分支) / 自建 Token Quota 面板 | poc | Y | new-api,token-quota,config-leak,sk-token,验活判据,unavailable,modelsell |
| 429 | <target> 注册墙全链路程序化突破：mail.tm 临时邮箱 + GET /api/verification 真实收码 + go-captcha(register) 建号成功（新用户初始 quota=0） | auth-bypass | New API (modelsell 分支 modelsell-2026.08.24) | exp | Y | new-api,modelsell,register,mail.tm,临时邮箱,verification,go-captcha,quota,注册建号,<targ |
| 430 | Discuz <target> 后台 runquery/SQL 导入对非创始人 adminid=1 全线 403（isfounder 闸先于 checkpermission 生效，实测证否） | unauthorized | Discuz 7.2 admincp | poc | Y | discuz,<target>,admincp,runquery,isfounder,后台,SQL执行,负面情报,into-outfile |
| 431 | PHPWind 后台操作日志 data/bbscache/admin_record.php 零鉴权直读（含真实 adminjob 模块与 bakup 备份前缀） | info-leak | PHPWind 7.x (7.3.2/7.5) | poc | Y | phpwind,postlog,adminlog,info-leak,backup-prefix,unauthorized |
| 432 | New API 邮箱验证注册绕过（一次性邮箱 + 6 位十六进制验证码 + 无图形验证码） | auth-bypass | New API | poc | Y | new-api,注册绕过,邮箱验证,一次性邮箱,ondemand |
| 433 | PHPWind <target> 全量源码获取：<target> 官方下载站快照直取法（免 GitHub） | info-leak | PHPWind 7.3.2 | template | Y | phpwind,<target>,源码获取,wayback,archive,镜像,下载,升级包,offline-audit |
| 434 | melody-admin pageconfig tableSuffix 未授权布尔盲注（MySQL root，跨租户分表切换） | sqli | melody-admin / SpringBoot + MyBatis | poc | Y | sqli,blind,mysql,unauth,pageconfig,multi-tenant |
| 435 | 七牛云 Kodo AK/SK 接管验证与桶枚举（QBox 签名保留 base64 填充） | unauthorized | Qiniu Kodo | script | Y | qiniu,kodo,cloud,aksk,object-storage,info-leak |
| 436 | renren-fast 派生后台：SQL 复制角色绕过角色创建者校验 + 自建管理员规避 token IP 互踢 | auth-bypass | renren-fast / melody-admin | poc | Y |  |
| 437 | melody 会员面 pageconfig/sysparam 签名-only 未授权读不可提权（AG49 生产实例） | info-leak | melody 会员 API (Spring Boot + openresty) | poc | Y | melody,pageconfig,sysparam,unauthorized,negative,ag49 |
| 438 | 七牛 Kodo 泄露 AK/SK：全量枚举、精确计量与公开读桶匿名外带 | info-leak | Qiniu Kodo | script | Y | qiniu,kodo,object-storage,public-bucket,ak-sk,info-leak,cdn |
| 439 | AG49 会员端 Spring Boot 2.x 框架层 N-day 与反序列化全否（Shiro/Spring4Shell/CVE-2018-1273/Fastjson/actuator/env） | deserialization | Spring Boot 2.x + springfox 2.8.0-SNAPSHOT + openresty 203.0.113.29 | markdown | Y |  |
| 440 | melody-admin(renren-fast 派生) 调度器 mode 语义与 MySQL UDF 面判否 + UDF 正控方法 | rce | melody-admin / com.invech.kitty / MySQL 8.0 | script | Y |  |
| 441 | openresty 1.29 + Spring Boot：请求走私(CL.TE/TE.CL/TE.TE)与 h2 降级全族阴性（勿再重打） | auth-bypass | openresty/203.0.113.29 + Spring Boot | payload | Y |  |
| 442 | melody-admin 彩种模板同步 SSRF：回显只认 data 数组信封；60s 时延不能当内网 DROP 判据 | ssrf | melody-admin / renren 风格 Spring Boot 管理端 | poc | Y | melody-admin,ssrf,time-based,oracle,readback,hutool,内网探测 |
| 443 | AG49 IM 节点 <target>:2585 配置未授权读不可提权（管理面/写面/注入全阴性 + 登录滑块阻断 staff 口令攻击） | info-leak | melody / lark IM (Spring Boot + nginx/1.14.1, WXIM 会员端) | poc | Y | melody,lark,im,ag49,unauthorized,info-leak,negative,slider-captcha,2585 |
| 444 | X-Forwarded-For 等代理头伪造客户端 IP | auth-bypass | Java Servlet / Spring Boot（通用 Web） | poc | Y | ip-spoof,x-forwarded-for,proxy-header,auth-bypass,rate-limit,audit-log |
| 445 | MySQL root 空口令远程直连检测（只读取证脚本） | weak-password | MySQL 5.x/8.x | poc | Y | mysql,空口令,未授权,数据库,caching_sha2_password |
| 446 | Laravel 图片上传（客户端后缀白名单+雪花ID重命名+Flysystem path 参数）无法转 RCE，残留 SVG 同源 XSS 与 4MiB FPM DoS | file-upload | Laravel + PHP-FPM + openresty | poc | Y | file-upload,laravel,flysystem,svg-xss,php-fpm,dos,negative |
| 447 | sendSmsCode 响应回显短信验证码 + modifyPasswordBySms 免短信接管任意账号（epay/Upay 会员 API） | auth-bypass | epay 会员 API（PHP/Yii2 + 前端 RSA） | exp | Y |  |
| 448 | 短信找回密码链：sendSmsCode 响应明文回显 OTP + modifyPasswordBySms 未认证改密 → 任意会员账号接管 | auth-bypass | epay / Hyperf-style JSON API | exp | Y | ATO,sms,sendSmsCode,modifyPasswordBySms,user-enumeration,account-takeover |
| 449 | PHP session.upload_progress 会话文件注入 + LFI → 预认证 RCE（Chain A） | rce | PHP | exp |  | php,session,rfc1867,upload_progress,lfi,rce,预认证 |
| 450 | PHP var_destroy 二次析构 UAF → `/proc/self/mem` 覆写 internal function handler → RCE（Chain B） | rce | PHP | exp |  | php,uaf,var_destroy,unserialize,proc-self-mem,handler-overwrite,rce,提权 |
| 451 | PHP session serialize_handler 不一致 → 解码失败仍触发 __destruct → POP gadget RCE（Chain C） | rce | PHP | exp |  | php,session,serialize_handler,pop-gadget,destruct,php-fpm,rce |
| 452 | CVE-2026-109618 Apache mod_php：multipart 解析期 OOM → longjmp 跳过 init_executor → shutdown 读脏 EG() → RCE（Chain D，零依赖） | rce | PHP / Apache mod_php | exp |  | php,oom,longjmp,mod_php,apache,eg-state,零依赖,rce |
| 453 | CVE-2026-109619 nginx + php-fpm：脏 BG(unserialize) 状态被下一请求消费 → __destruct → RCE（Chain E） | rce | PHP / php-fpm | exp |  | php,fpm,bg-state,unserialize,session,use-after-free,rce |
| 454 | PHP 环境存在性判别器（跑 PHP2Shell 五链前的必备预检） | app-audit | PHP2Shell / session.upload_progress 链族 | script | Y | php,php-fpm,mod_php,sapi,fingerprint,triage,php2shell,preflight |
| 455 | 授权总控类面板默认口令 admin:admin（HTTP Basic）→ 全量授权码+机器绑定接管 | weak-password | shouquan 授权总控（Python http.server 自制授权管理后台） | exp | Y | basic-auth,默认口令,授权码,后台接管,python |
| 457 | PHP2Shell 通用利用框架（CVE-2026-109618/619）：session.upload_progress UAF → Heap OOB 泄露绕 ASLR → pDestructor=zif_system 预认证 RCE | rce | PHP 8.x（mod_php / php-fpm）/ session.upload_progress | script | Y | php, php2shell, cve-2026-109618, cve-2026-109619, session.upload_progress, uaf,  |