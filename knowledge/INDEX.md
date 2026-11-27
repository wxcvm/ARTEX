# 知识库索引

共 **2672** 条，其中已验证 **378** 条。数据源 `knowledge.db`（SQLite + FTS5，`knowledge_dir` 直接加载）。

生成：`python knowledge/tools/export_docs.py`；体检：`python knowledge/tools/kb_check.py`

## 分类分布

| category | 条数 |
|---|---|
| `kev` | 1805 |
| `nvd` | 195 |
| `info-leak` | 86 |
| `auth-bypass` | 74 |
| `absorbed` | 60 |
| `gray-market` | 59 |
| `unauthorized` | 57 |
| `rce` | 42 |
| `recon` | 32 |
| `methodology` | 23 |
| `java-web` | 22 |
| `file-upload` | 20 |
| `web-cms` | 19 |
| `weak-password` | 18 |
| `frontend-reverse` | 16 |
| `edge-device` | 16 |
| `app-audit` | 15 |
| `cloud-container` | 14 |
| `exchange-ad` | 14 |
| `env` | 11 |
| `linux-priv` | 11 |
| `waf` | 10 |
| `path-traversal` | 9 |
| `java-deser` | 9 |
| `ssrf` | 8 |
| `tunnel` | 8 |
| `sqli` | 6 |
| `deserialization` | 5 |
| `protocol` | 4 |
| `privesc` | 4 |

## 形态分布

| kind | 条数 |
|---|---|
| `markdown` | 2201 |
| `script` | 217 |
| `poc` | 150 |
| `exp` | 46 |
| `template` | 44 |
| `payload` | 8 |
| `http` | 6 |

## 全量条目

| id | 标题 | category | kind | 已验证 | tags |
|---|---|---|---|---|---|
| 1 | ThinkPHP 精确版本与绝对路径识别（异常页法） | info-leak | script | Y | thinkphp,fingerprint,info-leak,路径泄露,版本识别 |
| 2 | 未鉴权接口回显完整 SQL —— 表结构/资金字段/拼接写法泄露 | info-leak | poc | Y | thinkphp,sql报错回显,信息泄露,资金接口,未鉴权,mysql |
| 3 | WG-lobby（新葡京系）前端加密模型还原 + OSS 预签名写入利用链 | unauthorized | exp | Y | WG-lobby,AES,密钥派生,oss,presigned,未授权,反混淆 |
| 4 | WG-lobby 游客 token 数据面拖库（getSysInfo 取 token + chipher 重放 + 响应解密 + 接口普查翻页） | unauthorized | exp | Y | wg-lobby,未授权,数据拖库,游客token,AES-CBC,翻页 |
| 5 | 暴露数据库/SSH 服务弱口令与未授权探测脚本（MySQL/PostgreSQL/SSH 三合一） | weak-password | script |  | mysql,postgresql,ssh,弱口令,暴露服务,未授权 |
| 6 | NextAuth v4 JWE 会话令牌解密与伪造（弱 NEXTAUTH_SECRET 检测） | auth-bypass | script |  | nextauth,nextjs,jwe,forge,session |
| 7 | New API / one-api 未授权面枚举与源码鉴权对照审计清单 | info-leak | script | Y | new-api,one-api,未授权,信息泄露,源码对照,审计清单 |
| 8 | 1Panel 安全入口（Security Entrance）旁路探测：未认证路由与 /assets/ 静态目录不受入口保护 | auth-bypass | poc | Y | 1panel,security-entrance,未授权,目录遍历,信息泄露 |
| 9 | 1Panel 验证码绕过（CVE-2025-66507）+ RSA 公钥加密口令在线爆破脚本 | weak-password | script | Y | 1panel,captcha-bypass,CVE-2025-66507,爆破,RSA |
| 10 | 1Panel 登录口令加密格式（RSA 保护 AES 密钥的三段式）与验证码绕过 | weak-password | script | Y | 1panel,登录加密,RSA,AES,CVE-2025-66507,爆破 |
| 11 | express-rate-limit 限速绕过：伪造 X-Forwarded-For 变换限速 key | methodology | script | Y | 限速绕过,rate-limit,trust-proxy,X-Forwarded-For,Express |
| 12 | 纯 Python AWS SigV4 实现：无 awscli/boto3 时实测 AWS 凭据有效性（STS GetCallerIdentity） | methodology | script | Y | aws,sts,sigv4,凭据验证,云凭据 |
| 13 | 直连本机 ChromeMcpServer 驱动用户真实 Chrome（MCP streamable-http 自建会话 + 裸 IP 导航缺陷绕法） | env | script | Y | browser,chrome,mcp,真实登录态,ChromeMcpServer,<target>,工具坑 |
| 14 | Express 鉴权中间件路径判定绕过（大小写前缀 /API + URL 编码 .html），配套变体清单与只读收割顺序 | auth-bypass | script | Y | express,鉴权绕过,路径判定,大小写,url编码,信息泄露 |
| 15 | Dujiao-Next 发卡系统：订单/游客订单鉴权模型、API 契约、订单号格式与 IDOR 判定速查 | unauthorized | template | Y | dujiao-next,发卡系统,IDOR,鉴权模型,订单号格式,源码审计,速查 |
| 16 | Dujiao-Next 用户中心与认证链路速查：/me·wallet·affiliate·2FA 鉴权事实 + 验证码节流与限流缺口 | auth-bypass | template | Y | dujiao-next,用户中心,2fa,鉴权模型,验证码节流,限流缺口,速查 |
| 17 | cloudflare_temp_email <target> 鉴权模型与匿名端点 oracle（地址存在性枚举） | unauthorized | template | Y | cloudflare_temp_email,临时邮箱,未授权,用户枚举,版本差异分析 |
| 18 | LLM 中转面板渠道 config.quota.url SSRF（触发三条件 + 绕过前置 Squid 直连 + 可达性判据表） | ssrf | exp | Y | ssrf,llm-panel,proxy-bypass,quota,writable-config |
| 19 | Microsoft/Google OAuth refresh_token 特征 grep 清单 + aigemini-org 13030 数据面阴性结论 | info-leak | script |  | oauth,refresh_token,microsoft,outlook,grep,negative-result |
| 20 | AIRouter 白标 AI 网关：管理面鉴权绕过尝试集（/API/ 大小写 + /admin%2Ehtml 编码变体）与被证伪的迁移性 | auth-bypass | payload |  | express,鉴权绕过,case-sensitive,url-encoding,airouter,白标网关,证伪,判别清单 |
| 21 | LLM API key 零消费只读验活（DeepSeek 余额 / OpenRouter 分层 / 通用 /models） | info-leak | script | Y | llm,apikey,balance,readonly,zero-cost |
| 22 | gptcdk 卡密兑换站群自托管后端未授权卡密查询/枚举探测 | unauthorized | poc | Y | 卡密兑换中心,gptcdk,未授权访问,枚举,fastapi,c ChatGPT卡台 |
| 23 | nps web 管理端 auth_key=md5(timestamp) 认证绕过（含隧道/客户端清单导出） | auth-bypass | poc | Y | nps,内网穿透,auth_key,认证绕过,内网入口 |
| 24 | 自研卡密兑换站未授权批量查询接口 /api/v1/cdk/batch（20 张/请求，响应含套餐/订单） | unauthorized | poc | Y | unauthorized,卡密,批量查询,gptcdk,接口枚举 |
| 25 | nps 管理端接口取数姿势：/client/list 与 /index/gettunnel 必须 POST（GET 返回模板页/空集），gettunnel 必须带 type | info-leak | script | Y | nps,取数姿势,POST,gettunnel,云安全组,内网穿透 |
| 26 | MSSQL TDS PRELOGIN 探针（47 字节，struct 构造 + assert）—— 长度写错会被误判为"伪造开放" | info-leak | script | Y | mssql,tds,prelogin,探针,踩坑,长度校验 |
| 27 | GET/POST 分支鉴权不一致导致 userId 水平越权遍历（含无盐 MD5 口令泄露） | unauthorized | script | Y | 越权,IDOR,未授权访问,无盐MD5,主键遍历 |
| 28 | 方法论陷阱：本地 HTTP 代理会按 Host 分流，制造"Host 头命中源站"的假象（必须无代理复核） | recon | script | Y | 方法论,源站定位,Host头,代理陷阱,Cloudflare,对照实验 |
| 29 | 灰产资金站追踪链：灰云判别 → 中转页统计JS → Workers跳板 → 会员系统JS → SPA bundle 逆向 API | unauthorized | script | Y | 灰产,资金链,域名轮换,SPA逆向,API提取,Workers跳板,未授权 |
| 30 | 对象存储凭据有效性验证器（七牛 AK/SK 只读端点四连测） | info-leak | poc | Y | 七牛,又拍云,对象存储,未授权配置泄露,凭据验证 |
| 31 | 未授权越权接口 · 受控批量导出用户档案（IDOR 遍历 + 弱口令反查） | unauthorized | poc | Y | 未授权,水平越权,IDOR,用户档案导出,弱口令,MD5 |
| 32 | Druid StatViewServlet 暴露检测（含从 actuator/prometheus 的 statViewServlet 指标反推） | info-leak | poc | Y | druid,spring-boot,actuator,监控台,信息泄露 |
| 33 | 移动 App 侦察：从 pcap 配置分发文件定位真实后端 + 认证头校难强度对照实验 | info-leak | script | Y | 移动App,抓包分析,配置泄露,认证绕过,签名校验,未授权 |
| 34 | SockJS + STOMP WebSocket 客户端实现要点（JSON 数组发送侧编码 / 版本协商 / CDP 抄帧法） | unauthorized | script | Y | websocket,sockjs,stomp,协议逆向,CDP |
| 35 | 用 Actuator prometheus 的 exception 标签做参数名/字段名 oracle（异常处理器统一返回 200 时的参数枚举法） | info-leak | script | Y | spring-boot,actuator,prometheus,参数枚举,oracle,micrometer |
| 36 | Android App 自研 API 加密协议还原（DEX 反汇编定位 AES 密钥 + param 密文协议 + 登录三层校验差分定位） | frontend-reverse | script | Y | android,dex,reverse,aes,protocol,app-api,param-encryption |
| 37 | 阿里云验证码 2.0（AliyunCaptcha）自动化跑通并取得 captchaVerifyParam —— App/Web 登录风控绕过通用打法 | auth-bypass | script | Y | 阿里云验证码,captchaVerifyParam,风控绕过,App登录,自动化,Burp不可见 |
| 38 | SPA bundle 端点半自动提取 + 只读白名单批量未授权探测（含写语义端点硬跳过） | unauthorized | script | Y | 未授权,SPA逆向,批量探测,只读白名单,Java,Spring,范围外来源,方法本身通用 |
| 39 | 自研 Java 资金站 TIMESTAMP 通道：安全凭证头 + AES-ECB 请求体加密与时间片密钥派生算法 | unauthorized | script |  | AES-ECB,SHA256,时间片密钥,前端逆向,Spring,资金站,传输层,范围外来源,方法本身通用 |
| 40 | 目标不可达的「三出口对照」诊断法 —— 区分目标侧故障 vs 我方出口被封 | env | poc | Y | 可达性诊断,出口对照,502,排障方法论,避免白跑 |
| 41 | 阿里云 WAF 封禁维度识别：敏感端点为「方法+路径+来源」持续封禁（405 拦截页），非频率限流；附三态判据 | waf | script | Y | 阿里云WAF,acw_tc,Anti-DDoS,封禁维度,方法论,别再踩坑,范围外来源,方法本身通用 |
| 42 | 代理分流造假样本：<target> 经 7892 返回伪造 502/404「hello」 | methodology | poc | Y | 代理造假,7892,判定方法论,负结论 |
| 43 | OCS 在线客服系统：STOMP over SockJS 握手真值 + 客服侧 protobuf 信封 UserMsgPack | protocol | template | Y | OCS,SockJS,STOMP,protobuf,客服侧,UserMsgPack |
| 44 | SockJS 文本通道无法承载 protobuf 二进制：成因、对照实验与判据 | protocol | poc | Y | SockJS,STOMP,protobuf,二进制,帧长度,负结论,方法论 |
| 45 | nginx location 绕过：%2e%2e 系列必然失败，及用"不存在路径对照"证伪 ..;/ 假穿越 | path-traversal | poc | Y | nginx,location绕过,路径归一化,..;,判据,负结论 |
| 46 | nps web 管理面板 auth_key=md5(timestamp) 认证绕过 + 隧道/客户端枚举 | unauthorized | exp | Y | nps,内网穿透,认证绕过,tunnel,面板接管 |
| 47 | 密钥分发范式：4 段密钥/盐/公钥 XOR 0x5A 隐写在 4 个未授权图片接口的尾部 2048 字节（标记 [P1][P2][S][AES]） | info-leak | script | Y | 隐写,密钥分发,图片尾部,XOR-0x5A,AES-ECB,RSA-2048,SPKI,前端逆向,响应信封,同团伙复用,范围外来源,方法本身 |
| 48 | 阿里云 WAF 封禁维度修正实测：换出口无效（非 IP 维度）、裸 curl/Python 必被拦、真实浏览器可通过 | waf | script | Y | 阿里云WAF,acw_tc,封禁维度修正,浏览器可过,方法论,405误判,范围外来源,方法本身通用 |
| 49 | MQTT over WebSocket 未授权接入与 ACL 绕过（`#` 被拒 → `+/#` 全量放行，EMQX） | unauthorized | poc | Y | mqtt,emqx,websocket,acl-bypass,unauth,iot |
| 50 | Spring Boot Whitelabel Error Page 参数名 oracle（Accept: text/html 读出必需 @RequestParam 名） | info-leak | script | Y | spring,java,参数名泄露,信息泄露,接口枚举,oracle |
| 51 | 1Panel 验证码绕过对照判定 + 三段式加密登录探针（可直接换字典做在线口令审计） | weak-password | script | Y | 1panel,验证码绕过,CVE-2025-66507,三段式加密,口令审计,对照实验 |
| 52 | Tomcat AJP 8009 探测 + Ghostcat(CVE-2020-1938) 文件读取（含「HTTP 代理 CONNECT 盲回 200」的识破判据） | rce | script | Y | tomcat,ajp,ghostcat,8009,代理盲回,探测判据 |
| 53 | PowerShell 下 curl 内联 JSON body 被吃双引号 → 服务端全 500 被误判为目标故障 | methodology | script | Y | powershell,curl,json,false-positive,探测陷阱 |
| 54 | 1Panel 版本指纹（路由集合差集法，离线）：把面板压到具体 LTS 版本 | info-leak | script | Y | 1panel,版本指纹,路由差集,离线,swagger |
| 55 | 教训：1Panel 版本区间误判避坑指南（GHSA first_patched_version 漏掉 LTS 分支；正确判据 = swagger basePath + router 源码集合差集） | info-leak | poc | Y | 1panel,版本判断,教训,swagger basePath,路由集合差集,LTS分支,GHSA陷阱 |
| 56 | 本机红队环境事实与替代方案（Windows/ELF 平台限制） | env | script | Y |  |
| 57 | 环境事实：本机出口走全局代理（127.0.0.1:7892）；仅代理不可用时才改走演练 VPS | env | script | Y | 环境事实,出口链路,TUN,Windows,工具白名单,VPS跳板,nuclei,FOFA_KEY |
| 58 | 1Panel 安全入口拦截页与 Swagger 拼接 ⇒ 未认证绕读完整 swagger（v1.10.28-LTS 实测）+ 无 swagger 时的版本夹逼探针 | info-leak | poc | Y | 1panel,swagger,security-entrance,绕读,版本指纹,路由存在性探针 |
| 59 | GM Collection 邮件工具站免登录开放接口未授权探测清单（Outlook/iCloud 邮件与凭据查询） | unauthorized | script | Y | unauthorized,email,outlook,icloud,temp-mail,cloudflare-worker,curl |
| 60 | New API 漏洞清单与未修复版本边界对照（CVE-2026-42339 SSRF / 33655 / 64865 / 30886） | ssrf | template |  | new-api,one-api,ssrf,idor,cve,版本边界,源码审计,计费绕过,yaml |
| 61 | Vue3/Vite SPA 接口测绘法：入口 JS 只给路由，真实 API 在懒加载 chunk 的数据层 use*.js 里 | frontend-reverse | script | Y | 前端接口提取,Vite,Vue,chunk,SPA,接口测绘,只读 |
| 62 | 商城游客订单自认证头机制：Authorization: Guest base64url(email+"\n"+order_password)（含空口令/畸形头边界实测） | unauthorized | poc | Y | 商城,游客订单,自认证头,Authorization,base64,越权验证,待验证 |
| 63 | 白标 SPA 商城溯源到上游开源项目（品牌名兜底值 + app_version）→ 转离线源码审计 | frontend-reverse | poc | Y | 溯源,白标商城,开源项目识别,源码审计,版本锁定,默认密钥,零负载 |
| 64 | Webpack 异步 chunk 收割枚举 SPA 全量 API 面（new-api 定制构建实测） | info-leak | script | Y | webpack,spa,api枚举,异步chunk,new-api,前端信息收集 |
| 65 | 通用方法：开源组件识别 → 拉源码离线审计 → 定向验证（找 advisory 之外的路由漏挂） | methodology | script | Y | 源码审计,中间件,漏挂认证,gin,new-api,codeload,方法论 |
| 66 | 凭据来源侦察：公开代码托管泄露检索入口清单 + 仓库全树密钥扫描法 | info-leak | script | Y | 信息泄露,凭据来源,代码托管,供应链,被动收集 |
| 67 | New API / one-api 网关令牌有效性判别法（Invalid token vs User has been banned） | info-leak | poc | Y | new-api,one-api,令牌判别,代码泄露,凭据,零爆破 |
| 68 | 被动 DNS / 证书透明可用源清单与 CDN 误判防护（实测可用性矩阵） | recon | template | Y | 被动DNS,证书透明,源站,CDN误判防护,信息收集 |
| 69 | CVE-2026-81027 one-api/New API channelid 路径参数越权锁定信道（任意有效 token，无需余额） | auth-bypass | exp |  | one-api,new-api,越权,CWE-862,channelid,提供商key,无需余额 |
| 70 | 环境事实（全体适用）：本机 Python urllib 走 IE 注册表代理 127.0.0.1:7892 导致请求超时，curl 正常 | env | script | Y | 环境事实,windows,urllib,代理,curl,超时,误判 |
| 71 | Discuz 7.x ajax.php checkusername 用户名枚举（零配额成本判定账号是否存在） | info-leak | poc | Y | discuz,用户名枚举,信息收集,零配额 |
| 72 | phpwind <target> 前台 SQL 注入（actions/ajax/leaveword.php） | sqli | poc |  | phpwind,sql注入,front-end |
| 73 | 智讯网站群管理系统 /api/ 未授权访问（form-urlencoded 传参 + action=index 泄露站点配置） | unauthorized | poc | Y | php,站群,未授权,api,form-urlencoded,智讯网 |
| 74 | 中文站群"网站防火墙"WAF 识别判据（403/1200B 中文页 + 按扩展名拦截 + 多长度分类法） | waf | poc | Y | WAF,网站防火墙,403,扩展名拦截,判据,phpwind,绕过矩阵未验证,站群通用 |
| 75 | phpwind 头像上传任意文件写入（job.php?action=uploadicon&step=2&filename=x.php） | file-upload | exp |  | phpwind,getshell,文件上传,头像,写入 |
| 76 | phpwind 附件上传扩展名黑名单审计结论（前台变形不可行，口子在后台 $db_uploadfiletype） | file-upload | poc | Y | phpwind,文件上传,黑名单,证伪,后台getshell |
| 77 | 六零导航页 LyLme Spage ≤<target> 后台 getshell（ajax_link.php submit=update 任意 zip 解压到站点根）+ 未授权 SSRF | rce | exp |  | lylme,六零导航页,导航站,后台getshell,任意文件下载,zip解压,SSRF,未授权,RCE |
| 78 | 智讯网站群：注册接口用户名枚举 + 口令大小写不敏感判定 | auth-bypass | poc | Y | php,站群,用户枚举,口令策略,注册接口,智讯网 |
| 79 | 六零导航页 LyLme Spage admin_token 离线伪造器（SYS_KEY 硬编码 lylme_key → 任意配置泄露即通吃会话） | auth-bypass | script | Y | lylme,六零导航页,SYS_KEY硬编码,admin_token,离线伪造,cookie伪造,会话绕过 |
| 80 | 网段存活探测纪律：nmap -Pn 会把全部 IP 判为 up，必须用 -sn 复核 | recon | script | Y | nmap,存活探测,假阳性,纪律,-Pn |
| 81 | FOFA_KEY 在本机需从 .dsh/.env 显式加载（shell 环境变量未继承） | env | script | Y | fofa,环境变量,信息收集,踩坑 |
| 82 | 主域 NXDOMAIN 后从姊妹域证书按机器命名规律挖跨段新资产 | recon | script | Y | crt.sh,证书透明,子域枚举,站群,被动信息收集,跨段扩面 |
| 83 | SSRF 判据三陷阱：外层码≠目标状态 / icon 非空≠端口开放 / 快速返回≠存活 | ssrf | poc | Y | SSRF,判据陷阱,get_head,误判,方法论,icon字段,端口探测 |
| 84 | 假阴性判据：路径级 403 ≠ IP 封禁 / curl 默认 UA 被拦 ≠ 不可达 / 按时间下线 ≠ 被封 | waf | poc | Y | 假阴性,封禁判定,路径级403,curlUA,关站,对照实验,方法论,全队踩坑 |
| 85 | 批量扫描三条对照实验：戳破假存活 / 假 403 / 假全开 | recon | template | Y | 对照实验,假存活,假全开,假403,nmap,端口验证,信息收集,方法论 |
| 86 | JWT HS256 秘密密钥组合爆破器（无依赖，含 0-99999 数字后缀与词组组合） | tunnel | script | Y | jwt,hs256,bruteforce,authententication |
| 87 | 自研 Vue/Vite SPA 面板：从 bundle 递归还原完整 API 面（含 /admin 与开放 API 协议） | frontend-reverse | template | Y | spa,vite,vue,recon,api-surface |
| 88 | 智讯网站群 root_user_index 免认证读取全站用户表（含明文 opass）｜附 curpage 分页被忽略的踩坑 | unauthorized | poc | Y | 站群,未授权,用户表泄露,uid令牌,分页陷阱 |
| 89 | 从单点证书反查整个站群（cert= 一把捞全网同证书资产） | recon | script | Y | 证书反查,fofa,crt.sh,站群扩面,通用证书,信息收集 |
| 90 | 第三方 IP 归属库不可靠：必须用响应特征自行验证 | recon | template | Y | IP归属,资产验证,误伤风险,合规,对照实验,信息收集 |
| 91 | 指令边界歧义处理：指定 IP 白名单 vs 禁扫网段（定向优先，显式确认） | methodology | template | Y | 授权边界,指令歧义,禁扫网段,白名单,协作规范,指挥意图 |
| 92 | LFI 注入点定位法：先找不经过守卫的可达文件目录，再用数组参数判定是否读参数 | path-traversal | script | Y | LFI,任意文件读取,注入点定位,数组参数,display_errors,方法论,Windows穿越,负结论 |
| 93 | 取证表述纪律：区分「实测」与「推断」—— 同一证据存在两个同样合理解释时，未排除的那个不得写进结论 | methodology | poc | Y | 取证纪律,实测与推断,对照实验,假阳性,自我纠错,表述规范,前缀拦截 |
| 94 | 从 TLS 证书 CN 反查不走 CDN 的子域 ⇒ 源站 IP 泄露 + 绕过 Cloudflare（含 HTTP.sys Host 头路由命中技巧） | info-leak | script | Y | recon,cloudflare,源站泄露,http.sys,host头路由,cert-cn |
| 95 | 同一 ApiInterface 跨控制器对照法：定位漏挂 [Authorize] 的控制器 | unauthorized | script | Y | unauthorized,aspnet,mvc,authorize,反射接口,对照实验 |
| 96 | nginx `/login` 通配重定向规则会造成 ffuf 假阳性：命中项必须用『任意词对照』核实 | methodology | script | Y | 假阳性,nginx重写,通配规则,目录爆破,ffuf,对照实验,Location相对路径,方法论 |
| 97 | <target>/.23 先验排除表：两条恒定假阳性 + /code 验证码端点（爆破前置） + 五处已穷尽面 | methodology | script | Y | 先验排除表,假阳性,nginx通配规则,整前缀认证,验证码端点,已穷尽面,449918,宝塔 |
| 98 | HTTP Range 精准抽取远端 zip 内任意单文件（不下载整包） | info-leak | script | Y | range,zip,源码泄露,目录浏览,取证 |
| 99 | vhost 真伪判定（SNI 双置位法）：curl --resolve 同时置位 SNI+Host + 逐 SNI 证书对比 | recon | script | Y | vhost,sni,iis,方法论,避坑 |
| 100 | MSSQL 1433 防火墙握手代理式伪装的判定法与 TDS PRELOGIN 探针 | waf | script | Y | mssql,1433,TDS,可达性判定,防火墙伪装,误判避免 |
| 101 | curl --resolve 端口陷阱：URL 不带端口会静默逃逸到真实 DNS（伪造出"目标有业务"的假象） | recon | script | Y | curl,取证陷阱,方法论,vhost,避坑,静默失败 |
| 102 | 证书域名提取陷阱：leaf SAN 与链内 OCSP/CRL/CA-Issuers 域名必须分开（否则误当业务域名） | recon | script | Y | 证书,SAN,AIA,CRL,取证陷阱,资产测绘,避坑 |
| 103 | 会话态不完整时后台功能页静默返回 200/0 字节 —— 必须把"0 字节 200"当"会话可疑"而非"页面为空" | methodology | script | Y | aspnet,后台,会话,误判,踩坑,geek |
| 104 | 浏览器实测必须用 hosts 映射：Chrome --host-resolver-rules 在本机 153 版不生效（附 Server 头复验判据与当即还原纪律） | env | script | Y | 浏览器实测,hosts,cloudflare,源站直连,踩坑,cdp |
| 105 | 带验证码且无回显差异的后台弱口令打法：验证码放大 6 倍识别 + 每组口令独立会话 + 以 302 为唯一成功信号 | weak-password | script | Y | 弱口令,验证码,ocr,后台,aspnet,独立会话,方法论 |
| 106 | 记分被 capped 时的正确处置：追加式更正（不改 target 规避去重，追加 chain + 报告显式呈现差异） | methodology | script | Y | 记分,去重,capped,数据完整性,自我纠错,方法论 |
| 107 | Discuz <target> failedlogins 登录失败计数绕过（伪造 Client-IP / X-Forwarded-For 实现无限爆破） | auth-bypass | poc | Y | discuz,brute-force,auth-bypass,ip-spoofing,failedlogins |
| 108 | Discuz 管理员 Cookie 伪造（authkey+UA 签名 / 由 SQL 注入取三要素） | auth-bypass | exp | Y | discuz,auth-bypass,cookie,authkey,authcode,sqli |
| 109 | SCLUB/168six：403 站点规则可被本站 Referer 绕过（探测假阴性） | auth-bypass | poc | Y | SCLUB,Discuz,Apache,403绕过,Referer,防盗链,信息收集假阴性,防盗链规则 |
| 110 | Discuz 论坛平台（SCLUB 免费论坛）连接层封禁阈值探测方法 —— 区分「并发触发」与「累计触发」 | waf | poc | Y | discuz,brute-force,rate-limit,waf-bypass,methodology,sclub |
| 117 | 未授权接口必须"重读+diff"：静态认知固化导致漏掉新增明文凭据（19 轮作业的真实教训） | info-leak | template | Y | 方法论,增量监控,未授权接口,认知固化,settings,明文凭据,作业纪律 |
| 118 | 配置类接口（settings/config/option）明文回显管理员口令与 API Key —— 字段特征与取用方法 | info-leak | template | Y | settings,明文口令,api_key,未授权访问,defaults泄露默认口令,自研面板 |
| 119 | 本地只读反代：把前端调用的 /api/ 改写为大写转发，让原版管理后台界面直接可用 | unauthorized | script | Y | 本地反代,大小写绕过,只读浏览,面板,Express,CaseSensitive,工具化 |
| 120 | 主列表漏报陷阱：明细统计必须逐对象走详情接口（本例少统计 30%） | info-leak | template | Y | 统计口径,漏报,主从结构,apikey_path,404误解,集合比对 |
| 121 | 凭据复用验证 vs 爆破的边界：拿到实证口令后，同集群各试 1~2 次的正确打法 | weak-password | template |  | 凭据复用,爆破边界,1Panel,RSA+AES三段式,Redis,横向,纪律 |
| 122 | SSRF 内网探测的两个陷阱：服务端自身卡顿伪装成"目标超时"，以及把公网可见服务误报成"内网突破" | info-leak | template | Y | SSRF,内网探测,超时误判,对照组,方法学,攻击面核对,报告口径 |
| 123 | 登录"成功"的真伪三重证伪：字段校验 / 鉴权中间件 / 浏览器注入会话+清 Cookie 对照 | auth-bypass | template | Y | 口令复用,假成功,会话验证,三重对照,计分纪律,凭据复用边界 |
| 124 | 路由与端点发现的四个高频误判：OPTIONS 不能证明路由存在、404 必须在三态对照下读、白名单形态特征、稳定性基线 | methodology | template | Y | 路由发现,OPTIONS,三态对照,404vs401,白名单识别,稳定性基线,收口判据 |
| 125 | cloud-mail 认证机制与探测要点（JWT + KV 双校验 / /init secret 预言机 / KV 键名规律） | auth-bypass | script |  | cloud-mail,cloudflare-workers,jwt,kv,认证机制,预判否定 |
| 126 | cloudflare_temp_email 认证机制与攻击面（四套认证头 / address JWT 需查库 / admin_login 免认证爆破入口 / DISABLE_ADMIN_PASSWORD_CHECK） | auth-bypass | script | Y | cloudflare_temp_email,cloudflare-workers,临时邮箱,认证机制,x-admin-auth,admin_ |
| 127 | 前端 chunk 硬编码凭据 → 共享钥匙 → 密钥库路径枚举批量导出 | unauthorized | script | Y | 前端硬编码,凭据泄露,密钥库,Cloudflare Workers,路径枚举,未授权访问 |
| 128 | Vite dev server 暴露（Astro 生产误跑 dev）—— /@fs/ + ?raw 读取项目源码 | info-leak | poc | Y | vite,dev-server,astro,源码泄露,任意文件读,/@fs,?raw,robots内部路径 |
| 129 | 1Panel v2 安全入口启用时的可打面与判定手册（匿名白名单路由 + 验证码是否为唯一前置） | info-leak | script | Y | 1panel,v2,安全入口,验证码,匿名路由,CVE-2025-66507,CVE-2025-66508,端口9527 |
| 130 | 1Panel 口令审计前置判据：必须先有 panel_public_key 与（安全入口启用时的）entrance，且无用户名 oracle | weak-password | script | Y | 1panel,口令审计,前置条件,panel_public_key,RSA三段式,安全入口,ddddocr,验证码,用户名oracle |
| 131 | 1Panel v2 登录验证码自动化攻破（ddddocr）+ ErrAuth/ErrEntrance 命中判据 + “有安全入口即不可行”结构性结论 | weak-password | script | Y | 1Panel,v2,验证码,ddddocr,OCR,ErrEntrance,安全入口,爆破可行性,端口9527 |
| 132 | 开放代理“穿内网”真伪判定法：内网 vs 公网同 SNI 对照 body hash（附 CONNECT+自定义SNI 手法） | tunnel | script | Y | squid,开放代理,内网突破,boundary,SNI,vhost,对照实验,判据,CONNECT |
| 133 | 1Panel 面板登录验证码绕过（ignoreCaptcha=true → ErrAuth） | auth-bypass | poc | Y | 1panel,验证码绕过,登录,弱口令,明文口令 |
| 134 | 端口定性三态判据：中间盒黑洞 vs 真 filtered vs 真服务（同机对照 + 12 载荷 + SNI 矩阵） | info-leak | script | Y | 端口定性,黑洞,filtered,对照实验,ECONNRESET,TLS,SNI,nmap,FOFA,8443 |
| 135 | Flask before_request 鉴权跳过 OPTIONS ⇒ 未授权枚举后台路由表与方法（含变量规则噪声判据） | info-leak | script | Y | flask,werkzeug,OPTIONS,before_request,未授权,路由枚举,Allow头,变量规则,噪声 |
| 136 | New API 支付回调/初始化免鉴权面判定清单（epay notify + stripe/creem/waffo webhook + /api/setup） | unauthorized | script | Y | new-api,epay,stripe,webhook,setup,未授权,支付回调,验签 |
| 137 | Discuz! <target> RELEASE 20101020 后台 Getshell 全路径（离线白盒审计手册） | rce | exp |  | discuz,<target>,getshell,后台,into outfile,模板编译,uc_key,config.inc.php,sq |
| 138 | 单请求型隧道下的内网渗透方法论（op=chk 主机发现 + 目标侧 PHP 绕过多轮协议限制） | tunnel | script | Y | 内网渗透,隧道,WebShell,op=chk,NFS,RPC,open_basedir绕过,主机发现,静默失败,eval马,冰蝎 |
| 139 | phpwind <target> 站群通用探测判据集（登录接口响应语义 + 额度机制实测 + 未授权接口 + 目录 deny 清单） | weak-password | template | Y | phpwind,<target>,判据,登录接口,额度,锁定,XFF,regcheck,未授权,目录deny,上传补丁,站群通用 |
| 140 | phpwind <target> 补丁面速查表（<target> 白盒产出）：已修复 9 类 + 未验证缺陷 6 条 + 未穷尽面 4 个（全附源码行号） | app-audit | template | Y | phpwind,<target>,补丁面,白盒审计,速查表,GBK,宽字节,addslashes,复用,同圈子,源码行号 |
| 141 | 博彩平台会话网关 /ssid1?url= 反射型 XSS + 开放重定向（可窃取 ssid1 会话 cookie，会员/代理/管理三端命中） | auth-bypass | poc | Y | XSS,开放重定向,会话网关,ssid1,博彩平台,会话劫持,三端,跨站复用 |
| 142 | phpwind <target> 存储型 XSS 选点判据：checkVar 豁免名单 + prosign 零转义写入 + apartment 校验绕过 + groupid=7 拦截真相 + 签名展示位对照 | auth-bypass | poc | Y | phpwind,<target>,存储型XSS,prosign,checkVar豁免,apartment绕过,groupid7,展示位,站群 |
| 143 | phpwind <target> 注册三必填绕过（regreason + _hexie + apartment 地区ID）→ 注册成功且自动登录 | auth-bypass | poc | Y |  |
| 144 | phpwind <target>「decode 错位链」白盒审计结论：无需账号的裸引号→convert(/e) 路径不存在；CMS 投稿链裸引号入库→展示 RCE（代码级成立）；forum 发帖链被 addslashes+escapeChar 双保险挡住 | rce | template |  | phpwind,<target>,decode错位,裸引号,code-e,convert,CMS投稿,白盒审计,负结论 |
| 145 | 入驻方资料接口未授权 IDOR（shop_user_id 任意遍历泄露商家微信/公众号） | unauthorized | poc | Y | idor,unauthorized,business-logic,pii,wechat |
| 146 | 域名清单接口未授权泄露全平台 20 个生产域名 + 后端回源 IP | info-leak | poc | Y | info-leak,unauthorized,recon,subdomain,backend-ip |
| 147 | agent_id 鉴权分支绕过 —— 免登录读取付费内容平台正文（goodsHistory/goodsList/goodsDetail） | auth-bypass | poc | Y | auth-bypass,idor,未授权,付费内容,agent_id,thinkphp,thinkphp5 |
| 148 | SCLUB 平台 <target> 注入客户端修正（繁体歡迎+hex左截断+FROM包层+限流特征） | frontend-reverse | script | Y |  |
| 149 | SCLUB 平台 <target> 注入回显客户端（含三个易踩根因与语句约束表） | auth-bypass | script | Y | sql注入,union,回显截断,繁体字符,old-password,sclub,客户分库站群 |
| 150 | acg-faka（异次元发卡）免登录用户 API 访问控制模型与判定基准（secret/query/commodity/card） | unauthorized | template | Y | acg-faka,异次元发卡,发卡网,未授权,越权,列过滤盲注,源码判定 |
| 151 | Dujiao-Next Admin 后台指纹与 /api/v1/admin/** 管理 API 全清单提取（含未授权定性探测） | info-leak | template | Y | dujiao-next,独角数卡NEXT,发卡网,后台,API枚举,未授权,Vue3,源站直连 |
| 152 | 若依系前台登录「请求加密守卫」绕过（isEncrypt/AES-ECB + RSA encrypt-key，前端硬编码密钥对） | auth-bypass | script | Y |  |
| 153 | Dujiao-Next 版本指纹 + 上游补丁差分法（app_version 免认证泄露 → 缺失安全修复清单 → 逐条可达性实测） | info-leak | script | Y | dujiao-next,独角数卡,发卡网,版本指纹,信息泄露,补丁差分,未授权 |
| 154 | agent 后台忘记密码流程逆向：未鉴权改密接口 getBackLoginPassword 完整字段契约与差分探测法 | auth-bypass | script | Y | agent,backoffice,password-reset,unauth,getBackLoginPassword,forgetPass |
| 155 | X-API-XXX WAF 签名头逆向：路径绑定规则 + 验签缺失（任意十六进制即可穿透） | env | script | Y | waf,signature,x-api-xxx,nextjs,moleculer,replay,bypass |
| 156 | Next.js SSR pageProps.t 密文可解密 —— 前端硬编码 AES 密钥 V/GO`(h$-?93a8jg | frontend-reverse | poc | Y | 前端加密,SSR加密,硬编码密钥,信息泄露 |
| 157 | 主站登录 X-API-FINGER 指纹算法还原 + 6022 真实成因纠正（FINGER 非决定因素） | auth-bypass | poc |  | 登录绕过,风控指纹,前端逆向,FingerprintJS,murmur3 |
| 158 | 代理后台 loginFlow 两步验证流程完全逆向（isValidate/checkCode/XCode/XCodeId），卡点精确定位为极验票据 | auth-bypass | poc | Y | 登录流程,极验,两步验证,代理后台,错误码差异 |
| 159 | glaxy 网关 (_glaxy_tus5r4_ 命名空间) API 探测与签名复现 —— 含 10050 证伪与 8 证书多租户识别 | auth-bypass | script | Y | python,springboot,签名网关,命名空间,地域封锁,信息收集,反证zabbix,多租户证书 |
| 160 | glaxy 业务网关签名复现（E=字符降序）+ 匿名未授权读取后端配置/密钥/内网IP | auth-bypass | exp | Y | 签名绕过,字符降序,未授权访问,glaxy网关,appid泄露,越权 |
| 161 | club-upload 未授权任意图片上传 + AES-ECB 响应解密（密钥 7EGWY67N31RFK7GY）+ FastDFS 落地回读 | unauthorized | poc | Y | unauthorized-upload,file-storage,fastdfs,aes-ecb,frontend-key-recovery |
| 162 | 纯 Python RTMP 握手与 chunk 层实现（修掉 message header 顺序与 S2 残留两个坑） | protocol | script | Y | rtmp,amf0,srs,protocol,fuzzing,harness,auth-probe |
| 163 | phpwind <target>（<target>:1852）补丁面黑盒实测矩阵：13 类已封堵 + 5 个门槛不可越 + write_config 落点存在但入口不可达 + 目录 deny 与限流判据 | recon | template | Y | phpwind,<target>,补丁面,黑盒实测,负结论,站群复用,write_config,速率限流,判据 |
| 164 | phpwind <target> download.php type=message 零鉴权越权读任意附件（游客可读） | unauthorized | poc | Y | phpwind,附件下载,越权,未授权访问,type参数,getDownloadFactory |
| 165 | <target> 字符串数组还原 + 前端 AES 密钥派生还原（Vite/Vue SPA 通用） | frontend-reverse | script | Y | <target>,JS混淆还原,RC4,字符串数组,AES,CryptoJS,前端加密,配置解密 |
| 166 | 响应/请求加解密完整还原（金沙904.com 系 SaaS 博彩门户前端） | frontend-reverse | script | Y |  |
| 167 | Vite SPA 多子应用 API 面递归提取（子应用发现 + manifest chunk 清单 + API 正则） | info-leak | script | Y | Vite,SPA,子应用,manifest,API提取,前端逆向,攻击面发现 |
| 168 | 绕过前端响应加密 - 读取 Vue/Pinia 运行时 store 获取明文数据 | info-leak | script | Y | 前端加密绕过,Vue,Pinia,运行时状态,响应解密,SPA,CDP |
| 169 | Vue Router 路由全量枚举 - 获取 SPA 完整页面地图 | info-leak | script | Y | Vue Router,SPA,路由枚举,页面测绘,前端逆向,CDP |
| 170 | Java Servlet 未授权数据读取探测（POST + userid 参数 + 404/500/200 判别） | unauthorized | poc | Y | java,servlet,tomcat,未授权访问,参数爆破,userid |
| 171 | SPA 管理后台源码未授权下载 + action 全量未授权测试（含子目录路径发现） | auth-bypass | poc | Y | spa,管理后台,源码泄露,未授权访问,action枚举,ecdsa |
| 172 | 菠菜站群 aocai-api 未授权读取加密站点配置（AES-ECB 密钥硬编码于前端 app.js） | unauthorized | poc | Y | 菠菜站群,prodSite5,aocai-api,Tomcat,未授权访问,AES,前端加密破解,信息泄露 |
| 173 | aocai 站群 API 网关未授权访问（userGuestRegister 签发 token + getWebConfigNew 配置泄露） | unauthorized | poc | Y | aocai,站群,未授权,spring-boot,mongodb,token签发,AES-ECB |
| 174 | 菠菜站群 prodSite5 未授权文件上传（attachment/upload + uploadChat，落地 Azure Blob 公网可读） | file-upload | poc | Y | 菠菜站群,prodSite5,aocai-api,未授权文件上传,Azure Blob,文件写入 |
| 175 | 前端混淆/Dart AOT 双路还原 WebSocket 上行协议（含心跳与可达性判定） | frontend-reverse | script | Y | websocket,混淆还原,flutter,dart-aot,协议逆向,<target>,心跳 |
| 176 | aocai 站群前端双 AES 密钥完全复现（webConfig.key 派生，离线解密响应） | info-leak | exp | Y | aocai,站群,前端加密,AES,ECB,密钥派生,离线解密,<target>,CryptoJS |
| 177 | WordPress XML-RPC system.multicall 批量口令爆破（含判定特征校准） | weak-password | script | Y | wordpress,xmlrpc,multicall,brute-force,旅 |
| 178 | WordPress REST API 未授权信息收割（用户枚举 + 全量页面媒体 + 版本指纹） | info-leak | script | Y | wordpress,rest-api,未授权,用户枚举,信息泄露 |
| 179 | 优客服 UKefu 路径 ACL 绕过（/chat/ 前缀放行 + sign 仅长度校验） | auth-bypass | poc | Y | ukefu,优客服,java,spring,403绕过,访问控制,sign,路径黑名单 |
| 180 | Sunflower 2.0 博彩系统 API 未授权访问（含 TLS 拒绝绕过法） | unauthorized | poc | Y | sunflower,白标博彩,未授权访问,api,gocache,tls绕过,ip直连,dotnet |
| 181 | 白标真人娱乐场 H5 客户端加密/签名方案完整还原（AES-128-CBC IV=key + zlib + HmacSHA1） | frontend-reverse | script | Y | 白标博彩,加密逆向,AES,HmacSHA1,game-http,网关,token,内联bundle |
| 182 | 尾斜杠触发 Traefik 301 清理泄露隐藏路由 BASE（通用绕过手法） | info-leak | script | Y | traefik,nginx,路径清理,隐藏路由,访问控制绕过,尾斜杠 |
| 183 | 执行式去混淆：伪装运行环境跑混淆脚本，直接读回 window 配置（通用手法） | info-leak | script | Y | 前端混淆,去混淆,rc4,字符串表,node,配置泄露,6686 |
| 184 | Jackson 未知字段报错泄露内部类名与 DTO 字段清单（Spring Boot 接口测绘手法） | info-leak | poc | Y | jackson,spring-boot,信息泄露,dto测绘,接口结构,未鉴权,6686 |
| 185 | 拼图滑块验证绕过：clientNonce 明文内含滑块正确坐标且可离线伪造 | auth-bypass | poc | Y | 滑块验证码,拼图,clientNonce,验证码绕过,签名,博彩,6686,自研网关 |
| 186 | gameplat-cms /api/img/upload 未授权任意文件上传（Content-Type 可控 → 存储型 XSS） | file-upload | poc | Y | 未授权,文件上传,对象存储,存储型XSS,博彩CMS |
| 187 | CN7 平台 /jumpto/img/&lt;url&gt; 未授权开放图片代理 SSRF | ssrf | poc | Y | ssrf,图片代理,CN7,未授权 |
| 188 | OpenSSH regreSSHion CVE-2024-6387 崩溃验证（7etsuo PoC 落盘 + Python 复刻 v1/v2） | methodology | exp |  | openssh,sshd,regresshion,race-condition,DoS |
| 189 | Vite + javascript-obfuscator SPA：还原 chunk 映射表 / 字符串表反混淆 / 服务层接口全量提取 | frontend-reverse | script | Y | vite,webpack,chunk,javascript-obfuscator,反混淆,接口提取,SPA |
| 190 | OCS 在线客服中台 访客会话越权读取（vid IDOR） | auth-bypass | poc | Y | idor,越权,在线客服,OCS,cid,x-v-token,snowflake,未授权 |
| 191 | 自研 Web 后台接管四步法：app.js 提 API 面 + 真实域名绕 Origin 白名单 + 无锁定弱口令 + 数据面遍历 | weak-password | exp | Y | 弱口令,origin绕过,sni,自研后台,app.js,信息泄露 |
| 192 | 站群 RN APP 分发链下载 + bundle 接口/密钥提取（域名轮换跳转 JS 型分发） | info-leak | script | Y | react-native,apk,逆向,接口提取,硬编码密钥,下载链 |
| 193 | PHPWind 零鉴权信息泄露双接口：sitemap.php 泄露全站帖子标题+作者，member.php 泄露全部会员名 | info-leak | poc | Y | phpwind,信息泄露,sitemap,member,零鉴权,用户名枚举,站群 |
| 194 | PHPWind v4/5/6/7/8 hack/bank/index.php 银行排行榜 bank_sort.php 代码注入 RCE（官方称 <target>+ 才修复 → <target> 受影响） | rce | exp |  |  |
| 195 | PHPWind 7.5 apps/groups/share/index.php $route/$basePath 未初始化任意文件包含（前台零登录、可远程包含 getshell） | rce | exp |  |  |
| 196 | PHPWind 7.5 pw_api.php → api/class_base.php $mode 未过滤任意文件包含（来自 $_POST+$_GET，不依赖 register_globals） | rce | poc |  |  |
| 197 | PHPWind 7.x/8.x 后台 getshell 手法全谱（SQL INTO OUTFILE / 模板包导入 / 备份恢复 / cache.php+level.php eval 点 / 附件扩展名）+ 后台入口定位 | rce | template |  |  |
| 198 | PHPWind <target> 负面情报清单（register_globals 已移除 / mawhole eval 证否 / setstyles $ 转义 / hack.php 白名单锁死 / api 全 404）+ 唯一命中的 nuclei 模板 | app-audit | template |  |  |
| 199 | PHPWind <target> 负面情报：member.php 的 orderway/asc 排序参数未进入 SQL（纯回显） | app-audit | template | Y | phpwind,<target>,负面情报,orderway,member.php,注入,判据,假阳性 |
| 200 | PHPWind 7.x hack/bank savearray → data/bbscache/bank_sort.php 用户名注入（写码+同请求include） | rce | poc |  | phpwind,savearray,变量插值,写码执行,bbscache,<target>未修 |
| 201 | PHPWind global.php InitGP(array('u','a')) 装载 $a 但被无条件 unset —— ${@eval($a)} 取不到值的判据 | rce | poc |  | phpwind,InitGP,register_globals,变量注入,死payload判据 |
| 202 | PHPWind 7.x /eis 注入点：copyctrl('$read[colour]') —— 需 pw_threads.colour 列存在 | rce | poc |  | phpwind,preg_replace,/e,复制控制,copyctrl,静态化 |
| 203 | PHPWind 7.x 入参清洗边界与 5 个 POST 豁免字段（payload 构造必须遵守） | app-audit | script |  | phpwind,参数清洗,CheckVar,Char_cv,payload构造约束 |
| 204 | PHPWind uploadicon 头像上传在 7.5 SP3 的校验实现（<target> 需实测是否有白名单） | file-upload | poc |  | phpwind,uploadicon,头像上传,任意文件写入,7.5已加校验 |
| 205 | PHPWind 7.x 后台入口机制：$admin_file = PHP_SELF，全站唯一入口 admin.php，install.php 存在则后台拒绝进入 | app-audit | script |  | phpwind,后台入口,admin_file,PHP_SELF,adminjob,入口爆破 |
| 206 | PHPWind passport_client.php 通行证接口：verify=md5(action+userdb+forward+db_pptkey)，可改写任意用户密码并下发登录 Cookie | auth-bypass | poc |  | phpwind,passport,通行证,认证绕过,账号接管,未授权面 |
| 207 | PHPWind 7.5 SP3 全站 /e 与代码求值点穷举（30 处 /e，仅 2 处可注入，均为 $read[colour]） | rce | script |  | phpwind,preg_replace,/e,eis,代码求值,全量清单,<target>未修 |
| 208 | colour 判死与 7.5 等价字段 titlefont 证据链（/eis 注入必须的 ' 与 ) 都被前台清洗吃掉） | rce | script |  | phpwind,colour,titlefont,/eis,判死,历史遗留字段,<target>未修 |
| 209 | PHPWind 7.5 SP3 上传扩展名变形的黑名单缺口（pht/phar/inc 未覆盖） | file-upload | poc |  | phpwind,附件上传,扩展名变形,pht,后台白名单 |
| 210 | PHPWind mode.php / apps.php 路由审计：$m 与 $q 都不进路径，两侧均无文件包含（含 mode/o 真实文件清单） | app-audit | script |  | phpwind,mode.php,apps.php,路由审计,selectMode,Pcv,文件包含,判死 |
| 211 | 照片上传（mode.php?m=o&q=photos&a=upload）判死：扩展名白名单硬编码 + 文件名全服务端生成 + GetImgSize 内容校验，.user.ini 不可行 | file-upload | poc |  | phpwind,照片上传,photoupload,扩展名白名单,user.ini,判死 |
| 212 | PHP /e 三类替换串（形态 A 死 / 形态 B 语法可行 / 形态 C 真危险）+ Char_cv 不处理括号与 CheckVar 白名单的实测印证 | rce | script | Y | php,/e,preg_replace,php_addslashes,preg_do_eval,引擎级判据,否证 |
| 213 | PHPWind 7.5 SP3 全站 30 处 /e 逐条判定（含 atc_content 最强候选链、prosign 与 titlefont 链路判死） | rce | script |  | phpwind,/e,逐条判定,atc_content,prosign,titlefont,CheckVar,判死 |
| 214 | PHPWind 7.5 SP3 全部 /e 替换串原文与形态判定：30 处无一为形态 B（参数位逃逸），全死 | rce | script |  | phpwind,/e,形态A,形态B,替换串原文,逐条判定,判死,参数位逃逸 |
| 215 | PHPWind <target> 实测判据集：prosign 写入姿势+引号实体化、admincheck 无自助入口、mode.php?m=o 各 q 可用性、载体长度硬上限 | recon | template | Y | phpwind,<target>,prosign,admincheck,mode.php,判据,负面情报,站群复用 |
| 216 | PHPWind 7.5 SP3 反序列化判死：全站 0 个魔术方法（无链起点）+ 无 cookie/请求参数来源 | deserialization | script |  | phpwind,反序列化,unserialize,POP chain,魔术方法,判死,版本错位 |
| 217 | PHPWind 7.5 SP3 cookie 登录态：StrCode 无 MAC 可逆，但伪造受 db_hash 封死（拿到 db_hash 即账号接管） | auth-bypass | script |  | phpwind,cookie,StrCode,PwdCode,User_info,db_hash,登录态,账号接管 |
| 218 | 站群同源判定：多 vhost 是否指向同一 DOCUMENT_ROOT（body diff 相似度 + 尺寸 = 4×主机名长度差） | recon | template | Y | 站群,同源判定,body,diff,DOCUMENT_ROOT,Cookie前缀,判据,负面情报 |
| 219 | 端口开放判定判据：代理 CONNECT 假成功陷阱 + 以应用层数据为准（宝塔 Pure-FTPd / nginx 403 实测） | recon | template | Y | 端口复验,CONNECT假成功,应用层判据,宝塔,Pure-FTPd,phpMyAdmin,IP白名单,免费代理不可用,换出口 |
| 220 | 宝塔 Pure-FTPd 口令面判据：存在/不存在用户同一条 530（无法枚举），只有 230 算通；字典来源与 PUT 马战术 | weak-password | script | Y | 宝塔,Pure-FTPd,MySQL虚拟用户,530,用户名枚举不可行,230判据,FTP字典,上传webshell |
| 221 | EyouCMS 前台任意文件写入 GetShell（api/Uploadify::preview 后缀绕过） | file-upload | exp |  | eyoucms,文件上传,getshell,前台未授权,Uploadify |
| 222 | EyouCMS 后台模板编辑 RCE（CVE-2026-7388，ThinkPHP {:} 标签绕过过滤） | rce | exp |  | eyoucms,后台getshell,模板写马,代码注入,thinkphp,标签绕过,FilemanagerLogic |
| 224 | EyouCMS 前台 pre-auth SQL 注入（common.php GetSortData，参数 sort_asc） | sqli | poc |  | eyoucms,sqli,前台未授权,pre-auth,GetSortData,sort_asc |
| 225 | EyouCMS 后台 Logo 上传不受限 → GetShell（CVE-2024-11211） | file-upload | exp |  | eyoucms,后台getshell,文件上传,不受限上传,Logo |
| 226 | EyouCMS 未授权任意文件读 / 文件包含致代码执行（CVE-2026-1107，Diyajax check_userinfo） | unauthorized | poc |  | eyoucms,文件包含,任意文件读,路径穿越,未授权,Diyajax,check_userinfo |
| 227 | EyouCMS 二次反序列化 RCE（CVE-2025-15375，api/Ajax.php arcpagelist 参数 attstr） | deserialization | exp |  | eyoucms,反序列化,php对象注入,thinkphp,gadget链,arcpagelist |
| 228 | EyouCMS 后台接管跳板：官方 setpwd.php 密码/密保重置工具（配合任意文件写） | auth-bypass | script |  | eyoucms,后台接管,密码重置,任意文件写,跳板,setpwd |
| 229 | trade.php 判死 + 全 NUL 成因（防护清零）+ 反推前人用的是 bank 链（真正该查 data/bbscache/bank_sort.php） | app-audit | script |  | phpwind,trade.php,全NUL,查杀清零,写码点,反推前人payload,bank_sort |
| 230 | PHPWind 7.5/<target> 根目录标准文件清单对照：trade.php 是标准，moneylog.php 与 userinfo.php 是非标准 | app-audit | script |  | phpwind,标准文件清单,非标准文件,moneylog,userinfo,版本识别 |
| 231 | EasyPay 支付回调签名走私（return_url 内嵌 trade_status 伪造支付成功） | auth-bypass | poc |  | payment,webhook,signature,parameter-smuggling,md5,go,business-logic |
| 232 | EyouCMS 后台任意文件复制 → 泄露数据库配置（CVE-2026-6561，Index::edit_adminlogo） | info-leak | poc |  | eyoucms,后台,任意文件复制,配置泄露,数据库口令,edit_adminlogo |
| 233 | EyouCMS 版本判定速查（4 判据 / 最多 4 包锁定版本区间） | recon | script |  | eyoucms,版本判定,指纹,踩点,省包,install.lock |
| 234 | EyouCMS 低价值 / 已修复条目汇总（反面清单，勿耗包数） | app-audit | template |  | eyoucms,低价值,已修复,反面清单,xss,csrf,信息泄露 |
| 235 | <target> 站群跳转机全机地图与负面情报（唯一动态程序 PHPWind <target> + 全静态页 + 三套异机真入口） | recon | template | Y | 站群地图,负面情报,宝塔,跳转机,域名池,每日域名生成器,线路选择页,泛解析,PHPWind<target>,判死 |
| 236 | 易支付回调伪造可行性判定（签名域是否含攻击者可控值）——new-api 反面 / sub2api 正面 | auth-bypass | poc |  | payment,webhook,signature,parameter-smuggling,easypay,new-api,sub2api, |
| 237 | 得分系统 G1 去重陷阱：同档位重复记分会被累加，禁止为补 step_id 关联而重记得分 | methodology | template | Y | 得分系统,G1,去重陷阱,台账纪律,负面情报,禁止重记,step关联 |
| 238 | IDOR 判定铁律：必须做"同身份换 ID"对照实验，否则不算确认（含一次真实误报复盘） | auth-bypass | template | Y | IDOR,越权,判定纪律,误报,对照实验,负面情报,phpwind |
| 239 | PHPWind 7.x 后台入口名定位：读 pw_config.db_adminfile（前台 $db_adminfile 渲染），不做文件名爆破 | info-leak | script |  | phpwind,后台入口,配置项,信息泄露,白盒审计 |
| 240 | Sub2API EasyPay 支付回调验签绕过（参数走私 / 签名串拼接不一致 → 伪造回调入账） | auth-bypass | exp | Y | 易支付,easypay,sub2api,验签绕过,参数走私,业务逻辑漏洞,伪造入账 |
| 241 | Sub2API 易支付(EasyPay) 回调验签绕过 —— 通用参数化利用脚本（零密钥伪造回调入账） | auth-bypass | exp | Y | sub2api,easypay,易支付,支付回调,验签绕过,参数走私,签名拼接不一致,零密钥,业务逻辑 |
| 242 | PHPWind hack 插件入口 hack.php?H_name=&lt;插件&gt; —— 绕过 hack/&lt;插件&gt;/index.php 的 9B Forbidden 守卫 | frontend-reverse | poc | Y | phpwind,hack插件,入口发现,readover绕过,9B守卫 |
| 243 | PHPWind 7.x pw_ajax.php?action=pcdelimg fieldname 布尔盲注（wooyun-2010-0417，影响 phpwind 7/8，7.5 SP3 已修） | sqli | poc |  | phpwind,<target>,sql注入,盲注,pcdelimg,pw_ajax,后台入口,db_adminfile,apikey |
| 244 | PHPWind 7.x pw_api.php → api/class_Other.php threadscateGory 写码 RCE（wooyun-2010-0417 第二洞，7.5 SP3/8.x 已加 intval 修复） | rce | poc |  | phpwind,<target>,rce,写码,pw_api,class_Other,threadscateGory,info_class, |
| 245 | PHPWind 7.x pw_api.php → api/class_base.php $mode 任意文件包含（CN-VA10-02 点名文件，不依赖 register_globals） | rce | poc |  | phpwind,<target>,任意文件包含,LFI,pw_api,class_base,mode,CN-VA10-02,不依赖regis |
| 246 | PHPWind <target> 情报审计负面清单 + 7.5已修/<target>未修补丁对照（tag写码/hack_rate LFI/faq/vote/apps 全部源码级证否） | app-audit | poc |  | phpwind,<target>,负面情报,证否,pw_var_export,hack/rate,faq,vote,apps,补丁差异,7. |
| 247 | PHPWind <target> 生日显示写码链（6.3.2 前台 getshell 手法）成立性判定：机制成立但被注册名 Char_cv 引号实体化阻断 | rce | poc |  | phpwind,<target>,生日缓存,brith_cache,写码,前台getshell,Char_cv,引号实体化,负面情报,6.3 |
| 248 | PHPWind 7.x 探测方法论：CC 防护 pwDefendCc 导致的「0 字节假阴性」与四条规避规范 | recon | script | Y | phpwind,方法论,CC防护,假阴性,探测规范,盲注,布尔判定 |
| 249 | PHPWind 7.x 入口 × action 全量枚举方法论（251 个 action + 10 类危险标记）+ 写点判死结论集（pw_var_export 转义 / Pcv 封穿越 / uploadicon 数学不可能 / mutiupload 三重锁） | path-traversal | script | Y | phpwind,action枚举,攻击面测绘,写点判死,pw_var_export,Pcv,P_unlink,uploadicon,muti |
| 250 | PHPWind 7.x 后台风格 CSS 编辑 → data/style/*_css.htm → updatecache_sy eval 代码注入 RCE | rce | poc |  | phpwind,后台,风格模板,代码注入,eval,${},stripslashes抵消,data/style |
| 251 | PHPWind 写 PHP 文件 sink 的三面叠加扫描法（数据来源分类 + 变量路径补扫） | rce | script | Y | phpwind,白盒审计,sink扫描,写文件,数据来源分类,变量回溯 |
| 252 | PHPWind 写 .php 缓存文件判死边界速查表（含 bank_sort 例外与字段长度表） | app-audit | payload | Y | phpwind,白盒审计,判死边界,缓存写点,bank_sort,字段长度 |
| 253 | PHPWind <target> admin/bakup.php 任意 SQL 执行（bakindata 多语句 + bakuptable 表名注入） | sqli | poc |  | phpwind,bakup,任意SQL,多语句,备份恢复,权限绕过判据 |
| 254 | PHPWind <target> 系统性加固判死清单（bank_sort RCE 链 + pw_hack 写入面 + 权限模型全量审计） | rce | poc |  | phpwind,<target>,bank_sort,rce,判死清单,权限模型,pw_hack,复用价值 |
| 255 | GBK 站点后台表单「全量回放」必须走 latin-1 字节级往返（PHPWind <target> 实测） | app-audit | script | Y | phpwind,gbk,utf-8,form-replay,encoding,latin-1,config-restore |
| 256 | PHPWind <target> 后台 Getshell 优先级手册（权限自检 148xx 判据 + P0 SQL OUTFILE → P5 附件扩展名 + 源码级判死清单） | rce | template |  | phpwind,<target>,后台getshell,操作手册,权限判据,判死清单,宝塔 |
| 257 | PHPWind <target> 后台附件白名单放开（GBK 全量表单回放）与上传失败分层判据 | info-leak | script | Y | phpwind,后台,附件白名单,全量表单回放,GBK,latin-1,attachdir,宝塔WAF,判据 |
| 258 | 出口 IP 被宝塔 WAF 拉黑后的应急出口：免费 HTTP 代理池（抓取/验活/轮换 + schannel 坑） | tunnel | script | Y | 出口IP,宝塔WAF,IP黑名单,免费代理池,轮换,schannel,ssl-no-revoke,判据 |
| 259 | 宝塔网站防火墙 <target> 拦截边界与绕过（IP 层换出口 / URL 关键词混合大小写 / multipart 双拦 PHP 特征） | waf | template | Y | 宝塔WAF,BTWAF,IP黑名单,关键词绕过,multipart,上传拦截,换出口IP,PHPWind,判据 |
| 260 | 出口 IP 与 Host 头是两个独立变量：环境代理静默改道导致宝塔 WAF 封禁被误读（curl 必须 --noproxy "*"） | waf | template | Y | 宝塔WAF,出口IP,环境变量代理,noproxy,假阴性,判据,站群,PHPWind |
| 261 | PHPWind 7.x 前台/后台 POST 一律被 Showmsg('undefined_action') 拒绝 —— 根因是 PostCheck() 要求 Referer 的 host 严格等于 HTTP_HOST（用 IP 做 Referer 是致命陷阱） | app-audit | template | Y | phpwind,referer,undefined_action,PostCheck,POST失败根因,判据,GBK,IP直连陷阱 |
| 262 | 宝塔 WAF 来源封禁实测判据：multipart 上传后全出口 000/ProxyError，而同代理到 baidu 200（区分「代理失效」vs「目标拒我」） | waf | poc | Y |  |
| 263 | 宝塔网站防火墙 <target> 拦截页指纹（403/1200）＋ SQL 关键字载荷被整拦的实测判据＋代理出口两级验活法（--noproxy 自杀陷阱） | env | template | Y | 宝塔WAF,<target>,SQL注入判据,403,1200,代理验活,--noproxy陷阱,myip,出口被封 |
| 264 | PHPWind <target>/7.5 后台写文件 getshell 三处写入端判死（Pcv 封穿越 / pw_var_export 转义数组键 / setcss 非写入端）+ data 目录"是否写入成功"的可靠判据 | path-traversal | template | Y | phpwind,后台getshell,判死,Pcv,路径穿越,setcss,PostLog,admin_record,pw_var_expo |
| 265 | SPA+Spring Cloud 网关批量未授权接口探测（/xxa 前缀型） | unauthorized | script | Y | 未授权访问,Spring Cloud,SPA,接口枚举,配置泄露 |
| 266 | 点云客服 dotclouds 访客端未授权 multipart 任意文件上传（net.dotclouds / Spring Boot） | file-upload | poc | Y | 点云客服,dotclouds,udesk,未授权,文件上传,getshell |
| 267 | 点云客服 dotclouds 会话接口零鉴权读取（IDOR，queryChatHistory/queryConvMessage） | auth-bypass | poc | Y | 点云客服,dotclouds,未授权,IDOR,会话读取 |
| 268 | dysaas 平台家族统一上传接口任意内容写入（/xxa/longm/api/v1/file/upload） | file-upload | poc | Y | 文件上传,扩展名绕过,对象存储,S3,任意文件写入,dysaas |
| 269 | dysaas 代理后台 API 前缀双写规律 + 匿名 clientCode 商户配置泄露 | unauthorized | poc | Y | 代理后台,未授权,路径前缀双写,信息泄露,dysaas,AES |
| 270 | 点云客服 dotclouds sendConvMessageFile：companyId 未 normalize 导致任意目录文件写入（路径穿越、可逃出 /data 到根文件系统） | file-upload | exp | Y | 点云客服,dotclouds,udesk,未授权,文件上传,路径穿越,任意文件写入,spring-boot |
| 271 | dysaas 代理后台：CDN 边缘 IP 轮换绕过源站封禁 + 请求格式/AES + 匿名代理注册 appSave | unauthorized | script | Y | dysaas,代理后台,CDN边缘IP轮换,绕过IP封禁,未授权,匿名代理注册,AES |
| 272 | 直播平台未授权注册绕过短信验证（/api/user/simple_regist 空 mobile_hash 直接签发 token） | auth-bypass | exp | Y |  |
| 273 | 腾讯云 STS 临时凭证 + COS q-signature 手工签名（写入/边界判定）可用脚本 | unauthorized | script | Y |  |
| 274 | 多租户直播平台前端响应加密还原（AES-256-CBC 硬编码 key/IV + 时间戳尾部拼接） | info-leak | script | Y |  |
| 275 | 匿名注册/申请接口作账号存在性 oracle（用户名枚举通用手法 + dysaas 代理 SaaS 判定表） | info-leak | script | Y | 用户枚举,存在性oracle,未授权,SpringBoot,CDN边缘直连 |
| 276 | /xxa/longm/api/v1/file/upload 未授权文件上传（multipart，无 Authorization） | file-upload | exp | Y | 未授权,文件上传,对象存储,S3,白名单绕过 |
| 277 | 密码喷洒：批量生成的运营账号共用默认口令（先证明登录无需验证码 + 按前端算法构造口令密文） | weak-password | script | Y | 密码喷洒,弱口令,批量账号,验证码绕过,CDN边缘直连 |
| 278 | Swagger 驱动的水平越权批量探测（自己/他人 id 对照 + 未文档化参数注入 + 加密响应长度 oracle） | auth-bypass | poc | Y | idor,swagger,openapi,越权,加密响应,oracle |
| 279 | 多租户 SaaS 后台「已证否方向」坑位卡（上传无解析 / detail 非 IDOR / Bearer 前缀与单会话顶号 / 验证码漏校验 / 错误码语义表） | methodology | template | Y | 负向结论,坑位卡,多租户,SpringBoot,对象存储,单会话,Bearer |
| 280 | RN/Expo(Hermes) APK 逆向：提取后端签名密钥与加密算法 | info-leak | script | Y | android,apk,react-native,hermes,逆向,签名密钥,硬编码密钥 |
| 281 | CN7 /jumpto/img SSRF 内网/云元数据不可达（边缘代理阻断，实测负面） | ssrf | poc | Y | ssrf,CN7,图片代理,内网不可达,负面 |
| 282 | 博彩站群(<target>:10666) Java 反序列化/JNDI/Log4j RCE 全否 — 严格 Jackson 非 Fastjson + /ecb 加密网关可达真实后端 | deserialization | poc | Y | fastjson,Jackson,反序列化,JNDI,Log4j2,actuator,ecb,若依,SpringBoot,负面结论 |
| 283 | CN7/gameplat 10666 管理组件面（actuator/druid/swagger/若依）全否 — SPA 兜底 + Spring 全局鉴权过滤器 + 版本泄露 5.3.0-RELEASE | info-leak | poc | Y |  |
| 284 | CN7/gameplat 5.3.0-RELEASE 无公开版本级 Nday + 同产品实例全为静态 CDN 克隆 + nginx ^/api 前缀路由测绘 — 版本指纹 Nday 匹配判死清单 | info-leak | poc | Y | CN7,gameplat,博彩CMS,版本指纹,Nday,nignx路由,多租户 |
| 285 | 博彩 CMS「构建哈希」指纹：body="&lt;buildhash&gt;" 在 FOFA 一键枚举同产品同构建部署舰队（含每租户一端口划分） | info-leak | poc | Y | 指纹,FOFA,构建哈希,博彩CMS,资产测绘,CN7,gameplat |
| 286 | 博彩站群 <target>:10666 Java 反序列化/JNDI/Log4j2 全否（含 OOB 有效性复核与 /ecb 网关格式） | deserialization | markdown | Y | negative,deserialization,jndi,log4j,fastjson,jackson,spring-boot,ecb,w |
| 287 | CN7/gameplat 多租户SaaS：未授权 /api/tenant/setting/getPersonalCenter 跨租户泄露（含跨租户运营账号） | info-leak | poc | Y | 多租户,cross-tenant,unauth,info-leak,CN7,gameplat,SaaS |
| 288 | CN7/gameplat 多租户博彩 SaaS 未授权跨租户配置泄露 (getPersonalCenter/getConfig) | info-leak | http | Y | multi-tenant,broken-tenant-isolation,unauthorized,CN7,gameplat,info-le |
| 289 | CN7/gameplat /api/config/promote/* 是推广跳转(catch-all 302)而非受门控的配置管理组——无 host/租户绕过面 | auth-bypass | poc | Y |  |
| 290 | Web 缓存投毒/缓存欺骗面一站式测绘脚本（含缓存键归一化与 unkeyed 头判定） | recon | script | Y | cache-poisoning,web-cache-deception,host-header,smuggling,http-proxy,c |
| 291 | CDNGreat 缓存/防护层指纹识别与缓存范围、缓存键判定方法 | recon | script | Y | cdngreat,cache,fingerprint,waf,host-header,cl-cache-status |
| 292 | fakerToken 服务无鉴权探测：/fbsport/v1/getToken 与 /ob/v2/getTryToken 泄漏后端域名与 token | info-leak | script | Y | fakerToken,guest token,未授权,信息泄露,requestData,form编码,线路域名路由表,白标供应商,体育平台, |
| 293 | 对象存储未授权列举/写入判定（含 SPA 回落假阳性排除 + 写探针回滚） | unauthorized | script | Y | 对象存储,S3,OSS,未授权,桶列举,ListBucketResult,写探针,假阳性排除,SPA回落,静态托管 |
| 294 | Vite/Vue3 minify SPA 接口调用点与请求体还原器（三层别名映射） | frontend-reverse | script | Y | Vite,Vue3,minify,混淆,别名映射,三层映射,接口还原,请求体还原,前端逆向,SPA,静态分析 |
| 295 | 用户可控字段 → 无鉴权公开接口渲染 → 面向全站访客的投毒链（四步验证 + 如实边界） | unauthorized | script | Y | 投毒链,头像外链,未授权渲染,公开榜单,内容注入,XSS判定,清洗器,safe-url,博彩SPA |
| 296 | 前端硬编码 requestData 常量重放 → 未授权签发第三方平台 token 与线路域名路由表 | info-leak | script | Y | 硬编码常量,requestData,重放,guest token,未授权,fakerToken,线路域名路由表,白标供应商,体育平台 |
| 297 | fakerToken 硬编码 requestData 常量重放：未认证签发 OB 平台登录 token 与线路后端域名表（RB88/OB 站群通用） | unauthorized | poc | Y | fakerToken,rb88,博彩站群,未授权,token签发,硬编码常量,requestData,信息泄露 |
| 298 | Vite SPA 资金接口静态还原 + 二维码渲染定性（含带会话的页面内 fetch 验证） | info-leak | script | Y |  |
| 299 | 统一 20000 未登录拦截的绕过穷尽判定法（含 body×CT 归因控制实验，防 116 误报） | auth-bypass | script | Y | auth-bypass,拦截器绕过,20000,Content-Type,白名单前缀,对照实验,防误报 |
| 300 | cmwall「过移动墙」服务管理面板 <target> 未授权读写（38888）+ 零风险写验证流程 | auth-bypass | poc | Y | cmwall,过移动墙,38888,未授权,面板,流量劫持,红队 |
| 301 | 前端资源投毒可行性判定法：客户端渲染资金字段 + 无 SRI/无 CSP/资源域写方法（含正确分级口径） | frontend-reverse | script | Y | 前端投毒,SRI,CSP,资源域,CORS,二维码,资金链路,可行性判定 |
| 302 | 支付回调/通知端点探测法：静态零引用判定 + 未知路由三态判别 + 按站点命名风格派生字典（含负结论口径） | auth-bypass | script | Y | 支付回调,notify,伪造到账,路由三态判别,前端零引用,字典探测,负结论口径 |
| 303 | fakerToken/OB 的 api= 参数解密（硬编码 AES-ECB 密钥 OBTY20220712OBTY）与 requestData 不可解结论 | info-leak | script | Y | 硬编码密钥,AES-ECB,fakerToken,OB,api参数解密,离线,bundle审计 |
| 304 | cmwall「过移动墙」wizard 面板 静态分析速查 + 命令注入实现基础（含厂商公开包/部署路径/注入验证要点） | rce | script | Y | cmwall,<target>,wizard,过移动墙,命令注入,UPX,面板,红队 |
| 305 | WAF 绕过/能力修正：直连 HTTPS 清空 SNI（+任意 Host）→ 200；带真实 SNI（<target>）→ RST（按 SNI 丢包） | waf | script | Y | WAF绕过,SNI,直连,TLS,站群,API,能力修正 |
| 306 | 加固 APK（zxprotect 抽取壳）零依赖逆向 + config2.json 配置引导链提取 | info-leak | script | Y | android,apk,加固,抽取壳,zxprotect,dex切分,config引导链,C段指纹比对,信息收集 |
| 307 | LHZS/LHService v27 网关免签名面定位 + 两层签名校验拒绝码判据 | unauthorized | poc | Y | LHZS,网关,免签名,未授权,签名校验,拒绝码判据,routing-health,v27 |
| 308 | 原生 socket 端口普查器（绕代理假象）tcpscan.py | recon | poc | Y | 扫端口,原生socket,代理绕过,段普查,recon |
| 309 | LHZS/靶标节点签名探测器（严格判据）lhprobe.py | recon | poc | Y | 资产归属,严格判据,免签名接口,签名指纹,recon,LHZS |
| 310 | TLS 证书抓取 + 从 DER 提 CN/SAN（证书归属判定）certgrab.py | info-leak | poc | Y | 证书,归属判定,SAN,DER,无代理,tls |
| 311 | 多目标业务路径 + Host(vhost) 覆盖探测器 sweep.py | unauthorized | poc | Y | vhost,Host头,业务路径,归属判定,recon |
| 312 | 非标准服务四态识别器 rawprobe.py | recon | poc | Y | 服务识别,黑洞端口,协议判定,recon |
| 313 | 端口 banner 抓取器（自动判 SSH/MySQL/HTTP）banner.py | recon | poc | Y | banner,服务识别,SSH,MySQL,版本指纹,recon |
| 314 | Web 指纹探测（标题 + LHZS 特征签名）httpfp.py | recon | poc | Y | Web指纹,标题,UTF8,中文判定,recon |
| 315 | 单目标业务路径 + vhost 探测（先看正文再判归属）bizprobe.py | unauthorized | poc | Y | 业务路径,vhost,正文判读,归属判定,recon |
| 316 | DNS 解析 + 指定 SNI/Host 的 HTTPS 定向探测 sni72.py | recon | poc | Y | SNI,DNS,归属判定,vhost,历史节点 |
| 317 | nginx→Tomcat `..;/` 归一化差异：穿 location 白名单 + 跨 context 横向（含绕开 nginx 死上游 504）；Manager IP 阀 33 变体实测无效 | path-traversal | poc | Y | nginx,location绕过,路径归一化,..;,跨context,Tomcat,Manager,RemoteAddrValve,负结论 |
| 318 | Tomcat 拿 shell 路线的快速判死：PUT 405 + OPTIONS 无 PUT ⇒ CVE-2017-12615/12617/2025-24813 全灭；WEB-INF 25 种写法全 404；Tribes RCE 与 Ghostcat 版本不匹配 | rce | poc | Y | tomcat,负结论,readonly,PUT,405,CVE-2017-12615,CVE-2017-12617,CVE-2025-248 |
| 319 | 全量参数钥匙：批量候选参数名越过老 servlet 的 Integer.parseInt 门槛（分组探测 + 二分收敛，通用可复用） | unauthorized | script | Y | servlet,未授权访问,参数爆破,老代码,Tomcat,六合 |
| 320 | Tomcat Manager 快速判死：403 页字节级指纹（同版本 stock 页跨机器 sha1 完全一致）+ 端口面普查（只暴露 nginx ⇒ 头伪造必然无效，Tribes/8009 一并判死） | recon | poc | Y | tomcat,manager,RemoteAddrValve,负结论,stock403指纹,端口面,CVE-2026-34486,Ghost |
| 321 | 盲 SSRF 判定法：差分三值（死地址基线/内网/外带）+ A/B/A 交替复测识破抖动 —— 附"零差异=无 SSRF"实测样例与真 SSRF 对照特征 | ssrf | script | Y | ssrf,盲ssrf,判定法,差分,A/B/A复测,网络抖动,负结论,老servlet |
| 322 | 六助(LHZS) App <target> 客户端接口字段名与请求体结构全表 | protocol | payload |  |  |
| 323 | LHZS v27 管理台：认证链精确字段 + admin/query 信封 + release.publish SSRF + storage.configure 落盘可控 + action 全表 | auth-bypass | payload |  |  |
| 324 | 【技巧】Java 异常行号 oracle：用 Tomcat 错误页的 Class.java:NN 行号推进，零猜测定位 servlet 参数与鉴权门槛 | methodology | script | Y | java,servlet,参数发现,异常行号,oracle,judge,鉴权判定 |
| 325 | 【技巧】「参数钥匙串」批量命中老 Java servlet 数据分支：390 个候选参数各 =1，一次请求绕过猜参数名（单请求拿到 143KB/408 条用户记录） | methodology | script | Y | java,servlet,参数钥匙串,NumberFormatException,未认证,批量参数,老系统 |
| 326 | 【漏洞】nginx→Tomcat 归一化穿越（..;/）：绕开 nginx location 白名单，未认证直达 Tomcat ROOT/docs/examples 并确认 Manager/host-manager 存在（IP 阀纯头伪造绕不过） | path-traversal | poc | Y | nginx,tomcat,path-traversal,归一化,manager,RemoteAddrValve,绕过 |
| 327 | 原生 socket / 手工 TLS 探测骨架（绕本机 HTTP 代理污染 + binary_form 证书取证 + 定制404/Spring404/406 三态判存） | env | poc | Y | 原生socket,手写HTTP,绕本机代理,SNI=IP,证书取证,指纹判存,chunked |
| 328 | 老式 Java servlet「分块灌参 + 栈行号二分收敛」参数发现法（Tomcat 报错栈当 oracle） | info-leak | script | Y | java,servlet,参数爆破,tomcat,栈泄露,delta,二分 |
| 329 | UpLoadPic 三个入口对照表（upload2=唯一真写点/写 OSS；无扩展名分支=只回 URL 模板不落盘；upload 的 head 四字段=真字段但空响应、token 模式=无条件假成功） | file-upload | poc | Y | uploadpic,upload2,未认证文件上传,OSS,假成功,head字段,负结论,判据,LHZS |
| 330 | 【方法·负面知识】四步判死「无条件假成功」写端点 + 参数级 delta 证明：UploadPicServlet 只读 userid 一个参数、无文件名/路径可控、无文件落地 | methodology | script | Y | false-positive,假成功,写端点判定,UploadPicServlet,jar,参数级证明,delta,negative-kno |
| 331 | 行号 oracle 参数隔离法（附 PutTp 判死实测）：用 Tomcat 错误页的 `类名.java:行号` 逐级隔离参数、二分求最小集、判 POST 是否独立分支 | methodology | poc | Y | 老servlet,Tomcat,错误页,行号oracle,参数隔离,最小参数集,doPost包装,负结论,PutTp |
| 332 | 【技巧】把老 servlet 入参面压到 1 个：最小必需参数集二分 + 行号 oracle 定位 parseInt 参数 + 用「-操作过于频繁」当分支存在判据（实测 Login4 收敛到 flag 一个参数） | auth-bypass | script | Y | java,servlet,flag分支,最小参数集,二分,行号oracle,值域枚举,限流判据,登录入口 |
| 333 | 【方法论·五步判据】老 Java servlet 读写判别：参数要补齐（防假否定）+ 读回要跨时间窗（防缓存假否定）+ 值域枚举看文案推进而非限流触发（防假肯定）+ 探针换假号（防锁桶） | methodology | script | Y | java,servlet,判据,假否定,假肯定,Redis缓存,写后读回,限流,值域枚举,flag分支,老系统 |
| 334 | 【oracle】判定「验证码到底有没有真的发出去」：用错码探「用码」分支，「请先获取验证码」vs「验证码错误」两句文案即判据（零成本、不消耗真码、无需目标手机） | unauthorized | script | Y | 短信验证码,发码判定,oracle,状态机,存根,假成功,参数名对照,限流保护 |
| 335 | 复现冲突四步法：纯读稳定性 → 裸/戴钥匙对照 → 逐个剔除定最小集 → 跨时间窗读回（写侧+读侧都做对才不误判） | methodology | poc | Y | 未认证写,复现冲突,最小可复现请求,参数钥匙,Redis缓存,写后读回,判据口径,假成功 |
| 336 | jadx 卡死时的替代：最小 DEX 解析器 + 字符串表邻接/簇分析（无需反编译） | frontend-reverse | script | Y | android,apk,dex,静态分析,jadx替代,字符串邻接,不反编译 |
| 337 | 如何证明"我的会话确实生效了"：同一端点有无会话的响应差异自证法（400 动作级 vs 403 路径级） | methodology | script | Y | 会话自证,假阴性排除,方法论,黑盒验证,错误码分级 |
| 338 | Go 自研代理服务（管理口 + 代理口）免凭据探测与绕过穷尽核查脚本 | unauthorized | script | Y | go,代理面板,basic认证,open-proxy,未授权,指纹,穷尽测试 |
| 339 | FastAdmin CVE-2024-7928 任意文件读取（/ajax/lang lang 参数） | path-traversal | poc | Y | FastAdmin,ThinkPHP,任意文件读取,path-traversal,CVE-2024-7928 |
| 340 | PHP .user.ini 两用技巧：追加 open_basedir=/ 破除目录限制 + auto_prepend_file 绕过 .php 写入黑名单落马（冰蝎 v4 双通道） | rce | script | Y | php,user.ini,auto_prepend_file,open_basedir,绕过黑名单,冰蝎,落马,宝塔 |
| 341 | PHP .user.ini 覆盖导致 open_basedir 逃逸（可写子目录 → 整机文件系统访问） | privesc | payload | Y | php,open_basedir,user.ini,逃逸,宝塔,phpwind,提权 |
| 342 | 宝塔防篡改只锁 .php + 替代后缀不被解析的实测判据（.php5/.phtml/.inc 换后缀拿不到 webshell） | app-audit | template | Y | 宝塔,防篡改,php后缀,nginx,405,phtml,php5,webshell判死,phpwind |
| 343 | boyue(webman) 博彩站群会员端 —— upload-qrcode 任意扩展名写入 + 同源存储型 XSS（WithdrawController::uploadQrCode） | file-upload | exp | Y | webman,boyue,博彩站群,任意文件上传,同源XSS,file-upload |
| 344 | boyue(webman+Pear Admin) 博彩站群后台指纹 / 路由表 / 密钥 / 免登录路由 速查手册 | info-leak | script | Y | webman,boyue,pear-admin,博彩站群,gofly,指纹,路由表,免登录 |
| 345 | 本机 fake-ip DNS + 规则型本地代理下，用 CONNECT+自定义 SNI 直达 nginx vhost（绕过假性「目标封 IP」） | tunnel | script | Y | fake-ip,DNS劫持,SNI,CONNECT,代理,vhost,出口 |
| 346 | webman 插件路由指向非插件命名空间控制器 → AccessControl 中间件失效导致后台接口未授权 | unauthorized | script | Y | webman,PearAdmin,未授权访问,中间件失效,路由命名空间,后台接口 |
| 347 | webman 路由大小写不敏感 → 绕过按路径控制器段精确匹配的「功能关闭」开关（可救回代码生成器写文件 → Monitor 热重载 RCE） | auth-bypass | script | Y | webman,workerman,路由,大小写绕过,WAF绕过,功能开关,代码生成器,file_put_contents,热重载,RCE |
| 348 | webman 默认(自动)路由绕过路由级 AuthMiddleware —— 未授权直达控制器方法（通用） | auth-bypass | poc | Y | webman,默认路由,自动路由,鉴权绕过,未授权访问,auth-bypass,php,workerman |
| 349 | webman-admin TableController::crud 代码生成器 docblock 注入 → 写 app/controller/*.php → Monitor 热重载 RCE | rce | exp |  | webman,webman-admin,pear-admin,代码生成器,docblock注入,rce,php |
| 350 | webman/Pear Admin 目标「上传→执行」全判定手册（include/LFI/phar/模板 四路排除 + 扩展名名单测绘） | path-traversal | script | Y | webman,workerman,pear-admin,php,code-audit,lfj,lfi,phar,upload,rce,ora |
| 351 | webman 响应加密（X-Encrypted: 1）硬编码 AES 密钥解密 / 加密请求体构造 | frontend-reverse | script | Y | webman,响应加密,X-Encrypted,AES-256-CBC,硬编码密钥,固定IV |
| 352 | webman 自动路由绕过鉴权中间件（/{controller}/{action} 零凭证直连） | auth-bypass | script | Y | webman,自动路由,鉴权绕过,AuthMiddleware,未授权访问,RateLimit绕过 |
| 353 | webman 自动路由作为「方法存在性 oracle」：零凭证枚举部署版控制器方法（区分零鉴权直出 / 需登录 / 只判角色） | unauthorized | script | Y | webman,自动路由,未授权,方法枚举,oracle,鉴权测绘 |
| 354 | webman 自动路由未授权接口普查 + Crypto 硬编码密钥解密响应 | info-leak | script | Y | webman,未授权,自动路由,信息泄露,AES硬编码 |
| 355 | nginx vhost oracle 枚举（Host 头区分真实虚拟主机，识别同名后端多品牌） | recon | script | Y | nginx,vhost,信息收集,Host头,多品牌 |
| 356 | 跨 vhost 同源 XSS 可执行性判定矩阵（多域名站群：/uploads 可达性 × /app/admin 交叉对照） | unauthorized | script | Y | 同源判定,vhost,多域名站群,存储型XSS,web root,跨域会话 |
| 357 | webman-admin(Vue Art Design Pro) 后台会话与请求加密规范 + 403 开关大小写绕过 + crud docblock 注入开火参数 | auth-bypass | script | Y | webman-admin,ArtDesignPro,pinia,localStorage,AES-256-CBC,X-Encrypted,4 |
| 358 | boyue/PearAdmin 充值渠道链上收款地址未授权替换 POC（payset-update / configs_data[trc20]） | unauthorized | exp |  | webman,PearAdmin,未授权写入,充值渠道,链上地址,资金流,boyue |
| 359 | webman「显式路由 vs 自动路由」双姿势鉴权/闸门 oracle（可复用：一次判定整族方法被哪一层拦） | auth-bypass | script | Y | webman,自动路由,鉴权绕过,交付版闸门,oracle,hch7777 |
| 360 | webman 自动路由「形参绑定」未授权身份/ID 注入（$request->all() → 控制器形参） | auth-bypass | poc | Y | webman,php,参数绑定,未授权,身份注入,IDOR |
| 361 | webman「形参名 oracle」：用 MissingInput/InputType 错误页在线反推控制器方法签名 | info-leak | script | Y | webman,php,方法签名枚举,oracle,信息泄露,参数绑定 |
| 362 | SSRF 原语筛查法：从形参到外发 URL 的四步数据流溯源（含二阶配置写入者追踪） | ssrf | script | Y | ssrf,url参数,数据流溯源,白盒筛查,guzzle,二阶ssrf |
| 363 | 零副作用验证判据套件：身份注入 / 原密码强校验 / 验证码校验（用不存在 uid 与现值回写） | auth-bypass | script | Y | 判据,零副作用,身份注入,验证码,改密,风控 |
| 364 | 路径穿越有效性判定：归一化等价对照法（nginx alias 前缀映射场景） | path-traversal | script | Y | nginx,alias,路径穿越,%2f,归一化,判定方法 |
| 365 | 任意后缀上传 → 同源存储型 XSS → localStorage JWT 会话窃取 → 无限速资金密码 oracle（账号接管/盗资金组合链） | file-upload | exp | Y | 任意文件上传,存储型XSS,localStorage,JWT窃取,资金密码爆破,oracle无限试 |
| 366 | GoFly 客服系统鉴权模型与越权判定（/2/ 族 401vs403 判据 + ws_kefu/ws_visitor 子协议票据） | unauthorized | script | Y | gofly,goflylivechat,websocket,子协议鉴权,IDOR,客服系统,未授权 |
| 367 | 上游/供应链域名「可控性」判定流水线（whois+RDAP+DNS+HTTP+CT+TLS+家族挖掘） | recon | script | Y | 供应链,域名接管,whois,rdap,ct,recon |
| 368 | suo5 classic 隧道断链自愈 + PHP<target> 冰蝎马 openssl_decrypt("AES128") 返回 false 导致 200 空响应 | tunnel | script | Y | suo5,socks5,隧道,PHP<target>,冰蝎,openssl,运维 |
| 369 | ThinkPHP3 后台 Cookie 会话伪造（admin_user=username + admin_id=password哈希） | auth-bypass | poc | Y | thinkphp,会话伪造,cookie,auth-bypass,后台 |
| 370 | ThinkPHP3 后台 Api::article_save 任意文件写入（data:// 伪协议）→ WebShell | file-upload | exp | Y | thinkphp,任意文件写入,data协议,getshell,rce |
| 371 | 冰蝎马 PHP 5.6 适配：openssl_decrypt("AES128") 失效须改 "aes-128-ecb" | rce | payload | Y | 冰蝎,behinder,webshell,php5.6,openssl,AES128,aes-128-ecb,踩坑 |
| 372 | ThinkPHP <target> 站群 CMS：extends Action 的 Admin 控制器绕鉴权 + 微信签名可自算 → 未授权文件写入（Wechat::index_1） | auth-bypass | exp | Y | thinkphp,thinkphp<target>,鉴权绕过,extends Action,微信签名,任意文件写入,站点群,kindedit |
| 373 | PHP-FPM FastCGI 未授权访问 → PHP_ADMIN_VALUE 解除 open_basedir / 清空 disable_functions（纯 PHP 客户端，无需 shell） | info-leak | exp | Y | php-fpm,FastCGI,未授权,open_basedir,disable_functions,宝塔,unix socket,横向通用 |
| 374 | ThinkPHP <target> 站群 CMS：Home/Order::get_html1 零凭证覆盖模板文件（App/Tpl/Home/Index_index.html）→ 模板注入前置原语 | unauthorized | poc | Y | thinkphp,thinkphp<target>,未授权,模板注入,任意文件写入,模板覆盖,get_html1,extends Actio |
| 375 | ThinkPHP <target> 「未授权模板写入 → 编译执行」RCE 链（Home/Order::get_html1 → App/Tpl → App/Runtime/Cache） | rce | exp | Y | thinkphp,<target>,template-injection,rce,unauth,站群,abrdns,php |
| 376 | Admin/Wechat 免鉴权签名绕过（token 空串伪造 sha1 签名） | auth-bypass | poc | Y | wechat,签名绕过,未授权,站群,signature-bypass,thinkphp |
| 377 | ThinkPHP <target> 站群 CMS「模板烘焙」RCE 链完整打法 + 两个否证 oracle（get_html1 零凭证覆盖模板 / abrdns 抓取源陷阱 / Public/log 无回显失败判定） | rce | script | Y | thinkphp,thinkphp<target>,模板注入,模板烘焙,rce,get_html1,TMPL_DENY_PHP,parseK |
| 378 | 冰蝎马 Python 自动化客户端（PHP <target> 需把 cipher 改成 aes-128-ecb） | tunnel | script | Y | behinder,冰蝎,AES-128-ECB,PHP5.6,webshell,自动化,客户端 |
| 379 | ThinkPHP 3：VAR_FILTERS=filter_exp 让 `exp` 数组注入全线失效（静默丢字段、无痕）——判定向量 + 无效绕过清单 | auth-bypass | poc | Y | thinkphp,thinkphp3,VAR_FILTERS,filter_exp,exp数组注入,parseSet,静默丢弃,否证,判定向 |
| 380 | webman+nginx 站群：/app/admin/** 被 nginx 静态 location 遮蔽——路径/头/方法/差分编码/vhost 全不可绕（阴性） | unauthorized | markdown | Y | nginx,webman,vhost-admin,path-bypass,negative,static-location |
| 381 | 站群通用：article.ename = 创建者 admin 用户名 + Home/Index::bbs?post=&lt;id&gt; 零凭证渲染点（枚举后台账号名 / 内容可控渲染） | info-leak | payload | Y | 站群,thinkphp<target>,零凭证,渲染点,账号枚举,article_save,get_html1 |
| 382 | UEditor PHP 零鉴权面盲探（config 泄露 / listimage 路径枚举 / upload 目录存在性 / catchimage SSRF） | unauthorized | poc | Y | ueditor,unauthorized,file-upload,mkdir,ssrf,thinkphp |
| 383 | ThinkPHP<target> 站群 CMS 通杀 RCE：admin 表哈希即凭据 → Home/Article::edit_cont 写 article.title → get_html1 烘焙模板 → 编译 <?php 执行 | rce | poc | Y | ThinkPHP<target>,站群,RCE,模板烘焙,article.title,get_html1,open_basedir,Host |
| 384 | ThinkPHP<target> 站群通用链：后台 mass assignment 写 article.title → 零凭证 get_html1 模板烘焙 → 首页 RCE | rce | exp | Y | thinkphp,站群,模板注入,模板烘焙,mass-assignment,get_html1,凭据复用,rce,全站群通用 |
| 385 | PHP-FPM unix socket 未授权 + PHP_ADMIN_VALUE 清空 open_basedir/disable_functions → 全盘任意文件读写（宝塔环境通用） | privesc | exp | Y | php-fpm,fastcgi,open_basedir,宝塔,bypass,任意文件读写 |
| 386 | 对手检测能力画像：站群运营方 ThinkPHP3 后门自查脚本 jnbgg.sh（= 我们的规避指南） | recon | script | Y | 对手画像,检测规避,IOC,ThinkPHP3,站群,应急处置 |
| 387 | php-fpm FastCGI 沙箱逃逸：PHP_ADMIN_VALUE 清零 open_basedir 与 disable_functions（宝塔环境通用） | privesc | exp | Y | php-fpm,FastCGI,open_basedir绕过,disable_functions绕过,宝塔,aapanel,提权,沙箱逃逸, |
| 388 | 站群兄弟节点发现：从目标自己的 vhost 访问日志 Referer 反查（日志分块读 + 一条请求判定同站群） | recon | script | Y | 站群,信息收集,Referer,vhost日志,兄弟节点,DNS,指纹 |
| 389 | php-fpm FastCGI sendmail_path 劫持 → 绕过 disable_functions 拿命令执行（mail() 的 C 层 popen） | privesc | exp | Y | php-fpm,fastcgi,disable_functions,mail,sendmail_path,rce,linux |
| 390 | 熊猫体育 TY SDK /yewu40 签名通道无防重放 + checkId 不校验（试玩 token 实测） | auth-bypass | template | Y | TY SDK,熊猫体育,白标体育,签名通道,防重放缺失,checkId,试玩token,<target>,ue3238 |
| 391 | VerneMQ-over-WebSocket 匿名消息总线快速评估（无鉴权 CONNECT/SUB/PUB + 直连要求 + 空闲判定） | unauthorized | script | Y | mqtt,vernemq,websocket,anonymous,unauthorized,message-bus,recon |
| 392 | PHPWind 7.x/8.x 版本指纹识别（robots.txt + meta generator + 路径特征） | info-leak | poc | Y |  |
| 393 | FOFA 官方接口可用性切换（<target> 失效 → <target>） | recon | script | Y |  |
| 394 | 宝塔 WAF 对备份/源码文件扩展名的拦截：本站无「保留路径」的可用绕过（尾斜杠/%3f 只改路径，读不到文件） | info-leak | template | Y | 宝塔WAF,BTWAF,403-1200,备份文件,源码泄露,尾斜杠绕过,%3f,WAF绕过,db_hash,PHPWind |
| 395 | PHPWind 站点精确版本与数据库表前缀的无认证识别 | info-leak | poc | Y | phpwind,指纹,版本识别,表前缀,信息泄露 |
| 396 | 宝塔 per-source 封禁下的新出口：免费 HTTP 代理池可 CONNECT 任意端口(含 FTP:21)绕开源级拉黑 | tunnel | poc | Y | 宝塔,per-source封禁,出口IP,免费代理池,FTP,CONNECT任意端口,绕WAF,新出口 |
| 397 | PHPWind <target> 后台 admin.php 锁定被 X-Forwarded-For 绕过（db_xforwardip=1 时）+ per-source CC 封禁才是爆破真正瓶颈 | weak-password | poc | Y | phpwind,<target>,后台爆破,XFF,锁定绕过,per-source CC,db_xforwardip |
| 398 | 宝塔 WAF 按真实源IP CC 封禁下，用免费HTTP代理池逐请求轮换出口放大后台口令喷洒吞吐（含判据陷阱） | weak-password | script | Y | 宝塔WAF,CC封禁,出口IP,免费代理池,口令喷洒,XForwardedFor,PHPWind,admin.php,判据陷阱 |
| 399 | WordPress/WooCommerce 匿名攻击面清单（版本泄露+REST未授权+Store API非授权写入+XML-RPC爆破） | info-leak | script | Y | wordpress,woocommerce,unauthorized,info-leak,brute-force,rest-api,xmlr |
| 400 | 博彩 PHPWind 站群实战结论集（<target>/24）—— 判死路线 + 8 条可复用判据 | env | script | Y | 博彩站群,PHPWind<target>,SafeLineWAF,宝塔,判死结论,判据集,站群分群,出口被封,lgt,GBK |
| 401 | SSH 批量作业「banner 抖动 vs 真防护」判据 —— 单发复验可救回被误判的整批目标 | rce | script | Y | ssh,bruteforce,误判防护,fail2ban,抖动复验,paramiko |
| 402 | FastAPI/OpenAPI 未授权测绘 —— openapi.json 泄露全端点 + 按 security 标记自动分类实测 | info-leak | script | Y | fastapi,openapi,unauthorized,info-leak,swagger |
| 403 | ThinkPHP5 /kj 开奖数据接口未授权读取 | unauthorized | poc | Y | thinkphp5,unauthorized,api,json,博彩 |
| 404 | 自研BBS 站群 JSON API 硬编码密钥 + 接口协议逆向 | unauthorized | poc | Y | 自研cms,api,硬编码密钥,注册接口,站群 |
| 405 | PHPWind 后台登录判据纠正 —— 失败态=14540B+页内计数文案；失败也下发 cookie；必须带 cookie 复访验证 | auth-bypass | script | Y | phpwind,后台,登录判据,假阳性,gbk,尝试计数器 |
| 406 | 页面源码挖自动生成 vhost —— 从 HTML/JS 提取域名并批量解析（补字典爆破盲区） | info-leak | script | Y | recon,vhost,页面源码,自动生成域名,subdomain,方法论 |
| 407 | MinIO Console 两段式登录逆向 + Access Key 用户名枚举 | info-leak | script | Y | minio,console,用户枚举,用户名枚举,认证,两段式登录 |
| 408 | TLS 证书 SAN 反查域名（纯 IP 段挖真实业务域名 + vhost 路由） | info-leak | script | Y | 信息收集,证书,域名挖掘,vhost,TLS,SAN,纯IP目标 |
| 409 | AJ-Captcha 验证码协议复现 + secretKey 明文下发利用（芋道/RuoYi 系列） | auth-bypass | poc | Y | AJ-Captcha,anji-plus,验证码,芋道,yudao,RuoYi,secretKey,AES,协议逆向 |
| 410 | TLS 证书 SAN 反查域名 + 同 IP vhost 路由发现（纯 IP 段挖真实域名，避免错过 target） | info-leak | script | Y | 信息泄露,TLS证书,SAN,vhost发现,域名挖掘,纯IP段,内网服务发现,可复用 |
| 411 | 母舰域名编号枚举 —— 一次画出攻击者全部节点（56 节点/10 网段实测） | info-leak | script | Y | recon,dns,母舰,基础设施测绘,编号枚举,资产发现 |
| 412 | SPA 后台突破 —— 从前端 api.js 提取全接口清单 + 定向弱口令（实测 admin/admin123 拿到超管） | weak-password | script | Y | spa,react,api提取,弱口令,未授权,后台突破 |
| 413 | Synology DSM 未授权侦察 —— 两端点泄露 API 清单 / 主机名 / 公网状态 / 全部已装套件 | info-leak | poc | Y | synology,DSM,未授权,信息泄露,套件枚举,攻击面,群晖,侦察 |
| 414 | 方法论纠错：高速 SYN 扫描在 SYN-flood 防护设备上产生海量假开放端口（必须低速复核） | recon | script | Y | nmap,扫描,假阳性,方法论,SYN-flood,群晖,端口扫描,复核 |
| 415 | 群晖 DSM 套件前端 bundle 未鉴权扒取 + API/协议逆向（含容器终端 WebSocket 协议） | info-leak | script | Y | synology,DSM,群晖,套件,静态资源,未鉴权,JS逆向,API清单,WebSocket,攻击面 |
| 416 | 群晖 DSM code 120 判读纠错 —— 参数校验先于权限，不可误判为未授权 | methodology | script | Y | 群晖,DSM,错误码,判读,未授权,误报,参数校验,测试方法论 |
| 417 | webman 站群(无PHP-FPM) upload-qrcode 上传→RCE 判死：/uploads 的 .php 硬403且不可绕、替代后缀只回源码、全 vhost 无解释器 | file-upload | poc | Y | webman,file-upload,rce,negative,nginx,no-fpm,php后缀 |
| 418 | webman「写入 Monitor 监听目录→热重载 RCE」判死：落盘路径不可控 + 插件 admin 命名空间自动路由不可达 + /app/admin 被 nginx 遮蔽（hch7777 站群） | rce | poc | Y | webman,workerman,monitor,热重载,file-write,自动路由,命名空间,nginx,negative |
| 419 | shopnn(<target>) mock 登录认证绕过：alipay-mock 按 phone 无口令接管任意已知手机号账号 | auth-bypass | poc | Y | auth-bypass,account-takeover,mock-login,node,express,jwt,shopnn |
| 420 | shopnn (<target>) 可利用性判死集：mock 登录回显 MySQL 错误（=信息泄露）；img-proxy SSRF 无内网/写面绕过 | info-leak | poc | Y | info-leak,mysql-error,express,mock-login,ssrf,blocked_url,redirect,sel |
| 422 | GoFly 客服附件卡片 safeChatAssetPath(%2e%2e) 绕过 → 同源存储型 XSS 主动投递链 | file-upload | script | Y | gofly,kefu,xss,stored-xss,file-upload,sanitizer-bypass,phishing,same-o |
| 423 | SafeLine/宝塔站群 WAF 精确触发签名（union select 需带操作数、相邻 ../、sleep() 括号）与全站路径无豁免实测判据 | sqli | poc | Y | waf,safeline,宝塔,phpwind,sqli绕过,全站拦截,union select,路径穿越,代理池 |
| 424 | PHPDisk <target> adminlogin 账号枚举 + bind_user 无频控口令校验双通道 | auth-bypass | script | Y | phpdisk,auth-bypass,user-enum,brute-force,php |
| 425 | PHPDisk 6.x adminlogin 注入面三点连测 + 后台 gid=1 校验 | auth-bypass | script | Y | phpdisk,sqli-triage,dosafe,padding-bypass,auth,admincp |
| 426 | New API 定制构建「全关闭」实例：注册关/setup已初始化/SMTP缺省下的负面情报与速判清单 | info-leak | poc | Y | new-api,one-api,信息泄露,负面情报,未授权面,注册关闭,setup假象,限速预算 |
| 427 | go-captcha 拼图行为验证码破解：Sobel 梯度模板匹配定位 + /api/captcha/verify 直接换 token（4/4 实测通过） | auth-bypass | exp | Y | go-captcha,验证码绕过,行为验证码,模板匹配,sobel,new-api,modelsell |
| 428 | New API Token Quota 面板 /token-quota.config.json 配置泄露 + 令牌状态判别法（unavailable=存在但禁用） | info-leak | poc | Y | new-api,token-quota,config-leak,sk-token,验活判据,unavailable,modelsell |
| 429 | <target> 注册墙全链路程序化突破：mail.tm 临时邮箱 + GET /api/verification 真实收码 + go-captcha(register) 建号成功（新用户初始 quota=0） | auth-bypass | exp | Y | new-api,modelsell,register,mail.tm,临时邮箱,verification,go-captcha,quota, |
| 430 | Discuz <target> 后台 runquery/SQL 导入对非创始人 adminid=1 全线 403（isfounder 闸先于 checkpermission 生效，实测证否） | unauthorized | poc | Y | discuz,<target>,admincp,runquery,isfounder,后台,SQL执行,负面情报,into-outfile |
| 431 | PHPWind 后台操作日志 data/bbscache/admin_record.php 零鉴权直读（含真实 adminjob 模块与 bakup 备份前缀） | info-leak | poc | Y | phpwind,postlog,adminlog,info-leak,backup-prefix,unauthorized |
| 432 | New API 邮箱验证注册绕过（一次性邮箱 + 6 位十六进制验证码 + 无图形验证码） | auth-bypass | poc | Y | new-api,注册绕过,邮箱验证,一次性邮箱,ondemand |
| 433 | PHPWind <target> 全量源码获取：<target> 官方下载站快照直取法（免 GitHub） | info-leak | template | Y | phpwind,<target>,源码获取,wayback,archive,镜像,下载,升级包,offline-audit |
| 434 | melody-admin pageconfig tableSuffix 未授权布尔盲注（MySQL root，跨租户分表切换） | sqli | poc | Y | sqli,blind,mysql,unauth,pageconfig,multi-tenant |
| 435 | 七牛云 Kodo AK/SK 接管验证与桶枚举（QBox 签名保留 base64 填充） | unauthorized | script | Y | qiniu,kodo,cloud,aksk,object-storage,info-leak |
| 436 | renren-fast 派生后台：SQL 复制角色绕过角色创建者校验 + 自建管理员规避 token IP 互踢 | auth-bypass | poc | Y |  |
| 437 | melody 会员面 pageconfig/sysparam 签名-only 未授权读不可提权（AG49 生产实例） | info-leak | poc | Y | melody,pageconfig,sysparam,unauthorized,negative,ag49 |
| 438 | 七牛 Kodo 泄露 AK/SK：全量枚举、精确计量与公开读桶匿名外带 | info-leak | script | Y | qiniu,kodo,object-storage,public-bucket,ak-sk,info-leak,cdn |
| 439 | AG49 会员端 Spring Boot 2.x 框架层 N-day 与反序列化全否（Shiro/Spring4Shell/CVE-2018-1273/Fastjson/actuator/env） | deserialization | markdown | Y |  |
| 440 | melody-admin(renren-fast 派生) 调度器 mode 语义与 MySQL UDF 面判否 + UDF 正控方法 | rce | script | Y |  |
| 441 | openresty 1.29 + Spring Boot：请求走私(CL.TE/TE.CL/TE.TE)与 h2 降级全族阴性（勿再重打） | auth-bypass | payload | Y |  |
| 442 | melody-admin 彩种模板同步 SSRF：回显只认 data 数组信封；60s 时延不能当内网 DROP 判据 | ssrf | poc | Y | melody-admin,ssrf,time-based,oracle,readback,hutool,内网探测 |
| 443 | AG49 IM 节点 <target>:2585 配置未授权读不可提权（管理面/写面/注入全阴性 + 登录滑块阻断 staff 口令攻击） | info-leak | poc | Y | melody,lark,im,ag49,unauthorized,info-leak,negative,slider-captcha,258 |
| 444 | X-Forwarded-For 等代理头伪造客户端 IP | auth-bypass | poc | Y | ip-spoof,x-forwarded-for,proxy-header,auth-bypass,rate-limit,audit-log |
| 445 | MySQL root 空口令远程直连检测（只读取证脚本） | weak-password | poc | Y | mysql,空口令,未授权,数据库,caching_sha2_password |
| 446 | Laravel 图片上传（客户端后缀白名单+雪花ID重命名+Flysystem path 参数）无法转 RCE，残留 SVG 同源 XSS 与 4MiB FPM DoS | file-upload | poc | Y | file-upload,laravel,flysystem,svg-xss,php-fpm,dos,negative |
| 447 | sendSmsCode 响应回显短信验证码 + modifyPasswordBySms 免短信接管任意账号（epay/Upay 会员 API） | auth-bypass | exp | Y |  |
| 448 | 短信找回密码链：sendSmsCode 响应明文回显 OTP + modifyPasswordBySms 未认证改密 → 任意会员账号接管 | auth-bypass | exp | Y | ATO,sms,sendSmsCode,modifyPasswordBySms,user-enumeration,account-takeo |
| 449 | PHP session.upload_progress 会话文件注入 + LFI → 预认证 RCE（Chain A） | rce | exp |  | php,session,rfc1867,upload_progress,lfi,rce,预认证 |
| 450 | PHP var_destroy 二次析构 UAF → `/proc/self/mem` 覆写 internal function handler → RCE（Chain B） | rce | exp |  | php,uaf,var_destroy,unserialize,proc-self-mem,handler-overwrite,rce,提权 |
| 451 | PHP session serialize_handler 不一致 → 解码失败仍触发 __destruct → POP gadget RCE（Chain C） | rce | exp |  | php,session,serialize_handler,pop-gadget,destruct,php-fpm,rce |
| 452 | CVE-2026-109618 Apache mod_php：multipart 解析期 OOM → longjmp 跳过 init_executor → shutdown 读脏 EG() → RCE（Chain D，零依赖） | rce | exp |  | php,oom,longjmp,mod_php,apache,eg-state,零依赖,rce |
| 453 | CVE-2026-109619 nginx + php-fpm：脏 BG(unserialize) 状态被下一请求消费 → __destruct → RCE（Chain E） | rce | exp |  | php,fpm,bg-state,unserialize,session,use-after-free,rce |
| 454 | PHP 环境存在性判别器（跑 PHP2Shell 五链前的必备预检） | app-audit | script | Y | php,php-fpm,mod_php,sapi,fingerprint,triage,php2shell,preflight |
| 455 | 授权总控类面板默认口令 admin:admin（HTTP Basic）→ 全量授权码+机器绑定接管 | weak-password | exp | Y | basic-auth,默认口令,授权码,后台接管,python |
| 457 | PHP2Shell 通用利用框架（CVE-2026-109618/619）：session.upload_progress UAF → Heap OOB 泄露绕 ASLR → pDestructor=zif_system 预认证 RCE | rce | script | Y | php, php2shell, cve-2026-109618, cve-2026-109619, session.upload_progr |
| 458 | 3x-ui-x-ui-assets-unauth-directory-listing-and-dual-404-base | absorbed | script |  | 3x,ui,ui,assets,unauth,directory,listing,dual |
| 459 | cloudfront-behavior-prefix-vs-origin-normalization-bypass-co | absorbed | markdown |  | cloudfront,behavior,prefix,vs,origin,normalization,bypass,co |
| 460 | edge-host-127.0.0.1-baota-default-site-nginx-stub-status-una | absorbed | script |  | edge,host,127.0.0.1,baota,default,site,nginx,stub |
| 461 | edge-host-whitelist-third-party-vhost-isolation-config-and-p | absorbed | script |  | edge,host,whitelist,third,party,vhost,isolation,config |
| 462 | fofa-官方接口可用性切换-target417.example.com-失效-fofa.info | absorbed | script |  | fofa,官方接口可用性切换,target417.example.com,失效,fofa.info |
| 463 | h2c-grpc-日志摄取接口-pb.logingestservice-ingestlogevents-未授权调用-无 | absorbed | markdown |  | h2c,grpc,日志摄取接口,pb.logingestservice,ingestlogevents,未授权调用 |
| 464 | landing-page-localstorage-query-channelcode-iframe-src-path | absorbed | markdown |  | landing,page,localstorage,query,channelcode,iframe,src,path |
| 465 | mail-auth-missing-spf-dkim-dmarc-email-spoofing | absorbed | script |  | mail,auth,missing,spf,dkim,dmarc,email,spoofing |
| 466 | minio-admin-v3-streaming-unsigned-payload-trailer-bypass-cve | absorbed | markdown |  | minio,admin,v3,streaming,unsigned,payload,trailer,bypass |
| 467 | minio-streaming-unsigned-payload-trailer-sigv4-bypass-cve-20 | absorbed | markdown |  | minio,streaming,unsigned,payload,trailer,sigv4,bypass,cve |
| 468 | mysql-3306-exposed-host-not-allowed-ip-whitelist-three-egres | absorbed | script |  | mysql,3306,exposed,host,not,allowed,ip,whitelist |
| 469 | omnihelm-console-grpc-bff-host-oracle | absorbed | script |  | omnihelm,console,grpc,bff,host,oracle |
| 470 | oss-listbucket-denied-getobject-allowed-anon-read-matrix | absorbed | script |  | oss,listbucket,denied,getobject,allowed,anon,read,matrix |
| 471 | photon-edge-clte-desync-smuggle-0chunk-tail-forwarded | absorbed | script |  | photon,edge,clte,desync,smuggle,0chunk,tail,forwarded |
| 472 | postmessage-missing-origin-config-injection-js-sink | absorbed | markdown |  | postmessage,missing,origin,config,injection,js,sink |
| 473 | postmessage-source-identity-bypass-iframe-navigation-new-fun | absorbed | markdown |  | postmessage,source,identity,bypass,iframe,navigation,new,fun |
| 474 | postmessage-仅校验-source-的-executefunction-绕过-跨源导航被校验子帧-new | absorbed | markdown |  | postmessage,仅校验,source,executefunction,绕过,跨源导航被校验子帧,new |
| 475 | sclub-168six-up_photo-a_photo-头像与直播台图标上传-.php-文件名被丢弃-gd-重编码裁 | absorbed | markdown |  | sclub,168six,up,photo,photo,头像与直播台图标上传,.php,文件名被丢弃 |
| 476 | somcns-thinkphp3.1.3-站群-166-站-admin-autoart-控制器脱队-extends-ac | absorbed | markdown |  | somcns,thinkphp3.1.3,站群,166,admin,autoart,控制器脱队,extends |
| 477 | somcns-thinkphp3.1.3-站群-精简版-实测负面清单-控制器文件存在性预言机-wechat-auto | absorbed | markdown |  | somcns,thinkphp3.1.3,站群,精简版,实测负面清单,控制器文件存在性预言机,wechat,auto |
| 478 | somcns-站群前台会员登录-home-index-login-无验证码无频控-test-test-直登-后台登 | absorbed | markdown |  | somcns,站群前台会员登录,home,index,login,无验证码无频控,test,test |
| 479 | sub2api-0.2.19-activities-fork-reset-feature-gate-dead | absorbed | markdown |  | sub2api,0.2.19,activities,fork,reset,feature,gate,dead |
| 480 | sub2api-0.2.19-密码重置链审查-feature-gate-双开关-token-设计-reset-不清-to-cve-2026-27812 | absorbed | markdown |  | sub2api,0.2.19,密码重置链审查,feature,gate,双开关,token,设计 |
| 481 | sub2api-admin-login-weak-password-domain-derived-jwt-takeove | absorbed | markdown |  | sub2api,admin,login,weak,password,domain,derived,jwt |
| 482 | sub2api-invalid-auth-rate-limit-429-shared-bucket-dos | absorbed | script |  | sub2api,invalid,auth,rate,limit,429,shared,bucket |
| 483 | sub2api-photonedge-nginx-path-normalization-no-auth-bypass | absorbed | markdown |  | sub2api,photonedge,nginx,path,normalization,no,auth,bypass |
| 484 | sub2api-plain-http-no-redirect-missing-hsts-sslstrip | absorbed | markdown |  | sub2api,plain,http,no,redirect,missing,hsts,sslstrip |
| 485 | sub2api-validate-promo-code-unauth-enum-oracle | absorbed | script |  | sub2api,validate,promo,code,unauth,enum,oracle |
| 486 | sub2api-网关-api-key-未随凭据轮换失效-dump-泄露的-sk-key-仍可调用-v1-付费接口 | absorbed | markdown |  | sub2api,网关,api,key,未随凭据轮换失效,dump,泄露的,sk |
| 487 | sub2api-邮箱验证码接口账号存在性枚举-oracle-公开-api-v1-auth-send-verify-co | absorbed | script |  | sub2api,邮箱验证码接口账号存在性枚举,oracle,公开,api,v1,auth,send |
| 488 | target160.example.com-1852-补丁面黑盒实测矩阵-13-类已封堵-5-个门槛不可越-write_config-落点存在但入口不可达 | absorbed | markdown |  | target160.example.com,1852,补丁面黑盒实测矩阵,13,类已封堵,个门槛不可越,write,config |
| 489 | thinkphp-3.1.3-站群-cms-admin-控制器-403-假分叉-判别-实为-nginx-层危险函数 | absorbed | markdown |  | thinkphp,3.1.3,站群,cms,admin,控制器,403,假分叉 |
| 490 | thinkphp-3.1.3-站群-home-index-bbs-的-post-参数零鉴权渲染文章正文-游客态可批量拖 | absorbed | markdown |  | thinkphp,3.1.3,站群,home,index,bbs,post,参数零鉴权渲染文章正文 |
| 491 | thinkphp-3.1.3-站群两个零凭证原生语-think_page_trace-泄全部内部路径-app | absorbed | markdown |  | thinkphp,3.1.3,站群两个零凭证原生语,think,page,trace,泄全部内部路径,app |
| 492 | thinkphp-站群-kindeditor-php-版本-upload_json.php-file_manager | absorbed | markdown |  | thinkphp,站群,kindeditor,php,版本,upload,json.php,file |
| 493 | 博彩站群-203.0.113.1-10666-java-反序列化-jndi-log4j-rce-全否-严格 | absorbed | markdown |  | 博彩站群,203.0.113.1,10666,java,反序列化,jndi,log4j,rce |
| 494 | 博彩站群-203.0.113.28-10666-java-反序列化-jndi-log4j-rce-全否-严格 | absorbed | markdown |  | 博彩站群,203.0.113.28,10666,java,反序列化,jndi,log4j,rce |
| 495 | 邮件伪造类-finding-缺失-spf-dkim-dmarc-的核实方法与-不可利用-边界 | absorbed | script |  | 邮件伪造类,finding,缺失,spf,dkim,dmarc,的核实方法与,不可利用 |
| 496 | 203.0.113.1-.23-先验排除表-两条恒定假阳性-code-验证码端点-爆破前置-五处已穷尽面 | absorbed | script |  | 203.0.113.1,.23,先验排除表,两条恒定假阳性,code,验证码端点,爆破前置,五处已穷尽面 |
| 497 | 203.0.113.1-站群跳转机全机地图与负面情报-唯一动态程序-phpwind-7.3.2-全静态页-三套异机真入口 | absorbed | http |  | 203.0.113.1,站群跳转机全机地图与负面情报,唯一动态程序,phpwind,7.3.2,全静态页,三套异机真入口 |
| 498 | 203.0.113.10-.23-先验排除表-两条恒定假阳性-code-验证码端点-爆破前置-五处已穷尽面 | absorbed | script |  | 203.0.113.10,.23,先验排除表,两条恒定假阳性,code,验证码端点,爆破前置,五处已穷尽面 |
| 499 | 203.0.113.13-站群跳转机全机地图与负面情报-唯一动态程序-phpwind-7.3.2-全静态页-三套异机真入口 | absorbed | http |  | 203.0.113.13,站群跳转机全机地图与负面情报,唯一动态程序,phpwind,7.3.2,全静态页,三套异机真入口 |
| 500 | discuz-7.2-后台-runquery-sql-导入对非创始人-adminid-1-全线-403-isfounde | absorbed | markdown |  | discuz,7.2,后台,runquery,sql,导入对非创始人,adminid,全线 |
| 501 | phpwind-7.3.2-7.5-后台写文件-getshell-三处写入端判死-pcv-封穿越-pw_var_export-转义数组键-setcss-非写入端 | absorbed | poc |  | phpwind,7.3.2,7.5,后台写文件,getshell,三处写入端判死,pcv,封穿越 |
| 502 | phpwind-7.3.2-负面情报清单-register_globals-已移除-mawhole-eval-证否-setstyles-转义-hack.php | absorbed | http |  | phpwind,7.3.2,负面情报清单,register,globals,已移除,mawhole,eval |
| 503 | phpwind-8.7-decode-错位链-白盒审计结论-无需账号的裸引号-convert-e-路径不存在-cms-投稿链裸引号入库-展示-rce-代码级成立 | absorbed | poc |  | phpwind,8.7,decode,错位链,白盒审计结论,无需账号的裸引号,convert,路径不存在 |
| 504 | phpwind-8.7-补丁面速查表-target124.example.com-白盒产出-已修复-9-类-未验证缺陷-6-条-未穷尽面-4-个-全附源码行号 | absorbed | poc |  | phpwind,8.7,补丁面速查表,target124.example.com,白盒产出,已修复,未验证缺陷,未穷尽面 |
| 505 | sclub-平台-index.phtml-注入回显客户端-含三个易踩根因与语句约束表 | absorbed | script |  | sclub,平台,index.phtml,注入回显客户端,含三个易踩根因与语句约束表 |
| 506 | sclub-平台-index.phtml-注入客户端修正-繁体歡迎-hex左截断-from包层-限流特征 | absorbed | markdown |  | sclub,平台,index.phtml,注入客户端修正,繁体歡迎,hex左截断,from包层,限流特征 |
| 507 | tgt-1852-补丁面黑盒实测矩阵-13-类已封堵-5-个门槛不可越-write_config-落点存在但入口不可达 | absorbed | markdown |  | tgt,1852,补丁面黑盒实测矩阵,13,类已封堵,个门槛不可越,write,config |
| 508 | tgt2403-admin-后台指纹与-api-v1-admin-管理-api-全清单提取-含未授权定性探测 | absorbed | markdown |  | tgt2403,admin,后台指纹与,api,v1,admin,管理,api |
| 509 | tgt2403-发卡系统-订单-游客订单鉴权模型-api-契约-订单号格式与-idor-判定速查 | absorbed | markdown |  | tgt2403,发卡系统,订单,游客订单鉴权模型,api,契约,订单号格式与,idor |
| 510 | tgt2403-版本指纹-上游补丁差分法-app_version-免认证泄露-缺失安全修复清单-逐条可达性实测 | absorbed | markdown |  | tgt2403,版本指纹,上游补丁差分法,app,version,免认证泄露,缺失安全修复清单,逐条可达性实测 |
| 511 | tgt2403-用户中心与认证链路速查-me-wallet-affiliate-2fa-鉴权事实-验证码节流与限流缺口 | absorbed | markdown |  | tgt2403,用户中心与认证链路速查,me,wallet,affiliate,2fa,鉴权事实,验证码节流与限流缺口 |
| 512 | thinkphp-3.1.3-站群-cms-模板烘焙-rce-链完整打法-两个否证-oracle-get_html1-零凭证覆盖模板-abrdns-抓取源陷阱 | absorbed | markdown |  | thinkphp,3.1.3,站群,cms,模板烘焙,rce,链完整打法,两个否证 |
| 513 | thinkphp3.1.3-站群-cms-通杀-rce-admin-表哈希即凭据-home-article-edit_cont-写-article.title | absorbed | markdown |  | thinkphp3.1.3,站群,cms,通杀,rce,admin,表哈希即凭据,home |
| 514 | waf-绕过-能力修正-直连-https-清空-sni-任意-host-200-带真实-target320.example.com-rst-按-sni-丢包 | absorbed | script |  | waf,绕过,能力修正,直连,https,清空,sni,任意 |
| 515 | 代理分流造假样本-target30.example.com-经-7892-返回伪造-502-404-hello | absorbed | script |  | 代理分流造假样本,target30.example.com,7892,返回伪造,502,404,hello |
| 516 | 博彩-phpwind-站群实战结论集-203.0.113.1-24-判死路线-8-条可复用判据 | absorbed | http |  | 博彩,phpwind,站群实战结论集,203.0.113.1,24,判死路线,条可复用判据 |
| 517 | 博彩-phpwind-站群实战结论集-203.0.113.27-24-判死路线-8-条可复用判据 | absorbed | http |  | 博彩,phpwind,站群实战结论集,203.0.113.27,24,判死路线,条可复用判据 |
| 518 | Nexus Repository 未授权 RCE（CVE-2019-7238） | cloud-container | markdown |  | nexus,sonatype,2019-7238,未授权,rce,供应链,cloud-container,fujiang |
| 519 | PostgreSQL COPY PROGRAM 权限边界（CVE-2019-9193，厂商认定非漏洞） | cloud-container | markdown |  | postgresql,postgres,copy,program,命令执行,数据库,cloud-container,fujiang |
| 520 | vCenter Server 未授权上传 → RCE（vROPS 插件） | cloud-container | markdown |  | vcenter,vrops,uploadova,21972,未授权,cloud-container,fujiang |
| 521 | Nacos 认证绕过（User-Agent 白名单） | cloud-container | markdown |  | nacos,user-agent,auth-bypass,cve-2021-29441,cloud-container,fujiang |
| 522 | Grafana 插件路径穿越任意文件读（CVE-2021-43798） | cloud-container | markdown |  | grafana,2021-43798,路径穿越,文件读取,监控,cloud-container,fujiang |
| 523 | cgroups release_agent 容器逃逸 | cloud-container | markdown |  | cgroup,release_agent,escape,cap_sys_admin,container,cloud-container,fu |
| 524 | VMware Workspace ONE Access 模板注入 RCE | cloud-container | markdown |  | workspace one,freemarker,ssti,22954,未授权,cloud-container,fujiang |
| 525 | Veeam Backup 凭据泄露（CVE-2023-27532） | cloud-container | markdown |  | veeam,2023-27532,凭据,备份,cloud-container,fujiang |
| 526 | runc 文件描述符泄露逃逸（Leaky Vessels） | cloud-container | markdown |  | runc,leaky vessels,fd leak,workdir,escape,docker build,Docker,containe |
| 527 | VMware ESXi 认证绕过 → hypervisor 失守（CVE-2024-37085） | cloud-container | markdown |  | esxi,vmware,2024-37085,ad,认证绕过,勒索,cloud-container,fujiang |
| 528 | Veeam Backup 未认证 RCE（CVE-2024-40711） | cloud-container | markdown |  | veeam,2024-40711,反序列化,未认证,rce,cloud-container,fujiang |
| 529 | Docker Remote API 未授权（直接拿宿主） | cloud-container | markdown |  | docker,api,2375,privileged,escape,cloud-container,fujiang |
| 530 | MongoDB 未授权访问（直接拖库） | cloud-container | markdown |  | mongodb,未授权,拖库,数据库,nosql,cloud-container,fujiang |
| 531 | Redis 未授权（写文件 / 模块加载 / 主从复制） | cloud-container | markdown |  | redis,unauth,config-set,rdb,cron,module-load,replica,cloud-container,f |
| 532 | FortiOS SSL VPN 路径穿越（凭据读取） | edge-device | markdown |  | fortinet,fortigate,fortios,sslvpn,path-traversal,fgt_lang,edge-device, |
| 533 | Pulse Secure 任意文件读取（凭据与会话） | edge-device | markdown |  | pulse,ivanti,file-read,session,凭据,edge-device,fujiang |
| 534 | Citrix ADC 路径穿越 → 模板注入 RCE | edge-device | markdown |  | citrix,netscaler,path-traversal,newbm,cve-2019-19781,edge-device,fujia |
| 535 | FortiOS SSL VPN 双因素绕过（用户名大小写） | edge-device | markdown |  | fortigate,fortios,sslvpn,2fa bypass,ldap,大小写,edge-device,fujiang |
| 536 | F5 BIG-IP TMUI 未授权 RCE（..;/ 路径穿越） | edge-device | markdown |  | f5,bigip,tmui,5902,tmsh,路径穿越,edge-device,fujiang |
| 537 | F5 iControl REST 未授权 RCE（2021 变体） | edge-device | markdown |  | f5,icontrol,22986,bash,未授权,edge-device,fujiang |
| 538 | F5 BIG-IP iControl REST 认证绕过 → 命令执行 | edge-device | markdown |  | f5,bigip,icontrol,auth-bypass,tmsh,cve-2022-1388,edge-device,fujiang |
| 539 | FortiOS / FortiProxy 管理面认证绕过（假 Forwarded 头） | edge-device | markdown |  | fortinet,fortios,fortiproxy,auth bypass,forwarded header,edge-device,f |
| 540 | Ivanti EPMM 认证绕过（CVE-2023-35078） | edge-device | markdown |  | ivanti,epmm,mobileiron,2023-35078,auth bypass,edge-device,fujiang |
| 541 | Ivanti Connect Secure 认证绕过 + 命令注入（串联） | edge-device | markdown |  | ivanti,pulse,authentication-bypass,command-injection,cve-2024-21887,ed |
| 542 | FortiClient EMS SQL 注入 → 未认证 RCE（CVE-2023-48788） | edge-device | markdown |  | forticlient,ems,sqli,2023-48788,未授权,edge-device,fujiang |
| 543 | PAN-OS 管理面认证绕过 → root RCE（CVE-2024-0012 + CVE-2024-9474） | edge-device | markdown |  | panos,paloalto,2024-0012,2024-9474,auth bypass,管理面,edge-device,fujiang |
| 544 | GoAnywhere MFT 认证绕过 → 新建管理员（CVE-2024-0204） | edge-device | markdown |  | goanywhere,mft,2024-0204,认证绕过,管理员,edge-device,fujiang |
| 545 | Check Point Quantum 安全网关任意文件读（CVE-2024-24919） | edge-device | markdown |  | checkpoint,quantum,2024-24919,文件读取,vpn,网关,edge-device,fujiang |
| 546 | PAN-OS GlobalProtect 命令注入 | edge-device | markdown |  | panos,globalprotect,sessid,command-injection,cve-2024-3400,edge-device |
| 547 | FortiManager 未认证 RCE（FortiJump，CVE-2024-47575） | edge-device | markdown |  | fortimanager,fortijump,2024-47575,fgfm,未认证,edge-device,fujiang |
| 548 | AD CS 证书滥用（ESC1 / ESC8 等） | exchange-ad | markdown |  | adcs,esc1,esc8,certipy,pki,certificate-template,exchange-ad,fujiang |
| 549 | EternalBlue（MS17-010 SMBv1 远程溢出） | exchange-ad | markdown |  | eternalblue,ms17-010,smbv1,doublepulsar,cve-2017-0144,exchange-ad,fuji |
| 550 | Exchange ViewState 反序列化（静态 machineKey） | exchange-ad | markdown |  | exchange,viewstate,machinekey,ysoserial.net,ecp,exchange-ad,fujiang |
| 551 | Zerologon（Netlogon 密钥置零 → 域控接管） | exchange-ad | markdown |  | zerologon,netlogon,dcsync,cve-2020-1472,domain-controller,exchange-ad, |
| 552 | Exchange ProxyLogon（SSRF + 任意文件写入 → RCE） | exchange-ad | markdown |  | exchange,proxylogon,ssrf,ecp,autodiscover,cve-2021-26855,exchange-ad,f |
| 553 | Exchange ProxyShell（三个洞串联 → RCE） | exchange-ad | markdown |  | exchange,proxyshell,autodiscover,ews,mailboxexport,cve-2021-34473,exch |
| 554 | PrintNightmare（打印后台服务任意 DLL 加载） | exchange-ad | markdown |  | printnightmare,spooler,rpcaddprinterdriverex,dll,cve-2021-34527,exchan |
| 555 | noPac（sAMAccountName 伪装 + 无 PAC 票据 → 域管） | exchange-ad | markdown |  | nopac,samaccountname,kerberos,golden-ticket,cve-2021-42278,exchange-ad |
| 556 | Exchange 认证后反序列化 RCE（CVE-2021-42321） | exchange-ad | markdown |  | exchange,2021-42321,反序列化,ews,rce,exchange-ad,fujiang |
| 557 | MSDT 协议远程代码执行（Follina） | exchange-ad | markdown |  | follina,msdt,30190,word,rtf,钓鱼,exchange-ad,fujiang |
| 558 | Exchange ProxyNotShell（SSRF + PowerShell 远程） | exchange-ad | markdown |  | exchange,proxynotshell,owassrf,ssrf,powershell,exchange-ad,fujiang |
| 559 | Exchange EWS NTLM 中继提权（CVE-2024-21410） | exchange-ad | markdown |  | exchange,ews,2024-21410,ntlm,relay,提权,exchange-ad,fujiang |
| 560 | Outlook MonikerLink 绕过与 NTLM 泄露（CVE-2024-21413） | exchange-ad | markdown |  | outlook,monikerlink,2024-21413,ntlm,钓鱼,exchange-ad,fujiang |
| 561 | Zimbra postjournal 未授权命令执行（CVE-2024-45519） | exchange-ad | markdown |  | zimbra,2024-45519,postjournal,未授权,rce,邮件,exchange-ad,fujiang |
| 562 | WebLogic T3 反序列化（CommonsCollections） | java-deser | markdown |  | weblogic,t3,deserialization,commonscollections,ysoserial,java-deser,fu |
| 563 | Shiro rememberMe 默认密钥反序列化 RCE（CVE-2016-4437） | java-deser | markdown |  | shiro,rememberme,反序列化,默认密钥,java,java-deser,fujiang |
| 564 | WebLogic wls-wsat XMLDecoder 反序列化 RCE | java-deser | markdown |  | weblogic,wls-wsat,xmldecoder,soap,coordinatorporttype,java-deser,fujia |
| 565 | Struts2 S2-045 Content-Type OGNL 注入 | java-deser | markdown |  | struts2,s2-045,ognl,multipart,content-type,java-deser,fujiang |
| 566 | Apache Log4j2 JNDI 注入（Log4Shell） | java-deser | markdown |  | log4j,log4shell,jndi,ldap,rmi,dnslog,log4j-core,java-deser,fujiang |
| 567 | SnakeYAML 反序列化 RCE（CVE-2022-1471） | java-deser | markdown |  | snakeyaml,yaml,2022-1471,反序列化,rce,java,java-deser,fujiang |
| 568 | Java ECDSA 空签名绕过（Psychic Signatures） | java-deser | markdown |  | ecdsa,psychic-signatures,jwt,es256,saml,签名绕过,java-deser,fujiang |
| 569 | Spring Framework 数据绑定 RCE（Spring4Shell） | java-deser | markdown |  | spring,spring4shell,classloader,tomcat,accessLogValve,java-deser,fujia |
| 570 | Fastjson autoType 反序列化 RCE | java-deser | markdown |  | fastjson,autotype,jdbcrowsetimpl,jndi,json,java-deser,fujiang |
| 571 | Tomcat PUT 任意文件写入（JSP 落地） | java-web | markdown |  | tomcat,put,jsp,readonly,windows,12615,java-web,fujiang |
| 572 | Apache Solr Velocity 模板注入 RCE | java-web | markdown |  | solr,velocity,params.resource.loader,template,cve-2019-17558,java-web, |
| 573 | WebLogic wls9_async_response 未授权 RCE（2019 在野） | java-web | markdown |  | weblogic,2725,wls9_async,xmldecoder,未授权,java-web,fujiang |
| 574 | WebLogic 控制台未授权 RCE（console.portal 路径穿越） | java-web | markdown |  | weblogic,console,path traversal,shellsession,mvel,java-web,fujiang |
| 575 | Tomcat AJP 文件读取 / 包含（Ghostcat） | java-web | markdown |  | tomcat,ajp,ghostcat,web.xml,文件包含,java-web,fujiang |
| 576 | Apache Druid 未授权 RCE（CVE-2021-25646） | java-web | markdown |  | druid,apache,2021-25646,未授权,rce,时序库,java-web,fujiang |
| 577 | Confluence 未授权 OGNL 注入（pre-auth，ENCODED_TOKEN） | java-web | markdown |  | confluence,ognl,pre-auth,queryString,doenterpagevariables,java-web,fuj |
| 578 | Spring Cloud Gateway Actuator SpEL 注入 RCE（CVE-2022-22947） | java-web | markdown |  | spring,gateway,actuator,spel,2022-22947,rce,java-web,fujiang |
| 579 | Spring Cloud Function SpEL 注入 RCE（CVE-2022-22963） | java-web | markdown |  | spring,cloud,function,2022-22963,spel,rce,java-web,fujiang |
| 580 | Confluence 未授权 OGNL 注入（Confluence OGNL RCE） | java-web | markdown |  | confluence,ognl,cve-2022-26134,atlassian,rce,java-web,fujiang |
| 581 | Confluence 未授权管理员接管（setup 流程未关） | java-web | markdown |  | confluence,22515,setupadministrator,未授权,管理员,java-web,fujiang |
| 582 | Confluence 未授权恢复导致管理员接管 | java-web | markdown |  | confluence,setup-restore,未授权,管理员接管,java-web,fujiang |
| 583 | Confluence 模板注入 RCE（text-inline.vm） | java-web | markdown |  | confluence,ssti,velocity,text-inline,cve-2023-22527,java-web,fujiang |
| 584 | PaperCut NG/MF 认证绕过 → RCE（CVE-2023-27350） | java-web | markdown |  | papercut,2023-27350,auth bypass,setup,打印,java-web,fujiang |
| 585 | TeamCity 认证绕过 → 管理员令牌 → RCE | java-web | markdown |  | teamcity,auth-bypass,token,cve-2023-42793,java-web,fujiang |
| 586 | Apache ActiveMQ OpenWire 反序列化 RCE | java-web | markdown |  | activemq,openwire,deserialization,classinfo,61616,java-web,fujiang |
| 587 | Struts2 文件上传路径穿越（S2-066） | java-web | markdown |  | struts2,s2-066,upload,path-traversal,filename,java-web,fujiang |
| 588 | GitLab 任意账号接管（CVE-2023-7028） | java-web | markdown |  | gitlab,2023-7028,账号接管,密码重置,java-web,fujiang |
| 589 | Jenkins CLI 任意文件读取 | java-web | markdown |  | jenkins,cli,args4j,file-read,cve-2024-23897,java-web,fujiang |
| 590 | TeamCity 鉴权绕过（;.jsp 路径参数） | java-web | markdown |  | teamcity,auth bypass,27198,rest,token,java-web,fujiang |
| 591 | Apache OFBiz 未授权 RCE（CVE-2024-38856） | java-web | markdown |  | ofbiz,2024-38856,未授权,rce,电商,java-web,fujiang |
| 592 | Struts2 文件上传路径穿越 → RCE（S2-067，CVE-2024-53677） | java-web | markdown |  | struts2,s2-067,2024-53677,文件上传,路径穿越,rce,java-web,fujiang |
| 593 | Baron Samedit（sudo 堆溢出提权） | linux-priv | markdown |  | sudo,baron-samedit,sudoedit,heap-overflow,cve-2021-3156,linux-priv,fuj |
| 594 | Ubuntu overlayfs 本地提权（2021） | linux-priv | markdown |  | ubuntu,overlayfs,3493,userns,setuid,linux-priv,fujiang |
| 595 | PwnKit（polkit pkexec 本地提权） | linux-priv | markdown |  | pwnkit,pkexec,polkit,suid,cve-2021-4034,linux-priv,fujiang |
| 596 | fs_context 堆溢出提权（legacy_parse_param） | linux-priv | markdown |  | fs_context,heap overflow,unshare,userns,cap_sys_admin,linux-priv,fujia |
| 597 | Dirty Pipe（页缓存写覆盖提权） | linux-priv | markdown |  | dirtypipe,page-cache,splice,cve-2022-0847,linux-priv,fujiang |
| 598 | overlayfs copy-up 权限绕过提权 | linux-priv | markdown |  | overlayfs,copy-up,setuid,unshare,cve-2023-0386,linux-priv,fujiang |
| 599 | sudoedit 环境变量注入提权（EDITOR 带 --） | linux-priv | markdown |  | sudo,sudoedit,editor,22809,提权,linux-priv,fujiang |
| 600 | Ubuntu overlayfs 权限提升（GameOver(lay)） | linux-priv | markdown |  | ubuntu,overlayfs,gameoverlay,xattr,setuid,提权,linux-priv,fujiang |
| 601 | nf_tables 释放后重用提权（2023） | linux-priv | markdown |  | nftables,netfilter,uaf,32233,unshare,linux-priv,fujiang |
| 602 | Linux io_uring 本地提权（CVE-2024-0582） | linux-priv | markdown |  | linux,io_uring,2024-0582,uaf,本地提权,linux-priv,fujiang |
| 603 | nf_tables 释放后重用提权 | linux-priv | markdown |  | nftables,netfilter,uaf,unshare,cve-2024-1086,linux-priv,fujiang |
| 604 | 命令注入（分号 / 反引号 / $() 与盲注外带） | web-cms | markdown |  | 命令注入,rce,ping,反引号,盲注,shell,web-cms,fujiang |
| 605 | PHPUnit eval-stdin 远程代码执行 | web-cms | markdown |  | phpunit,eval-stdin,vendor,9841,rce,web-cms,fujiang |
| 606 | Drupal 表单 API 远程代码执行（Drupalgeddon2） | web-cms | markdown |  | drupal,drupalgeddon2,form-api,post_render,cve-2018-7600,web-cms,fujian |
| 607 | PHP-FPM 路径信息下溢 RCE | web-cms | markdown |  | phpfpm,nginx,fastcgi,cve-2019-11043,phuip-fpizdam,web-cms,fujiang |
| 608 | vBulletin widget 配置代码执行 | web-cms | markdown |  | vbulletin,widget_php,preg_replace,cve-2019-16759,web-cms,fujiang |
| 609 | Laravel Ignition 调试模式任意文件读 → RCE（CVE-2021-3129） | web-cms | markdown |  | laravel,ignition,2021-3129,rce,debug,php,web-cms,fujiang |
| 610 | Apache HTTP Server 路径穿越 / CGI RCE | web-cms | markdown |  | httpd,apache,path traversal,cgi,mod_cgi,printf,web-cms,fujiang |
| 611 | Joomla 未授权配置泄露（API） | web-cms | markdown |  | joomla,api,config leak,数据库口令,web-cms,fujiang |
| 612 | WordPress Backup Migration 未授权 RCE（CVE-2023-6553） | web-cms | markdown |  | wordpress,backup-migration,2023-6553,未授权,rce,web-cms,fujiang |
| 613 | WordPress Bricks 未授权 RCE（CVE-2024-25600） | web-cms | markdown |  | wordpress,bricks,2024-25600,未授权,rce,web-cms,fujiang |
| 614 | WordPress WP-Automatic 未授权 SQL 注入（CVE-2024-27956） | web-cms | markdown |  | wordpress,wp-automatic,2024-27956,sqli,未授权,web-cms,fujiang |
| 615 | PHP-CGI Windows 参数注入（best-fit 编码） | web-cms | markdown |  | php,cgi,windows,best-fit,cve-2024-4577,web-cms,fujiang |
| 616 | 本地文件包含 → RCE（日志包含 / pearcmd / PHP filter chain） | web-cms | markdown |  | lfi,文件包含,日志包含,pearcmd,filter chain,rce,web-cms,fujiang |
| 617 | Nginx alias 路径穿越（读任意文件/源码泄露） | web-cms | markdown |  | nginx,alias,穿越,目录,源码,配置,web-cms,fujiang |
| 618 | PHP 反序列化（POP 链 → RCE / phar 触发） | web-cms | markdown |  | php,反序列化,pop链,phar,unserialize,rce,web-cms,fujiang |
| 619 | SQL 注入（报错 / UNION / 堆叠）通用打法 | web-cms | markdown |  | sqli,sql注入,union,报错,堆叠,sqlmap,web-cms,fujiang |
| 620 | SSTI 模板注入 → RCE（Jinja2 / Twig / Freemarker） | web-cms | markdown |  | ssti,模板注入,jinja2,twig,freemarker,rce,web-cms,fujiang |
| 621 | ThinkPHP 5.x 方法注入 / 路由 RCE | web-cms | markdown |  | thinkphp,invokefunction,_method,__construct,filter,rce,web-cms,fujiang |
| 622 | XXE 外部实体注入（读文件 / SSRF / OOB 外带） | web-cms | markdown |  | xxe,xml,外部实体,文件读取,oob,saml,web-cms,fujiang |
| 623 | 后台管理 API 未授权（只在前端隐藏菜单） | gray-market | markdown |  | admin api,未授权,越权,前端隐藏,菜单,gray-market,fujiang |
| 624 | 后台 IP 白名单绕过（X-Forwarded-For 伪造内网来源） | gray-market | markdown |  | 白名单,ip,xff,绕过,后台,内网,gray-market,fujiang |
| 625 | 后台“重置任意用户密码”接口越权（直接接管会员/管理员） | gray-market | markdown |  | 重置密码,任意用户,越权,后台,接管,gray-market,fujiang |
| 626 | 代理层级越权（读全站下级 / 改佣金比例） | gray-market | markdown |  | 代理,agentId,层级越权,佣金,返水,会员列表,gray-market,fujiang |
| 627 | 接口越权 / BOLA（读他人订单、流水、提现） | gray-market | markdown |  | idor,bola,越权,水平越权,orderNo,uid,代理后台,gray-market,fujiang |
| 628 | 接口无 nonce / 时间戳 → 请求重放（充值、领奖、下注） | gray-market | markdown |  | 重放,replay,nonce,时间戳,签名,幂等,gray-market,fujiang |
| 629 | 注单 / 流水接口越权批量拉取（全站会员数据） | gray-market | markdown |  | 注单,流水,idor,越权,批量,报表,gray-market,fujiang |
| 630 | 下注/结算赔率与倍数篡改（资金面） | gray-market | markdown |  | bet,odds,赔率,倍数,结算,篡改,gray-market,fujiang |
| 631 | 支付 / 回调地址可控 → 内网 SSRF（含云元数据） | gray-market | markdown |  | 回调,支付,ssrf,内网,元数据,notify,gray-market,fujiang |
| 632 | 前端签名 / 加密可逆（伪造请求绕过风控） | gray-market | markdown |  | sign,hmac,aes,前端加密,可逆,nonce,timestamp,风控,gray-market,fujiang |
| 633 | CORS 配置错误（带凭据读他人接口数据） | gray-market | markdown |  | cors,跨域,带凭据,origin,泄露,gray-market,fujiang |
| 634 | 券码 / 彩金码可枚举（批量领取与套现） | gray-market | markdown |  | 券码,彩金,兑换,枚举,coupon,领取,gray-market,fujiang |
| 635 | 客服聊天记录导出未鉴权（全站会话与身份材料） | gray-market | markdown |  | 客服,聊天记录,导出,未鉴权,会话,im,gray-market,fujiang |
| 636 | CSRF（无 token 的提现/改密/绑卡） | gray-market | markdown |  | csrf,跨站,提现,改密,绑卡,token,gray-market,fujiang |
| 637 | 签到 / 任务奖励批量刷取（缺幂等与频率控制） | gray-market | markdown |  | 签到,任务,奖励,刷取,幂等,并发,gray-market,fujiang |
| 638 | phpMyAdmin / Adminer 暴露（弱口令直进数据库） | gray-market | markdown |  | phpmyadmin,adminer,pma,数据库,弱口令,脱库,gray-market,fujiang |
| 639 | Druid 监控未授权 → session 泄露/后台接管（若依系必备） | gray-market | markdown |  | druid,监控,session,未授权,若依,jeecg,gray-market,fujiang |
| 640 | 导出 / 下载接口任意文件读取 | gray-market | markdown |  | 任意文件读,download,export,path traversal,对账单,gray-market,fujiang |
| 641 | 财务导出接口越权（银行卡 + 身份证批量导出） | gray-market | markdown |  | 财务,导出,银行卡,身份证,提现,越权,gray-market,fujiang |
| 642 | 三方游戏 / 视讯接口 key 泄露与回调伪造 | gray-market | markdown |  | game vendor,api key,视讯,上分,callback,transfer,gray-market,fujiang |
| 643 | Host 头投毒（密码重置链接被劫持 → 任意账号接管） | gray-market | markdown |  | host,头,投毒,密码重置,链接,劫持,gray-market,fujiang |
| 644 | 客服 IM / WebSocket 未鉴权（订阅他人会话） | gray-market | markdown |  | websocket,socket.io,im,客服,未鉴权,room,订阅,gray-market,fujiang |
| 645 | 后台"切换用户/代登录"接口越权（任意账号接管） | gray-market | markdown |  | impersonate,switch,代登录,客服,越权,接管,gray-market,fujiang |
| 646 | 邀请返水自邀循环（小号矩阵刷返水） | gray-market | markdown |  | 邀请,返水,代理,自邀,循环,佣金,gray-market,fujiang |
| 647 | JWT 实现缺陷（alg:none / 弱密钥 / 可预测 / kid 注入） | gray-market | markdown |  | jwt,alg none,hmac,hashcat 16500,kid,越权,提权,gray-market,fujiang |
| 648 | 开奖结果可预测（弱 PRNG / 客户端算法 / 固定种子） | gray-market | markdown |  | lottery,prng,开奖,可预测,种子,快三,时时彩,gray-market,fujiang |
| 649 | Mass Assignment（注册/资料接口注入 role / vip / balance / agentId） | gray-market | markdown |  | mass assignment,批量赋值,role,vip,balance,agentId,gray-market,fujiang |
| 650 | 会员隐私批量查询 / 账号枚举（手机号、银行卡） | gray-market | markdown |  | member,enumeration,手机号,银行卡,隐私,批量查询,gray-market,fujiang |
| 651 | 包网多租户越权（tenantId/siteId 遍历，横穿几十个站点） | gray-market | markdown |  | 多租户,tenant,siteId,包网,白标,越权,gray-market,fujiang |
| 652 | 通知模板注入（标题/正文模板未转义 → SSTI / 存储型 XSS） | gray-market | markdown |  | 模板,ssti,通知,站内信,短信,注入,gray-market,fujiang |
| 653 | OAuth redirect_uri 劫持 / state 缺失（账号接管） | gray-market | markdown |  | oauth,redirect_uri,state,登录,接管,sso,gray-market,fujiang |
| 654 | 对象存储与数据存储裸露（OSS / S3 / MinIO / MySQL / MongoDB） | gray-market | markdown |  | oss,s3,bucket,minio,mongodb,mysql,弱口令,列举,kyc,gray-market,fujiang |
| 655 | 运营后台弱口令 + 验证码复用（博彩站点最常见入口） | gray-market | markdown |  | 弱口令,后台,验证码复用,运营商,包网,登录,gray-market,fujiang |
| 656 | 订单状态机篡改（未支付改成已支付 / 已取消再入账） | gray-market | markdown |  | 订单,状态机,篡改,未支付,充值,资金,gray-market,fujiang |
| 657 | 源站 IP 泄露 → 绕过 WAF 直打（博彩站必看） | gray-market | markdown |  | 源站,cdn,waf,绕过,真实ip,回源,gray-market,fujiang |
| 658 | 宝塔 / 1Panel 面板弱口令与未授权（黑灰站最常见后台） | gray-market | markdown |  | 宝塔,bt,1panel,面板,弱口令,未授权,服务器,gray-market,fujiang |
| 659 | 密码重置 token 可预测 / 未绑定用户（任意用户改密） | gray-market | markdown |  | reset,token,找回密码,任意用户,越权,gray-market,fujiang |
| 660 | 支付通道 / 商户号篡改（充值进攻击者账户） | gray-market | markdown |  | 支付通道,商户号,充值,路由,篡改,资金,gray-market,fujiang |
| 661 | 支付回调伪造与金额篡改（资金面） | gray-market | markdown |  | notify,callback,伪造回调,金额篡改,重放,验签,跑分,gray-market,fujiang |
| 662 | 活动 / 彩金 / 返水无限领取（幂等缺失） | gray-market | markdown |  | promotion,bonus,彩金,返水,幂等,并发领取,活动滥用,gray-market,fujiang |
| 663 | 余额并发与负数（竞态 / 金额符号） | gray-market | markdown |  | 竞态,race condition,并发提现,负数,金额校验,幂等,gray-market,fujiang |
| 664 | Redis 未授权 → 偷后台 session / token 直接进后台 | gray-market | markdown |  | redis,未授权,session,token,后台,接管,gray-market,fujiang |
| 665 | 源码 / 备份 / 配置泄露 | gray-market | markdown |  | .git,.svn,env,backup,zip,sourcemap,swagger,源码泄露,gray-market,fujiang |
| 666 | 报表/财务接口未授权（全站流水与代理佣金） | gray-market | markdown |  | report,finance,报表,流水,佣金,未授权,gray-market,fujiang |
| 667 | 若依（RuoYi）系后台默认口令 / Druid / 定时任务 RCE | gray-market | markdown |  | ruoyi,若依,druid,admin123,quartz,定时任务,prod-api,gray-market,fujiang |
| 668 | 后台定时任务 / 计划任务注入（Quartz、XXL-JOB、若依） | gray-market | markdown |  | 定时任务,quartz,xxl-job,若依,rce,计划任务,gray-market,fujiang |
| 669 | 短信轰炸与账号枚举（验证码面） | gray-market | markdown |  | sms,轰炸,枚举,验证码爆破,频控,user enumeration,gray-market,fujiang |
| 670 | 短信验证码回显 / 复用 / 爆破（任意账号接管） | gray-market | markdown |  | sms,code,验证码回显,复用,爆破,接管,gray-market,fujiang |
| 671 | sourcemap / 前端 JS 泄露接口与密钥（SPA 包网站第一入口） | gray-market | markdown |  | sourcemap,js,前端,密钥,接口,泄露,spa,gray-market,fujiang |
| 672 | Spring Boot Actuator 未授权（heapdump / Jolokia / Gateway） | gray-market | markdown |  | actuator,heapdump,jolokia,env,gateway,spring,gray-market,fujiang |
| 673 | 业务接口 SQL 注入（报表 / 订单 / 代理查询） | gray-market | markdown |  | sqli,注入,orderNo,agentId,报表,盲注,union,gray-market,fujiang |
| 674 | SSRF（图片抓取 / 回调校验 / 推广链接检测）打内网与云元数据 | gray-market | markdown |  | ssrf,内网,元数据,169.254.169.254,redis,gopher,gray-market,fujiang |
| 675 | 会员资料存储型 XSS 打后台审核（管理员会话劫持） | gray-market | markdown |  | stored xss,存储型,后台,审核,cookie,csrf,gray-market,fujiang |
| 676 | 客服工单附件越权下载（证件照 / KYC 材料） | gray-market | markdown |  | 工单,附件,越权,证件,kyc,下载,gray-market,fujiang |
| 677 | Swagger / Knife4j 接口文档暴露 → 未授权接口直打 | gray-market | markdown |  | swagger,knife4j,springdoc,文档,未授权,api,gray-market,fujiang |
| 678 | 上传点 getshell（头像 / KYC 证件 / 客服附件） | gray-market | markdown |  | upload,getshell,头像,kyc,二次渲染,双扩展名,svg xss,gray-market,fujiang |
| 679 | VIP/等级接口参数篡改（免费拿高等级权益） | gray-market | markdown |  | vip,level,等级,权益,参数篡改,会员,gray-market,fujiang |
| 680 | 提现审核接口越权（读他人银行卡 / 改审核状态） | gray-market | markdown |  | 提现,withdraw,审核,approve,银行卡,越权,gray-market,fujiang |
| 681 | 提现银行卡越权替换（把别人提现收进自己卡） | gray-market | markdown |  | 提现,银行卡,越权,替换,资金,withdraw,gray-market,fujiang |
| 4682 | SAP NetWeaver Remote Code Execution Vulnerability | kev | markdown |  | CVE-2010-5326,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,SAP,N |
| 4683 | HP OpenView Network Node Manager Remote Code Execution Vulnerability | kev | markdown |  | CVE-2005-2773,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Hewle |
| 4684 | Alcatel OmniPCX Enterprise Remote Code Execution Vulnerability | kev | markdown |  | CVE-2007-3010,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Alcat |
| 4685 | phpMyAdmin Remote Code Execution Vulnerability | kev | markdown |  | CVE-2009-1151,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,phpMy |
| 4686 | Oracle JRE Unspecified Vulnerability | kev | markdown |  | CVE-2010-0840,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Oracl |
| 4687 | Adobe ColdFusion Directory Traversal Vulnerability | kev | markdown |  | CVE-2010-2861,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cri |
| 4688 | Exim Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2010-4344,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Exim |
| 4689 | Microsoft Forefront TMG Remote Code Execution Vulnerability | kev | markdown |  | CVE-2011-1889,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Micro |
| 4690 | Adobe Reader and Acrobat Universal 3D Memory Corruption Vulnerability | kev | markdown |  | CVE-2011-2462,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4691 | Oracle Java SE Runtime Environment (JRE) Arbitrary Code Execution Vulnerability | kev | markdown |  | CVE-2011-3544,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Oracl |
| 4692 | Apache Struts 2 Improper Input Validation Vulnerability | kev | markdown |  | CVE-2012-0391,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Apach |
| 4693 | Oracle Fusion Middleware Unspecified Vulnerability | kev | markdown |  | CVE-2012-1710,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cri |
| 4694 | Oracle Java SE Runtime Environment (JRE) Arbitrary Code Execution Vulnerability | kev | markdown |  | CVE-2012-4681,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cri |
| 4695 | Adobe ColdFusion Authentication Bypass Vulnerability | kev | markdown |  | CVE-2013-0625,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4696 | Adobe ColdFusion Authentication Bypass Vulnerability | kev | markdown |  | CVE-2013-0632,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4697 | Apache Struts Improper Input Validation Vulnerability | kev | markdown |  | CVE-2013-2251,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Apach |
| 4698 | Adobe Reader and Acrobat Arbitrary Integer Overflow Vulnerability | kev | markdown |  | CVE-2013-2729,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4699 | Adobe Reader and Acrobat Memory Corruption Vulnerability | kev | markdown |  | CVE-2013-3346,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4700 | HP Multiple Products Remote Code Execution Vulnerability | kev | markdown |  | CVE-2013-4810,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Hewle |
| 4701 | Adobe Flash Player Integer Underflow Vulnerablity | kev | markdown |  | CVE-2014-0497,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4702 | Adobe Reader and Acrobat Sandbox Bypass Vulnerability | kev | markdown |  | CVE-2014-0546,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4703 | InduSoft Web Studio NTWebServer Directory Traversal Vulnerability | kev | markdown |  | CVE-2014-0780,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,InduS |
| 4704 | Multi-Router Looking Glass (MRLG) Buffer Overflow Vulnerability | kev | markdown |  | CVE-2014-3931,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Looki |
| 4705 | GNU Bourne-Again Shell (Bash) Arbitrary Code Execution Vulnerability | kev | markdown |  | CVE-2014-6271,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,GNU,B |
| 4706 | Rejetto HTTP File Server (HFS) Remote Code Execution Vulnerability | kev | markdown |  | CVE-2014-6287,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Rejet |
| 4707 | GNU Bourne-Again Shell (Bash) Arbitrary Code Execution Vulnerability | kev | markdown |  | CVE-2014-7169,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,GNU,B |
| 4708 | Realtek SDK Improper Input Validation Vulnerability | kev | markdown |  | CVE-2014-8361,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Realt |
| 4709 | Adobe Flash Player Remote Code Execution Vulnerability | kev | markdown |  | CVE-2015-0311,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4710 | Adobe Flash Player Use-After-Free Vulnerability | kev | markdown |  | CVE-2015-0313,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4711 | D-Link and TRENDnet Multiple Devices Remote Code Execution Vulnerability | kev | markdown |  | CVE-2015-1187,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,D-Lin |
| 4712 | Elasticsearch Groovy Scripting Engine Remote Code Execution Vulnerability | kev | markdown |  | CVE-2015-1427,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Elast |
| 4713 | Oracle Java SE and Java SE Embedded Remote Code Execution Vulnerability | kev | markdown |  | CVE-2015-2590,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Oracl |
| 4714 | Adobe Flash Player Memory Corruption Vulnerability | kev | markdown |  | CVE-2015-3043,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4715 | Adobe Flash Player Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2015-3113,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4716 | Oracle WebLogic Server Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2015-4852,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Oracl |
| 4717 | Adobe Flash Player Use-After-Free Vulnerability | kev | markdown |  | CVE-2015-5119,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4718 | Adobe Flash Player Use-After-Free Vulnerability | kev | markdown |  | CVE-2015-5122,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4719 | Adobe Flash Player Use-After-Free Vulnerability | kev | markdown |  | CVE-2015-5123,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4720 | IBM WebSphere Application Server and Server Hypervisor Edition Code Injection. | kev | markdown |  | CVE-2015-7450,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,IBM,W |
| 4721 | Juniper ScreenOS Improper Authentication Vulnerability | kev | markdown |  | CVE-2015-7755,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Junip |
| 4722 | PHPMailer Command Injection Vulnerability | kev | markdown |  | CVE-2016-10033,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,PHP, |
| 4723 | NETGEAR WNR2000v5 Router Buffer Overflow Vulnerability | kev | markdown |  | CVE-2016-10174,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,NETG |
| 4724 | Adobe Flash Player Arbitrary Code Execution Vulnerability | kev | markdown |  | CVE-2016-1019,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cri |
| 4725 | NETGEAR Multiple WAP Devices Command Injection Vulnerability | kev | markdown |  | CVE-2016-1555,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,NETGE |
| 4726 | SAP NetWeaver SQL Injection Vulnerability | kev | markdown |  | CVE-2016-2386,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,SAP,N |
| 4727 | Apache ActiveMQ Improper Input Validation Vulnerability | kev | markdown |  | CVE-2016-3088,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Apach |
| 4728 | Oracle Java SE and JRockit Unspecified Vulnerability | kev | markdown |  | CVE-2016-3427,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Oracl |
| 4729 | Adobe Flash Player Arbitrary Code Execution Vulnerability | kev | markdown |  | CVE-2016-4117,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cri |
| 4730 | Adobe Flash Player Remote Code Execution Vulnerability | kev | markdown |  | CVE-2016-4171,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4731 | Apache Shiro Code Execution Vulnerability | kev | markdown |  | CVE-2016-4437,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Apach |
| 4732 | SKYSEA Client View Improper Authentication Vulnerability | kev | markdown |  | CVE-2016-7836,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,SKYSE |
| 4733 | Apache Tomcat Remote Code Execution Vulnerability | kev | markdown |  | CVE-2016-8735,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Apach |
| 4734 | Jenkins Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-1000353,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Je |
| 4735 | Primetek Primefaces Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-1000486,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Pr |
| 4736 | Telerik UI for ASP.NET AJAX Unrestricted File Upload Vulnerability | kev | markdown |  | CVE-2017-11317,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Tele |
| 4737 | Telerik UI for ASP.NET AJAX Insecure Direct Object Reference Vulnerability | kev | markdown |  | CVE-2017-11357,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cr |
| 4738 | Red Hat JBoss Application Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-12149,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cr |
| 4739 | Cisco IOS and IOS XE Software DHCP Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-12240,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Cisc |
| 4740 | Palo Alto Networks PAN-OS Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-15944,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Palo |
| 4741 | Kaseya VSA SQL Injection Vulnerability | kev | markdown |  | CVE-2017-18362,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cr |
| 4742 | Adobe ColdFusion Deserialization Vulnerability | kev | markdown |  | CVE-2017-3066,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4743 | Cisco IOS and IOS XE Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-3881,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Cisco |
| 4744 | Apache Struts Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-5638,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cri |
| 4745 | Intel Active Management Technology (AMT), Small Business Technology (SBT), and Standard Manageability Privilege Escalation Vulnerability | kev | markdown |  | CVE-2017-5689,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Intel |
| 4746 | NETGEAR DGN2200 Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-6077,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,NETGE |
| 4747 | Citrix Multiple Products Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-6316,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Citri |
| 4748 | NETGEAR Multiple Devices Buffer Overflow Vulnerability | kev | markdown |  | CVE-2017-6862,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,NETGE |
| 4749 | Microsoft Windows Server Buffer Overflow Vulnerability | kev | markdown |  | CVE-2017-7269,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Micro |
| 4750 | Samba Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-7494,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cri |
| 4751 | Hikvision Multiple Products Improper Authentication Vulnerability | kev | markdown |  | CVE-2017-7921,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Hikvi |
| 4752 | Microsoft Windows Search Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-8543,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Micro |
| 4753 | Progress Telerik UI for ASP.NET AJAX and Sitefinity Cryptographic Weakness Vulnerability | kev | markdown |  | CVE-2017-9248,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Progr |
| 4754 | Apache Struts 1 Improper Input Validation Vulnerability | kev | markdown |  | CVE-2017-9791,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Apach |
| 4755 | PHPUnit Command Injection Vulnerability | kev | markdown |  | CVE-2017-9841,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,PHPUn |
| 4756 | Cisco VPN Routers Remote Code Execution Vulnerability | kev | markdown |  | CVE-2018-0125,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Cisco |
| 4757 | Cisco Secure Access Control System Java Deserialization Vulnerability | kev | markdown |  | CVE-2018-0147,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Cisco |
| 4758 | Cisco IOS Software and Cisco IOS XE Software Quality of Service Remote Code Execution Vulnerability | kev | markdown |  | CVE-2018-0151,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Cisco |
| 4759 | Cisco IOS and IOS XE Software Smart Install Remote Code Execution Vulnerability | kev | markdown |  | CVE-2018-0171,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Cisco |
| 4760 | Jenkins Stapler Web Framework Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2018-1000861,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Je |
| 4761 | Dasan GPON Routers Authentication Bypass Vulnerability | kev | markdown |  | CVE-2018-10561,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Dasa |
| 4762 | Dasan GPON Routers Command Injection Vulnerability | kev | markdown |  | CVE-2018-10562,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cr |
| 4763 | Quest KACE System Management Appliance Remote Command Execution Vulnerability | kev | markdown |  | CVE-2018-11138,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cr |
| 4764 | VMware Tanzu Spring Data Commons Property Binder Vulnerability | kev | markdown |  | CVE-2018-1273,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cri |
| 4765 | Tenda AC7, AC9, and AC10 Routers Command Injection Vulnerability | kev | markdown |  | CVE-2018-14558,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Tend |
| 4766 | Red Hat JBoss RichFaces Framework Expression Language Injection Vulnerability | kev | markdown |  | CVE-2018-14667,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Red  |
| 4767 | NUUO NVRmini Devices OS Command Injection Vulnerability | kev | markdown |  | CVE-2018-14933,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,NUUO |
| 4768 | Adobe ColdFusion Unrestricted File Upload Vulnerability | kev | markdown |  | CVE-2018-15961,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adob |
| 4769 | GIGABYTE Multiple Products Privilege Escalation Vulnerability | kev | markdown |  | CVE-2018-19323,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cr |
| 4770 | Paessler PRTG Network Monitor Local File Inclusion Vulnerability | kev | markdown |  | CVE-2018-19410,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Paes |
| 4771 | ThinkPHP "noneCms" Remote Code Execution Vulnerability | kev | markdown |  | CVE-2018-20062,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Thin |
| 4772 | Kaseya VSA Remote Code Execution Vulnerability | kev | markdown |  | CVE-2018-20753,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cr |
| 4773 | Oracle WebLogic Server Unspecified Vulnerability | kev | markdown |  | CVE-2018-2628,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Oracl |
| 4774 | Adobe ColdFusion Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2018-4939,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Adobe |
| 4775 | D-Link Multiple Routers OS Command Injection Vulnerability | kev | markdown |  | CVE-2018-6530,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cri |
| 4776 | Exim Buffer Overflow Vulnerability | kev | markdown |  | CVE-2018-6789,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cri |
| 4777 | MikroTik RouterOS Stack-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2018-7445,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Mikro |
| 4778 | Drupal Core Remote Code Execution Vulnerability | kev | markdown |  | CVE-2018-7600,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cri |
| 4779 | Drupal Core Remote Code Execution Vulnerability | kev | markdown |  | CVE-2018-7602,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,cvss-cri |
| 4780 | Oracle Fusion Middleware Unspecified Vulnerability | kev | markdown |  | CVE-2012-3152,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Oracl |
| 4781 | Arcserve Unified Data Protection (UDP) Directory Traversal Vulnerability | kev | markdown |  | CVE-2015-4068,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Arcse |
| 4782 | MikroTik Router OS Directory Traversal Vulnerability | kev | markdown |  | CVE-2018-14847,cve-arsenal,cisa-kev,kev,in-the-wild,cvss-critical,Mikr |
| 4783 | Microsoft Windows Privilege Escalation Vulnerability | kev | markdown |  | CVE-2002-0367,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4784 | Microsoft Windows Privilege Escalation Vulnerability | kev | markdown |  | CVE-2004-0210,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4785 | Cisco IOS Denial-of-Service Vulnerability | kev | markdown |  | CVE-2004-1464,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS |
| 4786 | Apache Struts 1 ActionForm Denial-of-Service Vulnerability | kev | markdown |  | CVE-2006-1547,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Struts 1 |
| 4787 | Microsoft Word Malformed Object Pointer Vulnerability | kev | markdown |  | CVE-2006-2492,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Word |
| 4788 | Microsoft Office Excel Remote Code Execution Vulnerability | kev | markdown |  | CVE-2007-0671,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4789 | Adobe Acrobat and Reader Buffer Overflow Vulnerability | kev | markdown |  | CVE-2007-5659,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat and R |
| 4790 | Microsoft Windows Video ActiveX Control Remote Code Execution Vulnerability | kev | markdown |  | CVE-2008-0015,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4791 | Adobe Acrobat and Reader Unspecified Vulnerability | kev | markdown |  | CVE-2008-0655,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat and R |
| 4792 | Adobe Reader and Acrobat Input Validation Vulnerability | kev | markdown |  | CVE-2008-2992,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Adobe,Ac |
| 4793 | Oracle VirtualBox Insufficient Input Validation Vulnerability | kev | markdown |  | CVE-2008-3431,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,VirtualBox |
| 4794 | Cisco IOS Cross-Site Request Forgery Vulnerability | kev | markdown |  | CVE-2008-4128,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS |
| 4795 | Microsoft Windows Buffer Overflow Vulnerability | kev | markdown |  | CVE-2008-4250,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4796 | Microsoft Office Remote Code Execution | kev | markdown |  | CVE-2009-0238,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4797 | Microsoft Office PowerPoint Code Injection Vulnerability | kev | markdown |  | CVE-2009-0556,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4798 | Microsoft Office Object Record Corruption Vulnerability | kev | markdown |  | CVE-2009-0557,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4799 | Microsoft Office Buffer Overflow Vulnerability | kev | markdown |  | CVE-2009-0563,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4800 | Adobe Reader and Adobe Acrobat Stack-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2009-0927,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Reader and Ac |
| 4801 | Microsoft Windows Improper Input Validation Vulnerability | kev | markdown |  | CVE-2009-1123,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4802 | Microsoft DirectX NULL Byte Overwrite Vulnerability | kev | markdown |  | CVE-2009-1537,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,DirectX |
| 4803 | Adobe Acrobat and Reader, Flash Player Unspecified Vulnerability | kev | markdown |  | CVE-2009-1862,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat and R |
| 4804 | Cisco IOS XR Border Gateway Protocol (BGP) Denial-of-Service Vulnerability | kev | markdown |  | CVE-2009-2055,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS XR |
| 4805 | Microsoft Excel Featheader Record Memory Corruption Vulnerability | kev | markdown |  | CVE-2009-3129,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Excel |
| 4806 | Adobe Acrobat and Reader Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2009-3459,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat and R |
| 4807 | Adobe Acrobat and Reader Universal 3D Remote Code Execution Vulnerability | kev | markdown |  | CVE-2009-3953,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat and R |
| 4808 | Adobe BlazeDS Information Disclosure Vulnerability | kev | markdown |  | CVE-2009-3960,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Adobe,Bl |
| 4809 | Adobe Acrobat and Reader Use-After-Free Vulnerability | kev | markdown |  | CVE-2009-4324,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat and R |
| 4810 | Adobe Reader and Acrobat Arbitrary Code Execution Vulnerability | kev | markdown |  | CVE-2010-0188,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Adobe,Re |
| 4811 | Microsoft Windows Kernel Exception Handler Vulnerability | kev | markdown |  | CVE-2010-0232,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4812 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | markdown |  | CVE-2010-0249,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4813 | Red Hat JBoss Authentication Bypass Vulnerability | kev | markdown |  | CVE-2010-0738,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Red Hat, |
| 4814 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | markdown |  | CVE-2010-0806,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4815 | Adobe Flash Player Memory Corruption Vulnerability | kev | markdown |  | CVE-2010-1297,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4816 | Red Hat JBoss Information Disclosure Vulnerability | kev | markdown |  | CVE-2010-1428,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Red Hat, |
| 4817 | Red Hat Linux JBoss Seam 2 Remote Code Execution Vulnerability | kev | markdown |  | CVE-2010-1871,cve-arsenal,cisa-kev,kev,in-the-wild,Red Hat,JBoss Seam  |
| 4818 | Microsoft Windows Remote Code Execution Vulnerability | kev | markdown |  | CVE-2010-2568,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4819 | Microsoft PowerPoint Buffer Overflow Vulnerability | kev | markdown |  | CVE-2010-2572,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,PowerPoin |
| 4820 | Adobe Acrobat and Reader Stack-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2010-2883,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat and R |
| 4821 | Cisco IOS XR Border Gateway Protocol (BGP) Denial-of-Service Vulnerability | kev | markdown |  | CVE-2010-3035,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS XR |
| 4822 | Microsoft Office Stack-based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2010-3333,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4823 | Mozilla Multiple Products Remote Code Execution Vulnerability | kev | markdown |  | CVE-2010-3765,cve-arsenal,cisa-kev,kev,in-the-wild,Mozilla,Multiple Pr |
| 4824 | Linux Kernel Improper Input Validation Vulnerability | kev | markdown |  | CVE-2010-3904,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 4825 | Microsoft Internet Explorer Uninitialized Memory Corruption Vulnerability | kev | markdown |  | CVE-2010-3962,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4826 | Exim Privilege Escalation Vulnerability | kev | markdown |  | CVE-2010-4345,cve-arsenal,cisa-kev,kev,in-the-wild,Exim |
| 4827 | Microsoft Windows Kernel Stack-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2010-4398,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4828 | Ubiquiti AirOS Command Injection Vulnerability | kev | markdown |  | CVE-2010-5330,cve-arsenal,cisa-kev,kev,in-the-wild,Ubiquiti,AirOS |
| 4829 | Adobe Flash Player Unspecified Vulnerability | kev | markdown |  | CVE-2011-0609,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4830 | Adobe Flash Player Remote Code Execution Vulnerability | kev | markdown |  | CVE-2011-0611,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4831 | Android OS Privilege Escalation Vulnerability | kev | markdown |  | CVE-2011-1823,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Android OS |
| 4832 | Microsoft Ancillary Function Driver (afd.sys) Improper Input Validation Vulnerability | kev | markdown |  | CVE-2011-2005,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Ancillary |
| 4833 | Microsoft Windows Remote Code Execution Vulnerability | kev | markdown |  | CVE-2011-3402,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4834 | D-Link DIR-300 Router Cleartext Storage of a Password Vulnerability | kev | markdown |  | CVE-2011-4723,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DIR-300 Rout |
| 4835 | Microsoft Windows Authenticode Signature Verification Remote Code Execution Vulnerability | kev | markdown |  | CVE-2012-0151,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4836 | Microsoft MSCOMCTL.OCX Remote Code Execution Vulnerability | kev | markdown |  | CVE-2012-0158,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4837 | Oracle Java SE Runtime Environment (JRE) Arbitrary Code Execution Vulnerability | kev | markdown |  | CVE-2012-0507,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Oracle,J |
| 4838 | Oracle Fusion Middleware Unspecified Vulnerability | kev | markdown |  | CVE-2012-0518,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,Fusion Middl |
| 4839 | Adobe Flash Player Memory Corruption Vulnerability | kev | markdown |  | CVE-2012-0754,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4840 | Adobe Flash Player Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2012-0767,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4841 | Adobe Flash Player Arbitrary Code Execution Vulnerability | kev | markdown |  | CVE-2012-1535,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4842 | Oracle Java SE Runtime Environment (JRE) Arbitrary Code Execution Vulnerability | kev | markdown |  | CVE-2012-1723,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Oracle,J |
| 4843 | PHP-CGI Query String Parameter Vulnerability | kev | markdown |  | CVE-2012-1823,cve-arsenal,cisa-kev,kev,in-the-wild,PHP |
| 4844 | Microsoft Visual Basic for Applications Insecure Library Loading Vulnerability | kev | markdown |  | CVE-2012-1854,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Visual Ba |
| 4845 | Microsoft Office MSCOMCTL.OCX Remote Code Execution Vulnerability | kev | markdown |  | CVE-2012-1856,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4846 | Microsoft XML Core Services Memory Corruption Vulnerability | kev | markdown |  | CVE-2012-1889,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,XML Core  |
| 4847 | Adobe Flash Player Memory Corruption Vulnerability | kev | markdown |  | CVE-2012-2034,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4848 | Microsoft Word Remote Code Execution Vulnerability | kev | markdown |  | CVE-2012-2539,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Word |
| 4849 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | markdown |  | CVE-2012-4792,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4850 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | markdown |  | CVE-2012-4969,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4851 | Adobe Flash Player Integer Overflow Vulnerability | kev | markdown |  | CVE-2012-5054,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4852 | Oracle Java SE Sandbox Bypass Vulnerability | kev | markdown |  | CVE-2012-5076,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,Java SE |
| 4853 | Microsoft Silverlight Double Dereference Vulnerability | kev | markdown |  | CVE-2013-0074,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4854 | Oracle JRE Remote Code Execution Vulnerability | kev | markdown |  | CVE-2013-0422,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Oracle,J |
| 4855 | Oracle JRE Sandbox Bypass Vulnerability | kev | markdown |  | CVE-2013-0431,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Oracle,J |
| 4856 | Adobe ColdFusion Directory Traversal Vulnerability | kev | markdown |  | CVE-2013-0629,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,ColdFusion |
| 4857 | Adobe ColdFusion Information Disclosure Vulnerability | kev | markdown |  | CVE-2013-0631,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,ColdFusion |
| 4858 | Adobe Reader and Acrobat Memory Corruption Vulnerability | kev | markdown |  | CVE-2013-0640,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Reader and Ac |
| 4859 | Adobe Reader Buffer Overflow Vulnerability | kev | markdown |  | CVE-2013-0641,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Reader |
| 4860 | Adobe Flash Player Incorrect Default Permissions Vulnerability | kev | markdown |  | CVE-2013-0643,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4861 | Adobe Flash Player Code Execution Vulnerability | kev | markdown |  | CVE-2013-0648,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4862 | Microsoft Office Buffer Overflow Vulnerability | kev | markdown |  | CVE-2013-1331,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4863 | Microsoft Internet Explorer Remote Code Execution Vulnerability | kev | markdown |  | CVE-2013-1347,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4864 | Mozilla Firefox Information Disclosure Vulnerability | kev | markdown |  | CVE-2013-1675,cve-arsenal,cisa-kev,kev,in-the-wild,Mozilla,Firefox |
| 4865 | Mozilla Firefox and Thunderbird Denial-of-Service Vulnerability | kev | markdown |  | CVE-2013-1690,cve-arsenal,cisa-kev,kev,in-the-wild,Mozilla,Firefox and |
| 4866 | Linux Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2013-2094,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 4867 | Oracle JRE Unspecified Vulnerability | kev | markdown |  | CVE-2013-2423,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,Java Runtime |
| 4868 | Oracle Java SE Unspecified Vulnerability | kev | markdown |  | CVE-2013-2465,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Oracle,J |
| 4869 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | markdown |  | CVE-2013-2551,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4870 | Linux Kernel Integer Overflow Vulnerability | kev | markdown |  | CVE-2013-2596,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 4871 | Code Aurora ACDB Audio Driver Stack-based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2013-2597,cve-arsenal,cisa-kev,kev,in-the-wild,Code Aurora,ACDB Au |
| 4872 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | markdown |  | CVE-2013-3163,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4873 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2013-3660,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 4874 | Microsoft Internet Explorer Resource Management Errors Vulnerability | kev | markdown |  | CVE-2013-3893,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4875 | Microsoft Silverlight Information Disclosure Vulnerability | kev | markdown |  | CVE-2013-3896,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Silverlig |
| 4876 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | markdown |  | CVE-2013-3897,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4877 | Microsoft WinVerifyTrust function Remote Code Execution | kev | markdown |  | CVE-2013-3900,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,WinVerify |
| 4878 | Microsoft Graphics Component Memory Corruption Vulnerability | kev | markdown |  | CVE-2013-3906,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Graphics  |
| 4879 | Microsoft Windows Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2013-3918,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4880 | IBM InfoSphere BigInsights Invalid Input Vulnerability | kev | markdown |  | CVE-2013-3993,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,IBM,Info |
| 4881 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2013-5065,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4882 | D-Link DSL-2760U Gateway Cross-Site Scripting Vulnerability | kev | markdown |  | CVE-2013-5223,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DSL-2760U |
| 4883 | Linux Kernel Improper Input Validation Vulnerability | kev | markdown |  | CVE-2013-6282,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 4884 | Microsoft Internet Explorer Information Disclosure Vulnerability | kev | markdown |  | CVE-2013-7331,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4885 | Ruby on Rails Directory Traversal Vulnerability | kev | markdown |  | CVE-2014-0130,cve-arsenal,cisa-kev,kev,in-the-wild,Rails,Ruby on Rails |
| 4886 | OpenSSL Information Disclosure Vulnerability | kev | markdown |  | CVE-2014-0160,cve-arsenal,cisa-kev,kev,in-the-wild,OpenSSL |
| 4887 | Linux Kernel Race Condition Vulnerability | kev | markdown |  | CVE-2014-0196,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 4888 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | markdown |  | CVE-2014-0322,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4889 | Adobe Reader and Acrobat Use-After-Free Vulnerability | kev | markdown |  | CVE-2014-0496,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Reader and Ac |
| 4890 | Adobe Flash Player Double Free Vulnerablity | kev | markdown |  | CVE-2014-0502,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4891 | D-Link DIR-600 Router Cross-Site Request Forgery (CSRF) Vulnerability | kev | markdown |  | CVE-2014-100005,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DIR-600 Ro |
| 4892 | Microsoft Word Memory Corruption Vulnerability | kev | markdown |  | CVE-2014-1761,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Word |
| 4893 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | markdown |  | CVE-2014-1776,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4894 | Microsoft Windows Group Policy Preferences Password Privilege Escalation Vulnerability | kev | markdown |  | CVE-2014-1812,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4895 | Cisco Adaptive Security Appliance (ASA) Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2014-2120,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Adaptive Secu |
| 4896 | Microsoft Internet Explorer Privilege Escalation Vulnerability | kev | markdown |  | CVE-2014-2817,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4897 | Elasticsearch Remote Code Execution Vulnerability | kev | markdown |  | CVE-2014-3120,cve-arsenal,cisa-kev,kev,in-the-wild,Elastic,Elasticsear |
| 4898 | Linux Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2014-3153,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 4899 | Microsoft IME Japanese Privilege Escalation Vulnerability | kev | markdown |  | CVE-2014-4077,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Input Met |
| 4900 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2014-4113,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 4901 | Microsoft Windows Object Linking & Embedding (OLE) Remote Code Execution Vulnerability | kev | markdown |  | CVE-2014-4114,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4902 | Microsoft Internet Explorer Privilege Escalation Vulnerability | kev | markdown |  | CVE-2014-4123,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4903 | Microsoft Windows Remote Code Execution Vulnerability | kev | markdown |  | CVE-2014-4148,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4904 | Apple OS X Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2014-4404,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,OS X |
| 4905 | GNU Bash OS Command Injection Vulnerability | kev | markdown |  | CVE-2014-6278,cve-arsenal,cisa-kev,kev,in-the-wild,GNU,GNU Bash |
| 4906 | Microsoft Kerberos Key Distribution Center (KDC) Privilege Escalation Vulnerability | kev | markdown |  | CVE-2014-6324,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Kerberos  |
| 4907 | Microsoft Windows Object Linking & Embedding (OLE) Automation Array Remote Code Execution Vulnerability | kev | markdown |  | CVE-2014-6332,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4908 | Microsoft Windows Code Injection Vulnerability | kev | markdown |  | CVE-2014-6352,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4909 | Adobe Flash Player Dereferenced Pointer Vulnerability | kev | markdown |  | CVE-2014-8439,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4910 | Adobe Flash Player Stack-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2014-9163,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4911 | Microsoft Windows TS WebProxy Directory Traversal Vulnerability | kev | markdown |  | CVE-2015-0016,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4912 | Microsoft Internet Explorer ASLR Bypass Vulnerability | kev | markdown |  | CVE-2015-0071,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4913 | Adobe Flash Player ASLR Bypass Vulnerability | kev | markdown |  | CVE-2015-0310,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4914 | Cisco Prime Data Center Network Manager (DCNM) Directory Traversal Vulnerability | kev | markdown |  | CVE-2015-0666,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Prime Data Ce |
| 4915 | Apple OS X Authentication Bypass Vulnerability | kev | markdown |  | CVE-2015-1130,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,OS X |
| 4916 | Microsoft HTTP.sys Remote Code Execution Vulnerability | kev | markdown |  | CVE-2015-1635,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,HTTP.sys |
| 4917 | Microsoft Office Memory Corruption Vulnerability | kev | markdown |  | CVE-2015-1641,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4918 | Microsoft Office Memory Corruption Vulnerability | kev | markdown |  | CVE-2015-1642,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4919 | Microsoft Windows Remote Code Execution Vulnerability | kev | markdown |  | CVE-2015-1671,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4920 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2015-1701,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4921 | Microsoft Windows Mount Manager Privilege Escalation Vulnerability | kev | markdown |  | CVE-2015-1769,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4922 | Microsoft Office Uninitialized Memory Use Vulnerability | kev | markdown |  | CVE-2015-1770,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4923 | D-Link DIR-645 Router Remote Code Execution Vulnerability | kev | markdown |  | CVE-2015-2051,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DIR-645 Rout |
| 4924 | Intel Ethernet Diagnostics Driver for Windows Denial-of-Service Vulnerability | kev | markdown |  | CVE-2015-2291,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Intel,Et |
| 4925 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2015-2360,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 4926 | Microsoft ATM Font Driver Privilege Escalation Vulnerability | kev | markdown |  | CVE-2015-2387,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,ATM Font  |
| 4927 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | markdown |  | CVE-2015-2419,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4928 | Microsoft PowerPoint Memory Corruption Vulnerability | kev | markdown |  | CVE-2015-2424,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,PowerPoin |
| 4929 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | markdown |  | CVE-2015-2425,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4930 | Microsoft Windows Adobe Type Manager Library Remote Code Execution Vulnerability | kev | markdown |  | CVE-2015-2426,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4931 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | markdown |  | CVE-2015-2502,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4932 | Microsoft Office Malformed EPS File Vulnerability | kev | markdown |  | CVE-2015-2545,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4933 | Microsoft Win32k Memory Corruption Vulnerability | kev | markdown |  | CVE-2015-2546,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4934 | TP-Link Multiple Archer Devices Directory Traversal Vulnerability | kev | markdown |  | CVE-2015-3035,cve-arsenal,cisa-kev,kev,in-the-wild,TP-Link,Multiple Ar |
| 4935 | Red Hat Libuser Race Condition Vulnerability | kev | markdown |  | CVE-2015-3246,cve-arsenal,cisa-kev,kev,in-the-wild,Red Hat,Libuser |
| 4936 | ProFTPD Improper Access Control Vulnerability | kev | markdown |  | CVE-2015-3306,cve-arsenal,cisa-kev,kev,in-the-wild,ProFTPD |
| 4937 | Mozilla Firefox Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2015-4495,cve-arsenal,cisa-kev,kev,in-the-wild,Mozilla,Firefox |
| 4938 | Oracle Java SE Integrity Check Vulnerability | kev | markdown |  | CVE-2015-4902,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,Java SE |
| 4939 | Red Hat Automatic Bug Reporting Tool Privilege Escalation Vulnerability | kev | markdown |  | CVE-2015-5287,cve-arsenal,cisa-kev,kev,in-the-wild,Red Hat,Automatic B |
| 4940 | Jenkins User Interface (UI) Information Disclosure Vulnerability | kev | markdown |  | CVE-2015-5317,cve-arsenal,cisa-kev,kev,in-the-wild,Jenkins,Jenkins Use |
| 4941 | ISC BIND Data Processing Errors Vulnerability | kev | markdown |  | CVE-2015-5477,cve-arsenal,cisa-kev,kev,in-the-wild,ISC,BIND |
| 4942 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2015-6175,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4943 | Adobe Flash Player Arbitrary Code Execution Vulnerability | kev | markdown |  | CVE-2015-7645,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Adobe,Fl |
| 4944 | Adobe Flash Player Integer Overflow Vulnerability | kev | markdown |  | CVE-2015-8651,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4945 | Microsoft Silverlight Runtime Remote Code Execution Vulnerability | kev | markdown |  | CVE-2016-0034,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4946 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2016-0040,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4947 | Microsoft Windows Secondary Logon Service Privilege Escalation Vulnerability | kev | markdown |  | CVE-2016-0099,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4948 | Microsoft Windows CSRSS Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2016-0151,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4949 | Microsoft Internet Explorer Information Disclosure Vulnerability | kev | markdown |  | CVE-2016-0162,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4950 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2016-0165,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 4951 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2016-0167,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4952 | Microsoft Windows Media Center Remote Code Execution Vulnerability | kev | markdown |  | CVE-2016-0185,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4953 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | markdown |  | CVE-2016-0189,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4954 | Ruby on Rails Directory Traversal Vulnerability | kev | markdown |  | CVE-2016-0752,cve-arsenal,cisa-kev,kev,in-the-wild,Rails,Ruby on Rails |
| 4955 | Adobe Flash Player and AIR Use-After-Free Vulnerability | kev | markdown |  | CVE-2016-0984,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player  |
| 4956 | Adobe Flash Player and AIR Integer Overflow Vulnerability | kev | markdown |  | CVE-2016-1010,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player  |
| 4957 | D-Link DCS-930L Devices OS Command Injection Vulnerability | kev | markdown |  | CVE-2016-11021,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DCS-930L De |
| 4958 | Google Chromium V8 Out-of-Bounds Read Vulnerability | kev | markdown |  | CVE-2016-1646,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 4959 | D-Link DSL-2750B Devices Command Injection Vulnerability | kev | markdown |  | CVE-2016-20017,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DSL-2750B D |
| 4960 | SAP NetWeaver Information Disclosure Vulnerability | kev | markdown |  | CVE-2016-2388,cve-arsenal,cisa-kev,kev,in-the-wild,SAP,NetWeaver |
| 4961 | Apache Struts Command Injection Vulnerability | kev | markdown |  | CVE-2016-3081,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Struts |
| 4962 | Microsoft Office OLE DLL Side Loading Vulnerability | kev | markdown |  | CVE-2016-3235,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4963 | Microsoft Internet Explorer Messaging API Information Disclosure Vulnerability | kev | markdown |  | CVE-2016-3298,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4964 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2016-3309,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4965 | Microsoft Internet Explorer and Edge Information Disclosure Vulnerability | kev | markdown |  | CVE-2016-3351,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4966 | Microsoft Windows Graphics Device Interface (GDI) Remote Code Execution Vulnerability | kev | markdown |  | CVE-2016-3393,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4967 | SolarWinds Virtualization Manager Privilege Escalation Vulnerability | kev | markdown |  | CVE-2016-3643,cve-arsenal,cisa-kev,kev,in-the-wild,SolarWinds,Virtuali |
| 4968 | ImageMagick Improper Input Validation Vulnerability | kev | markdown |  | CVE-2016-3714,cve-arsenal,cisa-kev,kev,in-the-wild,ImageMagick |
| 4969 | ImageMagick Arbitrary File Deletion Vulnerability | kev | markdown |  | CVE-2016-3715,cve-arsenal,cisa-kev,kev,in-the-wild,ImageMagick |
| 4970 | ImageMagick Server-Side Request Forgery (SSRF) Vulnerability | kev | markdown |  | CVE-2016-3718,cve-arsenal,cisa-kev,kev,in-the-wild,ImageMagick |
| 4971 | SAP NetWeaver Directory Traversal Vulnerability | kev | markdown |  | CVE-2016-3976,cve-arsenal,cisa-kev,kev,in-the-wild,SAP,NetWeaver |
| 4972 | Trihedral VTScada (formerly VTS) Denial-of-Service Vulnerability | kev | markdown |  | CVE-2016-4523,cve-arsenal,cisa-kev,kev,in-the-wild,Trihedral,VTScada ( |
| 4973 | Apple iOS Information Disclosure Vulnerability | kev | markdown |  | CVE-2016-4655,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS |
| 4974 | Apple iOS Memory Corruption Vulnerability | kev | markdown |  | CVE-2016-4656,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS |
| 4975 | Apple iOS Webkit Memory Corruption Vulnerability | kev | markdown |  | CVE-2016-4657,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS |
| 4976 | Linux Kernel Race Condition Vulnerability | kev | markdown |  | CVE-2016-5195,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 4977 | Google Chromium V8 Out-of-Bounds Memory Vulnerability | kev | markdown |  | CVE-2016-5198,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 4978 | NETGEAR Multiple Routers Remote Code Execution Vulnerability | kev | markdown |  | CVE-2016-6277,cve-arsenal,cisa-kev,kev,in-the-wild,NETGEAR,Multiple Ro |
| 4979 | Cisco Adaptive Security Appliance (ASA) SNMP Buffer Overflow Vulnerability | kev | markdown |  | CVE-2016-6366,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Adaptive Secu |
| 4980 | Cisco Adaptive Security Appliance (ASA) CLI Remote Code Execution Vulnerability | kev | markdown |  | CVE-2016-6367,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Adaptive Secu |
| 4981 | Cisco IOS, IOS XR, and IOS XE IKEv1 Information Disclosure Vulnerability | kev | markdown |  | CVE-2016-6415,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS, IOS XR,  |
| 4982 | Microsoft Office Memory Corruption Vulnerability | kev | markdown |  | CVE-2016-7193,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 4983 | Microsoft Edge Memory Corruption Vulnerability | kev | markdown |  | CVE-2016-7200,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Edge |
| 4984 | Microsoft Edge Memory Corruption Vulnerability | kev | markdown |  | CVE-2016-7201,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Edge |
| 4985 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2016-7255,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4986 | Microsoft Windows Open Type Font Remote Code Execution Vulnerability | kev | markdown |  | CVE-2016-7256,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4987 | Microsoft Office Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2016-7262,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Excel |
| 4988 | Adobe Flash Player Use-After-Free Vulnerability | kev | markdown |  | CVE-2016-7855,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4989 | Adobe Flash Player Use-After-Free Vulnerability | kev | markdown |  | CVE-2016-7892,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 4990 | Siemens SIMATIC CP 1543-1 Improper Privilege Management Vulnerability | kev | markdown |  | CVE-2016-8562,cve-arsenal,cisa-kev,kev,in-the-wild,Siemens,SIMATIC CP |
| 4991 | Mozilla Firefox, Firefox ESR, and Thunderbird Use-After-Free Vulnerability | kev | markdown |  | CVE-2016-9079,cve-arsenal,cisa-kev,kev,in-the-wild,Mozilla,Firefox, Fi |
| 4992 | SAP NetWeaver XML External Entity (XXE) Vulnerability | kev | markdown |  | CVE-2016-9563,cve-arsenal,cisa-kev,kev,in-the-wild,SAP,NetWeaver |
| 4993 | Microsoft Graphics Device Interface (GDI) Privilege Escalation Vulnerability | kev | markdown |  | CVE-2017-0001,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Graphics  |
| 4994 | Microsoft Windows Graphics Device Interface (GDI) Privilege Escalation Vulnerability | kev | markdown |  | CVE-2017-0005,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 4995 | Microsoft XML Core Services Information Disclosure Vulnerability | kev | markdown |  | CVE-2017-0022,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,XML Core  |
| 4996 | Microsoft Edge and Internet Explorer Type Confusion Vulnerability | kev | markdown |  | CVE-2017-0037,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Edge and  |
| 4997 | Microsoft Internet Explorer Information Disclosure Vulnerability | kev | markdown |  | CVE-2017-0059,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 4998 | Microsoft Windows Transaction Manager Privilege Escalation Vulnerability | kev | markdown |  | CVE-2017-0101,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 4999 | Microsoft Windows Server Message Block (SMBv1) Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-0143,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5000 | Microsoft SMBv1 Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-0144,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5001 | Microsoft SMBv1 Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-0145,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5002 | Microsoft Windows SMB Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-0146,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5003 | Microsoft Windows SMBv1 Information Disclosure Vulnerability | kev | markdown |  | CVE-2017-0147,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5004 | Microsoft SMBv1 Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-0148,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5005 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | markdown |  | CVE-2017-0149,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 5006 | Microsoft Office and WordPad Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-0199,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5007 | Microsoft Internet Explorer Privilege Escalation Vulnerability | kev | markdown |  | CVE-2017-0210,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 5008 | Microsoft Windows Privilege Escalation Vulnerability | kev | markdown |  | CVE-2017-0213,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5009 | Microsoft Internet Explorer Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-0222,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 5010 | Microsoft Office Use-After-Free Vulnerability | kev | markdown |  | CVE-2017-0261,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 5011 | Microsoft Office Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-0262,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 5012 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2017-0263,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 5013 | Linux Kernel PIE Stack Buffer Corruption Vulnerability | kev | markdown |  | CVE-2017-1000253,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Linux |
| 5014 | Oracle Corporation WebLogic Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-10271,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Oracle, |
| 5015 | Adobe Flash Player Type Confusion Vulnerability | kev | markdown |  | CVE-2017-11292,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 5016 | Microsoft Office Outlook Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2017-11774,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 5017 | Microsoft Office Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-11826,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 5018 | Microsoft Office Memory Corruption Vulnerability | kev | markdown |  | CVE-2017-11882,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5019 | Cisco IOS Software Network Address Translation Denial-of-Service Vulnerability | kev | markdown |  | CVE-2017-12231,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS software |
| 5020 | Cisco IOS Software for Cisco Integrated Services Routers Denial-of-Service Vulnerability | kev | markdown |  | CVE-2017-12232,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS software |
| 5021 | Cisco IOS Software Common Industrial Protocol Request Denial-of-Service Vulnerability | kev | markdown |  | CVE-2017-12233,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS software |
| 5022 | Cisco IOS Software Common Industrial Protocol Request Denial-of-Service Vulnerability | kev | markdown |  | CVE-2017-12234,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS software |
| 5023 | Cisco IOS Software for Cisco Industrial Ethernet Switches PROFINET Denial-of-Service Vulnerability | kev | markdown |  | CVE-2017-12235,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS software |
| 5024 | Cisco IOS and IOS XE Software Internet Key Exchange Denial-of-Service Vulnerability | kev | markdown |  | CVE-2017-12237,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS  |
| 5025 | Cisco Catalyst 6800 Series Switches VPLS Denial-of-Service Vulnerability | kev | markdown |  | CVE-2017-12238,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Catalyst 680 |
| 5026 | Cisco IOS XE Software Ethernet Virtual Private Network Border Gateway Protocol Denial-of-Service Vulnerability | kev | markdown |  | CVE-2017-12319,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS XE Softw |
| 5027 | Apache Tomcat on Windows Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-12615,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Apache, |
| 5028 | Apache Tomcat Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-12617,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Tomcat |
| 5029 | SAP NetWeaver Directory Traversal Vulnerability | kev | markdown |  | CVE-2017-12637,cve-arsenal,cisa-kev,kev,in-the-wild,SAP,NetWeaver |
| 5030 | Roundcube Webmail File Disclosure Vulnerability | kev | markdown |  | CVE-2017-16651,cve-arsenal,cisa-kev,kev,in-the-wild,Roundcube,Roundcub |
| 5031 | Embedthis GoAhead Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-17562,cve-arsenal,cisa-kev,kev,in-the-wild,Embedthis,GoAhead |
| 5032 | Zyxel P660HN-T1A Routers Command Injection Vulnerability | kev | markdown |  | CVE-2017-18368,cve-arsenal,cisa-kev,kev,in-the-wild,Zyxel,P660HN-T1A R |
| 5033 | Oracle WebLogic Server OS Command Injection Vulnerability | kev | markdown |  | CVE-2017-3506,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,WebLogic Ser |
| 5034 | Google Chromium V8 Memory Corruption Vulnerability | kev | markdown |  | CVE-2017-5030,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5035 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2017-5070,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5036 | NETGEAR Multiple Devices Exposure of Sensitive Information Vulnerability | kev | markdown |  | CVE-2017-5521,cve-arsenal,cisa-kev,kev,in-the-wild,NETGEAR,Multiple De |
| 5037 | Symantec Messaging Gateway Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-6327,cve-arsenal,cisa-kev,kev,in-the-wild,Symantec,Symantec M |
| 5038 | NETGEAR DGN2200 Devices OS Command Injection Vulnerability | kev | markdown |  | CVE-2017-6334,cve-arsenal,cisa-kev,kev,in-the-wild,NETGEAR,DGN2200 Dev |
| 5039 | Cisco IOS Software and Cisco IOS XE Software UDP Packet Processing Denial-of-Service Vulnerability | kev | markdown |  | CVE-2017-6627,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS X |
| 5040 | Cisco IOS Software and Cisco IOS XE Software Denial-of-Service Vulnerability | kev | markdown |  | CVE-2017-6663,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS X |
| 5041 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-6736,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS X |
| 5042 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-6737,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS X |
| 5043 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-6738,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS X |
| 5044 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-6739,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS X |
| 5045 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-6740,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS X |
| 5046 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-6742,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS X |
| 5047 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-6743,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS X |
| 5048 | Cisco IOS Software SNMP Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-6744,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS software |
| 5049 | Zyxel EMG2926 Routers Command Injection Vulnerability | kev | markdown |  | CVE-2017-6884,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Zyxel,EM |
| 5050 | Artifex Ghostscript Type Confusion Vulnerability | kev | markdown |  | CVE-2017-8291,cve-arsenal,cisa-kev,kev,in-the-wild,Artifex,Ghostscript |
| 5051 | Microsoft Windows Shell (.lnk) Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-8464,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5052 | Microsoft Malware Protection Engine Improper Restriction of Operations Vulnerability | kev | markdown |  | CVE-2017-8540,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Malware P |
| 5053 | Microsoft Office Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-8570,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 5054 | Microsoft .NET Framework Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-8759,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,.NET Fram |
| 5055 | Apache Struts Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2017-9805,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Struts |
| 5056 | DotNetNuke (DNN) Remote Code Execution Vulnerability | kev | markdown |  | CVE-2017-9822,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,DotNetNu |
| 5057 | Cisco IOS Software Integrated Services Module for VPN Denial-of-Service Vulnerability | kev | markdown |  | CVE-2018-0154,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS Software |
| 5058 | Cisco Catalyst Bidirectional Forwarding Detection Denial-of-Service Vulnerability | kev | markdown |  | CVE-2018-0155,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Catalyst 4500 |
| 5059 | Cisco IOS Software and Cisco IOS XE Software Smart Install Denial-of-Service Vulnerability | kev | markdown |  | CVE-2018-0156,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS Software  |
| 5060 | Cisco IOS and XE Software Internet Key Exchange Memory Leak Vulnerability | kev | markdown |  | CVE-2018-0158,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS Software  |
| 5061 | Cisco IOS and XE Software Internet Key Exchange Version 1 Denial-of-Service Vulnerability | kev | markdown |  | CVE-2018-0159,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS Software  |
| 5062 | Cisco IOS Software Resource Management Errors Vulnerability | kev | markdown |  | CVE-2018-0161,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS Software |
| 5063 | Cisco IOS, XR, and XE Software Buffer Overflow Vulnerability | kev | markdown |  | CVE-2018-0167,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS, XR, and  |
| 5064 | Cisco IOS and IOS XE Software Improper Input Validation Vulnerability | kev | markdown |  | CVE-2018-0172,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS X |
| 5065 | Cisco IOS and IOS XE Software Improper Input Validation Vulnerability | kev | markdown |  | CVE-2018-0173,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS X |
| 5066 | Cisco IOS Software and Cisco IOS XE Software Improper Input Validation Vulnerability | kev | markdown |  | CVE-2018-0174,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS XE Softwa |
| 5067 | Cisco IOS, XR, and XE Software Buffer Overflow Vulnerability | kev | markdown |  | CVE-2018-0175,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS, XR, and  |
| 5068 | Cisco IOS Software Denial-of-Service Vulnerability | kev | markdown |  | CVE-2018-0179,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS Software |
| 5069 | Cisco IOS Software Denial-of-Service Vulnerability | kev | markdown |  | CVE-2018-0180,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS Software |
| 5070 | Cisco Adaptive Security Appliance (ASA) Denial-of-Service Vulnerability | kev | markdown |  | CVE-2018-0296,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Adaptive Secu |
| 5071 | Microsoft Office Memory Corruption Vulnerability | kev | markdown |  | CVE-2018-0798,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 5072 | Microsoft Office Memory Corruption Vulnerability | kev | markdown |  | CVE-2018-0802,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5073 | Microsoft COM for Windows Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2018-0824,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5074 | Apache Struts Remote Code Execution Vulnerability | kev | markdown |  | CVE-2018-11776,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Struts |
| 5075 | Fortinet FortiOS and FortiADC Improper Access Control Vulnerability | kev | markdown |  | CVE-2018-13374,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortine |
| 5076 | Fortinet FortiOS SSL VPN Path Traversal Vulnerability | kev | markdown |  | CVE-2018-13379,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortine |
| 5077 | Fortinet FortiOS and FortiProxy Improper Authorization | kev | markdown |  | CVE-2018-13382,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortine |
| 5078 | Fortinet FortiOS and FortiProxy Out-of-bounds Write | kev | markdown |  | CVE-2018-13383,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortine |
| 5079 | Linux Kernel Integer Overflow Vulnerability | kev | markdown |  | CVE-2018-14634,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 5080 | LG N1A1 NAS Remote Command Execution Vulnerability | kev | markdown |  | CVE-2018-14839,cve-arsenal,cisa-kev,kev,in-the-wild,LG,N1A1 NAS |
| 5081 | Laravel Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2018-15133,cve-arsenal,cisa-kev,kev,in-the-wild,Laravel,Laravel Fr |
| 5082 | DotNetNuke (DNN) Inadequate Encryption Strength Vulnerability | kev | markdown |  | CVE-2018-15811,cve-arsenal,cisa-kev,kev,in-the-wild,DotNetNuke (DNN) |
| 5083 | Adobe Flash Player Use-After-Free Vulnerability | kev | markdown |  | CVE-2018-15982,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Adobe,F |
| 5084 | Google Chromium V8 Remote Code Execution Vulnerability | kev | markdown |  | CVE-2018-17463,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5085 | Google Chromium V8 Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2018-17480,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5086 | DotNetNuke (DNN) Inadequate Encryption Strength Vulnerability | kev | markdown |  | CVE-2018-18325,cve-arsenal,cisa-kev,kev,in-the-wild,DotNetNuke (DNN) |
| 5087 | TIBCO JasperReports Library Directory Traversal Vulnerability | kev | markdown |  | CVE-2018-18809,cve-arsenal,cisa-kev,kev,in-the-wild,TIBCO,JasperReport |
| 5088 | GIGABYTE Multiple Products Unspecified Vulnerability | kev | markdown |  | CVE-2018-19320,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,GIGABYT |
| 5089 | GIGABYTE Multiple Products Privilege Escalation Vulnerability | kev | markdown |  | CVE-2018-19321,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,GIGABYT |
| 5090 | GIGABYTE Multiple Products Code Execution Vulnerability | kev | markdown |  | CVE-2018-19322,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,GIGABYT |
| 5091 | QNAP NAS File Station Cross-Site Scripting Vulnerability | kev | markdown |  | CVE-2018-19943,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,QNAP,Ne |
| 5092 | QNAP NAS File Station Command Injection Vulnerability | kev | markdown |  | CVE-2018-19949,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,QNAP,Ne |
| 5093 | QNAP NAS File Station Cross-Site Scripting Vulnerability | kev | markdown |  | CVE-2018-19953,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,QNAP,Ne |
| 5094 | WinRAR Absolute Path Traversal Vulnerability | kev | markdown |  | CVE-2018-20250,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,RARLAB, |
| 5095 | SAP Customer Relationship Management (CRM) Path Traversal Vulnerability | kev | markdown |  | CVE-2018-2380,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SAP,Cust |
| 5096 | Sierra Wireless AirLink ALEOS Unrestricted Upload of File with Dangerous Type Vulnerability | kev | markdown |  | CVE-2018-4063,cve-arsenal,cisa-kev,kev,in-the-wild,Sierra Wireless,Air |
| 5097 | Apple Multiple Products Memory Corruption Vulnerability | kev | markdown |  | CVE-2018-4344,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Prod |
| 5098 | Adobe Flash Player Use-After-Free Vulnerability | kev | markdown |  | CVE-2018-4878,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Adobe,Fl |
| 5099 | Adobe Acrobat and Reader Double Free Vulnerability | kev | markdown |  | CVE-2018-4990,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat and R |
| 5100 | Adobe Flash Player Stack-based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2018-5002,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Flash Player |
| 5101 | TIBCO JasperReports Server Information Disclosure Vulnerability | kev | markdown |  | CVE-2018-5430,cve-arsenal,cisa-kev,kev,in-the-wild,TIBCO,JasperReports |
| 5102 | Google Chromium V8 Integer Overflow Vulnerability | kev | markdown |  | CVE-2018-6065,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5103 | Synacor Zimbra Collaboration Suite (ZCS) Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2018-6882,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Synacor, |
| 5104 | VMware SD-WAN Edge by VeloCloud Command Injection Vulnerability | kev | markdown |  | CVE-2018-6961,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,SD-WAN Edge |
| 5105 | Schneider Electric U.motion Builder SQL Injection Vulnerability | kev | markdown |  | CVE-2018-7841,cve-arsenal,cisa-kev,kev,in-the-wild,Schneider Electric, |
| 5106 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2018-8120,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5107 | Microsoft Windows VBScript Engine Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2018-8174,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5108 | ChakraCore Scripting Engine Type Confusion Vulnerability | kev | markdown |  | CVE-2018-8298,cve-arsenal,cisa-kev,kev,in-the-wild,ChakraCore,ChakraCo |
| 5109 | Microsoft Scripting Engine Memory Corruption Vulnerability | kev | markdown |  | CVE-2018-8373,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 5110 | Microsoft DirectX Graphics Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2018-8405,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5111 | Microsoft DirectX Graphics Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2018-8406,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5112 | Microsoft Windows Shell Remote Code Execution Vulnerability | kev | markdown |  | CVE-2018-8414,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5113 | Microsoft Windows Privilege Escalation Vulnerability | kev | markdown |  | CVE-2018-8440,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5114 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2018-8453,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5115 | Microsoft Exchange Server Privilege Escalation Vulnerability | kev | markdown |  | CVE-2018-8581,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5116 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2018-8589,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 5117 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2018-8611,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5118 | Microsoft Windows Win32k Improper Resource Shutdown or Release Vulnerability | kev | markdown |  | CVE-2018-8639,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5119 | Microsoft Internet Explorer Scripting Engine Memory Corruption Vulnerability | kev | markdown |  | CVE-2018-8653,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 5120 | Paessler PRTG Network Monitor OS Command Injection Vulnerability | kev | markdown |  | CVE-2018-9276,cve-arsenal,cisa-kev,kev,in-the-wild,Paessler,PRTG Netwo |
| 5121 | Apache Solr DataImportHandler Code Injection Vulnerability | kev | markdown |  | CVE-2019-0193,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Solr |
| 5122 | Apache HTTP Server Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-0211,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,HTTP Server |
| 5123 | SAP Commerce Cloud Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2019-0344,cve-arsenal,cisa-kev,kev,in-the-wild,SAP,Commerce Cloud |
| 5124 | Microsoft MSHTML Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-0541,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,MSHTML |
| 5125 | Microsoft Windows Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-0543,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5126 | Microsoft SharePoint Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-0604,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5127 | Microsoft Internet Explorer Information Disclosure Vulnerability | kev | markdown |  | CVE-2019-0676,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 5128 | Microsoft Windows SMB Information Disclosure Vulnerability | kev | markdown |  | CVE-2019-0703,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5129 | Microsoft Remote Desktop Services Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-0708,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5130 | Microsoft Internet Explorer Type Confusion Vulnerability | kev | markdown |  | CVE-2019-0752,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5131 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-0797,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 5132 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-0803,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5133 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-0808,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 5134 | Microsoft Windows AppX Deployment Service (AppXSVC) Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-0841,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5135 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-0859,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5136 | Microsoft Windows Error Reporting (WER) Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-0863,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5137 | Microsoft Windows Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-0880,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5138 | Microsoft GDI Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-0903,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Graphics  |
| 5139 | Jenkins Script Security Plugin Sandbox Bypass Vulnerability | kev | markdown |  | CVE-2019-1003029,cve-arsenal,cisa-kev,kev,in-the-wild,Jenkins,Script S |
| 5140 | Jenkins Matrix Project Plugin Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-1003030,cve-arsenal,cisa-kev,kev,in-the-wild,Jenkins,Matrix P |
| 5141 | Kentico Xperience Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2019-10068,cve-arsenal,cisa-kev,kev,in-the-wild,Kentico,Xperience |
| 5142 | Exim Mail Transfer Agent (MTA) Improper Input Validation | kev | markdown |  | CVE-2019-10149,cve-arsenal,cisa-kev,kev,in-the-wild,Exim,Mail Transfer |
| 5143 | Microsoft Windows AppX Deployment Service (AppXSVC) Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-1064,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5144 | Microsoft SQL Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-1068,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,SQL Serve |
| 5145 | Microsoft Task Scheduler Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-1069,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5146 | MongoDB mongo-express Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-10758,cve-arsenal,cisa-kev,kev,in-the-wild,MongoDB,mongo-expr |
| 5147 | Reolink Multiple IP Cameras OS Command Injection Vulnerability | kev | markdown |  | CVE-2019-11001,cve-arsenal,cisa-kev,kev,in-the-wild,Reolink,Multiple I |
| 5148 | PHP FastCGI Process Manager (FPM) Buffer Overflow Vulnerability | kev | markdown |  | CVE-2019-11043,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,PHP,Fas |
| 5149 | Microsoft Windows AppX Deployment Service (AppXSVC) Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-1129,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5150 | Microsoft Windows AppX Deployment Service Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-1130,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5151 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-1132,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 5152 | Ivanti Pulse Connect Secure Arbitrary File Read Vulnerability | kev | markdown |  | CVE-2019-11510,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Ivanti, |
| 5153 | Ivanti Pulse Connect Secure and Policy Secure Command Injection Vulnerability | kev | markdown |  | CVE-2019-11539,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Ivanti, |
| 5154 | Atlassian Crowd and Crowd Data Center Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-11580,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Atlassi |
| 5155 | Atlassian Jira Server and Data Center Server-Side Template Injection Vulnerability | kev | markdown |  | CVE-2019-11581,cve-arsenal,cisa-kev,kev,in-the-wild,Atlassian,Jira Ser |
| 5156 | Citrix Workspace Application and Receiver for Windows Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-11634,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Citrix, |
| 5157 | Mozilla Firefox and Thunderbird Type Confusion Vulnerability | kev | markdown |  | CVE-2019-11707,cve-arsenal,cisa-kev,kev,in-the-wild,Mozilla,Firefox an |
| 5158 | Mozilla Firefox and Thunderbird Sandbox Escape Vulnerability | kev | markdown |  | CVE-2019-11708,cve-arsenal,cisa-kev,kev,in-the-wild,Mozilla,Firefox an |
| 5159 | Microsoft Windows Privilege Common Log File System (CLFS) Escalation Vulnerability | kev | markdown |  | CVE-2019-1214,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5160 | Microsoft Windows Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-1215,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5161 | Microsoft Windows AppX Deployment Server Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-1253,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5162 | Microsoft Excel Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-1297,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Excel |
| 5163 | Citrix SD-WAN and NetScaler SQL Injection Vulnerability | kev | markdown |  | CVE-2019-12989,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,SD-WAN and  |
| 5164 | Citrix SD-WAN and NetScaler Command Injection Vulnerability | kev | markdown |  | CVE-2019-12991,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,SD-WAN and  |
| 5165 | Microsoft Windows Error Reporting Manager Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-1315,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5166 | Microsoft Windows Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-1322,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5167 | Linux Kernel Improper Privilege Management Vulnerability | kev | markdown |  | CVE-2019-13272,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 5168 | Citrix StoreFront Server XML External Entity (XXE) Processing Vulnerability | kev | markdown |  | CVE-2019-13608,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Citrix, |
| 5169 | Microsoft Internet Explorer Scripting Engine Memory Corruption Vulnerability | kev | markdown |  | CVE-2019-1367,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5170 | Google Chrome WebAudio Use-After-Free Vulnerability | kev | markdown |  | CVE-2019-13720,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chrome WebA |
| 5171 | Microsoft Windows AppX Deployment Extensions Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-1385,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5172 | Microsoft Windows Certificate Dialog Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-1388,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5173 | Microsoft Windows Universal Plug and Play (UPnP) Service Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-1405,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5174 | Microsoft Internet Explorer Scripting Engine Memory Corruption Vulnerability | kev | markdown |  | CVE-2019-1429,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 5175 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-1458,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5176 | Webmin Command Injection Vulnerability | kev | markdown |  | CVE-2019-15107,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Webmin |
| 5177 | Cisco RV Series Routers Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2019-15271,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,RV Series Ro |
| 5178 | Docker Desktop Community Edition Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-15752,cve-arsenal,cisa-kev,kev,in-the-wild,Docker,Desktop Com |
| 5179 | Palo Alto Networks PAN-OS Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-1579,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Palo Alt |
| 5180 | Nagios XI Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-15949,cve-arsenal,cisa-kev,kev,in-the-wild,Nagios,Nagios XI |
| 5181 | D-Link DNS-320 Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-16057,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,D-Link, |
| 5182 | SIMalliance Toolbox Browser Command Injection Vulnerability | kev | markdown |  | CVE-2019-16256,cve-arsenal,cisa-kev,kev,in-the-wild,SIMalliance,Toolbo |
| 5183 | Nostromo nhttpd Directory Traversal Vulnerability | kev | markdown |  | CVE-2019-16278,cve-arsenal,cisa-kev,kev,in-the-wild,Nostromo,nhttpd |
| 5184 | Cisco Small Business Routers Improper Input Validation Vulnerability | kev | markdown |  | CVE-2019-1652,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Small Busines |
| 5185 | Cisco Small Business RV320 and RV325 Routers Information Disclosure Vulnerability | kev | markdown |  | CVE-2019-1653,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Small Busines |
| 5186 | vBulletin PHP Module Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-16759,cve-arsenal,cisa-kev,kev,in-the-wild,vBulletin |
| 5187 | D-Link Multiple Routers Command Injection Vulnerability | kev | markdown |  | CVE-2019-16920,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,Multiple Ro |
| 5188 | Exim Out-of-bounds Write Vulnerability | kev | markdown |  | CVE-2019-16928,cve-arsenal,cisa-kev,kev,in-the-wild,Exim,Exim Internet |
| 5189 | Mozilla Firefox And Thunderbird Type Confusion Vulnerability | kev | markdown |  | CVE-2019-17026,cve-arsenal,cisa-kev,kev,in-the-wild,Mozilla,Firefox an |
| 5190 | Apache Solr VelocityResponseWriter Plug-In Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-17558,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Solr |
| 5191 | D-Link DIR-859 Router Command Execution Vulnerability | kev | markdown |  | CVE-2019-17621,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DIR-859 Rou |
| 5192 | Trend Micro OfficeScan Directory Traversal Vulnerability | kev | markdown |  | CVE-2019-18187,cve-arsenal,cisa-kev,kev,in-the-wild,Trend Micro,Office |
| 5193 | WhatsApp Cross-Site Scripting Vulnerability | kev | markdown |  | CVE-2019-18426,cve-arsenal,cisa-kev,kev,in-the-wild,Meta Platforms,Wha |
| 5194 | Progress Telerik UI for ASP.NET AJAX Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2019-18935,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Progres |
| 5195 | TeamViewer Desktop Bypass Remote Login Vulnerability | kev | markdown |  | CVE-2019-18988,cve-arsenal,cisa-kev,kev,in-the-wild,TeamViewer,Desktop |
| 5196 | Sangoma FreePBX Improper Authentication Vulnerability | kev | markdown |  | CVE-2019-19006,cve-arsenal,cisa-kev,kev,in-the-wild,Sangoma,FreePBX |
| 5197 | Netis WF2419 Devices Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-19356,cve-arsenal,cisa-kev,kev,in-the-wild,Netis,WF2419 Devic |
| 5198 | Citrix ADC, Gateway, and SD-WAN WANOP Appliance Code Execution Vulnerability | kev | markdown |  | CVE-2019-19781,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Citrix, |
| 5199 | TVT NVMS-1000 Directory Traversal Vulnerability | kev | markdown |  | CVE-2019-20085,cve-arsenal,cisa-kev,kev,in-the-wild,TVT,NVMS-1000 |
| 5200 | D-Link DWL-2600AP Access Point Command Injection Vulnerability | kev | markdown |  | CVE-2019-20500,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DWL-2600AP  |
| 5201 | Android Kernel Use-After-Free Vulnerability | kev | markdown |  | CVE-2019-2215,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Android Ker |
| 5202 | Oracle BI Publisher Unauthorized Access Vulnerability | kev | markdown |  | CVE-2019-2616,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,BI Publisher |
| 5203 | Oracle WebLogic Server, Injection | kev | markdown |  | CVE-2019-2725,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Oracle,W |
| 5204 | Oracle Solaris Privilege Escalation Vulnerability | kev | markdown |  | CVE-2019-3010,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,Solaris |
| 5205 | Atlassian Confluence Server and Data Center Server-Side Template Injection Vulnerability | kev | markdown |  | CVE-2019-3396,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Atlassia |
| 5206 | Atlassian Confluence Server and Data Center Path Traversal Vulnerability | kev | markdown |  | CVE-2019-3398,cve-arsenal,cisa-kev,kev,in-the-wild,Atlassian,Confluenc |
| 5207 | WhatsApp VOIP Stack Buffer Overflow Vulnerability | kev | markdown |  | CVE-2019-3568,cve-arsenal,cisa-kev,kev,in-the-wild,Meta Platforms,What |
| 5208 | Crestron Multiple Products Command Injection Vulnerability | kev | markdown |  | CVE-2019-3929,cve-arsenal,cisa-kev,kev,in-the-wild,Crestron,Multiple P |
| 5209 | IBM Planning Analytics Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-4716,cve-arsenal,cisa-kev,kev,in-the-wild,IBM,Planning Analyt |
| 5210 | Rails Ruby on Rails Path Traversal Vulnerability | kev | markdown |  | CVE-2019-5418,cve-arsenal,cisa-kev,kev,in-the-wild,Rails,Ruby on Rails |
| 5211 | VMware ESXi and Horizon DaaS OpenSLP Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2019-5544,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,VMware,V |
| 5212 | Fortinet FortiOS Default Configuration Vulnerability | kev | markdown |  | CVE-2019-5591,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortinet |
| 5213 | Google Chrome Blink Use-After-Free Vulnerability | kev | markdown |  | CVE-2019-5786,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chrome Blink |
| 5214 | Google Chromium V8 Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2019-5825,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5215 | Apple iOS and macOS Group Facetime Vulnerability | kev | markdown |  | CVE-2019-6223,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS and macOS |
| 5216 | Drupal Core Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-6340,cve-arsenal,cisa-kev,kev,in-the-wild,Drupal,Core |
| 5217 | Fortinet FortiOS Use of Hard-Coded Credentials Vulnerability | kev | markdown |  | CVE-2019-6693,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortinet |
| 5218 | QNAP Photo Station Improper Access Control Vulnerability | kev | markdown |  | CVE-2019-7192,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,QNAP,Pho |
| 5219 | QNAP QTS Improper Input Validation Vulnerability | kev | markdown |  | CVE-2019-7193,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,QNAP,QTS |
| 5220 | QNAP Photo Station Path Traversal Vulnerability | kev | markdown |  | CVE-2019-7194,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,QNAP,Pho |
| 5221 | QNAP Photo Station Path Traversal Vulnerability | kev | markdown |  | CVE-2019-7195,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,QNAP,Pho |
| 5222 | Sonatype Nexus Repository Manager Incorrect Access Control Vulnerability | kev | markdown |  | CVE-2019-7238,cve-arsenal,cisa-kev,kev,in-the-wild,Sonatype,Nexus Repo |
| 5223 | Nice Linear eMerge E3-Series OS Command Injection Vulnerability | kev | markdown |  | CVE-2019-7256,cve-arsenal,cisa-kev,kev,in-the-wild,Nice,Linear eMerge  |
| 5224 | Apple Multiple Products Memory Corruption Vulnerability | kev | markdown |  | CVE-2019-7286,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Prod |
| 5225 | Apple iOS Memory Corruption Vulnerability | kev | markdown |  | CVE-2019-7287,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS |
| 5226 | SonicWall SMA100 SQL Injection Vulnerability | kev | markdown |  | CVE-2019-7481,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SonicWal |
| 5227 | SonicWall SMA100 Directory Traversal Vulnerability | kev | markdown |  | CVE-2019-7483,cve-arsenal,cisa-kev,kev,in-the-wild,SonicWall,SMA100 |
| 5228 | Kibana Arbitrary Code Execution | kev | markdown |  | CVE-2019-7609,cve-arsenal,cisa-kev,kev,in-the-wild,Elastic,Kibana |
| 5229 | Zoho ManageEngine ServiceDesk Plus (SDP) File Upload Vulnerability | kev | markdown |  | CVE-2019-8394,cve-arsenal,cisa-kev,kev,in-the-wild,Zoho,ManageEngine |
| 5230 | Apple Multiple Products Type Confusion Vulnerability | kev | markdown |  | CVE-2019-8506,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Prod |
| 5231 | Apple macOS Use-After-Free Vulnerability | kev | markdown |  | CVE-2019-8526,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,macOS |
| 5232 | Apple Multiple Products Use-After-Free Vulnerability | kev | markdown |  | CVE-2019-8605,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Prod |
| 5233 | WebKitGTK Memory Corruption Vulnerability | kev | markdown |  | CVE-2019-8720,cve-arsenal,cisa-kev,kev,in-the-wild,WebKitGTK |
| 5234 | ThinkPHP Remote Code Execution Vulnerability | kev | markdown |  | CVE-2019-9082,cve-arsenal,cisa-kev,kev,in-the-wild,ThinkPHP |
| 5235 | Synacor Zimbra Collaboration Suite (ZCS) Server-Side Request Forgery (SSRF) Vulnerability | kev | markdown |  | CVE-2019-9621,cve-arsenal,cisa-kev,kev,in-the-wild,Synacor,Zimbra Coll |
| 5236 | Synacor Zimbra Collaboration Suite (ZCS) Improper Restriction of XML External Entity Reference | kev | markdown |  | CVE-2019-9670,cve-arsenal,cisa-kev,kev,in-the-wild,Synacor,Zimbra Coll |
| 5237 | Sitecore CMS and Experience Platform (XP) Deserialization Vulnerability | kev | markdown |  | CVE-2019-9874,cve-arsenal,cisa-kev,kev,in-the-wild,Sitecore,CMS and Ex |
| 5238 | Sitecore CMS and Experience Platform (XP) Deserialization Vulnerability | kev | markdown |  | CVE-2019-9875,cve-arsenal,cisa-kev,kev,in-the-wild,Sitecore,CMS and Ex |
| 5239 | WordPress Social Warfare Plugin Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2019-9978,cve-arsenal,cisa-kev,kev,in-the-wild,WordPress,Social Wa |
| 5240 | Android Kernel Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2020-0041,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Android Ker |
| 5241 | Mediatek Multiple Chipsets Insufficient Input Validation Vulnerability | kev | markdown |  | CVE-2020-0069,cve-arsenal,cisa-kev,kev,in-the-wild,MediaTek,Multiple C |
| 5242 | Microsoft Windows CryptoAPI Spoofing Vulnerability | kev | markdown |  | CVE-2020-0601,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5243 | Microsoft SQL Server Reporting Services Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-0618,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5244 | Microsoft Update Notification Manager Privilege Escalation Vulnerability | kev | markdown |  | CVE-2020-0638,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5245 | Microsoft .NET Framework Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-0646,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,.NET Fram |
| 5246 | Microsoft Internet Explorer Scripting Engine Memory Corruption Vulnerability | kev | markdown |  | CVE-2020-0674,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 5247 | Microsoft Windows Installer Privilege Escalation Vulnerability | kev | markdown |  | CVE-2020-0683,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5248 | Microsoft Exchange Server Validation Key Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-0688,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5249 | Microsoft Windows Background Intelligent Transfer Service (BITS) Improper Privilege Management Vulnerability | kev | markdown |  | CVE-2020-0787,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5250 | Microsoft SMBv3 Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-0796,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5251 | Microsoft Edge and Internet Explorer Memory Corruption Vulnerability | kev | markdown |  | CVE-2020-0878,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5252 | Microsoft Windows Adobe Font Manager Library Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-0938,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5253 | Microsoft Internet Explorer Scripting Engine Memory Corruption Vulnerability | kev | markdown |  | CVE-2020-0968,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5254 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2020-0986,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5255 | SolarWinds Orion Authentication Bypass Vulnerability | kev | markdown |  | CVE-2020-10148,cve-arsenal,cisa-kev,kev,in-the-wild,SolarWinds,Orion |
| 5256 | Sumavision EMR Cross-Site Request Forgery (CSRF) Vulnerability | kev | markdown |  | CVE-2020-10181,cve-arsenal,cisa-kev,kev,in-the-wild,Sumavision,Enhance |
| 5257 | Zoho ManageEngine Desktop Central File Upload Vulnerability | kev | markdown |  | CVE-2020-10189,cve-arsenal,cisa-kev,kev,in-the-wild,Zoho,ManageEngine |
| 5258 | Sonatype Nexus Repository Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-10199,cve-arsenal,cisa-kev,kev,in-the-wild,Sonatype,Nexus Rep |
| 5259 | Microsoft Windows Adobe Font Manager Library Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-1020,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5260 | rConfig OS Command Injection Vulnerability | kev | markdown |  | CVE-2020-10221,cve-arsenal,cisa-kev,kev,in-the-wild,rConfig |
| 5261 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2020-1027,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5262 | Microsoft Hyper-V RemoteFX vGPU Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-1040,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Hyper-V R |
| 5263 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2020-1054,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 5264 | Tenda AC1900 Router AC15 Model Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-10987,cve-arsenal,cisa-kev,kev,in-the-wild,Tenda,AC1900 Route |
| 5265 | JQuery Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2020-11023,cve-arsenal,cisa-kev,kev,in-the-wild,JQuery |
| 5266 | Qualcomm Multiple Chipsets Improper Input Validation Vulnerability | kev | markdown |  | CVE-2020-11261,cve-arsenal,cisa-kev,kev,in-the-wild,Qualcomm,Snapdrago |
| 5267 | Microsoft .NET Framework, SharePoint, and Visual Studio Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-1147,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,.NET Fram |
| 5268 | SaltStack Salt Authentication Bypass Vulnerability | kev | markdown |  | CVE-2020-11651,cve-arsenal,cisa-kev,kev,in-the-wild,SaltStack,Salt |
| 5269 | SaltStack Salt Path Traversal Vulnerability | kev | markdown |  | CVE-2020-11652,cve-arsenal,cisa-kev,kev,in-the-wild,SaltStack,Salt |
| 5270 | WordPress Snap Creek Duplicator Plugin File Download Vulnerability | kev | markdown |  | CVE-2020-11738,cve-arsenal,cisa-kev,kev,in-the-wild,WordPress,Snap Cre |
| 5271 | Treck TCP/IP stack Out-of-Bounds Read Vulnerability | kev | markdown |  | CVE-2020-11899,cve-arsenal,cisa-kev,kev,in-the-wild,Treck TCP/IP stack |
| 5272 | Apache Airflow Command Injection | kev | markdown |  | CVE-2020-11978,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Airflow |
| 5273 | Sophos SFOS SQL Injection Vulnerability | kev | markdown |  | CVE-2020-12271,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Sophos, |
| 5274 | Roundcube Webmail Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-12641,cve-arsenal,cisa-kev,kev,in-the-wild,Roundcube,Roundcub |
| 5275 | Fortinet FortiOS SSL VPN Improper Authentication Vulnerability | kev | markdown |  | CVE-2020-12812,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortine |
| 5276 | Microsoft Windows DNS Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-1350,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5277 | Drupal core Un-restricted Upload of File | kev | markdown |  | CVE-2020-13671,cve-arsenal,cisa-kev,kev,in-the-wild,Drupal,Drupal core |
| 5278 | Microsoft Internet Explorer Scripting Engine Memory Corruption Vulnerability | kev | markdown |  | CVE-2020-1380,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet  |
| 5279 | Apache Airflow's Experimental API Authentication Bypass | kev | markdown |  | CVE-2020-13927,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Airflow's E |
| 5280 | Roundcube Webmail Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2020-13965,cve-arsenal,cisa-kev,kev,in-the-wild,Roundcube,Webmail |
| 5281 | Microsoft Windows Spoofing Vulnerability | kev | markdown |  | CVE-2020-1464,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5282 | Oracle WebLogic Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-14644,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,WebLogic Se |
| 5283 | Microsoft Netlogon Privilege Escalation Vulnerability | kev | markdown |  | CVE-2020-1472,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5284 | Oracle WebLogic Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-14750,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,WebLogic Se |
| 5285 | Oracle Business Intelligence Enterprise Edition Path Transversal | kev | markdown |  | CVE-2020-14864,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,Intelligenc |
| 5286 | Oracle Solaris and Zettabyte File System (ZFS) Unspecified Vulnerability | kev | markdown |  | CVE-2020-14871,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,Solaris and |
| 5287 | Oracle WebLogic Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-14882,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,WebLogic Se |
| 5288 | Oracle WebLogic Server Unspecified Vulnerability | kev | markdown |  | CVE-2020-14883,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,WebLogic Se |
| 5289 | Sophos XG Firewall Buffer Overflow Vulnerability | kev | markdown |  | CVE-2020-15069,cve-arsenal,cisa-kev,kev,in-the-wild,Sophos,XG Firewall |
| 5290 | DrayTek Multiple Vigor Routers OS Command Injection Vulnerability | kev | markdown |  | CVE-2020-15415,cve-arsenal,cisa-kev,kev,in-the-wild,DrayTek,Multiple V |
| 5291 | Ivanti MobileIron Multiple Products Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-15505,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,MobileIron  |
| 5292 | Google Chrome FreeType Heap Buffer Overflow Vulnerability | kev | markdown |  | CVE-2020-15999,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chrome Free |
| 5293 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2020-16009,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5294 | Google Chrome for Android UI Heap Buffer Overflow Vulnerability | kev | markdown |  | CVE-2020-16010,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chrome for  |
| 5295 | Google Chromium V8 Incorrect Implementation Vulnerabililty | kev | markdown |  | CVE-2020-16013,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5296 | Google Chrome Use-After-Free Vulnerability | kev | markdown |  | CVE-2020-16017,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chrome |
| 5297 | Juniper Junos OS Path Traversal Vulnerability | kev | markdown |  | CVE-2020-1631,cve-arsenal,cisa-kev,kev,in-the-wild,Juniper,Junos OS |
| 5298 | SaltStack Salt Shell Injection Vulnerability | kev | markdown |  | CVE-2020-16846,cve-arsenal,cisa-kev,kev,in-the-wild,SaltStack,Salt |
| 5299 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2020-17087,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5300 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-17144,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Exchange |
| 5301 | Fuel CMS SQL Injection Vulnerability | kev | markdown |  | CVE-2020-17463,cve-arsenal,cisa-kev,kev,in-the-wild,Fuel CMS |
| 5302 | vBulletin PHP Module Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-17496,cve-arsenal,cisa-kev,kev,in-the-wild,vBulletin |
| 5303 | Apache Flink Improper Access Control Vulnerability | kev | markdown |  | CVE-2020-17519,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Flink |
| 5304 | Apache Struts Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-17530,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Struts |
| 5305 | Apache Tomcat Improper Privilege Management Vulnerability | kev | markdown |  | CVE-2020-1938,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Tomcat |
| 5306 | Apache Kylin OS Command Injection Vulnerability | kev | markdown |  | CVE-2020-1956,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Kylin |
| 5307 | Palo Alto Networks PAN-OS Authentication Bypass Vulnerability | kev | markdown |  | CVE-2020-2021,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Palo Alt |
| 5308 | TP-link TL-WA855RE Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2020-24363,cve-arsenal,cisa-kev,kev,in-the-wild,TP-Link,TL-WA855RE |
| 5309 | Trend Micro Multiple Products Improper Access Control Vulnerability | kev | markdown |  | CVE-2020-24557,cve-arsenal,cisa-kev,kev,in-the-wild,Trend Micro,Apex O |
| 5310 | QNAP Helpdesk Improper Access Control Vulnerability | kev | markdown |  | CVE-2020-2506,cve-arsenal,cisa-kev,kev,in-the-wild,QNAP Systems,Helpde |
| 5311 | D-Link DCS-2530L and DCS-2670L Devices Unspecified Vulnerability | kev | markdown |  | CVE-2020-25078,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DCS-2530L a |
| 5312 | D-Link DCS-2530L and DCS-2670L Command Injection Vulnerability | kev | markdown |  | CVE-2020-25079,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DCS-2530L a |
| 5313 | QNAP Network-Attached Storage (NAS) Command Injection Vulnerability | kev | markdown |  | CVE-2020-2509,cve-arsenal,cisa-kev,kev,in-the-wild,QNAP,QNAP Network-A |
| 5314 | WordPress File Manager Plugin Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-25213,cve-arsenal,cisa-kev,kev,in-the-wild,WordPress,File Man |
| 5315 | Sophos SG UTM Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-25223,cve-arsenal,cisa-kev,kev,in-the-wild,Sophos,SG UTM |
| 5316 | D-Link DNS-320 Device Command Injection Vulnerability | kev | markdown |  | CVE-2020-25506,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DNS-320 Dev |
| 5317 | Oracle Fusion Middleware Unspecified Vulnerability | kev | markdown |  | CVE-2020-2551,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,Fusion Middl |
| 5318 | Oracle Multiple Products Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-2555,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,Multiple Pro |
| 5319 | Netgear JGS516PE Devices Missing Function Level Access Control Vulnerability | kev | markdown |  | CVE-2020-26919,cve-arsenal,cisa-kev,kev,in-the-wild,NETGEAR,JGS516PE D |
| 5320 | Apple Multiple Products Memory Corruption Vulnerability | kev | markdown |  | CVE-2020-27930,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5321 | Apple Multiple Products Type Confusion Vulnerability | kev | markdown |  | CVE-2020-27932,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5322 | Apple Multiple Products Memory Initialization Vulnerability | kev | markdown |  | CVE-2020-27950,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5323 | Oracle WebLogic Server Unspecified Vulnerability | kev | markdown |  | CVE-2020-2883,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,WebLogic Ser |
| 5324 | PEAR Archive_Tar Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2020-28949,cve-arsenal,cisa-kev,kev,in-the-wild,PEAR,Archive_Tar |
| 5325 | D-Link DIR-825 R1 Devices Buffer Overflow Vulnerability | kev | markdown |  | CVE-2020-29557,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DIR-825 R1  |
| 5326 | CyberoamOS (CROS) SQL Injection Vulnerability | kev | markdown |  | CVE-2020-29574,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Sophos, |
| 5327 | Zyxel Multiple Products Use of Hard-Coded Credentials Vulnerability | kev | markdown |  | CVE-2020-29583,cve-arsenal,cisa-kev,kev,in-the-wild,Zyxel,Multiple Pro |
| 5328 | Cisco IOS XR Software Discovery Protocol Format String Vulnerability | kev | markdown |  | CVE-2020-3118,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS XR |
| 5329 | Cisco AnyConnect Secure Mobility Client for Windows Uncontrolled Search Path Vulnerability | kev | markdown |  | CVE-2020-3153,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Cisco,An |
| 5330 | Cisco IP Phones Web Server Remote Code Execution and Denial-of-Service Vulnerability | kev | markdown |  | CVE-2020-3161,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Cisco IP Phon |
| 5331 | Cisco ASA and FTD Information Disclosure Vulnerability | kev | markdown |  | CVE-2020-3259,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Cisco,Ad |
| 5332 | Cisco AnyConnect Secure Mobility Client for Windows DLL Hijacking Vulnerability | kev | markdown |  | CVE-2020-3433,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Cisco,An |
| 5333 | Cisco ASA and FTD Read-Only Path Traversal Vulnerability | kev | markdown |  | CVE-2020-3452,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Adaptive Secu |
| 5334 | Cisco IOS XR Software DVMRP Memory Exhaustion Vulnerability | kev | markdown |  | CVE-2020-3566,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS XR |
| 5335 | Cisco IOS XR Software DVMRP Memory Exhaustion Vulnerability | kev | markdown |  | CVE-2020-3569,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS XR |
| 5336 | Roundcube Webmail Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2020-35730,cve-arsenal,cisa-kev,kev,in-the-wild,Roundcube,Roundcub |
| 5337 | Cisco ASA and FTD Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2020-3580,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Cisco,Ad |
| 5338 | PEAR Archive_Tar Improper Link Resolution Vulnerability | kev | markdown |  | CVE-2020-36193,cve-arsenal,cisa-kev,kev,in-the-wild,PEAR,Archive_Tar |
| 5339 | Apple Multiple Products Memory Corruption Vulnerability | kev | markdown |  | CVE-2020-3837,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Prod |
| 5340 | VMware Multiple Products Privilege Escalation Vulnerability | kev | markdown |  | CVE-2020-3950,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,Multiple Pro |
| 5341 | VMware vCenter Server Information Disclosure Vulnerability | kev | markdown |  | CVE-2020-3952,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,vCenter Serv |
| 5342 | VMware ESXi OpenSLP Use-After-Free Vulnerability | kev | markdown |  | CVE-2020-3992,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,VMware,E |
| 5343 | Multiple VMware Products Command Injection Vulnerability | kev | markdown |  | CVE-2020-4006,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,Multiple Pro |
| 5344 | IBM Data Risk Manager Security Bypass Vulnerability | kev | markdown |  | CVE-2020-4427,cve-arsenal,cisa-kev,kev,in-the-wild,IBM,Data Risk Manag |
| 5345 | IBM Data Risk Manager Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-4428,cve-arsenal,cisa-kev,kev,in-the-wild,IBM,Data Risk Manag |
| 5346 | IBM Data Risk Manager Directory Traversal Vulnerability | kev | markdown |  | CVE-2020-4430,cve-arsenal,cisa-kev,kev,in-the-wild,IBM,Data Risk Manag |
| 5347 | SonicWall SonicOS Buffer Overflow Vulnerability | kev | markdown |  | CVE-2020-5135,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SonicWal |
| 5348 | VMware Tanzu Spring Cloud Config Directory Traversal Vulnerability | kev | markdown |  | CVE-2020-5410,cve-arsenal,cisa-kev,kev,in-the-wild,VMware Tanzu,Spring |
| 5349 | Grandstream Networks UCM6200 Series SQL Injection Vulnerability | kev | markdown |  | CVE-2020-5722,cve-arsenal,cisa-kev,kev,in-the-wild,Grandstream,UCM6200 |
| 5350 | Amcrest Cameras and NVR Stack-based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2020-5735,cve-arsenal,cisa-kev,kev,in-the-wild,Amcrest,Cameras and |
| 5351 | Plex Media Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-5741,cve-arsenal,cisa-kev,kev,in-the-wild,Plex,Media Server |
| 5352 | Unraid Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-5847,cve-arsenal,cisa-kev,kev,in-the-wild,Unraid |
| 5353 | Unraid Authentication Bypass Vulnerability | kev | markdown |  | CVE-2020-5849,cve-arsenal,cisa-kev,kev,in-the-wild,Unraid |
| 5354 | F5 BIG-IP Traffic Management User Interface (TMUI) Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-5902,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,F5,BIG-I |
| 5355 | SAP Solution Manager Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2020-6207,cve-arsenal,cisa-kev,kev,in-the-wild,SAP,Solution Manage |
| 5356 | SAP NetWeaver Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2020-6287,cve-arsenal,cisa-kev,kev,in-the-wild,SAP,NetWeaver |
| 5357 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2020-6418,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5358 | Google Chrome Media Use-After-Free Vulnerability | kev | markdown |  | CVE-2020-6572,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chrome Media |
| 5359 | Mozilla Firefox And Thunderbird Use-After-Free Vulnerability | kev | markdown |  | CVE-2020-6819,cve-arsenal,cisa-kev,kev,in-the-wild,Mozilla,Firefox and |
| 5360 | Mozilla Firefox And Thunderbird Use-After-Free Vulnerability | kev | markdown |  | CVE-2020-6820,cve-arsenal,cisa-kev,kev,in-the-wild,Mozilla,Firefox and |
| 5361 | OpenSMTPD Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-7247,cve-arsenal,cisa-kev,kev,in-the-wild,OpenBSD,OpenSMTPD |
| 5362 | Synacor Zimbra Collaboration Suite (ZCS) Server-Side Request Forgery Vulnerability | kev | markdown |  | CVE-2020-7796,cve-arsenal,cisa-kev,kev,in-the-wild,Synacor,Zimbra Coll |
| 5363 | Liferay Portal Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2020-7961,cve-arsenal,cisa-kev,kev,in-the-wild,Liferay,Liferay Por |
| 5364 | Citrix ADC, Gateway, and SD-WAN WANOP Appliance Authorization Bypass Vulnerability | kev | markdown |  | CVE-2020-8193,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,Application  |
| 5365 | Citrix ADC, Gateway, and SD-WAN WANOP Appliance Information Disclosure Vulnerability | kev | markdown |  | CVE-2020-8195,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,Application  |
| 5366 | Citrix ADC, Gateway, and SD-WAN WANOP Appliance Information Disclosure Vulnerability | kev | markdown |  | CVE-2020-8196,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,Application  |
| 5367 | Pulse Connect Secure Code Injection Vulnerability | kev | markdown |  | CVE-2020-8218,cve-arsenal,cisa-kev,kev,in-the-wild,Pulse Secure,Pulse  |
| 5368 | Ivanti Pulse Connect Secure Code Execution Vulnerability | kev | markdown |  | CVE-2020-8243,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Pulse Connec |
| 5369 | Ivanti Pulse Connect Secure Code Execution Vulnerability | kev | markdown |  | CVE-2020-8260,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Pulse Connec |
| 5370 | Trend Micro Apex One and OfficeScan Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-8467,cve-arsenal,cisa-kev,kev,in-the-wild,Trend Micro,Apex On |
| 5371 | Trend Micro Multiple Products Content Validation Escape Vulnerability | kev | markdown |  | CVE-2020-8468,cve-arsenal,cisa-kev,kev,in-the-wild,Trend Micro,Apex On |
| 5372 | Multiple DrayTek Vigor Routers Web Management Page Vulnerability | kev | markdown |  | CVE-2020-8515,cve-arsenal,cisa-kev,kev,in-the-wild,DrayTek,Multiple Vi |
| 5373 | Trend Micro Apex One and OfficeScan Authentication Bypass Vulnerability | kev | markdown |  | CVE-2020-8599,cve-arsenal,cisa-kev,kev,in-the-wild,Trend Micro,Apex On |
| 5374 | PlaySMS Server-Side Template Injection Vulnerability | kev | markdown |  | CVE-2020-8644,cve-arsenal,cisa-kev,kev,in-the-wild,PlaySMS |
| 5375 | EyesOfNetwork Improper Privilege Management Vulnerability | kev | markdown |  | CVE-2020-8655,cve-arsenal,cisa-kev,kev,in-the-wild,EyesOfNetwork |
| 5376 | EyesOfNetwork Use of Hard-Coded Credentials Vulnerability | kev | markdown |  | CVE-2020-8657,cve-arsenal,cisa-kev,kev,in-the-wild,EyesOfNetwork |
| 5377 | Pi-Hole AdminLTE Remote Code Execution Vulnerability | kev | markdown |  | CVE-2020-8816,cve-arsenal,cisa-kev,kev,in-the-wild,Pi-hole,AdminLTE |
| 5378 | Zyxel Multiple NAS Devices OS Command Injection Vulnerability | kev | markdown |  | CVE-2020-9054,cve-arsenal,cisa-kev,kev,in-the-wild,Zyxel,Multiple Netw |
| 5379 | D-Link DIR-610 Devices Remote Command Execution | kev | markdown |  | CVE-2020-9377,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DIR-610 Devi |
| 5380 | Adobe Acrobat Use-After-Free Vulnerability | kev | markdown |  | CVE-2020-9715,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat |
| 5381 | Apple iOS, iPadOS, and watchOS Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2020-9818,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS,  |
| 5382 | Apple iOS, iPadOS, and watchOS Memory Corruption Vulnerability | kev | markdown |  | CVE-2020-9819,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS,  |
| 5383 | Apple Multiple Products Code Execution Vulnerability | kev | markdown |  | CVE-2020-9859,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Prod |
| 5384 | Apple Multiple Products Memory Corruption Vulnerability | kev | markdown |  | CVE-2020-9907,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Prod |
| 5385 | Apple iOS, iPadOS, and macOS Input Validation Vulnerability | kev | markdown |  | CVE-2020-9934,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS,  |
| 5386 | Android Kernel Race Condition Vulnerability | kev | markdown |  | CVE-2021-0920,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Kernel |
| 5387 | Android Kernel Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-1048,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Kernel |
| 5388 | Cisco HyperFlex HX Installer Virtual Machine Command Injection Vulnerability | kev | markdown |  | CVE-2021-1497,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,HyperFlex HX |
| 5389 | Cisco HyperFlex HX Data Platform Command Injection Vulnerability | kev | markdown |  | CVE-2021-1498,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,HyperFlex HX |
| 5390 | Microsoft Defender Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-1647,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Defender |
| 5391 | Microsoft Windows Print Spooler Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-1675,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5392 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-1732,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microsof |
| 5393 | Apple Multiple Products Race Condition Vulnerability | kev | markdown |  | CVE-2021-1782,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Prod |
| 5394 | Apple Multiple Products Type Confusion Vulnerability | kev | markdown |  | CVE-2021-1789,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Prod |
| 5395 | Apple iOS, iPadOS, and macOS WebKit Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-1870,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS,  |
| 5396 | Apple iOS, iPadOS, and macOS WebKit Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-1871,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS,  |
| 5397 | Apple iOS, iPadOS, and watchOS WebKit Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2021-1879,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS,  |
| 5398 | Qualcomm Multiple Chipsets Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-1905,cve-arsenal,cisa-kev,kev,in-the-wild,Qualcomm,Multiple C |
| 5399 | Qualcomm Multiple Chipsets Detection of Error Condition Without Action Vulnerability | kev | markdown |  | CVE-2021-1906,cve-arsenal,cisa-kev,kev,in-the-wild,Qualcomm,Multiple C |
| 5400 | SonicWall SSLVPN SMA100 SQL Injection Vulnerability | kev | markdown |  | CVE-2021-20016,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SonicWa |
| 5401 | SonicWall Email Security Improper Privilege Management Vulnerability | kev | markdown |  | CVE-2021-20021,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SonicWa |
| 5402 | SonicWall Email Security Unrestricted Upload of File Vulnerability | kev | markdown |  | CVE-2021-20022,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SonicWa |
| 5403 | SonicWall Email Security Path Traversal Vulnerability | kev | markdown |  | CVE-2021-20023,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SonicWa |
| 5404 | SonicWall Secure Remote Access (SRA) SQL Injection Vulnerability | kev | markdown |  | CVE-2021-20028,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SonicWa |
| 5405 | SonicWall SMA100 Appliances OS Command Injection Vulnerability | kev | markdown |  | CVE-2021-20035,cve-arsenal,cisa-kev,kev,in-the-wild,SonicWall,SMA100 A |
| 5406 | SonicWall SMA 100 Appliances Stack-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2021-20038,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SonicWa |
| 5407 | Arcadyan Buffalo Firmware Path Traversal Vulnerability | kev | markdown |  | CVE-2021-20090,cve-arsenal,cisa-kev,kev,in-the-wild,Arcadyan,Buffalo F |
| 5408 | Draytek VigorConnect Path Traversal Vulnerability | kev | markdown |  | CVE-2021-20123,cve-arsenal,cisa-kev,kev,in-the-wild,DrayTek,VigorConne |
| 5409 | Draytek VigorConnect Path Traversal Vulnerability | kev | markdown |  | CVE-2021-20124,cve-arsenal,cisa-kev,kev,in-the-wild,DrayTek,VigorConne |
| 5410 | Adobe Acrobat and Reader Heap-based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2021-21017,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat and  |
| 5411 | Google Chromium V8 Heap Buffer Overflow Vulnerability | kev | markdown |  | CVE-2021-21148,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5412 | Google Chromium Race Condition Vulnerability | kev | markdown |  | CVE-2021-21166,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium |
| 5413 | Google Chromium Blink Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-21193,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium Bl |
| 5414 | Google Chromium Blink Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-21206,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium Bl |
| 5415 | Google Chromium V8 Improper Input Validation Vulnerability | kev | markdown |  | CVE-2021-21220,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5416 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2021-21224,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5417 | Adminer Server-Side Request Forgery Vulnerability | kev | markdown |  | CVE-2021-21311,cve-arsenal,cisa-kev,kev,in-the-wild,Adminer |
| 5418 | System Information Library for Node.JS Command Injection | kev | markdown |  | CVE-2021-21315,cve-arsenal,cisa-kev,kev,in-the-wild,Npm package,System |
| 5419 | Dell dbutil Driver Insufficient Access Control Vulnerability | kev | markdown |  | CVE-2021-21551,cve-arsenal,cisa-kev,kev,in-the-wild,Dell,dbutil Driver |
| 5420 | VMware vCenter Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-21972,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,VMware, |
| 5421 | VMware vCenter Server and Cloud Foundation Server Side Request Forgery (SSRF) Vulnerability | kev | markdown |  | CVE-2021-21973,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,vCenter Ser |
| 5422 | VMware Server Side Request Forgery in vRealize Operations Manager API | kev | markdown |  | CVE-2021-21975,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,VMware, |
| 5423 | VMware vCenter Server Improper Input Validation Vulnerability | kev | markdown |  | CVE-2021-21985,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,VMware, |
| 5424 | VMware vCenter Server File Upload Vulnerability | kev | markdown |  | CVE-2021-22005,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,VMware, |
| 5425 | VMware vCenter Server Improper Access Control | kev | markdown |  | CVE-2021-22017,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,vCenter Ser |
| 5426 | Omnissa Workspace ONE Server-Side Request Forgery | kev | markdown |  | CVE-2021-22054,cve-arsenal,cisa-kev,kev,in-the-wild,Omnissa,Workspace  |
| 5427 | GitLab Server-Side Request Forgery (SSRF) Vulnerability | kev | markdown |  | CVE-2021-22175,cve-arsenal,cisa-kev,kev,in-the-wild,GitLab |
| 5428 | ExifTool Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-22204,cve-arsenal,cisa-kev,kev,in-the-wild,Perl,Exiftool |
| 5429 | GitLab Community and Enterprise Editions Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-22205,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,GitLab, |
| 5430 | Micro Focus Operation Bridge Report (OBR) Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-22502,cve-arsenal,cisa-kev,kev,in-the-wild,Micro Focus,Operat |
| 5431 | Micro Focus Access Manager Information Leakage Vulnerability | kev | markdown |  | CVE-2021-22506,cve-arsenal,cisa-kev,kev,in-the-wild,Micro Focus,Micro  |
| 5432 | Linux Kernel Heap Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2021-22555,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 5433 | Linux Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-22600,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 5434 | Rockwell Multiple Products Insufficient Protected Credentials Vulnerability | kev | markdown |  | CVE-2021-22681,cve-arsenal,cisa-kev,kev,in-the-wild,Rockwell,Multiple  |
| 5435 | Ivanti Pulse Connect Secure Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-22893,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Ivanti, |
| 5436 | Ivanti Pulse Connect Secure Collaboration Suite Buffer Overflow Vulnerability | kev | markdown |  | CVE-2021-22894,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Pulse Conne |
| 5437 | Ivanti Pulse Connect Secure Command Injection Vulnerability | kev | markdown |  | CVE-2021-22899,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Pulse Conne |
| 5438 | Ivanti Pulse Connect Secure Unrestricted File Upload Vulnerability | kev | markdown |  | CVE-2021-22900,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Pulse Conne |
| 5439 | Citrix ShareFile Improper Access Control Vulnerability | kev | markdown |  | CVE-2021-22941,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Citrix, |
| 5440 | F5 BIG-IP and BIG-IQ Centralized Management iControl REST Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-22986,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,F5,BIG- |
| 5441 | F5 BIG-IP Traffic Management Microkernel Buffer Overflow | kev | markdown |  | CVE-2021-22991,cve-arsenal,cisa-kev,kev,in-the-wild,F5,BIG-IP Traffic  |
| 5442 | Ajax.NET Professional Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2021-23758,cve-arsenal,cisa-kev,kev,in-the-wild,Ajax.NET Professio |
| 5443 | McAfee Total Protection (MTP) Improper Privilege Management Vulnerability | kev | markdown |  | CVE-2021-23874,cve-arsenal,cisa-kev,kev,in-the-wild,McAfee,McAfee Tota |
| 5444 | Nagios XI OS Command Injection | kev | markdown |  | CVE-2021-25296,cve-arsenal,cisa-kev,kev,in-the-wild,Nagios,Nagios XI |
| 5445 | Nagios XI OS Command Injection | kev | markdown |  | CVE-2021-25297,cve-arsenal,cisa-kev,kev,in-the-wild,Nagios,Nagios XI |
| 5446 | Nagios XI OS Command Injection | kev | markdown |  | CVE-2021-25298,cve-arsenal,cisa-kev,kev,in-the-wild,Nagios,Nagios XI |
| 5447 | Samsung Mobile Devices Improper Access Control Vulnerability | kev | markdown |  | CVE-2021-25337,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,Mobile Dev |
| 5448 | Samsung Mobile Devices Improper Access Control Vulnerability | kev | markdown |  | CVE-2021-25369,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,Mobile Dev |
| 5449 | Samsung Mobile Devices Memory Corruption Vulnerability | kev | markdown |  | CVE-2021-25370,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,Mobile Dev |
| 5450 | Samsung Mobile Devices Unspecified Vulnerability | kev | markdown |  | CVE-2021-25371,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,Mobile Dev |
| 5451 | Samsung Mobile Devices Improper Boundary Check Vulnerability | kev | markdown |  | CVE-2021-25372,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,Mobile Dev |
| 5452 | Samsung Mobile Devices Race Condition Vulnerability | kev | markdown |  | CVE-2021-25394,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,Mobile Dev |
| 5453 | Samsung Mobile Devices Race Condition Vulnerability | kev | markdown |  | CVE-2021-25395,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,Mobile Dev |
| 5454 | Samsung Mobile Devices Out-of-Bounds Read Vulnerability | kev | markdown |  | CVE-2021-25487,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,Mobile Dev |
| 5455 | Samsung Mobile Devices Improper Input Validation Vulnerability | kev | markdown |  | CVE-2021-25489,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,Mobile Dev |
| 5456 | Atlassian Confluence Server and Data Center Object-Graph Navigation Language (OGNL) Injection Vulnerability | kev | markdown |  | CVE-2021-26084,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Atlassi |
| 5457 | Atlassian Confluence Server Pre-Authorization Arbitrary File Read Vulnerability | kev | markdown |  | CVE-2021-26085,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Atlassi |
| 5458 | Atlassian Jira Server and Data Center Path Traversal Vulnerability | kev | markdown |  | CVE-2021-26086,cve-arsenal,cisa-kev,kev,in-the-wild,Atlassian,Jira Ser |
| 5459 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | markdown |  | CVE-2021-26411,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5460 | OpenPLC ScadaBR Unrestricted Upload of File with Dangerous Type Vulnerability | kev | markdown |  | CVE-2021-26828,cve-arsenal,cisa-kev,kev,in-the-wild,OpenPLC,ScadaBR |
| 5461 | OpenPLC ScadaBR Cross-site Scripting Vulnerability | kev | markdown |  | CVE-2021-26829,cve-arsenal,cisa-kev,kev,in-the-wild,OpenPLC,ScadaBR |
| 5462 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-26855,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5463 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-26857,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5464 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-26858,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5465 | Microsoft Office Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-27059,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 5466 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-27065,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5467 | Microsoft Internet Explorer Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-27085,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet |
| 5468 | Accellion FTA SQL Injection Vulnerability | kev | markdown |  | CVE-2021-27101,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Accelli |
| 5469 | Accellion FTA OS Command Injection Vulnerability | kev | markdown |  | CVE-2021-27102,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Accelli |
| 5470 | Accellion FTA Server-Side Request Forgery (SSRF) Vulnerability | kev | markdown |  | CVE-2021-27103,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Accelli |
| 5471 | Accellion FTA OS Command Injection Vulnerability | kev | markdown |  | CVE-2021-27104,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Accelli |
| 5472 | DD-WRT Stack-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2021-27137,cve-arsenal,cisa-kev,kev,in-the-wild,DD-WRT |
| 5473 | Yealink Device Management Server-Side Request Forgery (SSRF) Vulnerability | kev | markdown |  | CVE-2021-27561,cve-arsenal,cisa-kev,kev,in-the-wild,Yealink,Device Man |
| 5474 | Arm Trusted Firmware Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2021-27562,cve-arsenal,cisa-kev,kev,in-the-wild,Arm,Trusted Firmwa |
| 5475 | Checkbox Survey Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2021-27852,cve-arsenal,cisa-kev,kev,in-the-wild,Checkbox,Checkbox  |
| 5476 | FatPipe WARP, IPVPN, and MPVPN Configuration Upload exploit | kev | markdown |  | CVE-2021-27860,cve-arsenal,cisa-kev,kev,in-the-wild,FatPipe,WARP, IPVP |
| 5477 | Veritas Backup Exec Agent File Access Vulnerability | kev | markdown |  | CVE-2021-27876,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Veritas |
| 5478 | Veritas Backup Exec Agent Improper Authentication Vulnerability | kev | markdown |  | CVE-2021-27877,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Veritas |
| 5479 | Veritas Backup Exec Agent Command Execution Vulnerability | kev | markdown |  | CVE-2021-27878,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Veritas |
| 5480 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-28310,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 5481 | Adobe Acrobat and Reader Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-28550,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat and  |
| 5482 | Arm Mali Graphics Processing Unit (GPU) Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-28663,cve-arsenal,cisa-kev,kev,in-the-wild,Arm,Mali Graphics  |
| 5483 | Arm Mali Graphics Processing Unit (GPU) Unspecified Vulnerability | kev | markdown |  | CVE-2021-28664,cve-arsenal,cisa-kev,kev,in-the-wild,Arm,Mali Graphics  |
| 5484 | QNAP NAS Improper Authorization Vulnerability | kev | markdown |  | CVE-2021-28799,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,QNAP,Ne |
| 5485 | Arm Mali GPU Kernel Driver Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-29256,cve-arsenal,cisa-kev,kev,in-the-wild,Arm,Mali Graphics  |
| 5486 | Kaseya Virtual System/Server Administrator (VSA) Information Disclosure Vulnerability | kev | markdown |  | CVE-2021-30116,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Kaseya, |
| 5487 | Google Chromium PopupBlocker Security Bypass Vulnerability | kev | markdown |  | CVE-2021-30533,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium Po |
| 5488 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2021-30551,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5489 | Google Chromium WebGL Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-30554,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium We |
| 5490 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2021-30563,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5491 | Google Chromium V8 Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2021-30632,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5492 | Google Chromium Indexed DB API Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-30633,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium In |
| 5493 | Apple macOS Unspecified Vulnerability | kev | markdown |  | CVE-2021-30657,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,macOS |
| 5494 | Apple Multiple Products WebKit Storage Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-30661,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5495 | Apple Multiple Products WebKit Integer Overflow Vulnerability | kev | markdown |  | CVE-2021-30663,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5496 | Apple Multiple Products WebKit Memory Corruption Vulnerability | kev | markdown |  | CVE-2021-30665,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5497 | Apple iOS WebKit Buffer Overflow Vulnerability | kev | markdown |  | CVE-2021-30666,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS |
| 5498 | Apple macOS Unspecified Vulnerability | kev | markdown |  | CVE-2021-30713,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,macOS |
| 5499 | Apple iOS WebKit Memory Corruption Vulnerability | kev | markdown |  | CVE-2021-30761,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS |
| 5500 | Apple iOS WebKit Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-30762,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS |
| 5501 | Apple Multiple Products Memory Corruption Vulnerability | kev | markdown |  | CVE-2021-30807,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5502 | Apple iOS, iPadOS, macOS Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-30858,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS, |
| 5503 | Apple Multiple Products Integer Overflow Vulnerability | kev | markdown |  | CVE-2021-30860,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5504 | Apple iOS, iPadOS, and macOS Type Confusion Vulnerability | kev | markdown |  | CVE-2021-30869,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS, |
| 5505 | Apple Multiple Products Memory Corruption Vulnerability | kev | markdown |  | CVE-2021-30883,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5506 | Apple iOS, iPadOS, and macOS Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2021-30900,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS, |
| 5507 | Apple Multiple Products Integer Overflow or Wraparound Vulnerability | kev | markdown |  | CVE-2021-30952,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5508 | Apple iOS and iPadOS Buffer Overflow Vulnerability | kev | markdown |  | CVE-2021-30983,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS and iPad |
| 5509 | Apple iOS, macOS, watchOS Sandbox Bypass Vulnerability | kev | markdown |  | CVE-2021-31010,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, macOS,  |
| 5510 | Microsoft HTTP Protocol Stack Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-31166,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,HTTP Pro |
| 5511 | Microsoft Exchange Server Information Disclosure Vulnerability | kev | markdown |  | CVE-2021-31196,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Exchange |
| 5512 | Microsoft Enhanced Cryptographic Provider Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-31199,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Enhanced |
| 5513 | Microsoft Enhanced Cryptographic Provider Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-31201,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Enhanced |
| 5514 | Microsoft Exchange Server Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2021-31207,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5515 | Laravel Ignition File Upload Vulnerability | kev | markdown |  | CVE-2021-3129,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Laravel, |
| 5516 | Sudo Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2021-3156,cve-arsenal,cisa-kev,kev,in-the-wild,Sudo |
| 5517 | Tenda AC11 Router Stack Buffer Overflow Vulnerability | kev | markdown |  | CVE-2021-31755,cve-arsenal,cisa-kev,kev,in-the-wild,Tenda,AC11 Router |
| 5518 | Microsoft Windows Kernel Information Disclosure Vulnerability | kev | markdown |  | CVE-2021-31955,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5519 | Microsoft Windows NTFS Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-31956,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5520 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-31979,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5521 | ONLYOFFICE Docs Server Path Traversal Vulnerability | kev | markdown |  | CVE-2021-3199,cve-arsenal,cisa-kev,kev,in-the-wild,ONLYOFFICE,Docs |
| 5522 | ASUS Routers Improper Authentication Vulnerability | kev | markdown |  | CVE-2021-32030,cve-arsenal,cisa-kev,kev,in-the-wild,ASUS,Routers |
| 5523 | October CMS Improper Authentication | kev | markdown |  | CVE-2021-32648,cve-arsenal,cisa-kev,kev,in-the-wild,October CMS |
| 5524 | Dahua IP Camera Authentication Bypass Vulnerability | kev | markdown |  | CVE-2021-33044,cve-arsenal,cisa-kev,kev,in-the-wild,Dahua,IP Camera Fi |
| 5525 | Dahua IP Camera Authentication Bypass Vulnerability | kev | markdown |  | CVE-2021-33045,cve-arsenal,cisa-kev,kev,in-the-wild,Dahua,IP Camera Fi |
| 5526 | Microsoft Desktop Window Manager (DWM) Core Library Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-33739,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5527 | Microsoft Windows MSHTML Platform Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-33742,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5528 | Microsoft Exchange Server Information Disclosure | kev | markdown |  | CVE-2021-33766,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Exchange |
| 5529 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-33771,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5530 | Microsoft Windows Scripting Engine Memory Corruption Vulnerability | kev | markdown |  | CVE-2021-34448,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5531 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-34473,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5532 | Microsoft Windows User Profile Service Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-34484,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5533 | Microsoft Windows Event Tracing Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-34486,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5534 | Microsoft Exchange Server Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-34523,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5535 | Microsoft Windows Print Spooler Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-34527,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5536 | Linux Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-3493,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 5537 | SolarWinds Serv-U Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-35211,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SolarWi |
| 5538 | SolarWinds Serv-U Improper Input Validation Vulnerability | kev | markdown |  | CVE-2021-35247,cve-arsenal,cisa-kev,kev,in-the-wild,SolarWinds,Serv-U |
| 5539 | Realtek Jungle SDK Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-35394,cve-arsenal,cisa-kev,kev,in-the-wild,Realtek,Jungle Sof |
| 5540 | Realtek AP-Router SDK Buffer Overflow Vulnerability | kev | markdown |  | CVE-2021-35395,cve-arsenal,cisa-kev,kev,in-the-wild,Realtek,AP-Router  |
| 5541 | ForgeRock Access Management (AM) Core Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-35464,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,ForgeRo |
| 5542 | Oracle Fusion Middleware Unspecified Vulnerability | kev | markdown |  | CVE-2021-35587,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,Fusion Midd |
| 5543 | Red Hat Polkit Incorrect Authorization Vulnerability | kev | markdown |  | CVE-2021-3560,cve-arsenal,cisa-kev,kev,in-the-wild,Red Hat,Polkit |
| 5544 | Hikvision Improper Input Validation | kev | markdown |  | CVE-2021-36260,cve-arsenal,cisa-kev,kev,in-the-wild,Hikvision,Security |
| 5545 | Sunhillo SureLine OS Command Injection Vulnerablity | kev | markdown |  | CVE-2021-36380,cve-arsenal,cisa-kev,kev,in-the-wild,Sunhillo,SureLine |
| 5546 | Trend Micro Multiple Products Improper Input Validation Vulnerability | kev | markdown |  | CVE-2021-36741,cve-arsenal,cisa-kev,kev,in-the-wild,Trend Micro,Apex O |
| 5547 | Trend Micro Multiple Products Improper Input Validation Vulnerability | kev | markdown |  | CVE-2021-36742,cve-arsenal,cisa-kev,kev,in-the-wild,Trend Micro,Apex O |
| 5548 | Microsoft Windows SAM Local Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-36934,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5549 | Microsoft Windows Local Security Authority (LSA) Spoofing Vulnerability | kev | markdown |  | CVE-2021-36942,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5550 | Microsoft Windows Update Medic Service Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-36948,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5551 | Microsoft Windows Common Log File System (CLFS) Driver Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-36955,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5552 | Zoho ManageEngine ServiceDesk Authentication Bypass Vulnerability | kev | markdown |  | CVE-2021-37415,cve-arsenal,cisa-kev,kev,in-the-wild,Zoho,ManageEngine  |
| 5553 | Google Chromium Portals Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-37973,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium Po |
| 5554 | Google Chromium V8 Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-37975,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5555 | Google Chromium Information Disclosure Vulnerability | kev | markdown |  | CVE-2021-37976,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium |
| 5556 | Google Chromium Intents Improper Input Validation Vulnerability | kev | markdown |  | CVE-2021-38000,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium In |
| 5557 | Google Chromium V8 Memory Corruption Vulnerability | kev | markdown |  | CVE-2021-38003,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5558 | SAP NetWeaver Unrestricted File Upload Vulnerability | kev | markdown |  | CVE-2021-38163,cve-arsenal,cisa-kev,kev,in-the-wild,SAP,NetWeaver |
| 5559 | Delta Electronics DOPSoft 2 Improper Input Validation Vulnerability | kev | markdown |  | CVE-2021-38406,cve-arsenal,cisa-kev,kev,in-the-wild,Delta Electronics, |
| 5560 | Microsoft Open Management Infrastructure (OMI) Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-38645,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Open Man |
| 5561 | Microsoft Office Access Connectivity Engine Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-38646,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5562 | Microsoft Open Management Infrastructure (OMI) Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-38647,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5563 | Microsoft Open Management Infrastructure (OMI) Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-38648,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Open Man |
| 5564 | Microsoft Open Management Infrastructure (OMI) Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-38649,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Open Man |
| 5565 | XStream Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-39144,cve-arsenal,cisa-kev,kev,in-the-wild,XStream |
| 5566 | Grafana Authentication Bypass Vulnerability | kev | markdown |  | CVE-2021-39226,cve-arsenal,cisa-kev,kev,in-the-wild,Grafana Labs,Grafa |
| 5567 | Google Pixel Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2021-39793,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Pixel |
| 5568 | GitLab Community and Enterprise Editions Server-Side Request Forgery (SSRF) Vulnerability | kev | markdown |  | CVE-2021-39935,cve-arsenal,cisa-kev,kev,in-the-wild,GitLab,Community a |
| 5569 | Red Hat Polkit Out-of-Bounds Read and Write Vulnerability | kev | markdown |  | CVE-2021-4034,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Red Hat, |
| 5570 | Reolink RLC-410W IP Camera OS Command Injection Vulnerability | kev | markdown |  | CVE-2021-40407,cve-arsenal,cisa-kev,kev,in-the-wild,Reolink,RLC-410W I |
| 5571 | Apache HTTP Server-Side Request Forgery (SSRF) | kev | markdown |  | CVE-2021-40438,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Apache |
| 5572 | Microsoft MSHTML Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-40444,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5573 | Microsoft Windows Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-40449,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5574 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-40450,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 5575 | Zoho ManageEngine ADSelfService Plus Authentication Bypass Vulnerability | kev | markdown |  | CVE-2021-40539,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Zoho,Ma |
| 5576 | D-Link DIR-605 Router Information Disclosure Vulnerability | kev | markdown |  | CVE-2021-40655,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DIR-605 Rou |
| 5577 | Aviatrix Controller Unrestricted Upload of File | kev | markdown |  | CVE-2021-40870,cve-arsenal,cisa-kev,kev,in-the-wild,Aviatrix,Aviatrix  |
| 5578 | Google Chromium V8 Use-After-Free Vulnerability | kev | markdown |  | CVE-2021-4102,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5579 | Metabase GeoJSON API Local File Inclusion Vulnerability | kev | markdown |  | CVE-2021-41277,cve-arsenal,cisa-kev,kev,in-the-wild,Metabase |
| 5580 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-41357,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 5581 | Microsoft Windows Installer Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-41379,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5582 | Apache HTTP Server Path Traversal Vulnerability | kev | markdown |  | CVE-2021-41773,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Apache, |
| 5583 | Apache HTTP Server Path Traversal Vulnerability | kev | markdown |  | CVE-2021-42013,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Apache, |
| 5584 | Sitecore XP Remote Command Execution Vulnerability | kev | markdown |  | CVE-2021-42237,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Sitecor |
| 5585 | BQE BillQuick Web Suite SQL Injection Vulnerability | kev | markdown |  | CVE-2021-42258,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,BQE,Bil |
| 5586 | Microsoft Active Directory Domain Services Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-42278,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5587 | Microsoft Active Directory Domain Services Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-42287,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5588 | Microsoft Excel Security Feature Bypass | kev | markdown |  | CVE-2021-42292,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 5589 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-42321,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5590 | Microsoft Windows Privilege Escalation Vulnerability | kev | markdown |  | CVE-2021-43226,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5591 | Grafana Path Traversal Vulnerability | kev | markdown |  | CVE-2021-43798,cve-arsenal,cisa-kev,kev,in-the-wild,Grafana Labs,Grafa |
| 5592 | Microsoft Windows AppX Installer Spoofing Vulnerability | kev | markdown |  | CVE-2021-43890,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5593 | Roundcube Webmail SQL Injection Vulnerability | kev | markdown |  | CVE-2021-44026,cve-arsenal,cisa-kev,kev,in-the-wild,Roundcube,Roundcub |
| 5594 | Zoho ManageEngine ServiceDesk Plus Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-44077,cve-arsenal,cisa-kev,kev,in-the-wild,Zoho,ManageEngine  |
| 5595 | Fortinet FortiOS Arbitrary File Download | kev | markdown |  | CVE-2021-44168,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,FortiOS |
| 5596 | Acclaim Systems USAHERDS Use of Hard-Coded Credentials Vulnerability | kev | markdown |  | CVE-2021-44207,cve-arsenal,cisa-kev,kev,in-the-wild,Acclaim Systems,US |
| 5597 | Apache Log4j2 Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-44228,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Apache, |
| 5598 | Zoho Desktop Central Authentication Bypass Vulnerability | kev | markdown |  | CVE-2021-44515,cve-arsenal,cisa-kev,kev,in-the-wild,Zoho,Desktop Centr |
| 5599 | Ivanti Endpoint Manager Cloud Service Appliance (EPM CSA) Code Injection Vulnerability | kev | markdown |  | CVE-2021-44529,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Ivanti, |
| 5600 | Apache Log4j2 Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2021-45046,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Apache, |
| 5601 | D-Link Multiple Routers Remote Code Execution Vulnerability | kev | markdown |  | CVE-2021-45382,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,Multiple Ro |
| 5602 | Palo Alto Networks PAN-OS Reflected Amplification Denial-of-Service Vulnerability | kev | markdown |  | CVE-2022-0028,cve-arsenal,cisa-kev,kev,in-the-wild,Palo Alto Networks, |
| 5603 | Linux Kernel Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2022-0185,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 5604 | Linux Kernel Improper Authentication Vulnerability | kev | markdown |  | CVE-2022-0492,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 5605 | Debian-specific Redis Server Lua Sandbox Escape Vulnerability | kev | markdown |  | CVE-2022-0543,cve-arsenal,cisa-kev,kev,in-the-wild,Redis,Debian-specif |
| 5606 | Google Chromium Animation Use-After-Free Vulnerability | kev | markdown |  | CVE-2022-0609,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium Ani |
| 5607 | Linux Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-0847,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 5608 | Linux Kernel Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2022-0995,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 5609 | Sophos Firewall Authentication Bypass Vulnerability | kev | markdown |  | CVE-2022-1040,cve-arsenal,cisa-kev,kev,in-the-wild,Sophos,Firewall |
| 5610 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2022-1096,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5611 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2022-1364,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5612 | F5 BIG-IP Missing Authentication Vulnerability | kev | markdown |  | CVE-2022-1388,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,F5,BIG-I |
| 5613 | Cisco Small Business RV Series Routers Stack-based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2022-20699,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Small Busine |
| 5614 | Cisco Small Business RV Series Routers Stack-based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2022-20700,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Small Busine |
| 5615 | Cisco Small Business RV Series Routers Stack-based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2022-20701,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Small Busine |
| 5616 | Cisco Small Business RV Series Routers Stack-based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2022-20703,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Small Busine |
| 5617 | Cisco Small Business RV Series Routers Stack-based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2022-20708,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Small Busine |
| 5618 | Cisco SD-WAN Path Traversal Vulnerability | kev | markdown |  | CVE-2022-20775,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,SD-WAN |
| 5619 | Cisco IOS XR Open Port Vulnerability | kev | markdown |  | CVE-2022-20821,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS XR |
| 5620 | Oracle ADF Faces Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2022-21445,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,ADF Faces |
| 5621 | Oracle E-Business Suite Unspecified Vulnerability | kev | markdown |  | CVE-2022-21587,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Oracle, |
| 5622 | Microsoft Win32k Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-21882,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5623 | Microsoft Windows User Profile Service Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-21919,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5624 | Microsoft Windows Runtime Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-21971,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5625 | Microsoft Windows Print Spooler Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-21999,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5626 | Microsoft Windows Client Server Runtime Subsystem (CSRSS) Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-22047,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5627 | Qualcomm Multiple Chipsets Use-After-Free Vulnerability | kev | markdown |  | CVE-2022-22071,cve-arsenal,cisa-kev,kev,in-the-wild,Qualcomm,Multiple  |
| 5628 | Samsung Mobile Devices Use-After-Free Vulnerability | kev | markdown |  | CVE-2022-22265,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,Mobile Dev |
| 5629 | SAP Multiple Products HTTP Request Smuggling Vulnerability | kev | markdown |  | CVE-2022-22536,cve-arsenal,cisa-kev,kev,in-the-wild,SAP,Multiple Produ |
| 5630 | Apple Memory Corruption Vulnerability | kev | markdown |  | CVE-2022-22587,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS and macO |
| 5631 | Apple iOS, iPadOS, and macOS Webkit Use-After-Free Vulnerability | kev | markdown |  | CVE-2022-22620,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS, |
| 5632 | Apple macOS Out-of-Bounds Read Vulnerability | kev | markdown |  | CVE-2022-22674,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,macOS |
| 5633 | Apple macOS Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2022-22675,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,macOS |
| 5634 | Arm Mali GPU Kernel Driver Unspecified Vulnerability | kev | markdown |  | CVE-2022-22706,cve-arsenal,cisa-kev,kev,in-the-wild,Arm,Mali Graphics  |
| 5635 | Microsoft Windows Print Spooler Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-22718,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5636 | WebRTC Heap Buffer Overflow Vulnerability | kev | markdown |  | CVE-2022-2294,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,WebRTC |
| 5637 | VMware Spring Cloud Gateway Code Injection Vulnerability | kev | markdown |  | CVE-2022-22947,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,Spring Clou |
| 5638 | VMware vCenter Server Incorrect Default File Permissions Vulnerability | kev | markdown |  | CVE-2022-22948,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,vCenter Ser |
| 5639 | VMware Workspace ONE Access and Identity Manager Server-Side Template Injection Vulnerability | kev | markdown |  | CVE-2022-22954,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,VMware, |
| 5640 | VMware Multiple Products Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-22960,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,Multiple Pr |
| 5641 | VMware Tanzu Spring Cloud Function Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-22963,cve-arsenal,cisa-kev,kev,in-the-wild,VMware Tanzu,Sprin |
| 5642 | Spring Framework JDK 9+ Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-22965,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,Spring Fram |
| 5643 | Zabbix Frontend Authentication Bypass Vulnerability | kev | markdown |  | CVE-2022-23131,cve-arsenal,cisa-kev,kev,in-the-wild,Zabbix,Frontend |
| 5644 | Zabbix Frontend Improper Access Control Vulnerability | kev | markdown |  | CVE-2022-23134,cve-arsenal,cisa-kev,kev,in-the-wild,Zabbix,Frontend |
| 5645 | WatchGuard Firebox and XTM Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-23176,cve-arsenal,cisa-kev,kev,in-the-wild,WatchGuard,Firebox |
| 5646 | NUUO NVRmini2 Devices Missing Authentication Vulnerability | kev | markdown |  | CVE-2022-23227,cve-arsenal,cisa-kev,kev,in-the-wild,NUUO,NVRmini2 Devi |
| 5647 | Dante Discovery Process Control Vulnerability | kev | markdown |  | CVE-2022-23748,cve-arsenal,cisa-kev,kev,in-the-wild,Audinate,Dante Dis |
| 5648 | Adobe Commerce and Magento Open Source Improper Input Validation Vulnerability | kev | markdown |  | CVE-2022-24086,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Commerce and |
| 5649 | Apache APISIX Authentication Bypass Vulnerability | kev | markdown |  | CVE-2022-24112,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,APISIX |
| 5650 | Microsoft Windows CLFS Driver Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-24521,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5651 | Synacor Zimbra Collaborate Suite (ZCS) Cross-Site Scripting Vulnerability | kev | markdown |  | CVE-2022-24682,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Synacor |
| 5652 | Apache CouchDB Insecure Default Initialization of Resource Vulnerability | kev | markdown |  | CVE-2022-24706,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,CouchDB |
| 5653 | OSGeo GeoServer JAI-EXT Code Injection Vulnerability | kev | markdown |  | CVE-2022-24816,cve-arsenal,cisa-kev,kev,in-the-wild,OSGeo,JAI-EXT |
| 5654 | TerraMaster OS Remote Command Execution Vulnerability | kev | markdown |  | CVE-2022-24990,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,TerraMa |
| 5655 | Linux Kernel Use-After-Free Vulnerability | kev | markdown |  | CVE-2022-2586,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 5656 | Atlassian Confluence Server and Data Center Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-26134,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Atlassi |
| 5657 | Atlassian Questions For Confluence App Hard-coded Credentials Vulnerability | kev | markdown |  | CVE-2022-26138,cve-arsenal,cisa-kev,kev,in-the-wild,Atlassian,Confluen |
| 5658 | MiCollab, MiVoice Business Express Access Control Vulnerability | kev | markdown |  | CVE-2022-26143,cve-arsenal,cisa-kev,kev,in-the-wild,Mitel,MiCollab, Mi |
| 5659 | D-Link DIR-820L Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-26258,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DIR-820L |
| 5660 | WatchGuard Firebox and XTM Appliances Arbitrary Code Execution | kev | markdown |  | CVE-2022-26318,cve-arsenal,cisa-kev,kev,in-the-wild,WatchGuard,Firebox |
| 5661 | dotCMS Unrestricted Upload of File Vulnerability | kev | markdown |  | CVE-2022-26352,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,dotCMS |
| 5662 | Mozilla Firefox Use-After-Free Vulnerability | kev | markdown |  | CVE-2022-26485,cve-arsenal,cisa-kev,kev,in-the-wild,Mozilla,Firefox |
| 5663 | Mozilla Firefox Use-After-Free Vulnerability | kev | markdown |  | CVE-2022-26486,cve-arsenal,cisa-kev,kev,in-the-wild,Mozilla,Firefox |
| 5664 | Veeam Backup & Replication Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-26500,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Veeam,B |
| 5665 | Veeam Backup & Replication Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-26501,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Veeam,B |
| 5666 | Trend Micro Apex Central Arbitrary File Upload Vulnerability | kev | markdown |  | CVE-2022-26871,cve-arsenal,cisa-kev,kev,in-the-wild,Trend Micro,Apex C |
| 5667 | Microsoft Windows User Profile Service Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-26904,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5668 | Microsoft Active Directory Domain Services Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-26923,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Active D |
| 5669 | Microsoft Windows LSA Spoofing Vulnerability | kev | markdown |  | CVE-2022-26925,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5670 | Citrix Application Delivery Controller (ADC) and Gateway Authentication Bypass Vulnerability | kev | markdown |  | CVE-2022-27518,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,Application |
| 5671 | QNAP Photo Station Externally Controlled Reference Vulnerability | kev | markdown |  | CVE-2022-27593,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,QNAP,Ph |
| 5672 | Synacor Zimbra Collaboration Suite (ZCS) Command Injection Vulnerability | kev | markdown |  | CVE-2022-27924,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Synacor |
| 5673 | Synacor Zimbra Collaboration Suite (ZCS) Arbitrary File Upload Vulnerability | kev | markdown |  | CVE-2022-27925,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Synacor |
| 5674 | Synacor Zimbra Collaboration Suite (ZCS) Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2022-27926,cve-arsenal,cisa-kev,kev,in-the-wild,Synacor,Zimbra Col |
| 5675 | Google Chromium Intents Insufficient Input Validation Vulnerability | kev | markdown |  | CVE-2022-2856,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium Int |
| 5676 | Zoho ManageEngine ADSelfService Plus Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-28810,cve-arsenal,cisa-kev,kev,in-the-wild,Zoho,ManageEngine |
| 5677 | SolarView Compact Command Injection Vulnerability | kev | markdown |  | CVE-2022-29303,cve-arsenal,cisa-kev,kev,in-the-wild,SolarView,Compact |
| 5678 | WSO2 Multiple Products Unrestrictive Upload of File Vulnerability | kev | markdown |  | CVE-2022-29464,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,WSO2,Mu |
| 5679 | Mitel MiVoice Connect Data Validation Vulnerability | kev | markdown |  | CVE-2022-29499,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Mitel,M |
| 5680 | Microsoft Windows Support Diagnostic Tool (MSDT) Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-30190,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5681 | RARLAB UnRAR Directory Traversal Vulnerability | kev | markdown |  | CVE-2022-30333,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,RARLAB, |
| 5682 | Google Chromium Network Service Use-After-Free Vulnerability | kev | markdown |  | CVE-2022-3038,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium Net |
| 5683 | Zyxel Multiple Firewalls OS Command Injection Vulnerability | kev | markdown |  | CVE-2022-30525,cve-arsenal,cisa-kev,kev,in-the-wild,Zyxel,Multiple Fir |
| 5684 | Google Chromium Mojo Insufficient Data Validation Vulnerability | kev | markdown |  | CVE-2022-3075,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium Moj |
| 5685 | Netwrix Auditor Insecure Object Deserialization Vulnerability | kev | markdown |  | CVE-2022-31199,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Netwrix |
| 5686 | Sophos Firewall Code Injection Vulnerability | kev | markdown |  | CVE-2022-3236,cve-arsenal,cisa-kev,kev,in-the-wild,Sophos,Firewall |
| 5687 | Apple iOS and macOS Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2022-32893,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS and macO |
| 5688 | Apple iOS and macOS Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2022-32894,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS and macO |
| 5689 | Apple iOS, iPadOS, and macOS Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-32917,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS, |
| 5690 | Apache Spark Command Injection Vulnerability | kev | markdown |  | CVE-2022-33891,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Spark |
| 5691 | Microsoft Windows Support Diagnostic Tool (MSDT) Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-34713,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5692 | Zoho ManageEngine Multiple Products Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-35405,cve-arsenal,cisa-kev,kev,in-the-wild,Zoho,ManageEngine |
| 5693 | Teclib GLPI Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-35914,cve-arsenal,cisa-kev,kev,in-the-wild,Teclib,GLPI |
| 5694 | ZK Framework AuUploader Unspecified Vulnerability | kev | markdown |  | CVE-2022-36537,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,ZK Fram |
| 5695 | Atlassian Bitbucket Server and Data Center Command Injection Vulnerability | kev | markdown |  | CVE-2022-36804,cve-arsenal,cisa-kev,kev,in-the-wild,Atlassian,Bitbucke |
| 5696 | Synacor Zimbra Collaboration Suite (ZCS) Authentication Bypass Vulnerability | kev | markdown |  | CVE-2022-37042,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Synacor |
| 5697 | D-Link Routers Buffer Overflow Vulnerability | kev | markdown |  | CVE-2022-37055,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,Routers |
| 5698 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2022-3723,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5699 | Microsoft Windows Common Log File System (CLFS) Driver Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-37969,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5700 | Microsoft Windows Print Spooler Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-38028,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5701 | Arm Mali GPU Kernel Driver Use-After-Free Vulnerability | kev | markdown |  | CVE-2022-38181,cve-arsenal,cisa-kev,kev,in-the-wild,Arm,Mali Graphics  |
| 5702 | Fortra Cobalt Strike Teamserver Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2022-39197,cve-arsenal,cisa-kev,kev,in-the-wild,Fortra,Cobalt Stri |
| 5703 | Trend Micro Apex One and Apex One as a Service Improper Validation Vulnerability | kev | markdown |  | CVE-2022-40139,cve-arsenal,cisa-kev,kev,in-the-wild,Trend Micro,Apex O |
| 5704 | Fortinet Multiple Products Authentication Bypass Vulnerability | kev | markdown |  | CVE-2022-40684,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortine |
| 5705 | Mitel MiVoice Connect Command Injection Vulnerability | kev | markdown |  | CVE-2022-40765,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Mitel,M |
| 5706 | D-Link DNR-322L Download of Code Without Integrity Check Vulnerability | kev | markdown |  | CVE-2022-40799,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DNR-322L |
| 5707 | Microsoft Windows COM+ Event System Service Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-41033,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows  |
| 5708 | Microsoft Exchange Server Server-Side Request Forgery Vulnerability | kev | markdown |  | CVE-2022-41040,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5709 | Microsoft Windows Mark of the Web (MOTW) Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2022-41049,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5710 | Microsoft Windows Print Spooler Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-41073,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5711 | Microsoft Exchange Server Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-41080,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5712 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-41082,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5713 | Microsoft Windows Mark of the Web (MOTW) Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2022-41091,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5714 | Microsoft Windows CNG Key Isolation Service Privilege Escalation Vulnerability | kev | markdown |  | CVE-2022-41125,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5715 | Microsoft Windows Scripting Languages Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-41128,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5716 | Mitel MiVoice Connect Code Injection Vulnerability | kev | markdown |  | CVE-2022-41223,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Mitel,M |
| 5717 | Fortinet FortiOS Path Traversal Vulnerability | kev | markdown |  | CVE-2022-41328,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,FortiOS |
| 5718 | Google Chromium GPU Heap Buffer Overflow Vulnerability | kev | markdown |  | CVE-2022-4135,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium GPU |
| 5719 | Synacor Zimbra Collaboration Suite (ZCS) Arbitrary File Upload Vulnerability | kev | markdown |  | CVE-2022-41352,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Synacor |
| 5720 | Fortinet FortiOS Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2022-42475,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortine |
| 5721 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2022-4262,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5722 | Apple iOS and iPadOS Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2022-42827,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS and iPad |
| 5723 | Apple iOS Type Confusion Vulnerability | kev | markdown |  | CVE-2022-42856,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS |
| 5724 | Fortra Cobalt Strike User Interface Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-42948,cve-arsenal,cisa-kev,kev,in-the-wild,Fortra,Cobalt Stri |
| 5725 | Hitachi Vantara Pentaho BA Server Special Element Injection Vulnerability | kev | markdown |  | CVE-2022-43769,cve-arsenal,cisa-kev,kev,in-the-wild,Hitachi Vantara,Pe |
| 5726 | Hitachi Vantara Pentaho BA Server Authorization Bypass Vulnerability | kev | markdown |  | CVE-2022-43939,cve-arsenal,cisa-kev,kev,in-the-wild,Hitachi Vantara,Pe |
| 5727 | Microsoft Defender SmartScreen Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2022-44698,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5728 | CWP Control Web Panel OS Command Injection Vulnerability | kev | markdown |  | CVE-2022-44877,cve-arsenal,cisa-kev,kev,in-the-wild,CWP,Control Web Pa |
| 5729 | Cacti Command Injection Vulnerability | kev | markdown |  | CVE-2022-46169,cve-arsenal,cisa-kev,kev,in-the-wild,Cacti |
| 5730 | Zoho ManageEngine Multiple Products Remote Code Execution Vulnerability | kev | markdown |  | CVE-2022-47966,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Zoho,Ma |
| 5731 | IBM Aspera Faspex Code Execution Vulnerability | kev | markdown |  | CVE-2022-47986,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,IBM,Asp |
| 5732 | Apple Multiple Products Unspecified Vulnerability | kev | markdown |  | CVE-2022-48503,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5733 | Apple Multiple Products Memory Corruption Vulnerability | kev | markdown |  | CVE-2022-48618,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5734 | Linux Kernel Use-After-Free Vulnerability | kev | markdown |  | CVE-2023-0266,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 5735 | Linux Kernel Improper Ownership Management Vulnerability | kev | markdown |  | CVE-2023-0386,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 5736 | Fortra GoAnywhere MFT Remote Code Execution Vulnerability | kev | markdown |  | CVE-2023-0669,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortra,G |
| 5737 | TP-Link Archer AX-21 Command Injection Vulnerability | kev | markdown |  | CVE-2023-1389,cve-arsenal,cisa-kev,kev,in-the-wild,TP-Link,Archer AX21 |
| 5738 | Sophos Web Appliance Command Injection Vulnerability | kev | markdown |  | CVE-2023-1671,cve-arsenal,cisa-kev,kev,in-the-wild,Sophos,Web Applianc |
| 5739 | Cisco IOS and IOS XE Group Encrypted Transport VPN Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2023-20109,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS  |
| 5740 | Cisco Small Business RV Series Routers Command Injection Vulnerability | kev | markdown |  | CVE-2023-20118,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Small Busine |
| 5741 | Cisco IOS XE Web UI Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-20198,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS XE Web U |
| 5742 | Cisco Adaptive Security Appliance and Firepower Threat Defense Unauthorized Access Vulnerability | kev | markdown |  | CVE-2023-20269,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Cisco,A |
| 5743 | Cisco IOS XE Web UI Command Injection Vulnerability | kev | markdown |  | CVE-2023-20273,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Cisco IOS XE |
| 5744 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2023-2033,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5745 | VMware Tools Authentication Bypass Vulnerability | kev | markdown |  | CVE-2023-20867,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,Tools |
| 5746 | Vmware Aria Operations for Networks Command Injection Vulnerability | kev | markdown |  | CVE-2023-20887,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,Aria Operat |
| 5747 | Android Framework Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-20963,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Framework |
| 5748 | Android Pixel Information Disclosure Vulnerability | kev | markdown |  | CVE-2023-21237,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Pixel |
| 5749 | Google Chrome Skia Integer Overflow Vulnerability | kev | markdown |  | CVE-2023-2136,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium Ski |
| 5750 | Samsung Mobile Devices Insertion of Sensitive Information Into Log File Vulnerability | kev | markdown |  | CVE-2023-21492,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,Mobile Dev |
| 5751 | Microsoft Exchange Server Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2023-21529,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5752 | Adobe Acrobat and Reader Use-After-Free Vulnerability | kev | markdown |  | CVE-2023-21608,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat and  |
| 5753 | Microsoft Windows Advanced Local Procedure Call (ALPC) Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-21674,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5754 | Microsoft Office Publisher Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2023-21715,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 5755 | Microsoft Windows Graphic Component Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-21823,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5756 | Oracle WebLogic Server Unspecified Vulnerability | kev | markdown |  | CVE-2023-21839,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,WebLogic Se |
| 5757 | Atlassian Confluence Data Center and Server Broken Access Control Vulnerability | kev | markdown |  | CVE-2023-22515,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Atlassi |
| 5758 | Atlassian Confluence Data Center and Server Improper Authorization Vulnerability | kev | markdown |  | CVE-2023-22518,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Atlassi |
| 5759 | Atlassian Confluence Data Center and Server Template Injection Vulnerability | kev | markdown |  | CVE-2023-22527,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Atlassi |
| 5760 | Strapi Cleartext Storage of Sensitive Information Vulnerability | kev | markdown |  | CVE-2023-22894,cve-arsenal,cisa-kev,kev,in-the-wild,Strapi |
| 5761 | Multiple SugarCRM Products Remote Code Execution Vulnerability | kev | markdown |  | CVE-2023-22952,cve-arsenal,cisa-kev,kev,in-the-wild,SugarCRM,Multiple  |
| 5762 | Microsoft Windows Common Log File System (CLFS) Driver Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-23376,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5763 | Microsoft Office Outlook Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-23397,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 5764 | Apple Multiple Products WebKit Type Confusion Vulnerability | kev | markdown |  | CVE-2023-23529,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5765 | Joomla! Improper Access Control Vulnerability | kev | markdown |  | CVE-2023-23752,cve-arsenal,cisa-kev,kev,in-the-wild,Joomla! |
| 5766 | Citrix Content Collaboration ShareFile Improper Access Control Vulnerability | kev | markdown |  | CVE-2023-24489,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,Content Col |
| 5767 | Microsoft Windows SmartScreen Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2023-24880,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5768 | Microsoft SharePoint Server Code Injection Vulnerability | kev | markdown |  | CVE-2023-24955,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5769 | D-Link DIR-820 Router OS Command Injection Vulnerability | kev | markdown |  | CVE-2023-25280,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DIR-820 Rou |
| 5770 | PaperCut NG/MF Cross-Site Request Forgery (CSRF) Vulnerability | kev | markdown |  | CVE-2023-2533,cve-arsenal,cisa-kev,kev,in-the-wild,PaperCut,NG/MF |
| 5771 | Multiple Ruckus Wireless Products CSRF and RCE Vulnerability | kev | markdown |  | CVE-2023-25717,cve-arsenal,cisa-kev,kev,in-the-wild,Ruckus Wireless,Mu |
| 5772 | Danfoss AK-SM 800A Series Path Traversal Vulnerability | kev | markdown |  | CVE-2023-25914,cve-arsenal,cisa-kev,kev,in-the-wild,Danfoss,AK-SM 800A |
| 5773 | Danfoss AK-SM 800A Series Improper Input Validation Vulnerability | kev | markdown |  | CVE-2023-25915,cve-arsenal,cisa-kev,kev,in-the-wild,Danfoss,AK-SM 800A |
| 5774 | Arm Mali GPU Kernel Driver Information Disclosure Vulnerability | kev | markdown |  | CVE-2023-26083,cve-arsenal,cisa-kev,kev,in-the-wild,Arm,Mali Graphics  |
| 5775 | Adobe ColdFusion Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2023-26359,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,ColdFusion |
| 5776 | Adobe ColdFusion Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2023-26360,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,ColdFusion |
| 5777 | Adobe Acrobat and Reader Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2023-26369,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat and  |
| 5778 | PaperCut MF/NG Improper Access Control Vulnerability | kev | markdown |  | CVE-2023-27350,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,PaperCu |
| 5779 | PaperCut NG/MF Improper Authentication Vulnerability | kev | markdown |  | CVE-2023-27351,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,PaperCu |
| 5780 | Apache Superset Insecure Default Initialization of Resource Vulnerability | kev | markdown |  | CVE-2023-27524,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Superset |
| 5781 | Veeam Backup & Replication Cloud Connect Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2023-27532,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Veeam,B |
| 5782 | Zyxel Multiple NAS Devices Command Injection Vulnerability | kev | markdown |  | CVE-2023-27992,cve-arsenal,cisa-kev,kev,in-the-wild,Zyxel,Multiple Net |
| 5783 | Fortinet FortiOS and FortiProxy SSL-VPN Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2023-27997,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortine |
| 5784 | Apple Multiple Products WebKit Out-of-Bounds Read Vulnerability | kev | markdown |  | CVE-2023-28204,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5785 | Apple Multiple Products WebKit Use-After-Free Vulnerability | kev | markdown |  | CVE-2023-28205,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5786 | Apple iOS, iPadOS, and macOS IOSurfaceAccelerator Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2023-28206,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS, |
| 5787 | Microsoft Windows CNG Key Isolation Service Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-28229,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows  |
| 5788 | Microsoft Windows Common Log File System (CLFS) Driver Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-28252,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5789 | MinIO Information Disclosure Vulnerability | kev | markdown |  | CVE-2023-28432,cve-arsenal,cisa-kev,kev,in-the-wild,MinIO |
| 5790 | MinIO Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2023-28434,cve-arsenal,cisa-kev,kev,in-the-wild,MinIO |
| 5791 | Array Networks AG and vxAG ArrayOS Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2023-28461,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Array N |
| 5792 | Barracuda Networks ESG Appliance Improper Input Validation Vulnerability | kev | markdown |  | CVE-2023-2868,cve-arsenal,cisa-kev,kev,in-the-wild,Barracuda Networks, |
| 5793 | Zyxel Multiple Firewalls OS Command Injection Vulnerability | kev | markdown |  | CVE-2023-28771,cve-arsenal,cisa-kev,kev,in-the-wild,Zyxel,Multiple Fir |
| 5794 | Adobe ColdFusion Improper Access Control Vulnerability | kev | markdown |  | CVE-2023-29298,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,ColdFusion |
| 5795 | Adobe ColdFusion Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2023-29300,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Adobe,C |
| 5796 | Microsoft Win32K Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-29336,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Win32k |
| 5797 | Microsoft SharePoint Server Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-29357,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5798 | Microsoft Streaming Service Untrusted Pointer Dereference Vulnerability | kev | markdown |  | CVE-2023-29360,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Streamin |
| 5799 | Novi Survey Insecure Deserialization Vulnerability | kev | markdown |  | CVE-2023-29492,cve-arsenal,cisa-kev,kev,in-the-wild,Novi Survey |
| 5800 | Service Location Protocol (SLP) Denial-of-Service Vulnerability | kev | markdown |  | CVE-2023-29552,cve-arsenal,cisa-kev,kev,in-the-wild,IETF,Service Locat |
| 5801 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2023-3079,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5802 | Microsoft Windows MSHTML Platform Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-32046,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5803 | Microsoft Windows Defender SmartScreen Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2023-32049,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5804 | Ignite Realtime Openfire Path Traversal Vulnerability | kev | markdown |  | CVE-2023-32315,cve-arsenal,cisa-kev,kev,in-the-wild,Ignite Realtime,Op |
| 5805 | Apple Multiple Products WebKit Use-After-Free Vulnerability | kev | markdown |  | CVE-2023-32373,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5806 | Apple Multiple Products WebKit Sandbox Escape Vulnerability | kev | markdown |  | CVE-2023-32409,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5807 | Apple Multiple Products Integer Overflow Vulnerability | kev | markdown |  | CVE-2023-32434,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5808 | Apple Multiple Products WebKit Memory Corruption Vulnerability | kev | markdown |  | CVE-2023-32435,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5809 | Apple Multiple Products WebKit Type Confusion Vulnerability | kev | markdown |  | CVE-2023-32439,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5810 | Zyxel Multiple Firewalls Buffer Overflow Vulnerability | kev | markdown |  | CVE-2023-33009,cve-arsenal,cisa-kev,kev,in-the-wild,Zyxel,Multiple Fir |
| 5811 | Zyxel Multiple Firewalls Buffer Overflow Vulnerability | kev | markdown |  | CVE-2023-33010,cve-arsenal,cisa-kev,kev,in-the-wild,Zyxel,Multiple Fir |
| 5812 | Qualcomm Multiple Chipsets Use-After-Free Vulnerability | kev | markdown |  | CVE-2023-33063,cve-arsenal,cisa-kev,kev,in-the-wild,Qualcomm,Multiple  |
| 5813 | Qualcomm Multiple Chipsets Use of Out-of-Range Pointer Offset Vulnerability | kev | markdown |  | CVE-2023-33106,cve-arsenal,cisa-kev,kev,in-the-wild,Qualcomm,Multiple  |
| 5814 | Qualcomm Multiple Chipsets Integer Overflow Vulnerability | kev | markdown |  | CVE-2023-33107,cve-arsenal,cisa-kev,kev,in-the-wild,Qualcomm,Multiple  |
| 5815 | Apache RocketMQ Command Execution Vulnerability | kev | markdown |  | CVE-2023-33246,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,RocketMQ |
| 5816 | TP-Link Multiple Routers Command Injection Vulnerability | kev | markdown |  | CVE-2023-33538,cve-arsenal,cisa-kev,kev,in-the-wild,TP-Link,Multiple R |
| 5817 | VMware vCenter Server Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2023-34048,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,vCenter Ser |
| 5818 | Synacor Zimbra Collaboration Suite (ZCS) Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2023-34192,cve-arsenal,cisa-kev,kev,in-the-wild,Synacor,Zimbra Col |
| 5819 | Progress MOVEit Transfer SQL Injection Vulnerability | kev | markdown |  | CVE-2023-34362,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Progres |
| 5820 | Ivanti Endpoint Manager Mobile Authentication Bypass Vulnerability | kev | markdown |  | CVE-2023-35078,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Ivanti, |
| 5821 | Ivanti Endpoint Manager Mobile (EPMM) Path Traversal Vulnerability | kev | markdown |  | CVE-2023-35081,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Endpoint Ma |
| 5822 | Ivanti Endpoint Manager Mobile (EPMM) and MobileIron Core Authentication Bypass Vulnerability | kev | markdown |  | CVE-2023-35082,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Ivanti, |
| 5823 | Citrix NetScaler ADC and NetScaler Gateway Code Injection Vulnerability | kev | markdown |  | CVE-2023-3519,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Citrix,N |
| 5824 | Microsoft Outlook Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2023-35311,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Outlook |
| 5825 | Android Framework Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-35674,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Framework |
| 5826 | Microsoft Windows SmartScreen Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2023-36025,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5827 | Microsoft Windows Desktop Window Manager (DWM) Core Library Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-36033,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5828 | Microsoft Windows Cloud Files Mini Filter Driver Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-36036,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5829 | Microsoft Windows Out-of-Bounds Read Vulnerability | kev | markdown |  | CVE-2023-36424,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5830 | Microsoft WordPad Information Disclosure Vulnerability | kev | markdown |  | CVE-2023-36563,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,WordPad |
| 5831 | Microsoft Windows Mark of the Web (MOTW) Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2023-36584,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5832 | Microsoft Word Information Disclosure Vulnerability | kev | markdown |  | CVE-2023-36761,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Word |
| 5833 | Microsoft Streaming Service Proxy Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-36802,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Streamin |
| 5834 | Juniper Junos OS EX Series PHP External Variable Modification Vulnerability | kev | markdown |  | CVE-2023-36844,cve-arsenal,cisa-kev,kev,in-the-wild,Juniper,Junos OS |
| 5835 | Juniper Junos OS EX Series and SRX Series PHP External Variable Modification Vulnerability | kev | markdown |  | CVE-2023-36845,cve-arsenal,cisa-kev,kev,in-the-wild,Juniper,Junos OS |
| 5836 | Juniper Junos OS SRX Series Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2023-36846,cve-arsenal,cisa-kev,kev,in-the-wild,Juniper,Junos OS |
| 5837 | Juniper Junos OS EX Series Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2023-36847,cve-arsenal,cisa-kev,kev,in-the-wild,Juniper,Junos OS |
| 5838 | Juniper Junos OS SRX Series Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2023-36851,cve-arsenal,cisa-kev,kev,in-the-wild,Juniper,Junos OS |
| 5839 | Microsoft Windows Error Reporting Service Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-36874,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5840 | Microsoft Windows Search Remote Code Execution Vulnerability | kev | markdown |  | CVE-2023-36884,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5841 | Apple Multiple Products WebKit Code Execution Vulnerability | kev | markdown |  | CVE-2023-37450,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5842 | Synacor Zimbra Collaboration Suite (ZCS) Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2023-37580,cve-arsenal,cisa-kev,kev,in-the-wild,Synacor,Zimbra Col |
| 5843 | Ivanti Sentry Authentication Bypass Vulnerability | kev | markdown |  | CVE-2023-38035,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Ivanti, |
| 5844 | Microsoft .NET Core and Visual Studio Denial-of-Service Vulnerability | kev | markdown |  | CVE-2023-38180,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,.NET Cor |
| 5845 | Adobe ColdFusion Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2023-38203,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Adobe,C |
| 5846 | Adobe ColdFusion Improper Access Control Vulnerability | kev | markdown |  | CVE-2023-38205,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,ColdFusion |
| 5847 | Apple Multiple Products Kernel Unspecified Vulnerability | kev | markdown |  | CVE-2023-38606,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5848 | RARLAB WinRAR Code Execution Vulnerability | kev | markdown |  | CVE-2023-38831,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,RARLAB, |
| 5849 | ZKTeco BioTime Path Traversal Vulnerability | kev | markdown |  | CVE-2023-38950,cve-arsenal,cisa-kev,kev,in-the-wild,ZKTeco,BioTime |
| 5850 | ASUS RT-AX55 Routers OS Command Injection Vulnerability | kev | markdown |  | CVE-2023-39780,cve-arsenal,cisa-kev,kev,in-the-wild,ASUS,RT-AX55 Route |
| 5851 | Progress WS_FTP Server Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2023-40044,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Progres |
| 5852 | Apple iOS, iPadOS, and watchOS Wallet Code Execution Vulnerability | kev | markdown |  | CVE-2023-41061,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS, |
| 5853 | Apple iOS, iPadOS, and macOS ImageIO Buffer Overflow Vulnerability | kev | markdown |  | CVE-2023-41064,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS, |
| 5854 | Trend Micro Apex One and Worry-Free Business Security Remote Code Execution Vulnerability | kev | markdown |  | CVE-2023-41179,cve-arsenal,cisa-kev,kev,in-the-wild,Trend Micro,Apex O |
| 5855 | Qlik Sense HTTP Tunneling Vulnerability | kev | markdown |  | CVE-2023-41265,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Qlik,Se |
| 5856 | Qlik Sense Path Traversal Vulnerability | kev | markdown |  | CVE-2023-41266,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Qlik,Se |
| 5857 | Microsoft Skype for Business Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-41763,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Skype fo |
| 5858 | Apple iOS and iPadOS Use-After-Free Vulnerability | kev | markdown |  | CVE-2023-41974,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS and iPad |
| 5859 | Apple Multiple Products Code Execution Vulnerability | kev | markdown |  | CVE-2023-41990,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5860 | Apple Multiple Products Improper Certificate Validation Vulnerability | kev | markdown |  | CVE-2023-41991,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5861 | Apple Multiple Products Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-41992,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5862 | Apple Multiple Products WebKit Code Execution Vulnerability | kev | markdown |  | CVE-2023-41993,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5863 | Arm Mali GPU Kernel Driver Use-After-Free Vulnerability | kev | markdown |  | CVE-2023-4211,cve-arsenal,cisa-kev,kev,in-the-wild,Arm,Mali GPU Kernel |
| 5864 | JetBrains TeamCity Authentication Bypass Vulnerability | kev | markdown |  | CVE-2023-42793,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,JetBrai |
| 5865 | Apple iOS and iPadOS Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2023-42824,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS and iPad |
| 5866 | Apple Multiple Products WebKit Out-of-Bounds Read Vulnerability | kev | markdown |  | CVE-2023-42916,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5867 | Apple Multiple Products WebKit Memory Corruption Vulnerability | kev | markdown |  | CVE-2023-42917,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5868 | Apple Multiple products Use-After-Free Vulnerability | kev | markdown |  | CVE-2023-43000,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5869 | NextGen Healthcare Mirth Connect Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2023-43208,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,NextGen |
| 5870 | KNX Association KNX Protocol Connection Authorization Option 1 Overly Restrictive Account Lockout Mechanism Vulnerability | kev | markdown |  | CVE-2023-4346,cve-arsenal,cisa-kev,kev,in-the-wild,KNX Association,KNX |
| 5871 | Roundcube Webmail Persistent Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2023-43770,cve-arsenal,cisa-kev,kev,in-the-wild,Roundcube,Webmail |
| 5872 | SonicWall SMA100 Appliances OS Command Injection Vulnerability | kev | markdown |  | CVE-2023-44221,cve-arsenal,cisa-kev,kev,in-the-wild,SonicWall,SMA100 A |
| 5873 | HTTP/2 Rapid Reset Attack Vulnerability | kev | markdown |  | CVE-2023-44487,cve-arsenal,cisa-kev,kev,in-the-wild,IETF,HTTP/2 |
| 5874 | Acronis Cyber Infrastructure (ACI) Insecure Default Password Vulnerability | kev | markdown |  | CVE-2023-45249,cve-arsenal,cisa-kev,kev,in-the-wild,Acronis,Cyber Infr |
| 5875 | North Grid Proself Improper Restriction of XML External Entity (XXE) Reference Vulnerability | kev | markdown |  | CVE-2023-45727,cve-arsenal,cisa-kev,kev,in-the-wild,North Grid,Proself |
| 5876 | Apache ActiveMQ Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2023-46604,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Apache, |
| 5877 | F5 BIG-IP Configuration Utility Authentication Bypass Vulnerability | kev | markdown |  | CVE-2023-46747,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,F5,BIG- |
| 5878 | F5 BIG-IP Configuration Utility SQL Injection Vulnerability | kev | markdown |  | CVE-2023-46748,cve-arsenal,cisa-kev,kev,in-the-wild,F5,BIG-IP Configur |
| 5879 | Ivanti Connect Secure and Policy Secure Authentication Bypass Vulnerability | kev | markdown |  | CVE-2023-46805,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Ivanti, |
| 5880 | SysAid Server Path Traversal Vulnerability | kev | markdown |  | CVE-2023-47246,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SysAid, |
| 5881 | QNAP VioStor NVR OS Command Injection Vulnerability | kev | markdown |  | CVE-2023-47565,cve-arsenal,cisa-kev,kev,in-the-wild,QNAP,VioStor NVR |
| 5882 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2023-4762,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5883 | Qlik Sense HTTP Tunneling Vulnerability | kev | markdown |  | CVE-2023-48365,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Qlik,Se |
| 5884 | Google Chromium WebP Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2023-4863,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium Web |
| 5885 | Fortinet FortiClient EMS SQL Injection Vulnerability | kev | markdown |  | CVE-2023-48788,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortine |
| 5886 | ownCloud graphapi Information Disclosure Vulnerability | kev | markdown |  | CVE-2023-49103,cve-arsenal,cisa-kev,kev,in-the-wild,ownCloud,ownCloud  |
| 5887 | ownCloud Improper Authentication Vulnerability | kev | markdown |  | CVE-2023-49105,cve-arsenal,cisa-kev,kev,in-the-wild,ownCloud |
| 5888 | GNU C Library Buffer Overflow Vulnerability | kev | markdown |  | CVE-2023-4911,cve-arsenal,cisa-kev,kev,in-the-wild,GNU,GNU C Library |
| 5889 | Citrix NetScaler ADC and NetScaler Gateway Buffer Overflow Vulnerability | kev | markdown |  | CVE-2023-4966,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Citrix,N |
| 5890 | FXC AE1021, AE1021PE OS Command Injection Vulnerability | kev | markdown |  | CVE-2023-49897,cve-arsenal,cisa-kev,kev,in-the-wild,FXC,AE1021, AE1021 |
| 5891 | TP-Link TL-WR841N Authentication Bypass by Spoofing Vulnerability | kev | markdown |  | CVE-2023-50224,cve-arsenal,cisa-kev,kev,in-the-wild,TP-Link,TL-WR841N |
| 5892 | Digiever DS-2105 Pro Missing Authorization Vulnerability | kev | markdown |  | CVE-2023-52163,cve-arsenal,cisa-kev,kev,in-the-wild,Digiever,DS-2105 P |
| 5893 | Google Chromium libvpx Heap Buffer Overflow Vulnerability | kev | markdown |  | CVE-2023-5217,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium lib |
| 5894 | Roundcube Webmail Persistent Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2023-5631,cve-arsenal,cisa-kev,kev,in-the-wild,Roundcube,Webmail |
| 5895 | Google Skia Integer Overflow Vulnerability | kev | markdown |  | CVE-2023-6345,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium Ski |
| 5896 | Unitronics Vision PLC and HMI Insecure Default Password Vulnerability | kev | markdown |  | CVE-2023-6448,cve-arsenal,cisa-kev,kev,in-the-wild,Unitronics,Vision P |
| 5897 | Citrix NetScaler ADC and NetScaler Gateway Code Injection Vulnerability | kev | markdown |  | CVE-2023-6548,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,NetScaler AD |
| 5898 | Citrix NetScaler ADC and NetScaler Gateway Buffer Overflow Vulnerability | kev | markdown |  | CVE-2023-6549,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,NetScaler AD |
| 5899 | Google Chromium WebRTC Heap Buffer Overflow Vulnerability | kev | markdown |  | CVE-2023-7024,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium Web |
| 5900 | GitLab Community and Enterprise Editions Improper Access Control Vulnerability | kev | markdown |  | CVE-2023-7028,cve-arsenal,cisa-kev,kev,in-the-wild,GitLab,GitLab CE/EE |
| 5901 | Spreadsheet::ParseExcel Remote Code Execution Vulnerability | kev | markdown |  | CVE-2023-7101,cve-arsenal,cisa-kev,kev,in-the-wild,Spreadsheet::ParseE |
| 5902 | Palo Alto Networks PAN-OS Management Interface Authentication Bypass Vulnerability | kev | markdown |  | CVE-2024-0012,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Palo Alt |
| 5903 | Google Chromium V8 Out-of-Bounds Memory Access Vulnerability | kev | markdown |  | CVE-2024-0519,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 5904 | D-Link DIR-859 Router Path Traversal Vulnerability | kev | markdown |  | CVE-2024-0769,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DIR-859 Rout |
| 5905 | Synology DiskStation Manager (DSM) Improper Encoding or Escaping of Output Vulnerability | kev | markdown |  | CVE-2024-10441,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Synolog |
| 5906 | Linux Kernel Use-After-Free Vulnerability | kev | markdown |  | CVE-2024-1086,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Linux,Ke |
| 5907 | GeoVision Devices OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-11120,cve-arsenal,cisa-kev,kev,in-the-wild,GeoVision,Multiple |
| 5908 | MDaemon Email Server Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2024-11182,cve-arsenal,cisa-kev,kev,in-the-wild,MDaemon,Email Serv |
| 5909 | Zyxel Multiple Firewalls Path Traversal Vulnerability | kev | markdown |  | CVE-2024-11667,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Zyxel,M |
| 5910 | ProjectSend Improper Authentication Vulnerability | kev | markdown |  | CVE-2024-11680,cve-arsenal,cisa-kev,kev,in-the-wild,ProjectSend |
| 5911 | Progress Kemp LoadMaster OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-1212,cve-arsenal,cisa-kev,kev,in-the-wild,Progress,Kemp LoadM |
| 5912 | BeyondTrust Privileged Remote Access (PRA) and Remote Support (RS) Command Injection Vulnerability | kev | markdown |  | CVE-2024-12356,cve-arsenal,cisa-kev,kev,in-the-wild,BeyondTrust,Privil |
| 5913 | BeyondTrust Privileged Remote Access (PRA) and Remote Support (RS) OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-12686,cve-arsenal,cisa-kev,kev,in-the-wild,BeyondTrust,Privil |
| 5914 | DrayTek Vigor Routers OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-12987,cve-arsenal,cisa-kev,kev,in-the-wild,DrayTek,Vigor Rout |
| 5915 | Ivanti Endpoint Manager (EPM) Absolute Path Traversal Vulnerability | kev | markdown |  | CVE-2024-13159,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Endpoint Ma |
| 5916 | Ivanti Endpoint Manager (EPM) Absolute Path Traversal Vulnerability | kev | markdown |  | CVE-2024-13160,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Endpoint Ma |
| 5917 | Ivanti Endpoint Manager (EPM) Absolute Path Traversal Vulnerability | kev | markdown |  | CVE-2024-13161,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Endpoint Ma |
| 5918 | ConnectWise ScreenConnect Path Traversal Vulnerability | kev | markdown |  | CVE-2024-1708,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,ConnectW |
| 5919 | ConnectWise ScreenConnect Authentication Bypass Vulnerability | kev | markdown |  | CVE-2024-1709,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,ConnectW |
| 5920 | Cisco ASA and FTD Denial of Service Vulnerability | kev | markdown |  | CVE-2024-20353,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Adaptive Sec |
| 5921 | Cisco ASA and FTD Privilege Escalation Vulnerability | kev | markdown |  | CVE-2024-20359,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Adaptive Sec |
| 5922 | Cisco NX-OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-20399,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,NX-OS |
| 5923 | Cisco Smart Licensing Utility Static Credential Vulnerability | kev | markdown |  | CVE-2024-20439,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Smart Licens |
| 5924 | Cisco ASA and FTD Denial-of-Service Vulnerability | kev | markdown |  | CVE-2024-20481,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Adaptive Sec |
| 5925 | Adobe ColdFusion Improper Access Control Vulnerability | kev | markdown |  | CVE-2024-20767,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,ColdFusion |
| 5926 | Oracle Agile Product Lifecycle Management (PLM) Deserialization Vulnerability | kev | markdown |  | CVE-2024-20953,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,Agile Produ |
| 5927 | Oracle WebLogic Server Unspecified Vulnerability | kev | markdown |  | CVE-2024-21182,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,WebLogic Se |
| 5928 | Oracle Agile Product Lifecycle Management (PLM) Incorrect Authorization Vulnerability | kev | markdown |  | CVE-2024-21287,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,Agile Produ |
| 5929 | Microsoft Windows Kernel Exposed IOCTL with Insufficient Access Control Vulnerability | kev | markdown |  | CVE-2024-21338,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5930 | Microsoft Windows SmartScreen Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2024-21351,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5931 | Microsoft Exchange Server Privilege Escalation Vulnerability | kev | markdown |  | CVE-2024-21410,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Exchange |
| 5932 | Microsoft Windows Internet Shortcut Files Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2024-21412,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5933 | Microsoft Outlook Improper Input Validation Vulnerability | kev | markdown |  | CVE-2024-21413,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office O |
| 5934 | Fortinet FortiOS Out-of-Bound Write Vulnerability | kev | markdown |  | CVE-2024-21762,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortine |
| 5935 | Ivanti Connect Secure and Policy Secure Command Injection Vulnerability | kev | markdown |  | CVE-2024-21887,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Ivanti, |
| 5936 | Ivanti Connect Secure, Policy Secure, and Neurons Server-Side Request Forgery (SSRF) Vulnerability | kev | markdown |  | CVE-2024-21893,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Ivanti, |
| 5937 | Fortinet Multiple Products Format String Vulnerability | kev | markdown |  | CVE-2024-23113,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,Multiple  |
| 5938 | Apple Multiple Products WebKit Type Confusion Vulnerability | kev | markdown |  | CVE-2024-23222,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5939 | Apple Multiple Products Memory Corruption Vulnerability | kev | markdown |  | CVE-2024-23225,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5940 | Apple Multiple Products Memory Corruption Vulnerability | kev | markdown |  | CVE-2024-23296,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5941 | Rejetto HTTP File Server Improper Neutralization of Special Elements Used in a Template Engine Vulnerability | kev | markdown |  | CVE-2024-23692,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Rejetto |
| 5942 | Jenkins Command Line Interface (CLI) Path Traversal Vulnerability | kev | markdown |  | CVE-2024-23897,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Jenkins |
| 5943 | Check Point Quantum Security Gateways Information Disclosure Vulnerability | kev | markdown |  | CVE-2024-24919,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Check P |
| 5944 | Microsoft Windows Error Reporting Service Improper Privilege Management Vulnerability | kev | markdown |  | CVE-2024-26169,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5945 | JetBrains TeamCity Authentication Bypass Vulnerability | kev | markdown |  | CVE-2024-27198,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,JetBrai |
| 5946 | JetBrains TeamCity Relative Path Traversal Vulnerability | kev | markdown |  | CVE-2024-27199,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,JetBrai |
| 5947 | Apache HugeGraph-Server Improper Access Control Vulnerability | kev | markdown |  | CVE-2024-27348,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,HugeGraph-S |
| 5948 | Synacor Zimbra Collaboration Suite (ZCS) Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2024-27443,cve-arsenal,cisa-kev,kev,in-the-wild,Synacor,Zimbra Col |
| 5949 | Apple Multiple Products Improper Authentication Vulnerability | kev | markdown |  | CVE-2024-27834,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 5950 | SolarWinds Web Help Desk Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2024-28986,cve-arsenal,cisa-kev,kev,in-the-wild,SolarWinds,Web Hel |
| 5951 | SolarWinds Web Help Desk Hardcoded Credential Vulnerability | kev | markdown |  | CVE-2024-28987,cve-arsenal,cisa-kev,kev,in-the-wild,SolarWinds,Web Hel |
| 5952 | SolarWinds Serv-U Path Traversal Vulnerability | kev | markdown |  | CVE-2024-28995,cve-arsenal,cisa-kev,kev,in-the-wild,SolarWinds,Serv-U |
| 5953 | Microsoft .NET Framework Information Disclosure Vulnerability | kev | markdown |  | CVE-2024-29059,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,.NET Fra |
| 5954 | Android Pixel Information Disclosure Vulnerability | kev | markdown |  | CVE-2024-29745,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Pixel |
| 5955 | Android Pixel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2024-29748,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Pixel |
| 5956 | Ivanti Endpoint Manager (EPM) SQL Injection Vulnerability | kev | markdown |  | CVE-2024-29824,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Endpoint Ma |
| 5957 | Microsoft SmartScreen Prompt Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2024-29988,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,SmartScr |
| 5958 | Microsoft Windows MSHTML Platform Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2024-30040,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5959 | Microsoft DWM Core Library Privilege Escalation Vulnerability | kev | markdown |  | CVE-2024-30051,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5960 | Microsoft Windows Kernel TOCTOU Race Condition Vulnerability | kev | markdown |  | CVE-2024-30088,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5961 | Apache OFBiz Path Traversal Vulnerability | kev | markdown |  | CVE-2024-32113,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,OFBiz |
| 5962 | D-Link Multiple NAS Devices Use of Hard-Coded Credentials Vulnerability | kev | markdown |  | CVE-2024-3272,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,Multiple NAS |
| 5963 | D-Link Multiple NAS Devices Command Injection Vulnerability | kev | markdown |  | CVE-2024-3273,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,Multiple NAS |
| 5964 | Android Pixel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2024-32896,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Pixel |
| 5965 | Palo Alto Networks PAN-OS Malicious DNS Packet Vulnerability | kev | markdown |  | CVE-2024-3393,cve-arsenal,cisa-kev,kev,in-the-wild,Palo Alto Networks, |
| 5966 | Palo Alto Networks PAN-OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-3400,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Palo Alt |
| 5967 | Adobe Commerce and Magento Open Source Improper Restriction of XML External Entity Reference (XXE) Vulnerability | kev | markdown |  | CVE-2024-34102,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Commerce and |
| 5968 | Microsoft Windows Kernel-Mode Driver Untrusted Pointer Dereference Vulnerability | kev | markdown |  | CVE-2024-35250,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5969 | OSGeo GeoServer GeoTools Eval Injection Vulnerability | kev | markdown |  | CVE-2024-36401,cve-arsenal,cisa-kev,kev,in-the-wild,OSGeo,GeoServer |
| 5970 | Android Kernel Remote Code Execution Vulnerability | kev | markdown |  | CVE-2024-36971,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Kernel |
| 5971 | Broadcom VMware vCenter Server Out-of-bounds Write Vulnerability | kev | markdown |  | CVE-2024-37079,cve-arsenal,cisa-kev,kev,in-the-wild,Broadcom,VMware vC |
| 5972 | VMware ESXi Authentication Bypass Vulnerability | kev | markdown |  | CVE-2024-37085,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,VMware, |
| 5973 | RoundCube Webmail Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2024-37383,cve-arsenal,cisa-kev,kev,in-the-wild,Roundcube,Webmail |
| 5974 | Microsoft Windows Installer Improper Privilege Management Vulnerability | kev | markdown |  | CVE-2024-38014,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5975 | Microsoft Windows Hyper-V Privilege Escalation Vulnerability | kev | markdown |  | CVE-2024-38080,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5976 | Microsoft SharePoint Deserialization Vulnerability | kev | markdown |  | CVE-2024-38094,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 5977 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | markdown |  | CVE-2024-38106,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5978 | Microsoft Windows Power Dependency Coordinator Privilege Escalation Vulnerability | kev | markdown |  | CVE-2024-38107,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5979 | Microsoft Windows MSHTML Platform Spoofing Vulnerability | kev | markdown |  | CVE-2024-38112,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5980 | Microsoft Windows Scripting Engine Memory Corruption Vulnerability | kev | markdown |  | CVE-2024-38178,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5981 | Microsoft Project Remote Code Execution Vulnerability | kev | markdown |  | CVE-2024-38189,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Project |
| 5982 | Microsoft Windows Ancillary Function Driver for WinSock Privilege Escalation Vulnerability | kev | markdown |  | CVE-2024-38193,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5983 | Microsoft Windows SmartScreen Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2024-38213,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5984 | Microsoft Windows Mark of the Web (MOTW) Protection Mechanism Failure Vulnerability | kev | markdown |  | CVE-2024-38217,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 5985 | Microsoft Publisher Protection Mechanism Failure Vulnerability | kev | markdown |  | CVE-2024-38226,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Publishe |
| 5986 | Apache HTTP Server Improper Escaping of Output Vulnerability | kev | markdown |  | CVE-2024-38475,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,HTTP Server |
| 5987 | VMware vCenter Server Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2024-38812,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,vCenter Ser |
| 5988 | VMware vCenter Server Privilege Escalation Vulnerability | kev | markdown |  | CVE-2024-38813,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,vCenter Ser |
| 5989 | Apache OFBiz Incorrect Authorization Vulnerability | kev | markdown |  | CVE-2024-38856,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,OFBiz |
| 5990 | Versa Director Dangerous File Type Upload Vulnerability | kev | markdown |  | CVE-2024-39717,cve-arsenal,cisa-kev,kev,in-the-wild,Versa,Director |
| 5991 | Twilio Authy Information Disclosure Vulnerability | kev | markdown |  | CVE-2024-39891,cve-arsenal,cisa-kev,kev,in-the-wild,Twilio,Authy |
| 5992 | CrushFTP VFS Sandbox Escape Vulnerability | kev | markdown |  | CVE-2024-4040,cve-arsenal,cisa-kev,kev,in-the-wild,CrushFTP |
| 5993 | Veeam Backup and Replication Deserialization Vulnerability | kev | markdown |  | CVE-2024-40711,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Veeam,B |
| 5994 | SonicWall SonicOS Improper Access Control Vulnerability | kev | markdown |  | CVE-2024-40766,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SonicWa |
| 5995 | Zyxel DSL CPE OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-40890,cve-arsenal,cisa-kev,kev,in-the-wild,Zyxel,DSL CPE Devi |
| 5996 | Zyxel DSL CPE OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-40891,cve-arsenal,cisa-kev,kev,in-the-wild,Zyxel,DSL CPE Devi |
| 5997 | Mitel SIP Phones Argument Injection Vulnerability | kev | markdown |  | CVE-2024-41710,cve-arsenal,cisa-kev,kev,in-the-wild,Mitel,SIP Phones |
| 5998 | Mitel MiCollab Path Traversal Vulnerability | kev | markdown |  | CVE-2024-41713,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Mitel,M |
| 5999 | RoundCube Webmail Cross-Site Scripting Vulnerability | kev | markdown |  | CVE-2024-42009,cve-arsenal,cisa-kev,kev,in-the-wild,Roundcube,Webmail |
| 6000 | Qualcomm Multiple Chipsets Use-After-Free Vulnerability | kev | markdown |  | CVE-2024-43047,cve-arsenal,cisa-kev,kev,in-the-wild,Qualcomm,Multiple  |
| 6001 | Android Framework Privilege Escalation Vulnerability | kev | markdown |  | CVE-2024-43093,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Framework |
| 6002 | Microsoft Windows NTLMv2 Hash Disclosure Spoofing Vulnerability | kev | markdown |  | CVE-2024-43451,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6003 | Microsoft Windows MSHTML Platform Spoofing Vulnerability | kev | markdown |  | CVE-2024-43461,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6004 | Microsoft Configuration Manager SQL Injection Vulnerability | kev | markdown |  | CVE-2024-43468,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Configur |
| 6005 | Microsoft Windows Management Console Remote Code Execution Vulnerability | kev | markdown |  | CVE-2024-43572,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6006 | Microsoft Windows MSHTML Platform Spoofing Vulnerability | kev | markdown |  | CVE-2024-43573,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6007 | Progress Telerik Report Server Authentication Bypass by Spoofing Vulnerability | kev | markdown |  | CVE-2024-4358,cve-arsenal,cisa-kev,kev,in-the-wild,Progress,Telerik Re |
| 6008 | Apple Multiple Products Code Execution Vulnerability | kev | markdown |  | CVE-2024-44308,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 6009 | Apple Multiple Products Cross-Site Scripting (XSS) Vulnerability | kev | markdown |  | CVE-2024-44309,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 6010 | Apache OFBiz Forced Browsing Vulnerability | kev | markdown |  | CVE-2024-45195,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,OFBiz |
| 6011 | Synacor Zimbra Collaboration Suite (ZCS) Command Execution Vulnerability | kev | markdown |  | CVE-2024-45519,cve-arsenal,cisa-kev,kev,in-the-wild,Synacor,Zimbra Col |
| 6012 | PHP-CGI OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-4577,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,PHP Grou |
| 6013 | Arm Mali GPU Kernel Driver Use-After-Free Vulnerability | kev | markdown |  | CVE-2024-4610,cve-arsenal,cisa-kev,kev,in-the-wild,Arm,Mali GPU Kernel |
| 6014 | Google Chromium Visuals Use-After-Free Vulnerability | kev | markdown |  | CVE-2024-4671,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium |
| 6015 | Fortinet FortiManager Missing Authentication Vulnerability | kev | markdown |  | CVE-2024-47575,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,FortiMana |
| 6016 | Google Chromium V8 Out-of-Bounds Memory Write Vulnerability | kev | markdown |  | CVE-2024-4761,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 6017 | NAKIVO Backup and Replication Absolute Path Traversal Vulnerability | kev | markdown |  | CVE-2024-48248,cve-arsenal,cisa-kev,kev,in-the-wild,NAKIVO,Backup and  |
| 6018 | ServiceNow Improper Input Validation Vulnerability | kev | markdown |  | CVE-2024-4879,cve-arsenal,cisa-kev,kev,in-the-wild,ServiceNow,Utah, Va |
| 6019 | Progress WhatsUp Gold Path Traversal Vulnerability | kev | markdown |  | CVE-2024-4885,cve-arsenal,cisa-kev,kev,in-the-wild,Progress,WhatsUp Go |
| 6020 | Microsoft Partner Center Improper Access Control Vulnerability | kev | markdown |  | CVE-2024-49035,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Partner  |
| 6021 | Microsoft Windows Task Scheduler Privilege Escalation Vulnerability | kev | markdown |  | CVE-2024-49039,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 6022 | Microsoft Windows Common Log File System (CLFS) Driver Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2024-49138,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6023 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2024-4947,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 6024 | Justice AV Solutions (JAVS) Viewer Installer Embedded Malicious Code Vulnerability | kev | markdown |  | CVE-2024-4978,cve-arsenal,cisa-kev,kev,in-the-wild,Justice AV Solution |
| 6025 | Linux Kernel Use of Uninitialized Resource Vulnerability | kev | markdown |  | CVE-2024-50302,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 6026 | Aviatrix Controllers OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-50603,cve-arsenal,cisa-kev,kev,in-the-wild,Aviatrix,Controlle |
| 6027 | Cleo Multiple Products Unrestricted File Upload Vulnerability | kev | markdown |  | CVE-2024-50623,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Cleo,Mu |
| 6028 | CyberPanel Incorrect Default Permissions Vulnerability | kev | markdown |  | CVE-2024-51378,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,CyberPe |
| 6029 | CyberPanel Incorrect Default Permissions Vulnerability | kev | markdown |  | CVE-2024-51567,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,CyberPe |
| 6030 | ServiceNow Incomplete List of Disallowed Inputs Vulnerability | kev | markdown |  | CVE-2024-5217,cve-arsenal,cisa-kev,kev,in-the-wild,ServiceNow,Utah, Va |
| 6031 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2024-5274,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 6032 | Linux Kernel Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2024-53104,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 6033 | Linux Kernel Out-of-Bounds Read Vulnerability | kev | markdown |  | CVE-2024-53150,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 6034 | Linux Kernel Out-of-Bounds Access Vulnerability | kev | markdown |  | CVE-2024-53197,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 6035 | SonicWall SonicOS SSLVPN Improper Authentication Vulnerability | kev | markdown |  | CVE-2024-53704,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SonicWa |
| 6036 | AMI MegaRAC SPx Authentication Bypass by Spoofing Vulnerability | kev | markdown |  | CVE-2024-54085,cve-arsenal,cisa-kev,kev,in-the-wild,AMI,MegaRAC SPx |
| 6037 | Mitel MiCollab Path Traversal Vulnerability | kev | markdown |  | CVE-2024-55550,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Mitel,M |
| 6038 | Fortinet FortiOS and FortiProxy Authentication Bypass Vulnerability | kev | markdown |  | CVE-2024-55591,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortine |
| 6039 | Cleo Multiple Products Unauthenticated File Upload Vulnerability | kev | markdown |  | CVE-2024-55956,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Cleo,Mu |
| 6040 | Craft CMS Code Injection Vulnerability | kev | markdown |  | CVE-2024-56145,cve-arsenal,cisa-kev,kev,in-the-wild,Craft CMS |
| 6041 | SimpleHelp Missing Authorization Vulnerability | kev | markdown |  | CVE-2024-57726,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SimpleH |
| 6042 | SimpleHelp Path Traversal Vulnerability | kev | markdown |  | CVE-2024-57727,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SimpleH |
| 6043 | SimpleHelp Path Traversal Vulnerability | kev | markdown |  | CVE-2024-57728,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SimpleH |
| 6044 | Advantive VeraCore Unrestricted File Upload Vulnerability | kev | markdown |  | CVE-2024-57968,cve-arsenal,cisa-kev,kev,in-the-wild,Advantive,VeraCore |
| 6045 | Yiiframework Yii Improper Protection of Alternate Path Vulnerability | kev | markdown |  | CVE-2024-58136,cve-arsenal,cisa-kev,kev,in-the-wild,Yiiframework,Yii |
| 6046 | Palo Alto Networks Expedition Missing Authentication Vulnerability | kev | markdown |  | CVE-2024-5910,cve-arsenal,cisa-kev,kev,in-the-wild,Palo Alto Networks, |
| 6047 | GeoVision Devices OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-6047,cve-arsenal,cisa-kev,kev,in-the-wild,GeoVision,Multiple  |
| 6048 | Progress WhatsUp Gold SQL Injection Vulnerability | kev | markdown |  | CVE-2024-6670,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Progress |
| 6049 | Kingsoft WPS Office Path Traversal Vulnerability | kev | markdown |  | CVE-2024-7262,cve-arsenal,cisa-kev,kev,in-the-wild,Kingsoft,WPS Office |
| 6050 | Samsung MagicINFO 9 Server Path Traversal Vulnerability | kev | markdown |  | CVE-2024-7399,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,MagicINFO 9 |
| 6051 | Ivanti Virtual Traffic Manager Authentication Bypass Vulnerability | kev | markdown |  | CVE-2024-7593,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Virtual Traf |
| 6052 | TeamT5 ThreatSonar Anti-Ransomware Unrestricted Upload of File with Dangerous Type Vulnerability | kev | markdown |  | CVE-2024-7694,cve-arsenal,cisa-kev,kev,in-the-wild,TeamT5,ThreatSonar  |
| 6053 | Google Chromium V8 Inappropriate Implementation Vulnerability | kev | markdown |  | CVE-2024-7965,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 6054 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2024-7971,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 6055 | Citrix Session Recording Improper Privilege Management Vulnerability | kev | markdown |  | CVE-2024-8068,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,Session Reco |
| 6056 | Citrix Session Recording Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2024-8069,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,Session Reco |
| 6057 | Ivanti Cloud Services Appliance OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-8190,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Cloud Servic |
| 6058 | PTZOptics PT30X-SDI/NDI Cameras Authentication Bypass Vulnerability | kev | markdown |  | CVE-2024-8956,cve-arsenal,cisa-kev,kev,in-the-wild,PTZOptics,PT30X-SDI |
| 6059 | PTZOptics PT30X-SDI/NDI Cameras OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-8957,cve-arsenal,cisa-kev,kev,in-the-wild,PTZOptics,PT30X-SDI |
| 6060 | Ivanti Cloud Services Appliance (CSA) Path Traversal Vulnerability | kev | markdown |  | CVE-2024-8963,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Cloud Servic |
| 6061 | Ivanti Cloud Services Appliance (CSA) SQL Injection Vulnerability | kev | markdown |  | CVE-2024-9379,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Cloud Servic |
| 6062 | Ivanti Cloud Services Appliance (CSA) OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-9380,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Cloud Servic |
| 6063 | Palo Alto Networks Expedition OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-9463,cve-arsenal,cisa-kev,kev,in-the-wild,Palo Alto Networks, |
| 6064 | Palo Alto Networks Expedition SQL Injection Vulnerability | kev | markdown |  | CVE-2024-9465,cve-arsenal,cisa-kev,kev,in-the-wild,Palo Alto Networks, |
| 6065 | Palo Alto Networks PAN-OS Management Interface OS Command Injection Vulnerability | kev | markdown |  | CVE-2024-9474,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Palo Alt |
| 6066 | ScienceLogic SL1 Unspecified Vulnerability | kev | markdown |  | CVE-2024-9537,cve-arsenal,cisa-kev,kev,in-the-wild,ScienceLogic,SL1 |
| 6067 | Mozilla Firefox Use-After-Free Vulnerability | kev | markdown |  | CVE-2024-9680,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Mozilla, |
| 6068 | Palo Alto Networks PAN-OS Authentication Bypass Vulnerability | kev | markdown |  | CVE-2025-0108,cve-arsenal,cisa-kev,kev,in-the-wild,Palo Alto Networks, |
| 6069 | Palo Alto Networks PAN-OS File Read Vulnerability | kev | markdown |  | CVE-2025-0111,cve-arsenal,cisa-kev,kev,in-the-wild,Palo Alto Networks, |
| 6070 | Ivanti Connect Secure, Policy Secure, and ZTA Gateways Stack-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-0282,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Ivanti,C |
| 6071 | 7-Zip Mark of the Web Bypass Vulnerability | kev | markdown |  | CVE-2025-0411,cve-arsenal,cisa-kev,kev,in-the-wild,7-Zip |
| 6072 | Trimble Cityworks Deserialization Vulnerability | kev | markdown |  | CVE-2025-0994,cve-arsenal,cisa-kev,kev,in-the-wild,Trimble,Cityworks |
| 6073 | Fortra GoAnywhere MFT Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2025-10035,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortra, |
| 6074 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2025-10585,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 6075 | Gladinet CentreStack and Triofox Files or Directories Accessible to External Parties Vulnerability | kev | markdown |  | CVE-2025-11371,cve-arsenal,cisa-kev,kev,in-the-wild,Gladinet,CentreSta |
| 6076 | React Native Community CLI OS Command Injection Vulnerability | kev | markdown |  | CVE-2025-11953,cve-arsenal,cisa-kev,kev,in-the-wild,React Native Commu |
| 6077 | Gladinet Triofox Improper Access Control Vulnerability | kev | markdown |  | CVE-2025-12480,cve-arsenal,cisa-kev,kev,in-the-wild,Gladinet,Triofox |
| 6078 | Edimax IC-7100 IP Camera OS Command Injection Vulnerability | kev | markdown |  | CVE-2025-1316,cve-arsenal,cisa-kev,kev,in-the-wild,Edimax,IC-7100 IP C |
| 6079 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2025-13223,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 6080 | Google Chromium Out of Bounds Memory Access Vulnerability | kev | markdown |  | CVE-2025-14174,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium |
| 6081 | Gladinet CentreStack and Triofox Hard Coded Cryptographic Vulnerability | kev | markdown |  | CVE-2025-14611,cve-arsenal,cisa-kev,kev,in-the-wild,Gladinet,CentreSta |
| 6082 | WatchGuard Firebox Out of Bounds Write Vulnerability | kev | markdown |  | CVE-2025-14733,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,WatchGu |
| 6083 | MongoDB and MongoDB Server Improper Handling of Length Parameter Inconsistency Vulnerability | kev | markdown |  | CVE-2025-14847,cve-arsenal,cisa-kev,kev,in-the-wild,MongoDB,MongoDB an |
| 6084 | Notepad++ Download of Code Without Integrity Check Vulnerability | kev | markdown |  | CVE-2025-15556,cve-arsenal,cisa-kev,kev,in-the-wild,Notepad++ |
| 6085 | Broadcom Brocade Fabric OS Code Injection Vulnerability | kev | markdown |  | CVE-2025-1976,cve-arsenal,cisa-kev,kev,in-the-wild,Broadcom,Brocade Fa |
| 6086 | Cisco Identity Services Engine Injection Vulnerability | kev | markdown |  | CVE-2025-20281,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Identity Ser |
| 6087 | Cisco Secure Firewall Adaptive Security Appliance (ASA) and Secure Firewall Threat Defense (FTD) Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-20333,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Secure Firew |
| 6088 | Cisco Identity Services Engine Injection Vulnerability | kev | markdown |  | CVE-2025-20337,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Identity Ser |
| 6089 | Cisco IOS and IOS XE Software SNMP Denial of Service and Remote Code Execution Vulnerability | kev | markdown |  | CVE-2025-20352,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,IOS and IOS  |
| 6090 | Cisco Secure Firewall Adaptive Security (ASA) Appliance and Secure Firewall Threat Defense (FTD) Missing Authorization Vulnerability | kev | markdown |  | CVE-2025-20362,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Secure Firew |
| 6091 | Cisco Multiple Products Improper Input Validation Vulnerability | kev | markdown |  | CVE-2025-20393,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Multiple Pro |
| 6092 | Samsung Mobile Devices Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2025-21042,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,Mobile Dev |
| 6093 | Samsung Mobile Devices Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2025-21043,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,Mobile Dev |
| 6094 | Microsoft Windows Hyper-V NT Kernel Integration VSP Heap-based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-21333,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6095 | Microsoft Windows Hyper-V NT Kernel Integration VSP Use-After-Free Vulnerability | kev | markdown |  | CVE-2025-21334,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6096 | Microsoft Windows Hyper-V NT Kernel Integration VSP Use-After-Free Vulnerability | kev | markdown |  | CVE-2025-21335,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6097 | Microsoft Windows Storage Link Following Vulnerability | kev | markdown |  | CVE-2025-21391,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6098 | Microsoft Windows Ancillary Function Driver for WinSock Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-21418,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6099 | Qualcomm Multiple Chipsets Incorrect Authorization Vulnerability | kev | markdown |  | CVE-2025-21479,cve-arsenal,cisa-kev,kev,in-the-wild,Qualcomm,Multiple  |
| 6100 | Qualcomm Multiple Chipsets Incorrect Authorization Vulnerability | kev | markdown |  | CVE-2025-21480,cve-arsenal,cisa-kev,kev,in-the-wild,Qualcomm,Multiple  |
| 6101 | Juniper Junos OS Improper Isolation or Compartmentalization Vulnerability | kev | markdown |  | CVE-2025-21590,cve-arsenal,cisa-kev,kev,in-the-wild,Juniper,Junos OS |
| 6102 | VMware ESXi and Workstation TOCTOU Race Condition Vulnerability | kev | markdown |  | CVE-2025-22224,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,ESXi and Wo |
| 6103 | VMware ESXi Arbitrary Write Vulnerability | kev | markdown |  | CVE-2025-22225,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,VMware, |
| 6104 | VMware ESXi, Workstation, and Fusion Information Disclosure Vulnerability | kev | markdown |  | CVE-2025-22226,cve-arsenal,cisa-kev,kev,in-the-wild,VMware,ESXi, Works |
| 6105 | Ivanti Connect Secure, Policy Secure, and ZTA Gateways Stack-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-22457,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Ivanti, |
| 6106 | SonicWall SMA1000 Appliances Deserialization Vulnerability | kev | markdown |  | CVE-2025-23006,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SonicWa |
| 6107 | Craft CMS Code Injection Vulnerability | kev | markdown |  | CVE-2025-23209,cve-arsenal,cisa-kev,kev,in-the-wild,Craft CMS |
| 6108 | Wazuh Server Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2025-24016,cve-arsenal,cisa-kev,kev,in-the-wild,Wazuh,Wazuh Server |
| 6109 | Microsoft Windows NTLM Hash Disclosure Spoofing Vulnerability | kev | markdown |  | CVE-2025-24054,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6110 | Apple Multiple Products Use-After-Free Vulnerability | kev | markdown |  | CVE-2025-24085,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 6111 | Apple iOS and iPadOS Incorrect Authorization Vulnerability | kev | markdown |  | CVE-2025-24200,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS and iPad |
| 6112 | Apple Multiple Products WebKit Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2025-24201,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 6113 | Fortinet FortiOS and FortiProxy Authentication Bypass Vulnerability | kev | markdown |  | CVE-2025-24472,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Fortine |
| 6114 | Apache Tomcat Path Equivalence Vulnerability | kev | markdown |  | CVE-2025-24813,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Tomcat |
| 6115 | XWiki Platform Eval Injection Vulnerability | kev | markdown |  | CVE-2025-24893,cve-arsenal,cisa-kev,kev,in-the-wild,XWiki,Platform |
| 6116 | Microsoft Windows Win32k Use-After-Free Vulnerability | kev | markdown |  | CVE-2025-24983,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6117 | Microsoft Windows NTFS Information Disclosure Vulnerability | kev | markdown |  | CVE-2025-24984,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6118 | Microsoft Windows Fast FAT File System Driver Integer Overflow Vulnerability | kev | markdown |  | CVE-2025-24985,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6119 | Microsoft Power Pages Improper Access Control Vulnerability | kev | markdown |  | CVE-2025-24989,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Power Pa |
| 6120 | Microsoft Windows Untrusted Pointer Dereference Vulnerability | kev | markdown |  | CVE-2025-24990,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6121 | Microsoft Windows NTFS Out-Of-Bounds Read Vulnerability | kev | markdown |  | CVE-2025-24991,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6122 | Microsoft Windows NTFS Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-24993,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6123 | Advantive VeraCore SQL Injection Vulnerability | kev | markdown |  | CVE-2025-25181,cve-arsenal,cisa-kev,kev,in-the-wild,Advantive,VeraCore |
| 6124 | Fortinet Multiple Products Heap-based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-25249,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,Multiple  |
| 6125 | Fortinet FortiWeb SQL Injection Vulnerability | kev | markdown |  | CVE-2025-25257,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,FortiWeb |
| 6126 | SolarWinds Web Help Desk Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2025-26399,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SolarWi |
| 6127 | Microsoft Windows Management Console (MMC) Improper Neutralization Vulnerability | kev | markdown |  | CVE-2025-26633,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 6128 | Qualcomm Multiple Chipsets Use-After-Free Vulnerability | kev | markdown |  | CVE-2025-27038,cve-arsenal,cisa-kev,kev,in-the-wild,Qualcomm,Multiple  |
| 6129 | FreeType Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2025-27363,cve-arsenal,cisa-kev,kev,in-the-wild,FreeType |
| 6130 | Kentico Xperience CMS Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | markdown |  | CVE-2025-2746,cve-arsenal,cisa-kev,kev,in-the-wild,Kentico,Xperience C |
| 6131 | Kentico Xperience CMS Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | markdown |  | CVE-2025-2747,cve-arsenal,cisa-kev,kev,in-the-wild,Kentico,Xperience C |
| 6132 | Kentico Xperience Path Traversal Vulnerability | kev | markdown |  | CVE-2025-2749,cve-arsenal,cisa-kev,kev,in-the-wild,Kentico,Kentico Xpe |
| 6133 | SysAid On-Prem Improper Restriction of XML External Entity Reference Vulnerability | kev | markdown |  | CVE-2025-2775,cve-arsenal,cisa-kev,kev,in-the-wild,SysAid,SysAid On-Pr |
| 6134 | SysAid On-Prem Improper Restriction of XML External Entity Reference Vulnerability | kev | markdown |  | CVE-2025-2776,cve-arsenal,cisa-kev,kev,in-the-wild,SysAid,SysAid On-Pr |
| 6135 | Google Chromium Mojo Sandbox Escape Vulnerability | kev | markdown |  | CVE-2025-2783,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium Moj |
| 6136 | Synacor Zimbra Collaboration Suite (ZCS) Cross-site Scripting Vulnerability | kev | markdown |  | CVE-2025-27915,cve-arsenal,cisa-kev,kev,in-the-wild,Synacor,Zimbra Col |
| 6137 | Srimax Output Messenger Directory Traversal Vulnerability | kev | markdown |  | CVE-2025-27920,cve-arsenal,cisa-kev,kev,in-the-wild,Srimax,Output Mess |
| 6138 | D-Link DIR-823X Command Injection Vulnerability | kev | markdown |  | CVE-2025-29635,cve-arsenal,cisa-kev,kev,in-the-wild,D-Link,DIR-823X |
| 6139 | Microsoft Windows Common Log File System (CLFS) Driver Use-After-Free Vulnerability | kev | markdown |  | CVE-2025-29824,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 6140 | tj-actions/changed-files GitHub Action Embedded Malicious Code Vulnerability | kev | markdown |  | CVE-2025-30066,cve-arsenal,cisa-kev,kev,in-the-wild,tj-actions,changed |
| 6141 | reviewdog/action-setup GitHub Action Embedded Malicious Code Vulnerability | kev | markdown |  | CVE-2025-30154,cve-arsenal,cisa-kev,kev,in-the-wild,reviewdog,action-s |
| 6142 | Microsoft Windows Scripting Engine Type Confusion Vulnerability | kev | markdown |  | CVE-2025-30397,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6143 | Microsoft Windows DWM Core Library Use-After-Free Vulnerability | kev | markdown |  | CVE-2025-30400,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6144 | Gladinet CentreStack and Triofox Use of Hard-coded Cryptographic Key Vulnerability | kev | markdown |  | CVE-2025-30406,cve-arsenal,cisa-kev,kev,in-the-wild,Gladinet,CentreSta |
| 6145 | Vite Vitejs Improper Access Control Vulnerability | kev | markdown |  | CVE-2025-31125,cve-arsenal,cisa-kev,kev,in-the-wild,Vite,Vitejs |
| 6146 | CrushFTP Authentication Bypass Vulnerability | kev | markdown |  | CVE-2025-31161,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,CrushFT |
| 6147 | Apple Multiple Products Memory Corruption Vulnerability | kev | markdown |  | CVE-2025-31200,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 6148 | Apple Multiple Products Arbitrary Read and Write Vulnerability | kev | markdown |  | CVE-2025-31201,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 6149 | Apple Multiple Products Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-31277,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 6150 | SAP NetWeaver Unrestricted File Upload Vulnerability | kev | markdown |  | CVE-2025-31324,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SAP,Net |
| 6151 | Craft CMS Code Injection Vulnerability | kev | markdown |  | CVE-2025-32432,cve-arsenal,cisa-kev,kev,in-the-wild,Craft CMS |
| 6152 | Erlang Erlang/OTP SSH Server Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2025-32433,cve-arsenal,cisa-kev,kev,in-the-wild,Erlang,Erlang/OTP |
| 6153 | Sudo Inclusion of Functionality from Untrusted Control Sphere Vulnerability | kev | markdown |  | CVE-2025-32463,cve-arsenal,cisa-kev,kev,in-the-wild,Sudo |
| 6154 | Langflow Missing Authentication Vulnerability | kev | markdown |  | CVE-2025-3248,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Langflow |
| 6155 | Microsoft Windows Common Log File System (CLFS) Driver Use-After-Free Vulnerability | kev | markdown |  | CVE-2025-32701,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6156 | Microsoft Windows Common Log File System (CLFS) Driver Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-32706,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6157 | Microsoft Windows Ancillary Function Driver for WinSock Use-After-Free Vulnerability | kev | markdown |  | CVE-2025-32709,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6158 | Fortinet Multiple Products Stack-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-32756,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,Multiple  |
| 6159 | Quest KACE Systems Management Appliance (SMA) Improper Authentication Vulnerability | kev | markdown |  | CVE-2025-32975,cve-arsenal,cisa-kev,kev,in-the-wild,Quest,KACE Systems |
| 6160 | Microsoft Windows External Control of File Name or Path Vulnerability | kev | markdown |  | CVE-2025-33053,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6161 | Microsoft Windows SMB Client Improper Access Control Vulnerability | kev | markdown |  | CVE-2025-33073,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6162 | Versa Concerto Improper Authentication Vulnerability | kev | markdown |  | CVE-2025-34026,cve-arsenal,cisa-kev,kev,in-the-wild,Versa,Concerto |
| 6163 | Commvault Command Center Path Traversal Vulnerability | kev | markdown |  | CVE-2025-34028,cve-arsenal,cisa-kev,kev,in-the-wild,Commvault,Command  |
| 6164 | Langflow Origin Validation Error Vulnerability | kev | markdown |  | CVE-2025-34291,cve-arsenal,cisa-kev,kev,in-the-wild,Langflow |
| 6165 | Craft CMS External Control of Assumed-Immutable Web Parameter Vulnerability | kev | markdown |  | CVE-2025-35939,cve-arsenal,cisa-kev,kev,in-the-wild,Craft CMS |
| 6166 | Hewlett Packard Enterprise (HPE) OneView Code Injection Vulnerability | kev | markdown |  | CVE-2025-37164,cve-arsenal,cisa-kev,kev,in-the-wild,Hewlett Packard En |
| 6167 | Linux Kernel Time-of-Check Time-of-Use (TOCTOU) Race Condition Vulnerability | kev | markdown |  | CVE-2025-38352,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 6168 | Commvault Web Server Unspecified Vulnerability | kev | markdown |  | CVE-2025-3928,cve-arsenal,cisa-kev,kev,in-the-wild,Commvault,Web Serve |
| 6169 | ConnectWise ScreenConnect Improper Authentication Vulnerability | kev | markdown |  | CVE-2025-3935,cve-arsenal,cisa-kev,kev,in-the-wild,ConnectWise,ScreenC |
| 6170 | Linux Kernel Improper Check for Unusual or Exceptional Conditions Vulnerability | kev | markdown |  | CVE-2025-39682,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 6171 | Linux Kernel Race Condition Vulnerability | kev | markdown |  | CVE-2025-39964,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 6172 | Smartbedded Meteobridge Command Injection Vulnerability | kev | markdown |  | CVE-2025-4008,cve-arsenal,cisa-kev,kev,in-the-wild,Smartbedded,Meteobr |
| 6173 | SolarWinds Web Help Desk Security Control Bypass Vulnerability | kev | markdown |  | CVE-2025-40536,cve-arsenal,cisa-kev,kev,in-the-wild,SolarWinds,Web Hel |
| 6174 | SolarWinds Web Help Desk Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2025-40551,cve-arsenal,cisa-kev,kev,in-the-wild,SolarWinds,Web Hel |
| 6175 | SonicWall SMA1000 Missing Authorization Vulnerability | kev | markdown |  | CVE-2025-40602,cve-arsenal,cisa-kev,kev,in-the-wild,SonicWall,SMA1000  |
| 6176 | Broadcom VMware Aria Operations and VMware Tools Privilege Defined with Unsafe Actions Vulnerability | kev | markdown |  | CVE-2025-41244,cve-arsenal,cisa-kev,kev,in-the-wild,Broadcom,VMware Ar |
| 6177 | Danfoss AK-SM 800A Series Improper Authentication Vulnerability | kev | markdown |  | CVE-2025-41450,cve-arsenal,cisa-kev,kev,in-the-wild,Danfoss,AK-SM 800A |
| 6178 | Qualitia Active! Mail Stack-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-42599,cve-arsenal,cisa-kev,kev,in-the-wild,Qualitia,Active! M |
| 6179 | SAP NetWeaver Improper Neutralization of Special Elements in Data Query Logic Vulnerability | kev | markdown |  | CVE-2025-42884,cve-arsenal,cisa-kev,kev,in-the-wild,SAP,NetWeaver |
| 6180 | SAP NetWeaver Deserialization Vulnerability | kev | markdown |  | CVE-2025-42999,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SAP,Net |
| 6181 | Apple Multiple Products Unspecified Vulnerability | kev | markdown |  | CVE-2025-43200,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 6182 | Apple iOS, iPadOS, and macOS Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2025-43300,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,iOS, iPadOS, |
| 6183 | Apple Multiple Products Improper Locking Vulnerability | kev | markdown |  | CVE-2025-43510,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 6184 | Apple Multiple Products Classic Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-43520,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 6185 | Apple Multiple Products Use-After-Free WebKit Vulnerability | kev | markdown |  | CVE-2025-43529,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 6186 | Ivanti Endpoint Manager Mobile (EPMM) Authentication Bypass Vulnerability | kev | markdown |  | CVE-2025-4427,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Endpoint Man |
| 6187 | Ivanti Endpoint Manager Mobile (EPMM) Code Injection Vulnerability | kev | markdown |  | CVE-2025-4428,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Endpoint Man |
| 6188 | Samsung MagicINFO 9 Server Path Traversal Vulnerability | kev | markdown |  | CVE-2025-4632,cve-arsenal,cisa-kev,kev,in-the-wild,Samsung,MagicINFO 9 |
| 6189 | TeleMessage TM SGNL Hidden Functionality Vulnerability | kev | markdown |  | CVE-2025-47729,cve-arsenal,cisa-kev,kev,in-the-wild,TeleMessage,TM SGN |
| 6190 | Wing FTP Server Improper Neutralization of Null Byte or NUL Character Vulnerability | kev | markdown |  | CVE-2025-47812,cve-arsenal,cisa-kev,kev,in-the-wild,Wing FTP Server |
| 6191 | Wing FTP Server Information Disclosure Vulnerability | kev | markdown |  | CVE-2025-47813,cve-arsenal,cisa-kev,kev,in-the-wild,Wing FTP Server |
| 6192 | IGEL OS Use of a Key Past its Expiration Date Vulnerability | kev | markdown |  | CVE-2025-47827,cve-arsenal,cisa-kev,kev,in-the-wild,IGEL,IGEL OS |
| 6193 | Git Link Following Vulnerability | kev | markdown |  | CVE-2025-48384,cve-arsenal,cisa-kev,kev,in-the-wild,Git |
| 6194 | Android Runtime Use-After-Free Vulnerability | kev | markdown |  | CVE-2025-48543,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Runtime |
| 6195 | Android Framework Privilege Escalation Vulnerability | kev | markdown |  | CVE-2025-48572,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Framework |
| 6196 | Android Framework Integer Overflow Vulnerability | kev | markdown |  | CVE-2025-48595,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Framework |
| 6197 | Android Framework Information Disclosure Vulnerability | kev | markdown |  | CVE-2025-48633,cve-arsenal,cisa-kev,kev,in-the-wild,Android,Framework |
| 6198 | Synacor Zimbra Collaboration Suite (ZCS) Cross-site Scripting Vulnerability | kev | markdown |  | CVE-2025-48700,cve-arsenal,cisa-kev,kev,in-the-wild,Synacor,Zimbra Col |
| 6199 | CWP Control Web Panel OS Command Injection Vulnerability | kev | markdown |  | CVE-2025-48703,cve-arsenal,cisa-kev,kev,in-the-wild,CWP,Control Web Pa |
| 6200 | TeleMessage TM SGNL Initialization of a Resource with an Insecure Default Vulnerability | kev | markdown |  | CVE-2025-48927,cve-arsenal,cisa-kev,kev,in-the-wild,TeleMessage,TM SGN |
| 6201 | TeleMessage TM SGNL Exposure of Core Dump File to an Unauthorized Control Sphere Vulnerability | kev | markdown |  | CVE-2025-48928,cve-arsenal,cisa-kev,kev,in-the-wild,TeleMessage,TM SGN |
| 6202 | RoundCube Webmail Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2025-49113,cve-arsenal,cisa-kev,kev,in-the-wild,Roundcube,Webmail |
| 6203 | Microsoft SharePoint Code Injection Vulnerability | kev | markdown |  | CVE-2025-49704,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 6204 | Microsoft SharePoint Improper Authentication Vulnerability | kev | markdown |  | CVE-2025-49706,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 6205 | Dassault Systèmes DELMIA Apriso Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2025-5086,cve-arsenal,cisa-kev,kev,in-the-wild,Dassault Systèmes,D |
| 6206 | SmarterTools SmarterMail Unrestricted Upload of File with Dangerous Type Vulnerability | kev | markdown |  | CVE-2025-52691,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Smarter |
| 6207 | F5 BIG-IP Stack-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-53521,cve-arsenal,cisa-kev,kev,in-the-wild,F5,BIG-IP |
| 6208 | Sitecore Multiple Products Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2025-53690,cve-arsenal,cisa-kev,kev,in-the-wild,Sitecore,Multiple  |
| 6209 | Microsoft SharePoint Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2025-53770,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 6210 | Laravel Livewire Code Injection Vulnerability | kev | markdown |  | CVE-2025-54068,cve-arsenal,cisa-kev,kev,in-the-wild,Laravel,Livewire |
| 6211 | Google Chromium V8 Out-of-Bounds Read and Write Vulnerability | kev | markdown |  | CVE-2025-5419,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 6212 | Adobe Commerce and Magento Improper Input Validation Vulnerability | kev | markdown |  | CVE-2025-54236,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Commerce and |
| 6213 | Adobe Experience Manager Forms Code Execution Vulnerability | kev | markdown |  | CVE-2025-54253,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Experience M |
| 6214 | CrushFTP Unprotected Alternate Channel Vulnerability | kev | markdown |  | CVE-2025-54309,cve-arsenal,cisa-kev,kev,in-the-wild,CrushFTP |
| 6215 | Prettier eslint-config-prettier Embedded Malicious Code Vulnerability | kev | markdown |  | CVE-2025-54313,cve-arsenal,cisa-kev,kev,in-the-wild,Prettier,eslint-co |
| 6216 | Trend Micro Apex One OS Command Injection Vulnerability | kev | markdown |  | CVE-2025-54948,cve-arsenal,cisa-kev,kev,in-the-wild,Trend Micro,Apex O |
| 6217 | Meta Platforms WhatsApp Incorrect Authorization Vulnerability | kev | markdown |  | CVE-2025-55177,cve-arsenal,cisa-kev,kev,in-the-wild,Meta Platforms,Wha |
| 6218 | Meta React Server Components Remote Code Execution Vulnerability | kev | markdown |  | CVE-2025-55182,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Meta,Re |
| 6219 | Citrix NetScaler ADC and Gateway Out-of-Bounds Read Vulnerability | kev | markdown |  | CVE-2025-5777,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Citrix,N |
| 6220 | Sangoma FreePBX Authentication Bypass Vulnerability | kev | markdown |  | CVE-2025-57819,cve-arsenal,cisa-kev,kev,in-the-wild,Sangoma,FreePBX |
| 6221 | Fortinet FortiWeb OS Command Injection Vulnerability | kev | markdown |  | CVE-2025-58034,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,FortiWeb |
| 6222 | OSGeo GeoServer Improper Restriction of XML External Entity Reference Vulnerability | kev | markdown |  | CVE-2025-58360,cve-arsenal,cisa-kev,kev,in-the-wild,OSGeo,GeoServer |
| 6223 | Microsoft Windows Improper Access Control Vulnerability | kev | markdown |  | CVE-2025-59230,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6224 | Microsoft Windows Server Update Service (WSUS) Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2025-59287,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6225 | ASUS Live Update Embedded Malicious Code Vulnerability | kev | markdown |  | CVE-2025-59374,cve-arsenal,cisa-kev,kev,in-the-wild,ASUS,Live Update |
| 6226 | Libraesva Email Security Gateway Command Injection Vulnerability | kev | markdown |  | CVE-2025-59689,cve-arsenal,cisa-kev,kev,in-the-wild,Libraesva,Email Se |
| 6227 | Fortinet Multiple Products Improper Verification of Cryptographic Signature Vulnerability | kev | markdown |  | CVE-2025-59718,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,Multiple  |
| 6228 | Microsoft Windows Link Following Vulnerability | kev | markdown |  | CVE-2025-60710,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 6229 | Fortinet Multiple Products Path Traversal Vulnerability | kev | markdown |  | CVE-2025-61624,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,Multiple  |
| 6230 | Oracle Fusion Middleware Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2025-61757,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,Fusion Midd |
| 6231 | Oracle E-Business Suite Unspecified Vulnerability | kev | markdown |  | CVE-2025-61882,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Oracle, |
| 6232 | Oracle E-Business Suite Server-Side Request Forgery (SSRF) Vulnerability | kev | markdown |  | CVE-2025-61884,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Oracle, |
| 6233 | Motex LANSCOPE Endpoint Manager Improper Verification of Source of a Communication Channel Vulnerability | kev | markdown |  | CVE-2025-61932,cve-arsenal,cisa-kev,kev,in-the-wild,Motex,LANSCOPE End |
| 6234 | Dassault Systèmes DELMIA Apriso Code Injection Vulnerability | kev | markdown |  | CVE-2025-6204,cve-arsenal,cisa-kev,kev,in-the-wild,Dassault Systèmes,D |
| 6235 | Dassault Systèmes DELMIA Apriso Missing Authorization Vulnerability | kev | markdown |  | CVE-2025-6205,cve-arsenal,cisa-kev,kev,in-the-wild,Dassault Systèmes,D |
| 6236 | RARLAB WinRAR Path Traversal Vulnerability | kev | markdown |  | CVE-2025-6218,cve-arsenal,cisa-kev,kev,in-the-wild,RARLAB,WinRAR |
| 6237 | Microsoft Windows Race Condition Vulnerability | kev | markdown |  | CVE-2025-62215,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6238 | Microsoft Windows Use After Free Vulnerability | kev | markdown |  | CVE-2025-62221,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6239 | Ray-Project Ray Code Injection Vulnerability | kev | markdown |  | CVE-2025-62593,cve-arsenal,cisa-kev,kev,in-the-wild,Ray-Project,Ray |
| 6240 | Sangoma FreePBX OS Command Injection Vulnerability | kev | markdown |  | CVE-2025-64328,cve-arsenal,cisa-kev,kev,in-the-wild,Sangoma,FreePBX |
| 6241 | Fortinet FortiWeb Path Traversal Vulnerability | kev | markdown |  | CVE-2025-64446,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,FortiWeb |
| 6242 | Citrix NetScaler ADC and Gateway Buffer Overflow Vulnerability | kev | markdown |  | CVE-2025-6543,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,NetScaler AD |
| 6243 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2025-6554,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 6244 | Google Chromium ANGLE and GPU Improper Input Validation Vulnerability | kev | markdown |  | CVE-2025-6558,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium |
| 6245 | Synacor Zimbra Collaboration Suite (ZCS) Cross-Site Scripting Vulnerability | kev | markdown |  | CVE-2025-66376,cve-arsenal,cisa-kev,kev,in-the-wild,Synacor,Zimbra Col |
| 6246 | Array Networks ArrayOS AG OS Command Injection Vulnerability | kev | markdown |  | CVE-2025-66644,cve-arsenal,cisa-kev,kev,in-the-wild,Array Networks,Arr |
| 6247 | Lantronix EDS5000 Code Injection Vulnerability | kev | markdown |  | CVE-2025-67038,cve-arsenal,cisa-kev,kev,in-the-wild,Lantronix,EDS5000 |
| 6248 | RoundCube Webmail Cross-site Scripting Vulnerability | kev | markdown |  | CVE-2025-68461,cve-arsenal,cisa-kev,kev,in-the-wild,Roundcube,Webmail |
| 6249 | n8n Improper Control of Dynamically-Managed Code Resources Vulnerability | kev | markdown |  | CVE-2025-68613,cve-arsenal,cisa-kev,kev,in-the-wild,n8n |
| 6250 | Synacor Zimbra Collaboration Suite (ZCS) PHP Remote File Inclusion Vulnerability | kev | markdown |  | CVE-2025-68645,cve-arsenal,cisa-kev,kev,in-the-wild,Synacor,Zimbra Col |
| 6251 | Fortinet FortiOS Exposure of Sensitive Information to an Unauthorized Actor Vulnerability | kev | markdown |  | CVE-2025-68686,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,FortiOS |
| 6252 | Citrix NetScaler Memory Overflow Vulnerability | kev | markdown |  | CVE-2025-7775,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,NetScaler |
| 6253 | RARLAB WinRAR Path Traversal Vulnerability | kev | markdown |  | CVE-2025-8088,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,RARLAB,W |
| 6254 | Gogs Path Traversal Vulnerability | kev | markdown |  | CVE-2025-8110,cve-arsenal,cisa-kev,kev,in-the-wild,Gogs |
| 6255 | N-able N-Central Insecure Deserialization Vulnerability | kev | markdown |  | CVE-2025-8875,cve-arsenal,cisa-kev,kev,in-the-wild,N-able,N-Central |
| 6256 | N-able N-Central Command Injection Vulnerability | kev | markdown |  | CVE-2025-8876,cve-arsenal,cisa-kev,kev,in-the-wild,N-able,N-Central |
| 6257 | WatchGuard Firebox Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2025-9242,cve-arsenal,cisa-kev,kev,in-the-wild,WatchGuard,Firebox |
| 6258 | TP-Link Archer C7(EU) and TL-WR841N/ND(MS) OS Command Injection Vulnerability | kev | markdown |  | CVE-2025-9377,cve-arsenal,cisa-kev,kev,in-the-wild,TP-Link,Multiple Ro |
| 6259 | Palo Alto Networks PAN-OS Authentication Bypass Vulnerability | kev | markdown |  | CVE-2026-0257,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Palo Alt |
| 6260 | Palo Alto Networks PAN-OS Out-of-bounds Write Vulnerability | kev | markdown |  | CVE-2026-0300,cve-arsenal,cisa-kev,kev,in-the-wild,Palo Alto Networks, |
| 6261 | Palo Alto Networks PAN-OS Code Injection Vulnerability | kev | markdown |  | CVE-2026-0322,cve-arsenal,cisa-kev,kev,in-the-wild,Palo Alto Networks, |
| 6262 | Langflow Inclusion of Functionality from Untrusted Control Sphere Vulnerability | kev | markdown |  | CVE-2026-0770,cve-arsenal,cisa-kev,kev,in-the-wild,Langflow |
| 6263 | Zammad GmbH Zammad Session Fixation Vulnerability | kev | markdown |  | CVE-2026-102489,cve-arsenal,cisa-kev,kev,in-the-wild,Zammad GmbH,Zamma |
| 6264 | Zammad GmbH Zammad Improper Privilege Management Vulnerability | kev | markdown |  | CVE-2026-102490,cve-arsenal,cisa-kev,kev,in-the-wild,Zammad GmbH,Zamma |
| 6265 | BeyondTrust Remote Support Authentication Bypass by Spoofing Vulnerability | kev | markdown |  | CVE-2026-102764,cve-arsenal,cisa-kev,kev,in-the-wild,BeyondTrust,Remot |
| 6266 | GPUStack Improper Authentication Vulnerability | kev | markdown |  | CVE-2026-102895,cve-arsenal,cisa-kev,kev,in-the-wild,GPUStack |
| 6267 | Hewlett Packard Enterprise (HPE) Juniper Session Smart Conductor and Session Smart Router OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-103406,cve-arsenal,cisa-kev,kev,in-the-wild,Hewlett Packard E |
| 6268 | Fortinet FortiMail Path Traversal Vulnerability | kev | markdown |  | CVE-2026-104286,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,FortiMai |
| 6269 | Ivanti Connect Secure Authentication Bypass by Spoofing Vulnerability | kev | markdown |  | CVE-2026-104935,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Connect Se |
| 6270 | Ahsay AhsayCBS Improper Authentication Vulnerability | kev | markdown |  | CVE-2026-105133,cve-arsenal,cisa-kev,kev,in-the-wild,Ahsay,AhsayCBS |
| 6271 | Ahsay AhsayCBS Command Injection Vulnerability | kev | markdown |  | CVE-2026-105134,cve-arsenal,cisa-kev,kev,in-the-wild,Ahsay,AhsayCBS |
| 6272 | Ivanti Sentry OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-10520,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Sentry |
| 6273 | EG Secure Solutions SiteGuard Server Edition OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-106931,cve-arsenal,cisa-kev,kev,in-the-wild,EG Secure Solutio |
| 6274 | PHP Memory Corruption Vulnerability | kev | markdown |  | CVE-2026-109618,cve-arsenal,cisa-kev,kev,in-the-wild,PHP Group,PHP |
| 6275 | PHP Stale Internal Request-State Vulnerability | kev | markdown |  | CVE-2026-109619,cve-arsenal,cisa-kev,kev,in-the-wild,PHP Group,PHP |
| 6276 | Collabora Online – Built-in CODE Server OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-110447,cve-arsenal,cisa-kev,kev,in-the-wild,Collabora,Collabo |
| 6277 | Siemens Software Center Uncontrolled Search Path Element Vulnerability | kev | markdown |  | CVE-2026-110936,cve-arsenal,cisa-kev,kev,in-the-wild,Siemens,Software  |
| 6278 | Ivanti Connect Secure Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2026-111357,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Connect Se |
| 6279 | Google Brotli Unchecked Return Value Vulnerability | kev | markdown |  | CVE-2026-111497,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Brotli |
| 6280 | Barracuda Networks Multiple Products OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-111560,cve-arsenal,cisa-kev,kev,in-the-wild,Barracuda Network |
| 6281 | Openwave Messaging Limited Webtop Media Transcoder Incorrect Resource Transfer Between Spheres Vulnerability | kev | markdown |  | CVE-2026-111793,cve-arsenal,cisa-kev,kev,in-the-wild,Openwave Messagin |
| 6282 | Arista VeloCloud Orchestrator SQL Injection Vulnerability | kev | markdown |  | CVE-2026-114625,cve-arsenal,cisa-kev,kev,in-the-wild,Arista,VeloCloud  |
| 6283 | UnrealIRCd Improper Neutralization of Delimiters Vulnerability | kev | markdown |  | CVE-2026-114713,cve-arsenal,cisa-kev,kev,in-the-wild,UnrealIRCd |
| 6284 | Perforce Puppet Enterprise Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2026-115731,cve-arsenal,cisa-kev,kev,in-the-wild,Perforce,Puppet E |
| 6285 | Google Chromium V8 Out-of-Bounds Read and Write Vulnerability | kev | markdown |  | CVE-2026-11645,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 6286 | Barracuda Networks Email Security Gateway (ESG) Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | markdown |  | CVE-2026-116606,cve-arsenal,cisa-kev,kev,in-the-wild,Barracuda Network |
| 6287 | SignalWire FreeSwitch Code Injection Vulnerability | kev | markdown |  | CVE-2026-116779,cve-arsenal,cisa-kev,kev,in-the-wild,SignalWire,FreeSw |
| 6288 | Ivanti Connect Secure Generation of Predictable Numbers or Identifiers Vulnerability | kev | markdown |  | CVE-2026-118527,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Connect Se |
| 6289 | WordPress WordPress Core Time-of-check Time-of-use (TOCTOU) Race Condition Vulnerability | kev | markdown |  | CVE-2026-119977,cve-arsenal,cisa-kev,kev,in-the-wild,WordPress,WordPre |
| 6290 | Check Point Multiple Products Incorrect Permission Assignment for Critical Resource Vulnerability | kev | markdown |  | CVE-2026-123039,cve-arsenal,cisa-kev,kev,in-the-wild,Check Point,Multi |
| 6291 | OpenBSD OpenSSH Unexpected Status Code or Return Value Vulnerability | kev | markdown |  | CVE-2026-123158,cve-arsenal,cisa-kev,kev,in-the-wild,OpenBSD,OpenSSH |
| 6292 | IBM Aspera Faspex Incorrect Authorization Vulnerability | kev | markdown |  | CVE-2026-123218,cve-arsenal,cisa-kev,kev,in-the-wild,IBM,Aspera Faspex |
| 6293 | VMWare vCenter Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | markdown |  | CVE-2026-123890,cve-arsenal,cisa-kev,kev,in-the-wild,Broadcom,VMware v |
| 6294 | Microsoft Office Outlook Missing User Interaction for Potentially Unsafe Action Vulnerability | kev | markdown |  | CVE-2026-124194,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 6295 | Google Chromium Use-After-Free Vulnerability | kev | markdown |  | CVE-2026-125350,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium |
| 6296 | PTC Windchill and FlexPLM Improper Input Validation Vulnerability | kev | markdown |  | CVE-2026-12569,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,PTC,Win |
| 6297 | Ivanti Endpoint Manager Mobile (EPMM) Code Injection Vulnerability | kev | markdown |  | CVE-2026-1281,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Endpoint Man |
| 6298 | Apache Roller Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2026-129448,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Roller |
| 6299 | Idunno AntiCsrf Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2026-131914,cve-arsenal,cisa-kev,kev,in-the-wild,Idunno,AntiCsrf |
| 6300 | Drupal Core Remote Code Execution Vulnerability | kev | markdown |  | CVE-2026-132165,cve-arsenal,cisa-kev,kev,in-the-wild,Drupal,Core |
| 6301 | Optimizely Web Experimentation External Control of Critical State Data Vulnerability | kev | markdown |  | CVE-2026-132224,cve-arsenal,cisa-kev,kev,in-the-wild,Optimizely,Web Ex |
| 6302 | Mitel MiCollab Unrestricted Upload of File with Dangerous Type Vulnerability | kev | markdown |  | CVE-2026-133275,cve-arsenal,cisa-kev,kev,in-the-wild,Mitel,MiCollab |
| 6303 | Mitel MiCollab XML Injection Vulnerability | kev | markdown |  | CVE-2026-133303,cve-arsenal,cisa-kev,kev,in-the-wild,Mitel,MiCollab |
| 6304 | Mitel MiCollab Weak Authentication Vulnerability | kev | markdown |  | CVE-2026-133311,cve-arsenal,cisa-kev,kev,in-the-wild,Mitel,MiCollab |
| 6305 | Mitel MiCollab OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-133323,cve-arsenal,cisa-kev,kev,in-the-wild,Mitel,MiCollab |
| 6306 | Mitel MiCollab Not Failing Securely Vulnerability | kev | markdown |  | CVE-2026-133324,cve-arsenal,cisa-kev,kev,in-the-wild,Mitel,MiCollab |
| 6307 | Vates Xen Orchestra Improper Validation of Specified Type of Input Vulnerability | kev | markdown |  | CVE-2026-133469,cve-arsenal,cisa-kev,kev,in-the-wild,Vates,Xen Orchest |
| 6308 | Project MONAI MONAI Label Path Traversal Vulnerability | kev | markdown |  | CVE-2026-133968,cve-arsenal,cisa-kev,kev,in-the-wild,Project MONAI,MON |
| 6309 | Ivanti Endpoint Manager Mobile (EPMM) Code Injection Vulnerability | kev | markdown |  | CVE-2026-1340,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Endpoint Man |
| 6310 | N-able MailAssure (SpamExperts) Server-Side Request Forgery (SSRF) Vulnerability | kev | markdown |  | CVE-2026-134624,cve-arsenal,cisa-kev,kev,in-the-wild,N-able,MailAssure |
| 6311 | SonicWall SMA1000 Appliances Server-Side Request Forgery Vulnerability | kev | markdown |  | CVE-2026-15409,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SonicWa |
| 6312 | SonicWall SMA1000 Appliances Code Injection Vulnerability | kev | markdown |  | CVE-2026-15410,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,SonicWa |
| 6313 | Ivanti Endpoint Manager (EPM) Authentication Bypass Vulnerability | kev | markdown |  | CVE-2026-1603,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Endpoint Man |
| 6314 | Check Point SmartConsole Improper Authentication Vulnerability | kev | markdown |  | CVE-2026-16232,cve-arsenal,cisa-kev,kev,in-the-wild,Check Point,SmartC |
| 6315 | Arista VeloCloud Orchestrator On-Prem OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-16812,cve-arsenal,cisa-kev,kev,in-the-wild,Arista,VeloCloud O |
| 6316 | BeyondTrust Remote Support (RS) and Privileged Remote Access (PRA) OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-1731,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,BeyondTr |
| 6317 | N-able N-central Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | markdown |  | CVE-2026-18556,cve-arsenal,cisa-kev,kev,in-the-wild,N-able,N-central |
| 6318 | N-able N-central Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | markdown |  | CVE-2026-18577,cve-arsenal,cisa-kev,kev,in-the-wild,N-able,N-central |
| 6319 | Citrix NetScaler Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | markdown |  | CVE-2026-19490,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,NetScaler |
| 6320 | SonicWall SonicOS SSLVPN Use of Single-factor Authentication Vulnerability | kev | markdown |  | CVE-2026-19868,cve-arsenal,cisa-kev,kev,in-the-wild,SonicWall,SonicOS  |
| 6321 | Cisco Unified Communications Products Code Injection Vulnerability | kev | markdown |  | CVE-2026-20045,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Unified Comm |
| 6322 | Cisco Firewall Management Center Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | markdown |  | CVE-2026-20079,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Secure Firew |
| 6323 | Cisco Catalyst SD-WAN Manager Incorrect Use of Privileged APIs Vulnerability | kev | markdown |  | CVE-2026-20122,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Catalyst SD- |
| 6324 | Cisco Catalyst SD-WAN Controller and Manager Authentication Bypass Vulnerability | kev | markdown |  | CVE-2026-20127,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Catalyst SD- |
| 6325 | Cisco Catalyst SD-WAN Manager Storing Passwords in a Recoverable Format Vulnerability | kev | markdown |  | CVE-2026-20128,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Catalyst SD- |
| 6326 | Cisco Secure Firewall Management Center (FMC) Software and Cisco Security Cloud Control (SCC) Firewall Management Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2026-20131,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Cisco,S |
| 6327 | Cisco Catalyst SD-WAN Manager Exposure of Sensitive Information to an Unauthorized Actor Vulnerability | kev | markdown |  | CVE-2026-20133,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Catalyst SD- |
| 6328 | Cisco Catalyst SD-WAN Controller Authentication Bypass Vulnerability | kev | markdown |  | CVE-2026-20182,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Catalyst SD- |
| 6329 | Cisco Unified Communications Manager Server-Side Request Forgery (SSRF) Vulnerability | kev | markdown |  | CVE-2026-20230,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Unified Comm |
| 6330 | Cisco Catalyst SD-WAN Manager Improper Encoding or Escaping of Output Vulnerability | kev | markdown |  | CVE-2026-20245,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Catalyst SD- |
| 6331 | Splunk Enterprise Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2026-20253,cve-arsenal,cisa-kev,kev,in-the-wild,Splunk,Enterprise |
| 6332 | Cisco Catalyst SD-WAN Manager Directory or Path Traversal Vulnerability | kev | markdown |  | CVE-2026-20262,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Catalyst SD- |
| 6333 | Cisco Secure Firewall Management Center Use of Hard-coded Password Vulnerability | kev | markdown |  | CVE-2026-20316,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Cisco,S |
| 6334 | Cisco BroadWorks Improper Restriction of XML External Entity Reference Vulnerability | kev | markdown |  | CVE-2026-20320,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,BroadWorks |
| 6335 | Cisco Secure Firewall Adaptive Security Appliance (ASA) and Secure Firewall Threat Defense (FTD) Heap Inspection Vulnerability | kev | markdown |  | CVE-2026-20349,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Secure Firew |
| 6336 | Apple Multiple Buffer Overflow Vulnerability | kev | markdown |  | CVE-2026-20700,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 6337 | Microsoft Windows Information Disclosure Vulnerability | kev | markdown |  | CVE-2026-20805,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6338 | Microsoft SharePoint Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2026-20963,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,SharePoi |
| 6339 | Qualcomm Multiple Chipsets Memory Corruption Vulnerability | kev | markdown |  | CVE-2026-21385,cve-arsenal,cisa-kev,kev,in-the-wild,Qualcomm,Multiple  |
| 6340 | Microsoft Office Security Feature Bypass Vulnerability | kev | markdown |  | CVE-2026-21509,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 6341 | Microsoft Windows Shell Protection Mechanism Failure Vulnerability | kev | markdown |  | CVE-2026-21510,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6342 | Microsoft MSHTML Framework Protection Mechanism Failure Vulnerability | kev | markdown |  | CVE-2026-21513,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6343 | Microsoft Office Word Reliance on Untrusted Inputs in a Security Decision Vulnerability | kev | markdown |  | CVE-2026-21514,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Office |
| 6344 | Microsoft Windows Type Confusion Vulnerability | kev | markdown |  | CVE-2026-21519,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6345 | Microsoft Windows NULL Pointer Dereference Vulnerability | kev | markdown |  | CVE-2026-21525,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6346 | Microsoft Windows Improper Privilege Management Vulnerability | kev | markdown |  | CVE-2026-21533,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6347 | Atlassian Multiple Products Files or Directories Accessible to External Parties Vulnerability | kev | markdown |  | CVE-2026-21589,cve-arsenal,cisa-kev,kev,in-the-wild,Atlassian,Multiple |
| 6348 | Fortinet FortiClient EMS SQL Injection Vulnerability | kev | markdown |  | CVE-2026-21643,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,FortiClie |
| 6349 | Oracle HTTP Server and Oracle Weblogic Server Proxy Plug-in Improper Access Control Vulnerability | kev | markdown |  | CVE-2026-21962,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,HTTP Server |
| 6350 | Broadcom VMware Aria Operations Command Injection Vulnerability | kev | markdown |  | CVE-2026-22719,cve-arsenal,cisa-kev,kev,in-the-wild,Broadcom,VMware Ar |
| 6351 | Dell RecoverPoint for Virtual Machines (RP4VMs) Use of Hard-coded Credentials Vulnerability | kev | markdown |  | CVE-2026-22769,cve-arsenal,cisa-kev,kev,in-the-wild,Dell,RecoverPoint  |
| 6352 | SmarterTools SmarterMail Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | markdown |  | CVE-2026-23760,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Smarter |
| 6353 | GNU InetUtils Argument Injection Vulnerability | kev | markdown |  | CVE-2026-24061,cve-arsenal,cisa-kev,kev,in-the-wild,GNU,InetUtils |
| 6354 | Google Chromium CSS Use-After-Free Vulnerability | kev | markdown |  | CVE-2026-2441,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium |
| 6355 | SmarterTools SmarterMail Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2026-24423,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Smarter |
| 6356 | Fortinet Multiple Products Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | markdown |  | CVE-2026-24858,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,Multiple  |
| 6357 | Fortinet FortiSandbox OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-25089,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,FortiSand |
| 6358 | Soliton Systems K.K FileZen OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-25108,cve-arsenal,cisa-kev,kev,in-the-wild,Soliton Systems K. |
| 6359 | SolarWinds Serv-U Uncontrolled Resource Consumption Vulnerability | kev | markdown |  | CVE-2026-28318,cve-arsenal,cisa-kev,kev,in-the-wild,SolarWinds,Serv-U |
| 6360 | Citrix NetScaler Out-of-Bounds Read Vulnerability | kev | markdown |  | CVE-2026-3055,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,NetScaler |
| 6361 | Linux Kernel Incorrect Resource Transfer Between Spheres Vulnerability | kev | markdown |  | CVE-2026-31431,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 6362 | Microsoft SharePoint Server Improper Input Validation Vulnerability | kev | markdown |  | CVE-2026-32201,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,SharePoi |
| 6363 | Microsoft Windows Protection Mechanism Failure Vulnerability | kev | markdown |  | CVE-2026-32202,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6364 | Langflow Code Injection Vulnerability | kev | markdown |  | CVE-2026-33017,cve-arsenal,cisa-kev,kev,in-the-wild,Langflow |
| 6365 | Aquasecurity Trivy Embedded Malicious Code Vulnerability | kev | markdown |  | CVE-2026-33634,cve-arsenal,cisa-kev,kev,in-the-wild,Aquasecurity,Trivy |
| 6366 | Microsoft Internet Key Exchange (IKE) Service Extensions Double Free Vulnerability | kev | markdown |  | CVE-2026-33824,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Internet |
| 6367 | Microsoft Defender Insufficient Granularity of Access Control Vulnerability | kev | markdown |  | CVE-2026-33825,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 6368 | Apache ActiveMQ Improper Input Validation Vulnerability | kev | markdown |  | CVE-2026-34197,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,ActiveMQ |
| 6369 | Apache Tomcat Missing Encryption of Sensitive Data Vulnerability | kev | markdown |  | CVE-2026-34486,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Tomcat |
| 6370 | Adobe Acrobat and Reader Prototype Pollution Vulnerability | kev | markdown |  | CVE-2026-34621,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Acrobat and  |
| 6371 | Ubiquiti UniFi OS Improper Access Control Vulnerability | kev | markdown |  | CVE-2026-34908,cve-arsenal,cisa-kev,kev,in-the-wild,Ubiquiti,UniFi OS |
| 6372 | Ubiquiti UniFi OS Path Traversal Vulnerability | kev | markdown |  | CVE-2026-34909,cve-arsenal,cisa-kev,kev,in-the-wild,Ubiquiti,UniFi OS |
| 6373 | Ubiquiti UniFi OS Improper Input Validation Vulnerability | kev | markdown |  | CVE-2026-34910,cve-arsenal,cisa-kev,kev,in-the-wild,Ubiquiti,UniFi OS |
| 6374 | Trend Micro Apex One (On-Premise) Directory Traversal Vulnerability | kev | markdown |  | CVE-2026-34926,cve-arsenal,cisa-kev,kev,in-the-wild,Trend Micro,Apex O |
| 6375 | TrueConf Client Download of Code Without Integrity Check Vulnerability | kev | markdown |  | CVE-2026-3502,cve-arsenal,cisa-kev,kev,in-the-wild,TrueConf,Client |
| 6376 | Oracle PeopleSoft Enterprise PeopleTools Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2026-35273,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Oracle, |
| 6377 | Fortinet FortiClient EMS Improper Access Control Vulnerability | kev | markdown |  | CVE-2026-35616,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,FortiClie |
| 6378 | Google Skia Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2026-3909,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Skia |
| 6379 | Google Chromium V8 Improper Restriction of Operations Within the Bounds of a Memory Buffer Vulnerability | kev | markdown |  | CVE-2026-3910,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 6380 | Fortinet FortiSandbox OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-39808,cve-arsenal,cisa-kev,kev,in-the-wild,Fortinet,FortiSand |
| 6381 | Marimo Remote Code Execution Vulnerability | kev | markdown |  | CVE-2026-39987,cve-arsenal,cisa-kev,kev,in-the-wild,Marimo |
| 6382 | BerriAI LiteLLM Unprotected Alternate Channel Vulnerability | kev | markdown |  | CVE-2026-40217,cve-arsenal,cisa-kev,kev,in-the-wild,BerriAI,LiteLLM |
| 6383 | Microsoft Defender Link Following Vulnerability | kev | markdown |  | CVE-2026-41091,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Defender |
| 6384 | Broadcom VMWare ESX Insufficient Logging Vulnerability | kev | markdown |  | CVE-2026-41709,cve-arsenal,cisa-kev,kev,in-the-wild,Broadcom,VMware ES |
| 6385 | WebPros cPanel & WHM and WP2 (WordPress Squared) Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2026-41940,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,WebPros |
| 6386 | JFrog Artifactory Incorrect Authorization Vulnerability | kev | markdown |  | CVE-2026-42016,cve-arsenal,cisa-kev,kev,in-the-wild,JFrog,Artifactory |
| 6387 | JFrog Artifactory Improper Authentication Vulnerability | kev | markdown |  | CVE-2026-42018,cve-arsenal,cisa-kev,kev,in-the-wild,JFrog,Artifactory |
| 6388 | BerriAI LiteLLM SQL Injection Vulnerability | kev | markdown |  | CVE-2026-42208,cve-arsenal,cisa-kev,kev,in-the-wild,BerriAI,LiteLLM |
| 6389 | BerriAI LiteLLM Command Injection Vulnerability | kev | markdown |  | CVE-2026-42271,cve-arsenal,cisa-kev,kev,in-the-wild,BerriAI,LiteLLM |
| 6390 | Microsoft Exchange Server Cross-Site Scripting Vulnerability | kev | markdown |  | CVE-2026-42897,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft |
| 6391 | Mirasvit Full Page Cache Warmer Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2026-45247,cve-arsenal,cisa-kev,kev,in-the-wild,Mirasvit,Mirasvit  |
| 6392 | TanStack Unspecified Vulnerability | kev | markdown |  | CVE-2026-45321,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,TanStac |
| 6393 | Microsoft Defender Denial of Service Vulnerability | kev | markdown |  | CVE-2026-45498,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Defender |
| 6394 | Microsoft SharePoint Server Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2026-45659,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Microso |
| 6395 | Linux Kernel Missing Initialization of a Variable Vulnerability | kev | markdown |  | CVE-2026-46300,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 6396 | Oracle E-Business Suite Improper Privilege Management Vulnerability | kev | markdown |  | CVE-2026-46817,cve-arsenal,cisa-kev,kev,in-the-wild,Oracle,E-Business  |
| 6397 | Nx Console Embedded Malicious Code Vulnerability | kev | markdown |  | CVE-2026-48027,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Nx,Nx C |
| 6398 | LiteSpeed cPanel Plugin Privilege Escalation Vulnerability | kev | markdown |  | CVE-2026-48172,cve-arsenal,cisa-kev,kev,in-the-wild,LiteSpeed,cPanel P |
| 6399 | Adobe ColdFusion Path Traversal Vulnerability | kev | markdown |  | CVE-2026-48282,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,ColdFusion |
| 6400 | Adobe Commerce and Magento Unrestricted Upload of File with Dangerous Type Vulnerability | kev | markdown |  | CVE-2026-48356,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Commerce and |
| 6401 | SimpleHelp Authentication Bypass Vulnerability | kev | markdown |  | CVE-2026-48558,cve-arsenal,cisa-kev,kev,in-the-wild,SimpleHelp |
| 6402 | Kludex Starlette HTTP Request/Response Smuggling Vulnerability | kev | markdown |  | CVE-2026-48710,cve-arsenal,cisa-kev,kev,in-the-wild,Kludex,Starlette |
| 6403 | Widget Factory Joomla Content Editor Improper Access Control Vulnerability | kev | markdown |  | CVE-2026-48907,cve-arsenal,cisa-kev,kev,in-the-wild,Widget Factory,Joo |
| 6404 | JoomShaper SP Page Builder Unrestricted Upload of File with Dangerous Type Vulnerability | kev | markdown |  | CVE-2026-48908,cve-arsenal,cisa-kev,kev,in-the-wild,JoomShaper,SP Page |
| 6405 | iCagenda Unrestricted Upload of File with Dangerous Type Vulnerability | kev | markdown |  | CVE-2026-48939,cve-arsenal,cisa-kev,kev,in-the-wild,iCagenda |
| 6406 | Kestra OSS OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-49869,cve-arsenal,cisa-kev,kev,in-the-wild,Kestra,Kestra OSS |
| 6407 | Microsoft SharePoint Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2026-50522,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,SharePoi |
| 6408 | Check Point Security Gateway Improper Authentication Vulnerability | kev | markdown |  | CVE-2026-50751,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Check P |
| 6409 | Google Dawn Use-After-Free Vulnerability | kev | markdown |  | CVE-2026-5281,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Dawn |
| 6410 | Linux Kernel Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2026-53266,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 6411 | Linux Kernel Unspecified Vulnerability | kev | markdown |  | CVE-2026-53362,cve-arsenal,cisa-kev,kev,in-the-wild,Linux,Kernel |
| 6412 | WSO2 Multiple Products Path Traversal Vulnerability | kev | markdown |  | CVE-2026-5430,cve-arsenal,cisa-kev,kev,in-the-wild,WSO2,Multiple Produ |
| 6413 | LiteSpeed cPanel Plugin UNIX Symbolic Link (Symlink) Following Vulnerability | kev | markdown |  | CVE-2026-54420,cve-arsenal,cisa-kev,kev,in-the-wild,LiteSpeed,cPanel P |
| 6414 | Microsoft SharePoint Weak Authentication Vulnerability | kev | markdown |  | CVE-2026-55040,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,SharePoi |
| 6415 | Langflow Authorization Bypass Through User-Controlled Key Vulnerability | kev | markdown |  | CVE-2026-55255,cve-arsenal,cisa-kev,kev,in-the-wild,Langflow |
| 6416 | Microsoft Active Directory Federation Services Insufficient Granularity of Access Control Vulnerability | kev | markdown |  | CVE-2026-56155,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Active D |
| 6417 | Microsoft SharePoint Server Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2026-56164,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,SharePoi |
| 6418 | Joomlack Page Builder Improper Access Control Vulnerability | kev | markdown |  | CVE-2026-56290,cve-arsenal,cisa-kev,kev,in-the-wild,Joomlack,Page Buil |
| 6419 | Balbooa Forms Unrestricted Upload of File with Dangerous Type Vulnerability | kev | markdown |  | CVE-2026-56291,cve-arsenal,cisa-kev,kev,in-the-wild,Balbooa,Forms |
| 6420 | Microsoft SharePoint Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2026-58644,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,SharePoi |
| 6421 | Google Pixel Improper Authorization Vulnerability | kev | markdown |  | CVE-2026-58704,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Pixel |
| 6422 | Broadcom VMware vCenter Path Traversal Vulnerability | kev | markdown |  | CVE-2026-59310,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,Broadco |
| 6423 | BerriAI LiteLLM Improper Authentication Vulnerability | kev | markdown |  | CVE-2026-59822,cve-arsenal,cisa-kev,kev,in-the-wild,BerriAI,LiteLLM |
| 6424 | Gitea Code Injection Vulnerability | kev | markdown |  | CVE-2026-60004,cve-arsenal,cisa-kev,kev,in-the-wild,Gitea |
| 6425 | WordPress Core SQL Injection Vulnerability | kev | markdown |  | CVE-2026-60137,cve-arsenal,cisa-kev,kev,in-the-wild,WordPress,Core |
| 6426 | WordPress Core Interpretation Conflict Vulnerability | kev | markdown |  | CVE-2026-63030,cve-arsenal,cisa-kev,kev,in-the-wild,WordPress,Core |
| 6427 | JetBrains TeamCity Deserialization of Untrusted Data Vulnerability | kev | markdown |  | CVE-2026-63077,cve-arsenal,cisa-kev,kev,in-the-wild,ransomware,JetBrai |
| 6428 | MLflow Server-Side Request Forgery Vulnerability | kev | markdown |  | CVE-2026-64849,cve-arsenal,cisa-kev,kev,in-the-wild,MLflow |
| 6429 | Apple macOS Improper Authentication Vulnerability | kev | markdown |  | CVE-2026-65400,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,macOS |
| 6430 | Microsoft SharePoint Code Injection Vulnerability | kev | markdown |  | CVE-2026-65660,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,SharePoi |
| 6431 | JFrog Artifactory Improper Limitation of a Pathname to a Restricted Directory Vulnerability | kev | markdown |  | CVE-2026-66384,cve-arsenal,cisa-kev,kev,in-the-wild,JFrog,Artifactory |
| 6432 | MikroTik RouterOS Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2026-67277,cve-arsenal,cisa-kev,kev,in-the-wild,MikroTik,RouterOS |
| 6433 | Mikrotik RouterOS Improper Enforcement of Behavioral Workflow Vulnerability | kev | markdown |  | CVE-2026-67279,cve-arsenal,cisa-kev,kev,in-the-wild,MikroTik,RouterOS |
| 6434 | Microsoft Windows Ancillary Function Driver for WinSock Use-After-Free Vulnerability | kev | markdown |  | CVE-2026-68820,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows  |
| 6435 | Ivanti Endpoint Manager Mobile (EPMM) Improper Input Validation Vulnerability | kev | markdown |  | CVE-2026-6973,cve-arsenal,cisa-kev,kev,in-the-wild,Ivanti,Endpoint Man |
| 6436 | Adobe Commerce and Magento Incorrect Authorization Vulnerability | kev | markdown |  | CVE-2026-71362,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Commerce and |
| 6437 | TrueConf Server Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2026-72529,cve-arsenal,cisa-kev,kev,in-the-wild,TrueConf,Server |
| 6438 | TrueConf Server Code Injection Vulnerability | kev | markdown |  | CVE-2026-72530,cve-arsenal,cisa-kev,kev,in-the-wild,TrueConf,Server |
| 6439 | Zyxel GS1900 Series Switches Stack-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2026-7273,cve-arsenal,cisa-kev,kev,in-the-wild,Zyxel,GS1900 Series |
| 6440 | Metabase SQL Injection Vulnerability | kev | markdown |  | CVE-2026-72898,cve-arsenal,cisa-kev,kev,in-the-wild,Metabase |
| 6441 | Zimbra Collaboration Suite (ZCS) OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-73570,cve-arsenal,cisa-kev,kev,in-the-wild,Synacor,Zimbra Col |
| 6442 | Arista Extensible Operating System Incomplete Comparison with Missing Factors Vulnerability | kev | markdown |  | CVE-2026-7473,cve-arsenal,cisa-kev,kev,in-the-wild,Arista,Extensible O |
| 6443 | Adobe Commerce and Magento Improper Neutralization of Special Elements Used in a Template Engine Vulnerability | kev | markdown |  | CVE-2026-75650,cve-arsenal,cisa-kev,kev,in-the-wild,Adobe,Commerce and |
| 6444 | Cisco Secure Firewall ASA Software and Cisco Secure FTD Software Unchecked Return Value to NULL Pointer Dereference Vulnerability | kev | markdown |  | CVE-2026-76418,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Secure Firew |
| 6445 | Cisco Identity Services Engine Incorrect Use of Privileged APIs Vulnerability | kev | markdown |  | CVE-2026-76460,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Identity Ser |
| 6446 | Cisco Secure Email Gateway SQL Injection Vulnerability | kev | markdown |  | CVE-2026-76461,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Secure Email |
| 6447 | Cisco Expressway Series Use of Hard-coded Credentials Vulnerability | kev | markdown |  | CVE-2026-76466,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Expressway S |
| 6448 | Cisco Catalyst SD-WAN Manager Hex Encoding Vulnerability | kev | markdown |  | CVE-2026-76504,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Catalyst SD- |
| 6449 | Cisco Secure Email Gateway Relative Path Traversal Vulnerability | kev | markdown |  | CVE-2026-76519,cve-arsenal,cisa-kev,kev,in-the-wild,Cisco,Secure Email |
| 6450 | Progress LoadMaster Command Injection Vulnerability | kev | markdown |  | CVE-2026-8037,cve-arsenal,cisa-kev,kev,in-the-wild,Progress,LoadMaster |
| 6451 | PaperCut NG/MF Missing Authentication for Critical Function Vulnerability | kev | markdown |  | CVE-2026-81578,cve-arsenal,cisa-kev,kev,in-the-wild,PaperCut,NG/MF |
| 6452 | Microsoft Windows Link Following Vulnerability | kev | markdown |  | CVE-2026-81963,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6453 | PaperCut NG/MF Unsafe Reflection Vulnerability | kev | markdown |  | CVE-2026-82078,cve-arsenal,cisa-kev,kev,in-the-wild,PaperCut,NG/MF |
| 6454 | JFrog Artifactory Improper Authentication Vulnerability | kev | markdown |  | CVE-2026-82329,cve-arsenal,cisa-kev,kev,in-the-wild,JFrog,Artifactory |
| 6455 | SonicWall SMA1000 Appliances Server-Side Request Forgery Vulnerability | kev | markdown |  | CVE-2026-83548,cve-arsenal,cisa-kev,kev,in-the-wild,SonicWall,SMA1000  |
| 6456 | SonicWall SMA1000 Appliances OS Command Injection Vulnerability | kev | markdown |  | CVE-2026-83549,cve-arsenal,cisa-kev,kev,in-the-wild,SonicWall,SMA1000  |
| 6457 | Daemon Tools Lite Embedded Malicious Code Vulnerability | kev | markdown |  | CVE-2026-8398,cve-arsenal,cisa-kev,kev,in-the-wild,Daemon,Daemon Tools |
| 6458 | Citrix NetScaler ADC and NetScaler Gateway Improper Restriction of Operations within the Bounds of a Memory Buffer Vulnerability | kev | markdown |  | CVE-2026-8452,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,NetScaler AD |
| 6459 | ConnectWise ScreenConnect Improper Privilege Management and Missing Authorization Vulnerability | kev | markdown |  | CVE-2026-84869,cve-arsenal,cisa-kev,kev,in-the-wild,ConnectWise,Screen |
| 6460 | Google Chromium V8 Type Confusion Vulnerability | kev | markdown |  | CVE-2026-85046,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 6461 | Check Point Multiple Products Improper Certificate Validation Vulnerability | kev | markdown |  | CVE-2026-85102,cve-arsenal,cisa-kev,kev,in-the-wild,Check Point,Multip |
| 6462 | GitLab Community Edition and Enterprise Edition Path Traversal Vulnerability | kev | markdown |  | CVE-2026-85706,cve-arsenal,cisa-kev,kev,in-the-wild,GitLab,Community E |
| 6463 | Microsoft Windows Heap-Based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2026-85880,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6464 | MikroTik RouterOS Improper Neutralization of Argument Delimiters in a Command Vulnerability | kev | markdown |  | CVE-2026-86060,cve-arsenal,cisa-kev,kev,in-the-wild,MikroTik,RouterOS |
| 6465 | N-able N-central Static Code Injection Vulnerability | kev | markdown |  | CVE-2026-86218,cve-arsenal,cisa-kev,kev,in-the-wild,N-able,N-central |
| 6466 | Apple Multiple Products Out-of-Bounds Write Vulnerability | kev | markdown |  | CVE-2026-86950,cve-arsenal,cisa-kev,kev,in-the-wild,Apple,Multiple Pro |
| 6467 | Google Chromium V8 Out of Bounds Write Vulnerability | kev | markdown |  | CVE-2026-87491,cve-arsenal,cisa-kev,kev,in-the-wild,Google,Chromium V8 |
| 6468 | Acronis Backup Incorrect Default Permissions Vulnerability | kev | markdown |  | CVE-2026-87886,cve-arsenal,cisa-kev,kev,in-the-wild,Acronis,Backup |
| 6469 | WordPress Core Remote File Inclusion Vulnerability | kev | markdown |  | CVE-2026-87902,cve-arsenal,cisa-kev,kev,in-the-wild,WordPress,Core |
| 6470 | Citrix NetScaler Improper Input Validation Vulnerability | kev | markdown |  | CVE-2026-88771,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,NetScaler |
| 6471 | Citrix NetScaler Improper Restriction of Operations within the Bounds of a Memory Buffer Vulnerability | kev | markdown |  | CVE-2026-88772,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,NetScaler |
| 6472 | Citrix NetScaler Improper Restriction of Operations within the Bounds of a Memory Buffer Vulnerability | kev | markdown |  | CVE-2026-88779,cve-arsenal,cisa-kev,kev,in-the-wild,Citrix,NetScaler |
| 6473 | Apache Axis Unsafe Reflection Vulnerability | kev | markdown |  | CVE-2026-90436,cve-arsenal,cisa-kev,kev,in-the-wild,Apache,Axis |
| 6474 | Vanderbilt University Research Electronic Data Capture (REDCap) Code Injection Vulnerability | kev | markdown |  | CVE-2026-90817,cve-arsenal,cisa-kev,kev,in-the-wild,Vanderbilt Univers |
| 6475 | Drupal Core SQL Injection Vulnerability | kev | markdown |  | CVE-2026-9082,cve-arsenal,cisa-kev,kev,in-the-wild,Drupal,Core |
| 6476 | IBM Langflow Code Injection Vulnerability | kev | markdown |  | CVE-2026-9198,cve-arsenal,cisa-kev,kev,in-the-wild,IBM,Langflow |
| 6477 | Check Point Multiple Products Path Traversal Vulnerability | kev | markdown |  | CVE-2026-93616,cve-arsenal,cisa-kev,kev,in-the-wild,Check Point,Multip |
| 6478 | Arista VeloCloud Orchestrator Improper Input Validation Vulnerability | kev | markdown |  | CVE-2026-93952,cve-arsenal,cisa-kev,kev,in-the-wild,Arista,VeloCloud O |
| 6479 | F5 BIG-IP APM Heap-based Buffer Overflow Vulnerability | kev | markdown |  | CVE-2026-94127,cve-arsenal,cisa-kev,kev,in-the-wild,F5,BIG-IP APM |
| 6480 | Sangoma Switchvox SQL Injection Vulnerability | kev | markdown |  | CVE-2026-9586,cve-arsenal,cisa-kev,kev,in-the-wild,Sangoma,Switchvox |
| 6481 | Plex Media Server Server-Side Request Forgery (SSRF) Vulnerability | kev | markdown |  | CVE-2026-96652,cve-arsenal,cisa-kev,kev,in-the-wild,Plex,Media Server |
| 6482 | Plex Media Server Improper Neutralization of Encoded URI Schemes in a Web Page Vulnerability | kev | markdown |  | CVE-2026-96654,cve-arsenal,cisa-kev,kev,in-the-wild,Plex,Media Server |
| 6483 | Plex Media Server Server-Side Request Forgery (SSRF) Vulnerability | kev | markdown |  | CVE-2026-96655,cve-arsenal,cisa-kev,kev,in-the-wild,Plex,Media Server |
| 6484 | Plex Media Server External Control of File Name or Path Vulnerability | kev | markdown |  | CVE-2026-96656,cve-arsenal,cisa-kev,kev,in-the-wild,Plex,Media Server |
| 6485 | Microsoft Windows Command Injection Vulnerability | kev | markdown |  | CVE-2026-96951,cve-arsenal,cisa-kev,kev,in-the-wild,Microsoft,Windows |
| 6486 | GitLab Community Edition and Enterprise Edition Integer Overflow or Wraparound Vulnerability | kev | markdown |  | CVE-2026-97229,cve-arsenal,cisa-kev,kev,in-the-wild,GitLab,Community E |
| 6487 | AdBlock before 2.21 allows remote attackers to block arbitrary resources on arbitrary websites and to disable arbitrary blocking filters | nvd | markdown |  | CVE-2015-2692,cve-arsenal,nvd,cvss-critical |
| 6488 | The Data Protection component in the VMware vSphere GUI in IBM Tivoli Storage Manager for Virtual Environments: Data Protection for VMware (aka Spectrum Protect | nvd | markdown |  | CVE-2015-7425,cve-arsenal,nvd,cvss-critical |
| 6489 | The Data Protection extension in the VMware GUI in IBM Tivoli Storage Manager for Virtual Environments: Data Protection for VMware (aka Spectrum Protect for Vir | nvd | markdown |  | CVE-2015-7426,cve-arsenal,nvd,cvss-critical |
| 6490 | The initialize method in the Histogram class in lib/colorscore/histogram.rb in the colorscore gem before 0.0.5 for Ruby allows context-dependent attackers to ex | nvd | markdown |  | CVE-2015-7541,cve-arsenal,nvd,cvss-critical |
| 6491 | SearchBlox 8.3 before 8.3.1 allows remote attackers to write to the config file, and consequently cause a denial of service (application crash), via unspecified | nvd | markdown |  | CVE-2015-7919,cve-arsenal,nvd,cvss-critical |
| 6492 | Adcon Telemetry A840 Telemetry Gateway Base Station has hardcoded credentials, which allows remote attackers to obtain administrative access via unspecified vec | nvd | markdown |  | CVE-2015-7930,cve-arsenal,nvd,cvss-critical |
| 6493 | The PasswordReset.Controllers.ResetController.ChangePasswordIndex method in PasswordReset.dll in Dovestones AD Self Password Reset before 3.0.4 allows remote  | nvd | markdown |  | CVE-2015-8267,cve-arsenal,nvd,cvss-critical |
| 6494 | Integer overflow in the ImageRegionReader::ReadIntoBuffer function in MediaStorageAndFileFormat/gdcmImageRegionReader.cxx in Grassroots DICOM (aka GDCM) before  | nvd | markdown |  | CVE-2015-8396,cve-arsenal,nvd,cvss-critical |
| 6495 | Adobe Flash Player before 18.0.0 and 19.x and 20.x before 20.0.0 on Windows and OS X and before 11.2.202 on Linux, Adobe AIR before 20.0.0, Adob | nvd | markdown |  | CVE-2015-8459,cve-arsenal,nvd,cvss-critical |
| 6496 | Local privilege escalation vulnerability in the Gentoo QEMU package before 2.5.0-r1 | nvd | markdown |  | CVE-2015-8556,cve-arsenal,nvd,cvss-critical |
| 6497 | The idle stream handling in nghttp2 before 1.6.0 allows attackers to have unspecified impact via unknown vectors, aka a heap-use-after-free bug | nvd | markdown |  | CVE-2015-8659,cve-arsenal,nvd,cvss-critical |
| 6498 | The multifilesystem storage backend in Radicale before 1.1 allows remote attackers to read or write to arbitrary files via a crafted component name | nvd | markdown |  | CVE-2015-8747,cve-arsenal,nvd,cvss-critical |
| 6499 | SQL injection vulnerability in the Group Promotions module in the admin control panel in MyBB (aka MyBulletinBoard) before 1.6.18 and 1.8.x before 1.8.6 and MyB | nvd | markdown |  | CVE-2015-8974,cve-arsenal,nvd,cvss-critical |
| 6500 | MailEnable before 8.60 allows XXE via an XML document in the request.aspx Options parameter | nvd | markdown |  | CVE-2015-9280,cve-arsenal,nvd,cvss-critical |
| 6501 | MySQL for PCF tiles 1.7.x before 1.7.10 were discovered to log the AWS access key in plaintext. These credentials were logged to the Service Backup component lo | nvd | markdown |  | CVE-2016-0898,cve-arsenal,nvd,cvss-critical |
| 6502 | An issue was discovered in Radisys MRF Web Panel (SWMS) 9.0.1. The MSM_MACRO_NAME POST parameter in /swms/ms.cgi was discovered to be vulnerable to OS command i | nvd | markdown |  | CVE-2016-10043,cve-arsenal,nvd,cvss-critical |
| 6503 | Adobe Reader and Acrobat before 11.0.16, Acrobat and Acrobat Reader DC Classic before 15.006.30172, and Acrobat and Acrobat Reader DC Continuous before 15.016.2 | nvd | markdown |  | CVE-2016-1038,cve-arsenal,nvd,cvss-critical |
| 6504 | Adobe Reader and Acrobat before 11.0.16, Acrobat and Acrobat Reader DC Classic before 15.006.30172, and Acrobat and Acrobat Reader DC Continuous before 15.016.2 | nvd | markdown |  | CVE-2016-1041,cve-arsenal,nvd,cvss-critical |
| 6505 | Adobe Reader and Acrobat before 11.0.16, Acrobat and Acrobat Reader DC Classic before 15.006.30172, and Acrobat and Acrobat Reader DC Continuous before 15.016.2 | nvd | markdown |  | CVE-2016-1044,cve-arsenal,nvd,cvss-critical |
| 6506 | The XML parser in Cisco Information Server (CIS) 6.2 allows remote attackers to read arbitrary files or cause a denial of service (CPU and memory consumption) v | nvd | markdown |  | CVE-2016-1343,cve-arsenal,nvd,cvss-critical |
| 6507 | The filesystem storage backend in Radicale before 1.1 on Windows allows remote attackers to read or write to arbitrary files via a crafted path, as demonstrated | nvd | markdown |  | CVE-2016-1505,cve-arsenal,nvd,cvss-critical |
| 6508 | Multiple unspecified vulnerabilities in the browser engine in Mozilla Firefox before 44.0 allow remote attackers to cause a denial of service (memory corruption | nvd | markdown |  | CVE-2016-1931,cve-arsenal,nvd,cvss-critical |
| 6509 | HPE Operations Manager 8.x and 9.0 on Windows allows remote attackers to execute arbitrary commands via a crafted serialized Java object, related to the Apache  | nvd | markdown |  | CVE-2016-1985,cve-arsenal,nvd,cvss-critical |
| 6510 | Pulse Connect Secure (PCS) 8.2 before 8.2r1, 8.1 before 8.1r2, 8.0 before 8.0r10, and 7.4 before 7.4r13.4 allow remote attackers to read sensitive system authen | nvd | markdown |  | CVE-2016-4787,cve-arsenal,nvd,cvss-critical |
| 6511 | General Electric (GE) Bently Nevada 3500/22M USB with firmware before 5.0 and Bently Nevada 3500/22M Serial have open ports, which makes it easier for remote at | nvd | markdown |  | CVE-2016-5788,cve-arsenal,nvd,cvss-critical |
| 6512 | IBM BigFix Platform could allow a remote attacker to execute arbitrary code on the system, caused by a use-after-free race condition. An attacker could exploit  | nvd | markdown |  | CVE-2016-6082,cve-arsenal,nvd,cvss-critical |
| 6513 | VMware vRealize Operations (aka vROps) 6.x before 6.4.0 allows remote authenticated users to gain privileges, or halt and remove virtual machines, via unspecifi | nvd | markdown |  | CVE-2016-7457,cve-arsenal,nvd,cvss-critical |
| 6514 | SQL injection vulnerability in core services in Intel Security McAfee ePolicy Orchestrator (ePO) 5.3.2 and earlier and 5.1.3 and earlier allows attackers to alt | nvd | markdown |  | CVE-2016-8027,cve-arsenal,nvd,cvss-critical |
| 6515 | An issue was discovered in Schneider Electric ConneXium firewalls TCSEFEC23F3F20 all versions, TCSEFEC23F3F21 all versions, TCSEFEC23FCF20 all versions, TCSEFEC | nvd | markdown |  | CVE-2016-8352,cve-arsenal,nvd,cvss-critical |
| 6516 | An issue was discovered in Moxa OnCell OnCellG3470A-LTE, AWK-1131A/3131A/4131A Series, AWK-3191 Series, AWK-5232/6232 Series, AWK-1121/1127 Series, WAC-1001 V2  | nvd | markdown |  | CVE-2016-8363,cve-arsenal,nvd,cvss-critical |
| 6517 | IBM UrbanCode Deploy could allow a user to execute code using a specially crafted file upload that would replace code on the server. This code could be executed | nvd | markdown |  | CVE-2016-8938,cve-arsenal,nvd,cvss-critical |
| 6518 | A hard-coded cryptographic key vulnerability was identified in Red Lion Controls Sixnet-Managed Industrial Switches running firmware Version 5.0.196 and Stride- | nvd | markdown |  | CVE-2016-9335,cve-arsenal,nvd,cvss-critical |
| 6519 | An issue was discovered in Rockwell Automation Logix5000 Programmable Automation Controller FRN 16.00 through 21.00 (excluding all firmware versions prior to FR | nvd | markdown |  | CVE-2016-9343,cve-arsenal,nvd,cvss-critical |
| 6520 | Vulnerability in the Oracle WebLogic Server component of Oracle Fusion Middleware (subcomponent: JNDI). Supported versions that are affected are 10.3.6.0 and 12 | nvd | markdown |  | CVE-2017-10137,cve-arsenal,nvd,cvss-critical |
| 6521 | Vulnerability in the Oracle Identity Manager component of Oracle Fusion Middleware (subcomponent: Default Account). Supported versions that are affected are 11 | nvd | markdown |  | CVE-2017-10151,cve-arsenal,nvd,cvss-critical |
| 6522 | Vulnerability in the Oracle Tuxedo component of Oracle Fusion Middleware (subcomponent: Core). Supported versions that are affected are 11.1.1, 12.1.1, 12.1.3 a | nvd | markdown |  | CVE-2017-10269,cve-arsenal,nvd,cvss-critical |
| 6523 | Vulnerability in the Oracle Hospitality Reporting and Analytics component of Oracle Hospitality Applications (subcomponent: Report). Supported versions that are | nvd | markdown |  | CVE-2017-10402,cve-arsenal,nvd,cvss-critical |
| 6524 | Vulnerability in the Oracle Hospitality Reporting and Analytics component of Oracle Hospitality Applications (subcomponent: Report). Supported versions that are | nvd | markdown |  | CVE-2017-10405,cve-arsenal,nvd,cvss-critical |
| 6525 | Xen through 4.8.x mishandles page transfer, which allows guest OS users to obtain privileged host OS access, aka XSA-217 | nvd | markdown |  | CVE-2017-10912,cve-arsenal,nvd,cvss-critical |
| 6526 | Xen through 4.8.x does not validate memory allocations during certain P2M operations, which allows guest OS users to obtain privileged host OS access, aka XSA-2 | nvd | markdown |  | CVE-2017-10918,cve-arsenal,nvd,cvss-critical |
| 6527 | The grant-table feature in Xen through 4.8.x mishandles a GNTMAP_device_map and GNTMAP_host_map mapping, when followed by only a GNTMAP_host_map unmapping, whic | nvd | markdown |  | CVE-2017-10920,cve-arsenal,nvd,cvss-critical |
| 6528 | The grant-table feature in Xen through 4.8.x does not ensure sufficient type counts for a GNTMAP_device_map and GNTMAP_host_map mapping, which allows guest OS u | nvd | markdown |  | CVE-2017-10921,cve-arsenal,nvd,cvss-critical |
| 6529 | An issue was discovered in Adobe Connect 9.6.2 and earlier versions. A Server-Side Request Forgery (SSRF) vulnerability exists that could be abused to bypass ne | nvd | markdown |  | CVE-2017-11291,cve-arsenal,nvd,cvss-critical |
| 6530 | An exploitable heap overflow vulnerability exists in the tinysvcmdns library version 2016-07-18. A specially crafted packet can make the library overwrite an ar | nvd | markdown |  | CVE-2017-12087,cve-arsenal,nvd,cvss-critical |
| 6531 | A authentication bypass and execution of code vulnerability in HPE Integrated Lights-out 4 (iLO 4) version prior to 2.53 was found | nvd | markdown |  | CVE-2017-12542,cve-arsenal,nvd,cvss-critical |
| 6532 | Analysis of the Bomgar Remote Support Portal JavaStart.jar Applet 52790 and earlier revealed that it is vulnerable to a path traversal vulnerability. The archiv | nvd | markdown |  | CVE-2017-12815,cve-arsenal,nvd,cvss-critical |
| 6533 | Server Side Request Forgery vulnerability in Vebto Pixie Image Editor 1.4 and 1.7 allows remote attackers to disclose information or execute arbitrary code via  | nvd | markdown |  | CVE-2017-12905,cve-arsenal,nvd,cvss-critical |
| 6534 | An Improper Authentication issue was discovered in iniNet Solutions iniNet Webserver, all versions prior to V2.02.0100. The webserver does not properly authenti | nvd | markdown |  | CVE-2017-13995,cve-arsenal,nvd,cvss-critical |
| 6535 | EMC RSA Authentication Agent API 8.5 for C and RSA Authentication Agent SDK 8.6 for C allow attackers to bypass authentication, aka an "Error Handling Vulnerabi | nvd | markdown |  | CVE-2017-14378,cve-arsenal,nvd,cvss-critical |
| 6536 | An exploitable OS Command Injection vulnerability exists in the Telnet, SSH, and console login functionality of Moxa AWK-3131A Industrial IEEE 802.11a/b/g/n wir | nvd | markdown |  | CVE-2017-14459,cve-arsenal,nvd,cvss-critical |
| 6537 | An exploitable access control vulnerability exists in the data, program, and function file permissions functionality of Allen Bradley Micrologix 1400 Series B F | nvd | markdown |  | CVE-2017-14469,cve-arsenal,nvd,cvss-critical |
| 6538 | An exploitable access control vulnerability exists in the data, program, and function file permissions functionality of Allen Bradley Micrologix 1400 Series B F | nvd | markdown |  | CVE-2017-14470,cve-arsenal,nvd,cvss-critical |
| 6539 | An exploitable access control vulnerability exists in the data, program, and function file permissions functionality of Allen Bradley Micrologix 1400 Series B F | nvd | markdown |  | CVE-2017-14471,cve-arsenal,nvd,cvss-critical |
| 6540 | An exploitable access control vulnerability exists in the data, program, and function file permissions functionality of Allen Bradley Micrologix 1400 Series B F | nvd | markdown |  | CVE-2017-14472,cve-arsenal,nvd,cvss-critical |
| 6541 | An exploitable access control vulnerability exists in the data, program, and function file permissions functionality of Allen Bradley Micrologix 1400 Series B F | nvd | markdown |  | CVE-2017-14473,cve-arsenal,nvd,cvss-critical |
| 6542 | The safe-eval module describes itself as a safer version of eval. By accessing the object constructors, un-sanitized user input can access the entire standard l | nvd | markdown |  | CVE-2017-16088,cve-arsenal,nvd,cvss-critical |
| 6543 | A Buffer Overflow issue was discovered in Rockwell Automation Allen-Bradley MicroLogix 1400 Controllers, Series B and C Versions 21.002 and earlier. The stack-b | nvd | markdown |  | CVE-2017-16740,cve-arsenal,nvd,cvss-critical |
| 6544 | hw/input/ps2.c in Qemu does not validate 'rptr' and 'count' values during guest migration, leading to out-of-bounds access | nvd | markdown |  | CVE-2017-16845,cve-arsenal,nvd,cvss-critical |
| 6545 | A vulnerability in Juniper Networks NorthStar Controller Application prior to version 2.1.0 Service Pack 1 may allow an unauthenticated, unprivileged, network-b | nvd | markdown |  | CVE-2017-2320,cve-arsenal,nvd,cvss-critical |
| 6546 | The Integrated User Firewall (UserFW) feature was introduced in Junos OS version 12.1X47-D10 on the Juniper SRX Series devices to provide simple integration of  | nvd | markdown |  | CVE-2017-2343,cve-arsenal,nvd,cvss-critical |
| 6547 | An exploitable buffer overflow exists in the psnotifyd application of the Pharos PopUp printer client version 9.0. A specially crafted packet can be sent to the | nvd | markdown |  | CVE-2017-2785,cve-arsenal,nvd,cvss-critical |
| 6548 | A buffer overflows exists in the psnotifyd application of the Pharos PopUp printer client version 9.0. A specially crafted packet can be sent to the victim's co | nvd | markdown |  | CVE-2017-2788,cve-arsenal,nvd,cvss-critical |
| 6549 | Adobe Digital Editions versions 4.5.4 and earlier have an exploitable memory corruption vulnerability in the PDF runtime engine. Successful exploitation could l | nvd | markdown |  | CVE-2017-3088,cve-arsenal,nvd,cvss-critical |
| 6550 | Vulnerability in the Primavera P6 Enterprise Project Portfolio Management component of Oracle Primavera Products Suite (subcomponent: Web Access). Supported ver | nvd | markdown |  | CVE-2017-3324,cve-arsenal,nvd,cvss-critical |
| 6551 | Vulnerability in the Solaris component of Oracle Sun Systems Products Suite (subcomponent: Kernel RPC). For supported versions that are affected see note. Easil | nvd | markdown |  | CVE-2017-3623,cve-arsenal,nvd,cvss-critical |
| 6552 | A vulnerability in the web-based GUI of Cisco Prime Home could allow an unauthenticated, remote attacker to bypass authentication and execute actions with admin | nvd | markdown |  | CVE-2017-3791,cve-arsenal,nvd,cvss-critical |
| 6553 | An issue was discovered in Carlo Gavazzi VMU-C EM prior to firmware Version A11_U05, and VMU-C PV prior to firmware Version A17. Successful exploitation of this | nvd | markdown |  | CVE-2017-5145,cve-arsenal,nvd,cvss-critical |
| 6554 | When executing a program via the bubblewrap sandbox, the nonpriv session can escape to the parent session by using the TIOCSTI ioctl to push characters into the | nvd | markdown |  | CVE-2017-5226,cve-arsenal,nvd,cvss-critical |
| 6555 | The Symantec Messaging Gateway can encounter an issue of remote code execution, which describes a situation whereby an individual may obtain the ability to exec | nvd | markdown |  | CVE-2017-6326,cve-arsenal,nvd,cvss-critical |
| 6556 | Zoho ManageEngine Desktop Central before build 100082 allows remote attackers to obtain control over all connected active desktops via unspecified vectors | nvd | markdown |  | CVE-2017-7213,cve-arsenal,nvd,cvss-critical |
| 6557 | Uploaded XML documents were not correctly validated in Apache OpenMeetings 3.1.0 | nvd | markdown |  | CVE-2017-7664,cve-arsenal,nvd,cvss-critical |
| 6558 | In SolarWinds Log & Event Manager (LEM) before 6.3.1 Hotfix 4, a menu system is encountered when the SSH service is accessed with "cmc" and "password" (the defa | nvd | markdown |  | CVE-2017-7722,cve-arsenal,nvd,cvss-critical |
| 6559 | An Improper Access Control issue was discovered in Schweitzer Engineering Laboratories (SEL) SEL-3620 and SEL-3622 Security Gateway Versions R202 and, R203, R20 | nvd | markdown |  | CVE-2017-7928,cve-arsenal,nvd,cvss-critical |
| 6560 | Zyxel WRE6505 devices have a default TELNET password of 1234 for the root and admin accounts, which makes it easier for remote attackers to conduct DNS hijackin | nvd | markdown |  | CVE-2017-7964,cve-arsenal,nvd,cvss-critical |
| 6561 | www.modified-shop.org modified eCommerce Shopsoftware 2.0.2 rev 10690 has XXE in api/it-recht-kanzlei/api-it-recht-kanzlei.php | nvd | markdown |  | CVE-2017-8110,cve-arsenal,nvd,cvss-critical |
| 6562 | An issue was discovered on Accellion FTA devices before FTA_9_12_180. Because a regular expression (intended to match local https URLs) lacks an initial ^ chara | nvd | markdown |  | CVE-2017-8794,cve-arsenal,nvd,cvss-critical |
| 6563 | A vulnerability in the Secure Sockets Layer (SSL) VPN functionality of the Cisco Adaptive Security Appliance (ASA) Software could allow an unauthenticated, remo | nvd | markdown |  | CVE-2018-0101,cve-arsenal,nvd,cvss-critical |
| 6564 | A vulnerability in Cisco Digital Network Architecture (DNA) Center could allow an unauthenticated, remote attacker to log in to an affected system by using an a | nvd | markdown |  | CVE-2018-0222,cve-arsenal,nvd,cvss-critical |
| 6565 | A vulnerability in the container management subsystem of Cisco Digital Network Architecture (DNA) Center could allow an unauthenticated, remote attacker to bypa | nvd | markdown |  | CVE-2018-0268,cve-arsenal,nvd,cvss-critical |
| 6566 | I Librarian I-librarian version 4.8 and earlier contains a XML External Entity (XXE) vulnerability in line 154 of importmetadata.php(simplexml_load_string) that | nvd | markdown |  | CVE-2018-1000124,cve-arsenal,nvd,cvss-critical |
| 6567 | Eclipse RDF4j version < 2.4.0 Milestone 2 contains a XML External Entity (XXE) vulnerability in RDF4j XML parser parsing RDF files that can result in the disclo | nvd | markdown |  | CVE-2018-1000644,cve-arsenal,nvd,cvss-critical |
| 6568 | Stroom version <5.4.5 contains a XML External Entity (XXE) vulnerability in XML Parser that can result in disclosure of confidential data, denial of service, se | nvd | markdown |  | CVE-2018-1000651,cve-arsenal,nvd,cvss-critical |
| 6569 | JabRef version <=4.3.1 contains a XML External Entity (XXE) vulnerability in MsBibImporter XML Parser that can result in disclosure of confidential data, denial | nvd | markdown |  | CVE-2018-1000652,cve-arsenal,nvd,cvss-critical |
| 6570 | neo4j-contrib neo4j-apoc-procedures version before commit 45bc09c contains a XML External Entity (XXE) vulnerability in XML Parser that can result in Disclosure | nvd | markdown |  | CVE-2018-1000820,cve-arsenal,nvd,cvss-critical |
| 6571 | MicroMathematics version before commit 5c05ac8 contains a XML External Entity (XXE) vulnerability in SMathStudio files that can result in Disclosure of confiden | nvd | markdown |  | CVE-2018-1000821,cve-arsenal,nvd,cvss-critical |
| 6572 | codelibs fess version before commit faa265b contains a XML External Entity (XXE) vulnerability in GSA XML file parser that can result in Disclosure of confident | nvd | markdown |  | CVE-2018-1000822,cve-arsenal,nvd,cvss-critical |
| 6573 | exist version <= 5.0.0-RC4 contains a XML External Entity (XXE) vulnerability in XML Parser for REST Server that can result in Disclosure of confidential data,  | nvd | markdown |  | CVE-2018-1000823,cve-arsenal,nvd,cvss-critical |
| 6574 | FreeCol version <= nightly-2018-08-22 contains a XML External Entity (XXE) vulnerability in FreeColXMLReader parser that can result in Disclosure of confidentia | nvd | markdown |  | CVE-2018-1000825,cve-arsenal,nvd,cvss-critical |
| 6575 | XR3Player version <= V3.124 contains a XML External Entity (XXE) vulnerability in Playlist parser that can result in Disclosure of confidential data, denial of  | nvd | markdown |  | CVE-2018-1000830,cve-arsenal,nvd,cvss-critical |
| 6576 | K9Mail version <= v5.600 contains a XML External Entity (XXE) vulnerability in WebDAV response parser that can result in Disclosure of confidential data, denial | nvd | markdown |  | CVE-2018-1000831,cve-arsenal,nvd,cvss-critical |
| 6577 | KeePassDX version <= 2.5.0.0beta17 contains a XML External Entity (XXE) vulnerability in kdbx file parser that can result in Disclosure of confidential data, de | nvd | markdown |  | CVE-2018-1000835,cve-arsenal,nvd,cvss-critical |
| 6578 | UML Designer version <= 8.0.0 contains a XML External Entity (XXE) vulnerability in XML parser for plugins that can result in Disclosure of confidential data, d | nvd | markdown |  | CVE-2018-1000837,cve-arsenal,nvd,cvss-critical |
| 6579 | autopsy version <= 4.9.0 contains a XML External Entity (XXE) vulnerability in CaseMetadata XML Parser that can result in Disclosure of confidential data, denia | nvd | markdown |  | CVE-2018-1000838,cve-arsenal,nvd,cvss-critical |
| 6580 | A vulnerability in Trend Micro Control Manager (versions 6.0 and 7.0) could allow an attacker to conduct a server-side request forgery (SSRF) attack on vulnerab | nvd | markdown |  | CVE-2018-10511,cve-arsenal,nvd,cvss-critical |
| 6581 | Stack-based buffer overflow in Activision Infinity Ward Call of Duty Modern Warfare 2 before 2018-04-26 allows remote attackers to execute arbitrary code via cr | nvd | markdown |  | CVE-2018-10718,cve-arsenal,nvd,cvss-critical |
| 6582 | A SQL injection vulnerability in the web administration and quarantine components of Micro Focus Secure Messaging Gateway allows an unauthenticated remote attac | nvd | markdown |  | CVE-2018-12464,cve-arsenal,nvd,cvss-critical |
| 6583 | A vulnerability has been identified in TIM 1531 IRC (All version < V2.0). The devices was missing proper authentication on port 102/tcp, although configured. Su | nvd | markdown |  | CVE-2018-13816,cve-arsenal,nvd,cvss-critical |
| 6584 | FasterXML jackson-databind 2.x before 2.9.7 might allow remote attackers to conduct server-side request forgery (SSRF) attacks by leveraging failure to block th | nvd | markdown |  | CVE-2018-14721,cve-arsenal,nvd,cvss-critical |
| 6585 | A command injection vulnerability in the apex-publish-static-files npm module version <2.0.1 which allows arbitrary shell command execution through a maliciousl | nvd | markdown |  | CVE-2018-16462,cve-arsenal,nvd,cvss-critical |
| 6586 | In FreeBSD before 11.2-STABLE(r341486) and 11.2-RELEASE-p6, insufficient bounds checking in one of the device models provided by bhyve can permit a guest operat | nvd | markdown |  | CVE-2018-17160,cve-arsenal,nvd,cvss-critical |
| 6587 | IBM Security Access Manager Appliance 9.0.4 and 9.0.5 could allow remote code execution when Advanced Access Control or Federation services are running. IBM | nvd | markdown |  | CVE-2018-1722,cve-arsenal,nvd,cvss-critical |
| 6588 | An earlier fix for an Inter-process Communication (IPC) vulnerability, CVE-2011-3079, added authentication to communication between IPC endpoints and server par | nvd | markdown |  | CVE-2018-18505,cve-arsenal,nvd,cvss-critical |
| 6589 | Sandboxie 5.26 allows a Sandbox Escape via an "import os" statement, followed by os.system("cmd") or os.system("powershell"), within a .py file. NOTE: the vendo | nvd | markdown |  | CVE-2018-18748,cve-arsenal,nvd,cvss-critical |
| 6590 | The Kubernetes integration in GitLab Enterprise Edition 11.x before 11.2.8, 11.3.x before 11.3.9, and 11.4.x before 11.4.4 has SSRF | nvd | markdown |  | CVE-2018-18843,cve-arsenal,nvd,cvss-critical |
| 6591 | mPDF through 7.1.6, if deployed as a web application that accepts arbitrary HTML, allows SSRF, as demonstrated by a '<img src="http://192.168' substring that tr | nvd | markdown |  | CVE-2018-19047,cve-arsenal,nvd,cvss-critical |
| 6592 | An issue was discovered in the MQTT server in Contiki-NG before 4.2. The function parse_publish_vhdr() that parses MQTT PUBLISH messages with a variable length  | nvd | markdown |  | CVE-2018-19417,cve-arsenal,nvd,cvss-critical |
| 6593 | Vulnerability in the Sun ZFS Storage Appliance Kit (AK) component of Oracle Sun Systems Products Suite (subcomponent: Core Services). The supported version that | nvd | markdown |  | CVE-2018-2611,cve-arsenal,nvd,cvss-critical |
| 6594 | Vulnerability in the Oracle GoldenGate component of Oracle GoldenGate (subcomponent: Monitoring Manager). Supported versions that are affected are 12.1.2, 1 | nvd | markdown |  | CVE-2018-2913,cve-arsenal,nvd,cvss-critical |
| 6595 | Incorrect parsing in url-parse <1.4.3 returns wrong hostname which leads to multiple vulnerabilities such as SSRF, Open Redirect, Bypass Authentication Protocol | nvd | markdown |  | CVE-2018-3774,cve-arsenal,nvd,cvss-critical |
| 6596 | An exploitable vulnerability exists in the REST parser of video-core's HTTP server of the Samsung SmartThings Hub STH-ETH-250 - Firmware version 0.20.17. The vi | nvd | markdown |  | CVE-2018-3907,cve-arsenal,nvd,cvss-critical |
| 6597 | An exploitable heap overflow vulnerability exists in the WkbProgramLow function of WibuKey Network server management, version 6.40.2402.500. A specially crafted | nvd | markdown |  | CVE-2018-3991,cve-arsenal,nvd,cvss-critical |
| 6598 | An issue was discovered in certain Apple products. macOS before 10.13.3 is affected. The issue involves the "Sandbox" component. It allows bypass of a sandbox p | nvd | markdown |  | CVE-2018-4091,cve-arsenal,nvd,cvss-critical |
| 6599 | An issue was discovered in certain Apple products. macOS before 10.13.5 is affected. The issue involves the "Grand Central Dispatch" component. It allows attack | nvd | markdown |  | CVE-2018-4229,cve-arsenal,nvd,cvss-critical |
| 6600 | An issue was discovered in Adobe Acrobat Reader 2018.009.20050 and earlier versions, 2017.011.30070 and earlier versions, 2015.006.30394 and earlier versions. T | nvd | markdown |  | CVE-2018-4872,cve-arsenal,nvd,cvss-critical |
| 6601 | Authentication Bypass vulnerability in the administrative user interface in McAfee Web Gateway 7.8.1 through 7.8.1 allows remote attackers to execute arbitr | nvd | markdown |  | CVE-2018-6667,cve-arsenal,nvd,cvss-critical |
| 6602 | Stack-based Buffer Overflow vulnerability in libUPnPHndlr.so in Belkin Wemo Insight Smart Plug allows remote attackers to bypass local security protection via a | nvd | markdown |  | CVE-2018-6692,cve-arsenal,nvd,cvss-critical |
| 6603 | The VMware AirWatch Agent for Android prior to 8.2 and AirWatch Agent for Windows Mobile prior to 6.5.2 contain a remote code execution vulnerability in real ti | nvd | markdown |  | CVE-2018-6968,cve-arsenal,nvd,cvss-critical |
| 6604 | Juniper ATP ships with hard coded credentials in the Web Collector instance which gives an attacker the ability to take full control of any installation of the  | nvd | markdown |  | CVE-2019-0020,cve-arsenal,nvd,cvss-critical |
| 6605 | Juniper ATP ships with hard coded credentials in the Cyphort Core instance which gives an attacker the ability to take full control of any installation of the s | nvd | markdown |  | CVE-2019-0022,cve-arsenal,nvd,cvss-critical |
| 6606 | Zoho ManageEngine ADSelfService Plus 5.x before build 5703 has SSRF | nvd | markdown |  | CVE-2019-3905,cve-arsenal,nvd,cvss-critical |
| 6607 | Multiple use-after-free vulnerabilities in vnc.c in the VNC server in QEMU 0.10.6 and earlier might allow guest OS users to execute arbitrary code on the host O | nvd | markdown |  | CVE-2009-3616,cve-arsenal,nvd,cvss-critical |
| 6608 | The VMX process in VMware ESXi 3.5 through 4.1 and ESX 3.5 through 4.1 does not properly handle RPC commands, which allows guest OS users to cause a denial of s | nvd | markdown |  | CVE-2012-1516,cve-arsenal,nvd,cvss-critical |
| 6609 | The portal client in IBM Tivoli Monitoring (ITM) 6.2.2 through FP9, 6.2.3 through FP5, and 6.3.0 through FP6 allows remote authenticated users to gain privilege | nvd | markdown |  | CVE-2015-7411,cve-arsenal,nvd,cvss-critical |
| 6610 | eWON devices with firmware before 10.1s0 omit RBAC for I/O server information and status requests, which allows remote attackers to obtain sensitive information | nvd | markdown |  | CVE-2015-7926,cve-arsenal,nvd,cvss-critical |
| 6611 | The GMS ViewPoint (GMSVP) web application in Dell SonicWALL GMS, Analyzer, and UMA EM5000 7.2, 8.0, and 8.1 before Hotfix 168056 allows remote authenticated use | nvd | markdown |  | CVE-2016-2396,cve-arsenal,nvd,cvss-critical |
| 6612 | lshell 0.9.16 allows remote authenticated users to break out of a limited shell and execute arbitrary commands | nvd | markdown |  | CVE-2016-6902,cve-arsenal,nvd,cvss-critical |
| 6613 | lshell 0.9.16 allows remote authenticated users to break out of a limited shell and execute arbitrary commands | nvd | markdown |  | CVE-2016-6903,cve-arsenal,nvd,cvss-critical |
| 6614 | An issue was discovered in Smiths-Medical CADD-Solis Medication Safety Software, Version 1.0; 2.0; 3.0; and 3.1. CADD-Solis Medication Safety Software grants an | nvd | markdown |  | CVE-2016-8355,cve-arsenal,nvd,cvss-critical |
| 6615 | Remote Command Execution in com.trend.iwss.gui.servlet.ManagePatches in Trend Micro Interscan Web Security Virtual Appliance (IWSVA) version 6.5-SP2_Build_Linux | nvd | markdown |  | CVE-2016-9269,cve-arsenal,nvd,cvss-critical |
| 6616 | PricewaterhouseCoopers (PwC) ACE-ABAP 8.10.304 for SAP Security allows remote authenticated users to conduct ABAP injection attacks and execute arbitrary code v | nvd | markdown |  | CVE-2016-9832,cve-arsenal,nvd,cvss-critical |
| 6617 | Vulnerability in the OJVM component of Oracle Database Server. Supported versions that are affected are 11.2.0, 12.1.0 and 12.2.0. Easily exploitable vuln | nvd | markdown |  | CVE-2017-10202,cve-arsenal,nvd,cvss-critical |
| 6618 | Vulnerability in the Oracle Tuxedo component of Oracle Fusion Middleware (subcomponent: Core). Supported versions that are affected are 11.1.1, 12.1.1, 12.1.3 a | nvd | markdown |  | CVE-2017-10272,cve-arsenal,nvd,cvss-critical |
| 6619 | Vulnerability in the Oracle WebLogic Server component of Oracle Fusion Middleware (subcomponent: WLS - Web Services). The supported version that is affected are | nvd | markdown |  | CVE-2017-10352,cve-arsenal,nvd,cvss-critical |
| 6620 | Vulnerability in the Oracle Hospitality Cruise AffairWhere component of Oracle Hospitality Applications (subcomponent: AffairWhere). Supported versions that are | nvd | markdown |  | CVE-2017-10396,cve-arsenal,nvd,cvss-critical |
| 6621 | Vulnerability in the Oracle Hospitality Reporting and Analytics component of Oracle Hospitality Applications (subcomponent: iQuery). Supported versions that are | nvd | markdown |  | CVE-2017-10404,cve-arsenal,nvd,cvss-critical |
| 6622 | A vulnerability in the web console of the Cisco Cloud Services Platform (CSP) 2100 could allow an authenticated, remote attacker to interact maliciously with th | nvd | markdown |  | CVE-2017-12251,cve-arsenal,nvd,cvss-critical |
| 6623 | IBM Security Guardium 10.0 could allow a remote authenticated attacker to execute arbitrary commands on the system. By sending a specially-crafted request, an a | nvd | markdown |  | CVE-2017-1253,cve-arsenal,nvd,cvss-critical |
| 6624 | Remote enabling and disabling admin interface in Gemalto's HASP SRM, Sentinel HASP and Sentinel LDK products prior to Sentinel LDK RTE version 7.55 leads to new | nvd | markdown |  | CVE-2017-12822,cve-arsenal,nvd,cvss-critical |
| 6625 | XML external entity (XXE) vulnerability in the import package functionality of the deployment module in Lansweeper before 6.0.100 allows remote authenticated | nvd | markdown |  | CVE-2017-13706,cve-arsenal,nvd,cvss-critical |
| 6626 | An exploitable buffer overflow vulnerability exists in Insteon Hub running firmware version 1012. The HTTP server implementation incorrectly handles the URL par | nvd | markdown |  | CVE-2017-14444,cve-arsenal,nvd,cvss-critical |
| 6627 | An exploitable buffer overflow vulnerability exists in Insteon Hub running firmware version 1012. The HTTP server implementation incorrectly handles the host pa | nvd | markdown |  | CVE-2017-14445,cve-arsenal,nvd,cvss-critical |
| 6628 | An exploitable stack-based buffer overflow vulnerability exists in Insteon Hub running firmware version 1012. The HTTP server implementation unsafely extracts p | nvd | markdown |  | CVE-2017-14446,cve-arsenal,nvd,cvss-critical |
| 6629 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01bad0 the value for the | nvd | markdown |  | CVE-2017-16338,cve-arsenal,nvd,cvss-critical |
| 6630 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01bb1c the value for the | nvd | markdown |  | CVE-2017-16339,cve-arsenal,nvd,cvss-critical |
| 6631 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c0e8 the value for the | nvd | markdown |  | CVE-2017-16340,cve-arsenal,nvd,cvss-critical |
| 6632 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c224 the value for the | nvd | markdown |  | CVE-2017-16341,cve-arsenal,nvd,cvss-critical |
| 6633 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c254 the value for the | nvd | markdown |  | CVE-2017-16342,cve-arsenal,nvd,cvss-critical |
| 6634 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c284 the value for the | nvd | markdown |  | CVE-2017-16343,cve-arsenal,nvd,cvss-critical |
| 6635 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c2c8 the value for the | nvd | markdown |  | CVE-2017-16344,cve-arsenal,nvd,cvss-critical |
| 6636 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c318 the value for the | nvd | markdown |  | CVE-2017-16345,cve-arsenal,nvd,cvss-critical |
| 6637 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c368 the value for the | nvd | markdown |  | CVE-2017-16346,cve-arsenal,nvd,cvss-critical |
| 6638 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01e7d4 the value for the | nvd | markdown |  | CVE-2017-16347,cve-arsenal,nvd,cvss-critical |
| 6639 | A design flaw issue was found in the Red Hat OpenStack Platform director use of TripleO to enable libvirtd based live-migration. Libvirtd is deployed by default | nvd | markdown |  | CVE-2017-2637,cve-arsenal,nvd,cvss-critical |
| 6640 | Vulnerability in the Primavera P6 Enterprise Project Portfolio Management component of Oracle Primavera Products Suite (subcomponent: Web Access (Apache Commons | nvd | markdown |  | CVE-2017-3503,cve-arsenal,nvd,cvss-critical |
| 6641 | Vulnerability in the Oracle Identity Manager component of Oracle Fusion Middleware (subcomponent: Rules Engine). The supported version that is affected is 11.1 | nvd | markdown |  | CVE-2017-3553,cve-arsenal,nvd,cvss-critical |
| 6642 | The drag-and-drop (DnD) function in VMware Workstation 12.x before version 12.5.4 and Fusion 8.x before version 8.5.5 has an out-of-bounds memory access vulnera | nvd | markdown |  | CVE-2017-4901,cve-arsenal,nvd,cvss-critical |
| 6643 | The WHMCS Reseller Module V2 2.0.2 in Softaculous Virtualizor before 2.9.1 does not verify the user correctly, which allows remote authenticated users to cont | nvd | markdown |  | CVE-2017-6513,cve-arsenal,nvd,cvss-critical |
| 6644 | NfSen before 1.3.8 allows remote attackers to execute arbitrary OS commands via shell metacharacters in the customfmt parameter (aka the "Custom output format"  | nvd | markdown |  | CVE-2017-7175,cve-arsenal,nvd,cvss-critical |
| 6645 | TP-Link C2 and C20i devices through firmware 0.9.1 4.2 v0032.0 Build 160706 Rel.37961n allow remote code execution with a single HTTP request by placing shell c | nvd | markdown |  | CVE-2017-8220,cve-arsenal,nvd,cvss-critical |
| 6646 | A vulnerability in the role-based resource checking functionality of the Cisco Unified Computing System (UCS) Director could allow an authenticated, remote atta | nvd | markdown |  | CVE-2018-0238,cve-arsenal,nvd,cvss-critical |
| 6647 | An issue was discovered in MyBiz MyProcureNet 5.0.0. A malicious file can be uploaded to the webserver by an attacker. It is possible for an attacker to upload  | nvd | markdown |  | CVE-2018-11091,cve-arsenal,nvd,cvss-critical |
| 6648 | An issue was discovered in Xen 4.7 through 4.10.x. libxl fails to pass the readonly flag to qemu when setting up a SCSI disk, due to what was probably an errone | nvd | markdown |  | CVE-2018-12892,cve-arsenal,nvd,cvss-critical |
| 6649 | In OnlineJudge 2.0, the sandbox has an incorrect access control vulnerability that can write a file anywhere. A user can write a directory listing to /tmp, and  | nvd | markdown |  | CVE-2018-16367,cve-arsenal,nvd,cvss-critical |
| 6650 | A sandbox escape issue was discovered in VyOS 1.1.8. It provides a restricted management shell for operator users to administer the device. By issuing various s | nvd | markdown |  | CVE-2018-18555,cve-arsenal,nvd,cvss-critical |
| 6651 | A privilege escalation issue was discovered in VyOS 1.1.8. The default configuration also allows operator users to execute the pppd binary with elevated (sudo)  | nvd | markdown |  | CVE-2018-18556,cve-arsenal,nvd,cvss-critical |
| 6652 | A vulnerability was discovered in the Java VM component of Oracle Database Server. Supported versions that are affected are 11.2.0, 12.1.0, 12.2.0 and 18 | nvd | markdown |  | CVE-2018-3110,cve-arsenal,nvd,cvss-critical |
| 6653 | An exploitable vulnerability exists in the smart cameras RTSP configuration of the Samsung SmartThings Hub STH-ETH-250 - Firmware version 0.20.17. The device in | nvd | markdown |  | CVE-2018-3856,cve-arsenal,nvd,cvss-critical |
| 6654 | On Samsung SmartThings Hub STH-ETH-250 devices with firmware version 0.20.17, the video-core process incorrectly extracts fields from a user-controlled JSON pay | nvd | markdown |  | CVE-2018-3863,cve-arsenal,nvd,cvss-critical |
| 6655 | An exploitable buffer overflow vulnerability exists in the samsungWifiScan handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250 - Firmware | nvd | markdown |  | CVE-2018-3866,cve-arsenal,nvd,cvss-critical |
| 6656 | An exploitable stack-based buffer overflow vulnerability exists in the samsungWifiScan callback notification of video-core's HTTP server of Samsung SmartThings  | nvd | markdown |  | CVE-2018-3867,cve-arsenal,nvd,cvss-critical |
| 6657 | An exploitable buffer overflow vulnerability exists in the credentials handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250 - Firmware ver | nvd | markdown |  | CVE-2018-3872,cve-arsenal,nvd,cvss-critical |
| 6658 | An exploitable buffer overflow vulnerability exists in the credentials handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250-Firmware versi | nvd | markdown |  | CVE-2018-3873,cve-arsenal,nvd,cvss-critical |
| 6659 | An exploitable buffer overflow vulnerability exists in the credentials handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250-Firmware versi | nvd | markdown |  | CVE-2018-3874,cve-arsenal,nvd,cvss-critical |
| 6660 | An exploitable buffer overflow vulnerability exists in the credentials handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250-Firmware versi | nvd | markdown |  | CVE-2018-3875,cve-arsenal,nvd,cvss-critical |
| 6661 | An exploitable buffer overflow vulnerability exists in the credentials handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250-Firmware versi | nvd | markdown |  | CVE-2018-3877,cve-arsenal,nvd,cvss-critical |
| 6662 | Multiple exploitable buffer overflow vulnerabilities exist in the credentials handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250 devices | nvd | markdown |  | CVE-2018-3878,cve-arsenal,nvd,cvss-critical |
| 6663 | An exploitable stack-based buffer overflow vulnerability exists in the database 'find-by-cameraId' functionality of video-core's HTTP server of Samsung SmartThi | nvd | markdown |  | CVE-2018-3880,cve-arsenal,nvd,cvss-critical |
| 6664 | An exploitable buffer overflow vulnerability exists in the camera "replace" feature of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250 devices w | nvd | markdown |  | CVE-2018-3902,cve-arsenal,nvd,cvss-critical |
| 6665 | On Samsung SmartThings Hub STH-ETH-250 devices with firmware version 0.20.17, the video-core process incorrectly extracts fields from a user-controlled JSON pay | nvd | markdown |  | CVE-2018-3903,cve-arsenal,nvd,cvss-critical |
| 6666 | An exploitable buffer overflow vulnerability exists in the camera 'update' feature of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250 - Firmware | nvd | markdown |  | CVE-2018-3904,cve-arsenal,nvd,cvss-critical |
| 6667 | An exploitable buffer overflow vulnerability exists in the camera "create" feature of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250 devices wi | nvd | markdown |  | CVE-2018-3905,cve-arsenal,nvd,cvss-critical |
| 6668 | On Samsung SmartThings Hub STH-ETH-250 devices with firmware version 0.20.17, the video-core process insecurely extracts the fields from the "shard" table of it | nvd | markdown |  | CVE-2018-3917,cve-arsenal,nvd,cvss-critical |
| 6669 | An exploitable stack-based buffer overflow vulnerability exists in the retrieval of database fields in video-core's HTTP server of Samsung SmartThings Hub STH-E | nvd | markdown |  | CVE-2018-3919,cve-arsenal,nvd,cvss-critical |
| 6670 | In browser editing in Atlassian Bitbucket Server from version 4.13.0 before 5.4.8 (the fixed version for 4.13.0 through 5.4.7), 5.5.0 before 5.5.8 (the fixed ve | nvd | markdown |  | CVE-2018-5225,cve-arsenal,nvd,cvss-critical |
| 6671 | The default permissions of /dev/kmem in Linux versions before 2.0.36 allows IP spoofing | nvd | markdown |  | CVE-1999-0426,cve-arsenal,nvd,cvss-critical |
| 6672 | VAXstations running Open VMS 5.3 through 5.5-2 with VMS DECwindows or MOTIF do not properly disable access to user accounts that exceed the break-in limit thres | nvd | markdown |  | CVE-1999-1324,cve-arsenal,nvd,cvss-critical |
| 6673 | CGI Script Center News Update 1.1 does not properly validate the original news administration password during a password change operation, which allows remote a | nvd | markdown |  | CVE-2000-0944,cve-arsenal,nvd,cvss-critical |
| 6674 | The default configuration for the domain name resolver for Microsoft Windows 98, NT 4.0, 2000, and XP sets the QueryIpMatching parameter to 0, which causes Wind | nvd | markdown |  | CVE-2000-1218,cve-arsenal,nvd,cvss-critical |
| 6675 | Buffer overflow in FTP server in HPUX 11 allows remote attackers to execute arbitrary commands by creating a long pathname and calling the STAT command, which u | nvd | markdown |  | CVE-2001-0248,cve-arsenal,nvd,cvss-critical |
| 6676 | Heap overflow in FTP daemon in Solaris 8 allows remote attackers to execute arbitrary commands by creating a long pathname and calling the LIST command, which u | nvd | markdown |  | CVE-2001-0249,cve-arsenal,nvd,cvss-critical |
| 6677 | Lightwave ConsoleServer 3200 does not disconnect users after unsuccessful login attempts, which could allow remote attackers to conduct brute force password gue | nvd | markdown |  | CVE-2001-0395,cve-arsenal,nvd,cvss-critical |
| 6678 | Format string vulnerability in Infodrom cfingerd 1.4.3 and earlier allows a remote attacker to gain additional privileges via a malformed ident reply that is pa | nvd | markdown |  | CVE-2001-0609,cve-arsenal,nvd,cvss-critical |
| 6679 | Apache on MacOS X Client 10.0.3 with the HFS+ file system allows remote attackers to bypass access restrictions via a URL that contains some characters whose ca | nvd | markdown |  | CVE-2001-0766,cve-arsenal,nvd,cvss-critical |
| 6680 | Knox Arkeia server 4.2, and possibly other versions, uses a constant salt when encrypting passwords using the crypt() function, which makes it easier for an att | nvd | markdown |  | CVE-2001-0967,cve-arsenal,nvd,cvss-critical |
| 6681 | Symantec LiveUpdate before 1.6 does not use cryptography to ensure the integrity of download files, which allows remote attackers to execute arbitrary code via  | nvd | markdown |  | CVE-2001-1125,cve-arsenal,nvd,cvss-critical |
