# SDP status and consolidation study — 23 October 2026

Current follow-up — Session0013 T002: the authorized WR2 round now has a private
Go 1.27.1 build and passing integrated Go/race/vet checks. Local read-only MCP use
is demonstrated through Codex 0.164.0; the execution controller remains pinned to
0.161.0 and refuses that host. See [execution evidence](Execution-Evidence.md) and
[the MCP adoption study](MCP-Adoption-Study.md). Original findings below retain
their initial observation boundary; F01/F02 now have a working isolated development
candidate rather than a stable installation/release. MVP1 is deferred and its
existing Ponsse pilot model is preserved.

## Conclusion and evidence boundary

The available evidence supports consolidation around one accountable completion
agent and a bounded Ponsse-readiness target. It does **not** establish that SDP is
broken, that ModelGovernance needs reimplementation, or that parallel work must
be discarded. Existing history reconciliation passes its checks. The immediate
practical gap is between the integrated source, the selected installed engine and
the workstation build environment. ProjectGovernance also has explicitly unfinished
workflow and verification gates.

The owner confirms that all remaining laptop work was handed over. This supersedes
the laptop-recovery uncertainty in Session0012 T002. No further laptop inventory is
required. Existing active records describe unfinished scope, not proof that old
agents are still executing.

This is the delivered status study, not an implementation or release acceptance.
[Snapshot](evidence/snapshot-2026-10-23.json) binds the inspection to source hashes,
SDP main bfe3b40, assessment HEAD 37836ac and XFMD 953d09e. Product files were
unchanged during the checks. Later documentation registration changes the management
inventory counts. [Session0013](../../Sessions/session-%230013--SDP_consolidation_and_Ponsse_readiness.md)
owns the proposed route; [Plan](Plan.md) owns this study/planning assignment.

## What has actually been organized

The shared KanBan is the intent/work-state entrypoint. Typed plans own bounded
execution and milestone evidence. Sessions carry goals and handoffs across turns;
ProjectManagement owns lifecycle history, and Traceability owns implementation/design
evidence. These are different responsibilities, not five independent competing plans.

Parallel branches carried PG, blueprints/actions and XFMD producer changes. The
October 18 Session0011 import reconciled conflicting IDs while preserving exact
source histories. October 23 commits consolidated workstation WIP and retained
superseded branch history. Preserved branches/commits are provenance, not evidence
that duplicate executable implementations are active.

The October 10 interim coordination overview was useful at its date but is stale
as an allocation guide. Its next-PGL2 and Actions-blocked descriptions have been
superseded by owning Sessions. The old absolute worktree paths are recovery history,
not required workstation layout.

## Current capability and work inventory

| Area | Evidence and observed state | Remaining work / consequence |
| --- | --- | --- |
| ModelGovernance | Session0006 goal achieved; KB049 completed; production owner is SDPTool/model | Preserve implementation. Rerun meaningful model lifecycle tests on the integrated workstation candidate before relying on it in Ponsse. |
| Semantic blueprints | Session0008: BPI1–3 and PLAN-SDP-0020 completed; Blueprint generation/catalogue/assessment/assignment source present | KB050 gate-review and owner-paused pilot remain. Reduced MVP1 specimen is not proof that the actual Ponsse workflow is ready. |
| Actions and JSON | Session0011 and KB053 completed; SDPTool/Actions.md describes production source interface | Source completion is distinct from installed availability. Current gh sdp engine rejects actions; native XFMD consumer remains a separate held scope. |
| ProjectGovernance | Session0007 T028 and PLAN-SDP-0018 active; KB038/052 in-progress | Report consumption, responsive leader, capacity/recovery, native parallel evidence, later SAD/PGL and retained PGI gates remain. No claim of full automatic coordination. |
| SDL/SDUI and Toolkit | Existing implementation, widget/program and distribution records; Toolkit Python checks pass here | No fresh Go/race/native GUI evidence on the consolidated candidate. Preserve compatibility and verify selected workflows. |
| Viewpoint producer and XFMD | SDP KB032 includes later N5–N15 handoffs. XFMD Session0007 T012 records VP02–VP11, deferred discovery and local XFMD 0.9/private producer installation | Retained acceptance, integration and resource guards remain under existing cards. Session0005 alone understates this progress. Installation is recorded prior evidence, not freshly reproduced here. |
| Dynamic XFMD actions | XFMD Session0006 / KB030 paused pending released producer API | Preserve hold; producer source delivery does not satisfy release/installed consumer dependency. |
| Broad navigation/document work | KB017 retains broader scope; KB010 gate-review is an old concrete editorial review | Do not treat all open cards as new defects or mandatory Ponsse blockers. Give retained scope an explicit disposition when selected. |

Before this assignment, management replay found 68 cards, 45 management records,
43 completed and two active management records. The two active records are the
blueprint DesignPlan and PG ImplementationPlan; a paused Session can retain an
active plan because the plan schema has no paused state. That is not inherently
an inconsistent lifecycle. New registration for this study changes these counts.

## Findings and proposed treatment

| ID | Finding / classification | Evidence | Proposed change and completion criterion |
| --- | --- | --- | --- |
| F01 | Selected runtime lags source — confirmed readiness gap | gh sdp --version --json returns 2.2.0 / ef8741c; gh sdp SDP actions --json exits 1 with arguments error | Establish intended engine selection and build a pinned current candidate. Demonstrate version and action catalogue through the selected executable; preserve the released fallback. Publication/installation authority is separate. |
| F02 | Workstation build prerequisite unavailable in this execution environment | go is absent from PATH and checked conventional locations; SDPTool declares Go 1.26.0, CI selects 1.27.1 | Locate or provision a compatible toolchain through workstation conventions; record actual version and run affected suites. This does not assert Go is absent everywhere on the machine. |
| F03 | Aggregate completion evidence missing — confirmed evidence gap | Consolidation commits combine previously separately tested work; this turn runs Python checks only | Verify integrated model→blueprint→assignment→implementation/report→review/disposition workflow, including stale source and interruption behavior. Select manual or automated execution explicitly. |
| F04 | Competing-looking assignment states — boundary risk, not proven duplicate implementation | Blueprint lifecycle in blueprintstate/lifecycle.go; prepared assignments/runs in governance; design_check.go reads canonical Blueprint and live model providers | Document a small identity/state mapping and test transitions across the boundary. Keep model artifacts, Blueprint obligations and execution attempts distinct. Consolidate code only if an actual conflicting authority or duplicated rule is demonstrated. |
| F05 | Stale coordination projections and former branch paths — confirmed documentation debt | October 10 overview versus PG T028 and XFMD T012; old worktree coordinates in handoffs | Update current entrypoints and mark historical coordinates explicitly. Preserve historical journals and hashes. Acceptance: a replacement agent can find current status and next authorized step without old worktrees. |
| F06 | Laptop uncertainty in prior assessment — corrected by owner | Current owner statement: all remaining work handed over | Remove it as a planning dependency. Do not repeat recovery requests or infer missing code from historical branch names. |
| F07 | Ponsse readiness not yet measured | No Ponsse checkout in the observed /home/warloc/git inventory; actual installed profile/model state not inspected | Recover actual Ponsse repository through available repository access when execution is selected. Inspect adoption constraints and choose one real MVP1 task before migration or application changes. |
| F08 | Scope can expand indefinitely if every backlog idea gates adoption — planning risk | PG broader cards, full timelines, remote hub handoff and native dynamic menu have distinct scopes | Adopt explicit must-have / later dispositions per workflow. Keep deferred work with its owning cards; do not silently cancel accepted requirements. |

### Duplicate-implementation investigation

A byte-identical scan of production Go files in SDPTool, SDL/go, SDUI/go and
Toolkit found no duplicate groups, excluding tests and testdata/vendor/fixtures.
That is a narrow static observation, not proof of semantic uniqueness. Package
inspection identifies clear distinct owners: model manages model artifacts;
blueprints generates and retains diagnostic documents; blueprintstate validates
assignment history; governance prepares and runs work; governancehost/controller/MCP
adapt the same core. governance/design_check.go explicitly uses model snapshots and
Blueprint assignment views. This supports integration across boundaries rather than
an immediate rewrite into one state machine.

The remaining review target is whether acceptance, readiness and completion are
mapped consistently between Blueprint assignment and PG execution. A run reporting
success must not promote model state or establish owner acceptance. A stale design
must invalidate dispatch at the correct boundary. Existing fixtures cover pieces;
a full integrated workflow and independent review are still needed.

## Checks executed on this workstation

| Check | Result | What it proves / limitation |
| --- | --- | --- |
| ProjectManagement validator | PASS: 68 cards, 45 records, 4 lineage operations, 626 events | Schema/history/current metadata/placement/grouping agree at baseline. |
| Session0011 history import verifier | PASS | Exact archives, release prefix, complete event representation and mappings remain valid. |
| ProjectManagement unittest suite | 6 tests pass | Existing management validation regression coverage. |
| Toolkit validator, base-ref main | PASS | Local Toolkit contracts against the local main baseline; no remote publication claim. |
| Toolkit unittest suite | 85 tests pass | Existing Python Toolkit coverage, not Go/runtime or native rendering acceptance. |
| gh sdp selected engine / actions probe | Version 2.2.0; actions rejected | Confirms the selected engine differs from newer source capability. |
| Go / native GUI / Ponsse workflow | Not executed | No Go in PATH; no fresh native GUI or Ponsse acceptance claimed. |

No failing executed integrity check currently supports the claim that integration
lost history. No new review was performed by an independent agent, and no prior
review is relabeled as review of this combined candidate.

## Organization alternatives

| Alternative | Benefit | Cost / risk | Assessment |
| --- | --- | --- | --- |
| Resume several implementation lanes now | Parallel throughput on truly disjoint scope | More integration and shared-state coordination before the combined baseline is verified | Defer until ownership and integrated checks are clear. |
| One completion owner, sequential bounded milestones | One current route and clear responsibility across SDP/producer handoffs | Less concurrent implementation; independent review still needs separate context | Recommended for the readiness round. Role changes do not mean independent self-review. |
| Finish all PG automation before any Ponsse use | Richer future owner/agent workflow | Largest remaining scope and delayed learning from real MVP1 work | Select only if automatic leadership is an explicit adoption requirement. |
| Use manual governed execution while finishing proven blockers | Earliest real feedback using existing skills/Sessions | Owner and agent maintain journals/handoffs manually | Recommended pilot option, subject to owner disposition; does not cancel PG commitments. |

## Proposed execution sequence and readiness decision

1. **C1 — Reproducible candidate and environment:** pin integrated source and engines,
   resolve F01/F02, run the three Go modules and bootstrap checks appropriate to
   changed boundaries, plus required race/vet and record validators. Record actual
   failures before selecting fixes. No unrequested bulk refactor.
2. **C2 — Consolidate authority and handoffs:** reconcile F04/F05 using current
   contracts and exact sources, update entrypoints and state mapping, disposition
   restored/open cards without closing unaccepted work. Fix confirmed defects in
   bounded existing plans, with a new plan only when their scope does not fit.
3. **C3 — Integrated readiness trial:** perform one isolated real workflow using
   ModelGovernance and Blueprint plus the chosen manual/PG execution route. Exercise
   source drift, bounded writes, interruption/resume, reporting and independent
   review. Report manual steps explicitly; no full automation claim.
4. **C4 — Ponsse MVP1 pilot:** inspect actual consuming-project state and installed
   process, recover owner pilot feedback, choose one useful bounded MVP1 change and
   verify the end-to-end workflow there. Consumer migration and domain implementation
   require their own actual scoped authority; this study does not supply it.
5. **C5 — Close readiness, retain broader roadmap:** identify accepted outcomes,
   deferred PG/XFMD/backlog work and any release/upgrade needed. Use the release
   checklist only if a release is selected. Resume MVP1 as normal work after the
   pilot provides the required evidence.

The owner review now concerns this sequence and the adoption target: whether a
manual governed pilot is sufficient initially, or automatic PG leadership is a
must-have. The study recommends the first option. Existing blueprint pause, XFMD
hold, merge/release gates and PG commitments remain unchanged until disposition.

A defensible Ponsse-ready claim requires a pinned usable engine, recoverable current
work, protected model/source behavior, reproducible blueprint context, explicit
assignment/evidence/review disposition and one successful real MVP1 workflow.
Backlog emptiness, universal diagram coverage, a remote hub and fully automatic
Session capture are not proposed acceptance criteria.

## Owner revision — Session0013 T002

The owner authorizes executing C1–C5 with C4 replaced by a study of starting
ProjectGovernance use through MCP. The MVP1 pilot is deferred; a pilot model already
exists in Ponsse. Earlier proposed pilot/new-model dependencies above are superseded
for this Session. [WR2](ExecutionPlan.md) owns execution and the revised result.
