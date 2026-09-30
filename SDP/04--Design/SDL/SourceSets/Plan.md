# SDL System and source sets — DesignPlan

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0010 |
| project | SDP |
| state | completed |
| PlanType | DesignPlan |
| BranchPolicy | phase |
| CommitPolicy | milestone |
| Systems | SDL, SDPTOOL |
| source | KB-SDL-005; owner request to continue Session 0001 on 2026-09-30 |

## Outcome and authority

Execute S1 of [Session 0001](../../../Sessions/session-%230001--SDL_expansion.md):
design one explicit System across named source files, with file-aware diagnostics
and one revision consumed by checking, navigation and document generation.
[KB-SDL-005](../../../KanBan/active/%23005--SDL--Change--System-and-source-sets.md)
remains the delivery card; design completion does not complete its implementation.
The owner authorized the next Session step. Product coding, broad model migration,
merging and release are outside this DesignPlan.

## Baseline and scope

Start at 84fbdb34f7b87f2ef23547cc4b1684322916a342. Reuse the
[gap study](../../../02--Requirements/XFMD-Gaps/Study.md),
[source-tree study](../../../../SDL/docs/studies/SDL-Source-Tree-and-Compilation-Study.md)
and [ecosystem boundaries](../../../03--Architecture/Ecosystems/Decisions.md).
Design-core 0.5 remains the implemented single-file profile. The untracked
SDL/go/sourceinput draft is inspected but neither adopted nor modified.

Resolve input/version dispatch, System and file membership, name resolution,
public/cross-system boundaries, canonical form, provenance, revision freshness,
consumer migration and negative cases. No action/class-language extension, new
SDUI widgets, external XFMD changes or replacement parser is included.

## Phase and milestone

| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| SSD1 — bounded design | SSD1-M1 | Contract and alternatives; physical proposed examples; current-code evidence; consumer/acceptance matrix; explicit migration limits; Session/card handoff | completed |

## Git policy

Honor KB-SDL-005's phase-branch commitment: sdp/sdl-source-sets/design is based
on sdp/ecosystem-models at the baseline above. Commit SSD1-M1 when the complete
design and record checks pass. Existing owner policy permits a phase push; no
merge or release is selected. Preserve unrelated sourceinput and Node files.

## Verification and deliverables

- [Contract](Contract.md): recommended next-profile design, not published grammar.
- [Acceptance](Acceptance.md): concrete positive/negative and consumer cases.
- [Evidence](Evidence.md): inspected source identities and executed checks, with limits.
- examples/: proposed multi-file source input, explicitly not current syntax.

Check against current Go APIs and run bounded baseline probes using the existing
parser. Check management schemas, relative links and diffs. No new parser is
written to make proposed examples pass. Product verification and independent
review belong to Session S4; this design receives author consistency checks only.

## Handoff

SSD1-M1 delivers the contract, fixture, 26 acceptance cases and baseline evidence.
Session S2 will select an ImplementationPlan with runnable
increments and the acceptance cases defined here. KB-SDP-020 remains backlog;
no existing model or installed template moves during this work.

## SSD1-M1 delivery — 2026-09-30

Design completed; product implementation has not started. The recommended first
increment is one closed System with explicit source membership, file-aware
provenance and aggregate revisions. Public linking/exports remain explicit later
scope of KB-SDL-005. Current authoritative models remain unchanged.

The active card returns to ready at this handoff: it still owns the unfinished
language/input/consumer delivery, and Session S2 is next. This is not an owner
review gate or card completion. EVT-PM-SDP-000108 records DesignPlan completion.

Closeout checks pass: 52 cards, 25 management records, 3 lineage operations,
406 management events; Toolkit repository validation; local file links in 18
changed/new Markdown files; preserved historical ledger prefixes; git diff --check.
The current Go Frontend baseline check/canonical probe passes; proposed 0.6 is
correctly rejected; all 57 baseline declarations/statements survive the split.
These are author design checks, not independent implementation verification.
