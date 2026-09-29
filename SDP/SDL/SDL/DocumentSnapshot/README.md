# SDL and SDUI documentation snapshot

| Field | Value |
| --- | --- |
| Ecosystem | SDL |
| System | DocumentSnapshot |
| Model role | Bounded capability and product-boundary catalog |
| Profile | design-core 0.5 |
| Entry | [System.design](System.design) |
| Management | [PLAN-SDP-0008](../../../03--Architecture/Ecosystems/Plan.md), E2-M2 |

## Boundary and current status

Implemented integration command sdl-document and snapshot library. It consumes SDUI services and optionally SDL viewpoints; it is not the SDUI renderer implementation itself.

A temporary SDUI session applies supplied accepted state, then the production layout/presentation path generates documentation. State JSON is caller-provided, not proof that domain callbacks accepted those values in a live system. Output is a static snapshot, not interactive Markdown.

The model's `ExistingCapability` activity means implementation was found in the
mapped source, not that behavior was retested for this catalog. `IndependentPackaging`
is planned product-boundary work; it does not claim a new executable, module,
release pipeline or cross-system source linker. Software-system identity lives in
this directory and catalog metadata because design-core has no `system` keyword.

## Responsibility and source map

| Unit | Responsibility | Current source |
| --- | --- | --- |
| SnapshotStateAdapter | `ApplyAcceptedWidgetState` — Validate an explicit value/label/enabled/visible state overlay against known widgets. | [snapshot/snapshot.go](../../../../SDL/go/snapshot/snapshot.go) |
| SnapshotComposer | `RenderUiSnapshot` — Use production SDUI geometry and SVG projection without invoking domain callbacks. | [snapshot/snapshot.go](../../../../SDL/go/snapshot/snapshot.go) |
| DesignNavigationAdapter | `AttachOptionalDesignNavigation` — Optionally validate a structural design and attach its navigation package. | [snapshot/snapshot.go](../../../../SDL/go/snapshot/snapshot.go) |
| SnapshotPublisher | `PublishSnapshotBundle` — Publish a provenance-bearing owned Markdown/SVG bundle. | [cmd/sdl-document/main.go](../../../../SDL/go/cmd/sdl-document/main.go) |

Current libraries: [`snapshot`](../../../../SDL/go/snapshot/).

Current command entry points: [`cmd/sdl-document`](../../../../SDL/go/cmd/sdl-document/).

## Contract and authority limits

No SDL callback, native window or simulation is run by snapshot generation. This cross-ecosystem integration remains under SDL because current implementation and CLI live there; future packaging can revisit ownership without copying the SDUI renderer.

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
/tmp/sdl-catalog check SDP/SDL/SDL/DocumentSnapshot/System.design
/tmp/sdl-catalog ast SDP/SDL/SDL/DocumentSnapshot/System.design
/tmp/sdl-catalog viewpoints SDP/SDL/SDL/DocumentSnapshot/System.design --format static --output /tmp/sdl-documentsnapshot-views
```

The entry is one self-contained model. Generated views belong outside this source
folder. The ecosystem [index](../README.md) records command coverage and catalog
validation. No automatic SDPTool or daemon registration is claimed.
