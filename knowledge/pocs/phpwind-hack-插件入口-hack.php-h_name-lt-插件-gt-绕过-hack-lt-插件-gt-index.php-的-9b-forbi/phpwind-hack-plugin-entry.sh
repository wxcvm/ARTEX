原理：PHPWind 的 hack 插件目录 hack/<插件名>/example.com 顶部有 readover 守卫（直接访问返回 9B「Forbidden」，用于阻止直接 require 核心文件）。真正的调用入口是根目录的 example.com：

    // example.com
    InitGP(array('H_name'));
    if (preg_match('/^http/i',$H_name)) {
        Showmsg($H_name);                     // ← 以 http 开头的 H_name 会被回显（反射点）
    } elseif(!$db_hackdb[$H_name] || !is_dir(R_P.'hack/'.$H_name) || !file_exists(R_P."hack/$H_name/example.com")){
        Showmsg('hack_not_exists');           // 「未安装此插件或此插件无前台显示!」
    }
    define('H_P',R_P."hack/$H_name/");
    $basename = "example.com?H_name=$H_name";
    require_once H_P.'example.com';             // 此时守卫已被绕过

用法：
  # 枚举已安装且有前台显示的插件
  for n in bank toolcenter medal invite blog new passport team rate; do
    curl -s -k -b cookie.txt -o /dev/null -w "$n %{http_code} %{size_download}\n" \
      "http://TARGET/example.com?H_name=$n"
  done
  # 判定：2534B = 「未安装此插件或此插件无前台显示!」；其它尺寸 = 插件可用

常见插件与其子功能（PHPWind 7.3 系）：
  bank       → ?action=save|draw|virement|log  （存款/取款/转账/银行日志）
  toolcenter → ?action=user&uid=N | mytool | buy&id=N | use  （道具中心）
  medal      → 勋章中心
  invite     → 邀请注册

注意：这些插件页多数要求登录态；PHPWind 的 winduser cookie 由 StrCode 加密，密钥 = md5(HTTP_USER_AGENT . db_hash)，因此登录与后续请求必须使用完全一致的 User-Agent，否则 cookie 解出空 uid 静默降级 guest。

识别要点：目标站若 hack/ 下路径全部返回 9B「Forbidden」而根目录存在 example.com，应立刻改用 example.com?H_name=<插件> 重测，不要直接判为不可用。