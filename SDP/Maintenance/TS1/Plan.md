# TS1 — Current project templates and SDL source guidance

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0009 |
| project | SDP |
| state | completed |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | phase |

## Outcome and authority

Owner request of 2026-09-28: clarify how agents organize SDP/SDL and overhaul the
installation template documentation. Direct Maintenance selection; no wrapper card
or Scrum. Follow the MVP1 pilot ownership boundaries, without migrating existing
models or implementing experimental language features. KB-SDP-020 retains that
separate migration. No release, consumer upgrade or XFMD changes are selected.

## Phase TS1

- TS1-M1: document the System-owned source convention, agent discovery, numbered
  document roles, source/output authority and migration limits; update the actual
  five-phase seeds and clearly identify retained legacy templates.
- TS1-M2: include the source guide in the release inventory; rebuild the derived
  legacy profile artifact; verify a development Go descriptor and fresh install,
  project-owned preservation, current profile and management consistency.

Use current branch sdp/mvp1-source-ui, one TS1 phase commit after verification.
Git metadata is read-only under the current session permissions; leave commit/push
pending rather than bypassing those permissions. Preserve unrelated sourceinput,
node_modules and package files. Evidence must identify the uncommitted candidate.

## Acceptance and remaining work

Installed guidance must work without links back into the SDP development repo.
An agent can select a System/container source home and distinguish shared libraries,
UI screens, authored prose and generated views. The inventory distributes both
SDL README and AGENTS as initialize-if-missing project files. Existing project
content is preserved. Current published releases remain unchanged; distribution
requires a later release and existing prose may require explicit reconciliation.

See [verification and handoff](Evidence.md). TS1-M1 and TS1-M2 are verified in the
working tree; the required phase commit remains pending because Git metadata is
read-only. This Maintenance
does not close KB-SDP-020 or assert full MVP1 parsing/navigation support.

## TS1-M3 — canonical template root (owner follow-up)

Owner selects relocation of the current five-phase seeds to Template/sdp-root and
removal of Template/profiles. Preserve the old sdp-root and project-root together
under Template/legacy/install-v1; retain current copies of the four still-used
neutral seeds. Update both inventories and live conformance source paths, keeping
installed destinations, archived payload bytes and published releases unchanged.
Verify current/legacy inventory consistency and identical Go descriptor payloads
before/after this authoring-path move. This extends TS1 on the same branch and
phase commit policy; existing source-model migration remains outside scope.

TS1-M3 is verified in the working tree; see Evidence.md for exact payload/archive
preservation and the limits of legacy verification. Git closeout remains pending.

[Release handoff](Release-handoff.md) records owner-requested next delivery, verified
XFMD receipt, predecessor digests and the exact session-access blocker.

## Closeout

2026-09-28: owner enabled unrestricted session access and authorized integration,
release and XFMD upgrade. TS1 source/template work is complete and is recorded in
the TS1 phase commit. Earlier access-blocker notes describe the previous session
state. Release execution follows a separate RP3 plan.
