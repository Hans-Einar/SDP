# SDL ecosystem

This ecosystem groups software systems around SDL language processing, execution
and documentation. It is a source organization and product-boundary proposal under
[PLAN-SDP-0008](../../03--Architecture/Ecosystems/Plan.md), not an SDL language
namespace, monolithic running system or completed package extraction.

## Systems

| System | Boundary | Current delivery |
| --- | --- | --- |
| [Frontend](Frontend/README.md) | Source parsing, checking, AST and formatting | `parser` library; commands in `sdl` |
| [CodeGeneration](CodeGeneration/README.md) | Action constructors and combined SDL/SDUI Go bundles | `codegen`; `sdl-gen` |
| [Viewpoints](Viewpoints/README.md) | Validated facts to selectable/static documents and diagrams | `viewpoint`, `documents`; commands in `sdl` |
| [DocumentService](DocumentService/README.md) | On-demand local delivery and reader leases | `broker`, `reader`; `sdl-viewsd`, `sdl-view-request` |
| [DevelopmentHost](DevelopmentHost/README.md) | Build and restart explicitly selected Go application | `devhost`; `sdl-dev` |
| [ActionRuntime](ActionRuntime/README.md) | Registered Go action execution and SDUI binding | `runtime`, `bridge`, `reload`; embedded library |
| [DocumentSnapshot](DocumentSnapshot/README.md) | Static documentation combining UI state, SVG and optional design navigator | `snapshot`; `sdl-document` |

Each directory contains a standalone design-core 0.5 `System.design`. System
identity is directory metadata because the current parser does not implement
`system`, imports or multi-file symbol linking. These are deliberately bounded
capability catalogs, not a copy of every detailed implementation/design fact.
The old [combined design](../../../SDUI/design/architecture.design) remains its
existing detailed authority. New ecosystem organization does not migrate it.

Frontend and Viewpoints are separate proposed system responsibilities but **share
one current binary**. ActionRuntime is an embedded library; manufacturing a daemon
would not improve the model. Existing binaries may remain compatibility entry
points if later product packaging is selected. Independent packaging remains
planned in all models; an existing command does not prove independently released
products, stable libraries or a versioned cross-language ABI.

## Complete tracked command inventory

Inventory basis: tracked `SDL/go/cmd/*/main.go` and `SDL/go/bridge/cmd/*/main.go`.

| Entry | Owner / disposition |
| --- | --- |
| [`cmd/sdl`](../../../SDL/go/cmd/sdl/) | Frontend: check/ast/format/action-check/class-check; Viewpoints: viewpoints/view/class-view |
| [`cmd/sdl-gen`](../../../SDL/go/cmd/sdl-gen/) | CodeGeneration: combined action/UI constructor bundle |
| [`cmd/sdl-document`](../../../SDL/go/cmd/sdl-document/) | DocumentSnapshot: accepted-state SVG/Markdown export |
| [`cmd/sdl-dev`](../../../SDL/go/cmd/sdl-dev/) | DevelopmentHost: Go rebuild/restart |
| [`cmd/sdl-viewsd`](../../../SDL/go/cmd/sdl-viewsd/) | DocumentService: Linux broker process |
| [`cmd/sdl-view-request`](../../../SDL/go/cmd/sdl-view-request/) | DocumentService: Linux socket client |
| [`cmd/sdl-simulate`](../../../SDL/go/cmd/sdl-simulate/) | ActionRuntime fixture host: interpreted model with simulated domain |
| [`cmd/sdl-demo`](../../../SDL/go/cmd/sdl-demo/) | ActionRuntime/SDUI integration fixture: native Fyne and model reload; requires desktop build tag |
| [`cmd/sdl-compiled`](../../../SDL/go/cmd/sdl-compiled/) | CodeGeneration/ActionRuntime fixture: generated model constructors with simulated domain |
| [`cmd/sdl-compiled-fyne`](../../../SDL/go/cmd/sdl-compiled-fyne/) | Generated-model native fixture; requires desktop build tag |
| [`bridge/cmd/smoke`](../../../SDL/go/bridge/cmd/smoke/) | Native SDL/SDUI binding verification harness, not independent product |

The [`sdl-design` launcher](../../../SDL/scripts/README.md) composes prebuilt tools
for the existing design-navigation workflow. It is an entrypoint convenience,
not another parser, document engine or generic automatic source registration.
Verification Python helpers under `SDL/go/tools` are evidence tooling, not a
production parser fallback. Untracked `SDL/go/sourceinput` is excluded entirely.

## Package reuse and system boundaries

The implementation remains in one
[Go module](../../../SDL/go/go.mod), whose historical import path includes
`SystemDesignLanguage/go`. This delivery does not rename it. Reusable packages do
not become runtime containers:

- `parser` owns source/profile semantics; consuming systems must reuse it.
- `viewpoint` owns projections of validated structural facts.
- `documents` owns document bundles, publication, optional diagram adaptation and
  explicit-class projection. CodeGeneration and DocumentSnapshot reuse publication.
- `broker` and `reader` own local document lifetime and reader delivery.
- `runtime`, `bridge` and `reload` provide action execution/binding/reload libraries.
- `devhost` owns process rebuild/restart, which differs from in-process model reload.
- `snapshot` composes SDUI rendering and optional SDL navigation without callbacks.
- `examples/application`, `examples/simulation` and `examples/generatedmodel` are
  fixture composition, handwritten simulated domain logic and generated fixture
  output; they are not the Ponsse production backend.

SDUI packages supply UI parsing, generation, layout, runtime, Markdown and Fyne.
Mermaid rendering and XFMD remain external integrations. SDPTool can reuse suitable
libraries or registered commands; this catalog does not establish a new stable ABI
or silently adopt independently published binaries.

## Verification and navigation

Every `System.design` is its own parser and projection entrypoint. Read each system
README for commands, limitations and mappings. Validation results are recorded in
[Verification.md](Verification.md). Generated files go to temporary/output areas;
authored source remains here. Catalog navigation is handled by the enclosing SDP
source index; these directories do not register themselves in a running service.

The models deliberately describe implementation observations and selected target
product boundaries separately. Structural validation checks declarations and
relations. It is not runtime acceptance, API conformance, native rendering evidence
or proof that the repository has seven independently released SDL products.
