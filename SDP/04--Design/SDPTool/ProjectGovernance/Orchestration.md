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

The T014 baseline had one checkout-local run store, run/attempt/host identities,
controller-only launch, restricted worker MCP, replay receipts, exact skill inputs
and correlated recovery. At that baseline both leases were checkout-wide, context and routine selection
were Maintenance-only, and the role validator accepted only Master. Subsequent
PGL2 deliveries add SAD/Master preparation and approved Master implementation;
PGL3-C1 adds per-run leases and C2 adds automatic reports/read-only inbox. See the
[current evidence](../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL3-C2.md).
Responsive leader consumption, SAD dispatch and native parallel/child mapping
remain unproven. Prompt changes alone do not implement those remaining contracts.

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
reconciled against the original attempt before any new launch. The current worker MCP retains its restricted authority. PGL3-C3 added attention
read/request to the original five operations; a fully bound worker now has seven
tools. It still cannot launch, resume, decide acceptance or publish canonical work.

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
blueprint store. Maintenance proven not to affect design remains usable without model prerequisites;
T039 below makes SDL WORK mandatory for every design-affecting change.
Integration of those separately developed services needs a refreshed preflight and
selected combined baseline before dependent code changes.

Shared cross-worktree claims/IDs need one selected authority, not branch-local
registries advertised as global exclusion. Until that capability is implemented,
use the existing manual preflight and explicit coordinator sequencing. The original branch-local identity collisions were reconciled through the
Session0011 history import, which preserves exact source histories. That does not
implement global cross-worktree identity allocation. Current candidates must still
pass management replay and the concurrent-work preflight.

See the [current state responsibility map](State-Boundaries.md) for provider,
assignment, run and acceptance boundaries after integration.

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

## T025 — Automatic reporting and attention, 2026-10-19

The owner resolves the T024 reporting question: maximize automatic operation and
status/result submission within authorized work, targeting projects that run for
hours. Routine reporting must not require an owner decision on every submission.
A genuine question, changed scope or reserved decision uses a durable attention
request. Independent work continues where its authority and dependencies allow it.
No response or timeout is permission. Host restrictions still need a supported,
measured integration; this decision does not rewrite the earlier native report refusal.

Continue PGL3 with automatic local reports, correlated inbox and structured attention/
reply semantics. Retain assigned capabilities, source/revision checks, receipt
identity, evidence-versus-acceptance distinction and explicit continuation scope.
Do not implement blanket acceptance of host tool requests as a substitute for the
bounded report contract.

The owner also requests multi-machine oversight/control and an installable desktop/
mobile web client. [KB054](../../../KanBan/backlog/%23054--Proposal--PG-attention-and-remote-coordination.md)
and the [architecture proposal](../../../03--Architecture/ProjectGovernance/Attention-and-Remote-Coordination.md)
register that extension. Local runtime, proposed self-hosted hub and responsive PWA
share one message identity and answer disposition. This is proposed architecture,
not an implemented remote product or revival of obsolete SDP_Analyzer. Select its
own bounded execution plan after the local contract; retain this plan's existing
completion boundaries. No production deployment or account changes occur here.

## PGL3-C2 — Automatic report delivery and local inbox (2026-10-19)

New launches may explicitly select automatic reporting. The controller records that
selection in the launch payload and configures only the assignment-bound
`sdp_governance.run_transition` tool with `approval_mode = "approve"`. The tool retains
its truthful write annotation and only report/fail/scope-change capabilities. This
uses the pinned 0.161.0 host's supported configuration; it does not change global
approval settings, answer arbitrary approval requests or accept a submitted result.
Existing launches retain their previous policy unless automatic reporting is selected.

A report submitted from a current running agent may tolerate intervening tool/message
observations only if every immutable revision proves unchanged decision state. It
must not skip approval requests, owner inputs, lifecycle transitions, reconnects or
source/approval drift. Controller calls retain strict revision comparison. Exact
operation retry remains idempotent.

The local report inbox projects report transition receipts, retaining operation ID,
reported revision, input revision and text separately from current run state and
coverage. Reading does not create or mutate state. Reports survive process restart
and later run transitions; repeated reads/retries do not create duplicate messages.
This initial view has no consumed/accepted state, leader acknowledgement or automatic
continuation. Structured attention questions/replies are the next local contract.

T026 owner boundary: the hub, web/mobile application and notification delivery belong
to a separate project. SDP/SDPTool owns local execution, durable reports and attention,
and its eventual adapter contract. KB054 retains the external coordination handoff;
no destination repository or implementation is invented.

Host configuration references: [official MCP configuration](https://learn.chatgpt.com/docs/extend/mcp)
and [pinned 0.161.0 configuration schema](https://github.com/openai/codex/blob/rust-v0.161.0/codex-rs/core/config.schema.json).

The initial C2 native experiment mapped automatic reporting to host `auto`; the host
correctly refused the write under `never`. Pinned host implementation distinguishes
annotation-based `auto` from per-tool `approve`. The owner's T025 automatic reporting
selection is represented by `approve` for this bound tool only on a fresh launch.
No rejected operation is replayed, existing attempt policy changed or arbitrary tool
preapproved. The failed trial remains separate evidence. Source:
[pinned tool approval semantics](https://github.com/openai/codex/blob/rust-v0.161.0/codex-rs/core/src/mcp_tool_call.rs).

## PGL3-C3 — Durable local questions and attributed replies (2026-10-21)

Owner T027 continues the selected local attention/reply step. Extend the existing
atomic run store and immutable operation receipts instead of introducing a competing
mailbox ledger. Explicit question writes select operational profile 0.6; old stores
remain readable and old readers refuse the extended store. All writers retain 0.6.

A launch-bound agent or controller may ask a question on a running, current run:
nonempty prompt, zero to five stable option IDs/labels, free text always available,
and a blocking flag. IDs are operation UUIDs. Questions retain source revision,
run/attempt binding, creation time and a decision-basis digest. Host observations
may advance the run revision without invalidating the question; owner steering,
source/design changes, lifecycle changes and replacement/disconnected attempts do
invalidate applicability. Sibling questions/replies do not invalidate each other.
The question/reply is project data, not a design approval or expanded mandate.

Only a trusted controller may attribute a reply captured from its input surface.
The request names the question and exact basis, never an owner identity. Retain text,
optional selected option, actor/source, receipt ID and server timestamp. One reply
per question; exact retries return immutable history, conflicting repeats refuse.
An answered record is not proof of agent consumption or application. A controller
may explicitly withdraw an obsolete question with a reason; preserve its history.

Open blocking questions prevent report/start/resume advancement for that run even
when stale, until answered or explicitly withdrawn. Other runs and stop/fail/scope
bookkeeping remain available. This is a governance advancement guard, not a claim
that arbitrary host file/tool execution is forcibly suspended. Agents are instructed
to stop dependent work and return when waiting; a reply does not automatically start
another host turn in this increment. Leader consumption/continuation remains next.

Add write-annotated attention_request and read-only attention_read MCP tools,
restricted to the launch-bound run/attempt. The read returns current applicability,
not just historical records; optional questionId bounds the result for large histories.
No owner-reply tool is given to agents. Automatic-mode launches preapprove this bounded reporting tool
alongside run_transition, retaining all core guards. The controller CLI lists current
applicability, captures a reply against a displayed question, and supports explicit
withdrawal. Terminal output quotes project-controlled content; JSON supports local
clients. No hub, notification transport or new global host configuration is added.

Verification: ask/read/reply across reopen; option/free-text input; idempotency and
wrong-run/epoch/actor refusal; source, steering and terminal-state staleness; sibling
independence; two-run progress; report guard; withdrawal; legacy reader refusal and
all writer profile preservation; real MCP transport and CLI capture. Native question/
reply continuation and independent review remain distinct unfulfilled evidence.

## PGL3-C4 — Guarded attention continuation (2026-10-21)

T028 continues the authorized local workflow. When a running assignment's latest
host attempt has completed cleanly while awaiting input, answered current questions
may select a fresh bounded Master attempt in the same run. Reuse the existing launch
intent/claim/acknowledgement protocol and fresh thread isolation. Extending an ended
turn in place would require new host correlation/recovery semantics; a fresh attempt
fits the existing capability boundary and preserves the prior attempt as history.
This increment is the controller primitive and optional bounded watcher for future
Project Leader coordination, not the complete conversational leader.

A trusted continuation selects exact question/reply/basis references, requires all
blocking questions resolved, unchanged sources/authority, no observation gap, completed
prior attempt and no unresolved outbound input. Preparation atomically records a new
launch input and makes this same run ready. Claim rechecks the question's pre-prepare
basis and current authority before any host effect. Inputs retain immutable references;
prepared/sending/uncertain requests are never blindly resent. Existing terminal-history
recovery may prove an accepted send but does not authorize another send.

The fresh payload retains the original owner request, current assignment/context,
question/reply history and explicit continuation instructions. Sandbox, approval/
reporting policy and selected skill bytes must match the original launch. Reply input
is not design approval, scope expansion or acceptance. A report remains verifying and
cannot trigger this continuation. New questions in a later attempt require new replies.

Operational profile 0.7 adds reply references to launch input; legacy readers refuse
and every writer preserves it. Attention inspection derives dispatch-prepared/sending/
uncertain/acknowledged from those inputs. Host acknowledgement is not proof of semantic
consumption or correct implementation; retain that distinction in views and evidence.

Expose explicit continuation and an opt-in bounded automatic watcher on the local
controller. Default launch behavior remains unchanged. A watcher stays alive after a
clean waiting turn, checks for current answers and starts at most its configured number
of fresh attempts. Waiting does not append receipts or poll a model. It stops on report,
failure, stale authority, unresolved delivery or lifetime exhaustion. It starts no new
assignment, grants no native delegation and does not implement the external hub.

Verification: durable exact reply binding; source/steering/stale/duplicate refusals;
claim before host effects; lost response/crash never causes duplicate launch; unchanged
policies/skills; native ask, delayed reply, fresh attempt, real code change and report;
protected files, prior lab and legacy reader refusal. Multi-hour capacity and full
leader/report consumption, direct SAD dialogue and independent review remain separate.


## PGL3-C5 — Durable report acknowledgement (2026-10-24)

Refinement of the selected durable inbox contract: a trusted local controller can
record that a particular report was processed by a named coordinator, with a bounded
note. This is an acknowledgement, not evidence review, owner acceptance, publication
or permission to launch another attempt. No automatic host dispatch is selected.
The worker MCP retains its seven tools and cannot acknowledge reports.

Reuse the atomic operational store rather than a second queue or canonical ledger.
Pin report operation ID, run ID, reported revision and input revision; require the
inspected store generation and current run revision at first write. Require current
source/assignment coverage and a verifying run. Delayed reports after steering,
source changes or run disposition refuse, with the historical report still visible.
An identical operation retry returns its recorded result even after later changes;
a different request using that operation ID conflicts. Each report has at most one
acknowledgement; a second operation cannot silently overwrite it. Caller-supplied
coordinator attribution is local-controller provenance, not human authentication.

Store acknowledgements separately from run progress so reading/acknowledging one
report does not change routine state, attempt identity or another run's revision.
Adopt operational profile 0.8 only at the explicit acknowledgement write; older
readers must refuse. Preserve 0.8 through later preparation/attention/continuation
writes. Extend the inbox with a version 2 projection only once acknowledgements
exist; retain version 1 for untouched stores. Semantic agent consumption remains
unobserved even when acknowledgement is recorded.

Expose `sdp-controller inbox acknowledge --project ROOT --state DIR --input FILE`.
The request pins the inspected basis; it does not silently refresh stale inputs.
Verification covers persistence/restart, exact retry and conflicting duplicate,
authority and stale-source/generation/state refusals without mutation, no run
advancement, coexistence with attention/continuation profiles, legacy reader
refusal and actual compiled CLI/schema output. Native leader reasoning, automatic
continuation, independent review and the remaining PGI/PGL gates stay separate.


## PGL3-C6 — Shared coordination view and responsive local input (2026-10-24)

Refine the existing controller-owned inbox/event-loop boundary into a usable local
entry point. Reuse registered work, report receipts/acknowledgements and attention
records; do not create another queue or acceptance state. A bounded combined view
pins one operational generation, lists blocking questions before pending reports,
and retains explicit stale coverage, unknown liveness and historical dispositions.
Stable item IDs and a content revision support refresh coalescing. Pagination binds
that revision and filter; a changed view refuses an old cursor. Inspecting creates
no work. Filesystem coverage remains separate from the operational snapshot.

Expose the view through `sdp-controller coordinate` and an explicitly configured
coordinator MCP adapter. Default inspection and bound worker adapters keep their
existing tool sets. `sdp-mcp --coordinator ID` is a launcher capability, mutually
exclusive with run/attempt/epoch. It exposes only `coordination_read` and
`report_acknowledge`, using C5's exact request guards and launcher-bound coordinator
attribution. No owner-reply, design-decision, launch, resume, completion or publication
tool is added. A payload role string cannot select this adapter. These restrictions
make it useful for a host conversation without giving that conversation owner approval.

The optional terminal watch loop polls snapshots while accepting bounded JSON lines
on stdin. Explicit refresh, pinned report acknowledgement and attributed attention
reply reuse existing services. Reply requires the launcher's owner attribution;
report acknowledgement requires its coordinator identity. No free-text inference of
approval or automatic model dispatch. Process one input at a time, prioritize queued
input over polling, and do not hold run-controller leases or mutation locks while
waiting. Failed/invalid input emits an error and preserves the loop; output failure
stops this observer, not worker execution. EOF permits watching until timeout; quit
and cancellation stop it. On restart, pending work comes from durable core records.

This console is an intermediate coordination interface, not yet an AI Project Leader
conversation scheduler or a native wake-up guarantee. Verify actual compiled CLI
and stdio MCP, independent worker progress during watch, replies/acknowledgements,
restart, pagination/stale refusals and authority separation. Full native leader
turn scheduling, uncertainty recovery and independent review remain required later.


## PGL3-C7 — Serialized Project Leader wake-up and mapped recovery (2026-10-24)

Use the existing atomic checkout store and host client. Add bounded durable owner
messages and leader dispatch receipts in operational profile 0.9, preserving all
older profiles and refusing old readers after the first explicit leader write.
These are operational delivery facts, not a second management ledger or acceptance.
One active leader dispatch per store; worker run/controller leases remain independent.

A short decision turn uses a fresh persistent thread with read-only sandbox, never
approval policy, explicit selected project skill inputs and C6 coordinator MCP only.
The launcher supplies coordinator identity, verified host version and selected skills.
The controller injects bounded pending owner messages before report/question notices.
Stable notice keys exclude polling timestamps and observation counters. Handled keys
are not requeued merely because a report remains verifying. Delivery is not semantic
consumption, acknowledgement or acceptance; the leader may use C5 acknowledgement.
The controller cannot execute a proposed worker launch or owner decision from prose.

Persist prepare, opening-thread, bound-thread and sending-turn intent before each
external effect. New owner input or changed selected notice/skill basis invalidates
an unsent preparation; supersede only prepared/bound work, then reselect pending
inputs. Never automatically replace opening/sending/running uncertainty. A missing
thread/start response remains unresolved because the thread is unmapped. For a
mapped thread, recovery requires the same host version, project cwd and exactly
one matching client message ID plus exact input/skill digest in complete history.
Missing, duplicated or conflicting acceptance refuses; an empty/idle thread is not
permission to resend. Capture only visible completed messages, never private reasoning.
A completed/failed/interrupted host turn releases the leader slot, not assignment
acceptance. Reports and owner messages arriving during it remain pending.

Bound state and turn counts explicitly; no sustained-hours capacity claim. Owner
messages use trusted local attribution; model tools cannot impersonate the owner.
Inspect/status and message enqueue do not launch a model. An explicit leader run
command enables bounded dispatch/polling; queued input takes precedence over new
report work. Preserve the local coordination console and worker capabilities.

Verify with service and actual compiled command/protocol fixtures: serialization,
owner input during active work, restart, exactly correlated lost-response recovery,
no resend after ambiguous send, stale unsent supersession, profile coexistence and
old-reader refusal. Native host scheduling remains unverified until an authorized
native trial; role skill input does not establish procedural compliance or review.

C7 implementation refinements: owner enqueue uses stable operation/content identity
with expectedRevision=0; leader transitions use CAS. This lets a local sender retry
without retaining a moving queue revision. Store and effect leases remain separate.
The exact host MCP inventory must contain only the coordinator server; other inherited
servers or an incomplete page refuse before turn submission. No global config edits.

Limits are 64 retained messages, 16 retained dispatches and 256 leader receipts,
including superseded preparations. Each dispatch selects up to eight messages and
eight notices within 64 KiB of encoded JSON; message text is at most 4096 bytes and
visible answers at most 64 KiB. Fit messages in original order before notices.
Normalize only notice fingerprints; retain actionable current revisions in payloads.
No compaction/reset path or sustained-hours readiness is implied. See the
[controller contract](../../../../SDPTool/governancecontroller/README.md#serialized-project-leader-turns--pgl3-c7).

### C8 native refinement — unsent thread lifetime

An owned native host does not necessarily retain an unsubmitted thread across
process restart. C7's recorded-host fixture did not exercise that lifetime. A new
controller instance therefore supersedes a bound thread it did not create locally;
the durable send claim is absent, so pending input can safely select a fresh thread.
Opening uncertainty still refuses, while sending/running always recover exactly in
the original mapping. This implements the existing safe-unsent-replacement boundary
without adding reset authority or widening the store profile.

The lab uses process-local Codex Apps/remote-plugin/multi-agent restrictions. The
native default Apps catalogue exceeded the existing 1 MiB frame limit and exposed
an extra server; the exclusive MCP guard and bounded transport remain intact.
A restricted launch supplied exactly the two coordinator tools. No global config
was edited and no native test agent was granted delegation.


## T039 — Mandatory governance admission and lifecycle supervision

Owner Session0007 T039, 2026-10-24, makes the following explicit requirements of
KB052/PLAN-SDP-0018. They refine the selected PG workflow; they do not claim new
runtime checks have been delivered. The reusable detailed procedure is canonical
[sdp-project-leader](../../../../Skills/sdp-project-leader/SKILL.md), version 1.1.0,
which every Project Leader coordination turn must load and apply.

1. **Admission:** inspect intent against existing Session goals, card outcomes,
   selected plans, scope and authority. Reuse covered work. A coherent goal revision
   updates coverage before dependent action. A logically separate change requires
   its own KB card(s), Session and selected plan before agent assignment; tell the
   owner why, prepare under existing authority, and retain unrelated running work.
   A separate existing workstream may be reused. Do not equate same repository with
   same goal. Pure read/status questions need no new work hierarchy.
2. **Assignment contract:** bind every SAD/Master/delegated assignment to cards,
   Session step, selected plan/milestone, exact revisions, authority, role, allowed
   and protected scope, candidate/worktree and return criteria before launch/resume.
   Children inherit this chain; Master maintains it. Missing records cannot be
   retrospectively supplied as launch authority. Use existing Session/plan schemas.
3. **Design assessment:** record design impact and Blueprint applicability as two
   distinct decisions. Every design-affecting change requires an actual SDL change
   in identified WORK before dependent implementation. Missing/inadequate models
   require model preparation, never a Maintenance-label bypass. Required Blueprint
   context covers overall system and relevant NOW/TARGET detail; genuinely local
   tasks may justify no separate Blueprint without waiving WORK. Classify contracts,
   dependencies, behavior and state changes by effect, not patch size or task label.
4. **Coordination:** inspect current and paused work, committed/dirty code, model
   ancestry/elements, plans, Session writers and shared IDs across known worktrees.
   Resolve overlaps with scoped ownership, sequencing and an integrator. Merge
   models into mutable WORK through existing ModelGovernance, preserve provenance,
   validate semantic consistency and refresh affected views/reviews/handoffs. A
   clean Git merge or separate branch is insufficient. Until shared atomic claims
   exist, explicit sequencing/preflight is required and coverage limits stay visible.
5. **Supervision:** recheck these obligations on input, phase reports, scope/design
   changes, dispatch/resume and integration. Retain current checks, reasons, revisions
   and pending model-release state in existing durable records. Invalid coverage
   stops only dependent work; the owner should not have to remind PG of the method.
6. **Closeout:** design-changing work retains a model-release obligation through
   implementation. Integrate/validate WORK, freeze CANDIDATE, obtain the actual
   required integration-owner acceptance and promote the exact candidate to RELEASE,
   linked to code verification and predecessor. Recheck the current accepted model
   head; competing heads/conflicts block promotion. Do not close the whole assignment
   while that obligation is pending. An explicit scoped transfer/defer names a
   successor and preserves the gap. Model-only evidence cannot satisfy implemented
   behavior, and model promotion grants no repository merge or software publication.

The model proposes semantic fit/design assessments with reasons; core/controller
checks must enforce required records, references/revisions, authority and transition
preconditions. A Boolean self-attestation is insufficient. Explicit skill input,
observed host acknowledgment and verified compliance remain distinct evidence.
Unavailable skill or required checks blocks dependent dispatch; it cannot silently
fall back to an ungoverned prompt. No polling/background guarantee comes from the
skill alone. Use existing ModelGovernance and Blueprint providers; no competing
model store or general workflow-family expansion is selected.

[PG-A19–24](Acceptance.md) define the positive/negative acceptance cases. PGF1
must account for pending model disposition; PGF2 implements admission/supervision;
PGF3 completes SDL integration/release; PGF5 tests interruption across these gates;
PGF7 verifies the combined native workflow. Before dependent runtime changes,
prepare the corresponding SDL WORK under this rule; existing coarse model coverage
must be inspected, not assumed sufficient. This requirement/skill-maintenance turn
is not evidence of that SDL handoff or runtime implementation.


## T040 — Discussion: executable Session workflows and KanBan governance

Owner Session0007 T040 asks how current PG relates to durable state machines that
supervise Session work, require correct card updates/moves, survive interruption,
and accept Project Leader/MCP and owner/CLI interaction with tool/activity insight.
This is an analysis and recommendation for discussion, not adoption of a new Session
schema or authorization to implement a generic workflow interpreter. KB052 owns the
current PG integration; KB037 already records the broader state-machine capability.
The T039 mandatory governance requirements remain selected.

### Observed baseline at 851001f

- [KanBan reader](../../../../SDPTool/kanban.go) builds navigable card/status nodes
  and checks CardState/folder/history consistency; the shell `kanban status|state`
  helper lists/groups cards. The canonical KanBan guide still prescribes manual
  update/move/history/link maintenance. No general checked card-lifecycle mutation
  API was found in those surfaces. One shared board has status directories; this
  discussion does not require a new board per Session.
- [Run core](../../../../SDPTool/governance/run.go) already implements durable,
  revision-checked start/interruption/resume/report/retry/scope-change transitions.
  Complete still refuses. The pinned bounded Maintenance routine is not execution
  of every authored Session roadmap step or arbitrary Markdown instructions.
- [Publication](../../../../SDPTool/governance/publication_records.go) builds checked
  result additions to Session/card/plan and histories. Its recoverable file intents
  preserve preimages, but deliberately retain card/plan/Session lifecycle and paths.
  This is a usable foundation for reliable effects, not implemented card movement.
- [Controller](../../../../SDPTool/governancecontroller/README.md) supplies overview,
  coordination/watch, questions/replies and bounded host dispatch/recovery/attach.
  Observations describe retained host/tool events with explicit freshness limits;
  they are not a universal observer of arbitrary external work or private reasoning.
- [MCP](../../../../SDPTool/governancemcp/README.md) uses the shared core with restricted
  worker operations and separate coordinator read/acknowledgement. The coordinator
  does not currently have the full owner/controller launch/resume capability.
- [Context](../../../../SDPTool/governance/context.go) pins selected documents and
  canonical management history. Current conservative invalidation on ledger changes
  must be reconciled before routine card updates can advance a long Session; blindly
  accepting any new input digest would erase the existing stale-authority protection.
- The [Session contract](../../../Sessions/README.md) remains an authored goal,
  roadmap and journal; operational state is local and checkout-bound. A Git checkout
  of Session Markdown alone cannot recover live host identities and pending effects.

### Recommended ownership and execution model — proposed

Extend the existing shared PG core rather than build a second engine or use the
Project Leader conversation as the persistent scheduler. Keep four meanings distinct:

1. Session/plan defines the goal, selected cards, step dependencies and acceptance.
2. A versioned workflow definition specifies allowed transitions, checks and effects.
3. A durable workflow instance binds that definition to one Session/step/assignment,
   remembers current state, evidence, pending actions and revisions, and survives
   replacement of any agent conversation. Multiple bounded instances can coexist.
4. Host attempts execute individual assignments; tool/activity observations and
   unknown execution state are distinct from business progress and card lifecycle.

Do not interpret free-form Markdown as executable instructions or create a second
editable roadmap that competes with it. Prefer a small explicitly structured,
validated execution binding to the existing Session/plan and selected versioned
routines. Exact serialization/location and revision/adoption policy remain design
choices. Authored narrative stays authored; runtime views should be generated or
clearly labeled projections. A plan edit needs explicit reconciliation to the running
instance, not silent rebinding or a new engine on every resume. Reuse existing
management history as the canonical card/plan history and the operational store for
execution; define effects between them instead of making two stores own CardState.

### Checked KanBan effects — proposed

Offer semantic lifecycle operations through a shared service: select work, begin
assignment, record progress/evidence, submit a concrete owner gate, rework, defer,
complete and explicitly reopen/transfer. These are operation concepts, not installed
CLI commands. Each checks current card/revision, coverage, authority, evidence and
remaining scope; then updates prose/worklog, CardState, physical location, index,
links and canonical ledger through one recoverable operation identity.

The engine should perform required record effects as part of a valid transition,
or leave the transition visibly pending/blocked; it must not merely remind an agent
to edit files. Cross-file effects need durable intent, preimage checks and recovery,
not a claim that several filesystem renames form one atomic transaction. Do not
launch dependent work until required effects are committed/reconciled. Scope claims
must account for shared ledgers/links and other Session writers.

Selection normally means active/ready; actual assignment work means in-progress.
A phase result updates progress without necessarily moving the card. gate-review
means a concrete owner review, not every automated reviewer turn. completed requires
all of that card's selected outcome, required evidence and model release/disposition.
A run crash does not mean onHold or cancellation; unfinished work retains its actual
management state. A card shared across phases/Sessions cannot close when only one
contribution finishes. Keep waiting-for-input, verification and host interruption
in execution state instead of inventing a CardState for every runtime condition.

Manual edits remain possible in Git projects. Detect them as revision drift, explain
what changed and reconcile explicitly; never silently overwrite them or advertise
filesystem-wide enforcement against arbitrary same-user shell writes.

### Shared entrypoints and resume — proposed

Owner CLI and Project Leader MCP must act on the same instance/revision and command
rules through role-appropriate capabilities. An observer reads the same events and
state; closing it does not stop the workflow. Leader requests may propose actions,
while the core validates prerequisites and the trusted controller performs host
launch/reconnect. CLI access is not a second engine or an approval bypass. A shared
long-lived local service is an option for scheduling, but a new daemon/database is
not assumed necessary before evaluating the existing store/controller loop.

Resume must load the existing instance and pinned workflow version, reconcile pending
file/host effects, verify project/model/plan revisions and scope claims, locate the
existing host attempt, then report runnable versus waiting/conflict/unknown. Never
start another Master merely because a response or observer was lost. Give an owner
an inspectable next action; an unmet design/review/release gate remains a gate from
CLI too. Cross-machine recovery additionally needs an explicit portable operational
handoff; repository documents alone are insufficient.

### Bounded follow-up and acceptance proposal

First agree the owner-facing transition table and map each transition to existing
PGF1/2/3/5/7. Then prepare SDL WORK for the selected design, extend the current core
and checked publication effects, and expose the same guarded actions through the
owner CLI and appropriately privileged leader adapter. Avoid a general-purpose
workflow language or all-routine-family expansion before one complete path works.

In pg_lab, prove one card from selection through progress, failed review/rework,
required model release and completion; force interruption between card move/history
writes and before/after host launch. Resume from owner CLI while the leader is absent,
then reconnect the leader to the same instance without duplicate launch or events.
Include a shared-card partial result, external card edit, concurrent Session conflict,
a routine version change and an unrelated ledger append. Observe real tool activity
with source/freshness; reject attempts to bypass gates through either entrypoint.
These are proposed acceptance additions; none is claimed tested or selected by this
analysis. Resolve this boundary with the owner before dependent implementation.


## T041 — Selected concretization

Owner T041 accepts the T040 recommendation and selects the bounded design in
[Session execution and checked KanBan lifecycle](Session-Execution.md). It specifies
state ownership, the execution binding, semantic transitions and required effects,
coverage reconciliation, common role-bound entrypoints and recovery. PG-A25–30 and
PLAN-SDP-0018 carry its delivery gates. The initial SDL WORK and diagnostic Blueprint
are recorded in [T041 evidence](Evidence-T041.md). This supersedes T040's pending
choice of direction; exact runtime delivery and independent review remain open.
