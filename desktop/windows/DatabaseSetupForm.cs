using System.Diagnostics;
using System.Drawing;
using System.Security.Cryptography;
using System.Text.Json;
using System.Text.Json.Nodes;
using System.Windows.Forms;

namespace Artex.Desktop;

// First-launch and repair helper for ARTEX's config.json. Existing settings
// outside the edited database fields are preserved when saving.
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
    private readonly Button installLocal = new() { Text = "安装 PostgreSQL", Width = 140, Height = 30 };
    private readonly Button provisionApp = new() { Text = "初始化专用账户", Width = 135, Height = 30 };
    private readonly CancellationTokenSource discoveryLifetime = new();
    private readonly bool editExisting;

    internal DatabaseSetupForm(bool editExisting = false)
    {
        this.editExisting = editExisting;
        Text = editExisting ? "ARTEX · 修改数据库配置" : "ARTEX · 首次数据库配置";
        Icon = System.Drawing.Icon.ExtractAssociatedIcon(Application.ExecutablePath) ?? SystemIcons.Application;
        StartPosition = FormStartPosition.CenterParent;
        AutoScaleMode = AutoScaleMode.Dpi;
        ClientSize = new Size(860, 610);
        MinimumSize = new Size(620, 440);
        FormBorderStyle = FormBorderStyle.Sizable;
        MaximizeBox = true;
        MinimizeBox = false;

        // Keep the primary actions visible when Windows text size/DPI grows.
        // All editable settings and secondary actions live in an independently
        // scrollable section instead of forcing eleven rows into a fixed height.
        var root = new TableLayoutPanel
        {
            Dock = DockStyle.Fill,
            ColumnCount = 1,
            RowCount = 2,
            Margin = Padding.Empty,
            Padding = Padding.Empty,
        };
        root.ColumnStyles.Add(new ColumnStyle(SizeType.Percent, 100));
        root.RowStyles.Add(new RowStyle(SizeType.Percent, 100));
        root.RowStyles.Add(new RowStyle(SizeType.AutoSize));
        var scroll = new Panel { Dock = DockStyle.Fill, AutoScroll = true };
        var layout = new TableLayoutPanel
        {
            Dock = DockStyle.Top,
            ColumnCount = 2,
            RowCount = 11,
            Padding = new Padding(18, 16, 18, 12),
            Margin = Padding.Empty,
            AutoSize = true,
            AutoSizeMode = AutoSizeMode.GrowAndShrink,
        };
        layout.ColumnStyles.Add(new ColumnStyle(SizeType.Absolute, 145));
        layout.ColumnStyles.Add(new ColumnStyle(SizeType.Percent, 100));
        for (var row = 0; row < layout.RowCount; row++)
            layout.RowStyles.Add(new RowStyle(SizeType.AutoSize));

        var intro = new Label
        {
            Text = (editExisting ? "修改 PostgreSQL 连接信息。保存后需要重启后端。" :
                "首次运行需要连接 PostgreSQL。请填写连接信息。") +
                "\nDocker 手册示例端口是 5433；ARTEX_PG_DSN 环境变量优先于此配置。",
            AutoSize = true,
            Dock = DockStyle.Top,
            Margin = new Padding(3, 0, 3, 12),
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
            Dock = DockStyle.Top,
            AutoSize = true,
            AutoSizeMode = AutoSizeMode.GrowAndShrink,
            FlowDirection = FlowDirection.LeftToRight,
            WrapContents = true,
            Margin = new Padding(3, 12, 3, 8),
        };
        detectLocal.Click += async (_, _) => await DetectLocalAsync();
        startLocal.Click += async (_, _) => await StartLocalAsync();
        installLocal.Click += (_, _) => InstallPostgreSQL();
        diagnostics.Controls.Add(detectLocal);
        diagnostics.Controls.Add(localServices);
        diagnostics.Controls.Add(startLocal);
        diagnostics.Controls.Add(installLocal);
        layout.Controls.Add(diagnostics, 0, 7);
        layout.SetColumnSpan(diagnostics, 2);
        layout.Controls.Add(discoveryStatus, 0, 8);
        layout.SetColumnSpan(discoveryStatus, 2);
        discoveryStatus.AutoSize = true;
        discoveryStatus.Dock = DockStyle.Top;
        discoveryStatus.MinimumSize = new Size(0, 26);
        discoveryStatus.Margin = new Padding(3, 3, 3, 10);

        var configLabel = new Label
        {
            Text = "配置保存位置：" + ConfigPath,
            AutoSize = false,
            AutoEllipsis = true,
            Dock = DockStyle.Top,
            MinimumSize = new Size(0, 26),
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
        test.Click += async (_, _) => await RunDatabaseActionAsync(test, "-check-db");
        var create = new Button { Text = "创建目标数据库", Width = 124, Height = 30 };
        create.Click += async (_, _) => await RunDatabaseActionAsync(create, "-create-db");
        provisionApp.Click += async (_, _) => await ProvisionDedicatedAccountAsync();
        var cancel = new Button { Text = "稍后配置", Width = 100, Height = 30, DialogResult = DialogResult.Cancel };
        buttons.Controls.Add(save);
        buttons.Controls.Add(test);
        buttons.Controls.Add(create);
        buttons.Controls.Add(provisionApp);
        buttons.Controls.Add(cancel);
        layout.Controls.Add(buttons, 0, 10);
        layout.SetColumnSpan(buttons, 2);

        AcceptButton = save;
        CancelButton = cancel;
        Controls.Add(layout);
        if (editExisting && File.Exists(ConfigPath)) LoadExistingConfiguration();
        Shown += async (_, _) => await DetectLocalAsync();
        FormClosed += (_, _) =>
        {
            discoveryLifetime.Cancel();
            discoveryLifetime.Dispose();
        };
    }

    private void InstallPostgreSQL()
    {
        if (MessageBox.Show(this,
            "将启动 PostgreSQL 16 的官方 Windows 交互安装流程。\\n" +
            "安装程序需要联网下载、Windows 管理员授权，并由你设置数据库管理员密码。\\n" +
            "ARTEX 不会代填密码，也不会修改已有 PostgreSQL。是否继续？",
            "安装 PostgreSQL",
            MessageBoxButtons.YesNo, MessageBoxIcon.Question) != DialogResult.Yes)
            return;
        try
        {
            var launched = PostgresInstaller.LaunchInstallation();
            if ((host.Text.Trim() is "127.0.0.1" or "localhost") && user.Text.Trim() == "artex")
                user.Text = "postgres";
            discoveryStatus.Text = launched
                ? "安装窗口已打开。完成安装、记住设置的 postgres 密码，再点击“检测本机”。"
                : "系统未找到 WinGet，已打开 PostgreSQL 官方下载页。安装完成后点击“检测本机”。";
        }
        catch (Exception ex)
        {
            discoveryStatus.Text = "无法打开安装器：" + ex.Message;
        }
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
                if (!editExisting && (host.Text.Trim() is "127.0.0.1" or "localhost") &&
                    (port.Text.Trim() is "5432" or "5433"))
                    port.Text = detectedPort.ToString();
                discoveryStatus.Text = $"发现本机 PostgreSQL，端口 {detectedPort}。" +
                    " 仅确认服务可连接，账号和密码仍需填写。";
            }
            else
            {
                discoveryStatus.Text = localServices.Items.Count > 0
                    ? "已发现已安装的 PostgreSQL 服务，但常用端口尚未响应。可选择服务启动，或填写自定义端口。"
                    : "未检测到 PostgreSQL。可点“安装 PostgreSQL”打开官方安装器，或填写远程数据库。";
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

    private async Task ProvisionDedicatedAccountAsync()
    {
        var adminDsn = CurrentDsn();
        if (adminDsn is null) return;
        var dbName = database.Text.Trim();
        if (MessageBox.Show(this,
            $"将使用当前填写的管理员账号“{user.Text.Trim()}”在 PostgreSQL 上创建专用账户 artex，" +
            $"并创建归它所有的数据库“{dbName}”。\\n" +
            "专用账户密码将随机生成；已有数据库的所有权与已有账户密码不会改变。" +
            "\\n完成后会将向导账号切换为 artex。确认继续？",
            "初始化 ARTEX 专用账户",
            MessageBoxButtons.YesNo, MessageBoxIcon.Question) != DialogResult.Yes)
            return;

        var executable = Path.Combine(AppContext.BaseDirectory, "artex.exe");
        if (!File.Exists(executable))
        {
            MessageBox.Show(this, "找不到 artex.exe，请使用完整 Windows 桌面包。",
                "初始化专用账户", MessageBoxButtons.OK, MessageBoxIcon.Error);
            return;
        }
        var newPassword = Convert.ToHexString(RandomNumberGenerator.GetBytes(24));
        provisionApp.Enabled = false;
        discoveryStatus.Text = "正在创建 ARTEX 专用账户及其数据库…";
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
            start.ArgumentList.Add("-setup-app-db");
            // Admin credentials and generated app password are never command
            // line arguments or written to a temporary configuration file.
            start.Environment["ARTEX_PG_DSN"] = adminDsn;
            start.Environment["ARTEX_APP_DB_USER"] = "artex";
            start.Environment["ARTEX_APP_DB_PASSWORD"] = newPassword;
            using var process = Process.Start(start) ??
                throw new InvalidOperationException("无法启动专用账户初始化程序。");
            using var timeout = CancellationTokenSource.CreateLinkedTokenSource(discoveryLifetime.Token);
            timeout.CancelAfter(TimeSpan.FromSeconds(20));
            try
            {
                await process.WaitForExitAsync(timeout.Token);
            }
            catch (OperationCanceledException)
            {
                if (!process.HasExited) process.Kill(entireProcessTree: true);
                throw;
            }
            if (IsDisposed || discoveryLifetime.IsCancellationRequested) return;
            if (process.ExitCode == 0)
            {
                user.Text = "artex";
                password.Text = newPassword;
                discoveryStatus.Text = "ARTEX 专用账户与数据库已就绪，已切换到专用账户。请点击“保存并连接”。";
                MessageBox.Show(this,
                    "已创建专用账户 artex，数据库归它所有。\\n" +
                    "向导已填入随机生成的专用密码，请点击“保存并连接”以保存配置。\\n" +
                    "PostgreSQL 管理员密码不会保存到 ARTEX 配置中。",
                    "专用账户已就绪", MessageBoxButtons.OK, MessageBoxIcon.Information);
                return;
            }
            var details = await process.StandardError.ReadToEndAsync(timeout.Token);
            if (password.Text.Length > 0)
                details = details.Replace(password.Text, "***", StringComparison.Ordinal)
                    .Replace(Uri.EscapeDataString(password.Text), "***", StringComparison.Ordinal);
            discoveryStatus.Text = "专用账户初始化未完成，请检查管理员账号及数据库状态。";
            MessageBox.Show(this,
                details.Length > 650 ? details[..650] : details,
                "专用账户初始化失败", MessageBoxButtons.OK, MessageBoxIcon.Warning);
        }
        catch (OperationCanceledException) when (discoveryLifetime.IsCancellationRequested) { }
        catch (OperationCanceledException)
        {
            if (!IsDisposed) discoveryStatus.Text = "初始化超时，请检查 PostgreSQL 服务。";
        }
        catch (Exception ex)
        {
            if (!IsDisposed) discoveryStatus.Text = "专用账户初始化未完成：" + ex.Message;
        }
        finally
        {
            if (!IsDisposed) provisionApp.Enabled = true;
        }
    }

    private string? CurrentDsn()
    {
        var hostname = host.Text.Trim();
        var username = user.Text.Trim();
        var dbName = database.Text.Trim();
        if (hostname.Length == 0 || username.Length == 0 || dbName.Length == 0 ||
            !int.TryParse(port.Text, out var portNumber) || portNumber is < 1 or > 65535)
        {
            MessageBox.Show(this, "请填写服务器、用户名、数据库名和有效端口（1–65535）。",
                "数据库连接", MessageBoxButtons.OK, MessageBoxIcon.Warning);
            return null;
        }
        var hostForUri = hostname.Contains(':') && !hostname.StartsWith('[')
            ? "[" + hostname + "]" : hostname;
        return "postgres://" + Uri.EscapeDataString(username) + ":" +
            Uri.EscapeDataString(password.Text) + "@" + hostForUri + ":" + portNumber +
            "/" + Uri.EscapeDataString(dbName) + "?sslmode=" +
            Uri.EscapeDataString(sslMode.Text.Trim() is "" ? "disable" : sslMode.Text.Trim());
    }

    private async Task RunDatabaseActionAsync(Button button, string command)
    {
        var dsn = CurrentDsn();
        if (dsn is null) return;
        var create = command == "-create-db";
        if (create && MessageBox.Show(this,
            $"要在 PostgreSQL 上创建数据库“{database.Text.Trim()}”吗？\n" +
            "仅在数据库不存在、当前账号具备 CREATEDB 权限时创建。不会修改已有数据库、账户或密码。",
            "确认创建数据库", MessageBoxButtons.YesNo, MessageBoxIcon.Question) != DialogResult.Yes)
            return;
        var executable = Path.Combine(AppContext.BaseDirectory, "artex.exe");
        if (!File.Exists(executable))
        {
            MessageBox.Show(this, "找不到 artex.exe，请使用完整 Windows 桌面包。",
                "数据库连接", MessageBoxButtons.OK, MessageBoxIcon.Error);
            return;
        }

        button.Enabled = false;
        discoveryStatus.Text = create ? "正在创建目标数据库（已存在时不修改）…" :
            "正在验证 PostgreSQL 账号与数据库连接…";
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
            start.ArgumentList.Add(command);
            // No password in process command line or temporary configuration.
            start.Environment["ARTEX_PG_DSN"] = dsn;
            using var process = Process.Start(start) ??
                throw new InvalidOperationException("无法启动数据库检查程序。");
            using var timeout = CancellationTokenSource.CreateLinkedTokenSource(discoveryLifetime.Token);
            timeout.CancelAfter(TimeSpan.FromSeconds(17));
            try
            {
                await process.WaitForExitAsync(timeout.Token);
            }
            catch (OperationCanceledException)
            {
                if (!process.HasExited) process.Kill(entireProcessTree: true);
                throw;
            }
            var details = process.ExitCode == 0
                ? await process.StandardOutput.ReadToEndAsync(timeout.Token)
                : await process.StandardError.ReadToEndAsync(timeout.Token);
            if (IsDisposed || discoveryLifetime.IsCancellationRequested) return;
            if (process.ExitCode == 0)
            {
                discoveryStatus.Text = create
                    ? (details.Contains("already exists") ? "目标数据库已存在，未做修改。" :
                        "数据库创建成功。可测试连接，再保存配置。")
                    : "PostgreSQL 连接验证成功。可以保存并连接。";
                return;
            }
            if (password.Text.Length > 0)
                details = details.Replace(password.Text, "***", StringComparison.Ordinal)
                    .Replace(Uri.EscapeDataString(password.Text), "***", StringComparison.Ordinal);
            discoveryStatus.Text = create ? "创建数据库未成功。" : "数据库连接失败，请检查配置。";
            MessageBox.Show(this, details.Length > 650 ? details[..650] : details,
                create ? "创建数据库失败" : "数据库连接检查失败",
                MessageBoxButtons.OK, MessageBoxIcon.Warning);
        }
        catch (OperationCanceledException) when (discoveryLifetime.IsCancellationRequested) { }
        catch (OperationCanceledException)
        {
            if (!IsDisposed) discoveryStatus.Text = "操作超时，请检查 PostgreSQL 服务和连接信息。";
        }
        catch (Exception ex)
        {
            if (!IsDisposed) discoveryStatus.Text = "数据库操作未完成：" + ex.Message;
        }
        finally
        {
            if (!IsDisposed) button.Enabled = true;
        }
    }

    private void LoadExistingConfiguration()
    {
        var json = JsonNode.Parse(File.ReadAllText(ConfigPath)) as JsonObject
            ?? throw new InvalidDataException("config.json 必须是有效的 JSON 对象。");
        var db = json["database"] as JsonObject
            ?? throw new InvalidDataException("config.json 的 database 字段应为 JSON 对象。");

        var dsn = db["dsn"]?.ToString();
        if (!string.IsNullOrWhiteSpace(dsn) && Uri.TryCreate(dsn, UriKind.Absolute, out var uri))
        {
            host.Text = uri.Host;
            port.Text = uri.Port.ToString();
            var credentials = uri.UserInfo.Split(':', 2);
            user.Text = Uri.UnescapeDataString(credentials[0]);
            if (credentials.Length == 2) password.Text = Uri.UnescapeDataString(credentials[1]);
            database.Text = Uri.UnescapeDataString(uri.AbsolutePath.TrimStart('/'));
            var ssl = uri.Query.TrimStart('?').Split('&')
                .FirstOrDefault(part => part.StartsWith("sslmode=", StringComparison.OrdinalIgnoreCase));
            sslMode.Text = ssl is null ? "disable" : Uri.UnescapeDataString(ssl[8..]);
        }
        else
        {
            host.Text = db["host"]?.ToString() ?? host.Text;
            port.Text = db["port"]?.ToString() ?? port.Text;
            user.Text = db["user"]?.ToString() ?? user.Text;
            password.Text = db["password"]?.ToString() ?? password.Text;
            database.Text = db["dbname"]?.ToString() ?? database.Text;
            sslMode.Text = db["sslmode"]?.ToString() ?? sslMode.Text;
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

        try
        {
            var path = ConfigPath;
            Directory.CreateDirectory(Path.GetDirectoryName(Path.GetFullPath(path))!);
            if (!editExisting || !File.Exists(path))
            {
                // First launch: do not overwrite a file created by another instance.
                using var output = new FileStream(path, FileMode.CreateNew, FileAccess.Write, FileShare.None);
                JsonSerializer.Serialize(output, new
                {
                    database = new {
                        host = hostname, port = portNumber, user = username,
                        password = password.Text, dbname = dbName, sslmode = sslMode.Text.Trim()
                    },
                    skill_dir = "",
                }, new JsonSerializerOptions { WriteIndented = true });
            }
            else
            {
                // Preserve all other settings (skill_dir and future extensions).
                // A previous database.dsn takes precedence over component fields;
                // remove that key when the operator replaces the connection here.
                var root = JsonNode.Parse(File.ReadAllText(path)) as JsonObject
                    ?? throw new InvalidDataException("配置文件内容不是 JSON 对象。");
                var dbNode = root["database"] as JsonObject
                    ?? throw new InvalidDataException("配置文件的 database 字段不是 JSON 对象。");
                dbNode.Remove("dsn");
                dbNode["host"] = hostname;
                dbNode["port"] = portNumber;
                dbNode["user"] = username;
                dbNode["password"] = password.Text;
                dbNode["dbname"] = dbName;
                dbNode["sslmode"] = sslMode.Text.Trim();
                var tmp = path + ".new";
                File.WriteAllText(tmp, root.ToJsonString(new JsonSerializerOptions { WriteIndented = true }) + "\n");
                File.Move(tmp, path, overwrite: true);
            }
            DialogResult = DialogResult.OK;
            Close();
        }
        catch (JsonException ex)
        {
            MessageBox.Show(this, "配置文件的 JSON 格式无效：" + ex.Message,
                "数据库配置", MessageBoxButtons.OK, MessageBoxIcon.Error);
        }
        catch (InvalidDataException ex)
        {
            MessageBox.Show(this, ex.Message, "数据库配置",
                MessageBoxButtons.OK, MessageBoxIcon.Error);
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
