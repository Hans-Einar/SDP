# Open SDL documentation

Run `sdl-design` when installed in PATH, `./sdl-design` from this directory, or `./SDL/scripts/sdl-design` from the repository root. Otherwise the working directory is irrelevant.

The launcher uses **prebuilt** programs. It does not compile Go, start a daemon or open the pregenerated full export. It:

1. Reads this project's `SDUI/design/architecture.design` using the registered SDL tool.
2. Generates navigator/overviews in a private temporary directory: 14 Markdown pages and one manifest, no detailed diagrams.
3. Opens `index.md` in XFMD's main panel and `navigator.md` in its side panel.
4. Registers tool, source, project and Mermaid renderer. `sdl-view://` clicks run the SDL tool against source as it exists at click time.
5. Removes startup files when the window closes. XFMD owns detailed-view lifetimes.

The navigator is a startup snapshot; restart to refresh its inventory after structural model changes. Detail selection always validates current SDL. Invalid sources report diagnostics and retain the last view.

Defaults and optional overrides:

| Variable | Default |
| --- | --- |
| SDP_SDL_TOOL | `$HOME/.local/lib/sdp/sdl` — prebuilt SDL CLI |
| SDP_XFMD | Sibling `xfmd-sdl-navigation/build/xfmd` |
| SDP_MMDR | Sibling `mermaid-rs-renderer/target/debug/mmdr` |
| SDP_SDL_SOURCE | This repository's `SDUI/design/architecture.design` |

`--help` displays usage. One file argument selects another SDL source. Missing executables produce explicit errors, without automatic builds/installations. XFMD's `--navigator` and `--sdl-tool` require file arguments; the launcher supplies them. It does not automatically select the ordinary installed XFMD, which may lack PR #38 navigation integration.

The full `SDUI/design/viewpoints` export is the saved G5 verification/export result, not the daily source-based browsing entry point.
