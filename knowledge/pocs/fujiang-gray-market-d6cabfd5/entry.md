# 源站 IP 泄露 → 绕过 WAF 直打（博彩站必看）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/origin-ip-leak-waf-bypass.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.0 ｜ CVE：-
- 产品：任何套 CDN/WAF 的站点
- 影响版本：任意
- 端口：80,443,8080
- 关键词：源站, cdn, waf, 绕过, 真实ip, 回源
- 触发点：子域/邮件头/证书历史/SSRF 泄露的 origin

## 定位要点

- 组件指纹：按 `任何套 CDN/WAF 的站点` 的版本声明判断候选面（任意），端口 80,443,8080 优先。
- 信号来源：源站、cdn、waf、绕过、真实ip、回源。
- 触发条件：子域/邮件头/证书历史/SSRF 泄露的 origin。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/origin-ip-leak-waf-bypass.md
