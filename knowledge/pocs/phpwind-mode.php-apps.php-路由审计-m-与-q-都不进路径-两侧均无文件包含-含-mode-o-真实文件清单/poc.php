// ================= mode.php（全 76 行）=================
2: define('SCR','mode');
3: require_once('global.php');
5: $m = GetGP('m');
7: selectMode($m);                       // ★ $m 只用于白名单校验，不进路径
17: InitGP(array('q','space'));
36: if (strpos($q,'..') !== false) Showmsg('undefined_action');
42: if ($m == 'o' && $q && !in_array($q, array('user','friend','browse','invite','board','myapp','home','app'))) {
43:     require_once 'apps.php';           // ★ 非内置 q 全部转交 apps.php
44: } else {
46:     empty($q) && $q = 'home';
58:     if (file_exists(M_P . "m_{$q}.php")) {
59:         @include_once Pcv(D_P.'data/bbscache/'.$db_mode.'_config.php');
67:         if (file_exists(M_P.'require/core.php')) { require_once(M_P.'require/core.php'); }
70:         require_once Pcv(M_P . "m_{$q}.php");     // ← $q 拼文件名，但过 Pcv
72:     } else { Showmsg('undefined_action'); }
}

// target116.example.com:300-317   ★ 关键：$m 与 M_P 的关系
function selectMode(&$m) {
	global $defaultMode,$db_mode,$db_modes,$db_modepages,$pwServer,$db_modedomain;
	if (defined('M_P')) return;
	if (in_array(SCR, array('index','cate','mode'))) {
		$db_mode = $defaultMode;
		if (!$m && $db_modedomain) { $m = array_search($pwServer['HTTP_HOST'],$db_modedomain); }
		if ($db_modes && isset($db_modes[$m]) && is_array($db_modes[$m]) && $db_modes[$m]['ifopen']) {
			$db_mode = $m;                                   // ★ 必须在 $db_modes 白名单里
		}
		if (!empty($db_mode) && $db_mode!='bbs' && file_exists(R_P."mode/$db_mode/")) {
			define('M_P',R_P."mode/$db_mode/");               // ★ 用的是 $db_mode，且再走 file_exists
			$db_modepages = $db_modepages[$db_mode];
			$GLOBALS['pwModeImg'] = "mode/$db_mode/images";
		}
	}
}
// → $m 想变成路径必须同时满足：isset($db_modes[$m]) 为真 + is_array + ['ifopen'] 非空 + file_exists
// → ?m=../../xxx 的 $db_modes['../../xxx'] 不存在 → $db_mode 保持 $defaultMode('bbs') → M_P 不定义或 =mode/bbs/
// → mode.php 侧【无文件包含】

// ================= apps.php（前 60 行）=================
2: define('SCR','app');
3: require_once('global.php');
7: InitGP(array('q'));
9: if (in_array($q,array('ajax','article','diary','galbum','group','groups','hot','photos','sharelink','share','stopic','write'))) {
11:     require_once R_P . 'require/app_route.php';
}

// ================= require/app_route.php:88-95 =================
function app_router($app) {
	(empty($app) || (!is_dir($basePath = A_P . $app))) && Showmsg('undefined_action');
	$baseFile = $basePath . '/example.com';
	if (!file_exists($baseFile)) { Showmsg("包含文件不存在，请创建example.com"); }
	return array($basePath, $baseFile);
}
// :71  require_once Pcv($baseFile);
// A_P = R_P.'apps/'（target116.example.com:16）
// → apps.php?q=../../hack/bank：is_dir("apps/../../hack/bank") 能过，但 Pcv($baseFile) 默认 $ifcheck=1
//   会 str_replace('..','') 比对 → 命中即 exit('Forbidden')  ★ 穿越被 Pcv 拦住
// → 编码绕过：%2e%2e%2f 不被 Pcv 识别为 '..'（Pcv 查字面量），但 is_dir 也不认 %2e → 同样过不了
// → apps.php 侧【无文件包含】
// 注：$q 由 InitGP 装载（cvtype=1）→ 先过 Char_cv；Char_cv 不处理 '.' 和 '/'，所以文件名拼接面完全交给 Pcv 兜底。

// ================= mode/o/ 目录实况（★ 管理员关注的"m_*.php 清单"）=================
// mode/o/ 下实际存在的 m_*.php（23 个文件里只有这 7 个 m_ 文件）：
//   m_app.php     应用列表/我的应用
//   m_board.php   留言板
//   m_browse.php  随便看看（聚合浏览）
//   m_friend.php  好友
//   m_home.php    个人空间首页
//   m_invite.php  邀请
//   m_myapp.php   我的应用管理
//   m_user.php    用户资料页
// ★ 注意：**没有 m_photos.php / m_write.php / m_diary.php / m_share.php / m_article.php / m_hot.php / m_groups.php**
//   → 这些 q 全部走 target212.example.com:42-43 转交 apps.php，再由 app_route.php:88 路由到 apps/<q>/example.com
// → 所以 mode.php?m=o&q=photos 实际执行的入口是 apps/photos/example.com，不是 mode/o/ 下的文件
// 每个 m_*.php 的头两行都是同一形态守卫：
//   m_app.php / m_board.php / m_browse.php / m_friend.php / m_home.php / m_invite.php / m_myapp.php / m_user.php
//   第 1~3 行统一为 !defined('M_P') && exit('Forbidden'); 一类判定 → 直接 URL 访问全部输出 Forbidden，不能独立当入口
//   （mode/o/example.com 是 !defined('M_P') && exit('Forbidden'); 且 :6 ObHeader("mode.php?m=o")）
// permission 检查：mode.php 段没有 admincheck，但各 m_*.php 内部按需 $winduid 判空（如 require/target148.example.com:12 !$winduid && Showmsg('not_login')）

// ================= mode/o/ 内的 /e 与 eval（⑤）=================
// mode/o/require/target148.example.com:96
//   preg_replace("/\[share\]([\w\.]+),([^\[,]+?),([\d]{1,10})\[\/share\]/eis","shareEvea('\\1','\\2','\\3')",$message);
//   → 替换串是常量 + 捕获组引用，且 \1 的字符类 [\w\.]、\2 排除了 [ 和 ,、\3 纯数字 → 不可断引 → 安全
// mode/o/ 内无 eval、无 create_function、无第二处 /e
// mode/o/example.com:73 / m_home.php:95 / m_friend.php:108 / m_browse.php:103 写 o_indexset_cache.php 与 mode_o_mosort.php
//   全部使用 pw_var_export() → 安全