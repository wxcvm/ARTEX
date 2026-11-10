// passport_client.php 全文（92 行，7.5 SP3）
1: <?php
2: require_once('global.php');
3: require_once(R_P.'require/checkpass.php');
4:
5: InitGP(array('action','forward','verify'));
6: InitGP(array('userdb'),'GP',0);                 // ★ cvtype=0 → 不过 Char_cv
7: if(!$db_pptifopen || $db_ppttype!='client'){
8: 	Showmsg('passport_close');                    // ← 功能未开则直接退出
9: }
10: if(empty($db_pptkey) || md5($action.$userdb.$forward.$db_pptkey) != $verify){
11: 	Showmsg('passport_safe');                     // ← 唯一的认证：md5(拼接,$db_pptkey)
12: }
13: $_db_hash=$db_hash;
14:
15: $db_hash=$db_pptkey;
16: parse_str(StrCode($userdb,'DECODE'),$userdb);     // ★ userdb = 异或+base64 编码的 query string
17:
18: if($action=='login'){
19: 	$userdb = Char_cv($userdb);
20: 	if(!$userdb['time'] || !$userdb['username'] || !$userdb['password']){
21: 		Showmsg('passport_data');
22: 	}
23: 	if($timestamp-$userdb['time']>3600){ Showmsg('passport_error'); }
24:
25: 	$member_field = array('username','password','email');
26: 	$memberdata_field = array('rvrc','money','credit','currency');
...
34: 	$rt=$db->get_one("SELECT uid $sql FROM pw_members WHERE username=".pwEscape($userdb['username']));
35: 	if($rt){
36: 		$sql=$sql2=array();
37: 		foreach($userdb as $key=>$val){
38: 			if(in_array($key,$member_field) && $rt[$key] != $val){ $sql[$key] = $val; }   // ★ 直接改写该用户 password/email
40: 			elseif(in_array($key,$memberdata_field) && strpos(",$db_pptcredit,",",$key,")!==false){ $sql2[$key] = $val; }
43: 		}
46: 		$db->update("UPDATE pw_members SET $sql WHERE uid=".pwEscape($rt['uid']));
50: 		$db->update("UPDATE pw_memberdata SET $sql2 WHERE uid=".pwEscape($rt['uid']));
53: 		$winduid = $rt['uid'];
54: 	} else{
55: 		// 用户不存在 → REPLACE INTO 新建用户，groupid=-1, memberid=8
63: 		$sql1 += array('groupid'=>-1,'memberid'=>8,'gender'=>0,'regdate'=>$timestamp);
69: 		$db->update("REPLACE INTO pw_members SET".pwSqlSingle($sql1));
70: 		$winduid = $db->insert_id();
78: 		$db->update("REPLACE INTO pw_memberdata SET".pwSqlSingle($sql2));
79: 		$db->update("UPDATE pw_bbsinfo SET newmember=".pwEscape($userdb['username']).",totalmember=totalmember+1 WHERE id='1'");
80: 	}
81: 	$db_hash=$_db_hash;
82: 	$windpwd = PwdCode($userdb['password']);
83: 	Cookie("winduser",StrCode($winduid."\t".$windpwd),$userdb['cktime']);   // ★ 直接下发登录 Cookie
84: 	Cookie('lastvisit','',0);
85: 	Loginipwrite();
86: 	ObHeader($forward ? $forward : $db_bbsurl);
87: } elseif($action=='quit'){
88: 	$db_hash=$_db_hash;
89: 	Loginout();
90: 	ObHeader($forward ? $forward : $db_bbsurl);
91: }
// ※ 根目录不存在 passport.php；hack/passport/ 是后台配置模块（hack/passport/example.com 只做 ifcheck($db_pptifopen,...) 写入）