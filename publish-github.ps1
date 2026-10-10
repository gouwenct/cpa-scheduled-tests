param(
    [string]$RepoName = "cpa-scheduled-tests",
    [ValidateSet("private", "public")]
    [string]$Visibility = "private",
    [switch]$CreateRelease
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

function Require-Command([string]$Name, [string]$Hint) {
    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "$Name was not found. $Hint"
    }
}

Require-Command git "Install Git for Windows first."
Require-Command gh "Install GitHub CLI first: winget install --id GitHub.cli"

Write-Host "Checking GitHub login..." -ForegroundColor Cyan
& gh auth status 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "GitHub CLI is not logged in. Starting login..." -ForegroundColor Yellow
    & gh auth login
    if ($LASTEXITCODE -ne 0) { throw "GitHub login failed." }
}

$login = (& gh api user --jq .login).Trim()
if (-not $login) { throw "Could not determine the current GitHub username." }

$repoSlug = "$login/$RepoName"
$repoUrl = "https://github.com/$repoSlug"
Write-Host "Target repository: $repoUrl" -ForegroundColor Green

# Keep plugin metadata and registry aligned with the real repository.
$mainPath = Join-Path $PSScriptRoot "main.go"
$registryPath = Join-Path $PSScriptRoot "registry.example.json"

$main = Get-Content $mainPath -Raw
$main = [regex]::Replace(
    $main,
    'GitHubRepository:\s*"[^"]+"',
    ('GitHubRepository: "' + $repoUrl + '"')
)
Set-Content -Path $mainPath -Value $main -Encoding UTF8

if (Test-Path $registryPath) {
    $registry = Get-Content $registryPath -Raw
    $registry = $registry.Replace("https://github.com/YOUR_ACCOUNT/cpa-scheduled-tests", $repoUrl)
    $registry = [regex]::Replace($registry, '"author"\s*:\s*"[^"]*"', ('"author": "' + $login + '"'), 1)
    Set-Content -Path $registryPath -Value $registry -Encoding UTF8
}

Push-Location $PSScriptRoot
try {
    if (-not (Test-Path ".git")) {
        & git init
        & git branch -M main
    }

    if (-not (& git config user.name)) {
        & git config user.name $login
    }
    if (-not (& git config user.email)) {
        $email = (& gh api user --jq '.email // empty').Trim()
        if (-not $email) { $email = "$login@users.noreply.github.com" }
        & git config user.email $email
    }

    & git add .
    $hasHead = $true
    & git rev-parse --verify HEAD *> $null
    if ($LASTEXITCODE -ne 0) { $hasHead = $false }

    & git diff --cached --quiet
    $hasChanges = ($LASTEXITCODE -ne 0)
    if ($hasChanges -or -not $hasHead) {
        & git commit -m "Initial release: CPA Scheduled 5H v0.1.6"
    }

    & gh repo view $repoSlug *> $null
    $repoExists = ($LASTEXITCODE -eq 0)

    if (-not $repoExists) {
        $visArg = if ($Visibility -eq "public") { "--public" } else { "--private" }
        & gh repo create $repoSlug $visArg --source . --remote origin --push
        if ($LASTEXITCODE -ne 0) { throw "Repository creation/push failed." }
    } else {
        $origin = (& git remote get-url origin 2>$null)
        if ($LASTEXITCODE -ne 0 -or -not $origin) {
            & git remote add origin "https://github.com/$repoSlug.git"
        } else {
            & git remote set-url origin "https://github.com/$repoSlug.git"
        }
        & git push -u origin main
        if ($LASTEXITCODE -ne 0) { throw "git push failed." }
    }

    if ($CreateRelease) {
        $tag = "v0.1.6"
        & gh release view $tag --repo $repoSlug *> $null
        if ($LASTEXITCODE -ne 0) {
            $assets = @()
            foreach ($platform in @("darwin_amd64", "darwin_arm64", "linux_amd64", "linux_arm64", "windows_amd64")) {
                $package = Join-Path $PSScriptRoot "dist\cpa-scheduled-tests_0.1.6_$platform.zip"
                if (-not (Test-Path -LiteralPath $package)) {
                    throw "Missing full-platform asset: $package. Use the platform-builds GitHub Actions workflow."
                }
                $assets += $package
            }
            $checksumPath = Join-Path $PSScriptRoot "dist\checksums.txt"
            $checksumLines = @($assets | ForEach-Object { (Get-FileHash -LiteralPath $_ -Algorithm SHA256).Hash.ToLowerInvariant() + "  " + [IO.Path]::GetFileName($_) })
            if ($checksumLines.Count -gt 0) {
                [IO.File]::WriteAllText($checksumPath, ($checksumLines -join "`n") + "`n", [Text.UTF8Encoding]::new($false))
                $assets += $checksumPath
            }

            $args = @("release", "create", $tag, "--repo", $repoSlug, "--title", "CPA Scheduled 5H v0.1.6", "--generate-notes")
            $args += $assets
            & gh @args
            if ($LASTEXITCODE -ne 0) { throw "GitHub Release creation failed." }
        }
    }

    Write-Host "" 
    Write-Host "Published successfully:" -ForegroundColor Green
    Write-Host $repoUrl -ForegroundColor Cyan
}
finally {
    Pop-Location
}
