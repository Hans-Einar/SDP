# PGI3-C3 — explicit owner approval for the bounded MCP report

2026-10-08 (Europe/Oslo). PLAN-SDP-0018, KB-SDP-038, Session0007 T009.
This is a bounded PGI3 checkpoint, not full milestone or product acceptance.

## Candidate and authority

Base 37a09a6 is C2. The owner continued after the explicit per-request interaction
recommendation. Work remains on sdp/project-governance-pilot in the persistent
ProjectGovernance worktree; shared ModelGovernance/blueprint/SDUI work is untouched.
No merge, push, publication, native delegation or independent review is selected.

C3 adds opt-in `--approvals ask`; default never remains. The current controller
presents the exact observed SDP MCP transition, accepts a request-bound
accept/decline/cancel response and records exact stdin before sending. Scope is one
request in the current thread/turn/epoch. No persistent grant, general permission
handler or automatic acceptance is provided. Observed closure remains separate
from successful tool execution and SDP outcome acceptance.

The [design refinement](../../../04--Design/SDPTool/ProjectGovernance/Design.md)
records why recording the approval itself must not stale the exact queued report.
Every intervening revision is checked; controller CAS, changed arguments, steering,
replanning, reconnect and external coverage changes retain their refusal behavior.
The [candidate/evidence record](live-PGI3-C3.json) identifies final source hashes,
binaries, protocol captures, exact scenario drivers and actual outcomes.

## Protocol and measured corrections

Codex 0.161.0 was inspected using existing managed authentication; no API key,
account mutation or raw auth response is retained. The actual report request is
mcpServer/elicitation/request, form mode, empty object schema, with the exact
transition in tool_params metadata. It is not the generic question method.
[Official app-server documentation](https://learn.chatgpt.com/docs/app-server)
describes elicitation responses and the distinction between replies and resolution.
The adapter supports only the measured bounded SDP form.

Initial actual launch refused granular policy with `askForApproval.granular
requires experimentalApi capability`, despite that shape appearing in the default
schema export. The final explicit ask path uses stable on-request with user review;
unsupported approval methods still receive an error. There is no experimental opt-in
or silent fallback. The failed launch remains an uncertain historical attempt.

The next actual trial observed the prompt and advanced the MCP report to verifying,
but concurrent core bookkeeping returned busy and ended the controller. This is
recorded as a failed journey, not success. The correction retries only pre-mutation
bookkeeping contention, bounded to 100 attempts with 20 ms delays. MCP retains the
exact operation ID/arguments while retrying busy; neither path repeats a host send.
A deterministic real file-lock test checks one reply despite contention, and an
MCP transport test checks the same operation receipt after a blocked mutation.

## Actual final-binary workflow

The managed-account trial created one attempt for run
`e34b852e-1609-4636-a1e9-54317ad10254`. The model fixed Unicode boundary trimming
and added tests. The driver inspected a separate status view, then sent one exact
accept response for the observed report request. The host emitted resolved and
recorded exactly one successful run_transition with no error. The routine reached
verifying; it was not completed or owner-accepted. The host turn completed and the
approval input became acknowledged only after observed closure.

All declared exclusions retained their preimage hashes. Only normalize.go and
normalize_test.go changed, and a separate same-context driver test rerun passed.
The frozen driver records are test scenarios under PGI3 verification authority,
not verbatim owner interaction or independent review. No other route forced the
report through the core after a refusal.

## Verification boundary

Final affected-package race tests, compiled core CLI, vet, 28 captured application
cases, 11 host protocol cases and 24 actual live wire records pass. Management
validation (61 cards, 38 records, 513 events), Toolkit validation, source hashes,
ledger prefixes and changed-document links pass. The accompanying JSON retains
outputs and identities. The broad
regression run began before the final stable-policy/contention corrections; those
historical passing results are not relabeled as exact final-source results. Final affected
core/controller/host/MCP race tests, compiled core CLI checks, vet, protocol schemas
and the actual final-binary workflow establish the correction's candidate boundary.
Management/Toolkit validators, Markdown links, ledger prefixes and source hashes
are checked separately.

Negative tests cover no answer, decline/cancel, wrong request/thread/epoch, resolved
or duplicate response, changed arguments, steering, scope correction, external
coverage drift and lost writes. Unsupported servers/tools/forms/owner actions are
refused. The compiled terminal test exercises an independent status observer and
actual stdin response capture, using a fake host explicitly labeled as such.

## Remaining work

Executing-turn reattachment and applicable native role/provenance evidence remain
PGI3 obligations. Approval closure lost with transport remains uncertain; history
recovery does not reconstruct a notification it did not observe. The controller
input pipe must remain attached for responses. PGI4 still owns checked publication
and evidence acceptance; PGI5 owns independent integrated verification.

The earlier private C1 event-ID collisions with shared history still need deliberate
reconciliation before integration. Neither ledger history was rewritten, and this
checkpoint makes no merge/readiness claim for the shared branch.
