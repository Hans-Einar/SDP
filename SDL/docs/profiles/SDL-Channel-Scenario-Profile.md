# SDL — Channel contracts and declared scenarios

Introduced in V3 (0.4), this profile is now part of **design-core 0.5** in Go. It provides structural protocol validation and sequence diagrams, not message execution. The broader workspace/Channel proposal in the language definition remains a study. [Go entry points](../../go/README.md); [G4 evidence](../../go/evidence/G4.md).

New types: `channel`, `message`, `scenario`. A Channel represents logical collaboration; participants may be Units or Containers. Library calls do not automatically imply network traffic. A Scenario describes an explicitly ordered example path.

| Statement | Meaning |
| --- | --- |
| `C upholds K.` | Channel → Contract; exactly one contract |
| `K permits M.` | Contract → Message/Datagram; permitted message types/families |
| `M upholds K.` | Message → record Contract; an action requires no Dataset |
| `M has message-kind = request/result/event.` | Message role; required for Message |
| `Result replies-to Request.` | Message → Message; one request type per result type |
| `P uses C as sender of M in mode Mode.` | Unit × Channel × sender/receiver × Message/Datagram × Mode |
| `S runs-in Mode.` | Scenario → Mode; exactly one |
| `S exercises Goal.` | Scenario → UseCase; explicit traceability, not execution evidence |
| `S has completeness = closed/open.` | Closed requires responses to all requests; open may show a bounded prefix |

```text
S step 1 sends Request from Client to Server via Service.
S step 2 sends Result from Server to Client via Service reply-to 1.
S step 3 sends Notices variant Changed from Server to Observer via Events.
```

Step ordinals are explicit positive integers, unique and consecutive from 1. Canonical text sorting does not determine time order; the generator uses ordinals. Datagram requires a variant belonging to its family's contract; ordinary Message must not specify a variant. All contracts used must be closed before generating the sequence. Open models may still declare Channels/Messages without executable scenarios.

A step must be permitted by its Channel contract and match sender/receiver participation for its message type in the scenario's explicit Mode. Participants, contracts and modes are not inherited. Self-send is allowed as an explicitly local call. Multiple receivers require separate steps; no broadcast is inferred.

A Result requires reply-to referencing an earlier request in the same scenario, on the same Channel, with reversed participants and matching result-type replies-to. A request step has at most one result. Request/event/Datagram cannot use reply-to. This checks example-path correlation, not runtime IDs, delivery guarantees, timing, reentrancy, transactions or actual execution.

Channel contracts contain permits, not record fields/variants. Message uses record fields. A Contract cannot serve as both Channel contract and payload shape. Every result Message type requires replies-to; other message kinds forbid it. MessageSet is only a derived per-Channel/Mode export from permits and participation; there is no independently authored MessageSet syntax.

EBNF additions:

```text
participation = identifier, "uses", identifier, "as", ("sender" | "receiver"),
                "of", identifier, "in", "mode", identifier, "." ;
step = identifier, "step", integer, "sends", identifier,
       [ "variant", identifier ], "from", identifier, "to", identifier,
       "via", identifier, [ "reply-to", integer ], "." ;
```

AST: `Participation(subject, channel, role, message, mode, span)` and `Step(subject, ordinal, message, variant, sender, receiver, channel, reply_to, span)`. Identifier/Integer nodes carry source positions. Sequence arrows and MessageSet rows carry source IDs for steps, participation, governing contracts and permits.

Model alternatives as separate named scenarios, such as accepted/rejected input. Branch/loop/parallel syntax awaits defined semantics. A diagram is an agreed example path, not the entire permitted protocol.

## Historical V3 delivery — 2026-09-22

These counts/statuses describe this milestone before the Go port, not newly run tests or current overall implementation status.

Milestones: V3-M1 language and negative contract tests; V3-M2 VP08/MessageSet; V3-M3 request/result and reload examples with generated rendering evidence.

V3-M1 verified: 59 parser/contract tests pass, covering requests/results, modes, permissions, correlation, step ordering and required Datagram variants.

V3-M2: VP08 and derived MessageSet implemented; 22 tool tests pass, checking numeric order, correlation and arrow contract evidence. Syntax follows [Mermaid sequence](https://mermaid.js.org/syntax/sequenceDiagram.html); backend layout does not alter message direction/order.

V3-M3 delivered: 273 declarations, 675 facts, four scenarios and 88 SVG diagrams. 59 SDL, 22 viewpoint and 36 SDUI tests pass. Source facts, arrow references, bit ranges and export fingerprints checked; repeated export identical. Accepted-action and rejected-reload diagrams visually spot-checked.
