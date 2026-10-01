# Publish SDP 2.1.0 and paired gh-sdp — MaintenancePlan

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0013 |
| project | SDP |
| state | active |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |

## Outcome and authority

Owner explicitly authorized publication on 2026-10-01 after the completed
[SR1 preparation](../SR1/Plan.md). Publish the exact reviewed SDP 2.1.0 candidate
93517ad98cd188c0debeb1d0f3d36d123c6e4a3b and a thin gh-sdp patch selecting its
immutable bootstrap. No main merge, global extension installation or live project
upgrade is selected. The owner will test the extension update themselves.

No new notifier implementation: GitHub CLI checks extension updates once daily
when an extension runs and emits the notice on stderr. A new gh-sdp client release
is necessary because gh does not track the separately pinned SDPTool release.
Source: https://cli.github.com/manual/gh_help_environment

## Milestones

| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| SR2 | SR2-M1 | Revalidate exact prepared bytes; publish annotated SDP tag/assets and verify downloads | in-progress |
| SR2 | SR2-M2 | Publish independently reviewed pinned client; reconcile actual identities and owner handoff | planned |

Root records use sdp/release-2.1-preparation with milestone commits/pushes. gh-sdp
owns its bounded release records and separate branch. Preserve unrelated files.
The clean prepared product checkout/package remain the publication source; later
records commits do not change the package's source identity. Existing CI/review
and original predecessor/archive evidence in SR1 remain applicable to those bytes.

## Verification

Use [ReleaseChecklist](ReleaseChecklist.md), [SR1 candidate](../SR1/Candidate.json)
and [independent review](../SR1/Review.md). Check tag absence and clean exact source
before tag creation. Compare downloaded assets to prepared SHA256SUMS, bootstrap
with production trust and isolated cache, and inspect real GitHub identities.
Paired client requires its own race/vet/package/independent review and remote
immutable dependency checks. Mark publication facts only after they exist.
