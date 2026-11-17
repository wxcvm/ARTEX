# PHP 反序列化（POP 链 → RCE / phar 触发）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/web-cms/php-object-injection.md`
- 分类：web-cms ｜ 严重度：critical ｜ CVSS：- ｜ CVE：-
- 产品：Laravel / Symfony / ThinkPHP / 自研 PHP
- 影响版本：任意
- 端口：80,443
- 关键词：php, 反序列化, pop链, phar, unserialize, rce
- 触发点：-

## 定位要点

- 组件指纹：按 `Laravel / Symfony / ThinkPHP / 自研 PHP` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：php、反序列化、pop链、phar、unserialize、rce。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/web-cms/php-object-injection.md
