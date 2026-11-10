#!/usr/bin/env bash
# 对象存储匿名读判定：ListBucket 被拒 ≠ 不可读；必须做「随机键」对象级分叉
# 用法：BUCKET=https://target500.example.com bash oss_anon_read_probe.sh
#      已知键用 REALKEY 传入，多个用空格
set -u
BUCKET="${BUCKET:?set BUCKET=https://<bucket>.<region>.endpoint}"
REALKEY="${REALKEY:-}"
C="curl -sS --noproxy * --max-time 20"
RAND="zz$(date +%s | md5sum | cut -c1-16).png"

echo "=== [1] ListBucket：403 的正文才是信息，别只看状态码 ==="
BODY=$($C "$BUCKET/?list-type=2&max-keys=1000")
CODE=$($C -o /dev/null -w '%{http_code}' "$BUCKET/?list-type=2&max-keys=1000")
echo "HTTP $CODE"
echo "$BODY" | head -c 400; echo
case "$BODY" in
  *NoSuchBucket*) echo "  => 桶不存在，换桶名" ;;
  *ListBucketResult*) echo "  => 桶可列举！直接枚举 <Key>" ;;
  *AccessDenied*|*ACL*|*permission*|*acl*)
      echo "  => 桶存在、匿名禁列（属正常）；**继续第 2 步判断对象级**" ;;
  *) echo "  => 非标准回显，人工判断" ;;
esac

echo "=== [2] 对象级分叉（决定性判据）==="
B2=$($C -o /dev/null -w '%{http_code}' "$BUCKET/$RAND" -D /tmp/h2); H2=$(head -c 300 /tmp/h2)
echo "GET /$RAND -> HTTP $B2"
case "$B2" in
  404) echo "  => NoSuchKey：GetObject 已放行，桶可匿名取对象（不可枚举而已）" ;;
  403) echo "  => 对象级也私有（bucket acl）。看正文是否含 'bucket acl' 佐证" ;;
  200) echo "  => 命中同名对象（换更长的随机名重测）" ;;
esac
[ "$B2" = "403" ] && { echo "  body:"; echo "$B2$H2" | head -c 400; echo; }

echo "=== [3] 真实键取证：HEAD 长度 + Range 首片 + 魔数 ==="
for K in $REALKEY; do
  echo "--- $K"
  $C -I "$BUCKET/$K" | egrep -i 'HTTP/|content-length|content-type'
  $C -r 0-1023 -o /tmp/head.bin -D /tmp/h3 "$BUCKET/$K"
  egrep -i 'HTTP/|content-type|content-range' /tmp/h3
  echo -n "  magic: "; xxd -l 8 /tmp/head.bin 2>/dev/null | cut -d' ' -f2-5
done
# 判据：200/206 + 正确 Content-Type + 魔数吻合（zip 504b0304 / png 89504e47 / jpg ffd8ff）= 实锤匿名读
