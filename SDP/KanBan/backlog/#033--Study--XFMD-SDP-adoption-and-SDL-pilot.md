# Plan XFMD SDP adoption through gh-sdp and a practical SDL design pilot

| Field | Value |
| --- | --- |
| id | KB-SDP-033 |
| project | SDP |
| type | Study |
| CardState | backlog |
| Systems | SDP, SDL, SDPTOOL |
| created | 2026-09-25T22:20:56Z |
| source | Owner conversation 2026-09-26: SDL context, implementation drift and post-main XFMD adoption |
| next_review | After MAINT-SDP-0005 main consolidation |

## Owner intent and baseline

After consolidation of SDP-vNow into main, install or upgrade SDP in XFMD through
gh-sdp and begin describing XFMD with SDL. This creates practical experience for
viewpoints and assignment bundles. XFMD already has a bootstrapped SDP area and
KanBan history; inspect its current state rather than treating it as an empty
project. The owner reports no existing SDL model of XFMD; verify before authoring.

## Study and plan

Inspect gh-sdp's actual installation delegation, supported versions and project
recognition. Establish how it selects the intended main/distribution artifact;
merging to main does not publish a Toolkit release. Compare installed manifests,
capabilities and managed/project-owned paths. Prepare a dry-run upgrade with
versioned artifact provenance, backup/recovery and a Maintenance record. Preserve
XFMD's existing project documents and KanBan/ledger history. Coordinate existing
XFMD adoption cards rather than starting a duplicate bootstrap.

Study XFMD code and existing native blueprint documents critically. Describe
observed architecture in SDL with explicit uncertainty and distinguish it from
proposed improvements. Start with one bounded real workflow, then derive useful
views and propose a subsequent assignment-bundle pilot. A native blueprint is
input evidence, not an already accepted SDL blueprint format.

## Acceptance and responsibility

Deliver an executable adoption/pilot plan naming the actual gh-sdp path, artifact,
current/target installation identities, preservation checks and one end-to-end
XFMD workflow to model. Record missing gh-sdp capabilities explicitly. Route XFMD
code/model authoring and gh-sdp changes to their repository agents/cards; do not
perform live migration or product changes during this intake or main merge.

Depends on main consolidation and coordinates with
[KB-SDP-031](../active/%23031--Study--SDL-assignment-bundles-and-blueprints.md),
[KB-SDP-032](%23032--Study--Viewpoint-navigation-feedback.md) and
[KB-SDP-018](%23018--Study--Toolkit-audit-and-organization.md).

## Worklog

2026-09-25T22:20:56Z: Registered in backlog before main integration; EVT-KB-SDP-000174.
