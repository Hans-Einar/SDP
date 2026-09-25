# MAINT-SDP-0003 implementation evidence

## IU1-M2

Candidate: IU1-M1 9a5c79a plus the profile build/template/schema diff.
Python 3 on Linux. build_process_profile.py emits a reproducible self-contained
artifact; two test groups cover byte reproducibility and ten negative inputs.
The 1.0 installed schema is unchanged; 2.0 has its own closed schema.
New neutral templates are excluded from the legacy v1 inventory explicitly.

The full Python suite exposed a stale expected scenario list: SK1 had added two
conformance scenarios without updating that set. Added their existing IDs;
no expectation about installation behavior was weakened. The other failure is
the 38 pre-existing Traceability compatibility errors recorded under SK1/KB011.

## IU2-M1

Candidate: IU1 89981fd plus versioned PowerShell planning diff.
Five integration test groups pass with PowerShell 7.6.6: clean deterministic
read-only plan, legacy phase/archive link mapping, manual XFMD-shaped history
transfer, mixed-layout/link conflicts and corrupt artifacts. Shared installed
facts parsing remains closed and explicit; legacy defaults stay schema 1.0.
Fixture provenance is documented in Process-Installation.md. No live consumer
files were changed. Apply/resume remains IU2-M2 work.

The owner added typed plans and a Planning skill during execution. Registered
KB-SDP-029 without altering this Maintenance's already selected Git policy.

## IU2-M2

Delivered journaled apply/forward resume with exact-plan comparison, exclusive
writer lock, byte backups, post-failure drift checks and retry-safe report/facts/
history finalization. Thirteen integration groups passed on a frozen candidate
(307.662s, PowerShell 7.6.6/Linux); the final reserved-path and balanced/escaped
Markdown refinements also passed focused tests. Full legacy PowerShell fixture
suite passed. The exhaustive forced-exit matrix remains an IU3 completion gate,
not a completed check here.

Independent review by iu_review required four corrections: history schema/chain
validation, broader Markdown rebasing, FIFO rejection and reserved artifact
paths. All were corrected; the reviewer approved the IU2 scope after five
focused tests. Final reviewed hashes:

- Process-Install.ps1: 3af5f6e066fa07ea4ba59bc4d7245c3c7761b0e66f87aae1f51f3ed78117e2cb
- Install-SDP.ps1: 2f79d1a1d34af856fbe64feedaf944daefedeba8938dcf0e59ef20e24122021b
- build_process_profile.py: e7a35d15ae6ee7d3851758c36a71c6f84042b9fe28e88b9cd57970a994fe680f
- artifact: 47a5363f77408a9a99803cd370ba7bbe7f2b3d8e3bdce377f9a5aa1d4890ab11

The distributed management payload schema generalizes the project namespace;
this repository keeps its stricter local SDP namespace validator. Installation
allocates target-project Maintenance/event IDs from existing history/documents.
Backlog review: KB014 still owns standalone KanBan packaging, KB018 the wider
Toolkit audit, KB029 the new plan/Planning-skill proposal. None is silently
marked implemented by this installation milestone.

## IU3-M1

SDPTool now reads installed schemas 1.0/2.0, exposes early/late incomplete
installation state and has a prebuilt version/capability protocol. Native
package.sh emits the executable, manifest and checksums without changing host
PATH or rebuilding during viewing. Go race tests and vet pass. A disposable
fresh install was discovered from both project root and SDP area; its model-free
KanBan tree and full consuming-project validator pass. Added neutral release
notes because that existing consuming-project obligation still applies.

Independent consumer review found and verified corrections for non-mapping YAML,
early interruption before navigation publication and unknown operation folders.
Review approves this consumer scope; Windows/live rollout and full exit-matrix
proof remain outside that approval. Package evidence is a dirty working candidate
based on 8e36a5d, not a published release.

## IU3-M2 — completion evidence

Candidate: IU3-M1 b3f0513 plus the conformance/CI/operator-documentation and
closeout diff. The installer and builder hashes above are unchanged. Final
profile artifact SHA256:
662107441efcefc028bda51927053766a901c754681620a97033ee3d7c069239.
Host: Linux, Python 3, PowerShell 7.6.6, Go 1.27.1.

### Recovery coverage

The [forced-exit matrix](fault-matrix.txt) passed every one of the reviewed
artifact's 62 steps at backup/write/journal boundaries, plus preparation and
completion. Four isolated fixtures partitioned those indices. Actual process
termination and forward resume proved preserved history and exactly-once
finalization. This is process interruption evidence, not hardware-failure or
whole-tree rollback evidence.

The final artifact differs only by adding project-owned RELEASE-NOTES.md.
The [final-artifact test](final-artifact-recovery.txt) passed explicit exit 97
at preparation, all three boundaries of that added step and completion. It
also checked the final configuration digest, preserved ledger prefix, unique
events, no-change repetition and owner release-note preservation under force.
This is aggregate coverage: the original matrix covers unchanged engine steps,
and the added test covers the new artifact step. It is not an exhaustive matrix
run against the final artifact. Independent reviewer iu_review explicitly
accepted this evidence combination after confirming the unchanged engine hashes.
The source runner now also requires exit 97 at preparation/completion.

### Consumer and migration evidence

[Three integration groups](consumer-tests.txt) passed fresh installed SDPTool
discovery from both project root and SDP, model-free KanBan navigation, consuming
project validation, actual legacy-v1 installation to v2 upgrade, downgrade and
unknown-schema rejection, and case collisions. A separate old-entrypoint probe
rejects installed schema 2.0 with unsupported-installed-schema before mutation.
That check is now part of the automated versioned-upgrade test.

The [XFMD snapshot test](xfmd-snapshot.txt) passed against 242 recorded files
from commit cf11709e4ec9d6925b0d17d95c71f72fbf011959.
[Provenance](xfmd-provenance.json) records original hashes and the disposable
operation result. Existing history remained a byte prefix, repetition produced
no changes, and original source bytes/Git status remained unchanged. This is
not a live XFMD upgrade.

Go race tests and vet passed for SDPTool; malformed installed YAML list/scalar
regressions passed. Independent IU3 review approved corrections for those inputs,
early interruption before navigation exists and unknown operation folders.
No material review finding remains in the selected scope.

### Reproduction and limits

Use the commands and environment variables in
[conformance guidance](../../../Toolkit/conformance/install-v2/README.md).
Run Go tests/vet from Toolkit/SDPTool with the selected Go toolchain. Native
package.sh was exercised with a new output directory; its dirty build identity
is truthful and is not a release/publication claim. The Linux CI job now supplies
PowerShell and a prebuilt SDPTool, runs integration tests and the complete matrix;
configuration is delivered, but a remote CI pass is not claimed here.

The full Python run had 100 tests, one known failure and 16 environment skips
before the final extra regression tests were added. Installation/consumer groups
were separately run with explicit PowerShell/prebuilt-tool paths as documented
above. The full Toolkit validator still reports exactly the same 38 legacy
Traceability findings as the SK1 baseline: zero added or removed findings.
KB-SDP-011 owns that existing failure; this delivery does not claim an entirely
green repository-wide suite. Windows profile execution remains experimental.

Backlog disposition: KB-SDP-014 retains standalone KanBan distribution;
KB-SDP-018 retains the broader Toolkit audit; KB-SDP-029 retains typed plans
and the Planning skill. Owner review of KB-SDP-010 remains independent.
All six selected Maintenance milestones are delivered; live rollout, product
release and those separate backlog proposals remain outside completion.

Final closeout checks pass: 36 cards, eight management records, three lineage
operations and 240 events; four management tests and 15 lineage negative cases.
Documentation checks preserve 105 frozen records/ledger prefixes and 574
generated outputs, resolving 2,414 local file links and 130 fragments.
git diff --check passes. These checks prove record integrity, not GUI behavior.
