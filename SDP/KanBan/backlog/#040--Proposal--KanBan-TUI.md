# KanBan terminal navigator

| Field | Value |
| --- | --- |
| id | KB-SDP-040 |
| project | SDP |
| type | Proposal |
| CardState | backlog |
| created | 2026-09-28T23:56:34.833964+00:00 |
| source | Owner conversation 2026-09-29 |
| next_review | At PLAN-SDP-0008 closeout |

## Scope

The owner proposes a KanBan terminal tool in the spirit of gh-tree. Model it as
its own System in ProjectGovernance, alongside the existing shell CLI. Start
with read-only board/state grouping, card details, history and opening a card in
the configured viewer. Reuse SDPTool services and durable board identities.

Do not select Go UI framework, mutation semantics, a second board database or
client-owned routine engine merely by drawing a model. gh-tree is inspiration,
not a verified dependency. Later design should inspect its implementation and
compare reuse with a small purpose-built client. Handle missing external cards
as informational references, as the current contract requires.

## Next selection

Initial architecture is included in [PLAN-SDP-0008](../../03--Architecture/Ecosystems/Plan.md).
A bounded DesignPlan should select the first user workflow and interaction
library after the owner reviews the catalog. Product implementation is unselected.

## Ecosystem modeling outcome — 2026-09-29

The [KanBanTUI model](../../SDL/ProjectGovernance/KanBanTUI/README.md) now parses and generates use case, architecture and sequence views. Its initial read-only workflow consumes SDPTool board services. The existing KanBanCLI model is separate and records the actual Bash helper. No TUI runtime/library choice was made. This proposal remains backlog for a bounded DesignPlan.
