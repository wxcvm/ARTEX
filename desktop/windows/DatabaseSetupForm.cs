using System.Diagnostics;
using System.Drawing;
using System.Text.Json;
using System.Windows.Forms;

namespace Artex.Desktop;

// First-launch helper for ARTEX's existing config.json format. PostgreSQL
// remains an external service; the dialog neither installs nor modifies it.
internal sealed class DatabaseSetupForm : Form
{
    internal static string ConfigPath
    {
        get
        {
            var custom = Environment.GetEnvironmentVariable("ARTEX_CONFIG");
            return !string.IsNullOrWhiteSpace(custom)
                ? Path.GetFullPath(custom.Trim())
                : Path.Combine(AppContext.BaseDirectory, "config.json");
        }
    }

    internal static bool NeedsConfiguration =>
        string.IsNullOrWhiteSpace(Environment.GetEnvironmentVariable("ARTEX_PG_DSN")) &&
        !File.Exists(ConfigPath);

    private readonly TextBox host = new() { Text = "127.0.0.1", Width = 300 };
    private readonly TextBox port = new() { Text = "5432", Width = 300 };
    private readonly TextBox user = new() { Text = "artex", Width = 300 };
    private readonly TextBox password = new() { UseSystemPasswordChar = true, Width = 300 };
    private readonly TextBox database = new() { Text = "artex", Width = 300 };
    private readonly TextBox sslMode = new() { Text = "disable", Width = 300 };
    private readonly Label discoveryStatus = new()
    {
        Text = "正在检查本机 PostgreSQL…",
        AutoSize = false,
        Dock = DockStyle.Fill,
        ForeColor = SystemColors.GrayText,
    };
    private readonly ComboBox localServices = new()
    {
        DropDownStyle = ComboBoxStyle.DropDownList,
        Width = 175,
        Enabled = false,
    };
    private readonly Button detectLocal = new() { Text = "检测本机", Width = 94, Height = 30 };
    private readonly Button startLocal = new() { Text = "启动所选服务", Width = 112, Height = 30, Enabled = false };
    private readonly CancellationTokenSource discoveryLifetime = new();

    internal DatabaseSetupForm()
    {
        Text = "ARTEX · 首次数据库配置";
        Icon = System.Drawing.Icon.ExtractAssociatedIcon(Application.ExecutablePath) ?? SystemIcons.Application;
        StartPosition = FormStartPosition.CenterParent;
        ClientSize = new Size(610, 485);
        MinimumSize = new Size(570, 445);
        FormBorderStyle = FormBorderStyle.FixedDialog;
        MaximizeBox = false;
        MinimizeBox = false;

        var layout = new TableLayoutPanel
        {
            Dock = DockStyle.Fill,
            ColumnCount = 2,
            RowCount = 11,
            Padding = new Padding(16),
            AutoSize = false,
        };
        layout.ColumnStyles.Add(new ColumnStyle(SizeType.Absolute, 145));
        layout.ColumnStyles.Add(new ColumnStyle(SizeType.Percent, 100));

        var intro = new Label
        {
            Text = "首次运行需要连接 PostgreSQL。请填写你已有的数据库连接信息。" +
                   "\n如果使用手册中的 Docker 示例，端口改为 5433。",
            AutoSize = true,
            Dock = DockStyle.Fill,
        };
        layout.Controls.Add(intro, 0, 0);
        layout.SetColumnSpan(intro, 2);

        AddField(layout, 1, "服务器地址", host);
        AddField(layout, 2, "端口", port);
        AddField(layout, 3, "用户名", user);
        AddField(layout, 4, "密码", password);
        AddField(layout, 5, "数据库名", database);
        AddField(layout, 6, "SSL 模式", sslMode);

        var diagnostics = new FlowLayoutPanel
        {
            Dock = DockStyle.Fill,
            FlowDirection = FlowDirection.LeftToRight,
            WrapContents = false,
        };
        detectLocal.Click += async (_, _) => await DetectLocalAsync();
        startLocal.Click += async (_, _) => await StartLocalAsync();
        diagnostics.Controls.Add(detectLocal);
        diagnostics.Controls.Add(localServices);
        diagnostics.Controls.Add(startLocal);
        layout.Controls.Add(diagnostics, 0, 7);
        layout.SetColumnSpan(diagnostics, 2);
        layout.Controls.Add(discoveryStatus, 0, 8);
        layout.SetColumnSpan(discoveryStatus, 2);

        var configLabel = new Label
        {
            Text = "配置保存位置：" + ConfigPath,
            AutoSize = false,
            AutoEllipsis = true,
            Dock = DockStyle.Fill,
            ForeColor = SystemColors.GrayText,
        };
        layout.Controls.Add(configLabel, 0, 9);
        layout.SetColumnSpan(configLabel, 2);

        var buttons = new FlowLayoutPanel
        {
            Dock = DockStyle.Fill,
            FlowDirection = FlowDirection.RightToLeft,
            WrapContents = false,
        };
        var save = new Button { Text = "保存并连接", Width = 110, Height = 30 };
        save.Click += (_, _) => Save();
        var test = new Button { Text = "测试数据库连接", Width = 124, Height = 30 };
        test.Click += async (_, _) => await TestConnectionAsync(test);
        var cancel = new Button { Text = "稍后配置", Width = 100, Height = 30, DialogResult = DialogResult.Cancel };
        buttons.Controls.Add(save);
        buttons.Controls.Add(test);
        buttons.Controls.Add(cancel);
        layout.Controls.Add(buttons, 0, 10);
        layout.SetColumnSpan(buttons, 2);

        AcceptButton = save;
        CancelButton = cancel;
        Controls.Add(layout);
        Shown += async (_, _) => await DetectLocalAsync();
        FormClosed += (_, _) =>
        {
            discoveryLifetime.Cancel();
            discoveryLifetime.Dispose();
        };
    }

    private async Task DetectLocalAsync()
    {
        if (discoveryLifetime.IsCancellationRequested) return;
        detectLocal.Enabled = false;
        startLocal.Enabled = false;
        discoveryStatus.Text = "正在探测本机 PostgreSQL（5432 / 5433）…";
        try
        {
            var found = await PostgresDiscovery.FindLocalPortAsync(discoveryLifetime.Token);
            if (IsDisposed || discoveryLifetime.IsCancellationRequested) return;
            localServices.Items.Clear();
            foreach (var service in PostgresDiscovery.InstalledWindowsServices())
                localServices.Items.Add(service);
            localServices.Enabled = localServices.Items.Count != 0;
            if (localServices.Enabled) localServices.SelectedIndex = 0;
            startLocal.Enabled = localServices.Enabled;

            if (found is int detectedPort)
            {
                if (host.Text.Trim() is "127.0.0.1" or "localhost" &&
                    port.Text.Trim() is "5432" or "5433")
                    port.Text = detectedPort.ToString();
                discoveryStatus.Text = $"发现本机 PostgreSQL，端口 {detectedPort}。" +
                    " 仅确认服务可连接，账号和密码仍需填写。";
            }
            else
            {
                discoveryStatus.Text = localServices.Items.Count > 0
                    ? "已发现已安装的 PostgreSQL 服务，但常用端口尚未响应。可选择服务启动，或填写自定义端口。"
                    : "未在 5432 / 5433 发现 PostgreSQL。可填写远程数据库，或先安装并启动本地服务。";
            }
        }
        catch (OperationCanceledException) when (discoveryLifetime.IsCancellationRequested) { }
        catch (Exception ex)
        {
            if (!IsDisposed)
                discoveryStatus.Text = "检测未完成：" + ex.Message + "。仍可手动填写连接信息。";
        }
        finally
        {
            if (!IsDisposed) detectLocal.Enabled = true;
        }
    }

    private async Task StartLocalAsync()
    {
        if (localServices.SelectedItem is not string service) return;
        if (MessageBox.Show(this, $"是否启动本机已安装的 PostgreSQL 服务“{service}”？" +
            "\\n不会创建或修改数据库和账户。",
            "启动本机数据库服务", MessageBoxButtons.YesNo, MessageBoxIcon.Question) != DialogResult.Yes)
            return;

        startLocal.Enabled = false;
        discoveryStatus.Text = $"正在请求启动 {service}…";
        try
        {
            var started = await PostgresDiscovery.StartWindowsServiceAsync(service, discoveryLifetime.Token);
            if (IsDisposed || discoveryLifetime.IsCancellationRequested) return;
            if (!started)
            {
                discoveryStatus.Text = "服务启动命令未成功。可在 Windows 服务管理器中确认服务状态或管理员权限。";
                return;
            }
            await Task.Delay(1000, discoveryLifetime.Token);
            await DetectLocalAsync();
        }
        catch (OperationCanceledException) when (discoveryLifetime.IsCancellationRequested) { }
        catch (Exception ex)
        {
            if (!IsDisposed) discoveryStatus.Text = "无法启动服务：" + ex.Message;
        }
        finally
        {
            if (!IsDisposed) startLocal.Enabled = localServices.SelectedItem is not null;
        }
    }

    private async Task TestConnectionAsync(Button button)
    {
        var hostname = host.Text.Trim();
        var username = user.Text.Trim();
        var dbName = database.Text.Trim();
        if (hostname.Length == 0 || username.Length == 0 || dbName.Length == 0 ||
            !int.TryParse(port.Text, out var portNumber) || portNumber is < 1 or > 65535)
        {
            MessageBox.Show(this, "请先填写服务器、用户名、数据库名和有效的端口。",
                "连接测试", MessageBoxButtons.OK, MessageBoxIcon.Warning);
            return;
        }
        var executable = Path.Combine(AppContext.BaseDirectory, "artex.exe");
        if (!File.Exists(executable))
        {
            MessageBox.Show(this, "找不到后端 artex.exe，请使用完整 Windows 桌面包。",
                "连接测试", MessageBoxButtons.OK, MessageBoxIcon.Error);
            return;
        }

        // Send the credentials to the bundled Go PostgreSQL driver via a child
        // process environment variable. Neither the command line nor a temp
        // config file contains the password.
        var hostForUri = hostname.Contains(':') && !hostname.StartsWith('[')
            ? "[" + hostname + "]" : hostname;
        var dsn = "postgres://" + Uri.EscapeDataString(username) + ":" +
            Uri.EscapeDataString(password.Text) + "@" + hostForUri + ":" + portNumber +
            "/" + Uri.EscapeDataString(dbName) + "?sslmode=" +
            Uri.EscapeDataString(sslMode.Text.Trim() is "" ? "disable" : sslMode.Text.Trim());

        button.Enabled = false;
        discoveryStatus.Text = "正在使用 ARTEX 后端验证 PostgreSQL 账号与数据库…";
        try
        {
            var start = new ProcessStartInfo(executable)
            {
                WorkingDirectory = AppContext.BaseDirectory,
                UseShellExecute = false,
                CreateNoWindow = true,
                RedirectStandardError = true,
                RedirectStandardOutput = true,
            };
            start.ArgumentList.Add("-check-db");
            start.Environment["ARTEX_PG_DSN"] = dsn;
            using var process = Process.Start(start) ??
                throw new InvalidOperationException("无法启动连接检查程序。");
            using var timeout = CancellationTokenSource.CreateLinkedTokenSource(discoveryLifetime.Token);
            timeout.CancelAfter(TimeSpan.FromSeconds(12));
            await process.WaitForExitAsync(timeout.Token);
            var errorText = await process.StandardError.ReadToEndAsync(timeout.Token);
            if (IsDisposed || discoveryLifetime.IsCancellationRequested) return;
            if (process.ExitCode == 0)
            {
                discoveryStatus.Text = "PostgreSQL 连接验证成功。现在可以保存并连接。";
                return;
            }
            if (password.Text.Length != 0)
                errorText = errorText.Replace(password.Text, "***", StringComparison.Ordinal)
                    .Replace(Uri.EscapeDataString(password.Text), "***", StringComparison.Ordinal);
            discoveryStatus.Text = "PostgreSQL 连接失败，请检查配置。";
            MessageBox.Show(this, errorText.Length > 600 ? errorText[..600] : errorText,
                "数据库连接检查失败", MessageBoxButtons.OK, MessageBoxIcon.Warning);
        }
        catch (OperationCanceledException) when (discoveryLifetime.IsCancellationRequested) { }
        catch (OperationCanceledException)
        {
            if (!IsDisposed) discoveryStatus.Text = "连接测试超时；检查数据库服务或网络连接。";
        }
        catch (Exception ex)
        {
            if (!IsDisposed) discoveryStatus.Text = "连接检查未完成：" + ex.Message;
        }
        finally
        {
            if (!IsDisposed) button.Enabled = true;
        }
    }

    private static void AddField(TableLayoutPanel layout, int row, string name, TextBox field)
    {
        layout.Controls.Add(new Label { Text = name, Dock = DockStyle.Fill, TextAlign = ContentAlignment.MiddleLeft }, 0, row);
        field.Dock = DockStyle.Fill;
        layout.Controls.Add(field, 1, row);
    }

    private void Save()
    {
        var hostname = host.Text.Trim();
        var username = user.Text.Trim();
        var dbName = database.Text.Trim();
        if (hostname.Length == 0 || username.Length == 0 || dbName.Length == 0 ||
            !int.TryParse(port.Text, out var portNumber) || portNumber is < 1 or > 65535)
        {
            MessageBox.Show(this, "请填写服务器、用户名、数据库名和有效的端口（1–65535）。",
                "数据库配置", MessageBoxButtons.OK, MessageBoxIcon.Warning);
            return;
        }

        var configuration = new
        {
            database = new
            {
                host = hostname,
                port = portNumber,
                user = username,
                password = password.Text,
                dbname = dbName,
                sslmode = sslMode.Text.Trim(),
            },
            skill_dir = "",
        };

        try
        {
            // Never overwrite an existing config, including on a second launch
            // while the setup dialog is still open.
            var path = ConfigPath;
            Directory.CreateDirectory(Path.GetDirectoryName(Path.GetFullPath(path))!);
            using (var output = new FileStream(path, FileMode.CreateNew, FileAccess.Write, FileShare.None))
                JsonSerializer.Serialize(output, configuration, new JsonSerializerOptions { WriteIndented = true });
            DialogResult = DialogResult.OK;
            Close();
        }
        catch (IOException ex)
        {
            MessageBox.Show(this, "无法保存配置文件：" + ex.Message +
                "\n可以在配置目录手动创建 config.json，或使用 ARTEX_PG_DSN 环境变量。",
                "数据库配置", MessageBoxButtons.OK, MessageBoxIcon.Error);
        }
        catch (UnauthorizedAccessException ex)
        {
            MessageBox.Show(this, "当前目录没有写入权限：" + ex.Message +
                "\n可以选择可写安装目录，或使用 ARTEX_CONFIG 指定配置路径。",
                "数据库配置", MessageBoxButtons.OK, MessageBoxIcon.Error);
        }
    }
}
