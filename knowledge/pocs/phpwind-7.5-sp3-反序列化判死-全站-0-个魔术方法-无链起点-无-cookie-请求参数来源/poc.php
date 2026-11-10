// ============================================================================
// 结论：PHPWind 7.5 SP3 反序列化对象注入【彻底不成立】——不是"找不到链"，
//       而是【全站没有任何魔术方法定义】，链的起点不存在。
// ============================================================================
//
// 【一】全站魔术方法普查（grep 结果）
//   grep 'function __(destruct|wakeup|toString|call|callStatic|get|set|isset|unset|invoke|set_state|clone)\s*\('
//   → 【0 命中】
//   grep 'function __(getKeyPrefix|construct)' 之类：
//   · __construct 共 9 处（uc_client/model/{user,notify,credit,app}.php:9、uc_client/control/{user,credit}.php、
//     lang/step/target220.example.com:17、uc_client/class_core.php:31、lib/target221.example.com:16、
//     require/target222.example.com:16/229、lib/target223.example.com:18）
//   · lib/target224.example.com:190 的 __getKeyPrefix 是【普通方法】，不是 __get 魔术方法
//   ★ 关键：unserialize() 只触发 __wakeup()，【不调用 __construct】。全站没有 __wakeup，
//     也没有任何其它反序列化相关魔术方法 → unserialize 一个对象字符串只会构造出普通对象，不会执行任何用户代码。
//
// 【二】unserialize( 全量（329 处命中，抽样全核）——数据来源分类
//   A. 数据库字段（绝大多数）：
//      pw_config.db_value（admin/target117.example.com:563、admin/target225.example.com:54、admin/rebang.php:9、lang/step/target226.example.com:41 等）
//      pw_forums.forumset / creditset（example.com:362-364、m/wap_global.php:22/44、apps/stopic/example.com:31-33）
//      pw_threads.aid（target194.example.com:1491/1506、admin/target227.example.com:201/220/435、require/target191.example.com:296、apps/diary/example.com:181/236/545/600）
//      pw_tools.conditions（hack/toolcenter/example.com:150/279/428、require/target228.example.com:13）
//      pw_tradeinfo / pw_upgrade / pw_rightset / pw_advert.config 等
//   B. 缓存文件内容：admin/target229.example.com:841-847（wordsfb_progress.php）、lib/target139.example.com:575（wordsfb.php）
//   C. 远程响应（需 SSRF/中间人，且多为后台功能）：
//      admin/target150.example.com:181/211/248（app.phpwind.net）、lib/target230.example.com:436、
//      uc_client/model/target150.example.com:132、uc_client/model/target231.example.com:77/80/101
//   D. API 请求体：api/class_base.php:98 unserialize($request['params'])  ← API 入口在目标上 404（已知情报）
//   ★ 【没有任何一处来源是 $_COOKIE】——cookie 全程只走 StrCode 的异或解密，不进 unserialize
//   ★ 也没有一处来源是 $_GET / $_POST 直接透传（唯一沾边的是 admin/target232.example.com:11
//     unserialize(base64_decode($query))，但那是后台且 $query 来自 DB）
//
// 【三】公开情报的版本错位（重要）
//   博客园《PHP wind反序列化分析》(https://www.cnblogs.com/fr09/p/18241189) 给出的 POP chain 为：
//     class PwDelayRun { private $_callback; private $_args;
//                        public function __destruct() { foreach($this->_callback as $k=>$v) call_user_func_array($v, $this->_args[$k]); } }
//     配合 Wind::autoLoad 的 spl_autoload_register 做类加载，入口是某控制器的 beforeAction 里 unserialize($this->getInput('var'))
//   ★ 这些全部属于【PHPWind 9.x 的 Wind 框架】（namespace src\library\...、WIND_PATH、spl_autoload_register）
//   ★ 在 7.5 SP3 里 grep 'PwDelayRun' / 'spl_autoload_register' / 'Wind::' / 'beforeAction' → 【全部 0 命中】
//   → 该公开链对目标（7.3.2，比 7.5 更早）不适用。不要照抄。
//
// 【四】结论
//   ① 魔术方法 0 → 无 POP chain 起点 → 死（语言层，与版本/清洗/权限无关）
//   ② 数据来源里没有 cookie / 没有直接请求参数 → 即使有链也没有入口 → 死
//   ③ 若 7.3.2 存在 7.5 没有的魔术方法类，需重审；但 7.3.2 更早，Wind 框架尚不存在，概率极低
//
// 【五】PHP 5.5.38 相关技巧说明（供参考，本次无用武之地）
//   · CVE-2016-7124（__wakeup 绕过，属性数大于实际数）影响 PHP 5.6.25 之前 / 7.0.10 之前 → 5.5.38 受影响
//     但本项目【没有任何 __wakeup】，该技巧无目标可打
//   · O:4:"Name":N:{} 的属性名可控（含 \0 前缀的 private/protected 属性名）+ 对象类型混淆在 PHP 5 可用
//     → 同样是"有链才能用"，本项目无链