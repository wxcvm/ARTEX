using System.Diagnostics;
using System.Net;
using System.Net.Sockets;
using Microsoft.Win32;

namespace Artex.Desktop;

// Probe the PostgreSQL wire protocol, not merely whether TCP port 5432 is open.
// No authentication or database mutation takes place during discovery.
internal static class PostgresDiscovery
{
    internal static async Task<int?> FindLocalPortAsync(CancellationToken cancel = default)
    {
        foreach (var port in new[] { 5432, 5433 })
        {
            cancel.ThrowIfCancellationRequested();
            using var timeout = CancellationTokenSource.CreateLinkedTokenSource(cancel);
            timeout.CancelAfter(TimeSpan.FromSeconds(1));
            try
            {
                using var client = new TcpClient();
                await client.ConnectAsync(IPAddress.Loopback, port, timeout.Token);
                using var network = client.GetStream();
                // PostgreSQL SSLRequest (length=8, request code=80877103).
                // The server responds 'S' (SSL supported) or 'N' (not supported).
                await network.WriteAsync(new byte[] { 0, 0, 0, 8, 4, 210, 22, 47 }, timeout.Token);
                var answer = new byte[1];
                if (await network.ReadAsync(answer, timeout.Token) == 1 &&
                    (answer[0] == 'S' || answer[0] == 'N'))
                    return port;
            }
            catch (OperationCanceledException) when (!cancel.IsCancellationRequested) { }
            catch (SocketException) { }
            catch (IOException) { }
        }
        return null;
    }

    internal static string[] InstalledWindowsServices()
    {
        using var services = Registry.LocalMachine.OpenSubKey(@"SYSTEM\CurrentControlSet\Services");
        return services?.GetSubKeyNames()
            .Where(name => name.StartsWith("postgresql", StringComparison.OrdinalIgnoreCase))
            .OrderBy(name => name, StringComparer.OrdinalIgnoreCase)
            .ToArray() ?? [];
    }

    // Only called after an explicit operator confirmation; never installs,
    // reinitializes, or changes an existing PostgreSQL role or database.
    internal static async Task<bool> StartWindowsServiceAsync(string serviceName, CancellationToken cancel)
    {
        if (!InstalledWindowsServices().Contains(serviceName, StringComparer.OrdinalIgnoreCase))
            throw new InvalidOperationException("未找到所选 PostgreSQL 服务，请重新检测。");

        var start = new ProcessStartInfo("sc.exe")
        {
            UseShellExecute = false,
            CreateNoWindow = true,
        };
        start.ArgumentList.Add("start");
        start.ArgumentList.Add(serviceName);
        using var process = Process.Start(start) ??
            throw new InvalidOperationException("无法启动 Windows 服务控制命令。");
        using var timeout = CancellationTokenSource.CreateLinkedTokenSource(cancel);
        timeout.CancelAfter(TimeSpan.FromSeconds(15));
        await process.WaitForExitAsync(timeout.Token);
        return process.ExitCode == 0;
    }
}
