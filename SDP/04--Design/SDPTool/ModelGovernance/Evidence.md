# MG2-M1 — design and SDL validation evidence

2026-10-03. This is design evidence, not a ModelGovernance implementation test.

## Candidate and method

Canonical SDL entry: SDP/SDL/ProjectGovernance/SDPTool/ModelGovernance.design.
Profile design-core/0.6; four reachable source files. Checked revision:
`ddeddf6050239ebbde41bf0a74d149bd3a98948983e4c83e16d8fdae81f2d707`.
Exact source hashes are recorded by generated review/sources.json and
[check.json](check.json). Toolchain: local Go 1.27.1 under
/home/warloc/.local/share/sdp-toolchains/go1.27.1/go/bin/go (not on shell PATH).

From SDL/go:

```sh
go run ./cmd/sdl check ../../SDP/SDL/ProjectGovernance/SDPTool/ModelGovernance.design
go run ./cmd/sdl viewpoints ../../SDP/SDL/ProjectGovernance/SDPTool/ModelGovernance.design --format static --viewpoint VP01,VP02,VP08 --output ../../SDP/04--Design/SDPTool/ModelGovernance/review --project model-governance
```

The frontend's format --file-map output was applied to all four files before the
successful final check. Initial drafts used an unsupported relation and an incomplete
dataset declaration; those were corrected before recording this candidate. The final
check is valid with no warnings. Generated static output: 9 diagrams, 38 files.
No hand-authored diagrams replace generated SDL facts. Mermaid has not been visually
verified in XFMD; SVG rendering was not requested for this milestone.

## Scope proven and not proven

The Go frontend accepts the composed source and validates its declared relations,
closed message contracts and two scenarios (accepted checkpoint and refused stale
request). Generated viewpoints retain source provenance. Activities explicitly
remain planned. Text-valued expected-head/input-digest fields model a conceptual
API; the parser does not enforce filesystem transactions or SHA domains.

Contract.md specifies proposed command grammar, local record representation,
rollback/merge/publication behavior and consumer boundaries. MG3 must test baseline
reconstruction, deletes, overlap, crash recovery and promotion payload stripping.
No runtime operation or complete machine schema is tested by these checks.

Management/Toolkit validators and git diff --check are run at delivery. Concurrent
ProjectGovernance changes are outside this milestone and must not enter its commit.
