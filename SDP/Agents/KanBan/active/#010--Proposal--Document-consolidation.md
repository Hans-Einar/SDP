# Consolidate documentation and separate active contracts from history

| Field | Value |
| --- | --- |
| id | KB-SDP-010 |
| CardState | in-progress |
| project | SDP |
| type | Proposal |
| created | 2026-09-24T13:54:42Z |
| source | owner-conversation-2026-09-24 |
| next_review | Before the next phase/template-profile delivery in KB-SDP-001 |
| owner | Codex |

## Request and scope

Inventory and consolidate documents in docs, checkpoint #1, SDL and SDUI. Overlapping studies/plans exist, and some older documents describe delivered Go functionality as unimplemented. Establish one active entry point per responsibility, distinguishing language profiles, process, implementation status, proposals and dated evidence.

Checkpoint #1 mainly covers SDL/SDUI but also contains SDP process and shared decisions. Do not move it into SDL and present it as a current language specification. Preserve original evidence/fingerprints and record ownership.

## Relationships and next handling

[KB-SDP-001](%23001--Proposal--Project-structure.md) owns physical structure, path migration and inventory. This card owns editorial consolidation; shared files alone do not justify merging the cards. Start from R1's inventory and consider smaller language/process deliveries with lineage where useful.

## Acceptance

Every source document has an owner and status. Review contradictory claims against current profiles, code and evidence. Superseded documents point to successors or are explicitly archived. Active documentation and generated viewpoints have different authority; tools regenerate generated files. Historical tests and candidates are not evidence for newer implementations.

## Active delivery R2

The [plan](../../../Maintenance/R2/Plan.md) bounds the first pass to entry points, authority and incorrect implementation claims. [Findings and follow-up](../../../Maintenance/R2/Findings.md) identify corrections and evidence. Renumbering/installation profiles and project registry remain in KB-SDP-001. No new language rules are introduced.

## Worklog and revisions

| Time | Actor / event | Work and findings | Evidence / next step |
| --- | --- | --- | --- |
| 2026-09-24T14:55:42Z | Codex; EVT-KB-SDP-000023 | R2-M1: activated after R1. Found obsolete SDL path, delivered G4/G6 still described as planned, and active Python instructions. | [Findings](../../../Maintenance/R2/Findings.md); R2-M2 corrects and checks them. |
| 2026-09-24T15:03:31Z | Codex; EVT-KB-SDP-000025 | R2-M2: addressed identified status conflicts, created SDUI documentation map, replaced runtime proposal as active contract with package contracts. | [Checks](../../../Maintenance/R2/Evidence.md); candidate semantics and SDP process/template profile remain. |

## Next handling after R2

R2's bounded status pass is complete. The card remains active: active and proposed SDL core material still share one large source file, and process proposals need harmonization with an adopted phase/template profile. Next, address that profile with KB-SDP-001 and KB-SDL-001, then consolidate affected documents. R2 adopts no new SDL keywords, A0–A5 levels or installation names.

## R3 activation

2026-09-24T17:55:07Z: EVT-KB-SDP-000061, ready → in-progress.
[Plan](../../../Maintenance/R3/Plan.md) bounds the remaining work and review result.
