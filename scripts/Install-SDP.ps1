[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$ProjectRoot,

    [switch]$ForceManagedFiles
)

$ErrorActionPreference = 'Stop'
$RepoRoot = (Resolve-Path (Split-Path -Parent $PSScriptRoot)).Path
$ProjectRoot = (Resolve-Path $ProjectRoot).Path
$SdpRoot = Join-Path $ProjectRoot 'SDP'
$SkillsTarget = Join-Path $ProjectRoot '.codex\skills'

if ($RepoRoot.StartsWith($ProjectRoot, [System.StringComparison]::OrdinalIgnoreCase)) {
    throw "The SDP toolkit repository must not be cloned inside the consuming project. Clone it beside the project and run this installer from there."
}

function Copy-SafeFile {
    param(
        [string]$Source,
        [string]$Destination,
        [bool]$Managed
    )

    $parent = Split-Path -Parent $Destination
    New-Item -ItemType Directory -Force -Path $parent | Out-Null

    if (Test-Path $Destination) {
        if ($Managed -and $ForceManagedFiles) {
            Copy-Item -Force $Source $Destination
            Write-Host "Updated managed file: $Destination"
        } else {
            Write-Host "Preserved existing file: $Destination"
        }
    } else {
        Copy-Item $Source $Destination
        Write-Host "Installed: $Destination"
    }
}

$existingSdp = Test-Path $SdpRoot
New-Item -ItemType Directory -Force -Path $SdpRoot | Out-Null

if ($existingSdp) {
    Write-Host "Using existing project SDP directory: $SdpRoot"
} else {
    Write-Host "Created project SDP directory: $SdpRoot"
}

Copy-SafeFile `
    (Join-Path $RepoRoot 'payload\project-root\AGENTS.md.template') `
    (Join-Path $ProjectRoot 'AGENTS.md') `
    $false

Copy-SafeFile `
    (Join-Path $RepoRoot 'payload\sdp-root\AGENT-REMINDERS.md.template') `
    (Join-Path $SdpRoot 'AGENT-REMINDERS.md') `
    $false

$frameworkSource = Join-Path $RepoRoot 'payload\sdp-root\Framework'
Get-ChildItem -Path $frameworkSource -Recurse -File | ForEach-Object {
    $relative = $_.FullName.Substring($frameworkSource.Length).TrimStart('\')
    Copy-SafeFile $_.FullName (Join-Path $SdpRoot "Framework\$relative") $true
}

$skillsSource = Join-Path $RepoRoot 'skills'
Get-ChildItem -Path $skillsSource -Directory | ForEach-Object {
    $skillFile = Join-Path $_.FullName 'SKILL.md'
    if (Test-Path $skillFile) {
        Copy-SafeFile $skillFile (Join-Path $SkillsTarget "$($_.Name)\SKILL.md") $true
    }
}

Write-Host ''
Write-Host 'SDP installation complete.'
Write-Host 'Existing project-specific SDP files were preserved.'
Write-Host 'Use -ForceManagedFiles to refresh only Framework and skill files.'