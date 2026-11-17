# 活动 / 彩金 / 返水无限领取（幂等缺失）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/promotion-abuse.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：7.9 ｜ CVE：-
- 产品：促销活动与奖励接口
- 影响版本：任意
- 端口：80,443
- 关键词：promotion, bonus, 彩金, 返水, 幂等, 并发领取, 活动滥用
- 触发点：/api/promotion/receive 重复或并发领取

## 定位要点

- 组件指纹：按 `促销活动与奖励接口` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：promotion、bonus、彩金、返水、幂等、并发领取、活动滥用。
- 触发条件：/api/promotion/receive 重复或并发领取。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/promotion-abuse.md
