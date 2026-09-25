# Scrum-0002 — SDPTool delivery sprint

| Field | Value |
| --- | --- |
| id | SCRUM-SDP-0002 |
| project | SDP |
| state | completed |
| source | Owner conversation 2026-09-25; KB-SDP-017 and its phase plan |
| participants | Project Owner via conversation; Codex review and recording |
| outcome | SPR-SDP-0001 planned; KB-SDP-027 and KB-SDP-028 registered |

## Actual review and decisions

The owner asks for a Scrum that decomposes KB-SDP-017 into a Sprint, plus separate
backlog registrations for skills and installer Maintenance Scrums. Reviewed the
source card, pipeline guide, phase plan, shared management workflow, source-set
queue and related Toolkit/compatibility cards. This is the actual conversational
review, not a fictional meeting or owner acceptance of unimplemented behavior.

Select [SPR-SDP-0001](../../Sprints/Sprint--%230001--SDPTool-preview-and-navigation.md) as a bounded producer delivery:
P0-M1, T1, T2, T3 and T4-M1. Six new cards KB-SDP-021–026 have concrete acceptance
and dependencies. A partial split (KBO-SDP-000004) preserves provenance from
[KB-SDP-017](../../KanBan/active/%23017--Proposal--sdptool-and-project-navigation.md); no scope is discarded.
P0-M2, T4-M2 and T5 remain explicitly with that active/ready parent, outside Sprint
membership. Full native XFMD work stays in XFMD. The saved-file operation is
first because it gives a usable result without waiting for full discovery or new
SDL source syntax; the remaining cards build a coherent project navigation path.

## Dispositions

| Record | Outcome |
| --- | --- |
| KB-SDP-021 | queued; first proposed Sprint work after MAINT-SDP-0001 |
| KB-SDP-022–026 | backlog with SprintId; ordered contracts, implementation and review |
| KB-SDP-017 | active/ready retained scope; six producer milestones transferred |
| KB-SDL-005 | backlog, replacing the previous queue priority; still prerequisite to KB-SDP-020 |
| KB-SDP-020 | backlog, unchanged source migration dependency |
| KB-SDP-010 | existing non-blocking owner gate-review unchanged |
| KB-SDP-018 | broad Toolkit audit retained; reuse future focused reviews below |
| KB-SDP-027 | future Scrum → Maintenance: root Skills/ and verified host activation |
| KB-SDP-028 | future Scrum → Maintenance: declarative process profile and safe upgrades |

The two Maintenance requests are not members of the SDPTool Sprint. Their
ScrumId denotes registration here, not execution of their future topic Scrums.
No extra Maintenance record is invented before that work is selected.

## Findings retained for later Maintenance

Installation is already configuration-driven: Toolkit/SDP-install.manifest.json
owns inventory/policy; SDP.manifest.yaml owns release/capability facts; project
manifests point to generated installed Toolkit facts. The future upgrade work
must reconcile the newly adopted local process profile with distributed templates,
not add a competing version registry. Manual XFMD adoption needs an explicit
baseline and preservation path. The requested old/new-version Maintenance entry
is captured in #028, including failure and unknown-baseline handling.

Toolkit/skills is versioned inventory; Toolkit/skills_v2 is an explicitly unadopted
draft. #027 reviews both toward root Skills/, with actual discovery/loading tests.
A root SKILL.md as activation mechanism remains a question, not a proven host rule.

## Next action and record limits

The Sprint is planned, not running. Start it by activating its six cards together
and recording start; then execute #021. The owner requested planning in this
turn, so this Scrum does not claim product implementation. No Traceability event
is appended for this management-only work. Existing uncommitted sourceinput drafts
remain outside this phase and outside implementation evidence.

## PM2-M1 verification

Branch sdp/phase-pm2-sdptool-sprint-planning from 8d85b39; documentation and
management changes only. Checks passed:

- Project-management validator: 35 cards, four management records, three lineage
  operations and 178 events; schemas, chains, metadata, placement and membership.
- Four management test groups and all 15 lineage negative cases.
- Sprint CLI filtering lists exactly #021–026, with only #021 queued.
- Documentation check: 105 frozen records/prefixes and 574 generated artifacts
  preserved; 2,377 local file links and 130 fragments resolve.
- Prior management-ledger bytes are an exact prefix; Traceability is unchanged.
- git diff --check passed. No product/runtime, native GUI, installer execution
  or agent skill activation is claimed by these checks.

The phase is ready for its milestone commit/push and combined review. Backlog
review found no held work to select; source-set migration, broad Toolkit audit,
compatibility and the two future Maintenance Scrums remain separately owned.
