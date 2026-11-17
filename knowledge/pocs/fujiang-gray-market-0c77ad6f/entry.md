# 提现审核接口越权（读他人银行卡 / 改审核状态）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/withdraw-approval-idor.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.1 ｜ CVE：-
- 产品：博彩包网 / 跑分代付
- 影响版本：任意
- 端口：80,443
- 关键词：提现, withdraw, 审核, approve, 银行卡, 越权
- 触发点：两形态二选一（Query/JSON），抓包为准

## 定位要点

- 组件指纹：按 `博彩包网 / 跑分代付` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：提现、withdraw、审核、approve、银行卡、越权。
- 触发条件：两形态二选一（Query/JSON），抓包为准。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/withdraw-approval-idor.md
