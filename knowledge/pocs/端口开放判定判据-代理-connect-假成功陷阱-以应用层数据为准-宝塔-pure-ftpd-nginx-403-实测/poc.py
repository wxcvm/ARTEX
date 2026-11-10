## 一、核心陷阱：HTTP 代理的 CONNECT "假成功"
很多 HTTP 代理（含常见机场客户端）对**任意** `CONNECT host:port` 都回 `HTTP/1.1 200 Connection established`（SOCKS5 模式则回 `request granted`），但**实际并不转发**，表现为"连上之后零字节"。

⇒ **绝对不要用"CONNECT 成功"判定端口开放**。唯一可靠的判据是：**目标端口有没有主动吐出应用层数据**。

| 端口 | 打开时必然出现的应用层数据（无需我们发任何请求） |
|---|---|
| 21 FTP | `220 ... ` Pure-FTPd/vsftpd banner（宝塔 = `220---------- Welcome to Pure-FTPd [privsep] [TLS] ----------`） |
| 22 SSH | `SSH-2.0-OpenSSH_x.y` banner |
| 3306 MySQL | MySQL greeting（首字节是可读协议版本 + 版本号字符串） |
| 8888 宝塔 | 需主动发 HTTP 请求，看是否回宝塔登录页 |
| 888 phpMyAdmin | 需主动发 HTTP 请求（宝塔默认配置对外返回 nginx `403 Forbidden`） |

## 二、可直接复用的探针（Python，经 HTTP 代理）
```python
import socket
def probe(host, port, http=False, proxy=("127.0.0.1",7892), timeout=12):
    s = socket.create_connection(proxy, timeout=timeout)
    s.sendall(("CONNECT %s:%d HTTP/1.1\r\nHost: %s:%d\r\nProxy-Connection: Keep-Alive\r\n\r\n"
               % (host, port, host, port)).encode())
    buf = b""
    while b"\r\n\r\n" not in buf:
        d = s.recv(1024)
        if not d: break
        buf += d
    head, _, rest = buf.partition(b"\r\n\r\n")
    if b"200" not in head.split(b"\r\n")[0]: return ("REFUSED", head[:80])
    data = rest
    if http:
        s.sendall(("GET / HTTP/1.1\r\nHost: %s:%d\r\nUser-Agent: Mozilla/5.0\r\nConnection: close\r\n\r\n"
                   % (host, port)).encode())
    s.settimeout(6 if not http else timeout)
    try:
        while len(data) < 4000:
            d = s.recv(1024)
            if not d: break
            data += d
    except socket.timeout:
        pass
    return ("OPEN" if data else "NO-DATA(可能代理假成功/目标过滤)", data[:300])
```
- 判定：`OPEN` + 能读出 banner/HTTP 头 = **真开**；`NO-DATA` = **不可确认**（不要写成"关闭"）。

## 三、FTP（宝塔 Pure-FTPd）的坑
- 宝塔的 Pure-FTPd 用 **MySQL 虚拟用户**后端，**存在与不存在的用户名都回同一条 `530 Login authentication failed`** ⇒ **无法用响应差异枚举用户名**，只能直接试口令。
- 实测 12 对失败（`target234.example.com`/同名、`tgt2840`、`tgt3090`、`admin`、`wwwtgt2840`、`tgt2840a0`、`example.com`、`target235.example.com`、`ftp`/123456、`www`/123456）；**判据只有 `230` 才算通**。
- 字典正确来源（按依据强度）：宝塔站点**根目录名**（如 `/www/wwwroot/target234.example.com` → 账号常同名）> 域名段 > 面板/系统名。
- FTP 一旦能登且能写 web 目录 ⇒ **直接 PUT `.php` 马 = WebShell**，绕过所有 HTTP 层白名单。

## 四、888（宝塔 phpMyAdmin 端口）的坑
整站 `403 Forbidden`（nginx 的 IP 白名单）。已实测无效的绕过：`/pma`、`/phpmyadmin/`、`/PMA/`、`/pma` 无斜杠、`X-Forwarded-For/X-Real-IP/Client-IP: 127.0.0.1`。
⇒ 这类 403 是**真白名单**（nginx `allow/deny` 不看 XFF，除非配了 real_ip）。**别在这上面继续烧包**。

## 五、换出口复核的正确做法
- **免费代理池不可用**：Squid 类代理默认**禁止 CONNECT 到非 443 端口**（返回 `403 Proxy Error`），其余多为超时；实测 2286 个候选、276 个"验活"，无一能用于非标端口探测。
- 正确做法：**用自有 VPS 做出口**，优先 VPS 直连 `nmap -sT -Pn -p 21,22,80,443,3306,888,8888 <host>`；VPS SSH 不通时用 `ssh -D 1080` 建 SOCKS 或经代理建 CONNECT 隧道（ProxyCommand）。
- 复核的价值：**区分"目标侧防火墙/宝塔 IP 白名单"与"本机出口路径问题"**——同一台机器上 21 与 888 能通、3306/8888 零字节，说明代理本身能转发非标端口，则 3306/8888 更可能是**目标侧只放行特定来源**（宝塔面板默认支持"面板 IP 白名单"）。