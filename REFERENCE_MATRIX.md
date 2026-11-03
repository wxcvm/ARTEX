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
