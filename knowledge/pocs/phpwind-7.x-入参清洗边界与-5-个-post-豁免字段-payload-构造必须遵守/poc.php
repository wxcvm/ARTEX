// target116.example.com:53-61  —— 所有前台请求必经
foreach ($_POST as $_key => $_value) {
	if (!in_array($_key,array('atc_content','atc_title','prosign','pwuser','pwpwd'))) {
		CheckVar($_POST[$_key]);
	}
}
foreach ($_GET as $_key => $_value) {
	CheckVar($_GET[$_key]);        // GET 无豁免，全量
}

// require/target105.example.com:667-679
function CheckVar(&$var) {
	if (is_array($var)) {
		foreach ($var as $key => $value) { CheckVar($var[$key]); }
	} elseif (P_W != 'admincp') {
		$var = str_replace(array('..',')','<','='),array('&#46;&#46;','&#41;','&#60;','&#61;'),$var);
	} elseif (str_replace(array('<iframe','<meta','<script'),'',$var)!=$var) {
		global $basename; $basename = 'javascript:history.go(-1);'; adminmsg('word_error');
	}
}
// 注意：只转义右括号 ) 与 = < ..，左括号 ( 保留。

// require/target105.example.com:646-666  Char_cv（InitGP 默认 cvtype=1 会调用）
'"'    → &quot;
"'"    → &#39;
"<"    → &lt;      ("%3C" 同样)
">"    → &gt;      ("%3E" 同样)
"\t"   → 4 空格 ; 连续两个空格 → &nbsp;&nbsp;
去掉 \0 %00 \r 与控制字符；$istrim=true 时先 trim
// 保留：$ { } ( [ ] ; . / : _ @ 以及全部 ASCII 字母数字

// target116.example.com:40-51 URL 重写模式（$_GET = $_NGET）发生在 CheckVar 之前，
// 所以重写出来的参数同样会被 CheckVar 处理，无法借此绕过。