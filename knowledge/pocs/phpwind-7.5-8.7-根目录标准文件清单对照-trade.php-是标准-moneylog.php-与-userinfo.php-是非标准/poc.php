// ============================================================================
// 目的：把 7.5 SP3 与 8.7 的根目录 PHP 文件清单固化下来，用于快速识别目标上的"非标准文件"
//       —— 非标准文件（站长手加 / 第三方插件 / 老版本遗留）写法上更可能带洞。
// ============================================================================
//
// 【7.5 SP3 根目录 .php 全量 51 个】（runs/pw-src-732/pw75sp3/）
// active.php example.com example.com alipay.php apps.php cate.php ck.php columns.php faq.php forumcp.php
// global.php example.com header.php hitcache.php example.com install.php example.com jobcenter.php left.php login.php
// masingle.php mawhole.php example.com message.php mode.php new.php notice.php operate.php passport_client.php
// pay99bill.php paypal.php post.php profile.php push.php pw_example.com pw_api.php read.php register.php
// remail.php rss.php search.php sendemail.php sendpwd.php show.php sitemap.php sort.php tenpay.php
// thread.php trade.php u.php userpay.php
// 根目录子目录：admin api apps attachment data hack htm_data images ipdata js lang lib m mode require simple template uc_client
//
// 【8.7 根目录 .php 全量 58 个】（本地对照样本 <HOME>\Desktop\168\runs\hklot\pw-src\upload）
// act_alipay_push.php act_alipay_receive.php active.php example.com example.com alipay.php apps.php cate.php
// ck.php ckquestion.php columns.php faq.php forumcp.php global.php group.php example.com header.php hitcache.php
// example.com install.php example.com jobcenter.php left.php link.php login.php masingle.php mawhole.php
// mawholecolony.php example.com message.php mode.php new.php notice.php operate.php passport_client.php
// pay99bill.php paypal.php post.php profile.php push.php pw_example.com pw_api.php pweditor.php read.php
// register.php remail.php rss.php search.php searcher.php sendemail.php sendpwd.php show.php sitemap.php
// slide.php sort.php tenpay.php thread.php trade.php u.php userpay.php y.php
// 8.7 子目录：actions admin api apps attachment connexion data hack hook html images ipdata js m mode
//             require simple template u uc_client
//
// 【两版共有】（7.3.2 极可能也有的高置信集合）
// active.php example.com example.com alipay.php apps.php cate.php ck.php columns.php faq.php forumcp.php global.php
// example.com header.php hitcache.php example.com install.php example.com jobcenter.php left.php login.php masingle.php
// mawhole.php example.com message.php mode.php new.php notice.php operate.php passport_client.php pay99bill.php
// paypal.php post.php profile.php push.php pw_example.com pw_api.php read.php register.php remail.php rss.php
// search.php sendemail.php sendpwd.php show.php sitemap.php sort.php tenpay.php thread.php trade.php u.php userpay.php
// ★ trade.php 在 7.5 与 8.7 都存在 → 目标上的 trade.php 是标准文件，不是站长加的 ✓
//
// 【7.5 与 8.7 的差异】
// 仅 7.5 有：无（7.5 的都在 8.7，除 left.php/header.php 这类公共文件）
// 仅 8.7 有：act_alipay_push.php act_alipay_receive.php ckquestion.php group.php link.php mawholecolony.php
//            pweditor.php searcher.php slide.php y.php（8.x 新增/拆分）
//
// 【★ 非标准文件判定：目标上发现但 7.5/8.7 都没有的】
//   · example.com —— 7.5 与 8.7 根目录【均无】；7.5 全站 grep "example.com" 0 命中。
//     同族最接近的是 admin/creditlog.php（后台积分/币操作记录，需后台）
//   · example.com —— 7.5 与 8.7 根目录【均无】；7.5 全站 grep "example.com" 0 命中。
//     同族最接近的是 u.php（用户资料页，7.5/8.7 都有）
//   ★ 这两个文件的两种可能：
//     (a) 7.3.2 时代的标准文件名（7.3.2 早于 7.5，可能用 example.com 显示用户资料、用 example.com 显示币收支，
//         7.5 重构时合并/改名）。若如此，它们的代码会非常老（7.3 分支），写法可能比 7.5 更原始 → 值得单独拿到源码审。
//     (b) 站长手加 / 第三方插件（第三方插件往根目录丢 .php 很常见）→ 更可能有洞（无 PHPWind 的权限框架保护）。
//   ★ 目标侧可自行区分：观察这两个文件的 mtime 是否与站内其它标准文件一致、是否带 PHPWind 的
//     `!function_exists('readover') && exit('Forbidden');` 守卫、页面页脚是否输出 "Powered by PHPWind"。
//     - 带 PW 守卫 + 页脚 → 属 7.3.2 标准文件（老版本代码）
//     - 无守卫、页脚无 PW 或样式脱节 → 站长/插件手加
//
// 【其它非标准/常被误解的文件】columns.php 与 faq.php 都是 7.5 与 8.7 的标准文件（不是站长加的）；
// alipay.php / tenpay.php / paypal.php / pay99bill.php / userpay.php 也都是标准支付文件
// —— 之前已确认它们不写文件、无变量 include，无需再投入。
// 另外 api/ 与 uc_client/ 两个目录在目标上整体 404（已知情报），对应文件不必再试。