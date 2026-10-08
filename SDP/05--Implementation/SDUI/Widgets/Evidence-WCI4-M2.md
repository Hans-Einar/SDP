# WCI4-M2 — integrated widget acceptance

Final integrated disposition: **implemented, verified and independently approved**. The bounded correction is committed at
`dcf2a7415ba20559c1babd891c0b65140c2ab934`, following M1 `a026a5d1`.
[Candidate delta](candidate-WCI4-M2-delta.json) pins eight source/test/doc paths;
[raw evidence](native/WCI4-M2-final/README.md) retains exact invocations and failures.

## Corrected behavior

Actual OS pixels exposed a defect missed by the earlier aspect-only assertion:
nonzero SVG viewBox origins produced displaced native artwork. Actual raster
regression also found the pinned decoder treating `scale(s)` as `scale(s,0)`.
Detached preparation now creates a validated zero-origin native representation and
expands uniform scale syntax, preserving original source bytes/digest, provider
identity, budgets and publication Guards. The composition is
`displayScale × rootTransform × originTranslation × childTransform`, as required
by [SVG2 root transforms](https://www.w3.org/TR/SVG2/coords.html#ViewBoxAttribute).
Native resource hashes identify the derived bytes separately. No provider work is
introduced into resize, measurement, gates, Commit or painting. Both original
resource admission and the derived representation have explicit bounds; wrapper
size or derived geometry failure resolves the declared label/reject before freeze.

Independent component review approves all eight exact paths and the source/native
identity boundary. Persistent actual-raster regressions cover positive, negative,
mixed and zero origins; noncommuting root/child transforms; single-argument scale;
empty/self-closing documents; and fit scales 0.5, 1 and 1.885. Worker race runs pass
70.509s (prepared previews) and 34.469s (final origin tests); independent empty-SVG
27-case overlay passes 1.105s. The unchanged candidate-two OS binary fails the new
pixel oracle at its right edge (717 rather than 753.1). Retain that failure as
causal evidence, not a final failure or a successful fidelity check.

## Matching package and consumer

Package three contains ten binaries built by existing packaging/native recipes,
34 compiled modules, 42 third-party notices, pinned GLFW source/patch policy,
connected fixture sources and native-library requirements. Build metadata identifies
dcf2a74 with `vcs.modified=true`; the checkout contains untracked worker artifacts.
Before/after inventories establish exact unchanged product bytes; no clean-VCS
build claim is made. Executables are retained as a local archive, not committed
binary blobs, installed upgrades or a published release.

The packaged protocol passes eight cases/50 actual commands, including .2/.3
composition/discovery, UTF-8 source spans, spaces, stale input and atomic explicit
SVG-export rejection. An earlier setup attempt omitted the explicit fixture path
and executed no commands; its diagnostic remains archived. Generated constructors
are exercised by the corresponding Go suites, not an invented packaged CLI.
The actual supplied XFMD GUI-test binary passes (6.863s), with exec trace proving
all three staged executable paths and unchanged consumer files. Its exact C++
build source is unknown; this certifies supplied-binary compatibility only.

## Evidence applicability

Fresh corrected-source verification consists of full SDUI race, SDL previews race,
all five preview workflows, an additional real-pixel origin workflow, and both
rebuilt native IME recipes. The pixel oracle inspects the >=50%-green-coverage contour against known white,
exact-green solid interiors and three-pixel exterior clip margins, using recorded FYNE_SCALE=1
and actual X window origins. Metadata-only aspect fitting is insufficient.

Twenty unchanged non-preview native workflows retain their exact candidate-two
binary/run identities. The final eight-path correction changes only explicit SVG
preparation and its fixture/tests/docs; it changes no parser, runtime, layout,
bridge, text editor, command, collection, pane, scalar or dependency implementation.
Independent review accepts this bounded applicability, not a claim that these
predecessor executions used rebuilt package-three binaries. Earlier fine-grained
vertical-split, collection-empty/lifecycle, value reentry and text failure/constraint
cases remain explicitly mapped in the [all-family map](WCI4-all-family-closeout-map.md).

The candidate-two combined SDL command timed out in text at 600 seconds. Its
other eight packages pass; the unchanged text package subsequently passes alone
in 488.948s total (487.802s test log), under the same default timeout. The new
previews run replaces the affected old preview proof; unchanged bridge/runtime/
codegen/collections/panes/commands/values/text and SDPTool results remain applicable.
No successful aggregate SDL invocation or established cause of its timeout is
invented. Candidate-two helper IME failed before editing because X focus preceded
mapping; IsViewable readiness within the original deadline fixes the harness,
passing eight checks on the same binary and then requiring rebuilt-helper proof.

## Supported outcome and limits

All sixteen matrix rows and A01–A08 / GAP001–010 retain their selected bounded
contract. Connected .3 fixtures execute through SDL/DocumentHost. Standalone
RuntimeView and XFMD Launch do not gain .3 adapters from static discovery alone.
Direct closed-subset SVG and bounded Markdown prose render; unsupported embedded
native diagrams preserve surrounding prose and use explicit label/reject, without
rendered-Mermaid capability. Actual mounted Fyne Accessible descriptions/status
are proved; Linux OS screen-reader delivery is not. Native acceptance is the light
X11/Xvfb/configured-IBus profile. Dark contrast/font/wrapping fidelity remains
KB-SDUI-004. No rich editor, source-set expansion, other IMEs/platforms, GPU-cache
purge, installation, release or external XFMD gap-register closure is claimed.

The selected baseline `3d265d1` already contains 73 commits ahead of observed
origin/main `9e4c173`. [PR #52](https://github.com/Hans-Einar/SDP/pull/52) preserves that history and remains draft
pending predecessor integration. Widget evidence verifies the bounded WCI work
on that baseline, not all earlier commits. No merge or publication is authorized.

## Recorded final results

- 28 distinct applicable native workflows: eight fresh corrected-binary runs and
  twenty explicitly applicable predecessor runs; 321 checks including teardown.
  All 25 published openings have 25 unique exact terminal receipts. Fixtures
  close with empty pending requests. The native audit names every run and binary.
- Fresh full SDUI race: 248.335s; fresh SDL previews race: 212.731s total
  (208.327s package log). Both before/after inventories match all 1,411 paths.
- Five preview workflows: 60 checks, ten exact results, empty fixture stderr.
  Origin/clip supplement: 39 checks and nine actual pixel records. SDL-root IME:
  seven checks; private helper-root IME: eight checks. Both rebuilt recipes pass.
- The first origin oracle used exact RGB at the outer edge and failed on a
  correctly antialiased root-transform edge. Its image, raw failure and exact
  harness are retained. The reviewed >=50% contour keeps the 1.01px geometric
  tolerance and exact-green inset interior. The old defective image still fails
  under this contour (right edge 718 versus 753.1); its earlier exact-RGB bound
  was 717. No product code changed to resolve that harness measurement error.
- Fresh fixture stderr is empty. IBus wrappers retain service-directory permission
  warnings while actual composition assertions pass. Applicable predecessor runs
  retain their preferences-load/EOF diagnostics. The supplied consumer's raw log
  remains authoritative; no universal zero-stderr claim is made.
- After code verification, two documentation-only changes mark earlier milestone
  status statements as historical. Their before/after hashes are separate in
  final-documentation-delta.json. Tested package/source hashes are not rewritten.

Local corrected artifact SHA256:
`87f480418a2eead424f07cfb9e51afe45a651862993a38f8ec7fc1503d307777`.
This identifies the exact ten-binary package; the receipt records local path and
size. The older a026a5d archive is explicitly superseded.

[Independent integrated review](Review-WCI4-M2.md) approves the complete bounded
assignment; the plan/card/Session record its delivery.
