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
is active. The owner continuation on 2026-10-07 selects PGI2-M1; PGI3 was selected by the next owner continuation. Subsequent owner continuations
activated PGI4 and PGL; PGI5 integrated verification remains open. The current tables
and T038 completion checkpoints supersede earlier next-step forecasts.

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
| PGI4-M1 Evidence and reconciliation | Candidate-bound evidence, review/rework and checked publication produce one durable Session outcome and correct card/plan disposition | PG-A06/07/09/13/16/17; crash and concurrent-change fixtures | in-progress; C1 evidence/rework, C2 observed review and C3 checked result publication/leader notices delivered; final run/lifecycle disposition remains |
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

## Completion roadmap

The [T038 completion checkpoints](#pg-completion-checkpoints--t038) consolidate the
remaining PGI/PGL work into one ordered route with explicit exit conditions. They
are planning checkpoints in this same plan, not replacement plans or delivered work.

## Current result

Latest delivery is [PGI4-C3](Evidence-PGI4-C3.md): checked bounded result publication
with crash/conflict recovery and compact review/publication notices through the
existing coordination/leader route. C1 evidence/rework and C2 observed review remain
prerequisites. PGI4-M1 stays in-progress for final run/lifecycle disposition; PGL2/4,
parallel/active-host/capacity and independent integrated verification remain.
Publication does not close broader cards/plans or imply owner acceptance.

### Historical PGI1–3 baseline

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
| PGL3 | Responsive Project Leader plus two bounded Masters through the trusted controller; durable report inbox and same-Session phase continuation | Scenarios 1, 3, 4; real host execution, no duplicate launch after uncertain reply, delayed report versus new steering | in-progress; [C1 per-run isolation delivered](Evidence-PGL3-C1.md), [C2 automatic reports and read-only inbox delivered](Evidence-PGL3-C2.md), [C3 durable attention/replies delivered](Evidence-PGL3-C3.md); C4 continuation, C5 acknowledgement, C6 console/MCP and [C7 serialized leader delivery/recovery](Evidence-PGL3-C7.md) delivered at bounded evidence scope; [C8 native bounded leader/worker and completed-history recovery](Evidence-PGL3-C8.md) pass; parallel and whole-workflow gates remain |
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

## PGL3-C1 selection — continued T024

Continue with per-run control isolation as the first prerequisite for two disjoint
Masters and a responsive leader. The [contract](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md#pgl3-c1--per-run-controller-isolation-2026-10-19)
and [preflight](preflight-PGL3-C1.json) preserve old/new exclusion and all authorization
checks. This does not select automatic answers to host approval requests. A focused
owner interaction question is pending: automatic status/result delivery within an
already-authorized assignment versus a decision for each submission. Existing never
policy remains intact while unaffected control isolation proceeds.

[Evidence](Evidence-PGL3-C1.md) and [candidate](candidate-PGL3-C1.json) establish the
C1 lease boundary with service/protocol tests, process-exit recovery and unchanged
legacy exclusion. The full SDPTool suite, affected race tests and vet pass. This
is not native parallel-agent or leader evidence. PM203 / KB375 record the delivery;
report interaction remains a pending owner choice, not an assumed authorization.

## T025 — Reporting decision and remote extension

Owner selects automatic status/result reporting and continued work within granted
mandates, with attention requests for actual questions/decisions. This resolves the
T024 pending preference; it does not waive host restrictions or accept unfinished
work. Next PGL3 increment must establish the supported automatic report path, durable
local inbox and revision-bound attention/reply behavior, then the responsive leader
and parallel native evidence. A waiting item blocks its dependencies, not all projects.

[KB054](../../../KanBan/backlog/%23054--Proposal--PG-attention-and-remote-coordination.md)
registers a self-hosted multi-machine hub and responsive/installable web client.
Its [architecture proposal](../../../03--Architecture/ProjectGovernance/Attention-and-Remote-Coordination.md)
recommends a local runtime plus hub plus PWA and an optional separately verified
ChatGPT MCP Events adapter. Follow-on runtime/hosting needs its own bounded plan;
registration does not silently enlarge PLAN-SDP-0018's completion criteria or revive
SDP_Analyzer. Existing branch/commit policy, remaining PGL/PGI gates and no-merge/
no-release boundaries remain. PM204 records the plan refinement, not implementation.

## PGL3-C2 — Automatic reporting and local report inbox

T026 owner confirms that the hub and web application belong to a separate project.
Continue local PG in SDP/SDPTool. KB054 is the coordination/handoff record until an
external project has a concrete destination; it does not authorize hub implementation
here. The current branch and milestone commit policies remain unchanged.

Select one bounded increment: an explicit automatic reporting mode for new controller
launches using the pinned host's documented per-tool policy, plus a read-only durable
report inbox projected from immutable transition receipts. Keep existing launches
compatible and all assignment, source, epoch and owner-decision checks intact. Reports
are submissions for verification, not acceptance, delivery acknowledgement or permission
to launch more work. No new persistence profile, network service or global host setting
is needed. Structured attention/replies and actual leader consumption remain subsequent
local work. Verify policy isolation, persistence, duplicate reports, stale authority,
unchanged approval refusal behavior and a fresh native implementation/report trial.

[Preflight](preflight-PGL3-C2.json) records the scoped worktree boundary. The native
trial must be separate from the previously refused report; retain that historical
failure. No merge, release or independent-review claim is authorized by this increment.


## PGL3-C2 result — T026

[Evidence](Evidence-PGL3-C2.md), [candidate](candidate-PGL3-C2.json) and the
[native trial](trial-PGL3-C2-native.json) establish automatic bound report submission
and a durable read-only inbox. The first incorrect host-policy mapping remains
[failed evidence](trial-PGL3-C2-failed.json). The corrected candidate passed actual
Go implementation, external tests and automatic report delivery into verifying;
no approval request was observed and no acceptance is inferred. Full tests/race/vet
pass. PM205 / KB378–379 record delivery and the separate hub/webapp project boundary.

PGL3 remains in progress. Next local delivery is typed attention/reply with current
revision binding, then leader consumption and responsive continuation. Receipt
capacity (currently 1024), restart/recovery and native parallel/multi-hour evidence
must be resolved before claiming sustained unattended operation. PGL2 scope closure/
SAD dispatch, PGL4/5 and the retained PGI gates remain. No merge/release is selected.

## PGL3-C3 selection — T027

Owner continues on 2026-10-21. Implement the local durable attention question/reply
contract in Orchestration: existing atomic store/receipts, scoped agent ask capability,
controller-attributed reply/withdrawal, readable CLI and revision applicability.
[Preflight](preflight-PGL3-C3.json) isolates PG from historical dirty primary copies
and newer main changes. Existing branch/milestone policy applies. No hub, automatic
new host turn, actual owner acceptance or full Session closeout is selected here.


## PGL3-C3 result — T027

[Evidence](Evidence-PGL3-C3.md), [candidate](candidate-PGL3-C3.json) and
[CLI trial](trial-PGL3-C3.json) record durable local questions, optional choices,
free-text replies, explicit withdrawal and current applicability. Blocking questions
guard only their run's advancement. Bound agents may ask/read; trusted controllers
capture attributed replies. No answer grants design approval or expanded scope.
Operational profile 0.6 is explicit; legacy readers refuse, all writers preserve it.

The isolated owner-lab ./attention helper displays one open question and supports
--reply --owner NAME. Tests cover restart/retry, stale authority/steering, capability
binding, sibling/other-run independence, MCP transport and compiled CLI input capture.
The synthetic CLI trial starts no model. Earlier native report evidence stays scoped
to C2; native attention/reply continuation is still required. PM206 / KB380 record
this delivery while the plan/card stay active.

Next: durable leader consumption and guarded continuation after a reply/report,
then responsive native parallel execution. Resolve bounded receipt/state capacity
and recovery before sustained-hour claims. SAD dispatch/scope closure, PGL4/5 and
retained PGI gates remain. Hub/webapp stay in the separate project; no merge/release.

## PGL3-C4 selection — T028

Continue with guarded reply delivery into a fresh attempt in the same authorized run,
using existing durable launch claims. Add an explicit continuation command and bounded
opt-in waiting loop; preserve launch policies/skills, prevent duplicate sends and retain
actual dispatch state in attention inspection. A fresh native trial must demonstrate
ask, delayed reply, continuation, code behavior and report. [Preflight](preflight-PGL3-C4.json)
selects only PG paths and current milestone branch/commit policy. Full Project Leader,
report acceptance, scope expansion, hub/webapp and sustained-hour claims remain outside
this increment. No merge/release or independent-review authority is inferred.


## PGL3-C4 result — T028

[Evidence](Evidence-PGL3-C4.md) and [candidate](candidate-PGL3-C4.json) bind guarded
continuation to exact recorded replies. A fresh attempt reuses the original mandate,
policies and skills, with durable prepare/claim/ack fencing uncertain or duplicate
sends. The explicit continue command and bounded opt-in watcher share this path.
Profile 0.7 retains dispatch references; older readers refuse and writers preserve it.

The native trial demonstrates ask, delayed fixture reply, one fresh attempt, selected
Go behavior, protected files and automatic report into verifying. Protocol tests
cover stale/changed authority, uncertain delivery, duplicate refusal without state
mutation and current source checks. Full tests, affected race/vet and actual JSON
schema checks are recorded in the candidate. The original lab and earlier previews
remain intact; ./continuation provides a read-only view. PM207 / KB381 record this
bounded delivery, with separate Traceability implementation/verification entries.

PGL3 remains active. Next: leader consumption of reports and responsive coordination,
with native parallel execution. Establish capacity/recovery before sustained-hour
claims. SAD dispatch/scope closure, PGL4/5 and retained PGI gates remain. No independent
review, owner acceptance, main integration or release is claimed. Hub/webapp remain
in the separate project; backlog/onHold dependency scopes are unchanged.

## PGC1 selection — Session0007 T029, 24 October 2026

Owner authorizes Codex 0.164.0 and compatible later versions, not a single-release
repin. [PGC1 contract/milestones](Codex-Compatibility.md) and KB057 bound this
compatibility increment. Retain existing BranchPolicy current / CommitPolicy
milestone. Use isolated equivalent branch sdp/codex-forward-compatibility from
consolidated c907e22; the historical pilot branch/commits remain preserved in main
history. No rewriting of earlier branch commitments or source evidence.
PGC1-A analysis/contract delivered; PGC1-B implementation starts. Later PG/leader,
independent review and publication gates remain unchanged. No merge/release selected.

## PGC1-B result — T029

[Evidence-PGC1](Evidence-PGC1.md) records stable minimum/forward eligibility,
actual version/attempt propagation, refusal/preservation checks and shared wire
schema alignment. Scoped package races pass. Final committed-candidate full affected
module checks and real no-model probe are PGC1-C next. Existing full native and
independent-review/leader gates remain open; no release/merge selected.

## PGC1-C result — T029

Code candidate 98b449c passes full affected SDPTool race/vet, Python Toolkit/
management/import checks, 60 actual wire-object schema checks and the real 0.164.0
no-model PG controller/MCP probe. Fresh private lab preparation reports the selected
source's correct minimum and stays ready. [Evidence](Evidence-PGC1.md) preserves
exact limits. KB057 closes for this bounded compatibility increment; PLAN-SDP-0018
and KB038/052 remain active for native workflow/independent-review, leader/report,
canonical publication and sustained recovery/capacity gates. Do not infer full PG
readiness, future-host certification, merge, release or new native task authority.

## T030 — SDL-only scope and model handoff

Owner defers further code/runtime context expansion and prioritizes completing the
existing PG/SDL/Blueprint route for first use. KB059 is onHold; no new runtime,
snapshot collector or source-analysis prerequisite. Preserve required existing
engineering/evidence checks and the existing Ponsse pilot model/MVP1 deferral.
Owner requests GPT-6 Astra high for continuation; this agent cannot change the
current client model. Authorization for existing work persists across the model
handoff; no new permission gate or implementation/release authority is inferred.


## PGL3-C5 selection — T031

Owner continues the existing completion route after the model handoff. Select
bounded durable report acknowledgement under the existing PGL3 inbox contract.
[Preflight](preflight-PGL3-C5.json) records the current clean single-worktree baseline;
[design](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md#pgl3-c5--durable-report-acknowledgement-2026-10-24)
defines identity, authority, retry, freshness and compatibility boundaries. Keep
BranchPolicy current / CommitPolicy milestone on sdp/codex-forward-compatibility.
Service and compiled CLI checks must prove acknowledgement without acceptance or
run advancement. No native model/subagent launch, Blueprint resumption, Ponsse
pilot, code/runtime context expansion, merge or release is selected by this step.


## PGL3-C5 result — T031

[Evidence](Evidence-PGL3-C5.md) records durable exact-report acknowledgement through
the trusted controller and compiled CLI. Full SDPTool race/vet, retained wire schema
checks, actual old-reader refusal, retry/restart/authority/freshness guards and
profile coexistence pass. All run state remains unchanged by acknowledgement.
Existing worker MCP tools/privileges remain unchanged. No native agent or independent
review is claimed. S6/PGL3 remain active: next is responsive leader coordination
using the durable report/reply path, followed by current-host integrated workflow,
capacity/recovery, SAD/scope closure and retained PGI4/5 evidence/publication gates.
Backlog/onHold review preserves KB036/037/042/043, external KB054, proposed KB058 and
SDL-only/KB059 deferral. No new scope selected; no merge or release.


## PGL3-C6 selection — T032

Owner continues. Deliver the existing coordinator inbox/event-loop boundary as one
bounded service/CLI/MCP workflow: combined paged work/report/attention view, explicit
coordinator MCP read/ack authority and a responsive terminal watch accepting pinned
acknowledgements and owner-attributed replies. Reuse C3/C5 stores and guards; no new
store profile, native model launch or second lifecycle queue. Follow the
[design](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md#pgl3-c6--shared-coordination-view-and-responsive-local-input-2026-10-24)
and [preflight](preflight-PGL3-C6.json). Branch/commit policy remains current/milestone.
Compiled CLI and stdio MCP must establish the useful operation and unchanged worker
capabilities. AI leader scheduling and whole-workflow/native/review gates remain open.


## PGL3-C6 result — T032

[Evidence](Evidence-PGL3-C6.md) records the delivered shared coordination view,
explicit coordinator MCP read/ack capability and responsive terminal watch with
attributed replies. Core pagination/freshness and authority guards, actual compiled
console/two-run service workflow, stdio MCP restart/retry, full affected race/vet
and ten final wire objects pass. Private development binaries are available; stable
installation/global MCP settings are unchanged. No new operational profile beyond
C5 acknowledgement; no native agent turn or independent review.

PGL3/S6 remain active. The local coordination interface is delivered; next is actual
AI Project Leader turn scheduling/wake-up with serialized delivery and uncertainty
recovery, followed by the retained current-host integrated/native, SAD/scope closure,
capacity, independent review and canonical evidence/publication gates. Manual
Sessions without registered runs remain outside the operational view. Broader
KB036/037/042/043, external hub KB054 and metrics KB058 stay in backlog; KB059 stays
onHold under SDL-only scope. Blueprint pause and Ponsse MVP1 deferral remain.


## PGL3-C7 selection — T033

Owner continues. Implement serialized short Project Leader turns and durable owner/
notice delivery under the existing PGL3 contract, with mapped lost-response recovery.
Use [design](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md#pgl3-c7--serialized-project-leader-wake-up-and-mapped-recovery-2026-10-24)
and [preflight](preflight-PGL3-C7.json). Existing current/milestone policy remains.
One active leader dispatch, explicit launch policy/skills, unchanged worker rights,
no prose-to-action execution. Operational profile changes preserve earlier facts.
Verify service and compiled protocol workflows; no native model/subagent launch or
independent-review claim in this execution. Retain whole-PG publication/review gates.


## PGL3-C7 result — T033

[Evidence](Evidence-PGL3-C7.md) records durable owner/notice delivery, one active
leader decision dispatch, explicit restricted skill/MCP launch and exact mapped
recovery without resend. Full SDPTool race/vet, compiled controller protocol restart,
owner input during uncertainty, actual prior-reader refusal and 20 output/context
schema checks pass. Profile 0.9 preserves earlier worker/report capabilities.
Private candidate binaries do not change the stable installation or global config.

No native model task, independent review or full PG acceptance is claimed. Next is
native leader/worker workflow evidence, followed by retained SAD/scope, capacity/
compaction, recovery/review and canonical publication gates. Bounded limits and
unmapped thread-start recovery remain explicit. Backlog/onHold review preserves
KB036/037/042/043/054/058, SDL-only KB059 deferral, Blueprint pause and Ponsse MVP1
deferral. Current-branch/milestone policy remains; no merge or release.


## PGL3-C8 selection — T034

Owner explicitly selects actual-model testing in /home/warloc/git/pg_lab and
authorizes its extension into a useful PG fixture. Keep the SDP current/milestone
policy; create a lab test branch from its clean main. Exercise native leader input,
actual bounded worker/report delivery, owner enqueue during work and mapped recovery.
Retain failures and fix scoped implementation defects without weakening the C7
identity/authority guards. No production owner acceptance, independent-review claim
or release is implied. See [preflight](preflight-PGL3-C8.json).


## PGL3-C8 result — T034

[Native evidence](Evidence-PGL3-C8.md) passes for three real leader turns and one
bounded worker in pg_lab, queued owner input, coordinator report acknowledgement and
two exact completed-history recoveries without resubmission. Native Codex 0.164.0
reports gpt-6-astra/high. External/project behavior checks, affected race/vet and
11 native output/context schemas pass. Corrected unsent-bound restart under C7's
existing safe-supersession boundary; source/binary identity and failed runs retained.

The isolated process disables unrelated Apps/plugins/delegation; global settings
are unchanged. The worker run remains verifying; the lab coverage is intentionally
pinned. Parent remains active for two-Master execution, active-host recovery,
SAD/scope closure, sustained capacity/compaction, independent review and canonical
publication. No merge/release, production owner acceptance or full PG readiness.


## PGI4-C1 selection — T035

Continue the authorized Session completion route at the evidence boundary. The
C8 run remains verifying: report acknowledgement cannot prove candidate correctness.
Implement bounded source capture, durable attributed evidence, current/historical
applicability and controller-only failed-evidence rework. Use the existing atomic
operational store, explicit profile extension and the current branch with a
milestone commit. [Preflight](preflight-PGI4-C1.json) records the scope.
PG-A06/07/08/16/17/18 govern this increment; completion/publication and independently
observed review remain gates for the next increment. PGL3 parallel/recovery/capacity
and PGL2/4 owner/SAD interaction remain tracked, not waived by this sequencing.

## PGI4-C1 result — T035

[Evidence](Evidence-PGI4-C1.md) records the bounded source/evidence/rework delivery.
Exact source inventories and attributed external artifacts survive restart and
retain historical meaning after relevant changes. A current failed test/review can
return verifying work to ready with the same Session and revoked old agent access.
The compiled pg_lab trial preserves C8, demonstrates real pass/fault/fail/rework/
restore behavior and provides `./pg evidence` for read-only inspection.

PGI4-M1 remains in-progress: independent review provenance and checked canonical
publication are still missing, so completion continues to refuse. Next delivery
must establish the review/publication contract and exact-candidate evidence without
promoting an attachment into an observed independent review. Retain PGL3's remaining
parallel/active-host/capacity gates and PGL2/4 SAD/owner interaction; neither this
increment nor the native trial closes the whole Session. Existing branch/milestone
policy and owner continuation authority remain; no merge or release selected.


## PGI4-C2 selection — T036

Continue the authorized completion route from C1 evidence to observed review
provenance. Implement a durable review assignment and read-only fresh-host dispatch,
pinning current candidate/report/test evidence and explicit reviewer skill. Retain
uncertain sends and recover only exact mapped history. Observe context separation;
do not infer it from attached review prose or claim native child hierarchy. Exercise
positive/negative actual-model reviews in pg_lab under existing test authority.
[Preflight](preflight-PGI4-C2.json) bounds this work. BranchPolicy current and
CommitPolicy milestone remain. Canonical publication and whole-PG closeout remain.


## PGI4-C2 result — T036

[Evidence](Evidence-PGI4-C2.md) records durable review assignments, exact instructions
and skill inputs, fresh read-only native contexts and recovered results. The final
pg_lab native reviewer passed correct code and identified the deliberate separator
defect. Lost turn/start and completed-history responses recovered through the same
original thread/turn, without replacement. The failed finding also returned to ready
through C1 rework, preserving its native review provenance and Session.

Observed review support is delivered at this boundary. This is not independent
review of the PG implementation or proof of native child hierarchy. PGI4-M1 stays
in-progress for checked canonical publication and its crash/conflict reconciliation.
Next: that publication boundary and review-result delivery into leader coordination,
then retained PGL2/4 owner/SAD, PGL3 parallel/active-host/capacity and PGI5 independent
verification. No completed card/Session, merge or release is inferred.

## PGI4-C3 selection — T037

Continue from ef757ae with the recorded green preflight. Implement checked bounded
result publication, exact crash/conflict recovery, and review-result delivery through
the existing coordination/leader route. Preserve active canonical lifecycle and
manual roadmap; routine completion and overall PG closeout remain separate gates.
Validate supported record sets through the existing management validator; incomplete
historical lab fixtures must refuse. Exercise a canonical disposable lab fixture and
retain C2 native review evidence as historical evidence, without relabeling fixture
reviews as native. No new model turns are required for filesystem publication tests.


### PGI4-C3 result and next handoff

Final full race/vet, 13 observed wire objects, actual compiled six-file publication
and repeat, canonical fixture validation and stdio MCP coordination pass. Existing
C2 native review results appear in the new notices without new model calls. Exact
source and limits are in the linked evidence; old readers refuse profile 0.12.

Next is guarded run completion and canonical work disposition using published
candidate/review facts, without confusing publication with owner acceptance or
closing broader cards. Retain orchestration and independent PGI5 gates. Backlog/
onHold reviewed: no new prerequisite from Blueprint expansion, timeline or hub.


## PG completion checkpoints — T038

Owner asks for a detailed route to finish PG. Baseline is b26340a, with PGI4-C3
bounded publication delivered. This turn plans and records the remaining work;
it does not implement these checkpoints or claim new test results. Existing
implementation authority persists; deployment/merge/release permissions remain
separate. BranchPolicy current and CommitPolicy milestone are unchanged.

### Completion definition and scope

The selected local PG workflow is finished when the owner can give work to Project
Leader, receive responsive coordination, use SAD/SDL design when relevant, authorize
bounded implementation, obtain testing and independent-context review, return
findings for rework, and record the correct completion/remaining-work disposition.
This must survive interruption and practical sustained use without duplicate work,
lost history, stale design authorization or manual copying between agent chats.

Current foundation: core/CLI/MCP, bounded assignment activation, questions/replies,
report delivery/acknowledgement, serialized leader delivery, exact candidate evidence,
native separate-context review and checked six-file result publication. These are
real increments, but their combination is not yet a complete accepted workflow.
In particular, complete still refuses; leader notices do not implement the full
trusted scheduling/continuation policy; SAD dispatch and the owner-facing integrated
workflow remain incomplete. Do not report a percentage from the count of increments.

| Checkpoint | Existing scope | Dependency | Current status |
| --- | --- | --- | --- |
| PGF1 Final work disposition | PGI4; PG-A06/07/09/13/16 | C3 | next; implementation required |
| PGF2 Leader-directed execution and continuation | PGL2/PGL3; handoff scenarios 1/3/4 | PGF1 | partial foundation; implementation and native evidence required |
| PGF3 SDL design handoff and replan | PGL2/PGL4; scenario 2, PG-A05/07/17 | PGF2 | partial preparation/provider guards; integrated SAD flow required |
| PGF4 Usable owner conversation and overview | PGL4; scenarios 4/5 | PGF2/3 | existing inspection/console; integrated workflow required |
| PGF5 Recovery across the complete workflow | PGI3/PGI5/PGL5; PG-A08/09/11/12/13/14 | PGF1–4 | bounded cases delivered; active-host/full-flow gaps remain |
| PGF6 Capacity and host compatibility | PGL3/PGL5; PG-A15/18 | PGF5 | limits explicit; sustained operation/rollover not proven |
| PGF7 Native integrated acceptance in pg_lab | PGI5/PGL5; PG-A01–18 and handoff 1–5 | PGF1–6 | not verified as one final candidate |
| PGF8 Independent product review and owner trial | PGI5/PGL4/PGL5 | PGF7 | PG implementation review and final owner workflow feedback open |
| PGF9 Scope disposition and development closeout | Session S3/S4/S6; handoff 6 | PGF8, explicit remote-Steering disposition | planned; no merge/release authority inferred |

### PGF1 — Final work disposition

Implement the guarded path from a verified/published result to completed routine
work. Resolve the fact that publication changes canonical input bytes: an original
review becoming historical must not block legitimate closure or authorize unrelated
new work. Keep candidate/test/review/publication identity pinned through that path.
Update only the milestones, cards, plans and Session steps actually satisfied;
retain broader work as active with named remaining scope. Support rework/successor
work without deleting old completion evidence. Record actual owner acceptance only
where required and actually supplied.

Exit: a small real assignment reaches a consistent terminal disposition across
runtime and canonical records. Missing evidence, changed candidate, wrong authority
and partial publication refuse. Interrupted closure/retry does not duplicate events
or leave a falsely completed projection. A completed subtask cannot close a broader
card/Session. Preserve original journal and ledger bytes.

### PGF2 — Leader-directed execution and continuation

Connect leader proposals to bounded controller-owned prepare/start/continue/stop
operations with scope and authority checks. Existing coordinator read/acknowledgement
capability does not confer those privileges. Carry reports, review results, questions
and owner steering through one durable coordination flow. Advance authorized phases
without unnecessary repeated permission requests; hold genuine design/scope gates.
Release or revise completed/retired scope claims safely, retaining the same Session
for continued work. Keep local and cross-worktree claim guarantees explicit.

Exit: two actual disjoint Masters progress while Project Leader remains responsive;
a conflicting third assignment is refused or queued with a reason. Phase completion
becomes idle and then continues only under current authority. A delayed report or
repeat delivery cannot override newer owner steering or start work twice. Until a
stronger shared allocator exists, concurrent-work preflight remains required.

### PGF3 — SDL design handoff and replan

Finish SAD dispatch using existing SDL, ModelGovernance and Blueprint services.
Exercise NOW -> WORK/TARGET -> validated views/Blueprint -> revision-bound owner
disposition -> Master assignment. Give the implementation agent both broad system
orientation and the relevant NOW/TARGET details. Keep design revision, task revision,
owned paths and return conditions bound together. Owner feedback received directly
by SAD/Master must reach Project Leader and the same durable work records.

Exit: an actual SDL design change in pg_lab produces a usable bounded implementation
handoff. A later TARGET/task/protected-boundary change invalidates the affected
handoff before dependent continuation, preserves unaffected work, and follows an
explicit replan/review path. No Blueprint redesign, code-derived context, runtime
snapshots or new SDL profile is required by this checkpoint.

### PGF4 — Usable owner conversation and overview

Join the existing local controller/console and inspection functions into a clear
supported entry path. Present work, assigned role, Session/phase, running/idle/waiting
status, questions, next action and evidence/review/publication/disposition separately.
Show observed tools/child mappings and explicit stale/unknown values. Support direct
SAD/Master dialogue without losing the leader overview or relevant feedback. Keep
revision hashes and JSON as inspection details, not routine owner interaction.

Exit: the owner can start/resume work, ask for status, answer a question, inspect a
proposal, give feedback and see the consequence without constructing JSON or copying
messages between chats. No unselected GUI framework, Codex fork or remote hub is a
prerequisite; use the existing client unless a concrete usability gap requires a
recorded scoped design choice.

### PGF5 — Complete-workflow recovery

Test observer closure, controller/process loss, an actually active host turn,
unknown thread/turn creation, lost responses, queued owner input, review interruption,
partial canonical publication and closure. Recover the exact mapped work or expose
an actionable unresolved state. Preserve queued messages and distinguish gaps from
observed activity. A restart must not implicitly expand scope or replay a model send.

Exit: each recovery case has exact-candidate evidence and a documented supported
recovery action. Core crash fixtures supplement real host recovery; completed-history
replay alone does not establish active-turn recovery. Unrelated dirty files, earlier
history and other assignments remain intact. Store backup/restore requirements are
explicit; Git documents alone are not claimed to restore an in-flight runtime.

### PGF6 — Sustained operation and compatible hosts

Resolve finite-store limits before claiming everyday unattended operation. Current
examples include 16 leader dispatches, 64 queued owner messages, 256 leader receipts,
1024 general receipts and a 16 MiB store. Select a bounded archive/rotation/compaction
policy that preserves retry identity, historical evidence and pending work; raising
limits alone does not prove long-term operation. Verify restart around that boundary.

Keep stable Codex 0.164.0 as the minimum and compatible later releases capability-
checked. Test the minimum, available selected later host versions, unsupported host
responses and strict same-attempt version recovery. Do not promise compatibility
with every unknown future host or change global host settings as part of a test.

Exit: define and record a representative multi-hour workload and event volume before
running it; demonstrate that it crosses current short-trial boundaries without lost
reports, exhausted queues or corruption. Use accelerated boundary/retry tests as well
as an actual sustained workflow. Record duration, load, limits and recovery results.

### PGF7 — Final native integration in pg_lab

Run the final combined candidate through a plain Maintenance task and an SDL-based
design/implementation task. Include positive completion, deliberate failed review
and repair, scope feedback, delayed owner response, two disjoint Masters, conflict
refusal and interruption/recovery. Validate both runtime and canonical work states.
Use one final acceptance matrix for PG-A01–18 and handoff scenarios 1–5, with case,
implementation/evidence source, candidate and limitations. Reuse valid lower-level
evidence; rerun scenarios whose affected candidate or contract changed.

Exit: the combined flow passes at the stated workflow level. Distinguish actual
native reviewer/worker observations from fixtures and attributed external tests.
There is no Ponsse MVP1 pilot in this checkpoint; pg_lab remains the authorized lab.

### PGF8 — Independent review and owner trial

Review PG itself in an implementation-independent context, following the governing
intent, design and final acceptance matrix. Testing that PG can launch a reviewer
is not an independent review of PG's implementation. Arrange that context under
then-applicable delegation authority; keep the gate open if independence is absent.
Fix blocking findings and repeat affected verification/review on the changed basis.

Prepare a concise owner walkthrough of normal operation, a question/decision, a
rejected or stale design, recovery, and the final work disposition. Capture actual
owner observations and consequential corrections. The reported earlier development
trial is retained but does not substitute for feedback on the completed workflow.

Exit: required independent findings are disposed with evidence and the owner can
assess a concrete working result. Do not invent owner feedback or acceptance.

### PGF9 — Scope disposition and development closeout

Reconcile PGI/PGL milestones, Session S3/S4/S6 and KB038/052 with the final evidence.
Broader KB036/037/042/043 work remains explicitly listed, not implicitly delivered.
Update maintained current instructions/walkthrough and remove contradictory current
status statements while retaining historical evidence. Prepare a reproducible local
build/start/recovery handoff and the combined PR under existing authority.

The selected handoff still includes revision-bound remote Steering with honest
committed-versus-dirty visibility (scenario 6). The separate remote hub project does
not automatically cancel this requirement. Demonstrate the selected online host
integration, or obtain and record an explicit owner disposition to transfer/defer it
before claiming that full handoff complete. This decision can be prepared while the
local implementation proceeds; it is not a reason to block unrelated PGF1–8 work.

Exit: every selected requirement is delivered/verified or explicitly dispositioned;
no active card is kept merely as a pointer to a successor, and no unfinished broad
capability is closed. Distinguish development-ready local PG from merged/released
or installed PG. Merge, release and consuming-project installation remain separately
scoped actions. The existing Ponsse model and deferred MVP1 pilot stay untouched.

### Execution discipline and remaining decisions

Use one accountable implementation agent and sequential checkpoints. Native agents
inside pg_lab are product tests, not parallel implementation owners. Independent
review is a distinct role/context. Before each implementation checkpoint refresh
concurrent-work preflight; at delivery update this plan, Session and evidence,
commit the bounded outcome and push under existing authority. New gaps must map to
an existing acceptance requirement or be registered separately before expanding scope.

The remaining design choices include closure versus successor-work semantics,
trusted leader operations, supported owner-dialogue surface and store rollover.
Resolve them from existing authority at their checkpoint; do not silently reduce
scope or ask the owner to approve routine implementation choices. Actual owner
design/acceptance input and the remote-Steering disposition cannot be fabricated.
No calendar promise or percentage is inferred from current component tests.
