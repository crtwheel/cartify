$ErrorActionPreference = 'Stop'

$CartifyFolderPath = "$env:LOCALAPPDATA\Cartify"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

function Write-Success { Write-Host ' > OK' -ForegroundColor 'Green' }
function Write-Unsuccess { Write-Host ' > ERROR' -ForegroundColor 'Red' }

function Test-Admin {
    Write-Host "Checking if not running as administrator..." -NoNewline
    $currentUser = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    $result = -not $currentUser.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    if ($result) { Write-Success } else { Write-Unsuccess }
    return $result
}

function Test-PowerShellVersion {
    Write-Host 'Checking PowerShell version...' -NoNewline
    $result = $PSVersionTable.PSVersion -ge [version]'5.1'
    if ($result) { Write-Success } else { Write-Unsuccess }
    return $result
}

function Add-ToPath {
    Write-Host 'Adding Cartify to PATH...' -NoNewline
    $user = [EnvironmentVariableTarget]::User
    $path = [Environment]::GetEnvironmentVariable('PATH', $user)
    if ($path -notlike "*$CartifyFolderPath*") {
        $path = "$path;$CartifyFolderPath"
        [Environment]::SetEnvironmentVariable('PATH', $path, $user)
    }
    $env:PATH = $path
    Write-Success
}

# --- Main ---
if (-not (Test-PowerShellVersion)) {
    Write-Warning 'PowerShell 5.1 or higher required.'
    Pause; exit
}
if (-not (Test-Admin)) {
    Write-Warning "Don't run as administrator. Continuing anyway."
}

Write-Host "Installing Cartify to $CartifyFolderPath"

if (-not (Test-Path -Path $ScriptDir\cartify.exe)) {
    Write-Host "cartify.exe not found in current directory. Run this script from the extracted Cartify folder." -ForegroundColor Red
    Pause; exit
}

New-Item -ItemType Directory -Path $CartifyFolderPath -Force | Out-Null
Copy-Item -Path "$ScriptDir\*" -Destination $CartifyFolderPath -Recurse -Force
Write-Host "Files copied." -ForegroundColor Green

Add-ToPath

Write-Host "`nRunning cartify install to finalize setup..." -ForegroundColor Cyan
& "$CartifyFolderPath\cartify.exe" install

Write-Host "`nCartify installed! Restart Spotify if it's running." -ForegroundColor Green
Write-Host "You can always re-apply with: cartify auto" -ForegroundColor Cyan
