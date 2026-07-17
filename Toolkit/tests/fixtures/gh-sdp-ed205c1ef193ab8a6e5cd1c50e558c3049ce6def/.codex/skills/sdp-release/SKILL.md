---
skillId: sdp-release
skillVersion: 1.0.0
minimumToolkitVersion: 0.2.0
capabilities:
- sdp.release.prepare
- sdp.release.gate
- sdp.release.reconcile
compatibilityNotes: Initial formal version for SDP Toolkit 0.2.0.
---

# SDP Release

Use this skill to prepare or reconcile a release. It stops before publication
unless explicit human authorization is present.

## Prepare

1. Read manifests, `RELEASE-NOTES.md`, release record, traceability, reviews and
   verification evidence.
2. Run `sdp-versioning` and confirm the target is greater than the previous
   release.
3. Create or confirm a dedicated release-preparation Slice.
4. Run the deterministic gate from `docs/Release-Lifecycle.md`.
5. Freeze the selected Unreleased entries into the target version section.
6. Set state to `prerelease`; leave real tag, release commit and GitHub Release
   fields null until they exist.
7. Append only real preparation, verification and approval events.
8. Prepare exact annotated-tag and GitHub Release commands.
9. Stop before executing publication commands without explicit authorization.

## Reconcile after authorized publication

1. Verify the annotated tag resolves to the approved release-preparation commit.
2. Verify the GitHub Release exists for that tag.
3. Append `release-tag-created` and `release-published` with real identities.
4. Set state to `released` and record the real tag and release commit in a small
   follow-up commit.

Never backdate or pre-claim publication.
