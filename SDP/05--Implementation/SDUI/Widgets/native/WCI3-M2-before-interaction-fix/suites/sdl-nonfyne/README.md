# WCI3-M2 SDL non-Fyne regression completion

Authorized follow-up: run remaining SDL runtime/codegen and other non-Fyne
regressions, preserve logs/exits, do not repeat frozen bridge/values/commands/panes/
text tests, and make no product writes. Reused SDP Worker role.

Candidate HEAD: `d6742742f0968661d15698518df688ffb1d28946`; `go version go1.27.1 linux/amd64`. Working directory: `/tmp/sdp-sdui-widgets/SDL/go`.
GOFLAGS=`-mod=readonly` protects module files, including nested parity builds. Tests
run serially by package (`-p 1`), with race detection and fresh execution (`-count=1`).

## Results

All 22 selected default-build non-Fyne packages completed with exit 0: 13 packages
with tests PASS; 9 no-test packages compiled. No package in the frozen exclusion
set was selected. Runtime and codegen timings are 1.048s and 2.919s respectively.
Codegen includes generated collection constructor/provenance tests in a temporary
module and real source-versus-compiled executable parity; it does not run a Fyne
window or repeat the collection fixture package.

### runtime-codegen

```sh
go test -race -count=1 -p 1 ./runtime ./codegen
```

Exit **0**; command elapsed 5.91s. Log `runtime-codegen.log`,
exit `runtime-codegen.exit`, exact argv/environment `runtime-codegen.command.json`.
Log SHA256 `e42980db4758f8fef9df8dee8807dbc3f1ea582c366589089fe37452fbd7ecdc`.

### remaining-nonfyne

```sh
go test -race -count=1 -p 1 ./broker ./cmd/sdl ./cmd/sdl-compiled ./cmd/sdl-dev ./cmd/sdl-document ./cmd/sdl-gen ./cmd/sdl-simulate ./cmd/sdl-view-request ./cmd/sdl-viewsd ./devhost ./documents ./examples/application ./examples/generatedmodel ./examples/simulation ./parser ./reader ./reload ./snapshot ./sourcegraph ./viewpoint
```

Exit **0**; command elapsed 68.087s. Log `remaining-nonfyne.log`,
exit `remaining-nonfyne.exit`, exact argv/environment `remaining-nonfyne.command.json`.
Log SHA256 `932fbbf0a5d440efaf252efc1dbafdfdaf23d26efce1e0e319404388f3e4177a`.

## Coverage gaps and limits

- No remaining non-Fyne default-build package gap: package inventory was derived
  from `go list -deps -test -json ./...`, including test dependency closure.
- `./examples/collections` is the only default-build SDL package outside both this
  pass and the five coordinator-specified frozen passes. It has Fyne-dependent
  tests and was deliberately not run under this authorization. Main should match
  it to current-candidate evidence or schedule that bounded Fyne regression.
- The following packages exist only with `desktop`; this read-only inventory did
  not build/run them. Their final compile/native evidence remains Main's scope:
  - `./bridge/cmd/smoke`
  - `./cmd/sdl-compiled-fyne`
  - `./cmd/sdl-demo`
  - `./examples/collections/cmd/native`
  - `./examples/commands/cmd/native`
  - `./examples/panes/cmd/native`
  - `./examples/text/cmd/native`
  - `./examples/values/cmd/native`

All SDL `.go`, go.mod and go.sum bytes present at test start were compared after
both runs: unchanged (`source-before.json`, `source-after.json`). No product,
module, branch, commit or management/Session writes. Evidence files live only in
this directory; the owned worker handoff receives this result reference. This
report does not claim Fyne or OS acceptance, nor aggregate M2 gate approval.
