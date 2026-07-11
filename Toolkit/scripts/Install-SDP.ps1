[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$ProjectRoot,

    [switch]$ForceManagedFiles,

    [switch]$InitializeProjectStructure
)

$ErrorActionPreference = 'Stop'
$ToolkitRoot = Split-Path -Parent $PSScriptRoot
$RepositoryRoot = Split-Path -Parent $ToolkitRoot
$ProjectRoot = (Resolve-Path $ProjectRoot).Path
$SdpRoot = Join-Path $ProjectRoot 'SDP'
$SkillsTarget = Join-Path $ProjectRoot '.codex\skills'

if ($ProjectRoot.StartsWith($RepositoryRoot, [System.StringComparison]::OrdinalIgnoreCase)) {
    throw "The consuming project must not be inside the SDP toolkit repository: $RepositoryRoot"
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

# Existing non-empty SDP directories are explicitly supported.
New-Item -ItemType Directory -Force -Path $SdpRoot | Out-Null

Copy-SafeFile `
    (Join-Path $ToolkitRoot 'payload\project-root\AGENTS.md.template') `
    (Join-Path $ProjectRoot 'AGENTS.md') `
    $false

Copy-SafeFile `
    (Join-Path $ToolkitRoot 'payload\sdp-root\AGENT-REMINDERS.md.template') `
    (Join-Path $SdpRoot 'AGENT-REMINDERS.md') `
    $false

$frameworkSource = Join-Path $ToolkitRoot 'payload\sdp-root\Framework'
Get-ChildItem -Path $frameworkSource -Recurse -File | ForEach-Object {
    $relative = $_.FullName.Substring($frameworkSource.Length).TrimStart('\')
    Copy-SafeFile $_.FullName (Join-Path $SdpRoot "Framework\$relative") $true
}

$skillsSource = Join-Path $ToolkitRoot 'skills'
Get-ChildItem -Path $skillsSource -Directory | ForEach-Object {
    $skillFile = Join-Path $_.FullName 'SKILL.md'
    if (Test-Path $skillFile) {
        Copy-SafeFile $skillFile (Join-Path $SkillsTarget "$($_.Name)\SKILL.md") $true
    }
}

if ($InitializeProjectStructure) {
    $templateFolders = @(
        '01--Mandate', '02--Study', '03--Requirements', '04--Architecture',
        '05--DesignAnalysis', '06--Design', '07--Implementation',
        'Sprints', 'Refactors', 'CodeReview', 'Verification', 'Traceability',
        'Instructions'
    )

    foreach ($folder in $templateFolders) {
        $source = Join-Path $RepositoryRoot $folder
        if (Test-Path $source) {
            Get-ChildItem -Path $source -Recurse -File | ForEach-Object {
                $relative = $_.FullName.Substring($source.Length).TrimStart('\')
                Copy-SafeFile $_.FullName (Join-Path $SdpRoot "$folder\$relative") $false
            }
        }
    }

    Copy-SafeFile `
        (Join-Path $RepositoryRoot 'SDP-DOCUMENT-GUIDE.md') `
        (Join-Path $SdpRoot 'SDP-DOCUMENT-GUIDE.md') `
        $false
}

Write-Host ''
Write-Host 'SDP installation complete.'
Write-Host 'Existing project-specific files were preserved.'
Write-Host 'Use -ForceManagedFiles to refresh Toolkit-managed Framework and skill files.'
Write-Host 'Use -InitializeProjectStructure to add only missing standard template documents.'