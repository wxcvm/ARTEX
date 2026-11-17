# sourcemap / 前端 JS 泄露接口与密钥（SPA 包网站第一入口）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/sourcemap-js-secret-leak.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.6 ｜ CVE：-
- 产品：前后端分离站点（Vue/React/uni-app）
- 影响版本：任意
- 端口：80,443
- 关键词：sourcemap, js, 前端, 密钥, 接口, 泄露, spa
- 触发点：/static/js/*.js.map 与打包进 bundle 的 AK/secret

## 定位要点

- 组件指纹：按 `前后端分离站点（Vue/React/uni-app）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：sourcemap、js、前端、密钥、接口、泄露、spa。
- 触发条件：/static/js/*.js.map 与打包进 bundle 的 AK/secret。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/sourcemap-js-secret-leak.md
