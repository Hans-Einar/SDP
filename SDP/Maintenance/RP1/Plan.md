# RP1 — Publish the Go installer and adopt XFMD

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0007 |
| project | SDP |
| state | completed |
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

Candidate cbb6a4f was packaged from a clean detached worktree with version 0.2.0.
Its descriptor was signed with the production publisher and verified by SDPTool
without test-key configuration. Clean install applied 63 actions; repeated upgrade
proposed zero. The source archive extracted without .git passes Toolkit validation;
the retained legacy installer applied to an isolated target and recorded
sourceCommit: null. No live project was touched by these tests. The normal Toolkit
unit suite passed 106 tests with 19 environment-dependent installer tests skipped;
those skips are not claimed as coverage. GIP PR39 separately passed all four CI jobs,
including Windows installer and Linux profile execution.

The publisher key ID is
bfc4838a7fbe7e7e3e84f86b9b435e69e6e503f4caf01a92319147e9ab701618.
The private key remains at the owner's private publisher configuration; only the
public anchor is committed. Preserve that key for future releases; loss requires
a newly reviewed trust distribution, not an unsigned fallback.


### RP1-M2 — published

All four PR40 CI jobs passed (contracts, Go installation, Windows legacy installer,
Linux process-profile/recovery). PR40 merged to 738e6c882daed85591311f18248dab2a48ce2076.
That clean exact candidate was packaged, production-signed, tested with a fresh
install/apply/zero-change repeat and tagged v0.2.0. Runtime sources are unchanged
from the tested cbb6a4f candidate. Actual GitHub archive extraction also validated
and rebuilt the identical 59-file payload inventory without Git metadata.

SDP release: https://github.com/Hans-Einar/SDP/releases/tag/v0.2.0

gh-sdp release: https://github.com/Hans-Einar/gh-sdp/releases/tag/v0.1.0

Downloaded SDP assets pass SHA256SUMS. The descriptor digest is
edc0c72101a437c6e12c40a081ef59ae41cf0db32bcaedb73824ec48495aaee5.
Client tag source is 8ad0fc2906fc52bd4ee4c214e801a5872a4b0faa. Independent client
review approved the corrected GitHub CLI asset name gh-sdp-linux-amd64 before
publication. Actual `gh extension install Hans-Einar/gh-sdp` selected v0.1.0 and
`gh sdp --version` downloaded/verified SDPTool 0.2.0 with no test configuration.

### RP1-M3 — live adoption completed

Fresh inspection of the clean xfmd-sdl-navigation worktree at b95a4bb recorded 371
paths. The released client produced the same inspected 138 actions and 194 preserved
paths, with signed provenance and no conflicts. Apply completed as
install-141ecdb474cd970124926c4e, recorded by MAINT-XFMD-0001. Original history bytes
remain an exact prefix and all preserved paths were verified immediately after apply.
All 350 recorded application/build/tooling file hashes remain unchanged.

MAINT-XFMD-0002 reconciles the project-owned board verifier and current path prose,
sets the project name without inventing XFMD product release metadata, and excludes
local recovery journals/backups from Git. Both prior and current board schemas are
understood; management schemas/predecessors/current states are checked. Live validation
passes 18 cards, 77 card events and two management records; 15 lineage negative cases
pass. Copy checks reject broken management predecessors and false completion states.
The complete process adoption is committed locally in XFMD as 41fa494 on its existing
active branch. No XFMD application changes or branch merge/push is included.

Released `gh sdp discover` reports valid. Repeated upgrade before and after local
reconciliation proposes zero changes. The original instructions and operation backups
remain available locally. Fresh manifest/plan and machine evidence are retained in
/home/warloc/.local/state/sdp/adoptions/xfmd-20260927; those inputs target only that root.

**Navigation limitation:** `gh sdp tree` returns KanBan unavailable for the existing
cross-project Ref KB-SDP-014. A successful process exit was insufficient evidence of
a usable tree; inspecting node states exposed this pre-existing local-only resolver
restriction. KB-SDP-034 records the follow-up. No reference was removed or renamed to
hide it. Native XFMD navigation acceptance is not claimed. SDL/SDUI roots remain absent
because XFMD's domain models are not yet registered; KB-XFMD-017 owns that next work.

The authorized merge/release/install/upgrade scope is complete. Published artifacts
remain immutable; this reconciliation commit records real outcomes after publication.

Publication reconciliation passes 106 Toolkit tests (19 environment-dependent
skips), Toolkit/schema validation, management replay and document/frozen-evidence
checks. No executable source changed after the released candidate.
