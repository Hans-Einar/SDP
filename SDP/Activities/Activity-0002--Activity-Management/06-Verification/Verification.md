# Activity0002 qualification

| Field | Value |
| --- | --- |
| id | SDPTOOL-VER-ACTIVITY-0002 |

Producer candidate: eb11ac77b096587c9d4cbc8e7c85b634131098a3 in SDP.
Consumer candidate: a558c31 in the isolated XFMD worktree.
Both use branch sdp/activity-management. Technical qualification and independent review passed;
owner acceptance, Activity completion, merge, installation and release are separate.

## Evidence scope and exact candidates

A6/source-candidate.json pins the final producer inputs. The full SDPTool module
passed on cc7eebcc808bca37e87e9219bd5555d53ef2158d: all packages, including the full
governance suite (499.562s). The subsequent f8faf9b and eb11ac7 changes are confined
to Activity projection/reference identity, the active execution-contract guard,
root integration tests and accompanying documentation. Unchanged package evidence
is reused; final root/Activity race tests exercise those final deltas. Do not label
the earlier broad run an unchanged-final-commit run.

| Check | Result / evidence |
| --- | --- |
| Full module, go test -count=1 -timeout 30m ./... | PASS on cc7eebc; A6/full-module.log |
| Final root/Activity correction suite | PASS after f8 projection correction; A6/root-and-projection.log |
| Manual pilot copy adoption and external references | PASS; A6/manual-adoption.log; real source remains unchanged |
| Shared domain/history/installer/application race checks | PASS; A6/domain-race.log |
| Final implementation admission/controller race checks | PASS; A6/admission-race.log |
| Final root/Activity race suite | PASS on final candidate: root216.193s, activities17.923s; A6/final-delta-race.log |
| Seven native protocol/GUI regressions | PASS on a558c31; A6/native-gates.log |
| Exact final producer with native disposable and manual-project browsing | PASS; A6/native-fixture.log and native-real.log; actual screenshots/Markdown/SVG retained |
| Go/Python generated-event parity, random UUID and atomic switch TARGET predecessors | PASS in full install suite and race suite |
| Python management and Activity schema/replay tests | PASS: 6 + 5; A6/management-tests.log and activity-python.log |
| Entire real mixed management history | PASS 761 events before qualification closeout; A6/management-validation.log |
| Canonical ActivityCommands SDL | PASS; A6/model-check.log |
| go vet ./... | PASS; A6/vet.log |
| Independent review | APPROVED on final candidate; A6/Review.md |

A6/initial-full-module.log preserves the unsuccessful integration run: existing
compiled catalog expectations did not include the new actions, and an immutable
project test counted the required empty canonical flock inode as a document change.
Both expectations were reconciled to the selected behavior without removing source
byte checks. That run is not passing qualification evidence.

Environment: Linux/amd64, Go 1.26.10; XFMD uses GCC 14.3.1 / FOX 1.6.57, isolated Xvfb
and private test homes. Native consumer source hashes and original worker evidence
are retained in XFMD SDP/05--Implementation/SDP-Navigation/evidence/Activity-A3-001.
Native tests use the prebuilt candidate producer, not a substitute renderer.
Other writer/native platforms are not qualified here.

## Selected case coverage

| Obligation | Observable evidence |
| --- | --- |
| Creation/standard layout/mandate/root Sessions | Actual CLI/action integration; six phases, generic README and no execution dispatch |
| Multiple phase, plan/milestone, predecessor/handoff Session context | Typed record validation and retained creation snapshot tests |
| Real manual0.1 preservation | Real read-only catalog; isolated manual-history upgrade/adoption copies preserve descriptor bytes and ledger prefix |
| Exact retry/changed request/stale revision | Create/session/update/transition receipt and CAS tests; changed operations refuse |
| Partial creation/descriptor replacement/interrupted switch | Fault injection before/after each retained effect; recovery finishes fixed bytes or refuses external conflict |
| Concurrent activation/single-active | Concurrent service attempts yield one active winner; complete switch batch replay |
| Hold while running/unknown stores/claimed launch | Pending gate blocks new admissions; existing sending claim remains in-flight; stop/observe/recovery remain available |
| Scope/membership changes beneath live work | Adopted active contract changes require observed quiescent hold; metadata/reordering and pre-adoption mapping remain supported |
| Completion/no-code outcome/contribution/shared card | Explicit outcome/evidence and contribution dispositions; whole-scope checks do not complete shared cards implicitly |
| Dependencies/cancel/reopen | Stable identities, cycles, explicit satisfaction and active-dependent guards |
| Invalid/duplicate/missing/manual/old records | Shared live/captured reducer, strict version-presence checks, schema corruption and legacy compatibility tests |
| External/unavailable references | Visible diagnostics; distinct projects and slash-delimited identity tuples do not collide or open local paths |
| Legacy Sessions and readable group/phase pages | Native Activities tree plus linked and unassigned legacy records; one Sessions home per Activity |
| Large navigation | Producer pages of 100 children; native fixture visits 5000 documents by selective expansion/collapse within visible limits |
| Lifecycle/Git/ledger timeline | Causal intervals, UTC, clock reversal/unknown time, explicit linked-event join; unrelated same-time event excluded; duplicate commits correlated |
| Marker filters/coarse zoom/source preservation | SVG/Markdown controls and aggregate markers; complete textual events and canonical source links retained |
| Changed reviewed input/Blueprint ownership | Typed reference hash refusal and existing PG model/Blueprint freshness gates; no parallel model authority |
| Read-only investigation/unrelated KB capture | Context/Route/Read and read-only Reviewer/Leader paths preserved; only implementation admission is gated |
| Installation/older clients | Temporary installer/upgrade tests distribute schemas/reader/skills without opting into admission; populated old-client trees fail unsupported protocol |

## Limits and disposition

New producer/consumer protocol requires paired capable clients. No universal
upgrade of already-installed old clients is claimed. Native visible-row/depth and
producer resource caps remain explicit; oversized expansions diagnose limits,
and phase pages retain source links. This is bounded lazy browsing, not unlimited
virtualization. The pre-existing XFMD SDP-basename selection ambiguity remains:
select /home/warloc/git/SDP/SDP as the process area for this repository.

Admission is explicit opt-in against inventoried initialized PG stores. Read-only
review/leader work is intentionally outside implementation quiescence. This does
not stop arbitrary manual agents or hostile same-user writers. No live admission
policy or installation was enabled in the real project; adoption was exercised in
a disposable copy of its real manual Activity history and in generated fixtures.

The preliminary Blueprint remains constraint-checked rather than automatically
code-conformant. ../04-Design/CodeMapping.md maps its responsibilities to code/tests;
the immutable original and revised handoffs remain retained. No accepted model
release is inferred. Real Activity0003 is not created. Activity0002 remains active
until the owner's actual need-level disposition.

## Subsequent real adoption — Activity3

The preceding limits describe the Activity2 qualification checkpoint. The owner
subsequently authorized autonomous Activity3 completion/publication without time
for personal review. SDPTool created ACTIVITY-SDP-0003 and SESSION-SDP-0019, the
canonical PG store was explicitly inventoried/adopted, and a tool-recorded atomic
switch held Activity2 and activated Activity3. See Activity3 Session0019/T001 and
canonical management history. This later real adoption supersedes the checkpoint's
statements that no real Activity3/admission existed; it does not rewrite earlier
evidence or claim that a ledger transition stopped arbitrary agents. Activity2's
personal owner review was not performed.

## Delivered through the authorized XFMD0.10 release

Activity3 exercised real tool creation, Session creation, adoption and guarded
hold/update/resume. The final generic README correction c3463b7 closes the scaffold
requirement missed at the first qualification checkpoint. Full Activity tests and
independent delta review pass. The exact private producer now ships in signed
XFMD0.10.0-76.el10_2 with actual public upgrade/fresh-install native Activity tests.
See [release evidence](../../Activity-0003--Readable-SDL-browsing/06-Verification/ReleaseChecklist.md).

The owner authorized continuation and public delivery without personal review.
This objectively verified delivery closes the selected Activity2 implementation
scope under that authority; it does not assert personal owner acceptance, main
merge, accepted SDL model release or a public Toolkit installer release. Earlier
active/gate-review/on-hold statements describe prior checkpoints. No new Activity2
implementation is selected by administrative closure.
