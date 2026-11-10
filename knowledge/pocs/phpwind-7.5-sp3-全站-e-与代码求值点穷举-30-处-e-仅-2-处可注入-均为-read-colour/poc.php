// 全站 preg_replace 带 /e 或 /eis 的调用点 = 30 处（7.5 SP3），逐条判定输入可控性
// 格式: 文件:行号 | 替换串 | 判定
// ---------------------------------------------------------------
// 【唯一 2 处替换串含变量插值 → 唯一有注入可能的】
// require/target191.example.com:345
//   preg_replace("/<br>/eis","copyctrl('$read[colour]')",$read['content']);
//   → $read[colour] 来自 SELECT t.* FROM pw_threads（该表无 colour 列，见另一条 POC）→ 7.5 死代码
// lib/target192.example.com:409
//   $this->forum->foruminfo['copyctrl'] && preg_replace("/<br>/eis","copyctrl('$read[colour]')",$read['content']);
//   → 同上
//
// 【其余 28 处：替换串为常量或仅 \\1 捕获组引用，不可注入】
// target107.example.com:78            "/\[attachment=([0-9]+)\]/eis","upload('\\1')"
// apps/diary/example.com:168,336  "/\[s:(.+?)\]/eis",''
// target121.example.com:228,481,539,735   "/{#(.+?)}/eis",'$lang[\'\\1\']'  /  '$\\1'
// lang/up_function.php:167      "/{#(.+?)}/eis",'$\\1'
// m/wap_mod.php:124,125         "/\[code\](.+?)\[\/code\]/eis","wap_getcode('\\1')"
// target144.example.com:831                  "/<br \/>/eis","copyctrl()"            ← 常量
// simple/mod_read.php:208       "/<br>/eis","copyctrl('#FFFFFF')"      ← 常量
// target198.example.com:59                    "/\<(.+?)\>/eis",''
// target199.example.com:756               "/\[code\](.+?)\[\/code\]/eis","code_check('\\1')"
// require/app_core.php:96,113   "/\[share\]([\w\.]+),([^\[,]+?),([\d]{1,10})\[\/share\]/eis","shareEvea('\\1','\\2','\\3')"
// lib/target143.example.com:525        "/\[code\](.+?)\[\/code\]/eis","\$this->code_check('\\1')"
// require/target137.example.com:10,14,58,87,90,93,97,115,450
//   "phpcode('\\1')" / "payto('\\1')" / "SetMusic('\\1')" / "post('\\1')"
//   / "hidden('\\1','\\2')" / "sell('\\1','\\2')" / "qoute('\\1')" / "tablefun('\\2','\\3','\\4')" / "postcache('\\1')"
// require/target200.example.com:300      "code_check('\\1')"
// mode/o/require/target148.example.com:96    "shareEvea('\\1','\\2','\\3')"
// simple/example.com:71-74        "Htm_cv('\\3','<a\\1href=\"')"    ← 见下注
//
// 【simple/example.com:70-74 的次级分析（唯一有理论风险的常量替换串）】
// $output = preg_replace(
//   "/\<a(\s*[^\>]+\s*)href\=([\"|\']?)((index|cate|thread|read|faq|rss)\.php\?[^\"\'>\s]+\s?)[\"|\']?/ies",
//   "Htm_cv('\\3','<a\\1href=\"')", $output);
// \\3 的字符类排除了引号 → 不可断引；
// \\1 = (\s*[^\>]+\s*) 可含 " ，会被放进双引号串 <a\1href=" 内 → 理论上可闭合该双引号。
// 但要凑出可执行 PHP 需要页面输出里存在 "<a <可控属性含双引号> href=example.com?..." 形状，
// 且上游 Char_cv 已把请求参数里的 " 转成 &quot; → 只有 DB 中已存在未转义的 " 才可能，概率极低。
// 列为"待实测备选"，不建议优先投入额度。
//
// ---------------------------------------------------------------
// 【其它代码求值类，全部证否】
// create_function(          → 全站 0 处
// preg_replace_callback     → 全站 0 处
// assert( 带字符串          → 仅 lib/sphinx.class.php 86~1140 行，参数全是表达式（is_numeric/is_int/...），
//                             且 Sphinx 未启用（$db_sphinx 默认关）→ 死
// array_map('pwJsonEncode') → require/target105.example.com:1127、apps/stopic/action/admin/target201.example.com:38，常量回调名 → 安全
// uasort/usort              → admin/target202.example.com:314 "cmp" / admin/target203.example.com:43,66,150 "orderCmp" /
//                              require/target204.example.com:72,77 "cacheTimeCmp","orderCmp" / target205.example.com:475 "cmp" /
//                              require/profile/example.com:42 "sort_cmp" / mode/area/require/target206.example.com:13
//                              "sortTagRelate" —— 全部常量函数名 → 安全
// call_user_func_array     → api/class_base.php:122（$this->classdb[$mode] + $method）、
//                              uc_client/uc_client.php:360 —— 类方法调度，API/UC 入口在目标上已 404 → 死
// eval( / ${...} 双引号拼接   → 见下一张 POC（"写入即执行"清单）