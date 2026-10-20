# PGL3-C3 — Durable local attention questions and replies

2026-10-21, Session0007 T027 / PLAN-SDP-0018 / KB-SDP-052.
[Contract](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md),
[preflight](preflight-PGL3-C3.json), [candidate](candidate-PGL3-C3.json) and
[compiled CLI trial](trial-PGL3-C3.json) identify this bounded delivery.

## Delivered behavior

Running, current assignments can retain questions with a prompt, up to five options,
free-text input and a blocking flag. Each question carries an operation UUID, source
revision, run/attempt identity, server time and decision-basis digest. Bound agents
can ask/read via MCP; controller-only reply capture supplies actor/source separately
from request JSON. This does not authenticate the human behind a local shell account.
A response records input; it does not approve design, expand scope or start a host turn.

Current applicability ignores host message bookkeeping and completed-turn status,
while owner steering, source/authority changes, routine changes and observation gaps
invalidate it. A reply requires the exact current question basis. Questions, replies
and explicit withdrawals share the atomic operational store and immutable receipts.
Exact retries retain history, conflicting retries refuse, and an obsolete question
can be withdrawn with an explanation. Consumption remains explicitly not-tracked.

Unresolved blocking questions guard report/start/resume for their run. Other runs and
stop/fail/scope-change remain available. This guards governance advancement; arbitrary
host file/tool execution is not forcibly suspended. Agents are instructed to stop
dependent work and return a waiting summary instead of inventing an answer or polling
indefinitely. Automatic-mode launches preapprove only the existing transition tool
and the new bounded attention_request; all core authority guards remain.

New question writes explicitly select operational profile 0.6. Older readers refuse;
preparation, lifecycle, design-decision and activation writers preserve it. The paired
controller/MCP now requires seven bound tools. Unbound read-only discovery retains
three tools. attention_read reports current applicability and can select one question
by ID when a full history would exceed the transport's response limit.

## Evidence by level

- Service tests cover actual persistence/reopen, immutable question/reply retries,
  conflicting attribution, second-response refusal, empty read/no write, untrusted
  no-write refusal, invalid options/revisions, current-source checks, cancellation,
  withdrawal and sibling independence. Blocking/informational questions differ at
  resume, and a second authorized run progresses while the first waits.
- Controller tests cross real recorded tool/message observations, distinguish
  completion from steering/scope/disconnect staleness, and reject wrong run/epoch
  capabilities and agent-authored owner replies. Automatic policy remains explicit
  and restricted to two tools; prior approval guards remain covered.
- MCP protocol tests negotiate an in-memory transport with the actual server,
  inspect seven tool annotations, ask, receive a controller reply, and read current
  applicability. Unknown/foreign questions and agent reply-tool access refuse.
  Existing compiled stdio negotiation remains part of the full suite.
- The isolated trial invokes actual compiled CLI processes for ask/list/reply/withdraw,
  report blocking, other-run progress, EOF and invalid-choice refusal, captured option
  and answer, retry without input, attribution conflict and stale-source refusal.
  An old C2 controller refuses profile 0.6 without rewriting HEAD. Nine synthetic
  project files end byte-identical. One open question remains for owner experimentation.
- Actual attention JSON emitted by tests and the CLI is checked against the maintained
  wire definitions. Terminal rendering quotes project-controlled escape sequences.
  Full SDPTool tests, affected race tests and vet are recorded with exact output/hashes
  in the candidate. Documentation/management checks and ledger-prefix checks are
  separate evidence, not substitutes for runtime behavior.

The original owner-lab project/runtime entries remain unchanged; earlier previews
are retained. The new ./attention helper is outside those records and displays a
separate labelled synthetic preview. --reply --owner NAME captures input but launches
no agent. The maintained reproduction scripts and binary identities are included.

During verification, an added test incorrectly expected resume to enter running;
the existing routine deliberately returns ready. The assertion was corrected without
changing the routine. Wire validation also caught nil options serializing as null;
optional request options now omit the field, retaining the strict non-null contract.
The shared Result operation enumeration now includes the three attention operations.
Earlier failed local test logs are retained; final checks cover the corrected source.

## Limits and next work

No new native model task was submitted for C3. C2's native implementation/report trial
remains historical evidence for C2 only; native attention/reply continuation is still
required. There is no leader consumption acknowledgement, automatic next-turn dispatch,
responsive leader or native parallel/multi-hour evidence in this increment. The store
still bounds receipts at 1024, state at 16 MiB and retained questions at 32 per run;
transport responses are bounded at 256 KiB. Capacity/recovery must be addressed before
sustained unattended-use claims. Current filesystem coverage is checked separately
from the operational generation and is not an atomic source snapshot.

SAD dispatch, explicit scope closure/replanning, PGL4/5 and retained PGI gates remain.
Hub/webapp and notification delivery belong to the owner's separate project. No
independent review, actual owner approval, merge or release is claimed. PM206 / KB380
record management disposition; separate Traceability events link system evidence.
Origin main at 5ca0e90 is observed, not integrated: the branch shares e8760d3 with it.
Future integration still requires the ordinary conflict and candidate checks.
