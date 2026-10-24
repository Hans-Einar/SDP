# PGL3-C6 — Shared coordination view and responsive local input

Session0007 T032, PLAN-SDP-0018 and KB-SDP-052.
[Design](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md#pgl3-c6--shared-coordination-view-and-responsive-local-input-2026-10-24),
[preflight](preflight-PGL3-C6.json).

## Candidate and changed operation

Baseline 59830cf79b621016f869b7dd9c09ad6656e0a8cc plus the retained candidate
source hashes. Work remains on sdp/codex-forward-compatibility under the selected
current-branch/milestone-commit policy. No competing local changes observed.

The read-only coordination projection combines registered runs, questions/replies
and reports/acknowledgements. It pins one operational generation, prioritizes
blocking questions and pending reports, bounds pages and rejects stale cursors.
Inspection does not create state or infer live host/acceptance status. Manual
Sessions/cards without a registered run are outside this view.

The explicitly launched coordinator MCP adapter has exactly two tools: read this
view and acknowledge a report under its launcher identity. Default inspection and
bound workers keep their existing tools. No owner reply, approval, launch, resume,
completion or publication tool is exposed to the coordinator model context.

The terminal watch accepts bounded JSON commands while polling changed views.
Pinned report acknowledgements and owner-attributed replies use existing atomic
services. Queued owner input takes priority over polling; invalid input produces
an error without ending the console. Closing the observer leaves workers untouched.
Pending information survives in the existing store; no new queue/profile is added.

## Evidence and limits

Core tests cover empty read without writes, stable refresh, priority and full
pagination without duplicates, bounded encoded pages, per-run filtering and stale
source/generation cursors. MCP tests use a compiled stdio server, checking explicit
launch authority, exact tools/annotations, spoofing and mixed-scope refusal,
acknowledgement/retry across server restart and unchanged run state.

The compiled console test uses two synthetic work runs: one asks a question while
the other progresses. It submits malformed input and then an attributed answer,
receives a subsequently reported result, acknowledges it and recovers both records
in another process. The unrelated project fixture file is preserved. Unit checks
cover command scope, terminal control characters, coalesced polling and output
failure. These are service/transport workflows, not native agent concurrency.

No native model/subagent task, AI Project Leader scheduler/wake-up, independent
review or whole-PG acceptance is claimed. Current-host task/recovery, automatic
leader decision turns, SAD/scope closure, capacity and canonical evidence/publication
remain on the parent plan. SDL-only, paused Blueprint and Ponsse MVP1 boundaries
are unchanged. Stable installation and global MCP configuration are not modified.


## Verification refinements

The first console evidence helper decoded/re-encoded events and introduced null
optional fields that the actual emitter omits. The helper now retains exact raw
JSON emitted by the process; the earlier capture is retained as rejected evidence,
not a runtime failure. Final full-module checks use the corrected capture helper.

Output-only schemas were separated from the shared MCP input definitions so the
new view contracts do not enlarge every worker tool description. They reference
the existing canonical work/report/attention/result/error schemas. Actual wire
validation caught an incomplete error shape (missing retryable); that output-only
schema now references the canonical error definition. It is not embedded in Go;
final wire validation covers this last schema-only correction.


## Final result

- Full SDPTool `go test -race -count=1 ./...`: passed, 15 tested packages, 205.41 seconds.
- Full SDPTool `go vet ./...`: passed.
- Compiled console and coordinator stdio MCP workflows: passed within that suite.
- Ten actual terminal/MCP objects validate against the final canonical schemas.
- Private controller/MCP development binaries built successfully; identities retained.

See [retained artifacts](evidence-PGL3-C6/README.md),
[candidate](evidence-PGL3-C6/candidate.json),
[checks](evidence-PGL3-C6/checks.json) and
[wire validation](evidence-PGL3-C6/schema-checks.json).
Same-context inspection/verification only. This delivers the bounded local
coordination entry point; the parent PG plan and Session remain active.


[Final record checks](final-record-checks-PGL3-C6.json) confirm 11 candidate source/
test/schema hashes, 12 retained artifact hashes, preserved ledger prefixes and
218 relative links. Management replay (73 cards, 47 records, 656 events), Toolkit
validation and Git whitespace checks pass.
