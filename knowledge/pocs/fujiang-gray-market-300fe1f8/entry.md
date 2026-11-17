# 代理层级越权（读全站下级 / 改佣金比例）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/agent-hierarchy-privilege.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：8.8 ｜ CVE：-
- 产品：博彩包网 / 代理分销系统
- 影响版本：任意
- 端口：80,443
- 关键词：代理, agentId, 层级越权, 佣金, 返水, 会员列表
- 触发点：改 agentId 读他人下级，改 rate 提佣金

## 定位要点

- 组件指纹：按 `博彩包网 / 代理分销系统` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：代理、agentId、层级越权、佣金、返水、会员列表。
- 触发条件：改 agentId 读他人下级，改 rate 提佣金。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/agent-hierarchy-privilege.md
