using System.Drawing;
using System.Windows.Forms;

namespace Artex.Desktop;

// Native startup and recovery UI. It stays usable even when the backend or WebView2 cannot start.
internal sealed class StartupView : UserControl
{
    private static readonly Color Background = Color.FromArgb(15, 23, 36);
    private static readonly Color Surface = Color.FromArgb(25, 38, 57);
    private static readonly Color Primary = Color.FromArgb(71, 132, 224);
    private static readonly Color TextColor = Color.FromArgb(236, 242, 251);
    private static readonly Color Muted = Color.FromArgb(171, 190, 215);

    private readonly Panel card = new() { BackColor = Surface };
    private readonly TableLayoutPanel content = new()
    {
        Dock = DockStyle.Top,
        AutoSize = true,
        AutoSizeMode = AutoSizeMode.GrowAndShrink,
        ColumnCount = 1,
        RowCount = 8,
        Padding = new Padding(30, 26, 30, 24),
        BackColor = Surface,
    };
    private readonly Label eyebrow = new()
    {
        AutoSize = true, Text = "ARTEX  /  DESKTOP", ForeColor = Color.FromArgb(119, 173, 255),
        Margin = new Padding(0, 0, 0, 14),
    };
    private readonly Label heading = new()
    {
        AutoSize = true, Font = new Font("Microsoft YaHei UI", 21f, FontStyle.Bold),
        ForeColor = TextColor, Margin = new Padding(0, 0, 0, 10),
    };
    private readonly Label description = new()
    {
        AutoSize = true, ForeColor = Muted, Margin = new Padding(0, 0, 0, 12),
        Font = new Font("Microsoft YaHei UI", 10f),
    };
    private readonly Label guidance = new()
    {
        AutoSize = true, ForeColor = Color.FromArgb(192, 211, 237),
        Margin = new Padding(0, 2, 0, 20),
    };
    private readonly FlowLayoutPanel actions = new()
    {
        AutoSize = true, AutoSizeMode = AutoSizeMode.GrowAndShrink,
        FlowDirection = FlowDirection.LeftToRight, WrapContents = true,
        Margin = new Padding(0, 2, 0, 14),
    };
    private readonly Button repair = MakeButton("修复数据库连接", true, 170);
    private readonly Button retry = MakeButton("重新连接", false, 122);
    private readonly Button log = MakeButton("查看日志", false, 112);
    private readonly LinkLabel detailToggle = new()
    {
        Text = "查看技术详情", AutoSize = true,
        LinkColor = Color.FromArgb(138, 187, 255), ActiveLinkColor = Color.White,
        VisitedLinkColor = Color.FromArgb(138, 187, 255),
        Margin = new Padding(0, 4, 0, 4),
    };
    private readonly TextBox diagnostics = new()
    {
        Multiline = true, ReadOnly = true, BorderStyle = BorderStyle.None,
        ScrollBars = ScrollBars.Vertical, Height = 118,
        BackColor = Color.FromArgb(17, 28, 44),
        ForeColor = Muted,
        Font = new Font("Consolas", 9.5f), Visible = false,
        Margin = new Padding(0, 6, 0, 0),
    };
    private bool arranging;

    internal event EventHandler? RepairRequested;
    internal event EventHandler? RetryRequested;
    internal event EventHandler? LogsRequested;

    internal StartupView()
    {
        Dock = DockStyle.Fill;
        BackColor = Background;
        AutoScroll = true;
        Font = new Font("Microsoft YaHei UI", 10f);
        for (var row = 0; row < content.RowCount; row++)
            content.RowStyles.Add(new RowStyle(SizeType.AutoSize));
        content.ColumnStyles.Add(new ColumnStyle(SizeType.Percent, 100));
        content.Controls.Add(eyebrow, 0, 0);
        content.Controls.Add(heading, 0, 1);
        content.Controls.Add(description, 0, 2);
        content.Controls.Add(guidance, 0, 3);
        content.Controls.Add(actions, 0, 4);
        content.Controls.Add(new Label
        {
            AutoSize = true,
            Text = "详细配置和诊断日志仅在需要时展开。",
            ForeColor = Color.FromArgb(127, 146, 173),
            Font = new Font("Microsoft YaHei UI", 8.5f),
            Margin = new Padding(0, 0, 0, 6),
        }, 0, 5);
        content.Controls.Add(detailToggle, 0, 6);
        content.Controls.Add(diagnostics, 0, 7);
        actions.Controls.Add(repair);
        actions.Controls.Add(retry);
        actions.Controls.Add(log);
        repair.Click += (_, _) => RepairRequested?.Invoke(this, EventArgs.Empty);
        retry.Click += (_, _) => RetryRequested?.Invoke(this, EventArgs.Empty);
        log.Click += (_, _) => LogsRequested?.Invoke(this, EventArgs.Empty);
        detailToggle.LinkClicked += (_, _) =>
        {
            diagnostics.Visible = !diagnostics.Visible;
            detailToggle.Text = diagnostics.Visible ? "收起技术详情" : "查看技术详情";
            ArrangeCard();
        };
        card.Controls.Add(content);
        Controls.Add(card);
        Resize += (_, _) => ArrangeCard();
        content.SizeChanged += (_, _) => ArrangeCard();
        Display("正在启动 ARTEX", "正在连接本机服务并准备界面…", false);
    }

    private static Button MakeButton(string label, bool primary, int width)
    {
        var button = new Button
        {
            Text = label, Width = width, Height = 42,
            FlatStyle = FlatStyle.Flat,
            BackColor = primary ? Primary : Color.FromArgb(37, 55, 80),
            ForeColor = Color.White,
            Cursor = Cursors.Hand,
            Font = new Font("Microsoft YaHei UI", 9.5f, primary ? FontStyle.Bold : FontStyle.Regular),
            Margin = new Padding(0, 0, 10, 8),
        };
        button.FlatAppearance.BorderSize = primary ? 0 : 1;
        button.FlatAppearance.BorderColor = Color.FromArgb(67, 93, 127);
        return button;
    }

    internal void Display(string title, string message, bool canRetry,
        bool databaseProblem = false, string? technical = null, string? nextStep = null)
    {
        heading.Text = title;
        description.Text = message;
        guidance.Text = nextStep ?? (canRetry
            ? "你可以直接重试；如果问题持续，打开日志查看具体原因。"
            : "这个过程可能需要片刻，完成后会自动进入工作台。");
        eyebrow.Text = canRetry ? "ARTEX  /  需要处理" : "ARTEX  /  正在准备";
        repair.Visible = canRetry && databaseProblem;
        retry.Visible = canRetry;
        log.Visible = canRetry;
        detailToggle.Visible = canRetry && !string.IsNullOrWhiteSpace(technical);
        diagnostics.Text = technical ?? "";
        diagnostics.Visible = false;
        detailToggle.Text = "查看技术详情";
        Visible = true;
        BringToFront();
        ArrangeCard();
    }

    private void ArrangeCard()
    {
        if (arranging || IsDisposed) return;
        arranging = true;
        try
        {
            var width = Math.Min(700, Math.Max(300, ClientSize.Width - 48));
            card.Width = width;
            content.Width = width;
            var textWidth = Math.Max(220, width - content.Padding.Horizontal - 6);
            foreach (var label in new[] { eyebrow, heading, description, guidance })
                label.MaximumSize = new Size(textWidth, 0);
            actions.Width = textWidth;
            diagnostics.Width = textWidth;
            content.PerformLayout();
            card.Height = content.PreferredSize.Height;
            card.Left = Math.Max(12, (ClientSize.Width - card.Width) / 2);
            card.Top = Math.Max(22, (ClientSize.Height - card.Height) / 2);
        }
        finally { arranging = false; }
    }
}
