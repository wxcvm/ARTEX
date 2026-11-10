## 1. 连接与握手（已实测成功）

WS 端点（SockJS websocket transport）：
  wss://<host>:3493/api/v1/v/ws/{3位随机数字}/{8位随机session}/websocket

必须参数（缺 UA 会 403；直连 IP 会 TLS 失败）：
  http_proxy_host=127.0.0.1, http_proxy_port=7892
  sslopt={'cert_reqs': ssl.CERT_NONE, 'check_hostname': False}
  header: Origin: https://<host>:3449 / Referer / UA / x-v-token / cid / vid
  origin: https://<host>:3449

SockJS 帧规则（真值）：
- 服务端首帧 "o"（open）
- **发送侧必须 JSON 数组包裹**：ws.send(json.dumps([frame]))，NUL 由 json 转义为 \u0000
- **帧尾必须带 NUL \x00**
- 裸发文本帧 => 服务端 c[1007,""]（invalid payload）并关闭

STOMP CONNECT 帧（真值，实测返回 CONNECTED）：
  CONNECT\nx-v-token:<tk>\ncid:<cid>\nvid:<vid>\naid:null\naccept-version:1.1,1.0\nheart-beat:5000,5000\n\n\x00
  => a["CONNECTED\nversion:1.1\nheart-beat:25000,25000\nuser-name:<cid>_<vid>\n\n\u0000"]

★ 身份即 user-name = cid_vid，说明服务端把 STOMP principal 绑定到该会话身份。

## 2. 访客 Web 端 destinations（从 chunk_sac_sar.8fbad5.js 逆向，前端在用）
  SUBSCRIBE /user/topic/visitor/receive      (onmessage: JSON.parse(e.body))
  SUBSCRIBE /app/visitor/ack
  SEND      /app/visitor/ack   body=JSON      <- 访客侧是 JSON，不是 protobuf

## 3. 客服侧 protobuf 信封（common.proto）
  message UserMsgPack { int32 MsgId = 1; bytes Data = 3; }   // MsgId = CsMsgId 枚举值
  message MsgAuthReq  { string Token = 1; uint32 Account = 2; }   // cs_msg_staff.proto

命令表（cs_msg_id.proto，注释明确划分两端）：
  // ---------------- 客服 web 端 ----------------
  100 CsMsgAuthReq / 101 CsMsgAuthResp
  102 CsMsgGetStaffSessionsReq    <- **无参**，客服拉会话列表（值钱：可直读全量客服会话）
  103 CsMsgGetStaffSessionsResp { ErrorCode Code=1; repeated MemberSession Sessions=2; }
  106 CsMsgStaffGetHistoryMessage  / 107 resp { repeated ChatMessage MessageList=2 }
  122 CsMsgStaffSwitchMemberNotify / 123,124 GetSingleStaffSessions
  // ---------------- 玩家 app 端 ----------------
  200 CsMsgPlayerAuthReq / 201 resp / 204 GetStaffList / 208-223 订单与支付

手工 protobuf 编码（无需 protoc）：
  def varint(n): ...
  MsgAuthReq   = pb_bytes(1, token) + pb_uint(2, account)
  UserMsgPack  = pb_uint(1, msg_id) + pb_bytes(3, inner)
  # MsgId=100 -> 08 64 ; MsgId=102 -> 08 66 ; field3 tag -> 1a

STOMP SEND 帧：
  SEND\ndestination:<dest>\ncontent-type:application/octet-stream\ncontent-length:<len>\n\n<body>\x00

## 4. 结果（本轮实测，如实记录）
- CONNECT 可稳定复现成功（user-name=cid_vid）✔
- SEND 到 /app/staff/ack、/app/staff/auth、/app/staff、/app/staff/login、
  /app/staff/sessions、/app/sessions 以及对照 /app/visitor/ack —— **全静默**（STOMP 对未知
  destination 静默丢弃，无法区分"不存在"与"存在但无响应"），未取得 101 AuthResp / 103 Sessions ✘
- 见配套 POC「SockJS 文本通道无法承载 protobuf 二进制」：含 >0x7F 字节的 body 必被服务端拒
