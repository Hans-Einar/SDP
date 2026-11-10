# SDPTOOL-VER-R221 — release preparation evidence

Candidate `c3377121acd24013aab01a33067b9b1d976b8713` implements the same runtime as
owner-selected local producer `34ee535fbf6257b73d2fa0c7fbe4f9d5131aaf5e`.
Changes are release metadata/documents, the supported predecessor inventory and
its matching regression expectation. No runtime implementation changed.
Candidate.json records the exact clean static Linux amd64 v1 package. Earlier
packages remain outside Git as superseded construction evidence; do not publish them.

## Completed local checks

- Complete SDPTool suite passes in evidence/sdptool-tests-final.log. Unchanged
  packages reuse Go's test cache; the corrected profile package ran separately.
  Initial full suite in sdptool-tests.log exposed the stale five-predecessor test;
  adding the original signed 2.2.0 digest corrected it without runtime changes.
- Complete SDL suite, bootstrap race tests, 85 Toolkit tests and SDPTool vet pass.
  Vet output is empty on success; the parent observed exit 0. Independent current
  Toolkit/management validation and all seven release logs are in review-checks.log.
- Actual candidate binary generates all 22 VP01–VP11 overviews across the two
  Ponsse SDL models. Available and empty states remain truthful; this does not
  claim that the models fully document the operational Ponsse system.
- Original signed 0.2.0, 0.2.1, 1.0.0, 2.0.0, 2.1.0 and 2.2.0 descriptors install;
  all upgrade to the explicitly unsigned development candidate, preserve owner
  Session/navigation/SDL files and repeat with zero actions. Original provenance
  comes from the verified descriptor cache. An initial helper using
  --previous-artifact was correctly rejected because that flag means local
  development provenance; no receipt or product was altered to bypass it.
- Actual packaged-child test with a test signer passes. This is supplemental
  mechanism evidence, not production trust. Git-less extracted candidate source
  passes all generated release-log checks.

Evidence files are retained here; full trial plans, replies, descriptors, generated
viewpoints and source archive remain in
/home/warloc/.local/state/sdp-release-2.2.1. No consuming project was upgraded.

## Pending publication gates

Fresh review is recorded in Review.md, with focused follow-ups in the same review
context. CI for candidate c337712 is run 42739852278 and was still running when
this preparation record was written. Earlier 224e667 CI is superseded.
Production signing, signed-target upgrade trials, published download verification,
main integration and final release acceptance remain pending. Draft PR57 retains
all work in the existing checkout; no new worktree was created.

The trusted publisher private key was not located locally or by bounded filename
inspection on Beinerud. Owner suspects an offline laptop. It is a separate base64
Ed25519 signing key, not a normal SSH key. No secret was copied into evidence and
no replacement trust root was generated. Public release fields remain unset.
ProjectGovernance completion, XFMD upgrade and gh-sdp client release are excluded;
explicit SDP_RELEASE selection is documented in the release notes.
