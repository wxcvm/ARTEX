# 支付通道 / 商户号篡改（充值进攻击者账户）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/payment-channel-swap.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.4 ｜ CVE：-
- 产品：自研充值/支付路由（博彩包网）
- 影响版本：任意
- 端口：80,443
- 关键词：支付通道, 商户号, 充值, 路由, 篡改, 资金
- 触发点：/api/pay/create 的 channelId / merchantId / notifyUrl

## 定位要点

- 组件指纹：按 `自研充值/支付路由（博彩包网）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：支付通道、商户号、充值、路由、篡改、资金。
- 触发条件：/api/pay/create 的 channelId / merchantId / notifyUrl。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/payment-channel-swap.md
