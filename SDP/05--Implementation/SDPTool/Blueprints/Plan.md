# BPI — bounded semantic blueprint producer

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0020 |
| project | SDP |
| state | planned |
| PlanType | ImplementationPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDL, SDPTOOL |
| source | KB-SDP-050 / PLAN-SDP-0001 BP2-C |

## Outcome and authority

Implement a Go blueprint producer from captured NOW/TARGET models and an authored
task: semantic differences, affected modeled context, explicit unknowns and preserved
obligations, packaged as deterministic Markdown/Mermaid plus machine data. Inputs and
ready-state authority must remain distinct from implementation-conformance proof.
This plan is prepared by the authorized BP2 design work, not yet executing.

Governing design: ../../../04--Design/SDPTool/Blueprints/Contract.md,
Selection-and-Evidence.md and Producer-and-Handoff.md. Use the reduced MVP1 fixtures
and negative cases, not an unsupported full experimental-profile migration.
No production Ponsse/XFMD changes, mandatory Git, grammar expansion, new renderer,
automatic code execution or release packaging. SDUI bytes are retained with explicit
coverage gaps; a semantic SDUI analyzer is not silently bundled into this increment.

Use a dedicated working branch based on the integrated dependencies when execution
starts; do not switch a concurrently used worktree. Per-milestone commits and phase
pushes; no main merge or publication without owner authorization. ModelGovernance
and sourcegraph APIs are prerequisites, not dependencies on their future product release.

## Phases and acceptance

| Phase | Milestone | Complete vertical outcome |
| --- | --- | --- |
| BPI1 Analysis | BPI1-M1 | Go pure analyzer accepts checked 0.6 snapshots/task; reports typed diff, context, coverage and constraints from actual parser facts |
| BPI2 Publication | BPI2-M1 | SDPTool captures artifact inputs and produces deterministic Markdown/Mermaid/JSON bundle through shared output/publisher boundaries |
| BPI3 Assignment evidence | BPI3-M1 | Explicit readiness/dispositions, pinned before/after check receipts and independent generated-bundle worker/reviewer trial |

### BPI1 — analysis

Implement SDL/go/blueprint with no filesystem or lifecycle dependence. Define strict
Task/Policy/Analysis contracts and version IDs before the adapter. Check System/profile
identity, canonical semantic keys, all typed fields, codepoint/path escaping and
per-side origins. Table-driven policy must cover every accepted AST statement or
explicitly reject unsupported semantics. Implement atomic contract closure, cycles,
cut/excluded reasons, constraints and unknown frontier. Never infer code conformance.

Tests: existing composed NOW/TARGET, same bytes repeated, source-only formatting/move,
renames/kind changes, removed-edge neighbors, conflicting CHANGE/PRESERVE, authored
exclusion through a mandatory contract, absent peer, cycles, cross-system/profile
mismatch and node/fact/byte overflow. Every successful output must explain its coverage.
Port only proven experiment logic; Python remains historical design evidence, not
production dependency. Treat injected normalized graphs separately from parser cases.

### BPI2 — end-to-end generation

Thin proposed `model create blueprint from ... to ...` facade captures owned source
maps and task bytes, compiles the exact captured sources with existing Go APIs and
routes results through existing human/--json presentation. No process-template or
Git requirement. No source registration list. Outputs as specified in the producer
contract, with deterministic hashes and per-side source links. Reuse document
publishing only after testing containment, generated/unmanaged ownership and failure
behavior. Avoid changes to unrelated navigation, model mutation or install protocols.

Tests: compiled CLI in a fresh non-Git area, preliminary WORK without implicit commit,
frozen candidate parity, new edit -> new identity, task/policy changes invalidate
identity, invalid sources, unsupported SDUI coverage surfaced, malicious labels/paths,
changed source/task before publish, cancellation, existing notes, edited managed
outputs, write/rename failure and unchanged prior successful output. Check every link.
Same inputs/tool/policy/evidence must yield byte-identical managed payloads. Render
selected diagrams through the existing renderer as compatibility evidence, not manual
replacements. Full blueprint generation remains useful without XFMD.

### BPI3 — readiness and evidence

Implement explicit readiness separately from preview. Version strict check receipts
and unknown dispositions; bind source/task/code/evidence hashes and scope. Validate
missing/stale/ambiguous code mapping as UNKNOWN. Do not execute referenced commands
or claim authenticated acceptance. Unknown disposition does not convert not-run to
pass. Integrate existing Traceability references rather than a second ledger.

Run authorized bounded Worker then fresh independent Reviewer using an actual
generated bundle and both compliant and violating controls. Include a code-only
violation with unchanged SDL to demonstrate the semantic check's limit. Keep general
mapping/runtime features in KB004/KB-SDL-006. Owner feedback on the design specimen
must be recorded before claiming owner acceptance of that particular Ponsse change.

## Verification and closeout

Run targeted Go tests at each milestone, then relevant SDL/SDPTool full suites with
race/vet and compiled CLI workflow on the final candidate. Independent review must
read owner intent before implementation narrative and have all material findings
resolved or explicitly scoped out before closure. Record exact hashes/toolchain,
failed checks and limits. Update source design projections only from validated SDL.
No parser pass substitutes for behavior or cross-project GUI acceptance.

Complete KB050 only after delivered producer, evidence and documentation agree.
Session0008 continues through implementation; its roadmap must show design gates
and production milestones separately. Installer manifests, releases and integration
into XFMD require their own authorized work. No milestone has started here.
