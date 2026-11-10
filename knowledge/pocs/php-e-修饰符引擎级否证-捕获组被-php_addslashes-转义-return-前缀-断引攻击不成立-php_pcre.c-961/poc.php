// ==================================================================
// preg_replace /e 三类替换串的命运（php-src PHP-5.5 官方源码 + 目标实测）
// 源码已下载：runs/pw-src-732/php-src-verify/php_pcre.c (62042B) / zend_execute_API.c (51408B)
// ==================================================================

// ---- 引擎语义 ----
// ext/pcre/php_pcre.c:922-997  preg_do_eval() —— /e 的唯一执行路径
945:         if ('\\' == *walk || '$' == *walk) {          // \1 与 $1 都算反向引用
961:                         esc_match = php_addslashes(match, match_len, &esc_match_len, 0 TSRMLS_CC);   // ★ 无条件
970:                 smart_str_appendl(&code, esc_match, esc_match_len);
988:     if (zend_eval_stringl(code.c, code.len, &retval, compiled_string_description TSRMLS_CC) == FAILURE) {

// Zend/zend_execute_API.c:1139-1151
1148:     Z_STRLEN(pv) = str_len + sizeof("return ;") - 1;
1150:     memcpy(Z_STRVAL(pv), "return ", sizeof("return ") - 1);
1151:     memcpy(Z_STRVAL(pv) + sizeof("return ") - 1, str, str_len);
// → /e 对每个匹配等价于 eval("return <替换后的串>;")

// ---- 形态 A：\1 落在引号内（死）----
// 替换串： phpcode('\\1')       实际字符串： phpcode('\1')
// 捕获组含 ' → php_addslashes 变 \' → 插进去是 \' → 已转义引号 → 无法闭合 → Parse error
// 【与站点清洗无关，纯引擎行为】

// ---- 形态 B：\1 落在函数参数位/表达式位置（语法上可行）----
// 替换串： func(\1)
// 捕获组： 1),phpinfo(),strtolower(1     → eval("return func(1),phpinfo(),strtolower(1);")
//   ✗ 这是语法错误：PHP 的 return 只接受单个表达式，逗号在表达式级不是运算符
// 捕获组： 1)+(phpinfo())+(1              → eval("return func(1)+(phpinfo())+(1);")
//   ✓ 合法：单个表达式，phpinfo() 在执行过程中被调用
// 捕获组： 1)+(system)(id)+(1             → 反引号版：1)+(`id`)+(1  → 系统命令执行
// 要点：形态 B 的 payload 不含 ' " \ NUL，因此绕过 php_addslashes；
//       但必须用运算符（+ . && || 等）把所有片段连成【一个】表达式，不能出现裸分号或逗号。

// ---- 形态 C：替换串在 preg_replace 调用前已有变量插值（真正危险）----
// reg = preg_replace($p, "func('$v')", $s)      ← $v 在前一次 PHP 双引号插值时就已进入，不过 addslashes
// PHPWind 7.5 全站只有 2 处：require/target191.example.com:345 与 lib/target192.example.com:409，均为 $read['colour']

// ---- 目标实测判据（7.3.2，prosign 字段，落库回读）----
// 写入 [b]prosign_probe_7k2[/b]<\>'")
// 回读 [b]prosign_probe_7k2[/b]&lt;\&gt;&#39;&quot;)
// → < > ' " 全部实体化；【( ) 原样保留】
// 源码对应：
//   require/target105.example.com:646-666 Char_cv：str_replace(array('%3C','<'),'&lt;')、('%3E','>')→'&gt;'、
//        ('"',"'","\t",'  ')→('&quot;','&#39;','    ','&nbsp;&nbsp;')  ← 【不碰 ( )】
//   require/target105.example.com:667-679 CheckVar：str_replace(array('..',')','<','='),array('&#46;&#46;','&#41;','&#60;','&#61;'))
//        ← 确实处理 )，但——
//   target116.example.com:53-57：if (!in_array($_key,array('atc_content','atc_title','prosign','pwuser','pwpwd'))) { CheckVar(...); }
//        ← prosign 在白名单里，【完全跳过 CheckVar】
// → 实测结果 = "Char_cv 不处理括号" + "prosign 豁免 CheckVar" 两条规则的叠加，与源码完全一致，无矛盾。
// → 推论：CheckVar 的 ) 实体化只对【非白名单字段】生效；白名单字段要单独看其 InitGP 的 cvtype。