# https://raw.githubusercontent.com/smashedr/install-release/refs/heads/master/scripts/install.ps1
param(
    [string]$bin="$env:LOCALAPPDATA\Microsoft\WindowsApps"
)

$ErrorActionPreference = "Stop"

$exeName = "ir.exe"

#Write-Host -ForegroundColor DarkCyan "bin: $bin"
#Write-Host -ForegroundColor DarkCyan "Architecture: $env:PROCESSOR_ARCHITECTURE"

$file = switch ($env:PROCESSOR_ARCHITECTURE) {
    "AMD64" { "ir_Windows_x86_64.zip" }
    "x86"   { "ir_Windows_i386.zip" }
    "ARM64" { "ir_Windows_arm64.zip" }
    default { throw "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE" }
}
#Write-Host -ForegroundColor DarkCyan "File: $file"

$url = "https://github.com/smashedr/install-release/releases/latest/download/$file"

Write-Host -ForegroundColor White "Current: $bin"
$userInput = Read-Host "Enter Path [press <enter> to accept]"
$binPath = if ($userInput) { $userInput } else { $binPath }
$bin = if ($userInput) { $userInput } else { $bin }

if (-not (Test-Path -IsValid $bin)) {
    Write-Host -ForegroundColor Red "Invalid path: $bin"
    exit 1
}
if (-not (Test-Path $bin)) {
    Write-Host -ForegroundColor Red "Directory does not exist: $bin"
    exit 1
}

$pathUser = [Environment]::GetEnvironmentVariable("Path", "User")
$paths = $pathUser -split ";" | ForEach-Object { $_.TrimEnd('\') }
$binPath = $bin.TrimEnd('\')
#Write-Host -ForegroundColor DarkCyan "binPath: $binPath"

if ($paths -notcontains $binPath) {
    Write-Host -ForegroundColor Yellow "Adding PATH: $binPath"
    [Environment]::SetEnvironmentVariable("Path", "$pathUser;$bin", "User")
} else {
    Write-Host -ForegroundColor DarkGreen "Already in PATH: $binPath"
}

$tempDir = Join-Path $env:TEMP "install_$(Get-Random)"
#Write-Host -ForegroundColor DarkCyan "tempDir: $tempDir"
$zipPath = Join-Path $tempDir "download.zip"
#Write-Host -ForegroundColor DarkCyan "zipPath: $zipPath"

try {
    Write-Host -ForegroundColor DarkCyan "Downloading: $url"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    Invoke-WebRequest -Uri $url -OutFile $zipPath
    Expand-Archive -Path $zipPath -DestinationPath $tempDir -Force
    $source = Join-Path $tempDir $exeName
    Move-Item -Path $source -Destination $binPath -Force
} catch {
    Write-Host -ForegroundColor Red "Error: $_"
    exit 1
} finally {
    if (Test-Path $tempDir) {
        #Write-Host -ForegroundColor DarkCyan "Cleaning Up: $tempDir"
        Remove-Item -Path $tempDir -Recurse -Force
    }
}

$location = Join-Path $binPath $exeName
Write-Host -ForegroundColor DarkCyan "Location: $location "
Write-Host -ForegroundColor Green "Installation Successful!"
Write-Host -ForegroundColor White "To get started, run: ir --help"
