# Producer reading pages

The Activity3 reading path implements the accepted
[reading design](../../../../HEOS/SDP/03--Architecture/SDLBrowsing/NavigationDesign.md).
SDL owns semantic membership and page composition. The host supplies checked
source context, retained model/revision identity, source bindings, and rendering.
The consumer follows the producer's supplied actions.

## Selection and shared membership

`Views.ReadingCollections(viewpoint)` returns ordered collections with stable
`ID`, reader-facing `Label`, and sorted declaration `Names`. Membership comes
from declarations and explicit facts, independently of the legacy diagram list.
Consumers, Participants, Deliverables, Properties, and Model gaps are presentation
groups. They introduce no SDL grammar or declaration kinds. Relevant isolated
declarations remain present. Empty collections can be omitted from navigation.
Within one viewpoint an object has one primary collection route; page links can
open the same declared object through another contextual question.

`projection=reading` supports VP01–VP11. No selector opens a viewpoint landing.
Exactly one optional selector opens `collection=ID`, `focus=DeclaredName`, or
`register=facts|declarations`. Registers are valid only under VP11. `mode=Name`
retains an explicit declared operating mode. It filters applicable scoped facts;
global declarations remain visible with honest applicability explanations.
Traversal, relation pruning, diagram selection, and abstraction-level selectors
cannot be combined with reading selections. Legacy selection URIs remain valid.

`page=N` is one-based and valid only on collection/register selections. A page
contains at most 24 object chapters or 100 register rows, states the total scope,
and supplies Previous/Next actions when needed. Out-of-range pages fail rather
than returning an empty or truncated result. Chapters include all their scoped
statements. Large optional diagram sets are deferred to contextual object/mode
selections while complete semantic content and links remain available.

## Page composition and host binding

`documents.Selected` dispatches reading selections to `Views.ReadingMarkdown`.
Landings explain the viewpoint question, current scope, relevant collections,
reported gaps, and interpretation limits. Collection chapters and object details
share the same content component. It describes explicit properties and relevant
relationships, contract fields/variants/presence, payload contracts, channel
participants and roles, scenario steps/correlation, and encoding bit positions.
Missing relationships, unknown conditions, unspecified layout, and absent mode
applicability remain explicit. Dependencies are not execution order; source
status and completeness are not implementation evidence. Generic browsing does
not invent a source-to-code mapping.

Complete facts and declarations appear only at their dedicated reading register
selections. Registers preserve every checked statement/declaration across pages;
selected mode is retained for subsequent navigation without hiding global rows.

The host sets `Views.ReadingProject` for exact `sdl-view` action authority. Without
it, composition emits plain labels rather than inventing a project. Host-retained
model/revision context binds those actions; IDs stay in protocol metadata.
`SourceBindings` maps original span source names to retained source handles.
`SourceLabels` supplies readable filenames, including the empty source name used
by legacy single-file spans. A statement link shows its source text, filename,
and line while preserving the original byte start/end in its `sdl-source` target.
Machine fact IDs and revision hashes remain in manifests rather than reader text.

The document layer integrates the structure renderer for supported VP02/VP08
scopes. Generated geometry depicts actual declared Unit/Container participants
and Channel attachments, preserving role, payload, and mode. It does not add ports
or infer unnamed internal participants. Diagram assets are generated on demand.
