# 3X-UI / x-ui 面板 /assets/ 匿名目录列举

## 0. 前置

- 目标是 3X-UI / x-ui 面板的**明文 HTTP**端口，默认常用 2096（2053/2087 顺带探，但本案例仅 2096 命中）。
- 不需要认证、不需要 Host 路由、不需要 webBasePath。
- 本机有代理环境时**必须** `--noproxy '*'`，否则回包会被代理替换成代理错误页，全部结论失真。

## 1. 主判据：匿名目录列表（一条命令定生死）

```bash
curl.exe --noproxy '*' -sS -D - -o assets.txt \
     http://192.0.2.18:2096/assets/
```

成立时：

```
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 9773

<!doctype html>
<meta name="viewport" content="width=device-width">
<pre>
<a href="AddClientsToGroupModal-B855RnDJ.js">AddClientsToGroupModal-B855RnDJ.js</a>
<a href="AmneziaWGLogModal-BMaeu7oM.css">AmneziaWGLogModal-BMaeu7oM.css</a>
...
```

## 2. 两条 404 基线（证明存在独立静态层，并给 webBasePath 定性）

```bash
# 根路由层：18 字节 text/plain
curl.exe --noproxy '*' -sS -D - http://192.0.2.18:2096/zz9x7q2k
# -> HTTP/1.1 404 Not Found / Content-Type: text/plain / Content-Length: 18
#    404 page not found

# 静态层：0 字节，空 body
curl.exe --noproxy '*' -sS -D - http://192.0.2.18:2096/assets/nonexist.js
# -> HTTP/1.1 404 Not Found / Content-Length: 0
```

两种 404 字节数/内容体完全不同 ⇒ 2096 上的 `/assets/*` 由独立静态文件 handler 服务，
根路由是另一套 SPA/API handler。

顺带：`GET /assets`（无尾斜杠）→ `301`，是标准目录归一化，不是拦截。

## 3. 产物解读

列举 body 本身就是一份完整 chunk 清单，指纹极强：

- `XrayPage-*.js`、`InboundsPage-*.js`、`ClientsPage-*.js`、`HostsPage-*.js`、`NodesPage-*.js`
- `AmneziaWGLogModal-*`、`amneziawg-obfuscation-*.js` ⇒ 带 AmneziaWG 的分支
- `ApiDocsPage-*`、`vendor-swagger-*` ⇒ 面板自带 API 文档页
- 大量 `*-zh-CN.json` 等 i18n 包

⇒ 结论：3X-UI（x-ui）面板，且为含 Xray + AmneziaWG 的较新分支。

## 4. 往下走

```bash
# 取入口 chunk，里面带 window.X_UI_BASE_PATH 与接口调用
curl.exe --noproxy '*' -sS -o idx.js \
     http://192.0.2.18:2096/assets/IndexPage-tZHM2Ebk.js
grep -oE 'window\.X_UI_BASE_PATH[^,;]*' idx.js
grep -oE '/panel/api/[a-zA-Z0-9/_-]+' idx.js
```

实测 idx.js 中存在 `window.X_UI_BASE_PATH||'/')+"assets/..."` 的 chunk 映射表与
`/panel/api/inbounds/options` 等接口线索 —— 随机 webBasePath 下，根路径 `/panel/api/...`
一律命中 18 字节 404，需要先用 bundle 里的 BASE_PATH 拼真实前缀。

## 5. 失败排错

| 现象 | 原因 |
|---|---|
| 根路径 200 但 body 是十几 KB HTML 模板 | 前置 nginx 欢迎页/反代，不是面板 |
| 所有请求回代理错误页 | 漏 `--noproxy '*'` |
| `/assets/` 403 | 面板关了目录列表、或反代加了 autoindex 限制 |
| 2053/2087 全 timeout | 只开了 2096，或只对外映射了 2096 |
| TLS 端口无响应 | 换 `https://` + `--insecure` |

## 6. 定级

仅泄露前端 chunk 名与少量接口线索，无凭据、无配置 → low（信息面收益，用于指纹与后续 bundle 逆向）。
