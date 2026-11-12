# Activity3 reading contract

Selected under the owner's autonomous Activity3 implementation/publication request.
Canonical intent remains [NavigationDesign](../../../../../HEOS/SDP/03--Architecture/SDLBrowsing/NavigationDesign.md).

## Ownership and protocol

SDL owns semantic membership, projections and document composition. SDPTool owns
source discovery/context, final navigation/actions, retained snapshots and bundle
publication. XFMD owns generic protocol validation, rendering and navigation state.
No second SDL parser or per-viewpoint rules belong in XFMD.

Add `presentation: "producer-reading/1"` to navigation Tree responses (including
model-only trees), and advertise `readingNavigation: "1"` in discovery capabilities.
On this path arrays are final sibling order, labels are final, and Node.target is
the primary activation. Node.expansionTarget alone requests deferred children.
Legacy producers retain the existing consumer path; producer old selection URIs
remain supported. A retained target carries the existing project/model/revision.

The new SDL reading selection uses `projection=reading`, with optional mutually
exclusive `collection`, `register`, or existing `focus`; optional `mode` preserves
explicit mode selection. No selectors means a viewpoint landing. Collections have
stable semantic keys supplied by SDL, with ID/label/ordered object names.
Registers `facts` and `declarations` are dedicated VP11 destinations. Normal pages
never append complete registers. Query combinations are strictly validated.
Expose `Views.ReadingCollections(vp string)` returning []ReadingCollection with
fields ID, Label, Names; the navigation host uses this shared semantic authority.
All eleven viewpoints retain meaningful landings, isolated declarations and
contextual object pages. Collections render chapters using the same object content;
large selections paginate with complete scope and supplied links, never truncate.

Source statement text plus filename/line replaces machine ID labels, retaining
exact sdl-source spans and all-source freshness. No machine identities in reader
text/diagrams/tooltips; protocol identity remains unchanged. Use actual model facts
and explicitly state unknowns. Diagram files remain generated on demand.

## Discovery and structure

Systems contains actual System entry files, then supported Standalone designs and
Unassigned fragments. Each entry exposes ordered VP01–VP11 and Source files for
its checked include closure. Shared fragments occur in each valid context. Broken
roots retain diagnostics and known includes. Orphan syntax gets a clearly labeled
local declarations/relations page, never falsely checked system semantics.
Entry activation opens a coherent chaptered structure page. Channel endpoints are
projection geometry on actual Unit/Container participants, not new SDL entities.
Preserve role, payload and mode; never infer runtime order from dependencies.

## Release

Owner selected XFMD0.10. Package the exact committed matching SDPTool privately,
retain provenance, and verify signed public DNF upgrade/install on AlmaLinux10.2.
Feature-branch publication is authorized; no main merge is inferred. Preserve
unrelated primary XFMD working files. Consolidated independent review covers the
final integrated candidate and release evidence.
