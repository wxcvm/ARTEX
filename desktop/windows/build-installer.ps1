param(
    [string]$Version = "0.0.0-preview",
    [string]$OutputDir = "dist"
)

$ErrorActionPreference = "Stop"
$Version = $Version -replace '^v(?=\d)', ''
$root = (Resolve-Path (Join-Path $PSScriptRoot "../..")).Path
$stage = Join-Path $root (Join-Path $OutputDir "ARTEX-Windows-Desktop-amd64")
$dist = Join-Path $root $OutputDir
# The installer consumes the whole stage. Reject local runtime state rather
# than accidentally distributing it from a reused or manually modified stage.
foreach ($relative in @("config.json", "data", "keys", ".env")) {
    if (Test-Path -LiteralPath (Join-Path $stage $relative)) {
        throw "Desktop stage contains local runtime state: $relative. Rebuild the portable package in a clean stage."
    }
}
foreach ($relative in @(
    "ARTEX.Desktop.exe", "artex.exe", "ARTEX.ico", "README-Windows.txt",
    "knowledge/knowledge.db", "docs/USER_GUIDE_ZH.md", "extension/manifest.json",
    "config.example.json", "start.bat"
)) {
    if (-not (Test-Path -LiteralPath (Join-Path $stage $relative) -PathType Leaf)) {
        throw "Desktop stage is incomplete: $relative. Run desktop/windows/package.ps1 first."
    }
}
if (-not (Test-Path -LiteralPath (Join-Path $stage "skills") -PathType Container)) {
    throw "Desktop stage is incomplete: skills. Run desktop/windows/package.ps1 first."
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
