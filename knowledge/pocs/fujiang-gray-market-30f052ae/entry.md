# 客服工单附件越权下载（证件照 / KYC 材料）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/support-ticket-attachment-idor.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.2 ｜ CVE：-
- 产品：自研客服工单系统
- 影响版本：任意
- 端口：80,443
- 关键词：工单, 附件, 越权, 证件, kyc, 下载
- 触发点：/api/support/attachment?fileId=

## 定位要点

- 组件指纹：按 `自研客服工单系统` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：工单、附件、越权、证件、kyc、下载。
- 触发条件：/api/support/attachment?fileId=。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/support-ticket-attachment-idor.md
