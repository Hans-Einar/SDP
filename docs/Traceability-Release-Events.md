# Traceability Release Events

Release IDs use `REL-MAJOR.MINOR.PATCH`, for example `REL-0.2.0`.
Fix IDs use `FIX-<target-version>-NNN`; emergency fixes may use `HOTFIX-...`.

`CurrentIndex.yaml` exposes current release and development coordinates.
`Relations.yaml` links releases to mandate, requirements, Sprints/Refactors,
Iterations, Slices/Fixes, verification, reviews, notes, commits, tags, GitHub
Releases and migrations.

Allowed append-only release transitions include:

- `release-planned`
- `release-version-selected`
- `release-candidate-opened`
- `release-verification-completed`
- `release-approved`
- `release-tag-created`
- `release-published`
- `release-yanked`
- `release-migration-applied`
- `release-notes-corrected`

Only real transitions are appended. `release-tag-created` requires an existing
tag and commit; `release-published` requires the real GitHub Release URL. Failed
attempts may be recorded as separate failure events but never rewritten into
success.
