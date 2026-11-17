# 后台定时任务 / 计划任务注入（Quartz、XXL-JOB、若依）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/scheduled-task-injection.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.1 ｜ CVE：-
- 产品：若依/XXL-JOB/Quartz 管理后台
- 影响版本：任意
- 端口：80,443,8080
- 关键词：定时任务, quartz, xxl-job, 若依, rce, 计划任务
- 触发点：任务执行内容 / cron 表达式可写

## 定位要点

- 组件指纹：按 `若依/XXL-JOB/Quartz 管理后台` 的版本声明判断候选面（任意），端口 80,443,8080 优先。
- 信号来源：定时任务、quartz、xxl-job、若依、rce、计划任务。
- 触发条件：任务执行内容 / cron 表达式可写。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/scheduled-task-injection.md
