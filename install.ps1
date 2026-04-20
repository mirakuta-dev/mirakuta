#Requires -Version 5.1
<#
.SYNOPSIS
    Mirakuta installer for Windows.

.DESCRIPTION
    Downloads the signed Mirakuta release binary matching the host architecture,
    verifies its SHA-256 checksum, installs it under %LOCALAPPDATA%\Programs\mirakuta,
    and appends the install directory to the user PATH.

.PARAMETER Version
    Release tag to install (e.g. v0.1.0). Defaults to 'latest'.

.PARAMETER InstallDir
    Target install directory. Defaults to %LOCALAPPDATA%\Programs\mirakuta.

.EXAMPLE
    irm https://raw.githubusercontent.com/mirakuta-dev/mirakuta/main/install.ps1 | iex

.EXAMPLE
    & ([scriptblock]::Create((irm https://raw.githubusercontent.com/mirakuta-dev/mirakuta/main/install.ps1))) -Version v0.1.0
#>
[CmdletBinding()]
param(
    [string]$Version = 'latest',
    [string]$InstallDir = (Join-Path $env:LOCALAPPDATA 'Programs\mirakuta')
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$ProgressPreference    = 'SilentlyContinue'

[Net.ServicePointManager]::SecurityProtocol = `
    [Net.SecurityProtocolType]::Tls12 -bor [Net.SecurityProtocolType]::Tls13

$Repo       = 'mirakuta-dev/mirakuta'
$BinaryName = 'mirakuta.exe'

function Write-Step($msg) { Write-Host "==> $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "[ok] $msg" -ForegroundColor Green }

$arch = switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { 'amd64' }
    'ARM64' { 'arm64' }
    default {
        throw "Unsupported architecture: $($env:PROCESSOR_ARCHITECTURE). Supported: AMD64, ARM64."
    }
}
Write-Step "Architecture: windows-$arch"

if ($Version -eq 'latest') {
    Write-Step "Resolving latest release..."
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
    $Version = $release.tag_name
}
if (-not $Version) { throw "Could not determine release version." }
Write-Step "Version: $Version"

$assetName   = "mirakuta-windows-$arch.exe"
$assetUrl    = "https://github.com/$Repo/releases/download/$Version/$assetName"
$checksumUrl = "https://github.com/$Repo/releases/download/$Version/checksums.txt"

$targetBinary = Join-Path $InstallDir $BinaryName
if (Test-Path $targetBinary) {
    $running = Get-Process -Name 'mirakuta' -ErrorAction SilentlyContinue |
        Where-Object { $_.Path -eq $targetBinary }
    if ($running) {
        throw "mirakuta is currently running at $targetBinary. Close it and retry."
    }
}

$tempDir = Join-Path $env:TEMP "mirakuta-install-$([Guid]::NewGuid())"
New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
try {
    $tempBinary   = Join-Path $tempDir $assetName
    $tempChecksum = Join-Path $tempDir 'checksums.txt'

    Write-Step "Downloading $assetName..."
    Invoke-WebRequest -Uri $assetUrl    -OutFile $tempBinary   -UseBasicParsing
    Write-Step "Downloading checksums.txt..."
    Invoke-WebRequest -Uri $checksumUrl -OutFile $tempChecksum -UseBasicParsing

    Write-Step "Verifying SHA-256..."
    $expectedLine = Get-Content $tempChecksum |
        Where-Object { $_ -match ("\s" + [regex]::Escape($assetName) + "\s*$") } |
        Select-Object -First 1
    if (-not $expectedLine) {
        throw "No checksum entry for $assetName in checksums.txt"
    }
    $expected = ($expectedLine -split '\s+')[0].ToLower()
    $actual   = (Get-FileHash -Path $tempBinary -Algorithm SHA256).Hash.ToLower()
    if ($actual -ne $expected) {
        throw "Checksum mismatch for ${assetName}: expected $expected, got $actual"
    }
    Write-Ok "Checksum verified"

    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    Move-Item -Path $tempBinary -Destination $targetBinary -Force
    Write-Ok "Installed $targetBinary"
}
finally {
    Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
}

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
$entries  = @()
if ($userPath) { $entries = $userPath -split ';' | Where-Object { $_ } }
if ($entries -notcontains $InstallDir) {
    $newPath = if ($userPath) { "$userPath;$InstallDir" } else { $InstallDir }
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    Write-Ok "Added to user PATH: $InstallDir"
} else {
    Write-Step "Already in user PATH"
}
if (($env:Path -split ';') -notcontains $InstallDir) {
    $env:Path = "$env:Path;$InstallDir"
}

Write-Host ""
& $targetBinary version
Write-Host ""
Write-Host "Open a new terminal, then run 'mirakuta install' to get started." -ForegroundColor Cyan
