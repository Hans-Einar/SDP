# PGI — One governed Maintenance workflow

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0018 |
| project | SDP |
| state | planned |
| PlanType | ImplementationPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDPTOOL |
| source | PLAN-SDP-0017; KB-SDP-038 with KB036/037/042 contributions |

## Outcome and authority

Deliver one useful owner-to-agent Maintenance workflow through a local Codex client,
shared SDPTool process core and agent adapter. Restore the same Session after an
interrupted attempt and return candidate-bound evidence with a truthful disposition.
This is the implementation handoff from the owner-selected PGD1 design. It is planned,
not started; the 2026-10-03 continuation selected preparation of this handoff.

The [design](../../../04--Design/SDPTool/ProjectGovernance/Design.md),
[acceptance cases](../../../04--Design/SDPTool/ProjectGovernance/Acceptance.md) and
[Session 0007](../../../Sessions/session-%230007--Project_governance.md) govern
scope. Read their owner constraints before implementation. Record material design
contradictions rather than choosing a new product boundary silently.

## Existing behavior and boundaries

Preserve SDPTool discovery/navigation/output and existing manual Session documents;
reuse existing project identity and record validators. Keep ModelGovernance files,
model history and KB050 semantic blueprints with their separately selected work.
Do not make their completion a prerequisite to a plain Maintenance assignment.
KB043's full timeline, all routine families, distribution, release, GUI, remote
hosting and complete shell/file interception remain outside this plan.

## Git and assignment policy

BranchPolicy current means one working branch for the selected implementation,
with milestone commits. At actual start, establish an isolated checkout/branch
`sdp/project-governance-pilot` from the committed PGD1 delivery, or document an
already isolated equivalent; do not switch the concurrent ModelGovernance shared
worktree. Check Git status before creation and staging. Preserve unrelated files.
Follow existing phase-push/PR authority; merge and release need their own explicit
authorization. No push, merge or release is part of this planned handoff.

Use the authorized role/host delegation rules. An independent review must use a
context independent of implementation. If the current host does not support the
required separation or trusted child attribution, retain that integration gap;
do not turn a self-check into an independent review by changing its label.

## Runnable increments

| Phase / milestone | Delivered behavior | Required evidence | State |
| --- | --- | --- | --- |
| PGI1-M1 Core and inspection | Selected project + status request yields context/route without creating work; authorized fixture Maintenance creates a pinned durable run visible through a small CLI read/transition path | PG-A01/02/03/08/10/18; exact JSON schemas; core restart and stale-input tests | planned |
| PGI2-M1 Adapters and protocol | Thin local MCP adapter and controller protocol adapter use the core; inspect compatible schemas, negotiate the selected SDK/host, capture bounded recorded events and reconcile mapped attempts | PG-A10/11/12/15/18; actual protocol negotiation and failure outputs; no simulated-as-live claims | planned |
| PGI3-M1 Owner workflow | Minimal terminal client recovers Session, submits/steers work, presents route/status and launches the authorized fresh Master with bounded context; applicable child/review mapping is observed | PG-A02/04/05/11/12/13/15; managed account mode and one useful bounded task | planned |
| PGI4-M1 Evidence and reconciliation | Candidate-bound evidence, review/rework and checked publication produce one durable Session outcome and correct card/plan disposition | PG-A06/07/09/13/16/17; crash and concurrent-change fixtures | planned |
| PGI5-M1 Integrated verification | The same isolated task survives observer close, interruption/recovery and scope correction with unrelated work preserved | All PG-A01–18 at their stated level; independent review and exact-candidate evidence | planned |

Each milestone is runnable at its stated boundary; no phase claims the whole client
is delivered because a library or mock passes. Core work can proceed with protocol
fixtures while real host compatibility is measured, but PGI3/PGI5 cannot pass with
only schema exports. A real integration trial is required for account and native
child/context claims. Log protocol gaps and adjust the bounded design before relying
on an unsupported field or fallback role mapping.

## Concrete first assignment

PGI1-M1 owns a new bounded governance core package and its CLI adapter under SDPTool,
plus isolated test fixtures and this plan/evidence. Do not refactor the existing
navigation facade wholesale. Read current AGENTS/skills and validate actual Go
availability first; `go` was absent from this session's PATH during preparation.
The module declares Go 1.26.0. Choose a compatible local toolchain through the
project's existing environment conventions; a README claim is not a tested build.

Implement strict application schemas, explicit checkout binding, immutable routine
definition/run identity, read-only routing evaluation, one guarded Maintenance run,
local operation receipts/CAS and restart recovery. Exercise these through the CLI
with a real temporary project fixture. Initially no model submission, account
configuration, existing project-record publication or broad catalog installation.
Do not add a second generic management ledger. Maintain current core documentation
and command help when behavior becomes available.

## Verification and closeout

Run focused state/recovery/adapter tests and existing SDPTool regression tests for
the affected commands, with race checks for shared state. Record unavailable
environment-dependent checks explicitly. Validate project-management and Toolkit
records, preserve original ledger bytes, and verify exact candidate identity.

Independent review assesses intent, real workflow, preservation and authority;
verification distinguishes unit fixtures, protocol probes and real host execution.
Completion requires all selected milestone evidence and explicit remaining-card
dispositions. A passing pilot may leave broader KB036/037/038/042 capability scope
in backlog. No card is closed merely because this plan is complete.

## Current result

Planned handoff only. No milestone has started and no production implementation or
integrated pilot result is claimed. PGD1 compatibility evidence is an input, not
PGI2/PGI3 acceptance.
