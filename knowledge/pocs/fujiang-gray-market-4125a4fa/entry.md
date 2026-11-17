# SSRF（图片抓取 / 回调校验 / 推广链接检测）打内网与云元数据

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/ssrf-to-internal.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.2 ｜ CVE：-
- 产品：自研后端（抓图、URL 预览、支付回调校验）
- 影响版本：任意
- 端口：80,443
- 关键词：ssrf, 内网, 元数据, 169.254.169.254, redis, gopher
- 触发点：url=http://169.254.169.254/latest/meta-data/

## 定位要点

- 组件指纹：按 `自研后端（抓图、URL 预览、支付回调校验）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：ssrf、内网、元数据、169.254.169.254、redis、gopher。
- 触发条件：url=http://169.254.169.254/latest/meta-data/。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/ssrf-to-internal.md
