# N3i — VP01/VP07 integration handoff

Prepared 2026-10-18 under XFMD Session0005 T016 / PLAN-XFMD-0009 / KB-SDP-032.

## Candidate and scope

Branch: `sdp/xfmd-viewpoint-overviews`, based on `48c919f`. This milestone
contains the previously uncommitted N3d/N3e overview implementation, a current
contract/requirement clarification and composed-source freshness coverage.
Use this file's Git commit as the exact integration unit. It does not include
the concurrent SDUI, widget, sourceinput, MCP or governance work.

Product files: SDL/go/viewpoint/{query,overview,overview_test,feature_overview,
feature_overview_test}.go; SDL/go/documents/selection.go; SDPTool/{tree,
overview_test}.go; SDPTool/Contract.md. Requirements and this work-unit folder
record the selected boundary and evidence. The contract contains only the final
overview behavior; retain the separate SDUI combined-preview section when
reconciling with the shared dirty checkout.

## Verification

- `cd SDL/go && go test ./viewpoint ./documents`: PASS.
- `cd SDPTool && go test .`: PASS, including all existing root-package tests and
  composed-model overview/source-provenance/freshness regression; no exclusions.
- `cd SDPTool && go build -o /tmp/xfmd-n3i-sdptool ./cmd/sdptool`: PASS.
- Native XFMD SdpOverviewGuiTest against that exact binary: PASS. Expand causes
  no generation; Enter/label selection and both mode targets generate correctly.
  Native screenshots were inspected; exported PDF text retains the expected labels.
- Actual XFMD model discovery/select: VP01 12 nodes/11 edges; VP07 logical and
  NormalDesktop each 23 nodes/19 edges, with respectively 113 and 137 proof facts.
- ProjectManagement validator: PASS (64 cards, 43 management records, 551 events);
  lineage examples and git diff --check pass on the isolated candidate.

Binary/source hashes are in candidate.json; detailed logs and captured bundles:
file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDL-Navigation/evidence/N3i-001/README.md

## Integrating with concurrent work

The original producer working files are intentionally retained in SDP-vNow.
Compare the scoped commit with those dirty copies before applying it; most
product hunks are already present there. The composed-model test and consolidated
contract are new in this branch. Do not import all dirty producer files, reset
the working tree, or apply duplicate overview implementations. Reconcile shared
changes before committing/merging on the integration branch.

KB032's active card and earlier lifecycle events are currently uncommitted in
the shared checkout. The isolated branch keeps its committed board baseline.
Copying the older KB032 events after later committed IDs fails the namespace's
monotonic rule; no old ledger bytes were reordered or rewritten. This code commit
therefore deliberately excludes board/ledger fragments. Preserve that complete
shared history when its owner integrates it; never synthesize a second transition.

## Remaining boundaries

This is a local integration candidate, not main integration, independent review,
a published release or an installed gh upgrade. Installed SDPTool 2.2.0 remains
unchanged. Verify the exact merged/released binary advertises overviewTarget and
executes its unchanged URI before claiming consumer availability.

KB032 stays in progress for focused use-case/scenario projections, conditions,
exact source navigation, readable overview publication and other viewpoints.
KB030/PLAN-SDP-0024 JSON action discovery remains a separate held workstream.
