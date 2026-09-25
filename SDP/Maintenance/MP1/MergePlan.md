# MAINT-SDP-0005 — MergePlan: consolidate SDP development into main

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0005 |
| project | SDP |
| state | planned |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| source | Owner request 2026-09-25: prepare a MergePlan and tag current main as old |
| Systems | SDP, SDL, SDUI, SDPTOOL |

## Outcome and authority

Bring the accepted, verified SDP development stack onto main without losing
milestone history, project records or earlier review references. MergePlan is
this MaintenancePlan's purpose, not an unimplemented seventh PlanType.
The owner requested plan preparation and the archive tag. This document does not
record authorization to execute the future target-branch merges. Existing phase
commit/push/PR authority remains applicable to preparing the work.

No new Scrum, Sprint or wrapper card is required. Dependencies KB-SDP-011 and
KB-SDP-030 retain their own bounded acceptance. This is integration, not a Toolkit
release: do not change unreleased manifests, create a version tag or publish a
GitHub Release merely because main advances.

## Preparation already performed

[Observed facts and PR/branch inventory](Preparation.md) and the
[machine-readable snapshot](Preparation.json) establish:

- main: 2cb49c02145621b099c47d05786716598e414e75;
- sdp-vNow: 9ad432407004080dd7f4f0ab06d107523f4316fd;
- initial candidate: 481cafdc0279d2b28781ccaef08112de3312818c, 121 commits ahead
  of main and 92 ahead of sdp-vNow, with neither target having unique commits;
- annotated, remotely verified old tag preserves that exact main commit;
- merge-tree of main and the candidate is conflict-free and equals the candidate
  tree. This is a dry inspection, not a completed merge;
- three CI jobs failed. A mergeable Git graph does not establish readiness.

The planning commit and subsequent fixes will change the candidate. Record its
full SHA/tree and exact CI run before merging; do not reuse the initial SHA as
if it includes later work. Never move or overwrite old. If main changes before
execution, inspect the new commits and revise the integration baseline.

## Scope and disposition

The current stack contains all observed phase branches and the head of draft
PR #4. PRs #12–#34 overlap; only the latest verified combined candidate needs to
be merged. Preserve their branches and historical commits. Do not squash or
rebase the stack, or cherry-pick all its already-contained phase commits.

Draft PR #8 / codex/issue-7-provisional-vnext-pilot is separate: 34 unique commits,
experimental process authority, and conflicts with the current folder layout.
Recommended disposition: preserve the draft/branch separately and exclude it
from this main promotion, documenting why. This is a recommendation requiring
owner disposition of what “all” includes before the integration milestone. If
selected for inclusion, map its useful content against adopted contracts, resolve
conflicts in a separate work branch and run its own checks plus current checks.
Do not merge it wholesale or call it incorporated/superseded merely from age.

SDL/go/sourceinput is unrelated untracked work. Preserve it without staging,
cleaning or stashing it as part of this plan. Use a disposable checkout for
candidate verification. XFMD and renderer repositories are outside scope.

## Selected integration route and Git policy

Use one maintenance branch, sdp/maintenance-mp1-main-integration, descended from
completed PL1. Commit each real milestone; no branch per phase is needed. Keep
fixes bounded and include their card/milestone IDs. Push completed milestones.

After readiness and owner merge authorization, use two explicit integrations:
latest combined candidate → sdp-vNow → main. This preserves the agreed staging
target and provides one final main PR. Use merge commits to preserve milestone
ancestry and traceability. A direct main PR would also contain the stack today,
but would leave the agreed staging target stale; do not choose both routes.

Pin the PR head with gh pr merge --merge --match-head-commit <verified-full-sha>.
Select the actual latest combined PR after readiness fixes; #34 is only the
current implementation baseline. Do not enable auto-merge while gates are red,
use admin bypass, force-push a target, or delete source branches. The final main
PR head is the verified staging result, not an old phase head. Refresh both target
and candidate identities immediately before each integration.

## Phases and milestones

| Phase | Milestone | Delivery and acceptance | State |
| --- | --- | --- | --- |
| MP1-R — readiness | MP1-R-M1 | Resolve explicit inclusion/exclusion of PR #8; refresh branch/PR/local-work inventory and freeze intended scope | Planned |
| MP1-R | MP1-R-M2 | Repair KB-SDP-011's 38 findings without rewriting historical evidence; full Toolkit validator and unit suite pass | Planned |
| MP1-R | MP1-R-M3 | Repair KB-SDP-030: portable historical artifact provisioning and audited v1 expected outcomes; Linux/Windows CI jobs pass including the previously skipped fault matrix | Planned |
| MP1-V — candidate verification | MP1-V-M1 | Fresh disposable checkout of the exact candidate; current contract, language, consumer and process checks pass; independent review covers fixes and candidate; record full SHA/tree and CI run | Planned |
| MP1-I — integration | MP1-I-M1 | With owner merge authorization, merge the verified combined candidate to sdp-vNow; verify remote ancestry/tree and record actual merge SHA | Planned |
| MP1-I | MP1-I-M2 | Open/verify sdp-vNow → main PR; require green checks for that exact result; merge and verify main contains all selected milestone commits and old still resolves unchanged | Planned |
| MP1-C — closeout | MP1-C-M1 | Reconcile included PRs, main-target working guidance and management/evidence records; preserve branch/tag provenance and unrelated work; final status names actual main SHA | Planned |

## Verification contract

Resolve existing failures rather than treating a known-red baseline as a green
merge gate. Do not weaken validators, delete checks or regenerate reference
outcomes only to remove red status. Historical IDs/events remain immutable;
compatibility mappings or current-reference repairs need explicit evidence.
Use KB-SDP-011's baseline and the actual CI diagnostics in Preparation.md.

At the exact candidate, run the workflow's contracts, Windows installer and
Linux process-profile jobs through completion. Also run current Go tests/race/vet
for SDL, SDUI and SDPTool from their module roots, the Python parser/toolkit suites
where still maintained, shared management/lineage checks and document verification.
Check generated viewpoint provenance with the owning SDL tools. Do not hand-edit
generated artifacts. Preserve previous evidence with its original candidate.

A local pass cannot replace the currently failing host/workflow results. If any
platform/capability is deliberately excluded, that changes accepted scope and
requires explicit owner disposition and truthful documentation before merging.
Freeze/retest affected claims after any code/configuration change. Record external
infrastructure failures separately from tested product failures; neither is a pass.

## PR reconciliation and completion

After verifying main, re-enumerate PR heads and ancestry. Close only superseded
PRs whose exact heads are proven included; GitHub may auto-close contained PRs.
Keep #8 open unless its separate disposition says otherwise. Existing instruction
not to send messages without authorization still applies to optional PR comments.
Do not delete branches to make the list shorter. No issue is closed automatically
just because its PR head is an ancestor.

Update AGENTS/current work guidance to main as the integration target after the
owner's main transition, while preserving dated sdp-vNow branch history. Keep
post-merge administrative edits on a working branch and include them through a
small reviewed follow-up PR if they cannot truthfully be recorded before merge.
Record management movements here; code/design fixes and their verification belong
in Traceability with this Maintenance ID. Completion requires actual remote merge
verification and truthful closeout, not merely creating this plan or a PR.

If verification fails before integration, leave target refs unchanged and repair
on the work branch. If a merged result needs reversal, preserve old and history
and prepare an explicit revert of the identified merge with its effect reviewed;
do not reset or force-push main back to old.

## Preparation verification

Git status, remote heads/tags, branch ancestry, open PRs, current check logs and
merge-tree inspection were performed read-only apart from the requested archive
tag and this planning delivery. No target branch has been advanced and no future
readiness check is marked passed. The plan remains planned until execution is
selected. The old archive tag is already published and remotely verified.

Planning record checks pass: 37 cards, ten management records, three lineage
operations and 250 events; six management test groups. Document verification
preserves 105 frozen records/ledger prefixes and 574 generated outputs and resolves
2,451 local links and 130 fragments. git diff --check passes. These results verify
this plan/history, not the future merge candidate or the still-failing CI gates.
