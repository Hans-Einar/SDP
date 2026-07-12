# Distribution And Upgrades

## Goals

- reuse SDP across existing projects without nested Git repositories
- preserve project-owned method and release records
- refresh versioned managed contracts predictably
- detect and migrate supported older installations additively
- expose exact installed versions for audit and SDP-Analyzer

## Ownership classes

### Toolkit-managed

- root `AGENTS.md`
- `.codex/skills/sdp-*/SKILL.md`
- `SDP/Framework/*`, including installed Toolkit facts and templates

The installer may replace these during a Toolkit upgrade. Changed managed files
are backed up first. Same-version local differences are preserved unless
`-ForceManagedFiles` is supplied.

### Project-owned

- `AGENTS-project.md`
- `SDP/AGENT-REMINDERS.md`
- `SDP/SDP-project.manifest.yaml`
- `SDP/RELEASE-NOTES.md`
- Mandate through Implementation, Sprints, Refactors and Fixes
- Releases, CodeReview, Verification, Traceability and Instructions

These files are created only when missing and are never overwritten by normal or
forced installation.

## Installed identity

`SDP/Framework/installed-toolkit.manifest.yaml` records installed Toolkit,
Framework, AGENTS contract, installer and skill versions, capabilities,
installation timestamp and source commit when available. Dynamic Git/build facts
are generated separately and never maintained by hand.

## Migration transaction

1. Run with `-Preview`.
2. Detect the existing installed-manifest schema before mutation.
3. Treat a missing manifest as the supported pre-versioning baseline.
4. Stop on malformed or unsupported schemas.
5. Preserve or migrate old project AGENTS instructions.
6. Back up and refresh managed files when upgrading.
7. Create missing project manifest, release notes and traceability contracts.
8. Report proposed, applied, preserved and unchanged files.
9. Repeating the same install produces no content changes.

Backups default to `SDP/.sdp-backups/<UTC timestamp>/`.

## Repository layout

Keep one upstream clone outside consuming repositories. Do not clone it into an
existing project `SDP/` directory. A sibling such as `SDP-Analyzer` is allowed;
a project physically inside the Toolkit repository is rejected using complete
path-segment comparison.

Submodules, subtrees and release ZIPs remain possible distribution alternatives,
but the installer is canonical because shared managed files and project-owned
records must coexist safely.
