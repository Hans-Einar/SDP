# Blueprint assignment lifecycle — BPI3-M2

Implementation refinement, 2026-10-09, PLAN-SDP-0020 / KB050 / Session0008.
The owner-selected outcome remains assignment progress and evidence-backed discovery.
This document separates the persistence prerequisite M2a from integrated M2b.
This design governs the M2b command and replay implementation; delivery evidence
is separate from owner acceptance and release.

## Existing authority and persistence

The sole project-management history is SDP/ProjectManagement/Ledger.ndjson.
Traceability remains the authority for reviewed implementation evidence. The
parallel ProjectGovernance worktree at 4044f37 implements a private local run store
outside the project, with Linux locking and atomic snapshots. Its ledgerCoverage
reader inspects selected card/plan subjects; it does not write canonical history.
Its operational HEAD/generations must not become blueprint assignment authority.

At M2 planning, no compatible canonical writer existed. Delivered M2a adds a reusable SDPTool/projecthistory append primitive, preserving all prior
bytes. This follows existing lock/stage/sync/rename patterns, but does not copy or
integrate the separately owned governance subsystem. M2b will use this primitive
for versioned x-blueprint assignment events in the canonical ledger. It must extend
all applicable history validators before emitting such events. Existing consumers
that filter selected card/plan subjects need no inferred blueprint state.

## Canonical append contract — M2a

Read returns exact ledger bytes and their SHA-256 revision. Append takes a pinned
expected revision and one complete generic-envelope event with a unique event ID.
The owning domain validates payload/transition semantics; the shared writer checks
JSON/envelope integrity, event identity, stream bounds and byte-preserving append.
It does not allocate domain counters or confer actor authority.

A cooperating writer holds a nonblocking OS lock on a stable sidecar. It reads and
validates the current stream, rejects a stale revision, stages the complete old
bytes plus the new event and newline in the same directory, fsyncs, rechecks the
original bytes and atomically renames the stage. Existing history bytes are never
normalized. This is append-only history with atomic file replacement, not a second
ledger. Readers see the old or new complete stream. Caller-supplied event time,
ID and payload remain fixed across retry.

An exact repeat of an existing event is idempotent even after later events; reuse
of its ID with different bytes fails. Failure before rename leaves the old stream;
failure after rename may be uncertain to the caller and is resolved by repeating
the same request. A torn/malformed existing stream requires explicit reconciliation;
the library never truncates, repairs or guesses records. Orphaned staging files
have no authority. Linux supports mutation; other platforms explicitly refuse
writes until their locking/durability behavior is implemented and verified. Read
is portable. The bound is 16 MiB per ledger/event capture.

All programmatic canonical writers must adopt this lock/CAS protocol before being
claimed safe together. The existing manual scripts and arbitrary editors do not
participate; the final recheck detects observed interference but cannot eliminate
the final same-user race. Do not claim protection against uncooperative writers,
physical power loss on every filesystem, or authenticated event authorship.

## Assignment identity and authority — M2b

An assignment is distinct from a blueprint task and a retained revision. Its event
chain pins assignment ID, task/System, blueprint and retained revisions, assignee,
plan/milestone, readiness assessment identity and selected live source references.
Multiple assignments may refer to the same revision. Current state is rebuilt from
canonical events; no manually maintained registry or mutable bundle status exists.

A trusted local controller creates/adopts and assigns work. A bound assignee may
start and submit only its own assignment. An independent reviewer may accept or
reject submitted evidence; the assignee cannot approve its own result. Controller
closure requires that accepted, exact-revision evidence and a recorded disposition.
Caller capabilities come from trusted adapters, not fields claiming a role in a
request. A CLI is a trusted local invocation boundary, not authenticated remote
identity; future MCP callers require a separately bound adapter.

| Action | Origin | Result | Gate |
| --- | --- | --- | --- |
| create | absent | draft | Valid retained revision and scoped plan identity |
| adopt-readiness | draft | ready | Recomputed ready assessment and current pinned sources |
| assign | ready | assigned | Controller names assignee; freshness rechecked |
| start | assigned | in-progress | Bound assignee; expected assignment revision |
| submit | in-progress | review | Bound implementation and verification evidence |
| reject | review | in-progress | Independent reviewer rationale |
| accept-review | review | review | Independent reviewer pins accepted evidence |
| complete | review | completed | Controller closure referencing accepted exact evidence |
| hold | nonterminal | on-hold | Preserve previous state and rationale |
| resume | on-hold | prior state | Required freshness/authority gates rechecked |
| cancel | nonterminal | canceled | Controller disposition |
| supersede | nonterminal | superseded | Controller and explicit successor reference |

Terminal states never reopen or silently transfer their approval to another revision.
A new attempt receives a new identity. Operation identity and expected revision
prevent duplicate transitions and stale concurrent updates. Missing/ambiguous links
stay unknown; no parent card state implies individual assignment completion.

## Evidence and discovery — M2b

BPI3-M1 ready means a bounded attributed assessment, not workflow permission.
Readiness adoption must re-evaluate inputs and source freshness; a changed assessment
cannot reuse an old disposition silently. Completion additionally needs applicable
Traceability evidence, reviewed implementation/code hashes, actual check results
and independent closure disposition. An SDL parser pass/model RELEASE is insufficient.
No general code-conformance guarantee is inferred from one receipt.

Discovery rebuilds draft/ready/assigned/in-progress/review/completed and exceptional
state groups from validated events and actual retained bundles. It exposes each
assignment separately and keeps work state, historical validation, readiness and
live-source freshness distinct. Missing/damaged bundles or evidence remain visible
with diagnostics, without fabricated open targets. A completed historical revision
stays completed when new source appears; freshness can change independently.

M2b acceptance remains PLAN-SDP-0020's complete non-Git lifecycle trial, including
two retained revisions, concurrent/stale requests, expected negative controls,
consumer refresh and unchanged retained bytes. Native XFMD integration and release
remain separately owned. M2a alone does not deliver any assignment command or group.


## M2b transport and replay contract

Use `sdptool SDP model assignment apply --request REQUEST.json --as ACTOR
--authority controller|assignee|reviewer`, or `model assignment list`. The explicit
area is the project SDP directory, not the model-artifact area used by snapshot
commands. The local CLI trusts its invoker; flags attribute a local action and
are not authentication. An MCP adapter must bind a Principal independently;
request JSON cannot supply a principal or role.

The request schema is sdp-blueprint-assignment/1. Required fields are schema,
eventId, occurredAt, assignmentId, action, expectedEvent and reason. A create
request includes binding (bundle, task, system, revision, retained, plan, milestone,
modelArea, from, to). All paths are relative to the SDP area without traversal;
modelArea alone may be `.`. Bundle identifiers are checked against retained bytes.
Assignment IDs start BPA-; operation IDs follow generic EVT- rules and must be
kept stable across retry. expectedEvent is empty only for creation. Timestamps
use RFC3339 with at most nine fractional digits; replay compares nanoseconds.

Only adopt-readiness and submit take evidence; only submit takes traceEvent; only
assign takes assignee; only supersede takes successor. Unknown/case-aliased,
duplicate, missing required and null fields are rejected. JSON Schema and pure Go
replay constrain the canonical payload; the Python process validator matches the
transition contract. Installer validation reuses the Go reducer. Existing generic
management readers filter their event namespaces and never infer assignment state.

An active canonical Plan is required for create/adopt/assign. The milestone is an
explicit scoped identifier, not a parsed Markdown task. Later cleanup remains
possible after plan closure. hold/cancel/supersede/reject can operate when source
or evidence is damaged; resume rechecks the relevant readiness/submission gates.
Successors must already exist and belong to the same System.

Submission requires an actual scoped TARGET code capture and at least one passing
TARGET receipt. A stale code capture cannot be waived by unknown dispositions.
Each TARGET receipt references the selected Traceability event. Its generic
envelope uses x-verification:recorded and the assignment subject; its payload
schema is sdp-blueprint-implementation/1 with assignmentId, retained, evidenceDigest,
codeDigest and the exact set of passing TARGET check IDs. Submission, independent
accept-review and controller complete recompute these pinned inputs. This tool
validates attributed evidence; it does not execute commands or claim complete
model-to-code proof.

A persisted x-blueprint event records the exact request, adapter authority and
computed proof pins. Replay enforces per-assignment predecessors, transitions and
independent attributed reviewer identity. Retrying an identical operation returns
current assignment state without repeating the transition, even after inputs change.
Different request/principal content under that ID fails. Source/evidence captures
are rechecked immediately before canonical CAS append; external arbitrary writers
are not covered by a cross-file transaction.

## Discovery resource and status policy

The Blueprints tab keeps its task/revision tree and adds nonempty assignment-state
groups. Multiple assignments remain distinct. Each assignment exposes its last
event, assignee, immutable revision, workState, sourceFreshness, readinessStatus
and evidenceStatus. Missing bundles preserve historical state/adoption/submission
and suppress open targets. Valid targets carry the verified index.md content hash.
A new retained revision starts unassigned and never inherits an older completion.

Catalogue verification retains its 64 MiB / 10,000-file scan budget. Assignment
projection is a separate bounded pass: at most 256 assignments, with one shared
64 MiB / 10,000-file budget for bundles, assessment inputs and Traceability reads.
Read exhaustion produces unknown/unavailable diagnostics, not fabricated validity.
Live ModelGovernance comparisons are separately limited to two distinct assignment
bindings per refresh, with cached results for repeated bindings. Each comparison
uses two model snapshots, each using the existing 128 MiB source/two-pass and
16 MiB metadata bounds. Further live freshness stays unknown with a budget
diagnostic. This is a conservative first projection budget, not the old catalogue
budget applied to unlimited live-source reads. Explicit transitions still validate
the selected assignment fully.
