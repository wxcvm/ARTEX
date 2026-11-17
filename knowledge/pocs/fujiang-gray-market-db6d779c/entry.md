# 宝塔 / 1Panel 面板弱口令与未授权（黑灰站最常见后台）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/panel-bt-1panel.md`
- 分类：gray-market ｜ 严重度：critical ｜ CVSS：9.2 ｜ CVE：-
- 产品：宝塔面板（BT-Panel）/ 1Panel / 小皮 phpstudy
- 影响版本：任意
- 端口：8888,888,8443,2087,9999
- 关键词：宝塔, bt, 1panel, 面板, 弱口令, 未授权, 服务器
- 触发点：:8888/<随机入口> 或 :8888/login

## 定位要点

- 组件指纹：按 `宝塔面板（BT-Panel）/ 1Panel / 小皮 phpstudy` 的版本声明判断候选面（任意），端口 8888,888,8443,2087,9999 优先。
- 信号来源：宝塔、bt、1panel、面板、弱口令、未授权、服务器。
- 触发条件：:8888/<随机入口> 或 :8888/login。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/panel-bt-1panel.md
