<!-- 来源：ARTEX 0.4.2 发行包（Telegram 分发）根目录文档，原文逐字保留，仅转码为 UTF-8 并统一 LF。
     适用版本：0.4.2 发行包；概览见 knowledge/ops/INDEX.md；代码级说明见 docs/USER_GUIDE_ZH.md。 -->

# ARTEX 故障排查手册

跑不起来、跑着跑着不动了、数据看不见了——按症状查。

---

# 一、启动类

## 启动就退出 / 秒退

**先看日志**，别猜：

```bash
# Linux / macOS
tail -50 data/artex-err.log

# Windows
type app\data\artex-err.log
```

**最常见的原因：单实例锁**

ARTEX 基于 PostgreSQL DSN 做单实例保护。同一个库、同一个数据目录只能跑一个实例。

```
[instance] ARTEX 已在运行（同一数据库/数据目录），本次启动退出
```

说明已经有一个在跑了。先找到它：

```bash
# Linux / macOS
ps aux | grep artex | grep -v grep

# Windows
tasklist | findstr artex
```

杀掉再启动。**注意**：`stop.bat` 只认 `artex.exe`，如果你把程序改名成 `artex-final.exe` 之类，`stop.bat` 会以为没在运行，得任务管理器手动结束。

## 卡在初始化数据库

第一次启动要建库建表，**1-3 分钟正常**。

超过 5 分钟没动静：

```bash
# 看 PostgreSQL 起来没
ss -lntp | grep 5432
systemctl status postgresql

# 看 PG 自己的错误日志
# Linux
tail -50 pgdata/log/*.log
# Windows
type app\data\pg_stderr.log
```

**常见原因**：

| 现象 | 原因 | 处理 |
| --- | --- | --- |
| `could not connect to server` | PG 没起 | 手动起 PG，看日志 |
| `password authentication failed` | config.json 密码不对 | 见下面「数据看不见」 |
| `database "artex" does not exist` | 库没建 | 重跑初始化 |
| `could not create shared memory` | 系统 shm 太小 | 调 PG 的 `shared_buffers` |

## `GLIBC_2.xx not found`

```bash
./artex: /lib/x86_64-linux-gnu/libc.so.6: version `GLIBC_2.28' not found
```

系统太老（CentOS 7 / Ubuntu 18.04 及更早）。

**处理**：换 Debian 11+ / Ubuntu 20.04+ / CentOS 8+，或者用 Docker。

ARTEX 用 Go 静态编译，理论上不依赖系统库，但某些版本的工具链会有 glibc 依赖。

## `Permission denied`

压缩包经过 Windows 打包会丢执行权限位：

```bash
cd /opt/artex
chmod +x artex start.sh start-artex.sh *.sh
find data/tools -type f -exec chmod +x {} \;
```

验证：

```bash
find data/tools -type f ! -perm -u+x | head
```

这条命令会列出**没有执行权限**的文件，如果有输出就说明还有漏的。

## 端口 8787 被占用

```bash
# Linux / macOS
lsof -i :8787

# Windows
netstat -ano | findstr :8787
```

杀掉占用的进程，或者换端口：

```bash
./artex -addr 0.0.0.0:9000
```

---

# 二、运行类

## Agent 不动 / 任务卡在"探索中"

按这个顺序查，**别跳步**：

### ① LLM 配好了吗

```
系统 → API KEY 配置 → 点「测试」
```

**没有 LLM key，agent 就是个空壳**——它不会自己思考，只会等你说下一步。

测试失败看报错：余额不足 / key 无效 / 网络不通 / 模型名不对。

### ② 能出网吗

```
执行 → IP 代理池 → 看「出口 IP」
```

如果出口 IP 是空的或者报错：

```bash
# 本地直连测
curl -x http://127.0.0.1:17890 https://api.ipify.org

# 不走代理测
curl https://api.ipify.org
```

**代理池启用后所有流量都走代理**。代理挂了，agent 连 LLM 都访问不了，表现就是"不动"。

### ③ 任务日志里有没有报错

任务详情页看实时日志。常见：

```
403 / 429 / captcha     → 目标封了，看代理池有没有自动换
Connection reset        → 被墙或被 RST，换节点
timeout                 → 网络慢或目标不可达
SSL_ERROR               → TLS 握手失败，可能被中间设备干扰
```

### ④ 意图是不是卡住了

系统有意图看门狗：某个意图领了太久没进展会报警。日志里搜：

```bash
grep -i 'watchdog\|意图.*超时\|intent' data/artex-err.log | tail -20
```

## 任务跑一半停了

看任务状态：`running` / `paused` / `done`。

**paused 是暂停**，不是失败。点「继续」会从断点接着跑。

**真的失败了**，看这几处：

```bash
# 最近的任务相关错误
grep -iE 'error|fail|panic' data/artex-err.log | tail -30
```

**panic 比较严重**——Go 的崩溃。把那次 panic 的完整堆栈发出来。

## 代理池异常

### 出口 IP 一直是空的

```
执行 → IP 代理池 → 看节点数和「出口 IP」
```

| 现象 | 原因 |
| --- | --- |
| 节点数 0 | 订阅没拉到 / 网关配置错 |
| 有节点但出口空 | mihomo 没起来，或前置代理不通 |
| 报 `forbidden ip` | 住宅网关拒绝你的出口 IP（大陆直连），要填前置代理 |

### mihomo 崩了

有进程守护，崩了 30 秒内自动拉起。日志：

```bash
grep -i mihomo data/artex-err.log | tail -20
```

看到连续重启（退避 5→15→30→60 秒，5 次后放弃），说明配置有问题，不是偶发崩溃：

```bash
# 看 mihomo 自己的日志
tail -50 data/proxy-pool/mihomo.log
```

### 自动换 IP 被关了

日志里找：

```
节点池疑似耗尽：3 次切换（3 节点）仍被封，已暂停自动切换
```

**这是保护机制，不是 bug**。意思是轮完一圈还是被封，继续换只是空转。

处理：换国家/地区，或者加节点数，然后手动开启自动切换。

## 数据看不见了

**症状**：能打开页面，但任务、漏洞、资产全是 0。

**99% 是 `config.json` 被覆盖了**——程序连的是另一个库（空的或者连不上）。

```bash
# 1. 看 config.json 里的数据库密码
grep -i pass config.json

# 2. 手动连一下试试
PGPASSWORD='你看到的密码' psql -h 127.0.0.1 -p 15433 -U artex -d artex -c '\dt'

# 3. 如果连不上，说明密码不对
# 找回备份
ls -la config.json*
# 或者从 pgdata 里改密码
```

**数据一直在 `pgdata/` 里，不会因为 config 被换而丢失。**

详见 `ARTEX升级指南.md` 里「启动正常但数据全空」那节。

---

# 三、性能类

## CPU 一直很高

```bash
top -p $(pgrep -f artex)
```

**正常情况**：跑任务时 30-80% 是正常的，agent 在并发干活。

**异常情况**：

| 现象 | 可能原因 |
| --- | --- |
| 空闲也 100% | 死循环 / 某个 worker 卡住 |
| 突然飙高不降 | 扫描任务开的并发太大 |
| 伴随内存暴涨 | 内存泄漏 |

看是哪个 goroutine 在烧：

```bash
# 如果有 pprof 端口（默认没开）
go tool pprof http://127.0.0.1:8787/debug/pprof/profile
```

更实际的办法：看任务列表，是不是某个任务开了巨量并发。

## 内存一直涨

```bash
ps -o pid,rss,cmd -p $(pgrep -f artex)
```

看 RSS。几百 MB 正常，涨到几个 G 就要查。

**常见原因**：

- 大文件读进内存（某个工具输出被全量捕获）
- 任务产物没落盘（都在内存里攒着）
- traffic 捕获开着但没清理

**处理**：

```bash
# 关掉不需要的流量捕获
# 系统设置里关 capture

# 清理任务产物
du -sh data/tasks/* | sort -rh | head
```

## 页面很卡 / 接口超时

```bash
# 测接口响应
time curl -s http://127.0.0.1:8787/api/health
```

如果某个接口特别慢，可能是：

- 数据库慢查询（任务/漏洞表大了）
- 出口探测在阻塞（`/api/proxy-pool` 曾经有过这个坑，现在改成异步了）
- mihomo 拉不起来，请求在等超时

**数据库维护**：

```bash
psql -h 127.0.0.1 -p 15433 -U artex -d artex -c 'VACUUM ANALYZE;'
```

## 磁盘满了

```bash
df -h /
du -sh data/* | sort -rh
```

**占空间的大头**：

| 目录 | 是什么 | 能删吗 |
| --- | --- | --- |
| `data/tasks/` | 任务产物（每个几百 MB） | 删不要的任务目录 |
| `data/traffic-tasks/` | 流量捕获 | 能 |
| `data/artex-err.log` | 运行日志（能涨到几十 MB） | 能轮转 |
| `pgdata/` | 数据库 | **不能删** |
| `data/tools/` | 工具链 | **不能删** |

**日志轮转**：

```bash
# 保留最后 5000 行
tail -5000 data/artex-err.log > /tmp/keep.log
cat /tmp/keep.log > data/artex-err.log
rm /tmp/keep.log
```

**/tmp 被写满**（服务器上常见）：

```bash
du -sh /tmp/* | sort -rh | head -10
# 清 24 小时前的目录（保留活跃的）
find /tmp -maxdepth 1 -mindepth 1 -type d -mtime +0 \
  ! -name 'systemd-private-*' ! -name '.X11-unix' -exec rm -rf {} +
```

> **注意** `find -mtime +1` 是**超过 48 小时**，`-mtime +0` 才是 24 小时以上。这个坑踩过。

---

# 四、网络类

## 完全出不去网

```bash
# 1. 系统能出网吗
curl -s --max-time 5 https://api.ipify.org

# 2. 走本地代理能出吗
curl -s --max-time 5 -x http://127.0.0.1:17890 https://api.ipify.org

# 3. 代理端口在监听吗
ss -lntp | grep 17890
```

三种结果对应三种问题：

| 现象 | 原因 |
| --- | --- |
| 1 通 2 不通 | 代理池配置有问题 |
| 1、2 都不通 | 系统网络问题（防火墙 / DNS / 路由） |
| 2 通但 3 没监听 | mihomo 没起来 |

## 被墙 / 连接重置

```
Connection reset by peer
SSL_ERROR_SYSCALL
```

这三类信号会触发**自动换 IP**：

```
目标封禁（403/验证码）  → 立刻切
连接被墙（reset/timeout）→ 连续 2 次才切
本地代理挂了            → 不切，报警
```

如果一直切不动，看代理池的节点数和「最近错误」。

## 前置代理不通

大陆机器连住宅网关必须走前置。诊断：

```
执行 → IP 代理池 → 看「最近错误」
```

会分两级提示：

```
前置代理 127.0.0.1:7892 不可达 —— 确认那个客户端在运行
前置代理可达，但网关连接失败 —— 检查余额/白名单/地区限制
```

---

# 五、诊断工具箱

## 一键收集环境信息

```bash
#!/bin/bash
echo "=== 系统 ==="
uname -a
cat /etc/os-release | head -3

echo "=== 版本 ==="
grep -o 'ARTEX [0-9.]*' data/artex-err.log | tail -1

echo "=== 进程 ==="
ps -o pid,etime,rss,%cpu,cmd -p $(pgrep -f 'artex -addr' | head -1) 2>/dev/null

echo "=== 端口 ==="
ss -lntp | grep -E '8787|17890|17891|5432'

echo "=== 磁盘 ==="
df -h /

echo "=== 最近错误 ==="
tail -100 data/artex-err.log | grep -iE 'error|panic|fail' | tail -20

echo "=== 数据库 ==="
psql -h 127.0.0.1 -p 15433 -U artex -d artex -t -c \
  "select count(*) from tasks" 2>&1 | head -2

echo "=== 任务状态 ==="
psql -h 127.0.0.1 -p 15433 -U artex -d artex -t -c \
  "select status, count(*) from tasks group by status" 2>&1
```

把输出发出来，基本能定位 80% 的问题。

## 关键日志位置

```
data/artex-err.log          主日志（错误 + 重要事件）
data/artex-out.log          标准输出（启动横幅、版本信息）
data/pg_stderr.log          PostgreSQL 日志（Windows）
pgdata/log/*.log            PostgreSQL 日志（Linux）
data/proxy-pool/mihomo.log  代理内核日志
data/oom-watchdog.log       内存守护
```

## 常用过滤

```bash
# 只看错误
grep -iE 'error|panic|fatal' data/artex-err.log | tail -50

# 代理相关
grep -iE 'proxy-pool|mihomo|proxy' data/artex-err.log | tail -50

# 蜜罐命中
grep -i 'honeypot' data/artex-err.log | tail -20

# 门控
grep -i 'gate' data/artex-err.log | tail -10

# 某个任务
grep 'task.*#123' data/artex-err.log | tail -50
```

---

# 六、还是解决不了

收集这些发出来：

```
1. 上面那个「一键收集」脚本的输出
2. 完整的问题现象（什么时候开始的、做了什么操作）
3. data/artex-err.log 的最后 200 行
4. 如果是崩溃：panic 的完整堆栈
```

**不要只发一句"用不了"**——没有日志没法查。
