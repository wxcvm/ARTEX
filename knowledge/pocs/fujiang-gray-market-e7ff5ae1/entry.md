# 若依（RuoYi）系后台默认口令 / Druid / 定时任务 RCE

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/ruoyi-admin-unauth.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.0 ｜ CVE：-
- 产品：若依 RuoYi（及二开的包网/代理后台）
- 影响版本：4.x / 3.x / 二开变体
- 端口：80,443,8080,9090
- 关键词：ruoyi, 若依, druid, admin123, quartz, 定时任务, prod-api
- 触发点：-

## 定位要点

- 组件指纹：按 `若依 RuoYi（及二开的包网/代理后台）` 的版本声明判断候选面（4.x / 3.x / 二开变体），端口 80,443,8080,9090 优先。
- 信号来源：ruoyi、若依、druid、admin123、quartz、定时任务、prod-api。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/ruoyi-admin-unauth.md
