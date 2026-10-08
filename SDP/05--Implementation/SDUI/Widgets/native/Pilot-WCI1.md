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

## Reproducible pilot checks and new findings

The basic actual-input harness passed ten checks on the second pilot binary.
The old binary failed empty-root first-Tab ordering (preview was focused before
the tree). The next fixture build fixed ordering: empty-root Escape, late delivery,
R restart and paused reload passed, but R after reload produced no native event.
The native stderr reported failure to focus an object not yet in the canvas;
forwarding-wrapper identity and focus preparation are being corrected.

The nested third pilot passed inner-to-outer wheel chaining and reverse scrolling
to zero. Its thumb assertion failed because the outer frame thumb occupied the
same screen strip as the inner list thumb: dragging at (993,100) changed outer
Y to 112 rather than inner offset. The input route worked; the thumbs were not
separately reachable. A bounded shared-gutter geometry correction is required;
the acceptance expectation is retained. Native teardown now reported a closed
Session and no pending provider barriers. All these results remain pilot feedback
and must be repeated against the final integrated candidate.

## Fourth pilot corrections

The rebuilt pre-gutter candidate passed the seven empty-root checks, including
actual R after compatible reload and old-completion rejection. The expanded basic
native run passed sixteen checks: prior selection/action/load/reload checks plus
Alt-Right without selection/action changes, tree parent navigation, selectable
branch R recovery, pointer Retry and complete teardown. OS inspection confirms
negative-X labels now remain partially visible instead of disappearing.

One launch failed because the separate Xvfb process had ended; GLFW reported no
available platform. Restarting the test display with -noreset corrected the test
environment. It is not recorded as a product defect. Final gutter/native candidate
verification still supersedes these evolving-source pilot results.

## Gutter and visible-height pilot

The `fc42424fcfc537f2a82a224ce2b41390d73e9d5376ecae34a3d0be16d3a9a4da`
binary passed all twenty expanded nested checks in
`/tmp/sdui-wci1-native-gutter-paced`. Separate inner and outer scrollbar drags,
visible-height PageDown, background wheel routing, stable-ID reorder, shrink
clamping, seven failed candidate stages and valid recovery passed. OS inspection
confirmed separate gutter strips and the fixed sibling.

An earlier background-wheel run used an instantaneous pointer move and wheel;
Fyne consumed its cached previous pointer position. The harness now moves first
and permits 100 ms native input delivery before wheel/drag gestures. This is
input pacing, not a provider/race synchronization assumption. The corrected
run passed. Final evidence will use this explicit pacing.

Independent review reproduced a remaining recovery defect: a long provider
diagnostic could violate the default horizontal error policy and prevent its
own error state from publishing. The host correction must fit only status
presentation, retain the full bounded runtime error and keep the R recovery
hint visible even at the admitted minimum width. Final acceptance awaits this
correction and final-candidate checks.
