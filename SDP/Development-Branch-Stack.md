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
history and owner proposals. [Plan/evidence](KanBan/completed/%23007--Change--KanBan-foundation.md).
At that point repository/template migration, sdptool, Traceability extensions and new
SDL keywords remained backlog.

K2 follows bb3728c on sdp/phase-k2-readable-metadata. M1 replaces YAML frontmatter with
visible tables. [Evidence](KanBan/completed/%23008--Change--Visible-card-metadata.md)
separates format changes from unchanged ledger/language contracts.

K3 follows 321e193 on sdp/phase-k3-card-lineage. M1 defines full/partial merges/splits
with typed lineage in payload 0.2, informed by XFMD's local proposal.
[Plan/evidence](KanBan/completed/%23009--Change--Card-merge-and-split.md).

R1 follows 431e47e on sdp/phase-r1-repository-organization. M1 d269bc7 inventories ownership
and activates KB-SDP-001; M2 f42859e gathers templates/project records; M3 f722dc2 gathers
SDL/document entry points with verification. [Plan](Maintenance/R1/Plan.md),
[evidence](Maintenance/R1/Evidence.md), [migration map](Maintenance/R1/Migration-map.json)
separate physical cleanup from pending profile/editorial work.

K4 follows f722dc2 on sdp/phase-k4-card-history. M1 1713778 adds
[worklogs/Git history](KanBan/completed/%23012--Change--Card-history.md);
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

## S1 — language source organization and system decisions

Branch sdp/phase-s1-language-source-organization starts at TF1 commit 1a3f5e8.
S1-M0 records the owner's three-system decision and the independently committed
XFMD SDP1-M1 bootstrap (fadf90c). It is a process/design milestone, not delivery
of System syntax, source input loading or model relocation. S1-M1/M2 remain open;
sourceinput draft code is excluded from this commit. The S1 phase is not complete.

## PM1 — shared project management

Branch sdp/phase-pm1-project-management starts at S1-M0 b2fb900; S1's uncommitted
sourceinput drafts are preserved outside this phase. PM1-M1 (088734f) consolidates
boards/history and records Scrum-0001. PM1-M2 delivers optional CLI grouping and
verification, completing MAINT-SDP-0001. This phase does not deliver S1's source
loader/model work; KB-SDL-005 and KB-SDP-020 now own those split responsibilities.
The completed phase is pushed and offered as the latest combined stack against
sdp-vNow, without merging or rewriting earlier phase branches.

## PM2 — SDPTool Sprint planning

Branch sdp/phase-pm2-sdptool-sprint-planning starts at PM1 8d85b39. PM2-M1
records Scrum-0002, planned SPR-SDP-0001, partial split KBO-SDP-000004 and two
future Maintenance-review cards. The implementation plan maps six delivery cards
to existing phases. This is planning only; uncommitted SDL sourceinput drafts
remain excluded. Evidence is recorded in the Scrum document.

## Sprint-0001 — SDPTool producer

From PM2 1ed9d54, phase branches stack as follows: p0-saved-design-preview
(4cdcc88), t1-discovery-contract (60b960a, 173ac7f), t2-project-context
(88aae42, a2e5c2b), t3-navigation-services (207fed2, 22da73d, f7d1947),
and t4-consumer-contract-review. Each uses the sdp/phase- prefix and retains
its milestone commits. P0-M2/T4-M2 remain unselected scope: delivery of the selected
producer milestones does not complete those larger phase scopes or T5. The final
combined review includes the full stack against sdp-vNow; no merge is performed.
SDP/Verification contains milestone evidence and REVIEW-SDP-0001. The prior
untracked SDL sourceinput draft remains excluded. Skills/installer Scrums remain
backlog and have received the Sprint's concrete contract findings.

## PM3 — skills Scrum

Branch sdp/phase-pm3-skills-scrum follows T4 c4aed09. PM3-M1 completes KB-SDP-027
as a Study, records Scrum-0003 and plans MAINT-SDP-0002. A disposable native Codex
catalog probe informs the adapter decision; actual skill adoption/loading remains
future work. #028 retains the next recommended Scrum for installer coordination.
This documentation/evidence phase preserves the untracked sourceinput draft and
adds no system implementation claim. See the Scrum for verification results.

## SK1 — canonical skills and activation

Branch sdp/phase-sk1-skills-activation follows PM3 aabb359. SK1-M1 (7df5afe)
prepares the profile-aware collection and migration contract. SK1-M2 (2d1c560) adopts root
Skills/ and the project adapter, migrates installed-source consumers and validates
metadata/upgrade behavior. SK1-M3 records observed native discovery and independent
behavior trials before Maintenance closure. The owner selected execution before
the installer Scrum; that broader review receives actual results afterward.

SK1-M3 completes MAINT-SDP-0002 after committed-candidate catalog checks and
four independent task trials. The complete phase is pushed for combined review
against sdp-vNow, without merge. Known baseline Traceability failures remain
separately recorded; no unrelated sourceinput changes enter these commits.

## PM4 — direct installer Maintenance plan

Branch sdp/phase-pm4-installer-maintenance-planning follows SK1 6cf74e0.
PM4-M1 completes KB-SDP-028's Study through direct selection of planned
MAINT-SDP-0003, without a Scrum or Sprint. Three phases/six milestones own the
future implementation. No installer/consumer mutation or Traceability event is
claimed by this planning milestone; the sourceinput draft remains excluded.
