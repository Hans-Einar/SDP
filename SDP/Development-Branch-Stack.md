# SDL/SDUI — phase branches and milestone commits

Owner decision, 2026-09-22: one branch per phase, stacked on its predecessor, and
separate commits for completed milestones. Push after each completed phase is
authorized in the same session. Combined PR base: sdp-vNow, originally
9ad432407004080dd7f4f0ab06d107523f4316fd.

| Phase branch | Parent | Milestones |
| --- | --- | --- |
| sdl-sdui/phase-baseline | sdp-vNow | B0-M1: existing SDUI 0.2 prototype/examples/evidence/Go direction and Git workflow |
| sdl/phase-v1-viewpoints | sdl-sdui/phase-baseline | M1 language/AST/validation; M2 projections; M3 shared model/export/checkpoint |
| sdl/phase-v2-data-contracts | sdl/phase-v1-viewpoints | M1 data/wire profile; M2 VP09/VP10; M3 verified example/push |
| sdl/phase-v3-channels | sdl/phase-v2-data-contracts | M1 checked scenarios; M2 sequences/MessageSets; M3 rendering/push |
| sdl/phase-v4-integrated-design | sdl/phase-v3-channels | M1 G1–G5 in SDL; M2 implementation report; M3 export/checkpoint/push/PR |
| sdl/phase-g6-navigation-design | sdl/phase-v4-integrated-design | D1 navigation/scenarios; D2 navigation-only, levels/notation; M1–M6 planned at that point |
| sdui/phase-g1-frontend | sdl/phase-g6-navigation-design | M1 parser/AST; M2 semantics; M3 ported exports |

The first two phases deliberately record already-completed local work retrospectively;
they do not claim earlier sessions had these commit boundaries. V0/V1 already shared
files and are recorded together on V1. No hypothetical old V0 version was reconstructed.
Baseline preserves SDUI and its direction; some links point forward to V1. V1 milestones
separate language, consumer and validated example/evidence; its head is the combined review.

From V2, milestones were followed as work proceeded. Check scope, tests and generated
artifacts before each commit; create the phase branch before changes and commit only
completed deliveries. [V2–V4 scope](History/checkpoint-1/08-SDL-Viewpoints-and-Implementation-Status.md),
[Go phases](../SDUI/docs/implementation-plan.md).

Use the last completed phase as PR head against sdp-vNow. Preserve milestone commits;
do not squash, delete phase branches or rewrite history as automatic cleanup.

| Implementation phase | Pushed head / branch | Predecessor |
| --- | --- | --- |
| G1 | ca5aa91 — sdui/phase-g1-frontend | G6 design |
| G2 | 6a4d968 — sdui/phase-g2-presentation | G1 |
| G3 | d29a48f — sdui/phase-g3-runtime | G2 |
| G4 | d5430f0 — sdl/phase-g4-runtime | G3 |
| G6 | 1d52677 — sdl/phase-g6-navigation | G4 |
| G5 | sdl-sdui/phase-g5-codegen; final M4 32fadca | G6 |

G6 precedes G5 because G5-M3 consumes navigation. XFMD consumer changes are isolated
in its PR #38; at that delivery both stacks were pushed, not merged. G5 comprises M1
constructors, M2 execution parity, M3 reproducible documents and M4 port cleanup/checkpoint.
The combined PR then used G5's head. No merge was part of delivery.

G7 follows 32fadca on sdl/phase-g7-launch. M1 provides source-based browsing through a
prebuilt-tool launcher, as a separate startup follow-up.

K1 follows G7 d03eb78 on sdp/phase-k1-kanban. M1 establishes project boards, append-only
history and owner proposals. [Plan/evidence](Agents/KanBan/completed/%23007--Change--KanBan-foundation.md).
At that point repository/template migration, sdptool, Traceability extensions and new
SDL keywords remained backlog.

K2 follows bb3728c on sdp/phase-k2-readable-metadata. M1 replaces YAML frontmatter with
visible tables. [Evidence](Agents/KanBan/completed/%23008--Change--Visible-card-metadata.md)
separates format changes from unchanged ledger/language contracts.

K3 follows 321e193 on sdp/phase-k3-card-lineage. M1 defines full/partial merges/splits
with typed lineage in payload 0.2, informed by XFMD's local proposal.
[Plan/evidence](Agents/KanBan/completed/%23009--Change--Card-merge-and-split.md).

R1 follows 431e47e on sdp/phase-r1-repository-organization. M1 d269bc7 inventories ownership
and activates KB-SDP-001; M2 f42859e gathers templates/project records; M3 f722dc2 gathers
SDL/document entry points with verification. [Plan](Maintenance/R1/Plan.md),
[evidence](Maintenance/R1/Evidence.md), [migration map](Maintenance/R1/Migration-map.json)
separate physical cleanup from pending profile/editorial work.

K4 follows f722dc2 on sdp/phase-k4-card-history. M1 1713778 adds
[worklogs/Git history](Agents/KanBan/completed/%23012--Change--Card-history.md);
integrated tooling/graph remained backlog.

R2 follows 1713778 on sdp/phase-r2-document-consolidation. Its
[plan](Maintenance/R2/Plan.md) activates KB-SDP-010 for authority/status conflicts,
without new language/installation rules. M1 68ba33e records ten conflicts; M2 918fa46
resolves them and consolidates runtime boundaries/SDUI entry points with CLI/link checks.
At R2 completion KB-SDP-010 remains active for candidate/process harmonization;
phase/template profile and registry remain in KB-SDP-001.

L1 follows 918fa46 on sdp/phase-l1-english-documentation: M1 English policy/boards,
M2 maintained narrative documentation, M3 generator language and regenerated outputs.
[Plan](Maintenance/L1/Plan.md) records verification and the subsequent CardState work.

K5 follows L1 32f6ae5 on sdp/phase-k5-card-state. M1 49310de adds CardState;
M2 delivers the read-only CLI, installer and tests, with owner-authorized concurrent
KB-SDP-014 registration. [Plan](Maintenance/K5/Plan.md).

K5-M3 0c73cb5 corrects the newly introduced duplicate event ID in a separate
commit, preserving the original in Git/evidence. R3 follows on
sdp/phase-r3-profile-and-authority: M1 aca12a6 supplies the reviewable phase/
discovery/extraction package; M2 separates active core from candidates and records
owner gates plus the queued successor. [R3 plan](Maintenance/R3/Plan.md).

K6 follows R3 84089ae on sdp/phase-k6-terminal-links. K6-M1 adds OSC 8 card
links on terminals with plain pipe output; [evidence](Maintenance/K6/Plan.md).

K7 follows 5249eb6 on sdp/phase-k7-owner-review. K7-M1 records owner closure
of KB-SDP-001 and non-blocking deferred review of KB-SDP-010, with updated links
and queued follow-up. [Plan](Maintenance/R3/Plan.md).

K8 follows K7 29d5828 on sdp/phase-k8-backlog-consolidation. K8-M1 reviews all
12 backlog cards, merges KB-SDP-002/016 into queued KB-SDP-017 with typed lineage,
and captures native XFMD tab/tree direction. [Plan](Maintenance/K8/Plan.md).

T0 follows K8 94eb052 on sdp/phase-t0-sdptool-foundation. T0-M1 establishes
Toolkit/SDPTool, generalizes navigation scope and separates XFMD-owned work;
KB-SDP-018 captures the wider Toolkit audit.
[Plan and evidence](05--Implementation/SDPTool.md).

T0-M2 records the verified existing Go/Rust document pipeline and a proposed
early direct-design-preview producer slice. XFMD KB-XFMD-015 owns its native
consumer; existing diagram generation is reused, not reimplemented.

TF1 follows fbd434a on sdp/phase-tf1-sdptool-feature-design. TF1-M1 applies
the accepted five-phase profile locally and activates SDPTool as an SDP feature;
[plan](05--Implementation/SDPTool.md), [evidence](Verification/VER-SDPTOOL-001.md).
