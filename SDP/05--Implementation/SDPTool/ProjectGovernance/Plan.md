# PGI / PGL — Governed work and project leadership

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0018 |
| project | SDP |
| state | active |
| PlanType | ImplementationPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDPTOOL |
| source | PLAN-SDP-0017; KB-SDP-038 and KB-SDP-052 with KB036/037/042 contributions |

## Outcome and authority

Deliver one useful owner-to-agent Maintenance workflow through a local Codex client,
shared SDPTool process core and agent adapter. Restore the same Session after an
interrupted attempt and return candidate-bound evidence with a truthful disposition.
The owner continuation on 2026-10-06 selects PGI1-M1 from this handoff. The plan
is active. The owner continuation on 2026-10-07 selects PGI2-M1; PGI3 is selected by the next owner continuation; PGI4–PGI5 remain planned.

The [design](../../../04--Design/SDPTool/ProjectGovernance/Design.md),
[acceptance cases](../../../04--Design/SDPTool/ProjectGovernance/Acceptance.md) and
[Session 0007](../../../Sessions/session-%230007--Project_governance.md) govern
scope. Read their owner constraints before implementation. Record material design
contradictions rather than choosing a new product boundary silently.

## Existing behavior and boundaries

Preserve SDPTool discovery/navigation/output and existing manual Session documents;
reuse existing project identity and record validators. Keep ModelGovernance files,
model history and KB050 semantic blueprints with their separately selected work.
Do not make their completion a prerequisite to a plain Maintenance assignment.
KB043's full timeline, all routine families, distribution, release, GUI, remote
hosting and complete shell/file interception remain outside this plan.

## Git and assignment policy

BranchPolicy current means one working branch for the selected implementation,
with milestone commits. At actual start, establish an isolated checkout/branch
`sdp/project-governance-pilot` from the committed PGD1 delivery, or document an
already isolated equivalent; do not switch the concurrent ModelGovernance shared
worktree. Check Git status before creation and staging. Preserve unrelated files.
Follow existing phase-push/PR authority; merge and release need their own explicit
authorization. No push, merge or release is part of this planned handoff.

Use the authorized role/host delegation rules. An independent review must use a
context independent of implementation. If the current host does not support the
required separation or trusted child attribution, retain that integration gap;
do not turn a self-check into an independent review by changing its label.

## Runnable increments

| Phase / milestone | Delivered behavior | Required evidence | State |
| --- | --- | --- | --- |
| PGI1-M1 Core and inspection | Selected project + status request yields context/route without creating work; authorized fixture Maintenance creates a pinned durable run visible through a small CLI read/transition path | PG-A01/02/03/08/10/18; exact JSON schemas; core restart and stale-input tests | delivered; bounded core/CLI verified |
| PGI2-M1 Adapters and protocol | Thin local MCP adapter and controller protocol adapter use the core; inspect compatible schemas, negotiate the selected SDK/host, capture bounded recorded events and reconcile mapped attempts | PG-A10/11/12/15/18; actual protocol negotiation and failure outputs; no simulated-as-live claims | delivered; bounded adapters and real idle protocol probe verified |
| PGI3-M1 Owner workflow | Minimal terminal client recovers Session, submits/steers work, presents route/status and launches the authorized fresh Master with bounded context; applicable child/review mapping is observed | PG-A02/04/05/11/12/13/15; managed account mode and one useful bounded task | in-progress; C6 explicit skill input/correlation verified; native child/review mapping remains unverified |
| PGI3-U1 Owner development trial | Owner can inspect current work and exercise the candidate before further feature expansion | Identified local engine, runnable trial, actual owner feedback and expectation/design reconciliation | owner reports trial performed; detailed acceptance unrecorded; KB052 now selected |
| PGI4-M1 Evidence and reconciliation | Candidate-bound evidence, review/rework and checked publication produce one durable Session outcome and correct card/plan disposition | PG-A06/07/09/13/16/17; crash and concurrent-change fixtures | planned |
| PGI5-M1 Integrated verification | The same isolated task survives observer close, interruption/recovery and scope correction with unrelated work preserved | All PG-A01–18 at their stated level; independent review and exact-candidate evidence | planned |

Each milestone is runnable at its stated boundary; no phase claims the whole client
is delivered because a library or mock passes. Core work can proceed with protocol
fixtures while real host compatibility is measured, but PGI3/PGI5 cannot pass with
only schema exports. A real integration trial is required for account and native
child/context claims. Log protocol gaps and adjust the bounded design before relying
on an unsupported field or fallback role mapping.

## Concrete first assignment

PGI1-M1 owns a new bounded governance core package and its CLI adapter under SDPTool,
plus isolated test fixtures and this plan/evidence. Do not refactor the existing
navigation facade wholesale. Read current AGENTS/skills and validate actual Go
availability first; `go` was absent from this session's PATH during preparation.
The module declares Go 1.26.0. Choose a compatible local toolchain through the
project's existing environment conventions; a README claim is not a tested build.

Implement strict application schemas, explicit checkout binding, immutable routine
definition/run identity, read-only routing evaluation, one guarded Maintenance run,
local operation receipts/CAS and restart recovery. Exercise these through the CLI
with a real temporary project fixture. Initially no model submission, account
configuration, existing project-record publication or broad catalog installation.
Do not add a second generic management ledger. Maintain current core documentation
and command help when behavior becomes available.

## Verification and closeout

Run focused state/recovery/adapter tests and existing SDPTool regression tests for
the affected commands, with race checks for shared state. Record unavailable
environment-dependent checks explicitly. Validate project-management and Toolkit
records, preserve original ledger bytes, and verify exact candidate identity.

Independent review assesses intent, real workflow, preservation and authority;
verification distinguishes unit fixtures, protocol probes and real host execution.
Completion requires all selected milestone evidence and explicit remaining-card
dispositions. A passing pilot may leave broader KB036/037/038/042 capability scope
in backlog. No card is closed merely because this plan is complete.

## Current result

PGI1-M1 is delivered in an isolated local clone based on 043c59c, which includes
the completed ModelGovernance delivery. [Evidence-PGI1](Evidence-PGI1.md) and its
[candidate manifest](candidate-PGI1.json) identify the tested core/CLI subset.
Go 1.27.1 full SDPTool race tests, vet and 21 observed wire-schema cases pass.
PGI2-M1 is also delivered: [adapter evidence](Evidence-PGI2.md),
[candidate manifest](candidate-PGI2.json) and [protocol pin](protocol-PGI2.json)
record shared-core MCP parity, bounded observation/recovery fixtures and actual
Codex 0.160.1 / MCP 2025-11-25 negotiation with SDK v1.8.0. No model turn is sent.
PGI3-C3 verifies explicit approval/reporting; C4 verifies same-host executing-turn
reattachment. [C5 evidence](Evidence-PGI3-C5.md) verifies bounded launch provenance
and offline inspection through an actual Maintenance trial. Assigned role,
host-reported context, local source digests and skill availability remain distinct;
skill use is explicitly unobserved. Final affected race/compiled CLI and schema
checks pass, with the post-live schema-only correction stated in the evidence.
[C6 evidence](Evidence-PGI3-C6.md) verifies explicit project skill inputs and exact
lost-response recovery through final binaries. Invocation/acceptance is distinct
from unobserved procedural compliance. Native child/review mapping keeps PGI3 open
under current delegation rules. PGI4/PGI5 publication and independent review remain
planned.

The milestone history now resides in the original repository on
`sdp/project-governance-pilot`, using the persistent worktree
`/home/warloc/git/SDP-project-governance`. The shared worktree's concurrent changes
remain untouched. Earlier filesystem restrictions are historical. Shared-history
integration is separate and must resolve the recorded private C1 event-ID
collisions without overwriting either history. No merge, push or release is claimed.


## PGI1 implementation refinements and next handoff

The state directory is located by canonical-root SHA-256 and contains the separate
checkout UUID binding. This makes read-only lookup possible before explicit bind;
no discovery read registers a checkout. Selected canonical management records are
checked without publishing them; the full ledger hash conservatively invalidates
advancement after any append. Fine-grained authority/event reconciliation belongs
to the adapter/evidence phases. The pilot ships one compiled routine definition,
not the complete method catalog. Requests and results use the exact PGI1 schema.

PGI2 should start from this delivered core, preserve CLI behavior, inspect the
current installed protocol/SDK and implement thin adapters with equivalent core
results and measured negotiation. Do not promote local Caller capability from
agent-provided role strings. Resume here proves core persistence only; PGI2/PGI3
must establish actual host-attempt mapping before claiming controller recovery.
Reviewed backlog dependencies KB036/037/042/043 remain separate and unchanged;
ModelGovernance, semantic blueprints and full KanBan TUI are not entry prerequisites.

## PGI2 selection — 2026-10-07

EVT-PM-SDP-000145: owner resumes Session0007. Reuse the clean isolated PGI1 clone on
`sdp/project-governance-pilot`. Original read-only Git limitations were historical;
continued isolation now preserves concurrent shared work. Implement thin adapters,
scoped agent authority and bounded observed attempt recovery; no model execution
is needed for protocol negotiation. Live owner/agent workflow remains PGI3.


## PGI2 delivery and next handoff — 2026-10-07

The live probe observes managed ChatGPT mode, initialization and MCP discovery in
an idle ephemeral thread. The core stores observed thread/attempt/epoch identities
and paged completed content; a disconnect requires explicit reconciliation.
Fixtures establish replay/authority guards, not recovery of an executing model.
The next bounded assignment is PGI3: the terminal owner workflow, durable outbound
intent and controller loop, fresh Master launch and actual host role mapping.
Keep PGI4 publication and PGI5 independent integrated verification distinct.
Backlog review leaves KB036/037/042/043 unchanged; no new prerequisite from
ModelGovernance, semantic blueprints or KanBan TUI is introduced.

## PGI3 selection — 2026-10-07

EVT-PM-SDP-000148: owner continues. Work from committed PGI2 `12542b2` in the
isolated clone. Current writable roots include /tmp and SDL only; shared SDPTool
and Session updates cannot be returned in this host. Keep the Session, card, plan
and evidence together in the isolated checkout and report the publication gap.


## PGI3 implementation checkpoint — 2026-10-07

PGI3-M1-C1 adds durable owner-input dispatch, a fresh assignment-bound Master
launch, managed-mode/MCP checks, line-oriented steering/interruption/scope stop,
read-only status and conservative mapped-history reconciliation. The complete
isolated race suite, vet, application schemas and captured RPC schemas pass.
The compiled workflow uses a recorded fake host, never a real model.

The actual rebuilt idle probe exits with EOF before initialization. PGI3 remains
in progress; no real model task, running-turn reconnection or native role evidence
is claimed. Deliberate resolution of an unmapped ambiguous send remains unavailable.
Current host rules prohibit delegation, so the candidate does not grant native
children or pretend to measure a hierarchy. The next work is to restore a permitted
host connection and establish the required live evidence, then resolve remaining
recovery gaps. PGI4/PGI5 acceptance is unchanged. The shared checkout is read-only
outside SDL; this checkpoint and its Session records are retained only in the clone.


## PGI3-C2 — live workflow and correlation recovery, 2026-10-07

The owner confirms restored Git access and requests continued plan execution.
PGI1, PGI2 and PGI3-C1 commits are now imported into this repository's
sdp/project-governance-pilot branch. Development uses the persistent worktree
/home/warloc/git/SDP-project-governance; concurrent changes in SDP-vNow stay intact.
The earlier filesystem restriction is historical. No merge, push or release is
selected. Shared-history integration must reconcile private C1 event-ID collisions
with later blueprint/SDUI records before merging; neither history is overwritten.

PGI3-C2 measures Codex 0.161.0, actual Maintenance execution and steering,
observer closure, interruption/recovery and response-loss correlation. Core recovery
uses a unique observed clientId plus exact payload; it never automatically resends.
The live MCP report was refused because approval was required under the current
never policy. PGI3 stays in progress pending a deliberate owner-interaction path,
remaining executing-turn reattachment and applicable role/provenance evidence.
PGI4 and PGI5 retain evidence publication and independent integrated verification.
[PGI3-C2 evidence](Evidence-PGI3-C2.md) distinguishes these results from C1 fixtures.

2026-10-07T22:07:07.491506+00:00: EVT-PM-SDP-000153 — PGI3-C2 real workflow and correlated lost-response recovery verified; approval interaction and remaining provenance/recovery obligations keep PGI3 active. Persistent worktree replaces historical Git restriction.


## PGI3-C3 selection and result — 2026-10-08

The owner continued after C2's recommendation, selecting the explicit interaction
and disposable verification. This does not authorize arbitrary future requests.
C3 implements and verifies the request/response path, exact approved-call revision
handling and bounded contention recovery. Initial policy and lock failures are
retained in evidence, not relabeled as passing workflows. The final actual report
reached verifying and all excluded fixture files were preserved. Next: remaining
PGI3 executing-turn reattachment and applicable provenance evidence, then PGI4/PGI5.
The same branch/commit policy applies; no merge, push or release is selected.

2026-10-07T23:07:29.127244+00:00: EVT-PM-SDP-000156 — PGI3-C3 selected and verified: explicit per-request MCP approval/report and bounded contention recovery. PGI3 remains active for executing-turn reattachment and applicable provenance.

## PGI3-C4 selection and result — 2026-10-08

Owner continuation selects the next executing-turn recovery boundary. Private
same-host attachment now preserves live execution and the original scoped MCP
capability after controlled detach or abrupt controller loss. Both actual final
binary Maintenance journeys pass with one turn and unchanged exclusions. Initial
socket-alias refusals remain recorded failures. See [C4 evidence](Evidence-PGI3-C4.md).
Next: applicable native role/skill/context provenance; do not treat instructions
or a missing installed skill collection as observed compliance. PGI4/PGI5 remain
planned. Backlog/onHold review leaves KB036/037/042/043 and the separate ModelGovernance,
SDUI and blueprint workstreams unchanged. No merge, push or release is selected.

2026-10-08T15:46:29.121589+00:00: EVT-PM-SDP-000169 — PGI3-C4 actual executing-turn reattachment verified after graceful detach and abrupt controller loss; applicable provenance keeps PGI3 in progress.

## PGI3-C5 selection and result — 2026-10-09

Owner continuation selects bounded provenance. Controller-only context capture
and offline inspection now distinguish assigned role, reported context and skill
availability. The actual fixture includes maintained SDP skills; discovery and local
hashes pass without treating agent declarations as invocation. See [C5 evidence](Evidence-PGI3-C5.md)
for the live/runtime and later schema-only correction boundary. Next: explicit skill
invocation/correlation and applicable native role mapping under current delegation
rules. PGI4/PGI5 remain planned. Backlog/onHold review retains broader KB036/037/042/043
scope and separate ModelGovernance, SDUI and blueprint workstreams. No push, merge
or release is selected.

2026-10-09T02:32:53.879846+00:00: EVT-PM-SDP-000170 — PGI3-C5 bounded launch provenance verified; explicit skill invocation/correlation and applicable native role mapping keep PGI3 active.

## PGI3-C6 selection and result — 2026-10-09

Owner continuation selects explicit skill invocation/correlation. Final binaries
complete a skill-equipped Maintenance workflow and a deliberately dropped-response
recovery with exact ordered skill identities and one attempt/turn. Offline projection
reports observed input acceptance, never compliance. See [C6 evidence](Evidence-PGI3-C6.md).
Native child/review mapping remains unverified under current no-delegation rules;
retain that boundary before full PGI3 acceptance. Next: resolve its applicability
and prepare the PGI4 candidate/evidence/publication handoff. PGI4/PGI5 remain planned.
Backlog/onHold review preserves broader KB036/037/042/043 and separate ModelGovernance,
SDUI and blueprint scope. No push, merge or release is selected.

2026-10-09T02:49:51.131951+00:00: EVT-PM-SDP-000171 — PGI3-C6 explicit skill input and exact lost-response correlation verified on final binaries; native child/review mapping remains unverified.

## PGI3-U1 priority correction — owner direction, 2026-10-09

The owner cannot yet assess the workflow from documents and agent-run tests.
Prepare a separate local PG lab first; include both project/work overview and
interaction during an actual task. This replaces native mapping/PGI4 handoff as
the immediate next step. Further feature expansion follows the first owner feedback
and design reconciliation; bounded trial preparation/fixes can proceed now.

[Owner trial procedure](Owner-Trials.md) and [preparation evidence](Evidence-PGI3-U1.md)
record the runnable candidate and honest boundary. Preparation is delivered; owner
use and acceptance are pending. This does not finish PGI3 or waive independent review.
KB051 retains generic action/invoke and shared consumer-test contract ownership.
Backlog review retains KB036/037/042/043 and separate blueprint/ModelGovernance/SDUI
scope. No release, shared installation replacement, push or merge is selected.

2026-10-09T10:02:29.063036+00:00: EVT-PM-SDP-000173 — Owner prioritizes local PG overview/interaction trial and feedback before feature expansion; dedicated lab prepared with pinned C6 engine. Owner experience remains pending.

## KB052 reconciliation — Session0007 T014, 2026-10-10

Owner explicitly prioritizes KB052 after trying the PG lab. Extend this active plan
with the [orchestration design](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md).
Preserve delivered PGI history and unfulfilled PGI3/4/5 gates; KB052 does not imply
acceptance, native delegation evidence or publication. PGL is now the primary route.
The previous demand for owner trial before feature expansion is satisfied by the
reported use and superseded by this explicit scope selection; no new permission is
required to begin bounded work. Role-skill changes in upstream PLR1 are guidance
inputs, not yet installed into this candidate. Model/blueprint integration remains
sequenced with those workstreams rather than copied into a second engine.

| Milestone | Runnable outcome | Required evidence / handoff scenarios | State |
| --- | --- | --- | --- |
| PGL1 | Readable/JSON overview of all registered work from existing operational state; lab inspection without launching work | Scenario 5 subset; empty store, multiple runs, stale records, historical host/skill/tool facts, no writes | delivered; [bounded evidence](Evidence-PGL1.md) |
| PGL2 | Prepare and inspect distinct SAD/implementation assignments; stale design or conflicting scope refuses dispatch eligibility | Scenario 2 guards; current authority, exact bundle revisions, explicit cross-worktree claim boundary, retry/CAS tests | in-progress; [C1 preparation delivered](Evidence-PGL2-C1.md), [C2a revision/retirement delivered](Evidence-PGL2-C2a.md), [C2b provider checks delivered](Evidence-PGL2-C2b.md), [C2c local decision adapter delivered](Evidence-PGL2-C2c.md), [C2d activation and native code trial delivered](Evidence-PGL2-C2d.md); reporting/coordination remain PGL3 |
| PGL3 | Responsive Project Leader plus two bounded Masters through the trusted controller; durable report inbox and same-Session phase continuation | Scenarios 1, 3, 4; real host execution, no duplicate launch after uncertain reply, delayed report versus new steering | planned |
| PGL4 | Direct SAD/Master conversation and project panel; revision-bound ModelGovernance/Blueprint handoff | Scenarios 2, 5; real TARGET change stales dependent implementation; roles/tools/children/freshness evidence; owner development trial | planned |
| PGL5 | Integrated recovery/evidence review and revision-bound remote Steering package | All six scenarios; independent review and selected online-host evidence; PGI4/5 obligations retained | planned |

PGL1 owns governance read projection, the standalone controller command and owning
PG records. It does not edit shared sdptool dispatch/output/dependencies, ModelGovernance,
Blueprints or XFMD. [Preflight](preflight-PGL1.json) identifies nine local worktrees;
shared untracked PG files match our prior committed PGI2 bytes. Leave them untouched.
Known KB051/private event collisions remain explicit integration prerequisites;
new events are allocated after scanning all known local ledgers. This is scoped
sequencing, not a distributed reservation guarantee. Before PGL2, refresh preflight,
adopt the canonical role delivery deliberately and finish its concrete mutation
contract. Git BranchPolicy current / CommitPolicy milestone remains unchanged.

2026-10-10: PGL1 delivers read-only human/JSON inspection and a separate lab preview.
Actual owner-trial state is preserved; controller/MCP and final core/command checks
pass at the stated levels. Next: PGL2 assignment preparation, current authority and
revision gates. Native role-skill adoption and shared integration remain explicit
prerequisites to the affected later changes. No full KB052 acceptance is claimed.

## PGL2 selection — Session0007 T015

Owner continues KB052. PGL2-C1 selects the runnable preparation/inspection/refusal
contract in Orchestration.md; PGL2 remains in progress until C2 supplies disposition,
provider-aware revision, revision/retirement and dispatch-boundary integration.
Existing Maintenance behavior and worker MCP privileges are preserved. Scoped
preflight reads the new coordination overview at 684b1f8; shared integration is
separate and its proposed worker allocation is not assumed delivered. Skill
adoption is sequenced with the shared distribution baseline before role launch.
C1 writes only the PG core, standalone controller and owning PG records; it uses a
new trial state and preserves the owner's C6 lab. No agents are launched.

2026-10-10: PGL2-C1 delivers the bounded preparation CLI/core and separate fixture
preview. [Evidence](Evidence-PGL2-C1.md) and [handoff](Handoff-PGL2-C1.md) retain
C2/provider/disposition/revision/retirement gates. No general scheduler is delivered.

## Integrated continuation — Session0007 T017, 2026-10-18

Remote main 8aa16c2 includes PGL2-C1 de4b795 and the ACT0 history/identity
reconciliation. The owning branch incorporated it without conflicts at c3877ae;
BranchPolicy current and CommitPolicy milestone remain unchanged. Preserve this
integrated ancestry and canonical records instead of replaying the old PG history.

Continue with PG-local assignment revision/retirement in governance and the standalone
controller after a fresh scoped preflight. Actions Session0011 owns generic action
dispatch, shared CLI/model dispatch and its own records. Provider integration,
shared contracts/dependencies and new event allocations require refreshed coordination
before affected writes. Main's concurrent-work procedure governs these checks.
Before delivery, incorporate current main, inspect concurrent committed and dirty
changes, test the combined candidate and prepare a PR limited to the new PG delta.
This is a clean integration route, not a guarantee about future conflicts or authority
to merge main. No C2 implementation or acceptance is claimed by this synchronization.

## PGL2-C2a — Preparation revision and retirement

Owner Session0007 T018 continues the bounded PG-local work. Implement stable-ID
revision and terminal retirement of unlaunched preparations through the core and
standalone controller. [Orchestration](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md#pgl2-c2a-revision-and-retirement-contract--2026-10-18)
owns exact behavior, compatibility and verification obligations. Shared Actions
dispatch, provider services and runtime agent launch remain outside this slice.
[Preflight](preflight-PGL2-C2a.json) records the integrated baseline and relevant
worktrees; stale primary PG copies remain excluded. Use canonical post-ACT0 IDs.

[Evidence](Evidence-PGL2-C2a.md) and [candidate](candidate-PGL2-C2a.json) record the
delivery boundary and checks. This delivers revision/retirement only. PGL2 remains
in progress: actual controller-bound design disposition and provider-aware revision/
dispatch prerequisites are next. No owner acceptance or multi-agent execution is
implied by a prepared, revised or retired assignment.

2026-10-18: EVT-PM-SDP-000198 records PGL2-C2a delivery; EVT-KB-SDP-000370 updates KB052.

## PGL2-C2b — Provider-backed review basis

Session0007 T020 continues the plan. The integrated provider APIs distinguish
technical readiness from owner authority. Select a read-only Master design check
first, binding the prepared assignment to canonical Blueprint and live ModelGovernance
revisions before a later actual owner-input disposition. Reuse exported provider
services without modifying their storage, lifecycle or shared Actions interface.
[Orchestration](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md#pgl2-c2b-provider-backed-design-check--2026-10-18)
defines the contract and conservative exact code-path check.
[Preflight](preflight-PGL2-C2b.json) records the scoped continuation boundary.
Verification requires real provider-created artifacts/bundles, uncommitted TARGET
drift, wrong binding/task/scope, canonical event drift, optimistic read recheck and
proof of no writes. Actual design disposition and dispatch remain unfinished C2 work.

[Evidence](Evidence-PGL2-C2b.md), [candidate](candidate-PGL2-C2b.json) and the actual
[CLI trial](trial-PGL2-C2b.json) record delivery. The provider-backed review basis is
implemented and verified; actual owner-input disposition and dispatch remain next.
PGL2 stays in progress. EVT-PM-SDP-000200 / EVT-KB-SDP-000372 record this boundary.

## PGL2-C2c — Controller-attributed owner design decisions

Session0007 T022 owner authorizes continuing while reviewing the walkthrough. Reuse
the existing local controller authority and atomic operational store for explicit
revision-bound decisions; do not select a final conversation or XFMD UI. The
[contract](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md#pgl2-c2c-owner-decision-records--2026-10-18)
separates package requests, trusted captured input and current applicability.
[Preflight](preflight-PGL2-C2c.json) bounds writes to PG and its own records.

[Evidence](Evidence-PGL2-C2c.md), [candidate](candidate-PGL2-C2c.json) and
[actual trial](trial-PGL2-C2c.json) record the delivered local adapter, immutable
history, stale refusal, supersession and 0.4 compatibility. This does not authenticate
humans against same-user shell processes or deliver conversation capture. The final
interaction remains open to owner feedback. No actual owner trial acceptance is
inferred from scripted fixture decisions. PGL2 remains in progress for dispatch
prerequisites; PGL3–5 remain planned. EVT-PM-SDP-000201 / EVT-KB-SDP-000373 record
this bounded delivery; BranchPolicy and CommitPolicy are unchanged.

## Session completion route — T024, 2026-10-19

The owner asks to continue the plan toward completion of Session0007. Existing
implementation authority covers the remaining selected milestones; it does not
waive their evidence, owner-feedback or independent-review boundaries. Preserve
BranchPolicy current and CommitPolicy milestone. Do not turn each internal check
into a new owner permission gate, and do not close the Session after a local slice.

| Remaining outcome | Current next delivery / evidence |
| --- | --- |
| Finish PGL2 handoff into execution | C2d activates an approved implementation package into the existing controller/MCP runtime; stale/revoked authority refuses advancement; protocol and native trials are distinguished |
| PGL3 leader and coordinated execution | Responsive leader conversation, two disjoint Masters, durable report delivery, phase/idle continuation and uncertain-send recovery; no global-exclusion claims from local locks |
| PGL4 integrated owner interaction | Direct SAD/Master dialogue, actual model handoff and readable project panel with tools/children/freshness; incorporate the owner's walkthrough/trial feedback |
| PGL5 and retained PGI gates | Candidate-bound review/evidence/publication, recovery and independent review; remote Steering capability tested on the selected host, or explicitly dispositioned with the owner |
| Session closeout | Reconcile S3/S4/S6 and related cards with actual evidence; retain the broader KB043 timeline as deferred scope with its explicit disposition, not an invented delivered feature |

C2d follows the [activation contract](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md#pgl2-c2d--activate-approved-implementation-work-2026-10-19)
and [scoped preflight](preflight-PGL2-C2d.json). The real-host trial must use the
verified complete host package; the installed Codex version is not silently replaced.
A host turn ending without implementation is a failed trial, not acceptance.

## PGL2-C2d result — T024

[Evidence](Evidence-PGL2-C2d.md), [candidate](candidate-PGL2-C2d.json) and the separate
protocol/native trial records establish approved implementation activation, current
run authority, immutable retry and bounded execution through the existing controller.
The native code and external behavioral tests pass; report submission is refused
under never. Host completion therefore does not advance or accept the assignment.
The first failed native attempt remains recorded rather than overwritten.

PGL2's Master activation boundary is delivered. PGL2 remains in-progress until its
remaining handoff restrictions are reconciled with PGL3 execution: scope release/
replanning, report integration and SAD dispatch are not silently declared complete.
Next is PGL3's responsive coordination/durable report contract and per-run controller
isolation, retaining the actual native report-policy limitation. PGL4/5 and retained
PGI gates still apply. EVT-PM-SDP-000202 / EVT-KB-SDP-000374 record this delivery.
