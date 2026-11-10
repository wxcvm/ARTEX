<!-- 来源：ARTEX 0.4.2 发行包（Telegram 分发）根目录文档，原文逐字保留，仅转码为 UTF-8 并统一 LF。
     适用版本：0.4.2 发行包；概览见 knowledge/ops/INDEX.md；代码级说明见 docs/USER_GUIDE_ZH.md。 -->

# ARTEX 安全加固

部署在公网 VPS 上必须做的事。**不做会出真事。**

---

# 一、先看三个真实事故

这些不是假设，是实测踩出来的。

## ① SSH 密码登录 = 每天几万次爆破

> 实测一台机器累计 **6.4 万次**暴力破解尝试，并出现过陌生 IP 的 root 会话，**挂了 22 小时没人发现**。

VPS 开了 22 端口，全世界都在扫。密码再复杂也架不住撞库和字典。

**必须改成密钥登录。**

## ② agent 起的靶场服务全在裸奔

> 渗透测试期间 agent 会起大量靶场 / OOB 服务，**实测 70+ 个**，它们默认绑 `0.0.0.0`、无鉴权，**全部直接对公网可达**。

这是最容易被忽略的一条——你在打别人，同时把自己的靶场、OOB 接收端、临时 HTTP 服务全挂在公网上，任何人扫到都能直接访问。

## ③ config.json 权限过宽，里面有数据库口令

默认 644 意味着**同机任何用户都能读**。多用户机器上等于把库交出去。

---

# 二、一键加固（推荐）

包里带了一个加固脚本：

```bash
sudo bash harden-server.sh          # 交互式，每步询问
sudo bash harden-server.sh --yes    # 全部默认执行（首次部署推荐）
sudo bash harden-server.sh --ssh-only  # 只做 SSH 加固
```

## 它做什么

```
① SSH 加固
   禁密码登录（先验证公钥能登录，绝不锁死你）
   禁 root 直接密码登录
   限最大并发连接数

② 防火墙
   只放行 SSH / ARTEX 端口 / HTTP(S)
   其余全部 DROP

③ 权限
   config.json → 600
   pgdata → 700

④ 清理
   /tmp 里超期的渗透残留
```

## 它的设计原则（很重要）

```
绝不锁死操作者   禁密码前先测公钥能不能登
幂等            可以重复跑，不会叠加规则
可回滚          所有改动前备份，结束打印回滚命令
不擅自杀服务    只动防火墙和 sshd，不碰运行中的业务
```

**回滚命令脚本结束会打印**，类似：

```bash
cp -a /etc/ssh/sshd_config.bak.20261101-140000 /etc/ssh/sshd_config
systemctl reload sshd
iptables-restore < /root/fw-backup-20261101-140000.rules
```

**跑之前先把这两行记下来。**

---

# 三、手工加固（逐项说明）

## ① SSH 密钥登录

**先确认你的公钥能登上**，再禁密码——顺序错了一步就进不去。

```bash
# 1. 本地生成密钥（如果还没有）
ssh-keygen -t ed25519 -C "artex-$(hostname)"

# 2. 传到服务器
ssh-copy-id -i ~/.ssh/id_ed25519.pub root@你的服务器IP

# 3. 验证不用密码能登
ssh -i ~/.ssh/id_ed25519 root@你的服务器IP
# 能进，才继续下一步

# 4. 改 sshd 配置
sudo sed -i 's/^#*PasswordAuthentication.*/PasswordAuthentication no/' /etc/ssh/sshd_config
sudo sed -i 's/^#*PermitRootLogin.*/PermitRootLogin prohibit-password/' /etc/ssh/sshd_config
sudo systemctl reload sshd

# 5. 不要断开当前会话！另开一个窗口测能登进去
# 确认能登之后再关旧会话
```

**没验证就关会话 = 把自己锁在外面。**

## ② 防火墙

### ufw（Debian / Ubuntu，最简单）

```bash
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp
sudo ufw allow 8443/tcp        # 如果用 nginx 反代
# 或者直接暴露 ARTEX
# sudo ufw allow 8787/tcp
sudo ufw enable
sudo ufw status verbose
```

### firewalld（CentOS / RHEL）

```bash
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --permanent --add-port=8443/tcp
sudo firewall-cmd --reload
sudo firewall-cmd --list-all
```

### iptables（通用，脚本用的这个）

```bash
# 先备份
sudo iptables-save > /root/fw-backup-$(date +%Y%m%d-%H%M%S).rules

# 基本规则
sudo iptables -A INPUT -i lo -j ACCEPT
sudo iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 22 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 8443 -j ACCEPT
sudo iptables -P INPUT DROP
```

**注意**：很多云厂商（阿里云 / 腾讯云 / AWS）有**安全组**，网络层的规则在控制台里配，iptables 改了也可能被安全组拦。两边都要配。

## ③ ARTEX 的端口策略

**不要直接把 8787 暴露到公网。**

推荐架构：

```
公网 → nginx (443/8443, HTTPS) → 127.0.0.1:8787
```

nginx 配置：

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate     /etc/ssl/artex/artex.crt;
    ssl_certificate_key /etc/ssl/artex/artex.key;
    ssl_protocols       TLSv1.2 TLSv1.3;
    ssl_ciphers         HIGH:!aNULL:!MD5;

    client_max_body_size 500m;
    client_body_timeout  600s;

    # 让后端看到真实 IP
    proxy_set_header Host              $host;
    proxy_set_header X-Real-IP         $remote_addr;
    proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;

    # SSE 流式输出必须关缓冲，否则 agent 的实时思考看不到
    proxy_buffering    off;
    proxy_cache        off;
    proxy_http_version 1.1;
    proxy_set_header   Connection "";

    # 长连接（WebSocket / 隧道）
    proxy_set_header Upgrade    $http_upgrade;
    proxy_read_timeout  3600s;
    proxy_send_timeout  3600s;

    location / {
        proxy_pass http://127.0.0.1:8787;
    }
}
```

然后 iptables/ufw **只放行 443**，8787 只监听 127.0.0.1：

```bash
./artex -addr 127.0.0.1:8787
```

**这样对外只有一个 HTTPS 端口，ARTEX 本身不暴露。**

## ④ 伪装门控

ARTEX 有个反测绘设计：直接访问根路径会返回一个 **nginx 默认页**，让扫描器以为这是台没配置的 web 服务器。

```
直接访问  http://IP:8787/          → nginx 默认页（假的）
正确入口  http://IP:8787/g-xxxxx   → 真正的登录页
```

入口路径在：

```bash
cat data/gate.path
```

**这个路径就是你的第二道密码**，别发到公开地方。

**注意**：门控判断的是**路径**不是来源 IP，所以走 nginx 反代照样生效。

## ⑤ 数据库只监听本地

```bash
# 检查
ss -lntp | grep 5432
```

应该是：

```
LISTEN  127.0.0.1:5432        ← 正确
LISTEN  0.0.0.0:5432          ← 错误，对公网开放
```

如果是 `0.0.0.0`，改 `pgdata/postgresql.conf`：

```
listen_addresses = '127.0.0.1'
```

然后重启 PG。

## ⑥ 文件权限

```bash
chmod 600 config.json          # 里面是数据库口令
chmod 700 pgdata               # 数据库本体
chmod 600 data/jwt.key         # 登录密钥
chmod 700 data/                # 数据目录
```

验证：

```bash
ls -la config.json data/jwt.key
# 应该是 -rw------- 
```

## ⑦ 定期备份

**必须备份的三样**：

```bash
# 1. 数据库（逻辑备份，可跨版本恢复）
pg_dump -h 127.0.0.1 -p 15433 -U artex -d artex \
  -Fc -f backup-$(date +%Y%m%d).dump

# 2. 配置和密钥
tar -czf config-backup-$(date +%Y%m%d).tar.gz \
  config.json jwt.key license.dat

# 3. 知识库（这个丢了最可惜，是打过的经验）
tar -czf kb-backup-$(date +%Y%m%d).tar.gz redteam/
```

**多久备一次**：按你干活频率。每天跑任务就每天备，一周跑一次就一周备。

**备份放哪**：**不要只放本机**。扔到另一台机器、对象存储、或者至少另一个磁盘。

## ⑧ 日志脱敏

日志里可能出现：目标 IP、凭据、抓到的密码、cookie。

```bash
# 分享日志前先扫一遍敏感信息
grep -iE 'password|passwd|secret|token|key.*=' data/artex-err.log | head
```

**给别人的日志，先确认里面没有你的客户信息。**

---

# 四、渗透期间的特殊风险

## agent 起的服务会绑 0.0.0.0

这是设计行为——靶场、OOB 接收端、临时 HTTP 服务需要能被外网访问（比如 DNS 回连、反弹 shell）。

**但它们没有鉴权。**

### 缓解办法

**① 限制来源**

```bash
# 只允许你的 IP 访问高位端口
iptables -A INPUT -p tcp --dport 8000:9999 -s 你的IP -j ACCEPT
iptables -A INPUT -p tcp --dport 8000:9999 -j DROP
```

**② 用完就关**

任务结束后检查有没有残留：

```bash
ss -lntp | grep -vE '127.0.0.1|::1' | grep -E ':(8|9)[0-9]{3}'
```

**③ 定期扫一遍**

```bash
# 列出所有对公网开放的端口
ss -lntp | awk 'NR>1 && $4 !~ /127.0.0.1|::1/ {print $4, $6}'
```

对着这个列表问自己：**每个端口我都认识吗？**

## /tmp 会被写满

渗透脚本大量用 `/tmp` 存中间结果。跑一天可能几个 G。

```bash
# 定期清（保留 24 小时内的）
find /tmp -maxdepth 1 -mindepth 1 -type d -mtime +0 \
  ! -name 'systemd-private-*' ! -name '.X11-unix' -exec rm -rf {} +
```

> `-mtime +1` 是**超过 48 小时**，`-mtime +0` 才是 24 小时以上。别搞错。

---

# 五、加固检查清单

部署完对着过一遍：

```
□ SSH 已禁密码登录，且验证过公钥能登
□ root 不能直接用密码登录
□ 防火墙开了，只放行必要端口
□ 8787 没有直接暴露公网（走 nginx 或只监听 127.0.0.1）
□ HTTPS 证书配好了（不要用自签给客户看）
□ 数据库只监听 127.0.0.1
□ config.json 权限 600
□ pgdata 权限 700
□ gate 路径没外传
□ 备份跑了第一次
□ 知道怎么回滚
```

---

# 六、万一被入侵了

## 判断是否被入侵

```bash
# 1. 谁在登录
who
last -20
lastb | head -20        # 失败登录

# 2. 有没有陌生 SSH key
cat ~/.ssh/authorized_keys

# 3. 有没有陌生进程
ps aux --sort=-%cpu | head -10
ps aux | grep -vE 'systemd|kernel|kworker|bash|sshd'

# 4. 有没有陌生网络连接
ss -tnp | grep ESTAB

# 5. 有没有被改过的文件
find /opt/artex /etc -newer /opt/artex/artex -type f 2>/dev/null | head -20

# 6. 计划任务
crontab -l
ls -la /etc/cron.*

# 7. 有没有陌生用户
cat /etc/passwd | grep -E 'sh$'

# 8. rootkit 检查（如果装了）
rkhunter --check --skip-keypress
```

## 处理顺序

```
1. 先断网（但保留现场）
   iptables -P INPUT DROP
   iptables -P OUTPUT DROP

2. 备份现场（内存、日志、进程列表）
   ps aux > /root/incident-ps.txt
   ss -tnp > /root/incident-net.txt
   tar -czf /root/incident-logs.tar.gz /var/log/

3. 不要急着重启——内存里的东西重启就没了

4. 换机器重建（比清理更靠谱）
   拿备份的数据 → 新机器 → 重新部署 → 改所有密码
```

**被入侵过的机器，清理干净的成本远高于重建。**

---

# 七、一句话总结

```
公网 VPS 三件事必做：
  1. SSH 只留密钥
  2. 防火墙只放必要端口
  3. ARTEX 不裸奔（nginx + HTTPS）

然后：
  备份，备份，备份。
```
