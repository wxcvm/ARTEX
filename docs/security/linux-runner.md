# Linux Docker 独立工具 Runner（首阶段，可选）

> 状态：首阶段集成，**不是完整沙箱**。默认的 `docker-compose.yml` 不变。仅在显式叠加 `docker-compose.runner.yml` 后启用。

## 为什么需要

原有 ARTEX 服务同时运行模型、数据库访问和 Bash/文件工具。Norma 的 Bash、后台任务和 PTY 会继承父进程环境；Read/Write 等文件工具也不限制绝对路径。单靠提示词输出标记不能阻止错误工具调用读取控制端凭证。

新增 Runner 把控制 Socket Broker 和工具 Executor 分成两个独立容器：
- 控制容器保留数据库、模型和 JWT 凭证；
- Broker 只挂载控制 Socket 和单独的执行 Socket，**不运行任何模型工具**；Executor 只挂载执行 Socket、`data/tasks` 和 `data/sessions`，**不挂载控制 Socket**，也不挂载 `keys`、整个 `data` 或 Docker Socket；
- Executor 只位于执行网络，Broker 只位于控制网络，二者不共享 Docker 网络；两者通过执行 Socket 通讯，ARTEX 与 Broker 通过另一条控制 Socket 通讯；
- 独立代理桥仅将 `artex-proxy-relay:8788` 转发到 ARTEX 的流量代理端口；
- 内置 Bash、Read/Write/Edit/MultiEdit、LS/Glob/Grep、交互式 Shell（如被启用）和 TaskOutput/TaskStop/TaskList/Monitor 都通过 Socket 执行；
- 命令执行失败或 Runner 无法连接时**不回退到控制容器本地执行**；
- Docker Broker 和 Executor 的 Socket 都通过 Linux `SO_PEERCRED` 拒绝本容器进程直接调用，减少工具命令绕过主程序审计的风险；执行容器里的 Bash 也无法访问或替换**控制 Socket**；
- 各 Agent 的权限、拦截规则、模型与任务范围不做改动。

## 启用与回退

先构建当前源码对应的 Linux 镜像：

```bash
./scripts/docker-build-fork.sh
```

保持原有 `.env`，启动时额外选择 Runner 配置：

```bash
docker compose -f docker-compose.yml -f docker-compose.runner.yml up -d
```

回退为原有模式：

```bash
docker compose -f docker-compose.yml -f docker-compose.runner.yml down
docker compose up -d
```

不额外购买 VPS。Windows 11 上的 Docker Desktop + WSL2 也可以运行 Linux 容器；**原生 Windows EXE 不在此阶段的隔离范围内**。

## 验收（在测试环境执行）

```bash
# 检查合成配置（密码仅设置为本地测试占位值，不要把正式密码粘贴进日志）
POSTGRES_PASSWORD=local-test-only docker compose -f docker-compose.yml -f docker-compose.runner.yml config --quiet

# 必须存在 Socket，Runner 不具有控制端数据库或模型密钥环境变量
docker compose -f docker-compose.yml -f docker-compose.runner.yml exec -T artex test -S /run/artex-runner/runner.sock
docker compose -f docker-compose.yml -f docker-compose.runner.yml exec -T tool-executor sh -ec '
  test -z "${ARTEX_PG_DSN:-}"
  test -z "${OPENAI_API_KEY:-}"
  test -z "${ANTHROPIC_API_KEY:-}"
  test -z "${ARTEX_LLM_SECRET_KEY:-}"
  test ! -e /app/keys/jwt.key
  test ! -e /run/artex-runner/runner.sock
'
```

CI 还会通过 `scripts/ci/runner-broker-smoke.sh` 在两个**没有共享网络**的 Docker 容器中验证：Broker 经执行 Socket 转发文件、Bash 与交互式 PTY 工具请求；Executor 看不到控制 Socket，也不能直接通过执行 Socket 伪造额外工具调用，且无工具管理 TCP 端口。这不等于正式 Compose 部署已验收。

CI 还会通过 `scripts/ci/runner-relay-smoke.sh` 使用两个临时 Docker 网络和真实 `artex` 代理转发二进制验证：执行网络经指定 relay 能访问模拟控制端代理，不能直接访问控制端。这是网络路径的容器回归，**不代表正式 Compose 及真实流量抓包已验收**。

CI 使用 `scripts/ci/runner-compose-smoke.sh` 在临时目录中按真实 `docker-compose.yml` + `docker-compose.runner.yml` 启动 Broker 与 Executor，验证 Compose 命名卷、健康检查和经 Broker 的读写操作。测试使用基于当前 Go 源码的**最小 CI 镜像**，不启动完整 ARTEX/数据库，也不验证正式镜像的全部工具依赖或真实代理抓包。

CI 另通过 `scripts/ci/runner-full-stack-smoke.sh` 以假数据库凭证和本地 HTTP 模拟目标，完整拉起 PostgreSQL、ARTEX API、Broker、Executor 和代理桥；测试执行网络的流量经过代理桥可达本地目标但不能直连。为避免大型镜像构建，它仍使用**精简 Go CI 镜像**，未覆盖正式工具镜像、HTTPS MITM 记录持久化或真实任务。

还需运行真实 ARTEX 任务，确认 Bash/文件操作、任务后台输出、交互式 Shell、代理抓包与安装的工具包能正常工作；只通过以上静态检查**不等于完成隔离验收**。

## 目前明确的限制

1. **Executor 暂时沿用容器内 root 用户**，主要是为了兼容主程序现有的 root 所有权工作目录；容器的 Linux capabilities 被移除，只有需要原始网络探测时保留 `NET_RAW`。后续还须实现统一文件属主、非 root UID 和每任务隔离。
2. `data/tasks` 与 `data/sessions` 目前是共享目录，而不是独立任务挂载。Executor 中的工具可能看到其他任务的文件。不得声称已经实现每任务文件系统隔离。
3. 控制端数据库领域工具、WebFetch/WebSearch、MCP/技能启动的自定义本地进程、浏览器工具以及其他非 SDK 内置工具不由此 RPC 自动隔离。启用这些能力时不能将此阶段当成完整凭证沙箱。
4. SDK 本地后台任务完成通知不能自动跨进程转发；可通过远程 TaskList/TaskOutput/Monitor 查询。会话结束时会请求清理 Runner 的后台进程和 PTY。
5. Runner 拥有常规探测网络权限。代理桥只是固定端口转发，不是严格的网络目标白名单；抓包与非 HTTP 探测需实机回归。部分工具依赖可写 HOME，Runner 默认映射到 tmpfs。
6. 控制 Socket 只挂载给 ARTEX 与 Broker，执行 Socket 只挂载给 Broker 与 Executor。Broker 与 Executor 不共享 Docker 网络；Executor 无法删除或替换**控制 Socket**，并且直接调用执行 Socket 会被 Linux 进程来源检查拒绝。但 Executor 内 root 工具仍可篡改、替换**执行 Socket 路径**造成拒绝服务或伪造响应，工作目录也被所有任务共享，因此仍需任务级 UID、独立文件权限与完整 Compose 现场测试。

**后续目标**：完善工具覆盖、后台事件回传、每任务独立 UID/挂载和真正的非 root 执行，并进行真实 Docker 环境的凭证读取对抗测试。以上完成前不将 #159 标记为安全隔离验收通过。
