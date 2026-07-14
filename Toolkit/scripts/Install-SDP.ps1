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
$IsWindowsPlatform = [System.IO.Path]::DirectorySeparatorChar -eq '\'
if ($IsWindowsPlatform -and (-not ('Sdp.Install.NativePath' -as [type]))) {
    Add-Type -TypeDefinition @'
using System;
using System.ComponentModel;
using System.Globalization;
using System.Runtime.InteropServices;
using System.Text;
using Microsoft.Win32.SafeHandles;

namespace Sdp.Install {
    public static class NativePath {
        [StructLayout(LayoutKind.Sequential)]
        private struct FileIdInformation {
            public ulong VolumeSerialNumber;

            [MarshalAs(UnmanagedType.ByValArray, SizeConst = 16)]
            public byte[] FileId;
        }

        [DllImport("kernel32.dll", CharSet = CharSet.Unicode, SetLastError = true,
            EntryPoint = "GetLongPathNameW")]
        public static extern uint GetLongPathName(
            string shortPath,
            StringBuilder longPath,
            uint bufferLength
        );

        [DllImport("kernel32.dll", CharSet = CharSet.Unicode, SetLastError = true,
            EntryPoint = "CreateFileW")]
        private static extern SafeFileHandle CreateFile(
            string fileName,
            uint desiredAccess,
            uint shareMode,
            IntPtr securityAttributes,
            uint creationDisposition,
            uint flagsAndAttributes,
            IntPtr templateFile
        );

        [DllImport("kernel32.dll", SetLastError = true)]
        [return: MarshalAs(UnmanagedType.Bool)]
        private static extern bool GetFileInformationByHandleEx(
            SafeFileHandle file,
            int informationClass,
            out FileIdInformation information,
            uint bufferSize
        );

        public static string GetFileSystemIdentity(string path) {
            const uint ShareReadWriteDelete = 0x00000001 | 0x00000002 | 0x00000004;
            const uint OpenExisting = 3;
            const uint FileFlagBackupSemantics = 0x02000000;

            using (SafeFileHandle handle = CreateFile(
                path,
                0,
                ShareReadWriteDelete,
                IntPtr.Zero,
                OpenExisting,
                FileFlagBackupSemantics,
                IntPtr.Zero
            )) {
                if (handle.IsInvalid) {
                    throw new Win32Exception(
                        Marshal.GetLastWin32Error(),
                        "Cannot open path for physical identity"
                    );
                }
                FileIdInformation information;
                uint bufferSize = (uint)Marshal.SizeOf(typeof(FileIdInformation));
                const int FileIdInfo = 18;
                if (!GetFileInformationByHandleEx(
                    handle,
                    FileIdInfo,
                    out information,
                    bufferSize
                )) {
                    throw new Win32Exception(
                        Marshal.GetLastWin32Error(),
                        "Cannot read physical path identity"
                    );
                }
                if (information.FileId == null || information.FileId.Length != 16) {
                    throw new InvalidOperationException("The operating system returned an invalid file ID");
                }
                return information.VolumeSerialNumber.ToString("X16", CultureInfo.InvariantCulture)
                    + ":"
                    + BitConverter.ToString(information.FileId).Replace("-", "");
            }
        }
    }
}
'@
}

function Assert-JsonObjectShape {
    param(
        $Value,
        [string[]]$Required,
        [string[]]$Allowed,
        [string]$Label
    )

    if (($null -eq $Value) -or ($Value -isnot [pscustomobject])) {
        throw "$Label must be a JSON object."
    }
    $names = @($Value.psobject.Properties | ForEach-Object { $_.Name })
    foreach ($name in $Required) {
        if ($names -cnotcontains $name) {
            throw "$Label is missing required property '$name'."
        }
    }
    foreach ($name in $names) {
        if ($Allowed -cnotcontains $name) {
            throw "$Label contains unknown property '$name'."
        }
    }
}

function Assert-JsonArray {
    param($Value, [int]$MinimumCount, [string]$Label)

    if ($Value -isnot [System.Array]) {
        throw "$Label must be a JSON array."
    }
    if (@($Value).Count -lt $MinimumCount) {
        throw "$Label must contain at least $MinimumCount item(s)."
    }
}

function Assert-JsonString {
    param($Value, [string]$Label)

    if (($Value -isnot [string]) -or [string]::IsNullOrWhiteSpace($Value)) {
        throw "$Label must be a non-empty JSON string."
    }
}

function ConvertFrom-StrictYamlScalar {
    param([string]$Text, [string]$Label)

    if ([string]::IsNullOrWhiteSpace($Text)) {
        throw "$Label must be a scalar value."
    }
    if ($Text -cmatch '^"(?:[^"\\]|\\["\\/bfnrt]|\\u[0-9A-Fa-f]{4})*"$') {
        try {
            # PowerShell 6+ coerces standalone RFC 3339 JSON strings to dates.
            # A non-date prefix preserves JSON unescaping without changing the value.
            $stringPrefix = 'sdp-yaml-string:'
            $prefixedText = '"' + $stringPrefix + $Text.Substring(1)
            $decoded = ConvertFrom-Json -InputObject $prefixedText
            if (($decoded -isnot [string]) -or
                (-not $decoded.StartsWith($stringPrefix, [System.StringComparison]::Ordinal))) {
                throw "$Label must decode to a string."
            }
            return [pscustomobject]@{
                Type = 'string'
                Value = $decoded.Substring($stringPrefix.Length)
            }
        } catch {
            throw "$Label contains an invalid double-quoted scalar."
        }
    }
    if ($Text -cmatch "^'(?:[^']|'')*'$" ) {
        return [pscustomobject]@{
            Type = 'string'
            Value = $Text.Substring(1, $Text.Length - 2).Replace("''", "'")
        }
    }
    switch -CaseSensitive ($Text) {
        'null' { return [pscustomobject]@{ Type = 'null'; Value = $null } }
        'true' { return [pscustomobject]@{ Type = 'boolean'; Value = $true } }
        'false' { return [pscustomobject]@{ Type = 'boolean'; Value = $false } }
        '[]' { return [pscustomobject]@{ Type = 'empty-sequence'; Value = @() } }
        '{}' { return [pscustomobject]@{ Type = 'empty-mapping'; Value = @{} } }
    }
    if ($Text -cnotmatch '^[A-Za-z0-9][A-Za-z0-9._/+@:-]*(?: [A-Za-z0-9][A-Za-z0-9._/+@:-]*)*$') {
        throw "$Label is outside the supported strict YAML scalar subset."
    }
    return [pscustomobject]@{ Type = 'string'; Value = $Text }
}

function ConvertFrom-StrictYamlDocument {
    param(
        [string]$Content,
        [string]$Label,
        [string[]]$RejectNestedKeys = @()
    )

    if ([string]::IsNullOrWhiteSpace($Content)) {
        throw "$Label is empty."
    }
    if ($Content.Contains("`t")) {
        throw "$Label contains a tab, which is outside the supported YAML subset."
    }

    $values = @{}
    $sequences = @{}
    $containers = @{}
    $seenPaths = @{}
    $actualPaths = New-Object System.Collections.ArrayList
    $containerPathByIndent = @{ 0 = '' }
    $containerKindByIndent = @{ 0 = 'mapping' }
    $previousIndent = 0
    $pendingContainerPath = $null
    $sawToken = $false
    $lineNumber = 0

    foreach ($rawLine in [regex]::Split($Content, "\r?\n")) {
        $lineNumber++
        if (($rawLine -cmatch '^\s*$') -or ($rawLine -cmatch '^\s*#')) {
            continue
        }
        if ($rawLine -cmatch '^\s*(?:---|\.\.\.|%)') {
            throw "$Label line $lineNumber uses an unsupported YAML document feature."
        }
        $indent = $rawLine.Length - $rawLine.TrimStart(' ').Length
        if (($indent % 2) -ne 0) {
            throw "$Label line $lineNumber must use two-space indentation."
        }
        $body = $rawLine.Substring($indent)
        $sequenceMatch = [regex]::Match($body, '^-\s+(.+)$')
        $mappingMatch = [regex]::Match($body, '^([A-Za-z0-9][A-Za-z0-9._-]*):(?:\s*(.*))?$')
        $isSequence = $sequenceMatch.Success
        $isMapping = $mappingMatch.Success
        if (-not ($isSequence -xor $isMapping)) {
            throw "$Label line $lineNumber is outside the supported mapping/scalar-sequence YAML subset."
        }

        if (-not $sawToken) {
            if ($indent -ne 0 -or $isSequence) {
                throw "$Label must be a root mapping."
            }
            $sawToken = $true
        } elseif ($indent -gt $previousIndent) {
            if (($indent -ne ($previousIndent + 2)) -or
                [string]::IsNullOrWhiteSpace([string]$pendingContainerPath)) {
                throw "$Label line $lineNumber has invalid indentation."
            }
            $containerPathByIndent[$indent] = $pendingContainerPath
            $containerKindByIndent[$indent] = if ($isSequence) { 'sequence' } else { 'mapping' }
            $containers[$pendingContainerPath] = $containerKindByIndent[$indent]
        } else {
            foreach ($knownIndent in @($containerPathByIndent.Keys)) {
                if ([int]$knownIndent -gt $indent) {
                    $containerPathByIndent.Remove($knownIndent)
                    $containerKindByIndent.Remove($knownIndent)
                }
            }
            if (-not $containerPathByIndent.ContainsKey($indent)) {
                throw "$Label line $lineNumber has invalid indentation."
            }
        }

        $expectedKind = [string]$containerKindByIndent[$indent]
        if (($isSequence -and $expectedKind -cne 'sequence') -or
            ($isMapping -and $expectedKind -cne 'mapping')) {
            throw "$Label line $lineNumber mixes mapping and sequence entries."
        }
        $containerPath = [string]$containerPathByIndent[$indent]
        $pendingContainerPath = $null

        if ($isSequence) {
            $scalar = ConvertFrom-StrictYamlScalar $sequenceMatch.Groups[1].Value "$Label line $lineNumber"
            if (-not $sequences.ContainsKey($containerPath)) {
                $sequences[$containerPath] = New-Object System.Collections.ArrayList
            }
            [void]$sequences[$containerPath].Add($scalar)
        } else {
            $key = [string]$mappingMatch.Groups[1].Value
            $text = [string]$mappingMatch.Groups[2].Value
            if (($indent -gt 0) -and ($RejectNestedKeys -contains $key)) {
                throw "$Label contains nested shadow key '$key' at line $lineNumber."
            }
            $path = if ([string]::IsNullOrWhiteSpace($containerPath)) {
                $key
            } else {
                "$containerPath/$key"
            }
            $pathKey = $path.ToLowerInvariant()
            if ($seenPaths.ContainsKey($pathKey)) {
                throw "$Label contains duplicate key '$key' at line $lineNumber."
            }
            $seenPaths[$pathKey] = $true
            [void]$actualPaths.Add($path)
            if ([string]::IsNullOrWhiteSpace($text)) {
                $pendingContainerPath = $path
                $containers[$path] = 'pending'
            } else {
                $values[$path] = ConvertFrom-StrictYamlScalar $text "$Label line $lineNumber key '$key'"
            }
        }
        $previousIndent = $indent
    }

    if (-not $sawToken) {
        throw "$Label contains no mapping entries."
    }
    return [pscustomobject]@{
        Values = $values
        Sequences = $sequences
        Containers = $containers
        Paths = @($seenPaths.Keys)
        ActualPaths = [string[]]$actualPaths
    }
}

function Get-StrictYamlString {
    param($Document, [string]$Path, [string]$Label)

    if (-not $Document.Values.ContainsKey($Path)) {
        throw "$Label is missing required scalar '$Path'."
    }
    $scalar = $Document.Values[$Path]
    if ([string]$scalar.Type -cne 'string') {
        throw "$Label field '$Path' must be a string scalar."
    }
    return [string]$scalar.Value
}

function ConvertTo-SemVer {
    param([string]$Version)

    $identifier = '(?:0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*)'
    $pattern = '^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-(' +
        $identifier + '(?:\.' + $identifier + ')*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$'
    $match = [regex]::Match($Version, $pattern)
    if (-not $match.Success) {
        throw "Invalid Toolkit SemVer: $Version"
    }
    return [pscustomobject]@{
        Major = $match.Groups[1].Value
        Minor = $match.Groups[2].Value
        Patch = $match.Groups[3].Value
        PreRelease = if ($match.Groups[4].Success) {
            [string[]]$match.Groups[4].Value.Split('.')
        } else {
            [string[]]@()
        }
        Build = if ($match.Groups[5].Success) { $match.Groups[5].Value } else { $null }
    }
}

function Compare-NumericIdentifier {
    param([string]$Left, [string]$Right)

    if ($Left.Length -lt $Right.Length) { return -1 }
    if ($Left.Length -gt $Right.Length) { return 1 }
    return [string]::Compare($Left, $Right, [System.StringComparison]::Ordinal)
}

function Compare-SemVer {
    param([string]$Left, [string]$Right)

    $leftVersion = ConvertTo-SemVer $Left
    $rightVersion = ConvertTo-SemVer $Right
    foreach ($field in @('Major', 'Minor', 'Patch')) {
        $comparison = Compare-NumericIdentifier $leftVersion.$field $rightVersion.$field
        if ($comparison -ne 0) { return $comparison }
    }
    $leftPre = @($leftVersion.PreRelease)
    $rightPre = @($rightVersion.PreRelease)
    if (($leftPre.Count -eq 0) -and ($rightPre.Count -eq 0)) { return 0 }
    if ($leftPre.Count -eq 0) { return 1 }
    if ($rightPre.Count -eq 0) { return -1 }
    $count = [Math]::Min($leftPre.Count, $rightPre.Count)
    for ($index = 0; $index -lt $count; $index++) {
        $leftPart = [string]$leftPre[$index]
        $rightPart = [string]$rightPre[$index]
        $leftNumeric = $leftPart -cmatch '^[0-9]+$'
        $rightNumeric = $rightPart -cmatch '^[0-9]+$'
        if ($leftNumeric -and $rightNumeric) {
            $comparison = Compare-NumericIdentifier $leftPart $rightPart
        } elseif ($leftNumeric) {
            $comparison = -1
        } elseif ($rightNumeric) {
            $comparison = 1
        } else {
            $comparison = [string]::Compare(
                $leftPart,
                $rightPart,
                [System.StringComparison]::Ordinal
            )
        }
        if ($comparison -ne 0) { return $comparison }
    }
    if ($leftPre.Count -lt $rightPre.Count) { return -1 }
    if ($leftPre.Count -gt $rightPre.Count) { return 1 }
    return 0
}

function Assert-PortableRelativePath {
    param(
        [string]$Value,
        [string]$Label,
        [switch]$Destination
    )

    if ([string]::IsNullOrWhiteSpace($Value)) {
        throw "$Label must be a non-empty relative path."
    }
    if ($Value.Contains('\')) {
        throw "$Label must use portable '/' separators: $Value"
    }
    if ($Value.StartsWith('/') -or $Value -match '^[A-Za-z]:' -or $Value.StartsWith('//')) {
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
    foreach ($segment in $segments) {
        foreach ($character in $segment.ToCharArray()) {
            $codePoint = [int][char]$character
            if (($codePoint -lt 32) -or ($codePoint -eq 127)) {
                throw "$Label contains a control character."
            }
        }
        if ($segment.IndexOfAny([char[]]'<>:"|?*~') -ge 0) {
            throw "$Label contains a Windows-invalid segment '$segment'."
        }
        if ($segment.EndsWith('.') -or $segment.EndsWith(' ')) {
            throw "$Label contains a Windows-aliased segment '$segment'."
        }
        $deviceBase = $segment.Split('.')[0]
        if ($deviceBase -match '^(?i:CON|PRN|AUX|NUL|CLOCK\$|CONIN\$|CONOUT\$|COM[1-9\u00B9\u00B2\u00B3]|LPT[1-9\u00B9\u00B2\u00B3])$') {
            throw "$Label contains reserved device segment '$segment'."
        }
        if ($segment.Equals('.git', [System.StringComparison]::OrdinalIgnoreCase)) {
            throw "$Label must not address Git administrative paths."
        }
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

function Assert-CanonicalWindowsExtendedPath {
    param([string]$Path, [string]$Label)

    if (-not $IsWindowsPlatform) { return }

    $rawDevicePrefix = [regex]::Match($Path, '^[\\/]{2}([?.])[\\/]')
    if ($rawDevicePrefix.Success -and
        ($Path.Substring(0, 4) -cne '\\?\')) {
        if ($rawDevicePrefix.Groups[1].Value -ceq '?') {
            $normalizedNamespace = $Path.Substring(4).Replace('/', '\')
            if ($normalizedNamespace -cmatch '^[A-Za-z]:') {
                throw "$Label uses a normalization-sensitive extended drive path: the device prefix and separators must be '\\?\X:\'."
            }
            if ($normalizedNamespace.StartsWith(
                'UNC\',
                [System.StringComparison]::OrdinalIgnoreCase
            )) {
                throw "$Label uses a normalization-sensitive extended UNC path: the device prefix and separators must be '\\?\UNC\'."
            }
        }
        throw "$Label uses an unsupported Windows device namespace."
    }

    $extendedKind = $null
    $remainder = $null
    if ($Path.StartsWith('\\?\UNC\', [System.StringComparison]::OrdinalIgnoreCase)) {
        $extendedKind = 'UNC'
        $remainder = $Path.Substring(8)
    } elseif ($Path.StartsWith('\\?\UNC/', [System.StringComparison]::OrdinalIgnoreCase)) {
        throw "$Label uses a normalization-sensitive extended UNC path: '/' cannot replace the namespace separator."
    } elseif ($Path -cmatch '^\\\\\?\\[A-Za-z]:') {
        $extendedKind = 'drive'
        if (($Path.Length -lt 7) -or ($Path[6] -ne '\')) {
            throw "$Label uses a normalization-sensitive extended drive path: the drive prefix must end in '\'."
        }
        $remainder = $Path.Substring(7)
    } elseif ($Path.StartsWith('\\?\', [System.StringComparison]::OrdinalIgnoreCase) -or
        $Path.StartsWith('\\.\', [System.StringComparison]::OrdinalIgnoreCase) -or
        $Path.StartsWith('\??\', [System.StringComparison]::OrdinalIgnoreCase)) {
        throw "$Label uses an unsupported Windows device namespace."
    } else {
        return
    }

    if ($Path.IndexOf('/') -ge 0) {
        throw "$Label uses a normalization-sensitive extended $extendedKind path: '/' separators are not supported."
    }

    if (($extendedKind -ceq 'drive') -and ($remainder.Length -eq 0)) {
        return
    }

    $segments = @($remainder.Split([char[]]'\'))
    if (($extendedKind -ceq 'UNC') -and ($segments.Count -lt 2)) {
        throw "$Label uses a normalization-sensitive extended UNC path: both server and share segments are required."
    }
    foreach ($segment in $segments) {
        if ($segment.Length -eq 0) {
            throw "$Label uses a normalization-sensitive extended $extendedKind path: empty, doubled, or trailing separators are not supported."
        }
        if (($segment -ceq '.') -or ($segment -ceq '..')) {
            throw "$Label uses a normalization-sensitive extended $extendedKind path: traversal segment '$segment' is not supported."
        }
        if ($segment.EndsWith(' ') -or $segment.EndsWith('.')) {
            throw "$Label uses a normalization-sensitive extended $extendedKind path: segment '$segment' ends in a space or dot."
        }
    }
}

function Get-ProviderCompatibleFullPath {
    param([string]$Path, [string]$Label)

    if ([string]::IsNullOrWhiteSpace($Path)) {
        throw "$Label must be a non-empty path."
    }
    Assert-CanonicalWindowsExtendedPath $Path $Label
    $candidate = $Path
    if ($IsWindowsPlatform) {
        if ($candidate.StartsWith('\\?\UNC\', [System.StringComparison]::OrdinalIgnoreCase)) {
            $candidate = '\\' + $candidate.Substring(8)
        } elseif ($candidate -cmatch '^\\\\\?\\[A-Za-z]:\\') {
            $candidate = $candidate.Substring(4)
        } elseif ($candidate.StartsWith('\\?\', [System.StringComparison]::OrdinalIgnoreCase)) {
            throw "$Label uses an unsupported extended device namespace."
        }
    }
    try {
        return [System.IO.Path]::GetFullPath($candidate)
    } catch {
        throw "$Label is not a valid filesystem path: $($_.Exception.Message)"
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
    if ((-not $full.Equals($fullBase, $PathComparison)) -and
        (-not $full.StartsWith($prefix, $PathComparison))) {
        throw "$Label escapes its root: $Relative"
    }
    return $full
}

function Get-FileSystemParentPath {
    param([string]$Path)

    $full = [System.IO.Path]::GetFullPath($Path)
    $root = [System.IO.Path]::GetPathRoot($full)
    if ([string]::IsNullOrWhiteSpace($root)) { return $null }

    $trimCharacters = [char[]]@(
        [System.IO.Path]::DirectorySeparatorChar,
        [System.IO.Path]::AltDirectorySeparatorChar
    )
    $fullKey = $full.TrimEnd($trimCharacters)
    $rootKey = $root.TrimEnd($trimCharacters)
    if ($fullKey.Length -eq 0) { $fullKey = $full }
    if ($rootKey.Length -eq 0) { $rootKey = $root }
    if ($fullKey.Equals($rootKey, $PathComparison)) { return $null }

    $parentInfo = [System.IO.Directory]::GetParent($full)
    if ($null -eq $parentInfo) { return $null }
    $parent = [System.IO.Path]::GetFullPath($parentInfo.FullName)
    $parentRoot = [System.IO.Path]::GetPathRoot($parent)
    $parentRootKey = $parentRoot.TrimEnd($trimCharacters)
    if ($parentRootKey.Length -eq 0) { $parentRootKey = $parentRoot }
    if (-not $parentRootKey.Equals($rootKey, $PathComparison)) {
        throw 'Filesystem parent resolution crossed its namespace root.'
    }
    if ($parent.Equals($full, $PathComparison)) {
        throw 'Filesystem parent resolution did not advance.'
    }
    return $parent
}

function Test-IsLinkOrReparsePoint {
    param([string]$Path)

    try {
        $attributes = [System.IO.File]::GetAttributes($Path)
    } catch [System.IO.FileNotFoundException] {
        return $false
    } catch [System.IO.DirectoryNotFoundException] {
        return $false
    }
    if (($attributes -band [System.IO.FileAttributes]::ReparsePoint) -ne 0) {
        return $true
    }
    $item = Get-Item -LiteralPath $Path -Force -ErrorAction Stop
    if (($item.psobject.Properties.Name -contains 'LinkType') -and
        (-not [string]::IsNullOrWhiteSpace([string]$item.LinkType))) {
        return $true
    }
    return $false
}

function Assert-NoLinkOrReparsePointInExistingAncestors {
    param([string]$Path, [string]$Label)

    $current = [System.IO.Path]::GetFullPath($Path)
    while (-not [string]::IsNullOrWhiteSpace($current)) {
        if (Test-IsLinkOrReparsePoint $current) {
            throw "$Label traverses a symlink or reparse point: $current"
        }
        $parent = Get-FileSystemParentPath $current
        if ([string]::IsNullOrWhiteSpace($parent) -or
            $parent.Equals($current, $PathComparison)) {
            break
        }
        $current = $parent
    }
}

function Assert-ContainedPhysicalPath {
    param(
        [string]$Root,
        [string]$Path,
        [string]$Label
    )

    $fullRoot = [System.IO.Path]::GetFullPath($Root).TrimEnd(
        [System.IO.Path]::DirectorySeparatorChar,
        [System.IO.Path]::AltDirectorySeparatorChar
    )
    $fullPath = [System.IO.Path]::GetFullPath($Path)
    $prefix = $fullRoot + [System.IO.Path]::DirectorySeparatorChar
    if ((-not $fullPath.Equals($fullRoot, $PathComparison)) -and
        (-not $fullPath.StartsWith($prefix, $PathComparison))) {
        throw "$Label escapes its root."
    }
    Assert-NoLinkOrReparsePointInExistingAncestors $fullRoot "$Label root"
    Assert-NoLinkOrReparsePointInExistingAncestors $fullPath $Label
}

function Get-PhysicalPathIdentity {
    param([string]$Path, [string]$Label)

    $full = Get-ProviderCompatibleFullPath $Path $Label
    if (-not (Test-Path -LiteralPath $full)) {
        throw "$Label does not exist for physical identity: $full"
    }
    Assert-NoLinkOrReparsePointInExistingAncestors $full $Label
    if (-not $IsWindowsPlatform) {
        return Get-AliasNormalizedFullPath $full $Label
    }
    try {
        return [Sdp.Install.NativePath]::GetFileSystemIdentity($full)
    } catch {
        throw "$Label physical/device identity is unavailable: $($_.Exception.Message)"
    }
}

function Get-PhysicalPathState {
    param(
        [string]$Path,
        [string]$Label,
        [switch]$RequireExistingRoot
    )

    $full = Get-ProviderCompatibleFullPath $Path $Label
    $rootExists = Test-Path -LiteralPath $full
    if ($RequireExistingRoot -and (-not $rootExists)) {
        throw "$Label must be an existing directory: $full"
    }
    if ($rootExists -and (-not (Test-Path -LiteralPath $full -PathType Container))) {
        throw "$Label must be a directory: $full"
    }

    $probe = $full
    while (-not (Test-Path -LiteralPath $probe)) {
        $parent = Get-FileSystemParentPath $probe
        if ([string]::IsNullOrWhiteSpace($parent) -or
            $parent.Equals($probe, $PathComparison)) {
            throw "$Label has no existing ancestor whose physical identity can be determined."
        }
        $probe = $parent
    }

    $ancestorIdentities = @{}
    $nearestExistingIdentity = $null
    $namespaceRootIdentity = $null
    $current = $probe
    while (-not [string]::IsNullOrWhiteSpace($current)) {
        if (-not (Test-Path -LiteralPath $current -PathType Container)) {
            throw "$Label has an existing non-directory ancestor: $current"
        }
        $identity = Get-PhysicalPathIdentity $current "$Label ancestor '$current'"
        if ($null -eq $nearestExistingIdentity) {
            $nearestExistingIdentity = $identity
        }
        $namespaceRootIdentity = $identity
        $ancestorIdentities[$identity] = $current
        $parent = Get-FileSystemParentPath $current
        if ([string]::IsNullOrWhiteSpace($parent) -or
            $parent.Equals($current, $PathComparison)) {
            break
        }
        $current = $parent
    }

    $rootIdentity = if ($rootExists) {
        Get-PhysicalPathIdentity $full $Label
    } else {
        $null
    }
    return [pscustomobject]@{
        FullPath = $full
        RootExists = [bool]$rootExists
        RootIdentity = $rootIdentity
        VolumeIdentity = ([string]$nearestExistingIdentity).Split(':')[0]
        NamespaceRootIdentity = $namespaceRootIdentity
        AncestorIdentities = $ancestorIdentities
    }
}

function Assert-NoPhysicalTreeOverlap {
    param(
        $SourceState,
        $OtherState,
        [string]$Message
    )

    if ([string]::IsNullOrWhiteSpace([string]$SourceState.RootIdentity)) {
        throw 'The SDP source root has no physical/device identity.'
    }
    $otherIsInsideSource = $OtherState.AncestorIdentities.ContainsKey(
        [string]$SourceState.RootIdentity
    )
    $sourceIsInsideOther = (-not [string]::IsNullOrWhiteSpace(
        [string]$OtherState.RootIdentity
    )) -and $SourceState.AncestorIdentities.ContainsKey(
        [string]$OtherState.RootIdentity
    )
    if ($otherIsInsideSource -or $sourceIsInsideOther) {
        throw $Message
    }
    if ($IsWindowsPlatform -and
        ([string]$SourceState.VolumeIdentity -ceq [string]$OtherState.VolumeIdentity) -and
        ([string]$SourceState.NamespaceRootIdentity -cne
        [string]$OtherState.NamespaceRootIdentity)) {
        throw "$Message The same device is exposed through incomparable filesystem roots."
    }
}

function Assert-DestinationTopology {
    param(
        [string]$Root,
        [string]$Destination,
        [string]$Label
    )

    $fullRoot = Get-ProviderCompatibleFullPath $Root "$Label root"
    $fullDestination = Get-ProviderCompatibleFullPath $Destination $Label
    Assert-ContainedPhysicalPath $fullRoot $fullDestination $Label

    $trimCharacters = [char[]]@(
        [System.IO.Path]::DirectorySeparatorChar,
        [System.IO.Path]::AltDirectorySeparatorChar
    )
    $fullRootKey = $fullRoot.TrimEnd($trimCharacters)
    if ($fullRootKey.Length -eq 0) { $fullRootKey = $fullRoot }

    if ([System.IO.Directory]::Exists($fullDestination) -or
        ((Test-Path -LiteralPath $fullDestination) -and
        (-not [System.IO.File]::Exists($fullDestination)))) {
        throw "$Label exists but is not a file: $fullDestination"
    }

    $current = Get-FileSystemParentPath $fullDestination
    while (-not [string]::IsNullOrWhiteSpace($current)) {
        if ([System.IO.File]::Exists($current) -or
            ((Test-Path -LiteralPath $current) -and
            (-not [System.IO.Directory]::Exists($current)))) {
            throw "$Label has an existing non-directory ancestor: $current"
        }
        $currentKey = $current.TrimEnd($trimCharacters)
        if ($currentKey.Length -eq 0) { $currentKey = $current }
        if ($currentKey.Equals($fullRootKey, $PathComparison)) {
            break
        }
        $parent = Get-FileSystemParentPath $current
        if ([string]::IsNullOrWhiteSpace($parent) -or
            $parent.Equals($current, $PathComparison)) {
            throw "$Label escapes its root."
        }
        $current = $parent
    }
}

function Get-AliasNormalizedFullPath {
    param([string]$Path, [string]$Label)

    $full = [System.IO.Path]::GetFullPath($Path)
    if (-not $IsWindowsPlatform) { return $full }

    $probe = $full
    $missingSegments = New-Object System.Collections.ArrayList
    while ($true) {
        try {
            [void][System.IO.File]::GetAttributes($probe)
            break
        } catch [System.IO.FileNotFoundException] {
        } catch [System.IO.DirectoryNotFoundException] {
        }
        $leaf = Split-Path -Leaf $probe
        $parent = Get-FileSystemParentPath $probe
        if ([string]::IsNullOrWhiteSpace($leaf) -or
            [string]::IsNullOrWhiteSpace($parent) -or
            $parent.Equals($probe, $PathComparison)) {
            throw "$Label cannot be normalized to an existing physical ancestor."
        }
        $missingSegments.Insert(0, $leaf)
        $probe = $parent
    }

    $buffer = [System.Text.StringBuilder]::new(32768)
    $length = [Sdp.Install.NativePath]::GetLongPathName(
        $probe,
        $buffer,
        [uint32]$buffer.Capacity
    )
    if (($length -eq 0) -or ($length -ge $buffer.Capacity)) {
        throw "$Label could not be expanded to its physical long-path identity."
    }
    $expanded = $buffer.ToString()
    foreach ($segment in $missingSegments) {
        $expanded = Join-Path $expanded ([string]$segment)
    }
    return [System.IO.Path]::GetFullPath($expanded)
}

function Test-PhysicalPathWithin {
    param([string]$Candidate, [string]$Parent, [string]$Kind)

    $candidateFull = [System.IO.Path]::GetFullPath($Candidate).TrimEnd(
        [System.IO.Path]::DirectorySeparatorChar,
        [System.IO.Path]::AltDirectorySeparatorChar
    )
    $parentFull = [System.IO.Path]::GetFullPath($Parent).TrimEnd(
        [System.IO.Path]::DirectorySeparatorChar,
        [System.IO.Path]::AltDirectorySeparatorChar
    )
    if ($candidateFull.Equals($parentFull, $PathComparison)) { return $true }
    return ($Kind -eq 'tree') -and $candidateFull.StartsWith(
        $parentFull + [System.IO.Path]::DirectorySeparatorChar,
        $PathComparison
    )
}

function Assert-NoPhysicalGitSegment {
    param([string]$Root, [string]$Path, [string]$Label)

    $physicalRoot = Get-AliasNormalizedFullPath $Root "$Label root"
    $physicalPath = Get-AliasNormalizedFullPath $Path $Label
    if (-not (Test-PhysicalPathWithin $physicalPath $physicalRoot 'tree')) {
        throw "$Label escapes its physical root."
    }
    $relative = $physicalPath.Substring($physicalRoot.TrimEnd(
        [System.IO.Path]::DirectorySeparatorChar,
        [System.IO.Path]::AltDirectorySeparatorChar
    ).Length).TrimStart(
        [System.IO.Path]::DirectorySeparatorChar,
        [System.IO.Path]::AltDirectorySeparatorChar
    )
    foreach ($segment in $relative.Split([char[]]@(
        [System.IO.Path]::DirectorySeparatorChar,
        [System.IO.Path]::AltDirectorySeparatorChar
    ), [System.StringSplitOptions]::RemoveEmptyEntries)) {
        if ($segment.Equals('.git', [System.StringComparison]::OrdinalIgnoreCase)) {
            throw "$Label resolves through a Git administrative path."
        }
    }
    return $physicalPath
}

function Get-PortablePathKey {
    param([string]$Path)
    return $Path.ToLowerInvariant()
}

function Test-PortablePathWithin {
    param([string]$Candidate, [string]$Parent, [string]$Kind)

    $candidateKey = Get-PortablePathKey $Candidate
    $parentKey = Get-PortablePathKey $Parent
    if ($candidateKey -ceq $parentKey) { return $true }
    return ($Kind -eq 'tree') -and $candidateKey.StartsWith($parentKey + '/')
}

function ConvertTo-YamlQuotedScalar {
    param([string]$Value)

    return '"' + $Value.Replace('\', '\\').Replace('"', '\"') + '"'
}

$InstallerScriptRoot = Get-ProviderCompatibleFullPath $PSScriptRoot 'installer script root'
$ToolkitRoot = Get-ProviderCompatibleFullPath (Split-Path -Parent $InstallerScriptRoot) 'Toolkit root'
$RepositoryRoot = Get-ProviderCompatibleFullPath (Join-Path $ToolkitRoot '..') 'SDP repository root'
Assert-NoLinkOrReparsePointInExistingAncestors $RepositoryRoot 'SDP repository root'
$InstallManifestPath = Join-Path $ToolkitRoot 'SDP-install.manifest.json'
Assert-ContainedPhysicalPath $RepositoryRoot $InstallManifestPath 'installation manifest source'
if (-not (Test-Path -LiteralPath $InstallManifestPath -PathType Leaf)) {
    throw "Canonical installation manifest is missing: $InstallManifestPath"
}
try {
    $InstallManifest = [System.IO.File]::ReadAllText(
        $InstallManifestPath,
        [System.Text.Encoding]::UTF8
    ) | ConvertFrom-Json
} catch {
    throw "Cannot parse canonical installation manifest: $($_.Exception.Message)"
}
Assert-JsonObjectShape $InstallManifest @(
    'schemaVersion', 'contractId', 'toolkitVersion', 'sources',
    'capabilities', 'generators', 'entries', 'exclusions'
) @(
    'schemaVersion', 'contractId', 'toolkitVersion', 'sources',
    'capabilities', 'generators', 'entries', 'exclusions'
) 'installation manifest'
Assert-JsonString $InstallManifest.schemaVersion 'installation manifest schemaVersion'
Assert-JsonString $InstallManifest.contractId 'installation manifest contractId'
Assert-JsonString $InstallManifest.toolkitVersion 'installation manifest toolkitVersion'
Assert-JsonObjectShape $InstallManifest.sources @(
    'repositoryRoot', 'pathStyle'
) @(
    'repositoryRoot', 'pathStyle'
) 'installation manifest sources'
Assert-JsonString $InstallManifest.sources.repositoryRoot 'installation manifest sources.repositoryRoot'
Assert-JsonString $InstallManifest.sources.pathStyle 'installation manifest sources.pathStyle'
Assert-JsonArray $InstallManifest.capabilities 1 'installation manifest capabilities'
Assert-JsonArray $InstallManifest.generators 2 'installation manifest generators'
Assert-JsonArray $InstallManifest.entries 1 'installation manifest entries'
Assert-JsonArray $InstallManifest.exclusions 1 'installation manifest exclusions'

if ($SupportedInstallManifestSchemas -cnotcontains [string]$InstallManifest.schemaVersion) {
    throw "Unsupported installation manifest schema '$($InstallManifest.schemaVersion)'."
}
if ([string]$InstallManifest.contractId -cne 'sdp-install') {
    throw "Unsupported installation contract '$($InstallManifest.contractId)'."
}
if ([string]$InstallManifest.sources.repositoryRoot -cne '..' -or
    [string]$InstallManifest.sources.pathStyle -cne 'forward-slash-relative') {
    throw 'Unsupported installation source-root or path-style contract.'
}

$ProjectRoot = Get-ProviderCompatibleFullPath $ProjectRoot 'consuming project root'
if (-not (Test-Path -LiteralPath $ProjectRoot -PathType Container)) {
    throw "The consuming project root must be an existing directory: $ProjectRoot"
}
Assert-NoLinkOrReparsePointInExistingAncestors $ProjectRoot 'consuming project root'
$RepositoryPhysicalState = Get-PhysicalPathState `
    $RepositoryRoot `
    'SDP repository root' `
    -RequireExistingRoot
$ProjectPhysicalState = Get-PhysicalPathState `
    $ProjectRoot `
    'consuming project root' `
    -RequireExistingRoot
Assert-NoPhysicalTreeOverlap `
    $RepositoryPhysicalState `
    $ProjectPhysicalState `
    'The consuming project and SDP source repository must be physically separate trees.'
$SdpRoot = Join-Path $ProjectRoot 'SDP'
$ExpectedInstalledManifestRelative = 'SDP/Framework/installed-toolkit.manifest.yaml'
$ExpectedProjectManifestRelative = 'SDP/SDP-project.manifest.yaml'

if (-not [string]::IsNullOrWhiteSpace($BackupRoot)) {
    Assert-CanonicalWindowsExtendedPath $BackupRoot 'backup root'
}
if ([string]::IsNullOrWhiteSpace($BackupRoot)) {
    $BackupRoot = Join-Path $SdpRoot ".sdp-backups\$RunStamp"
} elseif (-not [System.IO.Path]::IsPathRooted($BackupRoot)) {
    $BackupRoot = Join-Path $ProjectRoot $BackupRoot
}
$BackupRoot = Get-ProviderCompatibleFullPath $BackupRoot 'backup root'
Assert-NoLinkOrReparsePointInExistingAncestors $BackupRoot 'backup root'
$BackupPhysicalState = Get-PhysicalPathState $BackupRoot 'backup root'
Assert-NoPhysicalTreeOverlap `
    $RepositoryPhysicalState `
    $BackupPhysicalState `
    'The backup root and SDP source repository must be physically separate trees.'

$ToolkitVersion = [string]$InstallManifest.toolkitVersion
[void](ConvertTo-SemVer $ToolkitVersion)
$Capabilities = $InstallManifest.capabilities
$CapabilityById = @{}
foreach ($capabilityValue in $Capabilities) {
    Assert-JsonString $capabilityValue 'installation manifest capability'
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
$Generators = $InstallManifest.generators
$Entries = $InstallManifest.entries
$Exclusions = $InstallManifest.exclusions
$GeneratorById = @{}
$EntryById = @{}
$DestinationByKey = @{}
$ExclusionByKey = @{}
$PhysicalExclusionByKey = @{}
$ExclusionRows = New-Object System.Collections.ArrayList

foreach ($generator in $Generators) {
    Assert-JsonObjectShape $generator @('id', 'type', 'format') @(
        'id', 'type', 'format', 'facts', 'dynamicFacts', 'content'
    ) 'installation manifest generator'
    Assert-JsonString $generator.id 'installation manifest generator id'
    Assert-JsonString $generator.type "generator '$($generator.id)' type"
    Assert-JsonString $generator.format "generator '$($generator.id)' format"
    $generatorId = [string]$generator.id
    if (($generatorId -cnotmatch '^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$') -or
        $GeneratorById.ContainsKey($generatorId)) {
        throw "Installation contract contains a missing or duplicate generator ID '$generatorId'."
    }
    $generatorType = [string]$generator.type
    Assert-AllowedValue $generatorType @('installed-toolkit-manifest', 'empty-ledger') "generator '$generatorId' type"
    if ($generatorType -eq 'installed-toolkit-manifest') {
        Assert-JsonObjectShape $generator @(
            'id', 'type', 'format', 'facts', 'dynamicFacts'
        ) @(
            'id', 'type', 'format', 'facts', 'dynamicFacts'
        ) "generator '$generatorId'"
        if (($generatorId -cne 'installed-toolkit-manifest') -or
            ([string]$generator.format -cne 'yaml')) {
            throw "Generator '$generatorId' must be installed-toolkit-manifest in YAML format."
        }
    } else {
        Assert-JsonObjectShape $generator @(
            'id', 'type', 'format', 'content'
        ) @(
            'id', 'type', 'format', 'content'
        ) "generator '$generatorId'"
        if (($generatorId -cne 'empty-ledger') -or
            ([string]$generator.format -cne 'ndjson') -or
            ($generator.content -isnot [string]) -or
            ([string]$generator.content -cne '')) {
            throw "Generator '$generatorId' must be empty-ledger with empty NDJSON content."
        }
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
Assert-JsonObjectShape $InstalledGenerator.facts @(
    'schemaVersion', 'toolkitVersion', 'frameworkVersion',
    'agentsContractVersion', 'installerVersion', 'skills', 'capabilities'
) @(
    'schemaVersion', 'toolkitVersion', 'frameworkVersion',
    'agentsContractVersion', 'installerVersion', 'skills', 'capabilities'
) "generator 'installed-toolkit-manifest' facts"
foreach ($factName in @(
    'schemaVersion', 'toolkitVersion', 'frameworkVersion',
    'agentsContractVersion', 'installerVersion'
)) {
    Assert-JsonString $InstalledGenerator.facts.$factName "installed-manifest fact '$factName'"
}
$InstallerVersion = [string]$InstalledGenerator.facts.installerVersion
[void](ConvertTo-SemVer $InstallerVersion)
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
    [void](ConvertTo-SemVer ([string]$InstalledGenerator.facts.$factName))
}
$skillFacts = $InstalledGenerator.facts.skills
if (($null -eq $skillFacts) -or ($skillFacts -isnot [pscustomobject]) -or
    (@($skillFacts.psobject.Properties).Count -eq 0)) {
    throw "Installed-manifest generator must declare installed skill versions."
}
foreach ($skillProperty in $skillFacts.psobject.Properties) {
    if ($skillProperty.Name -cnotmatch '^sdp-[a-z0-9]+(?:-[a-z0-9]+)*$') {
        throw "Installed-manifest generator has invalid skill ID '$($skillProperty.Name)'."
    }
    Assert-JsonString $skillProperty.Value "installed skill '$($skillProperty.Name)' version"
    [void](ConvertTo-SemVer ([string]$skillProperty.Value))
}
Assert-JsonArray $InstalledGenerator.facts.capabilities 1 'installed-manifest capabilities'
foreach ($capabilityValue in $InstalledGenerator.facts.capabilities) {
    Assert-JsonString $capabilityValue 'installed-manifest capability'
}
if (($InstalledGenerator.facts.capabilities -join "`n") -cne
    ($Capabilities -join "`n")) {
    throw "Installed-manifest generator capabilities do not match the contract capabilities."
}
$dynamicFacts = $InstalledGenerator.dynamicFacts
Assert-JsonObjectShape $dynamicFacts @(
    'toolkitInstalledAt', 'sourceCommit'
) @(
    'toolkitInstalledAt', 'sourceCommit'
) "generator 'installed-toolkit-manifest' dynamicFacts"
Assert-JsonObjectShape $dynamicFacts.toolkitInstalledAt @(
    'source', 'sameToolkitVersionPolicy'
) @(
    'source', 'sameToolkitVersionPolicy'
) 'toolkitInstalledAt dynamic fact'
Assert-JsonObjectShape $dynamicFacts.sourceCommit @(
    'source', 'unavailableValue'
) @(
    'source', 'unavailableValue'
) 'sourceCommit dynamic fact'
Assert-JsonString $dynamicFacts.toolkitInstalledAt.source 'toolkitInstalledAt source'
Assert-JsonString $dynamicFacts.toolkitInstalledAt.sameToolkitVersionPolicy 'toolkitInstalledAt same-version policy'
Assert-JsonString $dynamicFacts.sourceCommit.source 'sourceCommit source'
if (([string]$dynamicFacts.toolkitInstalledAt.source -cne 'utc-now') -or
    ([string]$dynamicFacts.toolkitInstalledAt.sameToolkitVersionPolicy -cne
        'preserve-existing') -or
    ([string]$dynamicFacts.sourceCommit.source -cne 'repository-head') -or
    ($null -ne $dynamicFacts.sourceCommit.unavailableValue)) {
    throw "Installed-manifest generator declares unsupported dynamic-fact policies."
}

foreach ($exclusion in $Exclusions) {
    Assert-JsonObjectShape $exclusion @('path', 'kind', 'reason') @(
        'path', 'kind', 'reason'
    ) 'installation manifest exclusion'
    Assert-JsonString $exclusion.path 'installation manifest exclusion path'
    Assert-JsonString $exclusion.kind "exclusion '$($exclusion.path)' kind"
    Assert-JsonString $exclusion.reason "exclusion '$($exclusion.path)' reason"
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
    $excludedKey = Get-PortablePathKey $excludedPath
    if ($ExclusionByKey.ContainsKey($excludedKey)) {
        throw "Installation contract contains duplicate exclusion '$excludedPath'."
    }
    $ExclusionByKey[$excludedKey] = $excludedKind
    $excludedFull = Join-PortablePath $RepositoryRoot $excludedPath "exclusion '$excludedPath'"
    Assert-ContainedPhysicalPath $RepositoryRoot $excludedFull "exclusion '$excludedPath'"
    $excludedPhysical = Assert-NoPhysicalGitSegment `
        $RepositoryRoot `
        $excludedFull `
        "exclusion '$excludedPath'"
    $excludedPhysicalKey = if ($IsWindowsPlatform) {
        $excludedPhysical.ToLowerInvariant()
    } else {
        $excludedPhysical
    }
    if ($PhysicalExclusionByKey.ContainsKey($excludedPhysicalKey)) {
        throw "Installation contract contains a physical-alias exclusion collision '$excludedPath'."
    }
    $PhysicalExclusionByKey[$excludedPhysicalKey] = $excludedPath
    if (($excludedKind -eq 'file') -and
        (-not (Test-Path -LiteralPath $excludedFull -PathType Leaf))) {
        throw "Excluded file does not exist: $excludedPath"
    }
    if (($excludedKind -eq 'tree') -and
        (-not (Test-Path -LiteralPath $excludedFull -PathType Container))) {
        throw "Excluded tree does not exist: $excludedPath"
    }
    [void]$ExclusionRows.Add([pscustomobject]@{
        Path = $excludedPath
        Kind = $excludedKind
        FullPath = $excludedFull
        PhysicalPath = $excludedPhysical
    })
}

foreach ($requiredLiveExclusion in @(
    '01--Mandate', '02--Study', '03--Requirements', '04--Architecture',
    '05--DesignAnalysis', '06--Design', '07--Implementation', 'CodeReview',
    'Fixes', 'Instructions', 'Refactors', 'Releases', 'Sprints', 'Traceability',
    'Verification', 'RELEASE-NOTES.md', 'SDP.manifest.yaml',
    'SDP-DOCUMENT-GUIDE.md', 'payload', 'skills',
    'Toolkit/payload/project-root/AGENTS-project.md.template',
    'Toolkit/payload/sdp-root/AGENT-REMINDERS.md.template'
)) {
    if (-not $ExclusionByKey.ContainsKey((Get-PortablePathKey $requiredLiveExclusion))) {
        throw "Installation contract is missing required live-state exclusion '$requiredLiveExclusion'."
    }
}

$CanonicalCapabilityBySchema = @{
    'toolkit/schemas/fix-record.schema.json' = 'sdp.release.v1'
    'toolkit/schemas/release-record.schema.json' = 'sdp.release.v1'
    'toolkit/schemas/sdp-project-manifest.schema.json' = 'sdp.manifest.v1'
    'toolkit/schemas/installed-toolkit-manifest.schema.json' = 'sdp.manifest.v1'
    'toolkit/schemas/current-index.schema.json' = 'sdp.traceability.current-index.v1'
    'toolkit/schemas/relations.schema.json' = 'sdp.traceability.relations.v1'
    'toolkit/schemas/ledger-event.schema.json' = 'sdp.traceability.ledger-events.v1'
}

foreach ($entry in $Entries) {
    Assert-JsonObjectShape $entry @(
        'id', 'kind', 'destination', 'ownership', 'selectionPolicy',
        'installPolicy', 'refreshPolicy', 'backupPolicy', 'forcePolicy',
        'migrationPolicy', 'governing'
    ) @(
        'id', 'kind', 'source', 'generator', 'destination', 'ownership',
        'selectionPolicy', 'installPolicy', 'refreshPolicy', 'backupPolicy',
        'forcePolicy', 'migrationPolicy', 'governing'
    ) 'installation manifest entry'
    foreach ($propertyName in @(
        'id', 'kind', 'destination', 'ownership', 'selectionPolicy',
        'installPolicy', 'refreshPolicy', 'backupPolicy', 'forcePolicy',
        'migrationPolicy'
    )) {
        Assert-JsonString $entry.$propertyName "entry property '$propertyName'"
    }
    Assert-JsonObjectShape $entry.governing @(
        'capability', 'schema'
    ) @(
        'capability', 'schema'
    ) "entry '$($entry.id)' governing"
    Assert-JsonString $entry.governing.capability "entry '$($entry.id)' governing capability"
    if (($null -ne $entry.governing.schema) -and
        ($entry.governing.schema -isnot [string])) {
        throw "Entry '$($entry.id)' governing schema must be a string or null."
    }
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
    Assert-PortableRelativePath $destination "entry '$entryId' destination" -Destination
    $destinationKey = Get-PortablePathKey $destination
    if ($DestinationByKey.ContainsKey($destinationKey)) {
        throw "Installation contract contains duplicate destination '$destination'."
    }
    foreach ($existingDestinationKey in @($DestinationByKey.Keys)) {
        if ($destinationKey.StartsWith(
            [string]$existingDestinationKey + '/',
            [System.StringComparison]::Ordinal
        ) -or ([string]$existingDestinationKey).StartsWith(
            $destinationKey + '/',
            [System.StringComparison]::Ordinal
        )) {
            $existingDestination = $DestinationByKey[$existingDestinationKey]
            throw "Installation contract contains ancestor/descendant destination conflict " +
                "'$destination' and '$($existingDestination.Path)'."
        }
    }
    $DestinationByKey[$destinationKey] = [pscustomobject]@{
        Id = $entryId
        Path = $destination
    }
    [void](Join-PortablePath $ProjectRoot $destination "entry '$entryId' destination")

    if ($entryKind -eq 'copied') {
        if (($entry.psobject.Properties.Name -cnotcontains 'source') -or
            ($entry.psobject.Properties.Name -contains 'generator')) {
            throw "Copied entry '$entryId' must not declare a generator."
        }
        Assert-JsonString $entry.source "entry '$entryId' source"
        $source = [string]$entry.source
        Assert-PortableRelativePath $source "entry '$entryId' source"
        if ($ownership -eq 'toolkit-managed') {
            $inAllowedSourceClass = (Test-PortablePathWithin $source 'Toolkit/payload' 'tree') -or
                (Test-PortablePathWithin $source 'Toolkit/skills' 'tree')
        } else {
            $inAllowedSourceClass = Test-PortablePathWithin $source 'Toolkit/project-templates' 'tree'
        }
        if (-not $inAllowedSourceClass) {
            throw "Entry '$entryId' source is outside its declared ownership source class: $source"
        }
        foreach ($excluded in $ExclusionRows) {
            if (Test-PortablePathWithin $source $excluded.Path $excluded.Kind) {
                throw "Entry '$entryId' uses explicitly excluded source '$source'."
            }
        }
        $sourcePath = Join-PortablePath $RepositoryRoot $source "entry '$entryId' source"
        Assert-ContainedPhysicalPath $RepositoryRoot $sourcePath "entry '$entryId' source"
        $sourcePhysical = Assert-NoPhysicalGitSegment `
            $RepositoryRoot `
            $sourcePath `
            "entry '$entryId' source"
        foreach ($excluded in $ExclusionRows) {
            if (Test-PhysicalPathWithin $sourcePhysical $excluded.PhysicalPath $excluded.Kind) {
                throw "Entry '$entryId' uses a physical alias of excluded source '$($excluded.Path)'."
            }
        }
        if (-not (Test-Path -LiteralPath $sourcePath -PathType Leaf)) {
            throw "Entry '$entryId' source does not exist: $source"
        }
    } elseif ($entryKind -eq 'generated') {
        if (($entry.psobject.Properties.Name -cnotcontains 'generator') -or
            ($entry.psobject.Properties.Name -contains 'source')) {
            throw "Generated entry '$entryId' must not declare a source."
        }
        Assert-JsonString $entry.generator "entry '$entryId' generator"
        $generatorId = [string]$entry.generator
        if (-not $GeneratorById.ContainsKey($generatorId)) {
            throw "Entry '$entryId' references unknown generator '$generatorId'."
        }
        if ((($ownership -eq 'toolkit-managed') -and
            ($generatorId -cne 'installed-toolkit-manifest')) -or
            (($ownership -eq 'project-owned') -and ($generatorId -cne 'empty-ledger'))) {
            throw "Entry '$entryId' generator is outside its declared ownership source class."
        }
    }
    $schemaPath = $entry.governing.schema
    if ($null -ne $schemaPath) {
        $schemaRelative = [string]$schemaPath
        Assert-JsonString $schemaRelative "entry '$entryId' governing schema"
        Assert-PortableRelativePath $schemaRelative "entry '$entryId' governing schema"
        if (-not (Test-PortablePathWithin $schemaRelative 'Toolkit/schemas' 'tree')) {
            throw "Entry '$entryId' governing schema must be under Toolkit/schemas."
        }
        $schemaFull = Join-PortablePath $RepositoryRoot $schemaRelative "entry '$entryId' schema"
        Assert-ContainedPhysicalPath $RepositoryRoot $schemaFull "entry '$entryId' schema"
        [void](Assert-NoPhysicalGitSegment `
            $RepositoryRoot `
            $schemaFull `
            "entry '$entryId' schema")
        if (-not (Test-Path -LiteralPath $schemaFull -PathType Leaf)) {
            throw "Entry '$entryId' governing schema does not exist: $schemaRelative"
        }
        $schemaKey = Get-PortablePathKey $schemaRelative
        if (-not $CanonicalCapabilityBySchema.ContainsKey($schemaKey)) {
            throw "Entry '$entryId' uses unsupported governing schema '$schemaRelative'."
        }
        if ([string]$CanonicalCapabilityBySchema[$schemaKey] -cne $governingCapability) {
            throw "Entry '$entryId' governing schema/capability pairing is invalid."
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

# The complete installation contract is closed and validated before any target
# file is inspected. Physical containment checks then cover every selected or
# unselected destination so a later option cannot expose an unchecked path.
$PhysicalDestinationByKey = @{}
foreach ($entry in $Entries) {
    $destinationPath = Join-PortablePath $ProjectRoot ([string]$entry.destination) "entry '$($entry.id)' destination"
    Assert-DestinationTopology `
        $ProjectRoot `
        $destinationPath `
        "entry '$($entry.id)' destination"
    $physicalDestination = Assert-NoPhysicalGitSegment `
        $ProjectRoot `
        $destinationPath `
        "entry '$($entry.id)' destination"
    $physicalDestinationKey = if ($IsWindowsPlatform) {
        $physicalDestination.ToLowerInvariant()
    } else {
        $physicalDestination
    }
    if ($PhysicalDestinationByKey.ContainsKey($physicalDestinationKey)) {
        throw "Installation contract contains physical-alias destination collision '$($entry.destination)'."
    }
    $PhysicalDestinationByKey[$physicalDestinationKey] = [string]$entry.id
}
Assert-ContainedPhysicalPath $ProjectRoot $InstalledManifestPath 'installed manifest destination'
Assert-ContainedPhysicalPath $ProjectRoot $ProjectManifestPath 'project manifest destination'

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
        $resolvedTopIdentity = Get-PhysicalPathIdentity $resolvedTop 'Git repository root'
    } catch {
        return $null
    }
    if ($resolvedTopIdentity -cne [string]$RepositoryPhysicalState.RootIdentity) { return $null }
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

function ConvertFrom-InstalledManifestDocument {
    param([string]$Content)

    $label = 'installed Toolkit manifest'
    $document = ConvertFrom-StrictYamlDocument $Content $label @(
        'schemaVersion', 'toolkitVersion', 'toolkitInstalledAt'
    )
    $allowedRoot = @(
        'schemaVersion', 'toolkitVersion', 'frameworkVersion',
        'agentsContractVersion', 'installerVersion', 'toolkitInstalledAt',
        'sourceCommit', 'skills', 'capabilities'
    )
    $rootKeys = @($document.ActualPaths | Where-Object { $_ -cnotmatch '/' })
    foreach ($required in $allowedRoot) {
        if ($rootKeys -cnotcontains $required) {
            throw "$label is missing required root field '$required'."
        }
    }
    foreach ($rootKey in $rootKeys) {
        if ($allowedRoot -cnotcontains $rootKey) {
            throw "$label contains unknown root field '$rootKey'."
        }
    }
    if ([string]$document.Containers['skills'] -cne 'mapping') {
        throw "$label field 'skills' must be a mapping."
    }
    if ([string]$document.Containers['capabilities'] -cne 'sequence') {
        throw "$label field 'capabilities' must be a scalar sequence."
    }

    foreach ($path in $document.ActualPaths) {
        if ($path -cnotmatch '/') { continue }
        if ($path -cnotmatch '^skills/[a-z][a-z0-9]*(?:-[a-z0-9]+)*$') {
            throw "$label contains unsupported nested field '$path'."
        }
    }
    $skillPaths = @($document.ActualPaths | Where-Object { $_ -cmatch '^skills/' })
    if ($skillPaths.Count -eq 0) {
        throw "$label must declare at least one installed skill."
    }
    foreach ($path in $skillPaths) {
        $version = Get-StrictYamlString $document $path $label
        [void](ConvertTo-SemVer $version)
    }
    $capabilityRows = @($document.Sequences['capabilities'])
    if ($capabilityRows.Count -eq 0) {
        throw "$label must declare at least one capability."
    }
    $capabilityKeys = @{}
    foreach ($row in $capabilityRows) {
        if ([string]$row.Type -cne 'string') {
            throw "$label capabilities must be strings."
        }
        $capability = [string]$row.Value
        if (($capability -cnotmatch '^sdp\.[a-z0-9.-]+\.v[0-9]+$') -or
            $capabilityKeys.ContainsKey($capability)) {
            throw "$label contains an invalid or duplicate capability '$capability'."
        }
        $capabilityKeys[$capability] = $true
    }

    foreach ($versionField in @(
        'toolkitVersion', 'frameworkVersion', 'agentsContractVersion', 'installerVersion'
    )) {
        [void](ConvertTo-SemVer (Get-StrictYamlString $document $versionField $label))
    }
    [void](Get-StrictYamlString $document 'schemaVersion' $label)
    $installedAt = Get-StrictYamlString $document 'toolkitInstalledAt' $label
    $parsedInstalledAt = [DateTimeOffset]::MinValue
    $rfc3339DateTimePattern = '^\d{4}-\d{2}-\d{2}[Tt]\d{2}:\d{2}:\d{2}' +
        '(?:\.\d+)?(?:[Zz]|[+-]\d{2}:\d{2})$'
    if (($installedAt -cnotmatch $rfc3339DateTimePattern) -or
        (-not [DateTimeOffset]::TryParse(
        $installedAt,
        [System.Globalization.CultureInfo]::InvariantCulture,
        [System.Globalization.DateTimeStyles]::RoundtripKind,
        [ref]$parsedInstalledAt
    ))) {
        throw "$label toolkitInstalledAt must be an RFC 3339 timestamp."
    }
    $sourceCommit = $document.Values['sourceCommit']
    if (($null -eq $sourceCommit) -or
        ([string]$sourceCommit.Type -notin @('string', 'null')) -or
        (([string]$sourceCommit.Type -eq 'string') -and
        [string]::IsNullOrWhiteSpace([string]$sourceCommit.Value))) {
        throw "$label sourceCommit must be a non-empty string or null."
    }
    return $document
}

function Test-InstalledManifestEquivalent {
    param($Document)

    if ($null -eq $Document) { return $false }
    $facts = $InstalledGenerator.facts
    foreach ($field in @(
        'schemaVersion', 'toolkitVersion', 'frameworkVersion',
        'agentsContractVersion', 'installerVersion'
    )) {
        if ((Get-StrictYamlString $Document $field 'installed Toolkit manifest') -cne
            [string]$facts.$field) {
            return $false
        }
    }
    if ((Get-StrictYamlString $Document 'toolkitInstalledAt' 'installed Toolkit manifest') -cne
        $InstalledAt) {
        return $false
    }
    $actualSource = $Document.Values['sourceCommit']
    if ($null -eq $SourceCommit) {
        if ([string]$actualSource.Type -cne 'null') { return $false }
    } elseif (([string]$actualSource.Type -cne 'string') -or
        ([string]$actualSource.Value -cne $SourceCommit)) {
        return $false
    }

    $actualSkillPaths = @($Document.ActualPaths | Where-Object { $_ -cmatch '^skills/' })
    $expectedSkills = @($facts.skills.psobject.Properties)
    if ($actualSkillPaths.Count -ne $expectedSkills.Count) { return $false }
    foreach ($skill in $expectedSkills) {
        $path = 'skills/' + $skill.Name
        if ($actualSkillPaths -cnotcontains $path) { return $false }
        if ((Get-StrictYamlString $Document $path 'installed Toolkit manifest') -cne
            [string]$skill.Value) {
            return $false
        }
    }
    $actualCapabilities = @($Document.Sequences['capabilities'] | ForEach-Object {
        [string]$_.Value
    })
    return (($actualCapabilities -join "`n") -ceq (@($facts.capabilities) -join "`n"))
}

$BlockReason = $null
$BlockEntryId = $null
$BlockDestination = $null
$BlockOwnership = $null
$InstalledSchemaVersion = $null
$InstalledToolkitVersion = $null
$PreviousInstalledAt = $null
$InstalledManifestDocument = $null

if (Test-Path -LiteralPath $ProjectManifestPath -PathType Leaf) {
    Assert-ContainedPhysicalPath $ProjectRoot $ProjectManifestPath 'project manifest destination'
    try {
        $projectManifestContent = [System.IO.File]::ReadAllText($ProjectManifestPath)
        $projectDocument = ConvertFrom-StrictYamlDocument `
            $projectManifestContent `
            'project manifest' `
            @('schemaVersion')
        $projectRootKeys = @($projectDocument.ActualPaths | Where-Object { $_ -cnotmatch '/' })
        if ($projectRootKeys -cnotcontains 'schemaVersion') {
            throw 'Project manifest must contain exact root field schemaVersion.'
        }
        $projectSchemaVersion = Get-StrictYamlString $projectDocument 'schemaVersion' 'project manifest'
    } catch {
        $BlockReason = 'malformed-project-manifest'
    }
    if ((-not $BlockReason) -and
        ($SupportedProjectManifestSchemas -cnotcontains $projectSchemaVersion)) {
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
    Assert-ContainedPhysicalPath $ProjectRoot $InstalledManifestPath 'installed manifest destination'
    try {
        $installedContent = [System.IO.File]::ReadAllText($InstalledManifestPath)
        $installedRootDocument = ConvertFrom-StrictYamlDocument `
            $installedContent `
            'installed Toolkit manifest' `
            @('schemaVersion', 'toolkitVersion', 'toolkitInstalledAt')
        $installedRootKeys = @($installedRootDocument.ActualPaths | Where-Object { $_ -cnotmatch '/' })
        foreach ($requiredRootKey in @('schemaVersion', 'toolkitVersion')) {
            if ($installedRootKeys -cnotcontains $requiredRootKey) {
                throw "Installed Toolkit manifest must contain exact root field $requiredRootKey."
            }
        }
        $InstalledSchemaVersion = Get-StrictYamlString `
            $installedRootDocument `
            'schemaVersion' `
            'installed Toolkit manifest'
        $InstalledToolkitVersion = Get-StrictYamlString `
            $installedRootDocument `
            'toolkitVersion' `
            'installed Toolkit manifest'
    } catch {
        $BlockReason = 'malformed-installed-manifest'
    }
    if ((-not $BlockReason) -and
        ($SupportedInstalledManifestSchemas -cnotcontains $InstalledSchemaVersion)) {
        $BlockReason = 'unsupported-installed-schema'
    } elseif (-not $BlockReason) {
        try {
            $InstalledManifestDocument = ConvertFrom-InstalledManifestDocument $installedContent
            $PreviousInstalledAt = Get-StrictYamlString `
                $InstalledManifestDocument `
                'toolkitInstalledAt' `
                'installed Toolkit manifest'
            if ((Compare-SemVer $InstalledToolkitVersion $ToolkitVersion) -gt 0) {
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
    ($InstalledToolkitVersion -cne $ToolkitVersion)
if (($InstalledToolkitVersion -ceq $ToolkitVersion) -and
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
        Assert-ContainedPhysicalPath $RepositoryRoot $sourcePath "entry '$($Entry.id)' source"
        return [System.IO.File]::ReadAllText($sourcePath)
    }
    return Get-GeneratorContent ([string]$Entry.generator)
}

$script:PlanActions = New-Object System.Collections.ArrayList
$PlanReasonRules = @{
    'missing-target' = [pscustomobject]@{ Action = 'create'; Mutates = $true }
    'missing-generated-target' = [pscustomobject]@{ Action = 'generate'; Mutates = $true }
    'content-matches' = [pscustomobject]@{ Action = 'unchanged'; Mutates = $false }
    'missing-only-content' = [pscustomobject]@{ Action = 'preserve'; Mutates = $false }
    'managed-content-differs' = [pscustomobject]@{ Action = 'preserve'; Mutates = $false }
    'backup-before-replace' = [pscustomobject]@{ Action = 'backup'; Mutates = $true }
    'refresh-managed-content' = [pscustomobject]@{ Action = 'replace'; Mutates = $true }
    'refresh-generated-content' = [pscustomobject]@{ Action = 'generate'; Mutates = $true }
    'migrate-existing-agents' = [pscustomobject]@{ Action = 'migrate'; Mutates = $true }
    'preserve-existing-agents-conflict' = [pscustomobject]@{ Action = 'migrate'; Mutates = $true }
    'malformed-project-manifest' = [pscustomobject]@{ Action = 'block'; Mutates = $false }
    'unsupported-project-schema' = [pscustomobject]@{ Action = 'block'; Mutates = $false }
    'malformed-installed-manifest' = [pscustomobject]@{ Action = 'block'; Mutates = $false }
    'unsupported-installed-schema' = [pscustomobject]@{ Action = 'block'; Mutates = $false }
    'downgrade-blocked' = [pscustomobject]@{ Action = 'block'; Mutates = $false }
}
$BlockReasonEntryIds = @{
    'malformed-project-manifest' = 'project-manifest'
    'unsupported-project-schema' = 'project-manifest'
    'malformed-installed-manifest' = 'generated-installed-toolkit-manifest'
    'unsupported-installed-schema' = 'generated-installed-toolkit-manifest'
    'downgrade-blocked' = 'generated-installed-toolkit-manifest'
}
function Add-PlanAction {
    param(
        [string]$Action,
        [string]$EntryId,
        [AllowNull()][string]$Source,
        [AllowNull()][string]$Generator,
        [string]$Destination,
        [string]$Ownership,
        [string]$Reason,
        [bool]$MutatesTarget,
        [AllowNull()][string]$TargetSource = $null
    )

    if (-not $PlanReasonRules.ContainsKey($Reason)) {
        throw "Unknown installation-plan reason '$Reason'."
    }
    $reasonRule = $PlanReasonRules[$Reason]
    if (([string]$reasonRule.Action -cne $Action) -or
        ([bool]$reasonRule.Mutates -ne $MutatesTarget)) {
        throw "Installation-plan reason '$Reason' contradicts action '$Action'."
    }
    Assert-PortableRelativePath $Destination 'plan action destination' -Destination
    if (-not [string]::IsNullOrWhiteSpace($TargetSource)) {
        Assert-PortableRelativePath $TargetSource 'plan action targetSource' -Destination
    }
    $hasSource = -not [string]::IsNullOrWhiteSpace($Source)
    $hasGenerator = -not [string]::IsNullOrWhiteSpace($Generator)
    $hasTargetSource = -not [string]::IsNullOrWhiteSpace($TargetSource)
    if ($Action -eq 'migrate') {
        if ($hasSource -or $hasGenerator -or (-not $hasTargetSource)) {
            throw 'A migrate action must identify only a targetSource.'
        }
    } elseif ($Action -eq 'block') {
        if ($hasSource -or $hasGenerator -or $hasTargetSource) {
            throw 'A block action must not identify content sources.'
        }
        if ([string]$BlockReasonEntryIds[$Reason] -cne $EntryId) {
            throw "Block reason '$Reason' contradicts entry '$EntryId'."
        }
    } elseif ((-not ($hasSource -xor $hasGenerator)) -or $hasTargetSource) {
        throw "Plan action '$Action' must identify exactly one manifest source or generator."
    }

    $row = [pscustomobject][ordered]@{
        sequence = $script:PlanActions.Count + 1
        action = $Action
        entryId = $EntryId
        source = if ([string]::IsNullOrWhiteSpace($Source)) { $null } else { $Source }
        generator = if ([string]::IsNullOrWhiteSpace($Generator)) { $null } else { $Generator }
        targetSource = if ([string]::IsNullOrWhiteSpace($TargetSource)) { $null } else { $TargetSource }
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
    Assert-ContainedPhysicalPath $ProjectRoot $managedAgentsPath 'managed AGENTS destination'
    Assert-ContainedPhysicalPath $ProjectRoot $projectAgentsPath 'project AGENTS destination'
    $managedAgentsContent = Get-EntryContent $managedAgents
    $MigrationCreatesProjectAgents = $false
    if (([string]$managedAgents.migrationPolicy -eq 'preserve-existing-agents') -and
        (Test-Path -LiteralPath $managedAgentsPath -PathType Leaf)) {
        Assert-ContainedPhysicalPath $ProjectRoot $managedAgentsPath 'managed AGENTS destination'
        $existingAgents = [System.IO.File]::ReadAllText($managedAgentsPath)
        if ($existingAgents -cne $managedAgentsContent) {
            if (-not (Test-Path -LiteralPath $projectAgentsPath -PathType Leaf)) {
                Add-PlanAction `
                    'migrate' `
                    ([string]$projectAgents.id) `
                    $null `
                    $null `
                    ([string]$projectAgents.destination) `
                    'project-owned' `
                    'migrate-existing-agents' `
                    $true `
                    ([string]$managedAgents.destination)
                $MigrationCreatesProjectAgents = $true
            } else {
                $agentsHash = (Get-FileHash -Algorithm SHA256 -LiteralPath $managedAgentsPath).Hash.ToLowerInvariant()
                $migrationRelative = "AGENTS-project.migration-sha256-$agentsHash.md"
                $migrationPath = Join-PortablePath $ProjectRoot $migrationRelative 'AGENTS conflict migration destination'
                Assert-ContainedPhysicalPath $ProjectRoot $migrationPath 'AGENTS conflict migration destination'
                if (Test-Path -LiteralPath $migrationPath) {
                    if ((-not (Test-Path -LiteralPath $migrationPath -PathType Leaf)) -or
                        ((Get-FileHash -Algorithm SHA256 -LiteralPath $migrationPath).Hash.ToLowerInvariant() -cne
                        $agentsHash)) {
                        throw "Deterministic AGENTS migration destination is occupied: $migrationRelative"
                    }
                } else {
                    Add-PlanAction `
                        'migrate' `
                        ([string]$managedAgents.id) `
                        $null `
                        $null `
                        $migrationRelative `
                        'project-owned' `
                        'preserve-existing-agents-conflict' `
                        $true `
                        ([string]$managedAgents.destination)
                }
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
        Assert-ContainedPhysicalPath $ProjectRoot $destinationPath "entry '$($entry.id)' destination"
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

        Assert-ContainedPhysicalPath $ProjectRoot $destinationPath "entry '$($entry.id)' destination"
        $existingContent = [System.IO.File]::ReadAllText($destinationPath)
        $contentMatches = if ([string]$entry.id -ceq 'generated-installed-toolkit-manifest') {
            Test-InstalledManifestEquivalent $InstalledManifestDocument
        } else {
            $existingContent -ceq $desiredContent
        }
        if ($contentMatches) {
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

function Assert-PlannedDestinationTopology {
    foreach ($action in $script:PlanActions) {
        $destinationPath = Join-PortablePath `
            $ProjectRoot `
            ([string]$action.destination) `
            "planned action '$($action.entryId)' destination"
        Assert-DestinationTopology `
            $ProjectRoot `
            $destinationPath `
            "planned action '$($action.entryId)' destination"

        if ([string]$action.action -ceq 'backup') {
            $backupDestination = $BackupRoot
            foreach ($segment in ([string]$action.destination).Split('/')) {
                $backupDestination = Join-Path $backupDestination $segment
            }
            Assert-DestinationTopology `
                $BackupRoot `
                $backupDestination `
                "planned backup '$($action.entryId)' destination"
        }
    }
}

Assert-PlannedDestinationTopology

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

function Assert-InstallationPlanSemantics {
    param($Value)

    $blockCount = 0
    $actions = @($Value.actions)
    for ($index = 0; $index -lt $actions.Count; $index++) {
        $action = $actions[$index]
        if ([int]$action.sequence -ne ($index + 1)) {
            throw 'Installation-plan action sequences must be contiguous and ordered.'
        }
        if ([string]$action.newToolkitVersion -cne [string]$Value.toolkitVersion) {
            throw "Plan action '$($action.entryId)' new Toolkit version contradicts the plan."
        }
        $expectedOld = if ($null -eq $Value.installedToolkitVersion) {
            $null
        } else {
            [string]$Value.installedToolkitVersion
        }
        if (($null -eq $expectedOld) -xor ($null -eq $action.oldToolkitVersion)) {
            throw "Plan action '$($action.entryId)' old Toolkit version contradicts the plan."
        }
        if (($null -ne $expectedOld) -and
            ([string]$action.oldToolkitVersion -cne $expectedOld)) {
            throw "Plan action '$($action.entryId)' old Toolkit version contradicts the plan."
        }
        if (-not $EntryById.ContainsKey([string]$action.entryId)) {
            throw "Plan action references unknown entry '$($action.entryId)'."
        }
        $entry = $EntryById[[string]$action.entryId]
        if ([string]$action.action -ceq 'block') {
            $blockCount++
            if (([string]$action.destination -cne [string]$entry.destination) -or
                ([string]$action.ownership -cne [string]$entry.ownership)) {
                throw 'A block action must identify its canonical manifest entry.'
            }
            continue
        }
        if ([string]$action.action -ceq 'migrate') {
            if ([string]$action.reason -ceq 'migrate-existing-agents') {
                if (([string]$action.entryId -cne 'project-agents') -or
                    ([string]$action.targetSource -cne [string]$EntryById['managed-agents'].destination) -or
                    ([string]$action.destination -cne [string]$EntryById['project-agents'].destination) -or
                    ([string]$action.ownership -cne 'project-owned')) {
                    throw 'AGENTS migration action contradicts the manifest migration contract.'
                }
            } elseif ([string]$action.reason -ceq 'preserve-existing-agents-conflict') {
                if (([string]$action.entryId -cne 'managed-agents') -or
                    ([string]$action.targetSource -cne [string]$EntryById['managed-agents'].destination) -or
                    ([string]$action.destination -cnotmatch '^AGENTS-project\.migration-sha256-[0-9a-f]{64}\.md$') -or
                    ([string]$action.ownership -cne 'project-owned')) {
                    throw 'AGENTS conflict migration action contradicts the manifest migration contract.'
                }
            } else {
                throw "Unsupported migration reason '$($action.reason)'."
            }
            continue
        }
        $expectedSource = if ([string]$entry.kind -ceq 'copied') {
            [string]$entry.source
        } else {
            $null
        }
        $expectedGenerator = if ([string]$entry.kind -ceq 'generated') {
            [string]$entry.generator
        } else {
            $null
        }
        $sourceMismatch = (($null -eq $expectedSource) -xor ($null -eq $action.source)) -or
            (($null -ne $expectedSource) -and ([string]$action.source -cne $expectedSource))
        $generatorMismatch = (($null -eq $expectedGenerator) -xor ($null -eq $action.generator)) -or
            (($null -ne $expectedGenerator) -and
            ([string]$action.generator -cne $expectedGenerator))
        if (([string]$action.destination -cne [string]$entry.destination) -or
            ([string]$action.ownership -cne [string]$entry.ownership) -or
            $sourceMismatch -or $generatorMismatch) {
            throw "Plan action '$($action.entryId)' contradicts its manifest entry."
        }
    }
    if (([bool]$Value.canApply -and ($blockCount -ne 0)) -or
        ((-not [bool]$Value.canApply) -and
        (($blockCount -ne 1) -or ($actions.Count -ne 1)))) {
        throw 'Installation-plan canApply contradicts its block actions.'
    }
}

Assert-InstallationPlanSemantics $Plan

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
        [ValidateSet('PROPOSED', 'APPLIED', 'PRESERVED', 'UNCHANGED')]
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

    Assert-PortableRelativePath $Relative "action destination '$Relative'" -Destination
    $destination = Join-PortablePath $ProjectRoot $Relative "action destination '$Relative'"
    Assert-DestinationTopology $ProjectRoot $destination "action destination '$Relative'"
    $parent = Split-Path -Parent $destination
    New-Item -ItemType Directory -Force -Path $parent | Out-Null
    Assert-DestinationTopology $ProjectRoot $destination "action destination '$Relative'"
    [System.IO.File]::WriteAllText(
        $destination,
        $Content,
        [System.Text.UTF8Encoding]::new($false)
    )
}

function Backup-ProjectFile {
    param([string]$Relative)

    $source = Join-PortablePath $ProjectRoot $Relative "backup source '$Relative'"
    Assert-ContainedPhysicalPath $ProjectRoot $source "backup source '$Relative'"
    if (-not (Test-Path -LiteralPath $source -PathType Leaf)) { return }
    $destination = $BackupRoot
    foreach ($segment in $Relative.Split('/')) {
        $destination = Join-Path $destination $segment
    }
    Assert-DestinationTopology $BackupRoot $destination "backup destination '$Relative'"
    $parent = Split-Path -Parent $destination
    New-Item -ItemType Directory -Force -Path $parent | Out-Null
    Assert-DestinationTopology $BackupRoot $destination "backup destination '$Relative'"
    Assert-ContainedPhysicalPath $ProjectRoot $source "backup source '$Relative'"
    Copy-Item -LiteralPath $source -Destination $destination -Force
    Write-SdpAction 'APPLIED' "Backed up $Relative"
}

if ($Preview) {
    foreach ($action in $Plan.actions) {
        switch ([string]$action.action) {
            'preserve' { Write-SdpAction 'PRESERVED' "$($action.destination) ($($action.reason))" }
            'unchanged' { Write-SdpAction 'UNCHANGED' $action.destination }
            default { Write-SdpAction 'PROPOSED' "$($action.action) $($action.destination) ($($action.reason))" }
        }
    }
} else {
    $CurrentRepositoryPhysicalState = Get-PhysicalPathState `
        $RepositoryRoot `
        'SDP repository root' `
        -RequireExistingRoot
    $CurrentProjectPhysicalState = Get-PhysicalPathState `
        $ProjectRoot `
        'consuming project root' `
        -RequireExistingRoot
    $CurrentBackupPhysicalState = Get-PhysicalPathState $BackupRoot 'backup root'
    Assert-NoPhysicalTreeOverlap `
        $CurrentRepositoryPhysicalState `
        $CurrentProjectPhysicalState `
        'The consuming project and SDP source repository must be physically separate trees.'
    Assert-NoPhysicalTreeOverlap `
        $CurrentRepositoryPhysicalState `
        $CurrentBackupPhysicalState `
        'The backup root and SDP source repository must be physically separate trees.'
    Assert-PlannedDestinationTopology

    foreach ($action in $Plan.actions) {
        switch ([string]$action.action) {
            'backup' {
                Backup-ProjectFile ([string]$action.destination)
            }
            'create' {
                $sourcePath = Join-PortablePath $RepositoryRoot ([string]$action.source) "entry '$($action.entryId)' source"
                Assert-ContainedPhysicalPath $RepositoryRoot $sourcePath "entry '$($action.entryId)' source"
                $destinationPath = Join-PortablePath $ProjectRoot ([string]$action.destination) "entry '$($action.entryId)' destination"
                Assert-ContainedPhysicalPath $ProjectRoot $destinationPath "entry '$($action.entryId)' destination"
                if (Test-Path -LiteralPath $destinationPath) {
                    throw "Create destination appeared after planning: $($action.destination)"
                }
                Write-ProjectContent ([string]$action.destination) ([System.IO.File]::ReadAllText($sourcePath))
                Write-SdpAction 'APPLIED' "Created $($action.destination) ($($action.ownership))"
            }
            'replace' {
                $sourcePath = Join-PortablePath $RepositoryRoot ([string]$action.source) "entry '$($action.entryId)' source"
                Assert-ContainedPhysicalPath $RepositoryRoot $sourcePath "entry '$($action.entryId)' source"
                Write-ProjectContent ([string]$action.destination) ([System.IO.File]::ReadAllText($sourcePath))
                Write-SdpAction 'APPLIED' "Replaced $($action.destination)"
            }
            'generate' {
                $destinationPath = Join-PortablePath $ProjectRoot ([string]$action.destination) "entry '$($action.entryId)' destination"
                Assert-ContainedPhysicalPath $ProjectRoot $destinationPath "entry '$($action.entryId)' destination"
                if (([string]$action.reason -ceq 'missing-generated-target') -and
                    (Test-Path -LiteralPath $destinationPath)) {
                    throw "Generate destination appeared after planning: $($action.destination)"
                }
                Write-ProjectContent ([string]$action.destination) (Get-GeneratorContent ([string]$action.generator))
                Write-SdpAction 'APPLIED' "Generated $($action.destination)"
            }
            'migrate' {
                $sourcePath = Join-PortablePath $ProjectRoot ([string]$action.targetSource) 'migration targetSource'
                $destinationPath = Join-PortablePath $ProjectRoot ([string]$action.destination) 'migration destination'
                Assert-ContainedPhysicalPath $ProjectRoot $sourcePath 'migration targetSource'
                Assert-DestinationTopology $ProjectRoot $destinationPath 'migration destination'
                if (-not (Test-Path -LiteralPath $sourcePath -PathType Leaf)) {
                    throw "Migration targetSource is no longer a file: $($action.targetSource)"
                }
                if (Test-Path -LiteralPath $destinationPath) {
                    throw "Migration destination appeared after planning: $($action.destination)"
                }
                $parent = Split-Path -Parent $destinationPath
                New-Item -ItemType Directory -Force -Path $parent | Out-Null
                Assert-ContainedPhysicalPath $ProjectRoot $sourcePath 'migration targetSource'
                Assert-DestinationTopology $ProjectRoot $destinationPath 'migration destination'
                Copy-Item -LiteralPath $sourcePath -Destination $destinationPath
                Write-SdpAction 'APPLIED' "Migrated $($action.targetSource) to $($action.destination)"
            }
            'preserve' {
                Write-SdpAction 'PRESERVED' "$($action.destination) ($($action.reason))"
            }
            'unchanged' {
                Write-SdpAction 'UNCHANGED' $action.destination
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
