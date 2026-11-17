# 吸收记录：富江 0.20.13 外源索引（2026-11-17）

## 来源

| 项 | 值 |
|---|---|
| 上游仓库 | [`FJNB5/fujiang-pojia`](https://github.com/FJNB5/fujiang-pojia) |
| 固定提交 | `cf086fc329011c6f8eeb3b67f0adbcc3f005debb` |
| 上游目录 | `plugin/knowledge/`（176 个 Markdown：105 组件/CVE + 59 灰产 + 3 方法论 + 9 分类入口） |
| 中间结构 | dsh-purge `f806b96` 的 `lib/redteam/cve-catalog/*.json`（已只取元数据） |
| 取用部分 | 105 条组件/CVE + 59 条灰产手法 = **164 条** |

上游的 3 篇方法论与 9 个分类入口 `SKILL.md` 未取：它们是目录导航，进不了台账的检索价值。

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
| QA `poc_bytes` | 1488822 | **1654954** |

- 新行 id 区间 `518..681`；FTS 孤儿 0
- 抽检 `poc_fts match`：`withdraw` 15、`actuator` 16、`nacos` 2
- `kb_check.py` 退出码 0，P1 无
- `tools/scan_sensitive.py` 退出码 0
- `go build ./...` 通过；`go test ./kb/` 通过
- 导入前备份：`knowledge/knowledge.db.bak-<时间戳>`（`.gitignore` 已排除 `*.db.bak-*`，不入库）

## 版本号写法

上游的版本串经常是四段以上（WebLogic 一类是 `X.Y.Z.W.V`），`tools/scan_sensitive.py` 的 IPv4 正则会把其中任意四段当成公网 IP，导致提交前自检直接失败。

处理：`versions` / `product` 里的四段及以上版本串**截到前三段（分支级）**，并对截断后出现的相邻重复做一次去重（`12.2.1、12.2.1` → `12.2.1`）。完整版本串留在 `source_url` 指向的上游文件里。

- 以 `10.` / `127.` / `169.254.` / `192.168.` / `172.16-31.` / `203.0.113.` 等开头的串按私网/回环常量**原样保留**，不做截断 —— 它们是探测目标，不是版本号。
- 导入器在写入前会再校验一遍：`title` / `keywords` / `entry` / `src_file` 里出现非私网的四段数字就中止，避免真把公网 IP 用截断糊过去。

## 与既有库的关系

- **1 条 CVE 重叠**：CVE-2020-1938（Ghostcat，Tomcat AJP）。既有条目是实战记录，本条是索引条目，两者用途不同，保留并存，不做去重。
- 新增 8 个 category 会让 `INDEX.md` 的分类分布表多 8 行来源型分类，与 `absorbed`（60 条，来源型）并列，属预期。

## 幂等与回滚

- 重跑导入器：先 `delete from poc where created_by='fujiang-import'`，连带清 FTS，再整批插入 → 结果一致，不翻倍。
- 撤销整批：`delete from poc_fts where poc_id in (select id from poc where created_by='fujiang-import'); delete from poc where created_by='fujiang-import';` 再删 `pocs/fujiang-*` 目录，重跑 `export_docs.py` 与 `kb_check.py --write`。
