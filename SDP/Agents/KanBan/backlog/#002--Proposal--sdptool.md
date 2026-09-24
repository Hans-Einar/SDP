# sdptool: project discovery, implementation plan and viewer

| Field | Value |
| --- | --- |
| id | KB-SDP-002 |
| project | SDP |
| type | Proposal |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | 2026-09-30 |

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

SDL is the language; SDP is the process. SDP tools should consume SDL/SDUI libraries without duplicating parsers or runtimes. The process CLI should behave identically for monorepos and subprojects in separate repositories. [KB-SDP-001](../active/%23001--Proposal--Project-structure.md) owns directory/Git migration.

Support both proposing vertical slices from the SDL model and validating manually or agent-proposed slices. Each slice should deliver a bounded, verifiable capability across relevant layers/containers, with explicit prerequisites and acceptance. Tools can check coverage and dependencies; business priorities and value come from the project. Proposed plans are not automatically approved. `generate ip` must preserve approved plans and show proposed changes. Plan format, decision process and how goals/constraints are specified remain open.

Use the model, approved plan and Traceability to show a roadmap and incremental growth in use cases, features and functionality. [KB-SDP-004](%23004--Proposal--Design-traceability.md) owns status/evidence semantics. Weak `links` must never count as implementation evidence.

## Next delivery and acceptance

Agree on the command contract and project discovery before implementation. The first vertical trial uses an existing SDL model, a small slice plan and the ledger. It shows proposal/validation diagnostics and source-linked progress without invented status. `view ip` must work from a repository root, its SDP area and an explicit path. Unknown projects/invalid sources produce diagnostics; sources and approved plans are not overwritten.

## Addition: card revisions and diff

On 2026-09-24 the owner requested content history while cards are active. Future `sdptool` should offer history/diff by stable card ID: committed revision against revision, and working draft against latest commit. Command names remain undecided. Reuse Git rather than a separate repository/version engine for every card. Show status/review events alongside Git revisions and follow explicit source/target lineage for merge/split. Rename heuristics do not define card identity. The project registry must support this in monorepos and after extraction. Acceptance: correct historical paths/content and diffs across moves, same-status revisions and merge/split; explicit diagnostics for unavailable history. [Manual workflow](../History.md).

## Current tool status

This is planning, not an implemented `sdptool` command. Existing foundations are `sdl-design` and SDL's Go CLI. [SDL scripts](../../../../SDL/scripts/README.md).

## Worklog and revisions

| Time | Actor / event | Handling | Remaining work |
| --- | --- | --- | --- |
| 2026-09-24T14:52:35Z | Codex; EVT-KB-SDP-000019 | Records integrated card history/diff as a future sdptool feature; K4 delivers the manual Git workflow. | Agree on a vertical tool delivery before implementation. |
