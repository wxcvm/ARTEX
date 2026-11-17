# CSRF（无 token 的提现/改密/绑卡）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/csrf-transfer.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.0 ｜ CVE：-
- 产品：自研会员中心（Cookie 会话）
- 影响版本：任意
- 端口：80,443
- 关键词：csrf, 跨站, 提现, 改密, 绑卡, token
- 触发点：关键操作缺 CSRF token 且 SameSite 未设

## 定位要点

- 组件指纹：按 `自研会员中心（Cookie 会话）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：csrf、跨站、提现、改密、绑卡、token。
- 触发条件：关键操作缺 CSRF token 且 SameSite 未设。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/csrf-transfer.md
