# RL1 — Sessions distribution and repeatable release preparation

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0011 |
| project | SDP |
| state | completed |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDPTOOL |

## Outcome and authority

Owner requests an automatic per-release log, Sessions in the next distribution,
and a ReleaseChecklist that verifies installation manifests before manual XFMD
upgrade. Execute this bounded Maintenance on sdp/release/sessions-and-release-log,
stacked on a0aa805. KB-SDP-044 owns delivery; KB-SDP-042 supplies the Session pilot.
Use Go for product tooling. Preserve unrelated untracked work and all project files.
Do not perform the owner's XFMD upgrade. Merge and publication remain separate
explicit decisions after concrete preparation; prior RP3 authorization is historical.

## Version and scope

SDL design-core 0.6 is a language profile, not an SDP product release. Actual latest
SDP is 1.0.0 (verified on GitHub). Recommend SDP 1.1.0 for additive source-composition,
Session guidance and release-log command; preserve 1.0.0 published history. This
work prepares capabilities and reviewed notes; production signing/tagging and the
client default are recorded as pending until publication is selected.

Adopt manual Sessions and distribute a guide/template, not project conversations.
No new Session ledger kind, automatic transcript capture or event timeline engine.
Generate per-release Markdown from canonical RELEASE-NOTES.md; do not infer accepted
behavior from Git commit subjects. Provide deterministic check mode for CI/gates.

## Milestones

| Milestone | Acceptance | State |
| --- | --- | --- |
| RL1-M1 | Session templates/guidance, payload inventory/capability, generated release logs and checklist, tests | completed |
| RL1-M2 | Real 1.0.0 descriptor upgrade rehearsal, preservation/no-op, candidate docs/evidence and manual handoff | completed |

## Verification and remaining gates

Run Go product/inventory tests and CLI tests, installer apply against temporary
projects using verified predecessor bytes, project/Toolkit/schema/link checks.
Check Session files and custom content preservation, repeat no-op and immutable
released notes. No required native GUI test: this changes process distribution.
Production signed assets, exact clean candidate CI/review, merge/release approval,
client selection and downloaded-asset acceptance must precede a real gh-sdp upgrade.
Record actual results in Evidence.md and ReleaseChecklist.md. Do not call the
release available merely because the local upgrade rehearsal passes.

## Delivered preparation

Both milestones are complete with [evidence](Evidence.md), independent bounded
review and [ReleaseChecklist](ReleaseChecklist.md). The proposed release remains
unpublished: freeze/version records, production signing/exact-candidate gates and
publication require the next selected release operation. The owner will run the
actual XFMD upgrade; [handoff](Manual-upgrade.md) distinguishes extension update
from project upgrade. This completed plan is not reopened for publication.
