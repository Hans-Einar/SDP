# WCI4 connected provider previews

This is an actual SDL/DocumentHost application for the reviewed bounded preview
routes. Run native input through the real window; stdin only arranges application
conditions and publication/lifetime barriers. No product runtime, bridge, transport
or provider registry is implemented here.

```sh
cd SDL/go
GOWORK=off go test -mod=readonly -race ./examples/previews -count=1
GOWORK=off go build -mod=readonly -tags desktop -o /tmp/wci4-previews ./examples/previews/cmd/native
/tmp/wci4-previews
/tmp/wci4-previews --nonmodal
```

The main title is `SDUI WCI4 Previews`, initial size 1100×850. Detail is
`SDUI WCI4 Preview Detail`. `--nonmodal` changes only its source modality. The
canonical module roots own native/IME dependencies; this fixture changes none.
The coordinator owns native execution, screenshots, package proof and final tests.
Fyne test-driver admission is component evidence, not native OS rendering proof.

## Source and application ownership

`source.go` has the actual `.3` source and canonical action-core program.
`ref: art "application-supplied SVG"` is symbolic and never opened. Hero,
inactive ReportSVG and closed DetailSVG each receive an explicit exact-path
`PreparedSVG` binding with matching `art.Chart.@resource`, provider identity,
full digest and dimensions. The default 400×200 SVG uses closed-subset shapes,
solid colors and a transform, with no text/image/external resource dependency.
Hero and Detail may have equal bytes, but the host owns independent native objects.

The Art page contains the image, explicit Markdown prose and a long report in an
actual outer scroll frame. Reports has two byte-identical Markdown documents
with distinct descriptions/renderer identities and an inactive SVG. Their diagrams
use declared per-block fallback; surrounding prose remains. The native composed
Mermaid route is unsupported and is never called rendered by this fixture.
Detail composes its own image, prose, editable child Save and Accept/Cancel/Close.

Load and the shared Mark command return text to distinct **extended readOnly**
fields. Mark's button, menu and Primary+M use one actual SDL action. SaveNote and
SaveDetail use typed ControlText exact self echo. AcceptDetail uses DialogFieldValue
`detailNote` and the closed boolean/text DialogAccept result. A child Save is a
real domain commit and survives later Cancel. Preparations invoke no SDL action.

Fresh request buffers are disposable loans, separate from the application catalogue.
`mutate-input` changes only the last request's loan. Staged resource/source/renderer
changes advance the application's revision, captured by Guard. Request-renderer
closures capture their own mode; changing conditions never changes an existing
prepared outcome. Renderer calls are counted only when Render actually runs.
The `prepared` event records its observed cumulative baseline; do not assume two
initial calls. Gates, caption changes, tabs, scrolling, resize and publication must
add zero calls. `mutate-return` reports an error if no output buffer was observed.
The fixture's deterministic supported-backend component test proves renderer-buffer
copying; that test is expressly not native Mermaid evidence.

## Stable paths

| Alias | Path |
| --- | --- |
| Hero / Prose / Tail | `page/view/body/workspace/art/viewport/hero`, `/prose`, `/tail` |
| ReportA / ReportB / ReportSVG | `page/view/body/workspace/reports/reportViewport/reportA`, `/reportB`, `/reportSVG` |
| Note / LoadStatus / MarkStatus | `page/view/body/aside/note`, `/loadStatus`, `/markStatus` |
| DetailSVG / DetailDoc | `page/detail/detailViewport/detailSVG`, `/detailDoc` |
| DetailNote | `page/detail/detailNote` |

Each shortened sibling suffix above has the same complete parent path as its row.
`Targets` in source.go contains all complete aliases. Tabs is `page/view/body/workspace`;
split is `page/view/body`. Toolbar is `page/view/header/toolbar/{loadButton,markButton,openButton}`.
Menu item is `page/view/header/actions/markItem`; command is `page/mark`.
Dialog buttons are `page/detail/detailActions/{acceptButton,cancelButton,closeButton}`.
Only ordinary source declarations are used; no hidden receiver or action simulation.

## Condition protocol

Input: one command per line, at most 65536 bytes. Output: NDJSON `{event,data}`.
Every processed line gets `command-result` with exact `command`, monotonic
`commandId`, `status` (`ok`/`error`), optional `error`, `candidateSequence`,
`bundleChanged`, actual `actionCalls`, and published `source`/`modelRevision` if any.
Action/change/host-state events can precede the correlated barrier. `state` after
that barrier obtains the actual settled inspection. No command simulates a gesture.

| Command | Condition only |
| --- | --- |
| `state` | Live host inspection plus application counters. |
| `resource Hero\|ReportSVG\|DetailSVG SLOT` | Stage good/alternate/wide/tall/negative-origin/missing/malformed/unsupported/bad-dimensions. |
| `renderer ReportA\|ReportB MODE` | Stage good/error/malformed/unsafe/missing. |
| `renderer-revision ReportA\|ReportB r1\|r2` | Change application identity, independent of source text. |
| `policy PREVIEW label\|reject` | Stage that source declaration's policy. |
| `document MARKDOWN prose\|diagram\|html\|image` | Stage real Markdown source; normal Report bindings are omitted when diagrams are removed. |
| `binding SVG digest\|source\|provider\|extra\|wrong-kind` | Inject one next-request fatal identity fault. |
| `binding ReportA\|ReportB provider\|revision\|nil\|extra\|wrong-kind\|unused` | Inject one next-request renderer identity fault. |
| `prepare`, `publish`, `abandon` | At most one detached candidate. Publish consumes it even on failure. |
| `reload` | Prepare then publish; failed admission preserves the current bundle. |
| `mutate-input SVG`, `mutate-return ReportA\|ReportB` | Corrupt only observed caller-owned loan/return buffers. |
| `caption SVG short\|long\|empty` | Existing checked Label update, without resource preparation. |
| `action Load\|Mark\|SaveNote\|SaveDetail success\|error` | Arm the next real action. |
| `accept true\|false\|error` | Arm next actual AcceptDetail. |
| `fail resource-ticket` | Fail exactly the next PrepareResources ticket; it never invokes a renderer. |
| `resize W H`, `parent-hide`, `parent-show`, `close` | Actual window/host lifecycle conditions. |

Alias matching is case-sensitive. Configuration faults fail without invoking SDL.
Restore conditions and reload before resuming a scenario after staging identities.
There is no auto-retry/replay. Source-root mismatch is tested with actual immutable
provider Check, not by adding an injectable outcome to DocumentRequest.

`state.data` retains the host inspector verbatim and adds:

- `actionCalls`: SDL action name → count; `changeCounts`: normalized field path → count.
- `rendererCalls`: normalized Markdown path → observed cumulative Render calls.
- `domain`: `commits`, `values` keyed by action, `marked`, `detailCommits`, `detailNote`.
- `candidateSequence`, `candidateHeld`, `stagedRevision`, staged `sourceSHA256`,
  `actionsSHA256`, `sourceRequests`, and `loanSHA256` keyed by normalized SVG path.
- `publishedIdentity` and optional `heldIdentity`: `sourceSHA256`, application
  `revision`, exact-path `resources` (provider:digest) and `renderers` (provider:revision).
- `pending`: actual outstanding collection requests observed from the session.

`provider-event` records real renderer invocation path/count/mode, source/output
hashes, byte count/dimensions and optional error. `action`, `change`, `dialog-result`,
`error`, `prepared`, `ready`, `state`, `command`, `command-result`, `closed` retain
separate meanings. `closed` contains actual `sessionClosed`, `pending`, `domain`
and `actionCalls`; it is not a fabricated empty teardown acknowledgement.

Host `previews[path]` supplies prepared kind/status/description/diagnostic/identity,
actual `mounted`, `visible`, canvas/title/rect/clip/image/caption/statusRect and
mounted Accessible label/role when present. Closed declarations have prepared
outcomes but no invented window or mounted accessibility. Linux OS screen-reader
delivery and immediate GPU cache disposal are not claimed. Negative Box positions
come from real outer scrolling; negative viewBox is a separate resource variant.
