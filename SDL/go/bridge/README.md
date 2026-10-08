# Explicit SDL callback bridge

Bind validates declared aliases, actual SDL action/Go-registration signatures,
selector ownership and results before installing handlers. It never opens ref
paths or invokes domain code during preflight. Existing basic buttons and legacy
Event/Widget conversions retain their Handler route; canonical commands, tabs
and dialog Accept use InteractionHandler. No event travels through both routes.

| EventField | Owner/event | SDL destination |
| --- | --- | --- |
| CollectionItemID | Collection Activate | text |
| TabPageID / TabPreviousPageID | Tabs ActivatePage | text |
| CommandContextItemID | Item-context InvokeCommand | text |
| CommandChecked | Toggle InvokeCommand; proposed state | boolean |
| DialogFieldValue + FieldPath | Dialog Accept; captured owned input | text |

Exactly one source selector is required. Only DialogFieldValue permits FieldPath;
its nonempty named relative path is resolved by runtime/frontend in preflight.
Nested-dialog/foreign/anonymous paths reject. Per-owner resolved handles are kept
in copied plans, so reused dialogs capture independently. Accept reads the event's
captured String value, never later live input state. WCI3 typed/select mappings
are not implemented here.

Plan.ResultMode is closed:

- TextResult is zero/default: text OutputField and explicit setHandle to input;
  existing optional revision result/context rules remain.
- DialogAcceptResult requires boolean AcceptField and text MessageField. It
  forbids OutputField, setHandle, RevisionField and RevisionContext. Only dialog
  Accept owns this result. It maps the validated returned record to Reply.Accept
  and returns no extra widget updates. Capture count is preflight-bounded by 256;
  runtime owns the combined atomic batch and revision checks.

False returns rejected without domain commit, with a default message if empty.
True returns succeeded before post-execution UI checks. Execute/record-validation
errors retain unknown. Successful domain work followed by stale draft/resource
publication retains succeeded; the runtime blocks replay for that opening and
permits Cancel/Close while preserving the outcome. The bridge adds no retry or
reconciliation API. Runtime owns dispatch eligibility, captured revisions, atomic
publication, local command effects and exactly-once terminal dialog receipts.

The commands fixture exercises real registrations, signature/FieldPath negatives,
false/true/malformed/error outputs and post-domain conflicts. Native focus/menu/
window behavior requires the coordinator's separate actual-input evidence.

## WCI3-M1 scalar Commit

ScalarResult is the explicit mode for checkbox/slider/number/select Commit. A
single self setHandle receiver is required because its one checked update must
accept the captured source proposal. OutputField matches boolean/integer/text-ID;
AcceptField/MessageField are forbidden. Existing optional result revision/context
mapping remains. TextResult is never inferred from a scalar widget or repurposed.

ControlBoolean reads checkbox Boolean, ControlNumber reads slider/number Number
into SDL integer, and ChoiceOptionID reads a captured select ID into text. The
closed ControlText selector supports the existing basic input String Commit with
TextResult. Explicit extended input uses the WCI3-M2 typed capture below. Legacy
Event/Widget conversions are unchanged. No scalar selector stringifies/coerces arbitrary input.

Numeric preflight uses frontend NumericArguments and the shared numeric package:
original min/max/step/value lexemes must be exact safe53 integers with round trips.
Captured numeric raw draft is independently checked before conversion, including
exact grid membership; typed output int64 is bounded before float64 conversion.
Fractional Go-only grids reject SDL connectivity with path/span context. No private
numeric parser, membership tolerance or rounded decimal-provenance reconstruction.

Scalar handlers recheck capture/receiver/model/value/draft/option generation and
post-consumption state after Execute, then return one ExpectedValueRevision /
ExpectedDraftRevision / ExpectedOptionGeneration accepting update. Stale returned
work is not replayed; the legacy Handler protocol makes no domain-outcome promise.
Select uses dedicated OptionTarget validation, never a collection/index/label lookup.

DialogAcceptResult remains text-only, including hidden owned fields. A mixed form
requires explicit Go InteractionHandler; no automatic heterogeneous SDL Accept
record mapping is introduced. Existing text-only M2 dialogs remain supported.

## WCI3-M2 extended text Commit

Only source-explicit multiline/readOnly/placeholder/required presence opts into
extended input, including false or empty arguments. The bridge consumes the
runtime FieldState.Input projection; it does not infer policy from profile,
selector, handler, or the presence of a Control envelope.

TextResult on an extended callback owner requires a self setHandle at complete
preflight, regardless of the selected request source (ControlText, Event, Widget,
Literal, or Context). ControlText consumes exact String plus Control.ValueRevision;
RawDraft and Option must be absent from that event. Basic input retains its legacy
String event with no Control envelope and may still target another input.

After Execute, extended Commit rechecks the full source FieldTarget and exact
text bytes returned by SDL. Non-echo, stale, malformed and failed results never
accept the proposal or replay domain work. Its single update checks both value
and draft revisions and explicitly accepts the captured draft. Unicode and CR/LF
are never normalized by the bridge. Runtime owns validation and line policy.

Legacy Load and other non-input owners can still write extended readOnly inputs
using checked value/draft revisions. Programmatic updates remain muted. Existing
legacy Handler results do not claim domain-outcome reporting or rollback. The text
fixture exposes independent domain counters so a successful domain write followed
by a UI conflict remains observable.

DialogAcceptResult continues to capture only text inputs, now including extended
single/multiline fields; mixed scalar forms still require Go Accept. Native Entry
history, clipboard, IME, wheel routing and focus belong to the host and require
separate actual OS evidence. See `../examples/text/README.md` for the real SDL
fixture and condition-only protocol.
