# MVP1 design evolution and the proposed SDP skills

Status: **study updated with the owner's selected design direction; detailed contracts and migration remain to be reconciled**  
Prepared: 2026-09-15  
Purpose: explain the evolving MVP1 design, locate its latest reasoning, and show
how the candidate SDP skills can support its continued design and implementation.

## 1. Findings

MVP1 already has substantial architecture, constituent designs and two recent
discussion summaries. The immediate problem is not an absence of documentation.
It is the distance between **accepted design, later discussion, implementation
status and the entry points an agent is likely to read first**.

The latest summaries are:

- [Domain architecture discussion][domain-discussion]: domain/service ownership,
  a possible Go UI Host, APT and stem-profile responsibilities, and extraction
  of reusable mechanisms.
- [UI domain design discussion][ui-discussion]: Representation, Composition,
  Presentation, delivery and Renderer responsibilities, including corrections
  made during the discussion.

Both were committed on September 10 as working discussions. They capture
developments beyond the original Issue #72 body. After reading this study,
the owner selected that design as the direction to use, explicitly choosing
**Representation -> Composition -> Presentation -> Renderer**, and asked that
relevant earlier design be carried forward. Section 1.1 records this later
decision. The historical source files still describe their September 10 status;
their older status does not override the owner's subsequent instruction.

The candidate skills fit this situation well in principle: recover context,
analyse a change, make responsibility and compatibility decisions explicit,
maintain current SDP documents, and verify the resulting user workflow. Their
effectiveness in MVP1 has **not** been demonstrated by this study. Reading skill
files is not installing or activating them.

### 1.1 Owner decision: September 10 design is the target direction

Source: the owner's clarification in this conversation on 2026-09-15, after
publication of the original study on `sdp-vNow`. This is a direct owner
instruction, not an inferred disposition from the GitHub archive. No GitHub
comment URL exists for this clarification in the inspected sources.

The owner selects the September 10 design as the basis for further MVP1 design,
with Representation -> Composition -> Presentation -> Renderer as the explicit
architectural direction. Earlier design that the newer summaries do not cover
must be brought forward where it remains applicable.

Consequences:

- **Selected:** the separation of semantic Representation, instance/group
  Composition, declarative Presentation and replaceable Renderer. Future design
  work should make this structure concrete, not repeatedly ask whether to use it.
- **Superseded for the target design:** the earlier rule that an explicit common
  Renderer boundary is unnecessary for ordinary views. React can still implement
  a Renderer; the new boundary need not mean four processes or a wrapper class
  around every widget. This changes the architectural direction, not the status
  of already integrated code.
- **Carried forward:** existing system responsibilities, observable behaviour,
  contracts, visual intent and operational constraints wherever compatible.
  Omission from the newer discussion is not deletion. Section 5.1 identifies
  concrete inherited obligations and where they fit.
- **Still to resolve:** questions explicitly left open in the September 10
  design, including exact process/package boundaries, Go UI Host realization,
  APT/profile ownership splits, schemas and migration sequencing. These are
  design elaboration within the selected direction, not reasons to reopen the
  four-part separation.

The arrow expresses responsibility relationships, not a mandatory global FIFO
or a one-way-only runtime pipeline. Renderers send typed intents back through
the normal command path; subscriptions, snapshots and delivery scheduling have
their own explicit contracts.

This decision is recorded here now. The pinned MVP1 documents have not been
rewritten by this update. Their target Architecture/Design and current entry
points need a corresponding reconciliation so an agent entering only through
MVP1 can discover the decision. Choosing the design direction does not by itself
claim migration completion or resume the held Garden assignments.

## 2. Sources, chronology and limits

The inspected repository is `Hans-Einar/ponsse`. MVP1 was absent from the active
local Concept1 branch and from the fetched `origin/main`. Its sources were read
from an isolated archive of the integration branch without switching the user's
checkout or changing Ponsse files.

| Source | Inspected revision / state |
|---|---|
| `steering/mvp1-foundation`, PR #20 | `33e22b2427b281816d840077812cb55774ca203c`; PR still open |
| `steering/mvp1-wave1-integration` | `7589a5811a21ec6bb13979612060ceff0f690df8` |
| `origin/main` | `9cc5a691c60c977bdf4fea493a13437017013fa7` |
| SDP working baseline | `d611b8bf72aeb30d86c5ef28902469462b06a803`, with local candidate documents |
| Skills examined | [All 13 draft.2 candidates](../Toolkit/skills_v2/README.md), not an installed MVP1 contract |

The GitHub inventory covered all 72 issue/PR records available at inspection:
30 issues and 42 pull requests. Their issue-thread comments were retrieved and
screened for design relevance. The central architecture, review, Garden and UI
threads were followed in timestamp order, including dispositions after apparently
final reports. The directly relevant MVP1 PR review and inline-comment endpoints
contained no additional discussion; review reports instead appear in issue
threads and repository records. Section 11 gives the coverage map.

Chronology below uses UTC. Issue bodies are their current edited versions; the
API retrieval does not reconstruct every earlier body revision. Timestamped
comments and Git history establish identifiable changes. A comment labelled
MASTER, REVIEWER or STEERING is assessed by its content and disposition, not
assumed to be human approval from its account name. Linked private Cloud review
conversations and unrecorded owner conversations were not independently inspected.

The two September discussion documents report owner follow-up reasoning. That
reasoning is available as a repository summary, not as a complete independently
recoverable conversation transcript. Historical verification reports are cited
as reports; this study did not rerun product tests or inspect a running MVP1 UI.

## 3. How the design evolved

| Time | Development and evidence | Consequence for an agent now |
|---|---|---|
| Aug 29, before MVP1 foundation | [Issue #5 production direction](https://github.com/Hans-Einar/ponsse/issues/5#issuecomment-5463244356) separates serial transport, machine semantics, backend computation and browser presentation. Profiling had already challenged a premature language rewrite. | Recover the user/performance reason for a boundary, not just its current package name. |
| Aug 29, 16:49 | [Issue #19 UI boundary](https://github.com/Hans-Einar/ponsse/issues/19#issuecomment-5463645438) introduces narrow UI stores/selectors, stable rendering and separation from domain truth. | High-rate updates must not make the React root the whole-system state owner. |
| Aug 29, 18:28–20:26 | [Clean-start foundation](https://github.com/Hans-Einar/ponsse/issues/19#issuecomment-5464119725), [reuse disposition](https://github.com/Hans-Einar/ponsse/issues/19#issuecomment-5464357727), then [Refinement 01 disposition](https://github.com/Hans-Einar/ponsse/issues/19#issuecomment-5464705505). Concept1 becomes selective evidence/reuse material. Its visual language is retained; its runtime composition is reimplemented. | A clean start does not authorize discarding established operator experience. Ordinary React views already serve as renderers; an extra abstraction is optional for specialized graphics. |
| Aug 29, 22:02–22:43 | [Refinement 02][ref02] distinguishes direct semantic simulation from wire-faithful simulation, gives BWEB dual Machine/Bucking upstreams, and federates command catalogues. [Refinement 03][ref03] adds BuckingLab and clarifies domain versus diagnostic history ownership. | Do not force all browser information through BKS or turn Lab storage into mandatory domain event sourcing. |
| Aug 30 | [Constituent/Garden foundation](https://github.com/Hans-Einar/ponsse/issues/19#issuecomment-5465666095), [Proto v1 freeze](https://github.com/Hans-Einar/ponsse/issues/19#issuecomment-5467455742), [Refinement 04][ref04], and [Wave-1 design lock](https://github.com/Hans-Einar/ponsse/issues/19#issuecomment-5469485574). | There are local SDP owners and centrally governed contracts. Later scoped, concurrent server-granted authority supersedes earlier globally exclusive-writer wording. |
| Aug 30–Sep 1 | [Issue #24](https://github.com/Hans-Einar/ponsse/issues/24) starts NOT_READY and repeatedly reopens defects, including after an initial READY. Findings affect delivery semantics, requirement mapping, assignment boundaries and verification. [Final acceptance](https://github.com/Hans-Einar/ponsse/issues/24#issuecomment-5491049819) opens Point 5 at foundation `33e22b2`. | The first favourable review is not the final disposition. Wire conformance does not prove complete runtime behaviour. |
| Sep 1 | [Issue #30](https://github.com/Hans-Einar/ponsse/issues/30) starts bounded Garden work. A [shared manifest ownership collision](https://github.com/Hans-Einar/ponsse/issues/30#issuecomment-5494868395) produces an integration-stewardship disposition. | Local workers need a defined route for shared changes. Broadening their ownership incidentally is not integration design. |
| Sep 1, 15:16–20:58 | [Missing diagnostic-loss carrier](https://github.com/Hans-Einar/ponsse/issues/30#issuecomment-5496167210) leads to [Issue #44](https://github.com/Hans-Einar/ponsse/issues/44) and CTR-SLC-004. | A discovered consequence of an existing requirement can require a shared contract correction. Local passing tests cannot replace that correction. |
| Sep 2 | UIR's React adapter needs a [bounded dependency/lockfile companion](https://github.com/Hans-Einar/ponsse/issues/30#issuecomment-5503720006). Later, the [owner usage hold](https://github.com/Hans-Einar/ponsse/issues/30#issuecomment-5507290208) stops further dispatch, continuation and merges. | Accepted architecture, reviewed branch work, integrated work and permission to continue are separate facts. |
| Sep 2–9 | Concept1 [#59](https://github.com/Hans-Einar/ponsse/issues/59), [#69](https://github.com/Hans-Einar/ponsse/issues/69) and [#71](https://github.com/Hans-Einar/ponsse/issues/71) develop UI composition and expose acceptance/maintainability problems. | These are relevant lessons and references. Their Concept1-specific dispositions do not silently replace MVP1 authority. |
| Sep 9–10 | [Issue #72](https://github.com/Hans-Einar/ponsse/issues/72) proposes stronger Representation/Composition/Presentation/Renderer separation. The [owner clarification](https://github.com/Hans-Einar/ponsse/issues/72#issuecomment-5617593056) adds semantic inputs/outputs, capabilities, headless interaction and unhandled-event diagnostics. | This is a substantive candidate design change, not merely new names for UIR/BOX. |
| Sep 10, 18:20–18:23 | The two [repository discussion summaries](https://github.com/Hans-Einar/ponsse/issues/72#issuecomment-5623436090) become the primary place for further reasoning. They include Go UI Host, APT/profile ownership, scheduling and corrected StanForD assumptions. | Continue from these summaries and their open questions. Do not restart from the older issue body or treat strong discussion agreement as adopted architecture. |
| Sep 15, owner clarification after this study | The owner selects the September 10 design, explicitly Representation -> Composition -> Presentation -> Renderer, and requests preservation of relevant earlier design. See Section 1.1. | The direction is now selected. Elaborate its contracts and reconcile inherited obligations; do not treat the superseded optional-renderer rule as a veto. |

### Current lifecycle cannot be read from one status file

At the inspected integration revision, [MVP1 README][mvp-readme],
[ROADMAP_STATUS][roadmap] and [Traceability CurrentIndex][current-index] still
describe the seventeenth pre-Garden corrective candidate with Point 5 blocked.
The [precedence guide][precedence] directs readers to that roadmap. However,
the later explicit Issue #24 acceptance and Issue #30 activity establish that
Garden opened. [Garden CurrentWave][current-wave] records the subsequent
`WAVE1_GARDEN_OWNER_USAGE_HOLD`.

This is a concrete currentness defect in the entry points, not evidence that
all accepted technical content in those files is invalid. The constituent
registry still supplies ownership; its old lifecycle header supplies stale
status. A future reconciliation should make the current entry points agree
while retaining the historical corrective records. It must also inspect any
validators that encode the old status before changing files mechanically.

The three pending product PRs — [#58 UIR-SLC-004](https://github.com/Hans-Einar/ponsse/pull/58),
[#57 BOX-SLC-003](https://github.com/Hans-Einar/ponsse/pull/57), and
[#56 SIM-SLC-002](https://github.com/Hans-Einar/ponsse/pull/56) — were still open
at inspection. Garden records distinguish favourable reviews from the remaining
current-base synchronization and separate verification gates. September's design
discussion does not itself resume this implementation work.

## 4. The prior accepted baseline and inherited obligations

This section describes the inspected baseline before the Section 1.1 decision.
Preserve its applicable behaviour and responsibility rules in the target design;
implementation structure is subject to the owner's selected direction.

The [constituent registry][constituents] defines **17 governable constituents**.
These are not 17 C4 containers: contracts and shared/pure libraries are also
constituents, without necessarily being independently running applications.

| Responsibility | Accepted owner / boundary |
|---|---|
| Shared contracts | CTR owns schemas, compatibility, generated bindings and conformance; domain owners retain semantic authority. |
| Generic infrastructure | SVK: service communication; LBK: diagnostic history/replay mechanics; UIR: browser-local runtime; BOX: presentation models, views and visual primitives. |
| Pure forestry/format packages | TPR: taper models; BCK: bucking algorithms; SFC: StanForD Classic compatibility. BKS consumes these. |
| Simulator and real transport | SIM owns simulated truth and P1000 emulation. PGW owns physical serial transport. They expose the Target boundary to MCH; hidden simulator truth must not leak through normal machine paths. |
| Machine and bucking domains | MCH owns P1000 semantics, machine state and calibration. BKS owns profile/taper/bucking/APT/policy/economics/production. |
| Labs | MLB and BLB observe/support their domain services; they do not become production-state authorities. |
| Browser applications | BWEB is a thin Machine+Bucking adapter. BUI and SUI own client presentation and interaction, not machine/bucking decisions. |

Existing reusable communication boundaries include Target v1, Machine v1,
Bucking v1, simulator-control v1, browser runtime and Lab v1. Their actual
Protobuf/JSON schemas and transport standards should anchor further work.
The earlier [language study](Feature-Functionality-and-Channel-Study.md) does
not justify introducing a new API language or one channel per Feature here.
BOX/UIR interfaces are shared in-process contracts, not automatically network
channels.

The accepted [UIR architecture][uir-architecture] and [design][uir-design]
already separate typed pub/sub, immutable selector stores, command correlation,
session recovery and explicitly bounded presentation scheduling. Source
`event_seq` identifies events; per-subscription `delivery_seq` establishes
delivery continuity. Filtering and legitimate latest-wins presentation cannot
be treated as accidental durable-event loss.

The accepted [BOX architecture][box-architecture] and [design][box-design]
keep presentation models headless and map trusted model kinds to React views.
[Refinement 01][ref01] explicitly avoids requiring an extra renderer layer for
ordinary React boxes. The [SUI design][sui-design] supplies concrete behavioural
constraints: same-ID command recovery, no optimistic mutation of simulator
truth, distinguishable submission/outcome states, and high-rate updates isolated
from unrelated boxes. These remain relevant even if implementation technology
later changes.

Physical transmission and machine operation remain outside this study. A
simulator-only UI or contract test cannot establish machine suitability.

## 5. What the selected September 10 design adds

The September 10 documents supply the target design basis following Section 1.1.
The table separates that reasoning from details still needing design and
evidence. Compare with the prior baseline to preserve obligations and identify
concrete changes, rather than to reconsider the already selected direction.

| Concern | Latest reasoning | Decision or evidence still needed |
|---|---|---|
| Domain ownership | Explicit StemProfile and APT-domain responsibilities; format codecs below UI; domain-valid editable truth retained by a domain owner. | Separate services versus modules inside current owners; persistent library versus active/draft ownership; revised BKS/SFC boundaries. |
| UI runtime location | A possible Go UI Host feeds web/React, Fyne and potentially Flutter; native and web deployment need not use identical process boundaries. | Concrete consumers, deployment and latency requirements, compatibility with UIR/BOX, and the cost/benefit of moving runtime responsibility. |
| Representation | Generic semantic state/actions with identities, revisions, units, quality and input/output descriptors; no Ponsse-specific generic class hierarchy. | Minimal stable kinds, identity rules, command and patch contracts. Go embedding/interfaces remain a candidate technique. |
| Composition | Selects, groups and shares Representation instances and their lifetimes. One instance may appear in several views. | Exact boundary for semantic importance hints versus Presentation policy. Avoid duplicated subscriptions or copied authoritative state. |
| Presentation | Maps meaning into a view schema: layout, formatting, localization and abstract styling. Structural schema revisions differ from value patches. | Schema contract, validation, revision compatibility, draft/focus preservation and failure behaviour during replacement. |
| Delivery | Separates logical ordering/replay from physical scheduling; permits bounded lanes and independently paced bulk/realtime presentation. | Ordering scopes, bootstrap/watermarks, backpressure and durability classes. Illustrative update frequencies are not adopted performance requirements. |
| Renderer | Replaceable consumer of schema/state; emits typed user intent through the normal command path. Headless interaction is part of the proposal. | Required renderer set, transport versus in-process adapters, and real browser/native evidence alongside headless tests. |
| Collections | List, Table, Matrix and Series express different semantic structures. Cells can have lightweight capabilities without full independent objects. | Whether Table is primitive or a constrained List of records; collection patch/identity rules and measured large-data behaviour. |
| File standards | Normalize common APT meaning and preserve standard-specific extensions; missing information is not zero. Different standards may require different presentation. | Named target standards/versions, representative files and preservation/edit/export acceptance. |

Source: [system discussion][domain-discussion] and [UI discussion][ui-discussion].

Two corrections are particularly important to preserve. First, normalization
does **not** imply an identical UI across all StanForD versions. Second,
deterministic replay does **not** require one blocking physical FIFO for every
kind of work. Both corrections are already recorded in the summaries; an agent
should not revive the earlier interpretations through a simplified diagram.

The latest UI discussion also distinguishes unhandled event types, failed
consumers, intentional ignoring and representations with no visible view. Those
are different observations. A generic callback bus alone does not provide the
proposed diagnostic contract.

### 5.1 Earlier decisions to carry into the selected design

This is an initial reconciliation of the reviewed system and UI boundaries,
not a claim that every system requirement has already been mapped or verified.
The [normative system requirements][system-requirements] remain the checklist
for the complete mapping. Each material obligation must eventually have a
retained location, an explicit replacement, or a justified non-applicability
disposition; silence in a newer summary is insufficient.

| Earlier decision / reason | Placement in the selected design | Evidence to retain or establish |
|---|---|---|
| MCH owns machine/calibration semantics; BKS owns the current forestry/APT/economics/production domain; UI never owns those truths. [Registry][constituents] | Domain contracts feed the UI domain and Representations. Composition and Renderer do not acquire domain authority. Any intended APT/StemProfile split gets an explicit new owner rather than leaving responsibility between modules. | Ownership/import checks and domain-command tests; no calculation or validation authority moved into a view. |
| SIM hidden truth is distinct from measured machine state; direct semantic injection and wire-faithful simulation serve different tests. [Registry][constituents], [Refinement 02][ref02] | Preserve provenance in domain inputs and Representation quality/identity. Presentation can display different evidence sources without merging their meaning. | Test fixtures cannot accidentally use hidden truth as measured input; displayed source remains distinguishable. |
| BWEB has Machine and Bucking upstreams; Labs observe without taking production authority. [Refinements 02][ref02] and [03][ref03] | Reconcile these responsibilities explicitly when defining Go UI Host and domain boundaries. A new host does not silently turn every flow into a BKS dependency or a Lab ledger into production truth. | Consumer/owner map and integration cases for both upstreams and diagnostic sidecars. |
| Shared contracts are centrally governed and reusable across clients. [Registry][constituents] | Transport adapters preserve existing Target/Machine/Bucking/browser semantics. A Renderer protocol, if needed, has its own bounded responsibility; it does not replace all service contracts. | Generated-binding/conformance evidence plus explicit compatibility analysis for changed interfaces. |
| Concurrent authority is server-granted and scoped, with authoritative mutations serialized by the owning reducer. [Refinement 04][ref04] | Representation action availability is descriptive. Renderer intent follows normal UI command and server validation paths; Composition identity is not security identity. | Multiple clients, observer-only access and permission changes; UI `edit` capability cannot grant control authority. |
| An uncertain command retains its ID; timeout is not proof of cancellation. [UIR design][uir-design], [SUI design][sui-design] | Command lifecycle and correlation live outside view lifetime. Schema replacement or reconnect does not submit a new action or lose the initiating intent's result. | Same-ID retry/lookup, duplicate results, rejection and unknown outcome during presentation changes. |
| Sparse source sequence differs from contiguous delivery sequence; snapshot watermarks and session changes determine recovery. [UIR design][uir-design] | Transport/session processing establishes a coherent baseline before Representations consume changes. Preserve those semantics if the runtime moves to Go or crosses a new Renderer connection. | Filtering does not trigger false resync; real delivery gaps and changed sessions do; late consumers receive coherent state. |
| Durable events are not coalesced; high-rate replaceable presentation is bounded; diagnostic loss has explicit accounting. [UIR design][uir-design], [Issue #44](https://github.com/Hans-Einar/ponsse/issues/44) | Delivery mechanisms implement declared policies around Representation updates and Renderer delivery. Presentation pacing does not discard commands, production events or required history. | Independent bulk/realtime progress, bounded queues, preserved durable ordering, diagnostic notices and no false delivery gap. |
| Narrow subscriptions prevent unrelated updates from redrawing the application; plot entities keep stable identity. [UIR design][uir-design], [SUI design][sui-design] | Shared Representation identity and scoped Composition subscriptions feed incremental Renderer updates. Stable domain/representation/view identities remain distinct. | Two views share one value/subscription source; unrelated controls stay stable; value changes do not rebuild the whole schema or plot. |
| Concept1 provides the accepted operator visual reference, not its monolithic state architecture. [Refinement 01][ref01], [BOX architecture][box-architecture] | Presentation schemas and Renderer tokens retain density, hierarchy, layout intent and recognizable controls while the internal model changes. | Real viewport/workflow checks. Import any more detailed Concept1 geometry decision only with an explicit MVP1 applicability assessment. |
| Software Global Stop has distinct propagation, priority and identity; it is not a certified hardware E-stop. [Refinement 04][ref04] | Preserve the domain stop path and its command-specific authority rules. Presentation exposes the selected operator control; an ordinary action-disable rule must not accidentally remove stop availability. | Stop identity/deduplication, software automation cancellation and retained diagnostic recording; simulator evidence does not establish hardware suitability. |
| Generic mechanisms contain no Ponsse domain hierarchy; reuse must be demonstrated. [Registry][constituents], [UI discussion][ui-discussion] | Representation primitives, schema machinery and Renderer interfaces remain generic; domain-specific Composition/adapters stay with their owners. | Negative dependency checks and concrete consumers before extracting shared libraries. |

The explicit Renderer boundary is the identified structural supersession; the
old behavioural guarantees above are inputs to the new design. Additional
conflicts should be resolved individually, preserving their reasons and owner
intent instead of either copying the old structure wholesale or discarding it.

### Features, functionality and the several meanings of capability

MVP1's [component workflow][component-workflow] already defines a **local
Feature** as an independently valuable capability owned by one constituent.
Cross-system outcomes coordinate local Features through shared requirements
and contracts; an execution Slice remains bounded inside one constituent.

The owner's newer model uses Feature for a use-case-related system outcome and
Functionality for reusable behaviour within a unit. These can be mapped without
rewriting existing IDs:

| View | How to express it during this study |
|---|---|
| User/system outcome | A descriptive cross-system outcome linked to applicable system requirements and participating local Features. No new Feature registry is assumed. |
| Constituent contribution | Existing local Feature and requirement IDs, with their offered functionality described explicitly. |
| Implementation work | Existing constituent Slices and their integration/verification dependencies. |
| Communication | Existing shared contracts used by multiple outcomes. |

Likewise, a product capability, negotiated protocol capability, UI Representation
capability such as `edit`, and server-granted command authority are different
concepts. A visible editable field does not grant permission to perform an
authoritative command. Skills should demand a qualified meaning when ambiguity
affects design, without inventing a new taxonomy for every occurrence.

## 6. How each candidate skill fits MVP1

The skills supply an operating method. MVP1 documents supply project knowledge
and authority. Encoding the entire changing architecture inside skill text
would create another stale copy of the system.

| Skill | MVP1-specific use and expected result |
|---|---|
| [sdp](../Toolkit/skills_v2/sdp/SKILL.md) | Resolve the actual MVP1 branch, system and constituent SDP roots; include the Section 1.1 owner decision alongside repository dispositions; distinguish selected direction, unresolved contracts and held implementation; route only relevant roles. |
| [sdp-change-analysis](../Toolkit/skills_v2/sdp-change-analysis/SKILL.md) | Translate a short observation into an affected workflow and inherited constraints. For a UI symptom, follow state ownership through transport, UI domain, presentation and rendering before choosing a layer to edit. |
| [sdp-architect](../Toolkit/skills_v2/sdp-architect/SKILL.md) | Elaborate the selected four-part UI model and reconcile UIR/BOX obligations using Section 5.1. Record what is retained, changed and why; distinguish services from modules; resolve Go Host and schema details without reopening the selected direction. |
| [sdp-master](../Toolkit/skills_v2/sdp-master/SKILL.md) | Turn an adopted design into bounded constituent work and coordinated integration. Preserve MVP1's granular assignments, shared ownership rules and current hold. A cross-system outcome does not make a worker's scope unbounded. |
| [sdp-worker](../Toolkit/skills_v2/sdp-worker/SKILL.md) | Implement the assigned contribution, preserve inherited invariants, and return out-of-scope discoveries with evidence and proposed document targets. Do not bury an architectural choice in a new helper or dependency. |
| [sdp-reviewer](../Toolkit/skills_v2/sdp-reviewer/SKILL.md) | Challenge whether the assignment solves the intended workflow before evaluating its implementation. For UI, inspect ownership, schema/value separation, command recovery and actual interaction consequences. Review a design proposal without treating it as approved implementation. |
| [sdp-verifier](../Toolkit/skills_v2/sdp-verifier/SKILL.md) | Match evidence to claims: pure model, contract, provider/consumer, real integration or user workflow. Headless schema tests cannot prove layout, focus, accessibility or smooth browser behaviour. Prior branch results do not establish a newly combined candidate. |
| [sdp-traceability](../Toolkit/skills_v2/sdp-traceability/SKILL.md) | Link the reason, disposition, current requirement/design, work and evidence. Reconcile supersession prospectively; do not make an issue comment or ledger entry the only current explanation. |
| [sdp-steering](../Toolkit/skills_v2/sdp-steering/SKILL.md) | Present unresolved owner choices in outcome language, with a recommendation and consequences. Record actual dispositions, including a hold. Do not request confirmation for routine choices already covered by authority. |
| [sdp-auditor](../Toolkit/skills_v2/sdp-auditor/SKILL.md) | Detect stale roadmap/index pointers, contradictory status, missing discovery/installation and misleading completion claims. Report the narrow inconsistency rather than invalidating the entire historical baseline. |
| [sdp-vertical-refactor](../Toolkit/skills_v2/sdp-vertical-refactor/SKILL.md) | If migration is adopted, move one complete interaction through the new boundaries with explicit adapters and removal conditions. Separate preserved behaviour from any separately approved behavioural redesign. |
| [sdp-versioning](../Toolkit/skills_v2/sdp-versioning/SKILL.md) | Evaluate compatibility when schemas, shared interfaces or published packages change. Folder names, a Go implementation or the label `skills_v2` do not determine a release version. |
| [sdp-release](../Toolkit/skills_v2/sdp-release/SKILL.md) | Apply only when an actual release is in scope. Check the selected integrated contents and publication authority; a design document or reviewed constituent is not a released MVP1 system. |

Not every row requires a separate agent, a new document or an approval round.
Use the existing role contract and a proportionate handoff. Independent review
must be independent where required; naming a second section “Reviewer” does not
create that independence.

## 7. A concrete requirement discovered during implementation

Issue #44 is an unusually useful example of the process we want to preserve.

1. **Existing requirement:** diagnostic loss must be visible to affected
   clients, not just counted privately inside a service queue.
2. **Implementation discovery:** SVK-SLC-003 could count loss locally, but the
   frozen CTR carrier had no field through which the remote client could learn it.
3. **Rejected shortcut:** weaken the requirement to queue-local accounting.
4. **Scoped disposition:** authorize CTR-SLC-004 to amend the shared contract,
   generate bindings, update fixtures and prove Go/TypeScript conformance.
5. **Concrete semantics:** `diagnostic_drops_before` reports intentionally
   dropped diagnostic deliveries before the next delivery. It does not count
   filtering, legitimate latest-wins replacement, or create a delivery-sequence
   gap requiring resynchronization.
6. **Integration:** accept the correction, then resume the dependent SVK work.

Evidence: [discovery](https://github.com/Hans-Einar/ponsse/issues/30#issuecomment-5496167210),
[disposition and contract scope](https://github.com/Hans-Einar/ponsse/issues/44),
[integration](https://github.com/Hans-Einar/ponsse/issues/30#issuecomment-5500339461).

This is not a new user Feature simply because implementation uncovered it.
It is a missing mechanism required to satisfy an existing observable obligation.
Worker identifies it; Architect/CTR resolve the contract; Master coordinates
affected owners; Reviewer/Verifier check the consequence; Traceability keeps the
current requirement and design connected to the reason. The shared
[document workflow](../Toolkit/skills_v2/sdp/references/document-workflow.md)
already makes this distinction. MVP1 gives us a real evaluation case for it.

## 8. A practical design workflow for the next UI discussion

Consider a short owner prompt:

> Show the same measured length in two places, and let me change the layout
> without interrupting measurement updates or losing an edit.

This example is a proposed design exercise, not an active MVP1 assignment.

**Recover the relevant contract.** The agent locates the accepted UIR/BOX and
application designs, the latest discussion sections on identity/schema changes,
and the current implementation state. It describes the observable outcome:
two views share one semantic value; layout replacement preserves continuity;
typed-but-unsubmitted text has an explicit preservation rule.

**Separate obligations from choices.** Domain truth and command authority remain
with their owners. A schema change cannot silently duplicate actions or turn
unknown state into a valid zero. However, whether an unfinished edit survives
every structural change, or asks the operator to finish/cancel, is a behaviour
choice if existing authority does not settle it. The agent asks only about such
material ambiguity, using a concrete example and recommendation.

**Compare realizations within the selected direction.** Decide how existing
BoxModel state maps into Representation and Presentation, how Composition owns
shared instances, and how a React Renderer consumes their contracts. Compare
in-process versus network adapters and migration increments where those remain
open. Keeping the old structure unchanged is no longer a competing target
architecture merely because it could solve this one example. Existing code may
still be reused or temporarily adapted where it satisfies the new boundaries.

**Make the design reviewable.** Define identity, lifetime, schema and value
revisions, invalid-schema recovery, command correlation and ownership of local
text drafts. Link retained requirements and identify prospective supersession.
Use a small state/interaction example; do not first generate a universal framework.

**Plan evidence at the correct levels.** If implementation is later authorized,
verify shared identity/subscription behaviour, schema rejection retaining the
last valid view, command recovery and independent high-rate/bulk scheduling.
Then exercise the real UI: two visible values, schema replacement during updates,
focus/draft behaviour, reachability at the required viewport and no unwanted
whole-page redraw. Any timing budget must come from the chosen acceptance
contract, not illustrative discussion numbers.

**Maintain the project record by consequence.** Update the relevant local
Requirements, Architecture and Design where their meaning changes. A changed
domain owner also affects the system registry/architecture; a shared wire change
also affects CTR and consumers. Record actual review and verification against
the candidate. The report and ledger point to these current explanations rather
than replacing them.

Current authorization covers the study and the owner's selection of the design
direction. The historical Garden hold does not prevent elaborating that design;
the example does not itself release implementation workers or resume pending PRs.

## 9. What should be improved before adopting the skills in MVP1

The current candidates already contain most of the necessary principles. MVP1
exposes four focused improvements to evaluate, rather than a reason to add more
roles or copy the whole architecture into a skill.

| Improvement | Practical implementation to evaluate |
|---|---|
| Reliable current-context entry | Reconcile MVP1's existing lifecycle entry points and add direct links to current discussion versus accepted designs. Teach/test the router to follow later explicit dispositions when an entry point is demonstrably stale. Avoid a second competing status registry. |
| Explicit vocabulary bridge | Add a brief project glossary mapping constituent, runtime container, local Feature, cross-system outcome and the relevant meanings of capability. Resolve project vocabulary before applying generic skill terminology. |
| Reviewable design evolution | For each proposed change, keep a concise delta: current decision and reason; proposed replacement and reason; affected owners/contracts; adoption status; required evidence. Store it in the existing discussion/design records and reconcile accepted documents when adopted. |
| Proportionate behaviour tests for skills | Test whether an agent actually finds the relevant decision and updates the right durable document. Do not judge success by frontmatter validity, report length or ledger volume. |

The first two are primarily project-document work. Generic skill guidance should
remain short and refer to local mappings. Reviewer guidance should make
design-only review explicit enough that an agent does not demand implementation
evidence to discuss a proposal, or claim implementation readiness from prose.
These are recommendations from this study; the skills have not been revised or
installed as part of it.

Useful evaluation scenarios are: the stale pre-Garden entry point; the #44
missing carrier; the older optional-renderer decision followed by the owner's
explicit selection of the new Renderer boundary; a held PR with favourable earlier reviews; and the owner's
short UI prompt above. Supply raw source evidence to a fresh context without
the expected answer. Observe authority recovery, scope, questions, document
maintenance and the specificity of verification claims. No observed trial in
this study establishes that these skills already work in MVP1.

## 10. Recommended next design step

Reconcile the owner decision into MVP1's target Architecture/Design and the two
existing discussion documents. Use Section 5.1 to carry forward earlier
obligations. First work one concrete interaction through the selected ownership:
shared live measurement plus
layout replacement is small enough to expose Representation identity, scheduling,
Renderer commands and draft continuity. A separate APT example can then test
format-specific meaning and domain ownership using explicitly selected standards
and representative files.

Produce a small baseline-to-target delta and evidence plan that resolves the
remaining Go Host, service-boundary and migration questions. The four-part UI
direction is already chosen. This allows the owner to describe the desired experience briefly while
the agent performs the architectural translation and exposes the few choices
that actually change the product.

The missing link is a maintained transition from **reasoning to current design
to verified behaviour**. MVP1 already demonstrates parts of that transition,
especially the shared-contract correction. The skill adoption should make it
repeatable and easy to discover, without reproducing every historical discussion
on every assignment.

## 11. Source coverage and verification of this study

| Source group | Treatment |
|---|---|
| Direct MVP1 discussion: #19, #24, #30, #44, #72 | Primary chronology, accepted decisions, corrective findings, hold and newer proposals; comments inspected in timestamp order. |
| Direct MVP1 PRs: #20, #31–33, #35–38, #40–43, #46, #50–58 | Foundation and constituent implementation context; issue bodies and discussion endpoints screened. Current PR state checked specifically for #20, #31 and pending #56–58. This is not a fresh code review of every PR. |
| Earlier Concept1 design: #1, #2, #3, #5 and related integration PRs | Taper/measurement evidence, economic/APT meaning and the performance/ownership origins of MVP1. Their application to MVP1 is mediated by its clean-start/reuse disposition. |
| Later Concept1 UI/process: #21, #59, #62, #65–71 | Related presentation, geometry, lifecycle and acceptance lessons. #59, #69 and #71 are especially relevant to the later UI discussion. Project-specific implementation decisions are not promoted to MVP1 authority. |
| Protocol, firmware and emulator threads, remaining records | Screened for system-design links; retained as domain evidence or separate work where relevant. No hardware finding or independent emulator disposition is adopted as an MVP1 UI decision. |
| Repository documents | Accepted constituent/workflow/contract authority, UIR/BOX/SUI design, refinements, precedence/status, Garden records and the two latest discussion summaries compared at the pinned integration revision. |
| Subsequent owner instruction, 2026-09-15 | Direct conversation decision recorded in Section 1.1; this updates the study's recommendation and target direction without pretending that the pinned September 10 files already contained that decision. |

Only this study and its SDP README index link are added for the current request.
Ponsse source, project authority, pending PRs and skill files are unchanged.
Verification is documentary: source/status comparison and local Markdown link
checks. Product builds, browser behaviour, native skill discovery and deployment
remain outside the evidence established here.

[domain-discussion]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/Domains/DetailedArchitectureDiscussion.md
[ui-discussion]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/Domains/UIDomain/DetailedDesignDiscussion.md
[constituents]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/CONSTITUENTS.md
[component-workflow]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/COMPONENT_SDP_WORKFLOW.md
[precedence]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/doc/MVP1_AUTHORITY_PRECEDENCE.md
[mvp-readme]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/README.md
[roadmap]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/ROADMAP_STATUS.md
[current-index]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/SDP/Traceability/CurrentIndex.yaml
[current-wave]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/SDP/Garden/CurrentWave.yaml
[ref01]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/doc/MVP1_ARCHITECTURE_REFINEMENT_01.md
[ref02]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/doc/MVP1_ARCHITECTURE_REFINEMENT_02.md
[ref03]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/doc/MVP1_ARCHITECTURE_REFINEMENT_03.md
[ref04]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/doc/MVP1_ARCHITECTURE_REFINEMENT_04.md
[uir-architecture]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/shared/ui-runtime/SDP/04--Architecture/architecture.md
[uir-design]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/shared/ui-runtime/SDP/05--Design/design.md
[box-architecture]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/shared/box-ui/SDP/04--Architecture/architecture.md
[box-design]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/shared/box-ui/SDP/05--Design/design.md
[sui-design]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/web/simulator-ui/SDP/05--Design/design.md
[system-requirements]: https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/SDP/03--Requirements/requirements.md
