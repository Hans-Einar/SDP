# SDPTool installation design contract — draft 1

Status: design for implementation, not an implemented or published API.
Authority: REQ-SDPTOOL-007, owner-selected SDPTool/gh-sdp ownership and
[PLAN-SDP-0002](Plan.md). The [canonical SDL model](../../../03--Architecture/SDPTool.design)
provides identities, responsibilities and validated exchanges. This document defines
wire/file semantics that the present SDL type system cannot fully express.

## Ownership and alternatives

| Owner | Responsibility | SDL identity |
| --- | --- | --- |
| SDPTool CLI | Project root, orchestration, deterministic plan/apply/resume and truthful outcome | InstallationCoordinator |
| SDPTool core | Observe facts, preserve project content, validate adoption, calculate transitions | InstallationBaselineInspector, InstallationPlanner |
| SDPTool core | Resolve/verify process release and artifact input | InstallationReleaseResolver |
| SDPTool core | Execute reviewed steps, lock, journal, recover and publish installed facts/history | InstallationExecutor, InstallationJournal, InstallationRecorder |
| gh-sdp | Obtain a compatible verified SDPTool executable and forward argument vector, standard streams and exit status | GhSdpLauncher |
| External artifact service | Distribute versioned catalog/manifests/payloads and executable assets | ReleaseArtifactService |

Selected design: use one SDPTool executable and library implementation in root
SDPTool. gh-sdp caches and launches the executable. Linking the same Go library into
gh-sdp could avoid an extra binary fetch, but couples engine upgrades to client
releases; select executable delegation for the first delivery. A separate Go
Toolkit engine or client-side migration logic would duplicate ownership and is
rejected. No required daemon, shell evaluation or PowerShell runtime.

Internal model channels denote Go API calls, not separate services. Only gh-sdp ↔
SDPTool crosses a local process boundary; artifact retrieval crosses the network.
Direct sdptool commands enter the same coordinator, bypassing GhSdpLauncher.
ReleaseRepositoryProcess is the existing external publishing service, not a server
this project must implement. Databases in this model denote persistent storage,
including files/cache/journal; none implies SQL.

## Proposed command interface

Preserve existing SDPTool commands and their output. Add the following forms:

```text
sdptool [PROJECT] install [--release RELEASE | --artifact PATH] [--plan-output FILE] [--json]
sdptool [PROJECT] upgrade [--release RELEASE | --artifact PATH] [--manifest FILE] [--plan-output FILE] [--json]
sdptool [PROJECT] install --apply PLAN [--json]
sdptool [PROJECT] upgrade --apply PLAN [--json]
sdptool [PROJECT] upgrade --resume OPERATION [--json]
gh sdp upgrade --manifest xfmd-upgrade.yaml --plan-output /tmp/xfmd-plan.json
gh sdp upgrade --apply /tmp/xfmd-plan.json
```

Default install/upgrade is **preview only**. Apply requires a saved, reviewed plan;
there is no implicit confirmation prompt, automatic build or hidden mutation.
`--manifest` means adoption input, not release inventory or an executable script.
Its selected target must agree with an explicit --release/--artifact if supplied.
Without an adoption manifest, upgrade requires valid installed facts. Empty install
must reject an existing unversioned SDP area and recommend explicit adoption.

The optional project defaults to cwd. Accept an existing project root or its SDP
area, without climbing arbitrary parents. Canonicalize the physical root and reject
symlink traversal. Installation selection must not require navigation.json: a clean
install cannot use the current navigation-only Discover function as its validity
test. Reuse path safety while retaining separate operation eligibility.

One positional project only. Flags are mutually exclusive across preview/apply/
resume; conflicting inputs fail before writes. Paths inside adoption documents
are relative to that document where explicitly marked; migration paths are always
project-relative. No shell command interpolation. Resolve a saved plan path before
changing cwd; the plan is bound to the physical project root and is not portable
between copies. Each disposable/live target needs its own fresh plan.

`--json` stdout is exactly one UTF-8 object plus LF, with schemaVersion,
operation, status, projectRoot and the operation-specific result or error code.
Diagnostics/progress go to stderr. Human mode summarizes conflicts, changes and
preservation and identifies the plan output; it does not print payload bytes.
Preview stdout alone may serialize a plan; --plan-output writes it outside the
managed project inventory using create-new semantics. Existing unequal outputs
require a new path, not silent replacement.

Exit codes: 0 successful preview/apply/resume or no-change; 2 invalid invocation,
unsupported input/schema/protocol; 3 plan conflict or drift; 4 verification/retrieval
failure; 5 lock contention or pending operation; 6 interrupted/failed mutation with
recoverable operation ID. A non-applicable preview returns 3 and a structured
canApply=false result. A completed operation with warnings still returns 0 and
lists them explicitly. Interrupt signals request stop at a journal boundary;
uncatchable exit is recovered from the journal, never mislabeled success.

## Three identities and their records

| Record | Authority / storage | Required meaning |
| --- | --- | --- |
| Release descriptor/inventory | Publisher-controlled signed asset; cache outside project | Immutable process release, source commit, inventory and payload digest, file types/ownership/hashes, explicit transitions and engine capabilities |
| Installed receipt | Project-local lookup facts at the existing installed-manifest location | Observed installed release and descriptor digest, process/profile identity and completing operation; not proof of authenticity by itself |
| Adoption manifest | Reviewed local YAML/JSON supplied with --manifest | Unknown/manual baseline, expected observations and explicit permitted mappings to a chosen release; cannot redefine the publisher inventory |

Keep three independent version axes: gh-sdp client, SDPTool engine/protocol and
installed process release/profile. Never infer a released version from folder
names, current Git branch or today's Toolkit 0.2.0 unreleased declaration.

Candidate identifiers are draft contracts: sdp-install-command/1,
sdp-release-descriptor/1, sdp-adoption/1, sdp-install-plan/1 and sdp-install-journal/1.
These do not change deployed schemas now. Installed facts need an explicit new
schema revision when adding release provenance; do not add unknown keys to the
closed installed-facts 2.0 reader. The implementation must upgrade that reader and
its tests before writing the new receipt, while preserving reading of 1.0/2.0.

An adoption manifest contains schemaVersion, project selector, baseline kind,
observed Git commit (provenance only), inventory entries and target release identity.
Each inventory entry has project-relative path, type, expected SHA-256 for files
or expected absence, and ownership claim. Mappings list source/destination with
preconditions; no glob commands, scripts, executable hooks, network URLs or escape
paths. Require the exact snapshot for the affected tree and incoming-link scope;
additional project-owned files are observed/preserved, not silently excluded from
plan identity. Wrong hashes, missing expected paths, overlaps, duplicate/case-folded
paths, traversal, unsupported objects or ambiguous ownership yield conflicts.
The manifest author cannot reclassify arbitrary project files as removable managed
files: classification must match verified old inventory or an explicit adoption
mapping reviewed as project-content preservation/relocation, with byte backups.

Release entries carry normalized relative path, file/directory type, ownership
(managed or initialize-if-missing), payload digest and optional declared retired
managed paths. Explicit migration records define rename destinations and supported
baseline predicates. An absent target entry is not permission to delete a project
file. Inventory covers root AGENTS/skills as well as SDP/. No implicit parent repo,
subrepo rewrite or execution of project-local code.

## Distribution, trust and offline behavior

Design choice: publisher signs the exact release-descriptor bytes with an Ed25519
key; supported public keys/key IDs ship in the wrapper/standalone distribution.
The descriptor pins platform binary, profile payload and inventory hashes and
protocol/capability requirements. Reject unknown signing keys, bad signatures,
wrong digest/size/platform and unsupported protocol before executing a downloaded
binary. A descriptor cannot supply its own new trusted key. Key rotation requires
a client update carrying an overlap of trusted keys; publishing actual keys and
releases belongs to the later release work, not this design.

Release selection defaults to eligible stable process releases from the configured
canonical publisher. Engine bootstrap selects a compatible platform binary from
a verified descriptor. Cache under the user's platform cache directory by digest,
using exclusive creation and atomic promotion after verification; private files,
no executable content from the target project. Child invocation uses a fixed
argument vector. It must advertise sdp-install-command/1 before delegation; do not
infer support from a version string. The existing --version response is extended
with an optional installationProtocol capability only when implemented.

SDPTool itself verifies process input; the wrapper's checks do not replace that
boundary. Apply/resume require all exact blobs locally and must not select a newer
release or download mutable input mid-operation. An explicit offline mode uses
only matching verified cache entries and fails clearly if missing. A hash is an
integrity check, not publisher authentication. Cached bytes are checked again
before use; a local installed receipt never authorizes a replacement release.

A local development artifact may be selected only by explicit --artifact together
with --allow-unreleased explicitly set. It is recorded as unverified
local development provenance and cannot claim signed release support. The current
PowerShell-specific profile prerequisite must not be silently ignored: a future
engine-neutral profile needs an explicit capability mapping/new release identity.
Use existing artifacts as legacy fixtures, not as proof a Go release already exists.

## Plan, ownership and deterministic execution

The plan contains schema/version, operation, canonical root, verified descriptor
and payload digests, observed old receipt, adoption digest or explicit absence,
complete observed path/hash/type snapshot, policies, ordered actions, preserved
paths, conflicts and warnings. Stable identifiers refer to exact input bytes.
It records no timestamp in its identity. Define planDigest as SHA-256 of the exact
canonical JSON bytes excluding planDigest: UTF-8/LF, recursively sorted object keys,
no insignificant whitespace, integer-only numeric fields, deterministic arrays.
Duplicate JSON/YAML keys, aliases/tags and non-finite values are rejected.
Initial limits: 16 MiB descriptor/adoption inputs, 64 MiB serialized plan/journal,
100,000 inventory entries, 64 MiB per payload file and 512 MiB total payload.
Reject over-limit inputs before allocating/decoding full content; validate sizes
again after decoding. Changes to limits require an explicit capability revision. Explicit
null denotes absent files; zero hashes are not an absence marker.

Policy: writes sorted by normalized destination ordinal byte order, then source
removals sorted the same way. Every relocation writes and verifies destination
before removing source. Reject dependency cycles/collisions that cannot satisfy
this order instead of improvising destructive moves. Backups precede each mutation.
Defer installed receipt and final Maintenance/history publication to journaled
finalization after ordinary actions verify. Pin the ordering policy in the plan;
PowerShell fixture comparisons normalize only documented representation differences.

Read-only planning observes inputs twice and rejects changed snapshots. It makes no
project directory/lock/journal. Apply obtains the project lock and recomputes the
entire plan with the same policy/input bytes; require exact identity. The reviewed
plan is not an authorization token for another root. Before each step, verify the
expected old bytes/type again; concurrent changes stop progress. Catchable drift
before any mutation yields conflict without an installation journal outcome claim.

Use old inventory, actual bytes and target inventory for managed edits. Preserve
project-owned documents; initialize missing templates only. Unmodified managed
files can be replaced with backups. Modified managed files are conflicts unless
an explicit reviewed refresh policy names them; no unconditional global overwrite.
Preserve unknown root instructions in the established AGENTS-project convention,
reject unequal preservation collisions, and verify repeat/no-change behavior.

Check confinement/case collisions before writes; reject symlinks, reparse points
and special files in the affected tree. Record scan exclusions and non-Markdown
references as warnings requiring explicit disposition for relocation. Reuse the
existing link-rebasing behavior without interpreting arbitrary source code.

## Journal, recovery and finalization

Reuse the existing .sdp-operations location with a versioned journal. Store exact
plan/input digests, before/after bytes/hashes, backups, next/completed step state,
operation time and preallocated Maintenance/event identities. Acquire one exclusive
project lock; a pending operation blocks other plans. A legacy v2 pending journal
is recognized and blocks Go mutation; resume it with the retained matching legacy
engine. Do not reinterpret old journals as the new protocol.

Per step: durable backup, atomically replace a file, verify output, publish journal
checkpoint. Resume accepts only the recorded before/after bytes and verifies every
completed step; post-failure edits cause a conflict. Reuse reserved finalization
bytes/IDs, never allocate another completion event on retry. Forward recovery is
the guarantee; no whole-tree rollback, filesystem power-loss guarantee or automatic
backup pruning is promised. Keep a failure report/operation ID even before the
Maintenance report can be safely published.

Finalization publishes the new installed receipt, Maintenance report and append-only
project-management events as journaled writes, then marks the journal completed.
Discovery reports incomplete while a journal is active/failed, even if receipt
publication already occurred. A repeat against an unchanged release is no-change:
no duplicate report/events, no reset timestamps and no byte churn. Unknown old
versions remain unknown in migration provenance. Existing ledger prefix and IDs
remain byte-identical; no release facts or work records from this repository are
seeded into consuming projects.

## SDL correspondence and limits

Planning calls coordinate baseline inspection, release resolution and plan creation.
Apply calls coordinate execution, checkpoints and publication. Resume uses the
same journal/recorder ownership. Bootstrap calls only retrieve the executable.
Model scenarios illustrate these paths and one pre-mutation drift rejection.

Model payload contracts are logical envelopes. The Inst...RecordBytes fields carry
opaque serialized content; SDL does not validate the JSON schemas, signature rules,
path confinement, hash comparisons or idempotence. Reused message payload shapes
are not assertions that each internal request has identical Go structs. Contract-
scoped field identities satisfy SDL's single-owner rule. AdoptionDigest has an
explicit textual `none` sentinel for normal install/upgrade; the wire schema uses
null and the adapter maps it. New activities remain planned even when the design
parses. The domain has no fixed-bit datagram encoding, so VP10 is intentionally
not used. Database storage is file-based and the generated data view is logical.

See [scenario acceptance](Scenarios.md) for the design walkthrough and planned
implementation evidence. Contract choices above are the working design produced
under owner authorization; independent review and owner acceptance are separate.
