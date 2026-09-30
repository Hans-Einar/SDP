# SDL source composition — ImplementationPlan

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0012 |
| project | SDP |
| state | completed |
| PlanType | ImplementationPlan |
| BranchPolicy | phase |
| CommitPolicy | milestone |
| Systems | SDL, SDPTOOL |
| source | Owner authorization to execute Session 0001 S2–S5 in one turn, 2026-09-30; KB-SDL-005; KB-SDL-007 |

## Selected outcome

Deliver source-owned design-core/0.6 composition through check/AST/format,
viewpoints, broker and SDPTool. Use the [SSD2 contract](../../../04--Design/SDL/SourceComposition/Contract.md)
and [KB007](../../../KanBan/completed/%23007--SDL--Proposal--Composable-file-ASTs-and-contextual-analysis.md).
No external manifest, second parser, native XFMD edits, action-runtime extension,
merge or release. Root-relative includes, path contains, one System per compiled
context, preserved file ASTs and complete source revisions are required.

Select canonical separate System declaration and membership sentences; compact
System-plus-contains syntax is deferred spelling, not needed for this delivery.
Select a pure syntax cache reusable before root arrival, context-specific full
semantic rebuilding and explicit fragment inspection. This proves reusable syntax,
not optimized dependency-selective semantic recompilation or runtime hot reload.

## Phases and milestones

| Phase | Milestone | Observable acceptance | State |
| --- | --- | --- | --- |
| SSI0 | SSI0-M1 (S2) | Plan, scope, original/updated acceptance and card activation | completed |
| SSI1 | SSI1-M1 (S3) | Go syntax, pure file-AST cache, source graph, validation and CLI with original spans; cycles/late root/negative tests | completed |
| SSI2 | SSI2-M1 (S3) | Same checked result in views, broker and SDPTool; aggregate revision, source maps and stale/output guards; consumer tests | completed |
| SSI3 | SSI3-M1 (S4–S5) | Bounded real-model workflow, regressions, independent review, documentation and truthful Session/card dispositions | completed |

## Git and evidence

Start at 32aa68a. Phase branches stack as sdp/sdl-source-sets/implementation-plan,
frontend, consumers and verification. Commit each milestone and push its completed
phase under existing authorization. Preserve unrelated untracked SDL/go/sourceinput
and Node files. New code uses a separately named sourcegraph package, not the
unadopted manifest draft. Run Go with GOMAXPROCS=2 and -p 2; never PowerShell.
[Evidence](Evidence.md) records commands/candidates, results and limitations.

Independent review uses a fresh Reviewer context as required by the loaded
Reviewer skill and Session S4. It may challenge scope and consumers; no invented
approval or self-review substitution. Fix material findings before closeout.

## Acceptance and card boundary

SSD2 replacement cases and KB007's six scenarios govern. Retain SSD1 original-span,
identity, formatting, limits, profile isolation, all-input revision, scenario order,
consumer consistency and publication protection checks; manifest-only cases are
superseded. Retain existing 0.5/action/class tests. CLI, broker and SDPTool evidence
is required; no claim of native XFMD acceptance or full experimental MVP1 support.

Use a copied real Frontend model as the positive system pilot; leave current
registration and model authority in place. KB-SDP-020 stays backlog. Public
cross-System exports/imports were explicitly excluded from this Session's bounded
System goal; KB-SDL-005 retains that remaining obligation after this increment.
Do not mark its entire historical scope complete: return it to backlog with the
completed increment and remaining work explicit. KB-SDL-007 may complete when
its selected late-root/partial-analysis behavior and limits are verified. Session
can close its bounded goal with these actual dispositions, not fictitious closure
of the broader language card. KB043 timeline automation remains separate backlog.

## Outcome

SSI0–SSI3 delivered the bounded System/source-composition increment and its
consumer workflow. Independent product review accepted the corrected candidate.
Full tracked headless Go regressions (including race detection) and the regenerated
constructor check passed; see evidence for the initial failures and corrections.
KB-SDL-007 completes its selected bounded scope; KB-SDL-005 returns to backlog
for public cross-System linking. KB-SDP-020 remains separately selectable migration.
Session 0001 S2–S5 are completed; merge/release and native XFMD acceptance are not
claimed. Further work selects a new plan instead of reopening this delivery.
