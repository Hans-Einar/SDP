# Independent SDP 2.2.0 final release review

APPROVED — exact Linux amd64 source/package readiness under MAINT-SDP-0015 RSP2-M1. No remaining blocking, high or medium findings within this review boundary. This is not a publication, paired-client completion or main-merge claim.

## Candidate and disposition

Source: ef8741c433c7fde67e689f94a74e7afdf938f24f
Package: /tmp/sdp22-package-final
Inventory: /tmp/sdp22-final-candidate.json

The sole source delta from rejected 62516212b4cf5ca891fad27c7f1859a0d81a229c is the generated Releases/2.2.0.md correction from PLAN-SDP-0017 to PLAN-SDP-0019. Independently inspected the actual diff. The old candidate and /tmp/sdp22-package remain rejected; this approval does not transfer to them.

Reviewed AGENTS.md, SDP entrypoint/document workflow, Reviewer 2.0.0, Release 2.0.0 and Versioning 2.0.0; RSP2 contract/checklist, installed release lifecycle and version contract; Session0010 T011 authority and S8 roadmap; program and ModelGovernance scope/evidence/review. Publication authority already exists; no main merge is selected. This read-only assignment overrides journal mutation: coordinating agent must reconcile Session0010 and RSP2 records with this result.

## Independently verified

- Working checkout HEAD equals the exact source above and git status --porcelain is empty, checked at start and closeout. The remote v2.2.0 tag was absent and actual latest public release remained v2.1.0 during review.
- Rebuilt with SDP_VERSION=2.2.0 through SDPTool/package.sh into /tmp/sdp22-final-independent-rebuild. Executable is byte-identical to final signed package. Runtime --version --json and manifest identify 2.2.0, exact ef8741c revision and sdp-install-command/1 without a dirty suffix.
- SHA256SUMS and all final inventory hashes agree. Production Ed25519 signature independently verified using the publisher public key compiled into source, with matching key ID. Descriptor source/release, linux/amd64 asset size/hash, executable and manifest agree.
- All 64 embedded payload bytes, SHA-256 values, destinations and ownership agree with actual source and payload inventory. Inventory is unchanged from v2.1.0. No private Session history is distributed. All five predecessor descriptors have original production signatures and hashes explicitly allowed by the final profile.
- Final signed and archive trials each contain ten completed cases: fresh/custom owner content for 0.2.0, 0.2.1, 1.0.0, 2.0.0 and 2.1.0. Inspected /tmp/sdp-22-trials.py assertions rather than relying on labels. It exercises install/apply, upgrade/apply, preserved Sessions/programs/SDUI/legacy navigation/ledger prefix, Session-template presence, discovery and repeated no-change/no-actions. Both trial inventories bind the final executable, actual target descriptor and exact source; each adds a clean install without navigation registry. Archive receipts correctly retain local-development provenance.
- The extracted final source archive has no .git and all 5222 tracked files/symlinks match the reviewed checkout. Final archive validation log passes.
- Independently exercised actual installed gh-sdp with explicit SDP_RELEASE pointing at the final production-signed local descriptor and fresh cache/project, without test-key or offline overrides. Install/apply/discover succeeded, receipt recorded exact final source and signed provenance, and Maintenance/history were created. Separately reran TestPackagedSignedInstall against the actual final binary; its test-signed fixture is supplementary evidence, not production-signature evidence.
- Independently reran all six release-log checks and Toolkit/management validation successfully. All earlier released note sections and five generated release logs are byte-preserved. REL-2.1.0 differs from its tagged preparation state only through historical verification/publication reconciliation commits ed55e56 and 746afa8; it is unchanged since that reconciliation.
- PR54 head is ef8741c433c7fde67e689f94a74e7afdf938f24f. GitHub run 38013745278 completed successfully in both contracts and go-installation-linux. Race suites, immutable release logs, packaging, actual packaged child, descriptor build and retired-entrypoint rejection all passed; none of the previously skipped packaging gates remains skipped.

CI: https://github.com/Hans-Einar/SDP/actions/runs/38013745278
PR: https://github.com/Hans-Einar/SDP/pull/54

## Exact hashes (SHA-256)

- sdptool: 4b6eca0a6db7953efa037db846f6ac34370847fa695c2de49f35af8eb15c0ef9
- sdp-release.json: d47ce2be184c680422958f6c161e1e4ed19dff06cfd471c0ed0872d1445ec61b
- sdp-release.json.sig: 99d631ff201339ff7eeabef41dd0c087858d3209ec5eacec1f7ec6b3be1580e9
- sdptool.manifest.json: a27bedd0059528aacd858eac913bdb6c630178ba01af92b73e91665195dd8162

## Reused evidence and limits

Reused the prior independent review process at /tmp/sdp22-release-review-process.log, which rejected only the immutable-log mismatch after its broader checks, and /tmp/sdp-rsp2-fix-review.md for terminal/SIGTERM corrections. All six runtime/test file hashes listed in that correction approval match current source. SDL/SDUI are unchanged from accepted 04f88ff; ModelGovernance module is unchanged from reviewed fd7033b. Examined the supplied passing SDPTool/bootstrap/SDL/SDUI/Toolkit/validation/management suite logs; did not redundantly rerun the full unchanged product suites. Empty vet logs alone are not proof, but the prior review process records successful independent vet completion. Final CI additionally supplies fresh SDPTool/bootstrap race evidence.

The trial script does not itself test conflicts, every receipt field, or all native UI behavior. Conflict coverage comes from unchanged installer tests (including the prior independent TestKnownUpgradeAndLocalEdits run), and source/receipt binding is checked independently here. Existing bounded native/widget and program integration evidence is reused, not claimed as rerun on this signed binary. Non-Linux execution, physical power-loss durability, authenticated model acceptance, semantic blueprints, XFMD Run adoption and every native widget combination are not newly established.

Only Linux amd64 is approved. Version 2.2.0 is consistent with the selected additive capability scope; Framework remains 2.0.0 and language/process/schema identities remain separate. Public download verification, annotated tag/release identity, paired gh-sdp 0.2.2 compiled-default verification and local extension upgrade necessarily remain publication-stage work. In particular verify the public default after upstream SDP publication and before client publication. No publication occurred in this review.

Pending review/evidence/checklist/release scaffolds truthfully represent the pre-review state and are not defects. Coordinating agent must record this disposition and exact evidence, reconcile preparation identity, then record only actually observed publication results. No governed record or source was edited, and nothing was published or merged.

## Temporary independent outputs

- /tmp/sdp22-final-independent-integrity.py and .log
- /tmp/sdp22-final-bootstrap-review.py
- /tmp/sdp22-final-independent-bootstrap.log
- /tmp/sdp22-final-independent-child.log
- /tmp/sdp22-final-independent-ci.json
- /tmp/sdp22-final-independent-rebuild/

## Publication reconciliation

Review-publication.md independently approves both actual releases, asset identity,
global extension/default and the bounded native lab evidence. This completes the
release review; original source-candidate dispositions above remain intact.
