# 前端签名 / 加密可逆（伪造请求绕过风控）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/client-side-signature-forgery.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.1 ｜ CVE：-
- 产品：H5 / App 自研接口（自定义 sign、AES 请求体）
- 影响版本：任意
- 端口：80,443
- 关键词：sign, hmac, aes, 前端加密, 可逆, nonce, timestamp, 风控
- 触发点：app.js 里抄出算法与 key，重算 sign

## 定位要点

- 组件指纹：按 `H5 / App 自研接口（自定义 sign、AES 请求体）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：sign、hmac、aes、前端加密、可逆、nonce、timestamp、风控。
- 触发条件：app.js 里抄出算法与 key，重算 sign。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/client-side-signature-forgery.md
