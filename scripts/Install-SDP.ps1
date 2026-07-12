[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$ProjectRoot,

    [switch]$ForceManagedFiles,

    [switch]$InitializeProjectStructure,

    [switch]$Preview,

    [string]$BackupRoot
)

$ToolkitInstaller = Join-Path (Split-Path -Parent $PSScriptRoot) 'Toolkit\scripts\Install-SDP.ps1'
Write-Warning 'scripts\Install-SDP.ps1 is a compatibility entry point. Prefer Toolkit\scripts\Install-SDP.ps1.'

$arguments = @{
    ProjectRoot = $ProjectRoot
    ForceManagedFiles = $ForceManagedFiles
    InitializeProjectStructure = $InitializeProjectStructure
    Preview = $Preview
}
if (-not [string]::IsNullOrWhiteSpace($BackupRoot)) {
    $arguments.BackupRoot = $BackupRoot
}

& $ToolkitInstaller @arguments
