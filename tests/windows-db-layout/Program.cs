using System.Drawing;
using System.Reflection;
using System.Text.Json;
using System.Text.Json.Nodes;
using System.Windows.Forms;
using Artex.Desktop;

// Run only on a Windows desktop CI runner. This is a real WinForms layout test,
// not a source-code assertion or a screenshot approximation.
internal static class Program
{
    [STAThread]
    private static int Main()
    {
        try
        {
            ApplicationConfiguration.Initialize();
            CheckHealthResponseShape();
            Run();
            Console.WriteLine("[ARTEX] Large-font/narrow-window database layout and password key repeat passed.");
            return 0;
        }
        catch (Exception error)
        {
            Console.Error.WriteLine(error);
            return 1;
        }
    }

    private static void CheckHealthResponseShape()
    {
        var check = typeof(MainForm).GetMethod("IsHealthDocument", BindingFlags.Static | BindingFlags.NonPublic)!;
        foreach (var (json, expected) in new (string, bool)[]
        {
            ("""{"service":"artex","ok":true,"version":"test"}""", true),
            ("""{"service":"other","ok":true}""", false),
            ("""{"service":"artex","ok":false}""", false),
            ("""{"service":"artex","ok":"true"}""", false),
            ("""{"service":"artex"}""", false),
            ("""{"ok":true}""", false),
            ("""{"service":123,"ok":true}""", false),
            ("""{"service":{},"ok":true}""", false),
            ("""{"service":null,"ok":true}""", false),
            ("[]", false),
            ("null", false),
            ("true", false),
            ("123", false),
        })
        {
            using var document = JsonDocument.Parse(json);
            if (!Equals(check.Invoke(null, new object[] { document.RootElement }), expected))
                throw new Exception("Incorrect desktop health result for: " + json);
        }
        Console.WriteLine("[ARTEX] Desktop health response shape regression passed.");
    }

    private static void Run()
    {
        var type = typeof(MainForm).Assembly.GetType("Artex.Desktop.DatabaseSetupForm", throwOnError: true)!;
        using var dialog = (Form)Activator.CreateInstance(
            type, BindingFlags.Instance | BindingFlags.NonPublic | BindingFlags.Public,
            binder: null, args: new object[] { false }, culture: null)!;
        var fontSize = float.TryParse(Environment.GetEnvironmentVariable("ARTEX_TEST_FONT_SIZE"), System.Globalization.CultureInfo.InvariantCulture, out var requested) ? requested : 18f;
        using var largeFont = new Font(SystemFonts.MessageBoxFont.FontFamily, fontSize);

        dialog.Font = largeFont;
        dialog.ClientSize = new Size(680, 460);
        dialog.StartPosition = FormStartPosition.Manual;
        dialog.Location = new Point(
            Screen.PrimaryScreen!.WorkingArea.Left + 10,
            Screen.PrimaryScreen.WorkingArea.Top + 10);
        dialog.Show();
        Application.DoEvents();
        dialog.PerformLayout();
        Application.DoEvents();

        var children = AllChildren(dialog).ToArray();
        Button GetButton(string name) => children.OfType<Button>()
            .First(button => button.Text == name);
        var buttons = new[]
        {
            "检测本机", "启动所选服务", "安装 PostgreSQL",
            "创建目标数据库", "初始化专用账户",
            "稍后配置", "测试数据库连接", "自动准备本机数据库",
        }.Select(GetButton).ToArray();

        foreach (var button in buttons)
        {
            if (!button.AutoSize)
                throw new Exception("Button did not autosize: " + button.Text);
            var textWidth = TextRenderer.MeasureText(button.Text, button.Font).Width;
            if (button.Width < textWidth)
                throw new Exception($"Clipped text: {button.Text}, {button.Width}px < {textWidth}px");
        }

        foreach (var panel in children.OfType<FlowLayoutPanel>())
        {
            var items = panel.Controls.OfType<Button>().ToArray();
            for (var i = 0; i < items.Length; i++)
                for (var j = i + 1; j < items.Length; j++)
                    if (items[i].Visible && items[j].Visible && items[i].Bounds.IntersectsWith(items[j].Bounds))
                        throw new Exception($"Buttons overlap: {items[i].Text}, {items[j].Text}");
        }

        var root = dialog.Controls.OfType<TableLayoutPanel>().Single();
        var scroll = root.Controls.OfType<Panel>().Single(panel => panel.AutoScroll);
        var footer = root.Controls.OfType<FlowLayoutPanel>().Single();
        foreach (var name in new[] { "自动准备本机数据库", "测试数据库连接", "稍后配置" })
            if (!footer.Controls.Contains(GetButton(name)))
                throw new Exception("Primary action escaped pinned footer: " + name);

        // Recommended setup must be simple by default, with advanced fields
        // and the original save/test operations available on demand.
        var radioButtons = children.OfType<RadioButton>().ToArray();
        var recommended = radioButtons.Single(r => r.Text == "自动准备本机数据库 · 推荐");
        var custom = radioButtons.Single(r => r.Text == "已有数据库 / 高级设置");
        if (!recommended.Checked || GetButton("自动准备本机数据库").Visible == false)
            throw new Exception("Recommended setup is not the default");
        if (children.OfType<TextBox>().Any(box => box.UseSystemPasswordChar && box.Visible))
            throw new Exception("Default setup should not require a database password");
        var screenshot = Environment.GetEnvironmentVariable("ARTEX_LAYOUT_SCREENSHOT");
        if (!string.IsNullOrWhiteSpace(screenshot))
        {
            using var bitmap = new Bitmap(dialog.Width, dialog.Height);
            dialog.DrawToBitmap(bitmap, new Rectangle(Point.Empty, bitmap.Size));
            bitmap.Save(screenshot);
        }
        custom.Checked = true;
        Application.DoEvents();
        if (!GetButton("验证并保存连接").Visible || !GetButton("测试数据库连接").Visible)
            throw new Exception("Custom database options were not restored");
        recommended.Checked = true;
        Application.DoEvents();
        if (!GetButton("自动准备本机数据库").Visible)
            throw new Exception("Recommended setup was not restored");

        var frame = dialog.RectangleToScreen(dialog.ClientRectangle);
        var bottom = footer.RectangleToScreen(footer.ClientRectangle);
        if (bottom.Bottom > frame.Bottom || bottom.Top < frame.Top)
            throw new Exception("Pinned footer is outside client bounds");
        if (scroll.AutoScrollPosition.Y != 0)
            throw new Exception("Wizard opens scrolled past introduction");

        var password = children.OfType<TextBox>().Single(box => box.UseSystemPasswordChar);
        var down = typeof(Control).GetMethod("OnKeyDown", BindingFlags.Instance | BindingFlags.NonPublic)!;
        var up = typeof(Control).GetMethod("OnKeyUp", BindingFlags.Instance | BindingFlags.NonPublic)!;
        var initial = new KeyEventArgs(Keys.A);
        down.Invoke(password, new object[] { initial });
        var repeated = new KeyEventArgs(Keys.A);
        down.Invoke(password, new object[] { repeated });
        if (initial.SuppressKeyPress || !repeated.SuppressKeyPress)
            throw new Exception("Password held-character suppression failed");
        up.Invoke(password, new object[] { new KeyEventArgs(Keys.A) });
        var fresh = new KeyEventArgs(Keys.A);
        down.Invoke(password, new object[] { fresh });
        if (fresh.SuppressKeyPress)
            throw new Exception("Password key did not recover after KeyUp");
        up.Invoke(password, new object[] { new KeyEventArgs(Keys.A) });

        var paste = new KeyEventArgs(Keys.Control | Keys.V);
        var pasteRepeat = new KeyEventArgs(Keys.Control | Keys.V);
        down.Invoke(password, new object[] { paste });
        down.Invoke(password, new object[] { pasteRepeat });
        if (paste.SuppressKeyPress || !pasteRepeat.SuppressKeyPress)
            throw new Exception("Held Ctrl+V was not suppressed");
        up.Invoke(password, new object[] { new KeyEventArgs(Keys.Control | Keys.V) });

        // Never offer new-account provisioning as the default when a process
        // environment variable would override the saved connection anyway.
        var previousDsn = Environment.GetEnvironmentVariable("ARTEX_PG_DSN");
        try
        {
            Environment.SetEnvironmentVariable("ARTEX_PG_DSN", "postgres://example@127.0.0.1:5432/example");
            using var overridden = (Form)Activator.CreateInstance(
                type, BindingFlags.Instance | BindingFlags.NonPublic | BindingFlags.Public,
                binder: null, args: new object[] { false }, culture: null)!;
            if (!AllChildren(overridden).OfType<RadioButton>()
                .Single(r => r.Text == "已有数据库 / 高级设置").Checked)
                throw new Exception("Environment override should default to custom connection mode");
            var note = (Label)type.GetField("overrideNotice", BindingFlags.NonPublic | BindingFlags.Instance)!
                .GetValue(overridden)!;
            if (!note.Text.Contains("ARTEX_PG_DSN", StringComparison.Ordinal))
                throw new Exception("Database connection override is not explained");
        }
        finally
        {
            Environment.SetEnvironmentVariable("ARTEX_PG_DSN", previousDsn);
        }
        dialog.Close();
        CheckConfigRepair(type);
    }

    private static void CheckConfigRepair(Type type)
    {
        var fixture = Path.Combine(AppContext.BaseDirectory, "ARTEX-layout-config-" + Guid.NewGuid().ToString("N") + ".json");
        var previousConfig = Environment.GetEnvironmentVariable("ARTEX_CONFIG");
        var previousDsn = Environment.GetEnvironmentVariable("ARTEX_PG_DSN");
        var fields = BindingFlags.NonPublic | BindingFlags.Instance;
        Form RepairForm() => (Form)Activator.CreateInstance(
            type, fields, binder: null, args: new object[] { true }, culture: null)!;
        try
        {
            Environment.SetEnvironmentVariable("ARTEX_CONFIG", fixture);
            Environment.SetEnvironmentVariable("ARTEX_PG_DSN", null);
            const string initial = "{\"skill_dir\":\"keep-me\"}";
            File.WriteAllText(fixture, initial);
            var save = type.GetMethod("Save", fields)!;
            var changed = type.GetProperty("ConfigurationChanged", fields)!;
            var confirmations = 0;
            Func<IWin32Window, bool> refuse = _ => { confirmations++; return false; };
            using (var refused = (Form)Activator.CreateInstance(
                type, fields, binder: null, args: new object[] { true, refuse }, culture: null)!)
            {
                save.Invoke(refused, null);
                if (confirmations != 1 || changed.GetValue(refused) is true || File.ReadAllText(fixture) != initial)
                    throw new Exception("Declining confirmation changed saved configuration");
            }
            using (var repair = RepairForm())
            {
                // Exercise only JSON persistence with a temporary file. Connection
                // validation belongs to the separate managed PostgreSQL integration.
                save.Invoke(repair, null);
                if (changed.GetValue(repair) is not true)
                    throw new Exception("Saving a missing database object must persist a change");
            }
            var saved = JsonNode.Parse(File.ReadAllText(fixture))!.AsObject();
            if (saved["skill_dir"]!.GetValue<string>() != "keep-me" || saved["database"] is not JsonObject)
                throw new Exception("Database repair did not preserve unrelated configuration");

            const string malformed = "{\"database\":";
            File.WriteAllText(fixture, malformed);
            using var invalid = RepairForm();
            if (type.GetField("configurationLoadError", fields)!.GetValue(invalid) is not string)
                throw new Exception("Malformed config has no actionable repair state");
            if (File.ReadAllText(fixture) != malformed)
                throw new Exception("Opening repair overwrote malformed configuration");

            Environment.SetEnvironmentVariable("ARTEX_PG_DSN", "postgresql://fixture@example.invalid/app");
            using var overridden = RepairForm();
            var primary = (Button)type.GetField("saveButton", fields)!.GetValue(overridden)!;
            if (primary.Text != "验证环境变量连接" || changed.GetValue(overridden) is true ||
                !AllChildren(overridden).OfType<LinkLabel>().Any(link => link.Text.Contains("打开配置目录修复")))
                throw new Exception("External override and malformed config must offer separate verification and repair");
        }
        finally
        {
            Environment.SetEnvironmentVariable("ARTEX_CONFIG", previousConfig);
            Environment.SetEnvironmentVariable("ARTEX_PG_DSN", previousDsn);
            File.Delete(fixture);
            File.Delete(fixture + ".new");
        }
    }

    private static IEnumerable<Control> AllChildren(Control parent)
    {
        foreach (Control child in parent.Controls)
        {
            yield return child;
            foreach (var descendant in AllChildren(child))
                yield return descendant;
        }
    }
}
