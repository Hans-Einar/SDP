# Traceability Release Events

Release IDs use `REL-<version>`, for example `REL-0.2.0`; the version is the
full SemVer identity when prerelease or build metadata is present. Fix IDs use
`FIX-<target-version>-NNN`; emergency fixes may use `HOTFIX-...`.

`CurrentIndex.yaml` exposes current release and development coordinates.
`Relations.yaml` links releases to mandate, requirements, Sprints/Refactors,
Iterations, Slices/Fixes, verification, reviews, notes, commits, tags, GitHub
Releases and migrations.

Governed record paths use canonical `.yaml` names and must resolve. Slice/Fix
review and verification links, review resolution links, and release/migration
links are reciprocal; every Ledger subject resolves to a governed entity. A
release record ID is exactly
`REL-<version>`; publication tags are exactly `v<version>` and must agree across
the release record, Relations and release events.

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
