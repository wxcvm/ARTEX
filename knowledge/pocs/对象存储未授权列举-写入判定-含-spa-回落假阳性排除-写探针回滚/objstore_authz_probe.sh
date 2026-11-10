#!/usr/bin/env bash
# 对象存储未授权「列举 + 写入」判定（换目标只改 BUCKET）
# 判据：列举必须以是否返回 <ListBucketResult> XML 为准，不能只看 200 ——
#       很多 S3 静态托管会把任意路径回落成 SPA 首页（200 + text/html），那是假阳性。
set -u
BUCKET="${BUCKET:-example.com}"
PROBE="write-probe-$RANDOM$RANDOM.txt"

echo "=== [1] 列举判定（看是否 ListBucketResult XML）==="
for Q in '/?list-type=2' '/?list-type=2&max-keys=1000' '/?delimiter=/' '/?prefix=' '/?acl' '/?location' '/?versioning' '/?uploads'; do
  BODY=$(curl -sk --max-time 20 "https://$BUCKET$Q")
  CODE=$(curl -sk -o /dev/null -w '%{http_code}' --max-time 20 "https://$BUCKET$Q")
  if echo "$BODY" | grep -q 'ListBucketResult'; then
    echo "  [命中] $Q -> $CODE  桶可列举！对象数：$(echo "$BODY" | grep -o '<Key>' | wc -l)"
    echo "$BODY" | grep -o '<Key>[^<]*</Key>' | head -30
  else
    echo "  [未命中] $Q -> $CODE （无 ListBucketResult，可能回落 SPA）"
  fi
done

echo "=== [2] 写权限判定（无害探针，成功立刻删除回滚）==="
PC=$(curl -sk -o /dev/null -w '%{http_code}' -X PUT --max-time 20 \
     -H 'Content-Type: text/plain' --data 'AUTHZ_PROBE_INERT_CONTENT' "https://$BUCKET/$PROBE")
echo "  PUT /$PROBE -> $PC"
GC=$(curl -sk -o /dev/null -w '%{http_code}' --max-time 20 "https://$BUCKET/$PROBE")
echo "  GET /$PROBE -> $GC"
if [ "$PC" = "200" ] || [ "$PC" = "204" ]; then
  echo "  [命中] 桶可写！立即回滚删除"
  DC=$(curl -sk -o /dev/null -w '%{http_code}' -X DELETE --max-time 20 "https://$BUCKET/$PROBE")
  echo "  DELETE /$PROBE -> $DC（回滚完成）"
else
  echo "  [未命中] PUT 被拒，未产生对象，无需删除"
fi

echo "=== [3] 排除 nginx 健康点假象 ==="
echo "  若根路径返回极短纯文本（如 3 字节 'ok'）且 PUT 返回 405，那是 nginx 健康检查点，不是对象存储"
curl -sk --max-time 15 "https://$BUCKET/" | head -c 120; echo
