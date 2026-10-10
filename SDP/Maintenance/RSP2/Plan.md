# Publish runnable-program SDP and paired gh-sdp

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0015 |
| project | SDP |
| state | active |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| source | Owner Session0010 T011 |

## Outcome and authority

Owner explicitly requests a gh-sdp release followed by local gh extension upgrade,
then intends to try SDUI from XFMD. This authorizes the necessary upstream signed
SDP release and paired client default update, publication and extension upgrade.
No XFMD application edits or project SDP-directory migration are selected. Existing
local XFMD source still has the prototype launcher; do not promise UI Run adoption.

Select SDP 2.2.0: additive runnable programs, the accepted SDUI 0.3 widget inventory
and ModelGovernance since 2.1.0. Retain existing language/profile compatibility,
Framework 2.0.0 and Linux amd64 publication scope. New gh-sdp 0.2.2 is a default-
engine patch; client forwarding API is unchanged. Include the small managed SDP
skill update already accepted in main. Preserve all older release sections/records.

Publish exact reviewed branches sdp/release-2.2.0 and sdp/release-0.2.2. A main merge
is not required and is not inferred; PR #53 remains a separately reviewable change.
Canonical Session/management stay synchronized with SDP-vNow; no dirty concurrent
work joins the clean release. gh-sdp uses its installed SPS-009 lifecycle.

## Milestones

| Milestone | Acceptance | State |
| --- | --- | --- |
| RSP2-M1 | Version/notes/inventory, exact signed clean SDP package, predecessor/archive tests, CI and fresh independent review | in-progress |
| RSP2-M2 | Publish verified SDP 2.2.0 and client 0.2.2, verify downloads/default and local gh extension upgrade | planned |

## Gates and evidence

Instantiate Toolkit/docs/ReleaseChecklist.md. Keep publication fields null until
real tag/release exist. Verify every supported predecessor descriptor, including
2.1.0, preserving owner Sessions, program declarations and history with repeat
no-op. Release, client package, production default and native Run evidence must
identify exact candidates. Independent fresh-context reviewers report actual
findings and disposition; no self-review substitution. Publication waits for gates.
