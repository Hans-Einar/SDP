# Interim governance agent prompts

Send each block to the indicated existing agent. These prompts have not been sent by the coordinator. The recommended first handoff is **A to the Actions/Blueprint agent**, making its integration assignment concrete. B and C establish the cooperating boundaries. D is a short release handoff after completed publication, not a request to repeat the release.

The latest overview is here:
file:///home/warloc/git/SDP-coordination/SDP/ProjectManagement/Coordination/Overview.md

## Common return format

Every prompt below requests a **coordination handoff** with:

1. As-of time, repository/worktree, branch, exact HEAD and dirty/staged state. Distinguish current candidate from tested candidate.
2. Owning Session, step, plan/milestone and card identity. Qualify any colliding identity with subject/path/source revision.
3. Delivered since the previous report, evidence/tests actually run, review disposition and remaining acceptance. Separate reported, inspected, tested, reviewed, accepted, merged, published and installed facts.
4. Next bounded action, dependencies and anything requiring a decision outside existing authority.
5. Exact files/shared contracts intended for the next writes; any new overlap or ID allocation conflict.
6. Durable handoff/evidence location and a concise response to return to the coordinator. Do not claim automatic delivery.

Use an existing Session evidence/handoff location. Update your own Session and canonical records; do not duplicate the coordinator's overview or another worker's milestone plan. Preserve records when another agent has advanced them. No live-agent state is inferred from a clean worktree, a running record or a stale Session.

## A Actions and Blueprint agent as bounded SDP integrator

Copy this assignment and include the common return format above:

> I appoint you as the bounded integrator for the current SDP workstreams, coordinated by the interim Project Leader conversation. Recover Session0011 T008 and PLAN-SDP-0024 ACT0, the current coordination overview and applicable SDP skills/instructions. This assignment authorizes a reviewable combined integration candidate in a new isolated worktree; it does not authorize a main merge, release, destructive cleanup, forced history rewrite or changing other agents' working copies.
>
> Refresh remote main, local worktrees, exact clean candidate heads, dirty/staged/untracked changes and current Session records. The inspected starting points were main 04f88ff, release 3475367, Blueprint/Actions 732e85c and PG 02e8c9d. Pin actual current inputs rather than assuming those hashes remain current. #53 is contained in #54; preserve the complete release lineage, including terminal/SIGTERM corrections and publication closeout.
>
> Own shared integration of SDPTool CLI dispatch, output arguments, presentation, Contract/README, dependency files and the combined model/management/traceability consistency. Suggested order: main plus release lineage, Blueprint/Actions role guidance, then PG. Adjust based on dependencies and record why. Preserve all three behaviors; choosing one side wholesale is not a resolution. Apply the freshest owning Session/card facts, especially PG's active KB038/052 and Actions' paused Session0008/current Session0011.
>
> Reconcile the two KB-SDP-051 subjects (runnable programs versus discoverable actions), the two MAINT-SDP-0015 subjects (RSP2 versus PLR1), and the 14 observed management event-ID collisions. Identify records by source revision/path/subject and preserve original immutable bytes/provenance. Deduplicate genuinely identical imported events, not different events sharing an ID. Do not silently renumber immutable history. Produce an explicit mapping/import proposal and validate it against the adopted schemas; if preserving both histories needs a new policy decision, return that concrete decision before the affected mutation and continue independent integration work.
>
> Do not import the whole dirty SDP-vNow tree. Much of it matches current main. Inventory residual deltas, obtain recorded ownership/evidence for the XFMD-driven producer/viewpoint work and any other distinct package, and include only bounded reviewed inputs. Keep an explicit disposition for excluded work. Never reset/stash another worker's tree.
>
> Coordinate the shared-file boundary through durable handoffs; the owner will convey this assignment to PG and XFMD. Do not claim a distributed lock or that other agents received it. An observed competing write requires a fresh comparison before dependent work.
>
> Verify the combined candidate with existing management/Toolkit validators and relevant ModelGovernance, Blueprint assignment, PG/controller/MCP and runnable-program/CLI tests. Use candidate-bound evidence and the existing required review process; do not call a self-check independent review. Keep published release artifacts immutable. Deliver integration commits/evidence, an unresolved-findings list, preserved/excluded work and a fresh ACT0 verdict. Resume Actions product implementation only when ACT0 is actually green and existing authority covers the next step. Session0008 remains paused.

## B ProjectGovernance agent

Copy this assignment and include the common return format above:

> Continue as the implementation owner of Session0007, PLAN-SDP-0018 and KB-SDP-052. The separate interim Project Leader conversation coordinates the overview; it does not replace your bounded Master role. Recover PGL1 evidence at 02e8c9d and refresh actual state before PGL2.
>
> The proposed integration assignment is to the Actions/Blueprint agent in its own integration worktree. Read its durable handoff before relying on that allocation; do not infer that a prompt has already been delivered. Supply a pinned PG candidate and the invariants that integration must preserve: authority/revision checks, correlated recovery/retry behavior, explicit skill provenance, registered-work observation limits, and the distinction between historical state and current host activity.
>
> Maintain your PG design/Session ownership. Your imported KB052/PLR1 records overlap the Actions branch; identify their original provenance and your later changes so the integrator preserves both without duplicate lifecycle events or stale KB038 state. Report every conflicting ID by source/path/subject.
>
> Continue authorized PGL2 contract/preparation work only after scoped concurrent-work preflight. Keep PG-specific changes isolated where possible. Shared CLI, output, presentation, dependency, model and history changes require the recorded integration boundary and fresh preflight; do not create a second competing integration. Do not expand into Actions or XFMD product implementation.
>
> Return exact PGL2 scope, intended writes, dependencies on ModelGovernance/Blueprints and the behavior actually available now. Preserve unfinished PGI/PGL review, recovery and owner-trial gates. Do not present PGL1 inspection as a completed multi-agent scheduler, or stored running state as verified live execution.

## C XFMD navigation and actions agent

Copy this assignment and include the common return format above:

> Recover Session0006 T001 and KB-XFMD-030, preserving the owner's pause of Session0005 T012 and installed XFMD 0.8. Continue S1 inspection and bounded consumer design under the existing authority. Read the current interim coordination overview and the actual Actions contract in the Blueprint/Actions worktree. Its executable catalogue/API is not yet delivered; SDP 2.2.0 program discovery/run is a distinct existing capability.
>
> Return a producer-consumer matrix: supported now, proposed, missing, and what can proceed without new producer behavior. Cover stable action identity, typed parameters/results, selected context, absent capability versus temporarily disabled action, JSON errors/revisions/retries, generic icons and retained personal toolbar configuration. Preserve KB029's unfinished scope; do not treat this Session as completion of all toolbar work.
>
> The Actions/Blueprint agent is proposed as the isolated SDP integrator. Identify the producer-side changes made for XFMD in the shared SDP-vNow checkout, including viewpoint overview/detail/source contracts and their exact paths/evidence. Provide a bounded candidate/patch manifest with source hashes and required consumer pairing; do not stage or sweep unrelated SDP changes. Explicitly identify any producer files whose ownership you cannot establish.
>
> Your XFMD checkout has substantial uncommitted work and existing staged card changes. Preserve it and its index. Record logical delivery boundaries and a safe future commit/branch plan consistent with existing phase commitments before isolating new implementation. Do not switch/reset the shared tree or claim HEAD alone captures installed 0.8.
>
> Select a proportionate consumer plan after actual contract inspection. Native menu/toolbar implementation that depends on new Actions support waits for the agreed producer baseline. Separately state whether XFMD Run still launches the SDUI prototype host or now consumes SDPTool's declared connected-program run; the installed gh-sdp upgrade alone does not prove UI parity. Do not resume paused diagram work or install a new version as part of this handoff.

## D Release and widget lab agent

Copy this assignment and include the common return format above:

> Recover the completed Session0010 T011/T012, RSP2-M2 and gh-sdp SPS-009. Supply a compact integration handoff for the interim Project Leader and Actions/Blueprint integrator. Do not repeat publication or installation: SDP 2.2.0, gh-sdp 0.2.2 and the extension upgrade were observed complete.
>
> Pin the release/tag source, branch closeout head, signed asset/publication verification, installed engine/client and lab/native evidence. Explain that #53 is contained in #54, and distinguish published source commits from later documentation commits. Preserve terminal input and SIGTERM lifecycle fixes in the eventual integration.
>
> Identify your KB-SDP-051 Runnable SDUI programs and MAINT-SDP-0015 RSP2 records by path/subject/revision because the Actions/role branch uses those IDs for different subjects. Supply source event provenance and the exact files copied into SDP-vNow; do not renumber historical records or resolve other workstreams yourself.
>
> State remaining PR/main integration and XFMD Run adoption boundaries. Report any still-uncommitted owned files; otherwise mark this release assignment complete at its authorized scope. Do not treat the open draft PRs as an uncompleted publication.

## Short prompt for later status refresh

> Give a coordination handoff for your current assignment using the common return format in the interim governance prompt document. Refresh actual Git and Session state first. Update the owning Session with your current outcome and next step. Report only changes since your previous handoff, remaining dependencies and intended shared-file writes; keep candidate-bound evidence and owner decisions distinct. Do not expand implementation scope merely to prepare this report.
