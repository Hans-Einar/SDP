# Bounded owner controller — PGI3 implementation candidate

This package coordinates the shared governance core and the pinned app-server
adapter for an **existing authorized Maintenance run**. It adds durable exact owner
input, one fresh Master turn, in-turn steering, interruption, scope-change stop and
read-only status. The PGI3 live-workflow milestone remains **in progress**: real task/steering/interruption and correlated recovery now have
evidence. Explicit per-request SDP MCP approval, executing-turn reattachment and explicit
skill-input correlation are implemented; native role/reviewer mapping remains open. The implemented controller executes directly without spawning children.

Build from SDPTool:

```sh
go build -o /tmp/sdp-controller ./cmd/sdp-controller
go build -o /tmp/sdp-mcp ./cmd/sdp-mcp
```

Use the [core CLI](../governance/README.md) to inspect coverage, explicitly bind the
checkout and create the selected run. Selection must reference an active Maintenance
plan, card, Session/step and authority; missing work records never authorize a launch.
Then invoke:

```sh
/tmp/sdp-controller run --project /absolute/project --state /private/state --run RUN_UUID --mcp-binary /tmp/sdp-mcp --message /path/owner-input.txt
```

The message file supplies exact owner text. The default sandbox is read-only.
`--sandbox workspace-write` is an explicit launcher choice for an already authorized
write assignment; inherited enclosing host restrictions still apply. The pilot
defaults to `--approvals never` and refuses host approval requests. An explicit
`--approvals ask` enables the bounded interaction below. It does not enable
network access, change account mode, sign in, choose a model or grant delegation.
Only observed ChatGPT account mode permits submission. This is not a claim that
account inspection proves model entitlement.

After launch, enter one JSON object per line:

```json
{"action":"status","text":""}
{"action":"steer","text":"Keep the existing scope and include the negative case."}
{"action":"interrupt","text":"Stop this attempt."}
{"action":"scope-change","text":"The requested change now exceeds the current assignment."}
```

Steer is for input within the existing authorization; it cannot change the stored
assignment or override its governing instructions. Scope-change records the new
input, changes the routine to replan-required and asks the host to interrupt.
An interrupt request is not proof that external effects have stopped. Observe the
host's outcome and reconcile committed effects before replanning. Natural-language
scope classification and the broader Project Leader conversation remain future work.

Output is newline-delimited JSON: bounded run summaries, completed visible messages
and explicit errors. Input EOF becomes observation-only; output closure cannot kill
the controller through SIGPIPE. Neither means cancel. Ctrl-C/deadline/controller
loss records uncertainty; a separate observer can inspect the durable state:

```sh
/tmp/sdp-controller status --project /absolute/project --state /private/state --run RUN_UUID
/tmp/sdp-controller recover --project /absolute/project --state /private/state --run RUN_UUID
```

Status is read-only and starts no host. When observed unanswered prompts exist,
it emits an additional `pendingApprovals` JSON object after the run summary. Treat
output as a JSON stream, including status. It reports the last recorded observation,
including lastObservedAt, not an asserted live process state. Recover reads mapped
host history, records the new connection epoch/gap and replays completed items
without duplicates. It never resumes or resends a model turn. For launch/steer with a lost acknowledgment, recovery requires a unique observed
userMessage.clientId and exact payload match inside the mapped thread. It records
the observed message before acknowledging the durable input. Missing/ambiguous or
unmapped history remains unresolved; no input is resent. Executing-turn reconnection
uses the separate attach path below. Native child/reviewer attribution remains
unverified. The default controller lifetime is 30 minutes (`--timeout`).

Current connections admit stable Codex 0.164.0 and compatible later releases.
Probe output and attempt observations retain the actual version. This is not a
claim that every future release or native workflow has been tested; unsupported
required wire methods/shapes refuse. See [PGC1](../../SDP/05--Implementation/SDPTool/ProjectGovernance/Codex-Compatibility.md).

## Durable dispatch contract

RecordInput is controller-only and uses existing immutable generations, receipts,
checkout identity, source revisions and operation CAS. Exact text and serialized
assignment payload are each bounded to 64 KiB; reject oversized context plus input
rather than truncate. Invalid UTF-8 refuses. At most 64 inputs per run are retained,
also subject to checkout receipt and 16 MiB snapshot limits. Run/MCP summaries omit
text and payload with contentOmitted=true; trusted OwnerInput reads recover them.
No generic input mutation tool is exposed to agents.

The sequence is prepared -> sending -> acknowledged or uncertain. Persist sending
**before** any host launch/submission, with a stable clientUserMessageId. Correlation
is not a host exactly-once guarantee. A crash after claim or ambiguous RPC outcome
never causes automatic resubmission. A fresh launch/steer refuses while any prior
input remains unresolved. Prepared records can be inspected, but this CLI does not
yet expose a deliberate discard/reconciliation action. Do not edit the state store
or create duplicate runs to bypass uncertainty.

A controller lease serializes operations on one run and a second lease covers that
run's terminal session, preventing takeover by another terminal. Distinct runs can
progress independently with distinct host connections. New per-run leases also hold
a shared lock on the legacy checkout lease, so an older exclusive controller cannot
race them. The short exclusive store mutation lock remains shared with MCP. The launch pins source references,
Session, assignment and sandbox in its payload; after thread binding it requires
the seven assignment-bound MCP tools before starting a turn. The fresh Master is
instructed to read governing documents and preserve bounds. No native child or
independent-review provenance is established by that instruction string.

Host completion records execution only. Agent report may move the routine to
verifying; completion, review acceptance and canonical Session/card publication
remain unavailable. Publication belongs to PGI4, integrated independent review to
PGI5. Same-user arbitrary shell access remains outside the API authority boundary.

## Verification

Core tests inject crashes before and after HEAD, preserve input across a fresh
Service, reject repeat claims/content changes and block dispatch after source drift.
Controller fixtures test launch-before-send persistence, steering, scope stop,
uncertain-send non-repetition, account refusal, leases and cross-thread events.
A compiled terminal test uses a **recorded fake app-server subprocess**, exercises
status/history recovery and a closed output observer, and preserves every fixture
project file. Protocol requests are captured and checked against the installed
0.160.1 default schemas. None of these invokes Codex or constitutes live execution.


## Measured PGI3-C2 boundary

Codex 0.161.0 executed the isolated Unicode-normalization Maintenance task, accepted
in-scope steering and left all excluded files unchanged. An actual output observer
closed without terminating execution; interrupt and new-process recovery preserved
the same run/Session. A later fresh attempt continues the same assignment.
A local transport proxy dropped the actual turn/start response after submission;
recovery matched the recorded clientId/text and acknowledged the one existing turn.
The proxy did not replace model execution with a fixture or resubmit the request.

The C2 MCP run_transition call returned an approval-required error under never.
C3 adds the explicitly selected interaction below. A model turn or tool grant still
does not complete the SDP routine or establish independent review.

## Explicit per-request approval

Add `--approvals ask` to a new `run` command. This selects stable host
`approvalPolicy=on-request` with `approvalsReviewer=user`, while the sandbox stays
as explicitly selected. The adapter supports only the measured SDP MCP
run_transition form; other servers/tools, URL/secret forms, command/file permission
requests and persistent grants are unsupported and refused. No experimental API,
auto-review child, account configuration or global policy is enabled.

A pending request produces an `approval` object containing requestId, epoch,
threadId, turnId, message and the exact transition arguments. Review those arguments
and copy its identities into one stdin response:

```json
{"action":"approval","approval":{"requestId":"0","epoch":"COPIED_EPOCH_UUID","threadId":"COPIED_THREAD","turnId":"COPIED_TURN","decision":"accept"}}
```

Choose `accept`, `decline` or `cancel`. Each answer applies once to that request;
there is no session/always choice. Input is recorded exactly before sending. No
answer, EOF or elapsed time grants anything. The input pipe must remain connected
to answer; an independent status reader cannot take over the controller lease.
Duplicate, cross-attempt, resolved and old-epoch replies refuse. Closing stdout
still does not kill execution, but work can naturally wait for an unanswered prompt.

A successful response write remains uncertain until serverRequest/resolved is
observed. That establishes closure, not tool success: the tool may still refuse
its arguments or state. Report succeeds only through the assignment-bound core.
Approval bookkeeping can advance revisions; only the exact accepted call can cross
exclusively observational/approval revisions with its decision state and coverage
intact. Steering, replanning, argument changes and external coverage changes still
refuse. Controller transitions retain strict CAS. Local observation/input bookkeeping
waits briefly for busy/stale pre-mutation conflicts (up to 100 attempts, 20 ms
delays). The MCP adapter retries busy only with its unchanged operation ID and
arguments. Host effects are never retried by this contention handling.

Transport loss never resends approval. History recovery does not reconstruct an
unobserved resolved notification; such a decision stays explicitly uncertain.
No owner acceptance or canonical Session/card publication follows from an approval.

## Reattach to an executing turn (Linux private socket profile)

The default stdio host is owned by the controller and ends with that process. For
work that must survive controller loss, explicitly run a compatible Codex >= 0.164.0 in a
separate terminal, using a private directory owned by your user:

```sh
mkdir -m 700 /absolute/private-host
codex app-server --listen unix:///absolute/private-host/host.sock
```

Start the existing covered run with its normal arguments plus
`--host-socket /absolute/private-host/host.sock`. The controller connects to that
server and verifies its reported version; it does not start or stop the server.
The measured host may create a direct alias into its private runtime directory.
The alias and target are checked, including private parent ownership, socket type,
peer credentials and process lifetime. TCP endpoints, indirect aliases, public
parents and replacement servers are refused. A same-user arbitrary process remains
outside the local API trust boundary.

Send `{"action":"detach"}` on stdin to record an observation gap and exit while
the server continues. Unexpected client loss can leave the last snapshot marked
running; it does not prove continued activity. A new terminal can rejoin the mapped
turn without submitting a new task:

```sh
sdp-controller attach --project /absolute/project --state /absolute/private-state \
  --run RUN_UUID --host-socket /absolute/private-host/host.sock --approvals ask
```

`attach` verifies the same server, current source coverage and mapped executing
turn, then subscribes using thread/resume with no configuration overrides. It
preserves the launch-time agent capability and creates a new observation epoch.
Owner steering, interruption and new per-request approval replies use that new
connection. No stale reply or uncertain send is replayed. `--approvals ask` opts
this connection into presenting supported requests; it does not change the thread's
host policy or inherit a grant. Select it explicitly again when needed.

Prepared/sending/uncertain inputs or an observed approval with unresolved closure
prevent attachment. Completed, unloaded or missing work requires `recover` for
history inspection; it is not silently restarted. A host process restart is a new
identity and cannot claim to preserve live execution. Recover retains its earlier
capability-revocation behavior. An attempt originally launched through stdio has no
socket identity and cannot be retrofitted into this path.

Running snapshots do not finalize incomplete items. Terminal history is replayed
for completed content, with conflicts refused and duplicates suppressed; the gap
stays visible. A late terminal transition during attachment is an explicit
reconciliation condition. The live profile demonstrates immediate reattachment,
not indefinite unattended execution, approval replay or survival of server failure.

## Launch provenance

Read a persisted launch snapshot without starting or connecting to Codex:

```sh
sdp-controller provenance --project /absolute/project --state /absolute/private-state --run RUN_UUID
```

The view identifies the run, attempt, thread, host version, capture/current epochs
and receive timestamp. It separates controller-assigned `Master` from the native
role/session/source fields actually returned by the host. Empty native values mean
no value was reported; they are not inferred from the assigned role. The host source
is retained as compact JSON data. In the measured app-server trial it is `"vscode"`;
that value alone does not prove a particular user interface or role hierarchy.

`instructionSources` contains paths the host reports as loaded at thread/start.
The governing sources/digests are the controller's separate local context snapshot;
they do not prove the exact instruction bytes the model consumed. Skill inventory
comes from `skills/list` with forceReload, resolved within the selected project.
It records enabled state and readable source-file hashes, with external entries
counted without reading their bodies. Unavailable metadata and discovery errors
remain explicit. The fixed `useEvidence: not-observed` must not be relabeled as
activation or compliance because a skill was discovered or the model said it used it.
Agent claims remain captured messages. Explicit skill inputs are an opt-in extension
described below; the availability snapshot retains its original meaning.

The view is an as-of snapshot, not a live cursor. After reconciliation/reattachment
it is labeled `historical-epoch`; historical attempts without a context observation
return unavailable. Reads do not change the run. Metadata does not permit an agent
to grant itself authority, publish records or satisfy independent review. Native
child provenance and procedural compliance remain unverified.


## Explicit project skill invocation

Select each project-local skill explicitly when launching an authorized run:

```sh
sdp-controller run --project /absolute/project --state /absolute/private-state --run RUN_UUID --mcp-binary /absolute/sdp-mcp --message /absolute/owner.txt --skill Skills/sdp/SKILL.md --skill Skills/sdp-master/SKILL.md
```

`--skill` is repeatable (maximum eight) and accepted only by `run`. Canonical source
paths and project-local discovery symlinks resolve through fresh host discovery.
Missing, disabled, unreadable, ambiguous, duplicate or external skills refuse before
an intent or thread is created. No skill is automatically selected from the role or
message text. Existing launches without the option retain their text-only behavior.

The durable input retains exact owner text separately from the outbound assignment,
which includes explicit skill-name markers. Selected names, absolute canonical
paths and local SHA-256 snapshots are recorded before submission. The controller
rechecks file bytes before claiming and immediately before turn/start; this is not
an atomic attestation of the bytes consumed by the host. A failure after claim stays
uncertain and is never resent automatically, even when the controller has not yet
submitted the turn. Ordinary recovery boundaries still apply.

`provenance.invocation` shows the latest explicit selection and input state.
Acceptance is unconfirmed, response-acknowledged or host-message-observed. The last
requires a unique client ID and an exact ordered content fingerprint. None proves
procedural compliance, which remains not-observed. The inventory's useEvidence field
continues to describe discovery only. Local changes after dispatch cannot rewrite
historical acceptance or establish current byte usage.

Recovery compares full input content, including each skill name/path and order;
matching text alone is insufficient. Missing/extra/reordered skills, unknown content
metadata, duplicate client messages and partial history refuse. The measured empty
text_elements field is normalized; nonempty values are not dropped. Steering remains
text-only. No native delegation, role authority or independent review is added.

For release-independent owner use, see the [PG lab](../labs/project-governance/README.md).

## Registered project work overview — PGL1

```sh
sdp-controller overview --project /absolute/project --state /private/state
sdp-controller overview --project /absolute/project --state /private/state --json
```

This read-only command starts no host and requires no run ID. It lists the runs
registered in the selected checkout store, including selected Session/step/card/plan,
task, coverage, routine state and last observed execution. It does not discover
other stores, unregistered tasks, all project cards or an active Project Leader.
Readable output is the default; JSON emits one
[`sdp-governance-overview/1` object](../governance/overview.schema.json).

The generation identifies one operational snapshot. Source coverage is checked
per run during inspection; it is not an atomic snapshot of project files. No state
is created for an unbound checkout. Invalid persisted state fails inspection.
Observed role, tool completions and explicit skill acceptance are historical facts;
current liveness, mapped children, compliance and owner acceptance are not observed.
Tool counts cover only retained events of the latest attempt and include failed or
declined completions; the last item's status is shown. No transcript or owner-input
content is included. Terminal output quotes text to neutralize control characters.
Multi-assignment launch, SAD and Project Leader turns remain later KB052 increments.

## Prepare SAD/Master work — PGL2-C1

```sh
sdp-controller assignment context --project ROOT --state DIR --input selection.json
sdp-controller assignment prepare --project ROOT --state DIR --input preparation.json
sdp-controller assignment inspect --project ROOT --state DIR
sdp-controller assignment inspect --project ROOT --state DIR --assignment UUID
```

`selection.json` contains role (`SAD` or `Master`) and the existing Selection fields.
Use explicit `sdptool ROOT governance bind --state DIR --json` to create the initial
binding, then context supplies coverage/input revision and expected generation.
The [preparation schema](../governance/preparation.schema.json) defines request and
single-object response shapes. Preparation changes operational state only; no host
is opened, no worker MCP privilege is added, no task is launched or accepted.

The operation ID must remain stable for an uncertain retry with identical content.
A successful retry returns the original immutable receipt; inspect to assess current
source freshness. Updated state requires a fresh context before a new operation.
Conflicting local scope or stale source bytes fail explicitly. Shared-worktree claims,
provider validation and actual design approval are separate unfinished gates.

The [separate preparation trial](../labs/project-governance/README.md#assignment-preparation-trial--pgl2-c1)
exercises these commands, including old-reader refusal after the explicit store
profile upgrade. It does not replace the original one-Master owner experiment.

## PGL3-C1 per-run isolation boundary

The control and terminal-session leases now have validated run scope while retaining
mutual exclusion with legacy checkout-wide controllers. A blocked host call for one
run no longer holds the other run's control lease. The same-run exclusion, immutable
input/retry checks and source/scope guards remain. File existence does not indicate
liveness; process exit releases OS ownership. Read-only observers take no lease.

Service/protocol tests exercise two disjoint runs while one host response waits,
legacy/new exclusion in both directions, failed-acquisition cleanup and subprocess
exit. These are not evidence of two real native agents, a responsive Project Leader,
a report inbox or phase continuation. Existing native reporting under never remains
refused; this change does not alter host policies or answer approval prompts.

## PGL3-C2 automatic reporting and report inbox

For a **new** authorized launch, add `--reporting automatic`. The default `manual`
retains the prior policy. Automatic mode sets only the thread-local MCP policy
`sdp_governance.tools.run_transition.approval_mode = "approve"` and, since PGL3-C3,
the same policy for `attention_request`; the launch payload
records the selected mode. `--approvals` and all other tools retain their previous
behavior. Attach/recover cannot retrofit this launch policy. The tool's capabilities
remain report/fail/scope-change for the current assignment, not approval or dispatch.

A current agent report can cross only receipt-proven tool/message observations
with unchanged decision state. Steering, approval requests, disconnect/reconnect,
lifecycle and source changes still invalidate it. This avoids making a completed
`run_read` observation invalidate its own next report. Controller CAS stays strict.

```sh
sdp-controller inbox --project ROOT --state DIR
sdp-controller inbox --project ROOT --state DIR --run UUID --json
```

The [inbox schema](../governance/inbox.schema.json) describes a read-only projection
of immutable report receipts. Operation IDs identify messages; reported revisions
and text remain historical even if a run later fails or is canceled. Current state
and source coverage are separate. Exact retries do not duplicate a report. Empty
inspection creates no store. The generation pins the operational snapshot; filesystem
coverage is checked separately and is not an atomic snapshot of source files.

This local view includes report text and should be handled as project data. Terminal
rendering quotes control characters. It does not claim author authentication beyond
existing local controller/capability binding, remote delivery, leader consumption,
acceptance or automatic continuation. The current store's 1024-receipt capacity still
applies; this is not evidence of multi-hour throughput. Structured attention/replies are described below. Hub/web/mobile implementation belongs to a separate
project under the owner's T026 boundary.


## PGL3-C5 report acknowledgement

After inspecting `inbox --json`, a trusted controller can record that one exact
report was processed:

```sh
sdp-controller inbox acknowledge --project ROOT --state DIR --input acknowledgement.json
```

The input follows `ReportAcknowledgementRequest` in
[the application schema](../governance/schema.json): `schema` is
`sdp-governance-report-acknowledgement/1`; copy `projectRef`, `checkoutId` and
`generation` (as `expectedGeneration`) from the inspected inbox. Supply a fresh
UUID `operationId`, the report's `runId`, `id` (as `reportId`), `reportedRevision`,
`inputRevision`, and `currentRevision` (as `expectedRevision`), plus bounded
`coordinator` and `note` text. The command emits the durable receipt as JSON.
Retry the identical input after a lost response; do not generate a new operation
ID for a retry. A changed generation requires a new inspection and decision.

Acknowledgement means recorded processing. It does not verify evidence, accept
work, publish Session/card records or launch another attempt. The run remains
verifying, with unchanged run revision. Workers cannot invoke this through MCP;
coordinator attribution comes from the trusted local controller, not authenticated
human identity. The same-user shell boundary is unchanged.

The first write selects operational profile 0.8. Older binaries refuse that store;
use the current binary to continue. Reads of stores without acknowledgements keep
inbox profile 1; profile 2 adds each report's optional `acknowledgementId` and
`acknowledgement`. Unacknowledged reports are displayed as pending. Current coverage
and current run state remain separate from the historical acknowledgement. Changed
source/authority, report decision context, run revision or store generation prevents
a first acknowledgement; an identical retry returns the original historical receipt.
An acknowledgement consumes one existing receipt slot; capacity is still bounded.


## PGL3-C3 durable attention

The bound agent may call `attention_request` with a current run/input revision,
operation UUID, prompt, up to five options and `blocking`. Use `attention_read`
with optional `questionId` to read the reply and current applicability. A question
may cross receipt-proven tool/message observations, with the same strict authority
checks as reporting. The seven-tool inventory is required at launch; pair this
controller with the matching MCP binary. Older five-tool binaries are refused.
No agent tool can answer on behalf of an owner, dispatch work or grant approval.

```sh
sdp-controller attention list --project ROOT --state DIR
sdp-controller attention list --project ROOT --state DIR --run UUID --json
sdp-controller attention request --project ROOT --state DIR --input QUESTION.json
sdp-controller attention reply --project ROOT --state DIR --input REPLY.json --owner NAME
sdp-controller attention withdraw --project ROOT --state DIR --input REPLY.json --reason TEXT
```

The [wire definitions](../governance/schema.json) include AttentionRequest,
AttentionDisposition and AttentionInbox. Request/disposition JSON never supplies
owner attribution. Reply separately captures an option ID (blank allows free text)
and required answer/comment from stdin, then rechecks the exact displayed basis.
EOF and invalid input write nothing. Exact operation retries recover the immutable
reply without prompting again; attribution conflicts refuse. This is local trusted
controller attribution, not authentication of the human behind a shell account.
Mutations return JSON on stdout; input prompts use stderr.

A blocking unresolved question prevents report/start/resume on that run. It does
not freeze unrelated runs, or claim to stop arbitrary host file/tool actions. The
agent is instructed to stop dependent work and return a waiting summary. New owner
steering, authority changes, lifecycle changes or an observation gap make a question
stale; ordinary host messages and turn completion do not. Explicit withdrawal retains
obsolete history. Answering is neither design approval nor expanded scope.

The first question write upgrades only the selected store to operational profile
0.6. Older readers refuse; preparation, design and activation writers preserve it.
Limits remain 32 retained questions per run, 4096 prompt/answer bytes, five options,
1024 operation receipts and 16 MiB state. MCP responses retain the 256 KiB bound;
use questionId for an individual record when the full history exceeds that bound.
Current source coverage is checked separately from the pinned store generation.
No leader consumption acknowledgement, automatic new host turn, multi-hour capacity
claim, hub or notification delivery is included. Native question/reply continuation
has not yet been verified; the new mailbox has service, MCP and CLI evidence.

## PGL3-C4 guarded attention continuation

Use `run ... --attention-continuations N` (0..8, default 0) to opt into a bounded
local watcher. After a clean completed waiting turn it checks stored replies once
per second without calling a model or appending idle receipts. When all blocking
questions are resolved and at least one current blocking question is answered, it
may start a fresh Master thread/attempt in the same authorized run. It stops at the
configured count or lifetime, on report/failed work, stale authority or unresolved
outbound delivery. This is not a cross-project scheduler or conversational leader.

A separate `continue` command handles an already completed waiting run after a reply:

```sh
sdp-controller continue --project ROOT --state DIR --run UUID --codex PINNED_HOST --mcp-binary MATCHING_MCP --sandbox workspace-write --approvals never --reporting automatic
```

Use exactly the original sandbox and approval/reporting policy. Original owner input
and skill selections are recovered from the recorded launch; --message and --skill
are not accepted for continue. A fresh attempt receives the original assignment,
current context, exact selected reply references and question/reply history. Previous
attempts remain historical. A reply does not approve design or change scope.

Launch preparation atomically binds the exact question/reply/basis references and
makes the run ready. Claim rechecks authority before host effects. Source changes,
open blocking questions, mismatched policy/skills, observation gaps, a non-completed
prior attempt or unresolved input refuse. Prepared/sending/uncertain inputs are never
automatically replayed. Existing history recovery may prove an exact accepted send;
it does not send work again. Prepared-but-unsent recovery is still a manual reconciliation
gap, not a claim of automatic crash recovery or exactly-once remote execution.

New continuation inputs select operational profile 0.7; all writers preserve it and
older tools refuse. Attention's consumption field now reports dispatch-prepared,
dispatch-sending, dispatch-uncertain or dispatch-acknowledged when selected. The last
means host acknowledgement, not semantic consumption or implementation correctness.
Original-question applicability may become historical after a fresh attempt starts;
its immutable reply references remain in the new launch input. Current launch/run
context still governs execution.

The CLI marks late events from known completed predecessor threads as ignored instead
of applying them to the fresh attempt. An expired wait retains the clean completed
attempt and question for a later explicit continuation. Loss during actual execution
still requires reconciliation. Waiting-controller settings are local and transient;
reply capture warns that an authorized waiting controller may continue after a reply.
No global host settings or default automatic behavior changes. Receipt, input, attempt,
payload and state bounds remain; sustained multi-hour execution is not established.


## Coordination console — PGL3-C6

```sh
sdp-controller coordinate --project ROOT --state DIR
sdp-controller coordinate --project ROOT --state DIR --json --limit 20
sdp-controller coordinate --project ROOT --state DIR --watch --json --coordinator project-leader --owner OWNER
```

The combined view shows registered work, questions and reports, prioritizing open
blocking questions and pending reports. Inspect subsequent pages with `--cursor`,
or narrow the view with `--run`. A cursor binds the current content revision and
filter; changed operational facts or displayed coverage require fresh inspection.
`generation` pins operational facts; filesystem coverage is checked separately.
Observed host status is historical, not a live heartbeat.

Watch polls once per second by default (`--interval` at least 100ms) until its
`--timeout` (default 30 minutes), interruption or quit. Unchanged snapshots are
coalesced. Stdin remains available for one bounded JSON command per line:

```json
{"action":"refresh"}
{"action":"quit"}
```

For acknowledgement, send `action: "acknowledge"` and `acknowledgement` containing
the exact `ReportAcknowledgementRequest` described above, with the coordinator ID
selected at launch. For a reply, send `action: "reply"`, `disposition` containing
an `AttentionDisposition` from the displayed question, and `reply` containing
`text` and optional `optionId`. The configured `--owner` supplies attribution; input
JSON cannot choose another actor/source. With `--run`, commands cannot affect a
different run. Omit coordinator/owner flags for a read-only observer.

JSON watch output uses `sdp-governance-coordination-event/1`: snapshot (`view`),
result (`result`) or error (`error`). Definitions live in
[the coordination schema](../governance/coordination.schema.json). Bad commands emit an error and
leave the loop available. EOF stops input but leaves observation active until its
bound; a read failure ends the input stream. Quit/timeout/output failure stops this
observer and does not terminate workers. No lease or state lock is held while
waiting. Replies and acknowledgements use the existing atomic store and are visible
after restart; snapshots alone do not create or consume records.

This is an intermediate local interface. It launches no model, interprets no
free-text approval and does not implement native Project Leader wake-up/turn
scheduling. Its optional owner reply is attributed input, not design approval.
