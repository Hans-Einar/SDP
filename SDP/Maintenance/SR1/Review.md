# REV-SDPTOOL-SR1 — Independent preparation review

## Disposition and scope

Approved for release preparation, candidate
`93517ad98cd188c0debeb1d0f3d36d123c6e4a3b`.
Reviewer: fresh read-only `/root/sessions_release_review` context; no implementation
or file edits by that reviewer. This is the verification review within
MAINT-SDP-0012, not a separate project-management CodeReview work package.

## Initial findings and resolution

- P2: unreadable directories received duplicate canonical IDs and an incorrect
  empty Sessions tab. The error callback now updates the existing node, clears
  its target and reports unavailable. Actual permission regressions cover root
  and nested directories; independent compiled CLI probes and race tests pass.
- P2: predecessor regression expected three hashes after adding the fourth original
  2.0.0 descriptor. The expected list now retains all four; test passes.

No remaining material findings.

## Independently inspected evidence

The reviewer verified the clean exact checkout, executable identity and all five
recorded artifact hashes, production Ed25519 signature against the compiled
publisher key, all 64 embedded payload hashes/source bytes, unchanged payload
inventory relative to v2.0.0, and generated release logs. Exact-candidate CI
36842826943 passed both jobs. Eight signed and eight archive predecessor trial
records include preservation and repeated no-op receipts.

Existing roots/node kinds support dynamic viewer tabs. Framework remains 2.0.0;
existing Sessions browsing requires no project SDP upgrade.

Final approval requires truthful reconciliation of evidence, checklist and
preparation references; these are recorded in this closeout. Approval does not
claim native XFMD UI acceptance, publication, main merge, client release or owner
authorization for those actions. See [Evidence](Evidence.md) and [Candidate](Candidate.json).

The review ID follows the existing release-record review schema (REV prefix);
verification retains the system-first SDPTOOL-VER prefix. No schema migration is
included in this preparation.
