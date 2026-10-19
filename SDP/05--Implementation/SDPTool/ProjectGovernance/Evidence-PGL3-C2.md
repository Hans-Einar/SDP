# PGL3-C2 — Automatic reports and local report inbox

2026-10-19, Session0007 T026 / PLAN-SDP-0018 / KB-SDP-052.
[Contract](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md),
[preflight](preflight-PGL3-C2.json), [candidate](candidate-PGL3-C2.json) and
[native trial](trial-PGL3-C2-native.json) identify this bounded delivery.

## Delivered behavior

New launches explicitly selected with `--reporting automatic` configure only the
bound `sdp_governance.run_transition` tool with the pinned host's `approve` policy.
The launch payload records the choice; default/manual launches retain the old policy.
Global configuration, sandbox, account, other tools and installed host are unchanged.
Reports/failures/scope notices remain bound to the current attempt; no start/resume,
owner design approval, evidence acceptance or publication capability is added.

Agent reports can cross receipt-proven tool/message observations while all decision
state stays unchanged. Every intervening revision must exist and add one permitted
observation in the same running attempt. Steering, approval requests, disconnects,
reconnects, lifecycle and authority changes are not ignored. Controller revisions
remain strict. A report enters verifying; completion remains unavailable.

The read-only inbox projects immutable report receipts, with stable operation IDs
and historical text/revision separate from current run state and coverage. It survives
restart, retains reports after cancellation and deduplicates exact retries. Empty
reads create no state. A changed operational generation during inspection refuses;
filesystem coverage is checked separately, not claimed as an atomic source snapshot.
Terminal rendering quotes control characters. No new persistence profile is introduced.

## Verification and failed experiment

The complete SDPTool test suite, affected race tests and vet pass. Tests exercise
explicit/default/invalid policy, unchanged approval arguments and exact per-tool
configuration; report delivery across observations; strict controller/stale authority,
steering, scope, approval, disconnect and non-report refusal; existing approval
accept/decline/cancel/replay cases; inbox restart, duplicate, cancellation and stale
coverage; missing-run/read-only behavior and terminal safety.

The [first native experiment](trial-PGL3-C2-failed.json) failed report delivery because
it used `auto`, which the pinned source implements as annotation-based classification.
The host correctly refused the write under never. Its operation was not replayed and
its result is not described as successful reporting. Source inspection distinguished
`approve`, the supported per-tool setting matching the owner's prior automatic-report
selection. A new project/run then tested the corrected final candidate.

The final native trial used the complete Codex 0.161.0 package and existing managed
account. It demonstrated the original Unicode-trimming defect, changed only the two
permitted Go files, passed project and external behavioral tests, and actually stored
one report without an observed approval request. The run is verifying with known
coverage; the completed host attempt is not acceptance. Stale-design launch, duplicate
activation and duplicate launch still refuse; an older reader still refuses profile
0.5. Scripted fixture-owner decisions are test authority, not actual owner acceptance.
The source/binary identities and native transcript hash are retained with the candidate.

The original 176 owner-lab project/runtime entries remain unchanged. A new read-only
`./inbox` helper selects the separate final preview, labels it synthetic and supports
`--json`. It starts no agents and makes no decision. The earlier ./implementation and
design-review previews retain their own candidate identity and historical results.

## Remaining boundaries

PGL3 remains in progress: structured attention questions/replies, leader consumption/
acknowledgement, responsive dialogue, phase continuation and native parallel evidence
are not delivered by this increment. The 1024-receipt store capacity remains and must
be addressed before claiming sustained multi-hour throughput. SAD dispatch, explicit
scope closure/replanning and retained evidence/review/Steering gates remain open.
No independent review, owner acceptance, merge or release is claimed.

Owner T026 assigns hub/web/mobile and notification delivery to a separate project.
KB054 retains the external handoff and local adapter relationship until a concrete
external destination exists. No hub/client implementation or deployment was added.
Management transitions are PM205 / KB378–379; system evidence is separately recorded
in Traceability. The owner's capacity/reset question is recorded in the Session;
it did not cause the native report failure and no current outage is inferred.
