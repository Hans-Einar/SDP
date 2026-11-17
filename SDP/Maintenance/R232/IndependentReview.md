# Independent combined R232 / SPS-013 review

Status: approved. Source, production artifacts and actual public distribution/consumer
reconciliation pass. No blocking findings or unresolved release evidence conditions.
Earlier conditional dispositions below retain their historical review stage.

Scope: owner-authorized compatible capacity correction, canonical SDPTool and
same-byte gh-sdp 2.3.2 distribution. No merge, global extension replacement or
consumer Framework change. Read-only independent review; no product repairs,
publication, further agents or duplicated broad tests.

Candidates: canonical f7a6c9f4793b55de5d0017ae74210135f2a4d365
(/tmp/sdp-r232-source), including BC1 delivery 4d12b27; extension
c1d86bd148e6ad24da55b03482000df68e3fa3d6 (/home/warloc/git/gh-sdp).
Compared both against their actual published v2.3.1 tags.

Loaded canonical sdp entrypoint, shared document workflow, sdp-reviewer and
sdp-release; extension managed/project instructions and installed Reviewer/release
instructions. Read BC1 and R232 contracts/evidence, SPS-013 contract/evidence,
changed source and tests, version/predecessor/distribution metadata and notes.

The shared helper enforces inclusive 512-file/128-MiB limits after sealing both
preview and source-reading output, includes retained sources and metadata, and
returns existing limit error classification with actual counts/bytes. The check
precedes publication, preserving ownership, validation and failure behavior.
This directly addresses the requested fourfold generated-resource capacity.
No SDL/source-graph limit changes or runtime ownership changes were introduced.

Evidence inspected: BC1 focused test log covers file/byte boundary acceptance,
overflow rejection, composed 128-source model producing >128 generated files
through both publication paths, freshness and failed-refresh preservation. Inspected
new tests directly. Existing source/publication invariants remain unchanged.
Reuse is adequate for this narrow correction; no native behavior expansion claimed.
Canonical profile/install race tests (33.467s/38.426s), bootstrap (1.201s), 90 Toolkit
tests, extracted-source Toolkit/management validation and all 10 generated release
logs pass in /tmp/sdp-r232. Extension focused packaging race/vet evidence passes;
production provenance remains separate from fixture tests.

Independent read-only checks executed: historical notes from the 2.3.1 heading
onward are byte-identical in both repositories; each existing management/traceability
ledger remains a byte prefix. All six retained Activity Session snapshots match
their SHA-256 filenames. All four pre-existing Session snapshots, eight descriptors,
22 operation records and nine generated released logs are byte-preserved against
canonical v2.3.1. No extension runtime, packaging or dependency source change exists.

Patch version 2.3.2, bootstrap default, package/development identity and both manifests
agree. Predecessor inventory preserves existing eight entries and adds signed2.3.1
digest 9c83225b3a5eae1762e5609965378b8c14d89fa676ae242856b323c1dd9facf6.
Release records correctly leave real publication identities unset at this stage.

Pending: exact production package signature/checksum/source identity; signed original
nine-predecessor upgrade/preservation/conflict/no-op evidence; required canonical CI;
extension exact-engine package identity/provenance; final ready disposition. Actual
public downloads, tag identities and isolated real gh upgrade require post-publication
reconciliation in this same independent review context.

## Production artifact review

Production artifacts approved on the same exact candidates, conditional only on
required exact-candidate CI success before publication. Canonical and extension
working trees are clean. No new blocking finding.

Independently verified every canonical SHA256SUMS entry; production Ed25519
signature using the compiled trusted Linux publisher public key and OpenSSL;
every embedded descriptor file-content hash; descriptor and signature byte equality
in extension distribution; canonical/extension executable byte equality; and
extension manifest digest. No private key was accessed.

Canonical descriptor SHA-256:
34e2fa36ad45047ca398612cbc99f0a30430c44250a5990e98bbaed2f13efac8.
Canonical/extension engine SHA-256:
7eab2b30286fe1608e95cc8011d07632dca7b73750486979130cfa8094b5f5bf.
Production signature key ID:
2c947457cffe015f870168cb332242c49606fdea378a13d875579544c2989c03.
The descriptor identifies canonical f7a6c9f4793b55de5d0017ae74210135f2a4d365,
73 inventory entries, nine supported predecessors and linux/amd64 engine2.3.2.
Extension manifest independently identifies distribution c1d86bd148e6ad24da55b03482000df68e3fa3d6.

Inspected all nine signed predecessor outcomes and the executing rehearsal script:
original signed descriptor installation, target upgrade, owner-document/Session/SDL
preservation, managed-file conflict refusal, signed receipt and repeat no-op pass
for0.2.0,0.2.1,1.0.0,2.0.0,2.1.0,2.2.0,2.2.1,2.3.0,2.3.1. Packaged signed-child
install test passes. Extracted-source descriptor differs from production only in
release identity. Public2.3.1 rejects the composed fixture at the former128 limit;
new packaged engine publishes135files from128sources.

Native helper source/dependency inputs remain unchanged, but binary bytes differ
with Go1.26.10 Red Hat versus Go1.27.1 build toolchains. Inspected metadata and the
isolated Xvfb startup smoke: declared controls entry remains alive after5seconds.
Initial default-page fixture rejection was corrected and retained in the result.
This establishes startup only; prior unchanged native behavior evidence is reused,
without claiming a fresh visual/interaction matrix or platform expansion.

Required CI44044114546 is still pending as reported by coordinator; publication
approval is conditional on actual success for the exact canonical candidate.
Post-publication public downloads, tag identities and real isolated gh upgrade
remain mandatory reconciliation conditions, not pre-publication facts.


## Final public release reconciliation

Final disposition: approved for R232/SPS-013 release closeout. No blocking findings.
CI44044114546 completed successfully: go-installation-linux and contracts both
succeeded; receipt identifies f7a6c9f4793b55de5d0017ae74210135f2a4d365. Coordinator
observed the PR synthetic merge d9849020f21ad52d06ad691b1dbf081670208ad7 tree equals
the production source tree69f74305e56e16e39905160f9325933cc810164d; this is source-tree
CI evidence, separate from production build identity/signature checks above.

Independently queried both public remote tag refs. Canonical annotated tag object
26fd52eec8b49e62735f905e04e0f8b0ec300b64 peels to the reviewed canonical commit;
extension annotated tag object4d3d683871903a368f0addbc7717a24010c8ff80 peels to the
reviewed extension commit. Both public v2.3.2 releases are non-draft/non-prerelease,
published canonical08:37:05Z then extension08:37:25Z on2026-11-17. Release IDs are
448347770 and448348099 respectively. API target_commitish=main is not the actual
tag target: peeled public tag refs verify the selected branch commits; no merge
is inferred.

Independently compared all downloaded public canonical and extension assets to
the prepared production packages and their GitHub API SHA-256/size fields. All
match, including the previously verified trusted production descriptor/signature.
Public extension manifest uses the immutable public2.3.2 descriptor URL and pins
the exact reviewed distribution/engine commits and same-byte engine digest.

Inspected consumer qualification script, all23command results and final outcomes:
normal unpinned `gh extension upgrade sdp` from verified actual2.3.1 bytes replaces
the registered engine before another tool/setup invocation; public clean install,
setup and direct/gh version/byte identity pass. Direct offline version/discovery/
SDUI check/preview pass in isolated network namespaces without gh or credentials.
Project bytes and real live extension remain unchanged. These are actual public
binary-release install/upgrade checks in isolated homes, not local symlink fixtures.

Public default Framework preview/apply independently reports signed2.3.2 descriptor
34e2fa36ad45047ca398612cbc99f0a30430c44250a5990e98bbaed2f13efac8 without release
or test-trust overrides; repeat upgrade is a no-op and owner file bytes are preserved.
All three recorded Framework commands succeed. This isolated qualification does
not claim an upgrade of any real consuming project's Framework.

Evidence: /tmp/sdp-r232-public/{CI,canonical-release,extension-release}.json and
public asset directories; /tmp/sdp-r232-public-consumers/{commands,results}.json;
/tmp/sdp-r232-public-framework/{commands,results}.json. Review is read-only apart
from this report; canonical Session, plan/checklist and truthful publication
closeout updates remain the coordinating agent's responsibility. Owner acceptance,
main merge and replacement of the owner's live extension are not implied.
