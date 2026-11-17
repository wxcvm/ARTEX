# 会员资料存储型 XSS 打后台审核（管理员会话劫持）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/stored-xss-to-admin.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.2 ｜ CVE：-
- 产品：自研后台（客服/审核/风控页）
- 影响版本：任意
- 端口：80,443
- 关键词：stored xss, 存储型, 后台, 审核, cookie, csrf
- 触发点：昵称/备注/工单里插 <img onerror>

## 定位要点

- 组件指纹：按 `自研后台（客服/审核/风控页）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：stored xss、存储型、后台、审核、cookie、csrf。
- 触发条件：昵称/备注/工单里插 <img onerror>。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/stored-xss-to-admin.md
