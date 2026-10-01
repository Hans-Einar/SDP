# SDL viewpoints and document projection

| Field | Value |
| --- | --- |
| Ecosystem | SDL |
| System | Viewpoints |
| Model role | Bounded capability and product-boundary catalog |
| Profile | design-core 0.5 |
| Entry | [System.design](System.design) |
| Management | [PLAN-SDP-0008](../../../03--Architecture/Ecosystems/Plan.md), E2-M2 |

## Boundary and current status

Implemented projection libraries reached through sdl viewpoints, view and class-view. A separate projection binary is an intended product boundary; the current executable is shared with Frontend.

Navigation export is the default and does not materialize all details. Static export and selected view requests produce source-derived content. Mermaid text generation and optional SVG backend invocation are separate operations. documents/publish.go is the single reusable owned-file publisher.

The model's `ExistingCapability` activity means implementation was found in the
mapped source, not that behavior was retested for this catalog. `IndependentPackaging`
is planned product-boundary work; it does not claim a new executable, module,
release pipeline or cross-system source linker. Software-system identity lives in
this directory and catalog metadata because design-core has no `system` keyword.

## Responsibility and source map

| Unit | Responsibility | Current source |
| --- | --- | --- |
| ViewpointModel | `BuildValidatedViewFacts` — Validate structural source and derive the eleven viewpoint families. | [viewpoint/model.go](../../../../SDL/go/viewpoint/model.go) |
| ViewSelector | `SelectBoundedView` — Validate sdl-view URI selection, relation direction, depth, viewpoint and focus. | [viewpoint/query.go](../../../../SDL/go/viewpoint/query.go) |
| DocumentProjector | `BuildNavigationOrStaticBundle` — Generate navigation indexes or selected/static Markdown and Mermaid bundles with provenance. | [documents/bundle.go](../../../../SDL/go/documents/bundle.go) |
| ClassProjector | `ProjectExplicitClasses` — Generate class documentation from the separate explicit class profile. | [documents/classes.go](../../../../SDL/go/documents/classes.go) |
| DiagramAdapter | `RenderSelectedDiagrams` — Optionally use registered mmdr and semantic diagram decoration. | [documents/renderer.go](../../../../SDL/go/documents/renderer.go) |

Current libraries: [`viewpoint`](../../../../SDL/go/viewpoint/), [`documents`](../../../../SDL/go/documents/).

Current command entry points: [`cmd/sdl`](../../../../SDL/go/cmd/sdl/).

## Contract and authority limits

Passing model checks proves the declared model is supported, not that viewpoints cover every intended design fact. Folder placement does not register sources with SDPTool or a running broker.

Ports express dependencies, not complete transport contracts. Reusing a helper does
not transfer its implementation ownership. Modes express operating contexts, not
an implemented application state machine. This standalone catalog introduces no
cross-file symbol resolution. Full older design authority remains at
[the combined architecture model](../../../../SDUI/design/architecture.design)
and [the shared implementation plan](../../../../SDUI/docs/implementation-plan.md);
implementation observations should be read against the current mapped sources.
This overview neither replaces nor mechanically extracts that detailed model.

## Reproduce model projections

From the repository root, with the Go toolchain on PATH:

```sh
(cd SDL/go && go build -o /tmp/sdl-catalog ./cmd/sdl)
/tmp/sdl-catalog check SDP/SDL/SDL/Viewpoints/System.design
/tmp/sdl-catalog ast SDP/SDL/SDL/Viewpoints/System.design
/tmp/sdl-catalog viewpoints SDP/SDL/SDL/Viewpoints/System.design --format static --output /tmp/sdl-viewpoints-views
```

The entry is one self-contained model. Generated views belong outside this source
folder. The ecosystem [index](../README.md) records command coverage and catalog
validation. No automatic SDPTool or daemon registration is claimed.
