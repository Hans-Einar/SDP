# SDP 2.1.0 Sessions release preparation — MaintenancePlan

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0012 |
| project | SDP |
| state | completed |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |

## Outcome and authority

Owner asks whether XFMD can adopt dynamic tabs using the installed producer while
we prepare the Sessions release. Prepare a concrete, verified SDP 2.1.0 candidate;
publication, a new gh-sdp client release and live project upgrades are separate.
Direct owner request selects Maintenance without another wrapper card/Scrum.
Completed KB047 / PLAN-SDP-0015 supply the implementation; do not reopen them.

Latest actual release: SDP 2.0.0. Version 2.1.0 is an additive Sessions capability,
optional inventory.sessions and navigation tab under unchanged sdptool/0.2.
Existing root IDs, language and installed process profiles remain unchanged.
Framework stays 2.0.0: no distributed template/payload changes are required.

## Consumer evidence and release boundary

Read-only installed gh-sdp in XFMD reports engine 2.0.0/d304261. Actual discovery
already returns root IDs whose nodes have kind tab (SDL/KanBan/SDUI); Files is a
directory. XFMD can implement dynamic tabs now by resolving roots through nodes.
No project migration is needed for an existing SDP/Sessions directory. A new
client selecting the future 2.1.0 descriptor is needed for the normal extension
update route. Until publication, installed clients remain on 2.0.0.

## Milestones

| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| SR1 | SR1-M1 | Version, frozen log, descriptor predecessor and consumer handoff prepared | completed |
| SR1 | SR1-M2 | Clean signed candidate, tests/upgrades/archive, CI and independent review evidence | completed |

Use sdp/release-2.1-preparation; commit/push each milestone. Preserve unrelated
files and frozen records. No tag/publication/main merge/global install/XFMD edit.
Verification follows the instantiated ReleaseChecklist. Publication items stay
unchecked at preparation closeout; record exact candidate and package hashes.

## Closeout

SR1-M2 completed on 2026-10-01: [Evidence](Evidence.md), [Checklist](ReleaseChecklist.md),
[Candidate](Candidate.json) and [Independent review](Review.md). Product candidate
93517ad98cd188c0debeb1d0f3d36d123c6e4a3b is signed and verified; subsequent records
commit preserves that exact identity. Preparation complete; publication and client
selection remain a separate future delivery. No active card is held open solely
for that future work.
