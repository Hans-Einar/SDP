# RP3 independent source review

Review date: 2026-09-28. Role: SDP Reviewer, in an independent delegated context.
Authority: owner scope in MAINT-SDP-0010 and TS1; Skills/sdp-reviewer/SKILL.md;
Toolkit/conformance/install-v1/README.md requires independent review of changed
expected outcomes. This review changed no product source.

## Candidate and disposition

Bounded source approval for candidate
56919a19edcbd6e672c8a26661d87d16828be4f8. Initial inspection covered TS1
8296aede7a89164d8d0e89dbc9b759f370c87d37 and RP3 preparation
591920b69534fa12e587f5f95a55b6bce9b13da4 against integration baseline
17c5dfb. Follow-up inspection covered the complete narrow deltas in
6fc9ba84ea9d23cafeec1df2d14fd29eff6ae2b0 and
56919a19edcbd6e672c8a26661d87d16828be4f8.

No unresolved material source finding remains within this scope. This approval
does not establish successful publication, signed upgrade, XFMD navigation or
completion of the outstanding release evidence below.

## Independently inspected evidence

- Compared all 22 files under the preceding Template/sdp-root and
  Template/project-root trees with their archive destinations: every byte is
  preserved. The four neutral current manifest/release/Traceability seeds also
  retain their preceding bytes.
- Compared current Template/sdp-root files with the explicit five-phase
  inventory: no missing source, uninventoried file or duplicate destination.
  Template/profiles is absent. Every Template source in the retained install-v1
  inventory points into Template/legacy/install-v1.
- Compared current Go inventory with baseline: only SDP/SDL/README.md and
  SDP/SDL/AGENTS.md add destinations; no destination removal or ownership change.
  Inspected descriptor construction: project files become initialize-if-missing.
  The guide and release notes accurately describe preservation of local prose.
- Checked all 61 derived artifact payloads against real source bytes and SHA256:
  all match. The authoring relocation does not itself imply unchanged guide
  content: deliberate guide and managed-instruction edits are included.
- Structurally compared every install-v1 expected JSON file with the pre-TS1
  candidate. All 751 changed scalar values are relocated Template source paths
  or current-version values changing from 0.2.1 to 0.2.2. No action semantics,
  preservation expectations or scenario outcomes were weakened. Inspected the
  scenario and test changes separately.
- Checked both declared predecessor hashes against recorded published release
  evidence: edc0c72101a437c6e12c40a081ef59ae41cf0db32bcaedb73824ec48495aaee5
  for 0.2.0 and
  66590e8e967ede6b36d8fa45cdbee1cd80f69505cca698b0e4a9bd960842735a
  for 0.2.1. The follow-up regression assertion checks exactly these two entries
  in the actual generated descriptor; it does not accept arbitrary additions.
- Verified baseline management and Traceability ledger bytes remain prefixes
  of their current files.
- Inspected installed source guidance, current template documentation, release
  compatibility claims, bootstrap selection and descriptor builder. Guidance
  separates System/container ownership and experimental model organization
  from supported parser grammar, navigation registration and runtime behavior.
  Earlier MVP1 repository examples are not presented as new engine grammar.

## Findings resolved before approval

1. The current-version edit initially left the newer-prerelease downgrade
   fixture at 0.2.2-alpha.1, which is older than final 0.2.2. Candidate 591920b
   uses 0.2.3-alpha.1, restoring the intended newer-core rejection scenario.
2. SDPTool/install/Records.md still named the old v0.2.0 compiled default.
   Commit 6fc9ba8 replaces this with an immutable-selector explanation and
   directs readers to bootstrap.DefaultRelease for the candidate version.
3. The existing predecessor test asserted exactly one published digest after
   the approved configuration added the second. Commit 56919a1 updates the
   strict expected inventory to both published predecessors without changing
   implementation behavior.

## Runtime and release evidence boundaries

The coordinator owns the full test and release runs. At report preparation,
the reviewed Go race log showed successful SDPTool and install package runs,
but an overall failure caused by the now-corrected exact-one predecessor test.
The coordinator reports the targeted corrected test passes. This reviewer
attempted a separate targeted run, but plain `go` was unavailable on its shell
PATH; no independent successful Go rerun is claimed here.

Runtime conformance replay and the broader Toolkit checks were still being
collected. No completed-suite or CI success is inferred from partial logs.
Before publication, attach passing current-candidate runtime conformance,
appropriate Go/Toolkit checks and CI; establish clean exact-source signed
packaging and actual predecessor upgrade rehearsal. Package and descriptor
identities must match the merged source used for release. Record actual
publication, client selection, XFMD preservation, repeat no-op and navigation
results before declaring RP3 complete. This review does not substitute for
those outcomes or confer new merge/publication authority.
