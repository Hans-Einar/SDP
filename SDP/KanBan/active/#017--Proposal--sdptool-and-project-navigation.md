# sdptool and integrated project navigation

| Field | Value |
| --- | --- |
| id | KB-SDP-017 |
| project | SDP |
| type | Proposal |
| CardState | ready |
| ScrumId | SCRUM-SDP-0002 |
| Systems | SDPTOOL |
| created | 2026-09-25T01:41:26Z |
| source | KB-SDP-002 and KB-SDP-016; owner clarification 2026-09-25 |
| next_review | First discovery contract and navigation delivery; coordinate with XFMD |
| tags | tooling, discovery, navigation, SDL, SDUI, XFMD |

## Current disposition — Scrum-0002

The feature remains open; this card is not a Sprint member. Typed partial split
KBO-SDP-000004 transfers the selected producer milestones to the six cards below.
Their acceptance now belongs there; earlier sequence descriptions below remain
background and do not create duplicate implementation assignments.

- [KB-SDP-021](../completed/%23021--Change--SDPTool-saved-design-preview.md)
- [KB-SDP-022](../completed/%23022--Change--SDPTool-discovery-contract.md)
- [KB-SDP-023](../completed/%23023--Change--SDPTool-project-viewer-bridge.md)
- [KB-SDP-024](../completed/%23024--Change--SDPTool-model-navigation.md)
- [KB-SDP-025](../completed/%23025--Change--SDPTool-KanBan-and-SDUI-inventory.md)
- [KB-SDP-026](../completed/%23026--CodeReview--SDPTool-consumer-contract-review.md)

Explicit work retained here: P0-M2 unsaved-buffer/snapshot input; T4-M2 interoperability with the separately delivered native XFMD consumer; T5 slice proposal/validation, evidence-aware implementation plans/roadmaps, and Git-backed card history/diff.

CardState remains ready for the retained feature scope; no work is in progress
on this parent. The planned Sprint selects the child cards only. No source
semantics, advanced plan generation or native XFMD implementation is silently
added to the Sprint. The [Sprint](../../Sprints/Sprint--%230001--SDPTool-preview-and-navigation.md)
records dependencies, exclusions and start/completion rules.

## Earlier assignment — owner decision 2026-09-25

SDPTool is one SDP product feature, **SdpTool**, designed through this project's
numbered 01--Mandate through 05--Implementation directories. No SDP directory,
board or implementation ledger belongs under Toolkit/SDPTool. The owner explicitly
requests first real use of the five-phase structure accepted in KB-SDP-001.

The [project entry](../../README.md), [requirements](../../02--Requirements/SDPTool.md),
[architecture/model](../../03--Architecture/SDPTool.md),
[detailed design](../../04--Design/SDPTool.md) and
[single implementation plan](../../05--Implementation/SDPTool.md) own the work.
CardState is ready at the design milestone handoff; no owner gate is requested.
Next work is the saved-file preview contract. Code stays in Toolkit/SDPTool. Shared Traceability registers DES-SDPTOOL-001,
requirements and design-check evidence; implementation status remains unimplemented.

The [pipeline explanation](../../../Toolkit/SDPTool/Navigation-and-Design-Preview.md)
records existing Go SDL parsing/projection, Rust mmdr layout, Go SDL symbols,
sdl-view:// selection, temporary-file lifetime and the optional broker. Reuse this
chain through the future facade; do not create another parser or renderer now.

Selected bounded work: establish the feature in the numbered process homes, parse
its structural model and record its requirements/ownership/verification. Then
specify the saved-file preview operation before implementing it. Full discovery,
native tree, unsaved preview and advanced process services remain later milestones.
XFMD application changes stay in its own KB-XFMD-014/015. KB-SDP-010 review remains
non-blocking; KB-SDP-018 retains the later Toolkit audit.

## Consolidated scope and authority

One project context should connect SDP process tools, SDL/SDUI tools and a
configured viewer. The owner clarified on 2026-09-25 that native XFMD should
detect a valid project/SDP area and offer an SDP sidebar with **KanBan, SDL and
SDUI** subtabs. `sdptool` is the common facade for useful existing and future
operations. It delegates to the responsible tools/libraries rather than
implementing their parsers, runtimes or projections again.

This card consolidates two proposals; it does not deliver a command, metadata
schema or XFMD feature. The first delivery is discovery and navigation. Plan
synthesis and history/diff remain later milestones, not prerequisites for browsing.

| Source | Transferred scope | Resolution and remaining questions |
| --- | --- | --- |
| [KB-SDP-002](../superseded/%23002--Proposal--sdptool.md) | Full scope: project-aware CLI, viewer bridge, slice proposal/validation, roadmap consumers, Git-backed card history and queue rationale | No scope discarded. Transfer queue here; implement browsing before advanced process commands. Command details and planning policy remain open. |
| [KB-SDP-016](../superseded/%23016--Proposal--SDP-discovery-and-viewer-capabilities.md) | Full scope: discovery/version/capability metadata, manifest reconciliation, host registration, XFMD protocol and acceptance cases | No scope discarded. Separate registrations are replaced by this shared contract/delivery. Filename/schema remain undecided; dated inspection stays in the source. |

Typed ledger operation `KBO-SDP-000001` records both full transfers. Source cards
retain their history in superseded; this does not mean their features are complete.
The earlier suggestion to handle discovery in separate cards is superseded.
The owner's subtab/tree direction replaces the XFMD companion's provisional
single generated-navigator panel sketch for planning here.

## Early standalone design preview

Owner discussion on 2026-09-25 proposes direct .design preview before the full
native tree. Add a bounded sdptool facade operation over the existing Go SDL
parser/projector/document packages; standalone supported files need not wait for
project discovery. Preserve Mermaid export and the current SDL-specific SVG
symbols. Exact command name, default views and unsaved-buffer input contract are
to be defined. Native source/preview integration belongs to
[KB-XFMD-015](../../../../xfmd-sdl-navigation/SDP/Agents/KanBan/backlog/%23015--Proposal--SDL-design-file-preview.md).
The [current pipeline guide](../../../Toolkit/SDPTool/Navigation-and-Design-Preview.md)
records existing tools, URI, temporary files and the proposed early slice.

## Project discovery and common commands

Owner examples, still proposed commands:

```sh
sdptool ~/git/XFMD generate ip
sdptool generate ip
sdptool view ip
```

An optional project path precedes the command; otherwise use `.`. Check whether
the selected directory is itself a valid SDP area, otherwise check its `SDP/`
child. Do not silently search parents, switch to a nested project or rely on Git
layout. Explicit selection must work the same for monorepos and separate repos.
Folder existence alone does not establish validity. Define missing, malformed,
unsupported and valid results, with actionable diagnostics and explicit path bases.
`ip` is proposed shorthand for `implementation-plan`.

`view ip` opens the project's plan main page and navigation in its configured
viewer. Other commands may expose useful SDL/SDUI operations through delegation;
inventory existing APIs/CLIs before settling command names. Domain logic, parser
semantics and projection ownership stay in SDL/SDUI. Existing direct language
tools remain usable. Define common errors and supported capabilities rather than
passing arbitrary unvalidated repository commands through the facade.

The existing [sdl-design launcher](../../../SDL/scripts/sdl-design) uses
prebuilt tools and supports a model argument/environment overrides, but defaults
to this repository's SDUI design and fixes `--project sdui`. General project
discovery/registration is the missing bridge; the launcher is not that bridge.
This is source inspection, not new GUI verification.

## Native XFMD consumer requirements — owned in XFMD

On successful project/SDP-area recognition, expose the SDP tab alongside XFMD's
ordinary sidebar tabs. The owner-selected inner tabs are:

| Tab | Navigation and selected action |
| --- | --- |
| KanBan | Tree of lifecycle statuses and card filenames; show CardState where helpful. Open the selected card in the main document viewer. Distinguish a Ref from its primary and follow current paths by stable identity. |
| SDL | Tree of supported ViewPoints, typed collections and model objects. For example, VP01 → UseCase → each actual UseCase. A selection asks SDL tooling for the corresponding current-source projection and shows the result in the main viewer. |
| SDUI | Entry point for project SDUI sources and supported documentation/preview operations. Exact nodes/actions are a later bounded decision; keep current Go/Fyne/SVG ownership. This tab does not establish an XFMD widget runtime. |

The ViewPoint catalog and child objects must come from SDL tooling and validated
model facts, not names extracted from rendered Markdown or a manually maintained
XFMD copy. Represent every catalog viewpoint; define disabled/empty/unsupported
states and preserve the established abstraction-phase grouping. The VP01 example
is illustrative, not a UseCase-only implementation or a fixed three-level depth.
The owner explicitly requires the same grouping for every supported groupable
concept: features, functionality, capabilities, activities, modes/states, actors,
containers, channels, contracts and data concepts where the selected profile
provides them. Group relationship selections where useful too. Candidate concepts
such as Stakeholder/UserStory remain profile-dependent, not assumed implemented.
Use variable nesting, explicit references and bounded/cycle-aware expansion for
shared objects and graph relationships. This does not assert a native tree API
exists today.
Use stable object/viewpoint identity separately from display names and file paths.

Generate only the small overview/navigation inventory initially. Generate details
on selection from current sources; full static export remains explicit. Define
refresh after model structure changes so the tree cannot silently imply current
facts from an old snapshot. A native tree and Markdown navigator should consume
the same semantic inventory. Decide a structured adapter/contract after inspecting
existing exports; do not introduce a second projector.

The KanBan status tree is the initial browsing feature; the time-axis/zoom/lineage
graph belongs to [KB-SDP-003](../backlog/%23003--Idea--KanBan-graph.md). Folder lifecycle,
CardState, ledger history and Git content revisions are distinct data. Specify
refresh after card moves/edits and behavior for inaccessible subproject boards.

XFMD owns tabs, selection, main-pane display, window targeting and lifecycle.
Its companion is [KB-XFMD-014](../../../../xfmd-sdl-navigation/SDP/Agents/KanBan/backlog/%23014--Proposal--SDP-sidebar-and-generated-navigation.md).
SDP owns common discovery and coordination; SDL owns model/projection facts;
SDUI owns UI-language services. Owner clarification on 2026-09-25: all XFMD
implementation belongs to its own cards and development session. SDP-vNow owns
`sdptool` in [Toolkit/SDPTool](../../../Toolkit/SDPTool/README.md). The existing
XFMD companion was updated with this scope and current producer link; only its
KanBan documentation and ledger were changed. XFMD now has its own SDP area and board at `SDP/Agents/KanBan`; its separate
process-adoption card owns further adoption. The former standalone path is historical.

## Metadata decision before implementation

Reconcile the [R3 discovery proposal](../../../docs/process/Project-Discovery-and-Extraction-Proposal.md),
[distribution manifest](../../../SDP.manifest.yaml),
[project manifest template](../../../Template/sdp-root/SDP-project.manifest.yaml)
and [project manifest contract](../../../Toolkit/docs/Project-Manifest.md).
`SDP/project.json`, `SDP_Version.yaml` and `SDL_version.yaml` are earlier proposed
names, not three adopted files. Choose reuse/extension or a justified separate
descriptor without competing version sources or silent release-semantic changes.

Separate descriptor schema, SDP process/profile, product release and authored
SDL/SDUI language/profile versions. Candidate facts include stable project ID,
SDP/model roots, model entries, child projects, supported viewpoint/navigation
identifiers, initial main view, navigation generation/path, KanBan location/contract
and required language-tool/renderer capabilities. Declare an explicit base for
relative paths, CLI override precedence and shared model selection.

Distinguish absent, declared, validated and unsupported optional capabilities.
Not every project has SDUI, SDL or KanBan. Agree whether unavailable subtabs are
hidden or disabled with explanation; do not mistake a boolean directory check for
version compatibility. [KB-SDP-014](../backlog/%23014--Proposal--KanBan-version-contract-and-distribution.md)
owns reusable KanBan compatibility/distribution, including standalone consumers.

Executable resolution belongs to deliberate host/tool registration. Durable
project descriptors do not contain window/client/request/lease identities or
implicitly authorize arbitrary repository executables. Reuse XFMD's existing
`--navigator`, `--sdl-tool`, `--sdl-source`, `--project`, `--renderer` and window
registration, plus `sdl-view://` generation and addressed window/pane delivery.
Keep client/request ordering, cancellation, dirty-buffer admission and resource
ownership. Broker/lease support is optional; direct generation requires no daemon.
Tools are prebuilt; opening documentation does not compile tools at startup.

## Later process capabilities retained from #002

Support both proposing vertical slices from SDL and validating owner/agent plans.
Slices need bounded value, prerequisites and acceptance across relevant layers.
The project supplies priorities; coverage/dependency checks cannot invent business
value. `generate ip` must preserve approved plans and present proposed changes.
Plan format, approval process and goal/constraint input remain open.

Use model, approved plan and Traceability for roadmaps and incremental use-case,
feature and functionality coverage. [KB-SDP-004](../backlog/%23004--Proposal--Design-traceability.md)
owns evidence/status semantics; weak links never prove implementation. Its contract
is a prerequisite for evidence-aware reporting, not for ordinary navigation.

Expose card history/diff by stable ID using Git: revision against revision, working
draft against latest commit, and lifecycle events alongside text history. Follow
explicit merge/split lineage; rename guesses do not define identity. Correctly
resolve moved historical paths, same-status changes and unavailable Git objects
across monorepos/extracted repos. Reuse the [manual workflow](../History.md);
command names remain open. Do not create a second version-control engine.

## Proposed delivery sequence and acceptance

| Step | Bounded outcome and checks |
| --- | --- |
| N1 — shared contract | Decide descriptor ownership/name/version, validity, path bases, overrides and capability behavior. Inventory existing APIs and define the model-backed tree boundary with XFMD. Include valid, minimal, unsupported and invalid examples and argument mapping. |
| N2 — resolver and viewer bridge | Implement the selected contract and first `view ip` path using prebuilt tools. Verify explicit project path, repo root, SDP area, unrelated cwd, missing/nested projects, optional capabilities, unavailable tools and preserved source/plan files. |
| N3 — general navigation producer | SDPTool supplies a versioned model-derived tree inventory for all supported groupable concepts, stable targets and selected generation, with SDL/SDUI-owned semantics. Test non-UseCase collections, variable depth, shared references/cycles, revision/refresh and KanBan inventory. |
| External consumer — KB-XFMD-014 | Native SDP / KanBan / SDL / SDUI tabs and all XFMD integration/GUI checks are planned and implemented in XFMD. They are not SDP-vNow implementation milestones. |
| Later milestones | Separately select slice synthesis/validation, evidence-aware roadmap and integrated history/diff. Each needs its owning contract and concrete trial; navigation completion does not complete these. |

For the process trial, use one existing model, a small proposed/approved slice plan
and real ledger evidence. Demonstrate diagnostics and source-linked status without
invented implementation claims. The first browsing delivery does not wait for this
later trial, a full export, the KanBan timeline graph or repository extraction.

The concrete SDPTool phase/milestone plan is maintained in
[SDP/05--Implementation/SDPTool.md](../../05--Implementation/SDPTool.md).
The wider [Toolkit audit](../backlog/%23018--Study--Toolkit-audit-and-organization.md) is
separate and does not block the bounded first tool delivery.

## Historical queue — activation supersedes this selection

Transferred from KB-SDP-002 by the owner's consolidation request on 2026-09-25.
[KB-SDP-001](../completed/%23001--Proposal--Project-structure.md) is accepted;
[KB-SDP-010](%23010--Proposal--Document-consolidation.md) remains at
non-blocking owner review. Explicit project/source identity unlocks navigation for
all three projects, so this remains the next proposed work.

Start with N1 contract reconciliation and a bounded resolver/viewer plan, jointly
aligned with XFMD. No schema filename has been selected by closing #001 or merging
these cards. Queued means selected next for consideration, not currently in progress.

## Worklog and revisions

| Time | Actor / event | Handling | Remaining work |
| --- | --- | --- | --- |
| 2026-09-25T01:41:26Z | Codex; EVT-KB-SDP-000072 | Consolidated full scope of #002/#016 and captured native tab/tree direction; queue transferred and SDL/SDUI Refs updated. | N1 contract, then bounded implementation; see the [all-card review](../../Maintenance/K8/Plan.md). |
| 2026-09-25T09:09:52Z | Codex; EVT-KB-SDP-000083 | Established SDPTool source home/plan, generalized tree scope and updated XFMD companion. | T1 discovery/delegation contract next; no executable delivered. |
| 2026-09-25T09:35:14Z | Codex; EVT-KB-SDP-000084 | Verified existing generation and recorded early standalone preview plus XFMD-owned consumer card. | Define the bounded producer operation; no Rust rewrite or native implementation selected. |
| 2026-09-25T10:11:01Z | Codex; EVT-KB-SDP-000085 | Owner-authorized activation: queued → in-progress; establish SDP feature design without a nested SDP project. | Feature model, plan placement and shared traceability. |
| 2026-09-25T10:16:50Z | Codex; EVT-KB-SDP-000087 | Five-phase local feature design and parser/projection checks delivered; in-progress → ready at this handoff. | Next: define saved-file preview operation, then implement it. Overall feature remains active. |

## Consolidated local consequence — KB-SDL-003

SDPTool reuses SDL-owned Go parser/projection services and preserves source/profile errors. The routing Ref is complete because this primary now
lives on the shared board. Product acceptance remains with this primary.

## Consolidated local consequence — KB-SDUI-001

SDPTool reuses SDUI library/preview capabilities; do not duplicate parser/layout or require an independent process area. The routing Ref is complete because this primary now
lives on the shared board. Product acceptance remains with this primary.

## Scrum-0001 review

Keep SDPTool as one feature delivery with its existing phase plan. Include reading the shared board/management history and optional SprintId/ScrumId in the navigation contract; do not absorb language changes or the timeline UI.

EVT-KB-SDP-000123; next review at the next selection or relevant dependency delivery.

## Scrum-0002 worklog

2026-09-25T12:38:14Z — EVT-KB-SDP-000130: recorded partial transfer and remaining scope; no product delivery claimed.

## Sprint-0001 outcome

2026-09-25T13:42:41Z — EVT-KB-SDP-000154: all six transferred cards are completed with
producer verification. Remaining scope above stays open; no additional feature
implementation is in progress on this parent. Skills/install Maintenance Scrums
are separately registered in #027/#028 and remain unstarted.
