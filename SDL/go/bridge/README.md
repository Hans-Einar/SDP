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
