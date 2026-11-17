# 余额并发与负数（竞态 / 金额符号）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/race-condition-balance.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.0 ｜ CVE：-
- 产品：自研钱包与投注接口
- 影响版本：任意
- 端口：80,443
- 关键词：竞态, race condition, 并发提现, 负数, 金额校验, 幂等
- 触发点：同一请求并发 20 次 / amount=-100

## 定位要点

- 组件指纹：按 `自研钱包与投注接口` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：竞态、race condition、并发提现、负数、金额校验、幂等。
- 触发条件：同一请求并发 20 次 / amount=-100。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/race-condition-balance.md
