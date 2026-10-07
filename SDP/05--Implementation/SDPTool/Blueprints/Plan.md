# BPI — bounded semantic blueprint producer

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0020 |
| project | SDP |
| state | active |
| PlanType | ImplementationPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDL, SDPTOOL |
| source | KB-SDP-050 / PLAN-SDP-0001 BP2-C |

## Outcome and authority

Implement a Go blueprint producer from captured NOW/TARGET models and an authored
task: semantic differences, affected modeled context, explicit unknowns and preserved
obligations, packaged as deterministic Markdown/Mermaid plus machine data. Include
a retained blueprint catalogue, assignment progress and source-derived discovery
for a consumer Blueprint tab. Owner selected this scope extension on 2026-10-07
in Session0008 T005. Inputs and
ready-state authority must remain distinct from implementation-conformance proof.
This plan is prepared by the authorized BP2 design work, BPI1 execution selected by the owner on 2026-10-07.

Governing design: ../../../04--Design/SDPTool/Blueprints/Contract.md,
Selection-and-Evidence.md and Producer-and-Handoff.md. Use the reduced MVP1 fixtures
and negative cases, not an unsupported full experimental-profile migration.
No production Ponsse/XFMD changes, mandatory Git, grammar expansion, new renderer,
automatic code execution or release packaging. SDUI bytes are retained with explicit
coverage gaps; a semantic SDUI analyzer is not silently bundled into this increment.

Execution uses branch sdp/blueprint-implementation in the isolated worktree
/tmp/sdp-blueprint-implementation based on 3d265d1; do not switch a concurrently used worktree. Per-milestone commits and phase
pushes; no main merge or publication without owner authorization. ModelGovernance
and sourcegraph APIs are prerequisites, not dependencies on their future product release.

## Phases and acceptance

| Phase | Milestone | Complete vertical outcome |
| --- | --- | --- |
| BPI1 Analysis | BPI1-M1 | Go pure analyzer accepts checked 0.6 snapshots/task; reports typed diff, context, coverage and constraints from actual parser facts |
| BPI2 Publication | BPI2-M1 | SDPTool captures artifact inputs and produces deterministic Markdown/Mermaid/JSON bundle through shared output/publisher boundaries |
| BPI2 Catalogue | BPI2-M2 | Retained revisions discovered without registration, with a typed Blueprints tab and verified document targets |
| BPI3 Assignment evidence | BPI3-M1 | Explicit readiness/dispositions, pinned before/after check receipts and independent generated-bundle worker/reviewer trial |
| BPI3 Lifecycle projection | BPI3-M2 | Explicit revision-bound assignment transitions and evidence-backed progress groups, including stale/missing diagnostics |

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


### BPI2-M2 — retained catalogue and discovery

Before coding, finalize the additive discovery contract and retention schema in
Producer-and-Handoff.md and the corresponding SDL design sources. Preserve the
existing discovery contract for other roots. Retain bundles under
SDP/Blueprints/<blueprint-id>/<revision>/, with stable task identity and exact
revision digests. Scratch previews may use an explicit temporary output destination.
Retained revisions are immutable; repeated publication of identical content is
idempotent, while a conflicting existing revision is rejected.

Discover actual bundle metadata and project-management links; do not require a
manual registry, navigation.json, generation during discovery, or status-folder moves.
Expose a root of kind tab, stable node IDs, task/revision children and typed open
targets for generated Markdown/Mermaid. BPI2 can expose the catalogue before
BPI3 lifecycle support, but must label unavailable work status as unknown.
The CLI tree and discovery JSON use the same projection. XFMD application changes
are external; deliver a documented consumer fixture for its agent.

Acceptance: empty/absent catalogue, multiple Systems/tasks/revisions, duplicate IDs,
malformed or unsupported metadata, missing documents, modified retained bytes,
symlinks/path escape and bounded scan behavior. Refresh observes additions/removals
and content edits. Invalid entries remain diagnosable without fabricated targets.
Every valid open target resolves to the intended revision; unrelated roots retain
their behavior. A headless consumer reconstructs the proposed tab from discovery
alone. This proves the producer contract, not native XFMD GUI integration.

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


### BPI3-M2 — assignment lifecycle and grouped navigation

Finalize the transition table, actor/authority requirements and versioned records
before implementation. Use canonical ProjectManagement history for assignment
events and existing Traceability for implementation evidence. Reuse compatible
governance primitives where available; do not invent a second writable ledger or
store mutable workflow status inside generated bundle files. A projection/cache
must be rebuildable. Coordinate schema changes with parallel ProjectGovernance work.

Implement bounded SDPTool operations for readiness disposition, assignment, start,
review and closure, plus hold/resume, cancel and supersede. Final CLI spelling is
part of the milestone design. Bind each assignment to blueprint revision, task,
assignee and applicable plan milestone; a parent card is only a link. Retain
separate assignment identities for multiple attempts or assignees.

Expose draft, ready, assigned, in-progress, review and completed groups, with
on-hold, canceled and superseded dispositions. A task with multiple assignments
must expose them individually rather than guess a single aggregate completion.
Readiness, source freshness, validation and evidence are separate fields.
Completion requires applicable reviewed implementation evidence and a recorded
closure disposition. A parser pass, model RELEASE or completed parent card is
insufficient. New live inputs may stale the comparison, but must not rewrite a
historically completed revision or transfer its approval to a new revision.

Acceptance: full ready -> assigned -> in-progress -> review -> completed flow;
review rejection -> rework; hold/resume and supersession; assignment to stale or
unready revision rejected; wrong assignee/authority or stale expected-state update
rejected; repeated operation handled without duplicate transitions; missing or
ambiguous links reported as unknown. Verify multiple assignments, new revision
after completion, evidence mismatch and interrupted append/recovery using the
selected canonical history mechanism. Discovery is read-only and deterministically
projects groups from the same records; it never infers progress from filenames.

## Verification and closeout

Run targeted Go tests at each milestone, then relevant SDL/SDPTool full suites with
race/vet and compiled CLI workflow on the final candidate. Independent review must
read owner intent before implementation narrative and have all material findings
resolved or explicitly scoped out before closure. Record exact hashes/toolchain,
failed checks and limits. Update source design projections only from validated SDL.
No parser pass substitutes for behavior or cross-project GUI acceptance.

Complete KB050 only after delivered producer, catalogue discovery, assignment
lifecycle, evidence and documentation agree. Include an end-to-end non-Git fixture
that generates two revisions, assigns one, records review/closure and refreshes
the discovery tree without altering the retained bundles.
Session0008 continues through implementation; its roadmap must show design gates
and production milestones separately. Installer manifests, releases and integration
into XFMD require their own authorized work. BPI1-M1 is delivered and independently reviewed; BPI2-M1 is delivered; BPI2-M2 catalogue/discovery is next. See
[Evidence-BPI1.md](Evidence-BPI1.md). Later milestones have not started.

## BPI1-M1 delivery — 2026-10-07

Pure Go analyzer and regression tests delivered; independent review approves the
bounded diagnostic library after four findings were fixed. No public CLI or
assignment readiness is claimed. BPI2-M1 publication is the next milestone.


## BPI2-M1 execution — 2026-10-07/08

Implemented the explicit diagnostic generation command, captured source compilation,
deterministic typed metadata, source-linked Markdown and Mermaid, and shared owned-file
publication. Task format is finalized as strict JSON rather than the earlier proposed
YAML. BPI2-M2 catalogue and BPI3 remain unstarted. Independent review found and then
approved fixes for symlinked-source overlap, explicit compiler/entry identity and
PRESERVE labels for protection IDs containing slashes. See Evidence-BPI2-M1.md. Final integration must resolve the recorded installer
race-test timeout; focused producer verification and independent review passed.
