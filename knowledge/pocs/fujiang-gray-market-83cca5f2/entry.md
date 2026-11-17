# Swagger / Knife4j 接口文档暴露 → 未授权接口直打

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/swagger-knife4j-unauth.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.3 ｜ CVE：-
- 产品：Spring Boot 系（swagger2/3、knife4j、springdoc）
- 影响版本：任意
- 端口：80,443,8080
- 关键词：swagger, knife4j, springdoc, 文档, 未授权, api
- 触发点：/v2/api-docs  /doc.html  /swagger-ui/index.html

## 定位要点

- 组件指纹：按 `Spring Boot 系（swagger2/3、knife4j、springdoc）` 的版本声明判断候选面（任意），端口 80,443,8080 优先。
- 信号来源：swagger、knife4j、springdoc、文档、未授权、api。
- 触发条件：/v2/api-docs  /doc.html  /swagger-ui/index.html。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/swagger-knife4j-unauth.md
