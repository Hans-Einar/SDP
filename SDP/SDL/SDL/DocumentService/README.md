# SDL on-demand document service

| Field | Value |
| --- | --- |
| Ecosystem | SDL |
| System | DocumentService |
| Model role | Bounded capability and product-boundary catalog |
| Profile | design-core 0.5 |
| Entry | [System.design](System.design) |
| Management | [PLAN-SDP-0008](../../../03--Architecture/Ecosystems/Plan.md), E2-M2 |

## Boundary and current status

Implemented Linux document-service and request-client executables. It is an optional local broker, not the future request/routine governance state machine.

The source registry is explicit. Client/window/pane sequence numbers reject superseded requests; the broker rechecks source revision before opening a document. Store leases survive daemon restarts; ambiguous opens retain leases conservatively. Runtime directory files are real files, not a promise of disk-free memory hosting.

The model's `ExistingCapability` activity means implementation was found in the
mapped source, not that behavior was retested for this catalog. `IndependentPackaging`
is planned product-boundary work; it does not claim a new executable, module,
release pipeline or cross-system source linker. Software-system identity lives in
this directory and catalog metadata because design-core has no `system` keyword.

## Responsibility and source map

| Unit | Responsibility | Current source |
| --- | --- | --- |
| ViewRequestClient | `SubmitViewSelection` — Send select/release/sweep calls over the local broker socket. | [cmd/sdl-view-request/main.go](../../../../SDL/go/cmd/sdl-view-request/main.go) |
| ProjectionBroker | `ResolveCurrentView` — Resolve a registered project, validate a bounded query, generate or reuse a revision-specific bundle. | [broker/broker.go](../../../../SDL/go/broker/broker.go) |
| LeaseStore | `RetainPublishedView` — Persist immutable packages and leases; sweep only released packages. | [broker/store.go](../../../../SDL/go/broker/store.go) |
| ReaderAdapter | `OpenCapturedReaderTarget` — Open the published entry through an explicitly registered reader with window/pane identity. | [reader/xfmd.go](../../../../SDL/go/reader/xfmd.go) |

Current libraries: [`broker`](../../../../SDL/go/broker/), [`reader`](../../../../SDL/go/reader/).

Current command entry points: [`cmd/sdl-viewsd`](../../../../SDL/go/cmd/sdl-viewsd/), [`cmd/sdl-view-request`](../../../../SDL/go/cmd/sdl-view-request/).

## Contract and authority limits

The abstract ViewCalls protocol below is a selected request/result view of broker.Call and Result, not a complete JSON wire schema. Release and sweep are modeled as responsibilities in prose; this scenario focuses on successful selection. No daemon or XFMD process was launched for this catalog.

The channel payload contracts are closed **for this selected metadata projection**.
They deliberately omit application values and other API fields. They must not be
used to generate or validate the real Go API/JSON wire contract; the modeled
scenario checks only this declared projection, not a live transaction.

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
/tmp/sdl-catalog check SDP/SDL/SDL/DocumentService/System.design
/tmp/sdl-catalog ast SDP/SDL/SDL/DocumentService/System.design
/tmp/sdl-catalog viewpoints SDP/SDL/SDL/DocumentService/System.design --format static --output /tmp/sdl-documentservice-views
```

The entry is one self-contained model. Generated views belong outside this source
folder. The ecosystem [index](../README.md) records command coverage and catalog
validation. No automatic SDPTool or daemon registration is claimed.
