# Evaluate SDL viewpoints through integrated XFMD navigation

| Field | Value |
| --- | --- |
| id | KB-SDP-032 |
| project | SDP |
| type | Study |
| CardState | backlog |
| Systems | SDL, SDPTOOL |
| created | 2026-09-25T22:20:56Z |
| source | Owner conversation 2026-09-26: SDL context, implementation drift and post-main XFMD adoption |
| next_review | After MAINT-SDP-0005 main consolidation |

## Owner intent

Current generated viewpoints are an initial proposal requiring practical design
feedback. The owner is waiting for integrated SDP navigation in XFMD to browse
the model and assess what the views reveal or obscure. Generated output alone
does not establish a useful viewpoint catalog.

## Scope and outcome

Prepare a Study and refinement plan using real SDL source models. Inventory the
actual generators and navigation contracts before proposing new views. Evaluate
requirements/use cases, features/functionality, system/container structure,
channels/contracts and cross-cutting relationships across abstraction levels.
Capture navigation tasks, missing context, ambiguous notation, duplication and
useful drill-down/filter behavior. Include an assignment-authoring task so results
feed [KB-SDP-031](../active/%23031--Study--SDL-assignment-bundles-and-blueprints.md).

Keep the SDL model as source of truth. Reproduce all diagrams via the tools;
record model revision, generation parameters and derived-versus-authored facts.
Do not hand-author purported generated diagrams or assume all viewpoints are
implemented. Test both overview and selected on-demand generation.

## Acceptance and ownership

A repeatable walkthrough records actual owner feedback, source-linked examples,
limitations and a prioritized, bounded refinement plan. Pending owner feedback
must remain pending rather than being inferred by an agent. Coordinate with
[KB-SDP-017](../active/%23017--Proposal--sdptool-and-project-navigation.md): that
card owns the navigation bridge; this card evaluates viewpoint usefulness.
XFMD owns its native sidebar implementation and its own KanBan work. No XFMD
application change or new broad renderer integration is authorized here.

## Worklog

2026-09-25T22:20:56Z: Registered in backlog before main integration; EVT-KB-SDP-000173.
