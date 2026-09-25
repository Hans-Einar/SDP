# SDPTool — detailed design work

Design record: DES-SDPTOOL-001. Status: producer contract implemented in Sprint-0001; advanced services and
native XFMD integration remain later scope.
[Requirements](../02--Requirements/SDPTool.md) ·
[Architecture/model](../03--Architecture/SDPTool.md) ·
[Implementation plan](../05--Implementation/SDPTool.md).

## Existing implementation to reuse

The [pipeline study](../../Toolkit/SDPTool/Navigation-and-Design-Preview.md) records
Go parsing/projection/document publication, Rust mmdr layout, Go SDL SVG symbols,
sdl-view selection and temporary-resource ownership. These are implemented SDL
services reused by the implemented SDPTool facade. The native .design
source/preview adapter remains XFMD-owned.

## First contract: saved design preview

Specify a bounded producer operation for REQ-SDPTOOL-001/004 before code:

| Decision | Required design outcome |
| --- | --- |
| Invocation | Command name, saved source path, supported profile, source base, optional view selection, output destination and registered renderer |
| Default projection | Relevant compact overview/diagrams and explicit limits; no invented facts or unconditional full export |
| Result | Entry document, Mermaid/SVG resources, source revision, selection and resource ownership; reuse existing bundle schema where possible |
| Errors | Positioned SDL diagnostics, unsupported profile, missing renderer and publication failure; source remains intact |
| Consumer contract | Native XFMD preview keeps source path, undo/dirty/save behavior; consumer implementation belongs to KB-XFMD-015 |
| Follow-on snapshot input | Unsaved source and revision identity, cancellation/stale-result handling and lifetime; current CLI reads saved files |

The first standalone preview does not require a project descriptor, daemon or
the complete native SDP tree. Model interfaces identify collaboration boundaries;
they do not yet define these concrete fields or executable signatures.

## Later contracts

Reconcile project/install manifests for directory validity, capabilities and
explicit path bases. Specify model-derived tree identity/children/references,
revision and selection across all supported concepts. KanBan history, process
evidence and SDUI exports retain their own contracts. Viewpoint source/model work
stays in SDL; native tree work stays in XFMD KB-XFMD-014.

Add A4 Channel/scenario/data/class details only when their semantics and actual
consumer are defined. Do not introduce unsupported imports, Requirement or State
keywords to make this feature look more complete. Requirement prose is currently
Markdown; the validated architectural model uses design-core 0.5.

## P0-M1 implemented contract

The [producer contract](../../Toolkit/SDPTool/Contract.md) resolves the saved-file
operation above with a Go facade, bounded selection, language diagnostics and
caller-owned revision-tagged bundles. Unsaved-buffer and native viewer behavior
remain separately owned. P0-M1 evidence is in the implementation plan.
