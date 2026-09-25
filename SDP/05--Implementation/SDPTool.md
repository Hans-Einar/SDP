# SDPTool implementation plan

Owner direction: 2026-09-25. Primary: [KB-SDP-017](../Agents/KanBan/active/%23017--Proposal--sdptool-and-project-navigation.md).
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
