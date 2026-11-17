# 邀请返水自邀循环（小号矩阵刷返水）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/invite-referral-selfdealing.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.0 ｜ CVE：-
- 产品：自研代理/邀请体系
- 影响版本：任意
- 端口：80,443
- 关键词：邀请, 返水, 代理, 自邀, 循环, 佣金
- 触发点：/api/user/invite 或注册时的 inviteCode

## 定位要点

- 组件指纹：按 `自研代理/邀请体系` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：邀请、返水、代理、自邀、循环、佣金。
- 触发条件：/api/user/invite 或注册时的 inviteCode。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/invite-referral-selfdealing.md
