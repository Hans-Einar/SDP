# SDP Toolkit Release Notes

## [Unreleased]

Release-Date: unreleased

### Added

- [REL-0.2.0] First-class Toolkit, project, skill, release and development identity manifests.
- [REL-0.2.0] Release preparation, versioning, auditing and verification skills.
- [REL-0.2.0] Deterministic release gates, traceability event schemas and two-phase Git/GitHub publication.
- [REL-0.2.0] Safe additive installer migration, preview mode, backups and fixture-based validation.
- [REL-0.2.0] Portable generated build metadata suitable for Vite and non-Vite applications.

### Changed

- [REL-0.2.0] Canonical skills now carry machine-readable version metadata.
- [REL-0.2.0] Traceability templates model releases and bounded Fix records as first-class entities.
- [REL-0.2.0] Toolkit-managed AGENTS and Framework guidance now require release/version awareness.

### Migration

- [REL-0.2.0] Existing consuming projects gain missing manifests and release templates without replacing populated project-owned files.
- [REL-0.2.0] Managed-file changes are backed up before replacement; unsupported manifest schemas stop safely.
