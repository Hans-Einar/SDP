# Scrum-0003 — skills and agent activation

| Field | Value |
| --- | --- |
| id | SCRUM-SDP-0003 |
| project | SDP |
| state | completed |
| source | Owner request for the next backlog Scrum after Sprint-0001; KB-SDP-027 |
| participants | Project Owner via conversation; Codex inspection and recording |
| outcome | MAINT-SDP-0002 planned |

## Review and selected outcome

Reviewed [KB-SDP-027](../../KanBan/completed/%23027--Study--Skills-review-and-project-activation.md)
after the SDPTool Sprint. Compared the ten distributed roles with thirteen draft
roles, current project/management contracts, installation consumers and current
Codex documentation. A fresh local catalog probe in disposable repositories checked
the proposed discovery path without activating project skills or executing agent
work. This is a real backlog review, not an invented meeting or release approval.

Select [MAINT-SDP-0002](../../Maintenance/SK1/Plan-and-Evidence.md) directly, without
an additional wrapper card or Sprint. It contains the adoption map, migration
boundary, three implementation milestones and observable acceptance criteria.
The thirteen draft roles are the editing base; adaptation is required before
promotion. Root Skills/ will be the sole authored collection, with a Codex
.agents/skills adapter. A root SKILL.md is unnecessary for the tested mechanism.

The installed Codex CLI 0.156.1 discovered the symlinked draft router. It did not
discover root Skills/ or root SKILL.md alone and rejected the old Master's missing
native metadata. This measures catalog discovery, not successful loading or
behavior. The plan keeps those later checks separate and retains the evidence.

## Dispositions and ordering

| Record | Disposition |
| --- | --- |
| KB-SDP-027 | Complete the Study after delivering this Scrum, inventory, plan and verified discovery finding; it does not close Maintenance execution |
| MAINT-SDP-0002 | Planned canonical skills adoption and demonstrated activation, with one phase and three milestones |
| KB-SDP-028 | Keep backlog; recommended next Scrum to settle the coupled installer migration boundary before skills distribution changes |
| KB-SDP-018 | Keep broad Toolkit audit separate; reuse this inventory and avoid repeating the skills review |
| KB-SDP-014 | Keep reusable process/KanBan version contract separately owned; no new published profile here |
| KB-SDP-010 | Existing owner gate-review unchanged and non-blocking |
| KB-SDP-017 | Existing ready producer remainder unchanged; native XFMD work stays in its repository |
| KB-SDL-005 / KB-SDP-020 | Source contract and dependent model migration remain separate backlog work |

Do not merge #027 and #028: one owns skill content/discovery, the other owns safe
process installation/upgrades. They share a concrete manifest/adapter boundary,
recorded in both primary documents. No held cards were selected. Planning does not
activate skills, update consumer installations or authorize an unrelated release.

## PM3-M1 evidence

Base c4aed09944235f5da1d247001688f6f162ec124d; phase branch
sdp/phase-pm3-skills-scrum. The local probe and its limitations are linked from the
Maintenance plan. Shared management and documentation checks are recorded below
when completed. No SDL/SDUI product change or Traceability transition belongs to
this management-only milestone. The untracked SDL/go/sourceinput draft is preserved.

PM3-M1 checks passed: 35 cards, seven management records, three lineage operations
and 219 events; all four management test groups and 15 lineage negative cases.
The previous 212 management events remain an exact byte prefix; Traceability is
unchanged. Documentation verification preserves 105 frozen records/prefixes and
574 generated artifacts and resolves 2,408 file links and 130 fragments.
git diff --check passed. The four-case native catalog probe is recorded in the
Maintenance evidence; it ran no model/agent turn and made no project activation.

The Scrum and Study are complete; Maintenance is planned. The next recommended
review is #028, using this handoff to avoid incompatible installer/skill changes.
