# PGL3-C1 — Per-run controller isolation

2026-10-19, continuation of Session0007 T024 / PLAN-SDP-0018 / KB-SDP-052.
[Contract](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md),
[preflight](preflight-PGL3-C1.json) and [candidate](candidate-PGL3-C1.json) identify
this prerequisite for responsive coordination. It follows committed C2d activation.

## Delivered boundary

Control operations and terminal-session ownership use an exclusive run lock plus a
shared lock on the existing legacy checkout lease. Distinct registered runs can own
separate controllers concurrently. Same-run control and legacy checkout ownership
still exclude one another. Failed partial acquisition releases its shared lock;
process exit releases both locks. Workers, invalid paths and absent runs refuse
before creating lock directories. Store mutations remain exclusively serialized.
No state schema, host policy, worker privilege, source/scope guard or uncertain-send
retry contract changes. Read-only observers remain outside control ownership.

## Evidence and limits

Affected race tests and vet, complete SDPTool tests (including compiled controller
journeys), management/Toolkit and document/history integrity checks pass. New tests
exercise old/new exclusion in both acquisition orders for both lease kinds, same-run
refusal, two disjoint locks, cleanup after failed acquisition, no writes on invalid
requests and actual subprocess exit without explicit unlock. A controller-level
protocol test blocks one host's turn/start response while a second run launches,
then releases the first and confirms exactly one submission per host.

The first concurrency fixture reused a thread ID across both fake hosts and correctly
failed the existing ownership guard. The final fixture gives the second host its own
thread ID; the guard was not weakened. These checks are service/protocol evidence,
not two native models or a responsive leader conversation. C2d's native binary and
trial remain a separately pinned earlier candidate. This change does not replay its
refused report or answer any approval prompt.

No independent review, owner acceptance or publication is claimed. The owner choice
between automatic bounded reports and a decision on every report is still pending;
no answer is inferred from elapsed time. Leader/inbox/phase continuation and native
parallel evidence remain PGL3 work. Direct SAD interaction, evidence publication,
closure/replanning, integrated recovery/review and remote Steering remain in the plan.
