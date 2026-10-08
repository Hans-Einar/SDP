# Sessions discovery — ImplementationPlan

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0015 |
| project | SDP |
| state | completed |
| PlanType | ImplementationPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDPTOOL |

## Outcome, analysis and authority

Owner requests discovery of Sessions for browsing. Existing discovery enumerates
Session files under Files, but exposes no Sessions capability or dedicated root.
KB-SDP-047 and Session 0004 select this bounded extension of ProjectNavigation.
KB042/043 retain broader Session metadata, capture and timeline automation.

Use the existing bounded SDP file scan. Derive inventory.sessions and a sessions
capability from a real SDP/Sessions directory. Add a Sessions navigation tab whose
children reuse the canonical file nodes, preserving directories, filename labels
and typed open targets. Include guides/templates; browsing is not Session validation.
No registry, required naming convention, duplicate lifecycle store or filesystem
watcher. Missing/empty/unreadable/non-directory/symlink paths remain distinguishable;
symlinks are never followed. Hash readable Session Markdown bytes for refresh.
Metadata limit is 1 MiB per document and the shared input budget is 64 MiB.

## Delivery

| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| SN1 | SN1-M1 | discover human/JSON and tree expose Sessions; nested files open; changes refresh; safety/absence tests and current docs aligned | completed |

## Git and verification

Small follow-up on sdp/sdptool/session-navigation, preserving completed Session3.
One milestone commit and push. No merge, release, extension install or XFMD change
is selected. Extend existing contract0.2 additively. Test real CLI JSON/readable
navigation, root/SDP equivalence, nested files, empty/missing/symlink paths,
creation/edit/rename/removal, content revision and no writes. Run SDPTool suite
and project validators. Record actual evidence here before closure.

## Evidence

Candidate: db796a0 plus the SN1-M1 commit contents. Linux amd64, Go 1.27.1,
GOMAXPROCS=2, -p 2. `go -C SDPTool test -p 2 ./...` passes; focused
`go -C SDPTool test -race -p 2 -run TestSession ./...` passes. Toolkit and
ProjectManagement validators and git diff --check pass. SDPTool.design is
formatted and parsed by the owning SDL Go tool.

The real compiled CLI discovers this repository's four Session documents plus
README/template under Sessions. Tests exercise human/JSON discover and tree,
typed open paths with spaces/#, nested nodes, no duplicate/dangling node IDs,
root/area equivalence, no registry/document writes, preserved-mtime content edits,
rename/removal, missing/empty/non-directory/symlink paths and oversized documents.
The consumer fixture adds only the authorized Sessions root; existing roots and
source services remain tested. No independent review or native XFMD acceptance
is claimed for this bounded author-tested change.

Requirements REQ-SDPTOOL-003 and its SDL responsibility, producer contract,
consumer guide and Unreleased notes now describe the extension. KB042/043 remain
backlog for broader automation; KB047 completes this bounded browsing work.
No release has been published for SN1: installed gh-sdp still selects SDP 2.0.0.
