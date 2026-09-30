# SDP Toolkit Release Notes

## [Unreleased]

Release-Date: unreleased

### Added

- [PLAN-SDP-0012] Go design-core 0.6 supports source-owned System composition, includes/path membership, reusable file ASTs and graph-aware CLI/viewpoints/broker/SDPTool consumers. Existing 0.5, action and class profiles remain separate.
- [MAINT-SDP-0011] Manual Sessions guide/template are included in installation, with sdp.sessions.manual.v1 capability and preservation of existing project Sessions. Automatic transcript capture and event timelines remain separate work.
- [MAINT-SDP-0011] sdptool release-log generates immutable per-version Markdown from canonical release notes, with check mode. ReleaseChecklist covers inventories, predecessor descriptors, signed artifacts and consumer readiness.

### Changed

- [PLAN-SDP-0013] Go SDPTool renders readable command output by default, including the navigation tree. --json explicitly preserves machine schemas; compiled-in adapters need no shell formatter. Machine consumers must opt in.

### Migration

- [PLAN-SDP-0013] The default CLI output change is incompatible with existing implicit JSON clients. The next combined candidate is proposed as 2.0.0 under the current version contract, superseding the earlier additive 1.1.0 proposal. Update gh-sdp bootstrap and XFMD machine invocations before rollout. No new release is published.
- [MAINT-SDP-0011] The earlier additive preparation proposed SDP 1.1.0; SDL 0.6 is a language profile. The candidate declares the published 1.0.0 descriptor as an additional supported predecessor. Publication and a gh-sdp default update are not yet performed.

## [1.0.0] - 2026-09-28

### Removed

- [MAINT-SDP-0010] Retired installation scripts, profile-artifact generator and shell-driven test/recovery jobs are removed, including the executable bootstrap archive copy. Go SDPTool is the only current install/upgrade engine. Pending legacy journals still block unsafe migration and require a separately assessed recovery; they are not silently resumed.

### Changed

- [MAINT-SDP-0010] Go release authoring owns its explicit payload inventory in SDPTool/profiles/payload.json. Current CI uses Go for installation, signatures and recovery; Python remains for schema/document checks.

### Fixed

- [MAINT-SDP-0009] Current project templates have one canonical Template/sdp-root home; legacy install-v1 templates are archived separately with preserved payloads.
- [MAINT-SDP-0009] Installed source guidance defines SDP/SDL/<System> ownership, container and shared-library boundaries, SDUI screen homes, and authored versus generated documentation.

### Added

- [PLAN-SDP-0005] Repository-only MVP1 source organization and three static SDUI examples with reproducible previews. This does not add experimental MVP1 parsing or runtime integration to the released tools.

### Migration

- [MAINT-SDP-0010] Signed upgrades support SDP 0.2.0 and 0.2.1. Missing SDL guides initialize; existing project-owned documents and model locations remain unchanged. Managed agent instructions link the source convention. gh-sdp 0.1.2 selects this release by default.

## [0.2.1] - 2026-09-27

### Fixed

- [KB-SDP-034] Foreign KanBan primary references are explicitly external and unverified. They no longer disable an otherwise valid board. Local references and history remain validated; no external checkout or network lookup is performed.

### Migration

- [MAINT-SDP-0008] Update the engine to receive the navigation fix. The gh-sdp client patch pins SDP 0.2.1. No project document or receipt-format migration is required; existing 0.2.0 installations remain valid.

## [0.2.0] - 2026-09-27

### Added

- [PLAN-SDP-0003] Go SDPTool install/upgrade preview, root-bound apply, journaled recovery and receipt 3.0. Shared signed distribution bootstrap supports the thin gh-sdp client.
- [MAINT-SDP-0007] First signed Linux amd64 binary/descriptor distribution, compiled publisher trust and exact release selection. Other native platforms are not advertised.
- [MAINT-SDP-0007] Shared five-phase SDP layout, typed plans, project-management history and canonical agent skills are delivered through the Go profile.


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

- [SK1] One canonical root Skills/ collection, native metadata, profile-aware routing and complete installed references replace both maintained legacy collections.

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

- [SK1] Existing installed skill paths remain .codex/skills. Same-Toolkit-version skill contract changes require reviewed ForceManagedFiles application with backups; unforced transitions fail before writes.

- [REL-0.2.0] Existing consuming projects gain missing manifests and release templates without replacing populated project-owned files.
- [REL-0.2.0] Managed-file changes are backed up before replacement; unsupported manifest schemas stop safely.
- [SPS-001] Project-owned content remains missing-only and is preserved even under force; extracted archives may truthfully generate installed facts with `sourceCommit: null`.

- [MAINT-SDP-0007] Manual installations require a fresh explicit adoption manifest. Project files are preserved; a legacy pending journal must be recovered by its original installer. Production signatures require the published client/engine; test keys never establish production provenance.
