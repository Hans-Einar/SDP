# SDUI in Go

G1-M1–M3 implement SDUI 0.2 lexer/parser and source-positioned AST without GUI/I/O in the parser package. Tagged JSON AST matches Python fixtures. CLI validates the local profile by default. Normalization preserves regions, reuse and instance paths; 93 port cases compare against Python evidence.

Module: `github.com/Hans-Einar/SDP/SDUI/go`; Go 1.26 baseline required by x/text 0.42.0. Verified with local Go 1.27.1 from go.dev, checked against published SHA-256. Build tools are installed outside the repository.

From this directory:

```sh
go test ./...
go run ./cmd/sdui ../examples/concept1-bucking.sdui
```

Parser opens no symbolic SDL refs and runs no callbacks. G5-M4 removed the Python frontend; frozen port fixtures preserve comparisons. [Phase plan](../docs/implementation-plan.md); [evidence](evidence/G1.md).

CLI supports `--format ast|dump|markdown|svg`, `--entry`, `--columns`, `--syntax-only` and `-o`. Flags may precede/follow source. File export publishes atomically after validation and refuses to overwrite source. Console/Markdown match Concept1 fixtures byte-for-byte. G5-M4 removed temporary prototype-svg/html commands; general SVG uses shared layout.

G2 supplies shared geometry (`layout`), SVG (`svg`), Fyne adapter (`host/fynehost`) and bounded Markdown (`markdown`). [Evidence](evidence/G2.md). General export: `--format svg --width 1920 --height 1200`, optionally `--mermaid-renderer /path/to/mmdr --resources DIR`. [Profile](../docs/markdown-provider.md).

Native trial: `go run -tags desktop ./cmd/sdui-fyne -entry bucking ../examples/concept1-bucking.sdui`. Native Go-flag entry points require flags before filenames. Linux desktop builds require OpenGL/X11/C compiler. Ordinary package tests use Fyne's memory driver without a display server.

Combined native-client preview: build `./cmd/sdui-preview`. Pass `-source`,
`-entry`, `-revision` and `-output` (an existing empty caller-owned directory).
It writes entry.md, the exact source.sdui snapshot and sdui.json. `-check` returns
sdptool/0.2 operation sdui-check with prototype/prototype-unbound readiness and
an explicit no-SDL-execution diagnostic. Fyne accepts `-revision` to reject stale
launches before opening a window. XFMD 0.7 bundles these tools without replacing
its user's gh sdp extension. Source/profile syntax is unchanged.
