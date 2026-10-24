# PGF1-C2c2a — Checked selected-step disposition and claim release

Session0007 T047 continues PLAN-SDP-0018 / KB038 / KB052. This checkpoint implements
no-design step completion where owner acceptance is explicitly not required. Whole
cards can complete; contributions retain explicit remaining work. Plan/Session
closure, owner decisions, rework/dependencies and later leader/model gates remain.

## Candidate and design basis

[Candidate](evidence-PGF1-C2c2a/candidate.json) pins baseline
e377d88327842a083ceca304570bfe5600c47747 plus fifteen product/test/schema hashes and
the compiled executable. Checks apply to that uncommitted candidate; maintained
documents are finalized afterward. The full race candidate is pinned separately in
[race candidate](evidence-PGF1-C2c2a/race-candidate.json): all executable sources and
the binary match the final candidate. Final schema inspection made `remaining`
optional/empty for whole cards, matching the existing core; contribution remaining
work stays mandatory/nonempty. Six schema boundary cases and both lab captures pass
against that final schema. [Checks](evidence-PGF1-C2c2a/checks.json) records
commands, outcomes and limits. Verification is same-context, not independent review.

Before product edits, the preliminary scoped Blueprint was generated at revision
683b351e8bec404e426515fa158dab7eaccbdc8746469d430f9742cb8622d77b; constraints pass.
[Task](evidence-PGF1-C2c2a/blueprint-task.json),
[hashes](evidence-PGF1-C2c2a/blueprint-checks.json) and
[SDL validation](evidence-PGF1-C2c2a/sdl-check.json) retain the basis. WORK :00002
already defines CheckWholeCardDisposition; its sources remain unchanged TARGET.
ObservedBaseline is not accepted NOW and aggregate delivery/model release stays open.

## Proven behavior and limits

- Exact committed passing verification and required review precede disposition.
  The current Run revision, candidate, artifacts, skills and canonical inputs must
  still match. Required owner acceptance refuses; no approval flag bypasses it.
- Ordered per-card outcomes must match whole/contribution binding coverage. Whole
  completion refuses when another bound step retains that card. Contributions stay
  active/in-progress with explicit remaining work; eligible whole cards move to
  completed with metadata, supported selected links and board index updates.
- Preparation does not publish records or complete the Run. Resume publishes an
  immutable disposition artifact, journals, append-only management history and a
  separate Traceability evidence event. Plan and Session remain active. Removals
  occur after creates/updates; scope-claim removal is last. Operational commitment
  then completes the Run and exposes the selected step as done.
- The shared transaction reservation survives until operational commitment. Current
  legacy admission observes that reservation, so partial claim removal cannot expose
  the scope as committed. Exact pre/postimages allow recovery at every effect; old
  retries cannot duplicate events. Manual writers and old binaries remain outside
  cooperative admission guarantees.
- Tests cover whole/contribution outcomes, not-required review, required owner holds,
  wrong/stale identities, missing coverage, overlapping bound steps, changed candidate/
  test/skill/card, occupied destinations, outside references, competing claims and
  stored tampering. Lost preparation replies, transaction reservation, every canonical
  fixture effect and both sides of operational head publication are exercised.
  Mid-publication candidate/card/destination/trace edits survive and keep work pending.
- Operational profile 0.18 validates retained dispositions with the same deterministic
  transformation used to prepare them. The [new executable reads the existing 0.17
  store](evidence-PGF1-C2c2a/predecessor-read.json) without mutation, and the [old reader
  refuses 0.18](evidence-PGF1-C2c2a/old-reader.json) without changing its HEAD.

The initial second-step test fixture omitted that step's Session heading and correctly
failed admission before reaching the disposition test. The fixture now supplies the
authored heading; the overlapping whole-card guard then refuses as intended. No
product gate or test obligation was weakened. Inspection also found an unchecked
history-payload type assertion in retained-state validation. The final candidate
returns a controlled error and tests that malformed payload. The first race run was
interrupted (exit 143); the full suite was restarted on the corrected candidate.
The expanded recovery matrix uses a
15-minute package timeout for the full race suite; prior governance race runs already
took several minutes. Actual outcomes, rather than timeout settings, govern delivery.

## Compiled pg_lab journey

/home/warloc/git/pg_lab/session-c2c2a-disposition uses actual compiled processes for
CLI preparation/read/resume/retry/list and trusted JSON invocation parity. Nine wire
objects validate. Instance revision 16 reaches done; its Run is completed. KB-SDP-001
remains active/in-progress as a contribution; KB-SDP-002 moves to completed. No scope
claim or transaction reservation remains. All 264 earlier tracked paths are unchanged.
The fixture is committed on test/pgf1-c2c2a-disposition; [lab checks](evidence-PGF1-C2c2a/lab-checks.json)
retain the exact commit. Private operational state remains outside Git. The initial journey used
[the preceding candidate](evidence-PGF1-C2c2a/lab-candidate.json) before the defensive
validation correction. The final executable reads/lists and retries both original
operations with identical receipts and unchanged operational HEAD/fixture bytes;
[four final wire objects](evidence-PGF1-C2c2a/final-lab-wire.json) validate. The final
full suite also runs the compiled journey against a fresh temporary fixture.

[Process output](evidence-PGF1-C2c2a/lab-output.json) and
[wire captures](evidence-PGF1-C2c2a/disposition-wire.json) retain the observations.
Agent reports, test evidence and review observations are synthetic service inputs.
They are not native agent execution, independent PG review or owner acceptance.
PG checks exact evidence identity/currentness and controller-attributed card scope;
semantic sufficiency still requires the applicable verifier/reviewer procedures.

## Remaining work

C2c2 still owns actual owner-decision/gate handling, failed-result rework, dependency
progression/replan and explicit plan/Session closure. PGF2/3 retain common privileged
leader operations and SDL WORK/model handoff/release. Native whole-workflow integration,
recovery/capacity evidence and independent PG review remain later gates. KB038/052
remain in-progress, KB050/Session0008 stays separately paused and MVP1/Ponsse remains
deferred. No merge, installation, software/model release or complete-PG claim.
