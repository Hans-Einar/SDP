# N3j — compatibility with current main

Session0005 T019, 2026-10-18. Owner requests resolving our main conflicts before
continuing feature work. Branch `sdp/xfmd-viewpoint-overviews-main` starts at
`e8760d331933361ab631ef0d52c20e87268a22b3`; it carries only original cc7b52f's
changes plus additive conflict resolution, an Actions compatibility regression
and current evidence. The original branch/commit remains preserved.

## Resolution and verification

- `SDPTool/tree.go`: retain all main fields (including assignment identity,
  readiness/freshness/evidence and blueprint revisions), then add overviewTarget.
- `SDPTool/Contract.md`: retain main's Runnable Programs, Blueprints and Actions
  contracts, then add the overview section. No registry/dispatch schema changes.
- `cd SDL/go && go test ./...`: PASS, complete module including viewpoint,
  documents, source graph, blueprint and generated example workflows.
- `cd SDPTool && go test ./...`: PASS, complete module including blueprint,
  assignment, Actions, governance, install and presentation tests; no exclusions.
- New `TestInvokedDiscoveryPreservesOverviewTargets`: PASS. Actual JSON invocation
  and direct discovery payloads agree; returned VP01/VP07 logical/mode targets
  execute unchanged while original detail targets remain distinct. This test was
  added after the full suite started and was run separately with `-count=1`.
- `cd SDPTool && go vet ./...`: PASS.
- `go build -o /tmp/xfmd-n3j-sdptool ./cmd/sdptool`: PASS.
- Native XFMD SdpOverviewGuiTest against this binary: PASS. Expansion, explicit
  Enter/pointer selection, both modes and PDF export exercised. VP01 and VP07
  native screenshots inspected; PDF extraction retains all expected fixture labels.
- Actual XFMD model VP01/VP07 logical/NormalDesktop diagram bytes match N3i.
- ProjectManagement validator: PASS, 66 cards, 45 management records, 604 events.
  Main's card/ledger files were preserved exactly; no historical event import.

`N3j-candidate.json` records tested code/dependency, binary and native-test hashes.
Adjacent N3j logs contain exact suite output. Detailed preflight, native captures,
actual-model output and post-commit merge forecast are retained in:

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N3j-001/README.md

The final integration commit is this document's introducing commit. A clean
forecast applies only to the recorded fetched main and this committed candidate.
It does not cover future changes, unpublished remote work or the dirty shared
checkout. No main merge, push, release, install, independent review or broad
KB032 completion is claimed. See Integration.md for the selected integration path.
