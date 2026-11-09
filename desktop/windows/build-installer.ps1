param(
    [string]$Version = "0.0.0-preview",
    [string]$OutputDir = "dist"
)

$ErrorActionPreference = "Stop"
$root = (Resolve-Path (Join-Path $PSScriptRoot "../..")).Path
$stage = Join-Path $root (Join-Path $OutputDir "ARTEX-Windows-Desktop-amd64")
$dist = Join-Path $root $OutputDir
if (-not (Test-Path (Join-Path $stage "ARTEX.Desktop.exe")) -or
    -not (Test-Path (Join-Path $stage "artex.exe")) -or
    -not (Test-Path (Join-Path $stage "ARTEX.ico"))) {
    throw "Desktop stage is incomplete. Run desktop/windows/package.ps1 first."
}
New-Item -ItemType Directory -Path $dist -Force | Out-Null

$compiler = Get-Command ISCC.exe -ErrorAction SilentlyContinue | Select-Object -First 1 -ExpandProperty Source
if (-not $compiler) {
    foreach ($path in @(
        (Join-Path ([Environment]::GetFolderPath("ProgramFilesX86")) "Inno Setup 6/ISCC.exe"),
        (Join-Path $env:ProgramFiles "Inno Setup 6/ISCC.exe"),
        "C:\ProgramData\chocolatey\bin\ISCC.exe"
    )) {
        if (Test-Path $path) { $compiler = $path; break }
    }
}
if (-not $compiler) {
    throw "Inno Setup 6 (ISCC.exe) not found. Install from https://jrsoftware.org/isdl.php and retry."
}
& $compiler "/DStageDir=$stage" "/DOutputDir=$dist" "/DAppVersion=$Version" (Join-Path $PSScriptRoot "installer.iss")
if ($LASTEXITCODE -ne 0) { throw "Inno Setup compiler failed (exit $LASTEXITCODE)" }

$installer = Join-Path $dist "ARTEX-Windows-Setup-amd64.exe"
if (-not (Test-Path $installer)) { throw "Installer was not produced: $installer" }
Write-Host "[ARTEX] Windows installer: $installer"
