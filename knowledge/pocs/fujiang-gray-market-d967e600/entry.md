# 注单 / 流水接口越权批量拉取（全站会员数据）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/bet-history-idor.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.6 ｜ CVE：-
- 产品：自研注单与报表接口（博彩）
- 影响版本：任意
- 端口：80,443
- 关键词：注单, 流水, idor, 越权, 批量, 报表
- 触发点：/api/bet/list?uid=  /api/report/bet

## 定位要点

- 组件指纹：按 `自研注单与报表接口（博彩）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：注单、流水、idor、越权、批量、报表。
- 触发条件：/api/bet/list?uid=  /api/report/bet。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/bet-history-idor.md
