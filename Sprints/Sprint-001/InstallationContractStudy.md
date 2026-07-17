# Installation contract v1 study

Status: decided for `SPS-001`
Decision owner: SDP Master
Decision date: `2026-07-13`

## Evidence

- `Toolkit/scripts/Install-SDP.ps1` currently embeds versions, managed source
  discovery, destinations, ownership, generated installed-manifest content and
  initialization sources in PowerShell.
- `-InitializeProjectStructure` recursively copies live root folders. The exact
  `gh-sdp` PR #1 evidence shows that a repeat preview proposes
  `SDP/Releases/REL-0.2.0.yaml` after the consuming project removed the
  Toolkit-specific seed.
- The current validator assumes the SDP Toolkit repository layout and only
  validates CurrentIndex and Relations as generic YAML objects.
- PowerShell 5.1 has a built-in JSON parser but no built-in YAML parser. Requiring
  Python merely to install would regress the supported compatibility entry point.
- A normal GitHub source archive preserves repository-relative files but has no
  `.git` directory or trustworthy source commit.

## Alternatives

### 1. Explicit manifest entry for every installed file

Strengths: deterministic inventory, stable IDs, simple implementations, simple
archive resolution, direct ownership/policy review and strong detection of
unlisted managed files. Cost: adding a payload file requires an intentional
manifest edit. This cost is desirable because distribution changes are public
contract changes.

### 2. Controlled directory or glob expansion

Strengths: concise for skills and Framework payload. Risks: glob syntax,
separator, case, symlink and ordering differences become part of the portable
contract; a newly added repository file may become installable accidentally.
Rejected for v1. The small current inventory does not justify the ambiguity.

### 3. Separate manifests for managed payload, project templates and generators

Strengths: physical categories are obvious. Risks: external clients must merge
several versioned authorities and cross-manifest policy can drift. Rejected as
the canonical model. Physical source trees will remain separated, while one
manifest lists every entry and category.

### 4. PowerShell remains canonical and a manifest only describes it

Strengths: smallest installer change. Risks: this is ceremonial metadata; an
external client still has to infer executable behavior and disagreement can
remain silent. Rejected.

### 5. Manifest is authoritative and PowerShell consumes it

Strengths: one source of truth for file inventory and portable policy; both
PowerShell and future Go can implement the same contract. Requires a focused
installer refactor. Selected.

### 6. Generate implementation-specific plans from one manifest

Strengths: normalized actions are directly conformance-testable. A generated
script or checked-in derived file would introduce synchronization risk. Selected
only as a runtime deterministic preview-plan representation; no checked-in
implementation-specific inventory will be generated.

## Decision

Create `Toolkit/SDP-install.manifest.json` with schema version `1.0` and an
accompanying Draft 2020-12 JSON Schema. JSON is selected deliberately because it
is platform-neutral, directly schema-validatable, natively readable by supported
PowerShell 5.1 and Go, and does not add an installer-time parser dependency.

The manifest will explicitly list every copied or generated entry. It will
define stable ID, source or generator, normalized destination, ownership,
installation phase/policy, refresh behavior, backup behavior, force behavior and
governing schema/capability where applicable. It will also declare excluded live
repository areas. Paths resolve from the archive repository root identified
relative to the manifest, never from the process working directory.

Physical inputs will be separated as follows:

- Toolkit-managed copies remain under `Toolkit/payload/` and `Toolkit/skills/`.
- Neutral, project-owned initialization inputs live under
  `Toolkit/project-templates/`.
- Dynamic installed facts and the empty ledger are represented as generators.
- Live root repository records are never installer sources.

PowerShell will enumerate entries from the manifest and expose a deterministic
JSON preview plan. The default human interface remains compatible. The plan is a
normalized result, not a second source list, and identifies action, entry ID,
source or generator, destination, ownership, reason, mutation flag and relevant
old/new Toolkit versions.

## Automatic disagreement detection

Toolkit validation and tests will fail when:

- the manifest or an entry violates its schema;
- Toolkit versions disagree;
- an entry ID or destination is duplicated;
- a path is absolute, traverses, is not normalized or references a missing
  source;
- a managed file under the canonical managed roots is not listed;
- an installer source enters a declared live-record exclusion;
- PowerShell installation/plan behavior differs from manifest fixtures; or
- an initialized project contains repository-owned active state.

The installer therefore has no independent file list to compare manually.

## Project validation schema boundary

Add strict schemas for CurrentIndex's stable outer shape and the Ledger event
envelope. Keep Relations extensible at category and relation-detail boundaries
while validating known path-bearing fields. Retain release events as a stricter
specialization. Work, review and verification events use the generic envelope in
v1; richer payload schemas remain future work until repeated semantics exist.
Project-defined event types are permitted only through a documented namespaced
extension convention and still must satisfy the envelope.

Markdown validation is intentionally limited to canonical-path existence,
release-note structure and deterministic references. No claim of full arbitrary
Markdown semantic validation is made.

## Release archive decision

The normal GitHub source archive is sufficient because all contract sources are
repository-relative and `.git` is optional. When Git metadata is unavailable,
the generated installed manifest records `sourceCommit: null`. No dedicated
release asset will be added in this Slice unless implementation evidence
contradicts this decision.
