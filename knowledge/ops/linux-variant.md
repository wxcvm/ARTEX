<!-- 来源：ARTEX 0.4.2 发行包（Telegram 分发）根目录文档，原文逐字保留，仅转码为 UTF-8 并统一 LF。
     适用版本：0.4.2 发行包；概览见 knowledge/ops/INDEX.md；代码级说明见 docs/USER_GUIDE_ZH.md。 -->

# Linux 该下哪个包？

ARTEX 的 Linux 版分两个架构：

```
ARTEX-0.4.1-Linux-amd64.tar.gz    530 MB
ARTEX-0.4.1-Linux-arm64.tar.gz    525 MB
```

**下错了装不上**，所以先看这一页再下载。

---

## 三十秒判断

在你要装的那台 Linux 上执行一条命令：

```bash
uname -m
```

输出对照：

| 输出 | 下这个包 |
| --- | --- |
| `x86_64` | **amd64** |
| `amd64` | **amd64** |
| `aarch64` | **arm64** |
| `arm64` | **arm64** |
| `armv7l` / `armv6l` | 不支持（32 位 ARM） |

**`x86_64` 和 `amd64` 是同一个东西**，只是叫法不同（Linux 用前者，打包用后者）。

**`aarch64` 和 `arm64` 也是同一个东西**。

---

## 两个包差在哪

不是"完整版"和"精简版"的关系，是**给不同 CPU 指令集编译的**。

| | amd64 | arm64 |
| --- | --- | --- |
| CPU 指令集 | x86-64 | AArch64 |
| 主程序 | 41 MB，Linux x86-64 编译 | 39 MB，Linux ARM64 编译 |
| 工具链 | 34 件，全部 amd64 版 | 34 件，全部 arm64 版 |
| mihomo | `mihomo-linux-amd64` | `mihomo-linux-arm64` |
| 功能 | **完全一样** | **完全一样** |

包里所有二进制都按架构对应，没有一个是"凑数"的。两个包解压出来的目录结构、配置文件、文档完全相同，只有二进制不同。

---

## 常见场景对照

### 选 amd64 的

```
阿里云 ECS 通用型 / 计算型
腾讯云 CVM 标准型
AWS EC2 t3 / m5 / c5 系列
Vultr / DigitalOcean / Linode 普通 VPS
自建服务器（Intel / AMD CPU）
VMware / VirtualBox 虚拟机
WSL2
几乎所有国内 VPS
```

**一句话**：只要不是明确标注 ARM 的，基本都是 amd64。市面上 95% 以上的 Linux 服务器是 x86。

### 选 arm64 的

```
AWS EC2 Graviton 系列（t4g / m6g / c7g）
阿里云 倚天 710 实例
腾讯云 ARM 标准型
Oracle Cloud Ampere A1（免费套餐就是 ARM）
树莓派 4 / 5（64 位系统）
Apple Silicon 上的 Linux 虚拟机（UTM / Parallels / Lima）
部分 NAS（群晖 ARM 机型）
```

**判断依据**：云厂商控制台里实例规格带 `g`（Graviton）、`ARM`、`Ampere`、`倚天` 字样，或者价格明显比同配置 x86 便宜（ARM 实例通常便宜 20-40%）。

---

## 拿不准时的两条路

### 路线一：直接问系统

```bash
# 看内核架构
uname -m

# 看 CPU 型号（x86 会显示 Intel/AMD，ARM 会显示 aarch64 相关）
lscpu | grep -E "Architecture|Model name|Vendor"
```

`Architecture: aarch64` → arm64
`Architecture: x86_64` → amd64

### 路线二：两个都下，试装

包不大（500 多 MB），下错了删掉重下就行。装错的报错长这样：

```bash
./artex: cannot execute binary file: Exec format error
```

或者：

```bash
bash: ./artex: No such file or directory
```

**看到这两个报错，八成是架构下错了**（第二个报错很反直觉——文件明明在，但内核认不出这个格式，所以报"找不到"）。

---

## 包内含什么 / 不含什么

### 完整包含（两个架构都有）

```
侦察    subfinder  dnsx  naabu  httpx  katana  gospider  ksubdomain
Web     ffuf  nuclei(14,365 模板)  dirsearch  oneforall
隧道    chisel  frp(frpc+frps)  suo5  ligolo  penelope
提权    linpeas  peass
其他    mihomo  mcp  wordlists  python 工具脚本
知识库  pocs 392 条 / skills 38 个 / knowledge.db + cve.db
```

### Windows 版有、Linux/macOS 版**没有**

| 工具 | 作用 | 为什么没有 |
| --- | --- | --- |
| **fscan** | 内网综合扫描 | 只有 Windows 编译版随包 |
| **gogo** | 内网自动化引擎 | 同上 |
| **nxc**（netexec） | SMB/WinRM/LDAP 横向 | 同上 |
| **sqlmap** | SQL 注入 | 同上 |
| AntSword / Behinder_v4.1 | WebShell 图形管理 | GUI 工具，Linux 服务器用不上 |
| jdk11 | Java 运行时（跑上两个） | 同上 |

**这五个是 Windows 包独有，Linux/macOS 版不带。**

### 需要的话自己装

```bash
# fscan（Go 写的，下对应架构的二进制）
wget https://github.com/shadow1ng/fscan/releases/latest/download/fscan_linux_amd64
chmod +x fscan_linux_amd64 && mv fscan_linux_amd64 /opt/artex/data/tools/fscan

# gogo（chainreactors 出品）
# 从 https://github.com/chainreactors/gogo/releases 下 linux_amd64 版本

# sqlmap（Python，跨平台）
apt install sqlmap
# 或
pip3 install sqlmap

# netexec（nxc，Python）
pip3 install netexec
# 或
pipx install netexec
```

装好后放到 `data/tools/` 下，ARTEX 启动时会自动扫描并注册成可用工具。

**注意**：文件名不能带架构后缀（`fscan_linux_amd64` → 改成 `fscan`），平台的工具注册逻辑按名字匹配。

### 不影响的部分

- **内网渗透**：没有 fscan/gogo，但可以用 `nmap` + `naabu` + `nuclei` 组合，或者走 `fscan-intranet` skill 的降级方案
- **WebShell 管理**：有 GSL5 哥斯拉 MCP，不依赖 AntSword/Behinder 的 GUI
- **SQL 注入**：`web-sqli` skill 里有手工 payload，不是必须用 sqlmap

---

## 装法

两个架构的安装步骤**完全一样**：

```bash
# 1. 解压
tar -xzf ARTEX-0.4.1-Linux-amd64.tar.gz
cd ARTEX-0.4.1-Linux-amd64

# 2. 给执行权限（压缩包有时会丢权限位）
chmod +x artex
find data/tools -type f -exec chmod +x {} \;

# 3. 启动
./artex

# 或者后台跑
nohup ./artex > artex.log 2>&1 &
```

默认监听 `0.0.0.0:8787`。启动后日志里会打印入口路径：

```
[gate] 伪装门控已启用，入口路径: /g-xxxxxxxxxxxx
```

浏览器访问 `http://你的IP:8787/g-xxxxxxxxxxxx` 进控制台。

> **注意**：直接访问 `http://IP:8787/` 会看到一个 nginx 默认页——那是**伪装门控**，故意的。别人扫到你的端口，看到的也是一样的页面。必须带 `/g-xxx` 路径才能进。

---

## 常见问题

**Q：装完启动报 `GLIBC_2.xx not found`**

系统太老。ARTEX 用 Go 静态编译，正常情况下不依赖系统库，但如果报这个，说明系统 glibc 版本低于 2.28（CentOS 7 及更早）。换 Debian 11+ / Ubuntu 20.04+ / CentOS 8+，或者在 Docker 里跑。

**Q：ARM 机器上某些工具跑不起来**

先确认工具权限：

```bash
ls -la data/tools/ | head
find data/tools -type f ! -perm -u+x
```

压缩包经过 Windows 打包时权限位可能丢失，第二条命令会列出没有执行权限的文件，统一 `chmod +x` 即可。

**Q：能不能跨架构跑？**

可以，但没必要。用 `qemu-user-static` 能模拟，但性能掉一半以上，而且部分工具（尤其是需要特定指令集的）会崩。直接下对应架构的包。

**Q：以后换机器了怎么办？**

数据全在 `data/` 目录和 PostgreSQL 里。换架构只需要：
1. 新机器上下对应架构的包
2. 把旧机器的 `data/` 目录拷过去
3. 导出导入数据库

二进制不能跨架构复制，但**数据可以**。

---

## 一句话总结

```bash
uname -m
```

- 看到 `x86_64` → **amd64**
- 看到 `aarch64` → **arm64**

其他都是小概率情况，真遇到了按上面"拿不准"那节处理。
