# Host 头投毒（密码重置链接被劫持 → 任意账号接管）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/host-header-poisoning.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.1 ｜ CVE：-
- 产品：自研找回密码/邀请链接（用 Host 拼 URL）
- 影响版本：任意
- 端口：80,443
- 关键词：host, 头, 投毒, 密码重置, 链接, 劫持
- 触发点：POST /api/user/findpwd 的 Host / X-Forwarded-Host

## 定位要点

- 组件指纹：按 `自研找回密码/邀请链接（用 Host 拼 URL）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：host、头、投毒、密码重置、链接、劫持。
- 触发条件：POST /api/user/findpwd 的 Host / X-Forwarded-Host。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/host-header-poisoning.md
