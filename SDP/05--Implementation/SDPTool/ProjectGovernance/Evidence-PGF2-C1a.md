# PGF2-C1a — Launch-scoped Session MCP operations

Session0007 T048 implements the common MCP entry point for the Session operations
already supported by PGF1. This is partial PGF2-C1: host scheduling/continuation,
new-goal intake and native multi-Master acceptance remain open. It deliberately does
not imply that all PGF1 transfers/replan or PGF3 model lifecycle work has completed.

The existing SessionWorkflow WORK :00002 defines ExternalLeaderMcp using the shared
WorkflowCommands boundary. [Scoped task](evidence-T048/mcp-blueprint-task.json)
generates diagnostic Blueprint af9d2bc773d1501c470a4c21592288ef2b6de6ff2cdfddb7ae475f80ac07d723
with passing constraints before product edits. This is unreleased TARGET, not model
acceptance. The existing CLI/core remain the workflow writer.

## Implemented boundary

The `sdp-mcp --session-controller INSTANCE --leader-skill PATH` launcher grants one
registered instance, pins the project-local mandatory skill and rechecks it on every
call. It cannot be mixed with coordinator/worker launch scopes. The twelve tools
inspect the definition/instance or call the same strict selection, assessment,
start/progress, verification/review, disposition, owner-request, continuation and
resume services as CLI. Core locks, exact revisions, canonical coverage and retry
receipts remain authoritative. Requests cannot switch instance or supply authority.

Owner decisions are refused. The existing coordinator and worker adapters gain no
privilege. No host dispatch, automatic next-phase execution or new Session adoption
is exposed. A skill snapshot proves availability only; native consumption is not
claimed. The local launcher is not an isolation boundary against arbitrary commands
by the same OS user.

## Actual verification

Compiled MCP process tests pass for start through whole-card/plan/Session closure,
core receipt parity, process reconnect/retry, tool inventory, strict argument refusal,
wrong-instance refusal, changed skill bytes, owner-decision refusal and mixed launch
configuration. Host report/test/review observations are synthetic service inputs.
Targeted race checks for the adapter, existing MCP contracts and new Session protocol
tests pass, as does vet. The durable pg_lab fixture `session-pgf2-c1a-mcp` reaches a
completed Session through actual MCP calls. This is protocol/integration evidence,
not execution by a native Project Leader or independent review of PG.

[Candidate](evidence-T048/mcp-candidate.json), [final race output](evidence-T048/mcp-final-race-output.json),
[vet](evidence-T048/mcp-vet-output.json), [final replay output](evidence-T048/lab-mcp-final-output.json)
and [lab identity](evidence-T048/mcp-lab-checks.json) retain the tested basis. Initial
fixture 153b263 is preserved; final named-skill replay 5301f9e leaves operational HEAD
and original receipts unchanged. [Replay](evidence-T048/mcp-final-replay.json) records
that result. Both commits are on the authorized lab branch.

Inspection found that a launcher could incorrectly label another skill as Project
Leader. The final adapter requires a regular SKILL.md with unambiguous frontmatter
name `sdp-project-leader`, checks this for API construction as well as CLI launch,
and pins those exact bytes. Wrong-name and changed-file cases refuse. The fixture's
named launcher skill is separate from synthetic assessment evidence and is not a
claim of native enabled inventory or actual consumption. The final compiled temporary
journey and final durable replay both pass. This bounded adapter increment is delivered. C2c2b broad race evidence applies to the
unchanged underlying lifecycle implementation; the adapter extension has its own
focused checks. No merge, release, installation or full-PG completion is asserted.

A supplemental [race check](evidence-T048/fresh-candidate-rework.json) changes actual
implementation bytes between next-step/rework executions and passes (122.111s).
Historical results remain readable while the new candidate receives fresh evidence.
This strengthens C2c2b verification without changing its executable lifecycle code.
