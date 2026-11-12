# PLAN-SDP-0025 — Producer-owned SDL browsing design and implementation handoff

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0025 |
| project | SDP |
| state | planned |
| PlanType | DesignPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDL, SDPTOOL; collaborating XFMD |
| source | Owner planning request, 2026-11-11; KB-SDP-060; HEOS Session0012 |

## Outcome and authority

Prepare an implementable design for readable, producer-owned SDL browsing across
VP01–VP11. The owner authorizes planning now. This plan is prepared, not executed;
requirements/model changes, product implementation and publication remain unselected.

This is a functional navigation improvement with architectural refactoring inside
it. It changes observable behavior and the producer/consumer contract, so a pure
behavior-preserving Refactor label is insufficient. The adopted plan contract has
no RefactorPlan type. Use DesignPlan now, followed by an ImplementationPlan whose
milestones are refined against the resulting model and Blueprint. A separate
Refactor record is unnecessary unless later work needs an independently bounded
behavior-preserving migration.

## Need, source and scope

Need: a reader should quickly understand system structure, responsibilities,
collaboration and inner mechanisms without reading implementation code. Existing
diagram inventories, duplicated tree routes and exhaustive proof tables obscure
that information. This need motivates testable requirements; it is not itself
an implementation specification.

Authoritative consumer analysis remains in HEOS:
[Study](../../../../../HEOS/SDP/03--Architecture/SDLBrowsing/Study.md) and
[NavigationDesign](../../../../../HEOS/SDP/03--Architecture/SDLBrowsing/NavigationDesign.md).
Do not maintain a competing copy here. These are locally authored, uncommitted
planning inputs at this checkpoint; capture exact revisions when producing the
governed model/Blueprint handoff.

Cover Systems/include membership, local fragment browsing, numeric viewpoint
order, consistent collection/object actions, meaningful overview/detail pages,
register separation, readable source links, truthful availability, channel endpoint
notation and generic XFMD consumption. Preserve source/model/mode identity,
freshness, on-demand generation and access to unique existing information.
No new SDL port entities and no hand-maintained generated viewpoints.

## Existing work and boundaries

- [KB060](../../../KanBan/active/%23060--Change--Producer-owned-SDL-browsing.md)
  owns this new bounded change. [KB032](../../../KanBan/active/%23032--Study--Viewpoint-navigation-feedback.md)
  remains predecessor/context, retaining its unresolved historical work.
- [Existing producer work](../../../05--Implementation/SDL/ViewpointNavigation/Plan.md)
  and [PLAN-XFMD-0009](../../../../../xfmd/SDP/05--Implementation/SDL-Navigation/Plan.md)
  retain their completed evidence and open obligations. Reconcile exact overlap
  before selecting new consumer implementation; do not restart old milestones.
- [Blueprint implementation](../../../05--Implementation/SDPTool/Blueprints/Plan.md)
  is marked completed for its bounded producer scope. Blueprint output is not
  evidence that the source model fully captures code or that PG execution is ready.
- [PG Session/model handoff](../ProjectGovernance/Session-Model-Handoff.md)
  governs exact WORK/task/Blueprint identity, design disposition and model release.
- Proposed management Activity work is [KB061](../../../KanBan/onHold/%23061--Proposal--Need-driven-activities.md).
  Owner decision 2026-11-11 now makes qualified Activity integration a predecessor
  of browsing work. This plan remains planned and KB060 queued until that outcome.

## Design milestones

| Phase | Milestone | Observable acceptance | State |
| --- | --- | --- | --- |
| Requirements and baseline | D1 Reading tasks and scope | Map NavigationDesign to existing/new requirements and reader acceptance tasks; inventory old targets/content and exact current producer/consumer baseline; distinguish bugs from new behavior | planned |
| Architecture and WORK design | D2 Model and protocol | Reconcile SDL/SDPTool/XFMD ownership; update applicable WORK models and prose contracts for membership, ordered tree, generic activation, collection/register scope, compatibility and failure handling | planned |
| Presentation design | D3 Reading and rendering | Define all eleven viewpoint dispositions, category membership and shared page rules; resolve representative VP01/VP04/VP08/fragment cases plus channel attachment feasibility in the native renderer | planned |
| Handoff | D4 Blueprint and delivery plan | Prepare reviewable NOW/TARGET context and authored assignment, expose UNKNOWN code mappings, obtain applicable design disposition, and create the bounded ImplementationPlan with testable milestones/cards | planned |

D3 examples validate the shared rules; they do not limit the implementation to
three viewpoints. Readability must be checked in actual XFMD during implementation,
not inferred from successful source parsing. No checks above are claimed passed.

## Proposed successive Sessions

Session routes select work; typed plans own detailed milestones. Do not copy
milestones into an additional authoritative SessionPlan. These are route slots,
not preallocated Session records. Create each next Session when that work starts.

| Route | Purpose and exit | Session / plan reference | Status |
| --- | --- | --- | --- |
| R0 | Diagnose browsing and establish reading intent | HEOS Session0012 and NavigationDesign | study delivered; this plan receives follow-up |
| R1 | Prepare route, initial plan and Activity proposal | [SDP Session0015](../../../Sessions/session-%230015--SDL_browsing_planning.md) | planning deliverable complete |
| R2 | Requirements, architecture and SDL WORK design; resolve presentation/protocol decisions | Future Session, created on selection; this plan D1–D3 | planned |
| R3 | Reviewed Blueprint context and detailed implementation planning | Future Session; D4, then link the actual ImplementationPlan | planned |
| R4 | Deliver one or more coherent increments from the ImplementationPlan | Successive execution Sessions referencing its selected milestones | proposed; no records allocated |
| R5 | Integrated verification, independent review and explicit disposition | Final execution/verification Session, or part of the last delivery Session | proposed |

R2/R3 may be combined if a bounded design is ready; R4 may need several Sessions.
Session boundaries reflect meaningful work/handoffs, not chat turns or compaction.
No automatic acceptance/release follows from a completed Session.

## Preliminary implementation themes, to refine at D4

Do not treat each horizontal subsystem area as a delivery milestone. Use complete
reader workflows and split further only when outcomes can be independently checked:

1. Select a System or included source in the correct context; browse ordered
   viewpoints through generic producer-advertised actions in XFMD.
2. Open a collection, read its chapters, open an object and return with context
   intact. Replace duplicate routes and reconcile content across all VP01–VP11.
3. Read useful overview/detail pages with source-statement links, honest availability
   and separately selected complete registers, including isolated declarations.
4. Follow container/unit/channel collaboration and inspect payload/contract details
   using verified native notation, large-model handling and print output.
5. Verify migration/content coverage, old producer/client compatibility, staleness,
   cancellation, history and whole-workflow readability; consolidate independent review.

These themes guide design, not a second executable implementation plan. New
technical KB cards should describe discovered independent obligations, not mirror
every file or duplicate each plan row.

## Git policy and preflight

Planning branch: sdp/sdl-browsing-planning from a96affd5378b1cee1a2926a3808bb9e2f2dccbb9.
Use milestone commits for coherent selected work. No commit/push/merge/release
performed by this planning assignment; publication authority is separately scoped.

Planning inspection found one registered SDP worktree, clean before branching.
Existing PG/Blueprint/navigation/release records retain their states. XFMD main
8893a32 has unrelated RPM/process changes; HEOS main has this conversation's
uncommitted study records. Leave both sets intact. Scoped verdict: documentation
planning only is clear; product-write readiness is not established. Rerun the
full concurrent-work preflight before D2 model edits or implementation, including
existing model-source holds and current native work. Known local records are not
proof of absent remote/concurrent work.

## Completion and verification

This DesignPlan completes only when D1–D4 have actual evidence and a concrete
implementation handoff. Preparing this document completes neither the plan nor
the product change. The current checkpoint verifies links, identities and ledger
registration only. Retain existing framework/schema versions; Activity adoption
is a separate proposal and does not silently change this plan's lifecycle.

## Owner sequence refinement — 2026-11-12

Activity0001 restructures SDPTool without new functionality. Activity support will
be the next manually created Activity. Browsing follows as the third and first
tool-created Activity. This plan remains planned and KB060 queued; no browsing
execution starts during restructuring.

## Selected execution

Owner now authorizes Activity3 implementation and public XFMD0.10 RPM delivery.
B0 of PLAN-SDP-0028 consolidates this design route; the Activity3 ImplementationContract
and ReadingWorkflow WORK/Blueprint pin the selected boundary. Historical planning
authority above is superseded for this selected scope; no main merge is authorized.
