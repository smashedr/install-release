#!/usr/bin/env pwsh
# https://raw.githubusercontent.com/smashedr/install-release/refs/heads/master/scripts/install.ps1
param(
    [string]$bin=""
)

$ErrorActionPreference = "Stop"

$exeName = "ir"
$repository = "smashedr/install-release"

Write-Host -ForegroundColor Cyan "Installing: $repository"

## ARCH
$platform = switch ($true) {
    $IsWindows { "Windows" }
    $IsLinux   { "Linux" }
    $IsMacOS   { "Darwin" }
    default    { "unknown" }
}
Write-Host -ForegroundColor DarkCyan "platform: $platform"
$osArchitecture = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture
Write-Host -ForegroundColor DarkCyan "osArchitecture: $osArchitecture"
$arch = switch ($osArchitecture) {
    "X64"   { "x86_64" }
    "Arm64" { "arm64" }
    "X86"   { "i386" }
    default { "unknown" }
}
Write-Host -ForegroundColor DarkCyan "arch: $arch"

## FILE
if ($IsWindows) {
    $exeName = "${exeName}.exe"
    $binPath = Join-Path $env:LOCALAPPDATA "Microsoft\WindowsApps"
    $file = "${exeName}_${platform}_${arch}.zip"
} else {
    $binPath = Join-Path $HOME ".local/bin"
    $file = "${exeName}_${platform}_${arch}.tar.gz"
}
Write-Host -ForegroundColor DarkCyan "exeName: $exeName"
Write-Host -ForegroundColor DarkCyan "binPath: $binPath"
Write-Host -ForegroundColor DarkCyan "file: $file"
$url = "https://github.com/$repository/releases/latest/download/$file"
Write-Host -ForegroundColor DarkCyan "url: $url"

## BIN
Write-Host -ForegroundColor White "Target Directory: $binPath"
if ($bin) {
    $binPath = $bin
} else {
    $userInput = Read-Host "Enter Path [press <enter> to accept]"
    $binPath = if ($userInput) { $userInput } else { $binPath }
}

if (-not (Test-Path -IsValid $binPath)) {
    Write-Host -ForegroundColor Red "Invalid path: $binPath"
    exit 1
}
if (-not (Test-Path $binPath)) {
    Write-Host -ForegroundColor Red "Directory does not exist: $binPath"
    exit 1
}

## PATH
if ($IsWindows) {
    # Windows
    $pathUser = [Environment]::GetEnvironmentVariable("Path", "User")
    $paths = $pathUser -split ";" | ForEach-Object { $_.TrimEnd('\') }
    $binPath = $binPath.TrimEnd('\')
    Write-Host -ForegroundColor DarkCyan "binPath: $binPath"
    if ($paths -notcontains $binPath) {
        Write-Host -ForegroundColor Yellow "Adding PATH: $binPath"
        [Environment]::SetEnvironmentVariable("Path", "$pathUser;$bin", "User")
    } else {
        Write-Host -ForegroundColor DarkGreen "Already in PATH: $binPath"
    }
} else {
    # Unix
    if ($env:PATH -split ':' -notcontains $binPath) {
        Write-Host -ForegroundColor Yellow "Adding PATH: $binPath"
        $env:PATH += [IO.Path]::PathSeparator + $binPath
        if (!(Test-Path -Path $PROFILE)) {
            New-Item -ItemType File -Path $PROFILE -Force | Out-Null
        }
        Add-Content -Path $PROFILE -Value "`$env:PATH += ':$binPath'"
    } else {
        Write-Host -ForegroundColor DarkGreen "Already in PATH: $binPath"
    }
}

## TEMP
$temp = [system.io.path]::GetTempPath()
Write-Host -ForegroundColor DarkCyan "temp: $temp"
$tempDir = Join-Path $temp "install_$(Get-Random)"
Write-Host -ForegroundColor DarkCyan "tempDir: $tempDir"
$zipPath = Join-Path $tempDir $file
Write-Host -ForegroundColor DarkCyan "zipPath: $zipPath"

## EXEC
try {
    # Download
    Write-Host -ForegroundColor DarkCyan "Downloading: $url"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    Invoke-WebRequest -Uri $url -OutFile $zipPath
    # Extract
    if ($zipPath -like "*.tar.gz") {
        tar -xzf $zipPath -C $tempDir
    } else {
        Expand-Archive -Path $zipPath -DestinationPath $tempDir -Force
    }
    # Install
    $source = Join-Path $tempDir $exeName
    Move-Item -Path $source -Destination $binPath -Force
} catch {
    Write-Host -ForegroundColor Red "Error: $_"
    exit 1
} finally {
    if (Test-Path $tempDir) {
        Write-Host -ForegroundColor DarkCyan "Cleaning Up: $tempDir"
        Remove-Item -Path $tempDir -Recurse -Force
    }
}

$location = Join-Path $binPath $exeName
Write-Host -ForegroundColor DarkCyan "Location: $location "
Write-Host -ForegroundColor Green "Installation Successful!"
Write-Host -ForegroundColor White "To get started run: ir --help"
