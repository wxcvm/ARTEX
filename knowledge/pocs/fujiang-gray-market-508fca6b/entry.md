# 导出 / 下载接口任意文件读取

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/export-file-read.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.6 ｜ CVE：-
- 产品：对账单 / 报表导出 / 凭证下载
- 影响版本：任意
- 端口：80,443
- 关键词：任意文件读, download, export, path traversal, 对账单
- 触发点：/api/report/export?file=../../../../etc/passwd

## 定位要点

- 组件指纹：按 `对账单 / 报表导出 / 凭证下载` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：任意文件读、download、export、path traversal、对账单。
- 触发条件：/api/report/export?file=../../../../etc/passwd。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/cf086fc329011c6f8eeb3b67f0adbcc3f005debb/plugin/knowledge/gray-market/export-file-read.md
