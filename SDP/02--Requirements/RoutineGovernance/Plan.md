# RGS1 — extended routine and project-governance study

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0006 |
| project | SDP |
| state | completed |
| PlanType | RequirementPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDPTOOL, SDL, SDUI |
| source | Owner request 2026-09-28; KB-SDP-037 and KB-SDP-036 |

## Outcome and authorization

The owner requests an extended study of project work before and after blueprint
assignment, including SDP, XFMD, Ponsse, HSX, agro-crm, TerrainAnalyzer,
GrassPhenology and the P1000 processor simulator/debugger. Identify reusable
request/routine coverage, unknown-routine detection, distributed procedure
ownership, MCP knowledge/action interfaces and observable execution state.
Deliver evidence-backed recommendations and a bounded next-step roadmap.

This authorizes research and maintained study/planning records, not production
engine/skill/schema changes, project migrations, MCP configuration, publication
or application edits. External projects are inspected read-only. Raw private
project files stay outside this public repository; record minimal process-level
findings, source identities and limitations without publishing domain content.
Owner-reported drift is evidence of the owner's experience, not an independently
established causal diagnosis of every repository.

## Current authority and Git policy

Use current five-phase SDP, shared KanBan and typed plans. Earlier Issue-first
and SDP 2.0 proposals are research inputs, not silently adopted rules. Reuse
KB-SDP-036/037 and planned BP2 rather than activate a full governance migration.
Work on sdp/request-routine; commit each milestone and push the study delivery.
Existing authorization covers pushes and a reviewable PR, not merging this study.
Preserve unrelated sourceinput/node files and all historical ledger bytes.

## Phases and milestones

| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| A — evidence | RGS1-A-M1 | Pin current local/remote project samples, reconcile historical studies, inspect relevant tool boundaries and official MCP/host documentation; distinguish findings from owner reports | Complete |
| B — proposed contract | RGS1-B-M1 | Cover intent-to-assignment and execution/release/learning; specify candidate categories, procedure gaps, procedure distribution, state/evidence ownership and observation/enforcement limits | Complete |
| C — challenge and handoff | RGS1-C-M1 | Walk concrete cross-project scenarios and negative cases; deliver testable requirements, alternatives, staged roadmap and explicit decisions/limits; validate links and management history | Complete |

## Verification and closeout

Use pinned read-only evidence and analytical scenario walkthroughs. A scenario
walkthrough is not execution of an implemented workflow engine or proof of
prevented drift. No agent trial or independent review is claimed unless actually
performed. Verify local links and management records. Study completion returns
KB-SDP-037 to backlog with its larger runtime/observer delivery explicitly open;
KB-SDP-036 remains backlog. The study itself can complete without claiming the
capability is delivered or forcing the owner to review everything immediately.

## Deliverables

[Study.md](Study.md) is the main reading entry. [Evidence.md](Evidence.md) identifies
inspected sources and limits, with [source pins](Source-pins.json).
[Routine-Catalog.md](Routine-Catalog.md) describes candidate coverage.
[Scenarios.md](Scenarios.md) records worked cases and future acceptance tests.
Recommendations remain proposed until
the owner selects a bounded implementation or design plan.

## Delivery log

- RGS1-A-M1: Evidence.md and Source-pins.json record 33 sampled source files,
  current PR observations and official interface boundaries. No product tests,
  independent review or host-hook enforcement are claimed.
- RGS1-B-M1: Study.md and Routine-Catalog.md specify 16 proposed requirements,
  16 routine families, gap handling, versioned distribution and separated state/
  authority/evidence. They propose a small Go SDPTool/MCP pilot; no production
  policy, engine, schema, host configuration or application changed.
- RGS1-C-M1: Scenarios.md challenges the proposal with 16 portfolio and failure
  cases. Same-author consistency checks resolve local links, all 16 requirement
  mappings and 33 acquired-source hashes. These checks prove record consistency,
  not an implemented engine or prevention of product drift.

## Outcome and remaining work

The selected study is complete. KB-SDP-037 returns to backlog because its engine,
MCP and observer capability remains unimplemented; KB-SDP-036 retains routing/
entry scope. The recommendation is a shared bounded DesignPlan, followed by one
complete maintenance pilot before broader blueprint/project rollout. BP2 remains
planned. No independent review, product execution, external repository mutation,
MCP/host setup or new process adoption is claimed.

Source and analysis deliveries: RGS1-A-M1 commit 525f59b; RGS1-B-M1 commit cc211ae.
RGS1-C-M1 is the closeout commit carrying this final record and the linked board
transitions. See its Git history for exact identity. Existing ledger bytes are
preserved; new management events record this research only.
