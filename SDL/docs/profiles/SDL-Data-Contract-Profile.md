# SDL — bounded data and wire profile

Introduced in V2 (0.3), the data/wire profile is now part of **design-core 0.5** in Go. The following rules are active; the broader data study remains a candidate. Structural facts are separate statements. This does not implement binary serialization. [Go entry points](../../go/README.md); [G4 evidence](../../go/evidence/G4.md).

New declaration kinds: `dataset`, `database`, `datagram`, `contract`, `variant`, `field`, `encoding`. All use `kind Name.` and one shared symbol table.

| Statement | Signature and rule |
| --- | --- |
| `D upholds C.` | Dataset/Datagram → Contract; exactly one contract per data definition |
| `M from D.` | Datagram → Dataset; exactly one logical source |
| `U owns B.` | Unit → Database; one immediate responsible owner, as for Functionality |
| `U holds D.` | Unit/Database → Dataset; one immediate holder in this definition profile |
| `C defines V.` | Contract → Variant; one contract owner per variant |
| `C has-field F.` | Contract/Variant → Field; one immediate field owner |
| `P projects D into M.` | Functionality × Dataset × Datagram; must match M's source |
| `E encodes V.` | Encoding → Variant; exactly one variant per encoding |
| `E places F at 0 bits 16.` | Encoding × Field × offset × bit width; explicit nonnegative integers, positive width |

`Contract has completeness = open/closed.` is mandatory. Open may describe incomplete work but produces no packet. A closed Dataset contract has at least one field and no variants. A closed Datagram contract has at least one variant; contract-level common fields are header fields, and each variant has at least one payload field. Multiple Datasets/families may reuse a contract when shape requirements match.

Every Field specifies `value-type = unsigned/signed/boolean/text/bytes/decimal` and `presence = required/optional`. Absence is not null, empty text, false or zero. These are shape contracts; units, value domains, transformations, version negotiation and execution of projections are not implemented.

Database means a place from which persistent data can be retrieved when needed. SQL, process boundaries, transactions and a particular storage engine are not inferred. A Dataset may be transient in a Unit; persistence requires an explicit Database holder. Definitions are not executable instances or a multiple-instance/deployment model.

Encoding requires a closed contract and explicit `byte-order = big-endian/little-endian` and `bit-order = most-significant-first/least-significant-first`. Place every header/variant field exactly once, contiguously from offset 0, without gaps/overlap. Bit positions are consecutive contract positions; labels are not encoder code. Boolean requires 1 bit; signed/unsigned allow at most 64; bytes require multiples of 8. Text, decimal and optional fields are rejected by this fixed packet profile. Maximum layout size: 65536 bits. Extensions must explicitly define padding, variable fields and presence encoding; the generator infers none.

Core EBNF additions:

```text
projection = identifier, "projects", identifier, "into", identifier, "." ;
placement = identifier, "places", identifier, "at", integer, "bits", integer, "." ;
integer = "0" | nonzeroDigit, { digit } ;
```

AST: `Projection(subject, dataset, datagram, span)` and `Placement(subject, field, offset, width, span)`; numbers have `Integer(value, span)` nodes. Other statements reuse Relation/PropertyAssignment. VP09 provides data/contract maps; VP10 provides packets for validated Encodings. Source maps include the projection's third argument and every explicit field placement.

## Historical V2 delivery — 2026-09-22

These counts/statuses describe this milestone before the Go port, not newly run tests or current overall implementation status.

Milestones: V2-M1 language/validation; V2-M2 generic VP09/VP10 and tests; V2-M3 ported example, regeneration and verification before phase push.

V2-M1 verified: 49 parser/validator/CLI tests pass. Later milestones deliver viewpoints and regenerated shared artifacts. 0.3 replaces active 0.2.

V2-M2: VP09/VP10 and source maps implemented; 17 tool tests pass. Packets use explicit bit ranges following [Mermaid packet](https://mermaid.js.org/syntax/packet.html), also checked against local mmdr; this is not binary serialization.

V2-M3 delivered: 227 declarations, 509 facts and 73 SVG diagrams. 49 SDL, 17 viewpoint and 36 SDUI tests pass; 15 partial allocation gaps remain visible. The tool generates packet/data examples from the shared SDL source.
