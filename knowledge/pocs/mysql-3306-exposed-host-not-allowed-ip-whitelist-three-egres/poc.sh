## 现象

目标主机监听 3306，公网 TCP 可建连，服务端主动回 MySQL 明文 ERR 包：

```
74B:  46 00 00 00 3f 6a 04 "Host '192.0.2.12' is not allowed to connect to this MySQL server"
```

错误文案里直接回显我方出口 IP —— 这是 mysql 的 ERROR 1130（鉴权前 Host ACL 拒绝），
既证明服务存在，又证明通道正常。

## 复现（Windows / curl.exe，Linux 用 curl 同理）

```bash
# 出口 A（境外 socks5）
curl.exe -s --http0.9 --proxy socks5://192.0.2.12:10808 \
  --connect-timeout 30 --max-time 25 http://<TARGET>:3306/ | xxd

# 出口 B（另一条 socks5，证明不是单出口临时封）
curl.exe -s --http0.9 --proxy socks5://192.0.2.13:10808 \
  --connect-timeout 30 --max-time 25 http://<TARGET>:3306/

# 本机直连（第三出口，绕过代理）
curl.exe -s --http0.9 --noproxy "*" \
  --connect-timeout 12 --max-time 18 http://<TARGET>:3306/

# 阳性对照：同一时刻打 Web 端口
curl.exe -s --noproxy "*" -o NUL -w "%{http_code} %{time_total}\n" http://<TARGET>/
```

## 判读

| 响应 | 结论 |
|---|---|
| 明文 ERR 含 `Host '<我方出口IP>' is not allowed to connect` | TCP 通 + 服务是 MySQL + IP ACL 拒绝（本条） |
| 多出口相同、Web 端口 200 且耗时正常 | 白名单基线成立，非网络故障 |
| `Access denied for user 'x'@'y' (using password: ...)` | ACL 已过，转口令攻击（1045） |
| 超时 / 0 字节 | 端口假开或被 DROP，非本条 |
| curl 报 `HTTP/0.9` 相关错误 | 缺 `--http0.9`，不是服务异常 |

## 判定后的处置

- 本阶段放弃 mysql 客户端暴破（结果必然同样是 1130），只登记「3306 暴露」资产。
- 走迂回：找白名单段内的出口（目标同 C 段/内网），或从目标 Web 应用的 SSRF / 命令执行
  回连 `127.0.0.1:3306`、内网 DB 主机。
- 打进去后的第一目标：`mysql` 库 `user` 表（看还有哪些 Host 放行）、业务库配置表
  （如 cdb_settings 一类，内含 authkey/第三方密钥）。

## 踩坑

- 不加 `--http0.9`：curl 认为响应不是 HTTP，直接报错，容易误判成服务异常。
- 只用一个出口就下结论：无法区分「单 IP 被临时封」与「IP 白名单基线」。
- 没做阳性对照：网络抖动/代理坏掉时同样零字节，容易误报成开放端口。
- 经代理时回显的是代理出口 IP；想命中白名单必须用白名单段内的出口。
