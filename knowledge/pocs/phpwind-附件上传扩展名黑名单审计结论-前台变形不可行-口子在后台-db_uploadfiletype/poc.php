## 结论（源码审计，逐行核对 runs/pw-src/）
**phpwind 8.x 的帖子附件上传，靠"扩展名变形"打不穿；真正的口子在后台上传类型设置项。**

## 一、黑名单正则原文（lib/upload/target51.example.com:30，attmutiupload.class.php 同款）
```php
$filename = "0_{$this->uid}_$prename." . preg_replace(
    '/(php|asp|jsp|cgi|fcgi|exe|pl|phtml|dll|asa|com|scr|inf)/i', "scp_\\1", $currUpload['ext']);
```
判定：
- **带 `/i`** ⇒ 大小写（`.PHP`/`.Php`）无效；
- **替换式**（`php`→`scp_php`）且 `preg_replace` 非重叠左到右 ⇒ **双写 `phpphp` → `scp_phpscp_php`，同样无效**；
- `.php3/.php4/.php5/.php7/.phps` 全含子串 `php` ⇒ 一律变 `.scp_phpX`，不可执行；
- 扩展名来源 `lib/target52.example.com:22` 只取**最后一个点之后** ⇒ `a.jpg/../../shell.php` 也只得 `php`（然后被替换）。

## 二、`.user.ini` 为什么也不行
`user.ini` 不在黑名单，但文件名前缀硬编码为 `"0_{$uid}_{$prename}."` ⇒ 上传 `.user.ini` 得到 **`0_1234_xxxxx.ini`**，**不是 `.user.ini`** ⇒ `auto_prepend_file` getshell 不成立。

## 三、唯一可行路径（推荐）
允许类型不是硬编码而是引用后台配置：
```php
// target51.example.com:18 / target53.example.com:51
$this->ftype =& $db_uploadfiletype;
```
⇒ **拿到后台后把附件允许类型加上 `php`，前台附件上传立刻 getshell**（比改模板更直接、更不易被 admin/filecheck.php 特征扫描拦）。

## 四、14 个上传类的白名单普查（runs/audit_ftype.py 自动化提取）
`faceupload/photoupload/pcupload/certificateupload/advupload/activityupload/activeupload/pushupload/stopicupload/spacebannerupload/imgupload` 等 **全部只允许 `gif/jpg/jpeg/bmp/png`**，`csvupload` 只允许 `csv`；**没有任何类把可执行扩展名放进白名单**。
仅 `actions/job/uploadicon.php` 是例外（见"phpwind 头像上传任意文件写入"那条），且本快照已含补丁。