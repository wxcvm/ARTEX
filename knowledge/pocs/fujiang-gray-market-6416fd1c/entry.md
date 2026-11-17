# 财务导出接口越权（银行卡 + 身份证批量导出）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/finance-export-pii.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.0 ｜ CVE：-
- 产品：自研财务/提现模块（资金面）
- 影响版本：任意
- 端口：80,443
- 关键词：财务, 导出, 银行卡, 身份证, 提现, 越权
- 触发点：/api/finance/export  /api/withdraw/export

## 定位要点

- 组件指纹：按 `自研财务/提现模块（资金面）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：财务、导出、银行卡、身份证、提现、越权。
- 触发条件：/api/finance/export  /api/withdraw/export。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/finance-export-pii.md
