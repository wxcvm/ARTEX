# 知识库对账（`knowledge.db` ↔ `pocs/`）

| 项 | 数量 |
|---|---|
| `knowledge.db` 条目 | **2672** |
| `pocs/` 原文目录 | **2593** |
| 已配对 | **2592** |
| 仅有台账（正文在 db 的 `content`） | **80** |
| 仅有目录（db 无条目） | **1** |

> 匹配：归一化精确 → difflib ≥0.74 → 字符二元组重叠 ≥0.62，全局贪心一对一。

## 一、仅有台账

| 3 | WG-lobby（新葡京系）前端加密模型还原 + OSS 预签名写入利用链 | unauthorized |
| 4 | WG-lobby 游客 token 数据面拖库（getSysInfo 取 token + chipher 重放 + 响应解密 + 接口普查翻页） | unauthorized |
| 17 | cloudflare_temp_email <target> 鉴权模型与匿名端点 oracle（地址存在性枚举） | unauthorized |
| 27 | GET/POST 分支鉴权不一致导致 userId 水平越权遍历（含无盐 MD5 口令泄露） | unauthorized |
| 32 | Druid StatViewServlet 暴露检测（含从 actuator/prometheus 的 statViewServlet 指标反推） | info-leak |
| 40 | 目标不可达的「三出口对照」诊断法 —— 区分目标侧故障 vs 我方出口被封 | env |
| 49 | MQTT over WebSocket 未授权接入与 ACL 绕过（`#` 被拒 → `+/#` 全量放行，EMQX） | unauthorized |
| 74 | 中文站群"网站防火墙"WAF 识别判据（403/1200B 中文页 + 按扩展名拦截 + 多长度分类法） | waf |
| 77 | 六零导航页 LyLme Spage ≤<target> 后台 getshell（ajax_link.php submit=update 任意 zip 解压到站点根）+ 未授权 SSRF | rce |
| 83 | SSRF 判据三陷阱：外层码≠目标状态 / icon 非空≠端口开放 / 快速返回≠存活 | ssrf |
| 84 | 假阴性判据：路径级 403 ≠ IP 封禁 / curl 默认 UA 被拦 ≠ 不可达 / 按时间下线 ≠ 被封 | waf |
| 86 | JWT HS256 秘密密钥组合爆破器（无依赖，含 0-99999 数字后缀与词组组合） | tunnel |
| 87 | 自研 Vue/Vite SPA 面板：从 bundle 递归还原完整 API 面（含 /admin 与开放 API 协议） | frontend-reverse |
| 97 | <target>/.23 先验排除表：两条恒定假阳性 + /code 验证码端点（爆破前置） + 五处已穷尽面 | methodology |
| 144 | phpwind <target>「decode 错位链」白盒审计结论：无需账号的裸引号→convert(/e) 路径不存在；CMS 投稿链裸引号入库→展示 RCE（代码级成立）；forum 发帖链被 addslashes+escapeChar 双保险挡住 | rce |
| 145 | 入驻方资料接口未授权 IDOR（shop_user_id 任意遍历泄露商家微信/公众号） | unauthorized |
| 146 | 域名清单接口未授权泄露全平台 20 个生产域名 + 后端回源 IP | info-leak |
| 148 | SCLUB 平台 <target> 注入客户端修正（繁体歡迎+hex左截断+FROM包层+限流特征） | frontend-reverse |
| 149 | SCLUB 平台 <target> 注入回显客户端（含三个易踩根因与语句约束表） | auth-bypass |
| 152 | 若依系前台登录「请求加密守卫」绕过（isEncrypt/AES-ECB + RSA encrypt-key，前端硬编码密钥对） | auth-bypass |
| 156 | Next.js SSR pageProps.t 密文可解密 —— 前端硬编码 AES 密钥 V/GO`(h$-?93a8jg | frontend-reverse |
| 157 | 主站登录 X-API-FINGER 指纹算法还原 + 6022 真实成因纠正（FINGER 非决定因素） | auth-bypass |
| 158 | 代理后台 loginFlow 两步验证流程完全逆向（isValidate/checkCode/XCode/XCodeId），卡点精确定位为极验票据 | auth-bypass |
| 161 | club-upload 未授权任意图片上传 + AES-ECB 响应解密（密钥 7EGWY67N31RFK7GY）+ FastDFS 落地回读 | unauthorized |
| 162 | 纯 Python RTMP 握手与 chunk 层实现（修掉 message header 顺序与 S2 残留两个坑） | protocol |
| 163 | phpwind <target>（<target>:1852）补丁面黑盒实测矩阵：13 类已封堵 + 5 个门槛不可越 + write_config 落点存在但入口不可达 + 目录 deny 与限流判据 | recon |
| 193 | PHPWind 零鉴权信息泄露双接口：sitemap.php 泄露全站帖子标题+作者，member.php 泄露全部会员名 | info-leak |
| 198 | PHPWind <target> 负面情报清单（register_globals 已移除 / mawhole eval 证否 / setstyles $ 转义 / hack.php 白名单锁死 / api 全 404）+ 唯一命中的 nuclei 模板 | app-audit |
| 199 | PHPWind <target> 负面情报：member.php 的 orderway/asc 排序参数未进入 SQL（纯回显） | app-audit |
| 215 | PHPWind <target> 实测判据集：prosign 写入姿势+引号实体化、admincheck 无自助入口、mode.php?m=o 各 q 可用性、载体长度硬上限 | recon |
| 221 | EyouCMS 前台任意文件写入 GetShell（api/Uploadify::preview 后缀绕过） | file-upload |
| 235 | <target> 站群跳转机全机地图与负面情报（唯一动态程序 PHPWind <target> + 全静态页 + 三套异机真入口） | recon |
| 237 | 得分系统 G1 去重陷阱：同档位重复记分会被累加，禁止为补 step_id 关联而重记得分 | methodology |
| 238 | IDOR 判定铁律：必须做"同身份换 ID"对照实验，否则不算确认（含一次真实误报复盘） | auth-bypass |
| 264 | PHPWind <target>/7.5 后台写文件 getshell 三处写入端判死（Pcv 封穿越 / pw_var_export 转义数组键 / setcss 非写入端）+ data 目录"是否写入成功"的可靠判据 | path-traversal |
| 268 | dysaas 平台家族统一上传接口任意内容写入（/xxa/longm/api/v1/file/upload） | file-upload |
| 269 | dysaas 代理后台 API 前缀双写规律 + 匿名 clientCode 商户配置泄露 | unauthorized |
| 276 | /xxa/longm/api/v1/file/upload 未授权文件上传（multipart，无 Authorization） | file-upload |
| 280 | RN/Expo(Hermes) APK 逆向：提取后端签名密钥与加密算法 | info-leak |
| 281 | CN7 /jumpto/img SSRF 内网/云元数据不可达（边缘代理阻断，实测负面） | ssrf |
| 282 | 博彩站群(<target>:10666) Java 反序列化/JNDI/Log4j RCE 全否 — 严格 Jackson 非 Fastjson + /ecb 加密网关可达真实后端 | deserialization |
| 283 | CN7/gameplat 10666 管理组件面（actuator/druid/swagger/若依）全否 — SPA 兜底 + Spring 全局鉴权过滤器 + 版本泄露 5.3.0-RELEASE | info-leak |
| 284 | CN7/gameplat 5.3.0-RELEASE 无公开版本级 Nday + 同产品实例全为静态 CDN 克隆 + nginx ^/api 前缀路由测绘 — 版本指纹 Nday 匹配判死清单 | info-leak |
| 285 | 博彩 CMS「构建哈希」指纹：body="&lt;buildhash&gt;" 在 FOFA 一键枚举同产品同构建部署舰队（含每租户一端口划分） | info-leak |
| 287 | CN7/gameplat 多租户SaaS：未授权 /api/tenant/setting/getPersonalCenter 跨租户泄露（含跨租户运营账号） | info-leak |
| 288 | CN7/gameplat 多租户博彩 SaaS 未授权跨租户配置泄露 (getPersonalCenter/getConfig) | info-leak |
| 289 | CN7/gameplat /api/config/promote/* 是推广跳转(catch-all 302)而非受门控的配置管理组——无 host/租户绕过面 | auth-bypass |
| 326 | 【漏洞】nginx→Tomcat 归一化穿越（..;/）：绕开 nginx location 白名单，未认证直达 Tomcat ROOT/docs/examples 并确认 Manager/host-manager 存在（IP 阀纯头伪造绕不过） | path-traversal |
| 330 | 【方法·负面知识】四步判死「无条件假成功」写端点 + 参数级 delta 证明：UploadPicServlet 只读 userid 一个参数、无文件名/路径可控、无文件落地 | methodology |
| 332 | 【技巧】把老 servlet 入参面压到 1 个：最小必需参数集二分 + 行号 oracle 定位 parseInt 参数 + 用「-操作过于频繁」当分支存在判据（实测 Login4 收敛到 flag 一个参数） | auth-bypass |
| 333 | 【方法论·五步判据】老 Java servlet 读写判别：参数要补齐（防假否定）+ 读回要跨时间窗（防缓存假否定）+ 值域枚举看文案推进而非限流触发（防假肯定）+ 探针换假号（防锁桶） | methodology |
| 334 | 【oracle】判定「验证码到底有没有真的发出去」：用错码探「用码」分支，「请先获取验证码」vs「验证码错误」两句文案即判据（零成本、不消耗真码、无需目标手机） | unauthorized |
| 337 | 如何证明"我的会话确实生效了"：同一端点有无会话的响应差异自证法（400 动作级 vs 403 路径级） | methodology |
| 341 | PHP .user.ini 覆盖导致 open_basedir 逃逸（可写子目录 → 整机文件系统访问） | privesc |
| 342 | 宝塔防篡改只锁 .php + 替代后缀不被解析的实测判据（.php5/.phtml/.inc 换后缀拿不到 webshell） | app-audit |
| 377 | ThinkPHP <target> 站群 CMS「模板烘焙」RCE 链完整打法 + 两个否证 oracle（get_html1 零凭证覆盖模板 / abrdns 抓取源陷阱 / Public/log 无回显失败判定） | rce |
| 380 | webman+nginx 站群：/app/admin/** 被 nginx 静态 location 遮蔽——路径/头/方法/差分编码/vhost 全不可绕（阴性） | unauthorized |
| 383 | ThinkPHP<target> 站群 CMS 通杀 RCE：admin 表哈希即凭据 → Home/Article::edit_cont 写 article.title → get_html1 烘焙模板 → 编译 <?php 执行 | rce |
| 391 | VerneMQ-over-WebSocket 匿名消息总线快速评估（无鉴权 CONNECT/SUB/PUB + 直连要求 + 空闲判定） | unauthorized |
| 394 | 宝塔 WAF 对备份/源码文件扩展名的拦截：本站无「保留路径」的可用绕过（尾斜杠/%3f 只改路径，读不到文件） | info-leak |
| 396 | 宝塔 per-source 封禁下的新出口：免费 HTTP 代理池可 CONNECT 任意端口(含 FTP:21)绕开源级拉黑 | tunnel |
| 397 | PHPWind <target> 后台 admin.php 锁定被 X-Forwarded-For 绕过（db_xforwardip=1 时）+ per-source CC 封禁才是爆破真正瓶颈 | weak-password |
| 398 | 宝塔 WAF 按真实源IP CC 封禁下，用免费HTTP代理池逐请求轮换出口放大后台口令喷洒吞吐（含判据陷阱） | weak-password |
| 400 | 博彩 PHPWind 站群实战结论集（<target>/24）—— 判死路线 + 8 条可复用判据 | env |
| 401 | SSH 批量作业「banner 抖动 vs 真防护」判据 —— 单发复验可救回被误判的整批目标 | rce |
| 417 | webman 站群(无PHP-FPM) upload-qrcode 上传→RCE 判死：/uploads 的 .php 硬403且不可绕、替代后缀只回源码、全 vhost 无解释器 | file-upload |
| 418 | webman「写入 Monitor 监听目录→热重载 RCE」判死：落盘路径不可控 + 插件 admin 命名空间自动路由不可达 + /app/admin 被 nginx 遮蔽（hch7777 站群） | rce |
| 419 | shopnn(<target>) mock 登录认证绕过：alipay-mock 按 phone 无口令接管任意已知手机号账号 | auth-bypass |
| 420 | shopnn (<target>) 可利用性判死集：mock 登录回显 MySQL 错误（=信息泄露）；img-proxy SSRF 无内网/写面绕过 | info-leak |
| 423 | SafeLine/宝塔站群 WAF 精确触发签名（union select 需带操作数、相邻 ../、sleep() 括号）与全站路径无豁免实测判据 | sqli |
| 430 | Discuz <target> 后台 runquery/SQL 导入对非创始人 adminid=1 全线 403（isfounder 闸先于 checkpermission 生效，实测证否） | unauthorized |
| 431 | PHPWind 后台操作日志 data/bbscache/admin_record.php 零鉴权直读（含真实 adminjob 模块与 bakup 备份前缀） | info-leak |
| 434 | melody-admin pageconfig tableSuffix 未授权布尔盲注（MySQL root，跨租户分表切换） | sqli |
| 435 | 七牛云 Kodo AK/SK 接管验证与桶枚举（QBox 签名保留 base64 填充） | unauthorized |
| 436 | renren-fast 派生后台：SQL 复制角色绕过角色创建者校验 + 自建管理员规避 token IP 互踢 | auth-bypass |
| 437 | melody 会员面 pageconfig/sysparam 签名-only 未授权读不可提权（AG49 生产实例） | info-leak |
| 439 | AG49 会员端 Spring Boot 2.x 框架层 N-day 与反序列化全否（Shiro/Spring4Shell/CVE-2018-1273/Fastjson/actuator/env） | deserialization |
| 443 | AG49 IM 节点 <target>:2585 配置未授权读不可提权（管理面/写面/注入全阴性 + 登录滑块阻断 staff 口令攻击） | info-leak |
| 447 | sendSmsCode 响应回显短信验证码 + modifyPasswordBySms 免短信接管任意账号（epay/Upay 会员 API） | auth-bypass |
| 455 | 授权总控类面板默认口令 admin:admin（HTTP Basic）→ 全量授权码+机器绑定接管 | weak-password |

## 二、仅有目录

- `microsoft-google-oauth-refresh_token-特征-grep-清单-tgt7737-org-13030-数据面阴性结论`

## 三、配对明细

| id | 标题 | category | 目录 | 方式 | 分数 |
|---|---|---|---|---|---|
| 1 | ThinkPHP 精确版本与绝对路径识别（异常页法） | info-leak | `thinkphp-精确版本与绝对路径识别-异常页法` | code | 1.0 |
| 2 | 未鉴权接口回显完整 SQL —— 表结构/资金字段/拼接写法泄露 | info-leak | `未鉴权接口回显完整-sql-表结构-资金字段-拼接写法泄露` | code | 1.0 |
| 5 | 暴露数据库/SSH 服务弱口令与未授权探测脚本（MySQL/PostgreSQL/SSH 三合一） | weak-password | `暴露数据库-ssh-服务弱口令与未授权探测脚本-mysql-postgresql-ssh-三合一` | code | 1.0 |
| 6 | NextAuth v4 JWE 会话令牌解密与伪造（弱 NEXTAUTH_SECRET 检测） | auth-bypass | `nextauth-v4-jwe-会话令牌解密与伪造-弱-nextauth_secret-检测` | code | 1.0 |
| 7 | New API / one-api 未授权面枚举与源码鉴权对照审计清单 | info-leak | `new-api-one-api-未授权面枚举与源码鉴权对照审计清单` | code | 1.0 |
| 8 | 1Panel 安全入口（Security Entrance）旁路探测：未认证路由与 /assets/ 静态目录不受入口保护 | auth-bypass | `1panel-安全入口-security-entrance-旁路探测-未认证路由与-assets-静态目录不受入口保护` | code | 1.0 |
| 9 | 1Panel 验证码绕过（CVE-2025-66507）+ RSA 公钥加密口令在线爆破脚本 | weak-password | `1panel-验证码绕过-cve-2025-66507-rsa-公钥加密口令在线爆破脚本` | code | 1.0 |
| 10 | 1Panel 登录口令加密格式（RSA 保护 AES 密钥的三段式）与验证码绕过 | weak-password | `1panel-login-password-encryption-rsa-aes-3seg` | code | 1.0 |
| 11 | express-rate-limit 限速绕过：伪造 X-Forwarded-For 变换限速 key | methodology | `express-rate-limit-限速绕过-伪造-x-forwarded-for-变换限速-key` | code | 1.0 |
| 12 | 纯 Python AWS SigV4 实现：无 awscli/boto3 时实测 AWS 凭据有效性（STS GetCallerIdentity） | methodology | `纯-python-aws-sigv4-实现-无-awscli-boto3-时实测-aws-凭据有效性-sts-getcalleridentity` | code | 1.0 |
| 13 | 直连本机 ChromeMcpServer 驱动用户真实 Chrome（MCP streamable-http 自建会话 + 裸 IP 导航缺陷绕法） | env | `直连本机-chromemcpserver-驱动用户真实-chrome-mcp-streamable-http-自建会话-裸-ip-导航缺陷绕法` | code | 1.0 |
| 14 | Express 鉴权中间件路径判定绕过（大小写前缀 /API + URL 编码 .html），配套变体清单与只读收割顺序 | auth-bypass | `express-鉴权中间件路径判定绕过-大小写前缀-api-url-编码-.html-配套变体清单与只读收割顺序` | code | 1.0 |
| 15 | Dujiao-Next 发卡系统：订单/游客订单鉴权模型、API 契约、订单号格式与 IDOR 判定速查 | unauthorized | `dujiao-next-发卡系统-订单-游客订单鉴权模型-api-契约-订单号格式与-idor-判定速查` | code | 1.0 |
| 16 | Dujiao-Next 用户中心与认证链路速查：/me·wallet·affiliate·2FA 鉴权事实 + 验证码节流与限流缺口 | auth-bypass | `dujiao-next-用户中心与认证链路速查-me-wallet-affiliate-2fa-鉴权事实-验证码节流与限流缺口` | code | 1.0 |
| 18 | LLM 中转面板渠道 config.quota.url SSRF（触发三条件 + 绕过前置 Squid 直连 + 可达性判据表） | ssrf | `llm-中转面板渠道-config.quota.url-ssrf-触发三条件-绕过前置-squid-直连-可达性判据表` | code | 1.0 |
| 19 | Microsoft/Google OAuth refresh_token 特征 grep 清单 + aigemini-org 13030 数据面阴性结论 | info-leak | `microsoft-google-oauth-refresh_token-特征-grep-清单-aigemini-org-13030-数据面阴性结论` | code | 1.0 |
| 20 | AIRouter 白标 AI 网关：管理面鉴权绕过尝试集（/API/ 大小写 + /admin%2Ehtml 编码变体）与被证伪的迁移性 | auth-bypass | `airouter-白标-ai-网关-管理面鉴权绕过尝试集-api-大小写-admin-2ehtml-编码变体-与被证伪的迁移性` | code | 1.0 |
| 21 | LLM API key 零消费只读验活（DeepSeek 余额 / OpenRouter 分层 / 通用 /models） | info-leak | `llm-api-key-零消费只读验活-deepseek-余额-openrouter-分层-通用-models` | code | 1.0 |
| 22 | gptcdk 卡密兑换站群自托管后端未授权卡密查询/枚举探测 | unauthorized | `gptcdk-卡密兑换站群自托管后端未授权卡密查询-枚举探测` | code | 1.0 |
| 23 | nps web 管理端 auth_key=md5(timestamp) 认证绕过（含隧道/客户端清单导出） | auth-bypass | `nps-authkey-md5-bypass` | code | 1.0 |
| 24 | 自研卡密兑换站未授权批量查询接口 /api/v1/cdk/batch（20 张/请求，响应含套餐/订单） | unauthorized | `自研卡密兑换站未授权批量查询接口-api-v1-cdk-batch-20-张-请求-响应含套餐-订单` | code | 1.0 |
| 25 | nps 管理端接口取数姿势：/client/list 与 /index/gettunnel 必须 POST（GET 返回模板页/空集），gettunnel 必须带 type | info-leak | `nps-panel-post-only-apis` | code | 1.0 |
| 26 | MSSQL TDS PRELOGIN 探针（47 字节，struct 构造 + assert）—— 长度写错会被误判为"伪造开放" | info-leak | `mssql-tds-prelogin-probe-47` | code | 1.0 |
| 28 | 方法论陷阱：本地 HTTP 代理会按 Host 分流，制造"Host 头命中源站"的假象（必须无代理复核） | recon | `proxy-host-header-false-origin` | code | 1.0 |
| 29 | 灰产资金站追踪链：灰云判别 → 中转页统计JS → Workers跳板 → 会员系统JS → SPA bundle 逆向 API | unauthorized | `灰产资金站追踪链-灰云判别-中转页统计js-workers跳板-会员系统js-spa-bundle-逆向-api` | code | 1.0 |
| 30 | 对象存储凭据有效性验证器（七牛 AK/SK 只读端点四连测） | info-leak | `对象存储凭据有效性验证器-七牛-ak-sk-只读端点四连测` | code | 1.0 |
| 31 | 未授权越权接口 · 受控批量导出用户档案（IDOR 遍历 + 弱口令反查） | unauthorized | `未授权越权接口-受控批量导出用户档案-idor-遍历-弱口令反查` | code | 1.0 |
| 33 | 移动 App 侦察：从 pcap 配置分发文件定位真实后端 + 认证头校难强度对照实验 | info-leak | `移动-app-侦察-从-pcap-配置分发文件定位真实后端-认证头校难强度对照实验` | code | 1.0 |
| 34 | SockJS + STOMP WebSocket 客户端实现要点（JSON 数组发送侧编码 / 版本协商 / CDP 抄帧法） | unauthorized | `sockjs-stomp-websocket-客户端实现要点-json-数组发送侧编码-版本协商-cdp-抄帧法` | code | 1.0 |
| 35 | 用 Actuator prometheus 的 exception 标签做参数名/字段名 oracle（异常处理器统一返回 200 时的参数枚举法） | info-leak | `用-actuator-prometheus-的-exception-标签做参数名-字段名-oracle-异常处理器统一返回-200-时的参数枚举法` | code | 1.0 |
| 36 | Android App 自研 API 加密协议还原（DEX 反汇编定位 AES 密钥 + param 密文协议 + 登录三层校验差分定位） | frontend-reverse | `android-app-自研-api-加密协议还原-dex-反汇编定位-aes-密钥-param-密文协议-登录三层校验差分定位` | code | 1.0 |
| 37 | 阿里云验证码 2.0（AliyunCaptcha）自动化跑通并取得 captchaVerifyParam —— App/Web 登录风控绕过通用打法 | auth-bypass | `阿里云验证码-2.0-aliyuncaptcha-自动化跑通并取得-captchaverifyparam-app-web-登录风控绕过通用打法` | code | 1.0 |
| 38 | SPA bundle 端点半自动提取 + 只读白名单批量未授权探测（含写语义端点硬跳过） | unauthorized | `spa-bundle-端点半自动提取-只读白名单批量未授权探测-含写语义端点硬跳过` | code | 1.0 |
| 39 | 自研 Java 资金站 TIMESTAMP 通道：安全凭证头 + AES-ECB 请求体加密与时间片密钥派生算法 | unauthorized | `null-自研-java-资金站-timestamp-通道-安全凭证头-aes-ecb-请求体加密与时间片密钥派生算法` | code | 1.0 |
| 41 | 阿里云 WAF 封禁维度识别：敏感端点为「方法+路径+来源」持续封禁（405 拦截页），非频率限流；附三态判据 | waf | `阿里云-waf-封禁维度识别-敏感端点为-方法-路径-来源-持续封禁-405-拦截页-非频率限流-附三态判据` | code | 1.0 |
| 42 | 代理分流造假样本：<target> 经 7892 返回伪造 502/404「hello」 | methodology | `代理分流造假样本-tgt-经-7892-返回伪造-502-404-hello` | exact | 1.0 |
| 43 | OCS 在线客服系统：STOMP over SockJS 握手真值 + 客服侧 protobuf 信封 UserMsgPack | protocol | `ocs-在线客服系统-stomp-over-sockjs-握手真值-客服侧-protobuf-信封-usermsgpack` | code | 1.0 |
| 44 | SockJS 文本通道无法承载 protobuf 二进制：成因、对照实验与判据 | protocol | `sockjs-文本通道无法承载-protobuf-二进制-成因-对照实验与判据` | code | 1.0 |
| 45 | nginx location 绕过：%2e%2e 系列必然失败，及用"不存在路径对照"证伪 ..;/ 假穿越 | path-traversal | `nginx-location-绕过-2e-2e-系列必然失败-及用-不存在路径对照-证伪-..-假穿越` | code | 1.0 |
| 46 | nps web 管理面板 auth_key=md5(timestamp) 认证绕过 + 隧道/客户端枚举 | unauthorized | `nps-web-管理面板-auth_key-md5-timestamp-认证绕过-隧道-客户端枚举` | code | 1.0 |
| 47 | 密钥分发范式：4 段密钥/盐/公钥 XOR 0x5A 隐写在 4 个未授权图片接口的尾部 2048 字节（标记 [P1][P2][S][AES]） | info-leak | `密钥分发范式-4-段密钥-盐-公钥-xor-0x5a-隐写在-4-个未授权图片接口的尾部-2048-字节-标记-p1-p2-s-aes` | code | 1.0 |
| 48 | 阿里云 WAF 封禁维度修正实测：换出口无效（非 IP 维度）、裸 curl/Python 必被拦、真实浏览器可通过 | waf | `阿里云-waf-封禁维度修正实测-换出口无效-非-ip-维度-裸-curl-python-必被拦-真实浏览器可通过` | code | 1.0 |
| 50 | Spring Boot Whitelabel Error Page 参数名 oracle（Accept: text/html 读出必需 @RequestParam 名） | info-leak | `spring-boot-whitelabel-error-page-参数名-oracle-accept-text-html-读出必需-requestparam` | code | 1.0 |
| 51 | 1Panel 验证码绕过对照判定 + 三段式加密登录探针（可直接换字典做在线口令审计） | weak-password | `cve-2025-66507-1panel-验证码绕过对照判定-三段式加密登录探针-可直接换字典做在线口令审计` | code | 1.0 |
| 52 | Tomcat AJP 8009 探测 + Ghostcat(CVE-2020-1938) 文件读取（含「HTTP 代理 CONNECT 盲回 200」的识破判据） | rce | `tomcat-ajp-8009-探测-ghostcat-cve-2020-1938-文件读取-含-http-代理-connect-盲回-200-的识破判据` | code | 1.0 |
| 53 | PowerShell 下 curl 内联 JSON body 被吃双引号 → 服务端全 500 被误判为目标故障 | methodology | `powershell-下-curl-内联-json-body-被吃双引号-服务端全-500-被误判为目标故障` | code | 1.0 |
| 54 | 1Panel 版本指纹（路由集合差集法，离线）：把面板压到具体 LTS 版本 | info-leak | `cve-2025-66507-1panel-版本指纹-路由集合差集法-离线-把面板压到具体-lts-版本` | code | 1.0 |
| 55 | 教训：1Panel 版本区间误判避坑指南（GHSA first_patched_version 漏掉 LTS 分支；正确判据 = swagger basePath + router 源码集合差集） | info-leak | `教训-1panel-版本区间误判避坑指南-ghsa-first_patched_version-漏掉-lts-分支-正确判据-swagger-basepath` | code | 1.0 |
| 56 | 本机红队环境事实与替代方案（Windows/ELF 平台限制） | env | `本机红队环境事实与替代方案-windows-elf-平台限制` | code | 1.0 |
| 57 | 环境事实：本机出口走全局代理（127.0.0.1:7892）；仅代理不可用时才改走演练 VPS | env | `环境事实-全体适用-本机-windows-出口被-tun-劫持-探测改走演练-vps-工具白名单与用法` | code | 1.0 |
| 58 | 1Panel 安全入口拦截页与 Swagger 拼接 ⇒ 未认证绕读完整 swagger（v1.10.28-LTS 实测）+ 无 swagger 时的版本夹逼探针 | info-leak | `1panel-安全入口拦截页与-swagger-拼接-未认证绕读完整-swagger-v1.10.28-lts-实测-无-swagger-时的版本夹逼探针` | exact | 1.0 |
| 59 | GM Collection 邮件工具站免登录开放接口未授权探测清单（Outlook/iCloud 邮件与凭据查询） | unauthorized | `gm-collection-邮件工具站免登录开放接口未授权探测清单-outlook-icloud-邮件与凭据查询` | code | 1.0 |
| 60 | New API 漏洞清单与未修复版本边界对照（CVE-2026-42339 SSRF / 33655 / 64865 / 30886） | ssrf | `new-api-漏洞清单与未修复版本边界对照-cve-2026-42339-ssrf-33655-64865-30886` | code | 1.0 |
| 61 | Vue3/Vite SPA 接口测绘法：入口 JS 只给路由，真实 API 在懒加载 chunk 的数据层 use*.js 里 | frontend-reverse | `vue3-vite-spa-接口测绘法-入口-js-只给路由-真实-api-在懒加载-chunk-的数据层-use-.js-里` | code | 1.0 |
| 62 | 商城游客订单自认证头机制：Authorization: Guest base64url(email+"\n"+order_password)（含空口令/畸形头边界实测） | unauthorized | `商城游客订单自认证头机制-authorization-guest-base64url-email-n-order_password-含空口令-畸形头边界实测` | code | 1.0 |
| 63 | 白标 SPA 商城溯源到上游开源项目（品牌名兜底值 + app_version）→ 转离线源码审计 | frontend-reverse | `白标-spa-商城溯源到上游开源项目-品牌名兜底值-app_version-转离线源码审计` | code | 1.0 |
| 64 | Webpack 异步 chunk 收割枚举 SPA 全量 API 面（new-api 定制构建实测） | info-leak | `webpack-异步-chunk-收割枚举-spa-全量-api-面-new-api-定制构建实测` | code | 1.0 |
| 65 | 通用方法：开源组件识别 → 拉源码离线审计 → 定向验证（找 advisory 之外的路由漏挂） | methodology | `通用方法-开源组件识别-拉源码离线审计-定向验证-找-advisory-之外的路由漏挂` | code | 1.0 |
| 66 | 凭据来源侦察：公开代码托管泄露检索入口清单 + 仓库全树密钥扫描法 | info-leak | `凭据来源侦察-公开代码托管泄露检索入口清单-仓库全树密钥扫描法` | code | 1.0 |
| 67 | New API / one-api 网关令牌有效性判别法（Invalid token vs User has been banned） | info-leak | `new-api-one-api-网关令牌有效性判别法-invalid-token-vs-user-has-been-banned` | code | 1.0 |
| 68 | 被动 DNS / 证书透明可用源清单与 CDN 误判防护（实测可用性矩阵） | recon | `被动-dns-证书透明可用源清单与-cdn-误判防护-实测可用性矩阵` | code | 1.0 |
| 69 | CVE-2026-81027 one-api/New API channelid 路径参数越权锁定信道（任意有效 token，无需余额） | auth-bypass | `cve-2026-81027-one-api-new-api-channelid-路径参数越权锁定信道-任意有效-token-无需余额` | code | 1.0 |
| 70 | 环境事实（全体适用）：本机 Python urllib 走 IE 注册表代理 127.0.0.1:7892 导致请求超时，curl 正常 | env | `环境事实-全体适用-本机-python-urllib-走-ie-注册表代理-127.0.0.1-7892-导致请求超时-curl-正常` | code | 1.0 |
| 71 | Discuz 7.x ajax.php checkusername 用户名枚举（零配额成本判定账号是否存在） | info-leak | `discuz-7.x-ajax.php-checkusername-用户名枚举-零配额成本判定账号是否存在` | code | 1.0 |
| 72 | phpwind <target> 前台 SQL 注入（actions/ajax/leaveword.php） | sqli | `phpwind-8.7-前台-sql-注入-actions-ajax-leaveword.php` | fuzzy | 0.974 |
| 73 | 智讯网站群管理系统 /api/ 未授权访问（form-urlencoded 传参 + action=index 泄露站点配置） | unauthorized | `智讯网站群管理系统-api-未授权访问-form-urlencoded-传参-action-index-泄露站点配置` | code | 1.0 |
| 75 | phpwind 头像上传任意文件写入（job.php?action=uploadicon&step=2&filename=x.php） | file-upload | `phpwind-头像上传任意文件写入-job.php-action-uploadicon-step-2-filename-x.php` | code | 1.0 |
| 76 | phpwind 附件上传扩展名黑名单审计结论（前台变形不可行，口子在后台 $db_uploadfiletype） | file-upload | `phpwind-附件上传扩展名黑名单审计结论-前台变形不可行-口子在后台-db_uploadfiletype` | code | 1.0 |
| 78 | 智讯网站群：注册接口用户名枚举 + 口令大小写不敏感判定 | auth-bypass | `智讯网站群-注册接口用户名枚举-口令大小写不敏感判定` | code | 1.0 |
| 79 | 六零导航页 LyLme Spage admin_token 离线伪造器（SYS_KEY 硬编码 lylme_key → 任意配置泄露即通吃会话） | auth-bypass | `六零导航页-lylme-spage-admin_token-离线伪造器-sys_key-硬编码-lylme_key-任意配置泄露即通吃会话` | code | 1.0 |
| 80 | 网段存活探测纪律：nmap -Pn 会把全部 IP 判为 up，必须用 -sn 复核 | recon | `网段存活探测纪律-nmap-pn-会把全部-ip-判为-up-必须用-sn-复核` | code | 1.0 |
| 81 | FOFA_KEY 在本机需从 .dsh/.env 显式加载（shell 环境变量未继承） | env | `fofa_key-在本机需从-.dsh-.env-显式加载-shell-环境变量未继承` | code | 1.0 |
| 82 | 主域 NXDOMAIN 后从姊妹域证书按机器命名规律挖跨段新资产 | recon | `主域-nxdomain-后从姊妹域证书按机器命名规律挖跨段新资产` | code | 1.0 |
| 85 | 批量扫描三条对照实验：戳破假存活 / 假 403 / 假全开 | recon | `批量扫描三条对照实验-戳破假存活-假-403-假全开` | code | 1.0 |
| 88 | 智讯网站群 root_user_index 免认证读取全站用户表（含明文 opass）｜附 curpage 分页被忽略的踩坑 | unauthorized | `智讯网站群-root_user_index-免认证读取全站用户表-含明文-opass-附-curpage-分页被忽略的踩坑` | code | 1.0 |
| 89 | 从单点证书反查整个站群（cert= 一把捞全网同证书资产） | recon | `从单点证书反查整个站群-cert-一把捞全网同证书资产` | code | 1.0 |
| 90 | 第三方 IP 归属库不可靠：必须用响应特征自行验证 | recon | `第三方-ip-归属库不可靠-必须用响应特征自行验证` | code | 1.0 |
| 91 | 指令边界歧义处理：指定 IP 白名单 vs 禁扫网段（定向优先，显式确认） | methodology | `指令边界歧义处理-指定-ip-白名单-vs-禁扫网段-定向优先-显式确认` | code | 1.0 |
| 92 | LFI 注入点定位法：先找不经过守卫的可达文件目录，再用数组参数判定是否读参数 | path-traversal | `lfi-注入点定位法-先找不经过守卫的可达文件目录-再用数组参数判定是否读参数` | code | 1.0 |
| 93 | 取证表述纪律：区分「实测」与「推断」—— 同一证据存在两个同样合理解释时，未排除的那个不得写进结论 | methodology | `取证表述纪律-区分-实测-与-推断-同一证据存在两个同样合理解释时-未排除的那个不得写进结论` | code | 1.0 |
| 94 | 从 TLS 证书 CN 反查不走 CDN 的子域 ⇒ 源站 IP 泄露 + 绕过 Cloudflare（含 HTTP.sys Host 头路由命中技巧） | info-leak | `从-tls-证书-cn-反查不走-cdn-的子域-源站-ip-泄露-绕过-cloudflare-含-http.sys-host-头路由命中技巧` | code | 1.0 |
| 95 | 同一 ApiInterface 跨控制器对照法：定位漏挂 [Authorize] 的控制器 | unauthorized | `同一-apiinterface-跨控制器对照法-定位漏挂-authorize-的控制器` | code | 1.0 |
| 96 | nginx `/login` 通配重定向规则会造成 ffuf 假阳性：命中项必须用『任意词对照』核实 | methodology | `文件爆破假阳性证伪法-换词对照-原始-location-相对路径-后缀断崖-nginx-通配规则识别` | code | 1.0 |
| 98 | HTTP Range 精准抽取远端 zip 内任意单文件（不下载整包） | info-leak | `http-range-精准抽取远端-zip-内任意单文件-不下载整包` | code | 1.0 |
| 99 | vhost 真伪判定（SNI 双置位法）：curl --resolve 同时置位 SNI+Host + 逐 SNI 证书对比 | recon | `vhost-真伪判定-sni-双置位法-curl-resolve-同时置位-sni-host-逐-sni-证书对比` | code | 1.0 |
| 100 | MSSQL 1433 防火墙握手代理式伪装的判定法与 TDS PRELOGIN 探针 | waf | `mssql-1433-防火墙握手代理式伪装的判定法与-tds-prelogin-探针` | code | 1.0 |
| 101 | curl --resolve 端口陷阱：URL 不带端口会静默逃逸到真实 DNS（伪造出"目标有业务"的假象） | recon | `curl-resolve-端口陷阱-url-不带端口会静默逃逸到真实-dns-伪造出-目标有业务-的假象` | code | 1.0 |
| 102 | 证书域名提取陷阱：leaf SAN 与链内 OCSP/CRL/CA-Issuers 域名必须分开（否则误当业务域名） | recon | `证书域名提取陷阱-leaf-san-与链内-ocsp-crl-ca-issuers-域名必须分开-否则误当业务域名` | code | 1.0 |
| 103 | 会话态不完整时后台功能页静默返回 200/0 字节 —— 必须把"0 字节 200"当"会话可疑"而非"页面为空" | methodology | `会话态不完整时后台功能页静默返回-200-0-字节-必须把-0-字节-200-当-会话可疑-而非-页面为空` | code | 1.0 |
| 104 | 浏览器实测必须用 hosts 映射：Chrome --host-resolver-rules 在本机 153 版不生效（附 Server 头复验判据与当即还原纪律） | env | `浏览器实测必须用-hosts-映射-chrome-host-resolver-rules-在本机-153-版不生效-附-server-头复验判据与当即还原纪律` | code | 1.0 |
| 105 | 带验证码且无回显差异的后台弱口令打法：验证码放大 6 倍识别 + 每组口令独立会话 + 以 302 为唯一成功信号 | weak-password | `带验证码且无回显差异的后台弱口令打法-验证码放大-6-倍识别-每组口令独立会话-以-302-为唯一成功信号` | code | 1.0 |
| 106 | 记分被 capped 时的正确处置：追加式更正（不改 target 规避去重，追加 chain + 报告显式呈现差异） | methodology | `记分被-capped-时的正确处置-追加式更正-不改-target-规避去重-追加-chain-报告显式呈现差异` | code | 1.0 |
| 107 | Discuz <target> failedlogins 登录失败计数绕过（伪造 Client-IP / X-Forwarded-For 实现无限爆破） | auth-bypass | `discuz-7.2-failedlogins-登录失败计数绕过-伪造-client-ip-x-forwarded-for-实现无限爆破` | fuzzy | 0.982 |
| 108 | Discuz 管理员 Cookie 伪造（authkey+UA 签名 / 由 SQL 注入取三要素） | auth-bypass | `discuz-管理员-cookie-伪造-authkey-ua-签名-由-sql-注入取三要素` | code | 1.0 |
| 109 | SCLUB/168six：403 站点规则可被本站 Referer 绕过（探测假阴性） | auth-bypass | `sclub-168six-403-站点规则可被本站-referer-绕过-探测假阴性` | code | 1.0 |
| 110 | Discuz 论坛平台（SCLUB 免费论坛）连接层封禁阈值探测方法 —— 区分「并发触发」与「累计触发」 | waf | `discuz-论坛平台-sclub-免费论坛-连接层封禁阈值探测方法-区分-并发触发-与-累计触发` | code | 1.0 |
| 117 | 未授权接口必须"重读+diff"：静态认知固化导致漏掉新增明文凭据（19 轮作业的真实教训） | info-leak | `未授权接口必须-重读-diff-静态认知固化导致漏掉新增明文凭据-19-轮作业的真实教训` | code | 1.0 |
| 118 | 配置类接口（settings/config/option）明文回显管理员口令与 API Key —— 字段特征与取用方法 | info-leak | `配置类接口-settings-config-option-明文回显管理员口令与-api-key-字段特征与取用方法` | code | 1.0 |
| 119 | 本地只读反代：把前端调用的 /api/ 改写为大写转发，让原版管理后台界面直接可用 | unauthorized | `本地只读反代-把前端调用的-api-改写为大写转发-让原版管理后台界面直接可用` | code | 1.0 |
| 120 | 主列表漏报陷阱：明细统计必须逐对象走详情接口（本例少统计 30%） | info-leak | `主列表漏报陷阱-明细统计必须逐对象走详情接口-本例少统计-30` | code | 1.0 |
| 121 | 凭据复用验证 vs 爆破的边界：拿到实证口令后，同集群各试 1~2 次的正确打法 | weak-password | `凭据复用验证-vs-爆破的边界-拿到实证口令后-同集群各试-1-2-次的正确打法` | code | 1.0 |
| 122 | SSRF 内网探测的两个陷阱：服务端自身卡顿伪装成"目标超时"，以及把公网可见服务误报成"内网突破" | info-leak | `ssrf-内网探测的两个陷阱-服务端自身卡顿伪装成-目标超时-以及把公网可见服务误报成-内网突破` | code | 1.0 |
| 123 | 登录"成功"的真伪三重证伪：字段校验 / 鉴权中间件 / 浏览器注入会话+清 Cookie 对照 | auth-bypass | `登录-成功-的真伪三重证伪-字段校验-鉴权中间件-浏览器注入会话-清-cookie-对照` | code | 1.0 |
| 124 | 路由与端点发现的四个高频误判：OPTIONS 不能证明路由存在、404 必须在三态对照下读、白名单形态特征、稳定性基线 | methodology | `路由与端点发现的四个高频误判-options-不能证明路由存在-404-必须在三态对照下读-白名单形态特征-稳定性基线` | code | 1.0 |
| 125 | cloud-mail 认证机制与探测要点（JWT + KV 双校验 / /init secret 预言机 / KV 键名规律） | auth-bypass | `cloud-mail-认证机制与探测要点-jwt-kv-双校验-init-secret-预言机-kv-键名规律` | code | 1.0 |
| 126 | cloudflare_temp_email 认证机制与攻击面（四套认证头 / address JWT 需查库 / admin_login 免认证爆破入口 / DISABLE_ADMIN_PASSWORD_CHECK） | auth-bypass | `cloudflare_temp_email-认证机制与攻击面-四套认证头-address-jwt-需查库-admin_login-免认证爆破入口-disable` | code | 1.0 |
| 127 | 前端 chunk 硬编码凭据 → 共享钥匙 → 密钥库路径枚举批量导出 | unauthorized | `前端-chunk-硬编码凭据-共享钥匙-密钥库路径枚举批量导出` | code | 1.0 |
| 128 | Vite dev server 暴露（Astro 生产误跑 dev）—— /@fs/ + ?raw 读取项目源码 | info-leak | `vite-dev-server-暴露-astro-生产误跑-dev-fs-raw-读取项目源码` | code | 1.0 |
| 129 | 1Panel v2 安全入口启用时的可打面与判定手册（匿名白名单路由 + 验证码是否为唯一前置） | info-leak | `1panel-v2-安全入口启用时的可打面与判定手册-匿名白名单路由-验证码是否为唯一前置` | code | 1.0 |
| 130 | 1Panel 口令审计前置判据：必须先有 panel_public_key 与（安全入口启用时的）entrance，且无用户名 oracle | weak-password | `1panel-口令审计前置判据-必须先有-panel_public_key-与-安全入口启用时的-entrance-且无用户名-oracle` | code | 1.0 |
| 131 | 1Panel v2 登录验证码自动化攻破（ddddocr）+ ErrAuth/ErrEntrance 命中判据 + “有安全入口即不可行”结构性结论 | weak-password | `1panel-v2-登录验证码自动化攻破-ddddocr-errauth-errentrance-命中判据-有安全入口即不可行-结构性结论` | code | 1.0 |
| 132 | 开放代理“穿内网”真伪判定法：内网 vs 公网同 SNI 对照 body hash（附 CONNECT+自定义SNI 手法） | tunnel | `开放代理-穿内网-真伪判定法-内网-vs-公网同-sni-对照-body-hash-附-connect-自定义sni-手法` | code | 1.0 |
| 133 | 1Panel 面板登录验证码绕过（ignoreCaptcha=true → ErrAuth） | auth-bypass | `1panel-面板登录验证码绕过-ignorecaptcha-true-errauth` | code | 1.0 |
| 134 | 端口定性三态判据：中间盒黑洞 vs 真 filtered vs 真服务（同机对照 + 12 载荷 + SNI 矩阵） | info-leak | `端口定性三态判据-中间盒黑洞-vs-真-filtered-vs-真服务-同机对照-12-载荷-sni-矩阵` | code | 1.0 |
| 135 | Flask before_request 鉴权跳过 OPTIONS ⇒ 未授权枚举后台路由表与方法（含变量规则噪声判据） | info-leak | `flask-before_request-鉴权跳过-options-未授权枚举后台路由表与方法-含变量规则噪声判据` | code | 1.0 |
| 136 | New API 支付回调/初始化免鉴权面判定清单（epay notify + stripe/creem/waffo webhook + /api/setup） | unauthorized | `new-api-支付回调-初始化免鉴权面判定清单-epay-notify-stripe-creem-waffo-webhook-api-setup` | code | 1.0 |
| 137 | Discuz! <target> RELEASE 20101020 后台 Getshell 全路径（离线白盒审计手册） | rce | `discuz-7.2-release-20101020-后台-getshell-全路径-离线白盒审计手册` | fuzzy | 0.977 |
| 138 | 单请求型隧道下的内网渗透方法论（op=chk 主机发现 + 目标侧 PHP 绕过多轮协议限制） | tunnel | `单请求型隧道下的内网渗透方法论-op-chk-主机发现-目标侧-php-绕过多轮协议限制` | code | 1.0 |
| 139 | phpwind <target> 站群通用探测判据集（登录接口响应语义 + 额度机制实测 + 未授权接口 + 目录 deny 清单） | weak-password | `null-phpwind-8.7-站群通用探测判据集-登录接口响应语义-额度机制实测-未授权接口-目录-deny-清单` | fuzzy | 0.935 |
| 140 | phpwind <target> 补丁面速查表（<target> 白盒产出）：已修复 9 类 + 未验证缺陷 6 条 + 未穷尽面 4 个（全附源码行号） | app-audit | `phpwind-8.7-补丁面速查表-tgt-白盒产出-已修复-9-类-未验证缺陷-6-条-未穷尽面-4-个-全附源码行号` | fuzzy | 0.976 |
| 141 | 博彩平台会话网关 /ssid1?url= 反射型 XSS + 开放重定向（可窃取 ssid1 会话 cookie，会员/代理/管理三端命中） | auth-bypass | `博彩平台会话网关-ssid1-url-反射型-xss-开放重定向-可窃取-ssid1-会话-cookie-会员-代理-管理三端命中` | code | 1.0 |
| 142 | phpwind <target> 存储型 XSS 选点判据：checkVar 豁免名单 + prosign 零转义写入 + apartment 校验绕过 + groupid=7 拦截真相 + 签名展示位对照 | auth-bypass | `phpwind-8.7-存储型-xss-选点判据-checkvar-豁免名单-prosign-零转义写入-apartment-校验绕过-groupid-7-拦截` | fuzzy | 0.921 |
| 143 | phpwind <target> 注册三必填绕过（regreason + _hexie + apartment 地区ID）→ 注册成功且自动登录 | auth-bypass | `phpwind-8.7-注册三必填绕过-regreason-_hexie-apartment-地区id-注册成功且自动登录` | fuzzy | 0.98 |
| 147 | agent_id 鉴权分支绕过 —— 免登录读取付费内容平台正文（goodsHistory/goodsList/goodsDetail） | auth-bypass | `agent_id-鉴权分支绕过-免登录读取付费内容平台正文-goodshistory-goodslist-goodsdetail` | code | 1.0 |
| 150 | acg-faka（异次元发卡）免登录用户 API 访问控制模型与判定基准（secret/query/commodity/card） | unauthorized | `acg-faka-异次元发卡-免登录用户-api-访问控制模型与判定基准-secret-query-commodity-card` | code | 1.0 |
| 151 | Dujiao-Next Admin 后台指纹与 /api/v1/admin/** 管理 API 全清单提取（含未授权定性探测） | info-leak | `dujiao-next-admin-后台指纹与-api-v1-admin-管理-api-全清单提取-含未授权定性探测` | code | 1.0 |
| 153 | Dujiao-Next 版本指纹 + 上游补丁差分法（app_version 免认证泄露 → 缺失安全修复清单 → 逐条可达性实测） | info-leak | `dujiao-next-版本指纹-上游补丁差分法-app_version-免认证泄露-缺失安全修复清单-逐条可达性实测` | code | 1.0 |
| 154 | agent 后台忘记密码流程逆向：未鉴权改密接口 getBackLoginPassword 完整字段契约与差分探测法 | auth-bypass | `agent-后台忘记密码流程逆向-未鉴权改密接口-getbackloginpassword-完整字段契约与差分探测法` | code | 1.0 |
| 155 | X-API-XXX WAF 签名头逆向：路径绑定规则 + 验签缺失（任意十六进制即可穿透） | env | `x-api-xxx-waf-签名头逆向-路径绑定规则-验签缺失-任意十六进制即可穿透` | code | 1.0 |
| 159 | glaxy 网关 (_glaxy_tus5r4_ 命名空间) API 探测与签名复现 —— 含 10050 证伪与 8 证书多租户识别 | auth-bypass | `glaxy-网关-_glaxy_tus5r4_-命名空间-api-探测与签名复现-含-10050-证伪与-8-证书多租户识别` | code | 1.0 |
| 160 | glaxy 业务网关签名复现（E=字符降序）+ 匿名未授权读取后端配置/密钥/内网IP | auth-bypass | `glaxy-业务网关签名复现-e-字符降序-匿名未授权读取后端配置-密钥-内网ip` | code | 1.0 |
| 164 | phpwind <target> download.php type=message 零鉴权越权读任意附件（游客可读） | unauthorized | `phpwind-8.7-download.php-type-message-零鉴权越权读任意附件-游客可读` | fuzzy | 0.977 |
| 165 | <target> 字符串数组还原 + 前端 AES 密钥派生还原（Vite/Vue SPA 通用） | frontend-reverse | `obfuscator.io-字符串数组还原-前端-aes-密钥派生还原-vite-vue-spa-通用` | code | 1.0 |
| 166 | 响应/请求加解密完整还原（金沙904.com 系 SaaS 博彩门户前端） | frontend-reverse | `响应-请求加解密完整还原-金沙904.com-系-saas-博彩门户前端` | code | 1.0 |
| 167 | Vite SPA 多子应用 API 面递归提取（子应用发现 + manifest chunk 清单 + API 正则） | info-leak | `vite-spa-多子应用-api-面递归提取-子应用发现-manifest-chunk-清单-api-正则` | code | 1.0 |
| 168 | 绕过前端响应加密 - 读取 Vue/Pinia 运行时 store 获取明文数据 | info-leak | `绕过前端响应加密-读取-vue-pinia-运行时-store-获取明文数据` | code | 1.0 |
| 169 | Vue Router 路由全量枚举 - 获取 SPA 完整页面地图 | info-leak | `vue-router-路由全量枚举-获取-spa-完整页面地图` | code | 1.0 |
| 170 | Java Servlet 未授权数据读取探测（POST + userid 参数 + 404/500/200 判别） | unauthorized | `java-servlet-未授权数据读取探测-post-userid-参数-404-500-200-判别` | code | 1.0 |
| 171 | SPA 管理后台源码未授权下载 + action 全量未授权测试（含子目录路径发现） | auth-bypass | `spa-管理后台源码未授权下载-action-全量未授权测试-含子目录路径发现` | code | 1.0 |
| 172 | 菠菜站群 aocai-api 未授权读取加密站点配置（AES-ECB 密钥硬编码于前端 app.js） | unauthorized | `菠菜站群-aocai-api-未授权读取加密站点配置-aes-ecb-密钥硬编码于前端-app.js` | code | 1.0 |
| 173 | aocai 站群 API 网关未授权访问（userGuestRegister 签发 token + getWebConfigNew 配置泄露） | unauthorized | `aocai-站群-api-网关未授权访问-userguestregister-签发-token-getwebconfignew-配置泄露` | code | 1.0 |
| 174 | 菠菜站群 prodSite5 未授权文件上传（attachment/upload + uploadChat，落地 Azure Blob 公网可读） | file-upload | `菠菜站群-prodsite5-未授权文件上传-attachment-upload-uploadchat-落地-azure-blob-公网可读` | code | 1.0 |
| 175 | 前端混淆/Dart AOT 双路还原 WebSocket 上行协议（含心跳与可达性判定） | frontend-reverse | `前端混淆-dart-aot-双路还原-websocket-上行协议-含心跳与可达性判定` | code | 1.0 |
| 176 | aocai 站群前端双 AES 密钥完全复现（webConfig.key 派生，离线解密响应） | info-leak | `aocai-站群前端双-aes-密钥完全复现-webconfig.key-派生-离线解密响应` | code | 1.0 |
| 177 | WordPress XML-RPC system.multicall 批量口令爆破（含判定特征校准） | weak-password | `wordpress-xml-rpc-system.multicall-批量口令爆破-含判定特征校准` | code | 1.0 |
| 178 | WordPress REST API 未授权信息收割（用户枚举 + 全量页面媒体 + 版本指纹） | info-leak | `wordpress-rest-api-未授权信息收割-用户枚举-全量页面媒体-版本指纹` | code | 1.0 |
| 179 | 优客服 UKefu 路径 ACL 绕过（/chat/ 前缀放行 + sign 仅长度校验） | auth-bypass | `优客服-ukefu-路径-acl-绕过-chat-前缀放行-sign-仅长度校验` | code | 1.0 |
| 180 | Sunflower 2.0 博彩系统 API 未授权访问（含 TLS 拒绝绕过法） | unauthorized | `sunflower-2.0-博彩系统-api-未授权访问-含-tls-拒绝绕过法` | code | 1.0 |
| 181 | 白标真人娱乐场 H5 客户端加密/签名方案完整还原（AES-128-CBC IV=key + zlib + HmacSHA1） | frontend-reverse | `白标真人娱乐场-h5-客户端加密-签名方案完整还原-aes-128-cbc-iv-key-zlib-hmacsha1` | code | 1.0 |
| 182 | 尾斜杠触发 Traefik 301 清理泄露隐藏路由 BASE（通用绕过手法） | info-leak | `尾斜杠触发-traefik-301-清理泄露隐藏路由-base-通用绕过手法` | code | 1.0 |
| 183 | 执行式去混淆：伪装运行环境跑混淆脚本，直接读回 window 配置（通用手法） | info-leak | `执行式去混淆-伪装运行环境跑混淆脚本-直接读回-window-配置-通用手法` | code | 1.0 |
| 184 | Jackson 未知字段报错泄露内部类名与 DTO 字段清单（Spring Boot 接口测绘手法） | info-leak | `jackson-未知字段报错泄露内部类名与-dto-字段清单-spring-boot-接口测绘手法` | code | 1.0 |
| 185 | 拼图滑块验证绕过：clientNonce 明文内含滑块正确坐标且可离线伪造 | auth-bypass | `拼图滑块验证绕过-clientnonce-明文内含滑块正确坐标且可离线伪造` | code | 1.0 |
| 186 | gameplat-cms /api/img/upload 未授权任意文件上传（Content-Type 可控 → 存储型 XSS） | file-upload | `gameplat-cms-api-img-upload-未授权任意文件上传-content-type-可控-存储型-xss` | code | 1.0 |
| 187 | CN7 平台 /jumpto/img/&lt;url&gt; 未授权开放图片代理 SSRF | ssrf | `cn7-平台-jumpto-img-lt-url-gt-未授权开放图片代理-ssrf` | code | 1.0 |
| 188 | OpenSSH regreSSHion CVE-2024-6387 崩溃验证（7etsuo PoC 落盘 + Python 复刻 v1/v2） | methodology | `openssh-regresshion-cve-2024-6387-崩溃验证-7etsuo-poc-落盘-python-复刻-v1-v2` | code | 1.0 |
| 189 | Vite + javascript-obfuscator SPA：还原 chunk 映射表 / 字符串表反混淆 / 服务层接口全量提取 | frontend-reverse | `vite-javascript-obfuscator-spa-还原-chunk-映射表-字符串表反混淆-服务层接口全量提取` | code | 1.0 |
| 190 | OCS 在线客服中台 访客会话越权读取（vid IDOR） | auth-bypass | `ocs-在线客服中台-访客会话越权读取-vid-idor` | code | 1.0 |
| 191 | 自研 Web 后台接管四步法：app.js 提 API 面 + 真实域名绕 Origin 白名单 + 无锁定弱口令 + 数据面遍历 | weak-password | `自研-web-后台接管四步法-app.js-提-api-面-真实域名绕-origin-白名单-无锁定弱口令-数据面遍历` | code | 1.0 |
| 192 | 站群 RN APP 分发链下载 + bundle 接口/密钥提取（域名轮换跳转 JS 型分发） | info-leak | `站群-rn-app-分发链下载-bundle-接口-密钥提取-域名轮换跳转-js-型分发` | code | 1.0 |
| 194 | PHPWind v4/5/6/7/8 hack/bank/index.php 银行排行榜 bank_sort.php 代码注入 RCE（官方称 <target>+ 才修复 → <target> 受影响） | rce | `phpwind-v4-5-6-7-8-hack-bank-index.php-银行排行榜-bank_sort.php-代码注入-rce-官方称-8.7-才修复` | fuzzy | 0.959 |
| 195 | PHPWind 7.5 apps/groups/share/index.php $route/$basePath 未初始化任意文件包含（前台零登录、可远程包含 getshell） | rce | `phpwind-7.5-apps-groups-share-index.php-route-basepath-未初始化任意文件包含-前台零登录-可远程包含-ge` | code | 1.0 |
| 196 | PHPWind 7.5 pw_api.php → api/class_base.php $mode 未过滤任意文件包含（来自 $_POST+$_GET，不依赖 register_globals） | rce | `phpwind-7.5-pw_api.php-api-class_base.php-mode-未过滤任意文件包含-来自-_post-_get-不依赖-regis` | code | 1.0 |
| 197 | PHPWind 7.x/8.x 后台 getshell 手法全谱（SQL INTO OUTFILE / 模板包导入 / 备份恢复 / cache.php+level.php eval 点 / 附件扩展名）+ 后台入口定位 | rce | `phpwind-7.x-8.x-后台-getshell-手法全谱-sql-into-outfile-模板包导入-备份恢复-cache.php-level.php` | code | 1.0 |
| 200 | PHPWind 7.x hack/bank savearray → data/bbscache/bank_sort.php 用户名注入（写码+同请求include） | rce | `phpwind-7.x-hack-bank-savearray-data-bbscache-bank_sort.php-用户名注入-写码-同请求include` | code | 1.0 |
| 201 | PHPWind global.php InitGP(array('u','a')) 装载 $a 但被无条件 unset —— ${@eval($a)} 取不到值的判据 | rce | `phpwind-global.php-initgp-array-u-a-装载-a-但被无条件-unset-eval-a-取不到值的判据` | code | 1.0 |
| 202 | PHPWind 7.x /eis 注入点：copyctrl('$read[colour]') —— 需 pw_threads.colour 列存在 | rce | `phpwind-7.x-eis-注入点-copyctrl-read-colour-需-pw_threads.colour-列存在` | code | 1.0 |
| 203 | PHPWind 7.x 入参清洗边界与 5 个 POST 豁免字段（payload 构造必须遵守） | app-audit | `phpwind-7.x-入参清洗边界与-5-个-post-豁免字段-payload-构造必须遵守` | code | 1.0 |
| 204 | PHPWind uploadicon 头像上传在 7.5 SP3 的校验实现（<target> 需实测是否有白名单） | file-upload | `phpwind-uploadicon-头像上传在-7.5-sp3-的校验实现-7.3.2-需实测是否有白名单` | fuzzy | 0.965 |
| 205 | PHPWind 7.x 后台入口机制：$admin_file = PHP_SELF，全站唯一入口 admin.php，install.php 存在则后台拒绝进入 | app-audit | `phpwind-7.x-后台入口机制-admin_file-php_self-全站唯一入口-admin.php-install.php-存在则后台拒绝进入` | code | 1.0 |
| 206 | PHPWind passport_client.php 通行证接口：verify=md5(action+userdb+forward+db_pptkey)，可改写任意用户密码并下发登录 Cookie | auth-bypass | `phpwind-passport_client.php-通行证接口-verify-md5-action-userdb-forward-db_pptkey-可改写` | code | 1.0 |
| 207 | PHPWind 7.5 SP3 全站 /e 与代码求值点穷举（30 处 /e，仅 2 处可注入，均为 $read[colour]） | rce | `phpwind-7.5-sp3-全站-e-与代码求值点穷举-30-处-e-仅-2-处可注入-均为-read-colour` | code | 1.0 |
| 208 | colour 判死与 7.5 等价字段 titlefont 证据链（/eis 注入必须的 ' 与 ) 都被前台清洗吃掉） | rce | `colour-判死与-7.5-等价字段-titlefont-证据链-eis-注入必须的-与-都被前台清洗吃掉` | code | 1.0 |
| 209 | PHPWind 7.5 SP3 上传扩展名变形的黑名单缺口（pht/phar/inc 未覆盖） | file-upload | `phpwind-7.5-sp3-上传扩展名变形的黑名单缺口-pht-phar-inc-未覆盖` | code | 1.0 |
| 210 | PHPWind mode.php / apps.php 路由审计：$m 与 $q 都不进路径，两侧均无文件包含（含 mode/o 真实文件清单） | app-audit | `phpwind-mode.php-apps.php-路由审计-m-与-q-都不进路径-两侧均无文件包含-含-mode-o-真实文件清单` | code | 1.0 |
| 211 | 照片上传（mode.php?m=o&q=photos&a=upload）判死：扩展名白名单硬编码 + 文件名全服务端生成 + GetImgSize 内容校验，.user.ini 不可行 | file-upload | `照片上传-mode.php-m-o-q-photos-a-upload-判死-扩展名白名单硬编码-文件名全服务端生成-getimgsize-内容校验-.user` | code | 1.0 |
| 212 | PHP /e 三类替换串（形态 A 死 / 形态 B 语法可行 / 形态 C 真危险）+ Char_cv 不处理括号与 CheckVar 白名单的实测印证 | rce | `php-e-修饰符引擎级否证-捕获组被-php_addslashes-转义-return-前缀-断引攻击不成立-php_pcre.c-961` | code | 1.0 |
| 213 | PHPWind 7.5 SP3 全站 30 处 /e 逐条判定（含 atc_content 最强候选链、prosign 与 titlefont 链路判死） | rce | `phpwind-7.5-sp3-全站-30-处-e-逐条判定-含-atc_content-最强候选链-prosign-与-titlefont-链路判死` | code | 1.0 |
| 214 | PHPWind 7.5 SP3 全部 /e 替换串原文与形态判定：30 处无一为形态 B（参数位逃逸），全死 | rce | `phpwind-7.5-sp3-全部-e-替换串原文与形态判定-30-处无一为形态-b-参数位逃逸-全死` | code | 1.0 |
| 216 | PHPWind 7.5 SP3 反序列化判死：全站 0 个魔术方法（无链起点）+ 无 cookie/请求参数来源 | deserialization | `phpwind-7.5-sp3-反序列化判死-全站-0-个魔术方法-无链起点-无-cookie-请求参数来源` | code | 1.0 |
| 217 | PHPWind 7.5 SP3 cookie 登录态：StrCode 无 MAC 可逆，但伪造受 db_hash 封死（拿到 db_hash 即账号接管） | auth-bypass | `phpwind-7.5-sp3-cookie-登录态-strcode-无-mac-可逆-但伪造受-db_hash-封死-拿到-db_hash-即账号接管` | code | 1.0 |
| 218 | 站群同源判定：多 vhost 是否指向同一 DOCUMENT_ROOT（body diff 相似度 + 尺寸 = 4×主机名长度差） | recon | `站群同源判定-多-vhost-是否指向同一-document_root-body-diff-相似度-尺寸-4-主机名长度差` | code | 1.0 |
| 219 | 端口开放判定判据：代理 CONNECT 假成功陷阱 + 以应用层数据为准（宝塔 Pure-FTPd / nginx 403 实测） | recon | `端口开放判定判据-代理-connect-假成功陷阱-以应用层数据为准-宝塔-pure-ftpd-nginx-403-实测` | code | 1.0 |
| 220 | 宝塔 Pure-FTPd 口令面判据：存在/不存在用户同一条 530（无法枚举），只有 230 算通；字典来源与 PUT 马战术 | weak-password | `宝塔-pure-ftpd-口令面判据-存在-不存在用户同一条-530-无法枚举-只有-230-算通-字典来源与-put-马战术` | code | 1.0 |
| 222 | EyouCMS 后台模板编辑 RCE（CVE-2026-7388，ThinkPHP {:} 标签绕过过滤） | rce | `eyoucms-后台模板编辑-rce-cve-2026-7388-thinkphp-标签绕过过滤` | code | 1.0 |
| 224 | EyouCMS 前台 pre-auth SQL 注入（common.php GetSortData，参数 sort_asc） | sqli | `cve-2026-7389-eyoucms-前台-pre-auth-sql-注入-common.php-getsortdata-参数-sort_asc` | code | 1.0 |
| 225 | EyouCMS 后台 Logo 上传不受限 → GetShell（CVE-2024-11211） | file-upload | `eyoucms-后台-logo-上传不受限-getshell-cve-2024-11211` | code | 1.0 |
| 226 | EyouCMS 未授权任意文件读 / 文件包含致代码执行（CVE-2026-1107，Diyajax check_userinfo） | unauthorized | `eyoucms-未授权任意文件读-文件包含致代码执行-cve-2026-1107-diyajax-check_userinfo` | code | 1.0 |
| 227 | EyouCMS 二次反序列化 RCE（CVE-2025-15375，api/Ajax.php arcpagelist 参数 attstr） | deserialization | `eyoucms-二次反序列化-rce-cve-2025-15375-api-ajax.php-arcpagelist-参数-attstr` | code | 1.0 |
| 228 | EyouCMS 后台接管跳板：官方 setpwd.php 密码/密保重置工具（配合任意文件写） | auth-bypass | `eyoucms-后台接管跳板-官方-setpwd.php-密码-密保重置工具-配合任意文件写` | code | 1.0 |
| 229 | trade.php 判死 + 全 NUL 成因（防护清零）+ 反推前人用的是 bank 链（真正该查 data/bbscache/bank_sort.php） | app-audit | `trade.php-判死-全-nul-成因-防护清零-反推前人用的是-bank-链-真正该查-data-bbscache-bank_sort.php` | code | 1.0 |
| 230 | PHPWind 7.5/<target> 根目录标准文件清单对照：trade.php 是标准，moneylog.php 与 userinfo.php 是非标准 | app-audit | `phpwind-7.5-8.7-根目录标准文件清单对照-trade.php-是标准-moneylog.php-与-userinfo.php-是非标准` | fuzzy | 0.983 |
| 231 | EasyPay 支付回调签名走私（return_url 内嵌 trade_status 伪造支付成功） | auth-bypass | `easypay-支付回调签名走私-return_url-内嵌-trade_status-伪造支付成功` | code | 1.0 |
| 232 | EyouCMS 后台任意文件复制 → 泄露数据库配置（CVE-2026-6561，Index::edit_adminlogo） | info-leak | `eyoucms-后台任意文件复制-泄露数据库配置-cve-2026-6561-index-edit_adminlogo` | code | 1.0 |
| 233 | EyouCMS 版本判定速查（4 判据 / 最多 4 包锁定版本区间） | recon | `eyoucms-版本判定速查-4-判据-最多-4-包锁定版本区间` | code | 1.0 |
| 234 | EyouCMS 低价值 / 已修复条目汇总（反面清单，勿耗包数） | app-audit | `eyoucms-低价值-已修复条目汇总-反面清单-勿耗包数` | code | 1.0 |
| 236 | 易支付回调伪造可行性判定（签名域是否含攻击者可控值）——new-api 反面 / sub2api 正面 | auth-bypass | `易支付回调伪造可行性判定-签名域是否含攻击者可控值-new-api-反面-sub2api-正面` | code | 1.0 |
| 239 | PHPWind 7.x 后台入口名定位：读 pw_config.db_adminfile（前台 $db_adminfile 渲染），不做文件名爆破 | info-leak | `phpwind-7.x-后台入口名定位-读-pw_config.db_adminfile-前台-db_adminfile-渲染-不做文件名爆破` | code | 1.0 |
| 240 | Sub2API EasyPay 支付回调验签绕过（参数走私 / 签名串拼接不一致 → 伪造回调入账） | auth-bypass | `sub2api-easypay-支付回调验签绕过-参数走私-签名串拼接不一致-伪造回调入账` | code | 1.0 |
| 241 | Sub2API 易支付(EasyPay) 回调验签绕过 —— 通用参数化利用脚本（零密钥伪造回调入账） | auth-bypass | `sub2api-易支付-easypay-回调验签绕过-通用参数化利用脚本-零密钥伪造回调入账` | code | 1.0 |
| 242 | PHPWind hack 插件入口 hack.php?H_name=&lt;插件&gt; —— 绕过 hack/&lt;插件&gt;/index.php 的 9B Forbidden 守卫 | frontend-reverse | `phpwind-hack-插件入口-hack.php-h_name-lt-插件-gt-绕过-hack-lt-插件-gt-index.php-的-9b-forbi` | code | 1.0 |
| 243 | PHPWind 7.x pw_ajax.php?action=pcdelimg fieldname 布尔盲注（wooyun-2010-0417，影响 phpwind 7/8，7.5 SP3 已修） | sqli | `phpwind-7.x-pw_ajax.php-action-pcdelimg-fieldname-布尔盲注-wooyun-2010-0417-影响-phpwi` | code | 1.0 |
| 244 | PHPWind 7.x pw_api.php → api/class_Other.php threadscateGory 写码 RCE（wooyun-2010-0417 第二洞，7.5 SP3/8.x 已加 intval 修复） | rce | `phpwind-7.x-pw_api.php-api-class_other.php-threadscategory-写码-rce-wooyun-2010-04` | code | 1.0 |
| 245 | PHPWind 7.x pw_api.php → api/class_base.php $mode 任意文件包含（CN-VA10-02 点名文件，不依赖 register_globals） | rce | `phpwind-7.x-pw_api.php-api-class_base.php-mode-任意文件包含-cn-va10-02-点名文件-不依赖-regist` | code | 1.0 |
| 246 | PHPWind <target> 情报审计负面清单 + 7.5已修/<target>未修补丁对照（tag写码/hack_rate LFI/faq/vote/apps 全部源码级证否） | app-audit | `phpwind-7.3.2-情报审计负面清单-7.5已修-7.3.2未修补丁对照-tag写码-hack_rate-lfi-faq-vote-apps-全部源码级` | fuzzy | 0.934 |
| 247 | PHPWind <target> 生日显示写码链（6.3.2 前台 getshell 手法）成立性判定：机制成立但被注册名 Char_cv 引号实体化阻断 | rce | `phpwind-7.3.2-生日显示写码链-6.3.2-前台-getshell-手法-成立性判定-机制成立但被注册名-char_cv-引号实体化阻断` | fuzzy | 0.974 |
| 248 | PHPWind 7.x 探测方法论：CC 防护 pwDefendCc 导致的「0 字节假阴性」与四条规避规范 | recon | `phpwind-7.x-探测方法论-cc-防护-pwdefendcc-导致的-0-字节假阴性-与四条规避规范` | code | 1.0 |
| 249 | PHPWind 7.x 入口 × action 全量枚举方法论（251 个 action + 10 类危险标记）+ 写点判死结论集（pw_var_export 转义 / Pcv 封穿越 / uploadicon 数学不可能 / mutiupload 三重锁） | path-traversal | `phpwind-7.x-入口-action-全量枚举方法论-251-个-action-10-类危险标记-写点判死结论集-pw_var_export-转义-pcv` | code | 1.0 |
| 250 | PHPWind 7.x 后台风格 CSS 编辑 → data/style/*_css.htm → updatecache_sy eval 代码注入 RCE | rce | `phpwind-7.x-后台风格-css-编辑-data-style-_css.htm-updatecache_sy-eval-代码注入-rce` | code | 1.0 |
| 251 | PHPWind 写 PHP 文件 sink 的三面叠加扫描法（数据来源分类 + 变量路径补扫） | rce | `phpwind-写-php-文件-sink-的三面叠加扫描法-数据来源分类-变量路径补扫` | code | 1.0 |
| 252 | PHPWind 写 .php 缓存文件判死边界速查表（含 bank_sort 例外与字段长度表） | app-audit | `phpwind-写-.php-缓存文件判死边界速查表-含-bank_sort-例外与字段长度表` | code | 1.0 |
| 253 | PHPWind <target> admin/bakup.php 任意 SQL 执行（bakindata 多语句 + bakuptable 表名注入） | sqli | `phpwind-7.3.2-admin-bakup.php-任意-sql-执行-bakindata-多语句-bakuptable-表名注入` | fuzzy | 0.972 |
| 254 | PHPWind <target> 系统性加固判死清单（bank_sort RCE 链 + pw_hack 写入面 + 权限模型全量审计） | rce | `phpwind-7.3.2-系统性加固判死清单-bank_sort-rce-链-pw_hack-写入面-权限模型全量审计` | fuzzy | 0.968 |
| 255 | GBK 站点后台表单「全量回放」必须走 latin-1 字节级往返（PHPWind <target> 实测） | app-audit | `gbk-站点后台表单-全量回放-必须走-latin-1-字节级往返-phpwind-7.3.2-实测` | fuzzy | 0.96 |
| 256 | PHPWind <target> 后台 Getshell 优先级手册（权限自检 148xx 判据 + P0 SQL OUTFILE → P5 附件扩展名 + 源码级判死清单） | rce | `phpwind-7.3.2-后台-getshell-优先级手册-权限自检-148xx-判据-p0-sql-outfile-p5-附件扩展名-源码级判死清单` | fuzzy | 0.975 |
| 257 | PHPWind <target> 后台附件白名单放开（GBK 全量表单回放）与上传失败分层判据 | info-leak | `phpwind-7.3.2-后台附件白名单放开-gbk-全量表单回放-与上传失败分层判据` | fuzzy | 0.958 |
| 258 | 出口 IP 被宝塔 WAF 拉黑后的应急出口：免费 HTTP 代理池（抓取/验活/轮换 + schannel 坑） | tunnel | `出口-ip-被宝塔-waf-拉黑后的应急出口-免费-http-代理池-抓取-验活-轮换-schannel-坑` | code | 1.0 |
| 259 | 宝塔网站防火墙 <target> 拦截边界与绕过（IP 层换出口 / URL 关键词混合大小写 / multipart 双拦 PHP 特征） | waf | `宝塔网站防火墙-9.8.9-拦截边界与绕过-ip-层换出口-url-关键词混合大小写-multipart-双拦-php-特征` | fuzzy | 0.969 |
| 260 | 出口 IP 与 Host 头是两个独立变量：环境代理静默改道导致宝塔 WAF 封禁被误读（curl 必须 --noproxy "*"） | waf | `出口-ip-与-host-头是两个独立变量-环境代理静默改道导致宝塔-waf-封禁被误读-curl-必须-noproxy` | code | 1.0 |
| 261 | PHPWind 7.x 前台/后台 POST 一律被 Showmsg('undefined_action') 拒绝 —— 根因是 PostCheck() 要求 Referer 的 host 严格等于 HTTP_HOST（用 IP 做 Referer 是致命陷阱） | app-audit | `phpwind-7.x-前台-后台-post-一律被-showmsg-undefined_action-拒绝-根因是-postcheck-要求-referer` | code | 1.0 |
| 262 | 宝塔 WAF 来源封禁实测判据：multipart 上传后全出口 000/ProxyError，而同代理到 baidu 200（区分「代理失效」vs「目标拒我」） | waf | `宝塔-waf-来源封禁实测判据-multipart-上传后全出口-000-proxyerror-而同代理到-baidu-200-区分-代理失效-vs-目标拒我` | code | 1.0 |
| 263 | 宝塔网站防火墙 <target> 拦截页指纹（403/1200）＋ SQL 关键字载荷被整拦的实测判据＋代理出口两级验活法（--noproxy 自杀陷阱） | env | `宝塔网站防火墙-9.8.9-拦截页指纹-403-1200-sql-关键字载荷被整拦的实测判据-代理出口两级验活法-noproxy-自杀陷阱` | fuzzy | 0.973 |
| 265 | SPA+Spring Cloud 网关批量未授权接口探测（/xxa 前缀型） | unauthorized | `spa-spring-cloud-网关批量未授权接口探测-xxa-前缀型` | code | 1.0 |
| 266 | 点云客服 dotclouds 访客端未授权 multipart 任意文件上传（net.dotclouds / Spring Boot） | file-upload | `点云客服-dotclouds-访客端未授权-multipart-任意文件上传-net.dotclouds-spring-boot` | code | 1.0 |
| 267 | 点云客服 dotclouds 会话接口零鉴权读取（IDOR，queryChatHistory/queryConvMessage） | auth-bypass | `点云客服-dotclouds-会话接口零鉴权读取-idor-querychathistory-queryconvmessage` | code | 1.0 |
| 270 | 点云客服 dotclouds sendConvMessageFile：companyId 未 normalize 导致任意目录文件写入（路径穿越、可逃出 /data 到根文件系统） | file-upload | `点云客服-dotclouds-sendConvMessageFile-companyId-路径穿越-任意目录文件写入` | code | 1.0 |
| 271 | dysaas 代理后台：CDN 边缘 IP 轮换绕过源站封禁 + 请求格式/AES + 匿名代理注册 appSave | unauthorized | `dysaas-代理后台-cdn边缘ip轮换绕过源站封禁-请求格式与aes-匿名代理注册appsave` | code | 1.0 |
| 272 | 直播平台未授权注册绕过短信验证（/api/user/simple_regist 空 mobile_hash 直接签发 token） | auth-bypass | `直播平台未授权注册绕过短信验证-api-user-simple_regist-空-mobile_hash-直接签发-token` | code | 1.0 |
| 273 | 腾讯云 STS 临时凭证 + COS q-signature 手工签名（写入/边界判定）可用脚本 | unauthorized | `腾讯云-sts-临时凭证-cos-q-signature-手工签名-写入-边界判定-可用脚本` | code | 1.0 |
| 274 | 多租户直播平台前端响应加密还原（AES-256-CBC 硬编码 key/IV + 时间戳尾部拼接） | info-leak | `多租户直播平台前端响应加密还原-aes-256-cbc-硬编码-key-iv-时间戳尾部拼接` | code | 1.0 |
| 275 | 匿名注册/申请接口作账号存在性 oracle（用户名枚举通用手法 + dysaas 代理 SaaS 判定表） | info-leak | `匿名注册-申请接口作账号存在性-oracle-用户名枚举通用手法-dysaas-代理-saas-判定表` | code | 1.0 |
| 277 | 密码喷洒：批量生成的运营账号共用默认口令（先证明登录无需验证码 + 按前端算法构造口令密文） | weak-password | `密码喷洒-批量生成的运营账号共用默认口令-先证明登录无需验证码-按前端算法构造口令密文` | code | 1.0 |
| 278 | Swagger 驱动的水平越权批量探测（自己/他人 id 对照 + 未文档化参数注入 + 加密响应长度 oracle） | auth-bypass | `swagger-驱动的水平越权批量探测-自己-他人-id-对照-未文档化参数注入-加密响应长度-oracle` | code | 1.0 |
| 279 | 多租户 SaaS 后台「已证否方向」坑位卡（上传无解析 / detail 非 IDOR / Bearer 前缀与单会话顶号 / 验证码漏校验 / 错误码语义表） | methodology | `多租户-saas-后台-已证否方向-坑位卡-上传无解析-detail-非-idor-bearer-前缀与单会话顶号-验证码漏校验-错误码语义表` | code | 1.0 |
| 286 | 博彩站群 <target>:10666 Java 反序列化/JNDI/Log4j2 全否（含 OOB 有效性复核与 /ecb 网关格式） | deserialization | `cn7-10666-java-deser-jndi-log4j-all-negative` | code | 1.0 |
| 290 | Web 缓存投毒/缓存欺骗面一站式测绘脚本（含缓存键归一化与 unkeyed 头判定） | recon | `web-缓存投毒-缓存欺骗面一站式测绘脚本-含缓存键归一化与-unkeyed-头判定` | code | 1.0 |
| 291 | CDNGreat 缓存/防护层指纹识别与缓存范围、缓存键判定方法 | recon | `cdngreat-缓存-防护层指纹识别与缓存范围-缓存键判定方法` | code | 1.0 |
| 292 | fakerToken 服务无鉴权探测：/fbsport/v1/getToken 与 /ob/v2/getTryToken 泄漏后端域名与 token | info-leak | `fakertoken-服务无鉴权探测-fbsport-v1-gettoken-与-ob-v2-gettrytoken-泄漏后端域名与-token` | code | 1.0 |
| 293 | 对象存储未授权列举/写入判定（含 SPA 回落假阳性排除 + 写探针回滚） | unauthorized | `对象存储未授权列举-写入判定-含-spa-回落假阳性排除-写探针回滚` | code | 1.0 |
| 294 | Vite/Vue3 minify SPA 接口调用点与请求体还原器（三层别名映射） | frontend-reverse | `vite-vue3-minify-spa-接口调用点与请求体还原器-三层别名映射` | code | 1.0 |
| 295 | 用户可控字段 → 无鉴权公开接口渲染 → 面向全站访客的投毒链（四步验证 + 如实边界） | unauthorized | `用户可控字段-无鉴权公开接口渲染-面向全站访客的投毒链-四步验证-如实边界` | code | 1.0 |
| 296 | 前端硬编码 requestData 常量重放 → 未授权签发第三方平台 token 与线路域名路由表 | info-leak | `前端硬编码-requestdata-常量重放-未授权签发第三方平台-token-与线路域名路由表` | code | 1.0 |
| 297 | fakerToken 硬编码 requestData 常量重放：未认证签发 OB 平台登录 token 与线路后端域名表（RB88/OB 站群通用） | unauthorized | `fakertoken-硬编码-requestdata-常量重放-未认证签发-ob-平台登录-token-与线路后端域名表-rb88-ob-站群通用` | code | 1.0 |
| 298 | Vite SPA 资金接口静态还原 + 二维码渲染定性（含带会话的页面内 fetch 验证） | info-leak | `vite-spa-资金接口静态还原-二维码渲染定性-含带会话的页面内-fetch-验证` | code | 1.0 |
| 299 | 统一 20000 未登录拦截的绕过穷尽判定法（含 body×CT 归因控制实验，防 116 误报） | auth-bypass | `统一-20000-未登录拦截的绕过穷尽判定法-含-body-ct-归因控制实验-防-116-误报` | code | 1.0 |
| 300 | cmwall「过移动墙」服务管理面板 <target> 未授权读写（38888）+ 零风险写验证流程 | auth-bypass | `cmwall-过移动墙-服务管理面板-v1.35-未授权读写-38888-零风险写验证流程` | fuzzy | 0.944 |
| 301 | 前端资源投毒可行性判定法：客户端渲染资金字段 + 无 SRI/无 CSP/资源域写方法（含正确分级口径） | frontend-reverse | `前端资源投毒可行性判定法-客户端渲染资金字段-无-sri-无-csp-资源域写方法-含正确分级口径` | code | 1.0 |
| 302 | 支付回调/通知端点探测法：静态零引用判定 + 未知路由三态判别 + 按站点命名风格派生字典（含负结论口径） | auth-bypass | `支付回调-通知端点探测法-静态零引用判定-未知路由三态判别-按站点命名风格派生字典-含负结论口径` | code | 1.0 |
| 303 | fakerToken/OB 的 api= 参数解密（硬编码 AES-ECB 密钥 OBTY20220712OBTY）与 requestData 不可解结论 | info-leak | `fakertoken-ob-的-api-参数解密-硬编码-aes-ecb-密钥-obty20220712obty-与-requestdata-不可解结论` | code | 1.0 |
| 304 | cmwall「过移动墙」wizard 面板 静态分析速查 + 命令注入实现基础（含厂商公开包/部署路径/注入验证要点） | rce | `cmwall-过移动墙-wizard-面板-静态分析速查-命令注入实现基础-含厂商公开包-部署路径-注入验证要点` | code | 1.0 |
| 305 | WAF 绕过/能力修正：直连 HTTPS 清空 SNI（+任意 Host）→ 200；带真实 SNI（<target>）→ RST（按 SNI 丢包） | waf | `waf-绕过-能力修正-直连-https-清空-sni-任意-host-200-带真实-tgt-rst-按-sni-丢包` | fuzzy | 0.966 |
| 306 | 加固 APK（zxprotect 抽取壳）零依赖逆向 + config2.json 配置引导链提取 | info-leak | `加固-apk-zxprotect-抽取壳-零依赖逆向-config2.json-配置引导链提取` | code | 1.0 |
| 307 | LHZS/LHService v27 网关免签名面定位 + 两层签名校验拒绝码判据 | unauthorized | `lhzs-lhservice-v27-网关免签名面定位-两层签名校验拒绝码判据` | code | 1.0 |
| 308 | 原生 socket 端口普查器（绕代理假象）tcpscan.py | recon | `原生-socket-端口普查器-绕代理假象-tcpscan.py` | code | 1.0 |
| 309 | LHZS/靶标节点签名探测器（严格判据）lhprobe.py | recon | `lhzs-靶标节点签名探测器-严格判据-lhprobe.py` | code | 1.0 |
| 310 | TLS 证书抓取 + 从 DER 提 CN/SAN（证书归属判定）certgrab.py | info-leak | `tls-证书抓取-从-der-提-cn-san-证书归属判定-certgrab.py` | code | 1.0 |
| 311 | 多目标业务路径 + Host(vhost) 覆盖探测器 sweep.py | unauthorized | `多目标业务路径-host-vhost-覆盖探测器-sweep.py` | code | 1.0 |
| 312 | 非标准服务四态识别器 rawprobe.py | recon | `非标准服务四态识别器-rawprobe.py` | code | 1.0 |
| 313 | 端口 banner 抓取器（自动判 SSH/MySQL/HTTP）banner.py | recon | `端口-banner-抓取器-自动判-ssh-mysql-http-banner.py` | code | 1.0 |
| 314 | Web 指纹探测（标题 + LHZS 特征签名）httpfp.py | recon | `web-指纹探测-标题-lhzs-特征签名-httpfp.py` | code | 1.0 |
| 315 | 单目标业务路径 + vhost 探测（先看正文再判归属）bizprobe.py | unauthorized | `单目标业务路径-vhost-探测-先看正文再判归属-bizprobe.py` | code | 1.0 |
| 316 | DNS 解析 + 指定 SNI/Host 的 HTTPS 定向探测 sni72.py | recon | `dns-解析-指定-sni-host-的-https-定向探测-sni72.py` | code | 1.0 |
| 317 | nginx→Tomcat `..;/` 归一化差异：穿 location 白名单 + 跨 context 横向（含绕开 nginx 死上游 504）；Manager IP 阀 33 变体实测无效 | path-traversal | `nginx-tomcat-..-归一化差异-穿-location-白名单-跨-context-横向-含绕开-nginx-死上游-504-manager-ip-阀` | code | 1.0 |
| 318 | Tomcat 拿 shell 路线的快速判死：PUT 405 + OPTIONS 无 PUT ⇒ CVE-2017-12615/12617/2025-24813 全灭；WEB-INF 25 种写法全 404；Tribes RCE 与 Ghostcat 版本不匹配 | rce | `tomcat-拿-shell-路线的快速判死-put-405-options-无-put-cve-2017-12615-12617-2025-24813-全灭` | code | 1.0 |
| 319 | 全量参数钥匙：批量候选参数名越过老 servlet 的 Integer.parseInt 门槛（分组探测 + 二分收敛，通用可复用） | unauthorized | `全量参数钥匙-批量候选参数名越过老-servlet-的-integer.parseint-门槛-分组探测-二分收敛-通用可复用` | code | 1.0 |
| 320 | Tomcat Manager 快速判死：403 页字节级指纹（同版本 stock 页跨机器 sha1 完全一致）+ 端口面普查（只暴露 nginx ⇒ 头伪造必然无效，Tribes/8009 一并判死） | recon | `tomcat-manager-快速判死-403-页字节级指纹-同版本-stock-页跨机器-sha1-完全一致-端口面普查-只暴露-nginx-头伪造必然无效` | code | 1.0 |
| 321 | 盲 SSRF 判定法：差分三值（死地址基线/内网/外带）+ A/B/A 交替复测识破抖动 —— 附"零差异=无 SSRF"实测样例与真 SSRF 对照特征 | ssrf | `盲-ssrf-判定法-差分三值-死地址基线-内网-外带-a-b-a-交替复测识破抖动-附-零差异-无-ssrf-实测样例与真-ssrf-对照特征` | code | 1.0 |
| 322 | 六助(LHZS) App <target> 客户端接口字段名与请求体结构全表 | protocol | `六助-lhzs-app-v4.0.8-客户端接口字段名与请求体结构全表` | fuzzy | 0.926 |
| 323 | LHZS v27 管理台：认证链精确字段 + admin/query 信封 + release.publish SSRF + storage.configure 落盘可控 + action 全表 | auth-bypass | `lhzs-v27-管理台-认证链精确字段-admin-query-信封-release.publish-ssrf-storage.configure-落盘可控` | code | 1.0 |
| 324 | 【技巧】Java 异常行号 oracle：用 Tomcat 错误页的 Class.java:NN 行号推进，零猜测定位 servlet 参数与鉴权门槛 | methodology | `技巧-java-异常行号-oracle-用-tomcat-错误页的-class.java-nn-行号推进-零猜测定位-servlet-参数与鉴权门槛` | code | 1.0 |
| 325 | 【技巧】「参数钥匙串」批量命中老 Java servlet 数据分支：390 个候选参数各 =1，一次请求绕过猜参数名（单请求拿到 143KB/408 条用户记录） | methodology | `技巧-参数钥匙串-批量命中老-java-servlet-数据分支-390-个候选参数各-1-一次请求绕过猜参数名-单请求拿到-143kb-408-条用户记录` | code | 1.0 |
| 327 | 原生 socket / 手工 TLS 探测骨架（绕本机 HTTP 代理污染 + binary_form 证书取证 + 定制404/Spring404/406 三态判存） | env | `原生-socket-手工-tls-探测骨架-绕本机-http-代理污染-binary_form-证书取证-定制404-spring404-406-三态判存` | code | 1.0 |
| 328 | 老式 Java servlet「分块灌参 + 栈行号二分收敛」参数发现法（Tomcat 报错栈当 oracle） | info-leak | `老式-java-servlet-分块灌参-栈行号二分收敛-参数发现法-tomcat-报错栈当-oracle` | code | 1.0 |
| 329 | UpLoadPic 三个入口对照表（upload2=唯一真写点/写 OSS；无扩展名分支=只回 URL 模板不落盘；upload 的 head 四字段=真字段但空响应、token 模式=无条件假成功） | file-upload | `uploadpic-三个入口对照表-upload2-唯一真写点-写-oss-无扩展名分支-只回-url-模板不落盘-upload-的-head-四字段-真字段但` | code | 1.0 |
| 331 | 行号 oracle 参数隔离法（附 PutTp 判死实测）：用 Tomcat 错误页的 `类名.java:行号` 逐级隔离参数、二分求最小集、判 POST 是否独立分支 | methodology | `行号-oracle-参数隔离法-附-puttp-判死实测-用-tomcat-错误页的-类名.java-行号-逐级隔离参数-二分求最小集-判-post-是否独立分` | code | 1.0 |
| 335 | 复现冲突四步法：纯读稳定性 → 裸/戴钥匙对照 → 逐个剔除定最小集 → 跨时间窗读回（写侧+读侧都做对才不误判） | methodology | `复现冲突四步法-纯读稳定性-裸-戴钥匙对照-逐个剔除定最小集-跨时间窗读回-写侧-读侧都做对才不误判` | code | 1.0 |
| 336 | jadx 卡死时的替代：最小 DEX 解析器 + 字符串表邻接/簇分析（无需反编译） | frontend-reverse | `jadx-卡死时的替代-最小-dex-解析器-字符串表邻接-簇分析-无需反编译` | code | 1.0 |
| 338 | Go 自研代理服务（管理口 + 代理口）免凭据探测与绕过穷尽核查脚本 | unauthorized | `go-自研代理服务-管理口-代理口-免凭据探测与绕过穷尽核查脚本` | code | 1.0 |
| 339 | FastAdmin CVE-2024-7928 任意文件读取（/ajax/lang lang 参数） | path-traversal | `fastadmin-cve-2024-7928-任意文件读取-ajax-lang-lang-参数` | code | 1.0 |
| 340 | PHP .user.ini 两用技巧：追加 open_basedir=/ 破除目录限制 + auto_prepend_file 绕过 .php 写入黑名单落马（冰蝎 v4 双通道） | rce | `php-.user.ini-两用技巧-追加-open_basedir-破除目录限制-auto_prepend_file-绕过-.php-写入黑名单落马-冰蝎-v` | code | 1.0 |
| 343 | boyue(webman) 博彩站群会员端 —— upload-qrcode 任意扩展名写入 + 同源存储型 XSS（WithdrawController::uploadQrCode） | file-upload | `boyue-webman-博彩站群会员端-upload-qrcode-任意扩展名写入-同源存储型-xss-withdrawcontroller-uploadqr` | code | 1.0 |
| 344 | boyue(webman+Pear Admin) 博彩站群后台指纹 / 路由表 / 密钥 / 免登录路由 速查手册 | info-leak | `boyue-webman-pear-admin-博彩站群后台指纹-路由表-密钥-免登录路由-速查手册` | code | 1.0 |
| 345 | 本机 fake-ip DNS + 规则型本地代理下，用 CONNECT+自定义 SNI 直达 nginx vhost（绕过假性「目标封 IP」） | tunnel | `本机-fake-ip-dns-规则型本地代理下-用-connect-自定义-sni-直达-nginx-vhost-绕过假性-目标封-ip` | code | 1.0 |
| 346 | webman 插件路由指向非插件命名空间控制器 → AccessControl 中间件失效导致后台接口未授权 | unauthorized | `webman-插件路由指向非插件命名空间控制器-accesscontrol-中间件失效导致后台接口未授权` | code | 1.0 |
| 347 | webman 路由大小写不敏感 → 绕过按路径控制器段精确匹配的「功能关闭」开关（可救回代码生成器写文件 → Monitor 热重载 RCE） | auth-bypass | `webman-路由大小写不敏感-绕过按路径精确匹配的-功能关闭-开关-可救回代码生成器写文件-monitor-热重载-rce` | code | 1.0 |
| 348 | webman 默认(自动)路由绕过路由级 AuthMiddleware —— 未授权直达控制器方法（通用） | auth-bypass | `webman-默认-自动-路由绕过路由级-authmiddleware-未授权直达控制器方法-通用` | code | 1.0 |
| 349 | webman-admin TableController::crud 代码生成器 docblock 注入 → 写 app/controller/*.php → Monitor 热重载 RCE | rce | `webman-admin-tablecontroller-crud-代码生成器-docblock-注入-写-app-controller-.php-monito` | code | 1.0 |
| 350 | webman/Pear Admin 目标「上传→执行」全判定手册（include/LFI/phar/模板 四路排除 + 扩展名名单测绘） | path-traversal | `webman-pear-admin-目标-上传-执行-全判定手册-include-lfi-phar-模板-四路排除-扩展名名单测绘` | code | 1.0 |
| 351 | webman 响应加密（X-Encrypted: 1）硬编码 AES 密钥解密 / 加密请求体构造 | frontend-reverse | `webman-响应加密-x-encrypted-1-硬编码-aes-密钥解密-加密请求体构造` | code | 1.0 |
| 352 | webman 自动路由绕过鉴权中间件（/{controller}/{action} 零凭证直连） | auth-bypass | `webman-自动路由绕过鉴权中间件-controller-action-零凭证直连` | code | 1.0 |
| 353 | webman 自动路由作为「方法存在性 oracle」：零凭证枚举部署版控制器方法（区分零鉴权直出 / 需登录 / 只判角色） | unauthorized | `webman-自动路由作为-方法存在性-oracle-零凭证枚举部署版控制器方法-区分零鉴权直出-需登录-只判角色` | code | 1.0 |
| 354 | webman 自动路由未授权接口普查 + Crypto 硬编码密钥解密响应 | info-leak | `webman-自动路由未授权接口普查-crypto-硬编码密钥解密响应` | code | 1.0 |
| 355 | nginx vhost oracle 枚举（Host 头区分真实虚拟主机，识别同名后端多品牌） | recon | `nginx-vhost-oracle-枚举-host-头区分真实虚拟主机-识别同名后端多品牌` | code | 1.0 |
| 356 | 跨 vhost 同源 XSS 可执行性判定矩阵（多域名站群：/uploads 可达性 × /app/admin 交叉对照） | unauthorized | `跨-vhost-同源-xss-可执行性判定矩阵-多域名站群-uploads-可达性-app-admin-交叉对照` | code | 1.0 |
| 357 | webman-admin(Vue Art Design Pro) 后台会话与请求加密规范 + 403 开关大小写绕过 + crud docblock 注入开火参数 | auth-bypass | `webman-admin-vue-art-design-pro-后台会话与请求加密规范-403-开关大小写绕过-crud-docblock-注入开火参数` | code | 1.0 |
| 358 | boyue/PearAdmin 充值渠道链上收款地址未授权替换 POC（payset-update / configs_data[trc20]） | unauthorized | `boyue-pearadmin-充值渠道链上收款地址未授权替换-poc-payset-update-configs_data-trc20` | code | 1.0 |
| 359 | webman「显式路由 vs 自动路由」双姿势鉴权/闸门 oracle（可复用：一次判定整族方法被哪一层拦） | auth-bypass | `webman-显式路由-vs-自动路由-双姿势鉴权-闸门-oracle-可复用-一次判定整族方法被哪一层拦` | code | 1.0 |
| 360 | webman 自动路由「形参绑定」未授权身份/ID 注入（$request->all() → 控制器形参） | auth-bypass | `webman-自动路由-形参绑定-未授权身份-id-注入-request-all-控制器形参` | code | 1.0 |
| 361 | webman「形参名 oracle」：用 MissingInput/InputType 错误页在线反推控制器方法签名 | info-leak | `webman-形参名-oracle-用-missinginput-inputtype-错误页在线反推控制器方法签名` | code | 1.0 |
| 362 | SSRF 原语筛查法：从形参到外发 URL 的四步数据流溯源（含二阶配置写入者追踪） | ssrf | `ssrf-原语筛查法-从形参到外发-url-的四步数据流溯源-含二阶配置写入者追踪` | code | 1.0 |
| 363 | 零副作用验证判据套件：身份注入 / 原密码强校验 / 验证码校验（用不存在 uid 与现值回写） | auth-bypass | `零副作用验证判据套件-身份注入-原密码强校验-验证码校验-用不存在-uid-与现值回写` | code | 1.0 |
| 364 | 路径穿越有效性判定：归一化等价对照法（nginx alias 前缀映射场景） | path-traversal | `路径穿越有效性判定-归一化等价对照法-nginx-alias-前缀映射场景` | code | 1.0 |
| 365 | 任意后缀上传 → 同源存储型 XSS → localStorage JWT 会话窃取 → 无限速资金密码 oracle（账号接管/盗资金组合链） | file-upload | `任意后缀上传-同源存储型-xss-localstorage-jwt-会话窃取-无限速资金密码-oracle-账号接管-盗资金组合链` | code | 1.0 |
| 366 | GoFly 客服系统鉴权模型与越权判定（/2/ 族 401vs403 判据 + ws_kefu/ws_visitor 子协议票据） | unauthorized | `gofly-客服系统鉴权模型与越权判定-2-族-401vs403-判据-ws_kefu-ws_visitor-子协议票据` | code | 1.0 |
| 367 | 上游/供应链域名「可控性」判定流水线（whois+RDAP+DNS+HTTP+CT+TLS+家族挖掘） | recon | `上游-供应链域名-可控性-判定流水线-whois-rdap-dns-http-ct-tls-家族挖掘` | code | 1.0 |
| 368 | suo5 classic 隧道断链自愈 + PHP<target> 冰蝎马 openssl_decrypt("AES128") 返回 false 导致 200 空响应 | tunnel | `suo5-classic-隧道断链自愈-php5.4-冰蝎马-openssl_decrypt-aes128-返回-false-导致-200-空响应` | fuzzy | 0.983 |
| 369 | ThinkPHP3 后台 Cookie 会话伪造（admin_user=username + admin_id=password哈希） | auth-bypass | `thinkphp3-后台-cookie-会话伪造-admin_user-username-admin_id-password哈希` | code | 1.0 |
| 370 | ThinkPHP3 后台 Api::article_save 任意文件写入（data:// 伪协议）→ WebShell | file-upload | `thinkphp3-后台-api-article_save-任意文件写入-data-伪协议-webshell` | code | 1.0 |
| 371 | 冰蝎马 PHP 5.6 适配：openssl_decrypt("AES128") 失效须改 "aes-128-ecb" | rce | `冰蝎马-php-5.6-适配-openssl_decrypt-aes128-失效须改-aes-128-ecb` | code | 1.0 |
| 372 | ThinkPHP <target> 站群 CMS：extends Action 的 Admin 控制器绕鉴权 + 微信签名可自算 → 未授权文件写入（Wechat::index_1） | auth-bypass | `thinkphp-3.1.3-站群-cms-extends-action-的-admin-控制器绕鉴权-微信签名可自算-未授权文件写入-wechat-index` | fuzzy | 0.969 |
| 373 | PHP-FPM FastCGI 未授权访问 → PHP_ADMIN_VALUE 解除 open_basedir / 清空 disable_functions（纯 PHP 客户端，无需 shell） | info-leak | `php-fpm-fastcgi-未授权访问-php_admin_value-解除-open_basedir-清空-disable_functions-纯-php` | code | 1.0 |
| 374 | ThinkPHP <target> 站群 CMS：Home/Order::get_html1 零凭证覆盖模板文件（App/Tpl/Home/Index_index.html）→ 模板注入前置原语 | unauthorized | `thinkphp-3.1.3-站群-cms-home-order-get_html1-零凭证覆盖模板文件-app-tpl-home-index_index.ht` | fuzzy | 0.904 |
| 375 | ThinkPHP <target> 「未授权模板写入 → 编译执行」RCE 链（Home/Order::get_html1 → App/Tpl → App/Runtime/Cache） | rce | `thinkphp-3.1.3-未授权模板写入-编译执行-rce-链-home-order-get_html1-app-tpl-app-runtime-cache` | fuzzy | 0.976 |
| 376 | Admin/Wechat 免鉴权签名绕过（token 空串伪造 sha1 签名） | auth-bypass | `admin-wechat-免鉴权签名绕过-token-空串伪造-sha1-签名` | code | 1.0 |
| 378 | 冰蝎马 Python 自动化客户端（PHP <target> 需把 cipher 改成 aes-128-ecb） | tunnel | `冰蝎马-python-自动化客户端-php-5.6.40-需把-cipher-改成-aes-128-ecb` | fuzzy | 0.949 |
| 379 | ThinkPHP 3：VAR_FILTERS=filter_exp 让 `exp` 数组注入全线失效（静默丢字段、无痕）——判定向量 + 无效绕过清单 | auth-bypass | `thinkphp-3-var_filters-filter_exp-让-exp-数组注入全线失效-静默丢字段-无痕-判定向量-无效绕过清单` | code | 1.0 |
| 381 | 站群通用：article.ename = 创建者 admin 用户名 + Home/Index::bbs?post=&lt;id&gt; 零凭证渲染点（枚举后台账号名 / 内容可控渲染） | info-leak | `站群通用-article.ename-创建者-admin-用户名-home-index-bbs-post-lt-id-gt-零凭证渲染点-枚举后台账号名-内容可` | code | 1.0 |
| 382 | UEditor PHP 零鉴权面盲探（config 泄露 / listimage 路径枚举 / upload 目录存在性 / catchimage SSRF） | unauthorized | `ueditor-php-零鉴权面盲探-config-泄露-listimage-路径枚举-upload-目录存在性-catchimage-ssrf` | code | 1.0 |
| 384 | ThinkPHP<target> 站群通用链：后台 mass assignment 写 article.title → 零凭证 get_html1 模板烘焙 → 首页 RCE | rce | `thinkphp3.1.3-站群通用链-后台-mass-assignment-写-article.title-零凭证-get_html1-模板烘焙-首页-rce` | fuzzy | 0.976 |
| 385 | PHP-FPM unix socket 未授权 + PHP_ADMIN_VALUE 清空 open_basedir/disable_functions → 全盘任意文件读写（宝塔环境通用） | privesc | `php-fpm-unix-socket-未授权-php_admin_value-清空-open_basedir-disable_functions-全盘任意文件` | code | 1.0 |
| 386 | 对手检测能力画像：站群运营方 ThinkPHP3 后门自查脚本 jnbgg.sh（= 我们的规避指南） | recon | `对手检测能力画像-站群运营方-thinkphp3-后门自查脚本-jnbgg.sh-我们的规避指南` | code | 1.0 |
| 387 | php-fpm FastCGI 沙箱逃逸：PHP_ADMIN_VALUE 清零 open_basedir 与 disable_functions（宝塔环境通用） | privesc | `php-fpm-fastcgi-沙箱逃逸-php_admin_value-清零-open_basedir-与-disable_functions-宝塔环境通用` | code | 1.0 |
| 388 | 站群兄弟节点发现：从目标自己的 vhost 访问日志 Referer 反查（日志分块读 + 一条请求判定同站群） | recon | `站群兄弟节点发现-从目标自己的-vhost-访问日志-referer-反查-日志分块读-一条请求判定同站群` | code | 1.0 |
| 389 | php-fpm FastCGI sendmail_path 劫持 → 绕过 disable_functions 拿命令执行（mail() 的 C 层 popen） | privesc | `php-fpm-fastcgi-sendmail_path-劫持-绕过-disable_functions-拿命令执行-mail-的-c-层-popen` | code | 1.0 |
| 390 | 熊猫体育 TY SDK /yewu40 签名通道无防重放 + checkId 不校验（试玩 token 实测） | auth-bypass | `ty-sdk-yewu40-signed-channel-no-replay-protection` | code | 1.0 |
| 392 | PHPWind 7.x/8.x 版本指纹识别（robots.txt + meta generator + 路径特征） | info-leak | `phpwind-7.x-8.x-版本指纹识别-robots.txt-meta-generator-路径特征` | code | 1.0 |
| 393 | FOFA 官方接口可用性切换（<target> 失效 → <target>） | recon | `fofa-官方接口可用性切换-tgt-失效-tgt` | exact | 1.0 |
| 395 | PHPWind 站点精确版本与数据库表前缀的无认证识别 | info-leak | `phpwind-站点精确版本与数据库表前缀的无认证识别` | code | 1.0 |
| 399 | WordPress/WooCommerce 匿名攻击面清单（版本泄露+REST未授权+Store API非授权写入+XML-RPC爆破） | info-leak | `wordpress-woocommerce-匿名攻击面清单-版本泄露-rest未授权-store-api非授权写入-xml-rpc爆破` | code | 1.0 |
| 402 | FastAPI/OpenAPI 未授权测绘 —— openapi.json 泄露全端点 + 按 security 标记自动分类实测 | info-leak | `fastapi-openapi-未授权测绘-openapi.json-泄露全端点-按-security-标记自动分类实测` | code | 1.0 |
| 403 | ThinkPHP5 /kj 开奖数据接口未授权读取 | unauthorized | `thinkphp5-kj-开奖数据接口未授权读取` | code | 1.0 |
| 404 | 自研BBS 站群 JSON API 硬编码密钥 + 接口协议逆向 | unauthorized | `自研bbs-站群-json-api-硬编码密钥-接口协议逆向` | code | 1.0 |
| 405 | PHPWind 后台登录判据纠正 —— 失败态=14540B+页内计数文案；失败也下发 cookie；必须带 cookie 复访验证 | auth-bypass | `phpwind-后台登录判据纠正-失败态-14540b-页内计数文案-失败也下发-cookie-必须带-cookie-复访验证` | code | 1.0 |
| 406 | 页面源码挖自动生成 vhost —— 从 HTML/JS 提取域名并批量解析（补字典爆破盲区） | info-leak | `页面源码挖自动生成-vhost-从-html-js-提取域名并批量解析-补字典爆破盲区` | code | 1.0 |
| 407 | MinIO Console 两段式登录逆向 + Access Key 用户名枚举 | info-leak | `minio-console-two-stage-login-enum` | code | 1.0 |
| 408 | TLS 证书 SAN 反查域名（纯 IP 段挖真实业务域名 + vhost 路由） | info-leak | `tls-cert-san-domain-harvest` | code | 1.0 |
| 409 | AJ-Captcha 验证码协议复现 + secretKey 明文下发利用（芋道/RuoYi 系列） | auth-bypass | `aj-captcha-protocol-and-secretkey` | code | 1.0 |
| 410 | TLS 证书 SAN 反查域名 + 同 IP vhost 路由发现（纯 IP 段挖真实域名，避免错过 target） | info-leak | `tls-证书-san-反查域名-同-ip-vhost-路由发现-纯-ip-段挖真实域名-避免错过-target` | code | 1.0 |
| 411 | 母舰域名编号枚举 —— 一次画出攻击者全部节点（56 节点/10 网段实测） | info-leak | `母舰域名编号枚举-一次画出攻击者全部节点-56-节点-10-网段实测` | code | 1.0 |
| 412 | SPA 后台突破 —— 从前端 api.js 提取全接口清单 + 定向弱口令（实测 admin/admin123 拿到超管） | weak-password | `spa-后台突破-从前端-api.js-提取全接口清单-定向弱口令-实测-admin-admin123-拿到超管` | code | 1.0 |
| 413 | Synology DSM 未授权侦察 —— 两端点泄露 API 清单 / 主机名 / 公网状态 / 全部已装套件 | info-leak | `synology-dsm-unauth-initdata-recon` | code | 1.0 |
| 414 | 方法论纠错：高速 SYN 扫描在 SYN-flood 防护设备上产生海量假开放端口（必须低速复核） | recon | `nmap-highrate-false-positive-on-appliance` | code | 1.0 |
| 415 | 群晖 DSM 套件前端 bundle 未鉴权扒取 + API/协议逆向（含容器终端 WebSocket 协议） | info-leak | `synology-dsm-static-bundle-unauth-harvest` | code | 1.0 |
| 416 | 群晖 DSM code 120 判读纠错 —— 参数校验先于权限，不可误判为未授权 | methodology | `synology-error-code-120-not-unauth` | code | 1.0 |
| 422 | GoFly 客服附件卡片 safeChatAssetPath(%2e%2e) 绕过 → 同源存储型 XSS 主动投递链 | file-upload | `gofly-attach-card-safechatassetpath-bypass-xss-delivery` | code | 1.0 |
| 424 | PHPDisk <target> adminlogin 账号枚举 + bind_user 无频控口令校验双通道 | auth-bypass | `phpdisk-6.8-adminlogin-账号枚举-bind_user-无频控口令校验双通道` | fuzzy | 0.975 |
| 425 | PHPDisk 6.x adminlogin 注入面三点连测 + 后台 gid=1 校验 | auth-bypass | `phpdisk-6.x-adminlogin-注入面三点连测-后台-gid-1-校验` | code | 1.0 |
| 426 | New API 定制构建「全关闭」实例：注册关/setup已初始化/SMTP缺省下的负面情报与速判清单 | info-leak | `new-api-定制构建-全关闭-实例-注册关-setup已初始化-smtp缺省下的负面情报与速判清单` | code | 1.0 |
| 427 | go-captcha 拼图行为验证码破解：Sobel 梯度模板匹配定位 + /api/captcha/verify 直接换 token（4/4 实测通过） | auth-bypass | `go-captcha-drag-sobel-edge-template-crack` | code | 1.0 |
| 428 | New API Token Quota 面板 /token-quota.config.json 配置泄露 + 令牌状态判别法（unavailable=存在但禁用） | info-leak | `token-quota-panel-config-leak-and-token-status-judgement` | code | 1.0 |
| 429 | <target> 注册墙全链路程序化突破：mail.tm 临时邮箱 + GET /api/verification 真实收码 + go-captcha(register) 建号成功（新用户初始 quota=0） | auth-bypass | `new-api-mailtm-tempmail-register-flow-api-modelsell-cc` | code | 1.0 |
| 432 | New API 邮箱验证注册绕过（一次性邮箱 + 6 位十六进制验证码 + 无图形验证码） | auth-bypass | `无-cve-设计缺陷-new-api-邮箱验证注册绕过-一次性邮箱-6-位十六进制验证码-无图形验证码` | code | 1.0 |
| 433 | PHPWind <target> 全量源码获取：<target> 官方下载站快照直取法（免 GitHub） | info-leak | `phpwind-732-source-obtain-wayback` | code | 1.0 |
| 438 | 七牛 Kodo 泄露 AK/SK：全量枚举、精确计量与公开读桶匿名外带 | info-leak | `七牛-kodo-泄露-ak-sk-全量枚举-精确计量与公开读桶匿名外带` | code | 1.0 |
| 440 | melody-admin(renren-fast 派生) 调度器 mode 语义与 MySQL UDF 面判否 + UDF 正控方法 | rce | `artx-melody-schedule-mode-udf-neg` | code | 1.0 |
| 441 | openresty 1.29 + Spring Boot：请求走私(CL.TE/TE.CL/TE.TE)与 h2 降级全族阴性（勿再重打） | auth-bypass | `openresty-1.29-spring-boot-请求走私-cl.te-te.cl-te.te-与-h2-降级全` | code | 1.0 |
| 442 | melody-admin 彩种模板同步 SSRF：回显只认 data 数组信封；60s 时延不能当内网 DROP 判据 | ssrf | `melody-admin-彩种模板同步-ssrf-回显只认-data-数组信封-60s-时延不能当内网-drop-判据` | code | 1.0 |
| 444 | X-Forwarded-For 等代理头伪造客户端 IP | auth-bypass | `ip-spoof-x-forwarded-for` | code | 1.0 |
| 445 | MySQL root 空口令远程直连检测（只读取证脚本） | weak-password | `mysql-root-空口令远程直连检测-只读取证脚本` | code | 1.0 |
| 446 | Laravel 图片上传（客户端后缀白名单+雪花ID重命名+Flysystem path 参数）无法转 RCE，残留 SVG 同源 XSS 与 4MiB FPM DoS | file-upload | `laravel-图片上传-客户端后缀白名单-雪花id重命名-flysystem-path-参数-无法转-rce-残留-s` | code | 1.0 |
| 448 | 短信找回密码链：sendSmsCode 响应明文回显 OTP + modifyPasswordBySms 未认证改密 → 任意会员账号接管 | auth-bypass | `短信找回密码链-sendsmscode-响应明文回显-otp-modifypasswordbysms-未认证改密` | code | 1.0 |
| 449 | PHP session.upload_progress 会话文件注入 + LFI → 预认证 RCE（Chain A） | rce | `php-session.upload_progress-会话文件注入-lfi-预认证-rce-chain-a` | code | 1.0 |
| 450 | PHP var_destroy 二次析构 UAF → `/proc/self/mem` 覆写 internal function handler → RCE（Chain B） | rce | `php-var_destroy-二次析构-uaf-proc-self-mem-覆写-internal-function-handler-rce-chain-b` | code | 1.0 |
| 451 | PHP session serialize_handler 不一致 → 解码失败仍触发 __destruct → POP gadget RCE（Chain C） | rce | `php-session-serialize_handler-不一致-解码失败仍触发-__destruct-pop-gadget-rce-chain-c` | code | 1.0 |
| 452 | CVE-2026-109618 Apache mod_php：multipart 解析期 OOM → longjmp 跳过 init_executor → shutdown 读脏 EG() → RCE（Chain D，零依赖） | rce | `cve-2026-109618-apache-mod_php-multipart-解析期-oom-longjmp-跳过-init_executor-shutdo` | code | 1.0 |
| 453 | CVE-2026-109619 nginx + php-fpm：脏 BG(unserialize) 状态被下一请求消费 → __destruct → RCE（Chain E） | rce | `cve-2026-109619-nginx-php-fpm-脏-bg-unserialize-状态被下一请求消费-__destruct-rce-chain-e` | code | 1.0 |
| 454 | PHP 环境存在性判别器（跑 PHP2Shell 五链前的必备预检） | app-audit | `php-环境存在性判别器-跑-php2shell-五链前的必备预检` | code | 1.0 |
| 457 | PHP2Shell 通用利用框架（CVE-2026-109618/619）：session.upload_progress UAF → Heap OOB 泄露绕 ASLR → pDestructor=zif_system 预认证 RCE | rce | `php2shell-exploit-framework-cve-2026-109618-619-heap-oob-asl` | code | 1.0 |
| 458 | 3x-ui-x-ui-assets-unauth-directory-listing-and-dual-404-base | absorbed | `3x-ui-x-ui-assets-unauth-directory-listing-and-dual-404-base` | code | 1.0 |
| 459 | cloudfront-behavior-prefix-vs-origin-normalization-bypass-co | absorbed | `cloudfront-behavior-prefix-vs-origin-normalization-bypass-co` | code | 1.0 |
| 460 | edge-host-127.0.0.1-baota-default-site-nginx-stub-status-una | absorbed | `edge-host-127.0.0.1-baota-default-site-nginx-stub-status-una` | code | 1.0 |
| 461 | edge-host-whitelist-third-party-vhost-isolation-config-and-p | absorbed | `edge-host-whitelist-third-party-vhost-isolation-config-and-p` | code | 1.0 |
| 462 | fofa-官方接口可用性切换-target417.example.com-失效-fofa.info | absorbed | `fofa-官方接口可用性切换-target417.example.com-失效-fofa.info` | code | 1.0 |
| 463 | h2c-grpc-日志摄取接口-pb.logingestservice-ingestlogevents-未授权调用-无 | absorbed | `h2c-grpc-日志摄取接口-pb.logingestservice-ingestlogevents-未授权调用-无` | code | 1.0 |
| 464 | landing-page-localstorage-query-channelcode-iframe-src-path | absorbed | `landing-page-localstorage-query-channelcode-iframe-src-path` | code | 1.0 |
| 465 | mail-auth-missing-spf-dkim-dmarc-email-spoofing | absorbed | `mail-auth-missing-spf-dkim-dmarc-email-spoofing` | code | 1.0 |
| 466 | minio-admin-v3-streaming-unsigned-payload-trailer-bypass-cve | absorbed | `minio-admin-v3-streaming-unsigned-payload-trailer-bypass-cve` | code | 1.0 |
| 467 | minio-streaming-unsigned-payload-trailer-sigv4-bypass-cve-20 | absorbed | `minio-streaming-unsigned-payload-trailer-sigv4-bypass-cve-20` | code | 1.0 |
| 468 | mysql-3306-exposed-host-not-allowed-ip-whitelist-three-egres | absorbed | `mysql-3306-exposed-host-not-allowed-ip-whitelist-three-egres` | code | 1.0 |
| 469 | omnihelm-console-grpc-bff-host-oracle | absorbed | `omnihelm-console-grpc-bff-host-oracle` | code | 1.0 |
| 470 | oss-listbucket-denied-getobject-allowed-anon-read-matrix | absorbed | `oss-listbucket-denied-getobject-allowed-anon-read-matrix` | code | 1.0 |
| 471 | photon-edge-clte-desync-smuggle-0chunk-tail-forwarded | absorbed | `photon-edge-clte-desync-smuggle-0chunk-tail-forwarded` | code | 1.0 |
| 472 | postmessage-missing-origin-config-injection-js-sink | absorbed | `postmessage-missing-origin-config-injection-js-sink` | code | 1.0 |
| 473 | postmessage-source-identity-bypass-iframe-navigation-new-fun | absorbed | `postmessage-source-identity-bypass-iframe-navigation-new-fun` | code | 1.0 |
| 474 | postmessage-仅校验-source-的-executefunction-绕过-跨源导航被校验子帧-new | absorbed | `postmessage-仅校验-source-的-executefunction-绕过-跨源导航被校验子帧-new` | code | 1.0 |
| 475 | sclub-168six-up_photo-a_photo-头像与直播台图标上传-.php-文件名被丢弃-gd-重编码裁 | absorbed | `sclub-168six-up_photo-a_photo-头像与直播台图标上传-.php-文件名被丢弃-gd-重编码裁` | code | 1.0 |
| 476 | somcns-thinkphp3.1.3-站群-166-站-admin-autoart-控制器脱队-extends-ac | absorbed | `somcns-thinkphp3.1.3-站群-166-站-admin-autoart-控制器脱队-extends-ac` | code | 1.0 |
| 477 | somcns-thinkphp3.1.3-站群-精简版-实测负面清单-控制器文件存在性预言机-wechat-auto | absorbed | `somcns-thinkphp3.1.3-站群-精简版-实测负面清单-控制器文件存在性预言机-wechat-auto` | code | 1.0 |
| 478 | somcns-站群前台会员登录-home-index-login-无验证码无频控-test-test-直登-后台登 | absorbed | `somcns-站群前台会员登录-home-index-login-无验证码无频控-test-test-直登-后台登` | code | 1.0 |
| 479 | sub2api-0.2.19-activities-fork-reset-feature-gate-dead | absorbed | `sub2api-0.2.19-activities-fork-reset-feature-gate-dead` | code | 1.0 |
| 480 | sub2api-0.2.19-密码重置链审查-feature-gate-双开关-token-设计-reset-不清-to-cve-2026-27812 | absorbed | `sub2api-0.2.19-密码重置链审查-feature-gate-双开关-token-设计-reset-不清-to-cve-2026-27812` | code | 1.0 |
| 481 | sub2api-admin-login-weak-password-domain-derived-jwt-takeove | absorbed | `sub2api-admin-login-weak-password-domain-derived-jwt-takeove` | code | 1.0 |
| 482 | sub2api-invalid-auth-rate-limit-429-shared-bucket-dos | absorbed | `sub2api-invalid-auth-rate-limit-429-shared-bucket-dos` | code | 1.0 |
| 483 | sub2api-photonedge-nginx-path-normalization-no-auth-bypass | absorbed | `sub2api-photonedge-nginx-path-normalization-no-auth-bypass` | code | 1.0 |
| 484 | sub2api-plain-http-no-redirect-missing-hsts-sslstrip | absorbed | `sub2api-plain-http-no-redirect-missing-hsts-sslstrip` | code | 1.0 |
| 485 | sub2api-validate-promo-code-unauth-enum-oracle | absorbed | `sub2api-validate-promo-code-unauth-enum-oracle` | code | 1.0 |
| 486 | sub2api-网关-api-key-未随凭据轮换失效-dump-泄露的-sk-key-仍可调用-v1-付费接口 | absorbed | `sub2api-网关-api-key-未随凭据轮换失效-dump-泄露的-sk-key-仍可调用-v1-付费接口` | code | 1.0 |
| 487 | sub2api-邮箱验证码接口账号存在性枚举-oracle-公开-api-v1-auth-send-verify-co | absorbed | `sub2api-邮箱验证码接口账号存在性枚举-oracle-公开-api-v1-auth-send-verify-co` | code | 1.0 |
| 488 | target160.example.com-1852-补丁面黑盒实测矩阵-13-类已封堵-5-个门槛不可越-write_config-落点存在但入口不可达 | absorbed | `target160.example.com-1852-补丁面黑盒实测矩阵-13-类已封堵-5-个门槛不可越-write_config-落点存在但入口不可达` | code | 1.0 |
| 489 | thinkphp-3.1.3-站群-cms-admin-控制器-403-假分叉-判别-实为-nginx-层危险函数 | absorbed | `thinkphp-3.1.3-站群-cms-admin-控制器-403-假分叉-判别-实为-nginx-层危险函数` | code | 1.0 |
| 490 | thinkphp-3.1.3-站群-home-index-bbs-的-post-参数零鉴权渲染文章正文-游客态可批量拖 | absorbed | `thinkphp-3.1.3-站群-home-index-bbs-的-post-参数零鉴权渲染文章正文-游客态可批量拖` | code | 1.0 |
| 491 | thinkphp-3.1.3-站群两个零凭证原生语-think_page_trace-泄全部内部路径-app | absorbed | `thinkphp-3.1.3-站群两个零凭证原生语-think_page_trace-泄全部内部路径-app` | code | 1.0 |
| 492 | thinkphp-站群-kindeditor-php-版本-upload_json.php-file_manager | absorbed | `thinkphp-站群-kindeditor-php-版本-upload_json.php-file_manager` | code | 1.0 |
| 493 | 博彩站群-203.0.113.1-10666-java-反序列化-jndi-log4j-rce-全否-严格 | absorbed | `博彩站群-203.0.113.1-10666-java-反序列化-jndi-log4j-rce-全否-严格` | code | 1.0 |
| 494 | 博彩站群-203.0.113.28-10666-java-反序列化-jndi-log4j-rce-全否-严格 | absorbed | `博彩站群-203.0.113.28-10666-java-反序列化-jndi-log4j-rce-全否-严格` | code | 1.0 |
| 495 | 邮件伪造类-finding-缺失-spf-dkim-dmarc-的核实方法与-不可利用-边界 | absorbed | `邮件伪造类-finding-缺失-spf-dkim-dmarc-的核实方法与-不可利用-边界` | code | 1.0 |
| 496 | 203.0.113.1-.23-先验排除表-两条恒定假阳性-code-验证码端点-爆破前置-五处已穷尽面 | absorbed | `203.0.113.1-.23-先验排除表-两条恒定假阳性-code-验证码端点-爆破前置-五处已穷尽面` | code | 1.0 |
| 497 | 203.0.113.1-站群跳转机全机地图与负面情报-唯一动态程序-phpwind-7.3.2-全静态页-三套异机真入口 | absorbed | `203.0.113.1-站群跳转机全机地图与负面情报-唯一动态程序-phpwind-7.3.2-全静态页-三套异机真入口` | code | 1.0 |
| 498 | 203.0.113.10-.23-先验排除表-两条恒定假阳性-code-验证码端点-爆破前置-五处已穷尽面 | absorbed | `203.0.113.10-.23-先验排除表-两条恒定假阳性-code-验证码端点-爆破前置-五处已穷尽面` | code | 1.0 |
| 499 | 203.0.113.13-站群跳转机全机地图与负面情报-唯一动态程序-phpwind-7.3.2-全静态页-三套异机真入口 | absorbed | `203.0.113.13-站群跳转机全机地图与负面情报-唯一动态程序-phpwind-7.3.2-全静态页-三套异机真入口` | code | 1.0 |
| 500 | discuz-7.2-后台-runquery-sql-导入对非创始人-adminid-1-全线-403-isfounde | absorbed | `discuz-7.2-后台-runquery-sql-导入对非创始人-adminid-1-全线-403-isfounde` | code | 1.0 |
| 501 | phpwind-7.3.2-7.5-后台写文件-getshell-三处写入端判死-pcv-封穿越-pw_var_export-转义数组键-setcss-非写入端 | absorbed | `phpwind-7.3.2-7.5-后台写文件-getshell-三处写入端判死-pcv-封穿越-pw_var_export-转义数组键-setcss-非写入端` | code | 1.0 |
| 502 | phpwind-7.3.2-负面情报清单-register_globals-已移除-mawhole-eval-证否-setstyles-转义-hack.php | absorbed | `phpwind-7.3.2-负面情报清单-register_globals-已移除-mawhole-eval-证否-setstyles-转义-hack.php` | code | 1.0 |
| 503 | phpwind-8.7-decode-错位链-白盒审计结论-无需账号的裸引号-convert-e-路径不存在-cms-投稿链裸引号入库-展示-rce-代码级成立 | absorbed | `phpwind-8.7-decode-错位链-白盒审计结论-无需账号的裸引号-convert-e-路径不存在-cms-投稿链裸引号入库-展示-rce-代码级成立` | code | 1.0 |
| 504 | phpwind-8.7-补丁面速查表-target124.example.com-白盒产出-已修复-9-类-未验证缺陷-6-条-未穷尽面-4-个-全附源码行号 | absorbed | `phpwind-8.7-补丁面速查表-target124.example.com-白盒产出-已修复-9-类-未验证缺陷-6-条-未穷尽面-4-个-全附源码行号` | code | 1.0 |
| 505 | sclub-平台-index.phtml-注入回显客户端-含三个易踩根因与语句约束表 | absorbed | `sclub-平台-index.phtml-注入回显客户端-含三个易踩根因与语句约束表` | code | 1.0 |
| 506 | sclub-平台-index.phtml-注入客户端修正-繁体歡迎-hex左截断-from包层-限流特征 | absorbed | `sclub-平台-index.phtml-注入客户端修正-繁体歡迎-hex左截断-from包层-限流特征` | code | 1.0 |
| 507 | tgt-1852-补丁面黑盒实测矩阵-13-类已封堵-5-个门槛不可越-write_config-落点存在但入口不可达 | absorbed | `tgt-1852-补丁面黑盒实测矩阵-13-类已封堵-5-个门槛不可越-write_config-落点存在但入口不可达` | code | 1.0 |
| 508 | tgt2403-admin-后台指纹与-api-v1-admin-管理-api-全清单提取-含未授权定性探测 | absorbed | `tgt2403-admin-后台指纹与-api-v1-admin-管理-api-全清单提取-含未授权定性探测` | code | 1.0 |
| 509 | tgt2403-发卡系统-订单-游客订单鉴权模型-api-契约-订单号格式与-idor-判定速查 | absorbed | `tgt2403-发卡系统-订单-游客订单鉴权模型-api-契约-订单号格式与-idor-判定速查` | code | 1.0 |
| 510 | tgt2403-版本指纹-上游补丁差分法-app_version-免认证泄露-缺失安全修复清单-逐条可达性实测 | absorbed | `tgt2403-版本指纹-上游补丁差分法-app_version-免认证泄露-缺失安全修复清单-逐条可达性实测` | code | 1.0 |
| 511 | tgt2403-用户中心与认证链路速查-me-wallet-affiliate-2fa-鉴权事实-验证码节流与限流缺口 | absorbed | `tgt2403-用户中心与认证链路速查-me-wallet-affiliate-2fa-鉴权事实-验证码节流与限流缺口` | code | 1.0 |
| 512 | thinkphp-3.1.3-站群-cms-模板烘焙-rce-链完整打法-两个否证-oracle-get_html1-零凭证覆盖模板-abrdns-抓取源陷阱 | absorbed | `thinkphp-3.1.3-站群-cms-模板烘焙-rce-链完整打法-两个否证-oracle-get_html1-零凭证覆盖模板-abrdns-抓取源陷阱` | code | 1.0 |
| 513 | thinkphp3.1.3-站群-cms-通杀-rce-admin-表哈希即凭据-home-article-edit_cont-写-article.title | absorbed | `thinkphp3.1.3-站群-cms-通杀-rce-admin-表哈希即凭据-home-article-edit_cont-写-article.title` | code | 1.0 |
| 514 | waf-绕过-能力修正-直连-https-清空-sni-任意-host-200-带真实-target320.example.com-rst-按-sni-丢包 | absorbed | `waf-绕过-能力修正-直连-https-清空-sni-任意-host-200-带真实-target320.example.com-rst-按-sni-丢包` | code | 1.0 |
| 515 | 代理分流造假样本-target30.example.com-经-7892-返回伪造-502-404-hello | absorbed | `代理分流造假样本-target30.example.com-经-7892-返回伪造-502-404-hello` | code | 1.0 |
| 516 | 博彩-phpwind-站群实战结论集-203.0.113.1-24-判死路线-8-条可复用判据 | absorbed | `博彩-phpwind-站群实战结论集-203.0.113.1-24-判死路线-8-条可复用判据` | code | 1.0 |
| 517 | 博彩-phpwind-站群实战结论集-203.0.113.27-24-判死路线-8-条可复用判据 | absorbed | `博彩-phpwind-站群实战结论集-203.0.113.27-24-判死路线-8-条可复用判据` | code | 1.0 |
| 518 | Nexus Repository 未授权 RCE（CVE-2019-7238） | cloud-container | `fujiang-cve-2019-7238` | code | 1.0 |
| 519 | PostgreSQL COPY PROGRAM 权限边界（CVE-2019-9193，厂商认定非漏洞） | cloud-container | `fujiang-cve-2019-9193` | code | 1.0 |
| 520 | vCenter Server 未授权上传 → RCE（vROPS 插件） | cloud-container | `fujiang-cve-2021-21972` | code | 1.0 |
| 521 | Nacos 认证绕过（User-Agent 白名单） | cloud-container | `fujiang-cve-2021-29441` | code | 1.0 |
| 522 | Grafana 插件路径穿越任意文件读（CVE-2021-43798） | cloud-container | `fujiang-cve-2021-43798` | code | 1.0 |
| 523 | cgroups release_agent 容器逃逸 | cloud-container | `fujiang-cve-2022-0492` | code | 1.0 |
| 524 | VMware Workspace ONE Access 模板注入 RCE | cloud-container | `fujiang-cve-2022-22954` | code | 1.0 |
| 525 | Veeam Backup 凭据泄露（CVE-2023-27532） | cloud-container | `fujiang-cve-2023-27532` | code | 1.0 |
| 526 | runc 文件描述符泄露逃逸（Leaky Vessels） | cloud-container | `fujiang-cve-2024-21626` | code | 1.0 |
| 527 | VMware ESXi 认证绕过 → hypervisor 失守（CVE-2024-37085） | cloud-container | `fujiang-cve-2024-37085` | code | 1.0 |
| 528 | Veeam Backup 未认证 RCE（CVE-2024-40711） | cloud-container | `fujiang-cve-2024-40711` | code | 1.0 |
| 529 | Docker Remote API 未授权（直接拿宿主） | cloud-container | `fujiang-cloud-container-615246bd` | code | 1.0 |
| 530 | MongoDB 未授权访问（直接拖库） | cloud-container | `fujiang-cloud-container-0b0a3d30` | code | 1.0 |
| 531 | Redis 未授权（写文件 / 模块加载 / 主从复制） | cloud-container | `fujiang-cve-2022-0543` | code | 1.0 |
| 532 | FortiOS SSL VPN 路径穿越（凭据读取） | edge-device | `fujiang-cve-2018-13379` | code | 1.0 |
| 533 | Pulse Secure 任意文件读取（凭据与会话） | edge-device | `fujiang-cve-2019-11510` | code | 1.0 |
| 534 | Citrix ADC 路径穿越 → 模板注入 RCE | edge-device | `fujiang-cve-2019-19781` | code | 1.0 |
| 535 | FortiOS SSL VPN 双因素绕过（用户名大小写） | edge-device | `fujiang-cve-2020-12812` | code | 1.0 |
| 536 | F5 BIG-IP TMUI 未授权 RCE（..;/ 路径穿越） | edge-device | `fujiang-cve-2020-5902` | code | 1.0 |
| 537 | F5 iControl REST 未授权 RCE（2021 变体） | edge-device | `fujiang-cve-2021-22986` | code | 1.0 |
| 538 | F5 BIG-IP iControl REST 认证绕过 → 命令执行 | edge-device | `fujiang-cve-2022-1388` | code | 1.0 |
| 539 | FortiOS / FortiProxy 管理面认证绕过（假 Forwarded 头） | edge-device | `fujiang-cve-2022-40684` | code | 1.0 |
| 540 | Ivanti EPMM 认证绕过（CVE-2023-35078） | edge-device | `fujiang-cve-2023-35078` | code | 1.0 |
| 541 | Ivanti Connect Secure 认证绕过 + 命令注入（串联） | edge-device | `fujiang-cve-2023-46805-plus2` | code | 1.0 |
| 542 | FortiClient EMS SQL 注入 → 未认证 RCE（CVE-2023-48788） | edge-device | `fujiang-cve-2023-48788` | code | 1.0 |
| 543 | PAN-OS 管理面认证绕过 → root RCE（CVE-2024-0012 + CVE-2024-9474） | edge-device | `fujiang-cve-2024-0012-plus2` | code | 1.0 |
| 544 | GoAnywhere MFT 认证绕过 → 新建管理员（CVE-2024-0204） | edge-device | `fujiang-cve-2024-0204` | code | 1.0 |
| 545 | Check Point Quantum 安全网关任意文件读（CVE-2024-24919） | edge-device | `fujiang-cve-2024-24919` | code | 1.0 |
| 546 | PAN-OS GlobalProtect 命令注入 | edge-device | `fujiang-cve-2024-3400` | code | 1.0 |
| 547 | FortiManager 未认证 RCE（FortiJump，CVE-2024-47575） | edge-device | `fujiang-cve-2024-47575` | code | 1.0 |
| 548 | AD CS 证书滥用（ESC1 / ESC8 等） | exchange-ad | `fujiang-exchange-ad-e6e03808` | code | 1.0 |
| 549 | EternalBlue（MS17-010 SMBv1 远程溢出） | exchange-ad | `fujiang-cve-2017-0144` | code | 1.0 |
| 550 | Exchange ViewState 反序列化（静态 machineKey） | exchange-ad | `fujiang-cve-2020-0688` | code | 1.0 |
| 551 | Zerologon（Netlogon 密钥置零 → 域控接管） | exchange-ad | `fujiang-cve-2020-1472` | code | 1.0 |
| 552 | Exchange ProxyLogon（SSRF + 任意文件写入 → RCE） | exchange-ad | `fujiang-cve-2021-26855` | code | 1.0 |
| 553 | Exchange ProxyShell（三个洞串联 → RCE） | exchange-ad | `fujiang-cve-2021-34473` | code | 1.0 |
| 554 | PrintNightmare（打印后台服务任意 DLL 加载） | exchange-ad | `fujiang-cve-2021-34527` | code | 1.0 |
| 555 | noPac（sAMAccountName 伪装 + 无 PAC 票据 → 域管） | exchange-ad | `fujiang-cve-2021-42287-plus2` | code | 1.0 |
| 556 | Exchange 认证后反序列化 RCE（CVE-2021-42321） | exchange-ad | `fujiang-cve-2021-42321` | code | 1.0 |
| 557 | MSDT 协议远程代码执行（Follina） | exchange-ad | `fujiang-cve-2022-30190` | code | 1.0 |
| 558 | Exchange ProxyNotShell（SSRF + PowerShell 远程） | exchange-ad | `fujiang-cve-2022-41040-plus2` | code | 1.0 |
| 559 | Exchange EWS NTLM 中继提权（CVE-2024-21410） | exchange-ad | `fujiang-cve-2024-21410` | code | 1.0 |
| 560 | Outlook MonikerLink 绕过与 NTLM 泄露（CVE-2024-21413） | exchange-ad | `fujiang-cve-2024-21413` | code | 1.0 |
| 561 | Zimbra postjournal 未授权命令执行（CVE-2024-45519） | exchange-ad | `fujiang-cve-2024-45519` | code | 1.0 |
| 562 | WebLogic T3 反序列化（CommonsCollections） | java-deser | `fujiang-cve-2015-4852` | code | 1.0 |
| 563 | Shiro rememberMe 默认密钥反序列化 RCE（CVE-2016-4437） | java-deser | `fujiang-cve-2016-4437` | code | 1.0 |
| 564 | WebLogic wls-wsat XMLDecoder 反序列化 RCE | java-deser | `fujiang-cve-2017-10271` | code | 1.0 |
| 565 | Struts2 S2-045 Content-Type OGNL 注入 | java-deser | `fujiang-cve-2017-5638` | code | 1.0 |
| 566 | Apache Log4j2 JNDI 注入（Log4Shell） | java-deser | `fujiang-cve-2021-44228` | code | 1.0 |
| 567 | SnakeYAML 反序列化 RCE（CVE-2022-1471） | java-deser | `fujiang-cve-2022-1471` | code | 1.0 |
| 568 | Java ECDSA 空签名绕过（Psychic Signatures） | java-deser | `fujiang-cve-2022-21449` | code | 1.0 |
| 569 | Spring Framework 数据绑定 RCE（Spring4Shell） | java-deser | `fujiang-cve-2022-22965` | code | 1.0 |
| 570 | Fastjson autoType 反序列化 RCE | java-deser | `fujiang-cve-2022-25845` | code | 1.0 |
| 571 | Tomcat PUT 任意文件写入（JSP 落地） | java-web | `fujiang-cve-2017-12615` | code | 1.0 |
| 572 | Apache Solr Velocity 模板注入 RCE | java-web | `fujiang-cve-2019-17558` | code | 1.0 |
| 573 | WebLogic wls9_async_response 未授权 RCE（2019 在野） | java-web | `fujiang-cve-2019-2725` | code | 1.0 |
| 574 | WebLogic 控制台未授权 RCE（console.portal 路径穿越） | java-web | `fujiang-cve-2020-14882` | code | 1.0 |
| 575 | Tomcat AJP 文件读取 / 包含（Ghostcat） | java-web | `fujiang-cve-2020-1938` | code | 1.0 |
| 576 | Apache Druid 未授权 RCE（CVE-2021-25646） | java-web | `fujiang-cve-2021-25646` | code | 1.0 |
| 577 | Confluence 未授权 OGNL 注入（pre-auth，ENCODED_TOKEN） | java-web | `fujiang-cve-2021-26084` | code | 1.0 |
| 578 | Spring Cloud Gateway Actuator SpEL 注入 RCE（CVE-2022-22947） | java-web | `fujiang-cve-2022-22947` | code | 1.0 |
| 579 | Spring Cloud Function SpEL 注入 RCE（CVE-2022-22963） | java-web | `fujiang-cve-2022-22963` | code | 1.0 |
| 580 | Confluence 未授权 OGNL 注入（Confluence OGNL RCE） | java-web | `fujiang-cve-2022-26134` | code | 1.0 |
| 581 | Confluence 未授权管理员接管（setup 流程未关） | java-web | `fujiang-cve-2023-22515` | code | 1.0 |
| 582 | Confluence 未授权恢复导致管理员接管 | java-web | `fujiang-cve-2023-22518` | code | 1.0 |
| 583 | Confluence 模板注入 RCE（text-inline.vm） | java-web | `fujiang-cve-2023-22527` | code | 1.0 |
| 584 | PaperCut NG/MF 认证绕过 → RCE（CVE-2023-27350） | java-web | `fujiang-cve-2023-27350` | code | 1.0 |
| 585 | TeamCity 认证绕过 → 管理员令牌 → RCE | java-web | `fujiang-cve-2023-42793` | code | 1.0 |
| 586 | Apache ActiveMQ OpenWire 反序列化 RCE | java-web | `fujiang-cve-2023-46604` | code | 1.0 |
| 587 | Struts2 文件上传路径穿越（S2-066） | java-web | `fujiang-cve-2023-50164` | code | 1.0 |
| 588 | GitLab 任意账号接管（CVE-2023-7028） | java-web | `fujiang-cve-2023-7028` | code | 1.0 |
| 589 | Jenkins CLI 任意文件读取 | java-web | `fujiang-cve-2024-23897` | code | 1.0 |
| 590 | TeamCity 鉴权绕过（;.jsp 路径参数） | java-web | `fujiang-cve-2024-27198` | code | 1.0 |
| 591 | Apache OFBiz 未授权 RCE（CVE-2024-38856） | java-web | `fujiang-cve-2024-38856` | code | 1.0 |
| 592 | Struts2 文件上传路径穿越 → RCE（S2-067，CVE-2024-53677） | java-web | `fujiang-cve-2024-53677` | code | 1.0 |
| 593 | Baron Samedit（sudo 堆溢出提权） | linux-priv | `fujiang-cve-2021-3156` | code | 1.0 |
| 594 | Ubuntu overlayfs 本地提权（2021） | linux-priv | `fujiang-cve-2021-3493` | code | 1.0 |
| 595 | PwnKit（polkit pkexec 本地提权） | linux-priv | `fujiang-cve-2021-4034` | code | 1.0 |
| 596 | fs_context 堆溢出提权（legacy_parse_param） | linux-priv | `fujiang-cve-2022-0185` | code | 1.0 |
| 597 | Dirty Pipe（页缓存写覆盖提权） | linux-priv | `fujiang-cve-2022-0847` | code | 1.0 |
| 598 | overlayfs copy-up 权限绕过提权 | linux-priv | `fujiang-cve-2023-0386` | code | 1.0 |
| 599 | sudoedit 环境变量注入提权（EDITOR 带 --） | linux-priv | `fujiang-cve-2023-22809` | code | 1.0 |
| 600 | Ubuntu overlayfs 权限提升（GameOver(lay)） | linux-priv | `fujiang-cve-2023-2640-plus2` | code | 1.0 |
| 601 | nf_tables 释放后重用提权（2023） | linux-priv | `fujiang-cve-2023-32233` | code | 1.0 |
| 602 | Linux io_uring 本地提权（CVE-2024-0582） | linux-priv | `fujiang-cve-2024-0582` | code | 1.0 |
| 603 | nf_tables 释放后重用提权 | linux-priv | `fujiang-cve-2024-1086` | code | 1.0 |
| 604 | 命令注入（分号 / 反引号 / $() 与盲注外带） | web-cms | `fujiang-web-cms-fc7b661a` | code | 1.0 |
| 605 | PHPUnit eval-stdin 远程代码执行 | web-cms | `fujiang-cve-2017-9841` | code | 1.0 |
| 606 | Drupal 表单 API 远程代码执行（Drupalgeddon2） | web-cms | `fujiang-cve-2018-7600` | code | 1.0 |
| 607 | PHP-FPM 路径信息下溢 RCE | web-cms | `fujiang-cve-2019-11043` | code | 1.0 |
| 608 | vBulletin widget 配置代码执行 | web-cms | `fujiang-cve-2019-16759` | code | 1.0 |
| 609 | Laravel Ignition 调试模式任意文件读 → RCE（CVE-2021-3129） | web-cms | `fujiang-cve-2021-3129` | code | 1.0 |
| 610 | Apache HTTP Server 路径穿越 / CGI RCE | web-cms | `fujiang-cve-2021-41773-plus2` | code | 1.0 |
| 611 | Joomla 未授权配置泄露（API） | web-cms | `fujiang-cve-2023-23752` | code | 1.0 |
| 612 | WordPress Backup Migration 未授权 RCE（CVE-2023-6553） | web-cms | `fujiang-cve-2023-6553` | code | 1.0 |
| 613 | WordPress Bricks 未授权 RCE（CVE-2024-25600） | web-cms | `fujiang-cve-2024-25600` | code | 1.0 |
| 614 | WordPress WP-Automatic 未授权 SQL 注入（CVE-2024-27956） | web-cms | `fujiang-cve-2024-27956` | code | 1.0 |
| 615 | PHP-CGI Windows 参数注入（best-fit 编码） | web-cms | `fujiang-cve-2024-4577` | code | 1.0 |
| 616 | 本地文件包含 → RCE（日志包含 / pearcmd / PHP filter chain） | web-cms | `fujiang-web-cms-e2998ea3` | code | 1.0 |
| 617 | Nginx alias 路径穿越（读任意文件/源码泄露） | web-cms | `fujiang-web-cms-8b06abba` | code | 1.0 |
| 618 | PHP 反序列化（POP 链 → RCE / phar 触发） | web-cms | `fujiang-web-cms-3986d73b` | code | 1.0 |
| 619 | SQL 注入（报错 / UNION / 堆叠）通用打法 | web-cms | `fujiang-web-cms-0280c2e0` | code | 1.0 |
| 620 | SSTI 模板注入 → RCE（Jinja2 / Twig / Freemarker） | web-cms | `fujiang-web-cms-2439a0f4` | code | 1.0 |
| 621 | ThinkPHP 5.x 方法注入 / 路由 RCE | web-cms | `fujiang-cve-2018-20062-plus2` | code | 1.0 |
| 622 | XXE 外部实体注入（读文件 / SSRF / OOB 外带） | web-cms | `fujiang-web-cms-1b19f77e` | code | 1.0 |
| 623 | 后台管理 API 未授权（只在前端隐藏菜单） | gray-market | `fujiang-gray-market-c28f79db` | code | 1.0 |
| 624 | 后台 IP 白名单绕过（X-Forwarded-For 伪造内网来源） | gray-market | `fujiang-gray-market-2ef09faa` | code | 1.0 |
| 625 | 后台“重置任意用户密码”接口越权（直接接管会员/管理员） | gray-market | `fujiang-gray-market-cc7b3103` | code | 1.0 |
| 626 | 代理层级越权（读全站下级 / 改佣金比例） | gray-market | `fujiang-gray-market-300fe1f8` | code | 1.0 |
| 627 | 接口越权 / BOLA（读他人订单、流水、提现） | gray-market | `fujiang-gray-market-214935b4` | code | 1.0 |
| 628 | 接口无 nonce / 时间戳 → 请求重放（充值、领奖、下注） | gray-market | `fujiang-gray-market-5fd82b60` | code | 1.0 |
| 629 | 注单 / 流水接口越权批量拉取（全站会员数据） | gray-market | `fujiang-gray-market-d967e600` | code | 1.0 |
| 630 | 下注/结算赔率与倍数篡改（资金面） | gray-market | `fujiang-gray-market-269a4e2e` | code | 1.0 |
| 631 | 支付 / 回调地址可控 → 内网 SSRF（含云元数据） | gray-market | `fujiang-gray-market-a1c33d56` | code | 1.0 |
| 632 | 前端签名 / 加密可逆（伪造请求绕过风控） | gray-market | `fujiang-gray-market-a8e64394` | code | 1.0 |
| 633 | CORS 配置错误（带凭据读他人接口数据） | gray-market | `fujiang-gray-market-086e458d` | code | 1.0 |
| 634 | 券码 / 彩金码可枚举（批量领取与套现） | gray-market | `fujiang-gray-market-76e00ef0` | code | 1.0 |
| 635 | 客服聊天记录导出未鉴权（全站会话与身份材料） | gray-market | `fujiang-gray-market-49c0a3f1` | code | 1.0 |
| 636 | CSRF（无 token 的提现/改密/绑卡） | gray-market | `fujiang-gray-market-fa5aed6d` | code | 1.0 |
| 637 | 签到 / 任务奖励批量刷取（缺幂等与频率控制） | gray-market | `fujiang-gray-market-77e09490` | code | 1.0 |
| 638 | phpMyAdmin / Adminer 暴露（弱口令直进数据库） | gray-market | `fujiang-gray-market-387d2767` | code | 1.0 |
| 639 | Druid 监控未授权 → session 泄露/后台接管（若依系必备） | gray-market | `fujiang-gray-market-893b29b8` | code | 1.0 |
| 640 | 导出 / 下载接口任意文件读取 | gray-market | `fujiang-gray-market-508fca6b` | code | 1.0 |
| 641 | 财务导出接口越权（银行卡 + 身份证批量导出） | gray-market | `fujiang-gray-market-6416fd1c` | code | 1.0 |
| 642 | 三方游戏 / 视讯接口 key 泄露与回调伪造 | gray-market | `fujiang-gray-market-dbf31f8d` | code | 1.0 |
| 643 | Host 头投毒（密码重置链接被劫持 → 任意账号接管） | gray-market | `fujiang-gray-market-fade6caf` | code | 1.0 |
| 644 | 客服 IM / WebSocket 未鉴权（订阅他人会话） | gray-market | `fujiang-gray-market-4f6d3954` | code | 1.0 |
| 645 | 后台"切换用户/代登录"接口越权（任意账号接管） | gray-market | `fujiang-gray-market-4a9998a4` | code | 1.0 |
| 646 | 邀请返水自邀循环（小号矩阵刷返水） | gray-market | `fujiang-gray-market-dc82f024` | code | 1.0 |
| 647 | JWT 实现缺陷（alg:none / 弱密钥 / 可预测 / kid 注入） | gray-market | `fujiang-gray-market-39a60016` | code | 1.0 |
| 648 | 开奖结果可预测（弱 PRNG / 客户端算法 / 固定种子） | gray-market | `fujiang-gray-market-61cc5c65` | code | 1.0 |
| 649 | Mass Assignment（注册/资料接口注入 role / vip / balance / agentId） | gray-market | `fujiang-gray-market-d4dc739a` | code | 1.0 |
| 650 | 会员隐私批量查询 / 账号枚举（手机号、银行卡） | gray-market | `fujiang-gray-market-9811f67d` | code | 1.0 |
| 651 | 包网多租户越权（tenantId/siteId 遍历，横穿几十个站点） | gray-market | `fujiang-gray-market-e1ef8859` | code | 1.0 |
| 652 | 通知模板注入（标题/正文模板未转义 → SSTI / 存储型 XSS） | gray-market | `fujiang-gray-market-793b8a5c` | code | 1.0 |
| 653 | OAuth redirect_uri 劫持 / state 缺失（账号接管） | gray-market | `fujiang-gray-market-bfcbf2a8` | code | 1.0 |
| 654 | 对象存储与数据存储裸露（OSS / S3 / MinIO / MySQL / MongoDB） | gray-market | `fujiang-gray-market-e7b37cd6` | code | 1.0 |
| 655 | 运营后台弱口令 + 验证码复用（博彩站点最常见入口） | gray-market | `fujiang-gray-market-1217f1c1` | code | 1.0 |
| 656 | 订单状态机篡改（未支付改成已支付 / 已取消再入账） | gray-market | `fujiang-gray-market-679fc3c8` | code | 1.0 |
| 657 | 源站 IP 泄露 → 绕过 WAF 直打（博彩站必看） | gray-market | `fujiang-gray-market-d6cabfd5` | code | 1.0 |
| 658 | 宝塔 / 1Panel 面板弱口令与未授权（黑灰站最常见后台） | gray-market | `fujiang-gray-market-db6d779c` | code | 1.0 |
| 659 | 密码重置 token 可预测 / 未绑定用户（任意用户改密） | gray-market | `fujiang-gray-market-0e1a7416` | code | 1.0 |
| 660 | 支付通道 / 商户号篡改（充值进攻击者账户） | gray-market | `fujiang-gray-market-5a2545a6` | code | 1.0 |
| 661 | 支付回调伪造与金额篡改（资金面） | gray-market | `fujiang-gray-market-66d638f4` | code | 1.0 |
| 662 | 活动 / 彩金 / 返水无限领取（幂等缺失） | gray-market | `fujiang-gray-market-c68b2427` | code | 1.0 |
| 663 | 余额并发与负数（竞态 / 金额符号） | gray-market | `fujiang-gray-market-8e14e2e4` | code | 1.0 |
| 664 | Redis 未授权 → 偷后台 session / token 直接进后台 | gray-market | `fujiang-gray-market-e418b17b` | code | 1.0 |
| 665 | 源码 / 备份 / 配置泄露 | gray-market | `fujiang-gray-market-f63dd78c` | code | 1.0 |
| 666 | 报表/财务接口未授权（全站流水与代理佣金） | gray-market | `fujiang-gray-market-1a807271` | code | 1.0 |
| 667 | 若依（RuoYi）系后台默认口令 / Druid / 定时任务 RCE | gray-market | `fujiang-gray-market-e7ff5ae1` | code | 1.0 |
| 668 | 后台定时任务 / 计划任务注入（Quartz、XXL-JOB、若依） | gray-market | `fujiang-gray-market-91afaa10` | code | 1.0 |
| 669 | 短信轰炸与账号枚举（验证码面） | gray-market | `fujiang-gray-market-a7909e0a` | code | 1.0 |
| 670 | 短信验证码回显 / 复用 / 爆破（任意账号接管） | gray-market | `fujiang-gray-market-f9c739b1` | code | 1.0 |
| 671 | sourcemap / 前端 JS 泄露接口与密钥（SPA 包网站第一入口） | gray-market | `fujiang-gray-market-dd337e8b` | code | 1.0 |
| 672 | Spring Boot Actuator 未授权（heapdump / Jolokia / Gateway） | gray-market | `fujiang-cve-2022-22947-2` | code | 1.0 |
| 673 | 业务接口 SQL 注入（报表 / 订单 / 代理查询） | gray-market | `fujiang-gray-market-2a36eeee` | code | 1.0 |
| 674 | SSRF（图片抓取 / 回调校验 / 推广链接检测）打内网与云元数据 | gray-market | `fujiang-gray-market-4125a4fa` | code | 1.0 |
| 675 | 会员资料存储型 XSS 打后台审核（管理员会话劫持） | gray-market | `fujiang-gray-market-c0d9c765` | code | 1.0 |
| 676 | 客服工单附件越权下载（证件照 / KYC 材料） | gray-market | `fujiang-gray-market-30f052ae` | code | 1.0 |
| 677 | Swagger / Knife4j 接口文档暴露 → 未授权接口直打 | gray-market | `fujiang-gray-market-83cca5f2` | code | 1.0 |
| 678 | 上传点 getshell（头像 / KYC 证件 / 客服附件） | gray-market | `fujiang-gray-market-d33cbf51` | code | 1.0 |
| 679 | VIP/等级接口参数篡改（免费拿高等级权益） | gray-market | `fujiang-gray-market-f76beb82` | code | 1.0 |
| 680 | 提现审核接口越权（读他人银行卡 / 改审核状态） | gray-market | `fujiang-gray-market-0c77ad6f` | code | 1.0 |
| 681 | 提现银行卡越权替换（把别人提现收进自己卡） | gray-market | `fujiang-gray-market-25410c01` | code | 1.0 |
| 4682 | SAP NetWeaver Remote Code Execution Vulnerability | kev | `cve-2010-5326` | code | 1.0 |
| 4683 | HP OpenView Network Node Manager Remote Code Execution Vulnerability | kev | `cve-2005-2773` | code | 1.0 |
| 4684 | Alcatel OmniPCX Enterprise Remote Code Execution Vulnerability | kev | `cve-2007-3010` | code | 1.0 |
| 4685 | phpMyAdmin Remote Code Execution Vulnerability | kev | `cve-2009-1151` | code | 1.0 |
| 4686 | Oracle JRE Unspecified Vulnerability | kev | `cve-2010-0840` | code | 1.0 |
| 4687 | Adobe ColdFusion Directory Traversal Vulnerability | kev | `cve-2010-2861` | code | 1.0 |
| 4688 | Exim Heap-Based Buffer Overflow Vulnerability | kev | `cve-2010-4344` | code | 1.0 |
| 4689 | Microsoft Forefront TMG Remote Code Execution Vulnerability | kev | `cve-2011-1889` | code | 1.0 |
| 4690 | Adobe Reader and Acrobat Universal 3D Memory Corruption Vulnerability | kev | `cve-2011-2462` | code | 1.0 |
| 4691 | Oracle Java SE Runtime Environment (JRE) Arbitrary Code Execution Vulnerability | kev | `cve-2011-3544` | code | 1.0 |
| 4692 | Apache Struts 2 Improper Input Validation Vulnerability | kev | `cve-2012-0391` | code | 1.0 |
| 4693 | Oracle Fusion Middleware Unspecified Vulnerability | kev | `cve-2012-1710` | code | 1.0 |
| 4694 | Oracle Java SE Runtime Environment (JRE) Arbitrary Code Execution Vulnerability | kev | `cve-2012-4681` | code | 1.0 |
| 4695 | Adobe ColdFusion Authentication Bypass Vulnerability | kev | `cve-2013-0625` | code | 1.0 |
| 4696 | Adobe ColdFusion Authentication Bypass Vulnerability | kev | `cve-2013-0632` | code | 1.0 |
| 4697 | Apache Struts Improper Input Validation Vulnerability | kev | `cve-2013-2251` | code | 1.0 |
| 4698 | Adobe Reader and Acrobat Arbitrary Integer Overflow Vulnerability | kev | `cve-2013-2729` | code | 1.0 |
| 4699 | Adobe Reader and Acrobat Memory Corruption Vulnerability | kev | `cve-2013-3346` | code | 1.0 |
| 4700 | HP Multiple Products Remote Code Execution Vulnerability | kev | `cve-2013-4810` | code | 1.0 |
| 4701 | Adobe Flash Player Integer Underflow Vulnerablity | kev | `cve-2014-0497` | code | 1.0 |
| 4702 | Adobe Reader and Acrobat Sandbox Bypass Vulnerability | kev | `cve-2014-0546` | code | 1.0 |
| 4703 | InduSoft Web Studio NTWebServer Directory Traversal Vulnerability | kev | `cve-2014-0780` | code | 1.0 |
| 4704 | Multi-Router Looking Glass (MRLG) Buffer Overflow Vulnerability | kev | `cve-2014-3931` | code | 1.0 |
| 4705 | GNU Bourne-Again Shell (Bash) Arbitrary Code Execution Vulnerability | kev | `cve-2014-6271` | code | 1.0 |
| 4706 | Rejetto HTTP File Server (HFS) Remote Code Execution Vulnerability | kev | `cve-2014-6287` | code | 1.0 |
| 4707 | GNU Bourne-Again Shell (Bash) Arbitrary Code Execution Vulnerability | kev | `cve-2014-7169` | code | 1.0 |
| 4708 | Realtek SDK Improper Input Validation Vulnerability | kev | `cve-2014-8361` | code | 1.0 |
| 4709 | Adobe Flash Player Remote Code Execution Vulnerability | kev | `cve-2015-0311` | code | 1.0 |
| 4710 | Adobe Flash Player Use-After-Free Vulnerability | kev | `cve-2015-0313` | code | 1.0 |
| 4711 | D-Link and TRENDnet Multiple Devices Remote Code Execution Vulnerability | kev | `cve-2015-1187` | code | 1.0 |
| 4712 | Elasticsearch Groovy Scripting Engine Remote Code Execution Vulnerability | kev | `cve-2015-1427` | code | 1.0 |
| 4713 | Oracle Java SE and Java SE Embedded Remote Code Execution Vulnerability | kev | `cve-2015-2590` | code | 1.0 |
| 4714 | Adobe Flash Player Memory Corruption Vulnerability | kev | `cve-2015-3043` | code | 1.0 |
| 4715 | Adobe Flash Player Heap-Based Buffer Overflow Vulnerability | kev | `cve-2015-3113` | code | 1.0 |
| 4716 | Oracle WebLogic Server Deserialization of Untrusted Data Vulnerability | kev | `cve-2015-4852` | code | 1.0 |
| 4717 | Adobe Flash Player Use-After-Free Vulnerability | kev | `cve-2015-5119` | code | 1.0 |
| 4718 | Adobe Flash Player Use-After-Free Vulnerability | kev | `cve-2015-5122` | code | 1.0 |
| 4719 | Adobe Flash Player Use-After-Free Vulnerability | kev | `cve-2015-5123` | code | 1.0 |
| 4720 | IBM WebSphere Application Server and Server Hypervisor Edition Code Injection. | kev | `cve-2015-7450` | code | 1.0 |
| 4721 | Juniper ScreenOS Improper Authentication Vulnerability | kev | `cve-2015-7755` | code | 1.0 |
| 4722 | PHPMailer Command Injection Vulnerability | kev | `cve-2016-10033` | code | 1.0 |
| 4723 | NETGEAR WNR2000v5 Router Buffer Overflow Vulnerability | kev | `cve-2016-10174` | code | 1.0 |
| 4724 | Adobe Flash Player Arbitrary Code Execution Vulnerability | kev | `cve-2016-1019` | code | 1.0 |
| 4725 | NETGEAR Multiple WAP Devices Command Injection Vulnerability | kev | `cve-2016-1555` | code | 1.0 |
| 4726 | SAP NetWeaver SQL Injection Vulnerability | kev | `cve-2016-2386` | code | 1.0 |
| 4727 | Apache ActiveMQ Improper Input Validation Vulnerability | kev | `cve-2016-3088` | code | 1.0 |
| 4728 | Oracle Java SE and JRockit Unspecified Vulnerability | kev | `cve-2016-3427` | code | 1.0 |
| 4729 | Adobe Flash Player Arbitrary Code Execution Vulnerability | kev | `cve-2016-4117` | code | 1.0 |
| 4730 | Adobe Flash Player Remote Code Execution Vulnerability | kev | `cve-2016-4171` | code | 1.0 |
| 4731 | Apache Shiro Code Execution Vulnerability | kev | `cve-2016-4437` | code | 1.0 |
| 4732 | SKYSEA Client View Improper Authentication Vulnerability | kev | `cve-2016-7836` | code | 1.0 |
| 4733 | Apache Tomcat Remote Code Execution Vulnerability | kev | `cve-2016-8735` | code | 1.0 |
| 4734 | Jenkins Remote Code Execution Vulnerability | kev | `cve-2017-1000353` | code | 1.0 |
| 4735 | Primetek Primefaces Remote Code Execution Vulnerability | kev | `cve-2017-1000486` | code | 1.0 |
| 4736 | Telerik UI for ASP.NET AJAX Unrestricted File Upload Vulnerability | kev | `cve-2017-11317` | code | 1.0 |
| 4737 | Telerik UI for ASP.NET AJAX Insecure Direct Object Reference Vulnerability | kev | `cve-2017-11357` | code | 1.0 |
| 4738 | Red Hat JBoss Application Server Remote Code Execution Vulnerability | kev | `cve-2017-12149` | code | 1.0 |
| 4739 | Cisco IOS and IOS XE Software DHCP Remote Code Execution Vulnerability | kev | `cve-2017-12240` | code | 1.0 |
| 4740 | Palo Alto Networks PAN-OS Remote Code Execution Vulnerability | kev | `cve-2017-15944` | code | 1.0 |
| 4741 | Kaseya VSA SQL Injection Vulnerability | kev | `cve-2017-18362` | code | 1.0 |
| 4742 | Adobe ColdFusion Deserialization Vulnerability | kev | `cve-2017-3066` | code | 1.0 |
| 4743 | Cisco IOS and IOS XE Remote Code Execution Vulnerability | kev | `cve-2017-3881` | code | 1.0 |
| 4744 | Apache Struts Remote Code Execution Vulnerability | kev | `cve-2017-5638` | code | 1.0 |
| 4745 | Intel Active Management Technology (AMT), Small Business Technology (SBT), and Standard Manageability Privilege Escalation Vulnerability | kev | `cve-2017-5689` | code | 1.0 |
| 4746 | NETGEAR DGN2200 Remote Code Execution Vulnerability | kev | `cve-2017-6077` | code | 1.0 |
| 4747 | Citrix Multiple Products Remote Code Execution Vulnerability | kev | `cve-2017-6316` | code | 1.0 |
| 4748 | NETGEAR Multiple Devices Buffer Overflow Vulnerability | kev | `cve-2017-6862` | code | 1.0 |
| 4749 | Microsoft Windows Server Buffer Overflow Vulnerability | kev | `cve-2017-7269` | code | 1.0 |
| 4750 | Samba Remote Code Execution Vulnerability | kev | `cve-2017-7494` | code | 1.0 |
| 4751 | Hikvision Multiple Products Improper Authentication Vulnerability | kev | `cve-2017-7921` | code | 1.0 |
| 4752 | Microsoft Windows Search Remote Code Execution Vulnerability | kev | `cve-2017-8543` | code | 1.0 |
| 4753 | Progress Telerik UI for ASP.NET AJAX and Sitefinity Cryptographic Weakness Vulnerability | kev | `cve-2017-9248` | code | 1.0 |
| 4754 | Apache Struts 1 Improper Input Validation Vulnerability | kev | `cve-2017-9791` | code | 1.0 |
| 4755 | PHPUnit Command Injection Vulnerability | kev | `cve-2017-9841` | code | 1.0 |
| 4756 | Cisco VPN Routers Remote Code Execution Vulnerability | kev | `cve-2018-0125` | code | 1.0 |
| 4757 | Cisco Secure Access Control System Java Deserialization Vulnerability | kev | `cve-2018-0147` | code | 1.0 |
| 4758 | Cisco IOS Software and Cisco IOS XE Software Quality of Service Remote Code Execution Vulnerability | kev | `cve-2018-0151` | code | 1.0 |
| 4759 | Cisco IOS and IOS XE Software Smart Install Remote Code Execution Vulnerability | kev | `cve-2018-0171` | code | 1.0 |
| 4760 | Jenkins Stapler Web Framework Deserialization of Untrusted Data Vulnerability | kev | `cve-2018-1000861` | code | 1.0 |
| 4761 | Dasan GPON Routers Authentication Bypass Vulnerability | kev | `cve-2018-10561` | code | 1.0 |
| 4762 | Dasan GPON Routers Command Injection Vulnerability | kev | `cve-2018-10562` | code | 1.0 |
| 4763 | Quest KACE System Management Appliance Remote Command Execution Vulnerability | kev | `cve-2018-11138` | code | 1.0 |
| 4764 | VMware Tanzu Spring Data Commons Property Binder Vulnerability | kev | `cve-2018-1273` | code | 1.0 |
| 4765 | Tenda AC7, AC9, and AC10 Routers Command Injection Vulnerability | kev | `cve-2018-14558` | code | 1.0 |
| 4766 | Red Hat JBoss RichFaces Framework Expression Language Injection Vulnerability | kev | `cve-2018-14667` | code | 1.0 |
| 4767 | NUUO NVRmini Devices OS Command Injection Vulnerability | kev | `cve-2018-14933` | code | 1.0 |
| 4768 | Adobe ColdFusion Unrestricted File Upload Vulnerability | kev | `cve-2018-15961` | code | 1.0 |
| 4769 | GIGABYTE Multiple Products Privilege Escalation Vulnerability | kev | `cve-2018-19323` | code | 1.0 |
| 4770 | Paessler PRTG Network Monitor Local File Inclusion Vulnerability | kev | `cve-2018-19410` | code | 1.0 |
| 4771 | ThinkPHP "noneCms" Remote Code Execution Vulnerability | kev | `cve-2018-20062` | code | 1.0 |
| 4772 | Kaseya VSA Remote Code Execution Vulnerability | kev | `cve-2018-20753` | code | 1.0 |
| 4773 | Oracle WebLogic Server Unspecified Vulnerability | kev | `cve-2018-2628` | code | 1.0 |
| 4774 | Adobe ColdFusion Deserialization of Untrusted Data Vulnerability | kev | `cve-2018-4939` | code | 1.0 |
| 4775 | D-Link Multiple Routers OS Command Injection Vulnerability | kev | `cve-2018-6530` | code | 1.0 |
| 4776 | Exim Buffer Overflow Vulnerability | kev | `cve-2018-6789` | code | 1.0 |
| 4777 | MikroTik RouterOS Stack-Based Buffer Overflow Vulnerability | kev | `cve-2018-7445` | code | 1.0 |
| 4778 | Drupal Core Remote Code Execution Vulnerability | kev | `cve-2018-7600` | code | 1.0 |
| 4779 | Drupal Core Remote Code Execution Vulnerability | kev | `cve-2018-7602` | code | 1.0 |
| 4780 | Oracle Fusion Middleware Unspecified Vulnerability | kev | `cve-2012-3152` | code | 1.0 |
| 4781 | Arcserve Unified Data Protection (UDP) Directory Traversal Vulnerability | kev | `cve-2015-4068` | code | 1.0 |
| 4782 | MikroTik Router OS Directory Traversal Vulnerability | kev | `cve-2018-14847` | code | 1.0 |
| 4783 | Microsoft Windows Privilege Escalation Vulnerability | kev | `cve-2002-0367` | code | 1.0 |
| 4784 | Microsoft Windows Privilege Escalation Vulnerability | kev | `cve-2004-0210` | code | 1.0 |
| 4785 | Cisco IOS Denial-of-Service Vulnerability | kev | `cve-2004-1464` | code | 1.0 |
| 4786 | Apache Struts 1 ActionForm Denial-of-Service Vulnerability | kev | `cve-2006-1547` | code | 1.0 |
| 4787 | Microsoft Word Malformed Object Pointer Vulnerability | kev | `cve-2006-2492` | code | 1.0 |
| 4788 | Microsoft Office Excel Remote Code Execution Vulnerability | kev | `cve-2007-0671` | code | 1.0 |
| 4789 | Adobe Acrobat and Reader Buffer Overflow Vulnerability | kev | `cve-2007-5659` | code | 1.0 |
| 4790 | Microsoft Windows Video ActiveX Control Remote Code Execution Vulnerability | kev | `cve-2008-0015` | code | 1.0 |
| 4791 | Adobe Acrobat and Reader Unspecified Vulnerability | kev | `cve-2008-0655` | code | 1.0 |
| 4792 | Adobe Reader and Acrobat Input Validation Vulnerability | kev | `cve-2008-2992` | code | 1.0 |
| 4793 | Oracle VirtualBox Insufficient Input Validation Vulnerability | kev | `cve-2008-3431` | code | 1.0 |
| 4794 | Cisco IOS Cross-Site Request Forgery Vulnerability | kev | `cve-2008-4128` | code | 1.0 |
| 4795 | Microsoft Windows Buffer Overflow Vulnerability | kev | `cve-2008-4250` | code | 1.0 |
| 4796 | Microsoft Office Remote Code Execution | kev | `cve-2009-0238` | code | 1.0 |
| 4797 | Microsoft Office PowerPoint Code Injection Vulnerability | kev | `cve-2009-0556` | code | 1.0 |
| 4798 | Microsoft Office Object Record Corruption Vulnerability | kev | `cve-2009-0557` | code | 1.0 |
| 4799 | Microsoft Office Buffer Overflow Vulnerability | kev | `cve-2009-0563` | code | 1.0 |
| 4800 | Adobe Reader and Adobe Acrobat Stack-Based Buffer Overflow Vulnerability | kev | `cve-2009-0927` | code | 1.0 |
| 4801 | Microsoft Windows Improper Input Validation Vulnerability | kev | `cve-2009-1123` | code | 1.0 |
| 4802 | Microsoft DirectX NULL Byte Overwrite Vulnerability | kev | `cve-2009-1537` | code | 1.0 |
| 4803 | Adobe Acrobat and Reader, Flash Player Unspecified Vulnerability | kev | `cve-2009-1862` | code | 1.0 |
| 4804 | Cisco IOS XR Border Gateway Protocol (BGP) Denial-of-Service Vulnerability | kev | `cve-2009-2055` | code | 1.0 |
| 4805 | Microsoft Excel Featheader Record Memory Corruption Vulnerability | kev | `cve-2009-3129` | code | 1.0 |
| 4806 | Adobe Acrobat and Reader Heap-Based Buffer Overflow Vulnerability | kev | `cve-2009-3459` | code | 1.0 |
| 4807 | Adobe Acrobat and Reader Universal 3D Remote Code Execution Vulnerability | kev | `cve-2009-3953` | code | 1.0 |
| 4808 | Adobe BlazeDS Information Disclosure Vulnerability | kev | `cve-2009-3960` | code | 1.0 |
| 4809 | Adobe Acrobat and Reader Use-After-Free Vulnerability | kev | `cve-2009-4324` | code | 1.0 |
| 4810 | Adobe Reader and Acrobat Arbitrary Code Execution Vulnerability | kev | `cve-2010-0188` | code | 1.0 |
| 4811 | Microsoft Windows Kernel Exception Handler Vulnerability | kev | `cve-2010-0232` | code | 1.0 |
| 4812 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | `cve-2010-0249` | code | 1.0 |
| 4813 | Red Hat JBoss Authentication Bypass Vulnerability | kev | `cve-2010-0738` | code | 1.0 |
| 4814 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | `cve-2010-0806` | code | 1.0 |
| 4815 | Adobe Flash Player Memory Corruption Vulnerability | kev | `cve-2010-1297` | code | 1.0 |
| 4816 | Red Hat JBoss Information Disclosure Vulnerability | kev | `cve-2010-1428` | code | 1.0 |
| 4817 | Red Hat Linux JBoss Seam 2 Remote Code Execution Vulnerability | kev | `cve-2010-1871` | code | 1.0 |
| 4818 | Microsoft Windows Remote Code Execution Vulnerability | kev | `cve-2010-2568` | code | 1.0 |
| 4819 | Microsoft PowerPoint Buffer Overflow Vulnerability | kev | `cve-2010-2572` | code | 1.0 |
| 4820 | Adobe Acrobat and Reader Stack-Based Buffer Overflow Vulnerability | kev | `cve-2010-2883` | code | 1.0 |
| 4821 | Cisco IOS XR Border Gateway Protocol (BGP) Denial-of-Service Vulnerability | kev | `cve-2010-3035` | code | 1.0 |
| 4822 | Microsoft Office Stack-based Buffer Overflow Vulnerability | kev | `cve-2010-3333` | code | 1.0 |
| 4823 | Mozilla Multiple Products Remote Code Execution Vulnerability | kev | `cve-2010-3765` | code | 1.0 |
| 4824 | Linux Kernel Improper Input Validation Vulnerability | kev | `cve-2010-3904` | code | 1.0 |
| 4825 | Microsoft Internet Explorer Uninitialized Memory Corruption Vulnerability | kev | `cve-2010-3962` | code | 1.0 |
| 4826 | Exim Privilege Escalation Vulnerability | kev | `cve-2010-4345` | code | 1.0 |
| 4827 | Microsoft Windows Kernel Stack-Based Buffer Overflow Vulnerability | kev | `cve-2010-4398` | code | 1.0 |
| 4828 | Ubiquiti AirOS Command Injection Vulnerability | kev | `cve-2010-5330` | code | 1.0 |
| 4829 | Adobe Flash Player Unspecified Vulnerability | kev | `cve-2011-0609` | code | 1.0 |
| 4830 | Adobe Flash Player Remote Code Execution Vulnerability | kev | `cve-2011-0611` | code | 1.0 |
| 4831 | Android OS Privilege Escalation Vulnerability | kev | `cve-2011-1823` | code | 1.0 |
| 4832 | Microsoft Ancillary Function Driver (afd.sys) Improper Input Validation Vulnerability | kev | `cve-2011-2005` | code | 1.0 |
| 4833 | Microsoft Windows Remote Code Execution Vulnerability | kev | `cve-2011-3402` | code | 1.0 |
| 4834 | D-Link DIR-300 Router Cleartext Storage of a Password Vulnerability | kev | `cve-2011-4723` | code | 1.0 |
| 4835 | Microsoft Windows Authenticode Signature Verification Remote Code Execution Vulnerability | kev | `cve-2012-0151` | code | 1.0 |
| 4836 | Microsoft MSCOMCTL.OCX Remote Code Execution Vulnerability | kev | `cve-2012-0158` | code | 1.0 |
| 4837 | Oracle Java SE Runtime Environment (JRE) Arbitrary Code Execution Vulnerability | kev | `cve-2012-0507` | code | 1.0 |
| 4838 | Oracle Fusion Middleware Unspecified Vulnerability | kev | `cve-2012-0518` | code | 1.0 |
| 4839 | Adobe Flash Player Memory Corruption Vulnerability | kev | `cve-2012-0754` | code | 1.0 |
| 4840 | Adobe Flash Player Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2012-0767` | code | 1.0 |
| 4841 | Adobe Flash Player Arbitrary Code Execution Vulnerability | kev | `cve-2012-1535` | code | 1.0 |
| 4842 | Oracle Java SE Runtime Environment (JRE) Arbitrary Code Execution Vulnerability | kev | `cve-2012-1723` | code | 1.0 |
| 4843 | PHP-CGI Query String Parameter Vulnerability | kev | `cve-2012-1823` | code | 1.0 |
| 4844 | Microsoft Visual Basic for Applications Insecure Library Loading Vulnerability | kev | `cve-2012-1854` | code | 1.0 |
| 4845 | Microsoft Office MSCOMCTL.OCX Remote Code Execution Vulnerability | kev | `cve-2012-1856` | code | 1.0 |
| 4846 | Microsoft XML Core Services Memory Corruption Vulnerability | kev | `cve-2012-1889` | code | 1.0 |
| 4847 | Adobe Flash Player Memory Corruption Vulnerability | kev | `cve-2012-2034` | code | 1.0 |
| 4848 | Microsoft Word Remote Code Execution Vulnerability | kev | `cve-2012-2539` | code | 1.0 |
| 4849 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | `cve-2012-4792` | code | 1.0 |
| 4850 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | `cve-2012-4969` | code | 1.0 |
| 4851 | Adobe Flash Player Integer Overflow Vulnerability | kev | `cve-2012-5054` | code | 1.0 |
| 4852 | Oracle Java SE Sandbox Bypass Vulnerability | kev | `cve-2012-5076` | code | 1.0 |
| 4853 | Microsoft Silverlight Double Dereference Vulnerability | kev | `cve-2013-0074` | code | 1.0 |
| 4854 | Oracle JRE Remote Code Execution Vulnerability | kev | `cve-2013-0422` | code | 1.0 |
| 4855 | Oracle JRE Sandbox Bypass Vulnerability | kev | `cve-2013-0431` | code | 1.0 |
| 4856 | Adobe ColdFusion Directory Traversal Vulnerability | kev | `cve-2013-0629` | code | 1.0 |
| 4857 | Adobe ColdFusion Information Disclosure Vulnerability | kev | `cve-2013-0631` | code | 1.0 |
| 4858 | Adobe Reader and Acrobat Memory Corruption Vulnerability | kev | `cve-2013-0640` | code | 1.0 |
| 4859 | Adobe Reader Buffer Overflow Vulnerability | kev | `cve-2013-0641` | code | 1.0 |
| 4860 | Adobe Flash Player Incorrect Default Permissions Vulnerability | kev | `cve-2013-0643` | code | 1.0 |
| 4861 | Adobe Flash Player Code Execution Vulnerability | kev | `cve-2013-0648` | code | 1.0 |
| 4862 | Microsoft Office Buffer Overflow Vulnerability | kev | `cve-2013-1331` | code | 1.0 |
| 4863 | Microsoft Internet Explorer Remote Code Execution Vulnerability | kev | `cve-2013-1347` | code | 1.0 |
| 4864 | Mozilla Firefox Information Disclosure Vulnerability | kev | `cve-2013-1675` | code | 1.0 |
| 4865 | Mozilla Firefox and Thunderbird Denial-of-Service Vulnerability | kev | `cve-2013-1690` | code | 1.0 |
| 4866 | Linux Kernel Privilege Escalation Vulnerability | kev | `cve-2013-2094` | code | 1.0 |
| 4867 | Oracle JRE Unspecified Vulnerability | kev | `cve-2013-2423` | code | 1.0 |
| 4868 | Oracle Java SE Unspecified Vulnerability | kev | `cve-2013-2465` | code | 1.0 |
| 4869 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | `cve-2013-2551` | code | 1.0 |
| 4870 | Linux Kernel Integer Overflow Vulnerability | kev | `cve-2013-2596` | code | 1.0 |
| 4871 | Code Aurora ACDB Audio Driver Stack-based Buffer Overflow Vulnerability | kev | `cve-2013-2597` | code | 1.0 |
| 4872 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | `cve-2013-3163` | code | 1.0 |
| 4873 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2013-3660` | code | 1.0 |
| 4874 | Microsoft Internet Explorer Resource Management Errors Vulnerability | kev | `cve-2013-3893` | code | 1.0 |
| 4875 | Microsoft Silverlight Information Disclosure Vulnerability | kev | `cve-2013-3896` | code | 1.0 |
| 4876 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | `cve-2013-3897` | code | 1.0 |
| 4877 | Microsoft WinVerifyTrust function Remote Code Execution | kev | `cve-2013-3900` | code | 1.0 |
| 4878 | Microsoft Graphics Component Memory Corruption Vulnerability | kev | `cve-2013-3906` | code | 1.0 |
| 4879 | Microsoft Windows Out-of-Bounds Write Vulnerability | kev | `cve-2013-3918` | code | 1.0 |
| 4880 | IBM InfoSphere BigInsights Invalid Input Vulnerability | kev | `cve-2013-3993` | code | 1.0 |
| 4881 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | `cve-2013-5065` | code | 1.0 |
| 4882 | D-Link DSL-2760U Gateway Cross-Site Scripting Vulnerability | kev | `cve-2013-5223` | code | 1.0 |
| 4883 | Linux Kernel Improper Input Validation Vulnerability | kev | `cve-2013-6282` | code | 1.0 |
| 4884 | Microsoft Internet Explorer Information Disclosure Vulnerability | kev | `cve-2013-7331` | code | 1.0 |
| 4885 | Ruby on Rails Directory Traversal Vulnerability | kev | `cve-2014-0130` | code | 1.0 |
| 4886 | OpenSSL Information Disclosure Vulnerability | kev | `cve-2014-0160` | code | 1.0 |
| 4887 | Linux Kernel Race Condition Vulnerability | kev | `cve-2014-0196` | code | 1.0 |
| 4888 | Microsoft Internet Explorer Use-After-Free Vulnerability | kev | `cve-2014-0322` | code | 1.0 |
| 4889 | Adobe Reader and Acrobat Use-After-Free Vulnerability | kev | `cve-2014-0496` | code | 1.0 |
| 4890 | Adobe Flash Player Double Free Vulnerablity | kev | `cve-2014-0502` | code | 1.0 |
| 4891 | D-Link DIR-600 Router Cross-Site Request Forgery (CSRF) Vulnerability | kev | `cve-2014-100005` | code | 1.0 |
| 4892 | Microsoft Word Memory Corruption Vulnerability | kev | `cve-2014-1761` | code | 1.0 |
| 4893 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | `cve-2014-1776` | code | 1.0 |
| 4894 | Microsoft Windows Group Policy Preferences Password Privilege Escalation Vulnerability | kev | `cve-2014-1812` | code | 1.0 |
| 4895 | Cisco Adaptive Security Appliance (ASA) Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2014-2120` | code | 1.0 |
| 4896 | Microsoft Internet Explorer Privilege Escalation Vulnerability | kev | `cve-2014-2817` | code | 1.0 |
| 4897 | Elasticsearch Remote Code Execution Vulnerability | kev | `cve-2014-3120` | code | 1.0 |
| 4898 | Linux Kernel Privilege Escalation Vulnerability | kev | `cve-2014-3153` | code | 1.0 |
| 4899 | Microsoft IME Japanese Privilege Escalation Vulnerability | kev | `cve-2014-4077` | code | 1.0 |
| 4900 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2014-4113` | code | 1.0 |
| 4901 | Microsoft Windows Object Linking & Embedding (OLE) Remote Code Execution Vulnerability | kev | `cve-2014-4114` | code | 1.0 |
| 4902 | Microsoft Internet Explorer Privilege Escalation Vulnerability | kev | `cve-2014-4123` | code | 1.0 |
| 4903 | Microsoft Windows Remote Code Execution Vulnerability | kev | `cve-2014-4148` | code | 1.0 |
| 4904 | Apple OS X Heap-Based Buffer Overflow Vulnerability | kev | `cve-2014-4404` | code | 1.0 |
| 4905 | GNU Bash OS Command Injection Vulnerability | kev | `cve-2014-6278` | code | 1.0 |
| 4906 | Microsoft Kerberos Key Distribution Center (KDC) Privilege Escalation Vulnerability | kev | `cve-2014-6324` | code | 1.0 |
| 4907 | Microsoft Windows Object Linking & Embedding (OLE) Automation Array Remote Code Execution Vulnerability | kev | `cve-2014-6332` | code | 1.0 |
| 4908 | Microsoft Windows Code Injection Vulnerability | kev | `cve-2014-6352` | code | 1.0 |
| 4909 | Adobe Flash Player Dereferenced Pointer Vulnerability | kev | `cve-2014-8439` | code | 1.0 |
| 4910 | Adobe Flash Player Stack-Based Buffer Overflow Vulnerability | kev | `cve-2014-9163` | code | 1.0 |
| 4911 | Microsoft Windows TS WebProxy Directory Traversal Vulnerability | kev | `cve-2015-0016` | code | 1.0 |
| 4912 | Microsoft Internet Explorer ASLR Bypass Vulnerability | kev | `cve-2015-0071` | code | 1.0 |
| 4913 | Adobe Flash Player ASLR Bypass Vulnerability | kev | `cve-2015-0310` | code | 1.0 |
| 4914 | Cisco Prime Data Center Network Manager (DCNM) Directory Traversal Vulnerability | kev | `cve-2015-0666` | code | 1.0 |
| 4915 | Apple OS X Authentication Bypass Vulnerability | kev | `cve-2015-1130` | code | 1.0 |
| 4916 | Microsoft HTTP.sys Remote Code Execution Vulnerability | kev | `cve-2015-1635` | code | 1.0 |
| 4917 | Microsoft Office Memory Corruption Vulnerability | kev | `cve-2015-1641` | code | 1.0 |
| 4918 | Microsoft Office Memory Corruption Vulnerability | kev | `cve-2015-1642` | code | 1.0 |
| 4919 | Microsoft Windows Remote Code Execution Vulnerability | kev | `cve-2015-1671` | code | 1.0 |
| 4920 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2015-1701` | code | 1.0 |
| 4921 | Microsoft Windows Mount Manager Privilege Escalation Vulnerability | kev | `cve-2015-1769` | code | 1.0 |
| 4922 | Microsoft Office Uninitialized Memory Use Vulnerability | kev | `cve-2015-1770` | code | 1.0 |
| 4923 | D-Link DIR-645 Router Remote Code Execution Vulnerability | kev | `cve-2015-2051` | code | 1.0 |
| 4924 | Intel Ethernet Diagnostics Driver for Windows Denial-of-Service Vulnerability | kev | `cve-2015-2291` | code | 1.0 |
| 4925 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2015-2360` | code | 1.0 |
| 4926 | Microsoft ATM Font Driver Privilege Escalation Vulnerability | kev | `cve-2015-2387` | code | 1.0 |
| 4927 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | `cve-2015-2419` | code | 1.0 |
| 4928 | Microsoft PowerPoint Memory Corruption Vulnerability | kev | `cve-2015-2424` | code | 1.0 |
| 4929 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | `cve-2015-2425` | code | 1.0 |
| 4930 | Microsoft Windows Adobe Type Manager Library Remote Code Execution Vulnerability | kev | `cve-2015-2426` | code | 1.0 |
| 4931 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | `cve-2015-2502` | code | 1.0 |
| 4932 | Microsoft Office Malformed EPS File Vulnerability | kev | `cve-2015-2545` | code | 1.0 |
| 4933 | Microsoft Win32k Memory Corruption Vulnerability | kev | `cve-2015-2546` | code | 1.0 |
| 4934 | TP-Link Multiple Archer Devices Directory Traversal Vulnerability | kev | `cve-2015-3035` | code | 1.0 |
| 4935 | Red Hat Libuser Race Condition Vulnerability | kev | `cve-2015-3246` | code | 1.0 |
| 4936 | ProFTPD Improper Access Control Vulnerability | kev | `cve-2015-3306` | code | 1.0 |
| 4937 | Mozilla Firefox Security Feature Bypass Vulnerability | kev | `cve-2015-4495` | code | 1.0 |
| 4938 | Oracle Java SE Integrity Check Vulnerability | kev | `cve-2015-4902` | code | 1.0 |
| 4939 | Red Hat Automatic Bug Reporting Tool Privilege Escalation Vulnerability | kev | `cve-2015-5287` | code | 1.0 |
| 4940 | Jenkins User Interface (UI) Information Disclosure Vulnerability | kev | `cve-2015-5317` | code | 1.0 |
| 4941 | ISC BIND Data Processing Errors Vulnerability | kev | `cve-2015-5477` | code | 1.0 |
| 4942 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | `cve-2015-6175` | code | 1.0 |
| 4943 | Adobe Flash Player Arbitrary Code Execution Vulnerability | kev | `cve-2015-7645` | code | 1.0 |
| 4944 | Adobe Flash Player Integer Overflow Vulnerability | kev | `cve-2015-8651` | code | 1.0 |
| 4945 | Microsoft Silverlight Runtime Remote Code Execution Vulnerability | kev | `cve-2016-0034` | code | 1.0 |
| 4946 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | `cve-2016-0040` | code | 1.0 |
| 4947 | Microsoft Windows Secondary Logon Service Privilege Escalation Vulnerability | kev | `cve-2016-0099` | code | 1.0 |
| 4948 | Microsoft Windows CSRSS Security Feature Bypass Vulnerability | kev | `cve-2016-0151` | code | 1.0 |
| 4949 | Microsoft Internet Explorer Information Disclosure Vulnerability | kev | `cve-2016-0162` | code | 1.0 |
| 4950 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2016-0165` | code | 1.0 |
| 4951 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2016-0167` | code | 1.0 |
| 4952 | Microsoft Windows Media Center Remote Code Execution Vulnerability | kev | `cve-2016-0185` | code | 1.0 |
| 4953 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | `cve-2016-0189` | code | 1.0 |
| 4954 | Ruby on Rails Directory Traversal Vulnerability | kev | `cve-2016-0752` | code | 1.0 |
| 4955 | Adobe Flash Player and AIR Use-After-Free Vulnerability | kev | `cve-2016-0984` | code | 1.0 |
| 4956 | Adobe Flash Player and AIR Integer Overflow Vulnerability | kev | `cve-2016-1010` | code | 1.0 |
| 4957 | D-Link DCS-930L Devices OS Command Injection Vulnerability | kev | `cve-2016-11021` | code | 1.0 |
| 4958 | Google Chromium V8 Out-of-Bounds Read Vulnerability | kev | `cve-2016-1646` | code | 1.0 |
| 4959 | D-Link DSL-2750B Devices Command Injection Vulnerability | kev | `cve-2016-20017` | code | 1.0 |
| 4960 | SAP NetWeaver Information Disclosure Vulnerability | kev | `cve-2016-2388` | code | 1.0 |
| 4961 | Apache Struts Command Injection Vulnerability | kev | `cve-2016-3081` | code | 1.0 |
| 4962 | Microsoft Office OLE DLL Side Loading Vulnerability | kev | `cve-2016-3235` | code | 1.0 |
| 4963 | Microsoft Internet Explorer Messaging API Information Disclosure Vulnerability | kev | `cve-2016-3298` | code | 1.0 |
| 4964 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | `cve-2016-3309` | code | 1.0 |
| 4965 | Microsoft Internet Explorer and Edge Information Disclosure Vulnerability | kev | `cve-2016-3351` | code | 1.0 |
| 4966 | Microsoft Windows Graphics Device Interface (GDI) Remote Code Execution Vulnerability | kev | `cve-2016-3393` | code | 1.0 |
| 4967 | SolarWinds Virtualization Manager Privilege Escalation Vulnerability | kev | `cve-2016-3643` | code | 1.0 |
| 4968 | ImageMagick Improper Input Validation Vulnerability | kev | `cve-2016-3714` | code | 1.0 |
| 4969 | ImageMagick Arbitrary File Deletion Vulnerability | kev | `cve-2016-3715` | code | 1.0 |
| 4970 | ImageMagick Server-Side Request Forgery (SSRF) Vulnerability | kev | `cve-2016-3718` | code | 1.0 |
| 4971 | SAP NetWeaver Directory Traversal Vulnerability | kev | `cve-2016-3976` | code | 1.0 |
| 4972 | Trihedral VTScada (formerly VTS) Denial-of-Service Vulnerability | kev | `cve-2016-4523` | code | 1.0 |
| 4973 | Apple iOS Information Disclosure Vulnerability | kev | `cve-2016-4655` | code | 1.0 |
| 4974 | Apple iOS Memory Corruption Vulnerability | kev | `cve-2016-4656` | code | 1.0 |
| 4975 | Apple iOS Webkit Memory Corruption Vulnerability | kev | `cve-2016-4657` | code | 1.0 |
| 4976 | Linux Kernel Race Condition Vulnerability | kev | `cve-2016-5195` | code | 1.0 |
| 4977 | Google Chromium V8 Out-of-Bounds Memory Vulnerability | kev | `cve-2016-5198` | code | 1.0 |
| 4978 | NETGEAR Multiple Routers Remote Code Execution Vulnerability | kev | `cve-2016-6277` | code | 1.0 |
| 4979 | Cisco Adaptive Security Appliance (ASA) SNMP Buffer Overflow Vulnerability | kev | `cve-2016-6366` | code | 1.0 |
| 4980 | Cisco Adaptive Security Appliance (ASA) CLI Remote Code Execution Vulnerability | kev | `cve-2016-6367` | code | 1.0 |
| 4981 | Cisco IOS, IOS XR, and IOS XE IKEv1 Information Disclosure Vulnerability | kev | `cve-2016-6415` | code | 1.0 |
| 4982 | Microsoft Office Memory Corruption Vulnerability | kev | `cve-2016-7193` | code | 1.0 |
| 4983 | Microsoft Edge Memory Corruption Vulnerability | kev | `cve-2016-7200` | code | 1.0 |
| 4984 | Microsoft Edge Memory Corruption Vulnerability | kev | `cve-2016-7201` | code | 1.0 |
| 4985 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2016-7255` | code | 1.0 |
| 4986 | Microsoft Windows Open Type Font Remote Code Execution Vulnerability | kev | `cve-2016-7256` | code | 1.0 |
| 4987 | Microsoft Office Security Feature Bypass Vulnerability | kev | `cve-2016-7262` | code | 1.0 |
| 4988 | Adobe Flash Player Use-After-Free Vulnerability | kev | `cve-2016-7855` | code | 1.0 |
| 4989 | Adobe Flash Player Use-After-Free Vulnerability | kev | `cve-2016-7892` | code | 1.0 |
| 4990 | Siemens SIMATIC CP 1543-1 Improper Privilege Management Vulnerability | kev | `cve-2016-8562` | code | 1.0 |
| 4991 | Mozilla Firefox, Firefox ESR, and Thunderbird Use-After-Free Vulnerability | kev | `cve-2016-9079` | code | 1.0 |
| 4992 | SAP NetWeaver XML External Entity (XXE) Vulnerability | kev | `cve-2016-9563` | code | 1.0 |
| 4993 | Microsoft Graphics Device Interface (GDI) Privilege Escalation Vulnerability | kev | `cve-2017-0001` | code | 1.0 |
| 4994 | Microsoft Windows Graphics Device Interface (GDI) Privilege Escalation Vulnerability | kev | `cve-2017-0005` | code | 1.0 |
| 4995 | Microsoft XML Core Services Information Disclosure Vulnerability | kev | `cve-2017-0022` | code | 1.0 |
| 4996 | Microsoft Edge and Internet Explorer Type Confusion Vulnerability | kev | `cve-2017-0037` | code | 1.0 |
| 4997 | Microsoft Internet Explorer Information Disclosure Vulnerability | kev | `cve-2017-0059` | code | 1.0 |
| 4998 | Microsoft Windows Transaction Manager Privilege Escalation Vulnerability | kev | `cve-2017-0101` | code | 1.0 |
| 4999 | Microsoft Windows Server Message Block (SMBv1) Remote Code Execution Vulnerability | kev | `cve-2017-0143` | code | 1.0 |
| 5000 | Microsoft SMBv1 Remote Code Execution Vulnerability | kev | `cve-2017-0144` | code | 1.0 |
| 5001 | Microsoft SMBv1 Remote Code Execution Vulnerability | kev | `cve-2017-0145` | code | 1.0 |
| 5002 | Microsoft Windows SMB Remote Code Execution Vulnerability | kev | `cve-2017-0146` | code | 1.0 |
| 5003 | Microsoft Windows SMBv1 Information Disclosure Vulnerability | kev | `cve-2017-0147` | code | 1.0 |
| 5004 | Microsoft SMBv1 Server Remote Code Execution Vulnerability | kev | `cve-2017-0148` | code | 1.0 |
| 5005 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | `cve-2017-0149` | code | 1.0 |
| 5006 | Microsoft Office and WordPad Remote Code Execution Vulnerability | kev | `cve-2017-0199` | code | 1.0 |
| 5007 | Microsoft Internet Explorer Privilege Escalation Vulnerability | kev | `cve-2017-0210` | code | 1.0 |
| 5008 | Microsoft Windows Privilege Escalation Vulnerability | kev | `cve-2017-0213` | code | 1.0 |
| 5009 | Microsoft Internet Explorer Remote Code Execution Vulnerability | kev | `cve-2017-0222` | code | 1.0 |
| 5010 | Microsoft Office Use-After-Free Vulnerability | kev | `cve-2017-0261` | code | 1.0 |
| 5011 | Microsoft Office Remote Code Execution Vulnerability | kev | `cve-2017-0262` | code | 1.0 |
| 5012 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2017-0263` | code | 1.0 |
| 5013 | Linux Kernel PIE Stack Buffer Corruption Vulnerability | kev | `cve-2017-1000253` | code | 1.0 |
| 5014 | Oracle Corporation WebLogic Server Remote Code Execution Vulnerability | kev | `cve-2017-10271` | code | 1.0 |
| 5015 | Adobe Flash Player Type Confusion Vulnerability | kev | `cve-2017-11292` | code | 1.0 |
| 5016 | Microsoft Office Outlook Security Feature Bypass Vulnerability | kev | `cve-2017-11774` | code | 1.0 |
| 5017 | Microsoft Office Remote Code Execution Vulnerability | kev | `cve-2017-11826` | code | 1.0 |
| 5018 | Microsoft Office Memory Corruption Vulnerability | kev | `cve-2017-11882` | code | 1.0 |
| 5019 | Cisco IOS Software Network Address Translation Denial-of-Service Vulnerability | kev | `cve-2017-12231` | code | 1.0 |
| 5020 | Cisco IOS Software for Cisco Integrated Services Routers Denial-of-Service Vulnerability | kev | `cve-2017-12232` | code | 1.0 |
| 5021 | Cisco IOS Software Common Industrial Protocol Request Denial-of-Service Vulnerability | kev | `cve-2017-12233` | code | 1.0 |
| 5022 | Cisco IOS Software Common Industrial Protocol Request Denial-of-Service Vulnerability | kev | `cve-2017-12234` | code | 1.0 |
| 5023 | Cisco IOS Software for Cisco Industrial Ethernet Switches PROFINET Denial-of-Service Vulnerability | kev | `cve-2017-12235` | code | 1.0 |
| 5024 | Cisco IOS and IOS XE Software Internet Key Exchange Denial-of-Service Vulnerability | kev | `cve-2017-12237` | code | 1.0 |
| 5025 | Cisco Catalyst 6800 Series Switches VPLS Denial-of-Service Vulnerability | kev | `cve-2017-12238` | code | 1.0 |
| 5026 | Cisco IOS XE Software Ethernet Virtual Private Network Border Gateway Protocol Denial-of-Service Vulnerability | kev | `cve-2017-12319` | code | 1.0 |
| 5027 | Apache Tomcat on Windows Remote Code Execution Vulnerability | kev | `cve-2017-12615` | code | 1.0 |
| 5028 | Apache Tomcat Remote Code Execution Vulnerability | kev | `cve-2017-12617` | code | 1.0 |
| 5029 | SAP NetWeaver Directory Traversal Vulnerability | kev | `cve-2017-12637` | code | 1.0 |
| 5030 | Roundcube Webmail File Disclosure Vulnerability | kev | `cve-2017-16651` | code | 1.0 |
| 5031 | Embedthis GoAhead Remote Code Execution Vulnerability | kev | `cve-2017-17562` | code | 1.0 |
| 5032 | Zyxel P660HN-T1A Routers Command Injection Vulnerability | kev | `cve-2017-18368` | code | 1.0 |
| 5033 | Oracle WebLogic Server OS Command Injection Vulnerability | kev | `cve-2017-3506` | code | 1.0 |
| 5034 | Google Chromium V8 Memory Corruption Vulnerability | kev | `cve-2017-5030` | code | 1.0 |
| 5035 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2017-5070` | code | 1.0 |
| 5036 | NETGEAR Multiple Devices Exposure of Sensitive Information Vulnerability | kev | `cve-2017-5521` | code | 1.0 |
| 5037 | Symantec Messaging Gateway Remote Code Execution Vulnerability | kev | `cve-2017-6327` | code | 1.0 |
| 5038 | NETGEAR DGN2200 Devices OS Command Injection Vulnerability | kev | `cve-2017-6334` | code | 1.0 |
| 5039 | Cisco IOS Software and Cisco IOS XE Software UDP Packet Processing Denial-of-Service Vulnerability | kev | `cve-2017-6627` | code | 1.0 |
| 5040 | Cisco IOS Software and Cisco IOS XE Software Denial-of-Service Vulnerability | kev | `cve-2017-6663` | code | 1.0 |
| 5041 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | `cve-2017-6736` | code | 1.0 |
| 5042 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | `cve-2017-6737` | code | 1.0 |
| 5043 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | `cve-2017-6738` | code | 1.0 |
| 5044 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | `cve-2017-6739` | code | 1.0 |
| 5045 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | `cve-2017-6740` | code | 1.0 |
| 5046 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | `cve-2017-6742` | code | 1.0 |
| 5047 | Cisco IOS and IOS XE Software SNMP Remote Code Execution Vulnerability | kev | `cve-2017-6743` | code | 1.0 |
| 5048 | Cisco IOS Software SNMP Remote Code Execution Vulnerability | kev | `cve-2017-6744` | code | 1.0 |
| 5049 | Zyxel EMG2926 Routers Command Injection Vulnerability | kev | `cve-2017-6884` | code | 1.0 |
| 5050 | Artifex Ghostscript Type Confusion Vulnerability | kev | `cve-2017-8291` | code | 1.0 |
| 5051 | Microsoft Windows Shell (.lnk) Remote Code Execution Vulnerability | kev | `cve-2017-8464` | code | 1.0 |
| 5052 | Microsoft Malware Protection Engine Improper Restriction of Operations Vulnerability | kev | `cve-2017-8540` | code | 1.0 |
| 5053 | Microsoft Office Remote Code Execution Vulnerability | kev | `cve-2017-8570` | code | 1.0 |
| 5054 | Microsoft .NET Framework Remote Code Execution Vulnerability | kev | `cve-2017-8759` | code | 1.0 |
| 5055 | Apache Struts Deserialization of Untrusted Data Vulnerability | kev | `cve-2017-9805` | code | 1.0 |
| 5056 | DotNetNuke (DNN) Remote Code Execution Vulnerability | kev | `cve-2017-9822` | code | 1.0 |
| 5057 | Cisco IOS Software Integrated Services Module for VPN Denial-of-Service Vulnerability | kev | `cve-2018-0154` | code | 1.0 |
| 5058 | Cisco Catalyst Bidirectional Forwarding Detection Denial-of-Service Vulnerability | kev | `cve-2018-0155` | code | 1.0 |
| 5059 | Cisco IOS Software and Cisco IOS XE Software Smart Install Denial-of-Service Vulnerability | kev | `cve-2018-0156` | code | 1.0 |
| 5060 | Cisco IOS and XE Software Internet Key Exchange Memory Leak Vulnerability | kev | `cve-2018-0158` | code | 1.0 |
| 5061 | Cisco IOS and XE Software Internet Key Exchange Version 1 Denial-of-Service Vulnerability | kev | `cve-2018-0159` | code | 1.0 |
| 5062 | Cisco IOS Software Resource Management Errors Vulnerability | kev | `cve-2018-0161` | code | 1.0 |
| 5063 | Cisco IOS, XR, and XE Software Buffer Overflow Vulnerability | kev | `cve-2018-0167` | code | 1.0 |
| 5064 | Cisco IOS and IOS XE Software Improper Input Validation Vulnerability | kev | `cve-2018-0172` | code | 1.0 |
| 5065 | Cisco IOS and IOS XE Software Improper Input Validation Vulnerability | kev | `cve-2018-0173` | code | 1.0 |
| 5066 | Cisco IOS Software and Cisco IOS XE Software Improper Input Validation Vulnerability | kev | `cve-2018-0174` | code | 1.0 |
| 5067 | Cisco IOS, XR, and XE Software Buffer Overflow Vulnerability | kev | `cve-2018-0175` | code | 1.0 |
| 5068 | Cisco IOS Software Denial-of-Service Vulnerability | kev | `cve-2018-0179` | code | 1.0 |
| 5069 | Cisco IOS Software Denial-of-Service Vulnerability | kev | `cve-2018-0180` | code | 1.0 |
| 5070 | Cisco Adaptive Security Appliance (ASA) Denial-of-Service Vulnerability | kev | `cve-2018-0296` | code | 1.0 |
| 5071 | Microsoft Office Memory Corruption Vulnerability | kev | `cve-2018-0798` | code | 1.0 |
| 5072 | Microsoft Office Memory Corruption Vulnerability | kev | `cve-2018-0802` | code | 1.0 |
| 5073 | Microsoft COM for Windows Deserialization of Untrusted Data Vulnerability | kev | `cve-2018-0824` | code | 1.0 |
| 5074 | Apache Struts Remote Code Execution Vulnerability | kev | `cve-2018-11776` | code | 1.0 |
| 5075 | Fortinet FortiOS and FortiADC Improper Access Control Vulnerability | kev | `cve-2018-13374` | code | 1.0 |
| 5076 | Fortinet FortiOS SSL VPN Path Traversal Vulnerability | kev | `cve-2018-13379` | code | 1.0 |
| 5077 | Fortinet FortiOS and FortiProxy Improper Authorization | kev | `cve-2018-13382` | code | 1.0 |
| 5078 | Fortinet FortiOS and FortiProxy Out-of-bounds Write | kev | `cve-2018-13383` | code | 1.0 |
| 5079 | Linux Kernel Integer Overflow Vulnerability | kev | `cve-2018-14634` | code | 1.0 |
| 5080 | LG N1A1 NAS Remote Command Execution Vulnerability | kev | `cve-2018-14839` | code | 1.0 |
| 5081 | Laravel Deserialization of Untrusted Data Vulnerability | kev | `cve-2018-15133` | code | 1.0 |
| 5082 | DotNetNuke (DNN) Inadequate Encryption Strength Vulnerability | kev | `cve-2018-15811` | code | 1.0 |
| 5083 | Adobe Flash Player Use-After-Free Vulnerability | kev | `cve-2018-15982` | code | 1.0 |
| 5084 | Google Chromium V8 Remote Code Execution Vulnerability | kev | `cve-2018-17463` | code | 1.0 |
| 5085 | Google Chromium V8 Out-of-Bounds Write Vulnerability | kev | `cve-2018-17480` | code | 1.0 |
| 5086 | DotNetNuke (DNN) Inadequate Encryption Strength Vulnerability | kev | `cve-2018-18325` | code | 1.0 |
| 5087 | TIBCO JasperReports Library Directory Traversal Vulnerability | kev | `cve-2018-18809` | code | 1.0 |
| 5088 | GIGABYTE Multiple Products Unspecified Vulnerability | kev | `cve-2018-19320` | code | 1.0 |
| 5089 | GIGABYTE Multiple Products Privilege Escalation Vulnerability | kev | `cve-2018-19321` | code | 1.0 |
| 5090 | GIGABYTE Multiple Products Code Execution Vulnerability | kev | `cve-2018-19322` | code | 1.0 |
| 5091 | QNAP NAS File Station Cross-Site Scripting Vulnerability | kev | `cve-2018-19943` | code | 1.0 |
| 5092 | QNAP NAS File Station Command Injection Vulnerability | kev | `cve-2018-19949` | code | 1.0 |
| 5093 | QNAP NAS File Station Cross-Site Scripting Vulnerability | kev | `cve-2018-19953` | code | 1.0 |
| 5094 | WinRAR Absolute Path Traversal Vulnerability | kev | `cve-2018-20250` | code | 1.0 |
| 5095 | SAP Customer Relationship Management (CRM) Path Traversal Vulnerability | kev | `cve-2018-2380` | code | 1.0 |
| 5096 | Sierra Wireless AirLink ALEOS Unrestricted Upload of File with Dangerous Type Vulnerability | kev | `cve-2018-4063` | code | 1.0 |
| 5097 | Apple Multiple Products Memory Corruption Vulnerability | kev | `cve-2018-4344` | code | 1.0 |
| 5098 | Adobe Flash Player Use-After-Free Vulnerability | kev | `cve-2018-4878` | code | 1.0 |
| 5099 | Adobe Acrobat and Reader Double Free Vulnerability | kev | `cve-2018-4990` | code | 1.0 |
| 5100 | Adobe Flash Player Stack-based Buffer Overflow Vulnerability | kev | `cve-2018-5002` | code | 1.0 |
| 5101 | TIBCO JasperReports Server Information Disclosure Vulnerability | kev | `cve-2018-5430` | code | 1.0 |
| 5102 | Google Chromium V8 Integer Overflow Vulnerability | kev | `cve-2018-6065` | code | 1.0 |
| 5103 | Synacor Zimbra Collaboration Suite (ZCS) Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2018-6882` | code | 1.0 |
| 5104 | VMware SD-WAN Edge by VeloCloud Command Injection Vulnerability | kev | `cve-2018-6961` | code | 1.0 |
| 5105 | Schneider Electric U.motion Builder SQL Injection Vulnerability | kev | `cve-2018-7841` | code | 1.0 |
| 5106 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2018-8120` | code | 1.0 |
| 5107 | Microsoft Windows VBScript Engine Out-of-Bounds Write Vulnerability | kev | `cve-2018-8174` | code | 1.0 |
| 5108 | ChakraCore Scripting Engine Type Confusion Vulnerability | kev | `cve-2018-8298` | code | 1.0 |
| 5109 | Microsoft Scripting Engine Memory Corruption Vulnerability | kev | `cve-2018-8373` | code | 1.0 |
| 5110 | Microsoft DirectX Graphics Kernel Privilege Escalation Vulnerability | kev | `cve-2018-8405` | code | 1.0 |
| 5111 | Microsoft DirectX Graphics Kernel Privilege Escalation Vulnerability | kev | `cve-2018-8406` | code | 1.0 |
| 5112 | Microsoft Windows Shell Remote Code Execution Vulnerability | kev | `cve-2018-8414` | code | 1.0 |
| 5113 | Microsoft Windows Privilege Escalation Vulnerability | kev | `cve-2018-8440` | code | 1.0 |
| 5114 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2018-8453` | code | 1.0 |
| 5115 | Microsoft Exchange Server Privilege Escalation Vulnerability | kev | `cve-2018-8581` | code | 1.0 |
| 5116 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2018-8589` | code | 1.0 |
| 5117 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | `cve-2018-8611` | code | 1.0 |
| 5118 | Microsoft Windows Win32k Improper Resource Shutdown or Release Vulnerability | kev | `cve-2018-8639` | code | 1.0 |
| 5119 | Microsoft Internet Explorer Scripting Engine Memory Corruption Vulnerability | kev | `cve-2018-8653` | code | 1.0 |
| 5120 | Paessler PRTG Network Monitor OS Command Injection Vulnerability | kev | `cve-2018-9276` | code | 1.0 |
| 5121 | Apache Solr DataImportHandler Code Injection Vulnerability | kev | `cve-2019-0193` | code | 1.0 |
| 5122 | Apache HTTP Server Privilege Escalation Vulnerability | kev | `cve-2019-0211` | code | 1.0 |
| 5123 | SAP Commerce Cloud Deserialization of Untrusted Data Vulnerability | kev | `cve-2019-0344` | code | 1.0 |
| 5124 | Microsoft MSHTML Remote Code Execution Vulnerability | kev | `cve-2019-0541` | code | 1.0 |
| 5125 | Microsoft Windows Privilege Escalation Vulnerability | kev | `cve-2019-0543` | code | 1.0 |
| 5126 | Microsoft SharePoint Remote Code Execution Vulnerability | kev | `cve-2019-0604` | code | 1.0 |
| 5127 | Microsoft Internet Explorer Information Disclosure Vulnerability | kev | `cve-2019-0676` | code | 1.0 |
| 5128 | Microsoft Windows SMB Information Disclosure Vulnerability | kev | `cve-2019-0703` | code | 1.0 |
| 5129 | Microsoft Remote Desktop Services Remote Code Execution Vulnerability | kev | `cve-2019-0708` | code | 1.0 |
| 5130 | Microsoft Internet Explorer Type Confusion Vulnerability | kev | `cve-2019-0752` | code | 1.0 |
| 5131 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2019-0797` | code | 1.0 |
| 5132 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2019-0803` | code | 1.0 |
| 5133 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2019-0808` | code | 1.0 |
| 5134 | Microsoft Windows AppX Deployment Service (AppXSVC) Privilege Escalation Vulnerability | kev | `cve-2019-0841` | code | 1.0 |
| 5135 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2019-0859` | code | 1.0 |
| 5136 | Microsoft Windows Error Reporting (WER) Privilege Escalation Vulnerability | kev | `cve-2019-0863` | code | 1.0 |
| 5137 | Microsoft Windows Privilege Escalation Vulnerability | kev | `cve-2019-0880` | code | 1.0 |
| 5138 | Microsoft GDI Remote Code Execution Vulnerability | kev | `cve-2019-0903` | code | 1.0 |
| 5139 | Jenkins Script Security Plugin Sandbox Bypass Vulnerability | kev | `cve-2019-1003029` | code | 1.0 |
| 5140 | Jenkins Matrix Project Plugin Remote Code Execution Vulnerability | kev | `cve-2019-1003030` | code | 1.0 |
| 5141 | Kentico Xperience Deserialization of Untrusted Data Vulnerability | kev | `cve-2019-10068` | code | 1.0 |
| 5142 | Exim Mail Transfer Agent (MTA) Improper Input Validation | kev | `cve-2019-10149` | code | 1.0 |
| 5143 | Microsoft Windows AppX Deployment Service (AppXSVC) Privilege Escalation Vulnerability | kev | `cve-2019-1064` | code | 1.0 |
| 5144 | Microsoft SQL Server Remote Code Execution Vulnerability | kev | `cve-2019-1068` | code | 1.0 |
| 5145 | Microsoft Task Scheduler Privilege Escalation Vulnerability | kev | `cve-2019-1069` | code | 1.0 |
| 5146 | MongoDB mongo-express Remote Code Execution Vulnerability | kev | `cve-2019-10758` | code | 1.0 |
| 5147 | Reolink Multiple IP Cameras OS Command Injection Vulnerability | kev | `cve-2019-11001` | code | 1.0 |
| 5148 | PHP FastCGI Process Manager (FPM) Buffer Overflow Vulnerability | kev | `cve-2019-11043` | code | 1.0 |
| 5149 | Microsoft Windows AppX Deployment Service (AppXSVC) Privilege Escalation Vulnerability | kev | `cve-2019-1129` | code | 1.0 |
| 5150 | Microsoft Windows AppX Deployment Service Privilege Escalation Vulnerability | kev | `cve-2019-1130` | code | 1.0 |
| 5151 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2019-1132` | code | 1.0 |
| 5152 | Ivanti Pulse Connect Secure Arbitrary File Read Vulnerability | kev | `cve-2019-11510` | code | 1.0 |
| 5153 | Ivanti Pulse Connect Secure and Policy Secure Command Injection Vulnerability | kev | `cve-2019-11539` | code | 1.0 |
| 5154 | Atlassian Crowd and Crowd Data Center Remote Code Execution Vulnerability | kev | `cve-2019-11580` | code | 1.0 |
| 5155 | Atlassian Jira Server and Data Center Server-Side Template Injection Vulnerability | kev | `cve-2019-11581` | code | 1.0 |
| 5156 | Citrix Workspace Application and Receiver for Windows Remote Code Execution Vulnerability | kev | `cve-2019-11634` | code | 1.0 |
| 5157 | Mozilla Firefox and Thunderbird Type Confusion Vulnerability | kev | `cve-2019-11707` | code | 1.0 |
| 5158 | Mozilla Firefox and Thunderbird Sandbox Escape Vulnerability | kev | `cve-2019-11708` | code | 1.0 |
| 5159 | Microsoft Windows Privilege Common Log File System (CLFS) Escalation Vulnerability | kev | `cve-2019-1214` | code | 1.0 |
| 5160 | Microsoft Windows Privilege Escalation Vulnerability | kev | `cve-2019-1215` | code | 1.0 |
| 5161 | Microsoft Windows AppX Deployment Server Privilege Escalation Vulnerability | kev | `cve-2019-1253` | code | 1.0 |
| 5162 | Microsoft Excel Remote Code Execution Vulnerability | kev | `cve-2019-1297` | code | 1.0 |
| 5163 | Citrix SD-WAN and NetScaler SQL Injection Vulnerability | kev | `cve-2019-12989` | code | 1.0 |
| 5164 | Citrix SD-WAN and NetScaler Command Injection Vulnerability | kev | `cve-2019-12991` | code | 1.0 |
| 5165 | Microsoft Windows Error Reporting Manager Privilege Escalation Vulnerability | kev | `cve-2019-1315` | code | 1.0 |
| 5166 | Microsoft Windows Privilege Escalation Vulnerability | kev | `cve-2019-1322` | code | 1.0 |
| 5167 | Linux Kernel Improper Privilege Management Vulnerability | kev | `cve-2019-13272` | code | 1.0 |
| 5168 | Citrix StoreFront Server XML External Entity (XXE) Processing Vulnerability | kev | `cve-2019-13608` | code | 1.0 |
| 5169 | Microsoft Internet Explorer Scripting Engine Memory Corruption Vulnerability | kev | `cve-2019-1367` | code | 1.0 |
| 5170 | Google Chrome WebAudio Use-After-Free Vulnerability | kev | `cve-2019-13720` | code | 1.0 |
| 5171 | Microsoft Windows AppX Deployment Extensions Privilege Escalation Vulnerability | kev | `cve-2019-1385` | code | 1.0 |
| 5172 | Microsoft Windows Certificate Dialog Privilege Escalation Vulnerability | kev | `cve-2019-1388` | code | 1.0 |
| 5173 | Microsoft Windows Universal Plug and Play (UPnP) Service Privilege Escalation Vulnerability | kev | `cve-2019-1405` | code | 1.0 |
| 5174 | Microsoft Internet Explorer Scripting Engine Memory Corruption Vulnerability | kev | `cve-2019-1429` | code | 1.0 |
| 5175 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2019-1458` | code | 1.0 |
| 5176 | Webmin Command Injection Vulnerability | kev | `cve-2019-15107` | code | 1.0 |
| 5177 | Cisco RV Series Routers Deserialization of Untrusted Data Vulnerability | kev | `cve-2019-15271` | code | 1.0 |
| 5178 | Docker Desktop Community Edition Privilege Escalation Vulnerability | kev | `cve-2019-15752` | code | 1.0 |
| 5179 | Palo Alto Networks PAN-OS Remote Code Execution Vulnerability | kev | `cve-2019-1579` | code | 1.0 |
| 5180 | Nagios XI Remote Code Execution Vulnerability | kev | `cve-2019-15949` | code | 1.0 |
| 5181 | D-Link DNS-320 Remote Code Execution Vulnerability | kev | `cve-2019-16057` | code | 1.0 |
| 5182 | SIMalliance Toolbox Browser Command Injection Vulnerability | kev | `cve-2019-16256` | code | 1.0 |
| 5183 | Nostromo nhttpd Directory Traversal Vulnerability | kev | `cve-2019-16278` | code | 1.0 |
| 5184 | Cisco Small Business Routers Improper Input Validation Vulnerability | kev | `cve-2019-1652` | code | 1.0 |
| 5185 | Cisco Small Business RV320 and RV325 Routers Information Disclosure Vulnerability | kev | `cve-2019-1653` | code | 1.0 |
| 5186 | vBulletin PHP Module Remote Code Execution Vulnerability | kev | `cve-2019-16759` | code | 1.0 |
| 5187 | D-Link Multiple Routers Command Injection Vulnerability | kev | `cve-2019-16920` | code | 1.0 |
| 5188 | Exim Out-of-bounds Write Vulnerability | kev | `cve-2019-16928` | code | 1.0 |
| 5189 | Mozilla Firefox And Thunderbird Type Confusion Vulnerability | kev | `cve-2019-17026` | code | 1.0 |
| 5190 | Apache Solr VelocityResponseWriter Plug-In Remote Code Execution Vulnerability | kev | `cve-2019-17558` | code | 1.0 |
| 5191 | D-Link DIR-859 Router Command Execution Vulnerability | kev | `cve-2019-17621` | code | 1.0 |
| 5192 | Trend Micro OfficeScan Directory Traversal Vulnerability | kev | `cve-2019-18187` | code | 1.0 |
| 5193 | WhatsApp Cross-Site Scripting Vulnerability | kev | `cve-2019-18426` | code | 1.0 |
| 5194 | Progress Telerik UI for ASP.NET AJAX Deserialization of Untrusted Data Vulnerability | kev | `cve-2019-18935` | code | 1.0 |
| 5195 | TeamViewer Desktop Bypass Remote Login Vulnerability | kev | `cve-2019-18988` | code | 1.0 |
| 5196 | Sangoma FreePBX Improper Authentication Vulnerability | kev | `cve-2019-19006` | code | 1.0 |
| 5197 | Netis WF2419 Devices Remote Code Execution Vulnerability | kev | `cve-2019-19356` | code | 1.0 |
| 5198 | Citrix ADC, Gateway, and SD-WAN WANOP Appliance Code Execution Vulnerability | kev | `cve-2019-19781` | code | 1.0 |
| 5199 | TVT NVMS-1000 Directory Traversal Vulnerability | kev | `cve-2019-20085` | code | 1.0 |
| 5200 | D-Link DWL-2600AP Access Point Command Injection Vulnerability | kev | `cve-2019-20500` | code | 1.0 |
| 5201 | Android Kernel Use-After-Free Vulnerability | kev | `cve-2019-2215` | code | 1.0 |
| 5202 | Oracle BI Publisher Unauthorized Access Vulnerability | kev | `cve-2019-2616` | code | 1.0 |
| 5203 | Oracle WebLogic Server, Injection | kev | `cve-2019-2725` | code | 1.0 |
| 5204 | Oracle Solaris Privilege Escalation Vulnerability | kev | `cve-2019-3010` | code | 1.0 |
| 5205 | Atlassian Confluence Server and Data Center Server-Side Template Injection Vulnerability | kev | `cve-2019-3396` | code | 1.0 |
| 5206 | Atlassian Confluence Server and Data Center Path Traversal Vulnerability | kev | `cve-2019-3398` | code | 1.0 |
| 5207 | WhatsApp VOIP Stack Buffer Overflow Vulnerability | kev | `cve-2019-3568` | code | 1.0 |
| 5208 | Crestron Multiple Products Command Injection Vulnerability | kev | `cve-2019-3929` | code | 1.0 |
| 5209 | IBM Planning Analytics Remote Code Execution Vulnerability | kev | `cve-2019-4716` | code | 1.0 |
| 5210 | Rails Ruby on Rails Path Traversal Vulnerability | kev | `cve-2019-5418` | code | 1.0 |
| 5211 | VMware ESXi and Horizon DaaS OpenSLP Heap-Based Buffer Overflow Vulnerability | kev | `cve-2019-5544` | code | 1.0 |
| 5212 | Fortinet FortiOS Default Configuration Vulnerability | kev | `cve-2019-5591` | code | 1.0 |
| 5213 | Google Chrome Blink Use-After-Free Vulnerability | kev | `cve-2019-5786` | code | 1.0 |
| 5214 | Google Chromium V8 Out-of-Bounds Write Vulnerability | kev | `cve-2019-5825` | code | 1.0 |
| 5215 | Apple iOS and macOS Group Facetime Vulnerability | kev | `cve-2019-6223` | code | 1.0 |
| 5216 | Drupal Core Remote Code Execution Vulnerability | kev | `cve-2019-6340` | code | 1.0 |
| 5217 | Fortinet FortiOS Use of Hard-Coded Credentials Vulnerability | kev | `cve-2019-6693` | code | 1.0 |
| 5218 | QNAP Photo Station Improper Access Control Vulnerability | kev | `cve-2019-7192` | code | 1.0 |
| 5219 | QNAP QTS Improper Input Validation Vulnerability | kev | `cve-2019-7193` | code | 1.0 |
| 5220 | QNAP Photo Station Path Traversal Vulnerability | kev | `cve-2019-7194` | code | 1.0 |
| 5221 | QNAP Photo Station Path Traversal Vulnerability | kev | `cve-2019-7195` | code | 1.0 |
| 5222 | Sonatype Nexus Repository Manager Incorrect Access Control Vulnerability | kev | `cve-2019-7238` | code | 1.0 |
| 5223 | Nice Linear eMerge E3-Series OS Command Injection Vulnerability | kev | `cve-2019-7256` | code | 1.0 |
| 5224 | Apple Multiple Products Memory Corruption Vulnerability | kev | `cve-2019-7286` | code | 1.0 |
| 5225 | Apple iOS Memory Corruption Vulnerability | kev | `cve-2019-7287` | code | 1.0 |
| 5226 | SonicWall SMA100 SQL Injection Vulnerability | kev | `cve-2019-7481` | code | 1.0 |
| 5227 | SonicWall SMA100 Directory Traversal Vulnerability | kev | `cve-2019-7483` | code | 1.0 |
| 5228 | Kibana Arbitrary Code Execution | kev | `cve-2019-7609` | code | 1.0 |
| 5229 | Zoho ManageEngine ServiceDesk Plus (SDP) File Upload Vulnerability | kev | `cve-2019-8394` | code | 1.0 |
| 5230 | Apple Multiple Products Type Confusion Vulnerability | kev | `cve-2019-8506` | code | 1.0 |
| 5231 | Apple macOS Use-After-Free Vulnerability | kev | `cve-2019-8526` | code | 1.0 |
| 5232 | Apple Multiple Products Use-After-Free Vulnerability | kev | `cve-2019-8605` | code | 1.0 |
| 5233 | WebKitGTK Memory Corruption Vulnerability | kev | `cve-2019-8720` | code | 1.0 |
| 5234 | ThinkPHP Remote Code Execution Vulnerability | kev | `cve-2019-9082` | code | 1.0 |
| 5235 | Synacor Zimbra Collaboration Suite (ZCS) Server-Side Request Forgery (SSRF) Vulnerability | kev | `cve-2019-9621` | code | 1.0 |
| 5236 | Synacor Zimbra Collaboration Suite (ZCS) Improper Restriction of XML External Entity Reference | kev | `cve-2019-9670` | code | 1.0 |
| 5237 | Sitecore CMS and Experience Platform (XP) Deserialization Vulnerability | kev | `cve-2019-9874` | code | 1.0 |
| 5238 | Sitecore CMS and Experience Platform (XP) Deserialization Vulnerability | kev | `cve-2019-9875` | code | 1.0 |
| 5239 | WordPress Social Warfare Plugin Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2019-9978` | code | 1.0 |
| 5240 | Android Kernel Out-of-Bounds Write Vulnerability | kev | `cve-2020-0041` | code | 1.0 |
| 5241 | Mediatek Multiple Chipsets Insufficient Input Validation Vulnerability | kev | `cve-2020-0069` | code | 1.0 |
| 5242 | Microsoft Windows CryptoAPI Spoofing Vulnerability | kev | `cve-2020-0601` | code | 1.0 |
| 5243 | Microsoft SQL Server Reporting Services Remote Code Execution Vulnerability | kev | `cve-2020-0618` | code | 1.0 |
| 5244 | Microsoft Update Notification Manager Privilege Escalation Vulnerability | kev | `cve-2020-0638` | code | 1.0 |
| 5245 | Microsoft .NET Framework Remote Code Execution Vulnerability | kev | `cve-2020-0646` | code | 1.0 |
| 5246 | Microsoft Internet Explorer Scripting Engine Memory Corruption Vulnerability | kev | `cve-2020-0674` | code | 1.0 |
| 5247 | Microsoft Windows Installer Privilege Escalation Vulnerability | kev | `cve-2020-0683` | code | 1.0 |
| 5248 | Microsoft Exchange Server Validation Key Remote Code Execution Vulnerability | kev | `cve-2020-0688` | code | 1.0 |
| 5249 | Microsoft Windows Background Intelligent Transfer Service (BITS) Improper Privilege Management Vulnerability | kev | `cve-2020-0787` | code | 1.0 |
| 5250 | Microsoft SMBv3 Remote Code Execution Vulnerability | kev | `cve-2020-0796` | code | 1.0 |
| 5251 | Microsoft Edge and Internet Explorer Memory Corruption Vulnerability | kev | `cve-2020-0878` | code | 1.0 |
| 5252 | Microsoft Windows Adobe Font Manager Library Remote Code Execution Vulnerability | kev | `cve-2020-0938` | code | 1.0 |
| 5253 | Microsoft Internet Explorer Scripting Engine Memory Corruption Vulnerability | kev | `cve-2020-0968` | code | 1.0 |
| 5254 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | `cve-2020-0986` | code | 1.0 |
| 5255 | SolarWinds Orion Authentication Bypass Vulnerability | kev | `cve-2020-10148` | code | 1.0 |
| 5256 | Sumavision EMR Cross-Site Request Forgery (CSRF) Vulnerability | kev | `cve-2020-10181` | code | 1.0 |
| 5257 | Zoho ManageEngine Desktop Central File Upload Vulnerability | kev | `cve-2020-10189` | code | 1.0 |
| 5258 | Sonatype Nexus Repository Remote Code Execution Vulnerability | kev | `cve-2020-10199` | code | 1.0 |
| 5259 | Microsoft Windows Adobe Font Manager Library Remote Code Execution Vulnerability | kev | `cve-2020-1020` | code | 1.0 |
| 5260 | rConfig OS Command Injection Vulnerability | kev | `cve-2020-10221` | code | 1.0 |
| 5261 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | `cve-2020-1027` | code | 1.0 |
| 5262 | Microsoft Hyper-V RemoteFX vGPU Remote Code Execution Vulnerability | kev | `cve-2020-1040` | code | 1.0 |
| 5263 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2020-1054` | code | 1.0 |
| 5264 | Tenda AC1900 Router AC15 Model Remote Code Execution Vulnerability | kev | `cve-2020-10987` | code | 1.0 |
| 5265 | JQuery Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2020-11023` | code | 1.0 |
| 5266 | Qualcomm Multiple Chipsets Improper Input Validation Vulnerability | kev | `cve-2020-11261` | code | 1.0 |
| 5267 | Microsoft .NET Framework, SharePoint, and Visual Studio Remote Code Execution Vulnerability | kev | `cve-2020-1147` | code | 1.0 |
| 5268 | SaltStack Salt Authentication Bypass Vulnerability | kev | `cve-2020-11651` | code | 1.0 |
| 5269 | SaltStack Salt Path Traversal Vulnerability | kev | `cve-2020-11652` | code | 1.0 |
| 5270 | WordPress Snap Creek Duplicator Plugin File Download Vulnerability | kev | `cve-2020-11738` | code | 1.0 |
| 5271 | Treck TCP/IP stack Out-of-Bounds Read Vulnerability | kev | `cve-2020-11899` | code | 1.0 |
| 5272 | Apache Airflow Command Injection | kev | `cve-2020-11978` | code | 1.0 |
| 5273 | Sophos SFOS SQL Injection Vulnerability | kev | `cve-2020-12271` | code | 1.0 |
| 5274 | Roundcube Webmail Remote Code Execution Vulnerability | kev | `cve-2020-12641` | code | 1.0 |
| 5275 | Fortinet FortiOS SSL VPN Improper Authentication Vulnerability | kev | `cve-2020-12812` | code | 1.0 |
| 5276 | Microsoft Windows DNS Server Remote Code Execution Vulnerability | kev | `cve-2020-1350` | code | 1.0 |
| 5277 | Drupal core Un-restricted Upload of File | kev | `cve-2020-13671` | code | 1.0 |
| 5278 | Microsoft Internet Explorer Scripting Engine Memory Corruption Vulnerability | kev | `cve-2020-1380` | code | 1.0 |
| 5279 | Apache Airflow's Experimental API Authentication Bypass | kev | `cve-2020-13927` | code | 1.0 |
| 5280 | Roundcube Webmail Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2020-13965` | code | 1.0 |
| 5281 | Microsoft Windows Spoofing Vulnerability | kev | `cve-2020-1464` | code | 1.0 |
| 5282 | Oracle WebLogic Server Remote Code Execution Vulnerability | kev | `cve-2020-14644` | code | 1.0 |
| 5283 | Microsoft Netlogon Privilege Escalation Vulnerability | kev | `cve-2020-1472` | code | 1.0 |
| 5284 | Oracle WebLogic Server Remote Code Execution Vulnerability | kev | `cve-2020-14750` | code | 1.0 |
| 5285 | Oracle Business Intelligence Enterprise Edition Path Transversal | kev | `cve-2020-14864` | code | 1.0 |
| 5286 | Oracle Solaris and Zettabyte File System (ZFS) Unspecified Vulnerability | kev | `cve-2020-14871` | code | 1.0 |
| 5287 | Oracle WebLogic Server Remote Code Execution Vulnerability | kev | `cve-2020-14882` | code | 1.0 |
| 5288 | Oracle WebLogic Server Unspecified Vulnerability | kev | `cve-2020-14883` | code | 1.0 |
| 5289 | Sophos XG Firewall Buffer Overflow Vulnerability | kev | `cve-2020-15069` | code | 1.0 |
| 5290 | DrayTek Multiple Vigor Routers OS Command Injection Vulnerability | kev | `cve-2020-15415` | code | 1.0 |
| 5291 | Ivanti MobileIron Multiple Products Remote Code Execution Vulnerability | kev | `cve-2020-15505` | code | 1.0 |
| 5292 | Google Chrome FreeType Heap Buffer Overflow Vulnerability | kev | `cve-2020-15999` | code | 1.0 |
| 5293 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2020-16009` | code | 1.0 |
| 5294 | Google Chrome for Android UI Heap Buffer Overflow Vulnerability | kev | `cve-2020-16010` | code | 1.0 |
| 5295 | Google Chromium V8 Incorrect Implementation Vulnerabililty | kev | `cve-2020-16013` | code | 1.0 |
| 5296 | Google Chrome Use-After-Free Vulnerability | kev | `cve-2020-16017` | code | 1.0 |
| 5297 | Juniper Junos OS Path Traversal Vulnerability | kev | `cve-2020-1631` | code | 1.0 |
| 5298 | SaltStack Salt Shell Injection Vulnerability | kev | `cve-2020-16846` | code | 1.0 |
| 5299 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | `cve-2020-17087` | code | 1.0 |
| 5300 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | `cve-2020-17144` | code | 1.0 |
| 5301 | Fuel CMS SQL Injection Vulnerability | kev | `cve-2020-17463` | code | 1.0 |
| 5302 | vBulletin PHP Module Remote Code Execution Vulnerability | kev | `cve-2020-17496` | code | 1.0 |
| 5303 | Apache Flink Improper Access Control Vulnerability | kev | `cve-2020-17519` | code | 1.0 |
| 5304 | Apache Struts Remote Code Execution Vulnerability | kev | `cve-2020-17530` | code | 1.0 |
| 5305 | Apache Tomcat Improper Privilege Management Vulnerability | kev | `cve-2020-1938` | code | 1.0 |
| 5306 | Apache Kylin OS Command Injection Vulnerability | kev | `cve-2020-1956` | code | 1.0 |
| 5307 | Palo Alto Networks PAN-OS Authentication Bypass Vulnerability | kev | `cve-2020-2021` | code | 1.0 |
| 5308 | TP-link TL-WA855RE Missing Authentication for Critical Function Vulnerability | kev | `cve-2020-24363` | code | 1.0 |
| 5309 | Trend Micro Multiple Products Improper Access Control Vulnerability | kev | `cve-2020-24557` | code | 1.0 |
| 5310 | QNAP Helpdesk Improper Access Control Vulnerability | kev | `cve-2020-2506` | code | 1.0 |
| 5311 | D-Link DCS-2530L and DCS-2670L Devices Unspecified Vulnerability | kev | `cve-2020-25078` | code | 1.0 |
| 5312 | D-Link DCS-2530L and DCS-2670L Command Injection Vulnerability | kev | `cve-2020-25079` | code | 1.0 |
| 5313 | QNAP Network-Attached Storage (NAS) Command Injection Vulnerability | kev | `cve-2020-2509` | code | 1.0 |
| 5314 | WordPress File Manager Plugin Remote Code Execution Vulnerability | kev | `cve-2020-25213` | code | 1.0 |
| 5315 | Sophos SG UTM Remote Code Execution Vulnerability | kev | `cve-2020-25223` | code | 1.0 |
| 5316 | D-Link DNS-320 Device Command Injection Vulnerability | kev | `cve-2020-25506` | code | 1.0 |
| 5317 | Oracle Fusion Middleware Unspecified Vulnerability | kev | `cve-2020-2551` | code | 1.0 |
| 5318 | Oracle Multiple Products Remote Code Execution Vulnerability | kev | `cve-2020-2555` | code | 1.0 |
| 5319 | Netgear JGS516PE Devices Missing Function Level Access Control Vulnerability | kev | `cve-2020-26919` | code | 1.0 |
| 5320 | Apple Multiple Products Memory Corruption Vulnerability | kev | `cve-2020-27930` | code | 1.0 |
| 5321 | Apple Multiple Products Type Confusion Vulnerability | kev | `cve-2020-27932` | code | 1.0 |
| 5322 | Apple Multiple Products Memory Initialization Vulnerability | kev | `cve-2020-27950` | code | 1.0 |
| 5323 | Oracle WebLogic Server Unspecified Vulnerability | kev | `cve-2020-2883` | code | 1.0 |
| 5324 | PEAR Archive_Tar Deserialization of Untrusted Data Vulnerability | kev | `cve-2020-28949` | code | 1.0 |
| 5325 | D-Link DIR-825 R1 Devices Buffer Overflow Vulnerability | kev | `cve-2020-29557` | code | 1.0 |
| 5326 | CyberoamOS (CROS) SQL Injection Vulnerability | kev | `cve-2020-29574` | code | 1.0 |
| 5327 | Zyxel Multiple Products Use of Hard-Coded Credentials Vulnerability | kev | `cve-2020-29583` | code | 1.0 |
| 5328 | Cisco IOS XR Software Discovery Protocol Format String Vulnerability | kev | `cve-2020-3118` | code | 1.0 |
| 5329 | Cisco AnyConnect Secure Mobility Client for Windows Uncontrolled Search Path Vulnerability | kev | `cve-2020-3153` | code | 1.0 |
| 5330 | Cisco IP Phones Web Server Remote Code Execution and Denial-of-Service Vulnerability | kev | `cve-2020-3161` | code | 1.0 |
| 5331 | Cisco ASA and FTD Information Disclosure Vulnerability | kev | `cve-2020-3259` | code | 1.0 |
| 5332 | Cisco AnyConnect Secure Mobility Client for Windows DLL Hijacking Vulnerability | kev | `cve-2020-3433` | code | 1.0 |
| 5333 | Cisco ASA and FTD Read-Only Path Traversal Vulnerability | kev | `cve-2020-3452` | code | 1.0 |
| 5334 | Cisco IOS XR Software DVMRP Memory Exhaustion Vulnerability | kev | `cve-2020-3566` | code | 1.0 |
| 5335 | Cisco IOS XR Software DVMRP Memory Exhaustion Vulnerability | kev | `cve-2020-3569` | code | 1.0 |
| 5336 | Roundcube Webmail Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2020-35730` | code | 1.0 |
| 5337 | Cisco ASA and FTD Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2020-3580` | code | 1.0 |
| 5338 | PEAR Archive_Tar Improper Link Resolution Vulnerability | kev | `cve-2020-36193` | code | 1.0 |
| 5339 | Apple Multiple Products Memory Corruption Vulnerability | kev | `cve-2020-3837` | code | 1.0 |
| 5340 | VMware Multiple Products Privilege Escalation Vulnerability | kev | `cve-2020-3950` | code | 1.0 |
| 5341 | VMware vCenter Server Information Disclosure Vulnerability | kev | `cve-2020-3952` | code | 1.0 |
| 5342 | VMware ESXi OpenSLP Use-After-Free Vulnerability | kev | `cve-2020-3992` | code | 1.0 |
| 5343 | Multiple VMware Products Command Injection Vulnerability | kev | `cve-2020-4006` | code | 1.0 |
| 5344 | IBM Data Risk Manager Security Bypass Vulnerability | kev | `cve-2020-4427` | code | 1.0 |
| 5345 | IBM Data Risk Manager Remote Code Execution Vulnerability | kev | `cve-2020-4428` | code | 1.0 |
| 5346 | IBM Data Risk Manager Directory Traversal Vulnerability | kev | `cve-2020-4430` | code | 1.0 |
| 5347 | SonicWall SonicOS Buffer Overflow Vulnerability | kev | `cve-2020-5135` | code | 1.0 |
| 5348 | VMware Tanzu Spring Cloud Config Directory Traversal Vulnerability | kev | `cve-2020-5410` | code | 1.0 |
| 5349 | Grandstream Networks UCM6200 Series SQL Injection Vulnerability | kev | `cve-2020-5722` | code | 1.0 |
| 5350 | Amcrest Cameras and NVR Stack-based Buffer Overflow Vulnerability | kev | `cve-2020-5735` | code | 1.0 |
| 5351 | Plex Media Server Remote Code Execution Vulnerability | kev | `cve-2020-5741` | code | 1.0 |
| 5352 | Unraid Remote Code Execution Vulnerability | kev | `cve-2020-5847` | code | 1.0 |
| 5353 | Unraid Authentication Bypass Vulnerability | kev | `cve-2020-5849` | code | 1.0 |
| 5354 | F5 BIG-IP Traffic Management User Interface (TMUI) Remote Code Execution Vulnerability | kev | `cve-2020-5902` | code | 1.0 |
| 5355 | SAP Solution Manager Missing Authentication for Critical Function Vulnerability | kev | `cve-2020-6207` | code | 1.0 |
| 5356 | SAP NetWeaver Missing Authentication for Critical Function Vulnerability | kev | `cve-2020-6287` | code | 1.0 |
| 5357 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2020-6418` | code | 1.0 |
| 5358 | Google Chrome Media Use-After-Free Vulnerability | kev | `cve-2020-6572` | code | 1.0 |
| 5359 | Mozilla Firefox And Thunderbird Use-After-Free Vulnerability | kev | `cve-2020-6819` | code | 1.0 |
| 5360 | Mozilla Firefox And Thunderbird Use-After-Free Vulnerability | kev | `cve-2020-6820` | code | 1.0 |
| 5361 | OpenSMTPD Remote Code Execution Vulnerability | kev | `cve-2020-7247` | code | 1.0 |
| 5362 | Synacor Zimbra Collaboration Suite (ZCS) Server-Side Request Forgery Vulnerability | kev | `cve-2020-7796` | code | 1.0 |
| 5363 | Liferay Portal Deserialization of Untrusted Data Vulnerability | kev | `cve-2020-7961` | code | 1.0 |
| 5364 | Citrix ADC, Gateway, and SD-WAN WANOP Appliance Authorization Bypass Vulnerability | kev | `cve-2020-8193` | code | 1.0 |
| 5365 | Citrix ADC, Gateway, and SD-WAN WANOP Appliance Information Disclosure Vulnerability | kev | `cve-2020-8195` | code | 1.0 |
| 5366 | Citrix ADC, Gateway, and SD-WAN WANOP Appliance Information Disclosure Vulnerability | kev | `cve-2020-8196` | code | 1.0 |
| 5367 | Pulse Connect Secure Code Injection Vulnerability | kev | `cve-2020-8218` | code | 1.0 |
| 5368 | Ivanti Pulse Connect Secure Code Execution Vulnerability | kev | `cve-2020-8243` | code | 1.0 |
| 5369 | Ivanti Pulse Connect Secure Code Execution Vulnerability | kev | `cve-2020-8260` | code | 1.0 |
| 5370 | Trend Micro Apex One and OfficeScan Remote Code Execution Vulnerability | kev | `cve-2020-8467` | code | 1.0 |
| 5371 | Trend Micro Multiple Products Content Validation Escape Vulnerability | kev | `cve-2020-8468` | code | 1.0 |
| 5372 | Multiple DrayTek Vigor Routers Web Management Page Vulnerability | kev | `cve-2020-8515` | code | 1.0 |
| 5373 | Trend Micro Apex One and OfficeScan Authentication Bypass Vulnerability | kev | `cve-2020-8599` | code | 1.0 |
| 5374 | PlaySMS Server-Side Template Injection Vulnerability | kev | `cve-2020-8644` | code | 1.0 |
| 5375 | EyesOfNetwork Improper Privilege Management Vulnerability | kev | `cve-2020-8655` | code | 1.0 |
| 5376 | EyesOfNetwork Use of Hard-Coded Credentials Vulnerability | kev | `cve-2020-8657` | code | 1.0 |
| 5377 | Pi-Hole AdminLTE Remote Code Execution Vulnerability | kev | `cve-2020-8816` | code | 1.0 |
| 5378 | Zyxel Multiple NAS Devices OS Command Injection Vulnerability | kev | `cve-2020-9054` | code | 1.0 |
| 5379 | D-Link DIR-610 Devices Remote Command Execution | kev | `cve-2020-9377` | code | 1.0 |
| 5380 | Adobe Acrobat Use-After-Free Vulnerability | kev | `cve-2020-9715` | code | 1.0 |
| 5381 | Apple iOS, iPadOS, and watchOS Out-of-Bounds Write Vulnerability | kev | `cve-2020-9818` | code | 1.0 |
| 5382 | Apple iOS, iPadOS, and watchOS Memory Corruption Vulnerability | kev | `cve-2020-9819` | code | 1.0 |
| 5383 | Apple Multiple Products Code Execution Vulnerability | kev | `cve-2020-9859` | code | 1.0 |
| 5384 | Apple Multiple Products Memory Corruption Vulnerability | kev | `cve-2020-9907` | code | 1.0 |
| 5385 | Apple iOS, iPadOS, and macOS Input Validation Vulnerability | kev | `cve-2020-9934` | code | 1.0 |
| 5386 | Android Kernel Race Condition Vulnerability | kev | `cve-2021-0920` | code | 1.0 |
| 5387 | Android Kernel Use-After-Free Vulnerability | kev | `cve-2021-1048` | code | 1.0 |
| 5388 | Cisco HyperFlex HX Installer Virtual Machine Command Injection Vulnerability | kev | `cve-2021-1497` | code | 1.0 |
| 5389 | Cisco HyperFlex HX Data Platform Command Injection Vulnerability | kev | `cve-2021-1498` | code | 1.0 |
| 5390 | Microsoft Defender Remote Code Execution Vulnerability | kev | `cve-2021-1647` | code | 1.0 |
| 5391 | Microsoft Windows Print Spooler Remote Code Execution Vulnerability | kev | `cve-2021-1675` | code | 1.0 |
| 5392 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2021-1732` | code | 1.0 |
| 5393 | Apple Multiple Products Race Condition Vulnerability | kev | `cve-2021-1782` | code | 1.0 |
| 5394 | Apple Multiple Products Type Confusion Vulnerability | kev | `cve-2021-1789` | code | 1.0 |
| 5395 | Apple iOS, iPadOS, and macOS WebKit Remote Code Execution Vulnerability | kev | `cve-2021-1870` | code | 1.0 |
| 5396 | Apple iOS, iPadOS, and macOS WebKit Remote Code Execution Vulnerability | kev | `cve-2021-1871` | code | 1.0 |
| 5397 | Apple iOS, iPadOS, and watchOS WebKit Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2021-1879` | code | 1.0 |
| 5398 | Qualcomm Multiple Chipsets Use-After-Free Vulnerability | kev | `cve-2021-1905` | code | 1.0 |
| 5399 | Qualcomm Multiple Chipsets Detection of Error Condition Without Action Vulnerability | kev | `cve-2021-1906` | code | 1.0 |
| 5400 | SonicWall SSLVPN SMA100 SQL Injection Vulnerability | kev | `cve-2021-20016` | code | 1.0 |
| 5401 | SonicWall Email Security Improper Privilege Management Vulnerability | kev | `cve-2021-20021` | code | 1.0 |
| 5402 | SonicWall Email Security Unrestricted Upload of File Vulnerability | kev | `cve-2021-20022` | code | 1.0 |
| 5403 | SonicWall Email Security Path Traversal Vulnerability | kev | `cve-2021-20023` | code | 1.0 |
| 5404 | SonicWall Secure Remote Access (SRA) SQL Injection Vulnerability | kev | `cve-2021-20028` | code | 1.0 |
| 5405 | SonicWall SMA100 Appliances OS Command Injection Vulnerability | kev | `cve-2021-20035` | code | 1.0 |
| 5406 | SonicWall SMA 100 Appliances Stack-Based Buffer Overflow Vulnerability | kev | `cve-2021-20038` | code | 1.0 |
| 5407 | Arcadyan Buffalo Firmware Path Traversal Vulnerability | kev | `cve-2021-20090` | code | 1.0 |
| 5408 | Draytek VigorConnect Path Traversal Vulnerability | kev | `cve-2021-20123` | code | 1.0 |
| 5409 | Draytek VigorConnect Path Traversal Vulnerability | kev | `cve-2021-20124` | code | 1.0 |
| 5410 | Adobe Acrobat and Reader Heap-based Buffer Overflow Vulnerability | kev | `cve-2021-21017` | code | 1.0 |
| 5411 | Google Chromium V8 Heap Buffer Overflow Vulnerability | kev | `cve-2021-21148` | code | 1.0 |
| 5412 | Google Chromium Race Condition Vulnerability | kev | `cve-2021-21166` | code | 1.0 |
| 5413 | Google Chromium Blink Use-After-Free Vulnerability | kev | `cve-2021-21193` | code | 1.0 |
| 5414 | Google Chromium Blink Use-After-Free Vulnerability | kev | `cve-2021-21206` | code | 1.0 |
| 5415 | Google Chromium V8 Improper Input Validation Vulnerability | kev | `cve-2021-21220` | code | 1.0 |
| 5416 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2021-21224` | code | 1.0 |
| 5417 | Adminer Server-Side Request Forgery Vulnerability | kev | `cve-2021-21311` | code | 1.0 |
| 5418 | System Information Library for Node.JS Command Injection | kev | `cve-2021-21315` | code | 1.0 |
| 5419 | Dell dbutil Driver Insufficient Access Control Vulnerability | kev | `cve-2021-21551` | code | 1.0 |
| 5420 | VMware vCenter Server Remote Code Execution Vulnerability | kev | `cve-2021-21972` | code | 1.0 |
| 5421 | VMware vCenter Server and Cloud Foundation Server Side Request Forgery (SSRF) Vulnerability | kev | `cve-2021-21973` | code | 1.0 |
| 5422 | VMware Server Side Request Forgery in vRealize Operations Manager API | kev | `cve-2021-21975` | code | 1.0 |
| 5423 | VMware vCenter Server Improper Input Validation Vulnerability | kev | `cve-2021-21985` | code | 1.0 |
| 5424 | VMware vCenter Server File Upload Vulnerability | kev | `cve-2021-22005` | code | 1.0 |
| 5425 | VMware vCenter Server Improper Access Control | kev | `cve-2021-22017` | code | 1.0 |
| 5426 | Omnissa Workspace ONE Server-Side Request Forgery | kev | `cve-2021-22054` | code | 1.0 |
| 5427 | GitLab Server-Side Request Forgery (SSRF) Vulnerability | kev | `cve-2021-22175` | code | 1.0 |
| 5428 | ExifTool Remote Code Execution Vulnerability | kev | `cve-2021-22204` | code | 1.0 |
| 5429 | GitLab Community and Enterprise Editions Remote Code Execution Vulnerability | kev | `cve-2021-22205` | code | 1.0 |
| 5430 | Micro Focus Operation Bridge Report (OBR) Remote Code Execution Vulnerability | kev | `cve-2021-22502` | code | 1.0 |
| 5431 | Micro Focus Access Manager Information Leakage Vulnerability | kev | `cve-2021-22506` | code | 1.0 |
| 5432 | Linux Kernel Heap Out-of-Bounds Write Vulnerability | kev | `cve-2021-22555` | code | 1.0 |
| 5433 | Linux Kernel Privilege Escalation Vulnerability | kev | `cve-2021-22600` | code | 1.0 |
| 5434 | Rockwell Multiple Products Insufficient Protected Credentials Vulnerability | kev | `cve-2021-22681` | code | 1.0 |
| 5435 | Ivanti Pulse Connect Secure Use-After-Free Vulnerability | kev | `cve-2021-22893` | code | 1.0 |
| 5436 | Ivanti Pulse Connect Secure Collaboration Suite Buffer Overflow Vulnerability | kev | `cve-2021-22894` | code | 1.0 |
| 5437 | Ivanti Pulse Connect Secure Command Injection Vulnerability | kev | `cve-2021-22899` | code | 1.0 |
| 5438 | Ivanti Pulse Connect Secure Unrestricted File Upload Vulnerability | kev | `cve-2021-22900` | code | 1.0 |
| 5439 | Citrix ShareFile Improper Access Control Vulnerability | kev | `cve-2021-22941` | code | 1.0 |
| 5440 | F5 BIG-IP and BIG-IQ Centralized Management iControl REST Remote Code Execution Vulnerability | kev | `cve-2021-22986` | code | 1.0 |
| 5441 | F5 BIG-IP Traffic Management Microkernel Buffer Overflow | kev | `cve-2021-22991` | code | 1.0 |
| 5442 | Ajax.NET Professional Deserialization of Untrusted Data Vulnerability | kev | `cve-2021-23758` | code | 1.0 |
| 5443 | McAfee Total Protection (MTP) Improper Privilege Management Vulnerability | kev | `cve-2021-23874` | code | 1.0 |
| 5444 | Nagios XI OS Command Injection | kev | `cve-2021-25296` | code | 1.0 |
| 5445 | Nagios XI OS Command Injection | kev | `cve-2021-25297` | code | 1.0 |
| 5446 | Nagios XI OS Command Injection | kev | `cve-2021-25298` | code | 1.0 |
| 5447 | Samsung Mobile Devices Improper Access Control Vulnerability | kev | `cve-2021-25337` | code | 1.0 |
| 5448 | Samsung Mobile Devices Improper Access Control Vulnerability | kev | `cve-2021-25369` | code | 1.0 |
| 5449 | Samsung Mobile Devices Memory Corruption Vulnerability | kev | `cve-2021-25370` | code | 1.0 |
| 5450 | Samsung Mobile Devices Unspecified Vulnerability | kev | `cve-2021-25371` | code | 1.0 |
| 5451 | Samsung Mobile Devices Improper Boundary Check Vulnerability | kev | `cve-2021-25372` | code | 1.0 |
| 5452 | Samsung Mobile Devices Race Condition Vulnerability | kev | `cve-2021-25394` | code | 1.0 |
| 5453 | Samsung Mobile Devices Race Condition Vulnerability | kev | `cve-2021-25395` | code | 1.0 |
| 5454 | Samsung Mobile Devices Out-of-Bounds Read Vulnerability | kev | `cve-2021-25487` | code | 1.0 |
| 5455 | Samsung Mobile Devices Improper Input Validation Vulnerability | kev | `cve-2021-25489` | code | 1.0 |
| 5456 | Atlassian Confluence Server and Data Center Object-Graph Navigation Language (OGNL) Injection Vulnerability | kev | `cve-2021-26084` | code | 1.0 |
| 5457 | Atlassian Confluence Server Pre-Authorization Arbitrary File Read Vulnerability | kev | `cve-2021-26085` | code | 1.0 |
| 5458 | Atlassian Jira Server and Data Center Path Traversal Vulnerability | kev | `cve-2021-26086` | code | 1.0 |
| 5459 | Microsoft Internet Explorer Memory Corruption Vulnerability | kev | `cve-2021-26411` | code | 1.0 |
| 5460 | OpenPLC ScadaBR Unrestricted Upload of File with Dangerous Type Vulnerability | kev | `cve-2021-26828` | code | 1.0 |
| 5461 | OpenPLC ScadaBR Cross-site Scripting Vulnerability | kev | `cve-2021-26829` | code | 1.0 |
| 5462 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | `cve-2021-26855` | code | 1.0 |
| 5463 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | `cve-2021-26857` | code | 1.0 |
| 5464 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | `cve-2021-26858` | code | 1.0 |
| 5465 | Microsoft Office Remote Code Execution Vulnerability | kev | `cve-2021-27059` | code | 1.0 |
| 5466 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | `cve-2021-27065` | code | 1.0 |
| 5467 | Microsoft Internet Explorer Remote Code Execution Vulnerability | kev | `cve-2021-27085` | code | 1.0 |
| 5468 | Accellion FTA SQL Injection Vulnerability | kev | `cve-2021-27101` | code | 1.0 |
| 5469 | Accellion FTA OS Command Injection Vulnerability | kev | `cve-2021-27102` | code | 1.0 |
| 5470 | Accellion FTA Server-Side Request Forgery (SSRF) Vulnerability | kev | `cve-2021-27103` | code | 1.0 |
| 5471 | Accellion FTA OS Command Injection Vulnerability | kev | `cve-2021-27104` | code | 1.0 |
| 5472 | DD-WRT Stack-Based Buffer Overflow Vulnerability | kev | `cve-2021-27137` | code | 1.0 |
| 5473 | Yealink Device Management Server-Side Request Forgery (SSRF) Vulnerability | kev | `cve-2021-27561` | code | 1.0 |
| 5474 | Arm Trusted Firmware Out-of-Bounds Write Vulnerability | kev | `cve-2021-27562` | code | 1.0 |
| 5475 | Checkbox Survey Deserialization of Untrusted Data Vulnerability | kev | `cve-2021-27852` | code | 1.0 |
| 5476 | FatPipe WARP, IPVPN, and MPVPN Configuration Upload exploit | kev | `cve-2021-27860` | code | 1.0 |
| 5477 | Veritas Backup Exec Agent File Access Vulnerability | kev | `cve-2021-27876` | code | 1.0 |
| 5478 | Veritas Backup Exec Agent Improper Authentication Vulnerability | kev | `cve-2021-27877` | code | 1.0 |
| 5479 | Veritas Backup Exec Agent Command Execution Vulnerability | kev | `cve-2021-27878` | code | 1.0 |
| 5480 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2021-28310` | code | 1.0 |
| 5481 | Adobe Acrobat and Reader Use-After-Free Vulnerability | kev | `cve-2021-28550` | code | 1.0 |
| 5482 | Arm Mali Graphics Processing Unit (GPU) Use-After-Free Vulnerability | kev | `cve-2021-28663` | code | 1.0 |
| 5483 | Arm Mali Graphics Processing Unit (GPU) Unspecified Vulnerability | kev | `cve-2021-28664` | code | 1.0 |
| 5484 | QNAP NAS Improper Authorization Vulnerability | kev | `cve-2021-28799` | code | 1.0 |
| 5485 | Arm Mali GPU Kernel Driver Use-After-Free Vulnerability | kev | `cve-2021-29256` | code | 1.0 |
| 5486 | Kaseya Virtual System/Server Administrator (VSA) Information Disclosure Vulnerability | kev | `cve-2021-30116` | code | 1.0 |
| 5487 | Google Chromium PopupBlocker Security Bypass Vulnerability | kev | `cve-2021-30533` | code | 1.0 |
| 5488 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2021-30551` | code | 1.0 |
| 5489 | Google Chromium WebGL Use-After-Free Vulnerability | kev | `cve-2021-30554` | code | 1.0 |
| 5490 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2021-30563` | code | 1.0 |
| 5491 | Google Chromium V8 Out-of-Bounds Write Vulnerability | kev | `cve-2021-30632` | code | 1.0 |
| 5492 | Google Chromium Indexed DB API Use-After-Free Vulnerability | kev | `cve-2021-30633` | code | 1.0 |
| 5493 | Apple macOS Unspecified Vulnerability | kev | `cve-2021-30657` | code | 1.0 |
| 5494 | Apple Multiple Products WebKit Storage Use-After-Free Vulnerability | kev | `cve-2021-30661` | code | 1.0 |
| 5495 | Apple Multiple Products WebKit Integer Overflow Vulnerability | kev | `cve-2021-30663` | code | 1.0 |
| 5496 | Apple Multiple Products WebKit Memory Corruption Vulnerability | kev | `cve-2021-30665` | code | 1.0 |
| 5497 | Apple iOS WebKit Buffer Overflow Vulnerability | kev | `cve-2021-30666` | code | 1.0 |
| 5498 | Apple macOS Unspecified Vulnerability | kev | `cve-2021-30713` | code | 1.0 |
| 5499 | Apple iOS WebKit Memory Corruption Vulnerability | kev | `cve-2021-30761` | code | 1.0 |
| 5500 | Apple iOS WebKit Use-After-Free Vulnerability | kev | `cve-2021-30762` | code | 1.0 |
| 5501 | Apple Multiple Products Memory Corruption Vulnerability | kev | `cve-2021-30807` | code | 1.0 |
| 5502 | Apple iOS, iPadOS, macOS Use-After-Free Vulnerability | kev | `cve-2021-30858` | code | 1.0 |
| 5503 | Apple Multiple Products Integer Overflow Vulnerability | kev | `cve-2021-30860` | code | 1.0 |
| 5504 | Apple iOS, iPadOS, and macOS Type Confusion Vulnerability | kev | `cve-2021-30869` | code | 1.0 |
| 5505 | Apple Multiple Products Memory Corruption Vulnerability | kev | `cve-2021-30883` | code | 1.0 |
| 5506 | Apple iOS, iPadOS, and macOS Out-of-Bounds Write Vulnerability | kev | `cve-2021-30900` | code | 1.0 |
| 5507 | Apple Multiple Products Integer Overflow or Wraparound Vulnerability | kev | `cve-2021-30952` | code | 1.0 |
| 5508 | Apple iOS and iPadOS Buffer Overflow Vulnerability | kev | `cve-2021-30983` | code | 1.0 |
| 5509 | Apple iOS, macOS, watchOS Sandbox Bypass Vulnerability | kev | `cve-2021-31010` | code | 1.0 |
| 5510 | Microsoft HTTP Protocol Stack Remote Code Execution Vulnerability | kev | `cve-2021-31166` | code | 1.0 |
| 5511 | Microsoft Exchange Server Information Disclosure Vulnerability | kev | `cve-2021-31196` | code | 1.0 |
| 5512 | Microsoft Enhanced Cryptographic Provider Privilege Escalation Vulnerability | kev | `cve-2021-31199` | code | 1.0 |
| 5513 | Microsoft Enhanced Cryptographic Provider Privilege Escalation Vulnerability | kev | `cve-2021-31201` | code | 1.0 |
| 5514 | Microsoft Exchange Server Security Feature Bypass Vulnerability | kev | `cve-2021-31207` | code | 1.0 |
| 5515 | Laravel Ignition File Upload Vulnerability | kev | `cve-2021-3129` | code | 1.0 |
| 5516 | Sudo Heap-Based Buffer Overflow Vulnerability | kev | `cve-2021-3156` | code | 1.0 |
| 5517 | Tenda AC11 Router Stack Buffer Overflow Vulnerability | kev | `cve-2021-31755` | code | 1.0 |
| 5518 | Microsoft Windows Kernel Information Disclosure Vulnerability | kev | `cve-2021-31955` | code | 1.0 |
| 5519 | Microsoft Windows NTFS Privilege Escalation Vulnerability | kev | `cve-2021-31956` | code | 1.0 |
| 5520 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | `cve-2021-31979` | code | 1.0 |
| 5521 | ONLYOFFICE Docs Server Path Traversal Vulnerability | kev | `cve-2021-3199` | code | 1.0 |
| 5522 | ASUS Routers Improper Authentication Vulnerability | kev | `cve-2021-32030` | code | 1.0 |
| 5523 | October CMS Improper Authentication | kev | `cve-2021-32648` | code | 1.0 |
| 5524 | Dahua IP Camera Authentication Bypass Vulnerability | kev | `cve-2021-33044` | code | 1.0 |
| 5525 | Dahua IP Camera Authentication Bypass Vulnerability | kev | `cve-2021-33045` | code | 1.0 |
| 5526 | Microsoft Desktop Window Manager (DWM) Core Library Privilege Escalation Vulnerability | kev | `cve-2021-33739` | code | 1.0 |
| 5527 | Microsoft Windows MSHTML Platform Remote Code Execution Vulnerability | kev | `cve-2021-33742` | code | 1.0 |
| 5528 | Microsoft Exchange Server Information Disclosure | kev | `cve-2021-33766` | code | 1.0 |
| 5529 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | `cve-2021-33771` | code | 1.0 |
| 5530 | Microsoft Windows Scripting Engine Memory Corruption Vulnerability | kev | `cve-2021-34448` | code | 1.0 |
| 5531 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | `cve-2021-34473` | code | 1.0 |
| 5532 | Microsoft Windows User Profile Service Privilege Escalation Vulnerability | kev | `cve-2021-34484` | code | 1.0 |
| 5533 | Microsoft Windows Event Tracing Privilege Escalation Vulnerability | kev | `cve-2021-34486` | code | 1.0 |
| 5534 | Microsoft Exchange Server Privilege Escalation Vulnerability | kev | `cve-2021-34523` | code | 1.0 |
| 5535 | Microsoft Windows Print Spooler Remote Code Execution Vulnerability | kev | `cve-2021-34527` | code | 1.0 |
| 5536 | Linux Kernel Privilege Escalation Vulnerability | kev | `cve-2021-3493` | code | 1.0 |
| 5537 | SolarWinds Serv-U Remote Code Execution Vulnerability | kev | `cve-2021-35211` | code | 1.0 |
| 5538 | SolarWinds Serv-U Improper Input Validation Vulnerability | kev | `cve-2021-35247` | code | 1.0 |
| 5539 | Realtek Jungle SDK Remote Code Execution Vulnerability | kev | `cve-2021-35394` | code | 1.0 |
| 5540 | Realtek AP-Router SDK Buffer Overflow Vulnerability | kev | `cve-2021-35395` | code | 1.0 |
| 5541 | ForgeRock Access Management (AM) Core Server Remote Code Execution Vulnerability | kev | `cve-2021-35464` | code | 1.0 |
| 5542 | Oracle Fusion Middleware Unspecified Vulnerability | kev | `cve-2021-35587` | code | 1.0 |
| 5543 | Red Hat Polkit Incorrect Authorization Vulnerability | kev | `cve-2021-3560` | code | 1.0 |
| 5544 | Hikvision Improper Input Validation | kev | `cve-2021-36260` | code | 1.0 |
| 5545 | Sunhillo SureLine OS Command Injection Vulnerablity | kev | `cve-2021-36380` | code | 1.0 |
| 5546 | Trend Micro Multiple Products Improper Input Validation Vulnerability | kev | `cve-2021-36741` | code | 1.0 |
| 5547 | Trend Micro Multiple Products Improper Input Validation Vulnerability | kev | `cve-2021-36742` | code | 1.0 |
| 5548 | Microsoft Windows SAM Local Privilege Escalation Vulnerability | kev | `cve-2021-36934` | code | 1.0 |
| 5549 | Microsoft Windows Local Security Authority (LSA) Spoofing Vulnerability | kev | `cve-2021-36942` | code | 1.0 |
| 5550 | Microsoft Windows Update Medic Service Privilege Escalation Vulnerability | kev | `cve-2021-36948` | code | 1.0 |
| 5551 | Microsoft Windows Common Log File System (CLFS) Driver Privilege Escalation Vulnerability | kev | `cve-2021-36955` | code | 1.0 |
| 5552 | Zoho ManageEngine ServiceDesk Authentication Bypass Vulnerability | kev | `cve-2021-37415` | code | 1.0 |
| 5553 | Google Chromium Portals Use-After-Free Vulnerability | kev | `cve-2021-37973` | code | 1.0 |
| 5554 | Google Chromium V8 Use-After-Free Vulnerability | kev | `cve-2021-37975` | code | 1.0 |
| 5555 | Google Chromium Information Disclosure Vulnerability | kev | `cve-2021-37976` | code | 1.0 |
| 5556 | Google Chromium Intents Improper Input Validation Vulnerability | kev | `cve-2021-38000` | code | 1.0 |
| 5557 | Google Chromium V8 Memory Corruption Vulnerability | kev | `cve-2021-38003` | code | 1.0 |
| 5558 | SAP NetWeaver Unrestricted File Upload Vulnerability | kev | `cve-2021-38163` | code | 1.0 |
| 5559 | Delta Electronics DOPSoft 2 Improper Input Validation Vulnerability | kev | `cve-2021-38406` | code | 1.0 |
| 5560 | Microsoft Open Management Infrastructure (OMI) Privilege Escalation Vulnerability | kev | `cve-2021-38645` | code | 1.0 |
| 5561 | Microsoft Office Access Connectivity Engine Remote Code Execution Vulnerability | kev | `cve-2021-38646` | code | 1.0 |
| 5562 | Microsoft Open Management Infrastructure (OMI) Remote Code Execution Vulnerability | kev | `cve-2021-38647` | code | 1.0 |
| 5563 | Microsoft Open Management Infrastructure (OMI) Privilege Escalation Vulnerability | kev | `cve-2021-38648` | code | 1.0 |
| 5564 | Microsoft Open Management Infrastructure (OMI) Privilege Escalation Vulnerability | kev | `cve-2021-38649` | code | 1.0 |
| 5565 | XStream Remote Code Execution Vulnerability | kev | `cve-2021-39144` | code | 1.0 |
| 5566 | Grafana Authentication Bypass Vulnerability | kev | `cve-2021-39226` | code | 1.0 |
| 5567 | Google Pixel Out-of-Bounds Write Vulnerability | kev | `cve-2021-39793` | code | 1.0 |
| 5568 | GitLab Community and Enterprise Editions Server-Side Request Forgery (SSRF) Vulnerability | kev | `cve-2021-39935` | code | 1.0 |
| 5569 | Red Hat Polkit Out-of-Bounds Read and Write Vulnerability | kev | `cve-2021-4034` | code | 1.0 |
| 5570 | Reolink RLC-410W IP Camera OS Command Injection Vulnerability | kev | `cve-2021-40407` | code | 1.0 |
| 5571 | Apache HTTP Server-Side Request Forgery (SSRF) | kev | `cve-2021-40438` | code | 1.0 |
| 5572 | Microsoft MSHTML Remote Code Execution Vulnerability | kev | `cve-2021-40444` | code | 1.0 |
| 5573 | Microsoft Windows Win32k Privilege Escalation Vulnerability | kev | `cve-2021-40449` | code | 1.0 |
| 5574 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2021-40450` | code | 1.0 |
| 5575 | Zoho ManageEngine ADSelfService Plus Authentication Bypass Vulnerability | kev | `cve-2021-40539` | code | 1.0 |
| 5576 | D-Link DIR-605 Router Information Disclosure Vulnerability | kev | `cve-2021-40655` | code | 1.0 |
| 5577 | Aviatrix Controller Unrestricted Upload of File | kev | `cve-2021-40870` | code | 1.0 |
| 5578 | Google Chromium V8 Use-After-Free Vulnerability | kev | `cve-2021-4102` | code | 1.0 |
| 5579 | Metabase GeoJSON API Local File Inclusion Vulnerability | kev | `cve-2021-41277` | code | 1.0 |
| 5580 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2021-41357` | code | 1.0 |
| 5581 | Microsoft Windows Installer Privilege Escalation Vulnerability | kev | `cve-2021-41379` | code | 1.0 |
| 5582 | Apache HTTP Server Path Traversal Vulnerability | kev | `cve-2021-41773` | code | 1.0 |
| 5583 | Apache HTTP Server Path Traversal Vulnerability | kev | `cve-2021-42013` | code | 1.0 |
| 5584 | Sitecore XP Remote Command Execution Vulnerability | kev | `cve-2021-42237` | code | 1.0 |
| 5585 | BQE BillQuick Web Suite SQL Injection Vulnerability | kev | `cve-2021-42258` | code | 1.0 |
| 5586 | Microsoft Active Directory Domain Services Privilege Escalation Vulnerability | kev | `cve-2021-42278` | code | 1.0 |
| 5587 | Microsoft Active Directory Domain Services Privilege Escalation Vulnerability | kev | `cve-2021-42287` | code | 1.0 |
| 5588 | Microsoft Excel Security Feature Bypass | kev | `cve-2021-42292` | code | 1.0 |
| 5589 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | `cve-2021-42321` | code | 1.0 |
| 5590 | Microsoft Windows Privilege Escalation Vulnerability | kev | `cve-2021-43226` | code | 1.0 |
| 5591 | Grafana Path Traversal Vulnerability | kev | `cve-2021-43798` | code | 1.0 |
| 5592 | Microsoft Windows AppX Installer Spoofing Vulnerability | kev | `cve-2021-43890` | code | 1.0 |
| 5593 | Roundcube Webmail SQL Injection Vulnerability | kev | `cve-2021-44026` | code | 1.0 |
| 5594 | Zoho ManageEngine ServiceDesk Plus Remote Code Execution Vulnerability | kev | `cve-2021-44077` | code | 1.0 |
| 5595 | Fortinet FortiOS Arbitrary File Download | kev | `cve-2021-44168` | code | 1.0 |
| 5596 | Acclaim Systems USAHERDS Use of Hard-Coded Credentials Vulnerability | kev | `cve-2021-44207` | code | 1.0 |
| 5597 | Apache Log4j2 Remote Code Execution Vulnerability | kev | `cve-2021-44228` | code | 1.0 |
| 5598 | Zoho Desktop Central Authentication Bypass Vulnerability | kev | `cve-2021-44515` | code | 1.0 |
| 5599 | Ivanti Endpoint Manager Cloud Service Appliance (EPM CSA) Code Injection Vulnerability | kev | `cve-2021-44529` | code | 1.0 |
| 5600 | Apache Log4j2 Deserialization of Untrusted Data Vulnerability | kev | `cve-2021-45046` | code | 1.0 |
| 5601 | D-Link Multiple Routers Remote Code Execution Vulnerability | kev | `cve-2021-45382` | code | 1.0 |
| 5602 | Palo Alto Networks PAN-OS Reflected Amplification Denial-of-Service Vulnerability | kev | `cve-2022-0028` | code | 1.0 |
| 5603 | Linux Kernel Heap-Based Buffer Overflow Vulnerability | kev | `cve-2022-0185` | code | 1.0 |
| 5604 | Linux Kernel Improper Authentication Vulnerability | kev | `cve-2022-0492` | code | 1.0 |
| 5605 | Debian-specific Redis Server Lua Sandbox Escape Vulnerability | kev | `cve-2022-0543` | code | 1.0 |
| 5606 | Google Chromium Animation Use-After-Free Vulnerability | kev | `cve-2022-0609` | code | 1.0 |
| 5607 | Linux Kernel Privilege Escalation Vulnerability | kev | `cve-2022-0847` | code | 1.0 |
| 5608 | Linux Kernel Out-of-Bounds Write Vulnerability | kev | `cve-2022-0995` | code | 1.0 |
| 5609 | Sophos Firewall Authentication Bypass Vulnerability | kev | `cve-2022-1040` | code | 1.0 |
| 5610 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2022-1096` | code | 1.0 |
| 5611 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2022-1364` | code | 1.0 |
| 5612 | F5 BIG-IP Missing Authentication Vulnerability | kev | `cve-2022-1388` | code | 1.0 |
| 5613 | Cisco Small Business RV Series Routers Stack-based Buffer Overflow Vulnerability | kev | `cve-2022-20699` | code | 1.0 |
| 5614 | Cisco Small Business RV Series Routers Stack-based Buffer Overflow Vulnerability | kev | `cve-2022-20700` | code | 1.0 |
| 5615 | Cisco Small Business RV Series Routers Stack-based Buffer Overflow Vulnerability | kev | `cve-2022-20701` | code | 1.0 |
| 5616 | Cisco Small Business RV Series Routers Stack-based Buffer Overflow Vulnerability | kev | `cve-2022-20703` | code | 1.0 |
| 5617 | Cisco Small Business RV Series Routers Stack-based Buffer Overflow Vulnerability | kev | `cve-2022-20708` | code | 1.0 |
| 5618 | Cisco SD-WAN Path Traversal Vulnerability | kev | `cve-2022-20775` | code | 1.0 |
| 5619 | Cisco IOS XR Open Port Vulnerability | kev | `cve-2022-20821` | code | 1.0 |
| 5620 | Oracle ADF Faces Deserialization of Untrusted Data Vulnerability | kev | `cve-2022-21445` | code | 1.0 |
| 5621 | Oracle E-Business Suite Unspecified Vulnerability | kev | `cve-2022-21587` | code | 1.0 |
| 5622 | Microsoft Win32k Privilege Escalation Vulnerability | kev | `cve-2022-21882` | code | 1.0 |
| 5623 | Microsoft Windows User Profile Service Privilege Escalation Vulnerability | kev | `cve-2022-21919` | code | 1.0 |
| 5624 | Microsoft Windows Runtime Remote Code Execution Vulnerability | kev | `cve-2022-21971` | code | 1.0 |
| 5625 | Microsoft Windows Print Spooler Privilege Escalation Vulnerability | kev | `cve-2022-21999` | code | 1.0 |
| 5626 | Microsoft Windows Client Server Runtime Subsystem (CSRSS) Privilege Escalation Vulnerability | kev | `cve-2022-22047` | code | 1.0 |
| 5627 | Qualcomm Multiple Chipsets Use-After-Free Vulnerability | kev | `cve-2022-22071` | code | 1.0 |
| 5628 | Samsung Mobile Devices Use-After-Free Vulnerability | kev | `cve-2022-22265` | code | 1.0 |
| 5629 | SAP Multiple Products HTTP Request Smuggling Vulnerability | kev | `cve-2022-22536` | code | 1.0 |
| 5630 | Apple Memory Corruption Vulnerability | kev | `cve-2022-22587` | code | 1.0 |
| 5631 | Apple iOS, iPadOS, and macOS Webkit Use-After-Free Vulnerability | kev | `cve-2022-22620` | code | 1.0 |
| 5632 | Apple macOS Out-of-Bounds Read Vulnerability | kev | `cve-2022-22674` | code | 1.0 |
| 5633 | Apple macOS Out-of-Bounds Write Vulnerability | kev | `cve-2022-22675` | code | 1.0 |
| 5634 | Arm Mali GPU Kernel Driver Unspecified Vulnerability | kev | `cve-2022-22706` | code | 1.0 |
| 5635 | Microsoft Windows Print Spooler Privilege Escalation Vulnerability | kev | `cve-2022-22718` | code | 1.0 |
| 5636 | WebRTC Heap Buffer Overflow Vulnerability | kev | `cve-2022-2294` | code | 1.0 |
| 5637 | VMware Spring Cloud Gateway Code Injection Vulnerability | kev | `cve-2022-22947` | code | 1.0 |
| 5638 | VMware vCenter Server Incorrect Default File Permissions Vulnerability | kev | `cve-2022-22948` | code | 1.0 |
| 5639 | VMware Workspace ONE Access and Identity Manager Server-Side Template Injection Vulnerability | kev | `cve-2022-22954` | code | 1.0 |
| 5640 | VMware Multiple Products Privilege Escalation Vulnerability | kev | `cve-2022-22960` | code | 1.0 |
| 5641 | VMware Tanzu Spring Cloud Function Remote Code Execution Vulnerability | kev | `cve-2022-22963` | code | 1.0 |
| 5642 | Spring Framework JDK 9+ Remote Code Execution Vulnerability | kev | `cve-2022-22965` | code | 1.0 |
| 5643 | Zabbix Frontend Authentication Bypass Vulnerability | kev | `cve-2022-23131` | code | 1.0 |
| 5644 | Zabbix Frontend Improper Access Control Vulnerability | kev | `cve-2022-23134` | code | 1.0 |
| 5645 | WatchGuard Firebox and XTM Privilege Escalation Vulnerability | kev | `cve-2022-23176` | code | 1.0 |
| 5646 | NUUO NVRmini2 Devices Missing Authentication Vulnerability | kev | `cve-2022-23227` | code | 1.0 |
| 5647 | Dante Discovery Process Control Vulnerability | kev | `cve-2022-23748` | code | 1.0 |
| 5648 | Adobe Commerce and Magento Open Source Improper Input Validation Vulnerability | kev | `cve-2022-24086` | code | 1.0 |
| 5649 | Apache APISIX Authentication Bypass Vulnerability | kev | `cve-2022-24112` | code | 1.0 |
| 5650 | Microsoft Windows CLFS Driver Privilege Escalation Vulnerability | kev | `cve-2022-24521` | code | 1.0 |
| 5651 | Synacor Zimbra Collaborate Suite (ZCS) Cross-Site Scripting Vulnerability | kev | `cve-2022-24682` | code | 1.0 |
| 5652 | Apache CouchDB Insecure Default Initialization of Resource Vulnerability | kev | `cve-2022-24706` | code | 1.0 |
| 5653 | OSGeo GeoServer JAI-EXT Code Injection Vulnerability | kev | `cve-2022-24816` | code | 1.0 |
| 5654 | TerraMaster OS Remote Command Execution Vulnerability | kev | `cve-2022-24990` | code | 1.0 |
| 5655 | Linux Kernel Use-After-Free Vulnerability | kev | `cve-2022-2586` | code | 1.0 |
| 5656 | Atlassian Confluence Server and Data Center Remote Code Execution Vulnerability | kev | `cve-2022-26134` | code | 1.0 |
| 5657 | Atlassian Questions For Confluence App Hard-coded Credentials Vulnerability | kev | `cve-2022-26138` | code | 1.0 |
| 5658 | MiCollab, MiVoice Business Express Access Control Vulnerability | kev | `cve-2022-26143` | code | 1.0 |
| 5659 | D-Link DIR-820L Remote Code Execution Vulnerability | kev | `cve-2022-26258` | code | 1.0 |
| 5660 | WatchGuard Firebox and XTM Appliances Arbitrary Code Execution | kev | `cve-2022-26318` | code | 1.0 |
| 5661 | dotCMS Unrestricted Upload of File Vulnerability | kev | `cve-2022-26352` | code | 1.0 |
| 5662 | Mozilla Firefox Use-After-Free Vulnerability | kev | `cve-2022-26485` | code | 1.0 |
| 5663 | Mozilla Firefox Use-After-Free Vulnerability | kev | `cve-2022-26486` | code | 1.0 |
| 5664 | Veeam Backup & Replication Remote Code Execution Vulnerability | kev | `cve-2022-26500` | code | 1.0 |
| 5665 | Veeam Backup & Replication Remote Code Execution Vulnerability | kev | `cve-2022-26501` | code | 1.0 |
| 5666 | Trend Micro Apex Central Arbitrary File Upload Vulnerability | kev | `cve-2022-26871` | code | 1.0 |
| 5667 | Microsoft Windows User Profile Service Privilege Escalation Vulnerability | kev | `cve-2022-26904` | code | 1.0 |
| 5668 | Microsoft Active Directory Domain Services Privilege Escalation Vulnerability | kev | `cve-2022-26923` | code | 1.0 |
| 5669 | Microsoft Windows LSA Spoofing Vulnerability | kev | `cve-2022-26925` | code | 1.0 |
| 5670 | Citrix Application Delivery Controller (ADC) and Gateway Authentication Bypass Vulnerability | kev | `cve-2022-27518` | code | 1.0 |
| 5671 | QNAP Photo Station Externally Controlled Reference Vulnerability | kev | `cve-2022-27593` | code | 1.0 |
| 5672 | Synacor Zimbra Collaboration Suite (ZCS) Command Injection Vulnerability | kev | `cve-2022-27924` | code | 1.0 |
| 5673 | Synacor Zimbra Collaboration Suite (ZCS) Arbitrary File Upload Vulnerability | kev | `cve-2022-27925` | code | 1.0 |
| 5674 | Synacor Zimbra Collaboration Suite (ZCS) Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2022-27926` | code | 1.0 |
| 5675 | Google Chromium Intents Insufficient Input Validation Vulnerability | kev | `cve-2022-2856` | code | 1.0 |
| 5676 | Zoho ManageEngine ADSelfService Plus Remote Code Execution Vulnerability | kev | `cve-2022-28810` | code | 1.0 |
| 5677 | SolarView Compact Command Injection Vulnerability | kev | `cve-2022-29303` | code | 1.0 |
| 5678 | WSO2 Multiple Products Unrestrictive Upload of File Vulnerability | kev | `cve-2022-29464` | code | 1.0 |
| 5679 | Mitel MiVoice Connect Data Validation Vulnerability | kev | `cve-2022-29499` | code | 1.0 |
| 5680 | Microsoft Windows Support Diagnostic Tool (MSDT) Remote Code Execution Vulnerability | kev | `cve-2022-30190` | code | 1.0 |
| 5681 | RARLAB UnRAR Directory Traversal Vulnerability | kev | `cve-2022-30333` | code | 1.0 |
| 5682 | Google Chromium Network Service Use-After-Free Vulnerability | kev | `cve-2022-3038` | code | 1.0 |
| 5683 | Zyxel Multiple Firewalls OS Command Injection Vulnerability | kev | `cve-2022-30525` | code | 1.0 |
| 5684 | Google Chromium Mojo Insufficient Data Validation Vulnerability | kev | `cve-2022-3075` | code | 1.0 |
| 5685 | Netwrix Auditor Insecure Object Deserialization Vulnerability | kev | `cve-2022-31199` | code | 1.0 |
| 5686 | Sophos Firewall Code Injection Vulnerability | kev | `cve-2022-3236` | code | 1.0 |
| 5687 | Apple iOS and macOS Out-of-Bounds Write Vulnerability | kev | `cve-2022-32893` | code | 1.0 |
| 5688 | Apple iOS and macOS Out-of-Bounds Write Vulnerability | kev | `cve-2022-32894` | code | 1.0 |
| 5689 | Apple iOS, iPadOS, and macOS Remote Code Execution Vulnerability | kev | `cve-2022-32917` | code | 1.0 |
| 5690 | Apache Spark Command Injection Vulnerability | kev | `cve-2022-33891` | code | 1.0 |
| 5691 | Microsoft Windows Support Diagnostic Tool (MSDT) Remote Code Execution Vulnerability | kev | `cve-2022-34713` | code | 1.0 |
| 5692 | Zoho ManageEngine Multiple Products Remote Code Execution Vulnerability | kev | `cve-2022-35405` | code | 1.0 |
| 5693 | Teclib GLPI Remote Code Execution Vulnerability | kev | `cve-2022-35914` | code | 1.0 |
| 5694 | ZK Framework AuUploader Unspecified Vulnerability | kev | `cve-2022-36537` | code | 1.0 |
| 5695 | Atlassian Bitbucket Server and Data Center Command Injection Vulnerability | kev | `cve-2022-36804` | code | 1.0 |
| 5696 | Synacor Zimbra Collaboration Suite (ZCS) Authentication Bypass Vulnerability | kev | `cve-2022-37042` | code | 1.0 |
| 5697 | D-Link Routers Buffer Overflow Vulnerability | kev | `cve-2022-37055` | code | 1.0 |
| 5698 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2022-3723` | code | 1.0 |
| 5699 | Microsoft Windows Common Log File System (CLFS) Driver Privilege Escalation Vulnerability | kev | `cve-2022-37969` | code | 1.0 |
| 5700 | Microsoft Windows Print Spooler Privilege Escalation Vulnerability | kev | `cve-2022-38028` | code | 1.0 |
| 5701 | Arm Mali GPU Kernel Driver Use-After-Free Vulnerability | kev | `cve-2022-38181` | code | 1.0 |
| 5702 | Fortra Cobalt Strike Teamserver Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2022-39197` | code | 1.0 |
| 5703 | Trend Micro Apex One and Apex One as a Service Improper Validation Vulnerability | kev | `cve-2022-40139` | code | 1.0 |
| 5704 | Fortinet Multiple Products Authentication Bypass Vulnerability | kev | `cve-2022-40684` | code | 1.0 |
| 5705 | Mitel MiVoice Connect Command Injection Vulnerability | kev | `cve-2022-40765` | code | 1.0 |
| 5706 | D-Link DNR-322L Download of Code Without Integrity Check Vulnerability | kev | `cve-2022-40799` | code | 1.0 |
| 5707 | Microsoft Windows COM+ Event System Service Privilege Escalation Vulnerability | kev | `cve-2022-41033` | code | 1.0 |
| 5708 | Microsoft Exchange Server Server-Side Request Forgery Vulnerability | kev | `cve-2022-41040` | code | 1.0 |
| 5709 | Microsoft Windows Mark of the Web (MOTW) Security Feature Bypass Vulnerability | kev | `cve-2022-41049` | code | 1.0 |
| 5710 | Microsoft Windows Print Spooler Privilege Escalation Vulnerability | kev | `cve-2022-41073` | code | 1.0 |
| 5711 | Microsoft Exchange Server Privilege Escalation Vulnerability | kev | `cve-2022-41080` | code | 1.0 |
| 5712 | Microsoft Exchange Server Remote Code Execution Vulnerability | kev | `cve-2022-41082` | code | 1.0 |
| 5713 | Microsoft Windows Mark of the Web (MOTW) Security Feature Bypass Vulnerability | kev | `cve-2022-41091` | code | 1.0 |
| 5714 | Microsoft Windows CNG Key Isolation Service Privilege Escalation Vulnerability | kev | `cve-2022-41125` | code | 1.0 |
| 5715 | Microsoft Windows Scripting Languages Remote Code Execution Vulnerability | kev | `cve-2022-41128` | code | 1.0 |
| 5716 | Mitel MiVoice Connect Code Injection Vulnerability | kev | `cve-2022-41223` | code | 1.0 |
| 5717 | Fortinet FortiOS Path Traversal Vulnerability | kev | `cve-2022-41328` | code | 1.0 |
| 5718 | Google Chromium GPU Heap Buffer Overflow Vulnerability | kev | `cve-2022-4135` | code | 1.0 |
| 5719 | Synacor Zimbra Collaboration Suite (ZCS) Arbitrary File Upload Vulnerability | kev | `cve-2022-41352` | code | 1.0 |
| 5720 | Fortinet FortiOS Heap-Based Buffer Overflow Vulnerability | kev | `cve-2022-42475` | code | 1.0 |
| 5721 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2022-4262` | code | 1.0 |
| 5722 | Apple iOS and iPadOS Out-of-Bounds Write Vulnerability | kev | `cve-2022-42827` | code | 1.0 |
| 5723 | Apple iOS Type Confusion Vulnerability | kev | `cve-2022-42856` | code | 1.0 |
| 5724 | Fortra Cobalt Strike User Interface Remote Code Execution Vulnerability | kev | `cve-2022-42948` | code | 1.0 |
| 5725 | Hitachi Vantara Pentaho BA Server Special Element Injection Vulnerability | kev | `cve-2022-43769` | code | 1.0 |
| 5726 | Hitachi Vantara Pentaho BA Server Authorization Bypass Vulnerability | kev | `cve-2022-43939` | code | 1.0 |
| 5727 | Microsoft Defender SmartScreen Security Feature Bypass Vulnerability | kev | `cve-2022-44698` | code | 1.0 |
| 5728 | CWP Control Web Panel OS Command Injection Vulnerability | kev | `cve-2022-44877` | code | 1.0 |
| 5729 | Cacti Command Injection Vulnerability | kev | `cve-2022-46169` | code | 1.0 |
| 5730 | Zoho ManageEngine Multiple Products Remote Code Execution Vulnerability | kev | `cve-2022-47966` | code | 1.0 |
| 5731 | IBM Aspera Faspex Code Execution Vulnerability | kev | `cve-2022-47986` | code | 1.0 |
| 5732 | Apple Multiple Products Unspecified Vulnerability | kev | `cve-2022-48503` | code | 1.0 |
| 5733 | Apple Multiple Products Memory Corruption Vulnerability | kev | `cve-2022-48618` | code | 1.0 |
| 5734 | Linux Kernel Use-After-Free Vulnerability | kev | `cve-2023-0266` | code | 1.0 |
| 5735 | Linux Kernel Improper Ownership Management Vulnerability | kev | `cve-2023-0386` | code | 1.0 |
| 5736 | Fortra GoAnywhere MFT Remote Code Execution Vulnerability | kev | `cve-2023-0669` | code | 1.0 |
| 5737 | TP-Link Archer AX-21 Command Injection Vulnerability | kev | `cve-2023-1389` | code | 1.0 |
| 5738 | Sophos Web Appliance Command Injection Vulnerability | kev | `cve-2023-1671` | code | 1.0 |
| 5739 | Cisco IOS and IOS XE Group Encrypted Transport VPN Out-of-Bounds Write Vulnerability | kev | `cve-2023-20109` | code | 1.0 |
| 5740 | Cisco Small Business RV Series Routers Command Injection Vulnerability | kev | `cve-2023-20118` | code | 1.0 |
| 5741 | Cisco IOS XE Web UI Privilege Escalation Vulnerability | kev | `cve-2023-20198` | code | 1.0 |
| 5742 | Cisco Adaptive Security Appliance and Firepower Threat Defense Unauthorized Access Vulnerability | kev | `cve-2023-20269` | code | 1.0 |
| 5743 | Cisco IOS XE Web UI Command Injection Vulnerability | kev | `cve-2023-20273` | code | 1.0 |
| 5744 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2023-2033` | code | 1.0 |
| 5745 | VMware Tools Authentication Bypass Vulnerability | kev | `cve-2023-20867` | code | 1.0 |
| 5746 | Vmware Aria Operations for Networks Command Injection Vulnerability | kev | `cve-2023-20887` | code | 1.0 |
| 5747 | Android Framework Privilege Escalation Vulnerability | kev | `cve-2023-20963` | code | 1.0 |
| 5748 | Android Pixel Information Disclosure Vulnerability | kev | `cve-2023-21237` | code | 1.0 |
| 5749 | Google Chrome Skia Integer Overflow Vulnerability | kev | `cve-2023-2136` | code | 1.0 |
| 5750 | Samsung Mobile Devices Insertion of Sensitive Information Into Log File Vulnerability | kev | `cve-2023-21492` | code | 1.0 |
| 5751 | Microsoft Exchange Server Deserialization of Untrusted Data Vulnerability | kev | `cve-2023-21529` | code | 1.0 |
| 5752 | Adobe Acrobat and Reader Use-After-Free Vulnerability | kev | `cve-2023-21608` | code | 1.0 |
| 5753 | Microsoft Windows Advanced Local Procedure Call (ALPC) Privilege Escalation Vulnerability | kev | `cve-2023-21674` | code | 1.0 |
| 5754 | Microsoft Office Publisher Security Feature Bypass Vulnerability | kev | `cve-2023-21715` | code | 1.0 |
| 5755 | Microsoft Windows Graphic Component Privilege Escalation Vulnerability | kev | `cve-2023-21823` | code | 1.0 |
| 5756 | Oracle WebLogic Server Unspecified Vulnerability | kev | `cve-2023-21839` | code | 1.0 |
| 5757 | Atlassian Confluence Data Center and Server Broken Access Control Vulnerability | kev | `cve-2023-22515` | code | 1.0 |
| 5758 | Atlassian Confluence Data Center and Server Improper Authorization Vulnerability | kev | `cve-2023-22518` | code | 1.0 |
| 5759 | Atlassian Confluence Data Center and Server Template Injection Vulnerability | kev | `cve-2023-22527` | code | 1.0 |
| 5760 | Strapi Cleartext Storage of Sensitive Information Vulnerability | kev | `cve-2023-22894` | code | 1.0 |
| 5761 | Multiple SugarCRM Products Remote Code Execution Vulnerability | kev | `cve-2023-22952` | code | 1.0 |
| 5762 | Microsoft Windows Common Log File System (CLFS) Driver Privilege Escalation Vulnerability | kev | `cve-2023-23376` | code | 1.0 |
| 5763 | Microsoft Office Outlook Privilege Escalation Vulnerability | kev | `cve-2023-23397` | code | 1.0 |
| 5764 | Apple Multiple Products WebKit Type Confusion Vulnerability | kev | `cve-2023-23529` | code | 1.0 |
| 5765 | Joomla! Improper Access Control Vulnerability | kev | `cve-2023-23752` | code | 1.0 |
| 5766 | Citrix Content Collaboration ShareFile Improper Access Control Vulnerability | kev | `cve-2023-24489` | code | 1.0 |
| 5767 | Microsoft Windows SmartScreen Security Feature Bypass Vulnerability | kev | `cve-2023-24880` | code | 1.0 |
| 5768 | Microsoft SharePoint Server Code Injection Vulnerability | kev | `cve-2023-24955` | code | 1.0 |
| 5769 | D-Link DIR-820 Router OS Command Injection Vulnerability | kev | `cve-2023-25280` | code | 1.0 |
| 5770 | PaperCut NG/MF Cross-Site Request Forgery (CSRF) Vulnerability | kev | `cve-2023-2533` | code | 1.0 |
| 5771 | Multiple Ruckus Wireless Products CSRF and RCE Vulnerability | kev | `cve-2023-25717` | code | 1.0 |
| 5772 | Danfoss AK-SM 800A Series Path Traversal Vulnerability | kev | `cve-2023-25914` | code | 1.0 |
| 5773 | Danfoss AK-SM 800A Series Improper Input Validation Vulnerability | kev | `cve-2023-25915` | code | 1.0 |
| 5774 | Arm Mali GPU Kernel Driver Information Disclosure Vulnerability | kev | `cve-2023-26083` | code | 1.0 |
| 5775 | Adobe ColdFusion Deserialization of Untrusted Data Vulnerability | kev | `cve-2023-26359` | code | 1.0 |
| 5776 | Adobe ColdFusion Deserialization of Untrusted Data Vulnerability | kev | `cve-2023-26360` | code | 1.0 |
| 5777 | Adobe Acrobat and Reader Out-of-Bounds Write Vulnerability | kev | `cve-2023-26369` | code | 1.0 |
| 5778 | PaperCut MF/NG Improper Access Control Vulnerability | kev | `cve-2023-27350` | code | 1.0 |
| 5779 | PaperCut NG/MF Improper Authentication Vulnerability | kev | `cve-2023-27351` | code | 1.0 |
| 5780 | Apache Superset Insecure Default Initialization of Resource Vulnerability | kev | `cve-2023-27524` | code | 1.0 |
| 5781 | Veeam Backup & Replication Cloud Connect Missing Authentication for Critical Function Vulnerability | kev | `cve-2023-27532` | code | 1.0 |
| 5782 | Zyxel Multiple NAS Devices Command Injection Vulnerability | kev | `cve-2023-27992` | code | 1.0 |
| 5783 | Fortinet FortiOS and FortiProxy SSL-VPN Heap-Based Buffer Overflow Vulnerability | kev | `cve-2023-27997` | code | 1.0 |
| 5784 | Apple Multiple Products WebKit Out-of-Bounds Read Vulnerability | kev | `cve-2023-28204` | code | 1.0 |
| 5785 | Apple Multiple Products WebKit Use-After-Free Vulnerability | kev | `cve-2023-28205` | code | 1.0 |
| 5786 | Apple iOS, iPadOS, and macOS IOSurfaceAccelerator Out-of-Bounds Write Vulnerability | kev | `cve-2023-28206` | code | 1.0 |
| 5787 | Microsoft Windows CNG Key Isolation Service Privilege Escalation Vulnerability | kev | `cve-2023-28229` | code | 1.0 |
| 5788 | Microsoft Windows Common Log File System (CLFS) Driver Privilege Escalation Vulnerability | kev | `cve-2023-28252` | code | 1.0 |
| 5789 | MinIO Information Disclosure Vulnerability | kev | `cve-2023-28432` | code | 1.0 |
| 5790 | MinIO Security Feature Bypass Vulnerability | kev | `cve-2023-28434` | code | 1.0 |
| 5791 | Array Networks AG and vxAG ArrayOS Missing Authentication for Critical Function Vulnerability | kev | `cve-2023-28461` | code | 1.0 |
| 5792 | Barracuda Networks ESG Appliance Improper Input Validation Vulnerability | kev | `cve-2023-2868` | code | 1.0 |
| 5793 | Zyxel Multiple Firewalls OS Command Injection Vulnerability | kev | `cve-2023-28771` | code | 1.0 |
| 5794 | Adobe ColdFusion Improper Access Control Vulnerability | kev | `cve-2023-29298` | code | 1.0 |
| 5795 | Adobe ColdFusion Deserialization of Untrusted Data Vulnerability | kev | `cve-2023-29300` | code | 1.0 |
| 5796 | Microsoft Win32K Privilege Escalation Vulnerability | kev | `cve-2023-29336` | code | 1.0 |
| 5797 | Microsoft SharePoint Server Privilege Escalation Vulnerability | kev | `cve-2023-29357` | code | 1.0 |
| 5798 | Microsoft Streaming Service Untrusted Pointer Dereference Vulnerability | kev | `cve-2023-29360` | code | 1.0 |
| 5799 | Novi Survey Insecure Deserialization Vulnerability | kev | `cve-2023-29492` | code | 1.0 |
| 5800 | Service Location Protocol (SLP) Denial-of-Service Vulnerability | kev | `cve-2023-29552` | code | 1.0 |
| 5801 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2023-3079` | code | 1.0 |
| 5802 | Microsoft Windows MSHTML Platform Privilege Escalation Vulnerability | kev | `cve-2023-32046` | code | 1.0 |
| 5803 | Microsoft Windows Defender SmartScreen Security Feature Bypass Vulnerability | kev | `cve-2023-32049` | code | 1.0 |
| 5804 | Ignite Realtime Openfire Path Traversal Vulnerability | kev | `cve-2023-32315` | code | 1.0 |
| 5805 | Apple Multiple Products WebKit Use-After-Free Vulnerability | kev | `cve-2023-32373` | code | 1.0 |
| 5806 | Apple Multiple Products WebKit Sandbox Escape Vulnerability | kev | `cve-2023-32409` | code | 1.0 |
| 5807 | Apple Multiple Products Integer Overflow Vulnerability | kev | `cve-2023-32434` | code | 1.0 |
| 5808 | Apple Multiple Products WebKit Memory Corruption Vulnerability | kev | `cve-2023-32435` | code | 1.0 |
| 5809 | Apple Multiple Products WebKit Type Confusion Vulnerability | kev | `cve-2023-32439` | code | 1.0 |
| 5810 | Zyxel Multiple Firewalls Buffer Overflow Vulnerability | kev | `cve-2023-33009` | code | 1.0 |
| 5811 | Zyxel Multiple Firewalls Buffer Overflow Vulnerability | kev | `cve-2023-33010` | code | 1.0 |
| 5812 | Qualcomm Multiple Chipsets Use-After-Free Vulnerability | kev | `cve-2023-33063` | code | 1.0 |
| 5813 | Qualcomm Multiple Chipsets Use of Out-of-Range Pointer Offset Vulnerability | kev | `cve-2023-33106` | code | 1.0 |
| 5814 | Qualcomm Multiple Chipsets Integer Overflow Vulnerability | kev | `cve-2023-33107` | code | 1.0 |
| 5815 | Apache RocketMQ Command Execution Vulnerability | kev | `cve-2023-33246` | code | 1.0 |
| 5816 | TP-Link Multiple Routers Command Injection Vulnerability | kev | `cve-2023-33538` | code | 1.0 |
| 5817 | VMware vCenter Server Out-of-Bounds Write Vulnerability | kev | `cve-2023-34048` | code | 1.0 |
| 5818 | Synacor Zimbra Collaboration Suite (ZCS) Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2023-34192` | code | 1.0 |
| 5819 | Progress MOVEit Transfer SQL Injection Vulnerability | kev | `cve-2023-34362` | code | 1.0 |
| 5820 | Ivanti Endpoint Manager Mobile Authentication Bypass Vulnerability | kev | `cve-2023-35078` | code | 1.0 |
| 5821 | Ivanti Endpoint Manager Mobile (EPMM) Path Traversal Vulnerability | kev | `cve-2023-35081` | code | 1.0 |
| 5822 | Ivanti Endpoint Manager Mobile (EPMM) and MobileIron Core Authentication Bypass Vulnerability | kev | `cve-2023-35082` | code | 1.0 |
| 5823 | Citrix NetScaler ADC and NetScaler Gateway Code Injection Vulnerability | kev | `cve-2023-3519` | code | 1.0 |
| 5824 | Microsoft Outlook Security Feature Bypass Vulnerability | kev | `cve-2023-35311` | code | 1.0 |
| 5825 | Android Framework Privilege Escalation Vulnerability | kev | `cve-2023-35674` | code | 1.0 |
| 5826 | Microsoft Windows SmartScreen Security Feature Bypass Vulnerability | kev | `cve-2023-36025` | code | 1.0 |
| 5827 | Microsoft Windows Desktop Window Manager (DWM) Core Library Privilege Escalation Vulnerability | kev | `cve-2023-36033` | code | 1.0 |
| 5828 | Microsoft Windows Cloud Files Mini Filter Driver Privilege Escalation Vulnerability | kev | `cve-2023-36036` | code | 1.0 |
| 5829 | Microsoft Windows Out-of-Bounds Read Vulnerability | kev | `cve-2023-36424` | code | 1.0 |
| 5830 | Microsoft WordPad Information Disclosure Vulnerability | kev | `cve-2023-36563` | code | 1.0 |
| 5831 | Microsoft Windows Mark of the Web (MOTW) Security Feature Bypass Vulnerability | kev | `cve-2023-36584` | code | 1.0 |
| 5832 | Microsoft Word Information Disclosure Vulnerability | kev | `cve-2023-36761` | code | 1.0 |
| 5833 | Microsoft Streaming Service Proxy Privilege Escalation Vulnerability | kev | `cve-2023-36802` | code | 1.0 |
| 5834 | Juniper Junos OS EX Series PHP External Variable Modification Vulnerability | kev | `cve-2023-36844` | code | 1.0 |
| 5835 | Juniper Junos OS EX Series and SRX Series PHP External Variable Modification Vulnerability | kev | `cve-2023-36845` | code | 1.0 |
| 5836 | Juniper Junos OS SRX Series Missing Authentication for Critical Function Vulnerability | kev | `cve-2023-36846` | code | 1.0 |
| 5837 | Juniper Junos OS EX Series Missing Authentication for Critical Function Vulnerability | kev | `cve-2023-36847` | code | 1.0 |
| 5838 | Juniper Junos OS SRX Series Missing Authentication for Critical Function Vulnerability | kev | `cve-2023-36851` | code | 1.0 |
| 5839 | Microsoft Windows Error Reporting Service Privilege Escalation Vulnerability | kev | `cve-2023-36874` | code | 1.0 |
| 5840 | Microsoft Windows Search Remote Code Execution Vulnerability | kev | `cve-2023-36884` | code | 1.0 |
| 5841 | Apple Multiple Products WebKit Code Execution Vulnerability | kev | `cve-2023-37450` | code | 1.0 |
| 5842 | Synacor Zimbra Collaboration Suite (ZCS) Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2023-37580` | code | 1.0 |
| 5843 | Ivanti Sentry Authentication Bypass Vulnerability | kev | `cve-2023-38035` | code | 1.0 |
| 5844 | Microsoft .NET Core and Visual Studio Denial-of-Service Vulnerability | kev | `cve-2023-38180` | code | 1.0 |
| 5845 | Adobe ColdFusion Deserialization of Untrusted Data Vulnerability | kev | `cve-2023-38203` | code | 1.0 |
| 5846 | Adobe ColdFusion Improper Access Control Vulnerability | kev | `cve-2023-38205` | code | 1.0 |
| 5847 | Apple Multiple Products Kernel Unspecified Vulnerability | kev | `cve-2023-38606` | code | 1.0 |
| 5848 | RARLAB WinRAR Code Execution Vulnerability | kev | `cve-2023-38831` | code | 1.0 |
| 5849 | ZKTeco BioTime Path Traversal Vulnerability | kev | `cve-2023-38950` | code | 1.0 |
| 5850 | ASUS RT-AX55 Routers OS Command Injection Vulnerability | kev | `cve-2023-39780` | code | 1.0 |
| 5851 | Progress WS_FTP Server Deserialization of Untrusted Data Vulnerability | kev | `cve-2023-40044` | code | 1.0 |
| 5852 | Apple iOS, iPadOS, and watchOS Wallet Code Execution Vulnerability | kev | `cve-2023-41061` | code | 1.0 |
| 5853 | Apple iOS, iPadOS, and macOS ImageIO Buffer Overflow Vulnerability | kev | `cve-2023-41064` | code | 1.0 |
| 5854 | Trend Micro Apex One and Worry-Free Business Security Remote Code Execution Vulnerability | kev | `cve-2023-41179` | code | 1.0 |
| 5855 | Qlik Sense HTTP Tunneling Vulnerability | kev | `cve-2023-41265` | code | 1.0 |
| 5856 | Qlik Sense Path Traversal Vulnerability | kev | `cve-2023-41266` | code | 1.0 |
| 5857 | Microsoft Skype for Business Privilege Escalation Vulnerability | kev | `cve-2023-41763` | code | 1.0 |
| 5858 | Apple iOS and iPadOS Use-After-Free Vulnerability | kev | `cve-2023-41974` | code | 1.0 |
| 5859 | Apple Multiple Products Code Execution Vulnerability | kev | `cve-2023-41990` | code | 1.0 |
| 5860 | Apple Multiple Products Improper Certificate Validation Vulnerability | kev | `cve-2023-41991` | code | 1.0 |
| 5861 | Apple Multiple Products Kernel Privilege Escalation Vulnerability | kev | `cve-2023-41992` | code | 1.0 |
| 5862 | Apple Multiple Products WebKit Code Execution Vulnerability | kev | `cve-2023-41993` | code | 1.0 |
| 5863 | Arm Mali GPU Kernel Driver Use-After-Free Vulnerability | kev | `cve-2023-4211` | code | 1.0 |
| 5864 | JetBrains TeamCity Authentication Bypass Vulnerability | kev | `cve-2023-42793` | code | 1.0 |
| 5865 | Apple iOS and iPadOS Kernel Privilege Escalation Vulnerability | kev | `cve-2023-42824` | code | 1.0 |
| 5866 | Apple Multiple Products WebKit Out-of-Bounds Read Vulnerability | kev | `cve-2023-42916` | code | 1.0 |
| 5867 | Apple Multiple Products WebKit Memory Corruption Vulnerability | kev | `cve-2023-42917` | code | 1.0 |
| 5868 | Apple Multiple products Use-After-Free Vulnerability | kev | `cve-2023-43000` | code | 1.0 |
| 5869 | NextGen Healthcare Mirth Connect Deserialization of Untrusted Data Vulnerability | kev | `cve-2023-43208` | code | 1.0 |
| 5870 | KNX Association KNX Protocol Connection Authorization Option 1 Overly Restrictive Account Lockout Mechanism Vulnerability | kev | `cve-2023-4346` | code | 1.0 |
| 5871 | Roundcube Webmail Persistent Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2023-43770` | code | 1.0 |
| 5872 | SonicWall SMA100 Appliances OS Command Injection Vulnerability | kev | `cve-2023-44221` | code | 1.0 |
| 5873 | HTTP/2 Rapid Reset Attack Vulnerability | kev | `cve-2023-44487` | code | 1.0 |
| 5874 | Acronis Cyber Infrastructure (ACI) Insecure Default Password Vulnerability | kev | `cve-2023-45249` | code | 1.0 |
| 5875 | North Grid Proself Improper Restriction of XML External Entity (XXE) Reference Vulnerability | kev | `cve-2023-45727` | code | 1.0 |
| 5876 | Apache ActiveMQ Deserialization of Untrusted Data Vulnerability | kev | `cve-2023-46604` | code | 1.0 |
| 5877 | F5 BIG-IP Configuration Utility Authentication Bypass Vulnerability | kev | `cve-2023-46747` | code | 1.0 |
| 5878 | F5 BIG-IP Configuration Utility SQL Injection Vulnerability | kev | `cve-2023-46748` | code | 1.0 |
| 5879 | Ivanti Connect Secure and Policy Secure Authentication Bypass Vulnerability | kev | `cve-2023-46805` | code | 1.0 |
| 5880 | SysAid Server Path Traversal Vulnerability | kev | `cve-2023-47246` | code | 1.0 |
| 5881 | QNAP VioStor NVR OS Command Injection Vulnerability | kev | `cve-2023-47565` | code | 1.0 |
| 5882 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2023-4762` | code | 1.0 |
| 5883 | Qlik Sense HTTP Tunneling Vulnerability | kev | `cve-2023-48365` | code | 1.0 |
| 5884 | Google Chromium WebP Heap-Based Buffer Overflow Vulnerability | kev | `cve-2023-4863` | code | 1.0 |
| 5885 | Fortinet FortiClient EMS SQL Injection Vulnerability | kev | `cve-2023-48788` | code | 1.0 |
| 5886 | ownCloud graphapi Information Disclosure Vulnerability | kev | `cve-2023-49103` | code | 1.0 |
| 5887 | ownCloud Improper Authentication Vulnerability | kev | `cve-2023-49105` | code | 1.0 |
| 5888 | GNU C Library Buffer Overflow Vulnerability | kev | `cve-2023-4911` | code | 1.0 |
| 5889 | Citrix NetScaler ADC and NetScaler Gateway Buffer Overflow Vulnerability | kev | `cve-2023-4966` | code | 1.0 |
| 5890 | FXC AE1021, AE1021PE OS Command Injection Vulnerability | kev | `cve-2023-49897` | code | 1.0 |
| 5891 | TP-Link TL-WR841N Authentication Bypass by Spoofing Vulnerability | kev | `cve-2023-50224` | code | 1.0 |
| 5892 | Digiever DS-2105 Pro Missing Authorization Vulnerability | kev | `cve-2023-52163` | code | 1.0 |
| 5893 | Google Chromium libvpx Heap Buffer Overflow Vulnerability | kev | `cve-2023-5217` | code | 1.0 |
| 5894 | Roundcube Webmail Persistent Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2023-5631` | code | 1.0 |
| 5895 | Google Skia Integer Overflow Vulnerability | kev | `cve-2023-6345` | code | 1.0 |
| 5896 | Unitronics Vision PLC and HMI Insecure Default Password Vulnerability | kev | `cve-2023-6448` | code | 1.0 |
| 5897 | Citrix NetScaler ADC and NetScaler Gateway Code Injection Vulnerability | kev | `cve-2023-6548` | code | 1.0 |
| 5898 | Citrix NetScaler ADC and NetScaler Gateway Buffer Overflow Vulnerability | kev | `cve-2023-6549` | code | 1.0 |
| 5899 | Google Chromium WebRTC Heap Buffer Overflow Vulnerability | kev | `cve-2023-7024` | code | 1.0 |
| 5900 | GitLab Community and Enterprise Editions Improper Access Control Vulnerability | kev | `cve-2023-7028` | code | 1.0 |
| 5901 | Spreadsheet::ParseExcel Remote Code Execution Vulnerability | kev | `cve-2023-7101` | code | 1.0 |
| 5902 | Palo Alto Networks PAN-OS Management Interface Authentication Bypass Vulnerability | kev | `cve-2024-0012` | code | 1.0 |
| 5903 | Google Chromium V8 Out-of-Bounds Memory Access Vulnerability | kev | `cve-2024-0519` | code | 1.0 |
| 5904 | D-Link DIR-859 Router Path Traversal Vulnerability | kev | `cve-2024-0769` | code | 1.0 |
| 5905 | Synology DiskStation Manager (DSM) Improper Encoding or Escaping of Output Vulnerability | kev | `cve-2024-10441` | code | 1.0 |
| 5906 | Linux Kernel Use-After-Free Vulnerability | kev | `cve-2024-1086` | code | 1.0 |
| 5907 | GeoVision Devices OS Command Injection Vulnerability | kev | `cve-2024-11120` | code | 1.0 |
| 5908 | MDaemon Email Server Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2024-11182` | code | 1.0 |
| 5909 | Zyxel Multiple Firewalls Path Traversal Vulnerability | kev | `cve-2024-11667` | code | 1.0 |
| 5910 | ProjectSend Improper Authentication Vulnerability | kev | `cve-2024-11680` | code | 1.0 |
| 5911 | Progress Kemp LoadMaster OS Command Injection Vulnerability | kev | `cve-2024-1212` | code | 1.0 |
| 5912 | BeyondTrust Privileged Remote Access (PRA) and Remote Support (RS) Command Injection Vulnerability | kev | `cve-2024-12356` | code | 1.0 |
| 5913 | BeyondTrust Privileged Remote Access (PRA) and Remote Support (RS) OS Command Injection Vulnerability | kev | `cve-2024-12686` | code | 1.0 |
| 5914 | DrayTek Vigor Routers OS Command Injection Vulnerability | kev | `cve-2024-12987` | code | 1.0 |
| 5915 | Ivanti Endpoint Manager (EPM) Absolute Path Traversal Vulnerability | kev | `cve-2024-13159` | code | 1.0 |
| 5916 | Ivanti Endpoint Manager (EPM) Absolute Path Traversal Vulnerability | kev | `cve-2024-13160` | code | 1.0 |
| 5917 | Ivanti Endpoint Manager (EPM) Absolute Path Traversal Vulnerability | kev | `cve-2024-13161` | code | 1.0 |
| 5918 | ConnectWise ScreenConnect Path Traversal Vulnerability | kev | `cve-2024-1708` | code | 1.0 |
| 5919 | ConnectWise ScreenConnect Authentication Bypass Vulnerability | kev | `cve-2024-1709` | code | 1.0 |
| 5920 | Cisco ASA and FTD Denial of Service Vulnerability | kev | `cve-2024-20353` | code | 1.0 |
| 5921 | Cisco ASA and FTD Privilege Escalation Vulnerability | kev | `cve-2024-20359` | code | 1.0 |
| 5922 | Cisco NX-OS Command Injection Vulnerability | kev | `cve-2024-20399` | code | 1.0 |
| 5923 | Cisco Smart Licensing Utility Static Credential Vulnerability | kev | `cve-2024-20439` | code | 1.0 |
| 5924 | Cisco ASA and FTD Denial-of-Service Vulnerability | kev | `cve-2024-20481` | code | 1.0 |
| 5925 | Adobe ColdFusion Improper Access Control Vulnerability | kev | `cve-2024-20767` | code | 1.0 |
| 5926 | Oracle Agile Product Lifecycle Management (PLM) Deserialization Vulnerability | kev | `cve-2024-20953` | code | 1.0 |
| 5927 | Oracle WebLogic Server Unspecified Vulnerability | kev | `cve-2024-21182` | code | 1.0 |
| 5928 | Oracle Agile Product Lifecycle Management (PLM) Incorrect Authorization Vulnerability | kev | `cve-2024-21287` | code | 1.0 |
| 5929 | Microsoft Windows Kernel Exposed IOCTL with Insufficient Access Control Vulnerability | kev | `cve-2024-21338` | code | 1.0 |
| 5930 | Microsoft Windows SmartScreen Security Feature Bypass Vulnerability | kev | `cve-2024-21351` | code | 1.0 |
| 5931 | Microsoft Exchange Server Privilege Escalation Vulnerability | kev | `cve-2024-21410` | code | 1.0 |
| 5932 | Microsoft Windows Internet Shortcut Files Security Feature Bypass Vulnerability | kev | `cve-2024-21412` | code | 1.0 |
| 5933 | Microsoft Outlook Improper Input Validation Vulnerability | kev | `cve-2024-21413` | code | 1.0 |
| 5934 | Fortinet FortiOS Out-of-Bound Write Vulnerability | kev | `cve-2024-21762` | code | 1.0 |
| 5935 | Ivanti Connect Secure and Policy Secure Command Injection Vulnerability | kev | `cve-2024-21887` | code | 1.0 |
| 5936 | Ivanti Connect Secure, Policy Secure, and Neurons Server-Side Request Forgery (SSRF) Vulnerability | kev | `cve-2024-21893` | code | 1.0 |
| 5937 | Fortinet Multiple Products Format String Vulnerability | kev | `cve-2024-23113` | code | 1.0 |
| 5938 | Apple Multiple Products WebKit Type Confusion Vulnerability | kev | `cve-2024-23222` | code | 1.0 |
| 5939 | Apple Multiple Products Memory Corruption Vulnerability | kev | `cve-2024-23225` | code | 1.0 |
| 5940 | Apple Multiple Products Memory Corruption Vulnerability | kev | `cve-2024-23296` | code | 1.0 |
| 5941 | Rejetto HTTP File Server Improper Neutralization of Special Elements Used in a Template Engine Vulnerability | kev | `cve-2024-23692` | code | 1.0 |
| 5942 | Jenkins Command Line Interface (CLI) Path Traversal Vulnerability | kev | `cve-2024-23897` | code | 1.0 |
| 5943 | Check Point Quantum Security Gateways Information Disclosure Vulnerability | kev | `cve-2024-24919` | code | 1.0 |
| 5944 | Microsoft Windows Error Reporting Service Improper Privilege Management Vulnerability | kev | `cve-2024-26169` | code | 1.0 |
| 5945 | JetBrains TeamCity Authentication Bypass Vulnerability | kev | `cve-2024-27198` | code | 1.0 |
| 5946 | JetBrains TeamCity Relative Path Traversal Vulnerability | kev | `cve-2024-27199` | code | 1.0 |
| 5947 | Apache HugeGraph-Server Improper Access Control Vulnerability | kev | `cve-2024-27348` | code | 1.0 |
| 5948 | Synacor Zimbra Collaboration Suite (ZCS) Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2024-27443` | code | 1.0 |
| 5949 | Apple Multiple Products Improper Authentication Vulnerability | kev | `cve-2024-27834` | code | 1.0 |
| 5950 | SolarWinds Web Help Desk Deserialization of Untrusted Data Vulnerability | kev | `cve-2024-28986` | code | 1.0 |
| 5951 | SolarWinds Web Help Desk Hardcoded Credential Vulnerability | kev | `cve-2024-28987` | code | 1.0 |
| 5952 | SolarWinds Serv-U Path Traversal Vulnerability | kev | `cve-2024-28995` | code | 1.0 |
| 5953 | Microsoft .NET Framework Information Disclosure Vulnerability | kev | `cve-2024-29059` | code | 1.0 |
| 5954 | Android Pixel Information Disclosure Vulnerability | kev | `cve-2024-29745` | code | 1.0 |
| 5955 | Android Pixel Privilege Escalation Vulnerability | kev | `cve-2024-29748` | code | 1.0 |
| 5956 | Ivanti Endpoint Manager (EPM) SQL Injection Vulnerability | kev | `cve-2024-29824` | code | 1.0 |
| 5957 | Microsoft SmartScreen Prompt Security Feature Bypass Vulnerability | kev | `cve-2024-29988` | code | 1.0 |
| 5958 | Microsoft Windows MSHTML Platform Security Feature Bypass Vulnerability | kev | `cve-2024-30040` | code | 1.0 |
| 5959 | Microsoft DWM Core Library Privilege Escalation Vulnerability | kev | `cve-2024-30051` | code | 1.0 |
| 5960 | Microsoft Windows Kernel TOCTOU Race Condition Vulnerability | kev | `cve-2024-30088` | code | 1.0 |
| 5961 | Apache OFBiz Path Traversal Vulnerability | kev | `cve-2024-32113` | code | 1.0 |
| 5962 | D-Link Multiple NAS Devices Use of Hard-Coded Credentials Vulnerability | kev | `cve-2024-3272` | code | 1.0 |
| 5963 | D-Link Multiple NAS Devices Command Injection Vulnerability | kev | `cve-2024-3273` | code | 1.0 |
| 5964 | Android Pixel Privilege Escalation Vulnerability | kev | `cve-2024-32896` | code | 1.0 |
| 5965 | Palo Alto Networks PAN-OS Malicious DNS Packet Vulnerability | kev | `cve-2024-3393` | code | 1.0 |
| 5966 | Palo Alto Networks PAN-OS Command Injection Vulnerability | kev | `cve-2024-3400` | code | 1.0 |
| 5967 | Adobe Commerce and Magento Open Source Improper Restriction of XML External Entity Reference (XXE) Vulnerability | kev | `cve-2024-34102` | code | 1.0 |
| 5968 | Microsoft Windows Kernel-Mode Driver Untrusted Pointer Dereference Vulnerability | kev | `cve-2024-35250` | code | 1.0 |
| 5969 | OSGeo GeoServer GeoTools Eval Injection Vulnerability | kev | `cve-2024-36401` | code | 1.0 |
| 5970 | Android Kernel Remote Code Execution Vulnerability | kev | `cve-2024-36971` | code | 1.0 |
| 5971 | Broadcom VMware vCenter Server Out-of-bounds Write Vulnerability | kev | `cve-2024-37079` | code | 1.0 |
| 5972 | VMware ESXi Authentication Bypass Vulnerability | kev | `cve-2024-37085` | code | 1.0 |
| 5973 | RoundCube Webmail Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2024-37383` | code | 1.0 |
| 5974 | Microsoft Windows Installer Improper Privilege Management Vulnerability | kev | `cve-2024-38014` | code | 1.0 |
| 5975 | Microsoft Windows Hyper-V Privilege Escalation Vulnerability | kev | `cve-2024-38080` | code | 1.0 |
| 5976 | Microsoft SharePoint Deserialization Vulnerability | kev | `cve-2024-38094` | code | 1.0 |
| 5977 | Microsoft Windows Kernel Privilege Escalation Vulnerability | kev | `cve-2024-38106` | code | 1.0 |
| 5978 | Microsoft Windows Power Dependency Coordinator Privilege Escalation Vulnerability | kev | `cve-2024-38107` | code | 1.0 |
| 5979 | Microsoft Windows MSHTML Platform Spoofing Vulnerability | kev | `cve-2024-38112` | code | 1.0 |
| 5980 | Microsoft Windows Scripting Engine Memory Corruption Vulnerability | kev | `cve-2024-38178` | code | 1.0 |
| 5981 | Microsoft Project Remote Code Execution Vulnerability | kev | `cve-2024-38189` | code | 1.0 |
| 5982 | Microsoft Windows Ancillary Function Driver for WinSock Privilege Escalation Vulnerability | kev | `cve-2024-38193` | code | 1.0 |
| 5983 | Microsoft Windows SmartScreen Security Feature Bypass Vulnerability | kev | `cve-2024-38213` | code | 1.0 |
| 5984 | Microsoft Windows Mark of the Web (MOTW) Protection Mechanism Failure Vulnerability | kev | `cve-2024-38217` | code | 1.0 |
| 5985 | Microsoft Publisher Protection Mechanism Failure Vulnerability | kev | `cve-2024-38226` | code | 1.0 |
| 5986 | Apache HTTP Server Improper Escaping of Output Vulnerability | kev | `cve-2024-38475` | code | 1.0 |
| 5987 | VMware vCenter Server Heap-Based Buffer Overflow Vulnerability | kev | `cve-2024-38812` | code | 1.0 |
| 5988 | VMware vCenter Server Privilege Escalation Vulnerability | kev | `cve-2024-38813` | code | 1.0 |
| 5989 | Apache OFBiz Incorrect Authorization Vulnerability | kev | `cve-2024-38856` | code | 1.0 |
| 5990 | Versa Director Dangerous File Type Upload Vulnerability | kev | `cve-2024-39717` | code | 1.0 |
| 5991 | Twilio Authy Information Disclosure Vulnerability | kev | `cve-2024-39891` | code | 1.0 |
| 5992 | CrushFTP VFS Sandbox Escape Vulnerability | kev | `cve-2024-4040` | code | 1.0 |
| 5993 | Veeam Backup and Replication Deserialization Vulnerability | kev | `cve-2024-40711` | code | 1.0 |
| 5994 | SonicWall SonicOS Improper Access Control Vulnerability | kev | `cve-2024-40766` | code | 1.0 |
| 5995 | Zyxel DSL CPE OS Command Injection Vulnerability | kev | `cve-2024-40890` | code | 1.0 |
| 5996 | Zyxel DSL CPE OS Command Injection Vulnerability | kev | `cve-2024-40891` | code | 1.0 |
| 5997 | Mitel SIP Phones Argument Injection Vulnerability | kev | `cve-2024-41710` | code | 1.0 |
| 5998 | Mitel MiCollab Path Traversal Vulnerability | kev | `cve-2024-41713` | code | 1.0 |
| 5999 | RoundCube Webmail Cross-Site Scripting Vulnerability | kev | `cve-2024-42009` | code | 1.0 |
| 6000 | Qualcomm Multiple Chipsets Use-After-Free Vulnerability | kev | `cve-2024-43047` | code | 1.0 |
| 6001 | Android Framework Privilege Escalation Vulnerability | kev | `cve-2024-43093` | code | 1.0 |
| 6002 | Microsoft Windows NTLMv2 Hash Disclosure Spoofing Vulnerability | kev | `cve-2024-43451` | code | 1.0 |
| 6003 | Microsoft Windows MSHTML Platform Spoofing Vulnerability | kev | `cve-2024-43461` | code | 1.0 |
| 6004 | Microsoft Configuration Manager SQL Injection Vulnerability | kev | `cve-2024-43468` | code | 1.0 |
| 6005 | Microsoft Windows Management Console Remote Code Execution Vulnerability | kev | `cve-2024-43572` | code | 1.0 |
| 6006 | Microsoft Windows MSHTML Platform Spoofing Vulnerability | kev | `cve-2024-43573` | code | 1.0 |
| 6007 | Progress Telerik Report Server Authentication Bypass by Spoofing Vulnerability | kev | `cve-2024-4358` | code | 1.0 |
| 6008 | Apple Multiple Products Code Execution Vulnerability | kev | `cve-2024-44308` | code | 1.0 |
| 6009 | Apple Multiple Products Cross-Site Scripting (XSS) Vulnerability | kev | `cve-2024-44309` | code | 1.0 |
| 6010 | Apache OFBiz Forced Browsing Vulnerability | kev | `cve-2024-45195` | code | 1.0 |
| 6011 | Synacor Zimbra Collaboration Suite (ZCS) Command Execution Vulnerability | kev | `cve-2024-45519` | code | 1.0 |
| 6012 | PHP-CGI OS Command Injection Vulnerability | kev | `cve-2024-4577` | code | 1.0 |
| 6013 | Arm Mali GPU Kernel Driver Use-After-Free Vulnerability | kev | `cve-2024-4610` | code | 1.0 |
| 6014 | Google Chromium Visuals Use-After-Free Vulnerability | kev | `cve-2024-4671` | code | 1.0 |
| 6015 | Fortinet FortiManager Missing Authentication Vulnerability | kev | `cve-2024-47575` | code | 1.0 |
| 6016 | Google Chromium V8 Out-of-Bounds Memory Write Vulnerability | kev | `cve-2024-4761` | code | 1.0 |
| 6017 | NAKIVO Backup and Replication Absolute Path Traversal Vulnerability | kev | `cve-2024-48248` | code | 1.0 |
| 6018 | ServiceNow Improper Input Validation Vulnerability | kev | `cve-2024-4879` | code | 1.0 |
| 6019 | Progress WhatsUp Gold Path Traversal Vulnerability | kev | `cve-2024-4885` | code | 1.0 |
| 6020 | Microsoft Partner Center Improper Access Control Vulnerability | kev | `cve-2024-49035` | code | 1.0 |
| 6021 | Microsoft Windows Task Scheduler Privilege Escalation Vulnerability | kev | `cve-2024-49039` | code | 1.0 |
| 6022 | Microsoft Windows Common Log File System (CLFS) Driver Heap-Based Buffer Overflow Vulnerability | kev | `cve-2024-49138` | code | 1.0 |
| 6023 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2024-4947` | code | 1.0 |
| 6024 | Justice AV Solutions (JAVS) Viewer Installer Embedded Malicious Code Vulnerability | kev | `cve-2024-4978` | code | 1.0 |
| 6025 | Linux Kernel Use of Uninitialized Resource Vulnerability | kev | `cve-2024-50302` | code | 1.0 |
| 6026 | Aviatrix Controllers OS Command Injection Vulnerability | kev | `cve-2024-50603` | code | 1.0 |
| 6027 | Cleo Multiple Products Unrestricted File Upload Vulnerability | kev | `cve-2024-50623` | code | 1.0 |
| 6028 | CyberPanel Incorrect Default Permissions Vulnerability | kev | `cve-2024-51378` | code | 1.0 |
| 6029 | CyberPanel Incorrect Default Permissions Vulnerability | kev | `cve-2024-51567` | code | 1.0 |
| 6030 | ServiceNow Incomplete List of Disallowed Inputs Vulnerability | kev | `cve-2024-5217` | code | 1.0 |
| 6031 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2024-5274` | code | 1.0 |
| 6032 | Linux Kernel Out-of-Bounds Write Vulnerability | kev | `cve-2024-53104` | code | 1.0 |
| 6033 | Linux Kernel Out-of-Bounds Read Vulnerability | kev | `cve-2024-53150` | code | 1.0 |
| 6034 | Linux Kernel Out-of-Bounds Access Vulnerability | kev | `cve-2024-53197` | code | 1.0 |
| 6035 | SonicWall SonicOS SSLVPN Improper Authentication Vulnerability | kev | `cve-2024-53704` | code | 1.0 |
| 6036 | AMI MegaRAC SPx Authentication Bypass by Spoofing Vulnerability | kev | `cve-2024-54085` | code | 1.0 |
| 6037 | Mitel MiCollab Path Traversal Vulnerability | kev | `cve-2024-55550` | code | 1.0 |
| 6038 | Fortinet FortiOS and FortiProxy Authentication Bypass Vulnerability | kev | `cve-2024-55591` | code | 1.0 |
| 6039 | Cleo Multiple Products Unauthenticated File Upload Vulnerability | kev | `cve-2024-55956` | code | 1.0 |
| 6040 | Craft CMS Code Injection Vulnerability | kev | `cve-2024-56145` | code | 1.0 |
| 6041 | SimpleHelp Missing Authorization Vulnerability | kev | `cve-2024-57726` | code | 1.0 |
| 6042 | SimpleHelp Path Traversal Vulnerability | kev | `cve-2024-57727` | code | 1.0 |
| 6043 | SimpleHelp Path Traversal Vulnerability | kev | `cve-2024-57728` | code | 1.0 |
| 6044 | Advantive VeraCore Unrestricted File Upload Vulnerability | kev | `cve-2024-57968` | code | 1.0 |
| 6045 | Yiiframework Yii Improper Protection of Alternate Path Vulnerability | kev | `cve-2024-58136` | code | 1.0 |
| 6046 | Palo Alto Networks Expedition Missing Authentication Vulnerability | kev | `cve-2024-5910` | code | 1.0 |
| 6047 | GeoVision Devices OS Command Injection Vulnerability | kev | `cve-2024-6047` | code | 1.0 |
| 6048 | Progress WhatsUp Gold SQL Injection Vulnerability | kev | `cve-2024-6670` | code | 1.0 |
| 6049 | Kingsoft WPS Office Path Traversal Vulnerability | kev | `cve-2024-7262` | code | 1.0 |
| 6050 | Samsung MagicINFO 9 Server Path Traversal Vulnerability | kev | `cve-2024-7399` | code | 1.0 |
| 6051 | Ivanti Virtual Traffic Manager Authentication Bypass Vulnerability | kev | `cve-2024-7593` | code | 1.0 |
| 6052 | TeamT5 ThreatSonar Anti-Ransomware Unrestricted Upload of File with Dangerous Type Vulnerability | kev | `cve-2024-7694` | code | 1.0 |
| 6053 | Google Chromium V8 Inappropriate Implementation Vulnerability | kev | `cve-2024-7965` | code | 1.0 |
| 6054 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2024-7971` | code | 1.0 |
| 6055 | Citrix Session Recording Improper Privilege Management Vulnerability | kev | `cve-2024-8068` | code | 1.0 |
| 6056 | Citrix Session Recording Deserialization of Untrusted Data Vulnerability | kev | `cve-2024-8069` | code | 1.0 |
| 6057 | Ivanti Cloud Services Appliance OS Command Injection Vulnerability | kev | `cve-2024-8190` | code | 1.0 |
| 6058 | PTZOptics PT30X-SDI/NDI Cameras Authentication Bypass Vulnerability | kev | `cve-2024-8956` | code | 1.0 |
| 6059 | PTZOptics PT30X-SDI/NDI Cameras OS Command Injection Vulnerability | kev | `cve-2024-8957` | code | 1.0 |
| 6060 | Ivanti Cloud Services Appliance (CSA) Path Traversal Vulnerability | kev | `cve-2024-8963` | code | 1.0 |
| 6061 | Ivanti Cloud Services Appliance (CSA) SQL Injection Vulnerability | kev | `cve-2024-9379` | code | 1.0 |
| 6062 | Ivanti Cloud Services Appliance (CSA) OS Command Injection Vulnerability | kev | `cve-2024-9380` | code | 1.0 |
| 6063 | Palo Alto Networks Expedition OS Command Injection Vulnerability | kev | `cve-2024-9463` | code | 1.0 |
| 6064 | Palo Alto Networks Expedition SQL Injection Vulnerability | kev | `cve-2024-9465` | code | 1.0 |
| 6065 | Palo Alto Networks PAN-OS Management Interface OS Command Injection Vulnerability | kev | `cve-2024-9474` | code | 1.0 |
| 6066 | ScienceLogic SL1 Unspecified Vulnerability | kev | `cve-2024-9537` | code | 1.0 |
| 6067 | Mozilla Firefox Use-After-Free Vulnerability | kev | `cve-2024-9680` | code | 1.0 |
| 6068 | Palo Alto Networks PAN-OS Authentication Bypass Vulnerability | kev | `cve-2025-0108` | code | 1.0 |
| 6069 | Palo Alto Networks PAN-OS File Read Vulnerability | kev | `cve-2025-0111` | code | 1.0 |
| 6070 | Ivanti Connect Secure, Policy Secure, and ZTA Gateways Stack-Based Buffer Overflow Vulnerability | kev | `cve-2025-0282` | code | 1.0 |
| 6071 | 7-Zip Mark of the Web Bypass Vulnerability | kev | `cve-2025-0411` | code | 1.0 |
| 6072 | Trimble Cityworks Deserialization Vulnerability | kev | `cve-2025-0994` | code | 1.0 |
| 6073 | Fortra GoAnywhere MFT Deserialization of Untrusted Data Vulnerability | kev | `cve-2025-10035` | code | 1.0 |
| 6074 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2025-10585` | code | 1.0 |
| 6075 | Gladinet CentreStack and Triofox Files or Directories Accessible to External Parties Vulnerability | kev | `cve-2025-11371` | code | 1.0 |
| 6076 | React Native Community CLI OS Command Injection Vulnerability | kev | `cve-2025-11953` | code | 1.0 |
| 6077 | Gladinet Triofox Improper Access Control Vulnerability | kev | `cve-2025-12480` | code | 1.0 |
| 6078 | Edimax IC-7100 IP Camera OS Command Injection Vulnerability | kev | `cve-2025-1316` | code | 1.0 |
| 6079 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2025-13223` | code | 1.0 |
| 6080 | Google Chromium Out of Bounds Memory Access Vulnerability | kev | `cve-2025-14174` | code | 1.0 |
| 6081 | Gladinet CentreStack and Triofox Hard Coded Cryptographic Vulnerability | kev | `cve-2025-14611` | code | 1.0 |
| 6082 | WatchGuard Firebox Out of Bounds Write Vulnerability | kev | `cve-2025-14733` | code | 1.0 |
| 6083 | MongoDB and MongoDB Server Improper Handling of Length Parameter Inconsistency Vulnerability | kev | `cve-2025-14847` | code | 1.0 |
| 6084 | Notepad++ Download of Code Without Integrity Check Vulnerability | kev | `cve-2025-15556` | code | 1.0 |
| 6085 | Broadcom Brocade Fabric OS Code Injection Vulnerability | kev | `cve-2025-1976` | code | 1.0 |
| 6086 | Cisco Identity Services Engine Injection Vulnerability | kev | `cve-2025-20281` | code | 1.0 |
| 6087 | Cisco Secure Firewall Adaptive Security Appliance (ASA) and Secure Firewall Threat Defense (FTD) Buffer Overflow Vulnerability | kev | `cve-2025-20333` | code | 1.0 |
| 6088 | Cisco Identity Services Engine Injection Vulnerability | kev | `cve-2025-20337` | code | 1.0 |
| 6089 | Cisco IOS and IOS XE Software SNMP Denial of Service and Remote Code Execution Vulnerability | kev | `cve-2025-20352` | code | 1.0 |
| 6090 | Cisco Secure Firewall Adaptive Security (ASA) Appliance and Secure Firewall Threat Defense (FTD) Missing Authorization Vulnerability | kev | `cve-2025-20362` | code | 1.0 |
| 6091 | Cisco Multiple Products Improper Input Validation Vulnerability | kev | `cve-2025-20393` | code | 1.0 |
| 6092 | Samsung Mobile Devices Out-of-Bounds Write Vulnerability | kev | `cve-2025-21042` | code | 1.0 |
| 6093 | Samsung Mobile Devices Out-of-Bounds Write Vulnerability | kev | `cve-2025-21043` | code | 1.0 |
| 6094 | Microsoft Windows Hyper-V NT Kernel Integration VSP Heap-based Buffer Overflow Vulnerability | kev | `cve-2025-21333` | code | 1.0 |
| 6095 | Microsoft Windows Hyper-V NT Kernel Integration VSP Use-After-Free Vulnerability | kev | `cve-2025-21334` | code | 1.0 |
| 6096 | Microsoft Windows Hyper-V NT Kernel Integration VSP Use-After-Free Vulnerability | kev | `cve-2025-21335` | code | 1.0 |
| 6097 | Microsoft Windows Storage Link Following Vulnerability | kev | `cve-2025-21391` | code | 1.0 |
| 6098 | Microsoft Windows Ancillary Function Driver for WinSock Heap-Based Buffer Overflow Vulnerability | kev | `cve-2025-21418` | code | 1.0 |
| 6099 | Qualcomm Multiple Chipsets Incorrect Authorization Vulnerability | kev | `cve-2025-21479` | code | 1.0 |
| 6100 | Qualcomm Multiple Chipsets Incorrect Authorization Vulnerability | kev | `cve-2025-21480` | code | 1.0 |
| 6101 | Juniper Junos OS Improper Isolation or Compartmentalization Vulnerability | kev | `cve-2025-21590` | code | 1.0 |
| 6102 | VMware ESXi and Workstation TOCTOU Race Condition Vulnerability | kev | `cve-2025-22224` | code | 1.0 |
| 6103 | VMware ESXi Arbitrary Write Vulnerability | kev | `cve-2025-22225` | code | 1.0 |
| 6104 | VMware ESXi, Workstation, and Fusion Information Disclosure Vulnerability | kev | `cve-2025-22226` | code | 1.0 |
| 6105 | Ivanti Connect Secure, Policy Secure, and ZTA Gateways Stack-Based Buffer Overflow Vulnerability | kev | `cve-2025-22457` | code | 1.0 |
| 6106 | SonicWall SMA1000 Appliances Deserialization Vulnerability | kev | `cve-2025-23006` | code | 1.0 |
| 6107 | Craft CMS Code Injection Vulnerability | kev | `cve-2025-23209` | code | 1.0 |
| 6108 | Wazuh Server Deserialization of Untrusted Data Vulnerability | kev | `cve-2025-24016` | code | 1.0 |
| 6109 | Microsoft Windows NTLM Hash Disclosure Spoofing Vulnerability | kev | `cve-2025-24054` | code | 1.0 |
| 6110 | Apple Multiple Products Use-After-Free Vulnerability | kev | `cve-2025-24085` | code | 1.0 |
| 6111 | Apple iOS and iPadOS Incorrect Authorization Vulnerability | kev | `cve-2025-24200` | code | 1.0 |
| 6112 | Apple Multiple Products WebKit Out-of-Bounds Write Vulnerability | kev | `cve-2025-24201` | code | 1.0 |
| 6113 | Fortinet FortiOS and FortiProxy Authentication Bypass Vulnerability | kev | `cve-2025-24472` | code | 1.0 |
| 6114 | Apache Tomcat Path Equivalence Vulnerability | kev | `cve-2025-24813` | code | 1.0 |
| 6115 | XWiki Platform Eval Injection Vulnerability | kev | `cve-2025-24893` | code | 1.0 |
| 6116 | Microsoft Windows Win32k Use-After-Free Vulnerability | kev | `cve-2025-24983` | code | 1.0 |
| 6117 | Microsoft Windows NTFS Information Disclosure Vulnerability | kev | `cve-2025-24984` | code | 1.0 |
| 6118 | Microsoft Windows Fast FAT File System Driver Integer Overflow Vulnerability | kev | `cve-2025-24985` | code | 1.0 |
| 6119 | Microsoft Power Pages Improper Access Control Vulnerability | kev | `cve-2025-24989` | code | 1.0 |
| 6120 | Microsoft Windows Untrusted Pointer Dereference Vulnerability | kev | `cve-2025-24990` | code | 1.0 |
| 6121 | Microsoft Windows NTFS Out-Of-Bounds Read Vulnerability | kev | `cve-2025-24991` | code | 1.0 |
| 6122 | Microsoft Windows NTFS Heap-Based Buffer Overflow Vulnerability | kev | `cve-2025-24993` | code | 1.0 |
| 6123 | Advantive VeraCore SQL Injection Vulnerability | kev | `cve-2025-25181` | code | 1.0 |
| 6124 | Fortinet Multiple Products Heap-based Buffer Overflow Vulnerability | kev | `cve-2025-25249` | code | 1.0 |
| 6125 | Fortinet FortiWeb SQL Injection Vulnerability | kev | `cve-2025-25257` | code | 1.0 |
| 6126 | SolarWinds Web Help Desk Deserialization of Untrusted Data Vulnerability | kev | `cve-2025-26399` | code | 1.0 |
| 6127 | Microsoft Windows Management Console (MMC) Improper Neutralization Vulnerability | kev | `cve-2025-26633` | code | 1.0 |
| 6128 | Qualcomm Multiple Chipsets Use-After-Free Vulnerability | kev | `cve-2025-27038` | code | 1.0 |
| 6129 | FreeType Out-of-Bounds Write Vulnerability | kev | `cve-2025-27363` | code | 1.0 |
| 6130 | Kentico Xperience CMS Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | `cve-2025-2746` | code | 1.0 |
| 6131 | Kentico Xperience CMS Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | `cve-2025-2747` | code | 1.0 |
| 6132 | Kentico Xperience Path Traversal Vulnerability | kev | `cve-2025-2749` | code | 1.0 |
| 6133 | SysAid On-Prem Improper Restriction of XML External Entity Reference Vulnerability | kev | `cve-2025-2775` | code | 1.0 |
| 6134 | SysAid On-Prem Improper Restriction of XML External Entity Reference Vulnerability | kev | `cve-2025-2776` | code | 1.0 |
| 6135 | Google Chromium Mojo Sandbox Escape Vulnerability | kev | `cve-2025-2783` | code | 1.0 |
| 6136 | Synacor Zimbra Collaboration Suite (ZCS) Cross-site Scripting Vulnerability | kev | `cve-2025-27915` | code | 1.0 |
| 6137 | Srimax Output Messenger Directory Traversal Vulnerability | kev | `cve-2025-27920` | code | 1.0 |
| 6138 | D-Link DIR-823X Command Injection Vulnerability | kev | `cve-2025-29635` | code | 1.0 |
| 6139 | Microsoft Windows Common Log File System (CLFS) Driver Use-After-Free Vulnerability | kev | `cve-2025-29824` | code | 1.0 |
| 6140 | tj-actions/changed-files GitHub Action Embedded Malicious Code Vulnerability | kev | `cve-2025-30066` | code | 1.0 |
| 6141 | reviewdog/action-setup GitHub Action Embedded Malicious Code Vulnerability | kev | `cve-2025-30154` | code | 1.0 |
| 6142 | Microsoft Windows Scripting Engine Type Confusion Vulnerability | kev | `cve-2025-30397` | code | 1.0 |
| 6143 | Microsoft Windows DWM Core Library Use-After-Free Vulnerability | kev | `cve-2025-30400` | code | 1.0 |
| 6144 | Gladinet CentreStack and Triofox Use of Hard-coded Cryptographic Key Vulnerability | kev | `cve-2025-30406` | code | 1.0 |
| 6145 | Vite Vitejs Improper Access Control Vulnerability | kev | `cve-2025-31125` | code | 1.0 |
| 6146 | CrushFTP Authentication Bypass Vulnerability | kev | `cve-2025-31161` | code | 1.0 |
| 6147 | Apple Multiple Products Memory Corruption Vulnerability | kev | `cve-2025-31200` | code | 1.0 |
| 6148 | Apple Multiple Products Arbitrary Read and Write Vulnerability | kev | `cve-2025-31201` | code | 1.0 |
| 6149 | Apple Multiple Products Buffer Overflow Vulnerability | kev | `cve-2025-31277` | code | 1.0 |
| 6150 | SAP NetWeaver Unrestricted File Upload Vulnerability | kev | `cve-2025-31324` | code | 1.0 |
| 6151 | Craft CMS Code Injection Vulnerability | kev | `cve-2025-32432` | code | 1.0 |
| 6152 | Erlang Erlang/OTP SSH Server Missing Authentication for Critical Function Vulnerability | kev | `cve-2025-32433` | code | 1.0 |
| 6153 | Sudo Inclusion of Functionality from Untrusted Control Sphere Vulnerability | kev | `cve-2025-32463` | code | 1.0 |
| 6154 | Langflow Missing Authentication Vulnerability | kev | `cve-2025-3248` | code | 1.0 |
| 6155 | Microsoft Windows Common Log File System (CLFS) Driver Use-After-Free Vulnerability | kev | `cve-2025-32701` | code | 1.0 |
| 6156 | Microsoft Windows Common Log File System (CLFS) Driver Heap-Based Buffer Overflow Vulnerability | kev | `cve-2025-32706` | code | 1.0 |
| 6157 | Microsoft Windows Ancillary Function Driver for WinSock Use-After-Free Vulnerability | kev | `cve-2025-32709` | code | 1.0 |
| 6158 | Fortinet Multiple Products Stack-Based Buffer Overflow Vulnerability | kev | `cve-2025-32756` | code | 1.0 |
| 6159 | Quest KACE Systems Management Appliance (SMA) Improper Authentication Vulnerability | kev | `cve-2025-32975` | code | 1.0 |
| 6160 | Microsoft Windows External Control of File Name or Path Vulnerability | kev | `cve-2025-33053` | code | 1.0 |
| 6161 | Microsoft Windows SMB Client Improper Access Control Vulnerability | kev | `cve-2025-33073` | code | 1.0 |
| 6162 | Versa Concerto Improper Authentication Vulnerability | kev | `cve-2025-34026` | code | 1.0 |
| 6163 | Commvault Command Center Path Traversal Vulnerability | kev | `cve-2025-34028` | code | 1.0 |
| 6164 | Langflow Origin Validation Error Vulnerability | kev | `cve-2025-34291` | code | 1.0 |
| 6165 | Craft CMS External Control of Assumed-Immutable Web Parameter Vulnerability | kev | `cve-2025-35939` | code | 1.0 |
| 6166 | Hewlett Packard Enterprise (HPE) OneView Code Injection Vulnerability | kev | `cve-2025-37164` | code | 1.0 |
| 6167 | Linux Kernel Time-of-Check Time-of-Use (TOCTOU) Race Condition Vulnerability | kev | `cve-2025-38352` | code | 1.0 |
| 6168 | Commvault Web Server Unspecified Vulnerability | kev | `cve-2025-3928` | code | 1.0 |
| 6169 | ConnectWise ScreenConnect Improper Authentication Vulnerability | kev | `cve-2025-3935` | code | 1.0 |
| 6170 | Linux Kernel Improper Check for Unusual or Exceptional Conditions Vulnerability | kev | `cve-2025-39682` | code | 1.0 |
| 6171 | Linux Kernel Race Condition Vulnerability | kev | `cve-2025-39964` | code | 1.0 |
| 6172 | Smartbedded Meteobridge Command Injection Vulnerability | kev | `cve-2025-4008` | code | 1.0 |
| 6173 | SolarWinds Web Help Desk Security Control Bypass Vulnerability | kev | `cve-2025-40536` | code | 1.0 |
| 6174 | SolarWinds Web Help Desk Deserialization of Untrusted Data Vulnerability | kev | `cve-2025-40551` | code | 1.0 |
| 6175 | SonicWall SMA1000 Missing Authorization Vulnerability | kev | `cve-2025-40602` | code | 1.0 |
| 6176 | Broadcom VMware Aria Operations and VMware Tools Privilege Defined with Unsafe Actions Vulnerability | kev | `cve-2025-41244` | code | 1.0 |
| 6177 | Danfoss AK-SM 800A Series Improper Authentication Vulnerability | kev | `cve-2025-41450` | code | 1.0 |
| 6178 | Qualitia Active! Mail Stack-Based Buffer Overflow Vulnerability | kev | `cve-2025-42599` | code | 1.0 |
| 6179 | SAP NetWeaver Improper Neutralization of Special Elements in Data Query Logic Vulnerability | kev | `cve-2025-42884` | code | 1.0 |
| 6180 | SAP NetWeaver Deserialization Vulnerability | kev | `cve-2025-42999` | code | 1.0 |
| 6181 | Apple Multiple Products Unspecified Vulnerability | kev | `cve-2025-43200` | code | 1.0 |
| 6182 | Apple iOS, iPadOS, and macOS Out-of-Bounds Write Vulnerability | kev | `cve-2025-43300` | code | 1.0 |
| 6183 | Apple Multiple Products Improper Locking Vulnerability | kev | `cve-2025-43510` | code | 1.0 |
| 6184 | Apple Multiple Products Classic Buffer Overflow Vulnerability | kev | `cve-2025-43520` | code | 1.0 |
| 6185 | Apple Multiple Products Use-After-Free WebKit Vulnerability | kev | `cve-2025-43529` | code | 1.0 |
| 6186 | Ivanti Endpoint Manager Mobile (EPMM) Authentication Bypass Vulnerability | kev | `cve-2025-4427` | code | 1.0 |
| 6187 | Ivanti Endpoint Manager Mobile (EPMM) Code Injection Vulnerability | kev | `cve-2025-4428` | code | 1.0 |
| 6188 | Samsung MagicINFO 9 Server Path Traversal Vulnerability | kev | `cve-2025-4632` | code | 1.0 |
| 6189 | TeleMessage TM SGNL Hidden Functionality Vulnerability | kev | `cve-2025-47729` | code | 1.0 |
| 6190 | Wing FTP Server Improper Neutralization of Null Byte or NUL Character Vulnerability | kev | `cve-2025-47812` | code | 1.0 |
| 6191 | Wing FTP Server Information Disclosure Vulnerability | kev | `cve-2025-47813` | code | 1.0 |
| 6192 | IGEL OS Use of a Key Past its Expiration Date Vulnerability | kev | `cve-2025-47827` | code | 1.0 |
| 6193 | Git Link Following Vulnerability | kev | `cve-2025-48384` | code | 1.0 |
| 6194 | Android Runtime Use-After-Free Vulnerability | kev | `cve-2025-48543` | code | 1.0 |
| 6195 | Android Framework Privilege Escalation Vulnerability | kev | `cve-2025-48572` | code | 1.0 |
| 6196 | Android Framework Integer Overflow Vulnerability | kev | `cve-2025-48595` | code | 1.0 |
| 6197 | Android Framework Information Disclosure Vulnerability | kev | `cve-2025-48633` | code | 1.0 |
| 6198 | Synacor Zimbra Collaboration Suite (ZCS) Cross-site Scripting Vulnerability | kev | `cve-2025-48700` | code | 1.0 |
| 6199 | CWP Control Web Panel OS Command Injection Vulnerability | kev | `cve-2025-48703` | code | 1.0 |
| 6200 | TeleMessage TM SGNL Initialization of a Resource with an Insecure Default Vulnerability | kev | `cve-2025-48927` | code | 1.0 |
| 6201 | TeleMessage TM SGNL Exposure of Core Dump File to an Unauthorized Control Sphere Vulnerability | kev | `cve-2025-48928` | code | 1.0 |
| 6202 | RoundCube Webmail Deserialization of Untrusted Data Vulnerability | kev | `cve-2025-49113` | code | 1.0 |
| 6203 | Microsoft SharePoint Code Injection Vulnerability | kev | `cve-2025-49704` | code | 1.0 |
| 6204 | Microsoft SharePoint Improper Authentication Vulnerability | kev | `cve-2025-49706` | code | 1.0 |
| 6205 | Dassault Systèmes DELMIA Apriso Deserialization of Untrusted Data Vulnerability | kev | `cve-2025-5086` | code | 1.0 |
| 6206 | SmarterTools SmarterMail Unrestricted Upload of File with Dangerous Type Vulnerability | kev | `cve-2025-52691` | code | 1.0 |
| 6207 | F5 BIG-IP Stack-Based Buffer Overflow Vulnerability | kev | `cve-2025-53521` | code | 1.0 |
| 6208 | Sitecore Multiple Products Deserialization of Untrusted Data Vulnerability | kev | `cve-2025-53690` | code | 1.0 |
| 6209 | Microsoft SharePoint Deserialization of Untrusted Data Vulnerability | kev | `cve-2025-53770` | code | 1.0 |
| 6210 | Laravel Livewire Code Injection Vulnerability | kev | `cve-2025-54068` | code | 1.0 |
| 6211 | Google Chromium V8 Out-of-Bounds Read and Write Vulnerability | kev | `cve-2025-5419` | code | 1.0 |
| 6212 | Adobe Commerce and Magento Improper Input Validation Vulnerability | kev | `cve-2025-54236` | code | 1.0 |
| 6213 | Adobe Experience Manager Forms Code Execution Vulnerability | kev | `cve-2025-54253` | code | 1.0 |
| 6214 | CrushFTP Unprotected Alternate Channel Vulnerability | kev | `cve-2025-54309` | code | 1.0 |
| 6215 | Prettier eslint-config-prettier Embedded Malicious Code Vulnerability | kev | `cve-2025-54313` | code | 1.0 |
| 6216 | Trend Micro Apex One OS Command Injection Vulnerability | kev | `cve-2025-54948` | code | 1.0 |
| 6217 | Meta Platforms WhatsApp Incorrect Authorization Vulnerability | kev | `cve-2025-55177` | code | 1.0 |
| 6218 | Meta React Server Components Remote Code Execution Vulnerability | kev | `cve-2025-55182` | code | 1.0 |
| 6219 | Citrix NetScaler ADC and Gateway Out-of-Bounds Read Vulnerability | kev | `cve-2025-5777` | code | 1.0 |
| 6220 | Sangoma FreePBX Authentication Bypass Vulnerability | kev | `cve-2025-57819` | code | 1.0 |
| 6221 | Fortinet FortiWeb OS Command Injection Vulnerability | kev | `cve-2025-58034` | code | 1.0 |
| 6222 | OSGeo GeoServer Improper Restriction of XML External Entity Reference Vulnerability | kev | `cve-2025-58360` | code | 1.0 |
| 6223 | Microsoft Windows Improper Access Control Vulnerability | kev | `cve-2025-59230` | code | 1.0 |
| 6224 | Microsoft Windows Server Update Service (WSUS) Deserialization of Untrusted Data Vulnerability | kev | `cve-2025-59287` | code | 1.0 |
| 6225 | ASUS Live Update Embedded Malicious Code Vulnerability | kev | `cve-2025-59374` | code | 1.0 |
| 6226 | Libraesva Email Security Gateway Command Injection Vulnerability | kev | `cve-2025-59689` | code | 1.0 |
| 6227 | Fortinet Multiple Products Improper Verification of Cryptographic Signature Vulnerability | kev | `cve-2025-59718` | code | 1.0 |
| 6228 | Microsoft Windows Link Following Vulnerability | kev | `cve-2025-60710` | code | 1.0 |
| 6229 | Fortinet Multiple Products Path Traversal Vulnerability | kev | `cve-2025-61624` | code | 1.0 |
| 6230 | Oracle Fusion Middleware Missing Authentication for Critical Function Vulnerability | kev | `cve-2025-61757` | code | 1.0 |
| 6231 | Oracle E-Business Suite Unspecified Vulnerability | kev | `cve-2025-61882` | code | 1.0 |
| 6232 | Oracle E-Business Suite Server-Side Request Forgery (SSRF) Vulnerability | kev | `cve-2025-61884` | code | 1.0 |
| 6233 | Motex LANSCOPE Endpoint Manager Improper Verification of Source of a Communication Channel Vulnerability | kev | `cve-2025-61932` | code | 1.0 |
| 6234 | Dassault Systèmes DELMIA Apriso Code Injection Vulnerability | kev | `cve-2025-6204` | code | 1.0 |
| 6235 | Dassault Systèmes DELMIA Apriso Missing Authorization Vulnerability | kev | `cve-2025-6205` | code | 1.0 |
| 6236 | RARLAB WinRAR Path Traversal Vulnerability | kev | `cve-2025-6218` | code | 1.0 |
| 6237 | Microsoft Windows Race Condition Vulnerability | kev | `cve-2025-62215` | code | 1.0 |
| 6238 | Microsoft Windows Use After Free Vulnerability | kev | `cve-2025-62221` | code | 1.0 |
| 6239 | Ray-Project Ray Code Injection Vulnerability | kev | `cve-2025-62593` | code | 1.0 |
| 6240 | Sangoma FreePBX OS Command Injection Vulnerability | kev | `cve-2025-64328` | code | 1.0 |
| 6241 | Fortinet FortiWeb Path Traversal Vulnerability | kev | `cve-2025-64446` | code | 1.0 |
| 6242 | Citrix NetScaler ADC and Gateway Buffer Overflow Vulnerability | kev | `cve-2025-6543` | code | 1.0 |
| 6243 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2025-6554` | code | 1.0 |
| 6244 | Google Chromium ANGLE and GPU Improper Input Validation Vulnerability | kev | `cve-2025-6558` | code | 1.0 |
| 6245 | Synacor Zimbra Collaboration Suite (ZCS) Cross-Site Scripting Vulnerability | kev | `cve-2025-66376` | code | 1.0 |
| 6246 | Array Networks ArrayOS AG OS Command Injection Vulnerability | kev | `cve-2025-66644` | code | 1.0 |
| 6247 | Lantronix EDS5000 Code Injection Vulnerability | kev | `cve-2025-67038` | code | 1.0 |
| 6248 | RoundCube Webmail Cross-site Scripting Vulnerability | kev | `cve-2025-68461` | code | 1.0 |
| 6249 | n8n Improper Control of Dynamically-Managed Code Resources Vulnerability | kev | `cve-2025-68613` | code | 1.0 |
| 6250 | Synacor Zimbra Collaboration Suite (ZCS) PHP Remote File Inclusion Vulnerability | kev | `cve-2025-68645` | code | 1.0 |
| 6251 | Fortinet FortiOS Exposure of Sensitive Information to an Unauthorized Actor Vulnerability | kev | `cve-2025-68686` | code | 1.0 |
| 6252 | Citrix NetScaler Memory Overflow Vulnerability | kev | `cve-2025-7775` | code | 1.0 |
| 6253 | RARLAB WinRAR Path Traversal Vulnerability | kev | `cve-2025-8088` | code | 1.0 |
| 6254 | Gogs Path Traversal Vulnerability | kev | `cve-2025-8110` | code | 1.0 |
| 6255 | N-able N-Central Insecure Deserialization Vulnerability | kev | `cve-2025-8875` | code | 1.0 |
| 6256 | N-able N-Central Command Injection Vulnerability | kev | `cve-2025-8876` | code | 1.0 |
| 6257 | WatchGuard Firebox Out-of-Bounds Write Vulnerability | kev | `cve-2025-9242` | code | 1.0 |
| 6258 | TP-Link Archer C7(EU) and TL-WR841N/ND(MS) OS Command Injection Vulnerability | kev | `cve-2025-9377` | code | 1.0 |
| 6259 | Palo Alto Networks PAN-OS Authentication Bypass Vulnerability | kev | `cve-2026-0257` | code | 1.0 |
| 6260 | Palo Alto Networks PAN-OS Out-of-bounds Write Vulnerability | kev | `cve-2026-0300` | code | 1.0 |
| 6261 | Palo Alto Networks PAN-OS Code Injection Vulnerability | kev | `cve-2026-0322` | code | 1.0 |
| 6262 | Langflow Inclusion of Functionality from Untrusted Control Sphere Vulnerability | kev | `cve-2026-0770` | code | 1.0 |
| 6263 | Zammad GmbH Zammad Session Fixation Vulnerability | kev | `cve-2026-102489` | code | 1.0 |
| 6264 | Zammad GmbH Zammad Improper Privilege Management Vulnerability | kev | `cve-2026-102490` | code | 1.0 |
| 6265 | BeyondTrust Remote Support Authentication Bypass by Spoofing Vulnerability | kev | `cve-2026-102764` | code | 1.0 |
| 6266 | GPUStack Improper Authentication Vulnerability | kev | `cve-2026-102895` | code | 1.0 |
| 6267 | Hewlett Packard Enterprise (HPE) Juniper Session Smart Conductor and Session Smart Router OS Command Injection Vulnerability | kev | `cve-2026-103406` | code | 1.0 |
| 6268 | Fortinet FortiMail Path Traversal Vulnerability | kev | `cve-2026-104286` | code | 1.0 |
| 6269 | Ivanti Connect Secure Authentication Bypass by Spoofing Vulnerability | kev | `cve-2026-104935` | code | 1.0 |
| 6270 | Ahsay AhsayCBS Improper Authentication Vulnerability | kev | `cve-2026-105133` | code | 1.0 |
| 6271 | Ahsay AhsayCBS Command Injection Vulnerability | kev | `cve-2026-105134` | code | 1.0 |
| 6272 | Ivanti Sentry OS Command Injection Vulnerability | kev | `cve-2026-10520` | code | 1.0 |
| 6273 | EG Secure Solutions SiteGuard Server Edition OS Command Injection Vulnerability | kev | `cve-2026-106931` | code | 1.0 |
| 6274 | PHP Memory Corruption Vulnerability | kev | `cve-2026-109618` | code | 1.0 |
| 6275 | PHP Stale Internal Request-State Vulnerability | kev | `cve-2026-109619` | code | 1.0 |
| 6276 | Collabora Online – Built-in CODE Server OS Command Injection Vulnerability | kev | `cve-2026-110447` | code | 1.0 |
| 6277 | Siemens Software Center Uncontrolled Search Path Element Vulnerability | kev | `cve-2026-110936` | code | 1.0 |
| 6278 | Ivanti Connect Secure Missing Authentication for Critical Function Vulnerability | kev | `cve-2026-111357` | code | 1.0 |
| 6279 | Google Brotli Unchecked Return Value Vulnerability | kev | `cve-2026-111497` | code | 1.0 |
| 6280 | Barracuda Networks Multiple Products OS Command Injection Vulnerability | kev | `cve-2026-111560` | code | 1.0 |
| 6281 | Openwave Messaging Limited Webtop Media Transcoder Incorrect Resource Transfer Between Spheres Vulnerability | kev | `cve-2026-111793` | code | 1.0 |
| 6282 | Arista VeloCloud Orchestrator SQL Injection Vulnerability | kev | `cve-2026-114625` | code | 1.0 |
| 6283 | UnrealIRCd Improper Neutralization of Delimiters Vulnerability | kev | `cve-2026-114713` | code | 1.0 |
| 6284 | Perforce Puppet Enterprise Deserialization of Untrusted Data Vulnerability | kev | `cve-2026-115731` | code | 1.0 |
| 6285 | Google Chromium V8 Out-of-Bounds Read and Write Vulnerability | kev | `cve-2026-11645` | code | 1.0 |
| 6286 | Barracuda Networks Email Security Gateway (ESG) Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | `cve-2026-116606` | code | 1.0 |
| 6287 | SignalWire FreeSwitch Code Injection Vulnerability | kev | `cve-2026-116779` | code | 1.0 |
| 6288 | Ivanti Connect Secure Generation of Predictable Numbers or Identifiers Vulnerability | kev | `cve-2026-118527` | code | 1.0 |
| 6289 | WordPress WordPress Core Time-of-check Time-of-use (TOCTOU) Race Condition Vulnerability | kev | `cve-2026-119977` | code | 1.0 |
| 6290 | Check Point Multiple Products Incorrect Permission Assignment for Critical Resource Vulnerability | kev | `cve-2026-123039` | code | 1.0 |
| 6291 | OpenBSD OpenSSH Unexpected Status Code or Return Value Vulnerability | kev | `cve-2026-123158` | code | 1.0 |
| 6292 | IBM Aspera Faspex Incorrect Authorization Vulnerability | kev | `cve-2026-123218` | code | 1.0 |
| 6293 | VMWare vCenter Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | `cve-2026-123890` | code | 1.0 |
| 6294 | Microsoft Office Outlook Missing User Interaction for Potentially Unsafe Action Vulnerability | kev | `cve-2026-124194` | code | 1.0 |
| 6295 | Google Chromium Use-After-Free Vulnerability | kev | `cve-2026-125350` | code | 1.0 |
| 6296 | PTC Windchill and FlexPLM Improper Input Validation Vulnerability | kev | `cve-2026-12569` | code | 1.0 |
| 6297 | Ivanti Endpoint Manager Mobile (EPMM) Code Injection Vulnerability | kev | `cve-2026-1281` | code | 1.0 |
| 6298 | Apache Roller Deserialization of Untrusted Data Vulnerability | kev | `cve-2026-129448` | code | 1.0 |
| 6299 | Idunno AntiCsrf Deserialization of Untrusted Data Vulnerability | kev | `cve-2026-131914` | code | 1.0 |
| 6300 | Drupal Core Remote Code Execution Vulnerability | kev | `cve-2026-132165` | code | 1.0 |
| 6301 | Optimizely Web Experimentation External Control of Critical State Data Vulnerability | kev | `cve-2026-132224` | code | 1.0 |
| 6302 | Mitel MiCollab Unrestricted Upload of File with Dangerous Type Vulnerability | kev | `cve-2026-133275` | code | 1.0 |
| 6303 | Mitel MiCollab XML Injection Vulnerability | kev | `cve-2026-133303` | code | 1.0 |
| 6304 | Mitel MiCollab Weak Authentication Vulnerability | kev | `cve-2026-133311` | code | 1.0 |
| 6305 | Mitel MiCollab OS Command Injection Vulnerability | kev | `cve-2026-133323` | code | 1.0 |
| 6306 | Mitel MiCollab Not Failing Securely Vulnerability | kev | `cve-2026-133324` | code | 1.0 |
| 6307 | Vates Xen Orchestra Improper Validation of Specified Type of Input Vulnerability | kev | `cve-2026-133469` | code | 1.0 |
| 6308 | Project MONAI MONAI Label Path Traversal Vulnerability | kev | `cve-2026-133968` | code | 1.0 |
| 6309 | Ivanti Endpoint Manager Mobile (EPMM) Code Injection Vulnerability | kev | `cve-2026-1340` | code | 1.0 |
| 6310 | N-able MailAssure (SpamExperts) Server-Side Request Forgery (SSRF) Vulnerability | kev | `cve-2026-134624` | code | 1.0 |
| 6311 | SonicWall SMA1000 Appliances Server-Side Request Forgery Vulnerability | kev | `cve-2026-15409` | code | 1.0 |
| 6312 | SonicWall SMA1000 Appliances Code Injection Vulnerability | kev | `cve-2026-15410` | code | 1.0 |
| 6313 | Ivanti Endpoint Manager (EPM) Authentication Bypass Vulnerability | kev | `cve-2026-1603` | code | 1.0 |
| 6314 | Check Point SmartConsole Improper Authentication Vulnerability | kev | `cve-2026-16232` | code | 1.0 |
| 6315 | Arista VeloCloud Orchestrator On-Prem OS Command Injection Vulnerability | kev | `cve-2026-16812` | code | 1.0 |
| 6316 | BeyondTrust Remote Support (RS) and Privileged Remote Access (PRA) OS Command Injection Vulnerability | kev | `cve-2026-1731` | code | 1.0 |
| 6317 | N-able N-central Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | `cve-2026-18556` | code | 1.0 |
| 6318 | N-able N-central Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | `cve-2026-18577` | code | 1.0 |
| 6319 | Citrix NetScaler Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | `cve-2026-19490` | code | 1.0 |
| 6320 | SonicWall SonicOS SSLVPN Use of Single-factor Authentication Vulnerability | kev | `cve-2026-19868` | code | 1.0 |
| 6321 | Cisco Unified Communications Products Code Injection Vulnerability | kev | `cve-2026-20045` | code | 1.0 |
| 6322 | Cisco Firewall Management Center Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | `cve-2026-20079` | code | 1.0 |
| 6323 | Cisco Catalyst SD-WAN Manager Incorrect Use of Privileged APIs Vulnerability | kev | `cve-2026-20122` | code | 1.0 |
| 6324 | Cisco Catalyst SD-WAN Controller and Manager Authentication Bypass Vulnerability | kev | `cve-2026-20127` | code | 1.0 |
| 6325 | Cisco Catalyst SD-WAN Manager Storing Passwords in a Recoverable Format Vulnerability | kev | `cve-2026-20128` | code | 1.0 |
| 6326 | Cisco Secure Firewall Management Center (FMC) Software and Cisco Security Cloud Control (SCC) Firewall Management Deserialization of Untrusted Data Vulnerability | kev | `cve-2026-20131` | code | 1.0 |
| 6327 | Cisco Catalyst SD-WAN Manager Exposure of Sensitive Information to an Unauthorized Actor Vulnerability | kev | `cve-2026-20133` | code | 1.0 |
| 6328 | Cisco Catalyst SD-WAN Controller Authentication Bypass Vulnerability | kev | `cve-2026-20182` | code | 1.0 |
| 6329 | Cisco Unified Communications Manager Server-Side Request Forgery (SSRF) Vulnerability | kev | `cve-2026-20230` | code | 1.0 |
| 6330 | Cisco Catalyst SD-WAN Manager Improper Encoding or Escaping of Output Vulnerability | kev | `cve-2026-20245` | code | 1.0 |
| 6331 | Splunk Enterprise Missing Authentication for Critical Function Vulnerability | kev | `cve-2026-20253` | code | 1.0 |
| 6332 | Cisco Catalyst SD-WAN Manager Directory or Path Traversal Vulnerability | kev | `cve-2026-20262` | code | 1.0 |
| 6333 | Cisco Secure Firewall Management Center Use of Hard-coded Password Vulnerability | kev | `cve-2026-20316` | code | 1.0 |
| 6334 | Cisco BroadWorks Improper Restriction of XML External Entity Reference Vulnerability | kev | `cve-2026-20320` | code | 1.0 |
| 6335 | Cisco Secure Firewall Adaptive Security Appliance (ASA) and Secure Firewall Threat Defense (FTD) Heap Inspection Vulnerability | kev | `cve-2026-20349` | code | 1.0 |
| 6336 | Apple Multiple Buffer Overflow Vulnerability | kev | `cve-2026-20700` | code | 1.0 |
| 6337 | Microsoft Windows Information Disclosure Vulnerability | kev | `cve-2026-20805` | code | 1.0 |
| 6338 | Microsoft SharePoint Deserialization of Untrusted Data Vulnerability | kev | `cve-2026-20963` | code | 1.0 |
| 6339 | Qualcomm Multiple Chipsets Memory Corruption Vulnerability | kev | `cve-2026-21385` | code | 1.0 |
| 6340 | Microsoft Office Security Feature Bypass Vulnerability | kev | `cve-2026-21509` | code | 1.0 |
| 6341 | Microsoft Windows Shell Protection Mechanism Failure Vulnerability | kev | `cve-2026-21510` | code | 1.0 |
| 6342 | Microsoft MSHTML Framework Protection Mechanism Failure Vulnerability | kev | `cve-2026-21513` | code | 1.0 |
| 6343 | Microsoft Office Word Reliance on Untrusted Inputs in a Security Decision Vulnerability | kev | `cve-2026-21514` | code | 1.0 |
| 6344 | Microsoft Windows Type Confusion Vulnerability | kev | `cve-2026-21519` | code | 1.0 |
| 6345 | Microsoft Windows NULL Pointer Dereference Vulnerability | kev | `cve-2026-21525` | code | 1.0 |
| 6346 | Microsoft Windows Improper Privilege Management Vulnerability | kev | `cve-2026-21533` | code | 1.0 |
| 6347 | Atlassian Multiple Products Files or Directories Accessible to External Parties Vulnerability | kev | `cve-2026-21589` | code | 1.0 |
| 6348 | Fortinet FortiClient EMS SQL Injection Vulnerability | kev | `cve-2026-21643` | code | 1.0 |
| 6349 | Oracle HTTP Server and Oracle Weblogic Server Proxy Plug-in Improper Access Control Vulnerability | kev | `cve-2026-21962` | code | 1.0 |
| 6350 | Broadcom VMware Aria Operations Command Injection Vulnerability | kev | `cve-2026-22719` | code | 1.0 |
| 6351 | Dell RecoverPoint for Virtual Machines (RP4VMs) Use of Hard-coded Credentials Vulnerability | kev | `cve-2026-22769` | code | 1.0 |
| 6352 | SmarterTools SmarterMail Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | `cve-2026-23760` | code | 1.0 |
| 6353 | GNU InetUtils Argument Injection Vulnerability | kev | `cve-2026-24061` | code | 1.0 |
| 6354 | Google Chromium CSS Use-After-Free Vulnerability | kev | `cve-2026-2441` | code | 1.0 |
| 6355 | SmarterTools SmarterMail Missing Authentication for Critical Function Vulnerability | kev | `cve-2026-24423` | code | 1.0 |
| 6356 | Fortinet Multiple Products Authentication Bypass Using an Alternate Path or Channel Vulnerability | kev | `cve-2026-24858` | code | 1.0 |
| 6357 | Fortinet FortiSandbox OS Command Injection Vulnerability | kev | `cve-2026-25089` | code | 1.0 |
| 6358 | Soliton Systems K.K FileZen OS Command Injection Vulnerability | kev | `cve-2026-25108` | code | 1.0 |
| 6359 | SolarWinds Serv-U Uncontrolled Resource Consumption Vulnerability | kev | `cve-2026-28318` | code | 1.0 |
| 6360 | Citrix NetScaler Out-of-Bounds Read Vulnerability | kev | `cve-2026-3055` | code | 1.0 |
| 6361 | Linux Kernel Incorrect Resource Transfer Between Spheres Vulnerability | kev | `cve-2026-31431` | code | 1.0 |
| 6362 | Microsoft SharePoint Server Improper Input Validation Vulnerability | kev | `cve-2026-32201` | code | 1.0 |
| 6363 | Microsoft Windows Protection Mechanism Failure Vulnerability | kev | `cve-2026-32202` | code | 1.0 |
| 6364 | Langflow Code Injection Vulnerability | kev | `cve-2026-33017` | code | 1.0 |
| 6365 | Aquasecurity Trivy Embedded Malicious Code Vulnerability | kev | `cve-2026-33634` | code | 1.0 |
| 6366 | Microsoft Internet Key Exchange (IKE) Service Extensions Double Free Vulnerability | kev | `cve-2026-33824` | code | 1.0 |
| 6367 | Microsoft Defender Insufficient Granularity of Access Control Vulnerability | kev | `cve-2026-33825` | code | 1.0 |
| 6368 | Apache ActiveMQ Improper Input Validation Vulnerability | kev | `cve-2026-34197` | code | 1.0 |
| 6369 | Apache Tomcat Missing Encryption of Sensitive Data Vulnerability | kev | `cve-2026-34486` | code | 1.0 |
| 6370 | Adobe Acrobat and Reader Prototype Pollution Vulnerability | kev | `cve-2026-34621` | code | 1.0 |
| 6371 | Ubiquiti UniFi OS Improper Access Control Vulnerability | kev | `cve-2026-34908` | code | 1.0 |
| 6372 | Ubiquiti UniFi OS Path Traversal Vulnerability | kev | `cve-2026-34909` | code | 1.0 |
| 6373 | Ubiquiti UniFi OS Improper Input Validation Vulnerability | kev | `cve-2026-34910` | code | 1.0 |
| 6374 | Trend Micro Apex One (On-Premise) Directory Traversal Vulnerability | kev | `cve-2026-34926` | code | 1.0 |
| 6375 | TrueConf Client Download of Code Without Integrity Check Vulnerability | kev | `cve-2026-3502` | code | 1.0 |
| 6376 | Oracle PeopleSoft Enterprise PeopleTools Missing Authentication for Critical Function Vulnerability | kev | `cve-2026-35273` | code | 1.0 |
| 6377 | Fortinet FortiClient EMS Improper Access Control Vulnerability | kev | `cve-2026-35616` | code | 1.0 |
| 6378 | Google Skia Out-of-Bounds Write Vulnerability | kev | `cve-2026-3909` | code | 1.0 |
| 6379 | Google Chromium V8 Improper Restriction of Operations Within the Bounds of a Memory Buffer Vulnerability | kev | `cve-2026-3910` | code | 1.0 |
| 6380 | Fortinet FortiSandbox OS Command Injection Vulnerability | kev | `cve-2026-39808` | code | 1.0 |
| 6381 | Marimo Remote Code Execution Vulnerability | kev | `cve-2026-39987` | code | 1.0 |
| 6382 | BerriAI LiteLLM Unprotected Alternate Channel Vulnerability | kev | `cve-2026-40217` | code | 1.0 |
| 6383 | Microsoft Defender Link Following Vulnerability | kev | `cve-2026-41091` | code | 1.0 |
| 6384 | Broadcom VMWare ESX Insufficient Logging Vulnerability | kev | `cve-2026-41709` | code | 1.0 |
| 6385 | WebPros cPanel & WHM and WP2 (WordPress Squared) Missing Authentication for Critical Function Vulnerability | kev | `cve-2026-41940` | code | 1.0 |
| 6386 | JFrog Artifactory Incorrect Authorization Vulnerability | kev | `cve-2026-42016` | code | 1.0 |
| 6387 | JFrog Artifactory Improper Authentication Vulnerability | kev | `cve-2026-42018` | code | 1.0 |
| 6388 | BerriAI LiteLLM SQL Injection Vulnerability | kev | `cve-2026-42208` | code | 1.0 |
| 6389 | BerriAI LiteLLM Command Injection Vulnerability | kev | `cve-2026-42271` | code | 1.0 |
| 6390 | Microsoft Exchange Server Cross-Site Scripting Vulnerability | kev | `cve-2026-42897` | code | 1.0 |
| 6391 | Mirasvit Full Page Cache Warmer Deserialization of Untrusted Data Vulnerability | kev | `cve-2026-45247` | code | 1.0 |
| 6392 | TanStack Unspecified Vulnerability | kev | `cve-2026-45321` | code | 1.0 |
| 6393 | Microsoft Defender Denial of Service Vulnerability | kev | `cve-2026-45498` | code | 1.0 |
| 6394 | Microsoft SharePoint Server Deserialization of Untrusted Data Vulnerability | kev | `cve-2026-45659` | code | 1.0 |
| 6395 | Linux Kernel Missing Initialization of a Variable Vulnerability | kev | `cve-2026-46300` | code | 1.0 |
| 6396 | Oracle E-Business Suite Improper Privilege Management Vulnerability | kev | `cve-2026-46817` | code | 1.0 |
| 6397 | Nx Console Embedded Malicious Code Vulnerability | kev | `cve-2026-48027` | code | 1.0 |
| 6398 | LiteSpeed cPanel Plugin Privilege Escalation Vulnerability | kev | `cve-2026-48172` | code | 1.0 |
| 6399 | Adobe ColdFusion Path Traversal Vulnerability | kev | `cve-2026-48282` | code | 1.0 |
| 6400 | Adobe Commerce and Magento Unrestricted Upload of File with Dangerous Type Vulnerability | kev | `cve-2026-48356` | code | 1.0 |
| 6401 | SimpleHelp Authentication Bypass Vulnerability | kev | `cve-2026-48558` | code | 1.0 |
| 6402 | Kludex Starlette HTTP Request/Response Smuggling Vulnerability | kev | `cve-2026-48710` | code | 1.0 |
| 6403 | Widget Factory Joomla Content Editor Improper Access Control Vulnerability | kev | `cve-2026-48907` | code | 1.0 |
| 6404 | JoomShaper SP Page Builder Unrestricted Upload of File with Dangerous Type Vulnerability | kev | `cve-2026-48908` | code | 1.0 |
| 6405 | iCagenda Unrestricted Upload of File with Dangerous Type Vulnerability | kev | `cve-2026-48939` | code | 1.0 |
| 6406 | Kestra OSS OS Command Injection Vulnerability | kev | `cve-2026-49869` | code | 1.0 |
| 6407 | Microsoft SharePoint Deserialization of Untrusted Data Vulnerability | kev | `cve-2026-50522` | code | 1.0 |
| 6408 | Check Point Security Gateway Improper Authentication Vulnerability | kev | `cve-2026-50751` | code | 1.0 |
| 6409 | Google Dawn Use-After-Free Vulnerability | kev | `cve-2026-5281` | code | 1.0 |
| 6410 | Linux Kernel Out-of-Bounds Write Vulnerability | kev | `cve-2026-53266` | code | 1.0 |
| 6411 | Linux Kernel Unspecified Vulnerability | kev | `cve-2026-53362` | code | 1.0 |
| 6412 | WSO2 Multiple Products Path Traversal Vulnerability | kev | `cve-2026-5430` | code | 1.0 |
| 6413 | LiteSpeed cPanel Plugin UNIX Symbolic Link (Symlink) Following Vulnerability | kev | `cve-2026-54420` | code | 1.0 |
| 6414 | Microsoft SharePoint Weak Authentication Vulnerability | kev | `cve-2026-55040` | code | 1.0 |
| 6415 | Langflow Authorization Bypass Through User-Controlled Key Vulnerability | kev | `cve-2026-55255` | code | 1.0 |
| 6416 | Microsoft Active Directory Federation Services Insufficient Granularity of Access Control Vulnerability | kev | `cve-2026-56155` | code | 1.0 |
| 6417 | Microsoft SharePoint Server Missing Authentication for Critical Function Vulnerability | kev | `cve-2026-56164` | code | 1.0 |
| 6418 | Joomlack Page Builder Improper Access Control Vulnerability | kev | `cve-2026-56290` | code | 1.0 |
| 6419 | Balbooa Forms Unrestricted Upload of File with Dangerous Type Vulnerability | kev | `cve-2026-56291` | code | 1.0 |
| 6420 | Microsoft SharePoint Deserialization of Untrusted Data Vulnerability | kev | `cve-2026-58644` | code | 1.0 |
| 6421 | Google Pixel Improper Authorization Vulnerability | kev | `cve-2026-58704` | code | 1.0 |
| 6422 | Broadcom VMware vCenter Path Traversal Vulnerability | kev | `cve-2026-59310` | code | 1.0 |
| 6423 | BerriAI LiteLLM Improper Authentication Vulnerability | kev | `cve-2026-59822` | code | 1.0 |
| 6424 | Gitea Code Injection Vulnerability | kev | `cve-2026-60004` | code | 1.0 |
| 6425 | WordPress Core SQL Injection Vulnerability | kev | `cve-2026-60137` | code | 1.0 |
| 6426 | WordPress Core Interpretation Conflict Vulnerability | kev | `cve-2026-63030` | code | 1.0 |
| 6427 | JetBrains TeamCity Deserialization of Untrusted Data Vulnerability | kev | `cve-2026-63077` | code | 1.0 |
| 6428 | MLflow Server-Side Request Forgery Vulnerability | kev | `cve-2026-64849` | code | 1.0 |
| 6429 | Apple macOS Improper Authentication Vulnerability | kev | `cve-2026-65400` | code | 1.0 |
| 6430 | Microsoft SharePoint Code Injection Vulnerability | kev | `cve-2026-65660` | code | 1.0 |
| 6431 | JFrog Artifactory Improper Limitation of a Pathname to a Restricted Directory Vulnerability | kev | `cve-2026-66384` | code | 1.0 |
| 6432 | MikroTik RouterOS Missing Authentication for Critical Function Vulnerability | kev | `cve-2026-67277` | code | 1.0 |
| 6433 | Mikrotik RouterOS Improper Enforcement of Behavioral Workflow Vulnerability | kev | `cve-2026-67279` | code | 1.0 |
| 6434 | Microsoft Windows Ancillary Function Driver for WinSock Use-After-Free Vulnerability | kev | `cve-2026-68820` | code | 1.0 |
| 6435 | Ivanti Endpoint Manager Mobile (EPMM) Improper Input Validation Vulnerability | kev | `cve-2026-6973` | code | 1.0 |
| 6436 | Adobe Commerce and Magento Incorrect Authorization Vulnerability | kev | `cve-2026-71362` | code | 1.0 |
| 6437 | TrueConf Server Missing Authentication for Critical Function Vulnerability | kev | `cve-2026-72529` | code | 1.0 |
| 6438 | TrueConf Server Code Injection Vulnerability | kev | `cve-2026-72530` | code | 1.0 |
| 6439 | Zyxel GS1900 Series Switches Stack-Based Buffer Overflow Vulnerability | kev | `cve-2026-7273` | code | 1.0 |
| 6440 | Metabase SQL Injection Vulnerability | kev | `cve-2026-72898` | code | 1.0 |
| 6441 | Zimbra Collaboration Suite (ZCS) OS Command Injection Vulnerability | kev | `cve-2026-73570` | code | 1.0 |
| 6442 | Arista Extensible Operating System Incomplete Comparison with Missing Factors Vulnerability | kev | `cve-2026-7473` | code | 1.0 |
| 6443 | Adobe Commerce and Magento Improper Neutralization of Special Elements Used in a Template Engine Vulnerability | kev | `cve-2026-75650` | code | 1.0 |
| 6444 | Cisco Secure Firewall ASA Software and Cisco Secure FTD Software Unchecked Return Value to NULL Pointer Dereference Vulnerability | kev | `cve-2026-76418` | code | 1.0 |
| 6445 | Cisco Identity Services Engine Incorrect Use of Privileged APIs Vulnerability | kev | `cve-2026-76460` | code | 1.0 |
| 6446 | Cisco Secure Email Gateway SQL Injection Vulnerability | kev | `cve-2026-76461` | code | 1.0 |
| 6447 | Cisco Expressway Series Use of Hard-coded Credentials Vulnerability | kev | `cve-2026-76466` | code | 1.0 |
| 6448 | Cisco Catalyst SD-WAN Manager Hex Encoding Vulnerability | kev | `cve-2026-76504` | code | 1.0 |
| 6449 | Cisco Secure Email Gateway Relative Path Traversal Vulnerability | kev | `cve-2026-76519` | code | 1.0 |
| 6450 | Progress LoadMaster Command Injection Vulnerability | kev | `cve-2026-8037` | code | 1.0 |
| 6451 | PaperCut NG/MF Missing Authentication for Critical Function Vulnerability | kev | `cve-2026-81578` | code | 1.0 |
| 6452 | Microsoft Windows Link Following Vulnerability | kev | `cve-2026-81963` | code | 1.0 |
| 6453 | PaperCut NG/MF Unsafe Reflection Vulnerability | kev | `cve-2026-82078` | code | 1.0 |
| 6454 | JFrog Artifactory Improper Authentication Vulnerability | kev | `cve-2026-82329` | code | 1.0 |
| 6455 | SonicWall SMA1000 Appliances Server-Side Request Forgery Vulnerability | kev | `cve-2026-83548` | code | 1.0 |
| 6456 | SonicWall SMA1000 Appliances OS Command Injection Vulnerability | kev | `cve-2026-83549` | code | 1.0 |
| 6457 | Daemon Tools Lite Embedded Malicious Code Vulnerability | kev | `cve-2026-8398` | code | 1.0 |
| 6458 | Citrix NetScaler ADC and NetScaler Gateway Improper Restriction of Operations within the Bounds of a Memory Buffer Vulnerability | kev | `cve-2026-8452` | code | 1.0 |
| 6459 | ConnectWise ScreenConnect Improper Privilege Management and Missing Authorization Vulnerability | kev | `cve-2026-84869` | code | 1.0 |
| 6460 | Google Chromium V8 Type Confusion Vulnerability | kev | `cve-2026-85046` | code | 1.0 |
| 6461 | Check Point Multiple Products Improper Certificate Validation Vulnerability | kev | `cve-2026-85102` | code | 1.0 |
| 6462 | GitLab Community Edition and Enterprise Edition Path Traversal Vulnerability | kev | `cve-2026-85706` | code | 1.0 |
| 6463 | Microsoft Windows Heap-Based Buffer Overflow Vulnerability | kev | `cve-2026-85880` | code | 1.0 |
| 6464 | MikroTik RouterOS Improper Neutralization of Argument Delimiters in a Command Vulnerability | kev | `cve-2026-86060` | code | 1.0 |
| 6465 | N-able N-central Static Code Injection Vulnerability | kev | `cve-2026-86218` | code | 1.0 |
| 6466 | Apple Multiple Products Out-of-Bounds Write Vulnerability | kev | `cve-2026-86950` | code | 1.0 |
| 6467 | Google Chromium V8 Out of Bounds Write Vulnerability | kev | `cve-2026-87491` | code | 1.0 |
| 6468 | Acronis Backup Incorrect Default Permissions Vulnerability | kev | `cve-2026-87886` | code | 1.0 |
| 6469 | WordPress Core Remote File Inclusion Vulnerability | kev | `cve-2026-87902` | code | 1.0 |
| 6470 | Citrix NetScaler Improper Input Validation Vulnerability | kev | `cve-2026-88771` | code | 1.0 |
| 6471 | Citrix NetScaler Improper Restriction of Operations within the Bounds of a Memory Buffer Vulnerability | kev | `cve-2026-88772` | code | 1.0 |
| 6472 | Citrix NetScaler Improper Restriction of Operations within the Bounds of a Memory Buffer Vulnerability | kev | `cve-2026-88779` | code | 1.0 |
| 6473 | Apache Axis Unsafe Reflection Vulnerability | kev | `cve-2026-90436` | code | 1.0 |
| 6474 | Vanderbilt University Research Electronic Data Capture (REDCap) Code Injection Vulnerability | kev | `cve-2026-90817` | code | 1.0 |
| 6475 | Drupal Core SQL Injection Vulnerability | kev | `cve-2026-9082` | code | 1.0 |
| 6476 | IBM Langflow Code Injection Vulnerability | kev | `cve-2026-9198` | code | 1.0 |
| 6477 | Check Point Multiple Products Path Traversal Vulnerability | kev | `cve-2026-93616` | code | 1.0 |
| 6478 | Arista VeloCloud Orchestrator Improper Input Validation Vulnerability | kev | `cve-2026-93952` | code | 1.0 |
| 6479 | F5 BIG-IP APM Heap-based Buffer Overflow Vulnerability | kev | `cve-2026-94127` | code | 1.0 |
| 6480 | Sangoma Switchvox SQL Injection Vulnerability | kev | `cve-2026-9586` | code | 1.0 |
| 6481 | Plex Media Server Server-Side Request Forgery (SSRF) Vulnerability | kev | `cve-2026-96652` | code | 1.0 |
| 6482 | Plex Media Server Improper Neutralization of Encoded URI Schemes in a Web Page Vulnerability | kev | `cve-2026-96654` | code | 1.0 |
| 6483 | Plex Media Server Server-Side Request Forgery (SSRF) Vulnerability | kev | `cve-2026-96655` | code | 1.0 |
| 6484 | Plex Media Server External Control of File Name or Path Vulnerability | kev | `cve-2026-96656` | code | 1.0 |
| 6485 | Microsoft Windows Command Injection Vulnerability | kev | `cve-2026-96951` | code | 1.0 |
| 6486 | GitLab Community Edition and Enterprise Edition Integer Overflow or Wraparound Vulnerability | kev | `cve-2026-97229` | code | 1.0 |
| 6487 | AdBlock before 2.21 allows remote attackers to block arbitrary resources on arbitrary websites and to disable arbitrary blocking filters | nvd | `cve-2015-2692` | code | 1.0 |
| 6488 | The Data Protection component in the VMware vSphere GUI in IBM Tivoli Storage Manager for Virtual Environments: Data Protection for VMware (aka Spectrum Protect | nvd | `cve-2015-7425` | code | 1.0 |
| 6489 | The Data Protection extension in the VMware GUI in IBM Tivoli Storage Manager for Virtual Environments: Data Protection for VMware (aka Spectrum Protect for Vir | nvd | `cve-2015-7426` | code | 1.0 |
| 6490 | The initialize method in the Histogram class in lib/colorscore/histogram.rb in the colorscore gem before 0.0.5 for Ruby allows context-dependent attackers to ex | nvd | `cve-2015-7541` | code | 1.0 |
| 6491 | SearchBlox 8.3 before 8.3.1 allows remote attackers to write to the config file, and consequently cause a denial of service (application crash), via unspecified | nvd | `cve-2015-7919` | code | 1.0 |
| 6492 | Adcon Telemetry A840 Telemetry Gateway Base Station has hardcoded credentials, which allows remote attackers to obtain administrative access via unspecified vec | nvd | `cve-2015-7930` | code | 1.0 |
| 6493 | The PasswordReset.Controllers.ResetController.ChangePasswordIndex method in PasswordReset.dll in Dovestones AD Self Password Reset before 3.0.4 allows remote  | nvd | `cve-2015-8267` | code | 1.0 |
| 6494 | Integer overflow in the ImageRegionReader::ReadIntoBuffer function in MediaStorageAndFileFormat/gdcmImageRegionReader.cxx in Grassroots DICOM (aka GDCM) before  | nvd | `cve-2015-8396` | code | 1.0 |
| 6495 | Adobe Flash Player before 18.0.0 and 19.x and 20.x before 20.0.0 on Windows and OS X and before 11.2.202 on Linux, Adobe AIR before 20.0.0, Adob | nvd | `cve-2015-8459` | code | 1.0 |
| 6496 | Local privilege escalation vulnerability in the Gentoo QEMU package before 2.5.0-r1 | nvd | `cve-2015-8556` | code | 1.0 |
| 6497 | The idle stream handling in nghttp2 before 1.6.0 allows attackers to have unspecified impact via unknown vectors, aka a heap-use-after-free bug | nvd | `cve-2015-8659` | code | 1.0 |
| 6498 | The multifilesystem storage backend in Radicale before 1.1 allows remote attackers to read or write to arbitrary files via a crafted component name | nvd | `cve-2015-8747` | code | 1.0 |
| 6499 | SQL injection vulnerability in the Group Promotions module in the admin control panel in MyBB (aka MyBulletinBoard) before 1.6.18 and 1.8.x before 1.8.6 and MyB | nvd | `cve-2015-8974` | code | 1.0 |
| 6500 | MailEnable before 8.60 allows XXE via an XML document in the request.aspx Options parameter | nvd | `cve-2015-9280` | code | 1.0 |
| 6501 | MySQL for PCF tiles 1.7.x before 1.7.10 were discovered to log the AWS access key in plaintext. These credentials were logged to the Service Backup component lo | nvd | `cve-2016-0898` | code | 1.0 |
| 6502 | An issue was discovered in Radisys MRF Web Panel (SWMS) 9.0.1. The MSM_MACRO_NAME POST parameter in /swms/ms.cgi was discovered to be vulnerable to OS command i | nvd | `cve-2016-10043` | code | 1.0 |
| 6503 | Adobe Reader and Acrobat before 11.0.16, Acrobat and Acrobat Reader DC Classic before 15.006.30172, and Acrobat and Acrobat Reader DC Continuous before 15.016.2 | nvd | `cve-2016-1038` | code | 1.0 |
| 6504 | Adobe Reader and Acrobat before 11.0.16, Acrobat and Acrobat Reader DC Classic before 15.006.30172, and Acrobat and Acrobat Reader DC Continuous before 15.016.2 | nvd | `cve-2016-1041` | code | 1.0 |
| 6505 | Adobe Reader and Acrobat before 11.0.16, Acrobat and Acrobat Reader DC Classic before 15.006.30172, and Acrobat and Acrobat Reader DC Continuous before 15.016.2 | nvd | `cve-2016-1044` | code | 1.0 |
| 6506 | The XML parser in Cisco Information Server (CIS) 6.2 allows remote attackers to read arbitrary files or cause a denial of service (CPU and memory consumption) v | nvd | `cve-2016-1343` | code | 1.0 |
| 6507 | The filesystem storage backend in Radicale before 1.1 on Windows allows remote attackers to read or write to arbitrary files via a crafted path, as demonstrated | nvd | `cve-2016-1505` | code | 1.0 |
| 6508 | Multiple unspecified vulnerabilities in the browser engine in Mozilla Firefox before 44.0 allow remote attackers to cause a denial of service (memory corruption | nvd | `cve-2016-1931` | code | 1.0 |
| 6509 | HPE Operations Manager 8.x and 9.0 on Windows allows remote attackers to execute arbitrary commands via a crafted serialized Java object, related to the Apache  | nvd | `cve-2016-1985` | code | 1.0 |
| 6510 | Pulse Connect Secure (PCS) 8.2 before 8.2r1, 8.1 before 8.1r2, 8.0 before 8.0r10, and 7.4 before 7.4r13.4 allow remote attackers to read sensitive system authen | nvd | `cve-2016-4787` | code | 1.0 |
| 6511 | General Electric (GE) Bently Nevada 3500/22M USB with firmware before 5.0 and Bently Nevada 3500/22M Serial have open ports, which makes it easier for remote at | nvd | `cve-2016-5788` | code | 1.0 |
| 6512 | IBM BigFix Platform could allow a remote attacker to execute arbitrary code on the system, caused by a use-after-free race condition. An attacker could exploit  | nvd | `cve-2016-6082` | code | 1.0 |
| 6513 | VMware vRealize Operations (aka vROps) 6.x before 6.4.0 allows remote authenticated users to gain privileges, or halt and remove virtual machines, via unspecifi | nvd | `cve-2016-7457` | code | 1.0 |
| 6514 | SQL injection vulnerability in core services in Intel Security McAfee ePolicy Orchestrator (ePO) 5.3.2 and earlier and 5.1.3 and earlier allows attackers to alt | nvd | `cve-2016-8027` | code | 1.0 |
| 6515 | An issue was discovered in Schneider Electric ConneXium firewalls TCSEFEC23F3F20 all versions, TCSEFEC23F3F21 all versions, TCSEFEC23FCF20 all versions, TCSEFEC | nvd | `cve-2016-8352` | code | 1.0 |
| 6516 | An issue was discovered in Moxa OnCell OnCellG3470A-LTE, AWK-1131A/3131A/4131A Series, AWK-3191 Series, AWK-5232/6232 Series, AWK-1121/1127 Series, WAC-1001 V2  | nvd | `cve-2016-8363` | code | 1.0 |
| 6517 | IBM UrbanCode Deploy could allow a user to execute code using a specially crafted file upload that would replace code on the server. This code could be executed | nvd | `cve-2016-8938` | code | 1.0 |
| 6518 | A hard-coded cryptographic key vulnerability was identified in Red Lion Controls Sixnet-Managed Industrial Switches running firmware Version 5.0.196 and Stride- | nvd | `cve-2016-9335` | code | 1.0 |
| 6519 | An issue was discovered in Rockwell Automation Logix5000 Programmable Automation Controller FRN 16.00 through 21.00 (excluding all firmware versions prior to FR | nvd | `cve-2016-9343` | code | 1.0 |
| 6520 | Vulnerability in the Oracle WebLogic Server component of Oracle Fusion Middleware (subcomponent: JNDI). Supported versions that are affected are 10.3.6.0 and 12 | nvd | `cve-2017-10137` | code | 1.0 |
| 6521 | Vulnerability in the Oracle Identity Manager component of Oracle Fusion Middleware (subcomponent: Default Account). Supported versions that are affected are 11 | nvd | `cve-2017-10151` | code | 1.0 |
| 6522 | Vulnerability in the Oracle Tuxedo component of Oracle Fusion Middleware (subcomponent: Core). Supported versions that are affected are 11.1.1, 12.1.1, 12.1.3 a | nvd | `cve-2017-10269` | code | 1.0 |
| 6523 | Vulnerability in the Oracle Hospitality Reporting and Analytics component of Oracle Hospitality Applications (subcomponent: Report). Supported versions that are | nvd | `cve-2017-10402` | code | 1.0 |
| 6524 | Vulnerability in the Oracle Hospitality Reporting and Analytics component of Oracle Hospitality Applications (subcomponent: Report). Supported versions that are | nvd | `cve-2017-10405` | code | 1.0 |
| 6525 | Xen through 4.8.x mishandles page transfer, which allows guest OS users to obtain privileged host OS access, aka XSA-217 | nvd | `cve-2017-10912` | code | 1.0 |
| 6526 | Xen through 4.8.x does not validate memory allocations during certain P2M operations, which allows guest OS users to obtain privileged host OS access, aka XSA-2 | nvd | `cve-2017-10918` | code | 1.0 |
| 6527 | The grant-table feature in Xen through 4.8.x mishandles a GNTMAP_device_map and GNTMAP_host_map mapping, when followed by only a GNTMAP_host_map unmapping, whic | nvd | `cve-2017-10920` | code | 1.0 |
| 6528 | The grant-table feature in Xen through 4.8.x does not ensure sufficient type counts for a GNTMAP_device_map and GNTMAP_host_map mapping, which allows guest OS u | nvd | `cve-2017-10921` | code | 1.0 |
| 6529 | An issue was discovered in Adobe Connect 9.6.2 and earlier versions. A Server-Side Request Forgery (SSRF) vulnerability exists that could be abused to bypass ne | nvd | `cve-2017-11291` | code | 1.0 |
| 6530 | An exploitable heap overflow vulnerability exists in the tinysvcmdns library version 2016-07-18. A specially crafted packet can make the library overwrite an ar | nvd | `cve-2017-12087` | code | 1.0 |
| 6531 | A authentication bypass and execution of code vulnerability in HPE Integrated Lights-out 4 (iLO 4) version prior to 2.53 was found | nvd | `cve-2017-12542` | code | 1.0 |
| 6532 | Analysis of the Bomgar Remote Support Portal JavaStart.jar Applet 52790 and earlier revealed that it is vulnerable to a path traversal vulnerability. The archiv | nvd | `cve-2017-12815` | code | 1.0 |
| 6533 | Server Side Request Forgery vulnerability in Vebto Pixie Image Editor 1.4 and 1.7 allows remote attackers to disclose information or execute arbitrary code via  | nvd | `cve-2017-12905` | code | 1.0 |
| 6534 | An Improper Authentication issue was discovered in iniNet Solutions iniNet Webserver, all versions prior to V2.02.0100. The webserver does not properly authenti | nvd | `cve-2017-13995` | code | 1.0 |
| 6535 | EMC RSA Authentication Agent API 8.5 for C and RSA Authentication Agent SDK 8.6 for C allow attackers to bypass authentication, aka an "Error Handling Vulnerabi | nvd | `cve-2017-14378` | code | 1.0 |
| 6536 | An exploitable OS Command Injection vulnerability exists in the Telnet, SSH, and console login functionality of Moxa AWK-3131A Industrial IEEE 802.11a/b/g/n wir | nvd | `cve-2017-14459` | code | 1.0 |
| 6537 | An exploitable access control vulnerability exists in the data, program, and function file permissions functionality of Allen Bradley Micrologix 1400 Series B F | nvd | `cve-2017-14469` | code | 1.0 |
| 6538 | An exploitable access control vulnerability exists in the data, program, and function file permissions functionality of Allen Bradley Micrologix 1400 Series B F | nvd | `cve-2017-14470` | code | 1.0 |
| 6539 | An exploitable access control vulnerability exists in the data, program, and function file permissions functionality of Allen Bradley Micrologix 1400 Series B F | nvd | `cve-2017-14471` | code | 1.0 |
| 6540 | An exploitable access control vulnerability exists in the data, program, and function file permissions functionality of Allen Bradley Micrologix 1400 Series B F | nvd | `cve-2017-14472` | code | 1.0 |
| 6541 | An exploitable access control vulnerability exists in the data, program, and function file permissions functionality of Allen Bradley Micrologix 1400 Series B F | nvd | `cve-2017-14473` | code | 1.0 |
| 6542 | The safe-eval module describes itself as a safer version of eval. By accessing the object constructors, un-sanitized user input can access the entire standard l | nvd | `cve-2017-16088` | code | 1.0 |
| 6543 | A Buffer Overflow issue was discovered in Rockwell Automation Allen-Bradley MicroLogix 1400 Controllers, Series B and C Versions 21.002 and earlier. The stack-b | nvd | `cve-2017-16740` | code | 1.0 |
| 6544 | hw/input/ps2.c in Qemu does not validate 'rptr' and 'count' values during guest migration, leading to out-of-bounds access | nvd | `cve-2017-16845` | code | 1.0 |
| 6545 | A vulnerability in Juniper Networks NorthStar Controller Application prior to version 2.1.0 Service Pack 1 may allow an unauthenticated, unprivileged, network-b | nvd | `cve-2017-2320` | code | 1.0 |
| 6546 | The Integrated User Firewall (UserFW) feature was introduced in Junos OS version 12.1X47-D10 on the Juniper SRX Series devices to provide simple integration of  | nvd | `cve-2017-2343` | code | 1.0 |
| 6547 | An exploitable buffer overflow exists in the psnotifyd application of the Pharos PopUp printer client version 9.0. A specially crafted packet can be sent to the | nvd | `cve-2017-2785` | code | 1.0 |
| 6548 | A buffer overflows exists in the psnotifyd application of the Pharos PopUp printer client version 9.0. A specially crafted packet can be sent to the victim's co | nvd | `cve-2017-2788` | code | 1.0 |
| 6549 | Adobe Digital Editions versions 4.5.4 and earlier have an exploitable memory corruption vulnerability in the PDF runtime engine. Successful exploitation could l | nvd | `cve-2017-3088` | code | 1.0 |
| 6550 | Vulnerability in the Primavera P6 Enterprise Project Portfolio Management component of Oracle Primavera Products Suite (subcomponent: Web Access). Supported ver | nvd | `cve-2017-3324` | code | 1.0 |
| 6551 | Vulnerability in the Solaris component of Oracle Sun Systems Products Suite (subcomponent: Kernel RPC). For supported versions that are affected see note. Easil | nvd | `cve-2017-3623` | code | 1.0 |
| 6552 | A vulnerability in the web-based GUI of Cisco Prime Home could allow an unauthenticated, remote attacker to bypass authentication and execute actions with admin | nvd | `cve-2017-3791` | code | 1.0 |
| 6553 | An issue was discovered in Carlo Gavazzi VMU-C EM prior to firmware Version A11_U05, and VMU-C PV prior to firmware Version A17. Successful exploitation of this | nvd | `cve-2017-5145` | code | 1.0 |
| 6554 | When executing a program via the bubblewrap sandbox, the nonpriv session can escape to the parent session by using the TIOCSTI ioctl to push characters into the | nvd | `cve-2017-5226` | code | 1.0 |
| 6555 | The Symantec Messaging Gateway can encounter an issue of remote code execution, which describes a situation whereby an individual may obtain the ability to exec | nvd | `cve-2017-6326` | code | 1.0 |
| 6556 | Zoho ManageEngine Desktop Central before build 100082 allows remote attackers to obtain control over all connected active desktops via unspecified vectors | nvd | `cve-2017-7213` | code | 1.0 |
| 6557 | Uploaded XML documents were not correctly validated in Apache OpenMeetings 3.1.0 | nvd | `cve-2017-7664` | code | 1.0 |
| 6558 | In SolarWinds Log & Event Manager (LEM) before 6.3.1 Hotfix 4, a menu system is encountered when the SSH service is accessed with "cmc" and "password" (the defa | nvd | `cve-2017-7722` | code | 1.0 |
| 6559 | An Improper Access Control issue was discovered in Schweitzer Engineering Laboratories (SEL) SEL-3620 and SEL-3622 Security Gateway Versions R202 and, R203, R20 | nvd | `cve-2017-7928` | code | 1.0 |
| 6560 | Zyxel WRE6505 devices have a default TELNET password of 1234 for the root and admin accounts, which makes it easier for remote attackers to conduct DNS hijackin | nvd | `cve-2017-7964` | code | 1.0 |
| 6561 | www.modified-shop.org modified eCommerce Shopsoftware 2.0.2 rev 10690 has XXE in api/it-recht-kanzlei/api-it-recht-kanzlei.php | nvd | `cve-2017-8110` | code | 1.0 |
| 6562 | An issue was discovered on Accellion FTA devices before FTA_9_12_180. Because a regular expression (intended to match local https URLs) lacks an initial ^ chara | nvd | `cve-2017-8794` | code | 1.0 |
| 6563 | A vulnerability in the Secure Sockets Layer (SSL) VPN functionality of the Cisco Adaptive Security Appliance (ASA) Software could allow an unauthenticated, remo | nvd | `cve-2018-0101` | code | 1.0 |
| 6564 | A vulnerability in Cisco Digital Network Architecture (DNA) Center could allow an unauthenticated, remote attacker to log in to an affected system by using an a | nvd | `cve-2018-0222` | code | 1.0 |
| 6565 | A vulnerability in the container management subsystem of Cisco Digital Network Architecture (DNA) Center could allow an unauthenticated, remote attacker to bypa | nvd | `cve-2018-0268` | code | 1.0 |
| 6566 | I Librarian I-librarian version 4.8 and earlier contains a XML External Entity (XXE) vulnerability in line 154 of importmetadata.php(simplexml_load_string) that | nvd | `cve-2018-1000124` | code | 1.0 |
| 6567 | Eclipse RDF4j version < 2.4.0 Milestone 2 contains a XML External Entity (XXE) vulnerability in RDF4j XML parser parsing RDF files that can result in the disclo | nvd | `cve-2018-1000644` | code | 1.0 |
| 6568 | Stroom version <5.4.5 contains a XML External Entity (XXE) vulnerability in XML Parser that can result in disclosure of confidential data, denial of service, se | nvd | `cve-2018-1000651` | code | 1.0 |
| 6569 | JabRef version <=4.3.1 contains a XML External Entity (XXE) vulnerability in MsBibImporter XML Parser that can result in disclosure of confidential data, denial | nvd | `cve-2018-1000652` | code | 1.0 |
| 6570 | neo4j-contrib neo4j-apoc-procedures version before commit 45bc09c contains a XML External Entity (XXE) vulnerability in XML Parser that can result in Disclosure | nvd | `cve-2018-1000820` | code | 1.0 |
| 6571 | MicroMathematics version before commit 5c05ac8 contains a XML External Entity (XXE) vulnerability in SMathStudio files that can result in Disclosure of confiden | nvd | `cve-2018-1000821` | code | 1.0 |
| 6572 | codelibs fess version before commit faa265b contains a XML External Entity (XXE) vulnerability in GSA XML file parser that can result in Disclosure of confident | nvd | `cve-2018-1000822` | code | 1.0 |
| 6573 | exist version <= 5.0.0-RC4 contains a XML External Entity (XXE) vulnerability in XML Parser for REST Server that can result in Disclosure of confidential data,  | nvd | `cve-2018-1000823` | code | 1.0 |
| 6574 | FreeCol version <= nightly-2018-08-22 contains a XML External Entity (XXE) vulnerability in FreeColXMLReader parser that can result in Disclosure of confidentia | nvd | `cve-2018-1000825` | code | 1.0 |
| 6575 | XR3Player version <= V3.124 contains a XML External Entity (XXE) vulnerability in Playlist parser that can result in Disclosure of confidential data, denial of  | nvd | `cve-2018-1000830` | code | 1.0 |
| 6576 | K9Mail version <= v5.600 contains a XML External Entity (XXE) vulnerability in WebDAV response parser that can result in Disclosure of confidential data, denial | nvd | `cve-2018-1000831` | code | 1.0 |
| 6577 | KeePassDX version <= 2.5.0.0beta17 contains a XML External Entity (XXE) vulnerability in kdbx file parser that can result in Disclosure of confidential data, de | nvd | `cve-2018-1000835` | code | 1.0 |
| 6578 | UML Designer version <= 8.0.0 contains a XML External Entity (XXE) vulnerability in XML parser for plugins that can result in Disclosure of confidential data, d | nvd | `cve-2018-1000837` | code | 1.0 |
| 6579 | autopsy version <= 4.9.0 contains a XML External Entity (XXE) vulnerability in CaseMetadata XML Parser that can result in Disclosure of confidential data, denia | nvd | `cve-2018-1000838` | code | 1.0 |
| 6580 | A vulnerability in Trend Micro Control Manager (versions 6.0 and 7.0) could allow an attacker to conduct a server-side request forgery (SSRF) attack on vulnerab | nvd | `cve-2018-10511` | code | 1.0 |
| 6581 | Stack-based buffer overflow in Activision Infinity Ward Call of Duty Modern Warfare 2 before 2018-04-26 allows remote attackers to execute arbitrary code via cr | nvd | `cve-2018-10718` | code | 1.0 |
| 6582 | A SQL injection vulnerability in the web administration and quarantine components of Micro Focus Secure Messaging Gateway allows an unauthenticated remote attac | nvd | `cve-2018-12464` | code | 1.0 |
| 6583 | A vulnerability has been identified in TIM 1531 IRC (All version < V2.0). The devices was missing proper authentication on port 102/tcp, although configured. Su | nvd | `cve-2018-13816` | code | 1.0 |
| 6584 | FasterXML jackson-databind 2.x before 2.9.7 might allow remote attackers to conduct server-side request forgery (SSRF) attacks by leveraging failure to block th | nvd | `cve-2018-14721` | code | 1.0 |
| 6585 | A command injection vulnerability in the apex-publish-static-files npm module version <2.0.1 which allows arbitrary shell command execution through a maliciousl | nvd | `cve-2018-16462` | code | 1.0 |
| 6586 | In FreeBSD before 11.2-STABLE(r341486) and 11.2-RELEASE-p6, insufficient bounds checking in one of the device models provided by bhyve can permit a guest operat | nvd | `cve-2018-17160` | code | 1.0 |
| 6587 | IBM Security Access Manager Appliance 9.0.4 and 9.0.5 could allow remote code execution when Advanced Access Control or Federation services are running. IBM | nvd | `cve-2018-1722` | code | 1.0 |
| 6588 | An earlier fix for an Inter-process Communication (IPC) vulnerability, CVE-2011-3079, added authentication to communication between IPC endpoints and server par | nvd | `cve-2018-18505` | code | 1.0 |
| 6589 | Sandboxie 5.26 allows a Sandbox Escape via an "import os" statement, followed by os.system("cmd") or os.system("powershell"), within a .py file. NOTE: the vendo | nvd | `cve-2018-18748` | code | 1.0 |
| 6590 | The Kubernetes integration in GitLab Enterprise Edition 11.x before 11.2.8, 11.3.x before 11.3.9, and 11.4.x before 11.4.4 has SSRF | nvd | `cve-2018-18843` | code | 1.0 |
| 6591 | mPDF through 7.1.6, if deployed as a web application that accepts arbitrary HTML, allows SSRF, as demonstrated by a '<img src="http://192.168' substring that tr | nvd | `cve-2018-19047` | code | 1.0 |
| 6592 | An issue was discovered in the MQTT server in Contiki-NG before 4.2. The function parse_publish_vhdr() that parses MQTT PUBLISH messages with a variable length  | nvd | `cve-2018-19417` | code | 1.0 |
| 6593 | Vulnerability in the Sun ZFS Storage Appliance Kit (AK) component of Oracle Sun Systems Products Suite (subcomponent: Core Services). The supported version that | nvd | `cve-2018-2611` | code | 1.0 |
| 6594 | Vulnerability in the Oracle GoldenGate component of Oracle GoldenGate (subcomponent: Monitoring Manager). Supported versions that are affected are 12.1.2, 1 | nvd | `cve-2018-2913` | code | 1.0 |
| 6595 | Incorrect parsing in url-parse <1.4.3 returns wrong hostname which leads to multiple vulnerabilities such as SSRF, Open Redirect, Bypass Authentication Protocol | nvd | `cve-2018-3774` | code | 1.0 |
| 6596 | An exploitable vulnerability exists in the REST parser of video-core's HTTP server of the Samsung SmartThings Hub STH-ETH-250 - Firmware version 0.20.17. The vi | nvd | `cve-2018-3907` | code | 1.0 |
| 6597 | An exploitable heap overflow vulnerability exists in the WkbProgramLow function of WibuKey Network server management, version 6.40.2402.500. A specially crafted | nvd | `cve-2018-3991` | code | 1.0 |
| 6598 | An issue was discovered in certain Apple products. macOS before 10.13.3 is affected. The issue involves the "Sandbox" component. It allows bypass of a sandbox p | nvd | `cve-2018-4091` | code | 1.0 |
| 6599 | An issue was discovered in certain Apple products. macOS before 10.13.5 is affected. The issue involves the "Grand Central Dispatch" component. It allows attack | nvd | `cve-2018-4229` | code | 1.0 |
| 6600 | An issue was discovered in Adobe Acrobat Reader 2018.009.20050 and earlier versions, 2017.011.30070 and earlier versions, 2015.006.30394 and earlier versions. T | nvd | `cve-2018-4872` | code | 1.0 |
| 6601 | Authentication Bypass vulnerability in the administrative user interface in McAfee Web Gateway 7.8.1 through 7.8.1 allows remote attackers to execute arbitr | nvd | `cve-2018-6667` | code | 1.0 |
| 6602 | Stack-based Buffer Overflow vulnerability in libUPnPHndlr.so in Belkin Wemo Insight Smart Plug allows remote attackers to bypass local security protection via a | nvd | `cve-2018-6692` | code | 1.0 |
| 6603 | The VMware AirWatch Agent for Android prior to 8.2 and AirWatch Agent for Windows Mobile prior to 6.5.2 contain a remote code execution vulnerability in real ti | nvd | `cve-2018-6968` | code | 1.0 |
| 6604 | Juniper ATP ships with hard coded credentials in the Web Collector instance which gives an attacker the ability to take full control of any installation of the  | nvd | `cve-2019-0020` | code | 1.0 |
| 6605 | Juniper ATP ships with hard coded credentials in the Cyphort Core instance which gives an attacker the ability to take full control of any installation of the s | nvd | `cve-2019-0022` | code | 1.0 |
| 6606 | Zoho ManageEngine ADSelfService Plus 5.x before build 5703 has SSRF | nvd | `cve-2019-3905` | code | 1.0 |
| 6607 | Multiple use-after-free vulnerabilities in vnc.c in the VNC server in QEMU 0.10.6 and earlier might allow guest OS users to execute arbitrary code on the host O | nvd | `cve-2009-3616` | code | 1.0 |
| 6608 | The VMX process in VMware ESXi 3.5 through 4.1 and ESX 3.5 through 4.1 does not properly handle RPC commands, which allows guest OS users to cause a denial of s | nvd | `cve-2012-1516` | code | 1.0 |
| 6609 | The portal client in IBM Tivoli Monitoring (ITM) 6.2.2 through FP9, 6.2.3 through FP5, and 6.3.0 through FP6 allows remote authenticated users to gain privilege | nvd | `cve-2015-7411` | code | 1.0 |
| 6610 | eWON devices with firmware before 10.1s0 omit RBAC for I/O server information and status requests, which allows remote attackers to obtain sensitive information | nvd | `cve-2015-7926` | code | 1.0 |
| 6611 | The GMS ViewPoint (GMSVP) web application in Dell SonicWALL GMS, Analyzer, and UMA EM5000 7.2, 8.0, and 8.1 before Hotfix 168056 allows remote authenticated use | nvd | `cve-2016-2396` | code | 1.0 |
| 6612 | lshell 0.9.16 allows remote authenticated users to break out of a limited shell and execute arbitrary commands | nvd | `cve-2016-6902` | code | 1.0 |
| 6613 | lshell 0.9.16 allows remote authenticated users to break out of a limited shell and execute arbitrary commands | nvd | `cve-2016-6903` | code | 1.0 |
| 6614 | An issue was discovered in Smiths-Medical CADD-Solis Medication Safety Software, Version 1.0; 2.0; 3.0; and 3.1. CADD-Solis Medication Safety Software grants an | nvd | `cve-2016-8355` | code | 1.0 |
| 6615 | Remote Command Execution in com.trend.iwss.gui.servlet.ManagePatches in Trend Micro Interscan Web Security Virtual Appliance (IWSVA) version 6.5-SP2_Build_Linux | nvd | `cve-2016-9269` | code | 1.0 |
| 6616 | PricewaterhouseCoopers (PwC) ACE-ABAP 8.10.304 for SAP Security allows remote authenticated users to conduct ABAP injection attacks and execute arbitrary code v | nvd | `cve-2016-9832` | code | 1.0 |
| 6617 | Vulnerability in the OJVM component of Oracle Database Server. Supported versions that are affected are 11.2.0, 12.1.0 and 12.2.0. Easily exploitable vuln | nvd | `cve-2017-10202` | code | 1.0 |
| 6618 | Vulnerability in the Oracle Tuxedo component of Oracle Fusion Middleware (subcomponent: Core). Supported versions that are affected are 11.1.1, 12.1.1, 12.1.3 a | nvd | `cve-2017-10272` | code | 1.0 |
| 6619 | Vulnerability in the Oracle WebLogic Server component of Oracle Fusion Middleware (subcomponent: WLS - Web Services). The supported version that is affected are | nvd | `cve-2017-10352` | code | 1.0 |
| 6620 | Vulnerability in the Oracle Hospitality Cruise AffairWhere component of Oracle Hospitality Applications (subcomponent: AffairWhere). Supported versions that are | nvd | `cve-2017-10396` | code | 1.0 |
| 6621 | Vulnerability in the Oracle Hospitality Reporting and Analytics component of Oracle Hospitality Applications (subcomponent: iQuery). Supported versions that are | nvd | `cve-2017-10404` | code | 1.0 |
| 6622 | A vulnerability in the web console of the Cisco Cloud Services Platform (CSP) 2100 could allow an authenticated, remote attacker to interact maliciously with th | nvd | `cve-2017-12251` | code | 1.0 |
| 6623 | IBM Security Guardium 10.0 could allow a remote authenticated attacker to execute arbitrary commands on the system. By sending a specially-crafted request, an a | nvd | `cve-2017-1253` | code | 1.0 |
| 6624 | Remote enabling and disabling admin interface in Gemalto's HASP SRM, Sentinel HASP and Sentinel LDK products prior to Sentinel LDK RTE version 7.55 leads to new | nvd | `cve-2017-12822` | code | 1.0 |
| 6625 | XML external entity (XXE) vulnerability in the import package functionality of the deployment module in Lansweeper before 6.0.100 allows remote authenticated | nvd | `cve-2017-13706` | code | 1.0 |
| 6626 | An exploitable buffer overflow vulnerability exists in Insteon Hub running firmware version 1012. The HTTP server implementation incorrectly handles the URL par | nvd | `cve-2017-14444` | code | 1.0 |
| 6627 | An exploitable buffer overflow vulnerability exists in Insteon Hub running firmware version 1012. The HTTP server implementation incorrectly handles the host pa | nvd | `cve-2017-14445` | code | 1.0 |
| 6628 | An exploitable stack-based buffer overflow vulnerability exists in Insteon Hub running firmware version 1012. The HTTP server implementation unsafely extracts p | nvd | `cve-2017-14446` | code | 1.0 |
| 6629 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01bad0 the value for the | nvd | `cve-2017-16338` | code | 1.0 |
| 6630 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01bb1c the value for the | nvd | `cve-2017-16339` | code | 1.0 |
| 6631 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c0e8 the value for the | nvd | `cve-2017-16340` | code | 1.0 |
| 6632 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c224 the value for the | nvd | `cve-2017-16341` | code | 1.0 |
| 6633 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c254 the value for the | nvd | `cve-2017-16342` | code | 1.0 |
| 6634 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c284 the value for the | nvd | `cve-2017-16343` | code | 1.0 |
| 6635 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c2c8 the value for the | nvd | `cve-2017-16344` | code | 1.0 |
| 6636 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c318 the value for the | nvd | `cve-2017-16345` | code | 1.0 |
| 6637 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01c368 the value for the | nvd | `cve-2017-16346` | code | 1.0 |
| 6638 | An attacker could send an authenticated HTTP request to trigger this vulnerability in Insteon Hub running firmware version 1012. At 0x9d01e7d4 the value for the | nvd | `cve-2017-16347` | code | 1.0 |
| 6639 | A design flaw issue was found in the Red Hat OpenStack Platform director use of TripleO to enable libvirtd based live-migration. Libvirtd is deployed by default | nvd | `cve-2017-2637` | code | 1.0 |
| 6640 | Vulnerability in the Primavera P6 Enterprise Project Portfolio Management component of Oracle Primavera Products Suite (subcomponent: Web Access (Apache Commons | nvd | `cve-2017-3503` | code | 1.0 |
| 6641 | Vulnerability in the Oracle Identity Manager component of Oracle Fusion Middleware (subcomponent: Rules Engine). The supported version that is affected is 11.1 | nvd | `cve-2017-3553` | code | 1.0 |
| 6642 | The drag-and-drop (DnD) function in VMware Workstation 12.x before version 12.5.4 and Fusion 8.x before version 8.5.5 has an out-of-bounds memory access vulnera | nvd | `cve-2017-4901` | code | 1.0 |
| 6643 | The WHMCS Reseller Module V2 2.0.2 in Softaculous Virtualizor before 2.9.1 does not verify the user correctly, which allows remote authenticated users to cont | nvd | `cve-2017-6513` | code | 1.0 |
| 6644 | NfSen before 1.3.8 allows remote attackers to execute arbitrary OS commands via shell metacharacters in the customfmt parameter (aka the "Custom output format"  | nvd | `cve-2017-7175` | code | 1.0 |
| 6645 | TP-Link C2 and C20i devices through firmware 0.9.1 4.2 v0032.0 Build 160706 Rel.37961n allow remote code execution with a single HTTP request by placing shell c | nvd | `cve-2017-8220` | code | 1.0 |
| 6646 | A vulnerability in the role-based resource checking functionality of the Cisco Unified Computing System (UCS) Director could allow an authenticated, remote atta | nvd | `cve-2018-0238` | code | 1.0 |
| 6647 | An issue was discovered in MyBiz MyProcureNet 5.0.0. A malicious file can be uploaded to the webserver by an attacker. It is possible for an attacker to upload  | nvd | `cve-2018-11091` | code | 1.0 |
| 6648 | An issue was discovered in Xen 4.7 through 4.10.x. libxl fails to pass the readonly flag to qemu when setting up a SCSI disk, due to what was probably an errone | nvd | `cve-2018-12892` | code | 1.0 |
| 6649 | In OnlineJudge 2.0, the sandbox has an incorrect access control vulnerability that can write a file anywhere. A user can write a directory listing to /tmp, and  | nvd | `cve-2018-16367` | code | 1.0 |
| 6650 | A sandbox escape issue was discovered in VyOS 1.1.8. It provides a restricted management shell for operator users to administer the device. By issuing various s | nvd | `cve-2018-18555` | code | 1.0 |
| 6651 | A privilege escalation issue was discovered in VyOS 1.1.8. The default configuration also allows operator users to execute the pppd binary with elevated (sudo)  | nvd | `cve-2018-18556` | code | 1.0 |
| 6652 | A vulnerability was discovered in the Java VM component of Oracle Database Server. Supported versions that are affected are 11.2.0, 12.1.0, 12.2.0 and 18 | nvd | `cve-2018-3110` | code | 1.0 |
| 6653 | An exploitable vulnerability exists in the smart cameras RTSP configuration of the Samsung SmartThings Hub STH-ETH-250 - Firmware version 0.20.17. The device in | nvd | `cve-2018-3856` | code | 1.0 |
| 6654 | On Samsung SmartThings Hub STH-ETH-250 devices with firmware version 0.20.17, the video-core process incorrectly extracts fields from a user-controlled JSON pay | nvd | `cve-2018-3863` | code | 1.0 |
| 6655 | An exploitable buffer overflow vulnerability exists in the samsungWifiScan handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250 - Firmware | nvd | `cve-2018-3866` | code | 1.0 |
| 6656 | An exploitable stack-based buffer overflow vulnerability exists in the samsungWifiScan callback notification of video-core's HTTP server of Samsung SmartThings  | nvd | `cve-2018-3867` | code | 1.0 |
| 6657 | An exploitable buffer overflow vulnerability exists in the credentials handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250 - Firmware ver | nvd | `cve-2018-3872` | code | 1.0 |
| 6658 | An exploitable buffer overflow vulnerability exists in the credentials handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250-Firmware versi | nvd | `cve-2018-3873` | code | 1.0 |
| 6659 | An exploitable buffer overflow vulnerability exists in the credentials handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250-Firmware versi | nvd | `cve-2018-3874` | code | 1.0 |
| 6660 | An exploitable buffer overflow vulnerability exists in the credentials handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250-Firmware versi | nvd | `cve-2018-3875` | code | 1.0 |
| 6661 | An exploitable buffer overflow vulnerability exists in the credentials handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250-Firmware versi | nvd | `cve-2018-3877` | code | 1.0 |
| 6662 | Multiple exploitable buffer overflow vulnerabilities exist in the credentials handler of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250 devices | nvd | `cve-2018-3878` | code | 1.0 |
| 6663 | An exploitable stack-based buffer overflow vulnerability exists in the database 'find-by-cameraId' functionality of video-core's HTTP server of Samsung SmartThi | nvd | `cve-2018-3880` | code | 1.0 |
| 6664 | An exploitable buffer overflow vulnerability exists in the camera "replace" feature of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250 devices w | nvd | `cve-2018-3902` | code | 1.0 |
| 6665 | On Samsung SmartThings Hub STH-ETH-250 devices with firmware version 0.20.17, the video-core process incorrectly extracts fields from a user-controlled JSON pay | nvd | `cve-2018-3903` | code | 1.0 |
| 6666 | An exploitable buffer overflow vulnerability exists in the camera 'update' feature of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250 - Firmware | nvd | `cve-2018-3904` | code | 1.0 |
| 6667 | An exploitable buffer overflow vulnerability exists in the camera "create" feature of video-core's HTTP server of Samsung SmartThings Hub STH-ETH-250 devices wi | nvd | `cve-2018-3905` | code | 1.0 |
| 6668 | On Samsung SmartThings Hub STH-ETH-250 devices with firmware version 0.20.17, the video-core process insecurely extracts the fields from the "shard" table of it | nvd | `cve-2018-3917` | code | 1.0 |
| 6669 | An exploitable stack-based buffer overflow vulnerability exists in the retrieval of database fields in video-core's HTTP server of Samsung SmartThings Hub STH-E | nvd | `cve-2018-3919` | code | 1.0 |
| 6670 | In browser editing in Atlassian Bitbucket Server from version 4.13.0 before 5.4.8 (the fixed version for 4.13.0 through 5.4.7), 5.5.0 before 5.5.8 (the fixed ve | nvd | `cve-2018-5225` | code | 1.0 |
| 6671 | The default permissions of /dev/kmem in Linux versions before 2.0.36 allows IP spoofing | nvd | `cve-1999-0426` | code | 1.0 |
| 6672 | VAXstations running Open VMS 5.3 through 5.5-2 with VMS DECwindows or MOTIF do not properly disable access to user accounts that exceed the break-in limit thres | nvd | `cve-1999-1324` | code | 1.0 |
| 6673 | CGI Script Center News Update 1.1 does not properly validate the original news administration password during a password change operation, which allows remote a | nvd | `cve-2000-0944` | code | 1.0 |
| 6674 | The default configuration for the domain name resolver for Microsoft Windows 98, NT 4.0, 2000, and XP sets the QueryIpMatching parameter to 0, which causes Wind | nvd | `cve-2000-1218` | code | 1.0 |
| 6675 | Buffer overflow in FTP server in HPUX 11 allows remote attackers to execute arbitrary commands by creating a long pathname and calling the STAT command, which u | nvd | `cve-2001-0248` | code | 1.0 |
| 6676 | Heap overflow in FTP daemon in Solaris 8 allows remote attackers to execute arbitrary commands by creating a long pathname and calling the LIST command, which u | nvd | `cve-2001-0249` | code | 1.0 |
| 6677 | Lightwave ConsoleServer 3200 does not disconnect users after unsuccessful login attempts, which could allow remote attackers to conduct brute force password gue | nvd | `cve-2001-0395` | code | 1.0 |
| 6678 | Format string vulnerability in Infodrom cfingerd 1.4.3 and earlier allows a remote attacker to gain additional privileges via a malformed ident reply that is pa | nvd | `cve-2001-0609` | code | 1.0 |
| 6679 | Apache on MacOS X Client 10.0.3 with the HFS+ file system allows remote attackers to bypass access restrictions via a URL that contains some characters whose ca | nvd | `cve-2001-0766` | code | 1.0 |
| 6680 | Knox Arkeia server 4.2, and possibly other versions, uses a constant salt when encrypting passwords using the crypt() function, which makes it easier for an att | nvd | `cve-2001-0967` | code | 1.0 |
| 6681 | Symantec LiveUpdate before 1.6 does not use cryptography to ensure the integrity of download files, which allows remote attackers to execute arbitrary code via  | nvd | `cve-2001-1125` | code | 1.0 |
