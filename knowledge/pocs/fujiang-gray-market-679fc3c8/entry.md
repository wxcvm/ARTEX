# 订单状态机篡改（未支付改成已支付 / 已取消再入账）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/order-status-tamper.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.3 ｜ CVE：-
- 产品：自研充值/订单系统（资金面）
- 影响版本：任意
- 端口：80,443
- 关键词：订单, 状态机, 篡改, 未支付, 充值, 资金
- 触发点：/api/order/status  /api/recharge/confirm

## 定位要点

- 组件指纹：按 `自研充值/订单系统（资金面）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：订单、状态机、篡改、未支付、充值、资金。
- 触发条件：/api/order/status  /api/recharge/confirm。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/order-status-tamper.md
