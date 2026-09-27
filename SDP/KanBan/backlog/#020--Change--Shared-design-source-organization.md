# Organize shared system design sources by phase

| Field | Value |
| --- | --- |
| id | KB-SDP-020 |
| project | SDP |
| type | Change |
| CardState | backlog |
| created | 2026-09-25T11:42:10Z |
| source | SCRUM-SDP-0001; owner-conversation-2026-09-25 |
| ScrumId | SCRUM-SDP-0001 |
| Systems | SDL, SDUI, SDPTOOL |
| next_review | After MAINT-SDP-0001, before selecting the next implementation |

## Bounded delivery

Apply one root SDP process to SDL, SDUI and SDPTool. Use system-specific content
within Requirements/Architecture/Design/Implementation; detailed design belongs
under actual containers, with shared libraries/contracts defined once. Phase
folders do not assign A0–A5. Move the interim plan from its language process area
and remove duplicate active authorities with explicit migration records.

Split SDUI/design/architecture.design into requirements, system/container topology,
container/library details, collaboration/contracts and implementation records.
Preserve symbol identities and all intended facts; distinguish semantic changes
needed for the three-system decision from a byte/placement-only migration.
Canonical model comparisons, original-source diagnostics, working CLI/navigation,
selected projections, historical evidence and links are required. Generated
viewpoints come from the SDL toolkit, never manual substitute diagrams.

## Dependency and acceptance

[KB-SDL-005](%23005--SDL--Change--System-and-source-sets.md) owns the supported
System/source input contract. Do not keep a competing handwritten aggregate model
or introduce an include preprocessor to bypass it. SDL and SDUI Go modules remain
source homes; no XFMD code or installer migration is included. Record remaining
profile gaps explicitly. The current KanBan migration is MAINT-SDP-0001 and does
not fulfill this design-source migration.

## Lineage and remaining impact

KBO-SDP-000003 splits [KB-SDL-004](../superseded/%23004--SDL--Change--Language-source-organization.md).
SDP/SDUI Ref consequences are owned here: shared phase placement, no duplicated
language process trees, preserved histories and one authoritative contract/model
per concern. The former [S1 plan](../../Maintenance/S1/Plan.md) is background and
must be updated to these two successor deliveries before implementation.

## Consolidated local consequence — KB-SDP-019

Shared root process and three-system organization belongs to KB-SDP-020. The routing Ref is complete because this primary now
lives on the shared board. Product acceptance remains with this primary.

## Consolidated local consequence — KB-SDUI-002

SDUI design sources move into shared phase homes; preserve language-owned internals and dated exports. The routing Ref is complete because this primary now
lives on the shared board. Product acceptance remains with this primary.

## Scrum-0001 review

Keep the dependent shared model/source migration separate from the prerequisite language/input contract.

EVT-KB-SDP-000121; next review at the next selection or relevant dependency delivery.

## MVP1 pilot follow-up — 2026-09-27

PLAN-SDP-0005 pilots a system-owned source tree in experiments/mvp1_sdl: Containers, Shared/Libraries, Shared/UI, Contracts, Scenarios, Governance and SDUI. This refines the earlier phase-folder placement proposal for MVP1; numbered folders remain process/document organization. The broader SDL/SDUI/SDPTool migration is still pending. Preserve one declaration/fact authority and use this pilot as migration evidence, not as multi-file compiler evidence.

[Pilot and navigation handoff](../../../experiments/mvp1_sdl/Navigation.md).

## Source convention and template follow-up — 2026-09-28

MAINT-SDP-0009 documents SDP/SDL/<System>/ as the home for new authored models,
with container details and system-owned SDUI screens. This supersedes the earlier
phase-folder source-placement proposal above; numbered folders remain process
and reader-document organization. Current five-phase installation seeds now carry
the convention. Existing models remain authoritative at their current paths until
this card's migration is selected and verified. Parser support and registration
are not inferred from directory structure. This card remains backlog.

[Source guide](../../SDL/README.md) and [template work](../../Maintenance/TS1/Plan.md).
