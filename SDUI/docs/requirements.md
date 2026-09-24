# SDUI — requirements and traceability

Updated 2026-09-22. Phase evidence: [G1–G3](../go/README.md), [G4–G6](../../SDL/go/README.md). The table follows the Go port; named Python baseline cases are frozen test data, not active tests.

| ID | Requirement | Owner / evidence | Status |
| --- | --- | --- | --- |
| SDUI-R01 | Named nested frames/groups | parser.Node; nested_groups_and_local_rows_preserved | Structure implemented; former Box profile replaced |
| SDUI-R02 | Horizontal elements/explicit rows | parser.rows; formatting_binds_before_separator | Implemented/tested |
| SDUI-R03 | Symbolic module/member references | ast.Reference; ported_example_bindings_and_instance_paths | Implemented; no external resolution |
| SDUI-R04 | Declarative setHandle | ast.Connection, validate; local_symbol_errors | Locally implemented |
| SDUI-R05 | Unambiguous names/instance paths | normalize; references_forward_reuse_and_cycles | Implemented; runtime handles in G3 |
| SDUI-R06 | Versioned EBNF/AST/source positions | raw_markdown_and_utf8_spans; version_is_exact_no_legacy | Implemented; no formal EBNF-equivalence proof |
| SDUI-R07 | Reject invalid profile semantics | formatting; shapes_ratio_and_relative_dimensions, widget_contracts | Local rules implemented |
| SDUI-R08 | Bounded work without source execution | bounded_work, truncations_and_malformed_inputs_are_structured | Implemented; no hard real-time guarantee |
| SDUI-R09 | Standalone AST without GUI/Mermaid | CLI; cli_source_protection_and_exit_codes | Implemented/tested |
| SDUI-R10 | Same normalized model from SDUI/Concept1 layout | Concept1 reference and Go normalization/layout | Reference structure verified; general React adapter outside profile |
| SDUI-R11 | Visible outer box, title, weights and widget layout in Fyne/SVG/Markdown | Go layout/svg/fynehost, G2 | Implemented/verified |
| SDUI-R12 | SDL object receives handle and updates text/value | Go runtime/bridge, G3/G4 | Implemented with action-core/typed Go functions |

Test names are abbreviated; source files contain full names. [Historical 0.1 verification](../evidence/verification.md) records the retired prototype. Semantic validation does not establish existence of external SDL files/functions.

## Delivery requirements and phase evidence

| ID | Requirement | Phase / evidence |
| --- | --- | --- |
| SDUI-R13 | General frames/regions, nested independently laid-out groups, reusable groups with instance paths | G1/G2 grammar/AST/group preservation, local separators, multilevel layout |
| SDUI-R14 | Markdown strings become source-mapped content widgets | G1/G2 multiline/escaping, width→height, content provider |
| SDUI-R15 | Ancestor-relative dimensions, x:y ratio, one driving scale axis with ratio, relative bounds/spacing, canonical arrows | G1/G2 reject px/two ratio axes; normalization, geometry oracles, overflow |
| SDUI-R16 | Typed runtime/backend ports with stable instances/atomic updates | G3 Go consumer, fake backend, real Fyne trials |
| SDUI-R17 | Coordinated geometry/clipping/printing for SVG, native controls and Markdown | G2 control maps, transforms, visual PDF/GUI checks |
| SDUI-R18 | Explicit SDL adapter port; labeled simulation | G3/G4 bidirectional correlation, revision conflicts, teardown; separate real SDL runtime |
| SDUI-R19 | One active frontend/API; port/remove legacy | G1-M3/G5-M4 target/import inventory, unsupported-version, no fallback |
| SDUI-R20 | Frame Markdown reuses agreed host profile/diagram coverage | G2/G3 profile matrix, links/labels, scroll/resources |
| SDUI-R21 | Per-component formatting before separators; absolute fonts without resize scaling | G1/G2/G3 suffixes on all component types, rows, stable fonts, rewrapping |

R13/R14/R15/R21 are verified through frontend, measured layout and native trials. R19 completed in G5-M4. R20 is bounded by markdown-provider.md: flowchart/graph, not all Mermaid types; scroll rejected. Existing XFMD BoxUI worktree remains separate; G6 adds a separate documentation consumer.

R10 initially means equivalent reference structure/geometry rules, not a general Concept1/React adapter. R11 revised on 2026-09-21: Fyne is the first interactive host, SVG-in-Markdown the documentation view. R12 uses real SDL action-core/bridge in G4; handwritten bucking domain remains a labeled simulation.

Owner clarification, 2026-09-20: R15 permits no source pixel width/height. Root references host layout area; children reference nearest ancestor. Direction-pair order preserves alignment, while formatting uses one canonical shape. Other lengths are also proposed relative. Clarification, 2026-09-21: `{16:9,<->}` fills width and derives height without contain fallback; header/footer stay inside ratio. Absolute fonts do not scale during resize. Shared units are logical DIP measured with Go Regular. Formatting follows components before separators; comma continues horizontally, semicolon starts below the preceding row. Canonical corners: ^<, >^, v<, >v; right-center: ->.

**SDUI-R22:** static GUI dump from shared model with raw Markdown and omitted Mermaid fences. Implemented/tested in go/presentation. Interactive console is planned in the [TUI direction](concept1-console.md).

**SDUI-R23:** static Markdown dump from the same model, with layout overview/renderable content. Preserve tables/lists/code, label widgets static, omit Mermaid as in text dumps. No interaction/new parser. Content shows reading order/nesting, not exact GUI geometry.

**SDUI-R24:** bounded SDUI prototype library draws static SVG button/input widgets in Markdown. The same descriptions may demonstrate editing, pressed/focused appearance and printing current values in local HTML. No SDL calls, network, FOX dependency or new parser. Distinguish prototype placement/demo state from later shared layout/runtime. Delivered as a bounded trial; [scope](prototype-widgets.md), [test/browser evidence](../evidence/prototype-widgets/README.md).

**SDUI-R25:** parse/validate/publish source changes without closing UI windows. Retain last valid model on errors and compatible value/focus on reload; reject stale events. SDL state needs migration/reset rules. Go-function changes rebuild/restart. Implemented/verified G3/G4.

**SDUI-R26:** Go generation shares development models/runtime and separates generated code from handwritten domain functions. Unknown/incomplete SDL execution semantics produce diagnostics. Implemented/verified G5.
