# Vite SPA 多子应用 API 面递归提取

## 适用场景
一个域名/IP 下挂多个**独立构建**的前端子应用（主站 + 代理后台 + 活动站等），
每个子应用有各自的 API 面。只分析主站 bundle 会漏掉大量接口。

## 识别独立子应用的三个信号
1. **构建产物哈希不同**：`/agent` 用 `commonChunk.GQwXmA27.css`，
   `/gossip` 用 `commonChunk.CDp6Q3sH.css`，主站用 `commonChunk.CMoJwNYm.js`
   → 哈希不同 = 独立打包 = 独立的懒加载 chunk 集合
2. **入口 chunk 命名带业务语义**：`index.admin-agent.WOmu9FLX.js`
3. **引用不同的 manifest**：读入口 chunk 里的
   `window.__MANIFEST_ASSETS = 'assets/theme-0/manifest.be2a1d70.json'`

## 核心步骤

### 1) 从主 bundle 提取顶层路由
```python
tops = set(re.findall(r'["\'`](/(?!api/)[A-Za-z][A-Za-z0-9_\-/]{1,60})["\'`]', main_bundle))
```
逐个 GET，返回 200 且带 `<title>` 的就是真实子应用入口。

### 2) 读入口 chunk 拿 manifest
```python
entry = re.search(r'import"\./(index\.[A-Za-z0-9_\-]+\.js)"', html)
entry_js = GET("/assets/theme-0/" + target163.example.com(1))
mf = re.search(r"__MANIFEST_ASSETS\s*=\s*'([^']+)'", entry_js)
manifest = json.loads(GET("/" + target164.example.com(1)))
```

### 3) **manifest 里的 chunk 名本身就是功能清单**
这是最大价值点 —— 不用下载就能知道后台有哪些功能：
```
AgentApplication / CommissionIndex / ExportIndex
ClubAgentCommissionIndex / ClubAgentPerformanceIndex
DirectBettingIndex / DirectDataIndex / DirectFinanceIndex / DirectTransferIndex
ExportIndex / ExportButtonIndex / EffectiveBetIndex
Devtools / DevtoolsIndex          ← 后台里的开发者工具
DiscountCodeIndex / DepositRewardsIndex / DistributionDetailsIndex
```

### 4) 按业务关键字筛 chunk 下载，正则提 API
```python
BIZ_KW = ["Agent","Commission","Export","Finance","Data","Transfer","Bet",
          "Withdraw","Deposit","Member","User","Report","Club","Distribution"]
apis = re.findall(r'["\'`](/api/[A-Za-z0-9_\-/]{2,90})["\'`]', chunk_src)
```

## 实测产出（203.0.113.48）
- 发现 4 个子应用：`/agent`（代理后台）、`/gossip`、`/promote-game`、`/pc`
- `/agent` manifest：**531 项资源 / 322 个功能 chunk**
- 提取出 **24 个主站清单外的 API**，含高危功能：
  - `/api/agent/promote/batchCreateAccount`（批量建号）
  - `/api/agent/promote/editUserRate` / `setDefaultUserRate`（改代理费率）
  - `/api/agent/promote/reportPc/export/{directInfo,directFin,directOrder,transfer}`（数据导出）
  - `/api/agent/promote/reportPc/{directReportPcV1,directFinV1,userDayDepositList,userDayWithdrawList}`（会员充值/提现/业绩报表）

## 踩坑提醒
- **"未鉴权返回 200 + 加密密文" ≠ 数据泄露**。必须做「无鉴权 vs 携带有效会话」的
  **响应长度对照**：若两者逐字节相同且 raw 长度为 16 的倍数，说明返回的是
  「无权限」的加密错误信封，**没有实际数据**，不能当成果记分。
- 写/导出类接口常返回 `41060 参数异常/解析包异常`（参数层拒绝），需构造加密请求体。