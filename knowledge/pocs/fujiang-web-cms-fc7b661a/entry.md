# 命令注入（分号 / 反引号 / $() 与盲注外带）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/web-cms/command-injection.md`
- 分类：web-cms ｜ 严重度：critical ｜ CVSS：- ｜ CVE：-
- 产品：任何调用 shell 的功能（ping/traceroute/转换/压缩）
- 影响版本：任意
- 端口：80,443
- 关键词：命令注入, rce, ping, 反引号, 盲注, shell
- 触发点：-

## 定位要点

- 组件指纹：按 `任何调用 shell 的功能（ping/traceroute/转换/压缩）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：命令注入、rce、ping、反引号、盲注、shell。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/web-cms/command-injection.md
