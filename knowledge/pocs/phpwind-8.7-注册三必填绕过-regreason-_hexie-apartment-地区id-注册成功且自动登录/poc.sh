## 背景
phpwind 8.7，配置为「新会员需管理员验证」(rg_ifcheck=1)。注册提交后只返回通用
「提示信息」页，不含字段名，定位困难。

## 三个必填项（逐一实测确定）
| 顺序 | 缺失时提示 | 来源字段 | 说明 |
|---|---|---|---|
| 1 | `请填写注册原因。` | `regreason` | `target134.example.com:105 setReason()`：`L::reg('rg_ifcheck')` 为真时强制 |
| 2 | `字段 '现居住地' 填写有误` | `apartment` | 自定义 area 类型字段，值为地区数字 ID |
| 3 | 通用反爬提示 | `_hexie` | 注册页隐藏域，值形如 `cn0zz`；**HTML 中 `value` 在 `name` 之前**，正则写反会取到空值 |

## 实测请求（成功）
```bash
B=http://203.0.113.6:1852
curl -s -c ck.txt "$B/register.php" -o reg.html
# 注意属性顺序：<input type="hidden" value="cn0zz" name="_hexie" />
HEXIE=$(grep -o 'value="[^"]*" name="_hexie"' reg.html | sed 's/value="//;s/" name="_hexie"//')

curl -s -b ck.txt -c ck.txt -e "$B/register.php" \
  --data-urlencode "forward=" --data-urlencode "step=2" --data-urlencode "_hexie=$HEXIE" \
  --data-urlencode "regname=kwx108345" --data-urlencode "regpwd=Abc@9911xy" \
  --data-urlencode "regpwdrepeat=Abc@9911xy" --data-urlencode "regemail=kwx108345@163.com" \
  --data-urlencode "rgpermit=1" --data-urlencode "regreason=测试" \
  --data-urlencode "apartment=820001" "$B/register.php?"
```
结果：`HTTP 302`，cookie 出现 `winduser`（= `target133.example.com:298 if ($rgyz == 1) Cookie("winduser",...)`）
⇒ 同时证明 `rg_emailcheck=0`（`_setEmailCheck()` 令 `yz=1`），账号无需激活邮件即可登录。

## apartment 合法值（从注册页内联 `var areas` 数组提取）
```
areas['820001']=['澳门','820000','0'];
areas['820002']=['澳门特别行政区','820001','0'];
areas['110101']=['东城区','110100','1'];
areas['810001']=['香港','810000','1'];
```
服务端校验链：`target106.example.com:149 checkCustomerField` → `target135.example.com:99 checkData`
→ `target136.example.com:245 checkData` → `type==area` 分支 → `checkArea($data)`。

## 重要限制（不要误判）
**注册成功 ≠ 能发帖。** `_setGroupid()`（`target134.example.com:285-289`）：
```php
if ($this->data['groupid'] == 0)
    $this->data['groupid'] = L::reg('rg_ifcheck') ? 7 : -1;
```
目标 `rg_ifcheck=1` ⇒ 新账号 `groupid=7`（验证会员）。实测 `post.php?fid=2` 返回
「您还没通过管理员验证…」提示页。`groupid==7` 只拦发帖（`PwPost::postcheck()` /
`banUser()`），不拦登录、不拦头像上传、不拦消息页。

## 通用性
换站只需替换 TARGET 与 `apartment` 的一个合法地区 ID（同站群其它 phpwind 8.7 站点应通用，
因 `areas` 数据由后台地区库生成，取值在注册页内联可读）。