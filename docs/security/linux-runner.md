# Linux Docker 独立工具 Runner（首阶段，可选）

> 状态：首阶段集成，**不是完整沙箱**。默认的 `docker-compose.yml` 不变。仅在显式叠加 `docker-compose.runner.yml` 后启用。

## 为什么需要

原有 ARTEX 服务同时运行模型、数据库访问和 Bash/文件工具。Norma 的 Bash、后台任务和 PTY 会继承父进程环境；Read/Write 等文件工具也不限制绝对路径。单靠提示词输出标记不能阻止错误工具调用读取控制端凭证。

新增 Runner 把控制 Socket Broker 和工具 Executor 分成两个独立容器：
- 控制容器保留数据库、模型和 JWT 凭证；
- Broker 只挂载控制 Socket，**不运行任何模型工具**；Executor 只挂载 `data/tasks` 和 `data/sessions`，**不挂载控制 Socket**，也不挂载 `keys`、整个 `data` 或 Docker Socket；
- Executor 只位于执行网络，不直接接入控制端／PostgreSQL 网络；Broker 连接两侧，但只通过私有 Unix Socket 接收控制指令并转给指定 Executor；
- 独立代理桥仅将 `artex-proxy-relay:8788` 转发到 ARTEX 的流量代理端口；
- 内置 Bash、Read/Write/Edit/MultiEdit、LS/Glob/Grep、交互式 Shell（如被启用）和 TaskOutput/TaskStop/TaskList/Monitor 都通过 Socket 执行；
- 命令执行失败或 Runner 无法连接时**不回退到控制容器本地执行**；
- Docker Broker 通过 Linux `SO_PEERCRED` 限制 Socket 的来源；执行容器里的 Bash 根本无法访问或替换控制 Socket，避免通过 Socket 伪造工具调用；
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

CI 还会通过 `scripts/ci/runner-broker-smoke.sh` 真实启动两个独立 Docker 容器，验证 Broker 可以转发工具调用，Executor 无法读取或替换控制 Socket；这不等于正式 Compose 部署已验收。

CI 还会通过 `scripts/ci/runner-relay-smoke.sh` 使用两个临时 Docker 网络和真实 `artex` 代理转发二进制验证：执行网络经指定 relay 能访问模拟控制端代理，不能直接访问控制端。这是网络路径的容器回归，**不代表正式 Compose 及真实流量抓包已验收**。

还需运行真实 ARTEX 任务，确认 Bash/文件操作、任务后台输出、交互式 Shell、代理抓包与安装的工具包能正常工作；只通过以上静态检查**不等于完成隔离验收**。

## 目前明确的限制

1. **Executor 暂时沿用容器内 root 用户**，主要是为了兼容主程序现有的 root 所有权工作目录；容器的 Linux capabilities 被移除，只有需要原始网络探测时保留 `NET_RAW`。后续还须实现统一文件属主、非 root UID 和每任务隔离。
2. `data/tasks` 与 `data/sessions` 目前是共享目录，而不是独立任务挂载。Executor 中的工具可能看到其他任务的文件。不得声称已经实现每任务文件系统隔离。
3. 控制端数据库领域工具、WebFetch/WebSearch、MCP/技能启动的自定义本地进程、浏览器工具以及其他非 SDK 内置工具不由此 RPC 自动隔离。启用这些能力时不能将此阶段当成完整凭证沙箱。
4. SDK 本地后台任务完成通知不能自动跨进程转发；可通过远程 TaskList/TaskOutput/Monitor 查询。会话结束时会请求清理 Runner 的后台进程和 PTY。
5. Runner 拥有常规探测网络权限。代理桥只是固定端口转发，不是严格的网络目标白名单；抓包与非 HTTP 探测需实机回归。部分工具依赖可写 HOME，Runner 默认映射到 tmpfs。
6. Unix Socket 只挂载给 ARTEX 与 Broker。Executor 的工具进程无法删除或替换其路径；Broker 仍要求 Linux `SO_PEERCRED` 检查。但 Executor 内的 HTTP 工具管理服务还没有按任务隔离身份或文件权限，root 工具仍可能干扰其他任务。Socket 服务重启与完整 Compose 集成尚需验收。

**后续目标**：完善工具覆盖、后台事件回传、每任务独立 UID/挂载和真正的非 root 执行，并进行真实 Docker 环境的凭证读取对抗测试。以上完成前不将 #159 标记为安全隔离验收通过。
