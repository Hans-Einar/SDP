# Ecosystem architecture and executable SDL model plan

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0008 |
| project | SDP |
| state | active |
| PlanType | ArchitecturePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDP, SDL, SDUI, SDPTOOL |
| source | Owner request 2026-09-29; KB-SDP-039 |

## Outcome and authority

The owner authorizes establishing ecosystem/system folders under SDP/SDL and
modeling the repository's tools in SDL in this session. ProjectGovernance groups
SDPTool, the MCP adapter, the Codex client and KanBan tools; SDL and SDUI group
their parser/compiler/runtime/presentation tools. Model actual implementation
and proposed boundaries separately. This is architecture/model authoring, not
implementation of the new products, migration of legacy models, or publication.

## Baseline and choices

Start from 3b4cff6. The current Go parser supports design-core 0.5 standalone
models, not System declarations or source-set linking. The untracked sourceinput
package is unrelated unfinished work and must remain untouched. Each new System
gets a bounded, independently valid System.design and a README with source
provenance, proposed/current status, public contracts and actual binary/library
mapping. Ecosystem membership is directory/index metadata, never a fake Container.
No handwritten concatenation, parser fork, or duplicate detailed legacy authority.
Existing authoritative models and pinned evidence remain in place. This delivery
supersedes the three coarse system groups only for new architecture organization;
KB-SDP-020 and KB-SDL-005 retain full migration and language support respectively.

## Phases and milestones

| Phase | Milestone | Acceptance | Status |
| --- | --- | --- | --- |
| E1 | E1-M1 | Source inventory, boundary decisions, active card and updated source convention | Delivered |
| E2 | E2-M1 | ProjectGovernance models: SDPTool, MCP, client, KanBan CLI and TUI | Delivered |
| E2 | E2-M2 | SDL and SDUI tool models grounded in current Go packages and commands | Delivered |
| E3 | E3-M1 | All model entries pass Go parser/semantic check and generate AST/viewpoints; reproducible catalog validation and navigation index | Planned |
| E3 | E3-M2 | Review scope/authority, record source hashes and evidence, update cards/traceability and handoff remaining work | Planned |

## Git policy

Use sdp/ecosystem-models from the completed research branch; commit each meaningful
milestone, push at phase completion. Existing owner authorization covers these
pushes and a draft PR against main. Do not merge or release. Keep unrelated
SDL/go/sourceinput and Node files outside commits. Do not modify XFMD or gh-sdp.

## Verification and completion boundary

Inspect source entrypoints and packages; do not infer implementation from names.
Use the real Go SDL CLI for each independent model and selected/all static Mermaid
viewpoint outputs. Generated Markdown must come from the tool. Keep authored
models separate from temporary exports and evidence. Check catalog uniqueness,
links, metadata and ledger history. Review status claims and the public boundaries
against studies and code. No GUI/native behavior, full source-set support, real
MCP/app-server run or implementation-drift prevention is claimed by these checks.

## Related work

- [RGS2 study and synthesis](../../02--Requirements/RoutineGovernance/StudyPlan.md)
- [Source organization](../../SDL/README.md)
- [KB039](../../KanBan/active/%23039--Change--Ecosystem-and-system-models.md)
- [KB020 migration](../../KanBan/backlog/%23020--Change--Shared-design-source-organization.md)
- [KB-SDL-005 source sets](../../KanBan/backlog/%23005--SDL--Change--System-and-source-sets.md)
- [KB040 KanBan TUI](../../KanBan/backlog/%23040--Proposal--KanBan-TUI.md)

## Execution evidence

Evidence and actual milestone outcomes will be linked here as delivered.

- E1-M1/E2-M1: five ProjectGovernance models and source maps delivered; Go structural checks, AST and viewpoint generation pass. New capabilities remain proposals.

- E2-M2: seven SDL and five SDUI system models delivered, grounded in tracked Go packages/commands. Independent model checks, AST and viewpoint generation pass; packaging is not extracted.
