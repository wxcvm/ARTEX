param([string]$Backend = "dist/ARTEX-Windows-Desktop-amd64/artex.exe")
$ErrorActionPreference = "Stop"
$root = (Resolve-Path (Join-Path $PSScriptRoot "../..")).Path
$backendPath = (Resolve-Path (Join-Path $root $Backend)).Path
$project = Join-Path $root "tests/windows-onboarding/WindowsOnboarding.csproj"
$principal = [Security.Principal.WindowsPrincipal]::new([Security.Principal.WindowsIdentity]::GetCurrent())
if (-not $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    & dotnet run --project $project --configuration Release -- $backendPath
    if ($LASTEXITCODE -ne 0) { throw "Onboarding integration failed" }
    exit
}
# Hosted Windows runners are elevated. Exercise initdb as a real ordinary user,
# not by weakening the application's normal-user requirement.
$stage = Join-Path $env:PUBLIC ("ARTEX-onboarding-ci-" + [Guid]::NewGuid().ToString("N"))
$name = "artex" + [Guid]::NewGuid().ToString("N").Substring(0, 10)
$password = ConvertTo-SecureString ([Guid]::NewGuid().ToString("N") + "!aA7") -AsPlainText -Force
$userCreated = $false
try {
    New-Item -ItemType Directory -Path $stage | Out-Null
    New-LocalUser -Name $name -Password $password -PasswordNeverExpires | Out-Null
    $userCreated = $true
    Add-LocalGroupMember -Group "Users" -Member $name
    $sid = (Get-LocalUser -Name $name).SID
    $acl = Get-Acl -LiteralPath $stage
    $acl.AddAccessRule([Security.AccessControl.FileSystemAccessRule]::new($sid, "FullControl", "ContainerInherit,ObjectInherit", "None", "Allow"))
    Set-Acl -LiteralPath $stage -AclObject $acl
    & dotnet publish $project -c Release -r win-x64 --self-contained true -o (Join-Path $stage "runner")
    if ($LASTEXITCODE -ne 0) { throw "Cannot publish onboarding test runner" }
    Copy-Item -LiteralPath $backendPath -Destination (Join-Path $stage "artex.exe")
    $credential = [PSCredential]::new("$env:COMPUTERNAME\$name", $password)
    $run = Start-Process -FilePath (Join-Path $stage "runner/WindowsOnboarding.exe") -WorkingDirectory $stage `
        -ArgumentList @('"' + (Join-Path $stage 'artex.exe') + '"', '"' + (Join-Path $stage '独立 数据库') + '"') `
        -Credential $credential -LoadUserProfile -WindowStyle Hidden -Wait -PassThru `
        -RedirectStandardOutput (Join-Path $stage "result.log") -RedirectStandardError (Join-Path $stage "error.log")
    Get-Content -LiteralPath (Join-Path $stage "result.log")
    Get-Content -LiteralPath (Join-Path $stage "error.log")
    if ($run.ExitCode -ne 0) { throw "Ordinary-user onboarding integration failed ($($run.ExitCode))" }
} finally {
    if ($userCreated) { Remove-LocalUser -Name $name }
    # Preserve fixture and logs for failure diagnosis. No recursive removal.
}
