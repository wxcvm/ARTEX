# ARTEX desktop icon is derived from web/public/logo.jpg.
# Keep a single brand image for the browser, app executable, taskbar and installer.
$ErrorActionPreference = "Stop"
Add-Type -AssemblyName System.Drawing

$root = (Resolve-Path (Join-Path $PSScriptRoot "../..")).Path
$sourcePath = Join-Path $root "web/public/logo.jpg"
$outDir = Join-Path $PSScriptRoot "generated"
$outFile = Join-Path $outDir "ARTEX.ico"

if (-not (Test-Path $sourcePath)) {
    throw "ARTEX brand image is missing: $sourcePath"
}

New-Item -ItemType Directory -Force -Path $outDir | Out-Null

$source = [System.Drawing.Image]::FromFile($sourcePath)
$entries = [System.Collections.Generic.List[object]]::new()
try {
    foreach ($size in @(16, 24, 32, 48, 64, 128, 256)) {
        $bitmap = [System.Drawing.Bitmap]::new(
            $size, $size, [System.Drawing.Imaging.PixelFormat]::Format32bppArgb)
        try {
            $graphics = [System.Drawing.Graphics]::FromImage($bitmap)
            try {
                $graphics.Clear([System.Drawing.Color]::Transparent)
                $graphics.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
                $graphics.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::HighQuality
                $graphics.PixelOffsetMode = [System.Drawing.Drawing2D.PixelOffsetMode]::HighQuality
                $graphics.DrawImage($source, 0, 0, $size, $size)
            }
            finally { $graphics.Dispose() }

            $png = [System.IO.MemoryStream]::new()
            try {
                $bitmap.Save($png, [System.Drawing.Imaging.ImageFormat]::Png)
                $entries.Add([pscustomobject]@{
                    Size = $size
                    Data = $png.ToArray()
                })
            }
            finally { $png.Dispose() }
        }
        finally { $bitmap.Dispose() }
    }
}
finally { $source.Dispose() }

$stream = [System.IO.File]::Create($outFile)
$writer = [System.IO.BinaryWriter]::new($stream)
try {
    # ICO directory: reserved, image type 1, image count.
    $writer.Write([uint16]0)
    $writer.Write([uint16]1)
    $writer.Write([uint16]$entries.Count)
    $offset = 6 + 16 * $entries.Count
    foreach ($entry in $entries) {
        $dimension = if ($entry.Size -eq 256) { 0 } else { $entry.Size }
        $writer.Write([byte]$dimension)
        $writer.Write([byte]$dimension)
        $writer.Write([byte]0)
        $writer.Write([byte]0)
        $writer.Write([uint16]1)
        $writer.Write([uint16]32)
        $writer.Write([uint32]$entry.Data.Length)
        $writer.Write([uint32]$offset)
        $offset += $entry.Data.Length
    }
    foreach ($entry in $entries) {
        $writer.Write([byte[]]$entry.Data)
    }
}
finally { $writer.Dispose() }

Write-Host "[ARTEX] Generated eagle-shield icon: $outFile"
