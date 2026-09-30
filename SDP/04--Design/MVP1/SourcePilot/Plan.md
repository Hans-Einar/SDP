# MVP1 source organization and SDUI pilot

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0005 |
| project | SDP |
| state | completed |
| PlanType | DesignPlan |
| Systems | SDL, SDUI |
| BranchPolicy | current |
| CommitPolicy | milestone |

## Outcome and authority

Owner request, 2026-09-27: restructure the MVP1 System example; separate BuckingUI
and SimulatorUI as actual containers; collect shared components; design their
interfaces with current SDUI. The intended consumer is future on-demand MVP1
navigation in XFMD. This authorizes the design pilot, not application changes in
Ponsse/XFMD, a new SDL grammar, runtime bindings, merge or release.

Work in [the existing experiment](../../../../experiments/mvp1_sdl/README.md).
Retain its experimental profile and all model identities/facts. The existing
inventory auditor is not a parser. Preserve dated evidence and pinned upstream
source hashes. Use SDUI 0.2 and the existing Go parser/layout/SVG exporter.

## Selected organization

Containers/ owns the nine registered application designs. Shared/Libraries/
contains library designs, including UIRuntime and BoxUI; Shared/UI/ contains
MVP1-wide UI grouping and common presentation definitions. Physical location is
not ownership or a promise of cross-project reuse. Keep generic candidates and
MVP1-specific composition distinguishable in a shared index; extraction is deferred.
SDUI/BuckingUI/ and SDUI/SimulatorUI/ hold application screens. No unsupported
cross-file SDUI imports or copied shared framework implementation.

BuckingUI/SimulatorUI own projection and intent routing; shared mechanisms own
representation, composition, presentation and renderer abstractions. Concrete
React/Fyne widgets stay in host implementations. Static samples never claim
machine connectivity, executable contracts or domain authority.

## Phases and acceptance

| Phase / milestone | Acceptance | Status |
| --- | --- | --- |
| MPV1-M1 | Move/split sources by ownership; update explicit includes, coverage and live links; exact normalized statement multiset unchanged; inventory passes | Delivered |
| MPV2-M1 | Operator, APT editing and simulator SDUI screens parse and generate AST, dump and SVG; inspect rendered results; document intended bindings and limitations | Delivered |
| MPV2-M2 | Reproducible export, source/tool provenance and navigation handoff; record remaining multi-file/profile support in existing backlog; verify docs/history and close plan | Delivered |

## Git and verification

Branch sdp/mvp1-source-ui from 17c5dfb37d2e8c25df9c697aa3722f4d90db73e7.
Commit MPV1-M1 and MPV2-M1/M2 deliveries; existing authorization permits phase
pushes. No merge/release. Exclude unrelated sourceinput and node/package files.
Record actual commands/results beside this plan and experiment. Audit model
facts before/after independent of paths; preserve the original inventory-audit.json.
Use current SDUI CLI, SVG inspection and source-to-output hashes. Runtime/GUI
interaction and production behavior are outside static-design evidence.

## Dependencies and follow-up

[KB-SDL-005](../../../KanBan/active/%23005--SDL--Change--System-and-source-sets.md)
owns actual Go System/source-set support;
[KB-SDP-020](../../../KanBan/backlog/%23020--Change--Shared-design-source-organization.md)
owns broader shared-model migration. This direct-owner pilot informs both without
claiming either complete. Full MVP1 uses additional experimental constructs that
need explicit profile decisions before dynamic SDL navigation. Do not flatten or
silently discard those facts to make a legacy parser accept the system.

MPV1-M1 evidence: [preservation and inventory](../../../../experiments/mvp1_sdl/evidence/MPV1/README.md). All 4,523 statements preserved; 68 files; management replay passes.

## MPV2 outcome

[Generated gallery](../../../../experiments/mvp1_sdl/preview/index.md) and
[SDUI binding guide](../../../../experiments/mvp1_sdl/SDL/MVP1/SDUI/README.md)
deliver the operator dashboard, APT editor and separate simulator console.
[MPV2 evidence](../../../../experiments/mvp1_sdl/evidence/MPV2/README.md) records
current Go parsing, AST/dump/SVG generation, deterministic export, visual inspection,
six additional viewport checks and passing relevant Go package tests.
All three registered screens also generate structural Markdown on demand through
the installed gh sdp producer with expected source revisions.

Management replay passes. The L1 document/preserved-history checker passes with
Markdown enumeration limited to Git-indexed files: 105 frozen records/prefixes,
574 existing generated outputs, 2705 links and 134 fragments at the pre-closeout
check. Its unfiltered scan encounters unrelated untracked node_modules documentation;
that dependency tree is excluded, untouched and not claimed verified.

The existing backlog cards were reviewed and retain their broader unfinished
scope. [The navigation handoff](../../../../experiments/mvp1_sdl/Navigation.md)
records required multi-file/profile support and a concrete first MVP1 workflow.
This plan is complete for source organization and proposed static UI designs.
Full SDL navigation, actual domain bindings, native interaction and owner visual
acceptance remain separate work. No Ponsse/XFMD application changes or release.
