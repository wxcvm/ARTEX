# 本地文件包含 → RCE（日志包含 / pearcmd / PHP filter chain）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/web-cms/lfi-to-rce.md`
- 分类：web-cms ｜ 严重度：critical ｜ CVSS：- ｜ CVE：-
- 产品：PHP 应用（include/require 带用户输入）
- 影响版本：任意
- 端口：80,443
- 关键词：lfi, 文件包含, 日志包含, pearcmd, filter chain, rce
- 触发点：-

## 定位要点

- 组件指纹：按 `PHP 应用（include/require 带用户输入）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：lfi、文件包含、日志包含、pearcmd、filter chain、rce。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/web-cms/lfi-to-rce.md
