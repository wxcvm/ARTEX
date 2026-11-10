## 现象
同一 HTTP payload：
- `curl.exe -sS --max-time 45 -H "Host: example.com" http://<ip>/...` → **HTTP 200，0.14s**
- Python `urllib.request.urlopen(req, timeout=25)` → **TimeoutError: timed out（重试 3 次全失败）**

导致误判为"目标限速/封禁"，实际是**本机 Python 走了本地代理**。

## 根因
Windows 当前用户注册表配了 IE 代理：
```
HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings
  ProxyEnable = 1
  ProxyServer = 127.0.0.1:7892
```
- **Python 的 `urllib.request` 会自动读取注册表代理**（`getproxies()` 返回 `{'http': 'http://127.0.0.1:7892', ...}`），即使环境变量 `HTTP_PROXY` 为空
- **`curl` 只读环境变量，不读注册表** → 直连，所以正常
- 该本地代理对目标请求超时/极不稳定 → 表现为"目标不可达"

## 修复（任选其一）
1. **Python 显式禁用代理（推荐，改代码）**
```python
import urllib.request
OPENER = urllib.request.build_opener(urllib.request.ProxyHandler({}))
with OPENER.open(req, timeout=25) as resp:
    body = resp.read()
```
2. 给 URL 加 `no_proxy`：`urllib.request.proxy_bypass` 或直接 `os.environ["no_proxy"] = "*"`
3. 改用 `curl.exe`（子进程调用），它天然不吃注册表代理

## 判据（避免误判目标封禁）
- **curl 通 + Python 不通** ⇒ 99% 是本机代理问题，不是目标封禁
- 诊断一行：`python -c "import urllib.request; print(urllib.request.getproxies())"`
- 出现 `127.0.0.1:7892` 之类即命中

## 附带坑（同类）
- `urllib.request.Request` 对象**不可复用**：一次 `urlopen` 后重试必须**新建 Request**，否则表现为挂起/超时。写成 `for _ in range(3): r = Request(...); urlopen(r)`。
- 用 IP + `Host:` 头访问虚拟主机时，务必加 `Connection: close`，避免长连接被服务端半开导致等待。

## 实测环境
Windows 11 / Python 3.13.15（`D:\Python313\python.exe`）/ curl.exe（系统自带）；目标为境外 IPv4 虚拟主机。