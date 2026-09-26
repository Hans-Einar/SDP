# Study SDL-derived assignment bundles and blueprints

| Field | Value |
| --- | --- |
| id | KB-SDP-031 |
| project | SDP |
| type | Study |
| CardState | in-progress |
| Systems | SDL, SDPTOOL |
| created | 2026-09-25T22:20:56Z |
| source | Owner conversation 2026-09-26: SDL context, implementation drift and post-main XFMD adoption |
| next_review | At BP1 study delivery |

## Owner problem and intent

Assignments alone have not prevented agent tunnel vision, implementation drift
and loss of maintainable system structure; ponsse/Concept1 is the owner's concrete
experience. SDL was introduced to make the enclosing design explicit and usable
by workers, not merely to produce attractive diagrams. A worker needs the current
system, intended change, permitted location and desired resulting design together.
SDUI supplies a deliberately bounded language for UI design and prototyping.

## Study and proposed output

Produce a Study and a proportionate DesignPlan/ImplementationPlan for a generated
assignment bundle. Inspect the existing SDL model, analyzer/viewpoint generators,
SDPTool and the accepted [workflow study](../../Studies/UsageAnalysis/ProposedSDPWorkflow.md).
Its CurrentAssignment and per-Issue YAML proposals are historical input, not an
instruction to adopt Issue #7 schemas. Do not merge or copy that pilot as authority.

Explore a bundle containing: explicit outcome/non-goals; pinned current design and
code/evidence baseline; system purpose and neighboring responsibilities; relevant
requirements, use cases, interfaces, channels and invariants; the before/after
model change; owned/shared write surfaces; dependencies; verification obligations;
known unknowns; and escalation conditions. Define which facts are model-derived,
which must be authored by the owner/designer, and which are observations.
An extractor cannot invent intended future behavior or infer that code conforms
to SDL merely because the design parses. Include provenance and stale-bundle
invalidation. Context selection must include affected neighbors and end-to-end
paths, with links to the full model, rather than only the assigned component.

Compare compact bundles, layered context and full-model navigation for cognitive
load and drift detection. Assess preflight consistency, change-impact checks and
post-work conformance without prematurely requiring a new assignment schema.
Treat XFMD's native blueprints as a case study; their format is not SDL authority.

## Acceptance and dependencies

Use one real bounded change in a system described with SDL. Show a reproducible
current/target bundle, excluded-context rationale and review questions. Challenge
it with a plausible locally correct change that violates a neighboring contract:
record whether worker/reviewer guidance or tooling detects it, and remaining gaps.
Describe generated versus authored content and propose phased implementation.
No blueprint generator or reduction in drift is claimed by registering this card.

[KB-SDP-004](../backlog/%23004--Proposal--Design-traceability.md) owns design/code/evidence
traceability. [KB-SDP-032](../backlog/%23032--Study--Viewpoint-navigation-feedback.md) supplies
navigation feedback. [KB-SDP-033](../backlog/%23033--Study--XFMD-SDP-adoption-and-SDL-pilot.md)
provides a prospective practical pilot; the study need not wait for its completion.

## Worklog

2026-09-25T22:20:56Z: Registered in backlog before main integration; EVT-KB-SDP-000172.

## Selected study execution — BP1

Owner request, 2026-09-26: study blueprints regarding SDL. This selects research
and a proposed follow-up plan, not production implementation or language adoption.
Work on sdp/study-bp1-sdl-blueprints, based on main 2ff71c4. BranchPolicy: current;
CommitPolicy: milestone. BP1-M1 records scope; BP1-M2 delivers the evidence-backed
study, reproducible bounded experiment and proposed phased successor. Push the
completed study phase and offer review against main; no merge is authorized.

Use SDPTool's planned unsaved-source preview as the bounded example. Validate
current and proposed target SDL with the existing Go parser and derive viewpoints
with existing tools. Challenge the context with an ownership-drift variant.
A research harness may collect evidence; do not ship a blueprint command, change
production SDL/SDPTool code, adopt Issue #7 schemas or modify XFMD. Preserve the
untracked SDL/go/sourceinput draft. Study recommendations remain proposed.

2026-09-26T09:22:59Z: Backlog → active/in-progress for the owner-selected BP1 study; EVT-KB-SDP-000180.
