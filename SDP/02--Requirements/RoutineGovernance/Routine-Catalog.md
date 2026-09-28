# Candidate routine coverage

This is a proposed catalog for [RGS1](Study.md), not executable definitions or
new process authority. IDs below are study-local families. Actual registry names,
versioning and installed paths belong to the next selected DesignPlan. Reuse
current skills, KanBan and typed plans rather than creating another documentation
hierarchy. A procedure is not a new KanBan card type.

Further coordinated study work and the result register are now collected in
[the shared RGS2 StudyPlan](StudyPlan.md). The original RGS1 delivery remains
recorded under completed PLAN-SDP-0006.

## Matching model

Match action and context together. A message can combine categories; select a
primary routine and its necessary nested operations. Reuse the current plan when
it covers the request. Domain facets specialize obligations without multiplying
families into one routine per language, repository and artifact type.

| Family | Trigger and exclusions | Required context / guard | Durable result |
| --- | --- | --- | --- |
| RT01 — Explain or inspect | Question, status or read-only navigation; excludes implied instructions to modify | Identify factual source and freshness; a trivial unrelated answer needs no project intake | Answer with source/uncertainty; normally no new work record |
| RT02 — Capture and relate | Idea, observation or future requirement; capture does not activate | Search existing cards/decisions; preserve owner wording and affected systems | Existing card enriched or one primary card with explicit external references |
| RT03 — Discover/adopt a project | New project, imported legacy repository, SDP install/adoption inquiry | Separate installation from working-method adoption; resolve nested roots and current authority | Bounded adoption study/plan and declared coverage, not automatic process migration |
| RT04 — Diagnose or study | Unknown cause, feasibility, reverse engineering or competing options | Evidence scope, candidate identity, uncertainty and domain restrictions | Findings, alternatives, recommendation and unresolved questions; no implicit fixes |
| RT05 — Assess fit and impact | Proposed capability or change affects consumers, ownership or scope | User outcome, existing commitments, SDL coverage, dependencies and unknown frontier | Bounded impact/decision basis feeding plan/blueprint; no requirement for a full model |
| RT06 — Consolidate and prioritize | Backlog review, competing ideas, merge/split, optional Scrum/Sprint | Preserve lineage and owner selection; record deferrals and remaining work | Current cards/management events; Sprint only when useful and selected |
| RT07 — Plan or revise | Selected work needs execution structure or active scope changes | Correct PlanType; authorization, evidence, phases/milestones, explicit Git policy | Proportionate selected/revised plan; one message need not create one plan |
| RT08 — Prepare an assignment | Worker/delegation requested within authorized host constraints | Ready boundary in Study.md; owner/role limits; NOW/TARGET and observed context | Revision-bound assignment/blueprint reference and explicit stop rules |
| RT09 — Implement a bounded change | Code, SDL/SDUI model or product behavior change | Applicable card/plan/milestone, candidate baseline, owned paths and invariants | Scoped change plus appropriate tests/evidence and recorded emerging gaps |
| RT10 — Maintain/migrate/refactor | Templates, dependencies, process cleanup, schema/layout migration, structural product change | Maintenance versus product Refactor distinction; preservation and rollback; consumer impact | Selected plan delivery and migration evidence, not merely moved files |
| RT11 — Verify | Establish whether an outcome is proven | Correct candidate, environment and function/service/workflow/application level | Reproducible result with limits; fake and real integration remain distinct |
| RT12 — Review and decide | Independent review or owner disposition | Assigned role, independent context where required, explicit decision authority | Findings/disposition/rework links; self-check cannot become independent approval |
| RT13 — Integrate and publish | Merge, tag, release, deploy/install or upgrade | Separate authorizations and exact candidate; preserve user work and external effects | Actual remote/installation observations, receipts and reconciliation; no invented success |
| RT14 — Resume, redirect or recover | Continue, context loss, correction, pause, cancellation, failure | Recover active objective and authorization; compare current inputs and last evidence | Resumed or bounded stopped work; stale evidence invalidated; unresolved scope retained |
| RT15 — Investigate a procedure gap | Needed route missing/conflicting/incompatible; recurrent exceptions | Distinguish unavailable from missing; deduplicate; safe fallback boundaries | Linked gap/disposition and candidate proposal only if selected |
| RT16 — Evolve/retire procedures | Selected process improvement or obsolete/overlapping routine | Trial cases, version compatibility, adoption authority and preservation of old runs | Reviewed versioned change; existing install/release process used for distribution |

RT09 and RT10 need current KanBan/plan coverage before substantive mutation.
RT01 must not become an excuse to perform a requested code change without that
coverage. RT12 may complete even though implementation findings remain in backlog.
RT13 groups related operations conceptually but does not combine their permissions.
RT14 preserves earlier authorization; it does not ask the owner to approve the
same bounded step again merely because a session changed.

## Domain and artifact refinements

| Facet | Additional concerns | Source of actual authority |
| --- | --- | --- |
| SDP process/tool distribution | Canonical versus generated sources, installed profile, signed inventory, old-run compatibility, consumer upgrades | Current SDP plan/release/install contracts |
| SDL source/model work | Supported language profile, declarations/links, model coverage, generated views, experimental syntax | Installed SDL contract and selected design authority |
| SDUI layout/presentation | Parser/AST versus static preview versus interactive runtime; layout invariants and backend boundaries | Selected SDUI specification and plan |
| XFMD application | Application/interpreter/renderer ownership; document navigation; build versus install and actual user workflow | XFMD's adopted working method and bounded assignment |
| Ponsse Concept1/MVP1 | Correct nested project, stopped/manual-review disposition, offline/physical-I/O limits, separate UI containers | Ponsse instructions and selected project/branch records |
| HSX / debugger work | Selected branch/track, shared contracts and runtime/tool consumers | Selected branch's process/design plus actual integration evidence |
| Data-oriented UI and planning | Derived views, migration and cancellation/restart semantics | Project contracts and observed affected consumers |
| Processor emulator/DAP | Generic engine versus board policy, paired revisions, fake versus real-runtime evidence | Provider/consumer contracts and project-specific review requirements |
| External project | Unverified reference permitted; required peer evidence must be explicit and accessible | Named source/revision, not a guessed local checkout |

Facets express questions to resolve. They do not import one project's rules into
another or make every source in Evidence.md current authority for all future work.

## Composite examples

**“Remove PowerShell and release the new installer.”** Resolve current
Maintenance/release work → assess public compatibility → revise the selected plan
if needed → implement → verify exact payload → review → publish only within
existing authority → reconcile consumer upgrade. Record an unexpected affected
entry point as a scope finding, not a silent extra fix.

**“Why does the XFMD preview not open this design file?”** Inspect/diagnose first.
Determine whether the failure belongs to source resolution, supported SDL syntax,
SDPTool generation or XFMD rendering. A reproducible failure may justify a bounded
fix under its project's plan. Opening a card in SDP does not authorize XFMD edits.

**“Continue Ponsse.”** Resolve project and active candidate. A documented stop
pending owner/manual acceptance is not erased by vague continuation; determine
whether the new instruction actually supplies the missing disposition. Continue
independent authorized analysis while the dependent action remains unresolved.

**“This is a repeated step with no routine.”** RT15 captures evidence and searches
for an existing gap. It may recommend RT16, but does not activate procedure
implementation merely because the agent has discovered a useful pattern.

## Minimum procedure definition to design

A candidate definition should be short enough to inspect and complete enough to
execute consistently:

- Stable ID, version/hash, schema, status and supersession.
- Purpose, positive matches, negative matches and required project capabilities.
- Authority references, required inputs and entry conditions.
- Roles, steps, dependencies, allowed transitions and nested routine references.
- Evidence and completion rules, including stale-input invalidation.
- Waiting, failure, recovery, exception, cancellation and gap-detection behavior.
- Monitoring labels and links derived from those same definitions.
- Trial cases and adoption/distribution record.

Reference existing role instructions instead of copying their prose into every
routine. Generate the diagram from the executable definition when that exists;
never maintain a separate diagram with different steps. A routine may defer a
semantic decision to a human/agent role, while making the required inputs and
permitted resulting transitions explicit.

## Proportionality and retirement

A routine earns its maintenance cost when it addresses observed risk, repeated
coordination or consequential uncertainty. Track false matches and exceptions.
Merge overlapping definitions with clear lineage; retire obsolete ones while
keeping old run versions readable. Do not require a Scrum, Sprint, new plan,
independent reviewer or owner prompt for every request. Required review is driven
by actual project authority and selected work, not by the existence of RT12.
