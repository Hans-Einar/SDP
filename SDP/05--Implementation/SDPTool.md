# SDPTool implementation plan

Owner direction: 2026-09-25. Primary: [KB-SDP-017](../KanBan/active/%23017--Proposal--sdptool-and-project-navigation.md).
This is the single SDPTool feature plan in the shared SDP project. It was moved
from Toolkit/SDPTool on owner instruction; implementation code remains there.
[Requirements](../02--Requirements/SDPTool.md), [architecture](../03--Architecture/SDPTool.md)
and [detailed design](../04--Design/SDPTool.md) provide the design basis.

This plan owns SDP-vNow deliveries. Native XFMD work is separately selected and
tracked by KB-XFMD-014, not counted as a milestone delivered in this repository.
Use stacked phase branches and one commit per completed milestone.

| Phase | Milestones and acceptance | State |
| --- | --- | --- |
| T0 — ownership and location | M1: establish this directory/plan, generalize navigation scope, update the XFMD card and register the Toolkit audit; verify boards and links | Delivered, documentation only |
| TF1 — feature design in SDP | M1: apply five-phase local structure; establish mandate, requirements and parser-checked SDL model; relocate plan and register shared design evidence | Delivered; design only, feature still active |
| P0 — early design-preview slice (proposed priority) | M1: define and implement a small saved-file preview facade over existing Go SDL projection/document services; choose default relevant views, preserve SDL symbols and return a resource bundle with diagnostics. M2: source snapshot/revision input for unsaved-buffer preview, coordinated with XFMD KB-XFMD-015. No dependency on full project discovery/native tree | Proposed before T1; operation contract still to define |
| T1 — discovery and delegation contract | M1: inventory existing manifests, validators, SDL/SDUI commands and viewer protocol; decide one project-recognition contract with valid/minimal/invalid examples. M2: specify path/override/tool-resolution rules, versioned machine-readable results and errors, plus first command mapping | Planned |
| T2 — sdptool project context | M1: implement CLI and read-only resolver for explicit path, repo root and SDP area with no parent guessing. M2: implement configured `view ip` through prebuilt tooling and current XFMD arguments; verify errors, temporary-resource lifetime and preserved sources/plans | Planned |
| T3 — general navigation services | M1: expose catalog/phase/typed collection/object nodes using existing SDL facts, including non-UseCase collections. M2: add bounded relationship expansion, stable IDs, revision/refresh semantics and tests for shared objects/cycles. M3: expose KanBan status/card inventory and explicitly selected SDUI services through their existing owners | Planned |
| T4 — consumer adapter validation | M1: publish executable examples/contract fixtures for XFMD; test on-demand generation, unsupported capabilities and stale results using a consumer harness. M2: verify interoperability against the separately delivered XFMD consumer; record its exact version and remaining native acceptance | Planned; native delivery is external |
| T5 — advanced process commands | Separately select milestones for proposing/validating slices and preserving approved plans, evidence-aware roadmaps using KB-SDP-004, and Git-backed card history/diff. Define concrete acceptance before implementation | Later scope retained from #017 |

T2 may use the existing Markdown navigator before the native tree is available.
T3 data must support all model kinds the selected SDL profile actually exposes;
future Stakeholder/UserStory or other candidates must not be advertised as current
parser capabilities. Grouping a graph for navigation does not change its semantics.

The general Toolkit audit (#018), time-axis graph (#003) and repository extraction
do not block the first resolver/viewer delivery. Reuse existing KanBan versions
explicitly and coordinate changes with #014 rather than inventing another contract.
Neither a required daemon, startup compilation nor new language implementation is
part of the first delivery.

## Scrum-0002 execution allocation

[SPR-SDP-0001](../Sprints/Sprint--%230001--SDPTool-preview-and-navigation.md) selects the following milestones, without starting
implementation. This table refines assignment ownership; the phase acceptance
above remains the single feature plan. P0-M1 is selected first, not all of P0.
P0 and T4 retain open milestones after this Sprint; mark only their selected
milestones delivered. Preserve their phase branches and stack later work from
the then-current delivery tip rather than rebasing earlier history.

| Card | Selected phase/milestone |
| --- | --- |
| [KB-SDP-021](../KanBan/completed/%23021--Change--SDPTool-saved-design-preview.md) | P0-M1 |
| [KB-SDP-022](../KanBan/completed/%23022--Change--SDPTool-discovery-contract.md) | T1-M1/M2 |
| [KB-SDP-023](../KanBan/completed/%23023--Change--SDPTool-project-viewer-bridge.md) | T2-M1/M2 |
| [KB-SDP-024](../KanBan/completed/%23024--Change--SDPTool-model-navigation.md) | T3-M1/M2 |
| [KB-SDP-025](../KanBan/completed/%23025--Change--SDPTool-KanBan-and-SDUI-inventory.md) | T3-M3 |
| [KB-SDP-026](../KanBan/active/%23026--CodeReview--SDPTool-consumer-contract-review.md) | T4-M1 |

The source card retains P0-M2 unsaved-buffer/snapshot input; T4-M2 interoperability with the separately delivered native XFMD consumer; T5 slice proposal/validation, evidence-aware implementation plans/roadmaps, and Git-backed card history/diff.
They require later selection and are not prerequisites for this producer Sprint.
KB-SDP-026 owns producer contract/harness review, not XFMD native implementation.
System/source-set and subsequent model migration remain KB-SDL-005/KB-SDP-020;
use supported existing saved-file inputs meanwhile. No milestone is newly delivered
by this planning allocation. See Scrum-0002 for PM2-M1 planning evidence.

## T0-M1 record

Base: `94eb052`; branch `sdp/phase-t0-sdptool-foundation`. Created the source home
and plan, clarified all-concept tree navigation and separated producer/consumer
delivery. Updated XFMD's existing card and ledger only, without application code,
blueprints, product requirements, build or installation changes. Registered #018
with audit/migration acceptance; no Toolkit migration is performed now.

Verification: SDP board schema/replay/placement checks passed for 22 cards and
98 events across three boards, preserving ledger prefixes from `94eb052`.
Documentation checks preserved 105 frozen records/ledger prefixes and 574
generated outputs and validated 2,186 local file links and 130 fragments.
XFMD replay/placement checks passed for 14 cards and 47 events, preserving its
ledger prefix from `d7b54e4`; current bare file URLs resolve. Its blueprint/link
validator passed (37 blueprints, 70 requirements), and its symbol checker found
279 implemented callees. Both repositories passed `git diff --check`.

At T0-M1, the K8 backlog review applied: #017 was queued, #010's
review was non-blocking, #018 was backlog, and no held cards were selected. At T0-M1,
T1 was the next producer delivery; the later T0-M2 priority proposal is below.
The phase does not claim an implemented `sdptool`
executable or native sidebar. XFMD's card update is committed locally in its
existing branch for its own development session; no XFMD push is part of T0.

## T0-M2 — pipeline review and direct-preview capture

Inspected the actual Go CLI/projector/documents code and XFMD invocation/file
routing. Recorded the [pipeline guide](../../Toolkit/SDPTool/Navigation-and-Design-Preview.md), a proposed
early producer slice, and XFMD-owned KB-XFMD-015. The guide records a successful
installed-tool navigator and one-diagram SVG generation trial and its build-provenance
limit. No source parser, renderer, runtime or native preview code was changed.
The P0 row is the proposed next priority; T1 remains the discovery work afterward.

Checks passed: SDP boards 22 cards/99 events with preserved prefixes from
10b97c4; 2,196 local links/130 fragments and the unchanged 105 frozen records/
574 generated outputs. XFMD boards 15 cards/49 events with preserved prefix from
e850431; blueprint/link and symbol checks still report 37 blueprints, 70
requirements and 279 callees. Both diff whitespace checks passed. No native GUI
verification is claimed; the command trial is generation-only.

## TF1-M1 — local process initialization and feature design

Owner direction 2026-09-25; branch sdp/phase-tf1-sdptool-feature-design from
fbd434a. Use the accepted five-phase profile in this SDP project, rather than a
separate feature project. The earlier uncommitted Features scaffold was replaced
with phase artifacts. KB-SDP-017 is active. Existing installation templates,
project discovery descriptors, runtime code and XFMD application code are unchanged.

The model covers feature/use-case links, eight owned responsibilities, external
interfaces and planned delivery activities. All implementation-status assertions
remain planned. Preview API details are the next bounded work. Evidence is in
[VER-SDPTOOL-001](../Verification/VER-SDPTOOL-001.md); the shared Relations/ledger
records this design delivery, not runtime verification.

## P0-M1 delivery

Go saved-file preview implemented using existing SDL projection, query and bundle services, with compact default selection, structured diagnostics, cancellation and revision checks, guarded publication and caller-owned resources.

[Verification](../Verification/SDPTOOL-VER-P0-M1.md). Remaining milestones retain their existing status;
this delivery is limited to the named milestone.

## T1-M1 delivery

Selected navigation-only registration with explicit local profile, source bindings and optional references to existing version authorities. Added schema and minimal/current-project registrations without changing installer contracts.

[Verification](../Verification/SDPTOOL-VER-T1-M1.md). Remaining milestones retain their existing status;
this delivery is limited to the named milestone.

## T1-M2 delivery

Defined read-only command mapping, host-only executable precedence, JSON response/typed target boundary, synchronous viewer lifetime and source/request revision rules.

[Verification](../Verification/SDPTOOL-VER-T1-M2.md). Remaining milestones retain their existing status;
this delivery is limited to the named milestone.

## T2-M1 delivery

Implemented read-only project recognition for explicit repo/SDP-area selection, strict bounded JSON registration, confined paths, capability declarations and referenced YAML installation facts.

[Verification](../Verification/SDPTOOL-VER-T2-M1.md). Remaining milestones retain their existing status;
this delivery is limited to the named milestone.

## T2-M2 delivery

Implemented configured view ip/implementation-plan with prebuilt host arguments, current-source navigation generation, unique window registration and scoped temporary resource cleanup.

[Verification](../Verification/SDPTOOL-VER-T2-M2.md). Remaining milestones retain their existing status;
this delivery is limited to the named milestone.

## T3-M1 delivery

Implemented versioned SDL catalog/phase/abstraction, typed collection, object and diagram navigation from validated language-owned facts, with stable IDs and explicit shared references.

[Verification](../Verification/SDPTOOL-VER-T3-M1.md). Remaining milestones retain their existing status;
this delivery is limited to the named milestone.

## T3-M2 delivery

Implemented finite relationship references, stable semantic edge identity and project-bound selected generation with mandatory expected revision; source and request refresh semantics are explicit.

[Verification](../Verification/SDPTOOL-VER-T3-M2.md). Remaining milestones retain their existing status;
this delivery is limited to the named milestone.

## T3-M3 delivery

Added common KanBan/SDL/SDUI roots, card status/work/sprint/scrum metadata with descriptor-selected history, inventory revisions, local Ref resolution and SDUI frame-entry static Markdown preview through existing Go libraries.

[Verification](../Verification/SDPTOOL-VER-T3-M3.md). Remaining milestones retain their existing status;
this delivery is limited to the named milestone.
