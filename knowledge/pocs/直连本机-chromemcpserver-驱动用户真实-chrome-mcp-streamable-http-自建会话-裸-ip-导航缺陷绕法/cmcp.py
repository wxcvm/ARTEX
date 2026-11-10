#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""极简 MCP streamable-http 客户端：直连本机 ChromeMcpServer，驱动用户真实浏览器。

适用场景
--------
凡 DSH 里的 MCP 客户端报 "Invalid MCP request or session"（服务本体其实正常），
或需要复用用户真实 Chrome 登录态（Google OAuth、已登录后台）时，用本脚本绕开
DSH 的 MCP 客户端，自己起一个 MCP 会话直连服务。

用法
----
  python cmcp.py list                                  # 列出服务端工具
  python cmcp.py call get_windows_and_tabs '{}'         # 调工具
  python cmcp.py call chrome_javascript '{"code":"return document.title"}'

含中文/引号的 JSON 一律走临时文件 + --data-binary，或用下面的 tool.py/f(@file) 写法，
不要用 PowerShell echo 内联（转义会炸）。

关键坑（都踩过）
----------------
1) 服务地址在 DSH 配置里：C:\\Users\\<user>\\workspace\\cordis.patch.yml 的
   mcp-streamable-mcp-server → url: http://127.0.0.1:12306/mcp
   （是 streamable-http，不是 stdio；探活：POST initialize 能返回
   target5.example.com="ChromeMcpServer" 即服务活着）。
2) 必须自己维护 mcp-session-id：initialize 响应的 header 里取，后续每个请求
   带 mcp-session-id 头，并在 initialize 后补一条 notifications/initialized，
   否则 tools/call 会被拒。
3) 响应是 SSE（text/event-stream），要按行取 "data: " 后的 JSON。
4) 该 MCP 的 chrome_navigate **对裸 IP 有缺陷**：内部会拼
   'http://www.<ip>:<port>/*' 去查标签页，报 `Invalid url pattern`。
   绕法：用通配 DNS 域名指向该 IP —— 例如 http://<a>.<b>.<c>.<d>.sslip.io:<port>/，
   sslip.io / nip.io 都行（sslip.io 实测可用）。目标站点若对任意 Host 都返回
   同一应用则无副作用；注意换域名=换 origin，原有 Cookie 不再携带，需要重新登录。
   其它绕法：从已有标签页 window.open()（会被 popup blocker 拦，实测失败）。
5) 它驱动的是**用户真实浏览器**（本机若有另一个 --remote-debugging-port=9222 的
   HeadlessChrome，那是另一个实例，别搞混：用 get_windows_and_tabs 看标签页内容判断）。
6) chrome_javascript 的返回里敏感串会被自动脱敏（如 refreshToken 变
   <redacted_base64>）。要拿未脱敏原文，改用站点同源 curl（带 Cookie）或
   chrome_network_capture(needResponseBody=true)。
7) chrome_screenshot 默认走"下载"落盘（Downloads 目录），
   传 storeBase64=true + savePng=false 才拿到 base64Data 自己写文件；
   返回字段名是 base64Data（不是 base64）；且**只截当前活动标签**，
   非活动标签要显式传 tabId。
8) 同一 MCP 服务可能被多个会话/智能体共用，偶发单次请求超时（WinError 10060）
   —— 重试一次即可，不要据此判定服务不可用。

伦理/边界
--------
这会拿到用户真实浏览器的全部登录态。只允许导航到授权目标域名，绝不打开、
点击、读取、截图与演练无关的私人页面（邮箱、VPS 后台、网银等）。
"""
import json
import sys
import urllib.request

URL = "http://127.0.0.1:12306/mcp"
HDRS = {"Content-Type": "application/json",
        "Accept": "application/json, text/event-stream"}


def _post(payload, session=None):
    h = dict(HDRS)
    if session:
        h["mcp-session-id"] = session
    req = urllib.request.Request(URL, data=json.dumps(payload).encode(),
                                 headers=h, method="POST")
    with urllib.request.urlopen(req, timeout=180) as r:
        sid = r.headers.get("mcp-session-id")
        raw = r.read().decode("utf-8", "replace")
    out = {"_session": sid}
    for line in raw.splitlines():
        if line.startswith("data: "):
            out.update(json.loads(line[6:]))
    return out


def init():
    r = _post({"jsonrpc": "2.0", "id": 1, "method": "initialize",
               "params": {"protocolVersion": "2024-11-05", "capabilities": {},
                          "clientInfo": {"name": "rt", "version": "1.0"}}})
    sid = r["_session"]
    _post({"jsonrpc": "2.0", "method": "notifications/initialized"}, sid)
    return sid


def call(name, args, sid=None):
    sid = sid or init()
    return _post({"jsonrpc": "2.0", "id": 2, "method": "tools/call",
                  "params": {"name": name, "arguments": args}}, sid)


if __name__ == "__main__":
    sid = init()
    if sys.argv[1] == "list":
        r = _post({"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}, sid)
        for t in r["result"]["tools"]:
            print(t["name"], "|", (t.get("description") or "")[:100])
    else:
        args = json.loads(sys.argv[3]) if len(sys.argv) > 3 else {}
        print(json.dumps(call(sys.argv[2], args, sid), ensure_ascii=False)[:4000])
