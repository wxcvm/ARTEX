# 短信验证码回显 / 复用 / 爆破（任意账号接管）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/sms-code-reuse-and-bruteforce.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.1 ｜ CVE：-
- 产品：自研注册/登录/改密接口
- 影响版本：任意
- 端口：80,443
- 关键词：sms, code, 验证码回显, 复用, 爆破, 接管
- 触发点：响应里带 code / 同码复用 / 4 位爆破

## 定位要点

- 组件指纹：按 `自研注册/登录/改密接口` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：sms、code、验证码回显、复用、爆破、接管。
- 触发条件：响应里带 code / 同码复用 / 4 位爆破。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/sms-code-reuse-and-bruteforce.md
