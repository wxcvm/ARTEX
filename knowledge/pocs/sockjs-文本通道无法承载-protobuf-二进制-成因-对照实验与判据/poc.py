## 结论
SockJS 的 websocket transport 是**文本**通道：客户端消息必须 JSON 数组包裹的字符串，
服务端（Spring SockJsMessageCodec）把每个字符串按 **UTF-8** 还原为字节再交给 STOMP 解码器。
因此**无法携带任意二进制 protobuf**：只要 body 含 >=0x80 的字节，content-length 与实际字节数必然错位。

## 实测证据（203.0.113.85:3493，OCS 客服系统）

构造 UserMsgPack{MsgId=100, Data=MsgAuthReq{Token=<32hex>, Account=101096}}：
  hex = 08 64 1a 26 0a 20 63 39 31 34 65 34 64 64 39 38 35 66 34 37 33 61 38 32 63 66 64 37 63 30 38 38 37 62 31 34 33 34 10 e8 95 06
                                        ^^^^^^^^^^^^^^^^^ Token(32B, ASCII)          ^^^^^^^^^ Account varint = e8 95 06 (非 ASCII)

以 latin-1 映射进 JSON 字符串后发送（content-length 声明 38 字节）：

  -> ERROR
     message:Frame must be terminated with a null octet
     content-length:0
  -> c[1002,""]            # 连接被服务端关闭

## 关键对照实验（证明成因是字节内容而非 destination）
对**同一个** body 只改 Account 取值：

| body | hex 尾部 | 是否全 ASCII | 服务端响应 |
|---|---|---|---|
| auth100_noAcct (Account 省略) | ... 34 | 是 | **无 ERROR（成功穿过解析层）** |
| auth100_acct1   (Account=1)   | ... 34 10 01 | 是 | **无 ERROR** |
| sess102_empty   (MsgId=102)   | 08 66 1a 00 | 是 | **无 ERROR** |
| auth100_bigAcct (Account=101096)| ... 10 e8 95 06 | 否 | **ERROR: Frame must be terminated with a null octet + CLOSE** |

并且该 ERROR 在**所有** destination 上同样出现（/app/staff/ack、/app/staff/auth、/app/staff、
/app/staff/login、/app/staff/sessions、/app/sessions、/app/visitor/ack 全部如此）
=> 说明报错点在 **STOMP 解码层（destination 路由之前）**，与 destination 是否存在无关。

## 推论（对攻击的限制）
- `MsgAuthReq.Account` 是 `uint32`。凡是取值落在 128..16383 之间的 account，
  varint 必含 >=0x80 字节 => **无法经访客 SockJS 通道完成客服认证**。
- 真实客服 ID 几乎必然 >127（本系统 cid=101096 已是 6 位数）。
- => 客服端 protobuf 必须走**原始 WebSocket 二进制帧**（另一端点），或该 SockJS 通道存在
  以 base64/hex 承载 body 的编码约定（未观察到证据）。

## 可复现脚本
C:\rt\dgha\runs\ut_19_staff_ascii.py  （ASCII-safe 与 非ASCII 两组 body × 7 destination，逐条输出对照）
