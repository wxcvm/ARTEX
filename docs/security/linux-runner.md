# Linux Docker 独立工具 Runner（首阶段，可选）

> 状态：首阶段集成，**不是完整沙箱**。默认的 `docker-compose.yml` 不变。仅在显式叠加 `docker-compose.runner.yml` 后启用。

## 为什么需要

原有 ARTEX 服务同时运行模型、数据库访问和 Bash/文件工具。Norma 的 Bash、后台任务和 PTY 会继承父进程环境；Read/Write 等文件工具也不限制绝对路径。单靠提示词输出标记不能阻止错误工具调用读取控制端凭证。

新增 Runner 采用单独容器和 Unix Socket：
- 控制容器保留数据库、模型和 JWT 凭证；
- Runner 只挂载 `data/tasks`、`data/sessions` 与私有 Socket，不挂载 `keys`、整个 `data` 或 Docker Socket；
- Runner 在另一个 Docker 网络，不能直接使用服务名 `postgres` 连接控制数据库；
- 独立代理桥仅将 `artex-proxy-relay:8788` 转发到 ARTEX 的流量代理端口；
- 内置 Bash、Read/Write/Edit/MultiEdit、LS/Glob/Grep、交互式 Shell（如被启用）和 TaskOutput/TaskStop/TaskList/Monitor 都通过 Socket 执行；
- 命令执行失败或 Runner 无法连接时**不回退到控制容器本地执行**；
- Docker 模式通过 Linux `SO_PEERCRED` 拒绝同一 PID 命名空间内的 Shell 直接调用 Runner RPC；控制端必须处在独立 PID 命名空间；
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
docker compose -f docker-compose.yml -f docker-compose.runner.yml exec -T tool-runner sh -ec '
  test -z "${ARTEX_PG_DSN:-}"
  test -z "${OPENAI_API_KEY:-}"
  test -z "${ANTHROPIC_API_KEY:-}"
  test -z "${ARTEX_LLM_SECRET_KEY:-}"
  test ! -e /app/keys/jwt.key
'
```

CI 还会通过 `scripts/ci/runner-relay-smoke.sh` 使用两个临时 Docker 网络和真实 `artex` 代理转发二进制验证：执行网络经指定 relay 能访问模拟控制端代理，不能直接访问控制端。这是网络路径的容器回归，**不代表正式 Compose 及真实流量抓包已验收**。

还需运行真实 ARTEX 任务，确认 Bash/文件操作、任务后台输出、交互式 Shell、代理抓包与安装的工具包能正常工作；只通过以上静态检查**不等于完成隔离验收**。

## 目前明确的限制

1. **Runner 暂时沿用容器内 root 用户**，主要是为了兼容主程序现有的 root 所有权工作目录；容器的 Linux capabilities 被移除，只有需要原始网络探测时保留 `NET_RAW`。后续还须实现统一文件属主、非 root UID 和每任务隔离。
2. `data/tasks` 与 `data/sessions` 目前是共享目录，而不是独立任务挂载。Runner 中的工具可能看到其他任务的文件。不得声称已经实现每任务文件系统隔离。
3. 控制端数据库领域工具、WebFetch/WebSearch、MCP/技能启动的自定义本地进程、浏览器工具以及其他非 SDK 内置工具不由此 RPC 自动隔离。启用这些能力时不能将此阶段当成完整凭证沙箱。
4. SDK 本地后台任务完成通知不能自动跨进程转发；可通过远程 TaskList/TaskOutput/Monitor 查询。会话结束时会请求清理 Runner 的后台进程和 PTY。
5. Runner 拥有常规探测网络权限。代理桥只是固定端口转发，不是严格的网络目标白名单；抓包与非 HTTP 探测需实机回归。部分工具依赖可写 HOME，Runner 默认映射到 tmpfs。
6. Unix Socket 采用文件权限 0600；隔离部署中还将拒绝本容器 PID 命名空间的直接调用，防止 Bash 通过 Socket 绕过控制端审计。此机制依赖 Linux PID 命名空间（不要设置 `pid: host` 或共享 PID 命名空间），不是独立的身份认证；Runner 中的 root 工具仍可能破坏 Socket 文件造成拒绝服务。服务重启、恶意符号链接与 Compose 部署仍需独立测试。

**后续目标**：完善工具覆盖、后台事件回传、每任务独立 UID/挂载和真正的非 root 执行，并进行真实 Docker 环境的凭证读取对抗测试。以上完成前不将 #159 标记为安全隔离验收通过。
