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
    # Incoming Markdown links can live outside SDP. Exclude only declared
    # administrative/dependency/build trees; never follow external symlinks.
    $queue = [Collections.Generic.Queue[string]]::new()
    $queue.Enqueue($script:ProcessRoot)
    while ($queue.Count) {
        $dir = $queue.Dequeue()
        foreach ($item in Get-ChildItem -LiteralPath $dir -Force) {
            $rel = [IO.Path]::GetRelativePath($script:ProcessRoot,$item.FullName).Replace('\','/')
            if ($rel -ceq 'SDP' -or $rel -ceq '.codex/skills') { continue }
            if ($item.Name -in @('.git','node_modules','.venv','vendor','build','.cache')) { continue }
            if ($item.Attributes -band [IO.FileAttributes]::ReparsePoint) { continue }
            if ($item.PSIsContainer) { $queue.Enqueue($item.FullName) }
            elseif ($item.Extension -ieq '.md') { $paths.Add($rel) }
        }
    }
    $directories = @{}
    foreach ($rel in ($paths | Sort-Object -CaseSensitive)) {
        if ($map.Contains($rel)) { throw "Case-colliding file: $rel" }
        $segments = $rel.Split('/')
        for ($i=1; $i -lt $segments.Length; $i++) {
            $prefix = ($segments[0..($i-1)] -join '/')
            $key = $prefix.ToLowerInvariant()
            if ($directories.ContainsKey($key) -and $directories[$key] -cne $prefix) { throw 'Case-colliding directory' }
            $directories[$key] = $prefix
        }
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
function Assert-ProcessReserved([string]$Path, [switch]$BoardMapping) {
    $protected = @('SDP/.sdp-operations','SDP/.sdp-backups','SDP/Framework/installed-toolkit.manifest.yaml','SDP/ProjectManagement/Ledger.ndjson','SDP/navigation.json','SDP/Traceability/Ledger.ndjson')
    if (-not $BoardMapping) { $protected += 'SDP/KanBan/board.json' }
    foreach ($p in $protected) {
        if ($Path.Equals($p,[StringComparison]::OrdinalIgnoreCase) -or $Path.StartsWith($p+'/',[StringComparison]::OrdinalIgnoreCase) -or $p.StartsWith($Path+'/',[StringComparison]::OrdinalIgnoreCase)) { throw 'Engine-owned path or ancestor/descendant in artifact' }
    }
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
        $d = $entry.destination
        Assert-ProcessReserved $d
        if ($d -cne 'AGENTS.md' -and -not $d.StartsWith('SDP/') -and -not $d.StartsWith('.codex/skills/')) { throw 'Artifact destination outside installation scope' }
        if ($entry.ownership -ceq 'managed' -and $d -cne 'AGENTS.md' -and -not $d.StartsWith('.codex/skills/') -and -not $d.StartsWith('SDP/Framework/')) { throw 'Invalid managed destination' }
        if ($entry.ownership -cnotin @('managed','project') -or (Get-ContentHash $entry.content) -cne $entry.sha256) { throw 'Invalid artifact file' }
        foreach ($key in $seen.Keys) {
            if ($entry.destination -ieq $key -or $entry.destination.StartsWith($key+'/',[StringComparison]::OrdinalIgnoreCase) -or $key.StartsWith($entry.destination+'/',[StringComparison]::OrdinalIgnoreCase)) { throw 'Overlapping artifact destinations' }
        }
        $seen[$entry.destination] = $true
    }
    foreach ($move in $a.relocations) {
        Assert-ProcessKeys $move @('from','to')
        Assert-ProcessReserved $move.from
        Assert-ProcessReserved $move.to -BoardMapping:($move.from -ceq 'SDP/Agents/KanBan' -and $move.to -ceq 'SDP/KanBan')
        Assert-PortableRelativePath $move.from 'relocation source' -Destination
        Assert-PortableRelativePath $move.to 'relocation destination' -Destination
        Assert-ContainedPhysicalPath $script:ProcessRoot (Join-Path $script:ProcessRoot $move.from) 'relocation source'
        Assert-ContainedPhysicalPath $script:ProcessRoot (Join-Path $script:ProcessRoot $move.to) 'relocation destination'
        if (-not $move.from.StartsWith('SDP/') -or -not $move.to.StartsWith('SDP/')) { throw 'Unsupported relocation root' }
    }
    return $a
}
function ConvertFrom-ProcessBase64([string]$Value) { return [Text.Encoding]::UTF8.GetString([Convert]::FromBase64String($Value)) }

function Assert-ProcessHistory($Artifact, [string]$History, $Files) {
    $schemas = @{}
    foreach ($entry in $Artifact.files) {
        if ($entry.destination -cmatch '^SDP/ProjectManagement/(.+\.schema\.json)$') {
            $schemas[$Matches[1]] = ConvertFrom-ProcessBase64 $entry.content
        }
    }
    $ids = @{}
    $latest = @{}
    foreach ($line in $History.Split("`n")) {
        if (-not $line.Trim()) { continue }
        $event = Read-ProcessJson $line
        if (-not (Test-Json -Json $line -Schema $schemas['ledger-event.schema.json'] -ErrorAction SilentlyContinue)) { throw 'Invalid ledger envelope' }
        if ($ids.ContainsKey($event.eventId)) { throw 'Duplicate ledger event identity' }
        $ids[$event.eventId] = $true
        $payload = $event.payload
        if ($event.eventType.StartsWith('x-kanban:')) {
            if ($payload.schemaVersion -cnotin @('0.1','0.2')) { throw 'Unsupported KanBan payload' }
            $schema = 'kanban-payload-'+$payload.schemaVersion+'.schema.json'
            if ($event.subjectId -cnotmatch '^KB-[A-Z][A-Z0-9]*-[0-9]{3,}$') { throw 'Invalid card identity' }
        } elseif ($event.eventType.StartsWith('x-management:')) {
            $schema = 'management-payload.schema.json'
        } else { throw 'Unsupported event type in shared management history' }
        if (-not (Test-Json -Json (ConvertTo-ProcessJson $payload) -Schema $schemas[$schema] -ErrorAction SilentlyContinue)) { throw 'Invalid management/card payload' }
        $prior = $latest[$event.subjectId]
        if ($null -eq $prior) {
            if ($null -ne $payload.previousEventId -or $null -ne $payload.from -or $null -ne $payload.fromPath) { throw 'Broken initial history chain' }
        } elseif ($payload.previousEventId -cne $prior.eventId -or $payload.from -cne $prior.payload.to -or $payload.fromPath -cne $prior.payload.toPath) { throw 'Broken history chain' }
        $latest[$event.subjectId] = $event
    }
    $cards = @{}
    $stateFolders = @{backlog=@('backlog','queued');active=@('ready','in-progress','gate-review');onHold=@('onHold');completed=@('completed');canceled=@('canceled');superseded=@('superseded');irrelevant=@('irrelevant')}
    foreach ($path in $Files.Keys) {
        if ($path -cnotmatch '^SDP/KanBan/(backlog|active|onHold|completed|canceled|superseded|irrelevant)/#.+\.md$') { continue }
        $folder = $Matches[1]
        $text = ConvertFrom-ProcessBase64 $Files[$path]
        $states = [regex]::Matches($text,'(?m)^\| CardState \| ([^|]+) \|\s*$')
        $ident = [regex]::Matches($text,'(?m)^\| id \| (KB-[A-Z][A-Z0-9]*-[0-9]{3,}) \|\s*$')
        if ($states.Count -ne 1 -or $ident.Count -ne 1 -or $stateFolders[$folder] -cnotcontains $states[0].Groups[1].Value.Trim()) { throw 'Invalid card metadata/state' }
        $id = $ident[0].Groups[1].Value
        if ($cards.ContainsKey($id) -or -not $latest.ContainsKey($id)) { throw 'Duplicate card or missing history' }
        $e = $latest[$id]
        if ($e.payload.to -cne $folder -or $e.payload.toPath -cne $path.Substring('SDP/KanBan/'.Length)) { throw 'Card/history location mismatch' }
        $cards[$id]=$true
    }
    foreach ($key in $latest.Keys) {
        if ($key.StartsWith('KB-') -and -not $cards.ContainsKey($key)) { throw 'History points to missing card' }
    }
}

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
    $projectManifest = 'SDP/SDP-project.manifest.yaml'
    if ($before.Contains($projectManifest)) {
        $pm = ConvertFrom-StrictYamlDocument (ConvertFrom-ProcessBase64 $before[$projectManifest]) 'project manifest'
        if ((Get-StrictYamlString $pm 'schemaVersion' 'project manifest') -cne '1.0') { $conflicts.Add('unsupported-project-schema') }
        if ((Get-StrictYamlString $pm 'installed/manifestPath' 'project manifest') -cne 'Framework/installed-toolkit.manifest.yaml') { $conflicts.Add('custom-installed-facts-path-needs-mapping') }
    }
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
    if (-not $after.Contains($ledgerPath)) {
        if ($before.Contains($boardPath) -or $before.Contains('SDP/Agents/KanBan/board.json')) { $conflicts.Add('missing-shared-history') }
        $after[$ledgerPath] = ''
    }
    $history = ConvertFrom-ProcessBase64 $after[$ledgerPath]
    if ($history -ne '' -and -not $history.EndsWith("`n")) { $conflicts.Add('history-missing-final-newline') }
    try { Assert-ProcessHistory $Artifact $history $after } catch { $conflicts.Add('invalid-management-history: '+$_.Exception.Message) }
    # Rebase supported inline Markdown links from original to final locations.
    foreach ($source in @($before.Keys)) {
        if (-not $source.EndsWith('.md',[StringComparison]::OrdinalIgnoreCase)) { continue }
        $dest = if ($moves.Contains($source)) { $moves[$source] } else { $source }
        if (-not $after.Contains($dest)) { continue }
        $text = ConvertFrom-ProcessBase64 $after[$dest]
        $rootForLinks = $script:ProcessRoot
        $rebase = {
            param($match)
            $urlGroup = $match.Groups['url']
            $url = $urlGroup.Value
            if ($url -match '^[a-zA-Z][a-zA-Z0-9+.-]*:' -or $url.StartsWith('#') -or $url.StartsWith('/')) { return $match.Value }
            $parts = $url.Split('#',2)
            $raw = [Uri]::UnescapeDataString($parts[0]).Replace('\(','(').Replace('\)',')')
            if (-not $raw) { return $match.Value }
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
            $offset = $urlGroup.Index - $match.Index
            return $match.Value.Substring(0,$offset)+$encoded+$match.Value.Substring($offset+$urlGroup.Length)
        }.GetNewClosure()
        $inline = '\]\(\s*(?:<(?<url>[^>\r\n]+)>|(?<url>(?:[^\s()\\]+|\\.|(?<paren>\()|(?<-paren>\)))+)(?(paren)(?!)))(?:[ \t]+(?:"[^"\r\n]*"|''[^''\r\n]*''|\([^\r\n)]*\)))?\s*\)'
        $reference = '(?m)^ {0,3}\[[^\]\r\n]+\]:[ \t]*(?:<(?<url>[^>\r\n]+)>|(?<url>[^\s]+))'
        $updated = [regex]::Replace($text,$inline,[Text.RegularExpressions.MatchEvaluator]$rebase)
        $updated = [regex]::Replace($updated,$reference,[Text.RegularExpressions.MatchEvaluator]$rebase)
        if ($updated -cne $text) { $after[$dest] = ConvertTo-ProcessBase64 $updated }
    }
    if ($moves.Count) { $warnings.Add('Review non-Markdown references, project scripts and excluded dependency/build/symlink trees for relocated paths; semantic source repartition is not performed.') }
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

function Get-ProcessFileHash([string]$Relative) {
    $path = Get-ProcessPath $Relative
    if (-not (Test-Path -LiteralPath $path)) { return $null }
    if ((Get-PathObjectState $path) -cne 'regular-file') { throw "Unsupported destination: $Relative" }
    return Get-Sha256Hex ([IO.File]::ReadAllBytes($path))
}
function Write-ProcessAtomic([string]$Relative, [byte[]]$Bytes) {
    $dest = Get-ProcessPath $Relative
    [void][IO.Directory]::CreateDirectory([IO.Path]::GetDirectoryName($dest))
    [void](Get-ProcessPath $Relative)
    $tempKey = Get-Sha256Hex (ConvertTo-ProcessBytes $Relative)
    $temp = Get-ProcessPath ('SDP/.sdp-operations/temp-'+$tempKey)
    # Internal staging is excluded from the project snapshot. Under the exclusive
    # lock, regenerate a partial staging file from the immutable journal payload.
    if (Test-Path -LiteralPath $temp) {
        if ((Get-PathObjectState $temp) -cne 'regular-file') { throw 'Unsafe temporary object' }
        [IO.File]::Delete($temp)
    }
    $stream = [IO.FileStream]::new($temp,[IO.FileMode]::CreateNew,[IO.FileAccess]::Write,[IO.FileShare]::None)
    try { $stream.Write($Bytes,0,$Bytes.Length); $stream.Flush($true) } finally { $stream.Dispose() }
    [void](Get-ProcessPath $Relative)
    [IO.File]::Move($temp,$dest,$true)
}
function Save-ProcessJournal($Journal) {
    $Journal.integrity = ''
    $Journal.integrity = Get-Sha256Hex (ConvertTo-ProcessBytes (ConvertTo-ProcessJson $Journal))
    Write-ProcessAtomic ('SDP/.sdp-operations/'+$Journal.operationId+'/journal.json') (ConvertTo-ProcessBytes (ConvertTo-ProcessJson $Journal))
}
function Read-ProcessJournal([string]$ID) {
    if ($ID -cnotmatch '^install-[a-f0-9]{24}$') { throw 'Invalid operation ID' }
    $path = Get-ProcessPath ('SDP/.sdp-operations/'+$ID+'/journal.json')
    $j = Read-ProcessJson ([IO.File]::ReadAllText($path))
    $integrity = $j.integrity
    $j.integrity = ''
    if ((Get-Sha256Hex (ConvertTo-ProcessBytes (ConvertTo-ProcessJson $j))) -cne $integrity) { throw 'Corrupt operation journal; preserve it and restore the verified journal backup' }
    $j.integrity = $integrity
    if ($j.operationId -cne $ID -or $j.plan.projectRoot -cne $script:ProcessRoot -or $j.next -lt 0 -or $j.next -gt $j.steps.Count) { throw 'Invalid journal identity/progress' }
    return $j
}
function Get-PendingProcessOperations {
    $root = Join-Path $script:ProcessRoot 'SDP/.sdp-operations'
    if (-not (Test-Path -LiteralPath $root)) { return @() }
    Assert-NoLinkOrReparsePointInExistingAncestors $root 'operations'
    $pending = [Collections.Generic.List[string]]::new()
    foreach ($dir in Get-ChildItem -LiteralPath $root -Directory -Force) {
        Assert-NoLinkOrReparsePointInExistingAncestors $dir.FullName 'operation'
        if (-not (Test-Path -LiteralPath (Join-Path $dir.FullName 'journal.json'))) {
            if (@(Get-ChildItem -LiteralPath $dir.FullName -Force).Count) { throw 'Incomplete preparation contains unknown files' }
            continue
        }
        $j = Read-ProcessJournal $dir.Name
        if ($j.status -cne 'completed') { $pending.Add($dir.Name) }
    }
    return @($pending)
}
function Invoke-ProcessFault([string]$Boundary, [int]$Index) {
    # Deterministic fault injection for disposable conformance fixtures.
    if ($env:SDP_INSTALL_INTERRUPT -ceq ($Boundary+':'+$Index)) {
        if ($env:SDP_INSTALL_HARD_EXIT -ceq '1') { [Environment]::Exit(97) }
        throw "Injected interruption at ${Boundary}:${Index}"
    }
}
function New-ProcessJournal($Plan, $Artifact) {
    $id = 'install-'+(Get-Sha256Hex (ConvertTo-ProcessBytes (ConvertTo-ProcessJson $Plan))).Substring(0,24)
    $when = [DateTimeOffset]::UtcNow.ToString('yyyy-MM-ddTHH:mm:ssZ')
    $steps = [Collections.Generic.List[object]]::new()
    $factsStep = $null
    $ledgerPath = 'SDP/ProjectManagement/Ledger.ndjson'
    $history = ''
    if ($Plan.snapshot.Contains($ledgerPath)) { $history = [IO.File]::ReadAllText((Get-ProcessPath $ledgerPath)) }
    foreach ($action in $Plan.actions) {
        if ($action.destination -ceq $ledgerPath -and $action.action -ceq 'write') {
            $history = ConvertFrom-ProcessBase64 $action.content
            continue
        }
        if ($action.destination -ceq 'SDP/Framework/installed-toolkit.manifest.yaml') {
            $text = (ConvertFrom-ProcessBase64 $action.content).Replace('@operation-time@',$when)
            $content = ConvertTo-ProcessBase64 $text
            $factsStep = [ordered]@{action='write';destination=$action.destination;before=$action.before;after=(Get-ContentHash $content);content=$content}
            continue
        }
        $steps.Add($action)
    }
    $usedMaintenance = $history
    foreach ($file in (Get-ProcessFiles).GetEnumerator()) {
        if ($file.Key.StartsWith('SDP/Maintenance/') -and $file.Key.EndsWith('.md')) { $usedMaintenance += ConvertFrom-ProcessBase64 $file.Value }
    }
    $n = 1
    while (($usedMaintenance -match ('\bMAINT-'+$Plan.managementProject+'-'+('{0:D4}' -f $n)+'\b')) -or (Test-Path -LiteralPath (Join-Path $script:ProcessRoot ('SDP/Maintenance/MAINT-'+$Plan.managementProject+'-'+('{0:D4}' -f $n)+'.md')))) { $n++ }
    $maint = 'MAINT-'+$Plan.managementProject+'-'+('{0:D4}' -f $n)
    $reportPath = 'Maintenance/'+$maint+'.md'
    $eventNumber = 1
    $pattern = '"eventId"\s*:\s*"EVT-PM-'+[regex]::Escape($Plan.managementProject)+'-([0-9]+)"'
    foreach ($match in [regex]::Matches($history,$pattern)) { $eventNumber = [Math]::Max($eventNumber, [int]$match.Groups[1].Value+1) }
    $eventFirst = 'EVT-PM-'+$Plan.managementProject+'-'+('{0:D6}' -f $eventNumber)
    $eventLast = 'EVT-PM-'+$Plan.managementProject+'-'+('{0:D6}' -f ($eventNumber+1))
    foreach ($index in @(1,2)) {
        $payload = [ordered]@{schemaVersion='0.1';kind='Maintenance';previousEventId=$(if($index -eq 1){$null}else{$eventFirst});from=$(if($index -eq 1){$null}else{'active'});to=$(if($index -eq 1){'active'}else{'completed'});fromPath=$(if($index -eq 1){$null}else{$reportPath});toPath=$reportPath;reason=$(if($index -eq 1){'Execute reviewed installation plan '+$id}else{'Verified profile installation '+$Artifact.profile+'; configuration '+$Artifact.configurationDigest});links=@($id)}
        $event = [ordered]@{schemaVersion='1.0';eventId=$(if($index -eq 1){$eventFirst}else{$eventLast});eventType=$(if($index -eq 1){'x-management:created'}else{'x-management:completed'});occurredAt=$when;actor='sdp-installer';commit=$null;subjectId=$maint;payload=$payload}
        $history += ConvertTo-ProcessJson $event
    }
    $report = "# $maint — SDP process installation`n`n| Field | Value |`n| --- | --- |`n| id | $maint |`n| project | $($Plan.managementProject) |`n| state | completed |`n| operation | $id |`n| source | Reviewed installation plan |`n`n"
    $report += "Old Toolkit: $($Plan.oldFacts.toolkitVersion). Old process profile: $($Plan.oldFacts.processProfile).`nObserved layout: $($Plan.baseline). Target Toolkit: $($Artifact.facts.toolkitVersion).`nNew profile: $($Artifact.profile). Management: $($Artifact.managementProfile).`nConfiguration SHA-256: $($Artifact.configurationDigest).`n`n"
    $report += "Operation journal and byte backups: SDP/.sdp-operations/$id.`nEvery recorded before/after hash is verified. This report is finalized only when`nthe operation journal says completed; an earlier interruption is incomplete.`nThe operation uses forward resume, not whole-tree rollback.`n`n## Changed paths`n`n"
    foreach ($a in $Plan.actions) { $report += '- '+$a.action+' '+$a.destination+"`n" }
    $report += "`n## Preserved paths`n`n"
    foreach ($p in $Plan.preserved) { $report += '- '+$p+"`n" }
    $report += "`n## Capabilities`n`n"
    foreach ($cap in $Artifact.facts.capabilities) { $report += '- '+$cap+"`n" }
    if ($Plan.warnings.Count) { $report += "`n## Operator review`n`n"+($Plan.warnings -join "`n")+"`n" }
    $content = ConvertTo-ProcessBase64 $report
    $steps.Add([ordered]@{action='write';destination='SDP/'+$reportPath;before=$null;after=(Get-ContentHash $content);content=$content})
    $content = ConvertTo-ProcessBase64 $history
    $steps.Add([ordered]@{action='write';destination=$ledgerPath;before=$(if($Plan.snapshot.Contains($ledgerPath)){$Plan.snapshot[$ledgerPath]}else{$null});after=(Get-ContentHash $content);content=$content})
    if ($null -ne $factsStep) { $steps.Add($factsStep) }
    return [ordered]@{schemaVersion='2.0';operationId=$id;configurationDigest=$Artifact.configurationDigest;plan=$Plan;createdAt=$when;maintenanceId=$maint;report=$reportPath;steps=@($steps);next=0;status='active';error=$null;integrity=''}
}
function Assert-ProcessProgress($Journal) {
    $expected = [ordered]@{}
    foreach ($key in $Journal.plan.snapshot.Keys) { $expected[$key] = $Journal.plan.snapshot[$key] }
    for ($i=0; $i -lt $Journal.next; $i++) {
        $step = $Journal.steps[$i]
        if ($step.action -ceq 'delete') { $expected.Remove($step.destination) } else { $expected[$step.destination]=$step.after }
    }
    # A crash after an atomic replacement but before progress publication is recoverable.
    if ($Journal.next -lt $Journal.steps.Count) {
        $step = $Journal.steps[$Journal.next]
        $hash = Get-ProcessFileHash $step.destination
        if ($hash -ceq $step.after) {
            if ($step.action -ceq 'delete') { $expected.Remove($step.destination) } else { $expected[$step.destination]=$step.after }
        } elseif ($hash -cne $step.before) { throw "Post-failure destination edit: $($step.destination)" }
    }
    $actual = Get-ProcessSnapshot (Get-ProcessFiles)
    if ($actual.Count -ne $expected.Count) { throw 'Project changed after plan/interruption (file inventory)' }
    foreach ($key in $expected.Keys) {
        if (-not $actual.Contains($key) -or $actual[$key] -cne $expected[$key]) { throw "Project changed after plan/interruption: $key" }
    }
}
function Invoke-ProcessJournal($Journal) {
    try {
        Assert-ProcessProgress $Journal
        for ($i=[int]$Journal.next; $i -lt $Journal.steps.Count; $i++) {
            $step = $Journal.steps[$i]
            $path = Get-ProcessPath $step.destination
            $actual = Get-ProcessFileHash $step.destination
            if ($actual -cne $step.before -and $actual -cne $step.after) { throw "Destination drift: $($step.destination)" }
            $backup = 'SDP/.sdp-operations/'+$Journal.operationId+'/backups/'+$i
            if ($null -ne $step.before) {
                if (Test-Path -LiteralPath (Get-ProcessPath $backup)) {
                    if ((Get-ProcessFileHash $backup) -cne $step.before) { throw 'Backup integrity mismatch' }
                } elseif ($actual -ceq $step.before) {
                    Write-ProcessAtomic $backup ([IO.File]::ReadAllBytes($path))
                } else { throw 'Missing recovery backup' }
            }
            Invoke-ProcessFault 'backup' $i
            if ($actual -cne $step.after) {
                # Recheck the precondition after backup and immediately before mutation.
                if ((Get-ProcessFileHash $step.destination) -cne $step.before) { throw 'Destination changed during backup' }
                if ($step.action -ceq 'delete') { [IO.File]::Delete($path) }
                elseif ($step.action -ceq 'write') { Write-ProcessAtomic $step.destination ([Convert]::FromBase64String($step.content)) }
                else { throw 'Unsupported journal action' }
            }
            Invoke-ProcessFault 'write' $i
            if ((Get-ProcessFileHash $step.destination) -cne $step.after) { throw 'Written content failed verification' }
            $Journal.next = $i+1
            $Journal.status = 'active'
            $Journal.error = $null
            Save-ProcessJournal $Journal
            Invoke-ProcessFault 'journal' $i
        }
        Assert-ProcessProgress $Journal
        Invoke-ProcessFault 'complete' $Journal.next
        $Journal.status='completed'
        Save-ProcessJournal $Journal
    } catch {
        $Journal.status='failed'
        $Journal.error=$_.Exception.Message
        Save-ProcessJournal $Journal
        throw
    }
    Write-Output (ConvertTo-ProcessJson ([ordered]@{schemaVersion='2.0';status='completed';operationId=$Journal.operationId;maintenanceId=$Journal.maintenanceId;report='SDP/'+$Journal.report}))
}
function Invoke-ProcessInstallation {
    $script:ProcessRoot = Get-ProviderCompatibleFullPath $ProjectRoot 'project root'
    if (-not (Test-Path -LiteralPath $script:ProcessRoot -PathType Container)) { throw 'Project root must exist' }
    Assert-NoLinkOrReparsePointInExistingAncestors $script:ProcessRoot 'project root'
    Assert-NoPhysicalTreeOverlap (Get-PhysicalPathState $RepositoryRoot 'source' -RequireExistingRoot) (Get-PhysicalPathState $script:ProcessRoot 'project' -RequireExistingRoot) 'Source and project must be separate'
    if ($InitializeProjectStructure -or $BackupRoot) { throw 'Profile installation always includes its layout and uses operation-local backups' }
    $artifact = Read-ProcessArtifact
    if ($PlanJson -or $Preview) {
        if ($ApplyPlan -or $ResumeOperation) { throw 'Preview cannot apply or resume' }
        $plan = Get-ProcessPlan $artifact
        $pending = @(Get-PendingProcessOperations)
        if ($pending.Count) { $plan.canApply=$false; $plan.conflicts += @('incomplete-operation: '+($pending -join ',')) }
        Write-Output (ConvertTo-ProcessJson $plan)
        return
    }
    if ((-not $ApplyPlan) -eq (-not $ResumeOperation)) { throw 'Select exactly one of ApplyPlan or ResumeOperation after reviewing PlanJson' }
    # Exclusive OS handle is released on normal exit or process death.
    $lockPath = Get-ProcessPath 'SDP/.sdp-operations/install.lock'
    [void][IO.Directory]::CreateDirectory([IO.Path]::GetDirectoryName($lockPath))
    $lock = [IO.FileStream]::new($lockPath,[IO.FileMode]::OpenOrCreate,[IO.FileAccess]::ReadWrite,[IO.FileShare]::None)
    try {
        $artifact = Read-ProcessArtifact
        $pending = @(Get-PendingProcessOperations)
        if ($ResumeOperation) {
            if ($pending.Count -and ($pending.Count -ne 1 -or $pending[0] -cne $ResumeOperation)) { throw 'Another incomplete operation exists' }
            $journal = Read-ProcessJournal $ResumeOperation
            if ($journal.configurationDigest -cne $artifact.configurationDigest) { throw 'Resume requires the original artifact' }
            if ($journal.status -ceq 'completed') {
                Write-Output (ConvertTo-ProcessJson ([ordered]@{schemaVersion='2.0';status='already-completed';operationId=$ResumeOperation}))
                return
            }
        } else {
            if ($pending.Count) { throw ('Resume incomplete operation first: '+($pending -join ',')) }
            $supplied = Read-ProcessJson ([IO.File]::ReadAllText([IO.Path]::GetFullPath($ApplyPlan)))
            $plan = Get-ProcessPlan $artifact
            if ((ConvertTo-ProcessJson $supplied) -cne (ConvertTo-ProcessJson $plan)) { throw 'Reviewed plan changed: re-plan source/target/force settings' }
            if (-not $plan.canApply) { throw ('Plan has conflicts: '+($plan.conflicts -join ', ')) }
            if (-not $plan.actions.Count) {
                Write-Output (ConvertTo-ProcessJson ([ordered]@{schemaVersion='2.0';status='no-change'}))
                return
            }
            $journal = New-ProcessJournal $plan $artifact
            $journalPath = Get-ProcessPath ('SDP/.sdp-operations/'+$journal.operationId+'/journal.json')
            if (Test-Path -LiteralPath $journalPath) { throw 'Operation identity already exists; use resume' }
            Save-ProcessJournal $journal
            Invoke-ProcessFault 'prepared' 0
        }
        Invoke-ProcessJournal $journal
    } finally { $lock.Dispose() }
}
