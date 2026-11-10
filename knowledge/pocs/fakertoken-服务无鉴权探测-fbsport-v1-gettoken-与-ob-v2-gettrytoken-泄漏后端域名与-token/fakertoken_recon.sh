#!/usr/bin/env bash
# fakerToken 服务无鉴权探测（换站点复用：改 BASE 与 BLOB）
set -u
BASE="${BASE:-https://example.com}"
BLOB="${BLOB:-JFAeOXkY7T76H6p5I6P1p4r1D7So14BcD6jCJ6ILPysVSswgEBBHFOp8bDlV/pEx/qQcoK2DFZVY9a0s2xr5oGxAUI+HBriQsoHKwTDfDL1luH+b/kEOcH0T1X1Sfln1jWG36yuZAS+AdM9bwv0s+HH6GMkBxsAYSziYo39vxEWcsjZDjtwzeENQ/4A0PHWSa7P7hHcwlu6qwwoR82yAGFP+sZLxbNtADDWjTXlskFcn0aYb7ftDeEaeoGPvWUuYhig9fxJ/XzltQSATPf8zO/dXLWTwyIgc7YcR1xYzjkhu069ejxgexqUjO1bhdbURl8DaQOVcB39txxQ2z4EQbVGe5K8zytputcvRijLTdVoKCFWzNTjzHE5JaDg06w5G3d9Zopbd4PK6637xHvQx5vA/OUr9dgVDIQGQDA6EoOWyZwyATU3fP1zeOj9TOTsa7Zo4ybmNGBAhB0najp5sSOxKa4lmqMLC4vAIBIcJinme7vnMMEehKhjUDDAvkuzxLfxiwzcKw/crbnNKKFvxDQG2XsbJHOLM/Dr4tK7puDscdG3iAkwiFdEKFvWq8+/DoqxeVstr53k0Qd+PejVjWeQWIHNOfzmDlv2KWoXaXfZ7eFaqg31KDw9mex5IKO9zsYXS7yzPzHg1l8/6/rr8K8my4IAKD3/RAk3jHJXvVwT+S4E1a1N+Sy+l6dex3gkb}"

echo "=== [1] guest FB token + 后端域名清单（必须是 form，json 会 50040）==="
curl -sk -X POST "$BASE/fbsport/v1/getToken" \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode "requestData=$BLOB" | tee /tmp/faker_fb.json | head -c 2000
echo

echo "=== [2] guest OB token + 可直连 loginUrl ==="
curl -sk -X POST "$BASE/ob/v2/getTryToken" \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode "product=rb" --data-urlencode "requestData=$BLOB" | head -c 2000
echo

echo "=== [3] 真实联赛数据（无鉴权）==="
curl -sk -X POST "$BASE/common/v1/hotLegues/getHotLegues" \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode "requestData=AAAA" | head -c 500
echo

echo "=== [4] 提取全部后端域名 ==="
grep -oE '"(https?|wss?)://[A-Za-z0-9.-]+"?' /tmp/faker_fb.json | tr -d '"' | sort -u
