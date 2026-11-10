// ===== 照片上传链（7.5 SP3）=====
// 入口：mode.php?m=o&q=photos&a=upload  →  target212.example.com:42 转交 apps.php  →  apps/photos/example.com
//   apps/photos/example.com:2  !defined('A_P') && exit('Forbidden');
//   apps/photos/example.com:3  !$winduid && Showmsg('not_login');          ★ 必须登录
//   apps/photos/example.com:14 InitGP(array('a','u','space','s'));          ← a=upload
//   apps/photos/example.com:18 !$db_phopen && Showmsg('photos_close');      ★ 相册功能开关

// 核心上传类 lib/upload/photoupload.class.php（全 88 行）
13: function PhotoUpload($aid) {
14: 	global $db,$o_maxfilesize;
19: 	!$o_maxfilesize && $o_maxfilesize = 1000;
21: 	$this->ftype = array(                       // ★★ 硬编码扩展名白名单，只有 5 种
22: 		'gif'  => $o_maxfilesize, 'jpg'  => $o_maxfilesize,
23: 		'jpeg' => $o_maxfilesize, 'bmp'  => $o_maxfilesize,
24: 		'png'  => $o_maxfilesize
25: 	);
28: function allowType($key) { return true; }
40: function getFilePath($currUpload) {
41: 	global $timestamp,$o_mkdir;
42: 	$prename  = randstr(4) . $timestamp . substr(md5($timestamp . $currUpload['id'] . randstr(8)),10,15);
43: 	$filename = $this->aid . "_$prename." . $currUpload['ext'];   // ★ 文件名 = {aid}_{服务端随机}.{ext}
44: 	$savedir  = 'photo/';
45: 	if ($o_mkdir == '2')  { $savedir .= 'Day_' . date('ymd') . '/'; }
47: 	elseif ($o_mkdir == '3') { $savedir .= 'Cyid_' . $this->aid . '/'; }
49: 	else { $savedir .= 'Mon_'.date('ym').'/'; }
52: 	return array($filename, $savedir, 's_' . $filename, $savedir);
55: function update($uploaddb) { ... INSERT INTO pw_cnphoto (aid,pintro,path,uploader,uptime,ifthumb) VALUES ... }

// 复用与校验链（lib/upload.class.php）
24: foreach ($_FILES as $key => $value) {
25: 	if (!PwUpload::if_uploaded_file($value['tmp_name']) || !$bhv->allowType($key)) { continue; }
29: 	$upload = PwUpload::initCurrUpload($key, $value);
17:   //   initCurrUpload 里： $arr['ext'] = strtolower(substr(strrchr($arr['name'],'.'),1));   ← 扩展名取自客户端原始文件名
31: 	if (empty($upload['ext']) || !isset($bhv->ftype[$upload['ext']])) { uploadmsg('upload_type_error'); }   // ★★ 第一道墙
34: 	if ($upload['size'] < 1 || $upload['size'] > $bhv->ftype[$upload['ext']]*1024) { ... uploadmsg('upload_size_error'); }
39: 	list($filename, $savedir, $thumbname, $thumbdir) = $bhv->getFilePath($upload);
42: 	$source = PwUpload::savePath($bhv->ifftp, $filename, $savedir);
45: 	if (!PwUpload::postupload($atc_attachment, $source)) { uploadmsg('upload_error'); }
50: 	if (in_array($upload['ext'], array('gif','jpg','jpeg','png','bmp','swf'))) {
52: 		if (!$img_size = GetImgSize($source, $upload['ext'])) { P_unlink($source); uploadmsg('upload_content_error'); }   // ★★ 第二道墙
56: 		if ($upload['ext'] != 'swf') {
57: 			if ($bhv->allowThumb())      { PwUpload::makeThumb($source, $thumburl, $bhv->getThumbSize(), $upload['ifthumb']); }
60: 			if ($bhv->allowWaterMark())  { PwUpload::waterMark($source, $upload['ext'], $img_size); }
64: 			$upload['type'] = 'img';

// postupload（lib/target52.example.com:151-170）—— 是的，复用了
152: if (strpos($filename,'..') !== false || strpos($filename,'.php.') !== false || eregi("\.php$",$filename)) { exit('illegal file type!'); }
// 只拦 .. / .php. / .php$（eregi 忽略大小写）；未拦 .pht/.phtml/.php5/.php3 —— 但本项目里 $filename 由 getFilePath 生成，
// 只带白名单扩展名，所以这些未拦后缀【到不了 postupload】。

// savePath（lib/target52.example.com:182-191）→ 实际落盘路径
182: function savePath($ifftp, $filename, $dir, $thumb = '') {
183: 	global $attachdir;
184: 	if ($ifftp) { ... $source = D_P . "data/tmp/{$filename}"; }
187: 	else { $source = $attachdir . '/' . $dir . $filename; }        // $attachdir = R_P.$db_attachname（target116.example.com:125）
190: 	return $source;
// → 落盘 = {DOCUMENT_ROOT}/attachment/photo/{Mon_YYMM|Day_YYMMDD|Cyid_<aid>}/{aid}_{randstr4}{ts}{md5(15)}.{jpg等}
// → @chmod($filename,0777)（lib/target52.example.com:157/160/166）