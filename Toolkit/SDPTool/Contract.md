# SDPTool producer contract 0.1

This is the implemented local facade contract, not a new Toolkit release.
Go language packages own parsing, validation, projection and presentation.

## Saved design preview

```sh
sdptool preview model.design --output /tmp/request-preview
sdptool preview model.design --output /tmp/request-svg --renderer /absolute/mmdr
sdptool preview model.design --output /tmp/request-detail --viewpoint VP01
sdptool preview model.design --output /tmp/request-detail --uri 'sdl-view://project/VP02?diagram=VP02-roots&target=main&consumer=xfmd' --revision SOURCE_SHA256
```

Source precedes flags. The current input profile is SDL design-core/0.5 with its
canonical-source validation; unsupported language/profile input returns the
language diagnostic. Class/action/SDUI inputs are not structural design previews.
No SDP folder is required. Default selection is one nonempty diagram, preferring
architecture then use cases/responsibilities; an empty model uses VP11. Explicit
selections retain SDL's query validation and depth bound (0–8). A request accepts
at most 24 diagrams, 128 files and 32 MiB of generated resources; larger requests
must select a diagram/focus. No automatic full export occurs.

Stdout is one JSON result with schema sdptool/0.1, operation, source, profile,
revision (SHA-256 of exact source bytes), entry, directory and ownership. The
bundle reuses SDL entry.md, diagrams/*.mmd, optional diagrams/*.svg, selection.json,
manifest.json and delivery.txt; sdptool.json records facade provenance. SVG uses
SDL's symbols and existing Rust renderer geometry. Without a renderer, Markdown
contains Mermaid fences. Renderers are prebuilt programs supplied by the host.

Errors return nonzero and one JSON error on stderr (schema, error.code/message
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

## Project recognition — T1 contract

The local `SDP/navigation.json` registration owns **navigation bindings only**.
Its [schema](navigation.schema.json) is independent of Toolkit release numbering.
It names a project ID, the adopted local sdp-five-phase/0.1 profile, optional
implementation-plan/board/project-manifest paths and explicit model/source lists.
The mixed SDL/SDUI model is registered honestly as a joint source until source-set
migration; `system` is a display/ownership label, not new SDL syntax. Empty model
and SDUI lists are valid and expose unavailable tabs, not nonexistent language
capabilities. IDs are unique across both lists. Unknown language profiles remain
visible as unsupported and are never dispatched through a different parser.

All registration paths are relative to the **parent of the SDP area**, including
when the caller selected the SDP area itself. They must remain inside that root
also after resolving symlinks. Model selection is explicit when multiple models
exist. No recursively discovered repositories, implicit parent search, Git-based
identity, child-directory scanning or executable commands are part of registration.
Selecting a separate repo explicitly follows the same rule as a monorepo area.

Recognition checks a navigation.json in the selected directory first, otherwise
its SDP child. A selected malformed/unsupported registration does not fall through
to another project. No registration is `missing`, invalid fields/paths are
`invalid`, an unknown registration schema/process profile is `unsupported`, and
accepted registration is `valid`. This validates navigation eligibility, **not
whole-project SDP conformance, installed-tool compatibility or system correctness**.
An existing SDP folder/legacy installation without registration is reported missing
with instructions to register it; it is never silently migrated. Declared source
paths must exist and be readable before the associated operation runs.

| Existing authority | Integration decision |
| --- | --- |
| Root SDP.manifest.yaml | Distribution facts; not a consuming-project marker |
| Toolkit/SDP-install.manifest.json | Inventory/install policy; no discovery rewrite |
| SDP-project.manifest.yaml | Optional project-owned facts, referenced by registration; product release is not the process version |
| Its installed.manifestPath | Installed facts relative to that manifest's directory; never copy Toolkit/skill versions into navigation.json |
| Generated installed-toolkit manifest | Report read facts as declared; full installer/schema conformance stays with its validator |
| navigation.json | Explicit local navigation eligibility and bindings; no release, install timestamp, executable or runtime/window/lease facts |

T2 reads referenced YAML manifests safely, rejects malformed/unsupported schema
versions and reports installation facts without asserting full validation. Missing
manifests are errors when explicitly referenced; omitted manifests yield unknown
installation facts. An installer can later distribute/create bindings under
KB-SDP-028, with ownership/preservation decisions there. This local descriptor
neither installs the five-phase profile elsewhere nor changes strict 1.0 manifest
schemas; that is why a separate, narrowly scoped binding file is justified.

## Commands, host policy and consumer protocol — T1-M2

| Invocation | Result and ownership |
| --- | --- |
| sdptool [PROJECT-OR-SDP-AREA] discover | JSON recognition and declared capabilities; read-only |
| sdptool preview FILE --output DIR | Standalone saved-file operation above |
| sdptool [PATH] view ip --model ID | Open the registered authored plan with generated navigation in the configured viewer |
| sdptool [PATH] tree --model ID | JSON navigation inventory; no detail rendering or writes |
| sdptool [PATH] select --model ID --uri URI --revision HASH --output DIR | Validate project binding/revision, then generate selected current-source detail |
| sdptool [PATH] sdui-preview --model ID --output DIR | Only the implemented, registered SDUI service selected in T3 |

`implementation-plan` is an alias for `ip`. `generate ip` remains unsupported;
opening a plan must not synthesize, overwrite or imply approval of one. Explicit
path precedes the command; omission means current directory. Flags follow the
operation's fixed positional arguments. Reject extra arguments and unknown flags.

Host executable precedence: explicit --viewer/--sdl-tool/--renderer options,
then SDP_XFMD/SDP_SDL_TOOL/SDP_MMDR, then viewer/sdl names on PATH when needed.
No renderer means Mermaid output, not an automatic build/install. Resolve programs
before launch, preserve argument boundaries and never invoke a shell. Registration
files cannot register executables. No startup tool compilation or required daemon.
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

JSON responses identify schema `sdptool/0.1` and operation. Recognition reports
status, root, area, registration and capabilities; language capabilities remain
`declared` until the owning parser is run. Tree replies include source revision,
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
selected registration, plus an explicit model when ambiguous. It delegates to the
same guarded preview path. Refresh and retry after a stale error. A foreign URI
never switches project/source selection implicitly.

## KanBan and SDUI services — T3-M3

`tree` returns KanBan, SDL and SDUI roots, with unavailable diagnostics rather than
hiding failures in optional services. `revision` is the selected SDL source hash;
`inventoryRevision` also changes with registration, card/ledger data and SDUI
sources. Consumers refresh on either appropriate revision. All targets retain
source/card-specific hashes. No UI state is persisted by the inventory operation.

The initial board reader supports board schema 0.2 with the local
sdp-project-management/0.1 profile, payload 0.1/0.2 history and the board descriptor's
single ledger. It checks card chains, current paths, unique metadata/IDs, CardState
placement and local Ref resolution. It never reads archived copies as extra events.
This is a read-side consistency check, not the full management validator. Other
pinned profiles, including XFMD's separate board contract, remain explicitly
unsupported until adapted. Optional SprintId/ScrumId are grouping facts, not proof
of implemented features. Metadata files/history are bounded to 1 MiB each.

SDUI registrations currently support the existing sdui/0.2 Go parser/normalizer
and structural Markdown export. Frame entry nodes include target.entry for
explicit selection. `sdui-preview --model ID --entry FRAME --output DIRECTORY
--revision HASH` returns a caller-owned entry.md/provenance/manifest bundle.
It is the existing structural layout dump plus Markdown content and static widget
labels, not interactive controls, SVG/Fyne runtime or full Mermaid rendering.
Other SDUI profiles are visible as unsupported. This intentionally narrow service
is not an assertion that every existing SDUI export is exposed by this facade.
