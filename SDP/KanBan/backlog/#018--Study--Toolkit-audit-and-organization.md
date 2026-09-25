# Audit Toolkit responsibilities, organization and distribution

| Field | Value |
| --- | --- |
| id | KB-SDP-018 |
| project | SDP |
| type | Study |
| CardState | backlog |
| ScrumId | SCRUM-SDP-0001 |
| Systems | SDP |
| created | 2026-09-25T09:09:52Z |
| source | Owner conversation 2026-09-25: SDPTool location and future Toolkit review |
| next_review | Before broader Toolkit restructuring or distribution changes |
| tags | toolkit, maintenance, ownership, distribution |

## Need and scope

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
- [KB-SDP-014](%23014--Proposal--KanBan-version-contract-and-distribution.md) owns
  KanBan compatibility/distribution; coordinate its canonical/distributed paths.
- [KB-SDP-011](%23011--Bug--Traceability-id-conformance.md) owns the known ID/CI
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
