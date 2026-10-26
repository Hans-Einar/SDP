# PGF4-C1 — Terminal Session operations and reply continuity

T048 adds a thin owner terminal adapter over the existing Session and attention core:
list/status, explicit record resume, selection, dependency progression, failed-result
rework, eligible completion, exact owner-gate decisions and current question/reply
capture. Canonical Session IDs resolve only when unambiguous. No JSON is required
for these operations. The adapter creates no second workflow store or host scheduler.

## Actual behavior and checks

Compiled process tests exercise prepared completion followed by a new-process resume,
unchanged old ledger bytes, duplicate read-only resume, stale revision/authority
refusal, exact owner accept/reject and wrong-gate refusal, rework with retained history,
dependency progression, broader-scope closure refusal, EOF/no-reply and actual stdin
question capture. Terminal rendering escapes untrusted text and propagates writer
errors. Focused race and complete affected adapter suites pass; candidate/check outputs
are recorded in [candidate](evidence-T048/owner-terminal-candidate.json).

[Initial native result](evidence-T048/native-attention-result.json) deliberately remains
failed: a real first attempt asked a blocking question and exited, the new terminal
command recorded its answer, and a fresh controller dispatched it to a continuation.
The existing attention projection then incorrectly described that selected answer as
stale. The native agent correctly stopped without editing. Synthetic tests had checked
delivery but not this observed applicability. No passing outcome is inferred from the
two completed host attempts or exit codes alone.

The correction preserves current applicability for an exact reply selected by the
latest claimed/acknowledged launch into the current gap-free attempt of a running
assignment, with current source coverage. Later steering, changed source/assignment,
wrong attempt/epoch/reply, uncertainty, withdrawal or terminal routine state cannot
inherit it. Admission/claim/retry and semantic-consumption boundaries are unchanged.
Service/controller regression tests exercise both the corrected case and invalidation.

[Final native result](evidence-T048/native-attention-final-result.json) passes on newly
built controller and MCP binaries with Codex 0.164.0: first attempt completed waiting
without changing target.txt; stdin reply was recorded; a separate controller continued
the same run in a second fresh attempt; only target.txt changed, to the exact requested
bytes; the assigned report left the run verifying. The answer is attributed to
fixture-owner and is test input, not actual owner feedback or acceptance. Original
failed observations are retained; no unknown send was retried to manufacture success.

## Design and delivery limits

The existing SessionWorkflow WORK :00002 supplies ExternalOwnerCli,
ObservationProjector and shared WorkflowCommands. The final
[scoped task](evidence-T048/owner-terminal-blueprint-task.json) generates diagnostic
Blueprint 44a838b0c1f3fce5c8729caa711044b9423c9f55228a01f077c9f8b828a4ace3
with passing constraints. This is unreleased TARGET, not model acceptance.

The native fixture has synthetic prerequisite history and real question/reply/attempts.
The result proves this bounded owner takeover/continuation path, not active-host crash
recovery, native Project Leader scheduling, the full no-design/SDL acceptance matrix,
independent PG product review or owner trial. Transfer/rebind, model lifecycle,
supervision, remaining usability and sustained capacity remain open in PLAN-SDP-0018.
