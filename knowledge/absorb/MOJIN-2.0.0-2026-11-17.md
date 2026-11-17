# 摸金校尉 2.0.0（Windows x64）差异吸收记录

> 来源：用户提供的 `mojin-2.0.0-Windows-x64.zip` 静态解包；SHA-256 `d2d7b35a4946519e38ef56f0909e56199ddfb0302995b06244d13cef5611c07b`。
> 对照对象：`wxcvm/ARTEX` 的 `main` 技能目录和现有运行机制（2026-11-17）。
> 本记录只反映静态目录与配置对照，不声称已经运行、验证或信任来源 EXE。

## 目录差异与判断

| 项目 | 摸金校尉发行包 | ARTEX 当前主线 | 解释 |
| --- | ---: | ---: | --- |
| 技能目录（不含隐藏的 `.eino`，含 1 个演示 Skill） | 70（69 个非演示） | 54（首批前）/57（首批合并后） | 与最初 54 个目录名 **45 个完全一致**；名字不重合不等于能力缺失 |
| 仅摸金校尉出现的技能目录名 | 25（含 1 个演示） | — | 包含 `code-audit` ↔ `codeaudit`、`active-directory-attack` ↔ `chains-ad` 等重叠主题 |
| 最初仅 ARTEX 出现的技能目录名 | — | 9 | 包括 `api-security`、`finding-triage`、`patch-diff`、`rea` 等 |
| 工具 YAML 定义 | 90 | 不采用同一配置目录模式 | 发行包 78 个标记 enabled、12 个 disabled；标记启用**不等于工具二进制已安装** |

摸金校尉的 `tools/*.yaml` 包含 `api-schema-analyzer`（`spectral`）、`checkov`、`trivy`、`graphql-scanner` 等工具的参数草案。原包没有提供这些工具的完整二进制，部分配置包含 `sh`、`/bin/bash` 或未经验证的自由文本参数；不能直接在 Windows 运行，也不能直接混入 ARTEX 的原生工具 handler / MCP 绑定。应逐项核验具体工具安装方式、参数语义、返回格式与现有功能覆盖，再决定是否扩充现有工具管理器。

## 本专项实际吸收

1. `skills/ai-llm-app-security/SKILL.md`：为 AI 应用、RAG、MCP 与 Agent 工具调用增加**证据驱动的检查方法**，避免只靠提示词回显判断成功。
2. `skills/component-advisory-triage/SKILL.md`：按组件版本、官方公告、受影响范围和实际证据进行漏洞情报核验；与现有 `finding-triage`、`patch-diff` 协同。
3. `skills/blockchain-contract-review/SKILL.md`：覆盖智能合约/DeFi 授权、资产会计、预言机、签名、升级和跨链信任边界的源码审查方法；复用现有 `codeaudit` 与 `finding-triage`。
4. `skills/mobile-binary-review/SKILL.md`：在现有 APK 静态提取、REA 二进制分析之上，补充 UniApp/Flutter、Manifest/深链、WebView、native 与运行时证据检查。
5. `skills/configuration-security-review/SKILL.md`：为 IaC/容器配置、API Schema 与锁定依赖增加证据型检查流程；可利用已安装的 Spectral/Trivy/Checkov，但不承诺这些工具已安装或可自动安装。

五份 Skill 均为 ARTEX 独立编写，只借鉴**主题与检查维度**；未直接复制原包的指令、脚本或风险配置。它们是可检索的 Skill 文档，不新增工具执行器、不装软件、不修改主程序、不改变批准策略或任务范围。

## 其余主题去重与不迁移原因

- **已有等价能力**：`active-directory-attack` → `chains-ad`；`attack-surface-recon` → `active-scan` / `recon-pipeline`；`code-audit` / `source-code-hunting` → `codeaudit` / `patch-diff`；`cloud-attack-methods` → `cloud-security`；`post-exploitation` → 现有主机与决策链技能；`web-attack-methods` → `api-security` / `chains-web`。
- **已有编排/证据能力**：`pentest-agent-os`、`pentest-blackboard`、`pentest-output-standards`、`pentest-verification`、`capability-primitive-search`，对照已有 planner/worker、资产和事实图、coverage ledger、`finding-triage` 等。**概念覆盖不代表字段/API 完全相同**，不直接引用摸金校尉专有 `upsert_project_fact` 等函数。
- **本批已新增对应主题**：AI/LLM、组件公告、合约、移动逆向；`trivy`、`checkov`、`spectral` 等对应的静态配置检查由新 Skill 提供方法，实际工具调用仍依赖现有运行环境。
- **不适合直接迁移**：`find-skills`（工具/技能发现已有管理入口）、`cyberstrike-eino-demo`（演示素材）、`wireless-hardware-attack`（当前没有已核实的硬件采集/执行环境）、`specialized-attack-playbooks`（场景专属目标/样本未确认）、`proxy-tool-bootstrap`（不能强制覆盖用户代理）、`initial-access-phishing`、`redteam-opsec`、`unlimited-attack-scope` 等激进或强依赖场景的约束性指令。

## 工具 YAML 处理结论

已枚举和读取发行包全部 **90** 份 `tools/*.yaml`；其中 **78** 份标记启用。它们主要定义名称、命令和参数，并**不**附带工具本体。ARTEX 现有 `server/external-tools.json` 提供 **7** 个官方来源的 Windows 工具探测与管理项，安装代码对官方 release 资产名称、摘要、目录结构有专门约定；不能将这 90 个 YAML 当成可安装列表。

- **已有路径或能力**：`nuclei`、`katana`、`ffuf`、`subfinder` 及其他既有技能调用，沿用 ARTEX 的工具来源、宿主和 MCP 机制。
- **有价值但暂不直接接入执行器**：`trivy`、`checkov`、`api-schema-analyzer` (Spectral)、`graphql-scanner`，按对应上游 CLI 实测版本、安装位置与输出格式后才能创建执行绑定。
- **不复制的危险/不可移植配置**：`exec` 使用 `sh`，还有 Linux 专属路径及多种自由格式 `additional_args` 参数。直接引用会在 Windows 上失效，或绕过 ARTEX 现有的命令处理和权限语义。

本次吸收的是**经过去重的指导流程**，并非宣称 90 个第三方工具已经在 ARTEX 可运行。未来是否接入某个新 CLI，应由具体用途和兼容验证驱动，而非为了凑数量。

## 未直接吸收的内容

- `unlimited-attack-scope`、强制绕过访问限制的代理/搜索流程、要求每轮穷尽扫描的指令：与既有授权范围及可控执行要求不相容，且没有证据支持效果收益。
- 大篇幅系统提示词、固定社群宣传或“模型拒答即自动重试三次”：不移植到 ARTEX。先在 Purge 现有 Provider / 流式 / 拒答评测框架中比较实际行为，确认有效且无回归再讨论实现。
- Web 管理界面、数据库、C2/WebShell、浏览器插件、MCP 主框架：ARTEX 已有相应子系统，不重复建造。
- 发行包中无法核实的来源和许可证条款：未作为可直接复制代码的授权依据。

## 本批验收边界

- 文档 frontmatter 可解析，Skill 名称互不重叠，且新文件在仓库中可读。
- 不涉及 Go/Next.js 运行代码；静态检查可做，Windows/容器运行时的 Skill 发现与实际应用效果仍需对应构建实测。
- **Purge 本批不修改运行时代码**：其现有拒答分析、Provider / 流式与 ARTEX 本地联动正在独立演进，不能以别的产品的重试配置取代基线或声称效果提升。
