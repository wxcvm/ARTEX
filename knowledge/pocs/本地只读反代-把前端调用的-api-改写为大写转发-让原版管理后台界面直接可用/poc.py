【适用前提】目标面板满足：GET /api/admin/xxx → 401，而 GET /API/admin/xxx → 200（Express 等框架的 caseSensitive 疏漏）。
注意：**该疏漏不可跨产品推广**——同一运营方的另一套自研面板（AIRouter）实测 /API/ 与小写 /api/ 行为完全相同（同为 302），必须逐个实测。

【六个关键实现点（缺一个就白屏或跳出）】
1. 路径改写：请求 path 以 /api/ 开头 → 改为 "/API" + path[4:] 再转发到目标。
2. 入口映射：本地 /、/example.com、/admin、/admin.html、/login.html 一律映射到**真后台页**。
   本例真后台页是 /admin%2Ehtml（上游 /admin.html 会 302 到登录页）；用 %2E 这种编码地址当入口在浏览器里容易被自动解码，映射成简单路径最稳。
3. 302 Location 改写：上游返回的 Location 若指向目标域或 /api/，同步改回本地/大写，否则浏览器会跳出反代（跳回真站或被弹到登录页）。
4. 写操作硬挡：do_POST / do_PUT / do_DELETE / do_PATCH 一律返回 403，**在服务端实现**（不给前端任何绕过余地）——避免误改生产配置，也让"浏览"变成纯只读。
5. 禁用代理：urllib 要显式 build_opener(urllib.request.ProxyHandler({}))，否则会被本机 HTTP_PROXY 环境变量带偏导致转发失败。
6. 探针 + 启动顺序：提供 /health（本地实现，不转发）；启动脚本必须"先起服务 → 轮询 /health 就绪 → 再打开浏览器"，**顺序反了必然白屏**（常见现成脚本的 bug 就是这个）。

【其它细节】
- 用 ThreadingMixIn + allow_reuse_address，避免并发请求被单线程阻塞。
- 转发时剔除 Host 与 hop-by-hop 头，重设 Host 为目标。
- 浏览器若挂了代理插件，必须把 127.0.0.1 加入直连白名单，否则页面请求被代理截住。
- 验证清单：curl -o NUL -w "%{http_code} %{size_download}" 逐个确认 ① 入口页返回后台页字节数 ② /api/admin/channels 返回全量数据 ③ 资源（css/js）全 200 ④ POST 返回 403。

【本次实测数据（203.0.113.105:13030）】
- 入口 /、/admin、/login.html、/admin%2Ehtml → 全部 200 / 18899 B 真后台页（菜单 7 项：dashboard, channels, channelOps, keys, logs, gens, settings）
- 前端 JS 调用的 /api/admin/channels → 200 / 310689 B（105 渠道含明文 key）
- /api/admin/stats、/api/channels/status、/api/admin/channel-types → 全 200
- POST /api/admin/channels → 403（本地挡下）

【交付物（本次实现）】交付物/17-一键查看器/{server.py, example.com, 一键查看.bat}