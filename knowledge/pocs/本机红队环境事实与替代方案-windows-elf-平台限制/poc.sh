# 本机（Windows）红队作业环境事实与替代方案

靶标作业前必读，避免反复踩坑。实测时间：2026-09-28。

## 1. 平台限制：toolkit 里绝大多数是 Linux ELF，Windows 上跑不了
- toolkit 路径：`<HOME>\workspace\redteam\toolkit\`
- 实测：该目录下 **`.exe` 只有 3 个**：`fscan.exe`(8774144B)、`gogo.exe`(4307456B)、`suo5-windows-amd64.exe`(11940352B)。
- `<HOME>\.local\bin\pd-httpx` 文件头 `7F 45 4C 46`（ELF）→ Windows 无法执行（会报 "无法在管道中间运行文档"）。
- `toolkit` 下 `dnsx` / `naabu` / `katana` / `httpx` / `subfinder` 等在本机**并未实际落地为可执行文件**（`Test-Path` 为假）。
- 本机 WSL 无发行版，无法用 Linux 侧工具。
- **验平台的方法**：`Format-Hex <file> | Select-Object -First 1`，前 4 字节 `7F 45 4C 46` = ELF = 不能用；`4D 5A` = PE = 可用。

## 2. Windows 上确定可用的工具
- `nuclei`：`<HOME>\go\bin\nuclei.exe`，模板库 `<HOME>\nuclei-templates`（实测 **13742 个 yaml**）。
- `fscan.exe` / `gogo.exe` / `suo5-windows-amd64.exe`（toolkit 内）。
- `curl`（Windows 自带 curl.exe，本机可用）、`python`（`D:\Python313\python.exe`）、`nmap`（`C:\Program Files (x86)\Nmap`，但见第 4 条）。
- 缺失/不可用：`jq`、`whatweb`、`masscan`（自研 Python/curl 替代）。

## 3. FOFA_KEY 位置
配置在 `<HOME>\workspace\.env`，键名 `FOFA_KEY`（同文件还有 `REDTEAM_VPS_HOST` / `REDTEAM_VPS_USER` / `REDTEAM_VPS_KEY`）。
`redteam_preflight` 能读到（报告 configured），但**子进程环境变量里不一定有** → 需要时直接读该文件取值，或改用可用的 MCP `fofa` 工具。

## 4. nmap 对境外目标的结果不可信
对两台不同云厂商的机器（203.0.113.29 AWS-SG、203.0.113.105 Hostinger）扫描，均只返回 `25/110/143 open` + 997 filtered，与 FOFA 境外观测（80/443/6379/3306…）完全矛盾，且两台分布一模一样。
⇒ 本机出口链路对境外高频 connect 扫描存在干扰/过滤。**端口面以被动测绘（FOFA）+ 单次 HTTP 实测为准，不要用 nmap 结果落库。**

## 5. 替代工具选型
- Web 存活/指纹/标题：自研限速 Python 探针（`runs/recon-tgt6403/http-probe.py`、`runs/e5-probe.py`），串行 + 固定间隔 + 超时即停，比 pd-httpx 更适合“少而准”。
- vhost / 源站命中（Cloudflare 后域名找不到源站时）：连源站 IP:443 但用 `Host` 头 + SNI 指定域名（等价 `curl --resolve`），见 `runs/e5-probe.py` 的 `vhost` 模式。
- 目录/接口枚举：Python 自研 + 白名单路径，或 `dirsearch`（若为 .py 则可运行，需验平台）。