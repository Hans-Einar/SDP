# SDL Go code generation

| Field | Value |
| --- | --- |
| Ecosystem | SDL |
| System | CodeGeneration |
| Model role | Bounded capability and product-boundary catalog |
| Profile | design-core 0.5 |
| Entry | [System.design](System.design) |
| Management | [PLAN-SDP-0008](../../../03--Architecture/Ecosystems/Plan.md), E2-M2 |

## Boundary and current status

Implemented sdl-gen command and codegen library. SDL code generation is distinct from compiling arbitrary structural descriptions into a finished domain system.

Generate produces Program constructors; Bundle combines actions_gen.go and ui_gen.go. Handwritten Go domain functions remain external registrations. The shared documents publication helper is reused, not owned a second time by this catalog.

The model's `ExistingCapability` activity means implementation was found in the
mapped source, not that behavior was retested for this catalog. `IndependentPackaging`
is planned product-boundary work; it does not claim a new executable, module,
release pipeline or cross-system source linker. Software-system identity lives in
this directory and catalog metadata because design-core has no `system` keyword.

## Responsibility and source map

| Unit | Responsibility | Current source |
| --- | --- | --- |
| ActionEmitter | `EmitActionConstructors` — Validate action source and emit formatted Go constructors with source digest. | [codegen/generate.go](../../../../SDL/go/codegen/generate.go) |
| UiGenerationAdapter | `RequestUiConstructors` — Delegate SDUI constructors to the SDUI codegen package. | [codegen/generate.go](../../../../SDL/go/codegen/generate.go) |
| GenerationPublisher | `PublishOwnedGeneration` — Publish a manifest-owned bundle of generated files while retaining foreign files. | [documents/publish.go](../../../../SDL/go/documents/publish.go) |

Current libraries: [`codegen`](../../../../SDL/go/codegen/), [`documents`](../../../../SDL/go/documents/).

Current command entry points: [`cmd/sdl-gen`](../../../../SDL/go/cmd/sdl-gen/).

## Contract and authority limits

No arbitrary-Go interpreter, whole-system compiler or new ABI is implied. The Go compiler remains an external tool. Generated constructors still depend on runtime/parser types.

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
/tmp/sdl-catalog check SDP/SDL/SDL/CodeGeneration/System.design
/tmp/sdl-catalog ast SDP/SDL/SDL/CodeGeneration/System.design
/tmp/sdl-catalog viewpoints SDP/SDL/SDL/CodeGeneration/System.design --format static --output /tmp/sdl-codegeneration-views
```

The entry is one self-contained model. Generated views belong outside this source
folder. The ecosystem [index](../README.md) records command coverage and catalog
validation. No automatic SDPTool or daemon registration is claimed.
