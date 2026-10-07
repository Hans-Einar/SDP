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
