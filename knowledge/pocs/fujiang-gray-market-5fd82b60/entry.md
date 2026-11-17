# 接口无 nonce / 时间戳 → 请求重放（充值、领奖、下注）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/api-replay-nonce-missing.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.2 ｜ CVE：-
- 产品：自研 App/小程序 API
- 影响版本：任意
- 端口：80,443
- 关键词：重放, replay, nonce, 时间戳, 签名, 幂等
- 触发点：充值/领奖/下注接口缺 nonce 校验

## 定位要点

- 组件指纹：按 `自研 App/小程序 API` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：重放、replay、nonce、时间戳、签名、幂等。
- 触发条件：充值/领奖/下注接口缺 nonce 校验。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/api-replay-nonce-missing.md
