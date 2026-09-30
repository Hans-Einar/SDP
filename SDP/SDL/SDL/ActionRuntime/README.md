# SDL action runtime and SDUI binding

| Field | Value |
| --- | --- |
| Ecosystem | SDL |
| System | ActionRuntime |
| Model role | Bounded capability and product-boundary catalog |
| Profile | design-core 0.5 |
| Entry | [System.design](System.design) |
| Management | [PLAN-SDP-0008](../../../03--Architecture/Ecosystems/Plan.md), E2-M2 |

## Boundary and current status

Implemented reusable Go libraries embedded in application hosts. There is no standalone production SDL runtime daemon or generic runtime executable in the current command inventory.

The model uses units without a fabricated runtime container. The examples/application package composes SDL and SDUI runtimes; sdl-simulate, sdl-demo and compiled variants are sample hosts. Domain state belongs to registered Go code. Request sequencing is per engine; consumed commands are not automatically replayed.

The model's `ExistingCapability` activity means implementation was found in the
mapped source, not that behavior was retested for this catalog. `IndependentPackaging`
is planned product-boundary work; it does not claim a new executable, module,
release pipeline or cross-system source linker. Software-system identity lives in
this directory and catalog metadata because design-core has no `system` keyword.

## Responsibility and source map

| Unit | Responsibility | Current source |
| --- | --- | --- |
| ActionEngine | `InvokeRegisteredHandler` — Validate revision, sequence and typed input; call registered Go logic and validate output. | [runtime/engine.go](../../../../SDL/go/runtime/engine.go) |
| ActionReloader | `InstallCompatibleActionModel` — Validate a replacement program against existing registrations and publish a new model revision. | [runtime/reload.go](../../../../SDL/go/runtime/reload.go) |
| UiBindingBridge | `BindUiEventsToActions` — Validate explicit SDUI bindings and translate typed widget events/results through engine signatures. | [bridge/bind.go](../../../../SDL/go/bridge/bind.go) |
| ActionSourceWatcher | `ObserveActionCandidates` — Publish source candidates for validation; malformed updates do not replace the accepted model. | [reload/actions.go](../../../../SDL/go/reload/actions.go) |

Current libraries: [`runtime`](../../../../SDL/go/runtime/), [`bridge`](../../../../SDL/go/bridge/), [`reload`](../../../../SDL/go/reload/).

Current command entry points: No standalone production binary; see host fixtures in the ecosystem index..

## Contract and authority limits

action-core 0.1 executes explicitly registered handlers, not arbitrary structural design-core statements. Callback symbols require signatures and host registration. The abstract ActionCalls request/result below captures action, revision and sequence, not a complete serialized ABI. A stable cross-language ABI remains outside current evidence.

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
/tmp/sdl-catalog check SDP/SDL/SDL/ActionRuntime/System.design
/tmp/sdl-catalog ast SDP/SDL/SDL/ActionRuntime/System.design
/tmp/sdl-catalog viewpoints SDP/SDL/SDL/ActionRuntime/System.design --format static --output /tmp/sdl-actionruntime-views
```

The entry is one self-contained model. Generated views belong outside this source
folder. The ecosystem [index](../README.md) records command coverage and catalog
validation. No automatic SDPTool or daemon registration is claimed.
