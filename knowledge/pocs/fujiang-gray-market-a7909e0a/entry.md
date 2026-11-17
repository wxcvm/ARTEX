# 短信轰炸与账号枚举（验证码面）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/sms-bomb-and-enum.md`
- 分类：gray-market ｜ 严重度：medium ｜ CVSS：5.3 ｜ CVE：-
- 产品：自研注册/登录/找回接口
- 影响版本：任意
- 端口：80,443
- 关键词：sms, 轰炸, 枚举, 验证码爆破, 频控, user enumeration
- 触发点：-

## 定位要点

- 组件指纹：按 `自研注册/登录/找回接口` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：sms、轰炸、枚举、验证码爆破、频控、user enumeration。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/sms-bomb-and-enum.md
