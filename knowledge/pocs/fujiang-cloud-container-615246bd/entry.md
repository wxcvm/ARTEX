# Docker Remote API 未授权（直接拿宿主）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/cloud-container/docker-api-2375.md`
- 分类：cloud-container ｜ 严重度：critical ｜ CVSS：- ｜ CVE：-
- 产品：Docker Engine
- 影响版本：任意（2375 裸露）
- 端口：2375,2376
- 关键词：docker, api, 2375, privileged, escape
- 触发点：-

## 定位要点

- 组件指纹：按 `Docker Engine` 的版本声明判断候选面（任意（2375 裸露）），端口 2375,2376 优先。
- 信号来源：docker、api、2375、privileged、escape。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/cloud-container/docker-api-2375.md
