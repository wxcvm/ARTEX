using System.Diagnostics;
using System.Net;
using System.Net.Sockets;
using System.Security.AccessControl;
using System.Security.Cryptography;
using System.Security.Principal;
using System.Text;
using System.Text.Json;
using Microsoft.Win32;

namespace Artex.Desktop;

// Owns only the per-user cluster recorded in its DPAPI-protected manifest.
// Installed Windows services, external clusters and their credentials are never changed.
internal static class ManagedPostgres
{
    internal static string Root { get; set; } = Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData), "ARTEX", "PostgreSQL", "16");
    private static string Manifest => Path.Combine(Root, "cluster.dpapi");
    private static string Data => Path.Combine(Root, "data");
    private static readonly SemaphoreSlim Gate = new(1, 1);
    private static bool startedHere;
    [System.Runtime.InteropServices.DllImport("kernel32.dll", CharSet = System.Runtime.InteropServices.CharSet.Unicode, SetLastError = true)]
    private static extern uint GetShortPathName(string path, StringBuilder buffer, uint length);
    private static string CliDirectory()
    {
        // PostgreSQL 16's Windows CLI converts its absolute child-process paths
        // through the system ANSI code page. Keep the actual data in the user's
        // Unicode directory, but launch through its existing NTFS short name.
        var buffer = new StringBuilder(32768);
        var size = GetShortPathName(Root, buffer, (uint)buffer.Capacity);
        return size > 0 && size < buffer.Capacity ? buffer.ToString() : Root;
    }
    internal sealed record Cluster(string Bin, int Port, string AdminPassword, string AppPassword)
    {
        internal string Dsn(bool admin = false) => $"postgres://{(admin ? "postgres" : "artex")}:{(admin ? AdminPassword : AppPassword)}@127.0.0.1:{Port}/artex?sslmode=disable";
    }

    internal static string? FindBinaries()
    {
        var candidates = new List<string>();
        using var installs = Registry.LocalMachine.OpenSubKey(@"SOFTWARE\PostgreSQL\Installations");
        foreach (var key in installs?.GetSubKeyNames() ?? [])
        {
            using var install = installs!.OpenSubKey(key);
            if (install?.GetValue("Version") is string version && version.StartsWith("16.", StringComparison.Ordinal) &&
                install.GetValue("Base Directory") is string path) candidates.Add(Path.Combine(path, "bin"));
        }
        candidates.Add(Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.ProgramFiles), "PostgreSQL", "16", "bin"));
        candidates.Add(Path.Combine(AppContext.BaseDirectory, "postgresql", "bin"));
        return candidates.FirstOrDefault(p => File.Exists(Path.Combine(p, "initdb.exe")) &&
            File.Exists(Path.Combine(p, "pg_ctl.exe")) && File.Exists(Path.Combine(p, "postgres.exe")));
    }

    private static Cluster? ReadManifest() => !File.Exists(Manifest) ? null :
        JsonSerializer.Deserialize<Cluster>(ProtectedData.Unprotect(File.ReadAllBytes(Manifest), null, DataProtectionScope.CurrentUser))
        ?? throw new InvalidDataException("独立数据库记录无效，请打开数据库日志目录检查；不会重置已有数据。");

    internal static async Task<Cluster> PrepareAsync(string backend, IProgress<string> progress, CancellationToken cancel)
    {
        await Gate.WaitAsync(cancel);
        try
        {
            using var identity = WindowsIdentity.GetCurrent();
            if (new WindowsPrincipal(identity).IsInRole(WindowsBuiltInRole.Administrator))
                throw new InvalidOperationException("请正常打开 ARTEX，不要使用“以管理员身份运行”。独立数据库使用当前普通用户身份准备和运行。");
            Directory.CreateDirectory(Root);
            // Password-file and manifest inherit an owner-only ACL from this directory.
            var acl = new DirectorySecurity();
            acl.SetAccessRuleProtection(isProtected: true, preserveInheritance: false);
            acl.AddAccessRule(new FileSystemAccessRule(identity.User!, FileSystemRights.FullControl,
                InheritanceFlags.ContainerInherit | InheritanceFlags.ObjectInherit, PropagationFlags.None, AccessControlType.Allow));
            new DirectoryInfo(Root).SetAccessControl(acl);
            using var lease = new FileStream(Path.Combine(Root, "prepare.lock"), FileMode.OpenOrCreate, FileAccess.ReadWrite, FileShare.None);
            var cluster = ReadManifest();
            if (cluster is null)
            {
                if (Directory.Exists(Data) && Directory.EnumerateFileSystemEntries(Data).Any())
                    throw new InvalidOperationException("发现未登记的非空数据目录，未进行初始化。请先保留该目录，使用高级设置连接或检查日志。");
                var bin = FindBinaries() ?? throw new FileNotFoundException("未找到完整的 PostgreSQL 程序。请点“安装 PostgreSQL”，完成官方安装后回到这里重试。");
                await RunAsync(Path.Combine(bin, "postgres.exe"), ["--version"], cancel, expectedText: "PostgreSQL) 16.");
                cluster = new Cluster(bin, ChoosePort(), Convert.ToHexString(RandomNumberGenerator.GetBytes(24)), Convert.ToHexString(RandomNumberGenerator.GetBytes(24)));
                var encrypted = ProtectedData.Protect(JsonSerializer.SerializeToUtf8Bytes(cluster), null, DataProtectionScope.CurrentUser);
                File.WriteAllBytes(Manifest + ".new", encrypted);
                File.Move(Manifest + ".new", Manifest, overwrite: true);
            }
            if (!File.Exists(Path.Combine(cluster.Bin, "pg_ctl.exe")))
                throw new FileNotFoundException("独立数据库使用的 PostgreSQL 程序已被移动。请恢复原程序位置后重试，数据不会被重新初始化。");
            if (!File.Exists(Path.Combine(Data, "global", "pg_control")))
            {
                progress.Report("1/5 正在准备独立数据库目录；原有 PostgreSQL 服务不变…");
                if (Directory.Exists(Data) && Directory.EnumerateFileSystemEntries(Data).Any())
                {
                    // Preserve incomplete attempts instead of recursively deleting anything.
                    var preserved = Data + ".incomplete-" + Guid.NewGuid().ToString("N");
                    if (!Path.GetFullPath(preserved).StartsWith(Path.GetFullPath(Root) + Path.DirectorySeparatorChar, StringComparison.OrdinalIgnoreCase))
                        throw new IOException("数据目录不在预期位置。");
                    Directory.Move(Data, preserved);
                }
                var passwordFile = Path.Combine(Root, "init-password.tmp");
                try
                {
                    File.WriteAllText(passwordFile, cluster.AdminPassword + "\n", new UTF8Encoding(false));
                    await RunAsync(Path.Combine(cluster.Bin, "initdb.exe"), ["-D", "data", "-U", "postgres", "--auth=scram-sha-256", "--encoding=UTF8", "--locale=C", "--pwfile=init-password.tmp"], cancel, workingDirectory: CliDirectory());
                }
                finally { File.Delete(passwordFile); }
            }
            progress.Report("2/5 正在启动仅本机可访问的数据库…");
            await EnsureStartedCoreAsync(cluster, cancel);
            progress.Report("3/5 正在准备 ARTEX 专用账户及数据库…");
            await RunAsync(backend, ["-setup-app-db"], cancel, new Dictionary<string, string>
            {
                ["ARTEX_PG_DSN"] = cluster.Dsn(admin: true),
                ["ARTEX_APP_DB_USER"] = "artex", ["ARTEX_APP_DB_PASSWORD"] = cluster.AppPassword,
            }, secrets: [cluster.AdminPassword, cluster.AppPassword]);
            progress.Report("4/5 正在验证专用账户连接…");
            await RunAsync(backend, ["-check-db"], cancel, new() { ["ARTEX_PG_DSN"] = cluster.Dsn() }, secrets: [cluster.AppPassword]);
            progress.Report("5/5 数据库已就绪，正在保存配置…");
            return cluster;
        }
        finally { Gate.Release(); }
    }

    internal static async Task EnsureForConfigurationAsync(string dsn, CancellationToken cancel)
    {
        var c = ReadManifest();
        if (c is null || !Uri.TryCreate(dsn, UriKind.Absolute, out var uri) || uri.Host != "127.0.0.1" || uri.Port != c.Port) return;
        await Gate.WaitAsync(cancel);
        try { await EnsureStartedCoreAsync(c, cancel); }
        finally { Gate.Release(); }
    }

    private static async Task EnsureStartedCoreAsync(Cluster c, CancellationToken cancel)
    {
        var ctl = Path.Combine(c.Bin, "pg_ctl.exe");
        if (await RunAsync(ctl, ["status", "-D", "data"], cancel, allowFailure: true, workingDirectory: CliDirectory()) == 0) return;
        await RunAsync(ctl, ["start", "-D", "data", "-l", "postgresql.log", "-w", "-t", "30", "-o", $"-h 127.0.0.1 -p {c.Port}"], cancel, captureOutput: false, workingDirectory: CliDirectory());
        startedHere = true;
    }

    internal static void StopOwned()
    {
        if (!startedHere) return;
        try
        {
            var c = ReadManifest();
            // FormClosed runs on the WinForms context. Do not block that context
            // while the async process waiter tries to resume back onto it.
            if (c is not null) Task.Run(() => RunAsync(Path.Combine(c.Bin, "pg_ctl.exe"), ["stop", "-D", "data", "-m", "fast", "-w", "-t", "15"], CancellationToken.None, workingDirectory: CliDirectory())).GetAwaiter().GetResult();
            startedHere = false;
        }
        catch (Exception ex) { Debug.WriteLine("独立数据库未能正常停止：" + ex.GetType().Name); }
    }

    internal static int ChoosePort()
    {
        for (var port = 15432; port <= 15442; port++)
        {
            var listener = new TcpListener(IPAddress.Loopback, port);
            try { listener.Start(); return port; }
            catch (SocketException) { }
            finally { listener.Stop(); }
        }
        throw new IOException("本机 15432–15442 端口均被占用。请关闭占用程序后重试，或使用高级设置连接已有数据库。");
    }

    internal static async Task<int> RunAsync(string executable, string[] arguments, CancellationToken cancel,
        Dictionary<string, string>? environment = null, bool allowFailure = false, string? expectedText = null, string[]? secrets = null, bool captureOutput = true, string? workingDirectory = null)
    {
        var info = new ProcessStartInfo(executable) { UseShellExecute = false, CreateNoWindow = true, RedirectStandardOutput = captureOutput, RedirectStandardError = captureOutput, WorkingDirectory = workingDirectory ?? AppContext.BaseDirectory };
        foreach (var a in arguments) info.ArgumentList.Add(a);
        if (environment is not null) foreach (var item in environment) info.Environment[item.Key] = item.Value;
        using var process = Process.Start(info) ?? throw new IOException("无法启动数据库程序。");
        using var timeout = CancellationTokenSource.CreateLinkedTokenSource(cancel);
        timeout.CancelAfter(TimeSpan.FromSeconds(60));
        var stdout = captureOutput ? process.StandardOutput.ReadToEndAsync() : Task.FromResult("");
        var stderr = captureOutput ? process.StandardError.ReadToEndAsync() : Task.FromResult("");
        try { await process.WaitForExitAsync(timeout.Token); }
        catch (OperationCanceledException) { if (!process.HasExited) process.Kill(entireProcessTree: true); throw; }
        var output = await stdout + "\n" + await stderr;
        if ((!allowFailure && process.ExitCode != 0) || (expectedText is not null && !output.Contains(expectedText, StringComparison.Ordinal)))
        {
            foreach (var secret in secrets ?? []) if (secret.Length > 0) output = output.Replace(secret, "***", StringComparison.Ordinal).Replace(Uri.EscapeDataString(secret), "***", StringComparison.Ordinal);
            throw new IOException("数据库操作未完成：" + output[..Math.Min(output.Length, 1200)] + "\n日志目录：" + Root);
        }
        return process.ExitCode;
    }
}
