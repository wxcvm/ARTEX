# 后台 IP 白名单绕过（X-Forwarded-For 伪造内网来源）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/admin-ip-whitelist-bypass.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.4 ｜ CVE：-
- 产品：自研后台网关（只校验 XFF 或只查首段 IP）
- 影响版本：任意
- 端口：80,443
- 关键词：白名单, ip, xff, 绕过, 后台, 内网
- 触发点：403 页面 + 只允许 127.0.0.1/内网段

## 定位要点

- 组件指纹：按 `自研后台网关（只校验 XFF 或只查首段 IP）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：白名单、ip、xff、绕过、后台、内网。
- 触发条件：403 页面 + 只允许 127.0.0.1/内网段。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/admin-ip-whitelist-bypass.md
