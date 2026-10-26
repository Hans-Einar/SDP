# PGF1-C2c3 — Explicit remaining-scope transfer

T048 implements a bounded successor transition after every predecessor step has a
committed disposition. It closes the predecessor's selected work while preserving
remaining card work in one explicitly identified manual successor. Active rebind,
in-flight transfer, autonomous host scheduling and integrated model gates remain open.

## Behavior and authority

`session-successor` and `governance.session.successor` inspect the exact successor
binding, authority, roadmap anchors, planned plan and remaining cards. The request
names the predecessor instance. Session-scoped MCP exposes `session_successor` under
the same launcher-selected instance and pinned Project Leader skill.

`session-advance` with `action: transfer` requires the exact disposition and revisions,
an empty step ID and the returned successor binding path/input revision. Preparation
pins source preimages; resume publishes the checked effects. Remaining cards become
ready with the successor PlanId; the predecessor plan/Session close, the successor
plan stays planned, and original result/remaining-work histories remain. Explicit
canonical links connect both plans and the transition retains the successor context.
A new adoption/selection/assessment is mandatory before successor execution. No old
review, owner decision or evidence is inherited as permission to implement.

Adoption now shares the canonical writer reservation. Even after all transfer files
have their postimages, a pending operational commit cannot authorize adoption. The
same operation receipt is reusable across CLI and MCP without duplicate effects.
Profile 0.21 prevents older readers from treating transfer as legacy state; the actual
predecessor binary refuses the new `successor` field. Existing supported profiles
remain readable and later operations do not downgrade 0.21.

## Candidate and checks

[Candidate](evidence-T048/transfer-candidate.json) pins baseline plus every changed
product file and the actual compiled tools. The [scoped task](evidence-T048/transfer-blueprint-task.json)
passes diagnostic Blueprint constraints against unchanged unreleased WORK :00002.
It remains a preliminary observed-baseline comparison, not model acceptance/release.
The adapter scope was extended before edits to include JSON invocation and its
catalogue assertion; canonical model responsibilities and host authority are unchanged.

Tests exercise transfer, preserved remaining scope, fresh successor adoption and
assessment, stale/extra/missing inputs, wrong plan/action/caller, profile downgrade,
forged retained context, unfinished steps, no remaining scope, every canonical effect
boundary and reservation/HEAD interruptions. A changed prepared target is preserved
and refused. The pending-adoption test reaches complete canonical postimages while
the transfer reservation still exists and proves adoption cannot race that boundary.
Compiled CLI, JSON invocation and MCP exercise the same context, receipt and resume;
observed wire objects pass the maintained schema validator.

Final per-package and race/vet outcomes are recorded in the candidate. An early full
run overlapped source changes while building subprocess tools and produced incompatible
postimages; it is not accepted evidence. The frozen-product full run exposed the old
catalogue test expectation, which was updated for the added read-only action and
rerun for the entire affected root package. This correction does not remove a behavior
obligation or relabel that initial invocation as passing.

## pg_lab

Lab commit e63188f preserves the intermediate and final compiled fixtures. The final
fixture is `session-pgf1-c2c3-transfer-final`, using external private state under
`/home/warloc/.cache/sdp-workstation/T048/lab-transfer-final-state`. The successor stays
unadopted in this durable fixture. Unit/service verification separately proves its fresh
adoption and assessment boundary.

[Result](evidence-T048/transfer-lab-result.json) and
[canonical/preservation checks](evidence-T048/transfer-lab-validation.json) retain the
observed outcome: one card, two plans, 20 canonical events; all 545 previously tracked
lab paths unchanged; the earlier reader refuses unsupported transfer state. Reports,
tests and review inputs in this particular fixture are synthetic. Actual native work,
question/reply and fresh native review remain the separately pinned earlier PGF2/4/7
results. This is not independent PG source review or human acceptance.

## Remaining boundaries

The old binding is never rewritten. Legitimate selected-source edits still require the
unimplemented reviewed rebind/replan path. Transfer does not abandon unfinished steps,
split one transfer across multiple successors, move an active host, grant design/model
release authority or provide distributed allocation. Concurrent-work preflight remains
required across worktrees. All broader PGF obligations in the completion map remain;
PLAN-SDP-0018 and KB038/052 stay active.
