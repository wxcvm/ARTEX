## 前提（源码级确证）
- `require/target137.example.com:10` `preg_replace("/\[code\](.+)?\[\/code\]/eisU","phpcode('\\1')",$message,$db_cvtimes)` —— **/e 在本源码中真实存在**（知识库 POC 140 A 区"已全面清除 /e"的表述与此矛盾，以本条为准；目标侧父代也已实证 /e 真实执行）。
- RCE 充要条件：eval 串 `phpcode('捕获')` 中出现**裸单引号**（或能让字符串在预期位置闭合的裸引号）。
- 两道输入保护（require/security.php）：
  - init `S::slashes` 对 $_GET/$_POST/$_COOKIE 做 addslashes（:154-158，magic_quotes 关闭时必走）→ 输入 `'` 变 `\'`；
  - `S::gp` 默认 cvtype=1 → `escapeChar→escapeStr` 把 `'`→`&#39;`（:210-234）。
- 还原裸引号只有两条路：**strip 在 escape 之前且之后不再转义**（`\'`→`'`），或 **escape 之后 decode**（`&#39;`/`%27`→`'`）。

## 方向一：无需账号的裸引号→convert 路径 —— 白盒证否
1. 游客可达的 convert 入口全站仅 2 个，均先转义：
   - `actions/job/target138.example.com:13` `escapeChar(stripslashes(S::getGP('atc_content','P')))` → `:39 convert`（`wordsConvert` 只做敏感词 Trie 替换，target139.example.com:86-105，不解实体）；
   - `actions/job/activity.php:9` `S::gp(...,'P')` 默认 escapeChar → `:22 convert`。
2. `pw_example.com:9-34` 游客白名单 21 个 action（tag/login/showface/…/bubble 等）无一调用 bbscode convert；`target47.example.com:43`、`target140.example.com:94` 的 convert 都不在游客白名单。
3. decode-after-escape 全部核过：
   - **`example.com:315-316`**（subject 动作）：`S::gp('atc_content','P')` → `html_entity_decode(urldecode(...))`。escapeStr 只拦 `%3C/%3E/%00`，**不拦 `%27`**；POST 双编码 `%2527` 经 PHP 一次解码后剩 `%27`，urldecode 还原成**裸 `'`** —— 这是全站唯一坐实的"转义后还原裸引号"点。但 sink 是 `pwQuery::update(pw_threads.subject)`（:327）+ echo（:339），且 `:293` 要求 GM/版主，**不进 convert、需高权** → 不构成 /e RCE（最多标题注入/XSS 面）。
   - `apps/ajax/example.com:506` `convert(stripslashes($title),$db_windpost)`：title 先经 S::gp escapeChar → `\&#39;` → stripslashes 只去反斜杠 → **`&#39;` 实体进 convert** → 安全（先转义后 strip 不还原实体）。`:253` 同理（且 `:200` 只还原 `&#61;`/`%26`，不含引号）。
   - `apps/kmd/example.com:208,245` decode 后仅 `parse_url`+preg 取 tid（整型），不进 convert。
   - `target141.example.com:80,82` ENT_QUOTES decode 只用于 strlen 判长；`target142.example.com:309` ENT_QUOTES decode 是 feed 输出（XSS 面）。
   - `pwHtmlspecialchars_decode`（target105.example.com:1906-1911）只认 `&#039;`（带前导零），**不还原 escapeStr 产出的 `&#39;`** → 对 escapeChar 输出天然无效；全部 11 个调用点（register/post.class tags/articlemodule subject/weibo 等）无一接 convert。
   - `example.com:141-161 getPcviewdata` 还原 `&#41;`→`)`，值只进自定义字段数组，不进 convert。
   - `stripslashes(...)` 与 `convert(...)` 的配对全站仅 `apps/ajax/example.com:506` 一处（已证安全）。
   - GBK：输入路径不存在"escape 之后再 iconv/mb 转码"的点（charset 即 GBK，无转码），宽字节救不出实体引号。
4. `passport_client.php`（S::gp userdb 不转义）被 `md5(action.userdb.forward.pptkey)` 与 `db_pptifopen` 双闸挡住（:7-13），且不进 convert。
⇒ **无需账号的裸引号→convert 路径：不存在。**

## 方向一附带修正：账号链预期要降温
- **论坛发帖链（postnew/postreply/postmodify/ajax:157 等 9 处 S::gp(...,0)）即使拿到账号也不会裸引号**：内容保留 init 的 `\'` → `lib/forum/target143.example.com:546` `convert()` 在 `:480 escapeChar` **之前**执行，但 eval 串里是 `\'`（单引号内转义，无 breakout）→ 之后 escapeChar → `:480` 入库为 `&#39;` → 展示 convert（target144.example.com:995）也是实体。**缺的不是账号，是"去斜杠且不再实体化"的点**。父代"补一个 groupid≠7 的账号即可打通 [code] 引号链"的假设，按本源码不成立。
- **CMS 投稿链才是真正的裸引号入库→展示 /e 链（代码级成立；前置①CMS 模式需后台开启、②账号组需 cms_post 右 —— 均未在目标验证）**：
  1. `mode/cms/m_post.php:30` `S::gp(array('cms_subject','atc_content','cms_descrip'),'P',0)` —— 不转义，内容 = `\'`；
  2. `mode/cms/lib/target145.example.com:80` `htmlspecialchars($content)`（PHP5.3 默认 **ENT_COMPAT 不编 `'`**），全程无 escapeChar；`:61` setSubject 还 `pwHtmlspecialchars_decode`；
  3. 落库 `target146.example.com:156 → articlecontentdb:7 → basedb:148 pwQuery::insert → target147.example.com:155 S::sqlEscape = addslashes(stripslashes($v))` —— **净效果剥掉保护斜杠，pw_cms_articlecontent 存的是裸 `'`**；
  4. 展示 `mode/cms/m_view.php:77` `convert($content,$db_windpost)` → `target137.example.com:10` /e → eval `phpcode('…裸\'…')` 逃逸 → **RCE**。
  - 投稿权限：`mode/cms/require/target148.example.com:38 getPostPurview` = isGM 或 `$_G['cms_post']`（按用户组右，`mode/cms/config/level.php:6` 注明"普通会员用户组可投稿"）；`m_post.php:6` 挡游客。
  - 模式开关：`target116.example.com:761` 要求 `$db_modes[$modeName]['ifopen']` 为真才会 define M_P 加载 mode/cms。
  - 同型分支：`target143.example.com:482-487` —— ifsign>=2（用户组 htmlcode 右 + atc_html）时**跳过 escapeChar**，内容 `\'` 经 sqlEscape 后同样以裸引号入库 → `target144.example.com:995` 展示 convert → 同样 RCE（需带 htmlcode 右的组）。
  - 另两处 /e sink（`target143.example.com:483`、`target145.example.com:87` 的 `/<(...)>/eis` _jscv 分支）在输入时引号仍带斜杠，未见 breakout。

## 方向二（apikey）要点摘要
- `pw_api.php:4` `define('P_W','admincp')`；`:50 $api->run($_POST+$_GET)`；签名 `api/class_base.php:93` `empty($this->apikey) || md5($arg.$apikey) != $request['sig']` → Error Sign 无 oracle 区分"空 key"与"错 key"；`unserialize(params)`（:98）在签名之后，无未签反序列化。
- `db_siteownerid` **不在** `data/sql_config.php`（require/target149.example.com:32 写入清单只有 db/manager/attach/slave），而在 **DB 表 pw_config**（target121.example.com:464 REPLACE INTO；admin/target150.example.com:110-111 重新生成）+ 运行时 `$GLOBALS`（api/class_base.php:77）+ 缓存 `data/bbscache/config.php`（pw_api.php:14 读取，HTTP 访问只会执行输出空白）。
- 打印点全部要登录或出站：target107.example.com:67/203（:174 `!$windid && Showmsg('not_login')`）、appclient（向 app.phpwind.net 出站）、admin 侧。无未授权 sig 样本。
- `generatestr`（target121.example.com:741-750）`mt_srand((double)microtime()*1000000)` 种子空间仅 10^6，但终值 = `substr(md5(str.microtime().HTTP_HOST.UA.db_hash),0,32)` —— 还需安装秒级时间、安装者 UA、db_hash，**离线不可爆**。
- md5(known_arg+apikey) 即使拿到样本也不可逆；知识库 137 条、tgt9215/tgt9942/tgt7217/q=apikey|siteownerid|uc_key|pw_api 全部检索 → **零命中**，无可复用 key 或签名材料。

## 方向三（开奖系统）检索记录
- 知识库 `q=开奖/kj.php/g_history/六合彩` → 0；`component=开奖系统` → 0；nuclei 13764 模板 matched=0。
- 互联网：DuckDuckGo 精确 `"g_history" "kj.php"` → **0 结果**；`kj.php g_history 开奖系统` 只有北京28/PC28 站点与无关 API 博客，**无该盗版开奖系统源码流通**。（web_search 工具 DeepSeek 402 余额不足不可用；grep.app 429。）
⇒ 白盒 demo 端暂无源码来源，等 demo 线自身取证（备份/报错回显）。