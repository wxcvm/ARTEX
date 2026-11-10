# 吸收索引（三期产出汇总）

本仓库从三个本地发行包吸收知识资产，全部**脱敏后**入库，二进制与客户数据不入库。明细见各期文档。

## 第一期：知识库 + 技能 + 提示词

| 产出 | 路径 | 规模 |
|---|---|---|
| 知识台账 | `knowledge/knowledge.db` | 448 条（FTS5） |
| POC 原文 | `knowledge/pocs/` | 429 目录 / 432 文件 |
| 条目索引 | `knowledge/INDEX.md` | 448 行 |
| 角色提示词 | `knowledge/prompts/*.md` | 14 份 |
| armor / 部署默认 | `knowledge/prompts/armor/` | 4 套预设 + `defaults.sql` |
| 技能 | `skills/` | 新增 43（合计 54） |
| 安全补丁 | `patches/0001-guard-hooks.patch` | planner/mainagent 挂 Guard |
| 过程记录 | [`knowledge/absorb/ABSORB-2026-11-10.md`](ABSORB-2026-11-10.md) | — |

## 第二期：运维文档 + 跨产品经验 + 对账 + 自检

| 产出 | 路径 | 说明 |
|---|---|---|
| 运维文档 6 篇 + 索引 | `knowledge/ops/` | 安全加固 / 部署 / 故障排查 / 升级 / Linux 选型 / 发行包更新日志 |
| dsh-redteam-mode 经验 | [`dsh-redteam-mode.md`](dsh-redteam-mode.md) | 刷分四条根因、口径一致性、并发闸门反例、提示词多层覆盖链 |
| 部署默认配置 | `knowledge/prompts/config/` | `settings.defaults.json`（31 项）+ `intercept-rules.json`（8 条） |
| 知识库对账 | `knowledge/MAP.md` | db 448 ↔ 目录 429 的一对一映射与双侧孤儿清单 |
| 敏感自检工具 | `tools/scan_sensitive.py` | 提交前跑一次，命中即退出码 1 |

## 三期小结：来源与取舍

| 来源包 | 取了 | 没取 | 原因 |
|---|---|---|---|
| `ARTEX-0.4.2-Windows64位.zip` | skills、deploy-defaults（提示词/armor/settings/拦截规则）、6 篇运维文档、POC 基线 411 | `app/data/tools`（1.23 GB 武器）、`pg/`、`nuclei-templates` | 二进制/大体积快照不进 git |
| `ARTEX-老魏二改带知识库-linux-amd64.tar.gz` | `knowledge.db`、POC 增量 20 条、`chains-pwn`、`defaults.sql`、patches 里的方法论 | `redteam/toolkit`（1.07 GB）、`redteam/patches` 原始目录（含 25 个 `assets.db`） | 客户资产快照与武器不入库 |
| `ARTEX-Friend-win.zip`（**下载截断**，987 MiB 无 EOCD） | 无（仅用于交叉校验 `knowledge.db` 哈希） | 全部 | 包不完整，不能作为安装源 |

### Friend 包差集结论（本期实测）

- 该包的 `redteam/pocs` 共 **377** 个目录名，与已入库的 429 条比对后：
- 归一化后看似"独有"的 **8 条，全部是已入库条目的未脱敏同版**（差异只在目录名里的主机段：`<原始靶标IP>` vs 已脱敏的 `203.0.113.10`、`hk000999.com` vs `tgt`），相似度 0.71–1.00；
- **结论：无新增条目**，不需要再并入任何 POC。
- 该包的 `redteam/engagements`（1717 条演练记录）不吸收：属于客户上下文，且包体截断无法完整校验。

## 第三期：知识库运行时 + 语料体检 + 武器清单

| 产出 | 路径 | 说明 |
|---|---|---|
| 知识库运行时补丁 | `patches/0002-knowledge-runtime.patch`（+350/−4，6 文件） | `knowledge_dir` 配置 + `kb/` 只读包 + `kb_search`/`kb_get` 两个工具 + 启动装配；FTS5 AND→OR→LIKE 三级降级 |
| 补丁说明 | `patches/0002-knowledge-runtime.README.md` | SQL 等价验证结果、启动日志、已验证 / 待 `go build` 清单 |
| 语料体检工具 | `knowledge/tools/kb_check.py` | 台账/FTS 一致性（P1）+ Python/Shell/重复/空文件（P2） |
| 体检报告 | `knowledge/QA.md` | 429 目录 / 432 文件 / 162 py（44 个是笔记或 py2，非可执行）；P1 无 |
| 武器清单 | `knowledge/tools/manifest.json` + `README.md` | 197 件工具 + 8 个运行时聚合（jdk11 / pgsql / Behinder / oneforall…），**二进制不入库** |

### 为什么需要 0002

本期开头核实过：**当前 fork 源码里 `knowledge_dir` / `kb_search` / `kb_get` / `knowledge.db` 零命中** ——
知识库只存在于 0.4.2 发行版二进制。也就是说 `knowledge/` 在这些改动之前对 fork 是**死资产**：
技能能从 `skills/` 加载，台账没人读。0002 把这条链补上（`config.KnowledgeDir()` → `kb.Open()` → `agent.SetKnowledge()` → 两个只读工具）。

### 本机能验证到的程度

- `git apply --check patches/0002-knowledge-runtime.patch` 通过
- 补丁里的两条 SQL 已在真库上跑过：`thinkphp` 5 条、`1panel 验证码` AND 3 条、`phpwind 写入` AND 5 条、`kb_get(1)` 正文 1686 字符
- 分档分布：ammo 356 / reference 60 / negative 31 / playbook 1
- 编译与运行验证需要本地 `go build`（本机无 Go 工具链）

## 未吸收清单（明确不做）

| 项 | 体量 | 原因 |
|---|---|---|
| `redteam/toolkit`、`app/data/tools` | ~2.3 GB | 武器二进制，杀软必报，不进 git |
| `redteam/patches` 原始目录 | 62 MB / 1031 文件 | 内含 25 个 `assets.db` 客户资产快照 |
| `nuclei-templates` | 30 MB / 14k 模板 | 用官方仓库更新，不存快照 |
| `knowledge.db` 的 `engagement_id` 关联 | — | 字段保留（404 条有值）但 `engagement_name`/`asset_target` 已被发行方清空；如要彻底匿名可整列置空 |

## 提交前自检

```bash
python tools/scan_sensitive.py        # 期望：[ok] 未发现…，退出码 0
git status --short                    # 只应有 knowledge/ docs/ skills/ tools/ 与 README/compose/update 的改动
```

## 脱敏规则（两期一致）

| 项 | 规则 |
|---|---|
| 公网 IPv4 | → `203.0.113.N`（19 条映射表，见 `knowledge/absorb/ABSORB-2026-11-10.md`） |
| 非白名单域名 | → `<target>`（目录名里用 `tgt`） |
| 本机路径 | `C:\Users\<name>\…` → `<HOME>` |
| `knowledge.db.path` | 整列置空（原指向作者本机） |
| `assets.db` / 武器 / 模板 | 不入库 |

