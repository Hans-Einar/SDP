# PGL3-C5 — Durable report acknowledgement

Owner continuation: Session0007 T031; PLAN-SDP-0018 / KB-SDP-052.
[Design](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md#pgl3-c5--durable-report-acknowledgement-2026-10-24),
[preflight](preflight-PGL3-C5.json).

## Candidate and scope

Baseline a84d71cdc36ca93fcf7fbedc4cbd2b7b8f125728 on
sdp/codex-forward-compatibility plus the source hashes in the evidence manifest.
No other agent or worktree writes observed. The implementation is one local
controller operation over the existing atomic store and a versioned inbox view.
It does not launch a model, add an MCP privilege or change the host wire contract.

The trusted controller pins the exact report and inspected generation/revision,
records coordinator attribution and a note, and gets a durable retry receipt.
Acknowledgement leaves all runs unchanged. Profile 0.8 is selected at that explicit
write; later attention, continuation, preparation and design writers preserve it.
Inbox profile 1 is retained before acknowledgement, profile 2 shows recorded
processing afterward. Existing report text/history remain immutable.

## Verification boundary

The service suite exercises restart, exact retry after later state/source changes,
conflicting retry/duplicate report delivery, concurrent acknowledgement, unchanged
runs, wrong authority/checkout/report/revision/generation/source/state refusals,
corrupt history, and coexistence with attention/continuation and design preparation/
activation. The compiled CLI fixture inspects, acknowledges, retries and reopens
through distinct processes and preserves every project fixture file. It uses a
scripted report, not a native agent or owner acceptance.

Actual old-reader verification uses a controller built from baseline a84d71c;
actual command output is validated against the application and both inbox profiles.
Full affected-module race/vet checks pass. Canonical record validation is recorded
in the final record check below.

## Limits and next work

Acknowledgement proves recorded processing, not semantic understanding, evidence
verification, owner approval, canonical publication or permission for another turn.
No responsive Project Leader conversation/event loop, native task on Codex 0.164.0,
parallel-agent trial, capacity extension or independent review is claimed. Source
coverage remains a separate filesystem check rather than an atomic source snapshot.
The existing receipt capacity remains bounded; this is not a sustained-hours test.
Full PG readiness, PGL2 SAD/scope closure, PGL3 responsive coordination and PGI4/5
remain open. Blueprint implementation pause, SDL-only direction, KB059 onHold and
Ponsse MVP1 deferral are preserved.


## Result

- Full SDPTool `go test -race -count=1 ./...`: passed, 15 tested packages; 200.82 seconds.
- Full SDPTool `go vet ./...`: passed.
- Compiled CLI flow and actual baseline-reader refusal: passed in that suite.
- Four actual wire objects (request, result, inbox 1 and inbox 2): passed schema validation.
- Evidence: [retained artifacts](evidence-PGL3-C5/README.md),
  [source candidate](evidence-PGL3-C5/candidate.json),
  [check results](evidence-PGL3-C5/checks.json),
  [wire objects](evidence-PGL3-C5/wire.json).

Same-context implementation inspection and verification only. PGL3-C5 is delivered
at this local service/CLI boundary; the parent plan and Session remain active.


[Final record checks](final-record-checks-PGL3-C5.json) confirm 12 tested source
hashes, seven retained artifact hashes, unchanged ledger prefixes and 205 relative
links. Management replay (73 cards, 47 records, 654 events), Toolkit validation and
Git whitespace checks pass.
