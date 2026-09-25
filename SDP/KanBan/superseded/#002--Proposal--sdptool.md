# sdptool: project discovery, implementation plan and viewer

| Field | Value |
| --- | --- |
| id | KB-SDP-002 |
| CardState | superseded |
| project | SDP |
| type | Proposal |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | Continue in KB-SDP-017 |

## Consolidation outcome — 2026-09-25

Full scope transferred to [KB-SDP-017](../active/%23017--Proposal--sdptool-and-project-navigation.md)
under `KBO-SDP-000001`. No work remains assigned to this source. The proposal
and queue statements below are historical; the successor owns current planning.
Recorded 2026-09-25T01:41:26Z, Codex, EVT-KB-SDP-000073. Superseded is not implemented.


Registered from the owner conversation on 2026-09-23. The timestamp records registration, not a reconstructed discussion time. The directory and ledger record lifecycle status.

## Owner's intended usage

```sh
sdptool ~/git/XFMD generate ip
sdptool generate ip
sdptool view ip
```

An optional project path precedes the command; otherwise use `.`. First check whether the selected directory is a valid SDP area; otherwise check its `SDP/` child. Do not infer a parent project or switch to a subproject without an explicit rule. A versioned project marker is proposed for recognition; its format is undecided. `ip` is proposed shorthand for `implementation-plan`.

`view ip` should open the plan's main page and navigator in the configured SDP/SDL viewer, such as XFMD, with the correct source/project/tool registration. Reuse current source-based document generation, viewer adapters and prebuilt executables. Generate details from current sources on demand; do not require a full export, startup compilation or daemon. The navigator may be a snapshot; model-structure changes need an explicit refresh mechanism.

## Responsibilities and proposed plans

SDL is the language; SDP is the process. SDP tools should consume SDL/SDUI libraries without duplicating parsers or runtimes. The process CLI should behave identically for monorepos and subprojects in separate repositories. [KB-SDP-001](../completed/%23001--Proposal--Project-structure.md) owns directory/Git migration.

Support both proposing vertical slices from the SDL model and validating manually or agent-proposed slices. Each slice should deliver a bounded, verifiable capability across relevant layers/containers, with explicit prerequisites and acceptance. Tools can check coverage and dependencies; business priorities and value come from the project. Proposed plans are not automatically approved. `generate ip` must preserve approved plans and show proposed changes. Plan format, decision process and how goals/constraints are specified remain open.

Use the model, approved plan and Traceability to show a roadmap and incremental growth in use cases, features and functionality. [KB-SDP-004](../backlog/%23004--Proposal--Design-traceability.md) owns status/evidence semantics. Weak `links` must never count as implementation evidence.

## Next delivery and acceptance

Agree on the command contract and project discovery before implementation. The first vertical trial uses an existing SDL model, a small slice plan and the ledger. It shows proposal/validation diagnostics and source-linked progress without invented status. `view ip` must work from a repository root, its SDP area and an explicit path. Unknown projects/invalid sources produce diagnostics; sources and approved plans are not overwritten.

## Addition: card revisions and diff

On 2026-09-24 the owner requested content history while cards are active. Future `sdptool` should offer history/diff by stable card ID: committed revision against revision, and working draft against latest commit. Command names remain undecided. Reuse Git rather than a separate repository/version engine for every card. Show status/review events alongside Git revisions and follow explicit source/target lineage for merge/split. Rename heuristics do not define card identity. The project registry must support this in monorepos and after extraction. Acceptance: correct historical paths/content and diffs across moves, same-status revisions and merge/split; explicit diagnostics for unavailable history. [Manual workflow](../History.md).

## Current tool status

This is planning, not an implemented `sdptool` command. Existing foundations are `sdl-design` and SDL's Go CLI. [SDL scripts](../../../SDL/scripts/README.md).

## Worklog and revisions

| Time | Actor / event | Handling | Remaining work |
| --- | --- | --- | --- |
| 2026-09-24T14:52:35Z | Codex; EVT-KB-SDP-000019 | Records integrated card history/diff as a future sdptool feature; K4 delivers the manual Git workflow. | Agree on a vertical tool delivery before implementation. |

## Queue

Selected 2026-09-24T17:59:27Z by Codex; EVT-KB-SDP-000065. Predecessors:
[KB-SDP-001](../completed/%23001--Proposal--Project-structure.md) and
[KB-SDP-010](../active/%23010--Proposal--Document-consolidation.md).

Why next: project/source identity is needed before tools can open the correct
implementation plan in a monorepo or extracted project. The
[R3 proposal](../../../docs/process/Project-Discovery-and-Extraction-Proposal.md)
now supplies a concrete recognition matrix and exposes missing local markers.

First scope: confirm the discovery contract after owner review, define tested
marker/resolver/configuration milestones and integrate view ip with existing
prebuilt projection/viewer tools. Preserve approved source/plan files. Full
generate-ip slice synthesis, roadmap/evidence semantics, integrated card history
and Git extraction remain later scope. Queued means proposed next, not running
or automatically approved. No sdptool executable is delivered by this selection.

## Queue update after owner review

2026-09-24T21:41:46Z: EVT-KB-SDP-000071. KB-SDP-001 is accepted and completed.
KB-SDP-010 may await review independently; its open gate does not block this
queue. Keep CardState=queued until bounded work actually starts. The next step
is to reconcile discovery metadata with
[KB-SDP-016](%23016--Proposal--SDP-discovery-and-viewer-capabilities.md), then
confirm the first resolver/view-ip deliverable. Do not assume that closing the
planning card chose a final metadata filename or implemented its proposed schema.
