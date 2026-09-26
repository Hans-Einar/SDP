# MVP1 SDL exercise: findings and local extensions

Date: 2026-09-17

Status: exercise-local proposal, recorded before authoring the expanded MVP1
model. The owner authorizes necessary extensions for this exercise, with findings
kept here before promotion into the general language. This document does not
modify the general language definition or the `design-core 0.1` parser.

## 1. Objective and evidence boundary

Describe the complete known MVP1 responsibility/contract topology, all registered
constituents, selected UI architecture, system/local obligations and major behavior
paths. Distinguish an inventory-complete design draft from an executable or
implementation-complete system. Unknowns must be named objects, not silently
filled in to make the model appear complete.

The source is Ponsse integration revision
`7589a5811a21ec6bb13979612060ceff0f690df8`; the remote branch head was checked on
2026-09-17 and still matches. Issue #72 comments were refreshed and point to the
September 10 discussions. The owner's subsequent selection of
Representation -> Composition -> Presentation -> Renderer takes precedence over
their historical discussion-only headers. The active local Concept1 branch lacks
MVP1 and contains unrelated changes, so this exercise lives in
`experiments/mvp1_sdl/SDL/MVP1/` in the SDP repository. It can later be placed at
`ponsse/MVP1/SDL/MVP1/` in the appropriate checkout.

Source classes: the constituent registry; all 99 system Requirements; all 17 local
Requirements documents; shared Protobuf/browser contracts; both September 10
discussions; existing UIR/BOX/SUI designs; and the previous chronological study.
Repository references use the pinned revision. Owner decisions made in this
conversation are identified separately from repository evidence.

## 2. Exercise language boundary

Each file identifies `language sdl-mvp1-exercise version 0.1.`. This is a temporary
candidate surface for this corpus, not a new supported dialect or a published
general language release. It combines agreed structural concepts with the local
extensions below. No file claims conformance to the current parser.

One System is declared in the explicitly selected entry `System.design`.
`includes "relative/path.design".` selects source files relative to that entry.
All declarations across included files are collected before resolving references.
For this exercise there is one system-global namespace: file moves preserve
identity; names are unique. Package-qualified lookup remains a separate future
extension rather than an implicit filesystem convention.

Identifiers retain `[A-Z][A-Za-z0-9]*`; type and relation keywords are lowercase.
Each sentence ends with a period. Quoted strings use JSON string escaping. Full
lines beginning `//` are nonsemantic comments; executable meaning cannot be hidden
inside comments or quoted prose. Annotation rules from the latest owner decision
apply to typed references; an annotation is checked against the declaration, not
used to create an object.

## 3. Findings and proposed extensions

| ID | Pressure found in MVP1 | Exercise extension | Limit / promotion question |
|---|---|---|---|
| E01 | Registry mixes nine runtime applications and eight shared/pure constituents. | `library` and `subsystem` alongside Container/Unit. | Library use must instantiate state per application; never imply a shared runtime singleton. |
| E02 | Horizontal UI concerns cross package ownership. | `layer`; `has-layer`, `groups`; ordinary `contains` remains accountable structural containment. | Layer membership is not another owner or deployment. |
| E03 | Simulator and physical gateway implement one Target contract. | Mode-qualified participation; `external-device`, `client-tool`, `port`, `adapter`. | Model membership does not enable physical output; selection must exclude accidental dual-target delivery. |
| E04 | Known wire schemas should not be rewritten into a second authority. | `contract`, `message-set`, `follows`, `allows`, `schema`, `member`, `through`. | Schema links are exact source references; new UI/APT contracts remain explicitly incomplete. |
| E05 | Sender/receiver alone cannot express independent subscribers or addressed results. | Named contract policies with `specifies`; per-message-set delivery classification. | Policy prose is an obligation, not an executable queue configuration. Numeric limits remain gaps where unsettled. |
| E06 | Needs and offers cross boundaries without transferring ownership. | Provided Capabilities, library `depends-on`, Interface `exposes`, mode-scoped `requires`. | Build dependency, logical participation and runtime transport are different facts. |
| E07 | UI's chosen direction supersedes one existing system and one BOX rule. | `decision`, `requirement`, status/disposition, `constrains`, `replaced-by`, source citations. | Mapping records intended coverage, not proof of implementation or acceptance. |
| E08 | Source-backed obligations and unbound target behavior coexist. | `gap`, `binding`, `evidence`; `affects`, `anchors`, status. | Unbound is not implemented; a source anchor is not a verification claim. |
| E09 | Many-to-many Feature realization and whole-system paths are absent from core. | `feature`, `scenario`, `illustrates`, `exercises`, `step`, `calls`. | Feature does not become an owning unit or require a dedicated channel. |
| E10 | Invalid/stale/concurrent paths need explicit control structure. | `condition`, `outcome`, `starts-at`, `proceeds-to ... when`, `finishes ... with`, `runs-alongside`, `preserves`. | Named conditions have recorded meanings but no executable predicate calculus yet. |
| E11 | Identity, domain revision and UI revision must remain distinct. | `data-type`, `state`, `field`, `typed-as`, `initially`, `tracks`. | Field schemas are conceptual where no authoritative schema exists; no invented units or reset values. |
| E12 | Two views of one value differ from two Representation instances. | `representation`, `composition`, `presentation`, `view`, `references`, `renders`, `binds-to`. | These are design specifications; runtime instance multiplicity/lifetime remains explicit contract work. |
| E13 | Same runtime design can be embedded or hosted separately. | Preserve logical UI subsystem and app contexts; record hosting choice as Gap. | Do not assert an accepted Go host Container or invent separate Stem/APT services. |
| E14 | Unknown quality, unsupported formats and unbound algorithms must stay visible. | Explicit gap and availability/state obligations; retained versus proposed versus superseded status. | Prose obligations cannot make interpreter completeness pass. |
| E15 | A flow must connect actual channel traffic to its responsible handler. | Step `receives`/`emits` a MessageSet `through` Channel in its declared call context. | Sender/receiver membership and contract routing must agree; this is a designed interaction, not observed execution. |
| E16 | State and DataType do not identify held data, logical origins or contract-defined message families. | Contract-bound Datasets; optional Database ownership; Dataset-linked Datagram families whose contracts define variants; producer refinement. | [Revised data contracts proposal](SDL-Datasets-Datagrams-and-Data-Contracts.md) records source identity, presence and snapshots. MessageSet is now derived under E17; this proposal is not part of the baseline corpus/profile. |
| E17 | Commands/Values and inter-layer interactions need explicit contracts without another authored message inventory. | Candidate ControlSet/Command contracts; generated MessageSet; one logical Channel concept with declared layer/Container boundary scope. | [ControlSet and boundary study](SDL-ControlSets-Layer-Boundaries-and-Data-Access.md) separates owner-selected generated catalogs from proposed Value membership and Channel naming; arguments/results and execution bindings remain explicit work. |
| E18 | On-demand access, unsolicited arrival, reset durability and current availability are easily conflated. | Database as declared addressable retained data; operation-specific query/store capabilities; live-only Dataset and explicit cache/replay distinctions. | [Data-access study](SDL-ControlSets-Layer-Boundaries-and-Data-Access.md) leaves exact persistence guarantees open and rejects inferred full CRUD or invented P1000 protocol support. |
| E19 | Values risk being duplicated between data, control and UI surfaces. | Candidate owned Value with optional Datagram binding, internal discovery and Representation observation; command-only ControlSet under evaluation. | [Checkpoint data model](../../../SDP/History/checkpoint-1/03-Values-Data-and-Control.md) keeps the owner's latest exploration provisional; neither Value nor its binding grammar is implemented. |
| E20 | Detailed Container work can obscure functional intent and cross-level realization. | Candidate abstraction profiles, explicit realization/allocation links and Functionality before deployment allocation. | [Checkpoint model](../../../SDP/History/checkpoint-1/02-SDL-Model-and-Abstraction-Levels.md) identifies the conflict with current unit-local Functionality; a selected versioned change and migration are required. |

## 4. Fixed sentence inventory for this exercise

This is the bounded syntax contract for the authoring exercise. It is not a parser.
E16–E20 are later candidates and are deliberately excluded from the syntax below.
`Name` denotes a resolved identifier, `Text` a quoted JSON string, and `Ref` a Name
optionally prefixed by its registered type. These are grammar metavariables.

```text
declaration: Kind Name.
source file: includes Text.
metadata: Name describes Text. | Name cites Text. | Name status Status.
containment: Name contains Name.
layer membership: Name has-layer Name. | Name groups Name.
responsibility: Name owns Name. | Name realizes Name. | Name provides Name.
dependency: Name depends-on Name. | Name consumes Name.
necessity: Name requires Name in mode Name.
boundary: Name exposes Name. | Name follows Name. | Name allows Name.
schema: Name schema Text. | Name member Text. | Name through Name.
participation: Ref uses Ref as sender of Ref [in mode Name].
participation: Ref uses Ref as receiver of Ref [in mode Name].
policy/body: Name specifies Text.
traceability: Name constrains Name. | Name replaced-by Name. | Name affects Name.
source binding: Name anchors Name. | Name at Text.
feature flow: Name illustrates Name. | Name exercises Name.
scenario: Name starts-when Name. | Name starts-at Name. | Name preserves Name.
step: Name calls Name [within Name]. | Name proceeds-to Name when Name.
traffic: Name receives Name through Name. | Name emits Name through Name.
termination: Name finishes Name with Name.
concurrency: Name runs-alongside Name.
state/data: Name tracks Name. | Name typed-as Name. | Name initially Text.
field: Name field Text as Name.
UI: Name references Name. | Name renders Name. | Name binds-to Name. | Name configures Name.
```

`Kind` is one of the core kinds plus `system`, `subsystem`, `library`, `layer`,
`channel`, `message-set`, `contract`, `client-tool`, `external-device`, `port`,
`adapter`, `feature`, `scenario`, `step`, `condition`, `outcome`, `invariant`,
`requirement`, `decision`, `gap`, `binding`, `evidence`, `data-type`, `state`,
`representation`, `composition`, `presentation`, `view`.

`Status` is one of `inherited`, `selected`, `proposed`, `superseded`, `open`,
`source-observed`. It describes the claim's authority/evidence state, not execution.
No alternative relation verbs are used. `describes`, `specifies` and condition
descriptions remain non-executable obligations; accepting their syntax would not
constitute validating their English meaning.

### Type and ownership rules

- Library, Container, Subsystem, ClientTool and Adapter are Unit-compatible for
  structural responsibility; they remain distinct for runtime allocation.
- System can contain Unit-compatible children; it cannot own Functionality.
  Unit-compatible objects can contain Unit-compatible children. Exactly one
  immediate parent per contained object; no cycles. This extends the narrower
  System-to-Container proposal because the registry includes real non-containers.
- Every Functionality has one immediate Unit-compatible owner. Layer groups
  reference these units without adding ownership. Library behavior is a reusable
  definition instantiated inside a consuming Container, not a separate server.
- `has-layer` takes Subsystem/Unit and Layer; `groups` takes Layer and Unit.
- `realizes` is Functionality to Capability; `provides` Unit to Capability;
  `consumes` Unit to Interface; `requires` Capability to Interface qualified by Mode.
- `depends-on` connects a Unit to a Library. `exposes` connects Unit to Interface;
  `exposes` also permits an explicit Port object; `follows` connects
  Channel/Interface/Port to Contract; `allows` connects Contract to
  MessageSet; `through` connects Port to Adapter. Channel participants are
  Container, ClientTool or ExternalDevice; adapter wiring is expressed separately.
- `uses` retains Container/Channel/MessageSet checking, expanded only for the
  declared external/tool endpoint kinds. Unqualified participation applies whenever
  the endpoint is instantiated. A mode-qualified sentence narrows that set, not
  an alternative spelling of unconditional participation.
- `constrains` connects Requirement/Invariant/Decision to an affected model object;
  `replaced-by` links a superseded obligation to its Decision; `affects` links Gap
  to the object whose completeness it prevents. Metadata never overrides typing.
- Binding/Evidence `anchors` a modeled responsibility and `at` supplies a pinned
  source or evidence location. A binding can remain open without a fabricated URL.
- Scenario `illustrates` Feature and `exercises` Capability; Scenario `starts-when`
  Condition, `starts-at` Step, `preserves` Invariant; Step `calls` Functionality.
  A call to library-owned behavior adds `within Container` to identify its
  consuming application context; the library is not a process or state singleton.
  `within` can also make a service-owned call context explicit. Step `receives`
  or `emits` takes MessageSet and Channel, checked against that context's
  participation and mode. It does not infer an extra channel from control flow.
  Step `proceeds-to` Step under Condition; Step `finishes` Scenario with Outcome.
  Each scenario has one entry, only reachable steps, explicit terminal outcomes
  and declared branches. A Step can be a control-only branch or terminal without
  a Functionality call. Separate alternative send Steps when only one MessageSet
  is sent; two `emits` facts on one Step would mean both effects, not implicit choice.
  `runs-alongside` records permitted overlapping scenarios,
  not an OS-thread allocation or execution schedule.
- Unit `tracks` State; State/Representation `typed-as` DataType; DataType `field`
  names a field and its DataType. Initial Text records a specified or explicitly
  unresolved starting state; it is not silently evaluated as host code.
- Composition `references` Representation/Composition; Presentation `renders`
  Composition; View `binds-to` Representation; Unit `configures` the app-specific
  Composition, Presentation or Representation definition. Design-level shared identity does
  not collapse distinct application/session instances.

## 5. Modeling choices and unresolved product decisions

Preserve all 17 registry constituents. Model StemProfile and AptDomain as explicit
BuckingService-owned Units; record possible service extraction as open. Model
the selected UI separation under a logical UI subsystem while retaining BuckingUI,
SimulatorUI and BWEB's dual upstream role. Do not silently replace BWEB with a
new UIHost process. Existing UIR/BOX names identify source/governance continuity;
their target responsibility allocation is a proposed realization of the selected
separation, not a claim that existing code has already migrated.

The target disables physical execution; PhysicalPassive describes observation
only. Protocol TX requirements are retained as obligations gated on a separately
authorized milestone. Simulator-control and Machine semantic injection remain
different paths. Captured bytes and simulator hidden truth cannot enter normal
domain paths under misleading provenance.

Unknown schema details, exact delivery budgets, stale-candidate/focus policy,
Go-host placement, newer StanForD codec/version support and complete production
provenance must remain explicit gaps. The model still describes their responsible
units and affected flows rather than omitting them from the system.

## 6. Validation and promotion plan

Before parsing, audit file membership, unique declarations, references, exact
requirement/constituent coverage and selected paths against source documents.
These are authoring/inventory checks, not an SDL parser or a language-conformance
claim. Record the results in the exercise README/coverage artifacts.

For promotion, select each E-item independently: confirm semantics, add positive
and negative language examples, amend the general definition, then implement
parsing/linking. Do not promote the entire exercise grammar merely because it
can express the example. In particular, executable predicates, typed field/value
operations and runtime scheduling need further definition before IR generation.

## 7. Source anchors

- [Constituent registry](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/CONSTITUENTS.md)
- [System Requirements](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/SDP/03--Requirements/requirements.md)
- [Domain discussion](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/Domains/DetailedArchitectureDiscussion.md)
- [UI discussion](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/Domains/UIDomain/DetailedDesignDiscussion.md)
- [Chronology and later owner decision](../../../docs/process/MVP1-Design-Evolution-and-SDP-Skills.md)

All `.design` extensions used by this exercise must be described above before
they appear in the model. New findings discovered during authoring are appended
here with their effect on the corpus and their remaining limits.

## 8. Findings from authoring the full system

The first cross-system flow review exposed a modeling trap: placing an optional
Lab observation after a production commit in one success chain makes diagnostics
look mandatory. The production scenario must branch around an absent Lab and
distinguish a committed fact with pending delivery from a rejected operation.
Neither failed observation nor failed delivery reverses an authoritative commit.
Existing Step/Condition/Outcome syntax can express this; no new synonym is needed.
This also demonstrates why a well-formed graph alone cannot verify a design.
The review applies the same rule to GlobalStop: direct Machine routing must not
depend on Bucking or Lab availability. A partially completed stop needs a distinct
outcome; failure must never imply that an already handled stop is undone.

The same separation applies to UI publication: valid/published is different from
displayed. A stale candidate can retain the old schema; renderer failure after
publication is a different outcome whose recovery policy remains open.

Responsibility review separated ProductionAuthority and BuckingCommandCoordinator
from EconomicPolicy. Shared Container ownership does not justify grouping unrelated
durable-state and stop-control responsibilities inside an economic-policy unit.

Further concrete completeness limits:

- Existing Bucking wire types do not describe the entire draft edit/store lifecycle.
  The import activation flag must be reconciled with explicit activation semantics.
- A compact ProductionEvent does not carry the full required production lineage.
- Representation/Composition/Presentation definitions need per-application and
  per-session allocation. A Library dependency alone does not specify these instances.
- Channels with request/response envelopes still need variant routing, authority,
  retention, capacity and recovery rules. Listing both endpoints is insufficient.
- Named conditions and English Functionality bodies are useful design obligations,
  but cannot produce executable IR without predicates, operations and bindings.
- Requirement-to-object links establish an obligation inventory. They do not prove
  that every clause has a sufficient design realization or a passing implementation.
- Whole-system shutdown exposes a further allocation/coordination gap: all nine
  Container lifecycles, optional Labs, target selection, quiescence order and
  failure-time cleanup must be accounted for. The scenario names the full runtime
  scope, but exact ordering/deadline/resource-instance semantics are not settled.
  A single sequential happy path cannot establish system-wide shutdown correctness.

No extension from this exercise has been promoted into the general definition.

## 9. Recorded exercise result

The [corpus and reading guide](../../../experiments/mvp1_sdl/README.md) contains 66
source files: one entry and 65 explicit includes. It covers all nine registered
Containers, all eight Libraries, 109 Functionalities, 27 Capabilities, 14 Features,
20 scenario graphs and 11 Channels. All 332 source requirement IDs are represented
with a disposition and design target; two conflicting historical UI obligations
are explicitly superseded. The model records 23 open Gaps.

The reproducible [inventory audit](../../../experiments/mvp1_sdl/inventory-audit.json)
checks authored references/ownership/flow structure and source inventory, including
31 source-file hashes. It passes. The audit is not a parser, executable semantic
check, design acceptance test or proof that broad requirement links are sufficient.
The source-linked requirement mirror and coverage index make omissions reviewable
without claiming that unknown implementation or runtime behavior has been supplied.
