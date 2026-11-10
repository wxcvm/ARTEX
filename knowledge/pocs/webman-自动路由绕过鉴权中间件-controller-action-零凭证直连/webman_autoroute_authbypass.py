## 1. 原理
webman 默认开启自动路由：请求路径 `/{controller}/{action}`（controller 为 `app/controller` 下相对路径的小写形式，action 为方法名原样）会被直接绑定到控制器方法。
如果项目把鉴权写在**显式路由组**上（`Route::group(...)->middleware([AuthMiddleware::class])`）而没有全局 `Route::fallback`，那么自动路由这条路径**不经过任何鉴权中间件**。

典型源码特征：
```php
// config/middleware.php —— 全局中间件里没有鉴权
return ['' => [ExceptionMiddleware::class, CorsMiddleware::class, RateLimitMiddleware::class, EncryptionMiddleware::class]];
// config/route.php —— 鉴权只挂在显式路由组
Route::group('/api/v1', function () { ... })->middleware([app\middleware\AuthMiddleware::class]);
```

## 2. 验证对照（最关键的证据形式）
同一业务方法，两条路由结果不同即证明：
```
GET /api/v1/agent/commission-rates   -> {"code":10003,"message":"缺少认证信息"}   # 显式路由，被中间件拦
GET /api/agent/commissionRates       -> 200 {"code":0,...,"data":{...}}          # 自动路由，直出数据
```
另一个佐证：`GET /api/xxxController/info` 返回「方法体自己抛的 401」而不是中间件的「缺少认证信息」，说明中间件已被跳过。
对照不存在的控制器（如 `/api/safetyController/info`）应返回框架 404。

## 3. 枚举脚本（静态 + 在线）
静态：抽取 `app/controller/**/*.php` 所有 public 方法，用大括号配平取方法体，正则找 `Db::table('x')`、`$request->post/get('p')`、以及 `userId|uid|token|未登录|请登录` 判断方法体内有无登录校验。
在线：
```python
def variants(origin):  # origin 形如 "api\\AgentController.php::commissionRates"
    cp, meth = origin.split("::")
    parts = cp.replace("\\", "/")[:-4].split("/")
    cls, dirs = parts[-1], parts[:-1]
    low = (cls[:-10] if cls.endswith("Controller") else cls).lower()
    pre = "/" + "/".join(dirs + [low])
    out = [pre + "/" + meth]
    sn = re.sub(r"(?<!^)(?=[A-Z])", "_", meth).lower()   # 驼峰/下划线两种都试
    if sn != meth: out.append(pre + "/" + sn)
    if dirs: out.append("/" + low + "/" + meth)          # 目录被吃掉的情况
    return out
```
判定：404=不存在；200 且 body 带 `data`/`list` 业务数据=命中；`401/10003`=被拦（换变体）。
注意：若项目开了响应加密（`X-Encrypted: 1`），**必须解密后再判定**，否则会把大量命中误当成空响应漏掉。

## 4. 副作用提示
自动路由同样覆盖写类方法（add/edit/delete/save/audit...），因此还能**绕过鉴权改数据**；
但探测时只应打读类方法，写类用「必然失败但无害」的不存在 id 验证可达性。
自动路由路径通常也不在限速中间件的 path 表里（RateLimit 多按精确 path 匹配），等于同时绕过 IP 级限速。
