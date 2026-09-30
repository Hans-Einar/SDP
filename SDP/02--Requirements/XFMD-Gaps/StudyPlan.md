# XFMD language and UI gap StudyPlan

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0009 |
| project | SDP |
| state | completed |
| PlanType | RequirementPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDL, SDUI, SDPTOOL |
| source | Owner authorization 2026-09-29; KB-SDP-041 |

## Outcome and authority

Execute [KB041](../../KanBan/completed/%23041--Study--XFMD-driven-SDL-and-SDUI-gaps.md):
evaluate all 14 SDL/SDUI gaps and six XFMD-local concerns in the external register.
Deliver evidence-qualified dispositions, alternatives, ownership and a proposed
phased implementation/verification sequence. StudyPlan uses the adopted
RequirementPlan record type. Research/probes are authorized; product implementation,
new grammar adoption, external project edits, merge and release are not.

## Phases and milestones

| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| XGS1 | XGS1-M1 | Pin current source/report identities and reproduce bounded capability probes; inspect existing runtime/host/profile facilities | Completed |
| XGS2 | XGS2-M1 | Disposition of every gap, alternatives and minimal pilot; distinguish grammar from host/provider/application responsibilities | Completed |
| XGS3 | XGS3-M1 | Proposed phased delivery and acceptance matrix, follow-up ownership, evidence/record validation and study closeout | Completed |

## Method and boundaries

Use tracked Go SDL/SDUI sources and minimal positive/negative probes, plus relevant
existing tests. Keep compiler caches/artifacts under /tmp. Do not use the unrelated
untracked SDL/go/sourceinput draft or create another parser. Read external XFMD
documents and selected source evidence without writing to its worktree. Preserve
all source gap IDs and existing language-profile authorities. External report
recommendations are input, not accepted syntax or implementation promises.

Visual claims require inspected artifacts; where no native host or renderer is
available, report the limit rather than implying static checks prove interaction.
Local source/tests are the primary evidence for current repository behavior;
no third-party API change or framework selection is made by this study.
Research and record checks are distinct from independent product acceptance.

## Git policy and environment

Current branch: sdp/ecosystem-models, starting at a274265. Intended milestone
commits were XGS1-M1, XGS2-M1, XGS3-M1. Git writes were unavailable during study
execution, so all three completed milestones accumulated as local files. On
2026-09-29 the owner restored unrestricted access and explicitly requested a
commit. Recover these together in one coherent commit naming all three
milestones; do not reconstruct fictitious intermediate execution states. This is
the documented recovery exception to CommitPolicy milestone. Later work retains
its own selected policy. Unrelated sourceinput and Node work remains excluded.

## Result register

- [Study](Study.md): findings and complete disposition matrix.
- [Evidence](Evidence.md) and [source pins](evidence/source-manifest.json): executed commands, fixtures, identities and limits.
- [Delivery proposal](Delivery-Proposal.md): recommended phases/acceptance; no activated implementation plan.

Completion means delivered research and a concrete successor proposal, not closed
product gaps. The source register remains external and its entries are not edited.

## Delivery record

- XGS1-M1 — EVT-PM-SDP-000104: 21 capability cases and five supplemental probes; ten Go
  packages pass after restoring tracked archive fixtures. One configured-renderer
  test skipped; initial failures retained. Two static typography exports inspected.
- XGS2-M1 — EVT-PM-SDP-000105: all 14 language/UI and six XFMD-local IDs dispositioned;
  recommended boundaries and alternatives recorded without adopting syntax.
- XGS3-M1 — EVT-PM-SDP-000106: phased proposal and two SDUI follow-up cards delivered;
  four existing producer cards updated with study input. KB041 closes as research.

No native GUI, XFMD integration or independent review was performed. No product
code or accepted language profile changed during the study. Git recovery and the
subsequent owner scope clarification are recorded above and in the study; the
original evidence and append-only events retain their execution-time context.

## Closeout verification

- ProjectManagement validation passes: 51 cards, 24 management records,
  3 lineage operations, 390 events.
- SDP Toolkit repository validation passes.
- Local Markdown links checked in the study, plan, evidence, index, completed
  card and six affected successor/owner cards: no missing file targets.
- All 20 external gap IDs have explicit study dispositions. Selected source
  hashes still match; evidence artifacts are recorded in
  [artifact-hashes.json](evidence/artifact-hashes.json).
- git diff --check passes. No unrelated sourceinput/Node work changed or staged.
- The formatted glyph fixture was rerun and its output remains byte-identical.

These are author checks, not independent review or owner acceptance.

## Owner follow-up and Git recovery

The owner reaffirmed limited UI concepts/layout plus basic SDL-connected execution,
with a possible tree addition. Study, delivery proposal and KB-SDUI-003 now make
rich-application parity optional and distinguish native XFMD integration from
SDUI development. No implementation was selected by this clarification. The
completed study remains completed; follow-up is recorded as a KanBan review.
