# PGI2-M1 — adapters and measured protocol evidence

Date: 2026-10-07. Plan: PLAN-SDP-0018. Design: SDPTOOL-DES-0001.
Result: the selected adapter increment is delivered and passes bounded verification. This increment does not finish KB038 or Session0007.

## Delivered boundary

The local MCP server delegates application decisions to the existing Go core.
Launcher-supplied run/attempt/epoch scope permits report, fail and scope-change;
model arguments cannot grant owner/controller authority or select another run.
The separate controller command probes the installed app-server and optionally
its thread-local MCP server. It sends no model turns and changes no account mode.

Durable attempt mappings, connection epochs, completed visible observations,
strict replay receipts, stale-cursor refusal and conservative reconciliation are
implemented. Disconnect leaves the attempt unknown until explicit observed-history
reconciliation. Replayed items never fabricate a final response or complete SDP
work. Routine snapshots/receipts summarize attempts; message retrieval uses
count- and byte-bounded pages. Oversize content refuses rather than truncates.

Build targets are `SDPTool/cmd/sdp-mcp` and `SDPTool/cmd/sdp-controller`. They are
not yet distributed installer entries. Their package READMEs describe usage and
limits. The unpublished 0.1 schema reads older PGI1 stores; downgrading a store
with PGI2 fields to the strict older reader is unsupported.

## Exact candidate and environment

[Candidate manifest](candidate-PGI2.json) identifies base commit d3e457b and the
SHA-256 inventory of changed SDPTool files, with commands, environment and actual
outputs. Tests run on that uncommitted candidate; the subsequent milestone commit
also includes lifecycle and evidence documents. The [protocol pin](protocol-PGI2.json)
records installed-schema identities and the sanitized actual probe result.

Implementation and tests use the isolated clone at `/tmp/sdp-pgi1-8SeezS/repo`,
branch `sdp/project-governance-pilot`, to preserve concurrent shared work. The prior
PGI1 host's read-only Git limitation is historical, not the present reason for
isolation. Go 1.27.1 linux/amd64; module directive remains Go 1.26.0. The official
MCP Go SDK v1.8.0 and its transitive modules are pinned in go.mod/go.sum. Dependencies
were downloaded during preparation; retained verification uses cached modules.

## Executed verification

| Check | Evidence boundary |
| --- | --- |
| Full SDPTool `go test -race -count=1 -p 2 ./...` | Core, adapter, compiled CLI and stdio tests, and existing regression packages |
| `go vet ./...` | Static checks of the same module |
| Application schema checks | 21 observed compiled CLI cases and three observed host/core cases |
| Rebuilt controller + real Codex 0.160.1 probe | Initialization, ChatGPT account-mode read, idle ephemeral read-only thread, MCP tool discovery |
| Management and Toolkit validation | Passed in both checkouts: isolated 61 cards / 38 management records / 506 events; shared 63 / 39 / 518; four lineage operations each |
| Candidate/preimage hashes and whitespace | 22 candidate files match in both checkouts; original ledger prefixes preserved; new links and whitespace pass |
| Returned-workspace adapter race tests | All three core/adapter packages pass, including compiled MCP tests; no full-regression claim for concurrent changes |

The live probe discovers context_read, route_evaluate and run_read using MCP
2025-11-25. It reports managedAccount=true and modelSubmitted=false. The stdio
server's authStatus=unsupported means it advertises no OAuth, not that tool
negotiation failed. Account credentials, email, raw auth responses and host stderr
are excluded. This proves account mode, not model entitlement or model execution.

Earlier runs exposed a copied fixture's incorrect Maintenance ledger path and
an unkeyed test literal; both were corrected. A superseded broad run overlapped
source edits and produced a compiler import failure, so it is not delivery
evidence. Self-checks also added encoded-byte paging and blocked late observations
from reviving a disconnected attempt without reconciliation. Retained results
come from the final frozen source candidate; no assertions were removed.
The record validator also caught an initially misplaced verification YAML block;
it was moved to the existing verification map and both validators then passed.

## Acceptance mapping

| Case | Established now | Still required |
| --- | --- | --- |
| PG-A10 | Core/MCP parity; launch-bound rights; forged role, cross-run/controller action and revoked epoch refuse | Live native child and independent-review attribution |
| PG-A11 | Persisted mapping/restart, unknown disconnect, history projection, epoch-gap and replay fixtures; transport-loss detection | Actual running-turn reconnect and PGI3 controller loop |
| PG-A12 | Recorded user/assistant phases and IDs, bounded attachment references, interrupted/partial history, duplicate/conflicting replay | Owner steering and complete automatic Session capture during live execution |
| PG-A15 | Real installed host + chosen SDK protocol initialization/discovery, ChatGPT account-mode read; unsupported version refuses | Model entitlement, fresh Master and native role mapping in live task |
| PG-A18 | Strict MCP arguments, duplicate host JSON rejection, frame/text/attachment/state limits, revision-pinned byte/count pages and cancellation | Wider workflow limits under integrated execution |

These are recorded fixture tests and a real idle protocol probe, not a live agent
execution/recovery trial. No independent review is claimed. The controller still
needs durable outbound intent, observer/control loop and actual launch/steering
(PGI3), candidate-bound evidence/publication (PGI4), and integrated independent
verification (PGI5). Missing full-loop evidence cannot be inferred from the probe.

## Disposition

PGI2-M1 is delivered at the stated adapter boundary. PLAN-SDP-0018 and Session0007
remain active. KB038 returns to ready for the next PGI3 owner-workflow assignment;
KB036/037/042/043 retain their broader backlog scope. ModelGovernance and Session0008
semantic blueprints remain separate. No merge, push, release or model execution.
