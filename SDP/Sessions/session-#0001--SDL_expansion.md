# Session 0001 — SDL expansion

## Session roadmap

Latest recorded local turn: **T009 — capture AST and timeline proposals**.
Next SDL delivery step: **S2 — ImplementationPlan**, incorporating KB-SDL-007.

**Sequence-only Gantt mockup; not yet an event-derived timeline.** The dates below are synthetic placement slots,
not working-day estimates, deadlines or authorization. S1 design is delivered; S2 is next.
The step table below owns the manual projection. No init directive is used.
[KB-SDP-043](../KanBan/backlog/%23043--Proposal--Event-derived-session-timelines.md)
records the requested turn/seconds axis, current-position line and activity/event
markers. The chart below does not measure turns or elapsed time.

```mermaid
gantt
    title SDL expansion - sequence only
    dateFormat YYYY-MM-DD
    section Proposed route
    DONE S1 Design :done, s1, 2000-01-01, 1d
    NEXT S2 Plan :s2, after s1, 1d
    PLANNED S3 Build :s3, after s2, 1d
    PLANNED S4 Verify :s4, after s3, 1d
    PLANNED S5 Close :s5, after s4, 1d
```

| State | Step | Work / plan | Prerequisites | Authorization | Completion evidence |
| --- | --- | --- | --- | --- | --- |
| completed | S1 | P1: design System and explicit source sets; reuse completed study | KB-SDL-005 and source/projection constraints | Owner selected S1 | [SSD2 current contract and correction](../04--Design/SDL/SourceComposition/Plan.md); SSD1 retained as history |
| next | S2 | P2: plan runnable increments and checks | SSD2 and KB-SDL-007 refinement | Unselected | Phases/milestones, explicit branch/commit policy and acceptance |
| planned | S3 | P2: implement Go frontend/input resolution and producer consumers | S2 and execution authorization | Unselected | Cross-file identity, diagnostics, deterministic model/revision and existing behavior verified |
| planned | S4 | P3 or P2 verification: check a real model through navigation/generation | S3 candidate | Unselected | Missing/duplicate/cyclic input cases, original spans, stale revision and actual consumer evidence; required independent review |
| planned | S5 | Record goal outcome and KB020 migration disposition | S4 evidence | Unselected | KB005 disposition, explicit remaining model migrations and successor; no implied release |

**Manual roadmap pilot**, created from the owner's 2026-09-30 Session proposal.
S1 is now authorized below; this does not retroactively place prior work in a Session.
[KB-SDP-042](../KanBan/backlog/%23042--Proposal--Goal-oriented-sessions-and-roadmaps.md)
owns the provisional format; [guide](README.md) defines its limits.

| Field | Value |
| --- | --- |
| Session reference | SESSION-SDP-0001 (provisional) |
| Status | active |
| Primary delivery card | KB-SDL-005 |
| Snapshot date | 2026-09-30 |
| Current execution step | None — S1 correction delivered; awaiting S2 selection |
| Proposed next step | S2 — create the ImplementationPlan from SSD2 source composition |
| Execution authority | Owner selected S1 and source-owned correction on 2026-09-30; implementation remains later |

## Goal

Make an SDL System readable across explicit source files while retaining one
validated model, original-source diagnostics and model-derived navigation.
Demonstrate the result on a bounded real model and preserve independent model
checking. Completion requires source-set and consumer evidence, not just syntax.

Excluded: all experimental MVP1 language features, rich SDUI parity, a new FOX
adapter, XFMD application changes, release/merge without their own authorization.
The separate text-fidelity and routine/client work can continue independently.

## Affected cards

Snapshot captured 2026-09-30; lifecycle and CardState agree for these rows.

| Card | Role | Initial state | Planned final disposition | Current snapshot | Actual final disposition |
| --- | --- | --- | --- | --- | --- |
| [KB-SDL-005](../KanBan/active/%23005--SDL--Change--System-and-source-sets.md) | Primary delivery | backlog | completed after language/input/consumer acceptance | active / ready | Pending |
| [KB-SDP-020](../KanBan/backlog/%23020--Change--Shared-design-source-organization.md) | Related migration, separately selected | backlog | Determine after supported source-set design; do not promise entire migration | backlog | Pending |
| [KB-SDP-041](../KanBan/completed/%23041--Study--XFMD-driven-SDL-and-SDUI-gaps.md) | Completed input, not reopened | completed | Retain completed | completed | Already completed before this pilot |
| [KB-SDL-007](../KanBan/backlog/%23007--SDL--Proposal--Composable-file-ASTs-and-contextual-analysis.md) | Frontend refinement to review with KB005 during S2 | backlog | Select/consolidate explicit partial-analysis scope in S2; implementation not promised by capture | backlog | Pending |
| [KB-SDP-043](../KanBan/backlog/%23043--Proposal--Event-derived-session-timelines.md) | Process proposal discovered here; outside SDL delivery | backlog | Separate Session/routine/client plan | backlog | Pending |

Other language extension cards stay outside this bounded goal unless a recorded
route change adds them. A Session is not a promise to empty the backlog.

## Plan register

| Ref | Plan type and document | Document readiness | Canonical lifecycle | Role / dependency |
| --- | --- | --- | --- | --- |
| P0 | [PLAN-SDP-0009, RequirementPlan/study](../02--Requirements/XFMD-Gaps/StudyPlan.md) | completed | completed | Prior input; not performed during this Session |
| P1 | [PLAN-SDP-0010, DesignPlan](../04--Design/SDL/SourceSets/Plan.md) | completed | completed | System identity, source membership, resolution and consumer contract |
| P1R | [PLAN-SDP-0011, source composition correction](../04--Design/SDL/SourceComposition/Plan.md) | completed | completed | Replaces P1 external-manifest recommendation after owner rejection |
| P2 | ImplementationPlan — not yet created | planned | Not registered | Depends on P1R source-owned design |
| P3 | VerificationPlan — not yet created; may instead use explicit verification milestones in P2 | planned | Not registered | Candidate and consumer acceptance; avoid a redundant wrapper |

No mandatory separate ArchitecturePlan: create one only if the selected design
changes a material system boundary beyond the established architecture.

## Recorded frontend direction — T008/T009

The owner endorsed documenting the component-first model in
[KB-SDL-007](../KanBan/backlog/%23007--SDL--Proposal--Composable-file-ASTs-and-contextual-analysis.md).
Keep independent file ASTs, a source dependency graph and a context-dependent
semantic System model. A later parent/root can reuse unchanged syntax; semantic
bindings and validation must be established in its new context. Decorated AST
means semantic enrichment (for example resolved declarations and checked kinds),
which may live in separate tables keyed to syntax nodes/context revisions.

No physical reparenting/copying of shared file trees is required. Partial inspection
must distinguish parsed syntax from unresolved dependencies and fully validated
System state. Root-relative source paths still require known root context before
loading. This refines SSD2 for S2 planning; it does not claim implemented caching,
incremental compilation or native XFMD partial preview.

## Requested Session timeline — T009

[KB-SDP-043](../KanBan/backlog/%23043--Proposal--Event-derived-session-timelines.md)
preserves the detailed visual proposal. Each Step is a lane with one planned Task.
Work can occupy multiple colored segments within that lane: ordinary work,
red KanBan activity and green document activity. Turn boundaries, event cursors,
`#card`, `@document` and segment ends expose what happened without creating a new
planned task for each activity. Show a vertical current-turn/current-position line.

Generate the roadmap and table from correlated history when implemented. Use turn
numbers for overview and seconds/turn-relative offsets for detailed timing, with
clearly marked turns. Preserve unknown historical timings and distinguish active
work, waiting and wall-clock gaps. ProjectManagement remains lifecycle authority;
new correlation/telemetry schema needs design. The current journal is manual and
its Gantt remains sequence-only. This proposal does not expand the SDL delivery goal.

## Route changes and decisions

| Date / source | Previous route | Change and reason | Authority | Affected steps |
| --- | --- | --- | --- | --- |
| 2026-09-30 owner request | Next steps distributed across chat and cards | Propose one persistent goal/roadmap/turn record | Owner proposed Session concept; this is a document pilot | S1–S5 |
| 2026-09-30 pilot preparation | No time estimates selected | Use synthetic sequence slots; keep execution unselected | Agent recommendation, not an approved schedule | S1–S5 |
| 2026-09-30 next-step request | S1 proposed | Execute and deliver PLAN-SDP-0010; S2 becomes next | Owner selected S1; implementation remains subsequent | S1–S2 |
| T008/T009, 2026-09-30 | Root-first source-graph entry in SSD2 | Capture late-root AST composition and context-dependent analysis in KB-SDL-007 for S2 | Owner discussion and capture request; implementation remains unselected | S2–S4 |
| T009, 2026-09-30 | Roadmap below metadata/plans; manual sequence chart | Move roadmap/table first; register event-derived timeline in KB-SDP-043 | Owner layout/capture request; generator remains proposed | Session presentation; SDL steps unchanged |

## Turn journal

This pilot starts now. Prior chat history is not reconstructed as a verified
transcript. The completed study is linked as an input above.

### T001 — 2026-09-30, propose Sessions

- Capture mode: manual; exact current owner prompt retained in
  [the proposal card](../KanBan/backlog/%23042--Proposal--Goal-oriented-sessions-and-roadmaps.md#source-prompt).
- Host thread/turn/item IDs: unavailable to this document writer.
- Category: process proposal and planning; no language implementation.
- Procedure used: existing SDP entrypoint and Planning guidance, agent-reported.
  Formal routine-run ID/version: not available; no governance engine claimed.
- Skills loaded: sdp, sdp-planning, openai-docs; agent-reported reads.
- Work: registered KB042, drafted guide/template and this roadmap, linked existing
  routine/client work, checked current app-server documentation.
- Assistant work summary (not a captured final chat message): Session connects a
  durable goal to cards, plans, next steps, decisions and turn history. App-server
  is a suitable future client capture point; MCP supplies explicit SDP operations.
  Existing plan/card authority remains intact; rich SDUI parity stays excluded.
- Next: review/adopt the small Session format through KB042, or select S1 for SDL
  delivery. Process tooling need not block language design. Neither is executing.

## Closeout

S1 delivered under PLAN-SDP-0010 and corrected by PLAN-SDP-0011 after owner
rejection of its manifest recommendation. S2 ImplementationPlan is next; product
implementation and full Session goal acceptance remain pending.

## Pilot checks — 2026-09-30

Management validation: 52 cards, 24 management records, 398 events; passed.
Toolkit validation, local link checks and git diff --check passed. Historical
ledger bytes remain unchanged. The Gantt source was rendered with local
mermaid-rs-renderer/target/release-fast/mmdr, rasterized and visually inspected.
Short labels avoid row overlap. This renderer displayed synthetic calendar dates
rather than the requested Slot axis labels; they are explicitly not a schedule.
No native XFMD inspection, automatic projection or transcript capture is claimed.

Current session permissions again make .git read-only. Drafts are saved locally;
no commit or push was attempted. Unrelated sourceinput/Node files were preserved.

### XFMD compatibility correction — 2026-09-30

Owner reported `Unsupported Mermaid statement: axisFormat Slot %d`. Source
inspection locates that diagnostic in XFMD's Mermaid semantic parser, whose
Gantt planning profile accepts dateFormat YYYY-MM-DD but neither axisFormat nor
todayMarker. Both directives were removed from this example. The earlier
standalone mmdr rendering was not XFMD-path acceptance: it rendered without
applying the requested axis labels. No XFMD code changed or native GUI acceptance
claimed. The remaining dates still represent synthetic sequence slots.

### Git recovery and external handoff — 2026-09-30

The owner restored write access and authorized commits. The earlier read-only
note describes the pilot's initial execution, not the current permission state.
External KB-XFMD-020 now owns the Mermaid compatibility follow-up; its registration
is recorded in KB-SDP-042. This housekeeping does not start SDL design or make
the Session format an adopted management profile.

### T002–T004 — manual follow-up index, 2026-09-30

These are manually summarized work intervals, not recovered app-server IDs or
verbatim final-response capture. They refine the Session pilot; S1 remains next.

| Local turn | Owner request | Handling and outcome | Procedure / evidence |
| --- | --- | --- | --- |
| T002 | Report unsupported axisFormat in the Gantt and identify the rejecting component | Inspected XFMD semantic parser; removed axisFormat/todayMarker from the local example; standalone renderer check was insufficient | Existing source inspection and card update guidance, agent-reported; compatibility correction above |
| T003 | Register a broader XFMD card covering missing Mermaid constructs and explain parsing ownership | Prepared full compatibility draft while external writes were unavailable; documented parser/model/layout boundary and pending registration | SDP capture guidance, agent-reported; external handoff pointer |
| T004 | Restore write access, commit pending work and register XFMD card | Registered external KB-XFMD-020 and validated both boards; SDP documentation commit selected. XFMD commit handling separated because earlier modeling work shares its uncommitted ledger | SDP entrypoint and existing Git policy, agent-reported; no routine engine or automatic transcript capture |

Exact visible owner prompt for T004:

> ok you are back in yolo mode so you can do git commits and also add the kb cards in xfmd

No step is marked completed merely because this housekeeping turn ends.

### T005 — manual handoff index, 2026-09-30

Owner asked to leave external KB-XFMD-020 staged for the XFMD agent's next commit.
The card alone was staged; no XFMD commit was created. Shared index/history changes
remain with that agent. This is a manual summary, not recovered host turn metadata.

### T006 — 2026-09-30, execute S1

- Exact submitted owner prompt: “ok. continue on next step in our session plan”
- Category: execute the next bounded design step. Capture mode: manual.
- Skills loaded: sdp, sdp-planning, sdp-architect, sdp-traceability; agent-reported.
  Existing plan/document workflow used; no automated routine engine is claimed.
- Host thread/turn/item IDs: unavailable to this document writer.
- Work: activate KB-SDL-005, create and execute PLAN-SDP-0010 on the design phase
  branch, inspect parser and consumers, deliver contract/fixture/26 acceptance
  cases, verify baseline behavior and record original-source/revision requirements.
- Assistant work summary, not a captured final chat message: S1 design is complete.
  A three-file fixture retains all 57 baseline declarations/statements. The current
  parser correctly rejects proposed 0.6; no implementation is claimed. One closed
  System is the recommended first increment. Public linking remains explicit
  later card scope. Authoritative models and unrelated work remain unchanged.
- Next: S2 ImplementationPlan. KB-SDL-005 is ready at the handoff; it is not complete.

### S1 route disposition

The first implementation increment excludes cross-System imports/exports. This
keeps source assembly bounded while preserving the historical obligation in the
card. S2 must plan that later scope or explicitly split it before card closure.
KB-SDP-020 remains backlog until a supported pilot and consumer evidence exist.
This is the design recommendation delivered under owner-authorized S1, not a
claim of published syntax or owner acceptance of the eventual implementation.

### T007 — 2026-09-30, source-owned composition correction

- Capture: manual owner-prompt summary, not a verbatim transcript or host turn ID.
- Owner direction: reject a separate design-set file; inclusion belongs in source;
  path-addressed System membership should load the corresponding .design file.
  Explain whether repeated definitions and circular includes can be handled better.
- Skills: sdp, sdp-planning, sdp-architect; existing plan/document workflow,
  agent-reported. No runtime routine engine or automatic capture is claimed.
- Work: recover earlier System/includes intent, record correction in KB-SDL-005,
  execute successor PLAN-SDP-0011, preserve SSD1 evidence, add supersession notices
  and update this roadmap. No product parser or external XFMD changes.
- Assistant work summary: propose a source dependency graph with parse-once files,
  link-after-discovery and typed path references. Shared/cyclic source inclusion
  does not duplicate definitions; genuine duplicates/model cycles still fail.
  System membership includes channels/contracts under the proposed type extension.
- Decision authority: source-owned inclusion and rejection of authored manifests
  come from the owner; exact syntax and graph policies are design recommendations.
- Next: S2 must use P1R. No .design-set.json implementation is authorized by SSD1.

### T008 — 2026-09-30, AST forest and decorated AST discussion

- Capture mode: manual summary of the preceding visible exchange; not a verbatim
  transcript. Host IDs, start/end and subturn timestamps are unavailable.
- Owner prompt summary: support parsing smaller source trees before a later
  enclosing root, joining the forest into a larger model; explain decorated AST.
- Procedure: bounded design explanation, local source inspection and a primary
  compiler-teaching reference; not a new executed DesignPlan or adopted grammar.
- Assistant summary: file ASTs can be reusable while source and semantic relations
  form graphs. Late context requires semantic revalidation. Decorations can be
  separate immutable/context-specific data. Root-relative dependencies await root
  context. This is useful for component preview and incremental analysis.
- Outcome: endorsed for capture by the owner's following turn; KB-SDL-007 records
  it now. No implementation or measured duration is claimed.

### T009 — 2026-09-30, record proposals and put roadmap first

- Capture mode: manual English summary of the owner prompt and this work interval;
  no automatic transcript capture or available host/timing identifiers.
- Owner request: document the AST discussion in this Session and cards. Generate
  a future Gantt/table from ledger/lineage, with one planned task per step lane,
  segmented activity, turn/second axes, current-turn line and event/card/document
  cursors. Place the roadmap and table at the very top of the Session document.
- Skills loaded: sdp, sdp-planning, sdp-traceability; agent-reported. Procedure:
  existing KanBan capture and manual Session upkeep; no enforced routine claimed.
- Work: register KB-SDL-007 and KB-SDP-043 in backlog; link KB005/KB042; move the
  roadmap/table first in this Session and its local template; update the guide.
- Assistant work summary (not a captured final response): the two proposals and
  their boundaries are durable. The manual roadmap is now first. Generator,
  live cursor, event telemetry and new parser behavior remain unimplemented.
- Next: S2 implementation planning must disposition KB-SDL-007. KB-SDP-043 follows
  the separate Session/routine/client planning path and does not block SDL work.

Recorded management references for T009 (journal correlation is manual; event
write timestamps do not reconstruct turn-relative seconds):

| Event | Outcome |
| --- | --- |
| EVT-KB-SDL-000038 | Registered KB-SDL-007 in backlog |
| EVT-KB-SDP-000251 | Registered KB-SDP-043 in backlog |
| EVT-KB-SDP-000252 | Recorded KB042 refinement and roadmap-first local layout |
| EVT-KB-SDL-000039 | Linked AST refinement to KB005; ready state retained |

T009 document checks: management validation passes (54 cards, 26 management
records, 3 lineage operations, 414 events); Toolkit repository validation, local
file links in eight changed/new Markdown documents, roadmap-first section order,
preserved ledger prefix and git diff --check pass. No new Gantt rendering or
live timeline behavior was tested; the existing Mermaid block remains unchanged.
