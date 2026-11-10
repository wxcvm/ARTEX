## 适用目标
`203.0.113.72` / `203.0.113.38` / `203.0.113.37`（本靶标 example.com 的 Windows 宝塔运维机组）。**任何人碰这三台机器，先看本条，可省两轮试错。**

## 一、203.0.113.38 / .23 —— 两条恒定假阳性（不要当成果）

### 假阳性 1：`/static/*` 的 401 是「整前缀级认证」
```
/static           → 404
/static/          → 401 (338B, Werkzeug 页面)
/static/zzz_nonexistent_9f3k  → 401   ← 随机不存在的文件名，同样 401
/static/aaa_random_777        → 401
/static/dump.sql              → 401   ← 不能据此认为该文件存在
/static/api/swagger.json      → 401   ← 同上
```
⇒ **401 只证明"该前缀被认证拦截"，不证明任何具体文件存在**。曾据此推断"面板 API 面下存在 dump.sql/swagger.json"，经随机文件名对照实验证伪并撤回。
**该前缀只存在于 8888 端口**：443/80/888 访问 `/static/dump.sql` 一律 404。

### 假阳性 2：任何字典里的 `login` 一词必产出 1 条 302 假阳性
该 nginx 存在通配规则：**凡以 `/login` 结尾的路径一律交给 PHP，PHP 统一回 `header('Location: ../example.com')`（相对路径）**。
```
/upload/login        → 302  Location: ../example.com
/abc/login           → 302  Location: ../example.com   ← 任意词
/foobar/login        → 302  Location: ../example.com   ← 任意词
/dtjs/zzqqxx/login   → 302  Location: ../example.com   ← 任意词
/upload/admin/login  → 302  Location: ../example.com
/dtjs/login2         → 404   ← 破坏形态即无
/dtjs/alogin         → 404
/upload/loginx       → 404
```
⇒ 扫目录时看到 `/<任意目录>/login` 命中 302，**不要当成"发现隐藏登录口"**。真登录文件只有 `/agent/login.php`、`/uxj/login.php`（二者同样被守卫拦）。
**注意**：`Location: ../example.com` 是**相对路径**，会被客户端按"请求所在目录"解析 —— `/upload/login` 解析成 `/example.com`、`/upload/agent/login` 解析成 `/upload/example.com`，**看起来目标不同，其实是同一条规则**。判读时务必打印**原始 Location 头**，不要只看 `curl -w %{redirect_url}`。

## 二、203.0.113.38 / .23 —— 两个可直接用的端点

### `/code` = 全站唯一未受守卫保护的动态 PHP 端点（验证码生成器）
```
GET /code  → 200  Content-Type: image/jpeg  X-Powered-By: PHP/7.4.33
3 次取样 1479B / 1435B / 1420B，MD5 各不相同 ⇒ 动态生成
图像内容实测（read_image 打开看过）：80×30 的 4 位数字验证码
```
- 所有 HTTP 方法（GET/POST/HEAD/OPTIONS/PUT）都出图；
- 读 `t` 参数（传 `?t=<任意值>` ⇒ 200 空 body；`type`/`x`/`file`/`id`/`code` 均不影响出图）；
- **已证非注入点**：`t` 喂 LFI payload（相对/绝对路径、`php://filter`、Windows 反斜杠、`%00`）、数组参数 `t[]=`、以及**提交真实验证码值** ⇒ 响应全部一致为空。
- **⭐ 实战用法**：**若任何战线拿到本站凭据要爆破登录，必须先来此处取验证码并 OCR，否则每个候选口令都会因验证码缺失而失败**（本队曾在 ng-api 后台靠 OCR 解决同类问题）。

### 免认证可达面（low，不计分）
`/agent/example.com`(200/370B)、`/uxj/example.com`(200/369B) 为唯二不重定向的业务 PHP 页（静态 frameset 骨架，故未 include 守卫），泄露内部页 `/agent/top.php?xtype=this`、`/uxj/xy.php` —— 这两页**确实存在但均被守卫拦 302**，`xtype` 取值无绕过。

## 三、已穷尽、不要再投入的面（7000+ 请求全空）
| 面 | 结论 |
|---|---|
| `/data/`（两轮 2640 请求） | 完整清单仅 4 文件（db.php 282B、var.php 145B、config.inc.php、session.php），**全部经 MD5 级 A/B 对照确认不读任何请求参数 ⇒ 无 LFI 注入点** |
| `/upload/`（1172 请求） | 空目录（不是被禁，是确实没有文件） |
| `/dtjs/`（1148 请求） | 仅正常前端资源（jsglue.js/base.js/common.js/jquery.min.js） |
| 站点根目录 + `/agent/` + `/uxj/`（3444 请求） | 无 webshell、无前人后门 |
| 8888 宝塔登录 | 凭据攻击已由专职战线做过 8 次全败；`/static/*` 认证面攻坚 12+20 变体全败 |

## 四、已排除
443 业务入口全被守卫拦（域名全 NXDOMAIN 所致，Host 头无效）；21 FTP FileZilla 0.9.46 beta（匿名 530、`AUTH TLS` 502 未启用、14147 不通）；3306 MySQL IP 白名单（`1130 Host not allowed`）；3389 不打；888/80 空站；URL 路径穿越 11 变体被 nginx 规范化拦死；10+ Host 变体无分流。

## 五、待探（交接）
`203.0.113.72` 的 `2103/2105/2107/2179` 四端口（与 25/110 同属一个「IP+MAC 白名单准入代理网关」，返回 `Server: Proxy` + `Auth Result: 无效用户.` 并回显客户端真实 IP 与期望 MAC）；`2122` FTP 已试 22 组口令全 530，勿重复爆破。