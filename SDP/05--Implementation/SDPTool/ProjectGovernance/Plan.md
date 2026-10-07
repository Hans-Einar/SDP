# PGI — One governed Maintenance workflow

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0018 |
| project | SDP |
| state | active |
| PlanType | ImplementationPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDPTOOL |
| source | PLAN-SDP-0017; KB-SDP-038 with KB036/037/042 contributions |

## Outcome and authority

Deliver one useful owner-to-agent Maintenance workflow through a local Codex client,
shared SDPTool process core and agent adapter. Restore the same Session after an
interrupted attempt and return candidate-bound evidence with a truthful disposition.
The owner continuation on 2026-10-06 selects PGI1-M1 from this handoff. The plan
is active. The owner continuation on 2026-10-07 selects PGI2-M1; PGI3 is selected by the next owner continuation; PGI4–PGI5 remain planned.

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
| PGI1-M1 Core and inspection | Selected project + status request yields context/route without creating work; authorized fixture Maintenance creates a pinned durable run visible through a small CLI read/transition path | PG-A01/02/03/08/10/18; exact JSON schemas; core restart and stale-input tests | delivered; bounded core/CLI verified |
| PGI2-M1 Adapters and protocol | Thin local MCP adapter and controller protocol adapter use the core; inspect compatible schemas, negotiate the selected SDK/host, capture bounded recorded events and reconcile mapped attempts | PG-A10/11/12/15/18; actual protocol negotiation and failure outputs; no simulated-as-live claims | delivered; bounded adapters and real idle protocol probe verified |
| PGI3-M1 Owner workflow | Minimal terminal client recovers Session, submits/steers work, presents route/status and launches the authorized fresh Master with bounded context; applicable child/review mapping is observed | PG-A02/04/05/11/12/13/15; managed account mode and one useful bounded task | in-progress; live workflow evidence required |
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

PGI1-M1 is delivered in an isolated local clone based on 043c59c, which includes
the completed ModelGovernance delivery. [Evidence-PGI1](Evidence-PGI1.md) and its
[candidate manifest](candidate-PGI1.json) identify the tested core/CLI subset.
Go 1.27.1 full SDPTool race tests, vet and 21 observed wire-schema cases pass.
PGI2-M1 is also delivered: [adapter evidence](Evidence-PGI2.md),
[candidate manifest](candidate-PGI2.json) and [protocol pin](protocol-PGI2.json)
record shared-core MCP parity, bounded observation/recovery fixtures and actual
Codex 0.160.1 / MCP 2025-11-25 negotiation with SDK v1.8.0. No model turn is sent.
PGI3 is in progress: its terminal controller candidate has fixture and compiled
command coverage, but the real host currently exits before initialization.
[Checkpoint evidence](Evidence-PGI3-C1.md) records the exact limits. PGI4/PGI5 remain
planned; host execution, publication and independent review are still required.
KB038 remains active/in-progress for the unfinished PGI3 milestone.

The isolated clone on `sdp/project-governance-pilot` retains milestone commits.
Current isolation preserves concurrent shared work; the earlier host's read-only
Git restriction is historical. Verified changes return to the shared worktree
without switching its branch or including unrelated work in the milestone commit.


## PGI1 implementation refinements and next handoff

The state directory is located by canonical-root SHA-256 and contains the separate
checkout UUID binding. This makes read-only lookup possible before explicit bind;
no discovery read registers a checkout. Selected canonical management records are
checked without publishing them; the full ledger hash conservatively invalidates
advancement after any append. Fine-grained authority/event reconciliation belongs
to the adapter/evidence phases. The pilot ships one compiled routine definition,
not the complete method catalog. Requests and results use the exact PGI1 schema.

PGI2 should start from this delivered core, preserve CLI behavior, inspect the
current installed protocol/SDK and implement thin adapters with equivalent core
results and measured negotiation. Do not promote local Caller capability from
agent-provided role strings. Resume here proves core persistence only; PGI2/PGI3
must establish actual host-attempt mapping before claiming controller recovery.
Reviewed backlog dependencies KB036/037/042/043 remain separate and unchanged;
ModelGovernance, semantic blueprints and full KanBan TUI are not entry prerequisites.

## PGI2 selection — 2026-10-07

EVT-PM-SDP-000145: owner resumes Session0007. Reuse the clean isolated PGI1 clone on
`sdp/project-governance-pilot`. Original read-only Git limitations were historical;
continued isolation now preserves concurrent shared work. Implement thin adapters,
scoped agent authority and bounded observed attempt recovery; no model execution
is needed for protocol negotiation. Live owner/agent workflow remains PGI3.


## PGI2 delivery and next handoff — 2026-10-07

The live probe observes managed ChatGPT mode, initialization and MCP discovery in
an idle ephemeral thread. The core stores observed thread/attempt/epoch identities
and paged completed content; a disconnect requires explicit reconciliation.
Fixtures establish replay/authority guards, not recovery of an executing model.
The next bounded assignment is PGI3: the terminal owner workflow, durable outbound
intent and controller loop, fresh Master launch and actual host role mapping.
Keep PGI4 publication and PGI5 independent integrated verification distinct.
Backlog review leaves KB036/037/042/043 unchanged; no new prerequisite from
ModelGovernance, semantic blueprints or KanBan TUI is introduced.

## PGI3 selection — 2026-10-07

EVT-PM-SDP-000148: owner continues. Work from committed PGI2 `12542b2` in the
isolated clone. Current writable roots include /tmp and SDL only; shared SDPTool
and Session updates cannot be returned in this host. Keep the Session, card, plan
and evidence together in the isolated checkout and report the publication gap.


## PGI3 implementation checkpoint — 2026-10-07

PGI3-M1-C1 adds durable owner-input dispatch, a fresh assignment-bound Master
launch, managed-mode/MCP checks, line-oriented steering/interruption/scope stop,
read-only status and conservative mapped-history reconciliation. The complete
isolated race suite, vet, application schemas and captured RPC schemas pass.
The compiled workflow uses a recorded fake host, never a real model.

The actual rebuilt idle probe exits with EOF before initialization. PGI3 remains
in progress; no real model task, running-turn reconnection or native role evidence
is claimed. Deliberate resolution of an unmapped ambiguous send remains unavailable.
Current host rules prohibit delegation, so the candidate does not grant native
children or pretend to measure a hierarchy. The next work is to restore a permitted
host connection and establish the required live evidence, then resolve remaining
recovery gaps. PGI4/PGI5 acceptance is unchanged. The shared checkout is read-only
outside SDL; this checkpoint and its Session records are retained only in the clone.
