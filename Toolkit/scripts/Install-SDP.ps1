[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$ProjectRoot,

    [switch]$ForceManagedFiles,

    [switch]$InitializeProjectStructure,

    [switch]$Preview,

    [switch]$PlanJson,

    [string]$BackupRoot
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$SupportedInstallManifestSchemas = @('1.0')
$SupportedInstalledManifestSchemas = @('1.0')
$SupportedProjectManifestSchemas = @('1.0')
$RunStamp = (Get-Date).ToUniversalTime().ToString('yyyyMMdd-HHmmssZ')
$PathComparison = if ([System.IO.Path]::DirectorySeparatorChar -eq '\') {
    [System.StringComparison]::OrdinalIgnoreCase
} else {
    [System.StringComparison]::Ordinal
}

function Get-YamlScalar {
    param([string]$Content, [string]$Name)

    $escaped = [regex]::Escape($Name)
    $pattern = '(?m)^\s*{0}\s*:\s*[''"]?([^''"#\r\n]+)' -f $escaped
    $match = [regex]::Match(
        $Content,
        $pattern,
        [System.Text.RegularExpressions.RegexOptions]::CultureInvariant
    )
    if (-not $match.Success) { return $null }
    return $match.Groups[1].Value.Trim()
}

function ConvertTo-SemVerCore {
    param([string]$Version)

    $match = [regex]::Match(
        $Version,
        '^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$'
    )
    if (-not $match.Success) {
        throw "Invalid Toolkit SemVer: $Version"
    }
    return [pscustomobject]@{
        Major = [int64]$match.Groups[1].Value
        Minor = [int64]$match.Groups[2].Value
        Patch = [int64]$match.Groups[3].Value
    }
}

function Compare-SemVerCore {
    param([string]$Left, [string]$Right)

    $leftVersion = ConvertTo-SemVerCore $Left
    $rightVersion = ConvertTo-SemVerCore $Right
    foreach ($field in @('Major', 'Minor', 'Patch')) {
        if ($leftVersion.$field -lt $rightVersion.$field) { return -1 }
        if ($leftVersion.$field -gt $rightVersion.$field) { return 1 }
    }
    return 0
}

function Assert-PortableRelativePath {
    param([string]$Value, [string]$Label)

    if ([string]::IsNullOrWhiteSpace($Value)) {
        throw "$Label must be a non-empty relative path."
    }
    if ($Value.Contains('\')) {
        throw "$Label must use portable '/' separators: $Value"
    }
    if ($Value.StartsWith('/') -or $Value -match '^[A-Za-z]:') {
        throw "$Label must be relative: $Value"
    }
    $segments = $Value.Split('/')
    if ($segments.Count -eq 0 -or $segments -contains '' -or
        $segments -contains '.' -or $segments -contains '..') {
        throw "$Label is not normalized or contains traversal: $Value"
    }
    if (($segments -join '/') -cne $Value) {
        throw "$Label is not normalized: $Value"
    }
}

function Assert-AllowedValue {
    param(
        [string]$Value,
        [string[]]$Allowed,
        [string]$Label
    )

    if ($Allowed -cnotcontains $Value) {
        throw "$Label has unsupported value '$Value'."
    }
}

function Join-PortablePath {
    param([string]$Base, [string]$Relative, [string]$Label)

    Assert-PortableRelativePath $Relative $Label
    $path = $Base
    foreach ($segment in $Relative.Split('/')) {
        $path = Join-Path $path $segment
    }
    $fullBase = [System.IO.Path]::GetFullPath($Base).TrimEnd(
        [System.IO.Path]::DirectorySeparatorChar,
        [System.IO.Path]::AltDirectorySeparatorChar
    )
    $full = [System.IO.Path]::GetFullPath($path)
    $prefix = $fullBase + [System.IO.Path]::DirectorySeparatorChar
    if (($full -cne $fullBase) -and
        (-not $full.StartsWith($prefix, $PathComparison))) {
        throw "$Label escapes its root: $Relative"
    }
    return $full
}

function Test-PortablePathWithin {
    param([string]$Candidate, [string]$Parent, [string]$Kind)

    if ($Candidate -ceq $Parent) { return $true }
    return ($Kind -eq 'tree') -and $Candidate.StartsWith($Parent + '/')
}

function ConvertTo-YamlQuotedScalar {
    param([string]$Value)

    return '"' + $Value.Replace('\', '\\').Replace('"', '\"') + '"'
}

$ToolkitRoot = Split-Path -Parent $PSScriptRoot
$InstallManifestPath = Join-Path $ToolkitRoot 'SDP-install.manifest.json'
if (-not (Test-Path -LiteralPath $InstallManifestPath -PathType Leaf)) {
    throw "Canonical installation manifest is missing: $InstallManifestPath"
}
try {
    $InstallManifest = Get-Content -Raw -LiteralPath $InstallManifestPath |
        ConvertFrom-Json
} catch {
    throw "Cannot parse canonical installation manifest: $($_.Exception.Message)"
}
if ($SupportedInstallManifestSchemas -notcontains [string]$InstallManifest.schemaVersion) {
    throw "Unsupported installation manifest schema '$($InstallManifest.schemaVersion)'."
}
if ([string]$InstallManifest.contractId -cne 'sdp-install') {
    throw "Unsupported installation contract '$($InstallManifest.contractId)'."
}
if ([string]$InstallManifest.sources.repositoryRoot -cne '..' -or
    [string]$InstallManifest.sources.pathStyle -cne 'forward-slash-relative') {
    throw 'Unsupported installation source-root or path-style contract.'
}

$RepositoryRoot = [System.IO.Path]::GetFullPath((Join-Path $ToolkitRoot '..'))
$ProjectRoot = [System.IO.Path]::GetFullPath((Resolve-Path $ProjectRoot).Path)
$SdpRoot = Join-Path $ProjectRoot 'SDP'
$ExpectedInstalledManifestRelative = 'SDP/Framework/installed-toolkit.manifest.yaml'
$ExpectedProjectManifestRelative = 'SDP/SDP-project.manifest.yaml'

$repositoryTrimmed = $RepositoryRoot.TrimEnd(
    [System.IO.Path]::DirectorySeparatorChar,
    [System.IO.Path]::AltDirectorySeparatorChar
)
$repositoryPrefix = $repositoryTrimmed + [System.IO.Path]::DirectorySeparatorChar
$projectIsRepository = $ProjectRoot.Equals($repositoryTrimmed, $PathComparison)
$projectIsInsideRepository = $ProjectRoot.StartsWith($repositoryPrefix, $PathComparison)
if ($projectIsRepository -or $projectIsInsideRepository) {
    throw "The consuming project must not be inside the SDP repository: $RepositoryRoot"
}

if ([string]::IsNullOrWhiteSpace($BackupRoot)) {
    $BackupRoot = Join-Path $SdpRoot ".sdp-backups\$RunStamp"
} elseif (-not [System.IO.Path]::IsPathRooted($BackupRoot)) {
    $BackupRoot = Join-Path $ProjectRoot $BackupRoot
}
$BackupRoot = [System.IO.Path]::GetFullPath($BackupRoot)

$ToolkitVersion = [string]$InstallManifest.toolkitVersion
[void](ConvertTo-SemVerCore $ToolkitVersion)
$Capabilities = @($InstallManifest.capabilities)
$CapabilityById = @{}
foreach ($capabilityValue in $Capabilities) {
    $capability = [string]$capabilityValue
    if (($capability -cnotmatch '^sdp\.[a-z0-9.-]+\.v[0-9]+$') -or
        $CapabilityById.ContainsKey($capability)) {
        throw "Installation contract contains an invalid or duplicate capability '$capability'."
    }
    $CapabilityById[$capability] = $true
}
if (-not $CapabilityById.ContainsKey('sdp.install.v1')) {
    throw "Installation contract is missing capability 'sdp.install.v1'."
}
$Generators = @($InstallManifest.generators)
$Entries = @($InstallManifest.entries)
$Exclusions = @($InstallManifest.exclusions)
$GeneratorById = @{}
$EntryById = @{}
$DestinationByKey = @{}
$ExclusionByKey = @{}
$ExclusionRows = New-Object System.Collections.ArrayList

foreach ($generator in $Generators) {
    $generatorId = [string]$generator.id
    if (($generatorId -cnotmatch '^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$') -or
        $GeneratorById.ContainsKey($generatorId)) {
        throw "Installation contract contains a missing or duplicate generator ID '$generatorId'."
    }
    $generatorType = [string]$generator.type
    Assert-AllowedValue $generatorType @('installed-toolkit-manifest', 'empty-ledger') "generator '$generatorId' type"
    if (($generatorType -eq 'installed-toolkit-manifest') -and
        (($generatorId -cne 'installed-toolkit-manifest') -or
        ([string]$generator.format -cne 'yaml'))) {
        throw "Generator '$generatorId' must use YAML format."
    }
    if (($generatorType -eq 'empty-ledger') -and
        (($generatorId -cne 'empty-ledger') -or
        ([string]$generator.format -cne 'ndjson'))) {
        throw "Generator '$generatorId' must use NDJSON format."
    }
    $GeneratorById[$generatorId] = $generator
}
foreach ($requiredGenerator in @('installed-toolkit-manifest', 'empty-ledger')) {
    if (-not $GeneratorById.ContainsKey($requiredGenerator)) {
        throw "Installation contract is missing generator '$requiredGenerator'."
    }
}
$InstalledGenerator = $GeneratorById['installed-toolkit-manifest']
if ([string]$InstalledGenerator.type -cne 'installed-toolkit-manifest') {
    throw "Generator 'installed-toolkit-manifest' has the wrong type."
}
if ([string]$GeneratorById['empty-ledger'].type -cne 'empty-ledger') {
    throw "Generator 'empty-ledger' has the wrong type."
}
$InstallerVersion = [string]$InstalledGenerator.facts.installerVersion
[void](ConvertTo-SemVerCore $InstallerVersion)
if ([string]$InstalledGenerator.facts.toolkitVersion -cne $ToolkitVersion) {
    throw "Installed-manifest generator Toolkit version does not match the contract Toolkit version."
}
if ($SupportedInstalledManifestSchemas -cnotcontains
    [string]$InstalledGenerator.facts.schemaVersion) {
    throw "Installed-manifest generator declares an unsupported schema version."
}
foreach ($factName in @(
    'toolkitVersion', 'frameworkVersion', 'agentsContractVersion', 'installerVersion'
)) {
    [void](ConvertTo-SemVerCore ([string]$InstalledGenerator.facts.$factName))
}
$skillFacts = $InstalledGenerator.facts.skills
if (($null -eq $skillFacts) -or (@($skillFacts.psobject.Properties).Count -eq 0)) {
    throw "Installed-manifest generator must declare installed skill versions."
}
foreach ($skillProperty in $skillFacts.psobject.Properties) {
    if ($skillProperty.Name -cnotmatch '^sdp-[a-z0-9]+(?:-[a-z0-9]+)*$') {
        throw "Installed-manifest generator has invalid skill ID '$($skillProperty.Name)'."
    }
    [void](ConvertTo-SemVerCore ([string]$skillProperty.Value))
}
if ((@($InstalledGenerator.facts.capabilities) -join "`n") -cne
    ($Capabilities -join "`n")) {
    throw "Installed-manifest generator capabilities do not match the contract capabilities."
}
$dynamicFacts = $InstalledGenerator.dynamicFacts
if (([string]$dynamicFacts.toolkitInstalledAt.source -cne 'utc-now') -or
    ([string]$dynamicFacts.toolkitInstalledAt.sameToolkitVersionPolicy -cne
        'preserve-existing') -or
    ([string]$dynamicFacts.sourceCommit.source -cne 'repository-head') -or
    ($null -ne $dynamicFacts.sourceCommit.unavailableValue)) {
    throw "Installed-manifest generator declares unsupported dynamic-fact policies."
}
if ([string]$GeneratorById['empty-ledger'].content -cne '') {
    throw "Empty-ledger generator content must be empty."
}

foreach ($exclusion in $Exclusions) {
    $excludedPath = [string]$exclusion.path
    $excludedKind = [string]$exclusion.kind
    Assert-PortableRelativePath $excludedPath "exclusion '$excludedPath'"
    if ($excludedKind -notin @('file', 'tree')) {
        throw "Unsupported exclusion kind '$excludedKind'."
    }
    Assert-AllowedValue ([string]$exclusion.reason) @(
        'repository-instance-state',
        'repository-release-metadata',
        'legacy-duplicate-root',
        'legacy-project-template',
        'development-only'
    ) "exclusion '$excludedPath' reason"
    $excludedKey = $excludedPath.ToLowerInvariant()
    if ($ExclusionByKey.ContainsKey($excludedKey)) {
        throw "Installation contract contains duplicate exclusion '$excludedPath'."
    }
    $ExclusionByKey[$excludedKey] = $excludedKind
    $excludedFull = Join-PortablePath $RepositoryRoot $excludedPath "exclusion '$excludedPath'"
    if (($excludedKind -eq 'file') -and
        (-not (Test-Path -LiteralPath $excludedFull -PathType Leaf))) {
        throw "Excluded file does not exist: $excludedPath"
    }
    if (($excludedKind -eq 'tree') -and
        (-not (Test-Path -LiteralPath $excludedFull -PathType Container))) {
        throw "Excluded tree does not exist: $excludedPath"
    }
    [void]$ExclusionRows.Add([pscustomobject]@{ Path = $excludedPath; Kind = $excludedKind })
}

foreach ($entry in $Entries) {
    $entryId = [string]$entry.id
    if (($entryId -cnotmatch '^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$') -or
        $EntryById.ContainsKey($entryId)) {
        throw "Installation contract contains a missing or duplicate entry ID '$entryId'."
    }
    $EntryById[$entryId] = $entry

    $entryKind = [string]$entry.kind
    $ownership = [string]$entry.ownership
    $selectionPolicy = [string]$entry.selectionPolicy
    $installPolicy = [string]$entry.installPolicy
    $refreshPolicy = [string]$entry.refreshPolicy
    $backupPolicy = [string]$entry.backupPolicy
    $forcePolicy = [string]$entry.forcePolicy
    $migrationPolicy = [string]$entry.migrationPolicy
    Assert-AllowedValue $entryKind @('copied', 'generated') "entry '$entryId' kind"
    Assert-AllowedValue $ownership @('toolkit-managed', 'project-owned') "entry '$entryId' ownership"
    Assert-AllowedValue $selectionPolicy @('default', 'initialize-only') "entry '$entryId' selectionPolicy"
    Assert-AllowedValue $installPolicy @('always', 'missing-only') "entry '$entryId' installPolicy"
    Assert-AllowedValue $refreshPolicy @('always', 'upgrade-or-force', 'never') "entry '$entryId' refreshPolicy"
    Assert-AllowedValue $backupPolicy @('before-replace', 'migration-aware', 'none') "entry '$entryId' backupPolicy"
    Assert-AllowedValue $forcePolicy @('replace-managed', 'preserve') "entry '$entryId' forcePolicy"
    Assert-AllowedValue $migrationPolicy @('none', 'preserve-existing-agents') "entry '$entryId' migrationPolicy"

    if ($ownership -eq 'project-owned') {
        if (($installPolicy -cne 'missing-only') -or
            ($refreshPolicy -cne 'never') -or
            ($backupPolicy -cne 'none') -or
            ($forcePolicy -cne 'preserve') -or
            ($migrationPolicy -cne 'none')) {
            throw "Entry '$entryId' has an unsafe project-owned policy combination."
        }
    } else {
        if ($forcePolicy -cne 'replace-managed') {
            throw "Entry '$entryId' has an unsupported toolkit-managed force policy."
        }
    }
    if (($selectionPolicy -eq 'initialize-only') -and ($ownership -ne 'project-owned')) {
        throw "Entry '$entryId' initialize-only content must be project-owned."
    }
    if ($migrationPolicy -eq 'preserve-existing-agents') {
        if (($ownership -ne 'toolkit-managed') -or
            ($backupPolicy -ne 'migration-aware') -or
            ([string]$entry.destination -cne 'AGENTS.md')) {
            throw "Entry '$entryId' has an unsupported AGENTS migration policy combination."
        }
    } elseif ($backupPolicy -eq 'migration-aware') {
        throw "Entry '$entryId' uses migration-aware backup without a migration policy."
    }
    if (($backupPolicy -eq 'before-replace') -and ($refreshPolicy -eq 'never')) {
        throw "Entry '$entryId' cannot back up before a refresh that is never allowed."
    }
    $governingCapability = [string]$entry.governing.capability
    if (-not $CapabilityById.ContainsKey($governingCapability)) {
        throw "Entry '$entryId' references undeclared capability '$governingCapability'."
    }

    $destination = [string]$entry.destination
    Assert-PortableRelativePath $destination "entry '$entryId' destination"
    $destinationKey = $destination.ToLowerInvariant()
    if ($DestinationByKey.ContainsKey($destinationKey)) {
        throw "Installation contract contains duplicate destination '$destination'."
    }
    $DestinationByKey[$destinationKey] = $entryId
    [void](Join-PortablePath $ProjectRoot $destination "entry '$entryId' destination")

    if ($entryKind -eq 'copied') {
        if ($entry.psobject.Properties.Name -contains 'generator') {
            throw "Copied entry '$entryId' must not declare a generator."
        }
        $source = [string]$entry.source
        Assert-PortableRelativePath $source "entry '$entryId' source"
        foreach ($excluded in $ExclusionRows) {
            if (Test-PortablePathWithin $source $excluded.Path $excluded.Kind) {
                throw "Entry '$entryId' uses explicitly excluded source '$source'."
            }
        }
        $sourcePath = Join-PortablePath $RepositoryRoot $source "entry '$entryId' source"
        if (-not (Test-Path -LiteralPath $sourcePath -PathType Leaf)) {
            throw "Entry '$entryId' source does not exist: $source"
        }
    } elseif ($entryKind -eq 'generated') {
        if ($entry.psobject.Properties.Name -contains 'source') {
            throw "Generated entry '$entryId' must not declare a source."
        }
        $generatorId = [string]$entry.generator
        if (-not $GeneratorById.ContainsKey($generatorId)) {
            throw "Entry '$entryId' references unknown generator '$generatorId'."
        }
    }
    $schemaPath = $entry.governing.schema
    if ($null -ne $schemaPath) {
        $schemaRelative = [string]$schemaPath
        $schemaFull = Join-PortablePath $RepositoryRoot $schemaRelative "entry '$entryId' schema"
        if (-not (Test-Path -LiteralPath $schemaFull -PathType Leaf)) {
            throw "Entry '$entryId' governing schema does not exist: $schemaRelative"
        }
    }
}

foreach ($requiredEntry in @(
    'managed-agents', 'project-agents', 'project-manifest',
    'generated-installed-toolkit-manifest', 'generated-empty-ledger'
)) {
    if (-not $EntryById.ContainsKey($requiredEntry)) {
        throw "Installation contract is missing required entry '$requiredEntry'."
    }
}

$InstalledManifestRelative = [string]$EntryById['generated-installed-toolkit-manifest'].destination
$ProjectManifestRelative = [string]$EntryById['project-manifest'].destination
if ($InstalledManifestRelative -cne $ExpectedInstalledManifestRelative) {
    throw "Installed-manifest entry destination drifted from '$ExpectedInstalledManifestRelative'."
}
if ($ProjectManifestRelative -cne $ExpectedProjectManifestRelative) {
    throw "Project-manifest entry destination drifted from '$ExpectedProjectManifestRelative'."
}
$InstalledManifestPath = Join-PortablePath $ProjectRoot $InstalledManifestRelative 'installed manifest destination'
$ProjectManifestPath = Join-PortablePath $ProjectRoot $ProjectManifestRelative 'project manifest destination'

$MigrationEntries = @($Entries | Where-Object {
    [string]$_.migrationPolicy -eq 'preserve-existing-agents'
})
if (($MigrationEntries.Count -ne 1) -or
    ([string]$MigrationEntries[0].id -cne 'managed-agents')) {
    throw "Installation contract must define managed-agents as its single AGENTS migration entry."
}

function Get-RepositorySourceCommit {
    if (-not (Test-Path -LiteralPath (Join-Path $RepositoryRoot '.git'))) { return $null }
    if (-not (Get-Command git -ErrorAction SilentlyContinue)) { return $null }
    $previousErrorActionPreference = $ErrorActionPreference
    try {
        $ErrorActionPreference = 'SilentlyContinue'
        $topLevel = & git -C $RepositoryRoot rev-parse --show-toplevel 2>$null
        $topLevelExitCode = $LASTEXITCODE
    } finally {
        $ErrorActionPreference = $previousErrorActionPreference
    }
    if (($topLevelExitCode -ne 0) -or (-not $topLevel)) { return $null }
    try {
        $resolvedTop = [System.IO.Path]::GetFullPath(($topLevel | Select-Object -First 1).Trim())
    } catch {
        return $null
    }
    if (-not $resolvedTop.Equals($repositoryTrimmed, $PathComparison)) { return $null }
    try {
        $ErrorActionPreference = 'SilentlyContinue'
        $head = & git -C $RepositoryRoot rev-parse HEAD 2>$null
        $headExitCode = $LASTEXITCODE
    } finally {
        $ErrorActionPreference = $previousErrorActionPreference
    }
    if (($headExitCode -ne 0) -or (-not $head)) { return $null }
    return ($head | Select-Object -First 1).Trim()
}

$BlockReason = $null
$BlockEntryId = $null
$BlockDestination = $null
$BlockOwnership = $null
$InstalledSchemaVersion = $null
$InstalledToolkitVersion = $null
$PreviousInstalledAt = $null

if (Test-Path -LiteralPath $ProjectManifestPath -PathType Leaf) {
    $projectManifestContent = [System.IO.File]::ReadAllText($ProjectManifestPath)
    $projectSchemaVersion = Get-YamlScalar $projectManifestContent 'schemaVersion'
    if ([string]::IsNullOrWhiteSpace($projectSchemaVersion)) {
        $BlockReason = 'malformed-project-manifest'
    } elseif ($SupportedProjectManifestSchemas -notcontains $projectSchemaVersion) {
        $BlockReason = 'unsupported-project-schema'
    }
    if ($BlockReason) {
        $BlockEntryId = 'project-manifest'
        $BlockDestination = $ProjectManifestRelative
        $BlockOwnership = 'project-owned'
    }
}

if ((-not $BlockReason) -and
    (Test-Path -LiteralPath $InstalledManifestPath -PathType Leaf)) {
    $installedContent = [System.IO.File]::ReadAllText($InstalledManifestPath)
    $InstalledSchemaVersion = Get-YamlScalar $installedContent 'schemaVersion'
    $InstalledToolkitVersion = Get-YamlScalar $installedContent 'toolkitVersion'
    $PreviousInstalledAt = Get-YamlScalar $installedContent 'toolkitInstalledAt'
    if ([string]::IsNullOrWhiteSpace($InstalledSchemaVersion) -or
        [string]::IsNullOrWhiteSpace($InstalledToolkitVersion)) {
        $BlockReason = 'malformed-installed-manifest'
    } elseif ($SupportedInstalledManifestSchemas -notcontains $InstalledSchemaVersion) {
        $BlockReason = 'unsupported-installed-schema'
    } else {
        try {
            if ((Compare-SemVerCore $InstalledToolkitVersion $ToolkitVersion) -gt 0) {
                $BlockReason = 'downgrade-blocked'
            }
        } catch {
            $BlockReason = 'malformed-installed-manifest'
        }
    }
    if ($BlockReason) {
        $BlockEntryId = 'generated-installed-toolkit-manifest'
        $BlockDestination = $InstalledManifestRelative
        $BlockOwnership = 'toolkit-managed'
    }
}

$IsUpgrade = [string]::IsNullOrWhiteSpace($InstalledToolkitVersion) -or
    ($InstalledToolkitVersion -ne $ToolkitVersion)
if (($InstalledToolkitVersion -eq $ToolkitVersion) -and
    (-not [string]::IsNullOrWhiteSpace($PreviousInstalledAt))) {
    $InstalledAt = $PreviousInstalledAt
} else {
    $InstalledAt = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')
}
$SourceCommit = Get-RepositorySourceCommit

function Get-InstalledManifestContent {
    $facts = $InstalledGenerator.facts
    $lines = New-Object System.Collections.ArrayList
    [void]$lines.Add('schemaVersion: ' + (ConvertTo-YamlQuotedScalar ([string]$facts.schemaVersion)))
    [void]$lines.Add('toolkitVersion: ' + (ConvertTo-YamlQuotedScalar ([string]$facts.toolkitVersion)))
    [void]$lines.Add('frameworkVersion: ' + (ConvertTo-YamlQuotedScalar ([string]$facts.frameworkVersion)))
    [void]$lines.Add('agentsContractVersion: ' + (ConvertTo-YamlQuotedScalar ([string]$facts.agentsContractVersion)))
    [void]$lines.Add('installerVersion: ' + (ConvertTo-YamlQuotedScalar ([string]$facts.installerVersion)))
    [void]$lines.Add('toolkitInstalledAt: ' + (ConvertTo-YamlQuotedScalar $InstalledAt))
    if ($null -eq $SourceCommit) {
        [void]$lines.Add('sourceCommit: null')
    } else {
        [void]$lines.Add('sourceCommit: ' + (ConvertTo-YamlQuotedScalar $SourceCommit))
    }
    [void]$lines.Add('skills:')
    foreach ($property in $facts.skills.psobject.Properties) {
        [void]$lines.Add('  ' + $property.Name + ': ' +
            (ConvertTo-YamlQuotedScalar ([string]$property.Value)))
    }
    [void]$lines.Add('capabilities:')
    foreach ($capability in @($facts.capabilities)) {
        [void]$lines.Add('  - ' + [string]$capability)
    }
    return ([string]::Join("`n", [string[]]$lines) + "`n")
}

function Get-GeneratorContent {
    param([string]$GeneratorId)

    $generator = $GeneratorById[$GeneratorId]
    switch ([string]$generator.type) {
        'installed-toolkit-manifest' { return Get-InstalledManifestContent }
        'empty-ledger' { return [string]$generator.content }
        default { throw "Unsupported generator type '$($generator.type)'." }
    }
}

function Get-EntryContent {
    param($Entry)

    if ([string]$Entry.kind -eq 'copied') {
        $sourcePath = Join-PortablePath $RepositoryRoot ([string]$Entry.source) "entry '$($Entry.id)' source"
        return [System.IO.File]::ReadAllText($sourcePath)
    }
    return Get-GeneratorContent ([string]$Entry.generator)
}

$script:PlanActions = New-Object System.Collections.ArrayList
function Add-PlanAction {
    param(
        [string]$Action,
        [string]$EntryId,
        [AllowNull()][string]$Source,
        [AllowNull()][string]$Generator,
        [string]$Destination,
        [string]$Ownership,
        [string]$Reason,
        [bool]$MutatesTarget
    )

    $row = [pscustomobject][ordered]@{
        sequence = $script:PlanActions.Count + 1
        action = $Action
        entryId = $EntryId
        source = if ([string]::IsNullOrWhiteSpace($Source)) { $null } else { $Source }
        generator = if ([string]::IsNullOrWhiteSpace($Generator)) { $null } else { $Generator }
        destination = $Destination
        ownership = $Ownership
        reason = $Reason
        mutatesTarget = $MutatesTarget
        oldToolkitVersion = if ([string]::IsNullOrWhiteSpace($InstalledToolkitVersion)) {
            $null
        } else {
            $InstalledToolkitVersion
        }
        newToolkitVersion = $ToolkitVersion
    }
    [void]$script:PlanActions.Add($row)
}

if ($BlockReason) {
    Add-PlanAction `
        'block' `
        $BlockEntryId `
        $null `
        $null `
        $BlockDestination `
        $BlockOwnership `
        $BlockReason `
        $false
} else {
    $managedAgents = $MigrationEntries[0]
    $projectAgents = $EntryById['project-agents']
    $managedAgentsPath = Join-PortablePath $ProjectRoot ([string]$managedAgents.destination) 'managed AGENTS destination'
    $projectAgentsPath = Join-PortablePath $ProjectRoot ([string]$projectAgents.destination) 'project AGENTS destination'
    $managedAgentsContent = Get-EntryContent $managedAgents
    $MigrationCreatesProjectAgents = $false
    if (([string]$managedAgents.migrationPolicy -eq 'preserve-existing-agents') -and
        (Test-Path -LiteralPath $managedAgentsPath -PathType Leaf)) {
        $existingAgents = [System.IO.File]::ReadAllText($managedAgentsPath)
        if ($existingAgents -cne $managedAgentsContent) {
            if (-not (Test-Path -LiteralPath $projectAgentsPath -PathType Leaf)) {
                Add-PlanAction `
                    'create' `
                    ([string]$projectAgents.id) `
                    ([string]$projectAgents.source) `
                    $null `
                    ([string]$projectAgents.destination) `
                    'project-owned' `
                    'migrate-existing-agents' `
                    $true
                $MigrationCreatesProjectAgents = $true
            } else {
                Add-PlanAction `
                    'backup' `
                    ([string]$projectAgents.id) `
                    ([string]$projectAgents.source) `
                    $null `
                    ([string]$projectAgents.destination) `
                    'project-owned' `
                    'preserve-existing-agents' `
                    $true
            }
        }
    }

    :entryLoop foreach ($entry in $Entries) {
        switch ([string]$entry.selectionPolicy) {
            'default' { }
            'initialize-only' {
                if (-not $InitializeProjectStructure) { continue entryLoop }
            }
            default { throw "Unsupported selectionPolicy '$($entry.selectionPolicy)'." }
        }
        if ($MigrationCreatesProjectAgents -and
            ([string]$entry.id -eq [string]$projectAgents.id)) {
            continue
        }

        $destination = [string]$entry.destination
        $destinationPath = Join-PortablePath $ProjectRoot $destination "entry '$($entry.id)' destination"
        $source = if ([string]$entry.kind -eq 'copied') { [string]$entry.source } else { $null }
        $generator = if ([string]$entry.kind -eq 'generated') { [string]$entry.generator } else { $null }
        $desiredContent = Get-EntryContent $entry

        if (-not (Test-Path -LiteralPath $destinationPath -PathType Leaf)) {
            switch ([string]$entry.installPolicy) {
                'always' { }
                'missing-only' { }
                default { throw "Unsupported installPolicy '$($entry.installPolicy)'." }
            }
            $action = if ([string]$entry.kind -eq 'generated') { 'generate' } else { 'create' }
            $reason = if ([string]$entry.kind -eq 'generated') {
                'missing-generated-target'
            } else {
                'missing-target'
            }
            Add-PlanAction `
                $action `
                ([string]$entry.id) `
                $source `
                $generator `
                $destination `
                ([string]$entry.ownership) `
                $reason `
                $true
            continue
        }

        $existingContent = [System.IO.File]::ReadAllText($destinationPath)
        if ($existingContent -ceq $desiredContent) {
            Add-PlanAction `
                'unchanged' `
                ([string]$entry.id) `
                $source `
                $generator `
                $destination `
                ([string]$entry.ownership) `
                'content-matches' `
                $false
            continue
        }
        if ([string]$entry.installPolicy -eq 'missing-only') {
            Add-PlanAction `
                'preserve' `
                ([string]$entry.id) `
                $source `
                $generator `
                $destination `
                ([string]$entry.ownership) `
                'missing-only-content' `
                $false
            continue
        }
        if ([string]$entry.installPolicy -cne 'always') {
            throw "Unsupported installPolicy '$($entry.installPolicy)'."
        }

        $forceRefresh = switch ([string]$entry.forcePolicy) {
            'replace-managed' { [bool]$ForceManagedFiles }
            'preserve' { $false }
            default { throw "Unsupported forcePolicy '$($entry.forcePolicy)'." }
        }

        $refreshManaged = switch ([string]$entry.refreshPolicy) {
            'always' { $true }
            'upgrade-or-force' { $IsUpgrade -or $forceRefresh }
            'never' { $false }
            default { throw "Unsupported refreshPolicy '$($entry.refreshPolicy)'." }
        }
        if (-not $refreshManaged) {
            Add-PlanAction `
                'preserve' `
                ([string]$entry.id) `
                $source `
                $generator `
                $destination `
                'toolkit-managed' `
                'managed-content-differs' `
                $false
            continue
        }
        switch ([string]$entry.backupPolicy) {
            'before-replace' {
                Add-PlanAction `
                    'backup' `
                    ([string]$entry.id) `
                    $source `
                    $generator `
                    $destination `
                    ([string]$entry.ownership) `
                    'backup-before-replace' `
                    $true
            }
            'migration-aware' {
                if ([string]$entry.migrationPolicy -cne 'preserve-existing-agents') {
                    throw "Entry '$($entry.id)' has unsupported migration-aware backup semantics."
                }
            }
            'none' { }
            default { throw "Unsupported backupPolicy '$($entry.backupPolicy)'." }
        }
        $refreshAction = if ([string]$entry.kind -eq 'generated') { 'generate' } else { 'replace' }
        $refreshReason = if ([string]$entry.kind -eq 'generated') {
            'refresh-generated-content'
        } else {
            'refresh-managed-content'
        }
        Add-PlanAction `
            $refreshAction `
            ([string]$entry.id) `
            $source `
            $generator `
            $destination `
            'toolkit-managed' `
            $refreshReason `
            $true
    }
}

$Plan = [pscustomobject][ordered]@{
    schemaVersion = '1.0'
    manifestSchemaVersion = [string]$InstallManifest.schemaVersion
    mode = 'plan'
    toolkitVersion = $ToolkitVersion
    installedToolkitVersion = if ([string]::IsNullOrWhiteSpace($InstalledToolkitVersion)) {
        $null
    } else {
        $InstalledToolkitVersion
    }
    options = [pscustomobject][ordered]@{
        initializeProjectStructure = [bool]$InitializeProjectStructure
        forceManagedFiles = [bool]$ForceManagedFiles
    }
    canApply = -not [bool]$BlockReason
    actions = [object[]]$script:PlanActions
}

if ($PlanJson) {
    Write-Output ($Plan | ConvertTo-Json -Depth 20)
    return
}

if ($BlockReason) {
    Write-Warning "SDP installation blocked: $BlockReason"
    throw 'Refusing to modify a project that failed installation preflight.'
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

function Write-ProjectContent {
    param([string]$Relative, [string]$Content)

    $destination = Join-PortablePath $ProjectRoot $Relative "action destination '$Relative'"
    $parent = Split-Path -Parent $destination
    New-Item -ItemType Directory -Force -Path $parent | Out-Null
    [System.IO.File]::WriteAllText(
        $destination,
        $Content,
        [System.Text.UTF8Encoding]::new($false)
    )
}

function Backup-ProjectFile {
    param([string]$Relative)

    $source = Join-PortablePath $ProjectRoot $Relative "backup source '$Relative'"
    if (-not (Test-Path -LiteralPath $source -PathType Leaf)) { return }
    $destination = $BackupRoot
    foreach ($segment in $Relative.Split('/')) {
        $destination = Join-Path $destination $segment
    }
    $parent = Split-Path -Parent $destination
    New-Item -ItemType Directory -Force -Path $parent | Out-Null
    Copy-Item -LiteralPath $source -Destination $destination -Force
    Write-SdpAction 'APPLIED' "Backed up $Relative"
}

if ($Preview) {
    foreach ($action in $Plan.actions) {
        switch ([string]$action.action) {
            'preserve' { Write-SdpAction 'PRESERVED' "$($action.destination) ($($action.reason))" }
            'unchanged' { Write-SdpAction 'UNCHANGED' $action.destination }
            'warn' { Write-SdpAction 'WARNING' "$($action.destination) ($($action.reason))" }
            default { Write-SdpAction 'PROPOSED' "$($action.action) $($action.destination) ($($action.reason))" }
        }
    }
} else {
    $managedAgentsDestination = [string]$EntryById['managed-agents'].destination
    foreach ($action in $Plan.actions) {
        if ([string]$action.reason -eq 'migrate-existing-agents') {
            $source = Join-PortablePath $ProjectRoot $managedAgentsDestination 'AGENTS migration source'
            $destination = Join-PortablePath $ProjectRoot ([string]$action.destination) 'AGENTS migration destination'
            Copy-Item -LiteralPath $source -Destination $destination
            Write-SdpAction 'APPLIED' "Migrated existing AGENTS.md to $($action.destination)"
            continue
        }
        if ([string]$action.reason -eq 'preserve-existing-agents') {
            $source = Join-PortablePath $ProjectRoot $managedAgentsDestination 'AGENTS preservation source'
            $migrationRelative = "AGENTS-project.migration-$RunStamp.md"
            $destination = Join-PortablePath $ProjectRoot $migrationRelative 'AGENTS preservation destination'
            Copy-Item -LiteralPath $source -Destination $destination
            Write-SdpAction 'APPLIED' "Preserved existing AGENTS.md as $migrationRelative"
            continue
        }

        switch ([string]$action.action) {
            'backup' {
                Backup-ProjectFile ([string]$action.destination)
            }
            'create' {
                $sourcePath = Join-PortablePath $RepositoryRoot ([string]$action.source) "entry '$($action.entryId)' source"
                Write-ProjectContent ([string]$action.destination) ([System.IO.File]::ReadAllText($sourcePath))
                Write-SdpAction 'APPLIED' "Created $($action.destination) ($($action.ownership))"
            }
            'replace' {
                $sourcePath = Join-PortablePath $RepositoryRoot ([string]$action.source) "entry '$($action.entryId)' source"
                Write-ProjectContent ([string]$action.destination) ([System.IO.File]::ReadAllText($sourcePath))
                Write-SdpAction 'APPLIED' "Replaced $($action.destination)"
            }
            'generate' {
                Write-ProjectContent ([string]$action.destination) (Get-GeneratorContent ([string]$action.generator))
                Write-SdpAction 'APPLIED' "Generated $($action.destination)"
            }
            'preserve' {
                Write-SdpAction 'PRESERVED' "$($action.destination) ($($action.reason))"
            }
            'unchanged' {
                Write-SdpAction 'UNCHANGED' $action.destination
            }
            'warn' {
                Write-SdpAction 'WARNING' "$($action.destination) ($($action.reason))"
            }
            'block' {
                throw "Unexpected blocked action during apply: $($action.reason)"
            }
            default {
                throw "Unsupported planned action '$($action.action)'."
            }
        }
    }
}

Write-Host ''
Write-Host 'SDP installation summary'
Write-Host "  Toolkit version: $ToolkitVersion"
Write-Host "  Manifest schema: $($InstallManifest.schemaVersion)"
Write-Host "  Mode: $(if ($Preview) { 'preview' } else { 'apply' })"
Write-Host "  Upgrade: $IsUpgrade"
Write-Host "  Proposed: $script:ProposedCount"
Write-Host "  Applied: $script:AppliedCount"
Write-Host "  Preserved: $script:PreservedCount"
Write-Host "  Unchanged: $script:UnchangedCount"
Write-Host "  Backup root: $BackupRoot"
Write-Host 'No release or Git tag was created.'
