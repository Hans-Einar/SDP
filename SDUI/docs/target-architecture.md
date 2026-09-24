# SDL/SDUI — target Go architecture

**ID:** SDUI-ARCH-003 · 2026-09-21 · Selected and implemented direction; see phase evidence/profile limits. Replaces ARCH-002's Rust/C-ABI and mandatory FOX/XFMD track. [Checkpoint supplement 07](../../SDP/History/checkpoint-1/07-SDUI-0.2-and-Go-Direction.md) records decisions; [architecture.md](architecture.md) describes current Go code.

**Detailed SDL model, 2026-09-22:** [parser/runtime design](../design/README.md) covers both languages, ports, reload, layout, host and generation. Go design-core validates it; action-core runtime is separate.

## Ownership and dependencies

| Area | Responsibility |
| --- | --- |
| SDL/go/parser | Source, AST, symbols/profile validation; no domain execution |
| SDL/go/runtime | Bounded execution/registered Go functions; no GUI dependency |
| SDUI/go/parser | SDUI 0.2, AST, diagnostics, normalization; single implementation after port |
| SDUI/go/runtime | UI instances, identity, properties, events, bindings, reload |
| SDUI/go/layout | Shared measured layout for interaction/export |
| SDUI/go/svg | Static documentation from shared geometry/explicit state |
| SDUI/go/host/fynehost | Window, widget lifetime, focus/input, UI-thread publication |
| SDUI/go/codegen and SDL/go/codegen | Model/binding code; separate handwritten domain functions |

Separate Go modules with Go 1.26 baseline, verified using Go 1.27.1/Fyne 2.8.1. [Dated implementation snapshot](../../SDP/History/checkpoint-1/11-Go-Implementation-and-Navigation.md).

Parser/runtime import no Fyne, FOX, XFMD or Mermaid. Hosts compose libraries; unbound SDUI designs require no particular SDL implementation. SDL owns domain state; SDUI owns widget identity/UI state. Typed Go interfaces connect them; add binary ABI only for a concrete consumer. libsdui names the library role, not a required .so/C header.

## Measurement and presentation

Preserve relative dimensions, ratios and ancestor references until layout. Root receives host area. `{16:9,<->}` derives height from width; resize changes geometry/wrapping without font scaling. Fonts use logical DIP and Go Regular for shared measurement/SVG. Regions remain inside ratio.

Layout outputs identities, rectangles, clipping, text measurements and resources. Fyne and SVG share geometry/theme; no second Fyne interpretation of language rules. Reuse native controls for text entry, focus and keyboard. One SVG image is not interactive without hit/event handling.

Fyne supports SVG/custom renderers, but actual image/text/clipping coverage must be tested; no pixel-identical rendering guarantee. [SVG images](https://docs.fyne.io/canvas/image/), [custom widgets](https://docs.fyne.io/extend/custom-widget/).

Markdown, including Mermaid, remains desired content. Define bounded provider profiles, measurement, resources and diagram coverage. RichText does not automatically equal XFMD compatibility. Do not build another Mermaid parser. Start with simple Markdown; full coverage is separate. Documentation can embed static SVG without SDUI runtime.

## State, binding and reload

Logical references contain session, instance path and generation. Movement preserves compatible identity; deletion/type changes invalidate references. Separate accepted values/drafts; programmatic updates are not clicks. Runtime checks type/enabled/revision before dispatch. [Runtime contract](runtime-contract.md).

Build/validate candidates before publication. Errors retain last valid model. Preserve compatible values/focus; incompatible changes need diagnostics and explicit reset/migration. Cancel/reject outstanding callbacks by generation/revision; reload must not repeat domain actions.

Running Go runtimes can replace SDL/SDUI models. Initially, Go-function changes require builds/restarts. Process switching/state transfer are later options, not loading modified plugins.

## Generation and porting

Initial generated Go constructs the same models/bindings as file-based development, sharing runtime. Later direct behavior compilation requires conformance against defined execution semantics; do not implement two SDL semantics. Generated files never overwrite handwritten domain Go.

Python frontends/tests provided port evidence. G5-M4 removed active implementations after verified ports; frozen data remain. No SDUI 0.1 compatibility. SDL design-core is separate, now 0.5; V1–V4 replaced older active SDL profiles independently of SDUI versions.

Prior FOX/XFMD/Mermaid work offers reuse evidence, not prerequisites. G6 changes XFMD's document host in a separate worktree/PR; old BoxUI/Mermaid work remains preserved.

## Document navigation — delivered G6

[Navigation design](../../SDL/docs/integration/SDL-Navigable-Viewpoints-Design.md) adds directory-based viewpoints, selected generation, temporary packages and optional background service. PR #38 provides XFMD navigator/main Markdown panels; this document host does not replace Fyne runtime. Reuse SDL projection; reader launch/IPC belongs to adapters. Both sides were implemented/native-tested; [supplement 11](../../SDP/History/checkpoint-1/11-Go-Implementation-and-Navigation.md) identifies profile/branch limits.
