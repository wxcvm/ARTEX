# PHPWind 站点精确版本与表前缀被动识别

## 版本号来源（无需认证）
```bash
curl -sk https://TARGET/bbs/ | grep -oP '<meta name="generator" content="\K[^"]+'
# → 7.3.2
```

## 表前缀泄露（无需认证）
任意一次请求的 Set-Cookie 里带 `<prefix>_lastvisit` / `<prefix>_AdminUser`：
```bash
curl -skI https://TARGET/bbs/ | grep -i 'set-cookie'
# → 604b2_lastvisit=...; 604b2_AdminUser=deleted
# 前缀即数据库表前缀：604b2_
```

## 后台入口探测（低噪声）
```bash
curl -sk -o /dev/null -w '%{http_code}\n' https://TARGET/bbs/example.com   # 200 = 存在
```
后台登录表单特征：`action="/bbs/example.com?&verify=<8位hex>&"`，
字段 `Login_f` / `admin_name` / `admin_pwd`，失败提示「密码或安全问题错误,您还可以尝试7次」。

## PHPWind 独有路径指纹（用于同源资产反查）
- `style/txcss.css`（PHPWind 默认样式）
- `js/pw_ajax.js`
- `/simple/`（静态归档）
- `/rss.php?fid=`
- 页面变量 `verifyhash = '<8位hex>'`
- `robots.txt` 内容含 `PHPWIND BOARD`

## 定位同类站点
用上述任一特征做 FOFA/搜索引擎语法，即可批量发现同源 PHPWind 站点：
`body="pw_ajax.js" && body="txcss.css"` / `body="/bbs/simple" && title="论坛"`