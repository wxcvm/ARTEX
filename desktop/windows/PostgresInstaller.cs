using System.ComponentModel;
using System.Diagnostics;

namespace Artex.Desktop;

internal static class PostgresInstaller
{
    internal const string OfficialDownloadPage = "https://www.postgresql.org/download/windows/";
    internal const string WingetPackage = "PostgreSQL.PostgreSQL.16";

    // Launch the publisher's interactive installer in a separate console.
    // Winget verifies the package according to its manifest. ARTEX never
    // downloads or runs a hard-coded third-party EXE or supplies passwords.
    internal static bool LaunchInstallation()
    {
        if (!IsWingetAvailable())
        {
            Process.Start(new ProcessStartInfo(OfficialDownloadPage) { UseShellExecute = true });
            return false;
        }
        Process.Start(new ProcessStartInfo("cmd.exe")
        {
            UseShellExecute = true,
            Arguments = "/k winget install --id " + WingetPackage +
                " --exact --interactive --source winget" +
                " --accept-package-agreements --accept-source-agreements",
        }) ?? throw new Win32Exception("无法启动 Windows 安装命令。");
        return true;
    }

    private static bool IsWingetAvailable()
    {
        var path = Environment.GetEnvironmentVariable("PATH");
        if (string.IsNullOrWhiteSpace(path)) return false;
        foreach (var folder in path.Split(Path.PathSeparator))
        {
            try
            {
                if (File.Exists(Path.Combine(
                    Environment.ExpandEnvironmentVariables(folder.Trim('"')), "winget.exe")))
                    return true;
            }
            catch (ArgumentException) { }
        }
        return false;
    }
}
