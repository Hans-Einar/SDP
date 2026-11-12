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
| A1 Creation and records | Activity and Session creation via typed CLI/actions, real manual compatibility, shared canonical history | Domain complete; adapters awaiting integrated qualification. Create/retry/stale/path/identity/crash checks pass; no dispatch on creation |
| A2 Lifecycle | Activate/hold/resume/switch/complete/cancel/reopen with truthful quiescence and single-active rule | Domain complete; concurrency, interrupted/retried switch, quiescence and dispositions pass; host integration qualification remains |
| A3 Navigation | Producer-owned Activities tree, readable landing/phase/history pages and legacy references; XFMD consumption | Native Activities tab, correct ordering and selection, no duplicate Sessions tab or lost legacy documents |
| A4 Timeline | Actual lifecycle intervals and explicitly attributed commit/ledger markers | Unknown time/clock anomaly, dense markers, multiple references, no effort or unrelated-commit attribution |
| A5 Governance and adoption | Activity-scoped execution admission across existing paths; read-only/capture allowed; installed readers/process integration | No out-of-scope/no-active dispatch after adoption; compatibility and pending-run/recovery evidence |
| A6 Qualification | Full integrated candidate, independent review, native workflow, truthful records/docs | All selected design cases disposed, no deferred transition adapters, owner review on concrete result |

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
