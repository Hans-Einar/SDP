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
