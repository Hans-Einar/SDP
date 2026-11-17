# DocumentHost prototype launch context

The standard SDPTool launcher selects DocumentHost for SDUI 0.3 or when an explicit
context is supplied. SDUI 0.2 without context retains the legacy local prototype.
This is a UI prototype: symbolic SDL callbacks and application actions are unbound;
local field edits do not execute FreeCAD, a solver, or another domain application.
Fields declaring symbolic callbacks retain their existing unbound commit behavior;
editing a draft does not claim the callback ran.

```sh
sdptool sdui check --source Screen.sdui --entry Main --context Main.context.json --json
sdptool sdui run --entry Main --context Main.context.json Screen.sdui
```

`gh sdp` exposes the same commands. A context is never discovered automatically.
Static `sdui preview` keeps its structural purpose and rejects context arguments.
The equivalent host actions are `sdui.source.check` and `sdui.source.run`, with
optional `context` and `contextRevision` parameters. Connected callers supply
absolute source/context paths, as with existing source actions.

## Explicit fixture input

```json
{
  "schema": "sdui-prototype-context/1",
  "viewport": {"width": 1680, "height": 1150},
  "choices": {
    "Main/settings/loads/pattern": [
      {"id": "uniform", "label": "Uniform", "enabled": true},
      {"id": "asymmetric", "label": "Asymmetric", "enabled": true}
    ]
  },
  "svg": {
    "Main/preview": {"file": "assets/frame.svg", "width": 640, "height": 500}
  }
}
```

Keys are exact normalized instance paths, including reused/hidden descendants;
short widget names or labels are not aliases. Every selected entry's select must
have explicit options. Unknown paths or wrong control kinds reject. An explicitly
empty options list is distinct from an absent provider. Other root entries do not
implicitly share this context. Values and enabled flags are fixture data; no
application or network provider is synthesized. Tree/list collection providers
are not supplied by this context profile.

SVG files are local paths relative to the context file's directory. Absolute paths,
parent traversal and symlink escapes reject. Explicit SVG data is associated with
the source model's reference without opening or executing its SDL target. The
context limit is 1 MiB, individual SVG limit 4 MiB and combined SVG limit 32 MiB.
Viewport dimensions must be finite positive values at most 32768; absent viewport
uses 1280×800. They specify initial native content size; user resizing is retained
on subsequent reloads.
The context does not alter source layout semantics or native content minima.

## Readiness and identity

For DocumentHost, `check` returns `status: preflight`, `host: document` and
`nativeLayout: pending`. It validates the source, explicit fixture/provider inputs,
capabilities and semantic preparation without opening a window. It does **not**
certify native backend rendering, native geometry or display availability. Launch
performs those checks using the actual DocumentHost before publishing the UI.
Legacy 0.2 checks retain their existing `prototype`/`prototype-unbound` results.
Consumers must treat the new preflight status as a separate result, not assume
that all successful source checks imply native geometry acceptance.

`revision` identifies the source. `contextRevision` identifies the context and
supplied resource contents together; pass it back as `--context-revision SHA` to
refuse changed fixtures between check and launch. It requires `--context`.
`--revision SHA` continues to pin the initial source. A missing, stale or invalid
input causes a diagnostic before that candidate is published.

With watch enabled, source/context/resource edits are observed and prepared as
coherent snapshots. Native publication occurs on the UI thread, retains compatible
state and keeps the last valid UI when a replacement fails. Context/resource file
reads happen before UI dispatch; native guards compare the latest observed
candidate identity. This is polling and guarded publication, not an atomic lock
on files being edited by arbitrary external programs. Pins apply to initial
launch; later watched changes are intentional. `--watch=false` disables reload.
