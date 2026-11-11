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
    private readonly WebView2 browser = new() { Dock = DockStyle.Fill };
    private readonly ToolStrip toolbar = new() { Dock = DockStyle.Top, GripStyle = ToolStripGripStyle.Hidden };
    private readonly StatusStrip statusBar = new() { Dock = DockStyle.Bottom };
    private readonly ToolStripButton fullScreenButton = new("全屏 (F11)");
    private readonly ToolStripStatusLabel status = new("正在启动 ARTEX…");
    private readonly HttpClient http = new() { Timeout = TimeSpan.FromSeconds(2) };
    private readonly CancellationTokenSource lifetime = new();
    private readonly NotifyIcon tray;
    private readonly object logLock = new();
    private Process? ownedBackend;
    private StreamWriter? backendLog;
    private bool exiting;
    private bool ready;
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
        toolbar.Items.Add(Button("后退", () =>
        {
            if (browser.CoreWebView2?.CanGoBack == true) browser.CoreWebView2.GoBack();
        }));
        toolbar.Items.Add(Button("前进", () =>
        {
            if (browser.CoreWebView2?.CanGoForward == true) browser.CoreWebView2.GoForward();
        }));
        toolbar.Items.Add(Button("刷新", () => browser.CoreWebView2?.Reload()));
        fullScreenButton.Click += (_, _) => ToggleFullScreen();
        toolbar.Items.Add(fullScreenButton);
        toolbar.Items.Add(new ToolStripSeparator());
        toolbar.Items.Add(Button("在浏览器打开", OpenInBrowser));
        toolbar.Items.Add(Button("数据库配置", EditDatabaseConfiguration));
        toolbar.Items.Add(Button("打开配置目录", () =>
            Process.Start(new ProcessStartInfo("explorer.exe", "\"" + AppContext.BaseDirectory.TrimEnd('\\') + "\"")
            { UseShellExecute = true })));
        toolbar.Items.Add(new ToolStripSeparator());
        toolbar.Items.Add(Button("重新连接", () => _ = ConnectAsync(lifetime.Token)));
        toolbar.Items.Add(Button("退出", ExitDesktop));

        statusBar.Items.Add(status);

        Controls.Add(browser);
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
        var button = new ToolStripButton(label);
        button.Click += (_, _) => click();
        return button;
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

                    if (DatabaseSetupForm.NeedsConfiguration)
                    {
                        using var setup = new DatabaseSetupForm();
                        if (setup.ShowDialog(this) != DialogResult.OK)
                        {
                            status.Text = "未配置数据库 · 点击“重新连接”继续设置";
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

                browser.AcceleratorKeyPressed += (_, args) =>
                {
                    if (args.KeyEventKind is not (CoreWebView2KeyEventKind.KeyDown or
                        CoreWebView2KeyEventKind.SystemKeyDown)) return;
                    if (args.VirtualKey != 0x7A && (!fullScreen || args.VirtualKey != 0x1B)) return;
                    args.Handled = true;
                    // Bit 30 marks a repeated keydown; do not toggle twice on a long press.
                    if ((args.KeyEventLParam & (1 << 30)) == 0 && !IsDisposed && !Disposing)
                        BeginInvoke((Action)ToggleFullScreen);
                };

                browser.CoreWebView2.NewWindowRequested += (_, args) =>
                {
                    if (Uri.TryCreate(args.Uri, UriKind.Absolute, out var uri) &&
                        (uri.Scheme == Uri.UriSchemeHttps || uri.Scheme == Uri.UriSchemeHttp))
                    {
                        Process.Start(new ProcessStartInfo(uri.AbsoluteUri) { UseShellExecute = true });
                        args.Handled = true;
                    }
                };
            }

            browser.Source = Service;
            ready = true;
            status.Text = ownedBackend is { HasExited: false }
                ? "已连接 · 本机 ARTEX 后端由桌面客户端管理 · 127.0.0.1:8787"
                : "已连接 · 正在使用已运行的 ARTEX 后端 · 127.0.0.1:8787";
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
            var message = ex.Message + "\n\n" +
                "可以点击工具栏“数据库配置”修改连接。请确认 PostgreSQL（config.json 或 ARTEX_PG_DSN），" +
                "并已安装 Microsoft Edge WebView2 Runtime。\n" +
                "如需排查，打开配置目录或查看：" + LogPath;
            MessageBox.Show(this, message, "ARTEX 启动失败", MessageBoxButtons.OK, MessageBoxIcon.Error);
        }
        finally
        {
            Interlocked.Exchange(ref connecting, 0);
        }
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
        lock (logLock)
        {
            backendLog?.WriteLine(args.Data);
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
        Show();
        if (WindowState == FormWindowState.Minimized)
            WindowState = fullScreen ? FormWindowState.Normal : restoreWindowState;
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
            beforeFullScreenBounds = Bounds;
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
        tray.Visible = false;
        tray.Dispose();

        // Only stop the process started by this window. A pre-existing ARTEX
        // service remains running when a desktop window closes.
        if (ownedBackend is not null)
        {
            try
            {
                if (!ownedBackend.HasExited)
                    ownedBackend.Kill(entireProcessTree: true);
            }
            catch (InvalidOperationException)
            {
                // The child can exit between HasExited and Kill.
            }
            ownedBackend.Dispose();
        }
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
