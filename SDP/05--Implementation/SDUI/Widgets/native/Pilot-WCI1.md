# WCI1 native pilot worklog

Session0010 T002, PLAN-SDP-0022 WCI1-M1. This worklog records development
observations. It is not final-candidate acceptance evidence; final runs must name
source/binary hashes and supersede each failed or unverified scenario explicitly.

## Environment

Separate Xvfb display :189 (1000x700x24, no TCP), actual Fyne desktop driver,
Go 1.27.1, libX11/libXtst via x11_input.py. Task-specific XDG configuration/cache
under /tmp/sdui-wci1-native-config and /tmp/sdui-wci1-native-cache. No owner
window or settings are used. Native executable from SDL/go/examples/collections.

## Initial startup

Built the concurrent development candidate with:

```sh
GOCACHE=/tmp/sdp-wci0-gocache go build -tags desktop -o /tmp/sdui-wci1-native ./examples/collections/cmd/native
```

Build passed. Launch failed before window creation: SDL CompileActions reported
NONCANONICAL_FORM at 6:10; main.go line 33 panicked. Therefore no product native
input or screenshot was obtained from this attempt. The fixture owner was sent
the failure for correction and real fixture-test verification. The executable is
an intermediate build, not an accepted candidate; logs remain under /tmp during
iteration. The earlier simple-button tooling smoke proves only the input route.

## Read-only integration findings

The layout Worker and coordinator found potential native inconsistencies before
acceptance: item EnsureVisible width, partially clipped PageUp/PageDown target,
row marker measurement, thumb travel/grab mapping and actual minimum geometry.
The host Worker is correcting these with shared metrics and regression checks.
No native acceptance is inferred from code changes. Required Alt-arrow scrolling,
Tab order/disabled handling and nested viewport pointer routing remain in the
actual-input matrix.

## Second native pilot

The corrected fixture started on :189. Its binary SHA-256 was
`be85562e5c5aaf24ddf66dee38599b93a91561b766bc9452bbb0475a4d94bb84`.
The concurrent source tree was still changing; this is development feedback, not
exact final-source evidence. Initial XDG cache-parent warnings were corrected by
creating the task directories before subsequent launches.

Observed through actual XTest input and fixture logs:

- A single list row click selected entry-001 with zero actions; Enter performed
  one real SDL call and displayed Preview: Entry 001 [entry-001].
- Double-click entry-002 performed exactly one additional action. End selected
  entry-098 (skipping the final separator) and revealed it by scrolling.
- Disclosure of folder-a started request 1 without an action. Controlled failure
  and R produced request 2; Escape canceled it. Late success returned with the
  canceled context and changed neither data nor generation. R produced request 3;
  successful delivery published 24 children with generation 2 and no action.
- A loaded child activated through SDL. Two compatible reloads preserved selection,
  sequence and viewport offsets; activation afterward increased the sequence and
  action count once. A deliberately invalid geometry reload reported native-minimum
  while retaining model revision 3.

OS screenshots from ImageMagick import showed native rows, selection, clipping,
scrollbars and the actual preview. Fyne Canvas.Capture returned black images on
this driver path and is not accepted as visual evidence. Use OS captures in final
verification. A parent label disappeared at negative horizontal position; this
paint/culling observation was sent to the host owner for correction and retesting.

Read-only integration checks also reproduced accepted runtime mutation without
native reconciliation after a stale-result error, a finished load stranded after
failed resize, and resource preparation seeing unclamped offsets. These remain
blocking integration findings until corrected and verified. The fixture close
command bypassed Fyne's OS CloseIntercept; the fixture owner is unifying teardown.

`verify_wci1.py` automates real-input and barrier-controlled checks with raw event
logs, action/check journal, binary identity and OS captures. Its creation is test
tooling, not a product acceptance pass. Final nested and complete-slice evidence
remains pending.
