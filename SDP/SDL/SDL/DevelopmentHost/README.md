# SDL development process host

| Field | Value |
| --- | --- |
| Ecosystem | SDL |
| System | DevelopmentHost |
| Model role | Bounded capability and product-boundary catalog |
| Profile | design-core 0.5 |
| Entry | [System.design](System.design) |
| Management | [PLAN-SDP-0008](../../../03--Architecture/Ecosystems/Plan.md), E2-M2 |

## Boundary and current status

Implemented development CLI and manager library. This host rebuilds and restarts a selected application; it does not perform native Go hot swapping.

Changing SDL/SDUI models within an already running application is a distinct library-level reload concern. sdl-dev watches Go changes and replaces the process. A compile failure retains the old child; a later process-start failure does not guarantee rollback to that child.

The model's `ExistingCapability` activity means implementation was found in the
mapped source, not that behavior was retested for this catalog. `IndependentPackaging`
is planned product-boundary work; it does not claim a new executable, module,
release pipeline or cross-system source linker. Software-system identity lives in
this directory and catalog metadata because design-core has no `system` keyword.

## Responsibility and source map

| Unit | Responsibility | Current source |
| --- | --- | --- |
| SourceWatcher | `WatchGoSourceChanges` — Watch a selected Go source root for changes. | [devhost/watch.go](../../../../SDL/go/devhost/watch.go) |
| BuildCoordinator | `BuildCandidateApplication` — Invoke an explicit Go toolchain/package/build-tags selection. | [devhost/manager.go](../../../../SDL/go/devhost/manager.go) |
| ProcessManager | `RestartBuiltApplication` — Restart after successful build; retain the running app on build failure. | [devhost/manager.go](../../../../SDL/go/devhost/manager.go) |

Current libraries: [`devhost`](../../../../SDL/go/devhost/).

Current command entry points: [`cmd/sdl-dev`](../../../../SDL/go/cmd/sdl-dev/).

## Contract and authority limits

In-memory domain state is not promised across process restart. Process-group details differ on Linux and other platforms. No new reload supervisor or persistent restart store is delivered.

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
/tmp/sdl-catalog check SDP/SDL/SDL/DevelopmentHost/System.design
/tmp/sdl-catalog ast SDP/SDL/SDL/DevelopmentHost/System.design
/tmp/sdl-catalog viewpoints SDP/SDL/SDL/DevelopmentHost/System.design --format static --output /tmp/sdl-developmenthost-views
```

The entry is one self-contained model. Generated views belong outside this source
folder. The ecosystem [index](../README.md) records command coverage and catalog
validation. No automatic SDPTool or daemon registration is claimed.
