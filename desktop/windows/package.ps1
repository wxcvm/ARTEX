param(
    [string]$Version = "dev",
    [string]$OutputDir = "dist",
    [switch]$SkipWebBuild,
    [switch]$Incremental
)

$ErrorActionPreference = "Stop"
$root = (Resolve-Path (Join-Path $PSScriptRoot "../..")).Path
Set-Location $root

# WSARTEX-inspired fast local builds. CI/release packaging keeps the default
# full-build behavior; SkipWebBuild still takes precedence.
function Test-FrontendFresh {
    $export = Join-Path $root "server/webui/dist"
    $index = Join-Path $export "index.html"
    if (-not (Test-Path $index) -or -not (Test-Path (Join-Path $export "_next"))) { return $false }
    $stamp = (Get-Item $index).LastWriteTimeUtc
    foreach ($relative in @(
        "web/src", "web/public", "web/package.json", "web/package-lock.json",
        "web/next.config.mjs", "web/postcss.config.mjs", "web/tsconfig.json",
        "web/tsconfig.scripts.json", "web/components.json", "web/biome.json",
        "web/.env", "web/.env.local", "web/.env.production"
    )) {
        $path = Join-Path $root $relative
        if (-not (Test-Path $path)) { continue }
        $item = Get-Item $path
        $files = if ($item.PSIsContainer) { Get-ChildItem $path -Recurse -File } else { @($item) }
        foreach ($file in $files) {
            if ($file.LastWriteTimeUtc -gt $stamp) { return $false }
        }
    }
    return $true
}

if (-not $SkipWebBuild -and $Incremental -and (Test-FrontendFresh)) {
    Write-Host "[ARTEX] Frontend unchanged; reusing server/webui/dist"
} elseif (-not $SkipWebBuild) {
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
        & npm run build:static
        if ($LASTEXITCODE -ne 0) { throw "web static export failed" }
    } finally {
        Pop-Location
    }
    $export = Join-Path $root "web/out"
    if (-not (Test-Path (Join-Path $export "index.html"))) { throw "web/out/index.html not found" }
    $embed = Join-Path $root "server/webui/dist"
    if (Test-Path $embed) { Remove-Item -Path $embed -Recurse -Force }
    Copy-Item -Path $export -Destination $embed -Recurse -Force
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
& dotnet publish (Join-Path $root "desktop/windows/ARTEX.Desktop.csproj") -c Release -r win-x64 --self-contained true -p:PublishSingleFile=true -p:IncludeNativeLibrariesForSelfExtract=true -o $stage
if ($LASTEXITCODE -ne 0) { throw "Windows desktop publish failed" }

Copy-Item -Path (Join-Path $root "skills") -Destination (Join-Path $stage "skills") -Recurse -Force
Copy-Item -Path (Join-Path $root "config.example.json") -Destination (Join-Path $stage "config.example.json")
Copy-Item -Path (Join-Path $root "start.bat") -Destination (Join-Path $stage "start.bat")
Copy-Item -Path (Join-Path $root "web/src/app/favicon.ico") -Destination (Join-Path $stage "ARTEX.ico")
Copy-Item -Path (Join-Path $root "desktop/windows/README-Windows.txt") -Destination (Join-Path $stage "README-Windows.txt")
New-Item -ItemType Directory -Path (Join-Path $stage "docs") -Force | Out-Null
Copy-Item -Path (Join-Path $root "docs/USER_GUIDE_ZH.md") -Destination (Join-Path $stage "docs/USER_GUIDE_ZH.md")

$archive = Join-Path $dist "ARTEX-Windows-Desktop-amd64.zip"
if (Test-Path $archive) { Remove-Item -Path $archive -Force }
Compress-Archive -Path (Join-Path $stage "*") -DestinationPath $archive -CompressionLevel Optimal
Write-Host "[ARTEX] Windows desktop bundle: $archive"
