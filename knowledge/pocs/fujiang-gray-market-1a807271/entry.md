# 报表/财务接口未授权（全站流水与代理佣金）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/report-api-unauth.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.6 ｜ CVE：-
- 产品：包网/代理后台报表模块
- 影响版本：任意
- 端口：80,443
- 关键词：report, finance, 报表, 流水, 佣金, 未授权
- 触发点：/api/report/* 免认证或普通 token

## 定位要点

- 组件指纹：按 `包网/代理后台报表模块` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：report、finance、报表、流水、佣金、未授权。
- 触发条件：/api/report/* 免认证或普通 token。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/report-api-unauth.md
