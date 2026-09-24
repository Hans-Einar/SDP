# SDUI 0.2 — implemented source profile

Updated 2026-09-22. The Go frontend parses SDUI 0.2 and builds an AST. Runtime/Fyne/SVG are separate implementations; parsing never invokes domain functions. The [EBNF](../grammar/sdui-0.2.ebnf) and this profile replace old 0.1 syntax. The [layout proposal](layout-language-proposal.md) also describes future geometry; accepting formatting does not imply implementing its layout behavior.

## Source and structure

Documents start with exact `sdui 0.2;`, followed by optional `ref:` declarations, one or more named component definitions and optional `setHandle` bindings. No legacy profile/fallback. Definitions may reference later component definitions. Consumers select an entry frame; definitions are not mounted automatically. CLI requires `--entry` when multiple frames are available.

```text
sdui 0.2;
mainBody = <<"Group 1", "More text"> {v<}, <button("OK")>>;
page = [header="## Bucking", body=mainBody, footer="Static prototype"]*b {16:9, <->, font=10};
```

`[]` is a frame. `<>` is a nestable widget group; explicit groups remain even with one child. Frames contain frames, groups, Markdown, widget calls and references. Widget groups accept all of these except frames. Empty frames/groups are valid. Component strings are Markdown. `button("OK")` uses plain label text; old-profile `text(...)` is removed.

Comma continues a row horizontally; semicolon starts another row below the entire preceding row. Separators belong to their own list. A trailing semicolon is optional inside containers but required after top-level definitions. Trailing commas and empty rows are rejected. Formatting `{...}` follows its component, after optional `*box`/`*b` and before the separator. Detached formatting blocks are rejected.

In frames, `header`, `body`, `footer` are region roles, at most one of each. Explicit body cannot mix with unmarked body content. Normalization extracts regions; composition order is header, body, footer regardless of source order. In widget groups these assignment names are ordinary names.

## Formatting and names

Locally validated properties: scale/scale-x/scale-y, x/y=content/fill/Nfr, min-x/min-y/max-x/max-y, gap/gap-x/gap-y, padding or four-value padding tuple, align-x/align-y, justify, items, overflow-x/overflow-y, wrap, font, enabled, visible. The [layout proposal](layout-language-proposal.md) lists allowed values. Reject unknown properties, duplicates and conflicting sizing rules. Layout dimensions are relative; `font` is a positive absolute text size. Go layout uses logical display units (DIP); see the [measurement contract](go-layout-contract.md). Dumps do not scale terminal fonts.

Canonical corners: ^<, >^, v<, >v; reversing a direction pair is equivalent. Normalization lowers shorthand to properties while AST preserves spelling. `<->` and `>-<` are distinct complete operators. Ratio x:y applies only to frames; one scale axis or one fill axis may drive it. Reject two driving axes, scale on both, content/fr with ratio, and nonpositive ratios. Parser validation does not prove content fits; separate Go layout solves geometry. Wrap applies to one group row without horizontal fill/fr on the group or children.

Names are case-sensitive ASCII `[A-Za-z_][A-Za-z0-9_]*`; reserved words: sdui/ref/true/false/null/setHandle. Property names may also contain hyphens. Local names are unique per definition. Reuse creates distinct instance paths: `left=mainBody` and `right=mainBody` yield different widget paths. Unaliased references use definition names. Cycles and ambiguous instance paths are rejected. Anonymous groups receive internal source-position-based path segments; public setHandle paths omit these. The Go runtime contract defines compatible identity/reload behavior.

## Widgets and symbolic bindings

| Call | First argument | Other arguments |
| --- | --- | --- |
| button | label: string | callback: symbolic reference |
| input | text: string | value: string, callback: symbolic reference |
| svg | source: symbolic reference | label: string |

The first argument is required, positional or named; subsequent arguments must be named. Reject duplicates/wrong types. Anonymous unbound widgets are valid. Callbacks require widget names; external handles address named widgets. Neither parser nor dump executes SVG sources; dumps show placeholders.

```text
ref: domain "missing-on-purpose.sdl";
# inside a component:
# field=input("Name", callback=domain.object.@changed)
# after definitions:
# domain.object.setHandle(page.form.field);
```

References have form module.object.@member. The alias must exist; external files, objects and methods are not inspected. setHandle must resolve a widget instance through its definition and named ancestors. One binding per SDL object and widget. No arbitrary method calls, imports, interpolation or eval.

## Strings, positions and limits

UTF-8 with Unicode scalar positions and half-open UTF-8 byte ranges. Line/column numbering starts at 1; LF starts a new line, including CRLF. `#` starts comments outside strings. Single/double-quoted strings support n/r/t, backslash, quotes and four hexadecimal digits after u. Reject raw control characters, NUL and surrogate escapes. Triple double quotes preserve multiline Markdown without dedent, escaping or interpolation. Numbers follow JSON-like syntax and finite binary64; no precision guarantee beyond 2^53.

Limits: 256 KiB source, 50000 tokens, 64 syntactic/expanded levels, 2048 source components, 32 widget arguments, 32 formatting rules per component, 8192 expanded components across all document definitions. Errors carry code/message/span; no recovery or partial CLI output. `--syntax-only` also preserves unknown widget/format names without profile approval. EOF is required. Resource limits bound work, not hard real-time execution.
