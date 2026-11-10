# 运维文档索引（来自 0.4.2 发行包）

这 6 篇是发行包自带的运维手册，按"我现在遇到什么"来查。代码层面的功能说明在
[`docs/USER_GUIDE_ZH.md`](../USER_GUIDE_ZH.md)，两者重叠的部分以代码文档为准。

| 我遇到的情况 | 看这篇 |
|---|---|
| 要部署到公网 VPS / 想知道最低限度要做哪些加固 | [hardening.md](hardening.md) |
| 第一次装：Windows / Linux / Docker 怎么选、端口目录是什么 | [deployment.md](deployment.md) |
| 起不来、连不上数据库、端口被占、工具"消失" | [troubleshooting.md](troubleshooting.md) |
| 要升级，但不想丢数据 / 想知道哪些会被保留 | [upgrade.md](upgrade.md) |
| Linux 上选哪个包、架构对不对、依赖缺什么 | [linux-variant.md](linux-variant.md) |
| 这个版本改了什么 | [changelog-package.md](changelog-package.md) |

## 三条最容易被忽略的（来自 hardening.md 的真实事故）

1. **SSH 开密码登录**：实测一台机器 6.4 万次爆破尝试，出现过陌生 IP 的 root 会话挂 22 小时没人发现。改密钥登录，且**先验证公钥能登再关密码**。
2. **agent 起的靶场/OOB 服务全绑 0.0.0.0**：实测一次渗透期间起过 70+ 个无鉴权服务直接对公网可达。用完就关，或用 `iptables` 限来源。
3. **`config.json` 是 644**：里面是数据库口令，同机任何用户可读。`chmod 600 config.json`、`chmod 700 pgdata`。

## 与代码仓库的关系

- 面板默认只监听 `127.0.0.1:8787`（`-addr` 可改），Docker Compose 默认也只发布到回环；
- `knowledge/ops/hardening.md` 里提到的 `harden-server.sh` 属于发行包，**不在本仓库**；
- 本仓库的 `docker-compose.yml` 默认镜像已改为 `artex:local`（不再是上游作者镜像），部署前先 `scripts/docker-build-fork.sh`。
