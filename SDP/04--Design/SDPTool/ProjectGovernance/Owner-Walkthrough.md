# ProjectGovernance — owner walkthrough and feedback

Snapshot: 2026-10-21, Session0007 T028. The development branch now includes
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
| Automatic reports and local report inbox | A fresh automatic-mode trial stores the report without a separate prompt; the inbox shows submitted text and current coverage | Reports await verification; consumption and automatic phase continuation are not yet delivered |
| Continue after an attention reply | A current answer can start a fresh attempt in the same assignment, explicitly or through a bounded opt-in watcher | Original scope, policies and skills remain; uncertain delivery blocks retries. This is not the responsive Project Leader |
| Activate approved Master implementation | Creates one ready run from the exact current approved design; the existing controller can execute it with bound MCP context | Native code and automatic-report trial passed with explicit automatic mode. SAD launch and leader conversation remain unfinished |

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

## Automatic work and calling for attention

Owner direction in T025 is now explicit: status/result reporting and work within an
existing mandate should proceed automatically. A report is not acceptance of the
reported result. A genuine missing decision should appear as an attention item with
context, options/free text and a clear statement of what is waiting. Independent
authorized work should continue; absence of a reply is not a decision.

The requested next extension is a shared inbox/overview for projects on multiple
PCs, accessible through a self-hosted web service and an installable mobile/desktop
client. [KB054](../../../KanBan/backlog/%23054--Proposal--PG-attention-and-remote-coordination.md)
and its [architecture proposal](../../../03--Architecture/ProjectGovernance/Attention-and-Remote-Coordination.md)
describe local execution, a proposed hub and a responsive PWA. ChatGPT MCP Events
is an optional documented possibility requiring an actual integration trial. These
are design directions, not features of the current lab. SDP_Analyzer is not reused.
The local report/inbox path comes first under KB052; the owner assigns the hub and
web application to a separate project. The broader hub/client has a
separate follow-on scope. Existing evidence/review gates remain.

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


## Inspect the automatic report — T026

From /home/warloc/git/pg_lab run `./inbox` for a readable report or `./inbox --json`
for the machine view. This reads a separate synthetic preview and starts no agent.
The final native trial changed the two permitted Go files, passed external tests,
and submitted one report automatically. The assignment is now verifying; that is
not owner acceptance or completed project work. The earlier ./implementation and
./design-review helpers retain their earlier previews and may show different states.

This is the first local report inbox, not the final communication interface. Local attention questions and replies are now available in the separate preview below;
leader acknowledgement and automatic continuation remain next. Hub/webapp and notifications belong to the separate
project confirmed by the owner. Feedback on report readability and the distinction
between submitted/verified/accepted is useful, but no new approval is required to
continue the already authorized local development.


## Try a local question and reply — T027

From /home/warloc/git/pg_lab:

```sh
./attention
./attention --reply --owner warloc
./attention
```

The first command shows three synthetic questions: one answered, one withdrawn and
one open. The second selects the open question, displays it again and asks for an
option ID (short/detailed, or blank for free text), then your answer/comment. The last
command shows the stored reply. `./attention --json` exposes the same records to a
local client. No agent is started and no production design approval is recorded.
Earlier previews and the original project/runtime retain their own states.

This demonstrates the local mailbox, including answering later after process restart.
A blocking question guards only the affected run; other runs can progress. A reply
must still match the question's assignment and authority. If these have changed,
the controller refuses the reply; an obsolete question can be explicitly withdrawn.
The next delivery must let the leader consume the reply and continue the appropriate
work. That automatic continuation is not delivered by this preview.

Useful feedback: can you see what is waiting, why an answer is needed and what it
applies to? Is choosing an option plus adding a comment natural, or should an option
alone be sufficient? This terminal interaction is a provisional local client; the
separate hub/webapp can later present the same concepts differently. You need not
answer this demonstration question for authorized development to continue.


## Inspect a real continuation — T028

From /home/warloc/git/pg_lab, run `./continuation`. This read-only view shows the
new isolated native trial: the first agent asked which label prefix to use and
ended its turn; a separate scripted fixture reply selected the short prefix; the
controller started one new attempt in the same assignment; that agent changed the
two permitted Go files, passed tests and submitted its report. The run is verifying.
The reply is test input, not your acceptance. This command starts no new model task.

The earlier `./attention` preview remains separate, including any answer you gave.
To reproduce the complete new trial with a development build, follow the maintained
[lab procedure](../../../../SDPTool/labs/project-governance/README.md#native-attention-and-continuation--pgl3-c4).
No sdptool/gh extension release or installation change is needed.

In ordinary use, the controller can wait with `--attention-continuations 1` (up to
8, default 0). You answer the stored question through the local attention CLI; a
waiting controller checks that the answer still applies before continuing. Without
a watcher, `sdp-controller continue` selects the current answer explicitly. Original
sandbox, reporting and approval settings must match. This is bounded continuation,
not a multi-hour scheduler or a complete leader conversation.

A delivered answer displays dispatch-acknowledged. The old question basis becomes
historical after the fresh attempt begins; the exact question/reply references are
retained in the launch. That is delivery evidence, not proof of semantic consumption
or approval. The native trial separately checks the resulting code behavior.

Useful feedback: is it clear which assignment a question belongs to, that answering
may continue already-authorized work, and which result is awaiting review? You can
inspect this while implementation continues; no new approval of PG itself is required.
