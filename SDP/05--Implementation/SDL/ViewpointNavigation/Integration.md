# N3i — VP01/VP07 integration handoff

Prepared 2026-10-18 under XFMD Session0005 T016 / PLAN-XFMD-0009 / KB-SDP-032.

## Current integration vehicle — N3j

Supersedes the N3i branch for integration: `ecd9dcd5cd95fb2b1f323cea96da1776659b2513` on
`sdp/xfmd-viewpoint-overviews-main`, based directly on main `e8760d331933361ab631ef0d52c20e87268a22b3`.
Only our bounded delta is carried. Both conflicts are reconciled while retaining
main's assignment/blueprint fields and Actions contract. Full SDL/SDPTool suites,
vet, Actions-to-overview regression and native XFMD/PDF checks pass. The final
merge-tree forecast is clean; the candidate worktree is clean and one commit ahead.
No actual main merge, push, release or install. Do not merge the old branch's
ancestor history or reset shared dirty copies. Shared lifecycle reconciliation
remains separate. Current record:

file:///home/warloc/git/SDP-xfmd-viewpoint-main/SDP/05--Implementation/SDL/ViewpointNavigation/N3j.md

## Historical N3i candidate and scope

Branch: `sdp/xfmd-viewpoint-overviews`, based on `48c919f`. This milestone
contains the previously uncommitted N3d/N3e overview implementation, a current
contract/requirement clarification and composed-source freshness coverage.
Exact integration commit: `cc7b52f84f05fdef63efc1a305023a12347026ad`. It does not include
the concurrent SDUI, widget, sourceinput, MCP or governance work.

Product files: SDL/go/viewpoint/{query,overview,overview_test,feature_overview,
feature_overview_test}.go; SDL/go/documents/selection.go; SDPTool/{tree,
overview_test}.go; SDPTool/Contract.md. Requirements and this work-unit folder
record the selected boundary and evidence. The contract contains only the final
overview behavior; retain the separate SDUI combined-preview section when
reconciling with the shared dirty checkout.

## Verification

- `cd SDL/go && go test ./viewpoint ./documents`: PASS.
- `cd SDPTool && go test .`: PASS, including all existing root-package tests and
  composed-model overview/source-provenance/freshness regression; no exclusions.
- `cd SDPTool && go build -o /tmp/xfmd-n3i-sdptool ./cmd/sdptool`: PASS.
- Native XFMD SdpOverviewGuiTest against that exact binary: PASS. Expand causes
  no generation; Enter/label selection and both mode targets generate correctly.
  Native screenshots were inspected; exported PDF text retains the expected labels.
- Actual XFMD model discovery/select: VP01 12 nodes/11 edges; VP07 logical and
  NormalDesktop each 23 nodes/19 edges, with respectively 113 and 137 proof facts.
- ProjectManagement validator: PASS (64 cards, 43 management records, 551 events);
  lineage examples and git diff --check pass on the isolated candidate.

Binary/source hashes are in candidate.json; detailed logs and captured bundles:
file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N3i-001/README.md

## Integrating with concurrent work

The original producer working files are intentionally retained in SDP-vNow.
Compare the scoped commit with those dirty copies before applying it; most
product hunks are already present there. The composed-model test and consolidated
contract are new in this branch. Do not import all dirty producer files, reset
the working tree, or apply duplicate overview implementations. Reconcile shared
changes before committing/merging on the integration branch.

KB032's active card and earlier lifecycle events are currently uncommitted in
the shared checkout. The isolated branch keeps its committed board baseline.
Copying the older KB032 events after later committed IDs fails the namespace's
monotonic rule; no old ledger bytes were reordered or rewritten. This code commit
therefore deliberately excludes board/ledger fragments. Preserve that complete
shared history when its owner integrates it; never synthesize a second transition.

## Remaining boundaries

This is a local integration candidate, not main integration, independent review,
a published release or an installed gh upgrade. Installed SDPTool 2.2.0 remains
unchanged. Verify the exact merged/released binary advertises overviewTarget and
executes its unchanged URI before claiming consumer availability.

KB032 stays in progress for focused use-case/scenario projections, conditions,
exact source navigation, readable overview publication and other viewpoints.
KB030/PLAN-SDP-0024 JSON action discovery remains a separate held workstream.

## XFMD Session0005 T021 — current N4 integration handoff

The main-based successor now contains ecd9dcd and
fde2239f0cd8202f6221394d2a3dae84828edb10 on
sdp/xfmd-viewpoint-overviews-main in /home/warloc/git/SDP-xfmd-viewpoint-main.
It delivers focused actor/use-case diagrams, explicit exercises-linked scenarios,
sourceLinks with checked original per-file snapshots, and exact native source
selection. Removes fixed 24/48 diagram-count guards; 50-diagram test passes.
Full SDL/SDPTool suites and native XFMD/PDF checks pass. Final merge-tree against
main e8760d3 and PG 3e407d6 is clean; no shared changed paths at precommit check.
This is a forecast, not combined-runtime proof or an actual merge.

Integrate the successor branch, not old cc7b52f ancestry. Preserve the shared
working copies and reconcile KB032 history separately; no ledger fragments were
imported/reordered. Product files in this shared checkout remain untouched.
KB032 stays active for remaining viewpoint proposals, legacy eager-export notation
and integration/release. Large-print readability and general renderer issues
remain explicit XFMD card obligations. No main merge, push, release or install.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N4-001/README.md

file:///home/warloc/git/xfmd-sdl-navigation/SDP/04--Design/SDL-Navigation/Viewpoint-Disposition.md

## XFMD Session0007 T001 — main integration and VP02 candidate

Owner authorized the existing delivery merge and continuation to the next
viewpoint. N3j/N4 is merged and pushed to SDP main as 5ca0e90. The merged tree
exactly matches tested fde2239, with targeted integration checks passing.

Next increment VP02 is committed/pushed separately as cdf4437 on
sdp/xfmd-vp02-overview: complete System/Container/Unit inventory and contains
edges, isolated roots and source-linked structural register. Original detail
targets remain. Full SDL/SDPTool suites/vet and 12 affected native consumer tests
pass. Actual XFMD model renders 23 nodes/19 edges; all labels survive PDF.
No new release/install or production model edit. N5 is not merged automatically.
VP03 is proposed next; KB032 retains its other viewpoint/integration scope.
Shared producer product files/history are preserved; only this handoff is appended.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230007--Architecture-viewpoint.md

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N5-001/README.md

## Session0007 T002 — VP03 candidate

N6 is committed/pushed as 63da7ae on sdp/xfmd-vp03-overview, stacked on VP02
cdf4437. It retains complete architecture/responsibility/capability inventories
with distinct ownership, contribution and provision annotations and exact source
proofs. Existing detail targets remain. Full SDL/SDPTool suites, vet and 12 native
consumer tests pass; actual XFMD 23 nodes/19 edges/28 responsibilities render and
export. No capability is invented for the current model. Long PDF identifiers
wrap; large-print readability remains open. VP04 interfaces is proposed next.

No N5/N6 merge/release/install, production model edit or shared producer product
change. Broader cards/plan stay active and independent review remains outstanding.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N6-001/README.md

## Session0007 T003 — VP04 candidate

N7 85d62f2 on sdp/xfmd-vp04-overview stacks on N6 63da7ae. Complete interface
inventory, explicit dashed consumption graph and exact-source register preserve
the ordinary Interface use table. No provider/transport/call/necessity inference.
Full SDL/SDPTool suites, vet and all 12 affected consumer tests pass. Actual XFMD:
14 interfaces, 12 consumers, 24 edges; all names retained in PDF, full-size reading
verified. Native Right-on-leaf activation is preserved and tested. Large-fit
packing/print readability remains open. VP05 mode-specific dependencies is next
proposed. The broader plan/cards remain active; no merge/release/install or
production model edit. Shared producer product files/history remain untouched.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230007--Architecture-viewpoint.md

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N7-001/README.md

## Session0007 T004 — VP05 candidate

N8 3b3d935 on sdp/xfmd-vp05-overview stacks on N7 85d62f2. VP05 now provides a
mode index plus isolated per-mode dependency diagrams/registers; all declared
capabilities and empty cases remain visible. No mode inheritance or combined
graph, provider or sufficiency inference. Original details and VP07 overlays
remain. Full producer suites/vet and 12 native consumer tests pass; native/PDF
mode isolation and actual XFMD's empty NormalDesktop output are verified.
VP06 activities/delivery planning is proposed next. Broader cards/plan stay active.
No merge/release/install, authored model/native product change or PG mutation.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230007--Architecture-viewpoint.md

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N8-001/README.md

## Session0007 T005 — VP06 candidate

N9 76fc87f on sdp/xfmd-vp06-overview stacks on N8 3b3d935. Complete
activity refinement and separate prerequisites retain isolated activities, with
source-linked status, responsibility, Feature and Scenario registers. Status is
source-claimed; missing addresses does not imply missing native implementation.
Full SDL/SDPTool suites, vet and 12 affected consumer tests pass. Native/PDF checks
verify both graphs and all statuses. Actual XFMD has no Activities: its explicit
empty plan and all 28 uncovered Functionality links are verified. VP08 channels/
scenario overview is proposed next; VP07 is already delivered. Broader plan/cards
remain active; N5–N9 await integration. No merge/release/install, authored model or
native product change, PG mutation, independent approval or large-print acceptance.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230007--Architecture-viewpoint.md

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N9-001/README.md

## Session0007 T006 — VP08 candidate

N10 d677d70 on sdp/xfmd-vp08-overview stacks on N9 76fc87f. VP08 has an
index and explicit mode targets, complete per-channel participation diagrams and
source-linked channel/binding/scenario registers. Original scenario sequences
remain on demand. The initial union failed Libavoid's routing budget; final
channel partitions retain all 60 unique nodes, 123 grouped edges and 202 uses
facts across 40 diagrams. All 22 scenarios and node names survive PDF; 850 source
spans are in bounds. Full producer suites/vet and 12 native consumer tests pass.

Actual model tree/select succeeds. Whole-project discovery hits the existing
20000-node/32 MiB guard on N9 and N10; KB025 retains that obligation. KB020 retains
the failed union renderer fixture. No generic large-print acceptance, independent
approval, native product/authored model change, merge/release/install or PG change.
Extended Session round complete; broader plan/cards remain active. N5–N10 await
integration. VP09 contract/storage overview is next proposed; VP07 already exists.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230007--Architecture-viewpoint.md

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N10-001/README.md

## Session0007 T007 — VP09 candidate

N11 6b80b40 on sdp/xfmd-vp09-overview stacks on d677d70. Complete data/storage/
contract inventories and Dataset neighborhoods retain source proofs and original
detail selections. Registers distinguish explicit persistence, source-claimed
completeness, optional fields and ternary projections without SQL/runtime inference.
Full producer suites/vet and 12 consumer tests pass. Actual XFMD's three Unit-held
Datasets render as 3 diagrams/9 nodes/6 edges; all 44 Contracts and 14 Fields survive
PDF, with 358 checked source spans. Persistent stores, Datagram/projection behavior
and unused shapes are verified by the native fixture. Contract count columns were
combined after PDF inspection; final artifacts were reverified.

Extended Session round complete; broader plan/cards stay active. N5–N11 await
integration. VP10 encodings/packets is next proposed. Existing discovery/routing
limits remain KB025/KB020. No merge/release/install, native product/authored model
change, PG mutation, independent approval or general large-print acceptance.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230007--Architecture-viewpoint.md

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N11-001/README.md

## Session0007 T008 — VP10 candidate

N12 f1c3342 on sdp/xfmd-vp10-overview stacks on 6b80b40. Complete Encoding/
Variant coverage and separate validated packet diagrams retain exact source links,
byte/bit order and ordered field positions/types/owners. No serialization or layout
is inferred. Original detail selection is unchanged. Full producer suites/vet and
12 consumer tests pass. Native fixture verifies two 29-bit packets and both orders;
34 source spans and all six placement statements check exactly. Actual XFMD has
no Encodings/Variants; its empty page is verified natively and in PDF.

Visual inspection reproduces clipped narrow-field labels and crowded adjacent bit
numbers in overview and original detail. KB020 retains this renderer follow-up;
new source registers remain complete. PDF text retention is not visual acceptance.
Initial native-test synchronization failure was corrected by waiting for the new
document/render token; final reruns pass. Evidence: XFMD N12-001.

Bounded projection round complete; broader plan/cards and renderer/discovery
residuals remain active. N5–N12 await integration. VP11 properties/traceability/gaps
is next proposed. No merge/release/install, native product/authored model change,
PG mutation, independent approval or general visual/large-print acceptance.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230007--Architecture-viewpoint.md

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N12-001/README.md

## Session0007 T009 — VP11 candidate and catalog handoff

N13 486b29d on sdp/xfmd-vp11-overview stacks on f1c3342. Complete properties,
snapshot counts, existing cross-viewpoint gaps and declaration/fact registries
retain exact source links and file/line context. Counts/properties are source
assertions, not implementation evidence; no new gap rules or inferred defaults.
Original VP11 selection remains unchanged. Full producer suites/vet and 12 consumer
tests pass, including composed-source freshness and native exact property selection.

Actual XFMD: 312 declarations, 1034 facts, 195 explicit properties, one reported
NO_ENCODINGS gap and 1346 exact source spans. Every fact sentence and declaration
name survives the 88-page PDF. Native fixture has 47 declarations/87 facts with
cross-viewpoint gap context. Inspected native Preview/source selection and sampled
PDF pages; no independent or general large-print acceptance. N13-001 retains inputs,
checks and artifacts. Source-link test restores Preview explicitly after Editor.

All catalog overviews now have selected development candidates. N5–N13 await
integration; N3j/N4 remains on main. Next: reconcile the stacked integration
candidate and prioritize retained KB020 acceptance/renderer work. KB025 discovery
guard and KB032 legacy eager-export alignment also remain open. Broader plan/cards
remain active. No merge/release/install, native product/authored model change or
PG mutation. Shared producer receives only this existing-card/Integration handoff.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230007--Architecture-viewpoint.md

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N13-001/README.md

## 2026-10-21 — N15 and local XFMD 0.9 delivery handoff

N15a producer e194721 on sdp/xfmd-lazy-discovery adds opt-in lazy-sdl discovery
and revision-bound model-only expansion, shared by CLI/JSON. XFMD Session0007
T012 implements the consumer and installs XFMD 0.9 plus the matching standalone
and private SDPTool development binary. Full final native suite passes 94/94,
including 47 GUI checks; actual XFMD discovery retains 63 sources and its native
model expansion exposes all eleven viewpoints and generates VP02 successfully.

This is local development installation, not main integration or a signed SDP
release. gh-sdp's signed 2.2.0 default and project process receipt remain unchanged.
N5–N13 plus N15a remain stacked for separate integration; PG concurrent work is
preserved. KB025 resource guards, KB020 acceptance and KB032 legacy eager-export
alignment remain explicit residuals. No shared producer product/history import.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N15b-001/README.md
