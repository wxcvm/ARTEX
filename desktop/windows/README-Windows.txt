ARTEX Windows 桌面客户端（便携版）
================================

运行方法：
1. 将此压缩包完整解压到可写目录（例如 D:\Apps\ARTEX）。
2. 如未设置 ARTEX_PG_DSN，将 config.example.json 复制为 config.json，
   按照 docs\USER_GUIDE_ZH.md 配置已有 PostgreSQL。
3. 双击 ARTEX.Desktop.exe。首次启动会尝试连接现有 ARTEX，
   或自动启动同目录的 artex.exe（默认 127.0.0.1:8787）。
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
- 桌面窗口不会安装或管理 PostgreSQL；数据库仍需按原项目方式准备。
- 不要只复制 ARTEX.Desktop.exe，桌面包中的 artex.exe、skills 及配置
  模板应一并保留。

完整中文操作指南：docs\USER_GUIDE_ZH.md
仓库：https://github.com/wxcvm/ARTEX
