# WCI4-M1 integration fixture Worker handoff

Status: bounded implementation complete and frozen for independent review/native
integration; not whole-stage acceptance. SDP 1.1.1, Worker 2.0.0 and the shared
document workflow loaded/reused. Session0010 S5 recovered read-only; Main owns its
journal, canonical documents, PM/Trace, package/module roots and OS evidence.

## Authority and scope

Coordinator selected WCI4-M1 after M2 delivery `69d0a33` and handoff `90a94b5`.
Working clone `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci4`, observed HEAD
`90a94b5daa4b36635be102858b2ce4a8012150c6` plus concurrent lane changes.
Read ORIGINAL Providers-and-packaging, the final frontend/provider/layout/host API
memos and WCI4-assignment-audit in SDP/05--Implementation/SDUI/Widgets.
Earlier M2 work was superseded and was not edited, tested or process-managed here.
No commits, branch changes, module/canonical/management/other-lane writes.

Only new `SDL/go/examples/previews/**`, root `WCI4-fixture-API.md` and this report
were written. API/path/protocol handoff was sent directly to Main/Noether/James
before implementation. Actual exports consumed; no competing stubs or material
contract departures. Independent review requested from Mendel after component pass.

## Delivered behavior

Actual action-core application and DocumentHost, title SDUI WCI4 Previews,
1100×850; `--nonmodal` selects the existing ordinary Detail window route.
Direct closed-subset SVG, explicit Markdown prose, per-diagram labelled native
fallback, inactive tabs/outer scroll/split, and a composed preview dialog coexist.
Every supplied SVG/renderer is keyed by exact normalized path, including hidden
and closed declarations. Equal Markdown text retains different source identities.

Load and shared Mark button/menu/key invoke real SDL with distinct extended
readOnly result fields. SaveNote/SaveDetail use typed ControlText self echo;
AcceptDetail uses the actual closed DialogAccept schema. Child Save survives Cancel.
No basic-receiver workaround, hidden output, new bridge/runtime state or transport.

Fresh request loan buffers are separate from authoritative resources. Source,
resource and renderer revision conditions invalidate the application Guard;
prepare/publish/abandon use existing Host/Bundle methods. Failed preparation or
publication retains the live bundle. Commands arrange conditions only, never fake
UI gestures or SDL dispatch. Action/renderer/domain counters and host Inspector
are observed, with correlated exact-command barriers and truthful closed state.
Renderer baselines are observed; no hardcoded native call-count requirement.
Known native diagram unavailability does not become a rendered-Mermaid claim.

## Verification

Environment: `go version go1.27.1 linux/amd64`, `GOWORK=off`,
`GOFLAGS=-mod=readonly`; working directory `/tmp/sdp-sdui-widgets/SDL/go`.
No native executable was launched by this worker.

| Actual command | Result |
| --- | --- |
| `go test -race ./examples/previews -run 'TestRequestIdentityAndRealSDLPreflight|TestRequestGuardAndLoanSeparation|TestNativeRequestAdmission' -count=1 -v` | PASS 55.576s, exit0; initial admission, actual bridge/preflight, loan/caption/resize, zero renderer increments. |
| `go test -race ./examples/previews -run 'TestRequest|TestImmutable' -count=1 -v` | PASS 1.188s, exit0; exact request paths/typed SDL, stale authority, supplied/renderer/accessor copies and foreign-root rejection. |
| `go test -race ./examples/previews -count=1 -v` | PASS 244.019s, exit0; all six tests. Fixture Go source before/after inventories identical. |
| `go build -tags desktop -o /tmp/wci4-previews-worker-freeze1 ./examples/previews/cmd/native` | PASS exit0; eight-file fixture/memo inventory unchanged across build. |

The full new-package run additionally proves actual Host rejection of bad closed
SVG identity/content and nil renderer, stale renderer-revision publication,
successful replacement revocation and failed native-resource ticket retention.
Its component event dispatch exercises actual SDL Load/shared Mark to extended
readOnly, nonmodal child Save followed by Cancel, then Accept false/true. These
are Fyne test-driver/component results, not OS input/rendering evidence.
The renderer-copy test deliberately uses a supported deterministic component
backend and proves both renderers actually ran; it does not claim native Mermaid.

Logs, command JSON, exit files and source inventories live in
`/tmp/wci4-fixture-tests`. `initial.log`/exit1 retain an early fixture-test compile
error (Normalize/New signatures corrected only in new tests); it is excluded from
pass evidence. `admission`, `provider-sdl`, `fixture-race`, `build-pilot1` and
`build-freeze` identify separate captures. Prior pilot binary remains immutable:
`/tmp/wci4-previews-worker-pilot1`, SHA256
`7defa0ae7446bdc1a625718120f00f1b8c3a8e697a61e90b6ad7452ad603185c`;
it predates explicit zero initialization of changeCounts and is not final source.

## Exact freeze

Final worker binary `/tmp/wci4-previews-worker-freeze1`, SHA256
`aba9faaa5b6cd61c83802185e65fb2cd2ae89a1b4da149c9e9860f017ff72ffd`.
`build-freeze.buildinfo` records actual `go version -m`, replacements and desktop
build settings. Main owns package/license/module/native recipe evidence.
Manifest `/tmp/wci4-fixture-tests/fixture-freeze.sha256`, SHA256
`1420ef0518da4cc100bec8473847dc6a4f5423cc7189b2abffd4f6669af999b3`.

```text
306e67bcd2287d526a813c105e8bfdcf327d32671bf243ffad3198583515327a  SDL/go/examples/previews/README.md
4b5a8cf1e157a0b6a4e9c58ad860c68c9e1bb107904c15d8b28a0e4027d21061  SDL/go/examples/previews/cmd/native/main.go
1b78e8a4ca54d7cf30fa510a08a4d4f02deab21a332ce8a25edaa48ad9ebea13  SDL/go/examples/previews/controls.go
57ba20e580d0425669262363b7f3ee4efce74082688ad9143b16bba6bfd88cd5  SDL/go/examples/previews/fixture.go
41014264129a99230d42ade7c35107074927d68e1e788d4b46f1f68254df218b  SDL/go/examples/previews/fixture_test.go
55dc69bf7dabfc107a8ad12ffe231b7a2471a5fb774082d7e3a2b459f919bd51  SDL/go/examples/previews/resources.go
e47b320314750577f407c94df909d974aac0365c7c983b596bfca615d12098fd  SDL/go/examples/previews/source.go
2a2baaec7181a51ab79204235058d800d713c26e969a14a8d343f3866847fd5e  WCI4-fixture-API.md
```

## Evidence hashes and remaining limits

| Evidence | SHA256 |
| --- | --- |
| `admission.log` | `ff7c39e81d6a4edae8a8e5d17adfc1cad69f9d78229aadc15275c9cfe18786c2` |
| `provider-sdl.log` | `94a70b0c6e7d44526a12bd6a9b1ba18b9ce06798c282bc130c017f3fa0d30603` |
| `fixture-race.log` | `906ad1b973249dff5a1746364fca00d724b4613b6143214fee2a5a1b1fe1c0d3` |
| `fixture-race.command.json` | `60f115ecc7b85a7a88a8d65ba7a95a70a56c8790e90f986432ba72d8c024bb34` |
| `fixture-race.before.json` | `31a93d8bcf934a36b6b3734d87ff44d1f76258c1083747282b6d5b436b973ae2` |
| `fixture-race.after.json` | `31a93d8bcf934a36b6b3734d87ff44d1f76258c1083747282b6d5b436b973ae2` |
| `build-freeze.command.json` | `74231c0c8dcc9f58b8cadc844c5a7b8abfc9843007b4b16825d4e6e313be7d56` |
| `build-freeze.buildinfo` | `1ea7fdb48b146b792402df3178a4e5d8d40ad7dcf9af228eab9d7c29977ccdc8` |

The shared tree evolved in other exclusive lanes during this work. The fixture
source checks identify this lane exactly; these receipts are not a claim that
all shared dependencies remained frozen for a whole-stage regression. Main owns
current-candidate tests, native binary selection and the full archive.

Remaining integration proof: actual direct image/prose/fallback screenshots and
mounted Accessible observations; pointer/key/menu and modal/nonmodal terminal
receipts; scroll/negative clip, tabs/split retention, copy and parent/reload/window
lifetimes; current all-family/source/export/discovery/package checks; independent
review. No Linux OS screen-reader, native Mermaid, immediate GPU erasure or
standalone .3/helper/package readiness is inferred. No WCI4-M2 release or publish.
No further fixture edits are planned absent a reviewed finding or bounded request.

Session handoff to Main: S5 advanced with a runnable real provider fixture,
fixed source/protocol, targeted six-test race/build pass and this exact freeze.
Canonical provider/design contracts remain accurate; journal/review/native/package
status belongs to the coordinator, not this lane's implementation note.
