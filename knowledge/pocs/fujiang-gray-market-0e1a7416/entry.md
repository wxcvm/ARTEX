# 密码重置 token 可预测 / 未绑定用户（任意用户改密）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/password-reset-token-abuse.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.1 ｜ CVE：-
- 产品：自研找回密码流程
- 影响版本：任意
- 端口：80,443
- 关键词：reset, token, 找回密码, 任意用户, 越权
- 触发点：token = md5(uid) / 响应回显 / 不校验 uid

## 定位要点

- 组件指纹：按 `自研找回密码流程` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：reset、token、找回密码、任意用户、越权。
- 触发条件：token = md5(uid) / 响应回显 / 不校验 uid。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/password-reset-token-abuse.md
