# phpMyAdmin / Adminer 暴露（弱口令直进数据库）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/db-admin-tool-exposed.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.0 ｜ CVE：-
- 产品：phpMyAdmin / Adminer / 宝塔 pma
- 影响版本：任意
- 端口：80,443,888,8080
- 关键词：phpmyadmin, adminer, pma, 数据库, 弱口令, 脱库
- 触发点：/phpmyadmin/  /pma/  /adminer.php

## 定位要点

- 组件指纹：按 `phpMyAdmin / Adminer / 宝塔 pma` 的版本声明判断候选面（任意），端口 80,443,888,8080 优先。
- 信号来源：phpmyadmin、adminer、pma、数据库、弱口令、脱库。
- 触发条件：/phpmyadmin/  /pma/  /adminer.php。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/db-admin-tool-exposed.md
