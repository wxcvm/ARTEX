# 签到 / 任务奖励批量刷取（缺幂等与频率控制）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/daily-reward-farming.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：7.8 ｜ CVE：-
- 产品：自研活动系统（博彩/包网）
- 影响版本：任意
- 端口：80,443
- 关键词：签到, 任务, 奖励, 刷取, 幂等, 并发
- 触发点：/api/activity/sign、/api/task/claim

## 定位要点

- 组件指纹：按 `自研活动系统（博彩/包网）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：签到、任务、奖励、刷取、幂等、并发。
- 触发条件：/api/activity/sign、/api/task/claim。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/daily-reward-farming.md
