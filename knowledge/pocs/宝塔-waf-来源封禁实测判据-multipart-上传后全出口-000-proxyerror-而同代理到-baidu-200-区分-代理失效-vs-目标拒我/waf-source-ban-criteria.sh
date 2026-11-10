# 判据：multipart 上传 POST 后目标立即拒绝我方来源（宝塔 WAF 行为）
# 1) 上传前：任一出口访问目标应正常
curl.exe -sk --ssl-no-revoke --proxy http://203.0.113.19:17981 \
  -H 'Host: example.com' -A "$UA" -o x \
  -w '%{http_code}|%{size_download}\n' --max-time 25 \
  'https://example.com/example.com'
#   => 200|151628  (正常)

# 2) 发 1 次 multipart 上传（PHPWind 头像上传真实端点）
#    python requests: POST https://example.com/example.com?action=uploadicon&step=2&verify=<vh>&
#    files={"upload": (name, content, "image/gif")}
#    data={"uploadurl[1]":"80","uploadurl[2]":"80","facetype":"3","step":"2"}
#   => 200|0  (请求被接受)

# 3) 上传后立刻复测：同一出口对目标全部失败
#   curl -w '%{http_code}|%{size_download}'  => 000|0
#   python requests                          => ProxyError (连 CONNECT 都建不起来)

# 4) 判据分离（关键）：同批代理到 baidu 仍通 => 代理活着，是目标在拒我们
curl.exe -sk --proxy http://203.0.113.21:17981 -o y -w '%{http_code}|%{size_download}\n' --max-time 12 'https://www.baidu.com/'   # => 200|29506
curl.exe -sk --proxy http://203.0.113.21:17981 -o z -w '%{http_code}|%{size_download}\n' --max-time 15 'https://example.com/example.com'  # => ProxyError

# 5) 直连特征：TCP 443 通，但 HTTP 被 nginx 444 主动断连
Test-NetConnection 203.0.113.77 -Port 443   # => True (TCP 可达，目标在线)
curl.exe -sk -H 'Host: example.com' -o w -w '%{http_code}|%{size_download}\n' --max-time 20 'http://203.0.113.77/example.com'   # => 444|0
curl.exe -sk --ssl-no-revoke -H 'Host: example.com' -o v -w '%{http_code}\n' --max-time 20 'https://203.0.113.77/example.com'  # => 000 (TLS 阶段即失败)
