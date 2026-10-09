# Session 0011 — SDPTool actions and JSON API

## Session roadmap

T001, 2026-10-09. Active Session; next S1. Sequence only: dates below are display
slots, not estimates or measured time. The table is the authoritative roadmap.

```mermaid
gantt
    title SDPTool actions and JSON API - sequence only
    dateFormat YYYY-MM-DD
    section Contract
    NEXT S1 Contract and plan :a,2000-01-01,1d
    section Delivery
    PLANNED S2 Catalogue and read invocation :b,after a,1d
    PLANNED S3 Blueprint actions :c,after b,1d
    PLANNED S4 Consumer and wrapper verification :d,after c,1d
    PLANNED S5 Handoff and closeout :e,after d,1d
```

| State | Step | Outcome / linked milestone | Prerequisites | Completion evidence |
| --- | --- | --- | --- | --- |
| next | S1 | Define action catalogue and JSON invocation contract; register bounded ImplementationPlan | KB051, current Go APIs/CLI, XFMD handoff | Selected schemas, scope, failure/authority rules, milestones and Git policy; pending |
| planned | S2 | One shared action registry, catalogue and JSON read operation | S1 | CLI/API equivalence, deterministic metadata and strict request/error tests; pending |
| planned | S3 | Blueprint generation/retention, assessment and revision-bound assignment actions | S2 | Real success/failure workflows preserve source, bundle and ledger contracts; pending |
| planned | S4 | Test actual consumer protocol and gh-sdp forwarding on development candidate | S3 | Subprocess input/output, structured failures, compatibility and exact binary identity; pending |
| planned | S5 | Documentation, review and XFMD handoff; reconcile KB051 and resume pointer | S4 | Evidence-backed outcome, residual scopes and Session0008 recovery note; pending |

| Field | Value |
| --- | --- |
| Session reference | SESSION-SDP-0011 |
| Status | active |
| Primary card | [KB-SDP-051](../KanBan/backlog/%23051--Proposal--Discoverable-SDPTool-actions.md) |
| Snapshot date | 2026-10-09 |
| Current step | S1 — next |
| Proposed next step | Contract and bounded execution plan before product edits |
| Execution authority | Owner selects a new Session for KB051 implementation; this turn establishes the Session and handoff |
| Predecessor | [Session0008 — paused](session-%230008--Semantic_blueprints.md) |

## Goal

Let native clients discover supported SDPTool actions and invoke them through a
versioned JSON request/response interface using the same typed services as the
human CLI. First delivery supports a complete blueprint workflow rather than
advertising arbitrary unimplemented operations. The contract should let XFMD build
menus/forms and configured toolbar actions without guessing command syntax.

Retain readable CLI commands, existing identity/revision/authority gates and
structured diagnostics. Start with one request per process. A persistent daemon,
full MCP transport, native XFMD menu implementation and automatic publication are
outside this Session's product scope. A future adapter can reuse the same services.

## Affected cards

| Card | Role | Initial | Planned final | Current | Actual final |
| --- | --- | --- | --- | --- | --- |
| [KB051](../KanBan/backlog/%23051--Proposal--Discoverable-SDPTool-actions.md) | Primary | backlog | completed after scoped delivery/evidence | queued for S1 | pending |
| [KB050](../KanBan/active/%23050--Proposal--Semantic-blueprints.md) | Context only | gate-review | unchanged by this Session | gate-review; Session0008 paused | not disposed here |
| external:KB-XFMD-030 | Consumer handoff | reported backlog at creation | XFMD agent selects native work | external, not locally authoritative | external |

## Plan register

| Local ref | Plan document/type | Readiness | Lifecycle | Outcome |
| --- | --- | --- | --- | --- |
| P1 | ImplementationPlan — create after contract refinement in S1; no ID allocated yet | planned | not registered | Phased catalogue/JSON and blueprint delivery with verification |

Keep detailed milestone status in P1 once registered; the Session tracks the route.
Contract decisions can live in one feature document without mandatory separate
Requirement/Architecture/Design plans. Follow adopted plan rules and update SDL
sources for material responsibility/contract changes before product implementation.

## Baseline and proposed execution policy

Source checkout: /tmp/sdp-blueprint-implementation, sdp/blueprint-implementation;
code baseline d2cc760, documentation baseline a4a81f8. Inspect Git status at resume.
The primary SDP-vNow and other worktrees contain unrelated work; preserve it.
S1 should select a dedicated stacked working branch from this verified baseline,
with milestone commits, unless integration facts justify a different explicit plan
policy. This Session does not authorize main merge or publication. Test distribution
preparation is distinct from publishing alpha/beta assets or replacing installations.

[Current integration handoff](../../SDPTool/XFMD-Blueprint-Integration.md) distinguishes
implemented CLI from proposed invocation syntax. gh-sdp forwards stdin; XFMD's
current job transport does not and must retain structured nonzero-exit responses.
That native consumer change belongs to external KB-XFMD-030, not this repository.

## Route changes and decisions

| Turn | Change | Authority | Impact |
| --- | --- | --- | --- |
| T001 | Separate KB051 delivery from blueprint pilot review | Owner request | Session0008 paused; this Session active; no implicit pilot acceptance |

## Turn journal

### T001 — Establish focused delivery Session (2026-10-09)

Owner input is summarized in predecessor Session0008 T015: create a new Session
for KB051 and leave an explicit resume note for the current one. Host turn IDs and
exact timings are unknown; no automatic capture is claimed. Loaded SDP/Planning
and shared document workflow (agent-reported), reusing current architecture analysis.

Work summary: allocated 0011 after checking parallel worktrees; established this
roadmap, queued KB051 and linked the paused predecessor. No product code or formal
execution plan was written this turn. Next S1 resolves the typed operation catalogue,
JSON envelope, stdin/file behavior, context/authority and error contract, then
registers the implementation milestones. This is not a captured final response.

## Closeout and return

Pending. On closeout, record actual delivery and remaining consumer/publication
work. Return to Session0008's Resume here block and S1/BP2-A owner pilot disposition.
Do not close KB050 or approve its pilot merely because KB051 is delivered.
