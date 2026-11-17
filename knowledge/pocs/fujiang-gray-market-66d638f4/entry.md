# 支付回调伪造与金额篡改（资金面）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/payment-notify-forgery.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.1 ｜ CVE：-
- 产品：自研充值/代付回调（博彩、跑分、发卡站高发）
- 影响版本：任意
- 端口：80,443
- 关键词：notify, callback, 伪造回调, 金额篡改, 重放, 验签, 跑分
- 触发点：-

## 定位要点

- 组件指纹：按 `自研充值/代付回调（博彩、跑分、发卡站高发）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：notify、callback、伪造回调、金额篡改、重放、验签、跑分。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/payment-notify-forgery.md
