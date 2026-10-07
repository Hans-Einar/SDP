# Bounded owner controller — PGI3 implementation candidate

This package coordinates the shared governance core and the pinned app-server
adapter for an **existing authorized Maintenance run**. It adds durable exact owner
input, one fresh Master turn, in-turn steering, interruption, scope-change stop and
read-only status. The PGI3 live-workflow milestone remains **in progress**: real task/steering/interruption and correlated recovery now have
evidence, while interactive host approvals and native role/reviewer mapping remain
open. The implemented controller executes directly without spawning children.

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
uses approvalPolicy=never and refuses host approval requests. It does not enable
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
scope classification and the broader Steering/PM conversation remain future work.

Output is newline-delimited JSON: bounded run summaries, completed visible messages
and explicit errors. Input EOF becomes observation-only; output closure cannot kill
the controller through SIGPIPE. Neither means cancel. Ctrl-C/deadline/controller
loss records uncertainty; a separate observer can inspect the durable state:

```sh
/tmp/sdp-controller status --project /absolute/project --state /private/state --run RUN_UUID
/tmp/sdp-controller recover --project /absolute/project --state /private/state --run RUN_UUID
```

Status is read-only and starts no host. It reports the last recorded observation,
including lastObservedAt, not an asserted live process state. Recover reads mapped
host history, records the new connection epoch/gap and replays completed items
without duplicates. It never resumes or resends a model turn. For launch/steer with a lost acknowledgment, recovery requires a unique observed
userMessage.clientId and exact payload match inside the mapped thread. It records
the observed message before acknowledging the durable input. Missing/ambiguous or
unmapped history remains unresolved; no input is resent. Full executing-turn reconnection and
native child/reviewer attribution are required live PGI3/PGI5 evidence, not delivered
by this command. The default controller lifetime is 30 minutes (`--timeout`).

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

A controller lease serializes dispatch and a second lease covers the whole terminal
session, preventing takeover by another terminal. Core mutation locks remain short
so assignment-bound MCP calls can proceed. The launch pins source references,
Session, assignment and sandbox in its payload; after thread binding it requires
the five assignment-bound MCP tools before starting a turn. The fresh Master is
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

The real MCP run_transition call returned an approval-required error under never.
Do not change host policy automatically or route around the refusal. A completed
model turn still leaves the routine unaccepted. The proposed next client behavior
is explicit per-request owner interaction; it remains a separate selection until
recorded. See the PGI3-C2 evidence for the exact tested boundaries and remaining work.
