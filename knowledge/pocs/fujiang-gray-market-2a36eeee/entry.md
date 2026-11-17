# 业务接口 SQL 注入（报表 / 订单 / 代理查询）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/sql-injection-business-api.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.8 ｜ CVE：-
- 产品：自研报表与查询接口
- 影响版本：任意
- 端口：80,443
- 关键词：sqli, 注入, orderNo, agentId, 报表, 盲注, union
- 触发点：orderNo=1' AND SLEEP(5)--

## 定位要点

- 组件指纹：按 `自研报表与查询接口` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：sqli、注入、orderNo、agentId、报表、盲注、union。
- 触发条件：orderNo=1' AND SLEEP(5)--。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/sql-injection-business-api.md
