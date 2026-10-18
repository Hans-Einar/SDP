# Blueprint document generation — sdp-blueprint/1

The first producer generates a diagnostic blueprint from two ModelGovernance
artifacts. It runs locally without Git, an installed SDP template, a viewer,
Python or automatic code execution.

    sdptool /path/to/models model create blueprint from release:0.1.0 to work:Calibration --entry System.design --task task.json --output /path/to/preview

Use --json for the result envelope; default output uses the common human presenter.
This command does not commit/freeze WORK or assign work to an agent.

Task input is a strict JSON object using the SDL blueprint Task fields, for example:

    {"ID":"calibration","Intent":"Extract calibration ownership","Context":[],"Exclude":[],"Rules":[],"AllowedChanges":[],"AllowedModelPaths":[],"AllowedCodePaths":[],"Protect":[]}

Empty permissions make changed facts fail structural constraints. Such results are
useful diagnostic previews, never ready assignments. See SDL/go/blueprint/README.md
for permissions, fact expectations and boundary/subtree protections. The actual
task bytes are retained as assignment.json. YAML is not accepted by this increment.

Output includes index.md, changes.md, context.md, context.mmd, changes.mmd, obligations.md,
evidence.md, blueprint.json, manifest.json and sources/NOW plus sources/TARGET.
Every source file captured by ModelGovernance is retained; parsing visits the
selected entry's reachable SDL graph. SDUI bytes remain visible with unsupported
semantic coverage. Generated source links are bundle-relative and checked.

The bundle revision binds both artifact identities/heads/metadata digests, complete
source digests, task bytes, entrypoint, compiler profile/sourcegraph version, executable SHA-256,
Go runtime version, analyzer identity, policy and producer protocol version.
It differs from the reachable SDL graph revisions, which are retained separately.
No timestamps or machine-local absolute input paths enter deterministic payloads.
The producer protocol version must change when generated-payload semantics change.

Inputs are recaptured and compared immediately before publication; observed edits
reject publication without replacing prior output. This is optimistic freshness,
not atomic exclusion of arbitrary external editors. Output must not overlap source
artifacts, task input or model operation storage. Symlink ancestry is rejected.

The existing document publisher stages output, rejects edits to generated files,
and preserves unmanaged notes. Rename/write failure tests verify the prior output
remains available; if rollback itself fails, the error identifies the retained
backup. There is no power-loss durability claim. Preview directories may be
regenerated; immutable catalogue mode is described below.

The producer does not persist assignment states or claim code conformance.
Those remain BPI3.

## Immutable retention and browsing — BPI2-M2

Replace --output with --catalogue /project/SDP/Blueprints to retain a revision.
No registration or manual copying is needed. The result path identifies
<catalogue>/<hash-of-System-and-task-ID>/<retainedRevision>/.
Folder names are safe full digests; human labels remain in metadata.
revision identifies blueprint inputs/analysis; retainedRevision binds every exact
retained file, including manifest bytes, using the ModelGovernance Digest framing.
Both identities are returned. Equal input/tool/task bytes are idempotent; new bytes
add a sibling. Logical immutability does not prevent manual edits; discovery
diagnoses them, including edits accompanied by a rewritten manifest.

Discovery verifies both captured source digests, metadata identities, complete
file inventory/digests, retainedRevision and links. It returns a Blueprints tab
with task/revision nodes and Markdown/Mermaid open targets. Retained source copies
are not added to the current SDL/SDUI inventory. No index or ledger is written.
Each revision has workState unknown until BPI3; integrity does not establish
assignment, implementation or live-source freshness. Preliminary is separate.

Malformed, missing, modified, unsupported or wrongly placed revisions stay visible
without a trusted open target. Duplicate copies outside their canonical path are
invalid. No status-folder moves occur. Scanning is bounded to 256 task/revision
candidates, 10,000 verified files/directory entries and 64 MiB of verified bytes.
Metadata has an 8 MiB per-file limit. Symlinks are rejected.

## Scoped readiness assessment — BPI3-M1

    sdptool model assess blueprint --bundle /project/SDP/Blueprints/KEY/REVISION --evidence /trial/evidence.json --json

The command verifies a retained bundle, reads strict sdp-blueprint-evidence/1 JSON,
and emits sdp-blueprint-assessment/1. Human output is the default. Exit 0 means
ready within the declared scope; exit 3 returns a valid blocked assessment;
invalid input/integrity errors use the normal nonzero diagnostic path. No command
in the evidence is executed. Nothing is written to the catalogue or a workflow
ledger. Discovery workState remains unknown until BPI3-M2.

The input fields are defined by Evidence and its nested types in assessment.go:

- schema, blueprintRevision, retainedRevision and taskDigest pin the exact bundle;
  scope names the bounded assessment. Struct keys require exact JSON spelling.
- code maps NOW/TARGET to revision, root, digest and files. Each files entry maps
  a relative code path to SHA-256; digest uses model.Digest on those captured bytes.
  root resolves beneath the evidence JSON directory. This explicitly scoped file
  inventory is not automatic proof of repository-wide completeness.
- mappings locate a selected element by side, role, path, symbol and fileDigest.
  Missing, stale or ambiguous mappings remain unknown. Symbols are opaque locators;
  the assessor does not interpret Go or authenticate inline tags.
- requiredChecks lists id, side and obligation (an actual structural obligation ID
  or an authored behavior:NAME reference). receipts bind those IDs and sides to
  sourceDigest, taskDigest, codeDigest, command, environment, inputsDigest, inputs,
  result, artifact, artifactDigest and traceability. inputs maps fixture/config
  paths to SHA-256; inputsDigest uses model.Digest with those same paths. Observed
  pass/fail requires real pinned input files and an evidence artifact. Commands
  and Traceability references are attributed metadata, never executed or remotely
  resolved. The assessor checks file hashes, not whether a log tells the truth.
- dispositions names an unknown ID from an earlier blocked assessment, actor,
  role (owner or reviewer), the exact assessment scope, rationale and traceability.
  Unknown IDs are stable within that blueprint. Unknown dispositions permit scoped
  deferral; they do not change missing/not-run/stale/unavailable into pass. Reported
  failures and structural constraint failures block readiness despite dispositions.

File access rejects symlinks and escape, shares a 64 MiB/10,000-read bound, and
limits evidence JSON to 1 MiB. Unsafe/oversized input fails closed. Missing code or
artifacts produce explicit unknown/stale results. Duplicated or conflicting record
identities are errors. Assessment identity binds exact evidence input and observed
file hashes; changed code, fixture or log bytes invalidate the old assessment.

A ready result is not authenticated owner approval, complete code conformance,
current-WORK freshness or permission to execute/merge/release. It assesses captured
historical models and scoped supplied evidence. General code mapping and executable
channel tests remain separate work. Keep assessment receipts/results with existing
Traceability evidence; lifecycle operations will reference their pinned identities.

Each relative input is captured once and cached; the evidence JSON and observed
files are rechecked before returning. This detects observed concurrent edits; it
is optimistic freshness, not atomic exclusion of arbitrary editors. The assessment
schema versions readiness policy: bump it when readiness semantics change.


## Revision-bound assignments

The local command `sdptool SDP model assignment apply --request REQUEST.json
--as ACTOR --authority controller|assignee|reviewer` accepts one checked request;
`model assignment list` returns a read-only projection. Add --json for consumers.
Apply uses the canonical ProjectManagement ledger, never a mutable blueprint file.
A caller keeps eventId/occurredAt/request bytes stable when retrying and uses the
last assignment event as expectedEvent for a new operation. Different assignments
may share one retained revision. The CLI is trusted local attribution, not remote
authentication; an adapter must bind its Principal independently of request data.

The lifecycle contract and schema are linked from
[Assignment-Lifecycle.md](../../SDP/04--Design/SDPTool/Blueprints/Assignment-Lifecycle.md).
The pure reducer is blueprintstate; lifecycle.go owns bundle/source/evidence checks.
Submission requires scoped TARGET code and passing check receipts tied to a pinned
Traceability record. Acceptance requires a different attributed actor from the
assignee; completion needs controller disposition and unchanged accepted evidence.
Neither a parser pass nor a model release completes an assignment.

Discovery adds assignment-state groups while retaining immutable task/revision
navigation. Historical work state survives missing or stale inputs, with separate
validation, freshness, readiness and evidence diagnostics. Catalogue and lifecycle
read budgets are separate and documented; exhausted live checks remain unknown.
No checks execute from receipts and no native XFMD widgets are implemented here.
