# OAuth redirect_uri 劫持 / state 缺失（账号接管）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/oauth-redirect-uri.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.6 ｜ CVE：-
- 产品：自研三方登录（Google/微信/Telegram）
- 影响版本：任意
- 端口：80,443
- 关键词：oauth, redirect_uri, state, 登录, 接管, sso
- 触发点：/oauth/authorize?redirect_uri= 可篡改

## 定位要点

- 组件指纹：按 `自研三方登录（Google/微信/Telegram）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：oauth、redirect_uri、state、登录、接管、sso。
- 触发条件：/oauth/authorize?redirect_uri= 可篡改。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/oauth-redirect-uri.md
