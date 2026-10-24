# Project Leader and SAD orchestration — ProjectGovernance handoff

| Field | Value |
| --- | --- |
| id | KB-SDP-052 |
| project | SDP |
| type | Change |
| CardState | in-progress |
| created | 2026-10-09T12:08:50.581137+00:00 |
| source | Owner Session0011 T004 role definition and T005 KanBan handoff request |
| Systems | SDPTOOL |
| next_review | C2c2 owner decisions, failed-result rework, dependencies/replan and overall closure, then PGF2/3 leader/model gates |
| PlanId | PLAN-SDP-0018 |
| tags | project-leader, steering, SAD, master, MCP, app-server, coordination |

## Outcome and owning workstream

Make Project Leader the normal owner-facing coordination agent, able to supervise
multiple bounded SAD/Master assignments while remaining available to the owner.
The owner can also talk directly with SAD/Master; consequential feedback must reach
shared records and the leader without requiring manual chat-to-chat copying.

This is a distinct runtime/UI increment related to
[KB-SDP-038](%23038--Proposal--Codex-app-server-development-client.md).
Its owning workstream is Session0007 / ProjectGovernance. The registration branch's KB038 copy was backlog, while the owning worktree has
it active. This imported card now belongs to that active workstream; its original
handoff does not reset KB038 or duplicate the controller implementation.

## Authoritative handoff and delivered foundation

Read [Project-Leader-Handoff.md](../../04--Design/SDPTool/ProjectGovernance/Project-Leader-Handoff.md)
for the full owner direction, source observations, required increments and acceptance
scenarios. Keep runtime design details there and in the owning design/plan, rather
than duplicating them in this card.

[MAINT-SDP-0016](../../Maintenance/PLR1/Plan.md) delivered local skills and the shared
[role contract](https://github.com/Hans-Einar/SDP/blob/edfac185bb11b06ae9c6dd39a1809d166ad9541f/Skills/sdp/references/roles.md): new Project Leader, strategic
Steering, Architect as SAD, and Master as the owner of a bounded implementation
Session/ImplementationPlan. Metadata/distribution checks pass; this does not deliver
runtime orchestration, online installation or automatic notification.

## T039 owner requirement — mandatory governance supervision

Session0007 T039 explicitly requires card-derived Session plans before agent
assignment, new card/Session/plan preparation for logically separate goals,
continuous process checks, a recorded Blueprint applicability decision, mandatory
SDL WORK for every design-affecting change, model release tracking and coordinated
model/plan/code integration. This refines this card's existing goal; it does not
start a parallel Session or select the entire KB036/037 backlog. Project Leader
must load and apply its canonical skill on every coordination turn. Runtime checks
and complete native evidence remain required, separately from skill publication.

## Required scope

- Project Leader prepares bounded SAD assignments, then implementation assignments
  pinned to reviewed NOW/TARGET/task/blueprint revisions and linked cards.
- SAD creates model WORK/TARGET with ModelGovernance and generates previews/blueprints;
  owner design feedback and approval bind exact revisions. Changed design invalidates
  the affected handoff rather than silently widening implementation scope.
- Masters own Sessions and phased plans, report phase evidence, and can remain idle
  for authorized continuation. Preserve direct owner access and leader visibility.
- Extend privileged controller/MCP capabilities deliberately; current worker MCP
  does not expose launch/resume rights. Reuse actual runtime services and preserve
  per-assignment authority, optimistic revision checks and retry safety.
- Durable event delivery and a responsive leader conversation: correlate reports,
  retain them while the leader is idle/busy, and process them without duplicate
  launch or losing newer owner steering. MCP alone is not a wake-up scheduler.
- A TUI/client overview exposes agent hierarchy, role, Session/phase, observed state,
  tools and child activity with freshness/unknowns. Select a supported UI approach;
  a Codex TUI plugin/fork is not assumed available or selected.
- Reconcile shared-work claims and ID allocation across worktrees. Use the existing
  preflight until stronger coordination is implemented; do not claim global exclusion
  from branch-local records. Package/test Steering in the selected online host as a
  separately verified integration, reusing canonical skill sources.

## Acceptance and next action

Before implementation, the ProjectGovernance agent should reconcile the handoff with
its actual current design and PLAN-SDP-0018, then select bounded milestones or a linked
successor plan. Exercise all six scenarios in the handoff: parallel disjoint work,
shared-scope refusal, changed design, phase/idle continuation, durable notification,
truthful dashboard and remote evidence boundaries. Host completion must not imply
assignment acceptance. No main merge, release or native XFMD changes are authorized
by this registration.

## Worklog

2026-10-09: Owner requests a concrete KanBan entry for the handoff. Registered as
KB-SDP-052 after inspecting IDs in all nine known local worktrees. Session0011
records this transfer; its KB051 product implementation remains behind the red
preflight. The existing duplicate KB051 is recorded separately and not renumbered.

2026-10-10: Session0007 T014 owner explicitly selects KB052 after trying the lab.
Imported from blueprint commit edfac185bb11b06ae9c6dd39a1809d166ad9541f; retain original registration history.
[Orchestration design](../../04--Design/SDPTool/ProjectGovernance/Orchestration.md)
and [PLAN-SDP-0018](../../05--Implementation/SDPTool/ProjectGovernance/Plan.md)
now own the staged delivery. PGL1 is the first selected runtime increment.

2026-10-10: PGL1 implementation starts within the scoped preflight boundary.

2026-10-10: PGL1 delivers registered-work overview in readable and JSON forms,
verified on existing owner lab state without mutation. See [evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL1.md).
PGL2 assignment preparation/revision guards are next. The card stays in-progress:
Project Leader dialogue, multi-assignment execution, durable inbox, SAD integration,
phase continuation, full dashboard and online Steering evidence remain outstanding.

2026-10-10 / T015: Owner continues. PGL2-C1 delivers immutable SAD/Master preparation,
exact source/generation checks, local scope refusal, retry recovery and stale
inspection. [Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL2-C1.md)
is CLI/service-level with fixture artifacts. PGL2 remains in-progress for actual
disposition/provider/revision/retirement gates; PGL3–5 and owner acceptance remain.

2026-10-18 / T018: PGL2-C2a delivers stable-ID package revision and terminal
retirement with preserved receipts and local scope release.
[Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL2-C2a.md)
records scoped tests and the separate compiled CLI trial. CardState remains
in-progress; actual disposition/provider checks and PGL3–5 remain unfinished.

2026-10-18 / T020: PGL2-C2b delivers read-only checks against canonical Blueprint
and live ModelGovernance captures, exact task/scope and preparation revision.
[Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL2-C2b.md)
records source-drift refusals, compiled CLI/schema checks and the isolated lab preview.
No owner acceptance is inferred. CardState stays in-progress for trusted disposition,
dispatch and PGL3–5. EVT-KB-SDP-000372 links this review to the active plan delivery.

2026-10-18 / T022: Owner authorizes continuing while reading the walkthrough.
PGL2-C2c delivers the local controller decision adapter with immutable revision-bound
records, supersession and current/stale applicability. [Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL2-C2c.md)
uses scripted fixture decisions, not actual owner approval. Final conversation/UI and
same-user identity authentication are not claimed. CardState remains in-progress;
dispatch prerequisites and PGL3–5 remain. EVT-KB-SDP-000373 records this update.

2026-10-19 / T024: Owner requests continued work toward Session completion.
[C2d evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL2-C2d.md)
delivers exact approved Master activation, live run authority and a real isolated Go
implementation trial. External behavioral tests pass and protected files remain intact.
Native report submission was refused under never; no complete reporting, acceptance
or independent review is claimed. Scope closure/replanning, SAD dispatch and PGL3–5
remain open. CardState stays in-progress. EVT-KB-SDP-000374 records this boundary.

2026-10-19 / continued T024: [PGL3-C1](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL3-C1.md)
delivers per-run controller isolation with legacy exclusion and process-exit recovery.
Tests demonstrate independent progress with protocol hosts, not native parallel agents
or the responsive leader. Report interaction preference is pending; host policy is
unchanged. CardState remains in-progress. EVT-KB-SDP-000375 records this milestone.

2026-10-19 / T025: Owner chooses automatic reporting and sustained execution within
authorized assignments, with durable attention requests when input is needed. This
resolves the reporting preference from T024; it does not bypass a host refusal or
accept implementation. Local report/inbox/attention contracts continue under this
card. [KB054](../backlog/%23054--Proposal--PG-attention-and-remote-coordination.md)
registers the proposed self-hosted multi-machine hub and web/PWA client separately.
See the [architecture proposal](../../03--Architecture/ProjectGovernance/Attention-and-Remote-Coordination.md).
CardState stays in-progress; runtime completion/evidence remain unfulfilled.


2026-10-19 / T026: Owner separates hub/webapp into another project and continues local
SDP/SDPTool. [PGL3-C2 evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL3-C2.md)
records explicit automatic reports and the durable read-only inbox. Final native
code/report trial passes; first incorrect host-policy mapping is retained as failed
evidence. Run remains verifying, not accepted. Next: typed attention/reply, leader
consumption/continuation and native parallel evidence; receipt capacity/recovery must
support the later multi-hour claim. CardState remains in-progress. KB378 records this
bounded delivery; KB379 records the separate-project handoff refinement in KB054.


2026-10-21 / T027: Owner continues local PG. [PGL3-C3 evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL3-C3.md)
delivers durable questions, revision-bound attributed replies, withdrawal and per-run
advancement guards. Agent MCP and compiled CLI evidence pass; the separate lab has
an open question for owner experimentation. No native attention continuation or
actual owner acceptance is claimed. Next: leader consumption, guarded continuation,
capacity/recovery and native parallel evidence. KB380 records this review; CardState
remains in-progress. Backlog/onHold review retains KB036/037/042/043 dependencies
and KB054's separate-project boundary; no new card activation or scope expansion.


2026-10-21 / T028: [PGL3-C4 evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL3-C4.md)
delivers guarded current-reply continuation through a fresh attempt in the same run,
with original policies/skills, durable dispatch claims and a bounded opt-in watcher.
The native ask/reply/code/report trial and duplicate/stale/uncertain guards pass.
A report remains verifying; no owner acceptance or full leader is claimed. PM207 /
KB381 record the milestone. CardState remains in-progress. Next: report consumption,
responsive coordination, native parallel evidence and capacity/recovery; SAD/PGL4/5
and retained PGI gates remain. Backlog/onHold review retains KB036/037/042/043 and the
KB054 separate-project boundary without activating additional work.

## Session0013 consolidation/MCP handoff — 2026-10-23T21:58:09.203458+00:00

WR2 verifies the consolidated source with passing Go race/vet checks and actual
read-only Codex 0.164.0 MCP calls. Current sdp-controller refuses that host because
it requires 0.161.0. Worker MCP still has restricted seven-tool scope; no leader
control API or report-consumption loop is inferred.
[MCP adoption study](../../Maintenance/WorkstationReadiness/MCP-Adoption-Study.md)
recommends a measured current-host compatibility increment, then trusted leader/
report and evidence/publication gates. This is a readiness handoff, not new PG
implementation, phase selection, native/independent-review acceptance or closure.
CardState remains in-progress; existing PLAN-SDP-0018 authority is preserved.
EVT-KB-SDP-000387. MVP1 is deferred; the existing Ponsse pilot model is preserved.

## T029 selected host compatibility — 2026-10-23T22:10:45.708406+00:00

Owner selects [KB057](../completed/%23057--Change--Codex-forward-compatibility.md) / PGC1 under
PLAN-SDP-0018, minimum 0.164.0 and compatible later hosts. This replaces exact
release admission while retaining observed identity, live wire checks and recovery
boundaries. PGL3 leader/report work remains open; CardState stays in-progress.
EVT-KB-SDP-000390.

2026-10-23T22:33:53.920551+00:00: PGC1 source 98b449c delivers stable >=0.164.0 admission, actual-version
provenance and guarded legacy/recovery behavior. Full affected race/vet, 60 wire
objects and real no-model 0.164.0 probe pass. Earlier exact-version blocker is
resolved at this boundary. Current-host full native task/review and leader/report/
evidence gates remain. CardState stays in-progress; EVT-KB-SDP-000397.


## T031 — PGL3-C5 durable report acknowledgement

Delivered trusted-controller acknowledgement of exact reports with bounded
attribution/note, durable retry and current-source/generation checks. Acknowledgement
leaves runs unchanged and does not verify/accept work. Full affected race/vet,
compiled CLI/restart/old-reader and wire schema checks pass; see
[evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL3-C5.md).
CardState remains in-progress for responsive coordination and the retained plan
gates. No native task, independent review, Blueprint pilot resumption or Ponsse work.
SDL-only direction and KB059 onHold remain. EVT-KB-SDP-000400 records this result.


## T032 — PGL3-C6 shared coordination entry point

Delivered combined registered-work/report/attention view with bounded paging and
stale-cursor checks, explicit coordinator MCP read/ack authority, and terminal watch
with owner-attributed replies while independent service work progresses. Full race/
vet, actual compiled CLI and stdio MCP/restart checks and ten final wire objects
pass; see [evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL3-C6.md).
CardState remains in-progress for AI leader scheduling and the parent plan's native,
review/recovery/publication gates. No model/subagent launch, stable/global config
change, Blueprint resumption or Ponsse work. EVT-KB-SDP-000401 records this boundary.


2026-10-24, Session0007 T033: PGL3-C7 delivers durable owner-message/notice dispatch,
serialized short leader turns and exact mapped lost-response recovery. Full affected
race/vet, compiled CLI/protocol restart and 20 output/context schemas pass; previous
C6 binary refuses profile 0.9. See [evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL3-C7.md).
CardState remains in-progress: native leader/worker workflow, SAD/scope closure,
capacity/compaction, independent review and canonical publication remain. No native
model or subagent launch in this increment, and no full PG acceptance. SDL-only,
Blueprint pause and Ponsse deferral remain. EVT-KB-SDP-000402 records this review.


T034 native trial findings: inherited Apps inventory exceeds the bounded host frame,
and an unsubmitted bound thread is not loaded after controller restart. The test
uses process-local tool isolation; correct unsent-bound restart by safe supersession
under C7, preserving input and the no-resend guard for sending/running states.
Native evidence and final disposition follow in PGL3-C8.


2026-10-24, T034 / EVT-KB-SDP-000403: [PGL3-C8 native evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL3-C8.md)
passes in owner-authorized pg_lab: three leader turns, one worker, owner enqueue,
report receipt and two exact recovery cases. Corrected safely unsent bound-thread
restart; affected race/vet and 11 native schema checks pass. CardState remains
in-progress for retained parallel/active-host, SAD/scope, capacity/review/publication
gates. No owner acceptance, independent review, merge/release or Ponsse work.


2026-10-24, T035 / EVT-KB-SDP-000405: [PGI4-C1 evidence and rework](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGI4-C1.md)
implements exact candidate capture, attributed artifact receipts, current/historical
applicability and controller-only failed-evidence rework. Final full SDPTool race/vet,
30 wire schemas and the compiled pg_lab pass/fault/fail/rework/restore trial pass.
The original C8 native state and all 43 baseline regular files remain preserved.
Same Session continues; old agent capability is revoked on rework. No attachment
establishes independent review, and completion/publication remain unavailable.
CardState stays in-progress. Next: retained review/publication contract and evidence,
with SAD/owner, parallel/active-host/capacity and integrated verification still open.
Backlog/onHold reviewed; no Blueprint restart, Ponsse, hub, merge or release selected.


2026-10-24, T036 / EVT-KB-SDP-000407: [PGI4-C2 review evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGI4-C2.md)
delivers observed fresh-context read-only review of exact candidate/test evidence,
structured results, current qualification and strict original-turn recovery. Final
full race/vet, 24 schemas and actual gpt-6-astra/high positive/negative lab reviews
pass. The reviewer identified the deliberate separator defect; its finding returned
the same test Session to ready through C1 rework. Original C8 state and 51 baseline
regular files remain preserved. No independent review of the PG implementation,
owner acceptance or native child hierarchy is inferred. CardState remains in-progress.
Next: checked canonical publication/recovery and review-result leader delivery;
retain SAD/owner, parallel/active-host/capacity and independent integrated gates.
Backlog/onHold reviewed without changing SDL-only, Blueprint, Ponsse or hub scope.


### T037 — PGI4-C3 checked result publication and coordination

Delivered bounded canonical publication with durable intent and exact crash/retry/
conflict handling. Review/publication facts are exposed through existing coordinator
MCP/CLI and serialized leader notices. Full final race/vet, 13 wire objects, compiled
lab publication, canonical validation and actual MCP exchange pass. The new review
fixture is synthetic; earlier native reviews are inspected read-only.
[Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGI4-C3.md) and
[Session0007 T037](../../Sessions/session-%230007--Project_governance.md) retain limits.
CardState remains in-progress: final run/lifecycle disposition, retained orchestration
and independent integrated verification still require delivery. No broader closeout.


### T038 — Completion plan consolidated

Owner requests a detailed completion plan. PLAN-SDP-0018 now maps remaining PGI/PGL
work to PGF1–9 with dependencies and exit conditions: final disposition, leader
execution, SDL handoff, owner dialogue, recovery/capacity, native integration,
independent review and scope closeout. Remote Steering needs verification or explicit
owner disposition; broader deferred capabilities are not silently included or closed.
CardState remains in-progress. This is planning, not a new implementation delivery.


### T039 — Mandatory PG supervision

Owner makes card-derived Session/plan admission, separate-goal routing, mandatory
SDL WORK for design changes, reasoned Blueprint applicability, model release and
model/code/plan conflict coordination explicit. Recorded in the existing KB052
scope, Orchestration contract, PG-A19–24 and PLAN-SDP-0018 PGF checkpoints.
Canonical Project Leader skill 1.1.0 and development distribution inventory now
state the detailed mandatory routine. This is requirement/skill maintenance;
runtime guards and native compliance evidence remain open. CardState remains
in-progress; next is PGF1 with refreshed preflight and required SDL design preparation.


### T040 — Session execution and KanBan governance discussion

Owner asks for programmatic card maintenance/transitions, resumable Session workflow
instances and common Project Leader/MCP plus owner/CLI access and activity insight.
[Current evidence and proposed boundary](../../04--Design/SDPTool/ProjectGovernance/Orchestration.md#t040--discussion-executable-session-workflows-and-kanban-governance)
compare existing run state/publication with missing full Session and card lifecycle
execution. Recommend extending the shared core with checked KanBan effects and
versioned Session bindings, retaining authored plans, canonical management history
and role authority. Exact format and operation semantics remain discussion proposals.
Current PG integration remains active; no new runtime implementation was selected.


### T041 — Selected Session/KanBan workflow concretization

Owner accepts the T040 direction and selects its bounded design. The
[Session execution contract](../../04--Design/SDPTool/ProjectGovernance/Session-Execution.md)
defines instance/binding ownership, stage transitions, mandatory card effects,
coverage/replan, shared CLI/MCP authority and recovery. SDL WORK and diagnostic
Blueprint validation are recorded in [evidence](../../04--Design/SDPTool/ProjectGovernance/Evidence-T041.md).
PG-A25–30 and PLAN-SDP-0018 PGF1-C1 onward define the implementation sequence.
This is design preparation; runtime acceptance and model release remain open.
Existing CardState is unchanged; broader KB037 scope is not activated wholesale.

### T042 — Durable Session binding and registration delivered

PGF1-C1 now validates an explicit card-derived Session binding, persists one pinned
instance, exposes stages/holds and recovers the exact registration journal through
CLI/JSON actions. Full SDPTool race/vet, strict wire captures, crash/stale/authority
checks and actual-process pg_lab trial pass. All 110 earlier lab paths are preserved.
[Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C1.md)
records source identity and limits. The aggregate model target and independent
product review/release obligations remain open. No host dispatch or lifecycle
completion is implied. CardState stays in-progress for the substantial remaining
outcome; next is PGF1-C2 actual card effects and scoped reconciliation.


### T043 — Checked ready-card selection and recovery

PGF1-C2a delivers real queued/backlog-to-active/ready card moves and planned-plan
activation with journals, current links/index and append-only history. Prepared
record effects are applied by resume; partial moves recover without duplicate
transitions. Current cooperating PG/history/Blueprint writers honor a durable
reservation, and selected work claims reject conflicting Session scope. Valid
unrelated history is reconciled after commit; selected changes hold continuation.
[Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2a.md)
records final candidate checks and the isolated pg_lab trial. CardState remains
in-progress: start/progress, verified disposition, common leader/model gates and
independent review remain. Next is C2b; selection alone launches no agent and
completes no card or Session.


2026-10-24 / Session0007 T044: PGF1-C2b1 adds explicit design/Blueprint assessment,
required procedure/evidence snapshots, conservative routing, immutable assessment
records and checked canonical recovery. Unknown design holds admission; explicit
model scope requires design. CardState remains in-progress for the broader outcome.
[Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2b1.md)
records the exact candidate, actual CLI/lab journey and final checks. C2b2 recorded
start/progress, C2c verified disposition and privileged leader/model gates remain.
No native skill-consumption, independent review, owner acceptance or release claim.
EVT-KB-SDP-000424 records this unchanged-lifecycle work review.


2026-10-24 / Session0007 T045: PGF1-C2b2 delivers no-design Session start/progress
through the existing run/controller. Canonical in-progress records precede one ready
run; assigned report/current candidate evidence permits verification-stage progress.
All selected cards/contributions and exact required skill snapshots remain bound.
[Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2b2.md)
records full race/vet, fault/tamper/multi-card tests and the actual CLI/controller
protocol trial in pg_lab. No native Codex/model execution or independent PG review.
CardState remains in-progress: next C2c owns verified disposition, rework/dependencies
and claim release/replan; privileged leader/model gates and release remain open.
EVT-KB-SDP-000426 records this unchanged-lifecycle work review.


2026-10-25 / Session0007 T046: C2c1 binds exact passing test/review evidence to
Session stages and canonical records. Review follows committed verification and
retains all selected whole/contribution context. Current candidate/test/skill checks
and recovery prevent stale intake. Cards remain in-progress; no owner acceptance,
card completion or claim release is inferred. The compiled two-card pg_lab fixture
passes with synthetic report/review observations, not native execution or independent
PG review. [Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2c1.md)
records final candidate and checks. C2c2 owns actual disposition, owner decisions,
rework/dependencies and claim release; later leader/model/release gates remain open.
EVT-KB-SDP-000428 records this unchanged-lifecycle work review.


## T047 — Checked selected-step disposition

2026-10-25 / Session0007 T047: C2c2a completes an eligible no-design step after exact
current verification/review where owner acceptance is explicitly not required.
Whole-card completion moves records/metadata/selected links/index; contributions stay
active with explicit remaining work. Canonical disposition and management/system
evidence precede claim release and Run completion. Plan/Session closure is separate.
Required owner decisions, external edits and overlapping scope continue to hold work.

Final race/vet, every-effect recovery, malformed-state/currentness/authority checks
and compiled pg_lab evidence pass. All 264 earlier lab paths are preserved; the
final executable replays completed receipts without changing state or fixture bytes.
Synthetic report/test/review observations do not prove native execution or independent
PG review. [Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2c2a.md)
records exact candidate identities, corrections and limits.

This card remains in-progress for owner decisions, failed-result rework, dependency
progression/replan, overall closure and privileged leader/model gates. Existing WORK
:00002 remains TARGET, not accepted/released. Backlog/onHold review preserves the
related deferred cards; KB050/Session0008 is separately paused and MVP1 is deferred.
