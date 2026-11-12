# ARTEX 中文使用手册（安装、操作与排障）

> 适用仓库：[wxcvm/ARTEX](https://github.com/wxcvm/ARTEX) 的 `main` 分支，依据 2026-11-09 的实际代码编写。界面目前主要使用英文，本手册保留按钮原名以便对照。版本升级后若菜单有变化，请优先参考当前界面与仓库源码。
>
> **重要：使用的是你的 fork 还是历史上游镜像？** 仓库默认的 `docker-compose.yml` 使用 `autumn27/artex:latest`，**并不自动包含本 fork 新合并的功能**。想使用本文提到的 Skills 缺口、Lesson 经验评估、SARIF 等最新功能，请采用**从本仓库源码构建**的方式，或配置并运行你自行构建的 fork 镜像。不要以为执行 `git pull` 后再运行默认 Docker 镜像就完成了升级。

## 目录

1. [快速选择安装方式](#1-快速选择安装方式)
2. [Windows 11 从源码安装（推荐使用 fork 最新代码）](#2-windows-11-从源码安装推荐使用-fork-最新代码)
3. [Linux / macOS 安装](#3-linux--macos-安装)
4. [首次登录和界面导航](#4-首次登录和界面导航)
5. [配置 AI 模型（云端、Ollama、本地兼容接口）](#5-配置-ai-模型云端ollama本地兼容接口)
6. [第一个任务：从创建到查看结果](#6-第一个任务从创建到查看结果)
7. [资产、流量、证据和漏洞复测](#7-资产流量证据和漏洞复测)
8. [报告导出](#8-报告导出)
9. [Agent、Skills、MCP 和工具管理](#9-agentskillsmcp-和工具管理)
10. [历史经验（Lesson）与效果评估](#10-历史经验lesson与效果评估)
11. [设置、数据、更新和备份](#11-设置数据更新和备份)
12. [常见问题与检查清单](#12-常见问题与检查清单)
13. [其他参考文档](#13-其他参考文档)

本手册中的示例目标请替换为**你拥有或已获明确授权**的应用或测试环境。创建任务时应如实填写目标与工作范围；文档不会擅自为 ARTEX 加上新的程序级限制。

---

## 1. 快速选择安装方式

| 使用场景 | 推荐入口 | 是否包含本 fork 的最新提交 |
| --- | --- | --- |
| Windows 11，想用最新代码 | [第 2 节](#2-windows-11-从源码安装推荐使用-fork-最新代码)：PowerShell 构建 | **是**，从当前仓库编译 |
| Linux / macOS，想用最新代码 | [第 3 节](#3-linux--macos-安装)：本地编译 | **是**，从当前仓库编译 |
| 只想快速体验上游历史镜像 | `docker compose up -d` | **否**，默认镜像不是本 fork |
| 已经构建或发布了 fork 镜像 | `ARTEX_IMAGE` + `ARTEX_TAG` | **取决于镜像内容** |
| 使用预编译压缩包 | [Releases](https://github.com/wxcvm/ARTEX/releases) | 取决于发布版本的构建提交 |

**依赖清单：**

- PostgreSQL：项目的任务、资产、漏洞等主要数据存在 PostgreSQL 中。仓库的 Docker Compose 使用 PostgreSQL 16。
- Go：源码 `go.mod` 当前要求 **Go 1.26.3**；使用该版本或满足要求的更新版本。
- Node.js + npm：构建网页前端所需；项目 CI 使用 Node.js **22**。
- Git：拉取和更新仓库。
- Docker / Docker Compose：可选，用来启动 PostgreSQL 或运行容器部署。
- AI 模型服务：可在首次启动后由网页配置；无模型时无法完成依赖 LLM 的任务执行。

Docker 与源码构建属于不同安装路径；**不要同时让两个实例占用 8787 端口**。

## 2. Windows 11 从源码安装（推荐使用 fork 最新代码）

以下使用 PowerShell，适合想运行最新 `wxcvm/ARTEX` 的用户。准备 Git、Go 1.26.3+、Node.js 22、Docker Desktop（启用后端/WSL 2）。如果 PostgreSQL 已安装，也可以使用自己的数据库，无须运行第 2.2 步。

### 2.1 下载源代码并检查环境

打开 PowerShell：

```powershell
git clone https://github.com/wxcvm/ARTEX.git
cd ARTEX
git branch --show-current
go version
node --version
npm --version
docker --version
```

确认在 `main` 分支上。以后用 `git pull --ff-only` 获取最新提交。

### 2.2 启动本地 PostgreSQL

以下命令**仅用于创建一个独立的本地 PostgreSQL 容器**，不启动历史上游 ARTEX 镜像。请把示例密码改成自己的值：

```powershell
docker run -d --name artex-guide-pg -e POSTGRES_USER=artex -e POSTGRES_PASSWORD=请修改为自己的密码 -e POSTGRES_DB=artex -p 127.0.0.1:5433:5432 -v artex-guide-pg:/var/lib/postgresql/data postgres:16-alpine
docker logs artex-guide-pg
```

首次下载镜像后，等待数据库初始化完成。第二次开机如果容器已存在但没启动，使用：

```powershell
docker start artex-guide-pg
```

这里用主机端口 **5433**，对应容器内的 5432，避免与已经安装的 PostgreSQL 冲突。**不要重复执行 `docker run` 创建同名容器。**

### 2.3 创建本地配置

在仓库根目录执行：

```powershell
Copy-Item config.example.json config.json
notepad config.json
```

将示例文件改为下面这种**最小有效配置**，并填入第 2.2 步设置的密码：

```json
{
  "database": {
    "host": "127.0.0.1",
    "port": 5433,
    "user": "artex",
    "password": "请替换成自己的数据库密码",
    "dbname": "artex",
    "sslmode": "disable"
  },
  "skill_dir": "skills"
}
```

注意：仓库原始的 `config.example.json` 里 `skill_dir` 示例为 Linux 路径 `/opt/artex/skills`。**在 Windows 上应改为 `skills`、改成 Windows 有效路径，或删掉此字段让程序使用默认目录**。`config.json` 是运行配置，不要提交到公共仓库。

也可以用环境变量 `ARTEX_PG_DSN` 配置数据库；它优先于 `config.json` 的数据库字段。

### 2.4 编译并启动带网页的程序

在仓库根目录 PowerShell 执行：

```powershell
cd web
npm ci
npm run build:static
cd ..

New-Item -ItemType Directory -Force server\webui | Out-Null
if (Test-Path server\webui\dist) { Remove-Item server\webui\dist -Recurse -Force }
Copy-Item web\out server\webui\dist -Recurse -Force

$env:CGO_ENABLED = "0"
go build -tags embedui -o artex.exe ./cmd/artex
.\start.bat
```

看到服务启动后，用浏览器打开 **http://localhost:8787**。`start.bat` 会在异常退出时按退避间隔重启程序；正常停止可在对应控制台中操作。不要同时再运行第二份 `artex.exe`。

如果只想调试前端，可参考 README 的 `./dev.sh` 开发流程；正式想看到页面，编译时需要 `-tags embedui` 和已经复制好的前端静态文件。

### 2.5 Windows 预编译包（可选）

在 [本 fork 的 Releases](https://github.com/wxcvm/ARTEX/releases) 查找实际已发布、适合 Windows 的版本；**若没有新版本，不代表最新源码已有对应压缩包**。解压后保留 `artex.exe`、`start.bat`、`skills/`、`config.example.json` 等文件，配置有效的 `config.json`，然后运行 `start.bat`。首次使用仍需 PostgreSQL。

### 2.6 Windows 桌面客户端（便携版）

本 fork 新增原生 Windows 窗口 **ARTEX.Desktop.exe**，使用 Microsoft WebView2 显示现有 ARTEX 前端，不需要将网站安装为浏览器应用。独立窗口内保留任务、模型、Skills、报告、站内操作指南等现有功能，也可以继续用浏览器访问后端。

1. 打开本仓库 [Actions](https://github.com/wxcvm/ARTEX/actions) 的最新成功 CI，下载 `artex-windows-desktop-preview` 制品；或从正式 [Releases](https://github.com/wxcvm/ARTEX/releases) 下载 `ARTEX-Windows-Desktop-amd64.zip`。**只有相关 CI / Release 已成功生成时才会存在该下载项。**
2. 将 ZIP **完整解压**到可写目录（例如 `D:\Apps\ARTEX`），不要只复制桌面 EXE。
3. 如果没有设置 `ARTEX_PG_DSN`，复制 `config.example.json` 为 `config.json`，按照第 2.2、2.3 节配置一个已有的 PostgreSQL。
4. 双击 `ARTEX.Desktop.exe`。它会连接本机已运行的 ARTEX；若没有，则启动同目录的 `artex.exe` 并等待健康检查。**如果没有 `config.json` 且没有 `ARTEX_PG_DSN`，桌面客户端会打开首次数据库配置向导**，填入已有 PostgreSQL 的主机、端口、用户、密码和数据库名即可生成 `config.json`。手册第 2.2 节的 Docker 示例使用主机端口 **5433**；普通 PostgreSQL 的默认端口是 **5432**。向导会自动探测本机 PostgreSQL 默认端口 **5432** 和常用 Docker 映射端口 **5433**，可列出已安装的 Windows PostgreSQL 服务；你可以在明确确认后启动所选服务，再自动检测。端口探测只确认 PostgreSQL 服务响应，不会修改密码或初始化数据库。向导中的 **测试数据库连接** 会调用内置的 Go PostgreSQL 驱动校验主机、账号、密码及现有数据库，但不会创建数据库、迁移表结构或保存密码；需要时可跳过测试直接保存。向导还提供 **创建目标数据库** 按钮：使用已填写的账号连接 PostgreSQL 的管理数据库，仅在目标库不存在且账号具备 CREATEDB 权限时创建该库；不会初始化角色、重置密码或修改已有数据库。若启动失败或密码有误，可点击桌面工具栏的 **数据库配置** 再次编辑。修改现有 config.json 时会保留其它设置（例如 skill_dir），若桌面客户端自行管理后端，会在你确认后重启；外部服务须自行重启。对于自定义端口、Docker 或远程数据库，仍可直接手动填写。向导不会覆盖已存在的配置，也不自动静默安装 PostgreSQL。如果电脑没有 PostgreSQL，可以点击 **安装 PostgreSQL**：有 WinGet 时打开 EDB PostgreSQL 16 的交互安装流程（可能需要联网下载和 UAC 管理员授权），无 WinGet 时打开 [PostgreSQL 官方 Windows 下载页](https://www.postgresql.org/download/windows/)。在官方安装器里自行设置 `postgres` 用户密码；安装完成后返回向导，点击 **检测本机**，填写 `postgres` 和该密码。**推荐首次安装点击“初始化专用账户”**：使用此管理员连接新建普通 `artex` 登录角色及归它所有的目标数据库（仅限新建环境），生成随机密码并自动切换向导为 `artex`；随后点击 **保存并连接**。该专用角色没有超级用户、建角色或任意建库权限。若已有同名角色或数据库属于别人，初始化不会重置密码或转移数据库所有权，此时按现有设置选择使用原账号或手动配置。也可以继续使用 **创建目标数据库**、**测试数据库连接** 的旧流程。这是用户主动选择的安装入口，不会自动运行，也不会覆盖已有数据库服务。随后在窗口里完成管理员初始化及模型配置。
5. 窗口右上角关闭或最小化会收进系统托盘；通过窗口的 **退出** 或托盘菜单 **退出 ARTEX** 才真正关闭客户端。如果后端是由此客户端启动，退出也会停止后端，正在执行的任务可能中断；连接到已经启动的后端时则不会关闭它。

桌面版压缩包包含自包含 .NET 客户端、WebView2 组件、内嵌前端的 Go 后端、`skills/`、配置模板和中文手册；Win11 通常已有 WebView2 Runtime，若没有请从微软官方安装。桌面版**仍需独立 PostgreSQL**；不负责替你安装数据库。WebView2 用户数据和后端日志在 `%LOCALAPPDATA%\ARTEX\`，后端监听地址默认为 `127.0.0.1:8787`。

如果只想使用原浏览器版，继续执行 `start.bat` 打开 `http://localhost:8787`，两种方式互不替代。

### 2.7 Windows 安装版（可选 D 盘、升级和卸载）

除了上一节的便携 ZIP，本 fork 还提供 **`ARTEX-Windows-Setup-amd64.exe`** 标准 Windows 安装程序，基于 Inno Setup 6。CI 成功后可在 [GitHub Actions](https://github.com/wxcvm/ARTEX/actions) 下载 `artex-windows-setup-preview`；正式 tag 发布时可在 [Releases](https://github.com/wxcvm/ARTEX/releases) 获得安装版。**请确认你下载的是成功构建的安装包。**

1. 双击 `ARTEX-Windows-Setup-amd64.exe`，按向导选择安装目录。默认安装在当前用户的 LocalAppData Programs 目录；可自行改为 `D:\Apps\ARTEX` 等位置。普通用户安装不需要管理员权限。
2. 安装器创建开始菜单 **ARTEX** 快捷方式，可在向导里选择是否创建桌面快捷方式；在 Windows 设置的“已安装的应用”里可以正常卸载。
3. 安装完成后按第 2.3 节配置 PostgreSQL 的 `config.json`（或者提供 `ARTEX_PG_DSN`），再从开始菜单打开 ARTEX。安装器不包含数据库，也不启动自动后台服务；首次配置向导可以检测或启动现有 PostgreSQL 服务，但不会创建数据库或覆盖已有账户。
4. 后续运行新安装包会复用同一安装标识和原安装目录，替换程序文件，但**不会覆盖现有 `config.json` 和用户修改过的 `skills/`**。要获取更新的内置 Skills，可单独检查新版本内容并手动同步所需文件。
5. 卸载时会移除主程序、开始菜单快捷方式和安装记录，但保留已有配置与 Skills（以及用户数据目录）。需要彻底清理时，请自行检查原安装目录和 `%LOCALAPPDATA%\ARTEX\` 后再决定是否删除。

**便携版和安装版可以择一使用**：无需同时安装；浏览器方式也继续可用。安装包目前未配置 Authenticode 代码签名，因此 Windows 可能提示“未知发布者”；这并不代表程序已被签名验证。安装、升级前请确认下载来源为本仓库。

本地源码生成两种 Windows 包（需 Go 1.26.3、Node.js 22、.NET 10 SDK；安装版额外需要 Inno Setup 6）：

```powershell
./desktop/windows/package.ps1
./desktop/windows/build-installer.ps1 -Version 0.0.0-preview
```

### 2.8 Windows 外部工具检测（Tool Manager 第一阶段）

在 ARTEX 中打开 **系统 → Tools → External binaries**，点击“重新检测”即可查看 nuclei、httpx、subfinder、katana、ffuf、naabu、dnsx 七种常用工具在**后端主机**上的状态、版本输出、发现路径及官方项目入口。界面会区分“版本检查通过”“已发现但版本检查失败”“未发现”。如果你通过手机浏览器连接 ARTEX，这里仍检测 Windows / Linux 后端，而不是手机内部的软件。

检测目录按顺序为：

1. 可选环境变量 `ARTEX_TOOL_DIRS` 指定的多个目录，Windows 使用分号分隔（例如 `D:\Tools\Recon;E:\Security\bin`）；
2. `%LOCALAPPDATA%\ARTEX\tools`（Windows）；
3. ARTEX 程序目录下的 `tools/` 或 `data/tools/`；
4. 现有系统 `PATH`。

支持直接放入 `tools\nuclei.exe`，也支持将多个版本分开放入 `tools\nuclei\v3.11.1\nuclei.exe` 等层级。设置新的环境变量后，重启 ARTEX 后端再刷新页面。程序只对找到的可执行文件运行清单中指定的版本参数（每项最多 3 秒），不会自动扫描资产、修改全局 PATH 或下载工具。

**官方按需安装（Windows x64）**：在任一工具卡片点击“查看官方最新版”，程序向该工具官方 GitHub 仓库查询 latest release，展示版本、ZIP 大小以及 GitHub API 返回的 SHA-256 摘要。点击“下载并安装”或“更新到此版本”，确认后 ARTEX 才会下载 `windows_amd64.zip`，逐字节验证 SHA-256，且只提取 ZIP 中指定的 EXE。安装路径为 `%LOCALAPPDATA%\ARTEX\tools\<工具名>\current\`，不需要管理员权限，也不修改全局 PATH。更新会保留上一版在同级 `.previous\` 中，下载失败或校验失败不会替换当前版本。

**已有工具继续保留**：手动安装目录（特别是 `ARTEX_TOOL_DIRS`）的优先级高于 ARTEX 管理目录。即使安装了新版，如果界面仍显示自定义目录下的旧版，检查该环境变量即可；不会偷偷删除或更改已有文件。版本查询与下载需要 GitHub 网络可访问；遇到 HTTP 403 / 404、磁盘空间不足或 EXE 被占用时，页面会报告错误。新版本下载成功后点击“重新检测”，即可看到优先匹配的版本和受管理的版本。

**版本回滚与卸载**：有受 ARTEX 管理的版本时，工具卡片会显示管理版本；如果存在上一版，还会提供“恢复上一版”。点击并确认后，`current/` 与 `.previous/` 交换，旧的当前版本仍保留，可以再次恢复。点击“卸载管理版本”并确认后，会移除 ARTEX 管理目录中的当前版本和上一版，但不会删除 `ARTEX_TOOL_DIRS`、独立手动安装目录或系统 PATH 里的工具。放在受管理目录里的额外个人文件也会保留。如果 EXE 正在运行导致 Windows 拒绝移动或删除，请先退出该程序，再重试。刷新页面即可重新查看状态。当前不提供后台自动更新。

**离线工具包（Windows x64）**：如果需要在无法访问 GitHub 的 Windows 上准备工具环境，可在联网电脑的 ARTEX **System → Tools → External binaries** 页面先按需安装所需工具，再点击 **“导出已管理工具 ZIP”**。ARTEX 会把管理目录中七款工具的**当前版本**打包为 `ARTEX-Tools-Offline-windows-amd64.zip`；未安装的、仅在系统 PATH 或自定义目录中的工具不会被自动收集。也不会导出数据库、任务、凭据、Agent 配置或旧版本备份。

将 ZIP 用 U 盘等方式拷贝至离线电脑，安装并启动 ARTEX，打开同一页面，点击 **“导入离线工具 ZIP”**，选择文件并确认。后台先检查工具 ID、Windows x64 包格式、每个程序的 SHA-256 和实际大小；**全部通过才开始安装**。导入成功后会刷新版本与安装状态，原先受管理的版本仍可以用于回滚；现有 `ARTEX_TOOL_DIRS`、手动工具、PATH 和 Agent 绑定不会被覆盖。如果导入多个工具时，某个 EXE 正在运行导致 Windows 不允许替换，前面已成功导入的工具仍会保留，页面会提示错误并刷新状态，可关闭程序后重试。

离线包属于**工具数据包，不是完整 ARTEX 安装器**：离线电脑仍需要 ARTEX 本体和 PostgreSQL；也不包含所列工具运行时可能需要的额外组件（例如某些网络抓包驱动或模板库）。首次导出需要在联网机器上安装已有工具；完成导入本身不需要 GitHub。为避免内存或磁盘被无限占用，当前单个离线 ZIP 最大 1 GiB，清单只接受内置七款工具。支持带当前版本回执的 ARTEX 管理目录，不会扫描其他路径生成离线包。

**离线校验边界**：导出时计算 EXE 的 SHA-256，导入时重新验证可以证明与**离线 ZIP 内清单记录一致**；但清单本身没有独立数字签名，因而**不能**证明离线包或上游发布者的身份。请使用自己从可信 ARTEX 实例导出的文件，不要导入来源不明的工具包。原有 GitHub Release 下载校验仍单独适用。

**校验边界**：GitHub Release 资产 SHA-256 用于校验下载内容与 GitHub API 记录一致，但**不等于发布者代码签名，也无法单独证明上游仓库没有被入侵**。只支持清单中的七个官方仓库，不接受任意下载源。按需安装、回滚和卸载都不会自动扫描目标、注册 Agent 工具或改变已有任务权限。需要 Agent 使用外部程序时，仍在 Custom tools 中显式配置并绑定。



## 3. Linux / macOS 安装

### 3.1 最省步骤：仓库安装脚本

```bash
git clone https://github.com/wxcvm/ARTEX.git
cd ARTEX
./install.sh
```

脚本提供两种模式：

- **1 · All-Docker**：配置 PostgreSQL 密码，启动 Docker Compose。默认运行**历史上游镜像**，除非 `.env` 的 `ARTEX_IMAGE` / `ARTEX_TAG` 已指向你自行构建的 fork 镜像。
- **2 · Run locally (go build)**：连接已有 PostgreSQL 或通过 Docker 启动数据库，生成 `config.json`，在有 npm 的情况下构建嵌入式网页和 Go 程序。**这是直接运行当前仓库源码的便捷选项。**

该脚本是 Bash 脚本；Windows 原生 PowerShell 请按第 2 节操作。脚本提示安装 Docker 时会执行相关安装操作，先确认环境和权限。

### 3.2 手动构建 fork 代码

前提：有已连通的 PostgreSQL，以及配置好的 `config.json` 或 `ARTEX_PG_DSN`。

```bash
git clone https://github.com/wxcvm/ARTEX.git
cd ARTEX
cd web
npm ci
npm run build:static
cd ..
mkdir -p server/webui
rm -rf server/webui/dist
cp -r web/out server/webui/dist
CGO_ENABLED=0 go build -tags embedui -o artex ./cmd/artex
./start.sh
```

启动地址默认是 `http://localhost:8787`。自定义 HTTP 监听地址可以使用 `./start.sh -addr :9000`；流量记录代理独立使用 `-proxy` 参数，例如 `./start.sh -addr :8787 -proxy 127.0.0.1:8788`。

### 3.3 Docker Compose（镜像部署）

```bash
cp .env.example .env
# 编辑 .env，填写 POSTGRES_PASSWORD，必要时设置 ARTEX_IMAGE / ARTEX_TAG
docker compose up -d
docker compose ps
docker compose logs -f artex
```

`docker-compose.yml` 包含 `postgres` 与 `artex` 两个服务。默认将 HTTP 8787 和代理端口 8788 映射到宿主机回环地址；HTTP 宿主机绑定可通过 `ARTEX_BIND` 自定义。数据目录 `./data` 和 `./skills` 挂载到容器，数据库本体保存在 Docker 命名卷 `pgdata`。

**核对 HTTP 暴露状态**：登录后访问 `/api/privacy/exposure` 可查看 `configured_address`（启动参数）、`bound_address`（实际监听地址）、`binding_source`（`runtime_listener` 表示已读取真实监听器）及 `scope`。`listen_address` 在读取成功时也返回实际地址，否则返回配置值。该地址属于 ARTEX 进程所在的网络命名空间；容器内部绑定所有接口并不等于宿主机对公网开放。请在宿主机使用 `docker compose port artex 8787`、`docker compose port artex 8788` 和 `docker compose ps` 核对实际端口发布，并检查防火墙与反向代理。`public_reachability` 保持 `not_tested`，不代表已经验证公网不可达。

**可选的 Purge 隐私联动**：在 ARTEX 服务进程设置 `ARTEX_PRIVACY_BRIDGE_TOKEN`（建议使用随机生成的独立密钥，至少 32 个随机字节的编码结果），并在 DSH Purge 主机进程设置同值的 `DSH_PURGE_ARTEX_PRIVACY_TOKEN`。此密钥仅用于 `GET /api/privacy/bridge-status` 的 `X-Artex-Privacy-Token` 请求头，返回服务标识、监听范围、状态来源和公网未检测标记，不可用于登录或其他 ARTEX API。Docker Compose 从环境变量读取此可选配置；默认不启用。请勿使用管理员 JWT 或把密钥写入提交、报告或聊天；修改后重启相应进程。Purge 仅通过本机 HTTP 回环地址访问，不会跟随重定向。

**Docker 重建后保留登录状态**：本 fork 的新后端使用 `ARTEX_JWT_KEY_DIR=/app/keys`，由 Compose 挂载仓库根目录的 `./keys`，单独持久化 JWT 签名密钥。`keys/` 不会被 Git 跟踪，也不进入 Docker 镜像构建上下文；备份时应保护该目录。Linux/macOS 可先执行 `mkdir -p keys`，Windows PowerShell 可执行 `New-Item -ItemType Directory -Force keys`。新后端会自动搬迁仍在挂载的 `data/jwt.key`，但**旧容器内部 `/app/jwt.key` 在容器重建后可能无法恢复**。如果已有旧容器，并希望保留现有登录令牌，应在删除或重新创建旧容器前将其密钥复制到宿主机的 `keys/jwt.key`（例如 `docker cp <旧容器名>:/app/jwt.key ./keys/jwt.key`）；只在原文件确实存在时操作，**不要生成不同的密钥覆盖它**。如果旧容器已经删除且密钥没有备份，只能重新登录。此行为仅适用于运行本 fork 新后端的 Docker 镜像；当前 Compose 的默认历史上游镜像不会因此自动更新到本 fork。


**特别提醒：** 只设置 `ARTEX_IMAGE=wxcvm/ARTEX` 并不意味着镜像仓库一定存在，也不意味着镜像已经包含最新代码。只有确认自行构建或发布的镜像名称与标签后，才能把它们写入 `.env`。

如将 AI 服务运行在宿主机，而 ARTEX 运行在 Docker 容器，容器内部的 `127.0.0.1` **不是宿主机**。请使用容器实际可访问的主机名/IP，并用模型连接测试验证。

## 4. 首次登录和界面导航

1. 浏览器打开 `http://localhost:8787`（或你配置的域名和端口）。
2. 首次安装进入 `/setup` 初始化管理员密码。之后通过登录页使用该账号；**不要寻找所谓“默认管理员密码”**。
3. 左侧菜单的英文名称与主要功能如下；可用浏览器内置翻译辅助阅读。
4. 在左侧 **System → 使用指南 / Help** 可以直接打开这份中文手册的站内版，使用顶部搜索框或章节目录定位说明。网页构建时从本文件生成，更新本文件后需要重新构建前端才能同步到已部署实例。
5. 在模型、任务、MCP、Skills 等页面的顶部点击 **操作指南**，会跳转到最相关的操作章节；任务详情的顶部问号也可直接打开对应说明。站内首页另有首次使用路线：安装 → 模型 → 任务 → 结果 → 导出。

| 菜单 | 功能 |
| --- | --- |
| **Dashboard** | 总览、活动和 Token 统计 |
| **Chat** | 与系统 Agent 对话、查看交互 |
| **Tasks** | 创建、管理、暂停/恢复任务，查看执行过程和结果 |
| **Findings** | 集中查看、筛选、管理漏洞并导出 |
| **Traffic** | 查看网络流量记录 |
| **Tool Execution** | 查看工具执行记录 |
| **LLM Recording** | 查看模型调用记录 |
| **Assets / Asset Sync** | 资产图谱、同步 ScopeSentry 等外部资产来源 |
| **Workspace** | 管理任务工作目录中的文件 |
| **System → LLM** | 配置模型 Profile、连接测试、激活模型与故障转移 |
| **System → Agent / Skill / MCP / Tools** | 管理 Agent、技能、MCP 服务和工具 |
| **System → Interception Rules / Approval Records** | 查看拦截规则与审批记录 |
| **System → Notifications / Logs / System Settings** | 通知、日志与系统设置 |

界面截图可参考根目录 `screenshots/`，例如 [任务列表](../screenshots/tasks.png)、[模型配置](../screenshots/llm.png)、[Agent 管理](../screenshots/agents.png)、[资产](../screenshots/assets.png)。

## 5. 配置 AI 模型（云端、Ollama、本地兼容接口）

### 5.1 创建模型 Profile

1. 打开 **System → LLM**（`/system/llm`）。
2. 点击新增模型配置（界面使用 **New model profile**）。
3. 填写 **Name**（便于识别的名字）、**Format**（协议类型）、**Model**（模型 ID）、**Base URL**（需要时填写）和 **API Key**（服务要求时填写）。
4. 使用 **Test connection** 测试连通性。
5. 如果中转网关需要自定义请求标识，在 **Custom User-Agent (optional)** 填写其要求的 User-Agent；留空保持原有默认值。该字段会随 Profile 保存，并用于连接测试，不接受无效 HTTP 控制字符。
6. 点击 **Create / Save** 保存，再在新模型卡片上点击 **Set as active** 激活。**创建 Profile 不会自动激活。**
7. 创建任务时还可以为任务选择自己的模型 Profile；Agent 也可以单独绑定模型。

支持的 Profile 格式包括 `anthropic`、`openai`、`openai-responses` 和 `ollama`。不同模型支持的上下文长度、思考参数和 Token 参数不同；先用对应服务的实际接口确认模型名与协议，再在 ARTEX 中测试。

### 5.2 云端模型示例

- OpenAI 兼容接口：`Format = openai`、`Model = 服务实际提供的模型 ID`、`Base URL = 服务实际的 API 前缀`、`API Key = 你的密钥`。
- Anthropic 接口：`Format = anthropic`，填写对应模型和 Key。
- Responses API：只有服务实际实现该接口时才选 `openai-responses`。

**Base URL 填写规则（常见连接失败原因）：**

| Profile Format | Base URL 示例 | ARTEX/模型 SDK 请求时追加的路径 |
| --- | --- | --- |
| `anthropic` | `https://api.anthropic.com`（末尾不要 `/v1`） | `/v1/messages` |
| `openai` | `https://api.openai.com/v1` | `/chat/completions` |
| `openai-responses` | `https://api.openai.com/v1` | `/responses` |
| `ollama` | `http://127.0.0.1:11434/v1` | `/chat/completions` |

例如 Anthropic 的 Base URL 若错误地填为 `https://api.anthropic.com/v1`，实际路径就可能变成 `/v1/v1/messages`。使用其他云平台或中转服务时，以它实际兼容的协议和 API 前缀为准，不要将以上域名当作所有服务的通用地址。

环境变量 `OPENAI_API_KEY` / `ANTHROPIC_API_KEY` 也可以用于初始化默认模型配置；网页 Profile 是另一个配置入口。不要把 Key 发在 Issues、截图或日志里。

**可选：模型 API Key 静态存储加密。** 在 ARTEX 后端进程的环境中设置 `ARTEX_LLM_SECRET_KEY`，值为独立随机生成的 32 字节密钥（64 位十六进制字符）。Docker Compose 也支持从环境或本机 `.env` 传入；不要提交该文件或把密钥发给他人。重启后，新建或明确重新填写 API Key 并保存的模型配置将使用 AES-256-GCM 加密存入 PostgreSQL；界面仍显示脱敏提示，模型调用时在后端解密。已有明文记录继续兼容，不自动迁移；编辑模型参数时留空 API Key 保持原有存储方式。需要迁移某个旧配置时，在确认已备份加密密钥后重新填写该模型的 API Key 并保存。

未配置此变量时保持原有明文写入方式；已加密记录必须使用对应密钥读取，密钥缺失或错误时不会返回密文作为 API Key，也不会自动改回明文。请单独备份原密钥并保持使用同一数据库的各 ARTEX 实例配置一致；错误格式会在启动时提示。不要将加密密钥与数据库备份放在一起。启用后如需回退到旧版本，请先评估旧版本无法读取加密模型密钥的影响。该功能仅保护模型配置字段的静态存储，不覆盖旧数据库备份、WAL、其他凭据、运行时内存或能够读取后端进程环境的主体。启用加密也不代表旧明文副本已经消失。

在「模型配置 → 密钥保护」点击「检查密钥保护」，可以查看加密密钥是否配置、每个模型配置的存储状态及当前密钥认证结果。仅明文、仅密文、明文与密文并存、未保存密钥分别统计；已验证、认证失败、未校验另行统计。未校验不代表数据损坏，空密钥可能是正常的本地模型配置。诊断只读取数据库，不修改模型，不返回密钥或密文，不触发模型服务请求。结果为本次检查快照，编辑模型配置后可再次检查。

### 5.3 Ollama 和本地模型

如果 Ollama 与 ARTEX **都在同一台原生系统上运行**，可以先在命令行确认模型服务已启动并已加载模型，再创建 Profile：

| 字段 | 示例 |
| --- | --- |
| Name | `Local Ollama` |
| Format | `ollama` |
| Base URL | `http://127.0.0.1:11434/v1` |
| Model | Ollama 实际列出的模型名称 |
| API Key | 本地 Ollama 可以为空 |

对于 LM Studio 或其他 OpenAI 兼容服务器，选择它真正实现的协议与 API 地址；具体模型 ID、Key 是否必需、端口以服务的设置页为准。

如果提示连接失败，先分清是 **服务未启动、模型未加载、模型名错误、端口不通、协议不匹配**，还是 **容器无法访问宿主机**。在 Docker 模式下不能照抄宿主机的 `127.0.0.1` 地址。

### 5.4 多模型与故障转移

**LLM Pool / Failover** 支持为未固定模型的 Agent 配置故障转移。至少有两个可用且参与池的模型配置时才有实际切换效果；显式绑定模型的 Agent 或任务不一定进入池。首次使用建议先验证单个 Profile 能正常响应，再设置模型池。

## 6. 第一个任务：从创建到查看结果

以下流程说明**如何使用界面**，不要求安装新插件。示例目标应替换成你拥有或获授权的对象。

### 6.1 创建任务

1. 进入左侧 **Tasks**（`/function/tasks`）。
2. 点击右上角 **New task**，打开右侧任务抽屉。
3. 填写任务名称（Name）、测试目标描述（Description）和希望完成的目标（Goal）。
4. 按需要选择任务类别（Category）、公司/资产关系（Company/Assets）、任务模板及模型 Profile。没有现成模板也可以直接填写。
5. 展开高级选项，按任务情况决定是否使用现有的时长、并发、初始任务、资产覆盖或其他已有设置；不熟悉时保留原值。
6. 提交创建后回到任务列表，点开对应任务观察状态和执行记录。

示例文本（用于**授权测试任务描述**，不是额外的程序限制）：

- **Name**：`我的测试站点 · API 评估`
- **Description**：`检查已授权的测试 API 服务，记录能验证的输入、响应和证据；测试范围由已绑定资产和任务描述确定。`
- **Goal**：`梳理接口和访问控制检查结果，输出可复现的发现与处理建议。`

### 6.2 状态与并发

任务列表显示 `Created`、`Queued`、`Running`、`Paused`、`Completed`、`Failed`、`Timed out` 等状态。**Queued** 往往意味着当前任务并发额度已满，等待空位；**Paused** 可使用相应恢复操作；**Failed** 或 **Timed out** 应打开任务详情与日志定位问题，而非直接推断目标存在漏洞。

任务列表提供搜索、状态/类别筛选、排序以及批量控制；任务并发上限可通过任务页面的 **Concurrency** 相关入口配置。注意不要把“任务数并发”与“每个任务的 Worker 数量”混淆，后者在 **System Settings** 中设置。

### 6.3 任务详情页逐项看

点开某个任务，右侧或详情页有这些标签（实际按英文显示）：

| 标签 | 应该看什么 |
| --- | --- |
| **Sessions** | Planner、Worker、模型交互、执行过程 |
| **Overview** | 任务概要和进展 |
| **Exploration graph** | 探索目标、意图、事实与发现之间的关联 |
| **Broadcast** | 任务内部的共享信息/广播 |
| **Findings** | 当前任务发现的问题 |
| **Retests** | 独立复测、结论和历史记录 |
| **Test assets** | 本任务涉及的资产 |
| **Credentials** | 任务使用的授权凭据配置 |
| **Asset coverage** | 资产覆盖关系及测试进度 |
| **Interception approval** | 需要审查的操作与审批记录 |
| **Report** | Markdown 任务报告 |

任务结果摘要卡显示发现数量、已确认项、待复核项、工具调用、调用过的 Skills 以及 Lesson 注入记录。对已完成的任务，可以在卡片上点击 **Compare similar completed tasks** 查看历史经验的观察性对比（见第 10 节）。

## 7. 资产、流量、证据和漏洞复测

### 7.1 Assets：资产图谱与范围

进入 **Assets** 查看域名、子域名、IP、服务、应用、接口等资产及关联信息。任务详情里的 **Test assets** 与 **Asset coverage** 更偏向该任务的实际关联资产和覆盖情况。

如果已有 ScopeSentry：打开 **Asset Sync**，填入 ScopeSentry 的地址与 API Key，根据页面选项选择项目或任务及要同步的资产类型，再执行导入。同步完成后，到 **Assets** 和相关任务详情核对结果。

### 7.2 Traffic：流量记录

ARTEX 使用独立流量记录代理，默认参数为 `-proxy 127.0.0.1:8788`，与 Web/API 的 8787 端口不是同一个端口。需要代理时先确认进程在监听，再按你的客户端设置相应代理地址。使用 HTTPS 流量记录时，证书信任与浏览器配置应依照实际部署环境进行。

另外仓库带有 `extension/` 浏览器扩展。可以在 Chrome/Edge 扩展管理页开启开发者模式后加载该目录，在扩展弹窗配置 ARTEX 地址及有效会话 Token，针对正在使用的标签页启停捕获。文本请求/响应会进入本地流量存储，二进制响应可能被跳过。**不要把浏览器 Token 写进文档、截图或公开仓库。**

### 7.3 Findings：漏洞记录与证据

进入 **Findings**，可以根据页面的过滤项定位漏洞。进入详情查看影响资产、描述、状态以及流量证据的关联；若系统尚未确认问题，不要将其当成已验证结论。

文档 [漏洞与流量证据关联](vulnerability-traffic-evidence.md) 介绍流量证据的详细处理方式。

### 7.4 Retests：复测和结论

在任务详情的 **Retests** 标签或单条漏洞详情的复测入口启动独立复测。可填写修复版本、测试条件等附加信息，随后前往相关 Session 查看过程。常见结果为 **Still reproducible**、**Fixed**、**Cannot confirm**。复测结果与原有漏洞记录关联，不能把“无法确认”自动解释为“漏洞已经修复”。

## 8. 报告导出

### 8.1 按任务查看 Markdown 报告

打开 **Tasks → 任务详情 → Report**。这里会尝试加载 Markdown 报告，可以使用 **Refresh report / Retry report** 更新展示，或者点击 **Copy** 复制当前报告。

**没有报告文件 ≠ 没有漏洞；报告加载失败 ≠ 任务失败。** 先看 Sessions、Findings 和日志是否有相应结果。

### 8.2 从 Findings 批量导出

打开 **Findings → Export findings**，按页面选项选择范围：

- 当前筛选结果（Filtered）；
- 全部记录（All）；
- 手动选中的条目（Selected）。

可选格式包括 **Markdown 单文件、Markdown ZIP、CSV、JSON、SARIF 2.1.0**。确认后点击 **Export**，浏览器下载生成的文件。CSV 适合表格分析，SARIF 适合支持该格式的代码安全/结果处理系统；它们不是漏洞复测的替代品。

## 9. Agent、Skills、MCP 和工具管理

### 9.1 Agent

在 **System → Agent** 管理 Planner、Worker 及其他 Agent 的配置，查看可编辑的提示词、模型绑定和工具等选项。修改配置后可以先创建一个独立任务确认实际行为；并不是所有已经在运行的 Session 都会立即使用新配置。

### 9.2 Skills：创建、上传、编辑与调用统计

在 **System → Skill**（`/system/skills`）：

1. 点击 **Create skill**；填写名称、描述、可选说明内容、关联 MCP 和可见 Agent，保存后检查该 Skill。
2. 如已有 Skill 包，点击 **Upload archive** 上传包含 `SKILL.md` 的 ZIP。
3. 点击左侧 Skill 树，查看/编辑技能文件，检查最近调用与 Agent 可见性。
4. 在缺失 Skills 区域查看模型实际尝试调用却不存在的名称；点击该名称可打开预填名称的创建草稿。**草稿不会自动成为已审核、已配置好的技能。**

代码库 `skills/<skill-name>/SKILL.md` 是单个 Skill 的主要描述文件；本地文件系统和界面编辑能力以当前部署模式为准。历史调用记录不会因为 Skill 安装或卸载而被清空，但“当前缺失”列表会排除已正常加载的 Skill。

### 9.3 MCP：外部工具接入

在 **System → MCP** 点击 **Add MCP**。页面支持本地 `stdio` 和远程 `http`（Streamable HTTP），并可查看/刷新已发现的 Tools、设置服务启用状态和 Agent 可见性。旧式 SSE 服务还可按后端支持情况配置 `sse` 传输；实际协议必须与 MCP 服务端匹配。

常见排查顺序：确认服务地址/命令正确 → 服务端确实启动 → 传输协议匹配 → 认证信息正确 → 在 ARTEX 中刷新工具列表。**工具已注册不代表所有 Agent 都能看见**，还要检查可见性设置。

### 9.4 其他系统入口

- **System → Tools**：管理工具目录。
- **System → Interception Rules / Approval Records**：查看现有审批相关设置和记录。
- **System → Notifications**：通知设置。
- **System → Logs**：检查后端日志；遇到模型连接、任务或工具错误优先查看此处。
- **System → System Settings**：Worker 数量、网络/搜索、流量捕获以及其他全局设置。

## 10. 历史经验（Lesson）与效果评估

本 fork 已有可追溯的经验存储，基于已有任务的 Findings 和 Retests 提取候选经验；**无需另装记忆插件**。

### 10.1 三种模式

环境变量 `ARTEX_RSI_LESSONS` 控制启动时的经验模式：

| 模式 | 行为 |
| --- | --- |
| 未设置 / `off` | 默认关闭经验选择与注入 |
| `shadow` | 选择并记录本来会显示的经验，但不注入提示词 |
| `inject` | 选择、记录并把经验块加入 Planner / Worker 提示词 |

在 Linux/macOS 本地启动前，例如：

```bash
export ARTEX_RSI_LESSONS=shadow
./start.sh
```

在 Windows PowerShell 本地启动前：

```powershell
$env:ARTEX_RSI_LESSONS = "shadow"
.\start.bat
```

**Docker 注意事项：** 当前 `docker-compose.yml` 没有直接转发 `ARTEX_RSI_LESSONS`。仅把这一行写进 `.env` 不会自动传入容器；若确实需要在 Docker 中使用，应由部署者把变量加入 `artex` 服务的 `environment` 并重建/重启容器。改动环境变量影响新建或重新构建 Agent 的任务，不会重写旧任务的使用记录。

### 10.2 任务经验与证据

项目会在任务到达结束状态后提取可溯源的候选 Lesson；原始事实可追溯到 Findings 或复测记录。可通过现有接口查看或重新提取：

```text
GET  /api/tasks/{id}/lessons
POST /api/tasks/{id}/lessons/extract
GET  /api/tasks/{id}/outcome
```

这些是需要登录授权的 API；`{id}` 替换成真实任务 ID。不要在公共命令或截图中粘贴会话 Token。

### 10.3 效果比较：怎样使用

1. 进入一个 **Completed** 的任务详情。
2. 在上方任务结果卡中点击 **Compare similar completed tasks**。
3. 系统按**相同公司与任务类别**查询已完成、未归档任务，并按使用记录分为：
   - **Lessons injected**：只记录了注入模式；
   - **Shadow only**：只记录了 shadow 模式；
   - **Both modes**：同任务出现两种模式；
   - **No exposure recorded**：没有 Lesson 使用记录。
4. 阅读每组任务数、每任务有效发现、误报占已审查结果的比例、每有效发现消耗的工具调用与 Token、每任务 Skill 加载数、任务耗时。
5. 没有可比任务时显示空结果；有效发现数为零时，每有效发现的消耗指标显示 **—**，不会伪造一个效率值。

对应 API：`GET /api/tasks/{id}/lessons/evaluation`。查询是**点击按钮后才进行**，不会随着结果卡的周期刷新自动重复计算。

**怎样理解结果？** 这些数字描述的是历史任务记录，不是随机对照实验；模型、任务难度、目标覆盖面、任务配置都可能不同。“没有记录”不等于真正的未使用组，“用过 Lesson”也不表示 Lesson 导致漏洞被发现。当前功能不会根据比较结果自动采用、淘汰或修改任何 Lesson。

详细原理参见 [RSI Lessons 技术说明](rsi-lessons.md)。

## 11. 设置、数据、更新和备份

### 11.1 主要文件和数据去向

| 位置 | 用途 |
| --- | --- |
| `config.json` | 本地数据库连接与 Skill 路径 |
| `.env` | Docker Compose 环境变量、数据库密码、镜像参数 |
| `data/` | 本地运行状态、文件与部分运行记录 |
| `skills/` | 本地技能目录 |
| PostgreSQL | 任务、资产、漏洞、经验等关系数据 |
| Docker 卷 `pgdata` | Compose 使用的 PostgreSQL 持久化数据 |

不要把密码、Token 或包含敏感测试证据的文件直接提交到公开仓库。

### 11.2 更新 fork 源码版

本地源码构建方式：

```bash
git pull --ff-only
cd web && npm ci && npm run build:static && cd ..
rm -rf server/webui/dist
cp -r web/out server/webui/dist
CGO_ENABLED=0 go build -tags embedui -o artex ./cmd/artex
# 停止旧进程，再重新运行 ./start.sh
```

Windows 依次运行 `git pull --ff-only`、`npm ci`、`npm run build:static`，重新复制静态文件，并按第 2.4 节编译 `artex.exe`，停止旧进程后运行 `start.bat`。

Docker 模式可使用 `./update.sh` 或手动 `docker compose pull artex` / `docker compose up -d artex`；**该过程更新的是已配置的镜像，不会自动把本地 Go 源码编进容器**。如果镜像仍是 `autumn27/artex`，升级得到的仍是历史上游镜像线。

数据库 schema 会在服务启动时按项目现有机制执行幂等初始化/更新。即使如此，涉及重大更新前也应保存数据库和目录副本。

### 11.3 备份和恢复

Docker Compose 部署可运行：

```bash
scripts/backup.sh --with-data
# 或只备份数据库并指定保留数量
scripts/backup.sh --keep 7
```

`--with-data` 保存数据库转储和 `data/` 压缩包；默认备份位置为 `backups/`。数据库恢复须在确认目标数据库与备份文件匹配后执行，可参考 [README 备份章节](../README.md#backups) 里的 `gunzip | docker compose exec -T postgres psql ...` 示例。**不要误用 `docker compose down -v`：`-v` 会移除数据库卷**。

## 12. 常见问题与检查清单

| 现象 | 先检查什么 |
| --- | --- |
| 打不开 `localhost:8787` | 服务进程/容器是否运行、8787 是否被占用、浏览器访问的是正确机器 |
| 启动报数据库错误 | `config.json` 或 `ARTEX_PG_DSN`、PostgreSQL 是否 ready、端口和账号是否正确 |
| Windows 提示 `artex.exe` 找不到 | 是否在仓库根目录编译、运行了 `go build`、是否与 `start.bat` 在同一目录 |
| 浏览器只有 API 或没有前端 | 是否执行 `npm run build:static`、复制 `web/out` 到 `server/webui/dist`，以及是否带 `-tags embedui` 编译 |
| Docker 看不到最新 fork 功能 | `.env` 的 `ARTEX_IMAGE` / `ARTEX_TAG` 仍指向历史上游；需使用实际从 fork 构建的镜像或本地源码版 |
| 模型显示 No Key | 非 Ollama Profile 是否填写了有效 Key；本地服务是否选择了匹配格式 |
| 模型 Test connection 失败 | API 地址、模型 ID、协议格式、服务运行状态及 Docker 内外网络 |
| 任务一直 Queued | 查看 Tasks 页面任务并发额度和已有 Running 任务 |
| 任务 Failed / Timed out | 看详情 Sessions、LLM Recording、Tool Execution、System → Logs |
| Skill 不出现在列表 | `SKILL.md` 是否可解析、目录位置正确、部署的 `skill_dir` 是否仍指向旧路径 |
| MCP 有服务但无工具 | 服务器运行情况、传输类型、刷新 Tools、Agent 可见性 |
| Traffic 没有记录 | 是否启动代理/扩展并设置正确端口、HTTPS 捕获配置、是否确实有请求进入 |
| Report 没内容 | 先核对 Findings、Sessions；无报告文件与无漏洞不是一回事 |
| 比较面板没数据 | 同公司+同类别下是否有 Completed 且未归档的任务 |
| 设置 `ARTEX_RSI_LESSONS` 后无变化 | 是否重建 Agent、是否只写入 `.env` 而没有传到 Docker 容器 |

**最小排查命令：**

Windows PowerShell：

```powershell
go version
docker ps
docker logs artex-guide-pg
Get-NetTCPConnection -LocalPort 8787 -ErrorAction SilentlyContinue
```

Linux/macOS 的 Compose 部署：

```bash
docker compose ps
docker compose logs --tail=100 artex
docker compose logs --tail=100 postgres
```

想确认服务版本，可打开页面并检查启动日志中的 `version` / 运行平台信息；源码编译在未传入 release 版本参数时可能标为 `dev`，**不能单凭 `dev` 就认定不是最新代码**。

## 13. 其他参考文档

- [README：架构、部署、开发说明](../README.md)
- [RSI Lessons：经验提取、使用模式与效果评估](rsi-lessons.md)
- [漏洞与流量证据关联](vulnerability-traffic-evidence.md)
- [自主化工程路线图](../AUTONOMY_ROADMAP.md)
- [项目 Releases](https://github.com/wxcvm/ARTEX/releases)
- [项目 Issues](https://github.com/wxcvm/ARTEX/issues)

文档面向实际操作，不宣称所有外部模型、MCP、插件或部署环境都已经在该 fork 上完成兼容验证。如遇行为与文档不一致，优先通过当前 `main` 代码、浏览器页面和运行日志定位原因。


### 漏洞详情 Markdown 导出

漏洞详情页顶部提供「导出 Markdown」，可打包漏洞基本信息、摘要、PoC、详细报告、HTTP 流量证据预览、复测记录与溯源节点。若关联记录读取失败，会提示错误而不生成缺漏报告；二进制或截断正文会明确标注。


### Next.js 构建标识可复现

前端的 Next.js `BUILD_ID` 默认为固定 `artex`，可通过环境变量 `ARTEX_BUILD_ID` 指定（例如发布时设置为 Git commit SHA），使相同源码和相同构建标识的重复构建减少差异。此项只控制 Next.js 的构建标识，不保证所有 ZIP 或安装器字节完全一致，也不改变 Agent 的逻辑。
