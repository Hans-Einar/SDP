# SDP 2.3.2 release qualification

| Field | Value |
| --- | --- |
| id | SDPTOOL-VER-R232 |
| Status | completed |
| Activity | ACTIVITY-SDP-0006 |
| Plan | MAINT-SDP-0022 |

Canonical source: f7a6c9f4793b55de5d0017ae74210135f2a4d365.
Extension packaging source: c1d86bd148e6ad24da55b03482000df68e3fa3d6.
Engine SHA256: 7eab2b30286fe1608e95cc8011d07632dca7b73750486979130cfa8094b5f5bf.
Signed descriptor SHA256: 34e2fa36ad45047ca398612cbc99f0a30430c44250a5990e98bbaed2f13efac8.
Both clean candidates identify the same production-signed engine. No SDL design
facts, wire schemas, native application implementation or payload ownership change.

## Candidate evidence

Evidence in evidence/ records focused profile/install race and bootstrap tests,
90 Toolkit tests, source/archive validation and generated release logs, actual
packaged signed-child execution, and nine supported original signed predecessors.
Each predecessor installs, rejects a managed-file conflict, upgrades with owner
files preserved and repeats as a no-op. The archive-generated descriptor differs
from production only in its explicitly unsigned development release identity.

A 128-source composed SDL model is refused by the published 2.3.1 engine at the
old limit; the actual signed 2.3.2 binary publishes 135 files and all 128 source
links. Exact 512-file / 128-MiB bounds and overflow are covered by BC1 evidence.

The native helper source and dependencies are unchanged, but its bytes differ:
the prior build used Red Hat Go 1.26.10 with FIPS tags; this build uses Go 1.27.1.
An isolated Xvfb startup smoke confirms the newly bundled helper stays running
with the controls fixture, with no panic. The first attempt correctly refused an
undeclared default entry; the retry selects the actual controls frame. This is
startup evidence, not a new visual/interaction campaign. Historical behavior
evidence is reused for unchanged source, without claiming byte identity.

## Final observed delivery

CI 44044114546 passes for the source-tree-equivalent candidate (see evidence/ci-source-tree.json). The independent review records
its exact source/tree relationship and production artifact checks. Annotated
v2.3.2 tags and both public releases are verified; all downloaded assets match the
reviewed local packages and GitHub API SHA256 digests. The extension is repackaged
from the public immutable canonical descriptor before publication.

An isolated ordinary `gh extension upgrade sdp` from 2.3.1 and clean extension
installation pass. The existing setup symlink resolves to the new binary before
any repeated setup; direct disconnected commands and project preservation pass.
A fresh isolated Framework install uses the compiled public 2.3.2 default, receives
a signed receipt and repeats as a no-op. The user's live extension and real
project Frameworks are unchanged. No main merge is selected.

Public receipts and reproducible consumer scripts are in public-evidence/.
