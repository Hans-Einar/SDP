# VP01 overview producer work unit — N3d

Owner continuation in XFMD Session0005 T004, 2026-10-08, selects this producer
work unit under authoritative PLAN-XFMD-0009. This is not a second lifecycle plan.
KB-SDP-032 owns producer feedback; XFMD KB028 owns native consumption. Use the
current non-main branch with milestone commits. Preserve the concurrent shared
dirty baseline; defer a root-wide milestone commit and preserve scoped patches
and candidate source hashes. No push, merge, release or extension installation.

## Design and acceptance

Refines REQ-SDPTOOL-003/004: a viewpoint has an optional overviewTarget using the
existing typed target. Only VP01 advertises it in N3d, using the new SDL query
`projection=overview`. Separate query syntax avoids collisions with existing
user-derived diagram IDs such as VP01-overview. Old target/detail semantics stay
unchanged; unsupported projections or overview modifiers are rejected explicitly.

SDL owns lazy VP01 projection: all declared actors and use cases, including isolated
items; pursues means an ordinary association, with a legend. Do not invent System
boundaries or UML include/extend. Functionality realization remains a companion
table of supports/contributes-to facts; no flowchart is called a use-case diagram.
Retain snapshot revision, model IDs and declaration/relation provenance. Empty
models return truthful inventory/gaps without a fictitious diagram. Existing limits
reject oversized aggregates visibly, never truncate and claim completeness.

Milestone N3d: implement SDL projection and documents, SDPTool advertisement and
XFMD typed-target activation. Verify 25+ use cases, isolated actors/use cases,
empty data, unsupported filters, foreign project/stale revisions, old producers,
no generation during expansion, native click/keyboard selection and real-tool
native/PDF output. VP07 mode/architecture projection and focused scenario/condition
views remain successor work; exact source-span navigation remains N4.

Session and evidence:
file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230005--SDL-viewpoint-diagrams.md
file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N3d-001/README.md

## Process ownership

The active cross-project goal is the linked XFMD Session0005; its T004 journal is
updated for this substantive producer turn as well as consumer work. This work
unit and KB032 link that authoritative Session rather than creating a competing
Session/roadmap in the producer repo. No host transcript or owner acceptance is
inferred. Other active producer Sessions, including governance and SDUI widgets,
retain their independent scope and concurrent work.

Visual refinement: native inspection of the real XFMD model rejected the wide
TB fan-out. VP01 overview uses LR to favor a vertical use-case column. Transport
query JSON remains in selection.json but is omitted from the overview's reading
page. Preserve revision/provenance and all facts; no source-model relation changes.

## N3d outcome

N3d delivers the first VP01 overview through SDL projection, SDPTool overviewTarget
and native XFMD activation. Complete actor/use-case inventory includes isolated
items; functionality realization stays separate. Real XFMD source retains one
actor, 11 use cases and 11 pursues relations. Native inspection selected LR after
rejecting the unreadably wide TB fan-out. Query metadata stays in selection.json.
The tall real diagram also exposed a one-ULP paged-height overshoot; DiagramPlacement
now bounds the final height and has a synthetic plus real PDF regression.
Go viewpoint/documents and SDPTool facade tests pass. Native and affected consumer
checks pass, including Enter/click, expansion-only, older/empty producer fallback,
stale/foreign targets and PDF export. The first new GUI attempt lacked explicit
X11 focus; its failing log is preserved alongside corrected passing evidence.

Evidence and exact candidate are in the external Session N3d-001 link above. No producer release or gh installation. VP07 and detailed views remain successors.

## N3e selection — Session0005 T005, 2026-10-08

Deliver VP07 as a complete logical architecture base (all declared System,
container and unit nodes, including isolated nodes), overlaid with every Feature
and its contributing Functionality. Preserve actual contains, owns and
contributes-to edges, with distinct feature/functionality shapes and relation
styles. Do not fabricate placement from ownership. The default overview is
explicitly logical and omits allocation edges; discovery adds one lazy overview
selection per declared mode. A selected mode includes only its allocated-to
facts and explicitly reports unallocated contributions. Never union modes.
Existing feature detail targets remain usable. This uses the N3d typed target;
no XFMD parser or hardcoded viewpoint list is introduced.

A single union of all allocations is rejected because modes may be incompatible.
Generating all mode diagrams on activation is rejected in favor of on-demand
selection. Full structural base means decomposition, not every unrelated
interface/data/scenario relation. All declaration/relation provenance remains.
Verify isolation, multiple modes, unallocated/isolated features, empty models,
old VP01 behavior, actual discovery-to-native selection and real XFMD PDF.
The current branch/shared-baseline policy from N3d remains; no release/install,
model edit or broad commit. Exact source links remain N4.

### N3e overlay refinement after real-model acceptance

The literal realization graph (57 nodes / 137 edges) exceeded the native capture
budget on the real XFMD model. Replace Functionality nodes/edges in the overview
with deterministic Feature keys (F01...) annotated on their actual logical owners
and, for a selected mode only, on explicit allocation containers. The complete
structural contains graph remains. Unowned/unrealized Features remain standalone
rounded nodes. The legend maps every key to its declared Feature; a companion
table retains all contributing functions, owner/allocation and exact proof IDs.
Diagram Elements retain annotation provenance (Feature, owner/container, role,
source-fact IDs), and SourceFacts includes these proofs. This is a traceable
projection, not invented direct SDL relationships. No synthetic edges are added.
Detailed feature selections retain their explicit graphs. The initial literal
candidate and timeout are retained; all-mode union remains rejected.

## N3e outcome — Session0005 T005

N3e implements lazy VP07 logical overview and named-mode targets. All declared
System/container/unit structure is retained. Feature keys annotate real owners and
selected-mode containers through proven Functionality paths; the legend and table
retain all Features, contributions and allocation gaps. No cross-mode union or
invented direct relation. Native Enter/selection exercises both mode targets.
Ten affected CTests pass (five GUI); SDL viewpoint/documents pass and the SDPTool
facade suite passes with one documented pre-existing SDUI test excluded. That test
also fails in an isolated copy with N3e edits removed: UIPreview uses Combined,
while the all-family test compares against Markdown. Do not claim full producer
suite success or change another workstream's acceptance in this increment.

The literal graph timed out; the adopted annotation projection renders the real
XFMD model (23 architecture nodes, 19 contains edges, 140 proof annotations and
137 facts) and exports an eight-page vector PDF. Small-fixture rendering is readable.
Real-model full-view labels remain small at fit-to-width / A4. Grouping disconnected
nodes was explored but not adopted: it improved native width while worsening PDF
pagination. Large-overview reading is still an acceptance gap owned by KB020/028;
no overall visual acceptance or Session completion is claimed.

Exact candidate/evidence: file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N3e-001/README.md

## N3i integration preparation — Session0005 T016, 2026-10-18

Owner continues after clarification that this XFMD workstream owns the producer
changes. Select a local milestone commit in isolated worktree
`/home/warloc/git/SDP-xfmd-viewpoint-overviews`, branch
`sdp/xfmd-viewpoint-overviews`, based on `48c919f`. This bounded exception to the
earlier current-branch policy separates overview code from concurrent SDUI and
governance edits. No merge, push, release or extension install is selected.

Carry only viewpoint overview/query code, selected-document integration, typed
tree advertisement, overview tests and the overview contract. Verify SDL packages,
the complete SDPTool root package, actual discovery/select and the native XFMD
overview workflow with this exact binary. Preserve old/new target compatibility.
Review composed-source input alongside single-file models before declaring the
additive contract ready for integration. Focused use cases, conditions and exact
source navigation remain successors, not claims of this commit.

The shared checkout remains the authority for KB032's in-progress state and its
uncommitted append-only history. This isolated branch retains the committed board
baseline: importing only the older KB032 events after later committed event IDs
fails the strictly increasing namespace rule. Do not reorder/rewrite history or
create duplicate transitions to disguise that divergence. Integrators must retain
and reconcile the shared board/ledger work separately; the code commit contains
no board move or lifecycle-ledger fragment. See Integration.md for exact scope.

### N3i outcome

The isolated candidate passes SDL viewpoint/documents and the complete SDPTool
root-package suite, including added composed-source freshness and provenance
coverage. Real discovery/select and native XFMD Enter/click/mode/PDF checks pass.
See [Integration](Integration.md) for the exact boundary and downstream merge
notes. No broad worktree staging, producer install or public publication. This
milestone prepares the overview delivery; the wider KB032/Session goal remains.
