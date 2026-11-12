# Session 0018 — Activity management

| State | Step | Outcome / authority |
| --- | --- | --- |
| completed | S1 Manual registration | Owner accepts predecessor and selects Activity0002 |
| completed | S2 Design and implementation | A1–A5 candidate implemented; integrated qualification follows |
| completed | S3 Qualification | Candidate verification and independent review passed |
| waiting | S4 Owner disposition | Concrete qualified result ready for owner review |

| Field | Value |
| --- | --- |
| Session reference | SESSION-SDP-0018 |
| Activity | ACTIVITY-SDP-0002 |
| Status | active |
| Current step | S4 |
| Primary card | KB-SDP-061 |
| Plan | PLAN-SDP-0027 |
| Phase contributions | study, requirements, design, implementation, verification |

## T001 — Owner-selected execution

Owner accepts Activity0001 and requests starting Activity0002. Use SDP Planning,
Architect, Master/Worker and Traceability; reuse the existing study/design and
reconcile implementation details before code changes. Register this Activity and
Session manually; the following Activity is the first intended real tool creation.

Preflight: SDP has one clean worktree at430207b before predecessor closure; closure
commit26d4879 records acceptance. New branch sdp/activity-management continues the
accepted source without merging main. Activity0002, Session0018 and PLAN0027 are
unallocated in known records. XFMD main8893a32 has unrelated dirty RPM/process work;
leave it untouched and use an isolated worktree for consumer changes. Remote
concurrent work is outside observation. Green for SDP records and A0 design scope;
recheck product and consumer boundaries before implementation.

A bounded read-only architecture reconnaissance is assigned separately to locate
shared transaction/admission contracts. No automatic Activity admission or capture
is claimed by these manual records. Next: A0 executable contract/model handoff,
then complete A1–A6 under the authorized plan.

## T002 — Executable design and bounded implementation

A0 completed: validated/checkpointed ActivityCommands :00002 and the bounded
BeforeActivityCommands :00001 baseline, generated the selected Blueprint handoff,
and pinned transaction, descriptor retention, explicit admission and projection
contracts. Generated views remain on demand; Blueprint generation is a selected
worker handoff, not authored viewpoint output or runtime proof.

Implementation is split across records/lifecycle, shared history/admission and
producer projection, with the main agent owning public adapters and integration.
Projecthistory transaction batch tests, race checks and vet passed. The shared
transaction now supports held-lock reads and atomic ordered publication without
reentrant locks; exact retries verify bytes and the original prefix.

A1–A5 remain in progress. Outstanding integration includes CLI/action end-to-end
checks, large-tree expansion, native consumer, schema/installed reader adoption,
and integrated independent qualification. No real Activity0003 is allocated.

## T003 — Integrated implementation and review corrections

CLI/actions now share Activity operation metadata and typed calls. Activity and
root Session creation, lifecycle/recovery, revisions/dependencies/dispositions,
producer pages/paging/calendar and explicit PG admission are implemented. Native
XFMD consumes the producer protocol in an isolated worktree; actual disposable and
real manual read-only browsing passes. The real third Activity remains unallocated.

Independent review identified generated UUID IDs rejected by the Python reader/
0.2 predecessor schema, and forbidden0.2 extensions accepted by Go's0.1 parser.
Both are corrected with real generated-history parity and raw-shape regressions.
The root bind/adopt journey also exposed and corrected confusion between the PG
browser project identity and canonical Activity namespace. No approval is inferred
from these fixes; final qualification remains in progress.

The old compiled action-catalog test is extended for the intentionally new actions.
The immutable-project test still checks every source byte; it excludes only the
empty stable canonical flock inode required to serialize implementation admission
with opt-in adoption. Removing that inode after unlock would make flock unsafe.

The retained Blueprint is regenerated with Template and SDP/ProjectManagement
code paths, reflecting the actual source homes of installed schemas/readers.
Original immutable artifacts/checkpoints remain. Source process skills now explain
adopted scope, unrelated capture and explicit admission without migrating legacy
projects. SDP Verifier is used for exact-candidate evidence and limitations.

## T004 — Qualified candidate and handoff

Producer eb11ac7 and isolated consumer a558c31 are technically qualified and
independently approved. Full module regression passes on cc7eebc; final Activity/
root deltas pass complete race checks. Native fixtures and the actual manual
project open the Activities tree, pages, source links and colored calendar SVG.
Installer/Python/Go reader parity and real manual-history copy adoption pass.

Final review additionally closed ambiguous external reference keys and active
execution-scope updates. An adopted active Activity must reach observed quiescence
through hold before scope/reference/dependency changes. Ledger state never pretends
to stop already-issued work. All findings and limits are retained in A6/Review.md.

KB061 is gate-review; PLAN0027 and Activity0002 remain active pending actual owner
need-level disposition. The real third Activity is still unallocated. No live policy
adoption, installation, merge or release occurred. Existing primary XFMD RPM/process
edits remain untouched. Final response will summarize this recorded work; this is
not a claim of an already-captured response.

## T005 — Review branch handoff

Under the established phase-push authorization and the owner's request to retrieve
review branches locally, sdp/activity-management is pushed to origin in both SDP
and XFMD. Initial published tips are SDP de3c53d and XFMD4c77944; qualified product
candidates remain eb11ac7 and a558c31 respectively. The only subsequent SDP change
is this handoff record. No main branch, installed binary or release was changed.
The owner can fetch the same feature-branch name in both repositories. Activity0002
remains active/gate-review; next is the owner's actual disposition.
