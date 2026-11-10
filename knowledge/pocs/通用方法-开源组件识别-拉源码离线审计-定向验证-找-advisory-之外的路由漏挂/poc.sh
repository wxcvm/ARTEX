# 开源组件识别 → 拉源码离线审计 → 定向验证

适用于任何**用开源产品搭的靶标**（New API / one-api、Nacos、Grafana、Jenkins、各类 OA、AI 网关…）。价值在于：**advisory 之外的问题 + 只有源码能回答的"哪些路由没挂认证"**，且未认证可达的接口可以直接在目标上验证 ⇒ 能直接得分。

## 第 0 步：先识别组件与精确版本（不要猜）
- 响应头 / 前端资源哈希 / 错误页 / `/api/status` 之类的指数字段。
- 例：New API 的 `X-New-Api-Version: target1.example.com`（`curl -sS -k -D - -o /dev/null https://<host>/api/status | grep -i x-new-api-version`）。

## 第 1 步：拿到对应版本源码（**不要用 git clone，改用 codeload tarball**）
```bash
# 先确认 tag 存在（rc.x 这类中间版本常常没有 tag）
git ls-remote --tags https://github.com/<org>/<repo>.git | grep '<目标tag>'
# 下载 tarball（比 clone 快一个数量级；实测 3.6MB / 1.3s）
curl -sSL -m 300 -o /tmp/src.tar.gz "https://codeload.github.com/<org>/<repo>/tar.gz/refs/tags/<tag>"
tar -xzf /tmp/src.tar.gz -C /tmp/
```
> **踩坑**：`git clone --depth 1` 在受限链路上经常卡在 `fetch-pack: unexpected disconnect` / `early EOF`；codeload 是单次 HTTP，可控得多。
> **版本对齐纪律**：若无该版本 tag，用最接近的 commit/相邻 tag，并在结论里**注明版本偏差**，不要把别版本代码当目标实现。

## 第 2 步：审计路由注册与中间件链（核心）
```bash
# 一次性拉出「路由方法 + 中间件链」，靠缩进/Use 位置判分组权限
grep -nE '(^\s*(self|admin|router|Group|Route)|\.Use\(|GET|POST|PUT|DELETE|PATCH)' router/*.go
```
判读要点：
1. **`group.Use(auth())` 之前注册的路由不受该中间件约束** —— 这是最高频的漏挂模式。gin 中 `g.METHOD(...)` 立即绑定，之后再加 `.Use()` 只影响后续注册。
   - 实例：`relayMjRouter.GET("/image/:id", ...)` 写在 `.Use(middleware.TokenAuth(), ...)` **之前** ⇒ 免认证。
2. **同一 group 被注册到多个前缀（含通配段）时漏挂面会放大**。instance：`/mj` 与 `/:mode/mj` 都注册同一套路由 ⇒ `/<任意段>/mj/image/:id` 也免认证。
3. **区分"设计默认公开"与"鉴权绕过"**：若中间件是 `TryUserAuth()` 或"按配置/选项回退"（如读一个空的配置项 → 回退 `RequireAuth:false`），那是**设计公开**，定级时必须说明，别当绕过报。
4. **注意 `XxxOrYyyAuth()` 这类复合中间件**（如 `TokenOrUserAuth()`）：它接受多种凭据 ⇒ 判"是否未认证可达"**必须做有无凭据对照**，只看单次请求很容易误报（一次 404 可能是"通过了认证但资源不存在"）。
5. 顺手核对**已知实测分界**：用源码解释"为什么 A 匿名可读、B 却 401"，然后**推导被枚举漏掉的同类路由**。

## 第 3 步：定向验证（每个候选 ≤2 请求，串行、间隔 ≥1s、走可信出口）
- **判据必须对照**：① 存在且未认证可达 → 返回真实数据/业务错误体；② 受保护 → 401/403；③ 路由不存在 → 404 `page not found`。三者响应形态要事先用已知样本校准。
- 别一次批量打——命中一处就够，过量请求只会换来封禁。

## 第 4 步：源码里的“隐藏捷径与假前提”
- 列表接口脱敏 **≠** 拿不到明文：找有没有**专用取明文接口**（New API 的 `POST /api/token/:id/key` → `GetFullKey()`）。
- “凭据正确但被拒”要读校验函数定位真因：New API 的 `model/target35.example.com:241` 把 `!UnlimitedQuota && RemainQuota<=0` 判成 `ErrTokenInvalid` ⇒ 看起来像 token 无效，其实只是额度为 0。
- 顺着“能改自己”的写接口找**越权/提权**：`PUT /api/token/` 仅挂 UserAuth，且 `unlimited_quota=true` 时跳过额度校验 ⇒ 普通用户自设无限额度（计费绕过）。
- 汇报纪律：**源码审出但未在目标验证的一律 `status=candidate` + `source=源码审计` + evidence 写清文件:行号与函数**；已实测的才 `confirmed`，并配原始请求。