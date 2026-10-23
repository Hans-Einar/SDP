# PGL3-C4 — Guarded attention continuation

2026-10-21, Session0007 T028 / PLAN-SDP-0018 / KB-SDP-052.
[Contract](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md),
[preflight](preflight-PGL3-C4.json), [candidate](candidate-PGL3-C4.json) and
[native trial](trial-PGL3-C4-native.json) identify the bounded delivery.

## Delivered behavior

A clean completed waiting attempt can continue through a new Master thread/attempt
in the same run after a current blocking question is answered and all blocking
questions are resolved. The original owner input, assignment, selected skill bytes,
sandbox, reporting and approval modes remain unchanged. The payload includes question
history and exact selected question/reply-operation/basis references. This reuses
existing launch prepare/claim/ack and host methods; it does not reinterpret ended
turns as resumable live turns or grant native delegation.

Preparation durably stores the references and readies the run in one generation.
Claim reconstructs and checks the preceding question basis and current authority
before any launch effect. Duplicate use refuses. Prepared, sending and uncertain
inputs block a blind retry; existing history reconciliation can establish observed
acceptance but is not permission to resend. Inspection projects dispatch-prepared,
dispatch-sending, dispatch-uncertain or dispatch-acknowledged from the launch input.
Acknowledgement establishes transport delivery, not semantic consumption or approval.

The explicit continue command and --attention-continuations N share this path. The
watcher is opt-in (default 0, maximum 8) and bounded by controller lifetime. Idle polls
are read-only and submit no model turn. A clean wait survives timeout/transport close
without inventing an observation gap. Late events from completed predecessor threads
do not control the fresh attempt. A report ends the workflow in verifying; it does
not automatically authorize a new phase. CLI answer prompts explain that an authorized
watcher may continue after the reply.

Operational profile 0.7 retains reply references. All relevant writers preserve it;
older readers refuse without rewriting state. Original profile 0.6 attention-only
stores remain supported. Existing limits still apply: 1024 receipts, 16 MiB state,
32 questions per run, 64 retained owner inputs and 64 KiB launch payload. Large reply
histories may therefore refuse safely rather than establish unlimited throughput.

## Evidence and corrections

Core/controller tests exercise persistence/reopen, exact current references, immutable
original text/skills, changed authority/policy/scope, observation gaps, claim-time source
changes, uncertain StartTurn delivery and duplicate prevention. Profile preservation
covers preparation, lifecycle, design decisions, activation and later attention writes.
A compiled CLI test continues a stored answer after controller restart, then verifies
that repeated continuation refuses without changing HEAD or the completed predecessor.

The first native harness collided with a concurrent read-only generation change and
cancelled its controller. Its [failed record](trial-PGL3-C4-failed.json) is retained.
That historical record's modelTaskSubmitted=true came from attempt presence, not an
observed turn ID; it is not proof of a model turn. The send was uncertain and was
never retried. The harness now retries only read-only stale-generation inspection.
A subsequent isolated native trial passed before the final CLI refusal correction;
its [record](trial-PGL3-C4-intermediate.json) is retained separately.

Inspection then found that a refused explicit continue inherited generic disconnect
handling, which could change a clean completed predecessor. The command now preserves
preflight refusal state; post-claim uncertainty remains the launcher's responsibility.
The compiled regression test checks exact HEAD equality. Final full tests and affected
race/vet checks cover this correction and the final CLI prompt wording.

The final native trial uses complete Codex 0.161.0 and the final compiled binaries.
It demonstrates a real agent asking a blocking prefix question and ending its turn,
unchanged project bytes before the answer, no waiting-state writes, a delayed reply
captured by a separate process, exactly one fresh attempt, and an automatic report.
External tests fail before implementation and pass after it; project tests pass.
Only label.go and label_test.go change; nine protected files remain unchanged. The
question is answered and dispatch-acknowledged, and the run remains verifying.
The old question's applicability becomes historical after the new attempt; its exact
reply correlation remains retained in the dispatch input. The C3 reader refuses the
new store without mutation. Actual JSON wire objects are schema-validated separately.

All 176 original owner-lab entries remain unchanged; earlier previews are retained.
The separate ./continuation view reads current attempts, the answer, dispatch and
result without starting a model. Its scripted fixture-owner answer is not actual
owner acceptance. Candidate hashes identify source, binaries, transcript and logs;
management/document/ledger validation is separate from runtime evidence.

## Remaining work

PGL3 and Session0007 remain active. Report consumption and responsive Project Leader
coordination, native parallel execution, SAD dispatch/scope closure, PGL4/5 and retained
PGI gates remain. Capacity/recovery must precede sustained multi-hour claims. Current
filesystem coverage is checked separately from the operational generation and is not
an atomic filesystem snapshot. No independent review or actual owner acceptance is
claimed. PM207 / KB381 record management disposition, with separate implementation/
verification Traceability entries. Hub/webapp remain a separate project.

Main 5ca0e90 is observed, not integrated; the shared ancestor remains e8760d3. This
milestone does not claim a conflict-free later merge, a merge or a release.
