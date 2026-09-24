# SDL symbol profile 1

G6-M5: each figure carries source identity/type; each edge carries relation and source fact. Mermaid-rs-renderer still owns placement/routing. SDL's presentation adapter reads its existing `--dumpLayout` contract and draws consistent SVG symbols using the same node/edge geometry. It copies neither parser nor layout engine. Backend node IDs, edges and geometry are checked before drawing.

Actor: human figure with head, torso, arms and legs. UseCase: ellipse. Feature: tabbed card. Functionality: rounded card with side marker. Capability: hexagon. Activity: rounded activity card. Mode: dashed context frame. Container: double border. Unit: single border. Interface: labeled port card. Database: cylinder representing a persistent data source. Other data/contract objects use cards with explicit type labels. Color alone does not encode type.

`consumes` uses a dashed dependency with an open arrow toward Interface. Other structural relations use labeled SDL arrows. `realizes` still means contribution to realization and receives no UML hollow triangle. `contains`/`owns` never become composition. Channel-contract `uses` is participation/role, not UML Usage. No System boundary or include/extend relation is inferred.

SVG is the symbol profile's authoritative image. Mermaid source remains available as a portable, labeled structural view; it does not claim identical figures in every reader. Sequence and packet use existing tested backend forms. Native `usecase-beta` was not selected: an earlier local backend probe produced incorrect figures despite exit code 0. Class diagrams require their own explicit profile.

Basis: [Mermaid flowchart](https://mermaid.js.org/syntax/flowchart.html), [Mermaid classDiagram](https://mermaid.js.org/syntax/classDiagram) and [project notation semantics](../integration/SDL-Viewpoint-Levels-and-Notation.md).
