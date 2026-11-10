<!-- 来源：ARTEX 0.4.2 发行包（Telegram 分发）根目录文档，原文逐字保留，仅转码为 UTF-8 并统一 LF。
     适用版本：0.4.2 发行包；概览见 knowledge/ops/INDEX.md；代码级说明见 docs/USER_GUIDE_ZH.md。 -->

# ARTEX 部署指南

从解压到跑起第一个任务。

---

# 一、前置要求

## 硬件

```
CPU     2 核起步，4 核舒适（扫描和 agent 并发吃 CPU）
内存    4 GB 最低，8 GB 推荐
磁盘    20 GB 起（工具链 + 数据库 + 任务产物；跑多了会涨）
网络    能出网即可（要调 LLM API）
```

## 软件

**Windows**：什么都不用装。包内自带 PostgreSQL 和 Python 运行时。

**Linux / macOS**：需要一个 PostgreSQL，`install-deps.sh` 会帮你装（需要 root 或 sudo）。

## 端口

| 端口 | 用途 |
| --- | --- |
| **8787** | 主服务，浏览器访问入口 |
| 17890 | 本地 SOCKS5（代理池启动后） |
| 17891 | mihomo 控制接口（仅本机） |
| 5432 或 15433 | PostgreSQL |

**8787 是唯一需要对外开放的端口**（如果要从别的机器访问）。

---

# 二、Windows 部署

## 包结构

```
ARTEX-0.4.1-Windows-x64-full\
├── start.bat            启动（双击）
├── stop.bat             停止
├── 覆盖升级.bat          升级到新版
├── 应用更新.bat          增量更新
├── 重置密码.bat          忘记密码时用
├── 部署说明.txt
├── scripts\
└── app\
    ├── artex.exe        主程序
    ├── config.json      配置（首次启动生成）
    ├── pgdata\          数据库（首次启动生成）
    └── data\            数据目录
```

## 步骤

**1. 解压**

解压到一个**路径不含中文和空格**的目录，比如：

```
D:\ARTEX\
```

不要放在 `C:\Program Files\` 下（权限问题），也不要放在桌面（路径太长容易出问题）。

**2. 启动**

双击 `start.bat`。

第一次启动会初始化数据库，要等 **1-3 分钟**。窗口里会滚动一堆 PostgreSQL 初始化日志，正常。

初始化完成后显示：

```
ARTEX 0.4.1 backend listening on 0.0.0.0:8787
[gate] 伪装门控已启用，入口路径: /g-xxxxxxxxxxxx
[custom-tool] 自动检测到 python 解释器: ...
```

**把那个 `/g-xxxxxxxxxxxx` 记下来，每次访问都要带它。**

**3. 打开控制台**

浏览器访问：

```
http://127.0.0.1:8787/g-xxxxxxxxxxxx
```

> **直接访问 `http://127.0.0.1:8787/` 会看到一个 nginx 默认页——那不是故障，是伪装门控。** 别人扫到你的端口看到的也是一样的页面。必须带 `/g-xxx` 才能进。

**4. 输入网关密码**

首次启动会生成一个随机密码，**就在 `start.bat` 窗口里**。找到这一行：

```
[gate] 访问密码: xxxxxxxxxxxx
```

输进去，进入控制台。

**密码丢了怎么办**：看 `app\data\gate.path`（里面有入口路径），密码可以用 `重置密码.bat` 重设。

---

# 三、Linux 部署

## 包结构

```
ARTEX-0.4.1-Linux-amd64\
├── artex                  主程序
├── start.sh               守护启动（崩溃自动重拉）
├── start-artex.sh         环境变量 + 启动
├── install-deps.sh        装依赖 + 初始化数据库
├── python-tools.sh        装 Python 工具链
├── reset-password.sh      重置密码
├── apply-update.sh        增量更新
├── harden-server.sh       服务器加固（可选）
├── defaults.sql
├── README-linux.txt
├── 如何选择Linux版本.txt
├── data\
└── redteam\
```

## 步骤

**1. 选对包**

```bash
uname -m
```

- `x86_64` → 下 `amd64` 包
- `aarch64` → 下 `arm64` 包

下错了报 `cannot execute binary file: Exec format error`。详见 `ARTEX-Linux版本选择.md`。

**2. 解压到固定位置**

```bash
mkdir -p /opt
tar -xzf ARTEX-0.4.1-Linux-amd64.tar.gz -C /opt
mv /opt/ARTEX-0.4.1-Linux-amd64 /opt/artex
cd /opt/artex
```

**3. 恢复执行权限**

压缩包经过 Windows 打包可能丢权限位：

```bash
chmod +x artex start.sh start-artex.sh install-deps.sh python-tools.sh
chmod +x *.sh
find data/tools -type f -exec chmod +x {} \;
```

**4. 装依赖 + 初始化数据库**

```bash
sudo ./install-deps.sh
```

它会做四件事：

```
装系统包（postgresql-16）
初始化数据库（initdb 到包内 pgdata）
生成 config.json（含随机数据库密码）
建库建表
```

**这一步会改系统 PostgreSQL 配置**，如果你机器上已经跑着别的 PostgreSQL，注意端口冲突。脚本会提示。

**5. 装 Python 工具链**

```bash
sudo ./python-tools.sh
```

装 impacket、pypykatz、sqlmap 这些 Python 工具到 `data/tools/python313/`。

约 5-10 分钟，看网络。

**6. 启动**

```bash
# 前台跑（Ctrl-C 停）
./start-artex.sh

# 或者后台常驻
nohup ./start-artex.sh > artex.log 2>&1 &
```

> **`start.sh` 和 `start-artex.sh` 的区别**：`start-artex.sh` 会先导出环境变量（`ARTEX_HOME`、`ARTEX_KB_DIR`、Python PATH），然后调用 `start.sh`。**日常用 `start-artex.sh`**。

启动日志会打印入口路径和密码。

**6.5 补装 Windows 独有工具（可选）**

Linux 版**不含** fscan / gogo / nxc / sqlmap（只有 Windows 包带）。需要的话：

```bash
# sqlmap
sudo apt install sqlmap

# netexec
pip3 install netexec

# fscan / gogo：从 GitHub releases 下 linux_amd64 版本
# 放到 data/tools/ 下，去掉架构后缀
```

详见 `ARTEX-Linux版本选择.md` 里「包内含什么 / 不含什么」。

**7. 配成系统服务（推荐）**

```bash
sudo tee /etc/systemd/system/artex.service > /dev/null <<'EOF'
[Unit]
Description=ARTEX redteam platform
After=network.target postgresql.service

[Service]
Type=simple
WorkingDirectory=/opt/artex
ExecStart=/opt/artex/artex -addr 0.0.0.0:8787
Restart=always
RestartSec=5
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now artex
sudo systemctl status artex
```

**8. 访问**

```
http://你的服务器IP:8787/g-xxxxxxxxxxxx
```

入口路径在 `data/gate.path` 里：

```bash
cat /opt/artex/data/gate.path
```

---

# 四、macOS 部署

和 Linux 基本一样，差异：

```bash
# 1. 确认架构
uname -m
# arm64   → Apple Silicon（M1/M2/M3）
# x86_64  → Intel

# 2. 解压
tar -xzf ARTEX-0.4.1-macOS-arm64.tar.gz
cd ARTEX-0.4.1-macOS-arm64

# 3. 权限
chmod +x artex start.sh start-artex.sh *.sh
xattr -dr com.apple.quarantine .     # 去掉 Gatekeeper 隔离标记

# 4. 装 PostgreSQL（用 Homebrew）
brew install postgresql@16
brew services start postgresql@16

# 5. 初始化
./install-deps.sh

# 6. 启动
./start-artex.sh
```

**macOS 特有的坑**：

- 下载的二进制会被 Gatekeeper 拦（"无法验证开发者"）→ `xattr -dr com.apple.quarantine` 或系统偏好设置里"仍要打开"
- 端口 8787 可能被其它服务占用 → `lsof -i :8787` 查

---

# 五、首次配置

进控制台后，三件事要配。

## ① LLM API（必配）

**位置**：系统 → API KEY 配置

填一个供应商的 key 就行。支持多家，不同任务可以挑不同模型。

```
常用：
  OpenAI 兼容接口（大多数中转站都是这个）
  Anthropic Claude
  国内各家（DeepSeek / 通义 / 智谱 …）
```

填完点「测试」验证能通。

**没有 LLM key，agent 跑不起来**——它是整个系统的脑子。

## ② IP 代理池（强烈建议）

**位置**：执行 → IP 代理池

两种类型选一种：

**订阅型**（机场节点）

```
贴订阅链接 → 系统自动拉节点 → 起本地 SOCKS5
```

**网关型**（住宅 IP 服务商，如 CLIProxy / IPIDEA）

```
网关地址 + 端口 + 账号 + 密码
国家地区：全球随机 / 指定
会话类型：粘性（同一 IP 干完一件事）/ 轮换（一个动作换一条）
```

**前置代理**：如果你在国内机器上跑，住宅网关会拒绝大陆 IP，需要填一个能出去的前置代理（比如本机 `http://127.0.0.1:7892`）。**海外服务器不用填。**

配完点「启用」，在「出口 IP」那栏看到非本机 IP 就通了。

## ③ 网关密码（可选改）

首次启动的随机密码在日志里。想改成好记的，用重置密码脚本：

```bash
# Linux / macOS
./reset-password.sh
```

```bat
:: Windows
双击 重置密码.bat
```

---

# 六、跑第一个任务

```
1. 工作台 → 仪表盘 → 快速开始
2. 填目标（域名或 IP，不要带 http://）
3. 选模式：
     我不确定全流程    ← 新手用这个
     先侦察
     拿权限
     打后台
     拖数据库
     留后门
4. 点「开始」
```

任务会自己跑：拆解目标 → 规划意图 → 派 worker 去干 → 验证 → 利用。

**你在「任务」页能看到每一步的实时日志**，在「漏洞发现」页看结果。

**第一个任务建议小范围**——单站、单域名。看它怎么侦察、怎么验证、怎么收口。看明白一次，后面就顺了。

---

# 七、常见问题

## 启动卡在初始化数据库

第一次启动要建库建表，1-3 分钟正常。超过 5 分钟没动静：

```bash
# Linux 看 PostgreSQL 是否起来
systemctl status postgresql
# 或
ss -lntp | grep 5432
```

Windows 上看 `app\data\pg_stderr.log`。

## 端口 8787 被占用

```bash
# Linux / macOS
lsof -i :8787
kill -9 <PID>

# Windows
netstat -ano | findstr :8787
taskkill /F /PID <PID>
```

或者换个端口启动：

```bash
./artex -addr 0.0.0.0:9000
```

## 访问显示 nginx 默认页

**这是正常的**。伪装门控在起作用。

必须带 `/g-xxxxxxxxxxxx` 路径访问，路径在：

```bash
cat data/gate.path
```

## 忘记网关密码

```bash
./reset-password.sh          # Linux / macOS
重置密码.bat                  # Windows
```

## 任务跑起来但 agent 不动

三个原因按顺序查：

```
1. LLM key 没配或配额用尽 → 系统 → API KEY 配置 → 测试
2. 代理池没启用，出不去网 → 执行 → IP 代理池 → 看出口 IP
3. 目标有 WAF/封了 IP → 看任务日志里的 403/验证码
```

## 磁盘越用越多

任务产物堆在 `data/tasks/`。每个任务可能几百 MB（扫出来的结果、下载的文件）。

清理：

```bash
du -sh data/tasks/* | sort -rh | head
```

删掉不要的任务目录即可，不影响其他任务。

## 想从公网访问

**不要直接把 8787 裸奔到公网。**

推荐加一层 nginx 反代 + HTTPS：

```nginx
server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate     /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    client_max_body_size 500m;

    location / {
        proxy_pass http://127.0.0.1:8787;
        proxy_buffering off;             # SSE 流式输出必须关
        proxy_http_version 1.1;
        proxy_read_timeout 3600s;
    }
}
```

伪装门控是靠**路径**判断的，不是靠来源 IP，所以走反代照样生效。

---

# 八、目录说明

```
config.json          数据库连接配置（备份它！）
jwt.key              登录密钥（备份它！）
pgdata\              数据库本体（备份它！）
license.dat          授权文件

data\tasks\          任务产物
data\traffic\        流量捕获
data\transcripts\    对话记录
data\tools\          工具链
data\gate.path       网关入口路径
data\artex-err.log   运行日志

redteam\             知识库（knowledge.db + cve.db）
```

**升级前必看**：`config.json`、`jwt.key`、`pgdata` 这三个不能被覆盖。详见 `ARTEX升级指南.md`。
