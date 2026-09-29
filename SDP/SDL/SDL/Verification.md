# SDL ecosystem model verification

| Field | Value |
| --- | --- |
| Plan / milestone | PLAN-SDP-0008 / E2-M2 |
| Baseline | `ac50d52a322d7d981ee317325343693a21615a0d` plus the new uncommitted SDL ecosystem catalog |
| Toolchain | `go version go1.27.1 linux/amd64` |
| Profile | design-core 0.5 |
| Evidence type | Source inspection, model validation and document generation |

## Actual commands and results

The worker built the repository's existing Go SDL CLI; no parser changes were
made. `go` was absent from the shell PATH, so the installed toolchain was selected
explicitly:

```sh
(cd SDL/go && /home/warloc/.local/share/sdp-toolchains/go1.27.1/go/bin/go build -o /tmp/eco-sdl-cli ./cmd/sdl)
```

For each listed system, the worker ran:

```sh
/tmp/eco-sdl-cli check SDP/SDL/SDL/<System>/System.design
/tmp/eco-sdl-cli ast SDP/SDL/SDL/<System>/System.design
/tmp/eco-sdl-cli viewpoints SDP/SDL/SDL/<System>/System.design --format static --output /tmp/eco-sdl-validation/<System>/views
```

| System | Parse and semantic check | AST output | Generated static bundle files |
| --- | --- | --- | --- |
| ActionRuntime | Pass | Pass | 54 |
| CodeGeneration | Pass | Pass | 48 |
| DevelopmentHost | Pass | Pass | 48 |
| DocumentService | Pass | Pass | 58 |
| DocumentSnapshot | Pass | Pass | 48 |
| Frontend | Pass | Pass | 48 |
| Viewpoints | Pass | Pass | 48 |

All commands exited zero. Exact source hashes, command arguments and CLI binary
identity are recorded in [Validation.json](Validation.json). The temporary AST
and generated Markdown/Mermaid bundles were inspected as outputs and are not
checked-in hand-maintained documentation. No external SVG renderer was requested.

The initial draft was corrected through the real formatter and validator:
canonical declaration/statement order is required; ordered scenarios require
closed payload contracts. The two scenarios therefore explicitly describe closed
metadata projections, not the full runtime/API payload or serialized ABI.

## Scope and authority

Current source mappings are in each system README and the ecosystem index. The
inventory covers all ten tracked SDL/go/cmd entrypoints plus bridge/cmd/smoke.
Native examples and smoke programs are fixtures, not additional proposed products.
The command inventory was derived with `git ls-files`, excluding unrelated
untracked `SDL/go/sourceinput`.

This verifies that the seven authored catalog models pass the current Go frontend
and produce source-derived viewpoints. It does not rerun runtime integration,
spawn XFMD, launch the broker, render native Fyne windows, execute domain handlers,
validate codegen parity or establish independently released system boundaries.
`ExistingCapability` statuses represent mapped implementation observations, not
new behavior verification. Prior detailed design and evidence remain in place.
