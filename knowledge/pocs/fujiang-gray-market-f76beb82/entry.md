# VIP/等级接口参数篡改（免费拿高等级权益）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/vip-level-escalation.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.1 ｜ CVE：-
- 产品：自研会员体系（博彩/包网）
- 影响版本：任意
- 端口：80,443
- 关键词：vip, level, 等级, 权益, 参数篡改, 会员
- 触发点：/api/user/level 或资料更新里的 vipLevel/levelId

## 定位要点

- 组件指纹：按 `自研会员体系（博彩/包网）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：vip、level、等级、权益、参数篡改、会员。
- 触发条件：/api/user/level 或资料更新里的 vipLevel/levelId。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/vip-level-escalation.md
