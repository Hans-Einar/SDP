# SDL/SDUI — Go implementation plan

**Status 2026-09-22:** SDL V0–V4 are delivered and pushed as phase branches. The owner authorized all G phases in one continuous session, with phase branches, milestone commits and end-of-phase pushes. G1–G6 are delivered within bounded profiles. [Implementation evidence](../go/evidence/G1.md) remains separate from the [generated design plan](../design/viewpoints/implementation.md). Model verified statuses record completed phase/milestone checks; checkpoint supplement 11 identifies semantic limits and remaining product work.

**ID:** SDUI-PLAN-003 · **Revision:** 2026-09-22. Replaces PLAN-002's Rust/C-ABI/FOX P0–P6 track. New milestones use G prefixes; older P-phase evidence is historical. [Target architecture](target-architecture.md); [checkpoint](../../SDP/History/checkpoint-1/07-SDUI-0.2-and-Go-Direction.md).

Scope: Go SDL/SDUI parsers/runtimes, shared SDUI layout, SVG, first Fyne host, model reload and Go generation. Python port sources were removed after verified consumer ports; frozen fixtures/historical evidence remain.

Design basis, 2026-09-22: the [shared SDL structural model](../design/README.md) describes G1–G6 responsibilities/dependencies and passes the existing parser. This establishes design coverage, not implementation completion by itself.

## G0 — updated foundations and directories

| Milestone | Acceptance | Status |
| --- | --- | --- |
| G0-M1 | Checkpoint describes SDUI 0.2, actual implementation, Go direction and historical differences | Delivered |
| G0-M2 | Active architecture/plan/handoff instructions agree | Delivered |
| G0-M3 | Separate SDL/SDUI Go areas with parser/runtime responsibilities | Directory foundation delivered; implemented later in G1–G6 |

## G1 — executable SDUI frontend in Go

| Milestone | Delivery and acceptance |
| --- | --- |
| G1-M1 | **Delivered:** Go module, syntax-only CLI, parser/AST/source positions; example ASTs match Python fixtures |
| G1-M2 | **Delivered:** Validation/normalization, relative rules, regions, instance paths and source limits; ported positive/negative Python cases |
| G1-M3 | **Delivered:** Concept1 AST/diagnostics and identical console/Markdown dumps; callback-free SVG/HTML gallery. General geometry in G2; temporary gallery formats removed in G5-M4 |

Requirements R01–R09, R13–R15, R19, R21–R24. Parser opens no SDL sources. No Go/Python SDUI 0.1 fallback. Port oracle frozen as test data; active Python consumers removed.

## G2 — shared layout, SVG and first Fyne window

Depends on G1-M2. Owner: SDUI implementation.

| Milestone | Delivery and acceptance |
| --- | --- |
| G2-M1 | **Delivered:** Measurement/font contract; relative axes, ratio, rows/groups, header/body/footer, gap/padding and explicit overflow limits |
| G2-M2 | **Delivered:** General measured SVG export; small frame/button/input, then Concept1 without manually placed fixture; geometry/visual checks |
| G2-M3 | **Delivered:** Fyne window sharing geometry, buttons/inputs, Tab/focus and registered local Go function; no SDL dependency |
| G2-M4 | **Delivered:** Bounded Markdown provider and separate Mermaid resources; documented profile/negative cases, no full diagram-coverage claim |

Requirements R10/R11/R14/R15/R17/R20/R21/R24. Start with a small vertical M1–M3 trial; full Concept1/Markdown follows. Measure startup, resize and memory before claiming lightweight operation. Fyne handles interaction; core/SVG work without GUI. SVG exports selected state without callbacks.

## G3 — UI runtime and SDUI hot reload

Depends on G1-M2/G2-M3. Owner: SDUI implementation.

| Milestone | Delivery and acceptance |
| --- | --- |
| G3-M1 | **Delivered:** Typed events/updates, stable widget instances/revisions; explicit unbound-button status |
| G3-M2 | **Delivered:** File change → parse/validate → publish; invalid source retains last valid UI and shows source diagnostics |
| G3-M3 | **Delivered:** Preserve compatible value/draft/focus; handle type changes/deletion; reject stale events/post-teardown callbacks; no repeated domain actions on reload |

Requirements R12/R16/R18/R25. Separate UI/domain state. Define event/UI-thread boundaries before asynchronous work; no complex plugin/process mechanism here.

## G4 — bounded SDL parser/runtime and SDUI bridge

SDL parsing may start alongside G1; integration depends on G3-M1. SDL owns semantics/runtime; SDUI owns its bridge endpoint.

| Milestone | Delivery and acceptance |
| --- | --- |
| G4-M1 | **Delivered:** Go design-core port with grammar/tests; checkpoint/MVP1 candidates not automatically supported |
| G4-M2 | **Delivered:** Bounded named-action profile with typed inputs/results and registered Go functions; invalid/missing bindings rejected |
| G4-M3 | **Delivered:** SDUI button/input → SDL action → Go function → UI update; source maps and shared contract with explicit simulated domain |
| G4-M4 | **Delivered:** SDL reload with last valid model, state rules and in-flight handling; Go changes rebuild/restart |

Requirements R12/R18/R25. Use bounded MVP1 EditAptCell after simple binding works. The full 66-file corpus is not yet a parser acceptance target. Infer no machine/domain algorithms from structure; language implementation does not automatically change Ponsse product code.

## G5 — Go generation and shared documentation

Depends on resolved G4 profile/shared runtime. Joint SDL/SDUI ownership.

| Milestone | Delivery and acceptance |
| --- | --- |
| G5-M1 | **Delivered:** Generated Go creates equivalent models/bindings and builds with separate handwritten domain functions |
| G5-M2 | **Delivered:** File-based development/generated programs yield equivalent event traces/state for agreed profile; incomplete semantics diagnosed |
| G5-M3 | **Delivered:** Reproducible SVG/Markdown from selected UI/state, with source/tool versions and rendering evidence |
| G5-M4 | **Delivered:** Port complete: obsolete active Python entry points/fixture locations removed or replaced, links/commands updated; one implementation per profile |

Requirements R19/R26. Reuse Go build cache on restart; dynamic machine-code replacement/separate workers are not hot-reload acceptance requirements.

## G6 — navigable documents and on-demand generation

**G6-M1–M6 delivered.** [Design contract/XFMD handoff](../../SDL/docs/integration/SDL-Navigable-Viewpoints-Design.md). SDL owns projection/publication; XFMD owns document panels/link routing. Document viewing is separate from G2's interactive Fyne host.

| Milestone | Delivery and acceptance |
| --- | --- |
| G6-M1 | **Delivered:** Navigator/overview without detailed diagrams or static package; A0–A5 directories, type inventory and stable links/anchors from the same model; ported projector; optional monolithic export; checked links/images/determinism |
| G6-M2 | **Delivered:** Typed relation/direction/depth/level/mode selection by click/CLI; revision and publication of selected document/resources only; equivalent to corresponding full export; errors retain last view |
| G6-M3 | **Delivered:** XFMD navigation/main panels, registered reader adapter and explicit window/panel; click, focus, multiple-window and closed-target tests |
| G6-M4 | **Delivered:** Optional Go service with local IPC, cache/invalidation, leases, request ordering, quotas and cleanup; resources survive document changes/reload |
| G6-M5 | **Delivered:** Consistent symbol/arrow profile, UML where semantics match; actor figures/use-case ellipses; actual shape/marker tests beyond exit codes |
| G6-M6 | **Delivered:** Explicit class/relation profile with multiplicity and aggregation/composition; grammar/validation before source-linked diagrams; no inference from contains |

G6-M1 depends on G4-M1's structural frontend, not SDL runtime or G5 generation. Viewpoint/source-map port moved from G5-M3 to G6-M1; G5-M3 consumes it. M2→M3→M4 are sequential; M5 depends on M1 and can accompany host work; M6 follows M5 with a resolved class contract. Earlier Python generator remains frozen port evidence. URI/IPC/XFMD flags were verified against the separate implementation; see G6 evidence. G6-D1/D2 were design deliveries; later M1–M6 implementation has separate evidence. [Levels/notation](../../SDL/docs/integration/SDL-Viewpoint-Levels-and-Notation.md) defines export forms, A0–A5, Mode/State and diagram semantics.

## Scope and reuse

Fyne is the first interactive host. G6 delivers XFMD document navigation; FOX-based SDUI widgets, C ABI and Bubble Tea remain deferred. Existing renderer/worktree code can provide ideas, tests and suitable algorithms with provenance; no extracted crate or Mermaid merge is required. Test Markdown/Mermaid against an agreed profile. Do not build extra renderers merely to preserve options.

SDL/SDUI each document code ownership in their READMEs. This is the shared plan; do not duplicate phase plans. [Go port inventory](../../SDL/docs/integration/SDL-Go-Port-Inventory.md). Git uses [one branch per phase and commit per milestone](../../SDP/Development-Branch-Stack.md); end-of-phase pushes are authorized.

## G7 — simple document-navigator launch

Owner clarification after G6: daily use starts from SDL source and generates selected details on click. Full static export is an explicit archival/export option.

| Milestone | Delivery and acceptance |
| --- | --- |
| G7-M1 | **Delivered:** Launcher without required arguments; prebuilt SDL tool, no compilation; regenerate navigator/overviews only; register source/project/renderer in correct XFMD; direct mode without daemon; checked errors/cleanup |
