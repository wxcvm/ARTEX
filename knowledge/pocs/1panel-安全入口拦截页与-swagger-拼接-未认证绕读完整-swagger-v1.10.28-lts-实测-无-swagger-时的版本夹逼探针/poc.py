## 现象（可复用）
启用「安全入口（Security Entrance）」的 1Panel 上，`GET /1panel/swagger/doc.json` 的响应体是：
`[1410 字节安全入口拦截页 HTML] + [完整 Swagger 2.0 JSON]` **拼接**（中间件写入了拦截页却没有中断 handler）。

实测字节账（v1.10.28-LTS）：
- 总长 813736 = 1410（拦截页）+ 812326（纯 swagger）
- 纯 swagger 的 md5：`6b6b816a6caee60eea43193ee0fc39e8`（basePath=`/api/v1`，405 paths / 415 路由）
- 首页 `GET /` 仍返回 1410 字节「暂时无法访问」⇒ 仅看首页会误判为"完全不可达"

⇒ **只要切掉前 1410 字节即可未认证获得完整路由表**，无需口令。

## 切片方法
```python
b = open('/tmp/raw.bin','rb').read()
i = b.find(b'{\n    "schemes"')   # 或 b.find(b'{"swagger"')
js = b[i:]
```
或 `tail -c +1411 raw.bin > swagger.json`。

## 对照事实（判"到底能不能拿 swagger"）
| 实例配置 | `/1panel/swagger/doc.json` 行为 |
|---|---|
| 启用安全入口 | 813736B = 1410B 拦截页 + 812326B swagger（**可绕读**） |
| 未启用安全入口 | 返回 `{"code":401,"message":"用户未登录"}` 52B（swagger 路由被 JWT 中间件接管，**拿不到**） |
⇒ 「没启用安全入口」不等于「能拿 swagger」，别搞反。

## 拿到 swagger 后怎么判版本
把纯 swagger 与官方各 tag 的 `backend/router/*.go` 做路由集合差集（方法+路径）：
- example.com：目标独有 11、tag 独有 25
- v1.10.29/30/31/32-lts：目标独有仍 11，但 tag 独有陡增到 35（多出 10 条 MCP 路由：`POST /mcp/search`、`POST /mcp/server`、`POST /ai/mcp/*`）
- v2.x：basePath 变 `/api/v2`，直接被排除

## 拿不到 swagger 时的替代：路由存在性探针（JSON vs SPA 对照）
1. 先建立对照：请求一个**确定不存在**的路径，看返回什么。1Panel 未启用安全入口时返回 3448~3449B 的 SPA HTML（`<!DOCTYPE html>`）；启用安全入口时返回 1410B 拦截页。
2. 再探增删路由，用"返回 1Panel JSON（401/400/参数错误）＝路由存在"来夹逼：
   - `GET /api/v1/databases/redis/exec`（v1.10.28 新增）→ 存在 ⇒ ≥v1.10.28
   - `GET /1panel/swagger/example.com`（v1.10.28 引入）→ 存在 ⇒ ≥v1.10.28
   - `GET /ai/gpu/load`（v1.10.28 移除）→ 不存在 ⇒ ≥v1.10.28
   - `POST /api/v1/mcp/search`、`POST /api/v1/ai/mcp/search`（v1.10.29 新增）→ 不存在 ⇒ <v1.10.29
   - `GET /api/v2/core/auth/captcha` → 不存在（SPA）⇒ V1 线
   **注意方法要对**（`/api/v1/dashboard/base/os` 只有 GET；用 POST 探会得到 SPA 而误判为不存在）。