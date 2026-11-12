# A0 executable architecture contract

Selected under the owner's Activity0002 execution authority. The existing Activity
Design/Commands documents own user outcomes. Activity0001's dependency boundaries
remain mandatory. No independent lifecycle implementation may appear in CLI, MCP
or XFMD. Code and tests supplement the SDL model's structural proof limits.

## Typed service and records

`SDPTool/activities` owns Activity requests/results, replay, record construction,
transition rules, recovery and domain projections. It imports projecthistory and
standard/parser utilities, never root SDPTool, governance or host adapters. Root
composition resolves project identity and wires an execution-inspection port.
One operation definition supplies CLI/action dispatch metadata and strict decoding.
Create and Session-create do not dispatch work. Callers supply trusted actor and
authority separately from JSON; owner dispositions cannot be self-asserted by
request fields. Existing read-only and capability-bound MCP contracts remain intact.

Read the manual descriptor/event 0.1 format unchanged. New operations may use
version 0.2 to add scope, preserved descriptor references, linking and atomic switch
metadata. Keep the existing 0.1 schemas as historical contracts. Each event pins
the exact relevant descriptor bytes. Before replacing a live descriptor, retain its
previous bytes under `Activities/.descriptors/<sha256>.yaml`; old path/hash pairs
resolve against that verified immutable snapshot when the live path has changed.
New descriptor snapshots use the same store. Never substitute current bytes for a
historical hash or infer status from the folder/descriptor.

An operation has a stable UUID, canonical request digest, expected history revision
and, when applicable, expected Activity predecessor/descriptor revision. Validate
portable relative paths, existing typed references, identity uniqueness and target
state. Source cards can be absent for a direct owner need in new records; 0.1 records
retain their actual schema. Session IDs are allocated across legacy and Activity
homes. New Session documents may precede a plan but are not execution bindings.

## Persistence and recovery

Use the shared canonical history lock/reservation. `Transaction.Read(owner)` and
`AppendBatch(owner, originalRevision, events)` operate under an already-held lock;
ordinary Append is never called recursively while holding that transaction.
The shared primitive preserves old bytes, checks the entire ordered batch and
publishes one atomic replacement. Domain validation remains in activities.

Activity mutation retains a bounded intent before effects: exact request digest,
allocated identities, fixed event bytes/time, file preimages/postimages and result.
Reserve the canonical writer before publishing files. Recovery reacquires the
same owner, verifies each current file against its before/after image, completes
remaining effects, publishes the fixed event batch and releases the reservation.
Refuse external conflicts instead of overwriting them. A completed receipt supports
exact retry after later events; different content with the same operation ID fails.
Intent/receipts are operation recovery metadata, not another authoritative ledger.
No cancellation after a commit may misreport that the committed effect was undone.

Switch holds the former active Activity and activates the selected one in one
batch with complete-batch validation. Only the final validated replay is exposed;
incomplete/forged batches are diagnostics. Never concatenate conflicting clone
histories or claim local locks coordinate disconnected canonical writers.

## Execution boundary and adoption

Read-only browsing of valid manual records is supported before runtime admission
adoption. Explicit adoption enables the Activity admission policy and verifies
scope/membership, existing work and the selected canonical writer. Existing projects
without adoption retain current governance behavior and Sessions navigation when
they have no Activity catalog. New create/adopt commands must explain which policy
is enabled; merely viewing documents cannot change execution policy.

Governance owns actual runs/attempts and their recovery. An injected inspection
port supplies linked execution and uncertainty; no inspection must never mean no
work. Active-exit/complete/switch operations require known execution disposition.
When work is still active, retain a pending lifecycle request that blocks new
admissions, permit observation/interrupt/recovery, and finish only after actual
quiescence. Do not hold a filesystem lock or canonical record reservation for the
lifetime of a running agent. Detailed admission serialization must be reconciled
with existing project-coordination → canonical-history → operational-state ordering
before A5 writes; never introduce the reverse lock order.

After adoption, prepare/start/activate and actual launch/retry check active Activity,
explicit work membership and scope, not merely the CLI command name. Bind admitted
work to a checked Activity revision and revalidate before dispatch. Existing Session
execution bindings remain under their current store; catalog-validated Activity-local
Session documents become valid references. Scope additions/owner completion require
their actual trusted authority. Unrelated read-only discussion and KB capture remain
possible without opening an implementation bypass.

## Consumer and evidence boundaries

Producer-owned navigation preserves fixed state and phase order and one Sessions
home per Activity. Ordinary source/reference documents are discovered, with readable
landing/phase/history/timeline content generated on demand. Old clients retain the
existing file/open protocol; Activity-specific lazy actions need explicit capability
advertisement and native consumer qualification. No silent truncation or guessed
status is allowed for invalid, unavailable or large catalogs.

Timeline intervals follow causal lifecycle events and actual timestamps. Unknown
time/clock reversal stays visible. Commits are associated only by explicit Activity,
linked work or source-event evidence, never proximity in time. Dense views group
markers without losing drill-down/source records; calendar activity is not effort.

Each implementation handoff includes this contract, the plan and pinned Blueprint.
Workers explain ownership/deletions before editing; return contradictions here.
Tests must cover real byte-preserving manual compatibility, stale/duplicate inputs,
interrupted effects, concurrent activation, descriptor history, pending work and
native selection. Final review must cover all A1–A6 outcomes, not just creation.
