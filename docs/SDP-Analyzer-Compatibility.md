# SDP-Analyzer Compatibility Contract

SDP-Analyzer should consume, but is not implemented by, this release.

Machine-readable inputs:

- root `SDP.manifest.yaml`
- project `SDP/SDP-project.manifest.yaml`
- generated `SDP/Framework/installed-toolkit.manifest.yaml`
- skill YAML front matter
- `RELEASE-NOTES.md`
- `Traceability/CurrentIndex.yaml`
- `Traceability/Relations.yaml`
- append-only `Traceability/Ledger.ndjson`
- release, Fix, verification and review records
- generated build-identity JSON

The Analyzer must navigate Toolkit/project/skill versions, release targets,
development identities, Sprint or Refactor, Iteration, Slice, Fix/revision,
verification, tags/commits, release events, Unreleased changes and migrations.
It must flag unsupported schemas, stale Toolkit versions, contradictory state,
missing evidence and publication claims without corresponding identities.

Unknown optional fields should be preserved or ignored safely. Unknown schema
major versions must be reported as unsupported rather than guessed.
