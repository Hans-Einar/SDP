# SDPTool producer contract 0.2

This is the implemented local facade contract, not a new Toolkit release.
Go language packages own parsing, validation, projection and presentation.

## Output presentation

All result-producing operations default to human-readable output. Add `--json`
for the machine contract; output does not switch implicitly when piped. This is
an intentional CLI compatibility break relative to published SDP 1.0.0. Clients
must opt in, including version probes. Help remains text.

`presentation.Registry` routes serialized results by schema and operation to
compiled-in Go adapters. The CLI retains operation/exit/stream ownership. No
Bash, jq, external formatter, dynamic plugin loader or shell process is required.
Unknown routes retain JSON bytes; `--json` bypasses adapters. A renderer failure
is an error, never a partially printed result followed by fallback JSON.
Tree rendering preserves inventory order, shows state and references, stops at
cycles/depth bounds, and escapes terminal controls in labels. Registered general
results render a deterministic field hierarchy. Installation keeps its summary
and existing error streams. Generated document/image files are unchanged.

Use `sdptool tree --json`, `sdptool --version --json`, or global
`sdptool --json PROJECT tree`. String option values (including a literal
`--json`) remain values. Required positionals still precede command flags;
`preview FILE --json` is valid. `--json=false` selects readable output.
`release-log --version VERSION` still exports Markdown; with `--json` it emits
an envelope containing that Markdown. Generation/check mode emits a receipt.

See the [plan and consumer handoff](../SDP/05--Implementation/SDPTool/Output/Plan.md).

## Saved design preview

```sh
sdptool preview model.design --output /tmp/request-preview
sdptool preview model.design --output /tmp/request-svg --renderer /absolute/mmdr
sdptool preview model.design --output /tmp/request-detail --viewpoint VP01
sdptool preview model.design --output /tmp/request-detail --uri 'sdl-view://project/VP02?diagram=VP02-roots&target=main&consumer=xfmd' --revision SOURCE_SHA256
```

Source precedes flags. Supported structural profiles are SDL design-core/0.5 and source-composed
design-core/0.6 with canonical-source validation; unsupported language/profile input returns the
language diagnostic. Class/action/SDUI inputs are not structural design previews.
No SDP folder is required. Default selection is one nonempty diagram, preferring
architecture then use cases/responsibilities; an empty model uses VP11. Explicit
selections retain SDL's query validation and depth bound (0–8). There is no fixed diagram-count cap. Publication is bounded by 128 files and
32 MiB of generated resources, including retained source snapshots; larger
requests must select a diagram/focus. No automatic full export occurs.

With --json, stdout is one JSON result with schema sdptool/0.2, operation, source, profile,
revision (SHA-256 of exact source bytes), entry, directory and ownership. The
bundle reuses SDL entry.md, diagrams/*.mmd, optional diagrams/*.svg, selection.json,
manifest.json and delivery.txt; sdptool.json records facade provenance. SVG uses
SDL's symbols and existing Rust renderer geometry. Without a renderer, Markdown
contains Mermaid fences. Renderers are prebuilt programs supplied by the host.

Errors return nonzero; with --json they emit one JSON error on stderr (schema, error.code/message
and positioned language diagnostic when available). Codes include arguments,
source, model, selection, tool, render, output, stale, canceled and limit.
Caller cancellation and stale expected/source-during-generation revisions prevent
publication. Requests are synchronous: callers assign their own request IDs and
must discard replies for an obsolete selection even when source hashes match.

Output is caller-owned. Allocate a distinct private request directory per
concurrent consumer request; release it only after the consumer releases all
resources. SDL's guarded publisher preserves unmanaged notes, rejects conflicts
and changed generated files, and replaces complete bundles. Source containment
checks include symlink ancestors. Do not use a project/source directory as output.
A failed generation leaves an existing valid preview intact. This is filesystem
publication, not a daemon or implicit in-memory service. A save occurring after
the final revision check is observed by the next request; consumers compare
revision/request identity rather than assuming a permanently current snapshot.

## Project recognition and source discovery — DS1

Select a project with an immediate real `SDP` directory, or select that directory
itself. No ancestor search or Git identity guessing occurs. A directory is a
discovery area, not proof of installed process conformance. Existing project and
installed manifests supply installation facts; an absent manifest is allowed.
An active installation journal returns `incomplete` before source enumeration.

`discover --json` returns one `sdptool/0.2` snapshot: project identity, installation
facts, derived `inventory`, `sources`, `programs`, `programDiscovery`, `plans` and
`navigation`. Nothing is written. Programs use the explicit contract below;
ordinary source discovery does not make every SDUI file runnable.
`navigation.json` is neither read nor created. Existing project-owned copies are
preserved as inert historical files during upgrade, even when malformed.

Enumerate the selected SDP area, including old examples and actual directories.
Do not infer which valid source is "current". Parse `.design` and `.sdui` through
the owning language packages. Supported standalone models and System roots gain
semantic navigation; fragments remain visible as `context-required`; invalid and
unsupported sources carry diagnostics. SDL 0.6 composes files using source-owned
includes/contains. Different roots are evaluated independently, never concatenated.
The returned `inventory.models` and `inventory.sdui` are observations, not editable
registries. A model ID combines its normalized project-relative path with a hash;
content changes preserve identity, renames change identity. Display names are not IDs.
Use returned IDs/targets rather than copying old registration names.

All paths are relative to the parent of SDP, except typed open targets which carry
absolute paths. Enumeration does not follow symlinks; symlinks remain visible.
Containment checks also apply when resolving sources and generating outputs.
Skip `.git`, `.sdp-operations` and `.sdp-backups`. Limits: 10,000 filesystem entries, depth 64,
256 source files, 2 MiB per source and 64 MiB aggregate source/Session input. The combined
semantic navigation retains its 20,000-node and 32 MiB result bounds.

The viewer stores this snapshot in memory. It watches the SDP filesystem or offers
manual Refresh, then invokes discovery again. It owns debounce, cancellation and
request ordering; an older response must not replace a newer one. Generation is
still on demand and rechecks source revisions. Discovery is not an atomic filesystem
transaction and does not claim a permanently current view while files are edited.

KanBan follows `SDP/KanBan`. Sessions follows `SDP/Sessions`. ImplementationPlan documents under
`05--Implementation` are reported in `plans`. A single plan can be selected
implicitly; otherwise `view ip --plan PROJECT_RELATIVE_PATH` selects it explicitly.
A single SDL model can be selected implicitly; multiple models require `--model`
for operations acting on one model. Unfiltered `tree` shows all discovered models.

## Commands, host policy and consumer protocol — T1-M2

| Invocation | Result and ownership |
| --- | --- |
| sdptool [PROJECT-OR-SDP-AREA] discover | Source discovery and navigation snapshot; read-only |
| sdptool preview FILE --output DIR | Standalone saved-file operation above |
| sdptool [PATH] view ip --plan PATH --model ID | Open the selected authored plan with generated navigation in the configured viewer |
| sdptool [PATH] tree --model ID | Navigation inventory; no detail rendering or writes |
| sdptool [PATH] select --model ID --uri URI --revision HASH --output DIR | Validate project binding/revision, then generate selected current-source detail |
| sdptool [PATH] sdui-preview --model ID --output DIR | Only the implemented, discovered SDUI service selected in T3 |

`implementation-plan` is an alias for `ip`. `generate ip` remains unsupported;
opening a plan must not synthesize, overwrite or imply approval of one. Explicit
path precedes the command; omission means current directory. Flags follow the
operation's fixed positional arguments. Reject extra arguments and unknown flags.

Host executable precedence: explicit --viewer/--sdl-tool/--renderer options,
then SDP_XFMD/SDP_SDL_TOOL/SDP_MMDR, then viewer/sdl names on PATH when needed.
No renderer means Mermaid output, not an automatic build/install. Resolve programs
before launch, preserve argument boundaries and never invoke a shell. Ordinary
source and preview metadata cannot register executables. The owner-selected
`SDP/programs.json` application contract below is the explicit exception: only
`run` executes its command. No preview-time compilation or required daemon.
Standalone library calls accept already chosen options and do not read host policy.

The initial bridge opens the actual plan in the main pane and the selected model's
generated navigator in the navigation pane. Register --navigator, --sdl-tool,
--sdl-source, --project and optional --renderer with existing XFMD conventions;
allocate a unique --window-id. A model-free project can open its plan without
SDL flags. The configured viewer is the long-lived process owning that window;
keep temporary navigation resources until that process exits, then clean up on
normal/error/canceled exit. A viewer that detaches must use a future explicit
lease adapter, not this synchronous bridge. Use XDG_RUNTIME_DIR when available,
otherwise the normal temporary directory; never hardcode a user ID.

JSON responses identify schema `sdptool/0.2` and operation. Recognition reports
status, root, area, inventory, sources and navigation. Capabilities indicate
`discovered` sources; individual source states distinguish validation from errors. Tree replies include source revision,
nodes and roots. Each node has stable id, kind, label, state and optional children,
reference or typed target. A target carries project/model identity, operation and
an SDL-owned URI or source path; it is not a shell command. Nodes are shared by ID
rather than recursively cloning graphs. Bound expansion with explicit references.
All catalog viewpoints appear, with empty/unsupported states where applicable.

On structure changes refresh the tree. Select with the exact source revision
from that tree; reject mismatches before rendering and again before publication.
A client tracks its own monotonically increasing request/selection identity and
ignores late replies, including same-revision replies for a former selection.
Cancellation kills owned child requests; it never signals unrelated consumers.
Diagnostics leave the previous successful bundle displayed. Consumers release
bundles explicitly according to the ownership contract; transport completion alone
is not a request to delete files still in use. Fixture/harness validation belongs
to T4-M1; actual XFMD GUI integration is T4-M2 and XFMD-owned.

## Implemented tree and selection limits — T3

The SDL tree groups requirements (A0/A1), architecture (A2/A3), design (A4) and
implementation (A5/delivery). These are viewpoint groupings, not an assignment of
every model object to a process folder. Every catalog viewpoint appears, even
when empty. Typed collections use diagram/source-fact membership; VP11 exposes
all declarations. Canonical object IDs combine model ID, kind and authored name;
a rename changes identity because this SDL profile has no separate persistent ID.

Relationship leaves reference canonical objects instead of recursively copying
them. Consumers follow references with a visited set and the declared expansion
limit of eight. The producer graph has a 20,000-node cap. SDL selected generation
also enforces its query depth limit. Relationship IDs hash semantic endpoints;
fact source positions remain owned by SDL. Inventory only builds data, not SVG.

`select` requires both expected revision and a URI whose project matches the
selected project, plus an explicit model when multiple sources make selection ambiguous. It delegates to the
same guarded preview path. Refresh and retry after a stale error. A foreign URI
never switches project/source selection implicitly.

## KanBan and SDUI services — T3-M3

`tree` returns Files, KanBan, SDL, SDUI and Sessions roots, with unavailable diagnostics rather than
hiding failures in optional services. `revision` is the selected SDL source hash;
`inventoryRevision` also changes with source inventory, filesystem entries, card/ledger data and SDUI
sources. Consumers refresh on either appropriate revision. All targets retain
source/card-specific hashes. No UI state is persisted by the inventory operation.

The initial board reader supports board schema 0.2 with the local
sdp-project-management/0.1 or /0.2 profile, payload 0.1/0.2 history and the board descriptor's
single ledger. It checks card chains, current paths, unique metadata/IDs, CardState
placement and local Ref resolution. The board's `projectId` and `namespaces`
declare the local namespace set (including SDL/SDUI on a shared board). A `primary`
ID must match `KB-<PROJECT>-<number>` with an uppercase alphanumeric project
starting with a letter and at least three digits. Within that local set it must
resolve to a card; `reference` then contains its canonical local node ID.

A primary outside that set is represented by optional `externalReference`, whose
value is the authored card ID. This is an **unverified external reference**, not a
node ID, path, URI, promise of existence or selectable foreign target. Consumers
may show an external badge/label; they must not expand it as a local `reference`.
The containing local card remains available and its `target` still opens that
local card. No checkout scan, filesystem traversal, network lookup or external
configuration is used. External references do not excuse invalid local references,
CardState, placement or history. Older consumers may ignore this additive field;
all emitted `reference` values still resolve locally.

It never reads archived copies as extra events.
This is a read-side consistency check, not the full management validator. Other
pinned profiles, including XFMD's separate board contract, remain explicitly
unsupported until adapted. Optional SprintId/ScrumId are grouping facts, not proof
of implemented features. Metadata files/history are bounded to 1 MiB each.

Discovered SDUI sources currently support the existing sdui/0.2 Go parser/normalizer
and structural Markdown export. Frame entry nodes include target.entry for
explicit selection. `sdui-preview --model ID --entry FRAME --output DIRECTORY
--revision HASH` returns a caller-owned entry.md/provenance/manifest bundle.
It is the existing structural layout dump plus Markdown content and static widget
labels, not interactive controls, SVG/Fyne runtime or full Mermaid rendering.
Other SDUI profiles are visible as unsupported. This intentionally narrow service
is not an assertion that every existing SDUI export is exposed by this facade.

Review limits: navigation accepts at most 2,000 model declarations, 10,000 model
facts, 20,000 combined nodes and 32 MiB of inventory JSON. These are facade limits,
not new language rules. No authored defaultModel binding exists.

## Installed profile facts — IU3

Discovery accepts installed manifest schemas 1.0 and 2.0; project manifests remain
1.0. Version 2.0 adds processProfile, managementProfile and configurationDigest.
The closed shape and profile identity are checked; the full Toolkit validator
owns installation conformance. Navigation is derived from source files and does not copy installed facts. An active/failed installation journal
makes installation.state incomplete (including when final facts are not present).
Before navigation has been published, discovery returns status incomplete with
operation identities and no invented inventory/models. Tree navigation is
unavailable until the installation is complete. A completed journal is a declared
status, not independent verification.

Metadata remains limited to 1 MiB per file; operation journals have a separate
64 MiB read limit because they contain exact plans and recovery payloads.
Unsupported journal states/identities are errors. No discovery call resumes or
repairs an installation.

sdptool --version --json returns JSON build version/revision and installedFactSchemas.
A consumer requiring profile 2.0 should check that array before invoking discovery.
A source build is marked unknown or dirty when exact committed provenance is
unavailable; it is never advertised as a published Toolkit release.

package.sh NEW_OUTPUT_DIRECTORY is a maintainer-only native build. It emits a
prebuilt sdptool, its version/capability manifest and SHA256SUMS. Select a Go
compiler with SDP_GO when needed. Verify checksums and --version --json before putting
the executable on the host's configured PATH. The package does not modify PATH,
overwrite an existing destination, install a viewer or build during viewing.
The Go installation engine now owns the new install/upgrade path. The separate
[installation record contract](install/Records.md) defines preview/apply/resume,
receipt 3.0, signature provenance and explicit development/test selection.
Legacy facts 1.0/2.0 and journals remain readable; a pending legacy operation must
receive a separately assessed recovery/migration before upgrade. The current
distribution does not execute a retired engine or reinterpret its journal. No discovery command performs installation.

Typed-planning update: installed management facts and KanBan descriptors accept
sdp-project-management/0.2 as well as 0.1. The facade still projects cards; it
does not claim to render plans or validate every management transition. The selected process folder convention remains sdp-five-phase/0.1. sdp.planning.v1 is an installed process
capability, not a new navigation command.

## Sessions browsing — SN1

A real SDP/Sessions directory adds `inventory.sessions: "SDP/Sessions"` and
`capabilities.sessions: "discovered"` to discovery. The Sessions root (`sessions`,
kind `tab`) reuses canonical Files nodes, with nested directories, filename labels
and typed `open` targets. Guides/templates and other actual entries remain
browsable; no filename convention, registry or Session metadata validation is
required. This is an additive contract 0.2 extension. Existing roots retain their IDs.

The tab reports absent, empty, available or unavailable. A file or symlink at
SDP/Sessions is unavailable as a Sessions root; the generic Files tree still
shows it. Enumeration never follows symlinks. Unreadable/oversized Markdown
remains visible with a diagnostic and no open target. Session Markdown is bounded
to 1 MiB per file; readable content hashes bind open targets and the inventory
revision, including when timestamps are preserved. Other files retain normal
filesystem revisions. Existing combined navigation and scan limits still apply.

Refresh discovers additions, edits, moves and removals. `discover --json` supplies
this tree to a consumer buffer; `tree` displays it in the console. Viewer-owned
opening/watching and Session lifecycle/capture/timeline automation are separate.

## Runnable SDUI programs — RSP1

[Programs.md](Programs.md) defines sdp-programs/1 declarations, additive program
discovery and `run --program ID [--revision HASH]`. Discovery stays read-only;
only explicit Run executes project code. This is separate from generic prototype
readiness and supplies a consumer contract for XFMD without modifying that client.

## Blueprint diagnostic producer — BPI2-M1

The model create blueprint command takes from/to refs, --entry, --task TASK.json
and --output DIR. It uses owned ModelGovernance captures and the pure SDL analyzer.
The result schema is sdp-blueprint/1, operation create-blueprint, status
diagnostic-preview; human output is default, --json selects the machine envelope.
See blueprints/README.md for pinned identity, source links, strict JSON task input
and optimistic publication freshness. Preview success does not imply assignment
readiness or implementation conformance. Retained catalogue discovery is specified below.

## Retained blueprint discovery — BPI2-M2

inventory.blueprints is SDP/Blueprints when a real directory exists;
capabilities.blueprints is discovered or absent. navigation.roots includes
blueprints, with kind tab and label Blueprints. Catalogue states include absent,
empty, available and unavailable. Children group blueprint tasks then revisions.

Revision nodes expose blueprintId, blueprintRevision, retainedRevision, preliminary
and workState. workState is unknown in this milestone. validated means recorded
inventory/identity integrity only. No lifecycle, code conformance or live-source
freshness is inferred. Generated Markdown/Mermaid targets use operation open,
absolute path and file-content hash. Invalid entries have diagnostics and no target.
Tree uses the same projection as discover; refresh is consumer-owned. Retained
source copies are not treated as current-model declarations.

Creation chooses --output for a replaceable preview or --catalogue for immutable
retention. The existing result revision identifies blueprint metadata; optional
retainedRevision binds the entire stored bundle and names its revision directory.
See blueprints/README.md. No manual index or automatic XFMD update is required.

## Blueprint evidence assessment — BPI3-M1

model assess blueprint --bundle RETAINED-DIR --evidence FILE evaluates a verified
retained revision against strict sdp-blueprint-evidence/1 input. The response schema
sdp-blueprint-assessment/1 has operation assess-blueprint, status ready/blocked,
identity, exact blueprint/retained/evidence digests, scope, independent unknowns,
reported/effective check states, blockers and observed file hashes. Human rendering
uses the existing presenter; --json emits the envelope. A blocked assessment exits
3 after emitting the result. Malformed or unpinned inputs fail with diagnostics.

No commands execute and no retained bundle or canonical workflow event changes.
Ready applies only to the attributed scope/dispositions; source freshness is
not-evaluated and authority is attributed-only. Blueprint workState in discovery
is unchanged. See blueprints/README.md for input fields, budgets and limitations.


## Blueprint assignment lifecycle — BPI3-M2b

`sdptool SDP model assignment apply --request REQUEST.json --as ACTOR --authority
controller|assignee|reviewer` records one revision-bound transition in canonical
ProjectManagement history. `sdptool SDP model assignment list` is read-only. Both
use sdp-blueprint-assignment/1, human presentation by default and --json for clients.
The trusted local CLI attributes callers; it is not remote authentication. Unknown
request keys cannot override authority. Stable operation IDs and expectedEvent
enforce idempotence and optimistic concurrency. See
[the lifecycle contract](../SDP/04--Design/SDPTool/Blueprints/Assignment-Lifecycle.md)
for the exact request, evidence, authority, resource and transition rules.

The Blueprints tab adds assignment-state groups without removing retained task and
revision nodes. Assignment nodes carry assignmentId, assignmentRevision (last event),
workState, assignee, sourceFreshness, readinessStatus and evidenceStatus separately.
They never aggregate multiple assignments into a guessed task status. The older
revision nodes retain workState=unknown. Broken history is a diagnostic, never an
empty successful assignment list; damaged bundles remain visible without targets.
Retained output bytes remain immutable and no manual navigation index is introduced.

Writes currently require Linux. Release packaging/installable capabilities and
native XFMD integration are separate from this source implementation.


## Actions and invocation — ACT

`actions` emits a readable catalogue or sdptool.actions/1 with --json. `invoke`
accepts strict sdptool.invoke/1 on stdin/regular file and always returns a machine
envelope, preserving domain results and meaningful nonzero exits. The same typed
handlers serve included human commands. [Actions.md](Actions.md) owns the public
protocol detail and first five supported operations. This is an additive local
subprocess/Go interface, not an HTTP/MCP endpoint or new caller authentication.

## Viewpoint overview targets — VP01 and VP07

A viewpoint node may carry `overviewTarget` with the same Target shape as `target`.
Its presence advertises an implemented aggregate projection, not bulk detail
selection. VP01 and VP07 provide select targets whose producer-owned URI contains
`projection=overview`, with the existing exact project/model/revision. Consumers
pass that URI unchanged. Expansion never selects or renders it. Absence means
overview unavailable; older consumers may ignore the optional field. Existing
viewpoint and detail targets retain their original semantics.

VP01 emits the complete declared actor/use-case inventory, including isolated
declarations, as usecase-beta notation. Associations mean modeled `pursues`.
Functionality/Feature realization stays in a companion table; no System boundary,
UML include/extend, scenario order or pre/postcondition is inferred.

VP07 retains every declared System/container/unit and structural contains edge.
Feature keys annotate actual owners and selected-mode containers through explicit
Functionality contributions. The legend covers every Feature, and the companion
table and Diagram.Elements retain contribution/proof IDs. Features without a
represented placement remain standalone rounded nodes. Annotations are derived
views of source facts, not invented direct SDL relations.

The VP07 parent target is logical and excludes allocations. Discovery adds a lazy
child target for each declared mode; only its explicit allocated-to facts appear
in that mode's projection. Missing allocation is reported as a gap. Never combine
modes or infer allocation from ownership. Feature keys are local to the snapshot.

SDL rejects unsupported projections/viewpoints and combinations with focus,
diagram, relations, level or non-default traversal filters. Only VP07 accepts a
mode, and that mode must be declared. Snapshot freshness, cancellation, publication
and resource ownership rules still apply. Native-profile limits reject a complete
oversized projection explicitly rather than truncating it. Discovery advertises
capability without generating output. Renderers must support the emitted notation;
this contract does not claim compatibility with every external Mermaid executable.

Implementation/integration record: [ViewpointNavigation](../SDP/05--Implementation/SDL/ViewpointNavigation/Plan.md).

## Generated SDL source links — N4

Saved-source select/preview results may add sourceLinks, an array of unique
bundle-local {id,path,snapshot,revision} records. Path is the absolute original
.design file; snapshot is a safe relative path inside the returned bundle;
revision is SHA-256 of that file's exact bytes, not the composed model revision.
Generated Markdown uses sdl-source://ID/START/END with zero-based UTF-8 bytes and
an exclusive end. Every checked source is retained so clients can reject an
outdated composed graph even when the clicked file itself did not change.
Source snapshots and full graph freshness are checked before publication.

Clients validate their retained project/bundle context, containment, identities,
byte bounds and unchanged original bytes before selecting text through normal
editor transactions. Missing metadata means no exact-link capability. No source
parser, blueprint lifecycle or owner-decision behavior is delegated to the viewer.
Pure language output without filesystem binding retains text provenance.

VP01 focused use cases use actor/ellipse notation with separate realization.
Only explicit exercises relations attach ordered VP08 scenarios. Missing scenario
associations and pre/postconditions remain unspecified, never inferred.

## VP02 architecture overview — N5

VP02 additionally advertises the same optional overviewTarget with
projection=overview. It selects all declared Systems, Containers and Units plus
all explicit contains facts from the complete checked source graph, including
isolated roots. A source-linked structural register accompanies the flowchart.
Empty architecture is explicit. No allocation/ownership/deployment inference;
no new schema or fixed diagram-count cap. Original VP02 roots/parent detail targets
remain unchanged. Like VP01, no mode/focus/relation/diagram/traversal filters are
accepted for this complete overview. Existing sourceLinks/freshness rules apply.

## VP03 responsibilities overview — N6

VP03 advertises the optional overviewTarget using projection=overview. Complete
System/Container/Unit decomposition carries sorted Rxx responsibility annotations
on immediate owners, Cxx contribution annotations derived from owns plus realizes,
and separate Cxx offered-capability annotations from explicit provides. The linked
register retains every Functionality and Capability; unplaced capabilities remain
standalone nodes. Realization means contribution, never completeness or verification;
provides does not imply realizes. Every annotation retains exact supporting facts.
No feature/allocation inference, schema change or count-only truncation. Empty
inventories are explicit. Original per-capability/provider detail targets remain
unchanged; mode/focus/relation/diagram/traversal filters are rejected. Existing
composed-source snapshots, exact links and freshness rules apply.

## VP04 interface overview — N7

VP04 advertises the optional overviewTarget with projection=overview. It must
retain every declared interface, every explicit consumes fact and its consumer
(including Containers). Unconsumed interfaces remain visible. Dashed consumes
arrows describe design-time use, not observed calls or required availability.
No provider, Channel, transport, call ordering or mode necessity is inferred.
Source-linked inventory/register accompanies the unchanged Interface use table.
Empty inventories are explicit. Other structural units belong to VP02. Existing
sourceLinks/composed-source freshness rules apply; no schema change or new
count-only cap. Complete overview rejects mode/focus/relation/diagram/traversal
filters and preserves ordinary VP04 table selection.

## VP05 mode dependency overview — N8

VP05 overviewTarget without mode returns a declared-mode index, counts and source
links, never a union dependency graph. Advertise an existing-style overview-mode
target for every declared mode. With mode, select all Capabilities and the exact
mode's requires facts/Interfaces. Empty capabilities/requirements are explicit.
No inheritance, sufficiency, provider, consumption or runtime-flow inference.
Preserve existing detail IDs, VP07 targets and composed source freshness. Reject
unknown modes and other partial-overview filters; add no count-only truncation.
Consumers must retain viewpoint as well as mode when identifying mode rows.

## VP06 delivery-plan overview — N9

Advertise VP06 overviewTarget for all Activities with refines, a separate explicit
depends-on graph, source-linked status and addresses/owns/delivers/illustrates
registers. Retain isolated/statusless activities and unaddressed Functionality
inventory. Implementation status is a source claim, never completion evidence;
no inherited status/dependencies, scheduling, dates or scenario execution is
inferred. Empty cases are explicit. Keep existing detail and implementation.md
behavior; reject partial overview filters, preserve composed-source freshness
and add no count-only truncation or native/schema changes.

## VP08 communication overview — N10

Advertise an on-demand VP08 overview index and explicit named-mode overview targets.
The index lists all channels/modes/scenarios without a merged graph. Per-mode output
retains every channel and only explicit participants of that mode, grouping topology
edges by channel/participant/role while preserving all supporting source facts and
individual message bindings. Scenario and MessageSet scope is mode-isolated. Existing
scenario sequence targets remain unchanged. Empty scope is explicit; no timeline,
mode inheritance, transport or runtime success is inferred. Composed-source links,
revision freshness and generic consumer contracts apply unchanged.

The selected mode page partitions topology into one diagram per Channel. Every
channel neighborhood and supporting uses fact is retained; participants repeated
across diagrams keep the same declaration identity. This is a semantic partition,
not a size threshold or truncated graph. It avoids an all-channel union that cannot
currently be routed within the native renderer budget. Sequence targets and the
complete source-linked registers remain unchanged.

## VP09 data/storage overview — N11

Advertise a complete on-demand VP09 overview. Partition diagrams by Dataset with
explicit origin, holders, Database owners and data contracts. Include all data,
Database, Contract, Variant and Field definitions in source-linked registers,
including idle/unused definitions; retain explicit shape membership and ternary
projection/owner proofs. Detail contract diagrams remain independently selectable.
Do not infer SQL/ER semantics, persistence for Unit-held data, transformations or
wire encodings. Missing inventories/properties are explicit. Existing source and
freshness guarantees apply; no silent content truncation or new consumer schema.

## VP10 encoding/packet overview — N12

Advertise complete on-demand VP10 overviewTarget using the existing model/revision/
URI contract. Generate one separate validated packet per Encoding plus source-linked
Encoding and ordered field registers (bit offsets/widths, type/presence and header/
variant owner). Include Variant coverage and explicit missing Encoding/empty states.
Byte/bit order and closed Contract are source facts; no serialization, inferred
padding, optional-field layout or union packet. No fixed diagram-count cap. Existing
detail targets remain unchanged. Complete overview rejects narrowing filters/mode.
Native product needs no specialized VP10 handling; exact source references use the
existing sdl-source contract. Source graph freshness applies to every included file.
