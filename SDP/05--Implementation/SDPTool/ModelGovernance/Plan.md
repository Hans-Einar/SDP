# MGI — bounded ModelGovernance implementation

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0019 |
| project | SDP |
| state | active |
| PlanType | ImplementationPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDPTOOL, SDL, SDUI |
| source | PLAN-SDP-0016 / MG4-M1; KB-SDP-049 |

## Outcome and authority

Deliver a usable, bounded local SDL/SDUI model workflow: create WORK, save commits,
recover the whole model, integrate changes, freeze optional proposals/candidates,
and publish accepted model releases with retained metadata lineage. Work without
project Git and without a separate export operation. Session0006 continues across
this plan; completing its predecessor DesignPlan does not complete the feature.

The owner selected MG4 planning on 2026-10-05. This ImplementationPlan is prepared,
not started. At execution start record the selected milestone and an active plan
event; do not interpret this document as merge/publication authorization.

## Governing inputs and implementation boundary

- [Study](../../../04--Design/SDPTool/ModelGovernance/Study.md): owner decisions and rejected alternatives.
- [Contract](../../../04--Design/SDPTool/ModelGovernance/Contract.md): candidate v1 behavior.
- [SDL model](../../../SDL/ProjectGovernance/SDPTool/ModelGovernance.design): supported composed concern view.
- [MG3 proof](../../../04--Design/SDPTool/ModelGovernance/Proof.md): feasibility only, with explicit missing controls.
- [Session0006](../../../Sessions/session-%230006--Model_governance.md): journey and current step.

Build reusable Go code under SDPTool/model (or a documented equivalent package),
with a thin `sdptool model` adapter and existing human/--json presentation routing.
No runtime plugin loader or second parser. Use SDL/SDUI libraries for validation.
Keep library operations usable on explicit directories without SDP installation or
Git; the CLI must route model operations before mandatory project-discovery checks.
A separately packaged executable is not required for this first delivery.

Do not copy the experimental store into production as if it already fulfilled the
contract. Reuse demonstrated algorithms only after adding schema, containment,
transactions and tests. Keep experiments as frozen evidence, not a runtime dependency.
Do not use the experiment's external git merge-file call as a production dependency.
Semantic blueprint analysis/rendering belongs to KB050. Provide source snapshots
and identity metadata for that consumer; do not advertise create blueprint yet.

## Git, concurrent work and evidence policy

Use one work branch for MGI, not a branch per small phase. At implementation start,
create sdp/model-governance-implementation from the committed MG4 delivery (or record
an already isolated equivalent). Use an isolated worktree if ProjectGovernance is
still sharing files; do not switch another agent's working branch. Subsequent phases
continue on that branch; each delivered milestone gets its own ID-bearing commit.
Keep one combined PR against main when agreed work is ready. Follow existing phase
push authority; main merge and product release require explicit authorization.

Inspect status before changes/staging; preserve unrelated sourceinput/node files.
Coordinate shared cli.go/discovery/presentation edits with ProjectGovernance. Keep
new model operation code in its own package/adapter rather than rewriting the CLI.
No installer template migration, live XFMD changes or gh-sdp release in this plan.

At every milestone: update evidence, plan, Session, card and relevant design/model
before committing; use Traceability for system changes and ProjectManagement for
lifecycle. Card stays active/ready between work, in-progress during execution.
Independent final review uses context separate from implementation; no claim of
independent approval if unavailable. No simulated owner acceptance.

## Phases and milestones

All states below are planned; acceptance lists are obligations, not passing results.
Each phase ends with a runnable vertical workflow and evidence on its actual candidate.

| Phase | Milestone | User-visible delivery | Depends on |
| --- | --- | --- | --- |
| MGI1 Safe WORK | MGI1-M1 | `create work:NAME --initial`, `status`, baseline persistence, human and JSON output | MG4 handoff |
| MGI2 Recovery | MGI2-M1 | `commit`, whole-state `restore`, `history` and explicit interrupted-operation recovery | MGI1 |
| MGI3 Integration | MGI3-M1 | `merge SOURCE into WORK` and `create work:Combined from A B`, conflict inventory and resolution | MGI2 |
| MGI4 Frozen delivery | MGI4-M1 | `create candidate:NAME` / optional `proposal:NAME`, parser validation and metadata-only lineage | MGI3 |
| MGI4 Frozen delivery | MGI4-M2 | `create release:VERSION from candidate:NAME`, accepted-head handling and default WORK from release | MGI4-M1 |
| MGI5 Consumer/closeout | MGI5-M1 | Artifact-aware discovery/tree and consistent preliminary source views without registering sources manually | MGI4 |
| MGI5 Consumer/closeout | MGI5-M2 | Full CLI acceptance, independent review, documentation and remaining-risk disposition | MGI5-M1 |

### MGI1 — safe artifact creation

Define executable schema sdp-model/0.1 and positive/negative fixtures for WORK,
commit, frozen artifact, lineage and transaction metadata. Adopt exact UUID,
sequence, hash and path rules from the contract. Validate raw YAML restrictions
before decoding: aliases/tags, duplicate/unknown fields, document count and size.
Domain checks enforce graph acyclicity, parent availability and kind transitions;
JSON/YAML shape alone is insufficient.

Expose typed errors and one stable result envelope through human and --json output.
Choose/document exit statuses and operation recovery identifiers before adapter code;
errors never print success. Reserve metadata/history names, reject path escapes,
symlinks, nonregular files and portable case/name collisions. Implement local writer
ownership and same-filesystem staged create before publishing any mutating command.

Acceptance: clean temporary directory, no .git and no SDP manifest; creation and
status succeed. Empty .merge recreation after copying is harmless. Existing artifact
is never overwritten; ambiguous selectors, unexpected files in reserved areas,
malformed YAML and unsafe paths leave inputs intact. Compare golden source/metadata
hash vectors across process restarts. Test two processes attempting the same create.

### MGI2 — recovery before integration

After-image/deletion commits compare with previous head; rename begins as delete/add.
No-op commit creates no record. Invalid intermediate model can be saved for recovery.
Restore reconstructs the whole selected local state, saves dirty content first and
appends a new recovery record without reusing counters. History distinguishes a
restorable checkpoint from metadata-only ancestry.

Implement transaction journal and explicit recovery operation for interrupted writes;
finalize its human-readable command syntax in Contract.md before exposing it. Expected
head/source checks prevent stale writes. Cooperative locks do not lock external
editors: document consistency limits and preserve conflicting observed edits.

Acceptance: add/edit/delete/empty-file/rename, dirty restore, repeated restore and
source-tree replacement checks against exact bytes. Inject process termination at
payload stage, journal publication, each replacement boundary and head update;
reopen in a new process and prove resumable/abortable state with no lost dirty data.
Corrupt pending records/payloads must block publication. Lock contention and abandoned
writer recovery need process tests, not only race instrumentation. Retain recovery
backups until a documented cleanup operation; no silent expiry.

### MGI3 — integration and archive lineage

Record exact input IDs/digests and the merge base. Reconcile dirty source bytes without
modifying the source head, checkpoint dirty target and archive both pre-merge states.
Stage outside traversal, reject self/descendant copies and deduplicate only matching
immutable identities. Repeated integration must not create another merge event for
an already incorporated input, not merely leave contents unchanged.

Select a bounded Go three-way text implementation with recorded license/dependency
review and adversarial tests; no shell or installed Git requirement. Retain conservative
conflicts for binary edits, ambiguous renames, deletion/modification and path collisions.
Stop on missing/multiple common bases. An initial staged conflict inventory must
survive restart. Define explicit resolution-completion checks before promoting a
commit as a resolved merge; no merely deleting markers to bypass domain validation.

Acceptance: same-file disjoint edits, overlapping edits, merge with deletion, two
independent initial histories, different release bases and missing common ancestors.
Exercise three generations of .merge, duplicate ancestor identities, cycles/depth
limits and same-name different-UUID origins. Restore the complete pre/post merge
states. Selective inverse merges/cherry-picks remain excluded.

### MGI4 — immutable candidates and releases

Derive validation entrypoints from source headers/dependencies using existing
language libraries; persist the observed validation set and hashes for the frozen
artifact. Do not introduce a manually maintained navigation or file-registration
list. Explicit entry selection handles ambiguous source graphs; fragments are not
independent complete models. Update Contract.md's validationTargets meaning to this
source-derived receipt before implementation if it could be read as manual registry.

Candidates/proposals freeze consistent validated sources; dirty WORK need not be
committed solely for promotion, but exact captured bytes and provenance are retained.
Four-character UUID suffix collisions regenerate a new unpublished identity or fail
without overwriting. Promotion drops undo payloads and embeds the complete metadata
lineage closure. Dangling references to deleted WORKs fail acceptance.

Release accepts exact candidate identity/content plus an explicit evidence disposition:
model-only or verified with referenced code digest/checks. This is attributed acceptance,
not authenticated remote access control. Persist predecessor and expected accepted
head. No release receives edits. Bootstrap first release from an initial candidate;
then test WORK creation with no `from` selects that accepted baseline.

Acceptance: invalid SDL/SDUI, stale capture and unresolved merge reject candidate;
valid composed SDL and SDUI pass through real libraries. Delete originating WORKs
and browse retained lineage. Duplicate lineage ID with changed content fails. Verify
candidate sources unchanged at release. Model-only is never described as implemented.
Create competing releases in two project Git clones, integrate ordinary files and
verify both identities survive while duplicate version/competing heads block defaults.
Git is a transport test only, not a runtime requirement. Final source acceptance is
not inferred from highest version. No automatic patch-number allocation in v1.

### MGI5 — consumer boundary and closeout

Prune history/staging only inside recognized model artifacts, not arbitrary user
folders named .merge/.commits. Expose artifact roles without presenting historical
copies as duplicate live systems. Existing discover/tree/select/navigation outputs
and source-owned refresh semantics remain compatible outside artifacts. Model-only
folders must remain usable without installing the process template.

Preliminary read API returns captured bytes/digest, source identity and WORK status
without locking/persisting a commit; detect changed inputs with bounded retries and
truthful diagnostics. No semantic blueprint command is delivered by this adapter.
Use existing presentation registry, no Bash helper and no dynamic plugin system.

Acceptance: compiled CLI journey (initial WORK -> edits -> commit -> dirty restore ->
merge/conflict resolution -> candidate -> release -> new WORK -> copied project
inspection), both human and --json output, interrupted/restarted operations and
ordinary projects unchanged. Validate source discovery excludes deep merge archives.
Independent reviewer checks owner requirements, model/contracts and evidence, with
all high-impact findings resolved before declaring this plan complete.

## Verification coverage and limits

| MG3 gap | Required milestone | Evidence level |
| --- | --- | --- |
| Complete schema, UUID/hash/path/size rules | MGI1-M1 | Negative schema/domain fixtures and process CLI tests |
| Transactional live restore, dirty editors and lock recovery | MGI2-M1 | Filesystem/process interruption at every journal boundary |
| Deep archives, merge bases, lineage/event idempotence | MGI3-M1 | Multiple WORKs, fresh-process replay and adversarial graphs |
| Production text merge without external Git | MGI3-M1 | Dependency decision and disjoint/conflicting same-file fixtures |
| Real language validation and dropped payloads | MGI4-M1 | Actual SDL/SDUI sources and deleted origin WORKs |
| Competing accepted releases and evidence claims | MGI4-M2 | Two-clone transport plus acceptance-policy negative cases |
| CLI, discovery, output and regression behavior | MGI5-M1/M2 | Compiled end-to-end commands and full relevant suites |

Run scoped tests during each change, then SDPTool go test -race ./... and go vet ./...
at integrated delivery. Run relevant SDL/SDUI tests if their adapters change; do not
rewrite their parsers. Recheck canonical SDL and regenerate changed viewpoints from
source only. Record exact candidate/hash, OS/toolchain, commands, results and failed
or unverified obligations. Planning does not require rerunning unchanged MG3 evidence.

Implementation should remain portable Go. Initial acceptance is Linux process-level
recovery; native Windows/macOS runtime acceptance and physical power-loss guarantees
are not claimed. Cross-compilation is a compatibility check, not native verification.
Unsupported platform mutation must be explicit before writes; a platform support
matrix belongs with implementation evidence. Expand claims only with measured tests.

## Migration, scope guard and completion

No automatic move of existing SDP/SDL source trees into WORK. First use is explicit
on a disposable or owner-selected model copy. Preserve original files and publish
clear copying/initialization instructions. Installer manifest/distribution changes,
release packaging and gh-sdp upgrade belong to a separately authorized release plan.
Do not silently change installed process version while adding a tool capability.

Completion requires all milestones, working docs/help, stable schemas and fixtures,
independent review and honest platform/evidence boundaries. Update supported SDL
model activities only when implementation evidence exists; diagram generation alone
is not delivery. KB049 then closes with outcome evidence. Session0006 may close when
its goal/dispositions are satisfied; KB050 remains its own blueprint workstream.

## MG4 preparation evidence

2026-10-05: created from Contract.md, the real SDPTool CLI/discovery boundary and MG3
proof/limitations. Management/Toolkit validation and diff checks are run on delivery.
No tests above have been executed for the production implementation, no implementation
milestone is delivered, and no new product commands or releases are claimed.

2026-10-05: Owner authorizes all five implementation phases overnight. Execution
started on sdp/model-governance-implementation; milestone evidence follows below.

MGI1-M1: WORK creation/status and strict metadata foundation implemented, with staged area-locked publication. SDL/SDUI validation and recovery library scaffolding included but later CLI operations not yet exposed. Full SDPTool test suite passes.

MGI2-M1: Commit/history/whole-state restore and explicit recover resume/abort CLI implemented. Race tests pass, including child-process interruption at prepared/backup/installed boundaries, dirty preservation, corrupt payload and external edit refusal. Backups retained; physical power-loss and non-Linux mutation not claimed.

MGI3-M1: Native bounded Go three-way merge, combined WORK creation, persistent conflict inventory and resolved commits implemented. Tests cover disjoint and overlapping same-file edits, unrelated bases, dirty inputs, repeat integration without extra events, three archive generations and whole-state rollback. No external Git merge dependency.

MGI4-M1: Frozen candidate/proposal CLI now validates real SDL source graphs and SDUI using existing parsers, preserves metadata lineage and drops undo payloads. Model tests pass. Independent review found dirty-capture identity and file/directory restore defects; regression fixes and stricter domain validation are included, with final re-review pending.
