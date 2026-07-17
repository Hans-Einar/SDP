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
The plan declares `orderingPolicy: "migration-first-manifest-order-v1"`:
migration is first, ordinary entries retain installation-manifest array order,
and every required backup is immediately adjacent to its matching replacement
or regeneration. Sequence numbers are assigned last and begin at one.

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
- compares OS-backed root identities, rejects source/project or source/backup
  overlap in either ancestor direction across local, UNC, extended and available
  short-name aliases, and fails closed when identity or share ancestry is
  insufficient;
- validates raw extended drive/UNC source, project and backup roots before
  removing `\\?\` prefixes, rejecting traversal, empty/doubled or mixed
  separators and trailing-space/dot segments that Win32 would reinterpret;
- rejects case-insensitive destination ancestor/descendant pairs and any
  existing non-directory destination ancestor before an applicable plan;
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
destination. The hash is lowercase SHA-256 of the exact source bytes, without
text decoding or line-ending normalization, and the action records it as
`targetSourceSha256` with `destinationPrecondition: "absent"`. Apply does not
invent a timestamped path. Existing project-owned files are never overwritten.

For the deterministic conflict destination:

- absent means one migration action;
- an ordinary file with identical exact bytes means preservation is already
  complete and no migration action is emitted;
- an ordinary file with different bytes is fatal
  `agents-migration-destination-content-mismatch`;
- a directory, link, symlink, reparse point or unsupported object is fatal
  `agents-migration-destination-unsupported-object`.

Immediately before migration apply, the installer rechecks the exact source
hash and absent destination. Source drift is fatal
`agents-migration-source-changed`; a destination that appeared or changed is
fatal `agents-migration-destination-changed`. Exclusive creation prevents a
race from overwriting project-owned content. Migration-first ordering ensures
these checks and the preservation write happen before ordinary installation
mutations. If a later ordinary action fails, an already completed preservation
file can remain; the installer does not promise full transactional rollback.

Supported target-derived conditions such as an unsupported installed/project
schema or downgrade return a valid blocked plan (`canApply: false`) with exactly
one block action at sequence 1. Malformed installation contracts and the fatal
AGENTS filesystem cases above emit no plan and use the stable machine-readable
failure classes documented in `Installation-Contract.md`.

An independently serialized installed manifest is not refreshed merely because
its mapping order or YAML quoting differs. The installer compares declared facts
semantically, while retaining ordered capability and exact dynamic identity
meaning.

An extracted normal GitHub source archive is supported without `.git`. In that
case installed facts use `sourceCommit: null`; the installer does not infer a
commit from the archive name. For a Git checkout, a non-null `sourceCommit` is
the available `HEAD` baseline. A dirty checkout can therefore install bytes that
differ from that commit; the field is not a content attestation.

The portable contract fixtures are in `Toolkit/conformance/install-v1/`.
Independent clients consume `scenarios.json` and the checked-in outcomes without
executing PowerShell. Maintainers can run the reference comparison with:

```powershell
python Toolkit\conformance\install-v1\run_conformance.py --powershell powershell
```

Expected outcomes are never regenerated during normal tests. Candidate updates
require the explicit `--write-candidates` option and review of the resulting
contract diff.
