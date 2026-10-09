# Blueprint assignment lifecycle — BPI3-M2

Implementation refinement, 2026-10-09, PLAN-SDP-0020 / KB050 / Session0008.
The owner-selected outcome remains assignment progress and evidence-backed discovery.
This document separates the persistence prerequisite M2a from integrated M2b.
It is a design contract, not a claim that lifecycle commands are implemented.

## Existing authority and persistence

The sole project-management history is SDP/ProjectManagement/Ledger.ndjson.
Traceability remains the authority for reviewed implementation evidence. The
parallel ProjectGovernance worktree at 4044f37 implements a private local run store
outside the project, with Linux locking and atomic snapshots. Its ledgerCoverage
reader inspects selected card/plan subjects; it does not write canonical history.
Its operational HEAD/generations must not become blueprint assignment authority.

No compatible canonical writer exists in the current implementation branch.
M2a adds a reusable SDPTool/projecthistory append primitive, preserving all prior
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
