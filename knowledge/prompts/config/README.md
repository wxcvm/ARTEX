# 部署默认配置（来自 0.4.2 发行包）

这里放的是**发行包 / 部署者**层面的默认值，不是代码内建默认。三者关系：

| 层 | 存在哪 | 谁改 | 重启是否覆盖 |
|---|---|---|---|
| 代码内建默认 | Go 源码（`db`/`server` 里的常量与 seed） | 开发者 | — |
| 部署默认（本目录） | `settings.defaults.json` / `intercept-rules.json` / `armor/defaults.sql` | 部署者 | 只在标记 `defaults.deployment_state` 不存在时应用一次 |
| 当前运行值 | 数据库 `settings` / `intercept_rules` 表 | 面板 / API | 持久保存，不被部署默认覆盖 |

## 文件

- `settings.defaults.json` — 0.4.2 发行包的 31 项 `settings` 初始值（面板、LLM、通知、调度、覆盖率等）
- `intercept-rules.json` — 8 条内置拦截规则（命令/资产两类），即 `guard` 在 `PreToolUse` 时求值的那批规则
- `../armor/defaults.sql` — 部署者把 armor 开关 + 4 套预设 + 拦截规则写成一条 SQL，首次启动后由启动脚本应用一次
- `../armor/*.json` — 4 套 armor 预设原文

## 注意

- 想让 ARTEX 用上这里的部署默认，需要自己把 SQL 或 JSON 写进库（发行包是用 `defaults.sql` 做的），**本仓库不会自动应用**。
- `intercept-rules.json` 里的规则是**用户可删**的普通记录，不是硬编码；删掉就失去对应的拦截。
- 改这里不影响代码默认值；升级代码也不会覆盖本目录。
