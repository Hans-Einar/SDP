# Restore installation conformance in CI

| Field | Value |
| --- | --- |
| id | KB-SDP-030 |
| project | SDP |
| type | Bug |
| CardState | backlog |
| PlanId | MAINT-SDP-0005 |
| Systems | SDP |
| created | 2026-09-25T17:58:59Z |
| source | Merge preparation; GitHub Actions run 36169301165 on PR 34 |
| next_review | MP1-R-M3, before candidate freeze or main integration |

## Observed failures

The [run](https://github.com/Hans-Einar/SDP/actions/runs/36169301165) failed in
both installation jobs. Windows passed Install-SDP.Tests.ps1, then portable
install-v1 conformance reported: empty-default reference outcome differs from
committed authority. Linux ran 19 integration tests but planning_profile_upgrade
could not resolve historical commit 28bf156; the later fault matrix was not run.
These are distinct from KB-SDP-011's 38 existing Traceability findings.

## Selected planning context

[MergePlan](../../Maintenance/MP1/MergePlan.md) selects repair before integration.
Execution is not started. For Windows, compare the structured expected and actual
outcomes to identify behavioral versus authorized inventory/version changes;
never regenerate normative expectations blindly. For Linux, provide the pinned
historical artifact reproducibly in the CI checkout (verified fixture or explicit
history provisioning), with provenance and absence/corruption diagnostics.
Do not remove the upgrade test or substitute the current artifact for the old one.

## Acceptance

The exact integration candidate passes portable v1 conformance on the Windows
runner and profile integration plus the complete interruption matrix on Linux.
Preserve old-format coverage and current distribution facts. Document actual
root causes, corrections and test candidates. Local earlier passes remain dated
evidence; they are not a remote CI pass. No live consumer migration is included.

## Worklog

Registered during MP1 preparation from actual job logs. Card remains backlog
until the authorized MergePlan execution activates its readiness work.
