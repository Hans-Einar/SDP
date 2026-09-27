# External KanBan references — ER1

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0004 |
| project | SDP |
| state | completed |
| PlanType | ImplementationPlan |
| Systems | SDPTOOL |
| BranchPolicy | current |
| CommitPolicy | milestone |

## Outcome and authority

Owner decision, 2026-09-27: foreign card references are permitted, explicitly
external and unverified. Do not locate them across worktrees, require a checkout,
contact a remote service or introduce externalConfig.yaml. Source:
[KB-SDP-034](../../../KanBan/completed/%23034--Bug--External-KanBan-references.md).

REQ-SDPTOOL-003 and the [consumer contract](../../../../SDPTool/Contract.md)
govern this bounded navigation fix. The board descriptor defines local namespaces;
a missing primary in those namespaces remains invalid. A syntactically valid
foreign primary is informational, even when its card may not exist. Local card
open targets retain their identity and revision. No XFMD application changes.

## Phase ER1 — informational external references

ER1-M1: add an optional externalReference card-ID field to navigation nodes.
Keep reference for resolvable local node IDs only. Validate primary ID syntax,
classify using projectId plus namespaces, preserve local/history failures and
update requirements and consumer guidance. Test local, foreign, malformed and
missing references, shared namespaces and invalid history. Build the CLI and run
read-only navigation against XFMD's active worktree; verify actual node states.

## Git and verification policy

Use sdp/kanban-external-refs, branched from 26529be. One milestone commit includes
code, contracts, tests and evidence. Existing authorization permits phase push;
this fix does not select another merge/release or replace the installed signed
binary. Preserve unrelated SDL/go/sourceinput work. Run SDPTool tests and vet,
management replay and documentation-link verification. Close only on actual
acceptance results; record release/distribution separately.

## Evidence and remaining work

ER1-M1 delivered. Go 1.27.1 on Linux: `go test -race ./...` and `go vet ./...`
pass from SDPTool. Eleven reference-scope regression cases cover foreign/unknown
projects, local and shared namespaces, missing local cards, malformed IDs/paths/URIs,
invalid namespace metadata and broken history. Existing movement, CardState and
consumer tests remain passing. See [test output](evidence/go-tests.txt).

The built CLI was run read-only against `/home/warloc/git/xfmd-sdl-navigation`:
`sdptool /home/warloc/git/xfmd-sdl-navigation tree`. Its KanBan tab is validated,
all 18 cards are present, and KB-XFMD-012 is available with
`externalReference: KB-SDP-014`, no local reference and its original local open
target. [Actual output](evidence/xfmd-tree.json) records these facts. XFMD Git
status remains clean. No GUI badge claim or external-card existence claim.

[Candidate hashes](evidence/candidate.json) identify base 26529be plus the tested
ER1 source diff and built CLI. Management replay, preserved-history/document-link
checks and Toolkit validation pass. Requirements clarification and consumer
contract are updated; architecture/SDL declarations are unchanged because no
responsibility, capability or model graph changed. No independent review claimed.

KB-SDP-034 is complete for implementation. Related KB-SDP-014 (reusable board
contract) and KB-SDP-032 (viewpoint feedback) remain separate backlog scope.
Subsequently merged and published as SDP 0.2.1, selected by gh-sdp 0.1.1.
[RP2](../../../Maintenance/RP2/Plan.md) records exact signed publication and the
installed client's successful read-only XFMD navigation. ER1's original test
candidate and implementation evidence remain unchanged.
