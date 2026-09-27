# MPV2-M1/M2 — SDUI design and navigation evidence

Candidate: 5521b42 plus the MPV2 source/export/navigation delivery diff. No SDL
parser/runtime or SDUI renderer implementation changes were made.

## Results

- Three SDUI 0.2 sources parse/normalize with local-profile validation and produce
  AST, console dump and SVG through the existing Go CLI. No syntax-only bypass.
- The [preview manifest](../../preview/manifest.json) records each source hash,
  executable hash, selected entry, 1600 x 900 viewport and output hashes. A second
  export into /tmp/mpv-repeat produced an identical manifest and output hashes.
- All three final SVGs were rasterized with rsvg-convert and visually inspected
  at 1600 x 900: separate application layouts, header controls, input fields,
  disabled activation, bounded text and visible sample/prototype labels. The
  operator plot remains an explicit provider placeholder. No native GUI test.
- [Additional viewport results](layout-check.json): 1280 x 720 and 1920 x 1080
  exports pass for all three screens. These six checks establish successful layout,
  not exhaustive responsive-design or pixel-equivalence evidence.
- Go package tests passed for parser, layout, svg, presentation and cmd/sdui with
  Go 1.27.1. [Executable build identity](tool.txt) records the tool used.
- The installed gh-sdp 0.1.1 / SDPTool 0.2.1 discovers all three registered screen
  sources and their frames. [Relevant actual tree nodes](navigation-nodes.json)
  preserve the source revisions returned by the producer.
- [On-demand checks](navigation-check.json): each page generated structural
  Markdown through `gh sdp . sdui-preview --model ID --entry page --revision HASH
  --output DIRECTORY`. Commands used actual tree/source revisions and succeeded.
  They do not show SVG through the SDPTool facade; its current service is structural
  Markdown. No XFMD GUI integration test or native sidebar implementation claim.

## Verification commands

From the repository root (the executable was built with GOTOOLCHAIN=local):

```sh
go -C SDUI/go build -o /tmp/sdp-mpv-sdui ./cmd/sdui
python3 experiments/mvp1_sdl/export_ui.py --sdui /tmp/sdp-mpv-sdui
python3 experiments/mvp1_sdl/export_ui.py --sdui /tmp/sdp-mpv-sdui --output /tmp/mpv-repeat
cmp experiments/mvp1_sdl/preview/manifest.json /tmp/mpv-repeat/manifest.json
go -C SDUI/go test ./parser ./layout ./svg ./presentation ./cmd/sdui
python3 experiments/mvp1_sdl/verify_restructure.py
python3 SDP/ProjectManagement/validate.py
```

The ordinary L1 document checker encountered unrelated untracked node_modules
Markdown with a broken upstream link. That directory was not changed. The same
checker was run with its Markdown enumeration restricted to Git-indexed files
(including this staged delivery), retaining its frozen-byte and generated-bundle
checks. The scoped result is recorded in the plan after completion.

## Limits and next work

[Navigation handoff](../../Navigation.md) owns the concrete boundary and first
integrated MVP1 acceptance workflow. Full experimental SDL profile/source-set
support remains KB-SDL-005; broader migration remains KB-SDP-020; owner navigation
feedback remains KB-SDP-032. Screens use sample data and unbound controls; no
Ponsse product code, domain calls, persistence or machine commands were executed.
No new shared SDUI import mechanism, fake aggregate SDL model or hand-authored
viewpoint diagrams were introduced.
