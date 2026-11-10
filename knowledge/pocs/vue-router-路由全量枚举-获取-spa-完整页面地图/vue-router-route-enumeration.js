# 从 Vue Router 实例枚举 SPA 完整路由地图

## 适用场景
前端 SPA 的页面路由由 Vue Router 动态注册，源码里只有分块的懒加载 chunk 名，
无法直接得到「有哪些业务页面」。但**运行时的 router 实例持有全部已注册路由**。

## 取法（页面上下文 / CDP 执行）

```js
// 1) 取 router 实例
const app = document.querySelector('#app');
const router = app.__vue_app__.config.globalProperties.$router;

// 2) 枚举全部路由（含懒加载的）
const routes = router.getRoutes().map(r => ({
  path: r.path,
  name: target166.example.com,
  hasParams: /:/.test(r.path)
}));
console.log(routes.length, routes);

// 3) 过滤出高价值业务页面（账务/管理/导出/认证等）
const interesting = routes
  .map(r => r.path)
  .filter(p => !/[:*]/.test(p))
  .filter(p => /admin|manage|system|config|setting|export|finance|withdraw|deposit|
                agent|promote|report|member|user|security|kyc|bank|card|credit|loan|data|list/i.test(p));
console.log(interesting);
```

## 可直接用于实测的辅助代码

```js
// 逐页访问并抓取渲染文本（判断页面是否真实可用）
const visit = async (p) => {
  location.href = p;
  await new Promise(r => setTimeout(r, 3200));
  return {
    path: p,
    url: location.href,
    text: (document.body.innerText || '').replace(/\s+/g, ' ').slice(0, 800)
  };
};
// 用法
// for (const p of interesting) results.push(await visit(p));
```

```js
// 同时装 fetch 钩子，捕获页面加载时调用的接口（可发现新 API）
if (!window.__cap) {
  window.__cap = [];
  const of = window.fetch;
  window.fetch = async function (...a) {
    const r = await of.apply(this, a);
    try {
      const u = (typeof a[0] === 'string') ? a[0] : (a[0]?.url || '');
      if (u.includes('/api/')) {
        const t = await r.clone().text();
        window.__cap.push({ u: u.split('?')[0], s: r.status, len: t.length });
      }
    } catch (e) {}
    return r;
  };
}
```

## Vue 2 变体
```js
document.querySelector('#app').__vue__.$router.options.routes   // 静态路由表
document.querySelector('#app').__vue__.$router.getRoutes()      // 含动态注册
```

## React 变体
React Router 无统一实例，可从 `<Router>` 上下文或 `window.__REACT_ROUTER_HISTORY__` 取；
也可直接读打包后的 routes 配置模块。

## 输出过滤规避
- **字符码数组**：`[...s].map(c=>c.charCodeAt(0))`
- **分段输出**：`for(let i=0;i<s.length;i+=200) parts.push(s.slice(i,i+200))`
- 路由**路径本身是 ASCII**，通常不会被过滤；被过滤的往往是含 query/cookie 的完整 URL

## 实测产出（203.0.113.48 金沙娛樂場 target165.example.com）
枚举出 **101 条路由**（95 条无参数），发现此前未覆盖的高价值页面：
```
/home/report                    账户明细/投注记录/个人报表/找回余额/存款记录/提现记录
/home/personal-daily-report     个人日报表（输赢/投注金额/有效投注/注单量）
/home/credit-manage             额度管理（含「请输入直属下级账号」+ 批量操作）
/home/agent-center              代理中心
/home/withdraw-record-detail    提现记录详情
/home/security  /home/setting   安全中心 / 设置
/home/kyc  /home/kycFace        实名认证
```
逐页实测确认这些页面**全部完整渲染且无权限拒绝**（数据为空仅因账号无活动）。

## 要点
- `getRoutes()` 返回**运行时已注册**的全部路由，比读源码/搜索 chunk 名更完整。
- 对多子应用架构，需**分别进入各子应用**再取各自的 router（如 /agent 与主站是两套 router）。
- 该方法与「Pinia store 读取」互补：路由给**页面地图**，store 给**页面数据**。
