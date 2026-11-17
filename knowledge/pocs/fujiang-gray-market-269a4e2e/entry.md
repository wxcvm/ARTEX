# 下注/结算赔率与倍数篡改（资金面）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/bet-odds-tampering.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.0 ｜ CVE：-
- 产品：博彩/彩票投注接口
- 影响版本：任意
- 端口：80,443
- 关键词：bet, odds, 赔率, 倍数, 结算, 篡改
- 触发点：下注体里改 odds/multiple

## 定位要点

- 组件指纹：按 `博彩/彩票投注接口` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：bet、odds、赔率、倍数、结算、篡改。
- 触发条件：下注体里改 odds/multiple。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/bet-odds-tampering.md
