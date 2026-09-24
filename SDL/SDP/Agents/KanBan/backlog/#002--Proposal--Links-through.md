# Weak links and visual paths through nodes

| Field | Value |
| --- | --- |
| id | KB-SDL-002 |
| CardState | backlog |
| project | SDL |
| type | Proposal |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | 2026-09-30 |

Registered from the owner conversation on 2026-09-23. The timestamp records registration, not a reconstructed time for earlier discussions. The directory and ledger record lifecycle status.

## Owner's clarified intent

`links` is a descriptive relationship; `uses`, `realizes` and other specific relations retain their own semantics. Weak links may appear as dashed lines with arrow direction determined by `to`/`from`. Validate references, but do not treat the link as an obligation, data flow, dependency or implementation evidence.

`through` specifies one continuous visual path through an ordered list of nodes. The path enters a visual port, remains visible over the node, and exits on the opposite side before continuing. This is deliberate traversal, not obstacle avoidance or a line hidden behind a box. It can show a Feature passing through architecture or class blocks. Visual ports do not automatically represent model Channel ports.

## Provisional syntax sketch

```text
Actor actor1 links to stakeholder1
Actor actor2 links through actor1 to stakeholder1
Feature feature1 links through container1, container2 to container3
```

Preserve a single relation with source, direction, ordered intermediate nodes and target in the AST/model. Do not lower it to independent semantic `uses`/`realizes` edges. Identify the complete path through a label and continuous marking, not color alone.

## Open questions and next work

Resolve final grammar, `from` combined with `through`, multiple relations, identity, permitted repeated nodes/cycles and diagnostics. Define port placement by diagram direction, text collisions, crossing paths and partially hidden nodes in viewpoints. An exporter unable to draw traversal must state that limitation. Inspect the existing SDL projector/render adapter before selecting a backend or changing the Mermaid repository; registering this card does not authorize renderer changes.

## Acceptance for later implementation

Parser/AST preserve order and source positions; unknown references produce diagnostics. A figure with at least two intermediate nodes shows entry ports, visible traversal over each node, exit ports on opposite sides and correct arrow direction. Viewpoint selection preserves the path with explicit handling of omitted nodes. Status/coverage calculations ignore the link as satisfaction evidence.
