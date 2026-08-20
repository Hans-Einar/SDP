# SPR-SDP-005 Iterations and Slices

## ITR-SDP-005-001 — Inventory and repository evidence

Status: complete

### SLC-SDP-005-001 — Exact inventory and study contract

Status: complete

#### Goal

Define the reproducible activity-window method, enumerate every repository owned
by `Hans-Einar`, and classify each repository using actual commit evidence.

#### Why now

The inventory determines the required report set and must be stable before broad
repository study begins.

#### Expected files

- `SDP-usage-analysis/README.md`
- `SDP-usage-analysis/RepositoryInventory.md`
- this Sprint's execution records
- `Traceability/CurrentIndex.yaml`
- `Traceability/Relations.yaml`
- `Traceability/Ledger.ndjson`

#### Contract

- The fixed study window is the four calendar months ending at Issue #5's
  creation time: `2026-04-20T21:08:52Z` through `2026-08-20T21:08:52Z`,
  inclusive.
- Enumerate all repositories owned by `Hans-Einar`, including private or
  archived repositories visible to the authenticated account.
- A repository is in scope only when a commit reachable from its advertised Git
  refs has a committer timestamp inside the fixed window.
- Record the exact default-branch study commit separately from the activity
  commits used for inclusion.
- Record inaccessible or ambiguous cases as limitations and continue.

#### Invariants

- Repository `pushed_at`, creation date and preliminary search results are
  discovery inputs, never final inclusion evidence.
- No studied repository is modified.
- Inventory evidence is reproducible from recorded commands and identities.

#### Non-goals

- Drawing cross-repository conclusions before individual reports exist.
- Treating directory names alone as proof of process usage.

#### Verification

- Re-run the inventory derivation from the captured repository/ref/commit data.
- Confirm every in-scope repository has one planned report path.
- Confirm every exclusion states its evidence-backed reason.

#### Completion signal

`RepositoryInventory.md` names the considered, in-scope and excluded sets with
exact evidence, and a separate verification pass reproduces the classification.

## ITR-SDP-005-002 — Per-repository reports

Status: complete

### SLC-SDP-005-002 — Evidence-backed repository corpus

Status: complete

Create and Master-review one contract-complete Markdown report for every
in-scope repository, committing progressively.

## ITR-SDP-005-003 — Cross-repository synthesis

Status: active

### SLC-SDP-005-003 — Proposed future SDP workflow

Status: active

Synthesize the completed reports; do not concatenate Worker recommendations.

## ITR-SDP-005-004 — Validation and independent review

Status: planned

### SLC-SDP-005-004 — Corpus verification, adversarial review and closure

Status: planned

Run deterministic validation, resolve every Blocking/High/Medium review finding,
update traceability/handoff and stop for Steering Group review.
