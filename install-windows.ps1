param(
    [Parameter(Mandatory = $true)]
    [string]$CpaDir,
    [string]$BaseUrl = "http://127.0.0.1:8317",
    [string]$ManagementKey = "",
    [switch]$SkipConfig
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$ProjectDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ToolDir = Join-Path $ProjectDir ".toolchain"
$DistDir = Join-Path $ProjectDir "dist"
$GoVersion = "1.23.2"
$ZigVersion = "0.16.0"
$GoZip = Join-Path $ToolDir "go$GoVersion.windows-amd64.zip"
$ZigZip = Join-Path $ToolDir "zig-x86_64-windows-$ZigVersion.zip"
$GoRoot = Join-Path $ToolDir "go"
$GoExe = Join-Path $GoRoot "bin\go.exe"
$ZigRoot = Join-Path $ToolDir "zig-x86_64-windows-$ZigVersion"
$ZigExe = Join-Path $ZigRoot "zig.exe"
$CCWrapper = Join-Path $ToolDir "zigcc.cmd"
$Dll = Join-Path $DistDir "cpa-scheduled-tests.dll"
$Package = Join-Path $DistDir "cpa-scheduled-tests_0.1.6_windows_amd64.zip"

function Download-IfMissing {
    param([string]$Url, [string]$Path)
    if (Test-Path $Path) { return }
    Write-Host "Downloading $Url"
    Invoke-WebRequest -Uri $Url -OutFile $Path
}

if (-not (Test-Path $CpaDir)) {
    throw "CPA directory does not exist: $CpaDir"
}

New-Item -ItemType Directory -Force -Path $ToolDir, $DistDir | Out-Null

if (-not (Test-Path $GoExe)) {
    Download-IfMissing "https://go.dev/dl/go$GoVersion.windows-amd64.zip" $GoZip
    if (Test-Path $GoRoot) { Remove-Item -Recurse -Force $GoRoot }
    Expand-Archive -Path $GoZip -DestinationPath $ToolDir -Force
}

if (-not (Test-Path $ZigExe)) {
    Download-IfMissing "https://ziglang.org/download/$ZigVersion/zig-x86_64-windows-$ZigVersion.zip" $ZigZip
    if (Test-Path $ZigRoot) { Remove-Item -Recurse -Force $ZigRoot }
    Expand-Archive -Path $ZigZip -DestinationPath $ToolDir -Force
}

$ccLine = '@echo off' + [Environment]::NewLine + '"' + $ZigExe + '" cc %*' + [Environment]::NewLine
Set-Content -LiteralPath $CCWrapper -Value $ccLine -Encoding Ascii

$oldGOOS = $env:GOOS
$oldGOARCH = $env:GOARCH
$oldCGO = $env:CGO_ENABLED
$oldCC = $env:CC
try {
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    $env:CGO_ENABLED = "1"
    $env:CC = $CCWrapper

    Push-Location $ProjectDir
    try {
        Write-Host "Running tests..."
        & $GoExe test ./...
        if ($LASTEXITCODE -ne 0) { throw "go test failed" }

        Write-Host "Building Windows DLL..."
        & $GoExe build -trimpath -buildmode=c-shared -ldflags="-s -w" -o $Dll .
        if ($LASTEXITCODE -ne 0) { throw "go build failed" }
    }
    finally {
        Pop-Location
    }
}
finally {
    $env:GOOS = $oldGOOS
    $env:GOARCH = $oldGOARCH
    $env:CGO_ENABLED = $oldCGO
    $env:CC = $oldCC
}

$Header = [System.IO.Path]::ChangeExtension($Dll, ".h")
if (Test-Path $Header) { Remove-Item -Force $Header }
if (Test-Path $Package) { Remove-Item -Force $Package }
Compress-Archive -Path $Dll -DestinationPath $Package -Force
$checksum = (Get-FileHash -LiteralPath $Package -Algorithm SHA256).Hash.ToLowerInvariant()
[IO.File]::WriteAllText((Join-Path $DistDir 'checksums.txt'), "$checksum  $([IO.Path]::GetFileName($Package))`n", [Text.UTF8Encoding]::new($false))

$PluginDir = Join-Path $CpaDir "plugins\windows\amd64"
New-Item -ItemType Directory -Force -Path $PluginDir | Out-Null
$InstalledDll = Join-Path $PluginDir "cpa-scheduled-tests.dll"
Copy-Item -Force $Dll $InstalledDll

$configApplied = $false
if (-not $SkipConfig -and -not [string]::IsNullOrWhiteSpace($ManagementKey)) {
    try {
        $headers = @{ Authorization = "Bearer $ManagementKey" }
        $body = @{ enabled = $true; priority = 1 } | ConvertTo-Json -Compress
        Invoke-RestMethod -Method Put -Uri "$BaseUrl/v0/management/plugins/cpa-scheduled-tests/config" -Headers $headers -ContentType "application/json" -Body $body | Out-Null
        $configApplied = $true
    }
    catch {
        Write-Warning "DLL installed, but CPA plugin config API was not applied: $($_.Exception.Message)"
    }
}

Write-Host ""
Write-Host "Installed DLL: $InstalledDll"
Write-Host "Release ZIP:  $Package"
if ($configApplied) {
    Write-Host "CPA config enabled through Management API."
} else {
    Write-Host "Ensure config.yaml contains:"
    Write-Host "plugins:"
    Write-Host "  enabled: true"
    Write-Host "  dir: plugins"
    Write-Host "  configs:"
    Write-Host "    cpa-scheduled-tests:"
    Write-Host "      enabled: true"
    Write-Host "      priority: 1"
}
Write-Host "Restart CPA once, then open the CPA Scheduled 5H menu."
