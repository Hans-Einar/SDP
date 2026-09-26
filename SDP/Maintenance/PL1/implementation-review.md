# Independent PL1 implementation review

Final disposition: **approved for the reviewed PL1 scope; no remaining material findings identified**. All three concrete findings below were corrected and independently retested. This is independent review evidence, not owner acceptance or release approval.

The chronological findings below record earlier candidate states; their final disposition is resolved. The initial review began at the working diff after c2d8edf.

## Scope and authority
Read sdp-reviewer, KB-SDP-029, MAINT-SDP-0004, the planning contract/template, and the relevant implementation diff plus surrounding callers. Reviewed management payload 0.2, type/Sprint semantics, installer profile transition, installed readers, portable skill/reference delivery and Git-policy adoption. No repository edits or agent delegation. The unrelated SDL/go/sourceinput draft was not read. Source hashes identify the reviewed working files; the parent is actively updating this candidate.

## Original consumer-validation finding — resolved

### P2 — Consumer installer accepts a prohibited change to an existing plan's type

Evidence location: Toolkit/scripts/Process-Install.ps1, Assert-ProcessHistory (approximately lines 146–179). The new 0.2 dispatch chooses the right payload schema, but after per-event schema validation the function only checks previousEventId/from/fromPath. It never compares planType with the preceding event. The local validator explicitly rejects this mutation in SDP/ProjectManagement/validate.py.

Reproduction: the isolated consumer-invalid-type fixture contains a valid profile 0.2 board and two linked management events for PLAN-SDP-0001: created/planned ImplementationPlan, then updated/planned DesignPlan. Both values are individually legal enum values and the predecessor chain is correct. The real Install-SDP.ps1 -ProfileArtifact … -PlanJson invocation returned canApply=true with no conflicts. The plan document agrees with the changed type, so this is a history-invariant failure, not malformed Markdown.

Contract violated: adopted project-management profile 0.2 states that typed plan identity/type cannot change during its lifecycle. PL1 explicitly includes installed workflow and validation compatibility; a consumer can currently obtain an applicable installation plan over typed history that the producer's validator rejects.

Consequence: schema-valid but semantically invalid typed histories pass the distribution's consumer check. Merely advertising/accepting management profile 0.2 does not establish equivalent typed-plan validation.

Remediation direction: enforce immutable typed-plan identity/type in the generic PowerShell history checker and add an installed negative fixture for this case. Review the other newly introduced 0.2 invariants in the same path (plan references, policy/current metadata and Sprint membership) before making broader validation claims. Keep any evidence claim scoped to what the installed checker actually verifies.

Artifacts: consumer-invalid-type/SDP/ProjectManagement/Ledger.ndjson; consumer-invalid-type/SDP/Plan/PLAN-SDP-0001.md; consumer-invalid-type-plan.json. The latter is the real plan output, with canApply=true, conflicts=[], baseline=local-five-phase. No plan was applied.

## Resolved during review

### P2 — Deferred plan initially prevented truthful Sprint closeout

The original contract allowed a deferred active plan with an explicit pause, but the validator rejected a completed Sprint retaining that plan, and the schema required nonempty membership on every event. A direct fictional history was rejected with `closed sprint has unfinished plans` (deferred-sprint-events.json and deferred-sprint-result.txt).

The parent clarified explicit removal-before-close, preserving the deferred plan's identity and history, and changed the schema to permit empty snapshots after selection. Creation/start still require work. I independently verified the revised path using the full validate_current function, including an active plan without SprintId, an explicit sole-plan removal event, and a completed Sprint with empty Members/Plans. Result: PASS (0 cards, 2 management records, 0 lineage operations). Evidence: deferred-sprint-fixed/events.json and result.txt. This finding is closed for those updated files.

## Positive evidence and limits

- Six management unit tests passed before the closeout update; their earlier coverage did not catch either finding. The full custom current-record check above verifies the added closeout route.
- Both skill metadata/portable-reference tests passed. The new skill and both references are present in install-v1 inventory and the profile artifact inventory; internal skill links resolve in a copied installation.
- The planning skill behavioral trial separately produced four correct reviewable outcomes; it does not prove installer validation.
- Inspection confirms explicit old/new management-profile reader dispatch, old 0.1 payload schema retention, an explicit board-profile adoption branch, old project-owned text preservation, and managed planning guidance delivery. The parent-reported upgrade integration test was still running; this review does not assume its outcome.
- No general release, live consumer rollout, merge, all-platform compatibility, or owner acceptance is established by this review. The temporary fixture received only read-only PlanJson evaluation.

## Follow-up review after consumer history hardening

The original type-mutation fixture was rerun through the actual installer CLI. It now returns canApply=false with `invalid-management-history: Plan type changed or removed`. The immutable-type finding is resolved for this implementation; see consumer-invalid-type-plan-fixed.json. Inspection also confirms new 0.2 lifecycle/kind, metadata/Git policy, Sprint references/state/reciprocity, and PlanId checks.

A new P2 regression remains in Get-PlanningMetadata: its regex skips standard empty Markdown rows such as `| Plans | |` and `| Members | |`, because the expression consumes both a leading and trailing space around an empty value. Only the nonstandard doubled-space spelling used in the new integration fixture matches. The existing local parser accepts the usual spelling. This affects ordinary plan-only Sprints and empty Sprints after all work is explicitly deferred.

Reproduction: copied the independently valid deferred-sprint-fixed current records to a modern consumer fixture and invoked the actual installer -PlanJson. It returned canApply=false with `The property 'Plans' cannot be found on this object`. See consumer-deferred-sprint/ and consumer-deferred-sprint-plan.json. Correct the metadata parser to accept ordinary empty cell whitespace, and retain a positive installed test for this path.

Observed logs (read, not assumed): /tmp/pl1-invalid-history.txt reports the new negative integration test passed; /tmp/pl1-integration.txt reports upgrade and installed consumer workflow tests passed before this hardening. Neither proves the empty-cell case above. No installer plan was applied by this review.

## Final correction verification

Reran the exact consumer-deferred-sprint fixture through the actual installer CLI after the whitespace parser correction. Result: canApply=true, conflicts=[]. The standard empty `| Plans | |` and `| Members | |` rows now work. Evidence: consumer-deferred-sprint-plan-fixed.json. Inspection confirms whitespace handling is restricted to each row, without crossing line boundaries, and empty Plans/Members access is safe. This resolves the empty-metadata regression.

Final scope: typed plan purpose/identity and policy adoption, realistic skill routing, portable references, local lifecycle validation, explicit profile compatibility/adoption handling, and the inspected consumer-history checks. No remaining material defect was identified in this scoped review. The original type-mutation fixture remains documented as correctly rejected after hardening; the subsequent parser change only affects metadata extraction and does not remove that invariant. Broad platform/release assertions remain outside this evidence. The parent's final complete validation run remains its responsibility; this review does not substitute for that run or for owner acceptance.
