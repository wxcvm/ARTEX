# 跨 vhost 同源 XSS 可执行性判定矩阵（多域名站群通用）

## 要判的问题
会员端上传点落在 `public/uploads/...`，但**只有与后台/客服/代理同一个 document root 的 vhost** 才能把该文件当同源页面执行 JS。
「共享同一台机器 / 同一份代码」≠「共享同一 web root」。必须逐个 vhost 实测。

## 判定命令（本机 Clash fake-ip 环境，禁用系统 DNS）
```python
# CONNECT 真实源站 IP:443，只改 Host 头；SNI 实测不影响路由
raw_req("<任意SNI>", "/uploads/qrcode/YYYYMMDD/<name>.<ext>", "GET",
        connect_ip="<源站IP>", host_header="<待判 vhost>")
```
判据：
| 结果 | 结论 |
|---|---|
| `200` + `text/html`（.html）或 `image/svg+xml`（.svg） | ✅ 该域下 XSS 可执行（同源页面） |
| `404` + `nginx/1.18.0` 默认页 | ❌ 该 vhost 的 root 不含 public/uploads（≠同源） |
| `200` + SPA 的 example.com（固定 763B） | ❌ 命中了 try_files 回落到前端 SPA，不是我们的文件 |
| 连接重置 / TLS EOF / 无响应 | ❌ 该 vhost 已被摘除或对我们不可达 |

## 关键补充判据（避免被「同一份代码」误导）
**再测一条后台/应用路径做交叉**：`/app/admin/api/game-category/index`、`/app/admin/account/info`。
- 站点 vhost 上 `/app/admin/*` = nginx 默认 404 ⇒ 该域**根本没有后台应用**，同源 XSS 也驱动不了后台。
- 只有当「/uploads 可执行」与「/app/admin/* 有应用响应」出现在**同一个 Host** 上，才是可用的同源组合。

## 实配结论（203.0.113.88，2026-10-12）
- `/uploads/...` 200+可执行 JS 的 vhost（9 个）：example.com、example.com、example.com、example.com、example.com、example.com、example.com、example.com、example.com（www.example.com / example.com 本次单发无响应，待复核）
- example.com / example.com / api.example.com：全路径 nginx 404（vhost 已被运营方摘除，真实 Chrome 直接「意外终止了连接」）
- **所有站点 vhost 上 `/app/admin/*` 全部 nginx 404** ⇒ 不存在「既有我们文件、又能调后台 API」的同源点 ⇒ 该站群不能靠这条链拿后台权限
- 会话存放：后台 Vue 面板 token 在 `localStorage['user'].accessToken`；代理后台 token 在 `localStorage['agent_token']`；GoFly 客服 `/` 与 `/login` **不下发任何 Set-Cookie**（无域级 Cookie）⇒ **跨子域偷会话不可能**，各 origin 的 token 互不相通
