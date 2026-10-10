$ErrorActionPreference = 'Stop'
Set-Location $PSScriptRoot

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    throw 'Go not found. Install Go 1.23+ first.'
}
if (-not (Get-Command gcc -ErrorAction SilentlyContinue)) {
    throw @'
CGO needs a C compiler, but gcc was not found.
Recommended: install MSYS2, open UCRT64 shell once and install mingw-w64-ucrt-x86_64-gcc,
then add C:\msys64\ucrt64\bin to PATH and rerun this script.
'@
}

$env:CGO_ENABLED = '1'
New-Item -ItemType Directory -Force -Path dist | Out-Null
$dll = Join-Path $PSScriptRoot 'dist\cpa-scheduled-tests.dll'
& go build -trimpath -buildmode=c-shared -ldflags='-s -w' -o $dll .
if ($LASTEXITCODE -ne 0) { throw 'go build failed' }
Remove-Item 'dist\cpa-scheduled-tests.h' -ErrorAction SilentlyContinue

$zip = Join-Path $PSScriptRoot 'dist\cpa-scheduled-tests_0.1.6_windows_amd64.zip'
Remove-Item $zip -ErrorAction SilentlyContinue
Compress-Archive -Path $dll -DestinationPath $zip
$checksum = (Get-FileHash -LiteralPath $zip -Algorithm SHA256).Hash.ToLowerInvariant()
[IO.File]::WriteAllText((Join-Path $PSScriptRoot 'dist\checksums.txt'), "$checksum  $([IO.Path]::GetFileName($zip))`n", [Text.UTF8Encoding]::new($false))
Write-Host "Built DLL: $dll"
Write-Host "Built ZIP: $zip"
