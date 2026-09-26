# Study SDL-derived assignment bundles and blueprints

| Field | Value |
| --- | --- |
| id | KB-SDP-031 |
| project | SDP |
| type | Study |
| CardState | backlog |
| Systems | SDL, SDPTOOL |
| created | 2026-09-25T22:20:56Z |
| source | Owner conversation 2026-09-26: SDL context, implementation drift and post-main XFMD adoption |
| next_review | After MAINT-SDP-0005 main consolidation |

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

[KB-SDP-004](%23004--Proposal--Design-traceability.md) owns design/code/evidence
traceability. [KB-SDP-032](%23032--Study--Viewpoint-navigation-feedback.md) supplies
navigation feedback. [KB-SDP-033](%23033--Study--XFMD-SDP-adoption-and-SDL-pilot.md)
provides a prospective practical pilot; the study need not wait for its completion.

## Worklog

2026-09-25T22:20:56Z: Registered in backlog before main integration; EVT-KB-SDP-000172.
