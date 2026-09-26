# SDL-derived blueprints and assignment bundles — BP1 study

Status: research and recommendations, 2026-09-26. Source: KB-SDP-031 and the
owner's experience of implementation drift in ponsse/Concept1. No proposed
bundle schema, language extension or product command is adopted by this study.

[Worked assignment](Pilot.md) · [Generated model extracts](Model-extracts.md) ·
[Machine evidence](Evidence.json) · [Proposed DesignPlan](Plan.md)

## Finding and recommendation

SDL already supplies validated structural facts and source-linked projections.
It does not yet supply an assignment compiler, a complete account of current code,
or the owner's desired change. A useful blueprint therefore combines generated
system context with explicitly authored intent and separately observed evidence.
Treating an SDL diagram or an entire model dump as the assignment would leave the
original tunnel-vision problem unresolved.

Recommend a layered bundle: a short task entry, a mandatory enclosing-system and
contract context, and source-linked detail on demand. Keep a small machine-readable
manifest beside readable Markdown. Generate facts with SDL; let SDPTool assemble
the assignment, resolve process records and verify provenance. Keep the original
SDL, approved task/plan and Traceability records authoritative. The bundle is a
revision-bound view of them, not another independently maintained design tree.

## Terms with distinct lifetimes

| Term | Proposed meaning | Authority / lifetime |
| --- | --- | --- |
| System design | Authored SDL plus contract/requirement documents | Long-lived desired design; not proof of implementation |
| Blueprint | A purpose-selected, source-linked explanation of a design boundary and its surroundings | Regenerable view; can serve development, review or exploration |
| Assignment | Authorized change, constraints, scope, desired result and acceptance | Existing card/plan/owner decision; not forced to be a GitHub Issue |
| Assignment bundle | Frozen blueprint context + assignment + baseline/evidence references | One revision of one bounded work episode |
| Work result | Code/model delta, verification and explicit exceptions | Commit/PR and Traceability; management lifecycle remains separate |

A blueprint is useful without assigning a worker. An assignment may concern a
bug with no intended SDL change, a proposed design change, a refactor preserving
contracts, or a study. Require an explicit target description even when it says
"same model, restore implementation conformance". Do not create a new SDL keyword
for each work-management concept.

## Evidence baseline and existing responsibilities

Inspected SDP main 2ff71c4a033aa1393e02657777e85f5fb67ccefd. The study branch changes
only study/process records; the untracked SDL/go/sourceinput draft is excluded.

| Existing surface | Reusable capability | Boundary / missing capability |
| --- | --- | --- |
| [Go parser](../../../../SDL/go/parser/ast.go), Validate and Check | Typed structural statements, declarations, source spans, canonical formatting | Design-core 0.5; no System/import/assignment grammar, no implementation proof |
| [Viewpoint model](../../../../SDL/go/viewpoint/model.go) | Facts, typed diagrams, source revision SHA-256 and source positions | f0000-style fact IDs are statement-order local; not stable cross-revision identity |
| [Selection](../../../../SDL/go/viewpoint/query.go) | Focus, direction, relation, depth and mode; validates complete source before selecting | Graph radius is not an assignment context policy; diagram and fact inventories need not have identical scope |
| [Catalog](../../../../SDL/go/viewpoint/catalog.go) | VP01–VP11: goals, responsibilities, structure, interfaces, plans, channels, data and gaps | No dedicated blueprint view; consumes alone does not identify provider or executable flow |
| [SDL document service](../../../../SDL/go/documents/selection.go) | Generated Markdown/Mermaid and optional SVG, source facts and guarded bundle publication | A document bundle is not an authorized work package |
| [SDPTool producer contract](../../../../Toolkit/SDPTool/Contract.md) | Project/model resolution, typed tree/select, source-revision checks, bounded output and consumer ownership | No assignment generation, proposed-target comparison or model-to-code evidence resolver |
| [Traceability proposal](../../../KanBan/backlog/%23004--Proposal--Design-traceability.md) | Existing ledger/relations are the place to connect work and evidence | Exact model-identity joins, invalidation and aggregation still need an adopted contract |
| [Earlier workflow study](../../../Studies/UsageAnalysis/ProposedSDPWorkflow.md) | Authored assignment versus observed state, owned/shared surfaces, integration baseline and invariants | Historical proposal, not adopted Issue #7 schema or mandatory per-Issue storage |

Structural, action-core and class-core are separate supported profiles. This first
pilot uses only structural SDL; it must not imply an automatic join to executable
actions, class implementation or SDUI widget callbacks. File paths and Go symbols
are not modeled by structural owns/contains, and must not be guessed from names.

The inspected SDPTool README still contains an older statement denying XFMD's SDP
bootstrap and future-tense navigation wording. The implemented Contract, Go source
and delivered evidence govern this study's capability claims. This documentation
inconsistency is referred to existing KB-SDP-010; it does not justify changing the
producer contract or the XFMD repository during blueprint research.

## What the experiment demonstrates

The bounded example is the already planned unsaved-source preview, associated
with REQ-SDPTOOL-004 and UnsavedSourcePreview. Current code accepts saved files.
The proposed target introduces ValidateSnapshotRevision as an explicit design
responsibility while retaining ViewerBridge's delivery responsibility. This is
an illustrative design refinement, not a claim that unsaved preview is implemented.

[probe.py](probe.py) invokes the existing Go parser, canonicalizer and viewpoint
CLI. It obtains the baseline from the pinned Git commit and reconstructs target
and drift inputs; it neither parses SDL in Python nor adds a production generator.

| Case | Existing SDL validation | Authored delivery-owner guard |
| --- | --- | --- |
| Current model | Pass | Pass |
| Proposed target | Pass | Pass |
| Target with DeliverViewBundle moved to PreviewCoordinator | Pass | Reject |

The narrow current/target VP03 selection around PreviewCoordinator omits
ViewerBridge's ownership. The container context exposes it. For the current model
the narrow selection contains 19 facts; the broader selection contains 93 facts,
even when the displayed diagram is restricted to DesignPreview. A manifest fact
inventory is not automatically a minimal explanation of the visible diagram.

All three models are language-valid. The rejection comes from a separately authored
assignment constraint evaluated against generated facts. This is the distinction
between valid design and permitted change. The probe does not observe a worker,
analyze changed Go code or demonstrate that a generated blueprint prevents drift.
The real saved-file stale/preservation/navigation tests also pass; they establish
existing behavior only, not the proposed unsaved-buffer workflow.

## Alternatives considered

| Approach | Benefit | Limitation | Disposition |
| --- | --- | --- | --- |
| Full model and all viewpoints for every task | Broad discoverability | Large reading burden; no explicit relevance or protected boundary; stale snapshots still possible | Keep accessible as a reference, not mandatory inline content |
| Fixed-depth graph neighborhood only | Simple reuse of current selection | Important actors, sibling responsibilities or atomic protocols can be outside radius; shared hubs can pull in irrelevant work | Useful query primitive, insufficient selection policy |
| Hand-maintained per-feature blueprints | Rich contracts, failures and code mappings | Duplicates SDL facts and can drift independently | Reuse authored contract detail; generate structural portions |
| Layered generated context plus explicit change contract | Shows purpose, neighbors, intended delta and evidence gaps | Requires provenance, context rules and honest unknowns | Recommended first implementation direction |

Do not add embeddings, an agent orchestrator, a new database or a generic graph
engine to validate this idea. Existing typed SDL facts, Go libraries, Markdown,
Git and project records are sufficient for a first vertical pilot.

## Proposed context-selection policy

Start with an explicitly selected model and work anchors: functionality/use case,
activity or named unit, plus the intended model delta where one exists. Validate
the complete source before projecting; never reparse a sliced text fragment and
present it as a complete model.

1. Include purpose: anchors → contributes-to feature → supported use cases and
   actors, with relevant capability/mode/allocation. Display missing links as gaps.
2. Include placement: owner and containing chain to the current model boundary;
   identify neighboring responsibilities that share touched contracts. Do not infer
   a System boundary from the container name or repository folder.
3. Include all incident typed contracts for touched interfaces/channels/data.
   Preserve complete scenario steps, request/reply relationships, packet fields,
   encoding and placement when such a contract is selected. These are structured
   multi-party facts; a traversal of only subject/object misses roles and endpoints.
4. Include delivery dependencies and explicit prerequisite activities, then link
   evidence through reviewed mappings. Status in SDL remains a declared claim.
5. Compute selection for both current and proposed target, taking their union.
   Otherwise a removed edge can also remove the very neighbor a reviewer needs.
6. Attach exact exclusions and frontier: what is omitted, why, which contracts
   cross the boundary, and where to expand. Depth/size limits produce an incomplete
   result with a reason, never an apparently complete truncated assignment.

An include needs a reason such as "owner of anchor", "consumer of changed
contract" or "accepted invariant". External/unknown providers remain unknown and
can block an assignment's readiness where that missing contract matters. Large
context requires splitting or explicit review, not silently dropping obligations.
Rendering is a consumer; semantic selection must not depend on diagram coordinates.

## Proposed bundle composition and authority

| Input/output | Source of truth | How the bundle should represent it |
| --- | --- | --- |
| Outcome, non-goals, allowed changes, protected behavior | Authored assignment and owner decisions | Link exact version; show concise required reading |
| Current design and proposed target | Separate complete model snapshots or reviewed target delta | Preserve both identities; highlight additions/removals and unchanged boundary |
| Model facts and diagrams | Existing SDL parser/projection | Derived, with profile, model revision, source locations and selection reasons |
| Paths/symbols and neighboring consumers | Reviewed design-to-code mappings, code inspection | Observed at commit/hash; distinguish observed code from planned mapping |
| Verification and acceptance | Contract documents, selected plan and Traceability | Required versus passed versus stale/unknown; never infer from card closure |
| Generated Markdown/resources | SDPTool-owned output bundle | Read-only view; publish atomically and preserve authored source |

A proposed manifest needs an assignment reference and revision, qualified model
key (project/model/profile), code baseline, exact source/contract/evidence hashes,
selection-policy and generator versions, target identity, included facts and
boundary frontier, authored constraints, code mappings, required checks and
artifact digests. This is a field study, not an installed JSON schema. Keep short
human metadata visible as tables; do not hide the task contract in frontmatter.

Use model revision plus declaration name/span for current identity. Fact ordinals
alone cannot survive formatting or insertion. A later semantic diff can compare
typed tuples and normalized source while preserving raw-byte provenance; renames
need explicit mapping rather than similarity guesses. Multi-file source identities
must follow KB-SDL-005 when implemented. The pilot can use today's single-file
models; source composition is not a prerequisite for demonstrating its value.

## Change lifecycle and drift checks

Prepare → inspect gaps/context → owner or authorized designer selects the change
contract → freeze bundle → worker preflight → bounded implementation → compare
actual changes → review evidence and integration. These are proposed operations,
not new KanBan directories or adopted process states.

Preflight verifies relevant hashes, profiles, model anchors, permitted/shared write
surfaces and declared dependencies. A newer unrelated commit does not by itself
prove semantic staleness; initially require exact baseline or explicit rebase and
regeneration. Relevant input changes always invalidate the frozen bundle. An
invalid refresh preserves the prior artifact but marks it stale and unusable for
new execution. A bundle cannot grant permissions beyond the originating assignment.

During work, discovery of an unlisted consumer or required contract change triggers
an explicit amendment and regeneration. Keep the prior bundle for comparison;
never rewrite it to make an unauthorized change appear pre-approved. Parallel
assignments identify shared touchpoints and ordering through existing process
records, not independent global CurrentAssignment pointers.

At handoff compare actual model delta with intended delta, changed paths with
reviewed scope, and required evidence with actual results. Unmodeled behavior and
handwritten domain code still require ordinary engineering review and integration
tests. A filename or symbol existing does not prove callers, error handling or
protocol compatibility. Output reports checked, violated, stale and unknown
obligations; avoid a single misleading "conforms" badge.

## Lessons from XFMD's native blueprints

Read-only inspection of xfmd-sdl-navigation at
cf11709e4ec9d6925b0d17d95c71f72fbf011959 examined its blueprint registry,
Functionality-031 Generated Document Navigation, validate_blueprints.py and
check_blueprint_symbols.py. No XFMD source or process record was changed.

Useful material to preserve: purpose, acceptance IDs, ownership, behavior/failure
paths, an explicit caller/callee plumbing table, dependencies, evidence limitations
and change impact. FUNC-031 clearly distinguishes best-effort lease release from
durable acknowledgment and assigns application versus SDL responsibilities. These
are precisely the surrounding contracts that a narrow implementation task can miss.

The structure validator checks metadata, links and coverage of named acceptance
references. The symbol check looks for declared callees in files and explicitly
does not establish call-graph semantics. Neither is SDL-derived design conformance.
Borrow the questions, not its project-specific FTR/FUNC taxonomy, fixed eight-section
format or storage tree. This inspection is not a current full XFMD blueprint audit.

## Decision boundary and next work

Recommend the [DesignPlan](Plan.md): freeze a minimal contract, specify typed context
selection/diff and demonstrate one complete review workflow before implementation.
Product ownership remains SDL for language/model semantics, SDPTool for assembly
and process/evidence coordination, and XFMD for a possible viewer. No new renderer
or separate blueprint runtime is needed for the proposed first delivery.

Open decisions for that plan: minimal mandatory human-authored contract; reviewed
model-to-code mapping format with KB-SDP-004; context budget and readiness treatment
of unknown neighbors; and trial acceptance with an actual worker/reviewer. These
are bounded design questions, not blockers to completing this study.

KB-SDP-032 supplies later owner navigation feedback; KB-SDP-033 can provide an XFMD
pilot after adoption. Neither must block the smaller SDPTool-based pilot here.
Issue #7 remains excluded. Study completion does not authorize the proposed plan,
claim a new Feature is implemented, or record owner approval of these recommendations.

## Delivery verification — BP1-M2

Local checks on BP1's document/experiment delta pass: full Toolkit validation,
management replay (40 cards, 11 management records, three lineage operations,
267 events), six management test groups, and document verification (105 frozen
records/prefixes, 574 pre-existing generated outputs, 2,501 local links and 130
fragments). The two new retained research artifacts are checked separately by
byte-for-byte reproduction; their probe hash matches Evidence.json. Three named
SDPTool regression tests pass as recorded in the worked assignment.

The author performed consistency checks; no independent review or worker trial
is claimed. Canonical SDL models, production code, historical ledger prefixes
and unrelated sourceinput work remain unchanged. KB-SDP-031 closes as a delivered
Study, and PLAN-SDP-0001 remains planned. The phase is submitted for review against
main; merging is outside this study authorization.
