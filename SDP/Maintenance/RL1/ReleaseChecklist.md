# ReleaseChecklist — proposed SDP 1.1.0

Status: **preparation complete; NOT PUBLISHED and not ready for a normal gh-sdp
upgrade yet**. SDL 0.6 is the included language profile. Latest published SDP
remains 1.0.0. Candidate: 6fc662d plus RL1-M2 corrections/closeout; exact final
source is the commit containing this checklist. Owner selected preparation and
will manually test XFMD. Merge/tag/publication is a separate decision.

Evidence: [Evidence.md](Evidence.md), [upgrade cases](Upgrade-rehearsal.json),
[plan](Plan.md). The historical logs are generated; proposed 1.1.0 contents stay
in Unreleased until final release scope, version and publication are selected.

Copy this checklist into the selected release's work area. Record version,
previous release, exact commit, owner publication authority, evidence paths and
actual outcomes. Unchecked means pending; use N/A only with a reason. A generated
release log is not proof of publication. The project release contract governs.

## Scope and content

- [ ] Select product version from compatibility impact. Keep language/profile,
  process/management schema and client versions distinct.
- [x] Enumerate included work and remaining exclusions; validate selected model,
  tests, review dispositions and links on the actual candidate.
- [ ] Update canonical RELEASE-NOTES.md Unreleased entries with stable work IDs.
  Freeze selected entries into a dated version section; preserve released bytes.
- [ ] Run `sdptool release-log --version X.Y.Z --output Releases/X.Y.Z.md`.
  Check with the same command plus --check. Use --all --output Releases --check
  to check every recorded version. Generator extracts reviewed notes, not Git guesses.

## Installation distribution (SDP publisher)

- [x] Check Template/sdp-root and project-root guidance. Include Sessions guide and
  blank template, never project conversations, boards or existing work records.
- [x] Review SDPTool/profiles/payload.json: every distributed file, destination,
  ownership and removed/moved path. Inventory completeness tests must pass.
- [x] Review SDPTool/profiles/five-phase.json: capabilities, schema/profile versions
  and exact supported predecessor descriptor digests, including current stable.
- [ ] Reconcile SDP.manifest.yaml, Toolkit/SDP-install.manifest.json compatibility
  metadata, package version, release record and notes. Retained legacy metadata
  does not replace the current Go inventory or signed descriptor.
- [ ] Build package and immutable sdp-release.json from a clean exact commit using
  the Go profile builder. Verify file hashes, sourceCommit and executable identity.
  Never edit the built descriptor or use a project-local manifest as authority.
- [ ] Sign descriptor with the trusted publisher key; verify signature, binary hash
  and SHA256SUMS. Store no private keys or secrets in evidence/repository.
- [ ] Rehearse install and each supported upgrade from original verified descriptor
  bytes. Check new Sessions files, existing Session/document preservation, receipt,
  Maintenance entry/history, conflict handling and repeated no-op.
- [ ] Test extracted-source operation, packaged signed child and normal bootstrap.
  Development artifacts or test signers are not production signed evidence.

## Publish and consumer readiness

- [ ] Exact candidate CI and required independent review pass; working tree clean;
  tag absent; merge/publication explicitly authorized.
- [ ] Publish annotated tag and GitHub Release with generated log, signed descriptor,
  executable and checksums; verify downloaded assets against exact release.
- [ ] Reconcile actual tag/commit/time and append truthful release events only after
  success. Do not rewrite a released log or silently correct released notes.
- [ ] Verify gh-sdp selection: its compiled default is pinned. Either publish/update
  the client default or document an explicit SDP_RELEASE selecting this release.
- [ ] Provide preview/apply commands and expected receipt. Owner-requested manual
  project upgrade remains manual; never confuse updating the gh extension with
  upgrading the project's SDP directory.

Consuming projects reuse the scope/content/publication checks; mark SDP publisher
inventory/signing steps N/A if they do not distribute SDP itself. Keep their own
application/package/installer acceptance checks alongside this list.

## Preparation evidence and remaining publication boundary

- [x] Go release-log extraction/check and preservation regressions pass; published
  0.2.0/0.2.1/1.0.0 note sections unchanged. CI checks all generated logs.
- [x] Bounded manual Session format adopted. Only guide/blank template are payloads;
  no conversation content or new management ledger kind is distributed.
- [x] Local-development upgrade rehearsals from all three original signed predecessor
  descriptors pass for fresh/custom content; receipt capability and no-op verified.
- [ ] Freeze actual 1.1.0 notes and generate Releases/1.1.0.md. Align release
  manifests/records and package identity to that candidate, not current 1.0.0 facts.
- [ ] Production-signed clean candidate/CI/extracted-source/download checks. Local
  development rehearsal is not equivalent to these final publication gates.
- [ ] Owner selects merge/publication; tag, assets and actual event reconciliation.
- [ ] Verify new client default or explicit SDP_RELEASE selection before XFMD test.

Read [Manual-upgrade.md](Manual-upgrade.md) for exact preview/apply commands **after**
the release is published. No changes have been applied to XFMD in this Maintenance.
