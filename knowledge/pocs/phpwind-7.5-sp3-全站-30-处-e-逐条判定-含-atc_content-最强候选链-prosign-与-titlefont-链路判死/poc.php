// ============================================================================
// 前提（见另一条 POC）：/e 的捕获组被 php_addslashes 转义 + 求值带 return 前缀
//   → "常量替换串 + 捕获组" 形态在引擎层就不可能断引。
// 因此下表"可否断引"一列在引擎层已被封死；此处仍逐条给出【来源 / 清洗链 / 入口 CheckVar】，
// 用于确认没有遗漏"替换串里有调用前变量插值"的点，以及验证判断是否与清洗无关。
//
// 入口 CheckVar 判定：target116.example.com:53-61（POST 逐键，5 个白名单字段豁免；GET 全量）
//   · 全部根目录前台文件（message.php / read.php / post.php / example.com / apps.php /
//     mode.php / example.com / simple/example.com / m/* 等）第一件事都是 require_once('global.php')
//     → CheckVar 生效
//   · example.com → admin/example.com  【无 CheckVar 循环，已逐行确证】→ 后台路径不过 CheckVar（但进不去后台）
//   · install.php → 只 require require/common.php，不经过 global.php → 不过 CheckVar（有 install.lock 保护）
//
// ---------------------------------------------------------------------------
// 【A. 唯一"入库时 ' 与 ) 都能原样保留"的输入 —— 发帖正文 atc_content】
// 链： $_POST['atc_content']
//   target116.example.com:53-57  CheckVar → 【白名单豁免，) 与 = 与 < 都保留】
//   require/target214.example.com:94   InitGP(array('atc_title','atc_content'), 'P', 0)  ← cvtype=0，【不过 Char_cv】
//   require/target215.example.com:143 / target216.example.com:315  同样 cvtype=0
//   example.com:123 InitGP(array('atc_title','atc_content'),'P',0)  ← 同上
//   pw_example.com:313 InitGP(array('atc_content'),'P',0)          ← 同上
//   apps/groups/action/m_group.php:814 InitGP(array('atc_title','atc_content'),'P',0) ← 同上
//   （对照：target199.example.com:377 InitGP(array('msg_title','atc_content',...),'P') cvtype=1 → Char_cv 实体化）
//   → 入库 pw_tmsgs.content 为【原文】，' 与 ) 都在
//   → 读帖：target144.example.com:774 / require/target191.example.com:309 / lib/target192.example.com:368
//           / simple/mod_read.php:206 / mode/area/area_thread.php:112
//           全部 convert($read['content'], $db_windpost)
//   → convert() 首段（require/bbscode.php:9-11）：
//        if (strpos($message,"[code]")!==false && strpos($message,"[/code]")!==false)
//            $message = preg_replace("/\[code\](.+?)\[\/code\]/eis","phpcode('\\1')",$message,$db_cvtimes);
//   结论：这是全站【最强候选】—— 输入侧 ' 和 ) 都干净、捕获组 (.+?) 不排除引号、
//         替换串 '\1' 落在引号内、phpcode() 是普通函数。
//         理论上 payload = [code]'.`id`.'[/code]  → phpcode('' . `id` . '')
//         但【php_pcre.c:961 的 php_addslashes 把 ' 变成 \'】→ 无法闭合 → 死。
//   ★ 这一条的判死【不由 PHPWind 的清洗决定】，纯引擎行为。
//
// ---------------------------------------------------------------------------
// 【B. 逐点判定表（30 处）】
// 格式：文件:行号 | 捕获组来源 | 入库清洗 | 引擎层 | 结论
//
// 1) require/target137.example.com:10  | 帖子/短消息/签名/日记正文／DB 或转换参数 $message | 见 A / 签名见 D | 被 addslashes 堵 | 死
// 2) require/target137.example.com:14  | $message 同上（[payto]）                        | 同上 | 同上 | 死
// 3) require/target137.example.com:58  | $message 同上（[music]）                        | 同上 | 同上 | 死
// 4) require/target137.example.com:87  | $message（[post] 权限隐藏）                     | 同上 | 同上 | 死
// 5) require/target137.example.com:90  | $message（[hide=]）                             | 同上 | 同上 | 死
// 6) require/target137.example.com:93  | $message（[sell=]）                             | 同上 | 同上 | 死
// 7) require/target137.example.com:97  | $message（[quote]）                             | 同上 | 同上 | 死
// 8) require/target137.example.com:115 | $message（[table]）                             | 同上 | 同上 | 死
// 9) require/target137.example.com:450 | $message（[s:表情]）                            | 同上 | 同上 | 死
// 10) require/target137.example.com:33 | "/\[size=(\d+)\]/eis" → 捕获组是【纯数字】      | — | — | 死（字符类即数字）
// 11) require/target137.example.com:53/55 | [img] 捕获组 [^\<\r\n\"']+? 【显式排除 " 与 '】| — | — | 死（字符类排除引号）
// 12) require/target137.example.com:121-128 | [wmv]/[rm] 捕获组 [^\<\r\n\"']+? 同样排除引号 | — | — | 死
// 13) require/target137.example.com:142-144 | [iframe] 捕获组排除引号 | — | — | 死
// 14) target199.example.com:756  | 短消息正文 $message（DB）                    | 写入时 target199.example.com:377 InitGP(...,'P') cvtype=1 → Char_cv 实体化 | 双重 | 死
// 15) m/wap_mod.php:124/125 | WAP 帖子正文 $string（DB）              | m/target132.example.com:58 InitGP(array('subject','content'),'P',0) cvtype=0 → 不过 Char_cv | addslashes | 死（引擎层）
// 16) apps/diary/example.com:168/336 | 日记正文 $rt['content']（DB）   | apps/diary/example.com:454/578 直接读 $_POST['atc_content']（CheckVar 豁免），随后同文件无 cvtype=1 覆盖 | addslashes | 死（引擎层）
// 17) lib/target143.example.com:525 | 发帖/编辑时对 $message 的转换（'$this->code_check(\'\\1\')'）| 同 A（cvtype=0） | addslashes | 死
// 18) require/target200.example.com:300 | 同上 | 同 A | addslashes | 死
// 19) target107.example.com:78 | "/\[attachment=([0-9]+)\]/eis" → 捕获组【纯数字】 | — | — | 死
// 20) require/app_core.php:96/113 | feed 动态 $message | 写入路径为 pwAddFeed（内容来自已入库字段） | addslashes | 死
// 21) mode/o/require/target148.example.com:96 | feed $message；\1=[\w\.] \2=[^\[,]+? \3=[\d]{1,10} → \1/\3 字符类已限死 | — | — | 死（字符类+引擎）
// 22) simple/example.com:71-74 | 页面输出 $output；\3 字符类 [^"'>\s]+ 显式排除引号 | — | — | 死（字符类+引擎）
// 23) mode/o/require/core.php（同 21，parseFeed 唯一一处）| — | — | — | 死
// 24) target121.example.com:228/481/539/735 | 语言包字符串（服务器本地文件）| — | — | 死
// 25) lang/up_function.php:167 | 同上 | — | — | 死
// 26) apps/diary/example.com:162/329 | $wordsfb->convert() 是 FilterUtil 的敏感词方法，非 /e | — | — | 死
// 27) require/target137.example.com:16/22/24/27-47 | 无 /e 修饰符（普通 preg_replace / str_replace）| — | — | 死
// 28) require/target137.example.com:50  | $db_windcode['replacearray'] 来自【后台配置】的 bbcode 规则表 | 后台 | — | 需后台（进不去）；且捕获组同样过 addslashes
// 29) require/target191.example.com:345 | copyctrl('$read[colour]') ← 【替换串里有调用前变量插值 ★唯一形态】| 见 colour 判死 POC | 无需 addslashes 就已注入 | 死（无 colour 列 + CheckVar/Char_cv）
// 30) lib/target192.example.com:409 | 同 29 | 同上 | — | 死
//
// ---------------------------------------------------------------------------
// 【C. 结论】
// 全站 30 处 /e：
//   · 28 处属"常量替换串 + 捕获组"或"捕获组字符类已排除引号" → 引擎层 php_addslashes 封死 → 【死】
//   · 2 处属"替换串含调用前变量插值"（唯一真危险形态）→ 输入变量 $read['colour'] 不存在（7.5 无该列）
//     且即使存在，其写入路径的 ' 与 ) 分别被 Char_cv / CheckVar 实体化 → 【死】
//   · 无一处可断引 → 无 payload，无权限需求可言
//
// 【D. prosign（签名）完整链路 —— 判死】
// 写入：require/profile/target217.example.com:218  InitGP(array('proicon','prosign','showsign',...),'P')
//        ← cvtype 默认 1 → Char_cv（require/target105.example.com:646-666）：' → &#39;  " → &quot;  < → &lt;
//      :248  $_G['signnum'] && strlen($prosign) > $_G['signnum'] && Showmsg('sign_limit')  ← 长度上限（用户组配置）
//      :377  $cksign = convert($prosign,$db_windpic,2)   ← ★ 入库前就调 convert（会走 target137.example.com:10 那一行）
//      :389  'signature' => $prosign  → 入库 pw_members.signature
//      注册侧同款：target133.example.com:342 InitGP(array(...,'regsign',...),'P') cvtype=1 → 同样 Char_cv
// 渲染：target144.example.com:671、require/target191.example.com:246、lib/target192.example.com:307、target110.example.com:246、
//       admin/target218.example.com:433、target133.example.com:356 —— 全部 convert($signature,$db_windpic,2)
// 判定：CheckVar 白名单让 ' 与 ) 能过，但 InitGP 默认 cvtype=1 的 Char_cv 在【入库前】就把 ' 实体化
//       → DB 里没有可断引的 ' → 死。且即便没有 Char_cv，引擎层 addslashes 仍会挡。
//       （注意 :377 的位置很有意思：改签名这个动作本身就会触发一次 convert，但被同样两道墙挡住。）
//
// 【E. hack/toolcenter/require/colortitle.php + titlefont 链路 —— 判死】
// 写入 hack/toolcenter/require/target208.example.com:22-24（前台道具"彩色标题"）
//      :22 InitGP(array('title1','title2','title3','title4','title5','title6'))  ← cvtype 默认 1 → Char_cv
//      :23 $titlefont = "$title1~$title2~$title3~$title4~$title5~$title6~"      ← 无颜色白名单、无字符过滤
//      :24 $db->update("UPDATE pw_threads SET titlefont=".pwEscape($titlefont).",toolinfo=..." WHERE tid=...)
//      对照 target108.example.com:1587-1607（批量高亮）:1601 有 if ($title1 && !preg_match('/#[0-9A-F]{6}/is',$title1)) 颜色白名单 + :1606 Char_cv
// 读取（全站 titlefont 95 处命中已核）：target111.example.com:530、example.com:318、target194.example.com:404、example.com:1062、
//      lib/target209.example.com:214、require/target210.example.com:191、require/target211.example.com:72、
//      apps/article/example.com:132/224、apps/share/action/m_share.php:81、target110.example.com:293/333、target108.example.com:1558
//      —— 全部 explode("~",$titlefont) 取 [0] 当颜色，其中 target108.example.com:1563-1565 还有 #RRGGBB 白名单
// 判定：
//   · 【无一处进 /e、无一处进 eval、无一处进模板编译】—— 全站逐处核对
//   · 输出目标是 HTML 属性 style="color:xxx"，而写入时 Char_cv 已把 " ' < > 实体化 → 连 XSS 都难，不是 RCE
//   · 结论：死。7.3.2 即便同时存在 colour 与 titlefont 两个字段、各有读写，也改变不了：
//     colour 那两处 /e 的输入要么字段不存在、要么被 Char_cv 实体化；且引擎层 addslashes 对所有 /e 生效。
//     置信度：colour 列存在 = 中高（外推）；可断引 = 不可能（引擎级确证）。