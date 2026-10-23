# PGL2-C2d — Approved implementation activation

2026-10-19, Session0007 T024 / PLAN-SDP-0018 / KB-SDP-052. The owner asks to continue
the plan toward completing the Session. This milestone connects a prepared approved
implementation package to the existing controller; it does not close the Session.
[Contract](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md),
[candidate](candidate-PGL2-C2d.json) and [preflight](preflight-PGL2-C2d.json) identify
the scope, baseline and tested artifacts.

## Behavior

The trusted local controller activates one exact approved Master preparation into
one ready implementation run. Atomic mapping and immutable receipts preserve retry
identity. A new operation cannot create a second run for the same preparation.
The existing controller separately launches it with the approved package and current
ImplementationPlan context. The same five assignment-bound MCP operations remain;
no new worker launch, approval or publication privilege is granted.

Live checks propagate model/task/decision drift through run inspection, MCP context,
advancement and launch. Interrupt/cancel/failure bookkeeping remains possible after
drift. A final optimistic check precedes the host send; this is not filesystem write
interception. Current same-user trust and external cross-worktree preflight remain.

First activation upgrades its store to 0.5. New readers retain older profiles and
writers preserve the new profile. The C2c reader refuses it without changing HEAD.
Activated preparations cannot be revised or retired: their exact historical package
and scope remain bound to the run. Explicit closure/replanning is still required to
release those claims. Checkout-wide control leases still prevent parallel controllers.

## Checks and native boundary

Affected governance/controller/MCP/CLI race tests, vet and the full SDPTool test suite
pass. Service tests cover source and decision drift, scope/binding guards, worker
refusal, immutable retry, duplicate activation, profile preservation, stop versus
resume behavior and injected failures before/after HEAD publication. Compiled trials
also verify stale-design launch refusal, duplicate launch refusal and old-reader
refusal. Schema validation checks 36 actual wire request/result objects. All are
same-context verification; no independent review is claimed.

[Protocol trial](trial-PGL2-C2d-protocol.json) uses a recorded host subprocess and
real provider services. [Native trial](trial-PGL2-C2d-native.json) uses the complete
official Codex 0.161.0 package and existing managed account, without changing the
installed 0.162.1 binary, model, credentials or global configuration. Release archive
and binary hashes are retained in the candidate. The complete package supplies the
execution helper absent from the earlier standalone download.

The first native attempt ended without code changes: its public final message said
the execution host was missing. [Failure evidence](native-failure-PGL2-C2d.json)
retains that failure. A fresh isolated run with the complete package performed the
actual Go task. Independent external golden cases failed before execution and passed
afterward; project Go tests passed. Only code/normalize.go and code/normalize_test.go
changed. All 176 pre-existing owner project/runtime entries remained unchanged.
Provider artifacts/bundles are real, but this is a synthetic task, not proof of
production code conformance to an SDL model.

The native host's attempt to submit run_transition was refused under the selected
never approval policy. Its public final message reports automatic approval rejection;
no successful governance report is observed. The run remains running while its host
attempt is completed. The trial's passed status covers activation, actual bounded
code execution, behavior and preservation, **not successful reporting or acceptance**.
No policy was weakened, no rejected call was replayed and no report was invented.
Existing explicit per-request interaction remains separate and requires an actual
answer; automatic report delivery is unfinished integration work.

## Owner inspection and remaining work

The separate preview is /home/warloc/git/pg_lab/previews/pgl2-c2d. The installed
read-only ./implementation helper shows the task, current authority coverage,
recorded run/host states, missing report and changed-file links. Reads preserve state.
Earlier previews and original project/runtime are unchanged. The owner walkthrough
now explains what this proves and what it does not.

PGL2 activation is delivered at this bounded level. PGL3 responsive leader, two
Masters, durable reports and continuation are next; SAD launch/direct dialogue,
closure/publication, recovery, independent review and remote Steering retain their
plan obligations. No actual owner acceptance, main merge, release or completed
Session is claimed.
