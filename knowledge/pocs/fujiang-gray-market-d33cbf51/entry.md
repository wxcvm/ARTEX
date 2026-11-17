# 上传点 getshell（头像 / KYC 证件 / 客服附件）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/upload-getshell.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.0 ｜ CVE：-
- 产品：自研站上传点（博彩包网、代理后台高发）
- 影响版本：任意
- 端口：80,443
- 关键词：upload, getshell, 头像, kyc, 二次渲染, 双扩展名, svg xss
- 触发点：-

## 定位要点

- 组件指纹：按 `自研站上传点（博彩包网、代理后台高发）` 的版本声明判断候选面（任意），端口 80,443 优先。
- 信号来源：upload、getshell、头像、kyc、二次渲染、双扩展名、svg xss。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/upload-getshell.md
