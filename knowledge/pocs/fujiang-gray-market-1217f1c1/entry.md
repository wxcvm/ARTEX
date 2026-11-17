# 运营后台弱口令 + 验证码复用（博彩站点最常见入口）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/operator-backend-weak-login.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.1 ｜ CVE：-
- 产品：自研运营/代理后台（包网系统）
- 影响版本：任意
- 端口：80,443,8080,8443
- 关键词：弱口令, 后台, 验证码复用, 运营商, 包网, 登录
- 触发点：/admin /manage /agent 的登录接口

## 定位要点

- 组件指纹：按 `自研运营/代理后台（包网系统）` 的版本声明判断候选面（任意），端口 80,443,8080,8443 优先。
- 信号来源：弱口令、后台、验证码复用、运营商、包网、登录。
- 触发条件：/admin /manage /agent 的登录接口。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/operator-backend-weak-login.md
