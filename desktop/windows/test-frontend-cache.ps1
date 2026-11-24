$ErrorActionPreference = "Stop"
. (Join-Path $PSScriptRoot "frontend-cache.ps1")
$root = Join-Path ([IO.Path]::GetTempPath()) ("artex-frontend-cache-" + [Guid]::NewGuid().ToString("N"))
try {
    foreach ($dir in @("web/src", "web/public", "server/webui/dist/_next", "web/.next")) {
        New-Item -ItemType Directory -Path (Join-Path $root $dir) -Force | Out-Null
    }
    $source = Join-Path $root "web/src/page.tsx"
    $asset = Join-Path $root "web/public/old.txt"
    $cache = Join-Path $root "web/.next/artex-desktop-inputs.json"
    Set-Content -LiteralPath $source -Value "page"
    Set-Content -LiteralPath $asset -Value "asset"
    Set-Content -LiteralPath (Join-Path $root "server/webui/dist/index.html") -Value "export"
    function Assert-Fresh([bool]$expected, [string]$label) {
        $actual = Test-FrontendFresh $root $cache (Get-FrontendInputs $root)
        if ($actual -ne $expected) { throw "Unexpected cache state for $label (expected $expected)" }
    }
    function Save-Inputs {
        Set-Content -LiteralPath $cache -Value (Get-FrontendInputs $root) -NoNewline -Encoding UTF8
    }
    Assert-Fresh $false "missing inventory"
    Save-Inputs
    Assert-Fresh $true "unchanged source"
    Rename-Item -LiteralPath $asset -NewName "renamed.txt"
    Assert-Fresh $false "rename"
    Save-Inputs
    Remove-Item -LiteralPath (Join-Path $root "web/public/renamed.txt")
    Assert-Fresh $false "deletion"
    Save-Inputs
    Set-Content -LiteralPath $source -Value "changed content"
    Assert-Fresh $false "edit"
    Save-Inputs
    $envFile = Join-Path $root "web/.env.production.local"
    Set-Content -LiteralPath $envFile -Value "NEXT_PUBLIC_EXAMPLE=synthetic"
    Assert-Fresh $false "added production environment file"
    Save-Inputs
    Remove-Item -LiteralPath $envFile
    Assert-Fresh $false "removed production environment file"
    Save-Inputs
    Remove-Item -LiteralPath (Join-Path $root "server/webui/dist/index.html")
    Assert-Fresh $false "missing export"
    Write-Host "[ARTEX] Incremental frontend cache regression passed."
} finally {
    if (Test-Path -LiteralPath $root) { Remove-Item -LiteralPath $root -Recurse -Force }
}
