---
skillId: sdp-versioning
skillVersion: 1.0.0
minimumToolkitVersion: 0.2.0
capabilities:
- sdp.versioning.select
- sdp.versioning.development-identity
- sdp.versioning.validate
compatibilityNotes: Initial formal version for SDP Toolkit 0.2.0.
---

# SDP Versioning

Use this skill to select release increments and manage development identity.

1. Inspect the latest released version and public compatibility impact.
2. Select MAJOR, MINOR or PATCH using SemVer; never use work size as the reason.
3. Keep Sprint/Refactor, Iteration, Slice/Fix and revision in separate fields.
4. Validate that exactly one of Sprint and Refactor is active when either applies.
5. Allow a revision only for a bounded correction inside the same planned work
   contract; require a new Slice/Fix when scope or architecture changes.
6. Generate display identity from machine-readable fields and Git state.
7. Prevent edits that silently alter a released version or released notes.
8. Record the selected target and rationale in traceability.

An unreleased build must never present itself as the official release.
