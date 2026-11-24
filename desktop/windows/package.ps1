param(
    [string]$Version = "dev",
    [string]$OutputDir = "dist",
    [switch]$SkipWebBuild,
    [switch]$Incremental
)

$ErrorActionPreference = "Stop"
# Match build.sh: tags may start with v, binary versions do not.
$Version = $Version -replace '^v(?=\d)', ''
# CI uses labels such as preview/dev; assembly versions must remain numeric.
$assemblyVersion = if ($Version -match '^(\d+\.\d+\.\d+)(?:[-+].*)?$') { $Matches[1] } else { '0.0.0' }
$root = (Resolve-Path (Join-Path $PSScriptRoot "../..")).Path
Set-Location $root

# WSARTEX-inspired fast local builds. CI/release packaging keeps the default
# full-build behavior; SkipWebBuild still takes precedence.
. (Join-Path $PSScriptRoot "frontend-cache.ps1")
$frontendCache = Join-Path $root "web/.next/artex-desktop-inputs.json"
$cacheAllowed = $Incremental -and -not $SkipWebBuild -and
    -not $env:ARTEX_BUILD_ID -and -not (Get-ChildItem Env: | Where-Object Name -Like "NEXT_PUBLIC_*")
$inputs = if ($cacheAllowed) { Get-FrontendInputs $root } else { $null }

if ($cacheAllowed -and (Test-FrontendFresh $root $frontendCache $inputs)) {
    Write-Host "[ARTEX] Frontend unchanged; reusing server/webui/dist"
} elseif (-not $SkipWebBuild) {
    Remove-Item -LiteralPath $frontendCache -Force -ErrorAction SilentlyContinue
    Push-Location (Join-Path $root "web")
    try {
        $marker = Join-Path $root "web/node_modules/.package-lock.json"
        $lock = Join-Path $root "web/package-lock.json"
        $manifest = Join-Path $root "web/package.json"
        $depsFresh = $Incremental -and (Test-Path $marker) -and
            ((Get-Item $marker).LastWriteTimeUtc -ge (Get-Item $lock).LastWriteTimeUtc) -and
            ((Get-Item $marker).LastWriteTimeUtc -ge (Get-Item $manifest).LastWriteTimeUtc)
        if (-not $depsFresh) {
            & npm ci
            if ($LASTEXITCODE -ne 0) { throw "npm ci failed" }
        } else {
            Write-Host "[ARTEX] Node dependencies unchanged; skipping npm ci"
        }
        # npm's build:static script uses POSIX inline environment assignment.
        # On native Windows invoke the installed Next CLI with an explicit
        # process environment instead; do not download a second copy of Next.
        $previousExport = [Environment]::GetEnvironmentVariable("NEXT_EXPORT", "Process")
        try {
            $env:NEXT_EXPORT = "1"
            & node (Join-Path $root "web/node_modules/next/dist/bin/next") build
            if ($LASTEXITCODE -ne 0) { throw "web static export failed" }
        } finally {
            [Environment]::SetEnvironmentVariable("NEXT_EXPORT", $previousExport, "Process")
        }
    } finally {
        Pop-Location
    }
    $export = Join-Path $root "web/out"
    if (-not (Test-Path (Join-Path $export "index.html"))) { throw "web/out/index.html not found" }
    $embed = Join-Path $root "server/webui/dist"
    if (Test-Path $embed) { Remove-Item -Path $embed -Recurse -Force }
    Copy-Item -Path $export -Destination $embed -Recurse -Force
    # Do not mark an export reusable if the source changed during its build.
    if ($cacheAllowed -and $inputs -ceq (Get-FrontendInputs $root)) {
        New-Item -ItemType Directory -Path (Split-Path $frontendCache) -Force | Out-Null
        Set-Content -LiteralPath $frontendCache -Value $inputs -NoNewline -Encoding UTF8
    }
}

if (-not (Test-Path (Join-Path $root "server/webui/dist/index.html"))) {
    throw "embedded web UI is missing; build web/out or download the frontend artifact first"
}

$dist = Join-Path $root $OutputDir
$stage = Join-Path $dist "ARTEX-Windows-Desktop-amd64"
if (Test-Path $stage) { Remove-Item -Path $stage -Recurse -Force }
New-Item -ItemType Directory -Path $stage -Force | Out-Null

# The existing Go backend serves the same embedded frontend as browser mode.
$env:CGO_ENABLED = "0"
& go build -tags embedui -ldflags "-s -w -X main.version=$Version" -o (Join-Path $stage "artex.exe") ./cmd/artex
if ($LASTEXITCODE -ne 0) { throw "Go backend build failed" }

# A native Windows window powered by WebView2, not an Edge app shortcut.
& dotnet publish (Join-Path $root "desktop/windows/ARTEX.Desktop.csproj") -c Release -r win-x64 --self-contained true -p:PublishSingleFile=true -p:IncludeNativeLibrariesForSelfExtract=true "-p:Version=$assemblyVersion" "-p:InformationalVersion=$Version" -o $stage
if ($LASTEXITCODE -ne 0) { throw "Windows desktop publish failed" }

Copy-Item -Path (Join-Path $root "skills") -Destination (Join-Path $stage "skills") -Recurse -Force
Copy-Item -Path (Join-Path $root "extension") -Destination (Join-Path $stage "extension") -Recurse -Force
# Only the read-only SQLite ledger is needed at runtime; keep source notes out of the ZIP.
New-Item -ItemType Directory -Path (Join-Path $stage "knowledge") -Force | Out-Null
Copy-Item -Path (Join-Path $root "knowledge/knowledge.db") -Destination (Join-Path $stage "knowledge/knowledge.db")
Copy-Item -Path (Join-Path $root "config.example.json") -Destination (Join-Path $stage "config.example.json")
Copy-Item -Path (Join-Path $root "start.bat") -Destination (Join-Path $stage "start.bat")
Copy-Item -Path (Join-Path $root "desktop/windows/generated/ARTEX.ico") -Destination (Join-Path $stage "ARTEX.ico")
Copy-Item -Path (Join-Path $root "desktop/windows/README-Windows.txt") -Destination (Join-Path $stage "README-Windows.txt")
New-Item -ItemType Directory -Path (Join-Path $stage "docs") -Force | Out-Null
Copy-Item -Path (Join-Path $root "docs/USER_GUIDE_ZH.md") -Destination (Join-Path $stage "docs/USER_GUIDE_ZH.md")

$archive = Join-Path $dist "ARTEX-Windows-Desktop-amd64.zip"
if (Test-Path $archive) { Remove-Item -Path $archive -Force }
Compress-Archive -Path (Join-Path $stage "*") -DestinationPath $archive -CompressionLevel Optimal
Write-Host "[ARTEX] Windows desktop bundle: $archive"
