# Widget delivery and acceptance matrix

PLAN-SDP-0021 / PLAN-SDP-0022. Initial baseline at Session0010 T001.
All **pending** entries are obligations, not evidence. Each delivered row needs
candidate identity, reproduction command, actual result and independent review.

| Family | Milestone | Frontend / normalized identity | Runtime / native acceptance | Export / integration | Status |
| --- | --- | --- | --- | --- | --- |
| Shared capability/preflight | WCI0-M1 | Separate profile, layout, provider and host facts; source spans | Missing capability/module/signature rejects before activation; last valid UI retained | Local prototype distinguished from connected readiness | WCI0 detached preparation and WCI1 native publication delivered |
| Tree | WCI1-M1 | 0.3 schema, reused instances, stable IDs | Expand/select/activate separate; loading/error/retry; stale/deleted/canceled requests | Composition/text, explicit SVG policy, generated constructors, SDL navigation fixture | WCI1 delivered; see implementation Evidence-WCI1 |
| List | WCI1-M1 | Typed ordered rows and nonselectable groups/separators | Keyboard/pointer, bounded incremental data, removal and refresh races | Mixed folder/file fixture, same identity across exports | WCI1 delivered; see implementation Evidence-WCI1 |
| Scroll viewports | WCI1-M1 | Existing overflow intent, source spans | Wheel/keyboard, clipping/hit tests, nested focus visibility, resize/clamp, reload | Snapshot offset/extent or explicit fallback; inspect output | WCI1 delivered; see implementation Evidence-WCI1 |
| Tabs | WCI2-M1 | Named reusable pages, enabled/visible | Selection, hidden state retention, removal, focus, no hidden actions | Composition and declared snapshots | WCI2-M1 delivered; Evidence-WCI2-M1 |
| Split pane | WCI2-M1 | Axis, relative proportion and minimum extents | Drag/keyboard, resize, collapse/restore, reload/focus | Shared geometry and explicit snapshot | WCI2-M1 delivered; Evidence-WCI2-M1 |
| Button/toggle | WCI2-M2 | Shared semantic command/icon, tooltip/accessibility | Disabled commands refuse invocation; checked/exclusive state, programmatic mute | Same command through toolbar/menu/keyboard | WCI2-M2 delivered; Evidence-WCI2-M2 |
| Menu/context | WCI2-M2 | Command items/groups/submenus, typed target | Pointer/keyboard invocation, stale target and Escape/cancel | Static labelled declaration, native transient evidence | WCI2-M2 delivered; Evidence-WCI2-M2 |
| Dialog | WCI2-M2 | Ordinary content composition, modality | Accept/cancel/close, parent lifetime, focus restore, domain draft untouched on cancel | Explicit static surface policy | WCI2-M2 delivered; Evidence-WCI2-M2 |
| Checkbox | WCI3-M1 | Boolean schema | Keyboard/pointer, disabled/read-only, no programmatic actions | Typed SDL event, static checked snapshot | WCI3-M1 delivered; Evidence-WCI3-M1 |
| Slider | WCI3-M1 | Finite min/max/step | Invalid bounds, increments, changing/commit, draft rollback | Typed numeric binding and labelled snapshot | WCI3-M1 delivered; Evidence-WCI3-M1 |
| Select | WCI3-M1 | Stable option IDs, empty/disabled options | Changed options cannot silently change accepted domain identity | Option IDs retained through generation and reload | WCI3-M1 delivered; Evidence-WCI3-M1 |
| Numeric input | WCI3-M1 | Typed number with editable draft | Invalid intermediate text distinct from accepted value; increment/decrement | Same numeric binding/range contract as slider | WCI3-M1 delivered; Evidence-WCI3-M1 |
| Single-line input | WCI3-M2 | Shared read-only/placeholder/validation | Preserve draft/revert/revision behavior; focus, Unicode, programmatic mute | Existing input baseline plus new properties | WCI3-M2 delivered; Evidence-WCI3-M2 |
| Multiline input | WCI3-M2 | input multiline property | Edit/select/copy/paste, undo/redo, line breaks, scroll, read-only, UTF-8/IME, explicit commit | SDL Load/Save fixture without widget/parser I/O | WCI3-M2 delivered; Evidence-WCI3-M2 |
| SVG/Markdown preview | WCI4-M1 | Provider/resource identity and accessible description | Disposal, missing provider and explicit fallback | Inspect actual previews; no unsupported richness claim | pending refinement; bounded baseline exists |

Every new family includes positive/negative profile tests, original-source spans,
definition reuse, normalization, layout, composition, static text, declared SVG
support/rejection, code generation and native host evidence. Do not treat Fyne's
headless test canvas as OS keyboard/IME evidence. Use a separate display/user area
for native tests, and capture source/binary identities. Preserve all frozen fixtures.

## Old gap disposition

| GAP-XFMD-SDUI suffix | Current destination |
| --- | --- |
| 001 | WCI1 tree/list, required |
| 002 | WCI2 tabs/split, required |
| 003 | WCI1 viewport, required |
| 004 | WCI2 menu/context/dialog, required |
| 005 | WCI3 basic multiline required; full editor remains optional/unselected |
| 006 | WCI3 typed values/forms, required; transactions application-owned |
| 007 | WCI2 shared commands/button extensions, required |
| 008 | WCI0 preparation and per-slice SDL binding acceptance, required |
| 009 | Existing local definition reuse required throughout; component source sets optional/unselected |
| 010 | WCI4 truthful provider/preview boundary required; broader rich content optional/unselected |

SDUI-011 remains KB-SDUI-004. These are delivery destinations, not closure claims
for the external XFMD register. Native consumer distribution remains WCI4 work.
