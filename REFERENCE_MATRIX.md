# ARTEX 外部方案参考矩阵

核查日期：2026-11-03。保留 Go + PostgreSQL + Next.js；本阶段没有引入下列项目的运行时、生产依赖或新服务。

## 状态与固定版本

“已核查”仅表示阅读了项目/文档并确认参考方向，不代表完整源码审计或复现其成绩。“已采用”指下面明确指出的设计契约，不表示移植了第三方模块。“仅供参考”不得成为当前功能声明。

| 来源 | 固定源码快照 | 状态 | 解决的问题、现有重叠、采用方式 | 验证与维护成本 |
| --- | --- | --- | --- | --- |
| [Shannon](https://github.com/KeygraphHQ/shannon) | `b9744f6cecc5429333873b86a6762f1d70c54fc2` | 已核查；结果区分契约已采用 | 参考报告与结果组织；复用 ARTEX 已有 outcome API、报告及证据存储，不另建评估运行时。明确区分完成无发现、未完成、读取失败。结构化导出后续独立提交。 | `task-outcome.test.mjs` 与 9101–9111 mock 页面；成本低，DTO 与状态文案变动时同步测试。 |
| [Strix](https://github.com/usestrix/strix) | `6b7ae2ebe4d47037a927aade23265ff4e4f27e84` | 已核查；仅供参考 | API 契约输入、执行隔离、协作记录展示；重叠于现有 MCP、资产覆盖图、Agent 会话。后续按接口导入/覆盖可视化模块评估，不启动第二个平台。 | 后续 OpenAPI 固定样例、输入错误和覆盖映射测试；成本中高，涉及契约版本与工具适配，本阶段不实现。 |
| [PentestGPT](https://github.com/GreyDGL/PentestGPT) | `e8b1bb77d1ac00329675cec3b060aba971ec1ac8` | 已核查；仅供参考 | 阶段状态、会话和模型后端组织；复用 ARTEX 模型池与生命周期，不用历史基准成绩推导当前表现。 | 后续会话中断与模型失败回归；成本中，避免重复适配层。 |
| [LangGraph](https://github.com/langchain-ai/langgraph) / [持久化文档](https://docs.langchain.com/oss/python/langgraph/persistence) | `3c916d750575b22c464790e9774494d810ccd322` | 已核查；持久化状态为准的原则已采用 | checkpoint 与执行状态分离；现有 PostgreSQL 生命周期镜像作为截止时间唯一来源。此修复不是完整 checkpoint、outbox 或崩溃恢复机制。 | `TestTaskTransitionDeadlineUsesCommittedSnapshot`、PG 恢复测试；当前成本低，未来持久化受理需要独立设计。 |
| [Temporal Go SDK](https://github.com/temporalio/sdk-go) | `17e2e8eae4ccc30db9063e3f056d12b5eef68028` | 已核查；显式取消/受理语义已采用 | 借鉴编排与执行分离、显式重试/取消；复用 Engine、taskStateMu、任务级转换锁，不引入 Temporal 服务。 | 暂停先赢、收尾先赢、提交失败、删除交错、协调器继续等待；成本中，锁顺序与终态测试需持续维护。 |
| [OpenTelemetry Go](https://github.com/open-telemetry/opentelemetry-go) / [文档](https://opentelemetry.io/docs/languages/go/) | `4aafbd5466847ad42e0d837f36d34cd8186190aa` | 已核查；仅供参考 | 调用关联与耗时观测；现有调用账本与证据是事实来源，遥测不是证据存储。后续可选 exporter，失败不得阻塞任务。 | 后续 exporter 不可用与关联 ID 测试；成本中，版本升级与标签基数需要约束。 |
| [Promptfoo](https://github.com/promptfoo/promptfoo) / [测试用例文档](https://www.promptfoo.dev/docs/configuration/test-cases/) | `06778a39b27cb63345c23872d19503da012636cf` | 已核查；固定输入、断言、重复运行方法已采用 | 扩展既有 benchmark 目录，新建独立确定性工程 runner；不调用真实模型、不混入 L3 分数。 | 10 个 loopback HTTP 用例 + 生命周期 + UI 断言，20 次重复，逐证据 SHA256；成本低，无新增 Python 包。 |
| [CAI](https://github.com/aliasrobotics/cai) | `6dc79257777f5f1c9500b4d2319935d34a47412e` | 已核查；仅历史参考 | 仓库已归档、停止维护；不成为新核心依赖。 | 无移植，无运行依赖；若后续考虑其中模块，重新审查维护状态。 |
| [AntiFox](https://github.com/AntiFox-AI/antifox) | `070f798bf590667bc2a0a13c687dd87d8a4413bf` | 2026-11-16 源码对照；部分设计已实现于专题 PR | `db/knowledge.go`、`knowledge/{service,extract,chunk,embed}.go`、`server/knowledge.go` 与全局/任务 UI；采用文档管理和作用域设计，保留 ARTEX 内置 SQLite KB。自有文档使用现有 PostgreSQL、字面关键词检索和 Go 原生 DOCX 文本提取，纳入任务归档。未移植 pgvector、外部 Embedding 或 Python 解析器。后续 #99 接入现有智能体工具和检索历史；本地语义索引专题采用 OpenAI-compatible 本机 Embedding、PostgreSQL 派生向量与 Go 精确相似度，默认关闭、手工建立/重建索引，失败保留关键词检索；后续新增批量文件导入、增量索引（精确文本复用）、可选本机 Poppler PDF 文本抽取、智能体词项排序。 | AGPL-3.0 来源；独立实现，未复制源码。测试包含作用域隔离、CRUD、字面查询、上传与归档往返；整合 CI 待验收。 |

以上 SHA 是本次 `git ls-remote <repository> HEAD` 固定的源码快照；网页文档为核查日阅读版本，不声称网页内容与源码提交严格同版。项目声明与 ARTEX 实测分开记录。

## 后续参考准入模板

每次扩展单独提交以下记录，不允许整仓覆盖：

1. 来源 URL、不可变 commit/tag、具体模块路径及必要的署名/许可证文件。
2. 当前缺陷或覆盖缺口、与模型池/MCP/资产图/结果卡片/benchmark 的重叠。
3. 采用方式：设计参考、适配器、独立模块移植，或暂不采用。
4. 固定输入、输出契约、失败语义、重复/取消/大输出测试。
5. 依赖与基础设施增量、升级责任、回退方法、可验证的维护成本。
6. 在相同 source/case version 下比较工程结果；真实模型评测另表记录。

下一批候选：OpenAPI 输入、结构化报告导出、可选观测。它们均不在本批生产执行流程中。

## AntiFox 平台与运行时差异（2026-11-16）

对照固定 AntiFox `070f798bf590667bc2a0a13c687dd87d8a4413bf` 的 `server/auth.go`、`db/audit.go`、`server/audit.go`、`server/update.go`、`selfupdate/selfupdate.go` 与 ARTEX 现有 `server/auth.go`、`server/assembly.go`、`server/engine.go`、`agent/cancelcause.go`、Windows 安装器及 release workflow：

- **更新**：AntiFox 的二进制 staging/smoke/回滚依赖外部 supervisor。ARTEX 已有 Docker/源码 update.sh、Windows Setup 升级烟测和独立桌面进程；本批仅采用明确版本比较、平台资产匹配与发布页入口，不移植退出码 75 或在线替换 EXE。运行中应用不自动退出，不自动下载/执行文件。
- **身份**：AntiFox 多用户 RBAC、token_version、滑动续期与角色权限映射与 ARTEX 单管理员及专用 Purge/隐私桥凭证不同；不在未授权情况下改变既有认证行为，完整 RBAC 另行评估。
- **审计**：AntiFox 持久化登录/平台操作日志与 ARTEX 现有任务 guard 审计、模型出站审计、工具使用账本不同。暂不重复记录请求正文/密钥或扩大日志范围。
- **运行时**：AntiFox 的 `engine/harness/streaming.go` 和 `server/engine.go` 属于另一套 harness/runtime；ARTEX 当前使用 Norma、任务级取消原因、模型预算及 lifecycle 回归测试。不整仓替换，后续按流式事件和取消/恢复契约逐项对照。


## BaaRai/ARTEX 执行器权限边界对照（2026-11-27）

来源：[BaaRai/ARTEX PR #510](https://github.com/BaaRai/ARTEX/pull/510)，固定上游实现提交 `b6748f8ec08ec61c80b6c251491d8da76c83c315`。

- 上游为 Linux 执行 shim 增加 fail-closed 身份切换与 `PR_SET_NO_NEW_PRIVS`，并用 `runtime.LockOSThread` 保证设置线程标志的线程继续执行 `exec`。
- 本 fork 的 opt-in `docker-compose.runner.yml` 和 `docker-compose.mode-b.yml` 已在 broker、executor 与 egress 边界采用容器级 `no-new-privileges:true`、`cap_drop: ALL` 和只读根文件系统；executor 仅挂载任务/会话目录，不获得 PostgreSQL、JWT 或模型密钥。
- 这不是逐任务 UID/GID 降权的等价实现：容器内执行器目前仍使用镜像默认用户。PR #510 的每线程 shim 机制不直接移植，因为当前 fork 使用独立 toolrunner 容器与 Norma 工具管理器，执行边界不同。
- 本次通过 Compose 静态回归测试锁定上述现有控制，防止后续配置误删。默认非隔离 Compose/native 部署不因此自动获得相同隔离级别；部署文档应明确 opt-in runner 的边界。

状态：已核查；容器级控制已采用；逐任务身份降权保留为独立设计缺口。


## BaaRai/ARTEX 模型连接测试安全对照（2026-11-27）

来源：[BaaRai/ARTEX PR #526](https://github.com/BaaRai/ARTEX/commit/1f558d8681be3f55cfcd2e296e21af2fcfbe570a)。

连接测试必须是纯文本探测：显式传空工具列表、禁用 SDK 后台任务工具，并拒绝模型返回工具调用或 tool_use/tool_calls 结束原因（即使同时带有文本）。本 fork 采用该约束，并用流式/非流式假服务回归测试覆盖。该路径不得执行文件、Shell、后台任务或任何外部工具。


## Planner empty-context guard (2026-11-27)

- Planner rounds now resolve the task goal from the explicit argument, persisted root goal, or first non-empty goal node, in that order.
- A nil exploration store or a task with no usable goal context is rejected before compaction or model invocation and emits an error activity; database read errors propagate.
- The graph overview omits blank goal nodes and empty root fields so the planner cannot mistake placeholder data for a real objective.
- Regression tests cover nil stores, empty roots, blank goal nodes, and trimmed non-empty goal fallback.
