# Claude Usage - Windows Installation Script
# Run as: powershell -ExecutionPolicy Bypass -File install-windows.ps1

param(
    [switch]$Uninstall,
    [switch]$NoAutostart
)

$AppName = "Claude Usage"
$ExeName = "claude-usage.exe"
$InstallDir = "$env:LOCALAPPDATA\Programs\claude-usage"
$StartMenuDir = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs"
$StartupDir = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup"

function Write-Status {
    param([string]$Message)
    Write-Host "[*] $Message" -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host "[+] $Message" -ForegroundColor Green
}

function Write-Error {
    param([string]$Message)
    Write-Host "[-] $Message" -ForegroundColor Red
}

function Create-Shortcut {
    param(
        [string]$ShortcutPath,
        [string]$TargetPath,
        [string]$Description,
        [string]$IconPath
    )
    
    $WshShell = New-Object -ComObject WScript.Shell
    $Shortcut = $WshShell.CreateShortcut($ShortcutPath)
    $Shortcut.TargetPath = $TargetPath
    $Shortcut.Description = $Description
    if ($IconPath) {
        $Shortcut.IconLocation = $IconPath
    }
    $Shortcut.Save()
}

function Install-App {
    Write-Status "Installing $AppName..."
    
    # Check if exe exists in current directory
    $SourceExe = Join-Path $PSScriptRoot $ExeName
    if (-not (Test-Path $SourceExe)) {
        $SourceExe = $ExeName
    }
    if (-not (Test-Path $SourceExe)) {
        Write-Error "Cannot find $ExeName in current directory"
        Write-Host "Please run this script from the directory containing $ExeName"
        exit 1
    }
    
    # Create install directory
    Write-Status "Creating install directory: $InstallDir"
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    
    # Copy executable
    Write-Status "Copying executable..."
    Copy-Item $SourceExe -Destination "$InstallDir\$ExeName" -Force
    
    # Copy icon if exists
    $IconPath = Join-Path $PSScriptRoot "claude-usage.ico"
    if (Test-Path $IconPath) {
        Copy-Item $IconPath -Destination "$InstallDir\claude-usage.ico" -Force
    }
    
    # Create Start Menu shortcut
    Write-Status "Creating Start Menu shortcut..."
    $StartMenuShortcut = Join-Path $StartMenuDir "$AppName.lnk"
    Create-Shortcut -ShortcutPath $StartMenuShortcut `
                    -TargetPath "$InstallDir\$ExeName" `
                    -Description "Monitor Claude Code token usage" `
                    -IconPath "$InstallDir\$ExeName"
    
    # Create Startup shortcut (unless disabled)
    if (-not $NoAutostart) {
        Write-Status "Creating Startup shortcut (autostart)..."
        $StartupShortcut = Join-Path $StartupDir "$AppName.lnk"
        Create-Shortcut -ShortcutPath $StartupShortcut `
                        -TargetPath "$InstallDir\$ExeName" `
                        -Description "Monitor Claude Code token usage" `
                        -IconPath "$InstallDir\$ExeName"
    }
    
    Write-Success "Installation complete!"
    Write-Host ""
    Write-Host "Installed to: $InstallDir"
    Write-Host "Start Menu:   $StartMenuShortcut"
    if (-not $NoAutostart) {
        Write-Host "Autostart:    $StartupShortcut"
    }
    Write-Host ""
    Write-Host "You can now launch '$AppName' from the Start Menu."
}

function Uninstall-App {
    Write-Status "Uninstalling $AppName..."
    
    # Stop running instance
    $process = Get-Process -Name "claude-usage" -ErrorAction SilentlyContinue
    if ($process) {
        Write-Status "Stopping running instance..."
        Stop-Process -Name "claude-usage" -Force
        Start-Sleep -Seconds 1
    }
    
    # Remove Start Menu shortcut
    $StartMenuShortcut = Join-Path $StartMenuDir "$AppName.lnk"
    if (Test-Path $StartMenuShortcut) {
        Write-Status "Removing Start Menu shortcut..."
        Remove-Item $StartMenuShortcut -Force
    }
    
    # Remove Startup shortcut
    $StartupShortcut = Join-Path $StartupDir "$AppName.lnk"
    if (Test-Path $StartupShortcut) {
        Write-Status "Removing Startup shortcut..."
        Remove-Item $StartupShortcut -Force
    }
    
    # Remove install directory
    if (Test-Path $InstallDir) {
        Write-Status "Removing install directory..."
        Remove-Item $InstallDir -Recurse -Force
    }
    
    Write-Success "Uninstallation complete!"
}

# Main
if ($Uninstall) {
    Uninstall-App
} else {
    Install-App
}
