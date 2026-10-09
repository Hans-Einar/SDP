# PLR1 — Separate project leadership, Steering and SAD roles

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0015 |
| project | SDP |
| state | completed |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |

## Authority and scope

Owner Session0011 T004 requests Project Leader as the normal owner-facing agent,
separate strategic Steering, SAD design work and task Masters. Deliver a local
skill/procedure refinement and a concrete handoff to the active ProjectGovernance
workstream. Reuse Architect for the SAD assignment; do not duplicate architecture
policy in another skill. Skills are guidance, not implemented host capabilities.

Current branch sdp/blueprint-implementation; commit PLR1-M1 after validation.
No app-server/MCP implementation, product dispatch edits, agent launch, main merge,
release or online skill installation is selected. Update canonical skill inventory
and distribution metadata together; do not call the changed payload released 2.1.0.
Next release identity and signed assets remain separate release work.

## Scoped concurrency disposition

The T002 red gate still applies to KB051 product code. For this explicit maintenance,
inspection found no committed divergence or dirty changes in Skills/, skill discovery
symlinks or skill distribution inventory in the governance/program/primary worktrees.
Therefore the isolated role/skill documentation write set may proceed. Leave the
other workstream's active Design.md, implementation plan and source unchanged; add
one handoff document here. Shared ledger additions remain append-only and require
integration reconciliation. No exclusive reservation service is claimed.

## Phase PLR1

| Milestone | State | Acceptance |
| --- | --- | --- |
| PLR1-M1 | completed | Distinct portable roles, new Project Leader skill, SAD handoff, metadata/discovery consistency, scoped runtime-gap handoff and honest verification |

## Verification and remaining work

Five touched/new skills pass skill-creator quick_validate; both canonical/distribution
skill tests pass. Toolkit validation, management replay, maintained-link checks and
git diff --check pass. Initial validation caught an outdated installed-manifest
example and a draft management link to an unregistered Session record type; both
were corrected before commit, preserving all pre-existing ledger bytes. Runtime implementation and actual independent behavioral/online-host
validation are outside this maintenance. Existing governance adapter restrictions
remain in force. Direct conversation with SAD/Master is recorded without making
local chat history an alternative authoritative scope or approval record.


## Same-context boundary review

Inspected the skills against these scenarios: a remote Steering assistant has only
committed evidence; a leader has no launch capability; SAD changes a reviewed TARGET;
a Master phase completes while owner input changes scope; host disconnect loses
observation; existing multi-phase authorization permits continuation. Guidance retains
unknowns, supported-tool boundaries, revision invalidation, durable feedback, explicit
freshness and existing authorization respectively. This is a same-context instruction
review, not a live independent agent evaluation or runtime enforcement evidence.

## Delivered records

- New Skills/sdp-project-leader/SKILL.md and shared Skills/sdp/references/roles.md.
- Existing Steering, Architect/SAD, Master and router clarified/versioned; canonical
  and install payload inventories plus discovery link aligned.
- Project-Leader-Handoff.md captures owner direction and observed runtime limitations
  for Session0007 / KB038 without editing that branch's active contract/source.
- Session0011 T004 records the scope change and preserves the red product-code gate.

No online host/plugin installation, main merge, release, runtime UI or privileged
MCP dispatch is delivered. Next: ProjectGovernance owner reconciles the handoff into
its active design/plan and verifies a real multi-assignment workflow. Our KB051 source
integration/identity conflict remains unresolved.
