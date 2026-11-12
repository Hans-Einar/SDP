# PLAN-SDP-0027 — Activity management

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0027 |
| project | SDP |
| state | active |
| PlanType | ImplementationPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDPTOOL, ProjectGovernance, collaborating XFMD |
| source | Owner accepts Activity0001 and authorizes Activity0002; KB-SDP-061 |

## Outcome and authority

Deliver the selected [mandate](../01-Mandate/Mandate.md) and existing Activity design
through complete runnable workflows. Owner execution authority covers this Activity;
merge, release, live installation and final owner acceptance are separate.
No SDL viewpoint browsing redesign belongs here (KB060 / future Activity0003).

Baseline: accepted Activity0001 code 527cc90, closure record 26d4879. Execute on
sdp/activity-management, based on the accepted restructuring branch without an
implicit merge to main. Commit milestones; preserve unrelated XFMD RPM work.
Use a separate XFMD worktree when product integration requires changes.

## Milestones and acceptance

| Milestone | Outcome | Acceptance / current state |
| --- | --- | --- |
| A0 Executable design | Reconcile canonical schema/readers, transaction/recovery, retained descriptors, scope and producer contracts; validate/checkpoint SDL and pin Blueprint handoff | Completed: pinned concern checkpoints and Blueprint; transaction/recovery and host boundary contract selected |
| A1 Creation and records | Activity and Session creation via typed CLI/actions, real manual compatibility, shared canonical history | Verified with integrated adapters. Create/retry/stale/path/identity/crash checks pass; no dispatch on creation |
| A2 Lifecycle | Activate/hold/resume/switch/complete/cancel/reopen with truthful quiescence and single-active rule | Verified: concurrency, interrupted/retried switch, quiescence, dispositions and host integration |
| A3 Navigation | Producer-owned Activities tree, readable landing/phase/history pages and legacy references; XFMD consumption | Verified in the actual native shell. Producer paging, model-free selection and legacy references retained |
| A4 Timeline | Actual lifecycle intervals and explicitly attributed commit/ledger markers | Implemented; clock/unknown/dense/explicit marker tests pass; colored calendar SVG inspected in native shell |
| A5 Governance and adoption | Activity-scoped execution admission across existing paths; read-only/capture allowed; installed readers/process integration | Implemented; direct/host admission, pending recovery, installed readers and explicit opt-in tests pass |
| A6 Qualification | Full integrated candidate, independent review, native workflow, truthful records/docs | Technical qualification and independent review complete; owner disposition pending on the concrete result |

Milestones may be refined into smaller runnable checks as source inspection resolves
contracts; they may not silently drop selected behavior. The existing Design.md
mandatory-case list is the acceptance denominator. One domain owns lifecycle and
projection; no competing GUI classification or private Activity ledger is allowed.

## Verification and continuity

Use isolated fixture projects for all mutation tests, especially Activity0003.
Keep manual Activity0001/0002 history byte-preserved and explicitly test adoption.
Test crash boundaries/idempotency/concurrency, existing Session/Blueprint/governance
regressions and actual XFMD behavior. Test each material change locally and consolidate
broad suites/independent review at coherent boundaries. Session0018 records work
and handoffs; this plan alone owns detailed milestone status. Final acceptance does
not follow automatically from a passing parser or generated Blueprint.

## Technical delivery disposition

A1–A6 technical obligations are verified on producer eb11ac7 and consumer a558c31.
See ../06-Verification/Verification.md and A6/Review.md. The plan remains active
solely for the actual owner disposition; no incomplete implementation is hidden
behind that boundary. Activity0002 is still the selected active Activity and
Activity0003 has not been allocated. No live adoption, installation or merge.
