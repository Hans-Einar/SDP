# WCI4 previews fixture API — selected M1

Coordinator-selected SDL integration Worker, 2026-10-08. Own only new
`SDL/go/examples/previews/**`, this memo and `WCI4-fixture-worker.md`.
M2 delivered `69d0a33`; handoff `90a94b5`. No M2 edits/tests/process management.
Authority: ORIGINAL Providers-and-packaging plus the four final API memos and
WCI4-assignment-audit under `SDP/05--Implementation/SDUI/Widgets/`.
Main owns Session/canonical records, native OS proof, package/dependencies and final
acceptance. Existing parser/provider/layout/host exports are consumed directly;
no runtime preview state, generic renderer framework, bridge or SDL transport change.

## Application and stable paths

Package `SDL/go/examples/previews`, desktop CLI `./examples/previews/cmd/native`.
Title **SDUI WCI4 Previews**, initial 1100×850. `--nonmodal` changes only Detail's
modality; title **SDUI WCI4 Preview Detail**. Source, resource acquisition and
publication are application-owned. SDL actions are genuinely compiled/registered.
Actual pointer/key/wheel/window-close input belongs to Main's native harness.

| Alias | Exact normalized path |
| --- | --- |
| Hero | `page/view/body/workspace/art/viewport/hero` |
| Prose | `page/view/body/workspace/art/viewport/prose` |
| Tail | `page/view/body/workspace/art/viewport/tail` |
| ReportA | `page/view/body/workspace/reports/reportViewport/reportA` |
| ReportB | `page/view/body/workspace/reports/reportViewport/reportB` |
| ReportSVG | `page/view/body/workspace/reports/reportViewport/reportSVG` |
| Note | `page/view/body/aside/note` |
| LoadStatus | `page/view/body/aside/loadStatus` |
| MarkStatus | `page/view/body/aside/markStatus` |
| DetailSVG | `page/detail/detailViewport/detailSVG` |
| DetailDoc | `page/detail/detailViewport/detailDoc` |
| DetailNote | `page/detail/detailNote` |

Tabs: `page/view/body/workspace` with `art` and `reports`. Split: `page/view/body`.
Scroll owners are the two page frames `viewport`/`reportViewport` above and
`page/detail/detailViewport`. Negative Box positions arise from real outer scroll;
negative resource viewBox origin is a separate variant. Shared geometry remains
owned by existing layout/host, including fitting then clipping and split minima.

Toolbar: `page/view/header/toolbar/loadButton`, `markButton`, `openButton`.
Mark command: `page/mark`; menu item: `page/view/header/actions/markItem`;
keyboard: Primary+M. Opener: `page/openDetail`. Dialog: `page/detail`.
Buttons: `page/detail/detailActions/acceptButton`, `cancelButton`, `closeButton`.

## Real SDL contract

| Action | Actual schema/result | Selector and receiver |
| --- | --- | --- |
| Load | Value:text → Value:text | Literal metadata token; LoadStatus |
| Mark | Checked:boolean → Value:text | CommandChecked; MarkStatus |
| SaveNote | Value:text → exact Value:text echo | ControlText; Note self receiver |
| SaveDetail | Value:text → exact Value:text echo | ControlText; DetailNote self receiver |
| AcceptDetail | Note:text → Accepted:boolean, Message:text | DialogFieldValue `detailNote`; closed DialogAcceptResult |

LoadStatus and MarkStatus are **extended readOnly inputs**, using the delivered
M2 receiver guard; no basic-receiver workaround. Outputs use distinct visible
receivers. Open/Cancel/Close are existing local effects. Child SaveDetail commits
to the domain independently of later Cancel. Resource bytes never travel in SDL.

## Actual resources and outcomes

All three SVG declarations, including inactive ReportSVG and closed DetailSVG,
receive exact-path PreparedSVG bindings for `art.Chart.@resource`, ProviderID
`fixture-svg/1`, full lowercase SHA256 and exact dimensions. The symbolic module
`application-supplied SVG` is never opened. Default400×200 montage uses supported
shapes, solid colors and transform; no text/fonts/CSS/images/external retrieval.
ReportA/B have byte-identical Markdown, distinct descriptions and counting renderers
`fixture-diagram-ReportA/1` and `fixture-diagram-ReportB/1`, each Revision `r1`.
Prose-only declarations have no renderer binding. Native direct SVG and prose are
positive routes; native Mermaid remains declared per-diagram fallback, preserving
prose and advertising no rendered-Mermaid fact.

Fresh request loan buffers are separate from the authoritative catalogue.
Mutating a loan never changes captured Guard identity. Changes to staged source,
resource or renderer conditions advance application revision; Guard checks it plus
SDL revision at publication. Every request records its actual identity vector.
Frozen outcome/source/renderer counters are observations, not hardcoded readiness.
The `prepared` event captures the observed cumulative renderer baseline; later
gates/caption/resize/tickets/Commit must add zero. A missing/unsupported route may
short-circuit Render; skipped calls are never proof of renderer-return copying.
Component tests use a deterministic supported backend for that distinct copy proof.

## Condition protocol, frozen for Main

Input is one line; output `{event,data}` NDJSON. Alias spelling is case-sensitive.
Commands arrange conditions only: no SDL invocation, draft edit, tab selection,
scroll, dialog action, drag, clipboard or key simulation.

| Command | Effect |
| --- | --- |
| `state` | Actual published host inspector plus counters/identities. |
| `resource SVG_ALIAS SLOT` | Stage good/alternate/wide/tall/negative-origin/missing/malformed/unsupported/bad-dimensions. |
| `renderer ReportA\|ReportB MODE` | Stage good/error/malformed/unsafe/missing. |
| `renderer-revision ReportA\|ReportB r1\|r2` | Change application renderer identity independently of source. |
| `policy PREVIEW label\|reject` | Stage real source policy. |
| `document MARKDOWN prose\|diagram\|html\|image` | Stage real content; normal Report renderer omitted when diagrams removed. |
| `binding SVG digest\|source\|provider\|extra\|wrong-kind` | One next-request fatal identity fault. |
| `binding ReportA\|ReportB provider\|revision\|nil\|extra\|wrong-kind\|unused` | One next-request renderer identity fault. |
| `reload` | Prepare and publish; reject retains live bundle. |
| `prepare`, `publish`, `abandon` | One candidate maximum; publish consumes even on error. |
| `mutate-input SVG`, `mutate-return ReportA\|ReportB` | Mutate observed caller-owned buffers; absent return buffer reports unavailable. |
| `caption SVG short\|long\|empty` | Checked Label update without provider preparation. |
| `action Load\|Mark\|SaveNote\|SaveDetail success\|error` | Arm next actual SDL action. |
| `accept true\|false\|error` | Arm next actual AcceptDetail. |
| `fail resource-ticket` | Fail exactly next existing PrepareResources ticket. |
| `resize W H`, `parent-hide`, `parent-show`, `close` | Actual host/window lifecycle conditions. |

Change identity between prepare/publish to test actual Guard rejection. Restore
conditions/reload before continuing another scenario. No automatic retry/replay.
Prepared-root mismatch is a component test of actual Previews.Check; no injectable
outcome escape is added to DocumentRequest. Identity negatives are distinct from
correct-digest content failures. Both hidden and closed declarations participate.

## Exact event/state observations

`command-result`: exact `command`, monotonic `commandId`, `status` (`ok`/`error`),
optional `error`, `candidateSequence`, `bundleChanged`, `actionCalls`, published
`source` and `modelRevision` when present. Action/state may precede the barrier;
request `state` afterward for settled inspection. No host outcome is fabricated.

`state.data` retains Bundle.Inspect and adds `actionCalls` (action-name keys),
`changeCounts` and `rendererCalls` (normalized-path keys), actual `pending`,
`candidateSequence`, `candidateHeld`, `stagedRevision`, staged `sourceSHA256`,
`actionsSHA256`, `sourceRequests`, `loanSHA256` by SVG path. `publishedIdentity`
and optional `heldIdentity` each contain `sourceSHA256`, application `revision`,
`resources` (path→provider:digest), `renderers` (path→provider:revision).
`domain`: `commits`, `values` by action; `marked`, `detailCommits`, `detailNote`.
`provider-event`: actual path/calls/mode/sourceSHA256/outputSHA256/bytes/width/height,
optional error. `prepared` includes observed `rendererBaseline`.
`closed`: actual `sessionClosed`, `pending`, `domain`, `actionCalls`.
Other events: ready/state/action/change/dialog-result/error/command/command-result.

Host `previews[path]` supplies frozen kind/status/description/diagnostic/providerID/
revision/sha256 and actual mounted/visible/canvas/title/rect/clip/image/caption/
statusRect/accessibleLabel/accessibleRole when mounted. Closed declarations have
no invented native geometry/accessibility. Main canvas is `main`; modal uses its
parent canvas; nonmodal uses surface path. No Linux OS screen-reader, immediate
GPU-cache disposal, native Mermaid or public rich SVG export claim.

## Evidence allocation

Worker: source/SDL admission; actual host publication and modal/nonmodal composition;
identity/failed reload, copy and root checks; observed zero post-freeze renderers;
targeted new-package race/build and exact hashes in WCI4-fixture-worker.md.
Main: actual image/prose/fallback screenshots, pointer/key/menu/dialog actions,
scroll/split/tab retention, resize/negative clip, resource/native lifetime/copy,
Guard/reload/terminal receipts, all-family integration and package evidence.
Prior family fixtures remain intact. This new fixture alone is not all-family,
package or OS acceptance. Full runnable/protocol details are in its README.
