# MAINT-SDP-0005 — MergePlan: consolidate SDP development into main

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0005 |
| project | SDP |
| state | active |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| source | Owner request 2026-09-25: prepare a MergePlan and tag current main as old |
| Systems | SDP, SDL, SDUI, SDPTOOL |

## Outcome and authority

Bring the accepted, verified SDP development stack onto main without losing
milestone history, project records or earlier review references. MergePlan is
this MaintenancePlan's purpose, not an unimplemented seventh PlanType.
The owner requested plan preparation and the archive tag on 2026-09-25. On
2026-09-26 the owner selected the current SDP-vNow lineage as the new SDP, excluded
Issue #7 explicitly, and requested backlog capture before merging to main. This
authorizes execution of this integration plan, including its readiness repairs
and target-branch merges after verification; it does not waive failed checks.

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
Owner disposition, 2026-09-26: exclude all Issue #7-specific work from main.
The former vNext pilot stopped when SDL/SDUI development began; the current
SDP-vNow lineage is the selected successor. Preserve the pilot branch and draft
as historical work; do not merge/cherry-pick its contracts, validator or records.
The already shared Issue #5 study remains part of the accepted baseline.
Capture new owner intent in KB-SDP-031/032/033 without adopting pilot schemas.

SDL/go/sourceinput is unrelated untracked work. Preserve it without staging,
cleaning or stashing it as part of this plan. Use a disposable checkout for
candidate verification. XFMD and renderer repositories are outside scope.

## Selected integration route and Git policy

Use one maintenance branch, sdp/maintenance-mp1-main-integration, descended from
completed PL1. Commit each real milestone; no branch per phase is needed. Keep
fixes bounded and include their card/milestone IDs. Push completed milestones.

With owner merge authorization now recorded, after readiness use two integrations:
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
| MP1-R — readiness | MP1-R-M1 | Resolve explicit inclusion/exclusion of PR #8; refresh branch/PR/local-work inventory and freeze intended scope | Complete |
| MP1-R | MP1-R-M2 | Repair KB-SDP-011's 38 findings without rewriting historical evidence; full Toolkit validator and unit suite pass | Complete |
| MP1-R | MP1-R-M3 | Repair KB-SDP-030: portable historical artifact provisioning and audited v1 expected outcomes; Linux/Windows CI jobs pass including the previously skipped fault matrix | In progress |
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

### MP1-P2 — preparation visualization

Added a condensed Mermaid gitGraph to Preparation.md, distinguishing observed
commits from proposed readiness and integration nodes. Explained merge-tree's
file-tree result and pinned the reproduction commands to the original snapshot.
Corrected the inventory count: 50 inspected refs include one symbolic alias,
leaving 49 actual branches. The raw observation JSON remains unchanged.

The pinned merge-tree and candidate tree IDs match. The embedded diagram was
extracted and rendered with the local mermaid-rs-renderer target/debug/mmdr,
then visually inspected as PNG: both merge connections and snapshot labels are
visible. Bare merge statements avoid that installed renderer's mishandling of
merge attributes; hidden automatic commit labels prevent invented hashes from
appearing as evidence. This is a document-rendering check, not an XFMD GUI test.
The MergePlan remains planned; no target branch was merged.

MP1-P2 checks pass: management validation covers 37 cards, ten management records,
three lineage operations and 251 events; document verification preserves 105
frozen records/ledger prefixes and 574 generated outputs and resolves 2,452 local
links and 130 fragments. git diff --check passes.

### MP1-P3 — XFMD-compatible graph

The owner reported that XFMD rejected MP1-P2's init directive. The earlier
standalone renderer check did not establish compatibility with XFMD's interpreter.
Replaced the graph with XFMD's supported subset: plain gitGraph header, explicit
IDs without spaces after id:, underscore branch aliases, and bare merges.
Removed tag/type attributes and the init directive. The adjacent legend preserves
the real branch names, archive tag and hypothetical merge status. Synthetic merge
labels are explicitly identified as diagram IDs, not Git commit hashes.

Verification on 6cac4e6 plus this document diff: extracted the exact Mermaid block
and passed it through xfmd-sdl-navigation/build/cargo/release/libxfmd_mermaid_ffi.a
using a temporary C++ ABI harness. xfmd_mermaid_parse_v1 returned status 0 (564
bytes); xfmd_diagram_layout_measured_v1 with Pango text metrics returned status 0
(6,019 bytes). An initial unmeasured layout call correctly required external text
metrics. Standalone rendering was also visually inspected. This checks XFMD's
built parser/layout library, not its running GUI or an unknown installed binary.
No XFMD source or target Git branch was changed. Merge execution remains planned.

### MP1-R-M1 — owner scope decision and backlog intake

The owner excludes Issue #7-specific work and selects current SDP-vNow as the new
SDP. Registered KB-SDP-031 (SDL-derived assignment bundles/blueprints), KB-SDP-032
(viewpoint navigation feedback) and KB-SDP-033 (gh-sdp/XFMD adoption and SDL pilot).
These remain backlog studies; they are not implementation prerequisites for this
merge. Readiness repairs KB-SDP-011 and KB-SDP-030 remain required. Integration
execution is now selected; earlier planned-state observations remain historical.

MP1-R-M1 verification: remote main/staging/pilot heads and old tag still match the
recorded snapshot; PR #35 is the combined candidate. Backlog and document checks
pass after activating readiness dependencies. No Issue #7-specific commits are
selected, and SDL/go/sourceinput remains untracked and excluded.

### MP1-R-M2 — compatibility repair

[Evidence](Evidence.md) records the resolution of all 38 Traceability findings.
KB-SDP-011 is completed; installation host tests and CI remain separate gates.
