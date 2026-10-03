# PGD1 evidence and limits

Date: 2026-10-03. Scope: design, source inspection and offline protocol schema
generation. This is not a runtime trial or independent product review.

## Inputs inspected

- Root AGENTS, maintained Skills/sdp, planning, architect, change-analysis and
  traceability routines; project management, Session and KanBan contracts.
- Completed RGS2 studies, synthesis and routine catalog; KB036/037/038/042/043;
  ProjectGovernance ecosystem and separate ModelGovernance ownership contract.
- SDPTool project discovery/strict JSON/path handling, CLI dispatch, Session tests,
  module and proposed CodexClient/SDPMCPAdapter boundaries. Current inspected source
  supplies discovery and browsing, not the new generic routine/client operations.
- [Official app-server documentation](https://learn.chatgpt.com/docs/app-server),
  retrieved during Session0007 T001 on 2026-10-03; reused in T002. Generated local
  schemas govern the inspected version-specific field inventory below.

## Installed protocol inspection

Three read/inspection commands ran successfully in a temporary directory with an
isolated empty Codex configuration and an allowlisted environment, without inherited
API-key credentials. No account inspection, authentication change, server launch,
conversation, model turn or MCP connection was performed.

| Command | Observed result |
| --- | --- |
| codex --version | codex-cli 0.160.0; exit 0 |
| codex app-server --help | stdio/unix/WebSocket transports advertised; exit 0 |
| codex app-server generate-json-schema --out <temporary directory> | 314 default-surface JSON files; exit 0 |

All commands reported a warning that PATH helper aliases would not be created
under the temporary Codex configuration directory. The commands still exited 0;
no claim about helper availability follows from this probe. Unlike the earlier
RGS2 study, this turn did not export the experimental surface.

[Protocol evidence](protocol-evidence.json) retains command results, tool hashes,
schema inventory digest and selected required/property/method inventories. Full
generated schemas are disposable probe outputs and are not copied into product
source. Regenerate from the pinned installed version when implementing bindings.
The exported structures establish shape only; server behavior remains untested.

Relevant observations: thread/start and thread/resume exist; turn/start and
turn/steer expose clientUserMessageId; turn/steer requires expectedTurnId; completion
and message/tool lifecycle notifications exist. Correlation field presence does
not prove deduplication or notification replay. Only the default schema surface
was inspected; do not use optional experimental documentation fields by assumption.

## Design reconciliation

The contract keeps one process core and separate app-server/MCP adapters. It
distinguishes local operational state from canonical project lifecycle/evidence,
and defines checked recoverable publication instead of claiming atomic multi-file
updates. It retains the desired native delegation mapping but requires actual
caller/context evidence before accepting an independent review claim.

The first implementation plan sequences a pure core and fixtures before live
integration, with 18 mapped acceptance scenarios. ModelGovernance and semantic
blueprint implementation are not prerequisites to that first bounded task.

## Verification record

Phase-closeout checks passed:

- `python3 SDP/ProjectManagement/validate.py`: 61 cards, 37 management records,
  four lineage operations and 484 events; schema/history/metadata/placement/grouping.
- `python3 Toolkit/scripts/validate_sdp.py --repo .`: Toolkit validation passed.
- Changed-document local links resolve; both ledgers preserve the current committed
  byte prefix, including the concurrent ModelGovernance delivery.
- `git diff --check`: passed; design/acceptance/protocol artifact hashes match
  EVT-SDPTOOL-DES-PGD1 and the SDPTOOL-DES-0001 relation resolves.

Concurrent ModelGovernance commits arrived during preparation. Unpublished local
ProjectGovernance KanBan event IDs were reallocated after the committed counter;
the committed history was preserved exactly. No ModelGovernance files are part of
the scoped PGD1 commit. The introducing Git commit identifies this design delivery.

No Go runtime test is claimed: no product code was changed, and `go` was not found
on PATH during preparation. Exact protocol behavior, authentication, MCP/SDK
intersection and integrated execution remain PGI acceptance work.
