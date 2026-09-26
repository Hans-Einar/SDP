# SPR-VIEW-0001 — Reliable preview

| Field | Value |
| --- | --- |
| id | SPR-VIEW-0001 |
| project | VIEW |
| state | planned |
| Plans | PLAN-VIEW-0001, PLAN-VIEW-0002 |
| source | Owner selection; completed planning Study KB-VIEW-003 |

Goal: preserve preview selection and recover service reconnection through the two selected plans. Their authoritative phases and acceptance remain in the plans. Direct card membership is empty; KB-VIEW-003 is a completed source reference, not an execution member.

Both plans remain planned and receive SprintId SPR-VIEW-0001. No execution starts in this request. Sprint management payload 0.2 snapshots contain members: [] and plans: [PLAN-VIEW-0001, PLAN-VIEW-0002]. Record the Sprint creation and both plan membership updates with their actual predecessors and timestamps in the existing management ledger; preserve all older event bytes.

When execution is subsequently authorized and actually begins, record Sprint and selected plan start events. There are no linked execution cards to activate in this fixture. Keep the completed Study closed. At closure, explicitly record each plan's delivered, deferred or canceled disposition and any remaining work; Sprint completion cannot substitute for plan evidence.
