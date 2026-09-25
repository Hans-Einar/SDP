# SDPTool implementation plan

Owner direction: 2026-09-25. Primary: [KB-SDP-017](../../SDP/Agents/KanBan/backlog/%23017--Proposal--sdptool-and-project-navigation.md).
This plan owns SDP-vNow deliveries. Native XFMD work is separately selected and
tracked by KB-XFMD-014, not counted as a milestone delivered in this repository.
Use stacked phase branches and one commit per completed milestone.

| Phase | Milestones and acceptance | State |
| --- | --- | --- |
| T0 — ownership and location | M1: establish this directory/plan, generalize navigation scope, update the XFMD card and register the Toolkit audit; verify boards and links | Delivered, documentation only |
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

The broader backlog review from K8 remains applicable: #017 is queued, #010's
review is non-blocking, #018 is backlog, and no held cards were selected. T1 is
the next producer delivery. The phase does not claim an implemented `sdptool`
executable or native sidebar. XFMD's card update is committed locally in its
existing branch for its own development session; no XFMD push is part of T0.
