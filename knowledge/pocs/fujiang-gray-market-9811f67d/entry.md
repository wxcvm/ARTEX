# 会员隐私批量查询 / 账号枚举（手机号、银行卡）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/member-data-enumeration.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：7.5 ｜ CVE：-
- 产品：客服查询 / 代理报表 / 会员中心接口
- 影响版本：任意
- 端口：80,443
- 关键词：member, enumeration, 手机号, 银行卡, 隐私, 批量查询
- 触发点：/api/member/info?phone= 逐个返回他人资料

## 定位要点

- 组件指纹：按 `客服查询 / 代理报表 / 会员中心接口` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：member、enumeration、手机号、银行卡、隐私、批量查询。
- 触发条件：/api/member/info?phone= 逐个返回他人资料。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/member-data-enumeration.md
