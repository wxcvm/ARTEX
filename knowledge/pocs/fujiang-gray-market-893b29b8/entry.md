# Druid 监控未授权 → session 泄露/后台接管（若依系必备）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/druid-monitor-unauth.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.0 ｜ CVE：-
- 产品：Alibaba Druid 连接池监控（若依/jeecg 等）
- 影响版本：任意
- 端口：80,443,8080
- 关键词：druid, 监控, session, 未授权, 若依, jeecg
- 触发点：/druid/index.html  /druid/websession.html

## 定位要点

- 组件指纹：按 `Alibaba Druid 连接池监控（若依/jeecg 等）` 的版本声明判断候选面（任意），端口 80,443,8080 优先。
- 信号来源：druid、监控、session、未授权、若依、jeecg。
- 触发条件：/druid/index.html  /druid/websession.html。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/druid-monitor-unauth.md
