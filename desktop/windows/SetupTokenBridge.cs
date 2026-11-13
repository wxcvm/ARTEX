namespace Artex.Desktop;

internal static class SetupTokenBridge
{
    // A valid browser token may remain after switching to a new database.
    // The database's initialization state, not that cached token, owns this step.
    internal static string InitialRoute(bool initialized) => initialized ? "/" : "/setup/";

    internal static bool IsSetupDocument(string source) =>
        Uri.TryCreate(source, UriKind.Absolute, out var uri) &&
        uri.Scheme == "http" && uri.Host == "127.0.0.1" && uri.Port == 8787 &&
        uri.UserInfo.Length == 0 && uri.Query.Length == 0 && uri.Fragment.Length == 0 &&
        uri.AbsolutePath is "/setup" or "/setup/";
}
