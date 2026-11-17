# 包网多租户越权（tenantId/siteId 遍历，横穿几十个站点）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/multi-tenant-cross-site.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.2 ｜ CVE：-
- 产品：包网/白标系统（一套后台托管多站点）
- 影响版本：任意
- 端口：80,443
- 关键词：多租户, tenant, siteId, 包网, 白标, 越权
- 触发点：请求体或 Header 里的 tenantId / siteId / merchantId

## 定位要点

- 组件指纹：按 `包网/白标系统（一套后台托管多站点）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：多租户、tenant、siteId、包网、白标、越权。
- 触发条件：请求体或 Header 里的 tenantId / siteId / merchantId。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/multi-tenant-cross-site.md
