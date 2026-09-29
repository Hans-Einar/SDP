# SDP ecosystem catalog

Start here to find the tools, their responsibilities and their SDL models.
The owner selected ecosystem grouping on 2026-09-29. [Architecture decisions](../03--Architecture/Ecosystems/Decisions.md)
and [PLAN-SDP-0008](../03--Architecture/Ecosystems/Plan.md) explain scope and evidence.

## Ecosystems and systems

Each linked README maps existing commands/libraries, target boundaries, dependencies,
and remaining decisions. `System.design` is a complete design-core 0.5 model,
not a claim that the corresponding new executable or runtime exists.

| Ecosystem | System and source guide | SDL model | SDPTool model ID |
| --- | --- | --- | --- |
| ProjectGovernance | [CodexClient](ProjectGovernance/CodexClient/README.md) | [System.design](ProjectGovernance/CodexClient/System.design) | `pg-codex-client` |
| ProjectGovernance | [KanBanCLI](ProjectGovernance/KanBanCLI/README.md) | [System.design](ProjectGovernance/KanBanCLI/System.design) | `pg-kanban-cli` |
| ProjectGovernance | [KanBanTUI](ProjectGovernance/KanBanTUI/README.md) | [System.design](ProjectGovernance/KanBanTUI/System.design) | `pg-kanban-tui` |
| ProjectGovernance | [SDPMCPAdapter](ProjectGovernance/SDPMCPAdapter/README.md) | [System.design](ProjectGovernance/SDPMCPAdapter/System.design) | `pg-mcp-adapter` |
| ProjectGovernance | [SDPTool](ProjectGovernance/SDPTool/README.md) | [System.design](ProjectGovernance/SDPTool/System.design) | `pg-sdptool` |
| SDL | [ActionRuntime](SDL/ActionRuntime/README.md) | [System.design](SDL/ActionRuntime/System.design) | `sdl-action-runtime` |
| SDL | [CodeGeneration](SDL/CodeGeneration/README.md) | [System.design](SDL/CodeGeneration/System.design) | `sdl-code-generation` |
| SDL | [DevelopmentHost](SDL/DevelopmentHost/README.md) | [System.design](SDL/DevelopmentHost/System.design) | `sdl-development-host` |
| SDL | [DocumentService](SDL/DocumentService/README.md) | [System.design](SDL/DocumentService/System.design) | `sdl-document-service` |
| SDL | [DocumentSnapshot](SDL/DocumentSnapshot/README.md) | [System.design](SDL/DocumentSnapshot/System.design) | `sdl-document-snapshot` |
| SDL | [Frontend](SDL/Frontend/README.md) | [System.design](SDL/Frontend/System.design) | `sdl-frontend` |
| SDL | [Viewpoints](SDL/Viewpoints/README.md) | [System.design](SDL/Viewpoints/System.design) | `sdl-viewpoints` |
| SDUI | [SDUICodegen](SDUI/SDUICodegen/README.md) | [System.design](SDUI/SDUICodegen/System.design) | `sdui-codegen` |
| SDUI | [SDUIFrontend](SDUI/SDUIFrontend/README.md) | [System.design](SDUI/SDUIFrontend/System.design) | `sdui-frontend` |
| SDUI | [SDUINativeHost](SDUI/SDUINativeHost/README.md) | [System.design](SDUI/SDUINativeHost/System.design) | `sdui-native-host` |
| SDUI | [SDUIPresentation](SDUI/SDUIPresentation/README.md) | [System.design](SDUI/SDUIPresentation/System.design) | `sdui-presentation` |
| SDUI | [SDUIRuntime](SDUI/SDUIRuntime/README.md) | [System.design](SDUI/SDUIRuntime/System.design) | `sdui-runtime` |

Read the ecosystem inventories for complete command/package coverage:
[ProjectGovernance](ProjectGovernance/README.md), [SDL](SDL/README.md),
[SDUI](SDUI/README.md). Development fixtures, packaging helpers and external
collaborators are classified explicitly rather than manufactured into products.

## Support tooling and non-product assets

| Repository entry | Ownership / disposition |
| --- | --- |
| [Profile builder](../../SDPTool/tools/profile/main.go), [package recipe](../../SDPTool/package.sh) and [bootstrap library](../../SDPTool/bootstrap/) | SDPTool release/build support; separate maintainer command, not a second installer engine or new user-facing system |
| [CLI helper installer](../../Toolkit/scripts/cli/install-cli.sh) | Installs shell helpers including KanBanCLI; not SDP project installation |
| [Schema/document validation](../../Toolkit/scripts/validate_sdp.py) and [management validation](../ProjectManagement/validate.py) | Existing verification tools; retained Python validators, no alternate production SDL parser |
| [SDL fixture commands](SDL/README.md#complete-tracked-command-inventory) | Demo, compiled demo and smoke harnesses, not Ponsse production systems |
| [SDUI native smoke](../../SDUI/go/host/fynehost/cmd/smoke/main.go) and [reload smoke](../../SDUI/go/host/fynehost/cmd/reloadsmoke/main.go) | SDUINativeHost verification harnesses, not separate products |
| [Skills](../../Skills/README.md), Template and Toolkit payload/schemas | Distributed instructions/data; not binaries or independently running governance services |
| `.github/workflows/validate.yml` | Repository CI coordination; model/runtime acceptance remains bounded by actual executed checks |

These classifications cover all tracked Go main entrypoints in SDL, SDUI and
SDPTool and the maintained shell CLI helpers. Source documentation and historical
experiments remain evidence/assets rather than newly modeled deployable tools.

## Cross-system collaboration

[Landscape.design](Landscape.design), registered as `sdp-landscape`, owns only the
**proposed integration context** between these boundaries. Its units are boundary
references, not copies of their internal declarations or deployed containers.
The current SDL checker does not resolve them across files. The name mapping below
and the catalog above provide human traceability; cross-model validation remains
KB-SDL-005 work. Ecosystems have no SDL declarations or containment semantics.

- `SdpTool`, `SdpMcpAdapter`, `CodexClient`, `KanBanCli`, `KanBanTui` refer to
  ProjectGovernance's five systems.
- `SdlFrontend`, `SdlCodeGeneration`, `SdlViewpoints`, `SdlDocumentService`,
  `SdlDevelopmentHost`, `SdlActionRuntime`, `SdlDocumentSnapshot` refer to SDL's seven systems.
- `SduiFrontend`, `SduiCodegen`, `SduiPresentation`, `SduiRuntime`, `SduiNativeHost`
  refer to SDUI's five systems.
- `ExternalCodexAppServer` is external. It is not a library owned by SDP.

The three declared sequences illustrate agent routine selection through MCP,
client design browsing through SDPTool/viewpoints, and KanBanTUI board inspection.
The MCP/client/TUI portions are proposed. Channel messages and their minimal
correlation envelope are abstract collaboration contracts, **not JSON-RPC wire
schemas** or implementation claims. Full payloads, errors, identity and transport
negotiation must be selected in later bounded designs. Internal model READMEs
own current implementation observations and more detailed local scenarios.

The collaboration model deliberately leaves unrelated units unconnected rather
than inventing interactions. Library dependencies are interfaces in local models;
no call means no claim of a network connection. RGS2 retains authority for the
Master/client/Steering role split and unverified app-server/MCP behavior.

## Browse with the existing tools

Registration is explicit in [navigation.json](../navigation.json). Existing
`sdptool` and `sdl-sdui` IDs/default keep their old detailed models. New IDs select
bounded tool designs; they do not replace those authorities or grant native
XFMD ecosystem grouping. No XFMD changes are part of this delivery.

From repository root with prebuilt tools:

```sh
sdptool . discover
sdptool . tree --model sdp-landscape
sdptool . tree --model pg-mcp-adapter
sdptool . view ip --model sdp-landscape --viewer /absolute/xfmd --sdl-tool /absolute/sdl
```

Use the IDs in the table verbatim. `tree` returns current revisions and selection
targets for the existing `select` operation. Documents are produced on demand;
there is no checked-in bulk generated documentation tree.

For a disposable, reproducible preview using the current Go tools:

```sh
SDP_GO=/absolute/path/to/go bash SDP/03--Architecture/Ecosystems/verify.sh /tmp/sdp-ecosystem-review
```

The output directory must be empty and outside the repository. The script builds
the current SDL/SDPTool executables, checks every registered catalog model,
exports AST/canonical form, obtains its navigation tree, generates one revision-bound selection, and generates VP01, VP02, VP06
and VP08 Markdown/Mermaid. Its index links generated bundles. Set SDP_VIEWPOINTS
to another comma-separated selection to inspect additional viewpoints. No model
runs or live app-server/MCP interaction occur. See the plan's evidence record for
actual observed results rather than assuming this recipe has passed.

## Pending work and authority

KB039 covers this initial model catalog. KB040 owns KanBanTUI implementation
planning. KB036–038 retain routine/MCP/client product work; KB-SDL-005 owns real
System/source-set semantics; KB020 retains migration of legacy detailed models.
Templates and consuming-project upgrades are separate release work. These small
initial models support critique and selection, not an exhaustive detailed design
or a promise of seventeen new executables.
