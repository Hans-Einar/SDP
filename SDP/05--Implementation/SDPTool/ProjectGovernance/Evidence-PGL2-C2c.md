# PGL2-C2c — Local controller design decisions

2026-10-18, Session0007 T022 / PLAN-SDP-0018 / KB-SDP-052. Owner authorizes continuing
while reading the walkthrough. This is authority to build the selected capability,
not acceptance of a test assignment. [Orchestration](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md)
owns the contract. [Candidate](candidate-PGL2-C2c.json) pins the tested source/binaries
and check outputs; [trial](trial-PGL2-C2c.json) retains actual subprocess evidence.

## Delivered behavior and limits

The trusted local controller records explicit approved, changes-requested or rejected
input against a current C2b review revision. The stored record retains the complete
checked package, actor/source/comment and recording time. Atomic generation writes
retain immutable operation receipts. Later decisions supersede the effective decision,
without erasing earlier records. Retry returns historical evidence, never renewed
acceptance. The core refuses workers, wrong bindings, changed operation content and
cross-operation reuse. No worker MCP decision privilege is added.

The standalone design-decide adapter shows the proposal and captures stdin separately
from the JSON package request. EOF, blank choice, cancel, invalid choice and missing
comment do not record anything. Exact CLI retry returns the saved receipt without
another prompt. A changed actor attribution cannot retrieve that receipt through
this adapter. Actor/source are controller-attributed observations, not authenticated
human identity. Same-user unrestricted shell access remains outside the existing
Caller boundary. This is not a completed conversation capture or XFMD integration.

Design-review reports current applicability independently of historical receipts.
Changed design, preparation, retirement or unavailable provider checks make a decision
stale. Exact source restoration may match again; unobserved transient filesystem
edits are not tracked. A new preparation revision always changes the review basis.
The technical-only design-check and preparation inspection retain their earlier
contracts; callers use design-review for owner disposition. Approval never launches
work, closes cards, promotes models or supplies a cross-worktree reservation.

Successful first decision upgrades only the selected store to 0.4. Earlier readers
refuse it, new readers retain old-profile support, and later preparation/lifecycle
operations preserve 0.4 and decision history. Existing local trials are not upgraded.

## Verification

Passing checks: affected governance/controller/MCP race tests and vet; the complete
SDPTool `go test ./...` suite; management and Toolkit validators; schema checks for
48 actual design request/result objects. These are same-context checks, not independent
review or an owner acceptance trial. Tests use real model/Blueprint provider fixtures.

Service tests cover persistence/reopen, current versus stale decision, supersession,
immutable retries after drift, changed preparation and retirement, subsequent prepare,
worker/no-write and attribution/input guards, cross-operation collisions, concurrent
decision exclusion, source drift before publication, and before/after-HEAD recovery.
The concurrent test accepts the existing nonblocking lock's busy refusal or the
generation guard's stale refusal; exactly one decision succeeds. No lock behavior
was weakened to satisfy a test.

The maintained decision_trial.py invokes built binaries, including a C2b old reader.
It supplies clearly labelled scripted fixture-owner decisions, verifies no writes on
cancel/EOF or inspection, changes TARGET while the real adapter waits for stdin,
and verifies rejection without a partial record. After successful writes it retries
without stdin and obtains the original receipt, even after a newer rejection.

The separate preview is `/home/warloc/git/pg_lab/previews/pgl2-c2c`. Read-only
`/home/warloc/git/pg_lab/design-review` shows a current design with the final rejected
fixture decision. An optional --decide --owner NAME interaction affects this fixture
only. All 176 original owner project/runtime file/symlink entries are unchanged.
The existing C2a/C2b previews and installed tools are untouched.

## Integration and remaining work

Scoped preflight starts from e25f13c, with main e8760d3 integrated. Changes are PG
core storage/services, standalone controller, lab drivers and owning records. Provider
services, generic Actions, host protocol and XFMD code remain unchanged. Shared event
IDs are rescanned before appending; earlier ledger bytes remain intact.

PGL2 remains in progress for dispatch prerequisites. Responsive leader orchestration,
parallel Masters, durable inbox, direct conversations and integrated owner trials
remain later milestones. The owner can give walkthrough feedback while work continues;
this delivery does not introduce a new mandatory approval gate for that continuation.
