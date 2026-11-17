# JWT 实现缺陷（alg:none / 弱密钥 / 可预测 / kid 注入）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/jwt-none-and-weak-key.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.1 ｜ CVE：-
- 产品：自研 Web/App 后端（JJWT、java-jwt、PyJWT 等）
- 影响版本：任意
- 端口：80,443
- 关键词：jwt, alg none, hmac, hashcat 16500, kid, 越权, 提权
- 触发点：-

## 定位要点

- 组件指纹：按 `自研 Web/App 后端（JJWT、java-jwt、PyJWT 等）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：jwt、alg none、hmac、hashcat 16500、kid、越权、提权。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/jwt-none-and-weak-key.md
