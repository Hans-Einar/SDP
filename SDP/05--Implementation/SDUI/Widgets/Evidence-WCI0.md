# WCI0-M1 — detached preparation and legacy admission

Status: implemented, verified at service/adapter level and independently reviewed.
This is not a new-widget, atomic native publication or whole-card completion claim.

## Candidate

Phase branch `sdui/widgets-wci0`, implementation commit `fac09f2`, dependency
baseline `00b105a` in `/tmp/sdp-sdui-widgets`. The dependency snapshot preserves
KB-SDUI-005 preview work with provenance; it is not WCI0 implementation.
[Candidate inventory](candidate-WCI0.json) and
[dependency inventory](kb005-dependency-manifest.json) identify tested files.
The coordinator integrated exactly those ten implementation/test files into the
shared working tree after comparing each target against the dependency baseline.

## Behavior and evidence

The pure preparation package normalizes internally, checks exact capabilities
including hidden descendants, validates provider/layout and creates an unmounted
session. Prototype mode reports unbound declarations without invoking a binder;
connected mode requires the real adapter and installed handlers. Failures close
the detached session. Caller-owned admission rejects source/model-stale and closed
candidates. Concrete Fyne facts distinguish interactive button/input from symbolic
SVG placeholders. Native admission rejects unsupported normalized models before
constructing controls, including on later model checks.

The [worker report](WCI0-worker.md) records commands, source hashes and limits.
Full SDUI race tests passed. SDL bridge/runtime tests passed, including the real
SDL fixture: missing module/signature/result fails with zero domain calls and
unchanged live handlers; subsequent explicit dispatch proves retained handlers.
Headless Fyne tests are adapter evidence, not OS interaction evidence.

Independent reviewer `01a11854-523d-7c11-adf4-f622b19b68c6` inspected the complete
implementation against `00b105a`, verified all ten hashes, ran targeted SDUI race
and SDL bridge/runtime tests, and reported scoped approval with no blocking
findings. Coordinator inspected the diff and preserved all pre-existing changes.

## Remaining work and discovered prerequisite

Atomic native bundle publication remains WCI1, together with profile 0.3,
collections, scrolling, native keyboard/pointer and connected collection evidence.
WCI0 does not close any widget-family requirement.

The additional SDL code-generation check fails in `TestGeneratedConstructors` at
line 33 on unchanged baseline inputs. Worker and reviewer both observed it; the
complete SDL suite is therefore not claimed passing. Root cause and a bounded
prerequisite repair are being investigated under this card before WCI1 generation
acceptance. Frozen fixtures and historical evidence remain unchanged.

## Subsequent prerequisite correction

[WCI1-P0](WCI1-prerequisite.md) diagnosed the mismatch as an active generated
constructor predating KB005 provenance fields. Regeneration through sdl-gen and
an executable reuse/provenance regression resolve it; SDL codegen/bridge and
SDUI codegen checks pass. This later repair does not rewrite the earlier observed
failure or imply a full SDL suite run. Phase commit: `d1c5c88`.
