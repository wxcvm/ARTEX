param(
    [string]$OutputDir = "dist"
)

$ErrorActionPreference = "Stop"
$root = (Resolve-Path (Join-Path $PSScriptRoot "../..")).Path
$setup = Join-Path $root (Join-Path $OutputDir "ARTEX-Windows-Setup-amd64.exe")
$installed = Join-Path $env:RUNNER_TEMP "ARTEX Installer Smoke"
$shortcut = Join-Path $env:APPDATA "Microsoft/Windows/Start Menu/Programs/ARTEX/ARTEX.lnk"

function Run-Setup([string]$path, [string]$arguments) {
    Write-Host "[ARTEX] Running setup: $path"
    $p = Start-Process -FilePath $path -ArgumentList $arguments -Wait -PassThru
    if ($p.ExitCode -ne 0) { throw "Setup exited $($p.ExitCode): $path" }
    Write-Host "[ARTEX] Setup finished successfully."
}

if (-not (Test-Path $setup)) { throw "Missing installer: $setup" }
if (Test-Path $installed) { Remove-Item -Recurse -Force $installed }

$installArgs = "/VERYSILENT /SUPPRESSMSGBOXES /NORESTART /SP- /DIR=" + '"' + $installed + '"'
Run-Setup $setup $installArgs

foreach ($file in @("ARTEX.Desktop.exe", "artex.exe", "ARTEX.ico", "README-Windows.txt", "knowledge/knowledge.db", "docs/USER_GUIDE_ZH.md", "unins000.exe")) {
    if (-not (Test-Path (Join-Path $installed $file))) {
        throw "Installed package is incomplete: $file"
    }
}
if (-not (Test-Path $shortcut)) { throw "Start menu shortcut missing: $shortcut" }

# A subsequent installation must preserve user configuration and edited skills.
$config = Join-Path $installed "config.json"
$data = Join-Path $installed "data"
New-Item -ItemType Directory -Path $data -Force | Out-Null
$dataFixture = Join-Path $data "installer-test.txt"
Set-Content -Path $dataFixture -Value "preserve user data" -NoNewline
Set-Content -Path $config -Value "installer-test-user-settings" -NoNewline
$skill = Get-ChildItem (Join-Path $installed "skills") -Recurse -File | Select-Object -First 1
if (-not $skill) { throw "No skills shipped in the installed package" }
Add-Content -Path $skill.FullName -Value "# installer-test-user-skill"
Run-Setup $setup $installArgs
if ((Get-Content -Raw $config) -ne "installer-test-user-settings") {
    throw "An upgrade overwrote the user's config.json"
}
if ((Get-Content -Raw $dataFixture) -ne "preserve user data") {
    throw "An upgrade overwrote user data"
}
if (-not (Select-String -Path $skill.FullName -Pattern "# installer-test-user-skill" -SimpleMatch -Quiet)) {
    throw "An upgrade overwrote an edited skill"
}

Run-Setup (Join-Path $installed "unins000.exe") "/VERYSILENT /SUPPRESSMSGBOXES /NORESTART"
if (Test-Path (Join-Path $installed "ARTEX.Desktop.exe")) {
    throw "Uninstall left the desktop executable behind"
}
if (-not (Test-Path $config) -or -not (Test-Path $skill.FullName) -or -not (Test-Path $dataFixture)) {
    throw "Uninstall removed user configuration or modified skills"
}
if (Test-Path $shortcut) { throw "Uninstall left the Start menu shortcut behind" }

# Runner cleanup only; the product's uninstaller intentionally keeps user data.
Remove-Item -Path $installed -Recurse -Force
Write-Host "[ARTEX] Installer install/upgrade/uninstall smoke passed."
