# Project Leader and SAD orchestration — KB052

For a readable explanation of current behavior, planned collaboration and owner
feedback topics, start with the [owner walkthrough](Owner-Walkthrough.md).

Owner direction: Session0007 T014, 2026-10-10, incorporating the
[revision-bound handoff](Project-Leader-Handoff.md). This refines the earlier
[PGD1 design](Design.md); it does not claim its new operations already exist.
[PLAN-SDP-0018](../../../05--Implementation/SDPTool/ProjectGovernance/Plan.md)
retains the implementation branch and milestone policy.

## Outcome and current gap

The owner normally talks to Project Leader, can also address SAD/Master directly,
and sees consequential feedback reflected once in shared work. Project Leader
commissions bounded work and stays available while workers execute. SAD prepares
ModelGovernance WORK/TARGET and generated blueprint evidence; Master owns one
implementation Session and its phases. Steering provides strategic, revision-bound
advice. Role guidance does not create runtime privileges or owner acceptance.

The owner reports having tried the lab and needing more familiarity. This is use
feedback, not rejection or acceptance. The old lab demonstrates one Maintenance
Master and cannot demonstrate the new coordination model. Keep its actual state and
candidate intact. Do not require further old-lab acceptance before selected KB052 work.

Current runtime has one checkout-local run store, run/attempt/host identities,
controller-only launch, restricted worker MCP, replay receipts, exact skill inputs
and correlated recovery. Both controller leases are checkout-wide. Master context,
coverage and routine selection are Maintenance-only; the role validator accepts
only Master. Simply changing a prompt or launching multiple current controllers
would not implement the proposed system. No native child mapping is proven.

## Selected architecture and alternatives

Reuse the Go governance services and existing app-server controller. Extend their
operational state and adapters in bounded increments; preserve existing Maintenance
clients and their retry/authority contracts. Extend the dedicated terminal client
rather than selecting an unverified stock-TUI plugin or maintaining a Codex fork.
The first view is a standalone read-only command; a later event-loop client owns
conversations and panels. No XFMD changes or general KB040 KanBan TUI are selected.

Use logical assignment identity independent of run/attempt/thread. Explicit
controller-issued capabilities distinguish Project Leader operations from worker
reports. A payload role string grants nothing. Future dispatch uses expected
assignment/bundle revisions and stable operation identity; uncertain sends must be
reconciled against the original attempt before any new launch. Existing worker MCP
keeps its five restricted operations until a separately tested change is selected.

The controller owns durable inbox delivery and per-thread turn serialization.
Project Leader uses short turns over that state; it is not the event loop. Prioritize
owner input, coalesce reports by stable identities, retain them across restarts,
and revalidate delayed reports against newer steering before dispatch. Host turn
completion, reported phase completion, reviewed evidence and owner acceptance
are distinct. Delivery exactly once means idempotent processing, not a promise
that transport never duplicates a message.

Design-bound implementation packages reference existing ModelGovernance and
Blueprint services, exact NOW/TARGET/task/blueprint revisions, allowed/protected
paths, source evidence and actual design disposition. A changed dependency makes
the package stale before continuation. Do not implement another model-history or
blueprint store. General Maintenance remains usable without model prerequisites.
Integration of those separately developed services needs a refreshed preflight and
selected combined baseline before dependent code changes.

Shared cross-worktree claims/IDs need one selected authority, not branch-local
registries advertised as global exclusion. Until that capability is implemented,
use the existing manual preflight and explicit coordinator sequencing. The current
known ID collisions remain integration work; preserving both histories is required.

## PGL1: observable work before orchestration

Add `sdp-controller overview --project ROOT --state DIR`, with readable default
and `--json`. It inspects every registered run in that one store, using a single
persisted generation. No host process, registration, new turn or approval occurs.
This is a registered-work view, not a complete project board or a live scheduler.

Show run/attempt/thread identity, task, selected Session/step/card/plan, routine state,
current coverage and pinned input revision. Show the last recorded host status and
timestamp as historical observations, with gaps and current liveness unknown.
Assigned role comes from controller provenance; never infer it from a run name.
Show observed tool-completion count/last identity, explicitly requested skill count
and acceptance separately from unobserved compliance. Child mapping and owner
acceptance stay unknown/unobserved. No transcript, owner-input content or capability
credentials belong in the overview. Terminal strings must not execute control codes.

Per-run source checks occur during inspection and are not an atomic filesystem
snapshot. Unavailable/stale records must be visible; operational-store corruption
fails the operation. Missing operational state yields an empty, unbound view without
creating a directory. Sort by stable run identity. Never derive live state from an
arbitrary observation-age threshold. This is the first part of handoff scenario 5;
all multi-agent orchestration scenarios remain open.

## Verification and next contract

PGL1 tests must distinguish two runs, stale coverage, no attempt, old/reconnected
observations, completed host versus unfinished routine, and terminal-control text.
Run the compiled command against the owner's existing lab state read-only and verify
both project and state bytes remain unchanged. This tests inspection, not a new
model workflow or independent review. Freeze a separate development binary for use.

Next PGL2 specifies and implements bounded SAD/implementation assignments with
revision gates, persisted identity, explicit capabilities and an executable
prepare/inspect/refuse path. Then dispatch/inbox/continuation and direct conversation
follow, with the six handoff scenarios retained in the plan. Online Steering needs
separate selected-host evidence and cannot be proved by local skill files.

## PGL2-C1 preparation contract — 2026-10-10

Owner T015 continues KB052. Select one runnable preparation increment before adding
launch authority: trusted local controller prepares a SAD or Master assignment,
inspects its pinned inputs, and refuses stale inputs or overlapping local scope.
Prepared assignments are not runs, attempts or host threads. They share the existing
checkout store, lock and atomic generation publisher; no second management ledger
or model/blueprint database is created.

SAD requires an active DesignPlan or ArchitecturePlan; Master requires an active
ImplementationPlan. Reuse canonical card/plan/Session/authority source checks and
management lineage. Keep the older Maintenance route unchanged. The API pins the
complete coverage revision and references a controller-selected authority source;
it does not interpret natural-language authorization or prove an owner disposition.
Worker capabilities cannot prepare assignments. A caller-supplied role is only the
assignment label; the trusted adapter determines the caller capability.

A SAD package pins NOW. A Master package pins NOW, TARGET, task, blueprint and a
design-disposition record. Each is an existing regular project file with SHA-256,
read through the core's confined/symlink-refusing file reader. These are exact file
references, not new ModelGovernance identities or an assertion that those services
validated or approved them. Provider-aware revision resolution and authenticated
owner design disposition remain PGL2-C2 prerequisites; inspection explicitly says
`dispatch: unavailable` and `designAcceptance: not-assessed`. A payload saying
approved cannot grant launch. The prepared package also retains bounded outcome,
writable/excluded paths and return criteria.

Prepare requires explicit checkout binding, expected current generation, a fresh
operation UUID and matching coverage/file digests. Under the existing store lock,
reject parent/child/equal writable-path overlaps with any prepared assignment or
nonterminal legacy run in this checkout. Paths are literal, not globs; excluded
subpaths do not subtract from this conservative overlap check. No global worktree
reservation is claimed. Other checkout preflight/claims remain external prerequisites.
Validate sources again before publishing; files are not locked against outside
writers, so inspection and future dispatch must always recheck them.

Persist immutable prepared identity/revision 1 plus the idempotency receipt in the
same generation. Same operation/content retry returns the original prepared record,
including after an uncertain committed save; conflicting content refuses. It does
not return a fresh eligibility claim. Recheck through inspect. The operation-ID
namespace is shared with legacy receipts; collisions refuse in both directions.
Inspect lists assignments in stable ID order and reports changed/missing sources.
No TTL expires a scope claim. Retirement/revision/continuation belongs to PGL2-C2;
this first increment never releases scope merely because a source is stale.

The first preparation explicitly upgrades that store to operational profile
`sdp-project-governance-state/0.2`. New readers retain 0.1 support. Old binaries must
refuse the new generation rather than drop prepared records; use a separate PGL2
trial store and never upgrade the owner's original C6 lab implicitly. Application
Maintenance wire profile and worker MCP stay unchanged. No stable distribution or
online host is installed. Role skill adoption from PLR1 is deliberately sequenced
with the integrator-owned distribution boundary before role launch, not copied over
release manifests during this preparation slice.

Commands: `sdp-controller assignment context|prepare|inspect --project ROOT --state
DIR`, JSON input via `--input FILE` for context/prepare, optional `--assignment UUID`
for inspect. Context reports canonical coverage plus the operational generation;
prepare persists only after the checks above. All responses identify their narrow
preparation boundary. No model/account access occurs. A reproducible separate
local trial must show SAD and Master preparation, retry after restart, conflicting
scope refusal and a changed TARGET producing stale inspection.

## PGL2-C2a revision and retirement contract — 2026-10-18

Owner Session0007 T018 selects the PG-local continuation. Reuse the existing
atomic generation and immutable receipts instead of allocating a replacement ID
on every edit: an assignment retains its ID, role and selection, with monotonically
increasing revision. Revision supplies the complete new package, current coverage
and artifact pins plus a bounded reason. Reject overlapping claims except the
assignment's own old claim; publish replacement and receipt together. The old scope
is released only when the new revision is committed. Stale preparations retain claims
until explicitly revised or retired; source drift is never implicit cancellation.

Retirement is terminal for this prepared identity. The trusted controller supplies
assignment ID, expected revision/generation and a reason. It can retire stale work
even when source files or canonical coverage are unavailable: releasing an unlaunched
claim must not require repairing obsolete inputs. Retirement preserves the last
package and its receipt history, increments revision and releases the local claim.
It is not completion, owner acceptance, cancellation of a running agent or deletion.
New work after retirement needs a new preparation ID. Retained records still count
toward the existing bounded store capacity; no garbage collection is added.

Both operations require the existing trusted controller/checkout binding; worker
capabilities cannot invoke them. Expected generation and assignment revision prevent
lost updates. Exact retries return the original immutable receipt even after later
revision/retirement; changed content or cross-operation UUID reuse refuses. Receipts
and generation history retain every prior version. Existing legacy run claims still
block conflicting replacement; retired preparations no longer block legacy starts.

`assignment revise|retire --input FILE` adds lifecycle wire profile
`sdp-governance-assignment/2`. First successful lifecycle mutation upgrades only the
selected store to `sdp-project-governance-state/0.3`; C1 binaries must refuse it.
New code reads 0.1/0.2/0.3 and never downgrades 0.3 on later preparation or legacy
writes. C1 prepare/context requests and immutable receipts retain wire profile 1;
inspection of a 0.3 store uses profile 2 and reports retired records explicitly.
Older stores remain unchanged by reads. Design acceptance stays not-assessed and
dispatch stays unavailable. Provider-backed disposition and dispatch remain C2 work.

Verification must cover restart/retry, stale input refresh, atomic scope replacement,
retirement with missing sources, terminal refusal, concurrent CAS, before/after-HEAD
fault recovery, binding/worker refusal, profile downgrade prevention, legacy scope
release and compiled CLI/schema behavior. Use a separate synthetic trial store;
preserve the owner's existing lab attempts.

## PGL2-C2b provider-backed design check — 2026-10-18

Owner Session0007 T020 continues the plan on integrated Actions main. Before recording
an owner disposition, establish an exact, reproducible review basis. Existing
Blueprint readiness is attributed evidence, not owner design approval; reusing it as
approval would violate the role contract. This slice supplies read-only checks and
a review digest; actual trusted owner input remains the next C2 increment.

`assignment design-check --input FILE` selects a prepared Master by ID, expected
revision/generation and a canonical Blueprint assignment ID/event. Resolve the
binding through the existing `blueprints.AssignmentViews` API and live models through
`model.Snapshot`. No second provider ledger, copied lifecycle or provider mutation
is introduced. Keep the provider's existing read budgets and unknown results.

Require the active PG plan to match the binding; the provider assignment must be
draft, ready or assigned. Hold, executing, review and terminal states require their
own workflow disposition. The retained bundle must validate and be current. Its
blueprint.json must be the prepared blueprint pin; NOW/TARGET pins must name each
resolved artifact's metadata file. Compare full model captures, not metadata alone:
an uncommitted WORK edit must stale the check even if its metadata file is unchanged.
The prepared task pin must match the retained task digest. Require passing structural
constraints and exact literal membership of each PG writable path in the task's
AllowedCodePaths; do not infer wildcard or directory expansion. This is a conservative
preflight bound, not a claim of code conformance or complete impact coverage.

Read and check twice, including operational generation, preparation revision and
canonical coverage. Any observed drift refuses a current result. A successful
`sdp-governance-design-check/1` response binds the exact prepared package, coverage,
provider assignment/event/binding, live model captures and task to a reviewRevision.
Return preliminary/model coverage and provider readiness separately. Blocked checks
have no reviewRevision. File inspection is optimistic; future disposition/dispatch
must run the check again. No source locks or global checkout exclusion are implied.

All results say designAcceptance not-recorded and dispatch unavailable. File text,
Blueprint assessment, frozen model, returned digest or request role cannot establish
owner approval. This read operation does not upgrade operational storage or change
the existing preparation inspection/worker MCP contract. A later owner-input increment
must bind its disposition to this review revision and invalidate it on change.

## PGL2-C2c owner decision records — 2026-10-18

Owner Session0007 T022 authorizes continuation while reading the walkthrough. This
is implementation authority, not an owner decision on a prepared test assignment.
Reuse the checkout's atomic generation and immutable operation receipts to record
approved, changes-requested or rejected against the exact C2b review revision.
Keep this operational record separate from canonical card/plan lifecycle: neither
approval nor rejection closes a card, changes a model or starts an agent.

Select a trusted local-controller adapter, consistent with the existing Caller
boundary. The core receives owner input separately from the package request and
rejects worker callers. The CLI displays the bounded proposal and captures an
explicit decision and comment from stdin, with a separately supplied owner label.
Blank input, EOF and cancel write no decision. No decision is inferred from source
files, technical readiness, silence or the request's role. Source/actor/comment are
attributed controller observations, not cryptographic identity verification. Same-user
shell access remains outside this boundary; an agent with unrestricted shell access
can invoke local commands. Conversation capture and stronger authentication are not
claimed. This adapter permits testing without selecting the final conversation/UI.

Store the complete checked review basis and the captured input, timestamp and operation
identity. A later decision supersedes the effective decision for that assignment;
immutable receipts preserve earlier decisions. Exact retries return their historical
receipt, never a renewed acceptance claim. Reusing an ID for changed content or any
other operation refuses. A retry of the local CLI retrieves a matching receipt before
prompting, so an uncertain successful write does not ask for another decision.

Before recording, require current generation, preparation revision, exact review digest
and a current provider check; repeat after collecting input and immediately before
publication. Failed or stale checks do not upgrade storage. A new read-only design-review
operation reports the latest recorded decision and its current applicability. It says
stale if the current checked basis differs or cannot be checked, including retirement.
Exact source restoration can match again: this checks current facts, not unobserved
transient file edits. A new preparation revision changes the basis even with equal
source bytes. Rejection/changes-requested on a current basis supersedes approval.

First successful decision upgrades only its chosen store to profile 0.4. New code
reads older profiles and preserves 0.4 on every later operation; older binaries refuse
it. Existing preparation inspection remains explicitly not-assessed; C2b design-check
retains its technical-only not-recorded result. Use design-review for effective owner
disposition. No new MCP privileges are exposed, and dispatch remains unavailable
even for an approved, current package. The external disposition file pin from C1 is
still a source reference; it cannot substitute for this controller record.

Verification: real provider fixtures, current decision/restart/retry, changed-source
and preparation invalidation, supersession, worker/identity/collision refusal,
concurrent decision CAS, before/after-HEAD recovery, stale input while prompting,
profile preservation and old-reader refusal. Use a separate synthetic lab and label
its scripted owner input as fixture evidence; never record the actual owner as its
approver. No broad dispatch or final interaction design is delivered here.

## PGL2-C2d — Activate approved implementation work, 2026-10-19

Owner T024 authorizes continuing the plan toward Session completion. Connect a
current approved Master preparation to the existing run/controller pipeline; do
not create a second execution engine or mislabel ImplementationPlan work as
Maintenance. SAD launch, responsive leader and multi-run control remain subsequent
orchestration work. This is a runnable boundary, not just another eligibility flag.

The privileged assignment activate operation atomically creates one ready run and
an immutable preparation-to-run binding. Require exact checkout/generation,
preparation revision, review revision and current approved decision operation.
Recheck the actual providers and local scope before publication. Exact retry returns
the original receipt; another operation cannot create a second run for that identity.
Preparation cannot be revised or retired after activation: the run owns execution
and interruption/replanning, and its preparation claim remains held. Do not release
scope merely because a host becomes idle or a design becomes stale.

Use an explicit implementation routine, with the existing implement/verify/review/
record progression and evidence gates. Run context resolves ImplementationPlan
coverage and the bound approved package. Read/status/overview expose stale sources;
launch, steer, report, resume, retry, bind and reattach reject changed design or a
superseding decision. Stop/fail/scope-change and host observation bookkeeping remain
possible after drift. The controller includes the exact approved package in launch
context and rechecks before host submission. Source checks remain optimistic; this
is not filesystem locking or a sandbox preventing all writes after a concurrent edit.

Reuse the existing owner run command and restricted worker MCP set. Launch-bound
context_read resolves the actual run's coverage; general Maintenance route selection
remains unchanged. No new worker launch or owner-decision privilege is granted.
Host finish still does not complete the routine or grant independent review.

Successful activation upgrades only its store to 0.5, preserving older profiles and
0.4 decisions; later writes must not downgrade. The binding stays in operational
state rather than adding caller-editable authority fields to Run. The additive
implementation routine is selected only by the new activation entry point; old
readers refuse the new store. Canonical management/model records remain untouched.
Local scope checks are not global reservations; the existing manual cross-worktree
preflight remains required. The checkout-wide controller leases still serialize
owner control until the PGL3 concurrency contract is delivered.

Verify real provider-created approval through activation and the built controller
with a protocol fixture, current/stale worker context, drift before launch/steer,
one-run retry/CAS and uncertain launch recovery. Record actual host trials separately
from protocol fixtures; fixture-owner decisions are not actual owner acceptance.

## PGL3-C1 — Per-run controller isolation, 2026-10-19

Continue the selected PGL3 execution path by removing the checkout-wide lifetime
bottleneck without weakening exclusion with existing controllers. New controllers
hold a shared lock on the existing controller/live-controller lock files, then an
exclusive per-run lock underneath them. Two disjoint runs can own different hosts;
two controllers for one run still refuse. Legacy exclusive checkout leases conflict
with every new shared lease in both acquisition orders. Preserve the short exclusive
store mutation lock and all source/scope/retry checks. No new operational schema is
needed because lease layout does not alter stored run meaning or allow old writers
to race a new active controller.

Use validated existing run IDs before creating per-run directories. A worker caller,
invalid ID or missing run cannot acquire a lease. If the second lock fails, release
the shared compatibility lock. Process exit releases both OS locks; directories are
not liveness records. Read-only observers still acquire no control lease. Distinct
runs do not share a host connection; the existing host client remains serialized.

This is a prerequisite for the responsive leader and two-Master scenario, not proof
of either. Verify same-run and legacy/new exclusion in both orders, independent-run
leases, process-exit release, no stale global lock after partial acquisition, and
controller integration. Keep automatic report-policy design separate: the actual
native refusal is retained and the owner's interaction preference is pending.
