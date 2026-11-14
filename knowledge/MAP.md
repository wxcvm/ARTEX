# 知识库对账（`knowledge.db` ↔ `pocs/`）

| 项 | 数量 |
|---|---|
| `knowledge.db` 条目 | **508** |
| `pocs/` 原文目录 | **429** |
| 已配对 | **428** |
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
