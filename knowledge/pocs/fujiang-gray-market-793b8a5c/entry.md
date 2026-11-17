# 通知模板注入（标题/正文模板未转义 → SSTI / 存储型 XSS）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/notification-template-ssti.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.0 ｜ CVE：-
- 产品：自研消息中心（站内信/邮件/短信模板）
- 影响版本：任意
- 端口：80,443
- 关键词：模板, ssti, 通知, 站内信, 短信, 注入
- 触发点：后台模板编辑 / 会员自定义通知标题

## 定位要点

- 组件指纹：按 `自研消息中心（站内信/邮件/短信模板）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：模板、ssti、通知、站内信、短信、注入。
- 触发条件：后台模板编辑 / 会员自定义通知标题。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/notification-template-ssti.md
