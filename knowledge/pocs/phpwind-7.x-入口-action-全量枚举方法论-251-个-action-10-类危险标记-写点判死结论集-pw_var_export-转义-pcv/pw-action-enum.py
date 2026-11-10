【增补 2026-10 · 写点汇点全量分类（不按 action 分派，按调用点扫）】
全源码扫 1830 处写/删/包含汇点（writeover/file_put_contents/fwrite/fopen写模式/unlink/P_unlink/rmdir/mkdir/chmod/move_uploaded_file/include/require）：类A（路径含请求变量）39 处、类B（内容含请求变量）18 处、类C（参数是文件名）0 处、类D 1773 处。

★★ 核心判据（第三条防线，此前未记录）：**`P_unlink` 也是 `@unlink(Pcv($filename))`**（target105.example.com:209-211）。加上 `writeover`(:231) 与 `readover`(:221) 内部的 `Pcv`，**PHPWind 三道文件操作统一防线全部内建 `Pcv`（拦 `://`/NUL/`..`）** → 类A 39 处路径穿越**全部判死**，连「任意文件删除」都不存在。这解释了为何历史上所有路径穿越尝试（setstyles 的 $sid、bakup 的 $pre、sort.php 的 $month、thread.php 的 $fid/$page、level.php 的 $delid/$gid、guestfunc.php 的 $page）全部失败。

★ 唯一例外：**`hack/bank/example.com` L362-373 `savearray()` 不走 `pw_var_export`** —— 它用 `'"'.addslashes($value2).'"'` 拼数组，而 `addslashes` **不转义 `$ { } [ ] ( )`** → 双引号字符串里的 `${...}` 会被 PHP 求值。写入 `data/bbscache/bank_sort.php` 后由 L76 `include` 同作用域执行 → 构成代码执行链（数据源为 `pw_members.username`，需先污染 DB）。

【增补 · 字符过滤规则速查】
- `Char_cv`(target105.example.com:646-665)：实体化 `" ' < >`，去 NUL/`%00`/`\r`，**不过滤 `$ { } [ ] ( )`**
- `CheckVar`(target105.example.com:667-679)：把 `.. ) < =` 实体化 —— **但对 `$_POST['pwuser']`/`pwpwd`/`atc_content`/`atc_title`/`prosign` 豁免**（target116.example.com:53-57 白名单）
- `Sql_cv`(target105.example.com:441-448)：**不做任何过滤**，只是 INSERT 到 pw_sqlcv 再 SELECT 回来的往返测试
- `Pcv`(target105.example.com:372-380)：拦 `://`、NUL、`..`（默认 $ifcheck=1）；**不拦 `/`**
- `changename.php` 的用户名黑名单：`& ' " / * , < > \r \t \n #` —— **不含 `$ { } [ ] ( )`**；长度上限 `$rg_regmaxname` 来自后台配置 `data/bbscache/dbreg.php` 的 `$rg_namelen`

【增补 · data/ 目录可读性判据】
`.php` 缓存文件一律被 nginx 交给 PHP 执行 → **200/0B（存在但无输出，无源码泄露）**；而 **`.htm`/`.js`/`.css` 不经 PHP handler → 内容原文返回**。实测 `data/style/wind_css.htm` = 14279B 可读（含 `<?php` 片段，`<!--css-->` 分隔符 2 次），`data/bbscache/face.js` = 639B 可读，`data/tmp/` = 403。**注意**：`wind_css.htm` 虽是 `updatecache_sy` 的 `eval("\$style_css = \"$style_css\";")` 输入源，但全文 `${...}` 结构 **0 命中**（只有 `$linkcolor`/`$tplpath` 等普通 `$var`）→ 只做变量插值，**eval 链判死**。