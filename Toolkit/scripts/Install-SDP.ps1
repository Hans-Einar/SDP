[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$ProjectRoot,

    [switch]$ForceManagedFiles,

    [switch]$InitializeProjectStructure,

    [switch]$Preview,

    [string]$BackupRoot
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$ToolkitVersion = '0.2.0'
$InstallerVersion = '0.2.0'
$FrameworkVersion = '1.0.0'
$AgentsContractVersion = '1.0.0'
$SupportedInstalledManifestSchemas = @('1.0')
$RunStamp = (Get-Date).ToUniversalTime().ToString('yyyyMMdd-HHmmssZ')

$ToolkitRoot = Split-Path -Parent $PSScriptRoot
$RepositoryRoot = [System.IO.Path]::GetFullPath((Split-Path -Parent $ToolkitRoot))
$ProjectRoot = [System.IO.Path]::GetFullPath((Resolve-Path $ProjectRoot).Path)
$SdpRoot = Join-Path $ProjectRoot 'SDP'
$SkillsTarget = Join-Path $ProjectRoot '.codex\skills'
$InstalledManifestPath = Join-Path $SdpRoot 'Framework\installed-toolkit.manifest.yaml'

if ([string]::IsNullOrWhiteSpace($BackupRoot)) {
    $BackupRoot = Join-Path $SdpRoot ".sdp-backups\$RunStamp"
} elseif (-not [System.IO.Path]::IsPathRooted($BackupRoot)) {
    $BackupRoot = Join-Path $ProjectRoot $BackupRoot
}
$BackupRoot = [System.IO.Path]::GetFullPath($BackupRoot)

# Compare complete path segments. A sibling such as SDP-Analyzer must not be
# rejected merely because its raw path starts with the same characters as SDP.
$repositoryTrimmed = $RepositoryRoot.TrimEnd(
    [System.IO.Path]::DirectorySeparatorChar,
    [System.IO.Path]::AltDirectorySeparatorChar
)
$repositoryPrefix = $repositoryTrimmed + [System.IO.Path]::DirectorySeparatorChar
$projectIsRepository = $ProjectRoot.Equals(
    $repositoryTrimmed,
    [System.StringComparison]::OrdinalIgnoreCase
)
$projectIsInsideRepository = $ProjectRoot.StartsWith(
    $repositoryPrefix,
    [System.StringComparison]::OrdinalIgnoreCase
)

if ($projectIsRepository -or $projectIsInsideRepository) {
    throw "The consuming project must not be inside the SDP repository: $RepositoryRoot"
}

$script:ProposedCount = 0
$script:AppliedCount = 0
$script:PreservedCount = 0
$script:UnchangedCount = 0

function Write-SdpAction {
    param(
        [ValidateSet('PROPOSED', 'APPLIED', 'PRESERVED', 'UNCHANGED', 'WARNING')]
        [string]$Kind,
        [string]$Message
    )

    switch ($Kind) {
        'PROPOSED' { $script:ProposedCount++ }
        'APPLIED' { $script:AppliedCount++ }
        'PRESERVED' { $script:PreservedCount++ }
        'UNCHANGED' { $script:UnchangedCount++ }
    }
    Write-Host "[$Kind] $Message"
}

function Get-YamlScalar {
    param(
        [string]$Content,
        [string]$Name
    )

    $escaped = [regex]::Escape($Name)
    $pattern = '(?m)^\s*{0}\s*:\s*[''"]?([^''"#\r\n]+)' -f $escaped
    $match = [regex]::Match(
        $Content,
        $pattern,
        [System.Text.RegularExpressions.RegexOptions]::CultureInvariant
    )
    if (-not $match.Success) {
        return $null
    }
    return $match.Groups[1].Value.Trim()
}

function Get-RelativeProjectPath {
    param([string]$Path)

    $full = [System.IO.Path]::GetFullPath($Path)
    if ($full.StartsWith($ProjectRoot, [System.StringComparison]::OrdinalIgnoreCase)) {
        return $full.Substring($ProjectRoot.Length).TrimStart([char[]]'\/')
    }
    return [System.IO.Path]::GetFileName($full)
}

function Backup-ExistingFile {
    param([string]$Path)

    if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) {
        return
    }

    $relative = Get-RelativeProjectPath $Path
    $destination = Join-Path $BackupRoot $relative
    if ($Preview) {
        Write-SdpAction 'PROPOSED' "Back up $relative to $destination"
        return
    }

    $parent = Split-Path -Parent $destination
    New-Item -ItemType Directory -Force -Path $parent | Out-Null
    Copy-Item -LiteralPath $Path -Destination $destination -Force
    Write-SdpAction 'APPLIED' "Backed up $relative to $destination"
}

function Install-TextFile {
    param(
        [string]$Content,
        [string]$Destination,
        [ValidateSet('Managed', 'ProjectOwned')]
        [string]$Ownership,
        [bool]$RefreshManaged,
        [bool]$SkipBackup = $false
    )

    $relative = Get-RelativeProjectPath $Destination
    if (-not (Test-Path -LiteralPath $Destination -PathType Leaf)) {
        if ($Preview) {
            Write-SdpAction 'PROPOSED' "Create $relative ($Ownership)"
            return
        }
        $parent = Split-Path -Parent $Destination
        New-Item -ItemType Directory -Force -Path $parent | Out-Null
        [System.IO.File]::WriteAllText($Destination, $Content, [System.Text.UTF8Encoding]::new($false))
        Write-SdpAction 'APPLIED' "Created $relative ($Ownership)"
        return
    }

    $existing = [System.IO.File]::ReadAllText($Destination)
    if ($existing -ceq $Content) {
        Write-SdpAction 'UNCHANGED' $relative
        return
    }

    if ($Ownership -eq 'ProjectOwned') {
        Write-SdpAction 'PRESERVED' "$relative (project-owned)"
        return
    }

    if (-not $RefreshManaged) {
        Write-SdpAction 'PRESERVED' "$relative (managed file differs; use -ForceManagedFiles to restore it)"
        return
    }

    if (-not $SkipBackup) {
        Backup-ExistingFile $Destination
    }
    if ($Preview) {
        Write-SdpAction 'PROPOSED' "Replace managed file $relative"
        return
    }

    [System.IO.File]::WriteAllText($Destination, $Content, [System.Text.UTF8Encoding]::new($false))
    Write-SdpAction 'APPLIED' "Replaced managed file $relative"
}

function Install-SourceFile {
    param(
        [string]$Source,
        [string]$Destination,
        [ValidateSet('Managed', 'ProjectOwned')]
        [string]$Ownership,
        [bool]$RefreshManaged,
        [bool]$SkipBackup = $false
    )

    $content = [System.IO.File]::ReadAllText($Source)
    Install-TextFile $content $Destination $Ownership $RefreshManaged $SkipBackup
}

# Detect compatibility before any mutation.
$InstalledSchemaVersion = $null
$InstalledToolkitVersion = $null
$PreviousInstalledAt = $null
if (Test-Path -LiteralPath $InstalledManifestPath -PathType Leaf) {
    $installedContent = [System.IO.File]::ReadAllText($InstalledManifestPath)
    $InstalledSchemaVersion = Get-YamlScalar $installedContent 'schemaVersion'
    $InstalledToolkitVersion = Get-YamlScalar $installedContent 'toolkitVersion'
    $PreviousInstalledAt = Get-YamlScalar $installedContent 'toolkitInstalledAt'

    if ([string]::IsNullOrWhiteSpace($InstalledSchemaVersion)) {
        Write-Warning "Installed manifest is malformed: $InstalledManifestPath"
        throw 'Refusing to modify an installation with an unreadable manifest.'
    }
    if ($SupportedInstalledManifestSchemas -notcontains $InstalledSchemaVersion) {
        Write-Warning "Unsupported installed manifest schema '$InstalledSchemaVersion'."
        throw 'Refusing to modify an unsupported SDP installation.'
    }
} else {
    Write-SdpAction 'PROPOSED' 'Migrate supported pre-versioning installation to manifest schema 1.0'
}

$IsUpgrade = [string]::IsNullOrWhiteSpace($InstalledToolkitVersion) -or
    ($InstalledToolkitVersion -ne $ToolkitVersion)
$RefreshManaged = $IsUpgrade -or $ForceManagedFiles

# AGENTS.md is always canonical. Preserve pre-existing project instructions before
# replacing it, without ever overwriting AGENTS-project.md.
$agentsSource = Join-Path $ToolkitRoot 'payload\project-root\AGENTS.md.template'
$agentsDestination = Join-Path $ProjectRoot 'AGENTS.md'
$projectInstructions = Join-Path $ProjectRoot 'AGENTS-project.md'
$agentsContent = [System.IO.File]::ReadAllText($agentsSource)

if (Test-Path -LiteralPath $agentsDestination -PathType Leaf) {
    $existingAgents = [System.IO.File]::ReadAllText($agentsDestination)
    if ($existingAgents -cne $agentsContent) {
        if (-not (Test-Path -LiteralPath $projectInstructions -PathType Leaf)) {
            if ($Preview) {
                Write-SdpAction 'PROPOSED' 'Migrate existing AGENTS.md to AGENTS-project.md'
            } else {
                Copy-Item -LiteralPath $agentsDestination -Destination $projectInstructions
                Write-SdpAction 'APPLIED' 'Migrated existing AGENTS.md to AGENTS-project.md'
            }
        } else {
            $migrationBackup = Join-Path $ProjectRoot "AGENTS-project.migration-$RunStamp.md"
            if ($Preview) {
                Write-SdpAction 'PROPOSED' "Preserve existing AGENTS.md as $migrationBackup"
            } else {
                Copy-Item -LiteralPath $agentsDestination -Destination $migrationBackup
                Write-SdpAction 'APPLIED' "Preserved existing AGENTS.md as $migrationBackup"
            }
        }
    }
}
Install-TextFile $agentsContent $agentsDestination 'Managed' $true $true

Install-SourceFile `
    (Join-Path $ToolkitRoot 'payload\project-root\AGENTS-project.md.template') `
    $projectInstructions `
    'ProjectOwned' `
    $false

Install-SourceFile `
    (Join-Path $ToolkitRoot 'payload\sdp-root\AGENT-REMINDERS.md.template') `
    (Join-Path $SdpRoot 'AGENT-REMINDERS.md') `
    'ProjectOwned' `
    $false

# Framework and skills are clearly Toolkit-managed. They refresh automatically on
# a version upgrade; same-version local differences require -ForceManagedFiles.
$frameworkSource = Join-Path $ToolkitRoot 'payload\sdp-root\Framework'
Get-ChildItem -Path $frameworkSource -Recurse -File | ForEach-Object {
    $relative = $_.FullName.Substring($frameworkSource.Length).TrimStart([char[]]'\/')
    Install-SourceFile $_.FullName (Join-Path $SdpRoot "Framework\$relative") 'Managed' $RefreshManaged
}

$skillsSource = Join-Path $ToolkitRoot 'skills'
Get-ChildItem -Path $skillsSource -Directory | Sort-Object Name | ForEach-Object {
    $skillFile = Join-Path $_.FullName 'SKILL.md'
    if (Test-Path -LiteralPath $skillFile -PathType Leaf) {
        Install-SourceFile $skillFile (Join-Path $SkillsTarget "$($_.Name)\SKILL.md") 'Managed' $RefreshManaged
    }
}

# Generated installed facts are Toolkit-owned. Preserve the timestamp on a
# same-version reinstall so repeat installation is idempotent.
if (($InstalledToolkitVersion -eq $ToolkitVersion) -and
    (-not [string]::IsNullOrWhiteSpace($PreviousInstalledAt))) {
    $InstalledAt = $PreviousInstalledAt
} else {
    $InstalledAt = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')
}

$SourceCommit = $null
if (Get-Command git -ErrorAction SilentlyContinue) {
    $gitOutput = & git -C $RepositoryRoot rev-parse HEAD 2>$null
    if (($LASTEXITCODE -eq 0) -and $gitOutput) {
        $SourceCommit = ($gitOutput | Select-Object -First 1).Trim()
    }
}
$sourceCommitYaml = if ($SourceCommit) { '"' + $SourceCommit + '"' } else { 'null' }

$installedManifest = @"
schemaVersion: "1.0"
toolkitVersion: "$ToolkitVersion"
frameworkVersion: "$FrameworkVersion"
agentsContractVersion: "$AgentsContractVersion"
installerVersion: "$InstallerVersion"
toolkitInstalledAt: "$InstalledAt"
sourceCommit: $sourceCommitYaml
skills:
  sdp-architect: "1.0.0"
  sdp-auditor: "1.0.0"
  sdp-master: "1.0.0"
  sdp-release: "1.0.0"
  sdp-reviewer: "1.0.0"
  sdp-traceability: "1.0.0"
  sdp-verifier: "1.0.0"
  sdp-versioning: "1.0.0"
  sdp-vertical-refactor: "1.0.0"
  sdp-worker: "1.0.0"
capabilities:
  - sdp.install.v1
  - sdp.manifest.v1
  - sdp.release.v1
  - sdp.traceability.release-events.v1
  - sdp.versioning.v1
"@
Install-TextFile $installedManifest $InstalledManifestPath 'Managed' $true

# Project-owned release state is created only when missing.
Install-SourceFile `
    (Join-Path $frameworkSource 'templates\SDP-project.manifest.yaml') `
    (Join-Path $SdpRoot 'SDP-project.manifest.yaml') `
    'ProjectOwned' `
    $false
Install-SourceFile `
    (Join-Path $frameworkSource 'templates\RELEASE-NOTES.md') `
    (Join-Path $SdpRoot 'RELEASE-NOTES.md') `
    'ProjectOwned' `
    $false

# Missing traceability contracts are safe additive migrations. Populated files
# remain project-owned and are never replaced.
foreach ($traceFile in @('README.md', 'CurrentIndex.yaml', 'Relations.yaml')) {
    $source = Join-Path $RepositoryRoot "Traceability\$traceFile"
    if (Test-Path -LiteralPath $source -PathType Leaf) {
        Install-SourceFile $source (Join-Path $SdpRoot "Traceability\$traceFile") 'ProjectOwned' $false
    }
}

if ($InitializeProjectStructure) {
    $templateFolders = @(
        '01--Mandate', '02--Study', '03--Requirements', '04--Architecture',
        '05--DesignAnalysis', '06--Design', '07--Implementation',
        'Sprints', 'Refactors', 'Fixes', 'Releases', 'CodeReview', 'Verification',
        'Traceability', 'Instructions'
    )

    foreach ($folder in $templateFolders) {
        $source = Join-Path $RepositoryRoot $folder
        if (Test-Path -LiteralPath $source -PathType Container) {
            Get-ChildItem -Path $source -Recurse -File | ForEach-Object {
                $relative = $_.FullName.Substring($source.Length).TrimStart([char[]]'\/')
                Install-SourceFile $_.FullName (Join-Path $SdpRoot "$folder\$relative") 'ProjectOwned' $false
            }
        }
    }

    Install-SourceFile `
        (Join-Path $RepositoryRoot 'SDP-DOCUMENT-GUIDE.md') `
        (Join-Path $SdpRoot 'SDP-DOCUMENT-GUIDE.md') `
        'ProjectOwned' `
        $false
}

Write-Host ''
Write-Host 'SDP installation summary'
Write-Host "  Toolkit version: $ToolkitVersion"
Write-Host "  Mode: $(if ($Preview) { 'preview' } else { 'apply' })"
Write-Host "  Upgrade: $IsUpgrade"
Write-Host "  Proposed: $script:ProposedCount"
Write-Host "  Applied: $script:AppliedCount"
Write-Host "  Preserved: $script:PreservedCount"
Write-Host "  Unchanged: $script:UnchangedCount"
Write-Host "  Backup root: $BackupRoot"
Write-Host 'No release or Git tag was created.'
