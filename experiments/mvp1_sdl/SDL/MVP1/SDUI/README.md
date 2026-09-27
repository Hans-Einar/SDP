# MVP1 presentation sources

Current profile: **SDUI 0.2**. These are proposed screen designs using the existing
Go parser, normalization, layout and SVG exporter. They do not implement Ponsse
logic. All numbers are illustrative; all buttons are unbound. No symbolic callback
is invented for a structural SDL Functionality that has no executable signature.
The APT activation control is deliberately disabled in the sample state.

[Generated review gallery](../../../preview/index.md) shows the screens at
1600 x 900 logical units. The source uses relative sizes and a width-driven 16:9
frame; fonts stay absolute. Headers contain controls, nested groups preserve local
rows, and body frames explicitly fill the remaining height. The stem plot and
editable matrix are provider/form placeholders, not new widget types.

| Source / entry | Application owner | Relevant SDL design |
| --- | --- | --- |
| [BuckingUI/MainPage.sdui](BuckingUI/MainPage.sdui), page | BuckingUI | [Operator projection](../Containers/BuckingUI/Domain.design), [representations](../Containers/BuckingUI/Views.design), [InspectMachine](../Scenarios/InspectMachine.design), [PlanCurrentStem](../Scenarios/PlanCurrentStem.design) |
| [BuckingUI/AptEditor.sdui](BuckingUI/AptEditor.sdui), page | BuckingUI | [AptDomain](../Containers/BuckingService/AptDomain.design), [EditAptCell](../Scenarios/EditAptCell.design), [StoreApt](../Scenarios/StoreApt.design), [ActivateApt](../Scenarios/ActivateApt.design) |
| [SimulatorUI/MainPage.sdui](SimulatorUI/MainPage.sdui), page | SimulatorUI | [Simulator projection](../Containers/SimulatorUI/Domain.design), [SimulatorWireAction](../Scenarios/SimulatorWireAction.design), [command recovery](../Scenarios/RecoverCommand.design) |

## Intended bindings, not implemented calls

Stable named controls are candidates for later typed host registration. SDL refs
and setHandle require actual registered objects/actions; this design does not
pretend that a `.design` source is already an executable action module.

| Screen/control group | Intended boundary | Required behavior before binding |
| --- | --- | --- |
| Operator advice and selection | OperatorUiDomain / RouteOperatorIntent; BuckingWeb to domain services | Use domain availability and coherent stem/input revisions; a requested plan is not an accepted plan |
| Measurement basis/cursor/profile controls | Presentation/context selection | Do not modify underlying machine measurements; unavailable/calibration state remains explicit |
| APT selected/price and selected/submit | RouteOperatorIntent → AptDomain through the selected channel contract | Carry expected draft revision; distinguish editable draft from accepted value and report rejected/conflicting edits |
| APT importFile/load/store/exportFile/activate | Domain-owned APT lifecycle | Import, store, export and activate are distinct actions; no renderer-owned codec, file write or activation authority |
| Simulator session fields/load/reset and scenario controls | RouteSimulatorIntent → HandleSimulatorControl | Validate seeded scenario/configuration; separate hidden truth from emitted wire observations |
| Simulator command lookup/retry | UiCommandCoordinator / RecoverUiCommand | Preserve command identity; show unknown, expired, rejected and completed separately |
| Navigation buttons | Application presentation selection | Change the selected screen without implying a domain operation |

Control names are local to their enclosing named frame/group. For example,
`page/body/matrix/selected/price` identifies the APT input. Host registration must
use the actual normalized instance paths, not infer paths from this table.

## Reuse and deployment

[Shared](../Shared/README.md) records the reusable mechanism candidates and
MVP1-specific shared vocabulary. SDUI 0.2 has no cross-file import facility;
these small screens are self-contained, with named local groups where useful.
Do not create a second template/parser layer to share a few labels.

The application owns projections and typed intents. A renderer owns native
controls, drawing and local interaction. SDUI is an application presentation
asset, not a renderer implementation. A future Fyne host can consume these
sources; a React application may implement the presentation contract or use an
explicit adapter. Neither integration is delivered by parsing/exporting this design.

## Reproduce

From the SDP repository root, using Go 1.26+:

```sh
go -C SDUI/go build -o /tmp/mvp1-sdui ./cmd/sdui
python3 experiments/mvp1_sdl/export_ui.py --sdui /tmp/mvp1-sdui
```

This replaces only the named generated review files in preview/. Use `--output`
for a temporary review directory. manifest.json records source and executable
hashes and every output hash. It is a review snapshot, not the future navigation
cache. Individual exports use the same executable and source with `--entry page`
and `--format ast`, `dump` or `svg`.

The existing Fyne host can also load an unbound screen for a visual trial:

```sh
go -C SDUI/go run -tags desktop ./cmd/sdui-fyne -entry page ../../experiments/mvp1_sdl/SDL/MVP1/SDUI/SimulatorUI/MainPage.sdui
```

Native interaction was not verified in this pilot. No callback or data binding
becomes functional merely by opening this window.
