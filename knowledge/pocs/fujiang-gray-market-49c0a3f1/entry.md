# 客服聊天记录导出未鉴权（全站会话与身份材料）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/cs-chat-export-unauth.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.4 ｜ CVE：-
- 产品：自研客服系统（在线客服 / 工单 / IM）
- 影响版本：任意
- 端口：80,443
- 关键词：客服, 聊天记录, 导出, 未鉴权, 会话, im
- 触发点：/api/cs/history/export  /api/im/sessions

## 定位要点

- 组件指纹：按 `自研客服系统（在线客服 / 工单 / IM）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：客服、聊天记录、导出、未鉴权、会话、im。
- 触发条件：/api/cs/history/export  /api/im/sessions。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/cs-chat-export-unauth.md
