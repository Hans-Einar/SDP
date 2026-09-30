# SDL System and explicit source-set contract

| Field | Value |
| --- | --- |
| id | KB-SDL-005 |
| project | SDL |
| type | Change |
| CardState | ready |
| created | 2026-09-25T11:42:10Z |
| source | SCRUM-SDP-0001; owner-conversation-2026-09-25 |
| ScrumId | SCRUM-SDP-0002 |
| Systems | SDL |
| next_review | Session 0001 S2: select ImplementationPlan from SSD1 design |
| PlanId | PLAN-SDP-0010 |

## Bounded delivery

Define and implement an explicit, versioned System/source-set contract in Go.
SDL, SDUI and SDPTool are separate modeled Systems; the shared SDP project can
manage all three. Resolve system membership, source membership, public cross-system
references, parser/tool boundary, canonical formatting, source provenance and
revision invalidation before moving live models. Current design-core 0.5 has no
System keyword; do not silently reinterpret it or introduce a second parser.

Inspect the historical source-tree study and unfinished SDL/go/sourceinput draft.
The draft is uncommitted/untested, not an accepted contract. Keep pure parsing
I/O-free; test duplicate/unresolved names, invalid membership, cycles where relevant,
limits, file diagnostics, tool projection and source-change handling. Independent
system entries and cross-system landscape references need explicit semantics.
Migrate consumers with the selected profile; preserve frozen evidence.

## Historical queue — superseded by Scrum-0002

CardState is queued after MAINT-SDP-0001: source semantics must be settled before
KB-SDP-020 moves the mixed model. First milestone: a bounded profile/acceptance
contract checked against existing Go consumers. No sprint is selected by this Scrum.

## Lineage

KBO-SDP-000003 splits [KB-SDL-004](../superseded/%23004--SDL--Change--Language-source-organization.md).
This card owns language/input/tool behavior; [KB-SDP-020](../backlog/%23020--Change--Shared-design-source-organization.md)
owns the subsequent document/model migration. Nothing from the original scope is
claimed implemented by this split. Work in a phase branch with milestone evidence.

## Scrum-0001 review

Queue the explicit System/input contract after this maintenance. A draft loader does not deliver it.

EVT-KB-SDL-000025; next review at the next selection or relevant dependency delivery.

## Scrum-0002 disposition

2026-09-25T12:38:14Z — EVT-KB-SDL-000026: backlog, no longer queued. KB-SDP-021 is next.
Existing single-file SDL services allow the first SDPTool Sprint to proceed;
future source-set support still depends on this contract before model migration.

## MVP1 pilot follow-up — 2026-09-27

The owner-authorized MVP1 pilot (PLAN-SDP-0005) now preserves 4,523 statements across 68 explicitly listed files with separate UI containers. Full navigation still requires System/source-set linking AND decisions on the candidate exercise profile; source assembly alone is insufficient. Revisions must cover all inputs. Do not use the inventory auditor as an SDL parser.

[Pilot and navigation handoff](../../../experiments/mvp1_sdl/Navigation.md).

## Ecosystem modeling outcome — 2026-09-29

[PLAN-SDP-0008](../../03--Architecture/Ecosystems/Plan.md) now provides 17 small independent system models and an integration context model using current design-core 0.5. Real System/source-set linking is still absent. This corpus is a useful additional acceptance input: preserve standalone diagnostics/navigation, link declared public boundaries explicitly, and permit coherent per-system file splitting without an ad hoc preprocessor. The untracked sourceinput draft remains untouched and unadopted.

## XFMD gap study input — 2026-09-29

[Study and proposed acceptance](../../02--Requirements/XFMD-Gaps/Study.md).

SDL-001: retain explicit membership, original-source diagnostics and all-input revisions; verify both checker and navigation against the same source graph. The untracked draft remains unadopted.

EVT-KB-SDL-000030: CardState remains backlog; this review selects no implementation.

## Session proposal — 2026-09-30

[KB042](../backlog/%23042--Proposal--Goal-oriented-sessions-and-roadmaps.md) and
[proposed SDL Session](../../Sessions/session-%230001--SDL_expansion.md).

The provisional SDL-expansion Session links this primary card to proposed design, implementation and verification steps. It is a roadmap, not execution authorization.

EVT-KB-SDL-000033: Existing CardState and execution selection unchanged.

## Session S1 execution — 2026-09-30

The owner selected the next Session step. Execute
[PLAN-SDP-0010](../../04--Design/SDL/SourceSets/Plan.md) on the design phase branch.
CardState is in-progress while defining the bounded System/source-set contract.
Language implementation and model migration remain later steps, not delivered.

EVT-KB-SDL-000034 records activation; the previous proposal notes are history.

## SSD1-M1 handoff — 2026-09-30

[Contract](../../04--Design/SDL/SourceSets/Contract.md),
[acceptance cases](../../04--Design/SDL/SourceSets/Acceptance.md) and
[evidence](../../04--Design/SDL/SourceSets/Evidence.md) complete Session S1.
Recommended initial delivery: design-core/0.6 System scope and explicit
sdl-source-set/0.1 inputs; preserve current standalone 0.5. No implementation exists
yet. Public cross-System linking/exports are deliberately not smuggled into file
assembly; they remain this card's obligation until implemented or explicitly split
with preserved lineage. KB-SDP-020 migration remains deferred.

CardState ready is the recorded handoff to Session S2 implementation planning.
This card remains active because its selected delivery is unfinished, not merely
to point at unrelated future work. EVT-KB-SDL-000035 records this handoff.
