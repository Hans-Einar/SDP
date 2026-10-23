# Blueprint owner-intent assessment — 24 October 2026

Discussion-only assessment of current source and evidence, owner steering during
Session0007 T029; linked continuity in Session0008. KB050 remains gate-review and
Session0008 remains owner-paused for execution. This is not acceptance of the
agent-selected pilot, a new implementation assignment or a Ponsse model change.

Current owner scope — 24 October 2026: keep the next context/Blueprint delivery
SDL-based, finish the existing system and try it. Code-analysis/runtime proposals
are deferred under KB059 and are not adoption gates. Existing engineering tests
and verification obligations remain; no extra data-source feature is selected.

## Owner clarification

The owner wants agents to understand a large system and the impact of a bounded
change without spending most of their time/context discovering the source tree.
The package should convey broad system orientation, relevant local detail, the
current design and the intended resulting design. The owner has not participated
in all Blueprint details and expects the product to evolve through actual use.
This clarification is not feedback approving the existing calibration specimen.

## Assessment

The foundation is aligned with that idea. It is not yet evidence that the complete
agent-understanding experience is achieved. Existing parser-backed NOW/TARGET,
conservative affected context, preserved obligations, unknowns and revision identity
are useful primitives; they should be retained rather than replaced.

Sources inspected: [contract](Contract.md), [selection rules](Selection-and-Evidence.md),
[producer/handoff](Producer-and-Handoff.md), SDL/go/blueprint/analyze.go and render.go,
SDPTool/blueprints/README.md and
[generated assignment trial](../../../05--Implementation/SDPTool/Blueprints/Evidence-BPI3-M1.md).
This is source/evidence inspection, not a newly executed user study or benchmark.

| Owner need | Current evidence | Gap / proposed refinement |
| --- | --- | --- |
| Exact current and target design | Complete captured sources, distinct revisions, artifact identities and source provenance | NOW is a selected modeled baseline. Whether it matches implemented code remains a separate evidence question. Expose that correspondence prominently. |
| Understand the change and its environment | Typed semantic diff; context from NOW/TARGET union, including removed neighbors, channel groups and authored inclusions | The generated entry has task intent and links, but little system orientation. Add a short sourced explanation of purpose, actors, major responsibilities and where this task sits. |
| Broad view plus relevant local detail | Changed elements seed a conservative bidirectional closure; System/mode labels stop propagation; raw full sources retained | The selected context is not automatically a broad architecture overview. Provide a small global orientation separately from detailed affected context; do not expand every sibling through a System label. |
| Distinguish before, after and change | Per-side provenance and complete snapshots; current render produces union-context and delta diagrams | Add explicit NOW and TARGET views or an equally clear side selector, alongside the delta. A mixed graph alone is harder to read as the two actual designs. |
| Know what must remain intact | Authored permissions/protection/rules, CHANGE/CONTEXT/PRESERVE labels and negative cases | Add the source and reason for important invariants to the readable entry; design intent and rationale cannot be inferred from graph edges alone. |
| Locate the necessary implementation | Scoped code mappings/check receipts can be assessed; missing/stale/ambiguous mappings stay unknown | The analyzer does not automatically discover verified model-to-code mappings. Provide honest file/symbol locators and code revision/freshness, with unresolved areas explicit. |
| Save context on large systems | Deterministic machine/readable bundle, retained catalogue, resource bounds and impact closure | A connected closure can still become large; default bounds are 2000 nodes/10000 facts/8 MiB. Failure on overflow is sound but does not create a good short entry. Add progressive delivery and explicit expansion. |
| Demonstrate agents can use it | BPI3-M1 synthetic generated-bundle Worker/Reviewer trial succeeds; protected-model and code-only regressions are rejected | This is positive bounded evidence, not proof of understanding/efficiency on a representative large codebase. A later owner-selected trial should measure necessary source reads, comprehension gaps and result correctness. |

## Recommended next design refinements

These are proposals within KB050/KB004 scope, not implemented features or adopted
new schemas. Do not activate them while the current PG compatibility task is active.

1. **Layer the handoff.** Start with task intent, a small broad system orientation,
   the task's place in it, important decisions/invariants and coverage status. Link
   to detailed affected NOW/TARGET views, contracts, exact changes and source/code
   locators. The first page should help an unfamiliar agent decide what to read next.
2. **Make expansion explicit.** Retain complete pinned inputs, but deliver a short
   initial package with element/contract/source handles for read-only expansion.
   Keep required affected closure visible as obligations even when details are
   expanded later. A token budget is a delivery constraint, not a claim that omitted
   relationships cannot matter. Do not silently summarize away protected contracts.
3. **Expose model/code trust.** Put modeled NOW, verified correspondence, stale
   mappings and unmodeled runtime dependencies on the entry page. Code-only drift
   cannot be detected by a model diff. Focused implementation-source reads and
   independent checks remain necessary when correspondence is unknown.
4. **Validate the user's experience.** When the owner selects the Blueprint pilot,
   walk a concrete package together before adding more workflow machinery. Then use
   a representative bounded task to see whether an agent can explain current design,
   target change, preserved behavior and uncertainties with fewer exploratory reads.
   Use the existing Ponsse pilot/model only when its separate deferral is lifted.

The current producer analyzes one checked design-core/0.6 System per comparison.
Cross-System/ecosystem context, runtime call order, behavioral protocol semantics
and SDUI analysis are not automatically supplied. Broader orientation may link
pinned neighboring models/documents while preserving explicit unavailable details;
it must not pretend the current single-System analyzer has complete global impact.

## Proposed readable entry

A useful entry should answer, in this order: What is this system for? Where does
this task fit? What is modeled NOW and how well is it tied to code? What changes
in TARGET and why? Which behavior/interfaces must remain? What exact details and
code should be read next? Which unknowns require investigation or an owner decision?
Generated facts and sourced authored rationale should be distinguishable.

That progression is the recommended improvement over the current minimal index,
not a replacement parser, another task authority or another management ledger.
A structurally passing bundle is one input to an assignment, not automatic proof
that the agent understands or that the implementation will conform.

## Disposition

Agent recommendation: continue on the existing technical foundation, prioritizing
orientation/progressive context and model-code correspondence before more lifecycle
features. No owner pilot acceptance or Blueprint implementation resumption is
inferred. The current authorized implementation remains Codex compatibility PGC1.
