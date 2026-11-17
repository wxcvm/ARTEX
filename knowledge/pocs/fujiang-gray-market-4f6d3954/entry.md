# 客服 IM / WebSocket 未鉴权（订阅他人会话）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/im-websocket-unauth.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.2 ｜ CVE：-
- 产品：自研客服系统 / socket.io / 原生 WS
- 影响版本：任意
- 端口：80,443
- 关键词：websocket, socket.io, im, 客服, 未鉴权, room, 订阅
- 触发点：ws://host/ws 无 token 直接订阅 room

## 定位要点

- 组件指纹：按 `自研客服系统 / socket.io / 原生 WS` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：websocket、socket.io、im、客服、未鉴权、room、订阅。
- 触发条件：ws://host/ws 无 token 直接订阅 room。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/im-websocket-unauth.md
