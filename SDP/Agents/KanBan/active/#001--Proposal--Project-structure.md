# Project structure, Template and studies per phase

| Field | Value |
| --- | --- |
| id | KB-SDP-001 |
| CardState | ready |
| project | SDP |
| type | Proposal |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | At the next consolidation milestone in KB-SDP-010 |
| owner | Codex |

Registered from the owner conversation on 2026-09-23. The timestamp records registration, not a reconstructed discussion time. The directory and ledger record lifecycle status.

## Background and owner direction

SDP should use its own process: the repository's `SDP/` owns development documentation for the SDP product. Following planning, `Template/` should own templates used by other projects. Inspect existing numbered root directories before moving them; do not mix templates with actual project history.

The eventual aim is separate GitHub repositories for SDL and SDUI, linked into SDP. Tools should work with identical subdirectories in a monorepo and separate repositories. Each project has its own SDP area, plan, ledger and KanBan. The Git integration mechanism is undecided.

## Process and numbering — proposals to resolve

Retain numbered directories, aligned more closely with existing A0–A5 abstraction levels. The mandate comes from the Project Owner, who need not be a developer. Mandate is a process input, not automatically another SDL abstraction level. Do not replace A0–A5 or equate implementation phases with abstraction levels without an explicit decision.

- `01--Mandate/Mandate.md`: the owner's assignment.
- Optional `01--Mandate/Mandate-study.md`: understanding, uncertainties and possible approaches. Interpretation does not automatically become an approved mandate or requirement.
- Together, mandate and optional study must identify stakeholders.
- Propose removing a mandatory standalone `02--Study`, with optional studies per phase instead.
- `02--Requirements/`: optional Requirements-study, `01--Actors.design`, `02--UserStories.design`, `03--UseCases.design`. Subsequent numbering remains open.
- Express actors/stories directly in readable SDL; [KB-SDL-001](../../../../SDL/SDP/Agents/KanBan/backlog/%23001--Proposal--Requirements-narrative.md) owns language decisions.

## Next bounded work

R1 delivered the inventory and physical migration. Next: a coherent phase/template profile with studies per phase, A0–A5 mapping and requirements modeling; project registration independent of Git boundaries; and a possible SDL/SDUI extraction plan. Start from the existing Toolkit installer and KB-SDL-001. Submodule/subtree/other integration remains undecided; physical housekeeping does not justify creating repositories. KB-SDP-010 consolidates the supporting documentation.

## Acceptance for later restructuring

Approved migration plan; verified links and build/installation; no competing active template or language implementation; identical project lookup across Git boundaries. Keep historical evidence as history. After R1, SDL code is in `SDL/go`; `SDL/SDP/` owns its process documents. Cross-boundary project lookup remains acceptance for future tools, not a delivered R1 feature.

## Related foundations

[Viewpoint levels and notation](../../../../SDL/docs/integration/SDL-Viewpoint-Levels-and-Notation.md)
[Toolkit installation manifest](../../../../Toolkit/SDP-install.manifest.json)
[KB-SDP-002 — sdptool](../backlog/%23002--Proposal--sdptool.md)

## Activation, 2026-09-24

The owner requests repository housekeeping and activation of this card. R1 inventories ownership, gathers templates, moves the project's operational documents into SDP, gathers SDL material and establishes documentation entry points/migration maps. The [R1 plan](../../../Maintenance/R1/Plan.md) records scope and limits. KB-SDP-010 separately owns editorial consolidation.

## After R1 — delivered and remaining

Templates are in Template; project records are in SDP; SDL code and language documentation are in SDL. Checkpoint #1 is shared, dated history in SDP/History. Earlier locations above describe the migration's background.

The card remains active for the final phase/template profile (per-phase studies and requirements model), project registration across repository boundaries, and a separate-repository plan for SDL/SDUI. R1 creates no submodules or GitHub repositories. KB-SDP-010 owns editorial consolidation; KB-SDP-011 owns existing ID mismatches. Keep these separate from physical migration. See the R1 plan and evidence.

## Worklog and revisions

This log starts after R1. Earlier revisions remain in Git; do not invent retroactive log entries. [History/diff](../History.md).

| Time | Actor / event | Work and outcome | Evidence / remaining work |
| --- | --- | --- | --- |
| 2026-09-24T14:52:35Z | Codex; EVT-KB-SDP-000021 | Reviewed R1: M1 `d269bc7`, M2 `f42859e`, M3 `f722dc2` delivered physical organization. | [R1 evidence](../../../Maintenance/R1/Evidence.md). Phase/template profile, project registry and repository extraction remain; KB-SDP-010 tracks consolidation. |
| 2026-09-24T15:03:31Z | Codex; EVT-KB-SDP-000024 | R2 corrects the current next steps and locations; no new migration/profile adopted. | KB-SDP-010 gathers the documentation basis; phase/template profile and project registry are next. |
