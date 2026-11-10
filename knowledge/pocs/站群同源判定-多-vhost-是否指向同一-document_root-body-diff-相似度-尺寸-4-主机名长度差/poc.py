## 方法（3 步，通常只需 2 发包）

### 1) 取两个不同 vhost 的完整响应体
```python
import ssl, urllib.request, difflib
CTX = ssl.create_default_context(); CTX.check_hostname=False; CTX.verify_mode=ssl.CERT_NONE
OP = urllib.request.build_opener(
    urllib.request.ProxyHandler({"http":"http://127.0.0.1:7892","https":"http://127.0.0.1:7892"}),
    urllib.request.HTTPSHandler(context=CTX))
UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/203.0.113.20 Safari/537.36"
def get(u):
    r = OP.open(urllib.request.Request(u, headers={"User-Agent":UA,"Accept":"text/html,*/*;q=0.8"}), timeout=25)
    return r.read()
a, b = get("https://www.EXAMPLE-DOMAIN/bbs/"), get("https://sub.EXAMPLE-DOMAIN/bbs/")
print(len(a), len(b), difflib.SequenceMatcher(None, a.decode("gbk","replace"), b.decode("gbk","replace")).ratio())
for l in difflib.unified_diff(a.decode("gbk","replace").splitlines(),
                              b.decode("gbk","replace").splitlines(), lineterm="", n=1):
    print(l[:200])
```

### 2) 三档判据（按强度递增）
| 现象 | 判定 |
|---|---|
| 相似度 ≈ 0.999x，差异行**全部是主机名字符串**（绝对链接/JS 跳转/archives/rss） | **同一 DOCUMENT_ROOT、同一数据库** = 同一个站，域名只是别名 |
| 尺寸差 = **4 × 主机名长度差** | 同一份文件（首页里主机名出现 4 次）⇒ 数学证据，比"看起来一样"硬 |
| 版块/帖子数/在线人数等**动态数据不同** | 不同实例（不同库） |

### 3) 战术含义（这条判据存在的意义）
- 站群横向（想在同源站群里找"没加固的实例"）**成立的前提是各 vhost 是独立部署**。用本判据先证伪，可省下整批额度。
- 本次实测：6 个 vhost 全部 `generator=7.3.2` + Cookie 前缀同一个 `752fb_`，其中两个 vhost 首页相似度 **0.9998**，19 行差异全是主机名，尺寸 67528 vs 67544（16 = 4×4）；主目标 67556 与 67528 差 **28 = 4×7** ⇒ 主目标与全族是同一个站 ⇒ **站群别名横向 = 0 收益**，加固是"部署模板级"而非逐站手工。
- 辅助判据：Cookie 前缀（PHPWind 的 `xxxx_` 前缀）、`<meta name=generator>` 版本、首页 `title`，任一项不同才说明是**另一套部署**，才值得单独立案。

### 4) 注意
- GBK 站用 `decode("gbk","replace")` 再 diff，否则编码错误会污染相似度。
- 一定要带正常浏览器 UA + Referer（本站首页），并按 1 req/2s 限速。