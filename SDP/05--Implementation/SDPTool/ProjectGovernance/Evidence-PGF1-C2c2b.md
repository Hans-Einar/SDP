# PGF1-C2c2b — Owner decisions and retained Session continuation

Session0007 T048 continues PLAN-SDP-0018 / KB038 / KB052 under the owner's request
to work independently through the remaining plan. This increment implements trusted
owner gates, dependency progression, failed-result rework and eligible whole-plan/
Session closure. It does not complete PGF1 or the full PG plan: remaining-scope
transfer/rebind, leader scheduling, SDL integration and the later gates stay open.

## Candidate and evidence

The [staged candidate](evidence-T048/candidate.json) pins baseline
8ca2d6238e0e1b866e65cf6f679c31483cfd8b04, executable/test sources, schema and compiled
binary. [Scoped task](evidence-T048/blueprint-task.json),
[Blueprint hashes](evidence-T048/blueprint-checks.json) and
[SDL check](evidence-T048/sdl-check.json) retain the preimplementation basis.
Blueprint 0b4f7e934b32056f2cbeafa05228a4364acba636654c0b732e747e664b178550 is a
preliminary diagnostic preview with passing constraints. WORK :00002 remains TARGET;
no accepted NOW, independent model disposition or release is claimed.

Focused owner/continuation tests, complete next-step and rework journeys, compiled
CLI/JSON schema checks, vet and SDL validation pass. The full SDPTool race suite passes (governance package: 1159.244 seconds).
[Checks](evidence-T048/checks.json) retain raw outputs and the explicit build boundary;
[race source](evidence-T048/race-candidate.json) differs only by supplemental advance
tests from the staged C2c2b sources. Those extra disjoint-card and external-edit/
remaining-scope tests pass targeted race checks. Separately developed MCP additions
have separate evidence; no exact-candidate whole-PG acceptance is inferred.
Verification is same-context.

## Behavior established

- Controller requests concrete owner review only for a current exact result. Cards
  enter gate-review when the request effects commit. Only a trusted local owner can
  accept or request changes; invocation authority propagates into the shared service.
  A controller JSON action cannot impersonate the owner through the CLI adapter.
- Decisions return selected cards to in-progress. Acceptance authorizes disposition
  only for its exact result/run revision; rejection requires fresh implementation.
  Owner requests and decisions use immutable artifacts, journals and management events.
- A next step must be present, unfinished and have completed dependencies. Original
  selection, assessment, work, owner decisions, result and disposition are retained.
  The new execution requires fresh selection, assessment, run, verification and review.
  A whole card can complete after its earlier bound contributions have been disposed.
- Rework requires an exact failed test, completed nonpassing review or committed owner
  changes request. Failed evidence stays pinned; earlier runs are retired after effects
  commit. Active/unknown host work and active review prevent the transition. Fresh work
  cannot inherit a prior owner acceptance or passing result.
- Plan activation makes all bound execution cards ready, including future steps;
  only the selected step acquires its claim. Overall closure requires all bound steps
  and card obligations completed. Broader contribution scope still refuses without
  explicit successor-transfer support. Plan and Session completion are separate from
  selected-step/Run completion.
- Owner and continuation effects use the same pre/postimage publisher and shared
  reservation. Recovery tests interrupt reservation, every canonical fixture effect,
  and both sides of operational publication, then retry without duplicate results or
  leaked reservations. Retained-state tests reject altered basis, pre/postimages,
  coverage, digest and receipts. Historical run and retry identities remain accessible.
- Profiles 0.19/0.20 retain owner records and flat execution history. New assessments
  preserve the newer profile. [Prior 0.18 read](evidence-T048/predecessor-read.json)
  succeeds without changing HEAD; [old reader](evidence-T048/old-reader.json) refuses
  the new store without changing it. Capacity remains bounded, not solved by history.

Initial testing found an assessment writer downgrading the new history profile and
synthetic host thread reuse across runs; both were corrected. The first owner-suite
parallelization conflicted with fixtures using process environment isolation and was
removed. These failures are not counted as passes or evidence of native behavior.

## Durable compiled lab

Two new isolated fixtures preserve earlier pg_lab records:

- `session-c2c2b-owner`: actual compiled CLI and trusted JSON invocation prove owner
  request/decision/disposition, controller refusal, replay parity and ten schema-valid
  wire objects. One contribution stays active; the other whole card completes.
- `session-c2c2b-complete`: actual compiled preparation/resume/retry and trusted JSON
  parity close a whole-card plan/Session; six wire objects validate. The view separates
  completed Session state from observed host status.

[Owner output](evidence-T048/lab-owner.json), [closure output](evidence-T048/lab-complete.json),
[owner wire](evidence-T048/owner-wire.json) and [closure wire](evidence-T048/complete-wire.json)
retain process observations. Agent reports, tests, reviews and owner decisions in
these fixtures are explicitly synthetic. They are not actual owner acceptance,
native multi-agent coordination, independent review of PG or final workflow acceptance.
Private operational state stays outside Git; [Lab checks](evidence-T048/lab-checks.json) identify committed fixture 831c060 and
280 preserved earlier tracked files. The authorized lab branch is pushed.

## Remaining selected work

PGF1-C2 retains explicit remaining-scope transfer, manual rebind/replan and model-gated
closure. PGF2 retains privileged leader Session operations plus trusted scheduling
and native parallel-Master evidence. PGF3 retains SDL handoff/integration/model release.
PGF4–9 retain owner usability, full recovery, sustained capacity/compatibility, native
integration, independent PG review, actual owner trial and explicit scope closeout.
KB038/052 remain in-progress; paused Session0008/KB050 and deferred MVP1 are unchanged.
No software/model release, merge or installation is claimed.
