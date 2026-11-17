---
name: payout-hunt
description: 支付/资金盘/博彩类目标面的批量发现与定级：三引擎 + favicon 指纹扩站群、六道 Gate 逐条判真、P0-P3 排序派单，产出可直接交给 payout-flow-analysis 的候选清单
whenToUse: 手上只有几个种子域名或关键词、需要从公开资产面里筛出"后台面可达且有真实资金面"的支付/资金盘类站点时；或已有一批候选、需要按可打性排序并剔除误报时
role: recon
enabled: true
---

# 目标发现与定级

一句话：**关键词找方向，指纹找同胞，Gate 判真伪，回显定生死。** 单站作业是浪费——见一个就要拿到它一窝。

本技能只做**发现与判定**：不猜后台密码、不提交验证请求、不改任何订单状态。定级到 P0/P1 之后交给 `payout-flow-analysis` 还原出款链路。

## 铁律（破一条，整份结论作废）

1. **存活不等于有肉**：首页 200 只说明域名活着。没有后台面、没有资金面，就是线索不是战果。
2. **结论必须有回显**：说"有后台"，就要给出那条 200 的登录页 URL 与状态码；说"托管资金"，就要给出那个链上地址。
3. **单引擎命中整体降一级**：只有 FOFA 命中的不叫交叉验证，两级以上引擎互斥才算独立证据。
4. **宁可漏，不误报**：一条假的 P0 会让后面十条真的没人看。
5. **分型先于评分**：先判 T1/T2/T3/T4，再谈分；型判错，分毫无意义。
6. **关键词永远不够**：靠品牌词搜出来的七成是导航站和过期站，指纹才是主体。
7. **一源一指纹，扩一批**：拿到一个活样本的源码指纹，就是拿到它背后整套站群的名单。
8. **回显可复现**：报告里每条都要能被别人用同一串命令原样跑出来。

时效：域名到期、站被关、后台被补都以天计。**超过 30 天的结论标 `stale`**，重跑前不得直接引用。

## 目标分型（决定后面二十步怎么走）

| 型 | 长相 | 特征信号 | 价值 |
|---|---|---|---|
| **T1 信用盘 / 棋牌** | 会员中心、额度、上下分 | 自有 `/api/` 源站、多家 vendor 接口、客服在线、有充值地址表 | 最高，常带热钱包与归集配置 |
| **T2 资金盘 / 量化 / 抢单 / 分红** | 讲收益、讲静态动态、有团队层级 | 记账型余额、改余额面板、提现走审核、域名寿命 3–8 个月 | 高，存货量大但单站浅 |
| **T3 空壳 / 仿站 / 快照** | 有皮无骨，点了没反应 | 静态 HTML、接口 404、logo 盗用、`/api/` 全空 | 低，略过 |
| **T4 招商 / 导航 / 测评** | 一堆外链、排名、点进去跳走 | 纯跳转、无自身业务接口、多为 SEO 站群 | 最低，**只当种子源** |

判型口诀：点"登录"看有没有真接口；点"充值"看有没有真地址。两者都没有 → T3/T4。

## 三引擎与指纹

| 引擎 | 端点 | 凭据 |
|---|---|---|
| FOFA | `https://fofa.info/api/v1/search/all`（`FOFA_ENDPOINT` 可覆盖） | `FOFA_KEY` |
| Shodan | `https://api.shodan.io/shodan/host/search` | `SHODAN_KEY` |
| Censys | `https://search.censys.io/api/v2/hosts/search`（POST + Basic） | `CENSYS_ID` / `CENSYS_SECRET` |

三家语法各不相同，**不要手搓 curl**，用 `scripts/bot_hunt.py`——它负责翻译、限速、退避、去重、落盘。缺哪家 key 就自动跳过哪家。

指纹扩线三条，按威力排序：

- `icon_hash` 相同 → 同产线全部站点（一次几十上百，最狠）
- `body="<某静态文件名>"` → 同版本源码（比品牌词准得多）
- 证书 SAN / 同 C 段 / 同 `server` + 同标题模板 → 同运维

取指纹与扩线：

```bash
# 取样本指纹（favicon mmh3 / 静态资源路径 / 接口前缀）
python scripts/bot_hunt.py --from-sample https://TARGET -o out/fp.json

# 用指纹铺站群
python scripts/bot_hunt.py --engine fofa,shodan --icon-hash <fp.json 里的 hash> -o out/squad.csv
python scripts/bot_hunt.py --engine fofa,shodan --body "/static/app.<hash>.js" -o out/squad2.csv

# 关键词铺面（引子，不是主体）
python scripts/bot_hunt.py --engine fofa,shodan,censys --keyword --geo sea --species money,usdt,agent -o out/breadth.csv
```

坑：FOFA/Shodan 免费额度按次扣，`--pages` 别开大；`--qps` 默认 1.0，被 429 就降到 0.3；Censys 的 cursor 分页与 FOFA 页号语义不同，**跨引擎别混用页码**。盾站（CF 挑战页）取到的 favicon 是云厂商的，换内页重取或人工看源码。

## 六道 Gate（顺序不能换，任一道不过直接 discard）

**Gate 1 存活**
- 过：首页 200（或 301/302 到站内）、body > 3KB、有 `<title>`
- 拒：502/503/522、CF 挑战页（`cf-challenge` / `Just a moment`）、空白页、停放页（`sedoparking` / `hugedomains` / `GoDaddy`）、`Index of /` 空目录
- 盾挡住**不算死**：记 `shielded` 进人工队列，过盾后补 body

**Gate 2 品类命中**（多信号互斥，单一关键词不放行）
- 过（任一）：① 命中 ≥2 个 vendor 接口（`/api/game/` `/api/ag/` `/api/mg/` `/api/pt/` `/api/pg/` `/api/bbin/` `/api/ky/` `/api/db/` `/api/og/` `/api/cq9/`）；② 关键词 ≥3 个且分属 ≥2 个语义簇（品牌/品类/资金）；③ 页脚带系统名 + 有 `/api/` 源站
- 拒：只出现一个"博彩"字样、只有导航外链、只有一张招商海报

**Gate 3 后台面**
- 过（任一）：① `/admin` 类返回 200 且 body 含 `password`/`登录`/`验证码`；② `/api-docs`、`/v2/api-docs`、`/swagger-ui` 返回 200 且是 JSON/UI；③ `/.env`、`/.git/config` 返回 200 且有内容；④ 代理零门槛注册可用（`/agent/register`、`/api/game/agent/register` 返回 200）
- 拒：全部 404、302 回首页、200 但是首页 HTML（伪路由）
- **代理注册这条路单独记星**：注册成功拿到 JWT 就是半个后台，命中直接记 P0

**Gate 4 资金面**
- 过（任一）：① 页面出现 ≥1 个合法 TRON 地址（`T` + 33 位 base58）；② 有充值/提现页且含地址表或银行卡收款信息；③ `/pay/`、`/cashier/`、`/recharge/`、`/withdraw/` 返回 200
- 拒：只写"支持 USDT"但一个地址都抓不到
- 抓到地址立刻查两件事：余额（`apilist.tronscanapi.com/api/account?address=`）、是否命中已知归集地址。**有余额的地址是 P0 的硬证据**

**Gate 5 真托管**
- 过：首页有自己的 `/api/` 请求、会员登录走本站接口、换 IP 访问行为一致
- 拒：整站 iframe 套别人、点任何按钮都 `window.open` 到第三方、只有静态壳——这类是贴牌站，权限在后面那家手里，打它没意义

**Gate 6 分级**：按下面表落 P0–P3 或 discard，证据（URL + 状态码 + 片段）写进行内，写不出证据的一律降级。

## 分级与派单

| 级 | 条件 | 动作 | 后续路径 |
|---|---|---|---|
| **P0** | Gate1–5 全过 + 抓到有余额的地址 + 后台面可达 | 立即打后台 | 私钥直转 / 改归集 |
| **P1** | Gate1–5 全过、后台面可达、地址无余额或未抓到 | 排队 | 改余额提现 |
| **P2** | 品类真、有资金面，但后台面未找到 | 补字典 / 子域深挖 | — |
| **P3** | 只过 Gate1–2 | 监控，择日重扫 | — |
| **discard** | 任一道 fail / 命中黑名单 | 归档，不占资源 | — |

排序：`P0 > P1 > P2 > P3`，同级按**地址余额降序**，再按**域名年龄升序**（新盘后台更松）。

## 黑名单（命中即 discard，别浪费时间判）

| 特征 | 是什么 |
|---|---|
| `sedoparking` / `hugedomains` / `dan.com` / `GoDaddy` | 停放域名 |
| `title="Index of /"` 且无 `api` | 裸目录 |
| 页面只有外链列表 / 排行榜 / "十大平台" | 导航测评站 |
| `iframe` 占满 + 点按 `window.open` | 贴牌套壳 |
| whois 剩余 < 30 天 + 首页打不开 | 到期弃站 |
| 全站 `cf-challenge` 且 3 次都不过 | 盾死，转人工 |
| 首页 200 但所有 `/api/` 全 404 | 静态壳 |
| 同一 HTML 里 vendor 接口 ≥5 家但无自有 `/api/` | 聚合跳转页 |
| `server: aiohttp` + 固定假 JSON | 蜜罐 |
| 中文站但 `country="CN"` 且已备案 | 大陆站，优先级最低 |

## 脚本

| 脚本 | 职责 | 边界 |
|---|---|---|
| `scripts/bot_hunt.py` | 三引擎检索、指纹扩线、去重落盘 | **只找**，不判定、不发探测请求 |
| `scripts/bot_verify.py` | 六道 Gate 逐条判定、抽地址、分级 | **只判**，探测只取状态码与标题特征，不发验证性 payload |

一条命令跑全链：

```bash
python scripts/bot_hunt.py --batch https://SEED1,https://SEED2 --engine fofa,shodan \
        --geo sea --conc 30 --verify -o out/run_01
```

分步版（要精细控制时）：

```bash
python scripts/bot_hunt.py --from-sample https://SEED -o out/fp.json
python scripts/bot_hunt.py --engine fofa,shodan --icon-hash <hash> -o out/squad.csv
python scripts/bot_hunt.py --engine fofa,shodan,censys --keyword --geo sea --species money,usdt,agent -o out/breadth.csv
python scripts/bot_verify.py -i out/candidates.txt -o out/probe.jsonl -c 40 --timeout 8
# 对 P0/P1 补后台面字典
python scripts/bot_verify.py -i out/top.txt -o out/backdoor.jsonl -c 15 --stage backdoor
```

## 产出契约

```
hunt/
  report.md          # 结论，按分级排序
  raw/               # 各引擎原始返回，不加工
  candidates.csv     # 合并去重后的候选
  probe.jsonl        # Gate1-5 判定流水（一行一条）
  p0.csv             # 只有 P0/P1
  discarded.csv      # 误报留档，用于复盘
  addresses.txt      # 全部抓到的链上地址 + 余额
  evidence/          # 关键响应片段
```

判定流水字段（严格 JSONL，一行一条，不加解释）：

```json
{"url":"...","gate1":"pass","gate2":"pass","gate3":"fail","gate4":"pass","gate5":"pass","type":"T1","level":"P2","evidence":["/admin/login 404"],"usdt":["T..."],"next":"path_C","note":""}
```

报告规则：不堆 emoji、不写编造的统计；`/admin` 这类路径原样写；地址不脱敏。

## 环境与限速

- 并发：探测 40 线程起；单站并发 ≤3
- 限速：FOFA 单账号 ≤1 QPS，429 就退避到 0.3
- 出口：准备 3 个以上出口 IP，同一 IP 连扫会被 403 记指纹
- 超时：连接 5s / 读 8s，**超时按 fail 处理，不许当 pass**
- 落盘：数据放工作目录下的 `out/`，不放系统盘

## 交给谁

定级到 P0/P1 之后，出款链路还原走 `payout-flow-analysis`（通道表 / 凭据 / dry-run 脚本）；要拿后台权限后的横向动作走 `chains-web`；代理越权那条线走 `web-attack-methods` 的越权段。

## 自检（出货前）

- [ ] 三把引擎 key 都在环境变量里，没硬编码
- [ ] 候选来自 ≥2 个引擎，已去重
- [ ] 至少 3 个活样本的 favicon hash 已算出，站群已扩
- [ ] 每条 P0 都有 `/admin` 的 200 回显
- [ ] 每条 P0 都有真实链上地址（不是"支持 USDT"文案）
- [ ] 黑名单全部过了一遍
- [ ] 单引擎来源的已降级
- [ ] 结论未超过 30 天，或已标 `stale`
