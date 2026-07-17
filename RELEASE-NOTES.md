# SDP Toolkit Release Notes

## [Unreleased]

Release-Date: unreleased

### Added

- [REL-0.2.0] First-class Toolkit, project, skill, release and development identity manifests.
- [REL-0.2.0] Release preparation, versioning, auditing and verification skills.
- [REL-0.2.0] Deterministic release gates, traceability event schemas and two-phase Git/GitHub publication.
- [REL-0.2.0] Safe additive installer migration, preview mode, backups and fixture-based validation.
- [REL-0.2.0] Portable generated build metadata suitable for Vite and non-Vite applications.
- [SPS-001] Canonical schema-validated, platform-neutral installation manifest and deterministic JSON plan for PowerShell and independent clients.
- [SPS-001] Separate Toolkit-repository and consuming-project validator modes with reusable CurrentIndex, Relations and generic Ledger schemas.
- [SPS-001] Source-archive installation contract that requires no `.git` and records an unavailable source commit as null.
- [SPS-001] Versioned language-neutral install-v1 conformance package with 17 portable scenarios and committed normalized plan/failure outcomes.

### Changed

- [REL-0.2.0] Canonical skills now carry machine-readable version metadata.
- [REL-0.2.0] Traceability templates model releases and bounded Fix records as first-class entities.
- [REL-0.2.0] Toolkit-managed AGENTS and Framework guidance now require release/version awareness.
- [SPS-001] PowerShell installation inventory and policy now derive from the canonical JSON contract; neutral project templates are physically separate from live repository records.
- [SPS-001] Installation plan v1 now declares and enforces migration-first, manifest-array canonical ordering with adjacent backup/mutation pairs.

### Fixed

- [SPS-001] Structure initialization no longer proposes or copies Toolkit `REL-0.2.0`, active Sprint/Refactor, populated Ledger, review or verification state into consuming projects.
- [SPS-001] Installation manifests, plans, portable paths, SemVer ordering, YAML preflight and physical source/target containment now fail closed before mutation.
- [SPS-001] AGENTS migration now has an explicit deterministic plan/apply contract, including content-hash conflict destinations.
- [SPS-001] AGENTS conflict preservation now hashes and compares exact bytes, handles existing deterministic destinations idempotently, and rejects collision or plan/apply races with stable failure classes.
- [SPS-001] Project validation now enforces governed-path resolution, canonical YAML names, reciprocal traceability, Ledger subjects and coherent publication identities while accepting the pinned supported `gh-sdp` extension surface.

### Migration

- [REL-0.2.0] Existing consuming projects gain missing manifests and release templates without replacing populated project-owned files.
- [REL-0.2.0] Managed-file changes are backed up before replacement; unsupported manifest schemas stop safely.
- [SPS-001] Project-owned content remains missing-only and is preserved even under force; extracted archives may truthfully generate installed facts with `sourceCommit: null`.
