# Checkpoint #1 — implemented Go, runtime and document navigation

Updated 2026-09-22 after G1–G6. Implementation snapshot for this date; supplements 07–09 retain decision/V-phase foundations. Go tooling does not adopt broader checkpoint candidates.

| Area | Delivered profile / evidence |
| --- | --- |
| SDUI frontend | 0.2 AST/spans, names/formatting/bindings, normalized instances. [G1](../../../SDUI/go/evidence/G1.md) |
| Layout/presentation | Relative sizing, ratio, rows/wrap, DIP fonts, clipping, SVG/Fyne controls, bounded Markdown. [G2](../../../SDUI/go/evidence/G2.md) |
| SDUI runtime | Typed handles/events/batches, accepted/draft/focus, compatible reload/last valid UI. [G3](../../../SDUI/go/evidence/G3.md) |
| SDL frontend | Design-core 0.5 structure/goals/data/Channels/scenarios/plans; 151 frozen tests. [G4](../../../SDL/go/evidence/G4.md) |
| SDL execution | Action-core 0.1 records/actions/Go functions, SDUI bridge, model reload/Go rebuild. [Profile](../../../SDL/docs/profiles/SDL-Executable-Action-Profile.md) |
| Go generation | Typed constructors, separate domain functions, equivalent source/compiled traces and state docs. [G5](../../../SDL/go/evidence/G5.md) |
| Documents | 11 viewpoints, static/light navigator, typed selection, symbols, cache/IPC/leases, real XFMD panels. [G6](../../../SDL/go/evidence/G6.md) |
| Classes | Class-core 0.1 members/roles/multiplicity/aggregation/composition/source-linked diagrams. [Profile](../../../SDL/docs/profiles/SDL-Class-Profile.md) |

G phases are complete within documented limits. No active Python SDUI/design-core/projector remains; frozen fixtures/dated evidence remain. Python still drives some native tests and independent SDP experiments.

## Open and run

[Generated plan](../../../SDUI/design/viewpoints/implementation.md), [viewpoint index](../../../SDUI/design/viewpoints/index.md), [navigator](../../../SDUI/design/navigation/navigator.md), [selected UI state](../../../SDUI/design/runtime-preview/entry.md). SDL tools generate diagrams/facts/reports from validated models, without manually added relations.

[SDL commands](../../../SDL/go/README.md): check/AST, viewpoints/view, service, generation, compiled execution/documentation. [SDUI commands](../../../SDUI/go/README.md): AST/dump/Markdown/SVG/Fyne. Linux desktop needs X11/OpenGL/C tools; headless/core execution needs no GUI. Go 1.26 baseline, tested with Go 1.27.1/Fyne 2.8.1.

XFMD [PR #38](https://github.com/Hans-Einar/xfmd/pull/38), based on main c245fd9 in a separate worktree, supplies panels, explicit window/panel targeting and lease release on switch/close. Tested clicks, two windows, invalid requests, restart and cleanup. Old xfmd-boxui/Mermaid work remained unchanged. Mermaid renders diagrams; SDUI layout belongs to its Go library.

## Continuing limits

- EditAptCell is a labeled Go simulation, not Ponsse machine logic/automatic product migration. [MVP1 corpus](../../../experiments/mvp1_sdl/README.md) remains a candidate exercise.
- Design-core, action-core and class-core are distinct structural, execution and class profiles, not alternate interpretations of one source.
- Mode is operating/allocation context; State/state-machine profile unadopted. Database means persistent source, not necessarily SQL.
- Markdown is bounded; Mermaid flowchart/graph tested with registered backend, not all diagrams/full XFMD content. Scroll rejected; SVG widget placeholder.
- Reload publishes validated models; changed Go rebuilds/restarts, not native code replacement. Persistent state across restarts needs storage.
- Executed domain transactions cannot automatically roll back after UI failure; no automatic replay. Uncertain-delivery leases require explicit cleanup after reader crashes.
- Document-service/XFMD IPC uses Linux private same-UID sockets. Fyne hosts SDUI; FOX widgets, TUI and C ABI are not additional deliveries.

## Traceability

Branches: G1→G2→G3→G4→G6→G5; G6 preceded G5 because G5-M3 consumes projection. One commit per verified milestone, phase pushes, combined PR against sdp-vNow. [History](../../Development-Branch-Stack.md), [acceptance plan](../../../SDUI/docs/implementation-plan.md). No automatic merge.

## Separate repository mismatch

The overall SDP validator reports nine existing Traceability/Relations.yaml ID-format mismatches (ITR/SLC/SPR). Toolkit/Traceability were unchanged from target 9ad4324 at this snapshot. [Evidence](../../../SDL/go/evidence/sdp-validator-baseline.txt). Outside language-port scope; overall contract checks are not green.
