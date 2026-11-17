# 接口越权 / BOLA（读他人订单、流水、提现）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/api-idor-bola.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.1 ｜ CVE：-
- 产品：自研 API（博彩/代理/跑分平台高发）
- 影响版本：任意
- 端口：80,443
- 关键词：idor, bola, 越权, 水平越权, orderNo, uid, 代理后台
- 触发点：-

## 定位要点

- 组件指纹：按 `自研 API（博彩/代理/跑分平台高发）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：idor、bola、越权、水平越权、orderNo、uid、代理后台。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/api-idor-bola.md
