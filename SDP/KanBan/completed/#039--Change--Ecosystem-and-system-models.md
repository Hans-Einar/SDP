# Establish the SDP ecosystem and system models

| Field | Value |
| --- | --- |
| id | KB-SDP-039 |
| project | SDP |
| type | Change |
| CardState | completed |
| created | 2026-09-28T23:56:34.833964+00:00 |
| source | Owner conversation 2026-09-29 |
| next_review | Closed for the bounded catalog delivery |
| PlanId | PLAN-SDP-0008 |

## Scope

Author a navigable SDL source catalog for ProjectGovernance, SDL and SDUI.
Describe current binaries/libraries and proposed parser/compiler/runtime/adapter
boundaries without implying that the proposed tools already exist. Preserve old
model authority until explicit migration and source-set support. The owner
explicitly authorizes this architecture/model delivery now.

## Plan and acceptance

[PLAN-SDP-0008](../../03--Architecture/Ecosystems/Plan.md) owns phases, validation
and closeout. Deliver independently parseable SDL models, system responsibility
and implementation maps, a reproducible viewpoint export, and clear next work.
The card covers this initial catalog, not every future system's implementation.

## Worklog

- Selected and activated on 2026-09-29; source inventory and modeling begin.

- E1-M1/E2-M1: five ProjectGovernance models and source maps delivered; Go structural checks, AST and viewpoint generation pass. New capabilities remain proposals.

- E2-M2: seven SDL and five SDUI system models delivered, grounded in tracked Go packages/commands. Independent model checks, AST and viewpoint generation pass; packaging is not extracted.

## Ecosystem modeling outcome — 2026-09-29

Delivered the [catalog](../../SDL/Catalog.md), three ecosystems, 17 system models plus one collaboration model, source/binary maps, current navigation registrations, and [checked evidence](../../03--Architecture/Ecosystems/Evidence.md). All 18 models parse, export AST and generate revision-bound navigation selections and static views. The initial independent review found no material issues; final closeout review is recorded in the evidence. This completes the authorized initial modeling work. Runtime implementation, separate binaries, source-set migration and external viewer changes remain with linked cards.
