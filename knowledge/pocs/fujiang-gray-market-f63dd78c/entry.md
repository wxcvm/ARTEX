# 源码 / 备份 / 配置泄露

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/repo-config-backup-leak.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：7.5 ｜ CVE：-
- 产品：通用 Web（PHP/Java/Node 自研站）
- 影响版本：任意
- 端口：80,443,8080
- 关键词：.git, .svn, env, backup, zip, sourcemap, swagger, 源码泄露
- 触发点：-

## 定位要点

- 组件指纹：按 `通用 Web（PHP/Java/Node 自研站）` 的版本声明判断候选面（任意），端口 80,443,8080 优先。
- 信号来源：.git、.svn、env、backup、zip、sourcemap、swagger、源码泄露。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/repo-config-backup-leak.md
