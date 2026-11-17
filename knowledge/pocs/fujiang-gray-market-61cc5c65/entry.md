# 开奖结果可预测（弱 PRNG / 客户端算法 / 固定种子）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/lottery-prng-predictable.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.1 ｜ CVE：-
- 产品：自营彩票 / 快三 / 28 类玩法
- 影响版本：任意
- 端口：80,443
- 关键词：lottery, prng, 开奖, 可预测, 种子, 快三, 时时彩
- 触发点：历史期号统计 + 前端 JS 算法复现

## 定位要点

- 组件指纹：按 `自营彩票 / 快三 / 28 类玩法` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：lottery、prng、开奖、可预测、种子、快三、时时彩。
- 触发条件：历史期号统计 + 前端 JS 算法复现。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/lottery-prng-predictable.md
