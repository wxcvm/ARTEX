# Mass Assignment（注册/资料接口注入 role / vip / balance / agentId）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/mass-assignment-role.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.8 ｜ CVE：-
- 产品：自研用户/代理接口
- 影响版本：任意
- 端口：80,443
- 关键词：mass assignment, 批量赋值, role, vip, balance, agentId
- 触发点：注册体里多带一个 role=admin

## 定位要点

- 组件指纹：按 `自研用户/代理接口` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：mass assignment、批量赋值、role、vip、balance、agentId。
- 触发条件：注册体里多带一个 role=admin。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/mass-assignment-role.md
