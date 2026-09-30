# SDL source composition — design-core 0.6

Implemented in the Go frontend under [PLAN-SDP-0012](../../../SDP/05--Implementation/SDL/SourceComposition/Plan.md).
This is a bounded language profile, independent of the SDP release version.
It extends [design-core 0.5](SDL-Structural-Core-Profile.md) with System and
source-owned composition. Existing 0.5 files retain single-file behavior.
Action-core, class-core and SDUI keep their separate profiles.

## Authoring a System

Select one entry file, conventionally System.design:

```text
language design-core version 0.6.
system MVP1.
includes "Features.design".
MVP1 contains Containers/MachineService.
```

Containers/MachineService.design:

```text
language design-core version 0.6.
container MachineService.
```

Features.design:

```text
language design-core version 0.6.
feature Operation.
```

Each file has its own profile header. Declarations come first, then includes,
then facts. Canonical formatting sorts declarations, includes and ordinary facts
within each file; it preserves file boundaries and source dependency directives.
Keywords are lower-case. Identifiers retain the structural profile's spelling.
The compact `system MVP1 contains ...` sentence is not supported.

`includes "relative/path.design".` loads that exact file. All paths, including
those written in child files, are relative to the selected entry's directory.
`contains Containers/MachineService` additionally loads that path plus `.design`
and resolves its declared MachineService. The final identifier must be declared
in that file; folder names do not determine declaration kinds. A path operand
has at least one slash and identifier segments; a bare identifier is a reference,
never an implicit file search. Use quoted includes for other filenames.

An expanded `includes "Containers/MachineService.design".` and `MVP1 contains
MachineService.` express the same semantic membership. They have different
source bytes and therefore different revisions. No separate design-set file,
folder scan, glob, text macro, pragma or header guard is required.

## Composition and model semantics

Exactly one System is declared, in the selected entry. All reachable declarations
share that System's namespace. Names must be unique across files. A file can hold
several declarations; explicit System membership need not repeat every loaded
name. System contains accepts any supported declaration kind except System.
It records design membership, not Unit ownership, allocation or execution.
Existing Unit relation/type rules remain in force.

Repeated includes and source cycles reuse each file once. All dependency edges
are retained. Forward references resolve after source discovery. Repeated edges
can produce nonfatal REUSED_SOURCE warnings. Two distinct files declaring the
same name still fail, as do duplicate facts and semantic containment cycles.
An included second System is an error. Cross-System imports, exports and visibility
are not implemented by this profile.

Paths must be clean relative `.design` paths, exact case, forward slashes, at most
1024 bytes. Absolute paths, parent traversal, backslashes, control characters,
colon, symlinked components beneath the selected root, hardlink aliases and
non-regular sources are rejected. The filesystem root is resolved once.
Per graph limits are 128 files, 2 MiB of source bytes and 250,000 grammar tokens.
They also apply to pure in-memory compilation. Failed graph compilation returns
at most 100 diagnostics with at most eight related locations each. Related
locations are a bounded cross-file neighbourhood, not an exhaustive causal trace.
Repeated-source warnings are capped at 100; graph edges are still retained.

## File ASTs and checked snapshots

`parser.SyntaxCache.Parse(name, text)` performs no file I/O. Its immutable File
stores original text and syntax; Model returns a copy with original-source spans.
The bounded cache reuses syntax by exact bytes and parser identity. It never
caches semantic bindings. Files parsed before a parent can be passed with the
later root to `sourcegraph.Compile(entry, files, canonical)`; only reachable files
participate. Unreachable cached declarations cannot satisfy unresolved references.
Every compilation establishes and validates its own context.

`sourcegraph.Load(entry, cache, canonical)` owns filesystem discovery and produces
a checked Snapshot. Syntax, source dependency graph and linked model are separate.
AST output exposes source files, dependency edges and resolved symbol decorations;
semantic consumers use the checked model. This is not dependency-selective semantic
recompilation, public cross-System linking or runtime hot reload.

Fragment inspection returns `sdl-fragment/1`, `validated: false` and
`state: parsed-context-required`. It reports syntax and dependencies without
pretending a missing System context is resolved. A normal check on such a fragment
fails. A fragment never searches for a parent project or System.

## Revisions and document consumers

0.6 revisions use SHA-256 over length-prefixed domain, profile, entry and sorted
root-relative path/content-hash pairs (domain SDL-source-graph-revision/1).
Any reachable source change invalidates the revision; unlisted neighbours do not.
Relocating the complete tree preserves revision. Moving an individual file changes
provenance and revision; semantic fact IDs depend on System and resolved fact,
not filename or source order. 0.5 keeps its prior raw-source revision and fact IDs.

Viewpoints, CLI export/selection, broker and SDPTool load the same graph. Bundles
include sources.json with hashes, dependency edges and original declaration spans;
facts retain original-source locations. SDPTool registration selects a `.design`
entry, design-core/0.6 and its declared System, checked by tree, select and ViewPlan.
No registration migration is automatic.

Consumers protect every input from output-directory replacement and recheck the
whole graph before publishing/delivering a current result, including broker cache
hits. Invalid or stale refresh leaves the previous visible bundle/lease usable.
Freshness is optimistic: the filesystem is not locked into an atomic snapshot.
Native XFMD rendering and UI behavior have not been verified by this delivery.

## Commands and evidence

From SDL/go, using the [three-file Frontend pilot](../../go/examples/source-composition/System.design):

```sh
go run ./cmd/sdl check examples/source-composition/System.design
go run ./cmd/sdl ast examples/source-composition/System.design
go run ./cmd/sdl fragment examples/source-composition/Features.design
go run ./cmd/sdl format examples/source-composition/System.design --file-map
go run ./cmd/sdl viewpoints examples/source-composition/System.design --output /tmp/frontend-navigation --project frontend
go run ./cmd/sdl view examples/source-composition/System.design --uri 'sdl-view://frontend/VP02' --output /tmp/frontend-selected
```

Multi-file formatting requires --file-map and emits sdl-formatted-files/1 JSON;
it never rewrites source files. Normal single-file formatting still emits text.
Filesystem dependencies require a named entry; stdin cannot resolve external files.
The full experimental MVP1 corpus contains additional unsupported language syntax.
The pilot preserves the existing Frontend facts without moving authoritative models.
See [verification and independent review](../../../SDP/05--Implementation/SDL/SourceComposition/Evidence.md).
