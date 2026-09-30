# SDL frontend

| Field | Value |
| --- | --- |
| Ecosystem | SDL |
| System | Frontend |
| Model role | Bounded capability and product-boundary catalog |
| Profile | design-core 0.5 |
| Entry | [System.design](System.design) |
| Management | [PLAN-SDP-0008](../../../03--Architecture/Ecosystems/Plan.md), E2-M2 |

## Boundary and current status

Implemented Go library and multipurpose CLI. A separately packaged frontend binary is a target boundary, not an additional executable delivered here.

check, ast, format, action-check and class-check are current sdl commands. class-view and structural viewpoints are owned by the Viewpoints catalog boundary. The parser package currently contains structural, executable-action and class profiles; splitting these implementations is not part of this model.

The model's `ExistingCapability` activity means implementation was found in the
mapped source, not that behavior was retested for this catalog. `IndependentPackaging`
is planned product-boundary work; it does not claim a new executable, module,
release pipeline or cross-system source linker. Software-system identity lives in
this directory and catalog metadata because design-core has no `system` keyword.

## Responsibility and source map

| Unit | Responsibility | Current source |
| --- | --- | --- |
| StructuralParser | `ParseStructuralModel` — Build source-positioned declarations and statements for design-core 0.5. | [parser/parse.go](../../../../SDL/go/parser/parse.go) |
| SemanticChecker | `CheckStructuralRelations` — Check names, typed relations, contracts, channels and scenario constraints. | [parser/validate.go](../../../../SDL/go/parser/validate.go) |
| CanonicalFormatter | `FormatStructuralSource` — Emit canonical structural source after validation. | [parser/ast.go](../../../../SDL/go/parser/ast.go) |
| ActionCompiler | `CompileActionModel` — Compile action-core 0.1 into typed records and registered action signatures; do not execute source. | [parser/actions.go](../../../../SDL/go/parser/actions.go) |
| ClassChecker | `CheckClassModel` — Check explicit class-core declarations without inferring classes from containment. | [parser/classes.go](../../../../SDL/go/parser/classes.go) |

Current libraries: [`parser`](../../../../SDL/go/parser/).

Current command entry points: [`cmd/sdl`](../../../../SDL/go/cmd/sdl/).

## Contract and authority limits

No source-set imports, ecosystem declarations or full experimental MVP1 profile are implemented. Untracked sourceinput work is excluded.

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
/tmp/sdl-catalog check SDP/SDL/SDL/Frontend/System.design
/tmp/sdl-catalog ast SDP/SDL/SDL/Frontend/System.design
/tmp/sdl-catalog viewpoints SDP/SDL/SDL/Frontend/System.design --format static --output /tmp/sdl-frontend-views
```

The entry is one self-contained model. Generated views belong outside this source
folder. The ecosystem [index](../README.md) records command coverage and catalog
validation. No automatic SDPTool or daemon registration is claimed.
