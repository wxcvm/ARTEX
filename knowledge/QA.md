# knowledge/ 语料体检

> 语料是记录而非工程：`.py` 里可能贴的是 JS / 用法说明 / Python 2。
> 只有台账与 FTS 不一致算阻断（P1）。

| 项 | 值 |
|---|---|
| poc_dirs | 429 |
| poc_files | 432 |
| poc_bytes | 1488773 |
| py_files | 162 |
| py_not_runnable | 44 |
| py_kind_detail | NOT_PYTHON=4, SYNTAX=40 |
| sh_files | 65 |
| duplicate_files | 0 |
| db_rows | 448 |
| fts_rows | 448 |

## P1（阻断）

（无）

## P2（记录）

| 类型 | 文件 | 备注 |
|---|---|---|
| SYNTAX | `1panel-login-password-encryption-rsa-aes-3seg\1panel_password_encrypt.py` | SyntaxError: invalid character '→' (U+2192) |
| SYNTAX | `1panel-安全入口拦截页与-swagger-拼接-未认证绕读完整-swagger-v1.10.28-lts-实测-无-swagger-时的版本夹逼探针\poc.py` | SyntaxError: invalid character '「' (U+300C) |
| SYNTAX | `android-app-自研-api-加密协议还原-dex-反汇编定位-aes-密钥-param-密文协议-登录三层校验差分定位\poc.py` | SyntaxError: invalid character '（' (U+FF08) |
| SYNTAX | `discuz-7.x-ajax.php-checkusername-用户名枚举-零配额成本判定账号是否存在\poc.py` | SyntaxError: invalid decimal literal |
| SYNTAX | `discuz-管理员-cookie-伪造-authkey-ua-签名-由-sql-注入取三要素\discuz_cookie_forge.py` | SyntaxError: invalid character '（' (U+FF08) |
| SYNTAX | `discuz-论坛平台-sclub-免费论坛-连接层封禁阈值探测方法-区分-并发触发-与-累计触发\poc.py` | SyntaxError: invalid syntax |
| SYNTAX | `fakertoken-ob-的-api-参数解密-硬编码-aes-ecb-密钥-obty20220712obty-与-requestdata-不可解结论\poc.py` | SyntaxError: invalid character '，' (U+FF0C) |
| SYNTAX | `go-captcha-drag-sobel-edge-template-crack\poc.py` | SyntaxError: unmatched ')' |
| SYNTAX | `mail-auth-missing-spf-dkim-dmarc-email-spoofing\poc.py` | SyntaxError: invalid character '（' (U+FF08) |
| SYNTAX | `ocs-在线客服中台-访客会话越权读取-vid-idor\ocs_idor_poc.py` | SyntaxError: invalid character '：' (U+FF1A) |
| SYNTAX | `ocs-在线客服系统-stomp-over-sockjs-握手真值-客服侧-protobuf-信封-usermsgpack\poc.py` | SyntaxError: invalid character '（' (U+FF08) |
| SYNTAX | `openresty-1.29-spring-boot-请求走私-cl.te-te.cl-te.te-与-h2-降级全\poc.py` | SyntaxError: invalid character '（' (U+FF08) |
| SYNTAX | `openssh-regresshion-cve-2024-6387-崩溃验证-7etsuo-poc-落盘-python-复刻-v1-v2\regresshion.py` | SyntaxError: invalid character '（' (U+FF08) |
| NOT_PYTHON | `phpwind-7.x-入口-action-全量枚举方法论-251-个-action-10-类危险标记-写点判死结论集-pw_var_export-转义-pcv\pw-action-enum.py` | file extension says .py but the body is not Python (notes / pasted JS / usage text) |
| NOT_PYTHON | `sockjs-stomp-websocket-客户端实现要点-json-数组发送侧编码-版本协商-cdp-抄帧法\poc.py` | file extension says .py but the body is not Python (notes / pasted JS / usage text) |
| SYNTAX | `sockjs-文本通道无法承载-protobuf-二进制-成因-对照实验与判据\poc.py` | SyntaxError: invalid character '：' (U+FF1A) |
| SYNTAX | `spa-spring-cloud-网关批量未授权接口探测-xxa-前缀型\poc.py` | SyntaxError: invalid character '，' (U+FF0C) |
| NOT_PYTHON | `sub2api-easypay-支付回调验签绕过-参数走私-签名串拼接不一致-伪造回调入账\easypay_forge_callback.py` | file extension says .py but the body is not Python (notes / pasted JS / usage text) |
| SYNTAX | `sub2api-邮箱验证码接口账号存在性枚举-oracle-公开-api-v1-auth-send-verify-co\sub2api_email_enum_oracle.py` | SyntaxError: invalid character '：' (U+FF1A) |
| SYNTAX | `tls-证书-san-反查域名-同-ip-vhost-路由发现-纯-ip-段挖真实域名-避免错过-target\tls_san_vhost_discovery.py` | SyntaxError: unterminated triple-quoted string literal (detected at line 70) |
| SYNTAX | `vite-spa-多子应用-api-面递归提取-子应用发现-manifest-chunk-清单-api-正则\vite-subapp-api-extract.py` | SyntaxError: invalid character '（' (U+FF08) |
| SYNTAX | `web-缓存投毒-缓存欺骗面一站式测绘脚本-含缓存键归一化与-unkeyed-头判定\cache_poison_probe.py` | SyntaxError: invalid syntax. Perhaps you forgot a comma? |
| SYNTAX | `webman-admin-tablecontroller-crud-代码生成器-docblock-注入-写-app-controller-.php-monito\table_crud_docblock_inject.py` | SyntaxError: invalid character '，' (U+FF0C) |
| SYNTAX | `webman-响应加密-x-encrypted-1-硬编码-aes-密钥解密-加密请求体构造\webman_xencrypted_bypass.py` | SyntaxError: invalid character '，' (U+FF0C) |
| SYNTAX | `webman-自动路由绕过鉴权中间件-controller-action-零凭证直连\webman_autoroute_authbypass.py` | SyntaxError: invalid character '：' (U+FF1A) |
| SYNTAX | `全量参数钥匙-批量候选参数名越过老-servlet-的-integer.parseint-门槛-分组探测-二分收敛-通用可复用\poc.py` | SyntaxError: invalid character '（' (U+FF08) |
| SYNTAX | `出口-ip-被宝塔-waf-拉黑后的应急出口-免费-http-代理池-抓取-验活-轮换-schannel-坑\proxy-pool-egress.py` | SyntaxError: invalid character '（' (U+FF08) |
| SYNTAX | `单请求型隧道下的内网渗透方法论-op-chk-主机发现-目标侧-php-绕过多轮协议限制\poc.py` | SyntaxError: invalid character '「' (U+300C) |
| SYNTAX | `商城游客订单自认证头机制-authorization-guest-base64url-email-n-order_password-含空口令-畸形头边界实测\poc.py` | SyntaxError: invalid character '，' (U+FF0C) |
| SYNTAX | `多租户直播平台前端响应加密还原-aes-256-cbc-硬编码-key-iv-时间戳尾部拼接\poc.py` | SyntaxError: invalid character '（' (U+FF08) |
| SYNTAX | `宝塔-pure-ftpd-口令面判据-存在-不存在用户同一条-530-无法枚举-只有-230-算通-字典来源与-put-马战术\poc.py` | SyntaxError: invalid character '（' (U+FF08) |
| SYNTAX | `尾斜杠触发-traefik-301-清理泄露隐藏路由-base-通用绕过手法\traefik_slash_hidden_base.py` | SyntaxError: unterminated string literal (detected at line 26) |
| SYNTAX | `支付回调-通知端点探测法-静态零引用判定-未知路由三态判别-按站点命名风格派生字典-含负结论口径\poc.py` | SyntaxError: invalid character '，' (U+FF0C) |
| NOT_PYTHON | `本地只读反代-把前端调用的-api-改写为大写转发-让原版管理后台界面直接可用\poc.py` | file extension says .py but the body is not Python (notes / pasted JS / usage text) |
| SYNTAX | `本机-fake-ip-dns-规则型本地代理下-用-connect-自定义-sni-直达-nginx-vhost-绕过假性-目标封-ip\connect_sni_vhost.py` | SyntaxError: invalid syntax |
| SYNTAX | `环境事实-全体适用-本机-python-urllib-走-ie-注册表代理-127.0.0.1-7892-导致请求超时-curl-正常\poc.py` | SyntaxError: invalid character '：' (U+FF1A) |
| SYNTAX | `盲-ssrf-判定法-差分三值-死地址基线-内网-外带-a-b-a-交替复测识破抖动-附-零差异-无-ssrf-实测样例与真-ssrf-对照特征\ssrf_diff_test.py` | SyntaxError: invalid character '，' (U+FF0C) |
| SYNTAX | `移动-app-侦察-从-pcap-配置分发文件定位真实后端-认证头校难强度对照实验\poc.py` | SyntaxError: invalid character '（' (U+FF08) |
| SYNTAX | `站群同源判定-多-vhost-是否指向同一-document_root-body-diff-相似度-尺寸-4-主机名长度差\poc.py` | SyntaxError: invalid character '≈' (U+2248) |
| SYNTAX | `端口开放判定判据-代理-connect-假成功陷阱-以应用层数据为准-宝塔-pure-ftpd-nginx-403-实测\poc.py` | SyntaxError: invalid character '（' (U+FF08) |
| SYNTAX | `统一-20000-未登录拦截的绕过穷尽判定法-含-body-ct-归因控制实验-防-116-误报\poc.py` | SyntaxError: invalid character '“' (U+201C) |
| SYNTAX | `腾讯云-sts-临时凭证-cos-q-signature-手工签名-写入-边界判定-可用脚本\poc.py` | SyntaxError: invalid character '（' (U+FF08) |
| SYNTAX | `跨-vhost-同源-xss-可执行性判定矩阵-多域名站群-uploads-可达性-app-admin-交叉对照\poc.py` | SyntaxError: invalid character '，' (U+FF0C) |
| SYNTAX | `阿里云验证码-2.0-aliyuncaptcha-自动化跑通并取得-captchaverifyparam-app-web-登录风控绕过通用打法\poc.py` | SyntaxError: invalid character '：' (U+FF1A) |
| NO_SHEBANG | `203.0.113.1-.23-先验排除表-两条恒定假阳性-code-验证码端点-爆破前置-五处已穷尽面\poc.sh` | ## 适用目标 |
| NO_SHEBANG | `203.0.113.10-.23-先验排除表-两条恒定假阳性-code-验证码端点-爆破前置-五处已穷尽面\poc.sh` | ## 适用目标 |
| NO_SHEBANG | `3x-ui-x-ui-assets-unauth-directory-listing-and-dual-404-base\poc.sh` | # 3X-UI / x-ui 面板 /assets/ 匿名目录列举 |
| NO_SHEBANG | `gm-collection-邮件工具站免登录开放接口未授权探测清单-outlook-icloud-邮件与凭据查询\poc.sh` | # 「GM Collection 管理面板」类邮件工具站 —— 免登录开放接口探测清单 |
| NO_SHEBANG | `go-自研代理服务-管理口-代理口-免凭据探测与绕过穷尽核查脚本\go_proxy_panel_probe.sh` | # 目的：对 Go 自研代理类服务（管理口 + 代理口）做"不需凭据即可判定"的探测与穷尽绕过核查 |
| NO_SHEBANG | `lfi-注入点定位法-先找不经过守卫的可达文件目录-再用数组参数判定是否读参数\poc.sh` | ## 目标 |
| NO_SHEBANG | `mysql-3306-exposed-host-not-allowed-ip-whitelist-three-egres\poc.sh` | ## 现象 |
| NO_SHEBANG | `new-api-漏洞清单与未修复版本边界对照-cve-2026-42339-ssrf-33655-64865-30886\poc.sh` | # New API 漏洞清单与「未修复版本边界」对照表（外部情报 + 靶标实测 + 源码审计修正） |
| NO_SHEBANG | `nginx-location-绕过-2e-2e-系列必然失败-及用-不存在路径对照-证伪-..-假穿越\poc.sh` | ## 结论 |
| NO_SHEBANG | `nginx-tomcat-..-归一化差异-穿-location-白名单-跨-context-横向-含绕开-nginx-死上游-504-manager-ip-阀\poc.sh` | ## 一句话 |
| NO_SHEBANG | `phpwind-7.x-pw_ajax.php-action-pcdelimg-fieldname-布尔盲注-wooyun-2010-0417-影响-phpwi\pcdelimg-blind-sqli.sh` | 未修形态源码（runs/pw-src-732/maddogfyg-itangyou/pw_example.com L19 |
| NO_SHEBANG | `phpwind-7.x-pw_api.php-api-class_base.php-mode-任意文件包含-cn-va10-02-点名文件-不依赖-regist\pwapi-mode-lfi.sh` | api/class_base.php L110-117（pw75sp3 源码逐行）： |
| NO_SHEBANG | `phpwind-7.x-pw_api.php-api-class_other.php-threadscategory-写码-rce-wooyun-2010-04\pwapi-threadscategory-rce.sh` | api/class_base.php L70-123（调用链）： |
| NO_SHEBANG | `phpwind-7.x-后台入口名定位-读-pw_config.db_adminfile-前台-db_adminfile-渲染-不做文件名爆破\poc.sh` | # PHPWind 7.x 后台入口名定位（取代无脑文件名爆破） |
| NO_SHEBANG | `phpwind-8.7-注册三必填绕过-regreason-_hexie-apartment-地区id-注册成功且自动登录\poc.sh` | ## 背景 |
| NO_SHEBANG | `phpwind-hack-插件入口-hack.php-h_name-lt-插件-gt-绕过-hack-lt-插件-gt-index.php-的-9b-forbi\phpwind-hack-plugin-entry.sh` | 原理：PHPWind 的 hack 插件目录 hack/<插件名>/example.com 顶部有 readover 守 |
| NO_SHEBANG | `phpwind-站点精确版本与数据库表前缀的无认证识别\poc.sh` | # PHPWind 站点精确版本与表前缀被动识别 |
| NO_SHEBANG | `proxy-host-header-false-origin\poc.sh` | # 方法论陷阱：用本地 HTTP 代理做 Host 头 vhost 探测会造出"命中源站"的假象 |
| NO_SHEBANG | `thinkphp-精确版本与绝对路径识别-异常页法\poc.sh` | ## 识别方法（三步，零风险、纯被动） |
| NO_SHEBANG | `tomcat-manager-快速判死-403-页字节级指纹-同版本-stock-页跨机器-sha1-完全一致-端口面普查-只暴露-nginx-头伪造必然无效\poc.sh` | ## 结论（一条即可判死） |
| NO_SHEBANG | `tomcat-拿-shell-路线的快速判死-put-405-options-无-put-cve-2017-12615-12617-2025-24813-全灭\poc.sh` | ## 判死清单与实测回显（省得后来者重踩） |
| NO_SHEBANG | `vue3-vite-spa-接口测绘法-入口-js-只给路由-真实-api-在懒加载-chunk-的数据层-use-.js-里\poc.sh` | # 目标画像 |
| NO_SHEBANG | `webman-路由大小写不敏感-绕过按路径精确匹配的-功能关闭-开关-可救回代码生成器写文件-monitor-热重载-rce\poc.sh` | # 通用技巧：webman（workerman 系）路由大小写不敏感 × 业务侧「功能开关/路径黑名单」按精确路径匹配  |
| NO_SHEBANG | `wordpress-woocommerce-匿名攻击面清单-版本泄露-rest未授权-store-api非授权写入-xml-rpc爆破\wordpress_woo_unauth_recon.sh` | # WooCommerce 站（匿名可打）完整信息收集与未授权写入清单 |
| NO_SHEBANG | `代理分流造假样本-target30.example.com-经-7892-返回伪造-502-404-hello\poc.sh` | 实测记录（2026-09-28）： |
| NO_SHEBANG | `代理分流造假样本-tgt-经-7892-返回伪造-502-404-hello\poc.sh` | 实测记录（2026-09-28）： |
| NO_SHEBANG | `出口-ip-与-host-头是两个独立变量-环境代理静默改道导致宝塔-waf-封禁被误读-curl-必须-noproxy\poc.sh` | # 出口 IP 决定成败：本机 curl 默认走环境代理 127.0.0.1:17992，一发多即被宝塔 WAF 全站封 |
| NO_SHEBANG | `前端硬编码-requestdata-常量重放-未授权签发第三方平台-token-与线路域名路由表\hardcoded_requestdata_replay.sh` | # 「前端硬编码 requestData 常量重放 → 未授权签发第三方平台 token」通用手法 |
| NO_SHEBANG | `前端资源投毒可行性判定法-客户端渲染资金字段-无-sri-无-csp-资源域写方法-含正确分级口径\poc.sh` | # 目标 |
| NO_SHEBANG | `宝塔-waf-来源封禁实测判据-multipart-上传后全出口-000-proxyerror-而同代理到-baidu-200-区分-代理失效-vs-目标拒我\waf-source-ban-criteria.sh` | # 判据：multipart 上传 POST 后目标立即拒绝我方来源（宝塔 WAF 行为） |
| NO_SHEBANG | `对手检测能力画像-站群运营方-thinkphp3-后门自查脚本-jnbgg.sh-我们的规避指南\adversary_tp3_backdoor_scan_jnbgg.sh` | # 对手自查脚本 = 我们的规避指南 |
| NO_SHEBANG | `文件爆破假阳性证伪法-换词对照-原始-location-相对路径-后缀断崖-nginx-通配规则识别\poc.sh` | ## 问题 |
| NO_SHEBANG | `智讯网站群-注册接口用户名枚举-口令大小写不敏感判定\poc_user_enum.sh` | # 前置：验证码 GET /code.php（60x30 PNG，4 位纯数字、无干扰、一次有效） |
| NO_SHEBANG | `智讯网站群管理系统-api-未授权访问-form-urlencoded-传参-action-index-泄露站点配置\poc_zhixun_api.sh` | # 指纹识别 |
| NO_SHEBANG | `本机红队环境事实与替代方案-windows-elf-平台限制\poc.sh` | # 本机（Windows）红队作业环境事实与替代方案 |
| NO_SHEBANG | `环境事实-全体适用-本机-windows-出口被-tun-劫持-探测改走演练-vps-工具白名单与用法\poc.sh` | ## 一、本机（Windows）链路事实（实测） |
| NO_SHEBANG | `白标-spa-商城溯源到上游开源项目-品牌名兜底值-app_version-转离线源码审计\poc.sh` | # 思路 |
| NO_SHEBANG | `直播平台未授权注册绕过短信验证-api-user-simple_regist-空-mobile_hash-直接签发-token\poc.sh` | ## 原理 |
| NO_SHEBANG | `网段存活探测纪律-nmap-pn-会把全部-ip-判为-up-必须用-sn-复核\poc.sh` | ## 现象 |
| NO_SHEBANG | `通用方法-开源组件识别-拉源码离线审计-定向验证-找-advisory-之外的路由漏挂\poc.sh` | # 开源组件识别 → 拉源码离线审计 → 定向验证 |

## 说明

- `PY2` / `NOT_PYTHON`：上游记录的原始形态，保留原样（改名会破坏目录与台账的对应）。
- `NO_SHEBANG`：不影响阅读，直接 `bash <file>` 或补首行即可。
- `DUPLICATE`：同一段脚本在多个目录各存一份，属历史记录。
