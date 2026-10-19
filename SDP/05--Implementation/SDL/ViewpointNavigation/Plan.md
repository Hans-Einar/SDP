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

## N3j current-main compatibility — Session0005 T019, 2026-10-18

Owner explicitly selects conflict reconciliation before further feature work.
BranchPolicy: isolated `sdp/xfmd-viewpoint-overviews-main` from fetched main
`e8760d3`; CommitPolicy: one scoped milestone commit. Preserve cc7b52f and its
branch, the main checkout and all shared dirty work. Apply only cc7b52f's delta,
not its old ancestor history. The owner instruction covers the two concrete
shared-contract conflicts; it does not select main merge, push or release.

Preflight follows Concurrent-Work-Preflight.md. Relevant clean main/Actions heads
are e8760d3, clean N3i is cc7b52f, and shared producer 0918098 retains dirty work.
Session0011 records completed Actions/main integration. Capture exact file hashes
and statuses in the external Session's N3j-001/preflight.json. Green for this
isolated reconciliation; repeat on changes to participating heads or write scope.

Retain main's blueprint/assignment Node fields and complete producer contract;
add overviewTarget and its projection behavior. Test full SDL and SDPTool modules,
vet, new discovery-through-Actions-to-selection compatibility, and the native XFMD
keyboard/pointer/mode/PDF workflow. Compare actual XFMD overview diagram bytes with
N3i. Record final diff/input/binary identity and a clean merge forecast against
freshly fetched main. Scope remains the existing VP01/VP07 overview behavior;
focused views, exact source links and other viewpoints remain successors.

### N3j outcome

The additive resolution passes full SDL and SDPTool module tests, vet, the new
Actions/discovery overview regression and native XFMD/PDF integration. Actual
XFMD diagrams retain identical bytes. [N3j](N3j.md) and N3j-candidate.json record
the exact checked boundary. Commit this main-based milestone; retain the original
branch as history. Shared board/ledger reconciliation and remaining viewpoint
scope are unchanged. No release or main mutation is selected.

## N4 execution — XFMD Session0005 T021, 2026-10-19

Owner requests completion of remaining Session work. Continue the main-based
successor branch with milestone commits. Select focused VP01 native notation,
explicit exercises-linked VP08 companions and source-bound provenance transport.
The adopted result/source link contract is in SDPTool/Contract.md; exact source
selection remains XFMD's responsibility. Preserve PG/controller, Actions registry,
blueprint lifecycle and production .design sources. Verify single/composed inputs,
stale included files, full affected modules and actual native consumer behavior.
No merge, publication or install follows from this implementation selection.

### N4 implementation outcome

[N4](N4.md) records completed selected detail/source transport, full producer and
native verification, recovered diagram-count behavior and explicit residuals.
This milestone does not close broader KB032 or integrate/release the producer.

## N5 execution — XFMD Session0007 T001, 2026-10-19

Owner authorizes integrating the existing delivery and continuing to the next
viewpoint. N3j/N4 merge 5ca0e90 is pushed to main; its tree equals tested fde2239.
BranchPolicy: sdp/xfmd-vp02-overview from 5ca0e90. CommitPolicy: scoped milestone.
Select VP02 complete logical structure (System/Container/Unit plus contains),
including isolated declarations; no allocation/ownership inference. Advertise
through existing overviewTarget and preserve old roots/parent detail identities.
Keep a source-linked structural register, explicit empty behavior and query
filter rejection. No new count-only cap or guessed parent. Production model hold
remains. Verify full inventory/proofs, composed freshness, Actions/discovery and
native XFMD/PDF. See N5.md for the design and final evidence. No N5 merge/release
or installation is selected by the existing-delivery integration authorization.

### N5 outcome

VP02's complete logical overview and source-linked register are implemented and
verified through producer and native XFMD. See [N5](N5.md) and its candidate record.
N3j/N4 is on main 5ca0e90; N5 remains its own milestone for later integration.
Other viewpoints and previously recorded publication/resource obligations remain.

## N6 execution — XFMD Session0007 T002, 2026-10-19

Owner continuation selects VP03 responsibilities and capabilities. BranchPolicy:
sdp/xfmd-vp03-overview stacked on N5 cdf4437; CommitPolicy: scoped milestone.
[N6](N6.md) selects the architecture annotation/register projection before code.
Keep owns, realizes and provides distinct, retain all inventories and proof links,
and preserve detail views. Concurrent preflight finds no PG/coordination overlap.
Production models remain held. Verify full modules, native generation/rendering
and actual XFMD responsibilities without invented capabilities. No merge/release
or install is authorized by this continuation.

### N6 outcome

VP03 complete architecture annotations and responsibility/capability registers
are implemented and verified through producer and native XFMD. [N6](N6.md)
records the exact candidate and evidence. N6 stacks on N5; neither is merged or
released. Remaining viewpoint/publication work stays open under the same plan.

## N7 selection — Session0007 T003, 2026-10-19

Owner continuation selects VP04 interface consumption. BranchPolicy:
sdp/xfmd-vp04-overview stacked on N6 63da7ae; CommitPolicy: scoped milestone.
Complete interface inventory with explicit consumes dependency graph/register,
no inferred provider/transport/necessity; preserve ordinary VP04 table selection.
N7 design selects the projection before code, with full producer/native evidence.
No authored model changes or merge/release/install; broader plan remains active.

### N7 outcome

VP04 complete interface/consumer dependency diagram and source-linked register
are implemented and verified through producer/native/PDF workflows. [N7](N7.md)
records candidate, tests and the preserved keyboard semantics. Other viewpoint
and large-publication obligations remain. N5/N6/N7 are stacked, unmerged candidates.

## N8 selection — Session0007 T004

Owner continuation selects VP05 mode-index and isolated per-mode dependencies.
BranchPolicy: sdp/xfmd-vp05-overview from N7 85d62f2; milestone commit. N8 design
precedes code, preserving original details, VP07 overlays and authored models.
Verify producer/native/PDF workflows and actual XFMD's empty NormalDesktop scope.
No merge/release/install; broader plan and cards remain active.

### N8 outcome

VP05 index and isolated named-mode overview are implemented and verified through
producer/native/PDF workflows. [N8](N8.md) retains exact evidence and limitations.
VP07 and original details are preserved. Broader scope remains; N5–N8 are stacked
candidates, without merge/release/install.
