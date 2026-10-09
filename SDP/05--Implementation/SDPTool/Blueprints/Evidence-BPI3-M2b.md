# BPI3-M2b — Assignment lifecycle and grouped discovery

2026-10-09. PLAN-SDP-0020; KB-SDP-050; Session0008 T012. Candidate is baseline
295b048 plus the source bytes listed in Pilot-BPI3-M2b.json. Source tests precede
the milestone commit; no exact-commit or published-release result is invented.
Toolchain: Go 1.27.1 linux/amd64, Python 3 with jsonschema.

## Delivered workflow

The local CLI applies explicit canonical assignment events through the shared
append primitive. The reducer rebuilds separate attempts from the project ledger,
with predecessor checks and trusted adapter attribution. Readiness adoption binds
current model/evidence captures; submission binds scoped TARGET code, actual check
receipts and a Traceability event; independent reviewer acceptance is required
before controller closure. Commands do not execute evidence text.

The generated non-Git trial has two retained revisions of one blueprint task and
four separate attempts. One completes after rejection/rework, one is superseded,
one canceled and the new revision stays draft. Compiled discover and assignment
list agree after copying the fixture into another project directory. See frozen
[trial outputs](Trial-BPI3-M2b/README.md). Assignment targets carry verified file
hashes. No mutable bundle status or manual navigation registry is added.

## Verification coverage

- Compiled CLI: create, adopt-readiness, assign, start, hold/resume, submit, reject,
  resubmit, independent accept-review and controller complete. Terminal reopening
  refused; missing acceptance, wrong assignee, self-review and stale predecessor
  refused without changing history. Multiple attempts and a second revision remain
  separate. Whole retained bundle digest is unchanged across transitions/refresh.
- Evidence: stale code, failed source binding and changed Traceability bytes block
  adoption/submission/closure. Actual synthetic Go behavior tests produce a passing
  receipt; this proves only that fixture, not general code conformance.
- Recovery/concurrency: forced subprocess exit after append before response retries
  without duplication; simultaneous assignment cancellations have one winner.
  M2a additionally covers process exit before/after rename and stable lock recovery.
- Discovery: completed state survives a corrupted/missing bundle; no invalid open
  target is returned. Adopted/submitted evidence remains historically present but
  unavailable. Known mismatch is stale; missing input/budget exhaustion is unknown.
  Shared read-budget and repeated/third-binding cache-limit controls pass.
- Canonical compatibility: Go/Python accept the same generated chain and reject
  missing/null/case-aliased fields, invalid predecessors, self-review, cross-System
  supersession and time reversal. Nanosecond ordering is checked independently.
  Installer validation preserves valid assignment history and refuses invalid roles.
- SDL overview is canonical and passes sdl check. Process/Toolkit validators pass.

## Review and test record

Independent review found historical projection loss on unavailable bundles,
cleanup incorrectly tied to active parent plans, missing open-target hashes,
reset read budgets and replay consistency gaps. All were corrected with regression
coverage. Subsequent review requested explicit stale/unknown distinction and
evidence diagnostics; these were corrected too. Review scope is source lifecycle
implementation, not owner pilot acceptance, remote authentication or native UI.

Focused normal tests passed, including compiled CLI, source/evidence controls,
installer and cross-language replay. Final race/vet outcomes are recorded below
and in the machine record. A broad `go test -race ./... -timeout 180s` attempt
exceeded its per-package limit in root, blueprints, installer and unchanged model
packages under concurrent host load; the remaining broad run was terminated after
those failures to avoid duplicate load. Its log is retained. This was not a pass.

## Limits

Writes currently require Linux. The CLI is a trusted local attribution boundary,
not authentication; MCP adapters must bind caller capabilities separately. Evidence
is scoped and attributed; the tool does not establish full model/code conformance
or execute receipts. Local arbitrary writers are outside the shared locking
protocol; source/evidence rechecks are not a cross-file transaction.

Catalogue and assignment scans have separate budgets; live model comparisons are
capped as specified in Assignment-Lifecycle.md. Unchecked freshness is explicit.
No native XFMD, release publication, installer manifest change or owner pilot
acceptance is claimed. Final plan closure depends on the remaining integration
verification, not merely this document's presence.


## Final results and disposition

Independent Reviewer /root/bpi3_m2b_review approves bounded M2b code delivery,
with no remaining material findings. Its full relevant race suite passed:
blueprintstate 37.789s, blueprints 546.027s, projecthistory 7.085s, install 494.870s
and root SDPTool 343.673s. The long runtime occurred under host contention; it
does not change the failed three-minute attempt into a pass. Independent vet passed.

The reviewer binaries preceded the final navigation unavailable/empty-tab handling
and stronger stale assertion; those were reviewed statically and the exact final
delta was additionally tested with race detection: blueprints 61.593s, root SDPTool
1.444s. This includes concurrent cancellation, closed-plan cleanup and missing-bundle
navigation. Final `go vet ./...` passed. Candidate hashes and the frozen generated
consumer outputs are in Pilot-BPI3-M2b.json.

The prior installer race-test evidence gap is resolved for this candidate. The
implementation plan is complete for its bounded delivery; card and Session remain
for the concrete BP2-A owner pilot review, not additional hidden implementation.
No broad release qualification or integration of parallel work is claimed.
