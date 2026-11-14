# SDP2.3.1 release checklist

Owner publication/HEOS rollout authority: MAINT-SDP-0020. Previous canonical and
extension release2.3.0; HEOS installed2.2.0. Exact source84d5e46, extensionae7aa21.
See Verification.md and IndependentReview.md. Checkboxes below cover release gates;
the actual live HEOS rollout is also verified complete.

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

## Applicable evidence and limits

- Version: compatible patch default/predecessor/metadata changes, no runtime API or
  schema change. Exact clean source84d5e46 contains integrated source and the missing
  immutable Session creation snapshot; initial b428a3a was rejected, not published.
- Payload inventory/template bytes are unchanged from2.3.0. Profile now accepts8
  original signed descriptors, including2.3.0. All8 actual upgrade/preservation/
  conflict/no-op checks pass; see evidence/predecessor-results.json and script.
- Source model semantics are unchanged;2.3.0's post-release preliminary model
  formatting correction and validation are included in the source. No new accepted
  model release is claimed. All9 release logs pass in source and Git-less extraction.
- Exact static engine and production Linux signature/checksums pass independent
  review. Explicit executed BuildRevision and signed descriptor bind source;
  this host emits no embedded Go VCS stamp. Native helper bytes are identical to
  the prior qualified payload. Production signed fixture tests are distinct from
  the additional test-signed packaged-child check.
-90 Toolkit tests, profile/install race, bootstrap tests and extension package tests
  pass. Runtime bootstrap-child behavior is N/A for the shared engine architecture;
  actual same-byte extension consumer installation/update is the replacement gate.
- Fresh independent review resolves the missing snapshot and inspects final sources,
  package and HEOS-copy proposal. CI uses a proven identical source tree under its
  synthetic merge identity; the final CI/public gates below remain unchecked until
  observed. Main merge is not selected; release tags preserve exact candidates.

## Final gate disposition

All applicable release gates above pass; Verification.md and IndependentReview.md
retain evidence and scope. CI43544849624 passed the source-tree-equivalent candidate;
public downloaded canonical/extension bytes, signatures and consumer workflows pass.
Actual local extension and HEOS public-default upgrade are verified. Immutable
original HEOS ledger is archived under the separately explicit owner exception.
The release log is unchanged after publication. Main merge remains unselected;
retired runtime bootstrap architecture is N/A, replaced by actual shared-engine
extension qualification. No Windows/macOS or new native behavior claim is added.
