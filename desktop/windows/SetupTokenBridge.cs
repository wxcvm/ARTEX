namespace Artex.Desktop;

internal static class SetupTokenBridge
{
    internal static bool IsSetupDocument(string source) =>
        Uri.TryCreate(source, UriKind.Absolute, out var uri) &&
        uri.Scheme == "http" && uri.Host == "127.0.0.1" && uri.Port == 8787 &&
        uri.UserInfo.Length == 0 && uri.Query.Length == 0 && uri.Fragment.Length == 0 &&
        uri.AbsolutePath is "/setup" or "/setup/";
}
