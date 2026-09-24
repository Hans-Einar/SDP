# Checkpoint #1 — generated G1–G5 design

**Updated implementation status: [supplement 11](11-Go-Implementation-and-Navigation.md).** Dated design/V-phase foundations remain below; Python commands are historical.

Date: 2026-09-22. V2–V4 delivered for review before Go implementation. Counts describe V4; the living model later gained [G6](10-SDL-Viewpoint-Navigation.md). **Start with the [generated implementation report](../../../SDUI/design/viewpoints/implementation.md).** SDL parser/validator/projector derive phases, responsibilities, dependencies and scenarios from [source](../../../SDUI/design/architecture.design), not manual report additions.

## Review scope

Five planned phases, 18 milestones and links for all 94 Functionalities cover frontend, layout/presentation, UI runtime/reload, SDL runtime/binding and native generation. All statuses were **planned**. Reports identify logical owners/source IDs.

The [combined report](../../../SDUI/design/viewpoints/viewpoints.md) contains 11 viewpoints: goals/contributions, architecture, responsibilities, ports, modes, work plan, allocation, sequences, data/contracts, packets and facts. Mermaid/individual diagram sources generate together.

Ten scenarios cover compilation, static/interactive presentation, unbound Go actions, accepted/rejected bindings, reload and native builds. 47 [MessageSet entries](../../../SDUI/design/viewpoints/message-sets.json) derive from Channels, modes, permits and participation.

## Delivery and verification

| Phase | Delivered capability |
| --- | --- |
| V2 | Typed data/fields, persistent Database, projections, explicit packets |
| V3 | Channel roles/permits, ordered scenarios, correlation, derived MessageSet |
| V4 | Traceable G1–G5 plan, generic implementation report, complete design printout |

SDL **design-core 0.5** and SDUI **0.2** have independent versions. No replaced-profile fallback; broader checkpoint/MVP1 candidates remain unadopted.

[Machine report](../../../SDL/tools/verification.json): 61 SDL + 24 tool + 36 SDUI = **121 passing tests**; 368 declarations, 1106 facts, 142 SVGs. Checks cover source positions/facts, SVG labels, bit ranges, sequence order/correlation, hashes and identical re-export. Phase/parser/binding figures visually spot-checked; no physical print/paginated PDF test.

Historical command, requiring the original pre-Go revision:

```sh
python3 SystemDesignLanguage/tools/verify_design.py --phase V4 --renderer /home/warloc/git/mermaid-rs-renderer/target/debug/mmdr
```

Uses an existing renderer; changed neither Mermaid nor XFMD and added no parser to them.

## Limits before G phases

This delivery contained no Go parser/runtime. Scenarios are validated design paths, not executed domain logic. Tokens/AST/normalized/prepared-frame boundary fields are opaque bytes; Go types/signatures, atomic publication, state migration, widget lifetime and font/measurement units awaited G phases.

VP07 reports 15 missing allocations, so this is not a complete deployment plan despite milestone coverage. Packets use an explicit prototype profile, not adopted runtime ABI or P1000/StanForD format. Database means on-demand persistent data, without SQL requirement.

[Branch stack](../../Development-Branch-Stack.md): three commits per V2–V4 phase, push after each, combined PR against sdp-vNow. Merge/start of G phases were outside this delivery.
