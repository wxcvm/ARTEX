// example.com:238-315 —— action=uploadicon（7.5 SP3）
} elseif ($action == 'uploadicon') {
	if (empty($_GET['step'])) {
		// 路径 1：走 FaceUpload 行为类 + PwUpload::upload()
		list($db_upload,$db_imglen,$db_imgwidth,$db_imgsize) = explode("\t",$db_upload);
		InitGP(array('uid','verify'));
		$swfhash = GetVerify($uid);
		checkVerify('swfhash');
		require_once(R_P . 'lib/upload/faceupload.class.php');
		$face = new FaceUpload($uid);
		PwUpload::upload($face);
	} else {
		// 路径 2：raw POST 写入
		$ext = strtolower(substr(strrchr($_GET['filename'],'.'),1));      // ← 扩展名来自 GET filename
		$udir = str_pad(substr($winduid,-2),2,'0',STR_PAD_LEFT);
		$source = PwUpload::savePath(0, "{$winduid}_tmp.$ext", "upload/$udir/");
261:		if (!in_array(strtolower($ext), array('gif','jpg','jpeg','png','bmp')) || !file_exists($source)) {
		    Showmsg('undefined_action');                                   // ← 7.5 的扩展名白名单
		}
264:		$data = $_SERVER['HTTP_RAW_POST_DATA'] ? $_SERVER['HTTP_RAW_POST_DATA'] : file_get_contents('php://input');
266:		if ($data) {
271:			$filename  = "{$winduid}.jpg";
277:			$middleFile = PwUpload::savePath($db_ifftp, $filename, "$middleDir", 'm_');
279:			writeover($middleFile, $data);              // ← 任意 RAW 内容先落盘
282:			if (!$img_size = GetImgSize($middleFile,'jpg')) {
283:				P_unlink($middleFile); Showmsg('upload_content_error');  // ← 7.5 的内容校验
			}
287:			$normalFile = PwUpload::savePath($db_ifftp, "{$winduid}.$ext", "$normalDir");
292:			} elseif (!PwUpload::movefile($source, $normalFile)) { Showmsg('undefined_action'); }
297:			$smallFile = PwUpload::savePath($db_ifftp, $filename, "$smallDir", 's_');
		}
	}
}

// lib/target52.example.com:151-170
function postupload($tmp_name,$filename) {
152:	if (strpos($filename,'..') !== false || strpos($filename,'.php.') !== false || eregi("\.php$",$filename)) {
		exit('illegal file type!');        // ← 只堵 .. / .php. / .php$（eregi 忽略大小写）
	}
	PwUpload::createFolder(dirname($filename));
	if (function_exists("move_uploaded_file") && @move_uploaded_file($tmp_name,$filename)) { @chmod($filename,0777); return true; }
	elseif (@copy($tmp_name, $filename)) { @chmod($filename,0777); return true; }
	elseif (is_readable($tmp_name)) { writeover($filename,readover($tmp_name)); if (file_exists($filename)) { @chmod($filename,0777); return true; } }
	return false;
}
// 未拦截的备用后缀：pht / phtml / php3 / php4 / php5 / php7 / phps / inc（是否被 nginx 当 PHP 执行取决于宝塔站点配置）

// lib/target52.example.com:31-38 扩展名与体积校验
if (empty($upload['ext']) || !isset($bhv->ftype[$upload['ext']])) { uploadmsg('upload_type_error'); }
if ($upload['size'] < 1 || $upload['size'] > $bhv->ftype[$upload['ext']]*1024) { ... uploadmsg(...); }