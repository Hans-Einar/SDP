# PGI3-C2 — live execution and correlated recovery

Started 2026-10-07; closeout 2026-10-08 (Europe/Oslo). PLAN-SDP-0018,
KB-SDP-038, Session0007 T008. **Live workflow checkpoint; PGI3 is not complete.**

## Candidate and actual environment

Base `a2580e9` is the PGI3-C1 controller checkpoint, following PGI2 `12542b2` and
PGI1 `d3e457b`. These commits now exist in the original repository's object database
on `sdp/project-governance-pilot`. Development uses the persistent worktree
`/home/warloc/git/SDP-project-governance`. The shared SDP-vNow worktree and other
agents' uncommitted changes are preserved. No merge, push or release is claimed.

[Live evidence](live-PGI3-C2.json) records changed product file hashes, binary and
schema identities, actual run/attempt/input mappings, sanitized MCP results,
fixture source pre/post hashes and exact test-driver scripts. It excludes raw
account data, credentials, reasoning and host stderr. Test-driver steering and
interruption are scenarios under the owner's PGI3 authorization, not additional
verbatim owner conversation. The scripts are frozen evidence with temporary paths;
use a fresh disposable root for reproduction, not an existing project.

Codex is now 0.161.0. Its exported default schemas were inspected before advancing
the new-connection pin. Core records still recognize historical 0.160.1, but a new
connection cannot silently relabel an old attempt's host version. This is local
application compatibility, not a published method/Toolkit version. The actual
probe observes managed ChatGPT mode, MCP 2025-11-25 and the expected three read tools
with MCP Go SDK v1.8.0. No API-key fallback or account reconfiguration occurred.

## Real workflow observations

| Scenario | Actual result | Evidence boundary |
| --- | --- | --- |
| Useful Maintenance task | Codex replaced ASCII-only boundary trimming with Unicode-aware trimming and added eight tests; received in-scope steering | Real managed-account turn and tools, not a recorded fake host |
| Preservation | Only normalize.go and normalize_test.go changed; all declared excluded files, authored Session and ledger retain their hashes | Whole disposable fixture inventory checked |
| Verification | Test driver reran `go test ./...`, exit 0 | Same-context verification, not independent review |
| MCP reads/report | run_read succeeded; run_transition failed with an approval-required error under never | Exact sanitized host result retained; the refusal was not bypassed |
| Observer closure | Closing the actual controller stdout reader left execution running | Separate status command observed the mapped running attempt |
| Interrupt/recovery | Actual interrupt produced interrupted status; a new process recovered the same run/Session with an epoch gap | Host history, not a fabricated final result |
| Resume | Same routine/Session/assignment continued in a new actual thread and finished verification | Two attempts: first interrupted, second completed; routine remained unaccepted |
| Lost acknowledgment | Transport proxy dropped the actual turn/start response, then closed after real execution completed | Input became uncertain; history matched exact clientId/payload; one existing attempt acknowledged without resubmission |

Host execution completion does not complete the SDP routine. The report refusal
leaves an explicit gap rather than granting owner acceptance or writing project
records. No native child or independent Reviewer was launched; instruction text
naming the fresh Master is not proof of an independently verified native role tree.

## Implemented correction

Recovery now correlates uncertain launch/steer input with a **unique** recorded
userMessage.clientId and exact payload inside the already bound thread. It persists
that observed message before controller-only input.resolve can acknowledge the
input. The core checks client identity, text, turn, attempt and recorded coverage;
no role string or similar-looking text can supply the missing proof.

Missing/duplicate/conflicting correlation refuses. Completely unmapped thread
creation, incomplete active history and interrupt-outcome resolution remain
uncertain rather than being resent. Older messages without clientId are not
retroactively promoted into proof. This closes the measured mapped-turn response
loss case while preserving the no-duplicate-dispatch invariant.

## Verification

- Targeted core/controller/host/MCP race tests pass.
- New regression tests reject text-only, wrong client, changed payload, duplicate
  correlation and a different mapped turn; recorded exact correlation resolves
  without advancing the routine. Existing receipt, authority and replay tests pass.
- Four captured work requests conform to the installed 0.161.0 default schemas.
  All 28 captured application cases (21 CLI, four host/core, three dispatch) pass.
- Full SDPTool `go test -race -count=1 -p 2 ./...` and `go vet ./...` pass.
  Application and host schema checks pass. Management and Toolkit validators pass;
  closeout outputs and their hashes are retained in the accompanying live evidence.

The initial closeout validation caught a verification relation inserted under
a nested slice field. It was moved to the root verification map and validation
rerun; no historical relations or ledger bytes were changed.

The first live task/interruption used the version-updated C1 binary. Lost-response
and replacement-attempt tests used the final C2 binary with correlation recovery;
both binary identities are retained. No historical test is relabeled as a test of
a later implementation. The final source inventory identifies this checkpoint.

## Remaining decision and work

The actual report error is `MCP tool call requires approval, but approval policy is
never`. Existing [host documentation](https://learn.chatgpt.com/docs/app-server)
describes request/response owner interaction; [configuration reference](https://learn.chatgpt.com/docs/config-file/config-reference)
distinguishes approval modes. This checkpoint preserves the selected policy.
The owner was asked whether the client should present the specific request for
explicit owner approval or keep refusal visible. No answer is inferred from elapsed
time, and no policy change or automatic grant is implemented here.

Recommended next design: an explicit opt-in owner interaction path, bound to the
current host request/thread/turn, with exact accept/decline recording and refusal
on missing response or identity mismatch. Never convert this into automatic
approval of unrelated filesystem/network requests or completed-work acceptance.
Executing-turn reattachment and applicable native role/provenance evidence also
remain open. PGI4 owns checked evidence/publication; PGI5 owns independent integrated
verification. PGI3 remains in progress until those PGI3 obligations are resolved.

Before shared-history integration, resolve private C1 event-ID collisions with
later shared blueprint/SDUI events. Both existing histories remain unchanged in
this checkpoint; no overwrite or implicit merge is performed.
