// require/target105.example.com:605-620 —— PHPWind 的“register_globals 等价物”
function InitGP($keys,$method=null,$cvtype=1){
	!is_array($keys) && $keys = array($keys);
	foreach ($keys as $key) {
		if ($key == 'GLOBALS') continue;
		$GLOBALS[$key] = NULL;
		if ($method != 'P' && isset($_GET[$key])) {
			$GLOBALS[$key] = $_GET[$key];          // ← 请求参数进裸全局变量
		} elseif ($method != 'G' && isset($_POST[$key])) {
			$GLOBALS[$key] = $_POST[$key];
		}
		if (isset($GLOBALS[$key]) && !empty($cvtype) || $cvtype==2) {
			$GLOBALS[$key] = Char_cv($GLOBALS[$key],$cvtype==2,true);   // 默认 cvtype=1 → Char_cv
		}
	}
}

// target116.example.com:226-241 —— 唯一装载 'a' 的地方，且立刻 unset
if ($db_bbsifopen==2 && SCR!='login' && !defined('CK')) { ... }
if ($db_ifsafecv && !$safecv && !defined('PRO') && strpos($db_safegroup,",$groupid,") !== false ) {
	Showmsg('safecv_prompt');
}
232: if ($db_ads && !$windid && (is_numeric($_GET['u']) || ($_GET['a'] && strlen($_GET['a'])<16)) && strpos($pwServer['HTTP_REFERER'],$pwServer['HTTP_HOST'])===false) {
233:     InitGP(array('u','a'));
234:     Cookie('userads',"$u\t$a\t".md5($pwServer['HTTP_REFERER']));
235: } elseif ( $db_ads=='1' && ($cookie_userads = GetCookie('userads')) ) {
236:     list($u,$a) = explode("\t",$cookie_userads);
237:     if ((int)$u>0 || ($a && strlen($a)<16)) {
238:         require_once(R_P.'require/userads.php');
239:     }
240: }
241: unset($u,$a,$cookie_userads);   // ← 无论走哪个分支都执行，$a 被删

// 其它 InitGP(array('a',...)) 的位置（均属 apps/ 路由，不与 bank_sort.php 同作用域）：
// apps/ajax/example.com:6  apps/article/example.com:12  apps/photos/example.com:14
// apps/diary/example.com:14  apps/groups/action/m_groups.php:17  m_group.php:9
// apps/groups/action/m_galbum.php:12  apps/share/action/m_share.php:11
// 全站没有任何 InitGP 装载 '_'  → ${@eval($_)} 必死