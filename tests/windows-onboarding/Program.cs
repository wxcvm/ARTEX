using System.Diagnostics;
using System.Net;
using System.Net.Sockets;
using System.Reflection;
using Artex.Desktop;

internal static class Program
{
    private sealed class UnpumpedContext : SynchronizationContext
    {
        public override void Post(SendOrPostCallback callback, object? state) { }
    }
    private static readonly Assembly Desktop = typeof(MainForm).Assembly;
    private static readonly Type Manager = Desktop.GetType("Artex.Desktop.ManagedPostgres", true)!;
    private static object? Call(string method, params object[] args) => Manager.GetMethod(method, BindingFlags.Static | BindingFlags.NonPublic)!.Invoke(null, args);
    private static void Check(bool value, string message) { if (!value) throw new Exception(message); }
    private static async Task<object> Prepare(string backend)
    {
        var task = (Task)Call("PrepareAsync", backend, new Progress<string>(s => Console.WriteLine(s)), CancellationToken.None)!;
        await task;
        return task.GetType().GetProperty("Result")!.GetValue(task)!;
    }
    private static T Field<T>(object record, string name) => (T)record.GetType().GetProperty(name)!.GetValue(record)!;

    private static async Task<int> Main(string[] args)
    {
        try
        {
            var bridge = Desktop.GetType("Artex.Desktop.SetupTokenBridge", true)!.GetMethod("IsSetupDocument", BindingFlags.Static | BindingFlags.NonPublic)!;
            foreach (var uri in new[] { "http://127.0.0.1:8787/setup", "http://127.0.0.1:8787/setup/" })
                Check((bool)bridge.Invoke(null, [uri])!, "Valid setup document rejected");
            foreach (var uri in new[] { "https://example.com/setup", "http://localhost:8787/setup", "http://127.0.0.1:8788/setup", "http://127.0.0.1:8787/welcome", "http://127.0.0.1:8787/setup?next=external", "http://127.0.0.1:8787/setup#external", "http://x@127.0.0.1:8787/setup" })
                Check(!(bool)bridge.Invoke(null, [uri])!, "Unexpected setup token recipient: " + uri);
            var listener = new TcpListener(IPAddress.Loopback, 15432);
            try
            {
                try { listener.Start(); }
                catch (SocketException ex) when (ex.SocketErrorCode == SocketError.AddressAlreadyInUse) { }
                Check((int)Call("ChoosePort")! > 15432, "Occupied port reused");
            }
            finally { listener.Stop(); }
            Console.WriteLine("Setup-token route boundary and occupied-port tests passed.");
            if (args.Length == 0) return 0;
            var backend = Path.GetFullPath(args[0]);
            var root = args.Length > 1 ? Path.GetFullPath(args[1]) : Path.Combine(Path.GetTempPath(), "ARTEX 引导测试 " + Guid.NewGuid().ToString("N"));
            Check(!Directory.Exists(root) || !Directory.EnumerateFileSystemEntries(root).Any(), "Test root must be empty");
            Manager.GetProperty("Root", BindingFlags.Static | BindingFlags.NonPublic)!.SetValue(null, root);
            Directory.CreateDirectory(Path.Combine(root, "data"));
            File.WriteAllText(Path.Combine(root, "data", "preserve.txt"), "existing data");
            try { await Prepare(backend); throw new Exception("Unregistered nonempty directory was initialized"); }
            catch (InvalidOperationException) { }
            Check(File.ReadAllText(Path.Combine(root, "data", "preserve.txt")) == "existing data", "Existing file changed");
            // No recursive cleanup: remove only the explicit fixture file and now-empty directory.
            File.Delete(Path.Combine(root, "data", "preserve.txt")); Directory.Delete(Path.Combine(root, "data"));
            var first = await Prepare(backend);
            var second = await Prepare(backend);
            Check(Field<string>(first, "AppPassword") == Field<string>(second, "AppPassword"), "Retry changed app password");
            var manifest = File.ReadAllBytes(Path.Combine(root, "cluster.dpapi"));
            Check(!System.Text.Encoding.UTF8.GetString(manifest).Contains(Field<string>(first, "AdminPassword")), "Admin secret stored unprotected");
            Check(!File.Exists(Path.Combine(root, "init-password.tmp")), "Password scratch file survived initialization");
            var bin = Field<string>(first, "Bin"); var port = Field<int>(first, "Port");
            await Psql(bin, port, Field<string>(first, "AppPassword"), "CREATE TABLE onboarding_fixture(value text); INSERT INTO onboarding_fixture VALUES ('persisted');");
            var previousContext = SynchronizationContext.Current;
            try
            {
                SynchronizationContext.SetSynchronizationContext(new UnpumpedContext());
                Call("StopOwned");
            }
            finally { SynchronizationContext.SetSynchronizationContext(previousContext); }
            var dsn = (string)first.GetType().GetMethod("Dsn", BindingFlags.Instance | BindingFlags.NonPublic)!.Invoke(first, [false])!;
            await (Task)Call("EnsureForConfigurationAsync", dsn, CancellationToken.None)!;
            var value = await Psql(bin, port, Field<string>(first, "AppPassword"), "SELECT value FROM onboarding_fixture;");
            Check(value.Contains("persisted"), "Data did not survive restart");
            var privileges = await Psql(bin, port, Field<string>(first, "AppPassword"), "SELECT rolsuper OR rolcreatedb OR rolcreaterole FROM pg_roles WHERE rolname=current_user;");
            Check(privileges.Trim() == "f", "Application role has administrative privileges");
            Console.WriteLine("Independent cluster prepare, repeat, DPAPI, ordinary role and restart persistence passed. Fixture preserved: " + root);
            return 0;
        }
        catch (Exception ex) { Console.Error.WriteLine(ex); return 1; }
        finally { Call("StopOwned"); }
    }

    private static async Task<string> Psql(string bin, int port, string password, string sql)
    {
        var info = new ProcessStartInfo(Path.Combine(bin, "psql.exe")) { UseShellExecute = false, CreateNoWindow = true, RedirectStandardOutput = true, RedirectStandardError = true };
        foreach (var a in new[] { "-X", "-A", "-t", "-h", "127.0.0.1", "-p", port.ToString(), "-U", "artex", "-d", "artex", "-v", "ON_ERROR_STOP=1", "-c", sql }) info.ArgumentList.Add(a);
        info.Environment["PGPASSWORD"] = password;
        using var p = Process.Start(info)!;
        var stdout = p.StandardOutput.ReadToEndAsync(); var stderr = p.StandardError.ReadToEndAsync();
        using var timeout = new CancellationTokenSource(TimeSpan.FromSeconds(20));
        await p.WaitForExitAsync(timeout.Token);
        Check(p.ExitCode == 0, "Fixture query failed: " + await stderr);
        return await stdout;
    }
}
