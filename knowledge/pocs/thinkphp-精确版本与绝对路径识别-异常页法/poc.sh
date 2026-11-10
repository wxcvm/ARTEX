## 识别方法（三步，零风险、纯被动）

### 1. 响应头初筛
```
curl -sS -k -i "https://<target>/"
```
命中特征头：`X-Powered-By: ThinkPHP`、`Server: Nginx`；并可能伴随第三方 CDN 私有头。
ThinkPHP 3.x 的 Set-Cookie 只有 `PHPSESSID`，无框架专属 cookie。

### 2. 路由格式判定（3.x vs 5.x 的关键区分）
ThinkPHP **3.x** 使用 `g=`/`m=`/`a=` 三元组：
```
/example.com?g=Home&m=Index&a=login          # 模块=Home 控制器=Index 方法=login
/example.com?g=Home&m=Pay&a=pay
```
ThinkPHP **5.x/6.x** 使用 `/example.com/<module>/<controller>/<action>` 或 `?s=/a/b/c`。
若站点同时支持 `PATH_INFO`（`/Index/login` 直接可访问），说明配置了 url_rewrite。

### 3. 异常页取精确版本 + 绝对路径（最有效）
请求一个控制器下的**不存在模板/方法**，触发模板不存在异常：
```
curl -sS -k "https://<target>/User"          # 模块/控制器存在但无模板
curl -sS -k "https://<target>/example.com?g=Home&m=Pay&a=index"
```
返回体特征：
```html
<title>系统发生错误</title>
<h1>模板不存在[/www/wwwroot/<项目名>/App/Tpl/Home/User_example.com]</h1>
<p><a title="官方网站" href="http://www.thinkphp.cn">ThinkPHP</a><sup>3.1.3</sup>
{ Fast & Simple OOP PHP Framework } -- [ WE CAN DO IT JUST THINK ]</p>
```
一次命中同时得到三样东西：**精确版本 3.1.3**、**服务器绝对路径**（含宝塔面板标准目录 `/www/wwwroot/`）、**模板目录结构**（`App/Tpl/<模块>/<控制器>_<方法>.html`）。

## 指纹库落库值
| 字段 | 值 |
| --- | --- |
| category | 框架 |
| product | ThinkPHP |
| version | 3.1.3 |
| evidence | `X-Powered-By: ThinkPHP`；异常页 footer `ThinkPHP 3.1.3`；模板路径 `App/Tpl/Home/*.html` |

## 该版本对应的已知高危面（拿到版本后即应逐条验证）
- ThinkPHP 3.x 缓存文件包含 / 模板渲染类 RCE（`display()` + 可控模板变量、`__APP__` 常量替换）
- ThinkPHP 3.x `where()` 数组 / `exp` 表达式 SQL 注入（`logic`/`where` 参数传入数组）
- ThinkPHP 3.1.3 多语言 `l=` 参数与 `s=` 参数解析缺陷
> 注意：本方法只确认版本与路径，**不确认漏洞存在**。必须在目标上实测回显才能记为 confirmed。

## 附加价值
- 异常页泄露的绝对路径可直接用于：文件包含 payload 构造、日志投毒路径、`/www/wwwroot/` 宝塔站点定位。
- 若站点开启 debug，同类请求会输出完整调用栈与 SQL，价值更高。
