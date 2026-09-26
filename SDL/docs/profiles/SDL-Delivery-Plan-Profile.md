# SDL — traceable delivery plans

Introduced in V4, this profile is part of **design-core 0.5** in Go. Four planning relations and an Activity property identify responsibilities an activity should deliver. A source assertion does not prove completion. [Go entry points](../../go/README.md); [G4 evidence](../../go/evidence/G4.md).

| Statement | Signature and meaning |
| --- | --- |
| `A addresses F.` | Activity → Functionality; activity covers this implementation responsibility |
| `A delivers F.` | Activity → Feature; planned contribution, not proof of a complete Feature |
| `A depends-on B.` | Activity → Activity; explicit prerequisite, without self-reference/cycles |
| `S illustrates A.` | Scenario → Activity; explicit example path linked to a milestone |
| `A has implementation-status = planned/implemented/verified.` | Activity property; a source assertion, not status inferred from tests |

G phases and milestones are Activity identities. Each milestone `refines` its phase. This is the example's display convention; the generator knows neither G names nor phase counts. It uses relations, not prefixes or a built-in plan. `refines` does not inherit status or dependencies. Scheduling, duration, resources and automatic work execution are outside this profile.

VP06 shows activity decomposition, deliveries, responsibilities and a separate dependency graph. `implementation.md` is generated from the same facts, with source IDs and ownership links. It reports uncovered Functionality elements without inventing milestones. G1–G6 have since been delivered within profile boundaries. The current [design source](../../../SDUI/design/architecture.design) owns explicit model status assertions; the [phase plan](../../../SDUI/docs/implementation-plan.md) links implementation evidence. The original V4 planned snapshot remains below as history.

## Historical V4 delivery — 2026-09-22

These counts/statuses describe this milestone before the Go port, not newly run tests or current overall implementation status.

V4 milestones: M1 traceable phase/responsibility model; M2 generated reader report and complete parser/presentation/binding/reload scenarios; M3 combined verification, checkpoint, phase push and review PR against sdp-vNow.

V4-M1 verified: 61 parser tests pass. The model declares five planned G phases, 18 milestones and responsibility links for all 94 Functionality elements. No future Go-code status was upgraded to implemented/verified.

V4-M2 delivered: generic implementation.md from phase/milestone relations, with responsibility owners and scenario diagrams. 24 tool tests pass. Ten model scenarios cover compilation, interactive/static presentation, local Go actions, SDL binding, UI/SDL reload and native builds. These are design paths, not executed Go programs.

V4-M3 verified: 61 SDL parser, 24 tool and 36 SDUI tests pass. 142 SVG diagrams from 368 declarations and 1106 facts checked for source links and byte-identical re-export. Fifteen incomplete mode allocations remain explicitly reported. Phase graph, parser and binding sequences visually spot-checked; neither Go runtime nor physical printing verified.
