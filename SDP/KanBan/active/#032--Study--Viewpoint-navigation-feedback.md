# Evaluate SDL viewpoints through integrated XFMD navigation

| Field | Value |
| --- | --- |
| id | KB-SDP-032 |
| project | SDP |
| type | Study |
| CardState | in-progress |
| Systems | SDL, SDPTOOL |
| created | 2026-09-25T22:20:56Z |
| source | Owner conversation 2026-09-26: SDL context, implementation drift and post-main XFMD adoption |
| next_review | Review/integrate stacked VP02–VP06 and VP08–VP11 candidates; prioritize retained KB020/KB025 and legacy eager-export work. Release, shared-history reconciliation and large-print acceptance remain |

## Owner intent

Current generated viewpoints are an initial proposal requiring practical design
feedback. The owner is waiting for integrated SDP navigation in XFMD to browse
the model and assess what the views reveal or obscure. Generated output alone
does not establish a useful viewpoint catalog.

## Scope and outcome

Prepare a Study and refinement plan using real SDL source models. Inventory the
actual generators and navigation contracts before proposing new views. Evaluate
requirements/use cases, features/functionality, system/container structure,
channels/contracts and cross-cutting relationships across abstraction levels.
Capture navigation tasks, missing context, ambiguous notation, duplication and
useful drill-down/filter behavior. Include an assignment-authoring task so results
feed [KB-SDP-031](../completed/%23031--Study--SDL-assignment-bundles-and-blueprints.md).

Keep the SDL model as source of truth. Reproduce all diagrams via the tools;
record model revision, generation parameters and derived-versus-authored facts.
Do not hand-author purported generated diagrams or assume all viewpoints are
implemented. Test both overview and selected on-demand generation.

## Acceptance and ownership

A repeatable walkthrough records actual owner feedback, source-linked examples,
limitations and a prioritized, bounded refinement plan. Pending owner feedback
must remain pending rather than being inferred by an agent. Coordinate with
[KB-SDP-017](../active/%23017--Proposal--sdptool-and-project-navigation.md): that
card owns the navigation bridge; this card evaluates viewpoint usefulness.
XFMD owns its native sidebar implementation and its own KanBan work. No XFMD
application change or new broad renderer integration is authorized here.

## Worklog

2026-09-25T22:20:56Z: Registered in backlog before main integration; EVT-KB-SDP-000173.

## MVP1 pilot follow-up — 2026-09-27

The MVP1 pilot supplies separate operator/APT/simulator SDUI screens, registered for current SDPTool structural Markdown preview, plus a generated SVG review gallery. The first full SDL navigation acceptance path is BuckingUI → OperatorInteraction → InspectMachine → MachineService, then the separate simulator path. Full corpus parsing/linking remains a prerequisite; no owner GUI feedback has yet occurred.

[Pilot and navigation handoff](../../../experiments/mvp1_sdl/Navigation.md).

## Concrete owner feedback — XFMD Session 0003, 2026-10-06

The owner has now browsed generated views. This supersedes the earlier statement
that practical owner feedback has not yet occurred, without inventing acceptance.

- Remove A0–Ax presentation nodes from the default navigation; show each viewpoint
  once per model/System, preserving abstraction facts in the model and useful detail.
- Viewpoint-title activation should request a real aggregate overview, while tree
  expansion and focused child selection remain separate on-demand actions. Avoid
  generating every detail diagram as the substitute overview or rejecting at 24.
- Use-case overview includes all declared actors/use cases and functionality
  realization context. Use proper Mermaid use-case notation; preserve the distinction
  between use cases and internal functionality. The current official syntax is
  usecase-beta; compatibility of XFMD's native Rust path is separate work (#020).
- Focused use cases should include sequence diagrams when ordered scenario facts
  exist, plus authored or explicitly derived pre/postconditions with provenance.
  Static dependency edges cannot establish sequence order or missing conditions.
- VP07 should overlay features on the complete selected architecture; selecting
  one feature provides realization detail. Inventory overview/detail semantics,
  source facts and missing data for every remaining viewpoint.
- Hide non-actionable reference/depth-limit placeholder rows while preserving
  valid cross-reference targets and cycle protection. Genuine truncation/error
  diagnostics must remain visible through an appropriate diagnostic surface.
- Every source ID should link to original revision-aware exact source spans;
  carry declaration/relation provenance across composed files and SDUI normalization.
  Agree position units, multi-span IDs and stale-source behavior with XFMD #021.

The supplied VP07 ConfigurableScrolling Mermaid fails in XFMD with libavoid
non-orthogonal edge 11. It is captured under XFMD #020 with the exact source;
root cause is not established, and this is not evidence of invalid SDL semantics.
Do not rewrite model facts to bypass a rendering failure.

Producer responsibilities: SDL projection semantics and SDPTool tree/view/source
contract. Consumer: XFMD #028 (navigation/views), #021 (source spans), #020
(renderer), #029 (action groups). New KB-SDUI-005 owns SDUI combined previews;
this existing Study remains the SDL viewpoint handoff. No implementation selected.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230003--SDL-navigation-and-SDUI.md

## XFMD design handoff — Session 0003 T002

XFMD PLAN-XFMD-0009 separates consumer tree cleanup (N2, compatible with current
producer) from new overview projection (N3). Its selected native design preserves
original reference actions; canonical object resolution must not replace a
viewpoint-specific target with generic VP11. Phase/A-level wrappers become
transparent display groups without deleting model abstraction facts.

Concrete producer proposal for review: optional viewpoint-node overviewTarget,
using the existing typed target shape and an explicit producer-owned select URI
for one aggregate overview. Presence advertises actual implemented overview scope;
absence keeps older consumer fallback. Do not reinterpret current bulk selection
as the new overview or invent identifiers in XFMD. The proposal is not yet adopted
in Contract.md and no schema/producer implementation is claimed.

VP01/VP07 aggregates are new projections; existing goals()/feature logic emits
detail graphs. Reuse VP08 ordered scenario sequences only through explicit model
associations; missing associations/preconditions/postconditions remain gaps.
The XFMD design includes a complete VP01–VP11 overview/detail disposition matrix
for coordinated review. Authored-model hold remains intact in XFMD.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/04--Design/SDL-Navigation/Design.md

## N3d execution — XFMD Session0005 T004

Owner continues the diagram goal after native use-case notation delivery. Activate
the VP01 overview work unit under PLAN-XFMD-0009; broader viewpoint feedback stays
open. See [producer work unit](../../05--Implementation/SDL/ViewpointNavigation/Plan.md).
Use optional overviewTarget and explicit projection=overview; preserve original
detail targets. Current gh extension/release is unchanged. No fabricated scenario
order, pre/postconditions or System boundaries. VP07 and exact provenance follow.

## N3d outcome — 2026-10-08

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

Continue with VP07 and source/detail work under the linked producer work unit and external Session. No independent review, producer release or owner viewpoint acceptance claimed.

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

## N3i integration candidate — 2026-10-18

XFMD Session0005 T016 prepares its own earlier VP01/VP07 producer work as local
commit `cc7b52f84f05fdef63efc1a305023a12347026ad` on `sdp/xfmd-viewpoint-overviews`.
[Integration handoff](../../05--Implementation/SDL/ViewpointNavigation/Integration.md)
lists scope, tests and retained dirty-copy reconciliation. SDL and the full
SDPTool root package pass, including composed-source freshness/provenance; actual
XFMD model discovery/select and native overview/PDF checks pass. No independent
approval, main merge, release or installed gh change. This card remains in-progress
for the explicitly listed broader viewpoint obligations.

## N3j current-main compatibility — 2026-10-18

XFMD Session0005 T019 delivers successor `ecd9dcd5cd95fb2b1f323cea96da1776659b2513` on
`sdp/xfmd-viewpoint-overviews-main` directly above fetched main e8760d3. Use this
candidate for integration, replacing the earlier cc7b52f branch recommendation.
Both scoped conflicts are resolved additively; full SDL/SDPTool tests, vet, new
Actions discovery/overview regression and native XFMD/PDF verification pass.
Final merge forecast is clean. No main merge, push, release or card transition;
broader scope and shared historical ledger reconciliation remain open. Details:
[Integration](../../05--Implementation/SDL/ViewpointNavigation/Integration.md).

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
