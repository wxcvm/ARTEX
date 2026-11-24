# Metadata inventory only: never store environment file contents or credentials.
function Get-FrontendInputs([string]$Root) {
    $files = foreach ($relative in @(
        "web/src", "web/public", "web/package.json", "web/package-lock.json",
        "web/next.config.mjs", "web/postcss.config.mjs", "web/tsconfig.json",
        "web/tsconfig.scripts.json", "web/components.json", "web/biome.json",
        "web/.env", "web/.env.local", "web/.env.production", "web/.env.production.local"
    )) {
        $path = Join-Path $Root $relative
        if (-not (Test-Path -LiteralPath $path)) { continue }
        $item = Get-Item -LiteralPath $path -Force
        if ($item.PSIsContainer) { Get-ChildItem -LiteralPath $path -Recurse -File -Force }
        else { $item }
    }
    $entries = @($files | Sort-Object FullName | ForEach-Object {
        [ordered]@{
            path = $_.FullName.Substring($Root.TrimEnd([IO.Path]::DirectorySeparatorChar, [IO.Path]::AltDirectorySeparatorChar).Length + 1)
            size = $_.Length
            ticks = $_.LastWriteTimeUtc.Ticks
        }
    })
    ConvertTo-Json -InputObject $entries -Compress
}

function Test-FrontendFresh([string]$Root, [string]$Cache, [string]$Inputs) {
    $export = Join-Path $Root "server/webui/dist"
    if (-not (Test-Path -LiteralPath (Join-Path $export "index.html") -PathType Leaf) -or
        -not (Test-Path -LiteralPath (Join-Path $export "_next") -PathType Container) -or
        -not (Test-Path -LiteralPath $Cache -PathType Leaf)) { return $false }
    return (Get-Content -LiteralPath $Cache -Raw) -ceq $Inputs
}
