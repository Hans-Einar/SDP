# WCI1-M1 — collections, scrolling and native publication

Status: implemented, integrated, verified and independently reviewed at the
bounded WCI1 level. Full KB-SDUI-003 inventory remains in progress.

## Candidate and integration

Phase branch `sdui/widgets-wci1`, delivery `c39b330`, prerequisite `d1c5c88`; the 88 delivered files are in
[candidate-WCI1.json](candidate-WCI1.json). The raw reviewed inventory SHA-256 is
`171dd5e724bfa3f24639a606c4d7067d8fcf67a6ce2308e0420ee955c58efb9f`.
The native binary was independently rebuilt byte-identically:
`d97eac538dbb80816da863e07e3ec4aa4bb6ef0416b62c2e37211591a8badbd8`.

The coordinator compared all original targets with the isolated prerequisite
baseline before integration. Only SDPTool/sdui.go diverged: its existing KB005
Combined/source-bundle changes were retained alongside profile-aware WCI1 changes.
All other implementation bytes match the reviewed clone. Unrelated governance
files and frozen evidence were preserved. No merge, release or publication.

## Delivered behavior

Exact development profile 0.3 adds typed tree/list, stable item identity,
separate selection/activation, bounded provider data, explicit retry/cancel and
late-completion rejection. Runtime owns offsets and revisions; shared measured
layout owns clipping, nested scroll routing and separate scrollbar gutters.
Native host stages complete document bundles and guards publication against
source/window/state/provider/SDL changes. Failed preparation keeps the old UI.
The actual SDL action-core fixture routes stable item IDs to preview text.

Generated 0.3 constructors retain profile/provenance and compile; valid 0.2
outputs remain stable. Structural exports identify provider absence; collection
SVG explicitly rejects with source diagnostics. Standalone preview never invents
providers or claims connected readiness. See lane reports for exported Go struct
additions and the intentionally tightened SkipControls inventory requirement.

## Evidence and acceptance

| Contract rows | Evidence |
| --- | --- |
| A01–A03 | Frontend/profile/provenance, compiled generated constructors and exact 0.2 tests; frontend report |
| A04–A07 | Runtime validation, atomic replacement, controlled load revocation, stable targets and typed events; runtime report |
| A08 | Actual native basic/empty input, separate select/activate, retry/cancel, stable sequence across two reloads |
| A09–A11 | Layout geometry/inset/overflow tests; native nested/lifecycle scroll, separate thumbs, visible-height paging, hide/disable, resize and clamp |
| A12–A13 | Actual SDL bridge/fixture tests and native preview actions; engine retention and callback target/revision guards |
| A14–A15 | Native failed-candidate stages and recovery; host lifecycle/resource/publication tests; teardown and late delivery |
| A16 | Source-linked export rejection, structural metadata and actual-profile consumer tests |
| A17 | Full integrated SDUI race suite, SDL affected race checks and full original SDPTool suite passed |

From the integrated original workspace on Go 1.27.1 linux/amd64:

- SDUI/go: `go test -race ./...` — PASS, including host (24.871s).
- SDL/go: `go test -race ./bridge ./runtime ./codegen ./examples/collections` — PASS.
- SDPTool: `go test ./...` — PASS, including retained governance and KB005 integration.

[Native artifacts](native/WCI1-final/README.md) retain 54 passing actual-input
checks, raw logs and OS screenshots on the exact binary: basic 16, empty 8,
nested 20, lifecycle 10. All four stderr files are empty. Unit/headless tests
cover adversarial states; they are not presented as OS-input evidence.

Independent reviewer `01a11854-523d-7c11-adf4-f622b19b68c6` checked all 88 hashes,
reviewed the native logs/captures, independently rebuilt the binary identically,
ran full SDUI race and targeted SDL tests, and approved with no remaining blockers.
Original workspace reconciliation was subsequently verified with the tests above.

## Findings, limits and next stage

[Pilot history](native/Pilot-WCI1.md) retains initial failures and corrections.
The final review correction fits only recovery presentation: full bounded error
text remains in runtime, compact Retry/R stays visible, and data-row overflow
policy remains unchanged. The final minimum-width regressions include fonts
14/20/28 and deep indentation. Hidden adapters retain cached diagnostic rows;
actual hidden painting and widget state, not cache absence, establish visibility.

Native evidence is X11/Xvfb without a window manager, not WM/IME conformance.
No collection SVG renderer, source provider I/O, richer editor or final package
publication is claimed. WCI2 panes/commands, WCI3 typed/text controls and WCI4
provider/package/integrated acceptance remain required. WCI1 pilot findings have
been reconciled into the independently reviewed WCI2 stage contract.
