## 一、本机（Windows）链路事实（实测）
1. **系统级代理/TUN 劫持**：curl/Python 直连境外目标全部 timeout；`curl -v -x http://<ip>:<port>` 连 "Connected to" 都不打印；`scp` 把 IPv6 目标 `2602:...` 误解析成内网地址 `0.0.10.42`（本机 TUN 网关）。⇒ **本机所有直接出网探测的结果不可信（含 nmap 的 "25/110/143 open + 997 filtered"、fscan 的大面积漏报）**。
2. 判据：同一批端口，本机 timeout、VPS 拿到真实 banner；或任何解析结果落在 `0.0.10.0/24`。
3. **不适用**：不要用本机直连结果做「端口关闭」结论（负向结论不可信）；正向拿到 banner 的结论仍可信（如本机 fscan 拿到 squid banner）。

## 一·补 fscan.exe 的实测边界（重要）
本机跑 `fscan.exe -h <ip> -p <ports> -np -nobr -nopoc -noredis -t 2 -time 5`：
- **正向可信**：成功取回 `203.0.113.29:5566` 的 `HTTP/1.1 400 Bad Request / Server: squid`，与 `203.0.113.97:6379` 的 `-NOAUTH Authentication required.`
- **负向不可信**：同一台 `203.0.113.29` 的 80/443/9527（FOFA 明确存在）**一条都没报**；`203.0.113.28` 直接报「0 个开放端口」。
- 用法约束：**必加 `-nobr`（禁爆破）、`-nopoc`、`-noredis`、`-t 2`**；`-time 5` 会被自适应成 Timeout=5000ms。
⇒ 结论：**fscan 报出来的可以信，它没报的不能当成不存在**；要下「端口关闭」结论必须换出口复核。

## 一·补 双出口对照法（判「目标侧不可达」的标准动作）
出现「本机不可达」时不要直接判死，按顺序做两条对照，两条都不通才写「目标侧不可达」：
1. **本机原生 socket**（`socket.create_connection((ip,port), timeout=10)`；**不要用 curl**，它会读 IE/系统代理）：不通只说明本机链路有问题。
2. **VPS 出口 socket**（`timeout 6 bash -c 'cat </dev/null >/dev/tcp/IP/PORT'`）：仍不通 ⇒ 换第三条。
3. **第三方视角**（FOFA 等被动数据 + 目标其它端口行为）：若同机其它端口也不可达、且被动数据陈旧 ⇒ 判「目标侧不可达」并**止损**，不再耗轮次。

## 二、可信跳板（首选）
演练 VPS：`REDTEAM_VPS_HOST=2602:fa59:2:a9b9::1`（IPv6）、`REDTEAM_VPS_USER=root`、`REDTEAM_VPS_KEY=<HOME>\workspace\redteam\toolkit\vps\id_rsa`，出口 IPv4 `203.0.113.35`（干净）。Ubuntu 6.8 / 1C / 961MB；自带 curl、wget、python3、nc、jq（**无 nmap/socat**）。
- 连接写法（PowerShell 下必须给 IPv6 加方括号，否则 scp 解析错）：`scp -i $k file "root@[2602:fa59:2:a9b9::1]:/tmp/f"`；`ssh -i $k -o IdentitiesOnly=yes -o StrictHostKeyChecking=no -o UserKnownHostsFile=NUL root@2602:fa59:2:a9b9::1 "<cmd>"`
- upload 后 `sed -i 's/\r$//' /tmp/f`（Windows 换行）。
- 用完清理：`rm -f /tmp/<你的脚本>`，不留持久化。
- 提示：VPS 上已有其它角色留下的文件（`/tmp/api_scan.json`、`b.*` 等），**不要删**。

## 三、Windows 本机可用工具白名单
可用：`<HOME>\workspace\redteam\toolkit\fscan.exe`(v2.2.1)、`gogo.exe`、`suo5-windows-amd64.exe`、`java`(`jdk11\`)、`godzilla.jar`、`nuclei` + 模板 `<HOME>\nuclei-templates`(13,764 条 / CVE 4,480)、`curl.exe`、`python`(3.13)。
**Linux ELF 二进制在本机不可用**（`dnsx/naabu/katana/httpx/subfinder/ksubdomain/frpc/frps/fscan(无后缀)/gogo(无后缀)/suo5-linux-amd64` 全是 ELF，跑不了）——**要跑就把它们传到 VPS 上跑**，不要在 Windows 上反复试。
`FOFA_KEY` 在 `<HOME>\workspace/.env`（也可用 MCP `fofa` 工具）。

## 四、PowerShell 执行远端命令的坑
- 远端命令含 `$var`、单引号、括号时，**写成脚本文件 scp 上去执行**，不要在 ssh 的 inline 字符串里拼（实测多次被 PowerShell 提前展开/截断，导致 exit 1 且无输出）。
- `curl` 在 PowerShell 里是 `Invoke-WebRequest` 别名，**必须写 `curl.exe`**。
- ssh/scp 的 `Warning: Permanently added ...` 会被 PowerShell 当 stderr 报 RemoteException，**不代表命令失败**（用 `2>$null` 或看实际输出判断）。