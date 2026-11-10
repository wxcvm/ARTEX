# 目标
在"客户端下单已被服务端强校验挡死"的资金站上，判断能否靠**伪造第三方支付通知**充值到账。
第一步必须定位"本站接收上游通知的回调端点"。

# 一、静态先行（零请求，能砍掉大半工作量）
1. 前端全量 grep（主包 + 全部 chunk）：
   `notify` `callback` `callBack` `returnUrl` `backUrl` `payNotify` `rechargeNotify`
   `orderNotify` `payCallback` `asyncNotify` `syncNotify` `handleNotify` `payResult`
   `payFinish` `paySuccess` `tradeNotify` `payReturn` `thirdNotify` `<厂商>Notify`
   → 命中要逐条看是不是**噪声**：UI 组件名（base-notify/notify-slide-fade/useNotify）、
     浏览器 API（requestIdleCallback）、弹窗回调参数（callback:()=>{}）、棋牌转移参数（callback）。
2. 通道/商户配置接口的**全字段枚举**：把响应 JSON 递归收集所有 key，
   看有没有 `notifyUrl`/`callbackUrl`/`returnUrl`/`gateway`/`host`（服务端配的回调址会直接暴露路径）。
3. 服务注册表：找到 axios 多实例配置（`baseURL`/`successCode`），确认**本地后端前缀**是哪些；
   非本地前缀（第三方 sports/faker 服务）不用去猜回调。
4. 签名函数：`md5`/`sha1`/`sha256`/`hmac`/`CryptoJS`/`JSEncrypt`/`sign(`/`encrypt(`。
   命中要看是不是混淆 chunk 的随机字节或无关功能（本例 `sign(` 全是签到 signRecord/signAmount）。
   **前端无签名实现 = 回调验签 100% 在服务端 = 无算法无密钥，构造合法签名不可行**（不要在线爆破密钥）。

# 二、动态探测的三态判别（关键：不能把"路由不存在"误判成"被拦"）
先做对照，确定三态：
  · **未知 /api 路由** → 404 + SPA 404 页（本例 176394 字节 HTML）
  · **存在但未登录** → 200 + 统一业务码（本例 code 20000）
  · **非本地前缀路径** → 从不进后端：POST 得 nginx 405、GET 得 SPA 首页(13140B)
再打候选，判据：**只要不是这三态（例如业务码/500 异常/纯文本 SUCCESS/OK）即为命中候选**。

# 三、候选字典（按证据排优先级，而非随机）
1. 厂商前缀组：从通道配置的 `payCenterUrl`（如 /ubit/pay、/newc2c/rechargeEntrance）取前缀，
   拼 `<vendor>/{notify,callback,payNotify,notifyUrl}`（16 条左右）。
2. `/api/` 通用名组（20 条）：notify/callback/payNotify/payNotifyV2/notifyUrl/payResult/
   paySuccess/payFinish/payReturn/tradeNotify/orderNotify/rechargeNotify/thirdNotify/
   asyncNotify/handleNotify/payCallback/<vendor>Notify。
3. 顶层组（5 条）：/notify、/pay/notify、/payNotify、/callback、/pay/callback。
4. 遗留 `.php` 组：站点若已有 `/api/xxx.php`（本例有 getTokenPc.php），补 `.php` 变体。
5. **按站点自身命名风格派生**（最有价值、最容易被忽略）：把已知接口名当模板，
   例如站内叫 `submitThirdV4`/`confirmEbc2c`/`checkOnlineOrder` →
   试 `submitThirdV5`/`submitThirdNotify`/`notifyThirdV4`/`thirdNotifyV4`/`payNotifyV4`/
   `handlePayResult`/`queryPayResult`/`syncOrderStatus`/`thirdCallback`。
   json 体与 form 体**各打一遍**（存在性判定与 Ct 无关，但不排除只认 form 的端点）。

# 四、结论口径（必须写清的限制）
· 字典法只能证明"这些命名不存在"，**不能证明端点绝对不存在**；
· 第三方支付回调通常还有 **IP 白名单**（出口不在网关网段时，即使名字对也会 403）；
· 因此负结论应写成"按命名规律未发现可匿名触达的通知入口"，并列出已试路径数、
  判别依据、以及"需要什么才能继续"（网关网段出口 / 端点名泄露源：源码、日志、配置泄露）。
