# SDP Framework

Framework-Version: 1.0.0
AGENTS-Contract-Version: 1.0.0

This directory is Toolkit-managed. Project-specific Mandate, Study, Requirements,
Architecture, Design, Implementation, Sprints, Refactors, Fixes, Review,
Verification, release notes and Traceability remain project-owned.

Canonical installed facts are generated in
`Framework/installed-toolkit.manifest.yaml`. Project release and development
state lives in `SDP-project.manifest.yaml`. Dynamic Git/build facts are generated,
not manually duplicated.

Core rules:

- repository documents are authoritative
- design horizontally and implement vertically
- use SemVer only for public releases and Toolkit releases
- keep Sprint/Refactor, Iteration, Slice/Fix and revision as separate coordinates
- one active Slice or bounded Fix at a time
- Workers implement bounded work; Reviewers use fresh independent context
- completion requires real verification and current traceability
- tag and GitHub Release events may be recorded only after they exist
- publication requires explicit human authorization
- stop at the work boundary
