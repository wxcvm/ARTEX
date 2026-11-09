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

    internal DatabaseSetupForm()
    {
        Text = "ARTEX · 首次数据库配置";
        Icon = System.Drawing.Icon.ExtractAssociatedIcon(Application.ExecutablePath) ?? SystemIcons.Application;
        StartPosition = FormStartPosition.CenterParent;
        ClientSize = new Size(550, 385);
        MinimumSize = new Size(500, 380);
        FormBorderStyle = FormBorderStyle.FixedDialog;
        MaximizeBox = false;
        MinimizeBox = false;

        var layout = new TableLayoutPanel
        {
            Dock = DockStyle.Fill,
            ColumnCount = 2,
            RowCount = 9,
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

        var configLabel = new Label
        {
            Text = "配置保存位置：" + ConfigPath,
            AutoSize = false,
            AutoEllipsis = true,
            Dock = DockStyle.Fill,
            ForeColor = SystemColors.GrayText,
        };
        layout.Controls.Add(configLabel, 0, 7);
        layout.SetColumnSpan(configLabel, 2);

        var buttons = new FlowLayoutPanel
        {
            Dock = DockStyle.Fill,
            FlowDirection = FlowDirection.RightToLeft,
            WrapContents = false,
        };
        var save = new Button { Text = "保存并连接", Width = 110, Height = 30 };
        save.Click += (_, _) => Save();
        var cancel = new Button { Text = "稍后配置", Width = 100, Height = 30, DialogResult = DialogResult.Cancel };
        buttons.Controls.Add(save);
        buttons.Controls.Add(cancel);
        layout.Controls.Add(buttons, 0, 8);
        layout.SetColumnSpan(buttons, 2);

        AcceptButton = save;
        CancelButton = cancel;
        Controls.Add(layout);
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
