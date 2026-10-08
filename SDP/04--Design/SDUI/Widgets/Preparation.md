# WCI0 preparation contract

This stage preserves source/AST 0.2. No new widget grammar is implemented here.
The selected seam is a pure SDUI preparation package, used by prototype preflight
and native-host admission, with an SDL bridge fixture using its existing Bind API.

## Capability matching

Capabilities are exact identifiers with integer major versions; missing or
different versions reject. No optimistic range matching. Required dimensions:
frontend `sdui/0.2`, layout `relative/1`, widget `button/1`, `input/1`,
`svg-placeholder/1`, viewport `scroll-x/1` and `scroll-y/1` when requested.
Provider content and native interaction remain separately declared dimensions;
known 0.2 SVG content is a labelled placeholder, not real vector execution.
The actual baseline declares only existing support and rejects scroll. Unknown
normalized node/widget kinds must reject rather than rely on fallthrough drawing.

Check every node, including hidden descendants. A diagnostic identifies dimension,
missing capability, instance path and original source span. Host/provider facts
are passed by the concrete adapter; successful parsing is not evidence of them.
Shared package imports no Fyne or SDL and performs no filesystem I/O.

## Candidate lifecycle and binding

Preparation receives a validated document, entry, session identity, source revision,
concrete capability set, layout validator and explicit mode (prototype/connected).
It must normalize internally or otherwise validate the document/root correspondence;
callers cannot pair unrelated documents and roots. All preparation is synchronous
on the owner goroutine. Build a detached runtime Session and inspect required
capabilities before invoking binding preparation.

Connected mode requires an actual adapter function that binds the detached session
using the already loaded SDL modules and typed plans; it returns an error, not a
boolean readiness assertion. Validate every declared callback and connection in
the selected entry. The adapter must not execute domain actions. Reject a missing
adapter in connected mode, even when the document has no callbacks, rather than
infer connected readiness. The real bridge fixture proves missing module and wrong
signature leave both the existing Session and its handlers unchanged.

Prototype mode records unbound symbolic callbacks/connections and never invokes
the binding adapter. Candidate ownership is explicit: failure closes/disposes the
candidate; successful preparation transfers it to the composition caller, which
must close it if abandoned. No prepared candidate is mounted automatically.
Reject closed/stale candidates before publishing. A final source revision recheck
belongs to the I/O-owning caller; a source hash cannot be checked by the pure core.

WCI0 does not replace the live application reload implementation or claim a generic
atomic host swap has been delivered. Its no-partial-activation proof is that
failure occurs on a detached, unmounted session before any published state changes.
WCI1 must implement the full publication contract in Design.md before mounting new
widgets. The existing host admission uses capability checking to reject unknown
normalized models before constructing native controls; it does not claim connected
activation merely because legacy NewRuntime accepts an already supplied session.

## Required evidence

Tests cover exact capability mismatch, unknown widget/node, hidden unsupported
content, preserved diagnostic paths/spans, failed geometry, invalid entry/document,
prototype adapter non-invocation, missing connected adapter, disposal on adapter
failure, and real existing SDL Bind failures/success with zero actions during
preparation. Prototype check and native admission retain current 0.2 behavior.
No native interaction or new-widget delivery is claimed by these preparation tests.
