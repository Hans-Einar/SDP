[CmdletBinding()]
param()

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$RepositoryRoot = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path
$Installer = Join-Path $RepositoryRoot 'Toolkit\scripts\Install-SDP.ps1'
$InstallManifestPath = Join-Path $RepositoryRoot 'Toolkit\SDP-install.manifest.json'
$PlanSchema = Join-Path $RepositoryRoot 'Toolkit\schemas\SDP-install-plan.schema.json'
$ProjectValidator = Join-Path $RepositoryRoot 'Toolkit\scripts\validate_sdp.py'
$TestRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("sdp-installer-tests-" + [guid]::NewGuid().ToString('N'))
$IsWindowsPlatform = [System.IO.Path]::DirectorySeparatorChar -eq '\'

if ($IsWindowsPlatform -and (-not ('Sdp.Tests.NativePath' -as [type]))) {
    Add-Type -TypeDefinition @'
using System.Runtime.InteropServices;
using System.Text;

namespace Sdp.Tests {
    public static class NativePath {
        [DllImport("kernel32.dll", CharSet = CharSet.Unicode, SetLastError = true,
            EntryPoint = "GetShortPathNameW")]
        public static extern uint GetShortPathName(
            string longPath,
            StringBuilder shortPath,
            uint bufferLength
        );
    }
}
'@
}

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

function Assert-Contains {
    param([object[]]$Values, $Actual, [string]$Message)
    if ($Values -notcontains $Actual) {
        throw "Assertion failed: $Message. Unsupported value '$Actual'."
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
    $rows = Get-ChildItem -LiteralPath $Root -Recurse -Force |
        Sort-Object FullName |
        ForEach-Object {
            $relative = $_.FullName.Substring($Root.Length).TrimStart([char[]]'\/')
            if ($_.PSIsContainer) {
                "D:$relative"
            } else {
                "F:$relative=$((Get-FileHash -Algorithm SHA256 -LiteralPath $_.FullName).Hash)"
            }
        }
    return ($rows -join "`n")
}

function Get-RelativeFiles {
    param([string]$Root)
    return @(
        Get-ChildItem -LiteralPath $Root -Recurse -Force -File |
            ForEach-Object {
                $_.FullName.Substring($Root.Length).TrimStart([char[]]'\/').Replace('\', '/')
            } |
            Sort-Object
    )
}

function Write-Utf8File {
    param([string]$Path, [string]$Content)
    $parent = Split-Path -Parent $Path
    if (-not [string]::IsNullOrWhiteSpace($parent)) {
        New-Item -ItemType Directory -Force -Path $parent | Out-Null
    }
    [System.IO.File]::WriteAllText($Path, $Content, [System.Text.UTF8Encoding]::new($false))
}

function Invoke-PlanJson {
    param(
        [string]$Target,
        [string]$InstallerPath = $Installer,
        [switch]$Initialize,
        [switch]$Force,
        [string]$BackupRoot
    )
    $parameters = @{
        ProjectRoot = $Target
        PlanJson = $true
    }
    if ($Initialize) { $parameters.InitializeProjectStructure = $true }
    if ($Force) { $parameters.ForceManagedFiles = $true }
    if (-not [string]::IsNullOrWhiteSpace($BackupRoot)) { $parameters.BackupRoot = $BackupRoot }
    $output = & $InstallerPath @parameters
    return [string]::Join("`n", [string[]]@($output))
}

function Get-ShortPathAlias {
    param([string]$Path)

    if (-not $IsWindowsPlatform) { return $null }
    $buffer = [System.Text.StringBuilder]::new(32768)
    $length = [Sdp.Tests.NativePath]::GetShortPathName(
        $Path,
        $buffer,
        [uint32]$buffer.Capacity
    )
    if (($length -eq 0) -or ($length -ge $buffer.Capacity)) { return $null }
    $shortPath = $buffer.ToString()
    if ($shortPath.Equals($Path, [System.StringComparison]::OrdinalIgnoreCase)) { return $null }
    return $shortPath
}

function Get-LoopbackAdministrativeUncPath {
    param([string]$Path)

    if (-not $IsWindowsPlatform) { return $null }
    $full = [System.IO.Path]::GetFullPath($Path)
    $root = [System.IO.Path]::GetPathRoot($full)
    if ($root -cnotmatch '^([A-Za-z]):\\$') { return $null }
    return '\\localhost\' + $Matches[1] + '$' + $full.Substring(2)
}

function Assert-PlanConforms {
    param([string]$PlanJsonContent, [string]$Label)
    $planFile = Join-Path $TestRoot ("plan-" + [guid]::NewGuid().ToString('N') + '.json')
    try {
        Write-Utf8File $planFile $PlanJsonContent
        $pythonCode = @'
import json
import sys
from jsonschema import Draft202012Validator

with open(sys.argv[1], encoding='utf-8') as stream:
    schema = json.load(stream)
with open(sys.argv[2], encoding='utf-8') as stream:
    plan = json.load(stream)
errors = sorted(
    Draft202012Validator(schema).iter_errors(plan),
    key=lambda error: [str(part) for part in error.path],
)
if errors:
    raise SystemExit('\n'.join(error.message for error in errors))
'@
        $validationOutput = & python -c $pythonCode $PlanSchema $planFile 2>&1
        if ($LASTEXITCODE -ne 0) {
            throw "Assertion failed: $Label plan schema validation failed: $validationOutput"
        }
    } finally {
        Remove-Item -LiteralPath $planFile -Force -ErrorAction SilentlyContinue
    }
}

$PlanReasonRules = @{
    'missing-target' = @('create', $true)
    'missing-generated-target' = @('generate', $true)
    'content-matches' = @('unchanged', $false)
    'missing-only-content' = @('preserve', $false)
    'managed-content-differs' = @('preserve', $false)
    'backup-before-replace' = @('backup', $true)
    'refresh-managed-content' = @('replace', $true)
    'refresh-generated-content' = @('generate', $true)
    'migrate-existing-agents' = @('migrate', $true)
    'preserve-existing-agents-conflict' = @('migrate', $true)
    'malformed-project-manifest' = @('block', $false)
    'unsupported-project-schema' = @('block', $false)
    'malformed-installed-manifest' = @('block', $false)
    'unsupported-installed-schema' = @('block', $false)
    'downgrade-blocked' = @('block', $false)
}

function Assert-PlanReasonSemantics {
    param($Plan, [string]$Label)

    $actions = @($Plan.actions)
    for ($index = 0; $index -lt $actions.Count; $index++) {
        $action = $actions[$index]
        Assert-Equal ($index + 1) ([int]$action.sequence) "$Label sequence"
        $reason = [string]$action.reason
        Assert-True ($PlanReasonRules.ContainsKey($reason)) "$Label unknown reason $reason"
        Assert-Equal ([string]$PlanReasonRules[$reason][0]) ([string]$action.action) "$Label action for $reason"
        Assert-Equal ([bool]$PlanReasonRules[$reason][1]) ([bool]$action.mutatesTarget) "$Label mutation for $reason"
        Assert-Equal ([string]$Plan.toolkitVersion) ([string]$action.newToolkitVersion) "$Label new version"
        if ($null -eq $Plan.installedToolkitVersion) {
            Assert-True ($null -eq $action.oldToolkitVersion) "$Label unexpected old version"
        } else {
            Assert-Equal ([string]$Plan.installedToolkitVersion) ([string]$action.oldToolkitVersion) "$Label old version"
        }
    }
}

function Assert-ManifestPreflightBlock {
    param(
        [string]$Name,
        [string]$RelativePath,
        [string]$Content,
        [string]$ExpectedReason
    )

    $target = New-FixtureProject $Name
    Write-Utf8File (Join-Path $target ($RelativePath.Replace('/', '\'))) $Content
    Write-Utf8File (Join-Path $target 'marker.txt') 'UNCHANGED'
    $before = Get-TreeFingerprint $target
    $planJson = Invoke-PlanJson $target
    Assert-Equal $before (Get-TreeFingerprint $target) "$Name PlanJson mutated target"
    Assert-PlanConforms $planJson $Name
    $plan = $planJson | ConvertFrom-Json
    Assert-True (-not [bool]$plan.canApply) "$Name plan was applicable"
    Assert-Equal 1 @($plan.actions).Count "$Name block action count"
    Assert-Equal 'block' ([string]$plan.actions[0].action) "$Name action"
    Assert-Equal $ExpectedReason ([string]$plan.actions[0].reason) "$Name reason"
    Assert-PlanReasonSemantics $plan $Name
    $failed = $false
    try { & $Installer -ProjectRoot $target | Out-Null } catch { $failed = $true }
    Assert-True $failed "$Name apply did not fail"
    Assert-Equal $before (Get-TreeFingerprint $target) "$Name apply mutated target"
}

function Assert-InvalidArchiveContract {
    param(
        [string]$Name,
        $Contract,
        [string]$ArchiveManifest,
        [string]$OriginalManifest,
        [string]$ArchiveInstaller
    )

    try {
        Write-Utf8File $ArchiveManifest ($Contract | ConvertTo-Json -Depth 60)
        $target = New-FixtureProject $Name
        Write-Utf8File (Join-Path $target 'marker.txt') 'UNCHANGED'
        $before = Get-TreeFingerprint $target
        $failed = $false
        try {
            Invoke-PlanJson $target -InstallerPath $ArchiveInstaller | Out-Null
        } catch {
            $failed = $true
        }
        Assert-True $failed "$Name invalid contract passed PowerShell preflight"
        Assert-Equal $before (Get-TreeFingerprint $target) "$Name invalid contract mutated target"
    } finally {
        Write-Utf8File $ArchiveManifest $OriginalManifest
    }
}

function Assert-InstallerRejectedWithoutMutation {
    param(
        [string]$Name,
        [string]$Target,
        [string]$InstallerPath,
        [string]$BackupRoot,
        [string[]]$WatchRoots,
        [string]$ExpectedErrorPattern
    )

    $before = @{}
    foreach ($root in $WatchRoots) {
        $before[$root] = Get-TreeFingerprint $root
    }

    $planFailed = $false
    $planError = $null
    try {
        Invoke-PlanJson `
            $Target `
            -InstallerPath $InstallerPath `
            -BackupRoot $BackupRoot | Out-Null
    } catch {
        $planFailed = $true
        $planError = $_.Exception.Message
    }
    Assert-True $planFailed "$Name PlanJson was not blocked"
    if (-not [string]::IsNullOrWhiteSpace($ExpectedErrorPattern)) {
        Assert-True `
            ($planError -match $ExpectedErrorPattern) `
            "$Name PlanJson failed for an unexpected reason: $planError"
    }
    foreach ($root in $WatchRoots) {
        Assert-Equal $before[$root] (Get-TreeFingerprint $root) "$Name PlanJson mutated $root"
    }

    $parameters = @{ ProjectRoot = $Target }
    if (-not [string]::IsNullOrWhiteSpace($BackupRoot)) {
        $parameters.BackupRoot = $BackupRoot
    }
    $applyFailed = $false
    $applyError = $null
    try { & $InstallerPath @parameters | Out-Null } catch {
        $applyFailed = $true
        $applyError = $_.Exception.Message
    }
    Assert-True $applyFailed "$Name apply was not blocked"
    if (-not [string]::IsNullOrWhiteSpace($ExpectedErrorPattern)) {
        Assert-True `
            ($applyError -match $ExpectedErrorPattern) `
            "$Name apply failed for an unexpected reason: $applyError"
    }
    foreach ($root in $WatchRoots) {
        Assert-Equal $before[$root] (Get-TreeFingerprint $root) "$Name apply mutated $root"
    }
}

function New-DirectoryLink {
    param([string]$Path, [string]$Target)

    if ([System.IO.Path]::DirectorySeparatorChar -eq '\') {
        New-Item -ItemType Junction -Path $Path -Target $Target | Out-Null
    } else {
        New-Item -ItemType SymbolicLink -Path $Path -Target $Target | Out-Null
    }
}

function Remove-DirectoryLink {
    param([string]$Path)

    if ([System.IO.Directory]::Exists($Path)) {
        [System.IO.Directory]::Delete($Path, $false)
    }
}

function Try-NewFileSymbolicLink {
    param([string]$Path, [string]$Target)

    try {
        New-Item -ItemType SymbolicLink -Path $Path -Target $Target -ErrorAction Stop | Out-Null
        return $true
    } catch {
        return $false
    }
}

try {
    New-Item -ItemType Directory -Force -Path $TestRoot | Out-Null

    # The portable manifest is the only install inventory and every policy is explicit.
    $contract = Get-Content -Raw -LiteralPath $InstallManifestPath | ConvertFrom-Json
    Assert-Equal '1.0' ([string]$contract.schemaVersion) 'installation contract schema'
    Assert-Equal '0.2.0' ([string]$contract.toolkitVersion) 'installation contract Toolkit version'
    $entryIds = @{}
    $destinations = @{}
    foreach ($entry in @($contract.entries)) {
        $entryId = [string]$entry.id
        $destination = [string]$entry.destination
        Assert-True (-not $entryIds.ContainsKey($entryId)) "duplicate entry ID $entryId"
        Assert-True (-not $destinations.ContainsKey($destination.ToLowerInvariant())) "duplicate destination $destination"
        $entryIds[$entryId] = $true
        $destinations[$destination.ToLowerInvariant()] = $true
        Assert-True (
            $destination -notmatch '(^/|^[A-Za-z]:|\\|//|(^|/)\.\.?(/|$)|/$)'
        ) "non-portable destination $destination"
        Assert-Contains @('copied', 'generated') ([string]$entry.kind) "$entryId kind"
        Assert-Contains @('toolkit-managed', 'project-owned') ([string]$entry.ownership) "$entryId ownership"
        Assert-Contains @('default', 'initialize-only') ([string]$entry.selectionPolicy) "$entryId selectionPolicy"
        Assert-Contains @('always', 'missing-only') ([string]$entry.installPolicy) "$entryId installPolicy"
        Assert-Contains @('always', 'upgrade-or-force', 'never') ([string]$entry.refreshPolicy) "$entryId refreshPolicy"
        Assert-Contains @('before-replace', 'migration-aware', 'none') ([string]$entry.backupPolicy) "$entryId backupPolicy"
        Assert-Contains @('replace-managed', 'preserve') ([string]$entry.forcePolicy) "$entryId forcePolicy"
        Assert-Contains @('none', 'preserve-existing-agents') ([string]$entry.migrationPolicy) "$entryId migrationPolicy"
        if ([string]$entry.kind -eq 'copied') {
            $source = [string]$entry.source
            Assert-True (
                $source -notmatch '(^/|^[A-Za-z]:|\\|//|(^|/)\.\.?(/|$)|/$)'
            ) "non-portable source $source"
            $sourcePath = Join-Path $RepositoryRoot ($source.Replace('/', '\'))
            Assert-True (Test-Path -LiteralPath $sourcePath -PathType Leaf) "missing contract source $source"
        }
        if ([string]$entry.ownership -eq 'project-owned') {
            Assert-Equal 'missing-only' ([string]$entry.installPolicy) "$entryId project-owned install policy"
            Assert-Equal 'never' ([string]$entry.refreshPolicy) "$entryId project-owned refresh policy"
            Assert-Equal 'none' ([string]$entry.backupPolicy) "$entryId project-owned backup policy"
            Assert-Equal 'preserve' ([string]$entry.forcePolicy) "$entryId project-owned force policy"
        }
    }

    # JSON plans are deterministic, schema-valid and strictly read-only.
    $planOnly = New-FixtureProject 'plan-json-empty'
    $beforePlan = Get-TreeFingerprint $planOnly
    $defaultPlanJson = Invoke-PlanJson $planOnly
    $repeatPlanJson = Invoke-PlanJson $planOnly
    $afterPlan = Get-TreeFingerprint $planOnly
    Assert-Equal $beforePlan $afterPlan 'PlanJson mutated the target'
    Assert-Equal $defaultPlanJson $repeatPlanJson 'PlanJson was not deterministic'
    Assert-PlanConforms $defaultPlanJson 'default'
    $defaultPlan = $defaultPlanJson | ConvertFrom-Json
    Assert-True ([bool]$defaultPlan.canApply) 'default plan unexpectedly blocked'
    Assert-PlanReasonSemantics $defaultPlan 'default plan'
    $defaultActions = @($defaultPlan.actions)
    Assert-True ($defaultActions.Count -gt 0) 'default plan has no actions'
    for ($index = 0; $index -lt $defaultActions.Count; $index++) {
        $action = $defaultActions[$index]
        Assert-Equal ($index + 1) ([int]$action.sequence) 'plan action sequence'
        Assert-True (-not [string]::IsNullOrWhiteSpace([string]$action.entryId)) 'plan action missing entry ID'
        Assert-True (-not [string]::IsNullOrWhiteSpace([string]$action.destination)) 'plan action missing destination'
        $hasSource = -not [string]::IsNullOrWhiteSpace([string]$action.source)
        $hasGenerator = -not [string]::IsNullOrWhiteSpace([string]$action.generator)
        Assert-True ($hasSource -xor $hasGenerator) "plan action $($action.entryId) lacks a single source/generator"
        Assert-True ($null -eq $action.targetSource) "plan action $($action.entryId) unexpectedly has targetSource"
    }
    Assert-True (
        @($defaultActions | Where-Object { $_.entryId -eq 'project-sprints-readme' }).Count -eq 0
    ) 'default plan selected initialize-only content'

    $initializePlanJson = Invoke-PlanJson $planOnly -Initialize
    Assert-PlanConforms $initializePlanJson 'initialize'
    $initializePlan = $initializePlanJson | ConvertFrom-Json
    Assert-PlanReasonSemantics $initializePlan 'initialize plan'
    Assert-True (
        @($initializePlan.actions | Where-Object { $_.entryId -eq 'project-sprints-readme' }).Count -eq 1
    ) 'initialize plan omitted neutral structure content'
    Assert-True (
        @($initializePlan.actions | Where-Object { $_.destination -eq 'SDP/Releases/REL-0.2.0.yaml' }).Count -eq 0
    ) 'initialize plan included Toolkit release state'

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
    Assert-True (Test-Path (Join-Path $empty 'SDP\Traceability\Ledger.ndjson')) 'empty ledger missing'
    Assert-True (Test-Path (Join-Path $empty '.codex\skills\sdp-release\SKILL.md')) 'release skill missing'
    Assert-True (-not (Test-Path (Join-Path $empty 'SDP\01--Mandate\mandate.md'))) 'normal install initialized project structure'
    Assert-True (-not (Test-Path (Join-Path $empty 'SDP\Releases\REL-0.2.0.yaml'))) 'normal install copied Toolkit release state'
    $expectedDefaultFiles = @(
        @($contract.entries) |
            Where-Object { [string]$_.selectionPolicy -eq 'default' } |
            ForEach-Object { [string]$_.destination } |
            Sort-Object
    )
    $actualDefaultFiles = Get-RelativeFiles $empty
    Assert-Equal (
        $expectedDefaultFiles -join "`n"
    ) (
        $actualDefaultFiles -join "`n"
    ) 'installer wrote files outside the canonical default inventory'
    $beforeRepeat = Get-TreeFingerprint $empty
    & $Installer -ProjectRoot $empty | Out-Host
    $afterRepeat = Get-TreeFingerprint $empty
    Assert-Equal $beforeRepeat $afterRepeat 'repeat install was not idempotent'

    # Project-owned files survive; managed files require Force on same version.
    $projectNotes = Join-Path $empty 'SDP\RELEASE-NOTES.md'
    $projectManifest = Join-Path $empty 'SDP\SDP-project.manifest.yaml'
    $projectAgents = Join-Path $empty 'AGENTS-project.md'
    $customManifest = @"
schemaVersion: "1.0"
project:
  name: Custom preserved fixture
  capabilities: []
installed:
  manifestPath: Framework/installed-toolkit.manifest.yaml
release:
  currentVersion: "0.0.0"
  nextTargetVersion: null
  state: unreleased
  latestTag: null
  latestCommit: null
development:
  sprintId: null
  refactorId: null
  iterationId: null
  sliceId: null
  fixId: null
  revision: null
migration:
  pendingWarnings: []
"@
    Set-Content -LiteralPath $projectNotes -Value 'CUSTOM NOTES' -NoNewline
    Set-Content -LiteralPath $projectManifest -Value $customManifest -NoNewline
    Set-Content -LiteralPath $projectAgents -Value 'CUSTOM AGENT RULES' -NoNewline
    $managedSkill = Join-Path $empty '.codex\skills\sdp-release\SKILL.md'
    Set-Content -LiteralPath $managedSkill -Value 'LOCAL MANAGED EDIT' -NoNewline
    & $Installer -ProjectRoot $empty | Out-Host
    Assert-Equal 'CUSTOM NOTES' (Get-Content -Raw -LiteralPath $projectNotes) 'release notes overwritten'
    Assert-Equal $customManifest (Get-Content -Raw -LiteralPath $projectManifest) 'project manifest overwritten'
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
    $legacyPlanJson = Invoke-PlanJson $legacy
    Assert-PlanConforms $legacyPlanJson 'AGENTS migration'
    $legacyPlan = $legacyPlanJson | ConvertFrom-Json
    Assert-PlanReasonSemantics $legacyPlan 'AGENTS migration'
    $legacyMigration = @($legacyPlan.actions | Where-Object { $_.reason -eq 'migrate-existing-agents' })
    Assert-Equal 1 $legacyMigration.Count 'AGENTS migration action count'
    Assert-Equal 'migrate' ([string]$legacyMigration[0].action) 'AGENTS migration action'
    Assert-Equal 'AGENTS.md' ([string]$legacyMigration[0].targetSource) 'AGENTS migration targetSource'
    Assert-Equal 'AGENTS-project.md' ([string]$legacyMigration[0].destination) 'AGENTS migration destination'
    & $Installer -ProjectRoot $legacy | Out-Host
    Assert-Equal 'OLD LOCAL AGENTS' (Get-Content -Raw -LiteralPath (Join-Path $legacy 'AGENTS-project.md')) 'old AGENTS not migrated'
    Assert-Equal 'LEGACY REQUIREMENTS' (Get-Content -Raw -LiteralPath (Join-Path $legacy 'SDP\03--Requirements\requirements.md')) 'legacy document overwritten'

    # If AGENTS-project.md already exists, migration preserves both rule sets.
    $legacyConflict = New-FixtureProject 'legacy-agents-conflict'
    Write-Utf8File (Join-Path $legacyConflict 'AGENTS.md') 'OLD CONFLICTING AGENTS'
    Write-Utf8File (Join-Path $legacyConflict 'AGENTS-project.md') 'EXISTING PROJECT AGENTS'
    $beforeConflictPlan = Get-TreeFingerprint $legacyConflict
    $conflictPlanJson = Invoke-PlanJson $legacyConflict
    Assert-Equal $conflictPlanJson (Invoke-PlanJson $legacyConflict) 'AGENTS conflict plan was not deterministic'
    Assert-PlanConforms $conflictPlanJson 'AGENTS migration conflict'
    Assert-Equal $beforeConflictPlan (Get-TreeFingerprint $legacyConflict) 'AGENTS migration plan mutated target'
    $conflictPlan = $conflictPlanJson | ConvertFrom-Json
    Assert-PlanReasonSemantics $conflictPlan 'AGENTS migration conflict'
    $conflictMigrations = @(
        $conflictPlan.actions | Where-Object { $_.reason -eq 'preserve-existing-agents-conflict' }
    )
    Assert-Equal 1 @(
        $conflictMigrations
    ).Count 'AGENTS conflict plan preservation action count'
    Assert-Equal 'migrate' ([string]$conflictMigrations[0].action) 'AGENTS conflict action'
    Assert-Equal 'AGENTS.md' ([string]$conflictMigrations[0].targetSource) 'AGENTS conflict targetSource'
    Assert-True (
        [string]$conflictMigrations[0].destination -cmatch '^AGENTS-project\.migration-sha256-[0-9a-f]{64}\.md$'
    ) 'AGENTS conflict deterministic destination'
    $plannedMigrationRelative = [string]$conflictMigrations[0].destination
    & $Installer -ProjectRoot $legacyConflict | Out-Host
    Assert-Equal 'EXISTING PROJECT AGENTS' (Get-Content -Raw -LiteralPath (Join-Path $legacyConflict 'AGENTS-project.md')) 'existing AGENTS-project overwritten during migration'
    $migrationFiles = @(Get-ChildItem -LiteralPath $legacyConflict -File -Filter 'AGENTS-project.migration-sha256-*.md')
    Assert-Equal 1 $migrationFiles.Count 'AGENTS conflict migration file count'
    Assert-Equal (Join-Path $legacyConflict $plannedMigrationRelative) $migrationFiles[0].FullName 'AGENTS apply differed from planned migration destination'
    Assert-Equal 'OLD CONFLICTING AGENTS' (Get-Content -Raw -LiteralPath $migrationFiles[0].FullName) 'AGENTS conflict content not preserved'
    Assert-Equal (
        (Get-Content -Raw -LiteralPath (Join-Path $RepositoryRoot 'Toolkit\payload\project-root\AGENTS.md.template'))
    ) (
        Get-Content -Raw -LiteralPath (Join-Path $legacyConflict 'AGENTS.md')
    ) 'managed AGENTS not installed after conflict migration'

    # gh-sdp-like bootstrap: preserve project state and install only neutral structure.
    $bootstrap = New-FixtureProject 'gh-sdp-bootstrap'
    New-Item -ItemType Directory -Force -Path (Join-Path $bootstrap '.git') | Out-Null
    $bootstrapManifest = @'
schemaVersion: "1.0"
project:
  name: Bootstrap fixture
  capabilities: []
installed:
  manifestPath: Framework/installed-toolkit.manifest.yaml
release:
  currentVersion: "0.0.0"
  nextTargetVersion: null
  state: unreleased
  latestTag: null
  latestCommit: null
development:
  sprintId: null
  refactorId: null
  iterationId: null
  sliceId: null
  fixId: null
  revision: null
migration:
  pendingWarnings: []
'@
    $bootstrapMandate = @'
# Mandate

Status: active

## Purpose

CUSTOM BOOTSTRAP MANDATE
'@
    $bootstrapCurrent = @'
project:
  name: Bootstrap fixture
  status: active-development
release:
  previousVersion: null
  activeReleaseId: null
  targetVersion: null
  state: unreleased
active:
  sprint: null
  refactor: null
  iteration: null
  slice: null
  fix: null
  revision: null
# CUSTOM CURRENT INDEX
'@
    $bootstrapRelations = @'
requirements: {}
designs: {}
sprints: {}
refactors: {}
iterations: {}
slices: {}
fixes: {}
reviews: {}
verification: {}
migrations: {}
releases: {}
# CUSTOM RELATIONS
'@
    $bootstrapNotes = @'
# Release Notes

## [Unreleased]

Release-Date: unreleased

### Added

- CUSTOM BOOTSTRAP NOTES
'@
    $bootstrapLedger = '{"schemaVersion":"1.0","eventId":"EVT-X-0001","eventType":"x-bootstrap:created","occurredAt":"2026-07-13T20:00:00Z","actor":"fixture","commit":null,"payload":{"source":"fixture"}}' + "`n"
    Write-Utf8File (Join-Path $bootstrap 'SDP\SDP-project.manifest.yaml') $bootstrapManifest
    Write-Utf8File (Join-Path $bootstrap 'SDP\01--Mandate\mandate.md') $bootstrapMandate
    Write-Utf8File (Join-Path $bootstrap 'SDP\Traceability\CurrentIndex.yaml') $bootstrapCurrent
    Write-Utf8File (Join-Path $bootstrap 'SDP\Traceability\Relations.yaml') $bootstrapRelations
    Write-Utf8File (Join-Path $bootstrap 'SDP\Traceability\Ledger.ndjson') $bootstrapLedger
    Write-Utf8File (Join-Path $bootstrap 'SDP\RELEASE-NOTES.md') $bootstrapNotes
    & $Installer -ProjectRoot $bootstrap -InitializeProjectStructure | Out-Host
    Assert-Equal $bootstrapManifest (Get-Content -Raw -LiteralPath (Join-Path $bootstrap 'SDP\SDP-project.manifest.yaml')) 'bootstrap manifest overwritten'
    Assert-Equal $bootstrapMandate (Get-Content -Raw -LiteralPath (Join-Path $bootstrap 'SDP\01--Mandate\mandate.md')) 'bootstrap Mandate overwritten'
    Assert-Equal $bootstrapCurrent (Get-Content -Raw -LiteralPath (Join-Path $bootstrap 'SDP\Traceability\CurrentIndex.yaml')) 'bootstrap CurrentIndex overwritten'
    Assert-Equal $bootstrapRelations (Get-Content -Raw -LiteralPath (Join-Path $bootstrap 'SDP\Traceability\Relations.yaml')) 'bootstrap Relations overwritten'
    Assert-Equal $bootstrapLedger (Get-Content -Raw -LiteralPath (Join-Path $bootstrap 'SDP\Traceability\Ledger.ndjson')) 'bootstrap Ledger overwritten'
    Assert-Equal $bootstrapNotes (Get-Content -Raw -LiteralPath (Join-Path $bootstrap 'SDP\RELEASE-NOTES.md')) 'bootstrap release notes overwritten'
    $gitDirectories = @(Get-ChildItem -LiteralPath $bootstrap -Recurse -Force -Directory | Where-Object { $_.Name -eq '.git' })
    Assert-Equal 1 $gitDirectories.Count 'bootstrap .git directory count'
    Assert-Equal (Join-Path $bootstrap '.git') $gitDirectories[0].FullName 'bootstrap .git location'
    Assert-True (-not (Test-Path (Join-Path $bootstrap 'SDP\Releases\REL-0.2.0.yaml'))) 'Toolkit release record contaminated bootstrap'
    Assert-Equal 'README.md' ((Get-ChildItem -LiteralPath (Join-Path $bootstrap 'SDP\Sprints') -File).Name -join ',') 'active Sprint history copied'
    Assert-Equal 'README.md' ((Get-ChildItem -LiteralPath (Join-Path $bootstrap 'SDP\Releases') -File).Name -join ',') 'release history copied'
    Assert-Equal 'README.md' ((Get-ChildItem -LiteralPath (Join-Path $bootstrap 'SDP\CodeReview') -File).Name -join ',') 'review history copied'
    Assert-Equal 'README.md' ((Get-ChildItem -LiteralPath (Join-Path $bootstrap 'SDP\Verification') -File).Name -join ',') 'verification history copied'
    $beforeInitializePreview = Get-TreeFingerprint $bootstrap
    $initializePreviewOutput = (& $Installer -ProjectRoot $bootstrap -InitializeProjectStructure -Preview 6>&1 | Out-String)
    $afterInitializePreview = Get-TreeFingerprint $bootstrap
    Assert-Equal $beforeInitializePreview $afterInitializePreview 'repeated initialize preview mutated project'
    Assert-True ($initializePreviewOutput -notmatch 'REL-0\.2\.0') 'repeated initialize preview proposed Toolkit release state'
    $bootstrapPlanJson = Invoke-PlanJson $bootstrap -Initialize
    Assert-PlanConforms $bootstrapPlanJson 'repeated initialize'
    Assert-True ($bootstrapPlanJson -notmatch 'REL-0\.2\.0') 'repeated initialize plan proposed Toolkit release state'
    & python $ProjectValidator --mode project --project-root $bootstrap | Out-Host
    Assert-Equal 0 $LASTEXITCODE 'initialized consuming-project validation'

    # Strict YAML root preflight rejects nested shadows, duplicates and malformed
    # scalar lookalikes for both project and installed manifests.
    Assert-ManifestPreflightBlock `
        'project-nested-schema-shadow' `
        'SDP/SDP-project.manifest.yaml' `
        "nested:`n  schemaVersion: `"1.0`"`nschemaVersion: `"9.0`"`n" `
        'malformed-project-manifest'
    Assert-ManifestPreflightBlock `
        'project-duplicate-schema' `
        'SDP/SDP-project.manifest.yaml' `
        "schemaVersion: `"1.0`"`nschemaVersion: `"1.0`"`n" `
        'malformed-project-manifest'
    Assert-ManifestPreflightBlock `
        'project-malformed-scalar' `
        'SDP/SDP-project.manifest.yaml' `
        "schemaVersion: `"1.0`" trailing-text`n" `
        'malformed-project-manifest'
    Assert-ManifestPreflightBlock `
        'installed-nested-version-shadow' `
        'SDP/Framework/installed-toolkit.manifest.yaml' `
        "nested:`n  toolkitVersion: `"0.2.0`"`nschemaVersion: `"1.0`"`ntoolkitVersion: `"0.3.0`"`n" `
        'malformed-installed-manifest'
    Assert-ManifestPreflightBlock `
        'installed-duplicate-schema' `
        'SDP/Framework/installed-toolkit.manifest.yaml' `
        "schemaVersion: `"1.0`"`nschemaVersion: `"1.0`"`ntoolkitVersion: `"0.2.0`"`n" `
        'malformed-installed-manifest'
    Assert-ManifestPreflightBlock `
        'installed-malformed-scalar' `
        'SDP/Framework/installed-toolkit.manifest.yaml' `
        "schemaVersion: `"1.0`"`ntoolkitVersion: [not-valid`n" `
        'malformed-installed-manifest'

    # Unsupported installed manifest schemas stop before mutation.
    $unsupported = New-FixtureProject 'unsupported'
    New-Item -ItemType Directory -Force -Path (Join-Path $unsupported 'SDP\Framework') | Out-Null
    Set-Content -LiteralPath (Join-Path $unsupported 'marker.txt') -Value 'UNCHANGED' -NoNewline
    Set-Content -LiteralPath (Join-Path $unsupported 'SDP\Framework\installed-toolkit.manifest.yaml') -Value "schemaVersion: `"9.0`"`ntoolkitVersion: `"9.0.0`"`n" -NoNewline
    $beforeBlockedPlan = Get-TreeFingerprint $unsupported
    $blockedPlanJson = Invoke-PlanJson $unsupported
    $afterBlockedPlan = Get-TreeFingerprint $unsupported
    Assert-Equal $beforeBlockedPlan $afterBlockedPlan 'blocked PlanJson mutated target'
    Assert-PlanConforms $blockedPlanJson 'blocked'
    $blockedPlan = $blockedPlanJson | ConvertFrom-Json
    Assert-PlanReasonSemantics $blockedPlan 'unsupported installed schema plan'
    Assert-True (-not [bool]$blockedPlan.canApply) 'unsupported schema plan was applicable'
    Assert-Equal 'block' ([string]@($blockedPlan.actions)[0].action) 'unsupported schema plan action'
    $failed = $false
    try { & $Installer -ProjectRoot $unsupported | Out-Host } catch { $failed = $true }
    Assert-True $failed 'unsupported schema did not fail'
    Assert-Equal 'UNCHANGED' (Get-Content -Raw -LiteralPath (Join-Path $unsupported 'marker.txt')) 'unsupported fixture mutated'
    Assert-True (-not (Test-Path (Join-Path $unsupported 'AGENTS.md'))) 'unsupported install wrote AGENTS.md'

    # Unsupported project manifest schemas also stop before mutation.
    $unsupportedProject = New-FixtureProject 'unsupported-project-manifest'
    New-Item -ItemType Directory -Force -Path (Join-Path $unsupportedProject 'SDP') | Out-Null
    Set-Content -LiteralPath (Join-Path $unsupportedProject 'SDP\SDP-project.manifest.yaml') -Value "schemaVersion: `"9.0`"`n" -NoNewline
    $projectSchemaFailed = $false
    try { & $Installer -ProjectRoot $unsupportedProject | Out-Host } catch { $projectSchemaFailed = $true }
    Assert-True $projectSchemaFailed 'unsupported project manifest schema did not fail'
    Assert-True (-not (Test-Path (Join-Path $unsupportedProject 'AGENTS.md'))) 'unsupported project manifest installation mutated project'

    # An older installer must never downgrade a newer installed Toolkit.
    $newer = New-FixtureProject 'newer-toolkit'
    New-Item -ItemType Directory -Force -Path (Join-Path $newer 'SDP\Framework') | Out-Null
    Set-Content -LiteralPath (Join-Path $newer 'SDP\Framework\installed-toolkit.manifest.yaml') -Value "schemaVersion: `"1.0`"`ntoolkitVersion: `"0.3.0`"`ntoolkitInstalledAt: `"2026-07-12T00:00:00Z`"`n" -NoNewline
    $downgradeFailed = $false
    try { & $Installer -ProjectRoot $newer | Out-Host } catch { $downgradeFailed = $true }
    Assert-True $downgradeFailed 'newer Toolkit installation was downgraded'
    Assert-True (-not (Test-Path (Join-Path $newer 'AGENTS.md'))) 'downgrade attempt mutated project'

    # SemVer 2.0 precedence distinguishes prereleases from finals and rejects
    # numeric prerelease identifiers with leading zeroes.
    $semverSeed = New-FixtureProject 'semver-seed'
    & $Installer -ProjectRoot $semverSeed | Out-Null
    $semverSeedContent = Get-Content -Raw -LiteralPath (Join-Path $semverSeed 'SDP\Framework\installed-toolkit.manifest.yaml')

    $prereleaseInstalled = New-FixtureProject 'installed-prerelease'
    $prereleaseContent = $semverSeedContent -replace '(?m)^toolkitVersion: "0\.2\.0"$', 'toolkitVersion: "0.2.0-rc.1"'
    Write-Utf8File (Join-Path $prereleaseInstalled 'SDP\Framework\installed-toolkit.manifest.yaml') $prereleaseContent
    $prereleasePlan = (Invoke-PlanJson $prereleaseInstalled) | ConvertFrom-Json
    Assert-True ([bool]$prereleasePlan.canApply) 'final Toolkit did not upgrade its prerelease'
    Assert-PlanReasonSemantics $prereleasePlan 'prerelease-to-final plan'

    $newerPrerelease = New-FixtureProject 'newer-prerelease'
    $newerPrereleaseContent = $semverSeedContent -replace '(?m)^toolkitVersion: "0\.2\.0"$', 'toolkitVersion: "0.2.1-alpha.1"'
    Write-Utf8File (Join-Path $newerPrerelease 'SDP\Framework\installed-toolkit.manifest.yaml') $newerPrereleaseContent
    $newerPrereleasePlan = (Invoke-PlanJson $newerPrerelease) | ConvertFrom-Json
    Assert-True (-not [bool]$newerPrereleasePlan.canApply) 'newer prerelease core did not block downgrade'
    Assert-Equal 'downgrade-blocked' ([string]$newerPrereleasePlan.actions[0].reason) 'newer prerelease block reason'

    $invalidPrerelease = New-FixtureProject 'invalid-prerelease'
    $invalidPrereleaseContent = $semverSeedContent -replace '(?m)^toolkitVersion: "0\.2\.0"$', 'toolkitVersion: "0.2.0-01"'
    Write-Utf8File (Join-Path $invalidPrerelease 'SDP\Framework\installed-toolkit.manifest.yaml') $invalidPrereleaseContent
    $invalidPrereleasePlan = (Invoke-PlanJson $invalidPrerelease) | ConvertFrom-Json
    Assert-True (-not [bool]$invalidPrereleasePlan.canApply) 'invalid numeric prerelease passed preflight'
    Assert-Equal 'malformed-installed-manifest' ([string]$invalidPrereleasePlan.actions[0].reason) 'invalid prerelease block reason'

    $invalidTimestamp = New-FixtureProject 'invalid-installed-timestamp'
    $invalidTimestampContent = $semverSeedContent -replace `
        '(?m)^toolkitInstalledAt: ".*"$', `
        'toolkitInstalledAt: "2026-07-14"'
    Write-Utf8File (Join-Path $invalidTimestamp 'SDP\Framework\installed-toolkit.manifest.yaml') $invalidTimestampContent
    $invalidTimestampPlan = (Invoke-PlanJson $invalidTimestamp) | ConvertFrom-Json
    Assert-True (-not [bool]$invalidTimestampPlan.canApply) 'non-RFC3339 installed timestamp passed preflight'
    Assert-Equal 'malformed-installed-manifest' ([string]$invalidTimestampPlan.actions[0].reason) 'invalid timestamp block reason'

    # Mapping order and quoting are serialization details: a semantically equal
    # installed manifest must not cause refresh churn or a backup.
    $semanticManifest = New-FixtureProject 'semantic-installed-manifest'
    & $Installer -ProjectRoot $semanticManifest | Out-Null
    $semanticPath = Join-Path $semanticManifest 'SDP\Framework\installed-toolkit.manifest.yaml'
    $canonicalInstalled = Get-Content -Raw -LiteralPath $semanticPath
    $canonicalLines = [regex]::Split($canonicalInstalled.TrimEnd("`r", "`n"), "\r?\n")
    $skillsStart = [Array]::IndexOf($canonicalLines, 'skills:')
    $capabilitiesStart = [Array]::IndexOf($canonicalLines, 'capabilities:')
    Assert-True (($skillsStart -gt 0) -and ($capabilitiesStart -gt $skillsStart)) 'canonical installed manifest sections'
    $rootByName = @{}
    foreach ($line in $canonicalLines[0..($skillsStart - 1)]) {
        if ($line -cmatch '^([A-Za-z]+):\s*(.+)$') { $rootByName[$Matches[1]] = $line }
    }
    $reorderedSkills = @($canonicalLines[($skillsStart + 1)..($capabilitiesStart - 1)])
    [array]::Reverse($reorderedSkills)
    $alternateLines = New-Object System.Collections.ArrayList
    foreach ($field in @(
        'toolkitVersion', 'schemaVersion', 'installerVersion', 'frameworkVersion',
        'agentsContractVersion', 'sourceCommit', 'toolkitInstalledAt'
    )) {
        [void]$alternateLines.Add(([string]$rootByName[$field]).Replace('"', "'"))
    }
    [void]$alternateLines.Add('capabilities:')
    foreach ($line in $canonicalLines[($capabilitiesStart + 1)..($canonicalLines.Count - 1)]) {
        [void]$alternateLines.Add($line)
    }
    [void]$alternateLines.Add('skills:')
    foreach ($line in $reorderedSkills) {
        [void]$alternateLines.Add($line.Replace('"', "'"))
    }
    $alternateInstalled = [string]::Join("`n", [string[]]$alternateLines) + "`n"
    Assert-True ($alternateInstalled -cne $canonicalInstalled) 'alternate installed serialization did not differ'
    Write-Utf8File $semanticPath $alternateInstalled
    $beforeSemanticPlan = Get-TreeFingerprint $semanticManifest
    $semanticPlanJson = Invoke-PlanJson $semanticManifest
    Assert-Equal $beforeSemanticPlan (Get-TreeFingerprint $semanticManifest) 'semantic installed plan mutated target'
    Assert-PlanConforms $semanticPlanJson 'semantic installed manifest'
    $semanticPlan = $semanticPlanJson | ConvertFrom-Json
    $installedAction = @($semanticPlan.actions | Where-Object { $_.entryId -eq 'generated-installed-toolkit-manifest' })
    Assert-Equal 1 $installedAction.Count 'semantic installed action count'
    Assert-Equal 'unchanged' ([string]$installedAction[0].action) 'semantic installed manifest caused refresh'
    & $Installer -ProjectRoot $semanticManifest | Out-Null
    Assert-Equal $beforeSemanticPlan (Get-TreeFingerprint $semanticManifest) 'semantic installed apply caused serialization churn'

    # A source-archive extraction with no .git installs successfully and records an unknown commit truthfully.
    $archiveRoot = Join-Path $TestRoot 'archive-source'
    New-Item -ItemType Directory -Force -Path $archiveRoot | Out-Null
    Get-ChildItem -LiteralPath $RepositoryRoot -Force |
        Where-Object { $_.Name -ne '.git' } |
        ForEach-Object {
            Copy-Item -LiteralPath $_.FullName -Destination $archiveRoot -Recurse -Force
        }
    Assert-Equal 0 @(
        Get-ChildItem -LiteralPath $archiveRoot -Recurse -Force |
            Where-Object { $_.Name -eq '.git' }
    ).Count 'archive fixture contains .git'
    $archiveInstaller = Join-Path $archiveRoot 'Toolkit\scripts\Install-SDP.ps1'
    $archiveTarget = New-FixtureProject 'archive-target'
    & $archiveInstaller -ProjectRoot $archiveTarget | Out-Host
    $archiveInstalledManifest = Get-Content -Raw -LiteralPath (Join-Path $archiveTarget 'SDP\Framework\installed-toolkit.manifest.yaml')
    Assert-True ($archiveInstalledManifest -match '(?m)^sourceCommit:\s*null\s*$') 'archive install claimed a source commit'
    Assert-True (-not (Test-Path (Join-Path $archiveTarget 'SDP\Releases\REL-0.2.0.yaml'))) 'archive install copied release state'
    & python $ProjectValidator --mode project --project-root $archiveTarget | Out-Host
    Assert-Equal 0 $LASTEXITCODE 'archive-installed project validation'

    # Root overlap is decided by stable filesystem identities, not path
    # spelling. Exact, ancestor and descendant forms are symmetric, and both
    # PlanJson and apply remain mutation-free when rejected.
    Assert-InstallerRejectedWithoutMutation `
        'exact local source/project overlap' `
        $archiveRoot `
        $archiveInstaller `
        '' `
        @($archiveRoot) `
        -ExpectedErrorPattern 'physically separate trees'
    Assert-InstallerRejectedWithoutMutation `
        'project below local source' `
        (Join-Path $archiveRoot 'Toolkit') `
        $archiveInstaller `
        '' `
        @($archiveRoot) `
        -ExpectedErrorPattern 'physically separate trees'

    $localBackupOverlapTarget = New-FixtureProject 'local-backup-overlap-target'
    Write-Utf8File (Join-Path $localBackupOverlapTarget 'marker.txt') 'UNCHANGED'
    Assert-InstallerRejectedWithoutMutation `
        'backup below local source' `
        $localBackupOverlapTarget `
        $archiveInstaller `
        (Join-Path $archiveRoot 'missing-backups') `
        @($archiveRoot, $localBackupOverlapTarget) `
        -ExpectedErrorPattern 'physically separate trees'

    $reverseOverlapContainer = New-FixtureProject 'reverse-overlap-container'
    $reverseOverlapSource = Join-Path $reverseOverlapContainer 'source'
    New-Item -ItemType Directory -Force -Path $reverseOverlapSource | Out-Null
    Get-ChildItem -LiteralPath $archiveRoot -Force | ForEach-Object {
        Copy-Item -LiteralPath $_.FullName -Destination $reverseOverlapSource -Recurse -Force
    }
    $reverseOverlapInstaller = Join-Path $reverseOverlapSource 'Toolkit\scripts\Install-SDP.ps1'
    Assert-InstallerRejectedWithoutMutation `
        'source below local project' `
        $reverseOverlapContainer `
        $reverseOverlapInstaller `
        '' `
        @($reverseOverlapContainer) `
        -ExpectedErrorPattern 'physically separate trees'

    $reverseBackupTarget = New-FixtureProject 'reverse-backup-overlap-target'
    Write-Utf8File (Join-Path $reverseBackupTarget 'marker.txt') 'UNCHANGED'
    Assert-InstallerRejectedWithoutMutation `
        'source below local backup' `
        $reverseBackupTarget `
        $reverseOverlapInstaller `
        $reverseOverlapContainer `
        @($reverseOverlapContainer, $reverseBackupTarget) `
        -ExpectedErrorPattern 'physically separate trees'

    if ($IsWindowsPlatform) {
        $extendedArchiveRoot = '\\?\' + $archiveRoot
        Assert-InstallerRejectedWithoutMutation `
            'extended local source/project overlap' `
            $extendedArchiveRoot `
            $archiveInstaller `
            '' `
            @($archiveRoot) `
            -ExpectedErrorPattern 'physically separate trees'

        $shortArchiveRoot = Get-ShortPathAlias $archiveRoot
        if (-not [string]::IsNullOrWhiteSpace($shortArchiveRoot)) {
            Assert-InstallerRejectedWithoutMutation `
                'short-name source/project overlap' `
                $shortArchiveRoot `
                $archiveInstaller `
                '' `
                @($archiveRoot) `
                -ExpectedErrorPattern 'physically separate trees'
        } else {
            $knownLongPath = $env:ProgramFiles
            $knownShortPath = if ([string]::IsNullOrWhiteSpace($knownLongPath)) {
                $null
            } else {
                Get-ShortPathAlias $knownLongPath
            }
            if (-not [string]::IsNullOrWhiteSpace($knownShortPath)) {
                Assert-Equal `
                    ([Sdp.Install.NativePath]::GetFileSystemIdentity($knownLongPath)) `
                    ([Sdp.Install.NativePath]::GetFileSystemIdentity($knownShortPath)) `
                    'short-name native identity'
                Write-Host '[SKIPPED] Integrated short-name source fixture is unavailable; native short/long identity passed.'
            } else {
                Write-Host '[SKIPPED] Short-name aliases are unavailable on this host.'
            }
        }

        $archiveUncRoot = Get-LoopbackAdministrativeUncPath $archiveRoot
        if ((-not [string]::IsNullOrWhiteSpace($archiveUncRoot)) -and
            (Test-Path -LiteralPath $archiveUncRoot -PathType Container)) {
            Assert-InstallerRejectedWithoutMutation `
                'local/UNC source/project overlap' `
                $archiveUncRoot `
                $archiveInstaller `
                '' `
                @($archiveRoot) `
                -ExpectedErrorPattern 'physically separate trees'

            $extendedUncArchiveRoot = '\\?\UNC\' + $archiveUncRoot.Substring(2)
            Assert-InstallerRejectedWithoutMutation `
                'extended UNC source/project overlap' `
                $extendedUncArchiveRoot `
                $archiveInstaller `
                '' `
                @($archiveRoot) `
                -ExpectedErrorPattern 'physically separate trees'

            $uncBackupOverlapTarget = New-FixtureProject 'unc-backup-overlap-target'
            Write-Utf8File (Join-Path $uncBackupOverlapTarget 'marker.txt') 'UNCHANGED'
            Assert-InstallerRejectedWithoutMutation `
                'backup below UNC source alias' `
                $uncBackupOverlapTarget `
                $archiveInstaller `
                (Join-Path $archiveUncRoot 'missing-unc-backups') `
                @($archiveRoot, $uncBackupOverlapTarget) `
                -ExpectedErrorPattern 'physically separate trees'

            $uncSiblingTarget = New-FixtureProject 'unc-physical-sibling'
            $uncSiblingAlias = Get-LoopbackAdministrativeUncPath $uncSiblingTarget
            $uncSiblingBefore = Get-TreeFingerprint $uncSiblingTarget
            $uncSiblingPlanJson = Invoke-PlanJson `
                $uncSiblingAlias `
                -InstallerPath $archiveInstaller
            $uncSiblingPlan = $uncSiblingPlanJson | ConvertFrom-Json
            Assert-True ([bool]$uncSiblingPlan.canApply) 'same-volume local/UNC sibling was rejected'
            Assert-Equal `
                $uncSiblingBefore `
                (Get-TreeFingerprint $uncSiblingTarget) `
                'same-volume local/UNC sibling plan mutated target'

        } else {
            Write-Host '[SKIPPED] Loopback administrative UNC fixtures are unavailable on this host.'
        }
    }

    # Runtime preflight rejects policy combinations and semantic-anchor drift before mutation.
    $archiveManifestPath = Join-Path $archiveRoot 'Toolkit\SDP-install.manifest.json'
    $archiveManifestOriginal = Get-Content -Raw -LiteralPath $archiveManifestPath

    # Manifest destinations are a case-insensitive file topology: no entry may
    # be another entry's ancestor or descendant.
    try {
        $prefixConflictContract = $archiveManifestOriginal | ConvertFrom-Json
        $prefixConflictEntry = @(
            $prefixConflictContract.entries | Where-Object { $_.id -eq 'project-agents' }
        )[0]
        $prefixConflictEntry.destination = 'agents.md/README.md'
        Write-Utf8File `
            $archiveManifestPath `
            ($prefixConflictContract | ConvertTo-Json -Depth 60)
        $prefixConflictTarget = New-FixtureProject 'destination-prefix-conflict'
        Write-Utf8File (Join-Path $prefixConflictTarget 'marker.txt') 'UNCHANGED'
        Assert-InstallerRejectedWithoutMutation `
            'destination prefix conflict' `
            $prefixConflictTarget `
            $archiveInstaller `
            '' `
            @($prefixConflictTarget) `
            -ExpectedErrorPattern 'ancestor/descendant destination conflict'
    } finally {
        Write-Utf8File $archiveManifestPath $archiveManifestOriginal
    }

    # Every manifest destination is preflighted, including uncreated children
    # beneath an existing file ancestor such as canonical `.codex`.
    $fileAncestorTarget = New-FixtureProject 'destination-file-ancestor'
    Write-Utf8File (Join-Path $fileAncestorTarget '.codex') 'PROJECT FILE'
    Write-Utf8File (Join-Path $fileAncestorTarget 'marker.txt') 'UNCHANGED'
    Assert-InstallerRejectedWithoutMutation `
        'destination file ancestor' `
        $fileAncestorTarget `
        $archiveInstaller `
        '' `
        @($fileAncestorTarget) `
        -ExpectedErrorPattern 'existing non-directory ancestor'

    $invalidPolicyContract = $archiveManifestOriginal | ConvertFrom-Json
    $invalidPolicyEntry = @($invalidPolicyContract.entries | Where-Object { $_.id -eq 'managed-framework-readme' })[0]
    $invalidPolicyEntry.migrationPolicy = 'preserve-existing-agents'
    Write-Utf8File $archiveManifestPath ($invalidPolicyContract | ConvertTo-Json -Depth 50)
    $invalidPolicyTarget = New-FixtureProject 'invalid-policy-target'
    $beforeInvalidPolicy = Get-TreeFingerprint $invalidPolicyTarget
    $invalidPolicyFailed = $false
    try { Invoke-PlanJson $invalidPolicyTarget -InstallerPath $archiveInstaller | Out-Null } catch { $invalidPolicyFailed = $true }
    Assert-True $invalidPolicyFailed 'unsafe policy combination passed installer preflight'
    Assert-Equal $beforeInvalidPolicy (Get-TreeFingerprint $invalidPolicyTarget) 'unsafe policy preflight mutated target'

    Write-Utf8File $archiveManifestPath $archiveManifestOriginal
    $driftContract = $archiveManifestOriginal | ConvertFrom-Json
    $driftEntry = @($driftContract.entries | Where-Object { $_.id -eq 'generated-installed-toolkit-manifest' })[0]
    $driftEntry.destination = 'SDP/Framework/drifted-installed-toolkit.manifest.yaml'
    Write-Utf8File $archiveManifestPath ($driftContract | ConvertTo-Json -Depth 50)
    $driftTarget = New-FixtureProject 'destination-drift-target'
    $beforeDrift = Get-TreeFingerprint $driftTarget
    $driftFailed = $false
    try { Invoke-PlanJson $driftTarget -InstallerPath $archiveInstaller | Out-Null } catch { $driftFailed = $true }
    Assert-True $driftFailed 'installed-manifest destination drift passed preflight'
    Assert-Equal $beforeDrift (Get-TreeFingerprint $driftTarget) 'destination drift preflight mutated target'
    Write-Utf8File $archiveManifestPath $archiveManifestOriginal

    # Closed-world PowerShell validation covers the complete v1 object shapes,
    # array minima, nested policies, ownership source classes and governing pairs.
    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $invalidContract | Add-Member -NotePropertyName unexpectedRoot -NotePropertyValue $true
    Assert-InvalidArchiveContract 'contract-unknown-root' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $invalidContract.psobject.Properties.Remove('exclusions')
    Assert-InvalidArchiveContract 'contract-missing-required-root' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $invalidContract.capabilities = @()
    Assert-InvalidArchiveContract 'contract-empty-capabilities' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $invalidContract.generators = @($invalidContract.generators[0])
    Assert-InvalidArchiveContract 'contract-generator-minimum' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $invalidContract.entries = @()
    Assert-InvalidArchiveContract 'contract-entry-minimum' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $invalidContract.exclusions = @()
    Assert-InvalidArchiveContract 'contract-exclusion-minimum' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $invalidContract.exclusions = @(
        $invalidContract.exclusions | Where-Object { $_.path -cne 'SDP-DOCUMENT-GUIDE.md' }
    )
    Assert-InvalidArchiveContract 'contract-required-legacy-exclusion' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $installedGeneratorContract = @($invalidContract.generators | Where-Object { $_.id -eq 'installed-toolkit-manifest' })[0]
    $installedGeneratorContract | Add-Member -NotePropertyName ignoredGeneratorPolicy -NotePropertyValue 'unsafe'
    Assert-InvalidArchiveContract 'contract-unknown-generator-policy' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $installedGeneratorContract = @($invalidContract.generators | Where-Object { $_.id -eq 'installed-toolkit-manifest' })[0]
    $installedGeneratorContract.facts.psobject.Properties.Remove('skills')
    Assert-InvalidArchiveContract 'contract-missing-generator-fact' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $installedGeneratorContract = @($invalidContract.generators | Where-Object { $_.id -eq 'installed-toolkit-manifest' })[0]
    $installedGeneratorContract.dynamicFacts.sourceCommit | Add-Member -NotePropertyName ignoredFallback -NotePropertyValue 'HEAD'
    Assert-InvalidArchiveContract 'contract-unknown-dynamic-policy' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $entryContract = @($invalidContract.entries | Where-Object { $_.id -eq 'managed-framework-readme' })[0]
    $entryContract | Add-Member -NotePropertyName ignoredPolicy -NotePropertyValue 'replace-everything'
    Assert-InvalidArchiveContract 'contract-unknown-entry-policy' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $entryContract = @($invalidContract.entries | Where-Object { $_.id -eq 'managed-framework-readme' })[0]
    $entryContract.source = 'Toolkit/project-templates/sdp-root/Traceability/README.md'
    Assert-InvalidArchiveContract 'contract-wrong-source-class' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $invalidContract.exclusions[0] | Add-Member -NotePropertyName ignoredScope -NotePropertyValue 'lexical-only'
    Assert-InvalidArchiveContract 'contract-unknown-exclusion-policy' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $entryContract = @($invalidContract.entries | Where-Object { $_.id -eq 'project-current-index' })[0]
    $entryContract.governing.capability = 'sdp.traceability.relations.v1'
    Assert-InvalidArchiveContract 'contract-governing-pair-drift' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    # Every portable manifest path class rejects Windows aliases, reserved names,
    # control/ADS forms and Git administration paths before target mutation.
    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $entryContract = @($invalidContract.entries | Where-Object { $_.id -eq 'managed-framework-readme' })[0]
    $entryContract.source = 'Toolkit/payload/sdp-root/Framework/README.md.'
    Assert-InvalidArchiveContract 'path-source-trailing-dot' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $entryContract = @($invalidContract.entries | Where-Object { $_.id -eq 'project-agents' })[0]
    $entryContract.destination = '.git/config'
    Assert-InvalidArchiveContract 'path-destination-git' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    foreach ($pathCase in @(
        @('reserved-device', 'SDP/CON/file.md'),
        @('trailing-dot', 'SDP/portable.'),
        @('trailing-space', 'SDP/portable '),
        @('ads-colon', 'SDP/portable.md:stream'),
        @('short-name-alias', 'SDP/RELEAS~1.MD'),
        @('superscript-device', "SDP/COM$([char]0x00B9).txt"),
        @('control-character', ('SDP/bad' + [char]1 + '.md'))
    )) {
        $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
        $entryContract = @($invalidContract.entries | Where-Object { $_.id -eq 'project-agents' })[0]
        $entryContract.destination = [string]$pathCase[1]
        Assert-InvalidArchiveContract ("path-destination-" + $pathCase[0]) $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller
    }

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $entryContract = @($invalidContract.entries | Where-Object { $_.id -eq 'project-current-index' })[0]
    $entryContract.governing.schema = 'Toolkit/schemas/current-index.schema.json:stream'
    Assert-InvalidArchiveContract 'path-governing-ads' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $invalidContract.exclusions[0].path = 'Releases.'
    Assert-InvalidArchiveContract 'path-exclusion-alias' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    $invalidContract = $archiveManifestOriginal | ConvertFrom-Json
    $entryContract = @($invalidContract.entries | Where-Object { $_.id -eq 'project-agents' })[0]
    $entryContract.destination = 'agents.md'
    Assert-InvalidArchiveContract 'path-destination-case-collision' $invalidContract $archiveManifestPath $archiveManifestOriginal $archiveInstaller

    # Final releases outrank prereleases; build metadata changes identity but not
    # precedence, so differing builds refresh rather than block as a downgrade.
    try {
        $prereleaseSourceContract = $archiveManifestOriginal | ConvertFrom-Json
        $prereleaseSourceContract.toolkitVersion = '0.2.0-rc.1'
        $prereleaseGenerator = @($prereleaseSourceContract.generators | Where-Object { $_.id -eq 'installed-toolkit-manifest' })[0]
        $prereleaseGenerator.facts.toolkitVersion = '0.2.0-rc.1'
        Write-Utf8File $archiveManifestPath ($prereleaseSourceContract | ConvertTo-Json -Depth 60)
        $finalTarget = New-FixtureProject 'final-over-prerelease-source'
        Write-Utf8File (Join-Path $finalTarget 'SDP\Framework\installed-toolkit.manifest.yaml') $semverSeedContent
        $finalPlan = (Invoke-PlanJson $finalTarget -InstallerPath $archiveInstaller) | ConvertFrom-Json
        Assert-True (-not [bool]$finalPlan.canApply) 'prerelease source attempted to downgrade final installation'
        Assert-Equal 'downgrade-blocked' ([string]$finalPlan.actions[0].reason) 'final-over-prerelease block reason'
    } finally {
        Write-Utf8File $archiveManifestPath $archiveManifestOriginal
    }

    try {
        $buildSourceContract = $archiveManifestOriginal | ConvertFrom-Json
        $buildSourceContract.toolkitVersion = '0.2.0+source.2'
        $buildGenerator = @($buildSourceContract.generators | Where-Object { $_.id -eq 'installed-toolkit-manifest' })[0]
        $buildGenerator.facts.toolkitVersion = '0.2.0+source.2'
        Write-Utf8File $archiveManifestPath ($buildSourceContract | ConvertTo-Json -Depth 60)
        $buildTarget = New-FixtureProject 'build-identity-change'
        $buildInstalled = $semverSeedContent -replace '(?m)^toolkitVersion: "0\.2\.0"$', 'toolkitVersion: "0.2.0+consumer.9"'
        Write-Utf8File (Join-Path $buildTarget 'SDP\Framework\installed-toolkit.manifest.yaml') $buildInstalled
        $buildPlan = (Invoke-PlanJson $buildTarget -InstallerPath $archiveInstaller) | ConvertFrom-Json
        Assert-True ([bool]$buildPlan.canApply) 'SemVer build metadata incorrectly determined precedence'
        Assert-Equal '0.2.0+consumer.9' ([string]$buildPlan.installedToolkitVersion) 'build identity old version'
        Assert-Equal '0.2.0+source.2' ([string]$buildPlan.toolkitVersion) 'build identity new version'
    } finally {
        Write-Utf8File $archiveManifestPath $archiveManifestOriginal
    }

    # Existing destination and source ancestors may not redirect I/O through
    # junctions or symlinks. These regressions run with junctions on Windows and
    # symbolic links on Unix-like hosts.
    $redirectOutside = New-FixtureProject 'redirect-outside'
    Write-Utf8File (Join-Path $redirectOutside 'marker.txt') 'UNCHANGED'
    $redirectProject = New-FixtureProject 'redirect-project'
    $redirectLink = Join-Path $redirectProject 'SDP'
    New-DirectoryLink $redirectLink $redirectOutside
    try {
        $outsideBefore = Get-TreeFingerprint $redirectOutside
        $redirectFailed = $false
        try { Invoke-PlanJson $redirectProject | Out-Null } catch { $redirectFailed = $true }
        Assert-True $redirectFailed 'destination ancestor link passed physical containment'
        Assert-Equal $outsideBefore (Get-TreeFingerprint $redirectOutside) 'destination ancestor link mutated outside target'
    } finally {
        Remove-DirectoryLink $redirectLink
    }

    $linkedProjectTarget = New-FixtureProject 'linked-project-target'
    Write-Utf8File (Join-Path $linkedProjectTarget 'marker.txt') 'UNCHANGED'
    $linkedProjectRoot = Join-Path $TestRoot 'linked-project-root'
    New-DirectoryLink $linkedProjectRoot $linkedProjectTarget
    try {
        $linkedBefore = Get-TreeFingerprint $linkedProjectTarget
        $linkedFailed = $false
        try { Invoke-PlanJson $linkedProjectRoot | Out-Null } catch { $linkedFailed = $true }
        Assert-True $linkedFailed 'linked project root passed physical containment'
        Assert-Equal $linkedBefore (Get-TreeFingerprint $linkedProjectTarget) 'linked project root mutated physical target'
    } finally {
        Remove-DirectoryLink $linkedProjectRoot
    }

    $linkedArchiveRoot = Join-Path $TestRoot 'source-link-archive'
    New-Item -ItemType Directory -Force -Path $linkedArchiveRoot | Out-Null
    Get-ChildItem -LiteralPath $archiveRoot -Force | ForEach-Object {
        Copy-Item -LiteralPath $_.FullName -Destination $linkedArchiveRoot -Recurse -Force
    }
    $linkedFramework = Join-Path $linkedArchiveRoot 'Toolkit\payload\sdp-root\Framework'
    Assert-True ($linkedFramework.StartsWith($TestRoot, [System.StringComparison]::OrdinalIgnoreCase)) 'source-link fixture escaped test root'
    Remove-Item -LiteralPath $linkedFramework -Recurse -Force
    New-DirectoryLink $linkedFramework (Join-Path $RepositoryRoot 'Toolkit\payload\sdp-root\Framework')
    try {
        $linkedSourceInstaller = Join-Path $linkedArchiveRoot 'Toolkit\scripts\Install-SDP.ps1'
        $linkedSourceTarget = New-FixtureProject 'source-link-target'
        $linkedSourceBefore = Get-TreeFingerprint $linkedSourceTarget
        $linkedSourceFailed = $false
        try { Invoke-PlanJson $linkedSourceTarget -InstallerPath $linkedSourceInstaller | Out-Null } catch { $linkedSourceFailed = $true }
        Assert-True $linkedSourceFailed 'linked source ancestor passed physical containment'
        Assert-Equal $linkedSourceBefore (Get-TreeFingerprint $linkedSourceTarget) 'linked source ancestor mutated target'
    } finally {
        Remove-DirectoryLink $linkedFramework
    }

    $danglingOutside = New-FixtureProject 'dangling-link-outside'
    $danglingTarget = Join-Path $danglingOutside 'escaped-agents.md'
    Write-Utf8File $danglingTarget 'OUTSIDE'
    $danglingProject = New-FixtureProject 'dangling-link-project'
    $danglingLink = Join-Path $danglingProject 'AGENTS.md'
    $danglingLinkCreated = Try-NewFileSymbolicLink $danglingLink $danglingTarget
    if ($danglingLinkCreated) {
        try {
            [System.IO.File]::Delete($danglingTarget)
            $danglingFailed = $false
            try { Invoke-PlanJson $danglingProject | Out-Null } catch { $danglingFailed = $true }
            Assert-True $danglingFailed 'dangling destination symlink passed physical containment'
            Assert-True (-not (Test-Path -LiteralPath $danglingTarget)) 'dangling destination symlink wrote outside target'
        } finally {
            [System.IO.File]::Delete($danglingLink)
        }
    } else {
        [System.IO.File]::Delete($danglingTarget)
        Write-Host '[SKIPPED] File symlink regression is unavailable on this host.'
    }

    # A sibling named SDP-Analyzer is valid; a child of the Toolkit repo is not.
    $sibling = Join-Path (Split-Path -Parent $RepositoryRoot) ("SDP-Analyzer-fixture-" + [guid]::NewGuid().ToString('N'))
    try {
        New-Item -ItemType Directory -Force -Path $sibling | Out-Host
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
