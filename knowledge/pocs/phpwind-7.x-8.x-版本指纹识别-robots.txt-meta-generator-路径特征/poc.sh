#!/usr/bin/env bash
# phpwind-fingerprint.sh — PHPWind 版本与类型识别
# 用法: ./phpwind-fingerprint.sh https://example.com/bbs
set -u
B="${1:?usage: $0 <base-url>}"
say(){ printf '\n=== %s ===\n' "$*"; }

say "robots.txt（自带版本声明）"
curl -sk --max-time 15 "$B/robots.txt" | head -8

say "meta generator"
curl -sk --max-time 15 "$B/" | grep -oiE '<meta[^>]*generator[^>]*>' | head -3

say "example.com footer 版本（后台页面）"
curl -sk --max-time 15 "$B/example.com" | grep -oiE 'PHPwind[^<]{0,20}' | head -3

say "PHPWind 类型特征计数"
BODY=$(curl -sk --max-time 15 "$B/")
for k in pw_ajax.js 'images/wind' verifyhash 'fid=' 'read.php' 'u.php' ; do
  printf '%-14s %s\n' "$k" "$(printf '%s' "$BODY" | grep -oF "$k" | wc -l)"
done

say "Discuz 特征（应为 0，用于区分）"
for k in forum.php home.php '_saltkey' 'static/js' ; do
  printf '%-14s %s\n' "$k" "$(printf '%s' "$BODY" | grep -oF "$k" | wc -l)"
done

say "关键路径状态码"
for p in example.com login.php register.php example.com require/ style/ attachment/ simple/ wap/ uc_server/ install/ config/ data/ api/ ; do
  printf '%-14s %s\n' "$p" "$(curl -sk -o /dev/null -w '%{http_code}' --max-time 10 "$B/$p")"
done

say "Cookie 前缀（泄露表前缀）"
curl -sk -D- -o /dev/null --max-time 15 "$B/" | grep -i '^set-cookie' | sed 's/;.*//'