# R240 — Publish DocumentHost prototype launch

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0023 |
| project | SDP |
| state | completed |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |

Owner explicitly authorizes a new release of the completed PLAN-SDP-0030 delivery.
Publish SDPTool and the identical gh-sdp engine as 2.4.0: a minor additive capability
with explicit context and DocumentHost support; legacy 0.2 without context remains.
SDUI 0.3 check now reports preflight/document/nativeLayout pending rather than
blanket rejection. No domain callbacks, XFMD context picker, main merge, live
extension replacement or consumer Framework migration is included.

Branch sdp/release-2.4.0 starts at 8f1fb9a. Scoped preflight: root has only the
preserved ledger lock; all five detached historical/build worktrees are clean.
Activity0007 and its source plan are completed; other PG/blueprint work does not
own these release paths. Current ledger has no MAINT0023. Extension worker owns
gh-sdp release metadata/packaging only, parent owns canonical metadata/signing and
paired publication. No new product behavior or model edits are planned. Unknown
remote writers remain outside observation coverage. Green for these writes.

| Milestone | State | Acceptance |
| --- | --- | --- |
| R240-M1 | completed | Exact clean signed candidate, changed profile tests, source archive and supported predecessor upgrades, native package check, CI and one combined review |
| R240-M2 | completed | Both public releases downloaded and verified; isolated ordinary extension upgrade and default Framework install pass |

Reuse exact implemented source and native behavior evidence from DocumentLaunch;
qualify packaging and the new signed artifact. Retain original signed predecessors
and add 2.3.2. Preserve release history. One bounded extension worker and one
consolidated release review; no repeated implementation review or broad agent tests.

Both milestones are complete; see Verification.md and IndependentReview.md.
