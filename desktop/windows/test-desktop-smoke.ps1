param(
    [string]$OutputDir = "dist",
    [ValidateRange(1, 50)] [int]$WindowCycles = 3
)

$ErrorActionPreference = "Stop"
$root = (Resolve-Path (Join-Path $PSScriptRoot "../..")).Path
$program = Join-Path $root (Join-Path $OutputDir "ARTEX-Windows-Desktop-amd64/ARTEX.Desktop.exe")
if (-not (Test-Path $program)) { throw "Desktop executable not found: $program" }
$extensionManifest = Join-Path $root (Join-Path $OutputDir "ARTEX-Windows-Desktop-amd64/extension/manifest.json")
if (-not (Test-Path $extensionManifest)) { throw "Bundled browser capture extension is missing: $extensionManifest" }

$tempDir = if ($env:RUNNER_TEMP) { $env:RUNNER_TEMP } else { [IO.Path]::GetTempPath() }
$requestLog = Join-Path $tempDir "artex-desktop-http-requests.txt"
Remove-Item -LiteralPath $requestLog -ErrorAction SilentlyContinue
$desktop = $null
$server = $null

try {
    # A tiny local HTTP fixture exercises the existing-health path without a
    # database, and verifies that WebView2 actually navigates to the UI.
    $server = Start-ThreadJob -ArgumentList $requestLog -ScriptBlock {
        param($file)
        $ErrorActionPreference = "Stop"
        $listener = [System.Net.Sockets.TcpListener]::new([System.Net.IPAddress]::Loopback, 8787)
        $listener.Start()
        try {
            $running = $true
            while ($running) {
                if (-not $listener.Pending()) { Start-Sleep -Milliseconds 50; continue }
                $client = $listener.AcceptTcpClient()
                try {
                    $stream = $client.GetStream()
                    $stream.ReadTimeout = 1000
                    $reader = [IO.StreamReader]::new($stream, [Text.Encoding]::UTF8, $false, 1024, $true)
                    $request = $reader.ReadLine()
                    if (-not $request) { continue }
                    $route = ($request -split ' ')[1]
                    while ($true) {
                        $header = $reader.ReadLine()
                        if ($null -eq $header -or $header.Length -eq 0) { break }
                    }
                    Add-Content -LiteralPath $file -Value $route
                    if ($route -eq "/__shutdown") { $running = $false }
                    if ($route -eq "/api/health") {
                        $body = '{"ok":true,"service":"artex","version":"smoke"}'
                        $type = "application/json"
                    } else {
                        $body = '<!doctype html><html><head><title>ARTEX desktop smoke</title></head><body>ARTEX</body></html>'
                        $type = "text/html"
                    }
                    $bytes = [Text.Encoding]::UTF8.GetBytes($body)
                    $head = "HTTP/1.1 200 OK" + [Environment]::NewLine +
                        "Content-Type: $type; charset=utf-8" + [Environment]::NewLine +
                        "Content-Length: $($bytes.Length)" + [Environment]::NewLine +
                        "Connection: close" + [Environment]::NewLine + [Environment]::NewLine
                    $headerBytes = [Text.Encoding]::UTF8.GetBytes($head)
                    $stream.Write($headerBytes, 0, $headerBytes.Length)
                    $stream.Write($bytes, 0, $bytes.Length)
                } catch {
                    # Chromium may open an idle speculative connection or close
                    # a socket before sending a request. A single client must
                    # not terminate the fixture (ReadLine has a 1s timeout).
                    $cause = $_.Exception
                    while ($cause.InnerException) { $cause = $cause.InnerException }
                    if ($cause -is [System.IO.IOException] -or $cause -is [System.Net.Sockets.SocketException]) { continue }
                    throw
                } finally {
                    $client.Dispose()
                }
            }
        } finally {
            $listener.Stop()
        }
    }

    $serverReady = $false
    for ($i = 0; $i -lt 40; $i++) {
        if ($server.State -eq "Failed") {
            throw "Local HTTP fixture failed: $($server.JobStateInfo.Reason)"
        }
        try {
            $health = Invoke-RestMethod "http://127.0.0.1:8787/api/health" -TimeoutSec 1
            if ($health.service -eq "artex") { $serverReady = $true; break }
        } catch {}
        Start-Sleep -Milliseconds 200
    }
    if (-not $serverReady) { throw "Local HTTP fixture did not start" }

    # Regression: a speculative idle connection must not kill the fixture.
    # Verify a real health request still succeeds after its read timeout.
    $idleClient = [System.Net.Sockets.TcpClient]::new()
    try {
        $idleClient.Connect("127.0.0.1", 8787)
        Start-Sleep -Milliseconds 1500
    } finally { $idleClient.Dispose() }
    if ($server.State -ne "Running") { throw "Local HTTP fixture died on an idle connection: $($server.JobStateInfo.Reason)" }
    $health = Invoke-RestMethod "http://127.0.0.1:8787/api/health" -TimeoutSec 5
    if ($health.service -ne "artex") { throw "Local HTTP fixture health failed after idle connection" }
    Write-Host "[ARTEX] HTTP fixture survives idle browser connections."

    $desktop = Start-Process -FilePath $program -WorkingDirectory (Split-Path $program) -PassThru
    $navigated = $false
    for ($i = 0; $i -lt 240; $i++) {
        if ($server.State -ne "Running") {
            throw "Local HTTP fixture stopped before WebView2 navigation: $($server.State); $($server.JobStateInfo.Reason)"
        }
        $desktop.Refresh()
        if ($desktop.HasExited) {
            throw "ARTEX desktop exited before navigating (code $($desktop.ExitCode))"
        }
        if ((Test-Path $requestLog) -and ((Get-Content -LiteralPath $requestLog) -contains "/")) {
            $navigated = $true
            break
        }
        Start-Sleep -Milliseconds 250
    }
    if (-not $navigated) {
        $seen = if (Test-Path $requestLog) { (Get-Content -LiteralPath $requestLog) -join ", " } else { "(no requests)" }
        throw "ARTEX desktop did not load the local UI in WebView2 within 60 seconds. HTTP routes seen: $seen"
    }
    Write-Host "[ARTEX] Native window, health probe and WebView2 navigation passed."

    # Exercise the actual packaged WinForms window, not just the HTTP backend.
    # A minimized window must retain WS_VISIBLE and be restorable without the tray.
    Add-Type -TypeDefinition @'
using System;
using System.Runtime.InteropServices;
public static class ArtexSmokeWindow {
    [DllImport("user32.dll")] public static extern bool ShowWindowAsync(IntPtr hwnd, int command);
    [DllImport("user32.dll")] public static extern bool IsWindowVisible(IntPtr hwnd);
    [DllImport("user32.dll")] public static extern bool IsIconic(IntPtr hwnd);
    [DllImport("user32.dll")] public static extern bool IsWindow(IntPtr hwnd);
    [DllImport("user32.dll")] public static extern IntPtr GetForegroundWindow();
    [StructLayout(LayoutKind.Sequential)]
    private struct POINT { public int X, Y; }
    [StructLayout(LayoutKind.Sequential)]
    private struct RECT { public int Left, Top, Right, Bottom; }
    [StructLayout(LayoutKind.Sequential)]
    private struct WINDOWPLACEMENT {
        public uint length, flags, showCmd;
        public POINT ptMinPosition, ptMaxPosition;
        public RECT rcNormalPosition;
    }
    [DllImport("user32.dll", SetLastError = true)]
    private static extern bool GetWindowPlacement(IntPtr hwnd, ref WINDOWPLACEMENT placement);
    public static string Describe(IntPtr hwnd) {
        var placement = new WINDOWPLACEMENT { length = (uint)Marshal.SizeOf<WINDOWPLACEMENT>() };
        var success = GetWindowPlacement(hwnd, ref placement);
        var error = success ? 0 : Marshal.GetLastWin32Error();
        return $"hwnd={hwnd}, valid={IsWindow(hwnd)}, visible={IsWindowVisible(hwnd)}, iconic={IsIconic(hwnd)}, foreground={GetForegroundWindow()}, placementOk={success}, showCmd={placement.showCmd}, flags={placement.flags}, placementError={error}";
    }
    [DllImport("user32.dll")] public static extern bool PostMessage(IntPtr hwnd, uint msg, IntPtr wParam, IntPtr lParam);
}
'@
    $hwnd = [IntPtr]::Zero
    for ($i = 0; $i -lt 60; $i++) {
        $desktop.Refresh()
        if ($desktop.HasExited) { throw "Desktop exited before window lifecycle test" }
        $hwnd = $desktop.MainWindowHandle
        if ($hwnd -ne [IntPtr]::Zero -and [ArtexSmokeWindow]::IsWindowVisible($hwnd)) { break }
        Start-Sleep -Milliseconds 100
    }
    if ($hwnd -eq [IntPtr]::Zero -or -not [ArtexSmokeWindow]::IsWindowVisible($hwnd)) {
        throw "Desktop main window not visible: $([ArtexSmokeWindow]::Describe($hwnd))"
    }

    function Assert-ArtexWindowState([bool]$minimized, [string]$action) {
        for ($i = 0; $i -lt 60; $i++) {
            $desktop.Refresh()
            if ($desktop.HasExited) { throw "Desktop exited after $action (code $($desktop.ExitCode))" }
            if ([ArtexSmokeWindow]::IsWindow($hwnd) -and [ArtexSmokeWindow]::IsWindowVisible($hwnd) -and
                ([ArtexSmokeWindow]::IsIconic($hwnd) -eq $minimized)) { return }
            Start-Sleep -Milliseconds 100
        }
        throw "Desktop window did not retain its visible taskbar state after $action (expected minimized=$minimized): $([ArtexSmokeWindow]::Describe($hwnd)); process exited=$($desktop.HasExited); current main hwnd=$($desktop.MainWindowHandle)"
    }

    for ($cycle = 1; $cycle -le $WindowCycles; $cycle++) {
        [ArtexSmokeWindow]::ShowWindowAsync($hwnd, 6) | Out-Null # SW_MINIMIZE
        Assert-ArtexWindowState $true "minimize (cycle $cycle)"
        [ArtexSmokeWindow]::ShowWindowAsync($hwnd, 9) | Out-Null # SW_RESTORE
        Assert-ArtexWindowState $false "restore (cycle $cycle)"
        # The title-bar X sends WM_SYSCOMMAND/SC_CLOSE. Posting WM_CLOSE
        # directly is interpreted by WinForms as TaskManagerClosing, not
        # UserClosing, and intentionally takes the actual shutdown path.
        if (-not [ArtexSmokeWindow]::PostMessage($hwnd, 0x0112, [IntPtr]0xF060, [IntPtr]::Zero)) {
            throw "Could not send SC_CLOSE to desktop (cycle $cycle): $([ArtexSmokeWindow]::Describe($hwnd))"
        }
        Assert-ArtexWindowState $true "close button (cycle $cycle)"
        [ArtexSmokeWindow]::ShowWindowAsync($hwnd, 9) | Out-Null
        Assert-ArtexWindowState $false "restore after close (cycle $cycle)"
        Write-Host "[ARTEX] Minimize, close-to-taskbar and restore passed (cycle $cycle/$WindowCycles)."
    }
}
finally {
    if ($desktop) {
        $desktop.Refresh()
        if (-not $desktop.HasExited) {
            $desktop.Kill($true)
            $desktop.WaitForExit(5000) | Out-Null
        }
        $desktop.Dispose()
    }
    if ($server) {
        if ($serverReady) {
            try { Invoke-WebRequest -Uri "http://127.0.0.1:8787/__shutdown" -TimeoutSec 2 | Out-Null } catch {}
        }
        $finished = Wait-Job -Job $server -Timeout 5
        if (-not $finished) { Stop-Job $server -ErrorAction SilentlyContinue }
        Remove-Job $server -Force -ErrorAction SilentlyContinue
    }
    Remove-Item -LiteralPath $requestLog -ErrorAction SilentlyContinue
}
