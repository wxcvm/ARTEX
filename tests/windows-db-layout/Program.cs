using System.Drawing;
using System.Reflection;
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

    private static void Run()
    {
        var type = typeof(MainForm).Assembly.GetType("Artex.Desktop.DatabaseSetupForm", throwOnError: true)!;
        using var dialog = (Form)Activator.CreateInstance(
            type, BindingFlags.Instance | BindingFlags.NonPublic | BindingFlags.Public,
            binder: null, args: new object[] { false }, culture: null)!;
        using var largeFont = new Font(SystemFonts.MessageBoxFont.FontFamily, 18f);

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
            .Single(button => button.Text == name);
        var buttons = new[]
        {
            "检测本机", "启动所选服务", "安装 PostgreSQL",
            "创建目标数据库", "初始化专用账户",
            "稍后配置", "测试数据库连接", "保存并连接",
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
                    if (items[i].Bounds.IntersectsWith(items[j].Bounds))
                        throw new Exception($"Buttons overlap: {items[i].Text}, {items[j].Text}");
        }

        var root = dialog.Controls.OfType<TableLayoutPanel>().Single();
        var scroll = root.Controls.OfType<Panel>().Single(panel => panel.AutoScroll);
        var footer = root.Controls.OfType<FlowLayoutPanel>().Single();
        foreach (var name in new[] { "保存并连接", "测试数据库连接", "稍后配置" })
            if (!footer.Controls.Contains(GetButton(name)))
                throw new Exception("Primary action escaped pinned footer: " + name);

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
        dialog.Close();
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
