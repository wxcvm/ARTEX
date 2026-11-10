// 附件/图片上传的扩展名"变形"防护（7.5 SP3）
// lib/upload/target53.example.com:175
$filename = $this->forum->fid . "_{$this->uid}_$prename." . preg_replace('/(php|asp|jsp|cgi|fcgi|exe|pl|phtml|dll|asa|com|scr|inf)/i', "scp_\\1", $currUpload['ext']);
// lib/upload/target51.example.com:30
$filename = "0_{$this->uid}_$prename." . preg_replace('/(php|asp|jsp|cgi|fcgi|exe|pl|phtml|dll|asa|com|scr|inf)/i', "scp_\\1", $currUpload['ext']);
// lib/upload/target53.example.com:98（另一处，与扩展名无关）
$rt['fileuploadurl'] = $filename = $thumbname = preg_replace('/^0_/', "{$this->forum->fid}_", $rt['attachurl']);

// 黑名单覆盖：php asp jsp cgi fcgi exe pl phtml dll asa com scr inf
// 注意：php5/php3/php4/php7/phps 会被子串 "php" 先命中 → 变 scp_php5，仍被拦。
// 未覆盖且常见的 PHP 可解析后缀：pht（还有 phtm、phar、inc、shtm、shtml）
//   → 若后台 $db_uploadfiletype（附件白名单）里出现 pht，则 .pht 能原样落地。
//   是否被 PHP-FPM 解析取决于宝塔 nginx 站点配置（默认只解析 .php，AddHandler 型配置则会解析 pht）。
// 另两道更前置的墙：
//   lib/target52.example.com:31  if (empty($upload['ext']) || !isset($bhv->ftype[$upload['ext']])) uploadmsg('upload_type_error');
//   lib/target52.example.com:152 if (strpos($filename,'..')!==false || strpos($filename,'.php.')!==false || eregi("\.php$",$filename)) exit('illegal file type!');
//   → postupload() 只堵 .. / .php. / .php$，并未堵 .pht / .php5 / .php3（大小写由 eregi 处理）
//   → 但真实文件名由 attupload/mutiupload 用白名单 ext 拼出，只有 $db_uploadfiletype 放开才有机会。
// 落盘目录：$attachdir/$dir（默认 attachment/），@chmod($dstfile,0777)（lib/target52.example.com:134/137/157/160/166）
// $db_uploadfiletype 的默认值集合不在 install_wind.sql 硬编码里（由后台 optimize/upload 设置生成）。