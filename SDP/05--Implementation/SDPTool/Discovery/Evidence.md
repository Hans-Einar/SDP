# Session 0003 evidence

## DS1-M1 — source discovery

Candidate: c78831b plus DS1-M1 commit contents; Go 1.27.1, Linux amd64.
GOMAXPROCS=2 and -p 2. Root SDPTool and bootstrap module tests pass. Added real
filesystem tests cover composed Systems, fragments, invalid/future SDL, broken
SDUI, duplicate filenames, source changes, rename/removal, stale registration
inertness, contained paths and source limits. Compiled consumer journey preserves
revision-bound selections and output guards using source-derived model IDs.
Existing semantic service tests explicitly construct Go inventories; they no
longer pretend that a persisted register supplies discovery. New compiled consumer
and filesystem tests own the actual discovery boundary.

Built /tmp/sdptool-session3. Read-only discovery against the actual XFMD checkout
succeeds without registration edits. /tmp/sdp-session3-xfmd-discovery.json contains
the temporary result; summary remains reproducible from the CLI. No XFMD file was
written and no generated navigation sidecar is required.

Independent review, installer removal, final packaging and publication remain
pending. This evidence is not a release-readiness claim.

## DS2-M1 — integration

Installer no longer creates or edits navigation.json; existing project-owned bytes
remain untouched, including malformed historical files. Root registration and its
schema are removed. Contract 0.2, consumer examples, installed source guides and
ecosystem verification now use source-discovered IDs. The SDL model assigns source
inventory discovery to ProjectContext.

Independent review found header-comment rejection and selected-model operations
being blocked by aggregate navigation limits. Both are corrected with regression
tests. Reserved .sdp-backups is excluded with .sdp-operations. Root Go suite and
Toolkit contract validator pass; full release verification follows in DS3.
Client 0.2.0 pins pushed bootstrap 247fb7a; exact final package pairing remains.

## DS3-M1 — verification

Candidate b999352: full SDPTool race tests and vet pass; bootstrap race tests pass;
SDL and SDUI Go suites pass. Toolkit Python suite: 85 tests pass. Toolkit and
project-management validators pass. Ecosystem verify.sh checks all 18 models
with parsing, AST, canonical reparse, source-discovered navigation, revision-bound
selection and static generation. Independent review is approved in Review.md.

Upgrade-rehearsal.json records six disposable upgrades from original signed
0.2.0, 0.2.1 and 1.0.0 descriptors, each with fresh/custom project content.
Target is explicitly unsigned local-development, not production
proof. Sessions, invalid historical navigation.json and ledger prefixes survive;
repeat upgrade has zero actions. Fresh installation creates no navigation.json.
Production-signed repetition and exact clean package/CI remain DS4 obligations.
No live XFMD write occurred.

## DS4-M1 — exact release candidate

Candidate d304261c90066a86b2d8ffcaa2115517ea339b05 is packaged from a clean
isolated worktree. Product code is unchanged from reviewed b999352. Release
metadata and one historical fixture test context were corrected after the first
CI attempt; frozen install-v1 bytes remain unchanged. Python suite again passes
all 85 tests. The source archive gate now describes the actual Go descriptor
contract: explicit verified commit, without requiring retired installer execution.

The final production-signed descriptor is tested with the compiled publisher trust.
Six signed predecessor upgrades (0.2.0/0.2.1/1.0.0, fresh/custom) and clean install
pass, including Session/ledger/inert navigation preservation and zero-action repeat.
The exact packaged signed-child test passes. The no-Git archive validates, builds
a development descriptor with its verified commit, installs successfully and
records that same commit. This archive trial is not a production signature claim.
Final hashes, CI, public downloads and client evidence are reconciled after publication.

Publication.json binds the real annotated tag/release, exact successful CI and
public asset hashes. Downloaded artifacts pass all four checksum checks. The
actual GitHub source archive also passes Toolkit validation without .git.
Production release is published; paired client publication remains a separate
verified operation. Linux amd64 only; native XFMD UI and non-Linux execution
are not claimed by producer and harness tests.

External SPS-007 published gh-sdp 0.2.0 from exact reviewed source72b4e04.
Its immutable bootstrap dependency247fb7a selects SDP2.0.0. The exact final
engine pairing, isolated gh route, fresh-cache public default (production trust),
human/JSON version and source discovery pass. Independent client review approved;
downloaded client bytes match the tested package. Publication.json records both
real release identities and hashes. Global extension and live XFMD remain unchanged.
