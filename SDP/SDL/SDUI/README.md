# SDUI ecosystem

This ecosystem groups the UI language toolchain and host products. Each child is
a software-system boundary with its own independently parseable SDL entry. The
owner selected separate tool systems on 2026-09-29; current implementation still
shares the `SDUI/go` module and several executable entrypoints. This catalog is an
architecture proposal grounded in implemented packages, not a completed extraction.

| System | Responsibility | Current binary/library mapping | Extraction status |
| --- | --- | --- | --- |
| [SDUIFrontend](SDUIFrontend/README.md) | SDUI 0.2 parsing, validation and normalized instances | `parser`; combined `sdui` CLI | Separate system boundary; packaging not extracted |
| [SDUICodegen](SDUICodegen/README.md) | SDUI document/instance Go constructors | `codegen`; combined SDL `sdl-gen` command | Dedicated compiler wrapper not implemented |
| [SDUIPresentation](SDUIPresentation/README.md) | Shared relative geometry, bounded Markdown, SVG and structural dumps | `layout`, `markdown`, `svg`, `presentation`; combined `sdui` CLI | Shared libraries implemented; dedicated packaging proposed |
| [SDUIRuntime](SDUIRuntime/README.md) | Typed sessions, events, atomic properties and state reconciliation | `runtime`, `reload`; embedded in native/combined apps | Library implementation; no standalone runtime service |
| [SDUINativeHost](SDUINativeHost/README.md) | Fyne controls, UI-thread ownership and source reload | `host/fynehost`; `cmd/sdui-fyne` | Existing desktop-tagged executable; reusable adapter |

## Integration boundaries

Frontend produces source-positioned documents and normalized trees. Presentation
consumes trees or runtime snapshots and yields shared layout geometry/static
exports. Runtime owns typed UI state; native hosts own native widgets and event
loops. Codegen produces Go constructors using the same model types. It does not
replace runtime behavior or invoke user callbacks during generation.

The native Fyne host is distinct because it already has an executable, event loop,
display dependencies and resource lifetime. It is not a new parser, layout engine
or runtime implementation. A FOX renderer and KanBan TUI do not become implemented
SDUI backends merely because earlier discussions considered them. XFMD remains an
external document viewer, not an obligatory SDUI runtime dependency.

SDL owns its action engine and the explicit UI/action bridge in `SDL/go/bridge`.
The bridge consumes loaded modules, an SDUI session and validated binding plans;
it never discovers modules by executing `ref` paths. Generated applications can
combine SDL/SDUI libraries with handwritten Go domain functions. These bindings
are Go APIs, not a delivered language-neutral ABI.

## Coverage and authority

| Current source responsibility | Catalog owner | Source evidence |
| --- | --- | --- |
| Parser, AST, profiles, reuse expansion | SDUIFrontend | [parser](../../../SDUI/go/parser) |
| Combined inspect/export CLI | Frontend and presentation share current process | [sdui main](../../../SDUI/go/cmd/sdui/main.go) |
| Model constructor generation | SDUICodegen | [Generate](../../../SDUI/go/codegen/generate.go) |
| Combined generator/manifest orchestration | SDL toolchain; SDUICodegen called as library | [sdl-gen](../../../SDL/go/cmd/sdl-gen/main.go) |
| Relative tracks, aspect, clipping, hit geometry | SDUIPresentation | [layout](../../../SDUI/go/layout) |
| Bounded Markdown and configured Mermaid resources | SDUIPresentation | [Markdown provider](../../../SDUI/go/markdown/provider.go) |
| SVG and console/Markdown export | SDUIPresentation | [SVG](../../../SDUI/go/svg), [structural dumps](../../../SDUI/go/presentation) |
| Session identity, draft/value, events/properties | SDUIRuntime | [runtime](../../../SDUI/go/runtime) |
| Candidate preparation and source observation | SDUIRuntime candidate services, host publication | [reload](../../../SDUI/go/reload/watch.go), [sourcewatch](../../../SDUI/go/sourcewatch/watch.go) |
| Desktop event loop and controls | SDUINativeHost | [native command](../../../SDUI/go/cmd/sdui-fyne/main.go), [adapter](../../../SDUI/go/host/fynehost) |
| Explicit SDL action binding | External SDL integration, consumed by hosts | [binding](../../../SDL/go/bridge/bind.go) |

The [existing detailed model](../../../SDUI/design/architecture.design),
[language](../../../SDUI/docs/language.md), [layout contract](../../../SDUI/docs/go-layout-contract.md)
and [runtime contract](../../../SDUI/go/runtime/README.md) remain authoritative
for existing behavior. New System.design entries describe bounded capability and
interaction projections without copying the legacy phase/activity graph. There
are no new UI screens in this catalog: these systems mostly provide tool/library
services, and an invented screen would not establish a useful verified behavior.

Each entry is design-core 0.5 and validates independently. Ecosystem and System
are directory metadata, not parser keywords. Separate systems do not imply IPC,
different implementation languages or mandatory processes. Each README records
reproduction commands, current/proposed status and remaining packaging decisions.
Generated AST/viewpoints belong in temporary exports or explicit evidence, not
as manually maintained source. Project navigation discovery is not automatic.
