# 摸金校尉 2.0.0（Windows x64）差异吸收记录

> 来源：用户提供的 `mojin-2.0.0-Windows-x64.zip` 静态解包；SHA-256 `d2d7b35a4946519e38ef56f0909e56199ddfb0302995b06244d13cef5611c07b`。
> 对照对象：`wxcvm/ARTEX` 的 `main` 技能目录和现有运行机制（2026-11-17）。
> 本记录只反映静态目录与配置对照，不声称已经运行、验证或信任来源 EXE。

## 目录差异与判断

| 项目 | 摸金校尉发行包 | ARTEX 当前主线 | 解释 |
| --- | ---: | ---: | --- |
| 技能目录（不含隐藏的 `.eino`） | 69 | 54 | 其中 **45 个目录名完全一致**；名称不重合不等于能力缺失 |
| 仅摸金校尉出现的技能目录名 | 24 | — | 包含 `code-audit` ↔ `codeaudit`、`active-directory-attack` ↔ `chains-ad` 等重叠主题 |
| 仅 ARTEX 出现的技能目录名 | — | 9 | 包括 `api-security`、`finding-triage`、`patch-diff`、`rea` 等 |
| 工具 YAML 定义 | 90 | 不采用同一配置目录模式 | 发行包 78 个标记 enabled、12 个 disabled；标记启用**不等于工具二进制已安装** |

摸金校尉的 `tools/*.yaml` 包含 `api-schema-analyzer`（`spectral`）、`checkov`、`trivy`、`graphql-scanner` 等工具的参数草案。原包没有提供这些工具的完整二进制，部分配置包含 `sh`、`/bin/bash` 或未经验证的自由文本参数；不能直接在 Windows 运行，也不能直接混入 ARTEX 的原生工具 handler / MCP 绑定。应逐项核验具体工具安装方式、参数语义、返回格式与现有功能覆盖，再决定是否扩充现有工具管理器。

## 本批实际吸收

1. `skills/ai-llm-app-security/SKILL.md`：为 AI 应用、RAG、MCP 与 Agent 工具调用增加**证据驱动的检查方法**，避免只靠提示词回显判断成功。
2. `skills/component-advisory-triage/SKILL.md`：按组件版本、官方公告、受影响范围和实际证据进行漏洞情报核验；与现有 `finding-triage`、`patch-diff` 协同。

两份文件都为 ARTEX 重写，只借鉴**主题与检查维度**；未直接复制原包的指令、脚本或风险配置。它们是可检索的 Skill 文档，不新增工具执行器、不装软件、不修改主程序、不改变批准策略或任务范围。

## 未直接吸收的内容

- `unlimited-attack-scope`、强制绕过访问限制的代理/搜索流程、要求每轮穷尽扫描的指令：与既有授权范围及可控执行要求不相容，且没有证据支持效果收益。
- 大篇幅系统提示词、固定社群宣传或“模型拒答即自动重试三次”：不移植到 ARTEX。先在 Purge 现有 Provider / 流式 / 拒答评测框架中比较实际行为，确认有效且无回归再讨论实现。
- Web 管理界面、数据库、C2/WebShell、浏览器插件、MCP 主框架：ARTEX 已有相应子系统，不重复建造。
- 发行包中无法核实的来源和许可证条款：未作为可直接复制代码的授权依据。

## 本批验收边界

- 文档 frontmatter 可解析，Skill 名称互不重叠，且新文件在仓库中可读。
- 不涉及 Go/Next.js 运行代码；静态检查可做，Windows/容器运行时的 Skill 发现与实际应用效果仍需对应构建实测。
- **Purge 本批不修改运行时代码**：其现有拒答分析、Provider / 流式与 ARTEX 本地联动正在独立演进，不能以别的产品的重试配置取代基线或声称效果提升。
