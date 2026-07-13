# Installer Migration Guide

Let `$SdpSource` be an independent clone or extracted source-archive root and
`$Project` the consuming project:

```powershell
$SdpSource = 'C:\path\to\extracted-or-cloned-SDP'
$Project = 'C:\path\to\Project'
```

Generate the portable, deterministic plan first:

```powershell
& "$SdpSource\Toolkit\scripts\Install-SDP.ps1" `
  -ProjectRoot $Project `
  -PlanJson
```

`-PlanJson` writes only JSON conforming to
`Toolkit/schemas/SDP-install-plan.schema.json`. It makes zero target mutations.
Use `-Preview` when human-readable mutation-free action output is preferred.

Apply the update:

```powershell
& "$SdpSource\Toolkit\scripts\Install-SDP.ps1" `
  -ProjectRoot $Project
```

Initialize missing neutral lifecycle structure when wanted:

```powershell
& "$SdpSource\Toolkit\scripts\Install-SDP.ps1" `
  -ProjectRoot $Project `
  -InitializeProjectStructure
```

The canonical inventory and behavior are declared by
`Toolkit/SDP-install.manifest.json`; PowerShell does not own a parallel file
list. Default entries install managed contracts and missing neutral project,
release and traceability foundations. Initialization additionally selects only
the missing neutral files under `Toolkit/project-templates/`.

Initialization never reads the repository's root lifecycle folders as payload.
It does not copy or propose `SDP/Releases/REL-0.2.0.yaml`, the Toolkit's active
Sprint/Refactor, populated release notes, CurrentIndex, Ledger, review or
verification evidence. Repeating initialization is idempotent.

Use `-ForceManagedFiles` only when intentionally restoring same-version
Toolkit-managed files. It never replaces project-owned release notes, project
manifest, AGENTS-project, lifecycle documents or traceability.

Before mutation the installer:

- validates the complete closed-world installation contract before inspecting
  target files;
- rejects symlinks/reparse points in the source, project and backup roots and
  every existing source/destination ancestor; this is rechecked at write time;
- parses installed/project manifests with a strict YAML subset, accepting only
  a root mapping of scalar mappings/lists and rejecting duplicate required keys,
  nested shadows, malformed scalars and unsupported document features;
- accepts empty and non-empty SDP directories;
- treats an installation without installed facts as the pre-versioning baseline;
- stops safely on malformed or unsupported schemas;
- refuses to downgrade by full SemVer 2.0 precedence; prereleases sort below the
  matching final, while build metadata changes identity but not precedence;
- compares complete paths so a sibling such as `SDP-Analyzer` is valid.

During apply it refreshes only entries whose policy permits it, backs up changed
managed files before replacement, creates missing project-owned files, reports
actions and preserves the original installation timestamp on a same-version
reinstall. Backups default to `SDP/.sdp-backups/<UTC timestamp>/` and may be
redirected with `-BackupRoot`.

During the one-time AGENTS migration, prior project-specific `AGENTS.md` content
is copied to `AGENTS-project.md`. If that target already exists, the prior
content is preserved at project root as
`AGENTS-project.migration-sha256-<content-hash>.md`. The JSON plan represents
both operations as `migrate` with `targetSource: "AGENTS.md"` and the exact apply
destination; apply does not invent a timestamped path. Existing project-owned
files are never overwritten.

An independently serialized installed manifest is not refreshed merely because
its mapping order or YAML quoting differs. The installer compares declared facts
semantically, while retaining ordered capability and exact dynamic identity
meaning.

An extracted normal GitHub source archive is supported without `.git`. In that
case installed facts use `sourceCommit: null`; the installer does not infer a
commit from the archive name.
