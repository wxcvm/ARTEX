# Nginx alias 路径穿越（读任意文件/源码泄露）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/web-cms/nginx-alias-traversal.md`
- 分类：web-cms ｜ 严重度：high ｜ CVSS：- ｜ CVE：-
- 产品：Nginx（location alias 缺尾斜杠）
- 影响版本：任意
- 端口：80,443
- 关键词：nginx, alias, 穿越, 目录, 源码, 配置
- 触发点：-

## 定位要点

- 组件指纹：按 `Nginx（location alias 缺尾斜杠）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：nginx、alias、穿越、目录、源码、配置。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/web-cms/nginx-alias-traversal.md
