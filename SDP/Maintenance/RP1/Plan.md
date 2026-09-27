# RP1 — Publish the Go installer and adopt XFMD

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0007 |
| project | SDP |
| state | active |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| SourceCards | KB-SDP-033 |

## Outcome and authority

Owner authorization on 2026-09-27 selects merging the delivered PRs, release
publication, installing gh-sdp and upgrading the active XFMD worktree.
PLAN-SDP-0003 is the completed predecessor. Use sdp/release-r1 with milestone
commits; preserve unrelated SDL/go/sourceinput work. RP1-M1 is the bounded
release-preparation work unit in the adopted typed-plan process. gh-sdp follows
its separately installed Slice and independent review lifecycle.

| Phase / milestone | Acceptance |
| --- | --- |
| RP1-M1 — preparation | Merge GIP; production public trust anchor and exact default descriptor; maintainer signing; Linux amd64 candidate and archive checks; frozen release notes and record. |
| RP1-M2 — publication | Annotated SDP v0.2.0 and gh-sdp v0.1.0 tags, observed Releases and verified downloaded binary/descriptor; reconcile only after publication. |
| RP1-M3 — adoption | Install released extension; fresh root-bound XFMD manifest/preview; inspect changes/references; apply via gh sdp; verify preservation, identity, discovery and repeat/no-change; retain recovery evidence. |

## Version and distribution decisions

No prior GitHub Release exists. Keep pending SDP 0.2.0 and first client 0.1.0;
Toolkit 0.1.0 is a migration baseline, not a published release. Five-phase profile
and receipt 3.0 have separate schema identities. Advertise Linux amd64 only.

The old archive-only guidance is superseded for the Go entry point: source
archives cannot run a prebuilt installer without a compiler. Publish exact signed
descriptor and binary assets alongside the archive. Archive extraction must still
resolve shared payloads and permit checks without Git. Descriptor sourceCommit
pins the released source; legacy archive facts keep the null sourceCommit rule.
Preserve legacy compatibility gates and frozen fixtures.

Keep the publisher private key outside repositories in private owner config,
mode 0600. Compile only the public key/digest into consumers. No test trust for
published installation. Default to immutable v0.2.0 descriptor URL.

## Evidence and remaining work

PR39 passed four CI jobs and merged as bdfca5f. XFMD initially clean at b95a4bb;
recheck before adoption. No application source changes selected. Independent
review covers the thin client; engine release checks are author verification.
Broader XFMD SDL modeling is deferred. Append observed progress below.

### RP1-M1 preparation progress

The production public anchor and immutable default are committed in 659e543.
The private key is in the owner's private publisher configuration (not Git).
Go race suites pass, including installation recovery (59.33 seconds), bootstrap
trust checks and maintainer signer rejection tests; vet passes. Production signing
requires a clean exact source checkout and refuses an unknown publisher or an
exposed private key. Package version is explicitly selectable for release builds.
The retained Toolkit artifact is rebuilt from sources after release metadata changes.
