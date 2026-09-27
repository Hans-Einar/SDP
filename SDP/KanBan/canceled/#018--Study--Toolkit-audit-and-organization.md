# Audit Toolkit responsibilities, organization and distribution

| Field | Value |
| --- | --- |
| id | KB-SDP-018 |
| project | SDP |
| type | Study |
| CardState | canceled |
| ScrumId | SCRUM-SDP-0002 |
| Systems | SDP |
| created | 2026-09-25T09:09:52Z |
| source | Owner conversation 2026-09-25: SDPTool location and future Toolkit review |
| next_review | None; withdrawn by owner |
| tags | toolkit, maintenance, ownership, distribution |

## Disposition — withdrawn by owner, 2026-09-26

The owner withdrew the separate Toolkit workstream. Root SDPTool is the common
Go command entry point and owns future installation/upgrade execution; gh-sdp is
a thin client. [MAINT-SDP-0006](../../Maintenance/ST1/Plan.md) performs the root
relocation and records retained Toolkit dependencies. [KB-SDP-033](../completed/%23033--Study--XFMD-SDP-adoption-and-SDL-pilot.md)
retains the concrete installation/adoption outcome. This audit was canceled,
not completed; no comprehensive Toolkit audit is claimed. Historical scope below
explains the withdrawal and must not be used to start another Toolkit engine.

## Historical proposed scope (withdrawn)

Perform a thorough review of Toolkit's current organization, supported tools,
contracts and documentation. The new owner-selected `Toolkit/SDPTool/` home must
fit into a coherent Toolkit without creating competing installers, validators,
version sources or process/language responsibilities.

Inventory every top-level area, its purpose, owner, consumers and actual status:
SDPTool, scripts (including shell CLI), schemas, docs, tests/fixtures, conformance,
payload, skills and install inventory. Map authoritative definitions against
distributed copies and generated/frozen evidence. Compare documentation with
actual commands, manifests and consumer dependencies before calling anything legacy.

Identify obsolete/duplicated material, gaps and misplaced responsibilities;
propose keep/move/merge/retire decisions with explicit destinations, redirects,
installation/update consequences and validation. Examine the boundary to Template,
SDL, SDUI, reusable KanBan and project-owned SDP records. Preserve existing
consumers, approved project data, append-only ledgers and historical evidence.

## Related work and boundaries

- [KB-SDP-017](../active/%23017--Proposal--sdptool-and-project-navigation.md) owns sdptool;
  this audit does not block its bounded discovery/viewer implementation.
- [KB-SDP-014](../backlog/%23014--Proposal--KanBan-version-contract-and-distribution.md) owns
  KanBan compatibility/distribution; coordinate its canonical/distributed paths.
- [KB-SDP-011](../completed/%23011--Bug--Traceability-id-conformance.md) owns the known ID/CI
  mismatch; record it as a baseline, not proof that reorganization fixed it.
- [KB-SDP-010](../active/%23010--Proposal--Document-consolidation.md) owns the
  existing editorial consolidation review. Reuse that work; this is a broader
  Toolkit/tool/installation responsibility audit, not a duplicate document cleanup.

## Acceptance and follow-up

Deliver an inventory with source evidence, authority/consumer map, current check
baseline and a prioritized migration plan with phases, milestones and rollback
considerations. Explicitly disposition uncertain or retained legacy material.
Recommend concrete follow-up cards or existing owners for implementation; a
completed Study means the review/plan is delivered, not that files were migrated.
Do not delete or move Toolkit contents merely because this request is registered.

## Worklog

| Time | Actor / event | Outcome | Remaining work |
| --- | --- | --- | --- |
| 2026-09-25T09:09:52Z | Codex; EVT-KB-SDP-000082 | Registered owner-requested audit; checked overlap with #010/#011/#014/#017. | Select the audit later; CardState remains backlog. |

## Scrum-0001 review

Keep Toolkit audit separate from the release contract and current maintenance; it covers ownership/installers/payload/scripts, not only KanBan.

EVT-KB-SDP-000120; next review at the next selection or relevant dependency delivery.

## Scrum-0002 coordination

2026-09-25T12:38:14Z — EVT-KB-SDP-000133: the owner selected future focused Scrums in
[KB-SDP-027](../completed/%23027--Study--Skills-review-and-project-activation.md) and
[KB-SDP-028](../completed/%23028--Study--Installer-upgrade-and-versioned-layout.md).
Reuse their eventual skills/install findings here; do not create competing
migration plans. This card retains the broader Toolkit inventory, ownership,
organization and remaining areas. No audit has been delivered by registration.

## Scrum-0003 focused review available

Reuse [MAINT-SDP-0002](../../Maintenance/SK1/Plan-and-Evidence.md) for the current
skills inventory, draft adoption map, native catalog probe and installer coupling.
The skills Study is complete; adoption remains planned. This broad Toolkit audit
still owns the other responsibilities, organization and distribution questions;
do not repeat the focused review or infer that the whole Toolkit is consolidated.

2026-09-25T14:18:18Z — EVT-KB-SDP-000158: handoff recorded; card remains backlog.

## SK1 completion available

[MAINT-SDP-0002](../../Maintenance/SK1/Plan-and-Evidence.md) now completes skills
consolidation, metadata/installation consumer migration and tested agent activation.
Reuse its current ownership map and evidence. The broad Toolkit audit still owns
other modules, packaging and organization; it need not repeat this delivered work.

2026-09-25T14:48:15Z — EVT-KB-SDP-000162.

2026-09-26T19:03:42Z: EVT-KB-SDP-000187. Owner withdraws separate Toolkit audit; root SDPTool owns future Go installation. Cancel unexecuted audit, retain historical evidence and concrete adoption in KB-SDP-033.
