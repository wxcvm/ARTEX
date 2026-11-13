using System.Diagnostics;
using System.Net.Sockets;
using System.Drawing;
using System.Text.Json;
using System.Windows.Forms;
using Microsoft.Web.WebView2.Core;
using Microsoft.Web.WebView2.WinForms;

namespace Artex.Desktop;

public sealed class MainForm : Form
{
    private static readonly Uri Service = new("http://127.0.0.1:8787/");
    private WebView2 browser = new() { Dock = DockStyle.Fill };
    private readonly Panel viewHost = new() { Dock = DockStyle.Fill };
    private readonly StartupView pageState = new();
    private readonly System.Windows.Forms.Timer pageLoadWatchdog = new() { Interval = 25000 };
    private readonly ToolStrip toolbar = new() { Dock = DockStyle.Top, GripStyle = ToolStripGripStyle.Hidden };
    private readonly StatusStrip statusBar = new() { Dock = DockStyle.Bottom };
    private readonly ToolStripButton fullScreenButton = new("全屏 (F11)");
    private readonly ToolStripStatusLabel status = new("正在启动 ARTEX…");
    private readonly HttpClient http = new() { Timeout = TimeSpan.FromSeconds(2) };
    private readonly CancellationTokenSource lifetime = new();
    private readonly NotifyIcon tray;
    private readonly object logLock = new();
    private Process? ownedBackend;
    private string? ownedSetupToken;
    private StreamWriter? backendLog;
    private long backendLogStartOffset;
    private bool exiting;
    private bool ready;
    private bool pageResourceFailed;
    private bool browserRecreationNeeded;
    private bool rendererRecoveryAttempted;
    private bool fullScreen;
    private FormWindowState restoreWindowState = FormWindowState.Maximized;
    private FormWindowState beforeFullScreenState;
    private Rectangle beforeFullScreenBounds;
    private int connecting;

    private static string LocalData =>
        Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData), "ARTEX");

    private static string LogPath => Path.Combine(LocalData, "desktop-backend.log");

    public MainForm()
    {
        Text = "ARTEX";
        Icon = System.Drawing.Icon.ExtractAssociatedIcon(Application.ExecutablePath) ?? SystemIcons.Application;
        MinimumSize = new Size(900, 620);
        Size = new Size(1360, 850);
        StartPosition = FormStartPosition.CenterScreen;
        WindowState = FormWindowState.Maximized;
        pageState.RetryRequested += (_, _) => RetryConnection();
        pageState.RepairRequested += (_, _) => EditDatabaseConfiguration();
        pageState.LogsRequested += (_, _) => OpenBackendLog();
        toolbar.BackColor = Color.FromArgb(20, 32, 49);
        toolbar.ForeColor = Color.FromArgb(232, 239, 250);
        toolbar.Padding = new Padding(12, 7, 12, 7);
        toolbar.Font = new Font("Microsoft YaHei UI", 9.5f);
        statusBar.BackColor = Color.FromArgb(19, 29, 45);
        statusBar.ForeColor = Color.FromArgb(171, 192, 218);
        status.ForeColor = Color.FromArgb(171, 192, 218);
        pageLoadWatchdog.Tick += (_, _) =>
        {
            pageLoadWatchdog.Stop();
            if (ready || lifetime.IsCancellationRequested) return;
            status.Text = "页面加载超时 · 可重新连接";
            LogDesktopFailure("WebView2 页面加载超过 25 秒");
            ShowPageMessage("界面加载超时", "网页未能按时加载，请点击“重新连接”重试。", true);
        };
        viewHost.Controls.Add(browser);
        viewHost.Controls.Add(pageState);

        toolbar.Items.Add(Button("后退", () =>
        {
            if (browser.CoreWebView2?.CanGoBack == true) browser.CoreWebView2.GoBack();
        }));
        toolbar.Items.Add(Button("前进", () =>
        {
            if (browser.CoreWebView2?.CanGoForward == true) browser.CoreWebView2.GoForward();
        }));
        toolbar.Items.Add(Button("刷新", RetryConnection));
        fullScreenButton.Click += (_, _) => ToggleFullScreen();
        toolbar.Items.Add(fullScreenButton);
        toolbar.Items.Add(new ToolStripSeparator());
        toolbar.Items.Add(Button("在浏览器打开", OpenInBrowser));

        // Keep common navigation visible at high DPI; secondary operations stay
        // reachable under one menu instead of crowding the entire title toolbar.
        var more = new ToolStripDropDownButton("更多")
        {
            ForeColor = Color.FromArgb(224, 235, 249),
            Margin = new Padding(2, 0, 2, 0),
            Padding = new Padding(5, 3, 5, 3),
        };
        more.DropDownItems.Add("数据库配置", null, (_, _) => EditDatabaseConfiguration());
        more.DropDownItems.Add("打开配置目录", null, (_, _) =>
            Process.Start(new ProcessStartInfo("explorer.exe", "\"" + AppContext.BaseDirectory.TrimEnd('\\') + "\"")
            { UseShellExecute = true }));
        more.DropDownItems.Add(new ToolStripSeparator());
        more.DropDownItems.Add("重新连接", null, (_, _) => RetryConnection());
        more.DropDownItems.Add("退出 ARTEX", null, (_, _) => ExitDesktop());
        toolbar.Items.Add(more);

        statusBar.Items.Add(status);

        Controls.Add(viewHost);
        Controls.Add(toolbar);
        Controls.Add(statusBar);

        var trayMenu = new ContextMenuStrip();
        trayMenu.Items.Add("显示 ARTEX", null, (_, _) => ShowDesktop());
        trayMenu.Items.Add("在浏览器打开", null, (_, _) => OpenInBrowser());
        trayMenu.Items.Add("退出 ARTEX", null, (_, _) => ExitDesktop());
        tray = new NotifyIcon
        {
            Icon = Icon,
            Text = "ARTEX Desktop",
            Visible = true,
            ContextMenuStrip = trayMenu,
        };
        tray.DoubleClick += (_, _) => ShowDesktop();
    }

    private static ToolStripButton Button(string label, Action click)
    {
        var button = new ToolStripButton(label)
        {
            ForeColor = Color.FromArgb(224, 235, 249),
            Margin = new Padding(2, 0, 2, 0),
            Padding = new Padding(5, 3, 5, 3),
        };
        button.Click += (_, _) => click();
        return button;
    }

    protected override void OnHandleCreated(EventArgs e)
    {
        base.OnHandleCreated(e);
        DesktopChrome.UseDarkTitleBar(this);
    }

    protected override async void OnShown(EventArgs e)
    {
        base.OnShown(e);
        await ConnectAsync(lifetime.Token);
    }

    private async Task ConnectAsync(CancellationToken cancel)
    {
        if (cancel.IsCancellationRequested || Interlocked.Exchange(ref connecting, 1) != 0) return;

        try
        {
            ShowPageMessage("正在连接 ARTEX…", "正在检查本机服务", false);
            status.Text = "正在连接本机 ARTEX 服务…";
            if (!await IsHealthyAsync(cancel))
            {
                if (ownedBackend is { HasExited: true })
                {
                    ownedBackend.Dispose();
                    ownedBackend = null;
                }

                if (ownedBackend is null)
                {
                    if (await IsPortOccupiedAsync(cancel))
                        throw new InvalidOperationException(
                            "本机 8787 端口已被其他服务占用，但它不是 ARTEX。请检查占用端口的程序后重试。");

                    if (DatabaseSetupForm.NeedsConfiguration || !await DatabaseSetupForm.CheckConfiguredAsync(cancel))
                    {
                        using var setup = new DatabaseSetupForm(editExisting: File.Exists(DatabaseSetupForm.ConfigPath));
                        if (setup.ShowDialog(this) != DialogResult.OK)
                        {
                            status.Text = "未配置数据库 · 点击“重新连接”继续设置";
                            ShowPageMessage("完成数据库设置", "连接到 PostgreSQL 后才能进入 ARTEX 工作台。", true,
                databaseProblem: true, nextStep: "点击“修复数据库连接”，按照引导完成安装或填写现有账户。");
                            return;
                        }
                    }

                    StartBackend();
                }

                status.Text = "正在等待本机 ARTEX 后端就绪…";
                var deadline = DateTime.UtcNow.AddSeconds(60);
                while (!await IsHealthyAsync(cancel))
                {
                    if (ownedBackend?.HasExited == true)
                        throw new InvalidOperationException(
                            $"ARTEX 后端已退出（代码 {ownedBackend.ExitCode}）。请检查 PostgreSQL 配置和 {LogPath}");
                    if (DateTime.UtcNow >= deadline)
                        throw new TimeoutException($"后端启动超过 60 秒，请检查 PostgreSQL 和日志：{LogPath}");
                    await Task.Delay(500, cancel);
                }
            }

            if (browser.CoreWebView2 is null)
            {
                Directory.CreateDirectory(LocalData);
                var profile = Path.Combine(LocalData, "WebView2");
                var environment = await CoreWebView2Environment.CreateAsync(userDataFolder: profile);
                cancel.ThrowIfCancellationRequested();
                await browser.EnsureCoreWebView2Async(environment);

                // WebView2 forwards accelerator keys through the WinForms KeyDown event.
                browser.KeyDown += (_, args) =>
                {
                    if (args.KeyCode != Keys.F11 && (!fullScreen || args.KeyCode != Keys.Escape)) return;
                    args.Handled = true;
                    args.SuppressKeyPress = true;
                    if (!IsDisposed && !Disposing) BeginInvoke((Action)ToggleFullScreen);
                };

                var core = browser.CoreWebView2;
                core.WebMessageReceived += async (_, args) =>
                {
                    if (ownedBackend is not { HasExited: false } || ownedSetupToken is null ||
                        !SetupTokenBridge.IsSetupDocument(args.Source) || !SetupTokenBridge.IsSetupDocument(core.Source)) return;
                    try
                    {
                        if (args.TryGetWebMessageAsString() != "artex:request-setup-token") return;
                        using var response = await http.GetAsync(new Uri(Service, "api/auth/status"), lifetime.Token);
                        response.EnsureSuccessStatusCode();
                        using var state = JsonDocument.Parse(await response.Content.ReadAsStringAsync(lifetime.Token));
                        if (!state.RootElement.TryGetProperty("initialized", out var initialized) ||
                            initialized.ValueKind is not (JsonValueKind.True or JsonValueKind.False)) return;
                        if (initialized.GetBoolean()) { ownedSetupToken = null; return; }
                        if (!SetupTokenBridge.IsSetupDocument(core.Source) || ownedBackend is not { HasExited: false } || ownedSetupToken is null) return;
                        // Re-check location inside the executed document too: navigation may race the status request.
                        var payload = JsonSerializer.Serialize(ownedSetupToken);
                        await core.ExecuteScriptAsync("if(location.origin==='http://127.0.0.1:8787' && /^\\/setup\\/?$/.test(location.pathname) && !location.search && !location.hash) window.dispatchEvent(new CustomEvent('artex:setup-token',{detail:" + payload + "}));");
                    }
                    catch (Exception ex) when (ex is HttpRequestException or JsonException or InvalidOperationException or ArgumentException or OperationCanceledException or System.Runtime.InteropServices.COMException) { }
                };
                core.NavigationStarting += (_, args) =>
                {
                    if (!args.Uri.StartsWith(Service.AbsoluteUri, StringComparison.OrdinalIgnoreCase)) return;
                    pageResourceFailed = false;
                    pageLoadWatchdog.Stop();
                    pageLoadWatchdog.Start();
                    ShowPageMessage("正在加载 ARTEX 界面…", "正在等待网页完成加载", false);
                };
                core.NavigationCompleted += (_, args) =>
                {
                    pageLoadWatchdog.Stop();
                    if (!args.IsSuccess)
                    {
                        ready = false;
                        status.Text = "界面加载失败 · 可重新连接";
                        var details = "WebView2 导航失败：" + args.WebErrorStatus;
                        LogDesktopFailure(details);
                        ShowPageMessage("界面加载失败", details + "\n请确认本机服务可访问，再点击“重新连接”。", true);
                        return;
                    }
                    if (pageResourceFailed) return;
                    ready = true;
                    pageState.Hide();
                    status.Text = ownedBackend is { HasExited: false }
                        ? "已连接 · 本机 ARTEX 后端由桌面客户端管理 · 127.0.0.1:8787"
                        : "已连接 · 正在使用已运行的 ARTEX 后端 · 127.0.0.1:8787";
                };
                core.WebResourceResponseReceived += (_, args) =>
                {
                    if (!Uri.TryCreate(args.Request.Uri, UriKind.Absolute, out var uri) ||
                        args.Response.StatusCode < 400) return;
                    var staticAsset = uri.AbsolutePath.StartsWith("/_next/", StringComparison.OrdinalIgnoreCase) &&
                        (uri.AbsolutePath.EndsWith(".js", StringComparison.OrdinalIgnoreCase) ||
                         uri.AbsolutePath.EndsWith(".css", StringComparison.OrdinalIgnoreCase));
                    var document = uri.Host == Service.Host && uri.Port == Service.Port &&
                        (uri.AbsolutePath == "/" ||
                         (uri.AbsolutePath.EndsWith('/') &&
                          !uri.AbsolutePath.StartsWith("/api/", StringComparison.OrdinalIgnoreCase)));
                    if (!staticAsset && !document) return;
                    pageResourceFailed = true;
                    ready = false;
                    pageLoadWatchdog.Stop();
                    var details = "网页请求失败：HTTP " + args.Response.StatusCode + " " + uri.AbsolutePath;
                    LogDesktopFailure(details);
                    status.Text = "ARTEX 网页资源请求失败 · 可重新连接";
                    ShowPageMessage("界面加载失败", details + "\n请点击“重新连接”重新加载。", true);
                };
                core.ProcessFailed += (_, args) =>
                {
                    var kind = args.ProcessFailedKind;
                    if (kind is not (CoreWebView2ProcessFailedKind.BrowserProcessExited or
                        CoreWebView2ProcessFailedKind.RenderProcessExited or
                        CoreWebView2ProcessFailedKind.RenderProcessUnresponsive)) return;
                    if (!IsHandleCreated || IsDisposed || Disposing) return;
                    BeginInvoke((Action)(() =>
                    {
                        if (ReferenceEquals(browser.CoreWebView2, core)) HandleBrowserFailure(kind);
                    }));
                };
                core.NewWindowRequested += (_, args) =>
                {
                    if (Uri.TryCreate(args.Uri, UriKind.Absolute, out var uri) &&
                        (uri.Scheme == Uri.UriSchemeHttps || uri.Scheme == Uri.UriSchemeHttp))
                    {
                        Process.Start(new ProcessStartInfo(uri.AbsoluteUri) { UseShellExecute = true });
                        args.Handled = true;
                    }
                };
            }

            // Backend health alone is not proof that the UI has rendered.
            // NavigationCompleted is responsible for marking the page ready.
            ready = false;
            status.Text = "正在加载 ARTEX 页面…";
            var destination = Service;
            if (ownedBackend is { HasExited: false })
            {
                using var response = await http.GetAsync(new Uri(Service, "api/auth/status"), cancel);
                response.EnsureSuccessStatusCode();
                using var initialization = JsonDocument.Parse(await response.Content.ReadAsStringAsync(cancel));
                if (!initialization.RootElement.TryGetProperty("initialized", out var initialized) ||
                    initialized.ValueKind is not (JsonValueKind.True or JsonValueKind.False))
                    throw new InvalidDataException("无法确认登录账号初始化状态，请重新连接。");
                destination = new Uri(Service, SetupTokenBridge.InitialRoute(initialized.GetBoolean()));
            }
            browser.CoreWebView2.Navigate(destination.AbsoluteUri);
        }
        catch (OperationCanceledException) when (cancel.IsCancellationRequested)
        {
            // The window is closing.
        }
        catch (Exception ex)
        {
            if (cancel.IsCancellationRequested) return;
            ready = false;
            status.Text = "连接失败 · 可点击“数据库配置”修正连接或“重新连接”重试";
            LogDesktopFailure("连接失败：" + ex);
            var backendFailure = ex.Message.Contains("后端", StringComparison.OrdinalIgnoreCase);
            var databaseFailure = backendFailure || ex.Message.Contains("PostgreSQL", StringComparison.OrdinalIgnoreCase);
            var message = backendFailure ? ExplainBackendFailure() :
                ex.Message.Contains("8787", StringComparison.OrdinalIgnoreCase)
                    ? "本机服务端口被占用，ARTEX 无法连接到正确的服务。"
                    : "ARTEX 未能完成启动或页面加载，请检查下面的建议。";
            var overridden = !string.IsNullOrWhiteSpace(Environment.GetEnvironmentVariable("ARTEX_PG_DSN"));
            ShowPageMessage("ARTEX 暂时无法启动", message, true,
                databaseProblem: databaseFailure,
                technical: ex.Message + "\n后端日志：" + LogPath,
                nextStep: databaseFailure && overridden
                    ? "已检测到 ARTEX_PG_DSN 环境变量，它会覆盖界面保存的配置。请先修改或移除该变量、重新启动 ARTEX，再检查数据库连接。"
                    : databaseFailure
                        ? "建议先点击“修复数据库连接”，检查 PostgreSQL 服务、账号及数据库。"
                        : "建议先点击“重新连接”；如仍失败，查看本机日志定位原因。");
        }
        finally
        {
            Interlocked.Exchange(ref connecting, 0);
        }
    }

    private void ShowPageMessage(string title, string details, bool canRetry,
        bool databaseProblem = false, string? technical = null, string? nextStep = null)
    {
        pageState.Display(title, details, canRetry, databaseProblem, technical, nextStep);
    }

    private string ExplainBackendFailure()
    {
        try
        {
            if (File.Exists(LogPath))
            {
                // Old failures stay in the append-only log, so only inspect this
                // startup's output rather than misdiagnosing from a previous run.
                using var stream = new FileStream(LogPath, FileMode.Open, FileAccess.Read,
                    FileShare.ReadWrite | FileShare.Delete);
                stream.Seek(Math.Max(backendLogStartOffset, stream.Length - 131072), SeekOrigin.Begin);
                using var reader = new StreamReader(stream);
                var recent = reader.ReadToEnd();
                if (recent.Contains("password authentication failed", StringComparison.OrdinalIgnoreCase))
                    return "PostgreSQL 拒绝了当前账号密码，请在数据库设置中更正后重试。";
                if (recent.Contains("role", StringComparison.OrdinalIgnoreCase) &&
                    recent.Contains("does not exist", StringComparison.OrdinalIgnoreCase))
                    return "配置的 PostgreSQL 用户不存在，可以通过引导初始化专用账户。";
                if (recent.Contains("database", StringComparison.OrdinalIgnoreCase) &&
                    recent.Contains("does not exist", StringComparison.OrdinalIgnoreCase))
                    return "目标数据库尚未建立，请在数据库向导中完成初始化。";
                if (recent.Contains("connection refused", StringComparison.OrdinalIgnoreCase) ||
                    recent.Contains("connectex", StringComparison.OrdinalIgnoreCase))
                    return "PostgreSQL 连接被拒绝。请确认本机数据库服务已启动且端口正确。";
            }
        }
        catch (IOException) { }
        catch (UnauthorizedAccessException) { }
        return "ARTEX 后端未能启动。数据库连接或启动依赖可能存在问题。";
    }

    private static void OpenBackendLog()
    {
        var target = File.Exists(LogPath) ? LogPath :
            Path.Combine(LocalData, "desktop-ui.log");
        if (File.Exists(target))
            Process.Start(new ProcessStartInfo("notepad.exe", "\"" + target + "\"") { UseShellExecute = true });
        else
        {
            Directory.CreateDirectory(LocalData);
            Process.Start(new ProcessStartInfo("explorer.exe", "\"" + LocalData + "\"") { UseShellExecute = true });
        }
    }

    private void RetryConnection()
    {
        if (lifetime.IsCancellationRequested || Volatile.Read(ref connecting) != 0) return;
        if (browserRecreationNeeded)
        {
            var previous = browser;
            browser = new WebView2 { Dock = DockStyle.Fill };
            viewHost.Controls.Add(browser);
            browser.SendToBack();
            pageState.BringToFront();
            previous.Dispose();
            browserRecreationNeeded = false;
        }
        rendererRecoveryAttempted = false;
        _ = ConnectAsync(lifetime.Token);
    }

    private void HandleBrowserFailure(CoreWebView2ProcessFailedKind kind)
    {
        if (lifetime.IsCancellationRequested || IsDisposed) return;
        ready = false;
        pageLoadWatchdog.Stop();
        LogDesktopFailure("WebView2 进程异常：" + kind);
        if (kind is CoreWebView2ProcessFailedKind.BrowserProcessExited or
            CoreWebView2ProcessFailedKind.RenderProcessUnresponsive)
        {
            browserRecreationNeeded = true;
            status.Text = "WebView2 异常 · 请重新连接";
            ShowPageMessage("网页进程已停止", "WebView2：" + kind + "\n点击“重新连接”重建网页视图。", true);
            return;
        }
        // A renderer crash can be recovered by reloading; try once per user retry.
        if (!rendererRecoveryAttempted)
        {
            rendererRecoveryAttempted = true;
            status.Text = "网页渲染进程退出 · 正在恢复…";
            ShowPageMessage("正在恢复界面…", "网页渲染进程意外退出，正在重新载入。", false);
            browser.CoreWebView2?.Reload();
        }
        else
        {
            status.Text = "网页渲染进程再次退出 · 请重新连接";
            ShowPageMessage("网页渲染进程异常", "自动恢复未成功，请点击“重新连接”。", true);
        }
    }

    private static void LogDesktopFailure(string message)
    {
        try
        {
            Directory.CreateDirectory(LocalData);
            File.AppendAllText(Path.Combine(LocalData, "desktop-ui.log"),
                "[" + DateTimeOffset.Now.ToString("O") + "] " + message + Environment.NewLine);
        }
        catch (IOException) { }
        catch (UnauthorizedAccessException) { }
    }

    private void EditDatabaseConfiguration()
    {
        if (lifetime.IsCancellationRequested || Volatile.Read(ref connecting) != 0) return;
        var owned = ownedBackend is { HasExited: false };
        if (owned && MessageBox.Show(this,
            "保存新的数据库连接后，将重新启动此桌面窗口管理的 ARTEX 后端。" +
            "\n正在执行的任务可能中断。是否继续？",
            "修改数据库配置", MessageBoxButtons.YesNo, MessageBoxIcon.Question) != DialogResult.Yes)
            return;
        try
        {
            using var setup = new DatabaseSetupForm(editExisting: true);
            if (setup.ShowDialog(this) != DialogResult.OK) return;
            if (owned)
            {
                ownedBackend!.Kill(entireProcessTree: true);
                ownedBackend.WaitForExit(5000);
                ownedBackend.Dispose();
                ownedBackend = null;
            }
            if (!owned && ready && ownedBackend is null)
            {
                MessageBox.Show(this, "数据库配置已保存。当前连接的是外部启动的 ARTEX 服务，" +
                    "桌面客户端不会关闭它；请自行重启外部后端使新配置生效。",
                    "数据库配置", MessageBoxButtons.OK, MessageBoxIcon.Information);
                return;
            }
            _ = ConnectAsync(lifetime.Token);
        }
        catch (Exception ex)
        {
            MessageBox.Show(this, "无法编辑数据库配置：" + ex.Message,
                "数据库配置", MessageBoxButtons.OK, MessageBoxIcon.Error);
        }
    }

    private static async Task<bool> IsPortOccupiedAsync(CancellationToken cancel)
    {
        using var client = new TcpClient();
        try
        {
            await client.ConnectAsync(System.Net.IPAddress.Loopback, 8787, cancel);
            return true;
        }
        catch (SocketException) { return false; }
    }

    private void StartBackend()
    {
        var executable = Path.Combine(AppContext.BaseDirectory, "artex.exe");
        if (!File.Exists(executable))
            throw new FileNotFoundException(
                "桌面版缺少 artex.exe。请解压完整的 ARTEX Desktop 压缩包，不要只复制桌面程序。", executable);

        Directory.CreateDirectory(LocalData);
        lock (logLock)
        {
            backendLog?.Dispose();
            backendLogStartOffset = File.Exists(LogPath) ? new FileInfo(LogPath).Length : 0;
            backendLog = new StreamWriter(LogPath, append: true) { AutoFlush = true };
        }

        var startInfo = new ProcessStartInfo(executable)
        {
            WorkingDirectory = AppContext.BaseDirectory,
            UseShellExecute = false,
            CreateNoWindow = true,
            RedirectStandardOutput = true,
            RedirectStandardError = true,
        };
        startInfo.ArgumentList.Add("-addr");
        startInfo.ArgumentList.Add("127.0.0.1:8787");
        ownedSetupToken = Convert.ToHexString(System.Security.Cryptography.RandomNumberGenerator.GetBytes(24));
        startInfo.Environment["ARTEX_SETUP_TOKEN"] = ownedSetupToken;
        if (string.IsNullOrWhiteSpace(Environment.GetEnvironmentVariable("ARTEX_KB_DIR")) &&
            File.Exists(Path.Combine(AppContext.BaseDirectory, "knowledge", "knowledge.db")))
        {
            using var configuration = File.Exists(DatabaseSetupForm.ConfigPath) ?
                JsonDocument.Parse(File.ReadAllText(DatabaseSetupForm.ConfigPath)) : null;
            if (configuration is null || !configuration.RootElement.TryGetProperty("knowledge_dir", out var knowledgeDir) ||
                string.IsNullOrWhiteSpace(knowledgeDir.ToString()))
                startInfo.Environment["ARTEX_KB_DIR"] = Path.Combine(AppContext.BaseDirectory, "knowledge");
        }

        var process = new Process { StartInfo = startInfo, EnableRaisingEvents = true };
        process.OutputDataReceived += WriteBackendLog;
        process.ErrorDataReceived += WriteBackendLog;
        if (!process.Start())
        {
            process.Dispose();
            throw new InvalidOperationException("无法启动 artex.exe");
        }
        process.BeginOutputReadLine();
        process.BeginErrorReadLine();
        ownedBackend = process;
    }

    private void WriteBackendLog(object sender, DataReceivedEventArgs args)
    {
        if (args.Data is null) return;
        // An automatically supplied initialization token never belongs in a persistent log.
        var line = ownedSetupToken is { Length: > 0 } ? args.Data.Replace(ownedSetupToken, "***", StringComparison.Ordinal) : args.Data;
        lock (logLock)
        {
            backendLog?.WriteLine(line);
        }
    }

    private async Task<bool> IsHealthyAsync(CancellationToken cancel)
    {
        try
        {
            using var response = await http.GetAsync(new Uri(Service, "api/health"), cancel);
            if (!response.IsSuccessStatusCode) return false;
            using var document = await JsonDocument.ParseAsync(
                await response.Content.ReadAsStreamAsync(cancel), cancellationToken: cancel);
            return document.RootElement.TryGetProperty("service", out var service) &&
                   service.GetString() == "artex" &&
                   document.RootElement.TryGetProperty("ok", out var ok) &&
                   ok.ValueKind == JsonValueKind.True;
        }
        catch (OperationCanceledException) when (cancel.IsCancellationRequested)
        {
            throw;
        }
        catch (HttpRequestException) { return false; }
        catch (TaskCanceledException) { return false; }
        catch (JsonException) { return false; }
    }

    private static void OpenInBrowser() =>
        Process.Start(new ProcessStartInfo(Service.AbsoluteUri) { UseShellExecute = true });

    private void ShowDesktop()
    {
        var wasMinimized = WindowState == FormWindowState.Minimized;
        var lastState = restoreWindowState;
        Show();
        if (wasMinimized)
            WindowState = fullScreen ? FormWindowState.Normal : lastState;
        if (fullScreen)
            Bounds = Screen.FromControl(this).Bounds;
        Activate();
    }

    private void ToggleFullScreen()
    {
        if (!fullScreen)
        {
            var screen = Screen.FromControl(this);
            beforeFullScreenState = WindowState;
            beforeFullScreenBounds = WindowState == FormWindowState.Maximized ? RestoreBounds : Bounds;
            fullScreen = true;
            WindowState = FormWindowState.Normal;
            FormBorderStyle = FormBorderStyle.None;
            Bounds = screen.Bounds;
            toolbar.Visible = false;
            statusBar.Visible = false;
            fullScreenButton.Text = "退出全屏 (F11)";
        }
        else
        {
            fullScreen = false;
            FormBorderStyle = FormBorderStyle.Sizable;
            WindowState = FormWindowState.Normal;
            Bounds = beforeFullScreenBounds;
            WindowState = beforeFullScreenState;
            toolbar.Visible = true;
            statusBar.Visible = true;
            fullScreenButton.Text = "全屏 (F11)";
        }
    }

    protected override bool ProcessCmdKey(ref Message msg, Keys keyData)
    {
        if (keyData == Keys.F11 || (fullScreen && keyData == Keys.Escape))
        {
            ToggleFullScreen();
            return true;
        }
        return base.ProcessCmdKey(ref msg, keyData);
    }

    private void ExitDesktop()
    {
        exiting = true;
        Close();
    }

    protected override void OnResize(EventArgs e)
    {
        base.OnResize(e);
        if (WindowState == FormWindowState.Minimized)
        {
            if (ready) Hide();
        }
        else if (!fullScreen)
        {
            restoreWindowState = WindowState;
        }
    }

    protected override void OnFormClosing(FormClosingEventArgs e)
    {
        if (!exiting && e.CloseReason == CloseReason.UserClosing)
        {
            e.Cancel = true;
            Hide();
            return;
        }
        base.OnFormClosing(e);
    }

    protected override void OnFormClosed(FormClosedEventArgs e)
    {
        lifetime.Cancel();
        pageLoadWatchdog.Stop();
        pageLoadWatchdog.Dispose();
        tray.Visible = false;
        tray.Dispose();

        // Only stop the process started by this window. A pre-existing ARTEX
        // service remains running when a desktop window closes.
        if (ownedBackend is not null)
        {
            try
            {
                if (!ownedBackend.HasExited)
                {
                    ownedBackend.Kill(entireProcessTree: true);
                    ownedBackend.WaitForExit(5000);
                }
            }
            catch (InvalidOperationException)
            {
                // The child can exit between HasExited and Kill.
            }
            ownedBackend.Dispose();
        }
        ManagedPostgres.StopOwned();
        lock (logLock)
        {
            backendLog?.Dispose();
            backendLog = null;
        }
        http.Dispose();
        lifetime.Dispose();
        browser.Dispose();
        base.OnFormClosed(e);
    }
}
