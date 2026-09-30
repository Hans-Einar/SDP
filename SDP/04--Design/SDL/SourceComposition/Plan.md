# Source-owned SDL composition — design correction

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0011 |
| project | SDP |
| state | completed |
| PlanType | DesignPlan |
| BranchPolicy | phase |
| CommitPolicy | milestone |
| Systems | SDL, SDPTOOL |
| source | Owner correction 2026-09-30; KB-SDL-005; Session 0001 |

## Outcome and authority

Replace the separate authored design-set recommendation with inclusion controlled
by SDL source. Define a concrete proposal for path-addressed System membership,
load-once graph assembly, original-source diagnostics and consumer entry selection.
[Contract](Contract.md) is the current recommendation for Session S2. The owner
selected source-owned inclusion; detailed spelling and error policies below remain
agent design recommendations, not implemented grammar or claimed owner approval.

[PLAN-SDP-0010](../SourceSets/Plan.md) remains a completed historical delivery.
Its external manifest recommendation is rejected. Its original evidence and
hash records remain unchanged; its contract/acceptance receive supersession notices; completion did not mean owner acceptance.
The new work is recorded separately under the plan contract's successor rule.

## Phase and Git policy

| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| SSD2 — source composition | SSD2-M1 | Recover earlier System/include intent; replace external manifest authority; specify repeated/cyclic loads and path semantics; reconcile Session/card and consumer acceptance | completed |

Use sdp/sdl-source-sets/source-composition based on cfd7085, stacked above the
previous design branch, preserving KB-SDL-005's phase policy. Commit at SSD2-M1;
phase push is authorized by existing policy. No implementation, merge or release.
Preserve unrelated sourceinput and Node work. No external XFMD changes.

## Evidence and verification

Inspected historical inputs:

- [Language candidates §13](../../../../SDL/docs/studies/SDL-Language-Candidates.md#13-planned-workspace-extension-one-declared-system): explicitly selected source entry and one declared System; bare symbolic contains does not load a file.
- [Source-tree study](../../../../SDL/docs/studies/SDL-Source-Tree-and-Compilation-Study.md): entry owns explicit source membership; avoid textual substitution and preserve file ASTs.
- [Experimental MVP1 entry](../../../../experiments/mvp1_sdl/SDL/MVP1/System.design): existing includes statements plus System membership, under sdl-mvp1-exercise/0.1, not released parser support.
- [SSD1 evidence](../SourceSets/Evidence.md): Go parser remains single-file design-core/0.5. This correction changes documents only; those product observations remain valid.

No new grammar probes or renderer runs are needed to claim this bounded document
correction. Check local links, management schema/history, Toolkit repository
validation, preserved ledger/evidence bytes and git diff --check. Do not present
proposed graph cases as passing parser tests. Independent implementation review
remains Session S4 work.

## Handoff

SSD2-M1 delivered the source-owned composition recommendation. Session S2 must use this source-owned design, not the rejected
external-manifest path. Keep KB-SDL-005 open for implementation and remaining
public cross-System boundary work. KB-SDP-020 migration still awaits a working
source/consumer pilot; current models and registrations are unchanged.

## SSD2-M1 closeout

Management validation passes: 52 cards, 26 management records, 3 lineage operations,
410 events. Toolkit repository validation and git diff --check pass. Local links
in all six changed/new Markdown documents resolve. Original append-only ledger
prefixes and SSD1 machine-evidence bytes are unchanged. The superseded contract
and acceptance have only historical notices added; their prior hashes still name
the cfd7085 delivery, not the notice-bearing working copies. No product tests or
new-language acceptance is claimed for this documentation correction.
