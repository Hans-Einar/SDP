# R240 / SPS-014 independent combined release review

Disposition: approved for paired2.4.0 release closeout. All required observed
release and public consumer gates pass; no blocking findings or unresolved release
evidence conditions. Earlier conditional checkpoints below remain historical.

Candidates: canonical52888aaf96c657c0d6656b104c9594e06f99d46f at
/tmp/sdp-r240-source; extension f7763acf27900c0fb51fd28d3f2b6edc8a40a5f7
at /home/warloc/git/gh-sdp. Both working trees clean during inspection.
Authority: owner-selected MAINT-SDP-0023 / SPS-014, Activity0007 / Session0027.
Read-only independent review; no publication, product repairs, new agents or broad
retests. Loaded/reused sdp, shared document workflow, reviewer, release and
versioning guidance; extension project/installed instructions from this context.

Reused completed PLAN-SDP-0030 independent implementation review and exact735c865
verification (including real native painting/reload/rejection/recovery evidence).
Confirmed no SDUI/SDL/product logic diff follows that candidate; SDPTool delta is
release version/default/predecessor configuration and predecessor test only.
Inspected release changes since canonical8f1fb9a and extension c0b6f0f.
Minor2.4.0 correctly describes additive explicit-context DocumentHost launch and
preflight/document/nativeLayout-pending checks, retaining legacy0.2 without context.
Notes state explicit context, unbound callbacks, no XFMD context picker and separate
Framework upgrades. Version/default/package/manifests and truthful prerelease/null
publication fields agree. Extension runtime, packaging and dependency source unchanged.

Independent checks executed:
- Every canonical SHA256SUMS entry, production Ed25519 signature against the
  compiled trusted Linux publisher public key (OpenSSL), all73embedded inventory
  file-content hashes, descriptor source/version, and all10supported predecessors.
- Identical canonical/extension engine, descriptor and signature bytes; extension
  distribution manifest commit and executable digest match exact candidates.
- Historical notes from2.3.2 onward and old management/traceability ledger bytes
  remain preserved in both repositories. All8existing Session snapshots,
 14descriptors,45operation records and10released generated logs are byte-identical
  to baseline; all retained Session snapshot SHA-256 filenames match content.

Production descriptor SHA-256:
7184efb8b5a4b1fa9280f61f5e9c839b4d6dc14ee58173a81478a7eaaa815908.
Production canonical/extension engine SHA-256:
2034c184c7443d7cb90bba6a30926ed0c7f99ae12063407c5ec0f282e415b9f9.
Publisher key ID:2c947457cffe015f870168cb332242c49606fdea378a13d875579544c2989c03.
No private key accessed. Linux amd64 only; no new platform qualification claimed.

Inspected passing profile/install race (35.495s/41.744s), bootstrap(1.207s),
90Toolkit tests, packaged signed-child install, extracted-source Toolkit/management
validation and all11generated release logs. Archive descriptor differs from
production only in release identity. Extension focused packaging race/vet pass.
Inspected ten original signed predecessor rehearsals and script:0.2.0,0.2.1,1.0.0,
2.0.0,2.1.0,2.2.0,2.2.1,2.3.0,2.3.1,2.3.2 all upgrade to signed2.4.0 with owner
files/Sessions/SDL preservation, managed-file conflict refusal and repeat no-op.

Still pending at this checkpoint: exact-production native result (running),
required CI and real paired-public-release/consumer reconciliation. Source-only
and local-package success do not establish public installability or owner acceptance.
Canonical Session/plan/checklist closeout belongs to coordinator after observed gates.


## Exact production native completion

Native release qualification approved. Inspected the exact signed-package harness,
/tmp/sdp-r240/native-smoke-result.json and native.log. Actual launcher opens a
native window, publishes ready, acknowledges both SVG resource reloads, diagnoses
invalid source, recovers at generation5 and exits0 through WM_DELETE_WINDOW.
No startup pane-minimum diagnostic. Independently viewed initial capture: visible
native controls, panes, choices and supplied timber-frame SVG. Independently ran
ImageMagick AE comparison initial/recovered:0. This supports actual painting and
last-valid recovery; reused host/watcher tests supply detailed state assertions,
without claiming fresh domain callbacks, exhaustive interaction or new platforms.
No source candidate change. Durable R240 Verification.md reflects this result.

Publication readiness: approved conditional only on actual successful required
CI44287504558 for the reviewed source candidate. No other pre-publication evidence
condition remains. Actual paired public tags/assets/install/upgrade/default-Framework
reconciliation must follow observed publication; no public success is preclaimed.


## Final public reconciliation

Final disposition: approved for R240/SPS-014 closeout. Required CI44287504558 is
completed/success, head52888aaf96c657c0d6656b104c9594e06f99d46f; contracts and
go-installation-linux both succeeded. Production identity remains established by
the signed exact artifact checks above, separately from CI source evidence.

Independently queried public remote annotated tags. Canonical v2.4.0 object
c3931676e50726eea5eab952ff7b8987077c9d33 peels to reviewed canonical52888aaf;
extension object2617b9c6d8ae3af66bb0bde94e034c9c5271ab3a peels to reviewed
f7763acf. Public releases449873361 and449873710 are non-draft/non-prerelease,
published canonical2026-11-18T11:13:32Z then extension11:13:52Z. Independently
compared every downloaded canonical/extension asset to its prepared production
package and API hash/size: all match. The trusted descriptor/signature therefore
remain the verified bytes; public extension manifest pins immutable public2.4.0
descriptor and exact canonical/distribution identities with the same engine hash.

Inspected actual public consumer qualification script, baseline and all23successful
command outcomes. Baseline is verified actual2.3.2 engine7eab2b…b5f5bf. Ordinary
`gh extension upgrade sdp` upgrades its isolated clone to2.4.0; prior setup
registration resolves to the new engine immediately before another tool/setup call.
Public clean install, setup and direct/gh identity succeed. Offline direct version,
discovery, SDUI check and preview pass in isolated network namespaces without gh or
authentication. Project bytes and the real live extension remain unchanged. These
exercise actual public binary-release update/install paths rather than local fixtures.

Public-default Framework qualification passes preview/apply against production
signed2.4.0 descriptor7184efb8…815908 without release/test-trust overrides, preserves
owner file bytes and repeats as no-op. All three commands succeed. This isolated
qualification does not migrate any real consuming project's Framework.

Evidence: /tmp/sdp-r240-public/{CI,canonical-release,extension-release}.json and
public downloaded asset directories; /tmp/sdp-r240-public-consumers/{baseline,
commands,results}.json; /tmp/sdp-r240-public-framework/{commands,results}.json.
Coordinator owns truthful durable release/Session/plan/checklist/Activity closeout.
No main merge, live extension replacement or subjective owner acceptance is inferred.
