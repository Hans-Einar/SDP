# SDUI — documentation and authority

Start with the [SDUI README](../README.md) for current executable scope and [Go commands](../go/README.md) for usage. Only SDUI 0.2 in Go is active.

## Current profile and implementation

| Document | Responsibility |
| --- | --- |
| [language](language.md) | Language profile 0.2, AST and local validation |
| [architecture](architecture.md) | Go packages, lifecycle and responsibility boundaries |
| [go-layout-contract](go-layout-contract.md) | Measurement, relative layout, DIP fonts, clipping and limits |
| [markdown-provider](markdown-provider.md) | Actual Markdown/Mermaid coverage; no full-Markdown claim |
| [runtime-contract](runtime-contract.md) | SDUI/SDL boundary and detailed package-contract links |
| [go-generation](go-generation.md) | Generated models/bindings; same runtime as source execution |
| [requirements](requirements.md) | Requirements, delivery boundaries and evidence references |
| [implementation-plan](implementation-plan.md) | Shared SDL/SDUI phases and milestones with evidence |
| [target-architecture](target-architecture.md) | Selected Go/Fyne direction, reuse and other-repository boundaries |
| [concept1-console](concept1-console.md) | Concept1 source, AST and structural dump formats |

Language, layout and runtime have separate responsibilities. A broader study idea does not extend the parser or runtime contract. Record actual code/contract discrepancies as defects; do not hide them by declaring all code authoritative.

## Design background and dated references

| Document | How to read it |
| --- | --- |
| [layout-language-proposal](layout-language-proposal.md) | Original proposal; language and measurement contracts bound active choices |
| [frame-composition-proposal](frame-composition-proposal.md) | Owner input and recommendations, not a parallel language profile |
| [concept1-relative-layout-study](concept1-relative-layout-study.md) | Dated investigation of Ponsse 882ad7c, not a claim about current Ponsse |
| [renderer-extraction-and-language-direction](renderer-extraction-and-language-direction.md) | Dated worktree/reuse inventory, with subsequent Go direction |
| [handoff-mermaid-extraction](handoff-mermaid-extraction.md) | Historical coordination, not an order for another extracted renderer |
| [handoff-xfmd-sdui](handoff-xfmd-sdui.md) | Superseded FOX handoff and separate document-navigation reference |
| [prototype-widgets](prototype-widgets.md) | Static SVG/HTML gallery predating shared Go layout, not another runtime |

All 17 subject documents in this directory are classified above. Generated [viewpoints](../design/README.md) come from the tool and model facts. [Checkpoint #1](../../SDP/History/checkpoint-1/README.md) and G-phase evidence are dated snapshots; do not update their claims to resemble current contracts. Open work belongs in [KanBan](../SDP/Agents/KanBan/README.md); [KB-SDP-010](../../SDP/Agents/KanBan/active/%23010--Proposal--Document-consolidation.md) tracks editorial consolidation.
