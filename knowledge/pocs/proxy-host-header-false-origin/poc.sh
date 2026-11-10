# 方法论陷阱：用本地 HTTP 代理做 Host 头 vhost 探测会造出"命中源站"的假象

## 现象（我踩了这个坑，并把错误结论写进了报告）
目标：找一个 Cloudflare 后面的站点的真实源站。
做法：对候选 IP 发 `GET http://<IP>:8080/login` 并带 `Host: <被 CF 代理的域名>`，
      期望返回该域名的页面（= 命中源站 vhost）。
      **经本机 HTTP 代理（127.0.0.1:7892，Clash/V2Ray 类）执行时 → 5/5 台都返回 200 + 该域名的登录页**
      → 我据此判定"源站定位成功"。

## 真相
该代理**按 Host 头/域名做 DNS 解析与分流**，直接把请求送到了**真实域名（经 Cloudflare）**，
目标 IP 被完全忽略。铁证是响应头：
```
Server: cloudflare
Cf-Cache-Status: DYNAMIC
Cf-Ray: a4210cf089b6b486-SIN
Alt-Svc: h3=":443"; ma=86400
```
这些头只可能来自 Cloudflare 边缘 —— 一个直连的阿里云 IP:8080 不可能返回它们。

## 无代理复核（决定性对照）
在 VPS 上用 `http.client` 直连（不经任何代理），对同样 5 台 IP 各试 4 个 Host：
```
Host=(默认) / adminapi.<domain> / app.<domain> / admin.<domain>
  /       -> 四种 Host 全部返回**完全相同**的响应（87 字节 JSON）
  /login  -> 一律 404
```
→ 这些实例忽略 Host、不承载目标 vhost，**根本不是源站**。

## 判定规则（务必遵守）
1. **任何 vhost / 源站判定，必须用"无代理直连"复核一次**；带 Host 头的探测尤其危险。
2. 看响应头做交叉验证：若出现 `cf-ray` / `cf-cache-status` / `server: cloudflare` / `alt-svc: h3`，
   说明你很可能打到了 Cloudflare 边缘而不是源站。
3. 对照组要包含"**不同 Host 是否产生不同响应**"——若多个 Host 响应字节完全一致，
   说明目标忽略 Host（或你根本没打到 vhost 分层的网关）。
4. curl 复现（确保不走系统/环境代理）：
   `curl --noproxy '*' -H 'Host: target.example.com' http://SOURCE_IP:8080/login -D -`
   或直接在跳板机（VPS）上跑，避免本机代理规则干扰。
