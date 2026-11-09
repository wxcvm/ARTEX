ARTEX Windows 桌面客户端（便携版）
================================

运行方法：
1. 将此压缩包完整解压到可写目录（例如 D:\Apps\ARTEX）。
2. 如未设置 ARTEX_PG_DSN，将 config.example.json 复制为 config.json，
   按照 docs\USER_GUIDE_ZH.md 配置已有 PostgreSQL。
3. 双击 ARTEX.Desktop.exe。首次启动会尝试连接现有 ARTEX，
   或自动启动同目录的 artex.exe（默认 127.0.0.1:8787）。
   若没有 config.json 且未设置 ARTEX_PG_DSN，会自动弹出数据库连接向导；
   保存后继续启动（不会覆盖已有配置）。
   向导会检测本机 PostgreSQL 的 5432 / 5433 端口，并列出已安装的
   PostgreSQL Windows 服务；可明确选择启动某一服务后重新检测。
   端口探测不验证账号密码；可点击“测试数据库连接”使用内置 Go 驱动
   核验当前填写的连接信息（不会保存或修改数据库数据）。
   可明确确认后创建缺失的目标数据库（需账号有 CREATEDB 权限，不覆盖已有库）。
   配置填错时点击窗口工具栏“数据库配置”重新打开向导，其他 JSON 设置会保留。
   如果尚未安装 PostgreSQL，点击“安装 PostgreSQL”：
   有 WinGet 时启动 EDB PostgreSQL 16 官方交互安装（需联网、管理员许可），
   没有 WinGet 时打开 PostgreSQL 官方 Windows 下载页。
   按安装器提示设置 postgres 超级用户密码，回到向导点击“检测本机”；
   使用用户名 postgres 和刚设置的密码，可“创建目标数据库”再验证并保存。
4. 完成首次管理员设置和模型配置，即可在独立桌面窗口操作。

使用说明：
- Windows 11 通常已包含 Microsoft Edge WebView2 Runtime。
  若缺少该运行时，请从微软官方获取并安装。
- 桌面版基于 .NET 10，随包发布的 ARTEX.Desktop.exe 为自包含版本，
  无须再安装单独的 .NET Desktop Runtime。
- 点击窗口右上角 X 或最小化会将程序收进托盘；
  通过托盘或窗口内“退出”真正结束客户端。
- 若桌面客户端自行启动了后端，退出时也会停止它；
  如果连接的是已经运行的后端，退出不会停止既有服务。
  正在执行的任务可能因后端停止而中断。
- WebView2 的浏览器数据及后端日志位于：
  %LOCALAPPDATA%\ARTEX\
- 默认 Web/API 仍可通过 http://127.0.0.1:8787 访问。
- “打开配置目录”打开程序所在目录，便于修改数据库配置；
  “在浏览器打开”继续保留传统网页使用方式。
- 桌面窗口提供可选的 PostgreSQL 官方安装入口；安装程序独立运行，
  不内置数据库、默认密码或自动更改数据库服务。
  可明确选择创建缺失的目标数据库；数据库角色需事先存在且具备权限。
  管理服务可能需要 Windows 管理员权限；原有数据库和账户不会被重置。
- 不要只复制 ARTEX.Desktop.exe，桌面包中的 artex.exe、skills 及配置
  模板应一并保留。

完整中文操作指南：docs\USER_GUIDE_ZH.md
仓库：https://github.com/wxcvm/ARTEX
