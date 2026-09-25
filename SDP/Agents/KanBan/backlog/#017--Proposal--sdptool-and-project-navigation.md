# sdptool and integrated project navigation

| Field | Value |
| --- | --- |
| id | KB-SDP-017 |
| project | SDP |
| type | Proposal |
| CardState | queued |
| created | 2026-09-25T01:41:26Z |
| source | KB-SDP-002 and KB-SDP-016; owner clarification 2026-09-25 |
| next_review | First discovery contract and navigation delivery; coordinate with XFMD |
| tags | tooling, discovery, navigation, SDL, SDUI, XFMD |

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

The existing [sdl-design launcher](../../../../SDL/scripts/sdl-design) uses
prebuilt tools and supports a model argument/environment overrides, but defaults
to this repository's SDUI design and fixes `--project sdui`. General project
discovery/registration is the missing bridge; the launcher is not that bridge.
This is source inspection, not new GUI verification.

## Native XFMD navigation

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
is the requested hierarchy, not an assertion that a native tree API exists today.
Use stable object/viewpoint identity separately from display names and file paths.

Generate only the small overview/navigation inventory initially. Generate details
on selection from current sources; full static export remains explicit. Define
refresh after model structure changes so the tree cannot silently imply current
facts from an old snapshot. A native tree and Markdown navigator should consume
the same semantic inventory. Decide a structured adapter/contract after inspecting
existing exports; do not introduce a second projector.

The KanBan status tree is the initial browsing feature; the time-axis/zoom/lineage
graph belongs to [KB-SDP-003](%23003--Idea--KanBan-graph.md). Folder lifecycle,
CardState, ledger history and Git content revisions are distinct data. Specify
refresh after card moves/edits and behavior for inaccessible subproject boards.

XFMD owns tabs, selection, main-pane display, window targeting and lifecycle.
Its companion is [KB-XFMD-014](../../../../../xfmd-sdl-navigation/Agents/KanBan/backlog/%23014--Proposal--SDP-sidebar-and-generated-navigation.md).
SDP owns common discovery and coordination; SDL owns model/projection facts;
SDUI owns UI-language services. No XFMD code or companion card is changed here.
The companion needs this updated producer ID and owner layout clarification when
XFMD work resumes; its older direct link to the moved #016 is historical.

## Metadata decision before implementation

Reconcile the [R3 discovery proposal](../../../../docs/process/Project-Discovery-and-Extraction-Proposal.md),
[distribution manifest](../../../../SDP.manifest.yaml),
[project manifest template](../../../../Template/sdp-root/SDP-project.manifest.yaml)
and [project manifest contract](../../../../Toolkit/docs/Project-Manifest.md).
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
version compatibility. [KB-SDP-014](%23014--Proposal--KanBan-version-contract-and-distribution.md)
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
feature and functionality coverage. [KB-SDP-004](%23004--Proposal--Design-traceability.md)
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
| N3 — native navigation consumer | XFMD companion implements SDP / KanBan / SDL / SDUI tabs with bounded supported actions. Verify source-derived ViewPoint/use-case nodes, card status/filenames, one on-demand click, refresh, root switching, correct window/pane, stale-result cancellation and dirty-document retention. Decide the initial SDUI actions explicitly. |
| Later milestones | Separately select slice synthesis/validation, evidence-aware roadmap and integrated history/diff. Each needs its owning contract and concrete trial; navigation completion does not complete these. |

For the process trial, use one existing model, a small proposed/approved slice plan
and real ledger evidence. Demonstrate diagnostics and source-linked status without
invented implementation claims. The first browsing delivery does not wait for this
later trial, a full export, the KanBan timeline graph or repository extraction.

## Queue

Transferred from KB-SDP-002 by the owner's consolidation request on 2026-09-25.
[KB-SDP-001](../completed/%23001--Proposal--Project-structure.md) is accepted;
[KB-SDP-010](../active/%23010--Proposal--Document-consolidation.md) remains at
non-blocking owner review. Explicit project/source identity unlocks navigation for
all three projects, so this remains the next proposed work.

Start with N1 contract reconciliation and a bounded resolver/viewer plan, jointly
aligned with XFMD. No schema filename has been selected by closing #001 or merging
these cards. Queued means selected next for consideration, not currently in progress.

## Worklog and revisions

| Time | Actor / event | Handling | Remaining work |
| --- | --- | --- | --- |
| 2026-09-25T01:41:26Z | Codex; EVT-KB-SDP-000072 | Consolidated full scope of #002/#016 and captured native tab/tree direction; queue transferred and SDL/SDUI Refs updated. | N1 contract, then bounded implementation; see the [all-card review](../../../Maintenance/K8/Plan.md). |
