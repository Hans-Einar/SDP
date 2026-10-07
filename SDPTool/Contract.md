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
selections retain SDL's query validation and depth bound (0–8). A request accepts
at most 24 diagrams, 128 files and 32 MiB of generated resources; larger requests
must select a diagram/focus. No automatic full export occurs.

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
facts, derived `inventory`, `sources`, `plans` and `navigation`. Nothing is written.
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
before launch, preserve argument boundaries and never invoke a shell. Source and metadata files cannot register executables. No startup tool compilation or required daemon.
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


## Blueprint diagnostic producer — BPI2-M1

The model create blueprint command takes from/to refs, --entry, --task TASK.json
and --output DIR. It uses owned ModelGovernance captures and the pure SDL analyzer.
The result schema is sdp-blueprint/1, operation create-blueprint, status
diagnostic-preview; human output is default, --json selects the machine envelope.
See blueprints/README.md for pinned identity, source links, strict JSON task input
and optimistic publication freshness. Preview success does not imply assignment
readiness or implementation conformance. Catalogue discovery is a later milestone.
