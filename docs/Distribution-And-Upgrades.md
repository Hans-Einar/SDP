# Distribution And Upgrades

## Goals

- reuse SDP across existing projects without nested Git repositories;
- preserve project-owned method, work, release and evidence records;
- refresh versioned managed contracts predictably;
- expose a portable contract to PowerShell and independent clients;
- support clones and ordinary extracted GitHub source archives;
- expose exact installed versions for audit and SDP-Analyzer.

## Canonical boundary

`Toolkit/SDP-install.manifest.json` is canonical for installable files,
generators, destinations, ownership and policies. The PowerShell installer is a
supported reference and compatibility consumer of that contract, not a second
source of truth. Future external clients must conform to the same schema and
observable behavior.

## Ownership and source classes

### Toolkit-managed

Explicit manifest entries install root `AGENTS.md`, selected `SDP/Framework/`
files, generated installed facts and `.codex/skills/sdp-*/SKILL.md`. Changed
managed files refresh only under their declared policy and are backed up when
the policy requires it. Same-version local differences are preserved for
`upgrade-or-force` entries unless `-ForceManagedFiles` is supplied.

### Project-owned

`AGENTS-project.md`, `SDP/AGENT-REMINDERS.md`, the project manifest and release
notes, lifecycle documents, work/release/review/verification records and all
project traceability are project-owned. They are created only when missing and
never overwritten by normal, forced or repeated installation.

### Neutral project templates

Canonical missing-only seeds live under `Toolkit/project-templates/`. Default
seeds include project/release/traceability foundations. Lifecycle, operating
README and document-guide seeds are selected only by
`-InitializeProjectStructure`. Once created they are project-owned.

### Repository instance records

Root `Sprints/`, `Releases/`, `Traceability/`, `CodeReview/`, `Verification/`
and related lifecycle folders are the SDP Toolkit repository's live state. They
are excluded installation sources. In particular, Toolkit `REL-0.2.0`, its
active Sprint and its Ledger history never seed a consuming project.

## Installed identity

`SDP/Framework/installed-toolkit.manifest.yaml` is generated from the installation
manifest's declared facts. It records Toolkit, Framework, AGENTS contract,
installer and skill versions, capabilities, installation time and source commit
when trustworthy. In a source archive without Git metadata, `sourceCommit` is
truthfully null. Dynamic build identity is generated separately.

## Migration transaction

1. Produce `-PlanJson` or run `-Preview`; neither mutates the target.
2. Validate source and installed/project schema compatibility before mutation.
3. Treat a missing installed manifest as the supported pre-versioning baseline.
4. Stop on malformed/unsupported schemas or a Toolkit downgrade.
5. Preserve or migrate pre-SDP project AGENTS instructions.
6. Back up and refresh only managed entries allowed by policy.
7. Create only missing neutral project-owned defaults or requested structure.
8. Report planned/applied/preserved/unchanged behavior deterministically.
9. Repeating the same operation produces no project-owned content change.

Backups default to `SDP/.sdp-backups/<UTC timestamp>/`. Timestamps and concrete
backup-directory names are intentionally outside the portable JSON plan.

## Source layout and archives

Keep a clone outside consuming repositories, or extract a normal GitHub Release
source archive into any temporary directory. Locate its root by
`Toolkit/SDP-install.manifest.json`; do not assume the archive folder name. A
conforming implementation requires no `.git`, fixed absolute path, Windows path
semantics or PowerShell to interpret the JSON contract.

A sibling such as `SDP-Analyzer` is allowed. A target equal to or physically
inside the Toolkit source is rejected using complete path comparison.

The ordinary GitHub source archive is sufficient for this contract. No custom
release asset is required unless future verification demonstrates a real gap.
