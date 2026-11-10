# 机制（来自前端 order-Cgx3wnzv.js 的反编译）
游客订单不用 JWT，而是把 (email, order_password) 直接拼进 Authorization 头：

```js
function mkGuestHeader(email, pwd) {
  const raw = email.trim() + "\n" + pwd.trim();          // 用换行分隔
  const b64 = btoa(raw).replace(/\+/g,'-').replace(/\//g,'_').replace(/=+$/,'');  // base64url 去填充
  return "Guest " + b64;
}
```

调用形如：
```
GET /api/v1/guest/orders?page=1&page_size=20
Authorization: Guest dGVzdEBleGFtcGxlLmNvbQp0ZXN0
```
（`dGVzdEBleGFtcGxlLmNvbQp0ZXN0` = base64url("test@example.com\ntest")）

# 边界实测（只读 GET，7 组，VPS 出口）
| 请求 | 结果 |
|---|---|
| 不带 Authorization | 200 `{"status_code":400,"msg":"游客邮箱不能为空"}` |
| `Guest <b64("test@example.com\n")>`（空口令） | 同上（400 邮箱不能为空）→ **空口令被拒** |
| `Guest <b64("test@example.com")>`（无换行） | 同上（解析失败即视为无邮箱） |
| `Guest !!!not-base64!!!` | 同上 |
| `Guest <b64("test@example.com\ntest")>` | 200 `{"status_code":0,"msg":"success","data":[],"pagination":{...,"total":0}}` |
| `GET /guest/payments/latest` 空口令 | 200 `{"status_code":400,"msg":"请求参数错误"}` |

# 结论（如实）
**未发现越权**：空口令/畸形头被明确拒绝；错误口令对不存在的邮箱返回空列表。
(email, order_password) 应共同参与过滤，口令错误即查不到数据（前端 useGuestOrderDetail 另有 `authInvalid`「认证无效」分支）。
**未验证项**：无法区分「邮箱不存在」与「不校验口令」，因为需要一个真实存在的测试邮箱 —— 本轮不使用任何真实用户数据，故不下结论。

# 拿到价值
一旦获得任一真实「邮箱 + 该邮箱下过单时设置的 order_password」，可直接：
```
GET /api/v1/guest/orders?page=1&page_size=50                 # 列出该邮箱全部订单
GET /api/v1/guest/orders/<order_no>                          # 订单详情
GET /api/v1/guest/orders/<order_no>/fulfillment/download     # 下载交付内容（卡密）
```
⇒ 在拿到凭据后，这是**无需登录**即可取数据的通道（注意下载交付内容是敏感动作，需授权）。