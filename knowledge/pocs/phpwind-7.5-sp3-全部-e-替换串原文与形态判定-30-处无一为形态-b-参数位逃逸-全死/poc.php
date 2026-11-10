// ============================================================================
// 目的：回答"哪一处 /e 的 \1 落在表达式位置（形态 B，可活）"
// 方法：把每处的正则原文与替换串原文一字不改贴出，判定 \1/\2/\3 的落点
// 形态定义（见另一条 POC）：
//   形态 A：\N 落在引号内        → php_addslashes 把 ' 变 \' → 不能断引 → 死
//   形态 B：\N 落在参数位/表达式 → payload 无引号可绕 addslashes → 需用运算符串成单表达式
//   形态 C：替换串在调用前已有变量插值（非 backref）→ 不受 addslashes 影响
// ============================================================================
//
// 【1】require/target137.example.com:10  convert() 首段 [code] 标签
// 正则  "/\[code\](.+?)\[\/code\]/eis"
// 替换  "phpcode('\\1')"          → 实际串 phpcode('\1')
// 判定  \1 在【单引号内】→ 形态 A → 死
//
// 【2】require/target137.example.com:14  [payto]
// 正则  "/\[payto\](.+?)\[\/payto\]/eis"
// 替换  "payto('\\1')"            → payto('\1')
// 判定  形态 A → 死
//
// 【3】require/target137.example.com:58  [music=n]...[/music]
// 正则  "/\[music=(\d+)\](.+?)\[\/music\]/eis"
// 替换  "SetMusic('\\1')"         → SetMusic('\1')     且 \1 是 (\d+) 纯数字
// 判定  形态 A + 数字字符类 → 死（双重）
//
// 【4】require/target137.example.com:87  [post]
// 正则  "/\[post\](.+?)\[\/post\]/eis"
// 替换  "post('\\1')"             → post('\1')
// 判定  形态 A → 死
//
// 【5】require/target137.example.com:90  [hide=x]
// 正则  "/\[hide=(.+?)\](.+?)\[\/hide\]/eis"
// 替换  "hidden('\\1','\\2')"     → hidden('\1','\2')
// 判定  形态 A → 死
//
// 【6】require/target137.example.com:93  [sell=x]
// 正则  "/\[sell=(.+?)\](.+?)\[\/sell\]/eis"
// 替换  "sell('\\1','\\2')"       → sell('\1','\2')
// 判定  形态 A → 死
//
// 【7】require/target137.example.com:97  [quote]
// 正则  "/\[quote\](.+?)\[\/quote\]/eis"
// 替换  "qoute('\\1')"            → qoute('\1')
// 判定  形态 A → 死
//
// 【8】require/target137.example.com:115  [table]
// 正则  '/\[table(=(\d{1,3}(%|px)?))?\](.*?)\[\/table\]/eis'
// 替换  "tablefun('\\2','\\3','\\4')"   → 三个 backref 全在单引号内
// 判定  形态 A → 死（另 \2/\3 字符类为数字/px）
//
// 【9】require/target137.example.com:450  [s:表情]
// 正则  "/\[s:(.+?)\]/eis"
// 替换  "postcache('\\1')"        → postcache('\1')
// 判定  形态 A → 死
//
// 【10】require/target137.example.com:33  [size=n]
// 正则  "/\[size=(\d+)\]/eis"
// 替换  "size('\\1','$allow[size]')"    ← 注意 $allow[size] 是【调用前插值】（形态 C 特征），
//        但它来自 $allow 数组（$db_windpost 配置），非用户可控；且 \1 是 (\d+)
// 判定  形态 A + 数字字符类 + 插值项不可控 → 死
//
// 【11】require/target137.example.com:53 / 55  [img]
// 正则  "/\[img\]([^\<\r\n\"']+?)\[\/img\]/eis"       ← 字符类【显式排除 " 与 '】
// 替换  允许图片： "cvpic('\\1','','$allow[picwidth]','$allow[picheight]')"
//        禁止图片： "nopic('\\1')"
// 判定  字符类排除引号 + 形态 A → 死（双重）
//
// 【12】require/target137.example.com:121-130  [wmv] / [rm]（$allow['mpeg'] 分支）
// 正则  "/\[wmv=(0|1)\]([^\<\r\n\"']+?)\[\/wmv\]/eis"
//       "/\[wmv(=([0-9]{1,3})\,([0-9]{1,3})\,(0|1))?\]([^\<\r\n\"']+?)\[\/wmv\]/eis"
//       "/\[rm(=([0-9]{1,3})\,([0-9]{1,3})\,(0|1))?\]([^\<\r\n\"']+?)\[\/rm\]/eis"
// 替换  "wplayer('\\2','314','53','\\1','wmv')" / "wplayer('\\5','\\2','\\3','\\4','wmv')" / "wplayer('\\5','\\2','\\3','\\4','rm')"
// 判定  字符类排除引号 + 形态 A → 死
//
// 【13】require/target137.example.com:142 / 144  [iframe]
// 正则  "/\[iframe\]([^\[\<\r\n\"']+?)\[\/iframe\]/is"   ← 无 /e（第 142/144 行是 /is）
// 判定  无 /e → 不适用
//
// 【14】target199.example.com:756  autourl() 的 [code]（短消息正文）
// 正则  "/\[code\](.+?)\[\/code\]/eis"
// 替换  "code_check('\\1')"       → code_check('\1')
// 判定  形态 A → 死
//       另外：写入侧 target199.example.com:377 InitGP(array('msg_title','atc_content','pwuser','ifsave'),'P') 是 cvtype=1
//             → Char_cv 在入库前就把 ' 实体化（与 prosign 实测一致）
//
// 【15】m/wap_mod.php:124 / 125  WAP 版 [code] / [quote]
// 正则  "/\[code\](.+?)\[\/code\]/eis"  /  "/\[quote\](.+?)\[\/quote\]/eis"
// 替换  "wap_getcode('\\1')"  /  "wap_getquote('\\1')"
// 判定  形态 A → 死
//       （WAP 写入侧 m/target132.example.com:58 InitGP(array('subject','content'),'P',0) 是 cvtype=0，不过 Char_cv，
//         但引擎层 addslashes 一样挡；且替换串是形态 A）
//
// 【16】apps/diary/example.com:168 / 336  日记 [s:表情]
// 正则  "/\[s:(.+?)\]/eis"
// 替换  ''                        ← 替换串是【空字符串】，没有任何 backref
// 判定  死（无 \1 可注入）
//
// 【17】lib/target143.example.com:525  发帖/编辑时 [code]
// 正则  "/\[code\](.+?)\[\/code\]/eis"
// 替换  "\$this->code_check('\\1')"   → $this->code_check('\1')
// 判定  形态 A → 死
//
// 【18】require/target200.example.com:300  [code]
// 正则  "/\[code\](.+?)\[\/code\]/eis"
// 替换  "code_check('\\1')"
// 判定  形态 A → 死
//
// 【19】target107.example.com:78  [attachment=n]
// 正则  "/\[attachment=([0-9]+)\]/eis"
// 替换  "upload('\\1')"          且 \1 是 ([0-9]+) 纯数字
// 判定  形态 A + 数字字符类 → 死
//
// 【20】require/app_core.php:96 / 113  [share=...][/share]
// 正则  "/\[share\]([\w\.]+),([^\[,]+?),([\d]{1,10})\[\/share\]/eis"
// 替换  :96  "shareEvea('\\1','\\2','\\3')"      :113  ""（空串）
// 判定  :96 形态 A → 死（且 \1=[\w\.]、\3 纯数字）；:113 无 backref → 死
//
// 【21】mode/o/require/target148.example.com:96  parseFeed 的 [share]
// 正则与替换同 【20】:96
// 判定  形态 A → 死
//
// 【22】simple/example.com:71-74  页面输出内的链接静态化
// 正则  "/\<a(\s*[^\>]+\s*)href\=([\"|\']?)((index|cate|thread|read|faq|rss)\.php\?[^\"\'>\s]+\s?)[\"|\']?/ies"
// 替换  "Htm_cv('\\3','<a\\1href=\"')"      → 实际串 Htm_cv('\3','<a\1href="')
// 判定  \3 落在单引号内（形态 A）；\1 也落在第二个参数的单引号内（形态 A）
//       ＋ \3 的字符类 [^\"\'>\s]+ 排除 " 与 '     → 死（双重）
//
// 【23】require/target191.example.com:345 / lib/target192.example.com:409  ← 唯一形态 C
// 替换  "copyctrl('$read[colour]')"   ← $read[colour] 由 PHP 双引号【在调用前】插值，非 backref
// 判定  形态 C（真危险形态），但 7.5 的 pw_threads 无 colour 列（建表语句 lang/install_wind.sql:3383-3423 逐列核对）
//       → 变量恒空 → 死。见 colour 判死 POC。
//
// 【24】target121.example.com:228/481/539/735、lang/up_function.php:167
// 替换  '$lang[\'\\1\']' / '$\\1'  ← 来源是本地语言包文件，非用户输入
// 判定  死
//
// ============================================================================
// 【总结论】
//   · 30 处里：形态 C 1 处（target191.example.com:345 + target219.example.com:409，同一变量算 1 组）→ 变量不存在 → 死
//              形态 A 28 组 → php_addslashes 引擎层挡 → 死
//              形态 B 0 处 → 【PHPWind 7.5 SP3 里不存在任何"\\N 落在表达式位置"的 /e 点】
//   · 形态 B 在语法上是成立的（用 1)+(phpinfo())+(1 绕开 return 的单表达式限制），
//     只是 PHPWind 的所有 bbcode 替换串都写成 func('\\1') 这种"把捕获组当字符串参数"的自然形态，
//     一处都没写成 func(\\1)。这是本 CMS 的一个偶然性保护，不是设计。
//   · 因此：/e 方向全死。可以收敛。
// ============================================================================