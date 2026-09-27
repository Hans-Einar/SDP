# IPD design evidence

Evidence distinguishes parser/model consistency, generated projections and authored
contract review from installer execution. No Go installation or live XFMD upgrade
is performed. The current tool was built from SDL/go at planning baseline db7d82d
with the existing Go 1.27.1 linux/amd64 toolchain into /tmp/ipd-sdl.

## IPD-1-M1 — installation responsibility model

The canonical model adds the project maintainer, four install/upgrade/adoption/
recovery use cases, SDPTool installation units, the external gh-sdp process and
external release artifact service. The latter is a collaborating service boundary,
not a new SDP-owned server. Five new delivery activities are explicitly planned.
Library/package boundaries are not modeled as deployed containers.

Commands: `sdl format`, `sdl check`, `sdl ast` on the canonical SDPTool.design.
The first draft exposed declaration-before-fact and canonical-order requirements;
the formatter normalized it, then check/AST passed with no diagnostics and 90
symbols. A set comparison preserves every statement in the db7d82d model.
The AST is parser output, not a manually constructed JSON model. Detailed channels,
records and scenarios belong to IPD-2; no protocol implementation is inferred.


## IPD-2-M1 — contracts and explicit exchanges

Working candidate: 9607fdc plus IPD-2 diff. Canonical check/AST passes with 188
symbols; every IPD-1 statement is preserved. The extended model contains twelve
channels, closed request/result protocols, five scenarios and four datasets.
The first draft reused field declarations across contracts; SDL correctly rejected
that cardinality. Contract-scoped field names fixed the model without parser changes.
No AST or generated diagram is handwritten. Signatures, JSON validation, path safety
and mutation semantics remain authored design requirements, not parser proofs.

Contract/scenario review records preview-only default, root-bound apply, explicit
manual adoption, signed descriptor design, thin binary delegation, incompatible
legacy journals, forward-only recovery and no fixed packet encoding. The first
implementation milestone must freeze bounded wire schemas and conformance cases.
Actual XFMD read-only status was refreshed to b95a4bb (clean); no files were changed.
Independent reviewer/owner acceptance is not claimed.

## IPD-3-M1 — reproducible generated review and handoff

Candidate: 275d815 plus IPD-3 delivery diff; model SHA-256
c62878ccf713fe1fd51d9a7d14e90dd7f8a70dd5deb346af2d8fb9edae9969d2.
The model did not change after IPD-2. [Generated review](review/index.md) contains
15 bounded selections / 18 diagrams, with 80 files including manifests and source
provenance. Only 16 are Markdown pages (index plus one per selection). No diagram
or generated Markdown was manually edited. [Provenance](review/provenance.json)
records exact model/tool/source/orchestrator identities and output hashes.

Reproduce using a compatible built SDL Go CLI:

```sh
python3 SDP/04--Design/SDPTool/Installation/review.py --sdl /path/to/sdl --check
```

The --check operation generates twice into fresh temporary directories, compares
all bytes, then compares with the retained review. It validates each bundle hash,
source-fact references, model revision and nonempty selection. A different binary
identity is evidence drift even if the model is unchanged; inspect before regenerating
into a new directory. This run used /tmp/ipd-sdl built at planning baseline db7d82d
with Go 1.27.1. The binary's digest and tracked SDL source digest are recorded;
no untracked SDL/go/sourceinput material is inspected.

Model evidence: 188 declarations, including 12 channels, 5 scenarios, 4 datasets
and 3 persistent stores. All pre-IPD statements are preserved. Five new installation
activities remain planned. Negative model variants reject an unpermitted message
(MESSAGE_NOT_PERMITTED) and invalid reply ordinal (CORRELATION_MISMATCH,
INCOMPLETE_SCENARIO). These are meaningful semantic checks, not merely formatting
failures. No signature/path/journal-runtime correctness is inferred from them.

A refreshed legacy reference trial ran:

```sh
SDP_TEST_PWSH=/tmp/sk1-pwsh/pwsh python3 Toolkit/conformance/install-v2/verify_xfmd_snapshot.py /home/warloc/git/xfmd-sdl-navigation --provenance /tmp/ipd-xfmd-reference.json
```

It passed on a disposable copy of XFMD b95a4bbd5ef43c9d7ed27270b58eeff2701a0cb8:
247 copied input files, preserved source hashes/Git status, byte-preserved history
prefix and no-change repetition. The [reference evidence](XFMD-reference-evidence.json)
retains its actual observations. Temporary operation paths are historical evidence,
not live navigation targets. This is PowerShell reference behavior; signed release
selection, the Go engine and wrapper remain unimplemented.

[Scenarios](Scenarios.md) records the author walkthrough for IC01–IC12 and coverage
limits. VP10 is inapplicable to the JSON wire representation; opaque record bytes
still need explicit implementation schemas. Mermaid source is generated; no native
XFMD/visual-rendering or independent-review acceptance is claimed.

[PLAN-SDP-0003](../../../05--Implementation/SDPTool/Installation/Plan.md) is the
planned implementation handoff: read-only plans, apply/recovery, distribution/client
and disposable XFMD adoption with milestone acceptance. It is not started by IPD.
Repository management, Toolkit and document checks plus SDPTool race tests are
run at closeout; final outcomes are recorded below. Product source remains unchanged.

## Final integrated checks

- Fresh `go test -race -count=1 ./...` in SDPTool passes. No production Go source
  changed; this checks the design/board update has not broken current consumers.
- The existing ST1 packaged SDPTool reads the new model through `tree` (1,568
  navigation nodes) and `select` for VP08-InstallationApplied, returning the exact
  model revision and a real generated entry.md. This verifies the producer path,
  not XFMD's native rendering or installed application behavior.
- review.py --check reproduces all 80 retained files with the same binary and
  verifies the two negative semantic cases. Fifteen selections contain 18 diagrams.
- Toolkit validation, management replay, document links and historical-prefix checks
  pass. The pre-existing 105 frozen records/prefixes and 574 generated outputs stay
  unchanged. New IPD output is checked by review.py's own manifests/provenance.
- Authored files pass `git diff --cached --check`. The SDL publisher emits a
  trailing blank line in generated entry.md files; staged checking exposed that
  existing formatting behavior. Generated bytes are retained exactly. Their
  whitespace check excludes only blank-at-eof; all other whitespace checks pass.
  No generator/product edit or hand-edit of generated output was made.
  The unrelated untracked SDL/go/sourceinput remains
  untouched. No live XFMD bytes/status changed in the reference-copy trial.

All IPD design deliverables are complete. This is an author-produced design with
machine evidence and a concrete owner review, not independent approval. The card
moves to gate-review for the design choices and proposed GIP implementation scope;
PLAN-SDP-0003 remains planned. Production trust keys, release numbering, platform
support and live upgrade selection remain later decisions.
