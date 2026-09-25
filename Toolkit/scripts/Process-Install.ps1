# Versioned profile operations within the Install-SDP PowerShell engine.
# Shared path, physical topology, hashing and SemVer helpers are loaded by entrypoint.

function ConvertTo-ProcessJson($Value) {
    return (ConvertTo-Json -InputObject $Value -Depth 100 -Compress) + "`n"
}
function ConvertTo-ProcessBytes([string]$Value) { return ,[Text.Encoding]::UTF8.GetBytes($Value) }
function ConvertTo-ProcessBase64([string]$Value) { return [Convert]::ToBase64String((ConvertTo-ProcessBytes $Value)) }
function Get-ContentHash([string]$Content) {
    return Get-Sha256Hex ([Convert]::FromBase64String($Content))
}
function Read-ProcessJson([string]$Text) {
    $doc = [Text.Json.JsonDocument]::Parse($Text)
    function Check-Keys($node) {
        if ($node.ValueKind -eq [Text.Json.JsonValueKind]::Object) {
            $keys = [Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
            foreach ($p in $node.EnumerateObject()) {
                if (-not $keys.Add($p.Name)) { throw "Duplicate/case-colliding JSON key: $($p.Name)" }
                Check-Keys $p.Value
            }
        } elseif ($node.ValueKind -eq [Text.Json.JsonValueKind]::Array) {
            foreach ($v in $node.EnumerateArray()) { Check-Keys $v }
        }
    }
    try { Check-Keys $doc.RootElement } finally { $doc.Dispose() }
    return ConvertFrom-Json -InputObject $Text -AsHashtable -Depth 100
}
function Assert-ProcessKeys($Object, [string[]]$Keys) {
    if ($Object -isnot [Collections.IDictionary] -or $Object.Count -ne $Keys.Count) { throw 'Invalid process object shape' }
    foreach ($key in $Keys) { if (-not $Object.Contains($key)) { throw "Missing field $key" } }
}
function Get-ProcessPath([string]$Relative) {
    Assert-PortableRelativePath $Relative 'profile destination' -Destination
    $path = Join-PortablePath $script:ProcessRoot $Relative 'profile destination'
    Assert-DestinationTopology $script:ProcessRoot $path 'profile destination'
    return $path
}
function Get-ProcessFiles {
    $map = [ordered]@{}
    $paths = [Collections.Generic.List[string]]::new()
    foreach ($name in @('SDP', '.codex/skills')) {
        $p = Join-Path $script:ProcessRoot $name
        if (Test-Path -LiteralPath $p) {
            Assert-NoLinkOrReparsePointInExistingAncestors $p 'snapshot root'
            $pending = [Collections.Generic.Queue[string]]::new()
            $pending.Enqueue($p)
            while ($pending.Count) {
                $dir = $pending.Dequeue()
                foreach ($item in Get-ChildItem -LiteralPath $dir -Force) {
                    $rel = [IO.Path]::GetRelativePath($script:ProcessRoot, $item.FullName).Replace('\','/')
                    if ($rel -in @('SDP/.sdp-operations','SDP/.sdp-backups')) {
                        Assert-NoLinkOrReparsePointInExistingAncestors $item.FullName 'operation directory'
                        continue
                    }
                    if ($item.Attributes -band [IO.FileAttributes]::ReparsePoint) { throw "Linked object in process area: $rel" }
                    if ($item.PSIsContainer) { $pending.Enqueue($item.FullName) } else { $paths.Add($rel) }
                }
            }
        }
    }
    foreach ($item in Get-ChildItem -LiteralPath $script:ProcessRoot -File -Force) {
        if ($item.Extension -eq '.md') { $paths.Add($item.Name) }
    }
    foreach ($rel in ($paths | Sort-Object -CaseSensitive)) {
        if ($map.Contains($rel)) { throw "Case-colliding file: $rel" }
        $p = Get-ProcessPath $rel
        if ((Get-PathObjectState $p) -cne 'regular-file') { throw "Unsupported file: $rel" }
        $map[$rel] = [Convert]::ToBase64String([IO.File]::ReadAllBytes($p))
    }
    return $map
}
function Get-ProcessSnapshot($Files) {
    $map = [ordered]@{}
    foreach ($key in ($Files.Keys | Sort-Object -CaseSensitive)) { $map[$key] = Get-ContentHash $Files[$key] }
    return $map
}
function Read-ProcessArtifact {
    $path = [IO.Path]::GetFullPath($ProfileArtifact)
    Assert-NoLinkOrReparsePointInExistingAncestors $path 'artifact'
    $text = [IO.File]::ReadAllText($path)
    $a = Read-ProcessJson $text
    Assert-ProcessKeys $a @('schemaVersion','profile','managementProfile','prerequisites','capabilities','files','relocations','facts','sourceDigest','configurationDigest')
    if ($a.schemaVersion -cne '2.0' -or $a.profile -cne 'sdp-five-phase/0.1' -or $a.managementProfile -cne 'sdp-project-management/0.1') { throw 'Unsupported artifact profile/schema' }
    if ($a.configurationDigest -cnotmatch '^[a-f0-9]{64}$') { throw 'Invalid configuration digest' }
    $withoutDigest = $text.Replace('"configurationDigest":"' + $a.configurationDigest + '",','')
    if ((Get-Sha256Hex (ConvertTo-ProcessBytes $withoutDigest)) -cne $a.configurationDigest) { throw 'Artifact digest mismatch (use canonical build output)' }
    if (($a.prerequisites -join ',') -cne 'powershell>=7.4' -or $PSVersionTable.PSVersion -lt [version]'7.4') { throw 'PowerShell 7.4 or newer required' }
    Assert-ProcessKeys $a.facts @('toolkitVersion','frameworkVersion','agentsContractVersion','installerVersion','skills','capabilities')
    foreach ($key in @('toolkitVersion','frameworkVersion','agentsContractVersion','installerVersion')) { [void](ConvertTo-SemVer $a.facts[$key]) }
    $seen = @{}
    foreach ($entry in $a.files) {
        Assert-ProcessKeys $entry @('source','destination','ownership','sha256','content')
        Assert-PortableRelativePath $entry.source 'artifact source'
        [void](Get-ProcessPath $entry.destination)
        if ($entry.ownership -cnotin @('managed','project') -or (Get-ContentHash $entry.content) -cne $entry.sha256) { throw 'Invalid artifact file' }
        foreach ($key in $seen.Keys) {
            if ($entry.destination -ieq $key -or $entry.destination.StartsWith($key+'/',[StringComparison]::OrdinalIgnoreCase) -or $key.StartsWith($entry.destination+'/',[StringComparison]::OrdinalIgnoreCase)) { throw 'Overlapping artifact destinations' }
        }
        $seen[$entry.destination] = $true
    }
    foreach ($move in $a.relocations) {
        Assert-ProcessKeys $move @('from','to')
        Assert-PortableRelativePath $move.from 'relocation source' -Destination
        Assert-PortableRelativePath $move.to 'relocation destination' -Destination
        Assert-ContainedPhysicalPath $script:ProcessRoot (Join-Path $script:ProcessRoot $move.from) 'relocation source'
        Assert-ContainedPhysicalPath $script:ProcessRoot (Join-Path $script:ProcessRoot $move.to) 'relocation destination'
        if (-not $move.from.StartsWith('SDP/') -or -not $move.to.StartsWith('SDP/')) { throw 'Unsupported relocation root' }
    }
    return $a
}
function ConvertFrom-ProcessBase64([string]$Value) { return [Text.Encoding]::UTF8.GetString([Convert]::FromBase64String($Value)) }

function Get-ProcessPlan($Artifact) {
    $before = Get-ProcessFiles
    $after = [ordered]@{}
    foreach ($key in $before.Keys) { $after[$key] = $before[$key] }
    $conflicts = [Collections.Generic.List[string]]::new()
    $warnings = [Collections.Generic.List[string]]::new()
    $moves = [ordered]@{}
    $old = [ordered]@{toolkitVersion='unknown';processProfile='unknown';schemaVersion='unknown'}
    $factsPath = 'SDP/Framework/installed-toolkit.manifest.yaml'
    $oldFacts = $null
    if ($before.Contains($factsPath)) {
        $oldFacts = ConvertFrom-InstalledManifestDocument (ConvertFrom-ProcessBase64 $before[$factsPath]) -AllowV2
        $version = Get-StrictYamlString $oldFacts 'schemaVersion' 'installed facts'
        if ($version -cnotin @('1.0','2.0')) { $conflicts.Add('unsupported-installed-schema') }
        $old.schemaVersion = $version
        $old.toolkitVersion = Get-StrictYamlString $oldFacts 'toolkitVersion' 'installed facts'
        if ((Compare-SemVer $old.toolkitVersion $Artifact.facts.toolkitVersion) -gt 0) { $conflicts.Add('downgrade-blocked') }
        if ($version -ceq '2.0') {
            $old.processProfile = Get-StrictYamlString $oldFacts 'processProfile' 'installed facts'
            if ($old.processProfile -cne $Artifact.profile) { $conflicts.Add('unsupported-process-transition') }
        }
    }
    $legacy = @($before.Keys | Where-Object { $_ -cmatch '^SDP/(02--Study|03--Requirements|04--Architecture|05--DesignAnalysis|06--Design|07--Implementation)/' }).Count -gt 0
    $modern = @($before.Keys | Where-Object { $_ -cmatch '^SDP/(02--Requirements|03--Architecture|04--Design|05--Implementation)/' }).Count -gt 0
    $baseline = if ($legacy) { 'legacy-seven-phase' } elseif ($before.Contains('SDP/Agents/KanBan/board.json')) { 'manual-five-phase-board' } elseif ($modern) { 'local-five-phase' } elseif (@($before.Keys | Where-Object { $_.StartsWith('SDP/') }).Count -eq 0) { 'clean' } else { 'unclassified' }
    if ($legacy -and $modern) { $conflicts.Add('mixed-phase-layout') }
    if ($baseline -eq 'unclassified' -and $null -eq $oldFacts) { $conflicts.Add('unclassified-baseline') }
    foreach ($move in $Artifact.relocations) {
        $sources = @($before.Keys | Where-Object { $_.StartsWith($move.from+'/',[StringComparison]::Ordinal) })
        if (-not $sources.Count) { continue }
        $destinations = @($before.Keys | Where-Object { $_.StartsWith($move.to+'/',[StringComparison]::OrdinalIgnoreCase) })
        if ($destinations.Count) { $conflicts.Add('relocation-destination-exists: '+$move.to); continue }
        foreach ($source in $sources) {
            $dest = $move.to + $source.Substring($move.from.Length)
            [void](Get-ProcessPath $dest)
            $after[$dest] = $before[$source]
            $after.Remove($source)
            $moves[$source] = $dest
        }
    }
    # Transfer one known board history without editing historical event bytes.
    $boardPath = 'SDP/KanBan/board.json'
    $ledgerPath = 'SDP/ProjectManagement/Ledger.ndjson'
    if ($after.Contains($boardPath)) {
        $board = Read-ProcessJson (ConvertFrom-ProcessBase64 $after[$boardPath])
        if ($board.schemaVersion -ceq '0.1' -and $board.ledger -ceq 'Ledger.ndjson') {
            $source = 'SDP/KanBan/Ledger.ndjson'
            if ($after.Contains($ledgerPath)) { $conflicts.Add('multiple-management-histories') }
            elseif (-not $after.Contains($source)) { $conflicts.Add('missing-board-history') }
            else {
                $after[$ledgerPath] = $after[$source]
                $after.Remove($source)
                $original = if ($before.Contains('SDP/Agents/KanBan/Ledger.ndjson')) { 'SDP/Agents/KanBan/Ledger.ndjson' } else { $source }
                $moves[$original] = $ledgerPath
                $board = [ordered]@{schemaVersion='0.2';projectId=$board.projectId;namespaces=@($board.projectId);ledger='../ProjectManagement/Ledger.ndjson';profile=$Artifact.managementProfile}
                $after[$boardPath] = ConvertTo-ProcessBase64 (ConvertTo-ProcessJson $board)
            }
        } elseif ($board.schemaVersion -cne '0.2' -or $board.profile -cne $Artifact.managementProfile -or $board.ledger -cne '../ProjectManagement/Ledger.ndjson') { $conflicts.Add('unsupported-board-history-contract') }
        $projectID = $board.projectId
    } else {
        $projectID = 'PROJECT'
        $after[$boardPath] = ConvertTo-ProcessBase64 (ConvertTo-ProcessJson ([ordered]@{schemaVersion='0.2';projectId=$projectID;namespaces=@($projectID);ledger='../ProjectManagement/Ledger.ndjson';profile=$Artifact.managementProfile}))
    }
    if ($projectID -cnotmatch '^[A-Z][A-Z0-9]*$') { $conflicts.Add('invalid-management-project-id') }
    if (-not $after.Contains($ledgerPath)) { $after[$ledgerPath] = '' }
    $history = ConvertFrom-ProcessBase64 $after[$ledgerPath]
    if ($history -ne '' -and -not $history.EndsWith("`n")) { $conflicts.Add('history-missing-final-newline') }
    $eventIDs = @{}
    foreach ($line in $history.Split("`n")) {
        if (-not $line.Trim()) { continue }
        $e = Read-ProcessJson $line
        if ($e.schemaVersion -cne '1.0' -or -not $e.eventId -or $eventIDs.ContainsKey($e.eventId)) { $conflicts.Add('invalid-history-envelope'); break }
        if ($e.eventType.StartsWith('x-kanban:') -and $e.payload.schemaVersion -cnotin @('0.1','0.2')) { $conflicts.Add('unsupported-card-history') }
        if ($e.eventType.StartsWith('x-management:') -and $e.payload.schemaVersion -cne '0.1') { $conflicts.Add('unsupported-management-history') }
        $eventIDs[$e.eventId] = $true
    }
    # Rebase supported inline Markdown links from original to final locations.
    foreach ($source in @($before.Keys)) {
        if (-not $source.EndsWith('.md',[StringComparison]::OrdinalIgnoreCase)) { continue }
        $dest = if ($moves.Contains($source)) { $moves[$source] } else { $source }
        if (-not $after.Contains($dest)) { continue }
        $text = ConvertFrom-ProcessBase64 $after[$dest]
        $rootForLinks = $script:ProcessRoot
        $rebase = {
            param($match)
            $url = $match.Groups[1].Value
            if ($url -match '^[a-zA-Z][a-zA-Z0-9+.-]*:' -or $url.StartsWith('#') -or $url.StartsWith('/')) { return $match.Value }
            $parts = $url.Split('#',2)
            $raw = [Uri]::UnescapeDataString($parts[0])
            if (-not $raw -or $raw.Contains(' ')) { return $match.Value }
            $abs = [IO.Path]::GetFullPath((Join-Path ([IO.Path]::GetDirectoryName((Join-Path $rootForLinks $source))) $raw))
            $target = [IO.Path]::GetRelativePath($rootForLinks,$abs).Replace('\','/')
            $newTarget = $target
            if ($moves.Contains($target)) { $newTarget = $moves[$target] }
            elseif ($moves.Count) {
                foreach ($move in $Artifact.relocations) {
                    if ($target -ceq $move.from -or $target.StartsWith($move.from+'/')) { $newTarget = $move.to+$target.Substring($move.from.Length); break }
                }
            }
            if ($dest -ceq $source -and $target -ceq $newTarget) { return $match.Value }
            $relative = [IO.Path]::GetRelativePath([IO.Path]::GetDirectoryName((Join-Path $rootForLinks $dest)),(Join-Path $rootForLinks $newTarget)).Replace('\','/')
            $encoded = ($relative.Split('/') | ForEach-Object { [Uri]::EscapeDataString($_) }) -join '/'
            if ($parts.Length -gt 1) { $encoded += '#'+$parts[1] }
            return ']('+$encoded+')'
        }.GetNewClosure()
        $updated = [regex]::Replace($text,'\]\(([^\s)]+)\)',[Text.RegularExpressions.MatchEvaluator]$rebase)
        if ($updated -cne $text) { $after[$dest] = ConvertTo-ProcessBase64 $updated }
    }
    if ($moves.Count) { $warnings.Add('Review preserved non-Markdown references and project scripts for relocated paths; semantic source repartition is not performed.') }
    foreach ($entry in $Artifact.files) {
        $dest = $entry.destination
        if (-not $after.Contains($dest)) { $after[$dest] = $entry.content }
        elseif ($entry.ownership -ceq 'managed' -and $after[$dest] -cne $entry.content) {
            if ($dest -ceq 'AGENTS.md' -and -not (ConvertFrom-ProcessBase64 $after[$dest]).Contains('Toolkit-managed by the System Design Process')) {
                if ($after.Contains('AGENTS-project.md') -and $after['AGENTS-project.md'] -cne $after[$dest]) { $conflicts.Add('agents-preservation-conflict') }
                else { $after['AGENTS-project.md'] = $after[$dest] }
            }
            if ($ForceManagedFiles) { $after[$dest] = $entry.content } else { $conflicts.Add('managed-refresh-requires-force: '+$dest) }
        }
    }
    if (-not $after.Contains('SDP/Traceability/Ledger.ndjson')) { $after['SDP/Traceability/Ledger.ndjson'] = '' }
    $navPath = 'SDP/navigation.json'
    if (-not $after.Contains($navPath)) {
        $after[$navPath] = ConvertTo-ProcessBase64 (ConvertTo-ProcessJson ([ordered]@{schemaVersion='1.0';projectId=$projectID.ToLowerInvariant();processProfile=$Artifact.profile;projectManifest='SDP/SDP-project.manifest.yaml';kanban='SDP/KanBan';models=@();sdui=@()}))
    } elseif ($moves.Contains('SDP/Agents/KanBan/board.json')) {
        $nav = Read-ProcessJson (ConvertFrom-ProcessBase64 $after[$navPath])
        if ($nav.Contains('kanban') -and $nav.kanban -ceq 'SDP/Agents/KanBan') {
            $nav.kanban='SDP/KanBan'
            $after[$navPath] = ConvertTo-ProcessBase64 (ConvertTo-ProcessJson $nav)
        }
    }
    $facts = [Collections.Generic.List[string]]::new()
    $facts.Add('schemaVersion: "2.0"')
    foreach ($key in @('toolkitVersion','frameworkVersion','agentsContractVersion','installerVersion')) { $facts.Add($key+': '+(ConvertTo-YamlQuotedScalar $Artifact.facts[$key])) }
    $installedAt = '@operation-time@'
    if ($null -ne $oldFacts -and $oldFacts.ActualPaths -ccontains 'configurationDigest' -and (Get-StrictYamlString $oldFacts 'configurationDigest' 'facts') -ceq $Artifact.configurationDigest) {
        $installedAt = Get-StrictYamlString $oldFacts 'toolkitInstalledAt' 'facts'
    }
    $facts.Add('toolkitInstalledAt: '+(ConvertTo-YamlQuotedScalar $installedAt))
    $facts.Add('sourceCommit: null')
    $facts.Add('processProfile: '+(ConvertTo-YamlQuotedScalar $Artifact.profile))
    $facts.Add('managementProfile: '+(ConvertTo-YamlQuotedScalar $Artifact.managementProfile))
    $facts.Add('configurationDigest: '+(ConvertTo-YamlQuotedScalar $Artifact.configurationDigest))
    $facts.Add('skills:')
    foreach ($key in ($Artifact.facts.skills.Keys | Sort-Object -CaseSensitive)) { $facts.Add('  '+$key+': '+(ConvertTo-YamlQuotedScalar $Artifact.facts.skills[$key])) }
    $facts.Add('capabilities:')
    foreach ($cap in $Artifact.facts.capabilities) { $facts.Add('  - '+$cap) }
    $after[$factsPath] = ConvertTo-ProcessBase64 (($facts -join "`n")+"`n")
    $actions = [Collections.Generic.List[object]]::new()
    $preserved = [Collections.Generic.List[string]]::new()
    foreach ($dest in ($after.Keys | Sort-Object -CaseSensitive)) {
        $hash = if ($before.Contains($dest)) { Get-ContentHash $before[$dest] } else { $null }
        if ($before.Contains($dest) -and $before[$dest] -ceq $after[$dest]) { $preserved.Add($dest); continue }
        $actions.Add([ordered]@{action='write';destination=$dest;before=$hash;after=(Get-ContentHash $after[$dest]);content=$after[$dest]})
    }
    foreach ($dest in ($before.Keys | Sort-Object -CaseSensitive)) {
        if (-not $after.Contains($dest)) { $actions.Add([ordered]@{action='delete';destination=$dest;before=(Get-ContentHash $before[$dest]);after=$null;content=$null}) }
    }
    return [ordered]@{schemaVersion='2.0';projectRoot=$script:ProcessRoot;configurationDigest=$Artifact.configurationDigest;forceManagedFiles=[bool]$ForceManagedFiles;baseline=$baseline;oldFacts=$old;targetProfile=$Artifact.profile;managementProject=$projectID;canApply=($conflicts.Count -eq 0);conflicts=@($conflicts);warnings=@($warnings);snapshot=(Get-ProcessSnapshot $before);actions=@($actions);preserved=@($preserved)}
}

function Invoke-ProcessInstallation {
    $script:ProcessRoot = Get-ProviderCompatibleFullPath $ProjectRoot 'project root'
    if (-not (Test-Path -LiteralPath $script:ProcessRoot -PathType Container)) { throw 'Project root must exist' }
    Assert-NoLinkOrReparsePointInExistingAncestors $script:ProcessRoot 'project root'
    Assert-NoPhysicalTreeOverlap (Get-PhysicalPathState $RepositoryRoot 'source' -RequireExistingRoot) (Get-PhysicalPathState $script:ProcessRoot 'project' -RequireExistingRoot) 'Source and project must be separate'
    if ($InitializeProjectStructure -or $BackupRoot) { throw 'Profile installation always includes its layout and uses operation-local backups' }
    $artifact = Read-ProcessArtifact
    if ($ResumeOperation) { throw 'Resume implementation pending' }
    $plan = Get-ProcessPlan $artifact
    if ($PlanJson -or $Preview) { Write-Output (ConvertTo-ProcessJson $plan); return }
    throw 'Use -PlanJson to inspect; profile apply implementation pending'
}
