# ProjectGovernance ecosystem

This grouping contains systems supporting the SDP process. It is catalog metadata,
not a deployed software system or an SDL `system` declaration. The shared process,
KanBan, plans and evidence remain at the root SDP area.

| System | Current delivery | Proposed next responsibility | Command/library boundary |
| --- | --- | --- | --- |
| [SDPTool](SDPTool/README.md) | Go project/navigation/installation facade | Shared routine decision and durable work-state core | Existing `sdptool` and Go libraries; routine APIs not implemented |
| [SDPMCPAdapter](SDPMCPAdapter/README.md) | Study only | Agent access through MCP to the shared core | Proposed standalone `sdp-mcp` plus adapter library |
| [CodexClient](CodexClient/README.md) | Study and exported Codex schemas only | One app-server controller and observable bounded assignments | Proposed `sdp-codex`; no client runtime |
| [KanBanCLI](KanBanCLI/README.md) | Existing Bash listing command | Preserve lightweight status/grouping access | `kanban`, from Toolkit/scripts/cli/kanban.sh; no Go library |
| [KanBanTUI](KanBanTUI/README.md) | Proposal only | Read-only interactive board/card navigator | Proposed `sdp-kanban`; reuse SDPTool board read APIs |

Executable names marked proposed are suggestions, not installed commands or release
commitments. Reusable libraries do not become containers merely because several
frontends use them. In particular routine policy belongs to SDPTool, not another
independently deployed engine. MCP is an agent-facing adapter; app-server is the
external Codex-facing client protocol. Neither is project acceptance authority.

## Dependencies and authority

Each System.design is independently valid design-core 0.5. No file imports another
and no generated projection may imply cross-file link resolution. External units
inside boundary scenarios are explicitly unowned placeholders. Proposed envelope
contracts describe local views of a future shared boundary; they are not multiple
published APIs. A later DesignPlan must select one versioned application contract
and its adapters. SDL source-set support is a separate language deliverable.

- SDPTool consumes SDL/SDUI libraries and provides project/process operations.
- SDPMCPAdapter calls SDPTool's shared application service and returns its decisions.
- CodexClient controls external Codex app-server and reads SDPTool durable state.
  Native Master delegation is desired; the client must not replace Master decisions.
- KanBanTUI consumes SDPTool's board read service; KanBanCLI reads card tables directly
  and does not claim full ledger validation.
- XFMD and gh-sdp are external consumers with their own ownership. Codex and gh-tree
  are external products, not systems to implement inside this source catalog.

The earlier [detailed SDPTool model](../../03--Architecture/SDPTool.design) remains
canonical for existing navigation and installation design. New entries describe
bounded tool responsibilities and proposed additions; they do not migrate or
supersede that detailed model. System READMEs link implementation evidence and
[RGS2](../../02--Requirements/RoutineGovernance/StudyPlan.md) research precisely.
This architecture authoring does not activate any of those product proposals.
