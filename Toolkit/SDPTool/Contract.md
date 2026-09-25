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
