# ReleaseChecklist — SDP 2.1.0 publication

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
- [ ] Verify gh-sdp selection: its compiled default is pinned. Either publish/update
  the client default or document an explicit SDP_RELEASE selecting this release.
- [x] Provide preview/apply commands and expected receipt. Owner-requested manual
  project upgrade remains manual; never confuse updating the gh extension with
  upgrading the project's SDP directory.

Consuming projects reuse the scope/content/publication checks; mark SDP publisher
inventory/signing steps N/A if they do not distribute SDP itself. Keep their own
application/package/installer acceptance checks alongside this list.

Preparation checks reuse SR1 exact-candidate evidence. Owner publication authority
is recorded in Plan.md. Main merge is not selected; publication is from the
reviewed branch commit. Paired client publication and owner update are distinct.
