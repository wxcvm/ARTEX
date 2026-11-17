# 提现银行卡越权替换（把别人提现收进自己卡）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/withdraw-bankcard-swap.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.3 ｜ CVE：-
- 产品：自研提现模块（博彩包网）
- 影响版本：任意
- 端口：80,443
- 关键词：提现, 银行卡, 越权, 替换, 资金, withdraw
- 触发点：/api/withdraw/apply 的 bankId / cardNo / memberId

## 定位要点

- 组件指纹：按 `自研提现模块（博彩包网）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：提现、银行卡、越权、替换、资金、withdraw。
- 触发条件：/api/withdraw/apply 的 bankId / cardNo / memberId。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/withdraw-bankcard-swap.md
