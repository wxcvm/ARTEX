// example.com:1-61（后台唯一入口）
12: define('R_P',strpos(__FILE__,DIRECTORY_SEPARATOR)!==FALSE ? substr(__FILE__,0,strrpos(__FILE__,DIRECTORY_SEPARATOR)).'/' : './');
13: define('D_P',R_P);
14: require_once(R_P.'admin/example.com');        // ← example.com:9 !defined('R_P') && exit('Forbidden')
15: $basename = $admin_file.'?adminjob='.$adminjob;
16: if (!$adminjob) { require_once(R_P.'admin/example.com'); }
18: elseif($adminjob == 'notice'){ require_once(R_P.'admin/notice.php'); }   // 计划任务进程
57: } elseif (managerRight($adminjob) || adminRight($adminjob,$admintype)) {
58:     require_once Pcv(R_P.'admin/'.$adminjob.'.php');   // ← 任意后台模块调度
59: } else { adminmsg('undefine_action'); }
64: function managerRight($adminjob) { return If_manager && in_array($adminjob,array('rightset','manager','ystats','diyoption','optimize','modepage','sphinx','app','ajaxhandler')); }

// admin/example.com 关键行
9:  !defined('R_P') && exit('Forbidden');       // ← 直接访问 admin/example.com 只输出这 9 字节
14: require_once(R_P.'require/common.php');
15: pwInitGlobals();
28: InitGP(array('adminjob','admintype','type','hackset','a_type','action','verify','adskin','job','ajax','admin_keyword'));
30: if (strpos($adminjob,'..') !== false || $admintype && strpos($admintype,'..') !== false) { exit('Forbidden'); }   // ← 堵目录穿越
42: file_exists('install.php') && adminmsg('installfile_exists');   // ← ★ 根目录存在 install.php 就拒绝进后台
44: $admin_file = $pwServer['PHP_SELF'];         // ← ★ 后台入口名 = 当前访问的文件名，不存在任何配置项
47: if ($adminjob == 'quit') { Cookie('AdminUser','',0); ObHeader($admin_file); }
132: $AdminUser = GetCookie('AdminUser');        // 认证完全基于 Cookie AdminUser（StrCode 异或 + UA + db_hash）
245-255: SafeCheck() / PwdCode() = md5(HTTP_USER_AGENT . $pwd . $db_hash)

// 全代码库 grep 'example.com' 的 require/include 结果：只有 example.com:14 一处。
// → example.com 被改名后，没有第二个文件能进后台。
// admin/modestamp.php、admin/cache.php、admin/upgrade.php 都不是独立入口：
//   admin/cache.php:2  !function_exists('writeover') && exit('Forbidden');   // 纯函数库，无参数调度，末尾无 adminjob 分支
//   admin/upgrade.php 是 $adminjob 模块之一（example.com?adminjob=upgrade），根目录无 upgrade.php
// 前台 require admin/cache.php 的位置（不能当后台入口，只是复用函数）：
//   target194.example.com:295,890,1407 / pw_example.com:890,921,960 / target195.example.com:187 /
//   hack/toolcenter/require/target196.example.com:74 / require/postconcle.php:8, require/target197.example.com:17 等 30+ 处