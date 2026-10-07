# BP2-C — producer API and assignment handoff

Contract candidate, 2026-10-07. PLAN-SDP-0001 / KB050. This specifies a future Go
producer; none of the proposed APIs or CLI commands below is shipped by this plan.
Read Contract.md and Selection-and-Evidence.md for authority and selection rules.

## Ownership and reuse

| Owner | Responsibility | Existing boundary |
| --- | --- | --- |
| SDPTool/model | Capture WORK/frozen artifact bytes and lifecycle identity | Snapshot returns SourceView with owned Files, LiveDigest and Preliminary |
| SDL sourcegraph/parser | Compile the captured selected source graph | SyntaxCache.Parse, sourcegraph.Compile, Snapshot.Model/Revision/Sources |
| Proposed SDL/go/blueprint package | Semantic delta, relation-policy closure, obligations, evidence associations | New pure library; no filesystem, Git, subprocess or SDUI parser |
| SDL documents/viewpoint | Render checked blueprint facts and annotations | Reuse source escaping, diagram primitives and bundle/link verification |
| SDPTool blueprint adapter | Resolve refs/entrypoints, capture task/check receipts, enforce publication constraints | Thin facade over libraries; use existing presentation registry |
| Consumer / reviewer | Accept useful scope, interpret unknowns, execute authorized work and review code | XFMD or Markdown reader; no required GUI integration |

Do not place semantic graph analysis in ModelGovernance or implement a second SDL
parser. SDPTool may use a standalone SDL adapter later; no extra binary is required
for the initial increment. SDUI source bytes remain in captured provenance, with
changed SDUI identified as outside the structural analyzer's coverage. Do not claim
semantic SDUI comparison until a separately scoped adapter exists.

## Proposed in-memory API

Conceptual Go API (names/types finalized during implementation):

```go
Analyze(ctx context.Context, now, target CheckedInput, task Task, policy Policy) (Analysis, error)
BuildBundle(ctx context.Context, analysis Analysis, receipts []CheckReceipt) (Bundle, error)
```

CheckedInput binds a read-only sourcegraph snapshot, complete captured source map,
artifact UUID/kind/head, model-area source digest and selected entrypoint. The
sourcegraph revision describes reachable SDL sources; the model-area digest covers
all captured files. They are different identities and must not be substituted.
Store parser/tool identity, profile and System identity as separate fields.
Reject incompatible Systems/profiles rather than treating names as cross-system IDs.
For the first release support design-core/0.6 only; report unsupported other profiles.

Task has stable task/card reference, revision/content digest, authored intent,
allowed semantic changes, exact allowed code/model paths, explicit context/exclusions,
PRESERVE obligations and unknown dispositions. A task file is an authored work
specification, not a manually maintained source inventory. It can be supplied as
bytes by a host. Every field affecting selection must participate in bundle identity.
CheckReceipt is evidence supplied as data: never execute a command just because it
appears in a blueprint, tag or SDL source. General code mapping remains KB004.

Analysis returns sorted typed changes, selected facts with per-side source spans,
inclusion reasons, cut boundaries, unknowns, obligation evaluations and coverage.
Fail closed for unsupported relation semantics. Node/edge IDs use stable semantic
keys; NOW-only removals remain visible. Annotation text includes CHANGE, PRESERVE,
CONTEXT and UNKNOWN; use color only as a supplement. JSON result is typed, with a
Go human-readable formatter selected unless --json is supplied.

## Status and authority

Keep these independent: preliminary input, source validation, analysis completeness,
obligation result, evidence applicability and authorization to implement.

- Invalid/incompatible sources, stale requested identity, unknown semantic policies,
  conflicting constraints, limit overflow or publication failures return a typed
  error/nonzero exit; never a successful complete assignment.
- A structurally valid diagnostic preview may include unknowns or failed obligations.
  It is labeled non-executable with explicit reasons. Source parse success cannot
  turn it into a ready work package.
- Assignment-ready requires every applicable structural obligation to pass and each
  relevant unknown to have a pinned owner/reviewer disposition with scope and rationale.
  Unavailable behavioral evidence may be deferred explicitly; it never becomes pass.
- A change to task, source, policy or evidence invalidates the prior readiness result.
  Local attribution is not authenticated authority. Readiness is not permission to
  modify protected repositories, merge to main or publish software.

## Proposed CLI surface

```text
sdptool AREA model create blueprint from release:0.1.0 to work:Calibration \
  --entry System.design --task assignment.yaml --output OUTPUT [--json]
```

Retain the existing natural-language model facade. Blueprint is a generated view,
not a fifth mutable/frozen lifecycle artifact. This spelling is a proposal; do not
advertise it in installed help until implementation. No automatic release/commit,
script execution, source migration, viewer launch or network lookup occurs.

The adapter captures both refs once; an editor change before consistent capture is
an error. Computation uses only captured bytes. Publication validates caller expected
identities and, for a live request, rechecks sources/task before publishing; observed
changes fail stale without replacing the previous result. Retained snapshot generation
may intentionally analyze old captures but must say so. Arbitrary editor atomicity is
not promised. Cancellation stops work before publish and never modifies source inputs.

## Bundle and publication

```text
index.md                 # intent, readiness, changes, boundaries and entry links
changes.md               # semantic differences with NOW/TARGET provenance
context.md               # annotated union graph and explicit cut/unknown inventory
obligations.md           # authored constraints and separate observed results
evidence.md              # applicability, not-run/stale/pass/fail references
sources/NOW/...          # complete captured source bytes; no .commits/.merge
sources/TARGET/...
assignment.yaml          # exact authored task bytes
blueprint.json           # typed analysis, identity, coverage, receipts
manifest.json            # output inventory/digests, publisher ownership
```

Generate every diagram from Analysis; do not hand-edit output. Escape labels and
links; no diagram/Markdown-driven command execution. Use relative source links and
check them before publication. Initial v1 renders Markdown/Mermaid; SVG is optional
through an explicitly selected existing renderer, not a new rendering engine.

Canonical payload identity hashes both source inventories, task, policy version,
producer version and evidence digests. Generated wall-clock timestamps and host
paths belong to execution receipts, not deterministic payload hashes. Manifest
excludes its own hash to avoid recursion. Evidence records bind the exact tested
candidate, never just the latest pathname or a claimed release number.

Reuse documents.Bundle.Publish only after verifying its ownership/containment/error
contract against blueprint needs. It already preserves unmanaged notes and rejects
changed managed outputs, but that does not prove full crash durability. Place new
metadata in blueprint.json and retain the existing publisher manifest schema where
possible; do not replace the documents.Manifest format incidentally. Validate limits
and output/source disjointness, link targets and path/symlink safety before writing.
Stage beside the destination, preserve last good output on observed failures and
never publish truncated analyses. Recovery guarantees must be tested on the actual
adapter; physical power-loss guarantees are outside the first implementation.

## Bounded trial and completion boundary

Trial-Assignment.md is an authored test work package. The coordinator acts as Worker
on temporary copies only; a fresh independent Reviewer evaluates actual candidates
against that assignment, including a deliberately injected violating control. The
trial establishes readability and structural boundary detection, not autonomous
agent drift prevention, real Ponsse implementation or production generator usability.
A product trial must repeat with a generated bundle after implementation.

The successor ImplementationPlan stages production Go analysis, deterministic
publication and assignment evidence integration. BP2-A's explicit owner disposition
of the proposed calibration extraction remains open; plan/document availability does
not fabricate that approval. Generic analyzer implementation need not authorize an
actual Ponsse refactor. KB050 stays active until its implemented outcome is verified.

## Blueprint catalogue — owner input, 2026-10-07

The owner requests an XFMD Blueprint sub-tab under SDP, with a tree grouped by
assignment/implementation progress. This is a new consumer-workflow requirement;
the producer design above did not yet define a persistent catalogue or its state
projection. No implementation or XFMD capability is claimed.

Owner selected the catalogue scope for KB050 / PLAN-SDP-0020 on 2026-10-07
(Session0008 T005). The following is the implementation direction; detailed
transition/schema design remains a required milestone activity:

- Retain reviewable generated bundles under SDP/Blueprints/<blueprint-id>/<revision>/.
  A blueprint ID identifies a scoped task; immutable published revisions pin exact
  source/task/policy/evidence identities. Scratch previews need not be retained.
  Do not confuse these revisions with ModelGovernance RELEASE identities.
- Discovery derives the catalogue from actual bundle metadata and linked canonical
  project-management records. No manually maintained navigation.json or central
  duplicate catalogue is required. Missing/malformed bundles or references remain
  visible as diagnostics, never guessed states. Discovery does not generate bundles.
- Work state belongs to an identified assignment or plan milestone in the existing
  project-management history, not a mutable status field in generated Markdown.
  A single KanBan card may span several blueprints; its state alone cannot establish
  each blueprint's assignment or implementation status. Multiple assignments/revisions
  must retain distinct identities.
- Suggested work states: draft, ready, assigned, in-progress, review, completed;
  on-hold, canceled and superseded are exceptional dispositions. Assigned requires
  an explicit assignee bound to a revision. Ready incorporates producer readiness
  plus required workflow disposition. Completed needs reviewed implementation evidence;
  source validity, a model release or a completed parent card is insufficient.
  Exact transition schema and authority remain to be designed.
- Show freshness (current/stale), structural validity, evidence status and preliminary
  input independently of work state. New source/task bytes make comparisons stale
  relative to live inputs; they do not rewrite or erase a completed historical result.
- Proposed discovery projection: a Blueprints root of kind tab, virtual work-state
  groups, task/blueprint nodes, revision/assignment children and typed document targets.
  Refresh recalculates it from sources and records. XFMD consumes this projection
  and opens generated Markdown/Mermaid. It does not infer implementation or own
  semantic analysis.

Before executing BPI2/BPI3, refine their plan/contract to include the selected
retention location, discovery schema, state-record owner and stale/missing-reference
tests. Native XFMD changes require their own XFMD card/agent. PLAN-SDP-0020 now includes BPI2-M2 catalogue/discovery and BPI3-M2 lifecycle
projection. Native XFMD implementation is not activated by this scope selection.


## BPI2-M1 concrete adapter contract — 2026-10-07

The first adapter finalizes task input as strict JSON, retained as assignment.json,
rather than the earlier proposed assignment.yaml. CLI spelling is
model create blueprint from REF to REF --entry FILE --task FILE --output DIRECTORY.
This is a diagnostic preview only. The implemented contract and source identity /
publication limits are in SDPTool/blueprints/README.md. The earlier catalogue
design remains BPI2-M2; no immutable retention or discovery root is inferred from
the preview publisher. Canonical task and analysis fields remain typed Go contracts.
