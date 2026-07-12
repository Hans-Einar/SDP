---
skillId: sdp-verifier
skillVersion: 1.0.0
minimumToolkitVersion: 0.2.0
capabilities:
- sdp.verification.validate
- sdp.release.evidence
compatibilityNotes: Initial formal version for SDP Toolkit 0.2.0.
---

# SDP Verifier

Use this skill to validate verification records and release-gate evidence.

1. Read the work contract and its required verification matrix.
2. Inspect each referenced command, environment, artifact and result.
3. Rerun deterministic checks when the environment permits.
4. Confirm evidence applies to the exact commit and release candidate under review.
5. Mark missing, stale, ambiguous or unverifiable evidence as failure, not pass.
6. Check that release-gate evidence covers manifests, notes, traceability,
   compatibility, clean tree and migration requirements.
7. Return pass/fail per requirement plus residual limitations.

Do not infer results from a Worker summary, screenshots without provenance or
successful checks on a different commit.
