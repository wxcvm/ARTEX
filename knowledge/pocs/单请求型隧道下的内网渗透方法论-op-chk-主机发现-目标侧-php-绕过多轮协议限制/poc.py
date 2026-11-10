# 场景
只有一条「一次 HTTP 请求 = 一次内网请求-响应」的自建隧道（无持久连接、单 Session 串行，
典型于 PHP/ASP 类 WebShell 隧道；suo5 只支持 Java/.NET 服务端时只能自建）。
此时 gogo/fscan 直接 -t 高并发会大面积超时并伪造"端口关闭"假象（实测 gogo 每端口约 9s）。

# 技巧一：用隧道端点的 op=chk 做低成本主机发现（0.4s/次）
隧道端点若支持"端口连通性判定"模式，直接 HTTP GET 它，用**错误类型**区分状态：
  OK                       -> 端口开放
  E:Connection refused     -> 主机存活、端口关闭   ★ 存活判据
  E:No route to host       -> 同段该 IP 无主机
  E:Connection timed out   -> 跨段被过滤（注意：默认路由同样 timeout，**不可据此判定新网段存在**）
批量脚本见 runs/rt_net.py（多线程 6 即可，每个请求在目标侧独立建连）。

# 技巧二：★ 把探测逻辑搬到目标侧 PHP 执行（绕开 op=req 限制的核心）
WebShell 所在主机的 PHP 往往具备完整网络扩展。先枚举能力：
  mysql_connect / mysqli_connect / ftp_connect / curl_init / socket_create / fsockopen / pfsockopen
PHP 进程在内网直连，**由 PHP 自己维持同一条 TCP 完成多轮往返**，结果一次性回显
=> FTP/MySQL/RPC/NFS/任意多轮协议全部可用，且内网 RTT≈0，远快于经隧道。
命令执行若被 safe_mode 禁掉不影响这条路（网络函数不属命令执行）。

# 技巧三：马受 open_basedir 限制时，改用 MySQL 的 LOAD_FILE 读系统文件
PHP open_basedir（如 docroot:/tmp:/some/ad/）会挡掉 /etc、/proc；
但 mysqld 的 LOAD_FILE 不受其限制，且 5.0 无 secure_file_priv：
  SELECT LOAD_FILE('/etc/hosts');        -- 主机名/别名
  SELECT LOAD_FILE('/etc/sysconfig/network-scripts/ifcfg-eth1');  -- ★ 网段铁证
  SELECT LOAD_FILE('/etc/hosts.allow');  -- 运维白名单
  SELECT GROUP_CONCAT(DISTINCT host) FROM mysql.user;             -- ★ 内网客户端清单
注意：/proc/* 用 LOAD_FILE 读回空（procfs size=0），不要依赖它。

# 技巧四：RPC/NFS 枚举（无 mountd 时）
PMAP_DUMP(prog=100000 vers=2 proc=4) 一次列出全部注册服务及端口；
RPC over TCP 需 4 字节 record marking（0x80000000|len）；XDR 头 28 字节后接链表。
若列表中无 mountd(100005)，则 showmount/MOUNT EXPORT 不通，只剩 NFSv4 伪根(PUTROOTFH=24)；
AUTH_NULL 可能 MSG_DENIED，可再试 AUTH_UNIX(flavor=1, stamp+ns(machine)+uid+gid+gids_len)。
