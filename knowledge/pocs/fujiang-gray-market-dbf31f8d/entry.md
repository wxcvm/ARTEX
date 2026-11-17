# 三方游戏 / 视讯接口 key 泄露与回调伪造

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/game-vendor-api-forgery.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：8.8 ｜ CVE：-
- 产品：接入 AG / BBIN / OB 等供应商的包网站
- 影响版本：任意
- 端口：80,443
- 关键词：game vendor, api key, 视讯, 上分, callback, transfer
- 触发点：前端拿 vendor key，伪造 /api/game/callback 上分

## 定位要点

- 组件指纹：按 `接入 AG / BBIN / OB 等供应商的包网站` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：game vendor、api key、视讯、上分、callback、transfer。
- 触发条件：前端拿 vendor key，伪造 /api/game/callback 上分。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/game-vendor-api-forgery.md
