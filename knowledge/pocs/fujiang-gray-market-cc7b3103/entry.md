# 后台“重置任意用户密码”接口越权（直接接管会员/管理员）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/admin-reset-any-user.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.2 ｜ CVE：-
- 产品：自研后台（会员管理/代理管理）
- 影响版本：任意
- 端口：80,443
- 关键词：重置密码, 任意用户, 越权, 后台, 接管
- 触发点：/api/admin/user/resetPassword  /api/agent/resetPwd

## 定位要点

- 组件指纹：按 `自研后台（会员管理/代理管理）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：重置密码、任意用户、越权、后台、接管。
- 触发条件：/api/admin/user/resetPassword  /api/agent/resetPwd。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/admin-reset-any-user.md
