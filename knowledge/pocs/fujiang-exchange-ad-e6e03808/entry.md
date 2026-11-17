# AD CS 证书滥用（ESC1 / ESC8 等）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/exchange-ad/adcs-esc1-esc8.md`
- 分类：exchange-ad ｜ 严重度：critical ｜ CVSS：- ｜ CVE：-
- 产品：Active Directory Certificate Services
- 影响版本：任意 AD CS 部署（配置缺陷）
- 端口：80,443,135,445
- 关键词：adcs, esc1, esc8, certipy, pki, certificate-template
- 触发点：-

## 定位要点

- 组件指纹：按 `Active Directory Certificate Services` 的版本声明判断候选面（任意 AD CS 部署（配置缺陷）），端口 80,443,135,445 优先。
- 信号来源：adcs、esc1、esc8、certipy、pki、certificate-template。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/exchange-ad/adcs-esc1-esc8.md
