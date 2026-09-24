# SDL in Go

Run `sdl-design` to browse the documentation. It opens the main page and navigator
in XFMD using prebuilt tools; details are generated when selected.
[Launcher and setup](../scripts/README.md).

G4-M1 provides parsing, source-positioned AST, symbol/type checking, data/wire and
Channel/scenario validation, and canonical formatting for **design-core 0.5**.
Structural statements are not executed. Checkpoint candidates and the complete
MVP1 corpus are outside this profile. G4-M2 adds the explicit
[action-core 0.1](../docs/profiles/SDL-Executable-Action-Profile.md), with a parser,
typed records, registered Go functions and runtime.

Module: `github.com/Hans-Einar/SDP/SystemDesignLanguage/go`. Shared baseline:
Go 1.26; verified with Go 1.27.1. The structural core uses the standard library.
From this directory:

```sh
go test -race ./...
go run ./cmd/sdl check ../../SDUI/design/architecture.design
go run ./cmd/sdl ast ../../SDUI/design/architecture.design
go run ./cmd/sdl format ../../SDUI/design/architecture.design
```

`parser.Data` preserves the Python profile's named JSON AST, including start/end
positions. Comparison evidence covers 151 port cases and the full SDUI model.
[Evidence](evidence/G4.md), [shared plan](../../SDUI/docs/implementation-plan.md),
[language definition](../docs/studies/Design-Language-Definition.md).
G5-M4 replaced the Python frontend and viewpoint generator. Frozen fixtures remain
historical oracles; there is no fallback.

Check actions with `go run ./cmd/sdl action-check examples/echo.sdl`.
`runtime.New` requires explicit signature registration; source is never executed
as Go. G4-M3 supplies a typed SDUI port in `bridge/`. Echo and explicitly simulated
EditAptCell pass through both runtimes. `go run ./cmd/sdl-simulate` emits correlated
events; `go run -tags desktop ./cmd/sdl-demo` opens the native Fyne prototype.
`examples/simulation` contains separate, handwritten Go domain logic.

G4-M4 watches both sources. Invalid source retains the last valid model and an
error status. SDL model reload preserves Go-owned domain state and rejects old
events. Changed Go code requires a build and restart:

```sh
go run ./cmd/sdl-dev -root . -package ./cmd/sdl-demo -tags desktop
```

The development host retains the running process on build failure. Successful
builds start a new process; state across process restarts requires storage.

G6-M1 ports all 11 structural viewpoints without changing their source basis:

```sh
go run ./cmd/sdl viewpoints ../../SDUI/design/architecture.design --output /tmp/sdl-navigation --project sdui-design
go run ./cmd/sdl viewpoints ../../SDUI/design/architecture.design --output /tmp/sdl-static --format static --monolithic
```

Optional `--renderer /absolute/path/to/mmdr` produces SVG for static export.
`--viewpoint VP02,VP08` limits the export. Navigation is the default and renders no
details. Its action links require the registered XFMD reader adapter.
Select a view on demand:

```sh
go run ./cmd/sdl view ../../SDUI/design/architecture.design --uri 'sdl-view://sdui-design/VP02?focus=SduiFrontend&relations=contains&direction=out&depth=1&level=A2' --output /tmp/sdl-selected
```

`--renderer` is optional here too. Publication completes before `entry.md` is
offered. The CLI resolves the source explicitly; host adapters register project IDs.
Optional document service (Linux):

```sh
go run ./cmd/sdl-viewsd -source ../../SDUI/design/architecture.design -project sdui-design -xfmd /absolute/xfmd -renderer /absolute/mmdr
go run ./cmd/sdl-view-request -socket /private/sdl/views.sock -uri 'sdl-view://sdui-design/VP02?diagram=VP02-roots' -window design-one -client terminal -sequence 1 -open
```

The service prints its socket on startup. Use a new sequence per client/window/panel.
Without `-open`, it returns a readable package and lease. `-release TOKEN` releases
it; `-sweep` removes only released packages. Leases survive daemon crashes; crashed
readers require explicit release. XFMD phase 050 releases on replacement/close.

G5 generates and builds models alongside handwritten Go domain code:

```sh
go run ./cmd/sdl-gen -actions examples/edit-apt-cell.sdl -ui examples/edit-apt-cell.sdui -output examples/generatedmodel -package generatedmodel
go run ./cmd/sdl-compiled
go run -tags desktop ./cmd/sdl-compiled-fyne
go run ./cmd/sdl-document -ui examples/edit-apt-cell.sdui -state examples/accepted-state.json -design ../../SDUI/design/architecture.design -output /tmp/sdl-ui-document
```

Generated constructors do not open source files. `-values` on sdl-compiled or
sdl-simulate selects simulated events. sdl-document invokes no callbacks; `-state`
specifies accepted values, labels and enabled/visible properties. The manifest
owns generated files only. Full layout exports as SVG and Markdown with provenance
and an SDL navigator. The class profile requires explicit declarations; it infers
nothing from contains/owns:

```sh
go run ./cmd/sdl class-check examples/runtime-classes.sdl
go run ./cmd/sdl class-view examples/runtime-classes.sdl --output /tmp/sdl-classes --renderer /absolute/mmdr
```

[Go generation](../../SDUI/docs/go-generation.md), [G5 evidence](evidence/G5.md),
[class-core](../docs/profiles/SDL-Class-Profile.md),
[consolidated status](../../SDP/History/checkpoint-1/11-Go-Implementation-and-Navigation.md).
