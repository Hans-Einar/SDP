# SDUI staged widget implementation

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0022 |
| project | SDP |
| state | active |
| PlanType | ImplementationPlan |
| BranchPolicy | phase |
| CommitPolicy | milestone |
| Systems | SDUI; SDL; SDPTOOL |
| source | KB-SDUI-003; PLAN-SDP-0021; Session0010 T001 |

## Outcome and authority

The owner requests the work in KB-SDUI-003. Deliver the full listed widget inventory
in runnable stages, preserving the limited UI composition purpose, I/O-free frontend
and application-owned transactions. The prerequisite [design](../../../04--Design/SDUI/Widgets/Design.md)
and [matrix](../../../04--Design/SDUI/Widgets/Acceptance.md) govern behavior and proof.
No widget family is silently deferred. Release/publication and merge are unselected.

## Git and candidate policy

Use stacked phase branches `sdui/widgets-wci0` through `sdui/widgets-wci4` in an
isolated checkout, with milestone commits. Do not change the shared worktree's
branch or import concurrent governance changes. Baseline is repository HEAD plus
an explicitly inventoried snapshot of the SDUI preview/normalization producer work
needed from KB-SDUI-005. Commit that dependency separately with provenance before
claiming a clean candidate; an uncommitted candidate needs hashes/diff inventory.
Current preview files are not this assignment's original implementation.

## Phases and milestones

| Phase / milestone | Runnable outcome and acceptance | Dependency | State |
| --- | --- | --- | --- |
| WCI0-M1 | Capability/preparation boundary with legacy fixtures, missing capability/module/signature rejection and no partial activation | WCD1 reviewed design | completed |
| WCI1-M1 | 0.3 tree/list/scroll from source through runtime, Fyne, exports and typed SDL fixture; inspect pointer/keyboard/reload behavior | WCI0; reviewed [collection contract](../../../04--Design/SDUI/Widgets/Collections.md) | completed |
| WCI2-M1 | Tabs/splits preserve page state, relative geometry and focus | WCI1 evidence and pilot findings; reviewed Panes-and-commands contract | completed |
| WCI2-M2 | Button/toggle shared commands, menu/context/dialog with cancel/lifetime proof | WCI2-M1; typed command/surface grammar | completed |
| WCI3-M1 | Checkbox/slider/select/numeric input and atomic typed drafts | WCI2; reviewed Values-and-text contract and API reconciliation | completed |
| WCI3-M2 | Extended single-line and basic multiline text with native editing/IME/undo evidence | WCI3-M1; reviewed explicit opt-in/native editing/IME contract | completed |
| WCI4-M1 | All-family producer discovery/composition/text/codegen/provider integration and matching consumer package preparation | WCI1–WCI3; reviewed Providers-and-packaging contract and four API handoffs | in-progress |
| WCI4-M2 | Exact-candidate compatibility tests, native workflow evidence and independent integrated review; explicit disposition of all gaps | WCI4-M1 | planned |

Each stage updates the matrix and relevant requirements/language/architecture/runtime
documents. Do not merely extend parser tables. Static exports must reject unsupported
behavior or use an explicitly requested labelled fallback. Broad optional research
is excluded as stated by the primary card; full editor/FOX parity is not selected.

## Verification

Before code for each milestone, review its concrete stage contract: capability
identifiers and exact version matching, applicable EBNF/AST, typed property/event
payloads, bridge mapping, publication owner and rejection/recovery cases. This
applies uniformly, including WCI0 and tabs/splits. Later grammar may remain
pending at WCD1 closeout, but cannot be invented during implementation without
updating and reviewing the stage contract first. Connected numeric transport is
the checked integer action-core boundary in Design.md; fractional SDL transport
is explicitly rejected. WCI0 is scoped by [Preparation](../../../04--Design/SDUI/Widgets/Preparation.md).

Run targeted meaningful tests during development and `go test -race ./...` in
SDUI/go at integration; test SDL bridge/runtime/codegen and SDPTool consumers when
touched. Add positive/negative 0.3 baselines; retain 0.2 and frozen XM-M2 evidence.
Use controlled native fixtures on a separate display, inspect previews and exercise
keyboard/pointer, cancellation, stale completions, reload and disposal. Packaging
evidence names matching revisions of SDPTool and privately bundled Fyne/preview
helpers. Actual publication is separate from package preparation.

Independent review must challenge the integrated user workflow and all-family
completion matrix. An implemented slice does not complete the card. Outstanding
native, packaging or review evidence stays explicit and keeps completion pending.

## WCI0 delivery

[WCI0 evidence](Evidence-WCI0.md) records verified detached preparation/admission
and independent review. Atomic native publication remains WCI1. The SDL generated
constructor baseline mismatch is a tracked prerequisite investigation before WCI1;
no all-SDL-suite pass is claimed.

WCI1-P0 prerequisite delivered: regenerate active model provenance with the existing
generator, retain frozen fixtures, compile/run reused-instance constructor checks.
See [report](WCI1-prerequisite.md); stage feature implementation remains pending.

## WCI1 delivery

[WCI1 evidence](Evidence-WCI1.md) records phase commit `c39b330`, all 88 reviewed
files, 54 passing native checks and passing integrated original-workspace suites.
Tree/list/scroll and guarded native publication are delivered at the bounded
contract level. WCI2–WCI4 remain required; no whole-card completion is implied.

## WCI2-M1 delivery

[Pane evidence](Evidence-WCI2-M1.md) records commit `403c540`, 72 exact integrated
files, 58 native checks, passing suites and independent approval. M2 is selected after independent contract review (hash 9a9c8710): preserve basic
button dispatch, explicit opener identity and synchronous native menu selection
scope. Five exclusive implementation lanes continue on the same phase branch.

## WCI2-M2 delivery and WCI3 selection

[Commands/surface evidence](Evidence-WCI2-M2.md) records `0fc15c8`, 113 files,
118 native checks, 26 exactly-once opening results and independent approval.
WCI3-M1 is selected on `sdui/widgets-wci3` under existing owner full-card authority.
The reviewed Values-and-text contract and concrete API handoffs preserve WCI2
closed-dialog successor semantics and bounded Go numeric admission. Five exclusive
implementation lanes cover frontend, runtime/numeric, layout, host and SDL bridge.
Extended text/IME remains WCI3-M2. KB004/KB005 retain separate existing dispositions.

## WCI3-M1 delivery and extended text selection

[Scalar evidence](Evidence-WCI3-M1.md) records `ea49991f`, 88 paths,
113 native checks, 12 exactly-once results and independent approval. WCI3-M2 is
selected on the same phase branch under existing full-card authority. Reviewed
explicit source opt-in, single text authority, exact Commit/self-echo, CR/LF and
required-empty successor semantics govern implementation. Five disjoint lanes
extend frontend, runtime, layout, host and SDL fixtures. Main owns the bounded
pinned X11 GLFW filter dependency, native IME evidence and canonical integration.
No further widget family is deferred; WCI4 remains required after text acceptance.

## WCI3-M2 delivery and final integration selection

[Text evidence](Evidence-WCI3-M2.md) records `69d0a332`, 224 paths, 92 native
checks, 5 exact results and independent approval. WCI4-M1 is selected under
existing full-card authority on `sdui/widgets-wci4`. The reviewed provider contract,
backend matrix, public Accessible boundary and final frontend/provider/layout/host
API handoffs govern implementation. Direct supplied SVG and bounded Markdown prose
need actual native proof; unsupported diagrams resolve to declared label/reject
before outcome freeze. Main owns matching package preparation and consumer evidence.
WCI4-M2 retains integrated all-family verification and independent closeout.
