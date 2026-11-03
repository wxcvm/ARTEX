# 首批自主化工程改进：验证记录

日期：2026-11-02（Asia/Shanghai）。本记录是本地验证，不是发布验收。

- 基线：`wxcvm/ARTEX@7f8c7cbcae7324e823d8952371ae2013521cc711`。
- 本地分支：`codex/autonomy-foundation`；未提交、未推送、未发布。
- 环境：Windows amd64；便携 Go 1.26.3；Node 24.20.0 / npm 11.19.0；竞态检测使用本机 GCC。
- 没有配置 PostgreSQL，没有调用真实模型或对外部目标执行评估。

## 1. 实际改动

1. `Complete` 在取消、超时、HTTP 400/422 等不应切换模型的错误上，保留原有健康状态；错误仍原样返回，不触发备用调用。只有成功才清除状态。
2. 模型健康持久化改为单个异步写入协程；按 profile 合并待写状态。慢数据库不阻塞健康状态读取和手动恢复，不再为每次写入单独创建协程；保存、删除按序完成。数据库写入有 5 秒上限，错误记录日志。
3. 报告导出和工作区下载共用 ASCII fallback + UTF-8 `filename*` 响应头；前端优先解析扩展字段，对坏编码回退，并处理引号、分号和路径分隔符。
4. CI 增加模型池竞态回归和前端单元测试，未改变发布工作流、依赖版本、模型配置或已有执行限制。

持久化仍是 **best-effort**：回调失败不会被误称为写入成功；`FlushPersistence` 只保证等待待执行回调结束，不代表跨进程事务、数据库提交成功或突然断电不丢状态。本轮没有改动进程关闭协议。

## 2. 已执行的验证

| 检查 | 实测结果 |
| --- | --- |
| `go build ./...` | 通过，退出码 0 |
| `go vet ./...` | 通过，退出码 0 |
| `go test ./llmpool -count=10` | 通过 |
| `go test -race -count=20 ./llmpool` | 通过；覆盖慢写入后的恢复、状态合并、不同 profile、后续新一轮写入及取消等待 |
| `go test ./server -run TestDownload -count=1` | 通过；7 个文件名边界子用例 |
| `go test server/download.go server/download_test.go` | 通过；独立标准库验证 |
| `npm test --prefix web` | 25 项通过，0 失败；18 项新文件名用例及 7 项既有用例 |
| 原版 `pool.go` 配合新增非成功回归测试 | 6 类错误均失败，复现旧版本错误清空状态；替换回本轮实现则通过 |
| `git diff --check` | 通过；存在 Git 的 LF/CRLF 提示，不是差异错误 |

新增测试没有连接模型服务、数据库或目标网络，结果不代表真实模型任务完成率。

## 3. 全量测试：没有全绿

执行了 `go test -p 1 -json ./...`，最终退出码 **1**。按 Go JSON 中测试完成事件计数（含子测试）：

- 通过：**614**。
- 跳过：**242**，不能计作通过。
- 失败：**36**，全部位于 `server` 包。

对每条失败输出的分类：

| 数量 | 实际错误 |
| --- | --- |
| 34 | 未设置 `ARTEX_PG_DSN`，也没有数据库配置；相关辅助函数直接失败，而另外一部分测试选择跳过 |
| 1 | `TestExecPython` 的子命令退出码为 9009，预期输出缺失 |
| 1 | `TestSkillRelPath` 期望 `/`，Windows 返回 `\` |

这不是“全量通过”。尚需在带 PostgreSQL 的 Linux CI 和真实部署环境复测；本轮不为让数字变绿而删除测试、修改既有跳过规则或隐去失败。完整 JSON、失败分类和审计输出保存在本地 `.git/artex-validation/`，不随源码发布。

## 4. 前端完整构建未验收

源码目录的 `npm ci` 在依赖重建阶段长时间没有推进，已仅终止本轮启动的安装进程。随后将源码和锁文件复制至 OneDrive 之外的隔离目录并核对关键文件 SHA-256，再次安装仍在 npm 缓存写入时出现 `EEXIST / UNKNOWN copyfile`，并伴随清理时 `EPERM`。

- 未删除或清理用户的全局 npm 缓存，未用强制覆盖掩盖该错误。
- 没有完成完整依赖安装，因此 **未完成全量 TypeScript 检查和 Next.js 生产构建**。
- 文件名纯函数测试使用 Node 内建测试器，不依赖该安装，25 项结果有效。
- 未做浏览器真实下载弹窗、文件落盘及内容核验。

## 5. 已复现但留待下一批的 P0 问题

暂停恢复时，数据库与 `Task.DeadlineAt` 已更新，但 `Engine.taskDeadline` 仍优先返回旧缓存。临时无数据库审计用例：先读 deadline=100，再更新生命周期 deadline=200，再读引擎，实测失败：

```text
engine deadline=100, want resumed deadline 200
```

临时失败用例已移出正常测试目录，保留在本地审计目录。截止时间协调器还没有暂停判断；完整修复应处理暂停、恢复和进入 settling 的原子竞争，不是只替换一个缓存读取。具体实施与验收要求见 `AUTONOMY_ROADMAP.md` 的 P0.2。

## 6. 工作区边界

首次检出后，以下 3 个原有跟踪文件先呈现不可读，随后在工作区缺失。未确定宿主原因，本轮未尝试关闭保护、恢复这些文件或改写它们；它们不是本轮优化的删除项：

- `scripts/poc/php2shell/poc/verify_handler_mismatch.py`
- `scripts/poc/php2shell/poc/verify_session_injection.py`
- `scripts/poc/php2shell/poc/verify_uaf_asan.php`

提交时必须按本轮文件清单选择，不应直接把整个工作区 `git add -A`。没有整批合并参考仓库，没有恢复其自更新、替换 Norma/Ollama 行为或更改正式发布渠道。
