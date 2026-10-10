# SDPTOOL-VER-RSP2 — release evidence

## Preparation and corrections

Release scope is MAINT-SDP-0015, authorized by owner Session0010 T011. Source
starts at RSP1 handoff 0652e174e4fbd126ca713c08ed5ec2465fa5e5f6. SDP 2.2.0 is
additive; Framework remains 2.0.0 and only Linux amd64 is published. gh-sdp 0.2.2
selects the new default without changing the thin client API. No XFMD source or
live project installation changes are selected.

Fresh independent review of RSP1 found terminal input could stop a background
runner and SIGTERM bypassed CLI cancellation. RSP2 transfers terminal foreground
ownership for interactive input, restores it afterwards and routes SIGTERM through
context cancellation. Real-executable Linux PTY and descendant-cleanup regressions
pass, alongside the original cancellation test. The reviewer reruns those checks
independently; the signed-package gate remains separate.

Repository validation, management validation and all 85 Toolkit tests pass after
aligning the maintained installation examples with 2.2.0. Original released notes
and descriptor digests are preserved. The payload inventory is unchanged; only
the managed SDP entry skill differs from 2.1.0 (previously accepted Session upkeep).

Exact clean package, production signature, predecessor trials, archive gate,
CI and final independent approval will be recorded below when observed.

## Final exact publication gate

Candidate ef8741c433c7fde67e689f94a74e7afdf938f24f; Candidate.json records actual
Linux amd64 hashes and executable identity. Initial 6251621 was rejected because
its generated log had the wrong ModelGovernance plan reference. Final candidate
changes only that generated line and was rebuilt/signed; both exact-head CI jobs
pass at https://github.com/Hans-Einar/SDP/actions/runs/38013745278.

Full product suites, bootstrap race/vet, Toolkit/management validators, all 85
Toolkit tests and final packaged-child check passed. Signed-upgrades.json and
Archive-upgrades.json each record ten predecessor cases plus fresh install.
The reviewed trial script verifies owner Sessions, programs.json, SDUI source,
legacy navigation and ledger prefix preservation, followed by repeat no-op.
The final archive has no Git metadata and passes validation. Production trust
is separate from the test signer used in the supplementary packaged-child test.

Review.md independently APPROVES this exact source/package, after reproducing
its binary, signature, all 64 payload entries, original predecessor signatures,
release-log checks, actual gh bootstrap/install and CI. Earlier rejected reviews
are retained with their dispositions. No unresolved blocking findings remain.
RSP2-M1 is complete; RSP2-M2 publication, public default and extension upgrade are
still pending. No change to the reviewed product is made by this evidence commit.
