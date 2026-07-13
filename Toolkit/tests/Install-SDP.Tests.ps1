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
        [switch]$Force
    )
    $parameters = @{
        ProjectRoot = $Target
        PlanJson = $true
    }
    if ($Initialize) { $parameters.InitializeProjectStructure = $true }
    if ($Force) { $parameters.ForceManagedFiles = $true }
    $output = & $InstallerPath @parameters
    return [string]::Join("`n", [string[]]@($output))
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
    }
    Assert-True (
        @($defaultActions | Where-Object { $_.entryId -eq 'project-sprints-readme' }).Count -eq 0
    ) 'default plan selected initialize-only content'

    $initializePlanJson = Invoke-PlanJson $planOnly -Initialize
    Assert-PlanConforms $initializePlanJson 'initialize'
    $initializePlan = $initializePlanJson | ConvertFrom-Json
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
    & $Installer -ProjectRoot $legacy | Out-Host
    Assert-Equal 'OLD LOCAL AGENTS' (Get-Content -Raw -LiteralPath (Join-Path $legacy 'AGENTS-project.md')) 'old AGENTS not migrated'
    Assert-Equal 'LEGACY REQUIREMENTS' (Get-Content -Raw -LiteralPath (Join-Path $legacy 'SDP\03--Requirements\requirements.md')) 'legacy document overwritten'

    # If AGENTS-project.md already exists, migration preserves both rule sets.
    $legacyConflict = New-FixtureProject 'legacy-agents-conflict'
    Write-Utf8File (Join-Path $legacyConflict 'AGENTS.md') 'OLD CONFLICTING AGENTS'
    Write-Utf8File (Join-Path $legacyConflict 'AGENTS-project.md') 'EXISTING PROJECT AGENTS'
    $beforeConflictPlan = Get-TreeFingerprint $legacyConflict
    $conflictPlanJson = Invoke-PlanJson $legacyConflict
    Assert-PlanConforms $conflictPlanJson 'AGENTS migration conflict'
    Assert-Equal $beforeConflictPlan (Get-TreeFingerprint $legacyConflict) 'AGENTS migration plan mutated target'
    $conflictPlan = $conflictPlanJson | ConvertFrom-Json
    Assert-Equal 1 @(
        $conflictPlan.actions | Where-Object { $_.reason -eq 'preserve-existing-agents' }
    ).Count 'AGENTS conflict plan preservation action count'
    & $Installer -ProjectRoot $legacyConflict | Out-Host
    Assert-Equal 'EXISTING PROJECT AGENTS' (Get-Content -Raw -LiteralPath (Join-Path $legacyConflict 'AGENTS-project.md')) 'existing AGENTS-project overwritten during migration'
    $migrationFiles = @(Get-ChildItem -LiteralPath $legacyConflict -File -Filter 'AGENTS-project.migration-*.md')
    Assert-Equal 1 $migrationFiles.Count 'AGENTS conflict migration file count'
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

    # Runtime preflight rejects policy combinations and semantic-anchor drift before mutation.
    $archiveManifestPath = Join-Path $archiveRoot 'Toolkit\SDP-install.manifest.json'
    $archiveManifestOriginal = Get-Content -Raw -LiteralPath $archiveManifestPath
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
