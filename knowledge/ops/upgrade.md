<!-- 来源：ARTEX 0.4.2 发行包（Telegram 分发）根目录文档，原文逐字保留，仅转码为 UTF-8 并统一 LF。
     适用版本：0.4.2 发行包；概览见 knowledge/ops/INDEX.md；代码级说明见 docs/USER_GUIDE_ZH.md。 -->

# ARTEX 升级指南

已有的 ARTEX 怎么升到新版本，数据一个都不能丢。

---

# 升级前必读

## 三个东西绝对不能动

```
config.json     数据库连接配置（含密码）—— 覆盖了就连不上自己的库
pgdata\         你的数据库本体 —— 删了数据全没
jwt.key         登录密钥 —— 换了所有人被踢下线
```

**为什么 `config.json` 这么要命**

新包里的 `config.json` 是**打包那台机器**的，数据库密码跟你的不一样。

整包覆盖会把它一起换掉。结果就是：

```
程序启动正常 · 能打开网页 · 任务、漏洞、资产全空了
```

数据其实还在 `pgdata\` 里，只是密码不对进不去。**看起来像数据丢了，实际只是连不上。**

这一条是覆盖升级最容易踩的坑。

---

# Windows 升级

## 方式一：覆盖升级（推荐，最省事）

**需要的文件**：新版完整包 `ARTEX-0.4.1-Windows-x64-full.zip`（1.2 GB）

**步骤**：

```
1. 停掉正在跑的 ARTEX
   （双击 stop.bat，或任务管理器结束 artex.exe）

2. 把新版完整包 zip 和 覆盖升级.bat 放到同一个文件夹

3. 双击 覆盖升级.bat

4. 按提示选择【你现在正在用的那个 ARTEX 目录】
   ⚠ 一定要选对——选错了会覆盖掉别的东西，脚本会检查 config.json 和 pgdata 来防呆

5. 等它跑完（几分钟）

6. 双击 start.bat 启动
```

**它会做什么**：

```
会覆盖    程序、技能库、工具链、默认配置
会保留    数据库（pgdata）、config.json、jwt.key、授权文件
          对话记录（transcripts）、任务归档（task-archives）、日志
```

**它怎么保 config.json**：

```
备份 config.json → 解压覆盖 → 立刻还原 config.json
```

中间如果有任何一步失败，会回滚。

---

## 方式二：增量更新（有 update 包时用）

**需要的文件**：`update-xxx.zip`（几百 MB，比完整包小）

**步骤**：

```
1. 停掉 ARTEX

2. 把 update-xxx.zip 放到 ARTEX 根目录（和 应用更新.bat 同一层）

3. 双击 应用更新.bat

4. 它会：校验 sha256 → 备份现有文件 → 覆盖 → 按清单删除废弃文件
```

**和覆盖升级的区别**：

| | 覆盖升级 | 应用更新 |
| --- | --- | --- |
| 需要的包 | 完整包 1.2 GB | 增量包，小 |
| 会不会删旧文件 | 不会 | 会（按清单） |
| 适合 | 跨度大的升级 | 小版本迭代 |

---

# Linux 升级

**没有一键脚本**，需要手动操作。但步骤也不复杂，关键是那三个文件。

## 步骤

```bash
# 0. 先备份（强烈建议）
cd /opt
cp -a artex artex.bak-$(date +%Y%m%d)
# 至少确认 config.json 和 pgdata 有备份

# 1. 停服务
systemctl stop artex
# 或者如果你用 nohup 跑的：
# pkill -f '/opt/artex/artex'

# 2. 解压新包到临时目录
cd /tmp
tar -xzf ARTEX-0.4.1-Linux-amd64.tar.gz
cd ARTEX-0.4.1-Linux-amd64

# 3. 把 config.json 换成你的
cp /opt/artex/config.json .
# 注意：有的版本 config 在 app/config.json，看你的目录结构

# 4. 把新文件覆盖到旧目录，但排除数据和配置
rsync -a --exclude 'pgdata' --exclude 'config.json' \
      --exclude 'jwt.key' --exclude 'license.dat' \
      --exclude 'data' \
      ./ /opt/artex/

# 5. 恢复权限
cd /opt/artex
chmod +x artex
find data/tools -type f -exec chmod +x {} \;

# 6. 启动
systemctl start artex
```

## 关键点

**`rsync` 的 exclude 列表是核心**——它就是"保住那三个东西"的实现。

`data` 目录也排除了：里面是任务记录、流量、对话，新包里没有这些，覆盖了只是浪费时间。

**如果你没有 rsync**：

```bash
# 手动备份关键的
cp /opt/artex/config.json /tmp/config.json.keep
cp /opt/artex/jwt.key /tmp/jwt.key.keep 2>/dev/null

# 解压覆盖
tar -xzf ARTEX-0.4.1-Linux-amd64.tar.gz -C /opt/artex --strip-components=1

# 还原
cp /tmp/config.json.keep /opt/artex/config.json
cp /tmp/jwt.key.keep /opt/artex/jwt.key 2>/dev/null
```

**`pgdata` 不会被覆盖**——新包里根本没有这个目录。但保险起见还是备份一下。

---

# macOS 升级

和 Linux 一样，把路径换掉即可：

```bash
cd ~/artex          # 或你的安装位置
cp -a . ../artex.bak-$(date +%Y%m%d)

# 停掉（如果是终端里跑的，Ctrl+C）
pkill -f 'artex'

# 解压新包
tar -xzf ARTEX-0.4.1-macOS-arm64.tar.gz
cd ARTEX-0.4.1-macOS-arm64

# 覆盖（排除数据和配置）
rsync -a --exclude 'pgdata' --exclude 'config.json' \
      --exclude 'jwt.key' --exclude 'license.dat' --exclude 'data' \
      ./ 原目录/

# 启动
./start.sh
```

**注意 macOS 的架构**：Apple Silicon（M 系列）下 `arm64`，Intel 下 `x86_64`。下错了报 `Bad CPU type in executable`。

---

# 升级后验证

## 1. 看版本号

```bash
# Linux / macOS
curl -s http://127.0.0.1:8787/api/health -H 'Cookie: 你的cookie'

# 或者在日志里看
grep 'listening' data/artex-err.log | tail -1
```

应该看到：

```
ARTEX 0.4.1 backend listening on 0.0.0.0:8787
```

## 2. 数据还在不在

打开控制台，看四个数：

```
任务数     应该和你升级前一样
漏洞数     一样
资产数     一样
历史对话   能翻到
```

**如果全是 0** → config.json 被覆盖了，看下面「出问题怎么办」。

## 3. 新功能在不在

0.4.1 应该有：

```
侧栏状态徽章   任务/C2/代理池 那几行有数字或圆点
系统资源卡片   仪表盘上有 CPU/内存/磁盘
MCP 页面       多一条 fofa 记录
```

---

# 出问题怎么办

## 症状：启动正常，但任务/漏洞/资产全空了

**原因**：`config.json` 被新包的覆盖了，程序连的是打包机器那个库（连不上或连到空库）。

**修法**：

```bash
# 1. 停服务
systemctl stop artex    # 或 stop.bat

# 2. 找回备份的 config.json
ls -la /opt/artex/config.json.bak*    # 覆盖升级脚本会自动备份
# 或者从你自己之前的备份里拿

# 3. 对比一下密码对不对
cat /opt/artex/config.json | grep -i pass

# 4. 换回去，重启
cp config.json.bak config.json
systemctl start artex
```

**Windows 上**：覆盖升级脚本会在 `app\config.json.bak-时间戳` 留备份。

## 症状：启动报端口被占用

```bash
ss -lntp | grep 8787
```

旧进程没停干净。杀掉重来：

```bash
pkill -f '/opt/artex/artex'
```

Windows 上 `stop.bat` 只认 `artex.exe`，如果你改过名字（比如 `artex-final.exe`），得任务管理器手动结束。

## 症状：`GLIBC_2.xx not found`

系统太老（CentOS 7 及更早）。换 Debian 11+ / Ubuntu 20.04+，或者在 Docker 里跑。

## 症状：工具跑不起来，报 Permission denied

压缩包经过 Windows 打包会丢执行权限位：

```bash
cd /opt/artex
chmod +x artex
find data/tools -type f -exec chmod +x {} \;
find . -name '*.sh' -exec chmod +x {} \;
```

## 完全回滚

```bash
# 停服务
systemctl stop artex

# 换回备份目录
mv /opt/artex /opt/artex.failed
mv /opt/artex.bak-20261101 /opt/artex

# 启动
systemctl start artex
```

**所以第 0 步那个 `cp -a` 备份很重要**，别省。

---

# 一页速查

```
Windows
  双击 覆盖升级.bat → 选旧目录 → 等 → 启动
  它会自动保住 config.json 和 pgdata

Linux / macOS
  1. cp -a 备份整个目录
  2. 停服务
  3. 解压新包
  4. rsync 覆盖，排除 pgdata / config.json / jwt.key / data
  5. chmod +x
  6. 启动

验证
  版本号 → 数据还在不在 → 新功能在不在

万一数据"没了"
  99% 是 config.json 被覆盖，找回备份换回去就行，数据一直在 pgdata 里
```
