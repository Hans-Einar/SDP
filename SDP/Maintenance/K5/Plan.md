# K5 — CardState and a small shell CLI

Owner request, 2026-09-24. Branch: sdp/phase-k5-card-state, stacked on L1 32f6ae5.

| Milestone | Delivery | Verification |
| --- | --- | --- |
| K5-M1 | One visible CardState per card, queue/review conventions, no separate gate/queue store | Metadata, board replay, links and historical prefix checks |
| K5-M2 | kanban status/state and installer; install in ~/bin | Scope/errors/installation tests, shell syntax, real-board listings |

Use Bash and ordinary Unix tools. No mutation command, database, daemon, custom
Git engine or dynamic graph. Existing Git/worklogs and append-only ledger retain
content/lifecycle history. This is a local additive convention, not a released
cross-repository compatibility profile. Larger distribution/version work remains
separate. Preserve concurrent uncommitted KB-SDP-014 registration.

After K5, R3 completes the remaining bounded structure/consolidation work for
KB-SDP-001/010 and records concrete owner review gates and proposed next work.

K5-M1 and K5-M2 are complete. Owner authorized including the concurrent registration
after its session finished; evidence records its append-only integration.

K5-M3 corrects the newly introduced event collision and records final replay checks.
