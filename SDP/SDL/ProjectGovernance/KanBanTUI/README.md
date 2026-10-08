# KanBanTUI — interactive read-only work navigator

## Boundary and implementation state

This proposed system is a terminal navigator inspired by the owner's gh-tree
workflow. No KanBanTUI source, binary, toolkit selection or gh-tree code reuse is
claimed. Candidate binary name `sdp-kanban` is provisional. `KanBanTuiHost` is one
proposed process; its query/view/reader units may become libraries, not additional
deployed containers. The [proposal card](../../../KanBan/backlog/%23040--Proposal--KanBan-TUI.md)
records follow-on work separately from this model delivery.

The first boundary is **read-only**: browse card states, filter/group work, inspect
a selected card, and show freshness/unavailability. Reuse the validated
[SDPTool board service](../../../../SDPTool/kanban.go) and typed
[navigation inventory](../../../../SDPTool/navigation.go) rather than create a
second card parser and workflow engine. Existing Go functions support useful read
facts; the proposed client/public binding and terminal interaction are not built.
The shell [KanBanCLI](../KanBanCLI/README.md) remains independently useful.

## Interaction and contracts

`BoardQuery` requests a selected project snapshot, `BoardView` groups/selects,
`CardReader` reads the selected bounded path, and `FreshnessDisplay` distinguishes
an unavailable board from a validated empty one. Group by CardState, Sprint and
Scrum when present; preserve card namespaces and unverified external references.
Use returned paths/revisions instead of guessing a checkout from an external ID.
Card movement changes paths, not stable identity. The view must refresh/reconcile
selection after revision changes rather than open an unrelated or stale target.

`BoardQueryArguments` and `BoardSnapshotEnvelope` are conceptual read envelopes,
not a new committed schema replacing SDPTool's `Tree`/`Node` types. Exact API
binding, limits, terminal navigation keys, card preview, Git diff and history
rendering are follow-on design choices. `BoardLoaded` and
`UnavailableBoardShown` are sample read paths; they are not product tests.

`HistoryReadDelivery` is explicitly planned after the initial read view. Reading
ProjectManagement events and Git revisions may later explain merges/splits and
document changes. This is not a second Git implementation or a new writable
management ledger. A timeline with expanding week/day/hour resolution remains a
separate feature candidate, not a promised part of the first TUI.

`ExternalSdpToolBoardService` is an unowned endpoint stub for the separate SDPTool
system. Calling a linked Go read library is compatible with the model's logical
Channel; a new server/daemon is not required. No channel offers card transitions,
release actions or approvals. Closing the TUI must have no effect on work state.

## Model and validation boundary

[System.design](System.design) is the independently parsed entry for **design-core
0.5**. Its declarations belong to this file only; directory names are catalog
metadata, not language namespaces or imports. Activities distinguish implemented
source responsibilities from planned delivery. A valid model or rendered scenario
is not runtime evidence, owner acceptance or authorization to implement a proposal.

From the repository root, with a built SDL CLI:

```sh
sdl check SDP/SDL/ProjectGovernance/KanBanTUI/System.design
sdl ast SDP/SDL/ProjectGovernance/KanBanTUI/System.design
sdl viewpoints SDP/SDL/ProjectGovernance/KanBanTUI/System.design --format static --output /tmp/sdp-kanbantui-views
```

Build the CLI from `SDL/go` with `go build -o /desired/bin/sdl ./cmd/sdl`.
The authoring verification used the repository CLI with Go 1.27.1; the enclosing
[Ecosystems plan](../../../03--Architecture/Ecosystems/Plan.md) records consolidated
candidate hashes and export evidence. Generated output belongs in temporary or
explicitly marked derived directories, never in this authored source folder.
The [ecosystem index](../README.md) identifies cross-system dependencies.
Navigation registration must select this entry explicitly; folder placement alone
does not make it available to a viewer.
