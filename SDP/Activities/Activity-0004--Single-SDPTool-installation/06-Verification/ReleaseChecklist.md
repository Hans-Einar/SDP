# Activity4 release checklist — SDP2.3.0 / gh-sdp2.3.0 / XFMD0.11

Owner authorizes Activity4 execution of the selected single-distribution migration,
continuing the requested XFMD delivery. No main merge or live project upgrades.
Predecessors: canonicalSDP2.2.1, extension0.2.3, XFMD0.10.0-76. Exact source,
artifacts and public outcomes are retained in the evidence below.

# ReleaseChecklist

Copy this checklist into the selected release's work area. Record version,
previous release, exact commit, owner publication authority, evidence paths and
actual outcomes. Unchecked means pending; use N/A only with a reason. A generated
release log is not proof of publication. The project release contract governs.

## Scope and content

- [x] Select product version from compatibility impact. Keep language/profile,
  process/management schema and client versions distinct.
- [x] Enumerate included work and remaining exclusions; validate selected model,
  tests, review dispositions and links on the actual candidate.
- [x] Update canonical RELEASE-NOTES.md Unreleased entries with stable work IDs.
  Freeze selected entries into a dated version section; preserve released bytes.
- [x] Run `sdptool release-log --version X.Y.Z --output Releases/X.Y.Z.md`.
  Check with the same command plus --check. Use --all --output Releases --check
  to check every recorded version. Generator extracts reviewed notes, not Git guesses.

## Installation distribution (SDP publisher)

- [x] Check Template/sdp-root and project-root guidance. Include Sessions guide and
  blank template, never project conversations, boards or existing work records.
- [x] Review SDPTool/profiles/payload.json: every distributed file, destination,
  ownership and removed/moved path. Inventory completeness tests must pass.
- [x] Review SDPTool/profiles/five-phase.json: capabilities, schema/profile versions
  and exact supported predecessor descriptor digests, including current stable.
- [x] Reconcile SDP.manifest.yaml, Toolkit/SDP-install.manifest.json compatibility
  metadata, package version, release record and notes. Retained legacy metadata
  does not replace the current Go inventory or signed descriptor.
- [x] Build package and immutable sdp-release.json from a clean exact commit using
  the Go profile builder. Verify file hashes, sourceCommit and executable identity.
  Never edit the built descriptor or use a project-local manifest as authority.
- [x] Sign descriptor with the trusted publisher key; verify signature, binary hash
  and SHA256SUMS. Store no private keys or secrets in evidence/repository.
- [x] Rehearse install and each supported upgrade from original verified descriptor
  bytes. Check new Sessions files, existing Session/document preservation, receipt,
  Maintenance entry/history, conflict handling and repeated no-op.
- [x] Test extracted-source operation, packaged signed child and normal bootstrap.
  Development artifacts or test signers are not production signed evidence.

## Publish and consumer readiness

- [x] Exact candidate CI and required independent review pass; working tree clean;
  tag absent; merge/publication explicitly authorized.
- [x] Publish annotated tag and GitHub Release with generated log, signed descriptor,
  executable and checksums; verify downloaded assets against exact release.
- [x] Reconcile actual tag/commit/time and append truthful release events only after
  success. Do not rewrite a released log or silently correct released notes.
- [x] Verify gh-sdp selection: its compiled default is pinned. Either publish/update
  the client default or document an explicit SDP_RELEASE selecting this release.
- [x] Provide preview/apply commands and expected receipt. Owner-requested manual
  project upgrade remains manual; never confuse updating the gh extension with
  upgrading the project's SDP directory.

Consuming projects reuse the scope/content/publication checks; mark SDP publisher
inventory/signing steps N/A if they do not distribute SDP itself. Keep their own
application/package/installer acceptance checks alongside this list.

## Gate evidence and applicability

- Scope/version/notes: canonical2.3.0, same-byte gh-sdp2.3.0 and XFMD0.11.0-81
  implement PLAN-SDP-0029 C0–C5. The dated released notes/log were frozen on
  source774e150; Candidate/activity4-extracted-check.log checks all8 release logs.
  Current requirements and ImplementationMap document the bounded model handoff.
  ModelValidation.json checks corrected preliminary checkpoint00002 at closeout;
  checkpoint00001 had no semantic-validation receipt. No accepted model release
  is claimed. Source formatting is a post-release model-only correction.
- Templates/inventory/profiles: Candidate/activity4-profile-final.log,
  activity4-toolkit-tests-final.log and activity4-toolkit-validation.log retain
  inventory, payload, template and schema checks. Actual signed7-predecessor
  preservation/conflict/no-op evidence is retained in gh-sdp SPS-011/774e150.
- Clean package/signing: Candidate/Provenance.json, SHA256SUMS, signed descriptor
  signature and manifest bind clean source774e150. Public/ asset receipts and the
  independent review verify production Ed25519 trust and exact engine bytes.
  No private key or secret is included. Engine Go VCS metadata is absent; explicit
  executed BuildRevision and signed descriptor bind source, not a fabricated stamp.
- Install/extracted source/child: Candidate/activity4-packaged-install-774e150.log,
  activity4-extracted-check.log and Public/framework-{results,receipt,commands}.json.
  The old runtime bootstrap child is retired by this migration (N/A as a shipping
  architecture); actual gh binary extension install/update and public-default
  signed Framework preview/apply replace that consumer gate. Existing bootstrap
  compatibility tests and packaged child tests passed where retained.
- Exact CI/review: Public/CI.json and IndependentReview.md. Run43398054812 passed
  exact source774e150. Clean source and absent tags were checked before publication;
  the publication tags now exist intentionally. Owner selected publication; main
  merge is N/A/unselected and technical approval does not claim personal review.
- Public delivery: Public/Publication.json and canonical/extension release receipts
  record actual annotated tags, server times, downloads and asset byte equality.
  Public/consumer-{results,commands}.json covers normal old0.2.3 extension upgrade,
  fresh install/setup and offline direct commands. The extension ships the exact
  engine; its public-default signed Framework selector is2.3.0.
- Native application: NativeCandidate/ records95 CTests and exact local package
  qualification. Public/Native retains trusted signature, identical payload,
  actual signed public predecessor upgrade and fresh-install native workflow logs.
  XFMD has no private SDP engine or helper payload. OS GUI dependencies are declared.
- User commands: `gh extension upgrade sdp`, one-time `gh sdp setup`, then
  `sudo dnf upgrade --refresh xfmd`. Framework changes remain a separate explicit
  `gh sdp PROJECT upgrade --plan-output PLAN.json` followed by
  `gh sdp PROJECT upgrade --apply PLAN.json`; Public/framework-commands.json
  receipts retain the exact tested CLI and signed resulting receipt. Existing
  user project bytes and the live extension were preserved during isolated tests.

The checked boxes above mean applicable gates have passed, including the explicit
N/A distinctions here. They do not assert a main merge, accepted model release,
automatic extension post-install hook or live user environment upgrade.
