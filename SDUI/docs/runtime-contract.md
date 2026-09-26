# SDUI ↔ SDL — implemented runtime boundary

Updated 2026-09-24 after G3/G4. SDUI 0.2 and SDL action-core 0.1 have separate Go runtimes with an explicit bridge. No C ABI, FOX pointer or automatic SDL loader. This document gathers the boundary; detailed contracts each have one home:

| Responsibility | Active contract / implementation | Evidence |
| --- | --- | --- |
| UI session, handles, draft/accepted, property batches | [SDUI runtime](../go/runtime/README.md) | [G3](../go/evidence/G3.md) |
| SDL actions, records and Go registration | [action-core 0.1](../../SDL/docs/profiles/SDL-Executable-Action-Profile.md) | [G4](../../SDL/go/evidence/G4.md) |
| Field bindings and result port | [bridge](../../SDL/go/bridge/bind.go), action-core SDUI bridge contract | [G4-M3](../../SDL/go/evidence/G4.md#g4-m3) |
| UI/SDL model reload and Go restart | [Commands](../../SDL/go/README.md), UI runtime | [G4-M4](../../SDL/go/evidence/G4.md#g4-m4) |
| Generated models using the same runtime | [Go generation](go-generation.md) | [G5](../../SDL/go/evidence/G5.md) |

## Connection

The SDUI parser stores `ref`, callbacks and `setHandle` as data. The host registers module aliases → SDL runtimes, Go function signatures and a typed `bridge.Plan`. Neither parser nor bridge opens ref paths. Action form: `module.Action.@invoke`; `module.Action.setHandle(page.input)` identifies the result receiver. Names alone do not activate domain execution.

Binding plans specify exactly one source per input field: widget draft, event value, typed literal or named context value. Validate signatures, aliases, actions and result widgets before installing handlers. The delivered result port writes a text field to an input handle; it is not a general widget API.

## Identity, ownership and updates

UI handles contain session, instance path, generation and widget kind; they are not native pointers. SDUI owns accepted/draft, focus and UI revisions; registered Go functions own domain state/revisions. External value updates must not silently overwrite dirty drafts. UI batches validate atomically.

User events pass handle/revision, binding and enabled checks. Enter commits draft; Escape restores accepted. Programmatic updates do not repeat user callbacks. Fyne publishes on the UI thread. SDL validates input/output records and correlates calls; rejected/duplicate commands produce observable outcomes. Later UI failures do not roll back executed Go domain actions or trigger automatic replay.

Compatible model reload preserves state according to package contracts; deleted/incompatible widgets invalidate old handles. Invalid candidates retain the last valid model. SDL reload validates registrations and preserves Go-owned domain state. Changed Go code requires build/restart and does not automatically preserve memory state.

## Limits and earlier proposals

This bridge does not execute design-core facts. EditAptCell is an explicit simulation, not a machine/Ponsse runtime, distributed transport or exactly-once guarantee. G3/G4 deliver no general asynchronous widget API, FOX/C ABI or executable SVG producer. The SVG widget is a placeholder.

The original SDUI-RUNTIME-001 proposal remains in Git before R2-M2. The package contracts above supersede it as the active contract; old proposed payloads/methods are not another supported API.
