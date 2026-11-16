# AntiFox source inventory / integration status

固定上游：AntiFox-AI/antifox @ `070f798bf590667bc2a0a13c687dd87d8a4413bf`（AGPL-3.0）。
ARTEX 对照：wxcvm/ARTEX @ `85373c2b185880826c15ebb3f2d86f90a37dbf2f`（含已合并 #106 及并行变更）。

[逐文件库存 TSV](antifox-source-inventory.tsv) 覆盖 Git tree 的全部 474 个 blob，两个 tree 均未截断。按相同路径与 Git blob SHA 比较：108 相同、188 不同、178 仅上游存在。该库存不是逐行语义审计；文件路径不同或哈希不同并不自动说明功能缺失，也不等于需要移植。

| 专题 | 参考来源 | ARTEX 交付与处理 |
|---|---|---|
| 报告与证据 | report/*、db/finding_evidence.go、server/finding_evidence.go、report_mgmt.go | #94/#96/#97：Go OOXML、截图/复现、选中发现、专业报告，已合并。#105：结构化文字模板。#109：实际文件快照、任务归档、全局报告中心；后两项原分支 CI 已通过，整合进 #109 后等待新 head CI。未移植 Python 渲染、任意外部 DOCX 模板或 AI 生成漏洞事实。 |
| 知识 | knowledge/*、db/knowledge.go、agent/tools_knowledge.go、server/knowledge.go | #98/#99/#100/#101 已合并：双层文档、统一检索、来源、检索历史、本地 embedding、增量索引及 PDF。保留 bundled SQLite 和 PostgreSQL，无强制 pgvector/Python/远端 embedding。 |
| 更新 | selfupdate/*、server/update.go | #102 已合并：稳定版本、官方资产和设置入口。现有 Windows 安装器与 Docker/源码流程继续负责升级。不复制 supervisor 自动替换。 |
| 运行时 | engine/harness/streaming.go、server/engine.go、agent/chat.go、agent/taskclock.go | 保留 Norma v0.4.6、ARTEX 具名取消、任务生命周期与预算。#106 已合并：活动流背压/顺序/重放。#110：聊天先落库、历史工具状态，原分支 CI 已通过，已整合至 #109，等待新 head CI。不替换另一套引擎。 |
| 身份与审计 | db/rbac.go、server/rbac*.go、db/audit.go、server/audit.go、server/loginlimiter.go | ARTEX 单管理员和专用桥认证、任务/模型/工具记录与上游多用户模型不同；完整 RBAC/平台审计/登录限流需要明确产品需求，不在本次更改认证契约。 |
| 其余上游差异 | engine、skills、web、db 等 | TSV 保留定位线索。相同文件不重复复制；异名、变更文件需按实际功能和测试核查。尚未逐行审计全部 474 个文件。未复制 AntiFox 专有品牌、图片、部署名或整套运行时。 |

验收原则：合并/CI 是实际交付依据；待验收 PR 不算已合并。源文件库存固定于上述 SHA，新提交后重新计算才是新的基线。没有根据文件数推算能力提升或吸收百分比。

## 集中整合（2026-11-16）

#105、#109、#110 原提交 CI 全绿。基于最新 main `a8c42d8` 的三分支整合提交 `b48d089` 保留原祖先与并行 HAR/流量证据变更；唯一共同变更文件 `server/server.go` 已按精确上下文合并。新 head 尚未完成 CI，不能以原分支 CI 代替。平台多用户/RBAC、审计、三语覆盖和任意外部 Word 模板仍为未完成差异；不宣称全部源码语义核验或已发布。
