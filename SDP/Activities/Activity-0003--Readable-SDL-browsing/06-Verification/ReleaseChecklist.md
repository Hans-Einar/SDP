# ReleaseChecklist — XFMD 0.10

Instantiated from Toolkit/docs/ReleaseChecklist.md and scoped to the collaborating
XFMD application RPM. The XFMD application-specific checklist is
`SDP/Maintenance/RPM/0.10/ReleaseChecklist.md` in its repository.

## Scope and content

- [x] Owner selected XFMD0.10 (minor public product version); RPM Version0.10.0.
  No SDL grammar, Toolkit release or management-schema version is inferred.
- [x] Scope: Activity2 management plus Activity3 producer-owned readable browsing,
  generic native consumption and a mandatory matching private SDPTool.
- [x] Independent product qualification: final producerc3463b7, native browsing
  d530aa7; final packaging8546e95 independently approved after its local exact-byte gate.
- [x] Canonical SDP RELEASE-NOTES.md Unreleased describes selected work. This is
  not a public Toolkit version freeze; XFMD package notes live in its release area.
- N/A: `sdptool release-log` Toolkit-version extraction. This delivery uses XFMD's
  application RPM release procedure, not a versioned SDP Toolkit release.

## Exact candidate and application installation

Private producer source: `c3463b78d53ebf984b0f553ebe93a556fb3ac6a8`.
Public predecessor: `xfmd-0.9.0-71.el10_2.x86_64.rpm`, source8893a32.
Target: AlmaLinux10.2 x86_64. Native source: `8546e956d1a941ece9195d335b7196abea905289`.
Candidate: `xfmd-0.10.0-76.el10_2.x86_64.rpm`.
Unsigned SHA256: `e7e39924edaad957355e40210f1b8a7e3989ca1d920dbf107e556588c90af563`.
[B5](B5/) retains exact source/build/probe provenance and local gate logs.

- [x] Exact committed native candidate, externally pinned producer archive and
  binary identities; source/module licenses and provenance.
- [x] Target-container build and complete applicable native suite; environment
  exclusions explicitly reconciled with prior native evidence.
- [x] Installed private helper selected by default; actual Activity/reading/source/
  history workflows tested against installed producer and native runtime.
- [x] Exact signed public0.9 predecessor install and candidate upgrade, plus fresh
  install/reinstall/removal and unrelated-file preservation.
- [x] Independent final package qualification before publication.

## Publication and consumer readiness

- [x] Owner explicitly authorizes feature-branch RPM publication to beinerud; main
  merge is not authorized and is not needed for this application release.
- [x] Publish exact verified payload, sign package/repository on the server and
  verify downloaded payload identity and trusted-key signatures.
- [x] Public DNF upgrade from signed0.9 to exact0.10 with both package and metadata
  signature checks, then a fresh public install with exact identity/workflow probes.
- [x] Retain actual times, source hashes, package hashes, signatures and results;
  reconcile Session, plans/card and Activity only after the required outcome.
- [x] Provide the user `sudo dnf upgrade --refresh xfmd` after public verification.

## Inapplicable inherited Toolkit gates

Template/project seeds, payload inventory/five-phase descriptors, SDP manifests,
`Go profile builder`/`sdp-release.json`, signed Toolkit descriptor, framework
predecessor upgrades/receipts, gh-sdp default and a GitHub Toolkit tag/Release are
N/A: this package ships an application and private helper, not a public Toolkit
installation release. Their status must not be described as passed. The applicable
RPM identity, signature, install/upgrade and publication checks above replace none
of those Toolkit claims; they establish only the selected application delivery.

## Current result

The first bdd0602 candidate passed88/89 container native checks. The remaining
StandaloneMermaid test raced its200ms file filter; the bounded readiness correction
preserves the final assertion. Successor8546e95 passed89/89 (122.79s), including
real SdpReadingGuiTest. Relative to the prior90-test host suite, this environment
omits FullscreenTest (wmaker) and ThemePopupGuiTest (Xtst), and adds the actual
producer-driven reading test. Earlier unchanged host evidence remains applicable.

The signed predecessor upgrade, installed default helper selection, all11 reading
landings, rendering/source/Back and real Planned Activity workflows passed. Reinstall,
remove/unrelated-file preservation and fresh install also passed. Independent
review reran the read-only exact-byte publication admission and approved this
candidate. Public signature/download/upgrade/install has now passed; see B5/published.json,
B5/signature.txt and B5/public-dnf.log. Independent final release review approves.

## Published result

Published and verified: [xfmd-0.10.0-76.el10_2.x86_64.rpm](https://beinerud.ddns.net/repo/AlmaLinux-10.2-apps/stable/x86_64/Packages/xfmd-0.10.0-76.el10_2.x86_64.rpm).
Signed SHA256: `8c5f2232d87c043bbae16b49e9fbe3815237d69fdb858efb7adb2d51f3587f2d`.
Trusted key ID901af029; RSA/SHA256 signature and payload digest pass. Signing
changes the package header; the payload matches the reviewed unsigned candidate.
Public DNF forced both gpgcheck=1 and repo_gpgcheck=1, upgraded the signed0.9.0-71
predecessor to0.10.0-76, then removed/reinstalled from the public repository.
Both outcomes passed exact app/private-producer identities and native workflows.
The user upgrade command is `sudo dnf upgrade --refresh xfmd`.
