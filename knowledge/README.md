# knowledge/ — 打过的经验（知识库 + 提示词 + 运维手册）

这个目录是 ARTEX 的知识层：`knowledge.db` 是可检索台账，`pocs/` 是逐条 PoC 原文，
`prompts/` 是 Agent 角色提示词与 armor 预设，`ops/` 是发行包自带的运维手册。

## 怎么让 ARTEX 用上

在 `config.json` 里加一行（绝对路径），或设置环境变量 `ARTEX_KB_DIR`：

```json
{
  "knowledge_dir": "D:/path/to/ARTEX/knowledge"
}
```

启动后日志出现「知识库已加载: N 条」即生效；agent 开局会自动检索它（按 弹药 / 打法 / 负面情报 三档返回）。
当前 fork 需要先应用 `patches/0002-knowledge-runtime.patch`（`kb_search` / `kb_get` 运行时）；
不应用时本目录只是人读的文档库。

## 目录

| 路径 | 内容 |
|---|---|
| `knowledge.db` | 448 条台账（SQLite + FTS5 全文索引） |
| `pocs/` | 429 个条目目录（脚本 / md / http 报文） |
| `INDEX.md` | 448 条扁平索引 |
| `MAP.md` | 台账 ↔ 目录 一对一映射，含双侧未配对清单 |
| `prompts/` | 14 份角色提示词 + `armor/`（4 套预设 + `defaults.sql`） |
| `prompts/config/` | 部署默认值（settings 31 项 + 拦截规则 8 条） |
| `ops/` | 运维手册 6 篇（加固 / 部署 / 故障 / 升级 / Linux 选型 / 更新日志） |
| `absorb/` | 吸收过程记录与方法论（含跨产品经验、脱敏规则） |
| `tools/` | 武器库清单（197 件工具 + 8 个运行时聚合，二进制不入库）+ `kb_check.py` |
| `QA.md` | 语料体检报告（`python knowledge/tools/kb_check.py --write` 生成） |

## 内容

| 路径 | 数量 | 说明 |
|---|---|---|
| `knowledge.db` | 448 条 | 表 `poc` + FTS5；字段 category / component / kind / verified / tags / content |
| `pocs/` | 429 个目录 | 每条一个目录：`poc.sh` / `poc.py` / `*.md` / `*.http` |
| `prompts/` | 14 份 | planner / worker / verifier / retester / exploit / pentest / cloud / avevasion / payout_analyst / kb_curator / goals / mainagent / auto / worker.intranet |
| `prompts/armor/` | 4 套 | armor 预设 + `defaults.sql` + 摘要 |

## 条目怎么读

每条 POC 的标题就是结论，格式 `<产品>-<手法>-<判定>`。看到「判死 / 负面 / 已证否」的条目，
意思是**这条路已经走死，别重复踩** —— 负面情报和成功 PoC 一样值钱。

## 来源与脱敏

- 基线：ARTEX 0.4.2 发行包 `redteam/pocs`（411 条，已脱敏）
- 增量：老魏二改包 20 条独有条目（原始包**未脱敏**，入库前处理过）
- 规则：公网 IPv4 → `203.0.113.x`；非白名单域名 → `<target>`（目录名用 `tgt`）；
  `C:\Users\<name>\...` → `<HOME>`；`knowledge.db` 的 `path` 列置空
- **未收录**：`redteam/toolkit`（武器二进制）、`redteam/patches`（含 assets.db）、`nuclei-templates`

## 自检

```bash
python knowledge/tools/kb_check.py           # 台账/FTS 一致性 + 语料统计
python tools/scan_sensitive.py               # 公网 IP / 本机路径 / 客户数据文件
```

## 许可

内容来自本地发行包与实战记录，随仓库按 AGPL-3.0 发布前请自行确认第三方素材
（nuclei 模板、开源 PoC 片段）的原始许可。
