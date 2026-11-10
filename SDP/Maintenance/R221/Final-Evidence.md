# R221 final release evidence

Exact producer: 6f33c68ce25fb344487f565d5611a51089357632.
Exact client: 03e45710db3ad845d05b26b19c286bd069911294 (0.2.3).
Bootstrap dependency: v0.0.0-20261110232832-35bfdf8469e0.
The final producer differs from that dependency only in CI budget and evidence
documentation; no runtime or distributed payload changed.

## Reproducible evidence

Linux evidence root: /home/warloc/.local/state/sdp-rollout-2.2.1.
Production artifacts are producer-final/ and client/; earlier producer/ is superseded.

| Artifact | SHA256 |
| --- | --- |
| sdptool | 1c112c24aee360cd6da406a0ebdfc58db6eaf96148e33b5065264f36b82a907f |
| sdp-release.json | 11e76cffd507e06aa36e07d56f3cf4f121c47393bf8383786c47ccbd202b5ea1 |
| sdp-release.json.sig | 81b4a0cc4f8f6a8939c774c30d555924bf7a11cf3e1c8c7e6276e89adbf42cb9 |
| gh-sdp-linux-amd64 | 5631dc2e849da5598606394319451be94f6023e0eb7e9b4edc6c82e1647e39c5 |

Linux publisher key ID: 2c947457cffe015f870168cb332242c49606fdea378a13d875579544c2989c03.
Private material remains outside repositories and evidence. The previous public
trust anchor is retained. Portable static Linux amd64/v1 packages use clean source.

- Full local engine suite, bootstrap race tests, vet, Toolkit validation and85
  Toolkit tests pass (product unchanged after35bfdf8).
- Client race tests, vet and module verification pass at03e4571.
- final-signed-upgrade-results.json: original production-signed0.2.0,0.2.1,
  1.0.0,2.0.0,2.1.0,2.2.0 installed and upgraded to final signed2.2.1;
  owner files preserved, modified managed-file conflicts refused, repeat no-op.
- final-client-packaged-tests.log: actual packaged client and final engine,
  isolated gh dispatch, human/JSON discovery, plans/apply, preservation,
  repeat and corrupt-cache cases pass.
- final-client-local-signed-version.json: final engine signature/default client
  chain verified with local descriptor selection. This is not public-default evidence.
- final-extracted-check.log: all seven generated logs checked from exact Git-less source.
- CI: https://github.com/Hans-Einar/SDP/actions/runs/42925278417
  is the exact6f33 run; publication waits for success.

## Independent review

REV-SDPTOOL-R221-FINAL / REV-SPS-010: independent reviewer
review_release_rollout inspected exact source, package provenance, public-key
signature (including independent OpenSSL verification), predecessor trials and
client evidence. No findings. Publication is conditional on exact CI success,
record reconciliation and actual downloaded release/default verification.
Review is distinct from owner publication authority, explicitly supplied2026-11-11.

## Consumer rollout scope

Inventory records seven schema3 signed installations under Linux ~/git.
Applicable previews: kitty2.1.0, sdui_widget_lab2.1.0 and xfmd2.0.0.
Original dirty files, Git index and history are backed up and hash-recorded.
HEOS, emuSA80535-N, ponsse and ponsse-service-continuity remain blocked by
pre-existing invalid management history; do not rewrite immutable event bytes.
The gh-sdp repository's schema1 manual Toolkit installation is excluded;
adoption is a separate change. No application builds or main merges are selected.

This record initially describes verified preparation. Actual publication,
installation and per-repository results must be recorded separately below.
