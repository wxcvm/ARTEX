using System.Threading;
using System.Windows.Forms;

namespace Artex.Desktop;

internal static class Program
{
    [STAThread]
    private static void Main()
    {
        ApplicationConfiguration.Initialize();

        // Prevent two windows from unintentionally managing the same
        // local service. Windows sessions have independent "Local" namespaces.
        using var singleInstance = new Mutex(true, @"Local\ARTEX.Desktop", out var firstInstance);
        if (!firstInstance)
        {
            MessageBox.Show("ARTEX 桌面客户端已经在运行。可从系统托盘打开。", "ARTEX",
                MessageBoxButtons.OK, MessageBoxIcon.Information);
            return;
        }

        Application.Run(new MainForm());
    }
}
