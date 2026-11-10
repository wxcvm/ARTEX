# 全量参数钥匙（Parameter-Key Batch）—— 打开老 JSP/servlet 的通用方法

## 原理
2015 年前后的 Java servlet（如本靶标的 `example.com.Service.*` / `example.com.Service_2015.*` / `example.com.Service2_8.*` / `example.com.Manager.*`）
到处写着：

```java
int userid = Integer.parseInt(request.getParameter("userid"));   // 缺参 → NumberFormatException: null
int start  = Integer.parseInt(request.getParameter("start"));
...
```

- 缺参数时**不是返回错误码，而是抛 `NumberFormatException: null` → HTTP 500 + 完整 Java 栈**（栈里直接给出**类名 + 行号**，等于白送一份 servlet 地图）。
- 猜参数名是线性成本；**一次性把一大批候选参数名全带上（都赋一个能被 parseInt 接受的数字 `1`）**，即可一次性越过所有 int 门槛，直接命中后面的数据分支。
- 由于 Tomcat 的 `getParameter(name)` 返回**第一次出现**的值，可以用「注入参数放前面 + 钥匙放后面（值=1）」的写法，对**同名参数实现覆写**（用于逐参数注入测试）。

## 可复用调用
```python
import socket
KEY = "userid=1&userId=1&uid=1&id=1&...&qishu_id=1"   # 296~2000 个候选名，全部 =1
req = ("POST /LHService/ZiLiao HTTP/1.1\r\nHost: <IP>\r\n"
       "Content-Type: application/x-www-form-urlencoded\r\n"
       "Content-Length: %d\r\nConnection: close\r\n\r\n%s" % (len(KEY), KEY))
# 注意：一律用原生 socket，不要用 curl（本机 http_proxy 会污染结论）
```

## 候选参数名来源（按性价比排序）
1. **从目标 App 的 APK dex 里抽全部标识符**（`re.findall(rb"[\x20-\x7e]{5,}", dex)` → 过滤 `^[A-Za-z_][A-Za-z0-9_]{1,19}$`），再补中文拼音与老 JSP 常见名 —— 本次由 65k 行 dex 字符串抽出 **34,491** 个候选名。
2. 补**单/双/三字符短名**（`a..z0-9_` 的 1~3 位组合，约 4,406 个）——本次 `userid`/`token`/`type`/`start`/`ip` 这类短名都被这一层覆盖到。
3. **分组探测**：每批 2000 个（避免 Tomcat `maxParameterCount=10000` 上限），比对响应差异；命中即**二分收敛**到具体参数名（约 log2(2000)≈11 次请求）。
4. 二分判据：**响应字节与基线不同 且 不再是 `NumberFormatException`**。

## 本次实战收益（同一把钥匙连开 7 个端点）
| 端点 | 钥匙前 | 钥匙后 |
|---|---|---|
| `/LHService/userlist` | 500 NFE | **200，5 条真实用户记录（40+ 字段）** |
| `/LHService/ZiLiao` | 500 NFE | **200，资料库列表 JSON** |
| `/LHService/ZiLiaoBody` | — | **200，文章全文** |
| `/LHService/pingbi` | 500 NFE（`PingbiServlet:44`） | **200，屏蔽词/博彩域名黑名单** |
| `/LHService/userinfo` | 200 `-1`（未登录） | **200 + `userid=<任意>` ⇒ IDOR 任意用户 PII** |
| `/LHService/userxunzhang` / `zj_info` | 500 NFE | **200，按 userid 返回勋章统计 / 专家推荐列表** |
| `/Tk_Kj/init`、`/LHService/init` | 200 空（LHService 甚至挂 45s） | **200 `ok`（参数 `ip`，int）** |

## 顺带产出的"servlet 地图"（靠 500 栈里的类名）
`example.com.Service.UploadPicServlet`、`example.com.Service.ZiLiaoServlet`、`example.com.Service_2015.Get_user_list`、
`example.com.Service_2015.UserXunZhangServlet`、`example.com.Service_2015.PingbiServlet`、`example.com.Service_2015.ZoushiServlet`、
`example.com.Service_2015.JubaoServlet`、`example.com.Service.PinglunServlet`、`example.com.Service.LiuyanServlet`、
`example.com.Service.GetDianpinListServlet`、`example.com.Service.GetUser_tiezi`、`example.com.Service2_8.GetTk`、
`example.com.Service2_8.GetKJ2_8`、`example.com.Service2_8.GetTj2`、`example.com.Service2_8.GetRenwu`、
`example.com.Service2_8.getRenwuPaihang`、`example.com.Service2_8.Zj_Dingyue`、`example.com.Manager.ConfigServlet`；
过滤器：`example.com.Manager.DeviceSecurityFilter`、`example.com.Service_2015.MoneyTransactionFilter`、`example.com.cache.v27.LegacyForumWrites`。

## 局限
- 只对**老式 `Integer.parseInt(getParameter(...))` 型 servlet** 有效。
- 若参数名不在候选集里（或服务端读的是**配置文件/Properties**而非请求参数，如本靶标 `target335.example.com:24`——34,491+4,406 个名全试仍 `NumberFormatException: null`，属"配置资源缺失"故障），则无效。