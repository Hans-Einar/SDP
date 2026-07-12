[CmdletBinding()]
param()

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$RepositoryRoot = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path
$Installer = Join-Path $RepositoryRoot 'Toolkit\scripts\Install-SDP.ps1'
$TestRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("sdp-installer-tests-" + [guid]::NewGuid().ToString('N'))

function Assert-True {
    param([bool]$Condition, [string]$Message)
    if (-not $Condition) { throw "Assertion failed: $Message" }
}

function Assert-Equal {
    param($Expected, $Actual, [string]$Message)
    if ($Expected -ne $Actual) {
        throw "Assertion failed: $Message. Expected '$Expected', got '$Actual'."
    }
}

function New-FixtureProject {
    param([string]$Name)
    $path = Join-Path $TestRoot $Name
    New-Item -ItemType Directory -Force -Path $path | Out-Null
    return $path
}

function Get-TreeFingerprint {
    param([string]$Root)
    $rows = Get-ChildItem -Path $Root -Recurse -File |
        Where-Object { $_.FullName -notmatch '[\\/]\.sdp-backups[\\/]' } |
        Sort-Object FullName |
        ForEach-Object {
            $relative = $_.FullName.Substring($Root.Length).TrimStart([char[]]'\/')
            "$relative=$((Get-FileHash -Algorithm SHA256 -LiteralPath $_.FullName).Hash)"
        }
    return ($rows -join "`n")
}

try {
    New-Item -ItemType Directory -Force -Path $TestRoot | Out-Null

    # Preview must be side-effect free.
    $preview = New-FixtureProject 'preview-empty'
    & $Installer -ProjectRoot $preview -Preview | Out-Host
    Assert-Equal 0 (Get-ChildItem -LiteralPath $preview -Force | Measure-Object).Count 'preview created files'

    # Empty installation and idempotent repeat.
    $empty = New-FixtureProject 'empty'
    & $Installer -ProjectRoot $empty | Out-Host
    Assert-True (Test-Path (Join-Path $empty 'AGENTS.md')) 'AGENTS.md missing'
    Assert-True (Test-Path (Join-Path $empty 'SDP\SDP-project.manifest.yaml')) 'project manifest missing'
    Assert-True (Test-Path (Join-Path $empty 'SDP\RELEASE-NOTES.md')) 'release notes missing'
    Assert-True (Test-Path (Join-Path $empty 'SDP\Framework\installed-toolkit.manifest.yaml')) 'installed manifest missing'
    Assert-True (Test-Path (Join-Path $empty '.codex\skills\sdp-release\SKILL.md')) 'release skill missing'
    $beforeRepeat = Get-TreeFingerprint $empty
    & $Installer -ProjectRoot $empty | Out-Host
    $afterRepeat = Get-TreeFingerprint $empty
    Assert-Equal $beforeRepeat $afterRepeat 'repeat install was not idempotent'

    # Project-owned files survive; managed files require Force on same version.
    $projectNotes = Join-Path $empty 'SDP\RELEASE-NOTES.md'
    $projectManifest = Join-Path $empty 'SDP\SDP-project.manifest.yaml'
    $projectAgents = Join-Path $empty 'AGENTS-project.md'
    Set-Content -LiteralPath $projectNotes -Value 'CUSTOM NOTES' -NoNewline
    Set-Content -LiteralPath $projectManifest -Value 'CUSTOM MANIFEST' -NoNewline
    Set-Content -LiteralPath $projectAgents -Value 'CUSTOM AGENT RULES' -NoNewline
    $managedSkill = Join-Path $empty '.codex\skills\sdp-release\SKILL.md'
    Set-Content -LiteralPath $managedSkill -Value 'LOCAL MANAGED EDIT' -NoNewline
    & $Installer -ProjectRoot $empty | Out-Host
    Assert-Equal 'CUSTOM NOTES' (Get-Content -Raw -LiteralPath $projectNotes) 'release notes overwritten'
    Assert-Equal 'CUSTOM MANIFEST' (Get-Content -Raw -LiteralPath $projectManifest) 'project manifest overwritten'
    Assert-Equal 'CUSTOM AGENT RULES' (Get-Content -Raw -LiteralPath $projectAgents) 'AGENTS-project overwritten'
    Assert-Equal 'LOCAL MANAGED EDIT' (Get-Content -Raw -LiteralPath $managedSkill) 'managed edit changed without Force'

    $backup = Join-Path $TestRoot 'explicit-backup'
    & $Installer -ProjectRoot $empty -ForceManagedFiles -BackupRoot $backup | Out-Host
    Assert-True ((Get-Content -Raw -LiteralPath $managedSkill) -match 'skillId: sdp-release') 'Force did not restore managed skill'
    Assert-True (Test-Path (Join-Path $backup '.codex\skills\sdp-release\SKILL.md')) 'managed backup missing'
    Assert-Equal 'CUSTOM NOTES' (Get-Content -Raw -LiteralPath $projectNotes) 'Force overwrote project release notes'

    # Pre-versioning migration preserves project documents and old AGENTS rules.
    $legacy = New-FixtureProject 'legacy'
    New-Item -ItemType Directory -Force -Path (Join-Path $legacy 'SDP\03--Requirements') | Out-Null
    Set-Content -LiteralPath (Join-Path $legacy 'AGENTS.md') -Value 'OLD LOCAL AGENTS' -NoNewline
    Set-Content -LiteralPath (Join-Path $legacy 'SDP\03--Requirements\requirements.md') -Value 'LEGACY REQUIREMENTS' -NoNewline
    & $Installer -ProjectRoot $legacy | Out-Host
    Assert-Equal 'OLD LOCAL AGENTS' (Get-Content -Raw -LiteralPath (Join-Path $legacy 'AGENTS-project.md')) 'old AGENTS not migrated'
    Assert-Equal 'LEGACY REQUIREMENTS' (Get-Content -Raw -LiteralPath (Join-Path $legacy 'SDP\03--Requirements\requirements.md')) 'legacy document overwritten'

    # Unsupported schemas stop before mutation.
    $unsupported = New-FixtureProject 'unsupported'
    New-Item -ItemType Directory -Force -Path (Join-Path $unsupported 'SDP\Framework') | Out-Null
    Set-Content -LiteralPath (Join-Path $unsupported 'marker.txt') -Value 'UNCHANGED' -NoNewline
    Set-Content -LiteralPath (Join-Path $unsupported 'SDP\Framework\installed-toolkit.manifest.yaml') -Value "schemaVersion: `"9.0`"`ntoolkitVersion: `"9.0.0`"`n" -NoNewline
    $failed = $false
    try { & $Installer -ProjectRoot $unsupported | Out-Host } catch { $failed = $true }
    Assert-True $failed 'unsupported schema did not fail'
    Assert-Equal 'UNCHANGED' (Get-Content -Raw -LiteralPath (Join-Path $unsupported 'marker.txt')) 'unsupported fixture mutated'
    Assert-True (-not (Test-Path (Join-Path $unsupported 'AGENTS.md'))) 'unsupported install wrote AGENTS.md'

    # A sibling named SDP-Analyzer is valid; a child of the Toolkit repo is not.
    $sibling = Join-Path (Split-Path -Parent $RepositoryRoot) ("SDP-Analyzer-fixture-" + [guid]::NewGuid().ToString('N'))
    try {
        New-Item -ItemType Directory -Force -Path $sibling | Out-Null
        & $Installer -ProjectRoot $sibling -Preview | Out-Host
    } finally {
        Remove-Item -LiteralPath $sibling -Recurse -Force -ErrorAction SilentlyContinue
    }

    $inside = Join-Path $RepositoryRoot ("installer-child-fixture-" + [guid]::NewGuid().ToString('N'))
    try {
        New-Item -ItemType Directory -Force -Path $inside | Out-Null
        $insideFailed = $false
        try { & $Installer -ProjectRoot $inside -Preview | Out-Host } catch { $insideFailed = $true }
        Assert-True $insideFailed 'project inside Toolkit repository was accepted'
    } finally {
        Remove-Item -LiteralPath $inside -Recurse -Force -ErrorAction SilentlyContinue
    }

    Write-Host 'Installer fixture tests passed.'
} finally {
    Remove-Item -LiteralPath $TestRoot -Recurse -Force -ErrorAction SilentlyContinue
}
