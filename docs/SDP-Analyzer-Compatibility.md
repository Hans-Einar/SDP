# SDP-Analyzer Compatibility Contract

SDP-Analyzer should consume, but is not implemented by, this release.

## Toolkit repository inputs

- root `SDP.manifest.yaml` release/capability facts;
- `Toolkit/SDP-install.manifest.json` and its install/plan schemas;
- schema files and skill YAML front matter;
- root release notes, release records and Toolkit traceability.

## Installed consuming-project inputs

- `SDP/SDP-project.manifest.yaml`;
- generated `SDP/Framework/installed-toolkit.manifest.yaml`;
- installed skill YAML front matter;
- `SDP/RELEASE-NOTES.md`;
- CurrentIndex, Relations and the append-only Ledger;
- release, Fix, verification and review records;
- generated build-identity JSON;
- project-mode validator results when available.

The Analyzer must navigate Toolkit/project/skill versions, capabilities, release
targets, development identities, Sprint or Refactor, Iteration, Slice,
Fix/revision, verification, tags/commits, events, Unreleased changes and
migrations. It must flag unsupported schemas, stale Toolkit versions,
contradictory state, missing evidence and publication claims without real
identities.

Every Ledger event obeys the generic envelope. Canonical `release-*` events also
obey the release specialization; generic canonical families are `work-*`,
`review-*` and `verification-*`. Project-defined types use exactly
`x-<namespace>:<event-name>`, for example
`x-acme.example:deployment-approved`.

Unknown data is handled according to the applicable schema, not a blanket rule.
Strict objects reject unknown fields; extensible Relations categories/details
and Ledger payloads accept only the extensions their schema permits. Unknown
schema versions are unsupported rather than guessed.

Source archives without `.git` legitimately produce `sourceCommit: null`.
Analyzer must not treat null as a defect or infer a commit from an archive name.
For a dirty checkout, a non-null `sourceCommit` is only the available `HEAD`
baseline and must not be treated as proof that installed bytes equal that commit.
Markdown analysis is limited to canonical structure and deterministic links; it
must not infer substantive completion from arbitrary prose.

All contract paths use normalized `/` separators. Paths beginning with `SDP/`
resolve from the project root; other record paths resolve from the project
`SDP/` directory, with a project-root fallback for repository-level records.
Canonical governed records use `.yaml`; a `.yml` shadow is invalid.

The supported `gh-sdp` extension surface includes additional Sprint/Iteration/
Slice status values, `acceptanceCriteria`, `requirementRefs`, `verificationRefs`,
`reviewRefs`, and `targetStatus`/list-shaped link values in Relations. Analyzer
must still require reciprocal Slice/Fix review and verification links,
review-resolution links, release/migration links, and
Ledger subjects that resolve to a governed entity. Publication state, record
identity, SemVer version and tag must agree rather than being inferred from one
field alone.
