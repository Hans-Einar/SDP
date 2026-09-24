# SDL describes the SDL/SDUI parsers and runtimes

**Daily source-based browsing:** run `sdl-design`. The [launcher](../../SDL/scripts/README.md)
opens the main page/navigation; details are generated on selection. Stored full exports
below are verification/export artifacts, not the live entry point.

Start with the [generated G1–G6 implementation plan](viewpoints/implementation.md): milestone
status, owners, dependencies and examples from SDL. [Viewpoints](viewpoints/viewpoints.md)
and [index](viewpoints/index.md) derive from this source; [tool documentation](../../SDL/tools/README.md)
explains regeneration/selection. Use-case/Feature allocation, data maps, packets and
contract-checked Channel sequences are tool-generated, without hand-drawn diagram facts.

Updated 2026-09-22. [architecture.design](architecture.design) is the shared target structure
for **both languages**, replacing the former shorter SDUI/host model. No duplicate exists
under SDL. It uses implemented design-core 0.5 and passes the SDL parser. This is the SDL
structural profile, not retired SDUI 0.1. [Implementation status](../../SDP/History/checkpoint-1/11-Go-Implementation-and-Navigation.md)
bounds executable support.

## Open the model

- [Source](architecture.design): authoritative ownership, structure and boundary use.
- [Responsibility overview](viewpoints/viewpoints/VP02/index.md): Units, responsibilities and ports.
- [AST/symbol table](architecture.ast.json): typed nodes and source positions.
- [Validation](architecture.validation.json): parser result; manifests record source/tool hashes.

From repository root:

```sh
go -C SDL/go run ./cmd/sdl check ../../SDUI/design/architecture.design
go -C SDL/go run ./cmd/sdl ast ../../SDUI/design/architecture.design
go -C SDL/go run ./cmd/sdl viewpoints ../../SDUI/design/architecture.design --format static --monolithic --output ../../SDUI/design/viewpoints --renderer /absolute/mmdr
```

CLI exports positioned AST and indexed viewpoint directories. The
[navigator](navigation/navigator.md) uses a registered XFMD tool to generate selected details.
The [UI/state example](runtime-preview/entry.md) uses the same Go toolchain with explicit state.
The model has 463 declarations, 1424 facts, 116 milestone-linked Functionality objects,
42 Activities (including six phases/24 milestones), 22 Channels, 48 Messages and 13 scenarios.
Status follows actual verification; model facts alone do not prove execution.

## Use cases and architectural contributions

DesignAuthor/DesignReviewer are external roles. Features supports UseCases; Functionality
contributes-to Features/use cases. TraceViewpointFacts contributes directly to InspectModels.
StructuralModelInspection supports inspection and editing; shared responsibilities retain
identity/ownership in every view. Allocation places responsibilities in CommandLineHost or
FyneHost in explicit modes without moving library ownership. Allocations are design claims,
not runtime proof. VP07 reports 18 incomplete mode slices; allocation is not exhaustive.
A shared Functionality allocation does not make every related Feature executable in that mode.

SdlViewpointGenerator under SdlLibrary owns ProjectSdlViewpoints, ExportViewpointMarkdown and
TraceViewpointFacts. Go viewpoint implements this bounded scope; its Python port is complete.

## Parser design

SdlLibrary/SduiLibrary are logical libraries, each with lexer, parser, validator and normalization.
They implement different source languages. Frontends coordinate; child Units own Functionality.

| Stage | SDL | SDUI |
| --- | --- | --- |
| Tokens/AST | SdlLexer, SdlParser | SduiLexer, SduiParser |
| Names/profile | SdlValidator: symbols, typed structural rules/profile | SduiValidator: names, arguments, relative formatting, symbolic bindings |
| Normalization | SdlNormalizer: validated facts/source maps | SduiNormalizer: expansion, header/body/footer/source maps |
| Execution gate | SdlExecutionGate checks the separately defined executable profile | SduiRuntime can display an unbound validated UI |

SourceLoader reads bounded revisioned snapshots; parsers receive loaded content. SDUI ref
opens no file while parsing. DiagnosticReporter owns source/binding diagnostics. Limits apply
throughout the chain. stateless/deterministic frontend requirements take source/profile as
inputs: temporary parser objects are allowed, retained session state is not. They do not imply
deterministic domain/network/font behavior.

## Runtime and domain binding

SdlRuntime owns SDL instances; SdlStateStore models domain state/snapshots. SdlFunctionRegistry
registers handwritten Go functions and checks signatures. SdlDispatcher validates action input,
invokes functions, correlates results and models cancellation. GoDomainImplementation is an
external logical Unit owning product algorithms, not a claimed process/machine implementation.
These are model responsibilities; actual runtime limits are in action-core and Go contracts.

SduiRuntime owns UI instances; SduiInstanceStore owns identities/generations. SduiPropertyStore
validates/publishes batches and tracks drafts; SduiDispatcher routes valid events and rejects old
ones. UiStateReconciler preserves/resets compatible state. Fyne owns concrete caret/focus while
runtime retains logical identity/value/draft. Toolkit pointers are not SDL handles.

SdlUiBindingAdapter is outside both cores. It resolves callbacks, binds typed handles, routes
calls and publishes domain results. Removing it does not invalidate static UI. DomainBindingPort
faces UI; SdlExecutionPort faces SDL. G3/G4 implement bounded event profiles; VP08 checks modeled
message order while runtime tests check effects:

1. Fyne sends activation/commit with identity, generation and revision.
2. SDUI checks widget, enabled state, binding and value.
3. Adapter resolves the action; SDL checks input and Go signature.
4. Registered Go performs domain work and returns a correlated result.
5. SDL publishes the accepted change; adapter constructs the UI update.
6. SDUI validates the entire batch and publishes snapshot/layout to the host.

Programmatic value updates do not retrigger input callbacks. Pending/late results require
generation checks. [Runtime contract](../docs/runtime-contract.md), action-core and Go specify
signatures, errors and delivery rules.

## Hot reload and Go generation

SourceWatcher collects changes. ReloadCoordinator builds/validates candidates before publication.
UiStateReconciler handles compatibility; modeled DomainStateMigrator uses explicit migration/reset
policy. Bindings switch at the agreed generation boundary. Errors retain the previous model.
An isolated UI change reparses SDUI only; changed SDL bindings can require UI revalidation without
changed UI source. Entire changed language models rebuild, not incremental machine code. Reload
must never replay completed actions; publication coordinates calls, cancellation, focus/drafts
and resources. General automatic domain migration is not implemented.

GoCodeGenerator models constructors/binding registration; actual generated constructors share the
file-based runtime, while explicit bridge plans remain handwritten Go. GoBuildRunner builds and
restarts changed Go code, without replacing code inside the process or overwriting domain code.

## Layout, content and hosts

SduiLayout measures, resolves ancestor-relative dimensions, allocates/clips geometry and prepares
shared frames. SduiPresentation provides SVG/console/Markdown export without domain execution.
MarkdownProvider, DiagramProvider and ResourceStore are separate services. Mermaid is an external
diagram engine, not SDUI grammar.

FyneBackend creates/updates controls, handles input/focus and releases native resources. FyneHost
owns interactive composition/UI-thread publication; CommandLineHost exports without a window.
These are modeled Containers; libraries are Units. contains is logical decomposition, not proof
that every library resides in each process. Fyne/SVG share measured geometry; native imports stay
outside parser/runtime. Different text/SVG backends still need visual verification.

## Ports and modes

consumes expresses use, not a runtime call or a requirement in every mode. provides applies to
Capability, **not Interface**, in the current parser. Concrete provider/signature relationships
below are explained contracts, not falsely encoded with unsupported SDL object types.

| Ports | Provider and contract boundary |
| --- | --- |
| SourceInputPort, FileChangePort | Host file/watcher adapter; source base, limits/errors |
| SourceSnapshotPort, DiagnosticPort | SourceLoader/DiagnosticReporter; identity/revision, positioned diagnostics |
| SdlFrontendPort, SduiFrontendPort | Respective frontend; source → validated model/diagnostics |
| SdlModelPort, SduiModelPort | Normalization/generated constructors; immutable model/source maps |
| ExecutionProfilePort | Explicit SDL execution profile; supported operations and complete binding |
| DomainFunctionPort, DomainStatePort | Registered Go/SdlStateStore; signature, owner, result/revision |
| SdlExecutionPort, DomainBindingPort | SdlRuntime/binding adapter; typed actions, correlation/lifetime |
| UiSessionPort, UiStatePort, UiSnapshotPort | Runtime/stores; logical state, events/snapshots |
| SdlReloadPort, UiReloadPort, BindingReloadPort, ReloadPort | Runtimes/adapter/coordinator; preparation, compatibility, publication/cancellation |
| MeasurementPort, ContentProviderPort | Host measurement/MarkdownProvider; width→height, font unit, bounded work |
| DiagramPort, DiagramEnginePort, ResourcePort | Providers/engine/store; inert content, identity/release |
| PreparedFramePort, WidgetBackendPort | SduiLayout/FyneBackend; geometry/clipping, lifetime/events |
| GeneratedArtifactPort, BuildToolPort, ExportSinkPort | Code generator/host build/file output; ownership and explicit execution |

Modes are independent named contexts, without implicit inheritance or executable activation:
SourceInspection, unbound UiPreview, BoundExecution requiring SDL/domain binding, UI LiveEditing,
coordinated BoundLiveEditing, StaticExport, NativeBuild and RichDocument. A host may activate
multiple contexts under explicit policy. Activity refines is refinement, not sequencing,
a state machine or scheduler.

## Language boundary and traceability

The model covers [G1–G6](../docs/implementation-plan.md). design-core 0.5 checks typed records,
explicit steps and correlation, not Go signatures/calls, atomicity, state machines, budgets,
thread rules or code/model conformance. V2–V4 evolved grammar/tests together; G4/G6 define
action-core/class-core separately.

| Model area | Implementation |
| --- | --- |
| SduiFrontend and children | G1; SDUI/go/parser |
| SduiRuntime and children | G3; SDUI/go/runtime |
| SdlFrontend and children | G4-M1; SDL/go/parser |
| SdlRuntime, SdlUiBindingAdapter | G4-M2–M4; SDL runtime/host binding |
| SduiLayout, content, Fyne, presentation | G2; SDUI/go/layout, markdown, svg, host/fynehost |
| ReloadCoordinator, SourceWatcher | G3/G4; shared development host, not duplicated parsers |
| GoCodeGenerator, GoBuildRunner | G5; both codegen packages and SDL/go/devhost |

These paths are manual traceability, not automatic package/file conformance checks. The
[target architecture](../docs/target-architecture.md) owns principles; this model owns detailed
responsibility allocation.

## V2 — data and selected packet example

DesignSourceDocuments resides at DesignSourceArchive, a planned persistent source owned by
SourceLoader. UiSessionState is transient at SduiInstanceStore. Contracts require source/revision
and UI generation with optional draft; absence is not an empty string. ProjectUiGeneration
projects to UiGenerationNotices. UiGenerationWire is an explicit **prototype** format: 16-bit
version, 64-bit generation, big-endian/MSB-first. It is not an adopted Go transport, SDL ABI,
P1000 or StanForD format; native Go calls need no such serialization.

## V3 — messages and scenarios

V3 introduced seven Channels/four scenarios for accepted/rejected actions/reload. BoundActionAccepted
runs Fyne → UI dispatch → adapter → SDL dispatch → registered Go responsibility with correlated
returns. BoundActionRejected stops at UI validation without domain invocation. These are modeled
collaboration, not observed execution. Reload separates published models/rejected candidates; only
the accepted path contains UiGenerationNotices. Optional draft/generation/diagnostics express
meaningful absence. Compatible UI state survives; action-core preserves Go state within a process,
not arbitrary automatic migration. [MessageSet](viewpoints/message-sets.json) is tool-derived.

## V4 — integrated review basis

addresses links milestones to Functionality; delivers links phases to Features; depends-on records
prerequisites; illustrates links scenarios to activities. The tool generates implementation.md from
facts without a hardcoded G1–G6 plan. The ten V4 scenarios cover frontend stages, shared Fyne/SVG
layout, local Go actions, accepted/rejected binding, both reloads and native builds. Bytes fields
for tokens/AST/model/scene are intentionally opaque boundaries. G-phase contracts define Go structs,
lifetime, atomic UI reload and font/dimension units. Message adjacency implies no transformation
or state change.

## G6 — navigable documents

The [navigation design](../../SDL/docs/integration/SDL-Navigable-Viewpoints-Design.md) is modeled as
G6NavigableDocumentation with six verified milestones. DocumentBroker, ViewArtifactStore and
ViewerLaunchAdapter separate generation/publication from XfmdDocumentHost's panes/display.
SelectedViewOpened, InvalidViewSelectionRejected and ViewProjectionFailed generate sequences.
Typed records, URI checks, Linux IPC, publication, leases and explicit windows/panes were tested
in G6. [G6-D2](../../SDL/docs/integration/SDL-Viewpoint-Levels-and-Notation.md) defines export forms,
A0–A5 and notation. Navigation builds overviews without details; M5 supplies symbols/arrows;
M6 class-core adds explicit classes/source-linked diagrams without reinterpreting contains/owns.
