# Interim project governance overview

Owner update — 23 October 2026: all remaining laptop work has been handed over.
The new status study and proposed consolidation route are maintained in
[Session0013](../../Sessions/session-%230013--SDP_consolidation_and_Ponsse_readiness.md)
and [WR1 Study](../../Maintenance/WorkstationReadiness/Study.md). Use those for
current readiness planning; the October 10 snapshot below remains historical.


Current recovery note — 23 October 2026: the snapshot below is historical. SDP main is now locally consolidated at bfe3b40; ModelGovernance and Actions are recorded complete, PG has advanced through Session0007 T028, and the selected XFMD viewpoint round is complete with residual handoffs. See [Session0012 T002](../../Sessions/session-%230012--Interim_project_coordination.md#t002--2026-10-23-workstation-recovery-and-completion-strategy) for evidence limits and the single-owner completion recommendation. No new execution allocation or Ponsse-readiness acceptance has been made.

Snapshot: **10 October 2026, 04:01:59 Europe/Oslo**. This is the owner's coordination view across SDP and XFMD. It projects existing Sessions, plans, Git candidates and observed publication facts; it is not another lifecycle ledger. Refresh it before allocating overlapping work. Agent activity and messages are not monitored in the background.

**Priority: establish one combined SDP candidate and reconcile identities before starting Actions product implementation.** PG has an independently bounded next increment; XFMD can prepare its consumer design. Blueprint and diagram work remain paused by the owner. Published release work is complete at its recorded scope.

## Integration baseline

The latest actual merge to SDP main is **PR #52**, `04f88ff918e7c25c9fa1883f061ef20b8c3923a2`, 8 October at 17:19:04 Oslo time. `git ls-remote` confirms that remote main still points there. The local branch named main is an old checkout at `2cb49c0`; it is unsuitable as the comparison baseline.

PR #52 incorporated the authorized inherited baseline and WCI0–WCI4 widget delivery. ModelGovernance and the earlier Sessions are therefore not all new post-merge work merely because their commits appear on a divergent local branch. Compare contents and merge ancestry as well as commit counts.

The actual post-merge focus is:

- Blueprint implementation and its subsequent Actions/JSON contract and role guidance.
- PG controller recovery, skill invocation, owner lab and Project Leader orchestration.
- Installed widget lab, program discovery/launch, SDP 2.2.0 and gh-sdp 0.2.2.
- XFMD viewpoint/renderer corrections, local 0.8 installation and the new dynamic-action UI Session.

## Workstream status

| Workstream | Best current source | Delivered or observed | Still open and next action |
| --- | --- | --- | --- |
| PG runtime and orchestration | PG worktree, Session0007 T014; PLAN-SDP-0018; KB038/052 | PGI core/adapters and bounded live/recovery increments recorded; PGL1 registered-work inspection delivered | PGL2 assignment preparation and revision/authority guards next. PGL3–5, broader review/recovery and owner-trial obligations remain. No implemented general multi-agent coordinator is inferred. |
| Semantic blueprints | Actions worktree, Session0008 T015; KB050; PLAN-SDP-0020 | BPI1–3 analysis, retained catalogue, assessment and revision-bound assignments recorded as delivered/reviewed | Owner paused Session. BP2-A pilot disposition remains; implementation completion does not close owner acceptance. |
| Discoverable actions and JSON API | Actions worktree, Session0011 T008; PLAN-SDP-0024 | Contract proposal and bounded implementation plan; Project Leader role skill and handoff | ACT0/S0 is RED for integration. ACT1 prose is ready; executable schemas/model and ACT2–5 are not delivered. Reconcile shared baseline and identities first. |
| Widgets, lab, runnable applications and release | Release worktree and shared Session0010 T012 | WCI already merged in #52; installed lab delivered; RSP1 program discovery/run; paired SDP 2.2.0 and gh-sdp 0.2.2 publication and local extension upgrade complete | Main integration remains separate. XFMD Run parity with the connected application is still a separate consumer obligation. No release/upgrade task remains selected in this Session. |
| XFMD diagrams and Mermaid | XFMD Session0005 T012; PLAN-XFMD-0009; KB020/028 | N3 work through N3g; ten actual VP07 details reported rendering; local XFMD 0.8 installed, recorded full 91/91 CTest | Owner paused diagrams. Large-overview readability, detail/scenario projections and N4 exact SDL sources remain. Installation did not close these cards. |
| XFMD discoverable menu and toolbar | XFMD Session0006 T001; KB-XFMD-030 | Session/queued card and intended consumer behavior | S1 producer capability/contract inspection next; no implementation plan selected yet. Depends on Actions contract, not on finishing the paused diagrams. |
| Older XFMD navigation, SDUI and modeling | XFMD Sessions0001–0004 and their linked plans/cards | Earlier navigation, standalone Mermaid, heading links and SDUI prototype launch installed at their recorded levels | Review dispositions, broader modeling, runtime limits and KB029 KanBan mutations remain distinct. Session completion does not imply every linked card is accepted. |

SDP Sessions0001–0004, 0006 and 0009 are earlier completed/bounded work; Session0005 transferred its discussion to ModelGovernance and semantic Blueprints. They provide context rather than new implementation assignments. Session0010 publication advanced during this inspection: its initial pending text was superseded by the owning agent's T011 closeout/T012 and observed public releases.

## Exact observed candidates

| Source | Branch | HEAD | Working copy at snapshot |
| --- | --- | --- | --- |
| SDP-vNow shared | sdp/model-governance-implementation | `48c919fda7` | 583 status entries including untracked files |
| PG | sdp/project-governance-pilot | `02e8c9dc4f` | clean |
| Blueprint and Actions | sdp/blueprint-implementation | `732e85cbb6` | clean |
| Runnable programs and release | sdp/release-2.2.0 | `3475367490` | clean |
| XFMD | sprint/009/phase/057-favorites | `41fa4948c5` | 1,885 status entries including untracked files; also staged changes |
| gh-sdp | sdp/release-0.2.2 | `d5097b97d2` | clean |
| Widget lab | See evidence branch field | `b7e15d5a34` | clean |

Status entry counts include evidence, generated files and possible caches; they are not counts of unfinished features. Of the shared SDP dirty regular files inspected, **455 are byte-identical to origin/main**, 27 differ and 99 are absent from main. This strongly suggests substantial already-integrated content remains represented as dirty/untracked in the old worktree. It does **not** authorize deleting or resetting it; the remaining files need content-based ownership/disposition.

## Confirmed conflict surfaces

Fresh `git merge-tree --write-tree` forecasts on the heads above returned exit 1 for all three pairs. These are predicted conflicts on committed candidates, not an unfinished merge. Dirty work is excluded from these forecasts.

| Pair | Conflict paths | Main boundaries |
| --- | ---: | --- |
| Release + PG | 7 | CLI dispatcher, output arguments, README, board, management and traceability |
| Release + Blueprint/Actions | 6 | Contract, output arguments, README, board and both ledgers |
| PG + Blueprint/Actions | 9 | Project Leader handoff, PLR1 maintenance record, KB038, output/presentation, board and ledgers |

There are also semantic identity collisions:

- **KB-SDP-051** means completed Runnable SDUI programs in release/shared, but queued Discoverable SDPTool actions in the Actions branch.
- **MAINT-SDP-0015** means RSP2 release preparation/publication in release/shared, but PLR1 Project Leader role guidance in Actions/PG.
- **14 ProjectManagement event IDs have different JSON records across the four SDP worktrees.** These cannot be concatenated or deduplicated by ID. The scan found no differing Traceability event IDs in that scope; Traceability still has textual merge conflicts.
- The PG branch deliberately imported KB052/PLR1 source records. Reconcile shared provenance rather than creating duplicate transitions or restoring the Actions branch's stale KB038 state.

Until reconciliation, qualify colliding identities by **subject + path + source revision** in all prompts and reports. A normal Git conflict resolution alone cannot make the management history valid.

## Publication and PR status

[SDP 2.2.0](https://github.com/Hans-Einar/SDP/releases/tag/v2.2.0) is published from `ef8741c`; [gh-sdp 0.2.2](https://github.com/Hans-Einar/gh-sdp/releases/tag/v0.2.2) is published from `9776f8a`. Observed extension inventory reports v0.2.2. Owning evidence records an installed default engine at 2.2.0 and an actual connected widget-lab run without development overrides. This audit inspected those records; it did not repeat the native GUI test.

[PR #53](https://github.com/Hans-Einar/SDP/pull/53) and [PR #54](https://github.com/Hans-Einar/SDP/pull/54) remain draft/open against main at observation. #53's branch is an ancestor of #54: #54 contains program work plus release corrections/records. Treat them as one integration lineage, not two independent feature patches. Exact GitHub check state is retained in the snapshot; earlier successful checks are not evidence for a subsequently changed head. Main merge authority remains separate from completed release publication.

## Proposed coordination and integration order

The owner has appointed this conversation for interim coordination and steering assistance. The following worker allocation is a **recommendation ready to send**, not a claim that workers have received it.

1. Give the Actions/Blueprint agent the bounded integrator prompt. Its T008 says the technical information is sufficient, but ownership/authorization was unresolved. Let it own a dedicated integration worktree; it must not absorb arbitrary shared dirty files.
2. Start the candidate from refreshed remote main and the complete published-release lineage through `3475367` or its verified successor. Include #53 once through that lineage.
3. Integrate Blueprint/Actions role guidance and PG from pinned commits, preserving the freshest owning Session/card status. One integrator resolves shared CLI/contract/presentation/dependencies and coordinates the explicit identity/history reconciliation.
4. Have the XFMD agent identify its uncommitted producer-side SDP changes and evidence, and any remaining unowned shared changes. Add only separately identified and reviewed residual packages needed for the candidate. Preserve byte-identical/main leftovers until their owner authorizes cleanup.
5. Verify the combined ModelGovernance, Blueprint assignment, PG and program-launch workflows plus management/history integrity. Then record a fresh green ACT0 baseline and hand it back for Actions implementation.
6. PG may continue scoped PGL2 preparation in its own packages after refreshing preflight; shared contract/dispatcher/history ownership needs explicit coordination. XFMD may progress consumer inspection/design now and integrate against the agreed actual producer capability later.

No existing owner pause is lifted. No main merge, release, force push or destructive cleanup is selected by this overview. The integrator must return a concrete reconciliation policy if existing schemas cannot preserve both histories; it must not silently rewrite immutable event bytes.

## Working agreement for future overview requests

Use the [agent prompts](Prompts.md) to establish scope and the common return format. Agents maintain their owning Sessions/plans/cards and return one candidate-bound report at a milestone or change in direction. This coordinator refreshes these files and Git facts on each requested “grand overview,” then records deltas, dependencies, conflicts and the next assignment. Chat memory alone is not the state store.

This coordination branch owns only this directory and Session0012. It does not own worker code, worker Session updates, card lifecycle or either canonical ledger. No agents were spawned or messaged. Prompt delivery and unrecorded chat decisions remain unknown until the owner forwards a prompt or the worker records a handoff. Reports can normally be read from disk without requiring the owner to copy their full contents.

## Evidence and source locations

The [snapshot](evidence/snapshot-2026-10-10.json) retains full heads/status, remote observations, merge forecasts and distinct colliding event records. The [initial snapshot](evidence/initial-snapshot-2026-10-10.json) preserves earlier observations before other agents committed their closeouts. Collection is sequential, not an atomic global snapshot. Freshness must be rechecked before execution.

Current PG Session:
file:///home/warloc/git/SDP-project-governance/SDP/Sessions/session-%230007--Project_governance.md

Current Actions Session and Blueprint Session:
file:///tmp/sdp-blueprint-implementation/SDP/Sessions/session-%230011--SDPTool_actions_and_JSON_API.md
file:///tmp/sdp-blueprint-implementation/SDP/Sessions/session-%230008--Semantic_blueprints.md

Project Leader role and concurrency procedure:
file:///tmp/sdp-blueprint-implementation/Skills/sdp-project-leader/SKILL.md
file:///tmp/sdp-blueprint-implementation/SDP/ProjectManagement/Concurrent-Work-Preflight.md

Release Session and evidence:
file:///tmp/sdp-runnable-programs/SDP/Sessions/session-%230010--SDUI_widgets.md
file:///tmp/sdp-runnable-programs/SDP/Maintenance/RSP2/Evidence.md

XFMD paused and current Sessions:
file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230005--SDL-viewpoint-diagrams.md
file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230006--Discoverable-SDPTool-actions.md
