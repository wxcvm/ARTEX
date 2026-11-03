# Web/API 优先第一阶段：实施与验收记录

日期：2026-11-03。基线 `7f8c7cbcae7324e823d8952371ae2013521cc711`，分支 `codex/autonomy-foundation`，改动尚未提交。上一批模型健康、持久化与 UTF-8 文件名改动保留。

**结论：A/B/C 实现及本机确定性回归已交付；不能宣告整个第一阶段验收完成。带 PostgreSQL 的 Linux CI 尚未运行，本机全量 Go 测试仍失败。没有推送、合并或发布。**

## A：参考选型

- `REFERENCE_MATRIX.md` 固定八个外部项目源码 SHA，区分已核查、设计契约已采用、仅供参考；记录现有重叠、验证和维护成本。
- 仅扩展路线和方法，不新增运行时/必需服务；OpenAPI 输入、结构化报告、可选遥测均留在后续独立变更。
- `RESTRICTION_AUDIT.md` 单列追加的内部限制检查；未重置用户配置、删规则、改密钥或取消生命周期保护。

## B：生命周期修复

- 移除 Engine 的第二份 deadline 缓存，`taskDeadline` 读取任务生命周期镜像。`100 → 200` 直接回归覆盖；已有 PG 暂停恢复测试也增加 Engine 断言。
- 每任务 `transitionMu` 串行化用户暂停提交/补偿与超时收尾受理。顺序为任务转换锁 → Manager.taskStateMu → Manager.mu → lifecycleMu；取任务指针时先释放 Manager.mu 再等待任务转换锁。
- 暂停先受理：先安装执行屏障，提交并同步镜像后才取消执行；收尾继续等待，恢复读取新截止时间。
- 收尾先受理：返回明确冲突。沿用单任务接口现有 HTTP 409 和错误格式；批量调用同一实现，仍为逐项结果。
- 提交失败：恢复原屏障，不取消原健康执行上下文；不移除删除流程拥有的屏障。
- 保留终态保护、删除屏障、任务预算和排队计时。虚拟时钟覆盖等待、恢复、取消、终态与 queued 行为；收尾只受理一次。
- PG 集成新增真实 CHECK constraint 失败补偿测试；原暂停恢复测试用数据库的固定十分钟 paused_at 替代 wall-clock sleep。

这不是持久化任务受理/outbox、跨进程锁或任意崩溃恢复的实现；不能据此推导这些保证。

## C：确定性 Web/API 基线与结果展示

- `server/webapi_engineering_test.go` 的 10 个本地用例复用实际 MCP HTTP 客户端、分页 helper、批量 ID 归一化和模型池；无外部目标、无真实模型。
- 场景包括正常、分页、401、429、实际 HTTP 超时、坏 JSON、大 UTF-8、取消、重复批量 ID、不完整结果。
- `scripts/benchmark/run_engineering.py` 独立于 L3 benchmark；报告记录源码/worktree 哈希、用例版本、pass/fail/skip、耗时、调用次数、证据路径与 SHA256；完整性逐文件重新核验。
- 复用现有 `/tasks/{id}/outcome`，没有新增控制端点/后端 DTO。结果卡片仅在任务及 outcome 都 done、且无 paused/queued 标记时允许“完成无发现”；空白、无效或失败响应不能补零。
- 报告请求失败显示独立错误及重试，不冒充空报告；保留上一成功快照并注明 stale。

## 本机实测

| 验证 | 结果与边界 |
| --- | --- |
| `go build ./...` / `go vet ./...` | 均 exit 0；Windows amd64，Go 1.26.3。 |
| 确定性 runner `--repeat 20 --race` | 通过：200 个 HTTP 场景记录、160 个顶层生命周期测试通过；Go 含父/子测试共 500 pass，0 fail，0 skip。前端 outcome 18 × 20 通过。 |
| 调用与证据 | 620 HTTP、220 工具、40 模拟模型调用；真实模型 0；620 份请求/响应证据经 SHA256 复核。不是漏洞发现成绩。 |
| `go test -race -count=20 ./llmpool` | 通过，保留首批并发回归。 |
| `npm test --prefix web` | 43 pass，0 fail，0 skip。 |
| `tsc --noEmit` | exit 0；在独立依赖目录执行。 |
| 生产静态构建 | `NEXT_EXPORT=1 NEXT_TELEMETRY_DISABLED=1 next build` exit 0；31 个静态页面。源码 src 下 201 个文件与工作区逐一哈希一致。 |
| Chromium 实际 UI | mock 模式下 9101–9111 共 11 个页面全部与预期 data-outcome 一致；9109 报告错误与重试验收通过，无空报告错误替代。保存 JSON 与截图。 |
| `go test -p 1 -json ./...` | exit 1：639 pass、242 skip、37 fail（含子测试）。35 个数据库环境失败、1 个 Python 执行环境失败、1 个既有 Windows 路径行为失败。没有新增 skip 来掩盖失败。 |
| Linux PostgreSQL CI | 工作流已接入，**未执行**；不得当作通过。无自动推送来触发远端。 |

便携工具与依赖位于用户 `.codex/tools`，未修改全局 Go/npm 配置。前端用独立 npm cache 解决上一批 EEXIST；首次构建因本机 Next 遥测配置写入 EXDEV 失败，设置仅本次进程的 `NEXT_TELEMETRY_DISABLED=1` 后通过。未修改该全局配置文件。

## 失败记录与补正

保留首轮失败，不改写为成功：

1. 第一轮 race 检出测试观察计数器的读写缺少同步；改为 atomic.Int64 后重复通过。生产转换锁未被 race 报告。
2. 第一轮 Go 子进程无法在本机 `.git` 下创建证据子目录；换独立可写 artifact 目录，不关闭断言或增加 skip。
3. 第二轮全部 Go/Node 测试通过，但 runner 不能解析 test2json 分片输出；按 package/test 重组日志行，增加分片解析自测，重跑 20 次通过。
4. 浏览器 dev 模式观察到本机 AdGuard 注入脚本引发 hydration 警告；11 个目标结果卡片和报告错误断言通过，截图保留。没有因此更改 AdGuard 或浏览器全局配置。

## 本地证据位置

- `.git/artex-validation/phase1-full-go.jsonl`、`phase1-failure-classification.json`：全量结果及环境分类。
- `.git/artex-validation/phase1-build-vet.log`、`phase1-llmpool-race.log`、`phase1-web-tests.log`。
- `.git/artex-validation/phase1-web-build.log`：首次环境失败；`phase1-web-build-retry.log`：成功生产构建。
- `C:/Users/Kerna/.codex/tools/artex-engineering-phase1/`：分次保存的 runner 报告、Go/Node 日志和证据；最终通过目录见交付 manifest。
- `C:/Users/Kerna/.codex/visualizations/2026/11/02/01a199ad-c340-7273-bd3d-eb39a71d16f1/artex-ui-acceptance.json` 及 `artex-*.png`：UI 验收。

## 差异、回退与剩余门槛

三个变更组的补丁和逐文件 SHA256 位于 `.git/artex-validation/phase1-delivery/`。补丁以**已保留的上一批工作区状态**为基准，不把上一批功能混入回退；实际 manifest 是变更清单权威。

回退时先备份后续修改，然后对所选补丁运行 `git apply --reverse --check <patch>`；检查通过再 `git apply --reverse <patch>`。C 包含 CI/回归与前端展示，B 包含生命周期，A 包含参考/路线/交付文档。全部回退按 C → B → A，随后复跑上一批测试。不要 `git reset --hard`，不要 `git add -A`。

首次检出即缺失的三个文件保持原状，未恢复、未删除、未纳入本批补丁：

- `scripts/poc/php2shell/poc/verify_handler_mismatch.py`
- `scripts/poc/php2shell/poc/verify_session_injection.py`
- `scripts/poc/php2shell/poc/verify_uaf_asan.php`

发布/阶段最终验收仍需在 Linux CI 中完成 PostgreSQL 全量及暂停失败补偿、恢复/超时重新运行的集成测试；本机 Python 和 Windows 路径失败也需独立处置或取得目标平台通过证据。源码实现完成不等于这些门槛已通过。
