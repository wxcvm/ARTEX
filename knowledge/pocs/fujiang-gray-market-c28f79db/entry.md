# 后台管理 API 未授权（只在前端隐藏菜单）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/admin-api-unauthenticated.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.1 ｜ CVE：-
- 产品：自研包网/代理后台
- 影响版本：任意
- 端口：80,443
- 关键词：admin api, 未授权, 越权, 前端隐藏, 菜单
- 触发点：普通 token 直接打 /api/admin/*

## 定位要点

- 组件指纹：按 `自研包网/代理后台` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：admin api、未授权、越权、前端隐藏、菜单。
- 触发条件：普通 token 直接打 /api/admin/*。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/admin-api-unauthenticated.md
