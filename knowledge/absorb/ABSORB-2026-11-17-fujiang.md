# 吸收记录：富江 0.20.13 外源索引（2026-11-17）

## 来源

| 项 | 值 |
|---|---|
| 上游仓库 | [`FJNB5/fujiang-pojia`](https://github.com/FJNB5/fujiang-pojia) |
| 固定提交 | `cf086fc329011c6f8eeb3b67f0adbcc3f005debb` |
| 上游目录 | `plugin/knowledge/`（176 个 Markdown：105 组件/CVE + 59 灰产 + 3 方法论 + 9 分类入口） |
| 中间结构 | dsh-purge `f806b96` 的 `lib/redteam/cve-catalog/*.json`（已只取元数据） |
| 取用部分 | 105 条组件/CVE + 59 条灰产手法 = **164 条** |

上游的 3 篇方法论与 9 个分类入口 `SKILL.md` 未取：它们是目录导航，未进入本批索引。

**来源限制（2026-11-17）：** 声明的原仓库目前无法通过 GitHub API 读取（404），第三方公开索引也没有提供可核实的许可证。保留固定提交与路径只为可追溯，**不代表原链接当前可访问、内容已独立验证或获得再分发授权**。本批只整理已在用户分支中的最小元数据/检查提示，不搬运原始利用脚本、正文与图片。

## 构成

| category（= 源目录） | 条数 |
|---|---|
| `cloud-container` | 14 |
| `edge-device` | 16 |
| `exchange-ad` | 14 |
| `java-deser` | 9 |
| `java-web` | 22 |
| `linux-priv` | 11 |
| `web-cms` | 19 |
| `gray-market` | 59 |

- 96 条带 CVE 编号（其中 7 条多 CVE）；68 条无 CVE（通用打法与灰产手法）
- 50 条带一句话触发点（`entry`），全部来自 `gray-market`
- 全部 `verified: false`

## 字段映射

| `poc` 列 | 取值 |
|---|---|
| `code` | `fujiang-<cve 小写>`；多 CVE 追加 `-plusN`；无 CVE 用 `fujiang-<src_dir>-<sha1(title)[:8]>`。**与目录名完全一致**，供 `export_docs.py` 的 pass0 精确配对 |
| `title` | 源标题 |
| `kind` | `markdown` |
| `category` | 源目录名（8 个来源型分类，与既有 `absorbed` 同风格） |
| `cve` | 源 `cve[]` 逗号拼接，空则 NULL |
| `component` / `versions` / `severity` | 源 `product` / `versions` / `severity` |
| `language` | NULL（元数据条目无代码） |
| `source` | `kb` |
| `source_url` | 源链接（指回上游固定提交） |
| `description` | `产品 · 版本 · ports:<端口>`，截断 400 |
| `usage` | 仅 `gray-market` 写源 `entry`；其余 NULL |
| `content` | `pocs/<slug>/entry.md` 全文 |
| `path` | `<slug>/entry.md` |
| `verified` | 0 |
| `verified_note` | `导入自富江 0.20.13 索引，未实测` |
| `tags` | `keywords` + 分类 + `fujiang`，去重 |
| `created_by` | `fujiang-import`（导入器的幂等键） |

## 落盘内容

`pocs/<slug>/entry.md` 每份包含：标题、来源路径、分类/严重度/CVSS/CVE、产品、影响版本、端口、关键词、触发点、**定位要点**（组件指纹 + 信号来源 + 触发条件 + 判断口径）、原始链接。

判断口径里写明：本条为**未实测**索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

**没写进去的**：上游的利用步骤、payload、复现命令、图片。台账保留的是「往哪看」而不是「怎么打」，与仓库既有 508 条的「打过什么记什么」不混层。

## 执行与结果

```bash
python knowledge/tools/import_fujiang.py --dry
python knowledge/tools/import_fujiang.py
python knowledge/tools/export_docs.py
python knowledge/tools/kb_check.py --write
```

| 项 | 前 | 后 |
|---|---|---|
| `poc` 行数 | 508 | **672** |
| `poc_fts` 行数 | 508 | **672** |
| `pocs/` 目录 | 429 | **593** |
| MAP 已配对 | 428 | **592** |
| MAP 仅有台账 | 80 | **80**（不变） |
| MAP 仅有目录 | 1 | **1**（不变） |
| QA `poc_files` | 432 | **596** |
| QA `poc_bytes` | 1488822 | **1654945** |

- 新行 id 区间 `518..681`；FTS 孤儿 0
- 抽检 `poc_fts match`：`withdraw` 15、`actuator` 16、`nacos` 2
- `kb_check.py` 退出码 0，P1 无
- `tools/scan_sensitive.py` 退出码 0
- `go build ./...` 通过；`go test ./kb/` 通过
- 导入前备份：`knowledge/knowledge.db.bak-<时间戳>`（`.gitignore` 排除该备份文件，不入库）

## 版本保真与来源修订

- **保留精确版本号**：`product`、`versions` 及生成的 `entry.md` 不再把四段以上的厂商构建号截成三段。受影响范围仅是原资料声明，具体资产仍需按厂商安全公告逐项复核。
- 原索引中 59 条 `/blob/HEAD/` 链接已替换成所声明的固定提交 `cf086fc329011c6f8eeb3b67f0adbcc3f005debb`。固定化不代表现时 URL 可访问；完整声明原件的独立核验仍待来源恢复。
- 导入器先在临时数据库和临时目录生成全部内容、校验 FTS，再替换正式文档与数据库；普通写入错误会恢复旧文档，原数据库在最后一步前不触动。断电或强制终止不能实现跨文件原子性，回退请使用 `knowledge.db.bak-*` 并重新执行导入。
- `tools/scan_sensitive.py` 保留原有敏感信息识别规则，不通过丢失版本精度来消除误报；对误报应核验确切字段与来源。

## 与既有库的关系

- **1 条 CVE 重叠**：CVE-2020-1938（Ghostcat，Tomcat AJP）。既有条目是实战记录，本条是索引条目，两者用途不同，保留并存，不做去重。
- 新增 8 个 category 会让 `INDEX.md` 的分类分布表多 8 行来源型分类，与 `absorbed`（60 条，来源型）并列，属预期。

## 幂等与回滚

- 重跑导入器：在临时 SQLite 副本中先删除 `created_by='fujiang-import'` 的旧行及 FTS，再插入 → 条数不翻倍；数据库行 ID 可能变化，外部引用请使用稳定的 `code`。
- 撤销整批：`delete from poc_fts where poc_id in (select id from poc where created_by='fujiang-import'); delete from poc where created_by='fujiang-import';` 再删 `pocs/fujiang-*` 目录，重跑 `export_docs.py` 与 `kb_check.py --write`。
