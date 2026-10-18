# ProjectGovernance — owner walkthrough and feedback

Snapshot: 2026-10-19, Session0007 T024. The development branch now includes
activation of an approved implementation package and a real isolated code trial; this is not a released feature set. This guide explains the existing
[orchestration design](Orchestration.md) and its delivered boundaries. It does not
introduce new approval requirements or record owner acceptance.

## What “owner approval bound to the design revision” means

There are two different conversations:

| Conversation | What the owner is deciding |
| --- | --- |
| Developing PG itself, as we are doing now | Whether this way of collaborating is understandable and matches how you want to work |
| Using PG later on a project task | Whether a particular proposed change is ready to be implemented within its stated scope |

The new local decision adapter concerns the second conversation. It retains an
explicit input and the exact proposal it concerned. The
previous progress report was not a request to approve the PG implementation or PR.

For example, suppose SAD proposes adding CSV export to an application. You inspect
the proposed behavior and say “Implement CSV export, but leave permissions unchanged.”
The intended PG workflow associates that decision with the reviewed proposal and
scope. If the proposal later changes to include permission changes, the previous
decision cannot silently authorize that different package. Relevant rework or a new
decision is needed before dependent implementation proceeds.

This is an illustrative complete workflow. The lab exercises decision capture and
change detection with scripted fixture input, not this complete project journey.
The readable proposal is what you review. Revision hashes are internal bookkeeping;
you should not have to inspect hashes or JSON to understand your decision.

## The intended collaboration

The role arrangement already comes from the
[owner's Project Leader/SAD handoff](Project-Leader-Handoff.md):

1. You discuss an outcome with Project Leader. The leader keeps the project overview
   and coordinates bounded assignments while remaining available to you.
2. SAD develops the proposed design using ModelGovernance. NOW is the starting model;
   TARGET is the proposed model. A Blueprint describes the model change and its task
   context. You can inspect the design in XFMD and talk directly with SAD.
3. You give feedback or approve the relevant design and scope. PG must retain that
   decision against the package you actually reviewed.
4. Master carries out the authorized implementation through its Session and plan,
   reporting results and phase evidence. Existing continuation authority may cover
   subsequent phases; this does not imply a new owner click for every tool call.
5. Significant design or scope changes are reconciled with your authority. You can
   talk directly with Master; relevant feedback must also reach Project Leader.

The responsive leader conversation, automatic delivery of feedback and multiple
running Masters are planned capabilities. The current native trial runs one bounded
implementation task; it does not demonstrate that coordination workflow.

## What works in the development version today

| Capability | Current behavior | Practical boundary |
| --- | --- | --- |
| Existing Maintenance pilot | A separate controller/lab route can run one bounded Maintenance task and expose supported steering and request decisions | This is the earlier pilot, not the new Project Leader/SAD orchestration |
| Registered-work overview | Shows stored runs, references, observations and unknowns | Stored observations are not proof that an agent is currently running |
| Prepare SAD/Master work | Retains the task context, selected inputs, writable scope and return criteria; refuses local scope conflicts | A prepared assignment is not a running agent; exclusion is local to its store |
| Revise or retire preparation | Revises the same assignment with retained history, or retires it and releases its unlaunched local claim | Retirement is not completion or cancellation of a running agent |
| Check the proposed design | Checks the assignment against the canonical Blueprint, actual NOW/TARGET contents, task and permitted paths | A successful check establishes a consistent review basis, not owner acceptance |
| Record a design decision for this new route | A local controller adapter displays the proposal and reads an explicit decision and comment; it retains earlier decisions and checks current applicability | Local attributed input only; conversation capture and the final UI remain unfinished |
| Activate approved Master implementation | Creates one ready run from the exact current approved design; the existing controller can execute it with bound MCP context | Native code trial passed; report submission was refused by the selected approval policy. SAD launch, leader conversation and inbox remain unfinished |

The technical design check can say **current** or **blocked**. “Current” means the
checked package still matches its referenced inputs. “Blocked” gives a reason, such
as changed TARGET contents. Neither result judges whether the proposed product is
useful or whether you want it built. Incomplete coverage and unknowns remain visible.

## What the current lab shows

The separate PGL2-C2b preview creates real ModelGovernance artifacts and a retained
Blueprint for a synthetic assignment. Its script checks an unchanged package, edits
TARGET, observes refusal, restores the exact sources and observes the same review
basis again. It deliberately leaves TARGET changed at the end.

Run `./design-check` from the existing PG lab directory to see that final blocked
result and its reason. `./design-check --json` exposes the technical response if
useful. This is a narrow demonstration of change detection, not an approval screen
or a complete application walkthrough. Your original lab project/runtime is preserved.

The [C2b evidence](../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL2-C2b.md)
links the actual captured requests/results and reproduction details.

A separate C2c preview now demonstrates decision records. Run `./design-review`
from the PG lab directory for its current state. It initially shows a technically
current design with a rejected decision from **fixture-owner**. This deliberately
shows that passing technical checks and accepting a proposal are different facts.
The scripted decision is not attributed to you and does not affect your original lab.

If you want to try the provisional interaction, `./design-review --decide --owner NAME`
displays this fixture's proposal and waits for a choice and comment. Use your chosen
label in place of NAME. `cancel` or ending input records nothing. This is optional;
you do not need to perform the trial before implementation can continue. No agent
launches. Your label identifies the input in the local record; it is not identity
authentication. An unrestricted process running as the same OS user is outside this
local trust boundary. Automatic capture of decisions from your normal conversation
has not been implemented.

The [C2c evidence](../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL2-C2c.md)
records these checks. The existing `design-check` remains technical-only; use
`design-review` when inspecting the recorded decision and whether it still applies.

### A real implementation trial

Run `./implementation` in the PG lab to inspect the separate C2d preview. It shows
one small Go task, the actual run and host states, changed files and behavioral test
results. The approved package was activated, a real Codex host changed only the two
permitted Go files, and external checks verified Unicode whitespace trimming while
preserving interior spacing and case. Your original project/runtime remain unchanged.

The synthetic fixture decision was not attributed to you. The host finished the
code task, but its MCP report was refused because that call requires approval and
the trial used `never`. The run therefore remains running, with a completed host
attempt and no recorded assignment report. This is an integration limitation, not
an accepted/completed assignment. The earlier read-only `design-review` remains a
separate rejected fixture; this trial does not overwrite your inputs there.

[Execution evidence](../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL2-C2d.md)
retains that distinction. The helper starts no work. Leader/report delivery, direct
SAD conversation and the full project panel are still subsequent plan work.

## Where your feedback is useful now

These are discussion points for shaping the next increment, not a pending approval
gate. No answer, documentation reading or successful technical check counts as approval.

| Topic | Concrete feedback to give | Current basis or suggestion |
| --- | --- | --- |
| What you need to see before deciding | Which information is missing before you can say “yes, implement this”? | Suggested review summary: requested outcome, before/after behavior, affected scope, excluded scope, unresolved questions and links to model views/previews |
| How you express the decision | Is a clear decision in your conversation with Project Leader/SAD sufficient, or would you also want a distinct confirmation control? | The role design permits direct conversation. The terminal input adapter is a first testable implementation; it does not select the final conversation/UI mechanism or adopt a mandatory confirmation button |
| What a decision covers | Should the displayed decision normally cover the whole bounded implementation assignment, or do you need finer checkpoints for a particular kind of work? | Existing authority can cover multiple phases. Changed design/scope and explicitly reserved decisions need reconciliation; routine execution need not repeatedly ask permission |

A useful response is a concrete example: “For this kind of change I want to see X,
then tell SAD Y, and expect Master to continue until Z.” Corrections to this guide
are equally useful. You do not need to approve the unfinished runtime or study the
implementation evidence to give that feedback.

## More detail when needed

- [Orchestration contract](Orchestration.md): exact behavior and delivered increments.
- [Implementation plan](../../../05--Implementation/SDPTool/ProjectGovernance/Plan.md): remaining milestones.
- [Session0007](../../../Sessions/session-%230007--Project_governance.md): owner direction, decisions and work history.
- [Draft PR 56](https://github.com/Hans-Einar/SDP/pull/56): current PG development delivery.

This Markdown document is the reading entry point for XFMD. No new XFMD approval
control or native PG dashboard is supplied by this documentation update.
