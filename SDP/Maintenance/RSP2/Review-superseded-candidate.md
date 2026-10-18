**CHANGES_REQUIRED** for exact candidate `62516212b4cf5ca891fad27c7f1859a0d81a229c` and its signed package.

**P2 — Generated release notes fail the release gate.** The candidate’s generated 2.2.0 log identifies ModelGovernance as `PLAN-SDP-0017`; canonical notes correctly identify `PLAN-SDP-0019`. Independently reproduced with `release-log --all --check`. The [exact-candidate CI job](https://github.com/Hans-Einar/SDP/actions/runs/38013612822/job/114098948394) failed at that same check; subsequent packaging checks were skipped.

file:///tmp/sdp22-independent-source/Releases/2.2.0.md

file:///tmp/sdp22-exact-release-log-check.log

The shared checkout and PR advanced during review to `ef8741c433c7fde67e689f94a74e7afdf938f24f`, correcting the generated log. That does not approve the requested older candidate/package. Complete exact-candidate CI and rebuild/sign against the selected corrected commit before publication.

Other checks passed:

- Production Ed25519 signature, checksums, all 64 payload entries, and byte-identical independent executable rebuild.
- All five original signed predecessor descriptors; both completed ten-case upgrade trials, preservation and repeated no-op checks.
- Archive identity across 5,222 tracked files; archive validation and trials.
- Actual installed gh-sdp bootstrap using the production-signed local package, temporary installation, discovery, receipt and Maintenance history.
- Repository/management validation, focused program regressions, inventory/conflict checks and independent vet. Exact-candidate CI passed both race suites before the release-log failure.
- Previous released notes/logs and ledger prefixes preserved; 2.1.0 record changes are publication reconciliation.
- Runtime correction hashes match prior approval. ModelGovernance implementation and SDL/SDUI widget sources match their accepted baselines; no unchanged UI rewrite is required.

file:///tmp/sdp22-independent-integrity.log

file:///tmp/sdp22-exact-focused.log

Publication, downloaded-asset verification, final public gh-sdp default, paired client publication and local extension upgrade remain pending transaction steps. The public default check belongs after upstream publication and before paired client publication; it is not this source blocker. No main merge is required.

No source or governed records were edited, and nothing was published or merged. The coordinator should record this disposition and completed evidence in RSP2 and Session0010.