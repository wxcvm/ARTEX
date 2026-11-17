# 支付 / 回调地址可控 → 内网 SSRF（含云元数据）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/callback-url-ssrf.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.5 ｜ CVE：-
- 产品：自研支付与通知模块
- 影响版本：任意
- 端口：80,443
- 关键词：回调, 支付, ssrf, 内网, 元数据, notify
- 触发点：支付配置/回调 URL/代理推广链接检测

## 定位要点

- 组件指纹：按 `自研支付与通知模块` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：回调、支付、ssrf、内网、元数据、notify。
- 触发条件：支付配置/回调 URL/代理推广链接检测。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/callback-url-ssrf.md
