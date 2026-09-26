# SDL navigation and SDPTool design preview

Baseline inspected 2026-09-25, before Sprint-0001. Current code, proposed SDPTool responsibilities and future
XFMD integration are distinguished below. Python is not in the active SDL
parsing/projection/rendering path.

## Sprint-0001 update

The Go SDPTool facade now implements saved-file preview, explicit project
recognition, configured plan viewing, typed SDL/KanBan/SDUI inventory and selected
current-source generation. [Contract](Contract.md) and [consumer examples](Consumer-Examples.md)
own the implemented wire/CLI behavior. Native XFMD sidebar/editor integration,
unsaved buffers and advanced process commands remain separately tracked. The
baseline explanation and original proposal below retain their dated context;
read the current contract for executable signatures and supported capabilities.

## What runs today

| Component | Current role |
| --- | --- |
| [sdl-design](../../SDL/scripts/sdl-design) | Bash launcher using prebuilt sdl, mmdr and XFMD; defaults to the SDUI architecture model and fixes project registration to sdui |
| [sdl CLI](../../SDL/go/cmd/sdl/views.go) | Go commands: viewpoints builds navigator/static exports; view generates a selected bundle |
| [parser and viewpoint](../../SDL/go/viewpoint/model.go) | Go validation, model facts, typed diagrams and Mermaid projection |
| [documents](../../SDL/go/documents/bundle.go) | Go Markdown/resource bundles, source revision, manifests and publication |
| [render adapter](../../SDL/go/documents/geometry.go) | Calls Rust mmdr; uses its layout dump for flowcharts and adds SDL symbols/arrows in Go; other diagram kinds use mmdr SVG |
| XFMD DocumentViews | Registered URI handling, child processes, main/navigation panes and resource lifetime; no SDL parser |
| sdptool | Go preview, discovery, viewer and navigation facade; see [current contract](Contract.md) |

The eleven VP01–VP11 definitions are a fixed supported catalog. Their diagram
instances, content, relations, gaps and selections come from the current validated
SDL model. A viewpoint definition is not a manually maintained diagram.

The launcher creates a private `sdl-design.XXXXXXXX` session below
`XDG_RUNTIME_DIR` (falling back to TMPDIR or /tmp). It runs:

```sh
sdl viewpoints MODEL --format navigator --project sdui --output SESSION/navigation
```

It opens generated index.md in XFMD's main pane and navigator.md in the separate
navigation pane, registering the source, project, SDL tool, renderer and window.
The navigator export includes overview/type-inventory pages, not rendered detail
diagrams. The present projector constructs diagram descriptions in memory during
model processing; this is not an incremental parser or a lazy model database.

An action link is **sdl-view://**, for example:

```text
sdl-view://sdui/VP02?diagram=VP02-roots&target=main&consumer=xfmd
```

It identifies a registered project and selection, not a script or filesystem URL.
Additional supported selection fields include focus, relations, direction, depth,
level and mode. XFMD validates the registration, then invokes:

```sh
sdl view MODEL --uri URI --output REQUEST_DIRECTORY --renderer MMDR
```

The command rereads and validates the saved source on each request. It publishes
entry.md, Mermaid sources, optional SVG assets, selection.json, manifest.json and
delivery.txt. XFMD opens the completed bundle in the addressed pane; the source
file is not replaced. The Markdown navigator is a snapshot and must be regenerated
when its inventory changes; selected detail generation already rereads the source.

Initial files belong to the launcher's session and are removed on exit. Direct
request bundles belong to XFMD's separate private views directory and are removed
on window teardown. These are ordinary temporary files, not a required daemon or
an in-process RAM document API. An XDG runtime directory may be memory-backed by
the operating system; the applications still use filesystem paths.

Optional Go commands sdl-viewsd and sdl-view-request provide a broker and addressed
delivery with leases. The standard launcher does not require them. Do not diagnose
an unregistered sdl-view URI as evidence that a daemon must be started.

## Existing design-file rendering

Both commands above already consume a supported structural `.design` file. Static
export can produce a complete report or selected viewpoints. Omitting --renderer
produces Markdown Mermaid fences; including it produces Markdown with linked SVG
assets and retains .mmd source. Class-profile projection has separate class-view
support; filename alone does not identify all supported language profiles.

There is an important presentation distinction: flowchart Mermaid text does not
encode every SDL symbol. Current SVG output uses mmdr geometry plus
[Go SDL symbol drawing](../../SDL/go/documents/symbols.go), including actor,
use-case, feature and container shapes. Giving only Mermaid to another renderer
may not preserve that appearance. Reuse the complete document pipeline for the
first preview; Mermaid remains an available export/intermediate representation.

Opening a readable `.design` file in current XFMD takes the plain-text route.
There is no automatic SDL preview adapter. Existing generated-document navigation
opens a generated document; that is not yet a source-editor/derived-preview pair.

## SDPTool and the native tree

SDPTool should resolve project context and host-registered tools and delegate to
the existing Go SDL services. Reconcile existing project/install manifests before
choosing discovery metadata. Selected valid SDP area or its SDP child are the two
recognition roots; folder existence alone is insufficient. Report unsupported or
invalid configuration separately from absence. No implicit parent selection.

Expose a versioned machine-readable tree inventory with stable node IDs, kinds,
labels, children/references, source revision and selectable targets. SDL supplies
catalog/model facts; SDPTool coordinates project capabilities and KanBan/SDUI
services. XFMD consumes this inventory using native widgets, not Markdown scraping.
All supported groupable concepts and useful relation selections qualify, not only
UseCases. Handle graph references/cycles explicitly and refresh after source edits.

Native SDP / KanBan / SDL / SDUI tabs belong to XFMD's KB-XFMD-014. Selecting a
node can retain the existing typed selection and generated-bundle mechanism;
no new URI scheme, shell-link execution or mandatory broker is necessary.

## Recommended early delivery: direct preview

Expose a bounded design-preview operation through sdptool, delegating to the SDL
Go parser/projector/documents packages. The command name and wire schema are not
adopted here. A standalone supported file should not require full SDP discovery
or native sidebar work; project context enriches it when available.

Default to a compact overview and relevant content-derived diagrams with explicit
viewpoint selection for detail. Define defaults per supported profile and bound
large outputs. Rendering a file still needs projection rules; it does not imply
one giant graph or all eleven viewpoints on every keystroke. Report unresolved
references/profile errors rather than inventing missing project context.

XFMD keeps the `.design` source buffer, path, dirty state, undo and save behavior.
Its preview adapter displays the derived bundle independently. Saved-file preview
is a first milestone; live unsaved preview needs a buffer snapshot/revision input
contract because the current sdl CLI reads files. Debounce/cancel generation,
discard stale results, retain the last valid preview on parse errors, and show
source-positioned diagnostics. Neither path executes SDL domain logic.

XFMD owns this new behavior in KB-XFMD-015. SDP owns the producer facade in
KB-SDP-017. Prefer this small end-to-end delivery before the complete native tree.
Keep semantic projection in Go and reuse Rust layout; a new sdl-rs-renderer is not
needed for this milestone. Extract or replace a rendering backend only for a
demonstrated capability/performance need, with the same typed model boundary.

## Observed command trial

On 2026-09-25, the installed Go sdl binary processed
SDUI/design/architecture.design: navigator mode wrote 15 files, of which 14 were
Markdown pages, with zero diagram files. A subsequent VP02-roots request using
the configured mmdr produced entry.md, selection.json, manifest.json, delivery.txt,
one .mmd and one .svg. Both outputs carried source revision
`e946313a6ae4a80603a8f13d3467373c8fa1990eb9c42ed8f82b759ee12845e1`.
This verifies existing CLI generation, not native design preview or visual fidelity.
The installed binary reports Go 1.27.1; build metadata references cbb0dc6 with a
modified source tree, so this trial does not certify an exact clean-commit build.
